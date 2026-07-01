// ABOUTME: Bridge egress-contract conformance — the Claude adapter (hooks.json +
// ABOUTME: spacedock-bridge-events.sh) must turn its host-shaped hook payload into the
// ABOUTME: harness-neutral egress contract: a canonical events.jsonl liveness line plus a
// ABOUTME: deterministic session→entity(+workflow) marker. DRC harness-agnostic FO events.
package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	bridgeAdapterBin     string
	bridgeAdapterBinErr  error
	bridgeAdapterBinOnce sync.Once
)

// runClaudeAdapter feeds ONE Claude-Code-shaped hook payload — the Claude adapter's input —
// to scripts/spacedock-bridge-events.sh (the Claude binding of the egress producer) and
// returns once it exits. cwd anchors the _bridge/ dir the adapter writes to.
func runClaudeAdapter(t *testing.T, payload string) {
	t.Helper()
	hook := filepath.Join("..", "..", "scripts", "spacedock-bridge-events.sh")
	if _, err := os.Stat(hook); err != nil {
		t.Fatalf("Claude adapter script not found at %s: %v", hook, err)
	}
	cmd := exec.Command("bash", hook)
	cmd.Stdin = strings.NewReader(payload)
	cmd.Env = append(os.Environ(), "SPACEDOCK_BIN="+bridgeAdapterBinary(t))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Claude adapter failed: %v\n%s", err, out)
	}
}

func bridgeAdapterBinary(t *testing.T) string {
	t.Helper()
	bridgeAdapterBinOnce.Do(func() {
		dir, err := os.MkdirTemp("", "spacedock-bridge-adapter-bin-*")
		if err != nil {
			bridgeAdapterBinErr = err
			return
		}
		bridgeAdapterBin = filepath.Join(dir, "spacedock")
		cmd := exec.Command("go", "build", "-o", bridgeAdapterBin, "./cmd/spacedock")
		cmd.Dir = repoRoot(t)
		if out, err := cmd.CombinedOutput(); err != nil {
			bridgeAdapterBinErr = err
			_ = os.RemoveAll(dir)
			bridgeAdapterBin = ""
			t.Logf("build output:\n%s", out)
		}
	})
	if bridgeAdapterBinErr != nil {
		t.Fatalf("build spacedock bridge adapter binary: %v", bridgeAdapterBinErr)
	}
	return bridgeAdapterBin
}

// claudeAdapterRead builds the Claude-Code-shaped PostToolUse/Read payload that Claude Code
// pipes to the hook. This CC-specific JSON is the Claude adapter's INPUT only — it is NOT the
// contract. The contract is what the adapter emits (asserted below); a future Codex or Pi
// producer would consume its own host payload but must emit the same two output shapes.
func claudeAdapterRead(cwd, sid, agentType, filePath string) string {
	return `{"cwd":"` + cwd + `","session_id":"` + sid + `","agent_type":"` + agentType +
		`","hook_event_name":"PostToolUse","tool_name":"Read","tool_input":{"file_path":"` + filePath + `"}}`
}

// egressLine is the Spacedock-owned, harness-neutral events.jsonl contract line every host
// adapter must produce. Pointer fields distinguish "key absent" (nil) from "present but empty",
// and the nested detail struct enforces detail.{tool,source} nesting — a producer that flattens
// tool/source to the top level fails to populate detail and is rejected.
type egressLine struct {
	TS        string  `json:"ts"`
	Event     string  `json:"event"`
	SessionID string  `json:"session_id"`
	AgentID   *string `json:"agent_id"`
	AgentType *string `json:"agent_type"`
	Detail    *struct {
		Tool   *string `json:"tool"`
		Source *string `json:"source"`
	} `json:"detail"`
}

