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
var rejectionRoundSuccess = regexp.MustCompile(`(?m)^round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=(?:2|4)$`)
var rejectionValidation2Command = regexp.MustCompile(`(?:spacedock|\$\{SPACEDOCK_BIN(?::-[^}]*)?\}|\$launcher|/[^ \t\r\n'";&|]+/spacedock)["']?\s+gate\s+record\s+["']?rejection-task["']?.*--round(?:=|\s+)["']?validation/2["']?`)
var rejectionPrepareCommand = regexp.MustCompile(`(?:spacedock|\$\{SPACEDOCK_BIN(?::-[^}]*)?\}|/[^ \t\r\n'";&|]+/spacedock)["']?\s+gate\s+prepare\s+["']?rejection-task["']?(?:\s|$)`)

const rejectionPreparedBriefingID = "briefing:rejection-task:validation:attempt-1:revision-1"

func rejectionRoundArtifactArg(flag, filename string) *regexp.Regexp {
	return regexp.MustCompile(`--` + regexp.QuoteMeta(flag) + `(?:=|\s+)['"]?(?:[^ \t\r\n'";&|]*/)?rejection-task/inputs/` +
		regexp.QuoteMeta(filename) + `['"]?(?:\s|[;&|]|$)`)
}
func commandRecordsRejectionValidation2(command string) bool {
	return rejectionValidation2Command.MatchString(command) && rejectionRoundArtifactArg("briefing", "briefing.json").MatchString(command) && rejectionRoundArtifactArg("log", "briefing.review.jsonl").MatchString(command)
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
	for _, command := range codexSuccessfulCommands(jsonl) {
		if commandRecordsRejectionRound(command) {
			return true
		}
	}
	return false
}
func codexSuccessfulCommands(jsonl string) (commands []string) {
	for _, line := range strings.Split(jsonl, "\n") {
		var entry struct {
			Type string `json:"type"`
			Item struct {
				Type, Command string
				ExitCode      *int `json:"exit_code"`
			} `json:"item"`
		}
		if json.Unmarshal([]byte(line), &entry) == nil && entry.Type == "item.completed" && entry.Item.Type == "command_execution" && entry.Item.ExitCode != nil && *entry.Item.ExitCode == 0 {
			commands = append(commands, entry.Item.Command)
		}
	}
	return commands
}
func codexRejectionGateSequence(jsonl string) error {
	stage, prepares := 0, 0
	for _, command := range codexSuccessfulCommands(jsonl) {
		if commandRecordsRejectionValidation2(command) {
			stage = 1
		}
		if rejectionPrepareCommand.MatchString(command) {
			prepares++
			if stage == 1 {
				stage = 2
			}
		}
	}
	if stage != 2 || prepares != 1 {
		return fmt.Errorf("Codex final-gate sequence stage/prepares = %d/%d, want 2/1", stage, prepares)
	}
	return nil
}
func assertRejectionRoundGateBoundary(entityPath, wantStatus string, required ...bool) error {
	doc, _, err := gates.Read(entityPath)
	if err != nil && strings.Contains(err.Error(), "entity has no gates record") {
		if len(required) != 0 && required[0] {
			return fmt.Errorf("required final validation gate is absent")
		}
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
	if record.Stage != "validation" || len(record.Attempts) == 0 {
		return fmt.Errorf("selected validation gate does not contain exactly one attempt")
	}
	for _, earlier := range record.Attempts[:len(record.Attempts)-1] {
		if earlier.Withdrawal == nil || earlier.ProviderEvidence != nil || earlier.Resolution != nil || earlier.Application != nil {
			return fmt.Errorf("earlier validation gate attempt is not cleanly withdrawn")
		}
	}
	attemptNumber, attempt := len(record.Attempts), record.Attempts[len(record.Attempts)-1]
	if attempt.Briefing.ID == rejectionBriefingID {
		return fmt.Errorf("validation/1 advisory round was retained as a gate attempt")
	}
	if attempt.ID != fmt.Sprintf("gate-attempt:rejection-task-validation-%d", attemptNumber) ||
		attempt.Briefing.ID != fmt.Sprintf("briefing:rejection-task:validation:attempt-%d:revision-1", attemptNumber) {
		return fmt.Errorf("selected validation gate is not bound to the expected prepared Briefing")
	}
	if attempt.Withdrawal != nil || attempt.ProviderEvidence != nil || attempt.Resolution != nil || attempt.Application != nil {
		return fmt.Errorf("final round-2 validation gate is not open")
	}
	if len(required) != 0 && required[0] {
		entity, readErr := os.ReadFile(entityPath)
		if readErr != nil || regexp.MustCompile(`(?m)^(?:completed|verdict):[^\S\n]*\S`).Match(entity) {
			return fmt.Errorf("final validation gate has terminal entity state: %v", readErr)
		}
	}
	return nil
}
func assertCodexRejectionFinalGate(entityPath, jsonl string) error {
	if err := assertRejectionRoundGateBoundary(entityPath, "validation", true); err != nil {
		return err
	}
	return codexRejectionGateSequence(jsonl)
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
	cycle := 1 + bytes.Count(entity, []byte("round:rejection-task:validation:2"))
	summary, err := gates.ValidateRoundFile(entityPath, fmt.Sprintf("validation/%d", cycle))
	if err != nil {
		return fmt.Errorf("validate retained round: %w", err)
	}
	if summary.ID != fmt.Sprintf("round:rejection-task:validation:%d", cycle) ||
		summary.Stage != "validation" ||
		summary.Cycle != cycle ||
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
	if err := assertRejectionRoundGateBoundary(entityPath, "validation", true); err == nil {
		t.Fatal("strict final-gate oracle accepted an absent gate")
	}
	writeFile(t, entityPath, readFile(t, entityPath)+"- Cycle 2: PASSED\n")
	if err := gates.RecordSemantic(entityPath, gates.RecordInput{Round: "validation/2", BriefingPath: filepath.Join(root, "rejection-task", "inputs", "briefing.json"), LogPath: filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"), WorkflowDir: root}); err != nil {
		t.Fatalf("record validation/2: %v", err)
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
		{strings.Replace(openGateEntity, "round:rejection-task:validation:2", "round:rejection-task:validation:9", 1), "validate retained round"},
	} {
		writeFile(t, entityPath, control.entity)
		if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err == nil ||
			!strings.Contains(err.Error(), control.want) {
			t.Fatalf("gate control diagnostic = %v, want %q", err, control.want)
		}
	}
	writeFile(t, entityPath, openGateEntity)
	if _, err := gates.Withdraw(entityPath, gates.WithdrawInput{Reason: "replace premature validation gate", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	if _, err := gates.Prepare(entityPath, gates.PrepareInput{WorkflowDir: root, Question: "Is the corrected candidate ready?", Artifact: gateReviewPath, Summary: "Ready for a decision."}); err != nil {
		t.Fatal(err)
	}
	recoveredGateEntity := readFile(t, entityPath)
	writeFile(t, entityPath, regexp.MustCompile(`(?m)^              withdrawal:\n(?:                .+\n){3}`).ReplaceAllString(recoveredGateEntity, ""))
	if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err == nil || !strings.Contains(err.Error(), "not cleanly withdrawn") {
		t.Fatalf("active prior-attempt control diagnostic = %v", err)
	}
	writeFile(t, entityPath, recoveredGateEntity)
	round1Briefing := filepath.Join(filepath.Dir(entityPath), "review", "validation", "round-1", "briefing.json")
	savedBriefing := readFile(t, round1Briefing)
	for _, mutation := range []string{"", "malformed"} {
		if mutation == "" {
			_ = os.Remove(round1Briefing)
		} else {
			writeFile(t, round1Briefing, mutation)
		}
		if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err == nil {
			t.Fatal("round oracle accepted missing or malformed validation/1 room")
		}
		writeFile(t, round1Briefing, savedBriefing)
	}
	validSequence := codexCommandOutput(`launcher=${SPACEDOCK_BIN:-spacedock}; "$launcher" gate record rejection-task --round validation/2 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`, "", 0, "completed") + "\n" +
		codexCommandOutput("${SPACEDOCK_BIN:-spacedock} gate prepare rejection-task validation", "", 0, "completed")
	if err := assertCodexRejectionFinalGate(entityPath, validSequence); err != nil {
		t.Fatalf("strict final-gate oracle rejected valid state: %v", err)
	}
	for _, mutation := range []string{
		codexCommandOutput("${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/2", "", 0, "completed"),
		codexCommandOutput("${SPACEDOCK_BIN:-spacedock} gate prepare rejection-task validation", "", 0, "completed") + "\n" +
			codexCommandOutput("${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/2", "", 0, "completed"),
		strings.Replace(validSequence, "rejection-task/inputs/briefing.json", "rejection-task/inputs/wrong.json", 1),
		strings.Replace(validSequence, "rejection-task/inputs/briefing.review.jsonl", "rejection-task/inputs/wrong.review.jsonl", 1),
		strings.Replace(validSequence, "gate record rejection-task", "gate record wrong-task", 1),
		strings.Replace(validSequence, "--round validation/2", "--round validation/9", 1),
		strings.Replace(validSequence, `"exit_code":0`, `"exit_code":1`, 1),
	} {
		if err := assertCodexRejectionFinalGate(entityPath, mutation); err == nil {
			t.Fatal("strict final-gate oracle accepted an invalid command sequence")
		}
	}
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
	command, result := "${SPACEDOCK_BIN:-spacedock} gate record rejection-task \\\n  --round validation/1 \\\n  --briefing /tmp/TestLiveCommonRejectionFlow2642046300/003/rejection-task/inputs/briefing.json \\\n  --log /tmp/TestLiveCommonRejectionFlow2642046300/003/rejection-task/inputs/briefing.review.jsonl \\\n  --workflow-dir /tmp/TestLiveCommonRejectionFlow2642046300/003", "round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=2\nentry=annotation:rejection-task:missing-marker type=Annotation advisory=false decision=\nentry=resolution:rejection-task:reviewer type=Resolution advisory=true decision=revise"
	entries4Command, entries4Result := "${SPACEDOCK_BIN:-spacedock} gate record rejection-task \\\n  --round validation/1 \\\n  --briefing /tmp/TestLiveCommonRejectionFlow1750522315/003/rejection-task/inputs/briefing.json \\\n  --log /tmp/TestLiveCommonRejectionFlow1750522315/003/rejection-task/inputs/briefing.review.jsonl \\\n  --workflow-dir /tmp/TestLiveCommonRejectionFlow1750522315/003 2>&1\necho \"EXIT: $?\"", "round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4\nentry=annotation:rejection-task:missing-marker type=Annotation advisory=false decision=\nentry=resolution:rejection-task:reviewer type=Resolution advisory=true decision=revise\nentry=annotation:rejection-task:fixed-marker type=Annotation advisory=false decision=\nentry=resolution:rejection-task:ensign type=Resolution advisory=true decision=revise\nEXIT: 0"
	claudeStream := strings.Join([]string{
		bashToolLine("toolu_round", command),
		toolResultLine("toolu_round", false, result),
	}, "\n")
	if !claudeRecordedRejectionRound(claudeStream) {
		t.Fatal("Claude extractor missed correlated round invocation with prefixed artifact paths")
	}
	if !claudeRecordedRejectionRound(strings.Join([]string{bashToolLine("toolu_round_entries4", entries4Command), toolResultLine("toolu_round_entries4", false, entries4Result)}, "\n")) {
		t.Fatal("Claude extractor missed retained Sonnet round invocation with four entries")
	}
	for name, invalid := range map[string]string{
		"wrong_suffix": strings.Replace(command, "briefing.json", "briefing.json.bak", 1),
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
		claudeRecordedRejectionRound(strings.Join([]string{bashToolLine("toolu_round", command), toolResultLine("toolu_round", false, strings.Replace(result, "round=round:rejection-task:validation:1", "round=malformed", 1))}, "\n")) ||
		claudeRecordedRejectionRound(strings.Join([]string{bashToolLine("toolu_round_entries4", entries4Command), toolResultLine("toolu_round_entries4", false, strings.Replace(entries4Result, "round=round:rejection-task:validation:1", "round=malformed", 1))}, "\n")) ||
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
