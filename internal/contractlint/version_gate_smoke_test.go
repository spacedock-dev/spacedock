// ABOUTME: Skill smoke test for the rewritten FO version gate (Startup step 1):
// ABOUTME: the lean in-core skeleton (OS-aware hint, sandbox check, deferred
// ABOUTME: trigger, launcher-invariant amendment) plus the deferred
// ABOUTME: references/fo-install-gate.md machinery (sentinel loop bound,
// ABOUTME: convergence, message grammar) — prose can drift, so tokens are
// ABOUTME: pinned against the on-disk skill files.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var sharedCorePath = filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md")
var installGatePath = filepath.Join("skills", "first-officer", "references", "fo-install-gate.md")

func readSkillFile(t *testing.T, rel string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(repoRoot(t), rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(body)
}

// TestVersionGateProseOSAwareHint pins the AC-1 prose shape in the lean core
// skeleton: the `uname -s` OS source (permanently load-bearing there — the
// absent binary produces no `--version` output) and the unsupported-OS
// source-build escape. The install commands themselves are pinned by
// TestInstallHintNoDrift as a relation against docs/site/get-started/install.md,
// so they are deliberately not frozen here as literals.
func TestVersionGateProseOSAwareHint(t *testing.T) {
	body := readSkillFile(t, sharedCorePath)
	for _, token := range []string{
		"uname -s",
		"go build -o spacedock ./cmd/spacedock",
		"unsupported OS",
	} {
		if !strings.Contains(body, token) {
			t.Fatalf("FO version-gate core skeleton is missing the OS-aware-hint token %q", token)
		}
	}
}

// TestVersionGateDeferredTrigger pins the core skeleton's deferred-read trigger:
// the binary-absent class names references/fo-install-gate.md INLINE, and the
// deferred-load-points inventory carries the matching one-line entry (mirroring
// the existing fo-dispatch-core.md pattern). The heavyweight machinery must NOT
// live in the boot-resident core — the byte-cap tests guard that.
func TestVersionGateDeferredTrigger(t *testing.T) {
	body := readSkillFile(t, sharedCorePath)
	if !strings.Contains(body, "read `references/fo-install-gate.md`") {
		t.Fatalf("core skeleton must name the deferred read of references/fo-install-gate.md in the binary-absent class")
	}
	if !strings.Contains(body, "- `references/fo-install-gate.md`") {
		t.Fatalf("deferred-load-points inventory must carry a references/fo-install-gate.md entry line")
	}
	// The deferred body exists and is non-trivial.
	gate := readSkillFile(t, installGatePath)
	if len(gate) < 500 {
		t.Fatalf("references/fo-install-gate.md = %d bytes, want the real install-offer machinery", len(gate))
	}
}

// TestInstallGateSentinelLoopBound pins the AC-2 guardrail shape in the deferred
// body: the sentinel path form, the `test -f` pre-offer check,
// create-before-run, the named identity env vars (claude + codex) with the pi
// working-directory-hash fallback, and the fallback message's sentinel-path +
// `rm` recovery obligation. Dropping any of these breaks the one-attempt bound
// the failure-fallback behavior fixture asserts on.
func TestInstallGateSentinelLoopBound(t *testing.T) {
	body := readSkillFile(t, installGatePath)
	for _, token := range []string{
		"${TMPDIR:-/tmp}/spacedock-install-attempted-<key>",
		"test -f <sentinel>",
		"BEFORE the install runs",
		"create-before-run",
		"CLAUDE_CODE_SESSION_ID",
		"CODEX_THREAD_ID",
		"rm <sentinel-path>",
		"no second install attempt",
	} {
		if !strings.Contains(body, token) {
			t.Fatalf("fo-install-gate.md is missing the sentinel-loop-bound token %q", token)
		}
	}
	// The pi fallback must be stated (pi exposes no session-identity env var).
	if !strings.Contains(body, "working directory") || !strings.Contains(body, "project-scoped") {
		t.Fatalf("fo-install-gate.md must state the pi project-scoped cwd-hash sentinel fallback")
	}
	// Convergence mechanics: session-scoped repoint, stderr-contract parse,
	// $HOME/.local/bin probe — and never persisted (D-2).
	for _, token := range []string{
		"install.sh: installed spacedock <version> to <dir>/spacedock",
		"$HOME/.local/bin/spacedock",
		"never persist it to a shell profile",
	} {
		if !strings.Contains(body, token) {
			t.Fatalf("fo-install-gate.md is missing the convergence token %q", token)
		}
	}
}

// insideRegistryRowRe extracts one full registry row from the safehouse source:
// the env var name AND its wantValue. A names-only assertion would stay green
// through a wantValue change while the prose check silently never fires.
var insideRegistryRowRe = regexp.MustCompile(`\{env: "([A-Z0-9_]+)", wantValue: "([^"]+)"`)

// TestVersionGateSandboxRegistry asserts BOTH the core check sentence and the
// deferred message body match EVERY row of the binary's insideRegistry — full
// name+value rows, read by source-grep of internal/safehouse/state.go (no new
// exported API), so a registry change the prose does not mirror fails here.
func TestVersionGateSandboxRegistry(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(repoRoot(t), "internal", "safehouse", "state.go"))
	if err != nil {
		t.Fatalf("read safehouse registry source: %v", err)
	}
	rows := insideRegistryRowRe.FindAllStringSubmatch(string(src), -1)
	if len(rows) == 0 {
		t.Fatalf("source-grep of internal/safehouse/state.go found no insideRegistry rows — the extraction regex must track the table's literal shape")
	}
	core := readSkillFile(t, sharedCorePath)
	gate := readSkillFile(t, installGatePath)
	for _, row := range rows {
		name, wantValue := row[1], row[2]
		for label, body := range map[string]string{"core": core, "fo-install-gate": gate} {
			if !strings.Contains(body, name) {
				t.Fatalf("%s prose does not check the sandbox env var %q", label, name)
			}
			if !strings.Contains(body, wantValue) {
				t.Fatalf("%s prose does not pin the registry VALUE %q for %q (matching is on value, not presence)", label, wantValue, name)
			}
		}
	}
	if !strings.Contains(core, "outside the sandbox") {
		t.Fatalf("core skeleton is missing the sandbox outcome sentence (run outside the sandbox)")
	}
	if !strings.Contains(gate, "outside the sandbox") {
		t.Fatalf("fo-install-gate.md is missing the human-run-outside-the-sandbox message body")
	}
	if !strings.Contains(gate, "^Sandbox: ") {
		t.Fatalf("fo-install-gate.md must name the ^Sandbox: line as corroboration-only")
	}
}

// TestVersionGateProseLauncherInvariantAmendment pins the amended invariant
// sentence in the core: a gate that fails the binary-absent class and then
// succeeds via the approved install performs its ONE launcher resolution at
// that point (the session-scoped SPACEDOCK_BIN repoint), still inside Startup
// step 1. Without this amendment the post-install repoint contradicts the
// invariant it ships alongside.
func TestVersionGateProseLauncherInvariantAmendment(t *testing.T) {
	body := readSkillFile(t, sharedCorePath)
	if !strings.Contains(body, "performs its ONE launcher resolution at that point") {
		t.Fatalf("launcher-invariant amendment missing: gate-internal post-install resolution not blessed")
	}
	if !strings.Contains(body, "Never drift to a bare `spacedock` mid-session") {
		t.Fatalf("launcher-invariant anti-drift sentence must survive the amendment")
	}
}
