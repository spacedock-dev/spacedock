# Pi First Officer Runtime

This file defines how the shared first-officer core executes on Pi. The host-neutral dispatch and merge procedures are in `references/fo-dispatch-core.md` / `fo-merge-core.md` (named by the boot-resident core); this file is the Pi parts those defer to.

## Runtime Shape

Pi is a first-class runtime target, but it does not expose Claude Code team-tool signatures. Do not call or ask workers to call Claude team tools. Pi dispatch uses a Pi-native substrate selected by the launch/test harness:

- default: `pi-subagents`, where the parent first officer uses the `subagent(...)` tool to run a bounded ensign assignment and observes the returned result as completion evidence;
- optional: `pi-agent-teams`, where an adapter maps Spacedock lifecycle intents to the `teams` tool actions (`member_spawn` or `delegate`, `message_dm`, `member_shutdown`, `team_done`).

The Spacedock contract talks in terms of dispatch, completion, follow-up, and shutdown. The Pi adapter owns how those lifecycle events map to the active substrate.

## Dispatch

Use `spacedock dispatch build` with `host: "pi"` in the input JSON. The build artifact is the assignment source of truth: it carries the entity slug/name, entity path, workflow directory, target stage, stage definition fetch command, worktree path when applicable, completion checklist, and completion-signal wording. Forward the emitted dispatch file prompt or dispatch file content to the Pi worker without composing a replacement assignment and without rewriting it into Claude syntax. For the first slice, dispatch via the `async-dispatch` capability: `subagent(... async: true)` returns a run id, the FO polls `subagent({action:"status", id:"<run-id>"})`, and `interrupt` / `resume` steer. Async dispatch keeps the FO event loop free to service `contact_supervisor` escalations within the 10-minute window (the `worker-back-channel` capability). `bare_mode: true` foreground blocking is the legacy default and is SUPERSEDED for back-channel dispatch; it remains the break-glass shape only when `async: true` is unavailable. No `team_name` is required.

A Pi first officer may dispatch via `subagent(...)` when that tool is available. The wrapper may add only Pi transport metadata around the build artifact: `context: "fresh"` plus optional human-facing `phase` / `label` values. Those wrapper fields are additive; they must not redefine or replace the builder-derived slug/name, workflow directory, entity path, target stage, worktree path, completion checklist, or completion contract. The child must load the Spacedock ensign skill and Pi ensign runtime adapter before working.

Manual Pi prompt composition is a break-glass fallback only when `spacedock dispatch build` is unavailable, exits non-zero, or the developer is explicitly debugging the dispatch builder. Any fallback must record the builder failure or unavailability reason and include the minimum canonical schema facts so the degraded dispatch is auditable.

For Spacedock stage dispatches through `pi-subagents`, call `subagent(...)` with explicit `context: "fresh"` and `cwd: <resolved repo root>` (sourced from the install-recorded / explicitly-resolved repo path — a working-directory concern, so the ensign's working directory is the repo regardless of the parent's launch cwd). Do not rely on the worker agent's default context; Spacedock stage workers must be seeded by the task prompt/dispatch content, workflow directory, entity file, and completion checklist rather than inherited first-officer transcript context.

For Spacedock stage dispatches through `pi-subagents`, do not use the `subagent(... acceptance: ...)` contract. Put acceptance requirements in the task prompt/dispatch content instead. Spacedock owns the independent implementation-to-validation workflow: the gate is verification via entity stage reports, product/state commits, and independent validation, not same-agent acceptance finalization by the child that did the work.

For the core's standing-injection call: `pi-subagents` has no standing surface (no-op); `pi-agent-teams` MAY map injection to `member_spawn` per the adapter's lifecycle mapping.

### Model Resolution (Pi-specific `model-resolution` rule)

The dispatch core (`fo-dispatch-core.md` step 4 + the `worker-identity-capture` capability) delegates the null-model case to each adapter's `model-resolution` rule — the core no longer hard-codes omit-on-null. **On Pi, omit-on-null is WRONG**: pi-subagents resolves an omitted/null model to `settings.json`'s `defaultModel` (e.g. `~openai/gpt-mini-latest`), NOT the parent FO session's live model. That is a silent tier/provider drift: the ensign runs on the settings default while the FO session is on a different model, and reuse-condition-4's model-match comparator then compares against a meaningless stamped value. Confirmed in run `0637e2ed` meta: `model: ~openai/gpt-mini-latest` while the FO session was `z-ai/glm-5.2`. So the Pi `model-resolution` rule stamps the parent's live model explicitly (below); Claude/Codex adapters instead omit the argument (the team/thread model inherits).

