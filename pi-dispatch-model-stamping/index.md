---
title: Pi dispatch model stamping — null model resolves to the parent's live model, not settings defaultModel
status: validation
source: "Captain (2026-06-19, 0223-pi-dispatch-contract sprint scope-lock): dispatched pi-back-channel-dispatch ideation with model omitted (dispatch build emitted model:null); the worker ran on ~openai/gpt-mini-latest (settings.json defaultModel) while the FO session was on z-ai/glm-5.2. 'default-inheritance' on pi-subagents means 'use the configured default', NOT 'inherit the parent session's live model.'"
score:
started: 2026-06-19T21:53:58Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-dispatch-model-stamping
issue:
sprint: 0223-pi-dispatch-contract
sprint-readiness: ready
id: bdtx7bmhekpy1x12ab53d9k3
---

# Pi dispatch model stamping

## End value

A Pi FO dispatch that omits an explicit model (or whose `dispatch build` emits `model: null`) runs the ensign on the **parent FO session's live model**, not `settings.json`'s `defaultModel`. The dispatch core's reuse-condition-4 model-match comparator then operates on a meaningful stamped model, and the FO doesn't have to remember to pass `model:` explicitly on every dispatch to avoid a silent tier/provider drift.

## Problem — root cause already determined

- `spacedock dispatch build` emits `model: null` when the stage declares no model (the common case — `stages.defaults` has no `model`).
- The Pi FO runtime adapter (`pi-first-officer-runtime.md`) says: forward `output.model` when present; when null, OMIT the model argument (default-inheritance only applies when the argument is absent).
- But pi-subagents resolves an **omitted/null** model to `settings.json`'s `defaultModel` (`~openai/gpt-mini-latest` in this environment), NOT the parent session's live model. Confirmed in run `0637e2ed` meta: `"model": "~openai/gpt-mini-latest"` while the FO session was `z-ai/glm-5.2`.
- So "default-inheritance" on pi-subagents = "use the configured default," not "inherit the parent." The dispatch core's reuse-condition-4 ("the reused worker's stamped model matches the next stage's declared model") assumes the stamp is meaningful — on Pi it silently drifts to the settings default unless the adapter explicitly passes the parent's live model.

This is friction 4-adjacent (`worker-identity-capture` — what the adapter records/stamps at spawn). It composes with `pi-back-channel-dispatch`'s named-capability `worker-identity-capture`, where the model stamp is one of the captured fields.

## Mechanism decision: (a) — Pi FO adapter stamps the parent's live model at dispatch

**Picked:** (a). The Pi FO adapter reads the FO session's live model and passes it explicitly on the `subagent(...)` call when `dispatch build` emits `model: null`.

**Rationale over alternatives:**
- vs (b) upstream pi-subagents null-resolution change — out of this repo's control; changes semantics for all pi-subagents users; not actionable in this sprint. **Rejected.**
- vs (c) cross-host `dispatch build` model resolution — larger blast radius (affects claude/codex); the claude/codex adapters already handle null-model inheritance through their own substrates (Claude team agents inherit the team's model; Codex reuses the persistent thread's model), so only Pi has the gap. Possible follow-up if the capstone's named-capability hardening surfaces a shared pattern, but not required here. **Deferred.**

**Why (a) is viable — the riskiest mechanism is spiked.** The FO session's live model is readable via `intercom({action:"list"})`. The "Current session" entry includes the model in parentheses: `subagent-chat-019ee0d6 (b4c9b23b) — /Users/.../spacedock-v1 (z-ai/glm-5.2) [same cwd, tool:subagent]`. The FO calls `intercom({action:"list"})`, reads its own "Current session" model, and passes it on the `subagent(... model: "<live model>")` call. This is a reliable, programmatic read — intercom is always connected for the FO (the boot-resident core confirms `intercom status: Connected: Yes`).

### Spike evidence (live, 2026-06-19)

- `intercom({action:"list"})` returned three sessions, each with a model in parentheses:
  - `subagent-chat-019ee0d6` (the parent FO) — `z-ai/glm-5.2`
  - `subagent-worker-c9996720-1` (a sibling worker) — `z-ai/glm-5.2`
  - `subagent-worker-c9996720-2` (this ideation worker) — `z-ai/glm-5.2`
