// ABOUTME: merge guard verb behavior — arm / blocked / finalize over the proven
// ABOUTME: --set/--archive paths, propagating (never bypassing) the guard refusals.
package status

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// driveMergeGuard runs `merge guard` against a fresh staged fixture, returning the
// staged root plus the three channels. The fixture is git-initialized (mutations
// resolve git_root), and --workflow-dir points at the staged copy so the verb
// drives the real --set/--archive paths against it.
func driveMergeGuard(t *testing.T, fixture string, args ...string) (root, stdout, stderr string, code int) {
	t.Helper()
	root = stageFixture(t, fixture)
	var out, errBuf bytes.Buffer
	full := append([]string{"--workflow-dir", root}, args...)
	code = MergeGuard(full, root, &out, &errBuf)
	return root, out.String(), errBuf.String(), code
}

// frontmatterField reads one top-level frontmatter field from a staged entity.
func frontmatterField(t *testing.T, path, field string) string {
	t.Helper()
	return strings.TrimSpace(ParseFrontmatter(path)[field])
}

// TestMergeGuardArmsUnderLocal (AC-1, Phase A): on a merge: local entity with no
// mod-block, the first guard run arms — it sets mod-block=merge:{hook}, signals
// `armed`/invoke-hook, exits 0, and does NOT terminalize or archive.
func TestMergeGuardArmsUnderLocal(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-local-workflow", "020-no-sentinel", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("arm should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "armed") {
		t.Fatalf("stdout should signal armed, got %q", out)
	}
	entity := filepath.Join(root, "020-no-sentinel.md")
	if got := frontmatterField(t, entity, "mod-block"); got != "merge:local-merge" {
		t.Fatalf("arm should set mod-block=merge:local-merge, got %q", got)
	}
	if got := frontmatterField(t, entity, "status"); got != "implementation" {
		t.Fatalf("arm must NOT terminalize, status=%q", got)
	}
	if fileExists(filepath.Join(root, "_archive", "020-no-sentinel.md")) {
		t.Fatal("arm must NOT archive")
	}
}

// TestMergeGuardAutoArmsUnderBothPolicies (AC-1): on the IDENTICAL precondition —
// empty mod-block, empty pr, `--verdict passed`, a merge hook registered — the verb
// AUTO-ARMS under BOTH policies. The merge: pr leg is the new behavior: entering
// terminal with an empty mod-block under merge: pr is the start of the merge
// ceremony, so the verb owns arming the mod-block (mod-block=merge:{hook}) and
// signaling the FO to invoke the hook — it no longer refuses at the merge-hook
// guard. The ceremony integrity is preserved downstream: the armed mod-block is
// what the FINALIZE step clears, and a merge: pr entity can only finalize once a
// merge sentinel records the landed PR (see the finalize-from-merged test).
func TestMergeGuardAutoArmsUnderBothPolicies(t *testing.T) {
	// merge: local leg — arms.
	localRoot, localOut, localErr, localCode := driveMergeGuard(t, "merge-local-workflow", "020-no-sentinel", "--verdict", "passed")
	if localCode != 0 {
		t.Fatalf("merge: local empty-mod-block --verdict passed must ARM (exit 0), got %d (stderr=%q)", localCode, localErr)
	}
	if !strings.Contains(localOut, "armed") {
		t.Fatalf("merge: local leg must signal armed, got %q", localOut)
	}
	if got := frontmatterField(t, filepath.Join(localRoot, "020-no-sentinel.md"), "mod-block"); got != "merge:local-merge" {
		t.Fatalf("merge: local arm must set mod-block, got %q", got)
	}

	// merge: pr leg — also arms. Same precondition, SAME outcome now (the fixture's
	// registered merge hook is named local-merge, so the arm value matches it; in
	// production the merge mod is pr-merge → mod-block=merge:pr-merge per AC-1).
	prRoot, prOut, prErr, prCode := driveMergeGuard(t, "merge-pr-workflow", "020-no-sentinel", "--verdict", "passed")
	if prCode != 0 {
		t.Fatalf("merge: pr empty-mod-block --verdict passed must AUTO-ARM (exit 0), got %d (stderr=%q)", prCode, prErr)
	}
	if !strings.Contains(prOut, "armed") {
		t.Fatalf("merge: pr leg must signal armed, got %q", prOut)
	}
	if got := frontmatterField(t, filepath.Join(prRoot, "020-no-sentinel.md"), "mod-block"); got != "merge:local-merge" {
		t.Fatalf("merge: pr auto-arm must set mod-block=merge:{hook}, got %q", got)
	}
	if got := frontmatterField(t, filepath.Join(prRoot, "020-no-sentinel.md"), "status"); got != "implementation" {
		t.Fatalf("merge: pr auto-arm must NOT terminalize, status=%q", got)
	}
	if fileExists(filepath.Join(prRoot, "_archive", "020-no-sentinel.md")) {
		t.Fatal("merge: pr auto-arm must NOT archive")
	}
}

// TestMergeGuardFinalizesFromMergedSentinelNonArmed (AC-2, the stranded case): an
// entity whose PR merged carries a merge sentinel (pr: pr-merge:99) and an EMPTY
// mod-block (a re-validation bounce cleared the block). `merge guard` must FINALIZE
// it — clear (no-op, already empty), terminalize, archive — from this non-armed
// state, keying off the merge sentinel rather than requiring an armed mod-block.
func TestMergeGuardFinalizesFromMergedSentinelNonArmed(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "080-pr-merged", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("merged-sentinel non-armed finalize should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "finalized") {
		t.Fatalf("stdout should signal finalized, got %q", out)
	}
	archived := filepath.Join(root, "_archive", "080-pr-merged.md")
	if !fileExists(archived) {
		t.Fatal("finalize should archive the merged entity")
	}
	if got := frontmatterField(t, archived, "status"); got != "done" {
		t.Fatalf("archived status=%q, want done", got)
	}
	if got := frontmatterField(t, archived, "verdict"); got != "passed" {
		t.Fatalf("archived verdict=%q, want passed", got)
	}
}

