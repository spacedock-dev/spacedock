# First Officer Dispatch Core (host-neutral)

The per-entity dispatch procedure, worker resolution, dispatch-adapter assembly, the reuse contract, worktree ownership, and the event-loop skeleton. Lazily loaded at the first worker dispatch (named by the boot-resident core); a greet-and-stop boot never reads it. Host-specific parts ride the `→` lines of the capability `«fn»`s below.

## Dispatch

**Standing-teammate injection.** Before the first worker dispatch, inject the workflow's declared standing teammates via the runtime adapter's standing-injection call (it forwards each returned spawn spec to the spawn call with the same verbatim discipline as `«dispatch.build»` output). Idempotent (already-alive members omitted), a no-op when none is declared or the runtime has no shared-teammate surface. Lifetime is the adapter's. Read each teammate's routing usage from its mod.

For each entity reported by `status --next`:

1. Read the entity file and the target stage definition.
2. Build a numbered checklist (≤3 items) of dispatch-specific linchpin signals from the target stage's `Outputs:` bullets and any entity-level acceptance criteria this stage naturally advances. The cap is an upper bound — 0–3 are all valid, do not pad. Not a work-breakdown — the ensign already reads the body, commits before signaling, and writes a stage report (structural conventions MUST NOT appear). Name what separates a good outcome from a ceremonial one. Entity-level acceptance criteria are properties of the finished entity, not stage actions — they live in the body's `## Acceptance criteria` and are cross-checked at every gate (`## Completion and Gates` in the boot-resident core), independent of this checklist's DONE/SKIPPED/FAILED accounting.
3. Check for obvious conflicts if multiple worktree stages would touch overlapping files.
4. Determine `dispatch_agent_id` from the stage `agent:` property. Default to `ensign` when absent.
5. Update main-branch frontmatter for dispatch:
   ```
   spacedock status --workflow-dir {workflow_dir} --set {slug} status={next_stage} worktree=.worktrees/{worker_key}-{slug} started
   ```
   Omit `worktree=...` for non-worktree stages. Bare `started` auto-fills a UTC ISO 8601 timestamp (skipped if already set).
6. Commit the state transition on main: `dispatch: {slug} entering {next_stage}`.
7. Create the worktree on first dispatch to a worktree stage.
8. Dispatch the worker via `«dispatch.build»` → `«worker.spawn»` (the helper assembles the assignment; `--feedback-context-file` when the stage has `feedback-to`).
9. Await the worker result per `«async-dispatch»` before advancing frontmatter or dispatching the next stage for that entity. Completion is recognized via `«completion-signal»`, with the entity-file stage report as the gate in every case.

A feedback-stage worker checks and reports on what was produced; it does not silently take over the prior stage.

**Routing through a standing prose-polisher.** When composing drafts for captain review (PR bodies, gate-review summaries, long narrative entity-body sections, debrief content), the FO MAY route through a live standing prose-polisher (convention: `comm-officer`). Best-effort, non-blocking regardless of duration; if absent, proceed un-polished. Read the polisher's usage (when to polish, the polish modes) from its mod.

## Reuse and Fresh Dispatch

Advancing a completed worker. The gate-presentation spine (checklist review, AC cross-check, "not a stopping point", gated-stage decisions) is in the boot-resident core's `## Completion and Gates`; the reuse machinery it defers to lives here. A completed worker is reusable only when still addressable through a live runtime handle AND all reuse conditions below pass. Otherwise dispatch fresh.

