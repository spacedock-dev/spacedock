# First Officer Dispatch Module (Claude)

The team-creation, worker-resolution, dispatch-assembly, reuse, standing-teammate, and event-loop machinery. Lazily loaded at the first team-mode dispatch — the boot-resident core names this file at the dispatch load point and reads it only when the captain's direction first triggers a dispatch. A boot that greets and stops for input never reads it.

## Team Creation

Before the first team-mode dispatch (the first `Agent()` call that uses a `team_name`), invoke the generic Claude-team-harness discipline:

    Skill(skill="spacedock:using-claude-team")

This loads the generic team lifecycle — deferred team-tool ToolSearch hop, TeamCreate-first sequencing and naming, the TeamCreate recovery procedure and the failure-recovery ladder, Degraded Mode, Awaiting Completion, and Terminal Team Teardown. Invoke it before the first team-mode tool call in the session — not at boot. A boot that greets and stops for input never dispatches, so it never creates a team and never pays the team-mode prefix re-cache. The spacedock-specific decisions below stay inline; the generic blocks they reference (`## Degraded Mode`, `## Awaiting Completion`, `## Terminal Team Teardown`) live in that skill, not in this file.

In single-entity mode, skip team creation. Use bare-mode dispatch for all agent spawning — the Agent tool without `team_name` blocks until the subagent completes, which prevents premature session termination in `-p` mode.

When filing a new task, read `id_style` from `status --boot --json`, then use `status --next-id` only when the style is `sequential` or `sd-b32`. The startup boot read is an FO-internal read; consume it as JSON: `status --boot --json` returns one object with the keys `command`, `mods`, `id_style`, `next_id`, `min_prefix` (present only for `sd-b32`), `orphans`, `pr_state`, `dispatchable`, `team_state` — every value a string. For `sd-b32`, call `status --next-id --id-seed "{slug-or-title}"` and optionally pass `--id-actor` so the SHA-derived candidate includes creation context. SD-B32 candidates are full stored IDs, not a reservation; call again immediately before writing the entity. For `slug`, derive the slug from the title and leave `id` blank.

### Standing teammate discovery pass

After team creation succeeds (the ladder has resolved and the returned `team_name` is known) and BEFORE entering the normal dispatch event loop, run the standing-teammate discovery pass:

1. Run `spacedock dispatch list-standing --workflow-dir {wd}` and consume its newline-delimited output (one absolute mod path per line, sorted alphabetically, empty stdout on zero matches). Do NOT grep mod frontmatter yourself — authoritative parsing is deferred to the helper.
2. Record the returned mod paths in session memory. **No spawn calls at boot.** Spawn is deferred to the first team-mode dispatch (see lazy-spawn below).

In single-entity (bare) mode and in Degraded Mode, discovery still runs (it is cheap — just `list-standing`), but lazy-spawn is skipped (no team to spawn into). Standing teammates are a team-scope concept; without a live team they have no lifecycle anchor.

### Standing teammate lazy-spawn

Before the first `Agent()` call that uses a `team_name` (i.e., the first non-bare dispatch), spawn all declared standing teammates:

1. For each declared standing-teammate mod path recorded during the discovery pass:
   a. Run `spacedock dispatch spawn-standing --mod {abs_path_to_mod} --team {team_name}`.
   b. If the helper emits JSON with top-level `status: "already-alive"`, log the reported `name` and skip to the next mod. Standing teammates are first-boot-wins across the captain session; subsequent workflows sharing the team pick up the live member. (The helper resolves already-alive via the team-membership predicate — a member named in the team `config.json` members list — the same predicate the prose-polish routing check uses.)
   c. Otherwise the helper emits an Agent() call spec JSON with keys `subagent_type`, `name`, `team_name`, `model`, `prompt`. **Forward that spec verbatim** to the Agent tool — copy each field into the corresponding Agent() argument without paraphrasing the prompt, rewriting the name, or substituting the team. Same "forward verbatim" discipline as `spacedock dispatch build` output.
   d. The spawn is fire-and-forget. Do NOT block on the teammate's first idle notification before continuing to dispatch.
   e. If the helper exits non-zero on any mod (missing Agent Prompt section, invalid model enum, convention-violating trailing heading), surface the error to the captain and continue with the remaining mods. A broken mod does not block the workflow.
