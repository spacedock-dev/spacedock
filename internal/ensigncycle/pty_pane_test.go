// ABOUTME: Resolves the first-officer tmux pane by title so the pty driver reads/sends to the FO,
// ABOUTME: never a dispatched-ensign sibling pane (the spike's FO-pane-capture hazard, AC-5).
package ensigncycle

import (
	"fmt"
	"strings"
	"testing"
)

// foPaneMarker is the title fragment a first-officer pane carries. When the FO
// dispatches a teams ensign the runtime can materialize the ensign as a sibling
// tmux pane (own index, own PID) that becomes ACTIVE — so reading/sending to the
// active pane would hit the ensign, not the FO (the spike's most important
// borrowed hazard). Resolving the FO pane by title is robust to both topologies:
// it resolves to pane 0 in the single-pane case and to the FO's pane in the split
// case. Borrowed from spacedock-gym's internal/driver leadPaneMarker.
const foPaneMarker = "spacedock:first-officer"

// firstOfficerPaneIndex parses `tmux list-panes -F
// '#{pane_index}\t#{pane_active}\t#{pane_title}'` output and returns the index of
// the pane whose title marks it as the first officer. It errors when no
// first-officer pane is present, so the driver refuses to read/send rather than
// steer a sub-agent. Borrowed (not imported) from spacedock-gym's
// firstOfficerPaneIndex — a separate Go module, so the mechanism is ported, not
// cross-imported (the ideation spacedock-gym decision).
func firstOfficerPaneIndex(listOutput string) (string, error) {
	for _, line := range strings.Split(strings.TrimRight(listOutput, "\n"), "\n") {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 3)
		if len(fields) < 3 {
			continue
		}
		index, title := fields[0], fields[2]
		if strings.Contains(title, foPaneMarker) {
			return index, nil
		}
	}
	return "", fmt.Errorf("no %q pane found among tmux panes; refusing to read/send to avoid steering a dispatched sub-agent. panes:\n%s", foPaneMarker, listOutput)
}

// TestFirstOfficerPaneIndex is the AC-5 offline proof of the FO-pane-capture
// hazard fix: the driver resolves the FO pane by #{pane_title}, NOT the active
// pane. It runs under default build tags (no model spend, no tmux) over synthetic
// `tmux list-panes` output, mirroring spacedock-gym's foreground table test. The
// two-pane case is the live reality the spike surfaced: an Agent dispatch
// materialized a second pane (the ensign) that became ACTIVE, so a read of the
// active pane would scrape the ensign, not the FO's TeamDelete/marker.
func TestFirstOfficerPaneIndex(t *testing.T) {
	const (
		// Real two-pane shape: the FO at index 0, an Agent-dispatched ensign at
		// index 1 that became active (#{pane_active}=1). Resolving by title must
		// still pick the FO pane (0), not the active ensign pane (1).
		twoPaneList = "0\t0\tnode ~ spacedock:first-officer\n1\t1\tnode ~ spacedock:ensign-live-team-mode-terminal-harness\n"
		// Single-pane case (no dispatch split): only the FO pane exists.
		onePaneList = "0\t1\tnode ~ spacedock:first-officer\n"
		// Degenerate: an ensign-only listing must NOT resolve — refuse rather than
		// steer a sub-agent.
		ensignOnlyList = "0\t1\tnode ~ spacedock:ensign-live-team-mode-terminal-harness\n"
	)

	t.Run("two_panes_resolves_fo_not_active_ensign", func(t *testing.T) {
		idx, err := firstOfficerPaneIndex(twoPaneList)
		if err != nil {
			t.Fatal(err)
		}
		if idx != "0" {
			t.Errorf("firstOfficerPaneIndex(twoPaneList) = %q, want \"0\" (the FO pane, not the active ensign pane at 1)", idx)
		}
	})

	t.Run("single_pane_resolves_fo", func(t *testing.T) {
		idx, err := firstOfficerPaneIndex(onePaneList)
		if err != nil {
			t.Fatal(err)
		}
		if idx != "0" {
			t.Errorf("firstOfficerPaneIndex(onePaneList) = %q, want \"0\"", idx)
		}
	})

	t.Run("no_fo_pane_errors", func(t *testing.T) {
		if idx, err := firstOfficerPaneIndex(ensignOnlyList); err == nil {
			t.Errorf("firstOfficerPaneIndex(ensignOnlyList) = (%q,nil), want an error (no FO pane → refuse to read/send)", idx)
		}
		if idx, err := firstOfficerPaneIndex(""); err == nil {
			t.Errorf("firstOfficerPaneIndex(\"\") = (%q,nil), want an error", idx)
		}
	})
}
