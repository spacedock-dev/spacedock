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

func assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(workflow string) error {
	for _, want := range []string{
		`SPACEDOCK_JOURNEY_METRICS_DIR: ${{ github.workspace }}/live-artifacts/journey-metrics/claude/${{ matrix.model }}`,
		`SPACEDOCK_JOURNEY_METRICS_DIR: ${{ github.workspace }}/live-artifacts/journey-metrics/codex`,
	} {
		if !hasExecutableYAMLLine(workflow, want) {
			return fmt.Errorf("runtime-live-e2e.yml missing active metrics env line %q", want)
		}
	}

	steps := parseWorkflowSteps(workflow)
	claudeRun := findExecutableStep(steps, "Run live Claude E2E", "TestLiveCommon")
	if claudeRun < 0 {
		return fmt.Errorf("runtime-live-e2e.yml has no executable Claude shared scenario run")
	}
	codexRun := findExecutableStep(steps, "Run live Codex shared scenarios", "TestLiveCommon")
	if codexRun < 0 {
		return fmt.Errorf("runtime-live-e2e.yml has no executable Codex shared scenario run")
	}
	piCoverageRun := findExecutableStep(steps, "Run live Pi common journeys", "TestLiveCommon")
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
	if workflowHasExecutableCommandContaining(workflow, `npm config get min-release-age`) || workflowHasExecutableCommandContaining(workflow, `npm config set min-release-age`) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job uses obsolete min-release-age probing")
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
	if hasExecutableYAMLLine(workflow, `SPACEDOCK_JOURNEY_METRICS_DIR: ${{ github.workspace }}/live-artifacts/journey-metrics/pi`) {
		return fmt.Errorf("runtime-live-e2e.yml retains a Pi journey-metrics path without a producer")
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
	if err := assertGoreleaserDoesNotNeedJourneyLedger(workflow); err != nil {
		return err
	}
	return nil
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

// assertReleaseLedgerStepsSkipWhenNoProducerRun proves POLICY 1's job-level
// skip: the journey-cost builder and its release-upload step must be GATED on
// the download step's producer-found output, so a producer-less / empty-dir cut
// SKIPS them (journey-ledger job green) instead of running `journey-costs` over
// an empty dir and REDding the job. It locates the download step by its emitted
// `found=` output and requires both the builder step and the upload step to
// carry an `if:` referencing that step's `found` output.
func assertReleaseLedgerStepsSkipWhenNoProducerRun(workflow string) error {
	steps := parseWorkflowSteps(workflow)

	downloadID := ""
	for _, step := range steps {
		if step.id == "" {
			continue
		}
		for _, command := range executableShellCommands(step.run) {
			if strings.Contains(command, `found=false`) && strings.Contains(command, `"$GITHUB_OUTPUT"`) {
				downloadID = step.id
			}
		}
	}
	if downloadID == "" {
		return fmt.Errorf("release.yml has no download step with an id that emits found=false to $GITHUB_OUTPUT (the non-fatal skip flag)")
	}

	gate := fmt.Sprintf("steps.%s.outputs.found == 'true'", downloadID)
	builderGated, publishGated := false, false
	for _, step := range steps {
		isBuilder, isPublish := false, false
		for _, command := range executableShellCommands(step.run) {
			if isJourneyCostBuilder(command) {
				isBuilder = true
			}
			if strings.HasPrefix(command, `gh release upload `) &&
				strings.Contains(command, `"$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`) {
				isPublish = true
			}
		}
		if isBuilder && step.ifCond == gate {
			builderGated = true
		}
		if isPublish && step.ifCond == gate {
			publishGated = true
		}
	}
	if !builderGated {
		return fmt.Errorf("release.yml journey-cost builder is not gated on %q; a producer-less cut would run journey-costs over an empty dir and RED the journey-ledger job", gate)
	}
	if !publishGated {
		return fmt.Errorf("release.yml journey-cost publish step is not gated on %q; a producer-less cut would attempt an upload of a missing ledger", gate)
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

// findExecutableStep returns the index of the step named name whose run block
// contains commandFragment, or -1.
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