**Null case — stamp the parent's live model.** When `dispatch build` emits `model: null` (the common case — `stages.defaults` declares no model), before calling `subagent(...)`, read the FO's own live model via `intercom({action:"list"})`. The "Current session" entry carries the live model in parentheses, e.g. `subagent-chat-019ee0d6 (b4c9b23b) — /Users/.../spacedock-v1 (z-ai/glm-5.2) [same cwd, tool:subagent]`. Pass that model explicitly on the spawn call: `subagent(... model: "<live model>")`. `intercom` is always connected for the boot-resident FO, so this is a single reliable programmatic read — no disk-grep, no session-file inference.

**Stage-declared model wins (override preserved).** When `dispatch build` emits a non-null `output.model` (the stage or `stages.defaults` declares a model), pass that value on the `subagent(... model: "<declared model>")` call unchanged. The stage-declared model overrides the parent's live model — the override path is the same as on other hosts. Stage-declared models on Pi MUST use Pi-valid model strings (see the model-space declaration below), not the Claude-centric enum.

**Q13 — resolved.** The core's `model-resolution` rule (under `worker-identity-capture`) now delegates the null case to each adapter: Pi stamps the parent's live model via `intercom({action:"list"})`; Claude inherits the team's model; Codex inherits the thread's model. The earlier "Q13 contradiction window" (core OMIT-on-null vs Pi stamp) was the planned transition state before the capstone generalized the rule; it is now closed — the core and this adapter agree.

### Canonical Model Space (Pi — for reuse-condition-4)

The dispatch core's reuse-condition-4 (`fo-dispatch-core.md` condition 4) says a member stamped with a captain-session fallback value — one outside the host's canonical model enum — never matches an enum value and forces a one-time fresh dispatch. The core's enum assumption is Claude-centric (`sonnet`/`opus`/`haiku`). **Pi has no such enum.** On Pi, any provider/model string is a valid model: `z-ai/glm-5.2`, `~openai/gpt-mini-latest`, `anthropic/claude-sonnet-4`, `google/gemini-2.5-pro`, etc. There is no Claude-centric `sonnet`/`opus`/`haiku` enum on Pi.

This adapter DECLARES the Pi canonical model space as: all valid pi-subagents model strings (provider-qualified or `~`-prefixed provider-relative strings). Reuse-condition-4's model-match comparator MUST operate on these Pi-native strings, not the Claude enum. Otherwise every Pi reuse would be a "captain-session fallback value outside the enum" and force fresh dispatch — defeating reuse entirely. The stamped model from the Model Resolution rule above is a Pi-native string in this declared space, so reuse-condition-4's comparator has a meaningful value to match against `next_stage.effective_model`.

This model-space declaration is the `worker-identity-capture` capability's model-space field (OWNER: this task). The capstone (`pi-back-channel-dispatch`) REFERENCES it but does not re-declare it — see the entity body's composition section.

## Awaiting Completion

For `pi-subagents`, completion is recognized via the `completion-signal` capability (PRIMARY: subagent return + optional advisory; file-verify gate). The primary completion signal is the subagent return value (`subagent({action:"status"})` reports `status: completed`). An ensign MAY additionally send a done-advisory via raw `intercom send` to the FO before returning (a non-blocking heads-up, not a completion carrier — `contact_supervisor` offers no completion reason, only `need_decision`/`interview_request`/`progress_update`). After the primary signal arrives, read the entity file and verify the stage report exactly as the shared core requires. The file-verify is the gate — neither the return nor the advisory alone advances state. Do not advance state based only on a cheerful worker summary.

Verifying the stage report is not the end of the parent's turn. Once the report is verified for a non-gated, non-terminal stage, the parent MUST continue the shared `## Completion and Gates` lifecycle in the same turn — advance the entity and dispatch the next stage (a fresh subagent when the next stage is `fresh: true`) — and only then return its final response. It does not return a completion-only result for the captain to resume from, unless the shared core's halt spans apply (next stage gated, terminal, blocked, or awaiting a captain decision).

For `pi-agent-teams`, completion is observed through the adapter's task/member notification and then verified against the entity file. The adapter should expose a clear completed/failed result to the first officer; the entity stage report remains the source of truth.

## Follow-up and Reuse

