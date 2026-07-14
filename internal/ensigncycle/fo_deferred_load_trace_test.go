// ABOUTME: Grades first-officer deferred-core reads from supported-host events.
// ABOUTME: It observes canonical paths and event order; it models no shell/runtime lifecycle.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	foSharedCore = "shared"
	foWriteCore  = "write"
	foMergeCore  = "merge"
	foAdapter    = "adapter"
)

var foCoreFiles = map[string]string{
	foSharedCore: "first-officer-shared-core.md",
	foWriteCore:  "fo-write-core.md",
	foMergeCore:  "fo-merge-core.md",
}

// foLoadSpec binds the observer to the loader-supplied base for one installed
// first-officer skill. Codex also uses the canonical file bodies to distinguish a
// full command result from path-only or partial output.
type foLoadSpec struct {
	paths  map[string]string
	bodies map[string]string
}

func readFOLoadSpec(firstOfficerBase, host string) (foLoadSpec, error) {
	firstOfficerBase = canonicalFilesystemPath(firstOfficerBase)
	files := make(map[string]string, len(foCoreFiles)+1)
	for core, name := range foCoreFiles {
		files[core] = name
	}
	files[foAdapter] = host + "-first-officer-runtime.md"

	spec := foLoadSpec{paths: map[string]string{}, bodies: map[string]string{}}
	for core, name := range files {
		path := filepath.Join(firstOfficerBase, "references", name)
		body, err := os.ReadFile(path)
		if err != nil {
			return foLoadSpec{}, fmt.Errorf("read canonical %s core %s: %w", core, path, err)
		}
		spec.paths[core] = path
		spec.bodies[core] = string(body)
	}
	return spec, nil
}

func mustFOLoadSpec(t testing.TB, firstOfficerBase, host string) foLoadSpec {
	t.Helper()
	spec, err := readFOLoadSpec(firstOfficerBase, host)
	if err != nil {
		t.Fatal(err)
	}
	return spec
}

func fixtureFOLoadSpec(t testing.TB, host string) foLoadSpec {
	t.Helper()
	return fixtureFOLoadSpecAt(t, host, filepath.Join(t.TempDir(), "skills", "first-officer"))
}

