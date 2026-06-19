// ABOUTME: merge guard verb behavior — arm / blocked / finalize over the proven
// ABOUTME: --set/--archive paths, propagating (never bypassing) the guard refusals.
package status

import (
	"bytes"
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

// TestMergeGuardArmIsPolicyGated is the load-bearing inversion pin (AC-1 vs AC-5).
// On the IDENTICAL precondition — empty mod-block, empty pr, `--verdict passed` —
// the ONLY differing variable is the `merge:` policy, and it must flip the
// outcome: merge: local ARMS (exit 0, sets mod-block, does not terminalize), while
// merge: pr REFUSES at the merge-hook guard (exit 1, no arm, entity unchanged).
// This resolves the policy-blindness contradiction between the Phase A prose and
// AC-5: an empty-mod-block / empty-pr entity is indistinguishable as "armable" vs
// "ceremony-skipping" by frontmatter, so a policy-blind arm cannot satisfy AC-5.
// Gating arm to merge: local (fully in-process, safe to auto-drive) while merge: pr
// routes the unarmed finalize to the guard refusal uniquely satisfies the oracle.
// Flipping arm back to policy-blind turns the merge: pr leg RED (it would arm,
// exit 0, instead of refusing).
func TestMergeGuardArmIsPolicyGated(t *testing.T) {
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

	// merge: pr leg — refuses, does NOT arm. Same precondition, opposite outcome.
	prRoot, prOut, prErr, prCode := driveMergeGuard(t, "merge-pr-workflow", "020-no-sentinel", "--verdict", "passed")
	if prCode != 1 {
		t.Fatalf("merge: pr empty-mod-block --verdict passed must REFUSE (exit 1, not arm), got %d (stdout=%q)", prCode, prOut)
	}
	if strings.Contains(prOut, "armed") {
		t.Fatalf("merge: pr leg must NOT arm, got %q", prOut)
	}
	if !strings.Contains(prErr, "cannot advance to terminal") {
		t.Fatalf("merge: pr leg must propagate the merge-hook refusal, got %q", prErr)
	}
	if got := frontmatterField(t, filepath.Join(prRoot, "020-no-sentinel.md"), "mod-block"); got != "" {
		t.Fatalf("merge: pr refusal must NOT arm a mod-block, got %q", got)
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

// TestMergeGuardPrNoSentinelRefuses (AC-5): under default merge: pr with a
// registered hook, a --verdict passed finalize on an entity with no pr and no
// mod-block propagates the merge-hook guard's `cannot advance to terminal` refusal
// (exit 1) verbatim, leaves the entity unchanged, and never passes --force.
func TestMergeGuardPrNoSentinelRefuses(t *testing.T) {
	root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", "020-no-sentinel", "--verdict", "passed")
	if code != 1 {
		t.Fatalf("merge: pr no-sentinel finalize must refuse (exit 1), got %d", code)
	}
	if !strings.Contains(errOut, "cannot advance to terminal") {
		t.Fatalf("stderr should propagate the merge-hook refusal verbatim, got %q", errOut)
	}
	if strings.Contains(out, "finalized") {
		t.Fatalf("stdout must not claim finalized on a refusal, got %q", out)
	}
	entity := filepath.Join(root, "020-no-sentinel.md")
	if got := frontmatterField(t, entity, "status"); got != "implementation" {
		t.Fatalf("refused entity must be unchanged, status=%q", got)
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
