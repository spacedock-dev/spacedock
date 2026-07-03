// ABOUTME: Deterministic Bridge inbox drain/ack/commit/check for FO sessions.
// ABOUTME: Moves cursor math, routing, heartbeat, and ack serialization out of FO prose into the binary.
package bridgeingress

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// The drain family packages the mechanism the bridge-inbox mod used to make the
// FO hand-execute as raw shell each tick: cursor read/advance, per-line routing,
// heartbeat stamping, and reply/ack serialization. The FO now calls these verbs
// and keeps only the judgment (interpreting a tell, resolving a gate). This is
// host-neutral; --host only selects which session-id token the heartbeat carries.

// DrainOptions controls one inbox drain (peek) pass.
type DrainOptions struct {
	Host      string
	Root      string
	Slug      string
	SessionID string
	Members   []string
	Now       func() time.Time
}

// DrainRecord is one addressed, not-yet-acked inbox record handed back to the FO.
type DrainRecord struct {
	Line       int      `json:"line"`
	ID         string   `json:"id"`
	TS         string   `json:"ts"`
	Kind       string   `json:"kind"`
	Text       string   `json:"text,omitempty"`
	Granted    *bool    `json:"granted,omitempty"`
	Target     string   `json:"target,omitempty"`
	TargetSet  []string `json:"target_set,omitempty"`
	Entity     string   `json:"entity,omitempty"`
	Field      string   `json:"field,omitempty"`
	Value      string   `json:"value,omitempty"`
	Verdict    string   `json:"verdict,omitempty"`
	Directives []string `json:"directives,omitempty"`
	RequestID  string   `json:"request_id,omitempty"`
}

// DrainResult is the JSON the FO reads after a drain. It never advances the
// cursor: the FO acts on records, acks them, then commits the high-water mark.
type DrainResult struct {
	Status    string        `json:"status"`
	Slug      string        `json:"slug"`
	Host      string        `json:"host,omitempty"`
	Cursor    int           `json:"cursor"`
	HighWater int           `json:"high_water"`
	Count     int           `json:"count"`
	Heartbeat bool          `json:"heartbeat"`
	Records   []DrainRecord `json:"records"`
	Error     string        `json:"error,omitempty"`
}

type fullInboxRecord struct {
	ID         string    `json:"id"`
	TS         time.Time `json:"ts"`
	RawTS      string    `json:"-"`
	Kind       string    `json:"kind"`
	Text       string    `json:"text"`
	Granted    *bool     `json:"granted"`
	Target     string    `json:"target"`
	TargetSet  []string  `json:"target_set"`
	Entity     string    `json:"entity"`
	Field      string    `json:"field"`
	Value      string    `json:"value"`
	Verdict    string    `json:"verdict"`
	Directives []string  `json:"directives"`
	RequestID  string    `json:"request_id"`
	Line       int       `json:"-"`
}

func (r fullInboxRecord) routing() inboxRecord {
	return inboxRecord{ID: r.ID, TS: r.TS, Kind: r.Kind, Target: r.Target, TargetSet: r.TargetSet, Line: r.Line}
}

