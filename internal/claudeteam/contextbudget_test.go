// ABOUTME: claudeteam unit tests — the AC-4 1M family-rule boundary table and the
// ABOUTME: context-budget envelope rendering, exercised directly in-package.
package claudeteam

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestContextLimitForModelBoundary is the AC-4 boundary table: the forward opus
// family rule (minor >= 7 -> 1M), the claude-{sonnet|fable|opus}-{major} family
// rule (major >= 5 -> 1M), and the [1m] suffix override, exercised at the
// version boundaries so a future release in either family stays correct
// without a code change.
func TestContextLimitForModelBoundary(t *testing.T) {
	cases := []struct {
		model string
		want  int
	}{
		{"claude-opus-4-8", extendedContextLimit},     // the live false-negative, now 1M
		{"claude-opus-4-8[1m]", extendedContextLimit}, // explicit suffix
		{"claude-opus-4-7", extendedContextLimit},     // first 1M-default minor
		{"claude-opus-4-6", defaultContextLimit},      // pre-default minor stays 200k
		{"claude-opus-4-6[1m]", extendedContextLimit}, // 4-6 with the suffix opts in
		{"claude-opus-4-10", extendedContextLimit},    // forward-safe: never goes stale
		{"claude-opus-4-100", extendedContextLimit},   // multi-digit minor
		{"claude-sonnet-4-6", defaultContextLimit},    // pre-5 sonnet, non-opus
		{"claude-haiku-4-5", defaultContextLimit},     // non-opus
		{"some-unknown-model", defaultContextLimit},   // safe fallback
		{"claude-opus-4", defaultContextLimit},        // no minor token -> no match

		// AC-2/AC-3/AC-4: the claude-{sonnet|fable|opus}-{major} family rule.
		{"claude-sonnet-5", extendedContextLimit},           // AC-2: sonnet-5 is 1M
		{"claude-sonnet-5[1m]", extendedContextLimit},       // explicit suffix, never observed live but consistent
		{"claude-sonnet-5-20260301", extendedContextLimit},  // hypothetical dated 5-family shape
		{"claude-sonnet-6", extendedContextLimit},           // forward-safe: next sonnet generation
		{"claude-fable-5", extendedContextLimit},            // AC-3: fable-5 is 1M
		{"claude-opus-5", extendedContextLimit},             // hypothetical opus 5th generation (no "4-" minor token)
		{"claude-sonnet-4-5-20250929", defaultContextLimit}, // dated pre-5 shape (major 4) stays 200k
		{"claude-haiku-5", defaultContextLimit},             // haiku deliberately excluded from the family rule
		{"claude-opus-4-20250514", extendedContextLimit},    // pre-existing quirk: the date token parses as minor >= 7 -> 1M (harness registry says 200k for this deprecated id; pinned as out of this entity's scope, not a change here)
	}
	for _, tc := range cases {
		if got := contextLimitForModel(tc.model); got != tc.want {
			t.Errorf("contextLimitForModel(%q) = %d, want %d", tc.model, got, tc.want)
		}
	}
}

