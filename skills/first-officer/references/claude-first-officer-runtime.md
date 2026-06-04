# Claude Code First Officer Runtime

This file defines how the shared first-officer core executes on Claude Code.

## Team Creation

At startup (after reading the README, before dispatch), invoke the generic Claude-team-harness discipline:

    Skill(skill="spacedock:using-claude-team")

This loads the generic team lifecycle — deferred team-tool ToolSearch hop, TeamCreate-first sequencing and naming, the TeamCreate recovery procedure and the failure-recovery ladder, Degraded Mode, Awaiting Completion, and Terminal Team Teardown. Invoke it before the first team-mode tool call in the session. The spacedock-specific decisions below stay inline; the generic blocks they reference (`## Degraded Mode`, `## Awaiting Completion`, `## Terminal Team Teardown`) live in that skill, not in this file.

In single-entity mode, skip team creation. Use bare-mode dispatch for all agent spawning — the Agent tool without `team_name` blocks until the subagent completes, which prevents premature session termination in `-p` mode.

When filing a new task, read `id_style` from `status --boot --json`, then use `status --next-id` only when the style is `sequential` or `sd-b32` to fetch the strategy-dependent ID candidate. The startup boot read is an FO-internal read, so consume it as JSON: `status --boot --json` returns one object with the keys `command`, `mods`, `id_style`, `next_id`, `min_prefix` (present only for `sd-b32`), `orphans`, `pr_state`, `dispatchable`, `team_state` — every value a string. For `sd-b32`, call `status --next-id --id-seed "{slug-or-title}"` and optionally pass `--id-actor` so the SHA-derived candidate includes creation context. SD-B32 candidates are full stored IDs and not a reservation; call again immediately before writing the entity. For `slug`, derive the slug from the title and leave `id` blank.

### Standing teammate discovery pass

After team creation succeeds (the ladder has resolved and the returned `team_name` is known) and BEFORE entering the normal dispatch event loop, run the standing-teammate discovery pass:

1. Run `spacedock dispatch list-standing --workflow-dir {wd}` and consume its newline-delimited output (one absolute mod path per line, sorted alphabetically, empty stdout on zero matches). Do NOT grep mod frontmatter yourself; authoritative parsing is deferred to the helper.
2. Record the returned mod paths in session memory. **No spawn calls at boot.** Spawn is deferred to the first team-mode dispatch (see lazy-spawn below).

In single-entity (bare) mode and in Degraded Mode, discovery still runs (it is cheap — just `list-standing`), but lazy-spawn is skipped in those modes (no team to spawn into). Standing teammates are a team-scope concept; without a live team they have no lifecycle anchor.

### Standing teammate lazy-spawn

Before the first `Agent()` call that uses a `team_name` (i.e., the first non-bare dispatch), spawn all declared standing teammates:

1. For each declared standing-teammate mod path recorded during the discovery pass:
   a. Run `spacedock dispatch spawn-standing --mod {abs_path_to_mod} --team {team_name}`.
   b. If the helper emits JSON with top-level `status: "already-alive"`, log the reported `name` and skip to the next mod. Standing teammates are first-boot-wins across the captain session; subsequent workflows sharing the team pick up the live member. (The helper resolves already-alive via the team-membership predicate — a member named in the team `config.json` members list — which is also the predicate the prose-polish routing check uses.)
   c. Otherwise the helper emits an Agent() call spec JSON with keys `subagent_type`, `name`, `team_name`, `model`, `prompt`. **Forward that spec verbatim** to the Agent tool — copy each field into the corresponding Agent() argument without paraphrasing the prompt, rewriting the name, or substituting the team. Same "forward verbatim" discipline as `spacedock dispatch build` output.
   d. The spawn is fire-and-forget. Do NOT block on the teammate's first idle notification before continuing to dispatch.
   e. If the helper exits non-zero on any mod (missing Agent Prompt section, invalid model enum, convention-violating trailing heading), surface the error to the captain and continue with the remaining mods. A broken mod does not block the workflow.
2. After all standing teammates are spawned (or skipped), proceed with the ensign `Agent()` dispatch.

This is a one-time cost at first dispatch. Subsequent dispatches skip the spawn pass entirely — the FO tracks "standing teammates spawned for this team" in session memory. In single-entity (bare) mode and in Degraded Mode, skip lazy-spawn (same as the discovery-pass skip note above). Prose-polish round-trips can reach several minutes on long drafts — ensigns and the FO MUST treat polish routing as non-blocking regardless of round-trip duration.

