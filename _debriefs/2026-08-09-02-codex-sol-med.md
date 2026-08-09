---
session-date: 2026-08-09
sequence: 2
harness: Codex
model: sol
model-version-build: med
first-commit: 1947aacb0d7c3481c18a846f3566645fd2cb89ee
last-commit: a929fcb60dab0deda4bdf1768af8dfc9d66fce8f
duration: ~5d2h12m
---

# Sprint Debrief — `live-test-truth`

This is the full `live-test-truth` sprint debrief. It excludes concurrent work
unless this sprint directly required it. The sprint delivered a desired-state
registry, portable common journeys, truthful lane selection, and executable
reconciliation. Product failures remain visible as exact TODO cells with active
owners in `test-behavior-completeness`.

## Shipped

### Sprint members

- **v7** `refresh-v1-pilot-manifest-after-archive-moves` — [#612](https://github.com/spacedock-dev/spacedock/pull/612). Restored current-checkout validation after pilot tasks moved to the archive.
- **nz6** `land-live-test-truth-planning-package` — [#614](https://github.com/spacedock-dev/spacedock/pull/614). Landed the sprint plan separately from implementation work.
- **3d** `make-live-test-results-truthful` — [#615](https://github.com/spacedock-dev/spacedock/pull/615). Made live results report truthful runtime outcomes.
- **15e** `make-live-lanes-buy-named-evidence` — [#626](https://github.com/spacedock-dev/spacedock/pull/626). Bound live-lane cost to named evidence.
- **ys** `deliver-portable-common-live-journey-surface` — [#633](https://github.com/spacedock-dev/spacedock/pull/633). Delivered one portable common journey surface across runtimes.

### Closure and direct follow-ons

- **0y** `make-sonnet-5-only-claude-pr-live-lane` — [#639](https://github.com/spacedock-dev/spacedock/pull/639). Kept Sonnet 5 as the pull-request Claude lane and moved Opus to pre-release evidence.
- **824** `align-claude-break-glass-dispatch-oracle` — [#637](https://github.com/spacedock-dev/spacedock/pull/637). Aligned the Claude break-glass oracle with the selected dispatch mode.
- **3z** `sonnet-gate-guardrail-no-authority` — [#640](https://github.com/spacedock-dev/spacedock/pull/640). Restored the Sonnet gate-guardrail proof without destructive room-state recovery.
- **qz** `keep-journey-delta-green-on-rerun` — [#643](https://github.com/spacedock-dev/spacedock/pull/643). Kept optional journey metrics green across failed-job reruns.
- **26n** `headless-recorded-gate-stop-stage-coherence` — [#583](https://github.com/spacedock-dev/spacedock/pull/583). Made the headless gate fixture and durable oracle coherent.
- **c1n** `run-common-live-journeys-two-at-a-time` — [#647](https://github.com/spacedock-dev/spacedock/pull/647). Runs at most two Claude common journeys concurrently.
- **Live-test close cleanup** — [#646](https://github.com/spacedock-dev/spacedock/pull/646). Reconciled stale registry and sprint references.
- **Live-test reconciliation close** — [#648](https://github.com/spacedock-dev/spacedock/pull/648). Rejects unclassified live tests and inactive TODO owners.
- **Final sprint record** — [#649](https://github.com/spacedock-dev/spacedock/pull/649). Records the final 17-journey reconciliation result and active TODO ownership.

## Filed (backlog)

- **0a** `restore-optional-manual-pi-common-live-ci` — Restore manual Pi evidence without making Pi a merge requirement.
- **98a** `codex-headless-implementation-worker-before-validation` — Repair the shared Sonnet and Codex worker-dispatch defect exposed by 26n.
- **zh** `publish-rejection-round-before-regate` — Replace archived rejected owner `zbc` for the rejection-flow gap.
- **xp6** `restore-live-evidence-after-completed-repairs` — Restore evidence whose original repair owners are complete.
- **5k** `cut-gate-guardrail-turn-and-tool-bloat` — Reduce gate-guardrail turns and tool calls. This task moved to the durable-decision lane.
- **5f6** `reject-stale-same-minor-launcher-before-fo-work` — Reject stale launchers that lack required capabilities. This task moved to the durable-decision lane.

## Non-PR commits (workflow-only)

State changes that do not belong to a code PR:

- `bf592d776` Archived the superseded broad restoration task.
- `78dc49ecc` Reconciled completed tasks and live TODO ownership.
- `1556b560a` Completed the archive moves without duplicate active entities.

Routine dispatch, gate, and state commits are omitted.

## Decisions

- The registry is the desired-state SSOT. It is not a gap ledger.
- The sprint was recarved into tasks that each delivered visible test-surface value.
- Source annotations own volatile function and fixture bindings.
- Common journeys apply to all runtimes unless an explicit exception exists.
- Runtime-specific proofs cover host substrate behavior only.
- Known product failures remain exact TODOs with active owners.
- Product repairs belong in `test-behavior-completeness`.
- Pull requests require Sonnet 5/max and Codex Luna/max evidence.
- Opus is a pre-release lane. Pi is an optional manual lane.
- Every task must land visible value. Pure component changes are banned.
- Tests that only test test infrastructure are banned.
- Stable CI does not depend on mutable workflow state.
- Sprint close and release preparation run the mutable TODO-owner join.
- Local subscription-backed live runs take priority over paid CI.
- 26n merged without another live run after both runtime defects were reproduced.
- c1n merged after its target Claude common step passed in 13m13s.
- A recorded gate attempt is the stopping point. Gate readiness is a calculated
  view, not a stored lifecycle state. Independent work can continue before the
  captain sees the gate.
- Delegated conn includes gate approval, CI approval, PR action, and merge when
  the captain grants those actions.
- Small, deterministic state and documentation corrections can land directly.
  Do not create a PR only to carry bookkeeping changes.

## Issues — Workflow

- The rejected `ys` candidate reached 43 files and approximately +2,750/-714
  lines across 21 cycles. It encoded runtime selection through 16-case host
  switches and added tests for the reconciler itself.
- The repair returned the entry point to ordinary Go functions and annotations.
- Net change size matters. Visible value does not justify a large support
  mechanism when a small, readable surface can deliver the same result.
- Several tasks initially separated mechanisms from visible value. The sprint required another carve.
- Archived tasks remained TODO owners after their implementation work completed.
- The new close join found and repaired these stale owner bindings.
- The c1n substrate step remained active after the target common suite passed.
- Codex failed c1n on an unrelated checklist-path retry. An existing backlog task owns this friction.
- The latest measured gate-guardrail journey used 22 assistant turns and 24
  tool calls. The v0.26 baseline used 11 turns and 11 tool calls. Task `5k`
  owns this durable-decision friction.
- The requested shallow-boot token comparison with v0.26, and the other token
  deltas, did not reach a durable report. Do not treat that audit as complete.

## Issues — Spacedock

- The installed First Officer and debrief skill paths disappeared after plugin rotation. The repository copy provided the debrief fallback.
- The state checker emits many warnings for legacy gate fields and lowercase verdicts. Existing tasks track this migration.
- A merged PR can report failure because a worktree still uses its local branch.
- Delegated conn did not consistently result in prompt CI approval or PR action.
- The ready-to-present versus presented gate distinction was not obvious during dispatch.

These items already have local task coverage. No new GitHub issue was filed.

## Observations

- The desired registry became useful when it could expose missing tests before code existed.
- The delivered registry has 17 common journeys and 68 runtime cells.
- Reconciliation derives 48 runnable cells and 20 exact TODO cells.
- The current TODO split is Opus 3, Sonnet 5, Codex 6, and Pi 6.
- The 20 TODO cells have these active owners:

  - `98a`: default-headless on Sonnet and Codex (2 cells).
  - `zh`: rejection-flow on Sonnet, Opus, Codex, and Pi (4 cells).
  - `9a`: smallest-sufficient-mechanism and keep-moving-posture on Sonnet,
    Codex, and Pi (6 cells).
  - `47g`: withdrawn-gate-recovery on Codex (1 cell).
  - `xp6`: gate-guardrail on Codex and Pi; recorded-gate-lifecycle on Opus;
    default-headless on Pi; and owned-conflict on Sonnet, Opus, and Pi
    (7 cells).
- TODOs are honest only when their owners remain active.
- The source classifier provides durable value without another scenario table.
- The mutable owner join belongs at sprint close and release preparation.
- Live tests exposed product behavior that deterministic checks did not reveal.
- Parallel Claude common journeys reduced the target step to 13m13s.
- The sprint is closed. Remaining gaps now belong to behavior completeness.

## Agent Testimonial

- Date: 2026-08-09
- Harness/runtime: Codex
- Model: sol
- Model version/build: med
- Session scale: 18 tasks touched, 34 recorded worker dispatches, and 14 PRs merged

Spacedock preserved decisions, task ownership, and evidence across a long
session and several context resets. Durable IDs and state records made recovery
possible. The main costs were ceremony, stale state, warning noise, and unclear
responsibility at approval boundaries. Without Spacedock, more decisions can be
lost.

A smaller state model keeps source facts that cannot be reconstructed. These
facts include the entity stage, gate attempts, decisions, evidence references,
PR identity, and an optional session claim. Status views calculate readiness,
eligibility, human waits, and the next permitted action. They do not store these
projections as more lifecycle state.

Clear action ownership maps each observed fact to one responsible actor. The
captain owns an open decision without delegated authority. The Commander owns
authorized CI approval and merge actions. The First Officer owns advancement
and dispatch. Workers own assigned changes, and validators own recommendations.
This mapping prevents a passive status from hiding an authorized action.

## What is next

- Finish **0a**, which is in implementation.
- Repair **98a** and **zh**.
- Run **xp6** after its required product repairs land.
- Continue **9a** and **47g** under `test-behavior-completeness`.
- Keep durable-decision entities outside the ownership of this session.
