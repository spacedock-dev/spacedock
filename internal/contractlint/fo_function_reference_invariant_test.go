// ABOUTME: Preserves First Officer component byte caps.
package contractlint

import (
	"path/filepath"
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
//
// Captain-approved cap raise (2026-08-16, "ok B then", entity
// prepare-initial-gated-stage-from-seed, id ra9qzfz94hzgsq938jz998mj):
// fo-gate-lifecycle/SKILL.md 7000 -> 7700. The gated initial-stage fix adds an
// initial-stage exception to the gate lifecycle, and the file was already
// saturated at 6993 bytes — seven bytes of headroom. Unlike the 2026-08-02
// raise above, minimize-first could not avoid this one, which is the evidence
// docs/roadmap/durable-decisions/staff-review-sprint-close.md requires: the
// absolute floor, keeping only the readiness-table and heading-scope edits and
// dropping the initial-stage paragraph entirely, still measures 7072; a pointer
// restructure lands at 7219-7340, saving only 193 bytes over the approved
// wording (7533) while relocating prose the gate approved in place. Since every
// option needed a raise, the captain chose the one that keeps the approved
// wording where it was approved and raised the cap by the measured remainder
// plus headroom.
func TestFOInstructionComponentCaps(t *testing.T) {
	for rel, cap := range map[string]int{
		"skills/first-officer/references/first-officer-shared-core.md": 26900,
		"skills/fo-gate-lifecycle/SKILL.md":                            7700,
	} {
		if got := len([]byte(readRepoFile(t, filepath.FromSlash(rel)))); got > cap {
			t.Errorf("%s = %d bytes, component cap %d", rel, got, cap)
		}
	}
}
