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
const execCalls = [];
const pi = {
  on(event, handler) { handlers.set(event, handler); },
  exec(command, args, options) {
    execCalls.push({ command, args, options });
    return Promise.resolve({ stdout: '{"command":"boot","mods":{},"ready_gates":[],"state_backend":"single-root"}', stderr: '', code: 0, killed: false });
  }
};
register(pi);

const marker = 'SPACEDOCK-FO-BOOTSTRAP-v1';
const bootRecordMarker = '[SPACEDOCK-FO-BOOT-v2]';
const textOf = (message) => Array.isArray(message?.content)
  ? message.content.filter((part) => part?.type === 'text').map((part) => String(part.text ?? '')).join('')
  : String(message?.content ?? '');
const countBootstraps = (messages) => messages.filter((message) => message?.role === 'user' && textOf(message).startsWith('<EXTREMELY_IMPORTANT>') && textOf(message).includes(marker)).length;
const countBootRecords = (messages) => messages.filter((message) => message?.role === 'user' && textOf(message).includes(bootRecordMarker)).length;
const assert = (condition, message) => { if (!condition) throw new Error(message); };

const resources = handlers.get('resources_discover')({ type: 'resources_discover', cwd: process.cwd(), reason: 'test' });
assert(resources.skillPaths.length === 1 && resources.skillPaths[0].endsWith('/skills'), 'resources_discover returns the package skills directory');

// --- session_start path (AC-3: unchanged) ---
handlers.get('session_start')({ type: 'session_start' });
let first = await handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'hello' }] }] });
assert(countBootstraps(first.messages) === 1, 'session_start injects exactly one structural bootstrap');
assert(textOf(first.messages[0]).includes('Load the $spacedock:first-officer skill'), 'bootstrap points at the shipped first-officer contract');
assert(textOf(first.messages[0]).includes('Pi tool mapping: read/write/edit/bash/grep/find/ls'), 'bootstrap includes the Pi tool mapping');

// --- session_compact path (AC-1 value-measuring, AC-2 mechanism) ---
handlers.get('agent_end')({ type: 'agent_end' }); // clear flags between paths
execCalls.length = 0; // reset exec call tracking
handlers.get('session_compact')({ type: 'session_compact' });
const compactSummary = { role: 'assistant', content: [{ type: 'text', text: 'Compaction summary mentions SPACEDOCK-FO-BOOTSTRAP-v1 but is not the bootstrap message.' }] };
let afterCompact = await handlers.get('context')({ messages: [compactSummary, { role: 'user', content: [{ type: 'text', text: 'continue' }] }] });
// AC-1: boot record present, contract re-injection absent
assert(countBootstraps(afterCompact.messages) === 0, 'compaction does NOT inject structural bootstrap text');
assert(countBootRecords(afterCompact.messages) === 1, 'compaction injects exactly one boot record');
assert(afterCompact.messages[0] === compactSummary, 'boot record is inserted after the leading compaction summary');
const bootRecordText = textOf(afterCompact.messages[1]);
assert(bootRecordText.includes('"command":"boot"'), 'compaction injects the boot record (command:boot)');
assert(!bootRecordText.includes(marker), 'compaction does NOT inject the FO bootstrap marker');
assert(!bootRecordText.includes('Load the $spacedock:first-officer skill'), 'compaction does NOT inject the contract pointer');
// AC-2: pi.exec called with the right args
assert(execCalls.length === 1, 'pi.exec called exactly once for the boot read');
assert(execCalls[0].command === 'spacedock', 'pi.exec called with spacedock');
assert(JSON.stringify(execCalls[0].args) === JSON.stringify(['status','--boot','--identify','--json']), 'pi.exec called with status --boot --identify --json');

// --- dedup (AC-5) ---
handlers.get('session_compact')({ type: 'session_compact' });
let deduped = await handlers.get('context')({ messages: afterCompact.messages });
assert(deduped === undefined, 'existing boot record de-duplicates context injection');

// --- agent_end suppresses further injection ---
handlers.get('agent_end')({ type: 'agent_end' });
let suppressed = await handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'new turn' }] }] });
assert(suppressed === undefined, 'agent_end suppresses further injection');
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
const pi = {
  on(event, handler) { handlers.set(event, handler); },
  exec() { return Promise.resolve({ stdout: '', stderr: '', code: 0, killed: false }); }
};
register(pi);

const assert = (condition, message) => { if (!condition) throw new Error(message); };
assert(process.env.PI_SUBAGENT_CHILD === '1', 'harness must run under the subagent-child marker');

const resources = handlers.get('resources_discover')({ type: 'resources_discover', cwd: process.cwd(), reason: 'test' });
assert(resources.skillPaths.length === 1 && resources.skillPaths[0].endsWith('/skills'), 'child sessions still discover the package skills directory');

handlers.get('session_start')({ type: 'session_start' });
const afterStart = await handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'hello' }] }] });
assert(afterStart === undefined, 'PI_SUBAGENT_CHILD=1 session_start injects zero FO bootstrap');

handlers.get('session_compact')({ type: 'session_compact' });
const afterCompact = await handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'continue' }] }] });
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
