---
title: FO contract keep-moving posture — approval is a trigger, parallel async when independent, no turn-end on async launch with work remaining, captain correction narrows not halts
status: implementation
source: "Shaping FO (2026-06-20, 0223 retrospective): repeatedly violated the contract's existing 'do obvious reversible work without ceremony' + 'keep dispatching other ready entities when one blocks' on pi. Patterns: (1) post-approval pause — asked 'want me to advance + dispatch?' after a gate approval (a reversible step the contract already permits); (2) sequential filing of independent followups instead of parallel; (3) turn-end on async launch with independent work remaining (the pi-subagents skill already forbids this; lift into the FO contract so it's host-neutral); (4) captain-question = full stop — conflated 'this member needs rework' with 'stop the session.' The async substrate exists to keep work moving while the FO coordinates; the contract should say so explicitly."
score:
started: 2026-07-04T10:38:18Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-fo-contract-keep-moving-posture
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

> A captain gate approval triggers the FO's next action, not the end of its turn: on approval the FO advances the entity and dispatches the next stage before yielding — unless that next stage is itself a gate or the captain directed otherwise. Asking "want me to advance + dispatch?" is a violation; advance+dispatch is reversible work the contract already permits. A verification-bearing decision (a merge, a triage call) still holds to the evidence-bar principle — keep-moving speeds the reversible dispatch, not a decision past the verification its change requires.

### S2 — "Parallel async when independent." (Working Principles)

> When two or more independent entities are ready for the same stage (or independent followups are ready to file + ideate), the FO dispatches them in parallel async, not sequentially — filing one, pausing for acknowledgment, then filing the next serializes work the captain already authorized. It dispatches the independent set in one motion and reports the batch, but only for units the smallest-sufficient-mechanism gate already judged worth a dispatch: keep-moving parallelizes chosen dispatches, it does not escalate the mechanism (a unit better done in-house stays in-house, alongside the dispatched ones).

### S3 — "Don't end turn on an async launch when independent work remains." (Working Principles — host-neutral lift of the pi-subagents skill rule)

> Launching an async ensign does not end the FO's turn while independent FO work remains: the FO keeps shaping the next entity, recording findings, filing followups, preparing the next dispatch. It yields only when blocked on the async result with no other work, or at a gate/captain decision. Ending a turn to "wait" on an async child while independent work sits undone is a violation.

### S4 — "A captain correction narrows scope; it does not halt the session." (Clarification and Communication)

> When the captain questions or corrects a mechanism, the FO re-shapes the affected member and keeps driving the unaffected ones: it pauses only the corrected member's dispatch and surfaces the re-shape when the correction is folded. A correction to one member is not a stop signal for the rest of the sprint — conflating "this member needs rework" with "stop the session" under-uses the parallel substrate.

### The line these preserve

The captain's interrogation is the primary de-risking mechanism (it retired two wrong mechanisms this session). S4 must not become "never stop for a captain question." The line: **stop for gates and genuine scope forks; don't stop for reversible execution steps, don't stop when independent work remains, don't end turn just because one async is in flight.**

### Reconciliation with z25 and zm — the sprint's named composition tension

The 0250 sprint doc sequences vcm LAST because keep-moving must be reconciled against z25 ("required verification follows from what changed") and zm ("justify before climbing a mechanism rung"); driven naively, the three read as contradictory edits to the same Working-Principles paragraphs. They compose without contradiction because **they govern three orthogonal axes of a single FO action, and keep-moving is subordinate to the other two on the axes they own:**

- **z25 owns the verification axis** — *what evidence a decision requires.* A merge / gate / triage decision moves only when the verification that follows from what changed actually ran and passed.
- **zm owns the mechanism axis** — *which rung to climb.* Before dispatching a worker / spinning a workflow / opening a PR, the FO names why the cheaper rung can't do it.
- **vcm owns the motion axis** — *don't idle or serialize once z25 and zm have chosen.* Keep-moving accelerates the CADENCE of correctly-verified, correctly-sized actions; it never sets what proof a decision needs (z25) or what mechanism it uses (zm).

Per strengthening:

