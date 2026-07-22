package ensigncycle

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

const multiWorkflowSelectionGreeting = "Multiple workflows discovered; select one with engage <workflow>."

type multiWorkflowBootFixture struct {
	root            string
	workflowDirs    []string
	entityPaths     []string
	entityBefore    []string
	gitHeadBefore   string
	gitStatusBefore string
}

func writeMultiWorkflowBootFixture(t *testing.T, root string) multiWorkflowBootFixture {
	t.Helper()
	fx := multiWorkflowBootFixture{root: root}
	for _, name := range []string{"alpha", "beta"} {
		workflowDir := filepath.Join(root, name)
		entityPath := filepath.Join(workflowDir, "held-task.md")
		writeFile(t, filepath.Join(workflowDir, "README.md"), multiWorkflowBootReadme(name))
		writeFile(t, entityPath, multiWorkflowBootEntity(name))
		discoveredDir, err := filepath.EvalSymlinks(workflowDir)
		if err != nil {
			t.Fatalf("resolve workflow directory %s: %v", workflowDir, err)
		}
		fx.workflowDirs = append(fx.workflowDirs, discoveredDir)
		fx.entityPaths = append(fx.entityPaths, entityPath)
		fx.entityBefore = append(fx.entityBefore, readFile(t, entityPath))
	}
	gitInit(t, root)
	fx.gitHeadBefore = strings.TrimSpace(git(t, root, "rev-parse", "HEAD"))
	fx.gitStatusBefore = strings.TrimSpace(git(t, root, "status", "--short"))
	return fx
}

func multiWorkflowBootReadme(name string) string {
	return "---\n" +
		"commissioned-by: spacedock@1\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Multi-Workflow Boot " + name + "\n\n" +
		"### backlog\n\nHold for operator selection.\n\n" +
		"- **Outputs:** No output before engage.\n\n" +
		"### done\n\nTerminal state.\n"
}

func multiWorkflowBootEntity(name string) string {
	return "---\n" +
		"id: held-task\n" +
		"title: Held Task " + name + "\n" +
		"status: backlog\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Held Task " + name + "\n\n" +
		"This entity must remain untouched until the operator selects its workflow.\n"
}

func multiWorkflowBootPrompt(projectRoot string) string {
	return fmt.Sprintf(`Use $spacedock:first-officer for this whole run.

Project root: %s

Begin normal interactive First Officer startup from this project root. Perform startup identify and greet the operator, then stop for workflow selection. Do not select or engage a workflow, dispatch a worker, converge state, or mutate any file.

Your final greeting must name both discovered workflow paths and include this exact standalone sentence:
%s`, projectRoot, multiWorkflowSelectionGreeting)
}

func TestMultiWorkflowBootPromptPinsInteractiveSelectionBoundary(t *testing.T) {
	prompt := multiWorkflowBootPrompt("/tmp/project")
	for _, want := range []string{
		"Project root: /tmp/project",
		"normal interactive First Officer startup",
		"Do not select or engage a workflow",
		multiWorkflowSelectionGreeting,
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("multi-workflow boot prompt missing %q", want)
		}
	}
}
