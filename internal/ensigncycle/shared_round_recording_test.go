package ensigncycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

var directRoundLauncher = regexp.MustCompile(`(?:^|[\s;&|])['"]?(?:spacedock|\$(?:\{SPACEDOCK_BIN(?::-[^}]*)?\}|SPACEDOCK_BIN)|/[^ \t\r\n'";&|]+/spacedock)['"]?\s+gate\s+record(?:\s|$)`)
var shellVariableReference = regexp.MustCompile(`^\$(?:[A-Za-z_][A-Za-z0-9_]*|\{[A-Za-z_][A-Za-z0-9_]*\})$`)

func commandRecordsRejectionRound(command string) bool {
	command = strings.ReplaceAll(command, "\\\n", " ")
	if !regexp.MustCompile(`(?:^|\s)['"]?rejection-task['"]?(?:\s|$)`).MatchString(command) ||
		commandFlagValue(command, "--round") != "validation/1" ||
		!commandFlagPathEndsAt(command, "--briefing", "rejection-task/inputs/briefing.json") ||
		!commandFlagPathEndsAt(command, "--log", "rejection-task/inputs/briefing.review.jsonl") ||
		!commandFlagPathEndsAt(command, "--feedback-cycle", "rejection-task/inputs/feedback-cycle.txt") {
		return false
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

func commandFlagValue(command, flag string) string {
	fields := strings.Fields(command)
	for i, field := range fields {
		if field == flag && i+1 < len(fields) {
			return strings.Trim(fields[i+1], `'"`)
		}
		if strings.HasPrefix(field, flag+"=") {
			return strings.Trim(strings.TrimPrefix(field, flag+"="), `'"`)
		}
	}
	return ""
}

func commandFlagPathEndsAt(command, flag, suffix string) bool {
	value := commandFlagValue(command, flag)
	if value == suffix || strings.HasSuffix(value, "/"+suffix) {
		return true
	}
	prefix, terminal, ok := strings.Cut(value, "/")
	return ok && shellVariableReference.MatchString(prefix) && terminal == filepath.Base(suffix)
}

func claudeRecordedRejectionRound(stream string) bool {
	for _, line := range strings.Split(stream, "\n") {
		var entry struct {
			Message *struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Input struct {
						Command string `json:"command"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &entry) != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type == "tool_use" && block.Name == "Bash" && commandRecordsRejectionRound(block.Input.Command) {
				return true
			}
		}
	}
	return false
}

func codexRecordedRejectionRound(jsonl string) bool {
	for _, line := range strings.Split(jsonl, "\n") {
		var entry codexCommandItem
		if json.Unmarshal([]byte(line), &entry) == nil &&
			entry.Item.Type == "command_execution" &&
			commandRecordsRejectionRound(entry.Item.Command) {
			return true
		}
	}
	return false
}

func rejectionRoundPhaseEntity(workflowRoot, entityPath string) ([]byte, error) {
	rel, err := filepath.Rel(workflowRoot, entityPath)
	if err != nil {
		return nil, err
	}
	rel = filepath.ToSlash(rel)
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return nil, fmt.Errorf("entity %s is outside workflow root %s", entityPath, workflowRoot)
	}
	log, err := exec.Command("git", "-C", workflowRoot, "log", "--reverse", "--format=%H", "--", rel).Output()
	if err != nil {
		return nil, fmt.Errorf("read round-record history: %w", err)
	}
	for _, commit := range strings.Fields(string(log)) {
		entity, showErr := exec.Command("git", "-C", workflowRoot, "show", commit+":"+rel).Output()
		if showErr == nil && bytes.Contains(entity, []byte(rejectionFeedbackCycle)) {
			return entity, nil
		}
	}
	return os.ReadFile(entityPath)
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
	roundPhaseEntity, err := rejectionRoundPhaseEntity(workflowRoot, entityPath)
	if err != nil {
		return err
	}
	if regexp.MustCompile(`(?m)^gates:`).Match(roundPhaseEntity) || regexp.MustCompile(`(?m)^\s+application:`).Match(roundPhaseEntity) {
		return fmt.Errorf("round recording introduced gate/application lifecycle state")
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

func TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl(t *testing.T) {
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
		"status: implementation",
		"workflow-state: preserve-me",
		"gate-state: preserve-me",
		"application-state: preserve-me",
	} {
		if strings.Count(before, line) != 1 || strings.Count(after, line) != 1 {
			t.Fatalf("round recorder did not preserve exact lifecycle line %q", line)
		}
	}
	if err := assertRejectionRecordedRound(root, entityPath, "implementation", true); err != nil {
		t.Fatalf("recorded-round durable oracle: %v", err)
	}
	if err := assertRejectionRecordedRound(root, entityPath, "implementation", false); err == nil {
		t.Fatal("inverted no-invocation control passed despite no observed launcher call")
	}

	gitCommitPathScoped(t, root, "rejection-task", "record rejection round")
	roundEntity := readFile(t, entityPath)
	writeFile(t, entityPath, roundEntity+"\ngates:\n  application: later-valid-gate-state\n")
	attributed, err := rejectionRoundPhaseEntity(root, entityPath)
	if err != nil {
		t.Fatalf("attribute round-record phase: %v", err)
	}
	if !bytes.Equal(attributed, []byte(roundEntity)) {
		t.Fatal("round-phase attribution included lifecycle state introduced after the round commit")
	}
}

func TestRejectionFlowRoundInvocationExtractors(t *testing.T) {
	command := `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --workflow-dir . --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl --feedback-cycle rejection-task/inputs/feedback-cycle.txt`
	if !claudeRecordedRejectionRound(claudeToolUse("Bash", `{"command":"`+command+`"}`)) {
		t.Fatal("Claude extractor missed resolved launcher round invocation")
	}
	captured := "B=${SPACEDOCK_BIN:-spacedock}\\n$B gate record rejection-task --workflow-dir . --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl --feedback-cycle rejection-task/inputs/feedback-cycle.txt"
	if !codexRecordedRejectionRound(codexCommand(captured)) {
		t.Fatal("Codex extractor missed captured resolved launcher round invocation")
	}
	absoluteMultiline := `/tmp/candidate/spacedock gate record \
  "rejection-task" --workflow-dir . --round="validation/1" \
  --briefing "rejection-task/inputs/briefing.json" \
  --log="rejection-task/inputs/briefing.review.jsonl" \
  --feedback-cycle "rejection-task/inputs/feedback-cycle.txt"`
	if !commandRecordsRejectionRound(absoluteMultiline) {
		t.Fatal("command recognizer missed absolute resolved launcher with multiline/quoted flags")
	}
	opusCaptured := `B=${SPACEDOCK_BIN:-spacedock}; WD=/tmp/rejection
$B gate record rejection-task --round validation/1 \
  --briefing $WD/rejection-task/inputs/briefing.json \
  --log $WD/rejection-task/inputs/briefing.review.jsonl \
  --feedback-cycle $WD/rejection-task/inputs/feedback-cycle.txt \
  --workflow-dir $WD`
	opusInput, err := json.Marshal(map[string]string{"command": opusCaptured})
	if err != nil {
		t.Fatal(err)
	}
	if !claudeRecordedRejectionRound(claudeToolUse("Bash", string(opusInput))) {
		t.Fatal("Claude extractor missed Opus variable-root round invocation")
	}
	opusInputCaptured := `WD=/tmp/rejection
IN="$WD/rejection-task/inputs"
${SPACEDOCK_BIN:-spacedock} gate record rejection-task \
  --round validation/1 \
  --briefing "$IN/briefing.json" \
  --log "$IN/briefing.review.jsonl" \
  --feedback-cycle "$IN/feedback-cycle.txt"`
	opusInputCapturedJSON, err := json.Marshal(map[string]string{"command": opusInputCaptured})
	if err != nil {
		t.Fatal(err)
	}
	if !claudeRecordedRejectionRound(claudeToolUse("Bash", string(opusInputCapturedJSON))) {
		t.Fatal("Claude extractor missed Opus variable-input-directory round invocation")
	}
	noInvocation := codexCommand("spacedock status rejection-task --workflow-dir .")
	if codexRecordedRejectionRound(noInvocation) {
		t.Fatal("no-invocation transcript falsely satisfied the round invocation extractor")
	}
}
