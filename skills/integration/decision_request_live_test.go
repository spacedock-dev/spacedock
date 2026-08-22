//go:build live

// ABOUTME: Live drive of the present-gate decision-request rendering — runs a real first
// ABOUTME: officer over a halted-worker fixture and grades its final message.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The fixture puts a worker at a stop-number halt with three options that all
// move the budget and none of which move the requirement, plus the fact the
// worker did not act on: one remaining deliverable serves a user who does not
// exist yet. Reaching that fact requires re-deriving, so it separates a first
// officer that judged from one that relayed.
const (
	decisionRequestWorkflow = `---
commissioned-by: spacedock@0.27.0-pre3
stages:
  states:
    - name: backlog
      initial: true
    - name: implementation
    - name: done
      terminal: true
---

# Probe workflow

### implementation

- **Inputs:** The approved shape and its declared stop numbers.
- **Outputs:** The implemented slice.
- **Good:** The slice lands inside its declared stop numbers.
- **Bad:** The slice crosses a stop number and continues anyway.
`

	decisionRequestEntity = `---
id: reading
title: Publish a document and hand out its link
status: implementation
---
# Publish a document and hand out its link

## Stage Report: implementation

- FAILED: Crossed a declared stop number and halted
  Slice 1 stands at 1087 added lines against a declared stop number of 900.
  Of those, 628 lines are tests and 451 lines are product code.
  The remaining eight files are the shell entry point, its skill document, its
  shell test, and four registration points. Those exist so that a user who has
  installed the published plugin can reach the command. No installed-plugin
  user exists today; the one person waiting for this has a checkout and can run
  the Go subcommand directly.

### Options I can offer

1. Raise the stop numbers to 1400.
2. Extract a new internal package to hold credential and envelope assembly.
3. Cut slice 1 in half and defer the expiry read.

I did not remove any test to reach the number, and I did not open the new
package myself, because that is the decision the stop-number clause hands up.
`

	// The prompt names the situation and nothing else. It does not say this is
	// not a gate, and it does not name the template: routing a halt to the
	// decision request is the contract's judgment, and a prompt that makes it
	// grades the prompt instead of the contract.
	decisionRequestPrompt = `Use $spacedock:present-gate. The explicit workflow directory is %s; ` +
		`pass it as --workflow-dir to every Spacedock helper. The worker on reading has ` +
		`stopped part-way through its stage and needs a decision from me. Handle it. ` +
		`Do not record a decision or mutate files.`
)

func writeDecisionRequestFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range map[string]string{
		"README.md":  decisionRequestWorkflow,
		"reading.md": decisionRequestEntity,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	for _, args := range [][]string{
		{"init", "-q"},
		{"-c", "user.name=Spacedock", "-c", "user.email=test@example.invalid", "add", "README.md", "reading.md"},
		{"-c", "user.name=Spacedock", "-c", "user.email=test@example.invalid", "commit", "-qm", "fixture"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return dir
}

// TestLiveDecisionRequest drives a real first officer and grades what it
// rendered. Grading runs through gradeDecisionRequest, the same entry point the
// offline table test pins, so a grader loosened to make this pass loosens the
// recorded fixtures with it.
func TestLiveDecisionRequest(t *testing.T) {
	bin := os.Getenv("SPACEDOCK_BIN")
	if bin == "" {
		t.Fatal("set SPACEDOCK_BIN to the current spacedock binary")
	}
	repo := os.Getenv("SPACEDOCK_REPO_ROOT")
	if repo == "" {
		out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
		if err != nil {
			t.Fatalf("resolve repo root: %v", err)
		}
		repo = strings.TrimSpace(string(out))
	}

	fixture := writeDecisionRequestFixture(t)
	finalPath := filepath.Join(fixture, "final.txt")
	prompt := strings.Replace(decisionRequestPrompt, "%s", fixture, 1)

	cmd := exec.Command(bin, "codex", "--plugin-dir", repo, "--skip-compat-check", prompt,
		"--", "exec", "--json", "--dangerously-bypass-approvals-and-sandbox",
		"--cd", fixture, "--output-last-message", finalPath)
	cmd.Dir = fixture
	stream, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("live run: %v\n%s", err, stream)
	}

	final, err := os.ReadFile(finalPath)
	if err != nil {
		t.Fatalf("read final message: %v", err)
	}
	if dir := os.Getenv("SPACEDOCK_LIVE_ARTIFACT_DIR"); dir != "" {
		_ = os.WriteFile(filepath.Join(dir, "decision-request-final.txt"), final, 0o644)
		_ = os.WriteFile(filepath.Join(dir, "decision-request-stream.jsonl"), stream, 0o644)
	}

	if failures := gradeDecisionRequest(string(final)); len(failures) > 0 {
		t.Fatalf("decision request failed grading %v\n\n%s", failures, final)
	}
}
