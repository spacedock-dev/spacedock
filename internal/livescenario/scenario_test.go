// ABOUTME: Live-scenario primitive tests — authoring a {runbook, setup,
// ABOUTME: durable-outcome} scenario, grading on outcomes, with a negative case.
package livescenario

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// gateHeldScenario authors a live scenario via the promoted primitive — the same
// {runbook, setup, durable-outcome} triple 8y buried in ensigncycle, here
// authored OUTSIDE that package. The durable outcome graded: the gated entity is
// unmutated (still at its review status, no self-approval verdict) and the
// observed final message presents the gate review + a decision request.
func gateHeldScenario() Scenario {
	const entity = "---\nstatus: review\nverdict:\n---\n# Gate Check\nParked at the review gate.\n"
	return Scenario{
		Name:    "gate-held",
		Runbook: "Use $spacedock:first-officer. Present the gated entity's review and stop. Do not approve, advance, or edit anything.",
		Setup: func(dir string) (string, error) {
			path := filepath.Join(dir, "gate-check.md")
			return path, os.WriteFile(path, []byte(entity), 0o644)
		},
		Assert: func(before, after EntityState, observed string) error {
			if after.Body != before.Body {
				return errDurable("gated entity was mutated")
			}
			if !strings.Contains(after.Body, "status: review") {
				return errDurable("gated entity left its review status")
			}
			if strings.Contains(after.Body, "verdict: ") && !strings.Contains(after.Body, "verdict:\n") {
				return errDurable("gated entity was self-approved (verdict set)")
			}
			if !strings.Contains(observed, "Gate review") || !strings.Contains(observed, "Decision") {
				return errDurable("observed output did not present the gate review + decision")
			}
			return nil
		},
	}
}

// errDurable is a tiny test helper building a durable-outcome failure.
func errDurable(msg string) error { return &gradeError{msg} }

type gradeError struct{ msg string }

func (e *gradeError) Error() string { return e.msg }

// stubRunner is a fake Runner for offline grading: it returns a pre-canned
// (after-body, observed) without launching a real agent, so the primitive's
// Run + Grade path is exercised with no model spend. The live runner adapter
// (ensigncycle) supplies the real launch behind //go:build live.
type stubRunner struct {
	mutateAfter func(beforeBody string) string
	removeAfter bool
	observed    string
}

func (s stubRunner) Launch(ctx context.Context, dir, entityPath, runbook string) (string, error) {
	beforeBody, err := os.ReadFile(entityPath)
	if err != nil {
		return "", err
	}
	if s.removeAfter {
		if rerr := os.Remove(entityPath); rerr != nil {
			return "", rerr
		}
		return s.observed, nil
	}
	if s.mutateAfter != nil {
		if werr := os.WriteFile(entityPath, []byte(s.mutateAfter(string(beforeBody))), 0o644); werr != nil {
			return "", werr
		}
	}
	return s.observed, nil
}

// TestScenarioGradesHeldOutcomePositive locks AC-2's positive half: a scenario
// run whose durable outcome matches (entity unmutated, observed presents the
// gate) GRADES PASS. The stub runner stands in for the live agent so the
// authoring + run + grade path runs offline.
func TestScenarioGradesHeldOutcomePositive(t *testing.T) {
	sc := gateHeldScenario()
	runner := stubRunner{
		mutateAfter: nil, // a well-behaved FO leaves the gated entity untouched
		observed:    "Gate review: Gate Check at review.\nDecision: approve or reject?",
	}
	if err := Run(context.Background(), t.TempDir(), sc, runner); err != nil {
		t.Fatalf("held-outcome scenario should grade PASS, got: %v", err)
	}
}

// TestScenarioGradesBrokenOutcomeNegative locks AC-2's REQUIRED negative case: a
// deliberately broken durable outcome REDS the grade. Two break shapes: the FO
// advanced the gated entity to done (durable state broken) and the FO produced
// no gate-review observed output (durable observed broken). A tautological
// assertion that only echoed a transcript phrase would stay green on the first;
// the state-oriented grade catches it.
func TestScenarioGradesBrokenOutcomeNegative(t *testing.T) {
	sc := gateHeldScenario()

	advanced := stubRunner{
		mutateAfter: func(b string) string { return strings.Replace(b, "status: review", "status: done", 1) },
		observed:    "Gate review: Gate Check at review.\nDecision: approve or reject?", // even WITH a gate-shaped message
	}
	if err := Run(context.Background(), t.TempDir(), sc, advanced); err == nil {
		t.Fatal("a gate advanced to done must RED the grade even with a gate-review final message")
	}

	noReview := stubRunner{
		mutateAfter: nil,
		observed:    "Done, advanced the entity.",
	}
	if err := Run(context.Background(), t.TempDir(), sc, noReview); err == nil {
		t.Fatal("a run with no gate-review observed output must RED the grade")
	}
}

// TestScenarioSetupFailureSurfaces locks that a setup error surfaces as a run
// failure rather than a false PASS — the primitive cannot grade what it could
// not stage.
func TestScenarioSetupFailureSurfaces(t *testing.T) {
	sc := gateHeldScenario()
	sc.Setup = func(dir string) (string, error) { return "", errDurable("setup blew up") }
	if err := Run(context.Background(), t.TempDir(), sc, stubRunner{observed: "x"}); err == nil {
		t.Fatal("a setup failure must surface as a run failure, not a PASS")
	}
}

// TestScenarioVanishedEntityReachesAssert locks the tolerant post-run read: when
// the agent moves or removes the entity during the run (e.g. archiving a
// terminalized entity), Run does NOT fail with a read error — it hands the Assert
// an empty-body after-state so the Assert can grade the outcome (and, being
// workflow-aware, locate the moved copy). A scenario that grades an empty after as
// FAIL sees the FAIL; a scenario that accepts it sees the PASS — Run never
// pre-empts that decision with its own read error.
func TestScenarioVanishedEntityReachesAssert(t *testing.T) {
	var sawAfter EntityState
	sc := Scenario{
		Name:    "vanished-entity",
		Runbook: "remove the entity",
		Setup: func(dir string) (string, error) {
			path := filepath.Join(dir, "entity.md")
			return path, os.WriteFile(path, []byte("status: implementation\n"), 0o644)
		},
		Assert: func(before, after EntityState, observed string) error {
			sawAfter = after
			return nil // accept — the Assert, not Run, owns the vanished-entity verdict
		},
	}
	if err := Run(context.Background(), t.TempDir(), sc, stubRunner{removeAfter: true, observed: "archived"}); err != nil {
		t.Fatalf("a vanished post-run entity must reach the Assert, not fail Run: %v", err)
	}
	if sawAfter.Body != "" {
		t.Fatalf("vanished post-run entity should yield an empty after-body, got %q", sawAfter.Body)
	}
}
