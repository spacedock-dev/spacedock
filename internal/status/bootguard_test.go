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
			var stderr strings.Builder
			stale, reason := bootGuardVerdict(e, gitRoot, &stderr)
			if stale != tc.wantStale {
				t.Fatalf("stale = %v (reason %q, stderr %q), want stale=%v", stale, reason, stderr.String(), tc.wantStale)
			}
			if tc.wantReasonSubstr != "" && !strings.Contains(reason, tc.wantReasonSubstr) {
				t.Fatalf("reason = %q, want substring %q", reason, tc.wantReasonSubstr)
			}
		})
	}
}

// TestBootReceiptSubSecondPrecisionSurvivesRoundTrip is the deferred
// sub-second-boundary risk's regression: it drives the REAL write path
// (writeBootReceipt), not a hand-built receipt line, because the bug lived in
// Format(time.RFC3339) dropping the fractional second on WRITE — Go's
// time.Parse already accepts a fractional second regardless of which layout
// constant it is given, so a read-side-only test cannot exercise this. A boot
// landing 278ms after the incident's own captured compaction, in the SAME
// wall-clock second, must read back as fresh: under the old second-precision
// Format this booted_at truncated to a whole second and read as BEFORE the
// compaction, wrongly refusing.
// TestBootGuardVerdictHandlesTranscriptPathWithSpaces is the deferred
// receipt-line-splitting risk's regression: strings.Fields would split a
// transcript path containing a space into extra fields and silently truncate
// it to the segment before the space, so latestCompactBoundary would open a
// nonexistent path and fail open — no detection at all, for any project path
// containing a space. SplitN(line, " ", 2) keeps the whole path as the second
// field regardless of embedded spaces.
func TestBootGuardVerdictHandlesTranscriptPathWithSpaces(t *testing.T) {
	gitRoot := t.TempDir()
	spacedDir := filepath.Join(gitRoot, "path with spaces")
	if err := os.MkdirAll(spacedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	transcript := bootGuardWriteTranscript(t, spacedDir, incidentCompactBoundaryLine)
	stale := incidentBoundaryTimestamp.Add(-time.Hour).UTC().Format(time.RFC3339Nano)
	bootGuardWriteReceipt(t, gitRoot, stale, transcript)

	var stderr strings.Builder
	got, reason := bootGuardVerdict(env{"CLAUDE_CODE_SESSION_ID": bootGuardTestSessionID}, gitRoot, &stderr)
	if !got || !strings.Contains(reason, "compacted after its last boot") {
		t.Fatalf("stale=%v reason=%q (stderr %q), want a correctly-detected stale verdict even with a spaced transcript path", got, reason, stderr.String())
	}
}

func TestBootReceiptSubSecondPrecisionSurvivesRoundTrip(t *testing.T) {
	gitRoot := t.TempDir()
	transcript := bootGuardWriteTranscript(t, gitRoot, incidentCompactBoundaryLine)
	transcriptProbe := func(home, sessionID string) string { return transcript }
	bootedAt := incidentBoundaryTimestamp.Add(278 * time.Millisecond) // same second, genuinely after

	var writeStderr strings.Builder
	writeBootReceipt(env{"CLAUDE_CODE_SESSION_ID": bootGuardTestSessionID}, gitRoot, transcriptProbe, &writeStderr, bootedAt)
	if writeStderr.Len() != 0 {
		t.Fatalf("writeBootReceipt stderr = %q, want no warning on a clean write", writeStderr.String())
	}

	var verdictStderr strings.Builder
	stale, reason := bootGuardVerdict(env{"CLAUDE_CODE_SESSION_ID": bootGuardTestSessionID}, gitRoot, &verdictStderr)
	if stale {
		t.Fatalf("stale = true (reason %q), want fresh: booted_at (%s) is 278ms after the compaction (%s), within the same second",
			reason, bootedAt.Format(time.RFC3339Nano), incidentBoundaryTimestamp.Format(time.RFC3339Nano))
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
	var stderr strings.Builder
	msg := BootGuardRefuse([]string{"CLAUDE_CODE_SESSION_ID=" + bootGuardTestSessionID}, gitRoot, &stderr)
	if !strings.Contains(msg, "no boot receipt") || !strings.Contains(msg, "status --boot --identify --json") {
		t.Fatalf("refusal message = %q, want the condition and the status --boot remedy", msg)
	}
	if !strings.Contains(msg, bootReceiptPath(gitRoot, bootGuardTestSessionID)) {
		t.Fatalf("refusal message = %q, want the resolved receipt path so all three verbs are self-diagnosing", msg)
	}
	if BootStaleExitCode != 4 {
		t.Fatalf("BootStaleExitCode = %d, want 4 (3 is already claimed by state commit / merge guard's rebase halt)", BootStaleExitCode)
	}
}

// skipIfRoot skips a permission-bit reproduction: root ignores Unix permission
// bits, so the setup below would silently fail to reproduce the condition
// under a root-run CI lane rather than genuinely proving the fail-open path.
func skipIfRoot(t *testing.T) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("permission bits are not enforced for root")
	}
}

