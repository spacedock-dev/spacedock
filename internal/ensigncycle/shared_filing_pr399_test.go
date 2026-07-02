// ABOUTME: AC-2 (fo-opus-behavioral-robustness) — the real archived PR #399
// ABOUTME: attempt-1 opus filing stream, checked in as a must-pass fixture.
package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAssertClaudeFilingViaNewAcceptsRealPR399Attempt1Stream replays the actual
// captured opus filing stream from PR #399 (run 27841539602 attempt 1, artifact
// runtime-live-e2e-claude-live-claude-opus-4-8,
// claude-shared-scenarios/filing/claude-stream.jsonl, 102,607 bytes) through
// assertClaudeFilingViaNew end to end. This pins the REAL deviation shape — the
// filing command splits the launcher resolution and the create call across two
// lines (`cd {root}\nB=${SPACEDOCK_BIN:-spacedock}\n$B new wire-the-thing
// --workflow-dir . <<'EOF' …`) — rather than a synthetic single-line fixture, so
// a regression in commandFilesViaNew's var-capture branch (PR #433) is caught
// against the shape that actually broke the live opus lane, not an approximation
// of it. Verified during ideation (Spike A) and re-verified here as a permanent
// regression test: reverting to the pre-#433 commandFilesViaNew (`newInvocation`
// alone — same-line spacedock/SPACEDOCK_BIN token required) reds this exact
// fixture, since the real command's `$B new` line carries neither token
// literally.
func TestAssertClaudeFilingViaNewAcceptsRealPR399Attempt1Stream(t *testing.T) {
	stream, err := os.ReadFile(filepath.Join("testdata", "claude_live_filing_pr399_attempt1.stream.jsonl"))
	if err != nil {
		t.Fatalf("open real-stream fixture: %v", err)
	}
	if err := assertClaudeFilingViaNew(string(stream), filingSlug); err != nil {
		t.Errorf("assertClaudeFilingViaNew red the real PR #399 attempt-1 opus filing stream — want nil (the FO filed correctly via the var-capture launcher idiom): %v", err)
	}
}
