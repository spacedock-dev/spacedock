// ABOUTME: Computes the effective binary display version — the linker-stamped
// ABOUTME: Version, or the embedded checkout manifest's version for a dev build.
package cli

import (
	"encoding/json"
	"runtime/debug"
	"strings"

	spacedock "github.com/spacedock-dev/spacedock"
)

var readBuildInfo = debug.ReadBuildInfo

// displayVersion returns the binary's effective display version, fed to every
// compatibility gate and version-bearing message. A linker-stamped release
// build returns Version unchanged. An unstamped `go build`/`go install` dev
// build (Version == "dev") instead reports the embedded checkout manifest's
// version with `+dev[.<sha7>[.dirty]]` build metadata: the dev binary always claims
// its checkout's minor, so a stale checkout is rejected by newer skills and a
// fresh one passes — no `dev` always-pass carve-out anywhere in the gate. A
// `go install …@vX.Y.Z` proxy build (module proxy, no ldflags) self-reports
// `X.Y.Z+dev`, since the tagged manifest equals the tag.
func displayVersion() string {
	if Version != "dev" {
		return Version
	}
	v, ok := embeddedManifestVersion()
	if !ok {
		return Version
	}
	version := v + "+dev"
	info, ok := readBuildInfo()
	if !ok {
		return version
	}
	var revision, modified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value
		}
	}
	if len(revision) < 7 {
		return version
	}
	for _, c := range revision {
		if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
			return version
		}
	}
	version += "." + strings.ToLower(revision[:7])
	if modified == "true" {
		version += ".dirty"
	}
	return version
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
