package gates

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRoundRecordCompleteReplayAndRefusalsAreByteClean(t *testing.T) {
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
	pendingBefore := treeDigest(t, root)
	if err := RecordSemantic(entity, input); err == nil {
		t.Fatal("findings-bearing reviewer-only log was persisted")
	}
	if got := treeDigest(t, root); got != pendingBefore {
		t.Fatal("incomplete-round refusal changed fixture bytes")
	}
	room := filepath.Join(root, "review", "implementation", "round-1")
	if _, err := os.Stat(room); !os.IsNotExist(err) {
		t.Fatalf("incomplete round left a room: %v", err)
	}

	if err := os.WriteFile(log, fullLog, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RecordSemantic(entity, input); err != nil {
		t.Fatalf("publish complete round: %v", err)
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
	if got := mustReadBytes(t, filepath.Join(room, "briefing.review.jsonl")); !bytes.Equal(got, fullLog) {
		t.Fatal("retained complete log does not equal supplied log")
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
		_, entity, briefing, log, _ := advisoryRoundFixture(t)
		noFindings := `{"type":"Resolution","id":"resolution:roborev-clear","briefing":"briefing:3j:implementation:round-1","by":"software:roborev","at":"2026-07-20T01:00:00Z","decision":"approve"}` + "\n"
		if err := os.WriteFile(log, []byte(noFindings), 0o644); err != nil {
			t.Fatal(err)
		}
		input := RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log}
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

	t.Run("malformed artifact URI", func(t *testing.T) {
		root, entity, briefing, log, feedback := advisoryRoundFixture(t)
		body := bytes.Replace(mustReadBytes(t, briefing), []byte("../../../candidate.patch"), []byte("%gh"), 1)
		if err := os.WriteFile(briefing, body, 0o644); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, root)
		err := RecordSemantic(entity, RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback})
		if err == nil || !strings.Contains(err.Error(), "URI") {
			t.Fatalf("malformed artifact URI error = %v", err)
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("malformed URI refusal changed fixture tree")
		}
	})
}

