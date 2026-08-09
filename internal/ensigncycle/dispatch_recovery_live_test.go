//go:build live

// ABOUTME: Live proof for AC-2 (degraded-bare) and AC-3 (break-glass-shim) — the
// ABOUTME: two fo-dispatch-recovery scenarios, intentionally Claude-only.
package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// writeDispatchRecoveryWorkflow seeds the shared one-entity fixture both
// dispatch-recovery scenarios drive.
//
//spacedock:live-fixture id=dispatch-recovery/base
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
//
//spacedock:live-fixture id=dispatch-recovery/failing-build
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

// TestLiveBareReachable is AC-2's behavioral proof (post-retirement): a plain
// `/spacedock bare` instruction (riding in the initial `-p` prompt) must produce
// bare-shaped Agent() calls (no `name`, no `run_in_background`) for the run, WITHOUT
// the retired Degraded Mode captain report and WITHOUT a
// Skill(skill="spacedock:fo-dispatch-recovery") load. A drive still emitting either is
// now a FAILURE.
// Run it against a real credential:
// `go test -tags live -run TestLiveBareReachable ./internal/ensigncycle -v -count=1`.
//
//spacedock:live-proof id=claude-bare-dispatch lane=claude-live
func TestLiveBareReachable(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	workflowRoot := t.TempDir()
	writeDispatchRecoveryWorkflow(t, workflowRoot)

	scenario := sharedRuntimeScenario{name: "bare-reachable"}
	result := runner.run(t, scenario, workflowRoot, bareReachablePrompt())
	if err := assertBareReachableObservables(result.stream); err != nil {
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
//
//spacedock:live-proof id=claude-dispatch-build-break-glass lane=claude-live
func TestLiveBreakGlassShimRecovery(t *testing.T) {
	role, err := claudeLiveRole(envOr("SPACEDOCK_LIVE_MODEL", "sonnet"))
	if err != nil {
		t.Fatal(err)
	}
	if role == "claude-sonnet" {
		t.Skipf("TODO(824ecawn5jttbykcgx82nbf4): %s/claude-dispatch-build-break-glass lacks passing live evidence", role)
	}
	runner := newClaudeLiveRunner(t)
	workflowRoot := t.TempDir()
	writeDispatchRecoveryWorkflow(t, workflowRoot)

	shimDir := writeStubBreakGlassSpacedock(t, runner.binary)
	runner.env = withSpacedockShimShellEnv(t, runner.env, shimDir)
	scenarioRunner := runner.withStubPATH(shimDir)

	scenario := sharedRuntimeScenario{name: "break-glass-shim"}
	result := scenarioRunner.run(t, scenario, workflowRoot, breakGlassShimPrompt())
	if err := assertBreakGlassObservables(result.stream); err != nil {
		t.Fatalf("%v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	emitClaudeScenarioMetrics(t, scenario, result, runner.model())
}
