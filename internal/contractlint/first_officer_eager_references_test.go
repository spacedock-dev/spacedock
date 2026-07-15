// ABOUTME: Pins the first-officer entry point's lazy write/merge boundary.
// ABOUTME: Only the shared core is eager; write and merge resolve at their triggers.
package contractlint

import (
	"os"
	"path/filepath"
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
	wantImports := []string{
		"@references/first-officer-shared-core.md",
	}
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
			t.Errorf("deferred load-point cue %q occurs %d times, want exactly once", cue, got)
		}
		path := filepath.Join(root, "first-officer", rel)
		body, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("deferred reference %s does not resolve to a readable file: %v", rel, err)
		} else if len(body) == 0 {
			t.Errorf("deferred reference %s resolves to an empty file", rel)
		}
	}

	for _, dir := range []string{"fo-merge-core", "fo-smallest-sufficient-mechanism", "fo-write-core"} {
		path := filepath.Join(root, dir)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("rejected promoted-skill directory exists at %s; want no separately callable capability", path)
		}
	}
}

func TestFirstOfficerDeferredWriteCoreHasSingleCanonicalBody(t *testing.T) {
	root := skillsRoot(t)
	canonical := filepath.Join(root, "first-officer", "references", "fo-write-core.md")
	body, err := os.ReadFile(canonical)
	if err != nil {
		t.Fatalf("read canonical deferred write core: %v", err)
	}
	if !strings.Contains(string(body), "## Mutation Gate") || !strings.Contains(string(body), "## FO Write Scope") {
		t.Fatalf("canonical deferred write core does not carry the write contract")
	}

	if _, err := os.Stat(filepath.Join(root, "fo-write-core")); !os.IsNotExist(err) {
		t.Fatalf("standalone fo-write-core wrapper remains: %v", err)
	}
}
