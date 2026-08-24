package ensigncycle

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Pi-dialect extractors for the rejection-flow scenario. The rejection-flow graders
// use Claude-dialect extractors by default (claudeRecordedRejectionRound,
// claudeRejectionRoundPublications, claudeRejectionRoutes) which parse the Claude
// stream-json format. Pi records its session in a different JSONL shape —
// `type:"message"` records whose `message.content` blocks are `toolCall`/`toolResult`
// rather than Claude's `tool_use`/`tool_result` — so the Claude extractors see nothing
// in a Pi stream. These extractors read the Pi session JSONL so the Pi XFAIL can XPASS
// after the timeout repair. They are test-harness code, not a new gate command, stored
// format, authority source, or CI lane.

// piRejectionBranch is the Pi branch key. Pi's «context-budget» is ABSENT (reuse
// condition 0 satisfied), but «addressable-worker»'s reuse-advance is DEFERRED —
// "Fresh redispatch remains the default first Pi slice; normal follow-up and retry
// dispatches are fresh assignment cycles, not context resumes. Reuse-advance over a
// kept-alive handle is deferred." The reuse-advance handle is therefore not exposed
// for reuse, so reuse condition 1 fails and the reuse route does not survive. The FO
// owes the fail-safe FRESH chain, the same shape Claude's probe-unavailable fail-safe
// produces.
const piRejectionBranch = rejectionBranchFresh

// piDispatchFileInTask extracts the worker handle from the dispatch file path in a
// subagent spawn task. The FO dispatches via `spacedock dispatch build`, whose emitted
// prompt points at `/tmp/spacedock-dispatch/spacedock-ensign-{slug}-{stage}.md`; the
// handle is the filename stem, the same (slug, stage)-derived identity Claude's Agent
// `name` and Codex's task path carry.
var piDispatchFileInTask = regexp.MustCompile(`spacedock-dispatch/(spacedock-ensign-[a-z0-9-]+-(?:implementation|validation))\.md`)

// piAsyncRunID extracts the async run id a subagent spawn result returns. Pi's
// `subagent(... async: true)` result text carries "Async workflow [run-id]"; the FO
// later polls `subagent({action:"status", id:"run-id"})`. The id is typically a UUID but
// may carry non-hex label characters in test fixtures, so the character class is
// permissive within the bracket-delimited token.
var piAsyncRunID = regexp.MustCompile(`\bAsync (?:workflow|run) \[([0-9a-zA-Z_-]+)\]`)

// piStatusCompleted reports whether a subagent status result text shows the worker
// completed. Pi's status result carries "worker completed, exit N, accept" or
// "worker completed, exit N" when the run finishes.
var piStatusCompleted = regexp.MustCompile(`worker completed,\s+exit\s+\d+`)

// piToolCallBlock is one content block in a Pi session assistant message.
type piToolCallBlock struct {
	Type      string          `json:"type"`
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// piToolCallArgs holds the argument fields the extractors read, parsed from the
// `arguments` raw JSON of a toolCall block.
type piToolCallArgs struct {
	Command string `json:"command"`
	Task    string `json:"task"`
	Action  string `json:"action"`
	ID      string `json:"id"`
}

// piTextContent extracts the concatenated text from a Pi content array (the
// `content` field of a toolResult message, which is an array of `{"type":"text","text":"…"}` blocks).
func piTextContent(raw json.RawMessage) string {
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &blocks) != nil {
		return ""
	}
	var sb strings.Builder
	for _, b := range blocks {
		if b.Type == "text" {
			sb.WriteString(b.Text)
		}
	}
	return sb.String()
}

// piSessionMessage is the `message` field of a Pi session JSONL record.
type piSessionMessage struct {
	Role       string          `json:"role"`
	ToolCallID string          `json:"toolCallId"`
	ToolName   string          `json:"toolName"`
	Content    json.RawMessage `json:"content"`
	IsError    *bool           `json:"isError"`
}

