//go:build live

// ABOUTME: Shared package-private isolated Paseo daemon, JSON runner, and split-root
// ABOUTME: fixture helpers consumed by Paseo live dispatch and adoption smoke tests.
package ensigncycle

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/cli"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

type paseoDaemonConfig struct{ artifactSubdir string }
type paseoDaemonOption func(*paseoDaemonConfig)

type paseoDaemon struct {
	Bin, Home, CanonicalRoot, ListenAddr, SpacedockBin string
	FOProvider, WorkerProvider, Thinking, CodexHome    string
	ArtifactDir                                        string
	Env                                                []string
	Timeout                                            time.Duration
}

type paseoSplitRootFixture struct {
	RepoRoot, WorkflowDir, StateRoot, EntityPath, Slug string
}

func startIsolatedPaseoDaemon(t *testing.T, opts ...paseoDaemonOption) *paseoDaemon {
	t.Helper()
	cfg := paseoDaemonConfig{artifactSubdir: "paseo-dispatch-smoke"}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	required := os.Getenv("SPACEDOCK_PASEO_LIVE_REQUIRED")
	paseoBin := strings.TrimSpace(os.Getenv("SPACEDOCK_PASEO_BIN"))
	var binErr error
	if paseoBin == "" {
		paseoBin, binErr = exec.LookPath("paseo")
	} else if _, err := os.Stat(paseoBin); err != nil {
		binErr = err
	}
	workerProvider := envOr("SPACEDOCK_PASEO_LIVE_PROVIDER", defaultPaseoLiveProvider())
	foProvider := envOr("SPACEDOCK_PASEO_LIVE_FO_PROVIDER", workerProvider)
	thinking := strings.TrimSpace(os.Getenv("SPACEDOCK_PASEO_LIVE_THINKING"))

	realHome, oauthJSON, openAIKey := os.Getenv("HOME"), os.Getenv("CODEX_AUTH_JSON"), os.Getenv("OPENAI_API_KEY")
	needsCodex := strings.HasPrefix(workerProvider, "codex/") || strings.HasPrefix(foProvider, "codex/")
	codexAuth := decideCodexLiveAuthForCI(oauthJSON, openAIKey, codexLocalAuthAvailable(realHome), required)
	checkProvider := workerProvider
	if needsCodex {
		checkProvider = "codex/check"
	}
	switch mode, msg := decidePaseoLivePreconditions(binErr == nil, checkProvider, codexAuth, required); mode {
	case "skip":
		t.Skip(msg)
	case "fatal":
		t.Fatal(msg)
	}
	repo := repoRoot(t)
	spacedockBin := piSpacedockBinary(t, repo)
	canonicalRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("canonicalize temp root: %v", err)
	}
	paseoHome := filepath.Join(canonicalRoot, "paseo-home")
	mustMkdirAll(t, paseoHome)
	artifactDir, listenAddr := codexLiveArtifactDir(t, cfg.artifactSubdir), pickFreePaseoListenAddr(t)
	var codexHome, envKey string
	if needsCodex {
		codexHome = setupIsolatedPaseoCodexHome(t, repo, artifactDir, realHome, oauthJSON, openAIKey, codexAuth)
		if codexAuth.mode == codexAuthAPIKey {
			envKey = openAIKey
		}
	}
	timeout := 12 * time.Minute
	if mins, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SPACEDOCK_PASEO_LIVE_TIMEOUT_MINUTES"))); err == nil && mins > 0 {
		timeout = time.Duration(mins) * time.Minute
	}
	d := &paseoDaemon{
		Bin: paseoBin, Home: paseoHome, CanonicalRoot: canonicalRoot, ListenAddr: listenAddr,
		SpacedockBin: spacedockBin, FOProvider: foProvider, WorkerProvider: workerProvider,
		Thinking: thinking, CodexHome: codexHome, ArtifactDir: artifactDir,
		Env: paseoLiveEnv(paseoHome, spacedockBin, codexHome, envKey), Timeout: timeout,
	}
	for _, kv := range paseoDaemonConfigEntries(listenAddr) {
		d.runJSON(t, nil, "daemon", "config", "set", kv[0], kv[1], "--json")
	}
	d.runJSON(t, nil, "daemon", "start", "--timeout", "30", "--json")
	t.Cleanup(func() {
		_, _ = runPaseoJSON(d.Bin, d.Home, d.Env, nil, "daemon", "stop", "--force", "--timeout", "5", "--json")
		var st struct{ LocalDaemon string }
		if _, err := runPaseoJSON(d.Bin, d.Home, d.Env, &st, "status", "--json"); err == nil && st.LocalDaemon != "stopped" {
			t.Errorf("isolated paseo daemon localDaemon = %q after stop, want stopped", st.LocalDaemon)
		}
		if t.Failed() {
			if data, err := os.ReadFile(filepath.Join(d.Home, "paseo.log")); err == nil {
				_ = os.WriteFile(filepath.Join(d.ArtifactDir, "paseo.log"), data, 0o644)
			}
		}
	})
	deadline := time.Now().Add(30 * time.Second)
	for {
		var st struct {
			ConnectedDaemon string
			WorkerPid       int
		}
		if _, err := runPaseoJSON(d.Bin, d.Home, d.Env, &st, "status", "--json"); err == nil && st.ConnectedDaemon == "reachable" && st.WorkerPid > 0 {
			return d
		}
		if time.Now().After(deadline) {
			t.Fatalf("isolated Paseo daemon at %s did not become reachable within 30s", d.ListenAddr)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (d *paseoDaemon) runJSON(t *testing.T, dst any, args ...string) string {
	t.Helper()
	stderr, err := runPaseoJSON(d.Bin, d.Home, d.Env, dst, args...)
	if err != nil {
		t.Fatal(err)
	}
	return stderr
}

func (d *paseoDaemon) stageNonCodexSkillsIfNeeded(t *testing.T, repoRoot string) {
	t.Helper()
	if strings.HasPrefix(d.FOProvider, "codex/") && strings.HasPrefix(d.WorkerProvider, "codex/") {
		return
	}
	if err := copyTree(filepath.Join(livePluginDir(t), "skills"), filepath.Join(repoRoot, ".agents", "skills")); err != nil {
		t.Fatalf("stage workflow-free skills into fixture .agents/skills: %v", err)
	}
}

func setupPaseoSplitRootFixture(t *testing.T, root string) paseoSplitRootFixture {
	t.Helper()
	repoRoot := filepath.Join(root, "repo")
	workflowDir := filepath.Join(repoRoot, "docs", "smoke")
	stateRoot := filepath.Join(workflowDir, ".spacedock-state")
	slug := "paseo-smoke"
	entityPath := filepath.Join(stateRoot, slug, "index.md")
	mustMkdirAll(t, workflowDir)
	testgit.InitRepo(t, repoRoot, "-q")

	mustWriteFile(t, filepath.Join(workflowDir, "README.md"), "---\ncommissioned-by: spacedock@1\nentity-type: task\nid-style: slug\nstate: .spacedock-state\nstages:\n  defaults:\n    worktree: false\n    concurrency: 1\n  states:\n    - name: implementation\n      initial: true\n      worktree: true\n    - name: done\n      terminal: true\n---\n# Paseo Smoke Workflow\n\n### implementation\n\nWrite receipt.txt on the worktree branch and append the stage report in .spacedock-state.\n\n- **Outputs:** Committed code file and path-scoped stage report commit.\n\n### done\n\nTerminal stage.\n")
	mustWriteFile(t, filepath.Join(repoRoot, ".gitignore"), "docs/smoke/.spacedock-state/\n.worktrees/\n.agents/\n")
	git(t, repoRoot, "add", ".gitignore", "docs/smoke/README.md")
	git(t, repoRoot, "commit", "-q", "-m", "init smoke workflow")

	git(t, repoRoot, "worktree", "add", "--orphan", "-b", "spacedock-state/smoke", stateRoot)
	mustMkdirAll(t, filepath.Dir(entityPath))
	mustWriteFile(t, entityPath, "---\nid: paseo-smoke\ntitle: Paseo Smoke Task\nstatus: implementation\nstarted:\ncompleted:\nverdict:\nworktree:\n---\n\n## Problem\n\nVerify end-to-end Paseo dispatch from FO to fresh ensign.\n")
	git(t, stateRoot, "add", filepath.Join(slug, "index.md"))
	git(t, stateRoot, "commit", "-q", "-m", "seed paseo-smoke entity")

	return paseoSplitRootFixture{RepoRoot: repoRoot, WorkflowDir: workflowDir, StateRoot: stateRoot, EntityPath: entityPath, Slug: slug}
}

func setupIsolatedPaseoCodexHome(t *testing.T, repo, artifactDir, realHome, oauthJSON, openAIKey string, auth codexLiveAuthDecision) string {
	t.Helper()
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Fatal("codex CLI not on PATH; install @openai/codex for codex/* Paseo live smoke")
	}
	codexHome := newCodexLiveIsolatedHome(t, repo, artifactDir)
	if err := seedCodexLiveConfig(codexHome); err != nil {
		t.Fatalf("seed codex config: %v", err)
	}
	if auth.mode == codexAuthOAuth {
		if err := seedCodexOAuthAuth(codexHome, oauthJSON); err != nil {
			t.Fatalf("seed codex OAuth auth: %v", err)
		}
	} else if auth.mode == codexAuthLocal {
		if err := seedCodexLocalAuth(codexHome, realHome); err != nil {
			t.Fatalf("seed local codex auth: %v", err)
		}
	}
	setupEnv := append(cleanEnviron("CODEX_HOME", "OPENAI_API_KEY", "CODEX_AUTH_JSON"), "CODEX_HOME="+codexHome)
	if auth.mode == codexAuthAPIKey {
		setupEnv = append(setupEnv, "OPENAI_API_KEY="+openAIKey)
		runCodexLiveCommand(t, artifactDir, "paseo-codex-login.txt", openAIKey+"\n", setupEnv, codexBin, "login", "--with-api-key")
	}
	install, err := cli.WriteCodexLocalMarketplace(filepath.Join(t.TempDir(), "codex-marketplace"), livePluginDir(t), "spacedock")
	if err != nil {
		t.Fatalf("stage workflow-free Codex local marketplace: %v", err)
	}
	runCodexLiveCommand(t, artifactDir, "paseo-codex-marketplace-add.txt", "", setupEnv, codexBin, "plugin", "marketplace", "add", install.MarketplaceRoot)
	runCodexLiveCommand(t, artifactDir, "paseo-codex-plugin-add.txt", "", setupEnv, codexBin, "plugin", "add", "spacedock@spacedock")
	return codexHome
}

func pickFreePaseoListenAddr(t *testing.T) string {
	t.Helper()
	for i := 0; i < 16; i++ {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen ephemeral port: %v", err)
		}
		port := ln.Addr().(*net.TCPAddr).Port
		_ = ln.Close()
		if port != 6767 && port != 6768 {
			return fmt.Sprintf("127.0.0.1:%d", port)
		}
	}
	t.Fatal("failed to allocate ephemeral Paseo port avoiding 6767/6768")
	return ""
}
