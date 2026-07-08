// ABOUTME: Value-binding checks for the fo-dispatch-recovery skill: the
// ABOUTME: live-scenario oracle's captain-report constant, and its two cross-file pointers.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// degradedModeCaptainReportConstRe extracts the Go string-literal value of
// degradedModeCaptainReportPrefix from its declaration in
// internal/ensigncycle/dispatch_recovery_assert_test.go. It is a single-line
// `const NAME = "..."` string literal, so a regex is sufficient — no AST needed
// for one literal.
var degradedModeCaptainReportConstRe = regexp.MustCompile(`degradedModeCaptainReportPrefix\s*=\s*"((?:[^"\\]|\\.)*)"`)

// TestDegradedModeCaptainReportPrefixBindsSkillBlockquote is a VALUE-BINDING
// check, not a prose-grep and not the banned CODE-BOUND-AS-BEHAVIOR-SUBSTITUTE
// (see the boundary-guard package doc): it never claims the skill's prose PROVES
// AC-2's behavior — the `degraded-bare` live scenario and its offline oracle
// (`assertDegradedBareObservables`) are what prove that, by running against a
// real captured stream. What this binds is narrower: two INDEPENDENTLY EDITABLE
// values — the Go constant the oracle matches stream text against, and the
// actual opening of `fo-dispatch-recovery/SKILL.md`'s `## Degraded Mode` Captain
// Report Template blockquote — must not silently drift apart. If either is
// paraphrased without the other, the oracle stops matching real output and every
// future live run either false-fails (a correct FO, paraphrased skill, stale
// constant) or, worse, the drift goes undetected because no offline test
// exercises the pairing. That is the same shape as a ref-closure or dedup
// check: a machine-checkable fact about two on-disk sources, not an
// interpretation of what the prose means.
func TestDegradedModeCaptainReportPrefixBindsSkillBlockquote(t *testing.T) {
	root := repoRoot(t)

	assertPath := filepath.Join(root, "internal", "ensigncycle", "dispatch_recovery_assert_test.go")
	assertSrc, err := os.ReadFile(assertPath)
	if err != nil {
		t.Fatalf("read %s: %v", assertPath, err)
	}
	m := degradedModeCaptainReportConstRe.FindSubmatch(assertSrc)
	if m == nil {
		t.Fatalf("degradedModeCaptainReportPrefix constant not found in %s — the value-binding check has nothing to bind", assertPath)
	}
	prefix := string(m[1])
	if prefix == "" {
		t.Fatal("extracted an empty degradedModeCaptainReportPrefix — extraction bug, the check would pass vacuously")
	}

	skillPath := filepath.Join(root, "skills", "fo-dispatch-recovery", "SKILL.md")
	skillBody, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read %s: %v", skillPath, err)
	}
	if !strings.Contains(string(skillBody), prefix) {
		t.Errorf("the oracle's degradedModeCaptainReportPrefix (%q) does not appear verbatim in %s's Captain Report Template — the skill's report sentence and the offline oracle's match constant have diverged", prefix, skillPath)
	}
}

// usingLegacyClaudeTeamDegradedModePointerMarkers are the two re-pointed
// Degraded-Mode pointer lines in using-legacy-claude-team/SKILL.md, identified by
// a stable lead phrase that survived the re-point (present in both the OLD bare
// pointer text and the CURRENT Skill-token text), so this pin locates the same
// two lines regardless of unrelated prose edits elsewhere in the file.
var usingLegacyClaudeTeamDegradedModePointerMarkers = []string{
	"**Failure handling:**",
	"**Fall back to Degraded Mode**",
}

// TestUsingLegacyClaudeTeamDegradedModePointersNameRealAnchor is a narrower,
// token-scoped alternative to walking using-legacy-claude-team/SKILL.md through
// the generic referenceProsePointerDanglers scanner (TestDeferredSkillProsePointersResolve).
// That generic walk was tried first and rejected: it also reds on :50/:52
// ("Degraded Mode cannot be entered" / "Degraded Mode is entered"), which name
// Degraded Mode as a STATE, not as a pointer at the moved section — the scanner
// has no way to tell the two apart, and those two lines are deliberately
// untouched prose the entity's design already reviewed as "checked and valid
// as-is." Generalizing the scanner to make that distinction is out of scope here
// (see the implementation cycle-1 amendment's residual observation).
//
// Instead, this pin locates the file's two KNOWN pointer lines by a stable
// marker phrase (not by token presence, so a regression to the OLD bare-pointer
// form — which drops the token entirely — is still caught) and asserts each one
// (1) still carries the Skill(skill="spacedock:fo-dispatch-recovery") token, and
// (2) names a section that is a real anchor in fo-dispatch-recovery/SKILL.md.
// The two marker-mentioned lines are exactly what the audit cared about; :50/:52
// carry neither marker, so they never enter this check.
func TestUsingLegacyClaudeTeamDegradedModePointersNameRealAnchor(t *testing.T) {
	root := repoRoot(t)

	recoverySkillRel := filepath.Join("skills", "fo-dispatch-recovery", "SKILL.md")
	anchors, ok := deferredSkillCores[recoverySkillRel]
	if !ok || len(anchors) == 0 {
		t.Fatalf("deferredSkillCores has no anchors for %s — the pin has nothing to bind against", recoverySkillRel)
	}
	recoveryBody, err := os.ReadFile(filepath.Join(root, recoverySkillRel))
	if err != nil {
		t.Fatalf("read %s: %v", recoverySkillRel, err)
	}

	legacyRel := filepath.Join("skills", "using-legacy-claude-team", "SKILL.md")
	legacyData, err := os.ReadFile(filepath.Join(root, legacyRel))
	if err != nil {
		t.Fatalf("read %s: %v", legacyRel, err)
	}
	lines := strings.Split(string(legacyData), "\n")

	for _, marker := range usingLegacyClaudeTeamDegradedModePointerMarkers {
		lineNo, line, found := findLineContaining(lines, marker)
		if !found {
			t.Errorf("%s: no line carries the marker %q — the pointer line this pin tracks has moved or been deleted", legacyRel, marker)
			continue
		}
		if !strings.Contains(line, "spacedock:fo-dispatch-recovery") {
			t.Errorf("%s:%d (marker %q) no longer carries Skill(skill=\"spacedock:fo-dispatch-recovery\") — it has regressed to a bare pointer with no resolving token: %q", legacyRel, lineNo, marker, strings.TrimSpace(line))
			continue
		}
		var named string
		for _, anchor := range anchors {
			name := strings.TrimLeft(anchor, "# ")
			if strings.Contains(line, name) {
				named = name
				break
			}
		}
		if named == "" {
			t.Errorf("%s:%d (marker %q) carries the recovery-skill token but names none of its section anchors %v: %q", legacyRel, lineNo, marker, anchors, strings.TrimSpace(line))
			continue
		}
		if !strings.Contains(string(recoveryBody), "## "+named) {
			t.Errorf("%s:%d (marker %q) names section %q, but %s carries no such heading — the pointer target has drifted", legacyRel, lineNo, marker, named, recoverySkillRel)
		}
	}
}

// findLineContaining returns the 1-indexed line number and text of the first
// line in lines containing marker.
func findLineContaining(lines []string, marker string) (lineNo int, line string, found bool) {
	for i, l := range lines {
		if strings.Contains(l, marker) {
			return i + 1, l, true
		}
	}
	return 0, "", false
}