func TestMergeGuardFinalizesTerminalBlockedEntityFromMergedSentinel(t *testing.T) {
	root := stageFixture(t, "merge-pr-workflow")
	body, err := os.ReadFile(filepath.Join(root, "070-pr-pending.md"))
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.NewReplacer("status: implementation", "status: done", `pr: "#42"`, "pr: pr-merge:42", "mod-block: merge:pr-merge", "mod-block:").Replace(string(body))
	writeFile(t, filepath.Join(root, "070-pr-pending.md"), updated)
	var out, errOut bytes.Buffer
	if code := MergeGuard([]string{"--workflow-dir", root, "070-pr-pending", "--verdict", "passed"}, root, &out, &errOut); code != 0 {
		t.Fatalf("terminal merged-sentinel finalize exit=%d stderr=%q", code, errOut.String())
	}
	if archived := filepath.Join(root, "_archive", "070-pr-pending.md"); !fileExists(archived) || frontmatterField(t, archived, "status") != "done" || frontmatterField(t, archived, "mod-block") != "" {
		t.Fatalf("terminal merged-sentinel did not finalize through existing guard: %s", out.String())
	}
}

// TestMergeGuardBlocksOnOpenPRNoModBlock (the premature-finalize gate): a bare,
// OPEN PR reference (pr: #42) with an EMPTY mod-block must signal blocked/await-pr
// and must NOT finalize or archive. The verb must NEVER finalize on pr-presence
// alone — only a merge sentinel (pr-merge:/local-merge:) finalizes. This pins the
// bug where an open-PR entity was archived before its PR landed.
func TestMergeGuardBlocksOnOpenPRNoModBlock(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "090-pr-open-unmerged", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("open-PR blocked should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "blocked") {
		t.Fatalf("open-PR with no mod-block must signal blocked, got %q", out)
	}
	if strings.Contains(out, "finalized") {
		t.Fatalf("open-PR must NOT finalize on pr-presence alone, got %q", out)
	}
	entity := filepath.Join(root, "090-pr-open-unmerged.md")
	if got := frontmatterField(t, entity, "status"); got != "implementation" {
		t.Fatalf("open-PR blocked must NOT terminalize, status=%q", got)
	}
	if fileExists(filepath.Join(root, "_archive", "090-pr-open-unmerged.md")) {
		t.Fatal("open-PR blocked must NOT archive (premature-finalize bug)")
	}
}

// TestMergeGuardDoesNotFinalizeOnMalformedSentinel (the fail-open gate): a sentinel
// whose prefix matches but whose suffix is garbage (pr: pr-merge:abc) must NOT
// finalize. A bare HasPrefix match would treat `abc` as a landed merge and drive a
// full finalize+archive — a fail-OPEN hole. The verb must instead leave the entity
// at its non-terminal status, unarchived; only a well-formed sentinel finalizes.
func TestMergeGuardDoesNotFinalizeOnMalformedSentinel(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "100-pr-malformed-sentinel", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("malformed-sentinel guard should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if strings.Contains(out, "finalized") {
		t.Fatalf("malformed sentinel (pr-merge:abc) must NOT finalize, got %q", out)
	}
	entity := filepath.Join(root, "100-pr-malformed-sentinel.md")
	if got := frontmatterField(t, entity, "status"); got != "implementation" {
		t.Fatalf("malformed sentinel must NOT terminalize, status=%q", got)
	}
	if fileExists(filepath.Join(root, "_archive", "100-pr-malformed-sentinel.md")) {
		t.Fatal("malformed sentinel must NOT archive (fail-open hole)")
	}
}

// TestPRIndicatesMerged pins the suffix-validation contract: a sentinel finalizes
// ONLY when well-formed — a pr-merge: suffix that parses as a positive integer, or
// a local-merge: suffix that is a non-empty SHA-like token. Every garbage form (a
// non-numeric or zero pr-merge suffix, an empty/quoted-empty suffix, a bare PR
// reference) must return false, the fail-CLOSED direction.
func TestPRIndicatesMerged(t *testing.T) {
	cases := []struct {
		pr   string
		want bool
	}{
		// Well-formed — these finalize.
		{"pr-merge:42", true},
		{"pr-merge:99", true},
		{"local-merge:abc1234", true},
		{"local-merge:0a1b2c3", true},
		// Malformed pr-merge suffix — must NOT finalize.
		{"pr-merge:abc", false},
		{"pr-merge:0", false},
		{"pr-merge:-1", false},
		{"pr-merge:12x", false},
		{"pr-merge:", false},
		{`"pr-merge:"`, false},
		// Malformed local-merge suffix — must NOT finalize.
		{"local-merge:", false},
		{`"local-merge:"`, false},
		// Bare / open PR references and empties — never finalize.
		{"", false},
		{"#42", false},
		{"owner/repo#42", false},
		{"https://github.com/o/r/pull/42", false},
	}
	for _, c := range cases {
		if got := prIndicatesMerged(c.pr); got != c.want {
			t.Errorf("prIndicatesMerged(%q) = %v, want %v", c.pr, got, c.want)
		}
	}
}

