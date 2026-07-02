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

type heartbeat struct {
	SessionID string    `json:"session_id"`
	TS        time.Time `json:"ts"`
	State     string    `json:"state"`
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

// Wake resumes live Codex FO sessions that are addressed by inbox records after
// the successful-wake cursor. It advances the cursor only after at least one
// resume succeeds; no-session records remain eligible for a future wake.
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

	cursor := readCursor(absRoot)
	records, lastLine, err := readPendingInbox(absRoot, cursor)
	if err != nil {
		return Result{Status: "failed", Error: err.Error()}
	}
	if len(records) == 0 {
		if lastLine > cursor {
			_ = writeCursor(absRoot, lastLine)
		}
		return Result{Status: "noop", Message: "no pending inbox records"}
	}

	sessions := map[string]*sessionWake{}
	targetsMissingSession := map[string]bool{}
	for _, rec := range records {
		targets := targetsFor(absRoot, rec, opts.Members)
		if len(targets) == 0 {
			continue
		}
		for _, target := range targets {
			hb, ok := loadHeartbeat(absRoot, target, now())
			if !ok || hb.SessionID == "" {
				targetsMissingSession[target] = true
				continue
			}
			w := sessions[hb.SessionID]
			if w == nil {
				w = &sessionWake{SessionID: hb.SessionID, TargetSet: map[string]bool{}}
				sessions[hb.SessionID] = w
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
			Message:   "no fresh FO heartbeat with a Codex session id",
		})
		return Result{Status: "skipped-no-session", Lines: recordLines(records), Targets: targets, Message: "no fresh FO heartbeat with a Codex session id"}
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

	if successes > 0 {
		_ = writeCursor(absRoot, lastLine)
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
	return nil
}

func wakePrompt(root string, w *sessionWake) string {
	return fmt.Sprintf(`Bridge queued captain intent for this Spacedock first-officer session.

Repo root: %s
Inbox: %s
Pending physical inbox lines for this session: %s
Addressed workflow slugs: %s

Run the Bridge inbox idle drain now for only those workflow slugs. Honor target_set routing, append per-slug FO replies or acknowledgements for valid captain intents, advance the matching _bridge/.inbox-cursor.<slug> files, then continue the normal first-officer event loop.
`, root, filepath.Join(root, "_bridge", "inbox.jsonl"), joinInts(w.Lines), strings.Join(w.Targets(), ","))
}

func acquireLock(root string) (func(), bool) {
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return func() {}, false
	}
	path := filepath.Join(dir, ".wake-lock.codex")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return func() {}, false
	}
	_, _ = fmt.Fprintf(f, "%d\n", os.Getpid())
	_ = f.Close()
	return func() { _ = os.Remove(path) }, true
}

func readCursor(root string) int {
	data, err := os.ReadFile(filepath.Join(root, "_bridge", ".wake-cursor.codex"))
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return n
}

func writeCursor(root string, line int) error {
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ".wake-cursor.codex"), []byte(strconv.Itoa(line)+"\n"), 0o600)
}

func readPendingInbox(root string, cursor int) ([]inboxRecord, int, error) {
	f, err := os.Open(filepath.Join(root, "_bridge", "inbox.jsonl"))
	if os.IsNotExist(err) {
		return nil, cursor, nil
	}
	if err != nil {
		return nil, cursor, err
	}
	defer f.Close()

	var out []inboxRecord
	lineNo := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNo++
		if lineNo <= cursor || strings.TrimSpace(scanner.Text()) == "" {
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
		return nil, lineNo, err
	}
	return out, lineNo, nil
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
	data, err := os.ReadFile(filepath.Join(root, "_bridge", "fo."+slug+".json"))
	if err != nil {
		return hb, false
	}
	if err := json.Unmarshal(data, &hb); err != nil || hb.TS.IsZero() {
		return hb, false
	}
	age := now.Sub(hb.TS)
	return hb, age >= 0 && age <= liveWindow
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
