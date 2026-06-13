// ABOUTME: Drives install.sh's checksum gate over a local dist fixture, proving
// ABOUTME: a tampered tarball aborts and that the gate lines are load-bearing.
package release

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// installFixture is a local goreleaser-shaped `dist/` directory: one
// `spacedock_<ver>_<os>_<arch>.tar.gz` holding a bare runnable `spacedock` at the
// archive root, plus a `checksums.txt` line matching that tarball — the same
// layout install.sh's SPACEDOCK_INSTALL_FROM=<dir> path consumes.
type installFixture struct {
	dir         string // the dist dir to point SPACEDOCK_INSTALL_FROM at
	tarballPath string // absolute path to the single os/arch tarball
	asset       string // the tarball's basename
	marker      string // the string the installed binary prints when run
}

// buildInstallFixture writes a dist/ fixture under a fresh temp dir. The bare
// `spacedock` payload is a tiny shell script that echoes a unique marker, so the
// happy-path test can exec the installed file and confirm a RUNNABLE binary
// landed (not just any file). checksums.txt is computed from the real tarball
// bytes, so the gate's expected hash is correct until a test mutates the tarball.
func buildInstallFixture(t *testing.T) installFixture {
	t.Helper()
	os, arch := goosArch(t)
	dist := t.TempDir()

	const marker = "spacedock-fixture-ran-ok"
	// A bare executable `spacedock` at the archive root. A shell script is enough
	// for install.sh (it `install`s the file 0755) and lets the test exec it.
	binary := "#!/bin/sh\necho " + marker + "\n"

	asset := "spacedock_0.0.0_" + os + "_" + arch + ".tar.gz"
	tarballPath := filepath.Join(dist, asset)
	writeTarGz(t, tarballPath, "spacedock", []byte(binary))

	sum := sha256OfFile(t, tarballPath)
	// goreleaser's checksums.txt format is `<sha256>␣␣<filename>`; install.sh
	// parses it with `awk '$2 == filename {print $1}'`, so two space-separated
	// fields suffice.
	checksums := sum + "  " + asset + "\n"
	if err := osWriteFile(filepath.Join(dist, "checksums.txt"), checksums); err != nil {
		t.Fatal(err)
	}

	return installFixture{dir: dist, tarballPath: tarballPath, asset: asset, marker: marker}
}

// runInstall runs the given install.sh script against a dist fixture via the
// SPACEDOCK_INSTALL_FROM local-dist override, installing into a fresh dir. It
// returns the install dir and the script's exit code (0 on success). The env is
// set explicitly (not scrubbed) because this path REQUIRES the override.
func runInstall(t *testing.T, script, distDir string) (installDir string, exitCode int) {
	t.Helper()
	installDir = filepath.Join(t.TempDir(), "bin")
	cmd := exec.Command("sh", script)
	cmd.Env = append(scrubInstallEnv(),
		"SPACEDOCK_INSTALL_FROM="+distDir,
		"SPACEDOCK_INSTALL_DIR="+installDir,
	)
	out, err := cmd.CombinedOutput()
	t.Logf("install.sh (%s) output:\n%s", filepath.Base(script), out)
	if err == nil {
		return installDir, 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return installDir, ee.ExitCode()
	}
	t.Fatalf("install.sh failed to launch: %v", err)
	return installDir, -1
}

