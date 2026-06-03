---
id: 0nadpgpzer0jhrvcxeg52az2
title: reconcile auto-discovery (no --team-name) is not session-aware — newest-mtime team glob picks stale prior-session or parallel-session configs, poisoning roster-derived drift classes
status: backlog
source: session-11 FO (2026-06-03) — observed bare `dispatch reconcile` resolve a two-day-old prior-session team config and report archived-entity agents as Class A; captain flagged it ("did reconcile not consider repeated or parallel sessions?")
score: "0.19"
worktree:
started:
completed:
verdict:
issue:
---

`spacedock dispatch reconcile` run without `--team-name` resolves the team by globbing `~/.claude/teams/*/config.json` and taking the **most-recently-modified** config containing a `spacedock:ensign` member (`internal/claudeteam/reconcile.go` `LoadReconcileTeam`, lines 43-87). The code comment calls this "a stable proxy for 'the live team in this session'" — but it is a newest-mtime heuristic, not a session-identity check, and it misfires on both repeated and parallel sessions.

## Problem

The roster-derived drift classes (A lingering, B superseded, C un-advanced-PR) are computed from the resolved team's member roster cross-referenced with entity state. When auto-discovery resolves the **wrong** team, those classes are poisoned:

- **Repeated (sequential) sessions.** Claude Code leaves `config.json` on disk after a session ends (registry-desync, anthropics/claude-code#36806). A session that does no dispatch never runs `TeamCreate`, so no fresh team dir exists — newest-mtime resolves to a **stale prior-session team**. Observed in session 11: bare reconcile picked `spacedock-v1-dev-20260601-0729-...` (two days old) and reported its archived-entity agents (`status-enumeration-and-validation-{implementation,validation}`) as Class A "lingering." The prescribed Class A action is `SendMessage shutdown_request` — meaningless against a dead session's agents.
- **Parallel (concurrent) sessions.** Two live sessions (e.g. a Claude FO + a codex peer FO, or two Claude FOs) each own a team dir. Newest-mtime **races** between them and can resolve to the *other* session's live team. Class A/B teardown then fires `shutdown_request` at a **parallel session's legitimately-live agents** — active cross-session interference, not just a false positive.

**Precision:** Class D (stale branch) and E (stale local main) are computed from git + filesystem state (worktree-behind, local-main-behind), independent of the team roster — they are session-safe and were accurate in the observed run. Only the **roster-derived classes (A/B/C)** are vulnerable.

**Unused affordance:** the config carries `leadSessionId`, loaded into `ReconcileTeamState`, with a comment saying it is "useful for narrowed discovery when --team-name is omitted" — but discovery never references it; it sorts by mtime only. The contract's event-loop step 0 also writes `[--team-name {team_name}]` as **optional**, which invites the unsafe bare invocation.

## Proposed approach

Make roster-derived drift classes provably session-scoped. Ideation picks among / combines:

1. **Degrade-to-git-only without a team identity.** When `--team-name` is omitted (or no session-scoped team resolves), emit only the git/filesystem classes (D/E) and suppress the roster classes (A/B/C), with a one-line note that roster reconciliation needs a team identity. Safe by construction: no team → no roster claims.
2. **Narrow auto-discovery by `leadSessionId`.** Match the discovered config's `leadSessionId` against the current session, using the affordance the struct already loads, instead of mtime. (Requires plumbing the current session id into the helper — confirm the env/source first; this is the riskiest unknown to exercise in ideation.)
3. **Contract update.** Reflect the safe usage in the FO event-loop step 0: roster reconciliation requires `--team-name` (the FO's own `TeamCreate` name); bare `reconcile` is git-only. Pair the prose with the code gate above so the guarantee is enforced, not just documented.

## Out of scope

- Changing the drift-class definitions themselves (A/B/C/D/E semantics stay).
- Deleting stale on-disk team dirs (forbidden by the runtime adapter's NEVER-delete constraint; this is about not *trusting* them, not removing them).
- The git-derived classes D/E — already session-safe.

## Acceptance criteria

Proof is Go tests over the reconcile loader with multiple team-config fixtures (the existing `reconcile_test.go` already builds hermetic `~/.claude/teams` fixtures). Ideation refines.

**AC-1 — Bare reconcile (no `--team-name`) never emits roster-derived classes from a non-current-session team.**
Verified by: a Go test seeding two team configs (a stale/foreign one with ensign members + entity state that would trip Class A, and the absence of a current-session team) and asserting the bare-invocation output contains no Class A/B/C entries sourced from the foreign roster.

**AC-2 — A team identity (explicit `--team-name` or session-matched discovery) still produces correct roster classes.**
Verified by: a Go test passing the live team's name and asserting Class A/B/C are computed against that roster as before (no regression to the explicit path).

**AC-3 — Git/filesystem classes (D, E) are still emitted without a team identity.**
Verified by: a Go test asserting bare reconcile over a worktree behind `origin/next` still reports Class D (the session-independent path is unaffected).

**AC-4 — The FO contract reflects the safe usage.**
Verified by: a presence check / oracle over the event-loop step-0 prose confirming it states roster reconciliation requires a team identity and bare reconcile is git-only (proof at the text's own level; paired with the AC-1 code gate that enforces it).

## Test plan

- Go unit tests for AC-1..AC-3 over `LoadReconcileTeam` + the reconcile assembly, using the existing hermetic team-config fixtures. Cost: low.
- AC-4 instruction-text invariant. Cost: trivial.
- If approach #2 (leadSessionId narrowing) is chosen, exercise the current-session-id source end-to-end first (the riskiest unknown) before committing to it.
- High-stakes surface (FO-event-loop teardown machinery that can `shutdown_request` a parallel session) → detached adversarial audit required before merge.

## Notes

- Surfaced session 11 while finishing the v0.19.4 release recovery; the FO ran bare reconcile with no live team and got a stale-team Class A report. Captain's question ("did reconcile not consider repeated or parallel sessions?") is the title.
- 0.19.5, test-improvement themed: the deliverable is the fix plus the regression tests that lock the session-scoping guarantee.
- The reconcile helper is named as a Phase 1 candidate in `docs/dev/_proposals/binary-simplification-roadmap.md` (FO-event-loop hot path) — this correctness fix should land before any consolidation of that surface.
