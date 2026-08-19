---
name: fo-dispatch-recovery
description: "Claude dispatch failure recovery — Break-Glass Manual Dispatch (the manual `Agent()` template) and Context Budget Failure/Dead Ensign Handling (the budget-unavailable stderr conditions, the recovery clause, dead-ensign bookkeeping). Read ONLY at its resident triggers inside `claude-fo-dispatch.md` — a non-zero or unavailable `spacedock dispatch build` (Break-Glass); or a budget-fail/zombie/dead-ensign replacement dispatch, including the dispatch-failure retry rung's fresh `-retry` re-dispatch under dead-ensign handling (Context Budget) — never at boot, never on the happy path."
user-invocable: false
---

# First Officer Dispatch Recovery (Claude)

The two Claude dispatch exception bodies, each read only at its own failure trigger in `claude-fo-dispatch.md` — never at boot, never on a session where dispatch never fails.

## Break-Glass Manual Dispatch

The resident trigger line already covers the first action (report the helper failure — command, exit code, stderr — to the captain before proceeding). The dispatch mode selected before `dispatch build` remains authoritative: do not probe another transport, retry in the other mode, or turn a selected bare dispatch into a named worker. Populate `{numbered checklist}` with the output of `«dispatch.checklist»(entity, stage)`; do not rebuild its rules here. In either arm, include `model="{effective_model}"` only when an effective model is set.

For selected bare mode, use this blocking call. Omit `name`, `team_name`, and `run_in_background` entirely; Claude's observable tool stream may preserve that omission or normalize it to `run_in_background=false`, and both mean blocking bare dispatch. Never pass or accept `run_in_background=true`. Omit the completion-message block because the blocking return is the completion signal:
```
Agent(
    subagent_type="spacedock:ensign",  // override with the stage's agent: field when the workflow README names one
    description="{entity title}: {stage}",
    model="{effective_model}",
    prompt="## First action\n\nBefore anything else, invoke your operating contract:\n\n    Skill(skill=\"spacedock:ensign\")\n\nThis loads the shared ensign discipline (stage-report format, background-task polling, worktree ownership, completion signal protocol). Do not paraphrase; call the tool.\n\nYou are working on: {entity title}\n\nStage: {stage}\n\n### Stage definition:\n\n{copy stage subsection from README verbatim}\n\nRead the entity file at {entity_file_path}.\n\n### Completion checklist\n\n{numbered checklist}\n\n### Summary\n{brief description of what was accomplished}\n\n### Stage report\n\nAppend a Stage Report section at the end of the entity file (per the shared-core Stage Report Protocol). Use the title `Stage Report: {stage}`. Account for every checklist item above with a `- DONE:` / `- SKIPPED:` / `- FAILED:` entry. Use the checklist item text verbatim when possible."
)
```

For selected team mode, use this named background call. Omit `team_name`; retain the completion message to the single `team-lead` target:
```
Agent(
    subagent_type="spacedock:ensign",  // override with the stage's agent: field when the workflow README names one
    description="{entity title}: {stage}",
    name="{worker_key}-{slug}-{stage}",  // if this exceeds 64 chars, cap it the way `spacedock dispatch build` does: keep the {worker_key} prefix and -{stage} suffix and, on id-style: sd-b32, replace the slug with a fixed-length prefix of the entity id (id-less slug workflows truncate the slug head instead)
    run_in_background=true,
    model="{effective_model}",
    prompt="## First action\n\nBefore anything else, invoke your operating contract:\n\n    Skill(skill=\"spacedock:ensign\")\n\nThis loads the shared ensign discipline (stage-report format, background-task polling, worktree ownership, completion signal protocol). Do not paraphrase; call the tool.\n\nYou are working on: {entity title}\n\nStage: {stage}\n\n### Stage definition:\n\n{copy stage subsection from README verbatim}\n\nRead the entity file at {entity_file_path}.\n\n### Completion checklist\n\n{numbered checklist}\n\n### Summary\n{brief description of what was accomplished}\n\n### Stage report\n\nAppend a Stage Report section at the end of the entity file (per the shared-core Stage Report Protocol). Use the title `Stage Report: {stage}`. Account for every checklist item above with a `- DONE:` / `- SKIPPED:` / `- FAILED:` entry. Use the checklist item text verbatim when possible.\n\n### Completion Signal\n\nSendMessage(to=\"team-lead\", message=\"Done: {entity title} completed {stage}. Report written to {entity_file_path}.\")"
)
```
This is the concrete Claude form of fo-dispatch-core.md's Break-Glass template; the contract (what it omits, the conditional `model=` slot, "use only when the helper is unavailable") is stated there. The canonical enum the conditional slot draws from is the resident `## Context Budget` section of `claude-fo-dispatch.md` (already loaded alongside this skill at the first dispatch).

## Context Budget Failure and Dead Ensign Handling

**Budget-unavailable is fail-safe (never silent-reuse).** The probe exits non-zero with no `reuse_ok` field in three conditions; the FO treats every one identically — fresh-dispatch:
- **missing jsonl** — no `agent-*.jsonl` exists for the named member (stderr: `no subagent jsonl found for '{name}'`).
- **unreadable/empty jsonl** — the jsonl exists but carries no assistant entry with non-zero `usage` (stderr: `no assistant entries with usage in {path}`).
- **agent-not-in-team-config** — no team `config.json` lists a member with that name (stderr: `no team config found for member '{name}'`).
A non-zero exit with no `reuse_ok: true` means the FO never silent-reuses on an absent reading.

**Recovery clause** (only when replacing a prior ensign): The prior ensign was shut down due to context budget limits. Its worktree may hold uncommitted changes. Run `git status` and `git diff` first; commit legitimate WIP or reset broken changes.

**Dead ensign handling:**

- `SendMessage(shutdown_request)` is cooperative — do NOT send to dead or unresponsive ensigns.
- Track dead ensigns in session memory; do not route work to dead names.
- Fresh-dispatch under a `-cycleN` suffix when replacing a zombie ensign.
- The post-dispatch config check does NOT detect zombies — zombies pass it. Session memory is the authoritative dead-vs-alive tracker.
