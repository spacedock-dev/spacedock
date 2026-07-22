// ABOUTME: Structural lint — a contract reference or deferred-skill file must not teach
// ABOUTME: a bare `spacedock` helper INVOCATION by example; resolved calls use ${SPACEDOCK_BIN:-spacedock}.
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
// contract reference or deferred-skill file teaches the drift the invariant forbids: a
// reader copies it and runs a different `$PATH` binary mid-session. This is a doc-
// AUTHORING rule — a defect a machine can see (a bare launcher token in a runnable
// example), not a prose-grep of a behavior claim and not a code-bound consistency
// check. It is the AC-2 structural arm; the behavior is proven by the SPACEDOCK_BIN-
// vs-PATH live drive in internal/ensigncycle (AC-3).
//
// The rule discriminates a runnable INVOCATION from the legitimate bare forms, applied
// per backtick SPAN (not per line, so one exemption span cannot swallow an invocation
// span sharing the same line):
//   - the span names a helper verb (status/state/dispatch/merge/new/contract) AND
//     carries ANY long flag (`--[a-z][a-z-]*`), so a bare command NAME mentioned in
//     prose ("`spacedock new <slug>` mints the id") is not flagged — naming a command
//     is not invoking it. The generic flag matcher means a newly introduced flag
//     (`--advance`, `--folder`, …) cannot silently escape the lint the way `--name`
//     once did;
//   - the LINE is NOT a `→` capability-binding line, which names the SHIPPED command
//     surface, not a call the FO emits;
//   - the SPAN does NOT already resolve the launcher (`${SPACEDOCK_BIN:-spacedock}`);
//   - the SPAN is NOT a `--help` reference — evaluated per span so a legitimate
//     `spacedock new --help` mention cannot exempt an invocation span on the same line;
//   - the LINE is NOT in a fallback/diagnostic/install context (`on $PATH`, `doctor`,
//     `brew install`, `go build`), where bare `spacedock` is correct (the version-gate
//     fallback probe and operator hints).
//
// Independently, a fenced code-block command-position line (`^\s*spacedock\s`) is
// always an invocation — no flag or verb requirement, since command position inside a
// block is executable by definition, not a mention.

// foReferenceDir is the first-officer reference surface this lint walks.
func foReferenceDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(skillsRoot(t), "first-officer", "references")
}

// ensignReferenceDir is the ensign reference surface this lint also walks.
func ensignReferenceDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(skillsRoot(t), "ensign", "references")
}

// deferredSkillPaths is the fixed set of deferred-skill SKILL.md files this lint also
// walks. These load on-trigger (feedback rejection, status query, gate
// presentation, dispatch-failure recovery) rather than at boot, but still teach the
// FO launcher invocations by example.
func deferredSkillPaths(t *testing.T) []string {
	t.Helper()
	root := skillsRoot(t)
	names := []string{
		"fo-status-viewer", "present-gate",
		"feedback-rejection-flow", "fo-dispatch-recovery",
	}
	paths := make([]string, len(names))
	for i, name := range names {
		paths[i] = filepath.Join(root, name, "SKILL.md")
	}
	return paths
}

// launcherSurfaceGroup is one independently-scanned walk of the launcher-invariant
// lint; each group carries its own scanned>0 guard so a moved or renamed directory
// cannot make the lint pass vacuously.
type launcherSurfaceGroup struct {
	name  string
	paths []string
}

func mdFilesIn(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir %s: %v", dir, err)
	}
	var paths []string
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".md") {
			continue
		}
		paths = append(paths, filepath.Join(dir, ent.Name()))
	}
	return paths
}

func launcherSurfaceGroups(t *testing.T) []launcherSurfaceGroup {
	t.Helper()
	return []launcherSurfaceGroup{
		{"first-officer references", mdFilesIn(t, foReferenceDir(t))},
		{"ensign references", mdFilesIn(t, ensignReferenceDir(t))},
		{"deferred-skill SKILL.md files", deferredSkillPaths(t)},
	}
}

// launcherSpanInvocation matches a bare `spacedock <verb> …` runnable example within a
// SINGLE backtick span: a helper verb followed (anywhere later in the span) by ANY long
// flag. The leading anchor keeps `${SPACEDOCK_BIN:-spacedock}` from matching — that
// form starts with `${`, not `spacedock`.
var launcherSpanInvocation = regexp.MustCompile(
	`^spacedock (?:status|state|dispatch|merge|new|contract)\b.*--[a-z][a-z-]*`)

// backtickSpan extracts single-backtick-delimited code spans from one line.
var backtickSpan = regexp.MustCompile("`([^`]+)`")

// launcherDiagnosticLineContext marks a LINE where a bare `spacedock` is legitimate
// regardless of span content: the version-gate PATH fallback probe, the `doctor`
// remedy, or an install hint. `--help` is deliberately NOT here — it is span-scoped
// (see lineHasBareLauncherHelperCall) so it cannot swallow a genuine invocation sharing
// the same line.
var launcherDiagnosticLineContext = regexp.MustCompile("on `?\\$?PATH|\\bdoctor\\b|brew install|go build")

