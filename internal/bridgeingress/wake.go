// ABOUTME: Bridge ingress wake-up for Codex first-officer sessions.
// ABOUTME: Reads the durable Bridge inbox and nudges live Codex sessions to drain it.
package bridgeingress

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const liveWindow = 30 * time.Minute

// maxScanLine raises the per-line scan limit above bufio's 64KB default. Bridge
// control records are small, but a large captain intent must not make the scan
// hard-fail (inbox) or silently stop short (replies/events).
const maxScanLine = 1 << 20

func lineScanner(f *os.File) *bufio.Scanner {
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 64*1024), maxScanLine)
	return s
}

var safeSlugPattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// ResumeFunc resumes a host session with the supplied prompt.
type ResumeFunc func(ctx context.Context, sessionID, prompt string) error

// Options controls one Bridge inbox wake pass.
type Options struct {
	Host     string
	Root     string
	Members  []string
	CodexBin string
	Now      func() time.Time
	Resume   ResumeFunc
}

// Result is the JSON shape printed by the hidden CLI for Bridge to display.
type Result struct {
	Status   string   `json:"status"`
	Lines    []int    `json:"lines,omitempty"`
	Sessions int      `json:"sessions,omitempty"`
	Targets  []string `json:"targets,omitempty"`
	Message  string   `json:"message,omitempty"`
	Error    string   `json:"error,omitempty"`
}

type inboxRecord struct {
	ID        string    `json:"id"`
	TS        time.Time `json:"ts"`
	Kind      string    `json:"kind"`
	Target    string    `json:"target"`
	TargetSet []string  `json:"target_set"`
	Line      int       `json:"-"`
}

type replyRecord struct {
	Schema        int       `json:"schema"`
	Kind          string    `json:"kind"`
	Target        string    `json:"target"`
	InReplyToID   string    `json:"in_reply_to_id"`
	InReplyToLine int       `json:"in_reply_to_line"`
	InReplyToTS   time.Time `json:"in_reply_to_ts"`
	IntentKind    string    `json:"intent_kind"`
	Status        string    `json:"status"`
}

type heartbeat struct {
	SessionID string    `json:"session_id"`
	TS        time.Time `json:"ts"`
	State     string    `json:"state"`
}

type sessionMarker struct {
	SessionID string `json:"session_id"`
	Workflow  string `json:"workflow"`
}

type eventRecord struct {
	Host      string    `json:"host"`
	TS        time.Time `json:"ts"`
	SessionID string    `json:"session_id"`
}

type wakeEvent struct {
	Timestamp string   `json:"timestamp"`
	TS        string   `json:"ts"`
	Host      string   `json:"host"`
	Event     string   `json:"event"`
	Status    string   `json:"status"`
	Line      int      `json:"line,omitempty"`
	Lines     []int    `json:"lines,omitempty"`
	IntentID  string   `json:"intent_id,omitempty"`
	Targets   []string `json:"targets,omitempty"`
	SessionID string   `json:"session_id,omitempty"`
	Message   string   `json:"message,omitempty"`
	Error     string   `json:"error,omitempty"`
}

