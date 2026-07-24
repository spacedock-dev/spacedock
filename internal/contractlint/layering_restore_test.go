// ABOUTME: Structural-absence guards restoring two FO-contract layering boundaries —
// ABOUTME: model-enum tokens out of the dispatch core, gh-pr-view out of the generic event loop.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// claudeModelTokens are the LITERAL Claude-model tokens that must NOT appear in the
// host-neutral dispatch core. `opus[1m]` is the captain-session fallback example and
// the bare enum words are the canonical model enum — both are Claude-specific (Codex
// and Pi load this core too and have no such enum). Their presence in the core IS the
// defect (structural-absence, same family as claudeTeamDispatchTokens), not a
// paraphrasable meaning: a host-neutral statement that defers to "the host's canonical
// model enum" carries none of these tokens and passes.
var claudeModelTokens = []string{
	"opus[1m]",
	"sonnet",
	"opus",
	"haiku",
}

// lineLeaksClaudeModelToken reports whether a line carries a literal Claude-model
// token. This is the scanner the host-neutral-core check and its discriminator
// control both drive.
func lineLeaksClaudeModelToken(line string) bool {
	for _, tok := range claudeModelTokens {
		if strings.Contains(line, tok) {
			return true
		}
	}
	return false
}

// TestDispatchCoreHasNoClaudeModelToken is a structural-ABSENCE check: the
// host-neutral dispatch core (fo-dispatch-core.md) must carry no LITERAL Claude-model
// token (`opus[1m]`, `sonnet`, `opus`, `haiku`) at EITHER leak site — reuse-condition-4
// or the Break-Glass template. The whole file is scanned, so both sites are covered.
// Codex and Pi load this core too, so a Claude model enum named here reads to them
// like a universal requirement they cannot satisfy. The expected value comes from the
// rule (these are literal host-coupled tokens that must not appear in the host-neutral
// core), not the file's own prose, so a host-neutral paraphrase that defers to "the
// host's canonical model enum" passes and a re-introduced enum fails — same family as
// TestDispatchCoreHasNoClaudeTeamImperative. The Claude realization legitimately lives
// in claude-fo-dispatch.md. The paired discriminator control keeps this non-vacuous.
func TestDispatchCoreHasNoClaudeModelToken(t *testing.T) {
	path := filepath.Join(skillsRoot(t), "first-officer", "references", "fo-dispatch-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read dispatch core %s: %v", path, err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		if lineLeaksClaudeModelToken(line) {
			t.Errorf("%s:%d carries a Claude-model token in the host-neutral core (Codex/Pi load this too); move the host realization to claude-fo-dispatch.md: %q", path, i+1, strings.TrimSpace(line))
		}
	}
}

// TestDispatchCoreModelTokenScannerDiscriminates is the DISCRIMINATOR control for the
// model-token absence check: it proves the scanner flags a genuinely host-coupled line
// and passes the host-neutral paraphrases the move produces — so
// TestDispatchCoreHasNoClaudeModelToken can never pass vacuously (e.g. by a typo'd
// token never matching anything).
func TestDispatchCoreModelTokenScannerDiscriminates(t *testing.T) {
	// Host-coupled line — the real pre-move leak shape. MUST flag.
	mustFlag := `stamped "opus[1m]" never matches sonnet/opus/haiku`
	if !lineLeaksClaudeModelToken(mustFlag) {
		t.Errorf("discriminator control: host-coupled model line was NOT flagged (the scanner would pass vacuously): %q", mustFlag)
	}

	// Host-neutral paraphrases the move produces — both MUST pass.
	reuseParaphrase := "a model stamped with a captain-session fallback value — one outside the host's canonical model enum — never matches an enum value"
	if lineLeaksClaudeModelToken(reuseParaphrase) {
		t.Errorf("discriminator control: the reuse-condition host-neutral paraphrase was wrongly flagged: %q", reuseParaphrase)
	}
	breakGlassParaphrase := "include it only when the stage declares a model in the host's canonical model enum; otherwise omit the entire model argument"
	if lineLeaksClaudeModelToken(breakGlassParaphrase) {
		t.Errorf("discriminator control: the break-glass host-neutral paraphrase was wrongly flagged: %q", breakGlassParaphrase)
	}
}

