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

	builder := "go run ./cmd/spacedock-release journey-costs"
	goreleaser := "goreleaser/goreleaser-action"
	if !strings.Contains(release, builder) {
		t.Fatalf("release.yml does not run the journey-cost builder")
	}
	if strings.Index(release, builder) > strings.Index(release, goreleaser) {
		t.Fatalf("release.yml runs journey-cost builder after goreleaser; builder must run before publish")
	}
	for _, want := range []string{
		"--metrics-dir \"$RUNNER_TEMP/journey-metrics\"",
		"--out \"$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json\"",
		"gh release upload \"$GITHUB_REF_NAME\" \"$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json\"",
	} {
		if !strings.Contains(release, want) {
			t.Fatalf("release.yml missing %q", want)
		}
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
