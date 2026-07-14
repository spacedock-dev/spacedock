// ABOUTME: Grades first-officer deferred-core reads from supported-host events.
// ABOUTME: It observes canonical paths and event order; it models no shell/runtime lifecycle.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const (
	foSharedCore = "shared"
	foWriteCore  = "write"
	foMergeCore  = "merge"
	foAdapter    = "adapter"
)

var foCoreSuffixes = map[string]string{
	foSharedCore: "/skills/first-officer/references/first-officer-shared-core.md",
	foWriteCore:  "/skills/first-officer/references/fo-write-core.md",
	foMergeCore:  "/skills/first-officer/references/fo-merge-core.md",
}

type foLoadTrace struct {
	completed  map[string][]int
	attempted  map[string][]int
	actions    []int
	violations []string
}

func newFOLoadTrace() foLoadTrace {
	return foLoadTrace{completed: map[string][]int{}, attempted: map[string][]int{}}
}

// observeLoad records only fields the hosts already emit: a canonical instruction
// path and whether its host event completed successfully. Counting basename and
// suffix occurrences rejects alternate-path retries without interpreting a shell.
func (tr *foLoadTrace) observeLoad(input, host string, line int, completed bool) []string {
	targets := make(map[string]string, len(foCoreSuffixes)+1)
	for core, suffix := range foCoreSuffixes {
		targets[core] = suffix
	}
	targets[foAdapter] = "/skills/first-officer/references/" + host + "-first-officer-runtime.md"

	var found []string
	for core, suffix := range targets {
		base := strings.TrimSuffix(suffix[strings.LastIndex(suffix, "/")+1:], ".md")
		mentions := strings.Count(input, base)
		if mentions == 0 {
			continue
		}
		if strings.Count(input, suffix) != mentions {
			tr.violations = append(tr.violations, fmt.Sprintf("line %d names %s outside its canonical first-officer reference path", line, base))
			continue
		}
		found = append(found, core)
		tr.attempted[core] = append(tr.attempted[core], line)
		if completed {
			tr.completed[core] = append(tr.completed[core], line)
		}
	}
	return found
}

