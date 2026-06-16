// ABOUTME: Behavioral proof that the live-CI test steps emit a clean step log plus
// ABOUTME: an archived -json detail from one gotestsum run, with the exit preserved.
package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// gotestsumBin resolves the pinned gotestsum the live steps run. It is on PATH in
// CI (the live jobs run .github/scripts/install-gotestsum.sh first; the offline
// gate installs it before `go test ./...` so this proof actually executes there).
// When absent — a local checkout without gotestsum — the AC tests skip rather than
// build-from-source, per the FO's AC-5 note; the binary-level guarantee is the
// real assurance and CI always exercises it.
func gotestsumBin(t *testing.T) string {
	t.Helper()
	if p, err := exec.LookPath("gotestsum"); err == nil {
		return p
	}
	if dir := os.Getenv("GOTESTSUM_BIN_DIR"); dir != "" {
		p := filepath.Join(dir, "gotestsum")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	t.Skip("gotestsum not on PATH or GOTESTSUM_BIN_DIR; CI installs the pinned binary before this runs")
	return ""
}

func fixtureDir(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join("testdata", "cleanoutputfixture"))
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// runGotestsum runs the live steps' command shape — `gotestsum --jsonfile
// <jsonl> --format pkgname -- <go test args>` — over the in-repo fixture, in a
// temp dir, returning combined stdout+stderr, the exit code, and the temp dir.
// The args mirror a live step minus `-tags live` (the fixture has no live tag);
// the rendering flags (`--jsonfile`, `--format pkgname`) are the workflow's.
func runGotestsum(t *testing.T, runFilter string, env ...string) (string, int, string) {
	t.Helper()
	bin := gotestsumBin(t)
	dir := t.TempDir()
	jsonl := filepath.Join(dir, "detail.jsonl")
	args := []string{
		"--jsonfile", jsonl,
		"--format", "pkgname",
		"--", "-count=1", "-run", runFilter, ".",
	}
	cmd := exec.Command(bin, args...)
	// Run inside the fixture module so `go test .` resolves its go.mod; the
	// jsonl archive is written to the temp dir via the absolute --jsonfile path.
	cmd.Dir = fixtureDir(t)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		} else {
			t.Fatalf("running gotestsum: %v\n%s", err, out)
		}
	}
	return string(out), code, dir
}

var runLineRe = regexp.MustCompile(`(?m)^\s*=== RUN`)

// TestLiveCIStepStdoutIsCleanOnGreen is AC-1 (green): the changed command's
// visible surface has no per-test === RUN verbosity, and the suite passes.
func TestLiveCIStepStdoutIsCleanOnGreen(t *testing.T) {
	out, code, _ := runGotestsum(t, ".*")

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
// visible surface shows the failing tests with file:line, still with no === RUN
// firehose, and the run goes red.
func TestLiveCIStepStdoutIsFailuresOnlyOnRed(t *testing.T) {
	out, code, _ := runGotestsum(t, ".*", "FIXTURE_FAIL=1")

	if code == 0 {
		t.Fatalf("red fixture run exited 0, want non-zero (rendering must not mask failures)\n%s", out)
	}
	if n := len(runLineRe.FindAllString(out, -1)); n != 0 {
		t.Errorf("red step stdout carries %d '=== RUN' lines (want 0):\n%s", n, out)
	}
	for _, want := range []string{
		"TestGammaFails",
		"a_test.go:", // the planted failure's file:line
		"TestZetaSubtests/case_bad",
		"Failed", // gotestsum's "=== Failed" recap header
	} {
		if !strings.Contains(out, want) {
			t.Errorf("red step stdout missing failures-only signal %q:\n%s", want, out)
		}
	}
}

// TestLiveCIStepArchivesJSONDetail is AC-2: the same run writes a non-empty -json
// archive that carries the failing test's event, sufficient for root cause.
func TestLiveCIStepArchivesJSONDetail(t *testing.T) {
	_, _, dir := runGotestsum(t, ".*", "FIXTURE_FAIL=1")

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
// exactly once. A counter test appends one byte per execution; a render-then-rerun
// regression would record two.
func TestLiveCIStepRunsSuiteOnce(t *testing.T) {
	bin := gotestsumBin(t)
	dir := t.TempDir()
	counter := filepath.Join(dir, "counter")
	cmd := exec.Command(bin,
		"--jsonfile", filepath.Join(dir, "detail.jsonl"), "--format", "pkgname",
		"--", "-count=1", "-run", "TestCounter", ".",
	)
	cmd.Dir = fixtureDir(t)
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
// on a failing one — gotestsum's rendering must not mask the test result.
func TestLiveCIStepPreservesExitCode(t *testing.T) {
	if _, code, _ := runGotestsum(t, ".*"); code != 0 {
		t.Errorf("passing fixture: step exit %d, want 0", code)
	}
	if _, code, _ := runGotestsum(t, ".*", "FIXTURE_FAIL=1"); code == 0 {
		t.Error("failing fixture: step exit 0, want non-zero")
	}
}
