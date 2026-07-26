package ensigncycle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// The keep-moving-posture scenario grades the FO's MOTION TRACE — its tool calls plus its
// turn-ending final message — against the 0223 Shaping FO false-stop patterns
// (docs/roadmap/0223-pi-dispatch-contract/debrief-shaping-2026-06-19.md): post-approval
// pause, sequential-instead-of-parallel dispatch, turn-end on an async launch with work
// remaining, and a captain correction halting the session. The durable end-state cannot
// tell keep-moving from a false stop, so the ACTIONS (advance, dispatch, re-shape) are read
// from the transcript. The post-approval pause is proven by those actions: an approved entity
// that was advanced AND dispatched did not pause for permission, so the structured
// advance/dispatch evidence is the no-false-stop proof — the free-form summary wording is not
// scanned for question phrasing. The one posture the actions cannot show — whether the
// corrected entity's re-shape was surfaced or silently parked — is read from the FINAL
// MESSAGE, not from all narration, so the FO's mid-stream reasoning quoting the shared-core
// keep-moving blockquote is not misread by a stream-wide phrase scan.
//
// A captain correction on one entity means an invariant TRIPLE (the shipped S4 clause and
// this grader key on the same thing): (1) the corrected entity does NOT drive FORWARD
// (advance to implementation/done, or merge) until the re-shape folds; (2) everything else
// keeps moving; (3) once folded, the re-shape SURFACES to the captain (a gate review /
// re-presentation) OR its rework is honestly still in flight and the stop condition names
// it — silent absorption and silent parking both fail. The re-shape MECHANISM is free: an
// in-house edit or a routed rework are both legitimate (host write-scopes differ — codex
// MUST route — and the smallest-sufficient-mechanism principle governs the choice), so
// correctedDriven keys on the FORWARD status only, never on a re-shape edit/route/dispatch.
// Like the smallest-mechanism trace the actions are host-specific; each host extractor
// fills the host-neutral keepMovingTrace the shared grader consumes. Default build tags so
// the offline negative exercises them without a model.

// The keep-moving fixture entities. The captain has just approved kmApprovedGate's gate
// (advance + dispatch its next stage); kmReadyOne / kmReadyTwo are independent and ready for
// the same stage (dispatch both in one motion); kmQuestioned's mechanism the captain has
// questioned (re-shape it and hold it from driving forward, then surface the re-shape —
// never halt the independent ones).
const (
	kmApprovedGate = "approved-gate"
	kmReadyOne     = "ready-one"
	kmReadyTwo     = "ready-two"
	kmQuestioned   = "questioned"
	kmNextStage    = "implementation"
	kmReopenStage  = "ideation"
)

func kmIndependent() []string { return []string{kmReadyOne, kmReadyTwo} }

// kmCorrectedDispositionRe matches, in the final message, a substantive disposition of the
// corrected entity's re-shape — SURFACED (a gate review / decision / re-presentation for the
// captain) or honestly IN-FLIGHT (folding / re-shaping / a routed rework named as running,
// to surface once the gate is presented). Paired with the corrected entity being named, it
// distinguishes an addressed re-shape from the silent wait/park failure specimen (an FO that
// dispatches everything and ends "no work remains, I'll wait" without naming the correction).
var kmCorrectedDispositionRe = regexp.MustCompile(`(?i)gate review|recommend|decision|present|surfac|approv|reject|fold|re-?shap|rework|dispatch|in[- ]flight|running|pending|gates? presented`)

// kmCorrectedAddressed reports whether the final message addresses the corrected entity's
// re-shape at the drive's stop — it names the corrected entity AND states a disposition
// (surfaced for review, or an honestly-in-flight rework). A final message that never names
// the corrected entity has silently parked or absorbed the re-shape.
func kmCorrectedAddressed(finalMessage, corrected string) bool {
	if !strings.Contains(strings.ToLower(finalMessage), corrected) {
		return false
	}
	return kmCorrectedDispositionRe.MatchString(finalMessage)
}

