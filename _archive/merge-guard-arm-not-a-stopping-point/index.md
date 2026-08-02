---
id: 85z12f0ywkzy47akg9gwh6hm
title: "merge-guard arm phase has no keep-moving clause: armed reads as a stopping point"
status: done
source: "FO self-diagnosis, 2026-07-08 live session. After the captain approved three validation gates and said \"push it,\" the FO ran `spacedock merge guard <slug> --verdict passed` for each entity, which only ARMS the merge (sets mod-block=merge:pr-merge) — then stopped to read the pr-merge.md hook file instead of immediately constructing and presenting the PR draft in the same turn. Before finishing even one entity's draft, the FO got pulled into an unrelated task and the arm sat untouched. When the captain later asked \"what did you do when I said push it,\" the honest answer was: armed three entities, pushed nothing."
started: 2026-07-20T03:29:40Z
completed: 2026-07-20T17:26:28Z
verdict: passed
score:
worktree: .worktrees/spacedock-ensign-merge-guard-arm-not-a-stopping-point
issue:
sprint: 0260-proportionality
group: verification
gates:
    version: 1
    current:
        gate: gate:docs-dev:85:validation
        attempt: gate-attempt:85-validation-2
    records:
        - id: gate:docs-dev:85:validation
          stage: validation
          current-attempt: gate-attempt:85-validation-2
          attempts:
            - id: gate-attempt:85-validation-2
              sequence: 2
              state: closed
              briefing:
                id: briefing:85-validation-2
                note: "Validation cycle 2 stage report (PASSED, zero material findings), read in full by the FO before resolving. The report IS the briefing. Cycle 1 REJECTED with three material findings; all three are resolved or dispositioned."
              resolution:
                type: Resolution
                id: resolution:fo-conn-85-validation-2
                briefing: briefing:85-validation-2
                by: agent:first-officer
                at: 2026-07-20T14:45:00Z
                decision: approve
                reason: "Approved by the FO under the captain's explicit conn grant (2026-07-20). Delegated approval, NOT a captain decision — recorded as agent:first-officer. Grounds: both shipped ACs verified against evidence the validator REPRODUCED. AC-1's blind A/B was re-run from a pre-registration written before any reader — 0/3 before-text vs 3/3 after-text surface-and-stop — and in doing so the validator found that the RECORDED probe had tested the parked 843-byte paragraph rather than the shipped clause (cmp separates them at char 435), so this is the first evidence that the text actually shipping moves the baseline. It replicates. AC-2 measured two independent ways agreeing exactly: net -24 bytes, ratchet green, full suite green uncached. The payload substitution survives its hardest test: the original AC-1 is byte-identical to a frozen ideation briefing committed nine hours BEFORE the substitution (an immutable independent source, not self-reference), and two blind readers given only the entity body each identified the shipped text, flagged the parked paragraph as unshippable, and answered NO to whether the entity delivered what it set out to — so the retreat reads as a retreat. Detached adversarial audit held on both misattribution and information-loss attacks. Zero material findings; three deferred risks and four polish items recorded with triggers and promote conditions."
              application:
                target-stage: done
                state: pending
              note: "MERGE PRECONDITION, unmet at approval: this diff touches skills/**/references/** (host-neutral FO contract), so offline plus every host lane is required. `offline` is RUN and GREEN — the cycle-1 red is fixed, and that red was caused by the parked clause's bytes. The three live lanes had NOT RUN at approval because the branch was unpushed and no PR existed; the validator stated that plainly rather than calling them green, and noted correctly that the captain's pi waiver covers a pi RED, not an UNRUN lane. FO note on its own conduct: the FO previously cited the cycle-2 probe's 3/3 result to the captain when recommending the payload substitution, without verifying which text that probe exercised. It had exercised the parked paragraph. The conclusion survived re-testing, but the evidence chain was broken and was caught only because the validation dispatch required re-running rather than accepting recorded numbers."
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
                target-stage: implementation
                state: consumed
              note: "FO applied the sweep directly under the captain's edit-directly grant; codex finding 5 and fable delta finding 19."
pr: pr-merge:537
archived: 2026-07-20T17:26:28Z
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

**PAYLOAD SUBSTITUTION, captain decision 2026-07-20.** What this entity delivers changed after validation cycle 1. The original payload — the "armed is not a stopping point" paragraph — is PARKED, NOT SHIPPED. What ships instead is the `--no-ff` conflict blocker, below. Grounds and the preserved original claim are in "## Acceptance criteria" and "### Review findings (validation cycle 1)". This is a design-reset event, not tidying: the entity now makes a narrower claim than the one it was approved on.

### What ships

Two edits in `skills/first-officer/references/fo-merge-core.md`: the "## Merge and Cleanup" paragraph — which describes the HOOK path, where the merge actually runs — gains the conflict instruction, and `«merge.guard»`'s `effect:` bullet drops a restatement of its own heading to fund it.

**Before:**

    `«merge.guard»` never invokes `«hooks.run»("merge")` and never local-merges. When armed, the FO invokes `«hooks.run»("merge")`: the `merge: local` registration performs `--no-ff merge`; `merge: pr` opens the captain-gated PR. …
    - **effect:** drive the terminal merge-finalize ceremony — auto-arm, block on an open PR, finalize on a merge sentinel, then archive (including the path-scoped archive commit) — the same under both `merge:` policies. Invoke it once per phase. …

**After** (this is the shipped text, byte-identical to `fo-merge-core.md`):

    `«merge.guard»` never invokes `«hooks.run»("merge")` and never local-merges. When armed, the FO invokes `«hooks.run»("merge")`: the `merge: local` registration performs `--no-ff merge` — on conflict, surface it and stop, never auto-resolve; `merge: pr` opens the captain-gated PR. …
    - **effect:** drive the terminal merge-finalize ceremony (including the path-scoped archive commit), the same under both `merge:` policies. Invoke it once per phase. …

