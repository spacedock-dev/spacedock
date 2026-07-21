// ABOUTME: Offline oracles over a Claude stream-json transcript for the two
// ABOUTME: fo-dispatch-recovery live scenarios (AC-2 degraded-bare, AC-3 break-glass).
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"strings"
)

// retiredDegradedModeReportPrefix is the fixed lead of the verbatim captain report
// the RETIRED `## Degraded Mode` section used to mandate on trip. It survives here
// only as a wrong-way check: after the retirement, no contract prose mandates it, so
// a bare-dispatch drive must NOT emit it. assertBareReachableObservables fails if it
// reappears.
const retiredDegradedModeReportPrefix = "Falling back to bare mode for the remainder of this session due to infrastructure failure."

// retrySuffix is the distinct axis a bounded dispatch-failure re-attempt carries on
// the `{worker_key}-{slug}-{stage}` name — never a `-cycleN` increment, so a retry
// cannot advance the feedback counter.
const retrySuffix = "-retry"

// recoverySkillArg is the Skill tool_use argument the break-glass scenario must
// observe after its trigger — the resident trigger line's `Skill(skill="spacedock:fo-dispatch-recovery")` load.
const recoverySkillArg = "spacedock:fo-dispatch-recovery"

// streamContentBlock is one tool_use/text content block of one assistant message
// delta, decoded generically: Input as a raw key-set (not a typed struct) so an
// oracle can test field PRESENCE (e.g. "was run_in_background emitted at all")
// rather than only its value — a bare-mode Agent() omits the key entirely, it does
// not emit `run_in_background: false`.
type streamContentBlock struct {
	Type  string                     `json:"type"`
	Name  string                     `json:"name"`
	Text  string                     `json:"text"`
	Input map[string]json.RawMessage `json:"input"`
}

type streamAssistantRow struct {
	Message *struct {
		Content []streamContentBlock `json:"content"`
	} `json:"message"`
}

// walkStreamBlocks decodes every assistant-row content block of stream in stream
// order and calls fn on each. Non-JSON lines (folded stderr) and non-assistant rows
// are skipped, matching journeymetrics.ParseClaudeTurns' tolerance.
func walkStreamBlocks(stream string, fn func(block streamContentBlock)) {
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var row streamAssistantRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			continue
		}
		if row.Message == nil {
			continue
		}
		for _, block := range row.Message.Content {
			fn(block)
		}
	}
}

// inputHasKey reports whether the tool_use input carried the named key AT ALL,
// distinct from carrying it with a false/empty value.
func inputHasKey(input map[string]json.RawMessage, key string) bool {
	_, ok := input[key]
	return ok
}

func inputStringField(input map[string]json.RawMessage, key string) string {
	raw, ok := input[key]
	if !ok {
		return ""
	}
	var s string
	_ = json.Unmarshal(raw, &s)
	return s
}

// assertBoundedRetryObservables is the AC-1 offline oracle: over the captured stream
// it asserts a dispatch failure of one `(entity, stage)` is retried exactly ONCE and
// no further. The first Agent() dispatch is the failed attempt; the bounded
// re-attempt is a fresh Agent() dispatch carrying the distinct `-retry` suffix on the
// same `{worker_key}-{slug}-{stage}` stem. No third Agent() call may appear for that
// stem — the bound. It reads only the dispatch surface (Agent() names), introducing
// no error-string classification.
func assertBoundedRetryObservables(stream string) error {
	var agentNames []string
	walkStreamBlocks(stream, func(block streamContentBlock) {
		if block.Type == "tool_use" && block.Name == "Agent" {
			agentNames = append(agentNames, inputStringField(block.Input, "name"))
		}
	})
	if len(agentNames) == 0 {
		return fmt.Errorf("no Agent() dispatch observed in the stream — the bounded-retry path never dispatched")
	}
	base := agentNames[0]
	if base == "" {
		return fmt.Errorf("the first Agent() dispatch carried no name — the retry stem cannot be derived")
	}
	retryName := base + retrySuffix
	retries := 0
	for _, name := range agentNames[1:] {
		if name == retryName {
			retries++
		}
	}
	if retries == 0 {
		return fmt.Errorf("no re-dispatch carrying the %q suffix (%q) after the first Agent() dispatch %q — a dispatch failure must be retried once", retrySuffix, retryName, base)
	}
	if len(agentNames) > 2 {
		return fmt.Errorf("%d Agent() calls for one (entity, stage) — the retry is bounded to the initial dispatch + one %q re-attempt; no third attempt may appear: observed %v", len(agentNames), retrySuffix, agentNames)
	}
	return nil
}

