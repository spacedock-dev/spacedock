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

// JourneyDeltaCommentMarker is stamped into every rendered comment. A caller
// finds its own prior comment by searching the PR's comments for one whose
// body starts with this marker (JourneyDeltaUpdateCommentArgs), rather than by
// editing "the poster's last comment" (`--edit-last`) — --edit-last targets the
// wrong comment if any OTHER automated comment from the same bot account lands
// on the PR in between.
const JourneyDeltaCommentMarker = "<!-- spacedock:journey-delta -->"

// RenderJourneyDeltaComment renders the AC-3 PR comment body: one table row per
// scenario/runtime/model. A row with a baseline shows the exact (PR value -
// baseline value) delta for turns, each token class, and cost; a row with NO
// baseline (a brand-new scenario/model this ledger has never seen) renders
// "n/a (new)" in every delta cell instead of a self-delta against an implicit
// zero baseline, which would otherwise read as a huge, meaningless "increase."
func RenderJourneyDeltaComment(deltas []JourneyDelta) string {
	var b strings.Builder
	b.WriteString(JourneyDeltaCommentMarker)
	b.WriteString("\n### Journey cost delta\n\n")
	if len(deltas) == 0 {
		b.WriteString("No journey metrics observations were produced by this run.\n")
		return b.String()
	}
	b.WriteString("| Scenario | Runtime | Model | Turns Δ | Cache Read Δ | Cache Creation Δ | Tokens Δ (total) | Cost Δ (USD) | Baseline |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- | --- |\n")
	for _, d := range deltas {
		if !d.HasBaseline {
			fmt.Fprintf(&b, "| %s | %s | %s | n/a (new) | n/a (new) | n/a (new) | n/a (new) | n/a (new) | none (new observation) |\n",
				d.ScenarioID, d.Runtime, d.Model)
			continue
		}
		baseline := "latest published"
		if d.BaselineRunURL != "" {
			baseline = fmt.Sprintf("[latest published](%s)", d.BaselineRunURL)
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %+d | %+d | %+d | %+d | %+.4f | %s |\n",
			d.ScenarioID, d.Runtime, d.Model, d.TurnsDelta, d.TokensDelta.CacheRead, d.TokensDelta.CacheCreation, d.TokensDelta.Total, d.CostDeltaUSD, baseline)
	}
	return b.String()
}

// JourneyDeltaCreateCommentArgs returns the `gh pr comment` argv that posts a
// NEW comment — used when no prior journey-delta comment exists on the PR yet.
func JourneyDeltaCreateCommentArgs(prNumber, bodyFile string) []string {
	return []string{"pr", "comment", prNumber, "--body-file", bodyFile}
}

// JourneyDeltaUpdateCommentArgs returns the `gh api` argv that PATCHes the
// EXACT existing comment by id — the find-by-marker replacement for
// `--edit-last` (which targets "the poster's last comment on the PR," not
// necessarily this job's own prior comment).
func JourneyDeltaUpdateCommentArgs(commentID, bodyFile string) []string {
	return []string{"api", "repos/{owner}/{repo}/issues/comments/" + commentID, "-X", "PATCH", "-f", "body=@" + bodyFile}
}