// assertEgressContractLine parses the LAST line of «root»/_bridge/events.jsonl and asserts it
// conforms to the harness-neutral egress contract — valid JSON, all canonical keys present with
// detail.{tool,source} genuinely NESTED, and the load-bearing values mapped correctly (event,
// session_id, agent_type — the field Bridge uses to tell FO from ensign — and a non-empty ts,
// which Bridge reads for freshness). Parsing (not substring matching) is what makes this a real
// cross-host guardrail: a future Codex/Pi producer emitting the same shape passes; one that
// drops a key, mis-nests detail, or mis-maps a value fails.
func assertEgressContractLine(t *testing.T, root, wantEvent, wantSession, wantAgentType, wantTool string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "events.jsonl"))
	if err != nil {
		t.Fatalf("egress contract line not written: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	raw := lines[len(lines)-1]
	var line egressLine
	if err := json.Unmarshal([]byte(raw), &line); err != nil {
		t.Fatalf("egress line is not valid JSON: %v\n%s", err, raw)
	}
	if line.AgentID == nil || line.AgentType == nil {
		t.Fatalf("egress line missing agent_id/agent_type keys\ngot: %s", raw)
	}
	if line.Detail == nil || line.Detail.Tool == nil || line.Detail.Source == nil {
		t.Fatalf("egress line must nest detail.{tool,source} (not flatten them to the top level)\ngot: %s", raw)
	}
	if line.TS == "" {
		t.Errorf("egress line has empty ts (Bridge reads freshness from ts)\ngot: %s", raw)
	}
	if line.Event != wantEvent {
		t.Errorf("event = %q, want %q\ngot: %s", line.Event, wantEvent, raw)
	}
	if line.SessionID != wantSession {
		t.Errorf("session_id = %q, want %q\ngot: %s", line.SessionID, wantSession, raw)
	}
	if *line.AgentType != wantAgentType {
		t.Errorf("agent_type = %q, want %q (FO-vs-ensign attribution)\ngot: %s", *line.AgentType, wantAgentType, raw)
	}
	if *line.Detail.Tool != wantTool {
		t.Errorf("detail.tool = %q, want %q\ngot: %s", *line.Detail.Tool, wantTool, raw)
	}
}

// assertSessionMarker parses «root»/_bridge/sessions/<sid>.json and asserts the marker contract
// {session_id,entity,workflow}. Parse-based (not substring) so a marker that records the state
// dir as the workflow, or mis-maps a field, fails. Reusable by a future Codex/Pi producer test.
func assertSessionMarker(t *testing.T, root, sid, wantEntity, wantWorkflow string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "sessions", sid+".json"))
	if err != nil {
		t.Fatalf("marker for %s not written: %v", sid, err)
	}
	var m struct {
		SessionID string `json:"session_id"`
		Entity    string `json:"entity"`
		Workflow  string `json:"workflow"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("marker for %s is not valid JSON: %v\n%s", sid, err, data)
	}
	if m.SessionID != sid {
		t.Errorf("marker session_id = %q, want %q", m.SessionID, sid)
	}
	if m.Entity != wantEntity {
		t.Errorf("marker entity = %q, want %q", m.Entity, wantEntity)
	}
	if m.Workflow != wantWorkflow {
		t.Errorf("marker workflow = %q, want %q", m.Workflow, wantWorkflow)
	}
}

// markerExists reports whether a session marker file was written for sid.
func markerExists(root, sid string) bool {
	_, err := os.Stat(filepath.Join(root, "_bridge", "sessions", sid+".json"))
	return err == nil
}

// TestClaudeAdapterConformsToEgressContract is the Claude adapter's conformance test against
// the harness-neutral egress contract. The unit under test is the CONTRACT, not the raw Claude
// Code payload: each step feeds the Claude adapter its host-shaped input and asserts the
// adapter's OUTPUT — (a) a canonical events.jsonl liveness line and (b) the session→entity
// marker {session_id,entity,workflow} — matches the Spacedock-owned shape. A future Codex or Pi
// producer is the same test with a different input builder reusing assertEgressContractLine /
// assertSessionMarker, and must satisfy these same OUTPUT assertions; that is what
// "harness-agnostic egress" means.
func TestClaudeAdapterConformsToEgressContract(t *testing.T) {
	root := t.TempDir()
	ent := filepath.Join(root, "docs", "spacedock", "linear-drc-review", "drc-3467.md")

	// 1. Ensign reads its own entity file. Assert BOTH contract outputs:
	//    (a) the canonical egress liveness line, and (b) the session→entity marker.
	runClaudeAdapter(t, claudeAdapterRead(root, "ses-1", "spacedock:ensign", ent))
	assertEgressContractLine(t, root, "PostToolUse", "ses-1", "spacedock:ensign", "Read")
	assertSessionMarker(t, root, "ses-1", "drc-3467", "linear-drc-review")

	// 2. First-write-wins: a later sibling Read in the same session must NOT overwrite.
	sibling := filepath.Join(root, "docs", "spacedock", "linear-drc-review", "drc-9999.md")
	runClaudeAdapter(t, claudeAdapterRead(root, "ses-1", "spacedock:ensign", sibling))
	assertSessionMarker(t, root, "ses-1", "drc-3467", "linear-drc-review")

	// 3. A non-ensign (FO) Read of an entity file writes no marker. The egress liveness line
	//    is still emitted (the FO is live) and carries the FO's agent_type, but no
	//    session→entity link is recorded.
	runClaudeAdapter(t, claudeAdapterRead(root, "fo-sess", "spacedock:first-officer", ent))
	assertEgressContractLine(t, root, "PostToolUse", "fo-sess", "spacedock:first-officer", "Read")
	if markerExists(root, "fo-sess") {
		t.Errorf("FO Read should not produce a session marker")
	}

	// 3b. RELATIVE entity path — the FO passes a repo-relative {entity_file_path}, so the
	// ensign's scoped Read carries "docs/spacedock/<wf>/<slug>.md" (no leading slash). The
	// adapter must still record it (regression: the absolute-only pattern missed every live
	// ensign).
	runClaudeAdapter(t, claudeAdapterRead(root, "ses-rel", "spacedock:ensign",
		"docs/spacedock/linear-drc-review/drc-7000.md"))
	assertSessionMarker(t, root, "ses-rel", "drc-7000", "linear-drc-review")

	// 4. _archive entity Reads are skipped.
	arch := filepath.Join(root, "docs", "spacedock", "linear-drc-review", "_archive", "drc-1.md")
	runClaudeAdapter(t, claudeAdapterRead(root, "ses-arch", "spacedock:ensign", arch))
	if markerExists(root, "ses-arch") {
		t.Errorf("_archive Read should not produce a session marker")
	}

	// 5. SPLIT-ROOT entity path: the entity now lives at <wf>/.spacedock-state/<slug>.md,
	// so the workflow must be derived from the segment after docs/spacedock/ — NOT the
	// entity's parent dir (which would wrongly be ".spacedock-state"). assertSessionMarker's
	// exact workflow match ("linear-drc-review") rejects the ".spacedock-state" mis-derivation.
	sr := filepath.Join(root, "docs", "spacedock", "linear-drc-review", ".spacedock-state", "drc-8000.md")
	runClaudeAdapter(t, claudeAdapterRead(root, "ses-sr", "spacedock:ensign", sr))
	assertSessionMarker(t, root, "ses-sr", "drc-8000", "linear-drc-review")

	// 6. SPLIT-ROOT _archive: an archived entity in the state checkout is still skipped.
	srArch := filepath.Join(root, "docs", "spacedock", "linear-drc-review", ".spacedock-state", "_archive", "drc-9.md")
	runClaudeAdapter(t, claudeAdapterRead(root, "ses-srarch", "spacedock:ensign", srArch))
	if markerExists(root, "ses-srarch") {
		t.Errorf("split-root _archive Read should not produce a session marker")
	}
}
