package ensigncycle

import (
	"path/filepath"
	"testing"
)

// writePiNonSelfDescribingSmokeWorkflow creates a split-root smoke workflow
// whose implementation stage-def names only the real work (append a marker
// line) — no "stage report" mention — so the worker's stage-report format
// source is the embedded dispatch body block, not the stage-def.
//
//spacedock:live-fixture id=pi/non-self-describing-smoke
func writePiNonSelfDescribingSmokeWorkflow(t *testing.T) (workflowRoot, stateRoot, entityPath string) {
	t.Helper()
	workflowRoot = t.TempDir()
	stateRoot = filepath.Join(workflowRoot, ".spacedock-state")
	writeFile(t, filepath.Join(workflowRoot, "README.md"), piNonSelfDescribingSmokeReadme())
	entityPath = filepath.Join(stateRoot, "pi-nonsd-smoke", "index.md")
	writeFile(t, entityPath, piNonSelfDescribingSmokeEntity())
	gitInit(t, workflowRoot)
	gitInit(t, stateRoot)
	return workflowRoot, stateRoot, entityPath
}

func piNonSelfDescribingSmokeReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"state: .spacedock-state\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: implementation\n" +
		"      initial: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Pi Non-Self-Describing Smoke\n\n" +
		"### implementation\n\n" +
		"Append the live Pi smoke marker line `PI-NONSD-SMOKE-MARKER` as a standalone line to the entity file, then commit only the entity path in the state checkout.\n\n" +
		"- **Outputs:** The marker line present in the entity file and a path-scoped state commit.\n\n" +
		"### done\n\nTerminal state.\n"
}

func piNonSelfDescribingSmokeEntity() string {
	return "---\n" +
		"id: pi-nonsd-smoke\n" +
		"title: Pi Non-Self-Describing Smoke\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Pi Non-Self-Describing Smoke\n\n" +
		"This entity is mutated only by the Pi subagent non-self-describing live smoke.\n"
}
