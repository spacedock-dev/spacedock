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
	dispatches, err := piNativeCompletions(stream, artifactDirs...)
	if err != nil {
		return -1, err
	}
	completed, latest := -1, -1
	for _, d := range dispatches {
		if stageToken(d.args.Task, stage) && d.index > latest {
			latest, completed = d.index, d.completed
		}
	}
	return completed, nil
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
