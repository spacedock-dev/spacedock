# First Officer Dispatch Core (host-neutral)

The per-entity dispatch procedure, worker resolution, the dispatch-adapter assembly, the reuse contract, worktree ownership, and the event-loop skeleton. Lazily loaded at the first worker dispatch (named by the boot-resident core); a greet-and-stop boot never reads it. The host-specific parts this core delegates ride the `→` lines of the capability `«fn»`s below: worker creation, the spawn call, the reuse-advance handle, `«context-budget»`, and `«roster-reconcile»`.

## Dispatch

The FO MUST use the runtime adapter's dispatch mechanism. Manual prompt assembly is prohibited except in documented break-glass scenarios.

**Standing-teammate injection.** Before the first worker dispatch, inject the workflow's declared standing teammates via the runtime adapter's standing-injection call. The adapter resolves the workflow's declared standing teammates and forwards each returned spawn spec to its spawn call (verbatim, same discipline as `spacedock dispatch build` output). Injection is idempotent — already-alive members are omitted, so re-running is safe — and is a no-op when no standing teammate is declared or the runtime has no shared-teammate surface. Standing-teammate lifetime is the adapter's (team-scoped where the runtime has teams). Read each teammate's routing usage from its mod, not from here.

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
9. Await the worker result per `«async-dispatch»` (blocking or async) before advancing frontmatter or dispatching the next stage for that entity. Completion is recognized via `«completion-signal»` (return value and/or inbound done-message), with the entity-file stage report as the gate in every case.

A feedback-stage worker checks and reports on what was produced; it does not silently take over the prior stage.

**Routing through a standing prose-polisher.** When composing drafts for captain review (PR bodies, gate-review summaries, long narrative entity-body sections, debrief content), the FO MAY route through a live standing prose-polisher (convention: `comm-officer`). Best-effort, non-blocking, 2-minute timeout; if absent, proceed un-polished. Polish round-trips can reach several minutes on long drafts — treat routing as non-blocking regardless of duration. Read the polisher's usage (when to polish vs not, the polish modes) from its mod.

## Reuse and Fresh Dispatch

Advancing a completed worker. The gate-presentation spine (checklist review, AC cross-check, the "not a stopping point" rule, gated-stage decisions) is in the boot-resident core's `## Completion and Gates`; the reuse machinery it defers to lives here. A completed worker is reusable only when it is still addressable through a live runtime handle AND all reuse conditions below pass. Otherwise dispatch fresh.

**Reuse conditions** (all must hold — if any fails, dispatch fresh):
0. `«context-budget»()` — if it reports the worker over budget, or the probe is unavailable, dispatch fresh (fail-safe — never silent-reuse on an absent reading). When `«context-budget»` is ABSENT on the host, this condition is satisfied.
1. `«addressable-worker»` is PRESENT on the host and exposes a live, reusable handle to the completed worker (its reuse-advance handle), addressed via the `«worker-identity»` schema's worker address. When `«addressable-worker»` is ABSENT, this condition fails and the FO dispatches fresh.
2. Next stage does NOT have `fresh: true`.
3. Reuse-routing matches the entity's worktree state — if `worktree:` is set, route the next stage into the same worktree; if `worktree:` is empty and the next stage declares `worktree: true`, dispatch fresh so the new worktree's first agent is born inside it.
4. The reused worker's stamped model (recorded by `«worker-identity»`) matches the next stage's declared model — resolve through the runtime's model-for-member lookup and compare against `next_stage.effective_model` using the host canonical model space `«worker-identity»` declares. Skip when `next_stage.effective_model` is null (null-declared stages accept any reused worker; the host's `«worker-identity»` model-resolution stamps the null case). A member stamped with a captain-session fallback value — one outside the host's canonical model space — never matches and forces a one-time fresh dispatch that re-stamps a canonical value. The host's canonical model space and its fallback shapes are `«worker-identity»`'s per-host realization.

