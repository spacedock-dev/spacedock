---
title: Pi back-channel dispatch — declare and wire the worker↔FO back-channel over pi-intercom
status: ideation
source: "Captain (2026-06-19): the pi-intercom supervisor-talkback capability is proven (archived spike cq9kb7cdpp9y48tn8gwzmqzq, PR #301 spike-only) but the Spacedock Pi FO/ensign adapters do not wire it. Treat the capability as already implemented for the host session; find the contract/adapter frictions and harden the dispatch core to runtime-neutral named capabilities."
score:
started: 2026-06-19T17:20:10Z
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness:
id: b23y61pgk93ph44pz506m2wy
---

# Pi back-channel dispatch

## End value

A Pi first officer that dispatches **addressable, steerable, reusable** workers over the pi-intercom back-channel — matching the worker-back-channel capability the host-neutral dispatch core already organizes around — instead of treating Pi as bare/one-shot with the subagent return value as the sole completion signal.

The host capability is **already proven and shipped**: the archived spike `pi-intercom-runtime-capability-probe` (verdict PASSED) demonstrated the full chain — child `contact_supervisor` (`progress_update` + `need_decision`) → FO reply → child resumes → durable marker written — and the `pi-intercom` SKILL documents orchestrator-side handling of all three escalation types (`need_decision` / `interview_request` / `progress_update`). This task does **not** implement host behavior; it wires the proven host capability into the Spacedock dispatch contract and the Pi FO/ensign adapters.

## Problem

The host-neutral dispatch core (`skills/first-officer/references/fo-dispatch-core.md`) organizes the entire dispatch model around one declared adapter capability — the **Worker back-channel**:

> The runtime adapter DECLARES whether it provides a live worker back-channel: a dispatched worker that can message the lead WHILE it is still running, and a lead that can message back to advance, steer, or query it. This declared handle is reuse-condition-1's "live, reusable handle" — the single capability the dispatch model organizes around.

Today the Pi FO adapter (`skills/first-officer/references/pi-first-officer-runtime.md`) and Pi ensign adapter (`skills/ensign/references/pi-ensign-runtime.md`) declare **no back-channel**: Pi is modeled as one-shot foreground blocking, the subagent return value is the sole completion signal, fresh redispatch is the default, and the ensign is told to "return one concise final result… stop. Do not idle waiting for another message." So reuse-condition-1 fails *by adapter declaration*, not by host capability. The capability ships; the contract ignores it.

The dispatch core's back-channel section is also currently written partly in host-specific terms (Claude team registry, Codex `send_input`). To admit the Pi back-channel cleanly, the core must be **hardened to runtime-neutral named capabilities**, where each runtime's adapter implements those capabilities by naming its concrete tools and runtime-specific logic. The frictions below are the concrete gaps.

## Frictions (the wire-up gaps — frictions 1–6 are in scope; 7–9 deferred)

