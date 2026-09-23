package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// A completion belongs to its exact dispatch call, not the latest worker at a
// stage. Both lifecycle and route observers use the same retained-child verifier.
type piNativeDispatch struct {
	parent, spawn        piSessionRecord
	args                 piToolCallArgs
	callID, runID, owner string
	index, completed     int
	resultSeen           bool
}

func piNativeCompletion(stream, stage string, artifactDirs ...string) (int, error) {
	completed, _, err := piNativeLifecycle(stream, stage, "", artifactDirs...)
	return completed, err
}

// A repaired gate must follow its own worker. Join actual prepare/withdraw
// results by tool-call ID, and retain every earlier attempt's ordering checks.
func piNativeLifecycle(stream, stage, nextSignal string, artifactDirs ...string) (int, int, error) {
	dispatches, err := piNativeCompletions(stream, artifactDirs...)
	if err != nil {
		return -1, -1, err
	}
	completed, latest, count := -1, -1, 0
	workers := map[int]*piNativeDispatch{}
	for _, d := range dispatches {
		if stageToken(d.args.Task, stage) {
			workers[d.index] = d
			count++
			if d.index > latest {
				latest, completed = d.index, d.completed
			}
		}
	}
	prepares := 0
	for _, call := range piBashCommands(stream) {
		if strings.Contains(call.command, "gate prepare") {
			prepares++
		}
	}
	// Withdrawal permits re-preparing the same validated work; it does not
	// itself require another worker. Still check every replacement attempt.
	if count == 0 || count < 2 && prepares < 2 || nextSignal != "gate prepare" {
		return completed, -1, nil
	}
	type gateCall struct {
		index    int
		withdraw bool
	}
	calls := map[string]gateCall{}
	active, prior := "", ""
	boundary := -1
	var worker *piNativeDispatch
	fail := func(reason string) (int, int, error) { return -1, -1, fmt.Errorf("Pi gate chronology: %s", reason) }
	for i, line := range strings.Split(stream, "\n") {
		if d := workers[i]; d != nil {
			if active != "" {
				return fail("repair dispatched before successful gate withdrawal")
			}
			worker = d
		}
		var rec piSessionRecord
		if json.Unmarshal([]byte(line), &rec) != nil {
			continue
		}
		if rec.Message.Role == "assistant" {
			var blocks []piToolCallBlock
			_ = json.Unmarshal(rec.Message.Content, &blocks)
			for _, b := range blocks {
				var args piToolCallArgs
				if b.Type != "toolCall" || b.Name != "bash" || json.Unmarshal(b.Arguments, &args) != nil {
					continue
				}
				prepare, withdraw := strings.Contains(args.Command, "gate prepare"), strings.Contains(args.Command, "gate withdraw")
				if !prepare && !withdraw {
					continue
				}
				if prepare && (worker == nil || worker.completed < 0 || worker.completed >= i) {
					return fail("gate precedes its worker completion")
				}
				if b.ID == "" || calls[b.ID].index != 0 || prepare && withdraw {
					return fail("ambiguous gate call")
				}
				calls[b.ID] = gateCall{index: i, withdraw: withdraw}
			}
		}
		call, ok := calls[rec.Message.ToolCallID]
		if !ok || rec.Message.Role != "toolResult" || rec.Message.ToolName != "bash" {
			continue
		}
		delete(calls, rec.Message.ToolCallID)
		if rec.Message.IsError == nil || *rec.Message.IsError {
			return fail("gate call lacks successful result")
		}
		fields := map[string]string{}
		for _, field := range strings.Fields(piTextContent(rec.Message.Content)) {
			if key, value, ok := strings.Cut(field, "="); ok {
				fields[key] = value
			}
		}
		briefing := fields["briefing"]
		parts := strings.Split(briefing, ":")
		if len(parts) != 5 || parts[0] != "briefing" || parts[2] != stage {
			return fail("gate result lacks matching briefing identity")
		}
		if call.withdraw {
			if active == "" || briefing != active || fields["state"] != "withdrawn" {
				return fail("withdrawal does not match open attempt")
			}
			prior, active = active, ""
		} else {
			if active != "" || briefing == prior || fields["state"] != "open" {
				return fail("replacement gate lacks withdrawn predecessor or new identity")
			}
			if prior != "" && !strings.HasPrefix(prior, strings.Join(parts[:3], ":")+":") {
				return fail("replacement gate belongs to another entity")
			}
			active, boundary = briefing, call.index
		}
	}
	if len(calls) != 0 || active == "" || boundary < latest || completed >= boundary {
		return fail("current gate is incomplete or precedes repair")
	}
	return completed, boundary, nil
}