// kmAdvancesToStatus reports whether the command sets the named entity to a specific status
// via `--set <entity> … status=<status>`. It scopes each `--set` to its OWN entity: a bulk
// command that sets several entities in one call (e.g. `--set approved-gate
// status=implementation … --set questioned status=ideation`) must not cross-attribute
// approved-gate's status to questioned — the entity must be the token right after `--set`,
// and the `status=` must fall within that one `--set`'s args (bounded at the next command
// separator). The codex 0.142.5 standing-loop dispatch surface advances a dispatched entity
// to `done` (the `wait` collab that accompanies it carries no per-entity binding), so
// `status=done` is the per-entity dispatch evidence on that dialect.
func kmAdvancesToStatus(command, entity, status string) bool {
	for _, seg := range strings.Split(command, "--set")[1:] {
		seg = strings.TrimSpace(seg)
		if seg != entity && !strings.HasPrefix(seg, entity+" ") {
			continue
		}
		if i := strings.IndexAny(seg, "\n;&|"); i >= 0 {
			seg = seg[:i]
		}
		if strings.Contains(seg, "status="+status) {
			return true
		}
	}
	return false
}

// kmAdvancesForward reports whether a command drives the named entity FORWARD through its
// pipeline — to the next working stage (implementation) or to terminal (done, which also
// covers a merge). A route BACK to an earlier stage (ideation) to re-open the design is a
// re-shape, not a forward drive, and is deliberately not matched here.
func kmAdvancesForward(command, entity string) bool {
	return kmAdvancesToStatus(command, entity, kmNextStage) || kmAdvancesToStatus(command, entity, "done")
}

func kmMergeGuardTerminalizes(command, entity string) bool {
	fields := strings.Fields(command)
	for i := 0; i+2 < len(fields); i++ {
		if fields[i] == "merge" && fields[i+1] == "guard" && fields[i+2] == entity {
			return true
		}
	}
	return false
}

func kmConsumesGate(command, entity string) bool {
	for _, segment := range strings.FieldsFunc(command, func(r rune) bool {
		return r == '\n' || r == ';' || r == '&' || r == '|'
	}) {
		fields := strings.Fields(segment)
		for i := 1; i+2 < len(fields); i++ {
			launcher := strings.Trim(fields[i-1], `"'`)
			if fields[i] != "gate" || fields[i+1] != "consume" || strings.Trim(fields[i+2], `"'`) != entity ||
				!(launcher == "spacedock" || strings.HasSuffix(launcher, "/spacedock") ||
					strings.Contains(launcher, "SPACEDOCK_BIN") || strings.HasPrefix(launcher, "$")) {
				continue
			}
			validPrefix := true
			for _, prefix := range fields[:i-1] {
				prefix = strings.Trim(prefix, `"'`)
				if prefix != "command" && prefix != "env" && prefix != "bash" && prefix != "/bin/bash" &&
					prefix != "sh" && prefix != "/bin/sh" && prefix != "-c" && prefix != "-lc" &&
					!strings.Contains(prefix, "=") {
					validPrefix = false
					break
				}
			}
			if validPrefix {
				return true
			}
		}
	}
	return false
}

func kmConsumeSucceeded(output, entity string) bool {
	gatePrefix := "gate=gate:" + entity + ":"
	for _, line := range strings.Split(output, "\n") {
		sawGate, sawConsumed := false, false
		for _, field := range strings.Fields(line) {
			if strings.HasPrefix(field, gatePrefix) {
				sawGate = true
			}
			if field == "consumed=true" {
				sawConsumed = true
			}
		}
		if sawGate && sawConsumed {
			return true
		}
	}
	return false
}

func kmClaudeToolResultText(raw json.RawMessage) string {
	var text string
	if json.Unmarshal(raw, &text) == nil {
		return text
	}
	var blocks []struct {
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &blocks) != nil {
		return ""
	}
	var texts []string
	for _, block := range blocks {
		texts = append(texts, block.Text)
	}
	return strings.Join(texts, "\n")
}

