package ensigncycle

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

const rejectionPreparedBriefingID = "briefing:rejection-task:validation:attempt-1:revision-1"

func assertRejectionRoundGateBoundary(entityPath, wantStatus string) error {
	doc, _, err := gates.Read(entityPath)
	if err != nil && strings.Contains(err.Error(), "entity has no gates record") {
		return nil
	}
	if err != nil {
		return fmt.Errorf("malformed final validation gate: %w", err)
	}
	if wantStatus != "validation" {
		return fmt.Errorf("round-only state contains an ordinary gate record; `gate record --round` must retain only advisory review-round state")
	}
	if len(doc.Records) != 1 || doc.Current.Gate != "gate:rejection-task:validation" {
		return fmt.Errorf("final gate selection does not identify exactly one rejection-task validation gate")
	}
	record := doc.Records[0]
	if record.ID != doc.Current.Gate || record.Stage != "validation" || len(record.Attempts) != 1 {
		return fmt.Errorf("selected validation gate does not contain exactly one attempt")
	}
	attempt := record.Attempts[0]
	if attempt.Briefing.ID == rejectionBriefingID {
		return fmt.Errorf("validation/1 advisory round was retained as a gate attempt")
	}
	if attempt.ID != "gate-attempt:rejection-task-validation-1" ||
		attempt.Briefing.ID != rejectionPreparedBriefingID {
		return fmt.Errorf("selected validation gate is not bound to the expected prepared Briefing")
	}
	if attempt.Resolution != nil || attempt.Application != nil {
		return fmt.Errorf("final round-2 validation gate is not open")
	}
	return nil
}

func assertRejectionRecordedRound(workflowRoot, entityPath, wantStatus string) error {
	entity, err := os.ReadFile(entityPath)
	if err != nil {
		return err
	}
	for _, preserved := range []string{
		"workflow-state: preserve-me",
		"gate-state: preserve-me",
		"application-state: preserve-me",
	} {
		if !bytes.Contains(entity, []byte(preserved)) {
			return fmt.Errorf("round recording corrupted lifecycle sentinel %q", preserved)
		}
	}
	if !regexp.MustCompile(`(?m)^status: ` + regexp.QuoteMeta(wantStatus) + `$`).Match(entity) {
		return fmt.Errorf("entity status changed or did not reach %s", wantStatus)
	}
	if err := assertRejectionRoundGateBoundary(entityPath, wantStatus); err != nil {
		return err
	}
	if got := bytes.Count(entity, []byte(rejectionFeedbackCycle)); got != 1 {
		return fmt.Errorf("Cycle 1 projection count = %d, want exactly 1", got)
	}

	summary, err := gates.ValidateRoundFile(entityPath, "validation/1")
	if err != nil {
		return fmt.Errorf("validate retained round: %w", err)
	}
	if summary.ID != "round:rejection-task:validation:1" ||
		summary.Stage != "validation" ||
		summary.Cycle != 1 ||
		summary.Briefing != rejectionBriefingID ||
		summary.Triage != "all-fixed" ||
		len(summary.Entries) != 4 {
		return fmt.Errorf("retained round summary = %#v", summary)
	}
	for _, entry := range summary.Entries {
		if entry.Type == "Resolution" && !entry.Advisory {
			return fmt.Errorf("retained Resolution %s is not advisory", entry.ID)
		}
	}

	room := filepath.Join(filepath.Dir(entityPath), "review", "validation", "round-1")
	entries, err := os.ReadDir(room)
	if err != nil {
		return err
	}
	if len(entries) != 2 ||
		entries[0].Name() != "briefing.json" ||
		entries[1].Name() != "briefing.review.jsonl" ||
		!entries[0].Type().IsRegular() ||
		!entries[1].Type().IsRegular() {
		return fmt.Errorf("round room entries = %#v, want exactly two regular canonical files", entries)
	}
	checks := []struct {
		path string
		want string
	}{
		{filepath.Join(room, "briefing.json"), rejectionBriefing()},
		{filepath.Join(room, "briefing.review.jsonl"), rejectionCompleteLog()},
		{filepath.Join(filepath.Dir(entityPath), "candidate.txt"), rejectionCandidate},
		{filepath.Join(workflowRoot, "README.md"), rejectionReadme()},
	}
	for _, check := range checks {
		got, readErr := os.ReadFile(check.path)
		if readErr != nil {
			return readErr
		}
		if !bytes.Equal(got, []byte(check.want)) {
			return fmt.Errorf("%s changed from its exact expected bytes", check.path)
		}
	}
	return nil
}

