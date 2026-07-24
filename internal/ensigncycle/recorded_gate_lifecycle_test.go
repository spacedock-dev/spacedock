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
	"decision-record",
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
	builds         int
	durableEffects int
	ordered        bool
	committed      bool
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
	if o.dispatch.builds != 1 {
		return fmt.Errorf("successful successor dispatch builds = %d, want 1", o.dispatch.builds)
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
		{"approval decision", "\n                decision: approve", 1},
		{"approval actor", "by: agent:first-officer", 1},
		{"approval reason", "\n                reason:", 1},
		{"delegated directive", recordedGateDirective, 1},
		{"application target", "target-stage: " + o.expectedNext, 1},
		{"consumed application", "\n                state: consumed", 1},
	} {
		if got := strings.Count(o.after, exact.value); got != exact.count || (exact.label == "approval reason" && strings.Trim(strings.TrimSpace(strings.SplitN(strings.SplitN(o.after, exact.value, 2)[1], "\n", 2)[0]), `"'`) == "") {
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
	trimmed := strings.TrimSpace(review)
	lower := strings.ToLower(trimmed)
	fields := []string{"capability/change:", "test and evidence:", "reviewed snapshot:", "findings:", "recommendation:", "decision ask:"}
	last := -1
	for i, field := range fields {
		if got := strings.Count(lower, field); got != 1 {
			return fmt.Errorf("gate review field %q count=%d, want 1", field, got)
		}
		at := strings.Index(lower, field)
		if at <= last {
			return fmt.Errorf("gate review field %q is out of order", field)
		}
		end := len(trimmed)
		if i+1 < len(fields) {
			end = strings.Index(lower, fields[i+1])
		}
		value := strings.TrimSpace(trimmed[at+len(field) : end])
		if value == "" {
			return fmt.Errorf("gate review field %q is blank", field)
		}
		last = at
	}
	if !strings.Contains(trimmed, recordedGateBriefingID) || !strings.Contains(trimmed, recordedGateDigest) {
		return fmt.Errorf("gate review does not name the exact retained snapshot identity and digest")
	}
	ask := lower[strings.Index(lower, "decision ask:"):]
	for _, decision := range []string{"approve", "revise", "hold"} {
		if !strings.Contains(ask, decision) {
			return fmt.Errorf("gate review decision ask omits %q", decision)
		}
	}
	if strings.HasPrefix(lower, "---") || strings.HasPrefix(lower, "gates:") ||
		strings.HasPrefix(lower, "{") || strings.HasPrefix(lower, "[") {
		return fmt.Errorf("gate review leads with a raw entity, Briefing, or room dump")
	}
	for _, forbidden := range []string{"\ngates:\n", `"type":"briefing"`, "\"artifacts\":[", "rebind ceremony", "supersede ceremony"} {
		if strings.Contains(lower, forbidden) {
			return fmt.Errorf("gate review contains raw or recorder-mechanics payload %q", forbidden)
		}
	}
	return nil
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
	assertCommandOutput(t, bind.stdout, "state=open", "briefing=briefing:docs-dev:3k:validation:attempt-1:revision-1")
	commitRecordedGateState(t, binary, fixture, "bind retained gate package")

	close := run("decision-record", "gate", "record", "recorded-gate-task",
		"--decision", "approve", "--actor", "agent:first-officer",
		"--reason", recordedGateReason,
		"--directive", recordedGateDirective,
		"--workflow-dir", fixture.root)
	assertCommandOutput(t, close.stdout, "state=closed", "decision=approve")
	closeCommit := commitRecordedGateState(t, binary, fixture, "record delegated gate decision")

	consume := run("consume", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
	assertCommandOutput(t, consume.stdout, "consumed=true", "target-stage=handoff")
	consumedCommit := commitRecordedGateState(t, binary, fixture, "consume gate authorization")

	events := successfulRecordedGateEvents(commands)
	dispatches := 0
	if err := authorizeRecordedGateDispatch(events, readFile(t, fixture.entity), "handoff"); err == nil {
		dispatches++
	} else {
		t.Fatalf("dispatch oracle refused complete lifecycle: %v", err)
	}
	writeRecordedGateEvidence(t, fixture.root, commands, before, readFile(t, fixture.entity), readFile(t, fixture.gateReview), dispatches)
	observation := recordedGateObservation{
		events: events,
		before: before,
		after:  readFile(t, fixture.entity),
		dispatch: recordedGateDispatchProof{
			builds: dispatches, durableEffects: dispatches, ordered: true,
			committed: recordedGateCommittedBeforeDispatch(t, fixture, closeCommit, consumedCommit, consumedCommit),
		},
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

func TestRecordedGateLifecycleRelativeBriefingMatchesAbsolute(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	relative, absolute := writeRecordedGateFixture(t), writeRecordedGateFixture(t)
	rel, err := filepath.Rel(relative.root, relative.briefing)
	if err != nil {
		t.Fatal(err)
	}
	mustRecordedGate(t, binary, relative.root, "gate", "record", "recorded-gate-task", "--briefing", rel, "--workflow-dir", relative.root)
	mustRecordedGate(t, binary, absolute.root, "gate", "record", "recorded-gate-task", "--briefing", absolute.briefing, "--workflow-dir", absolute.root)
	if got, want := readFile(t, relative.entity), readFile(t, absolute.entity); got != want {
		t.Fatalf("relative and absolute retained inputs bound different entity bytes")
	}
}

func TestRecordedGateLifecycleCapabilityStaleLauncherHaltsBeforeMutation(t *testing.T) {
	fixture := writeRecordedGateFixture(t)
	before := treeDigest(t, fixture.stateRoot)
	fresh := buildRecordedGateBinary(t)
	cache := map[string]bool{}
	shim := filepath.Join(t.TempDir(), "spacedock")
	probeLog := filepath.Join(t.TempDir(), "probe.log")
	writeFile(t, shim, fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' \"$*\" >> %q\nexec %q \"$@\"\n", probeLog, fresh))
	if err := os.Chmod(shim, 0o755); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err := probeRecordedGateCapability(shim, cache); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, shim, fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' \"$*\" >> %q\necho 'record validate eligibility consume'\n", probeLog))
	err := probeRecordedGateCapability(shim, cache)
	if err == nil {
		t.Fatal("same-path stale replacement passed readiness")
	}
	assertCommandOutput(t, err.Error(), "--briefing", "--result", "--decision", "refresh", "go build", "SPACEDOCK_BIN")
	if after := treeDigest(t, fixture.stateRoot); after != before {
		t.Fatal("capability preflight failure mutated the workflow")
	}
	if got := strings.Count(readFile(t, probeLog), "gate --help"); got != 2 {
		t.Fatalf("capability fingerprint probes=%d, want cached capable plus replaced stale", got)
	}
	capableDir, staleDir := filepath.Dir(fresh), t.TempDir()
	if err := os.Rename(shim, filepath.Join(staleDir, "spacedock")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", capableDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	if err := probeRecordedGateCapability("spacedock", cache); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", staleDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	if err := probeRecordedGateCapability("spacedock", cache); err == nil {
		t.Fatal("PATH target swap reused capable cache identity")
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
	mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task",
		"--briefing", fixture.briefing, "--workflow-dir", fixture.root)
}
func closeRecordedGate(t *testing.T, binary string, fixture recordedGateFixture, decision string) {
	mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task",
		"--decision", decision, "--actor", "agent:first-officer", "--reason", "evidence-backed route",
		"--directive", recordedGateDirective, "--workflow-dir", fixture.root)
}
func TestRecordedGateLifecycleAC5RefusalMatrix(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	t.Run("missing-briefing", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		before := treeDigest(t, fixture.stateRoot)
		result := runRecordedGateCommand(binary, fixture.root, "", "gate", "record", "recorded-gate-task",
			"--briefing", filepath.Join(filepath.Dir(fixture.briefing), "missing", "briefing.json"), "--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, result, "briefing")
		if after := treeDigest(t, fixture.stateRoot); after != before {
			t.Fatal("missing Briefing changed workflow bytes")
		}
	})
	t.Run("alternate-basename", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		alternate := filepath.Join(filepath.Dir(fixture.briefing), "manifest.json")
		writeFile(t, alternate, readFile(t, fixture.briefing))
		before := treeDigest(t, fixture.stateRoot)
		result := runRecordedGateCommand(binary, fixture.root, "", "gate", "record", "recorded-gate-task",
			"--briefing", alternate, "--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, result, "briefing.json")
		if after := treeDigest(t, fixture.stateRoot); after != before {
			t.Fatal("alternate basename changed workflow bytes")
		}
	})
	for _, tc := range []struct {
		name  string
		args  []string
		wants []string
	}{
		{"actor", []string{"--decision", "approve", "--reason", "evidence"}, []string{"actor"}},
		{"reason", []string{"--decision", "revise", "--actor", "agent:first-officer", "--directive", recordedGateDirective}, []string{"reason"}},
		{"directive", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence"}, []string{"directive"}},
	} {
		t.Run("invalid-"+tc.name, func(t *testing.T) {
			fixture := writeRecordedGateFixture(t)
			bindRecordedGate(t, binary, fixture)
			before := treeDigest(t, fixture.stateRoot)
			args := append([]string{"gate", "record", "recorded-gate-task"}, tc.args...)
			args = append(args, "--workflow-dir", fixture.root)
			result := runRecordedGateCommand(binary, fixture.root, "", args...)
			assertRecordedGateByteCleanFailure(t, fixture, result, tc.wants...)
			if after := treeDigest(t, fixture.stateRoot); after != before {
				t.Fatalf("invalid %s changed workflow bytes", tc.name)
			}
		})
	}
	t.Run("association", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		resultPath := filepath.Join(t.TempDir(), "result.json")
		association := filepath.Join(t.TempDir(), "association.json")
		testdata := filepath.Join(recordedGateRepoRoot(t), "internal", "gates", "testdata")
		writeFile(t, resultPath, readFile(t, filepath.Join(testdata, "exact-review-v1-result.json")))
		writeFile(t, association, readFile(t, filepath.Join(testdata, "exact-review-v1-association-truncated.json")))
		before := treeDigest(t, fixture.stateRoot)
		result := runRecordedGateCommand(binary, fixture.root, "", "gate", "record", "recorded-gate-task",
			"--result", resultPath, "--association", association, "--actor", "person:reviewer", "--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, result, "association")
		if after := treeDigest(t, fixture.stateRoot); after != before {
			t.Fatal("invalid association changed workflow bytes")
		}
	})
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
	for _, decision := range []string{"hold", "revise"} {
		t.Run(decision+"-consume", func(t *testing.T) {
			fixture := writeRecordedGateFixture(t)
			bindRecordedGate(t, binary, fixture)
			closeRecordedGate(t, binary, fixture, decision)
			closeCommit := commitRecordedGateState(t, binary, fixture, "durably record "+decision)
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
	t.Run("stale-consume", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		closeRecordedGate(t, binary, fixture, "approve")
		writeFile(t, fixture.briefing, strings.Replace(readFile(t, fixture.briefing), `"question": "Should`, `"question": "Must`, 1))
		beforeConsume := recordedGateTreeSnapshot(t, fixture.stateRoot)
		beforeEligibility := treeDigest(t, fixture.stateRoot)
		eligibility := mustRecordedGate(t, binary, fixture.root, "gate", "eligibility", "recorded-gate-task", "--workflow-dir", fixture.root)
		assertCommandOutput(t, eligibility.stdout, "condition=stale", "eligible=false")
		if after := treeDigest(t, fixture.stateRoot); after != beforeEligibility {
			t.Fatal("stale eligibility changed workflow bytes")
		}
		result := runRecordedGateCommand(binary, fixture.root, "", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		assertRecordedGateByteCleanFailure(t, fixture, result, "stale", "superseded")
		entityRel := strings.TrimPrefix(fixture.entity, fixture.stateRoot+string(os.PathSeparator))
		expectedEntity := strings.Replace(beforeConsume[entityRel], "state: pending", "state: superseded", 1)
		if expectedEntity == beforeConsume[entityRel] {
			t.Fatal("stale control could not construct the independent pending→superseded expectation")
		}
		beforeConsume[entityRel] = expectedEntity
		assertRecordedGateTreeSnapshot(t, fixture.stateRoot, beforeConsume)
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
		if after := readFile(t, fixture.entity); strings.Count(after, "gate-attempt:3k-validation-1") != 1 {
			t.Fatal("same-package open resume minted an attempt")
		}
	})
	t.Run("open-changed-package", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		before := readFile(t, fixture.entity)
		writeFile(t, fixture.briefing, strings.Replace(readFile(t, fixture.briefing), `"question": "Should`, `"question": "Must`, 1))
		replacementDigest := canonicalRecordedGateDigest(t, readFile(t, fixture.briefing))
		expectedTree := recordedGateTreeSnapshot(t, fixture.stateRoot)
		mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task",
			"--briefing", fixture.briefing, "--workflow-dir", fixture.root)
		after := readFile(t, fixture.entity)
		expected := strings.Replace(before, recordedGateDigest, replacementDigest, 1)
		if expected == before {
			t.Fatal("changed-package control could not construct an independent replacement binding")
		}
		if after != expected || strings.Count(after, recordedGateBriefingID) != 1 ||
			strings.Contains(after, recordedGateDigest) || strings.Count(after, replacementDigest) != 1 {
			t.Fatalf("changed-package resume changed fields beyond the exact replacement digest\n--- expected ---\n%s\n--- after ---\n%s", expected, after)
		}
		entityRel := strings.TrimPrefix(fixture.entity, fixture.stateRoot+string(os.PathSeparator))
		expectedTree[entityRel] = expected
		assertRecordedGateTreeSnapshot(t, fixture.stateRoot, expectedTree)
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
	t.Run("stale-replacement", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		closeRecordedGate(t, binary, fixture, "approve")
		writeFile(t, fixture.briefing, strings.Replace(readFile(t, fixture.briefing), `"question": "Should`, `"question": "Must`, 1))
		stale := runRecordedGateCommand(binary, fixture.root, "", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		if stale.exit == 0 || !strings.Contains(stale.stdout, "application=advance/superseded") {
			t.Fatalf("stale consume did not materialize supersession: %#v", stale)
		}
		writeFile(t, fixture.briefing, strings.Replace(readFile(t, fixture.briefing),
			"attempt-1:revision-1", "attempt-2:revision-1", 1))
		mustRecordedGate(t, binary, fixture.root, "gate", "record", "recorded-gate-task",
			"--briefing", fixture.briefing, "--workflow-dir", fixture.root)
		open := mustRecordedGate(t, binary, fixture.root, "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
		assertCommandOutput(t, open.stdout, "state=open")
		after := readFile(t, fixture.entity)
		if strings.Count(after, "gate-attempt:3k-validation-") != 2 ||
			strings.Count(after, "resolution:spacedock") != 1 ||
			!strings.Contains(after, "state: superseded") {
			t.Fatalf("stale resume did not supersede once and bind one replacement:\n%s", after)
		}
	})
	t.Run("consumed", func(t *testing.T) {
		fixture := writeRecordedGateFixture(t)
		bindRecordedGate(t, binary, fixture)
		closeRecordedGate(t, binary, fixture, "approve")
		mustRecordedGate(t, binary, fixture.root, "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
		afterFirst := readFile(t, fixture.entity)
		afterFirstTree := recordedGateTreeSnapshot(t, fixture.stateRoot)
		for pass := 0; pass < 3; pass++ {
			mustRecordedGate(t, binary, fixture.root, "gate", "validate", "recorded-gate-task", "--workflow-dir", fixture.root)
			repeat := runRecordedGateCommand(binary, fixture.root, "", "gate", "consume", "recorded-gate-task", "--workflow-dir", fixture.root)
			if repeat.exit == 0 || !strings.Contains(repeat.stdout, "condition=consumed") {
				t.Fatalf("consumed resume pass %d was not refused: %#v", pass, repeat)
			}
		}
		after := readFile(t, fixture.entity)
		assertRecordedGateTreeSnapshot(t, fixture.stateRoot, afterFirstTree)
		if after != afterFirst || strings.Count(after, "state: consumed") != 1 ||
			strings.Count(after, "resolution:spacedock") != 1 {
			t.Fatal("consumed resume duplicated transition, application, or resolution")
		}
	})
}
func recordedGateDiscovery(t *testing.T, binary, root string) []string {
	t.Helper()
	result := runRecordedGateCommand(binary, root, "", "status", "--boot", "--identify", "--json")
	if result.exit != 0 {
		t.Fatalf("discovery exit=%d stdout=%q stderr=%q", result.exit, result.stdout, result.stderr)
	}
	var boot struct {
		Discovery []string `json:"discovery"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &boot); err != nil {
		t.Fatal(err)
	}
	return boot.Discovery
}
func TestRecordedGateLifecycleWorkflowDiscoveryEquality(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	root := t.TempDir()
	workflow := filepath.Join(root, "workflow")
	writeFile(t, filepath.Join(workflow, "README.md"), strings.Replace(recordedGateReadme(), "state: .spacedock-state\n", "", 1))
	writeFile(t, filepath.Join(workflow, "seed.md"), recordedGateEntity())
	gitInit(t, root)
	before := recordedGateDiscovery(t, binary, root)
	fixture := writeRecordedGateFixtureAt(t, filepath.Join(root, "testdata", "recorded-gate-run"))
	bindRecordedGate(t, binary, fixture)
	after := recordedGateDiscovery(t, binary, root)
	if strings.Join(before, "\n") != strings.Join(after, "\n") {
		t.Fatalf("fixture execution polluted workflow discovery: before=%v after=%v", before, after)
	}
	planted := filepath.Join(root, "planted")
	writeFile(t, filepath.Join(planted, "README.md"), strings.Replace(recordedGateReadme(), "state: .spacedock-state\n", "", 1))
	writeFile(t, filepath.Join(planted, "seed.md"), recordedGateEntity())
	control := recordedGateDiscovery(t, binary, root)
	if strings.Join(before, "\n") == strings.Join(control, "\n") {
		t.Fatal("planted discoverable workflow did not turn equality red")
	}
}
func TestRecordedGateLifecycleShippedSkillMutantTurnsRed(t *testing.T) {
	root := recordedGateRepoRoot(t)
	original := readFile(t, filepath.Join(root, "skills", "fo-gate-lifecycle", "SKILL.md"))
	if events := procedureEvents(original); strings.Join(events, ",") != strings.Join(recordedGateRequiredEvents, ",") {
		t.Fatalf("shipped skill baseline trace=%v, want %v", events, recordedGateRequiredEvents)
	}
	commands := []struct {
		event, command string
	}{
		{"briefing-record", "gate record ENTITY --briefing BRIEFING"},
		{"decision-record", "gate record ENTITY --decision approve|revise|hold"},
		{"consume", "gate consume ENTITY"},
	}
	for _, tc := range commands {
		t.Run(tc.event, func(t *testing.T) {
			at := strings.Index(original, tc.command)
			if at < 0 {
				t.Fatalf("shipped skill has no %s command", tc.event)
			}
			mutant := original[:at] + original[at+len(tc.command):]
			if tc.event == "decision-record" {
				mutant = strings.ReplaceAll(original, tc.command, "")
			}
			copyPath := filepath.Join(t.TempDir(), "skills", "fo-gate-lifecycle", "SKILL.md")
			writeFile(t, copyPath, mutant)
			if events := procedureEvents(readFile(t, copyPath)); strings.Join(events, ",") == strings.Join(recordedGateRequiredEvents, ",") {
				t.Fatalf("copied shipped-skill %s deletion kept the three-event grader green: %v", tc.event, events)
			}
		})
	}
}
func TestRecordedGateLifecycleProvenanceAndPresentationMutants(t *testing.T) {
	valid := recordedGateObservation{
		events: append([]string(nil), recordedGateRequiredEvents...),
		before: "status: validation",
		after: "status: handoff\ngate: gate:docs-dev:3k:validation\nid: gate-attempt:3k-validation-1\n" +
			"id: " + recordedGateBriefingID + "\ndigest: " + recordedGateDigest + "\n" +
			"id: resolution:spacedock:docs-dev:3k:validation:1\nbriefing: " + recordedGateBriefingID + "\n" +
			"by: agent:first-officer\n                decision: approve\n                reason: " + recordedGateReason + "\n" +
			"adoption-note: '" + recordedGateDirective + "'\ntarget-stage: handoff\n                state: consumed\nreport repeats decision: approve",
		dispatch:     recordedGateDispatchProof{builds: 1, durableEffects: 1, ordered: true, committed: true},
		gateReview:   recordedGateReview(),
		expectedNext: "handoff",
	}
	if err := assertRecordedGateLifecycle(valid); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	for name, mutate := range map[string]func(*recordedGateObservation){
		"actor-swap": func(o *recordedGateObservation) {
			o.after = strings.Replace(o.after, "by: agent:first-officer", "by: person:captain", 1)
		},
		"blank-reason": func(o *recordedGateObservation) { o.after = strings.Replace(o.after, recordedGateReason, "", 1) },
		"altered-directive": func(o *recordedGateObservation) {
			o.after = strings.Replace(o.after, recordedGateDirective, "you have partial conn", 1)
		},
		"raw-entity-dump": func(o *recordedGateObservation) { o.gateReview = "---\nstatus: validation\n---\n" + o.gateReview },
		"raw-briefing-dump": func(o *recordedGateObservation) {
			o.gateReview = `{"type":"briefing","artifacts":[]}` + "\n" + o.gateReview
		},
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

func procedureEvents(skill string) []string {
	needles := []struct{ event, needle string }{
		{"briefing-record", "gate record ENTITY --briefing"},
		{"decision-record", "gate record ENTITY --decision"},
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
	return fmt.Sprintf("Use $spacedock:first-officer for this whole run.\n\nWorkflow directory: %s\n\nEngage only `recorded-gate-task`. Its retained validation package and concise gate review already exist; the canonical Briefing is `%s`. Copy the exact delegated conn bytes between these delimiters, including the final period immediately before END_CONN:\nBEGIN_CONN\n%s\nEND_CONN\nExercise the normal First Officer gate procedure, preserve those exact bytes as delegated provenance, and continue through one successor dispatch. Follow its capability preflight: successfully run exactly one `gate --help` before the first gate mutation. Before deciding, emit one assistant-text gate review with exactly these six nonblank labels in order: `Capability/change:`, `Test and evidence:`, `Reviewed snapshot:`, `Findings:`, `Recommendation:`, `Decision ask:`; copy the bound Briefing identity and digest from entity state after recording (do not calculate a file hash or use an artifact `rev`), and offer approve/revise/hold in the decision ask. Run `dispatch build` once successfully, then dispatch that exact artifact without rebuilding it. On Pi subagents use executable agent `worker`, not the artifact's semantic `subagent_type`. Stop after the handoff worker records %s in one new durable stage report/commit; do not advance to terminal.", workflowRoot, filepath.Join(workflowRoot, ".spacedock-state", "recorded-gate-task", "review", "validation", "briefing-1", "briefing.json"), recordedGateDirective, recordedGateDispatchMarker)
}

func writeRecordedGateLoggingShim(t *testing.T, binary, logPath string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	shim := filepath.Join(dir, "spacedock")
	script := fmt.Sprintf("#!/bin/sh\nprintf 'begin\\t%%s\\n' \"$*\" >> %q\n[ \"$1 $2\" = \"dispatch build\" ] && git -C .spacedock-state rev-parse HEAD | sed 's/^/dispatch-head\\t/' >> %q\n%q \"$@\"\ncode=$?\nprintf 'exit=%%s\\t%%s\\n' \"$code\" \"$*\" >> %q\n[ \"$code\" -eq 0 ] && [ \"$1 $2\" = \"state commit\" ] && git -C .spacedock-state rev-parse HEAD | sed 's/^/state-head\\t/' >> %q\nexit \"$code\"\n", logPath, logPath, binary, logPath, logPath)
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
			fixture := writeRecordedGateFixture(t)
			steps := [][]string{
				{"gate", "record", "recorded-gate-task", "--briefing", fixture.briefing, "--workflow-dir", fixture.root},
				{"gate", "record", "recorded-gate-task", "--decision", "approve", "--actor", "agent:first-officer", "--reason", recordedGateReason, "--directive", recordedGateDirective, "--workflow-dir", fixture.root},
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
func recordedGateReviewFromClaudeStream(stream string) string {
	var review string
	walkStreamBlocks(stream, func(block streamContentBlock) {
		if block.Type == "text" && strings.Contains(strings.ToLower(block.Text), "capability/change:") {
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
			strings.Contains(strings.ToLower(row.Item.Text), "capability/change:") {
			review = row.Item.Text
		}
	}
	return review
}

func recordedGateReviewFromPiSession(session string) string {
	var review string
	for _, line := range strings.Split(session, "\n") {
		var row struct {
			Message *struct {
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal([]byte(line), &row) != nil || row.Message == nil {
			continue
		}
		for _, block := range row.Message.Content {
			if block.Type == "text" && strings.Contains(strings.ToLower(block.Text), "capability/change:") {
				review = block.Text
			}
		}
	}
	return review
}

func recordedGateLiveObservation(t *testing.T, fixture recordedGateFixture, before, commandLog, review string) recordedGateObservation {
	t.Helper()
	log := readFile(t, commandLog)
	builds, consumed, ordered := 0, false, true
	ordered = strings.Index(log, "exit=0\tgate --help") >= 0 && strings.Index(log, "exit=0\tgate record ") > strings.Index(log, "exit=0\tgate --help")
	stateHeads, dispatchHead := []string{}, ""
	for _, line := range strings.Split(log, "\n") {
		if strings.HasPrefix(line, "exit=0\tgate consume ") {
			consumed = true
		}
		if strings.HasPrefix(line, "exit=0\tdispatch build ") {
			builds++
			ordered = ordered && consumed
		}
		if strings.HasPrefix(line, "state-head\t") {
			stateHeads = append(stateHeads, strings.TrimPrefix(line, "state-head\t"))
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
	if len(commits) == 1 && strings.Contains(after, recordedGateDispatchMarker) {
		effects = 1
	}
	closeCommit, consumedCommit := "", ""
	for _, head := range stateHeads {
		snapshot := recordedGateEntityAt(t, fixture, head)
		switch {
		case closeCommit == "" && strings.Contains(snapshot, "decision: approve") && strings.Contains(snapshot, "state: pending"):
			closeCommit = head
		case consumedCommit == "" && strings.Contains(snapshot, "state: consumed"):
			consumedCommit = head
		}
	}
	return recordedGateObservation{
		events: recordedGateEventsFromCommandLog(log), before: before, after: after,
		dispatch: recordedGateDispatchProof{
			builds: builds, durableEffects: effects, ordered: ordered,
			committed: recordedGateCommittedBeforeDispatch(t, fixture, closeCommit, consumedCommit, dispatchHead),
		},
		gateReview: review, expectedNext: "handoff",
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
	return writeRecordedGateFixtureAt(t, t.TempDir())
}
func writeRecordedGateFixtureAt(t *testing.T, root string) recordedGateFixture {
	t.Helper()
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
		"Reviewed snapshot: `" + recordedGateBriefingID + "` at `" + recordedGateDigest + "`.\n\n" +
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

func recordedGateCommittedBeforeDispatch(t *testing.T, fixture recordedGateFixture, close, consumed, dispatchHead string) bool {
	t.Helper()
	closed, spent := recordedGateEntityAt(t, fixture, close), recordedGateEntityAt(t, fixture, consumed)
	if close == "" || consumed == "" || dispatchHead == "" || close == consumed ||
		!strings.Contains(closed, "decision: approve") || !strings.Contains(closed, "state: pending") ||
		!strings.Contains(spent, "status: handoff") || !strings.Contains(spent, "state: consumed") {
		return false
	}
	return exec.Command("git", "-C", fixture.stateRoot, "merge-base", "--is-ancestor", close, consumed).Run() == nil &&
		exec.Command("git", "-C", fixture.stateRoot, "merge-base", "--is-ancestor", consumed, dispatchHead).Run() == nil
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

func probeRecordedGateCapability(binary string, cache map[string]bool) error {
	target, _ := exec.LookPath(binary)
	target, _ = filepath.EvalSymlinks(target)
	body, _ := os.ReadFile(target)
	key := fmt.Sprintf("%s:%x", target, sha256.Sum256(body))
	if cache[key] {
		return nil
	}
	missing := []string{}
	checks := []struct {
		argv  []string
		wants []string
	}{
		{[]string{"gate", "--help"}, []string{"record", "validate", "eligibility", "consume", "--briefing", "--result", "--association", "--decision", "--actor", "--directive"}},
	}
	for _, check := range checks {
		out, err := exec.Command(target, check.argv...).CombinedOutput()
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
	cache[key] = true
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
