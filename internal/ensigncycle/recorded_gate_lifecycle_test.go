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
	"strings"
	"testing"
)

var recordedGateRequiredEvents = []string{
	"briefing-record",
	"open-validate",
	"decision-record",
	"closed-validate",
	"eligibility",
	"consume",
}

const recordedGateDispatchMarker = "RECORDED-GATE-SUCCESSOR-DISPATCHED"

const (
	recordedGateBriefingID = "briefing:docs-dev:3k:validation:attempt-1:revision-1"
	recordedGateDigest     = "sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac"
	recordedGateReason     = "Captain directive: approved after reviewing the presented 3k validation gate."
	recordedGateDirective  = "you have the conn toward the sprint goal; authorized to approve gates, PR, relevant CI lanes, and merge; use your judgement."
)

type recordedGateDispatchProof struct {
	spawned      bool
	handle       string
	workerOutput bool
}

type recordedGateObservation struct {
	events       []string
	before       string
	after        string
	dispatch     recordedGateDispatchProof
	gateReview   string
	expectedNext string
}

func assertRecordedGateLifecycle(o recordedGateObservation) error {
	if len(o.events) != len(recordedGateRequiredEvents) {
		return fmt.Errorf("gate lifecycle recorded %d events, want %d before dispatch", len(o.events), len(recordedGateRequiredEvents))
	}
	for i, want := range recordedGateRequiredEvents {
		if got := o.events[i]; got != want {
			return fmt.Errorf("gate lifecycle event %d = %q, want %q", i+1, got, want)
		}
	}
	if !o.dispatch.spawned {
		return fmt.Errorf("no successful runtime successor spawn was observed")
	}
	if strings.TrimSpace(o.dispatch.handle) == "" {
		return fmt.Errorf("runtime successor spawn exposed no worker handle")
	}
	if !o.dispatch.workerOutput {
		return fmt.Errorf("spawned worker produced no correlated durable output")
	}
	if strings.Contains(o.before, recordedGateDispatchMarker) {
		return fmt.Errorf("successor marker already existed before the lifecycle began")
	}
	for _, want := range []string{"status: " + o.expectedNext, "state: consumed", "by: agent:first-officer"} {
		if !strings.Contains(o.after, want) {
			return fmt.Errorf("durable post-state missing %q", want)
		}
	}
	for _, exact := range []struct {
		label string
		value string
		count int
	}{
		{"gate identity", "gate: gate:docs-dev:3k:validation", 1},
		{"attempt identity", "id: gate-attempt:3k-validation-1", 1},
		{"briefing identity", "id: " + recordedGateBriefingID, 1},
		{"briefing resolution link", "briefing: " + recordedGateBriefingID, 1},
		{"briefing digest", "digest: " + recordedGateDigest, 1},
		{"resolution identity", "id: resolution:spacedock:docs-dev:3k:validation:1", 1},
		{"approval decision", "decision: approve", 1},
		{"approval actor", "by: agent:first-officer", 1},
		{"approval reason", recordedGateReason, 1},
		{"delegated directive", recordedGateDirective, 1},
		{"application target", "target-stage: " + o.expectedNext, 1},
		{"consumed application", "state: consumed", 1},
	} {
		if got := strings.Count(o.after, exact.value); got != exact.count {
			return fmt.Errorf("durable post-state %s count = %d, want %d for %q", exact.label, got, exact.count, exact.value)
		}
	}
	if o.before == o.after {
		return fmt.Errorf("gate lifecycle left entity byte-identical")
	}
	if err := assertConciseRecordedGateReview(o.gateReview); err != nil {
		return err
	}
	return nil
}

func assertConciseRecordedGateReview(review string) error {
	lower := strings.ToLower(review)
	for _, want := range []string{"capability", "evidence", "reviewed snapshot", "findings", "recommend", "decision"} {
		if !strings.Contains(lower, want) {
			return fmt.Errorf("gate review missing %q", want)
		}
	}
	for _, forbidden := range []string{"gates:\n", `"type":"briefing"`, "rebind ceremony", "supersede ceremony"} {
		if strings.Contains(lower, forbidden) {
			return fmt.Errorf("gate review leads with recorder mechanics %q", forbidden)
		}
	}
	return nil
}

func TestRecordedGateLifecycleContractDeclaresLoadBearingSequence(t *testing.T) {
	root := recordedGateRepoRoot(t)
	body := readFile(t, filepath.Join(root, "skills", "first-officer", "references", "first-officer-shared-core.md"))
	positions := []int{
		strings.Index(body, "gate record ENTITY --briefing"),
		strings.Index(body, "gate validate ENTITY"),
		strings.Index(body, "gate record ENTITY --decision"),
		strings.LastIndex(body, "gate validate ENTITY"),
		strings.Index(body, "gate eligibility ENTITY"),
		strings.Index(body, "gate consume ENTITY"),
	}
	for i, pos := range positions {
		if pos < 0 {
			t.Fatalf("shipped FO contract is missing %s", recordedGateRequiredEvents[i])
		}
		if i > 0 && pos <= positions[i-1] {
			t.Fatalf("shipped FO contract orders %s before %s", recordedGateRequiredEvents[i], recordedGateRequiredEvents[i-1])
		}
	}
}

