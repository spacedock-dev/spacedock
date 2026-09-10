// ABOUTME: Build byte-hazard parity — model precedence stderr/null literal,
// ABOUTME: no-HTML-escape, and shlex space-quoting locked against the oracle.
package dispatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// readmeModels declares model at three precedence sites so a fixture can pick a
// stage to toggle which one wins: stage-set (opus), defaults-set only (haiku),
// or neither (null). The states are non-worktree backlog stages.
const readmeModels = `---
entity-type: task
id-style: slug
stages:
  defaults:
    worktree: false
    concurrency: 1
    model: haiku
  states:
    - name: stagemodel
      initial: true
      model: opus
    - name: defaultsmodel
    - name: done
      terminal: true
---
# Models Fixture

### stagemodel

stage wins.

- **Outputs:** x.

### defaultsmodel

defaults win.

- **Outputs:** y.

### done

term.
`

// readmeModelsFable mirrors readmeModels but with fable at both the stage and
// defaults precedence sites, for the AC-1 fable-joins-the-enum cases.
const readmeModelsFable = `---
entity-type: task
id-style: slug
stages:
  defaults:
    worktree: false
    concurrency: 1
    model: fable
  states:
    - name: stagemodel
      initial: true
      model: fable
    - name: defaultsmodel
    - name: done
      terminal: true
---
# Fable Models Fixture

### stagemodel

stage wins.

- **Outputs:** x.

### defaultsmodel

defaults win.

- **Outputs:** y.

### done

term.
`

// readmeModelsNull mirrors readmeModels but with no model anywhere, so the
// effective model resolves to null (empty stderr, "model": null in stdout).
const readmeModelsNull = `---
entity-type: task
id-style: slug
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: nomodel
      initial: true
    - name: done
      terminal: true
---
# No Models Fixture

### nomodel

no model anywhere.

- **Outputs:** x.

### done

term.
`

// TestBuildModelPrecedence locks the effective_model precedence (stage >
// defaults > null): the [build] effective_model stderr line (with the U+2192
// arrow) and the "model" JSON value across all three sources.
func TestBuildModelPrecedence(t *testing.T) {
	type modelCase struct {
		name      string
		readme    string
		stage     string
		wantModel string // expected JSON model value, "" => null
	}
	cases := []modelCase{
		{name: "stage-wins-opus", readme: readmeModels, stage: "stagemodel", wantModel: "opus"},
		{name: "defaults-haiku", readme: readmeModels, stage: "defaultsmodel", wantModel: "haiku"},
		{name: "null", readme: readmeModelsNull, stage: "nomodel", wantModel: ""},
		{name: "stage-fable", readme: readmeModelsFable, stage: "stagemodel", wantModel: "fable"},
		{name: "defaults-fable", readme: readmeModelsFable, stage: "defaultsmodel", wantModel: "fable"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "README.md"), tc.readme)
			entityPath := filepath.Join(root, "thing.md")
			writeFile(t, entityPath, entityFM("Thing", tc.stage, ""))
			gitInit(t, root)

			stdin := strings.Join([]string{"- a"}, "\n")
			stdinArgs := []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", tc.stage, "--checklist-file", "-"}

			native := runNative(stdin, stdinArgs...)
			assertGolden(t, "build-model-"+tc.name, goldenEnvelope{res: normRun(native, root, home)})

			// Lock the "model" JSON value explicitly (the golden already covers it,
			// but this names the contract: null literal vs string).
			var out map[string]json.RawMessage
			if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
				t.Fatalf("native stdout not JSON: %v", err)
			}
			gotModel := string(out["model"])
			wantModel := "null"
			if tc.wantModel != "" {
				wantModel = `"` + tc.wantModel + `"`
			}
			if gotModel != wantModel {
				t.Errorf("model = %s, want %s", gotModel, wantModel)
			}
		})
	}
}

