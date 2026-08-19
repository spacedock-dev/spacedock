// ABOUTME: Codex dispatch-build host shape — the outer prompt and dispatch-file
// ABOUTME: body use Codex mailbox semantics instead of Claude Skill()/SendMessage.
package dispatch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCodexHostPromptShape(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
		"host":           "codex",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	var out struct {
		Prompt string `json:"prompt"`
	}
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}
	dispatchPath := dispatchFilePathFromStdout(t, native.stdout)
	assertCodexFreshPrompt(t, out.Prompt, dispatchPath)

	body := readDispatchBody(t, dispatchPath)
	for _, banned := range []string{"Skill(skill=\"spacedock:ensign\")", "SendMessage(to=\"team-lead\""} {
		if strings.Contains(body, banned) {
			t.Fatalf("codex dispatch body must omit %q:\n%s", banned, body)
		}
	}
	for _, want := range []string{
		"Read this dispatch file directly",
		"outer fresh-worker prompt invokes `$spacedock:ensign`",
		"file supplies the stage-specific assignment",
		"Codex final-status notification",
		"FO mailbox",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("codex dispatch body missing %q:\n%s", want, body)
		}
	}
}

func assertCodexFreshPrompt(t *testing.T, prompt, dispatchPath string) {
	t.Helper()
	want := fmt.Sprintf("$spacedock:ensign; then Read %s and treat its content as your assignment.", dispatchPath)
	if prompt != want {
		t.Fatalf("Codex fresh prompt = %q, want exact pointer bootstrap %q", prompt, want)
	}
}

func TestBuildCodexHostRejectsBareModeBeforeArtifactCreation(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeNonWorktreeStages())
	entityPath := filepath.Join(root, "f02codexbare.md")
	writeFile(t, entityPath, entityFM("Codex Bare", "implementation", ""))
	gitInit(t, root)
	artifactPath := filepath.Join(dispatchFileDir, "spacedock-ensign-f02codexbare-implementation.md")
	_ = os.Remove(artifactPath)
	t.Cleanup(func() { _ = os.Remove(artifactPath) })
	fields := map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- prove Codex host shape"},
		"bare_mode":      true,
		"host":           "codex",
	}
	bare := runNative(mergeStdin(fields, nil), "build", "--workflow-dir", root)
	if bare.exit != 2 {
		t.Fatalf("bare Codex build exit=%d, want 2; stderr=%q", bare.exit, bare.stderr)
	}
	const diagnostic = "bare_mode is unsupported on host codex; Codex worker.spawn requires a named spawn_agent task"
	if !strings.Contains(bare.stderr, diagnostic) {
		t.Fatalf("bare Codex diagnostic=%q, want %q", bare.stderr, diagnostic)
	}
	if _, err := os.Stat(artifactPath); !os.IsNotExist(err) {
		t.Fatalf("bare Codex created dispatch artifact %s; stat err=%v", artifactPath, err)
	}
	fields["bare_mode"] = false
	named := runNative(mergeStdin(fields, nil), "build", "--workflow-dir", root)
	if named.exit != 0 {
		t.Fatalf("named Codex build exit=%d, want 0; stderr=%q", named.exit, named.stderr)
	}
	var out struct{ Name, Prompt string }
	if err := json.Unmarshal([]byte(named.stdout), &out); err != nil {
		t.Fatalf("named Codex stdout is not JSON: %v", err)
	}
	if out.Name == "" || out.Prompt == "" {
		t.Fatalf("named Codex build must retain name and prompt: %s", named.stdout)
	}
}

// TestBuildCodexHostIgnoresModelWithNote is AC-5: dispatch build --host codex
// over a fable-declaring README exits 0, emits model: null, and prints the
// ignore-with-note stderr line in place of the effective_model line — model is
// outside codex's dispatch-settable model space (the thread's model is what
// runs), so a claude-enum value is dropped rather than rejected.
func TestBuildCodexHostIgnoresModelWithNote(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeModelsFable)
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "stagemodel", ""))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "stagemodel",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
		"host":           "codex",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	assertGolden(t, "build-host-codex-model-ignored", goldenEnvelope{res: normRun(native, root, home)})
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	if !strings.Contains(native.stdout, `"model": null`) {
		t.Errorf("stdout model must be null:\n%s", native.stdout)
	}
	if strings.Contains(native.stderr, "effective_model=") {
		t.Errorf("stderr must not contain the effective_model= line on host=codex:\n%s", native.stderr)
	}
}
