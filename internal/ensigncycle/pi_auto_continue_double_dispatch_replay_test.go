package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPiAutoContinueReplayDoubleDispatch is the correction-round replay for the
// Pi auto-continue single-root leg from GitHub Actions run 32749385148 (artifact
// runtime-live-e2e-pi-live, sha bb814f024). The FO dispatched the validation
// worker twice: the first completed via subagent_wait, then the FO hit a
// dispatch-build entity_path error, retried with the correct path, and
// re-dispatched. The assert's spawn counter correctly counted 2 (a real
// double-dispatch), but the old `spawns != 1` check rejected a legitimate
// retry. The fix tolerates `spawns` in {1, 2} when the rest of the lifecycle
// completes.
//
// This test reconstructs the full stream the assert sees (stdout + stderr +
// root session, matching piSharedLiveDriver.run) and feeds it through
// assertWorkerLifecycle with stage="validation". It grades GREEN after the
// spawns-tolerance fix. t.Skip keeps it hermetic when the artifact is absent.
func TestPiAutoContinueReplayDoubleDispatch(t *testing.T) {
	artifactDir := os.Getenv("SPACEDOCK_PI_AC_ARTIFACT_DIR")
	if artifactDir == "" {
		artifactDir = "/tmp/pi-live-10-art/live-artifacts/pi/pi-common/auto-continue-after-implementation--auto-continue/single-root"
	}
	stdout, err := os.ReadFile(filepath.Join(artifactDir, "pi-stdout.txt"))
	if err != nil {
		t.Skipf("artifact absent (pi-stdout.txt): %v", err)
	}
	stderr, err := os.ReadFile(filepath.Join(artifactDir, "pi-stderr.txt"))
	if err != nil {
		t.Skipf("artifact absent (pi-stderr.txt): %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(artifactDir, "sessions", "*.jsonl"))
	if err != nil || len(matches) != 1 {
		t.Skipf("artifact absent (root session): %v (%d matches)", err, len(matches))
	}
	rootSession, err := os.ReadFile(matches[0])
	if err != nil {
		t.Skipf("artifact absent (root session read): %v", err)
	}
	// Reconstruct exactly as piSharedLiveDriver.run() does:
	// stream = stdout + "\n" + (stderr + "\n" + rootSession)
	stream := string(stdout) + "\n" + string(stderr) + "\n" + string(rootSession)
	report := autoContinueGatedEndState(false)
	if err := assertWorkerLifecycle(stream, report, "validation", "gate prepare"); err != nil {
		t.Fatalf("double-dispatch replay graded RED after the spawns-tolerance fix: %v", err)
	}
}
