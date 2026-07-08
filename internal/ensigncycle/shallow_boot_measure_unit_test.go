package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// readMeasureFixture reads a committed claude-stream.jsonl fixture from testdata.
// These are trimmed copies of real stream-json (the journeymetrics precedent — a
// committed testdata/ jsonl, never a ~/.claude/projects per-machine artifact), so
// the AC-6 oracle has an offline positive and negative branch to validate against
// before the live shallow-boot run relies on it.
func readMeasureFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestAssertShallowBootMeasuredOffline validates that assertShallowBootMeasured's
// only remaining checks are structural (the stream produced a greet turn) now that
// the former ~60k ceiling/spike thresholds are removed at their actual location
// (assertShallowBootMeasuredTurns), not merely skipped by the caller. Both the
// shallow-boot positive fixture and the eager-team-boot fixture — whose ~89k
// pre-greet spike and over-ceiling greet used to FAIL this oracle — now pass,
// since neither fixture violates the remaining structural checks. The
// ceiling/spike signal itself now rides the shallow-boot-window journeymetrics
// observation instead of a pass/fail gate — see TestBuildShallowBootWindowRecord.
func TestAssertShallowBootMeasuredOffline(t *testing.T) {
	if err := assertShallowBootMeasured(readMeasureFixture(t, "shallow-boot-greet.stream.jsonl")); err != nil {
		t.Fatalf("shallow-boot positive fixture must pass the structural boot oracle: %v", err)
	}
	if err := assertShallowBootMeasured(readMeasureFixture(t, "eager-team-boot.stream.jsonl")); err != nil {
		t.Fatalf("eager-team fixture must pass now that the ceiling/spike thresholds are removed (it still produces a valid greet turn): %v", err)
	}
}

// TestAssertNoTeamCreateBeforeGreetOffline validates the AC-2 behavioral oracle
// against the committed MULTI-DELTA streams: the shallow-boot positive (no
// TeamCreate at all) passes; the eager-team negative (a TeamCreate before the greet)
// fails. Both fixtures are multi-delta — each message carries a `thinking` delta then
// a tool_use/text delta, the real runner shape — so the negative control genuinely
// exercises the path where the TeamCreate lands on a LATER delta. A first-delta-only
// parse would have FALSE-PASSED this negative (the hollow-AC-2 defect the forensics
// caught), so this is the positive control that the lazy-TeamCreate proof is real.
func TestAssertNoTeamCreateBeforeGreetOffline(t *testing.T) {
	if err := assertNoTeamCreateBeforeGreet(readMeasureFixture(t, "shallow-boot-greet.stream.jsonl")); err != nil {
		t.Fatalf("shallow-boot positive fixture (no TeamCreate) must pass AC-2: %v", err)
	}
	if err := assertNoTeamCreateBeforeGreet(readMeasureFixture(t, "eager-team-boot.stream.jsonl")); err == nil {
		t.Fatal("eager-team negative fixture (multi-delta TeamCreate before greet) must FAIL AC-2 — the TeamCreate lands on a later delta, so a first-delta-only parse would false-pass")
	}
}

// TestAssertNoTeamCreateBeforeGreetCatchesLaterDeltaTeamCreate is the AC-2
// positive control over the parser fix: a stream whose TeamCreate is on a LATER
// delta of its message (the real runner shape, NOT the synthetic single-delta one)
// must make assertNoTeamCreateBeforeGreet RED. Before the multi-delta merge this
// false-passed — the proof of no-team-at-boot was hollow.
func TestAssertNoTeamCreateBeforeGreetCatchesLaterDeltaTeamCreate(t *testing.T) {
	// msg_team: thinking on delta[0], TeamCreate on delta[1]; then a text greet.
	stream := `{"type":"assistant","message":{"id":"msg_team","model":"claude-opus-4-8","usage":{"input_tokens":8,"cache_read_input_tokens":16000,"cache_creation_input_tokens":89000},"content":[{"type":"thinking","thinking":"create the team"}]}}
{"type":"assistant","message":{"id":"msg_team","model":"claude-opus-4-8","usage":{"input_tokens":8,"cache_read_input_tokens":16000,"cache_creation_input_tokens":89000},"content":[{"type":"tool_use","id":"toolu_tc","name":"TeamCreate","input":{"team_name":"eager"}}]}}
{"type":"assistant","message":{"id":"msg_greet","model":"claude-opus-4-8","usage":{"input_tokens":100,"cache_read_input_tokens":5000,"cache_creation_input_tokens":0},"content":[{"type":"text","text":"Gate review: ... Decision: approve or reject?"}]}}`
	if err := assertNoTeamCreateBeforeGreet(stream); err == nil {
		t.Fatal("a pre-greet TeamCreate on a LATER delta must make assertNoTeamCreateBeforeGreet RED — the parser must merge later-delta tool_use, not read only the first delta")
	}
}