func TestRoundSharedCASAndRollbackBoundaries(t *testing.T) {
	t.Run("entity expectation is mandatory and full bytes are compared", func(t *testing.T) {
		entity := filepath.Join(t.TempDir(), "entity.md")
		original := []byte("---\nstatus: implementation\n---\n# Entity\n")
		if err := os.WriteFile(entity, original, 0o644); err != nil {
			t.Fatal(err)
		}
		buildCalled := false
		build := func(current []byte) ([]byte, error) {
			buildCalled = true
			return append(current, []byte("changed\n")...), nil
		}
		if err := mutateEntity(entity, entityExpectation{}, build, atomicWrite); err == nil {
			t.Fatal("zero entity expectation succeeded")
		}
		if buildCalled {
			t.Fatal("build ran without a mandatory expectation")
		}
		stale := append([]byte(nil), original...)
		stale[4] = 'X'
		if err := mutateEntity(entity, entityExpectation{Bytes: stale}, build, atomicWrite); err == nil {
			t.Fatal("stale full-entity expectation succeeded")
		}
		if buildCalled || !bytes.Equal(mustReadBytes(t, entity), original) {
			t.Fatal("stale entity CAS built or mutated output")
		}
	})

	t.Run("different retained room is immutable and refused before entity commit", func(t *testing.T) {
		room := filepath.Join(t.TempDir(), "review", "implementation", "round-1")
		if err := os.MkdirAll(room, 0o755); err != nil {
			t.Fatal(err)
		}
		old := roundRoomBytes{Exists: true, Briefing: []byte("briefing"), Log: []byte("old\n")}
		if err := os.WriteFile(filepath.Join(room, "briefing.json"), old.Briefing, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(room, "briefing.review.jsonl"), []byte("raced\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		committed := false
		err := publishRound(room, roundRoomBytes{Exists: true, Briefing: old.Briefing, Log: []byte("next\n")}, func(bool) error {
			committed = true
			return nil
		})
		if err == nil || committed {
			t.Fatalf("stale room err=%v committed=%v", err, committed)
		}
		if got := mustReadBytes(t, filepath.Join(room, "briefing.review.jsonl")); !bytes.Equal(got, []byte("raced\n")) {
			t.Fatal("stale-room refusal changed retained log")
		}
	})

	t.Run("entity failure removes a new room without mutable-room restoration", func(t *testing.T) {
		root := t.TempDir()
		injected := errors.New("injected entity replace failure")
		newRoom := filepath.Join(root, "review", "validation", "round-2")
		fresh := roundRoomBytes{Exists: true, Briefing: []byte("new briefing"), Log: []byte("new log\n")}
		if err := publishRound(newRoom, fresh, func(bool) error { return injected }); !errors.Is(err, injected) {
			t.Fatalf("new-room failure = %v, want injected error", err)
		}
		if _, err := os.Stat(newRoom); !os.IsNotExist(err) {
			t.Fatalf("entity failure left a new room: %v", err)
		}
	})
}

func TestRoundWorkerTriageRequiresFixedActorAndBackwardGraph(t *testing.T) {
	const briefingID = "briefing:test"
	reviewer := `{"type":"Annotation","id":"f1","briefing":"briefing:test","by":"software:roborev","at":"2026-07-20T01:00:00Z","body":"finding"}` + "\n" +
		`{"type":"Resolution","id":"r1","briefing":"briefing:test","by":"software:roborev","at":"2026-07-20T01:01:00Z","decision":"revise","includes":["f1"]}` + "\n"
	secondReviewer := reviewer +
		`{"type":"Annotation","id":"d2","briefing":"briefing:test","by":"software:roborev-2","at":"2026-07-20T01:02:00Z","includes":["f1"],"body":"class: correct-but-disproportionate; why-not-material: none; promotes-when: supported"}` + "\n" +
		`{"type":"Resolution","id":"r2","briefing":"briefing:test","by":"software:roborev-2","at":"2026-07-20T01:03:00Z","decision":"revise","includes":["d2"]}` + "\n"
	log, err := parseReviewLog([]byte(secondReviewer), briefingID)
	if err != nil {
		t.Fatal(err)
	}
	if class, err := classifyCompletedRound(log); err == nil {
		t.Fatalf("second reviewer classified as completed worker triage: %q", class)
	}

	worker := strings.ReplaceAll(secondReviewer, "software:roborev-2", "actor:ensign")
	log, err = parseReviewLog([]byte(worker), briefingID)
	if err != nil {
		t.Fatal(err)
	}
	class, err := classifyCompletedRound(log)
	if err != nil || class != "all-declines" {
		t.Fatalf("authorized worker graph = class=%q err=%v", class, err)
	}
	for _, broken := range []string{
		strings.Replace(worker, `"includes":["f1"]`, `"includes":[]`, 1),
		strings.Replace(worker, `"by":"actor:ensign"`, `"by":"software:roborev-2"`, 1),
	} {
		log, err := parseReviewLog([]byte(broken), briefingID)
		if err == nil {
			_, err = classifyCompletedRound(log)
		}
		if err == nil {
			t.Fatal("broken worker disposition graph completed triage")
		}
	}

	multiple := worker + `{"type":"Resolution","id":"r3","briefing":"briefing:test","by":"actor:ensign","at":"2026-07-20T01:04:00Z","decision":"revise","includes":["d2"]}` + "\n"
	log, err = parseReviewLog([]byte(multiple), briefingID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := classifyCompletedRound(log); err == nil {
		t.Fatal("multiple authorized worker triage Resolutions succeeded")
	}

	full := string(mustReadBytes(t, filepath.Join("testdata", "advisory-round", "briefing.review.jsonl")))
	allFixed := strings.Replace(full,
		"class: correct-but-disproportionate; why-not-material: released workflow and ACs remain correct at candidate 90aea55; promotes-when: a supported duplicate-member flow produces observable incorrect state",
		"class: material; disposition: fixed", 1)
	log, err = parseReviewLog([]byte(allFixed), "briefing:3j:implementation:round-1")
	if class, classifyErr := classifyCompletedRound(log); err != nil || classifyErr != nil || class != "all-fixed" {
		t.Fatalf("material-only triage class=%q parse=%v classify=%v", class, err, classifyErr)
	}
	lines := strings.Split(full, "\n")
	lines[3] = `{"type":"Annotation","id":"annotation:fixed","briefing":"briefing:3j:implementation:round-1","by":"actor:ensign","at":"2026-07-20T01:03:00Z","includes":["annotation:job-592"],"body":"class: material; disposition: fixed"}` + "\n" +
		`{"type":"Annotation","id":"annotation:declined","briefing":"briefing:3j:implementation:round-1","by":"actor:ensign","at":"2026-07-20T01:03:30Z","includes":["annotation:job-594"],"body":"class: correct-but-disproportionate; why-not-material: no released harm; promotes-when: supported harm"}`
	lines[4] = strings.Replace(lines[4], `["annotation:decline-duplicate-member"]`, `["annotation:fixed","annotation:declined"]`, 1)
	log, err = parseReviewLog([]byte(strings.Join(lines, "\n")), "briefing:3j:implementation:round-1")
	if class, classifyErr := classifyCompletedRound(log); err != nil || classifyErr != nil || class != "mixed" {
		t.Fatalf("mixed triage class=%q parse=%v classify=%v", class, err, classifyErr)
	}
}

func TestRoundValidateRequiresCanonicalRegularFileRoom(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(string, string)
	}{
		{"extra entry", func(room, _ string) {
			if err := os.WriteFile(filepath.Join(room, "extra"), []byte("unexpected"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{"symlinked log", func(room, log string) {
			if err := os.Remove(filepath.Join(room, "briefing.review.jsonl")); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(log, filepath.Join(room, "briefing.review.jsonl")); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, entity, briefing, log, feedback := advisoryRoundFixture(t)
			input := RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback}
			if err := RecordSemantic(entity, input); err != nil {
				t.Fatal(err)
			}
			tc.mutate(filepath.Join(root, "review", "implementation", "round-1"), log)
			if _, err := ValidateRoundFile(entity, input.Round); err == nil {
				t.Fatal("non-canonical retained room validated")
			}
		})
	}
}

func TestFeedbackCycleSpliceIsSectionScoped(t *testing.T) {
	line := "- Cycle 1: REJECTED — reviewer; surface 1/1 vs estimate 1 (100%); AC unchanged"
	entity := []byte("---\nstatus: implementation\n---\n# Entity\n\n" + line + "\n\n### Feedback Cycles\n\n### Next\n")
	out, err := spliceFeedbackCycle(entity, line, 1, true)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Count(out, []byte(line)) != 2 {
		t.Fatalf("projection did not distinguish same text outside section:\n%s", out)
	}
	feedback := bytes.Index(out, []byte("### Feedback Cycles"))
	next := bytes.Index(out, []byte("### Next"))
	projected := bytes.LastIndex(out, []byte(line))
	if projected < feedback || projected > next {
		t.Fatalf("projection is outside Feedback Cycles section:\n%s", out)
	}
}

func TestRoundCompleteOperationCASAndRollbackAreByteClean(t *testing.T) {
	t.Run("new room rolls back when entity replacement fails", func(t *testing.T) {
		root, entity, briefing, log, feedback := advisoryRoundFixture(t)
		before := treeDigest(t, root)
		injected := errors.New("injected entity failure")
		err := recordRoundLockedWith(entity, RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback}, nil,
			func(string, []byte) error { return injected })
		if !errors.Is(err, injected) {
			t.Fatalf("record error = %v, want injected failure", err)
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("complete operation left room or entity bytes after rollback")
		}
	})

	t.Run("exact room does not repair a missing pointer", func(t *testing.T) {
		root, entity, briefing, log, feedback := advisoryRoundFixture(t)
		input := RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback}
		if err := RecordSemantic(entity, input); err != nil {
			t.Fatal(err)
		}
		body := mustReadBytes(t, entity)
		withoutPointer, err := replaceTopLevels(body, false, topLevelReplacement{key: "review-round"})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(entity, withoutPointer, 0o644); err != nil {
			t.Fatal(err)
		}
		before := treeDigest(t, root)
		if err := RecordSemantic(entity, input); err == nil {
			t.Fatal("exact room repaired a missing entity pointer")
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("missing-pointer replay changed fixture bytes")
		}
	})

	t.Run("entity race after new-room publication rolls the room back", func(t *testing.T) {
		root, entity, briefing, log, feedback := advisoryRoundFixture(t)
		input := RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, FeedbackCyclePath: feedback}
		entityBefore := mustReadBytes(t, entity)
		var racedDigest string
		err := recordRoundLockedWith(entity, input, func(room string) {
			if err := os.WriteFile(entity, append(entityBefore, []byte("\nexternal race\n")...), 0o644); err != nil {
				t.Fatal(err)
			}
			racedDigest = treeDigest(t, root)
		}, atomicWrite)
		if err == nil || !strings.Contains(err.Error(), "entity changed") {
			t.Fatalf("stale entity error = %v", err)
		}
		if got := treeDigest(t, root); got != racedDigest {
			t.Fatal("stale-entity refusal changed bytes beyond the injected race")
		}
		room := filepath.Join(root, "review", "implementation", "round-1")
		if _, statErr := os.Stat(room); !os.IsNotExist(statErr) {
			t.Fatalf("stale-entity refusal left the new room: %v", statErr)
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
		"id: task\n" +
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
