// ABOUTME: Embeds the vendored .claude-plugin/plugin.json so an unstamped dev
// ABOUTME: build can report its own checkout's version (D3 dev-build gating).
package spacedock

import _ "embed"

// PluginManifest is the bytes of .claude-plugin/plugin.json, embedded here at
// the module root because go:embed cannot reach a file above the embedding
// package's directory. An unstamped `go build`/`go install` binary (no ldflags)
// reads this to report `<manifest-version>+dev` instead of the bare `dev`
// sentinel, so a source build always claims its checkout's minor for the
// compatibility gate.
//
//go:embed .claude-plugin/plugin.json
var PluginManifest []byte
