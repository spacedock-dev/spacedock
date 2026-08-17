package ensigncycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

var directRoundLauncher = regexp.MustCompile(`(?:^|[\s;&|])['"]*(?:spacedock|\$(?:\{SPACEDOCK_BIN(?::-[^}]*)?\}|SPACEDOCK_BIN)|/[^ \t\r\n'";&|]+/spacedock)['"]*\s+gate\s+record(?:\s|$)`)
var rejectionRoundSuccess = regexp.MustCompile(`(?m)^round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4$`)

// Quote runs here are `['"]*` for the same reason as rejectionRoundFlag below:
// Codex reaches the launcher through nested shell quoting and can wrap any
// argument, the entity operand included, in a multi-character run like `'"'`.
var rejectionRoundEntity = regexp.MustCompile(`(?:^|\s)['"]*rejection-task['"]*(?:\s|$)`)

// Quote runs are `['"]*`, not `['"]?`, because Codex emits multi-character runs
// like `'"'` around arguments. Missing a round would let the counter under-report,
// which is the one direction that could hide a second publication.
var rejectionRoundFlag = regexp.MustCompile(`--round(?:=|\s+)['"]*([A-Za-z][A-Za-z0-9_-]*/[0-9]+)['"]*(?:\s|[;&|]|$)`)

const rejectionPreparedBriefingID = "briefing:rejection-task:validation:attempt-1:revision-1"

// rejectionOpenAttemptID and rejectionOpenBriefingID match the prepared attempt and
// its Briefing at ANY attempt number. The happy path prepares attempt-1, but a
// sanctioned withdraw-and-re-prepare leaves the open gate at attempt-2, so the
// literal 1 encoded the happy path rather than the contract.
var (
	rejectionOpenAttemptID  = regexp.MustCompile(`^gate-attempt:rejection-task-validation-(\d+)$`)
	rejectionOpenBriefingID = regexp.MustCompile(`^briefing:rejection-task:validation:attempt-(\d+):revision-1$`)
)

// rejectionRoundArtifactArg matches one `--briefing`/`--log` argument pointing at
// the fixture's canonical artifact path. The quote runs are `*`, not `?`: codex
// emits the launcher through nested shell quoting, so a real invocation carries
// a RUN of quote characters against the path — `--briefing '"'rejection-task/...`
// — and a single optional quote reported a round that WAS recorded as missing.
// The flag name and the artifact path stay exact, so this relaxes only the
// quoting, not which arguments count.
func rejectionRoundArtifactArg(flag, filename string) *regexp.Regexp {
	return regexp.MustCompile(`--` + regexp.QuoteMeta(flag) + `(?:=|\s+)['"]*(?:[^ \t\r\n'";&|]*/)?rejection-task/inputs/` +
		regexp.QuoteMeta(filename) + `['"]*(?:\s|[;&|]|$)`)
}

// invokesRejectionRoundRecorder reports whether a command reaches `gate record`
// through a resolved launcher — named directly, or through a variable the same
// command captured it into.
func invokesRejectionRoundRecorder(command string) bool {
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

// rejectionRoundPublications returns every `--round` value a round-recorder
// command carries, pinning no round id. commandRecordsRejectionRound below pins
// validation/1 and so cannot see a validation/2 call at all; the counter must,
// because a second publication is the failure it exists to catch. Every value is
// returned, not just the first, because Codex reaches the recorder by chaining
// commands into one shell call and could chain both publications there.
func rejectionRoundPublications(command string) []string {
	command = strings.ReplaceAll(command, "\\\n", " ")
	if !rejectionRoundEntity.MatchString(command) || !invokesRejectionRoundRecorder(command) {
		return nil
	}
	var rounds []string
	for _, match := range rejectionRoundFlag.FindAllStringSubmatch(command, -1) {
		rounds = append(rounds, match[1])
	}
	return rounds
}

func commandRecordsRejectionRound(command string) bool {
	command = strings.ReplaceAll(command, "\\\n", " ")
	for _, required := range []*regexp.Regexp{
		rejectionRoundEntity,
		regexp.MustCompile(`--round(?:=|\s+)['"]*validation/1['"]*(?:\s|$)`),
		rejectionRoundArtifactArg("briefing", "briefing.json"),
		rejectionRoundArtifactArg("log", "briefing.review.jsonl"),
	} {
		if !required.MatchString(command) {
			return false
		}
	}
	return invokesRejectionRoundRecorder(command)
}

// assertSingleRejectionRoundPublication grades AC-2's invocation half: a rejection
// cycle publishes its round exactly once, at validation/1. It fails on zero calls
// and on the retired step-8 republication at validation/2 alike, so a flow that
// still publishes twice cannot pass by having one of the two calls look right.
func assertSingleRejectionRoundPublication(rounds []string) error {
	if len(rounds) != 1 {
		return fmt.Errorf("rejection cycle made %d successful `gate record --round` invocations %v, want exactly 1", len(rounds), rounds)
	}
	if rounds[0] != "validation/1" {
		return fmt.Errorf("rejection cycle published round %q, want validation/1", rounds[0])
	}
	return nil
}

func successfulRejectionRoundResult(content json.RawMessage, isError *bool) bool {
	var text string
	return isError != nil && !*isError && json.Unmarshal(content, &text) == nil &&
		rejectionRoundSuccess.MatchString(text)
}

type claudeToolBlock struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	ToolUseID string `json:"tool_use_id"`
	Input     struct {
		Command string `json:"command"`
	} `json:"input"`
	Content json.RawMessage `json:"content"`
	IsError *bool           `json:"is_error"`
}

