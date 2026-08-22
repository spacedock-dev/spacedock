//go:build live

// ABOUTME: Live drive of the present-gate decision-request rendering — runs a real first
// ABOUTME: officer over the halted-worker fixture and grades its final message.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeDecisionRequestFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range map[string]string{
		"README.md":  decisionRequestWorkflow(),
		"reading.md": decisionRequestEntity(),
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

// TestLiveDecisionRequest is the only check here that can establish the template
// works; the offline table can only establish that its graders were not
// loosened. Grading runs through gradeDecisionRequest, the same entry point that
// table pins, so a grader relaxed to turn this green breaks a case there.
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

	cmd := exec.Command(bin, "codex", "--plugin-dir", repo, "--skip-compat-check",
		decisionRequestPrompt(fixture),
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
