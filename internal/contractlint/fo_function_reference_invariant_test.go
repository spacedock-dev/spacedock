// ABOUTME: Preserves First Officer component byte caps and reference-file topology.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Captain-approved cap raise (2026-08-02, entity collapse-gate-approval-ceremony,
// id 7fhzvvk8d5smj858bp47xbjq): mechanism 1/2/3's doc diff (the gate
// record/consume sync=/phase= discriminator, --consume, and dispatch build
// --stamp) needed new load-bearing prose in both files at the same time an
// unrelated, independently-landed change (988163969, "Withdraw stale open gate
// attempt") was also growing fo-gate-lifecycle/SKILL.md. Per
// docs/roadmap/durable-decisions/staff-review-sprint-close.md's evidence
// requirement for changing this capped set: the content was first restructured
// to a pointer — the full explanation moved to
// skills/first-officer/references/fo-dispatch-core.md (uncapped), leaving only
// a one-line cross-reference in each capped file, cutting the overage 65-68%
// (SKILL.md 1003B over -> 357B over; shared-core.md 207B over -> 67B over)
// before the cap was raised at all. What remained had nowhere left to go
// without deleting the pointer itself or cutting unrelated content, so the
// captain raised both caps by the measured remainder plus headroom.
func TestFOInstructionComponentCaps(t *testing.T) {
	for rel, cap := range map[string]int{
		"skills/first-officer/references/first-officer-shared-core.md": 26900,
		"skills/fo-gate-lifecycle/SKILL.md":                            7000,
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