func inputHuntsDeferredCore(input string) bool {
	lower := strings.ToLower(input)
	if !strings.Contains(lower, "fo-write-core") && !strings.Contains(lower, "fo-merge-core") {
		return false
	}
	for _, marker := range []string{"find ", "rg ", "grep ", "glob", "locate ", "ls -r"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func claudeFOLoadTrace(stream, host string, action func(name, command string) bool) foLoadTrace {
	trace := newFOLoadTrace()
	pending := map[string][]string{}
	for lineNo, line := range strings.Split(stream, "\n") {
		var entry streamEntry
		if json.Unmarshal([]byte(line), &entry) != nil {
			continue
		}
		for _, block := range entry.toolUseBlocks() {
			input := strings.Join([]string{block.Input.FilePath, block.Input.Command, block.Input.Pattern, block.Input.Path}, " ")
			if inputHuntsDeferredCore(input) {
				trace.violations = append(trace.violations, fmt.Sprintf("line %d hunts for a deferred core", lineNo+1))
			}
			if strings.Contains(block.Input.Skill, "fo-write-core") || strings.Contains(block.Input.Skill, "fo-merge-core") {
				trace.violations = append(trace.violations, fmt.Sprintf("line %d invokes a deferred core through a wrapper skill", lineNo+1))
			}
			if cores := trace.observeLoad(input, host, lineNo+1, false); len(cores) > 0 {
				pending[block.ID] = cores
			}
			if action != nil && action(block.Name, block.Input.Command) {
				trace.actions = append(trace.actions, lineNo+1)
			}
		}
		for _, block := range entry.resultBlocks() {
			if !block.IsError {
				for _, core := range pending[block.ToolUseID] {
					trace.completed[core] = append(trace.completed[core], lineNo+1)
				}
			}
			delete(pending, block.ToolUseID)
		}
	}
	return trace
}

type codexFOTraceEntry struct {
	Type string `json:"type"`
	Item struct {
		Type     string `json:"type"`
		Command  string `json:"command"`
		Status   string `json:"status"`
		ExitCode *int   `json:"exit_code"`
	} `json:"item"`
}

func codexFOLoadTrace(jsonl, host string, action func(string) bool) foLoadTrace {
	trace := newFOLoadTrace()
	for lineNo, line := range strings.Split(jsonl, "\n") {
		var event codexFOTraceEntry
		if json.Unmarshal([]byte(line), &event) != nil || event.Item.Type != "command_execution" {
			continue
		}
		if event.Type == "item.started" {
			if inputHuntsDeferredCore(event.Item.Command) {
				trace.violations = append(trace.violations, fmt.Sprintf("line %d hunts for a deferred core", lineNo+1))
			}
			trace.observeLoad(event.Item.Command, host, lineNo+1, false)
			if action != nil && action(event.Item.Command) {
				trace.actions = append(trace.actions, lineNo+1)
			}
			continue
		}
		if event.Type == "item.completed" && event.Item.Status == "completed" && event.Item.ExitCode != nil && *event.Item.ExitCode == 0 {
			trace.observeLoad(event.Item.Command, host, lineNo+1, true)
		}
	}
	return trace
}

func traceViolation(trace foLoadTrace) error {
	if len(trace.violations) == 0 {
		return nil
	}
	return fmt.Errorf("non-canonical deferred-core discovery: %s", strings.Join(trace.violations, "; "))
}

func assertFOGateLoadBoundary(trace foLoadTrace) error {
	if err := traceViolation(trace); err != nil {
		return err
	}
	for _, core := range []string{foSharedCore, foAdapter} {
		if len(trace.completed[core]) == 0 {
			return fmt.Errorf("gate journey never completed the boot %s read", core)
		}
	}
	for _, core := range []string{foWriteCore, foMergeCore} {
		if len(trace.attempted[core]) > 0 {
			return fmt.Errorf("mutation-free gate journey attempted deferred %s-core read at line(s) %v", core, trace.attempted[core])
		}
	}
	return nil
}

func assertFOFilingLoadBoundary(trace foLoadTrace) error {
	if err := traceViolation(trace); err != nil {
		return err
	}
	if len(trace.actions) == 0 {
		return fmt.Errorf("filing journey emitted no supported atomic-create action")
	}
	if len(trace.completed[foWriteCore]) == 0 || trace.completed[foWriteCore][0] >= trace.actions[0] {
		return fmt.Errorf("write core did not complete in its own host event before filing (write=%v action=%v)", trace.completed[foWriteCore], trace.actions)
	}
	if len(trace.attempted[foMergeCore]) > 0 {
		return fmt.Errorf("filing journey attempted merge-core read at line(s) %v", trace.attempted[foMergeCore])
	}
	return nil
}

func assertFOTerminalLoadBoundary(trace foLoadTrace) error {
	if err := traceViolation(trace); err != nil {
		return err
	}
	if len(trace.actions) == 0 || len(trace.completed[foWriteCore]) == 0 || len(trace.completed[foMergeCore]) == 0 {
		return fmt.Errorf("terminal journey lacks action or deferred read (write=%v merge=%v action=%v)", trace.completed[foWriteCore], trace.completed[foMergeCore], trace.actions)
	}
	w, m, action := trace.completed[foWriteCore][0], trace.completed[foMergeCore][0], trace.actions[0]
	if !(w < m && m < action) {
		return fmt.Errorf("terminal load order must be write then merge then action in separate host events (write=%d merge=%d action=%d)", w, m, action)
	}
	return nil
}

func claudeFilingAction(slug string) func(string, string) bool {
	return func(name, command string) bool { return name == "Bash" && commandFilesViaNew(command, slug) }
}

func codexFilingAction(slug string) func(string) bool {
	return func(command string) bool { return commandFilesViaNew(command, slug) }
}

func commandSetsTerminalStatus(command, slug string) bool {
	return strings.Contains(command, slug) && strings.Contains(command, "status") && strings.Contains(command, "--set") && strings.Contains(command, "status=done")
}

func claudeTerminalAction(slug string) func(string, string) bool {
	return func(name, command string) bool { return name == "Bash" && commandSetsTerminalStatus(command, slug) }
}

func codexTerminalAction(slug string) func(string) bool {
	return func(command string) bool { return commandSetsTerminalStatus(command, slug) }
}

func codexBoundaryEvent(kind, command, status string, exitCode int) string {
	body, _ := json.Marshal(map[string]any{"type": kind, "item": map[string]any{"type": "command_execution", "command": command, "status": status, "exit_code": exitCode}})
	return string(body)
}

func claudeBoundaryEvent(id, kind, name string, input map[string]string) string {
	block := map[string]any{"type": kind}
	if kind == "tool_use" {
		block["id"], block["name"], block["input"] = id, name, input
	} else {
		block["tool_use_id"], block["is_error"] = id, false
	}
	role := "assistant"
	if kind == "tool_result" {
		role = "user"
	}
	body, _ := json.Marshal(map[string]any{"type": role, "message": map[string]any{"content": []any{block}}})
	return string(body)
}

func TestCodexFODeferredLoadTraceBoundaries(t *testing.T) {
	root := "/cache/skills/first-officer/references/"
	boot := codexBoundaryEvent("item.completed", "cat "+root+"first-officer-shared-core.md "+root+"codex-first-officer-runtime.md", "completed", 0)
	write := codexBoundaryEvent("item.completed", "cat "+root+"fo-write-core.md", "completed", 0)
	merge := codexBoundaryEvent("item.completed", "cat "+root+"fo-merge-core.md", "completed", 0)
	filing := codexBoundaryEvent("item.started", "spacedock new wire-the-thing", "in_progress", 0)
	terminal := codexBoundaryEvent("item.started", "spacedock status --set merge-check status=done", "in_progress", 0)

	if err := assertFOGateLoadBoundary(codexFOLoadTrace(boot, "codex", nil)); err != nil {
		t.Fatalf("cold gate: %v", err)
	}
	if err := assertFOFilingLoadBoundary(codexFOLoadTrace(strings.Join([]string{boot, write, filing}, "\n"), "codex", codexFilingAction("wire-the-thing"))); err != nil {
		t.Fatalf("filing: %v", err)
	}
	if err := assertFOTerminalLoadBoundary(codexFOLoadTrace(strings.Join([]string{boot, write, merge, terminal}, "\n"), "codex", codexTerminalAction("merge-check"))); err != nil {
		t.Fatalf("terminal: %v", err)
	}

	bad := []struct {
		name  string
		trace foLoadTrace
		check func(foLoadTrace) error
	}{
		{"eager gate", codexFOLoadTrace(boot+"\n"+write, "codex", nil), assertFOGateLoadBoundary},
		{"late write", codexFOLoadTrace(strings.Join([]string{boot, filing, write}, "\n"), "codex", codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"reversed terminal", codexFOLoadTrace(strings.Join([]string{boot, merge, write, terminal}, "\n"), "codex", codexTerminalAction("merge-check")), assertFOTerminalLoadBoundary},
		{"filesystem hunt", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEvent("item.started", `find / -iname "fo-write-core*"`, "in_progress", 0), write, filing}, "\n"), "codex", codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"alternate path", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEvent("item.started", "cat /tmp/fo-write-core.md", "in_progress", 0), write, filing}, "\n"), "codex", codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
	}
	for _, tc := range bad {
		if err := tc.check(tc.trace); err == nil {
			t.Errorf("%s trace must fail", tc.name)
		}
	}
}

