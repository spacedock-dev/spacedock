// ABOUTME: Computes the effective binary display version — the linker-stamped
// ABOUTME: Version, or the embedded checkout manifest's version for a dev build.
package cli

import (
	"encoding/json"

	spacedock "github.com/spacedock-dev/spacedock"
)

// displayVersion returns the binary's effective display version, fed to every
// compatibility gate and version-bearing message. Only an explicitly marked
// release build returns its stamped Version unchanged. Every unmarked
// `go build`/`go install` instead reports the embedded checkout manifest's
// version with a `+dev` build tag appended (D3): the source binary always claims
// its checkout's minor, so a stale checkout is rejected by newer skills and a
// fresh one passes — no `dev` always-pass carve-out anywhere in the gate. A
// `go install …@vX.Y.Z` proxy build (module proxy, no ldflags) self-reports
// `X.Y.Z+dev`, since the tagged manifest equals the tag.
func displayVersion() string {
	if releaseBuild == "true" && Version != "dev" {
		return Version
	}
	v, ok := embeddedManifestVersion()
	if !ok {
		return Version
	}
	return v + "+dev"
}

// embeddedManifestVersion reads the `version` field out of the module-root
// go:embed of .claude-plugin/plugin.json (spacedock.PluginManifest) — the
// checkout's own manifest, embedded at compile time. Returns ok=false when the
// embedded bytes do not parse (should not happen for a real checkout: the embed
// is a compile-time constant of a file this repo owns).
func embeddedManifestVersion() (string, bool) {
	var m struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(spacedock.PluginManifest, &m); err != nil || m.Version == "" {
		return "", false
	}
	return m.Version, true
}
