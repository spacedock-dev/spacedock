package ensigncycle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

type codexRolloutRecord struct {
	Type    string `json:"type"`
	Payload struct {
		Type       string          `json:"type"`
		Name       string          `json:"name"`
		CallID     string          `json:"call_id"`
		Input      string          `json:"input"`
		Arguments  string          `json:"arguments"`
		Output     json.RawMessage `json:"output"`
		Status     string          `json:"status"`
		Role       string          `json:"role"`
		Phase      string          `json:"phase"`
		Source     string          `json:"source"`
		Originator string          `json:"originator"`
		Cwd        string          `json:"cwd"`
		Message    string          `json:"message"`
		Content    []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"payload"`
}

type codexRolloutPendingCall struct {
	name   string
	input  string
	status string
}

// codexRolloutIsInteractive distinguishes the TUI's durable rollout from a
// headless `codex exec` rollout. The source/originator pair is written by Codex,
// not supplied by the scenario prompt.
func codexRolloutIsInteractive(stream string) bool {
	for _, record := range parseCodexRollout(stream) {
		if record.Type == "session_meta" {
			return record.Payload.Source == "cli" && record.Payload.Originator == "codex-tui"
		}
	}
	return false
}

// codexRolloutReachedIdle requires a completed task after the latest task start.
// In a TUI this marks the assistant's turn end; the independent resident-process
// probe proves the host then returned to its input loop rather than exiting.
func codexRolloutReachedIdle(stream string) bool {
	idle := false
	for _, record := range parseCodexRollout(stream) {
		if record.Type != "event_msg" {
			continue
		}
		switch record.Payload.Type {
		case "task_started":
			idle = false
		case "task_complete":
			idle = true
		}
	}
	return idle
}

func extractCodexInteractiveFinalMessage(stream string) (string, error) {
	var fallback, final string
	for _, record := range parseCodexRollout(stream) {
		switch {
		case record.Type == "response_item" && record.Payload.Type == "message" && record.Payload.Role == "assistant":
			var text []string
			for _, block := range record.Payload.Content {
				if block.Type == "output_text" && block.Text != "" {
					text = append(text, block.Text)
				}
			}
			if len(text) == 0 {
				continue
			}
			fallback = strings.Join(text, "\n")
			if record.Payload.Phase == "final_answer" {
				final = fallback
			}
		case record.Type == "event_msg" && record.Payload.Type == "agent_message" && record.Payload.Message != "":
			fallback = record.Payload.Message
			if record.Payload.Phase == "final_answer" {
				final = fallback
			}
		}
	}
	if final != "" {
		return final, nil
	}
	if fallback != "" {
		return fallback, nil
	}
	return "", fmt.Errorf("no assistant final message in Codex interactive rollout")
}

// normalizeCodexInteractiveFOReferenceEvents maps the TUI's durable rollout
// dialect to the same host-neutral order events as `codex exec --json`. A call is
// resolved only when its paired tool output reports success and contains the
// canonical reference anchors.
func normalizeCodexInteractiveFOReferenceEvents(stream string) []foReferenceEvent {
	records := parseCodexRollout(stream)
	skillBase := codexInteractiveFirstOfficerBase(records)
	pending := make(map[string]codexRolloutPendingCall)
	var events []foReferenceEvent
	for _, record := range records {
		if codexRolloutToolCall(record) {
			pending[record.Payload.CallID] = codexRolloutPendingCall{
				name:   record.Payload.Name,
				input:  codexRolloutCallInput(record),
				status: record.Payload.Status,
			}
			continue
		}
		if !codexRolloutToolOutput(record) {
			continue
		}
		call, ok := pending[record.Payload.CallID]
		if !ok {
			continue
		}
		delete(pending, record.Payload.CallID)
		statusText, output := codexRolloutOutputText(record.Payload.Output)
		succeeded := codexRolloutCallSucceeded(call, statusText)
		callEvents := classifyFOCommand(call.input, skillBase)
		if isCodexRolloutMutationTool(call.name, call.input) {
			callEvents = append(callEvents, foMutation)
		}
		events = append(events, resolvedFOCallEvents(callEvents, foCallSucceeded(callEvents, succeeded, output))...)
		if succeeded && containsMergeModBlock(output) {
			events = append(events, foModBlockSeen)
		}
	}
	return events
}

func parseCodexRollout(stream string) []codexRolloutRecord {
	var records []codexRolloutRecord
	scanner := bufio.NewScanner(strings.NewReader(stream))
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	for scanner.Scan() {
		var record codexRolloutRecord
		if json.Unmarshal(scanner.Bytes(), &record) == nil {
			records = append(records, record)
		}
	}
	return records
}

func codexRolloutToolCall(record codexRolloutRecord) bool {
	return record.Type == "response_item" && (record.Payload.Type == "custom_tool_call" || record.Payload.Type == "function_call") && record.Payload.CallID != ""
}

func codexRolloutToolOutput(record codexRolloutRecord) bool {
	return record.Type == "response_item" && (record.Payload.Type == "custom_tool_call_output" || record.Payload.Type == "function_call_output") && record.Payload.CallID != ""
}

func codexRolloutCallInput(record codexRolloutRecord) string {
	raw := record.Payload.Input
	if raw == "" {
		raw = record.Payload.Arguments
	}
	var args struct {
		Cmd     string `json:"cmd"`
		Command string `json:"command"`
	}
	if json.Unmarshal([]byte(raw), &args) == nil {
		if args.Cmd != "" {
			return args.Cmd
		}
		if args.Command != "" {
			return args.Command
		}
	}
	if match := codexRolloutJSCmdRE.FindStringSubmatch(raw); len(match) == 2 {
		var command string
		if json.Unmarshal([]byte(match[1]), &command) == nil {
			return command
		}
	}
	if commands := codexRolloutJSTemplateCommands(raw); len(commands) > 0 {
		return strings.Join(commands, "\n")
	}
	return raw
}

var (
	codexRolloutJSCmdRE         = regexp.MustCompile(`(?s)\bcmd\s*:\s*("(?:\\.|[^"\\])*")`)
	codexRolloutJSTemplateCmdRE = regexp.MustCompile("(?s)\\bcmd\\s*:\\s*`((?:\\\\.|[^`])*)`")
	codexRolloutJSStringBinding = regexp.MustCompile(`\bconst\s+([A-Za-z_][A-Za-z0-9_]*)\s*=\s*("(?:\\.|[^"\\])*")\s*;`)
	codexRolloutSkillPathRE     = regexp.MustCompile(`<path>[[:space:]]*([^<\r\n]+/skills/first-officer)/SKILL\.md[[:space:]]*</path>`)
)

