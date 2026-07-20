---
id: 85z12f0ywkzy47akg9gwh6hm
title: "merge-guard arm phase has no keep-moving clause: armed reads as a stopping point"
status: validation
source: "FO self-diagnosis, 2026-07-08 live session. After the captain approved three validation gates and said \"push it,\" the FO ran `spacedock merge guard <slug> --verdict passed` for each entity, which only ARMS the merge (sets mod-block=merge:pr-merge) — then stopped to read the pr-merge.md hook file instead of immediately constructing and presenting the PR draft in the same turn. Before finishing even one entity's draft, the FO got pulled into an unrelated task and the arm sat untouched. When the captain later asked \"what did you do when I said push it,\" the honest answer was: armed three entities, pushed nothing."
started: 2026-07-20T03:29:40Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-merge-guard-arm-not-a-stopping-point
issue:
sprint: 0260-proportionality
group: verification
gates:
    version: 1
    current:
        gate: gate:docs-dev:85:ideation
        attempt: gate-attempt:85-ideation-2
    records:
        - id: gate:docs-dev:85:ideation
          stage: ideation
          current-attempt: gate-attempt:85-ideation-2
          attempts:
            - id: gate-attempt:85-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:85-ideation-1
                digest: sha256:843de43172a64d695f0423dc81a357fd4b3af6c3c5653623e02480ffbe4983e7
              resolution:
                type: Resolution
                id: resolution:actor-1784520391265880000
                briefing: briefing:85-ideation-1
                by: person:reviewer
                at: 2026-07-20T04:06:31Z
                decision: approve
                reason: "based on our new principle, write a estimated change, so future stage can refer to it to judge deviation"
              application:
                action: advance
                target-stage: implementation
                state: superseded
              note: "Subspace advisory float, captain at the keyboard as person:reviewer; the resolution reason directs appending a declared expected surface to the body — applied post-approval as part of the approval's own terms, not drift. Superseded by attempt 2 (captain-approved staff-review vocabulary sweep); the approval itself stands."
            - id: gate-attempt:85-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:85-ideation-1
              state: closed
              briefing:
                id: briefing:85-ideation-2-chat
                digest: sha256:cff3a819597e625058429c606ac788046316d0e0e31f72941f210e493b2394ba
                note: "chat presentation; ADVISORY digest — it hashes the working file at recording time (sweep applied, this attempt's own record excluded), which no single committed tree reproduces because an entity cannot self-bind its gates record. For drift checking, diff the entity BODY against the state commit that introduced this attempt; do not re-hash the current file."
              resolution:
                type: Resolution
                id: resolution:captain-chat-85-ideation-2
                briefing: briefing:85-ideation-2-chat
                by: person:captain
                at: 2026-07-20T10:23:28Z
                decision: approve
                reason: "Staff-review sweep, captain-approved in chat: the banned coined vocabulary removed from the group field and four body spots (plain 'captain judgment' / 'check ordering' wording; group regrouped to verification). No design change."
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "FO applied the sweep directly under the captain's edit-directly grant; codex finding 5 and fable delta finding 19."
---

**Problem:** `skills/first-officer/references/fo-merge-core.md`'s `«merge.guard»` capability describes the phase-invocation mechanics ("invoke it directly per phase; its own stdout/stderr name the FO's next action") but never states that an "armed" result is not a stopping point. Contrast `fo-dispatch-core.md`, which is explicit for stage completions: "Implementation completion is not a stopping point... The FO does not park a completed implementation and wait." The merge ceremony is exactly as sequential and stateful as the dispatch stage sequence (arm → invoke hook → finalize), but only the dispatch side carries the "keep moving" clause. The asymmetry reads as an invitation to treat "armed" as a natural pause.

**Cause:** the FO's own behavior is only as good as what the contract prescribes at each state-machine transition. `fo-dispatch-core.md` earns correct keep-moving behavior at stage boundaries because it says so explicitly, in words, at the exact transition point. `fo-merge-core.md` never says the analogous thing at the arm→hook transition, so nothing in the contract text pushed the FO to continue past "armed" into constructing and presenting the PR draft in the same turn.

