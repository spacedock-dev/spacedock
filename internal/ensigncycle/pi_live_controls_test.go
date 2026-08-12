package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const piLiveSmokeMarker = "PI-LIVE-SUBAGENT-ENSIGN-SMOKE"

type piSmokeEnvelope struct {
	Agent        string `json:"agent"`
	Skill        string `json:"skill"`
	Prompt       string `json:"prompt"`
	DispatchFile string `json:"dispatch_file_path"`
}

func piLiveSmokePrompt(repo, workflowRoot, stateRoot, entityPath string, envelope piSmokeEnvelope) string {
	return fmt.Sprintf(`You are the Spacedock first officer for a live Pi smoke test.

An initial-dispatch artifact was assembled for the entity with `+"`spacedock dispatch build --host pi`"+`; forward it through pi-subagents exactly as emitted — this smoke exists to prove the build artifact's spawn fields drive the worker boot.

  agent: %[5]s
  skill: %[6]s
  task: %[7]s

Use the pi-subagents subagent(...) tool exactly once with those fields verbatim (context must be "fresh", working directory %[2]s). Do not use or mention Claude Agent, SendMessage, TeamCreate, or TeamDelete tools. Do not paraphrase, re-order, or extend the task string.

After subagent(...) returns, you as first officer must verify the entity file %[4]s contains %[8]s and verify the state checkout %[3]s git log contains 'ensign: pi live smoke' over pi-live-smoke/index.md. The stage report must use the exact heading '## Stage Report: implementation' — confirm that too. Exit successfully only after those durable checks pass; your final message names the agent and skill values you passed to subagent(...) and the child's run id.

Reference paths: ensign contract at %[1]s/skills/ensign/SKILL.md; Pi ensign adapter at %[1]s/skills/ensign/references/pi-ensign-runtime.md (the worker's dispatch artifact already points at them).`,
		repo, workflowRoot, stateRoot, entityPath, envelope.Agent, envelope.Skill, envelope.Prompt, piLiveSmokeMarker)
}

func TestPiLiveSmokePromptRequiresExactStageReportHeading(t *testing.T) {
	envelope := piSmokeEnvelope{Agent: "worker", Skill: "ensign", Prompt: "Read /tmp/spacedock-dispatch/x.md and treat its content as your assignment."}
	prompt := piLiveSmokePrompt("/repo", "/workflow", "/workflow/.spacedock-state", "/workflow/.spacedock-state/pi-live-smoke/index.md", envelope)
	want := "exact heading '## Stage Report: implementation'"
	if !strings.Contains(prompt, want) {
		t.Fatalf("pi live smoke prompt missing %q:\n%s", want, prompt)
	}
	source := readFile(t, "pi_live_runner_test.go")
	for _, contract := range []string{`filepath.Join(repoRoot(t), "skills", "ensign", "SKILL.md")`, `filepath.Join(repoRoot(t), "skills", "ensign", "references", "pi-ensign-runtime.md")`} {
		if !strings.Contains(source, contract) {
			t.Fatalf("Pi live smoke checklist missing boot-contract action %s", contract)
		}
	}
}

func piLiveEnv(piHome, sessionDir, cleanHome, binaryDir, piSubagentsRoot string) []string {
	env := cleanEnviron(
		"CODEX_THREAD_ID", "CLAUDECODE", "HOME", "PI_CODING_AGENT_DIR",
		"PI_CODING_AGENT_SESSION_DIR", "PI_INTERCOM_PACKAGE_ROOT",
		"PI_SUBAGENTS_PACKAGE_ROOT", "PI_OFFLINE", "OPENAI_API_KEY", "PI_OPENAI_CODEX_AUTH_JSON",
	)
	env = dropEnvPrefix(env, "PI_SUBAGENT_")
	env = append(env,
		"HOME="+cleanHome,
		"PI_CODING_AGENT_DIR="+piHome,
		"PI_CODING_AGENT_SESSION_DIR="+sessionDir,
		"PI_INTERCOM_PACKAGE_ROOT="+piIntercomPackageRoot(piSubagentsRoot),
		"PI_SUBAGENTS_PACKAGE_ROOT="+piSubagentsRoot,
		"PI_OFFLINE=1",
	)
	return withBinaryOnPath(env, filepath.Join(binaryDir, "spacedock"))
}

func piLiveEnvForAuth(piHome, sessionDir, cleanHome, binaryDir, piSubagentsRoot, openAIKey, mode string) []string {
	env := piLiveEnv(piHome, sessionDir, cleanHome, binaryDir, piSubagentsRoot)
	if mode == piAuthOAuth {
		env = withoutPiEnvKey(env, "OPENAI_API_KEY")
	} else if openAIKey != "" {
		env = append(env, "OPENAI_API_KEY="+openAIKey)
	}
	return env
}