func TestRecordedGateLifecycleSkippedStepMutantsRefuseDispatch(t *testing.T) {
	base := recordedGateObservation{
		events:       append([]string(nil), recordedGateRequiredEvents...),
		before:       "status: validation",
		after:        "status: done\nstate: consumed\nby: agent:first-officer\nreason: evidence\nadoption-note: captain grant",
		dispatch:     recordedGateDispatchProof{spawned: true, handle: "fixture-worker-1", workerOutput: true},
		expectedNext: "done",
		gateReview:   "Capability: recorder. Evidence: tests. Reviewed snapshot: abc. Findings: none. Recommend approve. Decision: approve?",
	}
	base.after = "status: done\ngate: gate:docs-dev:3k:validation\nid: gate-attempt:3k-validation-1\nid: " + recordedGateBriefingID + "\ndigest: " + recordedGateDigest + "\nid: resolution:spacedock:docs-dev:3k:validation:1\nbriefing: " + recordedGateBriefingID + "\nby: agent:first-officer\ndecision: approve\nreason: " + recordedGateReason + "\nadoption-note: '" + recordedGateDirective + "'\ntarget-stage: done\nstate: consumed"
	if err := assertRecordedGateLifecycle(base); err != nil {
		t.Fatalf("positive control: %v", err)
	}
	for skip := range recordedGateRequiredEvents {
		t.Run(recordedGateRequiredEvents[skip], func(t *testing.T) {
			mutant := base
			mutant.events = append([]string(nil), base.events[:skip]...)
			mutant.events = append(mutant.events, base.events[skip+1:]...)
			if err := assertRecordedGateLifecycle(mutant); err == nil {
				t.Fatal("skipped-step mutant graded PASS")
			}
		})
	}
}

type recordedGateFixture struct {
	root       string
	stateRoot  string
	entity     string
	briefing   string
	gateReview string
}

type recordedGateCommand struct {
	event  string
	argv   []string
	exit   int
	stdout string
	stderr string
}

