---
title: Pi dispatch model stamping — null model resolves to the parent's live model, not settings defaultModel
status: ideation
source: "Captain (2026-06-19, 0223-pi-dispatch-contract sprint scope-lock): dispatched pi-back-channel-dispatch ideation with model omitted (dispatch build emitted model:null); the worker ran on ~openai/gpt-mini-latest (settings.json defaultModel) while the FO session was on z-ai/glm-5.2. 'default-inheritance' on pi-subagents means 'use the configured default', NOT 'inherit the parent session's live model.'"
score:
started: 2026-06-19T21:53:58Z
completed:
verdict:
worktree:
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

## Approach (candidate fixes — ideation confirms and picks)

- **(a) The Pi FO adapter stamps the parent's live model at dispatch.** The adapter reads the FO session's current model (from the session/status) and passes it explicitly on the `subagent(...)` call when `dispatch build` emits `model: null`. The FO doesn't have to remember; the adapter owns it. Recommend.
- **(b) pi-subagents changes null-resolution to inherit the parent's live model.** Upstream pi-subagents change; out of this repo's control, and changes semantics for all pi-subagents users. Reject for this sprint.
- **(c) `dispatch build` resolves the model itself** when the stage declares none — emits the parent's live model rather than null. Cross-cutting (affects claude/codex too); larger blast radius. Possibly a follow-up.

Ideation picks one (recommend (a) — adapter-owned, Pi-scoped, non-breaking), records the decision, and plans verification. The fix likely lives in the Pi FO runtime adapter instruction + a verification that the stamped model matches the parent.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — A Pi FO dispatch with no stage-declared model runs the ensign on the parent FO session's live model.**
Verified by: a live dispatch whose run-meta `model` equals the FO session's live model (read from the pi status bar / session meta), NOT `settings.json` defaultModel. The inverse of the `0637e2ed` failure.

**AC-2 — The model stamp is captured in the worker-identity metadata.**
Verified by: the dispatch/identity metadata records the stamped model, so reuse-condition-4's comparator has a meaningful value. Structural check (two independent values — stamped model vs stage-declared model) is legitimate.

**AC-3 — An explicit stage-declared model still wins.**
Verified by: a stage with `model: sonnet` (or similar) stamps that model, not the parent's. The override path is preserved.

## Out of scope

- Upstream pi-subagents null-resolution change (option b) — out of repo control.
- Cross-host `dispatch build` model resolution (option c) — larger blast radius, possible follow-up.
- The back-channel / named-capability hardening — `pi-back-channel-dispatch` (this task feeds its `worker-identity-capture` capability).

## Test plan

- Live dispatch with no stage model (AC-1) — run-meta model == parent live model. Bounded probe.
- Identity-metadata schema check (AC-2) — structural, independent values.
- Stage-declared-model dispatch (AC-3) — run-meta model == declared, not parent.

## Related

- `pi-back-channel-dispatch` (`b23y61pgk93ph44pz506m2wy`) — capstone; this task's model stamp is the `worker-identity-capture` capability's model field.
- Run `0637e2ed` — the failure instance (gpt-mini on a glm-5.2 session).
- `skills/first-officer/references/pi-first-officer-runtime.md` — where the adapter stamping instruction lands.
- `0223-pi-dispatch-contract` sprint index.