// Drain reads newly-queued captain intent addressed to the given workflow slug,
// stamps the liveness heartbeat, and returns the records without advancing the
// cursor. It is idempotent: records already acked (per fo-replies.jsonl) are
// filtered out even if the cursor has not yet advanced past them.
func Drain(opts DrainOptions) DrainResult {
	slug := strings.TrimSpace(opts.Slug)
	if !validSlug(slug) {
		return DrainResult{Status: "failed", Slug: slug, Error: "invalid or missing --slug"}
	}
	host := normalizeHost(opts.Host)
	root := absRootOr(opts.Root)
	now := nowFunc(opts.Now)

	// Heartbeat is observe-only liveness and must be written on every tick, even
	// when no Bridge is attached or the inbox is empty.
	hbWritten := stampHeartbeat(root, slug, resolveSessionID(host, opts.SessionID), host, now())

	res := DrainResult{Status: "ok", Slug: slug, Host: host, Heartbeat: hbWritten, Records: []DrainRecord{}}

	inboxPath := filepath.Join(root, "_bridge", "inbox.jsonl")
	if _, err := os.Stat(inboxPath); err != nil {
		res.Status = "no-inbox"
		return res
	}

	adoptCursor(root, slug)
	cursor := inboxCursor(root, slug)
	records, highWater, err := readInboxFull(inboxPath)
	if err != nil {
		return DrainResult{Status: "failed", Slug: slug, Host: host, Heartbeat: hbWritten, Error: err.Error(), Records: []DrainRecord{}}
	}
	res.Cursor = cursor
	res.HighWater = highWater

	replies := loadReplies(root)
	members := opts.Members
	if len(members) == 0 {
		members = []string{slug}
	}
	for _, rec := range records {
		if rec.Line <= cursor {
			continue
		}
		if !addressedTo(root, rec.routing(), slug, members) {
			continue
		}
		if replies[replyKey(rec.routing(), slug)] {
			continue
		}
		res.Records = append(res.Records, toDrainRecord(rec))
	}
	res.Count = len(res.Records)
	return res
}

// addressedTo reports whether an inbox record routes to this workflow slug,
// honoring an authoritative frozen target_set over the legacy target field.
func addressedTo(root string, rec inboxRecord, slug string, members []string) bool {
	for _, t := range targetsFor(root, rec, members) {
		if t == slug {
			return true
		}
	}
	return false
}

// AckOptions controls one reply/ack append.
type AckOptions struct {
	Host       string
	Root       string
	Slug       string
	Line       int
	ID         string
	TS         string
	IntentKind string
	Status     string
	Text       string
	Granted    *bool
	Entity     string
	Field      string
	Value      string
	Verdict    string
	RequestID  string
	SessionID  string
	Now        func() time.Time
}

// AckResult is the JSON printed after an ack append.
type AckResult struct {
	Appended bool   `json:"appended"`
	Kind     string `json:"kind,omitempty"`
	Target   string `json:"target,omitempty"`
	Line     int    `json:"line,omitempty"`
	Error    string `json:"error,omitempty"`
}

type replyOut struct {
	Schema        int    `json:"schema"`
	TS            string `json:"ts"`
	Kind          string `json:"kind"`
	Target        string `json:"target"`
	InReplyToID   string `json:"in_reply_to_id,omitempty"`
	InReplyToLine int    `json:"in_reply_to_line"`
	InReplyToTS   string `json:"in_reply_to_ts,omitempty"`
	IntentKind    string `json:"intent_kind"`
	Status        string `json:"status"`
	Text          string `json:"text,omitempty"`
	Granted       *bool  `json:"granted,omitempty"`
	Entity        string `json:"entity,omitempty"`
	Field         string `json:"field,omitempty"`
	Value         string `json:"value,omitempty"`
	Verdict       string `json:"verdict,omitempty"`
	RequestID     string `json:"request_id,omitempty"`
	SessionID     string `json:"session_id,omitempty"`
	Host          string `json:"host,omitempty"`
}

// ackKindFor maps an inbox intent kind to its reply/ack record kind.
func ackKindFor(intentKind string) (string, bool) {
	switch intentKind {
	case "tell":
		return "reply", true
	case "conn":
		return "conn-ack", true
	case "decision":
		return "decision-ack", true
	case "permission-decision":
		return "permission-ack", true
	default:
		return "", false
	}
}