**Recommended fix:** add a line to `«merge.guard»`'s effect/done-when description mirroring the dispatch-core language: an "armed" result is not a stopping point — the FO proceeds to construct and present the PR draft (or invoke whatever hook action is named) in the same turn, not a later one, exactly as a completed non-gated dispatch stage routes immediately to the next action.

## Ideation: confirmed scope

Confirm-the-spec: the seed's Problem/Cause/Recommended-fix hold. This ideation pins the target site, the exact before/after wording, the recurrence evidence, and the spike determination. Thin by design — one keep-moving clause, generic FO-contract wording, smallest edit.

**Target site:** `skills/first-officer/references/fo-merge-core.md`, the `«merge.guard»` capability. The keep-moving imperative currently lives only on the dispatch/gate side — `first-officer-shared-core.md`: "A completed non-gated, non-terminal stage is not a stopping point... advancing is the FO's next action, not the captain's." The merge ceremony's arm→hook transition is as sequential and stateful as a stage boundary but carries no such clause. This one site's absence produced the recurrence.

**Recurrence evidence (two sessions, one root cause — armed, then did not proceed):**

1. 2026-07-08 live session (seed source): captain approved three gates and said "push it"; the FO ran `merge guard --verdict passed` per entity — which only ARMS — then read the pr-merge hook file instead of opening the PR that same turn, got pulled into an unrelated task, and left three arms untouched. Honest answer to "what did you do when I said push it": armed three, pushed nothing.
2. Session 6d175b2f: after the push was already granted, the FO re-asked the captain for permission to push — twice — instead of proceeding on the armed merge.

Both treat the arm→hook transition as a stopping point (park the arm, or re-ask for a push already granted). The dispatch side already forbids the analogue — `first-officer-shared-core.md`: "want me to advance + dispatch?" is the violation — the merge side does not.

## Proposed change (before / after)

One added clause in `«merge.guard»`, immediately after the `→ shipped` bullet. No deletions; the existing effect/done-when/block/shipped bullets are unchanged.

**Expected surface:** 1 file (`skills/first-officer/references/fo-merge-core.md`), ~1 clause / ~2 lines added (the bolded paragraph plus its blank separator), no deletions, no other files touched.

**Before** (`fo-merge-core.md`, `«merge.guard»` — nothing binds the FO's turn boundary; the done-when's "left it armed/blocked with its next step named" reads as an acceptable turn end):

    - **done-when:** the entity is archived terminal, or `«merge.guard»` left it armed/blocked with its next step named in its own output.
    ...
    - → **shipped**: `` `spacedock merge guard <slug>` `` — invoke it directly per phase (via `${SPACEDOCK_BIN:-spacedock}`, per the launcher invariant above).

**After** (append this one clause after the `→ shipped` bullet):

    **An armed result is not a stopping point.** Armed is a valid return from one `«merge.guard»` invocation, not a valid end to the FO's turn: the armed result names `«hooks.run»("merge")` as the next action, and the FO invokes it in the SAME turn — opening and presenting the captain-gated PR under `merge: pr`, or running the `--no-ff` merge under `merge: local` — not a later one. Parking an armed merge, or re-asking the captain for a push already granted, is the contract violation, exactly as stopping after a completion-only stage report is (the stage-completion keep-moving clause in `first-officer-shared-core.md`). The only legitimate halts after arming are contract-defined: a `«halt.rebase-conflict»`, an in-flight `mod-block` the guard directs you to resume, or — under `merge: pr` — the captain's decision on the presented PR.

**Why not edit the done-when instead?** The done-when is correct at the capability granularity — armed/blocked ARE legitimate returns from a single invocation. The bug is the FO conflating "this invocation returned" with "my turn is over." The added clause draws exactly that distinction and leaves the (correct) done-when untouched — the smaller, more precise edit.

## Acceptance criteria