func TestClaudeFODeferredLoadTraceBoundaries(t *testing.T) {
	root := "/cache/skills/first-officer/references/"
	read := func(id, name string) []string {
		return []string{
			claudeBoundaryEvent(id, "tool_use", "Read", map[string]string{"file_path": root + name}),
			claudeBoundaryEvent(id, "tool_result", "", nil),
		}
	}
	boot := append(read("shared", "first-officer-shared-core.md"), read("adapter", "claude-first-officer-runtime.md")...)
	write, merge := read("write", "fo-write-core.md"), read("merge", "fo-merge-core.md")
	terminal := claudeBoundaryEvent("terminal", "tool_use", "Bash", map[string]string{"command": "spacedock status --set merge-check status=done"})
	lines := append(append(append([]string{}, boot...), write...), merge...)
	lines = append(lines, terminal)
	if err := assertFOTerminalLoadBoundary(claudeFOLoadTrace(strings.Join(lines, "\n"), "claude", claudeTerminalAction("merge-check"))); err != nil {
		t.Fatalf("terminal: %v", err)
	}
	wrapper := claudeBoundaryEvent("wrapper", "tool_use", "Skill", map[string]string{"skill": "spacedock:fo-write-core"})
	if err := traceViolation(claudeFOLoadTrace(strings.Join(append(boot, wrapper), "\n"), "claude", nil)); err == nil {
		t.Fatal("wrapper-skill invocation must fail")
	}
}
