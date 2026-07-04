---
title: FO contract keep-moving posture — approval is a trigger, parallel async when independent, no turn-end on async launch with work remaining, captain correction narrows not halts
status: backlog
source: "Shaping FO (2026-06-20, 0223 retrospective): repeatedly violated the contract's existing 'do obvious reversible work without ceremony' + 'keep dispatching other ready entities when one blocks' on pi. Patterns: (1) post-approval pause — asked 'want me to advance + dispatch?' after a gate approval (a reversible step the contract already permits); (2) sequential filing of independent followups instead of parallel; (3) turn-end on async launch with independent work remaining (the pi-subagents skill already forbids this; lift into the FO contract so it's host-neutral); (4) captain-question = full stop — conflated 'this member needs rework' with 'stop the session.' The async substrate exists to keep work moving while the FO coordinates; the contract should say so explicitly."
score:
started:
completed:
verdict:
worktree:
issue:
sprint: 0250-fo-behavioral-discipline
sprint-readiness:
id: vcmbpj2dxq8bpkys2t0vkvq7
---

# FO contract keep-moving posture

## End value

The FO contract's Working Principles / Clarification sections explicitly require the keep-moving posture the async substrate enables, so every FO on every host holds it — not just pi sessions that happen to load the pi-subagents skill. A captain gate approval triggers the FO's next action; independent entities dispatch in parallel; an async launch doesn't end the turn when independent work remains; a captain correction narrows scope without halting the session. The false stops that serialized this session's work are contract violations, named and rejectable at the gate.

## Problem — four false-stop patterns observed this session

1. **Post-approval pause.** After captain gate approval, the FO asked "want me to advance + dispatch?" — a reversible step the contract already permits (the "do obvious reversible work without ceremony" principle). Approval was treated as a turn boundary instead of the trigger for the next action.
2. **Sequential filing of independent followups.** Filed `status-validate-determinism`, paused, `state-init-precommit-hook`, paused, `pi-launch-fnm-multishell-race`… when these were independent and could have been filed + ideation-dispatched in parallel async. Under-uses the async substrate and serializes work the captain already authorized.
3. **Turn-end on async launch with independent work remaining.** Fired one async ensign and ended the turn "waiting," despite independent FO work remaining (other followups to shape, findings to record). The pi-subagents skill already forbids this ("Do not end your turn immediately after launching an async child if you promised to keep working"); the FO contract is silent, so it's not host-neutral.
4. **Captain-question = full stop.** When the captain interrogated a mechanism (the symlink, the safehouse flags), the FO stopped *all* shaping. But the question retired one mechanism; the unaffected members kept being drivable. "This member needs rework" was conflated with "stop the session."

## Approach — four strengthenings to the shared core

Text edits to `skills/first-officer/references/first-officer-shared-core.md` (Working Principles + Clarification and Communication sections). This is the shipped FO contract — a high-stakes surface per the dev-workflow proof policy (the shipped contract and scaffolding).

### S1 — "Approval is a trigger, not a turn boundary." (Working Principles)

> A captain gate approval is the trigger for the FO's next action, not the end of its turn. On approval, the FO advances the entity and dispatches the next stage before ending its turn — unless the next stage is itself a gate or the captain directed otherwise. Asking "want me to advance + dispatch?" after an approval is a contract violation: advance+dispatch is reversible work the contract already permits.

### S2 — "Parallel async when independent." (Working Principles)

> When two or more independent entities are ready for the same stage (or independent followups are ready to file + ideate), the FO dispatches them in parallel async, not sequentially. Filing one task, pausing for acknowledgment, then filing the next under-uses the async substrate and serializes work the captain already authorized. The FO files + dispatches the independent set in one motion and reports the batch.

### S3 — "Don't end turn on an async launch when independent work remains." (Working Principles — host-neutral lift of the pi-subagents skill rule)

> Launching an async ensign does not end the FO's turn unless no independent FO work remains. While an async child runs, the FO continues: shapes the next entity, records findings, files followups, prepares the next dispatch. The FO ends its turn only when (a) it is blocked on the async result with no other work, or (b) it has reached a gate/captain decision. Ending a turn to "wait" on an async child while independent work sits undone is a contract violation.

### S4 — "A captain correction narrows scope; it does not halt the session." (Clarification and Communication)

