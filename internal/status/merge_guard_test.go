// ABOUTME: merge guard verb behavior — arm / blocked / finalize over the proven
// ABOUTME: --set/--archive paths, propagating (never bypassing) the guard refusals.
package status

import (
	"bytes"
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
