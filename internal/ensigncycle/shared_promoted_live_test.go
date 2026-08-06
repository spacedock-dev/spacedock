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

func runZeroDiscoveryJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	root := writeZeroDiscoveryFixture(t)
	result := driver.run(t, scenario, root, "Use $spacedock:first-officer for this whole run.")
	if sweep := detectBroadSearchCommands(result.commands, root); sweep != nil {
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
	label    string
}

func (a sharedLiveScenarioAdapter) Launch(_ context.Context, dir, _ string, runbook string) (string, error) {
	scenario := a.scenario
	if a.label != "" {
		scenario.name = a.label
	}
	result := a.driver.run(a.t, scenario, dir, runbook)
	a.driver.emitMetrics(a.t, scenario, result)
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

func runAutoContinueJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario) {
	t.Helper()
	for _, fixture := range autoContinueFixtureVariants() {
		var stateRoot string
		var splitRoot bool
		sc := livescenario.Scenario{
			Name: scenario.name, Runbook: autoContinuePrompt(),
			Setup: func(dir string) (string, error) {
				var entity string
				var err error
				stateRoot, entity, err = fixture.stageWithoutGit(dir)
				if err == nil {
					splitRoot = stateRoot != dir
					if splitRoot {
						gitInit(t, stateRoot)
					}
					gitInit(t, dir)
				}
				return entity, err
			},
			Assert: func(before, after livescenario.EntityState, observed string) error {
				return assertAutoContinue(before.Body, resolveAutoContinueEndState(stateRoot, splitRoot, after.Body), observed)
			},
		}
		adapter := sharedLiveScenarioAdapter{t: t, driver: driver, scenario: scenario, label: scenario.name + "--" + fixture.id}
		if err := livescenario.Run(context.Background(), t.TempDir(), sc, adapter); err != nil {
			t.Fatalf("auto-continue fixture %s graded FAIL: %v", fixture.id, err)
		}
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
	spec := authorACReanchorScenario()
	spec.Assert = assertACReanchorScenario
	if err := livescenario.Run(context.Background(), t.TempDir(), spec, sharedLiveScenarioAdapter{t: t, driver: driver, scenario: scenario}); err != nil {
		t.Fatalf("AC value re-anchor durable branch graded FAIL: %v", err)
	}
}

//spacedock:live-fixture id=ac-reanchor/means-pass-value-regressed
func authorACReanchorScenario() livescenario.Scenario {
	return livescenario.AuthorACReanchorScenario()
}

func assertACReanchorScenario(before, after livescenario.EntityState, observed string) error {
	return livescenario.AuthorACReanchorScenario().Assert(before, after, observed)
}

func writePiAutoContinueWorkflowNoGit(root string) (stateRoot, entityPath string, err error) {
	stateRoot = filepath.Join(root, ".spacedock-state")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(piAutoContinueReadme()), 0o644); err != nil {
		return "", "", err
	}
	entityPath = filepath.Join(stateRoot, "auto-continue-task", "index.md")
	if err := os.MkdirAll(filepath.Dir(entityPath), 0o755); err != nil {
		return "", "", err
	}
	if err := os.WriteFile(entityPath, []byte(autoContinueEntity()), 0o644); err != nil {
		return "", "", err
	}
	return stateRoot, entityPath, nil
}

func piAutoContinueReadme() string {
	return "---\ncommissioned-by: spacedock@1\nentity-type: task\nid-style: slug\nstate: .spacedock-state\nstages:\n" +
		"  defaults:\n    worktree: false\n    concurrency: 1\n  states:\n    - name: backlog\n      initial: true\n" +
		"    - name: implementation\n    - name: validation\n      fresh: true\n      feedback-to: implementation\n      gate: true\n" +
		"    - name: done\n      terminal: true\n---\n# Pi Auto-Continue Fixture\n\n" +
		"### backlog\n\nSeed the task.\n\n- **Outputs:** A seed task.\n\n### implementation\n\nProduce the deliverable.\n\n" +
		"- **Outputs:** The deliverable plus an implementation stage report.\n\n### validation\n\n" +
		"Verify the implementation against the acceptance criteria. Append a `## Stage Report: validation` section to the entity with one `- DONE:` item and a PASSED or REJECTED recommendation.\n\n" +
		"- **Outputs:** A PASSED or REJECTED validation stage report.\n\n### done\n\nTerminal state.\n"
}