// piSessionRecord is one line of a Pi session JSONL file.
type piSessionRecord struct {
	Type    string           `json:"type"`
	Message piSessionMessage `json:"message"`
}

// piBashCommands extracts every bash/shell toolCall id→command pair from a Pi session
// JSONL, in stream order. These are the FO's shell invocations — `gate record`,
// `status`, `dispatch build`, etc.
func piBashCommands(session string) []piCommandEntry {
	var entries []piCommandEntry
	for _, line := range strings.Split(session, "\n") {
		var rec piSessionRecord
		if json.Unmarshal([]byte(line), &rec) != nil {
			continue
		}
		if rec.Message.Role != "assistant" {
			continue
		}
		var blocks []piToolCallBlock
		if json.Unmarshal(rec.Message.Content, &blocks) != nil {
			continue
		}
		for _, b := range blocks {
			if b.Type != "toolCall" || (b.Name != "bash" && b.Name != "shell") {
				continue
			}
			var args piToolCallArgs
			if json.Unmarshal(b.Arguments, &args) != nil {
				continue
			}
			entries = append(entries, piCommandEntry{id: b.ID, command: args.Command})
		}
	}
	return entries
}

// piCommandEntry pairs a toolCall id with its shell command.
type piCommandEntry struct {
	id      string
	command string
}

// piToolResults extracts every toolResult toolCallId→(text, isError) pair from a Pi
// session JSONL. The result text is the concatenated content of the toolResult's
// content blocks.
func piToolResults(session string) map[string]piResultEntry {
	results := map[string]piResultEntry{}
	for _, line := range strings.Split(session, "\n") {
		var rec piSessionRecord
		if json.Unmarshal([]byte(line), &rec) != nil {
			continue
		}
		if rec.Message.Role != "toolResult" || rec.Message.ToolCallID == "" {
			continue
		}
		entry := piResultEntry{text: piTextContent(rec.Message.Content)}
		if rec.Message.IsError != nil {
			entry.failed = *rec.Message.IsError
		}
		results[rec.Message.ToolCallID] = entry
	}
	return results
}

// piResultEntry holds the text and error state of one tool result.
type piResultEntry struct {
	text   string
	failed bool
}

// piRecordedRejectionRound is the Pi half of claudeRecordedRejectionRound: it reports
// whether a `gate record --round validation/1` invocation ran and its correlated
// toolResult reported success (not an error) with the recorder's success line.
func piRecordedRejectionRound(session string) bool {
	commands := piBashCommands(session)
	results := piToolResults(session)
	for _, cmd := range commands {
		if !commandRecordsRejectionRound(cmd.command) {
			continue
		}
		result, ok := results[cmd.id]
		if !ok || result.failed {
			continue
		}
		if rejectionRoundSuccess.MatchString(result.text) {
			return true
		}
	}
	return false
}

// piRejectionRoundPublications is the Pi half of claudeRejectionRoundPublications: it
// returns the round id of every `gate record --round` invocation whose correlated
// toolResult did not report an error, in stream order.
func piRejectionRoundPublications(session string) []string {
	commands := piBashCommands(session)
	results := piToolResults(session)
	var rounds []string
	for _, cmd := range commands {
		published := rejectionRoundPublications(cmd.command)
		if len(published) == 0 {
			continue
		}
		result, ok := results[cmd.id]
		if ok && result.failed {
			continue
		}
		rounds = append(rounds, published...)
	}
	return rounds
}

// piSubagentSpawn is one subagent spawn observation: the toolCall id, the derived
// worker handle, and the stage.
type piSubagentSpawn struct {
	toolCallID string
	handle     string
	stage      string
}

