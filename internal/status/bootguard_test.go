// ABOUTME: AC-3 verdict matrix for the boot-at-compaction-boundary guard — the
// ABOUTME: durable-state-only table using the live-captured compact_boundary records as constants.
package status

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// incidentCompactBoundaryLine / spikeCompactBoundaryLine are verbatim byte-for-byte
// copies of the two live-captured Claude Code `compact_boundary` transcript
// records from force-boot-at-compaction-boundary/ideation-spike-evidence.md §1
// and §3 — the incident that motivated this guard, and the ideation spike that
// proved the boundary is durable, host-recorded state. Nothing here is a
// synthesized record shape; the capture is the oracle.
const incidentCompactBoundaryLine = `{"parentUuid":null,"logicalParentUuid":"6b99d305-cee5-406c-ae10-3bd1ec8d4eaf","isSidechain":false,"type":"system","subtype":"compact_boundary","content":"Conversation compacted","level":"info","compactMetadata":{"trigger":"manual","preTokens":801920,"postTokens":19650,"cumulativeDroppedTokens":782270,"durationMs":121692,"preCompactDiscoveredTools":["SendMessage","TaskList","TaskOutput","TaskStop"],"preservedSegment":{"headUuid":"4645e36d-8657-4b0e-9d4d-82667f3e4409","anchorUuid":"5c42641c-3b0b-41b3-97b7-822914a30b8e","tailUuid":"6b99d305-cee5-406c-ae10-3bd1ec8d4eaf"},"preservedMessages":{"anchorUuid":"5c42641c-3b0b-41b3-97b7-822914a30b8e","uuids":["4645e36d-8657-4b0e-9d4d-82667f3e4409","6b99d305-cee5-406c-ae10-3bd1ec8d4eaf"],"allUuids":["4645e36d-8657-4b0e-9d4d-82667f3e4409","8db20dca-608b-4f3e-860a-e292e4d1bf9b","6b99d305-cee5-406c-ae10-3bd1ec8d4eaf"]}},"uuid":"7b52ac75-04d8-4608-b7e1-07da94f5146c","timestamp":"2026-08-18T18:59:30.622Z","userType":"external","entrypoint":"cli","cwd":"/Users/clkao/git/spacedock-research/spacedock-v1","sessionId":"fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1","version":"2.1.226","gitBranch":"main","slug":"curried-brewing-kahn"}`

// incidentBoundaryTimestamp is the exact instant incidentCompactBoundaryLine
// carries, parsed once so test cases can express booted_at relative to it
// without re-deriving the literal.
var incidentBoundaryTimestamp = func() time.Time {
	ts, err := time.Parse(time.RFC3339Nano, "2026-08-18T18:59:30.622Z")
	if err != nil {
		panic(err)
	}
	return ts
}()

const bootGuardTestSessionID = "fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1"