- **S1 (approval is a trigger) yields to z25.** "Advance + dispatch without asking" is the reversible dispatch step after a gate approval (advance the entity to its next stage, fire the ensign) — NOT a merge/triage decision. A merge still holds to z25's bar; keep-moving never licenses "advance past the evidence" on a verification-bearing decision. The existing "unless the next stage is itself a gate" carve-out already routes verification-bearing decisions back through the gate/z25 path. (S1's shipped wording now says this inline.)
- **S2/S3 (parallel async, no turn-end) presuppose zm's gate.** They govern how to sequence dispatches the FO has ALREADY decided to make — not whether to dispatch. Parallel async is the right posture ONCE zm's smallest-sufficient-mechanism gate has judged that N independent units each warrant a worker; it does not escalate the mechanism. "Fire the independent set in parallel" means "for the units zm's gate already sent to a worker," never "turn every unit of work into a parallel worker." (S2's shipped wording now says this inline.)
- **zm's friction-#1 fold REINFORCES keep-moving.** zm's "don't re-run stage-owned verification at gate time" and vcm's "don't stall, don't add ceremony" are the same anti-busywork instinct from two sides — zm forbids the re-verification busywork, vcm forbids the pause/serialize busywork. They compose, they don't collide.
- **S4 (correction narrows, doesn't halt) preserves z25's captain-interrogation lever.** z25 makes the captain's gate judgment the primary de-risking mechanism; a mis-framed gate is its dangerous failure. S4 does NOT weaken that — it pauses the corrected member (the interrogation does its de-risking work there) while the independent members keep moving. AC-3 proves S4 preserves the lever rather than dulling it.

The one-line invariant the composition preserves: **keep-moving changes the CADENCE of action, never the BAR (z25) or the RUNG (zm).** A "keep-moving" action that skips required verification or over-climbs the mechanism ladder is not keep-moving — it is a z25 or zm violation wearing keep-moving's name, and the reconciled wording is written to make that mis-read impossible.

## Scope

In scope:
- The four strengthenings as text additions to `first-officer-shared-core.md` — S1–S3 in Working Principles, S4 in Clarification and Communication.
- Reconcile with the existing "do obvious reversible work without ceremony" and "keep dispatching other ready entities when one blocks" principles (these strengthenings make those explicit for the four observed patterns) AND with the sibling 0250 principles z25 (evidence bar) and zm (smallest-sufficient-mechanism gate), which land in the SAME Working-Principles section — see Reconciliation. The strict-serial Wave-1 order (z25 → zm → vcm, per the 0250 sprint doc) is what makes this composition safe: vcm lands into the shape z25/zm established and cross-references them, never contradicting them.
- Leanness: the four strengthenings are added in their terse resident one-line form (posture rules, not a lazy-reference candidate the way zm's gate is). The combined Startup + Working-Principles union is measured against the 0.24.0 baseline by the sprint's independent preflight staff review, not by this entity alone (the 0250 leanness constraint).

Out of scope:
- Changing the gate ceremony, dispatch core, or merge module — this is posture/Working-Principles text only.
- Editing the sibling entities' surfaces: vcm does NOT touch k74g's `«engage»` block / Startup edits, z25's evidence-bar principle + `present-gate` rule, or zm's mechanism gate. vcm adds its own four principles and cross-references the siblings; the keep-moving posture is what makes k74g's `«engage»` safe to lean on unattended (the 0250 theme), but vcm does not edit the engage block.
- Host-specific guidance — S3 is the host-neutral lift of a pi-subagents-skill rule; it does not edit that skill (which is Pi's external host skill, not a repo file — see Related).
- Enforcement via a binary guard — the contract is prose; the AC is a live FO scenario observing the posture, not a code gate. (A future task could add a contractlint structural check if a prose-property can be bound to an independent value, but that's out of scope here — see the proof-policy note below.)

## Acceptance criteria

**AC-1 (mechanism — serves the value AC-2/AC-3 measure) — the four strengthenings (S1–S4) land in `first-officer-shared-core.md`, composed without contradiction into the Working-Principles section z25 and zm also edit.**
The four additions are present in Working Principles (S1–S3) + Clarification and Communication (S4), host-neutral (no pi-specific framing in the shared core), and reconciled with the pre-existing reversible-work and keep-dispatching principles AND with the sibling sprint principles z25 (evidence bar) and zm (smallest-sufficient-mechanism gate) — the shipped wording subordinates keep-moving to those two on the axes they own (see Reconciliation).
- *Verified by:* the sprint's independent cross-member coherence review (the 0250 preflight staff review) reading the COMBINED Working-Principles section for non-contradiction + host-neutrality. This is a contract-text change; a contractlint "the phrase is present" check is the prose-grep tautology the dev-workflow proof policy bans, so the structural review is the gate, not a prose check. AC-1 asserts the mechanism shipped; per the ideation stage-def's value-AC rule it counts only paired with AC-2/AC-3, which measure the behavior it exists for.

**AC-2 (value, behavior — measured against a baseline that moved the wrong way) — a live FO drive holds the keep-moving posture where the current contract does not.**
On a constructed scenario reproducing the 0223 session's decision points (a just-approved gate; ≥2 independent entities ready for the same stage; an async launch with independent FO work remaining), a live FO drive on the branch: (a) after the approval, advances + dispatches without emitting a permission-question turn; (b) fires the ≥2 independent dispatches in one turn, not serially; (c) continues the turn past the async launch while independent work remains. The independent baseline that moved the wrong way is the documented production run (the 0223 Shaping FO session, `docs/roadmap/0223-pi-dispatch-contract/`) where the origin/main contract — silent on all four patterns — false-stopped every one; the branch flips each.
- *Verified by:* a live drive (any host lane) on the scenario, reading the transcript for the three categorical signals — post-approval advance+dispatch vs a permission-question turn; ≥2 dispatches in one assistant turn vs serial single-dispatch turns; turn-continues vs turn-ends-to-wait past the async launch. Per the k74g spike (see Spike), these are categorical signals at controlled decision points, not a raw count; the residual LLM-nondeterminism is the acknowledged legitimate shape for a prose-posture change (Proof-policy note), so the scenario is constructed to force each decision point.

**AC-3 (value, behavior — S4 preserves the captain's de-risking lever) — a captain correction narrows scope without halting the session.**
In the same (or a paired) live scenario, a captain correction to ONE member's mechanism re-shapes that member while an independent member keeps being driven — proving S4 pauses only the corrected member's dispatch, not the session, and preserves the captain's interrogation as the primary de-risking lever (z25) rather than dulling it.
- *Verified by:* the live drive observing, in one session, the corrected member's dispatch pause + re-shape AND the independent member's continued progress. The baseline that moved the wrong way is the 0223 pattern-4 false-stop (a captain question halted ALL shaping).

## Proof-policy note

A Working-Principles prose addition resists a clean code gate — a contractlint "the phrase is present" check is exactly the prose-grep tautology the dev-workflow proof policy bans (the matched text was written by the same implementer the check polices). The real assurance is AC-2's live drive observing the posture, plus the structural review (AC-1) binding two independent values — the named principle and the observed behavior — that can diverge. Ideation confirms this is the legitimate shape; if a binary-enforceable property emerges (e.g. a dispatch-log structural check that parallel dispatch occurred), record it, but don't force a prose gate.

## Spike (riskiest mechanism) — no fresh spike needed; inherits k74g's

The riskiest unverified mechanism AC-2/AC-3 rest on is "can a live FO drive observe the posture as a value proof against a baseline." k74g's ideation spike already exercised this (n=91 real boot transcripts): a passive raw tool-call count is NOT a viable value proof (the greet-stop boundary is un-recoverable from a passive log and the count is LLM-nondeterministic); the viable proof is categorical behavioral signals captured at a controlled stop on a constructed scenario. That finding transfers directly — vcm's decision points (post-approval action; parallel-vs-serial dispatch; turn-continues-past-async; correction-narrows) are the same class of categorical transcript signal, observed on a scenario constructed to force them, under both the origin/main contract (the documented 0223 production baseline that false-stopped all four) and the branch. Re-running k74g's 91-transcript exercise would be the re-verify-already-proven busywork zm's own gate forbids — a self-consistent reason not to. Proven mechanisms relied on: a live interactive/async FO drive is observable; categorical decision-point signals are the metric; the 0223 debrief is the independent production baseline. Residual risk (recorded, not hidden): a single drive is LLM-nondeterministic; the scenario forcing each decision point is the mitigation, and the residual is the acknowledged legitimate shape for a prose-posture change, not a code gate.

## Test plan

- Structural review (AC-1): the four additions present + reconciled in `first-officer-shared-core.md`; run as the sprint's independent cross-member coherence review over the COMBINED z25 + zm + vcm Working-Principles section.
- Live drive (AC-2): a `pi-live` (or claude-live/codex-live) scenario on the constructed fixture (just-approved gate + ≥2 independent ready entities + an async launch with work remaining), observing post-approval-advance, parallel dispatch, no-turn-end-on-async against the origin/main baseline. This is the dogfood — the contract change drives the scenario that proves it.
- AC-3: the live scenario's captain-correction-narrows element (one member re-shaped, an independent member kept moving, in one session).
- This is a shipped-contract change — high-stakes surface. Detached adversarial audit at validation; `claude-live` + `codex-live` + `pi-live` regression (the change is to the shared core every host adapter inherits).
- **Doc-diff determination:** the changed user-visible surface is the FO contract prose (the S1–S4 wording above), which IS the deliverable. No docs-site page describes the FO's Working-Principles posture — `docs/site/get-started/first-workflow.md`'s "Keep dispatching, approving, rejecting, steering" is captain-facing session guidance already consistent with a keep-moving FO, not a description of the FO's internal turn/parallel posture. So the S1–S4 before/after IS the documentation diff; no separate `docs/` file changes.

## Related

- `skills/first-officer/references/first-officer-shared-core.md` — the file to edit (Working Principles + Clarification and Communication).
- Pi's `pi-subagents` skill (external Pi host skill, not a repo file — referenced from `skills/first-officer/references/pi-first-officer-runtime.md`) carries the "don't end your turn after launching an async child while you still have work" rule; S3 is its host-neutral lift. The exact external wording is not repo-verifiable, so S3 stands on the repo-documented pattern-3 production failure (0223), not on that quote.
- The 0223 Shaping FO debrief (`docs/roadmap/0223-pi-dispatch-contract/debrief-shaping-2026-06-19.md`) — the retrospective that surfaced the four patterns.
- `status-validate-determinism`, `state-init-precommit-hook`, `pi-launch-fnm-multishell-race`, `pi-live-shared-scenario-runners` — the independent followups this session filed sequentially (the pattern S2 addresses).

## Stage Report: ideation

- DONE: Finalize the currently-provisional Acceptance criteria section into locked, ideation-complete form
  Dropped "(provisional…)"; AC-1 relabeled "(mechanism — serves the value AC-2/AC-3 measure)" per the stage-def value-AC rule; AC-2 pinned to a constructed scenario vs the origin/main baseline (the documented 0223 false-stops — a baseline that moved the wrong way), asserting categorical posture signals per the k74g spike; AC-3 pinned to the pattern-4 baseline. Each AC carries a `Verified by:`.
- DONE: Explicitly reconcile the four strengthenings against z25 ("verification follows from what changed") and zm ("justify before climbing a mechanism rung") per the 0250 named composition tension; do not ship wording that contradicts either
  Added a "Reconciliation with z25 and zm" subsection: the three govern orthogonal axes (verification / mechanism / motion) and keep-moving is subordinate on the two it doesn't own. Guard clauses added to the SHIPPED S1 (yields to the evidence bar on merge/triage) and S2 (presupposes zm's gate; parallelizes chosen dispatches, doesn't escalate the mechanism) so the wording is non-contradictory on its face. Invariant: keep-moving changes CADENCE, never the BAR or the RUNG.
- DONE: Confirm the Scope/Out-of-scope boundaries still hold given k74g/z25/zm's current bodies
  Boundaries hold, refined: in-scope reconcile bullet now names z25/zm as same-section sibling principles + the strict-serial Wave-1 order that makes the composition safe; a leanness bullet added (terse resident form; union measured by the sprint preflight, not per-entity). Out-of-scope now states vcm does NOT edit k74g's `«engage»` block/Startup, z25's evidence-bar text, or zm's gate — it cross-references them. Doc-diff determination recorded (contract prose is the deliverable; no docs-site page describes FO posture).

### Summary

Locked the ACs, resolved the sprint's named composition tension in vcm's final wording (the reason vcm sequences last), and confirmed the scope boundaries against the current sibling bodies. Key decision: the three disciplines compose because they govern orthogonal axes — z25 = what verification a decision needs, zm = which mechanism rung, vcm = don't idle/serialize once those two have chosen — so keep-moving is written to accelerate CADENCE without ever moving past the bar (z25) or over-climbing the rung (zm); guard clauses fold this into the shipped S1/S2 text, not just the entity body. Honesty note for the gate: `pi-subagents` is Pi's external host skill, so S3's verbatim quote isn't repo-verifiable — S3 now stands on the repo-documented 0223 pattern-3 failure, and the Related line records the provenance. No fresh spike: the live-drive-observability mechanism was proven by k74g's n=91 spike and transfers; re-running it would be the busywork zm forbids.

### Wording-economy pass (post-approval, best-effort)

Per the sprint staff review flagging vcm as the heaviest single resident addition, tightened S1–S4 prose without cutting substance: merged the redundant anti-pattern sentences into em-dash clauses (S2/S4), collapsed the "(a)/(b)" exit-condition list (S3), and dropped the doubled "next action / end of turn" framing (S1). Every reconciliation guard clause is preserved verbatim in meaning — S1's evidence-bar subordination, S2's zm-gate subordination, S3's two exit conditions, S4's pause-only-the-corrected-member. Measured on the four shipped `> ` blockquotes: 2,250 → 1,984 bytes, ~266 bytes (~12%) recovered. No fat left that isn't load-bearing.