// Wake resumes Codex FO sessions for inbox records that are not yet delivered by
// the addressed workflow cursors or an FO reply/ack. Starting a resume process is
// only a wake attempt; delivery is confirmed later by cursor advancement or ack.
func Wake(ctx context.Context, opts Options) Result {
	host := normalizeHost(opts.Host)
	if host == "" {
		host = "codex"
	}
	if host != "codex" {
		return Result{Status: "failed", Error: "bridge ingress wake currently supports host codex only"}
	}
	root := opts.Root
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Result{Status: "failed", Error: err.Error()}
	}
	now := time.Now().UTC
	if opts.Now != nil {
		now = func() time.Time { return opts.Now().UTC() }
	}

	unlock, ok := acquireLock(absRoot)
	if !ok {
		return Result{Status: "locked", Message: "another bridge ingress wake is running"}
	}
	defer unlock()

	allRecords, err := readInbox(absRoot)
	if err != nil {
		return Result{Status: "failed", Error: err.Error()}
	}
	replies := loadReplies(absRoot)
	var records []inboxRecord
	for _, rec := range allRecords {
		if len(pendingTargetsFor(absRoot, rec, opts.Members, replies)) > 0 {
			records = append(records, rec)
		}
	}
	if len(records) == 0 {
		return Result{Status: "noop", Message: "no pending inbox records"}
	}

	sessions := map[string]*sessionWake{}
	targetsMissingSession := map[string]bool{}
	for _, rec := range records {
		targets := pendingTargetsFor(absRoot, rec, opts.Members, replies)
		if len(targets) == 0 {
			continue
		}
		for _, target := range targets {
			sessionID, ok := resumableSessionID(absRoot, target)
			if !ok {
				targetsMissingSession[target] = true
				continue
			}
			w := sessions[sessionID]
			if w == nil {
				w = &sessionWake{SessionID: sessionID, TargetSet: map[string]bool{}}
				sessions[sessionID] = w
			}
			w.TargetSet[target] = true
			w.Lines = appendUniqueInt(w.Lines, rec.Line)
			if rec.ID != "" {
				w.IntentIDs = appendUniqueString(w.IntentIDs, rec.ID)
			}
		}
	}

	if len(sessions) == 0 {
		targets := keys(targetsMissingSession)
		appendWakeEvent(absRoot, wakeEvent{
			Timestamp: now().Format(time.RFC3339),
			TS:        now().Format(time.RFC3339),
			Host:      host,
			Event:     "wake",
			Status:    "skipped-no-session",
			Lines:     recordLines(records),
			Targets:   targets,
			Message:   "no resumable Codex session id",
		})
		return Result{Status: "skipped-no-session", Lines: recordLines(records), Targets: targets, Message: "no resumable Codex session id"}
	}

	resume := opts.Resume
	if resume == nil {
		resume = func(ctx context.Context, sessionID, prompt string) error {
			return execCodexResume(ctx, opts.CodexBin, sessionID, prompt)
		}
	}

	var successes int
	var firstErr error
	allTargets := map[string]bool{}
	for _, w := range sortedSessionWakes(sessions) {
		for target := range w.TargetSet {
			allTargets[target] = true
		}
		prompt := wakePrompt(absRoot, w)
		err := resume(ctx, w.SessionID, prompt)
		status := "woke"
		errText := ""
		if err != nil {
			status = "failed"
			errText = err.Error()
			if firstErr == nil {
				firstErr = err
			}
		} else {
			successes++
		}
		appendWakeEvent(absRoot, wakeEvent{
			Timestamp: now().Format(time.RFC3339),
			TS:        now().Format(time.RFC3339),
			Host:      host,
			Event:     "wake",
			Status:    status,
			Lines:     append([]int(nil), w.Lines...),
			Targets:   w.Targets(),
			SessionID: w.SessionID,
			Error:     errText,
		})
	}

	result := Result{
		Status:   "woke",
		Lines:    recordLines(records),
		Sessions: successes,
		Targets:  keys(allTargets),
		Message:  "resumed Codex FO session",
	}
	if firstErr != nil {
		result.Status = "partial"
		result.Error = firstErr.Error()
		if successes == 0 {
			result.Status = "failed"
			result.Message = ""
		}
	}
	return result
}

type sessionWake struct {
	SessionID string
	Lines     []int
	TargetSet map[string]bool
	IntentIDs []string
}

func (w *sessionWake) Targets() []string { return keys(w.TargetSet) }