// Native sync results and async notices are observed at their original parent
// indexes. Git/report/gate checks remain independent; historical forms need no
// artifact root. No synthetic notifications or child calls enter the stream.
func piNativeCompletions(stream string, artifactDirs ...string) (map[string]*piNativeDispatch, error) {
	dispatches := map[string]*piNativeDispatch{}
	if len(artifactDirs) == 0 || artifactDirs[0] == "" {
		return dispatches, nil
	}
	var parent piSessionRecord
	runs := map[string]string{}
	for i, line := range strings.Split(stream, "\n") {
		var rec piSessionRecord
		if json.Unmarshal([]byte(line), &rec) != nil {
			continue
		}
		if rec.Type == "session" {
			parent = rec
		}
		if rec.Message.Role == "assistant" {
			var blocks []piToolCallBlock
			_ = json.Unmarshal(rec.Message.Content, &blocks)
			for _, b := range blocks {
				var args piToolCallArgs
				if b.Type != "toolCall" || b.Name != "subagent" || json.Unmarshal(b.Arguments, &args) != nil || args.Task == "" {
					continue
				}
				if b.ID == "" || dispatches[b.ID] != nil {
					return nil, fmt.Errorf("missing or repeated Pi dispatch identity")
				}
				dispatches[b.ID] = &piNativeDispatch{parent: parent, spawn: rec, args: args, callID: b.ID, index: i, completed: -1}
			}
		}
		if rec.Message.Role == "toolResult" && rec.Message.ToolName == "subagent" {
			d := dispatches[rec.Message.ToolCallID]
			if d == nil {
				continue
			}
			if d.resultSeen {
				return nil, fmt.Errorf("multiple Pi results for dispatch %s", d.callID)
			}
			d.resultSeen = true
			details := rec.Message.Details
			sync := d.args.Async != nil && !*d.args.Async
			if !sync && details.RunID == "" {
				continue
			} // Historical status/text path.
			if rec.Message.IsError == nil || *rec.Message.IsError || details.RunID == "" {
				return nil, fmt.Errorf("Pi dispatch result lacks explicit success/run identity")
			}
			spawned, e1 := time.Parse(time.RFC3339Nano, d.spawn.Timestamp)
			returned, e2 := time.Parse(time.RFC3339Nano, rec.Timestamp)
			if e1 != nil || e2 != nil || returned.Before(spawned) {
				return nil, fmt.Errorf("Pi dispatch result precedes its spawn")
			}
			if prior := runs[details.RunID]; prior != "" {
				return nil, fmt.Errorf("Pi run reused by multiple dispatches")
			}
			runs[details.RunID] = d.callID
			d.runID, d.owner = details.RunID, details.Mission.OwnerSessionID
			if !sync {
				continue
			}
			if len(details.Results) != 1 {
				return nil, fmt.Errorf("Pi synchronous result must contain exactly one child")
			}
			result := details.Results[0]
			if result.Agent != d.args.Agent || result.ExitCode == nil || *result.ExitCode != 0 || result.OutputState != "present" || result.SessionFile == "" {
				return nil, fmt.Errorf("Pi synchronous child result lacks matching agent/explicit exit zero/output")
			}
			matched, err := piVerifyNativeChild(artifactDirs[0], d, rec, result.SessionFile)
			if err != nil {
				return nil, err
			}
			if !matched {
				return nil, fmt.Errorf("Pi synchronous child run does not match dispatch")
			}
			d.completed = i
		}
		if rec.Type != "custom_message" || rec.CustomType != "subagent-notify" {
			continue
		}
		locators := []string{}
		for _, text := range strings.Split(rec.Content, "\n") {
			if strings.HasPrefix(text, "Session file: ") {
				locators = append(locators, strings.TrimPrefix(text, "Session file: "))
			}
		}
		for _, d := range dispatches {
			if d.runID == "" {
				continue
			}
			if len(locators) != 1 {
				return nil, fmt.Errorf("Pi notification must name exactly one child session")
			}
			matched, err := piVerifyNativeChild(artifactDirs[0], d, rec, locators[0])
			if err != nil {
				return nil, err
			}
			if !matched {
				continue
			}
			if d.completed >= 0 {
				return nil, fmt.Errorf("duplicate Pi completion notification")
			}
			d.completed = i
		}
	}
	return dispatches, nil
}

