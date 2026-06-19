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
sprint: 0223-pi-dispatch-contract
sprint-readiness: ready
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

Two coordinated deliverables: (A) harden the dispatch core to runtime-neutral named capabilities, and (B) wire the Pi adapter's capability bindings.

### No spike needed

The riskiest mechanism — the async subagent + intercom back-channel round-trip — is **already proven live** by run `0637e2ed` (recorded in the Stage Report below): a foreground dispatch auto-detached for intercom coordination, the worker's `need_decision` reached the FO mid-run, and the FO replied within the 10-minute window. The archived spike `cq9kb7cdpp9y48tn8gwzmqzq` proved the full host talkback chain (`contact_supervisor` `progress_update` + `need_decision` → FO reply → child resume → durable marker). The pi-intercom SKILL documents the orchestrator-side handling of all three escalation types. The named-capability core rewrite is prose-structural (reorganizes existing contract text; no unverified runtime mechanism). The one integration detail — that `subagent(... async:true)` + `status`/`resume` frees the FO event loop to service escalations — is the async form of what `0637e2ed` already demonstrated in the foreground-detached form; the async path is strictly more capable (the FO is never blocked).

### A. Harden the dispatch core to runtime-neutral named capabilities

Rewrite the `## Dispatch Adapter` section's "Worker back-channel capability" subsection and the reuse/event-loop sections it touches in `skills/first-officer/references/fo-dispatch-core.md` so the core speaks in **named capabilities** rather than host-specific tools. Each capability is a behavioral contract the core references by name; each runtime adapter declares which it provides and binds each to concrete tools. No host tool call appears in the host-neutral core.

**Finalized named-capability set (7):**

| Capability | Declares | Core section that references it |
|---|---|---|
| `worker-back-channel` (organizing) | present/absent; when present names (a) worker→FO escalation call + message types, (b) FO→worker advance/steer/query call, (c) multiplexing or single-pending | Dispatch Adapter; reuse-condition-1 |
| `async-dispatch` | blocking or async; when async names the await/resume/interrupt mechanism | Dispatch Adapter (step 9 "wait for the worker result"); required when `worker-back-channel` is present |
| `inbound-message-service` | present/absent; when present names the listen call that drains pending worker messages | Event Loop (new step between step 0 and step 1); required when `worker-back-channel` is present |
| `worker-identity-capture` | the schema recorded at spawn: worker label, substrate, run/session handle, intercom address, entity slug, stage, state, completion epoch, stamped model | Reuse conditions (condition-1 handle + condition-4 model) |
| `completion-signal` | the set of signals treated as completion-equivalent (return value, inbound done-message), with file-verify as the gate | Dispatch Adapter (step 9) + Completion and Gates |
| `context-budget-probe` (already informally named) | present/absent; when present names the probe call | Reuse condition-0 (already references it) |
| `roster-reconcile` (already informally named) | present/absent; when present names the reconcile sweep call and drift classes | Event Loop step 0 (already references it) |

`context-budget-probe` and `roster-reconcile` are already informally named in the core (reuse-condition-0 says "the runtime adapter's context-budget probe"; event-loop step 0 says "the host's step-0 reconcile sweep"). The hardening formalizes them as named capabilities alongside the 5 new ones, so the core is uniformly capability-named.

**Before/after for the core's "Worker back-channel capability" section:**

