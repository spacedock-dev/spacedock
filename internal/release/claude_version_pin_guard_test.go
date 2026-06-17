// ABOUTME: Structural guard binding runtime-live-e2e.yml's Claude Code install to a
// ABOUTME: single pinned, self-documented version with the workflow_dispatch override preserved.
package release

import (
	"fmt"
	"strings"
	"testing"
)

// TestClaudeLiveWorkflowPinsClaudeVersion locks AC-1/AC-3/AC-4 structurally over
// the parsed workflow: the empty-input install path resolves to the pinned
// SPACEDOCK_PINNED_CLAUDE_VERSION (never to no-pin installer-latest, which is
// 2.1.178+ with the native team tools removed), the pin carries a rationale
// comment, and a maintainer's workflow_dispatch input still overrides it. AC-2
// (the resolved version equals the pin) is the LIVE proof on the Show tool
// versions step, not checkable offline.
func TestClaudeLiveWorkflowPinsClaudeVersion(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	if err := assertClaudeLiveWorkflowPinsClaudeVersion(live); err != nil {
		t.Fatal(err)
	}
}

// TestClaudeVersionPinGuardRejectsNoPinInstall is the adversarial twin for the
// no-float half (AC-1): restore the old unconditional no-pin `install.sh | bash`
// fallback branch and the guard must RED. A guard that stays green here is the
// exact hole this task closes — a PR live run takes the no-pin branch
// unconditionally (workflow_dispatch inputs are absent for pull_request) and
// silently floats to a team-tool-broken Claude Code.
func TestClaudeVersionPinGuardRejectsNoPinInstall(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	if err := assertClaudeLiveWorkflowPinsClaudeVersion(live); err != nil {
		t.Fatalf("real runtime-live-e2e.yml unexpectedly fails the pin guard before mutation: %v", err)
	}

	adversarial := strings.Replace(live,
		`          VERSION="${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}"
          echo "Installing Claude Code version: $VERSION"
          curl -fsSL https://claude.ai/install.sh | bash -s -- "$VERSION"`,
		`          if [ -n "$CLAUDE_VERSION" ]; then
            curl -fsSL https://claude.ai/install.sh | bash -s -- "$CLAUDE_VERSION"
          else
            curl -fsSL https://claude.ai/install.sh | bash
          fi`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing the pinned-fallback install lines to mutate")
	}

	if err := assertClaudeLiveWorkflowPinsClaudeVersion(adversarial); err == nil {
		t.Fatal("pin guard accepted an install step with an unconditional no-pin `install.sh | bash` fallback branch")
	}
}

// TestClaudeVersionPinGuardRejectsUndocumentedPin is the adversarial twin for the
// self-documenting half (AC-3): strip the rationale comment above the pin and the
// guard must RED. The comment is what stops a future reader treating 2.1.177 as a
// mystery magic number and ripping it out.
func TestClaudeVersionPinGuardRejectsUndocumentedPin(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	if err := assertClaudeLiveWorkflowPinsClaudeVersion(live); err != nil {
		t.Fatalf("real runtime-live-e2e.yml unexpectedly fails the pin guard before mutation: %v", err)
	}

	// Remove every comment line that sits in the env: block above the pin var.
	var kept []string
	for _, raw := range strings.Split(live, "\n") {
		trimmed := strings.TrimSpace(raw)
		if strings.HasPrefix(trimmed, "#") && (strings.Contains(trimmed, "team tool") ||
			strings.Contains(trimmed, "2.1.17") || strings.Contains(trimmed, "DISABLE_AUTOUPDATER") ||
			strings.Contains(trimmed, "Pinned Claude Code")) {
			continue
		}
		kept = append(kept, raw)
	}
	adversarial := strings.Join(kept, "\n")
	if adversarial == live {
		t.Fatal("fixture workflow missing the pin rationale comment lines to strip")
	}

	if err := assertClaudeLiveWorkflowPinsClaudeVersion(adversarial); err == nil {
		t.Fatal("pin guard accepted a SPACEDOCK_PINNED_CLAUDE_VERSION with no preceding rationale comment")
	}
}

// TestClaudeVersionPinGuardRejectsDroppedOverride is the adversarial twin for the
// manual-override half (AC-4): drop the `${CLAUDE_VERSION:-…}` fallback so the
// install hard-codes the pin and ignores a maintainer's workflow_dispatch input.
// The guard must RED — the deliberate-float escape hatch must survive.
func TestClaudeVersionPinGuardRejectsDroppedOverride(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	if err := assertClaudeLiveWorkflowPinsClaudeVersion(live); err != nil {
		t.Fatalf("real runtime-live-e2e.yml unexpectedly fails the pin guard before mutation: %v", err)
	}

	adversarial := strings.Replace(live,
		`VERSION="${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}"`,
		`VERSION="$SPACEDOCK_PINNED_CLAUDE_VERSION"`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing the `${CLAUDE_VERSION:-…}` override-precedence expression to drop")
	}

	if err := assertClaudeLiveWorkflowPinsClaudeVersion(adversarial); err == nil {
		t.Fatal("pin guard accepted an install step that ignores the claude_version workflow_dispatch override")
	}
}

