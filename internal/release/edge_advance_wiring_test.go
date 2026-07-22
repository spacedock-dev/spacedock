// ABOUTME: Structure guards for the edge-advance line-ordering decision gate (AC-2)
// ABOUTME: and the always-cut-pre0 auto-tag (AC-4), each with adversarial twins.
package release

import (
	"fmt"
	"strings"
	"testing"
)

// ifHasDecisionGate reports whether an `if:` guard conjoins the edge-advance
// decision output (`steps.decision.outputs.advance == 'true'`), the gate every
// mutating step in the job must carry so an old-line/patch tag skips the whole
// job.
func ifHasDecisionGate(ifCond string) bool {
	return strings.Contains(strings.ReplaceAll(ifCond, " ", ""), "steps.decision.outputs.advance=='true'")
}

// edgeAdvanceDecisionStep returns the step that runs the edge-advance-decision
// subcommand (the job's first, gating step), or nil when none does.
func edgeAdvanceDecisionStep(job workflowJob) *workflowStep {
	for i := range job.steps {
		for _, command := range executableShellCommands(job.steps[i].run) {
			if strings.Contains(command, "edge-advance-decision") {
				return &job.steps[i]
			}
		}
	}
	return nil
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

// assertEdgeAdvanceDecisionGating (AC-2) binds the line-ordering guard's job-level
// wiring: a `decision` step runs edge-advance-decision and writes BOTH
// advance=true and advance=false to $GITHUB_OUTPUT, and EVERY mutating step —
// both reconcile steps, the calendar-bump+push step, and the auto-pre0 step —
// gates on `steps.decision.outputs.advance == 'true'`. The calendar-bump gate is
// the one ASK 2 added: previously unconditional, an old-line patch would
// otherwise still churn every edge installer's re-pull (AC-2d / AC-3d).
func assertEdgeAdvanceDecisionGating(workflow string) error {
	job := edgeAdvanceJob(workflow)
	if job == nil {
		return fmt.Errorf("release.yml has no edge-advance job")
	}
	decision := edgeAdvanceDecisionStep(*job)
	if decision == nil {
		return fmt.Errorf("edge-advance has no step running edge-advance-decision")
	}
	if decision.id != "decision" {
		return fmt.Errorf("edge-advance decision step has id %q, want \"decision\" (else steps.decision.outputs.advance does not resolve)", decision.id)
	}
	emitsTrue, emitsFalse := false, false
	for _, command := range executableShellCommands(decision.run) {
		if strings.Contains(command, "advance=true") && strings.Contains(command, `"$GITHUB_OUTPUT"`) {
			emitsTrue = true
		}
		if strings.Contains(command, "advance=false") && strings.Contains(command, `"$GITHUB_OUTPUT"`) {
			emitsFalse = true
		}
	}
	if !emitsTrue || !emitsFalse {
		return fmt.Errorf("edge-advance decision step does not write both advance=true and advance=false to $GITHUB_OUTPUT (true=%v false=%v)", emitsTrue, emitsFalse)
	}
	for _, step := range edgeAdvanceReconcileSteps(*job) {
		if !ifHasDecisionGate(step.ifCond) {
			return fmt.Errorf("edge-advance reconcile step %q does not gate on the decision output; if: %q", step.name, step.ifCond)
		}
	}
	bump := edgeAdvanceBumpCalendarStep(*job)
	if bump == nil {
		return fmt.Errorf("edge-advance has no bump-calendar step")
	}
	if !ifHasDecisionGate(bump.ifCond) {
		return fmt.Errorf("edge-advance bump-calendar step is not decision-gated; if: %q — an old-line patch would still churn every edge installer's re-pull", bump.ifCond)
	}
	pre0 := edgeAdvanceAutoPre0Step(*job)
	if pre0 == nil {
		return fmt.Errorf("edge-advance has no auto-cut pre0 step")
	}
	if !ifHasDecisionGate(pre0.ifCond) {
		return fmt.Errorf("edge-advance auto-pre0 step is not decision-gated; if: %q", pre0.ifCond)
	}
	return nil
}

// TestReleaseWorkflowEdgeAdvanceDecisionGates locks AC-2 against the on-disk
// release.yml: the whole job gates on the line-ordering decision. The adversarial
// twins ungate the calendar-bump step, ungate the stable reconcile step, and
// break the decision step's advance=false emission; each must red, so the green
// is not vacuous.
func TestReleaseWorkflowEdgeAdvanceDecisionGates(t *testing.T) {
	workflow := readWorkflow(t, "release.yml")
	if err := assertEdgeAdvanceDecisionGating(workflow); err != nil {
		t.Fatalf("real release.yml edge-advance decision gating failed: %v", err)
	}

	// Ungate the calendar bump — an old-line patch would churn every edge
	// installer's re-pull (the AC-2d / AC-3d regression ASK 2 closed).
	ungatedBump := strings.Replace(workflow,
		"      - name: Bump the marketplace calendar key and push the edge line\n        if: \"steps.decision.outputs.advance == 'true'\"\n",
		"      - name: Bump the marketplace calendar key and push the edge line\n",
		1)
	if ungatedBump == workflow {
		t.Fatal("fixture workflow missing the decision-gated bump-calendar step to ungate")
	}
	if err := assertEdgeAdvanceDecisionGating(ungatedBump); err == nil {
		t.Fatal("decision-gating guard accepted an ungated calendar-bump step")
	}

	// Ungate the stable reconcile step — an old-line patch would `-X theirs`
	// clobber next's newer content.
	ungatedReconcile := strings.Replace(workflow,
		"      - name: Reconcile the edge line past the stable release\n        if: \"!contains(github.ref, '-') && steps.decision.outputs.advance == 'true'\"\n",
		"      - name: Reconcile the edge line past the stable release\n        if: \"!contains(github.ref, '-')\"\n",
		1)
	if ungatedReconcile == workflow {
		t.Fatal("fixture workflow missing the decision-gated stable reconcile step to ungate")
	}
	if err := assertEdgeAdvanceDecisionGating(ungatedReconcile); err == nil {
		t.Fatal("decision-gating guard accepted an ungated stable reconcile step")
	}

	// Break the decision step's advance=false emission — downstream gates would
	// never see a skip signal.
	brokenDecision := strings.Replace(workflow,
		"            echo \"advance=false\" >> \"$GITHUB_OUTPUT\"\n",
		"            echo \"skipped the edge line\"\n",
		1)
	if brokenDecision == workflow {
		t.Fatal("fixture workflow missing the decision step advance=false emission to break")
	}
	if err := assertEdgeAdvanceDecisionGating(brokenDecision); err == nil {
		t.Fatal("decision-gating guard accepted a decision step that never emits advance=false")
	}
}

// assertAlwaysCutPre0 (AC-4) binds the always-cut-pre0 auto-tag's wiring: the
// step runs ONLY on the stable path (a `-pre` tag must not recurse), is
// decision-gated, tags the GREENED release commit (`RELEASE_COMMIT` from
// `git rev-list -1 "$GITHUB_REF_NAME"`, not next's tip — so the pre0 run reuses
// the existing green e2e run), creates an ANNOTATED tag with a non-empty body (a
// lightweight/empty-body tag reds the release-notes extraction step before the
// binary builds), pushes it with a trigger-capable credential — the
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
