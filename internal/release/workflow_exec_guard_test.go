package release

import (
	"fmt"
	"strings"
)

type workflowStep struct {
	name string
	uses string
	run  string
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
		}
	}
	return steps
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