// TestContextBudgetEnvelopeWholePercent asserts a whole-number usage_pct renders
// with the trailing .0 Python's json.dumps emits (20.0, not 20) — the pyFloat
// rendering the parity harness depends on, exercised directly without a fixture.
func TestContextBudgetEnvelopeWholePercent(t *testing.T) {
	home := t.TempDir()
	// 40k resident on a 200k opus-4-6 member -> exactly 20.0%.
	writeBudgetFixture(t, home, "ensign-w", "claude-opus-4-6", 40000)

	var stdout, stderr bytes.Buffer
	if code := ContextBudget(home, "ensign-w", &stdout, &stderr); code != 0 {
		t.Fatalf("ContextBudget exit=%d stderr=%q", code, stderr.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte(`"usage_pct": 20.0`)) {
		t.Errorf("usage_pct did not render as 20.0 (Python float repr):\n%s", stdout.String())
	}
}

// TestContextBudgetHealthySuffixDropNoDrift is AC-1: a healthy member whose team
// config carries the captain-session [1m] suffix while the runtime jsonl reports
// the canonical id emits no config_drift_warning and reads the 1M window — the
// config's declared suffix names the window the session genuinely runs.
func TestContextBudgetHealthySuffixDropNoDrift(t *testing.T) {
	home := t.TempDir()
	// config claude-fable-5[1m], runtime jsonl claude-fable-5, 40k resident.
	writeBudgetFixtureModels(t, home, "ensign-sd", "claude-fable-5[1m]",
		[]modelEntry{{"claude-fable-5", 40000}})

	m := runBudgetJSON(t, home, "ensign-sd")
	if _, ok := m["config_drift_warning"]; ok {
		t.Errorf("AC-1: config_drift_warning present on healthy suffix-dropped member:\n%v", m)
	}
	if got := m["context_limit"]; got != float64(1000000) {
		t.Errorf("AC-1: context_limit = %v, want 1000000", got)
	}
	if got := m["model"]; got != "claude-fable-5[1m]" {
		t.Errorf("AC-1: model = %v, want claude-fable-5[1m]", got)
	}
	// 40k of 1M = 4% <= 60% threshold.
	if got := m["reuse_ok"]; got != true {
		t.Errorf("AC-1: reuse_ok = %v, want true", got)
	}
}

// TestContextBudgetSyntheticExcludedNoMixed is AC-2: a healthy member whose jsonl
// carries a harness-injected <synthetic> entry alongside the real model emits no
// mixed_models_warning, picks the real model, and reads the real model's window —
// the placeholder is excluded from the census.
func TestContextBudgetSyntheticExcludedNoMixed(t *testing.T) {
	home := t.TempDir()
	// config claude-fable-5, runtime: one real entry + one <synthetic> entry.
	writeBudgetFixtureModels(t, home, "ensign-sy", "claude-fable-5",
		[]modelEntry{{"claude-fable-5", 40000}, {"<synthetic>", 1000}})

	m := runBudgetJSON(t, home, "ensign-sy")
	if _, ok := m["mixed_models_warning"]; ok {
		t.Errorf("AC-2: mixed_models_warning present despite <synthetic> exclusion:\n%v", m)
	}
	if got := m["model"]; got != "claude-fable-5" {
		t.Errorf("AC-2: model = %v, want claude-fable-5", got)
	}
	// claude-fable-5 is a 1M-context model under the 5-generation family rule.
	if got := m["context_limit"]; got != float64(1000000) {
		t.Errorf("AC-2: context_limit = %v, want 1000000", got)
	}
}

// TestContextBudgetGenuineDriftStillWarns is AC-3 (over-suppression guard): a
// member whose config and runtime are genuinely different models (distinct bases)
// still emits config_drift_warning — the suffix normalization must not silence it.
func TestContextBudgetGenuineDriftStillWarns(t *testing.T) {
	home := t.TempDir()
	// config claude-opus-4-8, runtime jsonl claude-sonnet-4-6 (different bases).
	writeBudgetFixtureModels(t, home, "ensign-gd", "claude-opus-4-8",
		[]modelEntry{{"claude-sonnet-4-6", 40000}})

	m := runBudgetJSON(t, home, "ensign-gd")
	if _, ok := m["config_drift_warning"]; !ok {
		t.Errorf("AC-3: config_drift_warning absent on genuine cross-family drift:\n%v", m)
	}
}

// TestContextBudgetGenuineMixedStillWarns is AC-4 (over-suppression guard): a
// member whose jsonl carries two genuinely different REAL models (no <synthetic>)
// still emits mixed_models_warning and picks the smallest window.
func TestContextBudgetGenuineMixedStillWarns(t *testing.T) {
	home := t.TempDir()
	// Two real models, no <synthetic>.
	writeBudgetFixtureModels(t, home, "ensign-gm", "claude-opus-4-8",
		[]modelEntry{{"claude-opus-4-8", 1000}, {"claude-sonnet-4-6", 2000}})

	m := runBudgetJSON(t, home, "ensign-gm")
	if _, ok := m["mixed_models_warning"]; !ok {
		t.Errorf("AC-4: mixed_models_warning absent on two genuinely-mixed real models:\n%v", m)
	}
	if got := m["context_limit"]; got != float64(200000) {
		t.Errorf("AC-4: context_limit = %v, want 200000 (smallest window)", got)
	}
}

// TestContextBudgetSonnet5ReuseFlip is the entity's AC-2 value-measuring case: a
// 250k-resident claude-sonnet-5 member reads usage_pct 125.0 / reuse_ok false
// under the pre-family-rule 200k default, and 25.0 / reuse_ok true once the
// 5-generation family rule grants the 1M window.
func TestContextBudgetSonnet5ReuseFlip(t *testing.T) {
	home := t.TempDir()
	writeBudgetFixture(t, home, "ensign-s5", "claude-sonnet-5", 250000)

	m := runBudgetJSON(t, home, "ensign-s5")
	if got := m["context_limit"]; got != float64(1000000) {
		t.Errorf("context_limit = %v, want 1000000", got)
	}
	if got := m["usage_pct"]; got != 25.0 {
		t.Errorf("usage_pct = %v, want 25.0", got)
	}
	if got := m["reuse_ok"]; got != true {
		t.Errorf("reuse_ok = %v, want true", got)
	}
}

// TestContextBudgetFable5ReuseFlip is the entity's AC-3: a claude-fable-5 member
// likewise resolves the 1M window. A bare fable-5 transcript observed live at
// 561,372 resident tokens would otherwise read usage_pct 280.7 / reuse_ok
// permanently false under the 200k default.
func TestContextBudgetFable5ReuseFlip(t *testing.T) {
	home := t.TempDir()
	writeBudgetFixture(t, home, "ensign-f5", "claude-fable-5", 250000)

	m := runBudgetJSON(t, home, "ensign-f5")
	if got := m["context_limit"]; got != float64(1000000) {
		t.Errorf("context_limit = %v, want 1000000", got)
	}
	if got := m["reuse_ok"]; got != true {
		t.Errorf("reuse_ok = %v, want true", got)
	}
}

// TestContextBudgetAliasConfigDriftAdvisory pins the known advisory consequence
// (entity AC-8): a member whose team config stamps the short alias "fable" while
// the runtime jsonl reports the full 5-family id "claude-fable-5" trips
// config_drift_warning (config resolves 200k, runtime 1M) even though reuse_ok
// is computed from the runtime model and unaffected by the drift.
func TestContextBudgetAliasConfigDriftAdvisory(t *testing.T) {
	home := t.TempDir()
	writeBudgetFixtureModels(t, home, "ensign-al", "fable",
		[]modelEntry{{"claude-fable-5", 40000}})

	m := runBudgetJSON(t, home, "ensign-al")
	if _, ok := m["config_drift_warning"]; !ok {
		t.Errorf("config_drift_warning absent on alias/5-family drift:\n%v", m)
	}
	if got := m["context_limit"]; got != float64(1000000) {
		t.Errorf("context_limit = %v, want 1000000 (runtime-resolved)", got)
	}
	if got := m["reuse_ok"]; got != true {
		t.Errorf("reuse_ok = %v, want true (unaffected by config drift)", got)
	}
}

// modelEntry is one assistant jsonl entry: a runtime model and its resident token
// sum (placed in input_tokens).
type modelEntry struct {
	model  string
	tokens int
}

// writeBudgetFixtureModels writes a ~/.claude tree like writeBudgetFixture but
// with a multi-line transcript: one assistant entry per modelEntry (the resident
// extractor sums input + cache, and scans backward for the first non-zero usage).
func writeBudgetFixtureModels(t *testing.T, home, name, configModel string, entries []modelEntry) {
	t.Helper()
	writeTeamConfigWithSession(t, home, "fixture-team", "sess", map[string]string{name: configModel})
	var lines []string
	for _, e := range entries {
		entry := map[string]any{
			"type": "assistant",
			"message": map[string]any{
				"model": e.model,
				"usage": map[string]any{
					"input_tokens":                e.tokens,
					"cache_creation_input_tokens": 0,
					"cache_read_input_tokens":     0,
				},
			},
		}
		line, _ := json.Marshal(entry)
		lines = append(lines, string(line))
	}
	subagents := filepath.Join(home, ".claude", "projects", "p", "sess", "subagents")
	writeTestFile(t, filepath.Join(subagents, "agent-"+name+".meta.json"),
		`{"agentType": "`+name+`"}`)
	writeTestFile(t, filepath.Join(subagents, "agent-"+name+".jsonl"),
		strings.Join(lines, "\n")+"\n")
}

// runBudgetJSON drives real ContextBudget over home for name, asserts exit 0, and
// returns the parsed stdout JSON for field/warning-key assertions.
func runBudgetJSON(t *testing.T, home, name string) map[string]any {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := ContextBudget(home, name, &stdout, &stderr); code != 0 {
		t.Fatalf("ContextBudget exit=%d stderr=%q", code, stderr.String())
	}
	var m map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &m); err != nil {
		t.Fatalf("ContextBudget stdout not JSON: %v\n%s", err, stdout.String())
	}
	return m
}

