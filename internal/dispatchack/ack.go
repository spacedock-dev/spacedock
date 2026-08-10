// Package dispatchack binds one fresh CLI dispatch to one native worker start.
package dispatchack

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const Pending, Armed, Consumed = "pending", "armed", "consumed"

type Record struct {
	State          string `json:"state"`
	EntityID       string `json:"entity_id"`
	EntityPath     string `json:"entity_path"`
	Stage          string `json:"stage"`
	Host           string `json:"host"`
	Epoch          string `json:"epoch"`
	HostSessionID  string `json:"host_session_id,omitempty"`
	ToolUseID      string `json:"tool_use_id,omitempty"`
	NativeWorkerID string `json:"native_worker_id,omitempty"`
}
type lookup struct {
	Repo string `json:"repo"`
	Record
}

func (r Record) Ref() string { return activeRef(r.EntityID, r.Stage) }
func (r Record) Name(base string) string {
	suffix := "-sda-" + r.Epoch
	if len(base)+len(suffix) > 64 {
		base = strings.TrimRight(base[:64-len(suffix)], "-")
	}
	return base + suffix
}

var tokenRE = regexp.MustCompile(`sda-([0-9a-f]{32})`)

func Create(entityPath, entityID, stage, host string) (Record, error) {
	repo, rel, err := repository(entityPath)
	if err != nil {
		return Record{}, err
	}
	b := make([]byte, 16)
	if _, err = rand.Read(b); err != nil {
		return Record{}, err
	}
	r := Record{State: Pending, EntityID: entityID, EntityPath: rel, Stage: stage, Host: host, Epoch: hex.EncodeToString(b)}
	if err = transition(repo, "", r); err != nil {
		return Record{}, fmt.Errorf("dispatch acknowledgment already exists or cannot be created: %w", err)
	}
	if err = os.MkdirAll(lookupDir(), 0o700); err == nil {
		var f *os.File
		f, err = os.OpenFile(lookupPath(r.Epoch), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			err = json.NewEncoder(f).Encode(lookup{Repo: repo, Record: r})
			if closeErr := f.Close(); err == nil {
				err = closeErr
			}
		}
	}
	return r, err
}
func State(entityPath, entityID, stage string) (string, error) {
	repo, _, err := repository(entityPath)
	if err != nil {
		return "", nil
	}
	if r, _, err := read(repo, activeRef(entityID, stage)); err == nil {
		return r.State, nil
	}
	return "", nil
}
func Clear(entityPath, entityID, stage string) error {
	repo, _, err := repository(entityPath)
	if err != nil {
		return err
	}
	r, oid, err := read(repo, activeRef(entityID, stage))
	if err != nil || r.State != Consumed {
		return fmt.Errorf("dispatch acknowledgment is not consumed")
	}
	if _, err = git(repo, "update-ref", "-d", activeRef(entityID, stage), oid); err == nil {
		_ = os.Remove(lookupPath(r.Epoch))
	}
	return err
}
func HandleHook(stdin io.Reader, stdout, stderr io.Writer) int {
	var e struct {
		Event   string         `json:"hook_event_name"`
		Session string         `json:"session_id"`
		ToolID  string         `json:"tool_use_id"`
		AgentID string         `json:"agent_id"`
		Input   map[string]any `json:"tool_input"`
	}
	if err := json.NewDecoder(stdin).Decode(&e); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	switch e.Event {
	case "PreToolUse":
		epoch := eventEpoch(e.Input)
		if epoch == "" {
			return 0
		}
		if err := arm(epoch, e.Session, e.ToolID); err != nil {
			deny(stdout, err.Error())
		}
	case "SubagentStart":
		if err := consume(e.Session, e.AgentID); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	return 0
}
func arm(epoch, session, toolID string) error {
	l, err := loadLookup(epoch)
	if err != nil || session == "" || toolID == "" {
		return fmt.Errorf("dispatch acknowledgment does not match one pending envelope")
	}
	r, oid, err := read(l.Repo, activeRef(l.EntityID, l.Stage))
	if err != nil || r.State != Pending || r.Epoch != epoch {
		return fmt.Errorf("dispatch acknowledgment is stale or already armed")
	}
	r.State, r.HostSessionID, r.ToolUseID = Armed, session, toolID
	return transition(l.Repo, oid, r)
}
func consume(session, worker string) error {
	entries, _ := os.ReadDir(lookupDir())
	var found *lookup
	var oid string
	for _, entry := range entries {
		l, err := loadLookup(strings.TrimSuffix(entry.Name(), ".json"))
		if err != nil {
			continue
		}
		r, candidate, err := read(l.Repo, activeRef(l.EntityID, l.Stage))
		if err == nil && r.State == Armed && r.HostSessionID == session {
			if found != nil {
				return fmt.Errorf("more than one armed dispatch acknowledgment matches this session")
			}
			l.Record, found, oid = r, &l, candidate
		}
	}
	if found == nil {
		return nil
	}
	if worker == "" {
		return fmt.Errorf("SubagentStart did not supply a native worker ID")
	}
	found.State, found.NativeWorkerID = Consumed, worker
	return transition(found.Repo, oid, found.Record)
}
func transition(repo, old string, r Record) error {
	data, _ := json.Marshal(r)
	oid, err := gitInput(repo, data, "hash-object", "-w", "--stdin")
	if err != nil {
		return err
	}
	if old == "" {
		old = strings.Repeat("0", 40)
	}
	cmd := fmt.Sprintf("start\nupdate %s %s %s\ncreate %s %s\nprepare\ncommit\n", activeRef(r.EntityID, r.Stage), oid, old, auditRef(r.EntityID, r.Stage, r.Epoch, r.State), oid)
	_, err = gitInput(repo, []byte(cmd), "update-ref", "--stdin")
	return err
}
func read(repo, ref string) (Record, string, error) {
	oid, err := git(repo, "rev-parse", "--verify", ref)
	if err != nil {
		return Record{}, "", err
	}
	raw, err := git(repo, "cat-file", "blob", oid)
	var r Record
	if err == nil {
		err = json.Unmarshal([]byte(raw), &r)
	}
	return r, oid, err
}
func repository(entity string) (string, string, error) {
	if resolved, err := filepath.EvalSymlinks(entity); err == nil {
		entity = resolved
	}
	repo, err := git(filepath.Dir(entity), "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", err
	}
	rel, err := filepath.Rel(repo, entity)
	return repo, filepath.ToSlash(rel), err
}
func eventEpoch(input map[string]any) string {
	for _, key := range []string{"description", "task_name", "name", "prompt"} {
		if value, ok := input[key].(string); ok {
			if match := tokenRE.FindStringSubmatch(value); len(match) == 2 {
				return match[1]
			}
		}
	}
	return ""
}
func deny(w io.Writer, reason string) {
	fmt.Fprintf(w, `{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":%q}}`, reason)
}
func activeRef(id, stage string) string { return "refs/spacedock/dispatch-ack/" + id + "/" + stage }
func auditRef(id, stage, epoch, state string) string {
	return "refs/spacedock/dispatch-ack-audit/" + id + "/" + stage + "/" + epoch + "/" + state
}
func lookupDir() string              { return fmt.Sprint(os.TempDir(), "/spacedock-dispatch-ack-", os.Getuid()) }
func lookupPath(epoch string) string { return filepath.Join(lookupDir(), epoch+".json") }
func loadLookup(epoch string) (lookup, error) {
	var l lookup
	raw, err := os.ReadFile(lookupPath(epoch))
	if err == nil {
		err = json.Unmarshal(raw, &l)
	}
	return l, err
}
func git(repo string, args ...string) (string, error) { return gitInput(repo, nil, args...) }
func gitInput(repo string, input []byte, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir, cmd.Stdin = repo, bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