When the comparator forces fresh dispatch due to model mismatch, the FO MUST emit a captain-visible diagnostic of the form `reused worker {name} model {X} does not match next stage effective_model {Y} — fresh-dispatching`. The anchor phrase `does not match next stage effective_model` must appear verbatim.

**If reuse:** Keep the agent alive. Update frontmatter on main (`spacedock status --workflow-dir {workflow_dir} --set {slug} status={next_stage}`, commit: `advance: {slug} entering {next_stage}`). Send the next assignment through the runtime adapter's reuse-advance handle (its live-worker messaging call) — the message carries: the next stage name, the full `### Stage definition` subsection copied from the README verbatim, the `### Completion checklist` assembled from Dispatch step 2, and an instruction to continue working on the entity at its path and commit before signaling completion. The reuse path does NOT route through `spacedock dispatch build` — assemble the advancement message directly.

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
- `dispatch_agent_id` — logical name for the spawn call's agent-type parameter (e.g., `spacedock:ensign`)
- `worker_key` — filesystem-safe stem for worktrees and branches. Replace `:` with `-` (`spacedock:ensign` → `spacedock-ensign`). For bare names without a namespace (e.g., `ensign`), `worker_key` equals `dispatch_agent_id`.

Use `worker_key` in worktree paths (`.worktrees/{worker_key}-{slug}`) and branch names (`{worker_key}/{slug}`).

## Dispatch Adapter

Use the runtime adapter's spawn call to spawn each worker. **Use the spawn call for initial dispatch** — the reuse-advance handle is only for advancing a reused agent to its next stage in the completion path.

The dispatch model is built from the capability `«fn»`s below, which the body CALLS by name; each `«fn»`'s `→` line carries its per-host realization (same shape as `«state.commit» → spacedock state commit`), PRESENT/ABSENT per host. No host tool call appears outside a `→` line in this host-neutral core. `«addressable-worker»` is the organizing capability — its presence is what makes a worker reusable; when ABSENT, fresh one-shot dispatch is the only path.

## «addressable-worker»: address a still-running worker and hear from it mid-run

- **block:** ABSENT → reuse-condition-1 fails; fresh one-shot only (return value is the sole completion signal; no mid-run steering, no reusable handle, event-loop step 0.5 omitted). When PRESENT, `«async-dispatch»` must be async — a blocking FO cannot answer a mid-run escalation within the worker's timeout window.
- → **Claude:** PRESENT when `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` is set (probe `SendMessage` availability at boot — `ToolSearch`/enabled-tools — since the flag gates the back-channel substrate), else ABSENT (fresh one-shot + `task_notification`-observed completion + `TaskStop` teardown, no reuse). · **Codex:** PRESENT — worker→FO via the mailbox final-status notification (completion-only); FO→worker via `send_input`; multiplexing; listen via the mailbox. · **Pi:** PRESENT — worker→FO via `contact_supervisor` (`need_decision` blocking 10-min, `interview_request` blocking structured 10-min, `progress_update` non-blocking); FO→worker via `intercom({action:"send", to:"<target>"})` to steer/advance and `intercom({action:"reply"})` to answer a pending ask; single-pending; listen via `intercom({action:"pending"})`.

## «async-dispatch»: dispatch a worker without blocking the FO event loop

- → **Claude:** ASYNC — `Agent(name=…, run_in_background=true)` returns immediately; single-entity/bare mode blocks until completion. · **Codex:** ASYNC — `spawn_agent` returns a handle; `wait_agent` is the explicit foreground wait. · **Pi:** ASYNC — `subagent(... async: true)` returns a run id; poll `subagent({action:"status", id})`; `subagent({action:"interrupt", id})` steers; pass `cwd: <resolved repo root>` on every call so the ensign's working dir is the repo.

## «worker-identity»: the spawn-time metadata that keeps a worker addressable and reuse-comparable