func execCodexResume(ctx context.Context, bin, sessionID, prompt string) error {
	if strings.TrimSpace(bin) == "" {
		bin = "codex"
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	promptFile, err := os.CreateTemp("", "spacedock-bridge-wake-*.txt")
	if err != nil {
		return fmt.Errorf("codex exec resume prompt: %w", err)
	}
	promptPath := promptFile.Name()
	defer func() {
		_ = promptFile.Close()
		_ = os.Remove(promptPath)
	}()
	if _, err := promptFile.WriteString(prompt); err != nil {
		return fmt.Errorf("codex exec resume prompt: %w", err)
	}
	if _, err := promptFile.Seek(0, 0); err != nil {
		return fmt.Errorf("codex exec resume prompt: %w", err)
	}

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("codex exec resume output: %w", err)
	}
	defer devNull.Close()

	cmd := exec.Command(bin, "exec", "resume", sessionID, "-")
	cmd.Stdin = promptFile
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("codex exec resume: %w", err)
	}
	go func() {
		_ = cmd.Wait()
	}()
	return nil
}

// wakePrompt instructs the resumed Codex FO to drain the inbox using the direct
// _bridge/ file protocol (docs/seam-contract.md §3) — the FO writes files, not a
// spacedock verb. It deliberately names no `spacedock bridge inbox ...` command:
// draining, acking, and cursor advancement are file writes the FO performs
// itself (findings B.2). The three substrings the wake tests pin — the pending
// line list, the addressed slug list, and the inbox path — are preserved.
func wakePrompt(root string, w *sessionWake) string {
	return fmt.Sprintf(`Bridge queued captain intent for this Spacedock first-officer session.

Repo root: %s
Inbox: %s
Pending physical inbox lines for this session: %s
Addressed workflow slugs: %s

Drain the Bridge inbox now for ONLY those workflow slugs, writing _bridge/ files directly per the seam contract (no spacedock inbox command):
1. For each slug S, read the current cursor from _bridge/.inbox-cursor.<S> (treat a missing/unparseable file as 0).
2. Read _bridge/inbox.jsonl line by line, counting physical newline-terminated lines 1-based (a malformed line still consumes its number; skip a trailing fragment with no newline). For each line numbered N greater than the cursor whose routing addresses S (target "all"/absent, target == S, or S in target_set), act on the intent.
3. Before acting on an intent id X, scan _bridge/fo-replies.jsonl for a terminal ack whose in_reply_to_id is X and skip the intent if one already exists (dedup).
4. Append a terminal ack line to _bridge/fo-replies.jsonl for each intent you handle: schema:1, in_reply_to_id set to the intent id, intent_kind set to the intent kind, the matching reply kind, a valid terminal status, and your slug as target.
5. Write the inbox's new physical line count (wc -l) to _bridge/.inbox-cursor.<S> as a whole-file replace — monotonic, never lower it.
6. Refresh the heartbeat _bridge/fo.<S>.json (ts now, state "working" while acting).
Then continue the normal first-officer event loop.
`, root, filepath.Join(root, "_bridge", "inbox.jsonl"), joinInts(w.Lines), strings.Join(w.Targets(), ","))
}

// staleLockTTL bounds how long a wake lock may persist before another wake may
// reclaim it. A wake pass only starts resume processes and returns — it never
// blocks on Codex — so it holds the lock for well under a second. A lock older
// than this TTL was left by a crashed or killed wake that never ran its deferred
// unlock; reclaiming it keeps a single failure from wedging durable delivery
// permanently.
const staleLockTTL = 5 * time.Minute

func acquireLock(root string) (func(), bool) {
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return func() {}, false
	}
	path := filepath.Join(dir, ".wake-lock.codex")
	if unlock, ok := takeLock(path); ok {
		return unlock, true
	}
	// The lock exists. Reclaim it only if it is stale; otherwise another wake is
	// genuinely running.
	if info, err := os.Stat(path); err != nil || time.Since(info.ModTime()) <= staleLockTTL {
		return func() {}, false
	}
	_ = os.Remove(path)
	return takeLock(path)
}

func takeLock(path string) (func(), bool) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return func() {}, false
	}
	_, _ = fmt.Fprintf(f, "%d\n", os.Getpid())
	_ = f.Close()
	return func() { _ = os.Remove(path) }, true
}

