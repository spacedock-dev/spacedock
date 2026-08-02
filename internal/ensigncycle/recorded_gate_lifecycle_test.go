package ensigncycle

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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

var recordedGateRequiredEvents = []string{
	"prepare",
	"decision-record",
	"consume",
}

const recordedGateDispatchMarker = "RECORDED-GATE-SUCCESSOR-DISPATCHED"

const (
	recordedGateBriefingID = "briefing:recorded-gate-task:validation:attempt-1:revision-1"
	recordedGateDigest     = "sha256:61776a9cdacc5e71a977d72a3a6f81808e9cda4bb2d59df01ada38b0bf78f737"
	recordedGateReason     = "accepts-direction evidence: preserve the reviewed package after the presented 3k validation gate."
	recordedGateDirective  = "you have the conn toward the sprint goal; authorized to approve gates, PR, relevant CI lanes, and merge; use your judgement."
)

type recordedGateDispatchProof struct {
	builds, successfulBuilds int
	durableEffects           int
	ordered                  bool
	committed                bool
}

type recordedGateObservation struct {
	events       []string
	before       string
	after        string
	dispatch     recordedGateDispatchProof
	expectedNext string
	gateID       string
	attemptID    string
	briefingID   string
	digest       string
	resolutionID string
}

func assertRecordedGateLifecycle(o recordedGateObservation) error {
	gateID, attemptID := o.gateID, o.attemptID
	briefingID, digest, resolutionID := o.briefingID, o.digest, o.resolutionID
	if gateID == "" {
		gateID = "gate:recorded-gate-task:validation"
	}
	if attemptID == "" {
		attemptID = "gate-attempt:recorded-gate-task-validation-1"
	}
	if briefingID == "" {
		briefingID = recordedGateBriefingID
	}
	if digest == "" {
		digest = recordedGateDigest
	}
	if resolutionID == "" {
		resolutionID = "resolution:spacedock:recorded-gate-task:validation:1"
	}
	if len(o.events) != len(recordedGateRequiredEvents) {
		return fmt.Errorf("gate lifecycle recorded %d events, want %d before dispatch", len(o.events), len(recordedGateRequiredEvents))
	}
	for i, want := range recordedGateRequiredEvents {
		if got := o.events[i]; got != want {
			return fmt.Errorf("gate lifecycle event %d = %q, want %q", i+1, got, want)
		}
	}
	if o.dispatch.builds != 1 || o.dispatch.successfulBuilds != 1 {
		return fmt.Errorf("successor dispatch build attempts/successes = %d/%d, want 1/1", o.dispatch.builds, o.dispatch.successfulBuilds)
	}
	if o.dispatch.durableEffects != 1 {
		return fmt.Errorf("new durable successor effects = %d, want 1", o.dispatch.durableEffects)
	}
	if !o.dispatch.ordered || !o.dispatch.committed {
		return fmt.Errorf("successor dispatch was not observed after consume")
	}
	if strings.Contains(o.before, recordedGateDispatchMarker) {
		return fmt.Errorf("successor marker already existed before the lifecycle began")
	}
	for _, want := range []string{"status: " + o.expectedNext, "state: consumed", "by: agent:first-officer"} {
		if !strings.Contains(o.after, want) {
			return fmt.Errorf("durable post-state missing %q", want)
		}
	}
	authority := o.after
	if strings.HasPrefix(o.after, "---\n") {
		if end := strings.Index(o.after[len("---\n"):], "\n---\n"); end >= 0 {
			authority = o.after[:len("---\n")+end]
		}
	}
	for _, exact := range []struct {
		label string
		value string
		count int
	}{
		{"gate record identity", "id: " + gateID, 1},
		{"gate record stage", "stage: validation", 1},
		{"attempt identity", "id: " + attemptID, 1},
		{"briefing identity", "id: " + briefingID, 1},
		{"briefing resolution link", "briefing: " + briefingID, 1},
		{"briefing digest", "digest: " + digest, 1},
		{"resolution identity", "id: " + resolutionID, 1},
		{"approval decision", "\n                decision: approve", 1},
		{"approval actor", "by: agent:first-officer", 1},
		{"approval reason", "\n                reason:", 1},
		{"forged adoption note", "adoption-note:", 0},
		{"application target", "target-stage: " + o.expectedNext, 1},
		{"consumed application", "\n                state: consumed", 1},
	} {
		if got := strings.Count(authority, exact.value); got != exact.count || (exact.label == "approval reason" && strings.Trim(strings.TrimSpace(strings.SplitN(strings.SplitN(authority, exact.value, 2)[1], "\n", 2)[0]), `"'`) == "") {
			return fmt.Errorf("durable post-state %s count = %d, want %d for %q", exact.label, got, exact.count, exact.value)
		}
	}
	if report := strings.Split(o.after, "\n## Stage Report: handoff\n"); len(report) != 2 || !strings.Contains(strings.SplitN(report[1], "\n## ", 2)[0], "\n- DONE: ") {
		return fmt.Errorf("durable post-state lacks exactly one handoff Stage Report with DONE evidence")
	}
	if o.before == o.after {
		return fmt.Errorf("gate lifecycle left entity byte-identical")
	}
	return nil
}

type recordedGateFixture struct {
	root       string
	stateRoot  string
	entity     string
	gateReview string
	references []string
}

type recordedGateCommand struct {
	event  string
	argv   []string
	exit   int
	stdout string
	stderr string
}