func TestPiLiveAuthSelectionAndSeeding(t *testing.T) {
	oauth := `{"type":"oauth","access":"sentinel","refresh":"refresh"}`
	if got := decidePiLiveAuth(oauth, "key", ""); got.mode != piAuthOAuth || got.model != "openai-codex/gpt-5.6-luna:max" {
		t.Fatalf("OAuth decision = %#v", got)
	}
	home := t.TempDir()
	if err := seedPiOAuthAuth(home, oauth); err != nil {
		t.Fatal(err)
	}
	if mode := fileMode(t, filepath.Join(home, "auth.json")); mode != 0o600 {
		t.Fatalf("mode = %o", mode)
	}
	if got := readFile(t, filepath.Join(home, "auth.json")); !strings.Contains(got, `"openai-codex"`) || !strings.Contains(got, "sentinel") {
		t.Fatalf("seeded auth = %s", got)
	}
	if got := decidePiLiveAuth("", "key", ""); got.mode != piAuthAPIKey || got.model != "openai/gpt-5.6-luna:max" {
		t.Fatalf("key decision = %#v", got)
	}
}

func TestPiLiveEnvDropsForeignRuntimeMarkers(t *testing.T) {
	for key, value := range map[string]string{"CODEX_THREAD_ID": "codex", "CLAUDECODE": "claude", "PI_CODING_AGENT": "pi", "PI_CODING_AGENT_DIR": "/parent/pi",
		"PI_CODING_AGENT_SESSION_DIR": "/parent/sessions", "PI_INTERCOM_PACKAGE_ROOT": "/parent/intercom",
		"PI_SUBAGENTS_PACKAGE_ROOT": "/parent/package",
		"PI_OFFLINE":                "0", "HOME": "/parent/home", "OPENAI_API_KEY": "key", "PATH": "/parent/bin"} {
		t.Setenv(key, value)
	}
	env := piLiveEnv("/target/pi", "/target/sessions", "/target/home", "/spacedock/bin", "/target/package")
	want := map[string]string{"CODEX_THREAD_ID": "", "CLAUDECODE": "",
		"PI_CODING_AGENT": "pi", "PI_CODING_AGENT_DIR": "/target/pi",
		"PI_CODING_AGENT_SESSION_DIR": "/target/sessions", "PI_SUBAGENTS_PACKAGE_ROOT": "/target/package",
		"PI_INTERCOM_PACKAGE_ROOT": "/parent/intercom",
		"PI_OFFLINE":               "1", "HOME": "/target/home", "OPENAI_API_KEY": "",
		"PATH": "/spacedock/bin" + string(os.PathListSeparator) + "/parent/bin"}
	for key, value := range want {
		assertEnvValue(t, env, key, value)
	}
}

func dropEnvPrefix(env []string, prefix string) []string {
	kept := env[:0]
	for _, kv := range env {
		key := kv
		if i := strings.IndexByte(kv, '='); i >= 0 {
			key = kv[:i]
		}
		if strings.HasPrefix(key, prefix) {
			continue
		}
		kept = append(kept, kv)
	}
	return kept
}

func TestPiLiveEnvScrubsAmbientPiSubagentMarkers(t *testing.T) {
	t.Setenv("PI_SUBAGENT_CHILD", "1")
	t.Setenv("PI_SUBAGENT_RUN_ID", "ambient-run")
	t.Setenv("PI_SUBAGENT_DEPTH", "1")
	env := piLiveEnv("/target/pi", "/target/sessions", "/target/home", "/spacedock/bin", "/target/package")
	for _, kv := range env {
		key := kv
		if i := strings.IndexByte(kv, '='); i >= 0 {
			key = kv[:i]
		}
		if strings.HasPrefix(key, "PI_SUBAGENT_") {
			t.Fatalf("piLiveEnv leaked ambient marker %s", kv)
		}
	}
}

func piIntercomPackageRoot(piSubagentsRoot string) string {
	if p := os.Getenv("PI_INTERCOM_PACKAGE_ROOT"); p != "" {
		return p
	}
	return filepath.Join(filepath.Dir(piSubagentsRoot), "pi-intercom")
}

func TestPiIntercomPackageRootDefaultsBesideSubagents(t *testing.T) {
	t.Setenv("PI_INTERCOM_PACKAGE_ROOT", "")
	if got := piIntercomPackageRoot("/packages/pi-subagents"); got != "/packages/pi-intercom" {
		t.Fatalf("piIntercomPackageRoot() = %q, want sibling package", got)
	}
}
