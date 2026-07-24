---
name: fo-gate-lifecycle
description: Load at any engaged gate entry before action.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, and apply one recorded gate authorization

Load this skill in one host event before gate probe, mutation, presentation, route, replay, or dispatch. It grants no write authority; read `fo-write-core.md` before FO mutation.

**Boot projection.** Use only unresolved actionable `ready_gates` rows from `status --boot --identify --json`, fixed keys `id`, `slug`, `current`, `readiness`: `awaiting-captain` = selected current-stage open Briefing; `approved-awaiting-merge` = unblocked approve + advance/pending to terminal; `approved-awaiting-advance` = nonterminal target. Gate-stage/no selected attempt is omitted `validating`; without a supplied retained Briefing or decision authority, present that legacy gate from existing stage evidence without mutation and stop. Malformed/stale selection, blocked/held, feedback, consumed/superseded/not-applicable are omitted. Opt-in human/JSON `gate-readiness` summarizes; `gate-*` retains optional diagnostics. Engage row `slug` and read the entity; never infer readiness from status/stage.

**Capability preflight.** Immediately before every gate lifecycle, freshly resolve `${SPACEDOCK_BIN:-spacedock}` and run exactly one `gate --help`; do not cache it. Require `record`, `validate`, `eligibility`, `consume`, `--briefing`, `--room`, `--decision`, `--actor`, `--directive`. If absent, halt before mutation; prescribe refresh or a fresh build selected with `SPACEDOCK_BIN`. Never hand-edit `gates:`.

**Retain and bind.** Assemble `ROOM/briefing.json`: concise capability/change, tests/evidence, reviewed snapshot, material/deferred/polish findings, one recommendation, decision ask; link raw entity/spec/reports. Keep reproducible payloads as URI + SHA, else freeze a room copy. Relative retained-input paths resolve from launch cwd.

```text
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --briefing BRIEFING --workflow-dir WORKFLOW_DIR
```

Require exit 0, the expected gate/attempt/Briefing, and `state=open`; record already validates before atomic write. Present the bound Briefing identity and digest read from entity state, never a recomputed file hash or artifact `rev`. `«state.commit»(slug)` must commit the folder room and index before presentation. Then invoke `«gate.ac-cross-check»`, make the evidence judgment, and invoke `«gate.assemble-verdict»`; show the concise review, not raw JSON/YAML.

**Record and durably close.** Use exactly one semantic source:

```text
# Captain personally rendered the chat decision
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor person:captain [--reason REASON] --workflow-dir WORKFLOW_DIR

# FO rendered it under an explicit Captain conn
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor agent:first-officer --reason EVIDENCE_JUDGMENT --directive EXACT_QUOTED_CAPTAIN_GRANT --workflow-dir WORKFLOW_DIR

# Exact retained provider Result from its prepared room
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --room ROOM --workflow-dir WORKFLOW_DIR
```

`revise` and `hold` require a reason (or the provider's included same-Briefing Annotation). Delegated FO approval always carries both its nonblank evidence reason and the exact quoted grant; never relabel it `person:captain`. A provider Result requires the prepared room's complete retained presentation mapping and authorized Resolution.

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
