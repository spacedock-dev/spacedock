// ABOUTME: Merged-lane assertion predicates — the headless `-p` merged host shape (no TeamCreate,
// ABOUTME: named background Agent, in-process subagent meta) parsed over the stream-json + on-disk meta.
package ensigncycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mergedAgentDispatch is the parsed shape of the merged-mode dispatch the FO emits
// on a host where TeamCreate is gone (claude ≥2.1.178): an Agent tool_use carrying
// a member name, the spacedock:ensign subagent type, run_in_background true, and NO
// team_name. The merged dispatch's two inter-agent communication halves are `name` (lead→worker)
// and `run_in_background` (worker→lead); the absence of `team_name` is what
// distinguishes it from a legacy team dispatch.
type mergedAgentDispatch struct {
	name            string
	runInBackground bool
	hasTeamName     bool
}

// mergedEnsignDispatches scans a stream-json transcript for explicit ensign Agent
// calls and for the merged transport shape when Claude omits its defaulted
// subagent_type. Identity is proven independently from the dispatch artifact and
// on-disk member meta. team_name PRESENCE is what hasTeamName reports.
func mergedEnsignDispatches(lines []string) []mergedAgentDispatch {
	type rawInput struct {
		SubagentType    json.RawMessage `json:"subagent_type"`
		Name            string          `json:"name"`
		RunInBackground bool            `json:"run_in_background"`
		TeamName        json.RawMessage `json:"team_name"`
	}
	type rawBlock struct {
		Type  string   `json:"type"`
		Name  string   `json:"name"`
		Input rawInput `json:"input"`
	}
	type rawEntry struct {
		Type    string `json:"type"`
		Message *struct {
			Content []rawBlock `json:"content"`
		} `json:"message"`
	}
	var out []mergedAgentDispatch
	for _, line := range lines {
		var e rawEntry
		if json.Unmarshal([]byte(line), &e) != nil || e.Type != "assistant" || e.Message == nil {
			continue
		}
		for _, b := range e.Message.Content {
			if b.Type != "tool_use" || b.Name != "Agent" {
				continue
			}
			mergedTransport := b.Input.Name != "" && b.Input.RunInBackground && (len(b.Input.TeamName) == 0 || string(b.Input.TeamName) == "null")
			subagentType := ""
			subagentTypeOmitted := len(b.Input.SubagentType) == 0
			if !subagentTypeOmitted {
				_ = json.Unmarshal(b.Input.SubagentType, &subagentType)
			}
			if subagentType != "spacedock:ensign" && !(subagentTypeOmitted && mergedTransport) {
				continue
			}
			out = append(out, mergedAgentDispatch{
				name:            b.Input.Name,
				runInBackground: b.Input.RunInBackground,
				// team_name present iff the key decoded to a non-null JSON value.
				hasTeamName: len(b.Input.TeamName) > 0 && string(b.Input.TeamName) != "null",
			})
		}
	}
	return out
}

// streamHasTeamCreateToolUse reports whether any assistant entry in the stream
// emitted a TeamCreate (or TeamDelete) tool_use. On a merged host this must be
// false: the native team tools are gone, so the FO must reach the named-background
// dispatch without ever attempting TeamCreate (and must not fall to bare/sequential
// — proven separately by the presence of a merged Agent dispatch).
func streamHasTeamCreateToolUse(lines []string) bool {
	type rawBlock struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}
	type rawEntry struct {
		Type    string `json:"type"`
		Message *struct {
			Content []rawBlock `json:"content"`
		} `json:"message"`
	}
	for _, line := range lines {
		var e rawEntry
		if json.Unmarshal([]byte(line), &e) != nil || e.Type != "assistant" || e.Message == nil {
			continue
		}
		for _, b := range e.Message.Content {
			if b.Type == "tool_use" && (b.Name == "TeamCreate" || b.Name == "TeamDelete") {
				return true
			}
		}
	}
	return false
}

// initEventToolNames returns the `tools` list from the stream's `system`/`init`
// event (the first stream-json line of a `claude -p` run), or nil when absent. On a
// merged host the init tool surface carries SendMessage but NOT TeamCreate/TeamDelete
// — the static, parse-once analog of the contract's `ToolSearch(select:TeamCreate)`-
// empty discriminator: a Go test cannot run the ToolSearch hop, but the init event
// it reads tells the same truth (the native team tools are not on the merged host's
// tool surface).
func initEventToolNames(lines []string) []string {
	type initEntry struct {
		Type    string   `json:"type"`
		Subtype string   `json:"subtype"`
		Tools   []string `json:"tools"`
	}
	for _, line := range lines {
		var e initEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		if e.Type == "system" && e.Subtype == "init" {
			return e.Tools
		}
	}
	return nil
}

