// ABOUTME: Detect whether the running spacedock binary is the Homebrew edge
// ABOUTME: (spacedock@next) cask, so the too-old-binary remedy names that formula.
package cli

import (
	"path/filepath"
	"strings"
)

// runningEdgeCask reports whether the resolved running binary was installed from
// the Homebrew edge cask (`spacedock@next`). It anchors on the resolved executable
// (resolvedLauncherBin's os.Executable→EvalSymlinks), not the `spacedock` on PATH,
// so a session launched from a source checkout or the stable cask is not misread
// as edge.
func runningEdgeCask() bool {
	execPath, ok := resolvedLauncherBin()
	if !ok {
		return false
	}
	return isEdgeCaskPath(execPath)
}

// isEdgeCaskPath reports whether execPath resolves under a Homebrew
// `Caskroom/spacedock@next/` segment — the edge cask's staged location (e.g.
// `/opt/homebrew/Caskroom/spacedock@next/0.26.0-pre0/spacedock`). Any other path
// (the stable `spacedock` cask, a source checkout, an empty path) reports false,
// leaving the too-old-binary remedy on its unchanged `brew upgrade spacedock`
// default.
func isEdgeCaskPath(execPath string) bool {
	segments := strings.Split(filepath.ToSlash(execPath), "/")
	for i, seg := range segments {
		if seg == "Caskroom" && i+1 < len(segments) {
			return segments[i+1] == "spacedock@next"
		}
	}
	return false
}
