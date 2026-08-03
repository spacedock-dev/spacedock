---
name: fo-gate-lifecycle
description: Handle an engaged gate entry.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, apply one recorded authorization

Load before engaged gate action. It grants no writes; read `fo-write-core.md` before FO mutation.

The binary owns preparation, withdrawal, recording, and one-use consume; this skill only routes their observed results.

**Boot projection.** Use `ready_gates` from `status --boot --identify --json`. Engage `slug`; read its entity and record for current status. `needs-preparation` is mechanical: review the report before writing. `awaiting-captain` is open; `withdrawn-awaiting-prepare` needs its successor; approved routes are unblocked. Prior authority is history; malformed/ambiguous fails closed.

**Prepare and bind.** Resolve `${SPACEDOCK_BIN:-spacedock}`. Select a Markdown gate-review Artifact and References, author its concise summary, then commit selections. Supply judgment and paths; never author JSON, ids, digests, Git-root locators, or room coordinates. Paths use launch cwd.

```text
${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary SUMMARY [--reference FILE ...] --workflow-dir WORKFLOW_DIR
```

Preflight one lifecycle surface: `prepare`, `withdraw`, `record`, `validate`, `eligibility`, `consume`, and withdrawal's `--reason`. A nonzero command halts; surface its exact error and refresh or rebuild the version-gated bundle when unavailable. Never hand-edit `gates:` or replace binary-owned entity/room authority.

Require emitted `room`, `briefing`, `digest`, `state=open`; never reconstruct the absolute room. `«state.commit»(slug)` commits the entity and two-file room. Load `spacedock:present-gate`, cross-check ACs, assemble the verdict, and present entity/stage, bound Briefing, recommendation, and ask after commit.

**Cold report candidate.** For one `needs-preparation` row, load this lifecycle and re-read the entity, latest exact-stage report/checklist, and report-bearing commit. It is structural evidence only. If any obligation, evidence, Summary, or scope claim is insufficient, stop once with `report-incomplete: <concrete defect>`; do not prepare, record, commit, present, idle, or repeat `status --next`. Otherwise choose the question, one committed Markdown Artifact, summary, and References, and invoke `gate prepare` exactly once. Require `room`, `briefing`, `digest`, `state=open`; commit, re-read the envelope, and present only one same-slug `awaiting-captain`. On nonzero/mismatch, stop; never retry or use `gate record --briefing`.

**Withdraw stale open authority.** If a prepared room is stale before provider output or Captain decision, run:

```text
${SPACEDOCK_BIN:-spacedock} gate withdraw ENTITY --reason REASON --workflow-dir WORKFLOW_DIR
```

Require `state=withdrawn`, commit the entity, and stop unless replacement inputs are ready. Withdrawal never means approve, revise, or hold. On `withdrawn-awaiting-prepare`, prepare and commit N+1, present its emitted room, and stop open. Never record, consume, present, dispatch, or recover with `record --briefing` on withdrawn N.

**Record and durably close.** Use exactly one semantic source:

```text
# Captain-approve fast path: close, sync, consume, sync in one call
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve --actor person:captain [--reason REASON] --consume --workflow-dir WORKFLOW_DIR

# Captain's chat decision
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor person:captain [--reason REASON] --workflow-dir WORKFLOW_DIR

# FO decision under explicit Captain conn
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor agent:first-officer --reason EVIDENCE_JUDGMENT --workflow-dir WORKFLOW_DIR

# Recorder-ready prepared room
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --room ROOM --workflow-dir WORKFLOW_DIR
```

No explicit Captain grant in active conversation: bind/present; leave gate open. A grant issued later in that conversation permits delegation. Record an FO-rendered decision as `agent:first-officer` with a nonblank reason, never `person:captain`; reserve it for the Captain's own decision. Recorder authenticates/retains no grant. `revise`/`hold` need reasons; room-backed mappings must be complete.

Map Captain calls before recording: `approve` maps to `approve` with an accepts-direction evidence reason; `redo with feedback` maps to `revise` with an accepts-direction reason; `reject` with `feedback-to` maps to `revise` with a rejects-direction reason; `reject` without `feedback-to` maps to `hold` with a pause reason; `hold` maps to `hold` with a pause reason; `not yet` maps to `hold` with a pause reason naming what remains. Routed redo/reject reasons include concrete asks and invoke `«feedback.route»` after the close commit; hold decisions commit and stop at the gate.

Require exit 0, bound attempt/Briefing, `state=closed`, and decision; record validates Resolution/application before atomic write, commits, and syncs itself (split-root); `sync=`/`phase=` discriminator and recovery in fo-dispatch-core.md. Close/sync failure halts.

**Route fail-closed.** For approve, run:

```text
${SPACEDOCK_BIN:-spacedock} gate consume ENTITY --workflow-dir WORKFLOW_DIR
```

Consume rechecks currency, successor, blockers, and one-use state under lock. Nonterminal: require exit 0, `consumed=true`, expected successor; it writes successor status plus consumed state, commits, and syncs itself, then dispatch via `dispatch build --stamp` (fo-dispatch-core.md). Terminal: require `consumed=false`, `route=approved-awaiting-merge` — consume writes nothing; drive `«merge.guard»(slug)`, no successor dispatch. Never use `status --set` to advance a gate.

- `revise`: never consume after close commit; invoke `«feedback.route»`.
- `hold`: after close commit, remain at the gate and surface the reason.
- blocked/wrong-stage/unknown/ineligible approval: close is durable; halt and preserve status bytes.
- `stale`: consume exits nonzero, leaves status unchanged, changes only pending → superseded; commit it, bind a replacement Briefing, re-present.
- pending terminal approval: unspent — drive `«merge.guard»`.
- already `consumed`: authorization is spent. Nonterminal status resumes ordinary dispatch/recovery; terminal status resumes the existing merge ceremony. Never re-record, consume, or dispatch a terminal successor. Diagnostic repeat consume must be nonzero and byte-clean.

**Resume.** Recover durable-but-unsynced writes per fo-dispatch-core.md; never retry a failed verb blindly or repair frontmatter. Use current boot/entity state: open → present, `needs-preparation` → semantic review above, withdrawn → prepare successor, pending approval → consume, revise/hold → route/stop, consumed → dispatch or merge. Exact open prepare replay is idempotent; divergent binding refuses.