**AC-1 (value — behavior at the armed transition).** Given the armed-state scenario (captain has granted the push; `«merge.guard» <slug>` returns armed and names the merge hook as next action), the merge-core contract sanctions exactly two FO next-moves — invoke the merge hook this turn, or halt at the presented captain PR gate — and admits NO reading in which ending the turn at armed, or re-asking for the granted push, is legitimate.
- Test: captain-judgment check (adversarial read-through), not new machinery. A skeptic, given only the merge-core contract and the armed stdout, tries to justify parking at armed / re-asking permission; against the after-text it fails, against the before-text it succeeds. The "before" arm is already established by production — two real FOs parked (2026-07-08; session 6d175b2f). Cheap optional replay: hand the same armed prompt to a fresh reader-as-FO under each text and record proceed-to-hook vs end-turn.
- Baseline that can move the wrong way: the two parking incidents. If the clause is mis-placed (added as narrative but the arm→hook turn boundary left unbound), the skeptic still finds the park-license and AC-1 fails.

**AC-2 (mechanism — serves AC-1).** `«merge.guard»` in `fo-merge-core.md` carries the keep-moving clause above, attached to the armed result at the arm→hook transition and structurally mirroring the stage-completion clause in `first-officer-shared-core.md`.
- Test: the diff against `origin/main` adds only this clause — no deletions, no rewrites of existing bullets (smallest-edit + leanness constraint: net contract bytes are the added clause). Confirms the wording shipped at the right site.
- Counts only paired with AC-1, per the gate's mechanism→value re-anchor.

## Spike

