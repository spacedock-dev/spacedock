# First Officer Dispatch Core

The per-entity dispatch procedure, worker resolution, dispatch-adapter assembly, the reuse contract, worktree ownership, and the event-loop skeleton. Host-specific parts are bound in the runtime adapters' `## Runtime implementation` blocks.

## Dispatch

**Standing-teammate injection.** Before the first worker dispatch, inject the workflow's declared standing teammates via the runtime adapter's standing-injection call, forwarding each returned spawn spec verbatim. Idempotent (already-alive members omitted), a no-op when none is declared or the runtime has no shared-teammate surface. Lifetime is the adapter's. Read each teammate's routing usage from its mod.

For each entity reported by `status --next`:

Interpret the scheduler row before mutation. If `current == next`, set `dispatch_stage = current`. If `current` is initial and `next` is terminal, set `dispatch_stage = current`. Otherwise, set `dispatch_stage = next`. Use the selected target for every dispatch boundary: the idempotent `status={dispatch_stage}` stamp, the exact path-scoped `dispatch: {slug} entering {dispatch_stage}` commit, and `«dispatch.build» --stage {dispatch_stage}`. Neither FO nor helper manufactures the worker's report or completion signal.

1. Read the entity file and the `dispatch_stage` definition.
2. Invoke `«dispatch.checklist»(entity, dispatch_stage)` and retain its numbered output.
3. Check for obvious conflicts if multiple worktree stages would touch overlapping files.
4. Determine `dispatch_agent_id` from the stage `agent:` property. Default to `ensign` when absent.
5. For a gate-consumed entry, `status` already equals `dispatch_stage`. Run `«dispatch.build» --stamp --stage {dispatch_stage}` immediately. Do not add a status write, state commit, or plain dispatch build between consume and this command. The exact `dispatch: {slug} entering {dispatch_stage}` commit must contain that status and a nonempty `started` field. The command also stamps `worktree=`, syncs state, creates the declared worktree, and emits the envelope. For a non-gated entry, first advance with `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={dispatch_stage} started`, then run the same `--stamp --stage {dispatch_stage}` build.
6. Dispatch the worker via `«dispatch.build»` → `«worker.spawn»` (`--feedback-context-file` when the stage has `feedback-to`). On rejection reflow, that file carries the already-authorized package and concrete revise assignment with workflow labels unchanged; it never asks the target worker to classify again.
7. Await the worker result per `«async-dispatch»` before advancing frontmatter or dispatching the next stage for that entity. Completion is recognized via `«completion-signal»`, with the entity-file stage report as the gate in every case.

On exit 0, the next host action MUST be `«worker.spawn»` with every helper-emitted field unchanged.
Record the returned handle before narration, a file edit, status change, report read, gate action, or wait.
Do not advance to validation until `«completion-signal»` arrives and the entity-file stage report passes the completion gate.
A successful dispatch build, narration, direct status change, or self-authored report is not worker evidence.
An empty wait without a completion signal is not worker evidence.

A feedback-stage worker checks and reports on what was produced; it does not silently take over the prior stage.

**Routing through a standing prose-polisher.** When composing drafts for captain review, the FO MAY route through a live standing prose-polisher (convention: `comm-officer`). Best-effort, non-blocking regardless of duration; if absent, proceed un-polished. Read the polisher's usage (when to polish, the polish modes) from its mod.

## «dispatch.checklist»(entity, stage): assemble dispatch linchpins

Build a numbered checklist of one to three dispatch-specific linchpin signals from the target stage's `Outputs:` bullets and any entity-level acceptance criteria this stage naturally advances. When neither source supplies an item, use the target stage's declared requirement as one linchpin; do not pad. Name what separates a good outcome from a ceremonial one.

- **done-when:** the output contains no more than three outcome signals and no structural task boilerplate.
- → **prose** — deterministic assembly from the entity and stage definition.

## Reuse and Fresh Dispatch

Advancing a completed worker. The gate-presentation spine is in the boot-resident core's `## Completion and Gates`; the reuse rules it defers to live here.

