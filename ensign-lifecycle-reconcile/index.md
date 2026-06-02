---
id: qs87q0ca1wa3bhzfkj07mgwb
title: Harden FO reconciliation — teardown agents at terminal + supersede-shutdown + a state-keyed reconcile sweep (no lingering ensigns / stale branches / un-advanced PRs)
status: backlog
source: "captain (2026-06-02) — observed the FO leak 3 ensigns (the 7h implementer post-merge + the old 3g cycle agents post-rework) and reconcile in-flight branches only reactively at merge (02 hit #263 late); all keyed off drift-prone FO session memory"
completed:
verdict:
score: "0.30"
worktree:
issue:
---

The FO leaks agents and reconciles reactively because every lifecycle rule keys off **FO session memory**, which drifts under context pressure (this session: the FO tracked "5 alive" when 7 were — the 7h implementer lingered after its entity merged+archived, and the superseded old-3g `-implementation`/`-validation` lingered after 3g re-ideated). Same root cause as the prose-pin tests: a "remember to shut down / remember to rebase" rule can't self-check.

## Where it fails today (contract gaps)

1. **No teardown at terminal.** `first-officer-shared-core.md` `## Merge and Cleanup` step 9 removes the worktree + local branch but has NO step to shut down the entity's agents (implementer, validator, any `-N`/`-cycleN` variants, the detached auditor). → the 7h-class leak.
2. **No supersede-shutdown.** When an entity re-enters an earlier stage and a fresh / `-N`-suffixed agent supersedes a prior-cycle one (re-ideation, rework), nothing shuts down the superseded agent. → the old-3g-class leak.
3. **Reactive, memory-based reconciliation.** Merged-PR advancement, in-flight-branch freshness vs `next`, and agent liveness are all tracked in FO memory and reconciled only at merge time (02 hit the #263 conflict late — the pre-merge rebase caught it, but reactively).

## Proposed hardening (prefer a code gate over a prose rule)

- **`## Merge and Cleanup` — add a teardown step:** shut down every agent keyed to the entity (derive the set from the entity's worker-key/slug prefix over the LIVE roster, not memory), before/with worktree+branch removal.
- **`## Completion and Gates` / Feedback Rejection Flow — add a supersede-shutdown:** when a fresh dispatch supersedes a prior-cycle agent, shut the superseded one down at the moment of supersession.
- **A `spacedock dispatch reconcile` helper (the behavioral gate):** diff the live team `config.json` members against the workflow's active, non-terminal, worktree-bearing entities and EMIT the orphan set — agents whose entity is archived/terminal or whose cycle is superseded. Optionally also emit PR-pending entities whose PR is MERGED-but-not-advanced and in-flight branches BEHIND `origin/next`. The FO runs `reconcile` at **idle and each merge**. This replaces "remember to" (drifts) with "the tool computes the orphans from state" (reliable).

## Acceptance criteria (behavioral — no grep self-proof)

**AC-1 — `spacedock dispatch reconcile` computes the orphan set from state.** Given a team roster + the workflow's entity states, the helper lists agents whose entity is terminal/archived or whose cycle is superseded.
Verified by: a Go test feeding a synthetic roster + entity set and asserting the emitted orphan list (a real fixture in → expected orphans out), and a flip (mark an entity active → its agent drops off the orphan list).

**AC-2 — the contract teardown + supersede steps are enforced where they can be.** The Merge-and-Cleanup teardown and the supersede-shutdown are documented in the FO contract AND, where a code path exists, the reconcile helper would flag a violation (an archived entity with a live-named agent shows as orphaned).
Verified by: the AC-1 test covers the archived-entity-with-live-agent case; the prose steps are reviewed (doc-as-deliverable for the contract wording, with the helper as the behavioral backstop — same pattern as v1/2a).

**AC-3 — (optional, if cheap) reconcile also surfaces PR + branch drift.** The helper flags pr-set entities whose PR is MERGED but status is non-terminal, and in-flight worktree branches behind `origin/next`.

## Notes
- Scaffolding (FO contract) + a new `internal/dispatch` (or `internal/cli`) helper + tests — goes through a worktree. Coordinates with the `2a` opt-in-guard pattern (a helper the FO runs, not a universal mandate).
- The agent-teardown is the PRIMARY ask (lingering ensigns); the PR/branch-drift surfacing (AC-3) is the same idle-reconcile sweep generalized — fold in only if it doesn't bloat the helper.
- Captures the dogfooding lesson: this session's FO failures (leaked agents, reactive rebase) are precisely a memory-based-discipline ceiling; the fix is to compute from state.
