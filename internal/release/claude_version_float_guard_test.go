// ABOUTME: Structural guard binding runtime-live-e2e.yml's Claude Code install to float
// ABOUTME: to installer-latest with the workflow_dispatch override preserved and no pin re-introduced.
package release

import (
	"fmt"
	"strings"
	"testing"
)

// TestClaudeLiveWorkflowFloatsToLatest locks the float policy structurally over
// the parsed workflow: the empty-input install path resolves to installer-`latest`
// (the merged-team floor #396 ships on — what users actually run), a maintainer's
// workflow_dispatch input still overrides it, and the retired 2.1.177 pin var is
// gone, not hidden. AC-1 (the resolved version is the merged floor and the merged
// scenarios green) is the LIVE proof on the Show tool versions step, not checkable
// offline.
func TestClaudeLiveWorkflowFloatsToLatest(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	if err := assertClaudeLiveWorkflowFloatsToLatest(live); err != nil {
		t.Fatal(err)
	}
}

// TestClaudeVersionFloatGuardRejectsReintroducedPin is the adversarial twin (the
// real discriminator): mutate a SPACEDOCK_PINNED_CLAUDE_VERSION var and the
// `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` fallback back into the
// fixture and the guard must RED. A guard that stays green here is the exact hole
// this guard exists to close — a future PR silently re-pins the float lane and CI
// drifts back onto a frozen version, the #395 defect this task retired.
func TestClaudeVersionFloatGuardRejectsReintroducedPin(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	if err := assertClaudeLiveWorkflowFloatsToLatest(live); err != nil {
		t.Fatalf("real runtime-live-e2e.yml unexpectedly fails the float guard before mutation: %v", err)
	}

	adversarial := strings.Replace(live,
		`      DISABLE_AUTOUPDATER: "1"
      ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}`,
		`      DISABLE_AUTOUPDATER: "1"
      SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"
      ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}`,
		1)
	adversarial = strings.Replace(adversarial,
		`VERSION="${CLAUDE_VERSION:-latest}"`,
		`VERSION="${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}"`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing the float env block / latest resolution to mutate into a pin")
	}

	if err := assertClaudeLiveWorkflowFloatsToLatest(adversarial); err == nil {
		t.Fatal("float guard accepted an install step that re-pins via SPACEDOCK_PINNED_CLAUDE_VERSION")
	}
}

// assertClaudeLiveWorkflowFloatsToLatest is the shared predicate behind the float
// guard. It checks, over the parsed workflow:
//   - the install resolves `${CLAUDE_VERSION:-latest}` (empty input ⇒ installer
//     latest, the merged-team floor users run);
//   - the install step still reads inputs.claude_version into CLAUDE_VERSION, so a
//     non-empty workflow_dispatch input overrides the float (the escape hatch);
//   - the workflow declares NO SPACEDOCK_PINNED_CLAUDE_VERSION var anywhere — the
//     retired pin is gone, not hidden behind the float default.
func assertClaudeLiveWorkflowFloatsToLatest(workflow string) error {
	// The retired pin must be absent everywhere — not merely unreferenced from the
	// resolution. A lingering var declaration is a re-pin waiting to be wired back.
	if strings.Contains(workflow, "SPACEDOCK_PINNED_CLAUDE_VERSION") {
		return fmt.Errorf("runtime-live-e2e.yml still declares SPACEDOCK_PINNED_CLAUDE_VERSION — the retired 2.1.177 pin must be gone, not hidden")
	}

	step, ok := findStepByName(parseWorkflowSteps(workflow), "Install Claude Code")
	if !ok {
		return fmt.Errorf("runtime-live-e2e.yml has no Install Claude Code step")
	}

	// The override escape hatch: the step reads the workflow_dispatch input into
	// CLAUDE_VERSION (the `env:` mapping is parsed off the raw text since
	// parseWorkflowSteps drops `with:`/`env:` keys).
	if !stepEnvBinds(workflow, step, "CLAUDE_VERSION", "${{ inputs.claude_version }}") {
		return fmt.Errorf("runtime-live-e2e.yml Install Claude Code step does not read inputs.claude_version into CLAUDE_VERSION")
	}

	commands := executableShellCommands(step.run)

	// The install resolves the float via the override-precedence expression.
	// `${CLAUDE_VERSION:-latest}` means a non-empty input wins (the escape hatch)
	// and an empty input falls back to installer-latest (the float).
	resolvesFloat := false
	for _, c := range commands {
		if strings.Contains(c, "${CLAUDE_VERSION:-latest}") {
			resolvesFloat = true
		}
	}
	if !resolvesFloat {
		return fmt.Errorf("runtime-live-e2e.yml Install Claude Code step does not resolve ${CLAUDE_VERSION:-latest} (the override-precedence float fallback)")
	}

	return nil
}

// findStepByName returns the parsed step with the given name.
func findStepByName(steps []workflowStep, name string) (workflowStep, bool) {
	for _, s := range steps {
		if s.name == name {
			return s, true
		}
	}
	return workflowStep{}, false
}

// stepEnvBinds reports whether the named step's `env:` block maps key to value.
// parseWorkflowSteps does not capture step-level `env:`, so this re-walks the raw
// workflow text between the step's `- name:` header and the next step header,
// looking for the `key: value` line inside an `env:` mapping.
func stepEnvBinds(workflow string, step workflowStep, key, value string) bool {
	lines := strings.Split(workflow, "\n")
	start := -1
	for i, raw := range lines {
		if strings.TrimSpace(raw) == "- name: "+step.name {
			start = i
			break
		}
	}
	if start < 0 {
		return false
	}
	want := key + ": " + value
	inEnv := false
	for i := start + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "- name: ") {
			break
		}
		if trimmed == "env:" {
			inEnv = true
			continue
		}
		if inEnv && trimmed == want {
			return true
		}
	}
	return false
}