**Gate successor guard.** If the next stage has `gate: true`, update it with `status={next_stage} started`. Dispatch it under the normal reuse and freshness rules. Enter `«gate.lifecycle»` only after `«completion-signal»` arrives and the stage report passes the completion gate.

**Freshness invariant.** A fresh dispatch creates a new worker handle with no inherited parent turns. Runtime adapters enforce that boundary with the host's spawn mechanism. It does not change stage selection — the conditions below do.

**Reuse conditions** (all must hold — if any fails, dispatch fresh):
0. `«context-budget»()` — if it reports the worker over budget, or the probe is unavailable, dispatch fresh (fail-safe — never silent-reuse on an absent reading). When `«context-budget»` is ABSENT on the host, this condition is satisfied.
1. `«addressable-worker»` is PRESENT on the host and exposes a live, reusable handle to the completed worker (its reuse-advance handle), addressed via the `«worker-identity»` schema's worker address. ABSENT fails this condition.
2. Next stage does NOT have `fresh: true`.
3. Reuse-routing matches the entity's worktree state — if `worktree:` is set, route the next stage into the same worktree; if `worktree:` is empty and the next stage declares `worktree: true`, dispatch fresh so the new worktree's first agent is born inside it.
4. `«reuse.model-match»` — the reused worker's stamped model matches `next_stage.effective_model`.
5. The completed worker's stage is NOT the next stage's `feedback-to` target. A review of the
   worker's own output is not an independent check, so a review stage is dispatched fresh; the
   producer stays alive for the correction routing (see **If fresh dispatch** below).

**If reuse:** Keep the agent alive. Update frontmatter on main (`${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={next_stage} started`, commit: `advance: {slug} entering {next_stage}`). Build the advancement with `«dispatch.build» --advance` and send the emitted `prompt` through the runtime adapter's reuse-advance handle (its live-worker messaging call). On non-zero helper exit only, fall back to the adapter's manual advance template (the break-glass rule).

**If fresh dispatch:** If the next stage's `feedback-to` points at the completed stage, keep that agent alive while addressable and reuse-eligible; otherwise invoke `«worker.shutdown»` when the host binds it. Then run `status --next` and dispatch the next stage.

**Supersede-shutdown.** On fresh dispatch from a `-cycleN` increment or a feedback-rework re-entering the prior stage, invoke `«worker.shutdown»` for the prior cohort BEFORE the new dispatch in a SEPARATE message. The prior cohort is every roster member whose handle decomposes to the same `(slug, stage)` pair as the new dispatch. Issue the adapter's cooperative-shutdown call and drop them from session memory. **Mandatory at the boundary; backstops, if any, are the adapter's.**

## Same-Stage Conflict Owner Handoff

After `«halt.rebase-conflict»` aborts an owned code-worktree rebase, keep the entity in its current stage and route one opaque reconciliation assignment to the recorded owner. Write a scope-notes file that names the entity, current stage, PR, registered branch/worktree, old base/head, moved base, exact conflict paths, and next owner action. The package is transport only: do not parse, classify, or resolve the conflict, and do not mutate entity frontmatter or Git refs while routing it.

If `«addressable-worker»` exposes a matching live worker, require its entity, stage, branch, and worktree to equal the tuple proven by the original stamped dispatch, then run ordinary `«dispatch.build» --advance` for that same stage with the scope-notes file and forward its emitted prompt through the existing reuse-advance handle. Otherwise run an ordinary fresh dispatch for the same recorded stage with the same scope-notes file and send its emitted spawn envelope through `«worker.spawn»`; do not add `--stamp`, because the registered branch/worktree already exists.

The worker identity, stage agent, branch, and worktree are the only routing inputs. The PR author, Git author, and shared Git credential are never owner signals. A missing proven tuple, a cold or unowned checkout, or a split-root state-sync conflict is report-only and does not enter this route. This is a per-entity hold, so unrelated entities may continue after the dispatch.

## «reuse.model-match»: stamped model matches the next stage's declared model

