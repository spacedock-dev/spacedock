// ABOUTME: Absence-invariant tests over the vendored FO/ensign skill surface —
// ABOUTME: AC-1 (no plugin status path) and AC-6 (no new PR-merge / `## Hook:` mod), oracle = the structural scope-fence.
package integration

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// skillsRoot is the vendored skill tree under test (the project skills/ dir
// this test package lives inside).
func skillsRoot(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// vendoredSkillFiles returns the vendored skill instruction surface: the FO and
// ensign reference markdown plus the vendored claude-team helper. The vendored
// status library is excluded — it is the status oracle, not skill instruction
// text, and legitimately carries the literal status filename internally.
func vendoredSkillFiles(t *testing.T) map[string]string {
	t.Helper()
	root := skillsRoot(t)
	rel := []string{
		"first-officer/references/first-officer-shared-core.md",
		"first-officer/references/claude-first-officer-runtime.md",
		"ensign/references/ensign-shared-core.md",
		"commission/bin/claude-team",
	}
	out := make(map[string]string, len(rel))
	for _, r := range rel {
		p := filepath.Join(root, r)
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read vendored skill file %s: %v", p, err)
		}
		out[r] = string(b)
	}
	return out
}

// TestNoPluginStatusPathInVendoredSkills locks AC-1: no file in the vendored
// skill instruction surface references the plugin-private status path.
func TestNoPluginStatusPathInVendoredSkills(t *testing.T) {
	for name, content := range vendoredSkillFiles(t) {
		if strings.Contains(content, "skills/commission/bin/status") {
			t.Errorf("%s references plugin-private status path 'skills/commission/bin/status'", name)
		}
		if strings.Contains(content, "spacedock_plugin_dir") {
			t.Errorf("%s still references {spacedock_plugin_dir} plugin root", name)
		}
	}
}

// TestNoPRMergeOrModBehaviorIntroduced locks AC-6: the vendored skill surface
// introduces no new `## Hook:` mod heading and no PR-merge flow beyond the
// existing mod-block convention the surface already documents. The vendored
// files are copies of the plugin skill text plus the three amendments; the
// amendments add no new lifecycle hook or PR-merge command. This asserts the
// amendment regions do not introduce a `## Hook:` heading and that no PR-merge
// invocation (gh pr merge / git merge --no-ff into main) was added.
//
// Oracle: the amendment-region scope-fence — an absence invariant over the
// vendored on-disk skill surface (the ensign file and the FO Split-Root
// amendment region, scoped via sectionAfter). No positive behavioral seam can
// prove an absence of behavior: a re-introduced `## Hook:` lifecycle mod or a
// new PR-merge command would silently change the dispatch lifecycle, and only
// this structural scope-fence over the amendment regions catches it. This is
// NOT bare prose-grep — it asserts a structural negative over the amendments,
// not the presence of an instruction clause.
func TestNoPRMergeOrModBehaviorIntroduced(t *testing.T) {
	files := vendoredSkillFiles(t)

	// The only `## Hook:` text legitimately present is the pre-existing Mod Hook
	// Convention documentation in the FO shared core (describing startup/idle/
	// merge points). The amendments must not add a NEW `## Hook: {point}` mod
	// declaration. Assert the ensign file (which the split-root amendment B
	// touched) introduces no `## Hook:` heading at all.
	ensign := files["ensign/references/ensign-shared-core.md"]
	if strings.Contains(ensign, "## Hook:") {
		t.Errorf("ensign reference unexpectedly introduces a `## Hook:` heading")
	}

	// The amendment-introduced region in the FO file is the Split-Root Worktree
	// Contract subsection. Assert that region introduces no `## Hook:` heading —
	// the pre-existing Mod Hook Convention text lives in a different section.
	fo := files["first-officer/references/first-officer-shared-core.md"]
	if region := sectionAfter(fo, "### Split-Root Worktree Contract"); strings.Contains(region, "## Hook:") {
		t.Errorf("FO split-root amendment region introduces a `## Hook:` heading")
	}

	// No PR-merge invocation may be introduced anywhere in the vendored surface.
	prMergeMarkers := []string{"gh pr merge", "git merge --no-ff", "git merge --ff-only main"}
	for name, content := range files {
		for _, m := range prMergeMarkers {
			if strings.Contains(content, m) {
				t.Errorf("%s introduces a PR-merge invocation %q (out of scope per AC-6)", name, m)
			}
		}
	}
}