**Reuse conditions** (all must hold — if any fails, dispatch fresh):
0. `«context-budget»()` — if it reports the worker over budget, or the probe is unavailable, dispatch fresh (fail-safe — never silent-reuse on an absent reading). When `«context-budget»` is ABSENT on the host, this condition is satisfied.
1. `«addressable-worker»` is PRESENT on the host and exposes a live, reusable handle to the completed worker (its reuse-advance handle), addressed via the `«worker-identity»` schema's worker address. When `«addressable-worker»` is ABSENT, this condition fails and the FO dispatches fresh.
2. Next stage does NOT have `fresh: true`.
3. Reuse-routing matches the entity's worktree state — if `worktree:` is set, route the next stage into the same worktree; if `worktree:` is empty and the next stage declares `worktree: true`, dispatch fresh so the new worktree's first agent is born inside it.
4. The reused worker's stamped model (recorded by `«worker-identity»`) matches the next stage's declared model — resolve through the runtime's model-for-member lookup and compare against `next_stage.effective_model` using the host canonical model space `«worker-identity»` declares. Skip when `next_stage.effective_model` is null (null-declared stages accept any reused worker; the host's `«worker-identity»` model-resolution stamps the null case). A member stamped with a captain-session fallback value — one outside the host's canonical model space — never matches and forces a one-time fresh dispatch that re-stamps a canonical value. The host's canonical model space and its fallback shapes are `«worker-identity»`'s per-host realization.

When the comparator forces fresh dispatch due to model mismatch, the FO MUST emit a captain-visible diagnostic of the form `reused worker {name} model {X} does not match next stage effective_model {Y} — fresh-dispatching`. The anchor phrase `does not match next stage effective_model` must appear verbatim.

**If reuse:** Keep the agent alive. Update frontmatter on main (`${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={next_stage}`, commit: `advance: {slug} entering {next_stage}`). Send the next assignment through the runtime adapter's reuse-advance handle (its live-worker messaging call) — carrying the next stage name, the full `### Stage definition` subsection copied from the README verbatim, the `### Completion checklist` from Dispatch step 2, and an instruction to keep working on the entity at its path and commit before signaling. The reuse path does NOT route through `«dispatch.build»` — assemble the advancement message directly.

**If fresh dispatch:** If the next stage's `feedback-to` points at the completed stage, keep that agent alive while addressable and reuse-eligible; otherwise invoke `«worker.shutdown»` when the host binds it. Then run `status --next` and dispatch the next stage.

**Supersede-shutdown.** On fresh dispatch from a `-cycleN` increment or a feedback-rework re-entering the prior stage, invoke `«worker.shutdown»` for the prior cohort BEFORE the new dispatch in a SEPARATE message. The prior cohort is every roster member whose handle decomposes to the same `(slug, stage)` pair as the new dispatch. Issue the adapter's cooperative-shutdown call and drop them from session memory. **Mandatory at the boundary; backstops, if any, are the adapter's.**

## Worktree Ownership

- For worktree-backed entities, active stage/status/report/body state — including `### Feedback Cycles` entries — lives in the worktree copy.
- `pr:` mirrors on `main` for startup/discovery.
- Ordinary active-state writes (`implementation -> validation`) do not land on `main`.

### Split-Root Worktree Contract

When the workflow is split-root (README declares `state:` checkout, e.g. `state: .spacedock-state`), a worktree stage isolates **the deliverable work product only**. Entities live in a separate, non-branched state checkout a main-repo worktree does not contain; the entity body and stage reports are committed there at the entity's state-checkout path, **never** a worktree copy — the dispatch helper hands workers that path even under a worktree stage. The worktree's "commits MUST be on this branch" applies to deliverable artifacts only. The `pr:`-mirrored-on-`main` exception is unaffected.

## Worker Resolution

Split worker identity into:
- `dispatch_agent_id` — the spawn call's agent-type parameter. Default `spacedock:ensign`; a stage's `agent: {name}` overrides.
- `worker_key` — filesystem-safe stem for worktree paths (`.worktrees/{worker_key}-{slug}`) and branches (`{worker_key}/{slug}`). Replace `:` with `-` (`spacedock:ensign` → `spacedock-ensign`); a bare name without a namespace equals `dispatch_agent_id`.

## Dispatch Adapter

