// ABOUTME: Real-Claude AC-3 provider identity proof for post-fence additional plugins.
// ABOUTME: Reads Claude's agent-listing event and includes an impersonating-provider mutation.
package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	installedProviderMarker = "INSTALLED_PROVIDER_IDENTITY"
	mutantProviderMarker    = "MUTANT_ADDITIONAL_PROVIDER_IDENTITY"
)

// claudeProviderHost keeps the production resolver but captures the real host
// launch. A Claude authentication failure remains a host exit code, not a launch
// error, so the pre-API agent-listing event can be inspected deterministically.
type claudeProviderHost struct {
	execHost
	resolveCalls     int
	resolvedManifest string
	launchedArg      []string
	launchOutput     string
}

func (h *claudeProviderHost) ResolveManifest(host string) (string, error) {
	h.resolveCalls++
	manifest, err := h.execHost.ResolveManifest(host)
	h.resolvedManifest = manifest
	return manifest, err
}

func (h *claudeProviderHost) Launch(argv []string, env []string) (int, error) {
	h.launchedArg = append([]string(nil), argv...)
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	h.launchOutput = string(out)
	if err == nil {
		return 0, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), nil
	}
	return 1, err
}

// TestClaudeAdditionalPluginKeepsInstalledSpacedockProvider drives the real
// Claude 2.x loader in an isolated config. Claude records its resolved agents in
// an agent_listing_delta before the deliberately invalid API key is rejected,
// giving the test a host-owned provider-identity observation without an LLM
// response oracle.
//
// The second observation is the adversarial mutation: changing only the
// additional plugin's manifest name to `spacedock` makes it impersonate the
// provider. The installed compatibility gate still runs and the post-fence argv
// is byte-identical, but Claude's registered spacedock:first-officer description
// flips to the mutant marker. If the positive assertion were accidentally wired
// only to gate/argv again, this sensitivity check would fail.
func TestClaudeAdditionalPluginKeepsInstalledSpacedockProvider(t *testing.T) {
	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		t.Skip("claude not on PATH; provider-identity behavior test requires the host CLI")
	}
	savedBranch := devBranch
	devBranch = "main"
	defer func() { devBranch = savedBranch }()

	observed, positive := observeClaudeFirstOfficerProvider(t, claudeBin, "additional-tools")
	if !strings.Contains(observed, installedProviderMarker) {
		t.Fatalf("resolved provider = %q, want installed marker %q; launch=%v output=%s", observed, installedProviderMarker, positive.launchedArg, positive.launchOutput)
	}
	if positive.resolveCalls == 0 || positive.resolvedManifest == "" {
		t.Fatalf("installed compatibility gate was not observed: resolves=%d manifest=%q", positive.resolveCalls, positive.resolvedManifest)
	}
	if !argvContainsPair(positive.launchedArg, "--plugin-dir", positive.additionalDir) {
		t.Fatalf("additional plugin did not reach Claude unchanged: %v", positive.launchedArg)
	}

	mutated, mutant := observeClaudeFirstOfficerProvider(t, claudeBin, "spacedock")
	if !strings.Contains(mutated, mutantProviderMarker) {
		t.Fatalf("mutation did not make the additional directory win: provider=%q launch=%v output=%s", mutated, mutant.launchedArg, mutant.launchOutput)
	}
	if mutant.resolveCalls == 0 || mutant.resolvedManifest == "" || !argvContainsPair(mutant.launchedArg, "--plugin-dir", mutant.additionalDir) {
		t.Fatalf("mutation changed the gate/argv preconditions instead of provider identity: resolves=%d argv=%v", mutant.resolveCalls, mutant.launchedArg)
	}
	if !equalArgv(normalizeClaudeProviderArgv(positive.launchedArg), normalizeClaudeProviderArgv(mutant.launchedArg)) {
		t.Fatalf("mutation changed launch shape instead of only provider identity:\npositive=%v\nmutant=%v", positive.launchedArg, mutant.launchedArg)
	}
}

type claudeProviderObservation struct {
	*claudeProviderHost
	additionalDir string
}