// initEventSessionID returns the FO's session id from the stream's `system`/`init`
// event — the leadSessionId reconcile auto-discovery matches against (set by Claude
// Code at launch, the same value $CLAUDE_CODE_SESSION_ID carries). The merged lane
// reads it to run reconcile under the FO's own session identity and to locate the
// FO's subagents dir.
func initEventSessionID(lines []string) string {
	type initEntry struct {
		Type      string `json:"type"`
		Subtype   string `json:"subtype"`
		SessionID string `json:"session_id"`
	}
	for _, line := range lines {
		var e initEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		if e.Type == "system" && e.Subtype == "init" {
			return e.SessionID
		}
	}
	return ""
}

// stringInSlice is a tiny membership helper for the init-tools assertions.
func stringInSlice(s string, ss []string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

// mergedMemberMeta is the in-process teammate record Claude Code writes for a
// merged background Agent dispatch: projects/<encoded-cwd>/<session-id>/subagents/
// agent-*.meta.json. It carries agentType + the member name and NO team_name — the
// merged-host analog of a legacy team config.json members[] entry, which is NOT
// written on a merged host (verified: no teams/ dir is created headless). agentType
// is the same spacedock:ensign discriminator the reconcile roster loader keys on.
type mergedMemberMeta struct {
	AgentType   string `json:"agentType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TeamName    string `json:"team_name"`
}

func hasMergedEnsignMember(metas []mergedMemberMeta) bool {
	for _, m := range metas {
		if m.AgentType == "spacedock:ensign" && m.Name != "" && m.TeamName == "" {
			return true
		}
	}
	return false
}

// readMergedMemberMetas reads every subagents/agent-*.meta.json under the FO's
// session dir and returns the decoded records. The session dir is
// {configDir}/projects/{encodeProjectDir(resolvedCwd)}/{sessionID}/subagents. A
// missing dir yields an empty slice (the caller asserts the expected member is
// present, so absence fails loudly there with the scanned path).
func readMergedMemberMetas(configDir, resolvedCwd, sessionID string) ([]mergedMemberMeta, error) {
	subagentsDir := filepath.Join(configDir, "projects", encodeProjectDir(resolvedCwd), sessionID, "subagents")
	matches, err := filepath.Glob(filepath.Join(subagentsDir, "agent-*.meta.json"))
	if err != nil {
		return nil, err
	}
	var out []mergedMemberMeta
	for _, p := range matches {
		raw, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var m mergedMemberMeta
		if json.Unmarshal(raw, &m) != nil {
			continue
		}
		out = append(out, m)
	}
	return out, nil
}

// mergedMemberMetasPath returns the subagents dir the metas are read from, for
// failure messages naming exactly where the roster oracle looked.
func mergedMemberMetasPath(configDir, resolvedCwd, sessionID string) string {
	return filepath.Join(configDir, "projects", encodeProjectDir(resolvedCwd), sessionID, "subagents")
}

// TestMergedEnsignDispatchShape is the offline proof of the merged-lane stream
// predicates (no model spend): a merged Agent dispatch is recognized with its name +
// run_in_background + no-team_name shape, a legacy team dispatch is distinguished by
// its team_name, the TeamCreate-absence check is honest, and the init-event tools/
// session-id readers parse the `claude -p` init line. Built from the EXACT shapes the
// live probe captured (the Agent tool_use input + the init event), so the live
// assertions rest on a verified parser.
func TestMergedEnsignDispatchShape(t *testing.T) {
	const mergedAgentLine = `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_1","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"spacedock-ensign-make-it-work-implementation","description":"Make It Work: implementation","run_in_background":true}}]}}`
	const omittedSubagentTypeLine = `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_1b","name":"Agent","input":{"name":"spacedock-ensign-make-it-work-implementation","description":"Make It Work: implementation","run_in_background":true}}]}}`
	const wrongSubagentTypeLine = `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_1c","name":"Agent","input":{"subagent_type":"general-purpose","name":"spacedock-ensign-make-it-work-implementation","description":"Make It Work: implementation","run_in_background":true}}]}}`
	const legacyAgentLine = `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_2","name":"Agent","input":{"subagent_type":"spacedock:ensign","name":"spacedock-ensign-make-it-work-implementation","team_name":"proj-dir-20260618-spacedock","description":"x"}}]}}`
	const teamCreateLine = `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_3","name":"TeamCreate","input":{}}]}}`
	const initLine = `{"type":"system","subtype":"init","session_id":"9e9ec0b0-b998-4c86-b9c7-c3e70f78108c","tools":["Bash","Read","SendMessage","Skill"]}`

	t.Run("merged Agent dispatch recognized with the merged shape", func(t *testing.T) {
		ds := mergedEnsignDispatches([]string{initLine, mergedAgentLine})
		if len(ds) != 1 {
			t.Fatalf("mergedEnsignDispatches found %d, want 1", len(ds))
		}
		d := ds[0]
		if d.name != "spacedock-ensign-make-it-work-implementation" {
			t.Errorf("name = %q, want the derived ensign name", d.name)
		}
		if !d.runInBackground {
			t.Errorf("run_in_background = false, want true (the worker→lead inter-agent communication)")
		}
		if d.hasTeamName {
			t.Errorf("merged dispatch carries a team_name, want none")
		}
	})

	t.Run("legacy team dispatch distinguished by its team_name", func(t *testing.T) {
		ds := mergedEnsignDispatches([]string{legacyAgentLine})
		if len(ds) != 1 || !ds[0].hasTeamName {
			t.Fatalf("a team_name-bearing dispatch must report hasTeamName=true; got %+v", ds)
		}
	})

	t.Run("defaulted subagent_type may be omitted from merged transport", func(t *testing.T) {
		ds := mergedEnsignDispatches([]string{omittedSubagentTypeLine})
		if len(ds) != 1 || ds[0].name == "" || !ds[0].runInBackground || ds[0].hasTeamName {
			t.Fatalf("omitted subagent_type merged dispatch was not recognized: %+v", ds)
		}
	})

	t.Run("explicit non-ensign subagent_type is rejected", func(t *testing.T) {
		if ds := mergedEnsignDispatches([]string{wrongSubagentTypeLine}); len(ds) != 0 {
			t.Fatalf("explicit non-ensign dispatch must not be recognized: %+v", ds)
		}
	})

	t.Run("TeamCreate tool_use detected, absent when not present", func(t *testing.T) {
		if !streamHasTeamCreateToolUse([]string{teamCreateLine}) {
			t.Error("streamHasTeamCreateToolUse missed a TeamCreate tool_use")
		}
		if streamHasTeamCreateToolUse([]string{mergedAgentLine, initLine}) {
			t.Error("streamHasTeamCreateToolUse false-positived on a merged-only stream")
		}
	})

	t.Run("init event tools + session id parsed", func(t *testing.T) {
		tools := initEventToolNames([]string{initLine, mergedAgentLine})
		if stringInSlice("TeamCreate", tools) {
			t.Error("init tools should not contain TeamCreate on a merged host")
		}
		if !stringInSlice("SendMessage", tools) {
			t.Error("init tools should contain SendMessage")
		}
		if id := initEventSessionID([]string{initLine}); id != "9e9ec0b0-b998-4c86-b9c7-c3e70f78108c" {
			t.Errorf("initEventSessionID = %q, want the init session_id", id)
		}
	})

	t.Run("subagent meta decoded as the merged member record", func(t *testing.T) {
		dir := t.TempDir()
		configDir := filepath.Join(dir, "config")
		resolvedCwd := "/private/tmp/fixture"
		sessionID := "9e9ec0b0-b998-4c86-b9c7-c3e70f78108c"
		metaDir := filepath.Join(configDir, "projects", encodeProjectDir(resolvedCwd), sessionID, "subagents")
		if err := os.MkdirAll(metaDir, 0o755); err != nil {
			t.Fatal(err)
		}
		meta := `{"agentType":"spacedock:ensign","name":"spacedock-ensign-make-it-work-implementation","description":"Make It Work: implementation","toolUseId":"toolu_1"}`
		if err := os.WriteFile(filepath.Join(metaDir, "agent-abc123.meta.json"), []byte(meta), 0o644); err != nil {
			t.Fatal(err)
		}
		metas, err := readMergedMemberMetas(configDir, resolvedCwd, sessionID)
		if err != nil {
			t.Fatal(err)
		}
		if len(metas) != 1 {
			t.Fatalf("readMergedMemberMetas found %d, want 1 (looked under %s)", len(metas), mergedMemberMetasPath(configDir, resolvedCwd, sessionID))
		}
		if metas[0].AgentType != "spacedock:ensign" {
			t.Errorf("member agentType = %q, want spacedock:ensign", metas[0].AgentType)
		}
		if metas[0].Name == "" {
			t.Errorf("member name empty, want the derived ensign name")
		}
		if metas[0].TeamName != "" {
			t.Errorf("merged member meta carries a team_name %q, want none", metas[0].TeamName)
		}
		if !strings.HasPrefix(metas[0].Name, "spacedock-ensign-") {
			t.Errorf("member name %q does not look like a derived ensign name", metas[0].Name)
		}
		if !hasMergedEnsignMember(metas) {
			t.Fatal("complete member meta must independently prove ensign identity")
		}
		metas[0].AgentType = "general-purpose"
		if hasMergedEnsignMember(metas) {
			t.Fatal("member meta without agentType=spacedock:ensign must fail identity proof")
		}
	})

}
