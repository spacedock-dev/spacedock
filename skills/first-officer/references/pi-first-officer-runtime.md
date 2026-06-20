# Pi First Officer Runtime

This file defines how the shared first-officer core executes on Pi. The host-neutral dispatch and merge procedures are in `references/fo-dispatch-core.md` / `fo-merge-core.md` (named by the boot-resident core); this file is the Pi parts those defer to.

## Runtime Shape

Pi is a first-class runtime target, but it does not expose Claude Code team-tool signatures. Do not call or ask workers to call Claude team tools. Pi dispatch uses a Pi-native substrate selected by the launch/test harness:

- default: `pi-subagents`, where the parent first officer uses the `subagent(...)` tool to run a bounded ensign assignment and observes the returned result as completion evidence;
- optional: `pi-agent-teams`, where an adapter maps Spacedock lifecycle intents to the `teams` tool actions (`member_spawn` or `delegate`, `message_dm`, `member_shutdown`, `team_done`).

The Spacedock contract talks in terms of dispatch, completion, follow-up, and shutdown. The Pi adapter owns how those lifecycle events map to the active substrate.

## Dispatch

Use `spacedock dispatch build` with `host: "pi"` in the input JSON. The build artifact is the assignment source of truth: it carries the entity slug/name, entity path, workflow directory, target stage, stage definition fetch command, worktree path when applicable, completion checklist, and completion-signal wording. Forward the emitted dispatch file prompt or dispatch file content to the Pi worker without composing a replacement assignment and without rewriting it into Claude syntax. For the first slice, dispatch should normally run with `bare_mode: true`: the Pi subagent result is the completion signal, so no `team_name` is required.

A Pi first officer may dispatch via `subagent(...)` when that tool is available. The wrapper may add only Pi transport metadata around the build artifact: `context: "fresh"` plus optional human-facing `phase` / `label` values. Those wrapper fields are additive; they must not redefine or replace the builder-derived slug/name, workflow directory, entity path, target stage, worktree path, completion checklist, or completion contract. The child must load the Spacedock ensign skill and Pi ensign runtime adapter before working.

Manual Pi prompt composition is a break-glass fallback only when `spacedock dispatch build` is unavailable, exits non-zero, or the developer is explicitly debugging the dispatch builder. Any fallback must record the builder failure or unavailability reason and include the minimum canonical schema facts so the degraded dispatch is auditable.

For Spacedock stage dispatches through `pi-subagents`, call `subagent(...)` with explicit `context: "fresh"`. Do not rely on the worker agent's default context; Spacedock stage workers must be seeded by the task prompt/dispatch content, workflow directory, entity file, and completion checklist rather than inherited first-officer transcript context.

For Spacedock stage dispatches through `pi-subagents`, do not use the `subagent(... acceptance: ...)` contract. Put acceptance requirements in the task prompt/dispatch content instead. Spacedock owns the independent implementation-to-validation workflow: the gate is verification via entity stage reports, product/state commits, and independent validation, not same-agent acceptance finalization by the child that did the work.

For the core's standing-injection call: `pi-subagents` has no standing surface (no-op); `pi-agent-teams` MAY map injection to `member_spawn` per the adapter's lifecycle mapping.

### Model Resolution (Pi-specific — overrides core null handling)

The dispatch core (`fo-dispatch-core.md` step 4) instructs ALL adapters: when `dispatch build` emits `output.model` as null, OMIT the model argument entirely (default-inheritance applies only when the argument is absent). **On Pi this OMIT-on-null behavior is WRONG and this adapter overrides it.** pi-subagents resolves an omitted/null model to `settings.json`'s `defaultModel` (e.g. `~openai/gpt-mini-latest`), NOT the parent FO session's live model. That is a silent tier/provider drift: the ensign runs on the settings default while the FO session is on a different model, and reuse-condition-4's model-match comparator then compares against a meaningless stamped value. Confirmed in run `0637e2ed` meta: `model: ~openai/gpt-mini-latest` while the FO session was `z-ai/glm-5.2`.

**Null case — stamp the parent's live model.** When `dispatch build` emits `model: null` (the common case — `stages.defaults` declares no model), before calling `subagent(...)`, read the FO's own live model via `intercom({action:"list"})`. The "Current session" entry carries the live model in parentheses, e.g. `subagent-chat-019ee0d6 (b4c9b23b) — /Users/.../spacedock-v1 (z-ai/glm-5.2) [same cwd, tool:subagent]`. Pass that model explicitly on the spawn call: `subagent(... model: "<live model>")`. `intercom` is always connected for the boot-resident FO, so this is a single reliable programmatic read — no disk-grep, no session-file inference.

**Stage-declared model wins (override preserved).** When `dispatch build` emits a non-null `output.model` (the stage or `stages.defaults` declares a model), pass that value on the `subagent(... model: "<declared model>")` call unchanged. The stage-declared model overrides the parent's live model — the override path is the same as on other hosts. Stage-declared models on Pi MUST use Pi-valid model strings (see the model-space declaration below), not the Claude-centric enum.

