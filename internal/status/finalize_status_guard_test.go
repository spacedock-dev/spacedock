// ABOUTME: AC-3 finalize-status gate — --set/--archive refuse the incomplete-
// ABOUTME: finalize shape (verdict set, status never advanced to terminal).
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// Spike C (docs/dev fo-opus-behavioral-robustness) proved the residual hole
// behind the verdict gate: `--set {slug} completed verdict=PASSED worktree=`
// with no `status={terminal}` in the same call finalizes an entity that never
// advanced past its non-terminal stage, and a follow-on `--archive` then moves
// it into `_archive/` still carrying `status: implementation`. The verdict gate
// (verdict_guard_test.go) only checks that a verdict IS set — it never checks
// that the resulting status is terminal. This gate closes that gap.

// TestFinalizeStatusGateRefusesSpikeCSequence is the AC-3 red-first case:
// Spike C's exact command sequence, byte for byte. Before the gate this exits 0
// and mutates the entity; the gate must refuse it (exit 1, entity untouched,
// verbatim remediation text).
func TestFinalizeStatusGateRefusesSpikeCSequence(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "completed", "verdict=PASSED", "worktree=")
	out, errOut, code := runNative(t, root, env, args...)

	if code != 1 {
		t.Fatalf("finalize with no status={terminal} in the same call must refuse (exit 1), got %d (stdout=%q stderr=%q)", code, out, errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
	wantErr := "Error: entity 002-vendor-script cannot be finalized ('completed') while status 'ideation' is not the terminal stage. " +
		"Set status=done in the same --set (or run 'spacedock merge guard 002-vendor-script'), or use --force."
	if !strings.Contains(errOut, wantErr) {
		t.Fatalf("stderr = %q, want it to contain %q", errOut, wantErr)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: ideation") {
		t.Fatalf("entity must stay unmutated when the gate refuses:\n%s", fm)
	}
	if strings.Contains(fm, "verdict:") {
		t.Fatalf("entity must not have gained a verdict when the gate refuses:\n%s", fm)
	}

	// The --set step above never mutated the entity, so the follow-on --archive
	// half of Spike C's sequence has nothing incomplete to archive: it sees the
	// untouched, verdict-less, non-terminal entity and takes the sanctioned
	// verdict-less-non-terminal-archive path (still exit 0 — see
	// TestFinalizeStatusGateArchiveVerdictLessNonTerminalPasses). The dangerous
	// end-state — `status: implementation` + a false `verdict: PASSED` landing in
	// `_archive/` — is exactly what the --set refusal above prevents from ever
	// existing.
}

// TestFinalizeStatusGateSetWithTerminalStatusAllowed: the contract terminalize
// shape — status={terminal} in the SAME --set call as completed+verdict — passes.
func TestFinalizeStatusGateSetWithTerminalStatusAllowed(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done", "completed", "verdict=passed", "worktree=")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("finalize WITH status={terminal} in the same call must succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	// `verdict=passed` on the CLI stores as the schema's PASSED: updateFrontmatter
	// canonicalises conventional values on every write, so --set and merge guard
	// cannot disagree on stored case.
	if !strings.Contains(fm, "status: done") || !strings.Contains(fm, "verdict: PASSED") {
		t.Fatalf("entity should have terminalized:\n%s", fm)
	}
}

// TestFinalizeStatusGateAlreadyTerminalAllowed: an entity already sitting at the
// terminal stage (no status= in this call) finalizes fine — the gate reads the
// CURRENT status when this call sets none.
func TestFinalizeStatusGateAlreadyTerminalAllowed(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "seq-workflow", map[string]string{
		"already-done.md": "---\nid: \"053\"\ntitle: Already Done\nstatus: done\nscore: \"0.5\"\nsource: x\n---\n# Already Done\n",
	})

	args := append([]string{"--workflow-dir", root}, "--set", "already-done", "completed", "verdict=passed", "worktree=")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("finalize on an already-terminal entity must succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "already-done.md"))
	if !strings.Contains(fm, "verdict: PASSED") {
		t.Fatalf("entity should have gained a verdict:\n%s", fm)
	}
}

// TestFinalizeStatusGateForceBypasses: --force bypasses the gate, same idiom as
// the mod-block / merge-hook / verdict guards.
func TestFinalizeStatusGateForceBypasses(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "completed", "verdict=PASSED", "worktree=", "--force")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("--force must bypass the finalize-status gate (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: ideation") || !strings.Contains(fm, "verdict: PASSED") {
		t.Fatalf("entity should have finalized under --force despite the non-terminal status:\n%s", fm)
	}
}

