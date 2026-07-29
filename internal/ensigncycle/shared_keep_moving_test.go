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

var kmGateConsume = regexp.MustCompile(`(?:^|[\s;&|])['"]?(?:spacedock|\$(?:\{SPACEDOCK_BIN(?::-[^}]*)?\}|[A-Za-z_][A-Za-z0-9_]*)|/[^ \t\r\n'";&|]+/spacedock)['"]?\s+gate\s+consume\s+['"]?approved-gate['"]?(?:\s|$)`)
var kmEligibleTrue = regexp.MustCompile(`(?:^|\s)eligible=true(?:\s|$)`)
var kmConsumedTrue = regexp.MustCompile(`(?:^|\s)consumed=true(?:\s|$)`)
var kmImplementationTarget = regexp.MustCompile(`(?:^|\s)target-stage=implementation(?:\s|$)`)

func kmSuccessfulGateConsumeResult(content json.RawMessage, isError *bool) bool {
	var text string
	return isError != nil && !*isError && json.Unmarshal(content, &text) == nil &&
		kmConsumedTrue.MatchString(text) && kmImplementationTarget.MatchString(text)
}

type codexKeepMovingCommandItem struct {
	Type string `json:"type"`
	Item struct {
		Type             string `json:"type"`
		Command          string `json:"command"`
		AggregatedOutput string `json:"aggregated_output"`
		ExitCode         *int   `json:"exit_code"`
		Status           string `json:"status"`
	} `json:"item"`
}

