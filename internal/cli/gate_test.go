package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

func TestGateWithdrawCLIUsesExactGrammarAndImplicitFirstOfficerAttribution(t *testing.T) {
	workflow, state, artifact := gatePrepareCLIFixture(t)
	entity := filepath.Join(state, "task.md")
	if _, err := gates.Prepare(entity, gates.PrepareInput{
		WorkflowDir: workflow,
		Question:    "Advance?",
		Artifact:    artifact,
		Summary:     "candidate",
	}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{
		"gate", "withdraw", "task",
		"--reason", "Sprint re-scope replaced the candidate.",
		"--actor", "person:captain",
		"--workflow-dir", workflow,
	}, nil, filepath.Dir(workflow), nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 || !strings.Contains(errOut.String(), "gate withdraw accepts exactly one --reason") {
		t.Fatalf("actor refusal exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	after, _ := os.ReadFile(entity)
	if !bytes.Equal(before, after) {
		t.Fatal("usage refusal changed entity")
	}

	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{
		"gate", "withdraw", "task",
		"--reason", "Sprint re-scope replaced the candidate.",
		"--workflow-dir", workflow,
	}, nil, filepath.Dir(workflow), nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 || errOut.Len() != 0 {
		t.Fatalf("withdraw exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	want := "withdrawn gate=gate:task:validation attempt=gate-attempt:task-validation-1 state=withdrawn briefing=briefing:task:validation:attempt-1:revision-1\n"
	if out.String() != want {
		t.Fatalf("withdraw stdout=%q want=%q", out.String(), want)
	}
	doc, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	attempt := doc.Records[0].Attempts[0]
	if attempt.Withdrawal == nil || attempt.Withdrawal.By != "agent:first-officer" ||
		attempt.Withdrawal.Reason != "Sprint re-scope replaced the candidate." {
		t.Fatalf("withdrawal = %#v", attempt.Withdrawal)
	}
	if parsed, err := time.Parse(time.RFC3339Nano, attempt.Withdrawal.At); err != nil || parsed.Location() != time.UTC {
		t.Fatalf("withdrawal timestamp = %q (%v)", attempt.Withdrawal.At, err)
	}
	if attempt.Resolution != nil || attempt.Application != nil {
		t.Fatalf("withdrawal fabricated closure: %#v", attempt)
	}
}

func TestGatePrepareCLIPrintsExactRoomBindingAndCurrentV1HelpSurface(t *testing.T) {
	workflow, state, artifact := gatePrepareCLIFixture(t)
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{
		"gate", "prepare", "task",
		"--question", "Advance?",
		"--artifact", artifact,
		"--summary", "  Résumé — exact.  ",
		"--workflow-dir", workflow,
	}, nil, filepath.Dir(workflow), nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("prepare exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if len(lines) != 4 ||
		lines[0] != "room="+filepath.Join(state, "task", "review", "validation", "briefing-1") ||
		lines[1] != "briefing=briefing:task:validation:attempt-1:revision-1" ||
		!strings.HasPrefix(lines[2], "digest=sha256:") ||
		lines[3] != "state=open" {
		t.Fatalf("prepare stdout=%q", out.String())
	}
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "--help"}, nil, filepath.Dir(workflow), nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("help exit=%d stderr=%q", code, errOut.String())
	}
	// Published gate help is the intentional exact-text exception to the semantic
	// prose-oracle ban. Grammar additions/removals must change this fixture openly.
	wantHelp := "Usage: spacedock gate prepare <entity> --question TEXT --artifact REVIEW.md --summary TEXT [--reference FILE ...] [--workflow-dir DIR]\n" +
		"       spacedock gate withdraw <entity> --reason TEXT [--workflow-dir DIR]\n" +
		"       spacedock gate record <entity> --decision approve|revise|hold --actor ID [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]\n" +
		"       spacedock gate record <entity> --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl [--workflow-dir DIR]\n" +
		"       spacedock gate consume <entity> [--workflow-dir DIR]\n\n" +
		"On an approval whose target stage is terminal, consume spends nothing: it leaves the\n" +
		"application pending and reports route=approved-awaiting-merge. The terminal merge\n" +
		"ceremony (`spacedock merge guard <slug> --verdict passed|rejected`) is the sole terminal\n" +
		"consumer; `merge guard --rework` sends a failed delivery back through the declared\n" +
		"feedback-to (pending -> superseded, delivery state cleared).\n\n" +
		"`gate record --consume` is the captain-approve fast path: close, sync, consume, sync\n" +
		"in one call. `--consume` requires --decision approve and is rejected as a usage error\n" +
		"with --decision revise|hold. A delegated `--actor agent:first-officer` decision requires\n" +
		"`--conn-quote` (the grant verbatim) and `--conn-source` (where it was given); those flags\n" +
		"are refused with `--actor person:captain` or with `--round`. In a split-root workflow, a\n" +
		"successful close or consume ends with a machine-parseable `sync=.../phase=...` line;\n" +
		"branch on that final line plus the exit code, never on which prose lines printed.\n"
	if out.String() != wantHelp {
		t.Fatalf("gate help differs from the published contract:\n--- got ---\n%s--- want ---\n%s", out.String(), wantHelp)
	}
}

func TestGatePrepareCLIRejectsSummaryCardinalityAndEncodingBeforeMutation(t *testing.T) {
	for _, tc := range []struct {
		name    string
		summary []string
		want    string
	}{
		{"missing", nil, "gate prepare accepts --summary exactly once\n"},
		{"same repeated", []string{"same", "same"}, "gate prepare accepts --summary exactly once\n"},
		{"different repeated", []string{"one", "two"}, "gate prepare accepts --summary exactly once\n"},
		{"invalid utf8", []string{string([]byte{0xff})}, "--summary must be valid UTF-8\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflow, state, artifact := gatePrepareCLIFixture(t)
			entity := filepath.Join(state, "task.md")
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			args := []string{"gate", "prepare", "task", "--question", "Advance?", "--artifact", artifact}
			for _, summary := range tc.summary {
				args = append(args, "--summary", summary)
			}
			args = append(args, "--workflow-dir", workflow)
			var out, errOut bytes.Buffer
			code := run(context.Background(), args, nil, filepath.Dir(workflow), nil, &out, &errOut, &status.NativeRunner{}, nil)
			if code != 1 || errOut.String() != tc.want {
				t.Fatalf("exit=%d stdout=%q stderr=%q want stderr=%q", code, out.String(), errOut.String(), tc.want)
			}
			after, _ := os.ReadFile(entity)
			if !bytes.Equal(before, after) {
				t.Fatal("rejected summary changed entity")
			}
			if _, err := os.Stat(filepath.Join(state, "task", "review", "validation")); !os.IsNotExist(err) {
				t.Fatalf("rejected summary created a room: %v", err)
			}
		})
	}
}

func gatePrepareCLIFixture(t *testing.T) (workflow, state, artifact string) {
	t.Helper()
	root := t.TempDir()
	mainRoot := filepath.Join(root, "main")
	workflow = filepath.Join(mainRoot, "docs", "dev")
	if err := os.MkdirAll(workflow, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, mainRoot, "-q")
	writeFile(t, filepath.Join(workflow, "README.md"), "---\nid-style: slug\nstate: .state\nstages:\n  states:\n    - name: validation\n      initial: true\n      gate: true\n    - name: done\n      terminal: true\n---\n# Workflow\n")
	artifact = filepath.Join(mainRoot, "gate-review.md")
	writeFile(t, artifact, "# Review\n")
	git(t, mainRoot, "add", ".")
	git(t, mainRoot, "commit", "-q", "-m", "main fixture")

	state = filepath.Join(workflow, ".state")
	if err := os.MkdirAll(state, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, state, "-q")
	grandfatherFlatRooms(t, state, "task")
	writeFile(t, filepath.Join(state, "task.md"), "---\nid: task\nstatus: validation\ntitle: Task\n---\n# Task\n")
	git(t, state, "add", ".")
	git(t, state, "commit", "-q", "-m", "state fixture")
	// Mechanism 1 (implicit split-root sync in gate record/consume): a real
	// checkout is born on the state branch by `state init`/`new`, not `main` —
	// statesync preflight rightly refuses the mismatch. Spike-verified one-line
	// fixture alignment (AC-2's sole declared fixture change).
	git(t, state, "branch", "-M", "spacedock-state/dev")
	return workflow, state, artifact
}

func TestGateRoundRecordAndValidateCLI(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: implementation\n      initial: true\n---\n# Workflow\n")
	entityDir := filepath.Join(root, "task")
	if err := os.MkdirAll(entityDir, 0o755); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(entityDir, "index.md")
	writeFile(t, entity, "---\nid: task\nstatus: implementation\ntitle: Task\n---\n# Task\n")
	copyGateTestdata(t, filepath.Join(entityDir, "candidate.patch"), filepath.Join("advisory-round", "candidate.patch"))
	inputs := filepath.Join(entityDir, "inputs")
	if err := os.MkdirAll(inputs, 0o755); err != nil {
		t.Fatal(err)
	}
	briefing := filepath.Join(inputs, "briefing.json")
	log := filepath.Join(inputs, "briefing.review.jsonl")
	copyGateTestdata(t, briefing, filepath.Join("advisory-round", "briefing.json"))
	copyGateTestdata(t, log, filepath.Join("advisory-round", "briefing.review.jsonl"))
	var out, errOut bytes.Buffer
	invoke := func(args ...string) int {
		out.Reset()
		errOut.Reset()
		return run(context.Background(), args, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	}
	args := []string{"gate", "record", "task", "--workflow-dir", root, "--round", "implementation/1", "--briefing", briefing, "--log", log}
	if code := invoke(args...); code != 0 || !strings.Contains(out.String(), "round=round:task:implementation:1") || strings.Contains(out.String(), "triage=") {
		t.Fatalf("round record exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}

func TestGateRoundRejectsRemovedFeedbackCycleWithoutMutation(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: implementation\n      initial: true\n---\n# Workflow\n")
	entityDir := filepath.Join(root, "task")
	if err := os.MkdirAll(entityDir, 0o755); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(entityDir, "index.md")
	writeFile(t, entity, "---\nid: task\nstatus: implementation\ntitle: Task\n---\n# Task\n")
	inputs := filepath.Join(entityDir, "inputs")
	if err := os.MkdirAll(inputs, 0o755); err != nil {
		t.Fatal(err)
	}
	briefing := filepath.Join(inputs, "briefing.json")
	log := filepath.Join(inputs, "briefing.review.jsonl")
	copyGateTestdata(t, briefing, filepath.Join("advisory-round", "briefing.json"))
	copyGateTestdata(t, log, filepath.Join("advisory-round", "briefing.review.jsonl"))
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--round", "implementation/1", "--briefing", briefing, "--log", log, "--feedback-cycle", filepath.Join(inputs, "feedback-cycle.txt")}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 || !strings.Contains(errOut.String(), "unknown gate flag: --feedback-cycle") {
		t.Fatalf("removed flag exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("removed flag rejection changed entity bytes")
	}
}

// TestGateRoundRejectsConsumeFlagWithoutMutation pins finding 6 (roborev,
// branch_final): --round is an advisory correction-round publication with no
// close or consume attempt to sequence, and mechanism 1 deliberately excludes
// it from the implicit sync (Alternatives rejected 6) — so --consume combined
// with --round must be a usage error (exit 2) rather than silently ignored,
// and the round must not be published.
func TestGateRoundRejectsConsumeFlagWithoutMutation(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: implementation\n      initial: true\n---\n# Workflow\n")
	entityDir := filepath.Join(root, "task")
	if err := os.MkdirAll(entityDir, 0o755); err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(entityDir, "index.md")
	writeFile(t, entity, "---\nid: task\nstatus: implementation\ntitle: Task\n---\n# Task\n")
	copyGateTestdata(t, filepath.Join(entityDir, "candidate.patch"), filepath.Join("advisory-round", "candidate.patch"))
	inputs := filepath.Join(entityDir, "inputs")
	if err := os.MkdirAll(inputs, 0o755); err != nil {
		t.Fatal(err)
	}
	briefing := filepath.Join(inputs, "briefing.json")
	log := filepath.Join(inputs, "briefing.review.jsonl")
	copyGateTestdata(t, briefing, filepath.Join("advisory-round", "briefing.json"))
	copyGateTestdata(t, log, filepath.Join("advisory-round", "briefing.review.jsonl"))
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root,
		"--round", "implementation/1", "--briefing", briefing, "--log", log, "--consume"},
		nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 {
		t.Fatalf("--consume + --round exit=%d stdout=%q stderr=%q, want exit 2", code, out.String(), errOut.String())
	}
	if out.Len() != 0 {
		t.Fatalf("usage-error rejection wrote to stdout: %q", out.String())
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("--consume + --round usage-error rejection changed the entity")
	}
}

// TestGateRecordConnCitationGrammarMatrix pins AC-3: the two record shapes
// (a captain decision citing no grant, a delegated FO decision citing one) are
// disjoint by grammar. Every incoherent combination refuses at exit 2 with a
// byte-identical entity and no lock residue, before any mutation is attempted.
func TestGateRecordConnCitationGrammarMatrix(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"fo-missing-both", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence"}, "--conn-quote and --conn-source"},
		{"fo-missing-quote", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence", "--conn-source", "launch runbook"}, "--conn-quote and --conn-source"},
		{"fo-missing-source", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence", "--conn-quote", "you have the conn"}, "--conn-quote and --conn-source"},
		{"fo-blank-quote", []string{"--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence", "--conn-quote", " ", "--conn-source", "launch runbook"}, "--conn-quote and --conn-source"},
		{"captain-with-quote", []string{"--decision", "approve", "--actor", "person:captain", "--conn-quote", "you have the conn"}, "refused on a person:captain decision"},
		{"captain-with-source", []string{"--decision", "approve", "--actor", "person:captain", "--conn-source", "launch runbook"}, "refused on a person:captain decision"},
		{"captain-with-both", []string{"--decision", "approve", "--actor", "person:captain", "--conn-quote", "you have the conn", "--conn-source", "launch runbook"}, "refused on a person:captain decision"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, entity := semanticDecisionFixture(t)
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			var out, errOut bytes.Buffer
			args := append([]string{"gate", "record", "task", "--workflow-dir", root}, tc.args...)
			code := run(context.Background(), args, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
			if code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), tc.want) {
				t.Fatalf("%s exit=%d stdout=%q stderr=%q, want exit=2 empty stdout stderr containing %q", tc.name, code, out.String(), errOut.String(), tc.want)
			}
			after, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, after) {
				t.Fatalf("%s changed the entity", tc.name)
			}
			if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
				t.Fatalf("%s left lock residue: %v", tc.name, err)
			}
		})
	}

	t.Run("round-with-conn-flags", func(t *testing.T) {
		root := t.TempDir()
		writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: implementation\n      initial: true\n---\n# Workflow\n")
		entityDir := filepath.Join(root, "task")
		if err := os.MkdirAll(entityDir, 0o755); err != nil {
			t.Fatal(err)
		}
		entity := filepath.Join(entityDir, "index.md")
		writeFile(t, entity, "---\nid: task\nstatus: implementation\ntitle: Task\n---\n# Task\n")
		copyGateTestdata(t, filepath.Join(entityDir, "candidate.patch"), filepath.Join("advisory-round", "candidate.patch"))
		inputs := filepath.Join(entityDir, "inputs")
		if err := os.MkdirAll(inputs, 0o755); err != nil {
			t.Fatal(err)
		}
		briefing := filepath.Join(inputs, "briefing.json")
		log := filepath.Join(inputs, "briefing.review.jsonl")
		copyGateTestdata(t, briefing, filepath.Join("advisory-round", "briefing.json"))
		copyGateTestdata(t, log, filepath.Join("advisory-round", "briefing.review.jsonl"))
		before, err := os.ReadFile(entity)
		if err != nil {
			t.Fatal(err)
		}
		var out, errOut bytes.Buffer
		code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root,
			"--round", "implementation/1", "--briefing", briefing, "--log", log,
			"--conn-quote", "you have the conn", "--conn-source", "launch runbook"},
			nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
		if code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "--conn-quote and --conn-source are not valid with --round or --briefing") {
			t.Fatalf("--round + conn flags exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
		}
		after, err := os.ReadFile(entity)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(before, after) {
			t.Fatal("--round + conn flags usage-error rejection changed the entity")
		}
	})
}

func TestRemovedGateVerbsAreAbsentAndSideEffectFree(t *testing.T) {
	for _, verb := range []string{"review", "eligibility", "validate"} {
		t.Run(verb, func(t *testing.T) {
			root := t.TempDir()
			before, err := os.ReadDir(root)
			if err != nil {
				t.Fatal(err)
			}

			var out, errOut bytes.Buffer
			code := run(context.Background(), []string{"gate", verb, "task"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
			if code != 2 || !strings.Contains(errOut.String(), "unknown subcommand (want: prepare|withdraw|record|consume)") {
				t.Fatalf("gate %s exit=%d stdout=%q stderr=%q", verb, code, out.String(), errOut.String())
			}
			after, err := os.ReadDir(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(before) != len(after) {
				t.Fatalf("rejected %s verb changed working directory: before=%v after=%v", verb, before, after)
			}
		})
	}
}

func TestGateRecordRejectsNonRoundBriefing(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\ncommissioned-by: spacedock@1\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n---\n# Workflow\n")
	grandfatherFlatRooms(t, root, "task")
	writeFile(t, filepath.Join(root, "task.md"), "---\nstatus: validation\ntitle: Task\n---\n# Task\n")
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--briefing", "briefing.json", "--workflow-dir", root}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 || errOut.String() != "Error: gate record --briefing requires --round\n" {
		t.Fatalf("non-round Briefing exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}

func TestGateRequestLocatorCarriesArbitraryBriefingNameThroughRecordValidateAndEligibility(t *testing.T) {
	root, entity, room := unboundGateRoomFixture(t)
	originalBriefing := filepath.Join(room, "briefing.json")
	briefingBytes, err := os.ReadFile(originalBriefing)
	if err != nil {
		t.Fatal(err)
	}
	var briefing map[string]any
	if err := json.Unmarshal(briefingBytes, &briefing); err != nil {
		t.Fatal(err)
	}
	artifacts := briefing["artifacts"].([]any)
	artifacts[0].(map[string]any)["summary"] = "Exact canonical review summary"
	briefingBytes, err = json.MarshalIndent(briefing, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	briefingBytes = append(briefingBytes, '\n')
	nested := filepath.Join(room, "canonical", "decision-material.data")
	if err := os.MkdirAll(filepath.Dir(nested), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nested, briefingBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(originalBriefing); err != nil {
		t.Fatal(err)
	}
	digest, err := gates.CanonicalDigest(briefingBytes)
	if err != nil {
		t.Fatal(err)
	}
	requestPath := filepath.Join(room, "request.json")
	requestBytes, err := os.ReadFile(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	var request map[string]any
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		t.Fatal(err)
	}
	requestBriefing := request["briefing"].(map[string]any)
	requestBriefing["locator"] = "canonical/decision-material.data"
	requestBriefing["digest"] = digest
	requestBytes, err = json.MarshalIndent(request, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, requestPath, string(append(requestBytes, '\n')))

	var out, errOut bytes.Buffer
	invoke := func(args ...string) int {
		out.Reset()
		errOut.Reset()
		return run(context.Background(), args, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	}
	doc, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	attempt := &doc.Records[0].Attempts[0]
	requestDigest, err := gates.CanonicalDigest(requestBytes)
	if err != nil {
		t.Fatal(err)
	}
	entityBytes, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	entityBody := strings.Replace(string(entityBytes), "digest: '"+attempt.Briefing.Digest+"'", "digest: '"+digest+"'", 1)
	entityBody = strings.Replace(entityBody, "request-digest: '"+attempt.Briefing.RequestDigest+"'", "request-digest: '"+requestDigest+"'", 1)
	if entityBody == string(entityBytes) {
		t.Fatal("arbitrary-locator binding fixture was not updated")
	}
	writeFile(t, entity, entityBody)
	if code := invoke("gate", "record", "task", "--workflow-dir", root, "--decision", "approve", "--actor", "person:captain"); code != 0 {
		t.Fatalf("record arbitrary locator chat decision exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if code := invoke("gate", "consume", "task", "--workflow-dir", root); code != 0 ||
		!strings.Contains(out.String(), "consumed=false") || !strings.Contains(out.String(), "target-stage=done") ||
		!strings.Contains(out.String(), "route=approved-awaiting-merge") {
		t.Fatalf("consume arbitrary locator exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "briefing.json") ||
		!strings.Contains(string(body), "room-ref: ./review/validation/briefing-1") ||
		!strings.Contains(string(body), "status: validation") ||
		!strings.Contains(string(body), "state: pending") {
		t.Fatalf("binding inferred canonical basename, lost room, or spent the terminal-target approval:\n%s", body)
	}
}

func TestGatePreparedBriefingLocatorLifecycleAndRefusals(t *testing.T) {
	type fixture struct {
		workflow, entity, room, briefing string
	}
	invoke := func(t *testing.T, workflow string, args ...string) (int, string, string) {
		t.Helper()
		var out, errOut bytes.Buffer
		code := run(context.Background(), args, nil, filepath.Dir(workflow), nil, &out, &errOut, &status.NativeRunner{}, nil)
		return code, out.String(), errOut.String()
	}
	prepareClosed := func(t *testing.T) fixture {
		t.Helper()
		workflow, state, artifact := gatePrepareCLIFixture(t)
		entity := filepath.Join(state, "task.md")
		room := filepath.Join(state, "task", "review", "validation", "briefing-1")
		briefing := filepath.Join(room, "gate-briefing.json")
		if code, out, errOut := invoke(t, workflow,
			"gate", "prepare", "task",
			"--question", "Advance?",
			"--artifact", artifact,
			"--summary", "Exact prepared basename lifecycle.",
			"--workflow-dir", workflow,
		); code != 0 {
			t.Fatalf("prepare exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		if _, err := os.Stat(filepath.Join(room, "briefing.json")); !os.IsNotExist(err) {
			t.Fatalf("prepare unexpectedly published legacy briefing.json: %v", err)
		}
		if code, out, errOut := invoke(t, workflow,
			"gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", workflow,
		); code != 0 {
			t.Fatalf("close prepared gate exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		return fixture{workflow: workflow, entity: entity, room: room, briefing: briefing}
	}

	t.Run("prepared gate-briefing reaches consume", func(t *testing.T) {
		fixture := prepareClosed(t)
		if code, out, errOut := invoke(t, fixture.workflow,
			"gate", "consume", "task", "--workflow-dir", fixture.workflow,
		); code != 0 || !strings.Contains(out, "consumed=false") || !strings.Contains(out, "target-stage=done") ||
			!strings.Contains(out, "route=approved-awaiting-merge") {
			t.Fatalf("consume prepared gate exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		body, err := os.ReadFile(fixture.entity)
		if err != nil {
			t.Fatal(err)
		}
		// Terminal-target approvals route, they do not spend: merge guard's
		// delivery envelope is the only terminal consumer.
		if !strings.Contains(string(body), "status: validation") || !strings.Contains(string(body), "state: pending") {
			t.Fatalf("terminal-target consume must keep status and the pending application:\n%s", body)
		}
	})

	for _, tc := range []struct {
		name   string
		mutate func(*testing.T, fixture)
		wants  []string
	}{
		{
			name: "missing Briefing",
			mutate: func(t *testing.T, fixture fixture) {
				t.Helper()
				if err := os.Remove(fixture.briefing); err != nil {
					t.Fatal(err)
				}
			},
			wants: []string{"resolve canonical Briefing locator", "no such file"},
		},
		{
			name: "tampered Briefing",
			mutate: func(t *testing.T, fixture fixture) {
				t.Helper()
				body, err := os.ReadFile(fixture.briefing)
				if err != nil {
					t.Fatal(err)
				}
				changed := bytes.Replace(body,
					[]byte("Exact prepared basename lifecycle."),
					[]byte("Tampered prepared basename lifecycle."), 1)
				if bytes.Equal(changed, body) {
					t.Fatal("prepared Briefing summary fixture was not found")
				}
				if err := os.WriteFile(fixture.briefing, changed, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wants: []string{"bound canonical Briefing bytes", "frozen digest"},
		},
		{
			name: "tampered request locator",
			mutate: func(t *testing.T, fixture fixture) {
				t.Helper()
				requestPath := filepath.Join(fixture.room, "request.json")
				body, err := os.ReadFile(requestPath)
				if err != nil {
					t.Fatal(err)
				}
				changed := bytes.Replace(body,
					[]byte(`"locator": "gate-briefing.json"`),
					[]byte(`"locator": "missing-gate-briefing.json"`), 1)
				if bytes.Equal(changed, body) {
					t.Fatal("prepared request locator fixture was not found")
				}
				if err := os.WriteFile(requestPath, changed, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wants: []string{"retained request.json", "frozen digest"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := prepareClosed(t)
			before, err := os.ReadFile(fixture.entity)
			if err != nil {
				t.Fatal(err)
			}
			tc.mutate(t, fixture)
			for _, command := range []string{"consume"} {
				code, out, errOut := invoke(t, fixture.workflow,
					"gate", command, "task", "--workflow-dir", fixture.workflow)
				output := out + errOut
				if code == 0 {
					t.Fatalf("%s accepted %s: stdout=%q stderr=%q", command, tc.name, out, errOut)
				}
				for _, want := range tc.wants {
					if !strings.Contains(output, want) {
						t.Errorf("%s %s output missing actionable %q: %s", command, tc.name, want, output)
					}
				}
				if strings.Contains(output, "condition=ineligible") {
					t.Errorf("%s %s silently collapsed to condition=ineligible: %s", command, tc.name, output)
				}
				after, err := os.ReadFile(fixture.entity)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(after, before) {
					t.Fatalf("%s %s changed entity/application bytes", command, tc.name)
				}
				if _, err := os.Stat(fixture.entity + ".gates.lock"); !os.IsNotExist(err) {
					t.Fatalf("%s %s left lock residue: %v", command, tc.name, err)
				}
			}
		})
	}
}

func unboundGateRoomFixture(t *testing.T) (root, entity, room string) {
	t.Helper()
	root = t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n      gate: true\n    - name: done\n      terminal: true\n---\n# Workflow\n")
	entity = filepath.Join(root, "task.md")
	room = filepath.Join(root, "review", "validation", "briefing-1")
	if err := os.MkdirAll(room, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"briefing.json", "request.json"} {
		copyGateTestdata(t, filepath.Join(room, name), filepath.Join("gate-room", name))
	}
	briefingBytes, err := os.ReadFile(filepath.Join(room, "briefing.json"))
	if err != nil {
		t.Fatal(err)
	}
	briefingDigest, err := gates.CanonicalDigest(briefingBytes)
	if err != nil {
		t.Fatal(err)
	}
	requestBytes, err := os.ReadFile(filepath.Join(room, "request.json"))
	if err != nil {
		t.Fatal(err)
	}
	requestDigest, err := gates.CanonicalDigest(requestBytes)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, entity, "---\nstatus: validation\ntitle: Task\ngates:\n"+
		"  version: 1\n"+
		""+
		"  records:\n"+
		"    - id: gate:docs-dev:3k:validation\n"+
		"      stage: validation\n"+
		"      attempts:\n"+
		"        - id: gate-attempt:3k-validation-1\n"+
		"          briefing: {id: 'briefing:docs-dev:3k:validation:attempt-1:revision-1', digest: '"+briefingDigest+"', request-digest: '"+requestDigest+"', room-ref: ./review/validation/briefing-1}\n"+
		"---\n# Task\n")
	return root, entity, room
}

func TestGateRecordChatDecisionAndRejectsProvenanceAndOperationInterfaces(t *testing.T) {
	root, entity := semanticDecisionFixture(t)
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--decision", "approve", "--actor", "agent:first-officer", "--reason", "All retained ACs reproduced", "--conn-quote", "you have the conn toward the sprint goal", "--conn-source", "launch runbook for this headless session"}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("record decision exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"by: agent:first-officer", "decision: approve", "reason: All retained ACs reproduced", "quote: you have the conn toward the sprint goal", "source: launch runbook for this headless session"} {
		if !strings.Contains(string(body), want) {
			t.Fatalf("chat Resolution missing %q:\n%s", want, body)
		}
	}
	if strings.Contains(string(body), "adoption-note:") {
		t.Fatalf("new chat Resolution retained caller-controlled provenance:\n%s", body)
	}
	historical := strings.Replace(string(body),
		"reason: All retained ACs reproduced",
		"reason: All retained ACs reproduced\n                adoption-note: historical delegated context", 1)
	writeFile(t, entity, historical)
	if _, _, err := gates.Read(entity); err == nil || !strings.Contains(err.Error(), "adoption-note") {
		t.Fatalf("prototype adoption-note read error=%v, want unknown-field refusal", err)
	}

	for _, tc := range []struct {
		name, flag, value string
	}{
		{"exact-directive", "--directive", "you have the conn."},
		{"missing-period-directive", "--directive", "you have the conn"},
		{"directive-file", "--directive-file", "authority.txt"},
		{"operation", "--operation", "old.yml"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, entity := semanticDecisionFixture(t)
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			var out, errOut bytes.Buffer
			code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--decision", "approve", "--actor", "agent:first-officer", "--reason", "evidence", tc.flag, tc.value}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
			if code != 2 || !strings.Contains(errOut.String(), "unknown gate flag: "+tc.flag) {
				t.Fatalf("legacy %s exit=%d stdout=%q stderr=%q", tc.flag, code, out.String(), errOut.String())
			}
			after, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, before) {
				t.Fatalf("legacy %s changed entity", tc.flag)
			}
			if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
				t.Fatalf("legacy %s left lock residue: %v", tc.flag, err)
			}
		})
	}
}

func TestGateRecordRejectsProviderRoomBeforeMutation(t *testing.T) {
	root, entity := semanticDecisionFixture(t)
	room := filepath.Join(root, "task", "review", "validation", "briefing-1")
	if err := os.MkdirAll(filepath.Join(room, "provider"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(room, "provider", "result.json"), `{"decision":"approve"}`)
	beforeEntity, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	beforeRoom, err := os.ReadFile(filepath.Join(room, "provider", "result.json"))
	if err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{
		"gate", "record", "task", "--room", room, "--workflow-dir", root,
	}, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 2 || errOut.String() != "Error: unknown gate flag: --room\n" {
		t.Fatalf("provider room refusal exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	afterEntity, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	afterRoom, err := os.ReadFile(filepath.Join(room, "provider", "result.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(beforeEntity, afterEntity) || !bytes.Equal(beforeRoom, afterRoom) {
		t.Fatal("provider-only option refusal changed entity or room bytes")
	}
	if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
		t.Fatalf("provider-only option refusal reached the entity lock: %v", err)
	}
}

func TestGateRecordCLIRejectsIncoherentBriefingAndStageWithoutMutation(t *testing.T) {
	for _, tc := range []struct {
		name, currentStage, briefingID, stageFlags, want string
	}{
		{"cross-stage", "implementation", "briefing:task:validation:attempt-1:revision-1", "      gate: true\n", "Error: Briefing stage validation does not match current workflow stage implementation\n"},
		{"unqualified", "implementation", "briefing:legacy", "      gate: true\n", "Error: Briefing id briefing:legacy is not a canonical stage-qualified v1 identity\n"},
		{"malformed", "implementation", "briefing:task:implementation:attempt-0:revision-1", "      gate: true\n", "Error: Briefing id briefing:task:implementation:attempt-0:revision-1 is not a canonical stage-qualified v1 identity\n"},
		{"non-gated", "implementation", "briefing:task:implementation:attempt-1:revision-1", "", "Error: current workflow stage implementation is not an actionable gate:true stage\n"},
		{"terminal", "done", "briefing:task:done:attempt-1:revision-1", "      gate: true\n      terminal: true\n", "Error: current workflow stage done is not an actionable gate:true stage\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, entity := gateRecordCoherenceCLIFixture(t, tc.currentStage, tc.briefingID, tc.stageFlags)
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			args := []string{"gate", "record", "task", "--workflow-dir", root, "--decision", "hold", "--actor", "person:captain", "--reason", "wait"}
			var out, errOut bytes.Buffer
			code := run(context.Background(), args, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
			if code != 1 || out.Len() != 0 || errOut.String() != tc.want {
				t.Fatalf("exit=%d stdout=%q stderr=%q, want exit=1 empty stdout stderr=%q", code, out.String(), errOut.String(), tc.want)
			}
			after, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, before) || bytes.Contains(after, []byte("resolution:")) {
				t.Fatal("refused ordinary close changed entity or added Resolution")
			}
			if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
				t.Fatalf("refused ordinary close left lock residue: %v", err)
			}
		})
	}
}

func semanticDecisionFixture(t *testing.T) (root, entity string) {
	t.Helper()
	root = t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: validation\n      initial: true\n      gate: true\n    - name: done\n      terminal: true\n---\n# Workflow\n")
	entity = filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nstatus: validation\ngates:\n  version: 1\n  records:\n    - id: gate:docs-dev:3k:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          briefing:\n            id: briefing:docs-dev:3k:validation:attempt-1:revision-1\n            digest: sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac\n            room-ref: ./review/validation/briefing-1\ntitle: Task\n---\n# Task\n")
	briefing := filepath.Join(root, "review", "validation", "briefing-1", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	copyGateTestdata(t, briefing, "exact-validation-briefing.json")
	return root, entity
}

func gateRecordCoherenceCLIFixture(t *testing.T, currentStage, briefingID, stageFlags string) (root, entity string) {
	t.Helper()
	root = t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: "+currentStage+"\n"+stageFlags+"    - name: next\n---\n# Workflow\n")
	entity = filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nstatus: "+currentStage+"\ngates:\n  version: 1\n  records:\n    - id: gate:task:"+currentStage+"\n      stage: "+currentStage+"\n      attempts:\n        - id: gate-attempt:task-"+currentStage+"-1\n          briefing:\n            id: "+briefingID+"\n            digest: sha256:"+strings.Repeat("1", 64)+"\n            room-ref: ./review/"+currentStage+"/briefing-1\ntitle: Task\n---\n# Task\n")
	return root, entity
}

func copyGateTestdata(t *testing.T, destination, name string) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "gates", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, body, 0o644); err != nil {
		t.Fatal(err)
	}
}

func crossGateFixture(t *testing.T) (root, entity, briefing string) {
	t.Helper()
	root = t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n    - name: validation\n---\n# Workflow\n")
	entity = filepath.Join(root, "durable-gate-approval-pending-blockers.md")
	digest := func(char string) string { return "sha256:" + strings.Repeat(char, 64) }
	writeFile(t, entity, "---\nstatus: ideation\ngates:\n  version: 1\n  records:\n    - id: gate:docs-dev:3k:ideation\n      stage: ideation\n      attempts:\n        - id: gate-attempt:3k-ideation-9\n          briefing:\n            id: briefing:ideation:9\n            digest: "+digest("1")+"\n            room-ref: ./review/ideation/9\n          resolution:\n            type: Resolution\n            id: resolution:ideation:9\n            briefing: briefing:ideation:9\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: approve\n          application:\n            action: advance\n            target-stage: validation\n            state: pending\n            blockers: [{id: blocker:preserve-me, state: unsatisfied}]\n    - id: gate:docs-dev:3k:validation\n      stage: validation\n      attempts:\n        - id: gate-attempt:3k-validation-1\n          briefing:\n            id: briefing:validation:1\n            digest: "+digest("2")+"\n            room-ref: ./review/validation/1\n          resolution:\n            type: Resolution\n            id: resolution:validation:1\n            briefing: briefing:validation:1\n            by: person:captain\n            at: 2026-07-22T00:00:00Z\n            decision: revise\n            reason: Re-enter ideation.\nsprint: durable-decisions\ntitle: Task\n---\n# Task\n")
	briefing = filepath.Join(root, "review", "ideation", "briefing-18", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	briefingBytes, err := os.ReadFile(filepath.Join("..", "gates", "testdata", "revision-18.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(briefing, briefingBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	return root, entity, briefing
}

func fixtureSection(t *testing.T, body, start, end string) string {
	t.Helper()
	startAt := strings.Index(body, start)
	if startAt < 0 {
		t.Fatalf("fixture section start %q not found", start)
	}
	endAt := strings.Index(body[startAt+len(start):], end)
	if endAt < 0 {
		t.Fatalf("fixture section end %q not found", end)
	}
	return body[startAt : startAt+len(start)+endAt]
}

func outsideFixtureGates(t *testing.T, body string) string {
	t.Helper()
	start := strings.Index(body, "gates:\n")
	end := strings.Index(body, "sprint: durable-decisions\n")
	if start < 0 || end < 0 || end <= start {
		t.Fatal("fixture gates boundaries not found")
	}
	return body[:start] + body[end:]
}