func claudeToolBlocks(line string) []claudeToolBlock {
	var entry struct {
		Message *struct {
			Content []claudeToolBlock `json:"content"`
		} `json:"message"`
	}
	if json.Unmarshal([]byte(line), &entry) != nil || entry.Message == nil {
		return nil
	}
	return entry.Message.Content
}

func codexCompletedCommands(jsonl string) []codexCommandItem {
	var items []codexCommandItem
	for _, line := range strings.Split(jsonl, "\n") {
		var entry codexCommandItem
		if json.Unmarshal([]byte(line), &entry) == nil &&
			entry.Type == "item.completed" &&
			entry.Item.Type == "command_execution" {
			items = append(items, entry)
		}
	}
	return items
}

func claudeRecordedRejectionRound(stream string) bool {
	invocations := map[string]bool{}
	for _, line := range strings.Split(stream, "\n") {
		for _, block := range claudeToolBlocks(line) {
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

// claudeRejectionRoundPublications returns the round id of every round-recorder
// call in a Claude stream whose correlated tool_result did not report an error,
// in stream order. Success is the tool result's error flag, not the success line
// claudeRecordedRejectionRound pins, because a republication at validation/2 is
// still a publication and is exactly what must be counted. An absent flag counts
// as success: the counter is only ever allowed to over-report, since a missed
// second call would let a double-publishing flow pass.
func claudeRejectionRoundPublications(stream string) []string {
	pending := map[string][]string{}
	rounds := []string{}
	for _, line := range strings.Split(stream, "\n") {
		for _, block := range claudeToolBlocks(line) {
			if block.Type == "tool_use" && block.Name == "Bash" && block.ID != "" {
				if published := rejectionRoundPublications(block.Input.Command); len(published) > 0 {
					pending[block.ID] = published
				}
			}
			if block.Type == "tool_result" && (block.IsError == nil || !*block.IsError) {
				rounds = append(rounds, pending[block.ToolUseID]...)
				delete(pending, block.ToolUseID)
			}
		}
	}
	return rounds
}

func codexRecordedRejectionRound(jsonl string) bool {
	for _, entry := range codexCompletedCommands(jsonl) {
		if entry.Item.ExitCode != nil && *entry.Item.ExitCode == 0 &&
			commandRecordsRejectionRound(entry.Item.Command) {
			return true
		}
	}
	return false
}

// codexRejectionRoundPublications is the Codex half of the publication counter:
// the round id of every round-recorder call that exited 0, in stream order.
func codexRejectionRoundPublications(jsonl string) []string {
	rounds := []string{}
	for _, entry := range codexCompletedCommands(jsonl) {
		if entry.Item.ExitCode != nil && *entry.Item.ExitCode == 0 {
			rounds = append(rounds, rejectionRoundPublications(entry.Item.Command)...)
		}
	}
	return rounds
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
	// Select the ONE open attempt; withdrawn siblings are tolerated. A withdraw
	// followed by a re-prepare is the sanctioned recovery (the withdrawn-gate-recovery
	// journey in this same suite), and the end state it leaves — one open gate, the
	// stale attempt withdrawn and never presented — satisfies both this journey's
	// determined shape and the registry's required outcome. Counting attempts instead
	// accused a run that had recovered exactly as prescribed. A resolved or applied
	// attempt is neither open nor withdrawn and still fails: the journey must stop at
	// a gate nobody has decided.
	var open []gates.Attempt
	for _, candidate := range doc.Records[0].Attempts {
		switch {
		case candidate.Withdrawal != nil:
			continue
		case candidate.Resolution != nil || candidate.Application != nil:
			return fmt.Errorf("prepared validation gate is not open")
		}
		open = append(open, candidate)
	}
	if len(open) != 1 {
		return fmt.Errorf("selected validation gate holds %d open attempts, want exactly one (withdrawn attempts tolerated)", len(open))
	}
	attempt := open[0]
	if attempt.Briefing.ID == rejectionBriefingID {
		return fmt.Errorf("validation/1 advisory round was retained as a gate attempt")
	}
	// Both ids generalize from attempt-1 to attempt-N, and their N must AGREE. The
	// binding — this briefing belongs to THIS open attempt, not to a withdrawn
	// sibling — is what the pin was ever for; the literal 1 was an accident of the
	// happy path, and keeping it moved the CI red to a different line rather than
	// removing it.
	attemptNumber := rejectionOpenAttemptID.FindStringSubmatch(attempt.ID)
	briefingNumber := rejectionOpenBriefingID.FindStringSubmatch(attempt.Briefing.ID)
	if attemptNumber == nil || briefingNumber == nil || attemptNumber[1] != briefingNumber[1] {
		return fmt.Errorf("selected validation gate is not bound to the expected prepared Briefing")
	}
	return nil
}

// assertRejectionGatePrepared grades the requirement that the FO left the cycle-2
// validation gate PREPARED, not merely round-recorded. It is deliberately separate
// from assertRejectionRecordedRound: that oracle's boundary check
// (assertRejectionRoundGateBoundary, directly above) TOLERATES an entity with no
// gates record and returns nil. Folding this condition into the round oracle's
// result made one journey hold two contradictory positions on the same state and
// report both under `rejection-round-missing`, so a run whose round WAS recorded
// but whose gate was never prepared was diagnosed as a missing round. Each
// condition now carries its own code.
func assertRejectionGatePrepared(entityPath string) error {
	if _, _, err := gates.Read(entityPath); err != nil {
		return fmt.Errorf("FO never prepared the cycle-2 validation gate: %v", err)
	}
	return nil
}

// assertRejectionCycleLine grades the one durable Cycle line this fixture determines.
// The fixture tells the FO to copy `rejection-task/inputs/feedback-cycle.txt` verbatim
// for the rejection round, and says a re-validation that passes "closes its cycle
// without adding a line" — so the conforming end state is EXACTLY one entry in
// `### Feedback Cycles`, byte-equal to that file. Grading the count and the bytes
// together makes both directions falsifiable: an FO that never records the round
// leaves zero entries, and one that appends a cycle-2 line anyway leaves two. The
// audit found this line was fixture-determined but never asserted at all.
func assertRejectionCycleLine(entityPath string) error {
	entity, err := os.ReadFile(entityPath)
	if err != nil {
		return &gradedErr{code: "rejection-cycle-line", msg: fmt.Sprintf("cannot read entity for the Cycle line: %v", err)}
	}
	var entries []string
	for _, line := range strings.Split(feedbackCyclesSection(string(entity)), "\n") {
		if feedbackCycleEntry.MatchString(line) {
			entries = append(entries, strings.TrimRight(line, " \t\r"))
		}
	}
	want := strings.TrimRight(rejectionFeedbackCycle, "\n")
	// Honest label for the heading-drift case. The count alone reads as "the FO
	// never recorded the round", which is what a real tip-CI red implied about an FO
	// that HAD written the line byte-exact — under `## Feedback Cycles` instead of
	// the declared `### Feedback Cycles`. The grade does not soften (the declared
	// grammar is the contract), but the diagnostic must say what actually happened
	// or the next reader debugs the wrong thing.
	if len(entries) == 0 {
		if at := strings.Index(string(entity), want); at >= 0 {
			return &gradedErr{code: "rejection-cycle-line", msg: fmt.Sprintf(
				"the Cycle line is byte-exact but sits outside the declared `### Feedback Cycles` section, under %q — the round WAS recorded; the projection is declared at that exact heading level and a line under another level is not in it",
				headingAbove(string(entity), at))}
		}
	}
	if len(entries) != 1 {
		return &gradedErr{code: "rejection-cycle-line", msg: fmt.Sprintf(
			"`### Feedback Cycles` holds %d `- Cycle N:` entries, want exactly 1 (the rejection round; a passing re-validation adds none): %q", len(entries), entries)}
	}
	if entries[0] != want {
		return &gradedErr{code: "rejection-cycle-line", msg: fmt.Sprintf(
			"the Cycle line is not the fixture's verbatim text\n got: %q\nwant: %q", entries[0], want)}
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
	// The fixture's journey is ONE rejection cycle, so its durable end state is one
	// round room and a pointer that resolves it. The round-id fallback this replaced
	// accepted a validation/2 pointer as well, which absorbed the second publication
	// instead of failing on it. Pinning validation/1 is what makes a republication
	// falsifiable here; it is a fact about this one-cycle fixture, not a claim that
	// a workflow may hold only one room.
	// `review/validation/` holds the gate's own `briefing-N` rooms alongside the
	// round rooms, so only the `round-` prefix is counted here.
	stageRooms, err := os.ReadDir(filepath.Join(filepath.Dir(entityPath), "review", "validation"))
	if err != nil {
		return fmt.Errorf("read retained round rooms: %w", err)
	}
	var roundRooms []string
	for _, entry := range stageRooms {
		if strings.HasPrefix(entry.Name(), "round-") {
			roundRooms = append(roundRooms, entry.Name())
		}
	}
	if len(roundRooms) != 1 || roundRooms[0] != "round-1" {
		return fmt.Errorf("rejection cycle left round rooms %v, want exactly one round-1 room", roundRooms)
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
		t.Fatalf("recorded-round oracle rejected the later open validation gate: %v", err)
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
	// The sanctioned withdraw-and-re-prepare recovery, driven through the real gate
	// verbs rather than hand-edited YAML so the shape is the one the launcher
	// actually writes. This is the tip-CI claude run's shape: a premature attempt-1
	// withdrawn with a reason, then attempt-2 prepared and left open, with the round
	// already recorded. It must grade GREEN — the stale authority was withdrawn and
	// never presented, which is exactly what the registry's required outcome asks
	// for. Pinning attempt-1 red it.
	writeFile(t, entityPath, openGateEntity)
	if _, err := gates.Withdraw(entityPath, gates.WithdrawInput{
		WorkflowDir: root,
		Reason:      "prepared before the reviewer re-ran; withdrawing the stale attempt",
	}); err != nil {
		t.Fatalf("withdraw the stale validation attempt: %v", err)
	}
	reprepared, err := gates.Prepare(entityPath, gates.PrepareInput{
		WorkflowDir: root,
		Question:    "Does the second validation confirm the fix marker is present and PASS?",
		Artifact:    gateReviewPath,
		Summary:     "The corrected rejection-flow candidate is ready for a decision.",
	})
	if err != nil {
		t.Fatalf("re-prepare the validation gate after withdrawal: %v", err)
	}
	if !strings.HasSuffix(reprepared.Briefing, "attempt-2:revision-1") || reprepared.State != "open" {
		t.Fatalf("re-prepared gate = %#v, want an open attempt-2 briefing", reprepared)
	}
	if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err != nil {
		t.Fatalf("recorded-round oracle red the sanctioned withdraw-and-re-prepare recovery: %v", err)
	}
	recoveredEntity := readFile(t, entityPath)
	// A withdrawn sibling is tolerated, but a DECIDED attempt is not: the journey
	// stops at a gate nobody has resolved. Closing the re-prepared attempt-2 must
	// still red, or tolerating withdrawal would have blessed every non-open state.
	if err := gates.RecordSemantic(entityPath, gates.RecordInput{
		Decision: "hold", Actor: "person:captain", Reason: "exercise closed-gate counterexample", WorkflowDir: root,
	}); err != nil {
		t.Fatalf("close later validation gate for counterexample: %v", err)
	}
	if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err == nil ||
		!strings.Contains(err.Error(), "prepared validation gate is not open") {
		t.Fatalf("closed gate control diagnostic = %v", err)
	}
	// And an entity left with NO open attempt at all — every attempt withdrawn, none
	// re-prepared — must red too, so "tolerated" never becomes "unnecessary".
	writeFile(t, entityPath, recoveredEntity)
	if _, err := gates.Withdraw(entityPath, gates.WithdrawInput{
		WorkflowDir: root, Reason: "withdraw the replacement too, leaving nothing open",
	}); err != nil {
		t.Fatalf("withdraw the re-prepared attempt: %v", err)
	}
	if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err == nil ||
		!strings.Contains(err.Error(), "holds 0 open attempts") {
		t.Fatalf("no-open-attempt control diagnostic = %v", err)
	}

	// Republication counterexample, last because nothing undoes it: the retired
	// step-8 `validation/2` call is what the round-id fallback used to absorb, and
	// the oracle must now fail on it. Without this the single-publication claim is
	// unfalsifiable — a flow that publishes twice would still pass.
	writeFile(t, entityPath, openGateEntity)
	if err := gates.RecordSemantic(entityPath, gates.RecordInput{
		Round:        "validation/2",
		BriefingPath: filepath.Join(root, "rejection-task", "inputs", "briefing.json"),
		LogPath:      filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"),
		WorkflowDir:  root,
	}); err != nil {
		t.Fatalf("record second validation round: %v", err)
	}
	if err := assertRejectionRecordedRound(root, entityPath, "validation", true); err == nil ||
		!strings.Contains(err.Error(), "want exactly one round-1 room") {
		t.Fatalf("republication control diagnostic = %v, want the second round room rejected", err)
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

// TestRejectionRoundPublicationCounter pins the counter that grades AC-2's
// invocation half on both hosts. The claim is "exactly one publication, at
// validation/1", so the cases that must FAIL are the ones the retired flow
// produced: the second call at validation/2, and — on Codex, whose failure was a
// dropped tail — no call at all. A publication chained into one shell command is
// counted separately, because Codex reaches the recorder that way.
func TestRejectionRoundPublicationCounter(t *testing.T) {
	roundCommand := func(round string) string {
		return `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --workflow-dir . --round ` + round +
			` --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`
	}
	claudeCall := func(id, round string, failed bool) string {
		return strings.Join([]string{
			bashToolLine(id, roundCommand(round)),
			toolResultLine(id, failed, "round="+round),
		}, "\n")
	}
	for name, tc := range map[string]struct {
		claude, codex string
		want          string
	}{
		"one call at validation/1": {
			claude: claudeCall("toolu_1", "validation/1", false),
			codex:  codexCommandOutput(roundCommand("validation/1"), "", 0, "completed"),
		},
		"republished at validation/2": {
			claude: claudeCall("toolu_1", "validation/1", false) + "\n" + claudeCall("toolu_2", "validation/2", false),
			codex: codexCommandOutput(roundCommand("validation/1"), "", 0, "completed") + "\n" +
				codexCommandOutput(roundCommand("validation/2"), "", 0, "completed"),
			want: "made 2 successful",
		},
		"never published": {
			claude: bashToolLine("toolu_1", "spacedock status rejection-task --workflow-dir .") + "\n" +
				toolResultLine("toolu_1", false, "{}"),
			codex: codexCommandOutput("spacedock status rejection-task --workflow-dir .", "", 0, "completed"),
			want:  "made 0 successful",
		},
		"only publication is validation/2": {
			claude: claudeCall("toolu_1", "validation/2", false),
			codex:  codexCommandOutput(roundCommand("validation/2"), "", 0, "completed"),
			want:   `published round "validation/2"`,
		},
		"failed call is not a publication": {
			claude: claudeCall("toolu_1", "validation/1", true),
			codex:  codexCommandOutput(roundCommand("validation/1"), "", 1, "failed"),
			want:   "made 0 successful",
		},
		"codex multi-character quote run around the round": {
			claude: strings.Join([]string{
				bashToolLine("toolu_1", `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --workflow-dir . --round '"'validation/2'"' --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`),
				toolResultLine("toolu_1", false, "round=validation/2"),
			}, "\n"),
			codex: codexCommandOutput(`${SPACEDOCK_BIN:-spacedock} gate record rejection-task --workflow-dir . --round '"'validation/2'"' --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`, "", 0, "completed"),
			want:  `published round "validation/2"`,
		},
		"both publications chained into one command": {
			claude: strings.Join([]string{
				bashToolLine("toolu_1", roundCommand("validation/1")+"; "+roundCommand("validation/2")),
				toolResultLine("toolu_1", false, "round=validation/2"),
			}, "\n"),
			codex: codexCommandOutput(roundCommand("validation/1")+"; "+roundCommand("validation/2"), "", 0, "completed"),
			want:  "made 2 successful",
		},
	} {
		t.Run(name, func(t *testing.T) {
			for host, got := range map[string]error{
				"claude": assertSingleRejectionRoundPublication(claudeRejectionRoundPublications(tc.claude)),
				"codex":  assertSingleRejectionRoundPublication(codexRejectionRoundPublications(tc.codex)),
			} {
				switch {
				case tc.want == "" && got != nil:
					t.Fatalf("%s counter rejected the single publication: %v", host, got)
				case tc.want != "" && got == nil:
					t.Fatalf("%s counter accepted %s", host, name)
				case tc.want != "" && !strings.Contains(got.Error(), tc.want):
					t.Fatalf("%s counter diagnostic = %v, want %q", host, got, tc.want)
				}
			}
		})
	}
}

// TestRejectionRoundRecognizerAcceptsNestedShellQuoting pins the quoting shape a
// real codex FO emits. Captured verbatim from a live rejection-flow run: codex
// wraps the launcher in nested shell quoting, so `--briefing` is followed by the
// run `'"'` before the path. The recognizer previously allowed a single optional
// quote there, so this exact command — which exited 0 and printed the complete
// four-entry round summary — was graded `rejection-round-missing`. Narrowing the
// quote runs back to `?` reds this test. The wrong-artifact controls confirm the
// relaxation did not make the recognizer accept any artifact path.
func TestRejectionRoundRecognizerAcceptsNestedShellQuoting(t *testing.T) {
	// Verbatim capture, transcribed byte-for-byte from a live run's
	// codex-exec.jsonl including its own t.TempDir path. Retyping the shape by
	// hand risks silently normalizing the very quoting under test.
	const captured = `/bin/zsh -lc '${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 --briefing '"'rejection-task/inputs/briefing.json' --log 'rejection-task/inputs/briefing.review.jsonl' --workflow-dir '/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCommonRejectionFlow3029195893/002'"`
	if !commandRecordsRejectionRound(captured) {
		t.Fatal("recognizer rejected the nested-shell quoting a real codex run emits")
	}
	result := "round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4"
	if !codexRecordedRejectionRound(codexCommandOutput(captured, result, 0, "completed")) {
		t.Fatal("codex extractor missed the nested-shell-quoted round invocation")
	}
	// A quote run can land on any argument, so every operand the recognizer pins
	// must tolerate one. Each case below reds if its own site narrows back to `?`.
	for site, quoted := range map[string]string{
		"entity":   `${SPACEDOCK_BIN:-spacedock} gate record '"'rejection-task'"' --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`,
		"round":    `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round '"'validation/1'"' --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`,
		"briefing": `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 --briefing '"'rejection-task/inputs/briefing.json'"' --log rejection-task/inputs/briefing.review.jsonl`,
		"log":      `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 --briefing rejection-task/inputs/briefing.json --log '"'rejection-task/inputs/briefing.review.jsonl'"'`,
	} {
		t.Run("quote run on "+site, func(t *testing.T) {
			if !commandRecordsRejectionRound(quoted) {
				t.Fatalf("recognizer rejected a quote run around the %s operand", site)
			}
		})
	}
	for name, invalid := range map[string]string{
		"wrong_briefing_file": strings.Replace(captured, "briefing.json", "other.json", 1),
		"wrong_log_file":      strings.Replace(captured, "briefing.review.jsonl", "other.review.jsonl", 1),
		"wrong_entity":        strings.Replace(captured, "gate record rejection-task", "gate record other-task", 1),
		"wrong_round":         strings.Replace(captured, "--round validation/1", "--round validation/2", 1),
	} {
		t.Run(name, func(t *testing.T) {
			if commandRecordsRejectionRound(invalid) {
				t.Fatal("relaxed quoting made the recognizer accept an invalid argument")
			}
		})
	}
}

// TestRejectionUnpreparedGateReportsItsOwnCode pins AC-3's separation. The input is
// the exact durable end state both failing CI runs produced: the round IS recorded
// (the codex stream oracle says so) and the entity status reached validation, but
// the FO never prepared the cycle-2 gate, so the entity carries no gates record.
// Before this split, that state graded as `rejection-round-missing` — the round
// oracle was told a recorded round had not been recorded — and the real condition
// was invisible. The lane must now name the unprepared gate and must NOT claim the
// round is missing. Flipping either oracle back to the other's code fails this.
func TestRejectionUnpreparedGateReportsItsOwnCode(t *testing.T) {
	root := t.TempDir()
	entityPath := writeRejectionWorkflow(t, root)
	writeFile(t, filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"), rejectionCompleteLog())
	if err := gates.RecordSemantic(entityPath, gates.RecordInput{
		Round:        "validation/1",
		BriefingPath: filepath.Join(root, "rejection-task", "inputs", "briefing.json"),
		LogPath:      filepath.Join(root, "rejection-task", "inputs", "briefing.review.jsonl"),
		WorkflowDir:  root,
	}); err != nil {
		t.Fatalf("record rejection round: %v", err)
	}
	writeFile(t, entityPath, strings.Replace(readFile(t, entityPath), "status: backlog", "status: validation", 1))

	// The stream the FO actually produced: the round-recording invocation ran and
	// exited 0. Nothing about round recording failed in either CI run.
	stream := codexCommandOutput(
		"${SPACEDOCK_BIN:-spacedock} gate record rejection-task --workflow-dir . --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl",
		"round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4",
		0, "completed")
	recordedRound := codexRecordedRejectionRound(stream)
	if !recordedRound {
		t.Fatal("fixture stream must satisfy the round oracle; otherwise this test proves nothing about the gate")
	}
	if _, _, err := gates.Read(entityPath); err == nil {
		t.Fatal("fixture entity must carry no gates record; otherwise the unprepared-gate condition is absent")
	}

	grade := gradeLive(false,
		durableSemantic("rejection-round-missing", assertRejectionRecordedRound(root, entityPath, "validation", recordedRound)),
		durableSemantic("rejection-gate-not-prepared", assertRejectionGatePrepared(entityPath)))
	if !reflect.DeepEqual(grade.codes, []string{"rejection-gate-not-prepared"}) {
		t.Fatalf("grade codes = %v, want exactly [rejection-gate-not-prepared]", grade.codes)
	}
	if len(grade.details) != 1 || !strings.Contains(grade.details[0], "never prepared the cycle-2 validation gate") {
		t.Fatalf("grade details = %v, want the unprepared-gate finding's own message", grade.details)
	}

	// Control: prepare the gate and the code clears, so the code tracks the gate
	// rather than some unrelated property of the fixture.
	gateReviewPath := filepath.Join(root, "rejection-task", "inputs", "gate-validation", "gate-review.md")
	writeFile(t, gateReviewPath, "# Rejection Task — validation review\n\nReady for its prepared decision gate.\n")
	gateReviewRel, err := filepath.Rel(root, gateReviewPath)
	if err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", "--", gateReviewRel)
	git(t, root, "commit", "-q", "-m", "add prepared validation review", "--", gateReviewRel)
	if _, err := gates.Prepare(entityPath, gates.PrepareInput{
		WorkflowDir: root,
		Question:    "Does the second validation confirm the fix marker is present and PASS?",
		Artifact:    gateReviewPath,
		Summary:     "The corrected rejection-flow candidate is ready for a decision.",
	}); err != nil {
		t.Fatalf("prepare validation gate: %v", err)
	}
	if grade := gradeLive(false,
		durableSemantic("rejection-round-missing", assertRejectionRecordedRound(root, entityPath, "validation", recordedRound)),
		durableSemantic("rejection-gate-not-prepared", assertRejectionGatePrepared(entityPath))); grade.status != "pass" {
		t.Fatalf("prepared gate must grade pass; got %#v", grade)
	}
}