// Ack appends one newline-terminated, compact reply record to fo-replies.jsonl.
// Serialization is owned by the binary so the FO never hand-encodes JSONL — the
// non-compact / multi-line JSONL corruption class is eliminated.
func Ack(opts AckOptions) AckResult {
	slug := strings.TrimSpace(opts.Slug)
	if !validSlug(slug) {
		return AckResult{Error: "invalid or missing --slug"}
	}
	if opts.Line <= 0 {
		return AckResult{Error: "missing or non-positive --line"}
	}
	kind, ok := ackKindFor(opts.IntentKind)
	if !ok {
		return AckResult{Error: "unknown --kind (want tell|conn|decision|permission-decision)"}
	}
	if strings.TrimSpace(opts.Status) == "" {
		return AckResult{Error: "missing --status"}
	}
	host := normalizeHost(opts.Host)
	root := absRootOr(opts.Root)
	now := nowFunc(opts.Now)

	rec := replyOut{
		Schema:        1,
		TS:            now().UTC().Format(time.RFC3339),
		Kind:          kind,
		Target:        slug,
		InReplyToID:   opts.ID,
		InReplyToLine: opts.Line,
		InReplyToTS:   strings.TrimSpace(opts.TS),
		IntentKind:    opts.IntentKind,
		Status:        opts.Status,
		Text:          singleLine(opts.Text),
		Granted:       opts.Granted,
		Entity:        opts.Entity,
		Field:         opts.Field,
		Value:         opts.Value,
		Verdict:       opts.Verdict,
		RequestID:     opts.RequestID,
		SessionID:     resolveSessionID(host, opts.SessionID),
		Host:          host,
	}
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return AckResult{Error: err.Error()}
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return AckResult{Error: err.Error()}
	}
	f, err := os.OpenFile(filepath.Join(dir, "fo-replies.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return AckResult{Error: err.Error()}
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return AckResult{Error: err.Error()}
	}
	return AckResult{Appended: true, Kind: kind, Target: slug, Line: opts.Line}
}

// CommitOptions controls one cursor advance.
type CommitOptions struct {
	Root   string
	Slug   string
	Cursor int
}

// CommitResult is the JSON printed after a cursor advance.
type CommitResult struct {
	Cursor int    `json:"cursor"`
	Error  string `json:"error,omitempty"`
}

// Commit advances this workflow's inbox cursor to the supplied high-water mark.
// It is monotonic: a value below the current cursor is ignored so a stale commit
// can never re-deliver already-processed intent.
func Commit(opts CommitOptions) CommitResult {
	slug := strings.TrimSpace(opts.Slug)
	if !validSlug(slug) {
		return CommitResult{Error: "invalid or missing --slug"}
	}
	if opts.Cursor < 0 {
		return CommitResult{Error: "negative --cursor"}
	}
	root := absRootOr(opts.Root)
	adoptCursor(root, slug)
	current := inboxCursor(root, slug)
	target := opts.Cursor
	if target < current {
		target = current
	}
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return CommitResult{Error: err.Error()}
	}
	if err := os.WriteFile(filepath.Join(dir, ".inbox-cursor."+slug), []byte(strconv.Itoa(target)+"\n"), 0o644); err != nil {
		return CommitResult{Error: err.Error()}
	}
	return CommitResult{Cursor: target}
}

// CheckOptions controls a Stop-hook drain check.
type CheckOptions struct {
	Host      string
	Root      string
	Slug      string
	SessionID string
	Members   []string
	// StopHookActive mirrors the Claude Stop hook payload field. When true the
	// check never blocks again, so a session that fails to drain cannot loop.
	StopHookActive bool
}

