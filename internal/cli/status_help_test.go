// ABOUTME: AC-3 `status --help` coverage — the query-flag synopsis renders
// ABOUTME: instead of the entity listing that ran there before the wantsHelp guard.
package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TestStatusHelpRendersQuerySynopsisNotEntityListing pins AC-3: `spacedock
// status --help` exits 0, prints help instead of running the entity listing
// (no default-table header), and states the four discoverability facts —
// --where named as THE entity query, the one-clause-per-flag AND rule, the
// canonical known-field list, and --archived as active-plus-archived (not a
// scope swap). Before the wantsHelp guard, this invocation ran the listing
// (exit 0, entity rows) instead of printing help.
func TestStatusHelpRendersQuerySynopsisNotEntityListing(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"status", "--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{
		"THE entity query",
		"Repeat --where to AND clauses",
		"id, slug, status, title, score, source, worktree, pr, started, completed,\nverdict, mod-block, archived, issue",
		"active PLUS archived, not archived-only",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("status --help missing %q:\n%s", want, out)
		}
	}
	// Neither the default table header nor the top-level grouped menu should
	// leak through: this is `status`'s own help, not the entity listing or the
	// root's menu.
	if strings.Contains(out, "ID     SLUG") {
		t.Errorf("status --help ran the entity listing instead of printing help:\n%s", out)
	}
	if strings.Contains(out, "Launch") {
		t.Errorf("status --help leaked the top-level grouped menu (Launch header):\n%s", out)
	}
}

// TestStatusShortHelpFlagAlsoRendersHelp locks that -h is equivalent to --help
// for status, matching wantsHelp's shared recognition of both spellings.
func TestStatusShortHelpFlagAlsoRendersHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"status", "-h"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "THE entity query") {
		t.Errorf("status -h missing query synopsis:\n%s", stdout.String())
	}
}
