// ABOUTME: Lock-in oracle that the universal ensign contract carries no dev-only
// ABOUTME: authoring discipline, and that the dev disciplines survive in their dev homes.
package hostneutrality

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// devLeakageLiterals are dev-workflow-specific phrases that must NOT appear in
// the universal shared cores. They fingerprint the test-first authoring bullet
// and the code-substrate worktree noun that leaked into the portable ensign
// contract. "failing test" is deliberately absent — it appears legitimately at
// first-officer-shared-core.md ("enforced by the binary or a failing test"),
// and a flat all-files table would false-red on that legit FO usage.
var devLeakageLiterals = []struct {
	literal       string
	caseSensitive bool
}{
	// AC-1: the test-first authoring bullet's unique fingerprints.
	{"feature or bugfix", false},
	{"the test is what the gate judges", false},
	// AC-3: the code-substrate worktree noun, banned as the exact uppercased
	// form it appears in today.
	{"CODE only", true},
}

// devLeakageCorePaths are the universal cores swept for dev-leakage. The
// runtime adapters are policed separately by the AC-4 field-enumeration check.
var devLeakageCorePaths = []string{
	filepath.Join("..", "..", "skills", "first-officer", "references", "first-officer-shared-core.md"),
	filepath.Join("..", "..", "skills", "ensign", "references", "ensign-shared-core.md"),
}

// TestNoDevLeakageInUniversalCore locks AC-1 and AC-3: the universal shared
// cores must not assert dev-only authoring discipline (test-first) nor name the
// code substrate in the worktree-isolation clause. A re-introduction of the
// banned literal fails the test (negative proof of lock-in).
func TestNoDevLeakageInUniversalCore(t *testing.T) {
	for _, path := range devLeakageCorePaths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			body, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			text := string(body)
			lowerText := strings.ToLower(text)

			for _, banned := range devLeakageLiterals {
				hit := false
				if banned.caseSensitive {
					hit = strings.Contains(text, banned.literal)
				} else {
					hit = strings.Contains(lowerText, strings.ToLower(banned.literal))
				}
				// "feature or bugfix" and "the test is what the gate judges"
				// uniquely fingerprint the ensign TDD bullet; "CODE only" is in
				// both cores today. All three must be absent post-sweep.
				if hit {
					t.Errorf("%s contains banned dev-leakage literal %q — dev-only discipline re-homed into the universal core", path, banned.literal)
				}
			}
		})
	}
}

// TestWorktreeIsolationClauseSurvives locks AC-3's other half: removing the
// "CODE only" noun must NOT delete the worktree-isolation boundary — both cores
// must still carry an isolation clause naming the worktree.
func TestWorktreeIsolationClauseSurvives(t *testing.T) {
	for _, path := range devLeakageCorePaths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			body, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			lower := strings.ToLower(string(body))
			if !strings.Contains(lower, "worktree isolate") && !strings.Contains(lower, "isolate") {
				t.Errorf("%s lost its worktree-isolation clause — neutralizing the substrate noun must not delete the boundary", path)
			}
			if !strings.Contains(lower, "worktree") {
				t.Errorf("%s no longer mentions the worktree — isolation boundary gone", path)
			}
		})
	}
}

// fieldEnumerationRe captures the runtime-adapter sentence that enumerates the
// always-present assignment fields, from "authoritative for all assignment
// fields:" to the end of the sentence. The AC-4 check scopes ONLY to this
// sentence so the legitimate conditional "if you were given a worktree path"
// at ensign-shared-core.md:17 stays untouched.
var fieldEnumerationRe = regexp.MustCompile(`(?i)authoritative for all assignment fields:[^.]*\.`)

// runtimeAdapterFieldPaths are the two ensign runtime adapters whose
// field-enumeration sentence must use the neutral location vocabulary.
var runtimeAdapterFieldPaths = []string{
	filepath.Join("..", "..", "skills", "ensign", "references", "claude-ensign-runtime.md"),
	filepath.Join("..", "..", "skills", "ensign", "references", "codex-ensign-runtime.md"),
}

// TestRuntimeAdaptersUseNeutralLocationVocabulary locks AC-4: each runtime
// adapter's assignment-field enumeration lists "workflow location" and NOT
// "worktree path". Scoped to the field-enumeration sentence — a file-wide ban
// would false-fail on the legitimate conditional usage elsewhere.
func TestRuntimeAdaptersUseNeutralLocationVocabulary(t *testing.T) {
	for _, path := range runtimeAdapterFieldPaths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			body, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			loc := fieldEnumerationRe.FindString(string(body))
			if loc == "" {
				t.Fatalf("%s has no assignment-field enumeration sentence anchored on 'authoritative for all assignment fields:'", path)
			}
			lower := strings.ToLower(loc)
			if strings.Contains(lower, "worktree path") {
				t.Errorf("%s field enumeration still lists 'worktree path' — use the neutral 'workflow location'; sentence=%q", path, loc)
			}
			if !strings.Contains(lower, "workflow location") {
				t.Errorf("%s field enumeration does not list 'workflow location'; sentence=%q", path, loc)
			}
		})
	}
}

// devHomePresence maps each dev home to a required clause proving the lift
// RELOCATED rather than DELETED the dev discipline. Inverse polarity of the
// banned-literal tables: this errors on ABSENCE.
var devHomePresence = []struct {
	path      string
	required  string
	sectionRe *regexp.Regexp // when non-nil, the required clause must sit inside this section
}{
	{
		path:      filepath.Join("..", "..", "skills", "commission", "references", "templates", "development.md"),
		required:  "test-first",
		sectionRe: regexp.MustCompile(`(?is)## Recommended practices \(opt-in\).*`),
	},
	{
		path:     filepath.Join("..", "..", "docs", "dev", "README.md"),
		required: "real, checkable change",
	},
}

// TestDevDisciplinesSurviveInDevHomes locks AC-2: the re-homed dev guidance is
// present in its dev-specific home — development.md carries a test-first clause
// inside its Recommended-practices section, and docs/dev/README.md retains the
// "real, checkable change" deliverable-proof policy. Fails if a future edit
// strips a dev home.
func TestDevDisciplinesSurviveInDevHomes(t *testing.T) {
	for _, h := range devHomePresence {
		t.Run(filepath.Base(h.path), func(t *testing.T) {
			body, err := os.ReadFile(h.path)
			if err != nil {
				t.Fatalf("read %s: %v", h.path, err)
			}
			scope := string(body)
			if h.sectionRe != nil {
				scope = h.sectionRe.FindString(scope)
				if scope == "" {
					t.Fatalf("%s missing the section that must hold %q", h.path, h.required)
				}
			}
			if !strings.Contains(strings.ToLower(scope), strings.ToLower(h.required)) {
				t.Errorf("%s lost the re-homed dev clause %q — the lift must relocate, not delete", h.path, h.required)
			}
		})
	}
}
