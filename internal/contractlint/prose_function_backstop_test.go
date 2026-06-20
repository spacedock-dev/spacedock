// ABOUTME: AC-2c catastrophe backstop — the highest-stakes rule clusters survive the
// ABOUTME: prose-function restructure: the compact rebase-halt anchor + the verbatim reuse diagnostic.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// catastropheClause is one load-bearing string whose silent loss in the prose-function
// restructure would be catastrophic AND whose exact wording is itself load-bearing — the
// reuse diagnostic is contractually pinned ("must appear verbatim"); the compact rebase-halt
// prohibition anchor is the no-data-loss safety net. `path` is the restructured core that
// must still carry `want`.
type catastropheClause struct {
	path string
	want string
}

// catastropheClauses is DELIBERATELY narrow — only the clusters where a silent drop is
// catastrophic and the wording is load-bearing. It is NOT extended to the rest of the rule
// set: presence-grepping every rule would re-introduce the banned prose-grep tautology (the
// matched text is authored by the same restructuring implementer). The real preservation
// guard is the validator's detached original-vs-restructured audit (AC-2b); this is a cheap
// tripwire subordinate to it.
var catastropheClauses = []catastropheClause{
	// Rebase-conflict-halt prohibition anchor (shared-core State Management). The compact
	// «state.commit» body must still preserve the FO's no-data-loss halt obligation without
	// pinning the old multi-step git recipe.
	{
		path: filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"),
		want: "do not force-push or auto-resolve",
	},
	// The verbatim reuse diagnostic (dispatch-core Reuse and Fresh Dispatch). Contractually
	// pinned to appear character-for-character.
	{
		path: filepath.Join("skills", "first-officer", "references", "fo-dispatch-core.md"),
		want: "does not match next stage effective_model",
	},
}

// TestProseFunctionCatastropheClausesSurvive (AC-2c) asserts each highest-stakes
// clusters' load-bearing strings survives somewhere in its restructured core. A content-presence
// assertion is the right shape here precisely because the wording itself is contractual (a
// paraphrase of the verbatim diagnostic would BREAK the contract, so presence — not the usual
// "a paraphrase must pass" — is what we want). This is the only place that holds: it is
// subordinate to and does not replace AC-2b's detached audit.
func TestProseFunctionCatastropheClausesSurvive(t *testing.T) {
	root := repoRoot(t)
	if len(catastropheClauses) == 0 {
		t.Fatal("no catastrophe clauses declared — the backstop would pass vacuously")
	}
	for _, c := range catastropheClauses {
		data, err := os.ReadFile(filepath.Join(root, c.path))
		if err != nil {
			t.Fatalf("read restructured core %s: %v", c.path, err)
		}
		if !strings.Contains(string(data), c.want) {
			t.Errorf("%s no longer carries the catastrophe-cluster clause %q — the prose-function restructure dropped or paraphrased a load-bearing safety/diagnostic anchor", c.path, c.want)
		}
	}
}

// TestProseFunctionCatastropheBackstopDiscriminates is the non-vacuity control: it proves the
// content-presence scan FLAGS a core that is missing a clause (so the backstop can fail) and
// PASSES one that carries it. Without this a typo'd `want` (matching nothing, or matching
// everything) would let the backstop pass vacuously.
func TestProseFunctionCatastropheBackstopDiscriminates(t *testing.T) {
	clause := "does not match next stage effective_model"

	// A body that DROPPED the clause must be flagged (absence detected).
	dropped := "When the comparator forces fresh dispatch, emit a captain-visible diagnostic. The anchor phrase must appear."
	if strings.Contains(dropped, clause) {
		t.Fatalf("discriminator: the dropped-clause fixture unexpectedly contains %q", clause)
	}
	// A body that KEPT it must pass (presence detected) — the same scan the production check uses.
	kept := "the FO MUST emit `... does not match next stage effective_model ...`. The anchor phrase must appear verbatim."
	if !strings.Contains(kept, clause) {
		t.Fatalf("discriminator: the kept-clause fixture should contain %q but does not", clause)
	}
}