// TestMergeGuardArmThenFinalizeLocal (AC-1, the happy path): arm, then a second
// run finalizes — the entity ends terminal+passed, archived, mod-block cleared.
func TestMergeGuardArmThenFinalizeLocal(t *testing.T) {
	root := stageFixture(t, "merge-local-workflow")
	entity := filepath.Join(root, "020-no-sentinel.md")

	var out1, err1 bytes.Buffer
	if code := MergeGuard([]string{"--workflow-dir", root, "020-no-sentinel", "--verdict", "passed"}, root, &out1, &err1); code != 0 {
		t.Fatalf("arm exit=%d stderr=%q", code, err1.String())
	}
	if got := frontmatterField(t, entity, "mod-block"); got != "merge:local-merge" {
		t.Fatalf("post-arm mod-block=%q", got)
	}

	var out2, err2 bytes.Buffer
	code := MergeGuard([]string{"--workflow-dir", root, "020-no-sentinel", "--verdict", "passed"}, root, &out2, &err2)
	if code != 0 {
		t.Fatalf("finalize should exit 0, got %d (stderr=%q)", code, err2.String())
	}
	if !strings.Contains(out2.String(), "finalized") {
		t.Fatalf("stdout should signal finalized, got %q", out2.String())
	}
	archived := filepath.Join(root, "_archive", "020-no-sentinel.md")
	if !fileExists(archived) {
		t.Fatal("finalize should archive the entity")
	}
	if got := frontmatterField(t, archived, "status"); got != "done" {
		t.Fatalf("archived status=%q, want done", got)
	}
	if got := frontmatterField(t, archived, "verdict"); got != "passed" {
		t.Fatalf("archived verdict=%q, want passed", got)
	}
	if got := frontmatterField(t, archived, "mod-block"); got != "" {
		t.Fatalf("archived mod-block should be cleared, got %q", got)
	}
}

// TestMergeGuardBlockedOnPR (AC-3): with pr set, the verb signals blocked/await-pr,
// leaves mod-block + pr + non-terminal status intact, and does not archive.
func TestMergeGuardBlockedOnPR(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "070-pr-pending", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("blocked should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "blocked") {
		t.Fatalf("stdout should signal blocked, got %q", out)
	}
	entity := filepath.Join(root, "070-pr-pending.md")
	if got := frontmatterField(t, entity, "mod-block"); got != "merge:local-merge" {
		t.Fatalf("blocked must leave mod-block intact, got %q", got)
	}
	if got := frontmatterField(t, entity, "pr"); got != "#42" {
		t.Fatalf("blocked must leave pr intact, got %q", got)
	}
	if got := frontmatterField(t, entity, "status"); got != "implementation" {
		t.Fatalf("blocked must NOT terminalize, status=%q", got)
	}
	if fileExists(filepath.Join(root, "_archive", "070-pr-pending.md")) {
		t.Fatal("blocked must NOT archive")
	}
}

// TestMergeGuardVerdictOmissionUnreachable (AC-4): `merge guard` without --verdict
// exits non-zero BEFORE any terminal write — a verdict-less finalize is unreachable.
func TestMergeGuardVerdictOmissionUnreachable(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-local-workflow", "020-no-sentinel")
	if code == 0 {
		t.Fatalf("verdict-less guard must exit non-zero, got 0 (stdout=%q)", out)
	}
	if !strings.Contains(errOut, "verdict") {
		t.Fatalf("stderr should name the missing verdict, got %q", errOut)
	}
	// No terminal write: the entity is untouched (no mod-block armed either, since
	// the verdict gate fires before resolving the workflow).
	if got := frontmatterField(t, filepath.Join(root, "020-no-sentinel.md"), "status"); got != "implementation" {
		t.Fatalf("verdict-less guard must not mutate, status=%q", got)
	}
}

// TestMergeGuardRejectsBadVerdict (AC-4 companion): a --verdict value other than
// passed/rejected is a usage error, before any mutation.
func TestMergeGuardRejectsBadVerdict(t *testing.T) {
	_, out, errOut, code := driveMergeGuard(t, "merge-local-workflow", "020-no-sentinel", "--verdict", "maybe")
	if code != 1 {
		t.Fatalf("bad verdict must exit 1, got %d (stdout=%q)", code, out)
	}
	if !strings.Contains(errOut, "passed") || !strings.Contains(errOut, "rejected") {
		t.Fatalf("stderr should name the valid verdicts, got %q", errOut)
	}
}

// TestMergeGuardArmedButNoMergeRefusesFinalize (AC-5, the guard-not-bypassed pin):
// under merge: pr the verb auto-arms an empty-mod-block entity, but it must NEVER
// bypass the merge-hook guard to wrongly finalize an entity that has not actually
// merged. After auto-arm, if the FO re-runs the verb WITHOUT invoking the hook (no
// PR opened, no merge sentinel), the finalize step's terminalize --set hits the
// merge-hook guard and the verb propagates its `cannot advance to terminal` refusal
// (exit 1) verbatim — it never passes --force. This is the AC-5 spirit (the guard
// is propagated, never bypassed) carried into the new auto-arm contract: arming is
// safe because the merge sentinel — not arming — is what unlocks finalize.
func TestMergeGuardArmedButNoMergeRefusesFinalize(t *testing.T) {
	root := stageFixture(t, "merge-pr-workflow")
	entity := filepath.Join(root, "020-no-sentinel.md")

	// First run auto-arms.
	var out1, err1 bytes.Buffer
	if code := MergeGuard([]string{"--workflow-dir", root, "020-no-sentinel", "--verdict", "passed"}, root, &out1, &err1); code != 0 {
		t.Fatalf("auto-arm exit=%d stderr=%q", code, err1.String())
	}
	if got := frontmatterField(t, entity, "mod-block"); got != "merge:local-merge" {
		t.Fatalf("auto-arm must set mod-block, got %q", got)
	}

	// Second run with the hook NOT invoked (no pr, no sentinel): finalize must refuse
	// at the merge-hook guard, NOT bypass it. The mod-block clears (standalone --set),
	// then the terminalize --set refuses — so the entity ends with an empty mod-block
	// but is NOT terminalized or archived.
	var out2, err2 bytes.Buffer
	code := MergeGuard([]string{"--workflow-dir", root, "020-no-sentinel", "--verdict", "passed"}, root, &out2, &err2)
	if code != 1 {
		t.Fatalf("armed-but-no-merge finalize must refuse (exit 1), got %d (stdout=%q)", code, out2.String())
	}
	if !strings.Contains(err2.String(), "cannot advance to terminal") {
		t.Fatalf("stderr should propagate the merge-hook refusal verbatim, got %q", err2.String())
	}
	if strings.Contains(out2.String(), "finalized") {
		t.Fatalf("stdout must not claim finalized on a refusal, got %q", out2.String())
	}
	if got := frontmatterField(t, entity, "status"); got != "implementation" {
		t.Fatalf("refused entity must not terminalize, status=%q", got)
	}
	if fileExists(filepath.Join(root, "_archive", "020-no-sentinel.md")) {
		t.Fatal("refused entity must NOT archive")
	}
}

