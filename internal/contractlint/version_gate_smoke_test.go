// ABOUTME: Structural checks for the FO version gate (Startup step 1): the
// ABOUTME: deferred references/fo-install.md machinery resolves and stays
// ABOUTME: out of the boot-resident core, and the sandbox registry and install
// ABOUTME: hint match their sources outside the skill files.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var sharedCorePath = filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md")
var installRefPath = filepath.Join("skills", "first-officer", "references", "fo-install.md")

func readSkillFile(t *testing.T, rel string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(repoRoot(t), rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(body)
}

// TestVersionGateDeferredTrigger checks the structural shape of the core
// skeleton's deferred-read trigger: the deferred-load-points inventory carries
// the references/fo-install.md entry (mirroring the existing
// fo-dispatch-core.md pattern) and the deferred body resolves non-trivially.
// The heavyweight machinery must NOT live in the boot-resident core — the
// byte-cap tests guard that.
func TestVersionGateDeferredTrigger(t *testing.T) {
	body := readSkillFile(t, sharedCorePath)
	if !strings.Contains(body, "- `references/fo-install.md`") {
		t.Fatalf("deferred-load-points inventory must carry a references/fo-install.md entry line")
	}
	if strings.Contains(body, "gate --help") {
		t.Fatalf("boot-resident core must not carry an inline `gate --help` capability probe: the binary gate relies on the minor-version match alone, and a stale launcher fails at the first `gate withdraw`")
	}
	// The deferred body exists and is non-trivial.
	gate := readSkillFile(t, installRefPath)
	if len(gate) < 500 {
		t.Fatalf("references/fo-install.md = %d bytes, want the real install-offer machinery", len(gate))
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
	// The registry mirror is asserted against the CORE ALONE, because the core is what
	// PERFORMS the sandbox check (Startup step 1, every gate class). fo-install.md only
	// consumes the verdict, so duplicating the env name+value there bought nothing but a
	// second place to forget to update.
	core := readSkillFile(t, sharedCorePath)
	for _, row := range rows {
		name, wantValue := row[1], row[2]
		if !strings.Contains(core, name) {
			t.Fatalf("core prose does not check the sandbox env var %q", name)
		}
		if !strings.Contains(core, wantValue) {
			t.Fatalf("core prose does not pin the registry VALUE %q for %q (matching is on value, not presence)", wantValue, name)
		}
	}
	if !strings.Contains(core, "outside the sandbox") {
		t.Fatalf("core skeleton is missing the sandbox outcome sentence (run outside the sandbox)")
	}
	// fo-install.md still owns the outcome the human sees: no install here, run it outside.
	if !strings.Contains(readSkillFile(t, installRefPath), "outside the sandbox") {
		t.Fatalf("fo-install.md is missing the human-run-outside-the-sandbox message body")
	}
}
