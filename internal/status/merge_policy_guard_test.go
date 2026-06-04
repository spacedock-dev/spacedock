// ABOUTME: merge: local guard-relaxation parity — the registered-merge-hook
// ABOUTME: fixture exercises the terminal-guard branch the policy relaxes.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// assertMergeGolden runs args through native on a fresh fixture and freezes the
// normalized three channels to a golden keyed by name. It returns the native exit
// code/stdout/stderr so the caller can additionally assert the pass/refuse
// outcome. The root is normalized to a placeholder (the --archive dest carries
// an absolute path that differs per temp run).
func assertMergeGolden(t *testing.T, name, fixture string, args ...string) (int, string, string) {
	t.Helper()
	env := pinnedEnv(t)
	nativeRoot := stageFixture(t, fixture)

	nArgs := append([]string{"--workflow-dir", nativeRoot}, args...)
	nOut, nErr, nCode := runNative(t, nativeRoot, env, nArgs...)

	assertEnvelopeGolden(t, name, goldenEnvelope{
		stdout: normalize(nOut, nativeRoot), stderr: normalize(nErr, nativeRoot), exit: nCode,
	})
	return nCode, nOut, nErr
}

// TestMergeLocalNoSentinelTerminalSetSucceeds (AC-3): under merge: local, a
// terminal --set with empty pr and empty mod-block succeeds without --force even
// though a merge hook is registered — the policy exempts the pr-requirement.
func TestMergeLocalNoSentinelTerminalSetSucceeds(t *testing.T) {
	code, _, errOut := assertMergeGolden(t, "merge-local-nosentinel-set", "merge-local-workflow",
		"--set", "020-no-sentinel", "status=done", "verdict=passed")
	if code != 0 {
		t.Fatalf("merge: local terminal --set should succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
}

// TestMergeLocalNoSentinelArchiveSucceeds (AC-3): under merge: local, --archive
// with empty pr and empty mod-block succeeds without --force.
func TestMergeLocalNoSentinelArchiveSucceeds(t *testing.T) {
	code, _, errOut := assertMergeGolden(t, "merge-local-nosentinel-archive", "merge-local-workflow",
		"--archive", "020-no-sentinel")
	if code != 0 {
		t.Fatalf("merge: local --archive should succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
}

// TestSentinelSatisfiesGuardTerminalSet (AC-1): the post-merge pr=local-merge:{sha}
// sentinel satisfies the merge-hook guard with NO --force, so a terminal --set
// succeeds. This holds under merge: local too — the sentinel records the landed
// merge regardless of policy.
func TestSentinelSatisfiesGuardTerminalSet(t *testing.T) {
	code, _, errOut := assertMergeGolden(t, "merge-sentinel-set", "merge-local-workflow",
		"--set", "010-sentinel", "status=done", "verdict=passed")
	if code != 0 {
		t.Fatalf("sentinel terminal --set should succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
}

// TestSentinelSatisfiesGuardArchive (AC-1): the sentinel also satisfies the
// --archive merge-hook guard with NO --force.
func TestSentinelSatisfiesGuardArchive(t *testing.T) {
	code, _, errOut := assertMergeGolden(t, "merge-sentinel-archive", "merge-local-workflow",
		"--archive", "010-sentinel")
	if code != 0 {
		t.Fatalf("sentinel --archive should succeed (exit 0), got %d (stderr=%q)", code, errOut)
	}
}

// TestMergeLocalModBlockPendingStillBlocks (AC-3, the wrongful-terminalization
// mechanical survivor): merge: local relaxes the pr-requirement of the merge-hook
// check, but NOT the mod-block-pending guard. An in-flight mod-block must still
// refuse a terminal --set under merge: local — the set->clear->terminalize
// ceremony separation stays mechanically enforced. This is the named scenario
// that the safety invariant turns on: an FO must not collapse the ceremony, and
// the mechanism still catches a terminal transition that combines a live block.
func TestMergeLocalModBlockPendingStillBlocks(t *testing.T) {
	code, out, errOut := assertMergeGolden(t, "merge-local-modblock-pending", "merge-local-workflow",
		"--set", "030-pending", "status=done")
	if code != 1 {
		t.Fatalf("mod-block-pending terminal --set must refuse (exit 1) under merge: local, got %d", code)
	}
	if !strings.Contains(errOut, "pending mod-block (merge:local-merge)") {
		t.Fatalf("stderr should name the pending mod-block, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
}

// TestMergeLocalCombinedClearAndTerminalizeRefused (AC-3, the wrongful-
// terminalization mechanical survivor): merge: local does NOT permit collapsing
// the mandatory clear-then-terminalize ceremony into one call. Clearing mod-block
// in the SAME --set that terminalizes is refused regardless of policy, so the
// audit history must show the block resolving separately from terminalization.
func TestMergeLocalCombinedClearAndTerminalizeRefused(t *testing.T) {
	code, out, errOut := assertMergeGolden(t, "merge-local-combined-clear-terminalize", "merge-local-workflow",
		"--set", "030-pending", "mod-block=", "status=done")
	if code != 1 {
		t.Fatalf("combined clear+terminalize must refuse (exit 1) under merge: local, got %d", code)
	}
	if !strings.Contains(errOut, "combined mod-block clear with terminal transition") {
		t.Fatalf("stderr should name the combined-clear refusal, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
}

// TestMergePrDefaultNoSentinelStillRefuses (AC-1 companion / AC-4): the SAME
// no-sentinel/no-mod-block entity, under the DEFAULT policy (merge: key absent),
// still refuses a terminal --set — the merge-hook catch is preserved. This is the
// byte-identical-to-today guarantee for un-declared workflows.
func TestMergePrDefaultNoSentinelStillRefuses(t *testing.T) {
	code, out, errOut := assertMergeGolden(t, "merge-pr-default-nosentinel-set", "merge-pr-workflow",
		"--set", "020-no-sentinel", "status=done", "verdict=passed")
	if code != 1 {
		t.Fatalf("default-policy terminal --set must still refuse (exit 1), got %d", code)
	}
	if !strings.Contains(errOut, "cannot advance to terminal") {
		t.Fatalf("stderr should name the merge-hook refusal, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
}

// TestMergePrDefaultNoSentinelArchiveStillRefuses (AC-4): the default-policy
// --archive of the same no-sentinel entity also still refuses.
func TestMergePrDefaultNoSentinelArchiveStillRefuses(t *testing.T) {
	code, out, errOut := assertMergeGolden(t, "merge-pr-default-nosentinel-archive", "merge-pr-workflow",
		"--archive", "020-no-sentinel")
	if code != 1 {
		t.Fatalf("default-policy --archive must still refuse (exit 1), got %d", code)
	}
	if !strings.Contains(errOut, "cannot be archived") {
		t.Fatalf("stderr should name the merge-hook archive refusal, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
}

// TestRejectedVerdictArchiveMatchesSet (terminal-guard rejected-consistency AC-1):
// under the default merge: pr policy with a registered merge hook, a rejected
// entity (verdict: rejected, empty pr/mod-block) must be treated the SAME by --set
// and --archive. A rejected entity never ran the merge ceremony, so the merge-hook
// pr-requirement is vacuous for it: --set already exempts verdict=rejected, and
// --archive must match. The two commands run on the SAME staged root in sequence
// (reject → terminalize → archive), and the test asserts the surfaces AGREE.
// RED on today's code: --set → exit 0, --archive → exit 1 (the asymmetry).
// GREEN after the fix: both → exit 0.
func TestRejectedVerdictArchiveMatchesSet(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "merge-pr-workflow")

	setArgs := []string{"--workflow-dir", root, "--set", "040-rejected", "status=done"}
	setOut, setErr, setCode := runNative(t, root, env, setArgs...)
	assertEnvelopeGolden(t, "merge-pr-rejected-set", goldenEnvelope{
		stdout: normalize(setOut, root), stderr: normalize(setErr, root), exit: setCode,
	})

	archiveArgs := []string{"--workflow-dir", root, "--archive", "040-rejected"}
	archiveOut, archiveErr, archiveCode := runNative(t, root, env, archiveArgs...)
	assertEnvelopeGolden(t, "merge-pr-rejected-archive", goldenEnvelope{
		stdout: normalize(archiveOut, root), stderr: normalize(archiveErr, root), exit: archiveCode,
	})

	if setCode != archiveCode {
		t.Fatalf("--set and --archive must agree on verdict=rejected, got set=%d (stderr=%q) archive=%d (stderr=%q)",
			setCode, setErr, archiveCode, archiveErr)
	}
	if setCode != 0 {
		t.Fatalf("rejected-entity terminal --set should succeed (exit 0), got %d (stderr=%q)", setCode, setErr)
	}
	if archiveCode != 0 {
		t.Fatalf("rejected-entity --archive should succeed (exit 0), got %d (stderr=%q)", archiveCode, archiveErr)
	}
}

// TestRejectedVerdictModBlockPendingArchiveRefuses (terminal-guard rejected-
// consistency AC-2): the verdict=rejected escape relaxes ONLY the merge-hook
// pr-requirement, not the policy-independent mod-block-pending guard. A rejected
// entity with a live mod-block must still be refused by --archive (exit 1), naming
// the pending mod-block — the ceremony-separation invariant survives the verdict
// escape just as it survives merge: local.
func TestRejectedVerdictModBlockPendingArchiveRefuses(t *testing.T) {
	code, out, errOut := assertMergeGolden(t, "merge-pr-rejected-pending-archive", "merge-pr-workflow",
		"--archive", "050-rejected-pending")
	if code != 1 {
		t.Fatalf("rejected-but-mod-block-pending --archive must refuse (exit 1), got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(errOut, "pending mod-block (merge:local-merge)") {
		t.Fatalf("stderr should name the pending mod-block, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
}

// TestNonRejectedVerdictArchiveStillRefuses (terminal-guard rejected-consistency,
// over-wide-exemption complement): the verdict escape is scoped to verdict=rejected
// ONLY. A NON-rejected, non-empty verdict (verdict: passed) with empty pr/mod-block
// under the default merge: pr policy must STILL be refused by --archive — this is
// exactly the case the merge-hook guard exists for (an accepted outcome claimed
// without the merge ceremony). This pins the exemption against over-widening:
// changing `verdict != "rejected"` to `verdict == ""` (any non-empty verdict
// escapes) goes RED here. The existing default-refuse fixture (020-no-sentinel)
// carries an EMPTY verdict, so it still refuses under that widening and gives no
// signal — this case supplies the missing one.
func TestNonRejectedVerdictArchiveStillRefuses(t *testing.T) {
	code, out, errOut := assertMergeGolden(t, "merge-pr-passed-nosentinel-archive", "merge-pr-workflow",
		"--archive", "060-passed-nosentinel")
	if code != 1 {
		t.Fatalf("non-rejected-verdict --archive must still refuse (exit 1), got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(errOut, "cannot be archived") {
		t.Fatalf("stderr should name the merge-hook archive refusal, got %q", errOut)
	}
	if out != "" {
		t.Fatalf("stdout must be empty on rejection, got %q", out)
	}
}

// TestSentinelDisplaysAsLocal (AC-2): a pr field of local-merge:{short-sha}
// renders in the status table as `{short-sha} (local)`, distinguishable from a
// real PR reference. Native and oracle must agree.
func TestSentinelDisplaysAsLocal(t *testing.T) {
	env := pinnedEnv(t)
	nativeRoot := stageFixture(t, "merge-local-workflow")

	args := append([]string{"--workflow-dir", nativeRoot}, "--fields", "pr")
	nOut, nErr, nCode := runNative(t, nativeRoot, env, args...)

	if nCode != 0 {
		t.Fatalf("exit: native=%d (%q)", nCode, nErr)
	}
	assertEnvelopeGolden(t, "merge-sentinel-displays-local", goldenEnvelope{
		stdout: normalize(nOut, nativeRoot), stderr: normalize(nErr, nativeRoot), exit: nCode,
	})
	if !strings.Contains(nOut, "abc1234 (local)") {
		t.Fatalf("sentinel should render as 'abc1234 (local)', got:\n%s", nOut)
	}
	if strings.Contains(nOut, "local-merge:abc1234") {
		t.Fatalf("raw sentinel value should not appear verbatim in the table:\n%s", nOut)
	}
}

// TestMergeLocalEntityActuallyAdvances (AC-1): confirms the merge: local terminal
// --set truly mutated the entity to done (not a vacuous exit-0).
func TestMergeLocalEntityActuallyAdvances(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "merge-local-workflow")
	_, errOut, code := runNative(t, root, env,
		"--workflow-dir", root, "--set", "020-no-sentinel", "status=done", "verdict=passed")
	if code != 0 {
		t.Fatalf("terminal --set should succeed, got %d (stderr=%q)", code, errOut)
	}
	fm := readWhole(t, filepath.Join(root, "020-no-sentinel.md"))
	if !strings.Contains(fm, "status: done") {
		t.Fatalf("entity should have advanced to done:\n%s", fm)
	}
}
