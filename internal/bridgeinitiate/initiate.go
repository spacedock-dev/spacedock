// ABOUTME: FO-to-Bridge initiation writer for FO-authored feed lines and gates.
// ABOUTME: Append-only and best-effort; Bridge reads _bridge/fo-initiate.jsonl.
package bridgeinitiate

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const (
	maxHeadlineLen = 240
	maxBodyLen     = 2000

	// Truncation bounds mirror bridgeegress.truncateEvents. A still-open
	// gate-review is retained past the tail regardless of these bounds.
	defaultMaxLines  = 2000
	defaultKeepLines = 1000
	maxScanLine      = 1 << 20
)

// validKinds are the initiation decidability classes Bridge understands.
var validKinds = map[string]bool{
	"status":      true,
	"reco":        true,
	"gate-review": true,
}

// InitiationRecord is the stable JSONL shape Bridge reads from
// _bridge/fo-initiate.jsonl. The writer ALWAYS stamps status "open"; the reader
// overlays resolved/approved/rejected from decision intents.
type InitiationRecord struct {
	Schema    int    `json:"schema"`
	ID        string `json:"id"`
	TS        string `json:"ts"`
	Kind      string `json:"kind"`
	Workflow  string `json:"workflow,omitempty"`
	Entity    string `json:"entity,omitempty"`
	ShipID    string `json:"ship_id,omitempty"`
	Host      string `json:"host,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	Headline  string `json:"headline"`
	Body      string `json:"body,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Status    string `json:"status"`
}

// InitiationOptions controls one initiation append.
type InitiationOptions struct {
	Root      string
	Now       func() time.Time
	ID        string
	Kind      string
	Workflow  string
	Entity    string
	ShipID    string
	Host      string
	SessionID string
	Headline  string
	Body      string
	RequestID string
}