func bootGuardWriteTranscript(t *testing.T, dir string, lines ...string) string {
	t.Helper()
	path := filepath.Join(dir, "transcript.jsonl")
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func bootGuardWriteReceipt(t *testing.T, gitRoot, bootedAt, transcript string) {
	t.Helper()
	dir := filepath.Join(gitRoot, ".spacedock", "boot")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	line := bootedAt
	if transcript != "" {
		line += " " + transcript
	}
	if err := os.WriteFile(filepath.Join(dir, bootGuardTestSessionID), []byte(line+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestBootGuardVerdictMatrix is AC-3's table: the verdict derives from durable
// state alone. Every row that is not one of the two provable-stale conditions
// (no receipt at all; a boundary newer than booted_at) fails OPEN, per the
// declared fail-direction — a false refusal is the risk this guard must never
// take on an unproven condition.
func TestBootGuardVerdictMatrix(t *testing.T) {
	fresh := incidentBoundaryTimestamp.Add(time.Hour).UTC().Format(time.RFC3339)
	stale := incidentBoundaryTimestamp.Add(-time.Hour).UTC().Format(time.RFC3339)

	for _, tc := range []struct {
		name             string
		sessionID        string
		setup            func(t *testing.T, gitRoot string)
		wantStale        bool
		wantReasonSubstr string
	}{
		{
			name:             "no receipt at all: never booted this session",
			sessionID:        bootGuardTestSessionID,
			setup:            func(t *testing.T, gitRoot string) {},
			wantStale:        true,
			wantReasonSubstr: "no boot receipt",
		},
		{
			name:      "fresh: booted after the latest boundary",
			sessionID: bootGuardTestSessionID,
			setup: func(t *testing.T, gitRoot string) {
				transcript := bootGuardWriteTranscript(t, gitRoot, incidentCompactBoundaryLine)
				bootGuardWriteReceipt(t, gitRoot, fresh, transcript)
			},
			wantStale: false,
		},
		{
			name:      "stale: compacted after the last boot",
			sessionID: bootGuardTestSessionID,
			setup: func(t *testing.T, gitRoot string) {
				transcript := bootGuardWriteTranscript(t, gitRoot, incidentCompactBoundaryLine)
				bootGuardWriteReceipt(t, gitRoot, stale, transcript)
			},
			wantStale:        true,
			wantReasonSubstr: "compacted after its last boot",
		},
		{
			name:      "missing transcript: fails open, not closed",
			sessionID: bootGuardTestSessionID,
			setup: func(t *testing.T, gitRoot string) {
				bootGuardWriteReceipt(t, gitRoot, stale, filepath.Join(gitRoot, "does-not-exist.jsonl"))
			},
			wantStale: false,
		},
		{
			name:      "no resolvable session identity: no-op",
			sessionID: "",
			setup:     func(t *testing.T, gitRoot string) {}, // never consulted: no session id to resolve
			wantStale: false,
		},
		{
			name:      "malformed transcript lines are skipped, not fatal",
			sessionID: bootGuardTestSessionID,
			setup: func(t *testing.T, gitRoot string) {
				transcript := bootGuardWriteTranscript(t, gitRoot,
					`not json at all`, `{"type":"user","message":"hi"}`, incidentCompactBoundaryLine, `{broken`)
				bootGuardWriteReceipt(t, gitRoot, stale, transcript)
			},
			wantStale:        true,
			wantReasonSubstr: "compacted after its last boot",
		},
		{
			name:      "malformed receipt record: fails open, not closed",
			sessionID: bootGuardTestSessionID,
			setup: func(t *testing.T, gitRoot string) {
				dir := filepath.Join(gitRoot, ".spacedock", "boot")
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, bootGuardTestSessionID), []byte("not-a-timestamp garbage\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantStale: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			tc.setup(t, gitRoot)
			e := env{}
			if tc.sessionID != "" {
				e["CLAUDE_CODE_SESSION_ID"] = tc.sessionID
			}
			stale, reason := bootGuardVerdict(e, gitRoot)
			if stale != tc.wantStale {
				t.Fatalf("stale = %v (reason %q), want stale=%v", stale, reason, tc.wantStale)
			}
			if tc.wantReasonSubstr != "" && !strings.Contains(reason, tc.wantReasonSubstr) {
				t.Fatalf("reason = %q, want substring %q", reason, tc.wantReasonSubstr)
			}
		})
	}
}

// TestBootGuardRefuseMessageNamesRemedy pins the exported entry point end to
// end: rawEnv/dir resolution through FindGitRoot, and the refusal text naming
// the exact recovery a stuck FO reads — re-run boot.
func TestBootGuardRefuseMessageNamesRemedy(t *testing.T) {
	gitRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(gitRoot, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	msg := BootGuardRefuse([]string{"CLAUDE_CODE_SESSION_ID=" + bootGuardTestSessionID}, gitRoot)
	if !strings.Contains(msg, "no boot receipt") || !strings.Contains(msg, "status --boot --identify --json") {
		t.Fatalf("refusal message = %q, want the condition and the status --boot remedy", msg)
	}
	if BootStaleExitCode != 4 {
		t.Fatalf("BootStaleExitCode = %d, want 4 (3 is already claimed by state commit / merge guard's rebase halt)", BootStaleExitCode)
	}
}

// TestBootWritesSessionReceipt proves --boot's declared side effect: a
// resolvable Claude session identity gets a one-line receipt at
// .spacedock/boot/{session_id} carrying `{booted_at} {transcript_path}` — the
// exact shape bootGuardVerdict reads back. This exercises the real
// claudeteam.TranscriptPath probe against genuine file I/O, not a stub.
func TestBootWritesSessionReceipt(t *testing.T) {
	root := t.TempDir()
	testgit.InitRepo(t, root)
	writeFile(t, filepath.Join(root, "README.md"), "---\ncommissioned-by: spacedock@1\nid-style: slug\nstages:\n  states:\n    - name: build\n      initial: true\n---\n# Receipt Workflow\n")

	home := t.TempDir()
	sessionID := "receipt-session-0000"
	transcriptDir := filepath.Join(home, ".claude", "projects", "-fake-project")
	if err := os.MkdirAll(transcriptDir, 0o755); err != nil {
		t.Fatal(err)
	}
	transcriptPath := filepath.Join(transcriptDir, sessionID+".jsonl")
	writeFile(t, transcriptPath, "{}\n")

	env := []string{"HOME=" + home, "CLAUDE_CODE_SESSION_ID=" + sessionID, "PATH=" + os.Getenv("PATH")}
	runner := &NativeRunner{TeamStateProbe: claudeteam.Probe, TranscriptProbe: claudeteam.TranscriptPath}
	var stdout, stderr strings.Builder
	before := time.Now().UTC()
	code, err := runner.Run(context.Background(), Request{
		Args: []string{"--workflow-dir", root, "--boot"}, Dir: root, Env: env,
		Stdout: &stdout, Stderr: &stderr,
	})
	if err != nil || code != 0 {
		t.Fatalf("--boot exit=%d err=%v stderr=%q", code, err, stderr.String())
	}

	raw, rerr := os.ReadFile(filepath.Join(root, ".spacedock", "boot", sessionID))
	if rerr != nil {
		t.Fatalf("read receipt: %v", rerr)
	}
	fields := strings.Fields(string(raw))
	if len(fields) != 2 {
		t.Fatalf("receipt = %q, want exactly two fields (booted_at transcript_path)", raw)
	}
	bootedAt, perr := time.Parse(time.RFC3339, fields[0])
	if perr != nil || bootedAt.Before(before.Add(-time.Minute)) {
		t.Fatalf("receipt booted_at = %q, want a recent RFC3339 timestamp: %v", fields[0], perr)
	}
	if fields[1] != transcriptPath {
		t.Fatalf("receipt transcript = %q, want %q", fields[1], transcriptPath)
	}
}