// codexRolloutJSTemplateCommands normalizes the real Codex TUI tool-call shape:
// an exec custom tool may bind an exact path in JavaScript and interpolate it
// into one or more exec_command `cmd` templates. Resolve only const JSON-string
// bindings and only their exact ${name} occurrences; arbitrary expressions stay
// unresolved and therefore cannot satisfy the exact-base oracle.
func codexRolloutJSTemplateCommands(raw string) []string {
	bindings := make(map[string]string)
	for _, match := range codexRolloutJSStringBinding.FindAllStringSubmatch(raw, -1) {
		var value string
		if json.Unmarshal([]byte(match[2]), &value) == nil {
			bindings[match[1]] = value
		}
	}
	var commands []string
	for _, match := range codexRolloutJSTemplateCmdRE.FindAllStringSubmatch(raw, -1) {
		command := match[1]
		for name, value := range bindings {
			command = strings.ReplaceAll(command, "${"+name+"}", value)
		}
		commands = append(commands, command)
	}
	return commands
}

func codexInteractiveFirstOfficerBase(records []codexRolloutRecord) string {
	for _, record := range records {
		if record.Type == "response_item" && record.Payload.Type == "message" {
			for _, block := range record.Payload.Content {
				if match := firstOfficerBaseAnnouncementRE.FindStringSubmatch(block.Text); len(match) == 2 {
					return strings.TrimSpace(match[1])
				}
				if match := codexRolloutSkillPathRE.FindStringSubmatch(block.Text); len(match) == 2 {
					return strings.TrimSpace(match[1])
				}
			}
		}
	}
	pending := make(map[string]codexRolloutPendingCall)
	for _, record := range records {
		if codexRolloutToolCall(record) {
			pending[record.Payload.CallID] = codexRolloutPendingCall{
				name:   record.Payload.Name,
				input:  codexRolloutCallInput(record),
				status: record.Payload.Status,
			}
			continue
		}
		if !codexRolloutToolOutput(record) {
			continue
		}
		call, ok := pending[record.Payload.CallID]
		if !ok {
			continue
		}
		delete(pending, record.Payload.CallID)
		statusText, output := codexRolloutOutputText(record.Payload.Output)
		if !codexRolloutCallSucceeded(call, statusText) || !strings.Contains(strings.ToLower(output), "name: first-officer") {
			continue
		}
		if match := firstOfficerEntryPathRE.FindStringSubmatch(call.input); len(match) == 2 {
			return strings.TrimSpace(match[1])
		}
	}
	return ""
}