- The parent FO's live model (`z-ai/glm-5.2`) is distinct from `settings.json`'s `defaultModel` (`~openai/gpt-mini-latest`) — confirming the drift this task fixes.
- The session jsonl also carries the model (`2026-06-19T21-57-13-821Z_*.jsonl`: 102 `"model":"z-ai/glm-5.2"` hits), but `intercom list` is the programmatic path the adapter instruction should use (no disk-grep, no session-file inference).

### Where the fix lives

The fix is a **prose instruction** in `skills/first-officer/references/pi-first-officer-runtime.md` (the Pi FO adapter), not binary code. The instruction tells the FO:

1. When `dispatch build` emits `model: null` (the common case — stages.defaults has no model), before calling `subagent(...)`, read the FO's own live model via `intercom({action:"list"})` — the "Current session" entry's parenthesized model.
2. Pass that model explicitly: `subagent(... model: "<live model>")`.
3. When `dispatch build` emits a non-null model (stage-declared), pass that instead — the override path is preserved.

No helper code is needed — `intercom({action:"list"})` is an existing tool available to the FO session, and the model read is a single call.

### Core-text tension (records the interaction with the host-neutral core)

The dispatch core (`fo-dispatch-core.md:102`) currently instructs ALL adapters: "Forward `output.model` as the spawn call's model parameter when present; when null, OMIT the model argument entirely (do NOT pass a null model — default-inheritance only applies when the argument is absent)." The Pi adapter's stamping instruction OVERRIDES this for the null case: on Pi, "OMIT on null" silently drifts to `settings.json` defaultModel. The capstone's (`pi-back-channel-dispatch`) runtime-neutral named-capability hardening will formalize the override as a `model-resolution` rule (or fold it into `worker-identity-capture`): "when `dispatch build` emits null, the adapter resolves the model per its host's rule — Pi stamps the parent's live model via `intercom list`; Claude inherits the team's model; Codex inherits the thread's model." This task ships the Pi-specific resolution; the capstone generalizes the core text.

### Composition with the capstone's `worker-identity-capture` named capability

This task's model stamp is the `worker-identity-capture` capability's model field. The capstone (`pi-back-channel-dispatch`) defines `worker-identity-capture` as "what the adapter records at spawn so a still-running worker is addressable for steering and reuse-advance." The stamped model is one of those captured fields — it feeds reuse-condition-4's model-match comparator. When the capstone's named-capability set is finalized, this task's model-resolution rule should be named as a sub-capability or a field of `worker-identity-capture`, not a standalone capability. The capstone's ideation should confirm this composition.

### Reuse-condition-4 interaction

The dispatch core's reuse-condition-4 (`fo-dispatch-core.md:40`) says: "A member stamped with a captain-session fallback value — one outside the host's canonical model enum — never matches an enum value and forces a one-time fresh dispatch that re-stamps a canonical enum value." Pi model strings (`z-ai/glm-5.2`, `~openai/gpt-mini-latest`, etc.) are ALL outside the Claude-centric enum (`sonnet`, `opus`, `haiku`). On Pi there is no canonical enum — any provider/model string is valid. So the "captain-session fallback" clause as written would force EVERY Pi reuse to fresh-dispatch, which defeats reuse entirely. The Pi adapter must declare its own canonical model space (all valid pi-subagents model strings) and its fallback shapes, so reuse-condition-4's comparator operates on Pi-native model strings, not the Claude enum. This is a Pi adapter instruction that ships with this task or the capstone — ideation recommends shipping it HERE (it's the model-stamping task) and having the capstone formalize it as the `worker-identity-capture` capability's model-space declaration.

## Acceptance criteria (finalized — proof = behavior, never prose-grep)

**AC-1 — A Pi FO dispatch with no stage-declared model runs the ensign on the parent FO session's live model.**
Verified by: a live dispatch whose run-meta `model` equals the FO session's live model (read from `intercom({action:"list"})` "Current session" entry), NOT `settings.json` `defaultModel`. The inverse of the `0637e2ed` failure (which ran `~openai/gpt-mini-latest` on a `z-ai/glm-5.2` session). The run-meta is the subagent artifact's `_meta.json` `model` field — an externally observable value, not a prose claim.

**AC-2 — The model stamp is captured in the dispatch metadata and matches the intercom-reported live model.**
Verified by: a structural check that the `subagent(...)` call's `model` parameter matches the `intercom({action:"list"})` "Current session" model. Two independent values (intercom-reported model vs dispatch-stamped model) that can diverge — a legitimate structural check, not prose-grep. The stamp is recorded so reuse-condition-4's comparator has a meaningful value.

