package ensigncycle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// The keep-moving-posture scenario grades the FO's MOTION TRACE — its tool calls plus
// its turn-ending final message — against the four false-stop patterns the 0223 Shaping
// FO retrospective recorded (docs/roadmap/0223-pi-dispatch-contract/debrief-shaping-2026-06-19.md):
// post-approval pause, sequential-instead-of-parallel dispatch, turn-end on an async
// launch with independent work remaining, and a captain correction halting the whole
// session. The durable end-state cannot tell keep-moving from a false stop — the entity
// files look the same whether the FO advanced+dispatched or asked permission and stopped
// — so the ACTIONS (advance, dispatch, re-shape) are read from the transcript and the two
// turn-ending POSTURES (a permission question; a passive wait announcement) are read from
// the FINAL MESSAGE only. Reading the postures from the final message, not from all
// narration, avoids the contract-read false-positive: the FO's mid-stream reasoning quotes
// the shared-core keep-moving blockquote ("Yield only when blocked on the async result
// with no other work"), which a stream-wide phrase scan would misread as a wait. Like the
// smallest-mechanism trace the actions are host-specific (Claude tool_use blocks vs Codex
// command_execution / file_change / spawn_agent items), so each host has an extractor that
// fills the host-neutral keepMovingTrace the shared grader consumes. Default build tags
// (stdlib JSON only) so the offline negative exercises them without a model.

// The keep-moving fixture entities. The captain has just approved kmApprovedGate's gate
// (the FO must advance it + dispatch its next stage); kmReadyOne / kmReadyTwo are
// independent and ready for the same stage (dispatch both in one motion); kmQuestioned's
// mechanism the captain has questioned (re-shape its body, pause its dispatch — never
// halt the independent ones).
const (
	kmApprovedGate = "approved-gate"
	kmReadyOne     = "ready-one"
	kmReadyTwo     = "ready-two"
	kmQuestioned   = "questioned"
	kmNextStage    = "implementation"
)

func kmIndependent() []string { return []string{kmReadyOne, kmReadyTwo} }

// kmPermissionRe matches the post-approval permission question (pattern 1): the FO asking
// the captain's leave to take the reversible advance+dispatch step the approval already
// triggered. Drive-cue: the 2026-07-06/07 Commander session's self-imposed gate stall
// under a standing conn grant. Keyed on an ask verb near a proceed verb so a plain status
// report ("I advanced and dispatched") does not match.
var kmPermissionRe = regexp.MustCompile(`(?i)(want me to|shall i|should i|do you want me to|would you like me to|let me know (if|whether)|awaiting your (go|ok|approval|sign-?off))\b[^.?!]{0,48}?\b(advance|dispatch|proceed|continue|kick off|go ahead|engage|move forward|move on)`)

// kmWaitRe matches a turn-ending passive-wait posture (pattern 3): the FO ending its turn
// to wait on an async result rather than continuing with the independent work that
// remains. Drive-cue: the 2026-07-06/07 Commander session's FO ending its turn on a CI
// wait with no event source. Keyed on a WAIT-to-resume posture, NOT on "I'll report back
// once they finish" — reporting on completion is fine; suspending further work until the
// async returns is the violation.
var kmWaitRe = regexp.MustCompile(`(?i)(waiting (for|on) (the|it|them|these|those|that|results?|the ensign|the worker|the async|ci)|i'?ll wait (for|on|until)|i will wait (for|on|until)|await(ing)? (the|its|their) (results?|completion|return|outcome)|will (proceed|resume|continue|move on|do more|pick .* up) (once|when|after)|pausing until|stand(ing)? by (for|until)|hold(ing)? off until|no further action until)`)

// kmAdvancesEntity reports whether a command advances the named entity to a new stage via
// `status --set <entity> status=<stage>` — the FO's advance surface after a gate approval.
func kmAdvancesEntity(command, entity string) bool {
	return strings.Contains(command, "--set") &&
		strings.Contains(command, entity) &&
		strings.Contains(command, "status=")
}

// kmAdvancesToStatus reports whether a command sets the named entity to a specific status
// via `status --set <entity> status=<status>`. The codex 0.142.5 standing-loop dispatch
// surface advances a dispatched entity to `done` (the `wait` collab that accompanies it
// carries no per-entity binding), so `status=done` is the per-entity dispatch evidence on
// that dialect — there is no spawn_agent to read.
func kmAdvancesToStatus(command, entity, status string) bool {
	return strings.Contains(command, "--set") &&
		strings.Contains(command, entity) &&
		strings.Contains(command, "status="+status)
}