// lineHasBareLauncherHelperCall reports whether a doc line teaches a bare `spacedock`
// helper INVOCATION the launcher invariant forbids, scanning each backtick span on the
// line independently.
func lineHasBareLauncherHelperCall(line string) bool {
	if strings.HasPrefix(strings.TrimSpace(line), "- → ") {
		return false
	}
	if launcherDiagnosticLineContext.MatchString(line) {
		return false
	}
	for _, m := range backtickSpan.FindAllStringSubmatch(line, -1) {
		span := m[1]
		if strings.Contains(span, "${SPACEDOCK_BIN:-spacedock}") {
			continue
		}
		if strings.Contains(span, "--help") {
			continue
		}
		if launcherSpanInvocation.MatchString(span) {
			return true
		}
	}
	return false
}

// fenceMarker matches a fenced-code-block delimiter line (```` ``` ```` optionally
// followed by a language tag).
var fenceMarker = regexp.MustCompile("^\\s*```")

// fencedLauncherLine matches a bare `spacedock` command-position line inside a fenced
// code block: no flag or verb requirement, since command position in a block is
// executable by definition.
var fencedLauncherLine = regexp.MustCompile(`^\s*spacedock\s`)

// scanBareLauncherCalls returns the 0-based line indices in content that teach a bare
// `spacedock` launcher invocation: a backtick-span invocation outside any fence, or a
// fenced-code-block command-position line.
func scanBareLauncherCalls(content string) []int {
	var violations []int
	inFence := false
	for i, line := range strings.Split(content, "\n") {
		if fenceMarker.MatchString(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			if fencedLauncherLine.MatchString(line) {
				violations = append(violations, i)
			}
			continue
		}
		if lineHasBareLauncherHelperCall(line) {
			violations = append(violations, i)
		}
	}
	return violations
}

// TestLauncherSurfaceUsesResolvedLauncher is the AC-1 lint: no file across the three
// walked groups teaches a bare `spacedock` helper invocation by example. A flagged line
// must resolve the launcher (`${SPACEDOCK_BIN:-spacedock}`) instead.
func TestLauncherSurfaceUsesResolvedLauncher(t *testing.T) {
	for _, group := range launcherSurfaceGroups(t) {
		scanned := 0
		for _, path := range group.paths {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			scanned++
			lines := strings.Split(string(data), "\n")
			for _, idx := range scanBareLauncherCalls(string(data)) {
				t.Errorf("%s:%d teaches a bare `spacedock` launcher invocation; post-gate helper calls must resolve the pinned launcher as `${SPACEDOCK_BIN:-spacedock}` (launcher invariant): %q", path, idx+1, strings.TrimSpace(lines[idx]))
			}
		}
		if scanned == 0 {
			t.Fatalf("scanned zero files in group %q — extractor bug; the lint would pass vacuously", group.name)
		}
	}
}

// TestBareLauncherHelperScannerDiscriminates is the DISCRIMINATOR control for the
// backtick-span scanner: it proves the scanner flags a genuine bare-invocation span and
// PASSES every legitimate form, so TestLauncherSurfaceUsesResolvedLauncher can never
// pass vacuously (e.g. a typo'd verb that never matches, or an exemption swallowing the
// real leak).
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

	// A bare `--name` helper invocation — the imperative context-budget probe. MUST
	// flag: this is the class the lint missed at claude-fo-dispatch.md:132 because
	// `--name` was absent from the invocation-flag allowlist. The generic flag matcher
	// closes this class permanently, not just for `--name`.
	leakName := "**Context budget check:** Run `spacedock dispatch context-budget --name {ensign-name}`. Parse the JSON output."
	if !lineHasBareLauncherHelperCall(leakName) {
		t.Errorf("discriminator: a bare `spacedock dispatch context-budget --name` invocation was NOT flagged (the line-132 escape class): %q", leakName)
	}

	// An unenumerated `--advance` invocation — the exact flag the old allowlist
	// missed at claude-fo-dispatch.md:38/44. MUST flag under the generic matcher.
	leakAdvance := "run `spacedock dispatch build --advance`"
	if !lineHasBareLauncherHelperCall(leakAdvance) {
		t.Errorf("discriminator: a bare `spacedock dispatch build --advance` invocation was NOT flagged (the flag-allowlist escape class): %q", leakAdvance)
	}

	// The `→`-binding TWIN of the same `dispatch context-budget --name` text at
	// fo-dispatch-core.md:98 — it names the SHIPPED command surface per host, not an
	// FO-emitted call. MUST pass, so adding `--name` does not over-flag the binding line.
	bindingTwin := "- → **Claude:** PRESENT — `spacedock dispatch context-budget --name {name}`. · **Codex:** ABSENT. · **Pi:** ABSENT."
	if lineHasBareLauncherHelperCall(bindingTwin) {
		t.Errorf("discriminator: the `→`-binding twin of the `--name` invocation was wrongly flagged: %q", bindingTwin)
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

	// A same-line mix: a must-fix invocation span AND a legitimate `--help` span, the
	// claude-first-officer-runtime.md:35 shape. MUST flag — a line-scoped `--help`
	// exemption would swallow the invocation; the span-scoped exemption must not.
	mixedLine := "Use `spacedock new <slug> [--folder] [--id-seed S --id-actor A]` via Bash (see `spacedock new --help`)"
	if !lineHasBareLauncherHelperCall(mixedLine) {
		t.Errorf("discriminator: an invocation span was swallowed by a `--help` span on the same line (the claude-first-officer-runtime.md:35 escape): %q", mixedLine)
	}

	// The `--help` span alone, with no invocation elsewhere on the line. MUST pass.
	helpOnly := "see `spacedock new --help` for options"
	if lineHasBareLauncherHelperCall(helpOnly) {
		t.Errorf("discriminator: a lone `--help` span was wrongly flagged: %q", helpOnly)
	}
}

