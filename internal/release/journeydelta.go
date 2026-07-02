// ABOUTME: Per-PR journey-cost delta — selects the latest-by-captured_at baseline
// ABOUTME: observation per scenario/model from a published ledger and diffs it against a PR run.
package release

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// journeyDeltaKey identifies one scenario/runtime/model line the PR delta
// comment reports.
type journeyDeltaKey struct {
	ScenarioID string
	Runtime    string
	Model      string
}

// LatestObservations selects, for every distinct scenario/runtime/model in a
// published ledger, the single observation with the LATEST captured_at
// timestamp — the AC-3 delta baseline. A published ledger holds N observations
// per scenario/model once AC-2's widened aggregation ships, so "the baseline"
// is ambiguous without this reduction; latest-by-captured_at (not mean, not
// max, not array position) is the one this entity pins. An observation with an
// empty or unparseable captured_at sorts as the oldest possible time, so any
// genuinely-timestamped peer always wins over it.
func LatestObservations(ledger journeymetrics.Ledger) map[string]journeymetrics.Record {
	latest := map[string]journeymetrics.Record{}
	latestTime := map[string]time.Time{}
	for _, scenario := range ledger.Scenarios {
		for _, obs := range scenario.Observations {
			key := deltaKeyString(journeyDeltaKey{ScenarioID: obs.ScenarioID, Runtime: obs.Runtime, Model: obs.Model})
			t := parseCapturedAt(obs.CapturedAt)
			if cur, ok := latestTime[key]; !ok || t.After(cur) {
				latest[key] = obs
				latestTime[key] = t
			}
		}
	}
	return latest
}

func parseCapturedAt(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func deltaKeyString(k journeyDeltaKey) string {
	return strings.Join([]string{k.ScenarioID, k.Runtime, k.Model}, "\x00")
}

// JourneyDelta is one scenario/runtime/model line of the AC-3 PR comment: a
// PR run's fresh observation against the previously published release's
// latest-by-captured_at baseline for the same scenario/runtime/model.
type JourneyDelta struct {
	ScenarioID     string
	Runtime        string
	Model          string
	HasBaseline    bool
	BaselineRunURL string
	TurnsDelta     int
	TokensDelta    journeymetrics.TokenTotals
	CostDeltaUSD   float64
}

// ComputeJourneyDeltas computes, for every scenario/runtime/model the current
// PR run produced, the delta against the baseline ledger's latest-by-
// captured_at observation for that same scenario/runtime/model (see
// LatestObservations): delta = PR value - baseline value, exactly, across
// turns, the full token breakdown, and cost — whatever the ledger already
// tracks generically, not a boot-window special case. A scenario/model with no
// matching baseline observation is still reported (HasBaseline=false, delta
// equal to the PR's own value) rather than silently dropped, so a brand-new
// scenario doesn't vanish from the comment.
func ComputeJourneyDeltas(baseline journeymetrics.Ledger, current []journeymetrics.Record) []JourneyDelta {
	latest := LatestObservations(baseline)
	deltas := make([]JourneyDelta, 0, len(current))
	for _, cur := range current {
		key := deltaKeyString(journeyDeltaKey{ScenarioID: cur.ScenarioID, Runtime: cur.Runtime, Model: cur.Model})
		base, ok := latest[key]
		d := JourneyDelta{
			ScenarioID:  cur.ScenarioID,
			Runtime:     cur.Runtime,
			Model:       cur.Model,
			HasBaseline: ok,
		}
		if ok {
			d.BaselineRunURL = base.RunURL
		}
		d.TurnsDelta = cur.Turns - base.Turns
		d.TokensDelta = journeymetrics.TokenTotals{
			Input:         cur.Tokens.Input - base.Tokens.Input,
			Output:        cur.Tokens.Output - base.Tokens.Output,
			CacheCreation: cur.Tokens.CacheCreation - base.Tokens.CacheCreation,
			CacheRead:     cur.Tokens.CacheRead - base.Tokens.CacheRead,
			Total:         cur.Tokens.Total - base.Tokens.Total,
		}
		d.CostDeltaUSD = cur.TotalCostUSD - base.TotalCostUSD
		deltas = append(deltas, d)
	}
	sort.Slice(deltas, func(i, j int) bool {
		if deltas[i].ScenarioID != deltas[j].ScenarioID {
			return deltas[i].ScenarioID < deltas[j].ScenarioID
		}
		if deltas[i].Runtime != deltas[j].Runtime {
			return deltas[i].Runtime < deltas[j].Runtime
		}
		return deltas[i].Model < deltas[j].Model
	})
	return deltas
}

// journeyDeltaCommentMarker is stamped into every rendered comment. It has no
// behavioral role in THIS process — the sticky-update behavior comes from `gh
// pr comment --edit-last` editing the poster's own last comment — but it lets a
// human (or a future automated check) identify the comment's origin at a
// glance.
const journeyDeltaCommentMarker = "<!-- spacedock:journey-delta -->"

// RenderJourneyDeltaComment renders the AC-3 PR comment body: one table row per
// scenario/runtime/model, each cell an exact (PR value - baseline value) delta.
func RenderJourneyDeltaComment(deltas []JourneyDelta) string {
	var b strings.Builder
	b.WriteString(journeyDeltaCommentMarker)
	b.WriteString("\n### Journey cost delta\n\n")
	if len(deltas) == 0 {
		b.WriteString("No journey metrics observations were produced by this run.\n")
		return b.String()
	}
	b.WriteString("| Scenario | Runtime | Model | Turns Δ | Tokens Δ | Cost Δ (USD) | Baseline |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	for _, d := range deltas {
		baseline := "none (new observation)"
		switch {
		case d.HasBaseline && d.BaselineRunURL != "":
			baseline = fmt.Sprintf("[latest published](%s)", d.BaselineRunURL)
		case d.HasBaseline:
			baseline = "latest published"
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %+d | %+d | %+.4f | %s |\n",
			d.ScenarioID, d.Runtime, d.Model, d.TurnsDelta, d.TokensDelta.Total, d.CostDeltaUSD, baseline)
	}
	return b.String()
}

// JourneyDeltaCommentArgs returns the `gh pr comment` argv the AC-3 step
// invokes to post (or update) the PR delta comment. --edit-last edits the
// poster's own last comment on the PR instead of appending a new one on every
// push; --create-if-none falls back to creating the first comment when none
// exists yet, so the FIRST post on a PR still lands.
func JourneyDeltaCommentArgs(prNumber, bodyFile string) []string {
	return []string{"pr", "comment", prNumber, "--body-file", bodyFile, "--edit-last", "--create-if-none"}
}
