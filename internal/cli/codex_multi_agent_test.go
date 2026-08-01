// ABOUTME: Codex collaboration launch-layer argv, conflict, and host-support proofs.
// ABOUTME: The opt-in live test grades an isolated-home native worker lifecycle.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var wantCodexCollaborationLayer = []string{
	"-c", "agents.enabled=true",
	"-c", "features.multi_agent=true",
	"-c", `features.multi_agent_v2={max_concurrent_threads_per_session=16,tool_namespace="agents",hide_spawn_agent_metadata=false}`,
}

func wantCodexArgv(tail ...string) []string {
	return append(append([]string{"codex"}, wantCodexCollaborationLayer...), tail...)
}

func TestCodexCollaborationLayerCompleteArgv(t *testing.T) {
	checkout, _ := localPluginCheckout(t, "codex")
	tests := []struct {
		name string
		args []string
		dir  func(*testing.T) string
		want func(*testing.T) []string
	}{
		{
			name: "plain",
			dir:  func(t *testing.T) string { return t.TempDir() },
			want: func(*testing.T) []string {
				return wantCodexArgv("--ask-for-approval", "on-request", wantCodexBootstrapPrompt)
			},
		},
		{
			name: "local plugin",
			args: []string{"--plugin-dir", checkout},
			dir:  func(t *testing.T) string { return t.TempDir() },
			want: func(*testing.T) []string {
				return wantCodexArgv("--ask-for-approval", "on-request", wantCodexBootstrapPrompt)
			},
		},
		{
			name: "safehouse",
			dir:  safehouseFixtureDir,
			want: func(t *testing.T) []string {
				bin := executableFixture(t)
				withExecutablePath(t, bin, nil)
				return append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--"},
					wantCodexArgv("--dangerously-bypass-approvals-and-sandbox", wantCodexBootstrapPrompt)...)
			},
		},
		{
			name: "resume",
			args: []string{"--", "resume", "thread-123"},
			dir:  func(t *testing.T) string { return t.TempDir() },
			want: func(*testing.T) []string { return wantCodexArgv("resume", "thread-123") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CODEX_HOME", t.TempDir())
			host := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer
			want := tt.want(t)
			code := runCodex(context.Background(), tt.args, tt.dir(t), host, lookFound, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			if !slices.Equal(host.launchedArg, want) {
				t.Fatalf("launch argv = %#v, want %#v", host.launchedArg, want)
			}
		})
	}
}

func TestCodexRejectsOwnedCollaborationOverridesBeforeSideEffects(t *testing.T) {
	checkout, _ := localPluginCheckout(t, "codex")
	conflicts := [][]string{
		{"-c", "agents.enabled=false"},
		{"-c=features.multi_agent=false"},
		{"--config", "features.multi_agent_v2={}"},
		{"--config=features.multi_agent_v2.max_concurrent_threads_per_session=2"},
		{"-c", `features.multi_agent_v2.tool_namespace="collaboration"`},
		{"-c", "features.multi_agent_v2.hide_spawn_agent_metadata=true"},
		{"--enable", "multi_agent"},
		{"--disable=multi_agent"},
		{"--enable=multi_agent_v2"},
		{"--disable", "multi_agent_v2"},
	}
	const wantDiagnostic = "spacedock codex: collaboration settings are managed by Spacedock; remove the forwarded override\n"

	for _, conflict := range conflicts {
		name := strings.NewReplacer("/", "_", "=", "_", " ", "_").Replace(strings.Join(conflict, "_"))
		t.Run(name, func(t *testing.T) {
			host := &fakeHost{manifest: ""}
			var stdout, stderr bytes.Buffer
			args := append([]string{"--plugin-dir", checkout, "--"}, conflict...)
			code := runCodex(context.Background(), args, t.TempDir(), host, lookFound, &stdout, &stderr)
			if code != 1 {
				t.Fatalf("exit = %d, want 1", code)
			}
			if stderr.String() != wantDiagnostic {
				t.Fatalf("stderr = %q, want %q", stderr.String(), wantDiagnostic)
			}
			if len(host.installCmds) != 0 || host.launchedArg != nil {
				t.Fatalf("conflict caused side effects: install=%v launch=%v", host.installCmds, host.launchedArg)
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
	stderr         *bytes.Buffer
	sessionStarted bool
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
	if host.sessionStarted {
		t.Fatal("unsupported host opened an interactive session")
	}
	if !strings.Contains(stderr.String(), "unsupported configuration key: features.multi_agent_v2") {
		t.Fatalf("native config diagnostic was not preserved: %q", stderr.String())
	}
}

func TestCodexIsolatedHomeCollaborationLifecycle(t *testing.T) {
	if os.Getenv("SPACEDOCK_LIVE_CODEX_MULTI_AGENT") != "1" {
		t.Skip("set SPACEDOCK_LIVE_CODEX_MULTI_AGENT=1 to spend one short live Codex turn")
	}
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("codex is not on PATH")
	}
	authSource := filepath.Join(codexHome(), "auth.json")
	auth, err := os.ReadFile(authSource)
	if err != nil {
		t.Skipf("read Codex authentication: %v", err)
	}
	isolatedHome := t.TempDir()
	if err := os.WriteFile(filepath.Join(isolatedHome, "auth.json"), auth, 0o600); err != nil {
		t.Fatal(err)
	}
	prompt := "You must complete every step with native collaboration tools before your final answer: (1) spawn one worker whose task is to reply CHILD_READY and wait for follow-up; (2) wait for that exact worker to reply CHILD_READY; (3) follow up to the same worker and require CHILD_DONE; (4) list workers; (5) wait for that worker to finish; (6) only after observing CHILD_READY and CHILD_DONE, print PARENT_DONE. Do not print PARENT_DONE if any tool or marker is unavailable."
	args := append([]string{"exec", "--json", "--skip-git-repo-check"}, wantCodexCollaborationLayer...)
	args = append(args, prompt)
	cmd := exec.Command(codexBin, args...)
	cmd.Env = append(withoutEnv(os.Environ(), "CODEX_HOME"), "CODEX_HOME="+isolatedHome)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("isolated Codex lifecycle failed: %v\n%s", err, out)
	}
	transcript := string(out)
	_ = filepath.WalkDir(isolatedHome, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr == nil && !entry.IsDir() && strings.HasSuffix(path, ".jsonl") {
			if session, readErr := os.ReadFile(path); readErr == nil {
				transcript += "\n" + string(session)
			}
		}
		return nil
	})
	required := map[string][]string{
		"spawn":                {`"tool":"spawn"`, "spawn_agent"},
		"same-worker followup": {`"tool":"send_input"`, "followup_task"},
		"list":                 {`"tool":"list"`, "list_agents"},
		"wait":                 {`"tool":"wait"`, "wait_agent"},
		"child ready":          {"CHILD_READY"},
		"child done":           {"CHILD_DONE"},
		"parent done":          {"PARENT_DONE"},
		"v2 turn context":      {`"multi_agent_version":"v2"`, `"multi_agent_version": "v2"`},
	}
	for label, alternatives := range required {
		if !slices.ContainsFunc(alternatives, func(value string) bool { return strings.Contains(transcript, value) }) {
			t.Fatalf("isolated lifecycle omitted %s (%v); transcript:\n%s", label, alternatives, transcript)
		}
	}
}
