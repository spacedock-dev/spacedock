// ABOUTME: AC-1/AC-2 proof — the edge-line reconcile closes `next`'s drift to the
// ABOUTME: release commit in a temp git repo, force-free, on both tag paths.
package release

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// edgeReconcileResult carries the SHAs and merge output a reconcile-fixture run
// produces, so the subtests can assert the divergence math, the tree match, and
// the force-free (first-parent-ancestor) push shape against real git state.
type edgeReconcileResult struct {
	dir             string
	releaseVersion  string
	releaseCommit   string
	preMergeNextTip string
	mergeCommit     string
	finalNextTip    string
	mergeOutput     string
	behindBefore    int // git rev-list --count next..<release-commit> before the merge
	aheadBefore     int // git rev-list --count <release-commit>..next before the merge
}

// runEdgeReconcileFixture builds a temp repo shaped like the real incident —
// `next` diverged behind `main` with some `next`-exclusive commits — then
// reconstructs, in local git plumbing, the reconcile sequence release.yml's
// edge-advance job runs for `tag`: a `git merge` whose strategy option is read
// from the on-disk workflow step (so a drift to `-X ours` reds this fixture
// directly), then, on a stable tag, a dev-preversion stamp, then a calendar bump.
// The stamp/bump/dev-preversion steps call the SAME production functions
// release.yml's `spacedock-release` subcommands wrap; the surrounding git
// plumbing mirrors the workflow's shell rather than shelling out to it.
func runEdgeReconcileFixture(t *testing.T, tag string) edgeReconcileResult {
	t.Helper()
	dir := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return string(out)
	}
	write := func(rel, content string) {
		t.Helper()
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	countCommits := func(rangeSpec string) int {
		t.Helper()
		n, err := strconv.Atoi(strings.TrimSpace(git("rev-list", "--count", rangeSpec)))
		if err != nil {
			t.Fatalf("parse rev-list --count %s: %v", rangeSpec, err)
		}
		return n
	}

	pluginJSON := func(version string) string {
		return "{\n  \"name\": \"spacedock\",\n  \"version\": \"" + version + "\"\n}\n"
	}
	codexJSON := func(version string) string {
		return "{\n  \"name\": \"spacedock-codex\",\n  \"version\": \"" + version + "\"\n}\n"
	}
	const preCutCalendar = "0.0.2026050101"
	marketplaceJSON := "{\n  \"name\": \"spacedock-edge\",\n  \"plugins\": [\n    {\n      \"name\": \"spacedock\",\n      \"version\": \"" + preCutCalendar + "\"\n    }\n  ]\n}\n"

	testgit.InitRepo(t, dir, "-q")

	// Base commit — the shared ancestor `next` and `main` both descend from.
	write(".claude-plugin/plugin.json", pluginJSON("0.23.0-pre1"))
	write(".codex-plugin/plugin.json", codexJSON("0.23.0-pre1"))
	write(".claude-plugin/marketplace.json", marketplaceJSON)
	write("app.txt", "base\n")
	git("add", "-A")
	git("commit", "-q", "-m", "seed")
	git("branch", "-M", "main")
	git("branch", "next")

	// Advance `main` past the shared base with real content changes, then the
	// release-prep stamp — the tagged commit reads the release version, mirroring
	// the stamp-then-tag ordering.
	releaseVersion := strings.TrimPrefix(tag, "v")
	write("app.txt", "feature-1\n")
	git("commit", "-q", "-am", "feat: one")
	write("app.txt", "feature-2\n")
	git("commit", "-q", "-am", "feat: two")
	write(".claude-plugin/plugin.json", pluginJSON(releaseVersion))
	write(".codex-plugin/plugin.json", codexJSON(releaseVersion))
	git("commit", "-q", "-am", "release: stamp "+releaseVersion)
	releaseCommit := strings.TrimSpace(git("rev-parse", "HEAD"))
	git("tag", tag, releaseCommit)

	// Advance `next` with exclusive commits touching the SAME version lines main
	// finalizes — a genuine divergence the reconcile must absorb, resolved in
	// main's favor by -X theirs (matching the real 0.24.0-pre1 incident).
	git("switch", "-q", "next")
	write(".claude-plugin/plugin.json", pluginJSON("0.23.0-pre2"))
	git("commit", "-q", "-am", "next: wip plugin")
	write(".codex-plugin/plugin.json", codexJSON("0.23.0-pre2"))
	git("commit", "-q", "-am", "next: wip codex")
	preMergeNextTip := strings.TrimSpace(git("rev-parse", "HEAD"))

	res := edgeReconcileResult{
		dir:             dir,
		releaseVersion:  releaseVersion,
		releaseCommit:   releaseCommit,
		preMergeNextTip: preMergeNextTip,
		behindBefore:    countCommits("next.." + releaseCommit),
		aheadBefore:     countCommits(releaseCommit + "..next"),
	}

	// The reconcile: merge the release commit into `next`, favoring the release.
	// The `-X <strategy>` option is read from the on-disk release.yml edge-advance
	// reconcile step for this tag's path, so this fixture merges with EXACTLY the
	// strategy CI uses — a workflow drift to `-X ours` (the 0.24.0-pre1 incident)
	// reds the tree-match assertion below rather than passing against a stale
	// hardcoded `theirs`. Run from `next`, so the merge commit's FIRST parent is
	// the pre-merge next tip — the property that makes the push a fast-forward.
	workflow := readWorkflow(t, "release.yml")
	reconcileStep := edgeAdvanceReconcileStep(t, workflow, strings.Contains(tag, "-"))
	strategy := mergeStrategyOption(mergeCommand(reconcileStep))
	if strategy == "" {
		t.Fatalf("edge-advance reconcile step %q has no `-X <strategy>` merge option to mirror", reconcileStep.name)
	}
	mergeCmd := exec.Command("git", "merge", "-X", strategy, "--no-edit", releaseCommit,
		"-m", "next: reconcile edge line to "+tag)
	mergeCmd.Dir = dir
	mergeOut, err := mergeCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("reconcile merge failed (want clean auto-merge): %v\n%s", err, mergeOut)
	}
	res.mergeOutput = string(mergeOut)
	res.mergeCommit = strings.TrimSpace(git("rev-parse", "HEAD"))

	// Fixed date so the bumped calendar key is deterministic and strictly greater
	// than the pre-cut value (0.0.20260501NN).
	now := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)

	// Stable path only: stamp `next` PAST the release to the dev pre-version.
	if !strings.Contains(tag, "-") {
		dev, derr := DevPreVersion(releaseVersion)
		if derr != nil {
			t.Fatalf("DevPreVersion(%q): %v", releaseVersion, derr)
		}
		for _, rel := range []string{".claude-plugin/plugin.json", ".codex-plugin/plugin.json"} {
			path := filepath.Join(dir, rel)
			data, rerr := os.ReadFile(path)
			if rerr != nil {
				t.Fatal(rerr)
			}
			out, serr := StampVersion(data, dev)
			if serr != nil {
				t.Fatalf("StampVersion %s: %v", rel, serr)
			}
			if werr := os.WriteFile(path, out, 0o644); werr != nil {
				t.Fatal(werr)
			}
		}
		git("commit", "-q", "-m", "next: bump dev pre-version to "+dev,
			"--", ".claude-plugin/plugin.json", ".codex-plugin/plugin.json")
	}

	// Both paths: bump the marketplace calendar key so `plugin update` re-pulls.
	mktPath := filepath.Join(dir, ".claude-plugin/marketplace.json")
	mktData, rerr := os.ReadFile(mktPath)
	if rerr != nil {
		t.Fatal(rerr)
	}
	bumped, berr := BumpCalendarVersion(mktData, now)
	if berr != nil {
		t.Fatalf("BumpCalendarVersion: %v", berr)
	}
	if werr := os.WriteFile(mktPath, bumped, 0o644); werr != nil {
		t.Fatal(werr)
	}
	git("commit", "-q", "-m", "next: bump marketplace calendar version",
		"--", ".claude-plugin/marketplace.json")

	res.finalNextTip = strings.TrimSpace(git("rev-parse", "HEAD"))
	return res
}