**Attribution matters here and was gotten wrong once.** The instruction first shipped in `«merge.guard»`'s `block:` bullet, which describes when the GUARD refuses. But the same file states that `«merge.guard»` never invokes the merge hook and never local-merges — the `merge: local` registration performs the `--no-ff merge`. A conflict is therefore something the FO hits on the hook path, not something the guard blocks on. It now sits on the sentence that describes that path. See "Review findings (roborev, substituted payload)".

**Surface:** 1 file, +2/−2 lines, **net −24 bytes** of FO prompt surface (122231 → 122207 against a 122634 ceiling) — the change returns headroom rather than spending it. The `effect:` trim is genuine: the heading directly above it already reads "auto-arm → block-on-open-PR → finalize-on-merge-sentinel, then archive", so the bullet restated it verbatim in substance.

### PARKED — NOT SHIPPED — do not carry this forward

The paragraph below is the original payload, in its M2/D1-repaired form (843 bytes). It is recorded for history only. It is NOT in `fo-merge-core.md` and MUST NOT be copied into any contract file. A later member copying from this entity must copy the "What ships" block above, never this one.

    **An armed result is not a stopping point.** Armed is a valid return from one invocation, not a valid end to the FO's turn: invoke `«hooks.run»("merge")` in the SAME turn, not a later one. Parking an armed merge, or re-asking the captain for a push already granted, is the contract violation, exactly as stopping after a completion-only stage report is (the stage-completion clause in `first-officer-shared-core.md`). The only halts the contract licenses after arming are: a captain decision the contract requires — under `merge: pr`, the decision on the presented PR; the guard's blocked result, which waits for the merge sentinel; or an unresolved blocker — a `«halt.rebase-conflict»`, a `--no-ff` merge conflict, an unmet clarification. On a blocker, surface it and stop; never force, never auto-resolve, never discard either side.

**Why not edit the done-when instead?** (Retained from the original design, and still the reason the parked paragraph did not touch it.) The done-when is correct at the capability granularity — armed/blocked ARE legitimate returns from a single invocation. The bug is the FO conflating "this invocation returned" with "my turn is over."

## Acceptance criteria

These describe WHAT SHIPS. They are NARROWER than the criteria this entity was approved on. The original criteria are preserved verbatim below under "Original value claim — NOT DELIVERED"; they were not edited down into these, because they are a different claim.

**AC-1 (value — behavior on a `--no-ff` conflict after arming).** Given a `merge: local` workflow where the ceremony's own mandated `git merge --no-ff` hits a content conflict, the merge-core contract directs the FO to surface the conflict and stop, and admits no reading in which the FO auto-resolves it or loops back into `«merge.guard»` without reaching the captain.
- Test: blind A/B read-through, baseline established by production-shaped readers, no new machinery. RUN in implementation cycle 2 — 3 readers per arm, pre-registered coding rule. Before-text 0/3 surfaced to the captain (3/3 looped back to `«merge.guard»`, citing the boot-time Mod-Block Guard clause); after-text 3/3 surfaced and stopped. Baseline moved.
- Known limit, recorded rather than hidden: the same probe returned a NULL on the force-restraint dimension — 6/6 readers across both arms already refused to auto-resolve. The restraint half of the bullet is therefore belt-and-braces, not evidence-backed.

**AC-2 (mechanism — serves AC-1).** `fo-merge-core.md` directs the FO to surface a `--no-ff` merge conflict and stop, attached to the HOOK path that performs the merge (not to `«merge.guard»`, which never local-merges), and the change pays its own way on the FO prompt-surface ratchet.
- Test: `go test ./...` green including `TestFOFunctionPromptSurfaceShrinks`, with the FO surface total recorded before and after. Net cost must be ≤ the 81-byte `effect:` trim that funds it.

### Original value claim — NOT DELIVERED

Preserved because it is the reason this entity exists, and it is still unmet. It was NOT narrowed into AC-1 above; it was abandoned as this entity's payload by captain decision on 2026-07-20.

> **AC-1 (value — behavior at the armed transition).** Given the armed-state scenario (captain has granted the push; `«merge.guard» <slug>` returns armed and names the merge hook as next action), the merge-core contract sanctions exactly two FO next-moves — invoke the merge hook this turn, or halt at the presented captain PR gate — and admits NO reading in which ending the turn at armed, or re-asking for the granted push, is legitimate.

**Why it was not delivered — two grounds, both evidential.**
1. **No probe moved its baseline.** Validation's blind A/B on the arm→hook transition: 3/3 before-text readers already chose proceed-this-turn at HIGH confidence, citing text already in force. The re-anchored probe in cycle 2 did not rehabilitate it either — what moved there was the conflict-blocker claim, a different claim.
2. **It could not fund itself.** The paragraph costs 843 bytes against 203 bytes of headroom on a shrink-ratcheted surface; paying its own way needed 844 bytes of offsetting trim, and only 81 were defensible without cutting unrelated contracts (one of them a file z7 is concurrently editing).

**The null is a limit of the probe as much as a verdict on the clause.** Three careful readers asked "what do you do next?" answer correctly. Both production failures happened to real FOs deep in a session, under interruption and competing demands. A careful-reader probe structurally cannot detect a context-pressure failure — it presents the decision in isolation, which is exactly the condition under which the decision is easy. The null is therefore weak evidence that the problem is not real, and should not be read as proof of it.

**The two production incidents remain UNADDRESSED.** Nothing that ships here prevents either:
1. 2026-07-08 — captain said "push it"; the FO armed three entities, was pulled into an unrelated task, and pushed nothing.
2. Session 6d175b2f — the FO re-asked the captain for a push already granted, twice.

**Follow-up:** filed separately as `probe-armed-parking-under-context-pressure`, carrying this entity's probe design plus the requirement that a future probe model context pressure rather than careful reading.

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

## Stage Report: validation

