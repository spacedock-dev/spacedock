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
	// interactiveTransport is true only when the host ran its interactive TUI,
	// rather than a headless prompt that asked the model to imitate a greeting.
	interactiveTransport bool
	// sessionResidentAfterGreet proves the host returned to its input loop after
	// the greet instead of exiting like a headless invocation.
	sessionResidentAfterGreet bool
	// finalMessage is the FO's final greet output.
	finalMessage string
	// gateBefore / gateAfter is the gate-check entity frontmatter before and after
	// the boot — it must be byte-identical (the FO presented the gate, did not
	// dispatch or self-approve).
	gateBefore string
	gateAfter  string
	// mergedBefore / mergedAfter is the PR-bearing entity before and after the
	// read-only boot. It must remain byte-identical and active until engage.
	mergedBefore string
	mergedAfter  string
	// mergedArchived is true when the PR-bearing entity was moved to _archive/.
	// Read-only boot must leave it active.
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

// assertShallowBoot is the host-neutral AC-1 oracle over the run's durable on-disk
// state and final message. It grades, on independent on-disk facts (never a
// transcript phrase as the sole signal):
//
//	(a)  the greet names the PR-bearing entity's local state, then offers engage;
//	(a2) the PR-bearing entity remains byte-identical and active until engage;
//	(b)  NO team artifact on disk (lazy-TeamCreate) AND no worker dispatched (the
//	     gate entity is unchanged, not archived, no worktree created);
//	(c)  the FO stopped for input without converging or mutating the workflow.
//
// The absence-of-team-config is the lazy-TeamCreate proof; the unchanged gate
// frontmatter and unchanged PR-bearing entity are the read-only boot proof.
func assertShallowBoot(o shallowBootObservation) error {
	if !o.interactiveTransport {
		return fmt.Errorf("shallow boot did not run through an actual interactive transport")
	}
	if !o.sessionResidentAfterGreet {
		return fmt.Errorf("interactive host was not resident after the greet-and-stop turn")
	}
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
	// (a2) live PR discovery and advancement belong to engage, never read-only boot.
	if o.mergedArchived {
		return fmt.Errorf("the PR-bearing entity was archived during read-only boot — live PR advancement belongs to engage")
	}
	if o.mergedBefore == "" {
		return fmt.Errorf("the PR-bearing fixture was missing before boot")
	}
	if o.mergedBefore != o.mergedAfter {
		return fmt.Errorf("the PR-bearing entity changed during read-only boot — convergence belongs to engage")
	}
	// (a) the greet reports the local PR mirror and leaves convergence behind engage.
	lowerFinal := strings.ToLower(o.finalMessage)
	for _, required := range []string{"merged-pr", "engage"} {
		if !strings.Contains(lowerFinal, required) {
			return fmt.Errorf("the greet did not name %q while reporting read-only startup state", required)
		}
	}
	return nil
}

// gatherShallowBootObservation reads the run's durable on-disk state into a
// shallowBootObservation: both entities' post-boot state, the archive/worktree
// facts, and the team-config-on-disk check under the host's team root. It is
// host-neutral over the entity/archive/worktree state; teamRoot is the
// host's team-config root (Claude: {home}/.claude/teams) — an empty teamRoot means
// the host writes no team config there and the check is vacuously false.
func gatherShallowBootObservation(t *testing.T, workflowRoot, teamRoot string, fx shallowBootFixture, gateBefore, mergedBefore, finalMessage string) shallowBootObservation {
	t.Helper()
	o := shallowBootObservation{
		finalMessage: finalMessage,
		gateBefore:   gateBefore,
		gateAfter:    readFileAllowMissing(fx.gateEntityPath),
		mergedBefore: mergedBefore,
	}
	// Prefer the archive when present so a forbidden boot-time advancement is visible.
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
