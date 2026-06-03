// ABOUTME: P1 cross-equality — the grade's terminalTeardownMarker const is byte-identical to the
// ABOUTME: marker the FO contract mandates, so a grade-const+fixtures drift that skips the contract is caught.
package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGradeMarkerMatchesContract pins the grade's terminalTeardownMarker const to
// the marker the FO contract mandates VERBATIM. The marker lives as three
// independent verbatim copies — this grade const, the integration lint's const,
// and the contract prose — with no single source of truth. Without this check the
// equality is only TRANSITIVELY held (the integration lint asserts ITS const is
// in the contract; the grade asserts ITS const greens the fixtures), so a
// consistent grade-const + fixtures drift that skips the contract would slip
// through: the offline grade would green a drifted marker the live FO never emits
// because the contract still mandates the old string. Asserting the grade const
// appears VERBATIM in both contract files ties grade↔contract directly; the
// integration lint ties integration-const↔contract; together the three copies are
// pinned to one prose source.
//
// The contract files are read via the in-repo layout path (the ensigncycle
// package sits at internal/ensigncycle, so the FO references are two dirs up).
// This is a fixed repo-relative path, not a machine-specific dependency.
func TestGradeMarkerMatchesContract(t *testing.T) {
	contractFiles := []string{
		filepath.Join("..", "..", "skills", "first-officer", "references", "first-officer-shared-core.md"),
		filepath.Join("..", "..", "skills", "first-officer", "references", "claude-first-officer-runtime.md"),
	}
	for _, f := range contractFiles {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read contract file %s: %v", f, err)
		}
		if !strings.Contains(string(b), terminalTeardownMarker) {
			t.Errorf("%s does not carry the grade's terminalTeardownMarker verbatim — the grade const and the contract have drifted; the offline grade would green a marker the live FO never emits", f)
		}
	}
}
