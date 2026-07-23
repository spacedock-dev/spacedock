package gates

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRoundRecordPrefixReplayAndRefusalsAreByteClean(t *testing.T) {
	root, entity, briefing, log, feedback := advisoryRoundFixture(t)
	fullLog := mustReadBytes(t, log)
	prefixEnd := nthNewline(fullLog, 3)
	if err := os.WriteFile(log, fullLog[:prefixEnd], 0o644); err != nil {
		t.Fatal(err)
	}

	input := RecordInput{
		Round:             "implementation/1",
		BriefingPath:      briefing,
		LogPath:           log,
		FeedbackCyclePath: feedback,
	}
	beforeCandidate := mustReadBytes(t, filepath.Join(root, "candidate.patch"))
	beforeProduct := mustReadBytes(t, filepath.Join(root, "product", "status.txt"))
	beforeLifecycle := lifecycleBytes(t, entity)
	if err := RecordSemantic(entity, input); err != nil {
		t.Fatalf("record reviewer prefix: %v", err)
	}
	prefixSummary, err := ValidateRoundFile(entity, input.Round)
	if err != nil {
		t.Fatalf("validate reviewer prefix: %v", err)
	}
	if prefixSummary.Triage != "pending" || len(prefixSummary.Entries) != 3 {
		t.Fatalf("reviewer prefix summary = %#v, want pending three-entry round", prefixSummary)
	}
	if body := mustReadBytes(t, entity); bytes.Contains(body, []byte("- Cycle 1:")) {
		t.Fatal("reviewer-only prefix projected a completed Feedback Cycles line")
	}

	if err := os.WriteFile(log, fullLog, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RecordSemantic(entity, input); err != nil {
		t.Fatalf("append worker triage: %v", err)
	}
	summary, err := ValidateRoundFile(entity, input.Round)
	if err != nil {
		t.Fatalf("validate complete round: %v", err)
	}
	if summary.ID != "round:task:implementation:1" || summary.Triage != "all-declines" || len(summary.Entries) != 5 {
		t.Fatalf("complete summary = %#v", summary)
	}
	resolutions := 0
	for _, entry := range summary.Entries {
		if entry.Type == "Resolution" {
			resolutions++
			if !entry.Advisory {
				t.Fatalf("round Resolution %s was not advisory", entry.ID)
			}
		}
	}
	if resolutions != 2 {
		t.Fatalf("resolution count = %d, want reviewer and worker", resolutions)
	}
	entityBytes := mustReadBytes(t, entity)
	if got := bytes.Count(entityBytes, []byte("- Cycle 1:")); got != 1 {
		t.Fatalf("Feedback Cycles projection count = %d, want one", got)
	}
	room := filepath.Join(root, "review", "implementation", "round-1")
	if got := mustReadBytes(t, filepath.Join(room, "briefing.review.jsonl")); !bytes.Equal(got, fullLog) {
		t.Fatal("retained strict-prefix append does not equal supplied log")
	}
	if got := mustReadBytes(t, filepath.Join(room, "briefing.json")); !bytes.Equal(got, mustReadBytes(t, briefing)) {
		t.Fatal("retained Briefing bytes changed")
	}
	if !bytes.Equal(mustReadBytes(t, filepath.Join(root, "candidate.patch")), beforeCandidate) ||
		!bytes.Equal(mustReadBytes(t, filepath.Join(root, "product", "status.txt")), beforeProduct) {
		t.Fatal("round recording changed candidate or product bytes")
	}
	if got := lifecycleBytes(t, entity); !bytes.Equal(got, beforeLifecycle) {
		t.Fatalf("round recording changed status or gates:\n%s", got)
	}

	if err := os.RemoveAll(filepath.Join(root, ".derived-cache")); err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateRoundFile(entity, input.Round); err != nil {
		t.Fatalf("pointer readback after cache removal: %v", err)
	}
	replayBefore := treeDigest(t, root)
	if err := RecordSemantic(entity, input); err != nil {
		t.Fatalf("exact replay: %v", err)
	}
	if replayAfter := treeDigest(t, root); replayAfter != replayBefore {
		t.Fatal("exact replay changed fixture bytes")
	}

	divergent := bytes.Replace(fullLog, []byte("job 592"), []byte("job 591"), 1)
	if err := os.WriteFile(log, divergent, 0o644); err != nil {
		t.Fatal(err)
	}
	divergentBefore := treeDigest(t, root)
	if err := RecordSemantic(entity, input); err == nil {
		t.Fatal("divergent replay succeeded")
	}
	if got := treeDigest(t, root); got != divergentBefore {
		t.Fatal("divergent replay changed fixture bytes")
	}

	if err := os.WriteFile(log, fullLog, 0o644); err != nil {
		t.Fatal(err)
	}
	changedBriefing := append(mustReadBytes(t, briefing), '\n')
	if err := os.WriteFile(briefing, changedBriefing, 0o644); err != nil {
		t.Fatal(err)
	}
	briefingBefore := treeDigest(t, root)
	if err := RecordSemantic(entity, input); err == nil {
		t.Fatal("changed Briefing bytes replay succeeded")
	}
	if got := treeDigest(t, root); got != briefingBefore {
		t.Fatal("changed Briefing refusal changed fixture bytes")
	}
}

func TestRoundNoFindingsAndPreflightRefusals(t *testing.T) {
	t.Run("no findings", func(t *testing.T) {
		_, entity, briefing, log, feedback := advisoryRoundFixture(t)
		noFindings := `{"type":"Resolution","id":"resolution:roborev-clear","briefing":"briefing:3j:implementation:round-1","by":"software:roborev","at":"2026-07-20T01:00:00Z","decision":"approve"}` + "\n"
		if err := os.WriteFile(log, []byte(noFindings), 0o644); err != nil {
			t.Fatal(err)
		}
		input := RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback}
		if err := RecordSemantic(entity, input); err != nil {
			t.Fatal(err)
		}
		summary, err := ValidateRoundFile(entity, input.Round)
		if err != nil {
			t.Fatal(err)
		}
		if summary.Triage != "no-findings" {
			t.Fatalf("triage = %q, want no-findings", summary.Triage)
		}
		if bytes.Contains(mustReadBytes(t, entity), []byte("- Cycle 1:")) {
			t.Fatal("no-findings round received worker-triage projection")
		}
	})

	t.Run("bad artifact digest", func(t *testing.T) {
		root, entity, briefing, log, feedback := advisoryRoundFixture(t)
		if err := os.WriteFile(filepath.Join(root, "candidate.patch"), []byte("corrupt\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, root)
		err := RecordSemantic(entity, RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback})
		if err == nil || !strings.Contains(err.Error(), "artifact") {
			t.Fatalf("bad digest error = %v", err)
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("bad digest refusal changed fixture tree")
		}
		if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
			t.Fatalf("bad digest left lock residue: %v", err)
		}
	})

	t.Run("occupied target", func(t *testing.T) {
		root, entity, briefing, log, feedback := advisoryRoundFixture(t)
		room := filepath.Join(root, "review", "implementation", "round-1")
		if err := os.MkdirAll(room, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(room, "foreign"), []byte("occupied"), 0o644); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, root)
		if err := RecordSemantic(entity, RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback}); err == nil {
			t.Fatal("conflicting occupied target succeeded")
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("occupied-target refusal changed fixture tree")
		}
	})

	t.Run("lock contention", func(t *testing.T) {
		root, entity, briefing, log, feedback := advisoryRoundFixture(t)
		if err := os.WriteFile(entity+".gates.lock", []byte("held"), 0o600); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, root)
		if err := RecordSemantic(entity, RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback}); err == nil {
			t.Fatal("lock-contended record succeeded")
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("lock contention changed fixture tree")
		}
	})

	t.Run("cross briefing log", func(t *testing.T) {
		root, entity, briefing, log, feedback := advisoryRoundFixture(t)
		body := bytes.Replace(mustReadBytes(t, log), []byte(`"briefing":"briefing:3j:implementation:round-1"`), []byte(`"briefing":"briefing:other"`), 1)
		if err := os.WriteFile(log, body, 0o644); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, root)
		if err := RecordSemantic(entity, RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback}); err == nil {
			t.Fatal("cross-Briefing log succeeded")
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("cross-Briefing refusal changed fixture tree")
		}
	})
}