- DONE: Every AC-N is verified with evidence you REPRODUCE yourself, not evidence you quote: re-run AC-2's mechanism test independently (1 file, +2/-0, no existing bullet rewritten) and re-anchor AC-1 on its end value — the clause must actually foreclose the park-license, not merely be present in the file.
  AC-2 PASSES, reproduced: `git diff --numstat $(git merge-base HEAD main)` = `2 0 skills/first-officer/references/fo-merge-core.md`; no existing bullet touched; entity line 99 (indent stripped) and shipped `fo-merge-core.md:18` are byte-identical (sha256 `b8f900a7…`, 837 chars / 851 UTF-8 bytes), one occurrence each. AC-1 NOT ESTABLISHED — see M3.
- FAILED: Determine the required CI lanes from what the diff touches per docs/dev/README.md:78 … name the lanes the path->lane mapping requires, and run them.
  Mapping: the diff touches `skills/**/references/**` (host-neutral FO contract) → `offline` + every host lane (`claude-live`, `codex-live`, `pi-live`). `offline` RUN and RED (M1) — deterministic, reproduced 2/2, not a flake, so not re-runnable to green. The three live lanes HAVE NOT RUN: the branch is not pushed (`git ls-remote --heads origin …` empty), no PR exists, and the live jobs are environment-gated on maintainer approval. `pi-live` red is WAIVED per the captain waiver in force; that waiver covers pi only and does NOT cover M1, which is the secret-free offline lane.
- DONE: Detached adversarial audit on a THROWAWAY checkout (never this worktree)
  Ran in `scratchpad/audit/` on `git show`-extracted copies, never the worktree. The park attack FAILS (foreclosed in terms; the `done-when` and "captain-gated" escapes are both rebutted). The rebase-conflict attack FAILS as framed — `«halt.rebase-conflict»` is preserved as list item one. But the audit surfaced an adjacent hole the framing missed: M2.

### Material findings (block the gate)

**M1 — outcome defect: the offline CI lane is red, caused exactly by this diff.** `go test ./...` (the literal command of the CI `offline` job) fails `TestFOFunctionPromptSurfaceShrinks`: FO prompt surface = 123283 bytes, must be strictly below the post-#531 baseline 122634. The branch base 972129ac measures 122430 and PASSES. Delta = 853 bytes = the clause's 851 UTF-8 bytes + its newline + the blank separator — the entire regression is this entity's +2 lines, nothing else. Base headroom was 204 bytes; the clause needs 853. This is a scope collision, not a typo: the FO prompt surface is under an active *shrink* ratchet and this entity adds to it. Ideation's "## Spike" ("no new mechanism… no separate docs-site diff") never considered the byte ratchet, and the implementation stage report verified surface with `git diff --numstat` only — it never ran the suite. Repairing this in place is not validation's call: the captain must decide whether the clause earns 853 bytes on a ratcheted surface, and if so whether it is paid by an offsetting trim or by moving the baseline.

**M2 — outcome defect: the closed halt list under-covers the operation the same sentence mandates.** The clause orders "running the `--no-ff` merge under `merge: local`", then declares "The only legitimate halts after arming are contract-defined: a `«halt.rebase-conflict»`, …" — exhaustive in form. Verified against the text: `«halt.rebase-conflict»` fires on exactly three triggers (`state ready` exit 3, `«state.commit»` exit 3, a manual FO-held `pull --rebase` CONFLICT — `first-officer-shared-core.md:129`). A `git merge --no-ff` conflict is none of the three. So on the `merge: local` path a merge conflict falls outside every licensed halt, while the restraining ban ("Never `--force`/`--force-with-lease`, never `-X ours`/`-X theirs`, never discard either side — do not force-push or auto-resolve") sits *inside* `«halt.rebase-conflict»`'s `block:` bullet, scoped to the halt that never fires. Corroborated live, not just on paper: the open-ended after-text drive named that exact sentence as its governing sentence and read it as enumerating the halts "exhaustively." This is the same defect class as the ACCEPTED roborev finding — the narrowing closed the rebase and no-PR-under-local cases and left this one open.

**M3 — evidence defect: AC-1's value claim has no evidence that can fail.** AC-1 names a baseline that can move the wrong way ("against the after-text it fails, against the before-text it succeeds") and an optional replay. I ran the replay as a blind A/B on the throwaway checkout — before-text vs after-text, identical scenario, real armed stdout taken from `internal/status/merge.go:578`. Result: 3/3 before-text readers chose proceed-this-turn at HIGH confidence and answered "yes" to whether parking would violate the contract as written, citing text already in force (`fo-merge-core.md`'s "When armed, the FO invokes `«hooks.run»("merge")`" plus shared-core's keep-moving clause). The baseline did not move. One before-text reader independently spotted and resolved the `done-when` trap the implementation report credits the clause with closing. The clause is NOT shown to be useless — the same reader noted the conclusion is *derived* rather than *stated*, and the two production incidents are real — but AC-1's designated probe cannot distinguish the clause's contribution from the status quo, so the value AC is unproven either way.

### Deferred risks (recorded, not blocking)

- **D1 — the parity claim is false in content.** The clause asserts it works "exactly as" shared-core's stage-completion clause, but that clause's halt list (`first-officer-shared-core.md:76`) also admits "an unmet clarification" and a general "a captain decision the contract requires"; the shipped list admits neither, so it is strictly stricter than the clause it claims parity with. Trigger: a captain correction or unmet clarification arising after arming. Promote to material on an observed drive where an FO suppresses either because the list excluded it. Fold into M2's rewrite — it is the same sentence.
- **D2 — "an in-flight `mod-block` the guard directs you to resume" is circular at the point of use.** Arming itself sets `mod-block=merge:pr-merge` and says so in stdout; the disambiguating qualifiers ("on boot", "check what the hook left") live in a different section. The audit rated this PARTIAL — it dies on a careful reading. Trigger: a careless reader at arm time. Promote to material on an observed drive where an FO cites this list item to park at armed.

### Polish

- The frozen ideation briefing (`review/ideation/briefing-1.md:56`) still carries the pre-narrowing absolute sentence. Correct as an immutable record — flagged only because the later sprint member that copies this clause verbatim must copy from the entity body (line 99), never from the briefing.

### Recommendation