// keepMovingTrace is the host-neutral view of the FO's motion the grader consumes. Each host
// extractor fills it from its own transcript dialect plus the shared final message.
type keepMovingTrace struct {
	approvedAdvanced      bool            // the approved entity was advanced past its gate to the next stage
	approvedDispatched    bool            // the approved entity's next stage was dispatched
	independentDispatched map[string]bool // each independent ready entity the FO dispatched
	correctedReshaped     bool            // the corrected entity's design was re-shaped (in-house edit, routed rework, or re-open to ideation)
	correctedDriven       bool            // the corrected entity was driven FORWARD (to implementation/done/merge) before the re-shape folded
	correctedAddressed    bool            // at stop, the re-shape was surfaced or its rework honestly named as in-flight
}

func newKeepMovingTrace() keepMovingTrace {
	return keepMovingTrace{independentDispatched: map[string]bool{}}
}

// gradeKeepMoving is host-neutral: it holds the FO's motion trace to the keep-moving
// patterns. Approval triggers the advance AND the dispatch — performing both is the proof the
// FO did not pause for permission (S1); every independent ready entity is dispatched, not
// serialized behind a pause (S2); the corrected entity is re-shaped, held from driving
// forward, and its re-shape surfaced or honestly named as in-flight while the independent
// ones keep moving (S4). The re-shape mechanism is free — only a FORWARD drive of the
// corrected entity or a silent park fails.
func gradeKeepMoving(tr keepMovingTrace, independent []string) error {
	if !tr.approvedAdvanced {
		return fmt.Errorf("the FO did not advance the approved entity %q past its gate to the next stage after the captain's approval", kmApprovedGate)
	}
	if !tr.approvedDispatched {
		return fmt.Errorf("the FO advanced the approved entity %q but did not dispatch its next stage — approval triggers the advance AND the dispatch before yielding", kmApprovedGate)
	}
	for _, e := range independent {
		if !tr.independentDispatched[e] {
			return fmt.Errorf("the FO did not dispatch independent ready entity %q — independent entities dispatch in parallel in one motion, not one-at-a-time with a pause between", e)
		}
	}
	if !tr.correctedReshaped {
		return fmt.Errorf("the FO did not re-shape the corrected entity %q to fold the captain's correction — a correction narrows scope by re-shaping the affected entity, not by dropping it", kmQuestioned)
	}
	if tr.correctedDriven {
		return fmt.Errorf("the FO drove the corrected entity %q forward (to implementation/done/merge) before the re-shape folded — the corrected entity does not advance forward until the correction is folded", kmQuestioned)
	}
	if !tr.correctedAddressed {
		return fmt.Errorf("the FO did not surface the corrected entity %q's re-shape or honestly name its rework as in-flight at the drive's stop — a folded re-shape must surface to the captain (a gate review / re-presentation); silent absorption and a silent wait both fail", kmQuestioned)
	}
	return nil
}