func kmSuccessfulCodexGateConsume(cmd codexKeepMovingCommandItem) bool {
	return cmd.Type == "item.completed" &&
		cmd.Item.Type == "command_execution" &&
		cmd.Item.Status == "completed" &&
		cmd.Item.ExitCode != nil && *cmd.Item.ExitCode == 0 &&
		kmGateConsume.MatchString(cmd.Item.Command) &&
		kmEligibleTrue.MatchString(cmd.Item.AggregatedOutput) &&
		kmConsumedTrue.MatchString(cmd.Item.AggregatedOutput) &&
		kmImplementationTarget.MatchString(cmd.Item.AggregatedOutput)
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
// (actions) plus the final message. Claude advances via a successful recorded `gate consume`
// or `status --set`, re-opens via `status --set`, dispatches (including a routed re-shape)
// via the Agent/Task tools, and re-shapes in-house via the Edit/Write tools.
func claudeKeepMovingTrace(stream, finalMessage string, independent []string) keepMovingTrace {
	tr := newKeepMovingTrace()
	gateConsumeInvocations := map[string]bool{}
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type      string `json:"type"`
					Name      string `json:"name"`
					ID        string `json:"id"`
					ToolUseID string `json:"tool_use_id"`
					Input     struct {
						Command     string `json:"command"`
						FilePath    string `json:"file_path"`
						Prompt      string `json:"prompt"`
						Description string `json:"description"`
					} `json:"input"`
					Content json.RawMessage `json:"content"`
					IsError *bool           `json:"is_error"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type == "tool_result" && gateConsumeInvocations[block.ToolUseID] &&
				kmSuccessfulGateConsumeResult(block.Content, block.IsError) {
				tr.approvedAdvanced = true
			}
			if block.Type != "tool_use" {
				continue
			}
			switch block.Name {
			case "Bash":
				c := block.Input.Command
				if kmAdvancesToStatus(c, kmApprovedGate, kmNextStage) {
					tr.approvedAdvanced = true
				}
				if block.ID != "" && kmGateConsume.MatchString(c) {
					gateConsumeInvocations[block.ID] = true
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
// (actions) plus the final message. It reads structured item types: advances/re-opens and
// the standing-loop dispatch surface from command_execution; re-shapes from a `file_change`
// item (codex 0.142.5) or an apply_patch command; dispatches (including a routed re-shape)
// from a spawn_agent naming the entity.
func codexKeepMovingTrace(jsonl, finalMessage string, independent []string) keepMovingTrace {
	tr := newKeepMovingTrace()
	evidenceEntities := append([]string{kmApprovedGate, kmQuestioned}, independent...)
	dispatchEvidence := codexDispatchCompletionEvidenceFromJSONL(jsonl, evidenceEntities)
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var cmd codexKeepMovingCommandItem
		if err := json.Unmarshal([]byte(line), &cmd); err == nil && cmd.Item.Type == "command_execution" {
			c := cmd.Item.Command
			if kmAdvancesToStatus(c, kmApprovedGate, kmNextStage) || kmSuccessfulCodexGateConsume(cmd) {
				tr.approvedAdvanced = true
			}
			// codex 0.142.5 dispatches via the standing loop (no spawn_agent): the dispatched
			// worker advances the entity to `done`, so a `status --set <e> status=done` is the
			// per-entity dispatch evidence.
			if kmAdvancesToStatus(c, kmApprovedGate, "done") {
				tr.approvedDispatched = true
			}
			for _, e := range independent {
				if kmAdvancesToStatus(c, e, "done") {
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

func TestClaudeKeepMovingGateConsumeCorrelation(t *testing.T) {
	command := `${SPACEDOCK_BIN:-spacedock} gate consume approved-gate --workflow-dir "$WD"`
	result := "gate=gate:approved-gate:review application=advance/consumed consumed=true target-stage=implementation"
	for _, tt := range []struct {
		name, command, result, resultID string
		isError, want                   bool
	}{
		{"success", command, result, "consume", false, true},
		{"wrong_entity", strings.Replace(command, "approved-gate", "ready-one", 1), result, "consume", false, false},
		{"failed", command, result, "consume", true, false},
		{"missing_result", command, result, "", false, false},
		{"uncorrelated_result", command, result, "other", false, false},
		{"not_consumed", command, strings.Replace(result, "consumed=true", "consumed=false", 1), "consume", false, false},
		{"wrong_target", command, strings.Replace(result, "target-stage=implementation", "target-stage=validation", 1), "consume", false, false},
	} {
		stream := bashToolLine("consume", tt.command)
		if tt.resultID != "" {
			stream += "\n" + toolResultLine(tt.resultID, tt.isError, tt.result)
		}
		if got := claudeKeepMovingTrace(stream, "", nil).approvedAdvanced; got != tt.want {
			t.Errorf("%s: approvedAdvanced = %t, want %t", tt.name, got, tt.want)
		}
	}
}

func TestCodexKeepMovingGateConsumeCorrelation(t *testing.T) {
	// Runtime Live E2E run 30421227237 assigned the resolved launcher to SD,
	// then consumed the approval through "$SD". Preserve that exact command
	// dialect: the successful structured result is the advancement evidence.
	command := `/bin/bash -lc 'if [ -n "${SPACEDOCK_BIN:-}" ] && [ -x "${SPACEDOCK_BIN}" ]; then SD="${SPACEDOCK_BIN}"; else SD=spacedock; fi
"$SD" gate consume approved-gate --workflow-dir /tmp/TestLiveCodexSharedScenarioskeep-moving-posture3759480657/001'`
	result := "gate=gate:approved-gate:review application=advance/consumed condition=approved-pending eligible=true consumed=true target-stage=implementation"
	for _, tt := range []struct {
		name, command, result, status string
		exit, want                    int
	}{
		{"retained_success", command, result, "completed", 0, 1},
		{"failed", command, result, "failed", 1, 0},
		{"wrong_entity", strings.Replace(command, "gate consume approved-gate", "gate consume ready-one", 1), result, "completed", 0, 0},
		{"invocation_only", command, "", "completed", 0, 0},
		{"output_only", "sed -n '1,80p' fo-gate-lifecycle.md", result, "completed", 0, 0},
		{"not_eligible", command, strings.Replace(result, "eligible=true", "eligible=false", 1), "completed", 0, 0},
		{"not_consumed", command, strings.Replace(result, "consumed=true", "consumed=false", 1), "completed", 0, 0},
		{"wrong_target", command, strings.Replace(result, "target-stage=implementation", "target-stage=validation", 1), "completed", 0, 0},
	} {
		got := codexKeepMovingTrace(codexCommandOutput(tt.command, tt.result, tt.exit, tt.status), "", nil).approvedAdvanced
		if got != (tt.want == 1) {
			t.Errorf("%s: approvedAdvanced = %t, want %t", tt.name, got, tt.want == 1)
		}
	}

	advanceWithoutDispatch := strings.Join([]string{
		codexCommandOutput(command, result, 0, "completed"),
		codexSpawn("Dispatch implementation for " + kmReadyOne + "."),
		codexSpawn("Dispatch implementation for " + kmReadyTwo + "."),
		codexFileChange(kmQuestioned + ".md"),
	}, "\n")
	err := assertCodexKeepMoving(advanceWithoutDispatch, kmCorrectFinal(), kmIndependent())
	if err == nil || !strings.Contains(err.Error(), "dispatch") {
		t.Fatalf("successful consume without approved-gate dispatch must fail on dispatch, got: %v", err)
	}
}
