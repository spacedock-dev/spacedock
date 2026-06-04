// ABOUTME: The live-scenario primitive — author a {runbook, setup, durable-
// ABOUTME: outcome} scenario, run it via a real-agent Runner, grade on outcomes.
package livescenario

import (
	"context"
	"fmt"
	"os"
)

// EntityState is the durable on-disk state of the scenario's entity at a point
// in the run — its body bytes as text. A scenario grades on the BEFORE→AFTER
// transition of this state plus the agent's observed final output, NOT on
// transcript phrasing (which is non-deterministic). Carrying the whole body
// keeps the primitive agnostic to which frontmatter field a given scenario
// asserts on.
type EntityState struct {
	Body string
}

// Scenario is the promoted live-scenario primitive: a runtime claim authored as
// the triple {runbook, setup, durable-outcome assertions}. It is host-neutral —
// the same scenario runs against any Runner (the Claude or Codex live adapter).
// Authoring one is a plain Go value, reachable from any package, not buried
// behind a build tag in one internal test package.
type Scenario struct {
	// Name identifies the scenario in failures and coverage.
	Name string
	// Runbook is the descriptive prompt a real agent runs — the instructions,
	// not a deterministic script. It is what the agent DOES; grading is on the
	// durable outcome of doing it, not on the runbook's wording.
	Runbook string
	// Setup stages the scenario's starting state into dir and returns the path
	// of the entity the assertions read. A setup error aborts the run (the
	// primitive cannot grade what it could not stage).
	Setup func(dir string) (entityPath string, err error)
	// Assert grades the durable outcome: the entity state BEFORE→AFTER the run
	// plus the agent's observed final output. A nil return is PASS; a non-nil
	// error is a graded failure naming the broken durable outcome.
	Assert func(before, after EntityState, observed string) error
}

// Runner is the launch + observed-extract seam — the only host-specific surface.
// It launches a real agent against the staged workflow dir / entity with the
// runbook and returns the agent's observed final message. The Claude and Codex
// live adapters (in internal/ensigncycle, behind //go:build live) implement it;
// offline tests inject a stub so the authoring + grade path runs with no model
// spend.
type Runner interface {
	Launch(ctx context.Context, dir, entityPath, runbook string) (observed string, err error)
}

// Run executes a scenario end to end: stage the setup in dir, capture the BEFORE
// entity state, launch the real agent via runner, capture the AFTER state and
// the observed output, then grade via the scenario's Assert. It returns the
// grade — nil for PASS, an error naming the broken durable outcome otherwise. A
// setup or launch error surfaces as a run failure, never a silent PASS.
func Run(ctx context.Context, dir string, sc Scenario, runner Runner) error {
	entityPath, err := sc.Setup(dir)
	if err != nil {
		return fmt.Errorf("scenario %q setup failed: %w", sc.Name, err)
	}
	before, err := readEntityState(entityPath)
	if err != nil {
		return fmt.Errorf("scenario %q could not read pre-run state: %w", sc.Name, err)
	}
	observed, err := runner.Launch(ctx, dir, entityPath, sc.Runbook)
	if err != nil {
		return fmt.Errorf("scenario %q launch failed: %w", sc.Name, err)
	}
	after, err := readEntityState(entityPath)
	if err != nil {
		return fmt.Errorf("scenario %q could not read post-run state: %w", sc.Name, err)
	}
	if gerr := sc.Assert(before, after, observed); gerr != nil {
		return fmt.Errorf("scenario %q graded FAIL: %w", sc.Name, gerr)
	}
	return nil
}

// readEntityState reads the entity file's body bytes into an EntityState.
func readEntityState(entityPath string) (EntityState, error) {
	data, err := os.ReadFile(entityPath)
	if err != nil {
		return EntityState{}, err
	}
	return EntityState{Body: string(data)}, nil
}
