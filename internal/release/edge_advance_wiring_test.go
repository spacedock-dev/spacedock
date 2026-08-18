// ABOUTME: Structure guards for the edge-advance job's always-cut-pre0 auto-tag
// ABOUTME: (AC-4), with adversarial twins. The next-reconcile it used to share a
// ABOUTME: decision gate with is retired (the edge marketplace entry now tracks
// ABOUTME: main directly; see the entity's design-change note); the pre0 step's
// ABOUTME: replacement "latest line" gate is proven here via ifHasDecisionGate,
// ABOUTME: and the decision step's OWN candidate-pool/fail-closed behavior is
// ABOUTME: proven by real shell execution in edge_advance_decision_shell_test.go.
package release

import (
	"fmt"
	"strings"
	"testing"
)

// edgeAdvanceJob returns the parsed edge-advance job, or nil when the workflow
// has none — the single lookup every edge-advance guard shares.
func edgeAdvanceJob(workflow string) *workflowJob {
	for _, job := range parseWorkflowJobs(workflow) {
		if job.name == "edge-advance" {
			j := job
			return &j
		}
	}
	return nil
}

// ifSelectsStable reports whether an `if:` guard fires on a stable tag and not a
// prerelease — the negated `!contains(github.ref, '-')` form, optionally
// conjoined with a further gate.
func ifSelectsStable(ifCond string) bool {
	c := strings.ReplaceAll(ifCond, " ", "")
	return c == "!contains(github.ref,'-')" || strings.HasPrefix(c, "!contains(github.ref,'-')&&")
}

// ifHasDecisionGate reports whether an `if:` guard conjoins the edge-advance
// decision output (`steps.decision.outputs.advance == 'true'`) — the gate the
// auto-pre0 step must carry so an old-line/patch stable tag skips the cut.
func ifHasDecisionGate(ifCond string) bool {
	return strings.Contains(strings.ReplaceAll(ifCond, " ", ""), "steps.decision.outputs.advance=='true'")
}

// edgeAdvanceAutoPre0Step returns the always-cut-pre0 step (the one running
// edge-pre0-version to auto-create the vX.(Y+1).0-pre0 tag), or nil when none.
func edgeAdvanceAutoPre0Step(job workflowJob) *workflowStep {
	for i := range job.steps {
		for _, command := range executableShellCommands(job.steps[i].run) {
			if strings.Contains(command, "edge-pre0-version") {
				return &job.steps[i]
			}
		}
	}
	return nil
}

