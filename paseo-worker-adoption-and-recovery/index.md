---
id: x9mfvnkqtbwgm4tba7dy41tw
title: Paseo worker record, FO restart adoption, and dispatch recovery
status: ideation
source: captain (2026-11-21)
started: 2026-11-22T17:09:25Z
completed:
verdict:
score: 0.70
sprint: paseo-backend
worktree:
issue:
pr:
mod-block:
sprint-readiness: ready
gates:
    version: 1
    records:
        - id: gate:x9mfvnkqtbwgm4tba7dy41tw:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:x9mfvnkqtbwgm4tba7dy41tw-backlog-1
              briefing:
                id: briefing:x9mfvnkqtbwgm4tba7dy41tw:backlog:attempt-1:revision-1
                digest: sha256:d0790bd77f86aa073ddedea14685c10fa1a0e149a0c77a38c8314a151ab900c3
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:x9mfvnkqtbwgm4tba7dy41tw:backlog:1
                briefing: briefing:x9mfvnkqtbwgm4tba7dy41tw:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-22T17:09:15.302548Z"
                decision: approve
                reason: Captain approved backlog->ideation dispatch on 6769 with agy neon high
              application:
                target-stage: ideation
                state: consumed
---

Persist the Paseo worker handle in the task so a restarted or compacted first officer adopts a live worker instead of double-dispatching, and handle wait timeouts, lost create receipts, dead workers, and duplicate completion events with deterministic single-retry-then-hold recovery isolated per entity.

## Problem

With `paseo-dispatch-backend` (`936`) alone, a Paseo-hosted first officer (FO) that restarts, compacts, or crashes mid-stage has no durable record of the in-flight Paseo agent ID, workspace ID, or pre-dispatch state revision. On resume, a fresh FO risks spawning a second worker on the same stamped worktree/branch or mistaking a prior stage report for new completion evidence. Similarly, Paseo's `paseo wait --timeout` exits `0` with JSON `status=timeout` during normal live progress, a dropped `paseo run --no-wait` call can leave an unrecorded worker running on the daemon (lost create receipt), and an ended worker with a missing/incomplete report or a duplicate finish event needs exact, zero-ambiguity recovery rules so one failing task never wedges unrelated tasks.

## Proposed approach

Extend `skills/first-officer/references/paseo-fo-dispatch.md` (shipped by `936`) and `skills/first-officer/references/fo-write-core.md`, add a small pure-Go worker-record codec and deterministic reconciliation evaluator in `internal/dispatch/paseo_worker.go` bound to the skill contract via `internal/contractlint`, and prove restart adoption end-to-end against a real isolated `paseo` daemon in `TestLivePaseoAdoptionSmoke`. No new `spacedock` CLI subcommands or fake ACP stub servers are introduced.

### 1. Exact worker-record schema stored in the task (`## Paseo Worker`)

