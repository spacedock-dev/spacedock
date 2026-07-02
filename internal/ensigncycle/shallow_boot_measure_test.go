package ensigncycle

import (
	"fmt"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// shallowBootWindowScenarioID is the journeymetrics ScenarioID the boot-window
// observation (see buildShallowBootWindowRecord) publishes under. It MUST differ
// from the whole-run "shallow-boot" scenario ID the same run already publishes
// (see emitClaudeScenarioMetrics) — journeymetrics.recordFilename has no
// run-distinguishing component, so reusing "shallow-boot" here would silently
// overwrite that record instead of adding a sibling observation.
const shallowBootWindowScenarioID = "shallow-boot-window"

// deferredFOSkillNames are the first-officer-internal skills a greet-and-stop boot
// must NOT invoke before the greet: the status-viewer surface, the write/id-style
// surface, and the dispatch-failure-recovery surface. Each loads only at its trigger
// (the FIRST status query / --set / id lookup / issue filing, the FIRST write to
// main, or a dispatch failure — never at boot). present-gate is deliberately NOT
// here — the greet legitimately presents a ready gate via
// Skill(skill="spacedock:present-gate") (Startup step 8), so the oracle keys on the
// skill ARGUMENT, not on any Skill call.
var deferredFOSkillNames = []string{"fo-status-viewer", "fo-write-core", "fo-dispatch-recovery"}

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

// assertGreetInvokesNoDeferredFOSkill is the AC-2 behavioral oracle over the captured
// stream's tool-call sequence: no pre-greet turn invokes a deferred FO skill (a Skill
// tool_use whose skill argument names fo-status-viewer or fo-write-core, in the turns
// up to and including the greet turn). It is a behavioral observation over the real
// run's tool ordering, NOT a contract grep — the sibling of
// assertNoTeamCreateBeforeGreet. The greet must compose its summary from status --boot
// + README frontmatter; a regression that eagerly loaded the deferred status-viewer or
// write skill surfaces a pre-greet Skill(spacedock:fo-status-viewer|fo-write-core) and
// fails this. It keys on the skill ARGUMENT, not on any Skill call — the greet
// legitimately invokes Skill(skill="spacedock:present-gate") to present a ready gate,
// so a blanket "no pre-greet Skill" oracle would false-fail that.
func assertGreetInvokesNoDeferredFOSkill(stream string) error {
	turns, err := journeymetrics.ParseClaudeTurns([]byte(stream))
	if err != nil {
		return fmt.Errorf("parse stream for AC-2 deferred-skill check: %w", err)
	}
	if len(turns) == 0 {
		return fmt.Errorf("stream carried no assistant turns — nothing to check")
	}
	greet := greetTurnIndex(turns)
	if greet < 0 {
		return fmt.Errorf("every assistant turn dispatched — no greet turn produced")
	}
	for i := 0; i <= greet; i++ {
		for _, skill := range turns[i].SkillNames {
			for _, forbidden := range deferredFOSkillNames {
				if strings.Contains(skill, forbidden) {
					return fmt.Errorf("pre-greet turn %d invoked the deferred FO skill %q via Skill(skill=%q) — the greet must render from status --boot, not load a deferred FO skill (present-gate is allowed pre-greet)", i, forbidden, skill)
				}
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
// journeymetrics observation — see buildShallowBootWindowRecord.
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

// buildShallowBootWindowRecord builds the shallow-boot-window journeymetrics.Record
// from a parsed turn list: Turns is greetIndex+1, and Tokens is the greet turn's
// FULL TokenTotals (input, output, cache_read, and cache_creation — not
// cache_creation alone), so a future reader can reconstruct both the former
// ceiling signal (Context() = input+cache_read+cache_creation) and the former
// pre-greet-spike signal (cache_creation) from this one recorded observation.
func buildShallowBootWindowRecord(turns []journeymetrics.ClaudeTurn, model string) (journeymetrics.Record, error) {
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