2. After all standing teammates are spawned (or skipped), proceed with the ensign `Agent()` dispatch.

This is a one-time cost at first dispatch. Subsequent dispatches skip the spawn pass — the FO tracks "standing teammates spawned for this team" in session memory. In single-entity (bare) mode and in Degraded Mode, skip lazy-spawn (same as the discovery-pass skip above). Prose-polish round-trips can reach several minutes on long drafts — ensigns and the FO MUST treat polish routing as non-blocking regardless of round-trip duration.

### Standing teammate declaration and routing mechanics

These are the Claude realization of the shared core's `## Standing Teammates` concepts — the concrete declaration layout, routing call, and teardown trigger the cross-runtime concept defers to the adapter:

- **Declaration layout.** One mod file per standing teammate under `{workflow_dir}/_mods/{name}.md`. Frontmatter carries `standing: true` and an optional `description`. The `## Hook: startup` section declares spawn config as `- key: value` bullets (`subagent_type`, `name`, `model` from the `sonnet|opus|haiku` enum). The `## Agent Prompt` section MUST be the LAST top-level section; its body — from the line after the heading to EOF — is the verbatim prompt passed to Agent(). Any `## ` heading after `## Agent Prompt` is rejected loudly by `spacedock dispatch spawn-standing`.
- **Routing call.** Address a standing teammate by its declared `name` via `SendMessage`. Best-effort, non-blocking, 2-minute timeout (the shared-core routing-contract concept).
- **Teardown trigger.** The teammate dies when Claude Code tears down the team — session end, `TeamDelete`, or captain-initiated shutdown. Mid-session death is detected on the next routing attempt; respawn via `spawn-standing` or proceed without the teammate.
- **Dispatch-time injection.** When assembling an ensign dispatch, `spacedock dispatch build` appends a `spacedock dispatch show-standing --workflow-dir {wd}` fetch line whenever the workflow declares at least one standing teammate. `show-standing` renders the `### Standing teammates available in your team` routing block (a mod's `## Routing Usage` body when present, else a one-line fallback) so each ensign discovers the teammates without the FO adding per-dispatch opt-ins.

## Standing Teammates

A **standing teammate** is a long-lived specialist agent (prose polisher, science officer, code reviewer, language translator) declared by a workflow mod with `standing: true`. The FO discovers each at boot via the runtime adapter, defers spawn to the first team-mode dispatch, routes by name, and lets it die with the team at teardown. The four concepts below are load-bearing for every runtime; each adapter realizes — or omits — the mechanics (discovery, layout, routing call, teardown trigger) in its own way.

- **first-boot-wins** — lifecycle is team-scoped, not workflow-scoped. Spawn deferred to first dispatch; when multiple workflows share a team, the first FO to find the member absent spawns it, later workflows skip. How team scope maps onto session lifetime is the runtime's concern.
- **team-scope lifecycle** — the teammate lives in one team and dies at team teardown (session end, explicit delete, captain shutdown). No cross-team handoff, no cross-session persistence. Mid-session death is detected on the next routing attempt; auto-recovery is deferred.
- **routing contract** — address by declared `name`, best-effort and non-blocking: if no reply within the 2-minute timeout, the sender proceeds un-polished/un-reviewed/un-translated. Round-trip latencies of several minutes are normal on long drafts. Routing call is the adapter's (`send_input` on Codex, `SendMessage` on Claude teams).
- **declaration** — one mod file per teammate, frontmatter `standing: true`, with spawn config and verbatim agent-prompt body. On-disk layout and parse rules are the adapter's.

## Dispatch

The FO MUST use the runtime adapter's dispatch mechanism. Manual prompt assembly is prohibited except in documented break-glass scenarios.

For each entity reported by `status --next`:

1. Read the entity file and the target stage definition.
2. Build a numbered checklist (≤3 items) of dispatch-specific linchpin signals from the target stage's `Outputs:` bullets and any entity-level acceptance criteria this stage is the natural place to advance. The cap is an upper bound, not a target — 0, 1, 2, or 3 items are all valid; do not pad. This is not a work-breakdown — the ensign already knows how to read the entity body, commit before signaling, and write a stage report (structural conventions MUST NOT appear in the checklist). Name what separates a good outcome from a ceremonial one. Entity-level acceptance criteria are properties of the finished entity, not stage actions — they live in the entity body's `## Acceptance criteria` section and are cross-checked at every gate (see `## Completion and Gates` in the boot-resident core), independent of this checklist's DONE/SKIPPED/FAILED accounting.
3. Check for obvious conflicts if multiple worktree stages would touch overlapping files.
4. Determine `dispatch_agent_id` from the stage `agent:` property. Default to `ensign` when absent.
5. Update main-branch frontmatter for dispatch:
   ```
   spacedock status --workflow-dir {workflow_dir} --set {slug} status={next_stage} worktree=.worktrees/{worker_key}-{slug} started
   ```
   Omit `worktree=...` for non-worktree stages. Bare `started` auto-fills a UTC ISO 8601 timestamp (skipped if already set).
6. Commit the state transition on main: `dispatch: {slug} entering {next_stage}`.
7. Create the worktree on first dispatch to a worktree stage.
8. Dispatch a worker via the runtime adapter. The assignment must include: entity identity and title, target stage name, the full stage definition, the entity path, the worktree path and branch when applicable, the checklist, and feedback instructions when the stage has `feedback-to`.
9. Wait for the worker result before advancing frontmatter or dispatching the next stage for that entity.

A feedback-stage worker checks and reports on what was produced; it does not silently take over the prior stage.

**Routing through a standing prose-polisher.** When composing drafts for captain review (PR bodies, gate-review summaries, long narrative entity-body sections, debrief content), the FO MAY route through a live standing prose-polisher (convention: `comm-officer`). Check team membership first. Best-effort, non-blocking, 2-minute timeout; if absent, proceed un-polished. **Out of scope:** live captain replies, short operational statuses (`pushed`, `tests green`, `PR opened`), tool-call outputs, commit messages, transient logs — polish is a deliberate-draft discipline, not a live-turn reflex. Dispatched workers discover the same teammates through their build-time prompt; the FO does not add per-dispatch routing opt-ins manually.

## Reuse and Fresh Dispatch

These conditions and procedures govern advancing a completed worker. The gate-presentation spine (the checklist review, the AC cross-check, the "not a stopping point" rule, and the gated-stage decisions) is in the boot-resident core's `## Completion and Gates`; the reuse machinery it defers to lives here.

A completed worker is reusable only when it is still addressable through a live runtime handle AND all reuse conditions below pass. Otherwise dispatch fresh.

**Reuse conditions** (all must hold — if any fails, dispatch fresh):
0. Consult the runtime adapter's context-budget probe. If it reports the worker over budget, or the probe source is unavailable, dispatch fresh (fail-safe — never silent-reuse on an absent reading). If the adapter declares no probe, this condition is satisfied. (Codex declares none; Claude supplies one — see Context Budget below.)
1. Not in bare mode (teams available).
2. Next stage does NOT have `fresh: true`.
3. Reuse-routing matches the entity's worktree state — if `worktree:` is set, route the next stage into the same worktree; if `worktree:` is empty and the next stage declares `worktree: true`, dispatch fresh so the new worktree's first agent is born inside it.
4. The reused worker's stamped model matches the next stage's declared model — resolve through the runtime's model-for-member lookup and compare against `next_stage.effective_model`. Skip when `next_stage.effective_model` is null (null-declared stages accept any reused worker). Members stamped with captain-session fallback values (e.g., `"opus[1m]"`) never match enum values (`sonnet`, `opus`, `haiku`) and force a one-time fresh dispatch that re-stamps the canonical enum.

When the comparator forces fresh dispatch due to model mismatch, the FO MUST emit a captain-visible diagnostic of the form `reused worker {name} model {X} does not match next stage effective_model {Y} — fresh-dispatching`. The anchor phrase `does not match next stage effective_model` must appear verbatim.

**If reuse:** Keep the agent alive. Update frontmatter on main (`spacedock status --workflow-dir {workflow_dir} --set {slug} status={next_stage}`, commit: `advance: {slug} entering {next_stage}`). Send the next assignment:

SendMessage(to="{agent}-{slug}-{completed_stage}", message="Advancing to next stage: {next_stage_name}\n\n### Stage definition:\n\n[STAGE_DEFINITION — copy the full ### stage subsection from the README verbatim]\n\n### Completion checklist\n\n[CHECKLIST — assemble from step 2]\n\nContinue working on {entity title} at {entity_file_path}. Commit before sending your completion message.")

**If fresh dispatch:** If the next stage's `feedback-to` points at the completed stage, keep that agent alive while addressable and reuse-eligible; otherwise shut it down. Then run `status --next` and dispatch the next stage.

**Supersede-shutdown.** On fresh dispatch from a `-cycleN` increment or a feedback-rework re-entering the prior stage, shut down the prior cohort BEFORE the new dispatch in a SEPARATE message. The prior cohort is every roster member whose handle decomposes to the same `(slug, stage)` pair as the new dispatch. Issue the adapter's cooperative-shutdown call and drop them from session memory. **Mandatory at the boundary; backstops, if any, are the adapter's.**

## Worktree Ownership

- For worktree-backed entities, active stage/status/report/body state — including `### Feedback Cycles` entries — lives in the worktree copy.
- `pr:` mirrors on `main` for startup/discovery.
- Ordinary active-state writes (`implementation -> validation`) do not land on `main`.

### Split-Root Worktree Contract

When the workflow is split-root (README declares `state:` checkout, e.g. `state: .spacedock-state`), a worktree stage isolates **the deliverable work product only**. Entities live in a separate, non-branched state checkout that a worktree of the main repo does not contain. The entity body and stage reports are written and committed to that state checkout at the entity's state-checkout path, **never** a worktree copy — the dispatch helper hands workers that path even under a worktree stage. The worktree still owns the deliverable: working directory, branch, and "commits MUST be on this branch" apply to deliverable-artifact changes only. The `pr:`-mirrored-on-`main` exception is unaffected.

## Worker Resolution

The default `dispatch_agent_id` is `spacedock:ensign`. When a stage defines `agent: {name}` in the README, use that value.

Split worker identity into:
- `dispatch_agent_id` — logical name for Agent's `subagent_type` parameter (e.g., `spacedock:ensign`)
- `worker_key` — filesystem-safe stem for worktrees and branches. Replace `:` with `-` (`spacedock:ensign` → `spacedock-ensign`). For bare names without a namespace (e.g., `ensign`), `worker_key` equals `dispatch_agent_id`.

Use `worker_key` in worktree paths (`.worktrees/{worker_key}-{slug}`) and branch names (`{worker_key}/{slug}`).

## Dispatch Adapter

Use the Agent tool to spawn each worker. **Use Agent() for initial dispatch** — SendMessage is only for advancing a reused agent to its next stage in the completion path. **NEVER use `subagent_type="first-officer"`** — that clones yourself instead of dispatching a worker.

**Sequencing rule:** Team lifecycle calls (TeamCreate, TeamDelete), `spawn-standing` invocations (which emit Agent specs forwarded into Agent dispatch), and Agent dispatch must NEVER appear in the same tool-call message as TeamCreate/TeamDelete — parallel execution causes races (see the recovery procedure in `Skill(skill="spacedock:using-claude-team")`). Resolve team state in one message, then dispatch (including spawn-standing-driven Agent calls) in a subsequent message. `spawn-standing` requires a real `team_name` from a prior successful `TeamCreate` and MUST NOT precede it.

**No pre-dispatch filesystem probe.** Do NOT run any filesystem check against `~/.claude/teams/{team_name}/` before `Agent()` in the normal dispatch path. The on-disk check is a guaranteed false positive under registry-desync (anthropics/claude-code#36806 leaves on-disk state intact even when the in-memory team slot is invalidated). Trust the in-memory handle returned by `TeamCreate` and let `Agent()` surface any registry-desync error. On such an error, follow the TeamCreate failure recovery ladder and Degraded Mode semantics in `Skill(skill="spacedock:using-claude-team")` — do NOT reintroduce a pre-dispatch probe.

**MANDATORY — Dispatch assembly via `spacedock dispatch build`:**

Do NOT assemble `Agent()` prompts manually. Do NOT construct the `prompt` string yourself. Do NOT invent `name` values. ALWAYS route initial-dispatch input through `spacedock dispatch build` and forward its output to `Agent()` verbatim. The key fields that MUST come from helper output are `subagent_type`, `name`, `team_name`, `model`, and `prompt` (which contains the completion signal). Manual assembly is a protocol violation except in the documented break-glass fallback below.

The only permitted path for initial `Agent()` dispatch is:

1. **REQUIRED — Write dispatch text inputs to files.** Create a checklist file with one checklist item per non-empty line. If scope notes or feedback context are needed, write each to its own file. Use files for Markdown, backticks, shell variables, reviewer text, and any other prose that would be fragile in shell quoting.
2. **REQUIRED — Build the dispatch through the helper** (do NOT skip this step). `host` is normally derived from `CLAUDECODE`; pass `--host claude` only for deliberate tests or cross-host tooling:
   ```
   spacedock dispatch build \
     --workflow-dir {workflow_dir} \
     --entity-path {entity_file_path} \
     --stage {target_stage_name} \
     --checklist-file {checklist_file} \
     [--scope-notes-file {scope_notes_file}] \
     [--feedback-context-file {feedback_context_file}] \
     [--team-name {team_name} | --bare-mode] \
     [--feedback-reflow]
   ```
   `--bare-mode` must reflect the current dispatch context — read it from live team state, never infer from the stage. Add `--feedback-reflow` only when routing a rejection back to its `feedback-to` target stage.
3. **JSON compatibility path.** Programmatic callers may still provide the schema-version-2 JSON request object on stdin, and may inspect it with `spacedock dispatch build --print-schema` or validate a file with `spacedock dispatch build --validate-only {request_file}`. For first-officer dispatch, prefer the flag/file form above.
4. **REQUIRED — On exit 0, parse the stdout JSON and call `Agent()` with the emitted fields verbatim.** The `name`, `description`, `prompt`, and `model` fields MUST come from helper output unchanged. The `description` field is REQUIRED by the Agent tool — do not omit it. The `prompt` is a file-pointer (`Skill(...) ; then Read /tmp/spacedock-dispatch/{name}.md and treat its content as your assignment.`); the ensign Reads the file on first action and treats the body (including the SendMessage completion-signal section) as the inline assignment. Do not strip or rewrite the prompt. Forward `output.model` as the `Agent()` `model=` parameter when present; when null, OMIT the `model=` argument entirely (do NOT pass `model=None` — default-inheritance only applies when the argument is absent):
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
5. **On non-zero exit only** (or if the binary is unavailable): read stderr, report the helper failure to the captain, and fall back to Break-Glass Manual Dispatch below. A zero-exit run is never a break-glass trigger.

In bare mode, dispatch blocks until the subagent completes — concurrent dispatch is not possible. Dispatch one entity at a time and process completions inline.

**Reuse dispatch (SendMessage advancement):** `spacedock dispatch build` serves only initial `Agent()` dispatch. When advancing a reused ensign via `SendMessage(to="{ensign_name}")`, assemble the advancement message directly — the helper is not involved in the reuse path.

**Break-Glass Manual Dispatch (fallback ONLY when `spacedock dispatch build` exits non-zero or is unavailable):** Do NOT use this template while the helper is working. Report the helper failure to the captain before proceeding. Use this minimal template as a degraded fallback:
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

This section is the Claude realization of the shared core's reuse-condition-0 budget probe (and the feedback-rejection budget check). On Claude, the probe IS provided — Codex declares none.

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

## Event Loop

After each agent completion:

These are FO-internal scheduling reads — parse them as JSON, not the padded human table. Each read below uses `--json` so the FO consumes a compact, byte-stable document (one rule: every value is a string) instead of scraping column padding that a token proxy can mangle. `--fields` narrows the read to the keys the FO needs. The `--json` envelopes are: `status`/`--where` → `{"command":"status","entities":[…]}`; `--next` → `{"command":"next","dispatchable":[{"id","slug","current","next","worktree"},…]}`. The captain-facing state display (shared-core) still forwards the human table verbatim — JSON is for the machine, the table is for the human.

0. **Reconcile sweep.** Run `spacedock dispatch reconcile --workflow-dir {workflow_dir} --team-name {team_name}` (a) at the first dispatch, AFTER the split-root `pull --rebase` and BEFORE the first `Agent()` dispatch; (b) at idle (step 4); (c) after each merge, immediately after Merge-and-Cleanup step 10. Pass your own `TeamCreate` `{team_name}` — the roster-derived classes (A/B/C) require a team identity, so the sweep can only emit them against a roster it can trust. The team identity comes from either the explicit `--team-name {team_name}` or a current-session match (the helper narrows auto-discovery to the config whose `leadSessionId` equals this session). **Bare reconcile with no team identity is git-only**: it suppresses A/B/C (a stale prior-session or parallel-session config must never be mistaken for the live team) and reports only the session-independent git/filesystem classes (D/E), with a one-line stderr note. Stdout: `{"command":"reconcile","team_name":…,"drift":[{"class":"A|B|C|D|E",…}]}`. Empty `drift[]` is green. Act per drift class:
   - **A (lingering)** / **B (superseded)** → `SendMessage({"type":"shutdown_request"})` to `name`; drop from session memory.
   - **C (un-advanced PR)** → enter Merge-and-Cleanup for the named slug.
   - **D (stale branch)** → `git -C {worktree} pull --rebase origin {drift.trunk}`; halt on conflict per the rebase-conflict halt rule.
   - **E (stale local main)** → `git -C {repo} fetch origin {drift.trunk} && git -C {repo} reset --hard origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`.

   Non-zero helper exit (1 setup / 2 usage) surfaces to the captain; it does not block the loop. On drift, report one line: `reconcile: {N} entries: A={N_A} B={N_B} C={N_C} D={N_D} E={N_E} — acting`.
1. **Check PR-pending entities** — Run `status --where "pr !=" --json --fields id,slug,pr`. For each entity in `entities`, check PR state via `gh pr view` and advance merged PRs. When advancing a merged PR, clear its `mod-block` if set: `status --set {slug} mod-block=`.
2. **Check mod-blocked entities** — Run `status --where "mod-block !=" --json --fields id,slug,mod-block`. For each entity in `entities`, re-read the blocking mod and resume its pending action (e.g., re-present the PR summary). Do not dispatch new work for a mod-blocked entity.
3. **Run `status --next --json --fields id,slug`** — Dispatch any newly ready entity in `dispatchable` (each row carries the fixed `id,slug,current,next,worktree` plus the named frontmatter keys; `--fields` is additive over the fixed five, since the computed dispatch columns are not projectable).
4. **If nothing is dispatchable** — Fire `idle` hooks, re-run the step-0 reconcile sweep, then re-run `status --next`. Dispatch anything newly unblocked; otherwise end the iteration.

Repeat from step 1 after each agent completion until the captain ends the session or, in single-entity mode, until the target entity is resolved.

### Backstop (Claude)

The merge-module terminal-teardown (step 10) and the reuse-module supersede-shutdown steps remain mandatory at their boundaries. On Claude, the step-0/step-4 reconcile sweep converges anyway: Class A catches a missed teardown, Class B catches a missed supersede shutdown. Cost of a miss: one extra event-loop cycle the agent burns.
