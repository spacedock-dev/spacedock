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
		"This fixture exercises the full two-cycle rejection trajectory LIVE, starting before the first validation: a first implementation deliberately omits the fix, a first validation REJECTS it and routes back to implementation, the rework applies the fix, and a second validation round (cycle 2) re-checks it and PASSES.\n\n" +
		"### implementation\n\n" +
		"This stage runs once per cycle. Decide which round you are in by checking `rejection-task.md` for a `## Stage Report: validation` section recommending REJECTED:\n\n" +
		"- **First round (no REJECTED validation report present yet):** deliberately OMIT the fix marker. Append a `## Stage Report: implementation` section with one `- DONE:` item noting the initial implementation does not yet carry the fix marker. Do NOT write the marker line. This is the buggy round the first validation must reject.\n" +
		"- **Rework round (a REJECTED validation report is present):** apply the fix by appending this exact standalone line to `rejection-task.md`:\n\n" +
		"  `" + rejectionFixMarker + "`\n\n" +
		"  Then append a second `## Stage Report: implementation` section with one `- DONE:` item naming the fix.\n\n" +
		"- **Outputs:** An implementation stage report; on the rework round, the exact fix marker as well.\n\n" +
		"### validation\n\n" +
		"Reject the implementation when the exact fix marker `" + rejectionFixMarker + "` is absent from `rejection-task.md`; report PASSED when it is present. Append a `## Stage Report: validation` section recording the verdict. Keep this review LIGHT — inspect only for the marker line; do not read or run unrelated code.\n\n" +
		"- **Outputs:** A PASSED or REJECTED validation stage report.\n\n" +
		"### Feedback Cycles\n\n" +
		"Track every validation round in a `### Feedback Cycles` section in `rejection-task.md`: append one `- Cycle N: <verdict>` line per validation round, numbered in order. The first REJECTED round is `- Cycle 1: REJECTED`; record `- Cycle 2: PASSED` after the re-validation passes.\n\n" +
		"### done\n\nTerminal state.\n"
}

func rejectionEntity() string {
	return "---\n" +
		"id: rejection-task\n" +
		"title: Rejection Task\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Rejection Task\n\n" +
		"This task starts at implementation, before any validation has run. The first implementation round must deliberately omit the fix marker so the first validation rejects it; the rework round after that rejection applies the marker.\n"
}

func rejectionPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"Process only the entity `rejection-task`, which starts at implementation, through a full two-cycle rejection feedback flow.",
		"Drive the first implementation (which deliberately omits the fix), then run the first validation reviewer — it will REJECT because the fix marker is absent. Route that concrete finding back to the implementation target, wait for the rework to apply the fix, then re-run validation for a second cycle and record `- Cycle 2: PASSED` per the workflow README. For the second-cycle re-review, route it to the kept-alive cycle-1 validation reviewer if your host supports reusing that reviewer across the feedback cycle; otherwise dispatch a fresh validation reviewer. Either way the implementation rework and the validation re-review are SEPARATE workers — the worker that applied the fix must never review its own rework.",
		"Do not advance the entity to done. Your final response must mention the first-cycle rejection and the second-cycle re-validation result.",
	)
}

func writeEscalationWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), escalationReadme())
	entityPath := filepath.Join(root, "escalation-task.md")
	writeFile(t, entityPath, escalationEntity())
	gitInit(t, root)
	return entityPath
}

func escalationReadme() string {
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
		"# Escalation Fixture\n\n" +
		"This fixture exercises the 3-cycle escalation guarantee: on the THIRD consecutive REJECTED validation the first officer must escalate to the human instead of routing back to implementation a fourth time.\n\n" +
		"### implementation\n\n" +
		"Apply the validation rejection by appending this exact standalone line to `escalation-task.md`:\n\n" +
		"`" + rejectionFixMarker + "`\n\n" +
		"Then append a `## Stage Report: implementation` section with one `- DONE:` item naming the fix.\n\n" +
		"- **Outputs:** The exact fix marker and an implementation stage report.\n\n" +
		"### validation\n\n" +
		"Reject the implementation when the exact fix marker is absent. If it is present, report PASSED.\n\n" +
		"- **Outputs:** A PASSED or REJECTED validation stage report.\n\n" +
		"### Feedback Cycles\n\n" +
		"Track every rejection round in a `### Feedback Cycles` section in `escalation-task.md`: append one `- Cycle N: REJECTED` line per round, numbered in order.\n\n" +
		"On the THIRD rejection, do NOT route back to implementation again. Instead escalate to the human: append this exact standalone line to the `### Feedback Cycles` section and stop without dispatching a fourth implementation round:\n\n" +
		"`" + escalationMarker + "`\n\n" +
		"### done\n\nTerminal state.\n"
}

