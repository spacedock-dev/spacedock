// ABOUTME: Pins the first-officer entry point's lazy write/merge file topology.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestFirstOfficerEntryEagerlyImportsOnlySharedCore(t *testing.T) {
	root := skillsRoot(t)
	skillPath := filepath.Join(root, "first-officer", "SKILL.md")
	data, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read first-officer skill: %v", err)
	}

	var imports []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "@references/") {
			imports = append(imports, line)
		}
	}
	wantImports := []string{"@references/first-officer-shared-core.md"}
	if strings.Join(imports, "\n") != strings.Join(wantImports, "\n") {
		t.Fatalf("first-officer eager imports = %v, want exactly %v", imports, wantImports)
	}

	sharedPath := filepath.Join(root, "first-officer", "references", "first-officer-shared-core.md")
	shared, err := os.ReadFile(sharedPath)
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	block := deferredLoadPointsBlock(t, string(shared))
	for _, rel := range []string{
		"references/fo-write-core.md",
		"references/fo-merge-core.md",
	} {
		cue := "{first_officer_base}/" + rel
		if got := strings.Count(block, cue); got != 1 {
			t.Errorf("deferred load-point reference %q occurs %d times, want exactly once", cue, got)
		}
		path := filepath.Join(root, "first-officer", rel)
		info, err := os.Stat(path)
		if err != nil || info.Size() == 0 {
			t.Errorf("deferred reference %s does not resolve to a non-empty file: %v", rel, err)
		}
	}

	for _, dir := range []string{"fo-merge-core", "fo-smallest-sufficient-mechanism", "fo-write-core"} {
		path := filepath.Join(root, dir)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("redundant promoted-skill directory exists at %s", path)
		}
	}
}

func TestFirstOfficerDeferredWriteCoreHasSingleCanonicalFile(t *testing.T) {
	root := skillsRoot(t)
	canonical := filepath.Join(root, "first-officer", "references", "fo-write-core.md")
	info, err := os.Stat(canonical)
	if err != nil || info.Size() == 0 {
		t.Fatalf("canonical deferred write core does not resolve non-empty: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "fo-write-core")); !os.IsNotExist(err) {
		t.Fatalf("standalone fo-write-core wrapper remains: %v", err)
	}
}

func deferredLoadPointsBlock(t *testing.T, body string) string {
	t.Helper()
	loc := regexp.MustCompile(`(?m)^## Deferred load points$`).FindStringIndex(body)
	if loc == nil {
		t.Fatal("shared core has no deferred load-points section")
	}
	rest := body[loc[1]:]
	if end := regexp.MustCompile(`(?m)^## `).FindStringIndex(rest); end != nil {
		rest = rest[:end[0]]
	}
	return rest
}
