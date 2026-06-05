package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowsPreserveAndPublishJourneyCosts(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	release := readWorkflow(t, "release.yml")

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(live); err != nil {
		t.Fatal(err)
	}

	if err := assertReleaseWorkflowPublishesJourneyCosts(release); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeLiveWorkflowGuardRejectsCommentOnlyJourneyMetricUpload(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.ReplaceAll(live,
		`live-artifacts/journey-metrics/**`,
		`# live-artifacts/journey-metrics/**`)

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted journey metrics paths that were only inert upload text")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsMissingCodexJourneyMetricUpload(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	first := strings.Index(live, `live-artifacts/journey-metrics/**`)
	if first < 0 {
		t.Fatal("fixture workflow missing first journey metrics upload path")
	}
	second := strings.Index(live[first+1:], `live-artifacts/journey-metrics/**`)
	if second < 0 {
		t.Fatal("fixture workflow missing second journey metrics upload path")
	}
	second += first + 1
	adversarial := live[:second] + `live-artifacts/codex/**` + live[second+len(`live-artifacts/journey-metrics/**`):]

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted a workflow where only one live job uploads journey metrics")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsMissingSharedScenarioRun(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.Replace(live,
		`go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle/ -v 2>&1 | tee claude-shared-scenarios-transcript.txt`,
		`# go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle/ -v 2>&1 | tee claude-shared-scenarios-transcript.txt`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing Claude shared scenario run command")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted a workflow without an executable Claude shared scenario run")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsUnscopedPiPackage(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.Replace(live,
		`npm install -g @earendil-works/pi-coding-agent --before="$NPM_BEFORE" --ignore-scripts --no-audit --no-fund --omit=dev`,
		`npm install -g pi-coding-agent --before="$NPM_BEFORE" --ignore-scripts --no-audit --no-fund --omit=dev`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing scoped Pi CLI install command")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted the wrong unscoped Pi CLI package")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsMissingPiBeforeAgeGate(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.ReplaceAll(live, ` --before="$NPM_BEFORE"`, ``)
	if adversarial == live {
		t.Fatal("fixture workflow missing npm --before install flags")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted Pi npm installs without --before")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsObsoletePiMinReleaseAgeProbe(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.Replace(live,
		`NPM_BEFORE="$(node -e 'console.log(new Date(Date.now() - 24*60*60*1000).toISOString())')"
          echo "Using npm --before age gate for pi-live installs: $NPM_BEFORE"`,
		`npm config get min-release-age
          npm config set min-release-age 1440`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing npm --before age-gate timestamp")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted obsolete min-release-age probing")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsUnverifiedPiPackageInstall(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.Replace(live,
		`node -e "const p=require('$global_npm_root/@earendil-works/pi-coding-agent/package.json'); if (p.name !== '@earendil-works/pi-coding-agent') throw new Error('unexpected Pi package name '+p.name); if (!p.bin || p.bin.pi !== 'dist/cli.js') throw new Error('unexpected Pi bin '+JSON.stringify(p.bin)); console.log('verified '+p.name+'@'+p.version+' bin pi='+p.bin.pi)"`,
		`echo "skipping Pi package verification"`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing Pi CLI package verification command")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted an unverified Pi CLI package install")
	}
}

func TestReleaseWorkflowGuardRejectsCommentOnlyJourneyCostBuilder(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	adversarial := strings.Replace(release,
		`go run ./cmd/spacedock-release journey-costs "$RELEASE_VERSION" \`,
		`# go run ./cmd/spacedock-release journey-costs "$RELEASE_VERSION" \`,
		1)
	adversarial = strings.Replace(adversarial,
		`test -s "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`,
		`printf '{}' > "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"
          test -s "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`,
		1)

	if err := assertReleaseWorkflowPublishesJourneyCosts(adversarial); err == nil {
		t.Fatal("release workflow guard accepted a commented builder plus fake JSON output")
	}
}

func TestReleaseWorkflowGuardRejectsCommentOnlyJourneyCostPublish(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	adversarial := strings.Replace(release,
		`gh release upload "$GITHUB_REF_NAME" "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json" --clobber`,
		`# gh release upload "$GITHUB_REF_NAME" "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json" --clobber
          echo "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`,
		1)

	if err := assertReleaseWorkflowPublishesJourneyCosts(adversarial); err == nil {
		t.Fatal("release workflow guard accepted a commented publish command")
	}
}

func readWorkflow(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