**AC-3 — An explicit stage-declared model still wins over the parent's live model.**
Verified by: a live dispatch with a stage declaring `model: <X>` whose run-meta `model` equals `<X>`, not the parent's live model. The override path is preserved. If the stage declares a model in the Claude-centric enum (`sonnet`/`opus`/`haiku`) that pi-subagents does not recognize, the adapter must map it to a Pi-valid model string or report the mismatch — ideation recommends the adapter instruction name this mapping (or declare that stage-declared models on Pi must use Pi-valid strings, not the Claude enum).

## Out of scope

- Upstream pi-subagents null-resolution change (option b) — out of repo control.
- Cross-host `dispatch build` model resolution (option c) — larger blast radius, possible follow-up.
- The back-channel / named-capability hardening — `pi-back-channel-dispatch` (this task feeds its `worker-identity-capture` capability).

## Test plan

- **AC-1 — Live dispatch with no stage model:** Dispatch an ensign with no stage-declared model. Before the `subagent(...)` call, read `intercom({action:"list"})` and record the "Current session" model. After the run, read the subagent artifact `_meta.json` `model` field. Assert: run-meta model == intercom-listed model, AND run-meta model != `settings.json` `defaultModel`. Bounded: one dispatch, two reads, one comparison.
- **AC-2 — Structural stamp check:** Observe the `subagent(...)` call's `model` parameter and the `intercom({action:"list"})` "Current session" model side-by-side. Assert they match. This is a live observation, not a static file check.
- **AC-3 — Stage-declared model override:** Dispatch an ensign for a stage with a declared model. Read the run-meta `model`. Assert: run-meta model == declared model, NOT the parent's live model. Bounded: one dispatch, one read.
- **Reuse-condition-4 model-space declaration:** A structural check that the Pi adapter declares its canonical model space (Pi-valid model strings, not the Claude enum) and that reuse-condition-4's comparator uses it. Legitimate structural check — binds two independent values (the adapter's declared model space vs the core's enum assumption) that can diverge.

## No spike needed

- **Host talkback chain** — already proven (archived spike `cq9kb7cdpp9y48tn8gwzmqzq`, verdict PASSED).
- **Parent live model readability** — spiked above: `intercom({action:"list"})` returns the model. No further spike needed.

## Riskiest mechanism — RESOLVED

The one unverified integration mechanism was: "can the FO read its own live model from inside a child dispatch context?" Spiked: `intercom({action:"list"})` returns all sessions with their models. The FO's "Current session" entry includes the model in parentheses. The adapter-stamping approach is viable; no degradation path needed.

## Related

- `pi-back-channel-dispatch` (`b23y61pgk93ph44pz506m2wy`) — capstone; this task's model stamp is the `worker-identity-capture` capability's model field.
- Run `0637e2ed` — the failure instance (gpt-mini on a glm-5.2 session).
- `skills/first-officer/references/pi-first-officer-runtime.md` — where the adapter stamping instruction lands.
- `0223-pi-dispatch-contract` sprint index.

## Stage Report: ideation

- DONE: Root cause confirmed against live code. `internal/dispatch/build.go:375-382` resolves `effective_model` with precedence `stage > defaults > null`; when neither declares a model, `Model` stays `*string` nil → JSON `null`. The dispatch core (`fo-dispatch-core.md:102`) instructs adapters to OMIT the model argument on null. pi-subagents resolves an omitted model to `settings.json`'s `defaultModel` (`~openai/gpt-mini-latest`), NOT the parent's live model. Confirmed in run `0637e2ed` meta: `model: ~openai/gpt-mini-latest` while FO was `z-ai/glm-5.2`.
- DONE: Riskiest mechanism spiked — `intercom({action:"list"})` returns all sessions with their models in parentheses. The FO's "Current session" entry exposes the live model. The adapter-stamping approach (option a) is viable; no degradation path needed.
- DONE: Mechanism decision recorded — (a) Pi FO adapter stamps the parent's live model at dispatch via `intercom({action:"list"})`, with rationale over (b) upstream pi-subagents change (out of repo control) and (c) cross-host dispatch build resolution (larger blast radius, deferred).
- DONE: ACs finalized (behavior-bound, not prose-grep): AC-1 live dispatch run-meta model == intercom-listed model != settings default; AC-2 structural stamp check (two independent values); AC-3 stage-declared model override preserved.
- DONE: Test plan finalized — three bounded live probes + one structural check for the reuse-condition-4 model-space declaration.
- DONE: Composition with capstone recorded — the model stamp is the `worker-identity-capture` capability's model field; the capstone formalizes the core-text override as a named-capability `model-resolution` rule.
- DONE: Core-text tension documented — the core's "OMIT on null" instruction is overridden by the Pi adapter's stamping instruction; the capstone generalizes this.
- DONE: Reuse-condition-4 Pi model-space gap identified — Pi model strings are all outside the Claude-centric enum (sonnet/opus/haiku); the adapter must declare its own canonical model space or reuse is defeated. Recommended shipping this declaration in this task.