// assertAlwaysCutPre0 (AC-4) binds the always-cut-pre0 auto-tag's wiring: the
// step runs ONLY on the stable path (a `-pre` tag must not recurse), is gated on
// the decision step's "is this the latest known release line" output (an
// old-line patch tag would otherwise auto-cut a wrong, lower pre0 and bump the
// spacedock@next cask DOWN a minor), tags the GREENED release commit
// (`RELEASE_COMMIT` from `git rev-list -1 "$GITHUB_REF_NAME"` — so the pre0 run
// reuses the existing green e2e run), creates an ANNOTATED tag with a non-empty
// body (a lightweight/empty-body tag reds the release-notes extraction step
// before the binary builds), pushes it with a trigger-capable credential — the
// EDGE_RELEASE_DEPLOY_KEY SSH deploy key, NOT the workflow-suppressed default
// GITHUB_TOKEN (which lands the ref but fires no run) — and then verifies a
// release.yml run was created for the pre0 tag, failing loudly if none appears.
func assertAlwaysCutPre0(workflow string) error {
	job := edgeAdvanceJob(workflow)
	if job == nil {
		return fmt.Errorf("release.yml has no edge-advance job")
	}
	pre0 := edgeAdvanceAutoPre0Step(*job)
	if pre0 == nil {
		return fmt.Errorf("edge-advance has no auto-cut pre0 step (none runs edge-pre0-version)")
	}
	if !ifSelectsStable(pre0.ifCond) {
		return fmt.Errorf("auto-pre0 step must run only on the stable path (!contains(github.ref, '-')); if: %q would let a prerelease tag recurse", pre0.ifCond)
	}
	if !ifHasDecisionGate(pre0.ifCond) {
		return fmt.Errorf("auto-pre0 step is not decision-gated; if: %q — a non-advancing patch would attempt a colliding pre0 tag", pre0.ifCond)
	}

	var tagCmd, pushCmd string
	derivesGreenedSHA, assignsBody := false, false
	verifiesRun, failsOnMiss := false, false
	for _, command := range executableShellCommands(pre0.run) {
		switch {
		case strings.HasPrefix(command, "git tag "):
			tagCmd = command
		case strings.Contains(command, "git push"):
			pushCmd = command
		case command == `RELEASE_COMMIT="$(git rev-list -1 "$GITHUB_REF_NAME")"`:
			derivesGreenedSHA = true
		case strings.HasPrefix(command, `PRE0_BODY="`) && command != `PRE0_BODY=""`:
			assignsBody = true
		}
		if strings.Contains(command, "actions/workflows/release.yml/runs") && strings.Contains(command, "PRE0_TAG") {
			verifiesRun = true
		}
		if strings.Contains(command, "exit 1") {
			failsOnMiss = true
		}
	}
	if tagCmd == "" {
		return fmt.Errorf("auto-pre0 step runs no `git tag` command")
	}
	if !tagCmdIsAnnotatedWithBody(tagCmd) {
		return fmt.Errorf("auto-pre0 `git tag` is not annotated (-a) with a non-empty -m body: %q", tagCmd)
	}
	if !assignsBody {
		return fmt.Errorf("auto-pre0 step does not assign a non-empty PRE0_BODY — the annotated tag body must be non-empty or the release-notes extraction hard-errors")
	}
	if !derivesGreenedSHA {
		return fmt.Errorf("auto-pre0 step does not derive RELEASE_COMMIT from `git rev-list -1 \"$GITHUB_REF_NAME\"` (the e2e-gate-greened SHA)")
	}
	if tagCmdTarget(tagCmd) != `"$RELEASE_COMMIT"` {
		return fmt.Errorf("auto-pre0 `git tag` targets %q, want \"$RELEASE_COMMIT\" (the greened SHA); tagging next's tip places the pre0 on an ungreened commit e2e-gate would block", tagCmdTarget(tagCmd))
	}
	if pushCmd == "" {
		return fmt.Errorf("auto-pre0 step runs no `git push` command")
	}
	if strings.Contains(pushCmd, "GITHUB_TOKEN") {
		return fmt.Errorf("auto-pre0 push authenticates with the default GITHUB_TOKEN (workflow-suppressed — lands the ref but fires no release.yml run): %q", pushCmd)
	}
	if !strings.Contains(pushCmd, `git@github.com:${GITHUB_REPOSITORY}.git`) {
		return fmt.Errorf("auto-pre0 push does not use the trigger-capable deploy-key SSH transport git@github.com:${GITHUB_REPOSITORY}.git: %q", pushCmd)
	}
	if !verifiesRun {
		return fmt.Errorf("auto-pre0 step has no release.yml run-verification poll keyed on the pre0 tag — a suppressed credential would leave the edge binary silently behind")
	}
	if !failsOnMiss {
		return fmt.Errorf("auto-pre0 step's verify poll never exits non-zero when no run is found — a suppressed credential would pass silently")
	}
	return nil
}

// tagCmdIsAnnotatedWithBody reports whether a normalized `git tag …` command is
// annotated (`-a`) with a non-empty `-m` body (the first field after `-m` is not
// an empty quote). A lightweight tag (no `-a`) or an empty body both fail.
func tagCmdIsAnnotatedWithBody(tagCmd string) bool {
	fields := strings.Fields(tagCmd)
	hasA := false
	for _, f := range fields {
		if f == "-a" {
			hasA = true
		}
	}
	if !hasA {
		return false
	}
	for i, f := range fields {
		if f == "-m" && i+1 < len(fields) {
			body := fields[i+1]
			return body != `""` && body != `''` && body != ""
		}
	}
	return false
}

// tagCmdTarget returns the commit-ish positional argument of a `git tag -a
// <tagname> <commit> -m …` command (the SECOND positional, after the tag name),
// or "" when it carries fewer than two positionals before `-m`.
func tagCmdTarget(tagCmd string) string {
	fields := strings.Fields(tagCmd)
	var positional []string
	for i := 2; i < len(fields); i++ { // skip "git" "tag"
		f := fields[i]
		if f == "-m" {
			break
		}
		if strings.HasPrefix(f, "-") {
			continue
		}
		positional = append(positional, f)
	}
	if len(positional) < 2 {
		return ""
	}
	return positional[1]
}

