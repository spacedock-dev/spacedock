// ABOUTME: AC-7/AC-8 structural guards for the «fn»-binding consolidation — the rebase-conflict
// ABOUTME: halt is one «fn» referenced by name, and the deferred registry keeps all four load-points + greet-guards.
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

// deferredRegistryBlock returns the `## Deferred Modules (registry)` section body (heading to
// next `## `). The block is the unit AC-8 checks, so the survival assertions cannot be satisfied
// by a load-point or greet-guard that survives ELSEWHERE in the core after a row was dropped.
func deferredRegistryBlock(t *testing.T, body string) string {
	t.Helper()
	const heading = "## Deferred Modules (registry)"
	start := strings.Index(body, heading)
	if start < 0 {
		t.Fatal("shared core has no `## Deferred Modules (registry)` block — the four-fold consolidation is missing")
	}
	rest := body[start+len(heading):]
	if end := regexp.MustCompile(`(?m)^## `).FindStringIndex(rest); end != nil {
		rest = rest[:end[0]]
	}
	return rest
}

// TestDeferredRegistryFoldsFourModulesWithLoadPointsAndGreetGuards (AC-8) closes the gap the
// closure walk leaves open: TestBootResidentDeferredLoadPointsResolve / ...CarryCeremony assert
// only that the four reference FILENAMES resolve on disk, never that a folded row kept its
// load-point trigger or its greet-guard. This walks the single registry block and asserts all
// four module reference paths, all four load-point triggers, and a per-row greet-guard for each
// of the four modules survive by name — so dropping a row (or its trigger/guard) reds here while
// the closure tests stay green.
func TestDeferredRegistryFoldsFourModulesWithLoadPointsAndGreetGuards(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, sharedCorePath()))
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	block := deferredRegistryBlock(t, string(data))

	for _, ref := range []string{
		"references/fo-status-viewer.md",
		"references/fo-dispatch-core.md",
		"references/fo-write-core.md",
		"references/fo-merge-core.md",
	} {
		if !strings.Contains(block, ref) {
			t.Errorf("deferred registry does not fold %s — a module reference path is missing", ref)
		}
	}
	for _, lp := range []string{
		"FIRST status query",    // status-query load-point
		"FIRST worker dispatch", // dispatch load-point
		"FIRST write to main",   // write/new-entity load-point
		"terminal boundary",     // merge load-point
	} {
		if !strings.Contains(block, lp) {
			t.Errorf("deferred registry lost load-point trigger %q — a folded row's done-when vanished", lp)
		}
	}
	if n := strings.Count(block, "greet-guard:"); n < 4 {
		t.Errorf("deferred registry has %d per-row `greet-guard:` clauses, want >=4 (one per folded module) — a folded row's greet-guard vanished", n)
	}
}
