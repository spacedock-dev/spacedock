// ABOUTME: --advance mode golden fixtures — plain / split-root / feedback-reflow
// ABOUTME: / codex-host shapes, flag validation, filename suffix, and envelope shape.
package dispatch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// advanceCase is a positive advance-mode golden fixture: it builds a fixture
// tree, runs the native --advance build, and asserts golden envelope + body.
type advanceCase struct {
	name      string
	splitRoot bool
	stage     string
	feedback  bool
	host      string // "" = default (claude, via runNativeWithDefaultClaudeHost)
}

func TestBuildAdvanceGoldens(t *testing.T) {
	cases := []advanceCase{
		{name: "plain", splitRoot: false, stage: "validation"},
		{name: "split-root", splitRoot: true, stage: "validation"},
		{
			name: "feedback-reflow", splitRoot: false, stage: "implementation",
			feedback: true,
		},
		{name: "codex-host", splitRoot: false, stage: "validation", host: "codex"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()

			workflowDir := root
			if tc.splitRoot {
				workflowDir = filepath.Join(root, "state-checkout")
			}
			writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(tc.splitRoot))

			worktreeRel := ".worktrees/spacedock-ensign-thing"
			if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
				t.Fatal(err)
			}
			entityPath := filepath.Join(workflowDir, "thing.md")
			writeFile(t, entityPath, entityFM("Thing", tc.stage, worktreeRel))

			gitInit(t, root)
			if tc.splitRoot {
				stateCheckout := filepath.Join(workflowDir, "state-checkout")
				if err := os.MkdirAll(stateCheckout, 0o755); err != nil {
					t.Fatal(err)
				}
				gitInitBare(t, stateCheckout)
				gitAddOrigin(t, stateCheckout)
			}

			stdin := "- a\n- b"
			stdinArgs := []string{"build", "--workflow-dir", workflowDir, "--entity-path", entityPath, "--stage", tc.stage, "--checklist-file", "-", "--advance"}
			if tc.host != "" {
				stdinArgs = append(stdinArgs, "--host", tc.host)
			}
			if tc.feedback {
				feedback := filepath.Join(t.TempDir(), "feedback.md")
				writeFile(t, feedback, "REJECTED: do better.")
				stdinArgs = append(stdinArgs, "--feedback-reflow", "--feedback-context-file", feedback)
			}

			native := runNative(stdin, stdinArgs...)
			if native.exit != 0 {
				t.Fatalf("advance build exit=%d stderr=%s", native.exit, native.stderr)
			}
			if tc.host == "codex" {
				out := decodeBuildOutput(t, native.stdout)
				if strings.HasPrefix(out.Prompt, "$spacedock:ensign; then Read ") {
					t.Fatalf("Codex advance prompt must not repeat the fresh bootstrap: %q", out.Prompt)
				}
			}
			nativeBody := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

			env := goldenEnvelope{res: normRun(native, root, home), body: normPaths(nativeBody, root, home)}
			assertGolden(t, "build-advance-"+tc.name, env)
		})
	}
}

// TestBuildAdvanceFilenameSuffix asserts the dispatch filename carries the
// -advance suffix after the merged-mode session disambiguator prefix, so an
// advance file can never alias a fresh-dispatch file for the same slug+stage.
func TestBuildAdvanceFilenameSuffix(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_CODE_SESSION_ID", "sessionaaa")
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "validation", worktreeRel))
	gitInit(t, root)

	stdin := strings.Join([]string{"- a"}, "\n")
	stdinArgs := []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", "validation", "--checklist-file", "-", "--advance"}

	native := runNative(stdin, stdinArgs...)
	if native.exit != 0 {
		t.Fatalf("advance build exit=%d stderr=%s", native.exit, native.stderr)
	}
	path := dispatchFilePathFromStdout(t, native.stdout)
	base := filepath.Base(path)
	if base != "sessionaaa-spacedock-ensign-thing-validation-advance.md" {
		t.Fatalf("unexpected advance dispatch filename: %s", base)
	}

	// A fresh (non-advance) dispatch for the identical session/slug/stage must
	// write to the bare filename — no collision with the advance file.
	freshStdin := strings.Join([]string{"- a"}, "\n")
	freshStdinArgs := []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", "validation", "--checklist-file", "-"}
	freshNative := runNative(freshStdin, freshStdinArgs...)
	if freshNative.exit != 0 {
		t.Fatalf("fresh build exit=%d stderr=%s", freshNative.exit, freshNative.stderr)
	}
	freshPath := dispatchFilePathFromStdout(t, freshNative.stdout)
	if freshPath == path {
		t.Fatalf("fresh dispatch file collided with advance dispatch file: %s", freshPath)
	}
	if filepath.Base(freshPath) != "sessionaaa-spacedock-ensign-thing-validation.md" {
		t.Fatalf("unexpected fresh dispatch filename: %s", filepath.Base(freshPath))
	}
}

