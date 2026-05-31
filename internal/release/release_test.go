// ABOUTME: Release-pipeline version steps — AC-4 plugin.json stamp + AC-2d
// ABOUTME: marketplace-entry calendar bump, as pure functions a CI step invokes.
package release

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

// TestStampVersionRewritesPluginVersion locks AC-4: StampVersion rewrites the
// top-level plugin.json `version` to the release value, leaving every other
// field (and the file's formatting) intact.
func TestStampVersionRewritesPluginVersion(t *testing.T) {
	src := `{
  "name": "spacedock",
  "version": "0.1.0-dev",
  "skills": "./skills/",
  "requires-contract": ">=1,<2"
}
`
	out, err := StampVersion([]byte(src), "0.19.0")
	if err != nil {
		t.Fatalf("StampVersion: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("stamped manifest does not parse: %v\n%s", err, out)
	}
	if m["version"] != "0.19.0" {
		t.Errorf("version = %v, want 0.19.0", m["version"])
	}
	// Untouched fields survive.
	if m["name"] != "spacedock" {
		t.Errorf("name field lost: %v", m["name"])
	}
	if m["requires-contract"] != ">=1,<2" {
		t.Errorf("requires-contract field lost: %v", m["requires-contract"])
	}
	if m["skills"] != "./skills/" {
		t.Errorf("skills field lost: %v", m["skills"])
	}
}

// TestStampVersionLeavesMarketplaceCalendarUntouched locks AC-4's negative half:
// the stamp step targets the top-level plugin `version` only. Applied to a
// marketplace.json (whose meaningful version lives on the nested plugin ENTRY,
// not at top level), it must NOT rewrite the entry's calendar version. The stamp
// is a plugin.json operation; a marketplace.json has no top-level version to
// stamp, so the entry calendar key is left exactly as-is.
func TestStampVersionLeavesMarketplaceCalendarUntouched(t *testing.T) {
	src := `{
  "name": "spacedock",
  "plugins": [
    {
      "name": "spacedock",
      "version": "0.0.2026053101"
    }
  ]
}
`
	out, err := StampVersion([]byte(src), "0.19.0")
	if err != nil {
		t.Fatalf("StampVersion: %v", err)
	}
	if strings.Contains(string(out), "0.19.0") {
		t.Errorf("stamp wrote the release version into a marketplace.json entry; calendar key must be untouched:\n%s", out)
	}
	if !strings.Contains(string(out), "0.0.2026053101") {
		t.Errorf("marketplace entry calendar version was lost:\n%s", out)
	}
}

// TestStampVersionRewritesOnlyFirstVersionKey locks the replace-first contract:
// a manifest carrying a second nested `version` key (beyond the top-level one)
// must have ONLY its top-level version rewritten — the nested key is left
// exactly as-is. A blanket ReplaceAll over every `"version":` match would clobber
// the nested key too; the targeted first-match replace must not.
func TestStampVersionRewritesOnlyFirstVersionKey(t *testing.T) {
	src := `{
  "name": "spacedock",
  "version": "0.1.0-dev",
  "skills": "./skills/",
  "metadata": {
    "version": "schema-7"
  }
}
`
	out, err := StampVersion([]byte(src), "0.19.0")
	if err != nil {
		t.Fatalf("StampVersion: %v", err)
	}
	var m struct {
		Version  string `json:"version"`
		Metadata struct {
			Version string `json:"version"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("stamped manifest does not parse: %v\n%s", err, out)
	}
	if m.Version != "0.19.0" {
		t.Errorf("top-level version = %q, want 0.19.0", m.Version)
	}
	if m.Metadata.Version != "schema-7" {
		t.Errorf("nested metadata.version was clobbered: %q, want schema-7 (replace-first must leave it untouched)", m.Metadata.Version)
	}
}

// TestBumpCalendarVersionStrictlyIncreases locks AC-2d: invoking the bump
// function twice over the SAME marketplace.json produces a strictly increasing
// entry version (the `plugin update` re-pull key actually moves), not two
// hand-written literals. The second call (same day) must increment the per-day
// sequence so the value is monotonic even within a single publish day.
func TestBumpCalendarVersionStrictlyIncreases(t *testing.T) {
	src := `{
  "name": "spacedock",
  "plugins": [
    {
      "name": "spacedock",
      "source": { "source": "url", "url": "https://example/spacedock.git", "ref": "next" },
      "version": "0.0.2026053101",
      "category": "workflow"
    }
  ]
}
`
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)

	first, err := BumpCalendarVersion([]byte(src), now)
	if err != nil {
		t.Fatalf("first bump: %v", err)
	}
	second, err := BumpCalendarVersion(first, now)
	if err != nil {
		t.Fatalf("second bump: %v", err)
	}

	v1 := entryVersion(t, first)
	v2 := entryVersion(t, second)
	if !(v2 > v1) {
		t.Errorf("calendar version did not strictly increase: %q then %q", v1, v2)
	}
	// The bump is calendar-keyed: same day -> shared date prefix, incremented seq.
	if !strings.HasPrefix(v1, "0.0.20260601") {
		t.Errorf("first bump = %q, want 0.0.20260601NN prefix", v1)
	}
	if !strings.HasPrefix(v2, "0.0.20260601") {
		t.Errorf("second bump = %q, want 0.0.20260601NN prefix", v2)
	}
}

// TestBumpCalendarVersionRewritesOnlyFirstVersionKey locks the replace-first
// contract for the calendar bump: a marketplace.json whose plugin entry carries a
// trailing second `version` key (e.g. a nested source/metadata block) must have
// ONLY the entry's first `version` advanced — the trailing key stays intact. A
// blanket ReplaceAll would over-rewrite it.
func TestBumpCalendarVersionRewritesOnlyFirstVersionKey(t *testing.T) {
	src := `{
  "name": "spacedock",
  "plugins": [
    {
      "name": "spacedock",
      "version": "0.0.2026053101",
      "schema": { "version": "marketplace-2" }
    }
  ]
}
`
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	out, err := BumpCalendarVersion([]byte(src), now)
	if err != nil {
		t.Fatalf("bump: %v", err)
	}
	if v := entryVersion(t, out); v != "0.0.2026060101" {
		t.Errorf("entry calendar version = %q, want 0.0.2026060101", v)
	}
	if !strings.Contains(string(out), `"version": "marketplace-2"`) {
		t.Errorf("nested schema.version was clobbered; replace-first must leave it untouched:\n%s", out)
	}
}

// TestBumpCalendarVersionNewDayResetsSequence locks the cross-day behavior: a
// bump on a later date produces a strictly greater value than a prior-day bump
// even though the new day's sequence restarts at 01 — the date component
// dominates the ordering.
func TestBumpCalendarVersionNewDayResetsSequence(t *testing.T) {
	src := `{"name":"spacedock","plugins":[{"name":"spacedock","version":"0.0.2026053199"}]}`
	day2 := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	out, err := BumpCalendarVersion([]byte(src), day2)
	if err != nil {
		t.Fatalf("bump: %v", err)
	}
	v := entryVersion(t, out)
	if v != "0.0.2026060101" {
		t.Errorf("new-day bump = %q, want 0.0.2026060101 (seq resets to 01 on a new date)", v)
	}
	if !(v > "0.0.2026053199") {
		t.Errorf("new-day bump %q did not exceed prior-day max 0.0.2026053199", v)
	}
}

// TestFilterCommitLogDropsWorkflowNoise locks AC-1: FilterCommitLog removes the
// workflow-state commit classes the changelog prompt is told to ignore
// (`dispatch:`/`advance:`/`merge:`/`archive:` subjects plus the CI commits
// `release: stamp …` and `next: bump …`) while preserving real user-facing
// commits, so both the `claude` input and the no-claude fallback output are
// already clean. The input is the `git log --oneline` shape (`<sha> <subject>`).
func TestFilterCommitLogDropsWorkflowNoise(t *testing.T) {
	raw := strings.Join([]string{
		"a28bb9d feat(ci): wire approval-gated live-runtime e2e job",
		"ff98605 dispatch: release-notes-local-summary entering implementation",
		"41e76a7 fix(dispatch): bind resolved state checkout",
		"fcee0ed release: stamp plugin manifests to 0.19.1",
		"23aa4d7 next: bump marketplace calendar version",
		"e00d3f3 advance: cli-cobra-redesign entering validation",
		"22a42ec merge: release-pipeline (jf) — goreleaser homebrew_casks",
		"6eeaad1 archive: release-notes-local-summary",
		"080ec3e feat(plugin): vendor repo plugin manifest",
	}, "\n")

	got := FilterCommitLog(raw)

	// Real commits survive — including ones whose subject merely mentions a
	// noise word inside a scope (`fix(dispatch): …`) rather than as the prefix.
	keep := []string{
		"feat(ci): wire approval-gated live-runtime e2e job",
		"fix(dispatch): bind resolved state checkout",
		"feat(plugin): vendor repo plugin manifest",
	}
	for _, want := range keep {
		if !strings.Contains(got, want) {
			t.Errorf("FilterCommitLog dropped a real commit %q:\n%s", want, got)
		}
	}

	// Each workflow-noise class is gone.
	drop := []string{
		"dispatch: release-notes-local-summary",
		"release: stamp plugin manifests",
		"next: bump marketplace calendar",
		"advance: cli-cobra-redesign",
		"merge: release-pipeline",
		"archive: release-notes-local-summary",
	}
	for _, bad := range drop {
		if strings.Contains(got, bad) {
			t.Errorf("FilterCommitLog kept workflow-noise commit %q:\n%s", bad, got)
		}
	}
}

// TestBuildChangelogPromptShape locks AC-1's prompt half: the prompt names the
// release version, demands plain text (no markdown), and its ignore-list names
// this repo's workflow-noise classes so the LLM is told to drop them.
func TestBuildChangelogPromptShape(t *testing.T) {
	p := BuildChangelogPrompt("0.19.2")
	if !strings.Contains(p, "0.19.2") {
		t.Errorf("prompt does not name the version:\n%s", p)
	}
	if !strings.Contains(p, "Plain text only") {
		t.Errorf("prompt does not demand plain text:\n%s", p)
	}
	for _, noise := range []string{"dispatch", "advance", "merge", "archive"} {
		if !strings.Contains(p, noise) {
			t.Errorf("prompt ignore-list omits the %q noise class:\n%s", noise, p)
		}
	}
}

// TestGenerateNotesFallsBackWithoutClaude locks AC-1's fallback: when the claude
// runner returns an error (binary absent), GenerateNotes still produces notes —
// the filtered raw log — rather than failing.
func TestGenerateNotesFallsBackWithoutClaude(t *testing.T) {
	raw := "a28bb9d feat(ci): real change\nff98605 dispatch: noise"
	io := NotesIO{
		RawLog: func() (string, error) { return raw, nil },
		Claude: func(prompt, input string) (string, error) { return "", errors.New("claude: not found") },
	}
	notes, err := GenerateNotes("0.19.2", io)
	if err != nil {
		t.Fatalf("GenerateNotes: %v", err)
	}
	if !strings.Contains(notes, "feat(ci): real change") {
		t.Errorf("fallback notes lost the real commit:\n%s", notes)
	}
	if strings.Contains(notes, "dispatch: noise") {
		t.Errorf("fallback notes kept workflow noise:\n%s", notes)
	}
}

// TestGenerateNotesUsesClaudeOnFilteredLog locks AC-1: when claude is available,
// it is fed the FILTERED log (not the raw one) and its output is the notes.
func TestGenerateNotesUsesClaudeOnFilteredLog(t *testing.T) {
	raw := "a28bb9d feat(ci): real change\nff98605 dispatch: noise"
	var sawInput string
	io := NotesIO{
		RawLog: func() (string, error) { return raw, nil },
		Claude: func(prompt, input string) (string, error) {
			sawInput = input
			return "A summary release.\n- feat: real change", nil
		},
	}
	notes, err := GenerateNotes("0.19.2", io)
	if err != nil {
		t.Fatalf("GenerateNotes: %v", err)
	}
	if strings.Contains(sawInput, "dispatch: noise") {
		t.Errorf("claude was fed unfiltered noise:\n%s", sawInput)
	}
	if !strings.Contains(sawInput, "feat(ci): real change") {
		t.Errorf("claude was not fed the real commit:\n%s", sawInput)
	}
	if notes != "A summary release.\n- feat: real change" {
		t.Errorf("notes != claude output: %q", notes)
	}
}

// TestConfirmAndTagDeclineCutsNoTag locks AC-2's negative half: when the captain
// declines, the tag hook is never invoked.
func TestConfirmAndTagDeclineCutsNoTag(t *testing.T) {
	tagged := false
	io := TagIO{
		Confirm: func(proposed string) (string, bool) { return proposed, false },
		CutTag:  func(body string) error { tagged = true; return nil },
	}
	if err := ConfirmAndTag("notes body", io); err != nil {
		t.Fatalf("ConfirmAndTag: %v", err)
	}
	if tagged {
		t.Errorf("tag was cut despite the captain declining")
	}
}

// TestConfirmAndTagConfirmCutsEditedBody locks AC-2's positive half: on confirm,
// the tag hook is called with the captain's EDITED body, not the proposed one.
func TestConfirmAndTagConfirmCutsEditedBody(t *testing.T) {
	var gotBody string
	io := TagIO{
		Confirm: func(proposed string) (string, bool) { return "captain-edited body", true },
		CutTag:  func(body string) error { gotBody = body; return nil },
	}
	if err := ConfirmAndTag("proposed body", io); err != nil {
		t.Fatalf("ConfirmAndTag: %v", err)
	}
	if gotBody != "captain-edited body" {
		t.Errorf("tag cut with body %q, want the captain-edited body", gotBody)
	}
}

// entryVersion extracts plugins[0].version from a marketplace.json blob.
func entryVersion(t *testing.T, blob []byte) string {
	t.Helper()
	var m struct {
		Plugins []struct {
			Version string `json:"version"`
		} `json:"plugins"`
	}
	if err := json.Unmarshal(blob, &m); err != nil {
		t.Fatalf("parse marketplace blob: %v\n%s", err, blob)
	}
	if len(m.Plugins) == 0 {
		t.Fatalf("no plugins entry in blob:\n%s", blob)
	}
	return m.Plugins[0].Version
}
