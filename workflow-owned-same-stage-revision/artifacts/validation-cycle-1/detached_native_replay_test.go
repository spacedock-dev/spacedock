package ensigncycle

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetachedEightNativeReplays(t *testing.T) {
	evidence := "/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/workflow-owned-same-stage-revision/artifacts/implementation-cycle-1"
	seen := map[string]bool{}
	for _, name := range []string{"before", "after", "review-required", "separate-review-required", "round-required", "round-missing", "cycle-limit", "conventional"} {
		t.Run(name, func(t *testing.T) {
			root := filepath.Join(evidence, name)
			public := readFile(t, filepath.Join(root, "codex-exec.jsonl"))
			raw := readFile(t, filepath.Join(root, "codex-native-lifecycle.jsonl"))
			if !strings.HasPrefix(raw, public+"\n") {
				t.Fatal("public prefix differs")
			}
			var id string
			for _, line := range strings.Split(public, "\n") {
				var x struct{ Type, ThreadID string }
				var m map[string]any
				_ = x
				if json.Unmarshal([]byte(line), &m) == nil && m["type"] == "thread.started" {
					id, _ = m["thread_id"].(string)
				}
			}
			if id == "" || seen[id] {
				t.Fatal("missing/repeated native drive identity", id)
			}
			seen[id] = true
			native := strings.TrimPrefix(raw, public+"\n")
			var first struct {
				Type    string
				Payload struct{ ID string }
			}
			if err := json.Unmarshal([]byte(strings.Split(native, "\n")[0]), &first); err != nil || first.Type != "session_meta" || first.Payload.ID != id {
				t.Fatalf("parent source correlation failed %v %+v", err, first)
			}
			routes := codexRejectionRoutes(raw)
			out := t.TempDir()
			writeRejectionTopologyDigest(t, out, codexRejectionBranch, routes)
			if readFile(t, filepath.Join(out, "rejection-topology.tsv")) != readFile(t, filepath.Join(root, "rejection-topology.tsv")) {
				t.Fatal("derived topology differs")
			}
			if name == "conventional" {
				if err := assertRejectionWorkerTopology(codexRejectionBranch, routes); err != nil {
					t.Fatal(err)
				}
				state, _ := filepath.Abs(filepath.Join("..", "..", ".retained", name))
				entity := filepath.Join(state, "rejection-task", "index.md")
				if err := assertRejectionGatePrepared(entity); err != nil {
					t.Fatal(err)
				}
				if err := assertRejectionCycleLine(entity); err != nil {
					t.Fatal(err)
				}
				if err := assertRejectionRecordedRound(state, entity, "validation", codexRecordedRejectionRound(public)); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := assertSameStageWorkers(routes, len(routes) == 4); err != nil {
					t.Fatal(err)
				}
			}
			// Source evidence must carry completion: deleting only attributed completion events must make the exact topology incomplete.
			lines := strings.Split(raw, "\n")
			var modified []string
			for _, line := range lines {
				var e struct {
					Payload struct {
						Type    string
						Content json.RawMessage
					}
				}
				json.Unmarshal([]byte(line), &e)
				if e.Payload.Type == "agent_message" && strings.Contains(string(e.Payload.Content), "FINAL_ANSWER") && strings.Contains(string(e.Payload.Content), "Done:") {
					continue
				}
				modified = append(modified, line)
			}
			altered := codexRejectionRoutes(strings.Join(modified, "\n"))
			if len(altered) >= len(routes) {
				t.Fatal("completion falsifier did not remove native evidence")
			}
			if _, err := parseRejectionRounds(altered); err == nil {
				t.Fatal("missing completion accepted")
			}
			t.Logf("thread=%s exact routes=%d native source replay PASS", id, len(routes))
		})
	}
}
