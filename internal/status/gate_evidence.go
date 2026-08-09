// ABOUTME: Read-only committed gate-evidence bundle for one current gate entity.
// ABOUTME: Git pins replace agent-side path search while leaving evidence judgment to the FO.
package status

import (
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/gitsource"
)

type gateEvidenceBlob struct{ root, path, inputPath, objectID, digest, bytes string }

func runReadGateEvidence(roots roots, ref string, asJSON bool, stdout, stderr io.Writer) int {
	if !asJSON {
		return errExit(stderr, "--gate-evidence requires --json")
	}
	if strings.ContainsAny(ref, `/\\`) || filepath.IsAbs(ref) || isRegularFile(ref) {
		return errExit(stderr, "--gate-evidence requires a task reference, not a path")
	}
	entityPath, err := ResolveActivePath(roots.definitionDir, "", ref, stderr)
	if err != nil {
		return errExit(stderr, "invalid gate task reference: "+err.Error())
	}
	if realpathOf(FindGitRoot(roots.definitionDir)) == realpathOf(FindGitRoot(roots.entityDir)) {
		return errExit(stderr, "split workflow and entity state must use independent Git roots")
	}
	pins := gitsource.Roots{Main: roots.definitionDir, State: roots.entityDir}
	entity, err := pinGateEvidence(pins, entityPath)
	if err != nil {
		return errExit(stderr, err.Error())
	}
	stageName := parseFrontmatterContent([]byte(entity.bytes))["status"]
	readmePath := filepath.Join(roots.definitionDir, "README.md")
	readme, err := pinGateEvidence(pins, readmePath)
	if err != nil {
		return errExit(stderr, err.Error())
	}
	var gate bool
	stages, _ := ParseStagesWithDefaultsData([]byte(readme.bytes))
	for _, stage := range stages {
		if stage.Name == stageName {
			gate = stage.gate
			break
		}
	}
	if !gate {
		return errExit(stderr, fmt.Sprintf("current stage %q is not a gate stage", stageName))
	}
	spans, err := FindSectionSpans([]byte(readme.bytes), []string{stageName})
	if err != nil {
		return errExit(stderr, "resolve current gate stage prose: "+err.Error())
	}

	mainPaths, err := gateEvidencePaths(roots.definitionDir, "", "", true)
	if err != nil {
		return errExit(stderr, err.Error())
	}
	entityRel, err := filepath.Rel(roots.entityDir, entityPath)
	if err != nil || entityRel == ".." || strings.HasPrefix(entityRel, ".."+string(filepath.Separator)) {
		return errExit(stderr, "resolved entity is outside the entity root")
	}
	scope := filepath.Dir(entityRel)
	if scope == "." {
		scope = strings.TrimSuffix(entityRel, filepath.Ext(entityRel))
	}
	entityRel = filepath.ToSlash(entityRel)
	statePaths, err := gateEvidencePaths(roots.entityDir, filepath.ToSlash(scope)+"/", entityRel, false)
	if err != nil {
		return errExit(stderr, err.Error())
	}
	if len(mainPaths) == 0 || len(statePaths) == 0 {
		return errExit(stderr, "no committed Markdown candidates in both workflow and entity roots")
	}

	seen := map[string]string{}
	candidates := make([]gateEvidenceBlob, 0, len(mainPaths)+len(statePaths))
	for _, item := range append(mainPaths, statePaths...) {
		name := strings.ToLower(filepath.Base(item.path))
		if prior := seen[name]; prior != "" {
			return errExit(stderr, fmt.Sprintf("duplicate candidate name %q: %s and %s", name, prior, item.path))
		}
		seen[name] = item.path
		blob, pinErr := pinGateEvidence(pins, item.abs)
		if pinErr != nil {
			return errExit(stderr, pinErr.Error())
		}
		candidates = append(candidates, blob)
	}
	stage := newJSONObj().set("name", stageName).set("root", readme.root).set("path", readme.path).
		set("object_id", readme.objectID).set("digest", readme.digest).set("bytes", readme.bytes[spans[0].Start:spans[0].End])
	doc := newJSONObj().set("command", "gate-evidence").setValue("entity", gateEvidenceJSON(entity)).setValue("stage", stage)
	arr := make(jsonArr, 0, len(candidates))
	for _, blob := range candidates {
		arr = append(arr, gateEvidenceJSON(blob))
	}
	doc.setValue("candidates", arr)
	emitJSON(stdout, doc)
	return 0
}

type gateEvidencePath struct{ path, abs string }

func gateEvidencePaths(root, prefix, excluded string, topLevel bool) ([]gateEvidencePath, error) {
	out, err := gateEvidenceGitOutput(root, "ls-files", "-z", "--", "*.md")
	if err != nil {
		return nil, fmt.Errorf("list committed Markdown candidates: %w", err)
	}
	var paths []gateEvidencePath
	for _, path := range strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00") {
		path = filepath.ToSlash(path)
		if path == "" || path == "README.md" || path == excluded {
			continue
		}
		if topLevel && strings.Contains(path, "/") {
			continue
		}
		if !topLevel && !strings.HasPrefix(path, prefix) {
			continue
		}
		paths = append(paths, gateEvidencePath{path, filepath.Join(root, filepath.FromSlash(path))})
	}
	return paths, nil
}

func pinGateEvidence(roots gitsource.Roots, path string) (gateEvidenceBlob, error) {
	source, err := gitsource.Inspect(roots, path)
	if err != nil {
		return gateEvidenceBlob{}, err
	}
	body, err := gitsource.Resolve(roots, source.URI, source.Rev)
	if err != nil {
		return gateEvidenceBlob{}, err
	}
	u, err := url.Parse(source.URI)
	if err != nil {
		return gateEvidenceBlob{}, fmt.Errorf("parse committed source identity: %w", err)
	}
	parts := strings.Split(strings.TrimPrefix(u.EscapedPath(), "/"), "/")
	if len(parts) < 2 {
		return gateEvidenceBlob{}, fmt.Errorf("committed source identity is incomplete")
	}
	repoPath, err := url.PathUnescape(strings.Join(parts[1:], "/"))
	if err != nil {
		return gateEvidenceBlob{}, err
	}
	inputPath, err := filepath.Rel(roots.Main, path)
	if err != nil {
		return gateEvidenceBlob{}, err
	}
	return gateEvidenceBlob{u.Host, repoPath, filepath.ToSlash(inputPath), parts[0], source.Rev, string(body)}, nil
}

func gateEvidenceJSON(blob gateEvidenceBlob) *jsonObj {
	return newJSONObj().set("root", blob.root).set("path", blob.path).set("input_path", blob.inputPath).set("object_id", blob.objectID).set("digest", blob.digest).set("bytes", blob.bytes)
}

func gateEvidenceGitOutput(root string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}