The current section declares the concept but references host-specific patterns ("a named background worker, a team registry, a subagent mailbox"). The rewrite replaces it with a `## Named Capabilities` section that declares each capability by name with its behavioral contract, and references the adapter for the concrete binding. The two-bullet present/absent shape stays (it is the core's logic); the host-specific examples move to each adapter's `## Capability implementations` subsection.

**Before/after for reuse conditions:** condition-0 already says "the runtime adapter's context-budget probe" → formalize as `context-budget-probe` capability. Condition-1 already says "the runtime adapter's reuse-advance handle" → formalize as `worker-back-channel` + `worker-identity-capture`. Condition-4's "the host's canonical enum" stays adapter-declared. The host-specific parentheticals ("Codex declares none; Claude supplies one") move to the adapters.

**Before/after for the event loop:** step 0 already says "the host's step-0 reconcile sweep" → formalize as `roster-reconcile`. Add a new step 0.5 (or fold into step 0) for `inbound-message-service`: when the adapter declares it, the FO drains pending worker messages (`intercom pending` on Pi) at each iteration before checking dispatchables.

**Adapter bindings** — each FO runtime adapter gets a `## Capability implementations` subsection:

| Capability | Claude | Codex | Pi (this task wires) |
|---|---|---|---|
| `worker-back-channel` | PRESENT: `Agent(run_in_background=true)` + `SendMessage(to="team-lead")` (worker→FO), `SendMessage(to=name)` (FO→worker); multiplexing | PRESENT: `spawn_agent` + mailbox final-status (worker→FO), `send_input` (FO→worker); multiplexing | PRESENT: `contact_supervisor` (worker→FO), `intercom send/reply` (FO→worker); single-pending (friction 7 deferred) |
| `async-dispatch` | ASYNC: `run_in_background=true` returns immediately | ASYNC: `spawn_agent` returns handle; `wait_agent` is explicit foreground wait | ASYNC: `subagent(... async:true)` returns run id; `subagent({action:"status"})` polls; `interrupt`/`resume` steer |
| `inbound-message-service` | PRESENT: `system task_notification` entries + `SendMessage` inbox | PRESENT: mailbox final-status notification | PRESENT: `intercom({action:"pending"})` drains pending asks; reply via `intercom({action:"reply"})` |
| `worker-identity-capture` | agent name + `agentType` on disk + model from team config | task name + mailbox handle | `subagent-worker-{runId}-{childIdx}` intercom target + run id + stamped model (sibling: `pi-dispatch-model-stamping`) |
| `completion-signal` | DUAL: "Done:" inbox message OR `task_notification` OR captain shutdown | DUAL: mailbox final-status notification (sole) | DUAL: subagent return (status:completed) OR inbound done-message via `contact_supervisor`/`intercom send` |
| `context-budget-probe` | PRESENT: `spacedock dispatch context-budget --name {name}` | NONE | NONE |
| `roster-reconcile` | PRESENT: `spacedock dispatch reconcile` (drift classes: lingering/superseded/un-advanced-pr/stale-branch/local-main-drift) | NONE | NONE |

This is a contract/scaffolding change to a high-stakes surface (the shipped FO/ensign contract + host adapters). Per the dev-workflow proof policy, the path→lane mapping is the gate: a change to the host-neutral dispatch core requires `claude-live` AND `codex-live` AND `pi-live` green; a change to the Pi adapter requires `pi-live`. The structural contractlint for capability-name↔adapter-binding is legitimate (binds two independent values that can diverge, not prose-grep).

### B. Wire the Pi back-channel (frictions 1–6)

Implement the Pi adapter's `## Capability implementations` subsection in `skills/first-officer/references/pi-first-officer-runtime.md` with the bindings above, and add the talkback protocol to `skills/ensign/references/pi-ensign-runtime.md`. Concretely:

**Friction 1 — declare `worker-back-channel` PRESENT.** The adapter declares the back-channel: worker→FO via `contact_supervisor` (`need_decision` blocking 10-min, `progress_update` non-blocking, `interview_request` blocking 10-min structured); FO→worker via `intercom({action:"send", to:"<intercom-target>", message:"..."})` for steering/advance and `intercom({action:"reply", message:"..."})` for replying to a pending ask. Single-pending (intercom allows one pending `ask` per session — friction 7 deferred). The back-channel is available when `pi-subagents` supplies child bridge metadata (the `contact_supervisor` tool is injected into the child).

**Friction 2 — declare `async-dispatch` ASYNC.** Dispatch via `subagent(... async: true)` which returns a run id. Poll with `subagent({action:"status", id:"<run-id>"})`. Interrupt with `subagent({action:"interrupt", id:"<run-id>"})`. The FO event loop polls status between event-loop steps; when a worker sends `contact_supervisor`, the FO is free to reply because the dispatch is async. This replaces the current `bare_mode: true` foreground-blocking default.

**Friction 3 — declare `inbound-message-service` PRESENT.** The FO checks `intercom({action:"pending"})` at each event-loop iteration (between status polls and dispatchable checks). When a `need_decision` arrives, reply within 10 minutes via `intercom({action:"reply", message:"..."})`. When a `progress_update` arrives, read and acknowledge (no reply required). When an `interview_request` arrives, reply with the provided JSON shape. This is the new event-loop step 0.5.

**Friction 4 — declare `worker-identity-capture` schema.** At spawn, record: worker label (`worker`), substrate (`pi-subagents`), run id (from `subagent(... async:true)` return), intercom target (`subagent-worker-{runId}-1` — constructed from the run id; 1-indexed child), entity slug, stage, state, completion epoch, stamped model (from sibling task `pi-dispatch-model-stamping`). The intercom target is the FO→worker address for steering and reuse-advance.

**Friction 5 — add ensign talkback protocol.** Add a `## Clarification` section to `skills/ensign/references/pi-ensign-runtime.md`:

> If requirements are unclear or ambiguous, ask for clarification via `contact_supervisor` with `reason: "need_decision"` rather than guessing. Describe what you understand and what's ambiguous so the FO can route a quick answer. The FO replies via `intercom({action:"reply", message:"..."})`; after receiving the reply, resume working.
>
> For non-blocking plan-changing discoveries (a riskier mechanism panning out, a scope boundary discovered), send `contact_supervisor` with `reason: "progress_update"` — the FO reads and acknowledges; no reply is required and you continue working.

This mirrors the Claude ensign's `SendMessage(to="team-lead")` clarification and the Codex ensign's thread-based clarification, adapted to Pi's `contact_supervisor` surface.

**Friction 6 — declare `completion-signal` DUAL.** Completion may arrive as (a) the subagent return value (`subagent({action:"status"})` returns `status: completed`) or (b) an inbound done-message via `contact_supervisor` or `intercom send`. Both signals trigger the same verify path: read the entity file, verify the stage report. The file-verify is the gate — neither signal alone advances state. This replaces the current "subagent return is the sole completion signal" rule.

### Documentation changes

This task changes the FO/ensign contract (skill files) — these ARE the docs. `docs/runtime-support.md` line 188 mentions the Pi readiness check and supervisor-talkback probe but does not reference the dispatch back-channel by name; no doc diff needed there. The dev README does not reference the back-channel. The skill files are the deliverable AND the docs.

## Acceptance criteria (entity-level; proof = behavior, never prose-grep)

**AC-1 — The dispatch core's back-channel, reuse, and event-loop logic references named capabilities (not host tool calls), and each host adapter carries a capability→tool binding for each declared capability.**
Verified by: a structural contractlint test that extracts the named-capability set from `fo-dispatch-core.md` (the `## Named Capabilities` section) and from each adapter's `## Capability implementations` subsection, then compares them as sets — the adapter must bind every capability the core declares. Pattern: `reconcile_class_binding_test.go` (dual-extraction, set comparison, empty-set guard). This binds two independent values (the core's named capability and the adapter's binding) that can diverge, so it is legitimate structural contractlint, NOT prose-grep.

**AC-2 — A Pi FO dispatches an async worker that sends `contact_supervisor need_decision` mid-run, the FO replies within the 10-minute intercom timeout, and the ensign resumes and completes.**
Verified by: a live `pi-live` drive — a dispatched ensign hits a seeded ambiguity, `contact_supervisor need_decision`s the FO, the FO replies via `intercom({action:"reply"})`, the ensign resumes and completes. Durable evidence in the entity body and the state checkout (run id, escalation text, reply text, completion status).

**AC-3 — A Pi FO verifies completion from either signal (subagent return OR inbound done-message), with the entity-file stage report as the gate in both cases.**
Verified by: a live or harness run exercising both signal paths, each followed by the file-verify (`status --read <ref> --json` → last `## Stage Report` section). The ensign's final result alone never advances state.

**AC-4 — The Pi ensign adapter directs clarifications to `contact_supervisor` and the ensign resumes after the FO's reply.**
Verified by: the AC-2 live drive — the ensign's `need_decision` is the behavioral proof, not the prose in `pi-ensign-runtime.md`.

**AC-5 — Worker intercom identity is captured at spawn and the FO can address a still-running worker by that identity.**
Verified by: the AC-2 drive showing the FO addresses the worker by its captured intercom target (`subagent-worker-{runId}-1`), plus a structural schema check that the identity-capture fields are present in the adapter's declared schema.

**AC-6 — The runtime-neutral core rewrite does not regress Claude or Codex back-channel behavior.**
Verified by: `claude-live` and `codex-live` lanes green after the core rewrite (the dogfood — the change touches every host adapter). A red live lane is diagnosed by reading THIS run's failing test, not by inheriting a prior session's label.

## Out of scope

- Frictions 7–9 (concurrency serialization of `ask`, standing `comm-officer` on Pi, reuse-advance on Pi) — deferred to a follow-up sprint, depend on 1–6 landing first. The boundary is clean: friction 7 is the `single-pending` declaration in `worker-back-channel`; friction 8 is the `standing-injection` no-op; friction 9 is the `fresh-redispatch` default. Each has a sharp seam in the named-capability declarations.
- Host behavior implementation — the pi-intercom talkback is already proven and shipped; this task does not modify the host.

## Test plan

| AC | Test type | Cost/complexity | What it proves |
|---|---|---|---|
| AC-1 | Structural contractlint (Go test) | Low — dual-extraction + set comparison, following `reconcile_class_binding_test.go` pattern | The core references capabilities by name; each adapter binds them. Not prose-grep. |
| AC-2 | Live `pi-live` drive | Medium — one async dispatch, one seeded `need_decision`, one reply, one resume, one completion. Bounded, mirrors the archived spike's shape. | The back-channel round-trip works end-to-end on Pi. |
| AC-3 | Live `pi-live` or harness | Medium — both signal paths (subagent return + inbound done-message), each file-verified | Completion-signal duality with file-verify gate. |
| AC-4 | Live `pi-live` drive (AC-2 subsumes) | None additional — AC-2's `need_decision` is the proof | Ensign talkback protocol is exercised, not just prose. |
| AC-5 | Live (AC-2 subsumes) + structural schema check | Low — identity-capture fields present in adapter declaration | FO can address a still-running worker by captured identity. |
| AC-6 | Live `claude-live` + `codex-live` regression | Medium — both lanes must run green after the core rewrite | The runtime-neutral rewrite did not regress a host that already had a back-channel. |

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

## Stage Report: ideation (re-dispatch 2026-06-19)

- DONE: Finalized the named-capability set — 7 capabilities (5 new: `worker-back-channel`, `async-dispatch`, `inbound-message-service`, `worker-identity-capture`, `completion-signal`; 2 already informally named: `context-budget-probe`, `roster-reconcile`). The core's logic references capabilities by name; each adapter binds them to concrete tools.
- DONE: Specified the Pi adapter bindings for all 7 capabilities (frictions 1–6) — `worker-back-channel` PRESENT via `contact_supervisor` + `intercom send/reply`; `async-dispatch` ASYNC via `subagent(async:true)` + `status`/`interrupt`; `inbound-message-service` via `intercom pending`; `worker-identity-capture` schema including intercom target + stamped model; `completion-signal` DUAL (subagent return OR inbound done-message, file-verify gate); `context-budget-probe` NONE; `roster-reconcile` NONE.
- DONE: Specified the ensign talkback protocol (friction 5) — `## Clarification` section for `pi-ensign-runtime.md` using `contact_supervisor` (`need_decision` for blocking, `progress_update` for non-blocking), mirroring Claude's `SendMessage(to="team-lead")` and Codex's thread-based clarification.
- DONE: Finalized 6 behavior-bound ACs (AC-1 structural contractlint for capability-name↔adapter-binding; AC-2 live pi-live back-channel round-trip; AC-3 completion-signal duality; AC-4 ensign talkback exercised by AC-2; AC-5 worker identity capture; AC-6 claude-live/codex-live regression). Each AC has a test method.
- DONE: Recorded "no spike needed" — the riskiest mechanism (async subagent + intercom back-channel) is already proven live by run `0637e2ed` (cited, not re-spiked). The archived spike `cq9kb7cdpp9y48tn8gwzmqzq` proved the host talkback chain.
- DONE: Confirmed the deferred boundary (frictions 7–9) is clean — each has a sharp seam in the named-capability declarations (`single-pending` in `worker-back-channel`, `standing-injection` no-op, `fresh-redispatch` default).
- SKIPPED: Documentation diff — the skill files ARE the docs; `docs/runtime-support.md` and the dev README do not reference the dispatch back-channel by name; no separate doc diff needed.

### Summary

Ideation finalized. The entity body now carries the full approach (Deliverable A: 7 named capabilities with before/after for the core rewrite; Deliverable B: Pi adapter bindings for frictions 1–6 including the ensign talkback protocol wording), 6 behavior-bound ACs, and a test plan with cost/complexity. The spike evidence from the prior run is preserved. No product files were edited (ideation = design only). Ready for the ideation gate.
