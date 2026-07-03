// ABOUTME: AC-2 oracle (3) behavior fixture — a real spacedock stub whose
// ABOUTME: --version prints a chosen binary version drives the startup gate.
package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestStartupGateAbortsBeforeDiscover builds a real `spacedock` stub whose
// `--version` prints a chosen version on line 1, runs the stub for its version
// output, drives the startup gate against a required minor literal, and
// observes: an out-of-range minor aborts with the pinned message and the stub's
// `status --discover` / `--boot` subcommands are NEVER invoked; a same-minor
// version proceeds and the discover call fires exactly once.
func TestStartupGateAbortsBeforeDiscover(t *testing.T) {
	stub, marker := buildVersionStub(t)

	cases := []struct {
		name          string
		stubVersion   string // value the stub prints as its version token
		requiredMinor string // the FO prose's stamped major.minor literal
		wantProceed   bool
		wantPinned    string // abort-message substring (empty when proceeding)
	}{
		{"too-old-binary-aborts", "0.23.0", "0.24", false, "Upgrade the binary to continue."},
		{"too-old-plugin-aborts", "0.25.0", "0.24", false, "Update the plugin to continue."},
		{"compatible-proceeds", "0.24.3", "0.24", true, ""},
		{"compatible-prerelease-proceeds", "0.24.0-pre1", "0.24", true, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			discoverCalls := 0
			stubEnv := append(os.Environ(),
				"SD_STUB_VERSION="+c.stubVersion,
				"SD_STUB_MARKER="+marker,
			)

			// runVersion executes the real stub's --version.
			runVersion := func() (string, error) {
				cmd := exec.Command(stub, "--version")
				cmd.Env = stubEnv
				out, err := cmd.Output()
				return string(out), err
			}
			// runDiscover executes the stub's status --discover and counts the call.
			runDiscover := func() error {
				discoverCalls++
				cmd := exec.Command(stub, "status", "--discover")
				cmd.Env = stubEnv
				return cmd.Run()
			}

			proceed, msg := gateAndMaybeDiscover(runVersion, c.requiredMinor, "claude", runDiscover)

			if proceed != c.wantProceed {
				t.Fatalf("proceed = %v, want %v (msg=%q)", proceed, c.wantProceed, msg)
			}
			if c.wantProceed {
				if discoverCalls != 1 {
					t.Fatalf("compatible gate: discover called %d times, want 1", discoverCalls)
				}
			} else {
				if discoverCalls != 0 {
					t.Fatalf("aborting gate invoked discover %d times, want 0", discoverCalls)
				}
				if !strings.Contains(msg, c.wantPinned) {
					t.Fatalf("abort message %q missing pinned remedy %q", msg, c.wantPinned)
				}
			}

			// Sanity: the stub records which subcommands it actually ran. On an
			// aborting gate the marker file must show only the version call.
			if !c.wantProceed {
				record := readMarker(t, marker)
				if strings.Contains(record, "discover") {
					t.Fatalf("stub was invoked with discover on an aborting gate; marker:\n%s", record)
				}
			}
			// Reset the marker between cases.
			os.Remove(marker)
		})
	}
}

// gateAndMaybeDiscover realizes the FO Startup step-1 gate as a callable
// mechanism: run the version probe, parse the `spacedock <version>` line-1
// token, compare its major.minor against the release-stamped required minor,
// and only call discover when compatible. This is the Go realization of the
// prose the FO follows — driven here by a real stub process, not a mock.
func gateAndMaybeDiscover(runVersion func() (string, error), requiredMinor, host string, runDiscover func() error) (proceed bool, message string) {
	out, err := runVersion()
	if err != nil {
		return false, "spacedock --version unavailable: " + err.Error()
	}
	binaryVersion, ok := parseVersionLine(out)
	if !ok {
		return false, "could not parse `spacedock <version>` from `spacedock --version`: " + strings.TrimSpace(out)
	}
	// requiredMinor (e.g. "0.24") stands in for the plugin's declared version —
	// synthesized to major.minor.0 so Compare's ordinary major.minor comparison
	// realizes the FO prose gate exactly as the manifest-based doctor/front-door
	// gates do.
	res := Compare(host, requiredMinor+".0", binaryVersion)
	if res.Verdict != Compatible {
		return false, res.Message
	}
	_ = runDiscover()
	return true, res.Message
}

// buildVersionStub compiles a tiny stub binary that prints a version token from
// the SD_STUB_VERSION env var on `--version` and records every subcommand it is
// invoked with into the file named by SD_STUB_MARKER. Returns the binary path
// and the marker path.
func buildVersionStub(t *testing.T) (binPath, markerPath string) {
	t.Helper()
	dir := t.TempDir()
	src := filepath.Join(dir, "main.go")
	markerPath = filepath.Join(dir, "marker.txt")
	stubSrc := `package main

import (
	"fmt"
	"os"
)

func main() {
	if f, err := os.OpenFile(os.Getenv("SD_STUB_MARKER"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
		for _, a := range os.Args[1:] {
			fmt.Fprintln(f, a)
		}
		f.Close()
	}
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("spacedock %s (contract 3)\n", os.Getenv("SD_STUB_VERSION"))
	}
}
`
	if err := os.WriteFile(src, []byte(stubSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	binPath = filepath.Join(dir, "spacedock-stub")
	cmd := exec.Command("go", "build", "-o", binPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build version stub: %v\n%s", err, out)
	}
	return binPath, markerPath
}

func readMarker(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

// parseVersionLine extracts the <version> token from a `spacedock --version`
// line 1 of the shape `spacedock <version> (contract <N>)`, the same parse the
// FO Startup step-1 prose does.
func parseVersionLine(versionOut string) (string, bool) {
	line := versionOut
	if i := strings.IndexByte(line, '\n'); i >= 0 {
		line = line[:i]
	}
	const prefix = "spacedock "
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(line, prefix)
	end := strings.IndexByte(rest, ' ')
	if end < 0 {
		end = len(rest)
	}
	version := strings.TrimSpace(rest[:end])
	if version == "" {
		return "", false
	}
	return version, true
}
