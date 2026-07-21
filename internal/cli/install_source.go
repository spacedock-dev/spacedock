// ABOUTME: Detect how the running spacedock binary was installed (brew stable,
// ABOUTME: brew edge @next, or non-brew) so the too-old-binary remedy fits the source.
package cli

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// currentInstallSource detects the install source of the running binary, for the
// too-old-binary remedy. It anchors on the resolved running executable (not the
// `spacedock` on PATH, which can be a different install) and probes the real
// `brew` on PATH — the thin production wiring the front door and doctor callers
// share. Detection stays injectable via detectInstallSource for tests.
func currentInstallSource() contract.InstallSource {
	execPath, _ := resolvedLauncherBin()
	return detectInstallSource(execPath, exec.LookPath, devBranch)
}

// detectInstallSource classifies how the running binary was installed, so the
// too-old-binary remedy names the upgrade path that actually applies. It anchors
// on execPath — the resolved RUNNING binary (resolvedLauncherBin's os.Executable
// →EvalSymlinks result) — not `command -v spacedock`, because a session launched
// from a source checkout must classify as its own source, not the brew binary on
// PATH. A `…/Caskroom/<token>/…` path is a Homebrew install named by <token>
// (`spacedock`→stable, `spacedock@next`→edge); brewLookPath resolving `brew` is
// what separates an in-place upgrade from the run-on-host (HostOnly) sandbox case.
// A resolved non-Caskroom path is a non-brew build; an empty/unresolvable path is
// the generic fallback. devBranch is threaded for the caller's channel stamp but
// does NOT pick the formula — the real cask token does (a plain `go build` stamps
// `next`, which would misread a source build as edge).
func detectInstallSource(execPath string, brewLookPath func(string) (string, error), devBranch string) contract.InstallSource {
	token, isCask := caskToken(execPath)
	if !isCask {
		if strings.TrimSpace(execPath) == "" {
			return contract.InstallSource{Kind: contract.SourceUnknown}
		}
		return contract.InstallSource{Kind: contract.NonBrew}
	}
	var kind contract.SourceKind
	switch token {
	case "spacedock":
		kind = contract.BrewStable
	case "spacedock@next":
		kind = contract.BrewEdge
	default:
		// A Caskroom path with an unrecognized cask token: fall back to the
		// generic remedy rather than guessing a formula name.
		return contract.InstallSource{Kind: contract.SourceUnknown}
	}
	hostOnly := false
	if _, err := brewLookPath("brew"); err != nil {
		hostOnly = true
	}
	return contract.InstallSource{Kind: kind, HostOnly: hostOnly}
}

// caskToken extracts the Homebrew cask token that owns execPath: the path segment
// immediately after a `Caskroom` segment (e.g.
// `/opt/homebrew/Caskroom/spacedock@next/0.26.0-pre0/spacedock` → `spacedock@next`).
// It reports found=false when the path is empty, carries no `Caskroom` segment, or
// ends at `Caskroom` with no token following — every non-brew shape.
func caskToken(execPath string) (token string, found bool) {
	if strings.TrimSpace(execPath) == "" {
		return "", false
	}
	segments := strings.Split(filepath.ToSlash(execPath), "/")
	for i, seg := range segments {
		if seg == "Caskroom" && i+1 < len(segments) {
			return segments[i+1], true
		}
	}
	return "", false
}