// TestMergeGuardRejectedFinalizesNoPR (AC-6a): --verdict rejected finalizes an
// entity with empty pr to terminal+rejected+archived without --force.
func TestMergeGuardRejectedFinalizesNoPR(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "040-rejected", "--verdict", "rejected")
	if code != 0 {
		t.Fatalf("rejected finalize should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "finalized") {
		t.Fatalf("stdout should signal finalized, got %q", out)
	}
	archived := filepath.Join(root, "_archive", "040-rejected.md")
	if got := frontmatterField(t, archived, "status"); got != "done" {
		t.Fatalf("archived status=%q, want done", got)
	}
	if got := frontmatterField(t, archived, "verdict"); got != "rejected" {
		t.Fatalf("archived verdict=%q, want rejected", got)
	}
}

// TestMergeGuardRejectedClearsModBlockFirst (AC-6b): --verdict rejected on an
// entity with an in-flight mod-block clears the block in its own --set FIRST, then
// terminalizes+archives — the combined-clear-with-terminal guard never trips
// because the verb separates the steps.
func TestMergeGuardRejectedClearsModBlockFirst(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "050-rejected-pending", "--verdict", "rejected")
	if code != 0 {
		t.Fatalf("rejected-with-mod-block finalize should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "finalized") {
		t.Fatalf("stdout should signal finalized, got %q", out)
	}
	if strings.Contains(errOut, "combined mod-block clear with terminal transition") {
		t.Fatalf("verb must separate the clear from the terminalize, not trip the combined guard: %q", errOut)
	}
	archived := filepath.Join(root, "_archive", "050-rejected-pending.md")
	if got := frontmatterField(t, archived, "status"); got != "done" {
		t.Fatalf("archived status=%q, want done", got)
	}
	if got := frontmatterField(t, archived, "verdict"); got != "rejected" {
		t.Fatalf("archived verdict=%q, want rejected", got)
	}
	if got := frontmatterField(t, archived, "mod-block"); got != "" {
		t.Fatalf("archived mod-block should be cleared, got %q", got)
	}
}

// TestMergeGuardRefusesMissingMergeModNoSentinel (AC-5 row a, D5): a mod-block
// naming a merge mod that is no longer registered under `_mods/`, with no merge
// sentinel and a non-rejected verdict, must REFUSE — not silently finalize. This
// pins the bug found by code read (merge.go's pre-D5 default case): clearing the
// block here would archive the entity without its hook ever having run. The
// refusal exits 1, names the missing mod, and mutates nothing (status, mod-block,
// and pr all survive byte-identical).
func TestMergeGuardRefusesMissingMergeModNoSentinel(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "120-missing-mod-no-sentinel", "--verdict", "passed")
	if code != 1 {
		t.Fatalf("missing-mod no-sentinel must refuse (exit 1), got %d (stdout=%q)", code, out)
	}
	if !strings.Contains(errOut, "ghost-merge") {
		t.Fatalf("stderr should name the missing mod, got %q", errOut)
	}
	if !strings.Contains(errOut, "is missing") {
		t.Fatalf("stderr should say the mod is missing, got %q", errOut)
	}
	// Feedback cycle 1, T2: pin the full remediation tail, not just the mod-name +
	// "is missing" fragment — a change that drops the restore/--force remedy must
	// fail this test.
	if !strings.Contains(errOut, "Restore the mod file, or have the operator clear the block with --force.") {
		t.Fatalf("stderr should carry the full remediation tail, got %q", errOut)
	}
	if strings.Contains(out, "finalized") {
		t.Fatalf("stdout must not claim finalized on the missing-mod refusal, got %q", out)
	}
	entity := filepath.Join(root, "120-missing-mod-no-sentinel.md")
	if got := frontmatterField(t, entity, "status"); got != "implementation" {
		t.Fatalf("refused entity must not terminalize, status=%q", got)
	}
	if got := frontmatterField(t, entity, "mod-block"); got != "merge:ghost-merge" {
		t.Fatalf("refused entity must keep its mod-block intact, got %q", got)
	}
	if fileExists(filepath.Join(root, "_archive", "120-missing-mod-no-sentinel.md")) {
		t.Fatal("refused entity must NOT archive")
	}
}

