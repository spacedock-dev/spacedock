// ABOUTME: Boot guard at the compaction boundary — status --boot writes a
// ABOUTME: per-session receipt; gate record/consume and merge guard refuse when it is stale or absent.
package status

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// BootStaleExitCode is the distinct exit code the three authority verbs — gate
// record, gate consume, merge guard — return when BootGuardRefuse refuses
// (BOOT_STALE). Exit 3 is already state commit's / merge guard's own rebase-halt
// signal; 4 is otherwise unclaimed across the CLI.
const BootStaleExitCode = 4

// bootSessionEnvVar is the Claude Code per-session identity variable the guard
// resolves the running session from. It survives the compaction boundary
// unchanged (captured live: force-boot-at-compaction-boundary/
// ideation-spike-evidence.md §4), so it is the one durable handle the binary
// needs — no hook, no host callback. A host with no such variable (a bare
// terminal, Codex, Pi today) resolves no identity and the guard no-ops.
const bootSessionEnvVar = "CLAUDE_CODE_SESSION_ID"

// resolveBootSessionID reads the session id from e, refusing to resolve
// anything that could escape the receipt directory as a path component. An
// empty or unsafe value reports "" — every caller here treats "" as no
// resolvable identity and no-ops, exactly today's behavior on a host with no
// session env.
func resolveBootSessionID(e env) string {
	id := strings.TrimSpace(e.get(bootSessionEnvVar))
	if id == "" || id == "." || id == ".." || strings.ContainsAny(id, "/\\") {
		return ""
	}
	return id
}

// bootReceiptPath is the per-session receipt location: project-local host
// scratch (gitignored), never workflow or entity state.
func bootReceiptPath(gitRoot, sessionID string) string {
	return filepath.Join(gitRoot, ".spacedock", "boot", sessionID)
}

// writeBootReceipt writes the one-line session boot receipt described in
// docs/runtime-support.md's "Boot guard at the compaction boundary" section:
// `{booted_at} {transcript_path}` at .spacedock/boot/{session_id}. A no-op when
// the session id is unresolvable (workflow/entity state is untouched either
// way — this is boot's only side effect, and it never touches either). transcriptProbe
// resolves the session's transcript file (nil on a non-Claude host — internal/status
// carries no ~/.claude read itself; see claudeteam.TranscriptProbe); an unresolved
// transcript still writes a timestamp-only line, and the guard fails open on it
// (see bootGuardVerdict).
//
// A failed write is surfaced to stderr, not swallowed: `--boot` still exits 0
// (the read it primarily promises still succeeded), but a silently-lost receipt
// would make the guard's own remedy ("one cheap idempotent boot") permanently
// unable to clear — the operator needs to see WHY re-running boot never helps.
func writeBootReceipt(e env, gitRoot string, transcriptProbe claudeteam.TranscriptProbe, stderr io.Writer, now time.Time) {
	sessionID := resolveBootSessionID(e)
	if sessionID == "" {
		return
	}
	home := e.get("HOME")
	if home == "" {
		home = os.Getenv("HOME")
	}
	line := now.UTC().Format(time.RFC3339Nano)
	if transcriptProbe != nil {
		if transcript := transcriptProbe(home, sessionID); transcript != "" {
			line += " " + transcript
		}
	}
	dir := filepath.Join(gitRoot, ".spacedock", "boot")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(stderr, "spacedock: warning — could not create %s for the session boot receipt (%v); gate record, gate consume, and merge guard will refuse until this is fixed.\n", dir, err)
		return
	}
	receiptPath := bootReceiptPath(gitRoot, sessionID)
	if err := os.WriteFile(receiptPath, []byte(line+"\n"), 0o644); err != nil {
		fmt.Fprintf(stderr, "spacedock: warning — could not write boot receipt at %s (%v); gate record, gate consume, and merge guard will refuse until this is fixed.\n", receiptPath, err)
	}
}

