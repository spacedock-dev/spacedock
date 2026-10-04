// ABOUTME: The single source of truth for the live E2E lane model ids. The
// ABOUTME: workflow resolves them through `spacedock live-models` rather than
// ABOUTME: repeating literals.
package release

// The live lane model ids, each pinned once. The workflow matrix, the Codex
// exec shim, the Pi harness, and the docs all derive from these, so a model
// change is one edit here plus a print, never a hand-sync across five sites.
const (
	// ClaudeSonnetModel is the routine Claude lane model: pull requests and the
	// sonnet cadence.
	ClaudeSonnetModel = "claude-sonnet-5"
	// ClaudeOpusModel is the Claude pre-release lane model.
	ClaudeOpusModel = "claude-opus-4-8"
	// CodexExecModel is the bare id the Codex lane pins on `codex exec`.
	CodexExecModel = "gpt-6-luna"
	// PiOAuthModel is the Pi lane model for subscription (OAuth) auth. It passed
	// a local front-door smoke on 2026-10-04 (root and child both on it).
	PiOAuthModel = "openai-codex/gpt-6-luna:max"
	// PiAPIKeyModel is the Pi lane model for API-key auth: the `openai` provider
	// counterpart of PiOAuthModel, with the same `max` thinking suffix.
	PiAPIKeyModel = "openai/gpt-6-luna:max"
)

// LiveLaneModel is one lane/auth model entry. Key is the dotted lookup key the
// workflow passes to `spacedock live-models --get`; it names both the lane and,
// where a lane has more than one, the auth path or cadence.
type LiveLaneModel struct {
	Key string
	ID  string
}

// LiveModels returns the lane models in print order. It is the one list
// `spacedock live-models` prints, so the workflow and the Go harness resolve a
// lane model from one source.
func LiveModels() []LiveLaneModel {
	return []LiveLaneModel{
		{Key: "claude.sonnet", ID: ClaudeSonnetModel},
		{Key: "claude.opus", ID: ClaudeOpusModel},
		{Key: "codex.exec", ID: CodexExecModel},
		{Key: "pi.oauth", ID: PiOAuthModel},
		{Key: "pi.api-key", ID: PiAPIKeyModel},
	}
}

// LiveModel returns the id for key, or ("", false) when key names no lane.
func LiveModel(key string) (string, bool) {
	for _, m := range LiveModels() {
		if m.Key == key {
			return m.ID, true
		}
	}
	return "", false
}