func TestRejectionFlowRoundRecordingDurableOracle(t *testing.T) {
	root := t.TempDir()
	entityPath := writeRejectionWorkflow(t, root)
	writeFile(t, filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"), rejectionCompleteLog())
	before := readFile(t, entityPath)
	if err := gates.RecordSemantic(entityPath, gates.RecordInput{
		Round:             "validation/1",
		BriefingPath:      filepath.Join(root, "rejection-task", "inputs", "briefing.json"),
		LogPath:           filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"),
		FeedbackCyclePath: filepath.Join(root, "rejection-task", "inputs", "feedback-cycle.txt"),
		WorkflowDir:       root,
	}); err != nil {
		t.Fatalf("record completed triage: %v", err)
	}
	after := readFile(t, entityPath)
	for _, line := range []string{
		"status: backlog",
		"workflow-state: preserve-me",
		"gate-state: preserve-me",
		"application-state: preserve-me",
	} {
		if strings.Count(before, line) != 1 || strings.Count(after, line) != 1 {
			t.Fatalf("round recorder did not preserve exact lifecycle line %q", line)
		}
	}
	if err := assertRejectionRecordedRound(root, entityPath, "backlog"); err != nil {
		t.Fatalf("recorded-round durable oracle: %v", err)
	}

	writeFile(t, entityPath, strings.Replace(readFile(t, entityPath), "status: backlog", "status: validation", 1))
	if err := assertRejectionRecordedRound(root, entityPath, "validation"); err != nil {
		t.Fatalf("recorded-round oracle rejected valid final validation state without a gate: %v", err)
	}
	gateReviewPath := filepath.Join(root, "rejection-task", "inputs", "gate-validation", "gate-review.md")
	writeFile(t, gateReviewPath, "# Rejection Task — validation review\n\nThe corrected candidate is ready for its prepared decision gate.\n")
	gateReviewRel, err := filepath.Rel(root, gateReviewPath)
	if err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", "--", gateReviewRel)
	git(t, root, "commit", "-q", "-m", "add prepared validation review", "--", gateReviewRel)
	prepared, err := gates.Prepare(entityPath, gates.PrepareInput{
		WorkflowDir: root,
		Question:    "Does the second validation confirm the fix marker is present and PASS?",
		Artifact:    gateReviewPath,
		Summary:     "The corrected rejection-flow candidate is ready for a decision.",
	})
	if err != nil {
		t.Fatalf("prepare later validation gate: %v", err)
	}
	if prepared.Briefing != rejectionPreparedBriefingID || prepared.State != "open" {
		t.Fatalf("prepared later validation gate=%#v", prepared)
	}
	if err := assertRejectionRecordedRound(root, entityPath, "validation"); err != nil {
		t.Fatalf("recorded-round oracle rejected later open round-2 validation gate: %v", err)
	}
	openGateEntity := readFile(t, entityPath)
	for _, control := range []struct{ entity, want string }{
		{strings.Replace(openGateEntity, "              briefing:\n", "              state: open\n              briefing:\n", 1), "malformed final validation gate"},
		{strings.Replace(openGateEntity, rejectionPreparedBriefingID, rejectionBriefingID, 1), "validation/1 advisory round was retained as a gate attempt"},
		{strings.ReplaceAll(openGateEntity, "gate:rejection-task:validation", "gate:rejection-task:wrong"), "final gate selection does not identify"},
	} {
		writeFile(t, entityPath, control.entity)
		if err := assertRejectionRecordedRound(root, entityPath, "validation"); err == nil ||
			!strings.Contains(err.Error(), control.want) {
			t.Fatalf("gate control diagnostic = %v, want %q", err, control.want)
		}
	}
	writeFile(t, entityPath, openGateEntity)
	if err := gates.RecordSemantic(entityPath, gates.RecordInput{
		Decision: "hold", Actor: "person:captain", Reason: "exercise closed-gate counterexample", WorkflowDir: root,
	}); err != nil {
		t.Fatalf("close later validation gate for counterexample: %v", err)
	}
	if err := assertRejectionRecordedRound(root, entityPath, "validation"); err == nil ||
		!strings.Contains(err.Error(), "final round-2 validation gate is not open") {
		t.Fatalf("closed gate control diagnostic = %v", err)
	}
}
