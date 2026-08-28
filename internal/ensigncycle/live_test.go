//go:build live

// ABOUTME: Live BEHAVIORAL test of the dispatch->ensign->stage cycle driven by a
// ABOUTME: REAL model through the spacedock claude front door (gated, -tags live only).
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// spacedockBinary resolves the built v1 binary the test shells. SPACEDOCK_BIN
// (set by the CI job after `go build -o ./spacedock`) takes precedence; locally
// it falls back to a `spacedock` on PATH. The test fails loudly when neither
// resolves rather than silently shelling a stale or absent binary.
func spacedockBinary(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("SPACEDOCK_BIN"); p != "" {
		abs, err := filepath.Abs(p)
		if err != nil {
			t.Fatalf("SPACEDOCK_BIN=%q is not resolvable: %v", p, err)
		}
		if _, err := os.Stat(abs); err != nil {
			t.Fatalf("SPACEDOCK_BIN=%q does not exist: %v", abs, err)
		}
		return abs
	}
	p, err := exec.LookPath("spacedock")
	if err != nil {
		t.Fatal("no spacedock binary: set SPACEDOCK_BIN to the built binary or put spacedock on PATH")
	}
	return p
}

// repoRoot resolves the plugin-checkout root passed to --plugin-dir. The
// ensigncycle package lives at internal/ensigncycle, so the repo root is two
// directories up from the test's working directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("SPACEDOCK_REPO_ROOT"); p != "" {
		abs, err := filepath.Abs(p)
		if err != nil {
			t.Fatalf("SPACEDOCK_REPO_ROOT=%q is not resolvable: %v", p, err)
		}
		return abs
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join(wd, "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// livePluginDir stages an ISOLATED plugin checkout for `--plugin-dir` and returns
// its path. It exists to stop a wrong-root boot: the real repo root carries a
// discoverable `docs/dev` workflow (with live entities), so an FO that anchors its
// `git rev-parse --show-toplevel` + `status --discover` on the plugin path — instead
// of its isolated cwd fixture — finds and drives the REAL workflow. Staging copies
// ONLY the plugin scaffolding (`.claude-plugin/`, `.codex-plugin/`, `skills/`,
// `agents/`) into a temp
// dir with NO `docs/dev` sibling, then `git init`s it so a `rev-parse` from the
// plugin path resolves to a workflow-free root. An FO that boots from here discovers
// zero workflows and falls back to the cwd fixture. The result is cached per repo
// root so parallel scenarios share one staged copy.
func livePluginDir(t *testing.T) string {
	t.Helper()
	return cachedLivePluginDir(t, repoRoot(t))
}

var (
	livePluginOnce sync.Once
	livePluginPath string
	livePluginErr  error
)

func cachedLivePluginDir(t *testing.T, repo string) string {
	t.Helper()
	livePluginOnce.Do(func() {
		// MkdirTemp (not t.TempDir) so the staged plugin outlives the first test's
		// cleanup and the cached path stays valid for every scenario in the run.
		marketplace, err := os.MkdirTemp("", "spacedock-live-plugin-")
		if err != nil {
			livePluginErr = err
			return
		}
		staged := filepath.Join(marketplace, "spacedock")
		for _, sub := range []string{".claude-plugin", ".codex-plugin", "skills", "agents"} {
			src := filepath.Join(repo, sub)
			if _, statErr := os.Stat(src); statErr != nil {
				continue // optional members (e.g. a layout without a top-level agents/)
			}
			if copyErr := copyTree(src, filepath.Join(staged, sub)); copyErr != nil {
				livePluginErr = copyErr
				return
			}
		}
		manifestDir := filepath.Join(marketplace, ".claude-plugin")
		if err := os.MkdirAll(manifestDir, 0o755); err != nil {
			livePluginErr = err
			return
		}
		manifest := []byte("{\n  \"name\": \"spacedock\",\n  \"owner\": { \"name\": \"Spacedock live suite\" },\n  \"plugins\": [\n    { \"name\": \"spacedock\", \"source\": \"./spacedock\", \"description\": \"release candidate\", \"category\": \"workflow\" }\n  ]\n}\n")
		if err := os.WriteFile(filepath.Join(manifestDir, "marketplace.json"), manifest, 0o644); err != nil {
			livePluginErr = err
			return
		}
		// git init so the FO's `git rev-parse --show-toplevel` resolves to this
		// workflow-free root, not an enclosing checkout that has a docs/dev.
		testgit.InitRepo(t, staged, "-q")
		livePluginPath = staged
	})
	if livePluginErr != nil {
		t.Fatalf("stage isolated live plugin dir: %v", livePluginErr)
	}
	return livePluginPath
}

var (
	stableLiveBinaryOnce sync.Once
	stableLiveBinaryPath string
	stableLiveBinaryErr  error
)

// stableLiveRelease returns the current plugin packaged as the stable channel
// plus a binary stamped with that package's release version and channel. Every
// common live journey installs this package before using the ordinary front door.
func stableLiveRelease(t *testing.T) (binary, marketplace string) {
	t.Helper()
	plugin := livePluginDir(t)
	stableLiveBinaryOnce.Do(func() {
		manifestData, err := os.ReadFile(filepath.Join(plugin, ".claude-plugin", "plugin.json"))
		if err != nil {
			stableLiveBinaryErr = err
			return
		}
		var manifest struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(manifestData, &manifest); err != nil {
			stableLiveBinaryErr = err
			return
		}
		buildDir, err := os.MkdirTemp("", "spacedock-live-stable-")
		if err != nil {
			stableLiveBinaryErr = err
			return
		}
		stableLiveBinaryPath = filepath.Join(buildDir, "spacedock")
		stamp := fmt.Sprintf("-X github.com/spacedock-dev/spacedock/internal/cli.Version=%s -X github.com/spacedock-dev/spacedock/internal/cli.devBranch=main", manifest.Version)
		cmd := exec.Command("go", "build", "-ldflags", stamp, "-o", stableLiveBinaryPath, "./cmd/spacedock")
		cmd.Dir = repoRoot(t)
		if out, err := cmd.CombinedOutput(); err != nil {
			stableLiveBinaryErr = fmt.Errorf("build stable live binary: %w: %s", err, out)
		}
	})
	if stableLiveBinaryErr != nil {
		t.Fatal(stableLiveBinaryErr)
	}
	return stableLiveBinaryPath, filepath.Dir(plugin)
}

// copyTree recursively copies src to dst, preserving file modes. Symlinks are
// resolved to real files so the staged plugin has no path back into the real repo.
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		target := filepath.Join(dst, rel)
		info, infoErr := os.Stat(path) // Stat (not Lstat) resolves symlinks to real content
		if infoErr != nil {
			return infoErr
		}
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if mkErr := os.MkdirAll(filepath.Dir(target), 0o755); mkErr != nil {
			return mkErr
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}

// readmeRealisticLifecycle is the live cycle's workflow README: a realistic
// ≥3-stage lifecycle — backlog (initial) → implementation (work) → done
// (terminal) — matching every real Spacedock workflow. The distinct WORK stage
// between initial and terminal makes the FO's TERMINALIZE step
// (implementation→done with a verdict) DISTINCT from its DISPATCH step
// (backlog→implementation): the FO records a verdict naturally at terminalization
// and the M1 verdict gate has a real trigger. A 2-stage backlog→done fixture
// collapses dispatch and terminalize onto the same stage, so the FO never runs a
// distinct finalize and no verdict is recorded (the failure this fixture fixes).
// It is live-only (kept beside the //go:build live test) so the offline minimal
// fixture readmeNonWorktree — which the single-stage mechanical test pins — is
// untouched. All stages are non-worktree to keep the flat-entity path.
func readmeRealisticLifecycle() string {
	return "---\n" +
		"commissioned-by: spacedock@1\n" +
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
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Fixture Workflow\n" +
		"\n" +
		"### backlog\n\nseed.\n\n- **Outputs:** a one-line note.\n\n" +
		"### implementation\n\nDo the trivial work and write the note.\n\n- **Outputs:** the note recorded.\n\n" +
		"### done\n\nterm.\n"
}

// isTeamCreate matches the FO's TeamCreate assistant tool_use — the first
// progress beat of the teams-mode live cycle (the team is engaged).
func isTeamCreate(e streamEntry) bool {
	b := e.toolUseBlock()
	return b != nil && b.Name == "TeamCreate"
}

// isEnsignDispatch matches the FO's first ensign dispatch — an
// Agent(subagent_type="spacedock:ensign") assistant tool_use. The contract runs
// `spacedock dispatch spawn-standing-all` immediately before this dispatch, so its
// OPEN is the reliable barrier for "standing teammates have been injected" (used by
// the residency test, which only needs injection to have run, not the dispatch to
// close). It scans ALL tool_use blocks so an Agent dispatch riding as a second
// block in a multi-tool turn is not missed.
func isEnsignDispatch(e streamEntry) bool {
	for _, b := range e.toolUseBlocks() {
		// Claude can omit a defaulted subagent_type from a successful Agent
		// input. Treat that omission as a dispatch barrier only; the required
		// merged oracle separately proves named/background/no-team transport,
		// the ensign prompt, and on-disk agentType identity.
		if b.Name == "Agent" && (b.Input.SubagentType == "spacedock:ensign" || b.Input.SubagentType == "") {
			return true
		}
	}
	return false
}