// keepMovingTrace is the host-neutral view of the FO's motion the grader consumes. Each
// host extractor fills it from its own transcript dialect plus the shared final message.
type keepMovingTrace struct {
	approvedAdvanced      bool            // the approved entity was advanced past its gate to the next stage
	approvedDispatched    bool            // the approved entity's next stage was dispatched
	independentDispatched map[string]bool // each independent ready entity the FO dispatched
	correctedReshaped     bool            // the questioned entity's body was re-shaped to fold the correction
	correctedDriven       bool            // the questioned entity was advanced or dispatched despite the pending correction
	askedPermission       bool            // the final message asks the captain's leave to advance/dispatch (pattern 1)
	endedWaitingOnAsync   bool            // the final message ends the turn on a passive async wait (pattern 3)
}

func newKeepMovingTrace() keepMovingTrace {
	return keepMovingTrace{independentDispatched: map[string]bool{}}
}

// gradeKeepMoving is host-neutral: it holds the FO's motion trace to the four keep-moving
// patterns. Approval triggers the advance AND the dispatch with no permission question
// (S1); every independent ready entity is dispatched, not serialized behind a pause (S2);
// the questioned entity is re-shaped and its dispatch paused while the independent ones
// keep moving (S4); and the turn does not end on a passive async wait while work remains
// (S3). The scenario is constructed so independent work always remains, so a turn-ending
// wait is never correct here — the forced decision point per the k74g categorical-signal
// spike.
func gradeKeepMoving(tr keepMovingTrace, independent []string) error {
	if tr.askedPermission {
		return fmt.Errorf("the FO asked the captain's permission to advance/dispatch after the gate approval — approval triggers the next action, it is not a turn boundary (advance+dispatch is reversible work the contract already permits)")
	}
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
		return fmt.Errorf("the FO did not re-shape the questioned entity %q to fold the captain's correction — a correction narrows scope by re-shaping the affected entity, not by dropping it", kmQuestioned)
	}
	if tr.correctedDriven {
		return fmt.Errorf("the FO advanced or dispatched the questioned entity %q while its mechanism was under correction — pause only the corrected entity's dispatch until the re-shape is folded", kmQuestioned)
	}
	if tr.endedWaitingOnAsync {
		return fmt.Errorf("the FO ended its turn to wait on an async result while independent work remained — launching an async ensign does not end the turn while independent FO work remains")
	}
	return nil
}

// claudeKeepMovingTrace extracts the motion trace from a Claude stream-json transcript
// (actions) plus the final message (turn-ending postures). Claude advances via a `status
// --set` Bash command, dispatches via the Agent/Task tools, and re-shapes via the
// Edit/Write tools.
func claudeKeepMovingTrace(stream, finalMessage string, independent []string) keepMovingTrace {
	tr := newKeepMovingTrace()
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Input struct {
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
			if block.Type != "tool_use" {
				continue
			}
			switch block.Name {
			case "Bash":
				if kmAdvancesEntity(block.Input.Command, kmApprovedGate) {
					tr.approvedAdvanced = true
				}
				if kmAdvancesEntity(block.Input.Command, kmQuestioned) {
					tr.correctedDriven = true
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
				if strings.Contains(target, kmQuestioned) {
					tr.correctedDriven = true
				}
			}
		}
	}
	tr.askedPermission = kmPermissionRe.MatchString(finalMessage)
	tr.endedWaitingOnAsync = kmWaitRe.MatchString(finalMessage)
	return tr
}

// codexKeepMovingTrace extracts the motion trace from a `codex exec --json` transcript
// (actions) plus the final message (turn-ending postures). It reads structured item types:
// advances from a `status --set … status=` command_execution; re-shapes from a `file_change`
// item (codex 0.142.5) OR an apply_patch/shell-write command touching the questioned entity;
// dispatches from a spawn_agent whose prompt names the entity.
func codexKeepMovingTrace(jsonl, finalMessage string, independent []string) keepMovingTrace {
	tr := newKeepMovingTrace()
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var cmd codexCommandItem
		if err := json.Unmarshal([]byte(line), &cmd); err == nil && cmd.Item.Type == "command_execution" {
			c := cmd.Item.Command
			if kmAdvancesEntity(c, kmApprovedGate) {
				tr.approvedAdvanced = true
			}
			// codex 0.142.5 dispatches via the standing loop (no spawn_agent): the
			// dispatched worker advances the entity to `done`, so a `status --set <e>
			// status=done` is the per-entity dispatch evidence.
			if kmAdvancesToStatus(c, kmApprovedGate, "done") {
				tr.approvedDispatched = true
			}
			for _, e := range independent {
				if kmAdvancesToStatus(c, e, "done") {
					tr.independentDispatched[e] = true
				}
			}
			if kmAdvancesEntity(c, kmQuestioned) {
				tr.correctedDriven = true
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
				tr.correctedDriven = true
			}
		}
	}
	tr.askedPermission = kmPermissionRe.MatchString(finalMessage)
	tr.endedWaitingOnAsync = kmWaitRe.MatchString(finalMessage)
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