func observeClaudeFirstOfficerProvider(t *testing.T, claudeBin, additionalName string) (string, claudeProviderObservation) {
	t.Helper()
	tmp := t.TempDir()
	configDir := filepath.Join(tmp, "config")
	cacheDir := filepath.Join(tmp, "cache")
	homeDir := filepath.Join(tmp, "home")
	mustMkdir(t, configDir)
	mustMkdir(t, cacheDir)
	mustMkdir(t, homeDir)

	marketplace := buildClaudeProviderMarketplace(t, filepath.Join(tmp, "marketplace-root"))
	additional := buildClaudeProviderPlugin(t, filepath.Join(tmp, "additional"), additionalName, mutantProviderMarker)
	t.Setenv("HOME", homeDir)
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)
	t.Setenv("CLAUDE_CODE_PLUGIN_CACHE_DIR", cacheDir)
	t.Setenv("ANTHROPIC_API_KEY", "invalid-provider-identity-test-key")
	t.Setenv("CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", "1")
	withExecutablePath(t, filepath.Join(tmp, "not-adjacent"), os.ErrNotExist)

	env := os.Environ()
	runHost(t, claudeBin, env, "plugin", "marketplace", "add", marketplace)
	runHost(t, claudeBin, env, "plugin", "install", "spacedock@spacedock")

	debugLog := filepath.Join(tmp, "claude-debug.log")
	host := &claudeProviderHost{}
	var stdout, stderr bytes.Buffer
	code := runClaude(context.Background(), []string{
		"--", "--plugin-dir", additional,
		"--print", "--debug-file", debugLog, "--max-budget-usd", "0.01",
	}, t.TempDir(), host, lookFound, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("Claude unexpectedly accepted the deliberately invalid API key; output=%s", host.launchOutput)
	}

	provider := readClaudeAgentListing(t, configDir, "spacedock:first-officer")
	return provider, claudeProviderObservation{claudeProviderHost: host, additionalDir: additional}
}

func buildClaudeProviderMarketplace(t *testing.T, root string) string {
	t.Helper()
	marketplace := filepath.Join(root, "marketplace")
	plugin := filepath.Join(marketplace, "spacedock")
	mustMkdir(t, filepath.Join(marketplace, ".claude-plugin"))
	buildClaudeProviderPlugin(t, plugin, "spacedock", installedProviderMarker)
	mustWrite(t, filepath.Join(marketplace, ".claude-plugin", "marketplace.json"), `{
  "name": "spacedock",
  "owner": { "name": "test" },
  "plugins": [
    { "name": "spacedock", "source": "./spacedock", "description": "test", "category": "workflow" }
  ]
}
`)
	return marketplace
}

func buildClaudeProviderPlugin(t *testing.T, root, name, marker string) string {
	t.Helper()
	mustMkdir(t, filepath.Join(root, ".claude-plugin"))
	mustMkdir(t, filepath.Join(root, "agents"))
	mustMkdir(t, filepath.Join(root, "skills"))
	mustWrite(t, filepath.Join(root, ".claude-plugin", "plugin.json"),
		`{"name":"`+name+`","version":"`+displayVersion()+`","skills":"./skills/"}`+"\n")
	mustWrite(t, filepath.Join(root, "agents", "first-officer.md"),
		"---\nname: first-officer\ndescription: "+marker+"\n---\n\n"+marker+"\n")
	return root
}

func readClaudeAgentListing(t *testing.T, configDir, agentName string) string {
	t.Helper()
	var found string
	err := filepath.WalkDir(configDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".jsonl" {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var event struct {
				Attachment struct {
					Type       string   `json:"type"`
					AddedLines []string `json:"addedLines"`
				} `json:"attachment"`
			}
			if json.Unmarshal(scanner.Bytes(), &event) != nil || event.Attachment.Type != "agent_listing_delta" {
				continue
			}
			prefix := "- " + agentName + ":"
			for _, line := range event.Attachment.AddedLines {
				if strings.HasPrefix(line, prefix) {
					found = line
				}
			}
		}
		return scanner.Err()
	})
	if err != nil {
		t.Fatalf("read Claude agent listing: %v", err)
	}
	if found == "" {
		t.Fatalf("Claude session did not record provider %s under %s", agentName, configDir)
	}
	return found
}

func argvContainsPair(argv []string, flag, value string) bool {
	for i := 0; i+1 < len(argv); i++ {
		if argv[i] == flag && argv[i+1] == value {
			return true
		}
	}
	return false
}

func normalizeClaudeProviderArgv(argv []string) []string {
	out := append([]string(nil), argv...)
	for i := 0; i+1 < len(out); i++ {
		switch out[i] {
		case "--plugin-dir":
			out[i+1] = "<additional-plugin>"
		case "--debug-file":
			out[i+1] = "<debug-log>"
		}
	}
	return out
}
