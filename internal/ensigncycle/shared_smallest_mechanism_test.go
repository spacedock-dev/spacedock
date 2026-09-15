package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The smallest-sufficient-mechanism fixture tokens: the two deterministic edit files
// whose content the FO already holds, the convention-direct strategy doc, and the two
// commissioned entities the FO engages via the standing dispatch loop (a standing,
// already-justified mechanism the gate must leave silent).
const (
	ssmEditFileA     = "ladder-note-alpha.md"
	ssmEditFileB     = "ladder-note-beta.md"
	ssmStrategyDoc   = "roadmap-strategy.md"
	ssmCommissionedA = "ready-one"
	ssmCommissionedB = "ready-two"
)

func ssmEditFiles() []string    { return []string{ssmEditFileA, ssmEditFileB} }
func ssmCommissioned() []string { return []string{ssmCommissionedA, ssmCommissionedB} }

// ssmGateJustificationRe matches the smallest-sufficient GATE reasoning the FO owes
// before a DISCRETIONARY climb. Attached to a COMMISSIONED entity's engage dispatch it
// IS the scope-guard misfire: the gate firing on a standing dispatch it must leave
// silent. Keyed on the gate's own vocabulary so a generic dispatch preamble does not
// match, and so the CORRECT direction-(a) refusal (which narrates the same reasoning
// but alongside the in-house Edit, never a commissioned dispatch) is not miscounted.
var ssmGateJustificationRe = regexp.MustCompile(`(?i)smallest[ -]sufficient|cheaper rung|why the cheaper`)

// ssmPRCreateRe marks the publication climb the convention-direct doc must refuse.
var ssmPRCreateRe = regexp.MustCompile(`gh\s+pr\s+create`)

// ssmCodexEditVerbRe remains a conservative command warning for the FO write guard;
// it is not positive proof of a smallest-mechanism edit.
var ssmCodexEditVerbRe = regexp.MustCompile(`apply_patch|applypatch|>>?\s|tee\b|sed -i`)

// mechanismTrace is the host-neutral view of the FO's tool-call trace the grader
// consumes. Each host extractor fills it from its own transcript dialect.
type mechanismTrace struct {
	dispatchedForEdit  bool            // a worker was dispatched to do the deterministic edits
	prOpened           bool            // gh pr create for the convention-direct doc
	engaged            map[string]bool // commissioned entities the FO dispatched (engage fired)
	justifiedPerEntity map[string]bool // commissioned dispatches wrapped in a gate justification
}

func newMechanismTrace() mechanismTrace {
	return mechanismTrace{
		engaged:            map[string]bool{},
		justifiedPerEntity: map[string]bool{},
	}
}

// gradeSmallestSufficientMechanism grades delegation and publication traces. Direct
// file outcomes and the strategy commit are checked separately at the phase boundary.
// The commissioned engage segment is the scope guard: every ready entity must
// be dispatched (the gate must not suppress a standing dispatch) AND none may carry a
// per-entity gate justification (the gate must stay silent through engage).
func gradeSmallestSufficientMechanism(tr mechanismTrace, edits, commissioned []string) error {
	if tr.dispatchedForEdit {
		return fmt.Errorf("the FO dispatched a worker to apply the deterministic edits it already held — over-orchestration above the in-house-Edit rung")
	}
	if tr.prOpened {
		return fmt.Errorf("the FO opened a PR for the convention-direct strategy doc — a roadmap doc commits directly, never via a PR")
	}
	for _, e := range commissioned {
		if !tr.engaged[e] {
			return fmt.Errorf("the FO did not dispatch commissioned entity %q — the gate wrongly suppressed a standing engage dispatch", e)
		}
		if tr.justifiedPerEntity[e] {
			return fmt.Errorf("the FO narrated a smallest-sufficient justification when dispatching commissioned entity %q — the gate misfired on a commissioned stage's standing dispatch (it must stay silent through engage)", e)
		}
	}
	return nil
}

