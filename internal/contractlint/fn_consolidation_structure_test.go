// ABOUTME: AC-7/AC-8 structural guards for the «fn»-binding consolidation — the rebase-conflict
// ABOUTME: halt is one «fn» referenced by name, and the deferred load points keep all four modules + the shared greet-guard.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func sharedCorePath() string {
	return filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md")
}

// TestRebaseConflictHaltConsolidatedToFn (AC-7) is the structural guard the os.Stat closure
// cannot be: the rebase-conflict halt is DEFINED once as a `## «halt.rebase-conflict»` «fn» and
// the FO call sites resolve to it BY NAME, never by restating the abort recipe inline. The
// independent expected values are fixed counts (one definition, three-plus references, the
// abort verb exactly once), so a second definition, a dropped reference, or a call site that
// re-pastes `git rebase --abort` reds. claude-fo-dispatch.md must also reference it by name.
func TestRebaseConflictHaltConsolidatedToFn(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, sharedCorePath()))
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	body := string(data)

	defs := regexp.MustCompile(`(?m)^## «halt\.rebase-conflict»`).FindAllString(body, -1)
	if len(defs) != 1 {
		t.Errorf("expected exactly 1 `## «halt.rebase-conflict»` definition in shared core, found %d", len(defs))
	}
	if refs := strings.Count(body, "«halt.rebase-conflict»") - len(defs); refs < 3 {
		t.Errorf("expected >=3 by-name «halt.rebase-conflict» references at FO call sites in shared core, found %d", refs)
	}
	// 0 inline restatements at the call sites: the abort verb appears ONLY in the one «fn»
	// definition body. A call site that re-pastes it (or a definition that drops it) breaks != 1.
	if n := strings.Count(body, "git rebase --abort"); n != 1 {
		t.Errorf("expected `git rebase --abort` exactly once (the «halt.rebase-conflict» definition) in shared core, found %d — a call site restated the abort recipe or the definition dropped it", n)
	}

	cf := filepath.Join("skills", "first-officer", "references", "claude-fo-dispatch.md")
	cfdata, err := os.ReadFile(filepath.Join(root, cf))
	if err != nil {
		t.Fatalf("read claude-fo-dispatch: %v", err)
	}
	if !strings.Contains(string(cfdata), "«halt.rebase-conflict»") {
		t.Errorf("%s does not reference «halt.rebase-conflict» by name — its rebase-halt site did not resolve to the «fn»", cf)
	}
}

// deferredLoadPointsBlock returns the `## Deferred load points` section body (heading to
// next `## `). The block is the unit AC-8 checks, so the survival assertions cannot be satisfied
// by a load-point or greet-guard that survives ELSEWHERE in the core after a row was dropped.
func deferredLoadPointsBlock(t *testing.T, body string) string {
	t.Helper()
	// Anchor the heading at start-of-line so a prose MENTION of `## Deferred load points`
	// elsewhere in the core (e.g. an intro forward-reference) cannot be mistaken for the block.
	loc := regexp.MustCompile(`(?m)^## Deferred load points$`).FindStringIndex(body)
	if loc == nil {
		t.Fatal("shared core has no `## Deferred load points` block — the four-fold collapse is missing")
	}
	rest := body[loc[1]:]
	if end := regexp.MustCompile(`(?m)^## `).FindStringIndex(rest); end != nil {
		rest = rest[:end[0]]
	}
	return rest
}

// TestDeferredLoadPointsFoldModulesWithTriggersAndGreetGuard (AC-8) closes the gap the
// closure walk leaves open: TestBootResidentDeferredLoadPointsResolve / ...CarryCeremony /
// ...SkillCoresResolveAndCarryCeremony assert only that the deferred module targets resolve on disk,
// never that a folded row kept its load-point trigger or the shared greet-guard. This walks the
// single load-points block and asserts the remaining module tokens and triggers plus the single
// shared greet-guard survive by name. The merge core is eagerly imported by first-officer/SKILL.md,
// so it must not reappear here as a deferred read.
func TestDeferredLoadPointsFoldModulesWithTriggersAndGreetGuard(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, sharedCorePath()))
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	block := deferredLoadPointsBlock(t, string(data))

	for _, tok := range []string{
		"spacedock:fo-status-viewer",     // status-viewer skill
		"references/fo-dispatch-core.md", // dispatch reference
	} {
		if !strings.Contains(block, tok) {
			t.Errorf("deferred load points do not fold %s — a module load-point token is missing", tok)
		}
	}
	for _, lp := range []string{
		"first status query",    // status-query load-point
		"first worker dispatch", // dispatch load-point
	} {
		if !strings.Contains(block, lp) {
			t.Errorf("deferred load points lost load-point trigger %q — a folded row's trigger vanished", lp)
		}
	}
	// The collapse replaces the four per-row greet-guards with ONE shared clause. Its survival
	// is the invariant; a block that dropped it reads as if a greet-and-stop boot could load a
	// deferred module.
	if !strings.Contains(block, "greet-and-stop boot loads NONE") {
		t.Error("deferred load points lost the shared greet-guard clause (\"a greet-and-stop boot loads NONE of these\")")
	}
}
