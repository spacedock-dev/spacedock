// ABOUTME: Behavioral replay of the REAL release.yml steps that protect the
// ABOUTME: stable channel and the `main` stamp, against bare-origin fixtures.
package release

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// The names of the release.yml steps that these replays execute. Each test reads
// the step's run block from the on-disk workflow. A rename or a deletion
// therefore reds here, and not at the next real cut.
const (
	stableRegressionGateStepName = "Gate the cut on the tag not regressing the stable channel"
	stableAdvanceStepName        = "Advance the stable channel ref to the tagged commit"
	mainStampStepName            = "Stamp main to the next edge prerelease version"
	decisionStepName             = "Decide whether this stable tag is the latest known release line"
)

// stampedFiles are the three files that the `main` stamp step rewrites.
var stampedFiles = []string{
	".claude-plugin/plugin.json",
	".codex-plugin/plugin.json",
	"skills/first-officer/references/first-officer-shared-core.md",
}

// releaseBinaryPath builds cmd/spacedock-release and gives the absolute path.
// The replayed steps call `go run ./cmd/spacedock-release`, which resolves
// against the process working directory. These replays run with the working
// directory INSIDE the fixture. At the repo root, a replayed stamp rewrites this
// checkout's own manifests. The swap of that one token for this binary is the
// single declared change to the script.
// Each other byte of the step — the gating, the refspecs, the conditions — runs
// as written.
func releaseBinaryPath(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "spacedock-release")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/spacedock-release")
	cmd.Dir = repoRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build spacedock-release: %v\n%s", err, out)
	}
	return bin
}

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// gitIn runs git in dir and gives the trimmed output. It fails the test when git
// exits non-zero, so a broken fixture never looks like a broken guard.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
	return strings.TrimSpace(string(out))
}

