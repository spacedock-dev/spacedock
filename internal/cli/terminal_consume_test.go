// ABOUTME: Terminal-consume/delivery boundary, driven through the real CLI:
// ABOUTME: consume routes (never spends) terminal approvals; merge guard is the
// ABOUTME: sole terminal consumer (envelope | --rework); retryable trouble moves
// ABOUTME: nothing. AC-1/AC-2/AC-3 of resolution-consume-terminal-before-delivery.
package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// terminalWorkflowOpts shapes the terminal-boundary CLI fixture.
type terminalWorkflowOpts struct {
	hook       string // merge hook mod name ("" = none registered)
	localMerge bool   // merge: local instead of the default merge: pr
	feedbackTo string // validation's declared feedback-to; "" = declared none
	extraStage bool   // declare feedback-to as an undefined stage name
}

// terminalCLIWorkflow materializes an inline workflow
// implementation (initial) -> validation (gate) -> done (terminal) with one
// entity at validation, git-initialized (the ceremony's mutations resolve
// git_root). It drives the REAL gate prepare/record verbs so the approval is a
// recorded, digest-bound approval — the hole's shape, not a mock.
func terminalCLIWorkflow(t *testing.T, opts terminalWorkflowOpts) (root, entity string) {
	t.Helper()
	root = t.TempDir()
	var mergeLine string
	if opts.localMerge {
		mergeLine = "merge: local\n"
	}
	feedback := ""
	if opts.feedbackTo != "" {
		feedback = "      feedback-to: " + opts.feedbackTo + "\n"
	}
	if opts.extraStage && opts.feedbackTo == "" {
		feedback = "      feedback-to: neverland\n"
	}
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\n"+mergeLine+"stages:\n  states:\n"+
		"    - name: implementation\n      initial: true\n"+
		"    - name: validation\n      gate: true\n"+feedback+
		"    - name: done\n      terminal: true\n---\n# Workflow\n")
	if opts.hook != "" {
		writeFileWithDirs(t, filepath.Join(root, "_mods", opts.hook+".md"), "---\nname: "+opts.hook+"\ndescription: stub merge hook.\n---\n\n# "+opts.hook+"\n\n## Hook: merge\n\n(stub — registration only)\n")
	}
	writeFile(t, filepath.Join(root, "gate-review.md"), "# Review\n")
	entity = filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nid: task\nstatus: validation\ntitle: Task\n---\n# Task\n")
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "spacedock-test"},
		{"add", "-A"},
		{"commit", "-q", "-m", "seed"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return root, entity
}