func fixtureFOLoadSpecAt(t testing.TB, host, base string) foLoadSpec {
	t.Helper()
	refs := filepath.Join(base, "references")
	if err := os.MkdirAll(refs, 0o755); err != nil {
		t.Fatal(err)
	}
	files := make(map[string]string, len(foCoreFiles)+1)
	for core, name := range foCoreFiles {
		files[core] = name
	}
	files[foAdapter] = host + "-first-officer-runtime.md"
	for core, name := range files {
		body := []byte("# " + host + " " + core + " canonical body\n\nfull-file sentinel for " + core + "\n")
		if err := os.WriteFile(filepath.Join(refs, name), body, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return mustFOLoadSpec(t, base, host)
}

type foLoadTrace struct {
	completed  map[string][]int
	attempted  map[string][]int
	actions    []int
	mutations  []int
	violations []string
}

func newFOLoadTrace() foLoadTrace {
	return foLoadTrace{completed: map[string][]int{}, attempted: map[string][]int{}}
}

func isPathByte(b byte) bool {
	return b == '/' || b == '.' || b == '_' || b == '-' || b == '~' ||
		b >= '0' && b <= '9' || b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func canonicalFilesystemPath(path string) string {
	path = filepath.Clean(path)
	if absolute, err := filepath.Abs(path); err == nil {
		path = absolute
	}
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	return filepath.Clean(path)
}

func sameCanonicalFilesystemPath(left, right string) bool {
	return canonicalFilesystemPath(left) == canonicalFilesystemPath(right)
}

// corePathMentions returns every bounded path-shaped token that names base.
// This is a lexical scan over one supported host field, not shell parsing: it
// exists so one exact path cannot mask a second alternate-root occurrence.
func corePathMentions(input, base string) []string {
	var mentions []string
	for from := 0; from <= len(input)-len(base); {
		i := strings.Index(input[from:], base)
		if i < 0 {
			break
		}
		i += from
		start, end := i, i+len(base)
		for start > 0 && isPathByte(input[start-1]) {
			start--
		}
		for end < len(input) && isPathByte(input[end]) {
			end++
		}
		mentions = append(mentions, input[start:end])
		from = i + len(base)
	}
	return mentions
}

// mentionedCores validates exact installed paths without interpreting command
// grammar. Every basename occurrence is graded; a mention at any other root is
// a hard evidence violation even when the same field also has an exact read.
func (tr *foLoadTrace) mentionedCores(input string, spec foLoadSpec, line int) []string {
	var found []string
	for core, path := range spec.paths {
		base := strings.TrimSuffix(filepath.Base(path), ".md")
		mentions := corePathMentions(input, base)
		if len(mentions) == 0 {
			continue
		}
		exact := false
		for _, mention := range mentions {
			if sameCanonicalFilesystemPath(mention, path) {
				exact = true
				continue
			}
			tr.violations = append(tr.violations, fmt.Sprintf("line %d names %s at noncanonical path %s (loader path %s)", line, base, mention, path))
		}
		if exact {
			found = append(found, core)
		}
	}
	return found
}

func (tr *foLoadTrace) attempt(cores []string, line int) {
	for _, core := range cores {
		tr.attempted[core] = append(tr.attempted[core], line)
	}
}

func (tr *foLoadTrace) complete(cores []string, line int) {
	for _, core := range cores {
		tr.completed[core] = append(tr.completed[core], line)
	}
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

// commandStartsFOMutation recognizes only the repository's supported mutation
// surfaces. It is deliberately a fixed marker check over an already-emitted host
// field, not a shell parser or an attempt to model arbitrary commands.
func commandStartsFOMutation(command string) bool {
	lower := strings.ToLower(command)
	padded := " " + strings.Join(strings.Fields(lower), " ") + " "
	if strings.Contains(lower, "spacedock") {
		if strings.Contains(lower, "--set") && strings.Contains(lower, "status") {
			return true
		}
		if strings.Contains(padded, " new ") {
			return true
		}
		if strings.Contains(lower, "dispatch build") && strings.Contains(lower, "--advance") {
			return true
		}
		if strings.Contains(lower, "state ready") {
			return true
		}
		for _, marker := range []string{" state commit ", " merge guard "} {
			if strings.Contains(padded, marker) {
				return true
			}
		}
		if strings.Contains(padded, " status ") && strings.Contains(padded, " --archive ") {
			return true
		}
	}
	for _, marker := range []string{"apply_patch", "applypatch", "git commit", "git mv", "git rm"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return strings.Contains(lower, "_archive") && strings.Contains(padded, " mv ")
}

func claudeToolStartsFOMutation(block *streamBlock) bool {
	switch block.Name {
	case "Edit", "Write", "MultiEdit", "NotebookEdit":
		return true
	case "Bash":
		return commandStartsFOMutation(block.Input.Command)
	default:
		return false
	}
}

func claudeFOLoadTrace(stream string, spec foLoadSpec, action func(name, command string) bool) foLoadTrace {
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
			cores := trace.mentionedCores(input, spec, lineNo+1)
			if len(cores) > 0 {
				if block.Name != "Read" {
					trace.violations = append(trace.violations, fmt.Sprintf("line %d mentions a first-officer core outside Claude Read", lineNo+1))
				} else {
					var exact []string
					for _, core := range cores {
						if sameCanonicalFilesystemPath(block.Input.FilePath, spec.paths[core]) {
							exact = append(exact, core)
						}
					}
					trace.attempt(exact, lineNo+1)
					if block.Input.Offset != nil || block.Input.Limit != nil {
						trace.violations = append(trace.violations, fmt.Sprintf("line %d uses a partial Claude Read for %v", lineNo+1, exact))
					} else if block.ID == "" {
						trace.violations = append(trace.violations, fmt.Sprintf("line %d has an uncorrelatable Claude Read", lineNo+1))
					} else {
						pending[block.ID] = exact
					}
				}
			}
			if claudeToolStartsFOMutation(block) {
				trace.mutations = append(trace.mutations, lineNo+1)
			}
			if action != nil && action(block.Name, block.Input.Command) {
				trace.actions = append(trace.actions, lineNo+1)
			}
		}
		for _, block := range entry.resultBlocks() {
			if !block.IsError {
				trace.complete(pending[block.ToolUseID], lineNo+1)
			}
			delete(pending, block.ToolUseID)
		}
	}
	return trace
}

type codexFOTraceEntry struct {
	Type string `json:"type"`
	Item struct {
		Type             string `json:"type"`
		Command          string `json:"command"`
		Status           string `json:"status"`
		AggregatedOutput string `json:"aggregated_output"`
		ExitCode         *int   `json:"exit_code"`
	} `json:"item"`
}

func codexFOLoadTrace(jsonl string, spec foLoadSpec, action func(string) bool) foLoadTrace {
	trace := newFOLoadTrace()
	for lineNo, line := range strings.Split(jsonl, "\n") {
		var event codexFOTraceEntry
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		if event.Item.Type == "file_change" && (event.Type == "item.started" || event.Type == "item.completed") {
			trace.mutations = append(trace.mutations, lineNo+1)
			continue
		}
		if event.Item.Type != "command_execution" || (event.Type != "item.started" && event.Type != "item.completed") {
			continue
		}
		if inputHuntsDeferredCore(event.Item.Command) {
			trace.violations = append(trace.violations, fmt.Sprintf("line %d hunts for a deferred core", lineNo+1))
		}
		cores := trace.mentionedCores(event.Item.Command, spec, lineNo+1)
		trace.attempt(cores, lineNo+1)
		if commandStartsFOMutation(event.Item.Command) {
			trace.mutations = append(trace.mutations, lineNo+1)
		}
		if action != nil && action(event.Item.Command) {
			trace.actions = append(trace.actions, lineNo+1)
		}
		if event.Type != "item.completed" || event.Item.Status != "completed" || event.Item.ExitCode == nil || *event.Item.ExitCode != 0 {
			continue
		}
		for _, core := range cores {
			if strings.Contains(event.Item.AggregatedOutput, spec.bodies[core]) {
				trace.complete([]string{core}, lineNo+1)
				continue
			}
			trace.violations = append(trace.violations, fmt.Sprintf("line %d did not return the full canonical %s body", lineNo+1, core))
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

func assertFOWriteBeforeFirstMutation(trace foLoadTrace) error {
	if err := traceViolation(trace); err != nil {
		return err
	}
	if len(trace.mutations) == 0 {
		return fmt.Errorf("mutating journey emitted no supported FO-authored mutation event")
	}
	if len(trace.completed[foWriteCore]) == 0 || trace.completed[foWriteCore][0] >= trace.mutations[0] {
		return fmt.Errorf("write core did not complete before the first FO-authored mutation (write=%v mutation=%v)", trace.completed[foWriteCore], trace.mutations)
	}
	return nil
}

func assertFOWriteBeforeObservedMutation(trace foLoadTrace) error {
	if len(trace.mutations) == 0 {
		return traceViolation(trace)
	}
	return assertFOWriteBeforeFirstMutation(trace)
}

func assertFOFilingLoadBoundary(trace foLoadTrace) error {
	if err := assertFOWriteBeforeFirstMutation(trace); err != nil {
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
	if err := assertFOWriteBeforeFirstMutation(trace); err != nil {
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
	return codexBoundaryEventWithOutput(kind, command, status, exitCode, "")
}

func codexBoundaryEventWithOutput(kind, command, status string, exitCode int, output string) string {
	body, _ := json.Marshal(map[string]any{"type": kind, "item": map[string]any{"type": "command_execution", "command": command, "status": status, "exit_code": exitCode, "aggregated_output": output}})
	return string(body)
}

func claudeBoundaryEvent(id, kind, name string, input map[string]any) string {
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
	spec := fixtureFOLoadSpec(t, "codex")
	boot := codexBoundaryEventWithOutput("item.completed", "cat "+spec.paths[foSharedCore]+" "+spec.paths[foAdapter], "completed", 0, spec.bodies[foSharedCore]+spec.bodies[foAdapter])
	write := codexBoundaryEventWithOutput("item.completed", "cat "+spec.paths[foWriteCore], "completed", 0, spec.bodies[foWriteCore])
	merge := codexBoundaryEventWithOutput("item.completed", "cat "+spec.paths[foMergeCore], "completed", 0, spec.bodies[foMergeCore])
	filing := codexBoundaryEvent("item.started", "spacedock new wire-the-thing", "in_progress", 0)
	terminal := codexBoundaryEvent("item.started", "spacedock status --set merge-check status=done", "in_progress", 0)

	if err := assertFOGateLoadBoundary(codexFOLoadTrace(boot, spec, nil)); err != nil {
		t.Fatalf("cold gate: %v", err)
	}
	if err := assertFOFilingLoadBoundary(codexFOLoadTrace(strings.Join([]string{boot, write, filing}, "\n"), spec, codexFilingAction("wire-the-thing"))); err != nil {
		t.Fatalf("filing: %v", err)
	}
	if err := assertFOTerminalLoadBoundary(codexFOLoadTrace(strings.Join([]string{boot, write, merge, terminal}, "\n"), spec, codexTerminalAction("merge-check"))); err != nil {
		t.Fatalf("terminal: %v", err)
	}

	bad := []struct {
		name  string
		trace foLoadTrace
		check func(foLoadTrace) error
	}{
		{"eager gate", codexFOLoadTrace(boot+"\n"+write, spec, nil), assertFOGateLoadBoundary},
		{"late write", codexFOLoadTrace(strings.Join([]string{boot, filing, write}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"reversed terminal", codexFOLoadTrace(strings.Join([]string{boot, merge, write, terminal}, "\n"), spec, codexTerminalAction("merge-check")), assertFOTerminalLoadBoundary},
		{"filesystem hunt", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEvent("item.started", `find / -iname "fo-write-core*"`, "in_progress", 0), write, filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"alternate path", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEvent("item.started", "cat /tmp/fo-write-core.md", "in_progress", 0), write, filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"path only", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEventWithOutput("item.completed", "printf "+spec.paths[foWriteCore], "completed", 0, spec.paths[foWriteCore]), filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"partial read", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEventWithOutput("item.completed", "head -n 1 "+spec.paths[foWriteCore], "completed", 0, strings.SplitN(spec.bodies[foWriteCore], "\n", 2)[0]+"\n"), filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"same suffix alternate root", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEventWithOutput("item.completed", "cat /alternate/skills/first-officer/references/fo-write-core.md", "completed", 0, spec.bodies[foWriteCore]), filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"exact plus same suffix alternate root", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEventWithOutput("item.completed", "cat "+spec.paths[foWriteCore]+" /alternate/skills/first-officer/references/fo-write-core.md", "completed", 0, spec.bodies[foWriteCore]), filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"early ordinary mutation", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEvent("item.started", "spacedock status --set ordinary status=review", "in_progress", 0), write, filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"early state commit", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEvent("item.started", "spacedock state commit ordinary", "in_progress", 0), write, filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"early status archive", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEvent("item.started", "spacedock status --archive ordinary", "in_progress", 0), write, filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
		{"early merge guard", codexFOLoadTrace(strings.Join([]string{boot, codexBoundaryEvent("item.started", "spacedock merge guard ordinary --verdict passed", "in_progress", 0), write, filing}, "\n"), spec, codexFilingAction("wire-the-thing")), assertFOFilingLoadBoundary},
	}
	for _, tc := range bad {
		if err := tc.check(tc.trace); err == nil {
			t.Errorf("%s trace must fail", tc.name)
		}
	}
}

func TestCodexFODeferredLoadTraceAcceptsTmpCanonicalAlias(t *testing.T) {
	root, err := os.MkdirTemp("/tmp", "spacedock-fo-load-alias-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(root) })
	base := filepath.Join(root, "skills", "first-officer")
	spec := fixtureFOLoadSpecAt(t, "codex", base)
	lexical := filepath.Join(base, "references", foCoreFiles[foWriteCore])
	canonical, err := filepath.EvalSymlinks(lexical)
	if err != nil {
		t.Fatal(err)
	}
	filing := codexBoundaryEvent("item.started", "spacedock new wire-the-thing", "in_progress", 0)
	for name, path := range map[string]string{"lexical": lexical, "canonical": canonical} {
		t.Run(name, func(t *testing.T) {
			write := codexBoundaryEventWithOutput("item.completed", "cat "+path, "completed", 0, spec.bodies[foWriteCore])
			if err := assertFOFilingLoadBoundary(codexFOLoadTrace(write+"\n"+filing, spec, codexFilingAction("wire-the-thing"))); err != nil {
				t.Fatalf("loader path %q and observed path %q must identify the same read: %v", spec.paths[foWriteCore], path, err)
			}
		})
	}
}

func TestClaudeFODeferredLoadTraceBoundaries(t *testing.T) {
	spec := fixtureFOLoadSpec(t, "claude")
	root := filepath.Dir(spec.paths[foWriteCore]) + string(os.PathSeparator)
	read := func(id, name string) []string {
		return []string{
			claudeBoundaryEvent(id, "tool_use", "Read", map[string]any{"file_path": root + name}),
			claudeBoundaryEvent(id, "tool_result", "", nil),
		}
	}
	boot := append(read("shared", "first-officer-shared-core.md"), read("adapter", "claude-first-officer-runtime.md")...)
	write, merge := read("write", "fo-write-core.md"), read("merge", "fo-merge-core.md")
	terminal := claudeBoundaryEvent("terminal", "tool_use", "Bash", map[string]any{"command": "spacedock status --set merge-check status=done"})
	lines := append(append(append([]string{}, boot...), write...), merge...)
	lines = append(lines, terminal)
	if err := assertFOTerminalLoadBoundary(claudeFOLoadTrace(strings.Join(lines, "\n"), spec, claudeTerminalAction("merge-check"))); err != nil {
		t.Fatalf("terminal: %v", err)
	}
	wrapper := claudeBoundaryEvent("wrapper", "tool_use", "Skill", map[string]any{"skill": "spacedock:fo-write-core"})
	if err := traceViolation(claudeFOLoadTrace(strings.Join(append(boot, wrapper), "\n"), spec, nil)); err == nil {
		t.Fatal("wrapper-skill invocation must fail")
	}

	filing := claudeBoundaryEvent("filing", "tool_use", "Bash", map[string]any{"command": "spacedock new wire-the-thing"})
	bad := []struct {
		name  string
		lines []string
	}{
		{"path only", []string{claudeBoundaryEvent("fake", "tool_use", "Bash", map[string]any{"command": "printf " + root + "fo-write-core.md"}), claudeBoundaryEvent("fake", "tool_result", "", nil), filing}},
		{"partial read", []string{claudeBoundaryEvent("partial", "tool_use", "Read", map[string]any{"file_path": root + "fo-write-core.md", "limit": 1}), claudeBoundaryEvent("partial", "tool_result", "", nil), filing}},
		{"same suffix alternate root", []string{claudeBoundaryEvent("alternate", "tool_use", "Read", map[string]any{"file_path": "/alternate/skills/first-officer/references/fo-write-core.md"}), claudeBoundaryEvent("alternate", "tool_result", "", nil), filing}},
		{"early ordinary mutation", append([]string{claudeBoundaryEvent("early", "tool_use", "Bash", map[string]any{"command": "spacedock status --set ordinary status=review"})}, append(write, filing)...)},
		{"early state commit", append([]string{claudeBoundaryEvent("early-state-commit", "tool_use", "Bash", map[string]any{"command": "spacedock state commit ordinary"})}, append(write, filing)...)},
		{"early status archive", append([]string{claudeBoundaryEvent("early-status-archive", "tool_use", "Bash", map[string]any{"command": "spacedock status --archive ordinary"})}, append(write, filing)...)},
		{"early merge guard", append([]string{claudeBoundaryEvent("early-merge-guard", "tool_use", "Bash", map[string]any{"command": "spacedock merge guard ordinary --verdict passed"})}, append(write, filing)...)},
	}
	for _, tc := range bad {
		lines := append(append([]string{}, boot...), tc.lines...)
		if err := assertFOFilingLoadBoundary(claudeFOLoadTrace(strings.Join(lines, "\n"), spec, claudeFilingAction("wire-the-thing"))); err == nil {
			t.Errorf("%s trace must fail", tc.name)
		}
	}
}