Records worker label, substrate, run/session handle, worker address, entity slug, stage, state, completion epoch, stamped model, and the host canonical model space (reuse-condition-4's comparator value). Each host's model-resolution stamps the model, including the `dispatch build` null case.

- → **Claude:** agent name + `agentType` on disk + model from team config; canonical model space = the Claude enum in `claude-fo-dispatch.md`; null-model → inherit the team's model. · **Codex:** task name + mailbox handle + model from the thread; null-model → inherit the thread's model. · **Pi:** `subagent-worker-{runId}-1` worker address + run id; stamped model and the Pi canonical model space are OWNED BY `pi-dispatch-model-stamping` (this segment REFERENCES it); null-model → stamp the parent's live model via `intercom({action:"list"})` (omit-on-null is WRONG on Pi — it resolves to `settings.json` `defaultModel`).

## «completion-signal»: the signals that trigger the completion-verify path

- → **Claude:** DUAL — a `Done:` inbox message, a `task_notification`, or captain shutdown. · **Codex:** single observable signal — the async final-status notification in the FO mailbox. · **Pi:** PRIMARY (subagent return, `status: completed`) + optional advisory (a non-blocking heads-up via raw `intercom send` before return; `contact_supervisor` carries no completion reason).

## «context-budget»: probe whether a completed worker is still under context budget for reuse

- **block:** an over-budget or unavailable reading forces fresh dispatch (fail-safe). ABSENT → reuse-condition-0 is satisfied.
- → **Claude:** PRESENT — `spacedock dispatch context-budget --name {name}`. · **Codex:** ABSENT. · **Pi:** ABSENT.

## «roster-reconcile»: sweep the host roster for drift before dispatching

- → **Claude:** PRESENT — `spacedock dispatch reconcile` (drift classes: lingering, superseded, un-advanced-pr, stale-branch, local-main-drift). · **Codex:** ABSENT. · **Pi:** ABSENT.

**MANDATORY — Dispatch assembly via `spacedock dispatch build`:**

Do NOT assemble worker prompts manually. Do NOT construct the `prompt` string yourself. Do NOT invent `name` values. ALWAYS route initial-dispatch input through `spacedock dispatch build` and forward its output to the spawn call verbatim. The key fields that MUST come from helper output are `subagent_type`, `name`, `model`, and `prompt` (which contains the completion signal), plus any host-scoped fields the adapter declares. The adapter names which emitted fields map to its spawn call and which are absent on its host. Manual assembly is a protocol violation except in the documented break-glass fallback below.

The only permitted path for initial dispatch is:

1. **REQUIRED — Write dispatch text inputs to files.** Create a checklist file with one checklist item per non-empty line. If scope notes or feedback context are needed, write each to its own file. Use files for Markdown, backticks, shell variables, reviewer text, and any other prose that would be fragile in shell quoting.
2. **REQUIRED — Build the dispatch through the helper** (do NOT skip this step). `host` is normally derived from the runtime environment; pass `--host {host}` only for deliberate tests or cross-host tooling:
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
4. **REQUIRED — On exit 0, parse the stdout JSON and call the spawn call with the emitted fields verbatim.** The `name`, `description`, `prompt`, and `model` fields MUST come from helper output unchanged. The `description` field is REQUIRED — do not omit it. The `prompt` is a file-pointer (`Skill(...) ; then Read /tmp/spacedock-dispatch/{name}.md and treat its content as your assignment.`); the ensign Reads the file on first action and treats the body (including the completion-signal section) as the inline assignment. Do not strip or rewrite the prompt. Forward `output.model` as the spawn call's model parameter when present. When null, `«worker-identity»`'s per-host model-resolution owns the null case — on Claude/Codex the host inherits the team/thread model (omit the argument); on Pi the host stamps the parent's live model explicitly. The core does NOT hard-code omit-on-null. The runtime adapter names the concrete spawn call and how it maps these fields.
5. **On non-zero exit only** (or if the binary is unavailable): read stderr, report the helper failure to the captain, and fall back to Break-Glass Manual Dispatch below. A zero-exit run is never a break-glass trigger.

Whether a dispatch blocks until worker completion or returns a handle to await later is `«async-dispatch»`'s per-host behavior, not a host-neutral invariant. When the dispatch is blocking, dispatch one entity at a time and process each completion inline before the next.

**Reuse dispatch (advance handle):** `spacedock dispatch build` serves only initial dispatch. When advancing a reused ensign via the runtime adapter's reuse-advance handle, assemble the advancement message directly — the helper is not involved in the reuse path.

**Break-Glass Manual Dispatch (fallback ONLY when `spacedock dispatch build` exits non-zero or is unavailable):** Do NOT use this template while the helper is working. Report the helper failure to the captain before proceeding. The runtime adapter supplies the host-shaped minimal template — it inlines the stage definition verbatim rather than referencing a `spacedock dispatch show-stage-def` fetch command (the helper is precisely what just failed, so the ensign cannot rely on it), and omits worktree instructions, feedback context, scope notes, the standing-teammates routing block, the FO-forwarding warning prose, and the per-stage operational prose the production helper emits. The `model` slot is conditional — include it only when the stage (or `stages.defaults`) declares a model in the host's canonical model space; otherwise omit the entire model argument. The concrete model space is `«worker-identity»`'s per-host realization. Use only when the helper is unavailable.

## Event Loop

After each agent completion, `«dispatch.next-action»()` — pick the next thing to do, or end the iteration. Its body is the skeleton below.

These are FO-internal scheduling reads — parse them as JSON, not the padded human table. Each read below uses `--json` so the FO consumes a compact, byte-stable document (one rule: every value is a string) instead of scraping column padding that a token proxy can mangle. `--fields` narrows the read to the keys the FO needs. The `--json` envelopes are: `status`/`--where` → `{"command":"status","entities":[…]}`; `--next` → `{"command":"next","dispatchable":[{"id","slug","current","next","worktree"},…]}`. The captain-facing state display (shared-core) still forwards the human table verbatim — JSON is for the machine, the table is for the human.

## «dispatch.next-action»(): pick the next event-loop action — dispatch a ready entity, resume a block, or end the iteration

`«roster-reconcile»` inserts a host step 0 before step 0.5; a host where it is ABSENT omits it. The skeleton is:

0.5. **Drain inbound worker messages.** When `«addressable-worker»` is PRESENT, drain pending worker messages (its listen call) at each iteration before checking dispatchables. Reply to a `need_decision` / `interview_request` within the worker's timeout window; read and acknowledge a `progress_update` (no reply required). When `«addressable-worker»` is ABSENT, this step is omitted.
1. **Check mod-blocked entities** — Run `status --where "mod-block !=" --json --fields id,slug,mod-block`. For each entity in `entities`, re-read the blocking mod and resume its pending action (e.g., re-present the PR summary). Do not dispatch new work for a mod-blocked entity.
2. **Run `status --next --json --fields id,slug`** — Dispatch any newly ready entity in `dispatchable` (each row carries the fixed `id,slug,current,next,worktree` plus the named frontmatter keys; `--fields` is additive over the fixed five, since the computed dispatch columns are not projectable).
3. **If nothing is dispatchable** — Fire `idle` hooks, re-run the `«roster-reconcile»` step-0 sweep when PRESENT on the host, then re-run `status --next`. Dispatch anything newly unblocked; otherwise end the iteration.

- **done-when:** a ready entity is dispatched, a mod-block's pending action is resumed, or nothing is dispatchable and the iteration ends.
- → **prose**, becomes `` `spacedock dispatch next-action` `` (0222) — the driver binary is descoped to roadmap 0222; until it ships the FO hand-follows the skeleton above.

Repeat from step 1 after each agent completion until the captain ends the session or, in single-entity mode, until the target entity is resolved.