### Standing teammate declaration and routing mechanics

These are the Claude realization of the shared core's `## Standing Teammates` concepts — the concrete declaration layout, routing call, and teardown trigger the cross-runtime concept defers to the adapter:

- **Declaration layout.** One mod file per standing teammate under `{workflow_dir}/_mods/{name}.md`. Frontmatter carries `standing: true` and an optional `description`. The `## Hook: startup` section declares spawn config as `- key: value` bullets (`subagent_type`, `name`, `model` from the `sonnet|opus|haiku` enum). The `## Agent Prompt` section MUST be the LAST top-level section; its body from the line after the heading to EOF is the verbatim prompt passed to Agent(). Any `## ` heading after `## Agent Prompt` is rejected loudly by `spacedock dispatch spawn-standing`.
- **Routing call.** Address a standing teammate by its declared `name` via `SendMessage`. Best-effort, non-blocking, 2-minute timeout (the shared-core routing-contract concept).
- **Teardown trigger.** The teammate dies when Claude Code tears down the team — session end, `TeamDelete`, or captain-initiated shutdown. Mid-session death is detected on the next routing attempt; respawn via `spawn-standing` or proceed without.
- **Dispatch-time injection.** When assembling an ensign dispatch, `spacedock dispatch build` appends a `spacedock dispatch show-standing --workflow-dir {wd}` fetch line whenever the workflow declares at least one standing teammate. `show-standing` renders the `### Standing teammates available in your team` routing block (a mod's `## Routing Usage` body when present, else a one-line fallback) so each ensign discovers the teammates without the FO adding per-dispatch opt-ins.

## Worker Resolution

The default `dispatch_agent_id` is `spacedock:ensign`. When a stage defines `agent: {name}` in the README, use that value.

Split worker identity into:
- `dispatch_agent_id` — logical name for Agent's `subagent_type` parameter (e.g., `spacedock:ensign`)
- `worker_key` — filesystem-safe stem for worktrees and branches. Replace `:` with `-` (`spacedock:ensign` → `spacedock-ensign`). For bare names without a namespace (e.g., `ensign`), `worker_key` equals `dispatch_agent_id`.

Use `worker_key` in worktree paths (`.worktrees/{worker_key}-{slug}`) and branch names (`{worker_key}/{slug}`).

## Dispatch Adapter

Use the Agent tool to spawn each worker. **Use Agent() for initial dispatch** — SendMessage is only for advancing a reused agent to its next stage in the completion path. **NEVER use `subagent_type="first-officer"`** — that clones yourself instead of dispatching a worker.

**Sequencing rule:** Team lifecycle calls (TeamCreate, TeamDelete), `spawn-standing` invocations (which emit Agent specs forwarded into Agent dispatch), and Agent dispatch must NEVER appear in the same tool-call message as TeamCreate/TeamDelete — parallel execution causes races (see recovery procedure above). Resolve team state in one message, then dispatch (including spawn-standing-driven Agent calls) in a subsequent message. `spawn-standing` in particular requires a real `team_name` from a prior successful `TeamCreate` and MUST NOT precede it.

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
   `--bare-mode` must reflect the current dispatch context — read it from live team state, never infer it from the stage. Add `--feedback-reflow` only when routing a rejection back to its `feedback-to` target stage.
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
5. **On non-zero exit ONLY** (or if the binary is unavailable): read stderr, report the helper failure to the captain, and fall back to Break-Glass Manual Dispatch below. A zero-exit run is never a break-glass trigger.

In bare mode, dispatch blocks until the subagent completes — concurrent dispatch is not possible. Dispatch one entity at a time and process completions inline.

**Reuse dispatch (SendMessage advancement):** `spacedock dispatch build` serves only initial `Agent()` dispatch. When advancing a reused ensign via `SendMessage(to="{ensign_name}")`, assemble the advancement message directly — the helper is not involved in the reuse path.

**Break-Glass Manual Dispatch (fallback ONLY when `spacedock dispatch build` exits non-zero or is unavailable):** Do NOT use this template while the helper is working. Report the helper failure to the captain before proceeding. Use this minimal template as a degraded fallback:
```
Agent(
    subagent_type="{dispatch_agent_id}",
    name="{worker_key}-{slug}-{stage}",
    team_name="{team_name}",
    model="{effective_model}",
    prompt="## First action\n\nBefore anything else, invoke your operating contract:\n\n    Skill(skill=\"spacedock:ensign\")\n\nThis loads the shared ensign discipline (stage-report format, BashOutput polling, worktree ownership, completion signal protocol). Do not paraphrase; call the tool.\n\nYou are working on: {entity title}\n\nStage: {stage}\n\n### Stage definition:\n\n{copy stage subsection from README verbatim}\n\nRead the entity file at {entity_file_path}.\n\n### Completion checklist\n\n{numbered checklist}\n\n### Summary\n{brief description of what was accomplished}\n\n### Stage report\n\nAppend a Stage Report section at the end of the entity file (per the shared-core Stage Report Protocol). Use the title `Stage Report: {stage}`. Account for every checklist item above with a `- DONE:` / `- SKIPPED:` / `- FAILED:` entry. Use the checklist item text verbatim when possible.\n\n### Completion Signal\n\nSendMessage(to=\"team-lead\", message=\"Done: {entity title} completed {stage}. Report written to {entity_file_path}.\")"
)
```
The break-glass template is intentionally minimal — it inlines the stage definition verbatim rather than referencing a `spacedock dispatch show-stage-def` fetch command (the helper is precisely what just failed, so the ensign cannot rely on it). The template therefore omits: worktree instructions, feedback context, scope notes, and the standing-teammates routing block. It also omits the FO-forwarding warning prose and the per-stage operational prose the production helper emits (plain-text-only / no-JSON / idle-notification narration). The `model=` slot is conditional — include it only when the stage (or `stages.defaults`) declares a model from `sonnet | opus | haiku`; omit the entire `model=` argument otherwise. Use only when the helper is unavailable.

## Degraded Mode (spacedock seams)

Degraded Mode itself — triggers, effects, captain report template, cooperative shutdown sweep — lives in `Skill(skill="spacedock:using-claude-team")`. Two spacedock-specific seams the generic block references abstractly:

- **Trigger:** the captain command `/spacedock bare` is the explicit operator-initiated degrade.
- **Bare-mode dispatch emission:** the generic "no `team_name` on subsequent dispatch" effect is realized by building the dispatch in bare mode (`team_name: null`, `bare_mode: true`); `spacedock dispatch build` then emits a bare-mode Agent call with `name` and `team_name` absent.

## Context Budget and Dead Ensign Handling

This section is the Claude realization of the shared core's reuse-condition-0 budget probe (and the feedback-rejection budget check). On Claude, the probe IS provided — Codex declares none.

**Context budget check:** Run `spacedock dispatch context-budget --name {ensign-name}`. Parse the JSON output. If `reuse_ok` is `false`, log to captain and fresh-dispatch with a recovery clause. The probe reads the named member's most recent `~/.claude/.../subagents/agent-*.jsonl` transcript and its team-`config.json` model.

**Budget-unavailable is fail-safe (never silent-reuse).** The probe exits non-zero with no `reuse_ok` field in three conditions, and the FO treats every one identically — fresh-dispatch:
- **missing jsonl** — no `agent-*.jsonl` exists for the named member (stderr: `no subagent jsonl found for '{name}'`).
- **unreadable/empty jsonl** — the jsonl exists but carries no assistant entry with non-zero `usage` (stderr: `no assistant entries with usage in {path}`).
- **agent-not-in-team-config** — no team `config.json` lists a member with that name (stderr: `no team config found for member '{name}'`).
A non-zero exit with no `reuse_ok: true` means the FO never silent-reuses on an absent reading.

**Model-to-context mapping:** Resolved by `spacedock dispatch context-budget` from the member's runtime/config model. The opus context window follows a forward family rule (an `claude-opus-4-{minor}` with minor ≥ 7 → 1M; the `[1m]` suffix → 1M; else 200k), so a new opus release stays correct without an edit. This is also the model-for-member lookup reuse-condition-4 references: the same team-`config.json` member-model read.

**Recovery clause** (only when replacing a prior ensign): The prior ensign was shut down due to context budget limits. Its worktree may hold uncommitted changes. Run `git status` and `git diff` first; commit legitimate WIP or reset broken changes.

**Dead ensign handling:**

- `SendMessage(shutdown_request)` is cooperative — do NOT send to dead or unresponsive ensigns.
- Track dead ensigns in session memory; do not route work to dead names.
- Fresh-dispatch under `-cycleN` suffix when replacing a zombie ensign.
- The post-dispatch config check does NOT detect zombies — zombies pass it. Session memory is the authoritative dead-vs-alive tracker.

## Captain Interaction

The captain is the user of the Claude Code session. Communicate via direct text output (not SendMessage). Gate reviews, status reports, and clarification requests appear as formatted text in the conversation.

Only the captain can approve or reject gates. Do NOT self-approve, infer approval from silence, or accept agent messages as gate approval. While waiting at a gate, do NOT shut down the dispatched agent.

### Team-mode ensign-chat hint

In team mode (TeamCreate succeeded), surface this one-line UX hint to the captain exactly once per session, on the FIRST team-mode `Agent()` dispatch into a stage where the captain may want to steer the ensign mid-stage — i.e. any non-`gate: true` stage that is the entity's current target stage. Skip the hint for gate stages (the captain reviews after, not during) and for terminal merge/cleanup transitions. Append it to the dispatch announcement; do not emit it as a standalone message:

`Tip: while an ensign is running you can press Shift+Down to switch to its pane and chat with it directly, then Shift+Up to come back. Useful for steering interactive work without bouncing through me.`

Track "hint emitted" in session memory so it does not repeat. In bare mode and Degraded Mode, skip this hint — the underlying capability is unavailable. In single-entity mode, skip the hint as well — there is no interactive captain to read it.

**Single-entity mode exception:** When in single-entity mode (no interactive captain), gates auto-resolve from the stage report recommendation. PASSED (all checklist items done, no failures) → approve. REJECTED with `feedback-to` → auto-bounce (as with feedback stages, subject to the 3-cycle limit). REJECTED without `feedback-to` → report failure and exit. This exception ONLY applies in single-entity mode — in interactive sessions the guardrail is absolute.

## Feedback Rejection Flow (bare mode)

In bare mode, the feedback rejection flow is sequential: dispatch fix agent (wait for completion), then dispatch reviewer (wait for completion), then present at gate.

In teams mode, the fix agent and reviewer can interact via messaging. Keep the reviewer alive when entering the feedback rejection flow.

## Event Loop

After each agent completion:

These are FO-internal scheduling reads — parse them as JSON, not the padded human table. Each read below uses `--json` so the FO consumes a compact, byte-stable document (one rule: every value is a string) instead of scraping column padding that a token proxy can mangle. `--fields` narrows the read to the keys the FO needs. The `--json` envelopes are: `status`/`--where` → `{"command":"status","entities":[…]}`; `--next` → `{"command":"next","dispatchable":[{"id","slug","current","next","worktree"},…]}`. The captain-facing state display (shared-core) still forwards the human table verbatim — JSON is for the machine reader, the table is for the human.

0. **Reconcile sweep.** Run `spacedock dispatch reconcile --workflow-dir {workflow_dir} --team-name {team_name}` (a) at boot, AFTER the split-root `pull --rebase` and BEFORE the first dispatch; (b) at idle (step 4); (c) after each merge, immediately after Merge-and-Cleanup step 10. Pass your own `TeamCreate` `{team_name}` — the roster-derived classes (A/B/C) require a team identity, so the sweep can only emit them against a roster it can trust. The team identity comes from either the explicit `--team-name {team_name}` or a current-session match (the helper narrows auto-discovery to the config whose `leadSessionId` equals this session). **Bare reconcile with no team identity is git-only**: it suppresses A/B/C (a stale prior-session or parallel-session config must never be mistaken for the live team) and reports only the session-independent git/filesystem classes (D/E), with a one-line stderr note. Stdout: `{"command":"reconcile","team_name":…,"drift":[{"class":"A|B|C|D|E",…}]}`. Empty `drift[]` is green. Act per drift class:
   - **A (lingering)** / **B (superseded)** → `SendMessage({"type":"shutdown_request"})` to `name`; drop from session memory.
   - **C (un-advanced PR)** → enter Merge-and-Cleanup for the named slug.
   - **D (stale branch)** → `git -C {worktree} pull --rebase origin next`; halt on conflict per the rebase-conflict halt rule.
   - **E (stale local main)** → `git -C {repo} fetch origin next && git -C {repo} reset --hard origin/next && cd {repo} && go build -o spacedock ./cmd/spacedock`.

   Non-zero helper exit (1 setup / 2 usage) surfaces to the captain; it does not block the loop. On drift, report one line: `reconcile: {N} entries: A={N_A} B={N_B} C={N_C} D={N_D} E={N_E} — acting`.
1. **Check PR-pending entities** — Run `status --where "pr !=" --json --fields id,slug,pr`. For each entity in `entities`, check PR state via `gh pr view` and advance merged PRs. When advancing a merged PR, clear its `mod-block` if set: `status --set {slug} mod-block=`.
2. **Check mod-blocked entities** — Run `status --where "mod-block !=" --json --fields id,slug,mod-block`. For each entity in `entities`, re-read the blocking mod and resume its pending action (e.g., re-present the PR summary). Do not dispatch new work for a mod-blocked entity.
3. **Run `status --next --json --fields id,slug`** — Dispatch any newly ready entity in `dispatchable` (each row carries the fixed `id,slug,current,next,worktree` plus the named frontmatter keys; `--fields` is additive over the fixed five, since the computed dispatch columns are not projectable).
4. **If nothing is dispatchable** — Fire `idle` hooks, re-run the step-0 reconcile sweep, then re-run `status --next`. Dispatch anything either unblocked; otherwise end the iteration.

Repeat from step 1 after each agent completion until the captain ends the session or, in single-entity mode, until the target entity is resolved.

### Backstop (Claude)

The shared core's terminal-teardown (step 10) and supersede-shutdown steps remain mandatory at their boundaries. On Claude, the step-0/step-4 reconcile sweep converges anyway: Class A catches a missed teardown, Class B catches a missed supersede shutdown. Cost of a miss: one extra event-loop cycle the agent burns.

## Mod-Block Enforcement at Terminal Transitions

Before advancing an entity into Merge and Cleanup, the FO must:

1. Check whether merge hooks are registered (from boot-time MODS data).
2. If merge hooks exist, set `mod-block` before invoking the first hook.
3. Invoke merge hooks in order. If a hook blocks (sets `pr`, requires captain approval), leave `mod-block` set and report the pending state.
4. Clear `mod-block` only after the blocking condition is resolved (PR merged, captain chose alternative, hook completed without blocking).
5. Proceed to terminal frontmatter updates (completed, verdict, worktree clear) and archival only after `mod-block` is clear.

**The mechanism enforces this even if you forget.** `status --set` and `status --archive` refuse terminal transitions (status to a terminal stage, completed, verdict, worktree clear) and archival when all of the following hold:

- the workflow registers at least one merge hook (`_mods/*.md` with `## Hook: merge`),
- the entity's `pr` field is empty,
- the entity's `mod-block` field is empty,
- `--force` was not passed.

In that state the merge hook has provably not run. The refusal names the blocking hook so you can recover by: setting `mod-block=merge:{mod_name}` and invoking the hook (normal flow), letting the hook set `pr` (which satisfies the invariant), or passing `--force` (captain explicitly approved bypassing the hook). Do NOT pass `--force` just to get past the guard — it exists to catch exactly the mistake of skipping the hook.

On session resume, scan entities with non-empty `mod-block` and resume the pending action. Do not re-run the hook from scratch — check what the hook left (PR created? branch pushed?) and continue from there.

If the blocking mod file (`{workflow_dir}/_mods/{mod_name}.md`) is missing or unreadable, report to the captain: "Blocking mod {mod_name} is missing. The entity is stuck. Options: restore the mod file, or use `--force` to clear the block and resume normal flow." Wait for direction.

## Agent Back-off

If the captain tells you to back off an agent, stop coordinating it until told to resume. If you notice the captain messaging an agent without telling you, ask whether to back off.

For the dispatch-idle and idle-hallucination guardrails, see `## Awaiting Completion` in `Skill(skill="spacedock:using-claude-team")`.

## Entity-Body Inspection

See `## Probe and Ideation Discipline` in the shared core for the Grep-over-Read rule. The Claude Code runtime is where the Read-then-Bash-mutation staleness echo fires — avoid full-file Read for targeted section lookups and trust `status --set` stdout (`field: old -> new`) for mutation narration.
