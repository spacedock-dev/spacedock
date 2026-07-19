// ABOUTME: Behavioral input-mode tests for `dispatch build` — pins the
// ABOUTME: reuse-advance contradiction and runs every --help example through the parser.
package dispatch

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// readmeNonWorktreeStages is a workflow README whose implementation and
// validation stages are non-worktree, so the help examples parse without a
// worktree directory on disk.
func readmeNonWorktreeStages() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"    - name: validation\n" +
		"      feedback-to: implementation\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Fixture Workflow\n" +
		"\n" +
		"### backlog\n\nseed.\n\n- **Outputs:** x.\n\n" +
		"### implementation\n\nwork.\n\n- **Outputs:** y.\n\n" +
		"### validation\n\nverify.\n\n- **Outputs:** z.\n\n" +
		"### done\n\nterm.\n"
}

// helpExampleFixture materializes the fixture the rendered --help examples name:
// a non-worktree workflow README, a thing.md entity, and the two checklist files
// the examples reference. It returns a leaf-value → absolute-path map keyed by the
// exact leaf spellings the help prints, so a caller can rewrite each printed
// example's paths to real files without touching flag names, field names, stage
// names, or --advance presence.
func helpExampleFixture(t *testing.T) (workflowDir string, leaf map[string]string) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeNonWorktreeStages())
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	implChecklist := filepath.Join(root, "impl.checklist")
	writeFile(t, implChecklist, "- run tests\n")
	validationChecklist := filepath.Join(root, "validation.checklist")
	writeFile(t, validationChecklist, "- verify\n")
	gitInit(t, root)
	return root, map[string]string{
		"thing.md":             entityPath,
		".":                    root,
		"impl.checklist":       implChecklist,
		"validation.checklist": validationChecklist,
	}
}

// TestDispatchBuildAdvanceInputMode is AC-2: the reuse-advance form is
// unambiguous. The accepted flag/file --advance form exits 0 with an advance
// envelope on stdout and empty stderr; the contradictory stdin-JSON-plus-CLI
// --advance form (no flag/file trio) exits 2 with the exact trio-required error.
func TestDispatchBuildAdvanceInputMode(t *testing.T) {
	t.Run("flag/file --advance succeeds", func(t *testing.T) {
		workflowDir, leaf := helpExampleFixture(t)
		res := runNative("", "build",
			"--workflow-dir", workflowDir,
			"--entity-path", leaf["thing.md"],
			"--stage", "validation",
			"--checklist-file", leaf["validation.checklist"],
			"--advance")
		if res.exit != 0 {
			t.Fatalf("flag/file --advance exit=%d, want 0\nstderr=%q", res.exit, res.stderr)
		}
		if res.stderr != "" {
			t.Fatalf("flag/file --advance stderr=%q, want empty", res.stderr)
		}
		var env map[string]json.RawMessage
		if err := json.Unmarshal([]byte(res.stdout), &env); err != nil {
			t.Fatalf("stdout is not JSON: %v\n%s", err, res.stdout)
		}
		for _, want := range []string{"schema_version", "prompt", "model"} {
			if _, ok := env[want]; !ok {
				t.Errorf("advance envelope missing %q: %s", want, res.stdout)
			}
		}
		for _, banned := range []string{"subagent_type", "name"} {
			if _, ok := env[banned]; ok {
				t.Errorf("advance envelope must omit spawn field %q: %s", banned, res.stdout)
			}
		}
	})

	t.Run("stdin JSON + --advance is rejected", func(t *testing.T) {
		workflowDir, leaf := helpExampleFixture(t)
		stdin := mergeStdin(map[string]any{
			"schema_version": 2,
			"entity_path":    leaf["thing.md"],
			"workflow_dir":   workflowDir,
			"stage":          "validation",
			"checklist":      []string{"- verify"},
		}, nil)
		res := runNative(stdin, "build", "--workflow-dir", workflowDir, "--advance")
		if res.exit != 2 {
			t.Fatalf("stdin JSON + --advance exit=%d, want 2\nstdout=%q\nstderr=%q", res.exit, res.stdout, res.stderr)
		}
		if res.stdout != "" {
			t.Fatalf("stdin JSON + --advance stdout=%q, want empty", res.stdout)
		}
		wantErr := "error: flag/file input requires --entity-path, --stage, and --checklist-file"
		if strings.TrimSpace(res.stderr) != wantErr {
			t.Fatalf("stderr=%q, want %q", res.stderr, wantErr)
		}
	})
}

