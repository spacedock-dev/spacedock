package ensigncycle

import (
	"fmt"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// AC-6 boot-window measurement thresholds. The ceiling is the milestone's ~60k
// greet-turn context ceiling; the spike threshold is set below the ~89k team-mode
// prefix re-cache (TeamCreate re-caching the whole conversation prefix to the 1h
// cache) so any cache_creation on that order is caught, while the small
// per-turn cache_creation a healthy boot writes (a few thousand tokens) is not.
const (
	greetContextCeiling   = 60000
	teamRecacheSpikeFloor = 60000
)

// dispatchToolNames are the tool_use names that mark a worker dispatch / team
// creation — the boundary AC-6 uses to bound the pre-greet window. A turn that
// names one of these is NOT a greet turn.
var dispatchToolNames = map[string]bool{
	"Agent":      true,
	"Task":       true,
	"TeamCreate": true,
}

// greetTurnIndex returns the index of the greet turn in turns: the LAST assistant
// turn that emits no dispatch tool_use (the FO greets via text output and stops,
// it does not dispatch in the same turn). It returns -1 when every turn dispatches
// (no greet was produced). A shallow boot has no dispatch at all, so the greet turn
// is simply the final turn; an eager-team boot fires TeamCreate before the greet,
// so the greet turn is still the final non-dispatch turn and the TeamCreate turn
// stays inside the pre-greet window where the spike check can see it.
func greetTurnIndex(turns []journeymetrics.ClaudeTurn) int {
	idx := -1
	for i, t := range turns {
		dispatches := false
		for _, name := range t.ToolNames {
			if dispatchToolNames[name] {
				dispatches = true
				break
			}
		}
		if !dispatches {
			idx = i
		}
	}
	return idx
}

// teamToolNames are the Claude team-mode tool calls whose presence before the
// greet would mean an eager team was created. TeamCreate is the one the FO would
// fire; the others are listed so any team-lifecycle call pre-greet is caught.
var teamToolNames = map[string]bool{
	"TeamCreate": true,
	"TeamDelete": true,
}

// assertNoTeamCreateBeforeGreet is the AC-2 behavioral oracle over the captured
// stream's tool-call sequence: no team-mode tool call appears in the pre-greet
// window (the turns up to and including the greet turn). It is a behavioral
// observation over the real run's tool ordering, NOT a contract grep. A regression
// that re-introduced an eager team create surfaces a TeamCreate before the greet
// and fails this — the complement to AC-6's measured 89k-spike absence.
func assertNoTeamCreateBeforeGreet(stream string) error {
	turns, err := journeymetrics.ParseClaudeTurns([]byte(stream))
	if err != nil {
		return fmt.Errorf("parse stream for AC-2 team-call check: %w", err)
	}
	if len(turns) == 0 {
		return fmt.Errorf("stream carried no assistant turns — nothing to check")
	}
	greet := greetTurnIndex(turns)
	if greet < 0 {
		return fmt.Errorf("every assistant turn dispatched — no greet turn produced")
	}
	for i := 0; i <= greet; i++ {
		for _, name := range turns[i].ToolNames {
			if teamToolNames[name] {
				return fmt.Errorf("pre-greet turn %d emitted a %s tool call — a team was created before the greet (lazy-TeamCreate violated)", i, name)
			}
		}
	}
	return nil
}

// assertShallowBootMeasured is the AC-6 measured-saving oracle over a captured
// claude-stream.jsonl: it parses the stream per turn, identifies the greet turn and
// the pre-greet window (turns up to and including the greet turn), and asserts
//
//	(1) the greet-turn context (input + cache_read + cache_creation) is below the
//	    ~60k ceiling, and
//	(2) no pre-greet turn shows a cache_creation spike on the order of the ~89k
//	    team-mode prefix re-cache.
//
// It grades the host's emitted usage numbers — an independent source the contract
// cannot fake — never a prose match. A regression that re-introduced an eager team
// create or a heavy boot read pushes the greet context over the ceiling or surfaces
// the spike, failing this oracle.
func assertShallowBootMeasured(stream string) error {
	turns, err := journeymetrics.ParseClaudeTurns([]byte(stream))
	if err != nil {
		return fmt.Errorf("parse stream for boot-window measurement: %w", err)
	}
	return assertShallowBootMeasuredTurns(turns)
}

// assertShallowBootMeasuredTurns is the turn-level half of the AC-6 oracle, split
// out so the offline unit cases can drive the ceiling and spike checks directly
// without a stream fixture.
func assertShallowBootMeasuredTurns(turns []journeymetrics.ClaudeTurn) error {
	if len(turns) == 0 {
		return fmt.Errorf("stream carried no assistant turns — nothing to measure")
	}
	greet := greetTurnIndex(turns)
	if greet < 0 {
		return fmt.Errorf("every assistant turn dispatched — no greet turn produced")
	}
	if ctx := turns[greet].Context(); ctx >= greetContextCeiling {
		return fmt.Errorf("greet-turn context %d is not below the ~%dk ceiling — a heavy boot read or eager team create regressed the saving", ctx, greetContextCeiling/1000)
	}
	for i := 0; i <= greet; i++ {
		if cc := turns[i].Usage.CacheCreation; cc >= teamRecacheSpikeFloor {
			return fmt.Errorf("pre-greet turn %d shows a cache_creation spike of %d (>= %d) — the ~89k team-mode prefix re-cache fired before the greet", i, cc, teamRecacheSpikeFloor)
		}
	}
	return nil
}
