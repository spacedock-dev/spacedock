// ABOUTME: Offline oracles over a Claude stream-json transcript for the two
// ABOUTME: fo-dispatch-recovery live scenarios (AC-2 degraded-bare, AC-3 break-glass).
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// inputBoolField treats an omitted JSON boolean as false.
func inputBoolField(input map[string]json.RawMessage, key string) bool {
	return string(input[key]) == "true"
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
// bare-shaped Agent() call (no name and false-or-omitted `run_in_background`) — bare dispatch is
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
				if inputStringField(block.Input, "name") == "" && !inputBoolField(block.Input, "run_in_background") {
					sawBareAgent = true
				}
			}
		}
	})

	if !sawBareAgent {
		return fmt.Errorf("no bare-shaped Agent() call (no name and false-or-omitted `run_in_background`) observed — bare dispatch was not reached")
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
type dispatchMode string

const (
	dispatchModeBare dispatchMode = "bare"
	dispatchModeTeam dispatchMode = "team"
)

func assertBreakGlassObservables(stream string, selectedMode dispatchMode) error {
	sawReportBeforeAgent := false
	sawAgentCall := false
	sawRecoverySkill := false
	var agentCount int
	var matchingAgentCount int
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
				agentCount++
				var runInBackground bool
				runRaw, hasRunInBackground := block.Input["run_in_background"]
				if hasRunInBackground {
					_ = json.Unmarshal(runRaw, &runInBackground)
				}
				name := inputStringField(block.Input, "name")
				_, hasName := block.Input["name"]
				_, hasTeamName := block.Input["team_name"]
				prompt := inputStringField(block.Input, "prompt")
				description := inputStringField(block.Input, "description")
				subagentType := inputStringField(block.Input, "subagent_type")
				hasEnsignSkill := strings.Contains(prompt, `Skill(skill="spacedock:ensign")`)
				hasStageDef := strings.Contains(prompt, "### Stage definition:") && strings.Contains(prompt, dispatchRecoveryStageDefinition)
				hasStageReport := strings.Contains(prompt, "### Stage report")
				hasCompletionSignal := strings.Contains(prompt, `SendMessage(to="team-lead"`)
				nameShaped := strings.Count(name, "-") >= 2
				commonShape := subagentType == "spacedock:ensign" && description != "" && hasEnsignSkill && hasStageDef && hasStageReport
				modeShape := false
				switch selectedMode {
				case dispatchModeBare:
					modeShape = !hasName && !hasTeamName && !hasRunInBackground && !hasCompletionSignal
				case dispatchModeTeam:
					modeShape = nameShaped && runInBackground && hasRunInBackground && !hasTeamName && hasCompletionSignal
				}
				details := fmt.Sprintf("name=%q name-present=%v team_name-present=%v run_in_background=%v run-present=%v subagent_type=%q description-present=%v ensign-skill-in-prompt=%v stage-def-in-prompt=%v stage-report-in-prompt=%v completion-signal=%v", name, hasName, hasTeamName, runInBackground, hasRunInBackground, subagentType, description != "", hasEnsignSkill, hasStageDef, hasStageReport, hasCompletionSignal)
				breakGlassAgentDetails = append(breakGlassAgentDetails, details)
				if commonShape && modeShape {
					matchingAgentCount++
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
	if agentCount != 1 {
		return fmt.Errorf("observed %d Agent() calls, want exactly one; calls: %v", agentCount, breakGlassAgentDetails)
	}
	if matchingAgentCount != 1 {
		return fmt.Errorf("the only Agent() call did not preserve selected %s mode or carry the complete ensign prompt; observed Agent() calls: %v", selectedMode, breakGlassAgentDetails)
	}
	return nil
}

// assertBreakGlassDurableResult proves the worker result, report protocol, scoped
// commit, and clean entity path independently from the Claude stream.
func assertBreakGlassDurableResult(root, entityPath string) error {
	content, err := os.ReadFile(entityPath)
	if err != nil {
		return fmt.Errorf("read recovery entity: %w", err)
	}
	if err := assertCompleteRecoveryReport(string(content)); err != nil {
		return err
	}
	rel, err := filepath.Rel(root, entityPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("entity %q is outside repository %q", entityPath, root)
	}
	status, err := gitOutput(root, "status", "--short", "--", rel)
	if err != nil {
		return err
	}
	if strings.TrimSpace(status) != "" {
		return fmt.Errorf("entity has uncommitted changes: %s", strings.TrimSpace(status))
	}
	history, err := gitOutput(root, "log", "--format=%H", "--", rel)
	if err != nil {
		return err
	}
	for _, commit := range strings.Fields(history) {
		blob, showErr := gitOutput(root, "show", commit+":"+filepath.ToSlash(rel))
		if showErr != nil || assertCompleteRecoveryReport(blob) != nil {
			continue
		}
		files, diffErr := gitOutput(root, "diff-tree", "--root", "--no-commit-id", "--name-only", "-r", commit)
		if diffErr == nil && strings.TrimSpace(files) == filepath.ToSlash(rel) {
			return nil
		}
	}
	return fmt.Errorf("no path-scoped commit contains the complete recovery result")
}

func assertCompleteRecoveryReport(content string) error {
	for _, want := range []string{dispatchRecoveryMarker, "## Stage Report: implementation", "- DONE:", "### Summary"} {
		if !strings.Contains(content, want) {
			return fmt.Errorf("recovery entity missing %q", want)
		}
	}
	return nil
}

func gitOutput(root string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