// TestMergeGuardRefusesMissingMergeModWhenNoHookRegisteredAtAll (feedback cycle 1,
// T1): the likeliest real-world D5 trigger — the deleted mod file WAS the
// workflow's only registered merge hook, so `mergeHooks` is EMPTY, not merely
// missing the named entry. This is a distinct shape from
// TestMergeGuardRefusesMissingMergeModNoSentinel (which fixtures a workflow with
// a DIFFERENT hook still registered): a mutant that special-cases
// modBlockNamesMissingMergeMod on `len(mergeHooks)==0` (treating "no hooks at
// all" as vacuously not-missing) passes every other D5 test but must fail here.
func TestMergeGuardRefusesMissingMergeModWhenNoHookRegisteredAtAll(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-no-hook-workflow", "030-missing-mod-no-hooks-registered", "--verdict", "passed")
	if code != 1 {
		t.Fatalf("missing-mod refusal must fire even with zero hooks registered (exit 1), got %d (stdout=%q)", code, out)
	}
	if !strings.Contains(errOut, "ghost-merge") {
		t.Fatalf("stderr should name the missing mod, got %q", errOut)
	}
	if !strings.Contains(errOut, "is missing") {
		t.Fatalf("stderr should say the mod is missing, got %q", errOut)
	}
	if strings.Contains(out, "finalized") {
		t.Fatalf("stdout must not claim finalized — this is the exact silent-finalize shape D5 fixes, got %q", out)
	}
	entity := filepath.Join(root, "030-missing-mod-no-hooks-registered.md")
	if got := frontmatterField(t, entity, "status"); got != "implementation" {
		t.Fatalf("refused entity must not terminalize, status=%q", got)
	}
	if got := frontmatterField(t, entity, "mod-block"); got != "merge:ghost-merge" {
		t.Fatalf("refused entity must keep its mod-block intact, got %q", got)
	}
	if fileExists(filepath.Join(root, "_archive", "030-missing-mod-no-hooks-registered.md")) {
		t.Fatal("refused entity must NOT archive")
	}
}

// TestMergeGuardFinalizesMissingMergeModWithSentinel (AC-5 row b): a mod-block
// naming a missing merge mod does NOT refuse when a well-formed merge sentinel is
// already recorded — the sentinel honestly proves the ceremony ran before the mod
// file was deleted, so finalize proceeds exactly as it does today.
func TestMergeGuardFinalizesMissingMergeModWithSentinel(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "130-missing-mod-with-sentinel", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("missing-mod WITH sentinel should still finalize (exit 0), got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "finalized") {
		t.Fatalf("stdout should signal finalized, got %q", out)
	}
	archived := filepath.Join(root, "_archive", "130-missing-mod-with-sentinel.md")
	if !fileExists(archived) {
		t.Fatal("finalize should archive the entity")
	}
	if got := frontmatterField(t, archived, "status"); got != "done" {
		t.Fatalf("archived status=%q, want done", got)
	}
}

// TestMergeGuardFinalizesMissingMergeModRejected (AC-5 row c): a rejected verdict
// finalizes even when the mod-block names a missing merge mod — the entity never
// merged, so the rejected-verdict escape takes priority over the missing-mod
// refusal, exactly as it does today.
func TestMergeGuardFinalizesMissingMergeModRejected(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "140-missing-mod-rejected", "--verdict", "rejected")
	if code != 0 {
		t.Fatalf("missing-mod rejected should still finalize (exit 0), got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "finalized") {
		t.Fatalf("stdout should signal finalized, got %q", out)
	}
	archived := filepath.Join(root, "_archive", "140-missing-mod-rejected.md")
	if got := frontmatterField(t, archived, "verdict"); got != "rejected" {
		t.Fatalf("archived verdict=%q, want rejected", got)
	}
}

// TestMergeGuardArmedLineNamesNextStep (AC-3, D3): the armed phase's default
// prose names the hook's file path and the never-invoked-by-the-verb caveat, plus
// the re-run command — the FO's next action, carried at fire time.
func TestMergeGuardArmedLineNamesNextStep(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-local-workflow", "020-no-sentinel", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("arm should exit 0, got %d (stderr=%q)", code, errOut)
	}
	want := fmt.Sprintf(
		"armed: mod-block set to merge:local-merge — invoke the local-merge merge hook (%s/_mods/local-merge.md; merge guard never invokes it), then re-run `merge guard 020-no-sentinel`.\n",
		root)
	if out != want {
		t.Fatalf("armed line = %q, want %q", out, want)
	}
}

// TestMergeGuardBlockedLineNamesNextStep (AC-3, D3): the blocked phase's default
// prose names the never-finalize-on-open-PR invariant and the sentinel format
// that unlocks finalize, plus the re-run command.
func TestMergeGuardBlockedLineNamesNextStep(t *testing.T) {
	_, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "070-pr-pending", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("blocked should exit 0, got %d (stderr=%q)", code, errOut)
	}
	want := "blocked: PR #42 is pending — mod-block left intact, never finalize on an open PR. " +
		"When gh reports it MERGED, record the sentinel (pr=pr-merge:{number}) and re-run `merge guard 070-pr-pending`.\n"
	if out != want {
		t.Fatalf("blocked line = %q, want %q", out, want)
	}
}

// TestMergeGuardFinalizedLineNoWorktreeNoHookClause (AC-3, D3): an entity
// finalizing via a merged sentinel with NO worktree recorded gets the base
// finalized line only — no worktree-removal clause (nothing to remove) and no
// no-merge-hook clause (a hook IS registered and a sentinel IS recorded).
func TestMergeGuardFinalizedLineNoWorktreeNoHookClause(t *testing.T) {
	_, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "080-pr-merged", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("finalize should exit 0, got %d (stderr=%q)", code, errOut)
	}
	want := "finalized: 080-pr-merged -> done (verdict passed), archived.\n"
	if out != want {
		t.Fatalf("finalized line = %q, want %q (no worktree/no-hook clauses)", out, want)
	}
}

