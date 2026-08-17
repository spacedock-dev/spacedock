//go:build live

package ensigncycle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/livescenario"
	"github.com/spacedock-dev/spacedock/internal/status"
)

func runFullEnsignCycleJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario, build func(*testing.T) string, assert func(*testing.T, string, string) bool) {
	t.Helper()
	root := build(t)
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
	if !assert(t, root, "make-it-work") {
		t.Fatalf("full ensign cycle has no path-scoped entity commit; artifacts: %s", result.artifactDir)
	}
	driver.emitMetrics(t, scenario, result)
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

func runZeroDiscoveryJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario, build func(*testing.T) string, assert func([]string, string) error) {
	t.Helper()
	root := build(t)
	result := driver.run(t, scenario, root, "Use $spacedock:first-officer for this whole run.")
	if sweep := assert(result.commands, root); sweep != nil {
		t.Fatalf("zero-discovery broad-searched instead of stopping: %v\nArtifacts: %s", sweep, result.artifactDir)
	}
	if got := strings.TrimSpace(git(t, root, "status", "--short")); got != "" {
		t.Fatalf("zero-discovery mutated the empty root: %s\nArtifacts: %s", got, result.artifactDir)
	}
	driver.emitMetrics(t, scenario, result)
}

func detectBroadSearchCommands(commands []string, root string) error {
	for _, command := range commands {
		if signature, broad := broadSweepCommand(command, filepath.Clean(root)); broad {
			return fmt.Errorf("FO broad-searched the filesystem at boot: %s in %q", signature, command)
		}
	}
	return nil
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
	result   *liveResult
}

func (a sharedLiveScenarioAdapter) Launch(_ context.Context, dir, _ string, runbook string) (string, error) {
	result := a.driver.run(a.t, a.scenario, dir, runbook)
	if a.result != nil {
		*a.result = result
	} else {
		a.driver.emitMetrics(a.t, a.scenario, result)
	}
	return result.finalMessage + "\n" + result.stream, nil
}

type autoContinueFixtureVariant struct {
	id              string
	stageWithoutGit func(string) (string, string, error)
}

//spacedock:live-fixture id=auto-continue/single-root,auto-continue/split-root
func autoContinueFixtureVariants() []autoContinueFixtureVariant {
	return []autoContinueFixtureVariant{
		{id: "auto-continue/single-root", stageWithoutGit: func(root string) (string, string, error) {
			entity, err := writeAutoContinueWorkflowNoGit(root)
			return root, entity, err
		}},
		{id: "auto-continue/split-root", stageWithoutGit: writePiAutoContinueWorkflowNoGit},
	}
}

func initializeAutoContinueFixtureGit(t *testing.T, workflowRoot, stateRoot string) {
	t.Helper()
	if stateRoot != workflowRoot {
		gitInit(t, stateRoot)
		branch, err := status.StateBranch(workflowRoot)
		if err != nil {
			t.Fatal(err)
		}
		git(t, stateRoot, "branch", "-m", branch)
	}
	gitInit(t, workflowRoot)
}

func runAutoContinueJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario, build func() []autoContinueFixtureVariant, assert func(string, string, string) error) {
	t.Helper()
	for _, fixture := range build() {
		var stateRoot, entityPath string
		var splitRoot bool
		sc := livescenario.Scenario{
			Name: scenario.name, Runbook: autoContinuePrompt(),
			Setup: func(dir string) (string, error) {
				var entity string
				var err error
				stateRoot, entity, err = fixture.stageWithoutGit(dir)
				if err == nil {
					splitRoot = stateRoot != dir
					initializeAutoContinueFixtureGit(t, dir, stateRoot)
				}
				entityPath = entity
				return entity, err
			},
			Assert: func(before, after livescenario.EntityState, observed string) error {
				return assert(before.Body, resolveAutoContinueEndState(stateRoot, splitRoot, after.Body), observed)
			},
		}
		fixtureScenario := scenario
		fixtureScenario.name += "--" + fixture.id
		var result liveResult
		adapter := sharedLiveScenarioAdapter{t: t, driver: driver, scenario: fixtureScenario, result: &result}
		var semantic []error
		if err := livescenario.Run(context.Background(), t.TempDir(), sc, adapter); err != nil {
			semantic = append(semantic, durableSemantic("auto-continue-state", err))
		}
		semantic = append(semantic, durableSemantic("validation-worker-lifecycle",
			assertAutoContinueDispatchEvidence(t, driver.lifecycleStream(t, result), stateRoot, entityPath)))
		finishLiveScenario(t, driver, fixtureScenario, result, semantic...)
	}
}

func resolveAutoContinueEndState(stateRoot string, splitRoot bool, after string) string {
	if worktree := autoContinueWorktreeDir(after); !splitRoot && worktree != "" {
		if data, err := os.ReadFile(filepath.Join(stateRoot, worktree, "auto-continue-task.md")); err == nil {
			return string(data)
		}
	}
	if after != "" {
		return after
	}
	for _, archived := range []string{filepath.Join(stateRoot, "_archive", "auto-continue-task.md"), filepath.Join(stateRoot, "_archive", "auto-continue-task", "index.md")} {
		if data, err := os.ReadFile(archived); err == nil {
			return string(data)
		}
	}
	return after
}

func runACValueReanchorJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario, build func() livescenario.Scenario, assert func(livescenario.Scenario, livescenario.EntityState, livescenario.EntityState, string) error) {
	t.Helper()
	spec := build()
	authored := spec
	spec.Assert = func(before, after livescenario.EntityState, observed string) error {
		return assert(authored, before, after, observed)
	}
	if err := livescenario.Run(context.Background(), t.TempDir(), spec, sharedLiveScenarioAdapter{t: t, driver: driver, scenario: scenario}); err != nil {
		t.Fatalf("AC value re-anchor durable branch graded FAIL: %v", err)
	}
}

//spacedock:live-fixture id=ac-reanchor/means-pass-value-regressed
func authorACReanchorScenario() livescenario.Scenario {
	return livescenario.AuthorACReanchorScenario()
}

func assertACReanchorScenario(spec livescenario.Scenario, before, after livescenario.EntityState, observed string) error {
	return spec.Assert(before, after, observed)
}