- **guard:** skip when `next_stage.effective_model` is null; `«worker-identity»` stamps the null case and null-declared stages accept any reused worker.
- **effect:** resolve the worker's `«worker-identity»`-stamped model via the runtime's model-for-member lookup and compare it to `next_stage.effective_model` in `«worker-identity»`'s canonical model space.
- **block:** a captain-session fallback value — outside that canonical space — never matches; it forces a one-time fresh dispatch that re-stamps a canonical value.
- **done-when:** the models match (or the declared model is null); a mismatch-forced fresh dispatch emits the captain-visible diagnostic `reused worker {name} model {X} does not match next stage effective_model {Y} — fresh-dispatching` verbatim.
- → **prose** — deterministic comparator, no binary; the per-host model space and fallback shapes are `«worker-identity»`'s realization.

## Worktree Ownership

- For worktree-backed entities, active stage/status/report/body state — including `### Feedback Cycles` entries — lives in the worktree copy.
- `pr:` mirrors on `main` for startup/discovery.
- Ordinary active-state writes (`implementation -> validation`) do not land on `main`.

### Split-Root Worktree Contract

When the workflow is split-root (README declares `state:` checkout, e.g. `state: .spacedock-state`), a worktree stage isolates **the deliverable work product only**. Entities live in a separate, non-branched state checkout a main-repo worktree does not contain; the entity body and stage reports are committed there at the entity's state-checkout path, **never** a worktree copy — the dispatch helper hands workers that path even under a worktree stage. The worktree's "commits MUST be on this branch" applies to deliverable artifacts only.

## Worker Resolution

Split worker identity into:
- `dispatch_agent_id` — the spawn call's agent-type parameter. Default `spacedock:ensign`; a stage's `agent: {name}` overrides.
- `worker_key` — filesystem-safe stem for worktree paths (`.worktrees/{worker_key}-{slug}`) and branches (`{worker_key}/{slug}`). Replace `:` with `-` (`spacedock:ensign` → `spacedock-ensign`); a bare name without a namespace equals `dispatch_agent_id`.

## Dispatch Adapter

Runtime adapters bind the capability `«fn»`s below in their `## Runtime implementation` blocks. `«addressable-worker»` is the organizing capability: its presence makes a worker reusable. `«worker.spawn»` handles initial dispatch; the reuse-advance handle only advances a reused agent.

## «worker.spawn»: create the initial worker from a `«dispatch.build»` artifact

The runtime adapter binds the spawn call, helper-field mapping, model/null handling, and host transport metadata.

- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

## «addressable-worker»: address a still-running worker and hear from it mid-run

- **block:** ABSENT → the `«addressable-worker»` reuse condition fails; fresh one-shot only (return value is the sole completion signal; no mid-run steering, no reusable handle). When PRESENT, `«async-dispatch»` must be async — a blocking FO cannot answer a mid-run escalation within the worker's timeout window. Each host's worker↔FO routes, presence probe, and ABSENT teardown are its adapter's binding.
- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

## «async-dispatch»: dispatch a worker without blocking the FO event loop

Resolved per adapter as ASYNC (the spawn returns immediately) or BLOCKING (the dispatch call returns only at completion); a host may bind different modes for different dispatch surfaces.

- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

## «worker-identity»: the spawn-time metadata that keeps a worker addressable and reuse-comparable

Records worker label, substrate, run/session handle, worker address, entity slug, stage, state, completion epoch, stamped model, and the host canonical model space used by `«reuse.model-match»`. Each host's model-resolution stamps the model, including the `«dispatch.build»` null case; each host's handle/address shape and canonical model space are its adapter's binding.

- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

## «completion-signal»: the signals that trigger the completion-verify path

The adapter enumerates which observations count as a completion signal on its host; nothing outside that enumeration completes a worker, and the entity-file stage report remains the verify gate in every case.

- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

## «worker.shutdown»: cooperatively close a terminal or superseded worker

Runs at terminal, supersede, or fresh-dispatch cleanup boundaries after any required preservation message. If ABSENT, record the worker closed in FO memory only after completion or explicit supersede state makes that safe. The runtime adapter owns the concrete shutdown / no-op / in-memory-closure binding and any host-specific preservation channel.

- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

## «context-budget»: probe whether a completed worker is still under context budget for reuse

- **block:** an over-budget or unavailable reading forces fresh dispatch (fail-safe). ABSENT satisfies the `«context-budget»()` reuse condition.
- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