func piVerifyNativeChild(artifactDir string, dispatch *piNativeDispatch, rec piSessionRecord, locator string) (bool, error) {
	childPath, err := piRetainedChild(artifactDir, dispatch.parent, dispatch.owner, locator)
	if err != nil {
		return false, err
	}
	raw, err := os.ReadFile(childPath)
	if err != nil {
		return false, fmt.Errorf("read Pi child: %w", err)
	}
	var session, terminal piSessionRecord
	names, tasks := []string{}, []string{}
	for _, childLine := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		var child piSessionRecord
		if err := json.Unmarshal([]byte(childLine), &child); err != nil {
			return false, fmt.Errorf("invalid Pi child: %w", err)
		}
		if child.Type == "session" {
			session = child
		}
		if child.Type == "session_info" {
			names = append(names, child.Name)
		}
		if child.Message.Role == "user" {
			tasks = append(tasks, piTextContent(child.Message.Content))
		}
		terminal = child
	}
	// Other workers also notify this parent; only the exact run can complete this
	// dispatch. No labels, directory UUID guesses or abbreviated commit matching.
	prefix := "subagent-" + dispatch.args.Agent + "-" + dispatch.runID + "-"
	sameRun := false
	for _, name := range names {
		sameRun = sameRun || strings.HasPrefix(name, prefix)
	}
	if !sameRun {
		return false, nil
	}
	if len(names) != 1 || names[0] != prefix+"1" {
		return false, fmt.Errorf("Pi child has conflicting run/epoch identities")
	}
	if dispatch.args.Agent == "" || dispatch.args.Context != "fresh" || dispatch.args.Cwd == "" || session.Cwd != dispatch.args.Cwd || len(tasks) != 1 || tasks[0] != "Task: "+dispatch.args.Task || filepath.Base(filepath.Dir(childPath)) != "run-0" {
		return false, fmt.Errorf("Pi child assignment/epoch does not match dispatch")
	}
	if terminal.Type != "message" || terminal.Message.Role != "assistant" || terminal.Message.StopReason != "stop" {
		return false, fmt.Errorf("Pi child has no successful terminal stop")
	}
	started, e1 := time.Parse(time.RFC3339Nano, dispatch.spawn.Timestamp)
	childStarted, e2 := time.Parse(time.RFC3339Nano, session.Timestamp)
	stopped, e3 := time.Parse(time.RFC3339Nano, terminal.Timestamp)
	notified, e4 := time.Parse(time.RFC3339Nano, rec.Timestamp)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || !started.Before(childStarted) || stopped.Before(childStarted) || stopped.After(notified) {
		return false, fmt.Errorf("Pi child completion timestamps are out of order")
	}

	return true, nil
}

// Remap only the structured owner-session prefix. The locator must be inside
// that parent's retained tree, never a sibling session or an external/symlink file.
func piRetainedChild(artifactDir string, parent piSessionRecord, owner, locator string) (string, error) {
	if parent.ID == "" || !filepath.IsAbs(owner) || filepath.Clean(owner) != owner || !strings.HasSuffix(owner, "_"+parent.ID+".jsonl") || !filepath.IsAbs(locator) || filepath.Clean(locator) != locator {
		return "", fmt.Errorf("invalid Pi parent/child locator")
	}
	sourceTree := strings.TrimSuffix(owner, ".jsonl")
	rel, err := filepath.Rel(sourceTree, locator)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("Pi child locator escapes parent session tree")
	}
	sessions, err := filepath.EvalSymlinks(filepath.Join(artifactDir, "sessions"))
	if err != nil {
		return "", err
	}
	// The retained root must be the parent whose structured run result supplied
	// ownerSessionId, not merely a file with a matching child basename.
	raw, err := os.ReadFile(filepath.Join(sessions, filepath.Base(owner)))
	var retained piSessionRecord
	if err != nil || json.Unmarshal([]byte(strings.SplitN(string(raw), "\n", 2)[0]), &retained) != nil || retained.Type != "session" || retained.ID != parent.ID || retained.Cwd != parent.Cwd {
		return "", fmt.Errorf("Pi retained parent identity is missing or mismatched")
	}
	path := filepath.Join(sessions, filepath.Base(sourceTree), rel)
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil || resolved != path {
		return "", fmt.Errorf("Pi child missing or symlinked: %s", path)
	}
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("Pi child is not a regular file: %s", path)
	}
	return path, nil
}