// TestDispatchBuildHelpExamplesParse is AC-3: every positive example printed in
// the `dispatch build --help` Examples section is run through the real parser
// against a minimal fixture. Only the leaf path values are rewritten to the
// fixture's real files — flag names, field names, JSON shape, stage names, and
// --advance presence stay exactly as printed. If any printed example drops a
// required field, renames a flag, or advertises a stdin+--advance form the parser
// rejects, the matching run exits non-zero and this test fails.
func TestDispatchBuildHelpExamplesParse(t *testing.T) {
	var help bytes.Buffer
	printBuildUsage(&help)
	examples := helpExamples(t, help.String())
	if len(examples) < 3 {
		t.Fatalf("expected at least 3 examples in --help, found %d:\n%v", len(examples), examples)
	}
	sawAdvance := false
	for _, ex := range examples {
		if strings.Contains(ex, "--advance") {
			sawAdvance = true
		}
	}
	if !sawAdvance {
		t.Fatalf("no reuse-advance example found among rendered examples:\n%v", examples)
	}

	for _, ex := range examples {
		workflowDir, leaf := helpExampleFixture(t)
		res := runHelpExample(ex, workflowDir, leaf)
		if res.exit != 0 {
			t.Fatalf("rendered --help example failed to parse (exit=%d):\n  example: %s\n  stderr: %s",
				res.exit, ex, res.stderr)
		}
	}
}

// helpExamples returns the positive examples in the help's Examples section: a
// line beginning `{` is a stdin JSON example, a line beginning `spacedock
// dispatch build` is a flag/file example. Scoped to the Examples section so the
// placeholder Usage lines (DIR/FILE/STAGE) are not mistaken for examples.
func helpExamples(t *testing.T, help string) []string {
	t.Helper()
	idx := strings.Index(help, "Examples:")
	if idx < 0 {
		t.Fatalf("help has no Examples: section:\n%s", help)
	}
	var out []string
	for _, line := range strings.Split(help[idx:], "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "{") || strings.HasPrefix(line, "spacedock dispatch build") {
			out = append(out, line)
		}
	}
	return out
}

// pathValueFlags are the flag/file-mode flags whose following token is a path the
// help renders as a leaf spelling; the fixture rewrites those to real files.
var pathValueFlags = map[string]bool{
	"--workflow-dir":          true,
	"--entity-path":           true,
	"--checklist-file":        true,
	"--scope-notes-file":      true,
	"--feedback-context-file": true,
}

// runHelpExample rewrites the leaf path values in a single rendered example to the
// fixture's real files, then runs it through the native parser. A stdin JSON
// example runs in stdin mode; a `spacedock dispatch build` line runs as its
// flag/file argument vector. Everything except leaf path values is preserved
// verbatim, so a renamed flag or dropped field still fails the parse.
func runHelpExample(example, workflowDir string, leaf map[string]string) runResult {
	if strings.HasPrefix(example, "{") {
		var fields map[string]any
		if err := json.Unmarshal([]byte(example), &fields); err != nil {
			return runResult{stderr: "example is not JSON: " + err.Error(), exit: 99}
		}
		for _, key := range []string{"entity_path", "workflow_dir"} {
			if s, ok := fields[key].(string); ok {
				if real, mapped := leaf[s]; mapped {
					fields[key] = real
				}
			}
		}
		raw, _ := json.Marshal(fields)
		return runNative(string(raw), "build", "--workflow-dir", workflowDir)
	}

	tokens := strings.Fields(example)
	// Drop the `spacedock dispatch` prefix; runNative takes args from `build`.
	args := tokens[2:]
	for i := 0; i < len(args); i++ {
		if pathValueFlags[args[i]] && i+1 < len(args) {
			if real, mapped := leaf[args[i+1]]; mapped {
				args[i+1] = real
			}
			i++
		}
	}
	return runNative("", args...)
}
