---
name: fo-gate-lifecycle
description: Load at any engaged gate entry before action.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, and apply one recorded gate authorization

Load this skill in one host event before gate probe, mutation, presentation, route, replay, or dispatch. It grants no write authority; read `fo-write-core.md` before FO mutation.

**Boot projection.** Use unresolved actionable `ready_gates` rows from `status --boot --identify --json`: `awaiting-captain` is an open current-stage Briefing; `approved-awaiting-merge` and `approved-awaiting-advance` are unblocked pending approvals. No selected attempt is omitted `validating` and must be prepared before presentation. Engage row `slug` and read the entity; never infer readiness from stage alone.

**Capability preflight.** Per lifecycle, resolve `${SPACEDOCK_BIN:-spacedock}` and run exactly one fresh `gate --help`. Require `prepare`, `record`, `validate`, `eligibility`, `consume`, the prepare flags `--question`, `--artifact`, `--summary`, `--reference`, `--workflow-dir`, and the semantic record flags `--briefing`, `--room`, `--decision`, `--actor`, `--reason`; reject retired `--directive` exposure. On absence or exposure, halt before mutation: do not commit selected sources, prepare a room, or change state; prescribe refresh or a fresh build via `SPACEDOCK_BIN`. Never hand-edit `gates:`.

**Prepare and bind.** Select one Markdown gate-review Artifact and References, author its exact concise summary, and commit newly authored selections in their owning main/state histories after preflight. Supply judgment and paths only; never author JSON, ids, digests, Git-root locators, or room coordinates. Relative paths resolve from launch cwd.

```text
${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary SUMMARY [--reference FILE ...] --workflow-dir WORKFLOW_DIR
```

Require the emitted `room`, `briefing`, `digest`, and `state=open` lines. The emitted clean absolute room is the only room authority; never reconstruct it. Preparation binds a two-file recorder-ready room with no copied sources. `«state.commit»(slug)` commits the folder entity or flat Markdown-plus-companion room unit before presentation. Invoke `«gate.ac-cross-check»`, judge evidence, then `«gate.assemble-verdict»`. On chat, presentation completes only after one root review names entity/stage, exact bound Briefing id/digest, recommendation, and decision ask. It follows the bind commit and must precede decision record; delegated conn does not waive it.

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
