package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	piAuthAPIKey = "api-key"
	piAuthOAuth  = "oauth"
)

type piLiveAuthDecision struct{ mode, model, message string }

func decidePiLiveAuth(oauthJSON, openAIKey, required string) piLiveAuthDecision {
	if strings.TrimSpace(oauthJSON) != "" {
		return piLiveAuthDecision{mode: piAuthOAuth, model: "openai-codex/gpt-5.6-luna:max"}
	}
	if strings.TrimSpace(openAIKey) != "" {
		return piLiveAuthDecision{mode: piAuthAPIKey, model: "openai/gpt-5.6-luna:max"}
	}
	return piLiveAuthDecision{message: "Pi OAuth or OPENAI_API_KEY is required for the live lane"}
}

func seedPiOAuthAuth(piHome, oauthJSON string) error {
	var record map[string]any
	if err := json.Unmarshal([]byte(oauthJSON), &record); err != nil || record == nil {
		return fmt.Errorf("Pi OAuth auth is not a JSON object")
	}
	if err := os.MkdirAll(piHome, 0o700); err != nil {
		return fmt.Errorf("create isolated Pi home: %w", err)
	}
	payload, err := json.Marshal(map[string]any{"openai-codex": record})
	if err != nil {
		return fmt.Errorf("encode Pi OAuth auth: %w", err)
	}
	if err := os.WriteFile(filepath.Join(piHome, "auth.json"), payload, 0o600); err != nil {
		return fmt.Errorf("write Pi OAuth auth: %w", err)
	}
	return nil
}

func withoutPiEnvKey(env []string, key string) []string {
	prefix := key + "="
	out := env[:0]
	for _, value := range env {
		if !strings.HasPrefix(value, prefix) {
			out = append(out, value)
		}
	}
	return out
}
