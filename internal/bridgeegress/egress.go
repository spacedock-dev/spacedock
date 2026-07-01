// ABOUTME: Host-neutral Bridge egress writer for events.jsonl and session markers.
// ABOUTME: Observe-only: malformed payloads and filesystem failures degrade to no-op.
package bridgeegress

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/status"
)

const (
	defaultMaxLines  = 2000
	defaultKeepLines = 1000
)

var safeIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// Options controls one observe-only egress emission.
type Options struct {
	Host      string
	CWD       string
	Now       func() time.Time
	MaxLines  int
	KeepLines int
}

// Detail is the intentionally small, non-sensitive event detail block.
type Detail struct {
	Tool   string `json:"tool"`
	Source string `json:"source"`
}

// Event is one normalized Bridge activity line.
type Event struct {
	Timestamp string `json:"timestamp"`
	TS        string `json:"ts"`
	Host      string `json:"host"`
	Event     string `json:"event"`
	SessionID string `json:"session_id"`
	AgentID   string `json:"agent_id"`
	AgentType string `json:"agent_type"`
	ActorID   string `json:"actor_id"`
	Detail    Detail `json:"detail"`
}

// Marker maps a host actor to the workflow entity it is driving.
type Marker struct {
	Host      string `json:"host"`
	SessionID string `json:"session_id"`
	AgentID   string `json:"agent_id"`
	ActorID   string `json:"actor_id"`
	Entity    string `json:"entity"`
	Workflow  string `json:"workflow"`
}

type payload struct {
	CWD           string `json:"cwd"`
	Host          string `json:"host"`
	Event         string `json:"event"`
	HookEventName string `json:"hook_event_name"`
	SessionID     string `json:"session_id"`
	AgentID       string `json:"agent_id"`
	AgentType     string `json:"agent_type"`
	ActorID       string `json:"actor_id"`
	ToolName      string `json:"tool_name"`
	Source        string `json:"source"`
	FilePath      string `json:"file_path"`
	EntityPath    string `json:"entity_path"`
	Timestamp     string `json:"timestamp"`
	TS            string `json:"ts"`
	Detail        struct {
		Tool     string `json:"tool"`
		Source   string `json:"source"`
		FilePath string `json:"file_path"`
	} `json:"detail"`
	ToolInput struct {
		FilePath string `json:"file_path"`
	} `json:"tool_input"`
}

// EmitFromReader reads one JSON payload and writes normalized egress files. It
// never returns an operational error: Bridge egress is telemetry and must not
// break the host session.
func EmitFromReader(r io.Reader, opts Options) {
	data, err := io.ReadAll(r)
	if err != nil {
		return
	}
	Emit(data, opts)
}

// Emit writes one normalized event and, when the payload names an entity file
// for a child actor, a first-write-wins session marker.
func Emit(data []byte, opts Options) {
	var p payload
	if err := json.Unmarshal(data, &p); err != nil {
		return
	}

	host := firstNonEmpty(opts.Host, p.Host)
	eventName := firstNonEmpty(p.Event, p.HookEventName)
	cwd := firstNonEmpty(p.CWD, opts.CWD)
	if host == "" || eventName == "" || cwd == "" {
		return
	}
	cwdAbs, err := filepath.Abs(cwd)
	if err != nil {
		return
	}

	actorID := actorIDFor(host, p.SessionID, p.AgentID, p.ActorID)
	ts := timestampFor(p, opts)
	line := Event{
		Timestamp: ts,
		TS:        ts,
		Host:      host,
		Event:     eventName,
		SessionID: p.SessionID,
		AgentID:   p.AgentID,
		AgentType: p.AgentType,
		ActorID:   actorID,
		Detail: Detail{
			Tool:   firstNonEmpty(p.Detail.Tool, p.ToolName),
			Source: firstNonEmpty(p.Detail.Source, p.Source),
		},
	}

	bridgeDir := filepath.Join(cwdAbs, "_bridge")
	if err := os.MkdirAll(bridgeDir, 0o755); err != nil {
		return
	}
	eventsPath := filepath.Join(bridgeDir, "events.jsonl")
	appendJSONLine(eventsPath, line)
	truncateEvents(eventsPath, opts)

	if actorID == "" || !isChildActor(p) {
		return
	}
	entityPath, ok := markerEntityPath(host, p)
	if !ok {
		return
	}
	workflow, entity, ok := DeriveEntity(cwdAbs, entityPath)
	if !ok {
		return
	}
	writeMarker(filepath.Join(bridgeDir, "sessions", actorID+".json"), Marker{
		Host:      host,
		SessionID: p.SessionID,
		AgentID:   p.AgentID,
		ActorID:   actorID,
		Entity:    entity,
		Workflow:  workflow,
	})
}

// DeriveEntity maps a read entity file path to (workflow, entity). It supports
// legacy docs/spacedock layouts and workflow-local split-root state checkouts.
func DeriveEntity(cwd, entityPath string) (string, string, bool) {
	if strings.TrimSpace(entityPath) == "" {
		return "", "", false
	}

	cleanInput := filepath.Clean(entityPath)
	if hasArchiveSegment(cleanInput) {
		return "", "", false
	}

	abs := cleanInput
	if !filepath.IsAbs(abs) {
		if cwd == "" {
			return "", "", false
		}
		abs = filepath.Join(cwd, cleanInput)
	}
	abs, err := filepath.Abs(abs)
	if err != nil {
		return "", "", false
	}
	if cwd != "" {
		if _, ok := pathRelInside(cwd, abs); !ok {
			return "", "", false
		}
	}
	if workflow, entity, ok := deriveLegacyDocsSpacedock(abs); ok {
		return workflow, entity, true
	}
	if workflow, entity, ok := deriveFromReadmeState(abs); ok {
		return workflow, entity, true
	}
	return deriveDotStateSegment(abs)
}

