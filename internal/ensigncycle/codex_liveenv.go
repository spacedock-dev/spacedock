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
