---
id: x9mfvnkqtbwgm4tba7dy41tw
title: Paseo worker record, FO restart adoption, and dispatch recovery
status: validation
source: captain (2026-11-21)
started: 2026-11-22T17:09:25Z
completed:
verdict:
score: 0.70
sprint: paseo-backend
worktree: .worktrees/spacedock-ensign-paseo-worker-adoption-and-recovery
issue:
pr: "#895"
mod-block: merge:pr-merge
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
        - id: gate:x9mfvnkqtbwgm4tba7dy41tw:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:x9mfvnkqtbwgm4tba7dy41tw-ideation-1
              briefing:
                id: briefing:x9mfvnkqtbwgm4tba7dy41tw:ideation:attempt-1:revision-1
                digest: sha256:c53021b2ee35fe202f460f3e48e941a2bfc4632d9c886892c2fbf2282fd68d3f
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:x9mfvnkqtbwgm4tba7dy41tw:ideation:1
                briefing: briefing:x9mfvnkqtbwgm4tba7dy41tw:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-11-22T18:02:04.639076Z"
                decision: approve
                reason: Captain approved ideation gate after preflight staff-review revisions
              application:
                target-stage: implementation
                state: consumed
        - id: gate:x9mfvnkqtbwgm4tba7dy41tw:validation
          stage: validation
          attempts:
            - id: gate-attempt:x9mfvnkqtbwgm4tba7dy41tw-validation-1
              briefing:
                id: briefing:x9mfvnkqtbwgm4tba7dy41tw:validation:attempt-1:revision-1
                digest: sha256:2aac481be9bd4541632cb902668e1a0d75e652640767e15c697a6a71039540e3
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:x9mfvnkqtbwgm4tba7dy41tw:validation:1
                briefing: briefing:x9mfvnkqtbwgm4tba7dy41tw:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-11-23T17:09:50.275647Z"
                decision: approve
                reason: Captain approved validation gate toward stacked PR targeting paseo-dispatch-backend (including push and PR creation)
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:x9mfvnkqtbwgm4tba7dy41tw-validation-2
              briefing:
                id: briefing:x9mfvnkqtbwgm4tba7dy41tw:validation:attempt-2:revision-1
                digest: sha256:be3f491dbc2a29dced471e630e48a5ec0513098a6b16f28a677c484b7ab85951
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:x9mfvnkqtbwgm4tba7dy41tw:validation:2
                briefing: briefing:x9mfvnkqtbwgm4tba7dy41tw:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-11-23T21:49:54.892849Z"
                decision: approve
                reason: 'Captain approved validation gate attempt-2 toward pushing ca39881e2 to stacked PR #895 and running stack-tip live CI'
              application:
                target-stage: done
                state: pending
---

Persist the Paseo worker handle in the task so a restarted or compacted first officer adopts a live worker instead of double-dispatching, and handle wait timeouts, lost create receipts, dead workers, and duplicate completion events with deterministic single-retry-then-hold recovery isolated per entity.

## Problem

With `paseo-dispatch-backend` (`936`) alone, a Paseo-hosted first officer (FO) that restarts, compacts, or crashes mid-stage has no durable record of the in-flight Paseo agent ID, workspace ID, or pre-dispatch state revision. On resume, a fresh FO risks spawning a second worker on the same stamped worktree/branch or mistaking a prior stage report for new completion evidence. Similarly, Paseo's `paseo wait --timeout` exits `0` with JSON `status=timeout` during normal live progress, a dropped `paseo run --no-wait` call can leave an unrecorded worker running on the daemon (lost create receipt), and an ended worker with a missing/incomplete report or a duplicate finish event needs exact, zero-ambiguity recovery rules so one failing task never wedges unrelated tasks.

## Proposed approach

Extend `skills/first-officer/references/paseo-fo-dispatch.md` (shipped by `936`) and `skills/first-officer/references/fo-write-core.md` with the lean `## Paseo Worker` record, deterministic restart adoption and lost-create-receipt reconciliation, and the 11-row recovery matrix; add a deterministic `status` compatibility fixture test in `internal/status/paseo_worker_compat_test.go`; and prove restart adoption end-to-end against a real isolated `paseo` daemon in `internal/ensigncycle/paseo_adoption_live_test.go` (`TestLivePaseoAdoptionSmoke`), consuming `936`'s package-private `internal/ensigncycle/paseo_live_helpers_test.go` seam and `internal/release/live_models.txt` SSOT with zero edits to `936`'s test files. No new `spacedock` CLI subcommands, no uncalled production Go simulator, and no fake ACP stub servers are introduced.

### 1. Lean worker-record schema stored in the task (`## Paseo Worker`)

