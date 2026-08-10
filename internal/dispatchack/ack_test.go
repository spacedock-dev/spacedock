package dispatchack

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

func TestLifecycleAndReplay(t *testing.T) {
	for _, host := range []string{"claude", "codex"} {
		t.Run(host, func(t *testing.T) {
			entity := repo(t)
			rec, err := Create(entity, "n28423efmj358m5av61z2fxx", "implementation", host)
			if err != nil || rec.State != Pending || rec.Epoch == "" {
				t.Fatalf("Create() = %#v, %v", rec, err)
			}
			if _, err := Create(entity, rec.EntityID, rec.Stage, host); err == nil {
				t.Fatal("second Create succeeded")
			}
			key := "description"
			if host == "codex" {
				key = "task_name"
			}
			pre := map[string]any{"hook_event_name": "PreToolUse", "session_id": "session-1", "tool_name": "Agent", "tool_use_id": "call-1", "tool_input": map[string]any{key: "worker-sda-" + rec.Epoch}}
			hook(t, pre)
			armed := active(t, entity, rec)
			if armed.State != Armed || armed.ToolUseID != "call-1" || armed.NativeWorkerID != "" {
				t.Fatalf("armed = %#v", armed)
			}
			if code, out := hook(t, pre); code != 0 || !strings.Contains(out, `"permissionDecision":"deny"`) {
				t.Fatalf("replay hook = %d, %q", code, out)
			}
			start := map[string]any{"hook_event_name": "SubagentStart", "session_id": "session-1", "agent_id": "worker-1", "agent_type": "spacedock:ensign"}
			hook(t, start)
			consumed := active(t, entity, rec)
			if consumed.State != Consumed || consumed.NativeWorkerID != "worker-1" {
				t.Fatalf("consumed = %#v", consumed)
			}
			for _, state := range []string{Pending, Armed, Consumed} {
				ref := auditRef(rec.EntityID, rec.Stage, rec.Epoch, state)
				if out := gitTest(t, filepath.Dir(entity), "show", ref); !strings.Contains(out, `"state":"`+state+`"`) {
					t.Fatalf("audit %s = %q", state, out)
				}
			}
			if err := Clear(entity, rec.EntityID, rec.Stage); err != nil {
				t.Fatal(err)
			}
			if state, _ := State(entity, rec.EntityID, rec.Stage); state != "" {
				t.Fatalf("cleared state = %q", state)
			}
		})
	}
}
func TestDisabledAndMalformedHooksFailClosed(t *testing.T) {
	entity := repo(t)
	rec, err := Create(entity, "n28423efmj358m5av61z2fxx", "implementation", "codex")
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range []map[string]any{
		{"hook_event_name": "PreToolUse", "session_id": "s", "tool_name": "Agent", "tool_use_id": "c", "tool_input": map[string]any{"task_name": "worker-sda-00000000000000000000000000000000"}},
		{"hook_event_name": "SubagentStart", "session_id": "s", "agent_id": "worker"},
	} {
		hook(t, event)
	}
	if state, _ := State(entity, rec.EntityID, rec.Stage); state != Pending {
		t.Fatalf("malformed hooks changed pending state to %q", state)
	}
	if err := Clear(entity, rec.EntityID, rec.Stage); err == nil {
		t.Fatal("pending receipt was cleared")
	}
}
func hook(t *testing.T, event any) (int, string) {
	t.Helper()
	raw, _ := json.Marshal(event)
	var out, stderr bytes.Buffer
	code := HandleHook(bytes.NewReader(raw), &out, &stderr)
	return code, strings.TrimSpace(out.String())
}
func active(t *testing.T, entity string, want Record) Record {
	t.Helper()
	var got Record
	if err := json.Unmarshal([]byte(gitTest(t, filepath.Dir(entity), "show", activeRef(want.EntityID, want.Stage))), &got); err != nil {
		t.Fatal(err)
	}
	return got
}
func repo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	testgit.InitRepo(t, root, "-q")
	entity := filepath.Join(root, "task.md")
	if err := os.WriteFile(entity, []byte("---\nid: n28423efmj358m5av61z2fxx\nstatus: implementation\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitTest(t, root, "add", "task.md")
	gitTest(t, root, "commit", "-qm", "seed")
	return entity
}
func gitTest(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
	return string(out)
}