// HookDecision is the Claude Stop hook contract. An empty struct (marshals to
// "{}") lets the session stop; Decision "block" with Reason forces one more turn.
type HookDecision struct {
	Decision string `json:"decision,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type stopPayload struct {
	CWD            string `json:"cwd"`
	SessionID      string `json:"session_id"`
	StopHookActive bool   `json:"stop_hook_active"`
}

// CheckFromReader parses a Claude Stop hook payload from stdin (cwd, session_id,
// stop_hook_active), merges any explicit CheckOptions overrides, and returns the
// hook decision. It never errors: a Stop hook must be safe, so any failure to
// resolve state yields an empty decision (let the session stop).
func CheckFromReader(r io.Reader, opts CheckOptions) HookDecision {
	if data, err := io.ReadAll(r); err == nil && len(data) > 0 {
		var p stopPayload
		if json.Unmarshal(data, &p) == nil {
			if opts.Root == "" && p.CWD != "" {
				opts.Root = p.CWD
			}
			if opts.SessionID == "" {
				opts.SessionID = p.SessionID
			}
			if p.StopHookActive {
				opts.StopHookActive = true
			}
		}
	}
	return Check(opts)
}

// Check computes the Stop-hook decision for a session: block with a drain
// instruction when captain intent is queued for the session's workflow(s).
func Check(opts CheckOptions) HookDecision {
	if opts.StopHookActive {
		return HookDecision{}
	}
	root := absRootOr(opts.Root)
	if _, err := os.Stat(filepath.Join(root, "_bridge", "inbox.jsonl")); err != nil {
		return HookDecision{}
	}
	slugs := resolveSessionSlugs(root, strings.TrimSpace(opts.Slug), strings.TrimSpace(opts.SessionID))
	if len(slugs) == 0 {
		return HookDecision{}
	}

	records, _, err := readInboxFull(filepath.Join(root, "_bridge", "inbox.jsonl"))
	if err != nil {
		return HookDecision{}
	}
	replies := loadReplies(root)
	pendingBySlug := map[string]int{}
	for _, slug := range slugs {
		members := []string{slug}
		for _, rec := range records {
			if rec.Line <= inboxCursor(root, slug) {
				continue
			}
			if !addressedTo(root, rec.routing(), slug, members) {
				continue
			}
			if replies[replyKey(rec.routing(), slug)] {
				continue
			}
			pendingBySlug[slug]++
		}
	}
	total := 0
	var pendingSlugs []string
	for _, slug := range slugs {
		if pendingBySlug[slug] > 0 {
			total += pendingBySlug[slug]
			pendingSlugs = append(pendingSlugs, slug)
		}
	}
	if total == 0 {
		return HookDecision{}
	}
	return HookDecision{Decision: "block", Reason: drainReason(total, pendingSlugs, normalizeHost(opts.Host))}
}

func drainReason(total int, slugs []string, host string) string {
	if host == "" {
		host = "claude"
	}
	slug := slugs[0]
	plural := "record"
	if total != 1 {
		plural = "records"
	}
	return "Bridge has " + strconv.Itoa(total) + " queued captain-intent " + plural +
		" in _bridge/inbox.jsonl for workflow(s): " + strings.Join(slugs, ", ") +
		". Before stopping, drain them: run `spacedock bridge inbox drain --host " + host + " --slug " + slug +
		"` (repeat per slug), act on each record, ack each with `spacedock bridge inbox ack ...`, then advance the cursor with `spacedock bridge inbox commit ...`."
}

// resolveSessionSlugs finds the workflow slug(s) this session drives. An explicit
// --slug wins. Otherwise it matches the session id against heartbeat and session
// markers so a Stop hook only ever blocks its OWN session's pending intent, never
// a sibling FO's sharing the same repo root.
func resolveSessionSlugs(root, explicitSlug, sessionID string) []string {
	if explicitSlug != "" {
		if validSlug(explicitSlug) {
			return []string{explicitSlug}
		}
		return nil
	}
	if !safeSessionID(sessionID) {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, slug := range discoverHeartbeatSlugs(root) {
		if hb, ok := loadHeartbeatAnyAge(root, slug); ok && strings.TrimSpace(hb.SessionID) == sessionID {
			if !seen[slug] {
				seen[slug] = true
				out = append(out, slug)
			}
		}
	}
	// Session markers (sessions/<actor>.json) map an actor id to a workflow; the
	// FO's own actor id is its session id on Claude.
	if slug, ok := workflowForSession(root, sessionID); ok && !seen[slug] {
		seen[slug] = true
		out = append(out, slug)
	}
	return out
}

func workflowForSession(root, sessionID string) (string, bool) {
	matches, _ := filepath.Glob(filepath.Join(root, "_bridge", "sessions", "*.json"))
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var rec sessionMarker
		if err := json.Unmarshal(data, &rec); err != nil {
			continue
		}
		if strings.TrimSpace(rec.SessionID) == sessionID && validSlug(rec.Workflow) {
			return rec.Workflow, true
		}
	}
	return "", false
}

