// ABOUTME: Structural guard that the pi-live lane consumes the single Pi stamp
// ABOUTME: source (internal/pilive) and carries no version, integrity, or JS.
package release

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/pilive"
)

// nodeInvocation matches an executable `node` token in a shell script (the
// `node` in setup-node/node-version lives in `uses:`/`with:`, never a run).
var nodeInvocation = regexp.MustCompile(`(?:^|[\s;&|()])node(?:\s|$)`)

var nodeHeredoc = regexp.MustCompile(`<<-?'?NODE'?`)

// assertPiLiveLaneConsumesTheStampSource binds the pi-live lane to the single
// stamp source: its steps may only invoke the helper commands, the helper is
// built before the install step needs it, no run script contains a node
// invocation or heredoc, no line carries a pinned version or integrity, and the
// install step never exports either package-root override.
func assertPiLiveLaneConsumesTheStampSource(workflow string) error {
	job, ok := piLiveJob(workflow)
	if !ok {
		return fmt.Errorf("workflow lacks the pi-live job")
	}
	want := map[string]string{
		"Install Pi CLI and substrates":    "spacedock-pilive install",
		"Guard Pi substrate compatibility": "spacedock-pilive guard",
	}
	indexOf := map[string]int{}
	for i, step := range job.steps {
		indexOf[step.name] = i
	}
	build, ok := stepNamed(job.steps, "Build spacedock binary")
	if !ok {
		return fmt.Errorf("pi-live lacks the helper build step")
	}
	if !strings.Contains(build.run, "./cmd/spacedock-pilive") {
		return fmt.Errorf("pi-live build step does not build the helper (./cmd/spacedock-pilive)")
	}
	for stepName, wantCommand := range want {
		step, ok := stepNamed(job.steps, stepName)
		if !ok {
			return fmt.Errorf("pi-live lacks step %q", stepName)
		}
		if step.ifCond != "" {
			return fmt.Errorf("pi-live step %q must be unconditional", stepName)
		}
		if !containsCommand(step.run, wantCommand) {
			return fmt.Errorf("pi-live step %q does not invoke %q", stepName, wantCommand)
		}
	}
	buildIdx, installIdx := indexOf["Build spacedock binary"], indexOf["Install Pi CLI and substrates"]
	if buildIdx >= installIdx {
		return fmt.Errorf("pi-live must build the helper (step %d) before installing with it (step %d)", buildIdx, installIdx)
	}

	checkout, ok := stepNamed(job.steps, "Verify Pi current-checkout setup")
	if !ok {
		return fmt.Errorf("pi-live lacks the current-checkout setup step")
	}
	if !strings.Contains(checkout.run, "spacedock-pilive verify-manifest ") {
		return fmt.Errorf("current-checkout setup does not invoke spacedock-pilive verify-manifest")
	}

	for _, step := range job.steps {
		if nodeInvocation.MatchString(step.run) || nodeHeredoc.MatchString(step.run) {
			return fmt.Errorf("pi-live step %q still contains JavaScript; the lane may only invoke commands", step.name)
		}
	}

	for _, pkg := range pilive.Packages {
		for _, literal := range []string{pkg.Version, pkg.Integrity} {
			if lineWith(workflow, literal) != "" {
				return fmt.Errorf("pi-live workflow still writes the stamp literal %q; it belongs only in internal/pilive", literal)
			}
		}
	}

	install := job.steps[indexOf["Install Pi CLI and substrates"]]
	for _, command := range executableShellCommands(install.run) {
		if !strings.Contains(command, "PI_SUBAGENTS_PACKAGE_ROOT") && !strings.Contains(command, "PI_INTERCOM_PACKAGE_ROOT") {
			continue
		}
		if strings.Contains(command, "GITHUB_ENV") || strings.HasPrefix(command, "export ") {
			return fmt.Errorf("pi-live install step must not export package-root overrides: %s", command)
		}
	}
	return nil
}

func piLiveJob(workflow string) (workflowJob, bool) {
	for _, job := range parseWorkflowJobs(workflow) {
		if job.name == "pi-live" {
			return job, true
		}
	}
	return workflowJob{}, false
}

func containsCommand(script, want string) bool {
	for _, command := range executableShellCommands(script) {
		if command == want || strings.HasPrefix(command, want+" ") {
			return true
		}
	}
	return false
}

// lineWith returns the first non-comment workflow line containing literal, or ""
// when the literal does not appear in an active line.
func lineWith(workflow, literal string) string {
	for _, line := range strings.Split(workflow, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") && strings.Contains(line, literal) {
			return line
		}
	}
	return ""
}

func TestPiLiveLaneConsumesTheStampSource(t *testing.T) {
	if err := assertPiLiveLaneConsumesTheStampSource(readWorkflowFiles(t)["runtime-live-e2e.yml"]); err != nil {
		t.Fatal(err)
	}
}

func TestPiLiveLaneGuardRejectsMutations(t *testing.T) {
	workflow := readWorkflowFiles(t)["runtime-live-e2e.yml"]
	if err := assertPiLiveLaneConsumesTheStampSource(workflow); err != nil {
		t.Fatal(err)
	}
	installAnchor := "spacedock-pilive install"
	mutations := []struct{ name, old, replacement string }{
		{"install-command-removed", installAnchor, "echo skipped-install"},
		{"guard-command-removed", "spacedock-pilive guard", "echo skipped-guard"},
		{"verify-manifest-removed", "spacedock-pilive verify-manifest ", "echo skipped-verify "},
		{"helper-build-removed", "go build -o ./spacedock-pilive ./cmd/spacedock-pilive", "true"},
		{"node-reintroduced", installAnchor, installAnchor + "\n          node -e \"process.exit(0)\""},
		{"node-heredoc-reintroduced", installAnchor, installAnchor + "\n          node - <<'NODE'\n          NODE"},
		{"version-literal-reintroduced", installAnchor, installAnchor + "\n          PI_CODING_AGENT_VERSION=\"" + pilive.PiCodingAgentVersion + "\""},
		{"integrity-literal-reintroduced", installAnchor, installAnchor + "\n          PI_SUBAGENTS_INTEGRITY=\"" + pilive.PiSubagentsIntegrity + "\""},
		{"override-export-reintroduced", installAnchor, installAnchor + "\n          echo \"PI_SUBAGENTS_PACKAGE_ROOT=$HOME/.pi/agent/npm/node_modules/pi-subagents\" >> \"$GITHUB_ENV\""},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			changed := strings.Replace(workflow, mutation.old, mutation.replacement, 1)
			if changed == workflow {
				t.Fatal("mutation did not change workflow")
			}
			if err := assertPiLiveLaneConsumesTheStampSource(changed); err == nil {
				t.Fatal("mutation escaped guard")
			}
		})
	}
}
