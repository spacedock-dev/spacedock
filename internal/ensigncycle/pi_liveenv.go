package ensigncycle

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	piAuthAPIKey = "api-key"
	piAuthOAuth  = "oauth"

	piOAuthModel  = "openai-codex/gpt-5.6-luna:max"
	piAPIKeyModel = "openai/gpt-5.6-luna:max"
)

type piLiveAuthDecision struct{ mode, model, message string }

func decidePiLiveAuth(oauthJSON, openAIKey, required string) piLiveAuthDecision {
	if strings.TrimSpace(oauthJSON) != "" {
		return piLiveAuthDecision{mode: piAuthOAuth, model: piOAuthModel}
	}
	if strings.TrimSpace(openAIKey) != "" {
		return piLiveAuthDecision{mode: piAuthAPIKey, model: piAPIKeyModel}
	}
	return piLiveAuthDecision{message: "Pi OAuth or OPENAI_API_KEY is required for the live lane"}
}

type codexAuthFile struct {
	AuthMode string `json:"auth_mode"`
	Tokens   struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		AccountID    string `json:"account_id"`
	} `json:"tokens"`
}

type piOAuthRecord struct {
	Type      string `json:"type"`
	Access    string `json:"access"`
	Refresh   string `json:"refresh"`
	Expires   int64  `json:"expires"`
	AccountID string `json:"accountId"`
}

// codexOAuthRecord converts the complete Codex auth file into Pi's provider
// record. The source file is intentionally never returned or passed through to
// the child process; only the fields Pi consumes are retained.
func codexOAuthRecord(authJSON string) (piOAuthRecord, error) {
	var auth codexAuthFile
	if err := json.Unmarshal([]byte(authJSON), &auth); err != nil || auth.Tokens.AccessToken == "" || auth.Tokens.RefreshToken == "" || auth.Tokens.AccountID == "" {
		return piOAuthRecord{}, fmt.Errorf("Codex OAuth auth has an invalid credential record")
	}
	expires, err := jwtExpiryMillis(auth.Tokens.AccessToken)
	if err != nil {
		return piOAuthRecord{}, fmt.Errorf("Codex OAuth auth has an invalid access token")
	}
	return piOAuthRecord{
		Type:      "oauth",
		Access:    auth.Tokens.AccessToken,
		Refresh:   auth.Tokens.RefreshToken,
		Expires:   expires,
		AccountID: auth.Tokens.AccountID,
	}, nil
}

func jwtExpiryMillis(token string) (int64, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("not a JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, fmt.Errorf("decode JWT payload: %w", err)
	}
	var claims struct {
		ExpiresAt json.Number `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.ExpiresAt == "" {
		return 0, fmt.Errorf("JWT expiry is missing")
	}
	seconds, err := claims.ExpiresAt.Int64()
	if err != nil || seconds <= 0 || seconds > math.MaxInt64/1000 {
		return 0, fmt.Errorf("JWT expiry is invalid")
	}
	return seconds * 1000, nil
}

func seedPiOAuthAuth(piHome, codexJSON string) error {
	record, err := codexOAuthRecord(codexJSON)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(piHome, 0o700); err != nil {
		return fmt.Errorf("create isolated Pi home: %w", err)
	}
	payload, err := json.Marshal(map[string]piOAuthRecord{"openai-codex": record})
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