1. **Pi FO adapter declares no back-channel.** The adapter never declares a back-channel, so the dispatch core treats Pi as bare/one-shot. → The adapter must declare the back-channel: worker→FO via `contact_supervisor` (the `need_decision` / `interview_request` / `progress_update` escalations), FO→worker via `intercom` `send`/`reply`/`ask`, and map this declared handle to the reuse-advance handle (reuse-condition-1's "live, reusable handle").

2. **Foreground blocking `subagent(...)` cannot service inbound escalations.** A blocking subagent call occupies the FO thread; a worker's `need_decision` carries a **10-minute timeout** (per the pi-intercom SKILL) and will time out before a blocking FO can answer. `subagent(...)` supports `async: true` plus `status` / `interrupt` / `resume`. → The adapter must specify async dispatch so the FO event loop remains free to answer escalations within the 10-minute window, with `status`/`resume` as the await mechanism.

3. **No event-loop step for inbound worker messages.** The dispatch core event loop (steps 1–4) fires only "after each agent completion" — there is no step for servicing `contact_supervisor` from a *still-running* worker. The host step-0 "roster-reconcile sweep" is for roster reconciliation, not message handling. → The event loop needs a runtime-neutral "service inbound worker messages" step (the adapter names the concrete listen call — `intercom pending` on Pi); the core describes it as a named capability.

4. **Worker intercom identity not captured at spawn.** Escalations arrive as "From subagent-worker-78f659a3-1" with run/agent/child-index metadata. The adapter's reuse-prevention metadata (worker label, substrate, run/session handle, entity slug, stage, state, completion epoch) does not record the intercom address, so the FO cannot *proactively* steer a still-running worker. → Capture the worker's intercom identity at dispatch and carry it in the reuse metadata.

5. **Pi ensign has no talkback protocol.** Claude ensign: "ask for clarification via `SendMessage(to=\"team-lead\")`." Codex ensign: has a Clarification section. Pi ensign: **none** — it says return a final result and stop. So even if the FO listened, the worker would not ask. → Add a clarification/steering section to `skills/ensign/references/pi-ensign-runtime.md` that directs the ensign to use `contact_supervisor` (`need_decision` for blocking clarifications, `progress_update` for non-blocking plan-changing discoveries), and to resume after the reply.

6. **Completion-signal duality.** With a back-channel, "done" can arrive as a `contact_supervisor`/`send` *or* as the subagent return value. The adapter currently recognizes only the latter. The stage report remains the source of truth (the FO verifies against the entity file regardless), but the adapter needs a rule declaring both signals equivalent triggers for the verify path, with the file-verify as the gate.

### Deferred (address separately later — frictions 7–9)

7. **`ask` 10-min timeout + one-pending-ask-per-session vs `concurrency: 3`.** Two workers `need_decision` simultaneously → one blocks/fails (intercom allows one pending ask per session). The dispatch core's `concurrency: 3` assumes a multiplexing back-channel. → Adapter needs a serialization/queueing rule for concurrent decisions, or a Pi concurrency cap for decision-heavy stages. **Deferred.**

8. **Standing `comm-officer` declared no-op on Pi.** The workflow declares a prose-polisher; the Pi adapter says "pi-subagents has no standing surface (no-op)." Intercom makes a long-lived comm-officer session feasible, so the polish capability is unavailable on Pi despite the host supporting it. → Standing injection could be wired via a long-lived intercom-addressed session. **Deferred.**

9. **Reuse structurally off on Pi.** The adapter: "Fresh redispatch is the default… A non-fresh resume is only allowed as an explicit manual/debug exception." With a back-channel, reuse-advance (send the next-stage assignment to the kept-alive worker) becomes possible — but the adapter forbids it. The "fresh only" stance is itself a friction once a back-channel exists. **Deferred** (depends on 1–6 landing first).

## Approach

Two coordinated deliverables:

### A. Harden the dispatch core to runtime-neutral named capabilities

Rewrite the back-channel section of `skills/first-officer/references/fo-dispatch-core.md` (and the reuse/event-loop sections it touches) so the core speaks in **named capabilities** rather than host-specific tools. Each capability is a behavioral contract; each runtime adapter implements it by naming its concrete tools and logic. Candidate named capabilities (ideation finalizes the exact set and names):

- **`worker-back-channel`** (the organizing capability) — declare present/absent. When present, names: (a) the worker→FO escalation call and its message types, (b) the FO→worker advance/steer/query call, (c) whether the channel is multiplexing or single-pending. When absent, the core's one-shot/fresh-only path applies.
- **`async-dispatch`** — whether the spawn call is blocking or returns an awaitable handle; names the await/resume mechanism. Required when `worker-back-channel` is present (a blocking spawn cannot service inbound escalations).
- **`inbound-message-service`** — the event-loop step that drains pending worker messages; names the listen call.
- **`worker-identity-capture`** — what the adapter records at spawn so a still-running worker is addressable for steering and reuse-advance.
- **`completion-signal`** — the set of signals the adapter treats as completion-equivalent (return value, inbound done-message), with the file-verify as the gate.

The Claude, Codex, and Pi adapters each get a short "Capability implementations" subsection naming their concrete tools: Claude (`Agent` background teammate + `SendMessage`, team registry), Codex (`send_input` to a persistent thread), Pi (`contact_supervisor` + `intercom` send/reply/ask, `subagent(... async:true)` + `status`/`resume`). The core's logic references capabilities by name; adapters bind names to tools.

This is itself a contract/scaffolding change, so per the dev-workflow proof policy it must go through implementation + validation, and per the self-evidence-bar principle the live `claude-live` / `codex-live` / `pi-live` drives must be green before merge (the dogfood — the change touches every host adapter).

### B. Wire the Pi back-channel (frictions 1–6)

Implement the Pi adapter's capability bindings: declare `worker-back-channel` present (1), specify async dispatch via `subagent(... async:true)` + `status`/`resume` (2), add the `inbound-message-service` step using `intercom pending` (3), capture worker intercom identity at spawn (4), add the ensign talkback protocol to `pi-ensign-runtime.md` (5), and declare the completion-signal duality with file-verify gate (6).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — The dispatch core describes back-channel dispatch in runtime-neutral named capabilities, with each host adapter binding those names to concrete tools.**
Verified by: a structural check that the core references capabilities by name (not host tool calls), and each adapter (Claude/Codex/Pi) carries a capability→tool binding; a contractlint structural check is legitimate here (it binds two independent values — the core's named capability and the adapter's binding — that can diverge, so it is not prose-grep).

**AC-2 — A Pi FO dispatches an async worker and services an inbound `need_decision` escalation mid-run, replying within the 10-minute window, with the ensign resuming and completing.**
Verified by: a live Pi drive (the `pi-live` lane) where a dispatched ensign hits a seeded ambiguity, `contact_supervisor need_decision`s the FO, the FO replies via `intercom reply`, the ensign resumes and completes; durable evidence in the entity body and the state checkout.

**AC-3 — A Pi FO verifies completion from either signal (subagent return OR inbound done-message), with the entity-file stage report as the gate in both cases.**
Verified by: a live or harness run exercising both signal paths, each followed by the file-verify; the ensign's final result alone never advances state.

**AC-4 — The Pi ensign adapter directs clarifications to `contact_supervisor` and resumes after the reply.**
Verified by: the ensign talkback protocol in `pi-ensign-runtime.md`, exercised by the AC-2 live drive (the ensign's `need_decision` is the proof, not the prose).

**AC-5 — Worker intercom identity is captured at spawn and carried in reuse metadata.**
Verified by: a reuse-metadata schema check (structural, two independent values) plus the AC-2 drive showing the FO can address the still-running worker by captured identity.

## Out of scope

- Frictions 7–9 (concurrency serialization of `ask`, standing `comm-officer` on Pi, reuse-advance on Pi) — deferred, addressed separately, depend on 1–6 landing first.
- Host behavior implementation — the pi-intercom talkback is already proven and shipped; this task does not modify the host.
- The two open PRs (`#397` gate-extract-verbs, `#398` wrong-root-guard) — host-neutral, unrelated.

## Test plan

- **Structural / contractlint** (AC-1, AC-5): capability-name references in the core; capability→tool bindings in each adapter; reuse-metadata schema. These bind independent values and are legitimate structural checks, not prose-grep.
- **Live `pi-live` drive** (AC-2, AC-3, AC-4): a seeded-ambiguity ensign dispatch that exercises the full back-channel round-trip and both completion-signal paths. Bounded: one child, one `need_decision`, one reply, one resume, one completion. Mirrors the archived spike's bounded shape.
- **Live `claude-live` / `codex-live` regression** (AC-1 dogfood): the core rewrite touches every host adapter; the other two live lanes must stay green to prove the runtime-neutral rewrite did not regress a host that already had a back-channel.

## Related

- Archived spike `pi-intercom-runtime-capability-probe` (`cq9kb7cdpp9y48tn8gwzmqzq`, PR #301 closed-as-spike) — proved the host capability this task wires in.
- `skills/first-officer/references/fo-dispatch-core.md` — the dispatch core to harden (deliverable A).
- `skills/first-officer/references/pi-first-officer-runtime.md` — the Pi FO adapter to wire (deliverable B, frictions 1–4, 6).
- `skills/ensign/references/pi-ensign-runtime.md` — the Pi ensign adapter (friction 5).
- `skills/ensign/references/claude-ensign-runtime.md` / `codex-ensign-runtime.md` — the reference clarification protocols Pi ensign lacks.
- Deferred follow-up: a sibling task for frictions 7–9 (file after this lands).

## Stage Report: ideation

- DONE: Spike evidence (live, 2026-06-19) — run id `0637e2ed`; the worker sent `need_decision` during the dispatch, the supervisor reply arrived within the 10-minute window, and the substrate detached the foreground dispatch for intercom coordination (`Detached for intercom coordination`) without explicit `async:true`.
  This is durable evidence that the Pi back-channel works on the current `pi-subagents` substrate while the worker is still running.
- DONE: Contract rewrite remains DRAFT only; do not commit the skill-file edits from this run.
  Supervisor instruction: the draft edits were made against a skill-less / default-model run and are not authoritative. Keep the contract changes out of the repo state for this cycle and re-dispatch cleanly later.
- DONE: In-flight frictions surfaced for the follow-up re-dispatch.
  (a) null `model` resolved to `settings.json` defaultModel rather than the parent run's live model; (b) skill injection failed with `Skills not found: ensign`.

### Summary

Live spike evidence is recorded above and establishes the back-channel / mid-run escalation path on the current Pi substrate. The contract-doc edits are intentionally not committed in this cycle; they remain a draft artifact only, pending a clean re-dispatch under a corrected skill/model context.
