// ABOUTME: reconcile-specific team-roster loader — keeps the ~/.claude/teams
// ABOUTME: read inside the Claude seam so generic internal/dispatch stays neutral.
package claudeteam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ReconcileMember is one team member as the reconcile sweep needs it: name,
// agentType (the discriminator that says spacedock:ensign vs team-lead vs
// standing-teammate), and model (metadata). The roster shape lives behind the
// Claude seam so the generic dispatch helper never opens ~/.claude/teams itself.
type ReconcileMember struct {
	Name      string
	AgentType string
	Model     string
}

// ReconcileTeamState is what the reconcile loader returns: the resolved team
// name (which the helper echoes back into the JSON envelope), the leadSessionId
// (the session-identity key auto-discovery matches against the current session
// when --team-name is omitted), and the members list. ConfigPath is the absolute
// path the loader resolved, primarily for diagnostic stderr messages. An empty
// TeamName is the degrade-to-git-only sentinel: no session-scoped team resolved,
// so the caller emits only the git/filesystem drift classes and no roster claims.
type ReconcileTeamState struct {
	TeamName      string
	LeadSessionID string
	ConfigPath    string
	Members       []ReconcileMember
}

// LoadReconcileTeam loads the team config for the reconcile sweep. When
// teamName is non-empty, that path is loaded directly; missing/unreadable is
// an error. When teamName is empty, auto-discovery scans
// ~/.claude/teams/*/config.json for the config whose leadSessionId equals
// sessionID (the current lead session) among configs carrying a spacedock:ensign
// member. Exactly one session-id match resolves it. Zero matches, more than one
// match, or an empty sessionID return the degrade-to-git-only sentinel (a
// ReconcileTeamState with an empty TeamName and no error) rather than guessing —
// a roster derived from a non-session-matched config could be a stale prior
// session's or a parallel live session's team, so it is never trusted.
//
// Errors carry a stderr-ready message; the caller surfaces them as exit-1 (setup
// failure). The degrade sentinel is NOT an error: the sweep still runs and emits
// the session-independent git/filesystem classes.
func LoadReconcileTeam(home, teamName, sessionID string) (ReconcileTeamState, error) {
	if teamName != "" {
		cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
		if _, err := os.Stat(cfgPath); err != nil {
			return ReconcileTeamState{}, fmt.Errorf("team config not found: %s", cfgPath)
		}
		return loadReconcileConfigAt(cfgPath, teamName)
	}
	// Without a session id there is nothing to match against — degrade.
	if sessionID == "" {
		return ReconcileTeamState{}, nil
	}
	pattern := filepath.Join(home, ".claude", "teams", "*", "config.json")
	matches, _ := filepath.Glob(pattern)
	var hits []ReconcileTeamState
	for _, p := range matches {
		state, err := loadReconcileConfigAt(p, filepath.Base(filepath.Dir(p)))
		if err != nil {
			continue
		}
		if state.LeadSessionID != sessionID {
			continue
		}
		hasEnsign := false
		for _, m := range state.Members {
			if m.AgentType == "spacedock:ensign" {
				hasEnsign = true
				break
			}
		}
		if !hasEnsign {
			continue
		}
		hits = append(hits, state)
	}
	// Exactly one session-scoped team is the trusted roster. Zero or multiple
	// matches degrade — never guess which of several configs owns this session.
	if len(hits) != 1 {
		return ReconcileTeamState{}, nil
	}
	return hits[0], nil
}

// loadReconcileConfigAt decodes the config at path, returning the reconcile
// state. The teamName arg is what we report in the result — for direct loads
// it is the user-supplied name; for discovery it is the parent dir's basename.
func loadReconcileConfigAt(path, teamName string) (ReconcileTeamState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ReconcileTeamState{}, fmt.Errorf("team config unreadable: %s: %w", path, err)
	}
	var raw struct {
		LeadSessionID string `json:"leadSessionId"`
		Members       []struct {
			Name      string `json:"name"`
			AgentType string `json:"agentType"`
			Model     string `json:"model"`
		} `json:"members"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return ReconcileTeamState{}, fmt.Errorf("team config malformed: %s: %w", path, err)
	}
	state := ReconcileTeamState{
		TeamName:      teamName,
		LeadSessionID: raw.LeadSessionID,
		ConfigPath:    path,
	}
	for _, m := range raw.Members {
		state.Members = append(state.Members, ReconcileMember{
			Name: m.Name, AgentType: m.AgentType, Model: m.Model,
		})
	}
	return state, nil
}
