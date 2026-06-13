package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// shallowBootObservation is the host-neutral set of durable, on-disk facts the
// shallow-boot scenario grades, plus the FO's final message. The runner gathers
// these from THIS run's filesystem (per-run temp HOME for the team root, the
// workflow root for entity/archive/worktree state) so the assertion never reads a
// stale prior team or a transcript phrase as authoritative.
type shallowBootObservation struct {
	// finalMessage is the FO's final greet output.
	finalMessage string
	// gateBefore / gateAfter is the gate-check entity frontmatter before and after
	// the boot — it must be byte-identical (the FO presented the gate, did not
	// dispatch or self-approve).
	gateBefore string
	gateAfter  string
	// mergedAfter is the merged-PR entity's durable state read from wherever it
	// lives after the boot (the archive, once S7b advances it). Empty when the file
	// is gone from the active dir AND absent from the archive.
	mergedAfter string
	// mergedArchived is true when the merged-PR entity was moved to _archive/.
	mergedArchived bool
	// gateWorktreeCreated is true when a .worktrees/ dir was created for the gate
	// entity (it must NOT be — no dispatch happened).
	gateWorktreeCreated bool
	// gateArchived is true when the gate entity was archived (it must NOT be — it is
	// parked at its gate).
	gateArchived bool
	// teamConfigOnDisk is true when a team config.json exists under THIS run's team
	// root (it must NOT be — lazy-TeamCreate means no team is created at boot).
	teamConfigOnDisk bool
}

var (
	mergedTerminalStatus = regexp.MustCompile(`(?im)^status:\s*done\s*$`)
	mergedVerdictPassed  = regexp.MustCompile(`(?im)^verdict:\s*PASSED\s*$`)
	mergedModBlockClear  = regexp.MustCompile(`(?im)^mod-block:\s*$`)
)

// assertShallowBoot is the host-neutral AC-1 oracle over the run's durable on-disk
// state and final message. It grades, on independent on-disk facts (never a
// transcript phrase as the sole signal):
//
//	(a)  the greet presents a gate review + decision prompt;
//	(a2) the S7b merged-PR entity is advanced before-greet — terminal frontmatter
//	     (done / verdict PASSED / mod-block cleared) AND archived (the M3 proof);
//	(b)  NO team artifact on disk (lazy-TeamCreate) AND no worker dispatched (the
//	     gate entity is unchanged, not archived, no worktree created);
//	(c)  the FO stopped for input (it presented a gate, did not advance it).
//
// The absence-of-team-config is the lazy-TeamCreate proof; the unchanged gate
// frontmatter is the shallow-boot / no-dispatch proof; the advanced+archived merged
// entity is the S7b proof.
func assertShallowBoot(o shallowBootObservation) error {
	// (b) lazy-TeamCreate: no team artifact created at boot.
	if o.teamConfigOnDisk {
		return fmt.Errorf("a team config.json exists under this run's team root — the boot created a team (lazy-TeamCreate was not honored)")
	}
	// (b) no-dispatch: the gate entity must be byte-identical and not advanced.
	if o.gateBefore != o.gateAfter {
		return fmt.Errorf("the gated entity's frontmatter changed during boot — a worker was dispatched or the gate self-resolved")
	}
	if !reviewStatus.MatchString(o.gateAfter) {
		return fmt.Errorf("the gated entity is no longer at status: review — it was advanced past its gate")
	}
	if completedSet.MatchString(o.gateAfter) || verdictSetFM.MatchString(o.gateAfter) {
		return fmt.Errorf("the gated entity has completed/verdict set — the boot self-approved instead of presenting the gate")
	}
	if o.gateArchived {
		return fmt.Errorf("the gated entity was archived during boot — it must stay parked at its gate")
	}
	if o.gateWorktreeCreated {
		return fmt.Errorf("a worktree was created for the gated entity — a dispatch happened at boot")
	}
	// (a2) S7b: the merged-PR entity is advanced and archived before-greet.
	if !o.mergedArchived {
		return fmt.Errorf("the merged-PR entity was not archived — S7b did not advance it before the greet")
	}
	if !mergedTerminalStatus.MatchString(o.mergedAfter) {
		return fmt.Errorf("the merged-PR entity is not at the terminal stage (status: done) — S7b advancement is incomplete")
	}
	if !mergedVerdictPassed.MatchString(o.mergedAfter) {
		return fmt.Errorf("the merged-PR entity has no verdict: PASSED — S7b advancement is incomplete")
	}
	if !mergedModBlockClear.MatchString(o.mergedAfter) {
		return fmt.Errorf("the merged-PR entity still carries a mod-block — S7b did not clear it on advancement")
	}
	// (a) the greet presents a gate review + decision prompt.
	lowerFinal := strings.ToLower(o.finalMessage)
	if !strings.Contains(lowerFinal, "gate review:") || !strings.Contains(lowerFinal, "decision:") {
		return fmt.Errorf("the greet did not present a gate review and decision prompt")
	}
	return nil
}

// gatherShallowBootObservation reads the run's durable on-disk state into a
// shallowBootObservation: the gate entity's post-boot frontmatter, the merged
// entity's state (from the archive once S7b advances it, else its active path), the
// archive/worktree facts, and the team-config-on-disk check under the host's team
// root. It is host-neutral over the entity/archive/worktree state; teamRoot is the
// host's team-config root (Claude: {home}/.claude/teams) — an empty teamRoot means
// the host writes no team config there and the check is vacuously false.
func gatherShallowBootObservation(t *testing.T, workflowRoot, teamRoot string, fx shallowBootFixture, gateBefore, finalMessage string) shallowBootObservation {
	t.Helper()
	o := shallowBootObservation{
		finalMessage: finalMessage,
		gateBefore:   gateBefore,
		gateAfter:    readFileAllowMissing(fx.gateEntityPath),
	}
	// The merged entity lives in _archive once S7b advances it; before that it stays
	// at its active path. Read whichever exists so the assertion sees its real state.
	if data, err := os.ReadFile(fx.mergedArchive); err == nil {
		o.mergedAfter = string(data)
		o.mergedArchived = true
	} else {
		o.mergedAfter = readFileAllowMissing(fx.mergedEntityPath)
	}
	if _, err := os.Stat(fx.gateEntityArchivePath(workflowRoot)); err == nil {
		o.gateArchived = true
	}
	// A worktree created for the gate entity is a dispatch fingerprint. The default
	// ensign worker_key is spacedock-ensign, so the dir would be
	// .worktrees/spacedock-ensign-gate-check; glob loosely on the slug to catch any
	// worker_key.
	if matches, _ := filepath.Glob(filepath.Join(workflowRoot, ".worktrees", "*gate-check*")); len(matches) > 0 {
		o.gateWorktreeCreated = true
	}
	if teamRoot != "" {
		if matches, _ := filepath.Glob(filepath.Join(teamRoot, "*", "config.json")); len(matches) > 0 {
			o.teamConfigOnDisk = true
		}
	}
	return o
}

// gateEntityArchivePath is the path the gate entity would occupy if it were
// archived (it must NOT be — it is parked at its gate).
func (fx shallowBootFixture) gateEntityArchivePath(workflowRoot string) string {
	return filepath.Join(workflowRoot, "_archive", "gate-check.md")
}

// readFileAllowMissing returns the file's content, or "" when the file is absent
// (e.g. an entity moved to the archive leaves its active path empty).
func readFileAllowMissing(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
