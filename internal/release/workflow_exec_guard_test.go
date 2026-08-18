package release

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type workflowStep struct {
	name     string
	id       string
	ifCond   string
	uses     string
	run      string
	withPath string
}

// workflowJob is a top-level entry under `jobs:` — its name, its declared
// `needs:` edges, and the steps it owns. The job graph is parsed separately from
// the flat step list so the separation guard can bind a `needs:` edge to the
// OWNING job (the job carrying the goreleaser action) rather than matching the
// edge anywhere in the document.
type workflowJob struct {
	name  string
	needs []string
	steps []workflowStep
}

// assertGoreleaserDoesNotNeedJourneyLedger binds the separation invariant to the
// job DAG: no job carrying the goreleaser action may declare `needs:` on a job
// that builds the journey-cost ledger. That edge would re-block the cut on the
// never-fired Runtime-Live-E2E run (the exact bug this separation closes). It
// collects ALL goreleaser-action carriers and ALL ledger-builder carriers and
// rejects if ANY carrier→ledger edge exists — so the result does not depend on
// Go's randomized map-iteration order (a last-wins scan over a multi-carrier
// workflow would be latently flaky). The SAFE reverse edge (journey-ledger
// needs: goreleaser, required so `gh release upload` runs after the Release
// exists) points the other way and does not false-trip the check; a job that is
// itself both carriers cannot need itself, so the self-pair is skipped.
func assertGoreleaserDoesNotNeedJourneyLedger(workflow string) error {
	jobs := parseWorkflowJobs(workflow)
	ledgerJobs := map[string]bool{}
	var goreleaserCarriers []workflowJob
	for _, job := range jobs {
		isGoreleaser, isLedger := false, false
		for _, step := range job.steps {
			if strings.HasPrefix(step.uses, "goreleaser/goreleaser-action@") {
				isGoreleaser = true
			}
			for _, command := range executableShellCommands(step.run) {
				if isJourneyCostBuilder(command) {
					isLedger = true
				}
			}
		}
		if isLedger {
			ledgerJobs[job.name] = true
		}
		if isGoreleaser {
			goreleaserCarriers = append(goreleaserCarriers, job)
		}
	}
	if len(goreleaserCarriers) == 0 {
		return fmt.Errorf("release.yml has no job carrying the goreleaser action")
	}
	if len(ledgerJobs) == 0 {
		return fmt.Errorf("release.yml has no job carrying the journey-cost builder")
	}
	for _, carrier := range goreleaserCarriers {
		for _, need := range carrier.needs {
			if need != carrier.name && ledgerJobs[need] {
				return fmt.Errorf("release.yml goreleaser job %q declares needs: %q — re-blocking the cut on the journey-ledger job", carrier.name, need)
			}
		}
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
		case strings.HasPrefix(trimmed, "id: "):
			step.id = strings.Trim(strings.TrimPrefix(trimmed, "id: "), `"`)
		case strings.HasPrefix(trimmed, "if: "):
			step.ifCond = strings.TrimSpace(strings.TrimPrefix(trimmed, "if: "))
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

// parseWorkflowJobs reads the workflow's whole job graph — job names, their
// `needs:` edges, and the steps each owns — from a single gopkg.in/yaml.v3 pass.
// Job identity is the security-relevant axis for the separation guard, and a
// real YAML parse resolves what a line-walk cannot: anchored/aliased needs
// (`needs: *anchor`), quoted job keys (`"goreleaser":`), comment-trailed and
// blank-line-split sequences, and flow/block/scalar forms all normalize the same
// way GitHub Actions resolves them. A workflow that does not parse as YAML
// yields no jobs (the guard's text-level checks still run against the source).
func parseWorkflowJobs(workflow string) []workflowJob {
	var doc struct {
		Jobs map[string]struct {
			Needs needsList `yaml:"needs"`
			Steps []struct {
				Name string `yaml:"name"`
				ID   string `yaml:"id"`
				If   string `yaml:"if"`
				Uses string `yaml:"uses"`
				Run  string `yaml:"run"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal([]byte(workflow), &doc); err != nil {
		return nil
	}
	var jobs []workflowJob
	for name, job := range doc.Jobs {
		wj := workflowJob{name: name, needs: job.Needs}
		for _, step := range job.Steps {
			wj.steps = append(wj.steps, workflowStep{
				name:   step.Name,
				id:     step.ID,
				ifCond: step.If,
				uses:   step.Uses,
				run:    step.Run,
			})
		}
		jobs = append(jobs, wj)
	}
	return jobs
}

// needsList decodes a job's `needs:` — GitHub Actions allows a scalar
// (`needs: a`) or a sequence (flow `[a, b]` or block list), and either may be an
// anchor/alias. Decoding through this custom type lets yaml.v3 resolve aliases
// to their anchored value and normalize quoting before we read the names, so no
// `needs:` shape can hide a re-coupling edge from the separation guard.
type needsList []string

func (n *needsList) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		if value.Value != "" {
			*n = needsList{value.Value}
		}
		return nil
	case yaml.SequenceNode:
		var names []string
		if err := value.Decode(&names); err != nil {
			return err
		}
		*n = names
		return nil
	default:
		return nil
	}
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