// TestBuildAdvanceEnvelopeOmitsSpawnFields is AC-2/the delta-5 shape check: the
// advance envelope carries schema_version/description/fetch_commands/
// dispatch_file_path/prompt/model and nothing else — no subagent_type, name,
// team_name, or run_in_background, since nothing is spawned.
func TestBuildAdvanceEnvelopeOmitsSpawnFields(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "validation", worktreeRel))
	gitInit(t, root)

	stdin := strings.Join([]string{"- a"}, "\n")
	stdinArgs := []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", "validation", "--checklist-file", "-", "--advance"}

	native := runNative(stdin, stdinArgs...)
	if native.exit != 0 {
		t.Fatalf("advance build exit=%d stderr=%s", native.exit, native.stderr)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(native.stdout), &raw); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, native.stdout)
	}
	for _, banned := range []string{"subagent_type", "name", "team_name", "run_in_background"} {
		if _, ok := raw[banned]; ok {
			t.Errorf("advance envelope must omit %q entirely: %s", banned, native.stdout)
		}
	}
	for _, want := range []string{"schema_version", "description", "fetch_commands", "dispatch_file_path", "prompt", "model"} {
		if _, ok := raw[want]; !ok {
			t.Errorf("advance envelope missing required field %q: %s", want, native.stdout)
		}
	}
}

// TestBuildAdvanceBareModeConflict is delta-6: --advance + --bare-mode (via
// stdin fields, exercising the same runBuildFields guard the CLI-flag path
// hits) is a usage error, exit 2 — a reuse advance presupposes an addressable
// worker, and bare mode has none.
func TestBuildAdvanceBareModeConflict(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	gitInit(t, root)

	stdin := strings.Join([]string{"- a"}, "\n")
	stdinArgs := []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", "backlog", "--checklist-file", "-", "--bare-mode", "--advance"}

	native := runNative(stdin, stdinArgs...)
	if native.exit != 2 {
		t.Fatalf("advance+bare_mode exit=%d, want 2; stderr=%s", native.exit, native.stderr)
	}
}

// TestBuildAdvanceCLIFlagBareModeConflict is the CLI-flag-level half of
// delta-6: `--advance --bare-mode` on the command line fails fast (exit 2)
// before any field loading, matching the stdin-field guard above.
func TestBuildAdvanceCLIFlagBareModeConflict(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	gitInit(t, root)
	checklistFile := filepath.Join(root, "checklist.txt")
	writeFile(t, checklistFile, "- a\n")

	native := runNative("", "build", "--advance", "--bare-mode",
		"--workflow-dir", root, "--entity-path", entityPath, "--stage", "backlog",
		"--checklist-file", checklistFile)
	if native.exit != 2 {
		t.Fatalf("--advance --bare-mode exit=%d, want 2; stderr=%s", native.exit, native.stderr)
	}
}

// TestBuildAdvanceFeedbackReflowRequiresContext is the rule-5 error case
// carried into advance mode: --feedback-reflow without feedback_context is a
// loud failure, identical to the fresh-dispatch guard.
func TestBuildAdvanceFeedbackReflowRequiresContext(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)

	stdin := strings.Join([]string{"- a"}, "\n")
	stdinArgs := []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", "implementation", "--checklist-file", "-", "--advance", "--feedback-reflow"}

	native := runNative(stdin, stdinArgs...)
	if native.exit != 1 {
		t.Fatalf("advance feedback-reflow without context exit=%d, want 1; stderr=%s", native.exit, native.stderr)
	}
	wantFragment := "feedback_context is missing"
	if !strings.Contains(native.stderr, wantFragment) {
		t.Fatalf("stderr %q missing fragment %q", native.stderr, wantFragment)
	}
}