// TestChecksumGateInstallsAndRejectsTamper locks AC-1: install.sh's checksum gate
// installs a runnable binary on the happy path and aborts (installing nothing) on
// a byte-tampered tarball. Driven over a local dist fixture, no goreleaser.
func TestChecksumGateInstallsAndRejectsTamper(t *testing.T) {
	script := filepath.Join("..", "..", "install.sh")

	// Happy path: the fixture's checksums.txt matches the tarball, so the gate
	// passes and a runnable `spacedock` lands and prints its marker.
	t.Run("happy path installs a runnable binary", func(t *testing.T) {
		fx := buildInstallFixture(t)
		installDir, code := runInstall(t, script, fx.dir)
		if code != 0 {
			t.Fatalf("happy-path install exited %d, want 0", code)
		}
		out, err := exec.Command(filepath.Join(installDir, "spacedock")).CombinedOutput()
		if err != nil {
			t.Fatalf("installed spacedock did not run: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), fx.marker) {
			t.Errorf("installed binary printed %q, want it to contain %q", out, fx.marker)
		}
	})

	// Tamper case: append bytes to the tarball so its sha256 no longer matches the
	// (unchanged) checksums.txt line. The gate MUST abort non-zero, installing
	// nothing. This is the assertion that reds if install.sh:164-169 are deleted.
	t.Run("tampered tarball aborts installing nothing", func(t *testing.T) {
		fx := buildInstallFixture(t)
		appendBytes(t, fx.tarballPath, []byte{0xde, 0xad, 0xbe, 0xef})

		installDir, code := runInstall(t, script, fx.dir)
		if code == 0 {
			t.Fatal("install.sh accepted a tampered tarball (exit 0); the checksum gate is not fail-closed")
		}
		if _, err := os.Stat(filepath.Join(installDir, "spacedock")); err == nil {
			t.Error("install.sh installed a binary despite the checksum mismatch")
		}
	})
}

// TestChecksumGateGuardIsLoadBearing proves the tamper assertion above actually
// exercises the gate: it strips the checksum-gate lines (install.sh:164-169) to a
// temp copy of the script, runs the SAME tamper case against THAT copy, and
// asserts the gateless installer now WRONGLY exits 0 and installs the tampered
// binary. If stripping the gate did NOT change behavior, the live tamper test
// wasn't binding the gate — that is the hole this load-bearing check closes.
func TestChecksumGateGuardIsLoadBearing(t *testing.T) {
	original := readInstallScript(t)
	gateless, removed := stripChecksumGate(original)
	if !removed {
		t.Fatalf("could not locate the checksum-gate block in install.sh to strip; the load-bearing check cannot bind")
	}

	gatelessScript := filepath.Join(t.TempDir(), "install-gateless.sh")
	if err := osWriteFile(gatelessScript, gateless); err != nil {
		t.Fatal(err)
	}

	fx := buildInstallFixture(t)
	appendBytes(t, fx.tarballPath, []byte{0xde, 0xad, 0xbe, 0xef})

	installDir, code := runInstall(t, gatelessScript, fx.dir)
	if code != 0 {
		t.Fatalf("gateless install.sh exited %d on a tampered tarball; expected 0 (the strip should let the tamper through), so the live tamper test is NOT exercising the gate", code)
	}
	if _, err := os.Stat(filepath.Join(installDir, "spacedock")); err != nil {
		t.Errorf("gateless install.sh did not install the tampered binary (%v); the strip removed more than the gate", err)
	}
}

// stripChecksumGate removes install.sh's checksum-verification block — the
// `expected=…` extraction, its empty-check `die`, the `actual=…` hash, and the
// mismatch `if … die … fi` — yielding a script that extracts and installs WITHOUT
// verifying. It returns the rewritten text and whether the block was found, so a
// drift in install.sh's gate shape fails the load-bearing test loudly rather than
// silently checking nothing.
func stripChecksumGate(script string) (string, bool) {
	lines := strings.Split(script, "\n")
	start, end := -1, -1
	for i, line := range lines {
		if start == -1 && strings.HasPrefix(strings.TrimSpace(line), "expected=") {
			start = i
		}
		// The gate ends at the closing `fi` of the mismatch check.
		if start != -1 && strings.TrimSpace(line) == "fi" {
			end = i
			break
		}
	}
	if start == -1 || end == -1 {
		return script, false
	}
	kept := append(append([]string{}, lines[:start]...), lines[end+1:]...)
	return strings.Join(kept, "\n"), true
}

func readInstallScript(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// writeTarGz writes a gzip-compressed tar at path containing a single file with
// the given name (at the archive root) and contents, mode 0755.
func writeTarGz(t *testing.T, path, name string, content []byte) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	hdr := &tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
}

func sha256OfFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func appendBytes(t *testing.T, path string, extra []byte) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.Write(extra); err != nil {
		t.Fatal(err)
	}
}

func osWriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
