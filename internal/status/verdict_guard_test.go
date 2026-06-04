// ABOUTME: AC-4 verdict gate — --set/--archive refuse the FINALIZE action
// ABOUTME: (setting `completed`, or archiving) without a verdict; dispatch-into-terminal and --force pass.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// The verdict gate keys on the FINALIZE action — setting `completed`
// (empty→non-empty), the contract terminalize shape (shared-core step 7:
// `completed verdict={verdict} worktree=`) — NOT on `status==terminal` alone. A
// bare dispatch-into-terminal (`--set status=done … started worktree=…`, no
// `completed`) is a legitimate verdict-less transition (the verdict is the
// OUTCOME of the work that hasn't happened yet), so it must PASS; only the
// finalize without a verdict reds. This is captain option A — the fix for the
// live flake where the broad `status==terminal` gate blocked the dispatch step
// of a backlog→done workflow and made the cycle non-deterministic.

// TestVerdictGateDispatchIntoTerminalPasses (M1 case a): a dispatch that advances
// status INTO the terminal stage WITHOUT `completed` is not a finalize — it must
// pass even with no verdict. This is the transition the broad gate wrongly blocked.
func TestVerdictGateDispatchIntoTerminalPasses(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done", "started")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("dispatch-into-terminal (status=done started, no `completed`, no verdict) must PASS (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: done") {
		t.Fatalf("entity should have advanced to done:\n%s", fm)
	}
}

// TestVerdictGateFinalizeWithoutVerdictRefused (M1 case b): a finalize (`completed`
// set) with an empty post-state verdict is refused (exit 1, entity unmutated).
func TestVerdictGateFinalizeWithoutVerdictRefused(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done", "completed", "worktree=")
	out, errOut, code := runNative(t, root, env, args...)

	if code != 1 {
		t.Fatalf("finalize without a verdict must refuse (exit 1), got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(errOut, "verdict") {
		t.Fatalf("stderr must name the missing verdict, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if strings.Contains(fm, "status: done") {
		t.Fatalf("entity must stay unmutated when the verdict gate refuses:\n%s", fm)
	}
}

// TestVerdictGateFinalizeWithVerdictAllowed (M1): the contract terminalize shape
// WITH a verdict in the same call is allowed (exit 0, status advances, verdict set).
func TestVerdictGateFinalizeWithVerdictAllowed(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done", "completed", "verdict=passed", "worktree=")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("finalize WITH a verdict must succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: done") || !strings.Contains(fm, "verdict: passed") {
		t.Fatalf("entity should have terminalized with a verdict:\n%s", fm)
	}
}

// TestVerdictGateIdempotentTerminalSetPasses (M1 case c / M3): a finalize on an
// entity that ALREADY carries a verdict passes even when THIS --set names no
// verdict — the post-state verdict is non-empty (read from the entity). The
// fixture carries a pre-existing `verdict:` so this exercises the
// postUpdateVerdict := currentVerdict read (M3): zeroing that read reds this test.
func TestVerdictGateIdempotentTerminalSetPasses(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "seq-workflow", map[string]string{
		"already-verdicted.md": "---\nid: \"050\"\ntitle: Already Verdicted\nstatus: implementation\nverdict: passed\nscore: \"0.5\"\nsource: x\n---\n# Already Verdicted\n",
	})

	args := append([]string{"--workflow-dir", root}, "--set", "already-verdicted", "status=done", "completed", "worktree=")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("a finalize on an already-verdicted entity must PASS without re-naming verdict (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "already-verdicted.md"))
	if !strings.Contains(fm, "status: done") || !strings.Contains(fm, "verdict: passed") {
		t.Fatalf("entity should have terminalized keeping its verdict:\n%s", fm)
	}
}

// TestVerdictGateFinalizeForceBypasses (M1 case d): --force bypasses the gate.
func TestVerdictGateFinalizeForceBypasses(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=done", "completed", "worktree=", "--force")
	_, errOut, code := runNative(t, root, env, args...)

	if code != 0 {
		t.Fatalf("--force must bypass the verdict gate (exit 0), got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: done") {
		t.Fatalf("entity should have terminalized under --force:\n%s", fm)
	}
}

// TestVerdictGateNonDoneTerminalStage (M3 polish): the gate is data-driven over
// the README's terminal stage names, not hardcoded 'done'. The split-root
// workflow's terminal stage is `review` — a finalize into it without a verdict
// must red, proving the gate reads the declared terminal set.
func TestVerdictGateNonDoneTerminalStage(t *testing.T) {
	env := pinnedEnv(t)
	def, _ := buildSplitRoot(t, splitRootReadme, map[string]string{
		"add-login.md": "---\nstatus: implementation\n---\n",
	})

	out, errOut, code := runNative(t, def, env, "--workflow-dir", def, "--set", "add-login", "status=review", "completed", "worktree=")
	if code != 1 {
		t.Fatalf("finalize into the non-'done' terminal stage 'review' without a verdict must refuse (exit 1), got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(errOut, "verdict") {
		t.Fatalf("stderr must name the missing verdict, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
}

// TestVerdictGateArchiveWithoutVerdictRefused (M1 case e): finalization-by-archive
// also requires a verdict — a verdict-less terminal entity cannot be archived
// (the gate must not be routed around via --archive). --force bypasses.
func TestVerdictGateArchiveWithoutVerdictRefused(t *testing.T) {
	env := pinnedEnv(t)
	// An entity already in the terminal stage but with NO verdict — archiving it
	// finalizes it, so the gate must refuse without a verdict.
	root := stageFixtureWith(t, "seq-workflow", map[string]string{
		"no-verdict-terminal.md": "---\nid: \"051\"\ntitle: No Verdict Terminal\nstatus: done\nscore: \"0.5\"\nsource: x\n---\n# No Verdict Terminal\n",
	})

	out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--archive", "no-verdict-terminal")
	if code != 1 {
		t.Fatalf("--archive of a verdict-less terminal entity must refuse (exit 1), got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(errOut, "verdict") {
		t.Fatalf("stderr must name the missing verdict, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}

	// --force bypasses the archive verdict gate.
	_, errOut, code = runNative(t, root, env, "--workflow-dir", root, "--archive", "no-verdict-terminal", "--force")
	if code != 0 {
		t.Fatalf("--force must bypass the archive verdict gate (exit 0), got %d (stderr=%q)", code, errOut)
	}
}

// TestVerdictGateArchiveWithVerdictAllowed (M1): archiving a terminal entity that
// carries a verdict is allowed.
func TestVerdictGateArchiveWithVerdictAllowed(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "seq-workflow", map[string]string{
		"verdicted-terminal.md": "---\nid: \"052\"\ntitle: Verdicted Terminal\nstatus: done\nverdict: passed\nscore: \"0.5\"\nsource: x\n---\n# Verdicted Terminal\n",
	})

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--archive", "verdicted-terminal")
	if code != 0 {
		t.Fatalf("--archive of a verdicted terminal entity must succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
}