func advisoryRoundFixture(t *testing.T) (root, entity, briefing, log, feedback string) {
	t.Helper()
	root = t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "inputs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "product"), 0o755); err != nil {
		t.Fatal(err)
	}
	copyRoundFixture(t, filepath.Join(root, "candidate.patch"), "candidate.patch")
	briefing = filepath.Join(root, "inputs", "briefing.json")
	log = filepath.Join(root, "inputs", "briefing.review.jsonl")
	copyRoundFixture(t, briefing, "briefing.json")
	copyRoundFixture(t, log, "briefing.review.jsonl")
	feedback = filepath.Join(root, "inputs", "feedback-cycle.txt")
	if err := os.WriteFile(feedback, []byte("- Cycle 1: REJECTED — Roborev; surface 0/0 vs estimate 340 (0%); AC unchanged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "product", "status.txt"), []byte("candidate=90aea55\nstate=unchanged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entity = filepath.Join(root, "task.md")
	body := "---\n" +
		"id: \"\"\n" +
		"status: implementation\n" +
		"custom: preserve-me\n" +
		"gates:\n" +
		"  version: 1\n" +
		"  current:\n" +
		"    gate: gate:task:ideation\n" +
		"  records:\n" +
		"    - id: gate:task:ideation\n" +
		"      stage: ideation\n" +
		"      attempts:\n" +
		"        - id: gate-attempt:task-ideation-1\n" +
		"          briefing:\n" +
		"            id: briefing:task:ideation:1\n" +
		"            digest: sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n" +
		"            digest-domain: canonical-bytes\n" +
		"            room-ref: ./review/ideation/briefing-1\n" +
		"title: Task\n" +
		"---\n" +
		"# Task\n\nUnrelated body bytes stay fixed.\n"
	if err := os.WriteFile(entity, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, entity, briefing, log, feedback
}

func copyRoundFixture(t *testing.T, destination, name string) {
	t.Helper()
	body := mustReadBytes(t, filepath.Join("testdata", "advisory-round", name))
	if err := os.WriteFile(destination, body, 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustReadBytes(t *testing.T, path string) []byte {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func nthNewline(body []byte, n int) int {
	seen := 0
	for i, b := range body {
		if b == '\n' {
			seen++
			if seen == n {
				return i + 1
			}
		}
	}
	return len(body)
}

func lifecycleBytes(t *testing.T, entity string) []byte {
	t.Helper()
	body := mustReadBytes(t, entity)
	root, _, _, err := frontmatterNode(body)
	if err != nil {
		t.Fatal(err)
	}
	status, gates := mappingValue(root, "status"), mappingValue(root, "gates")
	gatesBytes, err := yaml.Marshal(gates)
	if err != nil {
		t.Fatal(err)
	}
	return append([]byte("status:"+status.Value+"\n"), gatesBytes...)
}

func treeDigest(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		hash.Write([]byte(rel))
		hash.Write([]byte{0})
		if entry.IsDir() {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hash.Write(body)
		hash.Write([]byte{0})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
