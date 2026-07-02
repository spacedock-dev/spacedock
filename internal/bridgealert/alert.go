// ABOUTME: FO-to-Bridge alert writer for top-level captain interrupts.
// ABOUTME: Append-only and best-effort; Bridge reads the resulting JSONL.
package bridgealert

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	ID     string `json:"id"`
	Queued bool   `json:"queued"`
}

// AppendPermission writes one open permission alert under root/_bridge.
func AppendPermission(opts PermissionOptions) (Result, error) {
	root := opts.Root
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Result{}, err
	}
	reason := strings.TrimSpace(opts.Reason)
	if reason == "" {
		return Result{}, fmt.Errorf("permission alert: reason is required")
	}
	id := strings.TrimSpace(opts.ID)
	if id == "" {
		id, err = newID()
		if err != nil {
			return Result{}, err
		}
	}
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
		Command:    strings.TrimSpace(opts.Command),
		PrefixRule: cleanedPrefixRule(opts.PrefixRule),
		Status:     "open",
	}
	dir := filepath.Join(absRoot, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Result{}, err
	}
	f, err := os.OpenFile(filepath.Join(dir, "fo-alerts.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = f.Close() }()
	data, err := json.Marshal(alert)
	if err != nil {
		return Result{}, err
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		return Result{}, err
	}
	return Result{ID: id, Queued: true}, nil
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
