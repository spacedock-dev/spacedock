package release

import (
	"fmt"
	"strings"
)

type workflowStep struct {
	name     string
	uses     string
	run      string
	withPath string
}

func assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(workflow string) error {
	for _, want := range []string{
		`SPACEDOCK_JOURNEY_METRICS_DIR: ${{ github.workspace }}/live-artifacts/journey-metrics/claude/${{ matrix.model }}`,
		`SPACEDOCK_JOURNEY_METRICS_DIR: ${{ github.workspace }}/live-artifacts/journey-metrics/codex`,
		`SPACEDOCK_JOURNEY_METRICS_DIR: ${{ github.workspace }}/live-artifacts/journey-metrics/pi`,
	} {
		if !hasExecutableYAMLLine(workflow, want) {
			return fmt.Errorf("runtime-live-e2e.yml missing active metrics env line %q", want)
		}
	}

	steps := parseWorkflowSteps(workflow)
	claudeRun := findExecutableStep(steps, "Run live Claude shared scenarios", "TestLiveClaudeSharedScenarios")
	if claudeRun < 0 {
		return fmt.Errorf("runtime-live-e2e.yml has no executable Claude shared scenario run")
	}
	codexRun := findExecutableStep(steps, "Run live Codex shared scenarios", "TestLiveCodexSharedScenarios")
	if codexRun < 0 {
		return fmt.Errorf("runtime-live-e2e.yml has no executable Codex shared scenario run")
	}
	piCoverageRun := findExecutableStep(steps, "Run Pi shared scenario coverage guard", "TestPiSharedScenarioCoverage")
	if piCoverageRun < 0 {
		return fmt.Errorf("runtime-live-e2e.yml has no executable Pi shared scenario coverage guard")
	}
	piSmokeRun := findExecutableStep(steps, "Run live Pi front-door smoke", "TestLivePiFrontDoorSmoke")
	if piSmokeRun < 0 {
		return fmt.Errorf("runtime-live-e2e.yml has no executable Pi front-door smoke")
	}
	if !hasExecutableYAMLLine(workflow, `OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}`) || !hasExecutableYAMLLine(workflow, `SPACEDOCK_PI_LIVE_REQUIRED: "1"`) || !hasExecutableYAMLLine(workflow, `name: CI-E2E-PI`) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job is missing its OpenAI secret, required flag, or CI-E2E-PI environment")
	}
	if !workflowHasExecutableCommandContaining(workflow, "pi install npm:pi-subagents") || !workflowHasExecutableCommandContaining(workflow, "pi install npm:pi-intercom") {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job does not install required Pi substrates")
	}
	if !workflowHasExecutableCommandContaining(workflow, `spacedock doctor --host pi --plugin-dir "$GITHUB_WORKSPACE"`) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job does not verify current-checkout Spacedock skills")
	}
	if !hasJourneyMetricsUploadAfter(steps, claudeRun, codexRun) {
		return fmt.Errorf("runtime-live-e2e.yml Claude shared scenario job does not upload raw journey metrics")
	}
	if !hasJourneyMetricsUploadAfter(steps, codexRun, piCoverageRun) {
		return fmt.Errorf("runtime-live-e2e.yml Codex shared scenario job does not upload raw journey metrics")
	}
	if !hasJourneyMetricsUploadAfter(steps, piSmokeRun, len(steps)) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job does not upload raw journey metrics")
	}
	return nil
}

func assertReleaseWorkflowPublishesJourneyCosts(workflow string) error {
	steps := parseWorkflowSteps(workflow)
	builderStep, goreleaserStep, publishStep := -1, -1, -1
	builderHasOutput := false
	builderChecksOutput := false

	for i, step := range steps {
		if strings.HasPrefix(step.uses, "goreleaser/goreleaser-action@") {
			goreleaserStep = i
		}
		for _, command := range executableShellCommands(step.run) {
			switch {
			case isJourneyCostBuilder(command):
				builderStep = i
				builderHasOutput = strings.Contains(command, `--out "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`)
			case command == `test -s "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`:
				builderChecksOutput = true
			case strings.HasPrefix(command, `gh release upload `) &&
				strings.Contains(command, `"$GITHUB_REF_NAME"`) &&
				strings.Contains(command, `"$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`):
				publishStep = i
			}
		}
	}

	if builderStep < 0 {
		return fmt.Errorf("release.yml has no executable journey-cost builder command")
	}
	if !builderHasOutput {
		return fmt.Errorf("release.yml journey-cost builder does not write journey-costs-v${RELEASE_VERSION}.json")
	}
	if !builderChecksOutput {
		return fmt.Errorf("release.yml does not check the generated journey-cost ledger is non-empty")
	}
	if goreleaserStep < 0 {
		return fmt.Errorf("release.yml has no goreleaser publish step")
	}
	if builderStep > goreleaserStep {
		return fmt.Errorf("release.yml builds journey costs after goreleaser")
	}
	if publishStep < 0 {
		return fmt.Errorf("release.yml has no executable journey-cost release upload command")
	}
	if publishStep <= builderStep {
		return fmt.Errorf("release.yml publishes journey costs before building them")
	}
	return nil
}