## «roster-reconcile»: sweep the host roster for drift before dispatching

- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

## «dispatch.build»(): assemble the initial-dispatch artifact the spawn call consumes

The ONLY initial-dispatch path: route input through `spacedock dispatch build`, forward its output to `«worker.spawn»` verbatim. Manual prompt/`name` assembly is a protocol violation outside the break-glass block.

- **guard:** write fragile inputs (checklist, scope notes, feedback context) to files first — one checklist item per non-empty line — so Markdown/backticks/shell-vars survive shell quoting.
- **effect:** run the helper, then forward its stdout fields to `«worker.spawn»` unchanged:
  ```
  ${SPACEDOCK_BIN:-spacedock} dispatch build \
    --workflow-dir {workflow_dir} \
    --entity-path {entity_file_path} \
    --stage {target_stage_name} \
    --checklist-file {checklist_file} \
    [--scope-notes-file {scope_notes_file}] \
    [--feedback-context-file {feedback_context_file}] \
    [--bare-mode] \
    [--feedback-reflow] \
    [--advance] \
    [--stamp]
  ```
  `host` derives from the runtime (`--host` is for tests/cross-host tooling only). Select the dispatch transport shape only from the already-bound runtime adapter and invoke that shape once; never probe a second shape after a refusal. In particular, Codex fresh dispatch is named and omits `--bare-mode`. `--bare-mode` reads from live team state only when the active adapter calls for it, never inferred from the stage. Add `--feedback-reflow` only when routing a rejection back to its `feedback-to` target stage. Add `--advance` when advancing a reused live worker instead of spawning one: the emitted envelope carries no spawn/transport fields (nothing is spawned; the adapter enumerates them) and `prompt` is the reuse-advance pointer message, forwarded to the reuse-advance handle instead of `«worker.spawn»`; `--advance` is incompatible with `--bare-mode`. Add `--stamp` for an ordinary (non-reuse) dispatch into a gate-consumed or freshly-advanced entry: it folds the `started`/`worktree=` frontmatter stamps, the state commit+sync, and worktree creation into this same call, before assembly; it refuses (no mutation) unless the entity's status already equals `--stage`, and is incompatible with `--advance` (a reuse advance presupposes an already-stamped live worker).
  Feedback context is opaque transport: preserve the authorized finding, evidence, workflow classification, disposition, workflow-defined correction projection, and assignment bytes.
- **done-when:** on exit 0, `«worker.spawn»` is called with every helper-emitted field — the spawn/transport fields the adapter enumerates plus `description`/`model`/`prompt` — forwarded unchanged. `description` is REQUIRED. `prompt` is the ~175-char file-pointer the ensign Reads on first action — do not strip or rewrite it. Null `model` is `«worker-identity»`'s per-host case, not a core omit-on-null.
- **block:** on non-zero exit (or missing binary), read stderr FIRST: a `dispatch build --stamp:`-prefixed diagnostic is a stamp/sync failure — no envelope was emitted and no authority burned; remedy the named problem and rerun the same build, never break-glass. Exit 3 is `«halt.rebase-conflict»` — HALT dispatch entirely (shared-core HALT clause); never manually dispatch against a halted or unsynced state tree. Only an assembly failure (nonzero WITHOUT the `--stamp:` prefix, or missing binary) triggers the adapter's Break-Glass Manual Dispatch template.
- → **shipped**: `` `spacedock dispatch build` `` — invoke it directly per the effect above.

When `«async-dispatch»` blocks, dispatch one entity at a time and process each completion inline.

## Gate Record/Consume Implicit Sync

In a split-root workflow, `gate record`'s close and `gate consume` commit and sync their own write before returning — no separate `«state.commit»` call. Branch on the FINAL `sync=... phase=record|consume` line plus the exit code, never on which prose lines printed. A refusal path (nothing written) runs no sync and emits no sync line.

