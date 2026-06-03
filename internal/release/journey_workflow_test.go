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

	for _, want := range []string{
		"SPACEDOCK_JOURNEY_METRICS_DIR",
		"live-artifacts/journey-metrics/**",
		"actions/upload-artifact@v4",
	} {
		if !strings.Contains(live, want) {
			t.Fatalf("runtime-live-e2e.yml missing %q", want)
		}
	}

	if err := assertReleaseWorkflowPublishesJourneyCosts(release); err != nil {
		t.Fatal(err)
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
