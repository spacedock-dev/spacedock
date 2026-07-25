// ABOUTME: Closed local Git-object coordinates for gate-selected source files.
// ABOUTME: Resolution never fetches, falls back to worktree bytes, or retains refs.
package gitsource

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var objectIDRE = regexp.MustCompile(`^[0-9a-f]{40,}$`)
var digestRE = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

type Roots struct {
	Main  string
	State string
}

type Source struct {
	URI string
	Rev string
}

// Inspect freezes the selected worktree file as a logical root, full HEAD
// commit, repository-relative path, and independent raw-byte SHA-256.
func Inspect(roots Roots, selectedPath string) (Source, error) {
	path, err := filepath.Abs(selectedPath)
	if err != nil {
		return Source{}, fmt.Errorf("resolve selected source: %w", err)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return Source{}, fmt.Errorf("selected source must be a readable non-symlink regular file: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return Source{}, fmt.Errorf("selected source must be a readable non-symlink regular file")
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return Source{}, fmt.Errorf("resolve selected source: %w", err)
	}
	worktree, rootName, err := classify(roots, path)
	if err != nil {
		return Source{}, err
	}
	repoPath, err := filepath.Rel(worktree, path)
	if err != nil || repoPath == "." || repoPath == ".." || strings.HasPrefix(repoPath, ".."+string(filepath.Separator)) {
		return Source{}, fmt.Errorf("selected source is outside its workflow Git root")
	}
	repoPath = filepath.ToSlash(repoPath)
	if _, err := git(worktree, "ls-files", "--error-unmatch", "--", ":(literal)"+repoPath); err != nil {
		return Source{}, fmt.Errorf("selected source is not the exact committed file; commit the exact source before preparation")
	}
	commitBytes, err := git(worktree, "rev-parse", "HEAD")
	if err != nil {
		return Source{}, fmt.Errorf("read selected source commit: %w", err)
	}
	commit := strings.TrimSpace(string(commitBytes))
	if !objectIDRE.MatchString(commit) {
		return Source{}, fmt.Errorf("selected source repository returned a non-full commit id")
	}
	object, err := readObject(worktree, commit, repoPath)
	if err != nil {
		return Source{}, fmt.Errorf("read selected source from local Git object: %w", err)
	}
	current, err := os.ReadFile(path)
	if err != nil {
		return Source{}, fmt.Errorf("read selected source: %w", err)
	}
	if !bytes.Equal(object, current) {
		return Source{}, fmt.Errorf("selected source differs from its committed Git object; commit the exact source before preparation")
	}
	escaped, err := escapeRepositoryPath(repoPath)
	if err != nil {
		return Source{}, err
	}
	return Source{
		URI: "git-root://" + rootName + "/" + commit + "/" + escaped,
		Rev: RawDigest(object),
	}, nil
}

// Resolve reopens only the URI's immutable local object and verifies its raw pin.
func Resolve(roots Roots, locator, rev string) ([]byte, error) {
	rootName, commit, repoPath, err := parse(locator)
	if err != nil {
		return nil, err
	}
	var worktree string
	switch rootName {
	case "main":
		worktree = roots.Main
	case "state":
		worktree = roots.State
	default:
		return nil, fmt.Errorf("unknown git-root logical root %q", rootName)
	}
	if strings.TrimSpace(worktree) == "" {
		return nil, fmt.Errorf("git-root logical root %q is unavailable", rootName)
	}
	worktree, err = filepath.Abs(worktree)
	if err != nil {
		return nil, fmt.Errorf("resolve git-root %s: %w", rootName, err)
	}
	resolved, err := git(worktree, "rev-parse", "--verify", commit+"^{commit}")
	if err != nil || strings.TrimSpace(string(resolved)) != commit {
		return nil, fmt.Errorf("git-root commit is not an exact full local commit object")
	}
	body, err := readObject(worktree, commit, repoPath)
	if err != nil {
		return nil, fmt.Errorf("read git-root local object: %w", err)
	}
	if !digestRE.MatchString(rev) || RawDigest(body) != rev {
		return nil, fmt.Errorf("git-root object does not match its raw SHA-256 revision")
	}
	return body, nil
}

// SameLogicalRevision reports whether two immutable coordinates name the same
// logical root, repository path, and raw bytes. Commit identity may differ when
// an unrelated state commit advances the selected file's repository.
func SameLogicalRevision(left, right Source) (bool, error) {
	leftRoot, _, leftPath, err := parse(left.URI)
	if err != nil {
		return false, err
	}
	rightRoot, _, rightPath, err := parse(right.URI)
	if err != nil {
		return false, err
	}
	return leftRoot == rightRoot && leftPath == rightPath && left.Rev == right.Rev, nil
}

func RawDigest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func classify(roots Roots, path string) (worktree, name string, err error) {
	selected, err := inspectRepository(filepath.Dir(path))
	if err != nil {
		return "", "", fmt.Errorf("resolve selected source Git root: %w", err)
	}
	mainRoot, err := inspectRepository(roots.Main)
	if err != nil {
		return "", "", fmt.Errorf("resolve main workflow Git root: %w", err)
	}
	var stateRoot repository
	if roots.State != "" {
		stateRoot, err = inspectRepository(roots.State)
		if err != nil {
			return "", "", fmt.Errorf("resolve state workflow Git root: %w", err)
		}
	}
	// Compare Git common directories rather than checkout paths. This classifies a
	// selected file from a linked implementation worktree as main while retaining
	// that selected worktree's own HEAD and repository-relative path.
	if stateRoot.Common != "" && stateRoot.Top != mainRoot.Top &&
		(selected.Top == stateRoot.Top || stateRoot.Common != mainRoot.Common && selected.Common == stateRoot.Common) {
		return selected.Top, "state", nil
	}
	if selected.Common == mainRoot.Common {
		return selected.Top, "main", nil
	}
	return "", "", fmt.Errorf("selected source is not owned by a workflow Git root")
}

type repository struct {
	Top    string
	Common string
}

func inspectRepository(root string) (repository, error) {
	if strings.TrimSpace(root) == "" {
		return repository{}, fmt.Errorf("root is empty")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return repository{}, err
	}
	top, err := git(abs, "rev-parse", "--show-toplevel")
	if err != nil {
		return repository{}, err
	}
	topPath := filepath.Clean(strings.TrimSpace(string(top)))
	if resolved, resolveErr := filepath.EvalSymlinks(topPath); resolveErr == nil {
		topPath = resolved
	}
	common, err := git(abs, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return repository{}, err
	}
	commonPath := filepath.Clean(strings.TrimSpace(string(common)))
	if resolved, resolveErr := filepath.EvalSymlinks(commonPath); resolveErr == nil {
		commonPath = resolved
	}
	return repository{Top: topPath, Common: commonPath}, nil
}

func readObject(root, commit, repoPath string) ([]byte, error) {
	return git(root, "cat-file", "blob", commit+":"+repoPath)
}

func git(root string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

func parse(locator string) (root, commit, repoPath string, err error) {
	u, err := url.Parse(locator)
	if err != nil || u.Scheme != "git-root" || u.Host == "" || u.Opaque != "" ||
		u.User != nil || u.Port() != "" || u.RawQuery != "" || u.Fragment != "" {
		return "", "", "", fmt.Errorf("invalid git-root locator")
	}
	if u.Host != "main" && u.Host != "state" {
		return "", "", "", fmt.Errorf("unknown git-root logical root %q", u.Host)
	}
	rawPath := u.EscapedPath()
	if !strings.HasPrefix(rawPath, "/") || strings.Contains(rawPath, "\\") {
		return "", "", "", fmt.Errorf("invalid git-root path")
	}
	parts := strings.Split(strings.TrimPrefix(rawPath, "/"), "/")
	if len(parts) < 2 || !objectIDRE.MatchString(parts[0]) {
		return "", "", "", fmt.Errorf("git-root locator requires a full commit id and repository path")
	}
	decoded := make([]string, 0, len(parts)-1)
	for _, raw := range parts[1:] {
		if raw == "" {
			return "", "", "", fmt.Errorf("git-root repository path has an empty segment")
		}
		segment, err := url.PathUnescape(raw)
		if err != nil || segment == "" || segment == "." || segment == ".." || strings.ContainsAny(segment, `/\`) {
			return "", "", "", fmt.Errorf("invalid git-root repository path segment")
		}
		if url.PathEscape(segment) != raw {
			return "", "", "", fmt.Errorf("git-root repository path is not canonically escaped")
		}
		decoded = append(decoded, segment)
	}
	return u.Host, parts[0], strings.Join(decoded, "/"), nil
}

func escapeRepositoryPath(repoPath string) (string, error) {
	parts := strings.Split(repoPath, "/")
	escaped := make([]string, 0, len(parts))
	for _, segment := range parts {
		if segment == "" || segment == "." || segment == ".." || strings.ContainsAny(segment, `/\`) {
			return "", fmt.Errorf("invalid selected repository path")
		}
		escaped = append(escaped, url.PathEscape(segment))
	}
	return strings.Join(escaped, "/"), nil
}