// assertBareReachableObservables is the AC-2 behavioral oracle, rewritten for the
// post-retirement expectation: over the captured stream it asserts (i) at least one
// bare-shaped Agent() call (neither `name` nor `run_in_background`) — bare dispatch is
// reached — AND, as the wrong-way check the retirement preserves, (ii) NO retired
// Degraded Mode captain report in any text block, and (iii) NO
// Skill(skill="spacedock:fo-dispatch-recovery") load. A drive that still emits the
// report or loads the recovery skill is now a FAILURE.
func assertBareReachableObservables(stream string) error {
	sawBareAgent := false
	sawRetiredReport := false
	sawRecoverySkill := false

	walkStreamBlocks(stream, func(block streamContentBlock) {
		switch block.Type {
		case "text":
			if strings.Contains(block.Text, retiredDegradedModeReportPrefix) {
				sawRetiredReport = true
			}
		case "tool_use":
			switch block.Name {
			case "Skill":
				if inputStringField(block.Input, "skill") == recoverySkillArg {
					sawRecoverySkill = true
				}
			case "Agent":
				if !inputHasKey(block.Input, "name") && !inputHasKey(block.Input, "run_in_background") {
					sawBareAgent = true
				}
			}
		}
	})

	if !sawBareAgent {
		return fmt.Errorf("no bare-shaped Agent() call (neither `name` nor `run_in_background`) observed — bare dispatch was not reached")
	}
	if sawRetiredReport {
		return fmt.Errorf("the retired Degraded Mode captain report (%q) still appeared — it was retired with Degraded Mode and must not fire on a bare drive", retiredDegradedModeReportPrefix)
	}
	if sawRecoverySkill {
		return fmt.Errorf("a Skill(skill=%q) load appeared — a post-retirement bare drive must NOT load the recovery skill", recoverySkillArg)
	}
	return nil
}

// assertBreakGlassObservables is the AC-3 behavioral oracle: over the captured
// stream it asserts (i) a captain-facing helper-failure report (a text block
// mentioning the failed helper command) appears BEFORE any Agent() tool_use, (ii) a
// Skill(skill="spacedock:fo-dispatch-recovery") tool_use appears, and (iii) at least
// one Agent() tool_use carries run_in_background=true, a name matching the
// `{worker_key}-{slug}-{stage}` shape, and a prompt containing both
// `Skill(skill="spacedock:ensign")` and `### Stage definition`.
func assertBreakGlassObservables(stream string) error {
	sawReportBeforeAgent := false
	sawAgentCall := false
	sawRecoverySkill := false
	var breakGlassAgentSeen bool
	var breakGlassAgentDetails []string

	walkStreamBlocks(stream, func(block streamContentBlock) {
		switch block.Type {
		case "text":
			if !sawAgentCall && strings.Contains(block.Text, "dispatch build") {
				sawReportBeforeAgent = true
			}
		case "tool_use":
			switch block.Name {
			case "Skill":
				if inputStringField(block.Input, "skill") == recoverySkillArg {
					sawRecoverySkill = true
				}
			case "Agent":
				sawAgentCall = true
				var runInBackground bool
				if raw, ok := block.Input["run_in_background"]; ok {
					_ = json.Unmarshal(raw, &runInBackground)
				}
				name := inputStringField(block.Input, "name")
				prompt := inputStringField(block.Input, "prompt")
				hasEnsignSkill := strings.Contains(prompt, `Skill(skill="spacedock:ensign")`)
				hasStageDef := strings.Contains(prompt, "### Stage definition")
				nameShaped := strings.Count(name, "-") >= 2
				details := fmt.Sprintf("name=%q run_in_background=%v ensign-skill-in-prompt=%v stage-def-in-prompt=%v", name, runInBackground, hasEnsignSkill, hasStageDef)
				breakGlassAgentDetails = append(breakGlassAgentDetails, details)
				if runInBackground && nameShaped && hasEnsignSkill && hasStageDef {
					breakGlassAgentSeen = true
				}
			}
		}
	})

	if !sawReportBeforeAgent {
		return fmt.Errorf("no captain-facing helper-failure report (a text block naming `dispatch build`) observed before the first Agent() call")
	}
	if !sawRecoverySkill {
		return fmt.Errorf("no Skill(skill=%q) tool_use observed in the stream — the Break-Glass trigger did not load the recovery skill", recoverySkillArg)
	}
	if !breakGlassAgentSeen {
		return fmt.Errorf("no Agent() call matched the break-glass shape (run_in_background=true, a {worker_key}-{slug}-{stage} name, a prompt carrying Skill(skill=\"spacedock:ensign\") and ### Stage definition); observed Agent() calls: %v", breakGlassAgentDetails)
	}
	return nil
}