func isJourneyCostBuilder(command string) bool {
	return strings.HasPrefix(command, `go run ./cmd/spacedock-release journey-costs `) &&
		strings.Contains(command, `"$RELEASE_VERSION"`) &&
		strings.Contains(command, `--metrics-dir "$RUNNER_TEMP/journey-metrics"`)
}

func parseWorkflowSteps(workflow string) []workflowStep {
	lines := strings.Split(workflow, "\n")
	var steps []workflowStep
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- name: ") {
			steps = append(steps, workflowStep{name: strings.Trim(strings.TrimPrefix(trimmed, "- name: "), `"`)})
			continue
		}
		if len(steps) == 0 {
			continue
		}
		step := &steps[len(steps)-1]
		switch {
		case strings.HasPrefix(trimmed, "uses: "):
			step.uses = strings.Trim(strings.TrimPrefix(trimmed, "uses: "), `"`)
		case strings.HasPrefix(trimmed, "run: |"):
			baseIndent := leadingSpaces(line)
			var block []string
			for i+1 < len(lines) {
				next := lines[i+1]
				if strings.TrimSpace(next) != "" && leadingSpaces(next) <= baseIndent {
					break
				}
				block = append(block, next)
				i++
			}
			step.run = strings.Join(block, "\n")
		case strings.HasPrefix(trimmed, "run: "):
			step.run = strings.Trim(strings.TrimPrefix(trimmed, "run: "), `"`)
		case strings.HasPrefix(trimmed, "path: |"):
			baseIndent := leadingSpaces(line)
			var block []string
			for i+1 < len(lines) {
				next := lines[i+1]
				if strings.TrimSpace(next) != "" && leadingSpaces(next) <= baseIndent {
					break
				}
				block = append(block, next)
				i++
			}
			step.withPath = strings.Join(block, "\n")
		}
	}
	return steps
}

func findExecutableStep(steps []workflowStep, name, commandFragment string) int {
	for i, step := range steps {
		if step.name != name {
			continue
		}
		for _, command := range executableShellCommands(step.run) {
			if strings.Contains(command, commandFragment) {
				return i
			}
		}
	}
	return -1
}

func workflowHasExecutableCommandContaining(workflow, want string) bool {
	for _, step := range parseWorkflowSteps(workflow) {
		for _, command := range executableShellCommands(step.run) {
			if strings.Contains(command, want) {
				return true
			}
		}
	}
	return false
}

func hasJourneyMetricsUploadAfter(steps []workflowStep, start, stop int) bool {
	for i := start + 1; i < stop && i < len(steps); i++ {
		step := steps[i]
		if !strings.HasPrefix(step.uses, "actions/upload-artifact@") {
			continue
		}
		if pathBlockContainsLine(step.withPath, "live-artifacts/journey-metrics/**") {
			return true
		}
	}
	return false
}

func pathBlockContainsLine(block, want string) bool {
	for _, raw := range strings.Split(block, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == want {
			return true
		}
	}
	return false
}

func hasExecutableYAMLLine(doc, want string) bool {
	for _, raw := range strings.Split(doc, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == want {
			return true
		}
	}
	return false
}

func executableShellCommands(script string) []string {
	var commands []string
	var continued strings.Builder
	for _, raw := range strings.Split(script, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if continued.Len() > 0 {
			continued.WriteByte(' ')
		}
		if strings.HasSuffix(line, `\`) {
			continued.WriteString(strings.TrimSpace(strings.TrimSuffix(line, `\`)))
			continue
		}
		continued.WriteString(line)
		commands = append(commands, strings.Join(strings.Fields(continued.String()), " "))
		continued.Reset()
	}
	if continued.Len() > 0 {
		commands = append(commands, strings.Join(strings.Fields(continued.String()), " "))
	}
	return commands
}

func leadingSpaces(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}