func readInbox(root string) ([]inboxRecord, error) {
	f, err := os.Open(filepath.Join(root, "_bridge", "inbox.jsonl"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []inboxRecord
	lineNo := 0
	scanner := lineScanner(f)
	for scanner.Scan() {
		lineNo++
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}
		var rec inboxRecord
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			continue
		}
		rec.Line = lineNo
		out = append(out, rec)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func pendingTargetsFor(root string, rec inboxRecord, members []string, replies map[string]bool) []string {
	var pending []string
	for _, target := range targetsFor(root, rec, members) {
		if inboxCursor(root, target) >= rec.Line {
			continue
		}
		if replies[replyKey(rec, target)] {
			continue
		}
		pending = append(pending, target)
	}
	return pending
}

func inboxCursor(root, slug string) int {
	if !safeSlugPattern.MatchString(slug) || slug == "." || slug == ".." {
		return 0
	}
	data, err := os.ReadFile(filepath.Join(root, "_bridge", ".inbox-cursor."+slug))
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	if n < 0 {
		return 0
	}
	return n
}

func loadReplies(root string) map[string]bool {
	out := map[string]bool{}
	f, err := os.Open(filepath.Join(root, "_bridge", "fo-replies.jsonl"))
	if err != nil {
		return out
	}
	defer f.Close()
	scanner := lineScanner(f)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}
		var rec replyRecord
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			continue
		}
		if rec.Schema != 1 || rec.Target == "" || rec.InReplyToLine <= 0 || rec.IntentKind == "" {
			continue
		}
		if !safeSlugPattern.MatchString(rec.Target) || rec.Target == "." || rec.Target == ".." {
			continue
		}
		out[replyKey(inboxRecord{ID: rec.InReplyToID, TS: rec.InReplyToTS, Kind: rec.IntentKind, Line: rec.InReplyToLine}, rec.Target)] = true
	}
	return out
}

func replyKey(rec inboxRecord, target string) string {
	id := rec.ID
	if id == "" {
		id = rec.TS.Format(time.RFC3339Nano)
	}
	return strconv.Itoa(rec.Line) + "\x00" + id + "\x00" + rec.Kind + "\x00" + target
}

func targetsFor(root string, rec inboxRecord, members []string) []string {
	if len(rec.TargetSet) > 0 {
		return cleanSlugs(rec.TargetSet)
	}
	target := strings.TrimSpace(rec.Target)
	if target == "" || target == "all" {
		if len(members) > 0 {
			return cleanSlugs(members)
		}
		return discoverHeartbeatSlugs(root)
	}
	return cleanSlugs([]string{target})
}

func discoverHeartbeatSlugs(root string) []string {
	matches, _ := filepath.Glob(filepath.Join(root, "_bridge", "fo.*.json"))
	var out []string
	for _, path := range matches {
		name := filepath.Base(path)
		slug := strings.TrimSuffix(strings.TrimPrefix(name, "fo."), ".json")
		out = append(out, slug)
	}
	return cleanSlugs(out)
}

func cleanSlugs(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, slug := range in {
		slug = strings.TrimSpace(slug)
		if slug == "" || slug == "." || slug == ".." || !safeSlugPattern.MatchString(slug) || seen[slug] {
			continue
		}
		seen[slug] = true
		out = append(out, slug)
	}
	sort.Strings(out)
	return out
}

func loadHeartbeat(root, slug string, now time.Time) (heartbeat, bool) {
	var hb heartbeat
	if !safeSlugPattern.MatchString(slug) || slug == "." || slug == ".." {
		return hb, false
	}
	hb, ok := loadHeartbeatAnyAge(root, slug)
	if !ok || hb.TS.IsZero() {
		return hb, false
	}
	age := now.Sub(hb.TS)
	return hb, age >= 0 && age <= liveWindow
}

func loadHeartbeatAnyAge(root, slug string) (heartbeat, bool) {
	var hb heartbeat
	if !safeSlugPattern.MatchString(slug) || slug == "." || slug == ".." {
		return hb, false
	}
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "fo."+slug+".json"))
	if err != nil {
		return hb, false
	}
	if err := json.Unmarshal(data, &hb); err != nil {
		return hb, false
	}
	return hb, true
}