// TestBootGuardVerdictFailsOpenNotClosedOnUnreadableReceipt reproduces two of
// the three shapes validation found: an unreadable receipt directory (mode
// 000) and a receipt path that exists as a directory instead of a file. Both
// are read errors OTHER than fs.ErrNotExist, so both must fail OPEN with a
// stderr warning — the fix for the prior os.ReadFile handling that refused on
// EVERY error, which made an unwritable .spacedock/boot refuse permanently
// while `status --boot` exited 0 and said nothing.
func TestBootGuardVerdictFailsOpenNotClosedOnUnreadableReceipt(t *testing.T) {
	t.Run("receipt directory mode 000", func(t *testing.T) {
		skipIfRoot(t)
		gitRoot := t.TempDir()
		dir := filepath.Join(gitRoot, ".spacedock", "boot")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(dir, 0o000); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(dir, 0o755) }) // let t.TempDir()'s own cleanup remove it
		e := env{"CLAUDE_CODE_SESSION_ID": bootGuardTestSessionID}
		var stderr strings.Builder
		stale, reason := bootGuardVerdict(e, gitRoot, &stderr)
		if stale {
			t.Fatalf("stale = true (reason %q), want fail-open on an unreadable (not missing) receipt dir", reason)
		}
		if !strings.Contains(stderr.String(), "unreadable") {
			t.Fatalf("stderr = %q, want a warning naming the unreadable receipt", stderr.String())
		}
	})

	t.Run("receipt path exists as a directory", func(t *testing.T) {
		gitRoot := t.TempDir()
		receiptAsDir := filepath.Join(gitRoot, ".spacedock", "boot", bootGuardTestSessionID)
		if err := os.MkdirAll(receiptAsDir, 0o755); err != nil {
			t.Fatal(err)
		}
		e := env{"CLAUDE_CODE_SESSION_ID": bootGuardTestSessionID}
		var stderr strings.Builder
		stale, reason := bootGuardVerdict(e, gitRoot, &stderr)
		if stale {
			t.Fatalf("stale = true (reason %q), want fail-open when the receipt path is a directory", reason)
		}
		if !strings.Contains(stderr.String(), "unreadable") {
			t.Fatalf("stderr = %q, want a warning naming the unreadable receipt", stderr.String())
		}
	})
}

// TestBootWriteSurfacesFailureOnReadOnlyProjectRoot reproduces validation's
// third shape: a read-only project root (the read-only-mount / full-disk /
// root-owned case), where the receipt can never be created. `status --boot`
// must say so on stderr rather than exiting 0 silently — the write's only
// failure signal before this fix — and the subsequent guard refusal names the
// same expected receipt path, so an operator sees one consistent diagnosis
// instead of a re-run loop with no explanation.
func TestBootWriteSurfacesFailureOnReadOnlyProjectRoot(t *testing.T) {
	skipIfRoot(t)
	gitRoot := t.TempDir()
	if err := os.Chmod(gitRoot, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(gitRoot, 0o755) })

	e := env{"CLAUDE_CODE_SESSION_ID": bootGuardTestSessionID}
	var writeStderr strings.Builder
	writeBootReceipt(e, gitRoot, nil, &writeStderr, time.Now())
	if !strings.Contains(writeStderr.String(), "could not create") {
		t.Fatalf("write stderr = %q, want a warning naming the failed create", writeStderr.String())
	}

	var readStderr strings.Builder
	stale, reason := bootGuardVerdict(e, gitRoot, &readStderr)
	if !stale || !strings.Contains(reason, "no boot receipt") || !strings.Contains(reason, bootReceiptPath(gitRoot, bootGuardTestSessionID)) {
		t.Fatalf("verdict after a failed write: stale=%v reason=%q, want a refusal naming the never-created receipt path", stale, reason)
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
	bootedAt, perr := time.Parse(time.RFC3339Nano, fields[0])
	if perr != nil || bootedAt.Before(before.Add(-time.Minute)) {
		t.Fatalf("receipt booted_at = %q, want a recent RFC3339Nano timestamp: %v", fields[0], perr)
	}
	if fields[1] != transcriptPath {
		t.Fatalf("receipt transcript = %q, want %q", fields[1], transcriptPath)
	}
}