func codexRolloutOutputText(raw json.RawMessage) (string, string) {
	if len(raw) == 0 {
		return "", ""
	}
	var blocks []struct {
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &blocks) == nil {
		var texts []string
		for _, block := range blocks {
			if block.Text != "" {
				texts = append(texts, block.Text)
			}
		}
		if len(texts) == 0 {
			return "", ""
		}
		return texts[0], strings.Join(texts, "\n")
	}
	var text string
	if json.Unmarshal(raw, &text) == nil {
		return text, text
	}
	return "", string(raw)
}

func codexRolloutCallSucceeded(call codexRolloutPendingCall, statusText string) bool {
	if strings.EqualFold(call.status, "failed") {
		return false
	}
	status := strings.ToLower(statusText)
	for _, failure := range []string{"script failed", "script exited", "process exited with code 1", "timed out", "terminated"} {
		if strings.Contains(status, failure) {
			return false
		}
	}
	if strings.Contains(status, "script completed") || strings.Contains(status, "process exited with code 0") {
		return true
	}
	return strings.EqualFold(call.status, "completed") && statusText != ""
}

func isCodexRolloutMutationTool(name, input string) bool {
	lowerName := strings.ToLower(name)
	lowerInput := strings.ToLower(input)
	return strings.Contains(lowerName, "apply_patch") || strings.Contains(lowerName, "write_file") ||
		strings.Contains(lowerInput, "tools.apply_patch(")
}

func TestCodexInteractiveRolloutRequiresTUIIdleAndOrderedReads(t *testing.T) {
	stream := codexInteractiveRolloutFixture("cli", "codex-tui", "Script completed\nWall time 0.0 seconds\nOutput:\n")
	if !codexRolloutIsInteractive(stream) {
		t.Fatal("TUI rollout was not recognized as interactive")
	}
	if !codexRolloutReachedIdle(stream) {
		t.Fatal("completed TUI turn was not recognized as idle")
	}
	message, err := extractCodexInteractiveFinalMessage(stream)
	if err != nil || !strings.Contains(message, "gate-check") {
		t.Fatalf("interactive final message = %q, %v", message, err)
	}
	if err := assertFOReferenceJourney(normalizeCodexInteractiveFOReferenceEvents(stream), "gate"); err != nil {
		t.Fatalf("interactive ordered gate stream: %v", err)
	}
}

