// ABOUTME: reconcile-specific team-roster loader — keeps the ~/.claude/teams
// ABOUTME: read inside the Claude seam so generic internal/dispatch stays neutral.
package claudeteam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
// (useful for narrowed discovery when --team-name is omitted), and the members
// list. ConfigPath is the absolute path the loader resolved, primarily for
// diagnostic stderr messages.
type ReconcileTeamState struct {
	TeamName      string
	LeadSessionID string
	ConfigPath    string
	Members       []ReconcileMember
}

// LoadReconcileTeam loads the team config for the reconcile sweep. When
// teamName is non-empty, that path is loaded directly; missing/unreadable is
// an error. When teamName is empty, the loader scans ~/.claude/teams/*/config.json
// for the most-recently-modified config that contains at least one
// spacedock:ensign member — a stable proxy for "the live team in this session".
//
// Errors carry a stderr-ready message; the caller decides whether to surface as
// exit-1 (setup failure) or exit-2 (usage).
func LoadReconcileTeam(home, teamName string) (ReconcileTeamState, error) {
	if teamName != "" {
		cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
		if _, err := os.Stat(cfgPath); err != nil {
			return ReconcileTeamState{}, fmt.Errorf("team config not found: %s", cfgPath)
		}
		return loadReconcileConfigAt(cfgPath, teamName)
	}
	pattern := filepath.Join(home, ".claude", "teams", "*", "config.json")
	matches, _ := filepath.Glob(pattern)
	type candidate struct {
		path string
		name string
		mod  int64
	}
	var cands []candidate
	for _, p := range matches {
		state, err := loadReconcileConfigAt(p, filepath.Base(filepath.Dir(p)))
		if err != nil {
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
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		cands = append(cands, candidate{p, state.TeamName, info.ModTime().UnixNano()})
	}
	if len(cands) == 0 {
		return ReconcileTeamState{}, fmt.Errorf(
			"no team config with a spacedock:ensign member found under %s",
			filepath.Join(home, ".claude", "teams"))
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].mod > cands[j].mod })
	return loadReconcileConfigAt(cands[0].path, cands[0].name)
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