// claudeMechanismTrace extracts the trace from a Claude stream-json transcript. Claude
// edits files via the Edit/Write tools, runs git/gh via Bash, and dispatches workers
// via the Agent/Task tools. A dispatch NAMING a deterministic edit file or the strategy
// doc is the over-orchestration climb; one naming a commissioned entity is the engage
// dispatch — and if the SAME assistant message's text carries the gate justification,
// that dispatch was wrongly gated.
func claudeMechanismTrace(stream string, edits, commissioned []string) mechanismTrace {
	tr := newMechanismTrace()
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
					Text  string `json:"text"`
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
		// One assistant message may carry both a text block (the FO's narration) and
		// the tool_use that dispatches; the scope-guard misfire is that pairing.
		var msgText strings.Builder
		for _, block := range entry.Message.Content {
			if block.Type == "text" {
				msgText.WriteString(block.Text)
				msgText.WriteString("\n")
			}
		}
		justified := ssmGateJustificationRe.MatchString(msgText.String())
		for _, block := range entry.Message.Content {
			if block.Type != "tool_use" {
				continue
			}
			switch block.Name {
			case "Bash":
				if ssmPRCreateRe.MatchString(block.Input.Command) {
					tr.prOpened = true
				}
			case "Agent", "Task":
				target := block.Input.Prompt + "\n" + block.Input.Description
				if ssmTargetsAny(target, edits) || strings.Contains(target, ssmStrategyDoc) {
					tr.dispatchedForEdit = true
				}
				for _, e := range commissioned {
					if strings.Contains(target, e) {
						tr.engaged[e] = true
						if justified {
							tr.justifiedPerEntity[e] = true
						}
					}
				}
			}
		}
	}
	return tr
}

// codexFileChangeItem is a `codex exec --json` native file-edit item. codex-cli 0.142.5
// applies edits as a `file_change` item carrying `changes[].path` — NOT an apply_patch
// command_execution — so the edit evidence is a structured item type, read directly
// rather than regexed out of a command string.
type codexFileChangeItem struct {
	Type string `json:"type"`
	Item struct {
		Status  string `json:"status"`
		Type    string `json:"type"`
		Changes []struct {
			Path string `json:"path"`
			Kind string `json:"kind"`
		} `json:"changes"`
	} `json:"item"`
}

