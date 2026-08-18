---
name: fo-gate-lifecycle
description: Handle an engaged gate entry.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, apply one recorded authorization

Load before engaged gate action. It grants no writes; read `fo-write-core.md` before FO mutation.

The binary owns preparation, withdrawal, recording, and one-use consume; this skill only routes their observed results.

**Boot.** `status --boot --identify --json` `ready_gates` gives `definition_dir`, `entity_dir`, slug/stage/readiness. Engage slug via `status --read <slug> --json`, never search. `needs-preparation`: review report, or at an `initial: true` stage the committed seed itself; `awaiting-captain`: open; `withdrawn-awaiting-prepare`: successor; approved: unblocked; malformed/ambiguous: fail closed.

**Prepare.** Resolve `${SPACEDOCK_BIN:-spacedock}`; keep supplied paths. Else, per distinct applicable retained absolute R=`definition_dir`/`entity_dir`, resolve intended root-relative P; state R uses the engaged task directory, never `.`. Run: `git -C "<R>" ls-tree -r --name-only HEAD -- "<P>" | awk 'tolower($0)~/\.(md|markdown)$/{print}'`. `<R>/<P>`: shell-quoted values. Read/use once; summarize. No harness logs/history/status/help probes. Supply judgment/cwd paths; no binary JSON/ids/digests/Git locators/room coords.

Run this sequence once and in this order:

```text
1. ${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary SUMMARY [--reference FILE ...] --workflow-dir WORKFLOW_DIR
2. ${SPACEDOCK_BIN:-spacedock} state commit ENTITY --workflow-dir WORKFLOW_DIR
3. In one shell event, run `${SPACEDOCK_BIN:-spacedock} status --read ENTITY --checklist --json --workflow-dir WORKFLOW_DIR` and the same command with `--ac-scan` instead of `--checklist`.
4. Load spacedock:present-gate and present once from that structured evidence.
```

**Initial-stage gate.** At a stage the workflow marks `initial: true`, the committed seed IS the reviewed artifact: no prior stage existed to write a report. Prepare directly from the seed with `--artifact` naming the entity; never author, request, or wait for a stage report, and never emit `report-incomplete`. Skip step 3 — both structured reads exit nonzero with no report. Present outcome, scope, exclusions, and the ideation proof the next stage owes.

**Cold report candidate (non-initial stages).** Structurally review the path-resolved entity's latest exact-stage report/checklist and commit. An insufficient obligation, claim, Summary, or scope stops once with `report-incomplete: <concrete defect>` and no prepare, mutation, presentation, idle, or repeat-next. Otherwise invoke prepare once; nonzero/mismatch stops with its exact error and no retry or `gate record --briefing`.
Require prepare to emit `room`, `briefing`, `digest`, and `state=open`; never reconstruct the room. On nonzero commit, structured-read failure, or stage mismatch, stop before presentation or any later lifecycle effect. Use checklist text/ranges and AC citations to cross-check the gate; do not full-read/grep the entity or project boot after prepare or presentation. No conn: ask and stop open. Explicit conn: immediately record and consume; never final after presentation.

**Withdraw stale open authority.** If a prepared room is stale before the Captain decision, run:

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
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor agent:first-officer --reason EVIDENCE_JUDGMENT --conn-quote GRANT --conn-source SOURCE --workflow-dir WORKFLOW_DIR

```

No Captain grant: bind/present and leave open; a later grant permits delegation. An existing conn requires present once, immediate record, then the route below. Record an FO-rendered decision as `agent:first-officer` with a nonblank reason, never `person:captain`; reserve it for the Captain. Cite the grant: quote it verbatim and name where it was given; the recorder refuses an FO decision without both. `revise`/`hold` need reasons.

Map Captain calls before recording: `approve` maps to `approve` with an accepts-direction evidence reason; `redo with feedback` maps to `revise` with an accepts-direction reason; `reject` with `feedback-to` maps to `revise` with a rejects-direction reason; `reject` without `feedback-to` maps to `hold` with a pause reason; `hold` maps to `hold` with a pause reason; `not yet` maps to `hold` with a pause reason naming what remains. Routed redo/reject reasons include concrete asks and invoke `«feedback.route»` after the close commit; hold decisions commit and stop at the gate.

Require exit 0, bound attempt/Briefing, `state=closed`, and decision; record validates Resolution/application before atomic write, commits, and syncs itself (split-root); `sync=`/`phase=` discriminator and recovery in fo-dispatch-core.md. Close/sync failure halts.

**Route fail-closed.** For approve, run:

```text
${SPACEDOCK_BIN:-spacedock} gate consume ENTITY --workflow-dir WORKFLOW_DIR
```

Consume rechecks currency, successor, blockers, and one-use state under lock. Nonterminal: require exit 0, `consumed=true`, expected successor; it writes successor/consumed state, commits, and syncs. Then call `dispatch build --stamp` exactly once in the bound adapter shape (Codex named: no `--bare-mode`/`--team-name`); never probe another shape. Require status and nonempty `started` in its commit. Terminal: require `consumed=false`, `route=approved-awaiting-merge`; drive `«merge.guard»(slug)` with no dispatch. Never use `status --set` to advance a gate.

- `revise`: never consume after close commit; invoke `«feedback.route»`.
- `hold`: after close commit, remain at the gate and surface the reason.
- blocked/wrong-stage/unknown/ineligible approval: close is durable; halt and preserve status bytes.
- `stale`: consume exits nonzero, leaves status unchanged, changes only pending → superseded; commit it, bind a replacement Briefing, re-present.
- pending terminal approval: unspent — drive `«merge.guard»`.
- already `consumed`: authorization is spent. Nonterminal status resumes ordinary dispatch/recovery; terminal status resumes the existing merge ceremony. Never re-record, consume, or dispatch a terminal successor. Diagnostic repeat consume must be nonzero and byte-clean.

**Resume.** Recover durable-but-unsynced writes per fo-dispatch-core.md; never retry failures or repair frontmatter. Routes: open → always `state commit`, then structured reads and presentation; failure stops; `needs-preparation` → review; withdrawn → prepare; pending approval → consume; revise/hold → route/stop; consumed → dispatch/merge. Prepare replay is idempotent; divergent binding refuses.
