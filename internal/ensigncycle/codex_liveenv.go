package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type codexAuthMode int

const (
	codexAuthSkip codexAuthMode = iota
	codexAuthAPIKey
	codexAuthLocal
	codexAuthFatal
)

type codexLiveAuthDecision struct {
	mode    codexAuthMode
	message string
}

const codexLiveMultiAgentConfig = `[features.multi_agent_v2]
max_concurrent_threads_per_session = 16
tool_namespace = "agents"
hide_spawn_agent_metadata = false
`

func decideCodexLiveAuth(openAIAPIKey string, localAuthAvailable bool, required string) codexLiveAuthDecision {
	if openAIAPIKey != "" {
		return codexLiveAuthDecision{mode: codexAuthAPIKey}
	}
	if required != "" {
		return codexLiveAuthDecision{
			mode:    codexAuthFatal,
			message: "OPENAI_API_KEY is required for the approval-gated codex-live lane",
		}
	}
	if localAuthAvailable {
		return codexLiveAuthDecision{mode: codexAuthLocal}
	}
	return codexLiveAuthDecision{
		mode:    codexAuthSkip,
		message: "no live Codex auth available: set OPENAI_API_KEY or log in with codex to run the live Codex shared suite",
	}
}

func codexLocalAuthAvailable(realHome string) bool {
	authPath := codexAuthPath(realHome)
	if authPath == "" {
		return false
	}
	b, err := os.ReadFile(authPath)
	return err == nil && strings.TrimSpace(string(b)) != ""
}

func seedCodexLiveConfig(codexHome string) error {
	if codexHome == "" {
		return fmt.Errorf("isolated CODEX_HOME is empty")
	}
	if err := os.MkdirAll(codexHome, 0o700); err != nil {
		return fmt.Errorf("create isolated CODEX_HOME: %w", err)
	}
	if err := os.WriteFile(filepath.Join(codexHome, "config.toml"), []byte(codexLiveMultiAgentConfig), 0o600); err != nil {
		return fmt.Errorf("write isolated Codex config: %w", err)
	}
	return nil
}

func seedCodexLocalAuth(codexHome, realHome string) error {
	authPath := codexAuthPath(realHome)
	if authPath == "" {
		return fmt.Errorf("real HOME is empty")
	}
	b, err := os.ReadFile(authPath)
	if err != nil {
		return fmt.Errorf("read Codex auth: %w", err)
	}
	if strings.TrimSpace(string(b)) == "" {
		return fmt.Errorf("Codex auth file is empty")
	}
	if err := os.MkdirAll(codexHome, 0o700); err != nil {
		return fmt.Errorf("create isolated CODEX_HOME: %w", err)
	}
	if err := os.WriteFile(filepath.Join(codexHome, "auth.json"), b, 0o600); err != nil {
		return fmt.Errorf("write isolated Codex auth: %w", err)
	}
	return nil
}

func codexAuthPath(realHome string) string {
	if realHome == "" {
		return ""
	}
	return filepath.Join(realHome, ".codex", "auth.json")
}

func codexLiveIsolatedHomeParent(cacheDir string) (string, error) {
	if cacheDir == "" {
		return "", fmt.Errorf("user cache dir is empty")
	}
	return filepath.Join(cacheDir, "spacedock-live-codex"), nil
}

func codexLiveIsolatedHomeParentCandidates(cacheDir, repoRoot, artifactRoot string) []string {
	var out []string
	if artifactRoot != "" && !pathIsUnderSystemTemp(artifactRoot) && !pathIsUnder(artifactRoot, repoRoot) {
		out = append(out, filepath.Join(artifactRoot, "_codex-home"))
	}
	if parent, err := codexLiveIsolatedHomeParent(cacheDir); err == nil {
		out = append(out, parent)
	}
	if parent, err := codexLiveRepoAdjacentHomeParent(repoRoot); err == nil {
		out = append(out, parent)
	}
	return out
}

func codexLiveRepoAdjacentHomeParent(repoRoot string) (string, error) {
	if repoRoot == "" {
		return "", fmt.Errorf("repo root is empty")
	}
	repoRoot = filepath.Clean(repoRoot)
	return filepath.Join(filepath.Dir(repoRoot), ".spacedock-live-codex", filepath.Base(repoRoot)), nil
}

func pathIsUnder(path, dir string) bool {
	if path == "" || dir == "" {
		return false
	}
	dir = filepath.Clean(dir)
	path = filepath.Clean(path)
	if path == dir {
		return true
	}
	rel, err := filepath.Rel(dir, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func pathIsUnderSystemTemp(path string) bool {
	if path == "" {
		return false
	}
	tmp := filepath.Clean(os.TempDir())
	p := filepath.Clean(path)
	if rel, err := filepath.Rel(tmp, p); err == nil {
		return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
	}
	return false
}