// TestReleaseWorkflowAlwaysCutPre0 locks AC-4 against the on-disk release.yml.
// The adversarial twins: (a) retarget the tag from the greened RELEASE_COMMIT to
// next's tip; (b) drop -a/-m to make it a lightweight tag; (c) move the step to
// the prerelease path (recursion); (d) push via the workflow-suppressed
// GITHUB_TOKEN instead of the trigger-capable deploy key; (e) neuter the
// verify-or-fail run poll. Each must red.
func TestReleaseWorkflowAlwaysCutPre0(t *testing.T) {
	workflow := readWorkflow(t, "release.yml")
	if err := assertAlwaysCutPre0(workflow); err != nil {
		t.Fatalf("real release.yml always-cut-pre0 guard failed: %v", err)
	}

	const realTag = `git tag -a "$PRE0_TAG" "$RELEASE_COMMIT" -m "$PRE0_BODY"`

	// (a) tag next's tip instead of the greened release commit → ungreened SHA.
	nextTip := strings.Replace(workflow, realTag,
		`git tag -a "$PRE0_TAG" origin/next -m "$PRE0_BODY"`, 1)
	if nextTip == workflow {
		t.Fatal("fixture workflow missing the auto-pre0 `git tag` line to retarget")
	}
	if err := assertAlwaysCutPre0(nextTip); err == nil {
		t.Fatal("always-cut-pre0 guard accepted a pre0 tag on next's tip (ungreened SHA)")
	}

	// (b) lightweight tag (no -a, no -m) → reds the release-notes extraction.
	lightweight := strings.Replace(workflow, realTag,
		`git tag "$PRE0_TAG" "$RELEASE_COMMIT"`, 1)
	if lightweight == workflow {
		t.Fatal("fixture workflow missing the auto-pre0 `git tag` line to make lightweight")
	}
	if err := assertAlwaysCutPre0(lightweight); err == nil {
		t.Fatal("always-cut-pre0 guard accepted a lightweight (non-annotated) pre0 tag")
	}

	// (c) move the auto-pre0 step to the prerelease path → recursion.
	recursion := strings.Replace(workflow,
		"      - name: Auto-cut the edge prerelease tag on the greened release commit\n        if: \"!contains(github.ref, '-') && steps.decision.outputs.advance == 'true'\"\n",
		"      - name: Auto-cut the edge prerelease tag on the greened release commit\n        if: \"contains(github.ref, '-') && steps.decision.outputs.advance == 'true'\"\n",
		1)
	if recursion == workflow {
		t.Fatal("fixture workflow missing the auto-pre0 step guard to flip to the prerelease path")
	}
	if err := assertAlwaysCutPre0(recursion); err == nil {
		t.Fatal("always-cut-pre0 guard accepted an auto-pre0 step on the prerelease path (recursion)")
	}

	// (d) push via the workflow-suppressed GITHUB_TOKEN instead of the deploy key.
	wrongToken := strings.Replace(workflow,
		`git push "git@github.com:${GITHUB_REPOSITORY}.git" "$PRE0_TAG"`,
		`git push "https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git" "$PRE0_TAG"`, 1)
	if wrongToken == workflow {
		t.Fatal("fixture workflow missing the auto-pre0 deploy-key push line to swap")
	}
	if err := assertAlwaysCutPre0(wrongToken); err == nil {
		t.Fatal("always-cut-pre0 guard accepted a pre0 tag pushed via GITHUB_TOKEN (would not fire the pre0 run)")
	}

	// (e) neuter the verify-or-fail poll → a suppressed credential passes silently.
	noVerify := strings.Replace(workflow,
		"actions/workflows/release.yml/runs",
		"actions/workflows/DISABLED.yml/runs", 1)
	if noVerify == workflow {
		t.Fatal("fixture workflow missing the auto-pre0 verify poll to neuter")
	}
	if err := assertAlwaysCutPre0(noVerify); err == nil {
		t.Fatal("always-cut-pre0 guard accepted an auto-pre0 step whose verify poll no longer queries release.yml runs")
	}
}
