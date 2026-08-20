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

func TestRoundRecordNeutralReplayAndRefusalsAreByteClean(t *testing.T) {
	root, entity, briefing, log, _ := advisoryRoundFixture(t)
	fullLog := mustReadBytes(t, log)
	input := RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log}

	// A truncated JSONL line is rejected before publication and leaves all
	// mutable state untouched.
	if err := os.WriteFile(log, fullLog[:len(fullLog)-1], 0o644); err != nil {
		t.Fatal(err)
	}
	beforeCandidate := mustReadBytes(t, filepath.Join(root, "candidate.patch"))
	beforeProduct := mustReadBytes(t, filepath.Join(root, "product", "status.txt"))
	beforeLifecycle := lifecycleBytes(t, entity)
	pendingBefore := treeDigest(t, root)
	if err := RecordSemantic(entity, input); err == nil {
		t.Fatal("incomplete review log was persisted")
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
	if summary.ID != "round:task:implementation:1" || summary.Stage != "implementation" ||
		summary.Cycle != 1 || summary.Briefing != "briefing:3j:implementation:round-1" || len(summary.Entries) != 5 {
		t.Fatalf("complete summary = %#v", summary)
	}
	resolutions := 0
	for _, entry := range summary.Entries {
		if entry.Type == "Resolution" {
			resolutions++
			if !entry.Advisory {
				t.Fatalf("round Resolution %s was not marked structurally", entry.ID)
			}
		}
	}
	if resolutions != 2 {
		t.Fatalf("resolution count = %d, want reviewer and worker", resolutions)
	}
	entityBytes := mustReadBytes(t, entity)
	if bytes.Contains(entityBytes, []byte("Feedback Cycles")) || bytes.Contains(entityBytes, []byte("Cycle 1:")) {
		t.Fatalf("neutral recorder projected workflow prose:\n%s", entityBytes)
	}
	if !bytes.Equal(mustReadBytes(t, filepath.Join(root, "candidate.patch")), beforeCandidate) ||
		!bytes.Equal(mustReadBytes(t, filepath.Join(root, "product", "status.txt")), beforeProduct) {
		t.Fatal("round recording changed candidate or product bytes")
	}
	if got := lifecycleBytes(t, entity); !bytes.Equal(got, beforeLifecycle) {
		t.Fatalf("round recording changed status or gates:\n%s", got)
	}
	if got := mustReadBytes(t, filepath.Join(room, "briefing.review.jsonl")); !bytes.Equal(got, fullLog) {
		t.Fatal("retained complete log does not equal supplied log")
	}
	if got := mustReadBytes(t, filepath.Join(room, "briefing.json")); !bytes.Equal(got, mustReadBytes(t, briefing)) {
		t.Fatal("retained Briefing bytes changed")
	}
	if got := mustReadBytes(t, filepath.Join(root, "candidate.patch")); !bytes.Equal(got, mustReadBytes(t, filepath.Join("testdata", "advisory-round", "candidate.patch"))) {
		t.Fatal("round recording changed candidate bytes")
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

func TestRoundAcceptsWorkflowForeignLabelsAndActors(t *testing.T) {
	_, entity, briefing, log, _ := advisoryRoundFixture(t)
	body := mustReadBytes(t, log)
	body = bytes.ReplaceAll(body, []byte("actor:ensign"), []byte("reviewer:outside-workflow"))
	body = bytes.ReplaceAll(body, []byte("class: correct-but-disproportionate; why-not-material: released workflow and ACs remain correct at candidate 90aea55; promotes-when: a supported duplicate-member flow produces observable incorrect state"), []byte("label: foreign-policy; disposition: deferred"))
	body = bytes.ReplaceAll(body, []byte("triage: 0 material fixed; 2 declined"), []byte("reviewer notes: arbitrary labels are opaque"))
	if err := os.WriteFile(log, body, 0o644); err != nil {
		t.Fatal(err)
	}
	input := RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log}
	if err := RecordSemantic(entity, input); err != nil {
		t.Fatalf("workflow-foreign review log rejected: %v", err)
	}
	if _, err := ValidateRoundFile(entity, input.Round); err != nil {
		t.Fatalf("workflow-foreign review log did not validate: %v", err)
	}
	if bytes.Contains(mustReadBytes(t, entity), []byte("Feedback Cycles")) {
		t.Fatal("neutral recorder wrote a workflow projection")
	}
}

func TestRoundRequiresFolderFormWithoutCrossEntityCollision(t *testing.T) {
	workflow := t.TempDir()
	for _, slug := range []string{"task-a", "task-b"} {
		entity := filepath.Join(workflow, slug+".md")
		if err := os.WriteFile(entity, []byte("---\nid: "+slug+"\nstatus: implementation\n---\n# Task\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(workflow, "unrelated"), []byte("preserve"), 0o644); err != nil {
		t.Fatal(err)
	}
	inputRoot := filepath.Join(workflow, "inputs")
	if err := os.MkdirAll(inputRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	briefing := filepath.Join(inputRoot, "briefing.json")
	log := filepath.Join(inputRoot, "briefing.review.jsonl")
	copyRoundFixture(t, briefing, "briefing.json")
	copyRoundFixture(t, log, "briefing.review.jsonl")
	for _, slug := range []string{"task-a", "task-b"} {
		before := treeDigest(t, workflow)
		entity := filepath.Join(workflow, slug+".md")
		err := RecordSemantic(entity, inputForRound(briefing, log))
		if err == nil || !strings.Contains(err.Error(), "folder-form entity") || treeDigest(t, workflow) != before {
			t.Fatalf("%s flat refusal error=%v or changed workflow bytes", slug, err)
		}
		if _, statErr := os.Stat(entity + ".gates.lock"); !os.IsNotExist(statErr) {
			t.Fatalf("%s flat refusal left a lock: %v", slug, statErr)
		}
	}
	if _, err := os.Stat(filepath.Join(workflow, "review")); !os.IsNotExist(err) {
		t.Fatalf("flat refusals created a shared review room: %v", err)
	}
}

func TestRoundNoFindingsAndPreflightRefusals(t *testing.T) {
	// This case is also the closure predicate's over-refusal control: a lone
	// `approve` Resolution IS the whole round, so refusing it would break a shape
	// the recorder has always accepted.
	t.Run("no findings is structurally valid and has no projection", func(t *testing.T) {
		_, entity, briefing, log, _ := advisoryRoundFixture(t)
		noFindings := `{"type":"Resolution","id":"resolution:roborev-clear","briefing":"briefing:3j:implementation:round-1","by":"reviewer:foreign","at":"2026-07-20T01:00:00Z","decision":"approve"}` + "\n"
		if err := os.WriteFile(log, []byte(noFindings), 0o644); err != nil {
			t.Fatal(err)
		}
		input := inputForRound(briefing, log)
		if err := RecordSemantic(entity, input); err != nil {
			t.Fatal(err)
		}
		summary, err := ValidateRoundFile(entity, input.Round)
		if err != nil || len(summary.Entries) != 1 {
			t.Fatalf("summary=%#v err=%v", summary, err)
		}
		if bytes.Contains(mustReadBytes(t, entity), []byte("Cycle 1:")) {
			t.Fatal("neutral no-findings round received a workflow projection")
		}
	})

	for _, tc := range []struct {
		name   string
		mutate func(root, entity, briefing, log string)
		want   string
	}{
		{"bad artifact digest", func(root, entity, briefing, log string) {
			if err := os.WriteFile(filepath.Join(root, "candidate.patch"), []byte("corrupt\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}, "artifact"},
		{"self-referential mutable entity artifact", func(root, entity, briefing, log string) {
			body := bytes.Replace(mustReadBytes(t, briefing), []byte("../../../candidate.patch"), []byte("../../../index.md"), 1)
			body = bytes.Replace(body, []byte("sha256:8e85d4c9523a617e05b17c92390b10b2f9892152ca348433311230ac3ad98dd3"), []byte(RawDigest(mustReadBytes(t, entity))), 1)
			if err := os.WriteFile(briefing, body, 0o644); err != nil {
				t.Fatal(err)
			}
		}, "mutable entity"},
		{"occupied target", func(root, entity, briefing, log string) {
			room := filepath.Join(root, "review", "implementation", "round-1")
			if err := os.MkdirAll(room, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(room, "foreign"), []byte("occupied"), 0o644); err != nil {
				t.Fatal(err)
			}
		}, "canonical room"},
		{"lock contention", func(root, entity, briefing, log string) {
			if err := os.WriteFile(entity+".gates.lock", []byte("held"), 0o600); err != nil {
				t.Fatal(err)
			}
		}, "lock"},
		{"cross briefing log", func(root, entity, briefing, log string) {
			body := bytes.Replace(mustReadBytes(t, log), []byte(`"briefing":"briefing:3j:implementation:round-1"`), []byte(`"briefing":"briefing:other"`), 1)
			if err := os.WriteFile(log, body, 0o644); err != nil {
				t.Fatal(err)
			}
		}, "same Briefing"},
		{"malformed artifact URI", func(root, entity, briefing, log string) {
			body := bytes.Replace(mustReadBytes(t, briefing), []byte("../../../candidate.patch"), []byte("%gh"), 1)
			if err := os.WriteFile(briefing, body, 0o644); err != nil {
				t.Fatal(err)
			}
		}, "URI"},
		// The three open-log shapes. Each `want` pins the tail the refusal names,
		// not just that it refused, so a message that says "not closed" without
		// saying what it saw fails here. Deleting the predicate records the
		// truncated room and reds all three.
		{"log ends at the reviewer's revise verdict", func(root, entity, briefing, log string) {
			writeRoundLog(t, log, roundLogPrefix(t, log, 3))
		}, "ends at the reviewer's revise Resolution resolution:roborev-3j-round-1"},
		{"log ends at a dangling disposition annotation", func(root, entity, briefing, log string) {
			writeRoundLog(t, log, roundLogPrefix(t, log, 4))
		}, "ends at Annotation annotation:decline-duplicate-member with no closing Resolution"},
		{"log ends at a hold verdict", func(root, entity, briefing, log string) {
			writeRoundLog(t, log, bytes.Replace(roundLogPrefix(t, log, 3),
				[]byte(`"decision":"revise"`), []byte(`"decision":"hold"`), 1))
		}, "ends at the reviewer's hold Resolution resolution:roborev-3j-round-1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, entity, briefing, log, _ := advisoryRoundFixture(t)
			tc.mutate(root, entity, briefing, log)
			before := treeDigest(t, root)
			err := RecordSemantic(entity, inputForRound(briefing, log))
			if err == nil || !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tc.want)) {
				t.Fatalf("error=%v, want %q", err, tc.want)
			}
			if got := treeDigest(t, root); got != before {
				t.Fatal("refusal changed fixture tree")
			}
		})
	}

	// The refusal is only useful if the state it leaves is still recordable: the FO
	// takes the recovery its message names — route the correction, append that
	// round's entries — and records the SAME round successfully. Without this the
	// precondition could be satisfied by a dead end.
	t.Run("open-log refusal leaves a recordable round", func(t *testing.T) {
		root, entity, briefing, log, _ := advisoryRoundFixture(t)
		complete := mustReadBytes(t, log)
		writeRoundLog(t, log, roundLogPrefix(t, log, 3))
		before := treeDigest(t, root)
		err := RecordSemantic(entity, inputForRound(briefing, log))
		if err == nil || !strings.Contains(err.Error(), "not closed") {
			t.Fatalf("open-log refusal error = %v, want the stable \"not closed\" token", err)
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("open-log refusal changed fixture bytes")
		}
		writeRoundLog(t, log, complete)
		if err := RecordSemantic(entity, inputForRound(briefing, log)); err != nil {
			t.Fatalf("record after the recovery the refusal names: %v", err)
		}
		summary, err := ValidateRoundFile(entity, "implementation/1")
		if err != nil || len(summary.Entries) != 5 {
			t.Fatalf("recovered round summary = %#v err = %v", summary, err)
		}
	})

	t.Run("stage taxonomy permits historical backfill", func(t *testing.T) {
		root, entity, briefing, log, _ := advisoryRoundFixture(t)
		body := bytes.Replace(mustReadBytes(t, entity), []byte("status: implementation"), []byte("status: validation"), 1)
		if err := os.WriteFile(entity, body, 0o644); err != nil {
			t.Fatal(err)
		}
		input := RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log, WorkflowDir: root}
		if err := RecordSemantic(entity, input); err != nil {
			t.Fatalf("historical implementation round at validation status: %v", err)
		}
	})

	t.Run("unknown workflow stage", func(t *testing.T) {
		root, entity, briefing, log, _ := advisoryRoundFixture(t)
		before := treeDigest(t, root)
		if err := RecordSemantic(entity, RecordInput{Round: "unknown/1", BriefingPath: briefing, LogPath: log, WorkflowDir: root}); err == nil {
			t.Fatal("round for undefined workflow stage recorded")
		}
		if treeDigest(t, root) != before {
			t.Fatal("undefined-stage refusal changed fixture tree")
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
		if err := mutateEntity(entity, entityExpectation{}, build, atomicWrite); err == nil || buildCalled {
			t.Fatal("zero entity expectation succeeded or built")
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
			root, entity, briefing, log, _ := advisoryRoundFixture(t)
			input := inputForRound(briefing, log)
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

func TestRoundCompleteOperationCASAndRollbackAreByteClean(t *testing.T) {
	t.Run("new room rolls back when entity replacement fails", func(t *testing.T) {
		root, entity, briefing, log, _ := advisoryRoundFixture(t)
		before := treeDigest(t, root)
		injected := errors.New("injected entity failure")
		err := recordRoundLockedWith(entity, inputForRound(briefing, log), nil, func(string, []byte) error { return injected })
		if !errors.Is(err, injected) {
			t.Fatalf("record error = %v, want injected failure", err)
		}
		if got := treeDigest(t, root); got != before {
			t.Fatal("complete operation left room or entity bytes after rollback")
		}
	})

	t.Run("exact room does not repair a missing pointer", func(t *testing.T) {
		root, entity, briefing, log, _ := advisoryRoundFixture(t)
		input := inputForRound(briefing, log)
		if err := RecordSemantic(entity, input); err != nil {
			t.Fatal(err)
		}
		body := mustReadBytes(t, entity)
		withoutPointer, err := replaceTopLevels(body, topLevelReplacement{key: "review-round"})
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
		root, entity, briefing, log, _ := advisoryRoundFixture(t)
		input := inputForRound(briefing, log)
		entityBefore := mustReadBytes(t, entity)
		var racedDigest string
		err := recordRoundLockedWith(entity, input, func(string) {
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

func inputForRound(briefing, log string) RecordInput {
	return RecordInput{Round: "implementation/1", BriefingPath: briefing, LogPath: log}
}

func advisoryRoundFixture(t *testing.T) (root, entity, briefing, log, feedback string) {
	t.Helper()
	return advisoryRoundFixtureAt(t, filepath.Join(t.TempDir(), "task"))
}

func advisoryRoundFixtureAt(t *testing.T, root string) (string, string, string, string, string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "inputs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "product"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("---\nstages:\n  states:\n    - name: implementation\n    - name: validation\n---\n# Workflow\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	copyRoundFixture(t, filepath.Join(root, "candidate.patch"), "candidate.patch")
	briefing := filepath.Join(root, "inputs", "briefing.json")
	log := filepath.Join(root, "inputs", "briefing.review.jsonl")
	copyRoundFixture(t, briefing, "briefing.json")
	copyRoundFixture(t, log, "briefing.review.jsonl")
	feedback := filepath.Join(root, "inputs", "feedback-cycle.txt")
	if err := os.WriteFile(feedback, []byte("- Cycle 1: REJECTED — Roborev; surface 0/0 vs estimate 340 (0%); AC unchanged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "product", "status.txt"), []byte("candidate=90aea55\nstate=unchanged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(root, "index.md")
	body := "---\n" +
		"id: task\n" +
		"status: implementation\n" +
		"custom: preserve-me\n" +
		"gates:\n" +
		"  version: 1\n" +
		"  records:\n" +
		"    - id: gate:task:ideation\n" +
		"      stage: ideation\n" +
		"      attempts:\n" +
		"        - id: gate-attempt:task-ideation-1\n" +
		"          briefing:\n" +
		"            id: briefing:task:ideation:1\n" +
		"            digest: sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n" +
		"            room-ref: ./review/ideation/briefing-1\n" +
		"title: Task\n" +
		"---\n" +
		"# Task\n\nUnrelated body bytes stay fixed.\n"
	if err := os.WriteFile(entity, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, entity, briefing, log, feedback
}

// roundLogPrefix returns the first `lines` complete JSONL lines of a review log,
// which is how the open-log shapes are built: the fixture's five-entry log
// truncated at the reviewer's verdict (3) or at the ensign's dangling disposition
// annotation (4).
func roundLogPrefix(t *testing.T, log string, lines int) []byte {
	t.Helper()
	all := bytes.SplitAfter(mustReadBytes(t, log), []byte{'\n'})
	if len(all) <= lines {
		t.Fatalf("review log has %d lines, want at least %d", len(all)-1, lines)
	}
	return bytes.Join(all[:lines], nil)
}

func writeRoundLog(t *testing.T, log string, body []byte) {
	t.Helper()
	if err := os.WriteFile(log, body, 0o644); err != nil {
		t.Fatal(err)
	}
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
