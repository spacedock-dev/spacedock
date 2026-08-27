package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// TestGateOneFileRoomJourneyPublishesAndCompletes is value AC-1. A gate journey
// publishes one file, and the gate still closes, consumes, and archives. Before
// this change the same journey published two files and bound a request-digest,
// in 499 of 499 rooms scanned. Make Prepare mint the request again and the
// entry-set assertion reds.
func TestGateOneFileRoomJourneyPublishesAndCompletes(t *testing.T) {
	workflow, state, artifact := gateOneFileCLIFixture(t)
	entity := filepath.Join(state, "task.md")

	var out, errOut bytes.Buffer
	invoke := func(args ...string) int {
		out.Reset()
		errOut.Reset()
		return run(context.Background(), args, nil, workflow, nil, &out, &errOut, &status.NativeRunner{}, nil)
	}

	if code := invoke("gate", "prepare", "task",
		"--question", "Advance?", "--artifact", artifact, "--summary", "one-file candidate",
		"--workflow-dir", workflow); code != 0 {
		t.Fatalf("prepare exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	room := filepath.Join(state, "task", "review", "validation", "briefing-1")
	if got := gateRoomEntryNames(t, room); !reflect.DeepEqual(got, []string{"index.json"}) {
		t.Fatalf("room entries=%v want exactly [index.json]", got)
	}
	doc, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	binding := doc.Records[0].Attempts[0].Briefing
	if binding.RequestDigest != "" {
		t.Fatalf("binding carries request-digest %q", binding.RequestDigest)
	}
	if binding.RoomRef != "@review/validation/briefing-1" {
		t.Fatalf("room-ref=%q", binding.RoomRef)
	}

	if code := invoke("gate", "record", "task", "--decision", "approve",
		"--actor", "person:captain", "--consume", "--workflow-dir", workflow); code != 0 {
		t.Fatalf("close and consume exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "state=closed") || !strings.Contains(out.String(), "consumed=") {
		t.Fatalf("close and consume stdout=%q want state=closed and consumed=", out.String())
	}
	if got := status.ParseFrontmatter(entity)["status"]; got != "implementation" {
		t.Fatalf("consumed approval left status=%q want implementation", got)
	}
	doc, _, err = gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	attempt := doc.Records[0].Attempts[0]
	if attempt.Resolution == nil || attempt.Resolution.Decision != "approve" ||
		attempt.Application == nil || attempt.Application.State != "consumed" {
		t.Fatalf("attempt did not close and consume: %#v", attempt)
	}
	// The archived one-file room stays readable to the read-only validator.
	if got := gateRoomEntryNames(t, room); !reflect.DeepEqual(got, []string{"index.json"}) {
		t.Fatalf("archived room entries=%v", got)
	}
	if code := invoke("status", "--workflow-dir", workflow, "--validate"); code != 0 {
		t.Fatalf("validate over the archived room exit=%d stderr=%q", code, errOut.String())
	}
}

func gateRoomEntryNames(t *testing.T, room string) []string {
	t.Helper()
	entries, err := os.ReadDir(room)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			t.Fatalf("gate room holds non-regular entry %s", entry.Name())
		}
		names = append(names, entry.Name())
	}
	return names
}

// gateOneFileCLIFixture is gatePrepareCLIFixture with a nonterminal advance
// target, so an approved gate really consumes instead of parking on the terminal
// merge route. Only the stage list differs, so it reuses the shared fixture and
// re-commits the one changed file.
func gateOneFileCLIFixture(t *testing.T) (workflow, state, artifact string) {
	t.Helper()
	workflow, state, artifact = gatePrepareCLIFixture(t)
	readme := filepath.Join(workflow, "README.md")
	body, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	staged := bytes.Replace(body, []byte("    - name: done\n"), []byte("    - name: implementation\n    - name: done\n"), 1)
	if bytes.Equal(staged, body) {
		t.Fatal("shared fixture lost the done stage this helper inserts before")
	}
	if err := os.WriteFile(readme, staged, 0o644); err != nil {
		t.Fatal(err)
	}
	mainRoot := filepath.Dir(filepath.Dir(workflow))
	git(t, mainRoot, "add", "-A")
	git(t, mainRoot, "commit", "-q", "-m", "nonterminal advance target")
	return workflow, state, artifact
}