**REJECTED.** The roborev decline is properly recorded (three grounds, a promote-to-material condition, zero-line diff) and AC-2 holds exactly — but the offline lane is red and attributable to this diff alone, the shipped halt list leaves a `merge: local` conflict with no licensed halt and no force ban, and the value AC's own baseline does not move. M1 is a scope question for the captain (853 bytes onto a shrinking surface), not another blind implementation cycle; M2 and D1 are one sentence to rewrite; M3 needs either a probe that can actually move or an honest re-anchor of AC-1.

### Summary

Validated the one-clause merge-guard change and found three material problems. Reproduced AC-2 independently — it holds exactly (1 file, +2/-0, entity and shipped text byte-identical). Ran the offline lane and it fails `TestFOFunctionPromptSurfaceShrinks` by 649 bytes, with the entire 853-byte regression traced to this clause against a green base; the FO prompt surface is under a shrink ratchet this entity was never scoped against. The detached adversarial audit confirmed the park is foreclosed but found the closed halt list does not cover the `--no-ff merge` conflict the same sentence mandates under `merge: local`, leaving that path with no licensed halt and the anti-force ban scoped to a halt that never fires. The blind A/B replay of AC-1 showed no baseline movement — 3/3 before-text readers already proceed same-turn at high confidence — so the value AC is unproven rather than disproven. The three required live lanes have not run; the pi waiver is cited and does not cover the offline failure.

### Feedback Cycles

- Cycle 2 (2026-07-20, third roborev round) — **decision: STOP THE REPAIR CYCLE. Residual finding declined here and routed to follow-up.**

  **The finding, for the third time.** Roborev job 322, its second review, and now its third have each independently flagged the same gap: `«merge.guard»`'s `done-when` bullet ("or `«merge.guard»` left it armed/blocked with its next step named in its own output") licenses terminating at armed, so an FO may stop before running the merge hook. Each read was cold, and the third one sharpens the earlier two: the permission lives in the `done-when` clause specifically, not in the section generally.

  **Why it is not repaired here.** Three grounds, in order of weight.
  1. This is the captain's parked payload. The paragraph making exactly this claim was parked on 2026-07-20 on two recorded grounds, and neither has been overturned: no probe moved its baseline (3/3 before-text readers already proceed correctly, twice measured), and it cannot fund itself — 843 bytes against a ceiling that, even after this entity returned 24 bytes and z7 returned 161, leaves roughly 587. The full paragraph still does not fit.
  2. Continuing would BE the runaway repair loop this sprint exists to stop. This entity has now run three implementation cycles and three review rounds on a two-line contract change. The forensic record's four HIGH-severity incidents were all contract-legal one round at a time; each individual round looked justified. Recognising the shape and halting is the behavior 0260 is buying, and an FO that repairs a fourth time because the reviewer asked a third time has learned nothing from the sprint it is driving.
  3. The evidence is genuinely split and the split is informative, not a stalemate to break by fiat. Three adversarial TEXTUAL reads say the wording permits the wrong action; two behavioral probes say careful readers do not take it. Both are true. The reconciliation — that the failure needs context pressure to appear, which no probe so far reproduces — is a research question with an entity already filed for it.

  **Where it goes instead.** `probe-armed-parking-under-context-pressure` (hjb4z38k1f70yj1psb3yeaap). The three cold roborev reads are corroborating evidence for that task and materially strengthen the "our instrument was too weak" reading over "the claim is false". Recorded there by name rather than left in a review comment.

  **What was NOT done, deliberately.** An affordable miniature of the parked clause (a `done-when` amendment of roughly 80-120 bytes) is technically within budget now that this entity runs at -24. It was not taken: re-introducing a captain-parked payload in smaller type, on the same null behavioral evidence, is the decision-by-attrition this sprint prices — the reviewer asking three times is not new evidence about FO behavior. Surfaced to the captain as a live option rather than actioned or buried; if the captain wants it, it is one small cycle, and the option is recorded here so the choice stays open and visible.

  **Promote-to-material condition.** An observed live drive in which an FO terminates its turn at an armed result, citing the `done-when` bullet. That is the evidence two probes could not produce, and it would reopen this immediately.



- Cycle 1 (2026-07-20, validation REJECTED, routed to implementation) — **captain decision: RE-ANCHOR AC-1, then decide.**

  **Findings routed.** M1 (outcome): the offline lane reds on `TestFOFunctionPromptSurfaceShrinks` — FO prompt surface 123283 bytes against a 122634 ceiling, with the entire 853-byte regression attributable to this clause alone; the branch base measures 122430 and passes, so headroom was 203 bytes and the clause needs 853. M2 (outcome): the narrowed halt list is exhaustive in form but does not cover a `git merge --no-ff` conflict, which the same sentence mandates under `merge: local` — `«halt.rebase-conflict»` fires on exactly three triggers and that is none of them, so the path has no licensed halt while the anti-force ban sits inside a halt that never fires. M3 (evidence): a blind A/B replay of AC-1's own probe showed no baseline movement — 3/3 before-text readers chose proceed-this-turn at high confidence, citing text already in force — so the value AC is unproven rather than disproven. D1 and D2 recorded as deferred risks; D1 folds into M2's rewrite.

  **Provenance of M2, recorded against this FO.** The absolute sentence M2 attacks is the one the FO directed in the pre-validation roborev round, to close the rebase-conflict and no-PR-under-local cases. That narrowing was correct on both counts and opened a third gap. Recorded here rather than quietly re-fixed: an FO-directed correction is not exempt from the evidence bar it imposes.

  **Provenance of M1, recorded against this FO.** The implementation dispatch required surface verification by `git diff --numstat` and never required `go test ./...`, on the unstated assumption that a prose-only change could not break a test. The contract's own rule is that required verification follows from what changed, not the FO's sense of relevance. Every subsequent validation dispatch this sprint now names the offline lane explicitly.

  **Sprint-wide consequence (captain ruling, 2026-07-20): EACH MEMBER PAYS ITS OWN WAY.** The ratchet caps 13 first-officer contract files and the sprint plans roughly +5,500 bytes into them across this member, z7, bw, and 02av. Every ratcheted member must land an offsetting trim in the measured files at least equal to what it adds. The sprint index's stated mitigation — prefer lazy-loaded references over boot-resident lines — does NOT satisfy this check, because the measured set includes the deferred cores (`fo-dispatch-core.md`, `claude-fo-dispatch.md`) as well as the boot-resident ones. Raising the baseline was considered and rejected: editing a check so the change passes is the anti-pattern this sprint exists to remove.

  **Assignment.** Re-anchor AC-1 on a probe whose baseline can actually move; fix M2 and fold D1; fund the clause's bytes with an offsetting trim in the measured set. Whether the clause ships at all is a captain decision deferred to the re-anchored evidence — the two production incidents behind it are real, and "the probe was the wrong instrument" and "the clause adds nothing" are both still live readings.