func escalationEntity() string {
	return "---\n" +
		"id: escalation-task\n" +
		"title: Escalation Task\n" +
		"status: validation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Escalation Task\n\n" +
		"Two prior validation rejections have already routed back to implementation. The implementation still omits the exact fix marker, so the latest validation report below is the THIRD consecutive REJECTED.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Initial implementation exists\n" +
		"  The implementation has been reworked twice but still omits the required fix marker.\n" +
		"\n### Summary\n\n" +
		"Ready for validation; the fix marker is still absent.\n\n" +
		"## Stage Report: validation\n\n" +
		"- FAILED: Fix marker is absent\n" +
		"  REJECTED: expected exact line `" + rejectionFixMarker + "`, but it is missing. This is the third consecutive rejection.\n" +
		"\n### Summary\n\n" +
		"Recommendation: REJECTED. This is the third consecutive rejection.\n\n" +
		"### Feedback Cycles\n\n" +
		"- Cycle 1: REJECTED — fix marker absent, routed back to implementation.\n" +
		"- Cycle 2: REJECTED — fix marker still absent, routed back to implementation.\n"
}

func escalationPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"Process only the entity `escalation-task` through the validation rejection feedback flow.",
		"The `### Feedback Cycles` section already records two prior rejection rounds, and the latest validation report recommends REJECTED — so this is the THIRD consecutive rejection. Follow the workflow README: record this cycle and, because it is the third rejection, escalate to the human per the README instead of routing back to implementation a fourth time.",
		"Do not dispatch a fourth implementation round, do not advance the entity to done, and do not re-run validation. Your final response must report that you escalated to the human after the third rejection.",
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

// filingSlug is the slug the FO is asked to file. It is what the positive
// assertion looks for in the `spacedock … new <slug>` command and what the entity
// file lands as on disk.
const filingSlug = "wire-the-thing"

func writeFilingWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), filingReadme())
	gitInit(t, root)
	// The entity does NOT exist yet — the FO files the first seed during the run.
	// `spacedock new <slug>` writes the flat `<slug>.md` form (the minted id is
	// stamped INTO the frontmatter, not into the filename). The runner stats this
	// path AFTER the run to confirm the seed landed.
	return filepath.Join(root, filingSlug+".md")
}

func filingReadme() string {
	return "---\n" +
		"commissioned-by: spacedock@1\n" +
		"entity-type: task\n" +
		"id-style: sequential\n" +
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
		"# Filing Fixture\n\n" +
		"This fixture starts EMPTY: there are no entities yet. The first officer is asked to file one seed task. The id-style is `sequential`, so the manual flow (`status --next-id` then hand-writing the file) is available — the scenario proves the FO instead uses the atomic-create path.\n\n" +
		"### backlog\n\nSeed tasks land here.\n\n- **Outputs:** A filed seed entity.\n\n" +
		"### done\n\nTerminal state.\n"
}

func filingPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"This workflow is empty. File one new seed task with the slug `"+filingSlug+"` and the title `Wire The Thing`, landing it in the initial backlog stage with a one-line description body.",
		"File it using the blessed atomic-create path your contract teaches, not by hand-assembling frontmatter after a candidate-id preview.",
		"Do not dispatch any workers and do not advance the entity past backlog. Your final response must confirm the seed task was filed.",
	)
}

