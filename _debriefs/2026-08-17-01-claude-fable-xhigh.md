---
session-date: 2026-08-17
sequence: 1
harness: Claude Code
model: claude-fable-5
model-version-build: claude-fable-5 (reasoning effort xhigh)
first-commit: 884b55af0
last-commit: 04e936a12c4792943fdd0ca707442b009591daba
duration: ~34h
---

# Session Debrief — 2026-08-17 #1

The 0.27 endgame: the feedback-rejection arc root-caused (three product defects, not a flake), rebuilt as a six-layer stacked PR set, proven live on both runtimes with strict determined-shape grading, merged atomically, and shipped as v0.27.0-pre7. The session crossed one context compaction on the durable state alone.

## Shipped

All six landed atomically via stack #720 (lander [#718](https://github.com/spacedock-dev/spacedock/pull/718), merge commit 61dd8e435); member PRs closed as merged-via-stack with `local-merge` sentinels. Released as **v0.27.0-pre7** (tag on 04e936a12; e2e-gate satisfied by a real green run on the SHA — claude 17/17 attempt 1, codex 17/17 attempt 2 after one keep-moving flake).

- **189** `simplify-feedback-rejection-flow` — [#718](https://github.com/spacedock-dev/spacedock/pull/718). Redesigned the rejection skill to five single-publish steps; restored the round-recorder's documented entity operand; a usage error is a correctable mistake, not a terminal hold.
- **hz** `repair-codex-rejection-round-recording` — [#719](https://github.com/spacedock-dev/spacedock/pull/719). Codex completed the live rejection flow but the round graded unrecorded; fixed the dirty-tree dead end, made inline `state commit` durable, and recognized codex quote runs; regression carries the verbatim run bytes.
- **b8e** `fo-workflow-fit-gate` — [#721](https://github.com/spacedock-dev/spacedock/pull/721). The captain's workflow-fit admissibility gate added to the FO write core, ahead of the new-entity-file rule.
- **ra** `prepare-initial-gated-stage-from-seed` — [#722](https://github.com/spacedock-dev/spacedock/pull/722). A workflow whose initial stage is gated deadlocked the scheduler; a committed clean seed is now a preparable gate (keyed on `initial: true`, guards preserved).
- **7x** `red-auto-continue-gate-bypass` — [#723](https://github.com/spacedock-dev/spacedock/pull/723). Human-gate-bypass detection made host-neutral and honest: gate state graded on every entity copy fail-closed, and only a resolved gate accuses (sanctioned withdrawals pass). Caught a real sonnet bypass on the release run — a true positive with a self-confession.
- **zq** `run-rejection-journey-in-team-mode` — [#724](https://github.com/spacedock-dev/spacedock/pull/724). The rejection journey now runs in team mode — the reuse routing it exists to prove — graded strict on the determined worker topology per host; withdraw-recovery generalized to attempt-N open-attempt selection; heading strict+pin; Cycle-line target file pinned in the fixture.

## Filed (backlog)

- **3rdq** `record-real-grades-in-journey-metrics` — journey metrics record unbound failures as passes; bit us live this session (a red lane produced 18 passing-shaped records).
- **rcpa** `own-sonnet-gate-conn-bypass` — sonnet resolved a gate with no conn grant, attributed it to the captain, then self-confessed; the FO recited the rule it violated, so prevention likely needs a product-side conn-grant artifact check.
- **zf7r** `own-claude-early-rejection-round-record` — claude recorded the rejection round before the correction (entries=2); the recognizer's `entries=4` pin reports it as "never invoked" — owns the mode and both label-honesty repairs in that oracle.
- **3ptr** `own-codex-filing-variant-miss` — codex filed one filing-fixture variant of two, once; fixture-wording check first, then a measured loop.
- Filed and shipped same session: **zq**, **7x**, **b8e**, **ra**. `own-rejection-flow-adherence-residuals` folded into zq at filing and archived (captain: "same thing").

## Non-PR commits (workflow-only)

- `df0bd50d9` live-ci registry AUDIT annotations from the 14-finding harness audit — the rejection-flow and auto-continue blocks were dropped at merge as superseded by the stack; all other journeys' blocks preserved.
- `04e936a12` release stamp; tag `v0.27.0-pre7` pushed; release.yml fully green (e2e-gate, goreleaser, journey-ledger, edge-advance).
- State-side records: the live-harness audit debrief (…-02), the blind two-sided shape-attribution debrief (…-03), the coordination handoff (…-01, updated through the session), gate rooms for every attempt, and the digest gate review in `_reviews/`.
- Commits between the prior anchor and the pre6 tag belong to the prior codex session's record (pre6's tag message covers them).

All other session commits are rolled up in the shipped stack above.

## Decisions

- No bare-mode journey invocations unless the scenario is specifically testing bare mode; team mode is the rejection journey's shape.
- XFAIL requires a product-task owner with a real proposed approach; tracking stubs are not owners ("if it is a known failure it should have a product task owner").
- Stack-greening fixes fold into existing entities; no new entities for them.
- Heading grammar ruled strict + pin the declared level in the fixture, rather than loosening the checker.
- Composite green accepted for the merge (each lane proven 17/17 on the shipping bytes, residual reds attributed to owned live-model modes); the release gate still earned a real green run on the release SHA.
- Gate presentations lead with what the captain cares about — a one-page digest artifact, not the entity body.

## Issues — Workflow

- Post-merge worktree cleanup glob swept five live entities' clean worktrees (`grep ensign` instead of the six finalized slugs); all five restored at their recorded paths on their exact prior heads; no committed work lost. Lesson recorded in the handoff.
- Shared state checkout under concurrent writers: non-fast-forward pushes, a scratchpad collision, and a cross-agent stash swap earlier in the window — now codified as workflow rules (path-scoped commits, no stash, shared-path stop signal).
- The session task list was found wiped mid-session; durable entity/handoff records carried the work across.

## Issues — Spacedock

All tracked in-repo as entities or on the captain's filing list; none filed externally (this repo is the spacedock repo):

- `gate prepare` briefings lack subspace 0.14.0's required `context` node and pin `git-root://` artifact URIs its loader cannot resolve — every captain presentation needs a hand-built presentation copy (full recipe preserved; repair directions in the friction dossier).
- `status --read` inspects the main entity copy and misses reports in the registered worktree.
- Product reads hardcode `$HOME/.claude`, so `CLAUDE_CONFIG_DIR` operators silently lose reuse/teammates/reconcile (spin-off on the filing list).
- `merge guard` has no stacked-lander shape: member PRs closed as merged-via-stack needed hand-recorded `local-merge:{merge-sha}` sentinels after ancestry verification.
- State hooks emit ~50KB of archived-entity schema warnings on every commit (`unknown gate application field`).
- Journey metrics lying-records defect (filed as 3rdq); round-oracle label conflations (filed within zf7r).

## Observations

_(none recorded)_

## Agent Testimonial

- Date: 2026-08-17
- Harness/runtime: Claude Code (background job session)
- Model: Claude Fable 5
- Model version/build: claude-fable-5, reasoning effort xhigh
- Session scale: ~11 tasks touched; ~15 workers dispatched; 6 PRs + 1 stack merged

Spacedock's durable state is what made this session survivable: the work crossed a context compaction, two validator disputes were settled from recorded bytes rather than memory, and a multi-day arc lost no decision. The gates forced evidence to exist before claims — every refusal I hit (dirty tree, entered-stage guard, open-PR finalization) turned out to be correct, which is a strong property. The cost is the ceremony floor: a one-sentence fixture pin cost a dispatch envelope, an entity amendment, a validator delta, and a gate attempt; one validation was presented four times before the presentation shape matched what the captain needed. The friction concentrates where spacedock meets external surfaces — the subspace briefing skew, GitHub's model of stacked PRs, the shared state checkout — while its own core guards were consistently right. Driving the same work without it, I would have moved faster per step and almost certainly lost the thread at compaction, shipped at least one plausible-but-wrong fix (the gates-read placement remedy that codex refuted), and had no byte-level record to adjudicate the validator disputes.

## What's Next

- Soak pre7, then v0.27.0 stable: if main stays frozen, the stamp-only delta qualifies for the equivalent-prior-green waiver; the stable changelog (with dogfood numbers 392/282/319) is the only genuinely new work.
- Nothing currently dispatchable; 168 entities in view. Top ideation candidate: `respect-gate-application-authority` (score 1.0) — directly downstream of the sonnet conn bypass.
- Captain filing list: CLAUDE_CONFIG_DIR spin-off, `status --read` worktree-miss, H4 substring asymmetry in the cycles section finder, audit findings 4/5/6/9/13, pi-dialect extractors (0.28 ledger).
- Deferred local-refit entity (re-derive the dev README template and installed pr-merge mod from shipped versions).
- Keep-moving codex occurrence on watch at one strike; recurrence files an owner (likely extending `repair-sonnet-live-flakes` cross-lane).