// TestFinalizeStatusGateDispatchIntoTerminalPasses: a bare dispatch-into-terminal
// (`status=done started`, no `completed`) is not a finalize at all — the gate
// must not fire on it, matching the verdict gate's identical carve-out.
func TestFinalizeStatusGateDispatchIntoTerminalPasses(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done", "started")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("dispatch-into-terminal (no `completed`) must PASS (exit 0), got %d (stderr=%q)", code, errOut)
	}
}

// TestFinalizeStatusGateArchiveRefusesNonTerminalVerdict: the --archive mirror.
// An entity that already carries a non-rejected verdict but sits at a
// non-terminal status — the shape a stale binary (predating the --set gate
// above) or a drifting FO can still produce by hand-editing — must be refused
// at --archive too, so the false-verdict end-state cannot land in `_archive/`
// via that route either.
func TestFinalizeStatusGateArchiveRefusesNonTerminalVerdict(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "seq-workflow", map[string]string{
		"stale-finalized.md": "---\nid: \"054\"\ntitle: Stale Finalized\nstatus: implementation\nverdict: PASSED\nscore: \"0.5\"\nsource: x\n---\n# Stale Finalized\n",
	})

	out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--archive", "stale-finalized")
	if code != 1 {
		t.Fatalf("--archive of a non-terminal entity carrying a verdict must refuse (exit 1), got %d (stderr=%q)", code, errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
	if !strings.Contains(errOut, "cannot be archived") {
		t.Fatalf("stderr must name the archive refusal, got %q", errOut)
	}
	if !strings.Contains(errOut, "PASSED") || !strings.Contains(errOut, "implementation") {
		t.Fatalf("stderr must name the verdict and the non-terminal status, got %q", errOut)
	}

	// --force bypasses.
	_, errOut, code = runNative(t, root, env, "--workflow-dir", root, "--archive", "stale-finalized", "--force")
	if code != 0 {
		t.Fatalf("--force must bypass the archive-side gate (exit 0), got %d (stderr=%q)", code, errOut)
	}
}

// TestFinalizeStatusGateArchiveRejectedNonTerminalPasses: the sanctioned
// reject-then-archive path — verdict: rejected — is exempt from the archive-side
// gate regardless of status, matching the merge-hook guard's identical
// exemption (a rejected entity never ran the merge ceremony).
func TestFinalizeStatusGateArchiveRejectedNonTerminalPasses(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "seq-workflow", map[string]string{
		"rejected-nonterminal.md": "---\nid: \"055\"\ntitle: Rejected Non-Terminal\nstatus: implementation\nverdict: rejected\nscore: \"0.5\"\nsource: x\n---\n# Rejected Non-Terminal\n",
	})

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--archive", "rejected-nonterminal")
	if code != 0 {
		t.Fatalf("reject-then-archive must pass regardless of terminal status (exit 0), got %d (stderr=%q)", code, errOut)
	}
}

// TestFinalizeStatusGateArchiveVerdictLessNonTerminalPasses: archiving a
// non-terminal entity that carries NO verdict at all (an abandoned/cancelled
// task, never finalized) is unaffected — the gate only fires on a non-empty,
// non-rejected verdict.
func TestFinalizeStatusGateArchiveVerdictLessNonTerminalPasses(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--archive", "002-vendor-script")
	if code != 0 {
		t.Fatalf("verdict-less non-terminal archive must pass (exit 0), got %d (stderr=%q)", code, errOut)
	}
}
