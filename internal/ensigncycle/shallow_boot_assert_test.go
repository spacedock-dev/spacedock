package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
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

// assertShallowBoot is the host-neutral AC-1 oracle over the run's durable on-disk
// state and final message. It grades, on independent on-disk facts (never a
// transcript phrase as the sole signal):
//
//	(a)  the greet names the gate and reports that it remains held for engage;
//	(b)  NO team artifact on disk (lazy-TeamCreate) AND no worker dispatched (the
//	     gate entity is unchanged, not archived, no worktree created);
//	(c)  the FO stopped for input without engaging or resolving the gate.
//
// The absence-of-team-config is the lazy-TeamCreate proof; the unchanged gate
// frontmatter is the shallow-boot / no-dispatch / no-mutation proof.
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
	// (a) the greet accurately names the held gate and the engage boundary.
	lowerFinal := strings.ToLower(o.finalMessage)
	for _, want := range []string{"gate check", "review", "engage"} {
		if !strings.Contains(lowerFinal, want) {
			return fmt.Errorf("the greet did not report the held gate accurately: missing %q", want)
		}
	}
	return nil
}

// gatherShallowBootObservation reads the run's durable on-disk state into a
// shallowBootObservation: the gate entity's post-boot frontmatter, its
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