// Result is the compact JSON result printed for Bridge.
type Result struct {
	ID        string `json:"id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Queued    bool   `json:"queued"`
	Error     string `json:"error,omitempty"`
}

// AppendInitiation writes one open initiation under root/_bridge. Unlike alert
// writes, id is REQUIRED with no random fallback: idempotency depends on a
// stable caller-supplied fold key.
func AppendInitiation(opts InitiationOptions) (Result, error) {
	id := strings.TrimSpace(opts.ID)
	result := Result{ID: id, RequestID: id}
	if id == "" {
		return result.withError("initiation: id is required"), nil
	}

	kind := strings.TrimSpace(opts.Kind)
	if !validKinds[kind] {
		return result.withError("initiation: kind must be one of status|reco|gate-review"), nil
	}

	headline := oneLineSummary(opts.Headline, maxHeadlineLen)
	if headline == "" {
		return result.withError("initiation: headline is required"), nil
	}

	requestID := strings.TrimSpace(opts.RequestID)
	if requestID == "" && kind == "gate-review" {
		requestID = id
	}
	result.RequestID = requestID

	root := opts.Root
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return result.withError(err.Error()), nil
	}

	now := time.Now().UTC
	if opts.Now != nil {
		now = func() time.Time { return opts.Now().UTC() }
	}

	rec := InitiationRecord{
		Schema:    1,
		ID:        id,
		TS:        now().Format(time.RFC3339),
		Kind:      kind,
		Workflow:  strings.TrimSpace(opts.Workflow),
		Entity:    strings.TrimSpace(opts.Entity),
		ShipID:    strings.TrimSpace(opts.ShipID),
		Host:      strings.TrimSpace(opts.Host),
		SessionID: strings.TrimSpace(opts.SessionID),
		Headline:  headline,
		Body:      oneLineSummary(opts.Body, maxBodyLen),
		RequestID: requestID,
		Status:    "open",
	}

	dir := filepath.Join(absRoot, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return result.withError(err.Error()), nil
	}
	path := filepath.Join(dir, "fo-initiate.jsonl")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return result.withError(err.Error()), nil
	}
	data, err := json.Marshal(rec)
	if err != nil {
		_ = f.Close()
		return result.withError(err.Error()), nil
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		_ = f.Close()
		return result.withError(err.Error()), nil
	}
	if err := f.Close(); err != nil {
		return result.withError(err.Error()), nil
	}

	truncateInitiate(path)
	result.Queued = true
	return result, nil
}

func (r Result) withError(msg string) Result {
	r.Queued = false
	r.Error = msg
	return r
}

// truncateInitiate mirrors bridgeegress.truncateEvents (temp-file + rename tail
// retention) but NEVER drops the latest record of a still-open gate-review id,
// so an open gate can't be truncated out of the read window. Resolution is
// determined the way the reader determines it — a captain decision intent in the
// sibling inbox.jsonl keyed by request_id — because the writer ALWAYS stamps
// on-disk status "open", so a resolved gate's line still reads "open" here.
// Without the inbox overlay every gate would be protected forever and the file
// would grow unbounded in historical (already-decided) gate-reviews.
func truncateInitiate(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), maxScanLine)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	_ = f.Close()
	if len(lines) <= defaultMaxLines || scanner.Err() != nil {
		return
	}

	keepLines := defaultKeepLines
	if keepLines > len(lines) {
		keepLines = len(lines)
	}
	tailStart := len(lines) - keepLines

	resolved := loadResolvedGateRequestIDs(filepath.Dir(path))

	// Latest record per id wins; retain the latest line of every still-open
	// gate-review id even when it falls before the tail window. A gate whose
	// request_id has a captain decision in the inbox is resolved and evictable.
	latestLine := make(map[string]int)
	for i, line := range lines {
		rec, ok := parseRecord(line)
		if !ok {
			continue
		}
		latestLine[rec.ID] = i
	}
	protected := make(map[int]bool)
	for _, i := range latestLine {
		rec, ok := parseRecord(lines[i])
		if !ok {
			continue
		}
		if rec.Kind == "gate-review" && rec.Status == "open" && !resolved[gateResolutionKey(rec)] && i < tailStart {
			protected[i] = true
		}
	}

	var kept []string
	for i := 0; i < tailStart; i++ {
		if protected[i] {
			kept = append(kept, lines[i])
		}
	}
	kept = append(kept, lines[tailStart:]...)

	out := strings.Join(kept, "\n") + "\n"
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp.*")
	if err != nil {
		return
	}
	tmpPath := tmp.Name()
	if _, err := tmp.WriteString(out); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
	}
}

func parseRecord(line string) (InitiationRecord, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return InitiationRecord{}, false
	}
	var rec InitiationRecord
	if err := json.Unmarshal([]byte(line), &rec); err != nil {
		return InitiationRecord{}, false
	}
	return rec, true
}

// gateResolutionKey is the request_id a decision intent correlates against,
// falling back to the id (the writer defaults request_id to id for gate-review).
func gateResolutionKey(rec InitiationRecord) string {
	if rec.RequestID != "" {
		return rec.RequestID
	}
	return rec.ID
}

// loadResolvedGateRequestIDs reads the sibling inbox.jsonl and returns the set
// of request_ids carried by captain "decision" intents — the same overlay the
// Bridge reader uses to flip a gate-review to resolved. Best-effort: a missing
// or unreadable inbox yields an empty set, so every open gate stays protected
// (the safe fallback — never evict a gate that might still be open).
func loadResolvedGateRequestIDs(dir string) map[string]bool {
	resolved := make(map[string]bool)
	f, err := os.Open(filepath.Join(dir, "inbox.jsonl"))
	if err != nil {
		return resolved
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), maxScanLine)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var rec struct {
			Kind      string `json:"kind"`
			RequestID string `json:"request_id"`
		}
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		if rec.Kind == "decision" && rec.RequestID != "" {
			resolved[rec.RequestID] = true
		}
	}
	return resolved
}

// oneLineSummary collapses control chars/whitespace to single spaces and bounds
// the result to maxLen runes. Mirrors bridgealert.oneLineSummary.
func oneLineSummary(in string, maxLen int) string {
	in = strings.TrimSpace(in)
	var b strings.Builder
	prevSpace := false
	for _, r := range in {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
			continue
		}
		b.WriteRune(r)
		prevSpace = false
		if b.Len() >= maxLen {
			break
		}
	}
	return strings.TrimSpace(b.String())
}
