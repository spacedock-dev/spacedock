# Handoff — sprint 0221-layered-fo, paused before the v0.22.1 cut

You are the incoming first officer / Commander for sprint `0221-layered-fo` in `docs/dev`. The prior session drove six members to `main` (`c25fee26`) and **paused before tagging `v0.22.1`** because the captain and the pre-cut audit found a **structural antipattern**. Your job this session is **not** to finish the tag. It is to **validate the findings below, then propose a plan to the captain (CL) before executing anything.**

Read `docs/dev/_debriefs/2026-06-19-01.md` first — it has the full detail. This is the action layer.

## Hard stop
- **Do NOT fire the `v0.22.1` tag.** The headline deliverable (Haiku-operable verbs the FO actually uses) is **not delivered**.
- **Do NOT merge anything new** until CL approves your proposed plan.
- PR **#403** (`strip-deferred-tier-vocabulary`, a pre-cut fix) is **open, CI running, unmerged** — decide its fate as part of your proposal, don't auto-merge it.

## State you inherit
- `main` = `c25fee26`; last tag `v0.22.0`; **not tagged.**
- Merged this sprint: 3e (#398), 6re (#397), rgq (#399), mz (#400), cleanup (#401), czw (#402) — all archived PASSED.
- In-flight: `ga` `strip-deferred-tier-vocabulary` at validation, #403 CI running.
- Deferred: `72` (validation), `kt` (backlog), `y2` (backlog).
- Resident workers: `comm-officer`, the `ga` impl + validation ensigns. Open worktree: `.worktrees/spacedock-ensign-strip-deferred-tier-vocabulary`.
- Unfiled follow-up: `status-set-worktree-clear-guard` (the `handlers.go:116-118` over-broad terminal guard).

## The finding to validate (don't take on faith — reproduce it)
The verb work was sliced by **tech-stack layer, not end value**, so the verbs ship but the contract doesn't operatively use them. **Independently confirm or refute** each of these before proposing:
1. **State verb unwired.** `skills/first-officer/references/first-officer-shared-core.md` `### Split-Root State Sync` still presents `git -C add/commit` as the path with a pre-verb "Preferred — when the status tool owns add+commit" that never names `spacedock state commit`. → Does the contract operatively instruct the verb, or the hand sequence?
2. **Merge verb overstated.** `fo-merge-core.md` `«merge.guard»` body claims the verb "invoke[s] the registered merge hook" and runs "as one call," but `internal/status/merge.go` shows a re-entrant partial envelope (`armed: invoke the hook, then re-run`). → Does the contract description match the shipped verb?
3. **Mis-targeted CI.** rgq/mz are subcommand-only; the live-e2e FO never calls those verbs. → Did the live lanes execute any line of those changes, or was the deterministic Go suite the real (and sufficient) proof?
4. **No end-value validation exists.** `kt` (the live Haiku drive) was the only proof that an FO drives via the verbs, and it was deferred. → Confirm nothing else validates verb *usage*.

## Recommendations (validate, then propose — adjust as your validation warrants)
1. **Re-decompose by end-value vertical slices.** Each verb becomes one task that ships: the binary **+** the operative contract rewiring (the verb *replaces* the hand sequence, not annotates it) **+** an end-value check that the FO uses it. Kill "binary-only" and "wiring-only" tasks.
2. **Un-defer `kt`** (or an equivalent) as the integration/end-value proof, and sequence it to actually run before the cut — not last-and-deferred.
3. **Rewire the contract to use the verbs:** `### Split-Root State Sync` → call `spacedock state commit <slug>` as the path; Merge-and-Cleanup → invoke `merge guard` with an **accurate** description of the re-entrant partial envelope and the FO-owned hook-invocation / worktree-removal / teardown.
4. **Add a contractlint guard binding each `«fn»` body to its shipped verb** — claimed scope/usage vs actual routing + behavior — so "declared-but-unwired" and "overstated" fail mechanically (the routing oracle proves the command exists, not that the prose is accurate).
5. **Right-size CI to the change:** subcommand-only → deterministic Go tests, no live-e2e; FO-drives-the-verb → live-e2e (the end-value gate). Stop waiving live flakes on changes the live lanes can't validate.
6. **Reconcile the byte-identical tension:** behavior-preservation and "make the verb the path" conflict — decide explicitly per region which wins (the rewire *should* change the prose; that's the point).
7. **Resolve #403** (`ga`) — merge it as the standalone tier-vocab fix, or fold it into the re-work.
8. **File `status-set-worktree-clear-guard`** if you keep it in scope.
9. **Decide the release identity:** if the foundation re-works cleanly, `v0.22.1` may still be right; if scope shifts, propose the version/theme to CL.
10. **Spacedock-framework issues to weigh** (from the debrief): the `dispatch build` flag-vs-stdin contract skew; the boot `team_state` vs `reconcile` disagreement; the standing-teammate hand-forwarding ergonomics.

## How to start
1. Boot the FO (`spacedock claude`), run `status --boot`, read this handoff + the debrief.
2. **Validate** findings 1–4 against the live tree (`c25fee26`) — reproduce, don't assume.
3. **Propose to CL**: a short plan — which findings hold, the re-decomposition, the sequence, and the release call — and **wait for approval before dispatching any worker.** The captain explicitly asked that you validate and propose on start rather than execute.

The prior Commander's own lesson, for the record: it executed a layer-sliced decomposition without challenging whether each slice delivered end value, ran six rigorous validations at the wrong granularity (layer, not capability), and fought irrelevant live-CI flakes — so a clean-looking cut hid an undelivered deliverable. Don't inherit that; validate the value, then build it as slices that can be validated.
