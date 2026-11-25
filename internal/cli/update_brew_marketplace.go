package cli

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var stableRepairVersion = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
var repairGitHubRepo = regexp.MustCompile(`^[a-zA-Z0-9_-]+/[a-zA-Z0-9_.-]+$`)
var repairCommit = regexp.MustCompile(`^[0-9a-f]{40}$`)

type repairEntry struct {
	Version string
	Local   string
	Repo    string
	Path    string
	Ref     string
}

func resolveRepairEntry(data []byte, name string) (repairEntry, error) {
	var listing struct {
		Name    string
		Plugins []struct {
			Name, Version string
			Source        json.RawMessage
		}
	}
	if err := json.Unmarshal(data, &listing); err != nil || listing.Name != extensionMarketplace {
		return repairEntry{}, fmt.Errorf("could not read the %s marketplace listing", extensionMarketplace)
	}
	var entries []repairEntry
	for _, plugin := range listing.Plugins {
		if plugin.Name != name {
			continue
		}
		entry := repairEntry{Version: plugin.Version}
		if !stableRepairVersion.MatchString(entry.Version) {
			return repairEntry{}, fmt.Errorf("the marketplace does not list a stable X.Y.Z release of %s", name)
		}
		if err := json.Unmarshal(plugin.Source, &entry.Local); err == nil && strings.HasPrefix(entry.Local, "./") {
			entry.Local = strings.TrimPrefix(entry.Local, "./")
			if !safeRepairSourcePath(entry.Local) || entry.Local == "." {
				return repairEntry{}, fmt.Errorf("the marketplace lists an invalid local path for %s", name)
			}
		} else {
			var source struct{ Source, URL, Path, Ref string }
			if err := json.Unmarshal(plugin.Source, &source); err != nil ||
				(source.Source != "git-subdir" && source.Source != "url") ||
				!strings.HasPrefix(source.URL, "https://github.com/") {
				return repairEntry{}, fmt.Errorf("the marketplace source for %s cannot be verified", name)
			}
			entry.Repo = strings.TrimSuffix(strings.TrimPrefix(source.URL, "https://github.com/"), ".git")
			entry.Ref, entry.Path = source.Ref, source.Path
			if source.Source == "url" {
				if entry.Path != "" {
					return repairEntry{}, fmt.Errorf("the whole-repository source for %s declares a subdirectory", name)
				}
				entry.Path = "."
			}
			if !repairGitHubRepo.MatchString(entry.Repo) || entry.Ref == "" || strings.HasPrefix(entry.Ref, "-") ||
				strings.ContainsAny(entry.Ref, "*?[: \\~^") || !safeRepairSourcePath(entry.Path) {
				return repairEntry{}, fmt.Errorf("the marketplace lists an invalid repository, ref or path for %s", name)
			}
		}
		entries = append(entries, entry)
	}
	if len(entries) != 1 {
		return repairEntry{}, fmt.Errorf("the marketplace must list exactly one %s release (found %d)", name, len(entries))
	}
	return entries[0], nil
}

func safeRepairSourcePath(name string) bool {
	return name != "" && !path.IsAbs(name) && path.Clean(name) == name &&
		name != ".." && !strings.HasPrefix(name, "../") && !strings.ContainsAny(name, "\\\x00")
}

func (r brewRepair) releaseSource() string {
	if channelMarketplace(devBranch) == extensionMarketplace {
		return channelMarketplaceSource(devBranch)
	}
	return extensionMarketplaceSource()
}

func (r brewRepair) readReleaseListing(source string) ([]byte, error) {
	if info, err := os.Stat(source); err == nil && info.IsDir() {
		return readRepairFile(source, ".claude-plugin/marketplace.json")
	}
	if !repairGitHubRepo.MatchString(source) {
		return nil, fmt.Errorf("the stable marketplace source must be a directory or a GitHub owner/repo")
	}
	return r.rawGitHub(source, "HEAD", ".claude-plugin/marketplace.json")
}

func readRepairFile(root, file string) ([]byte, error) {
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	filename, err := filepath.EvalSymlinks(filepath.Join(root, file))
	if err != nil {
		return nil, err
	}
	relative, err := filepath.Rel(root, filename)
	if err != nil || !safeRepairSourcePath(filepath.ToSlash(relative)) {
		return nil, fmt.Errorf("marketplace file leaves its registered directory")
	}
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 1<<20 {
		return nil, fmt.Errorf("marketplace file is not a supported JSON file")
	}
	return os.ReadFile(filename)
}

func (r brewRepair) verifyRepairManifest(entry repairEntry, root, host, name string) (string, error) {
	manifest := manifestSubpath(host)
	read := func(file string) ([]byte, error) {
		return readRepairFile(root, path.Join(entry.Local, file))
	}
	var commit string
	if entry.Local == "" {
		var err error
		commit, err = r.resolveRepairRef(entry.Repo, entry.Ref, false)
		if err != nil {
			return "", err
		}
		read = func(file string) ([]byte, error) {
			return r.rawGitHub(entry.Repo, commit, path.Join(entry.Path, file))
		}
	}
	data, err := read(manifest)
	if host == "codex" && os.IsNotExist(err) {
		data, err = read(manifestSubpath("claude"))
	}
	if err != nil {
		return "", fmt.Errorf("could not read the %s manifest for %s: %w", host, name, err)
	}
	var installed struct{ Name, Version string }
	if err := json.Unmarshal(data, &installed); err != nil || installed.Name != name || installed.Version != entry.Version {
		return "", fmt.Errorf("the %s manifest does not identify %s %s", host, name, entry.Version)
	}
	return fmt.Sprintf("%s:%x", commit, sha256.Sum256(data)), nil
}

func (r brewRepair) resolveRepairRef(repo, ref string, tagOnly bool) (string, error) {
	if !tagOnly && repairCommit.MatchString(ref) {
		return ref, nil
	}
	remote := "https://github.com/" + repo + ".git"
	out, err := r.git("", "ls-remote", "--exit-code", remote, "refs/heads/"+ref, "refs/tags/"+ref, "refs/tags/"+ref+"^{}")
	if err != nil {
		return "", fmt.Errorf("could not resolve %s at %s: %w", repo, ref, err)
	}
	refs := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		hash, name, ok := strings.Cut(line, "\t")
		if !ok || !repairCommit.MatchString(hash) {
			return "", fmt.Errorf("Git returned an unrecognized ref for %s", repo)
		}
		refs[name] = hash
	}
	head, tag := refs["refs/heads/"+ref], refs["refs/tags/"+ref]
	switch {
	case !tagOnly && head != "" && tag == "":
		return head, nil
	case head == "" && tag != "":
		if peeled := refs["refs/tags/"+ref+"^{}"]; peeled != "" {
			return peeled, nil
		}
		return tag, nil
	default:
		return "", fmt.Errorf("Git ref %q for %s is ambiguous or has the wrong ref type", ref, repo)
	}
}