func TestRecordedGateLifecycleRealCLIReplay(t *testing.T) {
	fixture := writePreparedRecordedGateFixture(t)
	commandLog := filepath.Join(fixture.root, "command.log")
	binary := filepath.Join(writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog), "spacedock")
	before := readFile(t, fixture.entity)
	var commands []recordedGateCommand
	run := func(event string, args ...string) recordedGateCommand {
		t.Helper()
		cmd := runRecordedGateCommand(binary, fixture.root, event, args...)
		commands = append(commands, cmd)
		if cmd.exit != 0 {
			t.Fatalf("%s exit=%d\nstdout=%s\nstderr=%s", event, cmd.exit, cmd.stdout, cmd.stderr)
		}
		return cmd
	}

	mustRecordedGate(t, binary, fixture.root, "gate", "--help")
	prepareArgs := []string{"gate", "prepare", "recorded-gate-task",
		"--question", "Should the recorded validation gate advance?",
		"--artifact", fixture.gateReview,
		"--summary", "Exact recorded gate validation summary.",
	}
	for _, reference := range fixture.references {
		prepareArgs = append(prepareArgs, "--reference", reference)
	}
	prepareArgs = append(prepareArgs, "--workflow-dir", fixture.root)
	prepared := run("prepare", prepareArgs...)
	assertCommandOutput(t, prepared.stdout, "state=open")
	preparedRoom := outputValue(prepared.stdout, "room")
	if preparedRoom == "" || outputValue(prepared.stdout, "briefing") == "" || outputValue(prepared.stdout, "digest") == "" {
		t.Fatalf("prepare output is incomplete: %q", prepared.stdout)
	}
	commitRecordedGateState(t, binary, fixture, "bind prepared recorder-ready room")

	close := run("decision-record", "gate", "record", "recorded-gate-task",
		"--decision", "approve", "--actor", "agent:first-officer",
		"--reason", recordedGateReason,
		"--workflow-dir", fixture.root)
	assertCommandOutput(t, close.stdout, "state=closed", "decision=approve")
	commitRecordedGateState(t, binary, fixture, "record delegated gate decision")

	consume := run("consume", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
	assertCommandOutput(t, consume.stdout, "consumed=true", "target-stage=handoff")
	commitRecordedGateState(t, binary, fixture, "consume gate authorization")
	durable, _, durableErr := gates.Read(fixture.entity)
	requireRecordedGate(t, durableErr == nil && durable.Records[0].Attempts[0].Resolution.By == "agent:first-officer" && durable.Records[0].Attempts[0].Resolution.Reason == recordedGateReason, "approve durable snapshot unreadable")
	assertRecordedGateDispatchRow(t, binary, fixture, "handoff", "handoff")

	events := successfulRecordedGateEvents(commands)
	dispatches := 0
	if err := authorizeRecordedGateDispatch(events, readFile(t, fixture.entity), "handoff"); err == nil {
		checklist := filepath.Join(fixture.root, "handoff.checklist")
		writeFile(t, checklist, "- Append the successor marker\n")
		mustRecordedGate(t, binary, fixture.root, "dispatch", "build", "--workflow-dir", fixture.root, "--entity-path", fixture.entity, "--stage", "handoff", "--checklist-file", checklist, "--host", "claude", "--bare-mode")
		dispatches++
	} else {
		t.Fatalf("dispatch oracle refused complete lifecycle: %v", err)
	}
	zero := recordedGateLiveObservation(t, fixture, before, commandLog)
	requireRecordedGate(t, zero.dispatch.builds == 1 && zero.dispatch.durableEffects == 0 && assertRecordedGateLifecycle(zero) != nil, "zero-effect executed build qualified")
	writeFile(t, fixture.entity, readFile(t, fixture.entity)+"\n"+recordedGateDispatchMarker+"\n\n## Stage Report: handoff\n\n- DONE: Successor dispatch followed decision: approve.\n  The one-use application was already consumed before dispatch.\n\n### Summary\n\nThe entered handoff completed after the consumed approval.\n")
	gitCommitPathScoped(t, fixture.stateRoot, "recorded-gate-task/index.md", "record successor effect")
	assertRecordedGateDispatchRow(t, binary, fixture, "handoff", "done")
	writeRecordedGateEvidence(t, fixture.root, commands, before, readFile(t, fixture.entity), dispatches)
	observation := recordedGateLiveObservation(t, fixture, before, commandLog)
	if err := assertRecordedGateLifecycle(observation); err != nil {
		t.Fatal(err)
	}
	validLog := readFile(t, commandLog)
	withoutHelp := strings.ReplaceAll(strings.ReplaceAll(validLog, "begin\tgate --help\n", ""), "exit=0\tgate --help\n", "")
	writeFile(t, commandLog, withoutHelp)
	requireRecordedGate(t, assertRecordedGateLifecycle(recordedGateLiveObservation(t, fixture, before, commandLog)) == nil, "optional gate help changed lifecycle ordering")
	writeFile(t, commandLog, validLog)
	for name, log := range map[string]string{"zero-build": strings.Replace(validLog, "begin\tdispatch build ", "begin\tignored build ", 1), "failed-build": strings.Replace(validLog, "exit=0\tdispatch build ", "exit=1\tdispatch build ", 1), "build-before-consume": strings.Replace(validLog, "exit=0\tgate consume ", "exit=0\tignored consume ", 1) + "\nexit=0\tgate consume late", "missing-ancestry": strings.Replace(validLog, "dispatch-head\t", "missing-head\t", 1)} {
		writeFile(t, commandLog, log)
		requireRecordedGate(t, assertRecordedGateLifecycle(recordedGateLiveObservation(t, fixture, before, commandLog)) != nil, "%s control qualified", name)
	}
	writeFile(t, commandLog, validLog)
	mustRecordedGate(t, binary, fixture.root, "dispatch", "build", "--workflow-dir", fixture.root, "--entity-path", fixture.entity, "--stage", "handoff", "--checklist-file", filepath.Join(fixture.root, "handoff.checklist"), "--host", "claude", "--bare-mode")
	two := recordedGateLiveObservation(t, fixture, before, commandLog)
	requireRecordedGate(t, two.dispatch.builds == 2 && assertRecordedGateLifecycle(two) != nil, "two-build control qualified")
	writeFile(t, commandLog, validLog)
	writeFile(t, fixture.entity, readFile(t, fixture.entity)+"\n"+recordedGateDispatchMarker+"-SECOND\n")
	gitCommitPathScoped(t, fixture.stateRoot, "recorded-gate-task/index.md", "record duplicate effect")
	two = recordedGateLiveObservation(t, fixture, before, commandLog)
	requireRecordedGate(t, two.dispatch.durableEffects == 2 && assertRecordedGateLifecycle(two) != nil, "two-effect control qualified")

	log := git(t, fixture.stateRoot, "show", "--name-only", "--format=", "HEAD~5..HEAD")
	for _, want := range []string{
		"recorded-gate-task/index.md",
		"recorded-gate-task/review/validation/briefing-1/gate-briefing.json",
		"recorded-gate-task/review/validation/briefing-1/request.json",
	} {
		if !strings.Contains(log, want) {
			t.Errorf("folder-form state commits omitted %s:\n%s", want, log)
		}
	}
	if strings.Contains(log, "dirty-sibling.md") {
		t.Fatalf("folder-form state commit swept dirty sibling:\n%s", log)
	}
}

