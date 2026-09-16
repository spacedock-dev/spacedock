package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// piNativeCompletion credits only a retained child of this parent, at the native
// notification's original index. Report commits and gate state remain the callers'
// independent durable checks. Historical status/wait evidence needs no artifacts.
func piNativeCompletion(stream, stage string, artifactDirs ...string) (int, error) {
	if len(artifactDirs) == 0 || artifactDirs[0] == "" {
		return -1, nil
	}
	var parent, spawn piSessionRecord
	var args piToolCallArgs
	var callID, runID, owner string
	completed := -1
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
				var a piToolCallArgs
				if b.Type == "toolCall" && b.Name == "subagent" && json.Unmarshal(b.Arguments, &a) == nil && stageToken(a.Task, stage) {
					spawn, args, callID, runID, owner = rec, a, b.ID, "", ""
					completed = -1 // A fresh assignment cannot inherit a previous completion.
				}
			}
		}
		if rec.Message.Role == "toolResult" && rec.Message.ToolName == "subagent" && callID != "" && rec.Message.ToolCallID == callID && rec.Message.IsError != nil && !*rec.Message.IsError {
			runID, owner = rec.Message.Details.RunID, rec.Message.Details.Mission.OwnerSessionID
		}
		if rec.Type != "custom_message" || rec.CustomType != "subagent-notify" || runID == "" {
			continue
		}
		// Locator is the sole notification field used: completion prose is not proof.
		locators := []string{}
		for _, text := range strings.Split(rec.Content, "\n") {
			if strings.HasPrefix(text, "Session file: ") {
				locators = append(locators, strings.TrimPrefix(text, "Session file: "))
			}
		}
		if len(locators) != 1 {
			return -1, fmt.Errorf("Pi notification must name exactly one child session")
		}
		childPath, err := piRetainedChild(artifactDirs[0], parent, owner, locators[0])
		if err != nil {
			return -1, err
		}
		raw, err := os.ReadFile(childPath)
		if err != nil {
			return -1, fmt.Errorf("read Pi child: %w", err)
		}
		var session, terminal piSessionRecord
		names, tasks := []string{}, []string{}
		for _, childLine := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
			var child piSessionRecord
			if err := json.Unmarshal([]byte(childLine), &child); err != nil {
				return -1, fmt.Errorf("invalid Pi child: %w", err)
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
		prefix := "subagent-" + args.Agent + "-" + runID + "-"
		sameRun := false
		for _, name := range names {
			sameRun = sameRun || strings.HasPrefix(name, prefix)
		}
		if !sameRun {
			continue
		}
		if len(names) != 1 || names[0] != prefix+"1" {
			return -1, fmt.Errorf("Pi child has conflicting run/epoch identities")
		}
		if args.Agent == "" || args.Context != "fresh" || args.Cwd == "" || session.Cwd != args.Cwd || len(tasks) != 1 || tasks[0] != "Task: "+args.Task || filepath.Base(filepath.Dir(childPath)) != "run-0" {
			return -1, fmt.Errorf("Pi child assignment/epoch does not match dispatch")
		}
		if terminal.Type != "message" || terminal.Message.Role != "assistant" || terminal.Message.StopReason != "stop" {
			return -1, fmt.Errorf("Pi child has no successful terminal stop")
		}
		started, e1 := time.Parse(time.RFC3339Nano, spawn.Timestamp)
		childStarted, e2 := time.Parse(time.RFC3339Nano, session.Timestamp)
		stopped, e3 := time.Parse(time.RFC3339Nano, terminal.Timestamp)
		notified, e4 := time.Parse(time.RFC3339Nano, rec.Timestamp)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil || !started.Before(childStarted) || stopped.Before(childStarted) || stopped.After(notified) {
			return -1, fmt.Errorf("Pi child completion timestamps are out of order")
		}
		if completed >= 0 {
			return -1, fmt.Errorf("duplicate Pi completion notification")
		}
		completed = i
	}
	return completed, nil
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