func TestCodexInteractiveRolloutResolvesExactJSBaseTemplates(t *testing.T) {
	const skillBase = "/plugin/skills/first-officer"
	input := `const base="` + skillBase + `";
const r = await tools.exec_command({cmd:` + "`" + `cat "${base}/references/first-officer-shared-core.md" && cat "${base}/references/codex-first-officer-runtime.md"` + "`" + `}); text(r.output);`
	stream := strings.Join([]string{
		`{"type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"<skill><path>` + skillBase + `/SKILL.md</path></skill>"}]}}`,
		`{"type":"response_item","payload":{"type":"custom_tool_call","name":"exec","call_id":"refs","status":"completed","input":` + mustJSONString(input) + `}}`,
		`{"type":"response_item","payload":{"type":"custom_tool_call_output","call_id":"refs","output":[{"type":"text","text":"Script completed\nWall time 0.0 seconds\nOutput:\n"},{"type":"text","text":"# First Officer Shared Core\n# Codex First Officer Runtime\n"}]}}`,
	}, "\n")
	records := parseCodexRollout(stream)
	events := normalizeCodexInteractiveFOReferenceEvents(stream)
	if err := assertFOReferenceJourney(events, "gate"); err != nil {
		t.Fatalf("exact const-base template reads: %v; base=%q input=%q", err, codexInteractiveFirstOfficerBase(records), codexRolloutCallInput(records[1]))
	}

	wrongBase := strings.Replace(input, skillBase, "/wrong/skills/first-officer", 1)
	wrongStream := strings.Replace(stream, mustJSONString(input), mustJSONString(wrongBase), 1)
	if err := assertFOReferenceJourney(normalizeCodexInteractiveFOReferenceEvents(wrongStream), "gate"); err == nil {
		t.Fatal("a JS template bound to a base other than the loader-supplied skill base passed")
	}

	unbound := strings.Replace(input, `const base="`+skillBase+`";`, `const base=process.env.UNRELATED;`, 1)
	unboundStream := strings.Replace(stream, mustJSONString(input), mustJSONString(unbound), 1)
	if err := assertFOReferenceJourney(normalizeCodexInteractiveFOReferenceEvents(unboundStream), "gate"); err == nil {
		t.Fatal("an unresolved JS template expression passed as an exact-base read")
	}
}

func TestCodexInteractiveRolloutRejectsHeadlessAndFailedPreconditions(t *testing.T) {
	headless := codexInteractiveRolloutFixture("exec", "codex_exec", "Script completed\nWall time 0.0 seconds\nOutput:\n")
	if codexRolloutIsInteractive(headless) {
		t.Fatal("headless codex exec rollout was accepted as interactive evidence")
	}
	busy := strings.Replace(headless, `{"type":"event_msg","payload":{"type":"task_complete"}}`, "", 1)
	if codexRolloutReachedIdle(busy) {
		t.Fatal("rollout without task_complete was accepted as an idle greet")
	}
	failed := codexInteractiveRolloutFixture("cli", "codex-tui", "Script exited with code 1\nWall time 0.0 seconds\nOutput:\n")
	if err := assertFOReferenceJourney(normalizeCodexInteractiveFOReferenceEvents(failed), "gate"); err == nil {
		t.Fatal("failed shared/runtime reads satisfied interactive gate preconditions")
	}
}

func codexInteractiveRolloutFixture(source, originator, readStatus string) string {
	return strings.Join([]string{
		`{"type":"session_meta","payload":{"source":"` + source + `","originator":"` + originator + `","cwd":"/tmp/workflow"}}`,
		`{"type":"event_msg","payload":{"type":"task_started"}}`,
		`{"type":"response_item","payload":{"type":"custom_tool_call","name":"exec","call_id":"skill","status":"completed","input":"const r = await tools.exec_command({cmd: \"sed -n '1,100p' /plugin/skills/first-officer/SKILL.md\"}); text(r.output);"}}`,
		`{"type":"response_item","payload":{"type":"custom_tool_call_output","call_id":"skill","output":[{"type":"text","text":"Script completed\nWall time 0.0 seconds\nOutput:\n"},{"type":"text","text":"name: first-officer\n"}]}}`,
		`{"type":"response_item","payload":{"type":"custom_tool_call","name":"exec","call_id":"refs","status":"completed","input":"const r = await tools.exec_command({cmd: \"cat /plugin/skills/first-officer/references/first-officer-shared-core.md /plugin/skills/first-officer/references/codex-first-officer-runtime.md\"}); text(r.output);"}}`,
		`{"type":"response_item","payload":{"type":"custom_tool_call_output","call_id":"refs","output":[{"type":"text","text":` + mustJSONString(readStatus) + `},{"type":"text","text":"# First Officer Shared Core\n# Codex First Officer Runtime\n"}]}}`,
		`{"type":"response_item","payload":{"type":"message","role":"assistant","phase":"final_answer","content":[{"type":"output_text","text":"gate-check is ready; merged-pr awaits a live check. Use engage to continue."}]}}`,
		`{"type":"event_msg","payload":{"type":"task_complete"}}`,
	}, "\n")
}