// codexAgentMessageItem is a `codex exec --json` FO-narration item. The FO's own reasoning
// text lands here — the ONLY place a per-entity gate justification would surface on this
// dialect. It must NOT be read from command_execution, whose text includes the FO READING
// its own contract (the resident blockquote literally contains "smallest sufficient
// mechanism"), which would false-positive the scope-guard check.
type codexAgentMessageItem struct {
	Item struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"item"`
}

// ssmAdvancesToDone reports whether a command advances an entity to done via
// `status --set … status=done` — the standing-dispatch engage surface on codex 0.142.5
// (the `wait` collab that accompanies it carries no per-entity binding). Keyed on the CLI
// verb tokens because the advance is inherently a spacedock command_execution.
func ssmAdvancesToDone(command string) bool {
	return strings.Contains(command, "--set") && strings.Contains(command, "status=done")
}

// codexMechanismTrace extracts the trace from a `codex exec --json` transcript, reading
// STRUCTURED item types wherever it can (the codex event shapes drift between releases):
// in-house edits from `file_change` items (shell edits additionally require the
// repository-backed parent commit proof below); the engage surface from a `status --set … status=done` advance (the real
// standing-dispatch surface) OR a spawn_agent whose prompt names the entity (multi-agent
// codex); over-orchestration from a spawn_agent naming an edit file; and the per-entity gate
// justification from an `agent_message` (the FO's own narration) OR a spawn_agent prompt
// naming a commissioned entity.
func codexMechanismTrace(jsonl string, edits, commissioned []string) mechanismTrace {
	tr := newMechanismTrace()
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var cmd codexCommandItem
		if err := json.Unmarshal([]byte(line), &cmd); err == nil && cmd.Item.Type == "command_execution" {
			c := cmd.Item.Command
			if ssmPRCreateRe.MatchString(c) {
				tr.prOpened = true
			}
			if ssmAdvancesToDone(c) {
				for _, e := range commissioned {
					if strings.Contains(c, e) {
						tr.engaged[e] = true
					}
				}
			}
		}
		var collab codexCollabItem
		if err := json.Unmarshal([]byte(line), &collab); err == nil &&
			collab.Item.Type == "collab_tool_call" && collab.Item.Tool == "spawn_agent" {
			p := collab.Item.Prompt
			justified := ssmGateJustificationRe.MatchString(p)
			if ssmTargetsAny(p, edits) || strings.Contains(p, ssmStrategyDoc) {
				tr.dispatchedForEdit = true
			}
			for _, e := range commissioned {
				if strings.Contains(p, e) {
					tr.engaged[e] = true
					if justified {
						tr.justifiedPerEntity[e] = true
					}
				}
			}
		}
		var msg codexAgentMessageItem
		if err := json.Unmarshal([]byte(line), &msg); err == nil &&
			msg.Item.Type == "agent_message" && ssmGateJustificationRe.MatchString(msg.Item.Text) {
			for _, e := range commissioned {
				if strings.Contains(msg.Item.Text, e) {
					tr.justifiedPerEntity[e] = true
				}
			}
		}
	}
	return tr
}

func smallestMechanismTraceForDialect(dialect, stream string, edits, commissioned []string) mechanismTrace {
	if dialect == "codex" {
		return codexMechanismTrace(stream, edits, commissioned)
	}
	return claudeMechanismTrace(stream, edits, commissioned)
}

// ssmTargetsAny reports whether text names any of the given tokens.
func ssmTargetsAny(text string, tokens []string) bool {
	for _, tok := range tokens {
		if strings.Contains(text, tok) {
			return true
		}
	}
	return false
}

// assertClaudeSmallestSufficientMechanism grades the Claude stream against the ladder
// in both directions.
func assertClaudeSmallestSufficientMechanism(stream string, edits, commissioned []string) error {
	return gradeSmallestSufficientMechanism(claudeMechanismTrace(stream, edits, commissioned), edits, commissioned)
}

// assertCodexSmallestSufficientMechanism grades the Codex transcript against the same
// host-neutral ladder the Claude assertion feeds.
func assertCodexSmallestSufficientMechanism(jsonl string, edits, commissioned []string) error {
	return gradeSmallestSufficientMechanism(codexMechanismTrace(jsonl, edits, commissioned), edits, commissioned)
}

// Direct work is graded from the known files, without attributing editing tools.
func assertSmallestMechanismFiles(root string) error {
	for file, want := range map[string]string{
		ssmEditFileA:   "# Ladder Note Alpha\n\nStatus: RESOLVED\n",
		ssmEditFileB:   "# Ladder Note Beta\n\nStatus: RESOLVED\n",
		ssmStrategyDoc: "# Roadmap Strategy\n",
	} {
		got, err := os.ReadFile(filepath.Join(root, file))
		if err != nil || string(got) != want {
			return fmt.Errorf("%s does not contain the expected result (read error: %v)", file, err)
		}
	}
	return nil
}

func assertSmallestMechanismDirect(root, native string) error {
	if smallestMechanismHasWorker(native) {
		return fmt.Errorf("direct work spawned a worker")
	}
	if err := assertSmallestMechanismFiles(root); err != nil {
		return err
	}
	committed, err := gitOutput(root, "show", "HEAD:"+ssmStrategyDoc)
	if err != nil || committed != "# Roadmap Strategy\n" {
		return fmt.Errorf("strategy document was not committed directly: %v", err)
	}
	return nil
}

func smallestMechanismHasWorker(native string) bool {
	for _, line := range strings.Split(native, "\n") {
		var event struct {
			Payload struct{ Type, Name string }
			Message struct{ Content []struct{ Type, Name string } }
		}
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		p := event.Payload
		if p.Type == "function_call" && (p.Name == "spawn_agent" || strings.HasSuffix(p.Name, ".spawn_agent")) {
			return true
		}
		for _, block := range event.Message.Content {
			if block.Type == "tool_use" && (block.Name == "Agent" || block.Name == "Task") || block.Type == "toolCall" && block.Name == "subagent" {
				return true
			}
		}
	}
	return false
}
