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
	t             *testing.T
	driver        liveDriver
	scenario      sharedRuntimeScenario
	artifactLabel string
}

func (a sharedLiveScenarioAdapter) Launch(_ context.Context, dir, _ string, runbook string) (string, error) {
	launchScenario := a.scenario
	if a.artifactLabel != "" {
		launchScenario.name = a.artifactLabel
	}
	result := a.driver.run(a.t, launchScenario, dir, runbook)
	return result.finalMessage + "\n" + result.stream, nil
}

type autoContinueFixtureVariant struct {
	id              string
	stageWithoutGit func(root string) (stateRoot, entityPath string, err error)
}

func autoContinueFixtureVariants() []autoContinueFixtureVariant {
	return []autoContinueFixtureVariant{
		{
			id: "auto-continue/single-root",
			stageWithoutGit: func(root string) (string, string, error) {
				entityPath, err := writeAutoContinueWorkflowNoGit(root)
				return root, entityPath, err
			},
		},
		{id: "auto-continue/split-root", stageWithoutGit: writePiAutoContinueWorkflowNoGit},
	}
}

func stageAutoContinueFixture(t *testing.T, root string, fixture autoContinueFixtureVariant) (stateRoot, entityPath string) {
	t.Helper()
	stateRoot, entityPath, err := fixture.stageWithoutGit(root)
	if err != nil {
		t.Fatal(err)
	}
	if stateRoot != root {
		gitInit(t, stateRoot)
	}
	gitInit(t, root)
	return stateRoot, entityPath
}

func runAutoContinueJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	for _, fixture := range autoContinueFixtureVariants() {
		fixture := fixture
		var stateRoot string
		var splitRoot bool
		sc := livescenario.Scenario{
			Name:    scenario.name,
			Runbook: autoContinuePrompt(),
			Setup: func(dir string) (string, error) {
				var entityPath string
				stateRoot, entityPath = stageAutoContinueFixture(t, dir, fixture)
				splitRoot = stateRoot != dir
				return entityPath, nil
			},
			Assert: func(before, after livescenario.EntityState, observed string) error {
				return assertAutoContinue(before.Body, resolveAutoContinueEndState(stateRoot, splitRoot, after.Body), observed)
			},
		}
		adapter := sharedLiveScenarioAdapter{
			t:             t,
			driver:        driver,
			scenario:      scenario,
			artifactLabel: scenario.name + "--" + fixture.id,
		}
		if err := livescenario.Run(context.Background(), t.TempDir(), sc, adapter); err != nil {
			t.Fatalf("auto-continue fixture %s graded FAIL: %v", fixture.id, err)
		}
	}
}

func resolveAutoContinueEndState(stateRoot string, splitRoot bool, afterBody string) string {
	if worktree := autoContinueWorktreeDir(afterBody); !splitRoot && worktree != "" {
		active := filepath.Join(stateRoot, worktree, "auto-continue-task.md")
		if data, err := os.ReadFile(active); err == nil {
			return string(data)
		}
	}
	if afterBody != "" {
		return afterBody
	}
	for _, archived := range []string{
		filepath.Join(stateRoot, "_archive", "auto-continue-task.md"),
		filepath.Join(stateRoot, "_archive", "auto-continue-task", "index.md"),
	} {
		if data, err := os.ReadFile(archived); err == nil {
			return string(data)
		}
	}
	return afterBody
}

func autoContinueWorktreeDir(body string) string {
	for _, line := range strings.Split(body, "\n") {
		value, found := strings.CutPrefix(line, "worktree:")
		if !found {
			continue
		}
		value = filepath.Clean(strings.TrimSpace(value))
		if value == "." || value == ".." || filepath.IsAbs(value) || strings.HasPrefix(value, ".."+string(filepath.Separator)) {
			return ""
		}
		return value
	}
	return ""
}

func runACValueReanchorJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	if err := livescenario.Run(context.Background(), t.TempDir(), livescenario.AuthorACReanchorScenario(), sharedLiveScenarioAdapter{t: t, driver: driver, scenario: scenario}); err != nil {
		t.Fatalf("AC value re-anchor durable branch graded FAIL: %v", err)
	}
}
