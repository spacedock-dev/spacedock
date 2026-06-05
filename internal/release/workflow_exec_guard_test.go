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
	if workflowHasExecutableCommandContaining(workflow, `npm config get min-release-age`) || workflowHasExecutableCommandContaining(workflow, `npm config set min-release-age`) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job uses obsolete min-release-age probing instead of npm --before")
	}
	if !workflowHasExecutableCommandContaining(workflow, `NPM_BEFORE="$(node -e 'console.log(new Date(Date.now() - 24*60*60*1000).toISOString())')"`) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job does not compute an npm --before age-gate timestamp")
	}
	if !workflowHasExecutableCommandContaining(workflow, `npm install -g @earendil-works/pi-coding-agent --before="$NPM_BEFORE" --ignore-scripts --no-audit --no-fund --omit=dev`) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job does not install the scoped Pi CLI with npm --before and safer npm flags")
	}
	if workflowHasExecutableCommandContaining(workflow, "npm install -g pi-coding-agent") || workflowHasExecutableCommandContaining(workflow, "npm install -g \"pi-coding-agent@") {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job installs the wrong unscoped Pi CLI package")
	}
	if !workflowHasExecutableCommandContaining(workflow, `npm install --prefix "$pi_npm_root" pi-subagents pi-intercom --before="$NPM_BEFORE" --ignore-scripts --no-audit --no-fund --omit=dev`) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job does not directly install required Pi substrates with npm --before and safer npm flags")
	}
	if !workflowHasExecutableCommandContaining(workflow, `test -x "$(command -v pi)"`) || !workflowHasExecutableCommandContaining(workflow, `p.name !== '@earendil-works/pi-coding-agent'`) || !workflowHasExecutableCommandContaining(workflow, `p.bin.pi !== 'dist/cli.js'`) || !workflowHasExecutableCommandContaining(workflow, `p.name !== 'pi-subagents'`) || !workflowHasExecutableCommandContaining(workflow, `p.name !== 'pi-intercom'`) || !workflowHasExecutableCommandContaining(workflow, `test -f "$pi_npm_root/node_modules/pi-subagents/src/extension/index.ts"`) || !workflowHasExecutableCommandContaining(workflow, `test -f "$pi_npm_root/node_modules/pi-subagents/skills/pi-subagents/SKILL.md"`) || !workflowHasExecutableCommandContaining(workflow, `test -f "$pi_npm_root/node_modules/pi-subagents/src/intercom/intercom-bridge.ts"`) || !workflowHasExecutableCommandContaining(workflow, `test -f "$pi_npm_root/node_modules/pi-intercom/skills/pi-intercom/SKILL.md"`) || !workflowHasExecutableCommandContaining(workflow, `echo "PI_INTERCOM_PACKAGE_ROOT=$pi_npm_root/node_modules/pi-intercom" >> "$GITHUB_ENV"`) {
		return fmt.Errorf("runtime-live-e2e.yml Pi live job does not verify installed Pi package names, versions, bin, resource paths, and intercom package env")
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
	if err := assertGoreleaserDoesNotNeedJourneyLedger(workflow); err != nil {
		return err
	}
	return nil
}