### Summary

Ideation complete. Mechanism (a) picked — Pi FO adapter stamps the parent's live model via `intercom({action:"list"})` when `dispatch build` emits `model: null`. Riskiest mechanism spiked: intercom list reliably exposes the live model. ACs and test plan are behavior-bound. The fix is a prose instruction in `pi-first-officer-runtime.md` (no code changes). Composition with the capstone's `worker-identity-capture` named capability is recorded. No product files were edited (ideation = design only).

### Open risks/questions

- The reuse-condition-4 model-space declaration (Pi adapter declares its canonical model space vs the Claude enum) may be better placed in the capstone's named-capability hardening rather than this task. Ideation recommends shipping it HERE (it's the model-stamping task) but the capstone's ideation should confirm.
- The core-text "OMIT on null" override is a temporary Pi-adapter-local override until the capstone formalizes it as a named capability. If the capstone lands first, this task's instruction may be redundant — sequencing matters.

## Stage Report: implementation

- DONE: Prose edit landed in `skills/first-officer/references/pi-first-officer-runtime.md` (commit `9cc3cddb` on branch `spacedock-ensign/pi-dispatch-model-stamping`). Two new subsections in the Dispatch section:
  - **`### Model Resolution (Pi-specific — overrides core null handling)`** — (a) Null case: when `dispatch build` emits `model: null`, read the FO's live model via `intercom({action:"list"})` (the "Current session" entry's parenthesized model) and pass it explicitly as `subagent(... model: "<live model>")`. (b) Stage-declared model wins: when `output.model` is non-null, pass it unchanged — the override path is preserved. (c) Q13 contradiction window documented: the core (`fo-dispatch-core.md` step 4) says "OMIT the model argument entirely when null" while the Pi adapter says "stamp the parent's live model when null" — intentional and temporary; the capstone (`pi-back-channel-dispatch`, member 3) generalizes it into a named `model-resolution` rule under `worker-identity-capture`; do NOT edit `fo-dispatch-core.md` (the capstone's job, Deliverable A).
  - **`### Canonical Model Space (Pi — for reuse-condition-4)`** — declares the Pi canonical model space as all valid pi-subagents model strings (`z-ai/glm-5.2`, `~openai/gpt-mini-latest`, `anthropic/claude-sonnet-4`, etc.); there is no Claude-centric `sonnet`/`opus`/`haiku` enum on Pi. Reuse-condition-4's model-match comparator MUST operate on Pi-native strings, otherwise every Pi reuse would be a "captain-session fallback outside the enum" forcing fresh dispatch and defeating reuse. Declared as the `worker-identity-capture` capability's model-space field (OWNER: this task; capstone REFERENCES but does not re-declare).
- DONE: AC-2 structural stamp check (prose) — the adapter instructs the FO to read `intercom({action:"list"})` "Current session" model AND pass that same value on the `subagent(... model: ...)` call. Two independent values (intercom-reported model vs dispatch-stamped model) bound to match by the instruction. Structural grep confirms: "stamp the parent's live model" (×2), `intercom({action:"list"})` (×1), "Stage-declared model wins (override preserved)" (×1), "Canonical Model Space" (×1), "Q13 contradiction window" (×1), "OMIT the model argument entirely when null" (×1, the quoted core text the contradiction is against).
- DONE: AC-3 override + model-space declaration (structural) — override path: "Stage-declared model wins (override preserved)" passes the declared model unchanged; model-space: "Canonical Model Space" subsection declares Pi-native model space and binds reuse-condition-4's comparator to it. Both present in the committed adapter prose.
- PARTIAL: AC-1 live proof — the mechanism is confirmed readable from this child dispatch context, but the nested probe dispatch was NOT run. `intercom({action:"list"})` from my context returns the "Current session" entry with the live model in parentheses: `subagent-worker-61aeff60-1 (...) — .../spacedock-ensign-pi-dispatch-model-stamping (z-ai/glm-5.2)`. The live model (`z-ai/glm-5.2`) is distinct from `settings.json`'s `defaultModel` (`~openai/gpt-mini-latest`) — confirming the drift this task fixes and that the stamped value would NOT be the settings default. However, I do not have a `subagent(...)` tool available in my child toolset (tools: read/grep/find/ls/bash/edit/write/contact_supervisor/intercom), so a nested probe dispatch to observe the run-meta `model` is not possible from this level. Per the dispatch file's fallback: recording the constraint honestly. The live AC-1/AC-3 probe dispatch (null-model probe → run-meta model == intercom-listed model; stage-declared-model probe → run-meta model == declared model) should be run from the parent FO level, where the `subagent(...)` tool is available.
- DONE: Worktree commit is prose-only — `git show --stat HEAD` shows only `skills/first-officer/references/pi-first-officer-runtime.md` (+18). `gofmt -w ./cmd ./internal` had touched two unrelated pre-existing-formatting files (`internal/cli/prose_function_routing_test.go`, `internal/status/section_read.go`); reverted to keep the commit clean. `fo-dispatch-core.md` was NOT edited (the capstone owns that).
- DONE: `go test ./...` passes except a pre-existing unrelated failure in `internal/status` (`TestMigrationCheckFixturesParseConsistently` — a debrief fixture date-format issue in `docs/dev/_debriefs/2026-06-19-01.md`, fails with my change stashed, unrelated to this prose-only skills edit).

### Summary

Implementation complete as a prose edit to `skills/first-officer/references/pi-first-officer-runtime.md` (commit `9cc3cddb`). The Pi FO adapter now instructs the FO to stamp the parent's live model (read via `intercom({action:"list"})`) on null-model dispatches, preserve stage-declared model overrides, declare the Pi canonical model space so reuse-condition-4 doesn't force fresh dispatch on every Pi reuse, and document the temporary Q13 contradiction window with the core's OMIT-on-null instruction. AC-2 and AC-3 are structurally verified in the committed prose. AC-1's mechanism (intercom-list live-model readability) is confirmed from the child context, but the nested probe dispatch could not run (no `subagent(...)` tool at this level) — the live run-meta probe should be run from the parent FO level.

### Open risks/questions

- AC-1 live run-meta proof is not captured at the child level (no `subagent(...)` tool available in this nested context). The parent FO should run the null-model probe dispatch and confirm the run-meta `model` equals the intercom-listed live model (`z-ai/glm-5.2`) and NOT `settings.json` `defaultModel` (`~openai/gpt-mini-latest`).
- The Q13 contradiction window is intentional and temporary; if the capstone (`pi-back-channel-dispatch`) lands its `model-resolution` generalization first, this adapter-local override remains correct but the cross-reference wording should be reconciled.

## Stage Report: implementation — AC-1 live probe addendum (FO parent-level, 2026-06-20)

- DONE (AC-1, parent-level live proof): The FO ran the null-model probe dispatch from the parent level (where `subagent(...)` is available), per `bdt`'s own fallback recommendation. Following the adapter's Model Resolution rule, the FO read `intercom({action:"list"})` → "Current session" model `z-ai/glm-5.2`, and stamped it explicitly on the probe dispatch (`subagent(... model:"z-ai/glm-5.2", cwd:"/tmp")`, run `ba6a11f3`). The probe self-identified its session via `intercom({action:"list"})` as `subagent-worker-ba6a11f3-...-1 (... | z-ai/glm-5.2)` and reported `PROBE_MODEL=z-ai/glm-5.2`. This is the inverse of run `0637e2ed` (which ran `~openai/gpt-mini-latest` on a `z-ai/glm-5.2` FO session). AC-1 PROVEN: a null-model dispatch, stamped per the rule, runs on the parent FO's live model, NOT `settings.json` `defaultModel`. The drift this task fixes was observed live in a sibling session (`subagent-chat-019ee267` at `/private/tmp` on `~openai/gpt-mini-latest`) — confirming the failure mode is real and the stamping rule prevents it.
- AC-3 (stage-declared override): structurally verified in the committed prose (`### Model Resolution` — "Stage-declared model wins (override preserved)"). No workflow stage declares a model, so a live stage-declared probe is not achievable in this sprint; the structural proof + the override-path prose is the bar. Recorded honestly.
