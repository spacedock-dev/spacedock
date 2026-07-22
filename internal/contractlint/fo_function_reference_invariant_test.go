// ABOUTME: Preserves First Officer component byte caps and reference-file topology.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFOInstructionComponentCaps(t *testing.T) {
	for rel, cap := range map[string]int{
		"skills/first-officer/references/first-officer-shared-core.md": 27194,
		"skills/fo-gate-lifecycle/SKILL.md":                            6600,
	} {
		if got := len([]byte(readRepoFile(t, filepath.FromSlash(rel)))); got > cap {
			t.Errorf("%s = %d bytes, component cap %d", rel, got, cap)
		}
	}
}

func TestFirstOfficerReferenceTopology(t *testing.T) {
	root := skillsRoot(t)
	entry := readRepoFile(t, filepath.Join("skills", "first-officer", "SKILL.md"))
	var imports []string
	for _, line := range strings.Split(entry, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "@references/") {
			imports = append(imports, strings.TrimSpace(line))
		}
	}
	want := []string{"@references/first-officer-shared-core.md"}
	if strings.Join(imports, "\n") != strings.Join(want, "\n") {
		t.Fatalf("eager imports = %v, want topology %v", imports, want)
	}
	for _, rel := range []string{
		"references/first-officer-shared-core.md",
		"references/fo-write-core.md",
		"references/fo-merge-core.md",
	} {
		path := filepath.Join(root, "first-officer", rel)
		if info, err := os.Stat(path); err != nil || info.Size() == 0 {
			t.Errorf("canonical first-officer body %s does not resolve non-empty: %v", rel, err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "fo-write-core")); !os.IsNotExist(err) {
		t.Errorf("redundant standalone fo-write-core entry surface remains: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "first-officer", "references", "fo-smallest-sufficient-mechanism.md")); !os.IsNotExist(err) {
		t.Errorf("duplicated smallest-sufficient eager body remains: %v", err)
	}
}
