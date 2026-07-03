package release

import (
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// TestLatestObservationsSelectsByCapturedAtNotArrayPosition is AC-3's proof that
// baseline selection is driven by the captured_at timestamp, not by "last in the
// array": the observations are deliberately given OUT-OF-ORDER captured_at
// values (the chronologically latest one is authored FIRST in the slice), so a
// naive "last element wins" implementation would pick the wrong observation.
func TestLatestObservationsSelectsByCapturedAtNotArrayPosition(t *testing.T) {
	older := journeymetrics.Record{
		ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
		Turns: 18, Tokens: journeymetrics.TokenTotals{CacheCreation: 1562, Total: 1562},
		CapturedAt: "2026-06-20T00:00:00Z", RunID: "27931963802",
		RunURL: "https://github.com/spacedock-dev/spacedock/actions/runs/27931963802",
	}
	newer := journeymetrics.Record{
		ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
		Turns: 20, Tokens: journeymetrics.TokenTotals{CacheCreation: 1576, Total: 1576},
		CapturedAt: "2026-06-27T00:00:00Z", RunID: "28432388663",
		RunURL: "https://github.com/spacedock-dev/spacedock/actions/runs/28432388663",
	}
	ledger := journeymetrics.Ledger{
		Scenarios: []journeymetrics.ScenarioLedgerEntry{
			{
				ScenarioID: "shallow-boot-window",
				// newer is authored FIRST, older SECOND — the reverse of
				// chronological order — so an array-position-based ("last wins")
				// selection would wrongly pick `older`.
				Observations: []journeymetrics.Record{newer, older},
			},
		},
	}

	latest := LatestObservations(ledger)
	got, ok := latest[deltaKeyString(journeyDeltaKey{ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6"})]
	if !ok {
		t.Fatal("no baseline observation selected for shallow-boot-window/claude/claude-sonnet-4-6")
	}
	if got.RunID != "28432388663" {
		t.Fatalf("selected baseline RunID = %q, want the chronologically-latest 28432388663 (got Turns=%d, CapturedAt=%s) — selection is not timestamp-based", got.RunID, got.Turns, got.CapturedAt)
	}
}

// TestLatestObservationsTreatsUnparseableCapturedAtAsOldest proves a record with
// no (or malformed) captured_at never wins over a genuinely-timestamped peer —
// it sorts as the oldest possible time rather than erroring or comparing equal.
func TestLatestObservationsTreatsUnparseableCapturedAtAsOldest(t *testing.T) {
	noTimestamp := journeymetrics.Record{
		ScenarioID: "gate-guardrail", Runtime: "claude", Model: "claude-sonnet-4-6", Turns: 1,
	}
	timestamped := journeymetrics.Record{
		ScenarioID: "gate-guardrail", Runtime: "claude", Model: "claude-sonnet-4-6", Turns: 2,
		CapturedAt: "2026-01-01T00:00:00Z",
	}
	ledger := journeymetrics.Ledger{
		Scenarios: []journeymetrics.ScenarioLedgerEntry{
			{ScenarioID: "gate-guardrail", Observations: []journeymetrics.Record{noTimestamp, timestamped}},
		},
	}
	latest := LatestObservations(ledger)
	got := latest[deltaKeyString(journeyDeltaKey{ScenarioID: "gate-guardrail", Runtime: "claude", Model: "claude-sonnet-4-6"})]
	if got.Turns != 2 {
		t.Fatalf("Turns = %d, want 2 (the timestamped observation must win over the unparseable one)", got.Turns)
	}
}

// TestComputeJourneyDeltasExactArithmetic is AC-3's primary proof: the delta for
// each scenario/model equals (PR value - the ledger observation with the latest
// captured_at for that scenario/model) EXACTLY, across turns, the full token
// breakdown, and cost.
func TestComputeJourneyDeltasExactArithmetic(t *testing.T) {
	baseline := journeymetrics.Ledger{
		Scenarios: []journeymetrics.ScenarioLedgerEntry{
			{
				ScenarioID: "shallow-boot-window",
				Observations: []journeymetrics.Record{
					{
						ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
						Turns: 18, Tokens: journeymetrics.TokenTotals{Input: 100, Output: 50, CacheCreation: 1562, CacheRead: 20000, Total: 21712},
						TotalCostUSD: 1.2500, CapturedAt: "2026-06-20T00:00:00Z",
						RunURL: "https://github.com/spacedock-dev/spacedock/actions/runs/27931963802",
					},
					// An OLDER observation for the same scenario/model — must be
					// ignored in favor of the one above.
					{
						ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
						Turns: 5, Tokens: journeymetrics.TokenTotals{Total: 500},
						TotalCostUSD: 0.05, CapturedAt: "2026-05-01T00:00:00Z",
					},
				},
			},
		},
	}
	current := []journeymetrics.Record{
		{
			ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
			Turns: 20, Tokens: journeymetrics.TokenTotals{Input: 110, Output: 60, CacheCreation: 1576, CacheRead: 21000, Total: 22746},
			TotalCostUSD: 1.4000,
		},
	}

	deltas := ComputeJourneyDeltas(baseline, current)
	if len(deltas) != 1 {
		t.Fatalf("deltas = %d, want 1", len(deltas))
	}
	d := deltas[0]
	if !d.HasBaseline {
		t.Fatal("expected a matching baseline observation")
	}
	if d.TurnsDelta != 2 {
		t.Fatalf("TurnsDelta = %d, want 2 (20-18)", d.TurnsDelta)
	}
	wantTokens := journeymetrics.TokenTotals{Input: 10, Output: 10, CacheCreation: 14, CacheRead: 1000, Total: 1034}
	if d.TokensDelta != wantTokens {
		t.Fatalf("TokensDelta = %+v, want %+v", d.TokensDelta, wantTokens)
	}
	if diff := d.CostDeltaUSD - 0.15; diff > 1e-9 || diff < -1e-9 {
		t.Fatalf("CostDeltaUSD = %v, want 0.15 (1.40-1.25)", d.CostDeltaUSD)
	}
	if d.BaselineRunURL != "https://github.com/spacedock-dev/spacedock/actions/runs/27931963802" {
		t.Fatalf("BaselineRunURL = %q, want the latest-by-captured_at observation's run_url", d.BaselineRunURL)
	}
}

// TestComputeJourneyDeltasReportsScenarioWithNoBaseline proves a scenario/model
// the PR run produced but the baseline ledger never saw is still reported
// (against a zero-value baseline) rather than silently dropped.
func TestComputeJourneyDeltasReportsScenarioWithNoBaseline(t *testing.T) {
	baseline := journeymetrics.Ledger{}
	current := []journeymetrics.Record{
		{ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6", Turns: 20, Tokens: journeymetrics.TokenTotals{Total: 100}},
	}
	deltas := ComputeJourneyDeltas(baseline, current)
	if len(deltas) != 1 {
		t.Fatalf("deltas = %d, want 1", len(deltas))
	}
	if deltas[0].HasBaseline {
		t.Fatal("expected HasBaseline=false for a scenario/model absent from the baseline ledger")
	}
	if deltas[0].TurnsDelta != 20 {
		t.Fatalf("TurnsDelta = %d, want 20 (delta against a zero baseline equals the PR's own value)", deltas[0].TurnsDelta)
	}
}

// TestRenderJourneyDeltaCommentIncludesExactDeltasAndMarker proves the rendered
// comment BODY (not just the underlying struct) carries the marker used to
// identify the comment's origin, and the exact computed delta values.
func TestRenderJourneyDeltaCommentIncludesExactDeltasAndMarker(t *testing.T) {
	deltas := []JourneyDelta{
		{
			ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
			HasBaseline: true, BaselineRunURL: "https://github.com/spacedock-dev/spacedock/actions/runs/27931963802",
			TurnsDelta: 2, TokensDelta: journeymetrics.TokenTotals{Total: 1034}, CostDeltaUSD: 0.15,
		},
	}
	body := RenderJourneyDeltaComment(deltas)
	if !strings.HasPrefix(body, JourneyDeltaCommentMarker) {
		t.Fatalf("comment body does not start with the sticky marker:\n%s", body)
	}
	for _, want := range []string{"shallow-boot-window", "+2", "+1034", "+0.1500", "27931963802"} {
		if !strings.Contains(body, want) {
			t.Fatalf("comment body missing %q:\n%s", want, body)
		}
	}
}

// TestRenderJourneyDeltaCommentShowsTokenClassBreakdown proves the comment
// renders Cache Read Δ and Cache Creation Δ as their OWN columns, in the
// correct (non-swapped) position — not collapsed into the Tokens Δ (total)
// figure. cache_read and cache_creation differ ~12x in cost and meaning for a
// boot metric, so folding them into one number hides the signal this Minor
// exists to surface. Uses distinguishable, easily-confused-if-swapped values
// (CacheRead=111, CacheCreation=222) and checks the actual table row's cell
// positions, not just substring presence anywhere in the body (which would
// not catch a column swap).
func TestRenderJourneyDeltaCommentShowsTokenClassBreakdown(t *testing.T) {
	deltas := []JourneyDelta{
		{
			ScenarioID: "shallow-boot-window", Runtime: "claude", Model: "claude-sonnet-4-6",
			HasBaseline: true, BaselineRunURL: "https://github.com/spacedock-dev/spacedock/actions/runs/27931963802",
			TurnsDelta: 2, TokensDelta: journeymetrics.TokenTotals{CacheRead: 111, CacheCreation: 222, Total: 1034}, CostDeltaUSD: 0.15,
		},
	}
	body := RenderJourneyDeltaComment(deltas)

	var row string
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "| shallow-boot-window") {
			row = line
			break
		}
	}
	if row == "" {
		t.Fatalf("comment body missing the data row:\n%s", body)
	}
	cells := strings.Split(row, "|")
	// | (empty) | Scenario | Runtime | Model | Turns Δ | Cache Read Δ | Cache Creation Δ | Tokens Δ (total) | Cost Δ (USD) | Baseline | (empty) |
	if len(cells) != 11 {
		t.Fatalf("data row has %d cells, want 11 (header shape changed?): %q", len(cells), row)
	}
	cacheReadCell := strings.TrimSpace(cells[5])
	cacheCreationCell := strings.TrimSpace(cells[6])
	tokensTotalCell := strings.TrimSpace(cells[7])
	if cacheReadCell != "+111" {
		t.Fatalf("Cache Read Δ cell = %q, want +111 (got %q for Cache Creation Δ) — columns may be swapped", cacheReadCell, cacheCreationCell)
	}
	if cacheCreationCell != "+222" {
		t.Fatalf("Cache Creation Δ cell = %q, want +222 (got %q for Cache Read Δ) — columns may be swapped", cacheCreationCell, cacheReadCell)
	}
	if tokensTotalCell != "+1034" {
		t.Fatalf("Tokens Δ (total) cell = %q, want +1034", tokensTotalCell)
	}
}

// TestRenderJourneyDeltaCommentRendersNoBaselineAsNewNotSelfDelta proves a
// scenario/model with no matching baseline observation renders "n/a (new)" in
// its delta cells rather than a self-delta against an implicit zero baseline
// (which would otherwise print the observation's own full value as if it were
// a huge, meaningless increase).
func TestRenderJourneyDeltaCommentRendersNoBaselineAsNewNotSelfDelta(t *testing.T) {
	deltas := []JourneyDelta{
		{
			ScenarioID: "brand-new-scenario", Runtime: "claude", Model: "claude-sonnet-4-6",
			HasBaseline: false, TurnsDelta: 20, TokensDelta: journeymetrics.TokenTotals{Total: 5000}, CostDeltaUSD: 2.5,
		},
	}
	body := RenderJourneyDeltaComment(deltas)
	if !strings.Contains(body, "brand-new-scenario") {
		t.Fatalf("comment body missing the scenario row:\n%s", body)
	}
	if strings.Contains(body, "+20") || strings.Contains(body, "+5000") || strings.Contains(body, "+2.5000") {
		t.Fatalf("comment body rendered a self-delta for a no-baseline row instead of n/a (new):\n%s", body)
	}
	if !strings.Contains(body, "n/a (new)") {
		t.Fatalf("comment body missing n/a (new) for the no-baseline row:\n%s", body)
	}
}

// TestJourneyDeltaCreateCommentArgs proves the gh argv used to post a brand new
// comment (no prior journey-delta comment found on the PR).
func TestJourneyDeltaCreateCommentArgs(t *testing.T) {
	args := JourneyDeltaCreateCommentArgs("42", "/tmp/comment.md")
	joined := strings.Join(args, " ")
	for _, want := range []string{"pr comment", "42", "--body-file /tmp/comment.md"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("gh argv %v missing %q", args, want)
		}
	}
	if strings.Contains(joined, "--edit-last") {
		t.Fatalf("gh argv %v must not carry --edit-last (marker-based lookup already decided this is a fresh comment)", args)
	}
}

// TestJourneyDeltaUpdateCommentArgs proves the gh argv used to PATCH the exact
// existing comment found by marker — the --edit-last replacement.
func TestJourneyDeltaUpdateCommentArgs(t *testing.T) {
	args := JourneyDeltaUpdateCommentArgs("987654321", "/tmp/comment.md")
	joined := strings.Join(args, " ")
	for _, want := range []string{"api", "issues/comments/987654321", "-X PATCH", "body=@/tmp/comment.md"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("gh argv %v missing %q", args, want)
		}
	}
}
