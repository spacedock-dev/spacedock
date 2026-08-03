// ABOUTME: AC-1 stage-then-score default-sort fixture — proves the descending
// ABOUTME: stage/score comparator and its discovery-order tie-break end-to-end.
package status

import (
	"path/filepath"
	"testing"
)

// sortFixtureReadme declares four working stages (backlog, ideation,
// implementation, validation, in that dispatch order) plus a terminal done, so
// the default listing's reversed stage-order comparator has a real
// multi-stage taxonomy to sort against.
const sortFixtureReadme = `---
commissioned-by: spacedock@1
id-style: slug
stages:
  states:
    - name: backlog
      initial: true
    - name: ideation
    - name: implementation
      worktree: true
    - name: validation
    - name: done
      terminal: true
---
# Sort Fixture Workflow
`

// buildSortFixture writes the sort-order fixture entities: two rows each for
// backlog/ideation/validation, plus three implementation rows — a high score
// and a tied pair (impl-lo-a/impl-lo-b at the same score) whose alphabetical
// slug order is the only tie-break the stable sort has to preserve.
func buildSortFixture(t *testing.T) string {
	t.Helper()
	def := t.TempDir()
	writeFile(t, filepath.Join(def, "README.md"), sortFixtureReadme)
	entities := map[string]string{
		"backlog-a.md":     "---\nstatus: backlog\nscore: \"0.80\"\n---\n",
		"backlog-b.md":     "---\nstatus: backlog\nscore: \"0.30\"\n---\n",
		"ideation-a.md":    "---\nstatus: ideation\nscore: \"0.60\"\n---\n",
		"ideation-b.md":    "---\nstatus: ideation\nscore: \"0.20\"\n---\n",
		"impl-hi.md":       "---\nstatus: implementation\nscore: \"0.90\"\n---\n",
		"impl-lo-a.md":     "---\nstatus: implementation\nscore: \"0.40\"\n---\n",
		"impl-lo-b.md":     "---\nstatus: implementation\nscore: \"0.40\"\n---\n",
		"validation-lo.md": "---\nstatus: validation\nscore: \"0.10\"\n---\n",
		"validation-hi.md": "---\nstatus: validation\nscore: \"0.95\"\n---\n",
	}
	for name, content := range entities {
		writeFile(t, filepath.Join(def, name), content)
	}
	return def
}

// TestSortDefaultStageThenScore (AC-1) proves the default listing sorts later
// stages first, score descending within a stage, and keeps the discovery
// (slug) order for an exact stage+score tie. Row order — not membership — is
// the assertion: a comparator that left stage order ascending, or that
// resolved the implementation/implementation tie some other way, fails this
// exact-order check. In particular: validation-lo (score 0.10) must precede
// impl-hi (score 0.90) — stage beats score — and impl-hi must precede the two
// tied 0.40 implementation rows, which must render in slug order.
func TestSortDefaultStageThenScore(t *testing.T) {
	def := buildSortFixture(t)
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, def, env, "--workflow-dir", def)
	if code != 0 {
		t.Fatalf("default status exit=%d stderr=%q", code, stderr)
	}

	got := tableSlugs(t, out)
	want := []string{
		"validation-hi", "validation-lo",
		"impl-hi", "impl-lo-a", "impl-lo-b",
		"ideation-a", "ideation-b",
		"backlog-a", "backlog-b",
	}
	if !equalStrings(got, want) {
		t.Fatalf("default sort order = %v, want %v\n%s", got, want, out)
	}
}
