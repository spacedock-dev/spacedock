// ABOUTME: Behavioral proof that the live-CI test steps emit a clean step log plus
// ABOUTME: an archived -json detail from one run, with the exit code preserved.
package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// liveCleanOutputStepNames are the runtime-live-e2e.yml steps transformed to the
// clean-stdout + archived-jsonl shape. The behavioral test derives its runnable
// script from the FIRST of these (re-targeted to the in-repo fixture), so it
// exercises the workflow's ACTUAL pipe shape rather than a hand-copied duplicate.
var liveCleanOutputStepNames = []string{
	"Run live ensign cycle",
	"Run live Claude shared scenarios",
	"Run live Codex shared scenarios",
	"Run Pi shared scenario coverage guard",
	"Run live Pi front-door smoke",
}

// fixtureScript extracts a transformed live step's run block from the workflow and
// re-targets it at the committed fixture: it strips the `-tags live` and
// `-test.timeout` (the fixture is neither live nor slow) and rewrites the
// ensigncycle package path, binary names, and -test.run filter to the fixture, so
// the script the test runs is the workflow's own shape with only the target
// swapped. extraFlags is appended to the test-binary invocation (e.g. an env-gated
// run filter); fail toggles the planted failures.
func fixtureScript(t *testing.T, fixtureDir, jsonl, raw, runFilter string) string {
	t.Helper()
	live := readWorkflow(t, "runtime-live-e2e.yml")
	steps := parseWorkflowSteps(live)
	var block string
	for _, s := range steps {
		if s.name == liveCleanOutputStepNames[0] {
			block = s.run
			break
		}
	}
	if block == "" {
		t.Fatalf("workflow has no step %q to derive the fixture script from", liveCleanOutputStepNames[0])
	}
	block = dedent(block)

	// Re-target the ensigncycle live step at the fixture: drop the live tag and
	// the long timeout, point the compile + run + archive at the fixture dir and
	// fixed local artifact names, and replace the live -test.run alternation with
	// the caller's filter.
	repl := strings.NewReplacer(
		"go test -c -tags live -o live-e2e.test ./internal/ensigncycle/",
		`go -C "`+fixtureDir+`" test -c -o "$PWD/fixture.test" .`,
		"./live-e2e.test", "./fixture.test",
		"-test.timeout=40m ", "",
		"-test.run 'TestLiveEnsignCycle|TestLiveDefaultHeadlessStopsAtGate|TestLiveZeroDiscoverReportsAndStops'", runFilter,
		"-p ensigncycle", "-p cleanoutputfixture",
		"live-e2e-raw.txt", raw,
		"live-e2e-detail.jsonl", jsonl,
		`mkdir -p "$SPACEDOCK_LIVE_ARTIFACT_DIR"`, ":",
	)
	script := repl.Replace(block)

	// Guard the re-target: every ensigncycle/live token must be gone, or the
	// derived script silently diverged from what we think the workflow says.
	for _, leftover := range []string{"ensigncycle", "-tags live", "40m", "TestLive"} {
		if strings.Contains(script, leftover) {
			t.Fatalf("fixture re-target left %q in the derived script:\n%s", leftover, script)
		}
	}
	if !strings.Contains(script, "go tool test2json") || !strings.Contains(script, "${PIPESTATUS[0]}") {
		t.Fatalf("derived script lost the one-run archive shape:\n%s", script)
	}
	return script
}

// runFixtureStep runs the derived script under bash (the ${PIPESTATUS[0]} exit
// capture is bash-specific, matching the GitHub Actions default shell) in a temp
// dir, returning stdout+stderr, the exit code, and the temp dir.
func runFixtureStep(t *testing.T, script string, env ...string) (string, int, string) {
	t.Helper()
	dir := t.TempDir()
	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		} else {
			t.Fatalf("running fixture step: %v\n%s", err, out)
		}
	}
	return string(out), code, dir
}

func fixtureDir(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join("testdata", "cleanoutputfixture"))
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

var runLineRe = regexp.MustCompile(`(?m)^\s*=== RUN`)