// latestCompactBoundary scans a Claude Code session transcript (JSONL) for the
// latest `{"type":"system","subtype":"compact_boundary",...,"timestamp":...}`
// record — the durable, host-recorded compaction marker (captured verbatim in
// ideation-spike-evidence.md §1 and §3). A line that fails to parse as JSON, or
// parses but is not a compact_boundary record, is skipped rather than treated
// as an error; only total unreadability (missing file, no permission) reports
// found=false.
func latestCompactBoundary(path string) (latest time.Time, found bool) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		var rec struct {
			Type      string `json:"type"`
			Subtype   string `json:"subtype"`
			Timestamp string `json:"timestamp"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			continue
		}
		if rec.Type != "system" || rec.Subtype != "compact_boundary" {
			continue
		}
		ts, err := time.Parse(time.RFC3339Nano, rec.Timestamp)
		if err != nil {
			continue
		}
		if !found || ts.After(latest) {
			latest, found = ts, true
		}
	}
	return latest, found
}

// bootGuardVerdict reports whether the running session's boot is stale, per the
// declared fail-direction: no resolvable identity, no transcript, or a
// malformed receipt record all fail OPEN (stale=false) — the guard refuses only
// the two conditions it can prove from durable state: no receipt file at all
// for this session (never booted), or a compact_boundary record newer than the
// receipt's booted_at (compacted since the last boot).
//
// An unreadable receipt (permission denied, a read-only mount, the path
// existing as a directory) is NOT the same as a missing one: os.ErrNotExist is
// the only read error that refuses. Every other read error proves LESS than a
// missing file — a malformed record two lines below already fails open, so an
// unreadable one failing CLOSED would be backwards — and is reported to stderr
// as a warning rather than silently downgraded, since it is exactly the
// unrecoverable-loop shape (status --boot exits 0, the write keeps failing, the
// guard keeps refusing) this fix exists to close.
func bootGuardVerdict(e env, gitRoot string, stderr io.Writer) (stale bool, reason string) {
	sessionID := resolveBootSessionID(e)
	if sessionID == "" {
		return false, ""
	}
	receiptPath := bootReceiptPath(gitRoot, sessionID)
	raw, err := os.ReadFile(receiptPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return true, fmt.Sprintf("no boot receipt for this session (never booted, or booted in a different session) — expected at %s", receiptPath)
		}
		fmt.Fprintf(stderr, "spacedock: warning — boot receipt at %s is unreadable (%v); failing open until this is fixed (re-running boot will not clear it).\n", receiptPath, err)
		return false, ""
	}
	// SplitN, not Fields: a project path containing a space would otherwise
	// split into extra fields and silently disable the transcript half of the
	// receipt (Fields treats any run of whitespace as a separator).
	parts := strings.SplitN(strings.TrimRight(string(raw), "\n"), " ", 2)
	bootedAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return false, "" // malformed record: fail open
	}
	if len(parts) < 2 || parts[1] == "" {
		return false, "" // no transcript recorded: cannot prove staleness, fail open
	}
	latest, found := latestCompactBoundary(parts[1])
	if !found {
		return false, "" // transcript missing/unreadable/no boundary recorded: fail open
	}
	if latest.After(bootedAt) {
		return true, fmt.Sprintf("this session compacted after its last boot (receipt at %s)", receiptPath)
	}
	return false, ""
}

// BootGuardRefuse is the shared preflight `gate record`, `gate consume`, and
// `merge guard` each run before any mutation. It returns "" when the running
// session's boot is fresh, or when identity/transcript/receipt cannot prove
// staleness (fail open — see bootGuardVerdict, which may still print a
// non-fatal stderr warning on that path); otherwise it returns one stderr
// paragraph naming the condition, the receipt path, and the exact remedy.
// Callers gate on a non-empty return and exit BootStaleExitCode without
// mutating anything.
func BootGuardRefuse(rawEnv []string, dir string, stderr io.Writer) string {
	stale, reason := bootGuardVerdict(envFromSlice(rawEnv), FindGitRoot(dir), stderr)
	if !stale {
		return ""
	}
	return fmt.Sprintf(
		"spacedock: refusing — %s. Re-run `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json`, then retry this command; a fresh boot clears the guard.",
		reason)
}
