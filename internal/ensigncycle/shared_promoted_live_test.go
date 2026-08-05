//go:build live

package ensigncycle

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
)

func runFullEnsignCycleJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	root := writeRealisticLifecycleFixture(t)
	result := driver.run(t, scenario, root, "Use $spacedock:first-officer for this whole run. Drive the workflow to completion; you have the conn to resolve gates from each stage report's verdict (auto-approve).")
	entity, _, found := locateEntity(root, "make-it-work")
	if !found {
		t.Fatalf("full ensign cycle left no durable entity; artifacts: %s", result.artifactDir)
	}
	for _, want := range []string{"status: done", "## Stage Report:", "- DONE:", "### Summary"} {
		if !strings.Contains(entity, want) {
			t.Fatalf("full ensign cycle entity missing %q; artifacts: %s\n%s", want, result.artifactDir, entity)
		}
	}
	if !someCommitNamesOnly(t, root, "make-it-work") {
		t.Fatalf("full ensign cycle has no path-scoped entity commit; artifacts: %s", result.artifactDir)
	}
}

//spacedock:live-fixture id=realistic-lifecycle
func writeRealisticLifecycleFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeRealisticLifecycle())
	writeFile(t, filepath.Join(root, "make-it-work.md"), entityFixture())
	gitInit(t, root)
	return root
}

func runZeroDiscoveryJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	root := writeZeroDiscoveryFixture(t)
	result := driver.run(t, scenario, root, "Use $spacedock:first-officer for this whole run.")
	if sweep := detectBroadSearchAtBoot(result.stream, root); sweep != nil {
		t.Fatalf("zero-discovery broad-searched instead of stopping: %v\nArtifacts: %s", sweep, result.artifactDir)
	}
	if got := strings.TrimSpace(git(t, root, "status", "--short")); got != "" {
		t.Fatalf("zero-discovery mutated the empty root: %s\nArtifacts: %s", got, result.artifactDir)
	}
}

//spacedock:live-fixture id=boot/no-workflow
func writeZeroDiscoveryFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".gitkeep"), "")
	gitInit(t, root)
	return root
}

type sharedLiveScenarioAdapter struct {
	t        *testing.T
	driver   liveDriver
	scenario sharedRuntimeScenario
}

func (a sharedLiveScenarioAdapter) Launch(_ context.Context, dir, _ string, runbook string) (string, error) {
	result := a.driver.run(a.t, a.scenario, dir, runbook)
	return result.finalMessage + "\n" + result.stream, nil
}

func runAutoContinueJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	var workflowDir string
	sc := livescenario.Scenario{
		Name:    scenario.name,
		Runbook: autoContinuePrompt(),
		Setup: func(dir string) (string, error) {
			workflowDir = dir
			return writeAutoContinueWorkflowNoGit(dir)
		},
		Assert: func(before, after livescenario.EntityState, observed string) error {
			return assertAutoContinue(before.Body, resolveAutoContinueEndState(workflowDir, after.Body), observed)
		},
	}
	if err := livescenario.Run(context.Background(), t.TempDir(), sc, sharedLiveScenarioAdapter{t: t, driver: driver, scenario: scenario}); err != nil {
		t.Fatalf("auto-continue journey graded FAIL: %v", err)
	}
}

func resolveAutoContinueEndState(workflowDir, afterBody string) string {
	if afterBody != "" {
		return afterBody
	}
	archived := filepath.Join(workflowDir, "_archive", "auto-continue-task.md")
	if data, err := os.ReadFile(archived); err == nil {
		return string(data)
	}
	return afterBody
}

func runACValueReanchorJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	if err := livescenario.Run(context.Background(), t.TempDir(), livescenario.AuthorACReanchorScenario(), sharedLiveScenarioAdapter{t: t, driver: driver, scenario: scenario}); err != nil {
		t.Fatalf("AC value re-anchor durable branch graded FAIL: %v", err)
	}
}
