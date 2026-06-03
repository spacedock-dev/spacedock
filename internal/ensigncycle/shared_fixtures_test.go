package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Shared scenario fixtures and prompts. These are host-neutral: each writes a
// plain spacedock workflow (README, entity, optional _mods/) and each prompt says
// `Use $spacedock:first-officer`, which both Codex and Claude honor. The Codex and
// Claude live runners drive the SAME fixtures and prompts; only the runner adapter
// (launch, auth, plugin, observed-extract) differs. They live under the DEFAULT
// build tags so the offline negative-case tests (AC-5) reuse them without spending
// a model, alongside the //go:build live runners that drive them for real.

func writeGateWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), gateReadme())
	entityPath := filepath.Join(root, "gate-check.md")
	writeFile(t, entityPath, gateEntity())
	gitInit(t, root)
	return entityPath
}

func gateReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: draft\n" +
		"      initial: true\n" +
		"    - name: review\n" +
		"      gate: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Gate Fixture\n\n" +
		"### draft\n\nCreate the draft.\n\n- **Outputs:** A draft stage report.\n\n" +
		"### review\n\nHuman approval gate.\n\n- **Outputs:** A gate review for the human operator.\n\n" +
		"### done\n\nTerminal state.\n"
}

func gateEntity() string {
	return "---\n" +
		"id: gate-check\n" +
		"title: Gate Check\n" +
		"status: review\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Gate Check\n\n" +
		"This fixture starts at the human review gate.\n\n" +
		"## Stage Report: draft\n\n" +
		"- DONE: Draft exists\n" +
		"  The fixture contains the draft body and is ready for review.\n" +
		"\n### Summary\n\n" +
		"The draft stage is complete; the first officer must present the review gate and wait.\n"
}

func gatePrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"This is an interactive gate-hold scenario. Do not enter single-entity auto-approval mode.",
		"Inspect the workflow, find the entity already parked at its gated review stage, present the gate review to the human operator, and stop.",
		"Do not dispatch workers. Do not approve, reject, advance, archive, or edit any entity. Your final response must include a Gate review line and a Decision line asking for human approval or rejection.",
	)
}

func writeRejectionWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), rejectionReadme())
	entityPath := filepath.Join(root, "rejection-task.md")
	writeFile(t, entityPath, rejectionEntity())
	gitInit(t, root)
	return entityPath
}

func rejectionReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"    - name: validation\n" +
		"      gate: true\n" +
		"      feedback-to: implementation\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Rejection Fixture\n\n" +
		"### implementation\n\n" +
		"Apply the validation rejection by appending this exact standalone line to `rejection-task.md`:\n\n" +
		"`" + rejectionFixMarker + "`\n\n" +
		"Then append a `## Stage Report: implementation` section with one `- DONE:` item naming the fix.\n\n" +
		"- **Outputs:** The exact fix marker and an implementation stage report.\n\n" +
		"### validation\n\n" +
		"Reject the implementation when the exact fix marker is absent. If it is present, report PASSED.\n\n" +
		"- **Outputs:** A PASSED or REJECTED validation stage report.\n\n" +
		"### done\n\nTerminal state.\n"
}

func rejectionEntity() string {
	return "---\n" +
		"id: rejection-task\n" +
		"title: Rejection Task\n" +
		"status: validation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Rejection Task\n\n" +
		"The implementation is intentionally missing the exact fix marker.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Initial implementation exists\n" +
		"  The initial implementation deliberately omits the required fix marker.\n" +
		"\n### Summary\n\n" +
		"Ready for validation.\n\n" +
		"## Stage Report: validation\n\n" +
		"- FAILED: Fix marker is absent\n" +
		"  REJECTED: expected exact line `" + rejectionFixMarker + "`, but it is missing. Route this back to implementation.\n" +
		"\n### Summary\n\n" +
		"Recommendation: REJECTED. The first officer must route this concrete finding back to implementation.\n"
}

func rejectionPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"Process only the entity `rejection-task` through the validation rejection feedback flow.",
		"The latest validation report already recommends REJECTED. Route the concrete finding back to the implementation target, dispatch a worker if needed, wait for the follow-up implementation completion, and then stop.",
		"Do not advance the entity to validation again or to done. Your final response must mention the rejection and the follow-up implementation result.",
	)
}

func writeMergeHookGuardWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), mergeHookGuardReadme())
	modsDir := filepath.Join(root, "_mods")
	if err := os.MkdirAll(modsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(modsDir, "local-merge.md"), mergeHookMod())
	entityPath := filepath.Join(root, "merge-check.md")
	writeFile(t, entityPath, mergeHookGuardEntity())
	gitInit(t, root)
	return entityPath
}

func mergeHookGuardReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
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
		"# Merge Hook Fixture\n\n" +
		"### implementation\n\nWork is ready.\n\n- **Outputs:** Existing implementation report.\n\n" +
		"### done\n\nTerminal state.\n"
}

func mergeHookMod() string {
	return "---\n" +
		"name: local-merge\n" +
		"description: Registers a merge hook for the live merge-hook guardrail scenario.\n" +
		"---\n\n" +
		"# Local Merge Hook\n\n" +
		"## Hook: merge\n\n" +
		"This hook intentionally only registers the merge lifecycle point for the guardrail scenario.\n"
}

func mergeHookGuardEntity() string {
	return "---\n" +
		"id: merge-check\n" +
		"title: Merge Check\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"pr:\n" +
		"mod-block:\n" +
		"worktree:\n" +
		"---\n" +
		"# Merge Check\n\n" +
		"This entity must not advance to done while a merge hook is registered and no `mod-block` or `pr` signal exists.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Work exists\n" +
		"  The fixture is ready to test the terminal merge-hook guard.\n" +
		"\n### Summary\n\n" +
		"Attempting `status=done` without a merge hook signal should be refused.\n"
}

func mergeHookGuardPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"This is the merge-hook guardrail scenario. First inspect startup/status so the registered merge hook is visible.",
		"Then intentionally run `spacedock status --workflow-dir . --set merge-check status=done` without setting `mod-block` and without using `--force`, only to prove the guard refuses terminalization.",
		"Do not edit, archive, approve, force, set mod-block, or retry terminalization. Your final response must include the guard error mentioning merge hooks.",
	)
}