No spike needed: the fix mirrors an in-force keep-moving clause (`first-officer-shared-core.md`'s stage-completion "not a stopping point") into an already-loaded contract file (`fo-merge-core.md`, the merge module the FO reads at merge time). No new mechanism, on-disk format, runtime handoff, or CLI/command-surface change — `spacedock merge guard` behavior is unchanged, and the edited contract file is itself the doc, so there is no separate docs-site diff. Enforcement is captain judgment per this sprint's check ordering (gate read-through / captain catch), not built machinery.

## Lane / merge note

The edit lands in `fo-merge-core.md`; the parallel check-ordering task (z7) edits `first-officer-shared-core.md` — different files, trivial merge. The clause cross-references the shared-core stage-completion clause by concept, not line number, so a z7-side move of that clause does not dangle the reference.

## Stage Report: ideation

- DONE: one keep-moving clause for the merge-guard armed phase (armed means proceed with the already-granted action, not a stopping point)
  Clause drafted for `«merge.guard»` in `fo-merge-core.md`, mirroring `first-officer-shared-core.md`'s stage-completion "not a stopping point"; see "Proposed change (before / after)".
- DONE: concrete before/after wording
  Verbatim `«merge.guard»` before-block quoted and the single added after-clause given; placement fixed (immediately after the `→ shipped` bullet), no deletions.
- DONE: incident evidence cited (session 6d175b2f — re-asked permission already given twice)
  Two sessions cited: 2026-07-08 "armed three, pushed nothing" (seed source) and session 6d175b2f (re-asked granted push twice) — see "Recurrence evidence".
- DONE: Record "no spike needed"; generic FO-contract wording only, no dev-specific content; smallest edit that kills the recurrence
  "## Spike" records no-spike (mirrors an in-force clause into an already-loaded file, no new mechanism); wording is generic merge-ceremony contract; one added clause, done-when left intact ("Why not edit the done-when instead?").

### Summary

Confirmed the seed spec and pinned it to a gate-ready ideation: a single keep-moving clause added to `fo-merge-core.md`'s `«merge.guard»` making an armed result not a turn-end, with verbatim before/after wording placed after the `→ shipped` bullet and no deletions. Value AC-1 measures the behavior (no contract reading licenses parking at armed or re-asking a granted push) against the two real parking incidents as baseline, proved at the captain-judgment level — no new machinery, consistent with this sprint's thesis. Edit is in a different file from the parallel check-ordering task (z7), so the merge is trivial.

## Stage Report: implementation

- DONE: The keep-moving clause lands in the «merge.guard» capability of skills/first-officer/references/fo-merge-core.md, word-for-word from the entity body's "**After**" block, placed immediately after the `→ shipped` bullet.
  Commit fb3058c6 on `spacedock-ensign/merge-guard-arm-not-a-stopping-point`; clause is `fo-merge-core.md:18`, directly after the `→ shipped` bullet (line 16) and before the `«worker.shutdown»()` paragraph. Verbatim confirmed by `diff` of entity line 99 (leading indent stripped) against shipped line 18: byte-for-byte identical, exit 0.
- DONE: `git diff origin/main` for this branch adds only that clause: zero deletions, zero rewrites of existing bullets (this is AC-2's own test, not a style preference).
  `git diff --numstat main` (this branch's base commit 972129ac) = `2	0	skills/first-officer/references/fo-merge-core.md` — one file, 2 insertions, 0 deletions. The wider `git diff origin/main` also lists 6 sprint-0260 files (.roborev.toml, docs/dev/README.md, docs/roadmap/0260-proportionality/*) — those are pre-existing commits already on local `main`, the branch base, not this branch's work; `git diff --numstat main` isolates it.
- DONE: The stage report states the ACTUAL surface (files touched, lines added/deleted) against the declared 1 file / ~2 lines added / no deletions.
  Actual: 1 file, +2 / -0 (the clause paragraph plus its blank separator). Declared: 1 file, ~2 lines added, no deletions. Exact match — 1.0x of declared surface, well inside the default 2x tolerance. No re-scope, park, or escalation needed.

### Summary

Added the one clause verbatim to `«merge.guard»` in `fo-merge-core.md` at the pinned site; nothing else touched. AC-2's mechanism test holds exactly (1 file, +2/-0, no bullet rewritten, done-when left intact per the ideation's "Why not edit the done-when instead?"). AC-1's adversarial read-through was run against the after-text: the four park-licenses a skeptic can build from the file alone — done-when's "left it armed... with its next step named", effect's "Invoke it once per phase", the mod-block guard's "on boot means a merge is mid-flight" (a crash-recovery path, not a sanction), and "captain-gated PR" read as permission-before-opening — each fails against the clause, which binds the turn boundary separately from the per-invocation done-when and closes the halt list to the contract-defined halts. Two sanctioned next-moves remain, as AC-1 requires. No new mechanism, check, or terminology introduced.

### Review findings (roborev job 322)

**ACCEPTED — final sentence too absolute, contradicted sibling contract text.** As shipped, the clause ended "The only legitimate halt after arming is the captain's decision on the PR once it is presented." Two defects: (a) `«halt.rebase-conflict»` (`first-officer-shared-core.md:129`) MANDATES halting on a rebase conflict and forbids `--force` / `--force-with-lease` / `-X ours` / `-X theirs`, so a literal reader hitting a mid-merge conflict was pushed past a mandated safety halt toward a forbidden force-push; (b) under `merge: local` there is no PR at all, so the sentence named a halt condition that can never occur — reading as "never halt" for every local-merge workflow, inconsistent with the preceding sentence's own `merge: local` branch. Narrowed to: "The only legitimate halts after arming are contract-defined: a `«halt.rebase-conflict»`, an in-flight `mod-block` the guard directs you to resume, or — under `merge: pr` — the captain's decision on the presented PR." Applied in BOTH places — the shipped clause (`fo-merge-core.md:18`) and the entity's "**After**" block (line 99, the verbatim source a later sprint member copies from) — so the two cannot drift. The parking/re-asking prohibition is untouched, verbatim. The narrowing also tightens AC-2's "structurally mirrors" test: `first-officer-shared-core.md:76` enumerates its legitimate halts in exactly this closed-list shape, which the absolute version did not.

**DECLINED — "add behavioral regression smoke tests driving the workflow through `armed`" (Medium).** Correct-but-disproportionate. Grounds: it requires new committed test machinery, which this sprint's captain ruling gates behind explicit captain approval and normally its own entity; the declared surface is 1 file / ~2 lines of contract prose, so standing up merge-workflow smoke harnesses is an order-of-magnitude scope inversion; and a cheaper check that can fail already exists under this sprint's check ordering — the dev workflow PASSES a contract/skill change only when a LIVE DRIVE observed the claimed behavior, and that live drive is this clause's proof. Promote-to-material condition: an observed live FO drive in which this clause fails to produce same-turn continuation after an armed result. Diff contribution: zero lines.

**Surface after the correction:** 1 file, +2 / −0 against the branch base — unchanged from the declared 1 file / ~2 lines added / no deletions. The narrowing edits a sentence inside the already-added paragraph, so it consumes no additional budget. Verbatim check re-run after the edit: entity line 99 (indent stripped) vs shipped line 18 diff clean, exit 0.
