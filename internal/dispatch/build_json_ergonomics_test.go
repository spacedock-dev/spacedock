// ABOUTME: dispatch build flag/file input and host-resolution ergonomics.
package dispatch

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildFlagFileInputModePreservesLiteralChecklist(t *testing.T) {
	t.Setenv("CODEX_THREAD_ID", "")
	t.Setenv("CLAUDECODE", "1")
	t.Setenv("PI_CODING_AGENT", "")
	t.Setenv("PI_CODING_AGENT_DIR", "")
	root, entityPath := buildHostFixture(t)
	checklistPath := filepath.Join(root, "checklist.txt")
	scopePath := filepath.Join(root, "scope.md")
	literalChecklist := "1. keep `sharedRuntimeScenarios()` and $CLAUDE_CODE_SESSION_ID literal"
	writeFile(t, checklistPath, "\n"+literalChecklist+"\n\n2. preserve Markdown exactly\n")
	writeFile(t, scopePath, "### Scope\nUse $CLAUDE_CODE_SESSION_ID with `code`.\n")

	native := runNativePreservingHostEnv("", "build",
		"--workflow-dir", root,
		"--entity-path", entityPath,
		"--stage", "backlog",
		"--checklist-file", checklistPath,
		"--scope-notes-file", scopePath,
	)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
	}
	var out struct {
		DispatchFilePath string `json:"dispatch_file_path"`
	}
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, native.stdout)
	}
	body := readDispatchBody(t, out.DispatchFilePath)
	for _, want := range []string{literalChecklist, "### Scope\nUse $CLAUDE_CODE_SESSION_ID with `code`."} {
		if !strings.Contains(body, want) {
			t.Fatalf("dispatch body missing literal %q:\n%s", want, body)
		}
	}
}

func TestBuildHostResolutionFromFlagAndEnv(t *testing.T) {
	t.Setenv("PI_CODING_AGENT", "")
	t.Setenv("PI_CODING_AGENT_DIR", "")

	t.Run("derived-codex", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		assertCodexFreshPrompt(t, out.Prompt, out.DispatchFilePath)
		body := readDispatchBody(t, out.DispatchFilePath)
		for _, banned := range []string{"Skill(skill=\"spacedock:ensign\")", "SendMessage(to=\"team-lead\""} {
			if strings.Contains(body, banned) {
				t.Fatalf("derived Codex dispatch body contains %q:\n%s", banned, body)
			}
		}
	})

	t.Run("derived-claude", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "1")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		if !strings.HasPrefix(out.Prompt, "Skill(skill=\"spacedock:ensign\")") {
			t.Fatalf("derived Claude prompt lost Skill wrapper: %q", out.Prompt)
		}
	})

	t.Run("derived-pi-from-PI_CODING_AGENT", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		t.Setenv("PI_CODING_AGENT", "true")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		assertPiBuildOutput(t, native.stdout)
	})

	t.Run("derived-pi-from-PI_CODING_AGENT_DIR", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		t.Setenv("PI_CODING_AGENT", "")
		t.Setenv("PI_CODING_AGENT_DIR", t.TempDir())
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		assertPiBuildOutput(t, native.stdout)
	})

	t.Run("host-flag", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-", "--host", "codex")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		assertCodexFreshPrompt(t, out.Prompt, out.DispatchFilePath)
	})

	t.Run("host-flag-overrides-pi-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		t.Setenv("PI_CODING_AGENT", "true")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-", "--host", "claude")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		if !strings.HasPrefix(out.Prompt, "Skill(skill=\"spacedock:ensign\")") {
			t.Fatalf("explicit Claude override lost Skill wrapper: %q", out.Prompt)
		}
	})

	t.Run("codex-flag-overrides-pi-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		t.Setenv("PI_CODING_AGENT", "true")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-", "--host", "codex")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		assertCodexFreshPrompt(t, out.Prompt, out.DispatchFilePath)
	})

	t.Run("unsupported-explicit-host", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-", "--host", "banana")
		assertBuildHostError(t, native, "unsupported host", "claude, codex, or pi")
	})

	t.Run("missing-source", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-")
		assertBuildHostError(t, native, "host source", "CODEX_THREAD_ID", "CLAUDECODE", "PI_CODING_AGENT", "PI_CODING_AGENT_DIR")
	})

	t.Run("ambiguous-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "1")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-")
		assertBuildHostError(t, native, "ambiguous", "CODEX_THREAD_ID", "CLAUDECODE")
	})

	t.Run("ambiguous-pi-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "")
		t.Setenv("PI_CODING_AGENT", "true")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-")
		assertBuildHostError(t, native, "ambiguous", "CODEX_THREAD_ID", "PI_CODING_AGENT", "--host claude, codex, or pi")
	})

	t.Run("explicit-overrides-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "")
		root, entity := buildHostFixture(t)

		native := runNativePreservingHostEnv("- a\n- b", "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", "-", "--host", "claude")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		if !strings.HasPrefix(out.Prompt, "Skill(skill=\"spacedock:ensign\")") {
			t.Fatalf("explicit Claude override lost Skill wrapper: %q", out.Prompt)
		}
	})
}

