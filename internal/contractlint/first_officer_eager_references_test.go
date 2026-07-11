// ABOUTME: Pins the first-officer entry point's small eager-reference split.
// ABOUTME: Merge and write preload; the larger dispatch core stays deferred.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFirstOfficerEagerReferencesKeepDispatchCoreDeferred(t *testing.T) {
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
	wantImports := []string{
		"@references/first-officer-shared-core.md",
		"@references/fo-merge-core.md",
		"@references/fo-write-core.md",
	}
	if strings.Join(imports, "\n") != strings.Join(wantImports, "\n") {
		t.Fatalf("first-officer eager imports = %v, want exactly %v (dispatch core remains deferred)", imports, wantImports)
	}
	for _, rel := range wantImports[1:] {
		path := filepath.Join(root, "first-officer", strings.TrimPrefix(rel, "@"))
		body, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("eager reference %s does not resolve to a readable file: %v", rel, err)
		} else if len(body) == 0 {
			t.Errorf("eager reference %s resolves to an empty file", rel)
		}
	}

	sharedPath := filepath.Join(root, "first-officer", "references", "first-officer-shared-core.md")
	shared, err := os.ReadFile(sharedPath)
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	if block := deferredLoadPointsBlock(t, string(shared)); strings.Contains(block, "references/fo-merge-core.md") {
		t.Error("preloaded fo-merge-core.md remains in the shared core's deferred load points")
	}
	for _, bare := range []string{
		"references/fo-merge-core.md",
		"references/fo-write-core.md",
	} {
		if strings.Contains(string(shared), bare) {
			t.Errorf("shared core retains bare eager-reference load cue %q", bare)
		}
	}

	for _, dir := range []string{"fo-merge-core", "fo-smallest-sufficient-mechanism", "fo-write-core"} {
		path := filepath.Join(root, dir)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("rejected promoted-skill directory exists at %s; want no separately callable capability", path)
		}
	}
}

func TestFirstOfficerEagerWriteCoreHasSingleCanonicalBody(t *testing.T) {
	root := skillsRoot(t)
	canonical := filepath.Join(root, "first-officer", "references", "fo-write-core.md")
	body, err := os.ReadFile(canonical)
	if err != nil {
		t.Fatalf("read canonical eager write core: %v", err)
	}
	if !strings.Contains(string(body), "## Mutation Gate") || !strings.Contains(string(body), "## FO Write Scope") {
		t.Fatalf("canonical eager write core does not carry the write contract")
	}

	if _, err := os.Stat(filepath.Join(root, "fo-write-core")); !os.IsNotExist(err) {
		t.Fatalf("standalone fo-write-core wrapper remains: %v", err)
	}
}

func TestPR495FilingHuntTargetIsEagerlyResolvable(t *testing.T) {
	const capturedHunt = `find / -path /proc -prune -o -iname "fo-write-core*" -print 2>/dev/null`
	const capturedHuntTarget = "fo-write-core"
	root := skillsRoot(t)
	entry, err := os.ReadFile(filepath.Join(root, "first-officer", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(entry), "@references/fo-write-core.md") {
		t.Fatalf("PR #495 filing hunt %q for %s remains possible: first-officer entry does not eagerly import write-core", capturedHunt, capturedHuntTarget)
	}
	if _, err := os.Stat(filepath.Join(root, "first-officer", "references", "fo-write-core.md")); err != nil {
		t.Fatalf("PR #495 filing hunt eager target does not resolve: %v", err)
	}
}