> When the captain questions or corrects a mechanism, the FO re-shapes the affected member — and continues driving the unaffected members. A correction to one member's design is not a stop signal for the rest of the sprint. The FO pauses only the corrected member's dispatch; it keeps the independent members moving and surfaces the re-shape when the correction is folded. Conflating "this member needs rework" with "stop the session" under-uses the parallel substrate.

### The line these preserve

The captain's interrogation is the primary de-risking mechanism (it retired two wrong mechanisms this session). S4 must not become "never stop for a captain question." The line: **stop for gates and genuine scope forks; don't stop for reversible execution steps, don't stop when independent work remains, don't end turn just because one async is in flight.**

## Scope

In scope:
- The four strengthenings (S1–S4) as text additions to `first-officer-shared-core.md` (Working Principles + Clarification and Communication).
- Reconcile with the existing "do obvious reversible work without ceremony" and "keep dispatching other ready entities when one blocks" principles — these strengthenings make those explicit for the four observed violation patterns; avoid contradiction.

Out of scope:
- Changing the gate ceremony, dispatch core, or merge module — this is posture/Working-Principles text only.
- Host-specific guidance — S3 is the host-neutral lift of a pi-subagents-skill rule; it does not edit the pi-subagents skill.
- Enforcement via a binary guard — the contract is prose; the AC is a live FO scenario observing the posture, not a code gate. (A future task could add a contractlint structural check if a prose-property can be bound to an independent value, but that's out of scope here — see the proof-policy note below.)

## Acceptance criteria (provisional — ideation finalizes; proof = behavior)

**AC-1 — The four strengthenings (S1–S4) land in `first-officer-shared-core.md`'s Working Principles + Clarification sections, reconciled with the existing reversible-work and keep-dispatching principles.**
Verified by: a structural review that the four additions are present, non-contradictory with the existing principles, and host-neutral (no pi-specific framing leaks into the shared core). (This is a contract-text change; a contractlint "the phrase is present" check is the prose-grep tautology the dev-workflow proof policy bans — so the structural review is the gate, not a prose check.)

**AC-2 — A live FO scenario observes the keep-moving posture.**
Verified by: a live drive (any host lane) where the FO, after a gate approval, advances + dispatches without asking; where independent entities dispatch in parallel; where an async launch does not end the turn while independent work remains. The scenario asserts the *behavior* (post-approval dispatch happens, parallel dispatch happens, turn continues past an async launch), not the prose. This is the dogfood — the contract change is proven by an FO that holds the posture.

**AC-3 — The change does not regress the captain's de-risking role.**
Verified by: the live scenario includes a captain correction that narrows scope (one member re-shaped) while an independent member keeps moving — proving S4 preserves the captain's interrogation as a de-risking lever without halting the session.

## Proof-policy note

A Working-Principles prose addition resists a clean code gate — a contractlint "the phrase is present" check is exactly the prose-grep tautology the dev-workflow proof policy bans (the matched text was written by the same implementer the check polices). The real assurance is AC-2's live drive observing the posture, plus the structural review (AC-1) binding two independent values — the named principle and the observed behavior — that can diverge. Ideation confirms this is the legitimate shape; if a binary-enforceable property emerges (e.g. a dispatch-log structural check that parallel dispatch occurred), record it, but don't force a prose gate.

## Test plan

- Structural review (AC-1): the four additions present + reconciled in `first-officer-shared-core.md`.
- Live drive (AC-2): a `pi-live` (or claude-live/codex-live) scenario observing post-approval-advance, parallel dispatch, no-turn-end-on-async. This is the dogfood — the contract change drives the scenario that proves it.
- AC-3: the live scenario's captain-correction-narrows element.
- This is a shipped-contract change — high-stakes surface. Detached adversarial audit at validation; `claude-live` + `codex-live` + `pi-live` regression (the change is to the shared core every host adapter inherits).

## Related

- `skills/first-officer/references/first-officer-shared-core.md` — the file to edit (Working Principles + Clarification and Communication).
- The pi-subagents skill's "Do not end your turn immediately after launching an async child if you promised to keep working" rule — S3 is the host-neutral lift.
- The 0223 Shaping FO debrief (`docs/roadmap/0223-pi-dispatch-contract/debrief-shaping-2026-06-19.md`) — the retrospective that surfaced the four patterns.
- `status-validate-determinism`, `state-init-precommit-hook`, `pi-launch-fnm-multishell-race`, `pi-live-shared-scenario-runners` — the independent followups this session filed sequentially (the pattern S2 addresses).