// piRejectionRoutes extracts the ordered worker topology from a Pi session JSONL.
// Pi's FO spawns workers via `subagent(... async: true)` (a toolCall with a `task`
// argument), polls `subagent({action:"status", id})`, and reads the status result
// to detect completion. Worker identity is derived from the dispatch file path in the
// spawn task — the same (slug, stage)-derived handle Claude's Agent `name` carries.
// The branch is always FRESH on Pi because reuse-advance is deferred (see
// piRejectionBranch).
func piRejectionRoutes(session string) ([]rejectionRoute, rejectionBranch) {
	var routes []rejectionRoute
	// spawnByRunID maps the async run id (from the spawn result text) back to the
	// (handle, stage) the spawn opened. The FO polls status with that run id.
	spawnByRunID := map[string]piSubagentSpawn{}
	// spawnByToolCallID maps the spawn toolCall id to the spawn, so we can read the
	// run id out of the spawn's toolResult.
	spawnByToolCallID := map[string]piSubagentSpawn{}
	// pendingSpawn tracks the most recent spawn that has not yet been completed, for
	// the sequential-flow fallback when run-id correlation is unavailable.
	var pendingSpawn *piSubagentSpawn

	for i, line := range strings.Split(session, "\n") {
		var rec piSessionRecord
		if json.Unmarshal([]byte(line), &rec) != nil {
			continue
		}

		// Assistant messages carry toolCall blocks — subagent spawns and status polls.
		if rec.Message.Role == "assistant" {
			var blocks []piToolCallBlock
			if json.Unmarshal(rec.Message.Content, &blocks) != nil {
				continue
			}
			for _, b := range blocks {
				if b.Type != "toolCall" || b.Name != "subagent" {
					continue
				}
				var args piToolCallArgs
				if json.Unmarshal(b.Arguments, &args) != nil {
					continue
				}
				// A spawn carries a `task`; a status poll carries `action: "status"`.
				if args.Task != "" {
					handle := piHandleFromTask(args.Task)
					if handle == "" {
						continue
					}
					spawn := piSubagentSpawn{
						toolCallID: b.ID,
						handle:     handle,
						stage:      rejectionStageOfHandle(handle),
					}
					spawnByToolCallID[b.ID] = spawn
					pendingSpawn = &spawn
					routes = append(routes, rejectionRoute{
						index: i, event: routeSpawn, stage: spawn.stage, target: spawn.handle,
					})
				} else if args.Action == "status" && args.ID != "" {
					// Status poll — the run id may correlate to a spawn; if so, arm the
					// completion lookup for the toolResult that follows.
					if spawn, ok := spawnByRunID[args.ID]; ok {
						spawnByToolCallID[b.ID] = spawn
						pendingSpawn = &spawn
					}
				}
			}
			continue
		}

		// toolResult messages carry the result of a toolCall. For a spawn, extract the
		// run id; for a status poll, check for completion.
		if rec.Message.Role != "toolResult" || rec.Message.ToolCallID == "" {
			continue
		}
		text := piTextContent(rec.Message.Content)
		spawn, isSpawn := spawnByToolCallID[rec.Message.ToolCallID]
		if !isSpawn {
			continue
		}
		// A spawn result carries the async run id; record it for status correlation.
		if m := piAsyncRunID.FindStringSubmatch(text); m != nil {
			spawnByRunID[m[1]] = spawn
			continue
		}
		// A status result showing completion closes the pending spawn.
		if piStatusCompleted.MatchString(text) {
			routes = append(routes, rejectionRoute{
				index: i, event: routeDone, stage: spawn.stage, target: spawn.handle,
			})
			// Clear the pending spawn so a later status poll for the same run id does
			// not double-count completion.
			delete(spawnByToolCallID, rec.Message.ToolCallID)
			pendingSpawn = nil
		}
	}
	// _ = pendingSpawn keeps the fallback variable referenced for future sequential
	// correlation without dead-code warnings; the run-id correlation above is the
	// primary path.
	_ = pendingSpawn
	return routes, piRejectionBranch
}

// piHandleFromTask extracts the worker handle from the dispatch file path in a
// subagent spawn task. Returns "" when the task carries no dispatch file.
func piHandleFromTask(task string) string {
	m := piDispatchFileInTask.FindStringSubmatch(task)
	if m == nil {
		return ""
	}
	return m[1]
}