func resumableSessionID(root, slug string) (string, bool) {
	if hb, ok := loadHeartbeatAnyAge(root, slug); ok {
		if id := strings.TrimSpace(hb.SessionID); safeSessionID(id) {
			return id, true
		}
	}
	if sessionID, ok := sessionIDFromMarkers(root, slug); ok && safeSessionID(sessionID) {
		return sessionID, true
	}
	if sessionID, ok := latestCodexEventSession(root); ok && safeSessionID(sessionID) {
		return sessionID, true
	}
	return "", false
}

// safeSessionID guards a session id read from _bridge/ state before it becomes a
// codex argv positional, so a poisoned marker/heartbeat/event line cannot inject
// a leading-dash token that codex would parse as a flag.
func safeSessionID(s string) bool {
	return s != "" && s != "." && s != ".." && safeSlugPattern.MatchString(s)
}

func sessionIDFromMarkers(root, slug string) (string, bool) {
	matches, _ := filepath.Glob(filepath.Join(root, "_bridge", "sessions", "*.json"))
	sort.Strings(matches)
	var bestPath string
	var bestMod time.Time
	var bestSession string
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var rec sessionMarker
		if err := json.Unmarshal(data, &rec); err != nil {
			continue
		}
		if rec.Workflow != slug || strings.TrimSpace(rec.SessionID) == "" {
			continue
		}
		info, err := os.Stat(path)
		mod := time.Time{}
		if err == nil {
			mod = info.ModTime()
		}
		if bestSession == "" || mod.After(bestMod) || (mod.Equal(bestMod) && path > bestPath) {
			bestPath = path
			bestMod = mod
			bestSession = strings.TrimSpace(rec.SessionID)
		}
	}
	return bestSession, bestSession != ""
}

func latestCodexEventSession(root string) (string, bool) {
	f, err := os.Open(filepath.Join(root, "_bridge", "events.jsonl"))
	if err != nil {
		return "", false
	}
	defer f.Close()
	var best eventRecord
	scanner := lineScanner(f)
	for scanner.Scan() {
		var rec eventRecord
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			continue
		}
		if normalizeHost(rec.Host) != "codex" || strings.TrimSpace(rec.SessionID) == "" {
			continue
		}
		if best.SessionID == "" || rec.TS.After(best.TS) {
			best = rec
		}
	}
	if best.SessionID == "" {
		return "", false
	}
	return strings.TrimSpace(best.SessionID), true
}

func appendWakeEvent(root string, event wakeEvent) {
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(filepath.Join(dir, "wake-events.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	_, _ = f.Write(append(data, '\n'))
}

func sortedSessionWakes(in map[string]*sessionWake) []*sessionWake {
	var sessions []string
	for session := range in {
		sessions = append(sessions, session)
	}
	sort.Strings(sessions)
	out := make([]*sessionWake, 0, len(sessions))
	for _, session := range sessions {
		w := in[session]
		sort.Ints(w.Lines)
		w.IntentIDs = cleanStrings(w.IntentIDs)
		out = append(out, w)
	}
	return out
}

func recordLines(records []inboxRecord) []int {
	lines := make([]int, 0, len(records))
	for _, rec := range records {
		lines = appendUniqueInt(lines, rec.Line)
	}
	sort.Ints(lines)
	return lines
}

func appendUniqueInt(in []int, v int) []int {
	for _, existing := range in {
		if existing == v {
			return in
		}
	}
	return append(in, v)
}

func appendUniqueString(in []string, v string) []string {
	for _, existing := range in {
		if existing == v {
			return in
		}
	}
	return append(in, v)
}

func cleanStrings(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func joinInts(in []int) string {
	parts := make([]string, 0, len(in))
	for _, n := range in {
		parts = append(parts, strconv.Itoa(n))
	}
	return strings.Join(parts, ",")
}

func normalizeHost(host string) string {
	return strings.ToLower(strings.TrimSpace(host))
}