// claudeKeepMovingTrace extracts the motion trace from a Claude stream-json transcript
// (actions) plus the final message. Claude advances through a successful gate consume or
// legacy `status --set`, re-opens via `status --set`, dispatches via Agent/Task, and
// re-shapes in-house via Edit/Write.
func claudeKeepMovingTrace(stream, finalMessage string, independent []string) keepMovingTrace {
	tr := newKeepMovingTrace()
	pendingConsumes := map[string]bool{}
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Type    string `json:"type"`
			Message *struct {
				Content []struct {
					Type      string          `json:"type"`
					ID        string          `json:"id"`
					Name      string          `json:"name"`
					ToolUseID string          `json:"tool_use_id"`
					IsError   bool            `json:"is_error"`
					Content   json.RawMessage `json:"content"`
					Input     struct {
						Command     string `json:"command"`
						FilePath    string `json:"file_path"`
						Prompt      string `json:"prompt"`
						Description string `json:"description"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if entry.Type == "user" && block.Type == "tool_result" && pendingConsumes[block.ToolUseID] &&
				!block.IsError && kmConsumeSucceeded(kmClaudeToolResultText(block.Content), kmApprovedGate) {
				tr.approvedAdvanced = true
			}
			if block.Type != "tool_use" {
				continue
			}
			switch block.Name {
			case "Bash":
				c := block.Input.Command
				if kmConsumesGate(c, kmApprovedGate) && block.ID != "" {
					pendingConsumes[block.ID] = true
				}
				if kmAdvancesToStatus(c, kmApprovedGate, kmNextStage) {
					tr.approvedAdvanced = true
				}
				if kmAdvancesForward(c, kmQuestioned) {
					tr.correctedDriven = true
				}
				if kmAdvancesToStatus(c, kmQuestioned, kmReopenStage) {
					tr.correctedReshaped = true
				}
			case "Edit", "Write":
				if strings.Contains(block.Input.FilePath, kmQuestioned) {
					tr.correctedReshaped = true
				}
			case "Agent", "Task":
				target := block.Input.Prompt + "\n" + block.Input.Description
				if strings.Contains(target, kmApprovedGate) {
					tr.approvedDispatched = true
				}
				for _, e := range independent {
					if strings.Contains(target, e) {
						tr.independentDispatched[e] = true
					}
				}
				// A dispatch naming the corrected entity is a routed re-shape/rework — a
				// legitimate re-shape mechanism, never a forward drive (that is graded from
				// the entity's FORWARD status only).
				if strings.Contains(target, kmQuestioned) {
					tr.correctedReshaped = true
				}
			}
		}
	}
	tr.correctedAddressed = kmCorrectedAddressed(finalMessage, kmQuestioned)
	return tr
}

// codexKeepMovingTrace extracts the motion trace from a `codex exec --json` transcript
// (actions) plus the final message. It reads successful gate consumes, legacy status
// advances/re-opens, and the standing-loop dispatch surface from command_execution;
// re-shapes from file_change or apply_patch; and dispatches from spawn_agent.
func codexKeepMovingTrace(jsonl, finalMessage string, independent []string) keepMovingTrace {
	tr := newKeepMovingTrace()
	evidenceEntities := append([]string{kmApprovedGate, kmQuestioned}, independent...)
	dispatchEvidence := codexDispatchCompletionEvidenceFromJSONL(jsonl, evidenceEntities)
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var cmd struct {
			Type string `json:"type"`
			Item struct {
				Type     string `json:"type"`
				Command  string `json:"command"`
				Status   string `json:"status"`
				Output   string `json:"aggregated_output"`
				ExitCode *int   `json:"exit_code"`
			} `json:"item"`
		}
		if err := json.Unmarshal([]byte(line), &cmd); err == nil && cmd.Item.Type == "command_execution" {
			c := cmd.Item.Command
			if cmd.Type == "item.completed" && cmd.Item.Status == "completed" && cmd.Item.ExitCode != nil &&
				*cmd.Item.ExitCode == 0 && kmConsumesGate(c, kmApprovedGate) &&
				kmConsumeSucceeded(cmd.Item.Output, kmApprovedGate) {
				tr.approvedAdvanced = true
			}
			if kmAdvancesToStatus(c, kmApprovedGate, kmNextStage) {
				tr.approvedAdvanced = true
			}
			// codex 0.142.5 dispatches via the standing loop (no spawn_agent): the dispatched
			// worker advances the entity to `done`, so a `status --set <e> status=done` is the
			// per-entity dispatch evidence.
			if kmAdvancesToStatus(c, kmApprovedGate, "done") || kmMergeGuardTerminalizes(c, kmApprovedGate) {
				tr.approvedDispatched = true
			}
			for _, e := range independent {
				if kmAdvancesToStatus(c, e, "done") || kmMergeGuardTerminalizes(c, e) {
					tr.independentDispatched[e] = true
				}
			}
			if kmAdvancesForward(c, kmQuestioned) || kmMergeGuardTerminalizes(c, kmQuestioned) {
				tr.correctedDriven = true
			}
			if kmAdvancesToStatus(c, kmQuestioned, kmReopenStage) {
				tr.correctedReshaped = true
			}
			if ssmCodexEditVerbRe.MatchString(c) && strings.Contains(c, kmQuestioned) {
				tr.correctedReshaped = true
			}
		}
		var fc codexFileChangeItem
		if err := json.Unmarshal([]byte(line), &fc); err == nil && fc.Item.Type == "file_change" {
			for _, ch := range fc.Item.Changes {
				if strings.Contains(ch.Path, kmQuestioned) {
					tr.correctedReshaped = true
				}
			}
		}
		var collab codexCollabItem
		if err := json.Unmarshal([]byte(line), &collab); err == nil &&
			collab.Item.Type == "collab_tool_call" && collab.Item.Tool == "spawn_agent" {
			p := collab.Item.Prompt
			if strings.Contains(p, kmApprovedGate) {
				tr.approvedDispatched = true
			}
			for _, e := range independent {
				if strings.Contains(p, e) {
					tr.independentDispatched[e] = true
				}
			}
			if strings.Contains(p, kmQuestioned) {
				tr.correctedReshaped = true
			}
		}
	}
	if dispatchEvidence.stageReport[kmApprovedGate] {
		tr.approvedDispatched = true
	}
	for _, e := range independent {
		if dispatchEvidence.stageReport[e] {
			tr.independentDispatched[e] = true
		}
	}
	if dispatchEvidence.stageReport[kmQuestioned] {
		tr.correctedReshaped = true
	}
	tr.correctedAddressed = kmCorrectedAddressed(finalMessage, kmQuestioned)
	return tr
}

// assertClaudeKeepMoving grades the Claude stream + final message against the keep-moving
// patterns.
func assertClaudeKeepMoving(stream, finalMessage string, independent []string) error {
	return gradeKeepMoving(claudeKeepMovingTrace(stream, finalMessage, independent), independent)
}

// assertCodexKeepMoving grades the Codex transcript + final message against the same
// host-neutral patterns the Claude assertion feeds.
func assertCodexKeepMoving(jsonl, finalMessage string, independent []string) error {
	return gradeKeepMoving(codexKeepMovingTrace(jsonl, finalMessage, independent), independent)
}

func kmClaudeCompletedBash(id, command string, failed bool) string {
	return kmClaudeCompletedBashOutput(id, command, "command finished", failed)
}

func kmClaudeCompletedBashOutput(id, command, output string, failed bool) string {
	return fmt.Sprintf(
		"{\"type\":\"assistant\",\"message\":{\"content\":[{\"type\":\"tool_use\",\"id\":%q,\"name\":\"Bash\",\"input\":{\"command\":%q}}]}}\n"+
			"{\"type\":\"user\",\"message\":{\"content\":[{\"type\":\"tool_result\",\"tool_use_id\":%q,\"content\":%q,\"is_error\":%t}]}}",
		id, command, id, output, failed,
	)
}

func kmCodexCompletedCommand(command string, exitCode int) string {
	return kmCodexCompletedCommandOutput(command, "command finished", exitCode)
}

func kmCodexCompletedCommandOutput(command, output string, exitCode int) string {
	return fmt.Sprintf(
		`{"type":"item.completed","item":{"type":"command_execution","command":%q,"aggregated_output":%q,"status":"completed","exit_code":%d}}`,
		command, output, exitCode,
	)
}

func TestKeepMovingGateConsumeAdvancement(t *testing.T) {
	independent := kmIndependent()
	claudeRemainder := strings.Join([]string{
		claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyOne+`."}`),
		claudeToolUse("Agent", `{"prompt":"Dispatch `+kmNextStage+` for `+kmReadyTwo+`."}`),
		claudeToolUse("Edit", `{"file_path":"`+kmQuestioned+`.md"}`),
	}, "\n")
	codexRemainder := strings.Join([]string{
		codexSpawn("Dispatch the " + kmNextStage + " stage for " + kmApprovedGate + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyOne + "."),
		codexSpawn("Dispatch " + kmNextStage + " for " + kmReadyTwo + "."),
		codexFileChange(kmQuestioned + ".md"),
	}, "\n")

	consumeOutput := "gate=gate:" + kmApprovedGate + ":review application=advance/consumed condition=approved-pending eligible=true consumed=true target-stage=" + kmNextStage
	claudeConsume := kmClaudeCompletedBashOutput("consume-approved", "spacedock gate consume "+kmApprovedGate+" --workflow-dir .", consumeOutput, false)
	if err := assertClaudeKeepMoving(claudeConsume+"\n"+claudeRemainder, kmCorrectFinal(), independent); err != nil {
		t.Fatalf("successful Claude gate consume must count as the approved transition: %v", err)
	}
	codexConsume := kmCodexCompletedCommandOutput(
		"spacedock gate record "+kmApprovedGate+" --decision approve; spacedock state commit "+kmApprovedGate+
			"; spacedock gate consume "+kmApprovedGate+" --workflow-dir .; spacedock state commit "+kmApprovedGate,
		consumeOutput,
		0,
	)
	if err := assertCodexKeepMoving(codexConsume+"\n"+codexRemainder, kmCorrectFinal(), independent); err != nil {
		t.Fatalf("successful Codex gate consume must count as the approved transition: %v", err)
	}

	claudeControls := map[string]string{
		"absent":                  claudeRemainder,
		"other entity":            kmClaudeCompletedBash("consume-other", "spacedock gate consume other-task --workflow-dir .", false) + "\n" + claudeRemainder,
		"failed":                  kmClaudeCompletedBash("consume-failed", "spacedock gate consume "+kmApprovedGate+" --workflow-dir .", true) + "\n" + claudeRemainder,
		"failed then success":     kmClaudeCompletedBash("consume-false-green", "spacedock gate consume "+kmApprovedGate+" --workflow-dir .; true", false) + "\n" + claudeRemainder,
		"failure suppressed true": kmClaudeCompletedBash("consume-suppressed", "spacedock gate consume "+kmApprovedGate+" --workflow-dir . || true", false) + "\n" + claudeRemainder,
	}
	for name, stream := range claudeControls {
		t.Run("claude/"+name, func(t *testing.T) {
			if err := assertClaudeKeepMoving(stream, kmCorrectFinal(), independent); err == nil || !strings.Contains(err.Error(), "advance") {
				t.Fatalf("control must fail on the missing approved advance, got %v", err)
			}
		})
	}
	codexControls := map[string]string{
		"absent":                  codexRemainder,
		"other entity":            kmCodexCompletedCommand("spacedock gate consume other-task --workflow-dir .", 0) + "\n" + codexRemainder,
		"failed":                  kmCodexCompletedCommand("spacedock gate consume "+kmApprovedGate+" --workflow-dir .", 1) + "\n" + codexRemainder,
		"failed then success":     kmCodexCompletedCommand("spacedock gate consume "+kmApprovedGate+" --workflow-dir .; true", 0) + "\n" + codexRemainder,
		"failure suppressed true": kmCodexCompletedCommand("spacedock gate consume "+kmApprovedGate+" --workflow-dir . || true", 0) + "\n" + codexRemainder,
	}
	for name, stream := range codexControls {
		t.Run("codex/"+name, func(t *testing.T) {
			if err := assertCodexKeepMoving(stream, kmCorrectFinal(), independent); err == nil || !strings.Contains(err.Error(), "advance") {
				t.Fatalf("control must fail on the missing approved advance, got %v", err)
			}
		})
	}

	claudeNoDispatch := strings.Replace(claudeRemainder, claudeToolUse("Agent", `{"prompt":"Dispatch the `+kmNextStage+` stage for `+kmApprovedGate+`."}`)+"\n", "", 1)
	if err := assertClaudeKeepMoving(claudeConsume+"\n"+claudeNoDispatch, kmCorrectFinal(), independent); err == nil || !strings.Contains(err.Error(), "dispatch") {
		t.Fatalf("Claude consume without a separate approved dispatch must fail on dispatch, got %v", err)
	}
	codexNoDispatch := strings.Replace(codexRemainder, codexSpawn("Dispatch the "+kmNextStage+" stage for "+kmApprovedGate+".")+"\n", "", 1)
	if err := assertCodexKeepMoving(codexConsume+"\n"+codexNoDispatch, kmCorrectFinal(), independent); err == nil || !strings.Contains(err.Error(), "dispatch") {
		t.Fatalf("Codex consume without a separate approved dispatch must fail on dispatch, got %v", err)
	}
}