## Stage Report: implementation (cycle 2)

- DONE: AC-1 is re-anchored on a probe whose baseline can actually move, the probe is RUN, and its real result is recorded — including an honest null if that is the outcome.
  Blind A/B, 3 fresh readers per arm, pre-registered prediction and coding rule written before any reader ran. Result is SPLIT: the primary outcome moved cleanly, the secondary outcome is a null. Detail below.
- DONE: M2 is fixed and D1 folded into the same sentence: the halt list covers the `merge: local` --no-ff conflict path it mandates, and the anti-force ban is no longer scoped to a halt that never fires there.
  Repaired text is in the "**After**" block above (843 bytes). The halt list now names "a `--no-ff` merge conflict" and carries its own inline restraint ("On a blocker, surface it and stop; never force, never auto-resolve, never discard either side"), so the ban no longer depends on `«halt.rebase-conflict»` firing. D1 folded: the list now admits "a captain decision the contract requires" and "an unmet clarification", matching `first-officer-shared-core.md:76`, so the parity claim is true in content. NOT SHIPPED — see the byte finding.
- DONE: `go test ./...` passes, TestFOFunctionPromptSurfaceShrinks included, funded by a genuine offsetting trim in the measured file set with its byte accounting recorded — or a recorded finding that the clause cannot afford itself.
  Second branch taken: RECORDED FINDING — the clause cannot afford itself. The paragraph is parked (branch restored to the pre-clause file); `go test ./...` is green, 17 packages, zero FAIL, FO surface 122430 = the base exactly.

### Review findings (validation cycle 1)

**M3 — AC-1 re-anchored. Probe RUN. Primary moved; secondary is a null.**
Probe site: the `merge: local` `--no-ff` conflict path the validator identified, chosen because the contract there is genuinely ambiguous (no licensed halt; the ban scoped to a halt that cannot fire). Arms differ only by the clause; both arms also got the shared-core stage-completion clause and the full `«halt.rebase-conflict»` section. Pre-registration, fixtures, and the coding rule are in the throwaway probe dir, never the worktree.
- PRIMARY (next action / is ending the turn licensed): BEFORE 0/3 halt-and-surface — 3/3 instead re-ran `«merge.guard»`, all three citing the **Mod-Block Guard** clause as governing. AFTER 3/3 halt-and-surface, all three citing the named `--no-ff` blocker. Clean 3/3 separation: the baseline moved.
- SECONDARY (the force ban): NULL. 6/6 readers — both arms — refused to self-resolve and cited the ban. BEFORE readers imported `«halt.rebase-conflict»`'s ban by analogy even though it does not fire on this path.
- What this means, stated against my own interest: M2's force-risk premise is NOT behaviorally supported — no reader was tempted to force or auto-resolve, so the gap is real on paper but did not produce the predicted behavior. The clause's demonstrated effect is narrower than the entity claims: it routes a `--no-ff` conflict to the captain instead of back into the guard loop. It is NOT shown to prevent forcing. The original "armed is not a stopping point" payload remains unsupported by two independent probes now.
- Side finding: BEFORE readers using the boot-time Mod-Block Guard clause as a mid-turn rule is D2's circularity showing up empirically.
- Caveats: n=3 per arm, one model, one scenario, readers are not production FOs. The probe was not re-run or retuned after seeing the result.

