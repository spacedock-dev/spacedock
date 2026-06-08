// ABOUTME: AC-3 survey install-detect probe-behavior test — extracts the exact step-1
// ABOUTME: one-liner from SKILL.md and runs it under present/absent PATH conditions.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// probeLineRe matches the step-1 install-detect one-liner in skills/survey/SKILL.md: a
// complete `if ! agentsview --version … then echo "AGENTSVIEW MISSING"; fi` statement on
// one line. It pins the execve form (`agentsview --version`) and the sentinel; an
// FS-access regression (`command -v` / `which` / `test -x` / `stat`) or a leaked banner
// would not match this shape, so extraction itself is part of the guard.
var probeLineRe = regexp.MustCompile(
	`^if ! agentsview --version >/dev/null 2>&1; then echo "AGENTSVIEW MISSING"; fi$`)

// extractProbeLine reads skills/survey/SKILL.md and returns the single runnable step-1
// probe statement (the artifact under test). The test EXECUTES the shipped line rather
// than a copy, so the SKILL.md probe and this test cannot drift; an FS-access regression
// fails extraction here (no matching line), and a leaked-output regression fails the
// behavioral assertions below.
func extractProbeLine(t *testing.T) string {
	t.Helper()
	path := filepath.Join(repoRoot(t), "skills", "survey", "SKILL.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read SKILL.md %s: %v", path, err)
	}
	var found []string
	for _, line := range strings.Split(string(data), "\n") {
		if probeLineRe.MatchString(strings.TrimSpace(line)) {
			found = append(found, strings.TrimSpace(line))
		}
	}
	if len(found) != 1 {
		t.Fatalf("expected exactly one runnable install-detect probe line in SKILL.md matching %q, found %d: %v",
			probeLineRe.String(), len(found), found)
	}
	return found[0]
}

// runProbe runs the extracted probe one-liner under bash with PATH set to exactly
// pathDir (nothing else), returning combined stdout+stderr and the process exit code.
// A single-entry PATH is the deterministic control: the probe reaches `agentsview` only
// if pathDir holds it, so the present/absent conditions are fully synthesized — no
// dependence on whatever `agentsview` the host machine happens to have installed.
func runProbe(t *testing.T, probe, pathDir string) (string, int) {
	t.Helper()
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not on PATH; the survey probe is a bash one-liner")
	}
	cmd := exec.Command(bash, "-c", probe)
	cmd.Env = []string{"PATH=" + pathDir}
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			t.Fatalf("run probe: %v\n%s", err, out)
		}
	}
	return strings.TrimRight(string(out), "\n"), code
}

// TestSurveyInstallProbe is AC-3: the deterministic install-detect probe-behavior test.
// It runs the exact step-1 one-liner from SKILL.md twice — once with a stub `agentsview`
// present on a synthesized PATH (expect empty output, exit 0), once with the name absent
// (expect the sole line `AGENTSVIEW MISSING`). The oracle is the two independent fixture
// CONDITIONS, never a SKILL.md grep: the present condition asserts the silent-success
// contract, the absent condition asserts the sentinel. A revert to an FS-access probe
// fails extraction (no matching line); a leaked `--version` banner fails the present case.
func TestSurveyInstallProbe(t *testing.T) {
	probe := extractProbeLine(t)

	// Present: a stub `agentsview` that exits 0 silently (mimicking `--version >/dev/null`)
	// is the ONLY thing on PATH. The probe must produce no output and exit 0.
	t.Run("present", func(t *testing.T) {
		dir := t.TempDir()
		stub := filepath.Join(dir, "agentsview")
		if err := os.WriteFile(stub, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("write stub agentsview: %v", err)
		}
		out, code := runProbe(t, probe, dir)
		if out != "" {
			t.Errorf("with agentsview present the probe must be silent (no banner leak), got %q", out)
		}
		if code != 0 {
			t.Errorf("with agentsview present the probe must exit 0, got %d", code)
		}
	})

	// Absent: PATH is an empty dir, so `agentsview` is not found. The probe must emit the
	// sole sentinel line `AGENTSVIEW MISSING` (the underlying not-found exit is suppressed
	// by 2>&1, so only the sentinel reaches the agent).
	t.Run("absent", func(t *testing.T) {
		dir := t.TempDir() // empty — no agentsview
		out, code := runProbe(t, probe, dir)
		if out != "AGENTSVIEW MISSING" {
			t.Errorf("with agentsview absent the probe must emit the sole line %q, got %q", "AGENTSVIEW MISSING", out)
		}
		if code != 0 {
			t.Errorf("the probe's `if`-guard absorbs the not-found exit; the statement itself exits 0, got %d", code)
		}
	})
}
