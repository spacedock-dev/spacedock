# 0203 — Handoff: watch #365 CI → merge → advance j9 → debrief

Cold-boot prompt for the next session. Boot the first officer (`spacedock claude`, then
`status --boot --json`). Sprint = the entities matching `sprint: 0203-fo-efficiency`. Read
`friction-log.md` (this session's dogfood record, FR-1..FR-13), `index.md` (goal/DoD), and this file.

## 1. Watch PR #365 → merge on green → advance j9

**j9 (`lazy-teamcreate-shallow-boot`) is the keystone — validated PASSED over 3 cycles, now PR-pending as #365 → `main`.** The live-e2e environments are captain-approved and running.

1. **Watch #365** — `gh pr checks 365 --repo spacedock-dev/spacedock`. Already PASS: offline, build, install (×2), pi-live. Running: **claude-live (opus `CI-E2E-OPUS` + sonnet `CI-E2E`)** and **codex-live** (~20–30 min on opus). Poll with `gh run watch 27485285428` or re-check `gh pr checks`.
2. **All-green → merge:** `gh pr merge 365 --repo spacedock-dev/spacedock --squash` (base `main`).
3. **Lane FAIL → triage first.** The **Codex `collab:wait` hang is a KNOWN live-infra flake** (3/5 runs this session → seeded as task `czza` / `codex-collab-wait-subagent-hang`) — **re-run that lane, do not treat the hang as a real failure.** A genuine assertion failure is a different matter: investigate before merging (do NOT merge red).
4. **After merge, advance j9** — the `pr-merge` hook (event-loop / idle / startup) sees `MERGED` and runs: two-step `mod-block=` clear → terminalize (`status --set lazy-teamcreate-shallow-boot completed verdict=PASSED worktree=`) → `status --archive` → `git worktree remove .worktrees/spacedock-ensign-lazy-teamcreate-shallow-boot` + `git branch -d spacedock-ensign/lazy-teamcreate-shallow-boot`. (A fresh session has a fresh team — the prior j9 impl-ensign + validator are already gone; just do the state + worktree cleanup. The remote branch was cleaned by the PR merge.)

## 2. Debrief

Invoke `spacedock:debrief`. Fold in:

- **j9 shipped (the keystone):** FO contract split (slim boot-resident core + lazily-loaded `claude-fo-dispatch.md` / `claude-fo-merge.md`) + lazy-TeamCreate + shallow-boot-then-greet off `status --boot --json` (with the S7b before-greet merged-PR sweep). **Measured live: greet ~43k context (vs ~160k peak deep-boot), zero TeamCreate, no ~89k cache spike.** 3 validation cycles, each a *real* catch: (1) a pre-existing single-entity bare-mode reviewer-reuse bug P2 *unmasked* (spun off as `js`); (2) the AC-3 fix's shared-prompt change broke Codex host-neutrality (fixed host-conditionally); (3) AC-2's no-TeamCreate oracle was a hollow multi-delta-parser false-pass (fixed: `ParseClaudeTurns` now merges later-delta tool_use). Rebased onto `origin/main` to fold #362 (`spacedock new` teaching) + #364 — avoiding a silent clobber the FO caught from the captain's "why next-id not new" question.
- **Dogfood friction log `friction-log.md` (FR-1..FR-13) — UNCOMMITTED on `main`; COMMIT it in the debrief.** Boot levers FR-1/2/3 (lazy-team ~89k / contract ~16k / README ~7.7k, all proven by j9); FR-4 entity-body cycle-report weight; FR-6 reconcile destructive Class-E remedy (→ `sr`); FR-8 per-transition state commit overhead; FR-9 the #344 budget-probe bug spuriously blocking the FO's own reuse (proves #344 load-bearing); FR-10 `spacedock new` untaught + the #362 stale-base catch; **FR-11 the #1 run-cost lever — dispatched subagents pinned to the 5m cache while top-level threads get 1h (~$14–73/run re-cache storm), a host-support ask**; FR-12 the RTK-hook rewrite-echo + harness `task_reminder` prefix tax (RTK/harness config); FR-13 the fix→review shared-worktree race.
- **Reorientation after 0.20.2 landed:** formed the post-flip stale-trunk cluster `sr` + `87` (both tagged into 0203); filed `rz` (`simplify-dev-readme-and-templates`, 0203), `xf` (`polish-on-demand-drop-standing`, 0203), `js` (`single-entity-reviewer-continuity`, e3z-family), `czza` (`codex-collab-wait-subagent-hang`).
- **Open captain items:** host-support asks (FR-11 cache-TTL-1h for subagents; FR-12 RTK-hook context-silence + `task_reminder` suppression — raise with Claude Code / RTK); merge-target precedent is **`main`** post-flip (the `pr-merge` mod's `next`-default is stale → `87`).

## 3. Next wave (off merged j9, on `main` — the contract-cleanup cluster)

Once j9 is on `main`, its slimmed refs are the trunk; the cluster unblocks:

- **T3** (`fo-contract-prose-audit`, ideation, staff-review READY) — audit the slimmed FO refs. Step-0 survey is the collapse fork (mechanical cut vs a recorded roadmap decision if nothing to cut). **Steer: KEEP the budget-probe reuse-condition-0 prose** (a deliberate cross-host split, not collapsible dup). Also absorb the 2 ship-local trunk-staleness lines now at `claude-fo-merge.md:31,34` (→ `main`).
- **`sr` + `87`** (stale-trunk cluster) — `sr` de-conflates `reconcile.go` git-hygiene from team-management + flips Class-D/E off `origin/next`; `87` refits the `pr-merge` mod base `next`→`main`. Settle **ONE** trunk-config source jointly. DoD: no helper/mod/ref/doc resolves the integration trunk to `next` post-flip.
- **`xf`** (`polish-on-demand-drop-standing`) — replace the standing-teammate lifecycle with on-demand one-shot polish; move usage prose into the `comm-officer` mod, leaving the contract a generic hook.
- **`#344`** (`context-budget-spurious-warnings`) — validated, held: **merge with the batch** (it's the FR-9 fix the FO's own reuse path needs).
- Fast-follows: `rz` (README+templates), `js` (reviewer continuity), `czza` (Codex hang triage).
- **Pre-cut antipattern audit before the `v0.20.3` tag:** independent staff-eng over the assembled sprint; confirm AC-6 measured the <60k/89k saving (it did: 43k). Then fire `v0.20.3`.

## State at handoff

| Entity | State |
|--------|-------|
| **j9** | PR #365 → `main`, live-e2e approved+running, **merge-on-green pending** |
| **#344** | validation, PASSED-held for the batch |
| **T3** | ideation (READY) — next-wave lead |
| **sr / 87 / xf** | backlog (0203 cluster) |
| **rz / js / czza** | backlog (fast-follow) |

Note: this session ran **as-if-implemented** — the FO operated as though j9 (lazy-TeamCreate + shallow-boot + split contract), #344 (clean probe), and T3 (slim refs) were already live. That dogfood is what surfaced FR-1..FR-13. The next session runs normally (j9 will actually be live once #365 merges).
