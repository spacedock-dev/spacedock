// ABOUTME: Codex argv/conflict/support proofs plus an isolated-home lifecycle oracle.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

var wantCodexCollaborationLayer = []string{"-c", "agents.enabled=true", "-c", "features.multi_agent=true", "-c", `features.multi_agent_v2={max_concurrent_threads_per_session=16,tool_namespace="agents",hide_spawn_agent_metadata=false}`}

func wantCodexArgv(tail ...string) []string {
	return append(append([]string{"codex"}, wantCodexCollaborationLayer...), tail...)
}

func TestCodexCollaborationLayerCompleteArgv(t *testing.T) {
	checkout, _ := localPluginCheckout(t, "codex")
	tests := []struct {
		name, dir  string
		args, tail []string
	}{{"plain", "", nil, []string{"--ask-for-approval", "on-request", wantCodexBootstrapPrompt}}, {"local plugin", "", []string{"--plugin-dir", checkout}, []string{"--ask-for-approval", "on-request", wantCodexBootstrapPrompt}}, {"safehouse", "safehouse", nil, []string{"--dangerously-bypass-approvals-and-sandbox", wantCodexBootstrapPrompt}}, {"resume", "", []string{"--", "resume", "thread-123"}, []string{"resume", "thread-123"}}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CODEX_HOME", t.TempDir())
			host := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer
			dir, want := t.TempDir(), wantCodexArgv(tt.tail...)
			if tt.dir == "safehouse" {
				dir = safehouseFixtureDir(t)
				withExecutablePath(t, executableFixture(t), nil)
				want = append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--"}, want...)
			}
			code := runCodex(context.Background(), tt.args, dir, host, lookFound, &stdout, &stderr)
			if code != 0 || !slices.Equal(host.launchedArg, want) {
				t.Fatalf("exit=%d stderr=%q argv=%#v, want %#v", code, stderr.String(), host.launchedArg, want)
			}
		})
	}
}