func TestRecordedGateLifecycleRealCLIReplay(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	fixture := writeRecordedGateFixture(t)
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

	bind := run("briefing-record", "gate", "record", "recorded-gate-task", "--briefing", fixture.briefing, "--workflow-dir", fixture.root)
	open := run("open-validate", "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
	assertCommandOutput(t, bind.stdout, "state=open", "briefing=briefing:docs-dev:3k:validation:attempt-1:revision-1")
	assertCommandOutput(t, open.stdout, "state=open")
	assertCommandOutputField(t, open.stdout, "decision=")
	commitRecordedGateState(t, binary, fixture, "bind retained gate package")

	close := run("decision-record", "gate", "record", "recorded-gate-task",
		"--decision", "approve", "--actor", "agent:first-officer",
		"--reason", recordedGateReason,
		"--directive", recordedGateDirective,
		"--workflow-dir", fixture.root)
	closed := run("closed-validate", "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
	assertCommandOutput(t, close.stdout, "state=closed", "decision=approve")
	assertCommandOutput(t, closed.stdout, "state=closed", "decision=approve")
	commitRecordedGateState(t, binary, fixture, "record delegated gate decision")

	eligible := run("eligibility", "gate", "eligibility", "recorded-gate-task", "--workflow-dir", fixture.root)
	assertCommandOutput(t, eligible.stdout, "condition=approved-pending", "eligible=true", "target-stage=handoff")
	consume := run("consume", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
	assertCommandOutput(t, consume.stdout, "consumed=true", "target-stage=handoff")
	commitRecordedGateState(t, binary, fixture, "consume gate authorization")

	events := successfulRecordedGateEvents(commands)
	dispatches := 0
	if err := authorizeRecordedGateDispatch(events, readFile(t, fixture.entity), "handoff"); err == nil {
		dispatches++
	} else {
		t.Fatalf("dispatch oracle refused complete lifecycle: %v", err)
	}
	writeRecordedGateEvidence(t, fixture.root, commands, before, readFile(t, fixture.entity), readFile(t, fixture.gateReview), dispatches)
	observation := recordedGateObservation{
		events:       events,
		before:       before,
		after:        readFile(t, fixture.entity),
		dispatch:     recordedGateDispatchProof{spawned: dispatches == 1, handle: "fixture-worker-1", workerOutput: dispatches == 1},
		gateReview:   readFile(t, fixture.gateReview),
		expectedNext: "handoff",
	}
	if err := assertRecordedGateLifecycle(observation); err != nil {
		t.Fatal(err)
	}

	log := git(t, fixture.stateRoot, "show", "--name-only", "--format=", "HEAD~3..HEAD")
	for _, want := range []string{
		"recorded-gate-task/index.md",
		"recorded-gate-task/review/validation/briefing-1/briefing.json",
		"recorded-gate-task/review/validation/briefing-1/gate-review.md",
	} {
		if !strings.Contains(log, want) {
			t.Errorf("folder-form state commits omitted %s:\n%s", want, log)
		}
	}
	if strings.Contains(log, "dirty-sibling.md") {
		t.Fatalf("folder-form state commit swept dirty sibling:\n%s", log)
	}
}

func TestRecordedGateLifecycleRelativeBriefingFailsThenAbsoluteSucceeds(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	fixture := writeRecordedGateFixture(t)
	before := readFile(t, fixture.entity)
	rel, err := filepath.Rel(fixture.root, fixture.briefing)
	if err != nil {
		t.Fatal(err)
	}
	failed := runRecordedGateCommand(binary, fixture.root, "briefing-record", "gate", "record", "recorded-gate-task", "--briefing", rel, "--workflow-dir", fixture.root)
	if failed.exit == 0 || !strings.Contains(failed.stderr, "can't make") {
		t.Fatalf("relative retained input did not expose normalization defect: exit=%d stderr=%q", failed.exit, failed.stderr)
	}
	if after := readFile(t, fixture.entity); after != before {
		t.Fatal("relative-path refusal changed entity bytes")
	}
	if _, err := os.Stat(fixture.entity + ".gates.lock"); !os.IsNotExist(err) {
		t.Fatalf("relative-path refusal left lock residue: %v", err)
	}
	succeeded := runRecordedGateCommand(binary, fixture.root, "briefing-record", "gate", "record", "recorded-gate-task", "--briefing", fixture.briefing, "--workflow-dir", fixture.root)
	if succeeded.exit != 0 || !strings.Contains(succeeded.stdout, "state=open") {
		t.Fatalf("absolute retained input did not bind: exit=%d stdout=%q stderr=%q", succeeded.exit, succeeded.stdout, succeeded.stderr)
	}
}

func TestRecordedGateLifecycleReviseAndHoldDoNotApprovalDispatch(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	for _, decision := range []string{"revise", "hold"} {
		t.Run(decision, func(t *testing.T) {
			fixture := writeRecordedGateFixture(t)
			mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task", "--briefing", fixture.briefing, "--workflow-dir", fixture.root)
			mustRecordedGate(t, binary, fixture.root, "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
			mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task", "--decision", decision,
				"--actor", "agent:first-officer", "--reason", "named control finding", "--directive", recordedGateDirective, "--workflow-dir", fixture.root)
			mustRecordedGate(t, binary, fixture.root, "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)

			beforeEligibility := readFile(t, fixture.entity)
			eligibility := mustRecordedGate(t, binary, fixture.root, "gate", "eligibility", "recorded-gate-task", "--workflow-dir", fixture.root)
			assertCommandOutput(t, eligibility.stdout, "eligible=false")
			if got := readFile(t, fixture.entity); got != beforeEligibility {
				t.Fatalf("read-only %s eligibility changed entity bytes", decision)
			}
			if !strings.Contains(beforeEligibility, "status: validation") || strings.Contains(beforeEligibility, "state: consumed") {
				t.Fatalf("%s control advanced or consumed:\n%s", decision, beforeEligibility)
			}
			events := []string{"briefing-record", "open-validate", "decision-record", "closed-validate", "eligibility"}
			if err := authorizeRecordedGateDispatch(events, beforeEligibility, "handoff"); err == nil {
				t.Fatalf("%s control authorized successor dispatch without approval consume", decision)
			}
		})
	}
}

func TestRecordedGateLifecycleCapabilityStaleLauncherHaltsBeforeMutation(t *testing.T) {
	fixture := writeRecordedGateFixture(t)
	before := treeDigest(t, fixture.stateRoot)
	shim := filepath.Join(t.TempDir(), "spacedock")
	probeLog := filepath.Join(t.TempDir(), "probe.log")
	writeFile(t, shim, fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' \"$*\" >> %q\nif [ \"$1\" = \"--version\" ]; then echo 'spacedock 0.26.0+dev'; exit 0; fi\nif [ \"$1\" = \"gate\" ]; then echo 'record validate eligibility consume'; exit 0; fi\nexit 2\n", probeLog))
	if err := os.Chmod(shim, 0o755); err != nil {
		t.Fatal(err)
	}
	err := probeRecordedGateCapability(shim)
	if err == nil {
		t.Fatal("capability-stale launcher passed readiness")
	}
	for _, want := range []string{"record", "validate", "eligibility", "consume", "refresh", "go build", "SPACEDOCK_BIN"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("capability failure missing %q remediation: %v", want, err)
		}
	}
	if after := treeDigest(t, fixture.stateRoot); after != before {
		t.Fatal("capability preflight failure mutated the workflow")
	}
	probeLines := "\n" + readFile(t, probeLog) + "\n"
	for _, want := range []string{"gate --help", "gate record --help", "gate validate --help", "gate eligibility --help", "gate consume --help"} {
		if !strings.Contains(probeLines, "\n"+want+"\n") {
			t.Errorf("capability preflight omitted %q; probes=%s", want, probeLines)
		}
	}
	for _, probe := range strings.Split(strings.TrimSpace(probeLines), "\n") {
		if !strings.HasSuffix(probe, "--help") {
			t.Errorf("capability preflight performed mutation-capable command %q", probe)
		}
	}
	if err := probeRecordedGateCapability(buildRecordedGateBinary(t)); err != nil {
		t.Fatalf("fresh task-local binary failed capability preflight: %v", err)
	}
}

func TestRecordedGateLifecycleResumeConsumesOnce(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	fixture := writeRecordedGateFixture(t)
	mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task", "--briefing", fixture.briefing, "--workflow-dir", fixture.root)
	firstOpen := readFile(t, fixture.entity)
	mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task", "--briefing", fixture.briefing, "--workflow-dir", fixture.root)
	if got := readFile(t, fixture.entity); got != firstOpen {
		t.Fatal("same-Briefing resume minted a duplicate or changed bytes")
	}
	mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task",
		"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence green", "--directive", "you have the conn", "--workflow-dir", fixture.root)
	closed := readFile(t, fixture.entity)
	secondClose := runRecordedGateCommand(binary, fixture.root, "decision-record", "gate", "record", "recorded-gate-task",
		"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence green", "--directive", "you have the conn", "--workflow-dir", fixture.root)
	if secondClose.exit == 0 || readFile(t, fixture.entity) != closed {
		t.Fatal("closed resume re-recorded or changed the decision")
	}
	consumes := 0
	for pass := 0; pass < 3; pass++ {
		result := runRecordedGateCommand(binary, fixture.root, "consume", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		if result.exit == 0 && strings.Contains(result.stdout, "consumed=true") {
			consumes++
		}
	}
	if consumes != 1 {
		t.Fatalf("successful consumes across three resume passes = %d, want 1", consumes)
	}
	after := readFile(t, fixture.entity)
	if !strings.Contains(after, "status: handoff") || !strings.Contains(after, "state: consumed") {
		t.Fatalf("resume did not preserve atomic successor+consumed state:\n%s", after)
	}
}

func TestRecordedGateLifecycleShippedSkillMutantTurnsRed(t *testing.T) {
	root := recordedGateRepoRoot(t)
	original := readFile(t, filepath.Join(root, "skills", "first-officer", "references", "first-officer-shared-core.md"))
	mutant := strings.Replace(original, "${SPACEDOCK_BIN:-spacedock} gate eligibility ENTITY --workflow-dir WORKFLOW_DIR\n", "", 1)
	events := procedureEvents(mutant)
	if err := authorizeRecordedGateDispatch(events, "status: handoff\nstate: consumed", "handoff"); err == nil {
		t.Fatal("copied shipped-skill mutant without eligibility authorized dispatch")
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

func procedureEvents(skill string) []string {
	needles := []struct{ event, needle string }{
		{"briefing-record", "gate record ENTITY --briefing"},
		{"open-validate", "gate validate ENTITY"},
		{"decision-record", "gate record ENTITY --decision"},
		{"closed-validate", "gate validate ENTITY"},
		{"eligibility", "gate eligibility ENTITY"},
		{"consume", "gate consume ENTITY"},
	}
	var events []string
	from := 0
	for _, item := range needles {
		at := strings.Index(skill[from:], item.needle)
		if at < 0 {
			continue
		}
		events = append(events, item.event)
		from += at + len(item.needle)
	}
	return events
}

func recordedGatePrompt(workflowRoot string) string {
	return fmt.Sprintf("Use $spacedock:first-officer for this whole run.\n\nWorkflow directory: %s\n\nEngage only `recorded-gate-task`. Its retained validation package and concise gate review already exist. You have the conn toward the sprint goal; authorized to approve gates, relevant CI lanes, and dispatch; use your judgement. Exercise the normal First Officer gate procedure, preserve that exact quoted grant as delegated provenance, and continue through one observed successor dispatch. The scenario grades an actual `«worker.spawn»` host-tool call and the worker's durable output; `dispatch build`, narration that a worker is live, or `wait` with no worker handle is a failing observation. Stop after the handoff worker records %s; do not advance to terminal.", workflowRoot, recordedGateDispatchMarker)
}

func writeRecordedGateLoggingShim(t *testing.T, binary, logPath string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	shim := filepath.Join(dir, "spacedock")
	script := fmt.Sprintf("#!/bin/sh\nprintf 'begin\\t%%s\\n' \"$*\" >> %q\n%q \"$@\"\ncode=$?\nprintf 'exit=%%s\\t%%s\\n' \"$code\" \"$*\" >> %q\nexit \"$code\"\n", logPath, binary, logPath)
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

func recordedGateEventsFromCommandLog(log string) []string {
	var events []string
	validates := 0
	started := false
	for _, line := range strings.Split(log, "\n") {
		if !strings.HasPrefix(line, "exit=0\tgate ") || strings.Contains(line, " --help") {
			continue
		}
		switch {
		case strings.Contains(line, "gate record ") && strings.Contains(line, " --briefing "):
			started = true
			events = append(events, "briefing-record")
		case started && strings.Contains(line, "gate record ") && (strings.Contains(line, " --decision ") || strings.Contains(line, " --result ")):
			events = append(events, "decision-record")
		case started && strings.Contains(line, "gate validate "):
			validates++
			if validates == 1 {
				events = append(events, "open-validate")
			} else if validates == 2 {
				events = append(events, "closed-validate")
			}
		case started && strings.Contains(line, "gate eligibility "):
			events = append(events, "eligibility")
		case started && strings.Contains(line, "gate consume "):
			events = append(events, "consume")
		}
	}
	return events
}

func recordedGateEventsFromClaudeStream(stream string) []string {
	commands := map[string]string{}
	var successful []string
	for _, line := range strings.Split(stream, "\n") {
		var row struct {
			Message *struct {
				Content []struct {
					Type      string `json:"type"`
					Name      string `json:"name"`
					ID        string `json:"id"`
					ToolUseID string `json:"tool_use_id"`
					IsError   *bool  `json:"is_error"`
					Input     struct {
						Command string `json:"command"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &row) != nil || row.Message == nil {
			continue
		}
		for _, block := range row.Message.Content {
			if block.Type == "tool_use" && block.Name == "Bash" && block.ID != "" {
				commands[block.ID] = block.Input.Command
			}
			if block.Type == "tool_result" && block.ToolUseID != "" && block.IsError != nil && !*block.IsError {
				if command := commands[block.ToolUseID]; command != "" {
					successful = append(successful, command)
				}
			}
		}
	}
	var log strings.Builder
	for _, command := range successful {
		if at := strings.Index(command, "gate "); at >= 0 {
			fmt.Fprintf(&log, "exit=0\t%s\n", command[at:])
		}
	}
	return recordedGateEventsFromCommandLog(log.String())
}

func recordedGateReviewFromClaudeStream(stream string) string {
	var review string
	walkStreamBlocks(stream, func(block streamContentBlock) {
		if block.Type == "text" && strings.Contains(strings.ToLower(block.Text), "gate review") {
			review = block.Text
		}
	})
	return review
}

func recordedGateReviewFromCodexJSONL(jsonl string) string {
	var review string
	for _, line := range strings.Split(jsonl, "\n") {
		var row struct {
			Type string `json:"type"`
			Item struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"item"`
		}
		if json.Unmarshal([]byte(line), &row) == nil && row.Item.Type == "agent_message" &&
			strings.Contains(strings.ToLower(row.Item.Text), "gate review") {
			review = row.Item.Text
		}
	}
	return review
}

func recordedGateClaudeDispatchProof(stream, after string) recordedGateDispatchProof {
	type toolUse struct {
		name   string
		handle string
	}
	uses := map[string]toolUse{}
	consumeSucceeded := false
	proof := recordedGateDispatchProof{}
	for _, line := range strings.Split(stream, "\n") {
		var row struct {
			Message *struct {
				Content []struct {
					Type      string          `json:"type"`
					Name      string          `json:"name"`
					ID        string          `json:"id"`
					ToolUseID string          `json:"tool_use_id"`
					IsError   *bool           `json:"is_error"`
					Content   json.RawMessage `json:"content"`
					Input     struct {
						Command     string `json:"command"`
						Name        string `json:"name"`
						Prompt      string `json:"prompt"`
						Description string `json:"description"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &row) != nil || row.Message == nil {
			continue
		}
		for _, block := range row.Message.Content {
			if block.Type == "tool_use" && block.ID != "" {
				if block.Name == "Bash" {
					uses[block.ID] = toolUse{name: block.Name, handle: block.Input.Command}
				}
				if (block.Name == "Agent" || block.Name == "Task") && consumeSucceeded &&
					strings.Contains(block.Input.Prompt+block.Input.Description, "recorded-gate-task") {
					uses[block.ID] = toolUse{name: block.Name, handle: block.Input.Name}
				}
			}
			if block.Type != "tool_result" || block.ToolUseID == "" || block.IsError == nil || *block.IsError {
				continue
			}
			use := uses[block.ToolUseID]
			if use.name == "Bash" && strings.Contains(use.handle, "gate consume ") {
				consumeSucceeded = true
			}
			if use.name == "Agent" || use.name == "Task" {
				proof.spawned = true
				proof.handle = use.handle
				if match := claudeAgentIDResult.FindStringSubmatch(string(block.Content)); len(match) > 1 {
					proof.handle = match[1]
				}
				proof.workerOutput = strings.Contains(after, recordedGateDispatchMarker)
			}
		}
	}
	return proof
}

func recordedGateCodexDispatchProof(jsonl, after string) recordedGateDispatchProof {
	spawned := map[string]bool{}
	consumeSucceeded := false
	proof := recordedGateDispatchProof{}
	for _, line := range strings.Split(jsonl, "\n") {
		var row struct {
			Type string `json:"type"`
			Item struct {
				Type              string   `json:"type"`
				Tool              string   `json:"tool"`
				Command           string   `json:"command"`
				Prompt            string   `json:"prompt"`
				Status            string   `json:"status"`
				ExitCode          *int     `json:"exit_code"`
				ReceiverThreadIDs []string `json:"receiver_thread_ids"`
			} `json:"item"`
		}
		if json.Unmarshal([]byte(line), &row) != nil {
			continue
		}
		item := row.Item
		if row.Type == "item.completed" && item.Type == "command_execution" &&
			item.ExitCode != nil && *item.ExitCode == 0 && strings.Contains(item.Command, "gate consume ") {
			consumeSucceeded = true
		}
		if row.Type == "item.completed" && item.Type == "collab_tool_call" && item.Tool == "spawn_agent" &&
			consumeSucceeded && strings.Contains(item.Prompt, "recorded-gate-task") {
			for _, handle := range item.ReceiverThreadIDs {
				if handle != "" {
					spawned[handle] = true
					proof.spawned = true
					proof.handle = handle
				}
			}
		}
		if row.Type == "item.completed" && item.Type == "collab_tool_call" &&
			(item.Tool == "wait" || item.Tool == "wait_agent") && item.Status == "completed" {
			for _, handle := range item.ReceiverThreadIDs {
				if spawned[handle] && strings.Contains(after, recordedGateDispatchMarker) {
					proof.workerOutput = true
				}
			}
		}
	}
	return proof
}

func recordedGatePiObservation(session, after string) (recordedGateDispatchProof, string) {
	type piUse struct {
		name    string
		handle  string
		command string
	}
	uses := map[string]piUse{}
	consumeSucceeded := false
	proof := recordedGateDispatchProof{}
	var review string
	for _, line := range strings.Split(session, "\n") {
		var row struct {
			Message *struct {
				Role       string `json:"role"`
				ToolCallID string `json:"toolCallId"`
				ToolName   string `json:"toolName"`
				IsError    bool   `json:"isError"`
				Content    []struct {
					Type      string         `json:"type"`
					Text      string         `json:"text"`
					ID        string         `json:"id"`
					Name      string         `json:"name"`
					Arguments map[string]any `json:"arguments"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &row) != nil || row.Message == nil {
			continue
		}
		for _, block := range row.Message.Content {
			if block.Type == "text" && strings.Contains(strings.ToLower(block.Text), "gate review") {
				review = block.Text
			}
			if block.Type == "toolCall" && block.Name == "subagent" && block.ID != "" {
				encoded, _ := json.Marshal(block.Arguments)
				if consumeSucceeded && strings.Contains(string(encoded), "recorded-gate-task") {
					uses[block.ID] = piUse{name: block.Name, handle: block.ID}
				}
			}
			if block.Type == "toolCall" && (block.Name == "bash" || block.Name == "Bash") && block.ID != "" {
				command, _ := block.Arguments["command"].(string)
				uses[block.ID] = piUse{name: "bash", command: command}
			}
		}
		if row.Message.Role == "toolResult" && !row.Message.IsError {
			use := uses[row.Message.ToolCallID]
			if use.name == "bash" && strings.Contains(use.command, "gate consume ") {
				consumeSucceeded = true
			}
			if row.Message.ToolName == "subagent" && use.name == "subagent" {
				var resultText strings.Builder
				for _, block := range row.Message.Content {
					resultText.WriteString(block.Text)
				}
				proof.spawned = true
				proof.handle = use.handle
				proof.workerOutput = strings.Contains(resultText.String(), recordedGateDispatchMarker) &&
					strings.Contains(after, recordedGateDispatchMarker)
			}
		}
	}
	return proof, review
}

func TestRecordedGateLifecycleStructuredRuntimeEvidence(t *testing.T) {
	claude := strings.Join([]string{
		bashToolLine("help", "spacedock gate validate --help"),
		toolResultLine("help", false, "usage"),
		bashToolLine("failed", "spacedock gate record recorded-gate-task --briefing /tmp/briefing.json"),
		toolResultLine("failed", true, "exit 1"),
		bashToolLine("bind", "spacedock gate record recorded-gate-task --briefing /tmp/briefing.json"),
		toolResultLine("bind", false, "state=open"),
		bashToolLine("open", "spacedock gate validate recorded-gate-task"), toolResultLine("open", false, "state=open"),
		bashToolLine("close", "spacedock gate record recorded-gate-task --decision approve"), toolResultLine("close", false, "state=closed"),
		bashToolLine("closed", "spacedock gate validate recorded-gate-task"), toolResultLine("closed", false, "state=closed"),
		bashToolLine("eligible", "spacedock gate eligibility recorded-gate-task"), toolResultLine("eligible", false, "eligible=true"),
		bashToolLine("consume", "spacedock gate consume recorded-gate-task"), toolResultLine("consume", false, "consumed=true"),
		`{"type":"assistant","message":{"content":[{"type":"text","text":"Gate review: Capability; Evidence; Reviewed snapshot; Findings; Recommend; Decision"},{"type":"tool_use","id":"spawn","name":"Agent","input":{"name":"worker-1","prompt":"handoff recorded-gate-task"}}]}}`,
		toolResultLine("spawn", false, "agentId: a123abc"),
	}, "\n")
	if got := recordedGateEventsFromClaudeStream(claude); strings.Join(got, ",") != strings.Join(recordedGateRequiredEvents, ",") {
		t.Fatalf("successful Claude event trace = %v", got)
	}
	proof := recordedGateClaudeDispatchProof(claude, recordedGateDispatchMarker)
	if !proof.spawned || proof.handle != "a123abc" || !proof.workerOutput {
		t.Fatalf("Claude dispatch proof = %#v", proof)
	}
	if review := recordedGateReviewFromClaudeStream(claude); !strings.Contains(review, "Reviewed snapshot") {
		t.Fatalf("Claude presented review = %q", review)
	}

	codex := strings.Join([]string{
		`{"type":"item.completed","item":{"type":"command_execution","command":"spacedock gate consume recorded-gate-task","exit_code":0,"status":"completed"}}`,
		`{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","receiver_thread_ids":["thread-1"],"prompt":"handoff recorded-gate-task"}}`,
		`{"type":"item.completed","item":{"type":"collab_tool_call","tool":"wait_agent","receiver_thread_ids":["thread-1"],"status":"completed"}}`,
		`{"type":"item.completed","item":{"type":"agent_message","text":"Gate review: Capability; Evidence; Reviewed snapshot; Findings; Recommend; Decision"}}`,
	}, "\n")
	proof = recordedGateCodexDispatchProof(codex, recordedGateDispatchMarker)
	if !proof.spawned || proof.handle != "thread-1" || !proof.workerOutput {
		t.Fatalf("Codex dispatch proof = %#v", proof)
	}
	if review := recordedGateReviewFromCodexJSONL(codex); !strings.Contains(review, "Reviewed snapshot") {
		t.Fatalf("Codex presented review = %q", review)
	}

	pi := strings.Join([]string{
		`{"type":"message","message":{"role":"assistant","content":[{"type":"toolCall","id":"pi-consume","name":"bash","arguments":{"command":"spacedock gate consume recorded-gate-task"}}]}}`,
		`{"type":"message","message":{"role":"toolResult","toolCallId":"pi-consume","toolName":"bash","isError":false,"content":[{"type":"text","text":"consumed=true"}]}}`,
		`{"type":"message","message":{"role":"assistant","content":[{"type":"text","text":"Gate review: Capability; Evidence; Reviewed snapshot; Findings; Recommend; Decision"},{"type":"toolCall","id":"pi-call-1","name":"subagent","arguments":{"task":"handoff recorded-gate-task"}}]}}`,
		`{"type":"message","message":{"role":"toolResult","toolCallId":"pi-call-1","toolName":"subagent","isError":false,"content":[{"type":"text","text":"worker wrote RECORDED-GATE-SUCCESSOR-DISPATCHED"}]}}`,
	}, "\n")
	proof, review := recordedGatePiObservation(pi, recordedGateDispatchMarker)
	if !proof.spawned || proof.handle != "pi-call-1" || !proof.workerOutput {
		t.Fatalf("Pi dispatch proof = %#v", proof)
	}
	if !strings.Contains(review, "Reviewed snapshot") {
		t.Fatalf("Pi presented review = %q", review)
	}
	earlyPi := strings.Join([]string{
		`{"type":"message","message":{"role":"assistant","content":[{"type":"toolCall","id":"early","name":"subagent","arguments":{"task":"handoff recorded-gate-task"}}]}}`,
		`{"type":"message","message":{"role":"toolResult","toolCallId":"early","toolName":"subagent","isError":false,"content":[{"type":"text","text":"RECORDED-GATE-SUCCESSOR-DISPATCHED"}]}}`,
		`{"type":"message","message":{"role":"assistant","content":[{"type":"toolCall","id":"late-consume","name":"bash","arguments":{"command":"spacedock gate consume recorded-gate-task"}}]}}`,
		`{"type":"message","message":{"role":"toolResult","toolCallId":"late-consume","toolName":"bash","isError":false,"content":[{"type":"text","text":"consumed=true"}]}}`,
	}, "\n")
	if earlyProof, _ := recordedGatePiObservation(earlyPi, recordedGateDispatchMarker); earlyProof.spawned {
		t.Fatalf("Pi dispatch before successful consume was accepted: %#v", earlyProof)
	}
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
	root := t.TempDir()
	stateRoot := filepath.Join(root, ".spacedock-state")
	writeFile(t, filepath.Join(root, "README.md"), recordedGateReadme())
	entity := filepath.Join(stateRoot, "recorded-gate-task", "index.md")
	writeFile(t, entity, recordedGateEntity())
	gitInit(t, root)
	gitInit(t, stateRoot)
	room := filepath.Join(filepath.Dir(entity), "review", "validation", "briefing-1")
	briefing := filepath.Join(room, "briefing.json")
	fixtureBriefing := filepath.Join(recordedGateRepoRoot(t), "internal", "gates", "testdata", "exact-validation-briefing.json")
	writeFile(t, briefing, readFile(t, fixtureBriefing))
	gateReview := filepath.Join(room, "gate-review.md")
	writeFile(t, gateReview, recordedGateReview())
	writeFile(t, filepath.Join(room, "entity-snapshot.md"), "# Frozen entity snapshot\n\nExact reviewed state before the gate decision.\n")
	writeFile(t, filepath.Join(room, "contract-snapshot.md"), "# Recorder contract snapshot\n\nThe retained 3k recorder and h1 one-use application contract.\n")
	writeFile(t, filepath.Join(stateRoot, "dirty-sibling.md"), "unrelated concurrent dirt\n")
	return recordedGateFixture{root: root, stateRoot: stateRoot, entity: entity, briefing: briefing, gateReview: gateReview}
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
		"id: recorded-gate-task\n" +
		"title: Recorded Gate Task\n" +
		"status: validation\n" +
		"completed:\nverdict:\nworktree:\n" +
		"---\n# Recorded Gate Task\n\n" +
		"## Acceptance criteria\n\n**AC-1** Successor dispatch requires consumed approval.\n\n" +
		"## Stage Report: validation\n\n- DONE: Replayed retained evidence\n  The real command fixture is green.\n\n" +
		"### Summary\n\nReady for the recorded decision gate.\n"
}

func recordedGateReview() string {
	return "# Gate review: recorded First Officer lifecycle\n\n" +
		"Capability/change: the FO now calls the landed recorder and one-use application commands.\n\n" +
		"Test and evidence: fresh-binary command replay, byte comparisons, and skipped-step mutants pass.\n\n" +
		"Reviewed snapshot: the retained 3k validation Briefing at its exact declared digest.\n\n" +
		"Findings: no material findings; CLI path normalization remains a named deferred product issue.\n\n" +
		"Recommendation: approve and consume the authorization once.\n\n" +
		"Decision ask: approve, revise with a concrete finding, or hold for a named prerequisite?\n\n" +
		"References: entity, recorder contract, and briefing.json are linked in the package.\n"
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

func commitRecordedGateState(t *testing.T, binary string, fixture recordedGateFixture, message string) {
	t.Helper()
	result := runRecordedGateCommand(binary, fixture.root, "state-commit", "state", "commit", "recorded-gate-task", "--workflow-dir", fixture.root, "-m", message)
	if result.exit != 0 {
		t.Fatalf("state commit exit=%d stdout=%q stderr=%q", result.exit, result.stdout, result.stderr)
	}
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

func assertCommandOutputField(t *testing.T, output, want string) {
	t.Helper()
	for _, field := range strings.Fields(output) {
		if field == want {
			return
		}
	}
	t.Fatalf("command output missing exact field %q: %s", want, output)
}

func probeRecordedGateCapability(binary string) error {
	missing := []string{}
	checks := []struct {
		argv  []string
		wants []string
	}{
		{[]string{"gate", "--help"}, []string{"record", "validate", "eligibility", "consume"}},
		{[]string{"gate", "record", "--help"}, []string{"--briefing", "--result", "--association", "--decision", "--actor", "--directive"}},
		{[]string{"gate", "validate", "--help"}, []string{"gate validate <entity>"}},
		{[]string{"gate", "eligibility", "--help"}, []string{"gate eligibility <entity>"}},
		{[]string{"gate", "consume", "--help"}, []string{"gate consume <entity>"}},
	}
	for _, check := range checks {
		out, err := exec.Command(binary, check.argv...).CombinedOutput()
		if err != nil {
			missing = append(missing, strings.Join(check.argv[:len(check.argv)-1], " "))
			continue
		}
		for _, want := range check.wants {
			if !strings.Contains(string(out), want) {
				missing = append(missing, strings.Join(check.argv[:len(check.argv)-1], " ")+":"+want)
			}
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("gate capability missing %s; refresh the launcher or go build -o <temp>/spacedock ./cmd/spacedock and set SPACEDOCK_BIN to it", strings.Join(missing, ", "))
	}
	return nil
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

func writeRecordedGateEvidence(t *testing.T, root string, commands []recordedGateCommand, before, after, review string, dispatches int) {
	t.Helper()
	dir := filepath.Join(root, "evidence")
	var log strings.Builder
	for _, command := range commands {
		fmt.Fprintf(&log, "event=%s exit=%d argv=%q\nstdout=%s\nstderr=%s\n", command.event, command.exit, command.argv[1:], command.stdout, command.stderr)
	}
	writeFile(t, filepath.Join(dir, "command.log"), log.String())
	writeFile(t, filepath.Join(dir, "entity.before.md"), before)
	writeFile(t, filepath.Join(dir, "entity.after.md"), after)
	writeFile(t, filepath.Join(dir, "gate-review.md"), review)
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