// assertGoreleaserDoesNotNeedJourneyLedger binds the separation invariant to the
// job DAG: the job that carries the goreleaser action must NOT declare `needs:`
// on the job that builds the journey-cost ledger. That edge would re-block the
// cut on the never-fired Runtime-Live-E2E run (the exact bug this separation
// closes). The edge is bound to the OWNING job — found by which job holds the
// goreleaser action vs. which holds the journey-cost builder — so the SAFE
// reverse edge (journey-ledger needs: goreleaser, required so `gh release upload`
// runs after the Release exists) does not false-trip the check.
func assertGoreleaserDoesNotNeedJourneyLedger(workflow string) error {
	jobs := parseWorkflowJobs(workflow)
	goreleaserJob, ledgerJob := "", ""
	for _, job := range jobs {
		for _, step := range job.steps {
			if strings.HasPrefix(step.uses, "goreleaser/goreleaser-action@") {
				goreleaserJob = job.name
			}
			for _, command := range executableShellCommands(step.run) {
				if isJourneyCostBuilder(command) {
					ledgerJob = job.name
				}
			}
		}
	}
	if goreleaserJob == "" {
		return fmt.Errorf("release.yml has no job carrying the goreleaser action")
	}
	if ledgerJob == "" {
		return fmt.Errorf("release.yml has no job carrying the journey-cost builder")
	}
	if goreleaserJob == ledgerJob {
		// Single-job layout: no cross-job needs edge to re-block the cut.
		return nil
	}
	for _, job := range jobs {
		if job.name != goreleaserJob {
			continue
		}
		for _, need := range job.needs {
			if need == ledgerJob {
				return fmt.Errorf("release.yml goreleaser job %q declares needs: %q — re-blocking the cut on the journey-ledger job", goreleaserJob, ledgerJob)
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

// parseWorkflowJobs partitions a workflow into its top-level jobs (the 2-space
// indented keys under `jobs:`), pairing each job's `needs:` edges (from a real
// YAML parse, parseJobNeeds) with the steps it owns (sliced from the job's text
// span and run through the flat step parser). This is the job-aware view the
// separation guard needs: the flat parseWorkflowSteps cannot tell which job a
// step belongs to, so a `goreleaser → journey-ledger` edge would otherwise be
// invisible.
func parseWorkflowJobs(workflow string) []workflowJob {
	lines := strings.Split(workflow, "\n")
	jobsStart := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "jobs:" && leadingSpaces(line) == 0 {
			jobsStart = i + 1
			break
		}
	}
	if jobsStart < 0 {
		return nil
	}

	// `needs:` edges come from a real YAML parse (parseJobNeeds), not the
	// line-walk: needs is the security-relevant field for the separation guard,
	// and a hand-rolled scan misses comment-trailed entries, blank-line-split
	// block lists, anchors, and quoting that a parser handles for free.
	needsByJob := parseJobNeeds(workflow)

	var jobs []workflowJob
	for i := jobsStart; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		// A top-level key at column 0 ends the jobs block (e.g. a trailing
		// document key). Job keys sit at exactly 2-space indent.
		if leadingSpaces(line) == 0 {
			break
		}
		if leadingSpaces(line) == 2 && strings.HasSuffix(trimmed, ":") {
			name := strings.TrimSuffix(trimmed, ":")
			jobs = append(jobs, workflowJob{name: name, needs: needsByJob[name]})
			continue
		}
	}

	// Attach steps by slicing each job's text span and reusing the flat parser.
	jobHeaderAt := func(name string) int {
		for i := jobsStart; i < len(lines); i++ {
			if leadingSpaces(lines[i]) == 2 && strings.TrimSpace(lines[i]) == name+":" {
				return i
			}
		}
		return -1
	}
	for j := range jobs {
		start := jobHeaderAt(jobs[j].name)
		end := len(lines)
		if j+1 < len(jobs) {
			end = jobHeaderAt(jobs[j+1].name)
		}
		if start >= 0 && end > start {
			jobs[j].steps = parseWorkflowSteps(strings.Join(lines[start:end], "\n"))
		}
	}
	return jobs
}

// parseJobNeeds reads every job's `needs:` edges via a real YAML parse, keyed by
// job name. GitHub Actions allows `needs:` as a scalar (`needs: a`) or a
// sequence (flow `[a, b]` or block list); yaml.v3 normalizes all of them — and
// strips inline `# comments`, tolerates blank lines between block-list entries,
// and resolves quoting and anchors — so no syntactic shape of a re-coupling
// edge can hide from the separation guard. A workflow that does not parse as
// YAML yields no edges (the guard's other checks still run against the text).
func parseJobNeeds(workflow string) map[string][]string {
	var doc struct {
		Jobs map[string]struct {
			Needs yaml.Node `yaml:"needs"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal([]byte(workflow), &doc); err != nil {
		return nil
	}
	out := make(map[string][]string, len(doc.Jobs))
	for name, job := range doc.Jobs {
		switch job.Needs.Kind {
		case yaml.ScalarNode:
			if job.Needs.Value != "" {
				out[name] = []string{job.Needs.Value}
			}
		case yaml.SequenceNode:
			var needs []string
			for _, item := range job.Needs.Content {
				if item.Kind == yaml.ScalarNode && item.Value != "" {
					needs = append(needs, item.Value)
				}
			}
			out[name] = needs
		}
	}
	return out
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