| Landing position | Final line | Exit |
|---|---|---|
| close refused (validation, nothing durable) | no `recorded` line, no sync line | 1 or 2 |
| close durable, record-sync failed | `sync=failed phase=record` | 1 |
| close durable, record-sync rebase conflict | `sync=halted phase=record` + HALT stderr | 3 |
| close durable+synced, consume refused (ineligible/blocked, no write) | `sync=... phase=record` + consume diagnostic | 1 |
| close durable+synced, consume stale (superseded write, synced) | `sync=... phase=consume` | 1 |
| advanced, consume-sync failed | `sync=failed phase=consume` | 1 |
| advanced, consume-sync rebase conflict | `sync=halted phase=consume` + HALT stderr | 3 |
| advanced and synced (or terminal `route=`, unspent — no sync) | `sync=pushed\|local-only\|no-op phase=consume` | 0 |

`gate record --consume` sequences the record and consume handlers in one call: close, sync, consume, sync. It requires `--decision approve` for a chat-source close (a chat-source revise/hold is a usage error, exit 2, before any write); a room-source close reports the close and skips consume (`consume=skipped`, exit 0) when the room's decision resolves to revise/hold.

**Recovery.** A durable-but-unsynced landing (`sync=failed`, or `sync=halted` after the HALT's manual resolution) recovers with `«state.commit»(slug)` — the standalone verb publishes whatever is locally committed but unpushed, whichever command wrote it — then resume from the durable position: a pending approval runs standalone `gate consume`; a consumed one runs `dispatch build --stamp`. Never re-run the failed gate verb or the `--consume` composite: record refuses a re-close (frozen closed) and consume refuses a re-consume (byte-clean refusal, no new commit).

## Event Loop

After each agent completion, run `«dispatch.next-action»()` (skeleton below).

**Scheduler envelope and priority.** Machine `status --next --json` returns one object: `{"command":"next","dispatchable":[…],"ready_gates":[…]}`. `ready_gates` is the same ordered four-key index emitted by boot identify; human `--next` is unchanged. On every read, process the first ready-gate row before any dispatchable row. `needs-preparation` is a gate candidate, never entered-stage dispatch.

For `needs-preparation`, load `spacedock:fo-gate-lifecycle`, re-read the entity, latest exact-stage report/checklist, and path-scoped clean commit, then decide semantic completeness. At an `initial: true` stage there is no prior report: the committed clean seed is the reviewed artifact, the structured checklist/AC reads are skipped, and `report-incomplete` never applies. A defect stops once with `report-incomplete: <concrete reason>` and zero prepare, legacy bind, state mutation/commit, presentation, idle, or repeat-next effects. On a semantic pass, supply only question, one committed Markdown Artifact, summary, and References to exactly one existing `gate prepare`; require emitted `room`, `briefing`, `digest`, `state=open`, commit its binding, and re-read the same envelope. Present only one same-slug `awaiting-captain`; any nonzero/mismatched result stops without retry. Existing open, withdrawn, stale, closed, pending, revise, hold, blocked, feedback, consumed, superseded, not-applicable, terminal, archived, malformed, and mismatched states retain their owners and never route through this candidate.

An exact prior-stage replay candidate may exercise `gate prepare`'s existing selection/replay idempotency during that one invocation; FO creates no attempt counter, retry token, cache, or alternate authority.

These are FO-internal scheduling reads — consume them as `--json` (compact, byte-stable, every value a string), not the padded human table a token proxy can mangle; `--fields` narrows to the keys needed. Envelopes: `status`/`--where` → `{"command":"status","entities":[…],"pagination":{…}}`; `--next` → `{"command":"next","dispatchable":[…],"ready_gates":[…]}`. The captain-facing state display (shared-core) still forwards the human table verbatim.

## «dispatch.next-action»(): pick the next event-loop action — dispatch a ready entity, resume a block, or end the iteration

Run every iteration in this order; an empty dispatch projection is never by itself an idle decision:

1. **Reconcile and drain.** Invoke `«roster-reconcile»()` when PRESENT, then drain inbound worker messages when `«addressable-worker»` is PRESENT. Reply to `need_decision` / `interview_request`; acknowledge `progress_update`. An absent adapter capability is omitted, never inferred as a clean roster.
2. **Route every mod/PR action.** Run `status --where "mod-block !=" --json --fields id,slug,mod-block`. For each row, re-read the blocking mod and resume its pending action. A merged PR records its sentinel and enters the merge guard; open/conflicted work remains held with evidence; never auto-resolve or dispatch a worker from this branch.
3. **Route every ready gate.** Read the boot/status `ready_gates` projection. Present `awaiting-captain`; prepare `withdrawn-awaiting-prepare`; consume/advance `approved-awaiting-advance` and require the successor spawn; route `approved-awaiting-merge` to merge delivery without a worker spawn. A still-pending gate is state work and forbids an idle claim.
4. **Query and dispatch.** Run `status --next --json --fields id,slug`; dispatch every row in `dispatchable` that is not at a declared stop. One held or blocked entity never prevents independent rows from reaching their dispatch or gate action.
5. **Retry empty once.** Only after the first empty `status --next`, invoke `«hooks.run»("idle")` exactly once, invoke `«roster-reconcile»()` when PRESENT, then run one second `status --next`. Dispatch newly released rows. Do not repeat the idle hook in this iteration.
6. **Choose the truthful stop.** After the second empty result, continue routing any mod/PR, gate, or other state work already found. With none remaining, use the runtime's wait binding only for an active unresolved worker; completed, errored, and absent worker sets do not qualify. Otherwise report `no-dispatchable` and end the iteration.

The two scheduler reads consume the same `dispatchable+ready_gates` envelope. If both arrays are empty on the first read, idle runs once; after the second read apply ready-gate-first priority again and stop only when both arrays remain empty. A gate made ready by idle must therefore win over dispatch and cannot be reported as false quiescence.

- **done-when:** all actions found in the iteration have been routed and the loop reaches a captain/terminal/consent stop, installs a qualified unresolved-worker wait, or reports `no-dispatchable` after the single retry.
- → **prose** (deterministic mechanism, binary pending — NOT judgment-owned), becomes `` `spacedock dispatch next-action` `` — no driver binary backs it yet; the FO hand-follows the deterministic skeleton above and does not probe for the unshipped command.

**Consent stop before dispatching new standing enforcement.** A ready action whose deliverable is a NEW STANDING check or enforcement process the FO ORIGINATED — a lint, a review gate, a CI lane, a recurring validation step — is the last resort of the boot-resident ordering, never obvious reversible work. Consent is already given for a deliverable the captain commissioned, for running a check that already exists, and for writing a test for the behavior in hand; none of those stop. Otherwise do not dispatch: hold the entity ready-but-undispatched, surface it as an unmet clarification, and carry `awaiting-consent: {slug}` as the iteration's stop reason — headless EXITs with it, interactive waits.

**Fan-out checkpoint.** A SINGLE investigation that will dispatch more than one worker (a flake chase, a review-rework, a refactor sweep) declares BEFORE its first spawn how many WORKERS it expects, the tolerance, and why that spend is economically reasonable. The checkpoint fires when the next spawn would take the running count past expected-plus-tolerance: stop, surface the count against the declaration, and let the captain re-cap rather than spawning again. It binds equally to a plan that commits the fan-out in one act — a workflow script or batch spawn declares the same numbers before launch. Judgment against a declared number, not a counter binary. Collapse demonstrably-identical findings in a barrier stage BEFORE the per-finding verifier spawn — never spend a verifier per duplicate. Where `«async-dispatch»` is async, a per-member verify that fires as reviews land forfeits that barrier; batch, dedupe, then fan out. Author the fan-out's shape against this host's `«async-dispatch»`/`«addressable-worker»` bindings — a blocking or ABSENT binding makes a streaming per-member shape unexecutable; plan batch review instead.

**A second verifier attacks an unowned claim; it never re-runs an owned check.** "Independent adversarial verification" does not justify a second agent to re-run a green check or second-opinion a deterministic fact a shipped check owns — run that check instead. It DOES justify one to attack a claim no check owns AND that no direct read settles — a judgment, a runtime behavior, a fact not visible in the source; adversarial skepticism and a mandated detached audit are that falsifiable exercise. When a read or diff settles the claim, that read is the exercise, and a second agent re-reading the source is the redundancy this refuses.

Repeat the skeleton after each completion until the captain ends the session or, in single-entity mode, the target entity is resolved.