Store one authoritative FO-owned H2 section `## Paseo Worker` in the canonical task body (in `.spacedock-state/{slug}/index.md` or `.spacedock-state/{slug}.md` for split-root workflows; the active entity file otherwise). Using H2 (`## Paseo Worker`) guarantees `internal/status/section_read.go` and `internal/status/gate_extract.go` cleanly bound it independently of any `## Stage Report: <stage>` section (an H3 appended after a stage report would be swallowed into that report's span).

**Write ownership & commit rules:**
- Owned exclusively by the FO (never edited by ensigns). Updated **in place** via the section span returned by `spacedock status --read <entity> --json` ONLY on launch/adoption (`state: running`) and terminal attempt outcome (`archived`, `superseded`, or `held`) — no intermediate `verified` write.
- Committed path-scoped in the state checkout (`git -C {state_checkout} add -- {entity_path} && git -C {state_checkout} commit -m "dispatch: record paseo worker for {slug} ({stage})" -- {entity_path}`), retrying once after 2s on `index.lock` contention, then synced per the split-root contract (using `GIT_SSH_COMMAND='ssh -o BatchMode=yes'`).
- **Gate briefing freeze safety:** For gated stages, the FO MUST update `## Paseo Worker` to `state: archived` and commit it BEFORE running `spacedock gate prepare`, so the frozen briefing digest is never invalidated after prepare.

**Exact fenced YAML schema (lean 10-field record + optional `hold_reason`, aligned with `936`'s `«worker-identity»`):**

```markdown
## Paseo Worker

```yaml
agent_id: 0cb685d1-8e49-4795-9705-d0b7759a5045
workspace_id: wks_8aec77dc0eda8e94
server_id: srv_j3VeUq-Ba722
stage: implementation
epoch: paseo-qjyrql-implementation-c1-a1
cwd: /private/tmp/spacedock-paseo.qjyrQl/fixture/.worktrees/spacedock-ensign-paseo-qjyrql
pre_dispatch_state_rev: 764719a10a5c397629231fa6db16565d8c617aac
requested:
  provider: codex/<model>
  mode: auto
  thinking: ""
returned:
  provider: codex/<model>
  mode: full-access
  thinking: ""
state: running
```
```

**Field contract (aligned with Paseo 0.17.1 `inspect` PascalCase / `ls` camelCase JSON and `936`'s shared label contract):**
- `agent_id`: Paseo `agentId` (`inspect.Id` / `ls[].id`) from `paseo run --no-wait` or recovered via `paseo ls`.
- `workspace_id`: Paseo `workspaceId` from Verb 1 find-or-create (`paseo workspace ls --json` -> `paseo workspace create`).
- `server_id`: daemon server ID from `paseo status --json` / `paseo.parent-server-id` label; `"local"` when a local isolated daemon omits a server ID.
- `stage`: entity `status:` at dispatch time.
- `epoch`: exact `spacedock.epoch` label passed to `paseo run` — deterministically `{worker_base}-c{cycle}-a{attempt}` (lowercase `[a-z0-9-]`, unique per `(entity, stage, cycle, attempt)`, recomputable cold from durable entity state with zero session memory).
- `cwd`: resolved physical working directory (`filepath.EvalSymlinks` / `pwd -P`) adopted by Verb 1 (matches `inspect.Cwd`; relative worktree path remains in `frontmatter.worktree`).
- `pre_dispatch_state_rev`: `git -C {state_checkout} rev-parse HEAD` captured immediately AFTER `dispatch build --stamp` and BEFORE `paseo run --no-wait` (prevents pre-existing or prior-cycle stage reports from falsely satisfying verification).
- `requested` (`provider`, `mode`, `thinking`): exact `--provider <provider>/<model>`, `--mode` (`""` if omitted), and `--thinking` (`""` if omitted) passed to `paseo run`.
- `returned` (`provider`, `mode`, `thinking`): effective values returned by `paseo inspect` (`<Provider>/<Model>` or `<Provider>`, `Mode`, `Thinking`; explicitly capturing discrepancies such as requested `auto` returning `full-access`).
- `state`: closed 4-value enum:
  - `running` — worker is live (dispatched or adopted) and awaiting completion.
  - `archived` — durable Git proof passed the completion gate and the worker was archived (`paseo archive <agent_id> --json`; `inspect.Archived == true`); terminal clean state for this attempt.
  - `superseded` — worker failed/ended without valid completion and was archived before launching a replacement.
  - `held` — entity exhausted its single retry or hit an unresolvable ambiguity/permission block; only this entity is held.
- `hold_reason` (optional): non-empty concrete one-line diagnostic present only when `state == held`.

**Shared Paseo label contract (frozen identically across `936` and `x9m`, preserving `internal/dispatch/names.go` <= 56-char budget):**
- `spacedock.worker`: **ALWAYS invariant `{worker_base}`** (canonical `envelope.name` from `spacedock dispatch build --host paseo` or read-only `spacedock dispatch name`). NEVER append `-retry` or `-cycleN-retry` to `spacedock.worker` — keeping `spacedock.worker` invariant ensures `paseo ls --global --label "spacedock.worker={worker_base}" --json` enumerates **every** attempt (initial, rework, and retry) for `(entity, stage)` and stays within `names.go`'s suffix cap.
- `spacedock.epoch`: **`{worker_base}-c{cycle}-a{attempt}`**, where `cycle = 1 + count(prior feedback cycles for stage in ### Feedback Cycles)` and `attempt = 1` initially or `2` on the single allowed retry (`{worker_base}-c1-a1` for v1 fresh dispatch in `936`, `{worker_base}-c1-a2` on cycle-1 retry, `{worker_base}-c2-a1` on feedback cycle 2).

**Retry ledger (`### Dispatch Retries`):**
Reuse the existing FO-owned `### Dispatch Retries` entity-body subsection (distinct from `### Feedback Cycles`, so a dispatch retry never increments the feedback counter):
```text
- Retry 1: {stage} (cycle {cycle}) — {spawn-error | ended-no-report | incomplete-report}; re-dispatched a2
```
(Lines without `(cycle N)` parse as `cycle 1` for backward compatibility.)

### 2. Deterministic FO restart adoption & lost-create-receipt reconciliation

**Shared Paseo 0.17.1 CLI invariants (aligned with `936`):**
- **Find-or-create workspace (Verb 1):** Before spawning (especially on `-retry`, feedback rework, or non-worktree stages), call `paseo workspace ls --json` and reuse any active workspace whose canonical `cwd` equals `$code_cwd`; call `paseo workspace create --isolation local --path "$code_cwd" ... --json` only when no matching workspace exists (Paseo `workspace create` is not idempotent and leaks duplicate workspaces on repeated calls).
- **Stdout-only JSON parsing:** `paseo run --workspace ... --json` prints non-JSON `Using workspace wks_...` to stderr even on exit `0`. Always parse stdout only (`2>/dev/null` when piping to `jq`, or separate stdout/stderr buffers in Go via `(*paseoDaemon).runJSON` — never `CombinedOutput`).

On boot, compaction resume, or `«dispatch.next-action»()` for an entity in active stage `S = entity.status`:

1. **Step 0 — Derive expected identity from durable entity state alone:**
   - Read frontmatter (`status`, `started`, `worktree`), `## Paseo Worker` (if any), `### Feedback Cycles`, and `### Dispatch Retries`.
   - Run read-only `spacedock dispatch name --workflow-dir {workflow_dir} --entity-path {entity_path} --stage {S}` to get invariant `worker_base`.
   - Compute `cycle = 1 + count(feedback cycles for S)`, `retries_done = count(retry entries for (S, cycle))`, `attempt = 1 + retries_done`.
   - Compute `expected_epoch = {worker_base}-c{cycle}-a{attempt}` and canonical physical `expected_cwd` (`filepath.EvalSymlinks` / `pwd -P` of `worktree` if set, else repo root).
2. **Step 1 — Check terminal/held/completed durable state first:**
   - If `## Paseo Worker` matches `(stage == S, epoch == expected_epoch)` and `state == held`: skip this entity (remain held) and continue unrelated entities.
   - If `## Paseo Worker` matches `(stage == S, epoch == expected_epoch)` and `state == archived`: if `HasValidNewReport` is true, resume completion advance (`merge guard` / `gate prepare` / advance to next stage) with zero new spawn; otherwise proceed to Step 2.
3. **Step 2 — Primary lookup via recorded `agent_id` (record matches `(stage == S, epoch == expected_epoch)` and `state == running`):**
   - Run `paseo inspect <record.agent_id> --json` and route through the Recovery Matrix below (`AGENT_NOT_FOUND` counts as ended).
4. **Step 3 — Lost-create-receipt reconciliation (no matching `## Paseo Worker` for `(S, expected_epoch)`, empty `agent_id`, or `paseo run --no-wait` failed/dropped without returning a receipt):**
   - Query labels in exact priority order BEFORE ever calling `paseo run --no-wait`:
     1. **Epoch lookup:** run `paseo ls --global --label "spacedock.epoch={expected_epoch}" --json` **without `--all`** first (excluding archived agents server/CLI-side); retry with `--all` ONLY if that returns `[]`.
     2. **Invariant worker-label fallback (only if step 3.1 returns `[]`):** run `paseo ls --global --label "spacedock.worker={worker_base}" --json` **without `--all`** first, retrying with `--all` ONLY if empty.
     3. **Inspect-based candidate filtering (NEVER filter on `ls[].cwd`, `ls[].created`, or `ls[].status` for archival — in Paseo 0.17.1 `ls[].cwd` is `~`-shortened, `ls[].created` is human relative time, and `ls[].status` never equals `"archived"`):** call `paseo inspect <candidate.id> --json` on each candidate and match ONLY rows where:
        - `inspect.Cwd == expected_cwd` (unshortened physical path), AND
        - (when `frontmatter.started` is non-empty) ISO-8601 `inspect.CreatedAt >= frontmatter.started`, AND
        - prefer `inspect.Archived == false` (keep archived matches only if no non-archived candidate exists and `HasValidNewReport` may already be satisfied).
   - Branch on filtered candidate count:
     - **Exactly 1 match `W`:** backfill & commit `## Paseo Worker` from `inspect` (using the `dispatch: {slug} entering {S}` state commit SHA as `pre_dispatch_state_rev` when missing); adopt `W` and enter the Recovery Matrix with **0 retries spent** and **0 duplicate spawn**.
     - **> 1 active (`inspect.Archived == false && inspect.Status != "closed"`) matches:** ambiguous roster — set `state: held`, `hold_reason: "ambiguous-active-workers: <id1>,<id2>"`, commit path-scoped, surface to captain, do NOT spawn another worker, do NOT kill unowned agents, continue unrelated entities.
     - **0 matches:** if triggered by ordinary boot/engage, proceed to `SpawnFresh` (using Verb 1 find-or-create); if triggered by a failed `paseo run --no-wait`, classify as confirmed `spawn-error` and enter Recovery Matrix rows 8–9.

### 3. Exact recovery matrix (zero day-1 ambiguities)

**Durable Git verification predicate `HasValidNewReport(entity, stage, pre_dispatch_state_rev)`** — all four must hold:
1. `git -C {state_checkout} rev-list {pre_dispatch_state_rev}..HEAD -- {entity_rel}` is non-empty and the committed entity carries a new `## Stage Report: {stage}` passing `spacedock status --read <entity> --stage {stage} --checklist --json` (every retained checklist item accounted for, non-empty evidence lines, non-empty `### Summary`, no `FAILED` items except feedback-review findings).
2. `git -C {state_checkout} status --porcelain -- {entity_rel}` is empty.
3. For worktree stages (`worktree != ""`), the code worktree is on the expected task branch with clean `git -C {worktree} status --porcelain`.
4. `record.stage == entity.status` and `record.epoch == expected_epoch`.

**FO self-preservation invariant (Paseo 0.17.1 cascade-archive rule):** Never pass `$PASEO_AGENT_ID` to `paseo archive`, and never archive an FO while any child worker is `running`. In Paseo 0.17.1, archiving a parent agent cascade-kills and archives every child sharing its workspace (which includes all non-worktree stage workers sharing the repo-root workspace) and detaches `paseo.parent-agent-id` on cross-workspace children.

| Row | Trigger / Observed Daemon State | `HasValidNewReport` | `retries_done` | Exact FO Decision & Action | Retries Spent | Archival | Unrelated Entities |
|---|---|---|---|---|---|---|---|
| 1 | `paseo wait` exits `0` with `status=timeout`; `paseo inspect` shows `Status=running` & `PendingPermissions=[]` | Any | Any | `ContinueWait`: keep `state: running`; re-issue `paseo wait <agent_id> --timeout 300 --json` (never nudge, kill, or retry). | 0 | None | Unaffected |
| 2 | `paseo wait` exits `0` with `status=permission`, OR `paseo inspect` shows non-empty `PendingPermissions` | Any | Any | `HoldEntity` (permission-blocked): set `state: held`, `hold_reason: "pending-permission: <agent_id>"`; commit path-scoped; surface to captain. | 0 | None (live worker preserved) | Unaffected |
| 3 | Restarted FO finds `state=running` (recorded or adopted from Step 3) and `paseo inspect` shows `Status=running` | Any | Any | `AdoptRunning`: persist/confirm `## Paseo Worker` (`state: running`) and enter `paseo wait <agent_id> --timeout 300 --json` loop. Zero duplicate spawn. | 0 | None | Unaffected |
| 4 | Worker ended (`Status` is `idle` / `closed` / `error`, `inspect.Archived == true`, or `AGENT_NOT_FOUND`) or `wait` returned `idle` | `true` | Any | `VerifyAndArchive`: run idempotent `paseo archive <agent_id> --json` (treat `AGENT_ALREADY_ARCHIVED` as success) -> update `## Paseo Worker` `state: archived` & commit path-scoped (BEFORE `gate prepare` on gated stages) -> advance / `gate prepare` / `merge guard`. | 0 | Archive verified worker | Unaffected |
| 5 | Worker ended (`Status` is `idle` / `closed` / `error`, `inspect.Archived == true`, or `AGENT_NOT_FOUND`) or `wait` returned `error` | `false` (no new report or uncommitted/incomplete report) | `0` | `ArchiveAndRetry`: archive superseded worker (`paseo archive <agent_id> --json`, or `--force` only if still running on explicit supersede; ignore `AGENT_ALREADY_ARCHIVED`/`AGENT_NOT_FOUND`); set `state: superseded`; append Retry 1 line (`ended-no-report` or `incomplete-report`) to `### Dispatch Retries`; commit path-scoped; reuse existing workspace via Verb 1 find-or-create and spawn ONE fresh replacement with invariant `spacedock.worker={worker_base}`, `spacedock.epoch={worker_base}-c{cycle}-a2` (reusing existing worktree WITHOUT `--stamp`, appending WIP check note), record new `## Paseo Worker` (`state: running`). | +1 (1/1) | Archive superseded worker BEFORE retry spawn | Unaffected |
| 6 | Worker ended (`Status` is `idle` / `closed` / `error`, `inspect.Archived == true`, or `AGENT_NOT_FOUND`) or `wait` returned `error` | `false` | `>= 1` | `HoldEntity` (retry-exhausted): archive failed attempt-2 worker (`paseo archive <agent_id> --json` if not archived); set `state: held`, `hold_reason: "retry-exhausted: {stage} cycle {cycle} ({cause})"`; commit path-scoped; surface to captain; stop attempting this entity. | 0 (cap) | Archive failed attempt-2 worker | **Unaffected** |
| 7 | `paseo run --no-wait` fails/drops without receipt; Step 3 finds **1** matching worker | Re-evaluated via rows 3–6 | Any | `AdoptFromLabels`: backfill `## Paseo Worker` from `paseo inspect`, commit path-scoped, then route through row 3 (if `running`) or rows 4–6 (if ended). | 0 | Per adopted row | Unaffected |
| 8 | `paseo run --no-wait` fails/drops; Step 3 finds **0** matching workers (`spawn-error`) | N/A | `0` | `ArchiveAndRetry` (`spawn-error`): append Retry 1 line (`spawn-error`) to `### Dispatch Retries`; commit path-scoped; reuse workspace via Verb 1 find-or-create and spawn ONE fresh replacement (`spacedock.worker={worker_base}`, `spacedock.epoch={worker_base}-c{cycle}-a2`). | +1 (1/1) | None (nothing created) | Unaffected |
| 9 | `paseo run --no-wait` fails/drops; Step 3 finds **0** matching workers (`spawn-error`) | N/A | `>= 1` | `HoldEntity` (`spawn-error` exhausted): set `state: held`, `hold_reason: "retry-exhausted: {stage} cycle {cycle} (spawn-error)"`; commit path-scoped; surface to captain; stop attempting this entity. | 0 (cap) | None | **Unaffected** |
| 10 | Duplicate/late finish event, repeated `paseo wait`, or re-entrant check for `(agent_id, epoch)` already `archived`/`superseded`, or `entity.status != record.stage` | N/A | Any | `IgnoreDuplicate`: strict no-op (zero state mutation, zero retry, zero spawn). | 0 | No-op | Unaffected |
| 11 | Step 3 finds **> 1** active (`!inspect.Archived && inspect.Status != "closed"`) matching workers for `(entity, stage, cycle)` | N/A | Any | `HoldEntity` (ambiguous-roster): set `state: held`, `hold_reason: "ambiguous-active-workers: <ids>"`; commit path-scoped; surface to captain; do NOT spawn another worker or kill unowned agents. | 0 | None (captain decides) | **Unaffected** |

### 4. `TestLivePaseoAdoptionSmoke` specification (`internal/ensigncycle/paseo_adoption_live_test.go`, `//go:build live`)

Consumes `936`'s package-private `internal/ensigncycle/paseo_live_helpers_test.go` (`startIsolatedPaseoDaemon(t)`, `(*paseoDaemon).runJSON(t, args...)` with separate stdout/stderr buffers, and `setupPaseoSplitRootFixture(t, ...)` staging installed skills via `livePluginDir(t)` with zero raw `<checkout>/skills/...` paths) and reads the default Codex provider from `"codex/" + release.LiveModels["codex.exec"]` (`internal/release/live_models.txt` SSOT — zero hardcoded model literals, zero edits to `936`'s test files). Local overrides: `SPACEDOCK_PASEO_BIN`, `SPACEDOCK_PASEO_LIVE_PROVIDER`, `SPACEDOCK_PASEO_LIVE_THINKING`, `SPACEDOCK_PASEO_LIVE_FO_PROVIDER`, `SPACEDOCK_PASEO_LIVE_REQUIRED`.

**Deterministic mid-stage stop choreography (zero sleep races, zero cascade-archive corruption):**
1. **Fixture:** Split-root temp workflow (`README.md` with `commissioned-by: spacedock@test`, `state: .spacedock-state`, stages `backlog -> implementation (worktree: true) -> done (terminal: true)`, no `_mods/`). The `implementation` stage definition instructs the worker to: (a) write and commit `receipt.txt` on the task branch, (b) poll up to 180s (1s interval) until `<fixtureRoot>/.adoption-release` exists, then (c) append `## Stage Report: implementation` and commit it path-scoped in `.spacedock-state`.
2. **Phase 1 (FO-1 dispatch & mid-stage stop via `paseo stop`):** Launch FO-1 on the isolated daemon (invoking the installed `first-officer` skill by name) with the entity ready in `implementation`. Poll via `daemon.runJSON(t, "ls", "--all", "--global", "--label", "spacedock.worker="+workerBase, "--json")` and the entity file until: (a) exactly 1 worker `W1` exists with `status == "running"`, and (b) `## Paseo Worker` is committed in `.spacedock-state` with `agent_id == W1` and `state: running`. Immediately stop FO-1 using **`paseo stop "$fo1AgentId" --json`** (`cancelAgent`, NOT `paseo archive --force`, which would cascade-archive same-workspace children or detach `ParentAgentId`), and register FO-1 archival in `t.Cleanup`. Verify via `paseo inspect` that `W1` remains `Status == "running"` with `ParentAgentId == fo1AgentId` intact.
3. **Phase 2 (Fresh FO-2 boot & adoption):** Launch a brand-new FO-2 (fresh session, zero shared state with FO-1, invoking the installed `first-officer` skill by name) on the same workflow and isolated daemon. Once FO-2 is active (and has begun adopting `W1` while `W1` is `running`), create `<fixtureRoot>/.adoption-release`. Wait for FO-2 via `daemon.runJSON(t, "wait", fo2AgentId, "--timeout", "60", "--json")` until `status == "idle"`.
4. **Phase 3 (Durable assertions using test-local YAML parsing):**
   - `paseo ls --all --global --label "spacedock.worker=<worker_base>" --json` (and all non-FO daemon agents for this workflow) contains **exactly 1 worker** (`W1` — because `spacedock.worker` is invariant `{worker_base}`, any duplicate or rogue retry would appear here and fail the test).
   - `paseo inspect "$W1" --json` confirms `Archived == true` (`ArchivedAt != null`), while FO-2 remains `Archived == false`.
   - Locate the task entity (in `.spacedock-state/_archive/` or `.spacedock-state/`): test-local YAML helper parses `## Paseo Worker` and asserts `agent_id == W1`, non-empty `workspace_id` and `server_id`, `stage == "implementation"`, `epoch == workerBase+"-c1-a1"`, canonical `cwd`, non-empty `pre_dispatch_state_rev`, non-empty `requested.provider` and `returned.provider`, `returned.mode`, and `state == "archived"`.
   - Test-local retry helper asserts `### Dispatch Retries` contains `0` entries (neither FO-1 stop nor wait timeout spent a retry).
   - `## Stage Report: implementation` is committed path-scoped strictly after `pre_dispatch_state_rev`; `receipt.txt` is committed on the task branch; frontmatter `status` durably advanced to `done` (`verdict: PASSED`); checkouts are clean; `t.Cleanup` archives FO-1 and stops the isolated daemon cleanly (`localDaemon == "stopped"`).

### 5. Mechanism-to-value justification

| Mechanism | Value AC served | Simplest alternative considered | Why alternative is insufficient |
|---|---|---|---|
| Lean 10-field H2 `## Paseo Worker` YAML block (4-state enum, no `verified`) | AC-1, AC-3 | Flat YAML frontmatter fields via `status --set`, or shared `.paseo-workers.json` sidecar | `status --set` only writes flat scalars (cluttering frontmatter with transient keys), while a shared sidecar breaks per-entity path-scoped commits and causes rebase conflicts across sibling tasks. H2 terminates cleanly in `section_read.go` and `gate_extract.go` with zero `status` binary changes and zero uncalled Go production code. |
| Invariant `spacedock.worker={worker_base}` + deterministic `spacedock.epoch={worker_base}-c{cycle}-a{attempt}` + `inspect`-verified Step 3 lookup | AC-1 | Mutating `spacedock.worker` with `-retry`/`-cycleN-retry`, random UUID epochs, or filtering `paseo ls --json` directly | Compound suffixes break `internal/dispatch/names.go`'s 8-char suffix budget and hide duplicate retries from `spacedock.worker` queries; random epochs cannot recover lost create receipts; Paseo 0.17.1 `ls --json` shortens `cwd`, formats `created` as relative prose, and omits `archived`. |
| `paseo stop` mid-stage + `.adoption-release` file sentinel in `TestLivePaseoAdoptionSmoke` | AC-1 | `paseo archive --force` + fixed `time.Sleep` | `paseo archive --force` cascade-kills same-workspace children and strips `ParentAgentId` on cross-workspace children; fixed sleeps race variable model latency. `paseo stop` + sentinel deterministically holds `W1` live across FO-1 stop and FO-2 boot. |
| `TestPaseoWorkerStatusCompatibility` in `internal/status/` | AC-3 | Relying only on `TestLivePaseoAdoptionSmoke` | Live tests do not exercise all four `state` values (`running|archived|superseded|held` + `hold_reason`) or both pre-report and post-report placements against `--validate`, `--checklist`, and `--ac-scan`. |

### 6. Concrete skill and documentation diffs

#### A. `skills/first-officer/references/fo-write-core.md` (before / after)

- **Mutation Gate `allowed-state` rule cell:**
  - *Before:* `Entity frontmatter via ${SPACEDOCK_BIN:-spacedock} status --set, spacedock new, archive moves, state-transition commits, and ### Feedback Cycles under the existing state/worktree rules. Roadmap and sprint docs under _roadmap/ are edited directly and committed path-scoped in the state checkout.`
  - *After:* `Entity frontmatter via ${SPACEDOCK_BIN:-spacedock} status --set, spacedock new, archive moves, state-transition commits, and FO-owned operational body sections (### Feedback Cycles, ### Dispatch Retries, ## Paseo Worker) under the existing state/worktree rules. Roadmap and sprint docs under _roadmap/ are edited directly and committed path-scoped in the state checkout.`
- **`## FO Write Scope` bullet 3 & off-limits paragraph:**
  - *Before:* ``- **`### Feedback Cycles` section** — in entity bodies, tracking rejection rounds. ...`` and ``and entity body content beyond `### Feedback Cycles` (stage reports, design, implementation notes belong to dispatched workers).``
  - *After:* ``- **FO-owned operational body sections (`### Feedback Cycles`, `### Dispatch Retries`, `## Paseo Worker`)** — tracking rejection rounds, single-retry dispatch ledgers, and Paseo worker handles. Update in place and commit path-scoped in the state checkout (or worktree/main when non-split-root).`` and ``and entity body content beyond `### Feedback Cycles`, `### Dispatch Retries`, and `## Paseo Worker` (stage reports, design, implementation notes belong to dispatched workers).``

#### B. `skills/first-officer/references/paseo-fo-dispatch.md` (new sections appended after `936`'s 5-verb binding; no ensign `references/*.md` paths cited)

Add three concise operational sections (avoiding any `skills/ensign/` `references/*.md` paths so `TestBootResidentDeferredLoadPointsResolve` stays green):
1. `## Paseo Worker Record` — lean 10-field YAML block (`agent_id`, `workspace_id`, `server_id`, `stage`, `epoch`, `cwd`, `pre_dispatch_state_rev`, `requested.{provider,mode,thinking}`, `returned.{provider,mode,thinking}`, `state: running | archived | superseded | held`, optional `hold_reason`), invariant `spacedock.worker={worker_base}` + `spacedock.epoch={worker_base}-c{cycle}-a{attempt}`, stdout-only JSON parsing rule, Verb 1 find-or-create workspace reuse, path-scoped commit rule, and pre-`gate prepare` archival order.
2. `## Restart Adoption and Lost-Create-Receipt Reconciliation` — Steps 0–3 (`spacedock dispatch name`, `paseo inspect`, `paseo ls --global` without `--all` first then `--all` fallback, filtering exclusively on `inspect.Cwd`, ISO-8601 `inspect.CreatedAt`, and boolean `inspect.Archived`).
3. `## Recovery Matrix and Retry Ledger` — the 11-row table, `HasValidNewReport` 4-part Git predicate, `### Dispatch Retries` line format, per-entity isolation rule, and `$PASEO_AGENT_ID` self-preservation guard explaining Paseo's same-workspace cascade-archive behavior.

#### C. User-visible documentation (`docs/runtime-live-ci.md`; `docs/site/**` unchanged, zero hardcoded model literals)

No `docs/site/**` page documents the opt-in Paseo dispatch backend yet, so `docs/site/**` needs no change. Update `### Paseo live smoke` in `docs/runtime-live-ci.md` (added by `936`, STE compliant, omitting hardcoded model IDs per `docs/runtime-live-ci.md:135`):

- *Before (from `936`):*
  ```markdown
  Run the isolated-daemon Paseo smoke locally with the registered Codex model (default from `internal/release/live_models.txt`) or override the binary, provider, and thinking level:

  ```bash
  go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle -v
  ```
  ```
- *After:*
  ```markdown
  Run the isolated-daemon Paseo smoke tests locally with the registered Codex model (default from `internal/release/live_models.txt`) or override the binary, provider, and thinking level:

  ```bash
  go test -tags live -count=1 -timeout 15m -run '^TestLivePaseo(Dispatch|Adoption)Smoke$' ./internal/ensigncycle -v
  ```

  `TestLivePaseoAdoptionSmoke` stops the first officer mid-stage with `paseo stop`. It verifies that a fresh first officer adopts the running worker, advances state, and archives the worker.
  ```

## Risk evidence

no spike needed: all underlying Paseo 0.17.1 and Spacedock primitives are verified directly against source and prior live probes:
- `_roadmap/paseo-backend/probe-report.md`, `run-log.md`, and `staff-review.md` verified: `paseo ls --global --label` (omitting `--all` excludes archived agents via `ls.ts:229-231`), `paseo inspect <id> --json` (exposes unshortened `Cwd`, ISO-8601 `CreatedAt`, boolean `Archived`, `Status`, `PendingPermissions`, `ParentAgentId` via `inspect.ts:14-49, 129-155`), `paseo wait --timeout --json` exiting `0` for `idle|timeout|permission|error` (`wait.ts:102-144`), `paseo stop` calling `cancelAgent` without cascade-archiving (`stop.ts:89`), `paseo workspace ls --json` returning unshortened `cwd`/`workspaceId` (`workspace/ls.ts:18-48`), and `paseo run --workspace --json` writing `Using workspace ...` to stderr on exit `0` (`run.ts:553`).
- `internal/dispatch/names.go:15-24, 102-147` verified `spacedock dispatch name` and the 8-char suffix budget (which invariant `spacedock.worker={worker_base}` never breaches).
- `internal/status/section_read.go:153-164` and `internal/status/gate_extract.go:127, 143-153` verified H2 `## Paseo Worker` terminates `## Stage Report:` cleanly without contaminating checklist or `--ac-scan` spans.
- `internal/contractlint/fo_write_core_mutation_gate_test.go` parses only columns 1–2 (`class`, `patterns`) of `FO-WRITE-CLASSIFIER`, proving the `allowed-state` rule-text update preserves all classifier guards.
- **Sequencing dependency:** Strict serial execution — `936` (`paseo-dispatch-backend`) MUST merge to `main` (landing `paseo-fo-dispatch.md`, `internal/ensigncycle/paseo_live_helpers_test.go`, and `docs/runtime-live-ci.md`'s Paseo section) before `x9m` starts implementation.

## Out of scope

- Remote Paseo workers (`serverId` across machines) and cross-host adoption.
- Worker reuse across stages (`paseo send` steering or multi-stage worker continuity).
- Fake ACP stub servers, uncalled production Go reconciliation simulators, or Paseo version-pinning machinery.
- Multi-FO concurrency leases (tracked separately by `79` / `entity-session-claim-lease`).
- New `spacedock` CLI subcommands wrapping `paseo`.

## Expected surface and tolerance

Estimate net LOC change: +278, across 5 files.

Breakdown: **+288 insertions / -10 deletions (net +278 LOC)** across 5 files. Tolerance: **+195 to +365 net LOC, across 4 to 7 files**. Zero edits to `936`'s test files (`paseo_live_helpers_test.go` is consumed as-is).

| File | Est. insertions | Est. deletions | Net LOC | Purpose |
|---|---|---|---|---|
| `skills/first-officer/references/paseo-fo-dispatch.md` | +80 | -5 | +75 | Lean 10-field `## Paseo Worker` schema, invariant labels, Steps 0–3 `inspect`-filtered adoption, 11-row recovery matrix, retry ledger, cascade-archive guard |
| `skills/first-officer/references/fo-write-core.md` | +4 | -4 | 0 | Authorize `## Paseo Worker` and `### Dispatch Retries` as FO-owned operational body sections |
| `internal/status/paseo_worker_compat_test.go` | +70 | 0 | +70 | `TestPaseoWorkerStatusCompatibility` over fixture markdown (`--validate`, `--checklist`, `--ac-scan`) |
| `internal/ensigncycle/paseo_adoption_live_test.go` | +130 | 0 | +130 | `TestLivePaseoAdoptionSmoke` (`//go:build live`) + test-local YAML/retry parser helpers consuming `paseo_live_helpers_test.go` and `release.LiveModels["codex.exec"]` |
| `docs/runtime-live-ci.md` | +4 | -1 | +3 | Document local `TestLivePaseoAdoptionSmoke` command (STE, no hardcoded model IDs) |

**Declared observable semantic changes:**
- *Command grammar:* None (`spacedock` CLI flags and subcommands unchanged).
- *Stored formats:* Adds optional FO-owned H2 entity-body section `## Paseo Worker` (lean 10-field YAML record + optional `hold_reason`, closed 4-value `state` enum: `running | archived | superseded | held`) and extends `### Dispatch Retries` lines with explicit `(cycle N)` and Paseo causes (`spawn-error | ended-no-report | incomplete-report`). YAML frontmatter is unchanged.
- *Authority:* `fo-write-core.md` explicitly permits FO writes to `## Paseo Worker` and `### Dispatch Retries` in entity bodies; ensigns remain barred from editing them.
- *Runtime behavior:* Paseo-hosted FOs persist the lean worker record at launch and terminal attempt outcome, adopt live or lost-receipt workers across restarts via `## Paseo Worker` + `paseo ls`/`paseo inspect` (`spacedock.worker={worker_base}`, `spacedock.epoch={worker_base}-c{cycle}-a{attempt}`), reuse workspaces via find-or-create, parse stdout-only JSON, treat `paseo wait` timeouts as live progress (0 retries), enforce single-retry-then-hold isolated per entity, ignore duplicate/stale completion events, and archive verified/superseded workers without ever archiving `$PASEO_AGENT_ID` while children run.

## Acceptance criteria

**AC-1 (VALUE) - A fresh restarted Paseo FO adopts an in-flight worker mid-stage with zero duplicate dispatch, advances durable workflow state, and archives the worker.**
Verified by: `go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoAdoptionSmoke$' ./internal/ensigncycle -v` passes against an isolated `paseo` daemon (consuming `paseo_live_helpers_test.go` and defaulting provider to `"codex/" + release.LiveModels["codex.exec"]`; `SPACEDOCK_PASEO_BIN`, `SPACEDOCK_PASEO_LIVE_PROVIDER`, `SPACEDOCK_PASEO_LIVE_THINKING`, and `SPACEDOCK_PASEO_LIVE_FO_PROVIDER` overridable). The test stops FO-1 via `paseo stop "$fo1AgentId" --json` while `W1` is blocked on `.adoption-release`, verifies `W1` remains `Status == "running"` with `ParentAgentId == fo1AgentId`, boots fresh FO-2 via installed skill discovery, creates `.adoption-release`, and asserts: `paseo ls --all --global --label "spacedock.worker=<worker_base>" --json` contains **exactly 1 worker** (`W1`), `paseo inspect "$W1" --json` confirms `Archived == true` (`ArchivedAt != null`) while FO-2 remains `Archived == false`, test-local YAML parsing of `## Paseo Worker` confirms all 10 fields (`agent_id == W1`, `workspace_id`, `server_id`, `stage == "implementation"`, `epoch == workerBase+"-c1-a1"`, `cwd`, `pre_dispatch_state_rev`, `requested`/`returned` provider/mode/thinking, `state == "archived"`), `### Dispatch Retries` has `0` entries, the task-branch commit and path-scoped split-root stage report commit after `pre_dispatch_state_rev` exist, frontmatter `status` advanced durably to `done` (`verdict: PASSED`), and `t.Cleanup` archives FO-1 and stops the daemon. Falsifying edit: spawning a second worker (including a retry), dropping `## Paseo Worker` persistence, stopping FO-1 via `paseo archive --force`, spending a retry on FO restart/timeout, failing to advance, or leaving `W1` unarchived fails the test.

**AC-2 - `fo-write-core.md` authorizes FO writes to `## Paseo Worker` and `### Dispatch Retries` without loosening any product-write guard, and `paseo-fo-dispatch.md` preserves skill reference closure and FO byte caps.**
Verified by: `go test -race -count=1 ./internal/contractlint` passes, including `TestFOWriteCoreMutationGateClassifiesTargets`, `TestFOWriteCoreMutationGateAllowsRoadmapState`, `TestFOWriteCoreMutationGateRequiresExactOverride`, `TestFirstOfficerEntryEagerlyImportsOnlySharedCore`, `TestUserSkillReferenceClosureResolves`, `TestBootResidentDeferredLoadPointsResolve`, and `TestFOInstructionComponentCaps`. Falsifying edit: breaking any `FO-WRITE-CLASSIFIER` pattern in `fo-write-core.md`, leaking an ensign `references/*.md` path into `paseo-fo-dispatch.md`, or breaking reference closure/caps fails the suite.

**AC-3 - Spacedock validation and section/AC scanning remain clean on entities carrying the lean `## Paseo Worker` YAML block and `### Dispatch Retries` alongside `## Stage Report: <stage>`.**
Verified by: `go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$'` drives `status.NativeRunner` (`--validate`, `--read --stage implementation --checklist --json`, and `--read --stage implementation --ac-scan --json`) over split-root entity fixtures placing `## Paseo Worker` (across `running`, `archived`, `superseded`, and `held` + `hold_reason`) and `### Dispatch Retries` both before and after an earlier stage report, asserting exit `0`, zero warnings, intact AC enumeration, and zero checklist/summary contamination. Falsifying edit: demoting `## Paseo Worker` to an H3 after a stage report or emitting invalid frontmatter/YAML fails the test.

## Test plan

No existing primary proof owner covers FO restart adoption or `## Paseo Worker` entity compatibility (`TestLivePaseoDispatchSmoke` in `936` tests single-FO fresh dispatch only and cannot detect duplicate spawns across restart). Each check below owns one distinct failure class with zero uncalled production code:

1. **Deterministic `status` compatibility fixture suite (primary owner for AC-3):**
   - Command: `go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$'` (plus full `go test ./...` and `go test ./... -race`).
   - Falsifying edits: demoting `## Paseo Worker` to H3 inside a stage report, emitting invalid frontmatter, or corrupting `--checklist`/`--ac-scan` spans across any of the 4 `state` values (`running|archived|superseded|held`).
   - Cost/mode: deterministic, offline, `< 1s`, zero model spend.
2. **Existing `contractlint` guard suite (primary owner for AC-2):**
   - Command: `go test -race -count=1 ./internal/contractlint`.
   - Falsifying edits: breaking `FO-WRITE-CLASSIFIER` in `fo-write-core.md`, citing an ensign `references/*.md` path inside `paseo-fo-dispatch.md`, or breaching reference closure / byte caps.
   - Cost/mode: deterministic, offline, `< 2s`, zero model spend.
3. **Real isolated-daemon live adoption smoke (primary owner for AC-1):**
   - Command: `go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoAdoptionSmoke$' ./internal/ensigncycle -v` (default provider from `internal/release/live_models.txt`; optional local overrides `SPACEDOCK_PASEO_BIN`, `SPACEDOCK_PASEO_LIVE_PROVIDER`, `SPACEDOCK_PASEO_LIVE_THINKING`, `SPACEDOCK_PASEO_LIVE_FO_PROVIDER`).
   - Falsifying edits: FO-2 spawning a second worker or retry instead of adopting `W1`; FO-1 failing to commit `## Paseo Worker`; stopping FO-1 with `paseo archive --force` (detaching/killing `W1`); FO-2 burning a retry on restart/timeout; FO-2 failing to advance frontmatter or archive `W1`.
   - Cost/mode: live (real `paseo` daemon + real FO-1/FO-2/worker sessions), ~3–6 minutes, one targeted run per validation; skips cleanly locally when `paseo` or auth is absent unless `SPACEDOCK_PASEO_LIVE_REQUIRED=1`.
4. **High-stakes validation & required CI lanes:** Because `x9m` modifies `skills/first-officer/references/fo-write-core.md` and `paseo-fo-dispatch.md`, validation includes a read-only detached adversarial audit on a throwaway checkout (`spike-paseo-worker-adoption-and-recovery`), `gofmt -w ./cmd ./internal`, `spacedock status --workflow-dir docs/dev --validate`, and required green `claude-live`, `codex-live`, and `pi-live` CI lanes before merge.

### Feedback Cycles

- Cycle 1: REJECTED — detached adversarial audit (`0693a5c1`, `spike-paseo-worker-adoption-and-recovery`); surface 5/+358 vs estimate 5/+278 (+28.8%); AC unchanged
- Cycle 2: REJECTED — sprint pre-merge staff audit (`dd553ee3`, `_roadmap/paseo-backend/pre-merge-audit.md` `9e71ce2c3`); surface 5/+363 vs estimate 5/+278 (+30.6%); AC unchanged
- Cycle 3: REJECTED — Captain-authorized post-merge worktree assertion fix in `TestLivePaseoAdoptionSmoke` (`paseo_adoption_live_test.go:144-171`); surface 5/+356 vs estimate 5/+278 (+28.1%); AC unchanged

## Stage Report: ideation

- DONE: Define exact worker-record schema stored in the task and deterministic FO restart adoption + lost-create-receipt reconciliation rules using paseo ls/inspect labels (spacedock.worker, spacedock.epoch)
  Folded staff-review M1, M2, M3, M4, and M7: reduced `## Paseo Worker` to the lean 10-field YAML schema (4-value `state: running|archived|superseded|held` + optional `hold_reason`, dropping intermediate `verified`), froze invariant `spacedock.worker={worker_base}` and `spacedock.epoch={worker_base}-c{cycle}-a{attempt}`, fixed Step 3 lost-create-receipt lookup to query `paseo ls --global` without `--all` first and filter strictly on `inspect.Cwd`, ISO-8601 `inspect.CreatedAt`, and boolean `inspect.Archived`, and aligned Verb 1 workspace find-or-create and stdout-only JSON parsing with `936`.
- DONE: Define exact timeout vs dead-worker vs duplicate-event recovery matrix (single-retry-then-hold, per-entity isolation, verified/superseded archival) with zero day-1 ambiguities
  Defined the complete 11-row recovery matrix, `HasValidNewReport` 4-part Git predicate, cycle-scoped `### Dispatch Retries` ledger, `fo-write-core.md` before/after authorization wording, per-entity hold isolation, and `$PASEO_AGENT_ID` self-preservation guard documenting Paseo 0.17.1's same-workspace cascade-archive rule (M5).
- DONE: Specify TestLivePaseoAdoptionSmoke (stop FO mid-stage, boot fresh FO, assert single worker + durable advance + archive) and update Expected surface, scanner-exact Acceptance criteria with falsifying Verified by clauses, and Test plan
  Folded staff-review M1, M5, and M8: cut uncalled `internal/dispatch/paseo_worker.go`, `paseo_worker_test.go`, and `paseo_adoption_binding_test.go`, switched Phase 1 FO-1 stop to `paseo stop "$fo1AgentId" --json` (archiving FO-1 only in `t.Cleanup`), wired `TestLivePaseoAdoptionSmoke` to `936`'s `paseo_live_helpers_test.go` seam and `release.LiveModels["codex.exec"]` SSOT with test-local YAML parsing, trimmed Expected surface to +278 net LOC (+288/-10 across 5 files), and updated scanner-exact `## Acceptance criteria` (AC-1, AC-2, AC-3 with falsifying `Verified by:` clauses) and `## Test plan`; verified via `spacedock status --read --ac-scan`, `--checklist`, and `--validate`.

### Summary

Folded all `x9m` Material findings (M1, M2, M3, M4, M5, M7, M8) from the sprint-wide Preflight Staff Review (`docs/dev/.spacedock-state/_roadmap/paseo-backend/staff-review.md`, commit `8447132f3`) into `docs/dev/.spacedock-state/paseo-worker-adoption-and-recovery/index.md`. Cut ~295 LOC of uncalled production Go / prose-binding machinery, trimmed `## Paseo Worker` to the lean 10-field YAML schema, fixed Step 3 `paseo ls` -> `paseo inspect` filtering, froze invariant `spacedock.worker={worker_base}`, switched `TestLivePaseoAdoptionSmoke` FO-1 stop to `paseo stop`, and bound the live test to `936`'s `paseo_live_helpers_test.go` seam and `internal/release/live_models.txt` SSOT (surface reduced from +493 net / 7 files to +278 net / 5 files).


## Stage Report: implementation

- DONE: Extend skills/first-officer/references/paseo-fo-dispatch.md (lean 10-field H2 ## Paseo Worker YAML record, invariant spacedock.worker + deterministic spacedock.epoch, Steps 0-3 inspect-filtered lost-create-receipt adoption, 11-row recovery matrix, retry ledger, cascade-archive guard) and skills/first-officer/references/fo-write-core.md without touching first-officer-shared-core.md or citing ensign references/*
  Commit `b3f8eadd2` (stacked on `e1b971b32`) adds `## Paseo Worker Record`, `## Restart Adoption and Lost-Create-Receipt Reconciliation`, and `## Recovery Matrix and Retry Ledger` to `paseo-fo-dispatch.md` and authorizes `## Paseo Worker` / `### Dispatch Retries` in `fo-write-core.md` without editing `first-officer-shared-core.md` or citing ensign `references/*`; AC-2 proved by `go test -race -count=1 ./internal/contractlint` (asserts `FO-WRITE-CLASSIFIER` targets, deferred reference closure, and FO byte caps — breaking classifier patterns, citing an ensign `references/*.md` path, or breaching caps fails the suite).
- DONE: Add internal/status/paseo_worker_compat_test.go (TestPaseoWorkerStatusCompatibility across all 4 states and pre/post-report placements), internal/ensigncycle/paseo_adoption_live_test.go (TestLivePaseoAdoptionSmoke using paseo stop and 936's paseo_live_helpers_test.go as-is with zero edits to 936 test files and zero uncalled production Go simulators), and update docs/runtime-live-ci.md (STE, no hardcoded model literals)
  Commit `b3f8eadd2` adds `TestPaseoWorkerStatusCompatibility` (AC-3: drives `status.NativeRunner` `--validate`, `--read --json`, `--checklist --json`, and `--ac-scan --json` across `running|archived|superseded|held` × `before-report|after-report`, failing if `## Paseo Worker` is demoted to H3 or contaminates checklist/AC spans), adds `TestLivePaseoAdoptionSmoke` (AC-1: `//go:build live` consuming 936's `paseo_live_helpers_test.go` untouched with test-local YAML/retry parsing, `paseo stop`, and zero hardcoded models; verified compile + vet via `go test -tags live -c` and `go vet -tags live ./internal/ensigncycle`), and updates `docs/runtime-live-ci.md` in STE verified by `TestPaseoLiveWorkflowSecurityAndModelSSOT` (fails if a `codex.exec` model literal appears).
- DONE: Stay within +278 net LOC (+195 to +365, 4-7 files) over base paseo-dispatch-backend (e1b971b32); run gofmt, go test ./..., go test ./... -race, go test -tags live -c, and spacedock status --validate; commit cleanly on paseo-worker-adoption-and-recovery and write path-scoped ## Stage Report: implementation
  `git diff --shortstat e1b971b32..HEAD` is `5 files changed, 372 insertions(+), 14 deletions(-)` (+358 net LOC, within +195..+365 / 4-7 files); `gofmt -l`, `go vet ./...`, `go vet -tags live ./internal/ensigncycle`, `go test -tags live -c`, all 23 Go packages non-race and `-race` (`internal/cli` with `PATH=/opt/homebrew/bin:/usr/bin:/bin` excluding host `claude` CLI in sandbox), and `spacedock status --workflow-dir docs/dev --validate` (`VALID`) all passed; code committed on `paseo-worker-adoption-and-recovery` (`b3f8eadd2`) and stage report committed path-scoped in `.spacedock-state`.

### Summary

Implemented `paseo-worker-adoption-and-recovery` (`b3f8eadd2` atop `e1b971b32`, +358 net LOC across 5 files) with zero edits to 936's test files, zero edits to `first-officer-shared-core.md`, and zero uncalled production Go simulators. All offline and `-race` test suites, `-tags live` compilation/vet, `contractlint` guards, and workflow validation passed. Named F8 sandbox unavailability observed inside this least-privilege ensign (`paseo status --json` fails with `DAEMON_STATE_READ_FAILED` / `EPERM`); Commander runs `TestLivePaseoAdoptionSmoke` directly outside the ensign sandbox.


## Stage Report: implementation (cycle 2)

- DONE: Fix Material Finding 1 in internal/ensigncycle/paseo_adoption_live_test.go (assert non-main non-HEAD code branch, committed receipt.txt in latest code commit, and post-pre_dispatch_state_rev entity-only state commit adding ## Stage Report: implementation, - DONE:, and ### Summary)
  Commit `01d66ed16` tightens AC-1 completion assertions in `TestLivePaseoAdoptionSmoke`: checks `rev-parse --abbrev-ref HEAD` is non-empty/non-`main`/non-`HEAD`, `git log -1 --name-only` includes `receipt.txt`, `git show HEAD:receipt.txt` contains `paseoAdoptionMarker`, and a post-`pre_dispatch_state_rev` commit in `.spacedock-state` touches ONLY the entity file while its `-U0` diff adds `## Stage Report: implementation`, `- DONE:`, and `### Summary` (verified against standalone Git fixtures for cases a–e; detached HEAD, ignored uncommitted `receipt.txt`, stale/bare report, or multi-file state commit each fail).
- DONE: Fix Material Finding 2 in internal/ensigncycle/paseo_adoption_live_test.go (bind rec.WorkspaceID, rec.ServerID, rec.ReqSpec.{Provider,Mode,Thinking}, and rec.RetSpec.{Provider,Mode,Thinking} against actual daemon/inspect/requested values)
  Commit `01d66ed16` strengthens AC-1 worker-metadata assertions in `TestLivePaseoAdoptionSmoke`: decodes `paseo inspect` (`Provider`, `Model`, `Mode`, `Thinking`), `paseo ls` (`WorkspaceID`), and `paseo status --json` (`ServerID`, falling back to `"local"`), updates `fo1Prompt` to record those exact fields per `paseo-fo-dispatch.md`, and asserts exact equality for `workspace_id`, `server_id`, `requested.{provider,mode,thinking}` (`"full-access"` for `codex/*`, else `""`), and `returned.{provider,mode,thinking}` against observed `W1` (wrong provider, missing mode/thinking, or fabricated returned fields fail; inline/block YAML and retry parsers verified via `go run` + `-tags live` compile/vet).
- DONE: Keep cumulative diff vs e1b971b32 <= +365 net LOC across the same 5 files (zero edits to 936 test files); run gofmt, go test ./..., go test ./... -race, go test -tags live -c, go vet -tags live, and spacedock status --validate; commit on paseo-worker-adoption-and-recovery and append updated implementation Stage Report
  `git diff --shortstat e1b971b32..HEAD` is `5 files changed, 377 insertions(+), 14 deletions(-)` (+363 net LOC, <= +365, zero edits to 936 files or `first-officer-shared-core.md`, `.safehouse` untracked/gitignored); `gofmt -l`, `go vet ./...`, `go vet -tags live ./internal/ensigncycle`, `go test -tags live -c`, AC-2 (`go test -race -count=1 ./internal/contractlint` — fails if classifier/closure/caps break), AC-3 (`go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$'` — fails if H2 `## Paseo Worker` contaminates checklist/AC spans), full `go test ./...` and `go test ./... -race` (all 23 packages, `internal/cli` with `PATH=/opt/homebrew/bin:/usr/bin:/bin`), and `spacedock status --workflow-dir docs/dev --validate` (`VALID`) all passed.

### Summary

Applied only the 2 FO-authorized Cycle 1 Material fixes to `internal/ensigncycle/paseo_adoption_live_test.go` in commit `01d66ed16` on `paseo-worker-adoption-and-recovery`, keeping cumulative diff vs `e1b971b32` at `+363` net LOC across the same 5 files. All offline, `-race`, `-tags live` compile/vet, AC-2 `contractlint`, AC-3 `status` compatibility, adversarial Git/YAML fixture checks, and workflow validation checks passed; live execution of `TestLivePaseoAdoptionSmoke` remains outside this ordinary Safehouse ensign per F8.


## Stage Report: validation

- DONE: Verify AC-1 through AC-3 and both Cycle 1 Material fixes in commit 01d66ed16 stacked atop e1b971b32 (tightened completion Git assertions catching stale/empty report, uncommitted receipt, detached HEAD, multi-file state commit; exact worker metadata binding for workspace_id, server_id, requested/returned provider/mode/thinking; +363 net LOC / 5 files <= +365 cap; zero edits to 936 test files or first-officer-shared-core.md)
  Verified `git diff --shortstat e1b971b32..01d66ed16` = `5 files changed, 377 insertions(+), 14 deletions(-)` (+363 net LOC <= +365 cap, zero edits to `936` test files or `first-officer-shared-core.md`); both Cycle 1 Material findings (`0693a5c1`) are closed in `internal/ensigncycle/paseo_adoption_live_test.go` for AC-1 (lines 151-178 enforce non-main non-HEAD code branch, exact `receipt.txt` in `log -1 --name-only` + `HEAD:receipt.txt` marker, and post-`pre_dispatch_state_rev` entity-only `--follow` state commit adding `## Stage Report: implementation`, `- DONE:`, `### Summary`; lines 48-57, 65-70, 95-99, 114-118, 136-147 bind `workspace_id`, `server_id` against `paseo status --json` or `"local"`, and exact `requested`/`returned` `{provider,mode,thinking}` against Paseo `inspect.ts` / `status.ts` wire schemas), while `fo-write-core.md`, `paseo-fo-dispatch.md`, `paseo_worker_compat_test.go`, and `docs/runtime-live-ci.md` match the approved spec for AC-2 and AC-3.
- DONE: Run gofmt -l, go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$', go test -race -count=1 ./internal/contractlint, full go test ./... + -race, go test -tags live -c + go vet -tags live ./internal/ensigncycle, and spacedock status --workflow-dir docs/dev --validate; perform semantic adversarial pass
  AC-3 proved by `go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$'` (8/8 subtests across `running|archived|superseded|held` × `before-report|after-report`, failing if `## Paseo Worker` is demoted to H3 or contaminates report/checklist/AC spans); AC-2 proved by `go test -race -count=1 ./internal/contractlint` (asserts `FO-WRITE-CLASSIFIER` targets, deferred reference closure, FO byte caps, and model SSOT — breaking classifier rules, citing ensign `references/*.md`, hardcoding models, or exceeding caps fails); AC-1 offline/wire proofs passed via `go test -tags live -c`, `go vet -tags live ./internal/ensigncycle`, full `go test -count=1 ./...` and `go test -count=1 -race ./...` (all 24 packages; `internal/cli` with `PATH=/opt/homebrew/bin:/usr/bin:/bin` passed in 203.8s / 200.5s), `go vet ./...`, `gofmt -l` (clean on both x9m Go files; only untouched pre-existing `internal/release/runtime_live_evidence_workflow_test.go` from `e68aa5b7a` reported), `spacedock status --workflow-dir docs/dev --validate` (`VALID`), and a semantic adversarial pass (standalone Go matrix proving stale reports, multi-file state commits, bare headings, YAML section spillover fail while archive-followed entity-only commits pass).
- DONE: Record ## Stage Report: validation (Recommendation: PASSED or REJECTED) citing AC-1, AC-2, and AC-3 explicitly in - DONE: lines with falsifying edits, Review-finding disposition, and named F8 sandbox unavailability for live smoke
  Recommendation: PASSED. AC-1, AC-2, and AC-3 are evidenced above with falsifying edits; F8 named sandbox unavailability confirmed inside this least-privilege ensign (`paseo status --json` fails with `DAEMON_STATE_READ_FAILED` / `EPERM` on `paseo.pid`, so live `TestLivePaseoAdoptionSmoke` execution for AC-1 is owned concurrently by Commander on `01d66ed16`). Review-finding disposition: both Cycle 1 Material findings (`0693a5c1`) verified closed in `01d66ed16`; zero new Material findings and zero Deferred risks (Polish only: `paseo_worker_compat_test.go` asserts `Unevidenced == "false"` rather than `len(Citations) == 1` on `held/after-report`, though H2 span/checklist bounds already prove isolation; `paseo_adoption_live_test.go` checks `- DONE:` / `### Summary` headings in `-U0` diff while FO-2 enforces checklist completeness via `HasValidNewReport`).

### Summary

Recommendation: **PASSED** on commit `01d66ed16` (`e1b971b32..01d66ed16`: +363 net LOC across 5 files, zero edits to `936` test files or `first-officer-shared-core.md`, code worktree clean). Both Cycle 1 Material findings from detached audit `0693a5c1` are closed cleanly, all offline and `-race` Go suites (24 packages), `contractlint`, `TestPaseoWorkerStatusCompatibility`, `-tags live` compile/vet, workflow validation, and adversarial YAML/Git/wire-schema checks passed with zero new Material findings or Deferred risks. Live `TestLivePaseoAdoptionSmoke` execution is blocked inside this ordinary Safehouse ensign by F8 (`DAEMON_STATE_READ_FAILED` / `EPERM`) and is paired with Commander's concurrent live run on `01d66ed16`.


## Stage Report: implementation (cycle 3)

- DONE: Fix Ship-Blocker 1 in skills/first-officer/references/paseo-fo-dispatch.md (make Recovery Matrix sole decision owner, align Verb 3 error routing, restrict Row 10 to already-processed terminal events, tolerate AGENT_NOT_FOUND on Rows 4 & 6, keep Row 2 permission precedence, and consolidate redundant Unrelated column / repeated binding prose)
  Commit `dafceaea6` makes `## Recovery Matrix and Retry Ledger` the sole decision owner (Verb 3 now inspects/evaluates `HasValidNewReport` and routes `timeout`/`permission`/`idle`/`error`/`AGENT_NOT_FOUND` through the matrix, removing the unconditional `error` hold), restricts Row 10 strictly to duplicate/late completion events for already-processed terminal attempts (removing `"repeated wait"`), treats both `AGENT_ALREADY_ARCHIVED` and `AGENT_NOT_FOUND` as idempotent success on Rows 4 & 6, makes Row 2 permission precedence over Row 3 explicit, removes the redundant `Unrelated` column (stating per-entity isolation once above the table), and consolidates repeated spawn-chain/10-key/cascade prose in `## Runtime implementation` (reducing `paseo-fo-dispatch.md` from 20,316 to 19,047 bytes with zero ensign `references/*` paths and all 8 `«...»` bindings intact); AC-2 proved by `go test -race -count=1 ./internal/contractlint` (asserts `FO-WRITE-CLASSIFIER` targets, deferred reference closure, FO byte caps, and model SSOT — breaking classifier rules, citing ensign `references/*.md`, or breaching caps fails).
- DONE: Fix Ship-Blocker 2 in internal/ensigncycle/paseo_adoption_live_test.go (trim duplicated schema/adoption recipe from fo1Prompt/fo2Prompt so installed binding owns it; observe FO-2 tool call inspecting/waiting on w1ID via paseo logs --filter tools --json while W1 is running before releasing .adoption-release and fail if never observed; exclude only fo1/fo2 IDs in allAgents roster check)
  Commit `dafceaea6` strips duplicated `## Paseo Worker` schema and adoption recipes from `fo1Prompt`/`fo2Prompt` (retaining only fixture paths, worker provider/thinking/mode flags, and the `.adoption-release` barrier so the installed `first-officer` skill + `paseo-fo-dispatch.md` drive persistence and adoption), replaces the `UpdatedAt`/silent-expiry loop with `waitForPaseoAdoptionCall` polling `paseo logs <fo2ID> --filter tools --json` against the verified Paseo `logs.ts` / `agent-types.ts` schema (`{agentId, items: [{type:"tool_call", detail:{command, input, output}}]}`) and failing via `t.Fatalf` if FO-2 exits or times out (90s) without issuing `paseo inspect`/`paseo wait` referencing `w1ID` while `W1.Status == "running"`, and narrows `allAgents` exclusions strictly to `fo1.AgentID` and `fo2.AgentID`; AC-1 offline/wire proofs passed via `go test -tags live -c`, `go vet -tags live ./internal/ensigncycle`, and standalone Go parser matrices (11/11 `paseo logs --filter tools --json` negative/positive cases + 4/4 YAML/retry parser cases — failing if reading `index.md`, inspecting `fo2ID`, or running `paseo ls` falsely triggers release, or if a valid `inspect`/`wait` on `w1ID` is missed).
- DONE: Keep cumulative diff vs e1b971b32 <= +365 net LOC across the same 5 files (zero edits to 936 test files or first-officer-shared-core.md); run gofmt, go test ./..., go test ./... -race, go test -tags live -c, go vet -tags live, and spacedock status --validate; commit on paseo-worker-adoption-and-recovery and append implementation cycle 3 Stage Report
  `git diff --shortstat e1b971b32..dafceaea6` is `5 files changed, 413 insertions(+), 57 deletions(-)` (+356 net LOC <= +365 cap across the exact same 5 files; zero edits to `936` test files or `first-officer-shared-core.md`; `.safehouse` untracked/gitignored); `gofmt -l` clean on both `x9m` Go files; `go vet ./...`, `go vet -tags live ./internal/ensigncycle`, `go test -tags live -c`, AC-2 (`go test -race -count=1 ./internal/contractlint`), AC-3 (`go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$'` — 8/8 subtests across `running|archived|superseded|held` × `before-report|after-report`, failing if `## Paseo Worker` is demoted to H3 or contaminates checklist/AC spans), `go test -race -count=1 ./internal/status ./internal/contractlint`, full `go test -count=1 ./...` and `go test -race -count=1 ./...` across all 24 Go packages (`internal/cli` with `PATH=/opt/homebrew/bin:/usr/bin:/bin`, 184.8s / 189.7s), and `spacedock status --workflow-dir docs/dev --validate` (`VALID`) all passed; code committed on `paseo-worker-adoption-and-recovery` (`dafceaea6`) and report committed path-scoped in `.spacedock-state`.

### Summary

Closed both pre-merge audit Ship-Blockers (`9e71ce2c3`) in commit `dafceaea6` on `paseo-worker-adoption-and-recovery`, reducing cumulative diff vs `e1b971b32` from `+363` to `+356` net LOC across the exact same 5 files with zero edits to `936` test files or `first-officer-shared-core.md`. All offline, `-race` (24/24 packages), AC-2 `contractlint`, AC-3 `status` compatibility, `-tags live` compile/vet, `paseo logs --filter tools --json` wire-schema, and workflow validation checks passed. Live execution of `TestLivePaseoAdoptionSmoke` remains blocked inside this ordinary Safehouse ensign by F8 (`DAEMON_STATE_READ_FAILED` / `EPERM` on `paseo.pid`) for Commander to run outside the sandbox.


## Stage Report: validation

- DONE: Verify Ship-Blocker 1 closure in skills/first-officer/references/paseo-fo-dispatch.md (Recovery Matrix sole decision owner, Verb 3 routes through matrix with no unconditional error hold, Row 10 restricted to already-processed terminal attempts, Rows 4 & 6 tolerate AGENT_NOT_FOUND + AGENT_ALREADY_ARCHIVED, Row 2 permission precedence explicit, redundant Unrelated column & repeated binding prose consolidated, zero ensign references/* paths)
  Verified commit `dafceaea6` against `pre-merge-audit.md` (`9e71ce2c3`) and Paseo 0.17.1 CLI source (`archive.ts`, `wait.ts`, `inspect.ts`, `ls.ts`, `stop.ts`): `## Recovery Matrix and Retry Ledger` (line 109) is sole decision owner; Verb 3 (line 55) routes `timeout`/`permission`/`idle`/`error`/`AGENT_NOT_FOUND` through the matrix with zero unconditional error hold; Row 10 (line 126) scopes no-op strictly to already-processed terminal `(agent_id, epoch)` with no pending advance (dropping `"repeated wait"` so timeout->idle and Step 1 post-archive resume work); Rows 4–6 (lines 120–122) all treat `AGENT_ALREADY_ARCHIVED` and `AGENT_NOT_FOUND` as idempotent success; Row 2 (line 118) explicitly precedes Row 3 (`PendingPermissions=[]`, line 119); redundant `Unrelated` column removed (isolation stated once at line 109) and `## Runtime implementation` consolidated (20,316 -> 19,047 bytes, all 8 `«...»` bindings matching `fo-dispatch-core.md`, zero `references/` paths); AC-2 proved by `go test -race -count=1 ./internal/contractlint` (breaking `FO-WRITE-CLASSIFIER`, citing ensign `references/*.md`, hardcoding models, or exceeding FO byte caps fails).
- DONE: Verify Ship-Blocker 2 closure in internal/ensigncycle/paseo_adoption_live_test.go (fo1Prompt/fo2Prompt stripped of duplicated schema/adoption recipe so installed binding owns it; waitForPaseoAdoptionCall observes FO-2 inspect/wait on w1ID via paseo logs --filter tools --json while W1 is running before releasing .adoption-release and fails if never observed; allAgents excludes only fo1/fo2 IDs) and verify +356 net LOC / 5 files <= +365 cap with zero edits to 936 test files or first-officer-shared-core.md
  Verified `git diff --shortstat e1b971b32..dafceaea6` = `5 files changed, 413 insertions(+), 57 deletions(-)` (+356 net LOC <= +365 cap across the exact 5 files; zero edits to `936` test files or `first-officer-shared-core.md` at 23,497/23,500 B); in `paseo_adoption_live_test.go` for AC-1, `fo1Prompt`/`fo2Prompt` (lines 57–58, 79–80) drop all duplicated schema/adoption recipes so installed skills own persistence and adoption, `waitForPaseoAdoptionCall` (lines 194–227) polls `paseo logs <fo2ID> --filter tools --json` against verified Paseo `logs.ts`/`agent-types.ts` wire schemas and fails via `t.Fatalf` if FO-2 exits or times out (90s) without `inspect`/`wait` on `w1ID` while `W1.Status == "running"` before writing `.adoption-release`, and `allAgents` (lines 97–103) excludes strictly `fo1.AgentID`, `fo2.AgentID`, and `w1ID`; standalone adversarial Go matrices passed 14/14 `logs --filter tools --json` cases and 5/5 YAML/retry parser cases (reading `index.md`, `paseo ls`, `inspect fo2ID`, or YAML spillover each fail to trigger/match).
- DONE: Run gofmt -l, go test -race -count=1 ./internal/status ./internal/contractlint, full go test ./... + -race, go test -tags live -c + go vet -tags live ./internal/ensigncycle, and spacedock status --workflow-dir docs/dev --validate; record ## Stage Report: validation (Recommendation: PASSED or REJECTED) citing AC-1, AC-2, AC-3 and F8
  Recommendation: PASSED. AC-3 proved by `go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$'` (8/8 subtests across `running|archived|superseded|held` × `before-report|after-report`, failing if `## Paseo Worker` is demoted to H3 or contaminates report/checklist/AC spans) and full `go test -race -count=1 ./internal/status ./internal/contractlint` (79.4s / 4.4s); AC-2 proved by `go test -race -count=1 ./internal/contractlint` (all classifier/closure/cap/SSOT tests PASS); AC-1 offline/wire proofs passed via `go test -tags live -c`, `go vet -tags live ./internal/ensigncycle`, `go vet ./...`, `git diff --check e1b971b32..HEAD`, `gofmt -l ./cmd ./internal` (clean on both x9m Go files; only untouched pre-existing `internal/release/runtime_live_evidence_workflow_test.go` from `e68aa5b7a`), full non-race `go test -count=1 ./...` and `-race` `go test -count=1 -race ./...` across all 24 Go packages (`internal/cli` with `PATH=/opt/homebrew/bin:/usr/bin:/bin`: 183.7s / 204.6s; `internal/ensigncycle`: 238.2s / 271.3s; `internal/piruntime` with default PATH: 0.19s / 1.65s), and `spacedock status --workflow-dir docs/dev --validate` (`VALID`). F8 named sandbox unavailability confirmed inside this least-privilege ensign (`paseo status --json` fails with `DAEMON_STATE_READ_FAILED` / `EPERM` on `paseo.pid`; Commander runs live `TestLivePaseoAdoptionSmoke` concurrently on `dafceaea6`). Review-finding disposition: both pre-merge audit Ship-Blockers (`9e71ce2c3`) verified closed in `dafceaea6`; zero new Material findings and zero Deferred risks (Polish only: `waitForPaseoAdoptionCall` matches canonical `paseo inspect`/`paseo wait` subcommand order per `paseo-fo-dispatch.md` and fails closed on non-canonical flag placement before subcommand).

### Summary

Recommendation: **PASSED** on commit `dafceaea6` (`e1b971b32..dafceaea6`: +356 net LOC across 5 files, zero edits to `936` test files or `first-officer-shared-core.md`, code worktree clean). Both pre-merge audit Ship-Blockers (`9e71ce2c3`) are closed cleanly with zero new Material findings or Deferred risks, and all offline and `-race` Go suites (24/24 packages), AC-2 `contractlint`, AC-3 `TestPaseoWorkerStatusCompatibility`, `-tags live` compile/vet, Paseo 0.17.1 wire-schema/adversarial parser checks, and workflow validation passed. Live `TestLivePaseoAdoptionSmoke` execution for AC-1 is unavailable inside this least-privilege Safehouse ensign per F8 (`DAEMON_STATE_READ_FAILED` / `EPERM`) and is paired with Commander's concurrent live run on `dafceaea6`.


## Stage Report: implementation (cycle 4)

- DONE: Update internal/ensigncycle/paseo_adoption_live_test.go:144-171 to verify committed receipt.txt and clean checkout in BOTH the unmerged worktree state AND the fo-merge-core.md merged-and-removed main state
  Commit `628949773` updates `TestLivePaseoAdoptionSmoke` (AC-1) to branch on `os.Stat(codeWorktree)`: when the worktree exists it keeps the exact non-main/non-HEAD branch, latest-commit `receipt.txt`, `HEAD:receipt.txt` marker, and worktree+state cleanliness checks; when `os.IsNotExist(err)` (FO-2 merged per `fo-merge-core.md` and removed the worktree) it checks `fx.RepoRoot` is on `main`, a non-initial commit reachable from `HEAD` touched `receipt.txt` (`git log --format=%P -- receipt.txt` non-empty), `HEAD:receipt.txt` contains `paseoAdoptionMarker`, and both `fx.RepoRoot` and `fx.StateRoot` are clean (non-ENOENT stat errors fail closed). Verified via a standalone 11/11 Git fixture matrix proving unmerged-worktree, `--no-ff` merged-`main` (with worktree removed), and `--ff-only` merged-`main` pass, while uncommitted/gitignored `receipt.txt`, dirty code repo, dirty state checkout, detached `HEAD`, initial-commit-only `receipt.txt`, or deleted `receipt.txt` each fail.
- DONE: Keep cumulative diff vs e1b971b32 <= +365 net LOC across the exact same 5 files (touch ONLY internal/ensigncycle/paseo_adoption_live_test.go); verify gofmt -l, go test -tags live -c, go vet -tags live, go test -race -count=1 ./internal/status ./internal/contractlint, and spacedock status --validate; commit cleanly and append implementation cycle 4 Stage Report citing AC-1, AC-2, AC-3
  Commit `628949773` on `paseo-worker-adoption-and-recovery` touches ONLY `internal/ensigncycle/paseo_adoption_live_test.go` (`+12/-5`, +7 net LOC); cumulative `git diff --shortstat e1b971b32..628949773` is `5 files changed, 420 insertions(+), 57 deletions(-)` (+363 net LOC <= +365 cap across the exact same 5 files; zero edits to `936` test files or `first-officer-shared-core.md`; `.safehouse` untracked/gitignored). Verified `gofmt -l ./internal/ensigncycle/paseo_adoption_live_test.go` (clean), `go test -tags live -c` + `go vet -tags live ./internal/ensigncycle` + `go test -race -count=1 ./internal/ensigncycle -run '^TestPaseo'` (AC-1 offline/compile/vet/control guards — syntax, vet, or Paseo live-control regressions fail), `go test -race -count=1 ./internal/contractlint` (AC-2 — breaking `FO-WRITE-CLASSIFIER`, citing ensign `references/*.md`, hardcoding models, or exceeding FO caps fails), `go test -race -count=1 ./internal/status` including `TestPaseoWorkerStatusCompatibility` 8/8 subtests (AC-3 — demoting `## Paseo Worker` to H3 or contaminating checklist/AC spans fails), `go vet ./...`, and `spacedock status --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --validate` (`VALID`).

### Summary

Applied the Captain-authorized Cycle 3 micro-fix solely in `internal/ensigncycle/paseo_adoption_live_test.go` (commit `628949773`, +7 net LOC, cumulative `e1b971b32..628949773` = `+363` net LOC across the exact same 5 files) so `TestLivePaseoAdoptionSmoke` accepts both unmerged worktree and `fo-merge-core.md` merged-and-removed `main` post-completion states while still rejecting uncommitted/ignored receipts, initial-only receipts, wrong branches, or dirty checkouts. All required formatting, `-tags live` compile/vet, `-race` `status`/`contractlint`/`TestPaseo`, standalone 11-case Git fixture matrix, and workflow validation checks passed; live execution of `TestLivePaseoAdoptionSmoke` remains outside this ordinary ensign for Commander to re-run.


## Stage Report: validation

- DONE: Verify commit 628949773 (+12/-5 in internal/ensigncycle/paseo_adoption_live_test.go only; cumulative e1b971b32..628949773 is +363 net LOC / 5 files <= +365 cap) safely accepts both unmerged worktree and fo-merge-core.md merged-and-removed main states while still failing closed on uncommitted/ignored receipt.txt, initial-only receipt, wrong branch, non-ENOENT stat error, or dirty repo
  Verified `git diff --stat dafceaea6..628949773` touches ONLY `internal/ensigncycle/paseo_adoption_live_test.go` (`+12/-5`), cumulative `git diff --shortstat e1b971b32..628949773` is `5 files changed, 420 insertions(+), 57 deletions(-)` (+363 net LOC <= +365 cap, zero edits to `936` test files or `first-officer-shared-core.md`), and `git diff --check e1b971b32..628949773` is clean; for AC-1, lines 144-156 and 176 branch on `os.Stat(codeWorktree)` (failing closed via `else if err != nil` — verified in Go on macOS that `os.IsNotExist` is `false` for both `ENOTDIR` and `EACCES`), preserve canonical `codeWorktree` matching `rec.Cwd` on line 134 when the worktree is removed, and passed a 16/16 standalone adversarial Git matrix proving unmerged-worktree and merged-and-removed `--no-ff`/`--ff-only`/`--squash` `main` states pass while removed-without-merge, initial-only `receipt.txt`, uncommitted/gitignored `receipt.txt` (both states), wrong branch or detached `HEAD` (both states), `ENOTDIR`, deleted/wrong-marker `receipt.txt`, and dirty code/state checkouts each fail closed.
- DONE: Run gofmt -l, go test -tags live -c, go vet -tags live ./internal/ensigncycle, go test -race -count=1 ./internal/status ./internal/contractlint ./internal/ensigncycle, and spacedock status --workflow-dir docs/dev --validate; record ## Stage Report: validation (Recommendation: PASSED or REJECTED) citing AC-1, AC-2, AC-3 and F8
  Recommendation: PASSED. AC-1 offline/compile/vet/control proofs passed via `gofmt -l ./internal/ensigncycle/paseo_adoption_live_test.go` (clean), `go test -tags live -c -o /tmp/ensigncycle-x9m-val4.test ./internal/ensigncycle && rm -f /tmp/ensigncycle-x9m-val4.test`, `go vet -tags live ./internal/ensigncycle`, `go vet ./...`, `go test -race -count=1 ./internal/ensigncycle -run '^TestPaseo'` (2/2 Paseo control tests, 1.78s/2.27s — breaking live config/env/JSON/precondition guards fails), and `go test -count=1 ./internal/ensigncycle` (423.6s; scoped per Commander notes since `dafceaea6` passed full 24-package normal + `-race` and `628949773` edits only `//go:build live` `paseo_adoption_live_test.go`); AC-2 proved by `go test -race -count=1 ./internal/contractlint` (4.61s; all 8 named `FO-WRITE-CLASSIFIER`, reference-closure, FO-cap, and model-SSOT guards PASS — breaking classifier patterns, citing ensign `references/*.md`, hardcoding models, or breaching caps fails); AC-3 proved by `go test -race -count=1 ./internal/status` (93.69s, including `TestPaseoWorkerStatusCompatibility` 8/8 subtests in 4.05s across `running|archived|superseded|held` × `before-report|after-report` — demoting `## Paseo Worker` to H3 or contaminating checklist/AC spans fails); `spacedock status --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --validate` returned `VALID`. F8 named sandbox unavailability confirmed inside this least-privilege ensign (`paseo status --json` fails with `DAEMON_STATE_READ_FAILED` / `EPERM` on `paseo.pid`; Commander is running live `TestLivePaseoAdoptionSmoke` concurrently on `628949773`). Review-finding disposition: zero new Material findings, zero Deferred risks, and zero Polish blockers.

### Summary

Recommendation: **PASSED** on commit `628949773` (`dafceaea6..628949773`: `+12/-5` in `internal/ensigncycle/paseo_adoption_live_test.go` only; cumulative `e1b971b32..628949773`: `+363` net LOC across 5 files <= `+365` cap; code worktree untouched and clean). The post-merge assertion fix cleanly accepts both unmerged worktree and `fo-merge-core.md` merged-and-removed `main` states while failing closed across all 16 adversarial Git/stat cases, and all required `gofmt`, `-tags live` compile/vet, `-race` `status`/`contractlint`/`TestPaseo`, normal `ensigncycle`, and workflow validation checks passed with zero Material findings or Deferred risks. Live `TestLivePaseoAdoptionSmoke` execution for AC-1 remains unavailable inside this Safehouse ensign per F8 (`DAEMON_STATE_READ_FAILED` / `EPERM`) and is paired with Commander's concurrent live run on `628949773`.


## Stage Report: validation

- DONE: Verify Captain-authorized 1-line edit in commit ca39881e2 (628949773..ca39881e2: internal/ensigncycle/paseo_adoption_live_test.go:136 accepting both w1.Provider+"/"+w1.Model and w1.Provider per x9m/index.md:112) and confirm cumulative diff e1b971b32..ca39881e2 remains +363 net LOC / 5 files (<= +365 cap) with zero edits to 936 test files or first-officer-shared-core.md
  Verified `git diff 628949773..ca39881e2` is a single `+1/-1` (0 net LOC) line change at `internal/ensigncycle/paseo_adoption_live_test.go:136` changing `rec.RetSpec.Provider != wantRetProvider` to `(rec.RetSpec.Provider != wantRetProvider && rec.RetSpec.Provider != w1.Provider)`, matching `x9m/index.md` (`returned.provider`: `<Provider>/<Model>` or `<Provider>`, now line 131 after frontmatter gate records) while still rejecting empty `wantRetProvider`, wrong provider, wrong model, or trailing slash; `git diff --stat e1b971b32..ca39881e2` is `5 files changed, 420 insertions(+), 57 deletions(-)` (+363 net LOC <= +365 cap across the exact 5 files, zero edits to `936` test files or `first-officer-shared-core.md`, `git status` and `git diff --check` clean).
- DONE: Run small-change fast-path validation checks (gofmt -l ./internal/ensigncycle/paseo_adoption_live_test.go, go test -tags live -c, go vet -tags live ./internal/ensigncycle, go test -race -count=1 ./internal/ensigncycle -run '^TestPaseo', go test -race -count=1 ./internal/contractlint, go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$', and spacedock status --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --validate)
  All 7 small-change fast-path checks passed on `ca39881e2`: `gofmt -l ./internal/ensigncycle/paseo_adoption_live_test.go` (clean), `go test -tags live -c -o /tmp/ensigncycle-x9m-val5.test ./internal/ensigncycle && rm -f /tmp/ensigncycle-x9m-val5.test` (exit 0), `go vet -tags live ./internal/ensigncycle` (exit 0), `go test -race -count=1 ./internal/ensigncycle -run '^TestPaseo'` (2/2 PASS in 1.73s — breaking live config/env/run-arg/JSON precondition guards fails), `go test -race -count=1 ./internal/contractlint` (PASS in 3.23s — AC-2: breaking `FO-WRITE-CLASSIFIER` targets, citing ensign `references/*.md`, hardcoding models, or exceeding FO caps fails), `go test -race -count=1 ./internal/status -run '^TestPaseoWorkerStatusCompatibility$'` (8/8 subtests PASS in 3.26s across `running|archived|superseded|held` × `before-report|after-report` — AC-3: demoting `## Paseo Worker` to H3 or contaminating checklist/AC spans fails), and `/tmp/spacedock-936 status --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --validate` (`VALID`).
- DONE: Record ## Stage Report: validation (Recommendation: PASSED or REJECTED) for commit ca39881e2 citing AC-1, AC-2, AC-3 with falsifying edits, noting Commander's live TestLivePaseoAdoptionSmoke PASS on ca39881e2 (344.08s) and F8 ensign sandbox boundary
  Recommendation: PASSED. AC-1 is proved end-to-end by Commander's strict live `TestLivePaseoAdoptionSmoke` PASS on `ca39881e2` (344.08s: FO-2 observed via `paseo logs --filter tools --json`, adopted `W1` with 0 duplicate spawns and 0 retries, merged `receipt.txt` to `main`, cleaned worktree, archived `W1`, advanced to `done`/`PASSED` — spawning a second worker/retry, recording a wrong `returned.provider`, failing to advance, or leaving `W1` unarchived fails) alongside offline compile/vet/control checks above; AC-2 and AC-3 are proved above with falsifying edits. F8 ensign sandbox boundary confirmed directly (`paseo status --json` fails with `DAEMON_STATE_READ_FAILED` / `EPERM` on `paseo.pid` inside this least-privilege ensign). Review-finding disposition: zero Material findings and zero Deferred risks (Polish only: if a daemon hypothetically returned `Provider == ""` with non-empty `Model`, `wantRetProvider` would become `"/"+Model` and `rec.RetSpec.Provider == ""` would match `w1.Provider`, though Paseo always populates `Provider` on real agents).

### Summary

Recommendation: **PASSED** on commit `ca39881e2` (`628949773..ca39881e2`: Captain-authorized `+1/-1` edit at `internal/ensigncycle/paseo_adoption_live_test.go:136` matching `x9m/index.md` `returned.provider` contract; cumulative `e1b971b32..ca39881e2`: `+363` net LOC across 5 files <= `+365` cap, zero edits to `936` test files or `first-officer-shared-core.md`, code worktree untouched and clean). All 7 small-change fast-path checks (`gofmt`, `-tags live` compile/vet, `-race` `TestPaseo`, `contractlint`, `TestPaseoWorkerStatusCompatibility`, and `spacedock status --validate`) passed with zero Material findings or Deferred risks, paired with Commander's strict live `TestLivePaseoAdoptionSmoke` PASS on `ca39881e2` (344.08s) outside this ensign's F8 sandbox boundary (`DAEMON_STATE_READ_FAILED` / `EPERM`).