// TestMergeGuardFinalizedLineWithWorktree (AC-3, D3): an entity finalizing via a
// merged sentinel WITH a recorded worktree gets the worktree-removal/branch-
// cleanup/teardown next-step clause, and — since a hook IS registered and a
// sentinel IS recorded — no no-merge-hook clause.
func TestMergeGuardFinalizedLineWithWorktree(t *testing.T) {
	_, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "150-worktree-finalize", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("finalize should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "finalized: 150-worktree-finalize -> done (verdict passed), archived.") {
		t.Fatalf("stdout should carry the base finalized line, got %q", out)
	}
	if !strings.Contains(out, "Next: push; remove the worktree (`git worktree remove .worktrees/150-worktree-finalize`") {
		t.Fatalf("stdout should carry the worktree-removal next-step clause, got %q", out)
	}
	if !strings.Contains(out, "delete the local branch (`git branch -d`)") {
		t.Fatalf("stdout should carry the branch-cleanup clause, got %q", out)
	}
	if !strings.Contains(out, "keep the remote branch while a PR references it") {
		t.Fatalf("stdout should carry the remote-branch-retention clause, got %q", out)
	}
	if !strings.Contains(out, "tear down the entity's workers per your runtime adapter") {
		t.Fatalf("stdout should carry the worker-teardown clause, got %q", out)
	}
	if strings.Contains(out, "no merge hook registered") {
		t.Fatalf("stdout must NOT carry the no-merge-hook clause (a hook IS registered and a sentinel IS recorded), got %q", out)
	}
}

// TestMergeGuardFinalizedLineNoHookRegisteredNoWorktree (AC-3, D3): under a
// workflow with NO merge hook registered at all, a Phase-C default finalize (no
// pr, no mod-block, no worktree) names the manual `--no-ff` merge onto trunk —
// nothing automated it — with no worktree-removal clause.
func TestMergeGuardFinalizedLineNoHookRegisteredNoWorktree(t *testing.T) {
	_, out, errOut, code := driveMergeGuard(t, "merge-no-hook-workflow", "010-no-hook-no-worktree", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("finalize should exit 0, got %d (stderr=%q)", code, errOut)
	}
	want := "finalized: 010-no-hook-no-worktree -> done (verdict passed), archived.\n" +
		"no merge hook registered — merge the stage branch onto main with --no-ff if not already merged.\n"
	if out != want {
		t.Fatalf("finalized output = %q, want %q", out, want)
	}
}

