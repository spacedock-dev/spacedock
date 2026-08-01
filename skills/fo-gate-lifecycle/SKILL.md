---
name: fo-gate-lifecycle
description: Handle an engaged gate entry.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, apply one recorded authorization

Load before engaged gate action. It grants no writes; read `fo-write-core.md` before FO mutation.

| Command | Responsibility |
|---------|----------------|
| `gate prepare` | Opens one attempt, freezes selected authority in its room, and binds its Briefing. |
| `gate withdraw` | Retires one stale open prepared attempt without a decision. |
| `gate record` | Closes that attempt from one semantic source; never advances status. |
| `gate consume` | Applies one eligible approval; terminal targets route unspent to merge guard. |

**Boot projection.** Use actionable `ready_gates` from `status --boot --identify --json`. Engage row `slug`; read its entity, never infer readiness from stage. `awaiting-captain` is open; `withdrawn-awaiting-prepare` needs its successor prepared; `approved-awaiting-merge`/`approved-awaiting-advance` are unblocked. Retained authority for resume/present must match current status; prior-stage authority is history. If no selected attempt (omitted `validating`) or it mismatches, prepare current stage and present only its emitted binding.

**Prepare and bind.** Resolve `${SPACEDOCK_BIN:-spacedock}`. Select a Markdown gate-review Artifact and References, author its concise summary, then commit selections. Supply judgment and paths; never author JSON, ids, digests, Git-root locators, or room coordinates. Paths use launch cwd.

```text
${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary SUMMARY [--reference FILE ...] --workflow-dir WORKFLOW_DIR
```

Preflight one lifecycle surface: `prepare`, `withdraw`, `record`, `validate`, `eligibility`, `consume`, and withdrawal's `--reason`. A nonzero command halts; surface its exact error and refresh or rebuild the version-gated bundle when unavailable. Never hand-edit `gates:` or replace binary-owned entity/room authority.

Require emitted `room`, `briefing`, `digest`, `state=open`; never reconstruct the absolute room. `«state.commit»(slug)` commits the entity and two-file room. Load `spacedock:present-gate`, cross-check ACs, assemble the verdict, and present entity/stage, bound Briefing, recommendation, and ask after commit.

**Withdraw stale open authority.** If a prepared room is stale before provider output or Captain decision, run:

```text
${SPACEDOCK_BIN:-spacedock} gate withdraw ENTITY --reason REASON --workflow-dir WORKFLOW_DIR
```

Require `state=withdrawn`, commit the entity, and stop unless replacement inputs are ready. Withdrawal never means approve, revise, or hold. On `withdrawn-awaiting-prepare`, prepare and commit N+1, present its emitted room, and stop open. Never record, consume, present, dispatch, or recover with `record --briefing` on withdrawn N.

**Record and durably close.** Use exactly one semantic source:

```text
# Captain's chat decision
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor person:captain [--reason REASON] --workflow-dir WORKFLOW_DIR

# FO decision under explicit Captain conn
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor agent:first-officer --reason EVIDENCE_JUDGMENT --workflow-dir WORKFLOW_DIR

# Recorder-ready prepared room
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --room ROOM --workflow-dir WORKFLOW_DIR
```

No explicit Captain grant in active conversation: bind/present; leave gate open. A grant issued later in that conversation permits delegation. Record an FO-rendered decision as `agent:first-officer` with a nonblank reason, never `person:captain`; reserve it for the Captain's own decision. Recorder authenticates/retains no grant. `revise`/`hold` need reasons; room-backed mappings must be complete.

Map Captain calls before recording: `approve` maps to `approve` with an accepts-direction evidence reason; `redo with feedback` maps to `revise` with an accepts-direction reason; `reject` with `feedback-to` maps to `revise` with a rejects-direction reason; `reject` without `feedback-to` maps to `hold` with a pause reason; `hold` maps to `hold` with a pause reason; `not yet` maps to `hold` with a pause reason naming what remains. Routed redo/reject reasons include concrete asks and invoke `«feedback.route»` after the close commit; hold decisions commit and stop at the gate.

Require exit 0, bound attempt/Briefing, `state=closed`, and decision; record validates Resolution/application before atomic write. After every successful close, `«state.commit»(slug)` must commit that exact Resolution before approve, revise, hold, or any consume attempt. Close/commit failure halts.

**Route fail-closed.** For approve, run:

```text
${SPACEDOCK_BIN:-spacedock} gate consume ENTITY --workflow-dir WORKFLOW_DIR
```

Consume rechecks currency, successor, blockers, and one-use state under lock. Nonterminal: require exit 0, `consumed=true`, expected successor; it atomically writes successor status plus consumed state — commit through `«state.commit»(slug)`, then ordinary dispatch. Terminal: require `consumed=false`, `route=approved-awaiting-merge` — consume writes nothing; drive `«merge.guard»(slug)`, no successor dispatch. Never use `status --set` to advance a gate.

- `revise`: never consume after close commit; invoke `«feedback.route»`.
- `hold`: after close commit, remain at the gate and surface the reason.
- blocked/wrong-stage/unknown/ineligible approval: close is durable; halt and preserve status bytes.
- `stale`: consume exits nonzero, leaves status unchanged, changes only pending → superseded; commit it, bind a replacement Briefing, re-present.
- pending terminal approval: unspent — drive `«merge.guard»`.
- already `consumed`: authorization is spent. Nonterminal status resumes ordinary dispatch/recovery; terminal status resumes the existing merge ceremony. Never re-record, consume, or dispatch a terminal successor. Diagnostic repeat consume must be nonzero and byte-clean.

**Resume.** Use boot/entity state and prior result; validate/eligibility are optional diagnostics. Exact open prepare replay is idempotent; divergent binding refuses. Require Resolution commit before routing closed state. Withdrawn open attempt → prepare its successor; pending approval → consume; revise/hold → route/stop; consumed → dispatch if nonterminal, else merge. Surface nonzero command, exit, and remedy; never repair frontmatter.
