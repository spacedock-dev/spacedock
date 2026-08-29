// ABOUTME: Terminal-consume/delivery boundary, driven through the real CLI:
// ABOUTME: consume routes (never spends) terminal approvals; merge guard is the
// ABOUTME: sole terminal consumer via the gates-owned locked operations
// ABOUTME: (finalize | --rework); retryable trouble moves nothing.
// ABOUTME: AC-1/AC-2/AC-3 of resolution-consume-terminal-before-delivery.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
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
	testgit.InitRepo(t, root, "-q")
	for _, args := range [][]string{
		{"add", "-A"},
		{"commit", "-q", "-m", "seed"},
	} {
		terminalRunGit(t, root, args...)
	}
	return root, entity
}

func terminalRunGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
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

// disturbRetainedGateRoom tampers the retained canonical Briefing of the room
// approvedTerminalGate published, so the frozen digest no longer matches. The
// Briefing is the whole retained authority of a one-file room.
func disturbRetainedGateRoom(t *testing.T, root string) {
	t.Helper()
	var briefingFile string
	if err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.Name() == "index.json" {
			briefingFile = p
		}
		return nil
	}); err != nil || briefingFile == "" {
		t.Fatalf("locate retained index.json: %v %q", err, briefingFile)
	}
	if err := os.WriteFile(briefingFile, []byte("{\"type\":\"tampered\"}\n"), 0o644); err != nil {
		t.Fatalf("disturb briefing room: %v", err)
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

// consumedNonterminalWorkflow creates a gate whose successor is nonterminal.
// The real consume command advances the entity and spends the approval.
func consumedNonterminalWorkflow(t *testing.T) (root, entity string) {
	t.Helper()
	root = t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "---\nid-style: slug\nmerge: local\nstages:\n  states:\n"+
		"    - name: review\n      initial: true\n      gate: true\n"+
		"    - name: implementation\n"+
		"    - name: done\n      terminal: true\n---\n# Workflow\n")
	writeFile(t, filepath.Join(root, "gate-review.md"), "# Review\n")
	entity = filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nid: task\nstatus: review\ntitle: Task\n---\n# Task\n")
	testgit.InitRepo(t, root, "-q")
	terminalRunGit(t, root, "add", "-A")
	terminalRunGit(t, root, "commit", "-q", "-m", "seed")
	if code, out, errOut := terminalInvoke(t, root,
		"gate", "prepare", "task", "--question", "Advance?", "--artifact", filepath.Join(root, "gate-review.md"),
		"--summary", "Enter implementation.", "--workflow-dir", root,
	); code != 0 {
		t.Fatalf("prepare exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := terminalInvoke(t, root,
		"gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", root,
	); code != 0 {
		t.Fatalf("record exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 || !strings.Contains(out, "consumed=true") {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	report := "\n## Stage Report: implementation\n\n- DONE: Complete the entered stage.\n  The worker commit records this report.\n\n### Summary\n\nThe implementation is complete.\n"
	writeFile(t, entity, string(body)+report)
	terminalRunGit(t, root, "add", "--", "task.md")
	terminalRunGit(t, root, "commit", "-q", "-m", "worker: complete implementation", "--", "task.md")
	return root, entity
}

// TestConsumedNonterminalApprovalAllowsOrdinaryTerminalFields pins AC-1 and
// AC-3. Removing the consumed-history classification makes the normal write
// fail. Requiring --force also makes this test fail because no such flag runs.
func TestConsumedNonterminalApprovalAllowsOrdinaryTerminalFields(t *testing.T) {
	root, entity := consumedNonterminalWorkflow(t)
	if fields := entityFields(t, entity); fields["status"] != "implementation" {
		t.Fatalf("consume status=%q, want implementation", fields["status"])
	}
	if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"consumed"}) {
		t.Fatalf("consume application states=%v, want [consumed]", got)
	}

	code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task",
		"status=done", "completed", "verdict=PASSED", "worktree=")
	if code != 0 || strings.Contains(out+errOut, "ineligible") {
		t.Fatalf("ordinary terminal fields exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	fields := entityFields(t, entity)
	if fields["status"] != "done" || fields["verdict"] != "PASSED" || strings.TrimSpace(fields["completed"]) == "" {
		t.Fatalf("terminal fields = status:%q verdict:%q completed:%q", fields["status"], fields["verdict"], fields["completed"])
	}
	if body, err := os.ReadFile(entity); err != nil || !bytes.Contains(body, []byte("## Stage Report: implementation")) {
		t.Fatalf("terminal write lost the worker report: read error=%v", err)
	}
}

// TestTerminalDeliveryFailureReworkRoundTrip is AC-1's value spine: approval
// recorded -> consume routes without spending (pending, approved-awaiting-merge)
// -> merge guard arms; delivery fails beyond retry -> --rework supersedes
// through the declared feedback-to with delivery state cleared; rework,
// re-enter, fresh approval, delivery proven -> merge guard --verdict passed
// lands terminal status+verdict+completed+spend, pr sentinel retained.
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
	if code != 0 || !strings.Contains(out, "reworked: task -> implementation") || !strings.Contains(out, "superseded") ||
		!strings.Contains(out, "re-enter validation as a successor attempt") {
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
	// finalizes — terminal status, verdict, completed, and the spend land in
	// the gates-owned locked write; the pr sentinel is RETAINED through
	// archive as durable delivery proof.
	if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "pr=pr-merge:7"); code != 0 {
		t.Fatalf("record merge sentinel exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	code, out, errOut = terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "finalized: task -> done (verdict passed)") {
		t.Fatalf("finalize exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	archived := filepath.Join(root, "_archive", "task.md")
	fields = entityFields(t, archived)
	// The stored verdict carries the schema's case (PASSED) while the signal line
	// asserted above stays lowercase — the CLI surface and the stored surface
	// differ by design, and the gates-owned locked write normalises like the
	// legacy path does.
	if strings.TrimSpace(fields["status"]) != "done" || strings.TrimSpace(fields["verdict"]) != "PASSED" ||
		strings.TrimSpace(fields["completed"]) == "" {
		t.Fatalf("finalized fields = status:%q verdict:%q completed:%q", fields["status"], fields["verdict"], fields["completed"])
	}
	if strings.TrimSpace(fields["pr"]) != "pr-merge:7" {
		t.Fatalf("pr merge sentinel must be retained through archive as delivery proof: pr=%q", fields["pr"])
	}
	if got := gateApplicationStates(t, archived); !slices.Equal(got, []string{"superseded", "consumed"}) {
		t.Fatalf("final application states = %v, want [superseded consumed]", got)
	}
}

// TestTerminalRetryLeavesAuthorityPending is AC-3: guard invoked before delivery
// proof writes no authority/status/terminal-field change and reports waiting;
// once delivery lands the locked write spends exactly once, and a re-invocation
// is a clean refusal with no second write.
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
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	code, out, errOut := terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
	if code != 0 || !strings.Contains(out, "blocked") {
		t.Fatalf("blocked exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	// Retryable trouble moves NOTHING: the complete entity bytes —
	// mod-block, pr, and every other field included — survive untouched.
	if after, _ := os.ReadFile(entity); !bytes.Equal(before, after) {
		t.Fatalf("blocked guard moved bytes in the retry window:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"pending"}) {
		t.Fatalf("retry window moved authority: %v", got)
	}

	// The retryable failure clears: delivery proven; the locked write lands once.
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

// TestTerminalSetRefusedWhileApprovalPending is the AC-2 sole-consumer refusal:
// a non-forced --set to terminal fails closed, naming merge guard, byte-clean;
// --force is the uniform escape hatch and never spends authority.
func TestTerminalSetRefusedWhileApprovalPending(t *testing.T) {
	root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
	approvedTerminalGate(t, root)
	if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
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
}

// TestTerminalSetRefusedOnUnclassifiableAuthority is the fail-closed half of
// the sole-consumer refusal, exercised HOOKLESS so no downstream merge-hook
// guard can mask it: a non-forced terminal --set whose gate authority cannot
// be classified (retained briefing room disturbed) is refused by default —
// never allowed through with authority left pending — byte-clean.
func TestTerminalSetRefusedOnUnclassifiableAuthority(t *testing.T) {
	root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{feedbackTo: "implementation"})
	approvedTerminalGate(t, root)
	if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	// Disturb the retained briefing room: digest-stale authority.
	disturbRetainedGateRoom(t, root)
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "status=done")
	if code == 0 || !strings.Contains(errOut, "cannot be classified") || !strings.Contains(errOut, "sole terminal consumer") {
		t.Fatalf("terminal --set with digest-stale authority must refuse fail-closed: exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if after, _ := os.ReadFile(entity); !bytes.Equal(before, after) {
		t.Fatal("digest-stale terminal --set refusal changed the entity")
	}
	if got := gateApplicationStates(t, entity); !slices.Equal(got, []string{"pending"}) {
		t.Fatalf("digest-stale terminal --set refusal moved authority: %v", got)
	}
}

// TestMergeGuardRefusesDigestStaleAuthorityByteClean pins the fail-closed
// writer selection: a REAL pending terminal application whose retained
// authority became unreadable (disturbed briefing room) must refuse the
// finalization byte-clean — never terminalize with authority left pending.
// The guard is ARMED before the tamper so the refusal must occur BEFORE any
// mod-block clear: a clear-before-classify regression strips the mod-block and
// turns the byte-clean assertion red.
func TestMergeGuardRefusesDigestStaleAuthorityByteClean(t *testing.T) {
	root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
	approvedTerminalGate(t, root)
	if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	// Arm the merge guard FIRST (mod-block=merge:pr-merge, as the other armed
	// legs do) and record the merged-PR sentinel.
	if code, out, errOut := terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root); code != 0 ||
		!strings.Contains(out, "armed") || !strings.Contains(out, "merge:pr-merge") {
		t.Fatalf("merge guard arm exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := terminalInvoke(t, root, "status", "--workflow-dir", root, "--set", "task", "pr=pr-merge:7"); code != 0 {
		t.Fatalf("record merge sentinel exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	// THEN disturb the retained briefing room recorded at prepare, so its
	// frozen digest no longer matches.
	disturbRetainedGateRoom(t, root)
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(before, []byte("mod-block: merge:pr-merge")) || !bytes.Contains(before, []byte("pr: pr-merge:7")) {
		t.Fatalf("armed fixture expectation wrong — before bytes must carry the mod-block and pr sentinel:\n%s", before)
	}
	code, out, _ := terminalInvoke(t, root, "merge", "guard", "task", "--verdict", "passed", "--workflow-dir", root)
	if code == 0 {
		t.Fatalf("merge guard must refuse digest-stale authority: stdout=%q", out)
	}
	if after, _ := os.ReadFile(entity); !bytes.Equal(before, after) {
		t.Fatalf("digest-stale refusal changed the entity (mod-block and pr sentinel must survive byte-clean):\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if _, err := os.Stat(filepath.Join(root, "_archive", "task.md")); !os.IsNotExist(err) {
		t.Fatalf("digest-stale refusal must not archive: %v", err)
	}
}

// TestRoutedTerminalApprovalSurfacesExistingDisplay: the routed entity (pending
// terminal application, status unchanged) shows through the EXISTING
// gate-readiness display — no new readiness states.
func TestRoutedTerminalApprovalSurfacesExistingDisplay(t *testing.T) {
	root, entity := terminalCLIWorkflow(t, terminalWorkflowOpts{hook: "pr-merge", feedbackTo: "implementation"})
	approvedTerminalGate(t, root)
	if code, out, errOut := terminalInvoke(t, root, "gate", "consume", "task", "--workflow-dir", root); code != 0 {
		t.Fatalf("consume exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	code, out, errOut := terminalInvoke(t, root, "status", "--fields", "id,gate-readiness", "--json", "--workflow-dir", root)
	var statusEnvelope struct {
		Entities []map[string]string `json:"entities"`
	}
	decodeErr := json.Unmarshal([]byte(out), &statusEnvelope)
	if code != 0 || decodeErr != nil || len(statusEnvelope.Entities) != 1 ||
		statusEnvelope.Entities[0]["gate-readiness"] != "approved-awaiting-merge" {
		t.Fatalf("routed entity must surface through status readiness: exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if fields := entityFields(t, entity); strings.TrimSpace(fields["status"]) != "validation" {
		t.Fatalf("routed entity status moved: %q", fields["status"])
	}
}
