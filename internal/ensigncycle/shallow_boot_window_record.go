// ABOUTME: Boot-window telemetry extraction shared by the live shallow-boot
// ABOUTME: scenario test and the release-ledger backfill CLI (AC-1/AC-4).
package ensigncycle

import (
	"fmt"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// shallowBootWindowScenarioID is the journeymetrics ScenarioID the boot-window
// observation (see BuildShallowBootWindowRecord) publishes under. It MUST differ
// from the whole-run "shallow-boot" scenario ID the same run already publishes
// (see emitClaudeScenarioMetrics) — journeymetrics.recordFilename has no
// run-distinguishing component, so reusing "shallow-boot" here would silently
// overwrite that record instead of adding a sibling observation.
const shallowBootWindowScenarioID = "shallow-boot-window"

// dispatchToolNames are the tool_use names that mark a worker dispatch / team
// creation — the boundary the boot-window measurement uses to bound the
// pre-greet window. A turn that names one of these is NOT a greet turn.
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

// BuildShallowBootWindowRecord builds the shallow-boot-window journeymetrics.Record
// from a parsed turn list: Turns is greetIndex+1, and Tokens is the greet turn's
// FULL TokenTotals (input, output, cache_read, and cache_creation — not
// cache_creation alone), so a future reader can reconstruct both the former
// ceiling signal (Context() = input+cache_read+cache_creation) and the former
// pre-greet-spike signal (cache_creation) from this one recorded observation.
// Exported (unlike the live-scenario test glue around it) so the AC-4 release-
// ledger backfill CLI can apply the SAME extraction logic to archived streams,
// standalone, rather than duplicating it.
func BuildShallowBootWindowRecord(turns []journeymetrics.ClaudeTurn, model string) (journeymetrics.Record, error) {
	if len(turns) == 0 {
		return journeymetrics.Record{}, fmt.Errorf("stream carried no assistant turns — nothing to record")
	}
	greet := greetTurnIndex(turns)
	if greet < 0 {
		return journeymetrics.Record{}, fmt.Errorf("every assistant turn dispatched — no greet turn produced")
	}
	return journeymetrics.BuildRecord(journeymetrics.JourneySpec{
		ScenarioID: shallowBootWindowScenarioID,
		Source:     "live-harness",
		Mode:       journeymetrics.ModeLLMLive,
		Runtime:    "claude",
		Executor:   "llm",
		Host:       "claude",
		Model:      model,
	}, journeymetrics.BehaviorResult{Passed: true}, journeymetrics.Observation{
		Turns:  greet + 1,
		Tokens: turns[greet].Usage,
	}), nil
}