Store one authoritative FO-owned H2 section `## Paseo Worker` in the canonical task body (in `.spacedock-state/{slug}/index.md` or `.spacedock-state/{slug}.md` for split-root workflows; the active entity file otherwise). Using H2 (`## Paseo Worker`) guarantees `internal/status/section_read.go` and `internal/status/gate_extract.go` cleanly bound it independently of any `## Stage Report: <stage>` section (an H3 appended after a stage report would be swallowed into that report's span).

**Write ownership & commit rules:**
- Owned exclusively by the FO (never edited by ensigns). Updated **in place** via the section span returned by `spacedock status --read <entity> --json` on launch, adoption, verification, archival, supersede, hold, or stage transition.
- Committed path-scoped in the state checkout (`git -C {state_checkout} add -- {entity_path} && git -C {state_checkout} commit -m "dispatch: record paseo worker for {slug} ({stage})" -- {entity_path}`), retrying once after 2s on `index.lock` contention, then synced per the split-root contract.
- **Gate briefing freeze safety:** For gated stages, the FO MUST update `## Paseo Worker` to `state: archived` and commit it BEFORE running `spacedock gate prepare`, so the frozen briefing digest is never invalidated after prepare.

**Exact fenced YAML schema (20 required canonical scalar keys in fixed order + 1 optional `hold_reason`):**

```markdown
## Paseo Worker

```yaml
host: paseo
server_id: srv_j3VeUq-Ba722
workspace_id: wks_8aec77dc0eda8e94
agent_id: 0cb685d1-8e49-4795-9705-d0b7759a5045
stage: implementation
cycle: 1
attempt: 1
worker_label: paseo-qjyrql-implementation
epoch_label: paseo-qjyrql-implementation-c1-a1
worktree: .worktrees/spacedock-ensign-paseo-qjyrql
cwd: /private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql
pre_dispatch_state_rev: 764719a10a5c397629231fa6db16565d8c617aac
requested_provider: codex/gpt-6-luna
requested_mode: auto
requested_model: gpt-6-luna
returned_provider: codex
returned_mode: full-access
returned_model: gpt-6-luna
state: running
recorded_at: 2026-11-22T17:10:00Z
```
```

**Field contract:**
- `host`: literal `"paseo"`.
- `server_id`: daemon server ID from `paseo inspect <agent_id> --json` / `paseo status --json` / `paseo.parent-server-id` label; `"local"` when a local isolated daemon omits a server ID.
- `workspace_id`: Paseo workspace ID from `paseo workspace create` (e.g., `wks_...`).
- `agent_id`: Paseo agent ID from `paseo run --no-wait` or recovered via `paseo ls`.
- `stage`: entity `status:` at dispatch time.
- `cycle`: positive integer — `1` initially; `1 + count(prior feedback cycles for this stage in ### Feedback Cycles)` on feedback rework.
- `attempt`: `1` (initial dispatch for `(stage, cycle)`) or `2` (single allowed `-retry` replacement).
- `worker_label`: exact `spacedock.worker` label passed to `paseo run`:
  - `cycle == 1, attempt == 1`: `{worker_base}` (canonical `envelope.name` from `spacedock dispatch build --host paseo`, or `spacedock dispatch name` during cold reconciliation).
  - `cycle == 1, attempt == 2`: `{worker_base}-retry`.
  - `cycle >= 2, attempt == 1`: `{worker_base}-cycle{cycle}`.
  - `cycle >= 2, attempt == 2`: `{worker_base}-cycle{cycle}-retry`.
- `epoch_label`: exact `spacedock.epoch` label passed to `paseo run` — deterministically `{worker_base}-c{cycle}-a{attempt}` (lowercase `[a-z0-9-]`, unique per `(entity, stage, cycle, attempt)`, recomputable cold from durable entity state with zero session memory).
- `worktree`: relative worktree path from entity `frontmatter.worktree` (`""` for worktree-free stages).
- `cwd`: resolved physical working directory (`filepath.EvalSymlinks` / `pwd -P`) adopted by `paseo workspace create`.
- `pre_dispatch_state_rev`: `git -C {state_checkout} rev-parse HEAD` captured immediately AFTER `dispatch build --stamp` and BEFORE `paseo run --no-wait` (prevents pre-existing or prior-cycle stage reports from falsely satisfying verification).
- `requested_provider`, `requested_mode`, `requested_model`: exact provider, mode (`""` if daemon default), and model requested on `paseo run`.
- `returned_provider`, `returned_mode`, `returned_model`: effective values returned by `paseo run` / `paseo inspect` (explicitly capturing discrepancies such as requested `auto` returning `full-access`).
- `state`: closed 5-value enum:
  - `running` — worker is live (dispatched or adopted) and awaiting completion.
  - `verified` — durable Git commits and exact-stage report passed the completion gate; archival pending.
  - `archived` — worker archived (`paseo archive <agent_id>`) after verification; terminal clean state for this attempt.
  - `superseded` — worker failed/ended without valid completion and was archived before launching a replacement.
  - `held` — entity exhausted its single retry or hit an unresolvable ambiguity/permission block; only this entity is held.
- `recorded_at`: UTC RFC3339 timestamp of the last record write.
- `hold_reason` (optional): non-empty concrete one-line diagnostic present only when `state == held`.

**Retry ledger (`### Dispatch Retries`):**
Reuse the existing FO-owned `### Dispatch Retries` entity-body subsection (distinct from `### Feedback Cycles`, so a dispatch retry never increments the feedback counter):
```text
- Retry 1: {stage} (cycle {cycle}) — {spawn-error | ended-no-report | incomplete-report}; re-dispatched -retry
```
(Lines without `(cycle N)` parse as `cycle 1` for backward compatibility.)

### 2. Deterministic FO restart adoption & lost-create-receipt reconciliation

On boot, compaction resume, or `«dispatch.next-action»()` for an entity in active stage `S = entity.status`:

1. **Step 0 — Derive expected identity from durable entity state alone:**
   - Read frontmatter (`status`, `started`, `worktree`), `## Paseo Worker` (if any), `### Feedback Cycles`, and `### Dispatch Retries`.
   - Run read-only `spacedock dispatch name --workflow-dir {workflow_dir} --entity-path {entity_path} --stage {S}` to get `worker_base`.
   - Compute `cycle = 1 + count(feedback cycles for S)`, `retries_done = count(retry entries for (S, cycle))`, `attempt = 1 + retries_done`.
   - Compute `expected_worker`, `expected_epoch = {worker_base}-c{cycle}-a{attempt}`, and physical `expected_cwd`.
2. **Step 1 — Check terminal/held/completed durable state first:**
   - If `## Paseo Worker` matches `(stage == S, cycle == cycle)` and `state == held`: skip this entity (remain held) and continue unrelated entities.
   - If `## Paseo Worker` matches `(stage == S, cycle == cycle)` and `state == archived`: if `HasValidNewReport` is true, resume completion advance (`merge guard` / `gate prepare` / advance to next stage) with zero new spawn; otherwise proceed to Step 2.
3. **Step 2 — Primary lookup via recorded `agent_id` (record matches `(S, cycle)` and `state` is `running` or `verified`):**
   - If `state == verified`: run idempotent `paseo archive <agent_id>`, update `state: archived`, commit path-scoped, and resume advance/gate. Zero new spawn.
   - If `state == running`: run `paseo inspect <record.agent_id> --json` (falling back to `paseo ls --all --global --json` on not-found) and route through the Recovery Matrix below.
4. **Step 3 — Lost-create-receipt reconciliation (no matching `## Paseo Worker` for `(S, cycle, attempt)`, empty `agent_id`, or `paseo run --no-wait` failed/dropped without returning a receipt):**
   - Query labels in exact priority order BEFORE ever calling `paseo run --no-wait`:
     1. `paseo ls --all --global --label "spacedock.epoch={expected_epoch}" --json`.
     2. If empty, fallback: `paseo ls --all --global --label "spacedock.worker={expected_worker}" --json`, filtered to rows with `realpath(row.cwd) == expected_cwd` and (when `started` is set) `row.createdAt >= started`, preferring non-archived rows.
   - Branch on candidate count:
     - **Exactly 1 match `W`:** call `paseo inspect <W.id> --json`; backfill & commit `## Paseo Worker` (using the `dispatch: {slug} entering {S}` state commit SHA as `pre_dispatch_state_rev` when missing); adopt `W` and enter the Recovery Matrix with **0 retries spent** and **0 duplicate spawn**.
     - **> 1 active (non-archived/non-closed) matches:** ambiguous roster — set `state: held`, `hold_reason: "ambiguous-active-workers: <id1>,<id2>"`, commit path-scoped, surface to captain, do NOT spawn a 3rd worker, do NOT kill unowned agents, continue unrelated entities.
     - **0 matches:** if triggered by ordinary boot/engage, proceed to `SpawnFresh`; if triggered by a failed `paseo run --no-wait`, classify as confirmed `spawn-error` and enter Recovery Matrix rows 8–9.

### 3. Exact recovery matrix (zero day-1 ambiguities)

**Durable Git verification predicate `HasValidNewReport(entity, stage, pre_dispatch_state_rev)`** — all four must hold:
1. `git -C {state_checkout} rev-list {pre_dispatch_state_rev}..HEAD -- {entity_rel}` is non-empty and the committed entity carries a new `## Stage Report: {stage}` passing `spacedock status --read <entity> --stage {stage} --checklist --json` (every retained checklist item accounted for, non-empty evidence lines, non-empty `### Summary`, no `FAILED` items except feedback-review findings).
2. `git -C {state_checkout} status --porcelain -- {entity_rel}` is empty.
3. For worktree stages (`worktree != ""`), the code worktree is on the expected task branch with clean `git -C {worktree} status --porcelain`.
4. `record.stage == entity.status` and `record.epoch_label` matches the current `(stage, cycle, attempt)`.

**FO self-preservation invariant:** Never pass `$PASEO_AGENT_ID` to `paseo archive`, and never archive the FO while any child worker is `running`.

| Row | Trigger / Observed Daemon State | `HasValidNewReport` | `retries_done` | Exact FO Decision & Action | Retries Spent | Archival | Unrelated Entities |
|---|---|---|---|---|---|---|---|
| 1 | `paseo wait` exits `0` with `status=timeout`; `paseo inspect` shows `status=running` & `pendingPermissions=[]` | Any | Any | `ContinueWait`: keep `state: running`; re-issue `paseo wait <agent_id> --timeout <T> --json` (never nudge, kill, or retry). | 0 | None | Unaffected |
| 2 | `paseo inspect` shows `status=running` with non-empty `pendingPermissions` (or permission `requiresAttention`) | Any | Any | `HoldEntity` (permission-blocked): set `state: held`, `hold_reason: "pending-permission: <agent_id>"`; commit path-scoped; surface to captain. | 0 | None (live worker preserved) | Unaffected |
| 3 | Restarted FO finds `state=running` (recorded or adopted from Step 3) and `paseo inspect` shows `status=running` | Any | Any | `AdoptRunning`: persist/confirm `## Paseo Worker` (`state: running`) and enter `paseo wait <agent_id>` loop. Zero duplicate spawn. | 0 | None | Unaffected |
| 4 | Worker ended (`idle` / `closed` / `archived` / `error` / not-found) or finish event observed | `true` | Any | `VerifyAndArchive`: set `state: verified` -> run idempotent `paseo archive <agent_id>` -> set `state: archived` & commit path-scoped -> advance / `gate prepare` / `merge guard`. | 0 | Archive verified worker | Unaffected |
| 5 | Worker ended (`idle` / `closed` / `archived` / `error` / not-found) | `false` (no new report or uncommitted/incomplete report) | `0` | `ArchiveAndRetry`: archive superseded worker (`paseo archive <agent_id>` if present & not archived); set `state: superseded`; append Retry 1 line (`ended-no-report` or `incomplete-report`) to `### Dispatch Retries`; commit path-scoped; spawn ONE fresh `-retry` worker (`attempt: 2`, `epoch_label: ...-a2`, reusing existing worktree WITHOUT `--stamp`, appending WIP check note), record new `## Paseo Worker` (`state: running`). | +1 (1/1) | Archive superseded worker BEFORE retry spawn | Unaffected |
| 6 | Worker ended (`idle` / `closed` / `archived` / `error` / not-found) | `false` | `>= 1` | `HoldEntity` (retry-exhausted): archive failed attempt-2 worker (`paseo archive <agent_id>` if present & not archived); set `state: held`, `hold_reason: "retry-exhausted: {stage} cycle {cycle} ({cause})"`; commit path-scoped; surface to captain; stop attempting this entity. | 0 (cap) | Archive failed attempt-2 worker | **Unaffected** |
| 7 | `paseo run --no-wait` fails/drops without receipt; Step 3 finds **1** matching worker | Re-evaluated via rows 3–6 | Any | `AdoptFromLabels`: backfill `## Paseo Worker` from `paseo inspect`, commit path-scoped, then route through row 3 (if `running`) or rows 4–6 (if ended). | 0 | Per adopted row | Unaffected |
| 8 | `paseo run --no-wait` fails/drops; Step 3 finds **0** matching workers (`spawn-error`) | N/A | `0` | `ArchiveAndRetry` (`spawn-error`): append Retry 1 line (`spawn-error`) to `### Dispatch Retries`; commit path-scoped; spawn ONE fresh `-retry` worker (`attempt: 2`, `epoch_label: ...-a2`). | +1 (1/1) | None (nothing created) | Unaffected |
| 9 | `paseo run --no-wait` fails/drops; Step 3 finds **0** matching workers (`spawn-error`) | N/A | `>= 1` | `HoldEntity` (`spawn-error` exhausted): set `state: held`, `hold_reason: "retry-exhausted: {stage} cycle {cycle} (spawn-error)"`; commit path-scoped; surface to captain; stop attempting this entity. | 0 (cap) | None | **Unaffected** |
| 10 | Duplicate/late finish event, repeated `paseo wait`, or re-entrant check for `(agent_id, epoch_label)` already `archived`/`superseded`, or `entity.status != record.stage` | N/A | Any | `IgnoreDuplicate`: strict no-op (zero state mutation, zero retry, zero spawn; if `state == verified`, finish idempotent archival and advance). | 0 | No-op | Unaffected |
| 11 | Step 3 finds **> 1** active matching workers for `(entity, stage, cycle)` | N/A | Any | `HoldEntity` (ambiguous-roster): set `state: held`, `hold_reason: "ambiguous-active-workers: <ids>"`; commit path-scoped; surface to captain; do NOT spawn a 3rd worker or kill unowned agents. | 0 | None (captain decides) | **Unaffected** |

### 4. `TestLivePaseoAdoptionSmoke` specification (`internal/ensigncycle/paseo_adoption_live_test.go`, `//go:build live`)

Reuses `936`'s isolated `paseo` daemon helper (`SPACEDOCK_PASEO_BIN` defaulting to `paseo` on `PATH`, `filepath.EvalSymlinks(t.TempDir())` home, free port, `relay`/`dictation`/`voiceMode` off, env scrubbed, auth seeded from `CODEX_AUTH_JSON` or local auth, skipped locally when `paseo` or auth is absent unless `SPACEDOCK_PASEO_LIVE_REQUIRED=1`, stopped in `t.Cleanup`; provider `SPACEDOCK_PASEO_LIVE_PROVIDER` defaulting to `codex/gpt-6-luna`, overridable to e.g. `antigravity/gemini-4-neon`).

**Deterministic mid-stage stop choreography (zero sleep races):**
1. **Fixture:** Split-root temp workflow (`README.md` with `commissioned-by: spacedock@test`, `state: .spacedock-state`, stages `backlog -> implementation (worktree: true) -> done (terminal: true)`, no `_mods/`). The `implementation` stage definition instructs the worker to: (a) write and commit `receipt.txt` on the task branch, (b) poll up to 180s (1s interval) until `<fixtureRoot>/.adoption-release` exists, then (c) append `## Stage Report: implementation` and commit it path-scoped in `.spacedock-state`.
2. **Phase 1 (FO-1 dispatch & mid-stage stop):** Launch FO-1 on the isolated daemon with the entity ready in `implementation`. Poll `paseo ls --all --global --label "spacedock.worker=<worker_base>" --json` and the entity file until: (a) exactly 1 worker `W1` exists with `status == "running"`, and (b) `## Paseo Worker` is committed in `.spacedock-state` with `agent_id == W1` and `state: running`. Immediately stop/archive FO-1 while `W1` remains `running` (blocked on `.adoption-release`). Verify FO-1 is closed and `paseo inspect <W1> --json` still reports `status == "running"`.
3. **Phase 2 (Fresh FO-2 boot & adoption):** Launch a brand-new FO-2 (fresh session, zero shared state with FO-1) on the same workflow and isolated daemon. Once FO-2 is active (and has begun adopting `W1` while `W1` is `running`), create `<fixtureRoot>/.adoption-release`. Wait for FO-2 to finish its run.
4. **Phase 3 (Durable assertions):**
   - `paseo ls --all --global --label "spacedock.worker=<worker_base>" --json` (and all non-FO daemon agents for this workflow) contains **exactly 1 worker** (`W1`, zero `-retry` or duplicate workers spawned).
   - `paseo inspect <W1> --json` / `paseo ls --all` confirms `W1` is archived/closed.
   - Locate the task entity (in `.spacedock-state/_archive/` or `.spacedock-state/`): `dispatch.ParsePaseoWorkerRecord` succeeds with `agent_id == W1`, `stage == "implementation"`, `cycle == 1`, `attempt == 1`, `worker_label == worker_base`, `epoch_label == worker_base+"-c1-a1"`, non-empty `server_id`, `workspace_id`, `worktree`, `cwd`, `pre_dispatch_state_rev`, requested & returned provider/mode/model, and `state == "archived"`.
   - `dispatch.ParseDispatchRetries` returns `0` (neither FO-1 stop nor wait timeout spent a retry).
   - `## Stage Report: implementation` is committed path-scoped strictly after `pre_dispatch_state_rev`; `receipt.txt` is committed on the task branch; frontmatter `status` durably advanced to `done` (`verdict: PASSED`); checkouts are clean; `t.Cleanup` stops the daemon cleanly.

### 5. Mechanism-to-value justification

| Mechanism | Value AC served | Simplest alternative considered | Why alternative is insufficient |
|---|---|---|---|
| H2 `## Paseo Worker` entity-body YAML block | AC-1, AC-2, AC-4 | Flat YAML frontmatter fields via `status --set`, or shared `.paseo-workers.json` sidecar | `status --set` only writes flat scalars (cluttering frontmatter with 20 transient fields), while a shared sidecar breaks per-entity path-scoped commits and causes rebase conflicts across sibling tasks. H2 cleanly terminates `## Stage Report:` in `section_read.go` and `gate_extract.go` with zero `status` binary changes. |
| Deterministic `spacedock.epoch` (`{worker_base}-c{cycle}-a{attempt}`) + dual `paseo ls` lookup | AC-1, AC-2 | Random UUID epoch persisted only after `paseo run` returns | A crash between `paseo run --no-wait` and writing `## Paseo Worker` (lost create receipt) leaves a random epoch unrecoverable, risking a duplicate spawn on the same worktree. Deterministic epoch + `spacedock.worker`/`cwd` fallback recovers cold with zero pre-launch commit churn. |
| Pure-Go `internal/dispatch/paseo_worker.go` + `internal/contractlint` dual extraction | AC-1, AC-2, AC-3 | Prose-only rules in `paseo-fo-dispatch.md`, or a new `spacedock dispatch paseo-reconcile` CLI wrapper | Prose-grep tests are banned by the Proof Policy and live LLM runs for all 11 failure/duplicate rows are slow and flaky; conversely, shelling out to `paseo` inside the `spacedock` binary couples releases to fast-moving Paseo CLI flags. The pure-Go codec/evaluator grades `TestLivePaseoAdoptionSmoke`, drives offline table tests, and binds `paseo-fo-dispatch.md` structurally without adding CLI surface. |
| File-sentinel `.adoption-release` in `TestLivePaseoAdoptionSmoke` | AC-1 | Fixed `time.Sleep` before killing FO-1 | Model latency varies 10–90s; a fixed sleep races either before worker launch or after worker exit. The sentinel deterministically holds `W1` in `running` across the FO-1 stop and FO-2 boot. |

### 6. Concrete skill and documentation diffs

#### A. `skills/first-officer/references/fo-write-core.md` (before / after)

- **Mutation Gate `allowed-state` rule cell:**
  - *Before:* `Entity frontmatter via ${SPACEDOCK_BIN:-spacedock} status --set, spacedock new, archive moves, state-transition commits, and ### Feedback Cycles under the existing state/worktree rules. Roadmap and sprint docs under _roadmap/ are edited directly and committed path-scoped in the state checkout.`
  - *After:* `Entity frontmatter via ${SPACEDOCK_BIN:-spacedock} status --set, spacedock new, archive moves, state-transition commits, and FO-owned operational body sections (### Feedback Cycles, ### Dispatch Retries, ## Paseo Worker) under the existing state/worktree rules. Roadmap and sprint docs under _roadmap/ are edited directly and committed path-scoped in the state checkout.`
- **`## FO Write Scope` bullet 3 & off-limits paragraph:**
  - *Before:* ``- **`### Feedback Cycles` section** — in entity bodies, tracking rejection rounds. ...`` and ``and entity body content beyond `### Feedback Cycles` (stage reports, design, implementation notes belong to dispatched workers).``
  - *After:* ``- **FO-owned operational body sections (`### Feedback Cycles`, `### Dispatch Retries`, `## Paseo Worker`)** — tracking rejection rounds, single-retry dispatch ledgers, and Paseo worker handles. Update in place and commit path-scoped in the state checkout (or worktree/main when non-split-root).`` and ``and entity body content beyond `### Feedback Cycles`, `### Dispatch Retries`, and `## Paseo Worker` (stage reports, design, implementation notes belong to dispatched workers).``

#### B. `skills/first-officer/references/paseo-fo-dispatch.md` (new sections appended after `936`'s 5-verb binding)

Add three delimited sections carrying machine-readable markers for `internal/contractlint` dual extraction:
1. `## Paseo Worker Record` bounded by `<!-- PASEO-WORKER-SCHEMA:START -->` / `<!-- PASEO-WORKER-SCHEMA:END -->` listing the exact 20 required scalar keys, optional `hold_reason`, closed 5-value `state` enum (`running | verified | archived | superseded | held`), deterministic label rules (`spacedock.worker`, `spacedock.epoch`), path-scoped commit rule, and pre-`gate prepare` archival order.
2. `## Restart Adoption and Lost-Create-Receipt Reconciliation` carrying exact Steps 0–3 above (`spacedock dispatch name`, `paseo inspect`, `paseo ls --all --global --label spacedock.epoch=...`, fallback `spacedock.worker=...` + physical `cwd`).
3. `## Recovery Matrix and Retry Ledger` bounded by `<!-- PASEO-RECOVERY-MATRIX:START -->` / `<!-- PASEO-RECOVERY-MATRIX:END -->` carrying the exact 11-row table, `HasValidNewReport` predicate, `### Dispatch Retries` line format, per-entity isolation rule, and `$PASEO_AGENT_ID` self-preservation guard.

#### C. User-visible documentation (`docs/runtime-live-ci.md`; `docs/site/**` unchanged)

No `docs/site/**` page documents the opt-in Paseo dispatch backend yet, so `docs/site/**` needs no change. Update `docs/runtime-live-ci.md` under the Paseo live section added by `936` (STE compliant):

- *Before (from `936`):*
  ```markdown
  Run `TestLivePaseoDispatchSmoke` locally against an isolated Paseo daemon:

  ```bash
  go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle -v
  ```
  ```
- *After:*
  ```markdown
  Run the Paseo live smoke tests locally against an isolated Paseo daemon:

  ```bash
  go test -tags live -count=1 -timeout 15m -run '^TestLivePaseo(Dispatch|Adoption)Smoke$' ./internal/ensigncycle -v
  ```

  `TestLivePaseoAdoptionSmoke` stops the first officer mid-stage and verifies that a fresh first officer adopts the single running worker, advances state, and archives the worker.
  ```

## Risk evidence

no spike needed: all underlying Paseo primitives and Spacedock parsers are already proven on this codebase:
- `_roadmap/paseo-backend/probe-report.md` and `run-log.md` verified `paseo ls --all --global --label 'spacedock.epoch=...' --json`, `paseo inspect <agent_id> --json`, `paseo wait <agent_id> --timeout 1 --json` returning exit `0` with `status=timeout` while live, `paseo archive` / `kill_agent` closing workers, and requested `modeId=auto` returning `currentModeId=full-access` (which the schema records explicitly).
- `internal/status/section_read.go` (`assignSectionLines`) and `internal/status/gate_extract.go` (line 144) already terminate `## Stage Report:` at any `## ` heading, proving H2 `## Paseo Worker` cannot contaminate stage reports or `--ac-scan`.
- `internal/contractlint/fo_write_core_mutation_gate_test.go` parses only columns 1–2 (`class`, `patterns`) of `FO-WRITE-CLASSIFIER`, proving the `allowed-state` rule-text update preserves all classifier guards.
- **Sequencing dependency:** `936` (`paseo-dispatch-backend`) must land `paseo-fo-dispatch.md` and the isolated-daemon live test helper in `internal/ensigncycle/` before this task's PR merges.

## Out of scope

- Remote Paseo workers (`serverId` across machines) and cross-host adoption.
- Worker reuse across stages (`paseo send` steering or multi-stage worker continuity).
- Fake ACP stub servers or Paseo version-pinning machinery.
- Multi-FO concurrency leases (tracked separately by `79` / `entity-session-claim-lease`).
- New `spacedock` CLI subcommands wrapping `paseo`.

## Expected surface and tolerance

Estimate net LOC change: +493, across 7 files.

Breakdown: **+503 insertions / -10 deletions (net +493 LOC)** across 7 files. Tolerance: **+350 to +620 net LOC, across 5 to 8 files**.

| File | Est. insertions | Est. deletions | Net LOC | Purpose |
|---|---|---|---|---|
| `skills/first-officer/references/paseo-fo-dispatch.md` | +85 | -5 | +80 | Worker-record schema, labels, Steps 0–3 adoption, 11-row recovery matrix, retry ledger |
| `skills/first-officer/references/fo-write-core.md` | +4 | -4 | 0 | Authorize `## Paseo Worker` and `### Dispatch Retries` as FO-owned operational body sections |
| `internal/dispatch/paseo_worker.go` | +110 | 0 | +110 | `PaseoWorkerRecord` codec, `PaseoLabels`, `ParseDispatchRetries`, `ReconcilePaseoWorker` |
| `internal/dispatch/paseo_worker_test.go` | +120 | 0 | +120 | Deterministic table tests for all 11 matrix rows, parser round-trips, `status` compatibility |
| `internal/contractlint/paseo_adoption_binding_test.go` | +65 | 0 | +65 | Dual-extraction binding between `paseo-fo-dispatch.md` and `internal/dispatch/paseo_worker.go` |
| `internal/ensigncycle/paseo_adoption_live_test.go` | +115 | 0 | +115 | `TestLivePaseoAdoptionSmoke` (`//go:build live`) |
| `docs/runtime-live-ci.md` | +4 | -1 | +3 | Document local `TestLivePaseoAdoptionSmoke` command (STE) |

**Declared observable semantic changes:**
- *Command grammar:* None (`spacedock` CLI flags and subcommands unchanged).
- *Stored formats:* Adds optional FO-owned H2 entity-body section `## Paseo Worker` (20 canonical YAML scalar keys + optional `hold_reason`, closed 5-value `state` enum) and extends `### Dispatch Retries` lines with explicit `(cycle N)` and Paseo causes (`spawn-error | ended-no-report | incomplete-report`). YAML frontmatter is unchanged.
- *Authority:* `fo-write-core.md` explicitly permits FO writes to `## Paseo Worker` and `### Dispatch Retries` in entity bodies; ensigns remain barred from editing them.
- *Runtime behavior:* Paseo-hosted FOs persist worker records at dispatch, adopt live or lost-receipt workers across restarts via `## Paseo Worker` + `paseo ls/inspect` labels, treat `paseo wait` timeouts as live progress (0 retries), enforce single-retry-then-hold isolated per entity, ignore duplicate/stale completion events, and archive verified/superseded workers without archiving `$PASEO_AGENT_ID` while children run.

## Acceptance criteria

**AC-1 (VALUE) - A fresh restarted Paseo FO adopts an in-flight worker mid-stage with zero duplicate dispatch, advances durable workflow state, and archives the worker.**
Verified by: `go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoAdoptionSmoke$' ./internal/ensigncycle -v` passes against an isolated `paseo` daemon (`codex/gpt-6-luna` default; `SPACEDOCK_PASEO_BIN` and `SPACEDOCK_PASEO_LIVE_PROVIDER` overridable), stopping FO-1 while the worker is blocked on `.adoption-release`, booting fresh FO-2, releasing the worker, and asserting: `paseo ls --all --global --json` contains exactly 1 worker for the task, `## Paseo Worker` parses with `state: archived` and matching `agent_id` plus requested/returned provider/mode/model, `### Dispatch Retries` has 0 entries, the task-branch commit and path-scoped split-root stage report commit after `pre_dispatch_state_rev` exist, frontmatter `status` advanced durably to `done`, and the worker is archived. Spawning a second worker, dropping `## Paseo Worker` persistence, spending a retry on FO restart/timeout, failing to advance, or leaving the worker unarchived fails the test.

**AC-2 - Worker-record parsing, deterministic label derivation, lost-create-receipt reconciliation, and the 11-row recovery matrix enforce zero-retry timeouts, single-retry-then-hold, duplicate-event idempotency, and per-entity isolation.**
Verified by: `go test -race -count=1 ./internal/dispatch -run 'TestPaseoWorker|TestReconcilePaseoWorker'` passes across table cases covering every matrix row: `status=timeout` while `running` returns `ContinueWait` with 0 retries; lost create receipt with 1 matching `spacedock.epoch` / `spacedock.worker` candidate adopts with 0 retries and 0 duplicate spawn; ended worker with valid post-`pre_dispatch_state_rev` report verifies and archives; ended worker with missing/incomplete report archives the superseded worker and retries once (`-retry`, `a2`), then holds only the failing entity on attempt 2 while an independent sibling entity returns `SpawnFresh`; duplicate/stale finish events on `archived`/`superseded` records return `IgnoreDuplicate`; `>1` active workers or `pendingPermissions` returns `HoldEntity` with 0 retries and no spawn; malformed `## Paseo Worker` YAML or missing required keys fails closed. Spending a retry on timeout, allowing a 3rd attempt, holding a sibling entity, or accepting a pre-dispatch report fails the test.

**AC-3 - `paseo-fo-dispatch.md`, `fo-write-core.md`, and `internal/dispatch/paseo_worker.go` remain structurally bound with zero schema, state-enum, label-key, or write-gate drift.**
Verified by: `go test -race -count=1 ./internal/contractlint` passes, including `TestPaseoAdoptionContractBinding` (dual-extracting the 20 schema keys, 5-value `state` enum, `spacedock.worker`/`spacedock.epoch` label keys, and recovery action set from `paseo-fo-dispatch.md` and asserting set equality against `internal/dispatch/paseo_worker.go`) and `TestFOWriteCoreMutationGateClassifiesTargets`. Renaming, adding, or dropping any schema field, state token, label key, or recovery action on either side without updating the other fails the test.

**AC-4 - Spacedock validation and section/AC scanning remain clean on entities carrying `## Paseo Worker` and `### Dispatch Retries` alongside `## Stage Report: <stage>`.**
Verified by: `go test -race -count=1 ./internal/dispatch -run '^TestPaseoWorkerStatusCompatibility$'` drives `status.NativeRunner` (`--validate`, `--read --stage implementation --checklist --json`, and `--read --stage implementation --ac-scan --json`) over split-root entity fixtures placing `## Paseo Worker` both before and after an earlier stage report, asserting exit `0`, zero warnings, intact AC enumeration, and zero checklist/summary contamination. Demoting `## Paseo Worker` to an H3 after a stage report or emitting invalid frontmatter fails the test.

## Test plan

No existing primary proof owner covers FO restart adoption or Paseo recovery (`TestLivePaseoDispatchSmoke` in `936` tests single-FO fresh dispatch only and cannot detect duplicate spawns across restart or failure-matrix bugs). Each check below owns one distinct failure class:

1. **Deterministic codec, recovery-matrix, and `status` compatibility suite (primary owner for AC-2 & AC-4):**
   - Command: `go test -race -count=1 ./internal/dispatch -run 'TestPaseoWorker|TestReconcilePaseoWorker'` (plus full `go test ./...` and `go test ./... -race`).
   - Falsifying edits: spending a retry on `status=timeout`; allowing attempt 3 after Retry 1; propagating `HoldEntity` to an unrelated sibling entity; accepting a report at or before `pre_dispatch_state_rev`; omitting a required schema key; placing `## Paseo Worker` as H3 inside a stage report.
   - Cost/mode: deterministic, offline, `< 1s`, zero model spend.
2. **Contractlint dual-extraction structural binding (primary owner for AC-3):**
   - Command: `go test -race -count=1 ./internal/contractlint`.
   - Falsifying edits: adding/renaming a `## Paseo Worker` field, `state` enum value, `spacedock.*` label key, or recovery decision token in `paseo-fo-dispatch.md` without updating `internal/dispatch/paseo_worker.go` (or vice versa), or breaking `FO-WRITE-CLASSIFIER` in `fo-write-core.md`.
   - Cost/mode: deterministic, offline, `< 2s`, zero model spend.
3. **Real isolated-daemon live adoption smoke (primary owner for AC-1):**
   - Command: `go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoAdoptionSmoke$' ./internal/ensigncycle -v` (default `codex/gpt-6-luna`; optional local override `SPACEDOCK_PASEO_LIVE_PROVIDER=antigravity/gemini-4-neon`).
   - Falsifying edits: FO-2 spawning a second worker instead of adopting `W1`; FO-1 failing to commit `## Paseo Worker`; FO-2 burning a retry on restart/timeout; FO-2 failing to advance frontmatter or archive `W1`.
   - Cost/mode: live (real `paseo` daemon + real FO-1/FO-2/worker sessions), ~3–6 minutes, one targeted run per validation; skips cleanly locally when `paseo` or auth is absent unless `SPACEDOCK_PASEO_LIVE_REQUIRED=1`.
4. **No one-off manual validation required** beyond reviewing the targeted local `TestLivePaseoAdoptionSmoke` artifacts (daemon logs, `paseo ls --all --json` snapshot, final entity diff) during validation.

### Feedback Cycles

## Stage Report: ideation

- DONE: Define exact worker-record schema stored in the task and deterministic FO restart adoption + lost-create-receipt reconciliation rules using paseo ls/inspect labels (spacedock.worker, spacedock.epoch)
  Specified H2 `## Paseo Worker` fenced YAML schema (20 required keys + optional `hold_reason`, closed 5-value `state` enum), deterministic `spacedock.worker`/`spacedock.epoch` (`{worker_base}-c{cycle}-a{attempt}`) derivation aligned with `936`, pre-`gate prepare` archival ordering, and Steps 0–3 adoption + lost-create-receipt reconciliation in `docs/dev/.spacedock-state/paseo-worker-adoption-and-recovery/index.md`.
- DONE: Define exact timeout vs dead-worker vs duplicate-event recovery matrix (single-retry-then-hold, per-entity isolation, verified/superseded archival) with zero day-1 ambiguities
  Defined the complete 11-row recovery matrix, `HasValidNewReport` 4-part Git verification predicate, `### Dispatch Retries` cycle-scoped ledger format, `fo-write-core.md` before/after authorization wording, `$PASEO_AGENT_ID` self-preservation guard, and per-entity hold isolation.
- DONE: Specify TestLivePaseoAdoptionSmoke (stop FO mid-stage, boot fresh FO, assert single worker + durable advance + archive) and update Expected surface, scanner-exact Acceptance criteria with falsifying Verified by clauses, and Test plan
  Specified `TestLivePaseoAdoptionSmoke` with `.adoption-release` sentinel choreography and durable assertions, updated Expected surface (+493 net LOC / +503 -10 across 7 files + declared semantic changes), wrote scanner-exact `## Acceptance criteria` (AC-1, AC-2, AC-3, AC-4 with falsifying `Verified by:` clauses), and added the Test plan and mechanism-to-value table; verified via `spacedock status --read --ac-scan`, `--checklist`, and `--validate`.

### Summary

Fleshed out the full ideation spec for `paseo-worker-adoption-and-recovery` aligned with `936` (`paseo-dispatch-backend`), defining an H2 `## Paseo Worker` YAML record schema, deterministic `spacedock.worker`/`spacedock.epoch` label reconciliation, an unambiguous 11-row recovery matrix with single-retry-then-hold and per-entity isolation, and `TestLivePaseoAdoptionSmoke` using a race-free `.adoption-release` file sentinel. Kept the design minimal with zero new CLI subcommands or fake ACP mocks, pairing a small pure-Go codec/reconciler in `internal/dispatch/paseo_worker.go` with dual-extraction `contractlint` checks and concrete `fo-write-core.md` / `docs/runtime-live-ci.md` diffs.
