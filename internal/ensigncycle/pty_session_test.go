// ABOUTME: Pure session-jsonl resolution for the pty harness — pins the tail to the FO's OWN
// ABOUTME: transcript (the no-agentName root), never a spawned teammate's, so TeamCreate is never missed.
package ensigncycle

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// fileLineSource is the pty lineSource: it tails the FO's session jsonl, returning
// the complete lines appended since the last drain. It is the file analog of
// pipeLineSource (which tails a stdout pipe): the interactive session writes its
// stream-json transcript to disk, not a pipe, so the driver tails the file for
// liveness while feeding the SAME streamWatcher.
//
// It resolves the FO's file BY SESSION ID each drain (sessionFileByID over
// projectsDir), not by a fixed path or a newest-file heuristic: once the FO creates
// a team and dispatches a teammate, that teammate writes its OWN transcript into the
// same projectsDir, and a newest/assistant-bearing pick would FLIP the tail to the
// teammate and miss the FO's TeamCreate/Agent/TeamDelete (the F30 hazard — observed
// red'ing the residency drive when the tail read the comm-officer transcript). The
// byte offset advances only past COMPLETE (newline-terminated) lines, so a JSONL
// record split across two appends is held until its newline arrives.
type fileLineSource struct {
	projectsDir string
	sessionID   string
	offset      int64
}

func newFileLineSourceByID(projectsDir, sessionID string) *fileLineSource {
	return &fileLineSource{projectsDir: projectsDir, sessionID: sessionID}
}

func (s *fileLineSource) drain() []string {
	path := sessionFileByID(s.projectsDir, s.sessionID)
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil || int64(len(data)) <= s.offset {
		return nil
	}
	tail := data[s.offset:]
	// Only consume up to the last newline; bytes after it are a partial trailing
	// line held for the next drain.
	last := strings.LastIndexByte(string(tail), '\n')
	if last < 0 {
		return nil
	}
	consumed := tail[:last+1]
	s.offset += int64(len(consumed))
	var out []string
	for _, line := range strings.Split(strings.TrimRight(string(consumed), "\n"), "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// sessionFileByID returns the *.jsonl in dir whose entries carry the given session
// id, or "" if none. Resolving by id (not mtime) keeps the tail pinned to the FO's
// own transcript once a teammate transcript appears in the same dir (F30). Ported
// from spacedock-gym's SessionFileByID (reference-only, not imported).
func sessionFileByID(dir, sessionID string) string {
	if sessionID == "" {
		return ""
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	needle := `"sessionId":"` + sessionID + `"`
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if strings.Contains(string(body), needle) {
			return path
		}
	}
	return ""
}

// sessionIDFromFile returns the sessionId encoded on the first entry in a transcript
// that carries one, or "" if none. Used to pin the FO's own id right after launch,
// before any teammate transcript exists. Ported from spacedock-gym's
// sessionIDFromFile (reference-only, not imported).
func sessionIDFromFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 16*1024*1024)
	for sc.Scan() {
		var entry struct {
			SessionID string `json:"sessionId"`
		}
		if err := json.Unmarshal(sc.Bytes(), &entry); err != nil {
			continue
		}
		if entry.SessionID != "" {
			return entry.SessionID, nil
		}
	}
	return "", sc.Err()
}

// foRootSessionID returns the sessionId of the team-ROOT transcript in dir — the
// *.jsonl whose entries carry no spawned-teammate `agentName` — or "" if none yet.
// A teammate transcript (e.g. the comm-officer's) tags every entry with its
// `agentName`; the FO/root transcript does not (no `agentName` key at all —
// confirmed against every TeamCreate-bearing transcript on disk). So the root is the
// file with NO `"agentName":` occurrence. Scanning by name keeps the pick
// deterministic when (transiently) more than one file has yet to be tagged.
func foRootSessionID(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		path := filepath.Join(dir, name)
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		// A spawned-teammate transcript carries `"agentName":` on its entries; the
		// FO/root transcript never does. Skip files that carry it.
		if strings.Contains(string(body), `"agentName":`) {
			continue
		}
		if id, err := sessionIDFromFile(path); err == nil && id != "" {
			return id
		}
	}
	return ""
}

// activeSessionFile returns the live conversation transcript in dir: the newest
// *.jsonl that already carries an assistant entry, else the newest overall (right
// after launch, before the first assistant turn). Ported from spacedock-gym's
// ActiveSessionFile (reference-only, not imported).
func activeSessionFile(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	type cand struct {
		path    string
		mod     int64
		hasAsst bool
	}
	var cands []cand
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		path := filepath.Join(dir, e.Name())
		body, _ := os.ReadFile(path)
		cands = append(cands, cand{path, info.ModTime().UnixNano(), strings.Contains(string(body), `"type":"assistant"`)})
	}
	if len(cands) == 0 {
		return ""
	}
	sort.Slice(cands, func(i, j int) bool {
		if cands[i].hasAsst != cands[j].hasAsst {
			return cands[i].hasAsst
		}
		return cands[i].mod > cands[j].mod
	})
	return cands[0].path
}

