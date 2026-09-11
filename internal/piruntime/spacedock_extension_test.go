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
	modulePath := filepath.Join(tmp, ".pi", "extensions", "spacedock-extension.mjs")
	if err := os.MkdirAll(filepath.Dir(modulePath), 0o755); err != nil {
		t.Fatalf("stage extension dir: %v", err)
	}
	if err := os.WriteFile(modulePath, extensionSource, 0o644); err != nil {
		t.Fatalf("write harness module: %v", err)
	}

	harnessPath := filepath.Join(tmp, "harness.mjs")
	harness := `
import register from './.pi/extensions/spacedock-extension.mjs';
import * as fs from 'node:fs';

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
// The real pi contract runs this extension under a frontdoor launch, where
// the spacedock launcher sets PI_SPACEDOCK_LAUNCH=1 in the child env. Pin it
// here so an inherited (or absent) marker on the host shell cannot flip the
// gate.
assert(process.env.PI_SPACEDOCK_LAUNCH === '1', 'harness must run under the frontdoor launch marker');

// The real pi contract loads this extension from the registered package root,
// where package.json declares pi.skills: ["./skills"]. The module is staged
// bare in tmp, so first exercise the undeclared arm, then stage the manifest
// and assert the manifest-skip (AC-4: no skillPaths for manifest-declared dirs).
const manifestPath = new URL('./package.json', import.meta.url).pathname;
let resources = handlers.get('resources_discover')({ type: 'resources_discover', cwd: process.cwd(), reason: 'test' });
assert(resources.skillPaths.length === 1 && resources.skillPaths[0].endsWith('/skills'), 'resources_discover returns the package skills directory when the manifest does not declare it');
fs.writeFileSync(manifestPath, JSON.stringify({ pi: { skills: ['./skills'] } }));
resources = handlers.get('resources_discover')({ type: 'resources_discover', cwd: process.cwd(), reason: 'test' });
assert(resources.skillPaths.length === 0, 'resources_discover skips the skills dir the package manifest already declares');

// --- session_start path (AC-3: unchanged) ---
handlers.get('session_start')({ type: 'session_start' });
let first = await handlers.get('context')({ messages: [{ role: 'user', content: [{ type: 'text', text: 'hello' }] }] });
assert(countBootstraps(first.messages) === 1, 'session_start injects exactly one structural bootstrap');
assert(textOf(first.messages[0]).includes('[SPACEDOCK-FO-BOOTSTRAP-v1]'), 'bootstrap carries the FO bootstrap marker');
assert(textOf(first.messages[0]).includes('<available_skills>'), 'bootstrap points at the resolvable <available_skills> trigger');
assert(!textOf(first.messages[0]).includes('$spacedock:'), 'bootstrap uses no unexpandable $spacedock: syntax');
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
assert(!bootRecordText.includes('<available_skills>'), 'compaction does NOT inject the contract pointer');
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
	// The unset subagent marker must be genuinely unset: a pi-subagents child
	// shell exports PI_SUBAGENT_CHILD, and an inherited marker would silently
	// flip the exemption inside the harness process. The frontdoor launch
	// marker is pinned ON: the real pi contract runs this extension under a
	// frontdoor launch, and the injection gate reads it at handler time.
	cmd.Env = harnessEnv(map[string]string{"PI_SPACEDOCK_LAUNCH": "1"}, "PI_SUBAGENT_CHILD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("extension behavior harness failed: %v\n%s", err, out)
	}
}

// TestSpacedockPiExtensionChildExemption is the AC-1 child-session half: a
// pi-subagents child (PI_SUBAGENT_CHILD=1) is a delegated worker, not a first
// officer, so the context hook must inject zero FO bootstrap across
// session_start and session_compact. Skill discovery is unaffected. The run
// pins PI_SPACEDOCK_LAUNCH=1 as well, proving the child exemption dominates
// the frontdoor marker and that neither marker state is inherited by accident.
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
import register from './.pi/extensions/spacedock-extension.mjs';

const handlers = new Map();
const pi = {
  on(event, handler) { handlers.set(event, handler); },
  exec() { return Promise.resolve({ stdout: '', stderr: '', code: 0, killed: false }); }
};
register(pi);

const assert = (condition, message) => { if (!condition) throw new Error(message); };
assert(process.env.PI_SUBAGENT_CHILD === '1', 'harness must run under the subagent-child marker');
assert(process.env.PI_SPACEDOCK_LAUNCH === '1', 'child harness must pin the frontdoor launch marker');

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
	cmd.Env = harnessEnv(map[string]string{"PI_SUBAGENT_CHILD": "1", "PI_SPACEDOCK_LAUNCH": "1"})
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("child-exemption harness failed: %v\n%s", err, out)
	}
}

// copyExtension stages the shipped .pi/extensions/spacedock.ts as a node
// module under tmp at the real package layout (<tmp>/.pi/extensions/), so the
// extension's repo-root resolution (two levels up from the module) matches the
// registered-package shape pi actually loads.
func copyExtension(t *testing.T, repoRoot, tmp string) {
	t.Helper()
	extensionSource, err := os.ReadFile(filepath.Join(repoRoot, ".pi", "extensions", "spacedock.ts"))
	if err != nil {
		t.Fatalf("read extension: %v", err)
	}
	modulePath := filepath.Join(tmp, ".pi", "extensions", "spacedock-extension.mjs")
	if err := os.MkdirAll(filepath.Dir(modulePath), 0o755); err != nil {
		t.Fatalf("stage extension dir: %v", err)
	}
	if err := os.WriteFile(modulePath, extensionSource, 0o644); err != nil {
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