// TestEventLoopCoreHasNoPRScan is a structural-ABSENCE check: the host-neutral
// dispatch core (fo-dispatch-core.md) must carry no `gh pr view` PR-scan. PR lifecycle
// is the pr-merge mod's domain; a workflow with no pr-merge mod (a `merge: local` or
// non-code workflow) must never reach for `gh` in its loop. The token is the
// host-coupled defect (such a workflow can't satisfy it), not a paraphrasable meaning.
// SCOPING: the token is `gh pr view`, NOT `pr !=` — `--where "pr !="` is a legitimate
// status-query primitive (first-officer-shared-core.md) and banning it would be
// over-broad. The paired discriminator control keeps this non-vacuous.
func TestEventLoopCoreHasNoPRScan(t *testing.T) {
	path := filepath.Join(skillsRoot(t), "first-officer", "references", "fo-dispatch-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read dispatch core %s: %v", path, err)
	}
	if strings.Contains(string(data), "gh pr view") {
		t.Errorf("%s carries a `gh pr view` PR-scan in the host-neutral event loop; PR lifecycle is the pr-merge mod's domain (a merge:local workflow can't satisfy `gh`)", path)
	}
}

// lineLeaksPRScan reports whether a line carries the `gh pr view` PR-scan token. The
// discriminator drives this against a planted PR-scan line and the idle-hook line that
// legitimately remains in the loop.
func lineLeaksPRScan(line string) bool {
	return strings.Contains(line, "gh pr view")
}

// TestEventLoopPRScanScannerDiscriminates is the DISCRIMINATOR control for the
// event-loop PR-scan absence check: it proves the scanner flags a planted PR-scan line
// and passes the idle-hook-firing line that legitimately remains — so
// TestEventLoopCoreHasNoPRScan can never pass vacuously (a typo'd token would pass).
func TestEventLoopPRScanScannerDiscriminates(t *testing.T) {
	// Planted PR-scan line — the real pre-move leak shape. MUST flag.
	mustFlag := "check PR state via gh pr view and advance merged PRs"
	if !lineLeaksPRScan(mustFlag) {
		t.Errorf("discriminator control: planted PR-scan line was NOT flagged (the scanner would pass vacuously): %q", mustFlag)
	}

	// The idle-hook-firing line that legitimately remains in the generic loop. MUST pass.
	idleHookLine := "invoke «hooks.run»(\"idle\"), then «roster-reconcile»()"
	if lineLeaksPRScan(idleHookLine) {
		t.Errorf("discriminator control: the legitimately-remaining idle-hook line was wrongly flagged: %q", idleHookLine)
	}
}

// allowedPRViewFiles is the allow-list of shipped files permitted to name `gh pr view`:
// the canonical pr-merge mod (its startup + idle hooks legitimately scan). Modeled on
// TestNoUnexpectedModHookOrPRMergeIntroduced's allowedPRMergeFiles — the same allow-list
// discipline, adding `gh pr view` to it.
var allowedPRViewFiles = map[string]bool{
	filepath.Join("mods", "pr-merge.md"): true,
}

// prViewLeaksOutsideAllowList is the REAL allow-list scan, factored so the production
// check AND the negative discriminator control drive the SAME logic — not two inlined
// copies. It walks the given files, reads each, and returns the repo-relative paths that
// name `gh pr view` while absent from `allow`. Defeating this one function (widening
// `allow`, dropping the token match) reds both callers — that is what makes the negative
// control load-bearing rather than a re-implementation that can't fail.
func prViewLeaksOutsideAllowList(t *testing.T, root string, files []string, allow map[string]bool) []string {
	t.Helper()
	var leaks []string
	for _, path := range files {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			t.Errorf("rel %s: %v", path, err)
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		if strings.Contains(string(data), "gh pr view") && !allow[rel] {
			leaks = append(leaks, rel)
		}
	}
	return leaks
}