func TestCodexRejectsOwnedCollaborationOverridesBeforeSideEffects(t *testing.T) {
	checkout, _ := localPluginCheckout(t, "codex")
	conflicts := [][]string{{"-c", "agents.enabled=false"}, {`-cagents.enabled=false`}, {`-c`, `"agents".enabled=false`}, {`-c"features"."multi_agent"=false`}, {`--config='features'.'multi_agent_v2'.tool_namespace="other"`}, {"-c=features.multi_agent=false"}, {"--config", "features.multi_agent_v2={}"}, {"--config=features.multi_agent_v2.max_concurrent_threads_per_session=2"}, {"-c", `features.multi_agent_v2.tool_namespace="collaboration"`}, {"-c", "features.multi_agent_v2.hide_spawn_agent_metadata=true"}, {"--enable", "multi_agent"}, {"--disable=multi_agent"}, {"--enable=multi_agent_v2"}, {"--disable", "multi_agent_v2"}}
	const wantDiagnostic = "spacedock codex: collaboration settings are managed by Spacedock; remove the forwarded override\n"

	for _, conflict := range conflicts {
		name := strings.NewReplacer("/", "_", "=", "_", " ", "_").Replace(strings.Join(conflict, "_"))
		t.Run(name, func(t *testing.T) {
			host := &fakeHost{manifest: ""}
			var stdout, stderr bytes.Buffer
			args := append([]string{"--plugin-dir", checkout, "--"}, conflict...)
			code := runCodex(context.Background(), args, t.TempDir(), host, lookFound, &stdout, &stderr)
			if code != 1 || stderr.String() != wantDiagnostic || len(host.installCmds) != 0 || host.launchedArg != nil {
				t.Fatalf("exit=%d stderr=%q install=%v launch=%v", code, stderr.String(), host.installCmds, host.launchedArg)
			}
		})
	}

	t.Run("unrelated config remains forwarded", func(t *testing.T) {
		host := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		forwarded := []string{"--config", `model_reasoning_effort="high"`}
		code := runCodex(context.Background(), append([]string{"--"}, forwarded...), t.TempDir(), host, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if !slices.Contains(host.launchedArg, forwarded[1]) {
			t.Fatalf("unrelated config was not forwarded: %v", host.launchedArg)
		}
	})
}

type rejectingCodexHost struct {
	fakeHost
	stderr *bytes.Buffer
}

func (h *rejectingCodexHost) Launch(argv []string, env []string) (int, error) {
	h.launchedArg = slices.Clone(argv)
	ownedLen := 1 + len(wantCodexCollaborationLayer)
	if len(argv) < ownedLen || !slices.Equal(argv[:ownedLen], wantCodexArgv()) {
		return 1, fmt.Errorf("owned collaboration layer missing from unsupported-host probe: %v", argv)
	}
	fmt.Fprintln(h.stderr, "error: unsupported configuration key: features.multi_agent_v2")
	return 78, nil
}

func TestCodexUnsupportedHostFailsBeforeSession(t *testing.T) {
	var stdout, stderr bytes.Buffer
	host := &rejectingCodexHost{fakeHost: fakeHost{manifest: compatibleManifest(t)}, stderr: &stderr}
	code := runCodex(context.Background(), nil, t.TempDir(), host, lookFound, &stdout, &stderr)
	if code != 78 {
		t.Fatalf("exit = %d, want native host exit 78", code)
	}
	if !strings.Contains(stderr.String(), "unsupported configuration key: features.multi_agent_v2") {
		t.Fatalf("native config diagnostic was not preserved: %q", stderr.String())
	}
}

func TestCodexIsolatedHomeCollaborationLifecycle(t *testing.T) {
	if os.Getenv("SPACEDOCK_LIVE_CODEX_MULTI_AGENT") != "1" {
		t.Skip("set SPACEDOCK_LIVE_CODEX_MULTI_AGENT=1 to spend one short live Codex turn")
	}
	if _, err := exec.LookPath("codex"); err != nil {
		t.Skip("codex is not on PATH")
	}
	auth, err := os.ReadFile(filepath.Join(codexHome(), "auth.json"))
	if err != nil {
		t.Skipf("read Codex authentication: %v", err)
	}
	isolatedHome := t.TempDir()
	if err := os.WriteFile(filepath.Join(isolatedHome, "auth.json"), auth, 0o600); err != nil {
		t.Fatal(err)
	}
	spacedockBin := filepath.Join(t.TempDir(), "spacedock")
	build := exec.Command("go", "build", "-o", spacedockBin, "./cmd/spacedock")
	build.Dir = repoRootForDevBuild(t)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build spacedock front door: %v\n%s", err, out)
	}
	ready, done, parent := "READY_7E91", "DONE_7E91", "PARENT_7E91"
	prompt := fmt.Sprintf("Use native collaboration tools only: spawn exactly one worker and ask it to reply %s then wait for follow-up; wait for it; follow up to the same worker asking it to run `sleep 5` before replying %s; immediately list agents; wait for that worker; then reply exactly %s.", ready, done, parent)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, spacedockBin, "codex", "--skip-compat-check", prompt, "--", "exec", "--json", "--skip-git-repo-check")
	cmd.Env = append(withoutEnv(os.Environ(), "CODEX_HOME"), "CODEX_HOME="+isolatedHome)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("isolated Codex lifecycle failed: %v\n%s", err, out)
	}
	sessions := map[string][]byte{}
	_ = filepath.WalkDir(isolatedHome, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr == nil && !entry.IsDir() && strings.HasSuffix(path, ".jsonl") {
			if session, readErr := os.ReadFile(path); readErr == nil {
				base := strings.TrimSuffix(entry.Name(), ".jsonl")
				if len(base) >= 36 {
					sessions[base[len(base)-36:]] = session
				}
			}
		}
		return nil
	})
	if err := gradeCodexLifecycle(out, sessions, ready, done, parent); err != nil {
		t.Fatalf("structured launcher lifecycle: %v\nstdout:\n%s", err, out)
	}

	disabled := exec.Command("codex", "-c", "agents.enabled=false", "exec", "--json", "--skip-git-repo-check", "Attempt one native worker spawn. If the tool is unavailable, reply exactly DISABLED_CONTROL.")
	disabled.Env = cmd.Env
	disabledOut, err := disabled.Output()
	if err != nil {
		t.Fatalf("disabled control: %v\n%s", err, disabledOut)
	}
	if n := countCodexCollaborationEvents(disabledOut); n != 0 {
		t.Fatalf("disabled control emitted %d collaboration events:\n%s", n, disabledOut)
	}
}

type codexStreamEvent struct {
	Type     string
	ThreadID string `json:"thread_id"`
	Item     struct{ Type, Text string }
}
type codexSessionEvent struct {
	Type    string
	Payload struct {
		Type, Role, Name, Arguments, Output string
		Version                             string `json:"multi_agent_version"`
		Content                             []struct{ Type, Text string }
	}
}