// TestAssertGreetInvokesNoDeferredFOSkillOffline validates the AC-2 deferred-skill
// oracle against committed MULTI-DELTA streams: the shallow-boot positive (the greet
// renders from status --boot and invokes present-gate to show the ready gate, but no
// deferred FO skill) passes — proving the oracle keys on the skill ARGUMENT and does
// NOT flag the legitimate pre-greet Skill(spacedock:present-gate); the negative (a
// pre-greet Skill(spacedock:fo-status-viewer), on a LATER delta of its message) fails.
// Multi-delta is mandatory — the Skill lands on a later delta, so a first-delta-only
// parse would FALSE-PASS (the documented hollow-AC defect).
func TestAssertGreetInvokesNoDeferredFOSkillOffline(t *testing.T) {
	if err := assertGreetInvokesNoDeferredFOSkill(readMeasureFixture(t, "shallow-boot-greet.stream.jsonl")); err != nil {
		t.Fatalf("shallow-boot positive fixture (renders from status --boot, invokes only present-gate pre-greet) must pass AC-2: %v", err)
	}
	if err := assertGreetInvokesNoDeferredFOSkill(readMeasureFixture(t, "greet-invokes-fo-status-viewer-skill.stream.jsonl")); err == nil {
		t.Fatal("the negative fixture (multi-delta pre-greet Skill(spacedock:fo-status-viewer)) must FAIL AC-2 — the Skill lands on a later delta, so a first-delta-only parse would false-pass")
	}
}

// TestAssertGreetInvokesNoDeferredFOSkillCatchesLaterDeltaInvoke is the AC-2 positive
// control over the parser extension: a stream whose Skill(spacedock:fo-status-viewer)
// is on a LATER delta of its message (the real runner shape, NOT a synthetic
// single-delta one) must make assertGreetInvokesNoDeferredFOSkill RED. A
// first-delta-only parse drops the later-delta Skill entirely — the skill argument
// would be invisible and the greet-independence proof hollow. This mirrors the retired
// Read-based later-delta control; a Skill block can land on a later delta exactly like
// a Read, so retiring the Read control without a Skill equivalent would be a coverage
// regression on the design's riskiest path.
func TestAssertGreetInvokesNoDeferredFOSkillCatchesLaterDeltaInvoke(t *testing.T) {
	// msg_skill: thinking on delta[0], Skill(spacedock:fo-status-viewer) on delta[1]; then a text greet.
	stream := `{"type":"assistant","message":{"id":"msg_skill","model":"claude-opus-4-8","usage":{"input_tokens":8,"cache_read_input_tokens":21000,"cache_creation_input_tokens":600},"content":[{"type":"thinking","thinking":"load the status viewer skill"}]}}
{"type":"assistant","message":{"id":"msg_skill","model":"claude-opus-4-8","usage":{"input_tokens":8,"cache_read_input_tokens":21000,"cache_creation_input_tokens":600},"content":[{"type":"tool_use","id":"toolu_skill","name":"Skill","input":{"skill":"spacedock:fo-status-viewer"}}]}}
{"type":"assistant","message":{"id":"msg_greet","model":"claude-opus-4-8","usage":{"input_tokens":100,"cache_read_input_tokens":5000,"cache_creation_input_tokens":0},"content":[{"type":"text","text":"Workflow overview: ... Gate review: ... Decision: approve or reject?"}]}}`
	if err := assertGreetInvokesNoDeferredFOSkill(stream); err == nil {
		t.Fatal("a pre-greet Skill(spacedock:fo-status-viewer) on a LATER delta must make assertGreetInvokesNoDeferredFOSkill RED — the parser must merge later-delta skill names, not read only the first delta")
	}
}