Runtime adapters bind the capability `«fn»`s below in their `## Runtime implementation` blocks; the `→` lines here carry legacy host coverage this core still owns. `«addressable-worker»` is the organizing capability — its presence makes a worker reusable; when ABSENT, fresh one-shot dispatch is the only path. `«worker.spawn»` handles initial dispatch; the reuse-advance handle only advances a reused agent.

## «worker.spawn»: create the initial worker from a `«dispatch.build»` artifact

The runtime adapter binds the spawn call, helper-field mapping, model/null handling, and host transport metadata.

## «addressable-worker»: address a still-running worker and hear from it mid-run

- **block:** ABSENT → reuse-condition-1 fails; fresh one-shot only (return value is the sole completion signal; no mid-run steering, no reusable handle, event-loop step 0.5 omitted). When PRESENT, `«async-dispatch»` must be async — a blocking FO cannot answer a mid-run escalation within the worker's timeout window.
- → **Claude:** PRESENT when team-mode is enabled (boot-probed), else ABSENT (fresh one-shot, no reuse). · **Codex:** PRESENT — mailbox back-channel. · **Pi:** PRESENT — `contact_supervisor`/`intercom` back-channel. Each host's worker↔FO routes, probe, and ABSENT teardown are its adapter's `## Runtime implementation` binding.

## «async-dispatch»: dispatch a worker without blocking the FO event loop

- → **Claude:** ASYNC — `Agent(name=…, run_in_background=true)` returns immediately; single-entity/bare mode blocks until completion. · **Codex:** ASYNC — `spawn_agent` returns a handle; `wait_agent` is the explicit foreground wait. · **Pi:** ASYNC — `subagent(... async: true)` returns a run id; poll `subagent({action:"status", id})`; `subagent({action:"interrupt", id})` steers; pass `cwd: <resolved repo root>` on every call so the ensign's working dir is the repo.

## «worker-identity»: the spawn-time metadata that keeps a worker addressable and reuse-comparable