// TestBuildNoHTMLEscape locks that <, >, & survive verbatim in the emitted JSON
// (Go's default json.Marshal escapes them to < / > / &). The
// entity title carries all three, so they flow into the stdout JSON description
// field — the only build-output JSON field that can carry caller prose. The
// feedback_context field carries them too, exercising the dispatch body (plain
// markdown, no escaping). Native and oracle stdout + body must be byte-identical.
func TestBuildNoHTMLEscape(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	entityPath := filepath.Join(root, "thing.md")
	htmlTitle := "Compare a < b && c > d in <Tag>"
	writeFile(t, entityPath, entityFM(htmlTitle, "validation", ".worktrees/spacedock-ensign-thing"))
	if err := os.MkdirAll(filepath.Join(root, ".worktrees/spacedock-ensign-thing"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitInit(t, root)

	stdinFeedbackContext := filepath.Join(t.TempDir(), "feedback_context.md")
	writeFile(t, stdinFeedbackContext, "compare a < b && c > d in <Tag> & raw &amp; ampersand")
	stdin := strings.Join([]string{"- a < b", "- c > d & e"}, "\n")
	stdinArgs := []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", "validation", "--checklist-file", "-", "--feedback-reflow", "--feedback-context-file", stdinFeedbackContext}

	native := runNative(stdin, stdinArgs...)
	nativeBody := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	env := goldenEnvelope{res: normRun(native, root, home), body: normPaths(nativeBody, root, home)}
	assertGolden(t, "build-html-escape", env)

	// Explicit: the raw <, >, & bytes survive in the native stdout description,
	// and no \uXXXX-escaped sequences appear (Go's default HTML escaping is off).
	if !strings.Contains(native.stdout, htmlTitle) {
		t.Errorf("native stdout HTML-escaped the description title away:\n%s", native.stdout)
	}
	for _, ch := range []rune{'<', '>', '&'} {
		esc := fmt.Sprintf("\\u%04x", ch) // the \uXXXX sequence Go emits when escaping is ON
		if strings.Contains(native.stdout, esc) {
			t.Errorf("native stdout contains %s (SetEscapeHTML not off):\n%s", esc, native.stdout)
		}
	}
}

// TestBuildSpaceBearingPath locks shell-safe launcher quoting while retaining a
// workflow and launcher path that both contain spaces.
func TestBuildSpaceBearingPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	// A workflow dir whose name contains a space.
	workflowDir := filepath.Join(root, "work dir")
	writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(false))
	entityPath := filepath.Join(workflowDir, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	gitInit(t, root)

	stdin := strings.Join([]string{"- a"}, "\n")
	stdinArgs := []string{"build", "--workflow-dir", workflowDir, "--entity-path", entityPath, "--stage", "backlog", "--checklist-file", "-"}

	workflowLauncher := filepath.Join(root, "launcher dir", "spacedock")
	native := runNativeWithLauncher(stdin, workflowLauncher, append(stdinArgs, "--host", "claude")...)
	nativeBody := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	env := goldenEnvelope{res: normRun(native, root, home), body: normPaths(nativeBody, root, home)}
	assertGolden(t, "build-space-path", env)
	// Explicit: the space-bearing absolute launcher is a literal shell-safe prefix.
	wantQuoted := "    " + shlexQuote(workflowLauncher)
	if !strings.Contains(nativeBody, wantQuoted) {
		t.Errorf("workflow launcher is not shell-quoted:\nwant contains: %s\ngot:\n%s", wantQuoted, nativeBody)
	}
}

func TestBuildQuotedHeredocTransport(t *testing.T) {
	root, entity := buildHostFixture(t)
	binary := filepath.Join(t.TempDir(), "spacedock")
	build := exec.Command("go", "build", "-o", binary, "./cmd/spacedock")
	build.Dir = "../.."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, out)
	}
	sentinel := filepath.Join(t.TempDir(), "SHOULD_NOT_EXIST")
	line := "  Keep \"quotes\", `echo expanded`, $(touch " + shlexQuote(sentinel) + "), $HOME, 雪  "
	command := shlexQuote(binary) + " dispatch build --host claude --workflow-dir " + shlexQuote(root) + " --entity-path " + shlexQuote(entity) + " --stage backlog --checklist-file - <<'CHECKLIST'\n" + line + "\nCHECKLIST\n"
	cmd := exec.Command("sh", "-c", command)
	var out, errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	if err := cmd.Run(); err != nil {
		t.Fatalf("quoted heredoc: %v stderr=%s", err, errout.String())
	}
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, out.String()))
	if !strings.Contains(body, "### Completion checklist\n\n"+line+"\n\n### Summary") {
		t.Fatal("shell changed literal checklist bytes")
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("shell expanded checklist payload: %v", err)
	}
}
