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

// TestAssertShallowBootMeasuredOffline is the AC-6 de-risk: it validates the
// measured-saving oracle against committed real-shape streams BEFORE the live
// shallow-boot run spends a model on it. The positive fixture (a greet-and-stop
// boot, no TeamCreate, greet context under the ceiling) passes; the negative
// fixture (an eager-team boot with the ~89k cache_creation spike before the greet
// and a greet context over the ceiling) fails — proving the measurement
// distinguishes the realized saving from its absence, the AC-6 negative control.
func TestAssertShallowBootMeasuredOffline(t *testing.T) {
	if err := assertShallowBootMeasured(readMeasureFixture(t, "shallow-boot-greet.stream.jsonl")); err != nil {
		t.Fatalf("shallow-boot positive fixture must pass the measured-saving oracle: %v", err)
	}
	if err := assertShallowBootMeasured(readMeasureFixture(t, "eager-team-boot.stream.jsonl")); err == nil {
		t.Fatal("eager-team negative fixture must FAIL the measured-saving oracle (89k spike before greet + greet context over ceiling) — else the measurement does not distinguish the realized saving from its absence")
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

// TestShallowBootMeasureSignalsAreIndependent isolates the two AC-6 signals so
// neither can be silently dropped: a stream that fails ONLY the ceiling check (a
// heavy greet, no spike) and a stream that fails ONLY the spike check (a pre-greet
// 89k cache_creation, but a light greet) must each go red.
func TestShallowBootMeasureSignalsAreIndependent(t *testing.T) {
	// Only the ceiling fails: a single text greet turn whose context exceeds the
	// ceiling, with no cache_creation spike anywhere.
	heavyGreet := []journeymetrics.ClaudeTurn{
		{ID: "greet", Usage: journeymetrics.TokenTotals{Input: 100, CacheRead: greetContextCeiling, CacheCreation: 0}},
	}
	if err := assertShallowBootMeasuredTurns(heavyGreet); err == nil {
		t.Fatal("a greet turn whose context exceeds the ceiling (no spike) must fail on the ceiling check")
	}

	// Only the spike fails: a pre-greet dispatch turn carrying the ~89k spike, then
	// a light text greet under the ceiling.
	spikeThenLightGreet := []journeymetrics.ClaudeTurn{
		{ID: "team", Usage: journeymetrics.TokenTotals{Input: 8, CacheCreation: 89000, CacheRead: 16000}, ToolNames: []string{"TeamCreate"}},
		{ID: "greet", Usage: journeymetrics.TokenTotals{Input: 100, CacheRead: 5000, CacheCreation: 0}},
	}
	if err := assertShallowBootMeasuredTurns(spikeThenLightGreet); err == nil {
		t.Fatal("a pre-greet ~89k cache_creation spike (with a light greet) must fail on the spike check")
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