// writeBudgetFixture writes a minimal ~/.claude tree: a team config listing the
// member with model, and a one-line transcript whose resident equals tokens.
func writeBudgetFixture(t *testing.T, home, name, model string, tokens int) {
	t.Helper()
	writeTeamConfigWithSession(t, home, "fixture-team", "sess", map[string]string{name: model})
	entry := map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"model": model,
			"usage": map[string]any{
				"input_tokens":                tokens,
				"cache_creation_input_tokens": 0,
				"cache_read_input_tokens":     0,
			},
		},
	}
	line, _ := json.Marshal(entry)
	subagents := filepath.Join(home, ".claude", "projects", "p", "sess", "subagents")
	writeTestFile(t, filepath.Join(subagents, "agent-"+name+".meta.json"),
		`{"agentType": "`+name+`"}`)
	writeTestFile(t, filepath.Join(subagents, "agent-"+name+".jsonl"), string(line)+"\n")
}

// writeTeamConfigWithSession writes a team config.json with an optional
// leadSessionId and the given name->model members.
func writeTeamConfigWithSession(t *testing.T, home, team, session string, members map[string]string) {
	t.Helper()
	cfg := map[string]any{}
	if session != "" {
		cfg["leadSessionId"] = session
	}
	var ms []map[string]string
	for name, model := range members {
		ms = append(ms, map[string]string{"name": name, "model": model})
	}
	cfg["members"] = ms
	b, _ := json.MarshalIndent(cfg, "", "  ")
	writeTestFile(t, filepath.Join(home, ".claude", "teams", team, "config.json"), string(b))
}

// writeTestFile writes content to path, creating parent dirs.
func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
