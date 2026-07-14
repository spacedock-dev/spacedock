// ABOUTME: Pins the first-officer entry point's one eager-reference topology.
// ABOUTME: Dispatch, merge, and write cores stay exact-path deferred.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFirstOfficerKeepsOperationalCoresDeferred(t *testing.T) {
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
	for _, delayed := range []string{
		"references/fo-dispatch-core.md",
		"references/fo-merge-core.md",
		"references/fo-write-core.md",
	} {
		if !strings.Contains(block, delayed) {
			t.Errorf("shared core lacks exact delayed-reference cue %q", delayed)
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

func TestPR495FilingHuntTargetHasDeterministicDelayedAddress(t *testing.T) {
	const capturedHunt = `find / -path /proc -prune -o -iname "fo-write-core*" -print 2>/dev/null`
	const capturedHuntTarget = "fo-write-core"
	root := skillsRoot(t)
	shared, err := os.ReadFile(filepath.Join(root, "first-officer", "references", "first-officer-shared-core.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Base directory for this skill", "references/fo-write-core.md", "never cwd", "search, or retry"} {
		if !strings.Contains(string(shared), want) {
			t.Fatalf("PR #495 filing hunt %q for %s remains possible: deferred rule lacks %q", capturedHunt, capturedHuntTarget, want)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "first-officer", "references", "fo-write-core.md")); err != nil {
		t.Fatalf("PR #495 filing hunt delayed target does not resolve: %v", err)
	}
}
