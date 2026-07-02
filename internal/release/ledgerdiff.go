// ABOUTME: Structural pre-upload safeguard for the AC-4 release-ledger backfill —
// ABOUTME: proves a rebuilt ledger equals the original plus exactly the added scenarios.
package release

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

// DiffAddedScenarios is the AC-4 pre-upload safeguard: it proves a rebuilt
// journey-cost ledger equals the backed-up original ledger's scenarios array
// PLUS EXACTLY the named added scenario ids, with every pre-existing scenario
// byte-for-byte unchanged. Rebuilding via `journey-costs` re-normalizes every
// historical record (normalizeRecord restamps schema_version and recomputes
// token totals on every record it reads), so a rebuild quirk could otherwise
// silently mutate already-published historical data that has nothing to do
// with the backfill. It returns a descriptive error on ANY deviation — a
// changed, missing, or unexpectedly-extra scenario — so the caller can abort
// BEFORE an upload instead of publishing a corrupted ledger.
func DiffAddedScenarios(original, rebuilt journeymetrics.Ledger, wantAdded ...string) error {
	originalByID := scenarioByID(original)
	rebuiltByID := scenarioByID(rebuilt)

	for id, origEntry := range originalByID {
		rebuiltEntry, ok := rebuiltByID[id]
		if !ok {
			return fmt.Errorf("scenario %q present in the original ledger is MISSING from the rebuilt ledger", id)
		}
		origJSON, err := json.Marshal(origEntry)
		if err != nil {
			return fmt.Errorf("marshal original scenario %q: %w", id, err)
		}
		rebuiltJSON, err := json.Marshal(rebuiltEntry)
		if err != nil {
			return fmt.Errorf("marshal rebuilt scenario %q: %w", id, err)
		}
		if string(origJSON) != string(rebuiltJSON) {
			return fmt.Errorf("scenario %q changed by the rebuild (a rebuild must only ADD the backfilled scenario, never alter an existing one) — original:\n%s\nrebuilt:\n%s", id, origJSON, rebuiltJSON)
		}
	}

	var extra []string
	for id := range rebuiltByID {
		if _, ok := originalByID[id]; !ok {
			extra = append(extra, id)
		}
	}
	sort.Strings(extra)
	wantSorted := append([]string(nil), wantAdded...)
	sort.Strings(wantSorted)
	if !slices.Equal(extra, wantSorted) {
		return fmt.Errorf("rebuilt ledger added scenarios %v, want exactly %v", extra, wantSorted)
	}
	return nil
}

func scenarioByID(ledger journeymetrics.Ledger) map[string]journeymetrics.ScenarioLedgerEntry {
	out := make(map[string]journeymetrics.ScenarioLedgerEntry, len(ledger.Scenarios))
	for _, entry := range ledger.Scenarios {
		out[entry.ScenarioID] = entry
	}
	return out
}
