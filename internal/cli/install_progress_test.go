package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestInstallVerbosePreservesOutputAndNamesSkippedWork(t *testing.T) {
	t.Setenv("PATH", "/usr/bin:/bin")
	stubClaudeInventory(t, "[]")
	var ordinaryOut, ordinaryErr, stdout, stderr bytes.Buffer
	args := []string{"install", "--check", "--host", "claude"}
	ordinaryCode := Run(args, &ordinaryOut, &ordinaryErr)
	code := Run(append(args, "--verbose"), &stdout, &stderr)

	if code != ordinaryCode || stdout.String() != ordinaryOut.String() {
		t.Fatalf("verbose install changed its result: (%d, %q) vs (%d, %q)", code, stdout.String(), ordinaryCode, ordinaryOut.String())
	}
	for _, line := range []string{
		"Skipped: Install Claude Code plugins (check only)",
		"Start: Resolve Claude Code plugin manifest",
		"Done: Resolve Claude Code plugin manifest",
		"Total:",
		"(exit 0)",
	} {
		if !strings.Contains(stderr.String(), line) {
			t.Errorf("stderr did not contain %q:\n%s", line, &stderr)
		}
	}
}

func TestInstallVerboseUsageDoesNotStartTiming(t *testing.T) {
	for _, args := range [][]string{
		{"install", "--verbose", "--host"},
		{"install", "--host", "--verbose"},
		{"install", "--verbose", "--plugin-dir"},
		{"install", "--plugin-dir", "--verbose"},
		{"install", "--verbose", "--unknown"},
		{"doctor", "--verbose"},
	} {
		var stdout, stderr bytes.Buffer
		if code := Run(args, &stdout, &stderr); code != 2 || strings.Contains(stderr.String(), "[+") {
			t.Errorf("%v: wanted untimed exit 2, got (%d, %q)", args, code, &stderr)
		}
	}
}

func sameInstallSteps(got, want []installStep) bool {
	return slices.EqualFunc(got, want, func(a, b installStep) bool {
		a.operation, b.operation = "", ""
		return reflect.DeepEqual(a, b)
	})
}

func TestInstallVerboseUsesCommandClock(t *testing.T) {
	t.Setenv("PATH", "/usr/bin:/bin")
	stubClaudeInventory(t, "[]")
	now := time.Unix(0, 0)
	clock := func() time.Time {
		now = now.Add(125 * time.Millisecond)
		return now
	}
	var stdout, stderr bytes.Buffer
	args := []string{"--verbose", "--host", "claude", "--check"}
	if code := runInstallCommand(context.Background(), args, os.Environ(), t.TempDir(), &stdout, &stderr, clock); code != 0 {
		t.Fatalf("install exited %d: %s", code, &stderr)
	}
	want := "[+0.125s] Skipped: Install Claude Code plugins (check only)\n" +
		"[+0.250s] Start: Resolve Claude Code plugin manifest\n" +
		"[+0.375s] Done: Resolve Claude Code plugin manifest (0.125s)\n"
	if !strings.HasPrefix(stderr.String(), want) {
		t.Errorf("trace did not use the command clock for offsets and durations:\n%s", &stderr)
	}
}

func TestFinalVersionProbeTracesChildOutcome(t *testing.T) {
	for _, inventory := range []struct {
		name   string
		report func(versionInputs) versionReport
	}{
		{"installed versions", collectVersions},
		{"Pi versions", func(in versionInputs) versionReport {
			report, _ := piVersionReport(in)
			return report
		}},
	} {
		for _, tc := range []struct {
			name, child, version, trace string
		}{
			{"exit", "echo 0.17.0; echo private >&2; exit 37", "unknown", "Warning: Check subspace program version (0.125s; child exit 37; continuing)"},
			{"signal", `kill -TERM "$$"`, "unknown", "Interrupted: Check subspace program version (0.125s; child ended by signal; continuing)"},
			{"empty", "exit 0", "unknown", "Warning: Check subspace program version (0.125s; version unavailable; continuing)"},
			{"version", "echo 0.17.0; echo private >&2", "0.17.0", "Done: Check subspace program version (0.125s)"},
		} {
			t.Run(inventory.name+"/"+tc.name, func(t *testing.T) {
				dir := t.TempDir()
				program := filepath.Join(dir, "spacedock-review")
				if err := os.WriteFile(program, []byte("#!/bin/sh\n[ \"$*\" = --version ] || exit 99\n"+tc.child+"\n"), 0o755); err != nil {
					t.Fatal(err)
				}
				var stderr bytes.Buffer
				now := time.Unix(0, 0)
				in := versionInputs{
					env:      []string{"HOME=" + dir, "PATH=" + dir},
					lookPath: func(string) (string, error) { return "", os.ErrNotExist },
					installProgress: &installProgress{newElapsedTrace(&stderr, func() time.Time {
						now = now.Add(125 * time.Millisecond)
						return now
					})},
				}
				verbose := inventory.report(in)
				in.installProgress = nil
				ordinary := inventory.report(in)
				if !reflect.DeepEqual(verbose, ordinary) || ordinary.version("subspace program") != tc.version {
					t.Errorf("inventory changed or returned the wrong version: verbose %q, ordinary %q, want %q", verbose.version("subspace program"), ordinary.version("subspace program"), tc.version)
				}
				if trace := stderr.String(); !strings.Contains(trace, tc.trace) || strings.Contains(trace, "private") {
					t.Errorf("wanted %q and no child output, got:\n%s", tc.trace, trace)
				}
			})
		}
	}
}

