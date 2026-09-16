package ensigncycle

import (
	"encoding/json"
	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/gitsource"
	"os"
	"path/filepath"
	"testing"
)

func TestDetachedSameStageRetainedState(t *testing.T) {
	base := filepath.Join("..", "..", ".retained")
	for _, tc := range []struct {
		name             string
		attempts, rounds int
	}{
		{"before", 1, 0}, {"after", 2, 0}, {"review-required", 1, 0}, {"separate-review-required", 1, 0}, {"cycle-limit", 3, 0}, {"round-missing", 1, 0}, {"round-required", 2, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state, _ := filepath.Abs(filepath.Join(base, tc.name))
			entity := filepath.Join(state, "recorded-gate-task", "index.md")
			records, _, err := gates.Read(entity)
			if err != nil {
				t.Fatal(err)
			}
			if len(records.Records) != 1 || len(records.Records[0].Attempts) != tc.attempts {
				t.Fatalf("wrong attempts: %+v", records)
			}
			attempts := records.Records[0].Attempts
			for i, a := range attempts {
				if i == len(attempts)-1 && tc.attempts == 2 {
					if a.Resolution != nil || a.Application != nil || a.Withdrawal != nil {
						t.Fatal("fresh gate closed")
					}
					room, err := gates.ResolveRoomRef(entity, a.Briefing.RoomRef)
					if err != nil {
						t.Fatal(err)
					}
					data, err := os.ReadFile(filepath.Join(room, "index.json"))
					if err != nil {
						t.Fatal(err)
					}
					if err := assertSameStageSelectedPlan(gitsource.Roots{Main: state, State: state}, data); err != nil {
						t.Fatal(err)
					}
				} else if a.Resolution == nil || a.Resolution.Decision != "revise" || a.Application != nil {
					t.Fatal("old authority changed")
				}
			}
			for _, p := range []string{"plan.md", "frozen-input.txt"} {
				data, err := os.ReadFile(filepath.Join(state, "recorded-gate-task", "selected", p))
				if err != nil || string(data) != sameStagePlan {
					t.Fatalf("%s: %q %v", p, data, err)
				}
			}
			rooms, _ := filepath.Glob(filepath.Join(state, "recorded-gate-task/review/*/round-*"))
			if len(rooms) != tc.rounds {
				t.Fatal("wrong round cardinality", rooms)
			}
			if tc.rounds == 1 {
				summary, err := gates.ValidateRoundFile(entity, "validation/1")
				if err != nil || len(summary.Entries) != 4 {
					t.Fatalf("round: %+v %v", summary, err)
				}
			}
		})
	}
}

func TestDetachedSameStageSelectionBoundaries(t *testing.T) {
	state, _ := filepath.Abs(filepath.Join("..", "..", ".retained", "after"))
	roots := gitsource.Roots{Main: filepath.Dir(state), State: state}
	source, err := gitsource.Inspect(roots, filepath.Join(state, "recorded-gate-task/selected/plan.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, kind, uri, rev string
		pass                 bool
	}{
		{"canonical-artifact", "artifact", source.URI, source.Rev, true},
		{"canonical-reference", "Reference", source.URI, source.Rev, true},
		{"wrong-type", "Note", source.URI, source.Rev, false},
		{"wrong-revision", "artifact", source.URI, "sha256:0000000000000000000000000000000000000000000000000000000000000000", false},
		{"wrong-identity", "artifact", source.URI + "x", source.Rev, false},
		{"unicode-identity", "artifact", source.URI + "\u200b", source.Rev, false},
		{"empty", "artifact", "", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			field := "context"
			if tc.kind == "artifact" {
				field = "artifacts"
			}
			b, _ := json.Marshal(map[string]any{field: []map[string]string{{"type": tc.kind, "uri": tc.uri, "rev": tc.rev}}})
			if err := assertSameStageSelectedPlan(roots, b); (err == nil) != tc.pass {
				t.Fatalf("pass=%v err=%v", tc.pass, err)
			}
		})
	}
	if assertSameStageSelectedPlan(roots, []byte("{\"artifacts\": [")) == nil {
		t.Fatal("truncated JSON accepted")
	}
}

func TestDetachedSameStageWorkerEventOrder(t *testing.T) {
	spawn := rejectionRoute{event: routeSpawn, stage: "validation", target: "fix"}
	done := rejectionRoute{event: routeDone, stage: "validation", target: "fix"}
	other := done
	other.target = "unowned"
	for i, routes := range [][]rejectionRoute{nil, {done, spawn}, {spawn, other}, {spawn, done, done}, {spawn, spawn, done}} {
		if assertSameStageWorkers(routes, false) == nil {
			t.Errorf("bad trace %d accepted", i)
		}
	}
}