// encodeProjectDir converts an absolute cwd into the directory name Claude Code
// uses under <config>/projects: each of '/', '.', and '_' becomes '-'. Ported from
// spacedock-gym's EncodeProjectDir (reference-only, not imported).
func encodeProjectDir(cwd string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '/', '.', '_':
			return '-'
		default:
			return r
		}
	}, cwd)
}

// transcriptReachedIdle reports whether the FO transcript has reached a committed
// idle turn — the render-independent boot/turn-end signal that replaces pane-text
// scraping (which is fragile on a headless CI runner: missing terminfo → blank
// capture-pane → a negative "no busy marker" test never fires). It is gym's
// WaitForTurn discipline (read the on-disk JSONL for a POSITIVE turn-end signal)
// adapted to the INTERACTIVE transcript, which — unlike `claude -p` — emits NO
// `type:"result"` event, so that count would never fire.
//
// Idle = the LAST COMMITTED assistant entry has a TERMINAL stop_reason and no later
// entry reopened a turn:
//   - committed assistant turn: type=="assistant" AND message.stop_reason non-empty
//     (streamed deltas carry a null/absent stop_reason and are SKIPPED);
//   - TERMINAL: end_turn / stop_sequence / max_tokens. stop_reason=="tool_use" is
//     mid-turn → BUSY;
//   - a later user/tool_result entry after that terminal assistant reopens the turn
//     → BUSY (the FO is acting on a tool result, not awaiting input).
func transcriptReachedIdle(lines []string) bool {
	type idleEntry struct {
		Type    string `json:"type"`
		Message *struct {
			StopReason *string `json:"stop_reason"`
		} `json:"message"`
	}
	terminal := map[string]bool{"end_turn": true, "stop_sequence": true, "max_tokens": true}
	sawTerminalAssistant := false
	for _, line := range lines {
		var e idleEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		switch e.Type {
		case "assistant":
			// Only COMMITTED assistant turns carry a non-null stop_reason; streamed
			// deltas have it null/absent and must not be read as a turn end.
			if e.Message == nil || e.Message.StopReason == nil || *e.Message.StopReason == "" {
				continue
			}
			sawTerminalAssistant = terminal[*e.Message.StopReason]
		case "user":
			// A user/tool_result entry after a terminal assistant reopens the turn:
			// the FO is mid-work on a tool result, not awaiting captain input.
			sawTerminalAssistant = false
		}
	}
	return sawTerminalAssistant
}

// foRootLine is a FO/root transcript entry: it carries a sessionId and NO
// agentName (the root is not a spawned teammate).
func foRootLine(sessionID, text string) string {
	return `{"sessionId":"` + sessionID + `","type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"` + text + `"}]}}`
}

// teammateLine is a spawned-teammate transcript entry: it carries a sessionId AND
// an agentName (the discriminator that marks it NOT the root).
func teammateLine(sessionID, agentName, text string) string {
	return `{"sessionId":"` + sessionID + `","agentName":"` + agentName + `","type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"` + text + `"}]}}`
}