// fixedStepClock advances 125ms per read, so every span lasts 0.125s.
func fixedStepClock() func() time.Time {
	now := time.Unix(0, 0)
	return func() time.Time {
		now = now.Add(125 * time.Millisecond)
		return now
	}
}

// traceLines drops the [+offset] prefix from each trace line, keeping the
// status, label and duration a user reads.
func traceLines(trace string) []string {
	var lines []string
	for _, line := range strings.Split(strings.TrimSpace(trace), "\n") {
		if _, rest, ok := strings.Cut(line, "] "); ok {
			lines = append(lines, rest)
		}
	}
	return lines
}

func TestPiInstallVerboseRuntimeSpanMatchesInstallReadiness(t *testing.T) {
	for _, tc := range []struct {
		name, args string
		ready      bool
		trace      string
		stdout     string
	}{
		{"install without auth", "--host pi", true, "Done: Check Pi runtime (", "Pi runtime ready."},
		{"install with missing packages", "--host pi", false, "Warning: Check Pi runtime (", "Pi runtime setup incomplete."},
		{"check without auth", "--host pi --check", true, "Failed: Check Pi runtime (", "MISSING Pi auth"},
	} {
		packageSpans := []string{"Done: Remove Pi Spacedock package (0.125s)", "Done: Install Pi Spacedock package (0.125s)"}
		if strings.Contains(tc.args, "--check") {
			packageSpans = []string{"Skipped: Install Pi packages (check only)"}
		}
		t.Run(tc.name, func(t *testing.T) {
			repo := t.TempDir()
			writePiSkillFixtures(t, repo)
			pkg := t.TempDir()
			writePiSubagentsFixtures(t, pkg)
			statOK := statOKForPiResources(repo, pkg)
			if !tc.ready {
				delete(statOK, pkg+"-intercom")
			}
			piOps := &fakePiRuntimeOps{lookPath: piHealthyPathFixtures(), statOK: statOK, packageStatus: healthyPiPackageStatus()}
			// No auth.json and no OPENAI_API_KEY: launch-ready is still possible, doctor-healthy is not.
			env := append(piTestEnv(pkg, t.TempDir()), "SPACEDOCK_REPO_ROOT="+repo)
			args := strings.Fields(tc.args)
			var stdout, stderr bytes.Buffer
			setup, code := parsePiSetupArgs("install", args, &stderr)
			if code != 0 {
				t.Fatalf("parse exited %d: %s", code, &stderr)
			}
			progress := &installProgress{newElapsedTrace(&stderr, fixedStepClock())}
			runInstallSetup(context.Background(), args, setup, &fakeHost{}, piOps, env, &stdout, &stderr, progress)
			if !strings.Contains(stdout.String(), tc.stdout) {
				t.Fatalf("stdout missing %q:\n%s", tc.stdout, &stdout)
			}
			assertTraceInOrder(t, stderr.String(), packageSpans)
			var span string
			for _, line := range traceLines(stderr.String()) {
				if strings.Contains(line, ": Check Pi runtime (") {
					span = line
				}
			}
			if !strings.HasPrefix(span, tc.trace) {
				t.Errorf("runtime span = %q, want prefix %q (stdout says %q)\n%s", span, tc.trace, tc.stdout, &stderr)
			}
		})
	}
}

