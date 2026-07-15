package ensigncycle

import (
	"fmt"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

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

// assertShallowBootMeasured is the AC-6 structural oracle over a captured
// claude-stream.jsonl: it parses the stream per turn and asserts a greet turn was
// produced. It no longer gates on the ~60k greet-context ceiling or the ~89k
// pre-greet cache_creation spike — those were absolute constants calibrated once
// against a since-superseded baseline (task j9, v0.20.3) that carried no
// information about whether a given boot got better or worse. The greet turn's
// full token usage rides instead as a recorded (not gated) shallow-boot-window
// journeymetrics observation — see BuildShallowBootWindowRecord.
func assertShallowBootMeasured(stream string) error {
	turns, err := journeymetrics.ParseClaudeTurns([]byte(stream))
	if err != nil {
		return fmt.Errorf("parse stream for boot-window measurement: %w", err)
	}
	return assertShallowBootMeasuredTurns(turns)
}

// assertShallowBootMeasuredTurns is the turn-level half of the boot oracle, split
// out so the offline unit cases can drive it directly without a stream fixture.
// Only the structural checks remain a hard failure — the two former threshold
// branches (greet-context ceiling, pre-greet cache_creation spike) were removed
// here, at their actual location, not merely skipped by the caller.
func assertShallowBootMeasuredTurns(turns []journeymetrics.ClaudeTurn) error {
	if len(turns) == 0 {
		return fmt.Errorf("stream carried no assistant turns — nothing to measure")
	}
	if greetTurnIndex(turns) < 0 {
		return fmt.Errorf("every assistant turn dispatched — no greet turn produced")
	}
	return nil
}