func gradeCodexLifecycle(stdout []byte, sessions map[string][]byte, ready, done, parent string) error {
	var parentID string
	parentDone := false
	for _, line := range bytes.Split(stdout, []byte("\n")) {
		var event codexStreamEvent
		if json.Unmarshal(line, &event) == nil {
			if event.Type == "thread.started" {
				parentID = event.ThreadID
			}
			parentDone = parentDone || event.Type == "item.completed" && event.Item.Type == "agent_message" && event.Item.Text == parent
		}
	}
	wantTools := []string{"spawn_agent", "wait_agent", "followup_task", "list_agents", "wait_agent"}
	var tools []string
	var target, lastCall string
	listSawTarget := false
	for _, line := range bytes.Split(sessions[parentID], []byte("\n")) {
		var event codexSessionEvent
		if json.Unmarshal(line, &event) != nil || event.Type != "response_item" {
			continue
		}
		if event.Payload.Type == "function_call" && slices.Contains(wantTools, event.Payload.Name) {
			tools, lastCall = append(tools, event.Payload.Name), event.Payload.Name
			if event.Payload.Name == "followup_task" {
				var args struct{ Target string }
				_ = json.Unmarshal([]byte(event.Payload.Arguments), &args)
				if args.Target != target {
					return fmt.Errorf("follow-up target %q != spawned worker %q", args.Target, target)
				}
			}
		}
		if event.Payload.Type == "function_call_output" {
			if lastCall == "spawn_agent" {
				var result struct {
					TaskName string `json:"task_name"`
				}
				_ = json.Unmarshal([]byte(event.Payload.Output), &result)
				target = result.TaskName
			}
			if lastCall == "list_agents" {
				listSawTarget = strings.Contains(event.Payload.Output, `"agent_name":"`+target+`"`)
			}
		}
	}
	var child []byte
	for id, session := range sessions {
		if id != parentID && bytes.Contains(session, []byte(`"parent_thread_id":"`+parentID+`"`)) && bytes.Contains(session, []byte(`"agent_path":"`+target+`"`)) {
			child = session
		}
	}
	if !slices.Equal(tools, wantTools) || target == "" || !listSawTarget || !parentDone || !sessionHasV2(sessions[parentID]) || !sessionHasAssistantTexts(child, ready, done) {
		return fmt.Errorf("tools=%v target=%q listed=%t parent-complete=%t v2/child outputs=%t", tools, target, listSawTarget, parentDone, sessionHasV2(sessions[parentID]) && sessionHasAssistantTexts(child, ready, done))
	}
	return nil
}

func countCodexCollaborationEvents(stdout []byte) (count int) {
	for _, line := range bytes.Split(stdout, []byte("\n")) {
		var event codexStreamEvent
		if json.Unmarshal(line, &event) == nil && event.Item.Type == "collab_tool_call" {
			count++
		}
	}
	return
}

func sessionHasV2(session []byte) bool {
	for _, line := range bytes.Split(session, []byte("\n")) {
		var event codexSessionEvent
		if json.Unmarshal(line, &event) == nil && event.Type == "turn_context" && event.Payload.Version == "v2" {
			return true
		}
	}
	return false
}

func sessionHasAssistantTexts(session []byte, want ...string) bool {
	next := 0
	for _, line := range bytes.Split(session, []byte("\n")) {
		var event codexSessionEvent
		if json.Unmarshal(line, &event) == nil && event.Type == "response_item" && event.Payload.Type == "message" && event.Payload.Role == "assistant" {
			for _, content := range event.Payload.Content {
				if next < len(want) && content.Type == "output_text" && strings.TrimSpace(content.Text) == want[next] {
					next++
				}
			}
		}
	}
	return next == len(want)
}

func TestCodexLifecycleOracleRejectsVocabularyOnly(t *testing.T) {
	fake := []byte(`{"type":"item.completed","item":{"type":"agent_message","text":"spawn_agent followup_task list_agents wait_agent READY_7E91 DONE_7E91 PARENT_7E91 multi_agent_version v2"}}`)
	if err := gradeCodexLifecycle(fake, nil, "READY_7E91", "DONE_7E91", "PARENT_7E91"); err == nil {
		t.Fatal("vocabulary-only fake host passed the structured lifecycle oracle")
	}
}
