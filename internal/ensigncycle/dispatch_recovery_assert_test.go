// ABOUTME: Offline oracles over a Claude stream-json transcript for the two
// ABOUTME: fo-dispatch-recovery live scenarios (AC-2 degraded-bare, AC-3 break-glass).
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"strings"
)

// degradedModeCaptainReportPrefix is the fixed lead of the verbatim captain report
// sentence Diff 2's `## Degraded Mode` skill section mandates on trip. Matching the
// fixed prefix (not the whole sentence) tolerates trailing whitespace/newline
// variation a real model turn can introduce without weakening the observable: the
// distinguishing content is the "Falling back to bare mode ... infrastructure
// failure" clause, which no other FO output plausibly emits verbatim.
const degradedModeCaptainReportPrefix = "Falling back to bare mode for the remainder of this session due to infrastructure failure."

// recoverySkillArg is the Skill tool_use argument both live scenarios must observe
// after their trigger — the resident trigger line's `Skill(skill="spacedock:fo-dispatch-recovery")` load.
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

// assertDegradedBareObservables is the AC-2 behavioral oracle: over the captured
// stream it asserts (i) a Skill(skill="spacedock:fo-dispatch-recovery") tool_use
// appears, (ii) the verbatim captain-report sentence appears in a text block, and
// (iii) every Agent() tool_use that follows the report text omits BOTH `name` and
// `run_in_background` (the bare-mode shape Degraded Mode's Effects mandate). It does
// not require a specific ordering between (i) and (ii)/(iii) beyond "after the
// trigger" — the trigger rides in the initial `-p` prompt (HEADLESS transport, per
// the M4 decision), so the whole stream is post-trigger.
func assertDegradedBareObservables(stream string) error {
	sawRecoverySkill := false
	sawCaptainReport := false
	var nonBareAgentCalls []string

	walkStreamBlocks(stream, func(block streamContentBlock) {
		switch block.Type {
		case "text":
			if strings.Contains(block.Text, degradedModeCaptainReportPrefix) {
				sawCaptainReport = true
			}
		case "tool_use":
			switch block.Name {
			case "Skill":
				if inputStringField(block.Input, "skill") == recoverySkillArg {
					sawRecoverySkill = true
				}
			case "Agent":
				if inputHasKey(block.Input, "name") || inputHasKey(block.Input, "run_in_background") {
					nonBareAgentCalls = append(nonBareAgentCalls, fmt.Sprintf("name-present=%v run_in_background-present=%v",
						inputHasKey(block.Input, "name"), inputHasKey(block.Input, "run_in_background")))
				}
			}
		}
	})

	if !sawRecoverySkill {
		return fmt.Errorf("no Skill(skill=%q) tool_use observed in the stream — the Degraded Mode trigger did not load the recovery skill", recoverySkillArg)
	}
	if !sawCaptainReport {
		return fmt.Errorf("the verbatim captain report sentence (prefix %q) was not observed in any text block", degradedModeCaptainReportPrefix)
	}
	if len(nonBareAgentCalls) > 0 {
		return fmt.Errorf("%d Agent() call(s) carried `name` or `run_in_background` after Degraded Mode tripped — bare mode requires BOTH omitted: %v", len(nonBareAgentCalls), nonBareAgentCalls)
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
