// ABOUTME: Behavioral input-mode tests for `dispatch build` — pins the
// ABOUTME: reuse-advance contradiction and runs every --help example through the parser.
package dispatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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

// TestDispatchBuildAdvanceInputMode exercises file and stdin advance envelopes,
// while a retired JSON caller fails with the exact migration diagnostic.
func TestDispatchBuildAdvanceInputMode(t *testing.T) {
	for _, source := range []string{"file", "stdin"} {
		t.Run(source, func(t *testing.T) {
			workflowDir, leaf := helpExampleFixture(t)
			checklist := leaf["validation.checklist"]
			if source == "stdin" {
				checklist = "-"
			}
			res := runNative("- verify\n", "build",
				"--workflow-dir", workflowDir,
				"--entity-path", leaf["thing.md"],
				"--stage", "validation",
				"--checklist-file", checklist,
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
	}

	t.Run("stdin JSON + --advance is rejected", func(t *testing.T) {
		workflowDir, leaf := helpExampleFixture(t)
		_ = leaf
		stdin := `{"schema_version":2,"checklist":["- verify"]}`
		res := runNative(stdin, "build", "--workflow-dir", workflowDir, "--advance")
		if res.exit != 2 {
			t.Fatalf("stdin JSON + --advance exit=%d, want 2\nstdout=%q\nstderr=%q", res.exit, res.stdout, res.stderr)
		}
		if res.stdout != "" {
			t.Fatalf("stdin JSON + --advance stdout=%q, want empty", res.stdout)
		}
		wantErr := "error: flag/file input requires --entity-path, --stage, and --checklist-file\nJSON request input is retired; pass checklist lines on stdin with --checklist-file -"
		if strings.TrimSpace(res.stderr) != wantErr {
			t.Fatalf("stderr=%q, want %q", res.stderr, wantErr)
		}
	})
}

// TestDispatchBuildHelpExamplesParse is AC-3: every positive example printed in
// the `dispatch build --help` Examples section is run through the real parser
// against a minimal fixture. Only the leaf path values are rewritten to the
// fixture's real files — flags, checklist text, stage names, and
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

// helpExamples returns complete executable examples, including their heredoc bodies.
func helpExamples(t *testing.T, help string) []string {
	t.Helper()
	idx := strings.Index(help, "Examples:")
	if idx < 0 {
		t.Fatalf("help has no Examples: section:\n%s", help)
	}
	lines := strings.Split(help[idx:], "\n")
	var out []string
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "spacedock dispatch build") {
			continue
		}
		if strings.HasSuffix(line, "<<'CHECKLIST'") {
			for i++; i < len(lines); i++ {
				line += "\n" + lines[i]
				if lines[i] == "CHECKLIST" {
					break
				}
			}
		}
		out = append(out, line)
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
// fixture's real files, then runs it through the native parser with the printed
// heredoc contents supplied as stdin. Everything except leaf paths is preserved
// verbatim, so a renamed flag or dropped field still fails the parse.
func runHelpExample(example, workflowDir string, leaf map[string]string) runResult {
	stdin := ""
	if i := strings.Index(example, " <<'CHECKLIST'"); i >= 0 {
		stdin = strings.TrimSuffix(strings.TrimPrefix(example[i+len(" <<'CHECKLIST'"):], "\n"), "CHECKLIST")
		example = example[:i]
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
	return runNative(stdin, args...)
}

func TestBuildChecklistStdin(t *testing.T) {
	root, entity := buildHostFixture(t)
	before, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{"DONE: verify", "\r\n  `ticks` $HOME 雪  \r\n\t\r\nsecond\r\n", strings.Repeat("x", 70000), `{"literal":"JSON"}`} {
		result := runNative(input, "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-")
		if result.exit != 0 {
			t.Fatalf("stdin build exit=%d: %s", result.exit, result.stderr)
		}
	}
	after, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatal("stdin dispatch added a scratch input file")
	}
	path := filepath.Join(root, "existing.checklist")
	writeFile(t, path, "DONE: verify")
	result := runNative("", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", path)
	files, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if result.exit != 0 || len(files) != len(before)+1 {
		t.Fatalf("file baseline requires one checklist input: %+v", result)
	}

}

func TestBuildRequestControlsRetired(t *testing.T) {
	for _, flag := range []string{"--print-schema", "--print-schema=true", "--validate-only", "--validate-only=missing.json"} {
		result := runNative("", "build", flag)
		if result.exit != 2 || result.stdout != "" || !strings.Contains(result.stderr, "JSON request input and --print-schema/--validate-only are retired") {
			t.Fatalf("%s: %+v", flag, result)
		}
	}
}

type failedChecklistReader struct{ called bool }

func (r *failedChecklistReader) Read([]byte) (int, error) {
	r.called = true
	return 0, fmt.Errorf("reader sentinel")
}

func TestBuildInputFailuresDoNotReadOrStamp(t *testing.T) {
	root, entity := buildHostFixture(t)
	before, _ := os.ReadFile(entity)
	head := gitOutput(t, root, "rev-parse", "HEAD")
	artifactPath := filepath.Join(dispatchFileDir, "spacedock-ensign-thing-backlog.md")
	artifactBefore, artifactErr := os.ReadFile(artifactPath)
	assertUnchanged := func(t *testing.T) {
		t.Helper()
		after, _ := os.ReadFile(entity)
		artifactAfter, afterErr := os.ReadFile(artifactPath)
		if string(after) != string(before) || gitOutput(t, root, "rev-parse", "HEAD") != head || string(artifactAfter) != string(artifactBefore) || (artifactErr == nil) != (afterErr == nil) {
			t.Fatal("input failure mutated entity, commit, or dispatch artifact")
		}
	}
	base := []string{"build", "--host", "claude", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-", "--stamp"}
	for _, tc := range []struct {
		name string
		args []string
		read bool
		code int
		want string
	}{
		{"reader failure", base, true, 1, "failed to read checklist stdin: reader sentinel"},
		{"file failure", append(append([]string{}, base...), "--checklist-file", filepath.Join(root, "missing")), false, 1, "failed to read checklist file"},
		{"missing entity", []string{"build", "--workflow-dir", root, "--stage", "backlog", "--checklist-file", "-"}, false, 2, "JSON request input is retired"},
		{"missing stage", []string{"build", "--workflow-dir", root, "--entity-path", entity, "--checklist-file", "-"}, false, 2, "JSON request input is retired"},
		{"missing checklist", []string{"build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog"}, false, 2, "JSON request input is retired"},
		{"required trio", []string{"build", "--workflow-dir", root}, false, 2, "JSON request input is retired"},
		{"required workflow", []string{"build"}, false, 2, "requires --workflow-dir"},
		{"stamp advance", append(append([]string{}, base...), "--advance"), false, 2, "--stamp is incompatible with --advance"},
		{"bare advance", append(append([]string{}, base...), "--bare-mode", "--advance"), false, 2, "--advance is incompatible with --bare-mode"},
		{"retired equals", []string{"build", "--validate-only=missing"}, false, 2, "are retired"},
		{"retired schema", []string{"build", "--print-schema"}, false, 2, "are retired"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := &failedChecklistReader{}
			var out, errout bytes.Buffer
			code := RunWithLauncher(nil, testWorkflowLauncher, tc.args, reader, &out, &errout)
			if code != tc.code || out.Len() != 0 || !strings.Contains(errout.String(), tc.want) || reader.called != tc.read {
				t.Fatalf("exit=%d read=%v stdout=%q stderr=%q", code, reader.called, out.String(), errout.String())
			}
			assertUnchanged(t)
		})
	}
	legacy := runNative(`{"schema_version":2,"stage":"backlog","checklist":["retired request"]}`, "build", "--workflow-dir", root)
	if legacy.exit != 2 || legacy.stdout != "" || !strings.Contains(legacy.stderr, "JSON request input is retired") {
		t.Fatalf("legacy request: %+v", legacy)
	}
	assertUnchanged(t)
	for _, input := range []string{"", " \t\r\n\n"} {
		for _, source := range []string{"-", filepath.Join(t.TempDir(), "empty")} {
			if source != "-" {
				writeFile(t, source, input)
			}
			result := runNative(input, "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", source)
			if result.exit != 1 || result.stdout != "" || result.stderr != "error: missing required field 'checklist'\n" {
				t.Fatalf("empty input: %+v", result)
			}
			assertUnchanged(t)
		}
	}
}

func TestChecklistFileDoesNotReadStdin(t *testing.T) {
	root, entity := buildHostFixture(t)
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	writeFile(t, "-", "DONE: literal dash file\n")
	reader := &failedChecklistReader{}
	var out, errout bytes.Buffer
	args := []string{"build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "./-", "--host", "claude", "--scope-notes-file", "-", "--feedback-context-file", "-"}
	code := RunWithLauncher(nil, testWorkflowLauncher, args, reader, &out, &errout)
	if code != 0 || reader.called {
		t.Fatalf("exit=%d read=%v stderr=%s", code, reader.called, errout.String())
	}
	for i, arg := range args {
		if arg == "--checklist-file" {
			args[i+1] = "-"
		}
	}
	out.Reset()
	errout.Reset()
	code = RunWithLauncher(nil, testWorkflowLauncher, args, strings.NewReader("DONE: actual stdin"), &out, &errout)
	if code != 0 {
		t.Fatalf("dash stdin with a false filename sentinel: %s", errout.String())
	}
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, out.String()))
	if !strings.Contains(body, "### Completion checklist\n\nDONE: actual stdin\n\n### Summary") {
		t.Fatal("dash opened the file named - instead of reading stdin")
	}

}
