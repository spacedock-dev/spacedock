// ABOUTME: session-scoping unit table for LoadReconcileTeam — proves discovery
// ABOUTME: follows config.leadSessionId, not newest-mtime, and degrades safely.
package claudeteam

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeSessionTeamConfig writes a config.json under home/.claude/teams/{team}/
// carrying the given leadSessionId and exactly one spacedock:ensign member, then
// stamps its mtime to modAge ago so a test can make a decoy NEWER than the match.
func writeSessionTeamConfig(t *testing.T, home, team, leadSessionID string, modAge time.Duration) string {
	t.Helper()
	cfgPath := filepath.Join(home, ".claude", "teams", team, "config.json")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := map[string]any{
		"name":          team,
		"leadSessionId": leadSessionID,
		"members": []map[string]string{
			{"name": "spacedock-ensign-x-implementation", "agentType": "spacedock:ensign", "model": "m"},
		},
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		t.Fatal(err)
	}
	mod := time.Now().Add(-modAge)
	if err := os.Chtimes(cfgPath, mod, mod); err != nil {
		t.Fatal(err)
	}
	return cfgPath
}

// TestLoadReconcileTeamMatchesSessionNotMtime is the load-bearing discovery
// test: a stale decoy config is written with a NEWER mtime than the
// session-matched config, so the retired newest-mtime loader would pick the
// decoy. The fixed loader must follow leadSessionId and resolve the match.
func TestLoadReconcileTeamMatchesSessionNotMtime(t *testing.T) {
	home := t.TempDir()
	const sessionID = "11111111-1111-1111-1111-111111111111"
	// Match: older mtime, leadSessionId == current session.
	writeSessionTeamConfig(t, home, "team-match", sessionID, 10*time.Minute)
	// Decoy: NEWER mtime, different leadSessionId. mtime would pick this one.
	writeSessionTeamConfig(t, home, "team-decoy", "99999999-9999-9999-9999-999999999999", 1*time.Minute)

	state, err := LoadReconcileTeam(home, "", sessionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.TeamName != "team-match" {
		t.Errorf("resolved team=%q, want team-match (session id, not newest mtime)", state.TeamName)
	}
	if state.LeadSessionID != sessionID {
		t.Errorf("resolved leadSessionId=%q, want %q", state.LeadSessionID, sessionID)
	}
}

// TestLoadReconcileTeamNoSessionMatchDegrades proves the safe floor: when no
// config's leadSessionId equals the injected session id, the loader returns a
// sentinel (TeamName == "", no error) so the assembly can degrade to git-only.
// This is the repeated/parallel-session case — the only on-disk team belongs to
// a foreign session and must NEVER be trusted as our roster.
func TestLoadReconcileTeamNoSessionMatchDegrades(t *testing.T) {
	home := t.TempDir()
	writeSessionTeamConfig(t, home, "team-foreign", "ffffffff-ffff-ffff-ffff-ffffffffffff", 1*time.Minute)

	state, err := LoadReconcileTeam(home, "", "current-session-no-match")
	if err != nil {
		t.Fatalf("degrade path must not error; got %v", err)
	}
	if state.TeamName != "" {
		t.Errorf("degrade sentinel must have empty TeamName; got %q", state.TeamName)
	}
	if len(state.Members) != 0 {
		t.Errorf("degrade sentinel must carry no roster; got %d members", len(state.Members))
	}
}

// TestLoadReconcileTeamEmptySessionDegrades proves the env-absent case (real
// shells where CLAUDE_CODE_SESSION_ID is unset): an empty session id can never
// match a config, so the loader degrades rather than guessing by mtime.
func TestLoadReconcileTeamEmptySessionDegrades(t *testing.T) {
	home := t.TempDir()
	writeSessionTeamConfig(t, home, "team-a", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", 1*time.Minute)

	state, err := LoadReconcileTeam(home, "", "")
	if err != nil {
		t.Fatalf("empty-session degrade must not error; got %v", err)
	}
	if state.TeamName != "" {
		t.Errorf("empty session must degrade (empty TeamName); got %q", state.TeamName)
	}
}

// TestLoadReconcileTeamMultipleSessionMatchesDegrades keeps the "never trust an
// unverifiable roster" invariant absolute: two configs sharing the injected
// session id (not expected in practice — one team per lead session) degrade to
// git-only rather than guessing which is authoritative.
func TestLoadReconcileTeamMultipleSessionMatchesDegrades(t *testing.T) {
	home := t.TempDir()
	const sessionID = "22222222-2222-2222-2222-222222222222"
	writeSessionTeamConfig(t, home, "team-dup-1", sessionID, 5*time.Minute)
	writeSessionTeamConfig(t, home, "team-dup-2", sessionID, 1*time.Minute)

	state, err := LoadReconcileTeam(home, "", sessionID)
	if err != nil {
		t.Fatalf("multi-match degrade must not error; got %v", err)
	}
	if state.TeamName != "" {
		t.Errorf("multi-match must degrade (empty TeamName); got %q", state.TeamName)
	}
}

// TestLoadReconcileTeamExplicitNameIgnoresSession confirms the explicit path is
// unchanged: --team-name loads that config directly regardless of session id.
func TestLoadReconcileTeamExplicitNameIgnoresSession(t *testing.T) {
	home := t.TempDir()
	writeSessionTeamConfig(t, home, "team-explicit", "some-other-session", 1*time.Minute)

	state, err := LoadReconcileTeam(home, "team-explicit", "unrelated-session")
	if err != nil {
		t.Fatalf("explicit path error: %v", err)
	}
	if state.TeamName != "team-explicit" {
		t.Errorf("explicit team=%q, want team-explicit", state.TeamName)
	}
	if len(state.Members) == 0 {
		t.Errorf("explicit path must load roster; got 0 members")
	}
}

// TestLoadReconcileTeamExplicitMissingErrors confirms a missing explicit config
// is still a setup failure (exit-1 territory) — the degrade must not mask it.
func TestLoadReconcileTeamExplicitMissingErrors(t *testing.T) {
	home := t.TempDir()
	if _, err := LoadReconcileTeam(home, "team-absent", "any"); err == nil {
		t.Errorf("missing explicit config must error, not degrade")
	}
}
