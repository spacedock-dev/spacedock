//go:build live

// ABOUTME: Live proof for AC-2 (degraded-bare) and AC-3 (break-glass-shim) — the
// ABOUTME: two fo-dispatch-recovery scenarios, Claude-only (not sharedRuntimeScenarios).
package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// writeDispatchRecoveryWorkflow seeds the shared one-entity fixture both
// dispatch-recovery scenarios drive.
func writeDispatchRecoveryWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), dispatchRecoveryReadme())
	entityPath := filepath.Join(root, "widget-task.md")
	writeFile(t, entityPath, dispatchRecoveryEntity())
	gitInit(t, root)
	return entityPath
}

// writeStubBreakGlassSpacedock writes a `spacedock` shim that fails ONLY
// `dispatch build` (stderr + exit 1) and delegates every other subcommand —
// including the `--version` contract gate, `status --boot`, `dispatch
// context-budget`, etc. — to the REAL built binary at realBinary. It calls
// realBinary by its absolute path rather than by PATH lookup, so prepending the
// shim's dir ahead of the real binary's dir on PATH (via withStubPATH) cannot
// recurse into itself. Other live scenarios use the same withStubPATH seam for
// their scenario-local executable shims.
func writeStubBreakGlassSpacedock(t *testing.T, realBinary string) string {
	t.Helper()
	dir := t.TempDir()
	script := fmt.Sprintf("#!/bin/sh\n"+
		"if [ \"$1\" = \"dispatch\" ] && [ \"$2\" = \"build\" ]; then\n"+
		"  echo \"stub: spacedock dispatch build unavailable (break-glass-shim scenario)\" >&2\n"+
		"  exit 1\n"+
		"fi\n"+
		"exec %q \"$@\"\n", realBinary)
	path := filepath.Join(dir, "spacedock")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// TestLiveDegradedBareRecovery is AC-2's behavioral proof: the `/spacedock bare`
// trigger (riding in the initial `-p` prompt per the M4 headless-transport
// decision) must produce the verbatim captain report, a
// Skill(skill="spacedock:fo-dispatch-recovery") load, and bare-shaped Agent() calls
// (no `name`, no `run_in_background`) for the remainder of the run.
// Run it against a real credential:
// `go test -tags live -run TestLiveDegradedBareRecovery ./internal/ensigncycle -v -count=1`.
func TestLiveDegradedBareRecovery(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	workflowRoot := t.TempDir()
	writeDispatchRecoveryWorkflow(t, workflowRoot)

	scenario := sharedRuntimeScenario{
		name:          "degraded-bare",
		oldPythonTest: "n/a (net-new; dispatch-exception-paths-deferred-module)",
		intent:        "The `/spacedock bare` trigger produces the verbatim captain report, loads spacedock:fo-dispatch-recovery, and every subsequent Agent() dispatch is bare-shaped.",
	}
	result := runner.run(t, scenario, workflowRoot, degradedBarePrompt())
	if err := assertDegradedBareObservables(result.stream); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}

// TestLiveBreakGlassShimRecovery is AC-3's behavioral proof: a PATH-shimmed
// `spacedock dispatch build` that fails for real must produce a captain-facing
// helper-failure report BEFORE any Agent() call, a
// Skill(skill="spacedock:fo-dispatch-recovery") load, and a break-glass-shaped
// Agent() call (run_in_background=true, a {worker_key}-{slug}-{stage} name, a
// prompt carrying Skill(skill="spacedock:ensign") and an inline ### Stage
// definition).
// Run it against a real credential:
// `go test -tags live -run TestLiveBreakGlassShimRecovery ./internal/ensigncycle -v -count=1`.
func TestLiveBreakGlassShimRecovery(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	workflowRoot := t.TempDir()
	writeDispatchRecoveryWorkflow(t, workflowRoot)

	shimDir := writeStubBreakGlassSpacedock(t, runner.binary)
	scenarioRunner := runner.withStubPATH(shimDir)

	scenario := sharedRuntimeScenario{
		name:          "break-glass-shim",
		oldPythonTest: "n/a (net-new; dispatch-exception-paths-deferred-module)",
		intent:        "A failing `spacedock dispatch build` produces a captain-facing report before any Agent() call, loads spacedock:fo-dispatch-recovery, and manually assembles a break-glass-shaped Agent() dispatch.",
	}
	result := scenarioRunner.run(t, scenario, workflowRoot, breakGlassShimPrompt())
	if err := assertBreakGlassObservables(result.stream); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}