func assertPiBuildOutput(t *testing.T, stdout string) {
	t.Helper()
	out := decodeBuildOutput(t, stdout)
	for _, banned := range []string{"Skill(skill=", "Agent(", "SendMessage", "TeamCreate", "TeamDelete"} {
		if strings.Contains(out.Prompt, banned) {
			t.Fatalf("derived Pi prompt contains Claude syntax %q: %q", banned, out.Prompt)
		}
	}
	if !strings.Contains(out.Prompt, "Read ") || !strings.Contains(out.Prompt, "treat its content as your assignment") {
		t.Fatalf("derived Pi prompt should be the read-dispatch-file form: %q", out.Prompt)
	}
	body := readDispatchBody(t, out.DispatchFilePath)
	for _, want := range []string{"read this dispatch file", "Pi subagent completion result", "Do not emit Claude team-tool calls"} {
		if !strings.Contains(body, want) {
			t.Fatalf("derived Pi dispatch body missing %q:\n%s", want, body)
		}
	}
}

func buildHostFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	gitInit(t, root)
	return root, entityPath
}

func decodeBuildOutput(t *testing.T, stdout string) struct {
	DispatchFilePath string `json:"dispatch_file_path"`
	Prompt           string `json:"prompt"`
	Name             string `json:"name"`
} {
	t.Helper()
	var out struct {
		DispatchFilePath string `json:"dispatch_file_path"`
		Prompt           string `json:"prompt"`
		Name             string `json:"name"`
	}
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, stdout)
	}
	return out
}

func assertBuildHostError(t *testing.T, native runResult, wants ...string) {
	t.Helper()
	if native.exit == 0 {
		t.Fatalf("build unexpectedly exited 0\nstdout=%s\nstderr=%s", native.stdout, native.stderr)
	}
	for _, want := range wants {
		if !strings.Contains(native.stderr, want) {
			t.Fatalf("stderr missing %q:\n%s", want, native.stderr)
		}
	}
}

func TestChecklistSourcesHaveIdenticalLiteralSections(t *testing.T) {
	root, entity := buildHostFixture(t)
	const hazard = "  Keep \"quotes\", `ticks`, $(touch SHOULD_NOT_EXIST), $HOME, 雪 & <tag>  "
	scope, feedback := filepath.Join(t.TempDir(), "scope.md"), filepath.Join(t.TempDir(), "feedback.md")
	const scopeText = "### Scope\n`code` $HOME\n\n"
	const feedbackText = "REJECTED:\n  $(literal) 雪\n\n"
	writeFile(t, scope, scopeText)
	writeFile(t, feedback, feedbackText)
	for _, tc := range []struct{ name, input, want string }{
		{"one terminal CR", "x\r\r\n", "x\r"},
		{"EOF", hazard, hazard}, {"LF", hazard + "\n", hazard},
		{"CRLF blanks", "\r\n" + hazard + "\r\n\t \r\nsecond\r\n", hazard + "\nsecond"},
		{"long line", strings.Repeat("x", 70000), strings.Repeat("x", 70000)},
		{"JSON literal", `{"checklist":"ordinary text"}`, `{"checklist":"ordinary text"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var previous string
			for _, source := range []string{"-", filepath.Join(t.TempDir(), "checklist")} {
				input := tc.input
				if source != "-" {
					writeFile(t, source, input)
					input = "must not consume this"
				}
				result := runNative(input, "build", "--workflow-dir", root, "--entity-path", entity, "--stage", "backlog", "--checklist-file", source, "--scope-notes-file", scope, "--feedback-context-file", feedback)
				if result.exit != 0 {
					t.Fatalf("%s: %s", source, result.stderr)
				}
				body := readDispatchBody(t, dispatchFilePathFromStdout(t, result.stdout))
				expected := "### Completion checklist\n\n" + tc.want + "\n\n### Summary\n"
				if !strings.Contains(body, expected) {
					t.Fatal("checklist section bytes changed")
				}
				if !strings.Contains(body, "### Feedback from prior review\n\n"+feedbackText+"\n\n"+scopeText+"\n\n### Completion checklist") {
					t.Fatal("supporting prose bytes changed")
				}
				if previous != "" && body != previous {
					t.Fatal("file/stdin artifacts differ")
				}
				previous = body
			}
		})
	}
}