**M1 — THE BYTE RATCHET. Recorded finding: this clause cannot afford itself.**
Accounting, all measured with `TestFOFunctionReferenceCheckpointMetrics`:
- base 122430; ceiling strictly below 122634, so headroom is 203 bytes.
- original shipped clause +853 -> 123283 (validation's number, reproduced).
- M2/D1-repaired clause 843 bytes -> 123274. The repair is 9 bytes cheaper than what it replaced; it does not change the problem.
- To merely pass: shed 641. To pay its own way (net <= 0 per the captain ruling): shed 844.
- Trim actually defensible: 81 bytes — `«merge.guard»`'s `effect:` bullet restates its own heading verbatim in substance ("auto-arm, block on an open PR, finalize on a merge sentinel, then archive"). In my own file.
- Trim found but REJECTED: ~145 bytes in `first-officer-shared-core.md:78`, which restates the reuse rule its own sentence says is already loaded. Rejected because z7 is concurrently editing that file and this entity's own lane note exists to keep the two out of it — taking it would manufacture the merge conflict the lane note prevents.
- Remainder (~400 bytes of near-duplicate prose across `claude-fo-dispatch.md`, `claude-first-officer-runtime.md`, and the legacy team skill) are judgment calls in subsystems this entity has no standing in. Cutting them to fund a merge clause is bytes-for-bytes, not a genuine reduction — the thing the ruling forbids.
- 844 needed, 81 defensible without collateral edits. The paragraph form does not fit and cannot be made to fit honestly, so it is parked rather than shipped red.

**Verified alternative, for the captain's decision (NOT shipped).** Extending the existing `block:` bullet with "A `--no-ff` merge conflict is a blocker: surface it and stop, never auto-resolve." (+82 bytes) plus the 81-byte heading-restatement trim measures **122431 — net +1 byte, contractlint green**, empirically run, not estimated. It delivers exactly the payload the re-anchored probe validated and drops the payload two probes could not support. It is not shipped because it redefines an approved deliverable, which is a gate decision, not mine.

**Recommendation.** Ship the verified alternative above and retire the paragraph, or park this member entirely. Shipping the paragraph requires either 844 bytes of collateral trim in unrelated contracts or moving the baseline, and the ruling forbids the second.

### Summary

Re-anchored AC-1 on the `merge: local` conflict path and ran it blind, 3 readers per arm, pre-registered. The primary outcome moved 3/3 vs 3/3 — before-text readers loop back to `«merge.guard»`, after-text readers surface to the captain — but the force-ban outcome is a null: both arms already refuse to self-resolve, so M2's force-risk premise is not behaviorally supported and the clause's real effect is narrower than claimed. Repaired M2 and folded D1 into the same sentence (the halt list now names the `--no-ff` conflict, carries its own restraint, and matches shared-core's list). On M1 I took the second branch the captain left open: 844 bytes are needed to pay the clause's way, 81 are defensible in my own file, and the ~545 remaining would mean cutting unrelated contracts — including one z7 is concurrently editing — so the paragraph is parked and the branch is green at exactly the base. A proportionate alternative measuring net +1 byte is verified green and recorded for the ship decision.

## Stage Report: implementation (cycle 3)

- DONE: Ship the verified alternative — `block:` bullet gains the `--no-ff` conflict blocker, `effect:` bullet drops the restatement of its own heading.
  Branch rebased onto current main (bdf39f01) FIRST, then applied and re-measured — the earlier 122431 figure was against the old base and is superseded. New base 122231, shipped 122232: **net +1 byte**, ceiling 122634. `git diff --numstat main` = `2 2 skills/first-officer/references/fo-merge-core.md`. Shipped exactly what was verified; not extended.
- DONE: `go test ./...` green including TestFOFunctionPromptSurfaceShrinks.
  17 packages ok, zero FAIL. `TestFOFunctionPromptSurfaceShrinks` PASS. FO surface recorded before (122231) and after (122232) via `TestFOFunctionReferenceCheckpointMetrics`.
- DONE: Record the payload substitution as a decision, preserving the original value claim.
  "## Proposed change" now opens with the captain decision and date, carries the shipped text under "What ships", and quarantines the paragraph under "### PARKED — NOT SHIPPED — do not carry this forward". "## Acceptance criteria" states up front that the new ACs are narrower, then preserves the original AC-1 verbatim as a block quote under "Original value claim — NOT DELIVERED" with both grounds, the probe-limit caveat, the two still-live incidents, and the follow-up name.
- DONE: Keep the entity's "**After**" block byte-identical to the shipped text.
  Verified by comparison against `fo-merge-core.md`: the `block:` bullet in the After block is byte-identical to the shipped line (exit-0 equality check, 194 bytes). The `effect:` bullet is shown with a trailing `…` elision of the unchanged tail; its quoted prefix matches the shipped line exactly.

### Surface: actual against declared

Declared at ideation: 1 file (`fo-merge-core.md`), ~2 lines added, no deletions. Actual: 1 file, +2/−2 lines, net +1 byte. Same file, same capability, but a DIFFERENT EDIT than declared — the declared payload is parked and a narrower one shipped in its place. That is captain-directed scope (decision 2026-07-20 on the cycle-2 recorded finding), not drift. The line count now shows deletions where the declaration said none; both are the `effect:`-bullet trim that funds the change, not a rewrite of the payload.

### Summary

Shipped the alternative the captain selected, exactly as verified and no wider: the `«merge.guard»` `block:` bullet now names a `--no-ff` merge conflict as a blocker to surface and stop on, funded by trimming the `effect:` bullet's verbatim restatement of the heading directly above it. Rebased onto current main before measuring, because the base had moved from 122430 to 122231 and the old figure no longer applied; the shipped state measures net +1 byte with the full suite green. The originally approved payload is parked and recorded as a design-reset event rather than tidied away: the original AC-1 is preserved verbatim as NOT DELIVERED, with both grounds (no probe moved its baseline; it could not fund its 844 bytes), the explicit caveat that a careful-reader probe structurally cannot detect the context-pressure failure both production incidents actually exhibited, and the two incidents recorded as still unaddressed. Follow-up work is referenced as `probe-armed-parking-under-context-pressure`.

## Stage Report: implementation (cycle 4)

- DONE: ACCEPTED — re-home the conflict clause from `«merge.guard»`'s `block:` bullet to the hook path that performs the merge.
  The finding was right and it was my line. `fo-merge-core.md:7` states that `«merge.guard»` never invokes `«hooks.run»("merge")` and never local-merges; the `merge: local` registration performs the `--no-ff merge`. A conflict there is not something the guard blocks on, so the `block:` bullet — which describes when the guard refuses — was the wrong owner. The instruction now rides the sentence that describes the hook path, unchanged in force: "on conflict, surface it and stop, never auto-resolve". Commit on branch; shipped text mirrored byte-identically into "### What ships".
- DONE: Stay byte-neutral or better; re-measure and confirm the ratchet.
  Better. Base 122231, shipped **122207 — net −24 bytes**, so the change now returns headroom instead of spending it (the re-homed wording is 57 bytes against the 81-byte `effect:` trim that funds it; the previous placement cost 82). `TestFOFunctionPromptSurfaceShrinks` PASS. `go test ./...` green: 17 packages, zero FAIL. Main had not moved (still bdf39f01), so no rebase was needed; confirmed rather than assumed.
- DONE: Record the two declines.
  Both below, with grounds, promote conditions, and zero-line diffs.

### Review findings (roborev, substituted payload)

**ACCEPTED — misattribution.** See the cycle-4 checklist item above. Precision defect, not behavioral: the 3/3 after-text probe result stands, because what readers acted on was the instruction to surface and stop, not the bullet it hung from. But a reader taking the old placement at face value would conclude the guard handles the conflict, and it demonstrably does not.

**DECLINED 1 — "the contract allows an armed guard result to appear to be a stopping condition; a normal local merge could therefore halt after arming."** This is the parked payload restated. Not reinstated: the parking decision stands on both recorded grounds, and neither has changed — no probe moved its baseline, and it cannot fund its 843 bytes. Diff contribution: zero lines.

  **But this is new evidence and it belongs on the record.** Roborev reached that conclusion INDEPENDENTLY and cold, reading only the current contract text, with no knowledge that the paragraph existed or had been parked. That is a third independent read finding the gap real — after the seed's two production incidents.

  What it corroborates, precisely, because the distinction is the point: the probes asked readers to ACT ("what do you do next?") and 3/3 acted correctly — behavioral evidence, and it was NULL. Roborev asked whether the TEXT is ambiguous and found that it is — textual evidence, and it is POSITIVE. Both are true at once. Together they support the reading that **the instrument was too weak**, not that the claim is false: a careful reader resolves the ambiguity correctly, so a careful-reader probe cannot surface it as a behavior failure. The failure needs context pressure to appear, which is what the two production incidents had and no probe so far has reproduced.

  This is corroborating evidence for `probe-armed-parking-under-context-pressure` and should be carried into it. It is NOT grounds to reverse the park — it strengthens the case that the problem is real while leaving both parking grounds untouched.

**DECLINED 2 — "add behavioral smoke tests for armed continuation and conflict termination."** Same disposition as the identical proposal in roborev's first review of this branch, on the same grounds: a new committed check needs explicit captain approval and normally its own entity; standing up merge-workflow smoke harnesses for a two-line contract edit is the scope inversion this sprint prices; and the cheapest check that can fail already exists — the dev workflow requires a live drive before a contract change is marked PASSED. Promote to material on an observed live drive where the shipped instruction fails to produce surface-and-stop on a `--no-ff` conflict. Diff contribution: zero lines.

### Summary

Fixed the one accepted finding: the conflict instruction was attributed to `«merge.guard»`'s `block:` bullet, but the guard never performs the merge, so it now rides the sentence describing the hook path that does. Force of the instruction is unchanged and the probe result is unaffected. The re-homing came in cheaper than the placement it replaced — the entity now measures net −24 bytes against its base, returning headroom rather than consuming it, with the suite green and the ratchet passing. Both declines are recorded with grounds and zero-line diffs. The lead decline matters beyond bookkeeping: roborev independently and cold reached the same conclusion as the parked paragraph, which is textual corroboration that the gap is real even though the behavioral probes returned null — evidence that the instrument was too weak rather than that the claim was false, carried by name to `probe-armed-parking-under-context-pressure` and explicitly not treated as grounds to un-park.

## Stage Report: validation (cycle 2)

- DONE: Each shipped AC verified with evidence you REPRODUCE. AC-1 is a behavioral claim about the `merge: local` --no-ff conflict path: re-run the blind A/B yourself rather than accepting the recorded 0/3 vs 3/3, and report your own numbers.
  Re-ran it, pre-registered before any reader (`probe-v2/PRE-REGISTRATION.md`, throwaway dir). **My numbers: BEFORE 0/3 surface-and-stop (2/3 hand-resolve, 1/3 loop back), AFTER 3/3 surface-and-stop, 3/3 HIGH confidence.** Baseline moved; prediction met. Material discovery: the recorded 3/3-vs-0/3 was run against the PARKED 843-byte paragraph, not the shipped text — `cmp` puts cycle 2's arm B and the shipped file apart at char 435. This re-run is the first test of the text that actually ships, and it replicates.
- DONE: AC-2 is the byte claim: re-measure the FO prompt surface and confirm net <= 0 against main with the ratchet green.
  Measured two independent ways, agreeing exactly: my own `git show | wc -c` sum over the 13 ratcheted paths gives main 122231 / HEAD 122207; `TestFOFunctionReferenceCheckpointMetrics` (`-count=1`, uncached) reports 122207. **Net −24, ceiling 122634, ratchet PASS.** Full offline lane re-run uncached: `go test ./... -count=1` → 17 packages ok, zero FAIL.
- DONE: The payload substitution is honestly recorded and cannot mislead a later reader.
  Strongest evidence is independent-source, not self-referential: entity line 130 vs the frozen `review/ideation/briefing-1.md:62` is a byte-for-byte `diff` match, 447 bytes, exit 0. The briefing was committed 2026-07-20 12:45 (2c9f205d) and never touched again; the substitution landed 22:04 — so the original claim is genuinely preserved, not rewritten. Corroborated behaviorally: 2 fresh readers given only the entity body both picked the shipped text, both flagged the parked paragraph as NOT shipped citing the quarantine heading, both answered NO to "did it deliver everything," and both named the two production incidents as still unaddressed. The substitution reads as a recorded retreat, not a clean pass.
- DONE: The conflict blocker is attributed to the actor that actually performs the merge (the hook path), not to «merge.guard». Confirm the placement is correct and the instruction still reads as binding where it now sits.
  Detached adversarial audit (throwaway `scratchpad/audit/`, `git show`-extracted copies, never the worktree) ran the misattribution attack and returned HELD: the sentence opens "«merge.guard» never invokes «hooks.run»("merge") and never local-merges", and "surface it and stop" is a bare imperative binding the document's addressee, whose nearest finite governor is "the FO invokes". Confirmed behaviorally rather than by reading: 3/3 AFTER readers quoted that sentence as governing and acted as the FO; none concluded the guard handles the conflict. Binding force intact.
- DONE: Detached adversarial audit on a THROWAWAY checkout.
  Ran on extracted copies, never the worktree. Misattribution attack HELD; information-loss attack HELD (all four trimmed concepts survive verbatim in the section heading, verified directly). Two attacks landed and are recorded as D3/D4 below.
- DONE: Required lanes are a function of the diff. Derive and state the mapping.
  Diff touches `skills/**/references/**` — a host-neutral FO contract — so per `docs/dev/README.md:78` the mapping is `offline` + every host lane (`claude-live`, `codex-live`, `pi-live`). **`offline` RUN and GREEN** (the cycle-1 red is fixed; the regression was the parked clause's bytes). **The three live lanes have NOT RUN**: `git ls-remote --heads origin` is empty and no PR exists, so CI has had nothing to run. This is an unmet merge precondition, stated plainly and not called green. The pi waiver covers a pi *red*, not an unrun lane.
- DONE: Classify on both axes; only material findings block. Record deferred risks with exact trigger and promote condition.
  Below. No material finding survives.

### Evidence classification (which of my evidence is which)

- **Existence facts — prove only that we wrote it:** the `+2/−2` diff shape; the prefix comparisons; the heading-contains-concept check. Cited as bookkeeping, never as proof of behavior.
- **Independent-source facts — can diverge, can fail:** the byte totals against a ceiling constant fixed before this change; the verbatim-preservation `diff` against an immutable briefing committed nine hours earlier by a different commit.
- **Behavioral-proxy evidence — readers acted, arms separated, could have gone either way:** the 6-reader blind A/B and the 2-reader entity-comprehension test.
- **Absent:** production-drive evidence. No observed live FO drive backs any claim here.

### Material findings

None.

### Deferred risks (recorded, not blocking)

- **D3 — the stop is turn-local; the Mod-Block Guard recovery path re-arms.** After a surface-and-stop there is no merge commit and no FINALIZE sentinel, so a later boot hits `mod-block=merge:*`, "Check what `«hooks.run»("merge")` left, then re-run `«merge.guard»`", and the guard returns armed again. Not an AC-1 violation — the FO has already reached the captain, which is what AC-1 requires — and the Mod-Block Guard text is untouched by this diff, so the risk is inherited, not introduced. Git itself refuses a re-merge on an unmerged tree, so it fails safe. Trigger: a boot boundary after an unresolved conflict. Promote on an observed drive where an FO re-invokes the merge hook onto a still-conflicted tree.
- **D4 — prose-only rule with no enforcing check.** The audit built an adversarial edit that defeats the guarantee while changing zero words: relocate the em-dash clause to the end of the sentence so it attaches to the `merge: pr` branch, which cannot produce a local conflict. Nothing described in either file catches it. This is real and is the honest ceiling of a contract-prose change; adding machinery was declined twice on proportionality grounds and this dispatch forbids it. Trigger: a future edit to that sentence. Promote on an observed live drive where the shipped instruction fails to produce surface-and-stop.
- **D2 (carried forward, re-corroborated) — Mod-Block Guard circularity.** 3/3 of my BEFORE readers cited the boot-scoped Mod-Block Guard clause as a mid-turn rule, independently replicating cycle 1's D2. Unchanged trigger and promote condition.

### Polish

- "The original criteria are preserved verbatim below" (line 117) is plural, but only original AC-1 is preserved; original AC-2 (`briefing-1.md:66`) appears nowhere in the entity. Derivative — the briefing scopes AC-2 as counting "only paired with AC-1" — so its non-delivery follows, but the sentence is literally false as written.
- "byte-identical to `fo-merge-core.md`" (line 98) overstates. The *prefixes* are byte-identical — verified over 288 and 165 bytes — but both quoted lines end in `…` elisions of unchanged tails. Both blind readers independently flagged this as a copy hazard.
- Cycle 1's polish note points a later copier at "the entity body (line 99)"; line 99 is now blank and the shipped text is line 100. Stale by one line, and it lands adjacent to the right block, not on the parked paragraph (line 111). The note lives inside an append-only prior stage report and should not be rewritten.
- AC-1's "Known limit" bullet understates its own evidence: the recorded 6/6 restraint null did not replicate here (2/3 BEFORE readers proposed hand-resolving, 0/3 AFTER). Erring against its own interest; recorded, no change requested.

### On the declined `done-when` finding

I did not re-open it and I produce no evidence that would. The evidence that would reopen it is an observed drive in which an FO terminates its turn at armed; my probe examined the conflict path, and no reader in either arm parked at armed — the BEFORE arm's failure mode was the opposite, over-continuing into hand-resolution. That is a further small data point that careful readers do not park, and it cuts with the park decision, not against it. I concur with declining the 80-120 byte `done-when` amendment on the FO's stated ground: re-introducing a captain-parked payload in smaller type on unchanged null evidence is decision-by-attrition.

### Recommendation

**PASSED**, with one unmet merge precondition the captain must decide on.

Both shipped ACs verified with evidence I reproduced, and AC-1 verified against the shipped text for the first time — the recorded probe had tested the parked paragraph instead. The substitution holds as a recorded retreat: the original claim survives byte-identical against an immutable independent source, and two blind readers both concluded the entity did not deliver what it set out to. The change returns 24 bytes of ratchet headroom. The blocker is not in the deliverable: the three required live lanes have not run because the branch is unpushed, which no local validation can satisfy.

### Summary

Re-ran AC-1's blind A/B from a pre-registration written before any reader and got 3/3 surface-and-stop on the after-text against 0/3 on the before-text, replicating the recorded separation — and found that the recorded run had tested the parked 843-byte paragraph rather than the shipped clause, so this is the first evidence that the text actually shipping moves the baseline. Reproduced AC-2 two independent ways that agree exactly: net −24 bytes, ratchet green, full suite green uncached. The payload substitution survives its hardest test — the original AC-1 is byte-identical to the frozen ideation briefing committed nine hours before the substitution, and two fresh readers given only the entity body both identified the shipped text correctly, flagged the parked paragraph as unshippable, and answered NO to whether the task delivered what it set out to. The detached adversarial audit held on misattribution and information loss and landed two attacks, both recorded as deferred risks: the stop is turn-local and the Mod-Block Guard recovery re-arms, and a zero-word relocation of the clause would defeat it with nothing to catch it. Four polish items recorded, none blocking. The three required live lanes remain unrun because the branch is unpushed — stated as an unmet merge precondition rather than waved through.
