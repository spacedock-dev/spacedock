package ensigncycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

var directRoundLauncher = regexp.MustCompile(`(?:^|[\s;&|])['"]*(?:spacedock|\$(?:\{SPACEDOCK_BIN(?::-[^}]*)?\}|SPACEDOCK_BIN)|/[^ \t\r\n'";&|]+/spacedock)['"]*\s+gate\s+record(?:\s|$)`)
var rejectionRoundSuccess = regexp.MustCompile(`(?m)^round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4$`)

const rejectionPreparedBriefingID = "briefing:rejection-task:validation:attempt-1:revision-1"

func rejectionRoundArtifactArg(flag, filename string) *regexp.Regexp {
	return regexp.MustCompile(`--` + regexp.QuoteMeta(flag) + `(?:=|\s+)['"]?(?:[^ \t\r\n'";&|]*/)?rejection-task/inputs/` +
		regexp.QuoteMeta(filename) + `['"]?(?:\s|[;&|]|$)`)
}

func commandRecordsRejectionRound(command string) bool {
	command = strings.ReplaceAll(command, "\\\n", " ")
	for _, required := range []*regexp.Regexp{
		regexp.MustCompile(`(?:^|\s)['"]?rejection-task['"]?(?:\s|$)`),
		regexp.MustCompile(`--round(?:=|\s+)['"]?validation/1['"]?(?:\s|$)`),
		rejectionRoundArtifactArg("briefing", "briefing.json"),
		rejectionRoundArtifactArg("log", "briefing.review.jsonl"),
	} {
		if !required.MatchString(command) {
			return false
		}
	}
	if directRoundLauncher.MatchString(command) {
		return true
	}
	match := launcherCapture.FindStringSubmatch(command)
	if match == nil {
		return false
	}
	varName := ""
	for _, candidate := range match[1:] {
		if candidate != "" {
			varName = candidate
			break
		}
	}
	if varName == "" {
		return false
	}
	captureEnd := strings.Index(command, match[0]) + len(match[0])
	segments := regexp.MustCompile(`\r?\n|;|&&|\|\||\|`).Split(command[captureEnd:], -1)
	call := regexp.MustCompile(`^(?:\$` + regexp.QuoteMeta(varName) + `|\$\{` + regexp.QuoteMeta(varName) + `\}|"\$` + regexp.QuoteMeta(varName) + `"|"\$\{` + regexp.QuoteMeta(varName) + `\}")\s+gate\s+record(?:\s|$)`)
	for _, segment := range segments {
		if call.MatchString(strings.TrimSpace(segment)) {
			return true
		}
	}
	return false
}

func successfulRejectionRoundResult(content json.RawMessage, isError *bool) bool {
	var text string
	return isError != nil && !*isError && json.Unmarshal(content, &text) == nil &&
		rejectionRoundSuccess.MatchString(text)
}

