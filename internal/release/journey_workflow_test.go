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
		`go test -tags live -run TestLiveClaudeSharedScenarios ./internal/ensigncycle/ -v 2>&1 | tee claude-shared-scenarios-transcript.txt`,
		`# go test -tags live -run TestLiveClaudeSharedScenarios ./internal/ensigncycle/ -v 2>&1 | tee claude-shared-scenarios-transcript.txt`,
		1)

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted a workflow without an executable Claude shared scenario run")
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
