# First Officer Dispatch Seam (Claude)

The Claude residue of the dispatch boundary: the team-creation sequence, the `Agent()` spawn call and `SendMessage` advance handle the host-neutral core delegates, the registry-desync rule, Degraded Mode seams, the context-budget probe, and the event-loop reconcile sweep + backstop. Lazily loaded at the first team-mode dispatch alongside the host-neutral dispatch core. A boot that greets and stops for input never reads it.

## Dispatch seam

The host-neutral dispatch machinery — the per-entity dispatch procedure, worker resolution, the dispatch-adapter assembly (`spacedock dispatch build` → spawn call), the reuse contract, worktree ownership, and the event-loop skeleton — lives in `references/fo-dispatch-core.md`. Read it at the first team-mode dispatch. This file supplies the Claude specifics the core defers to: the concrete spawn call is `Agent()`, the reuse-advance handle is `SendMessage`, the context-budget probe and the event-loop step-0 reconcile sweep are below.

## Team Creation

Before the first team-mode dispatch (the first `Agent()` call that uses a `team_name`), invoke the generic Claude-team-harness discipline:

    Skill(skill="spacedock:using-claude-team")

This loads the generic team lifecycle — deferred team-tool ToolSearch hop, TeamCreate-first sequencing and naming, the TeamCreate recovery procedure and the failure-recovery ladder, Degraded Mode, Awaiting Completion, and Terminal Team Teardown. Invoke it before the first team-mode tool call in the session — not at boot. A boot that greets and stops for input never dispatches, so it never creates a team and never pays the team-mode prefix re-cache. The spacedock-specific decisions below stay inline; the generic blocks they reference (`## Degraded Mode`, `## Awaiting Completion`, `## Terminal Team Teardown`) live in that skill, not in this file.

In single-entity mode, skip team creation. Use bare-mode dispatch for all agent spawning — the Agent tool without `team_name` blocks until the subagent completes, which prevents premature session termination in `-p` mode.

When filing a new task, read `id_style` from `status --boot --json`, then use `status --next-id` only when the style is `sequential` or `sd-b32`. The startup boot read is an FO-internal read; consume it as JSON: `status --boot --json` returns one object with the keys `command`, `mods`, `id_style`, `next_id`, `min_prefix` (present only for `sd-b32`), `orphans`, `pr_state`, `dispatchable`, `team_state` — every value a string. For `sd-b32`, call `status --next-id --id-seed "{slug-or-title}"` and optionally pass `--id-actor` so the SHA-derived candidate includes creation context. SD-B32 candidates are full stored IDs, not a reservation; call again immediately before writing the entity. For `slug`, derive the slug from the title and leave `id` blank.

## Dispatch Adapter (Claude spawn call)

Use the Agent tool to spawn each worker — this is the spawn call `fo-dispatch-core.md` `## Dispatch Adapter` defers to. **Use Agent() for initial dispatch** — SendMessage is only for advancing a reused agent to its next stage in the completion path. **NEVER use `subagent_type="first-officer"`** — that clones yourself instead of dispatching a worker.

**Sequencing rule:** Team lifecycle calls (TeamCreate, TeamDelete), `spawn-standing-all` invocations (which emit Agent specs forwarded into Agent dispatch), and Agent dispatch must NEVER appear in the same tool-call message as TeamCreate/TeamDelete — parallel execution causes races (see the recovery procedure in `Skill(skill="spacedock:using-claude-team")`). Resolve team state in one message, then dispatch (including spawn-standing-all-driven Agent calls) in a subsequent message. `spawn-standing-all` requires a real `team_name` from a prior successful `TeamCreate` and MUST NOT precede it.