func TestRecordedGateLifecycleWithdrawColdBootReplaceAndConsume(t *testing.T) {
	fixture := writePreparedRecordedGateFixture(t)
	binary := buildRecordedGateBinary(t)

	prepare := mustRecordedGate(t, binary, fixture.root,
		"gate", "prepare", "recorded-gate-task",
		"--question", "Should the stale candidate advance?",
		"--artifact", fixture.gateReview,
		"--summary", "Stale candidate.",
		"--reference", fixture.references[0],
		"--reference", fixture.references[1],
		"--workflow-dir", fixture.root)
	firstRoom := outputValue(prepare.stdout, "room")
	firstBriefing := readFile(t, filepath.Join(firstRoom, "gate-briefing.json"))
	firstRequest := readFile(t, filepath.Join(firstRoom, "request.json"))
	commitRecordedGateState(t, binary, fixture, "prepare stale attempt")

	withdraw := mustRecordedGate(t, binary, fixture.root,
		"gate", "withdraw", "recorded-gate-task",
		"--reason", "Sprint re-scope replaced the reviewed candidate.",
		"--workflow-dir", fixture.root)
	assertCommandOutput(t, withdraw.stdout, "state=withdrawn", "attempt=gate-attempt:recorded-gate-task-validation-1")
	commitRecordedGateState(t, binary, fixture, "withdraw stale attempt")

	boot := mustRecordedGate(t, binary, fixture.root,
		"status", "--workflow-dir", fixture.root, "--boot", "--identify", "--json")
	assertCommandOutput(t, boot.stdout,
		`"ready_gates":[{"id":"recorded-gate-task","slug":"recorded-gate-task","current":"validation","readiness":"withdrawn-awaiting-prepare"}]`)

	replacement := mustRecordedGate(t, binary, fixture.root,
		"gate", "prepare", "recorded-gate-task",
		"--question", "Should the replacement candidate advance?",
		"--artifact", fixture.gateReview,
		"--summary", "Replacement candidate.",
		"--reference", fixture.references[0],
		"--reference", fixture.references[1],
		"--workflow-dir", fixture.root)
	secondRoom := outputValue(replacement.stdout, "room")
	if secondRoom == firstRoom {
		t.Fatal("replacement reused withdrawn room")
	}
	assertCommandOutput(t, replacement.stdout, "briefing=briefing:recorded-gate-task:validation:attempt-2:revision-1")
	commitRecordedGateState(t, binary, fixture, "prepare replacement attempt")
	assertCommandOutput(t, mustRecordedGate(t, binary, fixture.root,
		"status", "--workflow-dir", fixture.root, "--boot", "--identify", "--json").stdout,
		`"readiness":"awaiting-captain"`)

	writeRecordedGateProviderDecision(t, secondRoom)
	close := mustRecordedGate(t, binary, fixture.root,
		"gate", "record", "recorded-gate-task", "--room", secondRoom, "--workflow-dir", fixture.root)
	assertCommandOutput(t, close.stdout, "state=closed", "attempt=gate-attempt:recorded-gate-task-validation-2", "decision=approve")
	commitRecordedGateState(t, binary, fixture, "record replacement provider decision")
	consume := mustRecordedGate(t, binary, fixture.root,
		"gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
	assertCommandOutput(t, consume.stdout, "consumed=true", "target-stage=handoff")
	commitRecordedGateState(t, binary, fixture, "consume replacement authorization")

	doc, _, err := gates.Read(fixture.entity)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Records) != 1 || len(doc.Records[0].Attempts) != 2 {
		t.Fatalf("attempt history = %#v", doc.Records)
	}
	stale, current := doc.Records[0].Attempts[0], doc.Records[0].Attempts[1]
	if stale.Withdrawal == nil || stale.Withdrawal.By != "agent:first-officer" ||
		stale.Resolution != nil || stale.ProviderEvidence != nil || stale.Application != nil {
		t.Fatalf("withdrawn authority was not retained cleanly: %#v", stale)
	}
	if current.Withdrawal != nil || current.Resolution == nil || current.ProviderEvidence == nil ||
		current.Application == nil || current.Application.State != "consumed" {
		t.Fatalf("replacement did not exclusively own consumed authority: %#v", current)
	}
	if readFile(t, filepath.Join(firstRoom, "gate-briefing.json")) != firstBriefing ||
		readFile(t, filepath.Join(firstRoom, "request.json")) != firstRequest {
		t.Fatal("replacement lifecycle changed withdrawn room bytes")
	}
	if !strings.Contains(readFile(t, fixture.entity), "status: handoff") {
		t.Fatal("replacement consumption did not advance workflow status")
	}
	if !strings.Contains(git(t, fixture.stateRoot, "status", "--short"), "?? dirty-sibling.md") {
		t.Fatal("path-scoped lifecycle commits swept dirty sibling")
	}

	entityBeforeRefusal := readFile(t, fixture.entity)
	firstRequestPath := filepath.Join(firstRoom, "request.json")
	writeFile(t, firstRequestPath, strings.Replace(firstRequest, `"actor": "person:captain"`, `"actor": "agent:other"`, 1))
	refused := runRecordedGateCommand(binary, fixture.root, "", "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
	if refused.exit == 0 || !strings.Contains(refused.stderr, "frozen digest") {
		t.Fatalf("withdrawn retained-authority drift validated: exit=%d stderr=%q", refused.exit, refused.stderr)
	}
	if readFile(t, fixture.entity) != entityBeforeRefusal {
		t.Fatal("withdrawn retained-authority refusal changed entity")
	}
	writeFile(t, firstRequestPath, firstRequest)

	resultPath := filepath.Join(secondRoom, "provider", "result.json")
	resultBytes := readFile(t, resultPath)
	writeFile(t, resultPath, strings.Replace(resultBytes, `"decision": "approve"`, `"decision": "hold"`, 1))
	refused = runRecordedGateCommand(binary, fixture.root, "", "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
	if refused.exit == 0 || !strings.Contains(refused.stderr, "provider/result.json") {
		t.Fatalf("closed provider drift validated: exit=%d stderr=%q", refused.exit, refused.stderr)
	}
	if readFile(t, fixture.entity) != entityBeforeRefusal {
		t.Fatal("closed provider refusal changed entity")
	}
	writeFile(t, resultPath, resultBytes)
}

func writeRecordedGateProviderDecision(t *testing.T, room string) {
	t.Helper()
	var briefing struct {
		ID        string `json:"id"`
		Artifacts []struct {
			ID, URI, Rev string
		} `json:"artifacts"`
		Context []struct {
			ID, URI, Rev string
		} `json:"context"`
	}
	if err := json.Unmarshal([]byte(readFile(t, filepath.Join(room, "gate-briefing.json"))), &briefing); err != nil {
		t.Fatal(err)
	}
	items := make([]map[string]string, 0, len(briefing.Artifacts)+len(briefing.Context))
	for _, item := range briefing.Artifacts {
		items = append(items, map[string]string{"type": "Artifact", "id": item.ID, "uri": item.URI, "rev": item.Rev})
	}
	for _, item := range briefing.Context {
		items = append(items, map[string]string{"type": "Reference", "id": item.ID, "uri": item.URI, "rev": item.Rev})
	}
	if len(items) == 0 {
		t.Fatal("prepared Briefing has no presented items")
	}
	provider := filepath.Join(room, "provider")
	if err := os.Mkdir(provider, 0o755); err != nil {
		t.Fatal(err)
	}
	inventory, err := json.MarshalIndent(map[string]any{"items": items}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(provider, "presented-inventory.json"), string(inventory)+"\n")
	result, err := json.MarshalIndent(map[string]any{
		"type":        "review-v1-result",
		"briefing":    briefing.ID,
		"artifact":    map[string]string{"id": items[0]["id"], "uri": items[0]["uri"], "rev": items[0]["rev"]},
		"annotations": []any{},
		"resolution": map[string]string{
			"type": "Resolution", "id": "resolution:captain-replacement-2",
			"briefing": briefing.ID, "by": "person:captain",
			"at": "2026-07-26T12:00:00Z", "decision": "approve",
		},
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(provider, "result.json"), string(result)+"\n")
}

func assertRecordedGateDispatchRow(t *testing.T, binary string, fixture recordedGateFixture, current, next string) {
	t.Helper()
	for _, args := range [][]string{
		{"status", "--workflow-dir", fixture.root, "--next", "--json"},
		{"status", "--workflow-dir", fixture.root, "--boot", "--identify", "--json"},
	} {
		result := mustRecordedGate(t, binary, fixture.root, args...)
		var envelope struct {
			Dispatchable []struct {
				Current string `json:"current"`
				Next    string `json:"next"`
			} `json:"dispatchable"`
		}
		if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
			t.Fatalf("parse spacedock %v: %v\n%s", args, err, result.stdout)
		}
		if len(envelope.Dispatchable) != 1 ||
			envelope.Dispatchable[0].Current != current ||
			envelope.Dispatchable[0].Next != next {
			t.Fatalf("spacedock %v dispatchable=%+v, want current=%s next=%s",
				args, envelope.Dispatchable, current, next)
		}
	}
}

func TestRecordedGateLifecycleTerminalConsumeHasNoDispatchableSuccessor(t *testing.T) {
	binary, fixture := buildRecordedGateBinary(t), writeRecordedGateFixture(t)
	writeFile(t, filepath.Join(fixture.root, "README.md"), strings.Replace(readFile(t, filepath.Join(fixture.root, "README.md")), "    - name: handoff\n", "", 1))
	bindRecordedGate(t, binary, fixture)
	commitRecordedGateState(t, binary, fixture, "bind terminal gate package")
	closeRecordedGate(t, binary, fixture, "approve")
	commitRecordedGateState(t, binary, fixture, "record terminal gate decision")
	// Terminal-target consume routes without spending: the approval stays
	// pending for merge guard's delivery envelope, and the routed entity (still
	// at its gated stage) has no dispatchable successor — dispatch never
	// selects a gated stage as stage work.
	assertCommandOutput(t, mustRecordedGate(t, binary, fixture.root, "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root).stdout, "consumed=false", "target-stage=done", "route=approved-awaiting-merge")
	commitRecordedGateState(t, binary, fixture, "route terminal gate approval to the merge ceremony")
	assertCommandOutput(t, mustRecordedGate(t, binary, fixture.root, "status", "--workflow-dir", fixture.root, "--next", "--json").stdout, `"dispatchable":[]`)
}

func TestRecordedGateLifecycleEnteredStageRecoveryMatrix(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	complete := "\n## Stage Report: handoff\n\n- DONE: Complete the entered handoff.\n  Commit abc123 contains the durable handoff evidence.\n\n### Summary\n\nThe entered handoff is complete and ready for done.\n"
	cases := []struct {
		name   string
		report string
		dirty  bool
		next   string
	}{
		{"committed heading only", "\n## Stage Report: handoff\n", false, "handoff"},
		{"committed item without evidence", "\n## Stage Report: handoff\n\n- DONE: Complete the handoff.\n\n### Summary\n\nIncomplete evidence.\n", false, "handoff"},
		{"committed failed item", "\n## Stage Report: handoff\n\n- FAILED: Complete the handoff.\n  The handoff remains broken.\n\n### Summary\n\nIncomplete work.\n", false, "handoff"},
		{"committed wrong stage", strings.Replace(complete, "Stage Report: handoff", "Stage Report: implementation", 1), false, "handoff"},
		{"later malformed masks older valid", complete + "\n## Stage Report: handoff (cycle 2)\n\n- DONE: Later report has no evidence.\n\n### Summary\n\nLater but malformed.\n", false, "handoff"},
		{"dirty valid report", complete, true, "handoff"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := writeRecordedGateFixture(t)
			bindRecordedGate(t, binary, fixture)
			commitRecordedGateState(t, binary, fixture, "bind retained gate package")
			closeRecordedGate(t, binary, fixture, "approve")
			commitRecordedGateState(t, binary, fixture, "record delegated gate decision")
			mustRecordedGate(t, binary, fixture.root, "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
			commitRecordedGateState(t, binary, fixture, "consume gate authorization")

			writeFile(t, fixture.entity, readFile(t, fixture.entity)+tc.report)
			gitCommitPathScoped(t, fixture.stateRoot, "recorded-gate-task/index.md", "record handoff report")
			if tc.dirty {
				writeFile(t, fixture.entity, readFile(t, fixture.entity)+"\nUncommitted entity dirt.\n")
			}
			assertRecordedGateDispatchRow(t, binary, fixture, "handoff", tc.next)
		})
	}
}
func assertRecordedGateByteCleanFailure(t *testing.T, fixture recordedGateFixture, result recordedGateCommand, wants ...string) {
	if result.exit == 0 {
		t.Fatalf("refusal unexpectedly exited 0: stdout=%q stderr=%q", result.stdout, result.stderr)
	}
	output := result.stdout + result.stderr
	for _, want := range wants {
		if !strings.Contains(strings.ToLower(output), strings.ToLower(want)) {
			t.Errorf("refusal output missing actionable %q: %s", want, output)
		}
	}
	if _, err := os.Stat(fixture.entity + ".gates.lock"); !os.IsNotExist(err) {
		t.Fatalf("refusal left lock residue: %v", err)
	}
}
func bindRecordedGate(t *testing.T, binary string, fixture recordedGateFixture) {
	args := []string{
		"gate", "prepare", "recorded-gate-task",
		"--question", "Should the recorded validation gate advance?",
		"--artifact", fixture.gateReview,
		"--summary", "Exact recorded gate validation summary.",
	}
	for _, reference := range fixture.references {
		args = append(args, "--reference", reference)
	}
	args = append(args, "--workflow-dir", fixture.root)
	assertCommandOutput(t, mustRecordedGate(t, binary, fixture.root, args...).stdout, "state=open")
}
func closeRecordedGate(t *testing.T, binary string, fixture recordedGateFixture, decision string) {
	mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task",
		"--decision", decision, "--actor", "agent:first-officer", "--reason", "evidence-backed route",
		"--workflow-dir", fixture.root)
}
func TestRecordedGateLifecycleAC5RefusalMatrix(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	for _, tc := range []struct {
		name  string
		args  []string
		wants []string
	}{
		{"actor", []string{"--decision", "approve", "--reason", "evidence"}, []string{"actor"}},
		{"unsupported-actor", []string{"--decision", "approve", "--actor", "agent:ensign", "--reason", "evidence"}, []string{"actor"}},
		{"approve-missing-reason", []string{"--decision", "approve", "--actor", "agent:first-officer"}, []string{"reason"}},
		{"approve-whitespace-reason", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", " \t"}, []string{"reason"}},
		{"reason", []string{"--decision", "revise", "--actor", "agent:first-officer"}, []string{"reason"}},
		{"retired-exact-directive", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence", "--directive", recordedGateDirective}, []string{"unknown gate flag", "--directive"}},
		{"retired-altered-directive", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence", "--directive", strings.TrimSuffix(recordedGateDirective, ".")}, []string{"unknown gate flag", "--directive"}},
		{"retired-directive-file", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence", "--directive-file", "authority.txt"}, []string{"unknown gate flag", "--directive-file"}},
	} {
		t.Run("invalid-"+tc.name, func(t *testing.T) {
			fixture := writeRecordedGateFixture(t)
			bindRecordedGate(t, binary, fixture)
			before := treeDigest(t, fixture.stateRoot)
			args := append([]string{"gate", "record", "recorded-gate-task"}, tc.args...)
			args = append(args, "--workflow-dir", fixture.root)
			result := runRecordedGateCommand(binary, fixture.root, "", args...)
			assertRecordedGateByteCleanFailure(t, fixture, result, tc.wants...)
			wantExit := map[string]int{"actor": 2, "unsupported-actor": 1}[tc.name]
			if exact := map[string]string{"actor": "Error: --decision requires --actor ID\n", "unsupported-actor": "Error: unsupported chat decision actor \"agent:ensign\"\n"}[tc.name]; exact != "" && (result.exit != wantExit || result.stderr != exact) {
				t.Fatalf("invalid %s exit/stderr = %d/%q, want %d/%q", tc.name, result.exit, result.stderr, wantExit, exact)
			}
			if after := treeDigest(t, fixture.stateRoot); after != before {
				t.Fatalf("invalid %s changed workflow bytes", tc.name)
			}
		})
	}
	t.Run("validate-and-eligibility-reads", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		closeRecordedGate(t, binary, fixture, "approve")
		before := treeDigest(t, fixture.stateRoot)
		mustRecordedGate(t, binary, fixture.root, "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
		mustRecordedGate(t, binary, fixture.root, "gate", "eligibility", "recorded-gate-task", "--workflow-dir", fixture.root)
		if after := treeDigest(t, fixture.stateRoot); after != before {
			t.Fatal("validate/eligibility read changed workflow bytes")
		}
	})
	t.Run("forced-close-validation-mismatch", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		closeRecordedGate(t, binary, fixture, "approve")
		body := readFile(t, fixture.entity)
		writeFile(t, fixture.entity, strings.Replace(body, "decision: approve", "decision: hold", 1))
		before := treeDigest(t, fixture.stateRoot)
		result := runRecordedGateCommand(binary, fixture.root, "", "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, result, "application")
		if after := treeDigest(t, fixture.stateRoot); after != before {
			t.Fatal("close-validation mismatch changed workflow bytes")
		}
	})
	calls := []string{"redo with feedback", "reject with feedback-to", "reject without feedback-to", "hold", "not yet"}
	reasons := []string{"accepts-direction: add the retry test", "rejects-direction: replace the rejected cache design", "rejects-direction: name a feedback owner", "pause: wait for security sign-off", "pause: rerun the failing CI lane"}
	for i, decision := range []string{"revise", "revise", "hold", "hold", "hold"} {
		t.Run(calls[i]+"-consume", func(t *testing.T) {
			fixture := writeRecordedGateFixture(t)
			bindRecordedGate(t, binary, fixture)
			commitRecordedGateState(t, binary, fixture, "bind "+calls[i])
			mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task", "--decision", decision, "--actor", "agent:first-officer", "--reason", reasons[i], "--workflow-dir", fixture.root)
			closeCommit := commitRecordedGateState(t, binary, fixture, "durably record "+decision)
			closed, _, err := gates.Read(fixture.entity)
			attempt := closed.Records[0].Attempts[0]
			requireRecordedGate(t, err == nil && readFile(t, fixture.entity) == recordedGateEntityAt(t, fixture, closeCommit) && attempt.Resolution.Decision == decision && attempt.Resolution.Reason == reasons[i] && attempt.Application.Action == map[string]string{"revise": "feedback", "hold": "none"}[decision], "%s close/route snapshot mismatch", calls[i])
			before := treeDigest(t, fixture.stateRoot)
			result := runRecordedGateCommand(binary, fixture.root, "", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
			assertRecordedGateByteCleanFailure(t, fixture, result, "condition")
			if after := treeDigest(t, fixture.stateRoot); after != before {
				t.Fatalf("%s consume refusal changed workflow bytes", decision)
			}
			if snapshot := recordedGateEntityAt(t, fixture, closeCommit); !strings.Contains(snapshot, "decision: "+decision) {
				t.Fatalf("%s close was not durable before route refusal", decision)
			}
		})
	}
	t.Run("blocked", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		closeRecordedGate(t, binary, fixture, "approve")
		closeCommit := commitRecordedGateState(t, binary, fixture, "durably record blocked approval")
		closeSnapshot := recordedGateEntityAt(t, fixture, closeCommit)
		if !strings.Contains(closeSnapshot, "decision: approve") || !strings.Contains(closeSnapshot, "state: pending") {
			t.Fatalf("blocked approval close was not durable before consume:\n%s", closeSnapshot)
		}
		body := readFile(t, fixture.entity)
		writeFile(t, fixture.entity, strings.Replace(body, "                blockers: []",
			"                blockers:\n                    - id: blocker:external\n                      state: unsatisfied", 1))
		before := treeDigest(t, fixture.stateRoot)
		result := runRecordedGateCommand(binary, fixture.root, "", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, result, "blocked")
		if after := treeDigest(t, fixture.stateRoot); after != before {
			t.Fatal("blocked consume changed workflow bytes")
		}
	})
	t.Run("repeat-consume", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		closeRecordedGate(t, binary, fixture, "approve")
		mustRecordedGate(t, binary, fixture.root, "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		before := treeDigest(t, fixture.stateRoot)
		result := runRecordedGateCommand(binary, fixture.root, "", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, result, "consumed")
		if after := treeDigest(t, fixture.stateRoot); after != before {
			t.Fatal("repeat consume changed workflow bytes")
		}
	})
}
func TestRecordedGateLifecycleAC7ResumeMatrix(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	t.Run("open-same-package", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		before := recordedGateTreeSnapshot(t, fixture.stateRoot)
		bindRecordedGate(t, binary, fixture)
		assertRecordedGateTreeSnapshot(t, fixture.stateRoot, before)
		if after := readFile(t, fixture.entity); strings.Count(after, "gate-attempt:recorded-gate-task-validation-1") != 1 {
			t.Fatal("same-package open resume minted an attempt")
		}
	})
	for _, decision := range []string{"revise", "hold"} {
		t.Run("closed-"+decision, func(t *testing.T) {
			fixture := writeRecordedGateFixture(t)
			bindRecordedGate(t, binary, fixture)
			closeRecordedGate(t, binary, fixture, decision)
			before := recordedGateTreeSnapshot(t, fixture.stateRoot)
			for pass := 0; pass < 3; pass++ {
				eligibility := mustRecordedGate(t, binary, fixture.root, "gate", "eligibility", "recorded-gate-task", "--workflow-dir", fixture.root)
				assertCommandOutput(t, eligibility.stdout, "eligible=false")
			}
			after := readFile(t, fixture.entity)
			assertRecordedGateTreeSnapshot(t, fixture.stateRoot, before)
			if strings.Count(after, "resolution:spacedock") != 1 ||
				strings.Contains(after, "state: consumed") {
				t.Fatalf("%s resume changed bytes, duplicated resolution, or consumed", decision)
			}
		})
	}
	t.Run("approval-close-commit-consume", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		commitRecordedGateState(t, binary, fixture, "bind retained gate package")

		// Fresh process 1 closes the gate and stops before its required state commit.
		closeRecordedGate(t, binary, fixture, "approve")
		closedUncommitted := recordedGateTreeSnapshot(t, fixture.stateRoot)
		entityRel := strings.TrimPrefix(fixture.entity, fixture.stateRoot+string(os.PathSeparator))
		if exec.Command("git", "-C", fixture.stateRoot, "diff", "--quiet", "--", entityRel).Run() == nil {
			t.Fatal("successful close was already committed")
		}
		if commits := strings.Fields(git(t, fixture.stateRoot, "log", "--format=%H", "-Sdecision: approve", "--", entityRel)); len(commits) != 0 {
			t.Fatalf("uncommitted close already has %d decision commits", len(commits))
		}
		repeatClose := runRecordedGateCommand(binary, fixture.root, "", "gate", "record", "recorded-gate-task",
			"--decision", "approve", "--actor", "agent:first-officer", "--reason", "duplicate",
			"--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, repeatClose, "closed")
		assertRecordedGateTreeSnapshot(t, fixture.stateRoot, closedUncommitted)

		// Fresh process 2 resumes the uncommitted close and commits the exact pending state.
		closeCommit := commitRecordedGateState(t, binary, fixture, "record delegated gate decision")
		committedPending := recordedGateTreeSnapshot(t, fixture.stateRoot)
		closeCommits := strings.Fields(git(t, fixture.stateRoot, "log", "--format=%H", "-Sdecision: approve", "--", entityRel))
		if len(closeCommits) != 1 || closeCommits[0] != closeCommit {
			t.Fatalf("close commits=%v, want exactly %s", closeCommits, closeCommit)
		}
		if parent := recordedGateEntityAt(t, fixture, closeCommit+"^"); strings.Contains(parent, "resolution:spacedock") {
			t.Fatal("close commit parent already contains a Resolution")
		}
		repeatClose = runRecordedGateCommand(binary, fixture.root, "", "gate", "record", "recorded-gate-task",
			"--decision", "approve", "--actor", "agent:first-officer", "--reason", "duplicate",
			"--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, repeatClose, "closed")
		assertRecordedGateTreeSnapshot(t, fixture.stateRoot, committedPending)

		// Fresh process 3 resumes the committed pending approval and consumes it once.
		consume := mustRecordedGate(t, binary, fixture.root, "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		assertCommandOutput(t, consume.stdout, "consumed=true", "target-stage=handoff")
		consumedCommit := commitRecordedGateState(t, binary, fixture, "consume gate authorization")
		consumed := recordedGateTreeSnapshot(t, fixture.stateRoot)
		if !recordedGateCommittedBeforeDispatch(t, fixture, closeCommit, consumedCommit, consumedCommit) {
			t.Fatal("consumed commit is not a descendant of the exact close commit before dispatch")
		}
		for _, pickaxe := range []string{"state: consumed", "status: handoff"} {
			if commits := strings.Fields(git(t, fixture.stateRoot, "log", "--format=%H", "-S"+pickaxe, "--", entityRel)); len(commits) != 1 || commits[0] != consumedCommit {
				t.Fatalf("%s commits=%v, want exactly %s", pickaxe, commits, consumedCommit)
			}
		}
		repeatConsume := runRecordedGateCommand(binary, fixture.root, "", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, repeatConsume, "consumed")
		assertRecordedGateTreeSnapshot(t, fixture.stateRoot, consumed)
		after := readFile(t, fixture.entity)
		if strings.Count(after, "resolution:spacedock") != 1 ||
			strings.Count(after, "state: consumed") != 1 ||
			strings.Count(after, "status: handoff") != 1 {
			t.Fatal("resume duplicated decision, consume, or transition")
		}
	})
}
func TestRecordedGateLifecycleProvenanceMutants(t *testing.T) {
	valid := recordedGateObservation{
		events: append([]string(nil), recordedGateRequiredEvents...),
		before: "status: validation",
		after: "status: handoff\nid: gate:recorded-gate-task:validation\nstage: validation\nid: gate-attempt:recorded-gate-task-validation-1\n" +
			"id: " + recordedGateBriefingID + "\ndigest: " + recordedGateDigest + "\n" +
			"id: resolution:spacedock:recorded-gate-task:validation:1\nbriefing: " + recordedGateBriefingID + "\n" +
			"by: agent:first-officer\n                decision: approve\n                reason: " + recordedGateReason + "\n" +
			"target-stage: handoff\n                state: consumed\n## Stage Report: handoff\n\n- DONE: Successor dispatch followed decision: approve.",
		dispatch:     recordedGateDispatchProof{builds: 1, successfulBuilds: 1, durableEffects: 1, ordered: true, committed: true},
		expectedNext: "handoff",
	}
	if err := assertRecordedGateLifecycle(valid); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	fenced := valid
	fenced.after = "---\n" + strings.Replace(valid.after, "\n## Stage Report: handoff", "\n---\n## Stage Report: handoff", 1)
	fenced.after = strings.Replace(fenced.after, "- DONE: Successor dispatch", "- DONE: target-stage: handoff successor dispatch", 1)
	if err := assertRecordedGateLifecycle(fenced); err != nil {
		t.Fatalf("frontmatter-scoped authority: %v", err)
	}
	for name, mutate := range map[string]func(*recordedGateObservation){
		"actor-swap": func(o *recordedGateObservation) {
			o.after = strings.Replace(o.after, "by: agent:first-officer", "by: person:captain", 1)
		},
		"blank-reason":         func(o *recordedGateObservation) { o.after = strings.Replace(o.after, recordedGateReason, "", 1) },
		"forged-adoption-note": func(o *recordedGateObservation) { o.after = "adoption-note: forged\n" + o.after },
		"heading-deleted": func(o *recordedGateObservation) {
			o.after = strings.ReplaceAll(o.after, "## Stage Report: handoff", "")
		},
		"mutated-handoff-done": func(o *recordedGateObservation) { o.after = strings.Replace(o.after, "- DONE:", "- FAILED:", 1) },
	} {
		t.Run(name, func(t *testing.T) {
			mutant := valid
			mutate(&mutant)
			if err := assertRecordedGateLifecycle(mutant); err == nil {
				t.Fatal("mutant graded PASS")
			}
		})
	}
}
func authorizeRecordedGateDispatch(events []string, entity, successor string) error {
	if len(events) != len(recordedGateRequiredEvents) {
		return fmt.Errorf("trace incomplete: got %v", events)
	}
	for i, want := range recordedGateRequiredEvents {
		if events[i] != want {
			return fmt.Errorf("trace event %d = %q, want %q", i+1, events[i], want)
		}
	}
	if !strings.Contains(entity, "status: "+successor) || !strings.Contains(entity, "state: consumed") {
		return fmt.Errorf("authorization was not atomically consumed into %s", successor)
	}
	return nil
}

func recordedGatePrompt(workflowRoot string) string {
	return fmt.Sprintf("Use $spacedock:first-officer for this whole run.\n\nWorkflow directory: %s\nEngage only `recorded-gate-task` under this delegated conn: %s\nDo not pass `--host` to `dispatch build`.\nUse the already-committed `.spacedock-state/recorded-gate-task/selected/gate-review.md` Artifact and both already-committed References, `.spacedock-state/recorded-gate-task/selected/entity-snapshot.md` and `recorder-contract.md`; do not replace or regenerate them. Prepare the recorder-ready validation room, approve it, and continue until the handoff worker records %s in durable state, then stop.", workflowRoot, recordedGateDirective, recordedGateDispatchMarker)
}

func writeRecordedGateLoggingShim(t *testing.T, binary, logPath string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	shim := filepath.Join(dir, "spacedock")
	stateRoot := filepath.Join(filepath.Dir(logPath), ".spacedock-state")
	if filepath.Base(filepath.Dir(logPath)) == "evidence" {
		stateRoot = filepath.Join(filepath.Dir(filepath.Dir(logPath)), ".spacedock-state")
	}
	script := fmt.Sprintf("#!/bin/sh\nprintf 'begin\\t%%s\\n' \"$*\" >> %q\n[ \"$1 $2\" = \"dispatch build\" ] && git -C %q rev-parse HEAD | sed 's/^/dispatch-head\\t/' >> %q\n%q \"$@\"\ncode=$?\nprintf 'exit=%%s\\t%%s\\n' \"$code\" \"$*\" >> %q\n[ \"$code\" -eq 0 ] && [ \"$1 $2\" = \"state commit\" ] && git -C %q rev-parse HEAD | sed 's/^/state-head\\t/' | tee -a %q\nexit \"$code\"\n", logPath, stateRoot, logPath, binary, logPath, stateRoot, logPath)
	writeFile(t, shim, script)
	if err := os.Chmod(shim, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func withRecordedGateEnv(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	for _, item := range env {
		if !strings.HasPrefix(item, prefix) {
			out = append(out, item)
		}
	}
	return append(out, prefix+value)
}

func TestSpacedockShimShellEnvOverridesLauncherPin(t *testing.T) {
	dir, env := t.TempDir(), []string{"SPACEDOCK_BIN=/stale/spacedock"}
	env = withSpacedockShimShellEnv(t, env, dir)
	env = withRecordedGateEnv(env, "SPACEDOCK_BIN", "/real/spacedock")
	for _, shell := range []string{"/bin/bash", "zsh"} {
		t.Run(filepath.Base(shell), func(t *testing.T) {
			if shell == "zsh" {
				var err error
				if shell, err = exec.LookPath("zsh"); err != nil {
					t.Skip("zsh unavailable")
				}
			}
			cmd := exec.Command(shell, "-c", `printf %s "$SPACEDOCK_BIN"`)
			cmd.Env = env
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("%s startup: %v\n%s", shell, err, output)
			}
			if value, want := string(output), filepath.Join(dir, "spacedock"); value != want {
				t.Fatalf("%s observed SPACEDOCK_BIN=%q, want shim %q", shell, value, want)
			}
		})
	}
}

func withSpacedockShimShellEnv(t *testing.T, env []string, shimDir string) []string {
	t.Helper()
	shellEnvDir := t.TempDir()
	bashEnv := filepath.Join(shellEnvDir, "recorded-gate-env.sh")
	writeFile(t, bashEnv, "export SPACEDOCK_BIN="+filepath.Join(shimDir, "spacedock")+"\n")
	writeFile(t, filepath.Join(shellEnvDir, ".zshenv"), readFile(t, bashEnv))
	env = withRecordedGateEnv(env, "BASH_ENV", bashEnv)
	return withRecordedGateEnv(env, "ZDOTDIR", shellEnvDir)
}

func recordedGateEventsFromCommandLog(log string) []string {
	var events []string
	started := false
	for _, line := range strings.Split(log, "\n") {
		if !strings.HasPrefix(line, "exit=0\tgate ") || strings.Contains(line, " --help") {
			continue
		}
		switch {
		case strings.Contains(line, "gate prepare "):
			started = true
			events = append(events, "prepare")
		case started && strings.Contains(line, "gate record ") && (strings.Contains(line, " --decision ") || strings.Contains(line, " --room ")):
			events = append(events, "decision-record")
		case started && strings.Contains(line, "gate consume "):
			events = append(events, "consume")
		}
	}
	return events
}

func TestRecordedGateLifecycleMissingEventControls(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	for skip, omitted := range recordedGateRequiredEvents {
		t.Run(omitted, func(t *testing.T) {
			fixture := writePreparedRecordedGateFixture(t)
			steps := [][]string{
				{"gate", "prepare", "recorded-gate-task", "--question", "Advance?", "--artifact", fixture.gateReview, "--summary", "Exact summary.", "--reference", fixture.references[0], "--workflow-dir", fixture.root},
				{"gate", "record", "recorded-gate-task", "--decision", "approve", "--actor", "agent:first-officer", "--reason", recordedGateReason, "--workflow-dir", fixture.root},
				{"gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root},
			}
			var commands []recordedGateCommand
			for i, args := range steps {
				if i == skip {
					continue
				}
				commands = append(commands, runRecordedGateCommand(binary, fixture.root, recordedGateRequiredEvents[i], args...))
			}
			if err := authorizeRecordedGateDispatch(successfulRecordedGateEvents(commands), readFile(t, fixture.entity), "handoff"); err == nil {
				t.Fatalf("real command replay without %s authorized dispatch", omitted)
			}
		})
	}
}

func recordedGateLiveObservation(t *testing.T, fixture recordedGateFixture, before, commandLog string) recordedGateObservation {
	t.Helper()
	log := readFile(t, commandLog)
	builds, successfulBuilds, consumed, ordered := 0, 0, false, true
	prepareAt := strings.Index(log, "exit=0\tgate prepare ")
	bindCommitAt := strings.Index(log, "exit=0\tstate commit recorded-gate-task")
	decisionAt := strings.Index(log, "exit=0\tgate record ")
	ordered = prepareAt >= 0 && bindCommitAt > prepareAt && decisionAt > bindCommitAt
	dispatchHead, buildCommand := "", ""
	for _, line := range strings.Split(log, "\n") {
		if strings.HasPrefix(line, "exit=0\tgate consume ") {
			consumed = true
		}
		if strings.HasPrefix(line, "begin\tdispatch build ") && !strings.Contains(line, " --help") {
			builds++
			buildCommand = strings.TrimPrefix(line, "begin\t")
			ordered = ordered && consumed
		}
		if strings.HasPrefix(line, "exit=0\tdispatch build ") && strings.TrimPrefix(line, "exit=0\t") == buildCommand {
			successfulBuilds++
		}
		if strings.HasPrefix(line, "dispatch-head\t") {
			dispatchHead = strings.TrimPrefix(line, "dispatch-head\t")
		}
	}
	after := resolveRecordedGateEntity(fixture)
	entityRel, err := filepath.Rel(fixture.stateRoot, fixture.entity)
	if err != nil {
		t.Fatal(err)
	}
	commits := strings.Fields(git(t, fixture.stateRoot, "log", "--format=%H", "-S"+recordedGateDispatchMarker, "--", entityRel))
	effects := 0
	if strings.Contains(after, recordedGateDispatchMarker) {
		effects = len(commits)
	}
	closeCommit := strings.SplitN(strings.TrimSpace(git(t, fixture.stateRoot, "log", "--reverse", "--format=%H", "-Sid: resolution:spacedock:recorded-gate-task:validation:1", "--", entityRel)), "\n", 2)[0]
	consumedCommit := strings.SplitN(strings.TrimSpace(git(t, fixture.stateRoot, "log", "--reverse", "--format=%H", "-S\n                state: consumed", "--", entityRel)), "\n", 2)[0]
	gateID := firstRecordedGateMatch(after, `(?m)^\s+gate: (gate:[^\s]+)$`)
	attemptID := firstRecordedGateMatch(after, `(?m)^\s+- id: (gate-attempt:[^\s]+)$`)
	briefingID := firstRecordedGateMatch(after, `(?m)^\s+id: (briefing:[^\s]+)$`)
	digest := firstRecordedGateMatch(after, `(?m)^\s+digest: (sha256:[0-9a-f]{64})$`)
	resolutionID := firstRecordedGateMatch(after, `(?m)^\s+id: (resolution:[^\s]+)$`)
	return recordedGateObservation{
		events: recordedGateEventsFromCommandLog(log), before: before, after: after,
		dispatch: recordedGateDispatchProof{
			builds: builds, successfulBuilds: successfulBuilds, durableEffects: effects, ordered: ordered,
			committed: recordedGateCommittedBeforeDispatch(t, fixture, closeCommit, consumedCommit, dispatchHead, strings.Join(commits, " ")),
		},
		expectedNext: "handoff",
		gateID:       gateID, attemptID: attemptID, briefingID: briefingID,
		digest: digest, resolutionID: resolutionID,
	}
}

func firstRecordedGateMatch(body, pattern string) string {
	match := regexp.MustCompile(pattern).FindStringSubmatch(body)
	if len(match) != 2 {
		return ""
	}
	return match[1]
}

func resolveRecordedGateEntity(fixture recordedGateFixture) string {
	for _, path := range []string{
		fixture.entity,
		filepath.Join(fixture.stateRoot, "_archive", "recorded-gate-task", "index.md"),
		filepath.Join(fixture.stateRoot, "_archive", "recorded-gate-task.md"),
	} {
		if body, err := os.ReadFile(path); err == nil {
			return string(body)
		}
	}
	return ""
}

func writeRecordedGateFixture(t *testing.T) recordedGateFixture {
	t.Helper()
	return writeRecordedGateFixtureAt(t, t.TempDir())
}

func writePreparedRecordedGateFixture(t *testing.T) recordedGateFixture {
	t.Helper()
	return writePreparedRecordedGateFixtureAt(t, t.TempDir())
}

func writePreparedRecordedGateFixtureAt(t *testing.T, root string) recordedGateFixture {
	t.Helper()
	stateRoot := filepath.Join(root, ".spacedock-state")
	writeFile(t, filepath.Join(root, "README.md"), recordedGateReadme())
	mainReference := filepath.Join(root, "recorder-contract.md")
	writeFile(t, mainReference, "# Recorder contract\n\nPrepare one provider-neutral room from exact local Git objects.\n")
	gitInit(t, root)
	git(t, root, "config", "user.name", "Spacedock Test")
	git(t, root, "config", "user.email", "spacedock@example.invalid")

	entity := filepath.Join(stateRoot, "recorded-gate-task", "index.md")
	writeFile(t, entity, recordedGateEntity())
	gateReview := filepath.Join(filepath.Dir(entity), "selected", "gate-review.md")
	writeFile(t, gateReview, recordedGateSourceReview())
	entityReference := filepath.Join(filepath.Dir(entity), "selected", "entity-snapshot.md")
	writeFile(t, entityReference, "# Entity snapshot\n\nThe validation Stage Report is complete and ready for a decision.\n")
	gitInit(t, stateRoot)
	git(t, stateRoot, "config", "user.name", "Spacedock Test")
	git(t, stateRoot, "config", "user.email", "spacedock@example.invalid")
	git(t, stateRoot, "branch", "-M", "spacedock-state/"+filepath.Base(root))
	writeFile(t, filepath.Join(stateRoot, "dirty-sibling.md"), "unrelated concurrent dirt\n")
	return recordedGateFixture{
		root: root, stateRoot: stateRoot, entity: entity, gateReview: gateReview,
		references: []string{entityReference, mainReference},
	}
}

func writeRecordedGateFixtureAt(t *testing.T, root string) recordedGateFixture {
	t.Helper()
	return writePreparedRecordedGateFixtureAt(t, root)
}

func recordedGateReadme() string {
	return "---\n" +
		"commissioned-by: spacedock@1\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"state: .spacedock-state\n" +
		"stages:\n" +
		"  defaults: {worktree: false, concurrency: 1}\n" +
		"  states:\n" +
		"    - name: implementation\n      initial: true\n" +
		"    - name: validation\n      gate: true\n      feedback-to: implementation\n" +
		"    - name: handoff\n" +
		"    - name: done\n      terminal: true\n" +
		"---\n# Recorded Gate Lifecycle Fixture\n\n" +
		"### validation\n\nValidate and present the retained package.\n\n" +
		"### handoff\n\nAppend the exact marker `" + recordedGateDispatchMarker + "` and a `## Stage Report: handoff` with one DONE item, then return completion. Do not advance or archive the entity.\n\n- **Outputs:** The marker and handoff stage report.\n"
}

func recordedGateEntity() string {
	return "---\n" +
		"title: Recorded Gate Task\n" +
		"status: validation\n" +
		"completed:\nverdict:\nworktree:\n" +
		"---\n# Recorded Gate Task\n\n" +
		"## Acceptance criteria\n\n**AC-1** Successor dispatch requires consumed approval.\n\n" +
		"## Stage Report: validation\n\n- DONE: Replayed retained evidence\n  The real command fixture is green.\n\n" +
		"### Summary\n\nReady for the recorded decision gate.\n"
}

func recordedGateSourceReview() string {
	return "# Recorded Gate Task — validation review\n\n" +
		"Capability/change: provider-neutral preparation retains committed source identities without copying payloads.\n\n" +
		"Test and evidence: fresh-binary command replay, byte comparisons, and skipped-step mutants pass.\n\n" +
		"Findings: no material findings.\n\n" +
		"Recommendation: approve and consume the authorization once.\n"
}

func buildRecordedGateBinary(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "spacedock")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/spacedock")
	cmd.Dir = recordedGateRepoRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build fresh task-local spacedock: %v\n%s", err, out)
	}
	return binary
}

func runRecordedGateCommand(binary, cwd, event string, args ...string) recordedGateCommand {
	cmd := exec.Command(binary, args...)
	cmd.Dir = cwd
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exit := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exit = ee.ExitCode()
		} else {
			exit = -1
		}
	}
	return recordedGateCommand{event: event, argv: append([]string{binary}, args...), exit: exit, stdout: stdout.String(), stderr: stderr.String()}
}

func mustRecordedGate(t *testing.T, binary, cwd string, args ...string) recordedGateCommand {
	t.Helper()
	result := runRecordedGateCommand(binary, cwd, "", args...)
	if result.exit != 0 {
		t.Fatalf("spacedock %v exit=%d stdout=%q stderr=%q", args, result.exit, result.stdout, result.stderr)
	}
	return result
}

func commitRecordedGateState(t *testing.T, binary string, fixture recordedGateFixture, message string) string {
	t.Helper()
	result := runRecordedGateCommand(binary, fixture.root, "state-commit", "state", "commit", "recorded-gate-task", "--workflow-dir", fixture.root, "-m", message)
	if result.exit != 0 {
		t.Fatalf("state commit exit=%d stdout=%q stderr=%q", result.exit, result.stdout, result.stderr)
	}
	return strings.TrimSpace(git(t, fixture.stateRoot, "rev-parse", "HEAD"))
}

func recordedGateEntityAt(t *testing.T, fixture recordedGateFixture, commit string) string {
	t.Helper()
	if commit == "" {
		return ""
	}
	rel, _ := filepath.Rel(fixture.stateRoot, fixture.entity)
	return git(t, fixture.stateRoot, "show", commit+":"+filepath.ToSlash(rel))
}

func recordedGateCommittedBeforeDispatch(t *testing.T, fixture recordedGateFixture, close, consumed, dispatchHead string, effects ...string) bool {
	t.Helper()
	opened, closed, spent := recordedGateEntityAt(t, fixture, strings.TrimPrefix(close+"^", "^")), recordedGateEntityAt(t, fixture, close), recordedGateEntityAt(t, fixture, consumed)
	if close == "" || consumed == "" || dispatchHead == "" || close == consumed ||
		!strings.Contains(opened, "digest: sha256:") || strings.Contains(opened, "resolution:") ||
		!strings.Contains(closed, "decision: approve") || !strings.Contains(closed, "state: pending") ||
		!strings.Contains(spent, "status: handoff") || !strings.Contains(spent, "state: consumed") {
		return false
	}
	ordered := exec.Command("git", "-C", fixture.stateRoot, "merge-base", "--is-ancestor", close, consumed).Run() == nil &&
		exec.Command("git", "-C", fixture.stateRoot, "merge-base", "--is-ancestor", consumed, dispatchHead).Run() == nil
	return ordered && (len(effects) == 0 || strings.Count(effects[0], " ") == 0 &&
		exec.Command("git", "-C", fixture.stateRoot, "merge-base", "--is-ancestor", dispatchHead, effects[0]).Run() == nil)
}

func successfulRecordedGateEvents(commands []recordedGateCommand) []string {
	var events []string
	for _, command := range commands {
		if command.exit == 0 {
			events = append(events, command.event)
		}
	}
	return events
}

func assertCommandOutput(t *testing.T, output string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(output, want) {
			t.Fatalf("command output missing %q: %s", want, output)
		}
	}
}

func outputValue(output, key string) string {
	prefix := key + "="
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimPrefix(line, prefix)
		}
	}
	return ""
}

func treeDigest(t *testing.T, root string) string {
	t.Helper()
	h := sha256.New()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || strings.Contains(path, string(filepath.Separator)+".git"+string(filepath.Separator)) {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		h.Write([]byte(path))
		h.Write(body)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func canonicalRecordedGateDigest(t *testing.T, body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		t.Fatal(err)
	}
	canonical, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:])
}
func recordedGateTreeSnapshot(t *testing.T, root string) map[string]string {
	snapshot := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot[strings.TrimPrefix(path, root+string(os.PathSeparator))] = string(body)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}
func assertRecordedGateTreeSnapshot(t *testing.T, root string, expected map[string]string) {
	actual := recordedGateTreeSnapshot(t, root)
	if len(actual) != len(expected) {
		t.Fatalf("workflow tree file count=%d, want %d\nactual=%v\nexpected=%v", len(actual), len(expected), actual, expected)
	}
	for path, want := range expected {
		if got, ok := actual[path]; !ok || got != want {
			t.Fatalf("workflow tree mismatch at %s (present=%t)\n--- expected ---\n%s\n--- actual ---\n%s", path, ok, want, got)
		}
	}
}

func writeRecordedGateEvidence(t *testing.T, root string, commands []recordedGateCommand, before, after string, dispatches int) {
	t.Helper()
	dir := filepath.Join(root, "evidence")
	var log strings.Builder
	for _, command := range commands {
		fmt.Fprintf(&log, "event=%s exit=%d argv=%q\nstdout=%s\nstderr=%s\n", command.event, command.exit, command.argv[1:], command.stdout, command.stderr)
	}
	writeFile(t, filepath.Join(dir, "command.log"), log.String())
	writeFile(t, filepath.Join(dir, "entity.before.md"), before)
	writeFile(t, filepath.Join(dir, "entity.after.md"), after)
	writeFile(t, filepath.Join(dir, "consume.txt"), commands[len(commands)-1].stdout)
	writeFile(t, filepath.Join(dir, "dispatch.txt"), fmt.Sprintf("successor-dispatches=%d\n", dispatches))
}

func recordedGateRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join(wd, "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}