// TestParserExtractsTeamCallsFromRealHangCapture is the validator-named ready
// oracle: the committed real-runner stream `sonnet_teamdelete_hang.stream.jsonl`
// (20/27 message ids multi-delta; its lone TeamCreate and TeamDelete each land on a
// non-first delta) must surface BOTH team calls through the fixed ParseClaudeTurns.
// Against the pre-fix first-delta-only parse this reported TeamCreate=false across
// all 27 turns (the proven false-pass); the merge reports TeamCreate=true. Driving
// the FULL committed fixture (not a trimmed copy) pins the fix to the exact stream
// the forensics verified the defect on.
func TestParserExtractsTeamCallsFromRealHangCapture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "sonnet_teamdelete_hang.stream.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	turns, err := journeymetrics.ParseClaudeTurns(data)
	if err != nil {
		t.Fatalf("ParseClaudeTurns: %v", err)
	}
	saw := map[string]bool{}
	for _, turn := range turns {
		for _, name := range turn.ToolNames {
			saw[name] = true
		}
	}
	if !saw["TeamCreate"] {
		t.Error("the real hang capture's TeamCreate (on a non-first delta) was not surfaced — the parser is still first-delta-only (the AC-2 false-pass)")
	}
	if !saw["TeamDelete"] {
		t.Error("the real hang capture's TeamDelete (on a non-first delta) was not surfaced")
	}
}