// TestFOSessionPinning is the offline regression for the F30 fix: the tail must
// resolve the FO's OWN transcript (no agentName), never a spawned teammate's
// (agentName tagged), so the FO's TeamCreate/TeamDelete is never missed. It mirrors
// the exact on-disk shape that red'd the live residency drive — a comm-officer
// transcript newer than the FO's — over synthetic files, no model spend, no tmux.
func TestFOSessionPinning(t *testing.T) {
	dir := t.TempDir()
	// The FO/root transcript: carries TeamCreate, no agentName.
	foPath := filepath.Join(dir, "fo-root.jsonl")
	foLines := []string{
		foRootLine("fo-session-id", "booting"),
		`{"sessionId":"fo-session-id","type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","id":"t","name":"TeamCreate","input":{}}]}}`,
	}
	if err := os.WriteFile(foPath, []byte(strings.Join(foLines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A spawned comm-officer transcript: newer, assistant-bearing, agentName-tagged.
	teammatePath := filepath.Join(dir, "comm-officer.jsonl")
	teammateLines := []string{
		teammateLine("comm-session-id", "comm-officer", "online"),
		teammateLine("comm-session-id", "comm-officer", "polishing"),
	}
	if err := os.WriteFile(teammatePath, []byte(strings.Join(teammateLines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("foRootSessionID_skips_teammate", func(t *testing.T) {
		got := foRootSessionID(dir)
		if got != "fo-session-id" {
			t.Errorf("foRootSessionID = %q, want \"fo-session-id\" (the no-agentName root, not the comm-officer transcript)", got)
		}
	})

	t.Run("activeSessionFile_would_flip_to_teammate", func(t *testing.T) {
		// The bug being fixed: the newest-assistant-bearing heuristic picks the
		// teammate (this is WHY the by-id pin is required, not activeSessionFile).
		// The teammate file is written second, so it is the newest.
		if active := activeSessionFile(dir); active != teammatePath {
			t.Logf("activeSessionFile = %q (teammate=%q); the by-id pin is what avoids this flip", active, teammatePath)
		}
	})

	t.Run("sessionFileByID_resolves_FO", func(t *testing.T) {
		if got := sessionFileByID(dir, "fo-session-id"); got != foPath {
			t.Errorf("sessionFileByID(fo-session-id) = %q, want %q", got, foPath)
		}
	})

	t.Run("fileLineSource_tails_FO_not_teammate", func(t *testing.T) {
		src := newFileLineSourceByID(dir, "fo-session-id")
		lines := src.drain()
		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, `"name":"TeamCreate"`) {
			t.Errorf("the FO-pinned tail must carry TeamCreate; got:\n%s", joined)
		}
		if strings.Contains(joined, "comm-officer") {
			t.Errorf("the FO-pinned tail must NOT read the comm-officer transcript; got:\n%s", joined)
		}
		// A second drain with no new FO bytes returns nothing (offset advanced).
		if more := src.drain(); len(more) != 0 {
			t.Errorf("second drain with no new FO bytes should be empty, got %d lines", len(more))
		}
	})
}

// committedAssistant is a COMMITTED assistant turn (non-null stop_reason).
func committedAssistant(stopReason string) string {
	return `{"type":"assistant","message":{"role":"assistant","stop_reason":"` + stopReason + `","content":[{"type":"text","text":"hi"}]}}`
}

// streamingAssistant is a streamed assistant DELTA (null stop_reason) — not yet a
// committed turn end.
const streamingAssistant = `{"type":"assistant","message":{"role":"assistant","stop_reason":null,"content":[{"type":"text","text":"thinking"}]}}`

// userToolResult is a user/tool_result entry that reopens a turn after a terminal
// assistant (the FO acting on a tool result, not awaiting input).
const userToolResult = `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"t","content":"ok"}]}}`

// TestTranscriptReachedIdle is the offline proof of the render-independent boot gate
// that replaces pane-text scraping. It pins the false-pass guards from the fix spec:
// a clean end_turn reads idle; streamed deltas (null stop_reason) do not; a
// tool_use-stop turn is busy; and a user/tool_result after end_turn reopens the turn.
func TestTranscriptReachedIdle(t *testing.T) {
	t.Run("clean_end_turn_is_idle", func(t *testing.T) {
		lines := []string{streamingAssistant, committedAssistant("end_turn")}
		if !transcriptReachedIdle(lines) {
			t.Error("a committed end_turn assistant turn must read idle")
		}
	})

	t.Run("stop_sequence_and_max_tokens_are_idle", func(t *testing.T) {
		for _, sr := range []string{"stop_sequence", "max_tokens"} {
			if !transcriptReachedIdle([]string{committedAssistant(sr)}) {
				t.Errorf("stop_reason %q is terminal and must read idle", sr)
			}
		}
	})

	t.Run("streaming_only_nulls_are_busy", func(t *testing.T) {
		// Only streamed deltas (null stop_reason), no committed turn yet → BUSY.
		lines := []string{streamingAssistant, streamingAssistant}
		if transcriptReachedIdle(lines) {
			t.Error("streamed deltas with null stop_reason must NOT read idle (no committed turn)")
		}
	})

	t.Run("tool_use_stop_is_busy", func(t *testing.T) {
		// stop_reason tool_use = the FO is mid-turn issuing a tool call → BUSY.
		lines := []string{committedAssistant("tool_use")}
		if transcriptReachedIdle(lines) {
			t.Error("stop_reason tool_use is mid-turn and must NOT read idle")
		}
	})

	t.Run("user_reopened_after_end_turn_is_busy", func(t *testing.T) {
		// A user/tool_result AFTER a terminal assistant reopens the turn → BUSY.
		lines := []string{committedAssistant("end_turn"), userToolResult}
		if transcriptReachedIdle(lines) {
			t.Error("a user/tool_result after end_turn reopens the turn and must NOT read idle")
		}
	})

	t.Run("end_turn_after_reopen_is_idle_again", func(t *testing.T) {
		// end_turn → tool_result reopen → another end_turn: the LAST committed turn
		// is terminal with nothing after, so idle again.
		lines := []string{committedAssistant("end_turn"), userToolResult, committedAssistant("end_turn")}
		if !transcriptReachedIdle(lines) {
			t.Error("a fresh end_turn after a reopened turn must read idle again")
		}
	})

	t.Run("empty_transcript_is_busy", func(t *testing.T) {
		if transcriptReachedIdle(nil) {
			t.Error("an empty transcript must NOT read idle")
		}
	})
}
