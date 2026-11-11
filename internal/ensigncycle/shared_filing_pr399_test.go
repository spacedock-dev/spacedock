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

// TestCommandFilesViaNewSlugTerminators pins the direct-launcher slug boundary
// against the recorded Haiku filing command (10-07 and 10-27 claude-live runs):
// a shell separator or closing quote directly after the slug still names the
// slug, while a longer slug sharing the prefix does not.
func TestCommandFilesViaNewSlugTerminators(t *testing.T) {
	for command, want := range map[string]bool{
		`printf -- 'x' | ${SPACEDOCK_BIN:-spacedock} new wire-the-thing; echo "exit=$?"`: true,
		`${SPACEDOCK_BIN:-spacedock} new wire-the-thing;echo`:                            true,
		`${SPACEDOCK_BIN:-spacedock} new wire-the-thing&& git status`:                    true,
		`${SPACEDOCK_BIN:-spacedock} new wire-the-thing|tee out`:                         true,
		`(${SPACEDOCK_BIN:-spacedock} new wire-the-thing)`:                               true,
		`bash -c "${SPACEDOCK_BIN:-spacedock} new wire-the-thing"`:                       true,
		`${SPACEDOCK_BIN:-spacedock} new wire-the-thing --workflow-dir . ;`:              true,
		`${SPACEDOCK_BIN:-spacedock} new wire-the-thing-two; echo`:                       false,
		`${SPACEDOCK_BIN:-spacedock} new wire-the-thingy`:                                false,
		// pi run 42916110241: `new` flags before the slug.
		`"${SPACEDOCK_BIN:-spacedock}" new --workflow-dir /tmp/wf wire-the-thing <<'ENTITY'`: true,
		`spacedock new --folder --workflow-dir=. wire-the-thing`:                             true,
		`spacedock new --id-seed "seed value" --id-actor 'fo' wire-the-thing`:                true,
		`spacedock new --id-seed wire-the-thing other-task`:                                  false,
		`spacedock new --workflow-dir /tmp/wire-the-thing other-task`:                        false,
	} {
		if got := commandFilesViaNew(command, filingSlug); got != want {
			t.Errorf("commandFilesViaNew(%q) = %v, want %v", command, got, want)
		}
	}
}

func TestCodexFilingCreateCountAcceptsNewFlags(t *testing.T) {
	if got := codexFilingCreateCount(`spacedock new --workflow-dir . wire-the-thing`, filingSlug); got != 1 {
		t.Fatalf("codexFilingCreateCount = %d, want 1", got)
	}
}
