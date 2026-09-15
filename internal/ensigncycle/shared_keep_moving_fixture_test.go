package ensigncycle

import (
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

func TestKeepMovingPreparedGate(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	root := writeKeepMovingWorkflow(t, t.TempDir())
	path := filepath.Join(root, kmApprovedGate+".md")
	prepared := readFile(t, path)
	doc, _, err := gates.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Records) != 1 || doc.Records[0].Stage != "review" || len(doc.Records[0].Attempts) != 1 {
		t.Fatalf("want one prepared review attempt, got %#v", doc)
	}
	attempt := doc.Records[0].Attempts[0]
	if attempt.Briefing.ID == "" || attempt.Briefing.Digest == "" || attempt.Briefing.RoomRef == "" || attempt.Resolution != nil || attempt.Application != nil {
		t.Fatalf("want bound open briefing without approval, got %#v", attempt)
	}
	if durableField(prepared, "status") != "review" || durableField(prepared, "started") != "" {
		t.Fatal("fixture advanced or dispatched before captain approval")
	}
	git(t, root, "cat-file", "-e", "HEAD:approved-gate/review/review/briefing-1/index.json")
	if committed := git(t, root, "show", "HEAD:"+kmApprovedGate+".md"); committed != prepared {
		t.Fatal("prepared entity is not committed")
	}
	for _, commit := range durableEntityHistory(t, root, kmApprovedGate, kmApprovedGate+".md") {
		if durableField(commit.blob, "status") != "review" || durableField(commit.blob, "started") != "" {
			t.Fatal("fixture history already contains successor work")
		}
	}
	refused := runRecordedGateCommand(binary, root, "", "gate", "consume", kmApprovedGate, "--workflow-dir", root)
	if refused.exit != 1 || readFile(t, path) != prepared {
		t.Fatalf("missing approval must refuse without mutation: %#v", refused)
	}
	mustRecordedGate(t, binary, root, "gate", "record", kmApprovedGate, "--decision", "approve", "--actor", "person:captain", "--consume", "--workflow-dir", root)
	doc, _, err = gates.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	attempt = doc.Records[0].Attempts[0]
	if attempt.Resolution == nil || attempt.Resolution.By != "person:captain" || attempt.Resolution.Decision != "approve" || attempt.Application == nil || attempt.Application.State != "consumed" || attempt.Application.TargetStage != kmNextStage || durableField(readFile(t, path), "status") != kmNextStage {
		t.Fatalf("captain approval did not consume to implementation: %#v", attempt)
	}
	t.Run("unprepared approval refuses", func(t *testing.T) {
		root := t.TempDir()
		writeFile(t, filepath.Join(root, "README.md"), keepMovingReadme())
		path := filepath.Join(root, kmApprovedGate+".md")
		writeFile(t, path, keepMovingApprovedEntity())
		gitInit(t, root)
		before := readFile(t, path)
		result := runRecordedGateCommand(binary, root, "", "gate", "record", kmApprovedGate, "--decision", "approve", "--actor", "person:captain", "--consume", "--workflow-dir", root)
		if result.exit != 1 || readFile(t, path) != before {
			t.Fatalf("unprepared approval must refuse without mutation: %#v", result)
		}
	})
}