**Q13 contradiction window — intentional and temporary.** This adapter's "stamp the parent's live model when null" directly contradicts the core's "OMIT the model argument entirely when null" (`fo-dispatch-core.md` step 4). The contradiction is intentional and temporary: this task ships the Pi-specific resolution; the capstone (`pi-back-channel-dispatch`, sprint `0223-pi-dispatch-contract` member 3) generalizes the override into a named `model-resolution` rule under the `worker-identity-capture` capability ("when `dispatch build` emits null, the adapter resolves the model per its host's rule — Pi stamps the parent's live model via `intercom list`; Claude inherits the team's model; Codex inherits the thread's model"). Do NOT edit `fo-dispatch-core.md` to resolve the contradiction here — that is the capstone's job (Deliverable A). A Commander reading both documents during the transition window should see a planned override, not a bug.

### Canonical Model Space (Pi — for reuse-condition-4)

The dispatch core's reuse-condition-4 (`fo-dispatch-core.md` condition 4) says a member stamped with a captain-session fallback value — one outside the host's canonical model enum — never matches an enum value and forces a one-time fresh dispatch. The core's enum assumption is Claude-centric (`sonnet`/`opus`/`haiku`). **Pi has no such enum.** On Pi, any provider/model string is a valid model: `z-ai/glm-5.2`, `~openai/gpt-mini-latest`, `anthropic/claude-sonnet-4`, `google/gemini-2.5-pro`, etc. There is no Claude-centric `sonnet`/`opus`/`haiku` enum on Pi.

This adapter DECLARES the Pi canonical model space as: all valid pi-subagents model strings (provider-qualified or `~`-prefixed provider-relative strings). Reuse-condition-4's model-match comparator MUST operate on these Pi-native strings, not the Claude enum. Otherwise every Pi reuse would be a "captain-session fallback value outside the enum" and force fresh dispatch — defeating reuse entirely. The stamped model from the Model Resolution rule above is a Pi-native string in this declared space, so reuse-condition-4's comparator has a meaningful value to match against `next_stage.effective_model`.

This model-space declaration is the `worker-identity-capture` capability's model-space field (OWNER: this task). The capstone (`pi-back-channel-dispatch`) REFERENCES it but does not re-declare it — see the entity body's composition section.

## Awaiting Completion

For `pi-subagents`, the completion signal is the subagent result returned to the parent. After the result arrives, read the entity file and verify the stage report exactly as the shared core requires. Do not advance state based only on a cheerful worker summary.

Verifying the stage report is not the end of the parent's turn. Once the report is verified for a non-gated, non-terminal stage, the parent MUST continue the shared `## Completion and Gates` lifecycle in the same turn — advance the entity and dispatch the next stage (a fresh subagent when the next stage is `fresh: true`) — and only then return its final response. It does not return a completion-only result for the captain to resume from, unless the shared core's halt spans apply (next stage gated, terminal, blocked, or awaiting a captain decision).

For `pi-agent-teams`, completion is observed through the adapter's task/member notification and then verified against the entity file. The adapter should expose a clear completed/failed result to the first officer; the entity stage report remains the source of truth.

## Follow-up and Reuse

Fresh redispatch is the default safe behavior for the first Pi slice. Normal follow-up and retry dispatches are fresh assignment cycles, not context resumes. If a Pi substrate exposes a resumable worker handle, record only the minimum metadata needed to prevent stale reuse mistakes: worker label, substrate, run/session handle, entity slug, stage, state, and completion epoch. A follow-up assignment must increment the epoch, and a previous completion must never satisfy the new epoch.

A non-fresh resume is only allowed as an explicit manual/debug exception. Mark the dispatch visibly as a manual/debug resume and tie it to durable metadata in the entity stage evidence, including worker label, substrate, run/session handle, entity slug, stage, state, and completion epoch.

## Shutdown

This is the Pi terminal teardown — fo-merge-core.md's Merge-and-Cleanup step 10, mandatory at the terminal boundary whether the merge ran locally or via a PR host.

For `pi-subagents`, a completed child invocation needs no mailbox shutdown. Mark the worker complete/closed in first-officer memory and continue.

For `pi-agent-teams`, use the adapter's lifecycle mapping to request member shutdown or end the team run. Do not emulate Claude team deletion.

## Live Harness Isolation

Live Pi tests should run with an isolated Pi config directory and an isolated session directory. The harness may copy the operator's existing Pi auth file into the isolated config directory so OAuth/subscription credentials are reused without sharing global sessions, packages, or settings. This mirrors the Codex live runner pattern: isolate runtime state, reuse credentials only.

The durable proof for Pi support is not transcript phrasing. A valid live proof dispatches a Pi ensign against a temp split-root workflow and verifies exit code, state checkout file changes, git log, and stage report content.