// TestLiveCIStepStdoutIsCleanOnGreen is AC-1 (green): the changed command's
// visible surface has no per-test === RUN verbosity, and the suite passes.
func TestLiveCIStepStdoutIsCleanOnGreen(t *testing.T) {
	script := fixtureScript(t, fixtureDir(t), "detail.jsonl", "raw.txt", "-test.run '.*'")
	out, code, _ := runFixtureStep(t, script)

	if code != 0 {
		t.Fatalf("green fixture run exited %d, want 0\n%s", code, out)
	}
	if n := len(runLineRe.FindAllString(out, -1)); n != 0 {
		t.Errorf("green step stdout carries %d '=== RUN' lines (want 0 — the firehose must not reach stdout):\n%s", n, out)
	}
	if strings.Contains(out, "--- PASS:") {
		t.Errorf("green step stdout carries per-test '--- PASS:' verbosity (want clean surface):\n%s", out)
	}
}

// TestLiveCIStepStdoutIsFailuresOnlyOnRed is AC-1 (red): on a failing run the
// visible surface shows the failing tests with file:line and the package result,
// still with no === RUN firehose.
func TestLiveCIStepStdoutIsFailuresOnlyOnRed(t *testing.T) {
	script := fixtureScript(t, fixtureDir(t), "detail.jsonl", "raw.txt", "-test.run '.*'")
	out, code, _ := runFixtureStep(t, script, "FIXTURE_FAIL=1")

	if code == 0 {
		t.Fatalf("red fixture run exited 0, want non-zero (rendering must not mask failures)\n%s", out)
	}
	if n := len(runLineRe.FindAllString(out, -1)); n != 0 {
		t.Errorf("red step stdout carries %d '=== RUN' lines (want 0):\n%s", n, out)
	}
	for _, want := range []string{
		"--- FAIL: TestGammaFails",
		"a_test.go:", // the planted failure's file:line
		"--- FAIL: TestZetaSubtests/case_bad",
		"FAIL",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("red step stdout missing failures-only signal %q:\n%s", want, out)
		}
	}
}

// TestLiveCIStepArchivesJSONDetail is AC-2: the same run writes a non-empty -json
// archive that carries the failing test's event, sufficient for root cause.
func TestLiveCIStepArchivesJSONDetail(t *testing.T) {
	script := fixtureScript(t, fixtureDir(t), "detail.jsonl", "raw.txt", "-test.run '.*'")
	_, _, dir := runFixtureStep(t, script, "FIXTURE_FAIL=1")

	data, err := os.ReadFile(filepath.Join(dir, "detail.jsonl"))
	if err != nil {
		t.Fatalf("archive detail.jsonl not written: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("archive detail.jsonl is empty")
	}
	body := string(data)
	if !strings.Contains(body, `"Action":"fail"`) {
		t.Errorf("archive carries no fail event (root cause unreachable):\n%s", body)
	}
	if !strings.Contains(body, "TestGammaFails") || !strings.Contains(body, "compute() = 7, want 42") {
		t.Errorf("archive does not carry the planted failure's test + output for root cause:\n%s", body)
	}
}

// TestLiveCIStepRunsSuiteOnce is AC-3: the changed command runs the test binary
// exactly once. A counter test appends one byte per execution; a "render clean,
// then re-run for json" regression would record two.
func TestLiveCIStepRunsSuiteOnce(t *testing.T) {
	dir := t.TempDir()
	counter := filepath.Join(dir, "counter")
	script := fixtureScript(t, fixtureDir(t), "detail.jsonl", "raw.txt", "-test.run TestCounter")

	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "FIXTURE_COUNTER_FILE="+counter)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("counter fixture run failed: %v\n%s", err, out)
	}

	data, err := os.ReadFile(counter)
	if err != nil {
		t.Fatalf("counter file not written — the test binary did not run: %v", err)
	}
	if got := len(data); got != 1 {
		t.Errorf("test binary executed %d times, want exactly 1 (a double run would archive by re-running)", got)
	}
}

// TestLiveCIStepPreservesExitCode is AC-4: exit 0 on a passing fixture, non-zero
// on a failing one — the trailing clean-view grep must not mask the test result.
func TestLiveCIStepPreservesExitCode(t *testing.T) {
	script := fixtureScript(t, fixtureDir(t), "detail.jsonl", "raw.txt", "-test.run '.*'")

	if _, code, _ := runFixtureStep(t, script); code != 0 {
		t.Errorf("passing fixture: step exit %d, want 0", code)
	}
	if _, code, _ := runFixtureStep(t, script, "FIXTURE_FAIL=1"); code == 0 {
		t.Error("failing fixture: step exit 0, want non-zero (grep no-match exit must not mask the failure)")
	}
}