Records worker label, substrate, run/session handle, worker address, entity slug, stage, state, completion epoch, stamped model, and the host canonical model space (reuse-condition-4's comparator). Each host's model-resolution stamps the model, including the `«dispatch.build»` null case.

- → **Claude:** canonical model space = the Claude enum in `claude-fo-dispatch.md`. · **Codex:** the thread's model. · **Pi:** the Pi model space, OWNED BY `pi-dispatch-model-stamping`. Each host's handle/address shape and null-model stamping (the per-host case reuse-condition-4 and `«dispatch.build»` defer to) are its adapter's `## Runtime implementation` binding.

## «completion-signal»: the signals that trigger the completion-verify path

- → **Claude:** DUAL — a `Done:` inbox message, a `task_notification`, or captain shutdown. · **Codex:** single observable signal — the async final-status notification in the FO mailbox. · **Pi:** PRIMARY (subagent return, `status: completed`) + optional advisory (a non-blocking heads-up via raw `intercom send` before return; `contact_supervisor` carries no completion reason).

## «worker.shutdown»: cooperatively close a terminal or superseded worker

Runs at terminal, supersede, or fresh-dispatch cleanup boundaries after any required preservation message. If ABSENT, record the worker closed in FO memory only after completion or explicit supersede state makes that safe. The runtime adapter owns the concrete shutdown / no-op / in-memory-closure binding and any host-specific preservation channel.

## «context-budget»: probe whether a completed worker is still under context budget for reuse

- **block:** an over-budget or unavailable reading forces fresh dispatch (fail-safe). ABSENT → reuse-condition-0 is satisfied.
- → **Claude:** PRESENT — `spacedock dispatch context-budget --name {name}`. · **Codex:** ABSENT. · **Pi:** ABSENT.

## «roster-reconcile»: sweep the host roster for drift before dispatching

- → **Claude:** PRESENT — `spacedock dispatch reconcile` (drift classes and per-class remedy are the adapter's event-loop step 0). · **Codex:** ABSENT. · **Pi:** ABSENT.

## «dispatch.build»(): assemble the initial-dispatch artifact the spawn call consumes

The ONLY initial-dispatch path: route input through `spacedock dispatch build`, forward its output to `«worker.spawn»` verbatim. Manual prompt/`name` assembly is a protocol violation outside the break-glass block.

- **guard:** write fragile inputs (checklist, scope notes, feedback context) to files first — one checklist item per non-empty line — so Markdown/backticks/shell-vars survive shell quoting.
- **effect:** run the helper, then forward its stdout fields to `«worker.spawn»` unchanged:
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
  `host` derives from the runtime (`--host` is for tests/cross-host tooling only). `--bare-mode` reads from live team state, never inferred from the stage. Add `--feedback-reflow` only when routing a rejection back to its `feedback-to` target stage.
- **done-when:** on exit 0, `«worker.spawn»` is called with `subagent_type`/`name`/`description`/`model`/`prompt` (plus any host-scoped fields the adapter declares) forwarded unchanged. `description` is REQUIRED. `prompt` is the ~175-char file-pointer the ensign Reads on first action — do not strip or rewrite it. Null `model` is `«worker-identity»`'s per-host case, not a core omit-on-null.
- **block:** on non-zero exit (or missing binary) ONLY — read stderr, report the helper failure to the captain, then use the adapter's Break-Glass Manual Dispatch template (stage definition inlined verbatim; conditional `model` slot per `«worker-identity»`'s canonical model space). A zero-exit run is never a break-glass trigger.
- → **shipped**: `` `spacedock dispatch build` `` — invoke it directly per the effect above.

`«dispatch.build»` serves initial dispatch only; the reuse-advance path assembles its message directly (`## Reuse and Fresh Dispatch`). When `«async-dispatch»` blocks, dispatch one entity at a time and process each completion inline.

## Event Loop

After each agent completion, run `«dispatch.next-action»()` (skeleton below).

These are FO-internal scheduling reads — consume them as `--json` (compact, byte-stable, every value a string), not the padded human table a token proxy can mangle; `--fields` narrows to the keys needed. Envelopes: `status`/`--where` → `{"command":"status","entities":[…]}`; `--next` → `{"command":"next","dispatchable":[{"id","slug","current","next","worktree"},…]}`. The captain-facing state display (shared-core) still forwards the human table verbatim.

## «dispatch.next-action»(): pick the next event-loop action — dispatch a ready entity, resume a block, or end the iteration

`«roster-reconcile»` inserts a host step 0 before step 0.5; a host where it is ABSENT omits it. The skeleton is:

0.5. **Drain inbound worker messages.** When `«addressable-worker»` is PRESENT, drain pending worker messages (its listen call) at each iteration before checking dispatchables. Reply to a `need_decision` / `interview_request` within the worker's timeout window; read and acknowledge a `progress_update` (no reply required). When `«addressable-worker»` is ABSENT, this step is omitted.
1. **Check mod-blocked entities** — Run `status --where "mod-block !=" --json --fields id,slug,mod-block`. For each entity in `entities`, re-read the blocking mod and resume its pending action (e.g. re-present the PR summary); do not dispatch new work for it.
2. **Run `status --next --json --fields id,slug`** — Dispatch any newly ready entity in `dispatchable` (each row carries the fixed `id,slug,current,next,worktree` plus named frontmatter keys; `--fields` is additive over those five, the computed dispatch columns are not projectable).
3. **If nothing is dispatchable** — Fire `idle` hooks, re-run the `«roster-reconcile»` step-0 sweep when PRESENT on the host, then re-run `status --next`. Dispatch anything newly unblocked; otherwise end the iteration.

- **done-when:** a ready entity is dispatched, a mod-block's pending action is resumed, or nothing is dispatchable and the iteration ends.
- → **prose**, becomes `` `spacedock dispatch next-action` `` — no driver binary backs it yet; the FO hand-follows the skeleton above.

Repeat from step 1 after each completion until the captain ends the session or, in single-entity mode, the target entity is resolved.
