---
name: fo-gate-lifecycle
description: Load at any engaged gate entry before action.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, and apply one recorded gate authorization

Load this skill in one host event before gate probe, mutation, presentation, route, replay, or dispatch. It grants no write authority; read `fo-write-core.md` before FO mutation.

**Boot projection.** Use only unresolved actionable `ready_gates` rows from `status --boot --identify --json`, fixed keys `id`, `slug`, `current`, `readiness`: `awaiting-captain` = selected current-stage open Briefing; `approved-awaiting-merge` = unblocked approve + advance/pending to terminal; `approved-awaiting-advance` = nonterminal target. Gate-stage/no selected attempt is omitted `validating`: retain and bind the selected Briefing before presentation, including when no decision authority is supplied. Malformed/stale selection, blocked/held, feedback, consumed/superseded/not-applicable are omitted. Opt-in human/JSON `gate-readiness` summarizes; `gate-*` retains optional diagnostics. Engage row `slug` and read the entity; never infer readiness from status/stage.

**Capability preflight.** Per lifecycle, resolve `${SPACEDOCK_BIN:-spacedock}` and run exactly one fresh `gate --help`. Require `record`, `validate`, `eligibility`, `consume`, `--briefing`, `--room`, `--decision`, `--actor`, `--reason`; reject retired `--directive` exposure. On absence or exposure, halt before mutation; prescribe refresh or a fresh build via `SPACEDOCK_BIN`. Never hand-edit `gates:`.

**Retain and bind.** Assemble `ROOM/briefing.json`: concise capability/change, tests/evidence, reviewed snapshot, material/deferred/polish findings, one recommendation, decision ask; link raw entity/spec/reports. Keep reproducible payloads as URI + SHA, else freeze a room copy. Relative retained-input paths resolve from launch cwd.

```text
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --briefing BRIEFING --workflow-dir WORKFLOW_DIR
```

Require exit 0, the expected gate/attempt/Briefing, and `state=open`; record already validates before atomic write. Present the bound Briefing identity and digest read from entity state, never a recomputed file hash or artifact `rev`. `«state.commit»(slug)` must commit the folder room and index before presentation. Then invoke `«gate.ac-cross-check»`, make the evidence judgment, and invoke `«gate.assemble-verdict»`. On chat, presentation completes only after `present-gate` emits one root review naming the entity/stage, exact bound Briefing id/digest, recommendation, and decision ask. It must follow the bind commit and precede decision record; delegated conn does not waive it.

**Record and durably close.** Use exactly one semantic source:

```text
# Captain personally rendered the chat decision
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor person:captain [--reason REASON] --workflow-dir WORKFLOW_DIR

# FO rendered it under an explicit Captain conn
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor agent:first-officer --reason EVIDENCE_JUDGMENT --workflow-dir WORKFLOW_DIR

# Exact retained provider Result from its prepared room
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --room ROOM --workflow-dir WORKFLOW_DIR
```

No explicit Captain grant in the active conversation: bind/present; leave the gate open. A grant including one issued later in that conversation permits delegation. Record an FO-rendered decision as `agent:first-officer` with a nonblank reason, never `person:captain`; reserve it for the Captain's own decision. Recorder authenticates/retains no grant. `revise`/`hold` need reasons; provider mappings must be complete.

Map Captain calls before recording: `approve` maps to `approve` with an accepts-direction evidence reason; `redo with feedback` maps to `revise` with an accepts-direction reason; `reject` with `feedback-to` maps to `revise` with a rejects-direction reason; `reject` without `feedback-to` maps to `hold` with a pause reason; `hold` maps to `hold` with a pause reason; `not yet` maps to `hold` with a pause reason naming what remains. Routed redo/reject reasons include concrete asks and invoke `«feedback.route»` after the close commit; hold decisions commit and stop at the gate.

Require exit 0, the bound attempt/Briefing, `state=closed`, and the decision; record already validates the Resolution/application before atomic write. After every successful close, `«state.commit»(slug)` must commit that exact Resolution before approve, revise, hold, or any consume attempt. Close/commit failure halts.

**Route fail-closed.** For approve, run:

```text
${SPACEDOCK_BIN:-spacedock} gate consume ENTITY --workflow-dir WORKFLOW_DIR
```

Consume itself rechecks currency, successor, blockers, and one-use state under lock. Require exit 0 with `approved-pending`, `eligible=true`, `consumed=true`, and the expected successor. It atomically writes successor status plus consumed state; commit that descendant through `«state.commit»(slug)`. A nonterminal target then enters ordinary reuse-or-fresh dispatch; a terminal target enters the existing merge guard/hook and has no successor dispatch. Never use `status --set` to advance a gate.

- `revise`: after its close commit, never consume; invoke `«feedback.route»`.
- `hold`: after its close commit, remain at the gate and surface the reason.
- blocked/wrong-stage/unknown/ineligible approval: its close is already durable; halt and preserve status bytes.
- `stale`: consume exits nonzero, leaves status unchanged, changes only pending → superseded; commit it, bind a replacement Briefing, re-present.
- already `consumed`: the authorization is spent. A nonterminal current status resumes ordinary dispatch/recovery; a terminal current status resumes the existing merge ceremony. Do not re-record, consume, or dispatch a terminal successor. A diagnostic repeat consume must be nonzero and byte-clean.

**Resume.** Use boot/entity state and prior result; `gate validate`/`gate eligibility` are optional diagnostics, never mandatory positive-path calls. Same-Briefing bind is idempotent; changed open package updates and re-presents. Before routing any closed state, ensure its exact Resolution commit exists. Pending approval resumes consume; revise/hold routes/stops; consumed resumes dispatch only when nonterminal and otherwise merge; stale consume materializes supersession then replacement. Surface nonzero command, exit, remedy; never repair frontmatter.