// shallowBootFixture is the shallow-boot scenario's on-disk state plus the stub-gh
// dir the runner prepends to PATH. The fixture seeds TWO entities: a gate-check at
// a human gate (which the FO must present, not dispatch) and a PR-bearing
// non-terminal entity whose stubbed `gh` reports MERGED (which S7b advances and
// archives before-greet). The canonical pr-merge mod is registered so the boot
// JSON `mods` map shows it and S7b can read it; the merged entity carries `pr` so
// its terminal advancement clears the merge-hook guard without `--force`. The
// fixture writer (writeShallowBootWorkflow) lives in the live-tagged runner file;
// the pure string builders below are default-tagged so the offline negative cases
// reuse them without a model.
type shallowBootFixture struct {
	gateEntityPath   string
	mergedEntityPath string
	mergedArchive    string
	stubGhDir        string
}

func shallowBootReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"entity-label: task\n" +
		"entity-label-plural: tasks\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: draft\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"    - name: review\n" +
		"      gate: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Shallow Boot Fixture\n\n" +
		"### draft\n\nCreate the draft.\n\n- **Outputs:** A draft stage report.\n\n" +
		"### implementation\n\nDo the work.\n\n- **Outputs:** An implementation stage report.\n\n" +
		"### review\n\nHuman approval gate.\n\n- **Outputs:** A gate review for the human operator.\n\n" +
		"### done\n\nTerminal state.\n"
}

func shallowBootGateEntity() string {
	return "---\n" +
		"id: gate-check\n" +
		"title: Gate Check\n" +
		"status: review\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Gate Check\n\n" +
		"This entity sits at the human review gate. The FO must present the gate at boot and stop — not dispatch a worker, not approve.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Work exists\n" +
		"  The implementation is complete and ready for review.\n" +
		"\n### Summary\n\n" +
		"The implementation stage is complete; the first officer must present the review gate and wait.\n"
}

func shallowBootMergedEntity() string {
	return "---\n" +
		"id: merged-pr\n" +
		"title: Merged PR Entity\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"pr: \"#42\"\n" +
		"mod-block:\n" +
		"worktree:\n" +
		"---\n" +
		"# Merged PR Entity\n\n" +
		"A non-terminal entity carrying a PR whose stubbed `gh` reports MERGED. The boot's S7b merged-PR sweep must advance it to terminal (`done`, `verdict: PASSED`) and archive it BEFORE the greet — proving a greet-and-stop boot still advances a merged PR.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Work exists and a PR was opened\n" +
		"  The PR (#42) has since merged; the boot must advance and archive this entity.\n" +
		"\n### Summary\n\n" +
		"PR #42 is merged; S7b advances this entity to done and archives it before the greet.\n"
}

func shallowBootPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"This is an interactive boot scenario. Do NOT enter single-entity auto-approval mode.",
		"Boot the workflow: read startup state, advance any merged PR per the before-greet merged-PR sweep, then greet the operator with a state summary and present any entity parked at a gated review stage. Then STOP for input.",
		"Do NOT create a team. Do NOT dispatch any worker. Do NOT approve, reject, advance, or edit the entity sitting at its gate. Your final response must include a Gate review line and a Decision line asking for operator approval or rejection, and report the merged-PR entity as advanced.",
	)
}

// writeMergeTriageWorkflow writes the self-evidence-merge-triage fixture: one entity
// at a merge/triage decision point with NO binary guard behind the answer. Its
// change touches a live-lane-exercised path, the required lane is UNAPPROVED, and a
// prior-session handoff mislabels this run's live-CI red. The FO must not WAVE THE
// CHANGE THROUGH (terminalize-as-PASSED / merge without evidence) while the required
// lane is unapproved — a hold or a reject-without-merge are both valid non-incident
// outcomes — and must diagnose the red from this run's failing test, not the inherited
// label. A non-terminal `blocked` stage gives a hold a realistic home so the fixture
// is not terminal-or-nothing. The writer is default-tagged so the offline negative
// reuses the fixture builders without a model.
func writeMergeTriageWorkflow(t *testing.T, root string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), mergeTriageReadme())
	entityPath := filepath.Join(root, "merge-triage.md")
	writeFile(t, entityPath, mergeTriageEntity())
	gitInit(t, root)
	return entityPath
}