func claudeRecordedRejectionRound(stream string) bool {
	invocations := map[string]bool{}
	for _, line := range strings.Split(stream, "\n") {
		var entry struct {
			Message *struct {
				Content []struct {
					Type      string `json:"type"`
					Name      string `json:"name"`
					ID        string `json:"id"`
					ToolUseID string `json:"tool_use_id"`
					Input     struct {
						Command string `json:"command"`
					} `json:"input"`
					Content json.RawMessage `json:"content"`
					IsError *bool           `json:"is_error"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &entry) != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type == "tool_use" && block.Name == "Bash" && block.ID != "" &&
				commandRecordsRejectionRound(block.Input.Command) {
				invocations[block.ID] = true
			}
			if block.Type == "tool_result" && invocations[block.ToolUseID] &&
				successfulRejectionRoundResult(block.Content, block.IsError) {
				return true
			}
		}
	}
	return false
}

func codexRecordedRejectionRound(jsonl string) bool {
	for _, line := range strings.Split(jsonl, "\n") {
		var entry struct {
			Type string `json:"type"`
			Item struct {
				Type     string `json:"type"`
				Command  string `json:"command"`
				ExitCode *int   `json:"exit_code"`
			} `json:"item"`
		}
		if json.Unmarshal([]byte(line), &entry) == nil &&
			entry.Type == "item.completed" &&
			entry.Item.Type == "command_execution" &&
			entry.Item.ExitCode != nil &&
			*entry.Item.ExitCode == 0 &&
			commandRecordsRejectionRound(entry.Item.Command) {
			return true
		}
	}
	return false
}

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
	if len(doc.Records) != 1 || doc.Records[0].Stage != "validation" {
		return fmt.Errorf("final gate selection does not identify exactly one rejection-task validation gate")
	}
	record := doc.Records[0]
	if record.Stage != "validation" || len(record.Attempts) != 1 {
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

func assertRejectionRecordedRound(workflowRoot, entityPath, wantStatus string, invoked bool) error {
	if !invoked {
		return fmt.Errorf("resolved launcher never invoked `gate record --round validation/1`")
	}
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
	summary, err := gates.ValidateRoundFile(entityPath, "validation/1")
	if err != nil {
		return fmt.Errorf("validate retained round: %w", err)
	}
	if summary.ID != "round:rejection-task:validation:1" ||
		summary.Stage != "validation" ||
		summary.Cycle != 1 ||
		summary.Briefing != rejectionBriefingID ||
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

func TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl(t *testing.T) {
	root := t.TempDir()
	entityPath := writeRejectionWorkflow(t, root)
	writeFile(t, filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"), rejectionCompleteLog())
	// Workflow policy is authored by the First Officer before invoking the
	// neutral producer. The recorder must retain this line byte-for-byte.
	cycleBefore := readFile(t, entityPath)
	writeFile(t, entityPath, cycleBefore+"\n### Feedback Cycles\n\n"+rejectionFeedbackCycle)
	before := readFile(t, entityPath)
	if err := gates.RecordSemantic(entityPath, gates.RecordInput{
		Round:        "validation/1",
		BriefingPath: filepath.Join(root, "rejection-task", "inputs", "briefing.json"),
		LogPath:      filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"),
		WorkflowDir:  root,
	}); err != nil {
		t.Fatalf("record completed triage: %v", err)
	}
	after := readFile(t, entityPath)
	if strings.Count(after, rejectionFeedbackCycle) != 1 {
		t.Fatalf("neutral round recorder changed workflow-owned Cycle projection")
	}
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
	if err := assertRejectionRecordedRound(root, entityPath, "backlog", true); err != nil {
		t.Fatalf("recorded-round durable oracle: %v", err)
	}
	if err := assertRejectionRecordedRound(root, entityPath, "implementation", false); err == nil {
		t.Fatal("inverted no-invocation control passed despite no observed launcher call")
	}

	writeFile(t, entityPath, strings.Replace(readFile(t, entityPath), "status: backlog", "status: validation", 1))
	if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err != nil {
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
	if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err != nil {
		t.Fatalf("recorded-round oracle rejected later open round-2 validation gate: %v", err)
	}
	openGateEntity := readFile(t, entityPath)
	for _, control := range []struct{ entity, want string }{
		{strings.Replace(openGateEntity, "              briefing:\n", "              state: open\n              briefing:\n", 1), "malformed final validation gate"},
		{strings.Replace(openGateEntity, rejectionPreparedBriefingID, rejectionBriefingID, 1), "validation/1 advisory round was retained as a gate attempt"},
		{strings.Replace(openGateEntity, "      stage: validation", "      stage: wrong", 1), "final gate selection does not identify"},
	} {
		writeFile(t, entityPath, control.entity)
		if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err == nil ||
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
	if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err == nil ||
		!strings.Contains(err.Error(), "final round-2 validation gate is not open") {
		t.Fatalf("closed gate control diagnostic = %v", err)
	}
}

func TestRejectionFlowRoundInvocationExtractors(t *testing.T) {
	command := `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --workflow-dir "$WD" --round validation/1 --briefing "$WD/rejection-task/inputs/briefing.json" --log "$WD/rejection-task/inputs/briefing.review.jsonl"`
	result := "round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4"
	claudeStream := strings.Join([]string{
		bashToolLine("toolu_round", command),
		toolResultLine("toolu_round", false, result),
	}, "\n")
	if !claudeRecordedRejectionRound(claudeStream) {
		t.Fatal("Claude extractor missed correlated round invocation with prefixed artifact paths")
	}
	for name, invalid := range map[string]string{
		"wrong_suffix": strings.Replace(command, "briefing.json\"", "briefing.json.bak\"", 1),
		"wrong_file":   strings.Replace(command, "briefing.review.jsonl", "other.review.jsonl", 1),
		"wrong_entity": strings.Replace(command, "gate record rejection-task", "gate record other-task", 1),
		"wrong_round":  strings.Replace(command, "validation/1", "validation/2", 1),
	} {
		t.Run(name, func(t *testing.T) {
			if commandRecordsRejectionRound(invalid) {
				t.Fatal("command recognizer accepted invalid rejection-round arguments")
			}
		})
	}
	if claudeRecordedRejectionRound(bashToolLine("toolu_round", command)) ||
		claudeRecordedRejectionRound(strings.Join([]string{
			bashToolLine("toolu_round", command),
			toolResultLine("toolu_round", true, result),
		}, "\n")) {
		t.Fatal("Claude extractor accepted a missing or failed correlated result")
	}
	captured := "B=${SPACEDOCK_BIN:-spacedock}\n$B gate record rejection-task --workflow-dir . --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl"
	if !codexRecordedRejectionRound(codexCommandOutput(captured, result, 0, "completed")) {
		t.Fatal("Codex extractor missed captured resolved launcher round invocation")
	}
	retainedCodexWrapped := `/bin/zsh -lc "rg -n '"'^shared-rejection-fix: applied$|''^## Stage Report: implementation|''^- DONE:'"' rejection-task/index.md; tail -n 4 rejection-task/inputs/briefing.review.jsonl; "'${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl --workflow-dir .; ${SPACEDOCK_BIN:-spacedock} state commit rejection-task'`
	if !codexRecordedRejectionRound(codexCommandOutput(retainedCodexWrapped, result, 0, "completed")) {
		t.Fatal("Codex extractor missed retained nested-shell resolved launcher round invocation")
	}
	if codexRecordedRejectionRound(codexCommandOutput(retainedCodexWrapped, result, 1, "failed")) {
		t.Fatal("Codex extractor accepted failed retained round invocation")
	}
	absoluteMultiline := `/tmp/candidate/spacedock gate record \
  "rejection-task" --workflow-dir . --round="validation/1" \
  --briefing "rejection-task/inputs/briefing.json" \
  --log="rejection-task/inputs/briefing.review.jsonl"`
	if !commandRecordsRejectionRound(absoluteMultiline) {
		t.Fatal("command recognizer missed absolute resolved launcher with multiline/quoted flags")
	}
	noInvocation := codexCommandOutput("spacedock status rejection-task --workflow-dir .", "", 0, "completed")
	if codexRecordedRejectionRound(noInvocation) {
		t.Fatal("no-invocation transcript falsely satisfied the round invocation extractor")
	}
}
