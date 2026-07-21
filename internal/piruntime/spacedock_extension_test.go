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
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("extension behavior harness failed: %v\n%s", err, out)
	}
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