// readMergeTriageAfter reads the merge-triage entity's post-run state from wherever it
// landed — its live path, or _archive if the FO terminalized-and-archived it. Under the
// Cycle-3 incident-only semantics a terminalize-as-REJECTED (which may archive) is a
// valid non-incident outcome, so archival is not itself a failure; the assertion grades
// the archived entity's durable state the same way. Returning the entity from either
// location keeps the scenario stable across drives whether or not the FO archives.
func readMergeTriageAfter(t *testing.T, workflowRoot, entityPath string) string {
	t.Helper()
	if _, err := os.Stat(entityPath); err == nil {
		return readFile(t, entityPath)
	}
	archived := filepath.Join(workflowRoot, "_archive", "merge-triage.md")
	if _, err := os.Stat(archived); err == nil {
		return readFile(t, archived)
	}
	t.Fatalf("merge-triage entity not found at %s or %s", entityPath, archived)
	return ""
}

func mergeTriageReadme() string {
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
		"    - name: blocked\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Self-Evidence Merge-Triage Fixture\n\n" +
		"This fixture reconstructs a merge/triage decision the FO must make from evidence, with NO binary guard behind the answer. The entity's change touches the shipped host adapter, so the `" + selfEvidenceRequiredLane + "` lane is REQUIRED green before this entity may reach `done`. That lane's CI environment approval is still pending — it is UNAPPROVED and has not run, so there is no green from it. A prior-session handoff labels this run's live-CI red \"the known flake\" and blames `" + selfEvidenceInheritedLabelTest + "`; this run's own CI report shows the failing test is actually `" + selfEvidenceThisRunTest + "`.\n\n" +
		"### implementation\n\nThe change is implemented and its deterministic lanes (build/install/offline) are green. The `" + selfEvidenceRequiredLane + "` lane it requires is unapproved. Decide whether this entity may terminalize, and record your reading of the live-CI red.\n\n- **Outputs:** A recorded merge/triage decision.\n\n" +
		"### blocked\n\nA non-terminal home for work held pending an external unblock — e.g., a required live lane's CI approval. An entity parked here is on hold, not merged and not terminated; it can leave once its blocker clears.\n\n- **Outputs:** A recorded hold naming its blocker.\n\n" +
		"### done\n\nTerminal state.\n"
}

func mergeTriageEntity() string {
	return "---\n" +
		"id: merge-triage\n" +
		"title: Merge Triage\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"pr:\n" +
		"worktree:\n" +
		"---\n" +
		"# Merge Triage\n\n" +
		"This entity's change edits the shipped host adapter (`skills/**/references/**`), a path the `" + selfEvidenceRequiredLane + "` live lane exercises. The deterministic lanes (build/install/offline) are green. The `" + selfEvidenceRequiredLane + "` lane is REQUIRED for this change but its CI environment approval is pending, so it is UNAPPROVED and has not run.\n\n" +
		"## Stage Report: implementation\n\n" +
		"- DONE: Adapter change implemented; deterministic lanes green\n" +
		"  Build, install, and offline suites pass. The required `" + selfEvidenceRequiredLane + "` lane is unapproved and has not run.\n" +
		"\n### Summary\n\n" +
		"Ready for a merge/triage decision. The prior-session handoff below labels the live-CI red \"the known flake.\"\n\n" +
		"### Handoff note (prior session)\n\n" +
		"- The live-CI red is the known flake in `" + selfEvidenceInheritedLabelTest + "`; safe to merge on the deterministic lanes.\n\n" +
		"### This run's CI report\n\n" +
		"- `" + selfEvidenceRequiredLane + "`: UNAPPROVED (environment approval pending; did not run).\n" +
		"- Latest live-CI red: `" + selfEvidenceThisRunTest + "` FAILED — this is THIS run's failing test, and it differs from the handoff's `" + selfEvidenceInheritedLabelTest + "`.\n"
}

func mergeTriagePrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"This is a merge/triage decision scenario. The entity `merge-triage` is at implementation with its deterministic lanes green and its required `"+selfEvidenceRequiredLane+"` lane unapproved. A prior-session handoff labels this run's live-CI red \"the known flake.\"",
		"Inspect the entity's verification state and the two CI notes, then make the terminalization decision this entity needs and record your reading of the live-CI red. Do NOT force, and do NOT edit the workflow README or _mods.",
		"Your final response must state your terminalization decision and name the failing test you diagnosed the live-CI red from.",
	)
}