// assertClaudeLiveWorkflowPinsClaudeVersion is the shared predicate behind the
// pin guard. It checks, over the parsed workflow:
//   - AC-1: the pin var is declared in the job env: and the install step
//     references it with NO unconditional no-pin `install.sh | bash` branch;
//   - AC-3: the pin var line is preceded by a rationale comment naming the
//     team-tool regression (so it is not a mystery magic number);
//   - AC-4: the install resolves `${CLAUDE_VERSION:-…}` from inputs.claude_version,
//     so a non-empty workflow_dispatch input overrides the pin.
func assertClaudeLiveWorkflowPinsClaudeVersion(workflow string) error {
	if err := assertPinVarDocumented(workflow); err != nil {
		return err
	}

	step, ok := findStepByName(parseWorkflowSteps(workflow), "Install Claude Code")
	if !ok {
		return fmt.Errorf("runtime-live-e2e.yml has no Install Claude Code step")
	}

	// AC-4: the step reads the workflow_dispatch input into CLAUDE_VERSION (the
	// `env:` mapping is parsed off the raw text since parseWorkflowSteps drops
	// `with:`/`env:` keys).
	if !stepEnvBinds(workflow, step, "CLAUDE_VERSION", "${{ inputs.claude_version }}") {
		return fmt.Errorf("runtime-live-e2e.yml Install Claude Code step does not read inputs.claude_version into CLAUDE_VERSION")
	}

	commands := executableShellCommands(step.run)

	// AC-1 + AC-4: the install resolves the pin via the override-precedence
	// expression. `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` means a
	// non-empty input wins (AC-4) and an empty input falls back to the pin (AC-1).
	resolvesPin := false
	for _, c := range commands {
		if strings.Contains(c, "${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}") {
			resolvesPin = true
		}
	}
	if !resolvesPin {
		return fmt.Errorf("runtime-live-e2e.yml Install Claude Code step does not resolve ${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION} (the override-precedence pin fallback)")
	}

	// AC-1 adversarial: no executable command may pipe install.sh into a bash
	// that takes no version arg — that is the no-pin float branch. Every install
	// invocation must pass the resolved version positional (`bash -s -- ...`).
	for _, c := range commands {
		if !strings.Contains(c, "install.sh") || !strings.Contains(c, "| bash") {
			continue
		}
		if !strings.Contains(c, "bash -s --") {
			return fmt.Errorf("runtime-live-e2e.yml Install Claude Code step has an unconditional no-pin install branch (install.sh | bash with no version arg): %q", c)
		}
	}

	return nil
}

// assertPinVarDocumented checks the SPACEDOCK_PINNED_CLAUDE_VERSION line exists,
// pins 2.1.177, and is immediately preceded (ignoring blank lines) by a comment
// block that names the team-tool regression — AC-3's self-documenting proof.
func assertPinVarDocumented(workflow string) error {
	lines := strings.Split(workflow, "\n")
	pinIdx := -1
	for i, raw := range lines {
		if strings.TrimSpace(raw) == `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` {
			pinIdx = i
			break
		}
	}
	if pinIdx < 0 {
		return fmt.Errorf(`runtime-live-e2e.yml does not declare SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177" in the claude-live job env:`)
	}

	// Walk upward over the contiguous comment block directly above the pin.
	var comment strings.Builder
	for i := pinIdx - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "#") {
			comment.WriteString(strings.TrimPrefix(trimmed, "#"))
			comment.WriteByte('\n')
			continue
		}
		break
	}
	rationale := comment.String()
	if rationale == "" {
		return fmt.Errorf("runtime-live-e2e.yml SPACEDOCK_PINNED_CLAUDE_VERSION is not preceded by a rationale comment")
	}
	if !strings.Contains(rationale, "team tool") {
		return fmt.Errorf("runtime-live-e2e.yml pin rationale comment does not name the team-tool regression")
	}
	if !strings.Contains(rationale, "DISABLE_AUTOUPDATER") {
		return fmt.Errorf("runtime-live-e2e.yml pin rationale comment does not note DISABLE_AUTOUPDATER is load-bearing for the pin")
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
