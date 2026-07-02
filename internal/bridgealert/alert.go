// ABOUTME: FO-to-Bridge alert writer for top-level captain interrupts.
// ABOUTME: Append-only and best-effort; Bridge reads the resulting JSONL.
package bridgealert

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const (
	maxReasonLen  = 240
	maxCommandLen = 500
)

// PermissionAlert is the stable JSONL shape Bridge reads from
// _bridge/fo-alerts.jsonl.
type PermissionAlert struct {
	Schema     int      `json:"schema"`
	ID         string   `json:"id"`
	TS         string   `json:"ts"`
	Kind       string   `json:"kind"`
	Severity   string   `json:"severity"`
	Workflow   string   `json:"workflow,omitempty"`
	Entity     string   `json:"entity,omitempty"`
	Host       string   `json:"host,omitempty"`
	SessionID  string   `json:"session_id,omitempty"`
	Reason     string   `json:"reason"`
	Command    string   `json:"command,omitempty"`
	PrefixRule []string `json:"prefix_rule,omitempty"`
	Status     string   `json:"status"`
}

// PermissionOptions controls one permission-alert append.
type PermissionOptions struct {
	Root       string
	Now        func() time.Time
	ID         string
	Workflow   string
	Entity     string
	Host       string
	SessionID  string
	Reason     string
	Command    string
	PrefixRule []string
}

type Result struct {
	ID        string `json:"id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Queued    bool   `json:"queued"`
	Error     string `json:"error,omitempty"`
}

// AppendPermission writes one open permission alert under root/_bridge.
func AppendPermission(opts PermissionOptions) (Result, error) {
	id := strings.TrimSpace(opts.ID)
	if id == "" {
		id = fallbackID()
		if generated, err := newID(); err == nil {
			id = generated
		}
	}
	result := Result{ID: id, RequestID: id}
	root := opts.Root
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return result.withError(err.Error()), nil
	}
	reason := oneLineSummary(opts.Reason, maxReasonLen)
	if reason == "" {
		return result.withError("permission alert: reason is required"), nil
	}
	command := oneLineSummary(opts.Command, maxCommandLen)
	now := time.Now().UTC
	if opts.Now != nil {
		now = func() time.Time { return opts.Now().UTC() }
	}
	alert := PermissionAlert{
		Schema:     1,
		ID:         id,
		TS:         now().Format(time.RFC3339),
		Kind:       "permission-request",
		Severity:   "blocked",
		Workflow:   strings.TrimSpace(opts.Workflow),
		Entity:     strings.TrimSpace(opts.Entity),
		Host:       strings.TrimSpace(opts.Host),
		SessionID:  strings.TrimSpace(opts.SessionID),
		Reason:     reason,
		Command:    command,
		PrefixRule: cleanedPrefixRule(opts.PrefixRule),
		Status:     "open",
	}
	dir := filepath.Join(absRoot, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return result.withError(err.Error()), nil
	}
	f, err := os.OpenFile(filepath.Join(dir, "fo-alerts.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return result.withError(err.Error()), nil
	}
	defer func() { _ = f.Close() }()
	data, err := json.Marshal(alert)
	if err != nil {
		return result.withError(err.Error()), nil
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		return result.withError(err.Error()), nil
	}
	result.Queued = true
	return result, nil
}

func (r Result) withError(msg string) Result {
	r.Queued = false
	r.Error = msg
	return r
}

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

func cleanedPrefixRule(in []string) []string {
	out := make([]string, 0, len(in))
	for _, part := range in {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func newID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "perm_" + hex.EncodeToString(b[:]), nil
}

func fallbackID() string {
	return "perm_" + hex.EncodeToString([]byte(time.Now().UTC().Format("20060102150405.000000000")))
}