The Pi adapter now declares `worker-back-channel` PRESENT, so reuse-condition-1's live-handle requirement is satisfiable. Fresh redispatch nonetheless remains the default safe behavior for the first Pi slice — reuse-advance over the kept-alive back-channel handle (sending the next-stage assignment to a still-running worker) is friction 9 and is DEFERRED. Normal follow-up and retry dispatches are fresh assignment cycles, not context resumes. Record the `worker-identity-capture` metadata at spawn to prevent stale reuse mistakes: worker label, substrate, run id, intercom target, entity slug, stage, state, completion epoch, and stamped model. A follow-up assignment must increment the epoch, and a previous completion must never satisfy the new epoch.

A non-fresh resume is only allowed as an explicit manual/debug exception. Mark the dispatch visibly as a manual/debug resume and tie it to durable metadata in the entity stage evidence, including the full `worker-identity-capture` schema.

## Capability implementations

The Pi adapter binds each named capability the core declares (`fo-dispatch-core.md` `## Named Capabilities`) to concrete Pi tools. The back-channel is available when `pi-subagents` supplies child bridge metadata (the `contact_supervisor` tool is injected into the child).

- `worker-back-channel` — PRESENT. Worker→FO via `contact_supervisor` (`reason: "need_decision"` blocking with a 10-minute timeout, `reason: "interview_request"` blocking structured with the same window, `reason: "progress_update"` non-blocking). FO→worker via `intercom({action:"send", to:"<intercom-target>", message:"..."})` for steering/advance and `intercom({action:"reply", message:"..."})` for replying to a pending ask. Single-pending: intercom allows one pending `ask` per session, so concurrent `need_decision` from multiple workers must serialize/queue (friction 7 deferred).
- `async-dispatch` — ASYNC. Dispatch via `subagent(... async: true)` which returns a run id. Poll with `subagent({action:"status", id:"<run-id>"})`. Interrupt with `subagent({action:"interrupt", id:"<run-id>"})`. The FO event loop polls status between steps; because dispatch is async, the FO is free to service `contact_supervisor` escalations within the 10-minute window. Pass `cwd: <resolved repo root>` on every `subagent(...)` call (sourced from the install-recorded / explicitly-resolved repo path) so the ensign's working directory is the repo.
- `inbound-message-service` — PRESENT. The FO checks `intercom({action:"pending"})` at each event-loop iteration (step 0.5). On a `need_decision`, reply within 10 minutes via `intercom({action:"reply", message:"..."})`. On a `progress_update`, read and acknowledge (no reply required). On an `interview_request`, reply with the provided JSON shape.
- `worker-identity-capture` — schema recorded at spawn: worker label (`worker`), substrate (`pi-subagents`), run id (from `subagent(... async:true)` return), intercom target (`subagent-worker-{runId}-1`, constructed from the run id; 1-indexed child), entity slug, stage, state, completion epoch, stamped model. The stamped model and the Pi canonical model space are OWNED BY member 2 (`pi-dispatch-model-stamping`); this binding REFERENCES that declaration and does not re-declare the model set. When `dispatch build` emits a null model, the Pi adapter stamps the parent's live model per the `model-resolution` rule — member 2 owns the stamping instruction.
- `completion-signal` — PRIMARY (subagent return) + optional advisory (raw intercom send); file-verify gate. The primary completion signal is the subagent return value (`subagent({action:"status"})` returns `status: completed`). An ensign MAY additionally send a done-advisory via raw `intercom send` before returning (a non-blocking heads-up, not a completion carrier — `contact_supervisor` offers no completion reason, only `need_decision`/`interview_request`/`progress_update`). In every case the entity-file stage report is the gate — neither the return nor the advisory alone advances state.
- `context-budget-probe` — NONE. The FO treats the probe as unavailable and prefers fresh dispatch unless the current task explicitly depends on the previous worker's context.
- `roster-reconcile` — NONE. No roster source; the shared-core boundary steps are the only enforcement.

## Shutdown

This is the Pi terminal teardown — fo-merge-core.md's Merge-and-Cleanup step 10, mandatory at the terminal boundary whether the merge ran locally or via a PR host.

For `pi-subagents`, a completed child invocation needs no mailbox shutdown. Mark the worker complete/closed in first-officer memory and continue.

For `pi-agent-teams`, use the adapter's lifecycle mapping to request member shutdown or end the team run. Do not emulate Claude team deletion.

## Live Harness Isolation

Live Pi tests should run with an isolated Pi config directory and an isolated session directory. The harness may copy the operator's existing Pi auth file into the isolated config directory so OAuth/subscription credentials are reused without sharing global sessions, packages, or settings. This mirrors the Codex live runner pattern: isolate runtime state, reuse credentials only.

The durable proof for Pi support is not transcript phrasing. A valid live proof dispatches a Pi ensign against a temp split-root workflow and verifies exit code, state checkout file changes, git log, and stage report content.