// TestMergeGuardFinalizedLineNoHookRegisteredWithWorktree (AC-3, D3): the same
// no-hook-registered path but WITH a recorded worktree carries BOTH next-step
// clauses — the two conditions are independent.
func TestMergeGuardFinalizedLineNoHookRegisteredWithWorktree(t *testing.T) {
	_, out, errOut, code := driveMergeGuard(t, "merge-no-hook-workflow", "020-no-hook-with-worktree", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("finalize should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if !strings.Contains(out, "Next: push; remove the worktree (`git worktree remove .worktrees/020-no-hook-with-worktree`") {
		t.Fatalf("stdout should carry the worktree-removal clause, got %q", out)
	}
	if !strings.Contains(out, "no merge hook registered — merge the stage branch onto main with --no-ff if not already merged.") {
		t.Fatalf("stdout should carry the no-merge-hook clause, got %q", out)
	}
}

// TestMergeGuardRequiresSlug (usage): a guard call with no slug is a usage error.
func TestMergeGuardRequiresSlug(t *testing.T) {
	root := stageFixture(t, "merge-local-workflow")
	var out, errBuf bytes.Buffer
	code := MergeGuard([]string{"--workflow-dir", root, "--verdict", "passed"}, root, &out, &errBuf)
	if code != 1 {
		t.Fatalf("missing slug must exit 1, got %d", code)
	}
	if !strings.Contains(errBuf.String(), "slug") {
		t.Fatalf("stderr should name the missing slug, got %q", errBuf.String())
	}
}

// gitOutput runs git in dir and returns trimmed stdout, failing the test on error.
func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// TestMergeGuardFinalizeCommitsArchivePathScoped (AC-3): the finalize archive move
// is committed BY THE VERB, path-scoped — the rename is staged correctly (old path
// deleted, new _archive path added) and a sibling entity left dirty in the same
// tree is NOT swept into the commit. This pins the toil the verb eliminates: the
// FO no longer hand-commits the rename (the newline-in-variable staging bug that
// failed the archive commit twice).
func TestMergeGuardFinalizeCommitsArchivePathScoped(t *testing.T) {
	root := stageFixture(t, "merge-local-workflow")

	// Dirty a SIBLING entity so a bare `git add -A` would sweep it into the verb's
	// archive commit. The path-scoped commit must leave it unstaged/uncommitted.
	sibling := filepath.Join(root, "030-pending.md")
	f, err := os.OpenFile(sibling, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open sibling: %v", err)
	}
	if _, err := f.WriteString("\nsibling edit that must NOT be swept into the archive commit\n"); err != nil {
		t.Fatalf("dirty sibling: %v", err)
	}
	f.Close()

	baseHead := gitOutput(t, root, "rev-parse", "HEAD")

	// Arm, then finalize 020-no-sentinel (merge: local empty mod-block → arm → finalize).
	var out1, err1 bytes.Buffer
	if code := MergeGuard([]string{"--workflow-dir", root, "020-no-sentinel", "--verdict", "passed"}, root, &out1, &err1); code != 0 {
		t.Fatalf("arm exit=%d stderr=%q", code, err1.String())
	}
	var out2, err2 bytes.Buffer
	if code := MergeGuard([]string{"--workflow-dir", root, "020-no-sentinel", "--verdict", "passed"}, root, &out2, &err2); code != 0 {
		t.Fatalf("finalize exit=%d stderr=%q", code, err2.String())
	}

	// The verb committed: HEAD advanced past the seed commit.
	if head := gitOutput(t, root, "rev-parse", "HEAD"); head == baseHead {
		t.Fatalf("finalize must commit the archive move (HEAD unchanged at %s)", baseHead)
	}

	// The HEAD commit touches ONLY the entity's two paths (the rename) — never the
	// dirty sibling. git records the rename as the old + new path.
	changed := gitOutput(t, root, "show", "--name-only", "--format=", "HEAD")
	for _, line := range strings.Split(changed, "\n") {
		p := strings.TrimSpace(line)
		if p == "" {
			continue
		}
		if !strings.Contains(p, "020-no-sentinel") {
			t.Fatalf("archive commit must be path-scoped to the entity, but also touched %q (full set:\n%s)", p, changed)
		}
	}
	if !strings.Contains(changed, "_archive/020-no-sentinel.md") {
		t.Fatalf("archive commit must add the _archive path, got:\n%s", changed)
	}

	// The sibling edit survives uncommitted — it was NOT swept.
	porcelain := gitOutput(t, root, "status", "--porcelain", "030-pending.md")
	if !strings.Contains(porcelain, "030-pending.md") {
		t.Fatalf("sibling edit must remain uncommitted (a bare add -A would have swept it), git status:\n%q", porcelain)
	}

	// The working tree no longer has the source file at its old path; it lives in _archive.
	if fileExists(filepath.Join(root, "020-no-sentinel.md")) {
		t.Fatal("source entity must be moved out of the workflow root")
	}
	if !fileExists(filepath.Join(root, "_archive", "020-no-sentinel.md")) {
		t.Fatal("entity must land in _archive")
	}
}

// TestMergeGuardFinalizesFolderFormEntity (FIX 1): a FOLDER-FORM entity
// ({slug}/index.md) finalizes — terminalize + archive the whole folder to
// _archive/{slug}/ — and the archive move is committed PATH-SCOPED by the verb. This
// pins the bug where commitArchiveMove hardcoded the flat {slug}.md paths and
// exit-128'd the `git add` on a folder-form entity, stranding it. A dirty sibling
// must NOT be swept into the path-scoped commit.
func TestMergeGuardFinalizesFolderFormEntity(t *testing.T) {
	root := stageFixture(t, "merge-pr-workflow")

	// Dirty a sibling so a bare `git add -A` would sweep it in.
	sibling := filepath.Join(root, "070-pr-pending.md")
	f, err := os.OpenFile(sibling, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open sibling: %v", err)
	}
	if _, err := f.WriteString("\nsibling edit that must NOT be swept into the archive commit\n"); err != nil {
		t.Fatalf("dirty sibling: %v", err)
	}
	f.Close()

	baseHead := gitOutput(t, root, "rev-parse", "HEAD")

	var out, errBuf bytes.Buffer
	code := MergeGuard([]string{"--workflow-dir", root, "110-pr-merged-folder", "--verdict", "passed"}, root, &out, &errBuf)
	if code != 0 {
		t.Fatalf("folder-form finalize should exit 0, got %d (stderr=%q)", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "finalized") {
		t.Fatalf("stdout should signal finalized, got %q", out.String())
	}

	// The folder moved into _archive/{slug}/index.md, terminal + verdict recorded.
	archived := filepath.Join(root, "_archive", "110-pr-merged-folder", "index.md")
	if !fileExists(archived) {
		t.Fatal("folder-form entity must archive to _archive/{slug}/index.md")
	}
	if got := frontmatterField(t, archived, "status"); got != "done" {
		t.Fatalf("archived status=%q, want done", got)
	}
	if got := frontmatterField(t, archived, "verdict"); got != "passed" {
		t.Fatalf("archived verdict=%q, want passed", got)
	}
	// The source folder is gone from the live root.
	if fileExists(filepath.Join(root, "110-pr-merged-folder", "index.md")) {
		t.Fatal("source folder must be moved out of the workflow root")
	}

	// The verb committed, and the commit touches ONLY the entity's folder paths.
	if head := gitOutput(t, root, "rev-parse", "HEAD"); head == baseHead {
		t.Fatalf("folder-form finalize must commit the archive move (HEAD unchanged at %s)", baseHead)
	}
	changed := gitOutput(t, root, "show", "--name-only", "--format=", "HEAD")
	for _, line := range strings.Split(changed, "\n") {
		p := strings.TrimSpace(line)
		if p == "" {
			continue
		}
		if !strings.Contains(p, "110-pr-merged-folder") {
			t.Fatalf("archive commit must be path-scoped to the folder, but also touched %q (full set:\n%s)", p, changed)
		}
	}
	if !strings.Contains(changed, "_archive/110-pr-merged-folder/index.md") {
		t.Fatalf("archive commit must add the _archive folder path, got:\n%s", changed)
	}
	// The sibling edit survives uncommitted — it was NOT swept.
	porcelain := gitOutput(t, root, "status", "--porcelain", "070-pr-pending.md")
	if !strings.Contains(porcelain, "070-pr-pending.md") {
		t.Fatalf("sibling edit must remain uncommitted, git status:\n%q", porcelain)
	}
}

// installFailingPreCommitHook writes a pre-commit hook into the staged repo that
// always exits 1, so the verb's archive commit fails — exercising the rollback path
// with a real git failure rather than a mock.
func installFailingPreCommitHook(t *testing.T, root string) {
	t.Helper()
	hookDir := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hookDir, 0o755); err != nil {
		t.Fatalf("mkdir hooks: %v", err)
	}
	hook := filepath.Join(hookDir, "pre-commit")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatalf("write hook: %v", err)
	}
}

// removePreCommitHook deletes the failing pre-commit hook so a recovery re-run can
// commit cleanly.
func removePreCommitHook(t *testing.T, root string) {
	t.Helper()
	if err := os.Remove(filepath.Join(root, ".git", "hooks", "pre-commit")); err != nil {
		t.Fatalf("remove hook: %v", err)
	}
}

// TestMergeGuardFinalizeRollsBackOnCommitFailure (FIX 3): finalize mutates disk in
// order — clear mod-block, terminalize, archive (rename), commit LAST. The commit
// stages the rename (`git add`) BEFORE it runs, so a failed commit (here a failing
// pre-commit hook) must roll back BOTH the working tree AND the git index. The
// rollback is proven by its END, not a worktree proxy: after it, (a) the working tree
// is byte-identical to pre-finalize and (b) the index is empty (no leaked staged
// rename to _archive that a later commit would sweep into HEAD, losing the live
// entity) and (c) RECOVERY works — a re-run with the commit now passing SUCCEEDS and
// archives, rather than exit-128ing on the leaked index. Asserted across BOTH flat
// and folder form.
func TestMergeGuardFinalizeRollsBackOnCommitFailure(t *testing.T) {
	cases := []struct {
		name       string
		slug       string
		liveRel    string // entity path relative to the workflow root, pre-archive
		archiveRel string // where it would land in _archive
	}{
		{"flat-form", "080-pr-merged", "080-pr-merged.md", "_archive/080-pr-merged.md"},
		{"folder-form", "110-pr-merged-folder", "110-pr-merged-folder/index.md", "_archive/110-pr-merged-folder/index.md"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := stageFixture(t, "merge-pr-workflow")
			installFailingPreCommitHook(t, root)

			entity := filepath.Join(root, tc.liveRel)
			preStatus := frontmatterField(t, entity, "status")
			prePR := frontmatterField(t, entity, "pr")
			preModBlock := frontmatterField(t, entity, "mod-block")
			preBytes, err := os.ReadFile(entity)
			if err != nil {
				t.Fatalf("read pre-state: %v", err)
			}

			var out, errBuf bytes.Buffer
			code := MergeGuard([]string{"--workflow-dir", root, tc.slug, "--verdict", "passed"}, root, &out, &errBuf)
			if code != 1 {
				t.Fatalf("a failing archive commit must exit 1, got %d (stdout=%q)", code, out.String())
			}
			if strings.Contains(out.String(), "finalized") {
				t.Fatalf("verb must not claim finalized on a commit failure, got %q", out.String())
			}

			// The entity is back at its live location — NOT in _archive.
			if !fileExists(entity) {
				t.Fatal("rollback must restore the entity to its live location")
			}
			if fileExists(filepath.Join(root, tc.archiveRel)) {
				t.Fatal("rollback must remove the entity from _archive")
			}
			// Frontmatter restored to pre-finalize values, no `archived` stamp leaked.
			if got := frontmatterField(t, entity, "status"); got != preStatus {
				t.Fatalf("rollback must restore status to %q, got %q", preStatus, got)
			}
			if got := frontmatterField(t, entity, "pr"); got != prePR {
				t.Fatalf("rollback must restore pr to %q, got %q", prePR, got)
			}
			if got := frontmatterField(t, entity, "mod-block"); got != preModBlock {
				t.Fatalf("rollback must restore mod-block to %q, got %q", preModBlock, got)
			}
			if got := frontmatterField(t, entity, "archived"); got != "" {
				t.Fatalf("rollback must not leave an archived stamp, got %q", got)
			}
			// Byte-for-byte identical to the pre-finalize file.
			postBytes, err := os.ReadFile(entity)
			if err != nil {
				t.Fatalf("read post-state: %v", err)
			}
			if string(postBytes) != string(preBytes) {
				t.Fatalf("rollback must restore the entity byte-for-byte\n--- pre ---\n%s\n--- post ---\n%s", preBytes, postBytes)
			}

			// The git INDEX is clean — no leaked staged rename to _archive. A non-empty
			// `git diff --cached` here is the entity-LOSS path: a later plain commit would
			// sweep the phantom rename into HEAD, leaving the live file an untracked orphan.
			if staged := gitOutput(t, root, "diff", "--cached", "--name-only"); staged != "" {
				t.Fatalf("rollback must leave the index CLEAN, but staged paths remain:\n%s", staged)
			}
			// The working tree is clean too (the seed commit + rollback leave nothing dirty).
			if porcelain := gitOutput(t, root, "status", "--porcelain"); porcelain != "" {
				t.Fatalf("rollback must leave the working tree CLEAN, git status:\n%s", porcelain)
			}

			// RECOVERY: with the commit now passing, a re-run finalizes — it does NOT
			// exit-128 on a leaked index, and the entity archives.
			removePreCommitHook(t, root)
			var rout, rerr bytes.Buffer
			rcode := MergeGuard([]string{"--workflow-dir", root, tc.slug, "--verdict", "passed"}, root, &rout, &rerr)
			if rcode != 0 {
				t.Fatalf("recovery re-run after rollback must SUCCEED (exit 0), got %d (stderr=%q)", rcode, rerr.String())
			}
			if !strings.Contains(rout.String(), "finalized") {
				t.Fatalf("recovery re-run must signal finalized, got %q", rout.String())
			}
			if !fileExists(filepath.Join(root, tc.archiveRel)) {
				t.Fatal("recovery re-run must archive the entity")
			}
			if fileExists(entity) {
				t.Fatal("recovery re-run must move the entity out of the live location")
			}
		})
	}
}

// TestMergeGuardRollbackPreservesFileMode (FIX 3 robustness): rollback restores the
// source file's ORIGINAL mode, not a hardcoded 0o644. A non-default mode set before
// the failed finalize survives the rollback.
func TestMergeGuardRollbackPreservesFileMode(t *testing.T) {
	root := stageFixture(t, "merge-pr-workflow")
	installFailingPreCommitHook(t, root)

	entity := filepath.Join(root, "080-pr-merged.md")
	if err := os.Chmod(entity, 0o600); err != nil {
		t.Fatalf("chmod: %v", err)
	}

	var out, errBuf bytes.Buffer
	if code := MergeGuard([]string{"--workflow-dir", root, "080-pr-merged", "--verdict", "passed"}, root, &out, &errBuf); code != 1 {
		t.Fatalf("failing commit must exit 1, got %d", code)
	}

	info, err := os.Stat(entity)
	if err != nil {
		t.Fatalf("stat restored entity: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("rollback must restore the original mode 0o600, got %o", got)
	}
}
