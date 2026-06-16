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