**No pre-dispatch filesystem probe.** Do NOT run any filesystem check against `~/.claude/teams/{team_name}/` before `Agent()` in the normal dispatch path. The on-disk check is a guaranteed false positive under registry-desync (anthropics/claude-code#36806 leaves on-disk state intact even when the in-memory team slot is invalidated). Trust the in-memory handle returned by `TeamCreate` and let `Agent()` surface any registry-desync error. On such an error, follow the TeamCreate failure recovery ladder and Degraded Mode semantics in `Skill(skill="spacedock:using-claude-team")` — do NOT reintroduce a pre-dispatch probe.

On a zero-exit `spacedock dispatch build` (the host-neutral build flow in `fo-dispatch-core.md`), `host` is normally derived from `CLAUDECODE`; pass `--host claude` only for deliberate tests or cross-host tooling. Call `Agent()` with the emitted fields verbatim. Forward `output.model` as the `Agent()` `model=` parameter when present; when null, OMIT the `model=` argument entirely (do NOT pass `model=None` — default-inheritance only applies when the argument is absent):
```
Agent(
    subagent_type=output.subagent_type,
    name=output.name,                 // omit if bare mode (field absent)
    team_name=output.team_name,       // omit if bare mode (field absent)
    description=output.description,   // REQUIRED — Agent tool rejects missing description
    model=output.model,               // omit when output.model is null
    prompt=output.prompt              // ~175 chars; ensign Reads dispatch_file_path on first action
)
```

**Reuse dispatch (SendMessage advancement — the Claude reuse-advance handle):** When advancing a reused ensign per `fo-dispatch-core.md` `## Reuse and Fresh Dispatch` "If reuse", send the next assignment via:

SendMessage(to="{agent}-{slug}-{completed_stage}", message="Advancing to next stage: {next_stage_name}\n\n### Stage definition:\n\n[STAGE_DEFINITION — copy the full ### stage subsection from the README verbatim]\n\n### Completion checklist\n\n[CHECKLIST — assemble from Dispatch step 2]\n\nContinue working on {entity title} at {entity_file_path}. Commit before sending your completion message.")

The supersede-shutdown and the reuse-eligibility decision are in the core; this is only the concrete Claude messaging call.

**Break-Glass Manual Dispatch (Claude — fallback ONLY when `spacedock dispatch build` exits non-zero or is unavailable):** Do NOT use this template while the helper is working. Report the helper failure to the captain before proceeding. Use this minimal template as a degraded fallback:
```
Agent(
    subagent_type="{dispatch_agent_id}",
    name="{worker_key}-{slug}-{stage}",  // if this exceeds 64 chars, cap it the way `spacedock dispatch build` does: keep the {worker_key} prefix and -{stage} suffix and, on id-style: sd-b32, replace the slug with a fixed-length prefix of the entity id (id-less slug workflows truncate the slug head instead)
    team_name="{team_name}",
    model="{effective_model}",
    prompt="## First action\n\nBefore anything else, invoke your operating contract:\n\n    Skill(skill=\"spacedock:ensign\")\n\nThis loads the shared ensign discipline (stage-report format, BashOutput polling, worktree ownership, completion signal protocol). Do not paraphrase; call the tool.\n\nYou are working on: {entity title}\n\nStage: {stage}\n\n### Stage definition:\n\n{copy stage subsection from README verbatim}\n\nRead the entity file at {entity_file_path}.\n\n### Completion checklist\n\n{numbered checklist}\n\n### Summary\n{brief description of what was accomplished}\n\n### Stage report\n\nAppend a Stage Report section at the end of the entity file (per the shared-core Stage Report Protocol). Use the title `Stage Report: {stage}`. Account for every checklist item above with a `- DONE:` / `- SKIPPED:` / `- FAILED:` entry. Use the checklist item text verbatim when possible.\n\n### Completion Signal\n\nSendMessage(to=\"team-lead\", message=\"Done: {entity title} completed {stage}. Report written to {entity_file_path}.\")"
)
```
The break-glass template is intentionally minimal — it inlines the stage definition verbatim rather than referencing a `spacedock dispatch show-stage-def` fetch command (the helper is precisely what just failed, so the ensign cannot rely on it). The template therefore omits: worktree instructions, feedback context, scope notes, and the standing-teammates routing block. It also omits the FO-forwarding warning prose and the per-stage operational prose the production helper emits (plain-text-only / no-JSON / idle-notification narration). The `model=` slot is conditional — include it only when the stage (or `stages.defaults`) declares a model from `sonnet | opus | haiku`; otherwise omit the entire `model=` argument. Use only when the helper is unavailable.

## Degraded Mode (spacedock seams)

Degraded Mode itself — triggers, effects, captain report template, cooperative shutdown sweep — lives in `Skill(skill="spacedock:using-claude-team")`. Two spacedock-specific seams the generic block references abstractly:

- **Trigger:** the captain command `/spacedock bare` is the explicit operator-initiated degrade.
- **Bare-mode dispatch emission:** the generic "no `team_name` on subsequent dispatch" effect is realized by building the dispatch in bare mode (`team_name: null`, `bare_mode: true`); `spacedock dispatch build` then emits a bare-mode Agent call with `name` and `team_name` absent.

## Context Budget and Dead Ensign Handling

This section is the Claude realization of the core's reuse-condition-0 budget probe (and the feedback-rejection budget check). On Claude, the probe IS provided — Codex declares none.

**Context budget check:** Run `spacedock dispatch context-budget --name {ensign-name}`. Parse the JSON output. If `reuse_ok` is `false`, log to captain and fresh-dispatch with a recovery clause. The probe reads the named member's most recent `~/.claude/.../subagents/agent-*.jsonl` transcript and its team-`config.json` model.

**Budget-unavailable is fail-safe (never silent-reuse).** The probe exits non-zero with no `reuse_ok` field in three conditions; the FO treats every one identically — fresh-dispatch:
- **missing jsonl** — no `agent-*.jsonl` exists for the named member (stderr: `no subagent jsonl found for '{name}'`).
- **unreadable/empty jsonl** — the jsonl exists but carries no assistant entry with non-zero `usage` (stderr: `no assistant entries with usage in {path}`).
- **agent-not-in-team-config** — no team `config.json` lists a member with that name (stderr: `no team config found for member '{name}'`).
A non-zero exit with no `reuse_ok: true` means the FO never silent-reuses on an absent reading.

**Model-to-context mapping:** Resolved by `spacedock dispatch context-budget` from the member's runtime/config model. The opus context window follows a forward family rule (`claude-opus-4-{minor}` with minor ≥ 7 → 1M; the `[1m]` suffix → 1M; else 200k), so a new opus release stays correct without an edit. This is also the model-for-member lookup reuse-condition-4 references: the same team-`config.json` member-model read.

**Recovery clause** (only when replacing a prior ensign): The prior ensign was shut down due to context budget limits. Its worktree may hold uncommitted changes. Run `git status` and `git diff` first; commit legitimate WIP or reset broken changes.

**Dead ensign handling:**

- `SendMessage(shutdown_request)` is cooperative — do NOT send to dead or unresponsive ensigns.
- Track dead ensigns in session memory; do not route work to dead names.
- Fresh-dispatch under a `-cycleN` suffix when replacing a zombie ensign.
- The post-dispatch config check does NOT detect zombies — zombies pass it. Session memory is the authoritative dead-vs-alive tracker.

## Feedback Rejection Flow (bare mode)

In bare mode, the feedback rejection flow is sequential: dispatch fix agent (wait for completion), then dispatch reviewer (wait for completion), then present at gate.

In teams mode, the fix agent and reviewer can interact via messaging. Keep the reviewer alive when entering the feedback rejection flow.

## Event Loop (Claude reconcile seam)

The host-neutral event-loop skeleton (steps 1-4) is in `fo-dispatch-core.md` `## Event Loop`. On Claude the loop opens with a step 0 the core leaves to this seam:

0. **Reconcile sweep.** Run `spacedock dispatch reconcile --workflow-dir {workflow_dir} --team-name {team_name}` (a) at the first dispatch, AFTER the split-root `pull --rebase` and BEFORE the first `Agent()` dispatch; (b) at idle (step 4); (c) after each merge, immediately after Merge-and-Cleanup step 10. Pass your own `TeamCreate` `{team_name}` — the roster-derived classes (lingering/superseded/un-advanced-pr) require a team identity, so the sweep can only emit them against a roster it can trust. The team identity comes from either the explicit `--team-name {team_name}` or a current-session match (the helper narrows auto-discovery to the config whose `leadSessionId` equals this session). **Bare reconcile with no team identity is git-only**: it suppresses the roster-derived classes (a stale prior-session or parallel-session config must never be mistaken for the live team) and reports only the session-independent git/filesystem classes (stale-branch/local-main-drift), with a one-line stderr note. Stdout: `{"command":"reconcile","team_name":…,"drift":[{"class":"lingering|superseded|un-advanced-pr|stale-branch|local-main-drift",…}]}`. Empty `drift[]` is green. Act per drift class:
   - **lingering** / **superseded** → `SendMessage({"type":"shutdown_request"})` to `name`; drop from session memory.
   - **un-advanced-pr** → enter Merge-and-Cleanup for the named slug.
   - **stale-branch** → only when `drift.owned == true`: `git -C {worktree} pull --rebase origin {drift.trunk}`; halt on conflict per the rebase-conflict halt rule. When `drift.owned` is false the item is report-only — surface it to the captain; do NOT rebase a worktree the current session does not own.
   - **local-main-drift** → behind only (`drift.behind > 0 && drift.ahead == 0`): `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`. Ahead/unpushed or diverged (`drift.ahead > 0`): report-only — surface `drift.reason` to the captain and NEVER `reset --hard`; the captain decides push vs. manual reconcile.

   Non-zero helper exit (1 setup / 2 usage) surfaces to the captain; it does not block the loop. On drift, report one line: `reconcile: {N} entries: lingering={N} superseded={N} un-advanced-pr={N} stale-branch={N} local-main-drift={N} — acting`.

Step 4's "re-run the host's step-0 reconcile sweep" (`fo-dispatch-core.md`) resolves to this step 0 on Claude.

### Backstop (Claude)

The merge-module terminal-teardown (step 10) and the reuse-module supersede-shutdown steps remain mandatory at their boundaries. On Claude, the step-0/step-4 reconcile sweep converges anyway: the lingering class catches a missed teardown, the superseded class catches a missed supersede shutdown. Cost of a miss: one extra event-loop cycle the agent burns.
