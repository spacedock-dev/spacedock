// ABOUTME: Structural lint — a FO reference file must not teach a bare `spacedock`
// ABOUTME: helper INVOCATION by example; post-gate helper calls use ${SPACEDOCK_BIN:-spacedock}.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The FO launcher invariant pins ONE launcher at the version gate and uses it for
// every later Spacedock helper call. A bare `spacedock <helper> …` INVOCATION in a
// FO reference doc teaches the drift the invariant forbids: a reader copies it and
// runs a different `$PATH` binary mid-session. This is a doc-AUTHORING rule — a
// defect a machine can see (a bare launcher token in a runnable example), not a
// prose-grep of a behavior claim and not a code-bound consistency check. It is the
// AC-2 structural arm; the behavior is proven by the SPACEDOCK_BIN-vs-PATH live
// drive in internal/ensigncycle (AC-3).
//
// The rule discriminates a runnable INVOCATION from the legitimate bare forms:
//   - it names a helper verb (status/state/dispatch/merge/new) AND carries an
//     invocation flag (`--workflow-dir`, `--discover`, `--set`, …), so a bare
//     command NAME mentioned in prose ("`spacedock new <slug>` mints the id") is
//     not flagged — naming a command is not invoking it;
//   - it is NOT a `→` capability-binding line, which names the SHIPPED command
//     surface, not a call the FO emits;
//   - it does NOT already resolve the launcher (`${SPACEDOCK_BIN:-spacedock}`);
//   - it is NOT in a fallback/diagnostic/install context (`on $PATH`, `doctor`,
//     `brew install`, `go build`, `--help`), where bare `spacedock` is correct
//     (the version-gate fallback probe and operator hints).

// foReferenceDir is the first-officer reference surface this lint walks.
func foReferenceDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(skillsRoot(t), "first-officer", "references")
}

// launcherHelperInvocation matches a bare `spacedock <helper> … --<flag>` runnable
// example inside a backtick span: a helper verb followed (anywhere before the next
// backtick) by an invocation flag. The leading boundary keeps `${SPACEDOCK_BIN:-spacedock}`
// from matching — that form has `:-` before `spacedock`, not a span/space boundary.
var launcherHelperInvocation = regexp.MustCompile(
	"`spacedock (?:status|state|dispatch|merge|new)\\b[^`]*?--(?:workflow-dir|boot|discover|set|next-id|validate|json|resolve|where|next|archived)\\b")

// launcherDiagnosticContext marks a line where a bare `spacedock` is legitimate: the
// version-gate PATH fallback probe, the `doctor` remedy, an install hint, or a
// `--help` reference.
var launcherDiagnosticContext = regexp.MustCompile("on `?\\$?PATH|\\bdoctor\\b|brew install|go build|--help")

// lineHasBareLauncherHelperCall reports whether a doc line teaches a bare
// `spacedock` helper INVOCATION the launcher invariant forbids.
func lineHasBareLauncherHelperCall(line string) bool {
	if !launcherHelperInvocation.MatchString(line) {
		return false
	}
	if strings.Contains(line, "${SPACEDOCK_BIN:-spacedock}") {
		return false
	}
	if strings.HasPrefix(strings.TrimSpace(line), "- → ") {
		return false
	}
	if launcherDiagnosticContext.MatchString(line) {
		return false
	}
	return true
}

// TestFOReferencesUseResolvedLauncher is the AC-2 lint: no FO reference file teaches
// a bare `spacedock` helper invocation by example. A flagged line must resolve the
// launcher (`${SPACEDOCK_BIN:-spacedock}`) instead.
func TestFOReferencesUseResolvedLauncher(t *testing.T) {
	dir := foReferenceDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read FO reference dir %s: %v", dir, err)
	}
	scanned := 0
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, ent.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		scanned++
		for i, line := range strings.Split(string(data), "\n") {
			if lineHasBareLauncherHelperCall(line) {
				t.Errorf("%s:%d teaches a bare `spacedock` helper invocation; post-gate helper calls must resolve the pinned launcher as `${SPACEDOCK_BIN:-spacedock}` (launcher invariant): %q", path, i+1, strings.TrimSpace(line))
			}
		}
	}
	if scanned == 0 {
		t.Fatal("scanned zero FO reference files — extractor bug; the lint would pass vacuously")
	}
}

// TestBareLauncherHelperScannerDiscriminates is the DISCRIMINATOR control: it proves
// the scanner flags a genuine bare-invocation and PASSES every legitimate form, so
// TestFOReferencesUseResolvedLauncher can never pass vacuously (e.g. a typo'd verb
// that never matches, or an exemption swallowing the real leak).
func TestBareLauncherHelperScannerDiscriminates(t *testing.T) {
	// The real leak shape — a bare runnable helper invocation. MUST flag.
	leak := "otherwise `spacedock status --discover`: one path → use it"
	if !lineHasBareLauncherHelperCall(leak) {
		t.Errorf("discriminator: a bare `spacedock status --discover` invocation was NOT flagged (scanner would pass vacuously): %q", leak)
	}

	// A bare state-mutation invocation. MUST flag.
	leakSet := "Entity frontmatter — via `spacedock status --set` for all field updates"
	if !lineHasBareLauncherHelperCall(leakSet) {
		t.Errorf("discriminator: a bare `spacedock status --set` invocation was NOT flagged: %q", leakSet)
	}

	// The resolved-launcher form. MUST pass — this is the contract-blessed invocation.
	resolved := "run `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --discover`"
	if lineHasBareLauncherHelperCall(resolved) {
		t.Errorf("discriminator: the resolved `${SPACEDOCK_BIN:-spacedock}` form was wrongly flagged: %q", resolved)
	}

	// A bare command NAME with no invocation flag — naming, not invoking. MUST pass.
	nameOnly := "`spacedock new <slug>` mints the id and writes the stamped entity"
	if lineHasBareLauncherHelperCall(nameOnly) {
		t.Errorf("discriminator: a bare command-name mention (no invocation flag) was wrongly flagged: %q", nameOnly)
	}

	// A `→` capability-binding line names the SHIPPED command surface. MUST pass.
	shippedLine := "- → **shipped**: `` `spacedock status --boot --json` ``."
	if lineHasBareLauncherHelperCall(shippedLine) {
		t.Errorf("discriminator: a `→ shipped:` capability-binding line was wrongly flagged: %q", shippedLine)
	}

	// The version-gate PATH fallback probe — bare `spacedock` is correct here. MUST pass.
	fallback := "If `SPACEDOCK_BIN` is unusable, retry once with bare `spacedock status --discover` on `$PATH`"
	if lineHasBareLauncherHelperCall(fallback) {
		t.Errorf("discriminator: the version-gate PATH fallback probe was wrongly flagged: %q", fallback)
	}
}