// TestBuildShallowBootWindowRecord is AC-1's primary offline fixture unit test: it
// proves EmitRecord writes a shallow-boot-window--claude--llm--llm-live--<model>--
// measured.json file DISTINCT from the pre-existing whole-run
// shallow-boot--...--measured.json record (both present after one scenario run,
// neither overwriting the other), carrying Turns == greetIndex+1 and Tokens equal
// to the greet turn's full actual TokenTotals (input, output, cache_read,
// cache_creation) from the fixture stream — not CacheCreation alone. It also
// proves the follow-on fields the fixture now carries: BaselineTokens equals
// the FIRST turn's full TokenTotals (so a reader can compute net boot cost
// without re-parsing the stream), ClaudeCodeVersion is threaded through from
// the fixture's system/init event's claude_code_version field, and ResolvedModel
// (the fixture's actual "claude-opus-4-8") differs from the caller-supplied
// Model alias ("claude-test-model") — proving the two fields serve distinct
// purposes rather than duplicating each other.
func TestBuildShallowBootWindowRecord(t *testing.T) {
	stream := readMeasureFixture(t, "shallow-boot-greet.stream.jsonl")
	turns, err := journeymetrics.ParseClaudeTurns([]byte(stream))
	if err != nil {
		t.Fatalf("ParseClaudeTurns: %v", err)
	}
	greet := greetTurnIndex(turns)
	if greet < 0 {
		t.Fatal("fixture must produce a greet turn")
	}
	claudeCodeVersion := journeymetrics.ParseClaudeCodeVersion([]byte(stream))
	if claudeCodeVersion == "" {
		t.Fatal("fixture must carry a claude_code_version on its system/init event")
	}
	resolvedModel := journeymetrics.ParseClaudeInitModel([]byte(stream))
	if resolvedModel == "" {
		t.Fatal("fixture must carry a model on its system/init event")
	}

	const model = "claude-test-model"
	record, err := BuildShallowBootWindowRecord(turns, model, claudeCodeVersion, resolvedModel)
	if err != nil {
		t.Fatalf("BuildShallowBootWindowRecord: %v", err)
	}
	if record.ScenarioID != "shallow-boot-window" {
		t.Fatalf("ScenarioID = %q, want shallow-boot-window", record.ScenarioID)
	}
	if record.Turns != greet+1 {
		t.Fatalf("Turns = %d, want %d (greetIndex+1)", record.Turns, greet+1)
	}
	want := turns[greet].Usage
	if record.Tokens.Input != want.Input || record.Tokens.Output != want.Output ||
		record.Tokens.CacheCreation != want.CacheCreation || record.Tokens.CacheRead != want.CacheRead {
		t.Fatalf("Tokens = %+v, want the greet turn's full TokenTotals %+v (not CacheCreation alone)", record.Tokens, want)
	}
	wantBaseline := turns[0].Usage
	if record.BaselineTokens.Input != wantBaseline.Input || record.BaselineTokens.Output != wantBaseline.Output ||
		record.BaselineTokens.CacheCreation != wantBaseline.CacheCreation || record.BaselineTokens.CacheRead != wantBaseline.CacheRead {
		t.Fatalf("BaselineTokens = %+v, want the FIRST turn's full TokenTotals %+v", record.BaselineTokens, wantBaseline)
	}
	if record.ClaudeCodeVersion != claudeCodeVersion {
		t.Fatalf("ClaudeCodeVersion = %q, want %q from the fixture's system/init event", record.ClaudeCodeVersion, claudeCodeVersion)
	}
	if record.Model != model {
		t.Fatalf("Model = %q, want the caller-supplied alias %q unchanged", record.Model, model)
	}
	if record.ResolvedModel != resolvedModel || record.ResolvedModel == record.Model {
		t.Fatalf("ResolvedModel = %q, want the fixture's actual model %q, distinct from Model %q", record.ResolvedModel, resolvedModel, record.Model)
	}

	// The whole-run "shallow-boot" record the same scenario run already publishes
	// (see emitClaudeScenarioMetrics) must survive untouched as a sibling file.
	dir := t.TempDir()
	wholeRun := journeymetrics.BuildRecord(journeymetrics.JourneySpec{
		ScenarioID: "shallow-boot",
		Source:     "live-harness",
		Mode:       journeymetrics.ModeLLMLive,
		Runtime:    "claude",
		Executor:   "llm",
		Host:       "claude",
		Model:      model,
	}, journeymetrics.BehaviorResult{Passed: true}, journeymetrics.Observation{})
	if err := journeymetrics.EmitRecord(dir, wholeRun); err != nil {
		t.Fatalf("emit whole-run shallow-boot record: %v", err)
	}
	if err := journeymetrics.EmitRecord(dir, record); err != nil {
		t.Fatalf("emit shallow-boot-window record: %v", err)
	}

	windowPath := filepath.Join(dir, "shallow-boot-window--claude--llm--llm-live--"+model+"--measured.json")
	if _, err := os.Stat(windowPath); err != nil {
		t.Fatalf("expected shallow-boot-window record file at %s: %v", windowPath, err)
	}
	wholeRunPath := filepath.Join(dir, "shallow-boot--claude--llm--llm-live--"+model+"--measured.json")
	if _, err := os.Stat(wholeRunPath); err != nil {
		t.Fatalf("the pre-existing shallow-boot record must survive unmodified at %s: %v", wholeRunPath, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected exactly 2 record files (shallow-boot and shallow-boot-window), got %d", len(entries))
	}
}

// TestShallowBootMeasureSignalsAreIndependent previously proved the ceiling and
// spike checks could each fail independently. Now that both threshold branches are
// removed from assertShallowBootMeasuredTurns (not merely bypassed by the caller),
// it proves the gate is gone at its actual location — none of the three cases that
// used to trip a threshold return an error — while BuildShallowBootWindowRecord
// still captures each case's full greet-turn TokenTotals, so the former
// ceiling/spike signal stays reconstructable from the recorded telemetry even
// though it no longer gates CI.
func TestShallowBootMeasureSignalsAreIndependent(t *testing.T) {
	// Formerly failed only the ceiling check: a single text greet turn whose
	// context exceeds the old ~60k ceiling, with no cache_creation spike anywhere.
	heavyGreet := []journeymetrics.ClaudeTurn{
		{ID: "greet", Usage: journeymetrics.TokenTotals{Input: 100, CacheRead: 60000, CacheCreation: 0}},
	}
	if err := assertShallowBootMeasuredTurns(heavyGreet); err != nil {
		t.Fatalf("a greet turn over the former ceiling must no longer fail now that the threshold gate is removed: %v", err)
	}
	if record, err := BuildShallowBootWindowRecord(heavyGreet, "test-model", "2.1.196", "claude-sonnet-4-6"); err != nil {
		t.Fatalf("BuildShallowBootWindowRecord(heavyGreet): %v", err)
	} else if record.Tokens.CacheRead != 60000 {
		t.Fatalf("recorded Tokens must preserve the full greet-turn usage that used to trip the ceiling check, got %+v", record.Tokens)
	} else if record.ClaudeCodeVersion != "2.1.196" {
		t.Fatalf("ClaudeCodeVersion = %q, want the passed-through version 2.1.196", record.ClaudeCodeVersion)
	} else if record.ResolvedModel != "claude-sonnet-4-6" {
		t.Fatalf("ResolvedModel = %q, want the passed-through resolved model claude-sonnet-4-6", record.ResolvedModel)
	}

	// Formerly failed only the spike check: a pre-greet dispatch turn carrying the
	// ~89k spike, then a light text greet under the old ceiling.
	spikeThenLightGreet := []journeymetrics.ClaudeTurn{
		{ID: "team", Usage: journeymetrics.TokenTotals{Input: 8, CacheCreation: 89000, CacheRead: 16000}, ToolNames: []string{"TeamCreate"}},
		{ID: "greet", Usage: journeymetrics.TokenTotals{Input: 100, CacheRead: 5000, CacheCreation: 0}},
	}
	if err := assertShallowBootMeasuredTurns(spikeThenLightGreet); err != nil {
		t.Fatalf("a pre-greet ~89k cache_creation spike (light greet) must no longer fail now that the threshold gate is removed: %v", err)
	}
	if record, err := BuildShallowBootWindowRecord(spikeThenLightGreet, "test-model", "", ""); err != nil {
		t.Fatalf("BuildShallowBootWindowRecord(spikeThenLightGreet): %v", err)
	} else if record.Turns != 2 {
		t.Fatalf("Turns = %d, want 2 (greetIndex+1, greet is turns[1])", record.Turns)
	} else if record.BaselineTokens.CacheCreation != 89000 {
		t.Fatalf("BaselineTokens must preserve the pre-greet ~89k spike turn's full usage (turns[0]), got %+v", record.BaselineTokens)
	} else if record.PreGreetPeakCacheCreation != 89000 {
		t.Fatalf("PreGreetPeakCacheCreation = %d, want 89000 (the former teamRecacheSpikeFloor signal)", record.PreGreetPeakCacheCreation)
	}

	// Both clean: a light greet, no pre-greet spike — the realized-saving end-state.
	clean := []journeymetrics.ClaudeTurn{
		{ID: "boot", Usage: journeymetrics.TokenTotals{Input: 6, CacheCreation: 900, CacheRead: 18000}, ToolNames: []string{"Bash"}},
		{ID: "greet", Usage: journeymetrics.TokenTotals{Input: 120, CacheRead: 42000, CacheCreation: 400}},
	}
	if err := assertShallowBootMeasuredTurns(clean); err != nil {
		t.Fatalf("a clean shallow boot (light greet, no spike) must pass: %v", err)
	}
}

// TestPreGreetPeakCacheCreationFindsSpikeNotOnFirstTurn is the regression proof
// for the captain-review finding that BaselineTokens (turns[0].Usage alone)
// does NOT reconstruct the former teamRecacheSpikeFloor signal for a REALISTIC
// multi-turn pre-greet window where the spike lands on a LATER pre-greet turn,
// not the first one. The original teamRecacheSpikeFloor check looped EVERY
// pre-greet turn (assertShallowBootMeasuredTurns's removed `for i := 0; i <=
// greet; i++` branch) — a max over the whole window, not just turns[0].
func TestPreGreetPeakCacheCreationFindsSpikeNotOnFirstTurn(t *testing.T) {
	turns := []journeymetrics.ClaudeTurn{
		{ID: "boot", Usage: journeymetrics.TokenTotals{Input: 4, CacheCreation: 500, CacheRead: 8000}, ToolNames: []string{"Bash"}},
		{ID: "team", Usage: journeymetrics.TokenTotals{Input: 8, CacheCreation: 89000, CacheRead: 16000}, ToolNames: []string{"TeamCreate"}},
		{ID: "greet", Usage: journeymetrics.TokenTotals{Input: 100, CacheRead: 5000, CacheCreation: 0}},
	}
	record, err := BuildShallowBootWindowRecord(turns, "test-model", "", "")
	if err != nil {
		t.Fatalf("BuildShallowBootWindowRecord: %v", err)
	}
	if record.BaselineTokens.CacheCreation != 500 {
		t.Fatalf("BaselineTokens (turns[0]) = %+v, want CacheCreation 500 — it is NOT the spike turn here, proving it alone cannot reconstruct the spike signal", record.BaselineTokens)
	}
	if record.PreGreetPeakCacheCreation != 89000 {
		t.Fatalf("PreGreetPeakCacheCreation = %d, want 89000 — the spike on turns[1] (NOT turns[0]) must still be found", record.PreGreetPeakCacheCreation)
	}
}