type heartbeatOut struct {
	SessionID string `json:"session_id"`
	Host      string `json:"host,omitempty"`
	TS        string `json:"ts"`
	State     string `json:"state"`
}

func stampHeartbeat(root, slug, sessionID, host string, now time.Time) bool {
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false
	}
	hb := heartbeatOut{SessionID: sessionID, Host: host, TS: now.UTC().Format(time.RFC3339), State: "idle"}
	data, err := json.Marshal(hb)
	if err != nil {
		return false
	}
	if err := os.WriteFile(filepath.Join(dir, "fo."+slug+".json"), append(data, '\n'), 0o644); err != nil {
		return false
	}
	return true
}

// adoptCursor performs the one-time migration from the pre-versioning shared
// cursor to this workflow's own cursor, so a freshly-slugged FO does not re-drain
// (and re-apply) intent already processed under the shared cursor.
func adoptCursor(root, slug string) {
	if !validSlug(slug) {
		return
	}
	slugPath := filepath.Join(root, "_bridge", ".inbox-cursor."+slug)
	if _, err := os.Stat(slugPath); err == nil {
		return
	}
	sharedPath := filepath.Join(root, "_bridge", ".inbox-cursor")
	data, err := os.ReadFile(sharedPath)
	if err != nil {
		return
	}
	_ = os.WriteFile(slugPath, data, 0o644)
}

func readInboxFull(path string) ([]fullInboxRecord, int, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	var out []fullInboxRecord
	lineNo := 0
	scanner := lineScanner(f)
	for scanner.Scan() {
		lineNo++
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}
		var rec fullInboxRecord
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			continue
		}
		rec.Line = lineNo
		out = append(out, rec)
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return out, lineNo, nil
}

func toDrainRecord(r fullInboxRecord) DrainRecord {
	ts := ""
	if !r.TS.IsZero() {
		ts = r.TS.UTC().Format(time.RFC3339)
	}
	return DrainRecord{
		Line: r.Line, ID: r.ID, TS: ts, Kind: r.Kind, Text: r.Text, Granted: r.Granted,
		Target: r.Target, TargetSet: r.TargetSet, Entity: r.Entity, Field: r.Field,
		Value: r.Value, Verdict: r.Verdict, Directives: r.Directives, RequestID: r.RequestID,
	}
}

func resolveSessionID(host, flag string) string {
	if s := strings.TrimSpace(flag); s != "" {
		return s
	}
	if s := strings.TrimSpace(os.Getenv("SD_SESSION_ID")); s != "" {
		return s
	}
	switch host {
	case "claude":
		return strings.TrimSpace(os.Getenv("CLAUDE_CODE_SESSION_ID"))
	case "codex":
		return strings.TrimSpace(os.Getenv("CODEX_THREAD_ID"))
	}
	return ""
}

func validSlug(slug string) bool {
	return slug != "" && slug != "." && slug != ".." && safeSlugPattern.MatchString(slug)
}

func absRootOr(root string) string {
	if root == "" {
		root = "."
	}
	if abs, err := filepath.Abs(root); err == nil {
		return abs
	}
	return root
}

func nowFunc(fn func() time.Time) func() time.Time {
	if fn != nil {
		return fn
	}
	return time.Now
}

func singleLine(s string) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}
