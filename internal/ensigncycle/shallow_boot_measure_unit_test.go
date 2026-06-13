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