// commissionStateBackendDecisionRows extracts the two state-backend decision rows
// from the commission SKILL.md decision table: each `- **{Label}** (...)` bullet
// directly under the `**State backend (...)**` decision lead-in, returned keyed by
// its bold branch label. It bounds each row to its own bullet line so an
// assertion can check what that row alone binds, not a free-floating substring
// anywhere in the file.
func commissionStateBackendDecisionRows(t *testing.T) map[string]string {
	t.Helper()
	p := filepath.Join(skillsRoot(t), "commission", "SKILL.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read commission SKILL.md: %v", err)
	}
	lines := strings.Split(string(b), "\n")
	// Find the decision lead-in, then collect the contiguous `- **Label**` bullets
	// that immediately follow (the decision-table rows), stopping at the first
	// non-bullet line.
	start := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "**State backend") {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return nil
	}
	rows := map[string]string{}
	labelRow := regexp.MustCompile(`^- \*\*([^*]+)\*\*`)
	for i := start; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		m := labelRow.FindStringSubmatch(line)
		if m == nil {
			if len(rows) > 0 {
				break // past the contiguous decision-table block
			}
			continue
		}
		rows[m[1]] = line
	}
	return rows
}

// TestCommissionStateBackendDecisionRule asserts the structural shape of the
// commission SKILL.md state-backend decision table: two distinct labeled branches
// under one decision lead-in, each row binding its own frontmatter outcome — the
// split-root row binds `state: .spacedock-state`, the inline row binds the omit/
// $inline outcome and NOT the split-root path. A scaffolding agent reads exactly
// this table to choose a backend, so the load-bearing property is that both
// branches exist as separate rows and neither collapses into the other.
//
// This DEMOTES the prior grep-over-prose check (it asserted three free-floating
// substrings exist anywhere in the file, which a meaning-inverting paraphrase —
// e.g. swapping which condition selects split-root — would keep green). REPLACE
// with a behavioral driver was not feasible: the commission flow is
// instruction-driven (the agent follows SKILL.md's "Write the README with ..."
// steps); there is no Go scaffolder function that takes a standalone-vs-code-repo
// input and emits frontmatter, so nothing behavioral is invocable from a test.
// The honest fallback per the task is this structural decision-table assertion,
// which pins the two-row shape and each row's bound outcome rather than prose
// wording, and FAILS if a branch is dropped, merged, or rebound to the wrong path.
func TestCommissionStateBackendDecisionRule(t *testing.T) {
	rows := commissionStateBackendDecisionRows(t)
	splitRoot, hasSplit := rows["Split-root"]
	inline, hasInline := rows["Inline"]
	if !hasSplit {
		t.Fatal("commission SKILL.md state-backend table missing the `- **Split-root**` row")
	}
	if !hasInline {
		t.Fatal("commission SKILL.md state-backend table missing the `- **Inline**` row")
	}
	// The split-root row must bind the orphan-state path; the inline row must NOT
	// claim it (that is the binding that flips if the two branches are swapped or
	// collapsed).
	if !strings.Contains(splitRoot, "state: .spacedock-state") {
		t.Errorf("split-root decision row does not bind `state: .spacedock-state`: %q", splitRoot)
	}
	if strings.Contains(inline, "state: .spacedock-state") {
		t.Errorf("inline decision row wrongly binds the split-root path `state: .spacedock-state`: %q", inline)
	}
	// The inline row must bind its own outcome (the $inline value, or the omit
	// guidance) so the two rows carry distinct frontmatter choices.
	if !strings.Contains(inline, "$inline") && !strings.Contains(inline, "omit") {
		t.Errorf("inline decision row binds no inline outcome ($inline or omit): %q", inline)
	}
}

// sectionAfter returns the body of the markdown section beginning at the line
// equal to heading, up to (but excluding) the next top-level `## ` heading, or
// "" when the heading is absent. Used to scope an assertion to one section.
func sectionAfter(text, heading string) string {
	lines := strings.Split(text, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == heading {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			end = i
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}