func appendJSONLine(path string, value any) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	_, _ = f.Write(append(data, '\n'))
}

func writeMarker(path string, marker Marker) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := json.Marshal(marker)
	if err != nil {
		return
	}
	_, _ = f.Write(append(data, '\n'))
}

func truncateEvents(path string, opts Options) {
	maxLines := opts.MaxLines
	if maxLines <= 0 {
		maxLines = defaultMaxLines
	}
	keepLines := opts.KeepLines
	if keepLines <= 0 {
		keepLines = defaultKeepLines
	}
	if keepLines > maxLines {
		keepLines = maxLines
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	_ = f.Close()
	if len(lines) <= maxLines || scanner.Err() != nil {
		return
	}
	if keepLines > len(lines) {
		keepLines = len(lines)
	}
	kept := strings.Join(lines[len(lines)-keepLines:], "\n") + "\n"
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp.*")
	if err != nil {
		return
	}
	tmpPath := tmp.Name()
	if _, err := tmp.WriteString(kept); err != nil {
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

func timestampFor(p payload, opts Options) string {
	if p.Timestamp != "" {
		return p.Timestamp
	}
	if p.TS != "" {
		return p.TS
	}
	now := time.Now
	if opts.Now != nil {
		now = opts.Now
	}
	return now().UTC().Format(time.RFC3339)
}

func actorIDFor(host, sessionID, agentID, explicit string) string {
	if explicit != "" {
		if safeID(explicit) {
			return explicit
		}
		return ""
	}
	if host == "claude" {
		if safeID(sessionID) {
			return sessionID
		}
		return ""
	}
	if safeID(sessionID) && safeID(agentID) {
		return sessionID + "." + agentID
	}
	if safeID(sessionID) {
		return sessionID
	}
	if safeID(agentID) {
		return agentID
	}
	return ""
}

func isChildActor(p payload) bool {
	if p.AgentID != "" {
		return true
	}
	return strings.Contains(strings.ToLower(p.AgentType), "ensign")
}

func markerEntityPath(host string, p payload) (string, bool) {
	if p.EntityPath != "" {
		return p.EntityPath, true
	}
	if host == "claude" && p.HookEventName == "PostToolUse" && p.ToolName == "Read" {
		path := firstNonEmpty(p.ToolInput.FilePath, p.FilePath, p.Detail.FilePath)
		return path, path != ""
	}
	return "", false
}

func deriveLegacyDocsSpacedock(path string) (string, string, bool) {
	segments := pathSegments(path)
	for i := 0; i+2 < len(segments); i++ {
		if segments[i] != "docs" || segments[i+1] != "spacedock" {
			continue
		}
		workflow := segments[i+2]
		rel := segments[i+3:]
		if len(rel) > 0 && rel[0] == ".spacedock-state" {
			rel = rel[1:]
		}
		if entity, ok := entitySlugFromRel(rel); ok && safeID(workflow) {
			return workflow, entity, true
		}
	}
	return "", "", false
}

func deriveFromReadmeState(absPath string) (string, string, bool) {
	dir := filepath.Dir(absPath)
	for {
		readme := filepath.Join(dir, "README.md")
		if isRegularFile(readme) {
			mode, relPath, err := status.ClassifyState(status.ParseFrontmatter(readme)["state"])
			if err == nil && mode == status.StateSplitRoot {
				stateRoot := filepath.Join(dir, relPath)
				if rel, ok := pathRelInside(stateRoot, absPath); ok {
					if entity, ok := entitySlugFromRel(pathSegments(rel)); ok {
						workflow := filepath.Base(dir)
						if safeID(workflow) {
							return workflow, entity, true
						}
					}
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", false
		}
		dir = parent
	}
}

func deriveDotStateSegment(absPath string) (string, string, bool) {
	segments := pathSegments(absPath)
	for i := 1; i < len(segments); i++ {
		if segments[i] != ".spacedock-state" {
			continue
		}
		workflow := segments[i-1]
		rel := segments[i+1:]
		if entity, ok := entitySlugFromRel(rel); ok && safeID(workflow) {
			return workflow, entity, true
		}
	}
	return "", "", false
}

func entitySlugFromRel(rel []string) (string, bool) {
	if len(rel) == 1 && strings.HasSuffix(rel[0], ".md") && rel[0] != "README.md" {
		slug := strings.TrimSuffix(rel[0], ".md")
		return slug, safeID(slug)
	}
	if len(rel) == 2 && rel[1] == "index.md" {
		return rel[0], safeID(rel[0])
	}
	return "", false
}

func pathRelInside(root, target string) (string, bool) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", false
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return "", false
	}
	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return rel, true
}

func hasArchiveSegment(path string) bool {
	for _, seg := range pathSegments(path) {
		if seg == "_archive" {
			return true
		}
	}
	return false
}

func pathSegments(path string) []string {
	clean := filepath.ToSlash(filepath.Clean(path))
	raw := strings.Split(clean, "/")
	segments := raw[:0]
	for _, seg := range raw {
		if seg != "" && seg != "." {
			segments = append(segments, seg)
		}
	}
	return segments
}

func safeID(s string) bool {
	return safeIDPattern.MatchString(s)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