// assertTraceInOrder fails unless want appears in lines as an ordered subsequence.
func assertTraceInOrder(t *testing.T, trace string, want []string) {
	t.Helper()
	lines, next := traceLines(trace), 0
	for _, line := range lines {
		if next < len(want) && line == want[next] {
			next++
		}
	}
	if next < len(want) {
		t.Errorf("trace is missing %q (or it is out of order):\n%s", want[next], trace)
	}
}

func TestInstallVerboseTracesEachClaudeInstallStep(t *testing.T) {
	t.Setenv("PATH", "/usr/bin:/bin")
	t.Setenv("HOME", t.TempDir())
	// subspace's edge copy is installed at user scope, so its removal runs and the
	// stub fails it; every other removal is skipped.
	stubClaudeInventory(t, `[{"id":"subspace@spacedock-edge","scope":"user"}]`)
	dir := t.TempDir()
	args := []string{"--host", "claude"}
	var ordinaryOut, ordinaryErr, stdout, stderr bytes.Buffer
	ordinaryCode := runInstallCommand(context.Background(), args, os.Environ(), dir, &ordinaryOut, &ordinaryErr, fixedStepClock())
	code := runInstallCommand(context.Background(), append(args, "--verbose"), os.Environ(), dir, &stdout, &stderr, fixedStepClock())

	if code != 0 || code != ordinaryCode || stdout.String() != ordinaryOut.String() {
		t.Fatalf("verbose install changed its result: (%d, %q) vs (%d, %q)", code, stdout.String(), ordinaryCode, ordinaryOut.String())
	}
	if ordinaryErr.Len() != 0 {
		t.Errorf("an ordinary install printed a trace:\n%s", &ordinaryErr)
	}
	assertTraceInOrder(t, stderr.String(), []string{
		"Start: List Claude Code user plugins before install",
		"Done: List Claude Code user plugins before install (0.125s)",
		"Skipped: Claude Code other-channel Spacedock plugin removal (not installed at user scope)",
		"Skipped: Claude Code retired Spacedock plugin removal (not installed at user scope)",
		"Start: Claude Code marketplace add",
		"Done: Claude Code marketplace add (0.125s)",
		"Start: Claude Code marketplace refresh",
		"Done: Claude Code marketplace refresh (0.125s)",
		"Skipped: Claude Code current-channel Spacedock plugin removal (not installed at user scope)",
		"Start: Claude Code Spacedock plugin install",
		"Done: Claude Code Spacedock plugin install (0.125s)",
		"Start: Claude Code subspace edge plugin removal",
		"Warning: Claude Code subspace edge plugin removal (0.125s; child exit 1; continuing)",
		"Start: Claude Code subspace plugin install",
		"Done: Claude Code subspace plugin install (0.125s)",
		"Start: Resolve Claude Code plugin manifest",
	})
	if lines := traceLines(stderr.String()); len(lines) == 0 || !strings.HasSuffix(lines[len(lines)-1], "(exit 0)") {
		t.Errorf("trace did not end with the exit 0 total:\n%s", &stderr)
	}
}

func TestInstallVerboseTracesAFailFastClaudeStep(t *testing.T) {
	dir := t.TempDir()
	script := "#!/bin/sh\ncase \"$*\" in\n" +
		"\"plugin list --json\") echo '[]' ;;\n" +
		"\"plugin marketplace add \"*) echo 'marketplace unreachable' >&2; exit 3 ;;\n" +
		"esac\nexit 0\n"
	mustWrite(t, filepath.Join(dir, "claude"), script)
	if err := os.Chmod(filepath.Join(dir, "claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+"/usr/bin:/bin")
	t.Setenv("HOME", t.TempDir())
	var stdout, stderr bytes.Buffer
	code := runInstallCommand(context.Background(), []string{"--verbose", "--host", "claude"}, os.Environ(), t.TempDir(), &stdout, &stderr, fixedStepClock())

	if code != 1 {
		t.Fatalf("install exited %d, want 1:\n%s", code, &stderr)
	}
	assertTraceInOrder(t, stderr.String(), []string{
		"Skipped: Claude Code retired Spacedock plugin removal (not installed at user scope)",
		"Start: Claude Code marketplace add",
		"Failed: Claude Code marketplace add (0.125s; child exit 3)",
		"Total: 0.875s (exit 1)",
	})
	if strings.Contains(stderr.String(), "marketplace refresh") {
		t.Errorf("the install continued past a fail-fast step:\n%s", &stderr)
	}
}