// writeFixtureVersion writes the two plugin manifests and the FO prose pin that
// the stamp step rewrites. version sets the manifest field, and the prose keeps
// the matching major.minor.
func writeFixtureVersion(t *testing.T, work, version string) {
	t.Helper()
	major, minor, ok := contract.ParseMajorMinor(version)
	if !ok {
		t.Fatalf("fixture version %q has no major.minor", version)
	}
	bodies := []string{
		fmt.Sprintf("{\"name\":\"spacedock\",\"version\":\"%s\"}\n", version),
		fmt.Sprintf("{\"name\":\"spacedock\",\"version\":\"%s\"}\n", version),
		fmt.Sprintf("These skills require binary minor %d.%d (same major.minor).\n", major, minor),
	}
	for i, rel := range stampedFiles {
		path := filepath.Join(work, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(bodies[i]), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// bareOriginFixture builds a throwaway bare origin and a work clone that carries
// it as `origin`, with `main` at version on both. The shape matches the release
// runner: a checkout with a real remote, so `ls-remote`, `fetch` and `push` all
// behave as they do in CI. The caller adds the tags and refs that each test
// needs.
func bareOriginFixture(t *testing.T, version string) (origin, work string) {
	t.Helper()
	base := t.TempDir()
	origin, work = filepath.Join(base, "origin.git"), filepath.Join(base, "work")
	for _, dir := range []string{origin, work} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	gitIn(t, origin, "init", "-q", "--bare", "-b", "main")
	testgit.InitRepo(t, work, "-q", "-b", "main")
	gitIn(t, work, "remote", "add", "origin", origin)
	writeFixtureVersion(t, work, version)
	gitIn(t, work, "add", "-A")
	gitIn(t, work, "commit", "-q", "-m", "release "+version)
	gitIn(t, work, "push", "-q", "origin", "main")
	return origin, work
}

// runReleaseStep executes the REAL run block of the named release.yml step, with
// dir as the working directory. It gives the exit code and the combined output.
func runReleaseStep(t *testing.T, stepName, dir string, env ...string) (int, string) {
	t.Helper()
	var run string
	for _, step := range parseWorkflowSteps(readWorkflow(t, "release.yml")) {
		if step.name == stepName {
			run = step.run
		}
	}
	if run == "" {
		t.Fatalf("release.yml has no step named %q", stepName)
	}
	script := strings.ReplaceAll(run, "go run ./cmd/spacedock-release", releaseBinaryPath(t))
	cmd := exec.Command("bash", "-c", script)
	cmd.Dir, cmd.Env = dir, append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	code := 0
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("run release.yml step %q: %v\n%s", stepName, err, out)
	}
	return code, string(out)
}

// assertStepInJob fails the test unless the named job owns the named step. Job
// identity is load-bearing for both moved steps. The gate must sit in e2e-gate,
// which goreleaser needs. The stamp must sit in edge-advance, which owns the
// latest-line decision.
func assertStepInJob(t *testing.T, jobName, stepName string) {
	t.Helper()
	for _, job := range parseWorkflowJobs(readWorkflow(t, "release.yml")) {
		if job.name != jobName {
			continue
		}
		for _, step := range job.steps {
			if step.name == stepName {
				return
			}
		}
	}
	t.Fatalf("release.yml job %q has no step named %q", jobName, stepName)
}

// oldLineFixture builds the version-inversion shape WITHOUT an ancestry
// inversion. `stable` points at the 0.28.0 commit. A CHILD commit carries
// 0.27.1 under the tag v0.27.1. The version goes DOWN while the ancestry goes
// FORWARD, so git accepts a push of the child to `stable`. A missing gate can
// therefore not hide behind a git refusal. This is the realistic mis-cut:
// docs/releasing.md tells the cutter to branch off `origin/main`, so a
// hand-stamped 0.27.1 lands as a child of today's `main`.
func oldLineFixture(t *testing.T) (origin, work string) {
	t.Helper()
	origin, work = bareOriginFixture(t, "0.28.0")
	gitIn(t, work, "push", "-q", "origin", "main:refs/heads/stable")
	writeFixtureVersion(t, work, "0.27.1")
	gitIn(t, work, "add", "-A")
	gitIn(t, work, "commit", "-q", "-m", "release 0.27.1")
	gitIn(t, work, "tag", "-a", "v0.27.1", "-m", "patch release")
	gitIn(t, work, "push", "-q", "origin", "main", "v0.27.1")
	return origin, work
}

// TestStableRegressionGateBlocksOlderLine is AC-1's block half. On the
// version-inversion fixture, the real gate step exits non-zero and names the tag
// and both versions. goreleaser therefore never starts. The job check is part of
// the proof. A gate in the goreleaser job starts the job that it must stop.
func TestStableRegressionGateBlocksOlderLine(t *testing.T) {
	assertStepInJob(t, "e2e-gate", stableRegressionGateStepName)
	_, work := oldLineFixture(t)
	code, out := runReleaseStep(t, stableRegressionGateStepName, work, "GITHUB_REF_NAME=v0.27.1")
	if code == 0 {
		t.Fatalf("stable-regression gate exited 0 on an old-line tag; the cut would publish\n%s", out)
	}
	for _, want := range []string{"v0.27.1", "0.28.0"} {
		if !strings.Contains(out, want) {
			t.Fatalf("gate output does not name %q:\n%s", want, out)
		}
	}
}

// TestUnguardedOldLineTagWouldReachStable is AC-1's independent baseline. On the
// SAME fixture, the real stable-advance step runs WITHOUT the gate and the push
// SUCCEEDS, and `stable` then carries 0.27.1. This is what today's pipeline does.
// It is also the proof that the gate is not redundant with git's own ancestry
// check. git raises no objection at all, so the block must come from the gate.
func TestUnguardedOldLineTagWouldReachStable(t *testing.T) {
	assertStepInJob(t, "goreleaser", stableAdvanceStepName)
	origin, work := oldLineFixture(t)
	code, out := runReleaseStep(t, stableAdvanceStepName, work, "GITHUB_REF_NAME=v0.27.1")
	if code != 0 {
		t.Fatalf("the stable-advance step exited %d; the fixture must let the regression through\n%s", code, out)
	}
	got := mustManifestVersion(t, []byte(gitIn(t, origin, "show", "stable:.claude-plugin/plugin.json")))
	if got != "0.27.1" {
		t.Fatalf("stable carries %q after the unguarded push, want 0.27.1", got)
	}
}

// TestStableRegressionGatePassesWhenStableRefAbsent covers the first stable
// release. `ls-remote --exit-code` gives exit 2 for a ref that does not exist.
// There is no release to regress, so the gate must let the run continue.
func TestStableRegressionGatePassesWhenStableRefAbsent(t *testing.T) {
	_, work := bareOriginFixture(t, "0.28.0")
	gitIn(t, work, "tag", "-a", "v0.28.0", "-m", "first stable release")
	code, out := runReleaseStep(t, stableRegressionGateStepName, work, "GITHUB_REF_NAME=v0.28.0")
	if code != 0 {
		t.Fatalf("gate exited %d with no stable ref; the first stable release must pass\n%s", code, out)
	}
}

// TestStableRegressionGateFailsClosedOnUnreadableRemote covers the other arm: an
// unreadable remote gives a non-zero code that is not 2, and the gate must stop
// the cut. An unreadable baseline is not a permission to publish, and the cost of
// a false block is one job re-run.
func TestStableRegressionGateFailsClosedOnUnreadableRemote(t *testing.T) {
	_, work := bareOriginFixture(t, "0.28.0")
	gitIn(t, work, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing.git"))
	code, out := runReleaseStep(t, stableRegressionGateStepName, work, "GITHUB_REF_NAME=v0.28.0")
	if code == 0 {
		t.Fatalf("gate exited 0 against an unreadable remote; it must fail closed\n%s", out)
	}
}

// runStampUnderDecision replays the edge-advance `decision` step and then, only
// when that step writes advance=true, the `main` stamp step. It gives the
// decision. GitHub Actions applies the same gate through the step `if:`. The
// wiring guard in edge_advance_wiring_test.go proves that the `if:` carries it.
// A replay drives a step directly and cannot prove that Actions gates it.
func runStampUnderDecision(t *testing.T, work, refName string) bool {
	t.Helper()
	outFile := filepath.Join(t.TempDir(), "github_output")
	if err := os.WriteFile(outFile, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	code, out := runReleaseStep(t, decisionStepName, work, "GITHUB_REF_NAME="+refName, "GITHUB_OUTPUT="+outFile)
	if code != 0 {
		t.Fatalf("decision step exited %d for %s\n%s", code, refName, out)
	}
	decision, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(decision), "advance=true") {
		return false
	}
	if code, out := runReleaseStep(t, mainStampStepName, work, "GITHUB_REF_NAME="+refName); code != 0 {
		t.Fatalf("main stamp step exited %d for %s\n%s", code, refName, out)
	}
	return true
}

// TestPatchTagDoesNotStampMain is AC-2: a patch tag that PASSES the regression
// gate still leaves `main` alone. The fixture carries v0.26.0, v0.27.0 and
// v0.28.0-pre0, so a v0.26.1 cut is two lines back and the decision skips. Both
// the tip SHA and the bytes of the three stamped files must not change. Today's
// pipeline instead stamps `main` DOWN to 0.26.1, which points each edge user's
// plugin at 0.26 while the binary is 0.28.0-pre0.
func TestPatchTagDoesNotStampMain(t *testing.T) {
	assertStepInJob(t, "edge-advance", mainStampStepName)
	origin, work := bareOriginFixture(t, "0.28.0-pre0")
	for _, tag := range []string{"v0.26.0", "v0.26.1", "v0.27.0", "v0.28.0-pre0"} {
		gitIn(t, work, "tag", tag)
	}
	before := map[string]string{"HEAD": gitIn(t, origin, "rev-parse", "main")}
	for _, rel := range stampedFiles {
		before[rel] = gitIn(t, origin, "show", "main:"+rel)
	}
	if runStampUnderDecision(t, work, "v0.26.1") {
		t.Fatal("the decision advanced an old-line patch tag; the stamp then rewrites main DOWN")
	}
	if after := gitIn(t, origin, "rev-parse", "main"); after != before["HEAD"] {
		t.Fatalf("main moved from %s to %s on an old-line patch tag", before["HEAD"], after)
	}
	for _, rel := range stampedFiles {
		if after := gitIn(t, origin, "show", "main:"+rel); after != before[rel] {
			t.Fatalf("main:%s changed on an old-line patch tag:\n%s", rel, after)
		}
	}
}

// TestLatestLineCutStampsMainToPre0 is AC-3: a latest-line stable cut needs zero
// human commits to restore the edge line. After the replay, `main`'s manifest
// carries the pre0 version that the same job's auto-cut tags, and the FO prose
// pin carries its major.minor. Today's pipeline instead leaves `main` at 0.27.0
// while the pre0 tag is 0.28.0-pre0. One hand-authored commit closed that
// mismatch 6m18s after the v0.27.0 cut.
func TestLatestLineCutStampsMainToPre0(t *testing.T) {
	origin, work := bareOriginFixture(t, "0.27.0")
	for _, tag := range []string{"v0.26.0", "v0.27.0"} {
		gitIn(t, work, "tag", tag)
	}
	if !runStampUnderDecision(t, work, "v0.27.0") {
		t.Fatal("the decision skipped a latest-line stable tag; main keeps the released version")
	}
	got := mustManifestVersion(t, []byte(gitIn(t, origin, "show", "main:.claude-plugin/plugin.json")))
	if got != "0.28.0-pre0" {
		t.Fatalf("main carries %q after a latest-line cut, want 0.28.0-pre0", got)
	}
	minor := mustProseMinorField(t, []byte(gitIn(t, origin, "show", "main:"+stampedFiles[2])))
	if minor != "0.28" {
		t.Fatalf("the FO prose pins minor %q, want 0.28", minor)
	}
}