// TestFencedLauncherBlockScannerDiscriminates is the DISCRIMINATOR control for the
// fenced-code-block rule: a bare `spacedock` command-position line inside a fence MUST
// flag regardless of flags, and the resolved-launcher form MUST pass.
func TestFencedLauncherBlockScannerDiscriminates(t *testing.T) {
	leak := "Update main-branch frontmatter for dispatch:\n```\nspacedock status --workflow-dir {workflow_dir} --set {slug} status={next_stage}\n```\n"
	if got := scanBareLauncherCalls(leak); len(got) != 1 {
		t.Errorf("discriminator: a bare `spacedock` command-position line inside a fenced code block was not flagged exactly once: got violations at line indices %v", got)
	}

	resolved := "Update main-branch frontmatter for dispatch:\n```\n${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={next_stage}\n```\n"
	if got := scanBareLauncherCalls(resolved); len(got) != 0 {
		t.Errorf("discriminator: a resolved-launcher line inside a fenced code block was wrongly flagged: violations at line indices %v", got)
	}

	// A bare `spacedock` mention OUTSIDE any fence, in ordinary prose (no backtick
	// span at all), is not a runnable example and MUST pass.
	proseOutsideFence := "spacedock is the CLI this workflow drives.\n"
	if got := scanBareLauncherCalls(proseOutsideFence); len(got) != 0 {
		t.Errorf("discriminator: unfenced bare prose was wrongly flagged: violations at line indices %v", got)
	}
}

// TestLauncherSurfaceGroupsScopeDiscriminates is the DISCRIMINATOR control for the
// walked-group list itself: `launcherSurfaceGroups` MUST return exactly the three named
// groups, each with a non-empty path set. Pinning `deferredSkillPaths` alone (below)
// does not catch a group deregistered from `launcherSurfaceGroups` — a per-group
// scanned>0 guard cannot fire for a group that was never returned, so deleting a group
// entry there is invisible to every other test in this file. This test calls
// `launcherSurfaceGroups` directly to close that gap.
func TestLauncherSurfaceGroupsScopeDiscriminates(t *testing.T) {
	wantNames := []string{"first-officer references", "ensign references", "deferred-skill SKILL.md files"}
	groups := launcherSurfaceGroups(t)
	if len(groups) != len(wantNames) {
		t.Fatalf("launcherSurfaceGroups returned %d groups, want %d (a deregistered group would silently stop being scanned): %v", len(groups), len(wantNames), groups)
	}
	for i, wantName := range wantNames {
		if groups[i].name != wantName {
			t.Errorf("launcherSurfaceGroups[%d].name = %q, want %q", i, groups[i].name, wantName)
		}
		if len(groups[i].paths) == 0 {
			t.Errorf("launcherSurfaceGroups[%d] (%q) has zero paths — a moved or renamed directory would make the lint pass vacuously", i, groups[i].name)
		}
	}
}

// TestDeferredSkillLauncherScopeDiscriminates is the DISCRIMINATOR control for the
// scope-extension escape class: the deferred-skill SKILL.md files that were
// previously unscanned MUST be present in the walked set, by exact name.
func TestDeferredSkillLauncherScopeDiscriminates(t *testing.T) {
	want := []string{
		"fo-status-viewer", "present-gate",
		"feedback-rejection-flow", "fo-dispatch-recovery",
	}
	paths := deferredSkillPaths(t)
	if len(paths) != len(want) {
		t.Fatalf("deferredSkillPaths returned %d paths, want %d: %v", len(paths), len(want), paths)
	}
	for i, name := range want {
		wantSuffix := filepath.Join(name, "SKILL.md")
		if !strings.HasSuffix(paths[i], wantSuffix) {
			t.Errorf("deferredSkillPaths[%d] = %q, want suffix %q (the unscanned-scope escape class this lint closes)", i, paths[i], wantSuffix)
		}
		if _, err := os.Stat(paths[i]); err != nil {
			t.Errorf("deferred-skill path %s does not exist: %v", paths[i], err)
		}
	}
}
