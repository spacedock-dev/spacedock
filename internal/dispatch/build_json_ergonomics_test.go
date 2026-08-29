// ABOUTME: dispatch build flag/file input and host-resolution ergonomics.
package dispatch

import (
	"encoding/json"
	"os"
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

func TestBuildHostResolutionFromFlagJSONAndEnv(t *testing.T) {
	t.Setenv("PI_CODING_AGENT", "")
	t.Setenv("PI_CODING_AGENT_DIR", "")

	t.Run("derived-codex", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root)
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
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root)
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
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root)
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
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root)
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		assertPiBuildOutput(t, native.stdout)
	})

	t.Run("host-flag", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root, "--host", "codex")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		assertCodexFreshPrompt(t, out.Prompt, out.DispatchFilePath)
	})

	t.Run("matching-explicit-sources", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, map[string]any{"host": "codex"}), "build", "--workflow-dir", root, "--host", "codex")
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
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root, "--host", "claude")
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		if !strings.HasPrefix(out.Prompt, "Skill(skill=\"spacedock:ensign\")") {
			t.Fatalf("explicit Claude override lost Skill wrapper: %q", out.Prompt)
		}
	})

	t.Run("json-host-overrides-pi-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		t.Setenv("PI_CODING_AGENT", "true")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, map[string]any{"host": "codex"}), "build", "--workflow-dir", root)
		if native.exit != 0 {
			t.Fatalf("build exit=%d stderr=%s", native.exit, native.stderr)
		}
		out := decodeBuildOutput(t, native.stdout)
		assertCodexFreshPrompt(t, out.Prompt, out.DispatchFilePath)
	})

	t.Run("conflicting-explicit-sources", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, map[string]any{"host": "codex"}), "build", "--workflow-dir", root, "--host", "claude")
		assertBuildHostError(t, native, "--host", "JSON host")
	})

	t.Run("unsupported-explicit-host", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root, "--host", "banana")
		assertBuildHostError(t, native, "unsupported host", "claude, codex, or pi")
	})

	t.Run("missing-source", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "")
		t.Setenv("CLAUDECODE", "")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root)
		assertBuildHostError(t, native, "host source", "CODEX_THREAD_ID", "CLAUDECODE", "PI_CODING_AGENT", "PI_CODING_AGENT_DIR")
	})

	t.Run("ambiguous-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "1")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root)
		assertBuildHostError(t, native, "ambiguous", "CODEX_THREAD_ID", "CLAUDECODE")
	})

	t.Run("ambiguous-pi-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "")
		t.Setenv("PI_CODING_AGENT", "true")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root)
		assertBuildHostError(t, native, "ambiguous", "CODEX_THREAD_ID", "PI_CODING_AGENT", "--host claude, codex, or pi")
	})

	t.Run("explicit-overrides-runtime", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "")
		root, _ := buildHostFixture(t)

		native := runNativePreservingHostEnv(buildHostStdin(t, root, nil), "build", "--workflow-dir", root, "--host", "claude")
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

func TestBuildSchemaAndValidateOnly(t *testing.T) {
	t.Setenv("PI_CODING_AGENT", "")
	t.Setenv("PI_CODING_AGENT_DIR", "")

	t.Run("print-schema", func(t *testing.T) {
		native := runNativePreservingHostEnv("", "build", "--print-schema")
		if native.exit != 0 {
			t.Fatalf("print-schema exit=%d stderr=%s", native.exit, native.stderr)
		}
		var schema map[string]any
		if err := json.Unmarshal([]byte(native.stdout), &schema); err != nil {
			t.Fatalf("schema is not valid JSON: %v\n%s", err, native.stdout)
		}
		props, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("schema missing properties object:\n%s", native.stdout)
		}
		host, ok := props["host"].(map[string]any)
		if !ok {
			t.Fatalf("schema missing host property:\n%s", native.stdout)
		}
		if got := strings.Join(anyStrings(host["enum"]), ","); got != "claude,codex,pi" {
			t.Fatalf("host enum = %q, want claude,codex,pi", got)
		}
		if containsAnyString(schema["required"], "host") {
			t.Fatalf("host must be optional in schema required list:\n%s", native.stdout)
		}
	})

	t.Run("validate-only-derived-host-does-not-write-dispatch", func(t *testing.T) {
		t.Setenv("CODEX_THREAD_ID", "codex-thread")
		t.Setenv("CLAUDECODE", "")
		root, _ := buildHostFixture(t)
		reqPath := filepath.Join(root, "request.json")
		writeFile(t, reqPath, buildHostStdin(t, root, nil))
		dispatchPath := filepath.Join(dispatchFileDir, "spacedock-ensign-thing-backlog.md")
		if err := os.Remove(dispatchPath); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}

		native := runNativePreservingHostEnv("", "build", "--validate-only", reqPath)
		if native.exit != 0 {
			t.Fatalf("validate-only exit=%d stderr=%s", native.exit, native.stderr)
		}
		if native.stdout != "" {
			t.Fatalf("validate-only success should not emit dispatch JSON, got %q", native.stdout)
		}
		if _, err := os.Stat(dispatchPath); !os.IsNotExist(err) {
			t.Fatalf("validate-only wrote deterministic dispatch file %s (stat err=%v)", dispatchPath, err)
		}
	})

	t.Run("validate-only-errors", func(t *testing.T) {
		cases := []struct {
			name  string
			env   map[string]string
			extra map[string]any
			body  string
			want  []string
		}{
			{name: "malformed-json", body: `{"schema_version":`, want: []string{"invalid JSON"}},
			{name: "missing-host", env: map[string]string{"CODEX_THREAD_ID": "", "CLAUDECODE": ""}, want: []string{"host source"}},
			{name: "ambiguous-env", env: map[string]string{"CODEX_THREAD_ID": "codex-thread", "CLAUDECODE": "1"}, want: []string{"ambiguous", "CODEX_THREAD_ID", "CLAUDECODE"}},
			{name: "unsupported-host", extra: map[string]any{"host": "banana"}, want: []string{"unsupported host"}},
			{name: "empty-checklist", extra: map[string]any{"host": "claude", "checklist": []string{}}, want: []string{"checklist must not be empty"}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				for _, key := range []string{"CODEX_THREAD_ID", "CLAUDECODE", "PI_CODING_AGENT", "PI_CODING_AGENT_DIR"} {
					t.Setenv(key, "")
				}
				for key, value := range tc.env {
					t.Setenv(key, value)
				}
				root, _ := buildHostFixture(t)
				body := tc.body
				if body == "" {
					body = buildHostStdin(t, root, tc.extra)
				}
				reqPath := filepath.Join(root, "request.json")
				writeFile(t, reqPath, body)

				native := runNativePreservingHostEnv("", "build", "--validate-only", reqPath)
				assertBuildHostError(t, native, tc.want...)
			})
		}
	})
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

func buildHostStdin(t *testing.T, root string, extra map[string]any) string {
	t.Helper()
	entityPath := filepath.Join(root, "thing.md")
	return mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "backlog",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
	}, extra)
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

func anyStrings(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func containsAnyString(v any, want string) bool {
	for _, s := range anyStrings(v) {
		if s == want {
			return true
		}
	}
	return false
}
