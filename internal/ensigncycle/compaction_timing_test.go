// ABOUTME: Offline fixture for Rule 1 (safe-to-compact timing) — pairs a durable and
// ABOUTME: an unrecoverable boundary with real commit OIDs and proves the suggestion timing.
package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// compactionSuggestionRe recognizes the safe-to-compact suggestion the after-compaction
// rule authorizes. It keys on the load-bearing "safe time to compact" clause so the
// oracle counts real suggestions, not the word "compact" appearing in narration.
var compactionSuggestionRe = regexp.MustCompile(`(?i)safe time to compact`)

// boundaryWorkItem is one unit of FO-owned work at a compaction boundary and whether
// it is durably committed. A committed item carries the real git OID of its commit —
// the "relevant state/worktree commit OID before the suggestion" AC-1 requires as
// evidence, not a boolean the test author asserts by hand.
type boundaryWorkItem struct {
	desc      string
	committed bool
	oid       string
}

// compactionBoundary is the durability state the FO reads before deciding whether to
// suggest compaction: every FO-owned change committed, and no received completion,
// gate decision, or transition awaiting reconciliation.
type compactionBoundary struct {
	workItems             []boundaryWorkItem
	pendingReconciliation bool // an unconsumed worker completion / gate decision
}

// durable reports whether the boundary is recoverable without conversational memory.
func (b compactionBoundary) durable() (bool, string) {
	if b.pendingReconciliation {
		return false, "an unconsumed completion/gate decision awaits reconciliation"
	}
	for _, it := range b.workItems {
		if !it.committed || it.oid == "" {
			return false, "uncommitted FO-owned work: " + it.desc
		}
	}
	return true, ""
}

func countCompactionSuggestions(captainMessage string) int {
	return len(compactionSuggestionRe.FindAllString(captainMessage, -1))
}

// assertCompactionSuggestionTiming is the Rule 1 oracle. The invariant it enforces is
// the timing value, not the presence of rule text: a safe-to-compact suggestion may
// appear only at a durable, recoverable boundary, and only riding on committed OID
// evidence for every FO-owned work item. A suggestion at an unrecoverable boundary —
// the baseline that moves the wrong way — is the failure.
func assertCompactionSuggestionTiming(b compactionBoundary, captainMessage string) error {
	suggestions := countCompactionSuggestions(captainMessage)
	durable, reason := b.durable()
	if suggestions > 0 && !durable {
		return fmt.Errorf("safe-to-compact suggestion emitted at an unrecoverable boundary (%s)", reason)
	}
	if suggestions > 0 {
		for _, it := range b.workItems {
			if it.oid == "" {
				return fmt.Errorf("suggestion emitted without commit-OID evidence for %q", it.desc)
			}
		}
	}
	if suggestions > 1 {
		return fmt.Errorf("emitted %d safe-to-compact suggestions; the cue is a single non-blocking hint", suggestions)
	}
	return nil
}

// commitWorkItem writes a file into a git repo and commits it path-scoped, returning
// the item with the real commit OID — the "relevant state/worktree commit OID before
// the suggestion" AC-1 requires as evidence. Reuses the package git helpers.
func commitWorkItem(t *testing.T, repo, name, content string) boundaryWorkItem {
	t.Helper()
	if err := os.WriteFile(filepath.Join(repo, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	gitCommitPathScoped(t, repo, name, "commit "+name)
	oid := strings.TrimSpace(git(t, repo, "rev-parse", "HEAD"))
	return boundaryWorkItem{desc: name, committed: true, oid: oid}
}

// newCompactionRepo initializes an empty git repo for the fixture; the first
// commitWorkItem lands the first commit.
func newCompactionRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	git(t, repo, "init", "-q")
	return repo
}

// TestCompactionSafeBoundarySuggestsOnce is AC-1's durable case: both FO-owned changes
// are committed (real OIDs), nothing awaits reconciliation, and the FO emits exactly
// one non-blocking safe-to-compact suggestion.
func TestCompactionSafeBoundarySuggestsOnce(t *testing.T) {
	repo := newCompactionRepo(t)
	boundary := compactionBoundary{workItems: []boundaryWorkItem{
		commitWorkItem(t, repo, "entity.md", "status: implementation\n"),
		commitWorkItem(t, repo, "report.md", "## Stage Report\n"),
	}}
	captainMessage := "Context is getting tight. Current Spacedock state is durable at a clean boundary; now is a safe time to compact."

	if got := countCompactionSuggestions(captainMessage); got != 1 {
		t.Fatalf("safe scenario should emit exactly one suggestion, got %d", got)
	}
	if err := assertCompactionSuggestionTiming(boundary, captainMessage); err != nil {
		t.Fatalf("durable boundary with one suggestion must pass: %v", err)
	}
}

// TestCompactionUnsafeBoundaryStaysSilent is AC-1's unsafe case: an unconsumed worker
// completion awaits reconciliation, so the FO emits zero suggestions.
func TestCompactionUnsafeBoundaryStaysSilent(t *testing.T) {
	repo := newCompactionRepo(t)
	boundary := compactionBoundary{
		workItems:             []boundaryWorkItem{commitWorkItem(t, repo, "entity.md", "status: implementation\n")},
		pendingReconciliation: true,
	}
	captainMessage := "A worker just reported; verifying its stage report before anything else."

	if got := countCompactionSuggestions(captainMessage); got != 0 {
		t.Fatalf("unsafe scenario must emit zero suggestions, got %d", got)
	}
	if err := assertCompactionSuggestionTiming(boundary, captainMessage); err != nil {
		t.Fatalf("unsafe boundary with no suggestion must pass: %v", err)
	}
}

// TestCompactionSuggestionAtUnrecoverableBoundaryRejected proves the oracle bites: a
// suggestion fired while an uncommitted change or unconsumed completion is pending is
// the exact failure AC-1 measures, and must be rejected.
func TestCompactionSuggestionAtUnrecoverableBoundaryRejected(t *testing.T) {
	repo := newCompactionRepo(t)
	suggestion := "Now is a safe time to compact."

	uncommitted := compactionBoundary{workItems: []boundaryWorkItem{
		commitWorkItem(t, repo, "entity.md", "x"),
		{desc: "uncommitted report edit", committed: false},
	}}
	if err := assertCompactionSuggestionTiming(uncommitted, suggestion); err == nil {
		t.Fatalf("a suggestion with an uncommitted change must be rejected")
	}

	pending := compactionBoundary{
		workItems:             []boundaryWorkItem{commitWorkItem(t, repo, "entity2.md", "y")},
		pendingReconciliation: true,
	}
	if err := assertCompactionSuggestionTiming(pending, suggestion); err == nil {
		t.Fatalf("a suggestion with an unconsumed completion must be rejected")
	}
}