// readEdgeVersion returns the top-level `version` of a committed plugin manifest
// in the fixture repo's working tree (which reflects the final `next` tip).
func readEdgeVersion(t *testing.T, dir, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, rel))
	if err != nil {
		t.Fatal(err)
	}
	v, err := ManifestVersion(data)
	if err != nil {
		t.Fatalf("parse %s: %v", rel, err)
	}
	return v
}

// readEdgeCalendar returns the marketplace entry's calendar version in the
// fixture repo's working tree.
func readEdgeCalendar(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, ".claude-plugin/marketplace.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Plugins []struct {
			Version string `json:"version"`
		} `json:"plugins"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse marketplace.json: %v", err)
	}
	if len(doc.Plugins) == 0 {
		t.Fatalf("marketplace.json has no plugin entry")
	}
	return doc.Plugins[0].Version
}

// gitCount runs `git rev-list --count <args...>` in dir, so callers can pass a
// bare range or a range plus flags (e.g. --first-parent to count only the
// commits authored on top of `next`, excluding the release side the merge pulls
// in through its second parent).
func gitCount(t *testing.T, dir string, args ...string) int {
	t.Helper()
	cmd := exec.Command("git", append([]string{"rev-list", "--count"}, args...)...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-list --count %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		t.Fatalf("parse rev-list count %q: %v", out, err)
	}
	return n
}

// isFirstParentAncestor reports whether ancestor is reachable from descendant
// (`git merge-base --is-ancestor`), the property that makes pushing descendant
// to the branch ancestor tips a fast-forward — never a force.
func isFirstParentAncestor(t *testing.T, dir, ancestor, descendant string) bool {
	t.Helper()
	cmd := exec.Command("git", "merge-base", "--is-ancestor", ancestor, descendant)
	cmd.Dir = dir
	err := cmd.Run()
	if err == nil {
		return true
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return false
	}
	t.Fatalf("git merge-base --is-ancestor: %v", err)
	return false
}

// TestEdgeLineReconcileClosesDivergence locks AC-1/AC-2a/AC-3/AC-4 on both tag
// paths: starting from a `next` genuinely diverged from the release commit (ahead
// AND behind, not a fast-forward case), the reconcile closes the behind-divergence
// to zero with a clean, conflict-free merge whose tree matches the release commit
// exactly, and the pre-merge `next` tip stays a first-parent ancestor of the
// pushed tip (fast-forward, never a force). The prerelease path leaves the
// manifests at the tag version; the stable path stamps them PAST it to the dev
// pre-version. Both bump the calendar key past its pre-cut value.
func TestEdgeLineReconcileClosesDivergence(t *testing.T) {
	cases := []struct {
		name          string
		tag           string
		wantEdge      string // expected plugin manifest version on the final `next` tip
		newAfterMerge int    // commits added on top of the pre-merge next tip (merge + stamps)
	}{
		{name: "prerelease", tag: "v0.24.0-pre1", wantEdge: "0.24.0-pre1", newAfterMerge: 2},
		{name: "stable", tag: "v0.24.0", wantEdge: "0.25.0-pre1", newAfterMerge: 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := runEdgeReconcileFixture(t, c.tag)

			// Precondition: genuinely diverged BOTH ways, mirroring the incident's
			// 40-commit drift with `next`-exclusive commits on top.
			if res.behindBefore == 0 {
				t.Fatalf("fixture is not behind the release before the merge (behind=%d); nothing to reconcile", res.behindBefore)
			}
			if res.aheadBefore == 0 {
				t.Fatalf("fixture is not ahead of the release before the merge (ahead=%d); not a genuine divergence, just a fast-forward", res.aheadBefore)
			}

			// AC-2a: the merge is conflict-free (auto-resolved by -X theirs).
			if strings.Contains(res.mergeOutput, "CONFLICT") {
				t.Fatalf("reconcile merge reported a conflict:\n%s", res.mergeOutput)
			}

			// AC-1 / precedent: the merge commit's tree is byte-for-byte the release
			// commit's tree (no `next`-only content survives the -X theirs favor).
			cmd := exec.Command("git", "diff", "--stat", res.releaseCommit, res.mergeCommit)
			cmd.Dir = res.dir
			diffOut, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("git diff: %v\n%s", err, diffOut)
			}
			if strings.TrimSpace(string(diffOut)) != "" {
				t.Fatalf("merged tree does not match the release commit's tree:\n%s", diffOut)
			}

			// AC-1 VALUE metric: the behind-divergence closes to zero — the release
			// commit is fully contained in `next`'s history after the merge.
			if behindAfter := gitCount(t, res.dir, "HEAD.."+res.releaseCommit); behindAfter != 0 {
				t.Fatalf("behind-divergence did not close: git rev-list --count HEAD..<release> = %d, want 0", behindAfter)
			}

			// Only the expected reconcile commits (merge + optional dev stamp +
			// calendar bump) are authored on top of the pre-merge `next` tip. The
			// first-parent walk excludes the release-side commits the merge pulls in
			// through its second parent.
			if got := gitCount(t, res.dir, "--first-parent", res.preMergeNextTip+"..HEAD"); got != c.newAfterMerge {
				t.Fatalf("reconcile commits added over the pre-merge next tip = %d, want %d", got, c.newAfterMerge)
			}

			// AC-2a: the pre-merge `next` tip is an ancestor of the pushed tip, so
			// `git push edge-advance:next` fast-forwards — never a force.
			if !isFirstParentAncestor(t, res.dir, res.preMergeNextTip, res.finalNextTip) {
				t.Fatalf("pre-merge next tip %s is not an ancestor of the pushed tip %s; the push would not fast-forward", res.preMergeNextTip, res.finalNextTip)
			}

			// AC-3 (prerelease) / AC-4 (stable): the manifests read the expected
			// edge version, and never the just-released stable version.
			for _, rel := range []string{".claude-plugin/plugin.json", ".codex-plugin/plugin.json"} {
				if got := readEdgeVersion(t, res.dir, rel); got != c.wantEdge {
					t.Errorf("%s version on next = %q, want %q", rel, got, c.wantEdge)
				}
			}
			if c.name == "stable" {
				if got := readEdgeVersion(t, res.dir, ".claude-plugin/plugin.json"); got == res.releaseVersion {
					t.Errorf("stable path left the edge manifest at the released stable version %q; it must be stamped past it", res.releaseVersion)
				}
			}

			// AC-3: the calendar key advances past its pre-cut value on both paths.
			if got := readEdgeCalendar(t, res.dir); !(got > "0.0.2026050101") {
				t.Errorf("marketplace calendar key = %q, did not advance past the pre-cut 0.0.2026050101", got)
			}
		})
	}
}

// TestReleaseWorkflowEdgeAdvanceNeverForces locks AC-2b: the real edge-advance
// job exists as a `needs: goreleaser` sibling (so a reconcile conflict cannot
// unwind or block the already-published release) and NONE of its steps
// force-pushes or resets `next` — the reconcile is a fast-forward merge, never a
// `--force`/`-f`/`reset --hard`. The adversarial twins inject each dangerous form
// into the job and assert the guard REDs, so a green result is not vacuous.
func TestReleaseWorkflowEdgeAdvanceNeverForces(t *testing.T) {
	workflow := readWorkflow(t, "release.yml")
	if err := assertEdgeAdvanceIsForceFreeSibling(workflow); err != nil {
		t.Fatalf("real release.yml edge-advance job fails the never-forces guard: %v", err)
	}

	forcePush := strings.Replace(workflow,
		"git push origin edge-advance:next",
		"git push --force origin edge-advance:next",
		1)
	if forcePush == workflow {
		t.Fatal("fixture workflow missing the `git push origin edge-advance:next` line to mutate")
	}
	if err := assertEdgeAdvanceIsForceFreeSibling(forcePush); err == nil {
		t.Fatal("never-forces guard accepted an edge-advance job with a --force push")
	}

	resetHard := strings.Replace(workflow,
		"git switch -c edge-advance origin/next",
		"git reset --hard origin/next",
		1)
	if resetHard == workflow {
		t.Fatal("fixture workflow missing the `git switch -c edge-advance origin/next` line to mutate")
	}
	if err := assertEdgeAdvanceIsForceFreeSibling(resetHard); err == nil {
		t.Fatal("never-forces guard accepted an edge-advance job with a reset --hard")
	}
}

// assertEdgeAdvanceIsForceFreeSibling binds AC-2b to the parsed job graph: the
// edge-advance job must exist, `needs:` the goreleaser job, and carry no step
// that force-pushes or resets `next`.
func assertEdgeAdvanceIsForceFreeSibling(workflow string) error {
	edge := edgeAdvanceJob(workflow)
	if edge == nil {
		return fmt.Errorf("release.yml has no edge-advance job")
	}
	needsGoreleaser := false
	for _, need := range edge.needs {
		if need == "goreleaser" {
			needsGoreleaser = true
		}
	}
	if !needsGoreleaser {
		return fmt.Errorf("edge-advance job does not declare needs: goreleaser; a reconcile failure could block or unwind the published release")
	}
	for _, step := range edge.steps {
		for _, command := range executableShellCommands(step.run) {
			if commandForcePushesOrResets(command) {
				return fmt.Errorf("edge-advance step force-pushes or resets next: %q", command)
			}
		}
	}
	return nil
}

// edgeAdvanceJob returns the parsed edge-advance job, or nil when the workflow
// has none — the single lookup every edge-advance guard and the fixture share.
func edgeAdvanceJob(workflow string) *workflowJob {
	for _, job := range parseWorkflowJobs(workflow) {
		if job.name == "edge-advance" {
			j := job
			return &j
		}
	}
	return nil
}

// edgeAdvanceReconcileSteps returns the edge-advance job's reconcile steps — the
// ones whose run block carries a `git merge`. Callers distinguish the two by
// their `if:` guard, not by order.
func edgeAdvanceReconcileSteps(job workflowJob) []workflowStep {
	var steps []workflowStep
	for _, step := range job.steps {
		if mergeCommand(step) != "" {
			steps = append(steps, step)
		}
	}
	return steps
}

// edgeAdvanceReconcileStep returns the reconcile step guarding the given tag path
// — prerelease selects the `contains(github.ref, '-')` step, stable the
// `!contains(...)` step — from the parsed workflow, so the fixture can mirror the
// exact merge strategy CI runs for that path.
func edgeAdvanceReconcileStep(t *testing.T, workflow string, prerelease bool) workflowStep {
	t.Helper()
	job := edgeAdvanceJob(workflow)
	if job == nil {
		t.Fatal("release.yml has no edge-advance job")
	}
	for _, step := range edgeAdvanceReconcileSteps(*job) {
		if prerelease && ifSelectsPrerelease(step.ifCond) {
			return step
		}
		if !prerelease && ifSelectsStable(step.ifCond) {
			return step
		}
	}
	t.Fatalf("edge-advance has no reconcile step guarding the prerelease=%v path", prerelease)
	return workflowStep{}
}

// edgeAdvanceBumpCalendarStep returns the edge-advance step that bumps the
// marketplace calendar key (and pushes the edge line), or nil when none does.
func edgeAdvanceBumpCalendarStep(job workflowJob) *workflowStep {
	for i := range job.steps {
		for _, command := range executableShellCommands(job.steps[i].run) {
			if strings.Contains(command, "bump-calendar .claude-plugin/marketplace.json") {
				return &job.steps[i]
			}
		}
	}
	return nil
}

// mergeCommand returns the first normalized `git merge …` command in a step's
// run block, or "" when the step runs none.
func mergeCommand(step workflowStep) string {
	for _, command := range executableShellCommands(step.run) {
		if strings.HasPrefix(command, "git merge ") {
			return command
		}
	}
	return ""
}

// mergeStrategyOption returns the `-X <strategy>` argument of a `git merge`
// command (e.g. "theirs"), or "" when the command carries none.
func mergeStrategyOption(command string) string {
	fields := strings.Fields(command)
	for i, f := range fields {
		switch {
		case f == "-X" && i+1 < len(fields):
			return fields[i+1]
		case strings.HasPrefix(f, "-X") && len(f) > len("-X"):
			return f[len("-X"):]
		case strings.HasPrefix(f, "--strategy-option="):
			return strings.TrimPrefix(f, "--strategy-option=")
		}
	}
	return ""
}

// ifSelectsPrerelease reports whether an `if:` guard fires on a prerelease
// (hyphenated) tag and not a stable one — the un-negated `contains(github.ref,
// '-')` form, optionally conjoined with the edge-advance decision gate
// (`&& steps.decision.outputs.advance == 'true'`), which every reconcile step
// now carries so an old-line tag skips the whole job.
func ifSelectsPrerelease(ifCond string) bool {
	c := strings.ReplaceAll(ifCond, " ", "")
	return c == "contains(github.ref,'-')" || strings.HasPrefix(c, "contains(github.ref,'-')&&")
}

// ifSelectsStable reports whether an `if:` guard fires on a stable tag and not a
// prerelease — the negated `!contains(github.ref, '-')` form, optionally
// conjoined with the decision gate. It is the exact logical complement of
// ifSelectsPrerelease on the tag-shape axis (the `!` prefix distinguishes them),
// so a step pair carrying both fires on exactly one tag shape each.
func ifSelectsStable(ifCond string) bool {
	c := strings.ReplaceAll(ifCond, " ", "")
	return c == "!contains(github.ref,'-')" || strings.HasPrefix(c, "!contains(github.ref,'-')&&")
}

// TestReleaseWorkflowEdgeAdvanceWiring locks the edge-advance job's step-level
// wiring against the on-disk release.yml, closing the coupling gaps the detached
// audit found: both reconcile steps must merge with `-X theirs` (favor the
// release — the `-X ours` flip is the exact 0.24.0-pre1 divergence incident) and
// their `if:` guards must be complementary on the tag-shape axis so EXACTLY one
// fires per tag (prerelease `contains(github.ref, '-')`, stable `!contains(...)`,
// each conjoined with the decision gate). The adversarial twins flip the strategy
// to `-X ours`, widen the stable guard to `always()`, and copy the prerelease
// guard onto the stable step; each must red, so a green result is not vacuous.
// The decision-gate conjunct itself is guarded by TestReleaseWorkflowEdgeAdvanceDecisionGates.
func TestReleaseWorkflowEdgeAdvanceWiring(t *testing.T) {
	workflow := readWorkflow(t, "release.yml")
	if err := assertEdgeAdvanceWiring(workflow); err != nil {
		t.Fatalf("real release.yml edge-advance wiring guard failed: %v", err)
	}

	// -X ours on BOTH reconcile steps — the 0.24.0-pre1 incident reintroduction.
	ours := strings.ReplaceAll(workflow, "git merge -X theirs", "git merge -X ours")
	if ours == workflow {
		t.Fatal("fixture workflow missing `git merge -X theirs` to mutate")
	}
	if err := assertEdgeAdvanceWiring(ours); err == nil {
		t.Fatal("wiring guard accepted edge-advance reconcile steps using -X ours")
	}

	// always() on the stable step — both branches fire on a prerelease, so the
	// stable dev-preversion stamp wrongly runs on a `-pre` cut. Anchored to the
	// stable step's name so it does not mutate the other `!contains(...)` guards.
	always := strings.Replace(workflow,
		"      - name: Reconcile the edge line past the stable release\n        if: \"!contains(github.ref, '-') && steps.decision.outputs.advance == 'true'\"\n",
		"      - name: Reconcile the edge line past the stable release\n        if: \"always()\"\n",
		1)
	if always == workflow {
		t.Fatal("fixture workflow missing the stable reconcile step guard to widen to always()")
	}
	if err := assertEdgeAdvanceWiring(always); err == nil {
		t.Fatal("wiring guard accepted a stable reconcile step guarded by always()")
	}

	// Copy-paste: the stable step reuses the prerelease `contains(...)` guard, so
	// the stable dev-preversion path would never fire on any tag.
	copyPaste := strings.Replace(workflow,
		"      - name: Reconcile the edge line past the stable release\n        if: \"!contains(github.ref, '-') && steps.decision.outputs.advance == 'true'\"\n",
		"      - name: Reconcile the edge line past the stable release\n        if: \"contains(github.ref, '-') && steps.decision.outputs.advance == 'true'\"\n",
		1)
	if copyPaste == workflow {
		t.Fatal("fixture workflow missing the stable reconcile step guard to copy-paste")
	}
	if err := assertEdgeAdvanceWiring(copyPaste); err == nil {
		t.Fatal("wiring guard accepted two reconcile steps with the same non-complementary if: guard")
	}
}

// assertEdgeAdvanceWiring consolidates the edge-advance job's step-level wiring
// invariants against the parsed workflow: it is a force-free `needs: goreleaser`
// sibling (assertEdgeAdvanceIsForceFreeSibling), both reconcile steps merge with
// `-X theirs`, and their `if:` guards are complementary on the tag-shape axis
// (exactly one fires per tag). The decision-gate conjunct each step now also
// carries — and the calendar-bump/auto-pre0 gating — is asserted separately by
// assertEdgeAdvanceDecisionGating.
func assertEdgeAdvanceWiring(workflow string) error {
	if err := assertEdgeAdvanceIsForceFreeSibling(workflow); err != nil {
		return err
	}
	job := edgeAdvanceJob(workflow)
	if job == nil {
		return fmt.Errorf("release.yml has no edge-advance job")
	}

	reconcile := edgeAdvanceReconcileSteps(*job)
	if len(reconcile) != 2 {
		return fmt.Errorf("edge-advance has %d reconcile (git merge) steps, want 2", len(reconcile))
	}
	prereleaseGuarded, stableGuarded := false, false
	for _, step := range reconcile {
		if strat := mergeStrategyOption(mergeCommand(step)); strat != "theirs" {
			return fmt.Errorf("edge-advance reconcile step %q merges with -X %q, want theirs (favor the release)", step.name, strat)
		}
		switch {
		case ifSelectsPrerelease(step.ifCond):
			prereleaseGuarded = true
		case ifSelectsStable(step.ifCond):
			stableGuarded = true
		default:
			return fmt.Errorf("edge-advance reconcile step %q has if: %q, want contains/!contains(github.ref, '-')", step.name, step.ifCond)
		}
	}
	if !prereleaseGuarded || !stableGuarded {
		return fmt.Errorf("edge-advance reconcile steps are not complementary (prerelease-guarded=%v stable-guarded=%v); exactly one must fire per tag", prereleaseGuarded, stableGuarded)
	}
	return nil
}

// commandForcePushesOrResets reports whether a shell command force-pushes
// (`--force`/`--force-with-lease`, or `git push … -f`) or hard-resets
// (`reset --hard`) — the destructive forms the edge-advance reconcile must never
// use on `next`.
func commandForcePushesOrResets(command string) bool {
	if strings.Contains(command, "--force") {
		return true
	}
	fields := strings.Fields(command)
	for i, f := range fields {
		if f == "reset" && i+1 < len(fields) && fields[i+1] == "--hard" {
			return true
		}
	}
	if strings.Contains(command, "git push") {
		for _, f := range fields {
			if f == "-f" {
				return true
			}
		}
	}
	return false
}