// TestNoUnexpectedPRViewScanIntroduced is a structural-ABSENCE check over the shipped
// surface: after the move, `gh pr view` appears ONLY in mods/pr-merge.md. Any other
// shipped file naming `gh pr view` fails. A no-pr-merge-mod workflow therefore loads no
// instruction that reaches `gh` in its loop. Two discriminator controls ship below.
func TestNoUnexpectedPRViewScanIntroduced(t *testing.T) {
	root := repoRoot(t)
	files := shippedInstructionMarkdown(t)
	if len(files) == 0 {
		t.Fatal("walked zero shipped instruction files — scope bug; the absence check would pass vacuously")
	}
	for _, rel := range prViewLeaksOutsideAllowList(t, root, files, allowedPRViewFiles) {
		t.Errorf("%s names `gh pr view` outside the canonical pr-merge mod — a no-pr-merge-mod workflow must reach no `gh` in its loop", rel)
	}
}

// TestPRViewAllowListIsLoadBearing is the POSITIVE discriminator control: it asserts
// mods/pr-merge.md DOES contain `gh pr view` (its startup + idle hooks legitimately
// scan), so the allow-list entry exempts a real occurrence rather than a vacuous one.
// If the mod ever stopped carrying the token this control reds. Modeled on
// TestPortabilityCheckDiscriminatesHostSpecific's load-bearing-exclusion assertion.
func TestPRViewAllowListIsLoadBearing(t *testing.T) {
	for path, want := range map[string]string{
		filepath.Join(repoRoot(t), "mods", "pr-merge.md"):                 "pr=pr-merge:{N}",
		filepath.Join(repoRoot(t), "docs", "dev", "_mods", "pr-merge.md"): "terminal status with `mod-block: merge:pr-merge`",
		filepath.Join(skillsRoot(t), "fo-gate-lifecycle", "SKILL.md"):     "terminal current status resumes the existing merge ceremony",
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
		} else if !strings.Contains(string(data), want) {
			t.Errorf("%s must contain %q", path, want)
		} else if strings.HasSuffix(path, filepath.Join("mods", "pr-merge.md")) && !strings.Contains(string(data), "gh pr view") {
			t.Errorf("%s no longer carries the load-bearing `gh pr view` scan", path)
		}
	}
}

// TestPRViewAllowListConstrains is the NEGATIVE discriminator control: it plants a real
// non-allowed shipped file carrying `gh pr view`, runs the REAL allow-list scan
// (prViewLeaksOutsideAllowList — the same function the production check drives), and
// asserts the scan FLAGS the plant. This is NOT a re-implementation of the scan: defeating
// the shared scan reds this control, so it cannot pass vacuously. Non-vacuity was verified
// by mutation: stubbing prViewLeaksOutsideAllowList to never flag turns this control RED via
// its scan-flagging assertion (recorded in the implementation stage report's cycle-2 section).
func TestPRViewAllowListConstrains(t *testing.T) {
	root := repoRoot(t)
	// Plant a real non-allowed shipped file under skills/ so the actual shipped-surface
	// walk picks it up; clean it up after.
	plantPath := filepath.Join(skillsRoot(t), "first-officer", "references", "zz-pr-view-negative-control.md")
	plantRel, err := filepath.Rel(root, plantPath)
	if err != nil {
		t.Fatalf("rel %s: %v", plantPath, err)
	}
	if allowedPRViewFiles[plantRel] {
		t.Fatalf("negative control invalid: planted path %q is on the allow-list", plantRel)
	}
	plantBody := "<!-- negative control -->\nFor each entity, check PR state via `gh pr view` and advance merged PRs.\n"
	if err := os.WriteFile(plantPath, []byte(plantBody), 0o644); err != nil {
		t.Fatalf("plant non-allowed file %s: %v", plantPath, err)
	}
	defer os.Remove(plantPath)

	files := shippedInstructionMarkdown(t)
	leaks := prViewLeaksOutsideAllowList(t, root, files, allowedPRViewFiles)
	found := false
	for _, rel := range leaks {
		if rel == plantRel {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("negative control: the REAL allow-list scan did NOT flag the planted non-allowed file %q (leaks=%v) — the allow-list does not constrain", plantRel, leaks)
	}
}
