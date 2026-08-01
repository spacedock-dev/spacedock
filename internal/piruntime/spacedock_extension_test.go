package piruntime

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSpacedockPiExtensionBootstrapBehavior(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required to execute the Pi extension behavior harness")
	}

	repoRoot := findRepoRoot(t)
	extensionPath := filepath.Join(repoRoot, ".pi", "extensions", "spacedock.ts")
	extensionSource, err := os.ReadFile(extensionPath)
	if err != nil {
		t.Fatalf("read extension: %v", err)
	}

	tmp := t.TempDir()
	modulePath := filepath.Join(tmp, "spacedock-extension.mjs")
	if err := os.WriteFile(modulePath, extensionSource, 0o644); err != nil {
		t.Fatalf("write harness module: %v", err)
	}

	harnessPath := filepath.Join(tmp, "harness.mjs")
	harness := `
import register from './spacedock-extension.mjs';

const handlers = new Map();
const pi = { on(event, handler) { handlers.set(event, handler); } };
register(pi);

const marker = 'SPACEDOCK-FO-BOOTSTRAP-v1';
const textOf = (message) => Array.isArray(message?.content)
  ? message.content.filter((part) => part?.type === 'text').map((part) => String(part.text ?? '')).join('')
  : String(message?.content ?? '');
const countBootstraps = (messages) => messages.filter((message) => message?.role === 'user' && textOf(message).startsWith('<EXTREMELY_IMPORTANT>') && textOf(message).includes(marker)).length;
const assert = (condition, message) => { if (!condition) throw new Error(message); };

const resources = handlers.get('resources_discover')({ type: 'resources_discover', cwd: process.cwd(), reason: 'test' });
assert(resources.skillPaths.length === 1 && resources.skillPaths[0].endsWith('/skills'), 'resources_discover returns the package skills directory');

handlers.get('session_start')({ type: 'session_start' });
let first = handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'hello' }] }] });
assert(countBootstraps(first.messages) === 1, 'session_start injects exactly one structural bootstrap');
assert(textOf(first.messages[0]).includes('Load the $spacedock:first-officer skill'), 'bootstrap points at the shipped first-officer contract');
assert(textOf(first.messages[0]).includes('Pi tool mapping: read/write/edit/bash/grep/find/ls'), 'bootstrap includes the Pi tool mapping');

handlers.get('session_compact')({ type: 'session_compact' });
const compactSummary = { role: 'assistant', content: [{ type: 'text', text: 'Compaction summary mentions SPACEDOCK-FO-BOOTSTRAP-v1 but is not the bootstrap message.' }] };
let afterCompact = handlers.get('context')({ messages: [compactSummary, { role: 'user', content: [{ type: 'text', text: 'continue' }] }] });
assert(countBootstraps(afterCompact.messages) === 1, 'summary marker mention does not suppress reinjection');
assert(afterCompact.messages[0] === compactSummary, 'bootstrap is inserted after the leading compaction summary');
assert(textOf(afterCompact.messages[1]).startsWith('<EXTREMELY_IMPORTANT>'), 'bootstrap is the first non-summary message after compaction');

handlers.get('session_compact')({ type: 'session_compact' });
let deduped = handlers.get('context')({ messages: afterCompact.messages });
assert(deduped === undefined, 'existing structural bootstrap de-duplicates context injection');

handlers.get('agent_end')({ type: 'agent_end' });
let suppressed = handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'new turn' }] }] });
assert(suppressed === undefined, 'agent_end suppresses further bootstrap injection');
`
	if err := os.WriteFile(harnessPath, []byte(harness), 0o644); err != nil {
		t.Fatalf("write harness: %v", err)
	}

	cmd := exec.Command(node, harnessPath)
	cmd.Dir = repoRoot
	// The unset marker run must be genuinely unset: a pi-subagents child shell
	// exports PI_SUBAGENT_CHILD, and an inherited marker would silently flip the
	// exemption inside the harness process.
	cmd.Env = harnessEnv(nil, "PI_SUBAGENT_CHILD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("extension behavior harness failed: %v\n%s", err, out)
	}
}

// TestSpacedockPiExtensionChildExemption is the AC-1 child-session half: a
// pi-subagents child (PI_SUBAGENT_CHILD=1) is a delegated worker, not a first
// officer, so the context hook must inject zero FO bootstrap across
// session_start and session_compact. Skill discovery is unaffected.
func TestSpacedockPiExtensionChildExemption(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required to execute the Pi extension behavior harness")
	}

	repoRoot := findRepoRoot(t)
	tmp := t.TempDir()
	copyExtension(t, repoRoot, tmp)

	harnessPath := filepath.Join(tmp, "harness-child.mjs")
	harness := `
import register from './spacedock-extension.mjs';

const handlers = new Map();
const pi = { on(event, handler) { handlers.set(event, handler); } };
register(pi);

const assert = (condition, message) => { if (!condition) throw new Error(message); };
assert(process.env.PI_SUBAGENT_CHILD === '1', 'harness must run under the subagent-child marker');

const resources = handlers.get('resources_discover')({ type: 'resources_discover', cwd: process.cwd(), reason: 'test' });
assert(resources.skillPaths.length === 1 && resources.skillPaths[0].endsWith('/skills'), 'child sessions still discover the package skills directory');

handlers.get('session_start')({ type: 'session_start' });
const afterStart = handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'hello' }] }] });
assert(afterStart === undefined, 'PI_SUBAGENT_CHILD=1 session_start injects zero FO bootstrap');

handlers.get('session_compact')({ type: 'session_compact' });
const afterCompact = handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'continue' }] }] });
assert(afterCompact === undefined, 'PI_SUBAGENT_CHILD=1 session_compact injects zero FO bootstrap');
`
	if err := os.WriteFile(harnessPath, []byte(harness), 0o644); err != nil {
		t.Fatalf("write child harness: %v", err)
	}

	cmd := exec.Command(node, harnessPath)
	cmd.Dir = repoRoot
	cmd.Env = harnessEnv(map[string]string{"PI_SUBAGENT_CHILD": "1"})
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("child-exemption harness failed: %v\n%s", err, out)
	}
}

// copyExtension stages the shipped .pi/extensions/spacedock.ts as a node
// module next to the harness scripts in tmp.
func copyExtension(t *testing.T, repoRoot, tmp string) {
	t.Helper()
	extensionSource, err := os.ReadFile(filepath.Join(repoRoot, ".pi", "extensions", "spacedock.ts"))
	if err != nil {
		t.Fatalf("read extension: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "spacedock-extension.mjs"), extensionSource, 0o644); err != nil {
		t.Fatalf("write harness module: %v", err)
	}
}

// harnessEnv returns the process environment with drop keys removed and extra
// overrides applied, so each harness run pins the marker state explicitly.
func harnessEnv(extra map[string]string, drop ...string) []string {
	dropped := map[string]bool{}
	for _, key := range drop {
		dropped[key] = true
	}
	var env []string
	for _, kv := range os.Environ() {
		key, _, _ := strings.Cut(kv, "=")
		if dropped[key] {
			continue
		}
		if _, ok := extra[key]; ok {
			continue
		}
		env = append(env, kv)
	}
	for key, value := range extra {
		env = append(env, key+"="+value)
	}
	return env
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir || strings.TrimSpace(dir) == "" {
			t.Fatal("could not find repo root")
		}
		dir = parent
	}
}