func terminalInvoke(t *testing.T, root string, args ...string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(context.Background(), args, nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
	return code, out.String(), errOut.String()
}

// approvedTerminalGate runs the real gate prepare + record against the
// terminalCLIWorkflow fixture so the entity carries a binding approved
// application targeting done.
func approvedTerminalGate(t *testing.T, root string) {
	t.Helper()
	if code, out, errOut := terminalInvoke(t, root,
		"gate", "prepare", "task",
		"--question", "Advance?",
		"--artifact", filepath.Join(root, "gate-review.md"),
		"--summary", "Terminal boundary fixture.",
		"--workflow-dir", root,
	); code != 0 {
		t.Fatalf("prepare exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := terminalInvoke(t, root,
		"gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", root,
	); code != 0 {
		t.Fatalf("record exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
}

func entityFields(t *testing.T, entity string) map[string]string {
	t.Helper()
	return status.ParseFrontmatter(entity)
}

func gateApplicationStates(t *testing.T, entity string) []string {
	t.Helper()
	doc, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	var states []string
	for ri := range doc.Records {
		for ai := range doc.Records[ri].Attempts {
			app := doc.Records[ri].Attempts[ai].Application
			if app == nil {
				states = append(states, "<nil>")
			} else {
				states = append(states, app.State)
			}
		}
	}
	return states
}

// TestTerminalDeliveryFailureReworkRoundTrip is AC-1: approval recorded ->
// consume routes without spending -> merge guard arms; delivery fails beyond
// retry -> --rework supersedes through the declared feedback-to with delivery
// state cleared; rework, re-enter, fresh approval, delivery proven ->
// merge guard --verdict passed finalizes status+verdict+completed+spend.
func TestTerminalDeliveryFailureReworkRoundTrip(t *testing.T) {
	root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
	approvedTerminalGate(t, root)

	// Consume routes, does not spend.
	code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "eligible=true") || !strings.Contains(out, "consumed=false") ||
		!strings.Contains(out, "target-stage=done") || !strings.Contains(out, "route=approved-awaiting-merge") {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	fields := entityFields(t, entity)
	if strings.TrimSpace(fields["status"]) != "validation" {
		t.Fatalf("terminal consume moved status: %q", fields["status"])
	}
	if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"pending"}) {
		t.Fatalf("terminal consume application states = %v, want [pending]", got)
	}

	// merge guard discovers and arms the registered pr-merge hook at delivery
	// time (consume armed nothing).
	code, out, errOut = terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "armed") || !strings.Contains(out, "merge:pr-merge") {
		t.Fatalf("merge guard arm exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	fields = entityFields(t, entity)
	if strings.TrimSpace(fields["mod-block"]) != "merge:pr-merge" {
		t.Fatalf("merge guard did not arm the registered hook: mod-block=%q", fields["mod-block"])
	}

	// The hook opened a PR; delivery then fails irreversibly (closed unmerged).
	if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "pr=#7"); code != 0 {
		t.Fatalf("record open PR exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	code, out, errOut = terminalInvoke(t, root, "merge", "guard", "task", "--rework", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "reworked: task -> implementation") || !strings.Contains(out, "superseded") {
		t.Fatalf("merge guard --rework exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	fields = entityFields(t, entity)
	if strings.TrimSpace(fields["status"]) != "implementation" {
		t.Fatalf("--rework did not route through the declared feedback-to: status=%q", fields["status"])
	}
	if strings.TrimSpace(fields["mod-block"]) != "" || strings.TrimSpace(fields["pr"]) != "" {
		t.Fatalf("--rework left delivery state behind: mod-block=%q pr=%q", fields["mod-block"], fields["pr"])
	}
	if strings.TrimSpace(fields["verdict"]) != "" || strings.TrimSpace(fields["completed"]) != "" {
		t.Fatalf("--rework wrote terminal fields pre-delivery: verdict=%q completed=%q", fields["verdict"], fields["completed"])
	}
	if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"superseded"}) {
		t.Fatalf("--rework application states = %v, want [superseded]", got)
	}
	doc, _, err := gates.Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Current.Gate != "gate:task:validation" {
		t.Fatalf("--rework moved gates.current: %q", doc.Current.Gate)
	}

	// Rework; re-enter the gated stage through the normal lifecycle.
	if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "status=validation"); code != 0 {
		t.Fatalf("re-enter validation exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := terminalInvoke(t, root,
		"gate", "prepare", "task",
		"--question", "Advance?",
		"--artifact", filepath.Join(root, "gate-review.md"),
		"--summary", "Reworked.",
		"--workflow-dir", root,
	); code != 0 {
		t.Fatalf("re-entry prepare exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := terminalInvoke(t, root,
		"gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", root,
	); code != 0 {
		t.Fatalf("re-entry record exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	// Superseded authority is never re-spent: attempt-1 stays superseded; the
	// FRESH approval is a successor attempt.
	if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"superseded", "pending"}) {
		t.Fatalf("re-entry application states = %v, want [superseded pending]", got)
	}

	code, out, errOut = terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "consumed=false") || !strings.Contains(out, "route=approved-awaiting-merge") {
		t.Fatalf("re-entry consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}

	// Delivery proven: the merge sentinel records the landed PR; merge guard
	// finalizes — terminal status, verdict, completed, and the spend land.
	if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "pr=pr-merge:7"); code != 0 {
		t.Fatalf("record merge sentinel exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	code, out, errOut = terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "finalized: task -> done (verdict passed)") {
		t.Fatalf("envelope finalize exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	archived := filepath.Join(root, "_archive", "task.md")
	fields = entityFields(t, archived)
	if strings.TrimSpace(fields["status"]) != "done" || strings.TrimSpace(fields["verdict"]) != "passed" ||
		strings.TrimSpace(fields["completed"]) == "" {
		t.Fatalf("envelope fields = status:%q verdict:%q completed:%q", fields["status"], fields["verdict"], fields["completed"])
	}
	if got := gateApplicationStates(t, archived); !slices.Equal(got, []string{"superseded", "consumed"}) {
		t.Fatalf("final application states = %v, want [superseded consumed]", got)
	}
}

// TestTerminalSpendOnlyInDeliveryEnvelope is AC-2: across merge: pr, registered
// merge: local, and manual local merge with NO registration, consume leaves the
// approval pending and the delivery envelope is the only terminal entry —
// terminal status, the verdict/completed pair, and the consumed application move
// together and only there. Plus the --set refusal and idempotent re-consume.
func TestTerminalSpendOnlyInDeliveryEnvelope(t *testing.T) {
	classes := []struct {
		name     string
		opts     terminalWorkflowOpts
		sentinel string
	}{
		{"merge-pr", terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"}, "pr-merge:7"},
		{"merge-local-registered", terminalWorkflowOpts{hook: "local-merge", localMerge: true, feedbackTo: "implementation"}, "local-merge:abc123f"},
	}
	for _, tc := range classes {
		t.Run(tc.name, func(t *testing.T) {
			root, entity := terminalCLIWorkflow(t, tc.opts)
			approvedTerminalGate(t, root)
			code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root)
			if code != 0 || !strings.Contains(out, "consumed=false") || !strings.Contains(out, "route=approved-awaiting-merge") {
				t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
			}
			// Pre-delivery invariant: nothing terminal anywhere.
			fields := entityFields(t, entity)
			if strings.TrimSpace(fields["status"]) != "validation" || strings.TrimSpace(fields["verdict"]) != "" ||
				strings.TrimSpace(fields["completed"]) != "" {
				t.Fatalf("pre-delivery terminal fields leaked: %v", fields)
			}
			if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"pending"}) {
				t.Fatalf("pre-delivery application states = %v", got)
			}
			// Delivery proven → the envelope.
			if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "pr="+tc.sentinel, "--workflow-dir", root); code != 0 {
				t.Fatalf("record sentinel exit=%d stdout=%q stderr=%q", code, out, errOut)
			}
			code, out, errOut = terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
			if code != 0 || !strings.Contains(out, "finalized: task -> done") {
				t.Fatalf("finalize exit=%d stdout=%q stderr=%q", code, out, errOut)
			}
			fields = entityFields(t, filepath.Join(root, "_archive", "task.md"))
			terminalDone := strings.TrimSpace(fields["status"]) == "done"
			fullTerminal := strings.TrimSpace(fields["verdict"]) != "" && strings.TrimSpace(fields["completed"]) != "" &&
				slices.Equal(gateApplicationStates(t, filepath.Join(root, "_archive", "task.md")), []string{"consumed"})
			if terminalDone != fullTerminal {
				t.Fatalf("AC-2 coherence: status terminal=%t but full terminal envelope=%t (fields=%v)", terminalDone, fullTerminal, fields)
			}
		})
	}

	t.Run("manual-local-merge-no-registration", func(t *testing.T) {
		root, _ := terminalCLIWorkflow(t, terminalWorkflowOpts{localMerge: true, feedbackTo: "implementation"})
		approvedTerminalGate(t, root)
		if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 ||
			!strings.Contains(out, "consumed=false") {
			t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		// The captain merges locally FIRST; the merge commit exists before the
		// guard runs. Record its sentinel, then let the guard observe it.
		mergeSHA := manualMergeCommit(t, root)
		if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "pr=local-merge:"+mergeSHA, "--workflow-dir", root); code != 0 {
			t.Fatalf("record local merge sentinel exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		code, out, errOut := terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
		if code != 0 || !strings.Contains(out, "finalized: task -> done") {
			t.Fatalf("manual-merge finalize exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		fields := entityFields(t, filepath.Join(root, "_archive", "task.md"))
		if strings.TrimSpace(fields["status"]) != "done" || strings.TrimSpace(fields["verdict"]) != "passed" ||
			strings.TrimSpace(fields["completed"]) == "" {
			t.Fatalf("manual-merge envelope fields = %v", fields)
		}
		if got := gateApplicationStates(t, filepath.Join(root, "_archive", "task.md")); !slices.Equal(got, []string{"consumed"}) {
			t.Fatalf("manual-merge application states = %v, want [consumed]", got)
		}
	})

	t.Run("terminal-set-refused-while-approval-pending", func(t *testing.T) {
		root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
		approvedTerminalGate(t, root)
		if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 {
			t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		before, err := os.ReadFile(entity)
		if err != nil {
			t.Fatal(err)
		}
		// Non-forced terminal --set: refused, naming merge guard.
		code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "status=done")
		if code == 0 || !strings.Contains(errOut, "sole terminal consumer") || !strings.Contains(errOut, "merge guard task") {
			t.Fatalf("terminal --set must refuse naming merge guard: exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		after, err := os.ReadFile(entity)
		if err != nil || !bytes.Equal(before, after) {
			t.Fatalf("refused terminal --set changed the entity (read err=%v)", err)
		}
		// --force remains the uniform escape hatch: the status writes, the
		// authority is untouched (never force-spent).
		if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "status=done", "--force"); code != 0 {
			t.Fatalf("forced terminal --set exit=%d stdout=%q stderr=%q", code, out, errOut)
		}
		fields := entityFields(t, entity)
		if strings.TrimSpace(fields["status"]) != "done" {
			t.Fatalf("forced terminal --set status = %q", fields["status"])
		}
		if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"pending"}) {
			t.Fatalf("--force must not spend authority: applications = %v", got)
		}
	})

	t.Run("re-consume-is-idempotent-routing", func(t *testing.T) {
		root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
		approvedTerminalGate(t, root)
		for pass := 0; pass < 2; pass++ {
			code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root)
			if code != 0 || !strings.Contains(out, "consumed=false") || !strings.Contains(out, "route=approved-awaiting-merge") {
				t.Fatalf("consume pass %d exit=%d stdout=%q stderr=%q", pass, code, out, errOut)
			}
		}
		if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"pending"}) {
			t.Fatalf("repeated consume moved authority: applications = %v", got)
		}
	})
}

// TestTerminalRetryLeavesAuthorityPending is AC-3: guard invoked before delivery
// proof writes no authority/status/terminal-field change and reports waiting;
// once delivery lands the envelope spends exactly once, and a re-invocation is
// a clean refusal with no second write.
func TestTerminalRetryLeavesAuthorityPending(t *testing.T) {
	root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
	approvedTerminalGate(t, root)
	if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	// Arm, then an open-but-unmerged PR: guard reports waiting, writes nothing.
	if code, out, errOut := terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root); code != 0 ||
		!strings.Contains(out, "armed") {
		t.Fatalf("arm exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "pr=#7"); code != 0 {
		t.Fatalf("open PR exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	code, out, errOut := terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "blocked") {
		t.Fatalf("blocked exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	fields := entityFields(t, entity)
	if strings.TrimSpace(fields["status"]) != "validation" || strings.TrimSpace(fields["verdict"]) != "" ||
		strings.TrimSpace(fields["completed"]) != "" {
		t.Fatalf("retry window moved terminal fields: %v", fields)
	}
	if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"pending"}) {
		t.Fatalf("retry window moved authority: %v", got)
	}

	// The retryable failure clears: delivery proven; the envelope lands once.
	if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "pr=pr-merge:7"); code != 0 {
		t.Fatalf("sentinel exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	code, out, errOut = terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "finalized: task -> done") {
		t.Fatalf("finalize exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	archived := filepath.Join(root, "_archive", "task.md")
	finalBytes, err := os.ReadFile(archived)
	if err != nil {
		t.Fatalf("finalized entity not archived: %v", err)
	}
	// Re-invoked after the spend: clean refusal, no double write.
	code, out, errOut = terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
	if code == 0 {
		t.Fatalf("re-invoked guard after the spend must not exit 0: stdout=%q", out)
	}
	after, err := os.ReadFile(archived)
	if err != nil || !bytes.Equal(finalBytes, after) {
		t.Fatalf("re-invoked guard wrote again (read err=%v)", err)
	}
}

// TestTerminalReworkRefusalMatrix: --rework refuses without a pending
// terminal-target approval, and when the record stage declares no feedback-to,
// an undefined one, or a terminal one. Every refusal leaves the entity
// byte-identical.
func TestTerminalReworkRefusalMatrix(t *testing.T) {
	root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	code, out, errOut := terminalInvoke(t, root, "merge", "guard", "task", "--rework", "--workflow-dir", root)
	if code == 0 || !strings.Contains(errOut, "pending terminal-target approval") {
		t.Fatalf("--rework without a pending approval must refuse: exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if after, _ := os.ReadFile(entity); !bytes.Equal(before, after) {
		t.Fatal("refused --rework changed the entity")
	}

	for _, tc := range []struct {
		name string
		opts terminalWorkflowOpts
		want string
	}{
		{"no-feedback-to", terminalWorkflowOpts{hook: "pr-merge"}, "declares no feedback-to"},
		{"undefined-feedback-to", terminalWorkflowOpts{hook: "pr-merge", extraStage: true}, "not a stage defined"},
		{"terminal-feedback-to", terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "done"}, "terminal stage"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, entity := terminalCLIWorkflow(t, tc.opts)
			approvedTerminalGate(t, root)
			if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 {
				t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
			}
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			code, out, errOut := terminalInvoke(t, root, "merge", "guard", "task", "--rework", "--workflow-dir", root)
			if code == 0 || !strings.Contains(errOut, tc.want) {
				t.Fatalf("--rework with %s must refuse naming it: exit=%d stdout=%q stderr=%q", tc.name, code, out, errOut)
			}
			if after, _ := os.ReadFile(entity); !bytes.Equal(before, after) {
				t.Fatal("refused --rework changed the entity")
			}
			if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"pending"}) {
				t.Fatalf("refused --rework moved authority: %v", got)
			}
		})
	}

	// --rework and --verdict are disjoint outcomes.
	code, out, errOut = terminalInvoke(t, root, "merge", "guard", "task", "--rework", "--verdict", "passed", "--workflow-dir", root)
	if code == 0 || !strings.Contains(errOut, "not both") {
		t.Fatalf("--rework with --verdict must be refused: exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
}

// TestRoutedTerminalApprovalSurfacesExistingDisplay: the routed entity (pending
// terminal application, status unchanged) shows through the EXISTING pending-
// application display — no new readiness states.
func TestRoutedTerminalApprovalSurfacesExistingDisplay(t *testing.T) {
	root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
	approvedTerminalGate(t, root)
	if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	code, out, errOut := terminalInvoke(t, root, "gate", "eligibility", "task", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "application=advance/pending") || !strings.Contains(out, "condition=approved-pending") ||
		!strings.Contains(out, "eligible=true") {
		t.Fatalf("routed entity must surface through the existing pending-application display: exit=%d stdout=%q stderr=%q",
			code, out, errOut)
	}
	if fields := entityFields(t, entity); strings.TrimSpace(fields["status"]) != "validation" {
		t.Fatalf("routed entity status moved: %q", fields["status"])
	}
}

// manualMergeCommit performs a real --no-ff merge on the fixture's code repo and
// returns the merge-commit SHA — the delivery proof a manual local merge leaves
// behind BEFORE merge guard runs.
func manualMergeCommit(t *testing.T, root string) string {
	t.Helper()
	gitC := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	gitC("checkout", "-q", "-b", "task-branch")
	writeFile(t, filepath.Join(root, "task-work.txt"), "work\n")
	gitC("add", "task-work.txt")
	gitC("commit", "-q", "-m", "task work")
	gitC("checkout", "-q", "main")
	gitC("merge", "--no-ff", "-q", "-m", "merge task branch", "task-branch")
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("rev-parse: %v\n%s", err, out)
	}
	return strings.TrimSpace(string(out))
}
