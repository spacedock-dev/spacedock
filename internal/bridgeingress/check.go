// ABOUTME: Synchronous Stop-hook inbox check for FO sessions.
// ABOUTME: Blocks a stopping session for one more turn when captain intent is queued; never errors.
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

// Check is the read-only half of the seam that used to sit inside the
// drain/ack/commit family: it decides whether a stopping FO session still has
// queued captain intent and, if so, forces one more turn so the FO drains it.
// The FO does the draining itself by writing _bridge/ files directly per
// docs/seam-contract.md §3 — there is no drain/ack/commit verb. Check shares the
// inbox/reply/heartbeat readers with wake.go (inboxCursor, loadReplies,
// replyKey, targetsFor, discoverHeartbeatSlugs, loadHeartbeatAnyAge,
// safeSessionID, lineScanner, safeSlugPattern).

// fullInboxRecord is the full-field parse of one inbox line. Check only needs
// its line number and routing projection, but readInboxFull returns this shape,
// so the type and its routing() projection travel with the reader.
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

// CheckOptions controls a Stop-hook drain check.
type CheckOptions struct {
	Host      string
	Root      string
	Slug      string
	SessionID string
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
		cursor := inboxCursor(root, slug)
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
	return HookDecision{Decision: "block", Reason: drainReason(total, pendingSlugs)}
}

// drainReason is the Stop-hook block reason. It instructs the FO to drain by
// writing _bridge/ files directly (docs/seam-contract.md §3 / the bridge-seam
// mod) and names no `spacedock bridge inbox` verb — draining, acking, and cursor
// advancement are FO-authored file writes, not CLI verbs (findings B.2).
func drainReason(total int, slugs []string) string {
	plural := "record"
	if total != 1 {
		plural = "records"
	}
	return "Bridge has " + strconv.Itoa(total) + " queued captain-intent " + plural +
		" in _bridge/inbox.jsonl for workflow(s): " + strings.Join(slugs, ", ") +
		". Before stopping, drain them by writing _bridge/ files directly, per the bridge-seam mod. For each slug: read _bridge/inbox.jsonl lines past _bridge/.inbox-cursor.<slug>, act on each addressed intent, append a terminal ack line to _bridge/fo-replies.jsonl, write the inbox's new physical line count to _bridge/.inbox-cursor.<slug> (monotonic whole-file replace), and refresh the heartbeat _bridge/fo.<slug>.json."
}

// resolveSessionSlugs finds the workflow slug(s) this session drives. An explicit
// --slug wins. Otherwise it matches the session id against heartbeat and session
// markers so a Stop hook only ever blocks its OWN session's pending intent, never
// a sibling FO's sharing the same repo root. A heartbeat that omits its
// session_id therefore resolves nothing (review B2): the check never blocks and
// the intent sits queued until the heartbeat is written with the correct id.
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
		// Preserve the exact on-disk ts string so a downstream reader can round-trip
		// it verbatim rather than reformatting to UTC.
		var raw struct {
			TS string `json:"ts"`
		}
		_ = json.Unmarshal(scanner.Bytes(), &raw)
		rec.RawTS = raw.TS
		rec.Line = lineNo
		out = append(out, rec)
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return out, lineNo, nil
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
