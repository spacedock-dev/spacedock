---
id: bwr6j6edkmfx5sbz73cr2952
title: Feedback-cycle record convention and design-reset decision — declared estimate, per-round actuals, recorded reframe
status: done
source: "captain (2026-06-04) — forked from xa (feedback-guarantee-binary-gate) per the roadmap-the-decision + separate-build-task call. xa's ideation determined Candidate 1 (3-cycle escalation) is mechanizable via a dedicated cycle-record command (a spike disproved a --set status guard) and Candidate 2 (budget-probe) is not. This task SHIPS the Candidate-1 guard; xa closed as a roadmap decision."
score: "0.30"
started: 2026-07-20T03:29:33Z
completed: 2026-07-21T04:41:32Z
verdict: passed
worktree: .worktrees/spacedock-ensign-feedback-cycle-record-command
issue:
sprint: 0260-proportionality
group: reframe
gates:
    version: 1
    current:
        gate: gate:docs-dev:bw:ideation
        attempt: gate-attempt:bw-ideation-3
    records:
        - id: gate:docs-dev:bw:ideation
          stage: ideation
          current-attempt: gate-attempt:bw-ideation-3
          attempts:
            - id: gate-attempt:bw-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:bw-ideation-1
                digest: sha256:51b38d47a1b15d0ea7b34d4908d20d257ae3e39f5bc587839075a78466d10167
              resolution:
                type: Resolution
                id: resolution:actor-1784521963753201000
                briefing: briefing:bw-ideation-1
                by: person:reviewer
                at: 2026-07-20T04:32:43Z
                decision: revise
                reason: "Four annotations: (1) briefing packaging should use separate artifacts (FO-side, 3k experiment); (2) is the new record command needed at all — apply the cheapest-check ordering; (3) add a final landing-spot review AC (core vs dev-specific) and propose the dev README change; (4) roborev-shaped in-stage AC coverage."
              note: "Subspace advisory float; four captain annotations included by id in the resolution. Annotation 1 is FO-owned (briefing packaging), 2-4 routed to the worker; next attempt opens at re-presentation."
            - id: gate-attempt:bw-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:bw-ideation-1
              state: closed
              briefing:
                id: briefing:bw-ideation-2
                digest: sha256:520f84901aae2b1e9fc0c78eaaf974cbec3a8b8cbb82db3d8eb7fecb9a778a6d
              resolution:
                type: Resolution
                id: resolution:actor-1784523098406183000
                briefing: briefing:bw-ideation-2
                by: person:reviewer
                at: 2026-07-20T04:51:38Z
                decision: approve
              application:
                target-stage: implementation
                state: superseded
              note: "Subspace advisory float, no annotations; third presentation (prose-only convention) approved. Superseded by attempt 3 (captain-approved staff-review folds); the approval itself stands."
            - id: gate-attempt:bw-ideation-3
              sequence: 3
              previous-attempt: gate-attempt:bw-ideation-2
              state: closed
              briefing:
                id: briefing:bw-ideation-3-chat
                digest: sha256:837779a0b96ebddc7e695106109a6026034c28337686a73bba3d8d18f2ff8c6f
                note: "chat presentation; ADVISORY digest — it hashes the working file at recording time (body folds applied, this attempt's own record excluded), which no single committed tree reproduces because an entity cannot self-bind its gates record. For drift checking, diff the entity BODY against the state commit that introduced this attempt; do not re-hash the current file."
              resolution:
                type: Resolution
                id: resolution:captain-chat-bw-ideation-3
                briefing: briefing:bw-ideation-3-chat
                by: person:captain
                at: 2026-07-20T10:18:41Z
                decision: approve
                reason: "Staff-review folds, captain-approved in chat with the concrete strengthened wording shown: past-tolerance deviation requires a recorded reconfirm/re-scope/park/escalate decision before further repair dispatch (no automatic re-dispatch); tolerance default named as 2x unless the entity declares otherwise; the testlint arm struck from AC-1 (one-off script, no committed check); the skill frontmatter description line added to the declared surface; title corrected to the prose-only convention the body ships."
              application:
                target-stage: implementation
                state: consumed
              note: "FO applied the folds directly under the captain's edit-directly grant; codex staff review finding 1 and fable delta findings 8-10."
mod-block:
pr: pr-merge:541
archived: 2026-07-21T04:41:32Z
---

Turn the feedback-rejection correction loop from a prose-only cycle count into a measured, calibrated loop. At ideation the entity declares an **expected surface** as part of the captain-approved design; every correction round records its **actuals** into the durable `### Feedback Cycles` section in a documented entry format; the **deviation** of actuals from the captain-approved estimate is **narrated** at each re-dispatch decision point, so the runaway loop becomes visible against the estimate rather than legal against the prior round.

**The first cut is a documented entry-format convention — prose, not a command.** Applying this entity's own ordering (the cheapest check that can fail; build new machinery last, consent-gated), the calibration value ships with zero new code: the `### Feedback Cycles` entry format is specified in the skill and dev-README prose; the FO and the in-stage review loop hand-append a conforming entry per round (actuals from a documented one-line `git diff --numstat`); deviation is arithmetic any reader does from the estimate and an entry. The `status --record-feedback-cycle` command, its in-binary git capture, AC-digest hashing, and the escalation backstop are **deferred machinery** — shipped only when live drives show hand-authored entries drifting from the format (which the format itself makes checkable). Three independent sources converged on this: the codex cross-review, this entity's own ordering, and the captain's gate annotation; and no falsifiable capability the convention lacks was found (see Problem). Hard dispatch-refusal on raw diff growth stays demoted for good — it false-refuses a legitimately growing fix, and a non-zero `dispatch build` exit is read by the FO as an infra failure that triggers Break-Glass **manual** dispatch (`fo-dispatch-core.md:145`), so a build-side refusal is routed around, not honored.

## First cut vs deferred

**Expected surface (first cut, prose-only — this entity practicing the convention it ships):** ~3 prose files — `feedback-rejection-flow` (frontmatter description + steps 2-3 + the entry-format spec, ~10 lines), `first-officer-shared-core` (`«feedback.route»` effect/done-when, ~2 lines), `docs/dev/README.md` (ideation expected-surface line + the implementation in-stage-round convention, ~5 lines); plus 1 fixture entity and 1 offline deviation/AC-drift check (the spike, formalized) for AC-1/AC-5/AC-6. **0 Go source files, 0 product LOC.** Declared tolerance: 2× — and a hard self-check: any Go/product code appearing in the first cut, or a fourth doc file, trips a reconfirm, because deferring the command IS the point.

**First cut (this entity), all prose + a fixture:**
- The documented `### Feedback Cycles` **entry-format convention**: each round appends `- {timestamp} — {reviewer/loop} {verdict}; surface {files}/{LOC} vs estimate {declared} ({P}%); AC {unchanged | narrowed: <note>}`. Specified in the skill and dev-README prose.
- The two skill-prose swaps: `feedback-rejection-flow` steps (hand-append a conforming entry; read the deviation + AC-drift; record a design-reset decision past tolerance) and the `first-officer-shared-core` `«feedback.route»` lines.
- The `docs/dev/README.md` changes: the ideation expected-surface line, and the implementation in-stage-round convention (in-stage review rounds append a conforming entry; the documented `git diff --numstat` one-liner for actuals).
- A fixture entity + an offline check proving deviation and AC-drift are computable from two conforming entries (AC-1 / AC-5 / AC-6).

**Deferred machinery (ships only when live drives show hand-authored entries drifting from the format):**
- The `status --record-feedback-cycle` command, its in-binary `git diff --numstat` capture, and section-scoped auto-count.
- AC-drift as a computed hash/digest (the convention ships the human-authored note; automation is deferred).
- The escalation marker + refuse-further-auto-bounce + `--force` + its `feedback-escalate` schema field.
- The `present-gate` surface-deviation evidence line; the template `Outputs` propagation (the `template` member's); a CLI docs page (no CLI surface until the command ships).
- The enforcement half of AC-drift: an AC-weakening edit after a rejection treated as a design-reset event requiring captain sign-off.

## Problem

- **The runaway loops were contract-legal round by round.** All four HIGH incidents in the 0260 forensics (`_evidence/0260-agent-derail-forensics/synthesis.md`) — e6j (2-defect fix → 10 roborev cycles, 26 files / +3,373, PR closed), dp (one-paragraph fix → 4-cycle ladder, discarded, ~38.5h), task-91 (16 roborev panels, own round-limit bypassed), 7h (harness repaired twice before park) — passed every round against the only baseline available: the prior round's accident. No single round stood out. The baseline has to be the entity's own captain-approved intent, not the last cycle's overrun.
- **AC-narrowing under validation pressure** (synthesis addendum, 2026-07-20, the 0.25.1 release; the addendum calls it repair-forward's paperwork twin). When validation correctly found a value claim unproven, the task narrowed its AC until a weaker claim passed — a real rejection converted into a paperwork pass; the failure then reproduced live. The gate cross-check compares against the CURRENT AC text, so silent narrowing defeats it by construction. A calibrated loop must make AC-drift across rounds visible to a reader, the same way it makes surface-growth visible.
- **In-stage rounds have no durable record.** `### Feedback Cycles` tracks only cross-stage gate bounces (dp, 7h). e6j's 10 roborev cycles and task-91's 16 panels were in-stage (roborev-at-end-of-implementation, never crossing a gate), so cycle tracking never saw them. Both loop shapes must land in one section, under one convention.
- **Do we need a command? No — the convention is the cheapest check that can fail.** Applying this entity's own ordering (and the falsifiability ladder it serves): the value is a declared baseline + per-round actuals + reader-computable deviation + narration; all of that is a documented format plus a one-line `git diff`. A Go command adds tamper-evidence and auto-capture, but the observed failures were nobody-looking and wrong-baseline, not forged counts — the convention + narration kills them. The falsifiable rebuttal test ("what can the command do that the convention cannot, that the first cut needs?") returns nothing: tamper-evidence guards an adversary that did not appear, and auto-capture replaces a one-liner whose mistyping is exactly the drift the format makes checkable. So the command is deferred until that drift is observed live. (xa's determination that the cycle-record WRITE is the right *hook* still stands; it just isn't first-cut.)
- The prose-only guarantee's old ceiling was "the wording is present" — that failure mode was a bare count with no actuals and no baseline. A *formatted* entry carrying actuals + estimate + deviation is falsifiable at the reader level (a missing or malformed field is visible), which is what the old prose lacked.

## Proposed approach

### Generic principle (workflow-agnostic — lands in core contract prose: the skills)

1. **Declare expected surface at ideation.** The gated design states the surface it expects to touch, in the workflow's own unit — dev: "~x files, ~y LOC"; a research workflow: "~10 external docs". Captain-approved at the ideation gate → this is the baseline.
2. **Record per-round actuals in a conforming entry.** Every correction round (whatever the loop shape) appends one `### Feedback Cycles` entry carrying its actuals, the declared estimate, the deviation%, and an AC-drift note.
3. **Deviation vs the declared baseline** (not the prior cycle), beyond a **declared tolerance** (default 2× unless the entity declares otherwise), is **narrated** and requires a **recorded decision** — reconfirm / re-scope / park / escalate — **before any further repair dispatch**. No automatic re-dispatch while that decision is absent. (Captain-approved strengthen, staff-review fold 2026-07-20.)
4. **Narration is the guard.** The deviation% and AC-drift, read from the entries, are what the FO looks at where it decides whether to re-dispatch; past tolerance, the recorded decision above is required before the next dispatch. The enforcement machinery stays deferred.
5. **Escalation is the backstop — DEFERRED.** A hard stop ships only if live drives show narration being ignored.

### Dev-specific realization (lands in dev prose: `docs/dev/README.md` + the template)

- **Surface unit:** files touched / LOC.
- **Actuals capture (documented one-liner, no binary):** `git diff --numstat "$(git merge-base main HEAD)"..HEAD` — cumulative surface of the entity's work vs its branch point; sum the file and line columns. Deterministic, no round-store (roborev has none; rounds are Git-versioned — spike-proven). The FO/loop reads the two numbers and writes them into the entry.
- **In-stage review-round instance:** the roborev re-panel (or the detached-audit pass) is the dev instance of the generic in-stage-round rule; it appends a conforming entry per re-panel.

### Both loop shapes, one convention (trigger paths named)

- **Cross-stage feedback bounce** (dp, 7h): the FO's `feedback-rejection-flow`, when a gate recommends REJECTED and routes to the `feedback-to` target, hand-appends a conforming entry. Generic; the skill carries it.
- **In-stage review round** (e6j's 10 roborev cycles, task-91's 16 panels): an in-stage reviewer loop that returns material findings and triggers another pass appends a conforming entry before the next pass. Generic shape (any in-stage review loop, any workflow); the dev instance is the roborev re-panel / detached audit, carried in `docs/dev/README.md`.

Both write the same section under the same format: `### Feedback Cycles` becomes the single durable record of correction rounds. Making in-stage rounds record here is precisely what would have surfaced e6j at all — and the dev README already writes `### Feedback Cycles` entries for detached-audit findings (README validation stage-def), so this extends an existing habit, not a new mechanism.

### The record command (DEFERRED)

`spacedock status --record-feedback-cycle {slug}` would own the append + section-scoped count + in-binary `git diff` capture + AC-digest + (later) the escalation marker. It composes proven machinery (`section_read.go` body-parse, `mutate.go` `atomicWrite`, `merge.go arm`), so building it is cheap **when justified** — i.e., when live drives show hand-authored entries drifting from the format. It is not first-cut: the convention delivers the same reader-computable calibration with no code, and the format makes hand-authoring self-checkable.

### Narration read (first cut)

- The FO reads the deviation% and AC-drift directly from the two most recent entries per the skill prose, at the re-dispatch decision point (its own awareness before routing another repair).
- The `present-gate` render is deferred; a rendered round-state string is a convenience the command would later emit.

### Design-reset routing (narration-driven)

When the deviation read from the entries crosses tolerance — or the AC narrowed — the FO weighs a **design-reset decision** (reconfirm the estimate / re-scope / park / escalate) instead of automatically routing the next repair; a recorded reframe re-baselines. This aligns with prose already in `docs/dev/README.md` (implementation stage-def: "Stop and request a design reset when the simpler … path still reaches the value"; validation stage-def: classify a rejection as mechanism-failure and "recommend a scope/design reset … do not send it through another automatic implementation feedback cycle"). The enforced backstop is deferred.

### AC-drift note (first cut, convention form)

Each conforming entry records whether the `## Acceptance criteria` text is unchanged or narrowed since the prior round, with a one-line note. This makes AC-narrowing under validation pressure (synthesis addendum) visible to a reader from two entries — the same reader-level check as surface deviation. The computed **hash/digest** form (automation) is deferred with the command; the first cut ships the human-authored note. Honest limit: the note's baseline is the prior recorded round, so it catches narrowing DURING the loop (the 0.25.1 shape); narrowing during initial implementation, before any entry, needs an ideation-gate AC snapshot — deferred. The ENFORCEMENT (an AC-weakening edit is a design-reset event requiring captain sign-off) is prose the backstop carries, deferred.

### Per-mechanism justification (value AC served / simplest alternative / why insufficient)

- **The entry-format convention (first cut):** serves AC-1 / AC-2. Alt: a bare prose cycle count (status quo). Insufficient — no actuals, no baseline, uncomputable deviation; the formatted entry is the minimum that makes deviation reader-computable and self-checkable.
- **Deviation vs the declared estimate (not the prior cycle):** serves AC-1. Alt: diff-growth vs the prior cycle (a superseded design). Insufficient — every e6j round passed round-by-round; the prior-cycle baseline never stands out (spike counterfactual: +5 files reads as ordinary growth).
- **AC-drift note (AC-5):** serves AC-1's calibration by covering the narrow-AC direction. Alt: rely on the gate's live cross-check of current AC text. Insufficient — the addendum proves that cross-check is defeated by construction when the AC is silently narrowed; a durable per-round note makes the edit comparable across rounds.
- **The command / auto-capture / hash / escalation marker (DEFERRED):** would add tamper-evidence + automation. Alt: the convention + a one-liner + a reader. Sufficient for the first cut — the deferred machinery is justified only by observed drift, per this entity's ordering; building it now is the over-engineering the sprint exists to stop.

## Landing-spot map (generic principle in core prose vs dev realization in dev prose; first-cut vs deferred)

| Piece | Kind | Cut | Lands in |
|---|---|---|---|
| Declare expected surface at ideation | generic | first | `docs/dev/README.md` ideation stage-def (the dev instance of a generic ideation output) |
| Entry-format shape; per-round actuals; deviation vs baseline; reconfirm-or-reframe; narration read; in-stage-round rule | generic | first | `feedback-rejection-flow/SKILL.md`; `first-officer-shared-core.md` effect line |
| files/LOC unit; the `git diff --numstat` one-liner; the roborev/detached-audit in-stage instance | dev | first | `docs/dev/README.md` (ideation + implementation stage-defs) |
| Fixture entity + offline deviation/AC-drift check | dev | first | state checkout / `internal/status` testdata |
| `status --record-feedback-cycle` command + in-binary git capture + section-scoped auto-count | dev/binary | deferred | `internal/status` (`section_read.go`, `mutate.go` `atomicWrite`, `runGitCmd`) |
| AC-drift as a computed hash | dev/binary | deferred | the command |
| escalation marker + refuse + `--force` + `feedback-escalate` schema field | dev/binary | deferred | `internal/status` (`merge.go arm`, `handlers.go`, embedded schema) |
| present-gate surface-deviation line; template `Outputs` propagation; CLI docs page | mixed | deferred | `present-gate/SKILL.md`; `template` member; docs site |
| AC-weakening = design-reset requiring captain sign-off | generic | deferred | `feedback-rejection-flow/SKILL.md` + backstop |

### Landing-placement audit (annotation 2 — verified at validation, AC-7)

The generic clauses must land in the fleet-loaded skill prose and the dev specifics in `docs/dev/README.md` / template — neither leaking into the other layer. AC-7 makes that a checked property after implementation.

## Riskiest-mechanism spike (done first, per ideation policy)

**Claim under test:** can per-round actuals be captured from durable evidence, and deviation computed vs the declared estimate (not the prior cycle), with NO command — from a documented one-liner and two text entries? Exercised end-to-end against real git state (throwaway repo, results recorded — not asserted):

- **Cumulative actuals from durable git** via `git diff --numstat {BASE}..HEAD`: round 1 = **2 files**, round 2 = **7 files** (base = the pre-fix commit; entity declares estimate 2 files / 40 LOC).
- **Deviation vs the declared estimate** (2 files): round 1 = **100%**, round 2 = **350%** → crosses a 200% tolerance **at round 2** (the e6j "surface visible by round 2" target from the DoD).
- **Counterfactual vs the prior cycle** (the accident baseline e6j effectively used): round 2 − round 1 = **+5 files**, reads as ordinary growth — proving the baseline must be the estimate, not the accident.
- **Deviation recomputes offline from two durable `### Feedback Cycles` text entries**, no live process, no command — confirming the first-cut convention is a reader/measurement, and that the command is not needed to realize the value.

**Determination:** the estimate/actual semantics are fully realized by the documented one-liner + the entry format + reader arithmetic. No command needed for the first cut; no further spike needed. The throwaway seeds the AC-1 fixture directly.

## Acceptance criteria

**AC-1 (VALUE) — Fed the archived e6j per-round surface shape as conforming `### Feedback Cycles` entries (2-file / 40-LOC estimate; round-2 actuals from e6j's history), the deviation prescribed by the documented format and arithmetic reads ≥ 200% of estimate at round 2 — surfacing the runaway two rounds before e6j's real history did.**
Verified by: a checked-in fixture entity (estimate + round-1 / round-2 entries) and an offline check (the spike, formalized as a one-off script whose output lands in the validation report — no committed check) that applies the documented deviation arithmetic and asserts ≥ 200% at round 2; the check fails if an entry omits actuals/estimate or the arithmetic is mis-specified. Independent baseline that moved the wrong way: e6j's real 10-round / 26-file history. This is the outcome the entity exists for — the runaway made visible at round 2, by prose + a reader, no code.

**AC-2 — The `### Feedback Cycles` entry-format convention is specified in the skill and dev-README prose and is sufficient for a reader to compute deviation vs the ideation baseline from the estimate and one entry.**
Verified by: a one-off prose grep run DURING validation (output pasted into the validation report, not committed as a test) confirming `feedback-rejection-flow` and `docs/dev/README.md` specify the format fields (actuals, estimate, deviation%, AC-drift note) and the `git diff --numstat` one-liner; plus the AC-1 fixture demonstrating a reader computing deviation from a conforming entry. Per the captain's prose-grep ruling, a grep over prose the implementer wrote cannot fail independently, so it is validation-time evidence, not a checked-in test. Serves AC-1.

**AC-3 (deferred-scope note) — The `status --record-feedback-cycle` command, its in-binary git capture and auto-count, AC-drift-as-hash, the escalation marker + refuse + `--force` + `feedback-escalate` schema field, the `present-gate` line, the template propagation, and a CLI docs page are NOT in this cut.** They ship only when live drives show hand-authored entries drifting from the format (which the format makes checkable). Recorded so the boundary is explicit; xa's determination that the cycle-record WRITE is the correct *hook* holds for when the command does ship.

**AC-4 — The skill prose directs the FO to read the round state (deviation%, AC-drift) from the entries at the re-dispatch decision point and weigh a design-reset decision when beyond tolerance.**
Verified by: a one-off prose grep during validation (output in the report, not a committed test) confirming `feedback-rejection-flow` carries the read-and-weigh step. `present-gate` render deferred. The FO ACTING on it is gq's live scenario, not this AC.

**AC-5 — Each conforming entry carries an AC-drift note (unchanged / narrowed + one-line reason), making AC-narrowing across two entries visible to a reader.**
Verified by: the AC-1 fixture includes a round whose AC narrowed and whose entry flags it, plus a control round with unchanged AC that does not; the offline check independently compares the AC text across the two rounds and asserts the entry's flag matches that comparison (so a wrong flag fails it — a real value against an independent source, not a prose-grep). The computed-hash form is deferred (AC-3).

**AC-6 — An in-stage review round produces a conforming entry in the same `### Feedback Cycles` section by the same convention — phrased so any similar-shape in-stage reviewer loop, in any workflow, is covered, not a roborev-only clause.**
Verified by: a fixture (or prose audit) showing an in-stage-round entry conforming to the format in the same section as a cross-stage entry; `docs/dev/README.md`'s implementation stage-def carries the generic in-stage-round trigger with the roborev re-panel / detached audit as the dev instances. (Annotation 3.)

**AC-7 — After implementation, the generic principle lives in core contract prose (the skills) and the dev realization in `docs/dev/README.md` / template; neither leaks into the other layer.**
Verified by: a validation-stage landing-placement audit — read/grep the two surfaces and assert the generic clauses (declare-surface, entry-format shape, deviation-vs-baseline, narration read, in-stage-round rule) are present in the skill and carry no dev-specific unit, and the dev specifics (files/LOC unit, the git one-liner, the roborev/detached-audit instance) are present in `docs/dev/README.md` and absent from the skill. (Annotation 2's "final landing-spot review".)

## Test plan

- **Value replay (fixture + offline check) → AC-1 / AC-5:** a checked-in fixture entity with conforming entries (e6j surface shape; a narrowed-AC round; controls) and a small offline check applying the documented arithmetic; it fails if a field is missing or the arithmetic is wrong. This is the spike formalized — cheap, falsifiable, no product code. Cost: low.
- **Validation-time prose grep (NOT committed) → AC-2 / AC-4:** a one-off grep during validation confirms the skill / dev-README prose specifies the entry format, the one-liner, and the read-and-weigh step; the output is pasted into the validation report. Per the captain's prose-grep ruling this is external evidence, not a checked-in test — a grep over implementer-written prose cannot fail independently (the tautological-test anti-pattern the sprint retires). No `sectionAfter` port is needed.
- **In-stage-round conformance (fixture / prose audit) → AC-6:** an in-stage entry conforms in the same section; the dev README carries the generic trigger.
- **Landing-placement audit (prose audit at validation) → AC-7:** the core-vs-dev grep/read check.
- **No Go unit tests, no product code in this cut.** The command's unit tests (append / count / git-capture / escalation / `--force`) ship with the deferred command.
- **High-stakes note:** deferred — the `status` mutation/guard paths are not touched until the command ships; the detached adversarial audit applies then.

## Documentation changes (first cut; concrete before/after — ideation proposes, implementation applies)

**`skills/feedback-rejection-flow/SKILL.md`** — steps 2-3:

> - Before (step 2): `2. Track cycles in `### Feedback Cycles` in the entity body.`
> - After: `2. Append one `### Feedback Cycles` entry for this round in the format: `- {timestamp} — {reviewer/loop} {verdict}; surface {actuals} vs estimate {declared} ({P}%); AC {unchanged | narrowed: <note>}`. Capture the surface actuals with the workflow's documented one-liner (dev: `git diff --numstat "$(git merge-base main HEAD)"..HEAD`).`
> - Before (step 3): `3. On cycle 3, escalate to the human instead of another round.`
> - After: `3. Read the deviation% and AC-drift from the two most recent entries. When deviation is beyond the declared tolerance, or the AC narrowed, record a design-reset decision — reconfirm / re-scope / park / escalate — before any further round; no automatic re-dispatch while that decision is absent.`
> - Before (frontmatter description, line 3): `… track `### Feedback Cycles`, escalate on cycle 3 …`
> - After: `… append conforming `### Feedback Cycles` entries; past tolerance, record a design-reset decision before another round …`

**`skills/first-officer/references/first-officer-shared-core.md:102-103`** — the `«feedback.route»` effect/done-when:

> - Before: `… track `### Feedback Cycles`, escalate on cycle 3, …` / `… (or escalated at cycle 3).`
> - After: `… append a conforming `### Feedback Cycles` entry (actuals + estimate + deviation + AC-drift note) and read the round state, …` / `… (or a recorded design-reset decision when deviation is beyond tolerance).`

**`docs/dev/README.md`** — ideation stage-def `Outputs` (add one sub-bullet under the Outputs list):

> `- The task body declares an **expected surface** — the files/LOC it expects to touch — as part of the gated design; the ideation gate approves it as the baseline the correction loop calibrates against.`

**`docs/dev/README.md`** — implementation stage-def (add one bullet after the design-reset bullet, line 123):

> `- When an in-stage review round (a roborev re-panel, a detached-audit pass, or any similar in-stage reviewer loop) returns material findings that trigger another pass, first append a conforming `### Feedback Cycles` entry: `- {timestamp} — {reviewer} {verdict}; surface {files}/{LOC} vs estimate {declared} ({P}%); AC {unchanged | narrowed: <note>}` — actuals from `git diff --numstat "$(git merge-base main HEAD)"..HEAD`. Read the deviation vs the ideation estimate from the entry; beyond the declared tolerance, record a design-reset decision (see the bullet above) before another pass. This is the same convention the cross-stage feedback flow uses; in-stage and cross-stage rounds share one section.`

Deferred doc changes (not proposed here): the `present-gate` evidence line, the template `Outputs` propagation, and a CLI docs page (no CLI surface until the command ships).

## Boundary and notes

- **gq (`feedback-nonhappy-live-coverage`) owns the live half:** that the FO acts on the narration is FO-LLM behavior with no in-process place to observe it — its scenario, not an AC here.
- **Out of scope:** the budget-probe fail-safe (xa Candidate 2); any `dispatch build` refusal (demoted). The command itself is deferred, not out of scope — its hook and machinery are recorded for when drift justifies it.
- **Leanness (0250/0260):** the first cut is prose + a fixture, no boot-resident additions beyond ~2 lines in the always-on shared-core effect line; the entry-format spec is lazy-loaded (`feedback-rejection-flow` loads at the rejection point) and the dev specifics ride `docs/dev/README.md` (not boot). Net contract-byte delta is small (the old cross-stage tracking prose is replaced, not added to). Deferring the command removes all Go / schema / CLI byte cost from this cut.
- **Lane calibration:** the convention settled here (the entry-format fields; the one-liner; the in-stage-round rule) is what ek and w0 build against after this gate; ht runs parallel.
- No minted identifier schemes or coined abstractions; bare ordinals throughout.
- Siblings: xa (`feedback-guarantee-binary-gate`, archived — determination + grounding), gq (`feedback-nonhappy-live-coverage` — live half). Provenance: a9 (`feedback-rejection-flow-skill-extraction`) detached audit, 2026-06-04.

## Stage Report: ideation

- DONE: Riskiest mechanism spiked first: estimate/actual surface semantics — what counts as declared expected surface, how per-round actuals are captured from durable evidence, what tolerance triggers reconfirmation — smallest end-to-end exercise recorded in the body, not asserted.
  Ran a throwaway git repo end-to-end (see `## Riskiest-mechanism spike`): `git diff --numstat BASE..HEAD` captured cumulative actuals (round 1 = 2 files, round 2 = 7 files); deviation vs the declared 2-file estimate hit 350% at round 2 (crosses 200% tolerance); the prior-cycle counterfactual (+5 files) never screams; deviation recomputed offline from two durable `### Feedback Cycles` entries.
- DONE: The design separates GENERIC principle from DEV-SPECIFIC realization and names where each lands.
  `## Proposed approach` splits the two; `## Landing-spot map` names a landing spot per piece (skill prose / shared-core / template ideation scaffold / `internal/status` binary / roborev-loop prose).
- DONE: One correction-round record covers BOTH loop shapes — cross-stage feedback bounce (dp, 7h) and in-stage review rounds like roborev-at-end-of-implementation (e6j's 10 rounds, task-91's 16 panels) — with the trigger paths named.
  `### Both loop shapes, one record` names the two trigger paths (FO `feedback-rejection-flow` for cross-stage; the implementation-stage roborev loop for in-stage) both writing the single `### Feedback Cycles` ledger.

### Summary

Reframed the entity to the captain's 2026-07-20 estimate-then-calibrate direction: ideation declares an expected surface (captain-approved baseline), each correction round records actuals into `### Feedback Cycles`, and deviation-vs-estimate is narrated at re-dispatch and gate decision points; the banked 3-cycle escalation survives as the threshold-driven backstop. Hard dispatch-refusal on raw diff growth is demoted (false-refuses growing fixes; a non-zero `dispatch build` exit is routed around as a Break-Glass infra failure, not honored) — the record-command write is the correct hook. Added a value AC (e6j replay: deviation crosses tolerance by round 2 against the archived 10-round/26-file baseline), the generic/dev landing-spot map, and concrete doc diffs for `feedback-rejection-flow`, `first-officer-shared-core`, `present-gate`, and the ideation stage-def. The riskiest mechanism was exercised against real git state, not asserted.

## Stage Report: ideation (cycle 2)

- DONE: Reduce to the narration-first first cut (gate ask 1-2).
  First cut = the record command core (append + section-scoped count + actuals + deviation + stdout round-state line) + the two skill-prose swaps; the escalation marker/refuse/`--force`/schema field, the `present-gate` line, the template `Outputs` propagation, and the CLI docs page are moved to `## First cut vs deferred` and marked deferred. AC-3 reduced to a deferred-scope note; AC-4 kept the stdout half and dropped `present-gate`; AC-5 kept the two swaps + the one dev-README line.
- DONE: Keep the declared-estimate baseline home (gate ask 3).
  The one-line expected-surface declaration stays in the dev README ideation doc-diff; template propagation is annotated as the `template` member's, not this entity's.
- DONE: Answer the AC-narrowing design question (gate ask 4), reasoning recorded.
  Read the 0.25.1 synthesis addendum. Recommend INCLUDING the AC-drift detection+narration in the first cut (AC-6, severable), DEFERRING the enforcement (captain sign-off on AC-weakening). Rationale in `### AC-drift detection`: it is the narration half of the paperwork-twin fix, rides the existing `readSections` + entry-write (one hash, one entry field, one narration clause — no new command/gate/frontmatter/schema), and covers a live-reproduced failure; honest limit recorded (baseline is round 1, so pre-loop narrowing needs a deferred ideation-gate snapshot).

### Summary

Applied the gate bounce: shrank to the narration-first core (record command + two skill swaps + one ideation baseline line), moved the hard backstop / present-gate line / template propagation / CLI page to a deferred set with an explicit `## First cut vs deferred` split, and restructured AC-1..AC-5 as directed. On the AC-narrowing question I recommend including the AC-drift digest as a severable AC-6 in the first cut (cheap detection, same narration class, covers the fresh 0.25.1 paperwork-twin) while deferring its enforcement half; reasoning recorded in the body and above. The e6j value AC and the spike are unchanged.

## Stage Report: ideation (cycle 3)

- DONE: Append the explicit written surface estimate (pre-gate captain directive), practicing the convention this entity ships.
  Added the **Expected surface** line under `## First cut vs deferred` with this entity's own first-cut numbers (~2 Go files + 1 test file, ~150 LOC source / ~200 LOC test; 2 skill-prose files, ~6 lines; 1 dev-README line) and a declared 2× tolerance.
- DONE: Terminology sweep for coined vocabulary.
  Replaced my coined "ledger" with plain "record"/"section" (3 design-body sites); demoted "paperwork twin" from a load-bearing label to a single attributed quote of the synthesis addendum, using plain "AC-narrowing under validation pressure" everywhere it does work. Kept existing system vocabulary as-is ("round-state line", "### Feedback Cycles"). Historical cycle-1/cycle-2 stage reports left verbatim.

### Summary

Pre-gate polish only: appended the written surface-estimate line in the standard form (the entity now declares its own baseline, eating its own dog food) and swept coined vocabulary out of the design body. No design change.

## Stage Report: ideation (cycle 4)

- DONE: Apply the entity's own ordering to the command (annotation 2: "do we need this?"). Answer: no — deferred.
  Reworked the first cut to a documented `### Feedback Cycles` entry-format CONVENTION (prose): the FO and the in-stage loop hand-append conforming entries, actuals from a documented `git diff --numstat` one-liner, deviation computed by any reader. The `status --record-feedback-cycle` command, in-binary capture, AC-hash, and escalation machinery are deferred, shipped only when live drives show hand-authored entries drifting from the format. Recorded the falsifiable rebuttal (nothing the convention cannot do that the first cut needs; tamper-evidence guards an absent adversary). Surface shrank to prose-only: 0 Go files.
- DONE: Add the core-vs-dev landing audit AC and the concrete dev-README diff (annotation 2's second half).
  AC-7 audits landing placement (generic in the skills, dev specifics in `docs/dev/README.md`) at validation; `## Documentation changes` now carries concrete before/after diffs for `docs/dev/README.md` — the ideation expected-surface line AND the implementation in-stage-round convention.
- DONE: Restore dedicated in-stage-loop AC coverage, generically (annotation 3).
  AC-6 covers an in-stage review round producing a conforming entry in the same section, phrased for any similar-shape in-stage reviewer loop in any workflow (roborev re-panel / detached audit are the dev instances), not a roborev-only clause.
- DONE: Update the written estimate to the reworked scope.
  Expected surface is now prose-only: ~3 doc files (~16 lines) + 1 fixture + 1 offline check; 0 Go source / 0 product LOC; tolerance 2× with a hard self-check that any product code in the first cut trips a reconfirm.

### Summary

Applied the REVISE. The design now ships a documented entry-format convention (prose + a fixture), not a command — the honest application of this entity's own cheapest-check-first ordering, converging with the codex cross-review and the captain. The command and all binary / schema / CLI machinery are deferred behind observed format-drift. Added AC-7 (core-vs-dev landing audit) with concrete `docs/dev/README.md` diffs, and AC-6 (generic in-stage-round coverage). ACs renumbered to the prose-only scope; the e6j value proof is now a fixture + offline check (no product code), directly validated by the spike. Estimate updated to 0 Go LOC.

## Stage Report: implementation

- DONE: ONE composed `### Feedback Cycles` entry format lands in skills/feedback-rejection-flow/SKILL.md carrying BOTH members' fields: bw's declared surface / estimate / per-round actuals / AC-drift AND 02av's findings disposition — including the all-declines case, which must render as a real recorded state and not as an empty field.
  `## Feedback Cycles entry` (commit ee4a23cf) specifies one line per round: `- Cycle {N}: {verdict} — {reviewer/loop}; surface {actuals} vs estimate {declared} ({P}%); findings {none | {F} fixed, {D} declined: <ref · class · why not material · promotes when>}; AC {unchanged | narrowed: <note>}`, with `findings none` (nothing arrived) held distinct from an all-declines round's `0 fixed` plus every decline named.
- DONE: Both members' contract prose lands in the files each declared and NOWHERE else.
  `git diff --numstat HEAD~1` = 3 files, +33/-8: `docs/dev/README.md` (+3), `feedback-rejection-flow/SKILL.md` (+28/-6), `first-officer-shared-core.md` (+2/-2). `skills/ensign/references/ensign-shared-core.md` untouched — 02av's always-loaded delta is zero.
- FAILED: BYTE-FUNDED and suite-green: `go test ./...` passes INCLUDING TestFOFunctionPromptSurfaceShrinks, with per-file byte accounting recorded.
  `go test ./...` is green on every suite EXCEPT the ratchet: surface = 125,564 bytes against baseline 122,634. Accounting below; the residual 2,903 bytes are not fundable from these files without deleting meaning, so this is reported rather than absorbed. Baseline constant untouched; no prose relocated to an unmeasured file.

### Byte accounting (per file, per member)

| File | Before | After | Δ |
|---|---|---|---|
| `skills/feedback-rejection-flow/SKILL.md` | 2,741 | 5,970 | +3,229 |
| `skills/first-officer/references/first-officer-shared-core.md` | 26,298 | 26,483 | +185 |
| `docs/dev/README.md` (unmeasured) | 27,533 | 29,551 | +2,018 |
| **Measured total** | **122,150** | **125,564** | **+3,414** |

Within the skill: frontmatter description + intro **trimmed −248** (self-funding: the description no longer re-enumerates the procedure the numbered steps carry, which also retires the fable staff-review finding that it prescribed the superseded rule); steps 2-3 +358; `## Feedback Cycles entry` +899; `## Finding-triage block` +2,032; the findings-routing amendment +160. **bw's measured net +1,194; 02av's +2,192.** Headroom was 483, so bw alone overruns by 711 and the composed change by 2,903.

Trim search, run before cutting anything: a sentence-level duplication scan across all 13 ratcheted files recovers **110 bytes** in total (one shutdown-sweep sentence shared by `using-legacy-claude-team` and `claude-fo-dispatch`). The measured set is already deduplicated by prior leanness sprints; the remaining candidates in `first-officer-shared-core.md` are all either bound by contractlint literals or landed by live sibling members this sprint. Two honest reductions were taken instead: the −248 above, and removing the dev `git diff --numstat` one-liner from the generic skill (an AC-7 layering leak the ideation doc-diff had proposed — the one-liner and the files/LOC unit now live only in `docs/dev/README.md`).

### Live-lane reconciliation (no fixture or assertion edited)

The shipped `feedback-3-cycle-escalation` scenario asserts `^- Cycle \d+:` entries plus an escalation-to-human handoff, both section-scoped. Leading the entry with `- Cycle {N}:` rather than a timestamp satisfies the new convention and that assertion at once, and it gives the reader the round ordinal the deviation arithmetic and the cycle-3 backstop both need. Exercised, not asserted: a one-off test fed a three-entry conforming body to the real `assertThirdCycleEscalation` — PASS — and the same body with the handoff clause removed — RED, so the pass is not vacuous. The one-off was deleted, not committed.

Two reconciliations against the ideation doc-diff, both recorded rather than silent: the cycle-3 escalation clause is KEPT in the flow (the ideation "after" text dropped it; AC-3 defers the escalation *machinery*, not the shipped prose rule, and both the live lane and the contractlint anchor read it), and the dev one-liner is out of the skill per AC-7.

### Summary

Landed the composed convention as prose in the three declared files — 0 Go, 0 product LOC, surface 33 added lines against a ~35-line combined estimate. One entry shape carries surface deviation and findings disposition as adjacent fields, so the healthy response to review pressure (a recorded decline) and the pathological one (a narrowed AC) sit side by side in one record. The change does NOT fund itself: it needs 2,903 bytes the ratcheted files cannot give up without losing meaning, which is a captain call on re-baselining or re-scoping — I did not touch the baseline, relocate prose to dodge the check, or force a trim.

## Stage Report: implementation (cycle 2)

- DONE: ONE composed `### Feedback Cycles` entry format lands in skills/feedback-rejection-flow/SKILL.md carrying bw's declared surface / estimate / per-round actuals / AC-drift.
  Re-cut to this entity's own fields after the captain parked 02av (commit 457b910d): `- Cycle {N}: {verdict} — {reviewer/loop}; surface {actuals} vs estimate {declared} ({P}%); AC {unchanged | narrowed: <note>}`. The `findings` disposition clause and the all-declines prose were withdrawn with the member that owned them; the design for both is preserved in 02av's body for a 3k-based redesign.
- DONE: bw's contract prose lands in the files it declared and NOWHERE else.
  `git diff --numstat HEAD~2` = 3 files, +15/-7: `docs/dev/README.md` (+2), `feedback-rejection-flow/SKILL.md` (+13/-5), `first-officer-shared-core.md` (+2/-2). 0 Go, 0 product LOC. Against the declared ~17 lines across 3 prose files, inside the 2× tolerance.
- FAILED: BYTE-FUNDED and suite-green: `go test ./...` passes INCLUDING TestFOFunctionPromptSurfaceShrinks, with per-file byte accounting recorded.
  Every suite green except the ratchet: surface **122,839** vs baseline 122,634 — **206 bytes over**. bw's measured net is **+689** (skill +547, shared-core +142) against 483 of headroom.

### Byte accounting after the re-scope

Removing 02av's contributions recovered 2,725 measured bytes (the standing block, the routing amendment, the `findings` clause, and the all-declines prose), taking the composed +3,414 down to +689. One further reduction was taken on layering grounds, not to buy bytes: the boot-resident `«feedback.route»` effect line no longer re-lists the entry's fields, since the skill owns the format and that file states it references this procedure by name (−38).

What remains is rule text with no padding left in it: the entry format itself, the deviation-vs-approved-estimate rule and the one sentence saying why the baseline is the estimate and not the prior round (the entity's central finding — without it a reader re-baselines to the prior round, which is the exact failure e6j had), the tolerance default, the design-reset decision and its no-automatic-re-dispatch clause, and the ideation expected-surface line. Per the captain's instruction I did not shave into that to close 206 bytes, and did not touch the baseline constant. The residual is a call about the estate's budget.

### Feedback Cycles

- Cycle 1: RE-SCOPE — captain, on the implementation stop; surface 3 files/15 lines vs estimate 3 files/~17 lines (88%); AC unchanged — the composed-entry scope was re-cut by captain decision (02av parked), not narrowed under review pressure; every AC of this entity stands as gated.
- Cycle 2: RECONFIRM — captain, on the escalated byte-ratchet stop; surface 4 files/18 lines vs estimate 3 files/~17 lines (106%); AC unchanged — the estimate holds; the single Go-file line is the ratchet baseline the captain approved re-setting, not scope growth.

### Summary

Re-scoped to bw alone and re-measured, as directed. Withdrawing 02av's fields took the overrun from 2,903 bytes to 206 — the re-scope recovered 94% of it, which is what the team lead predicted. The convention now ships as the correction-round record it was gated as: declared estimate, per-round actuals, deviation against that estimate, AC-drift, and a recorded design-reset decision past tolerance. The `- Cycle {N}:` leading form is unchanged, so the shipped 3-cycle escalation lane still reads a conforming entry with no fixture or assertion edited. The entry above is this entity practicing its own convention on the round that produced it.

## Stage Report: implementation (cycle 3)

- DONE: `foFunctionReferenceBaselineBytes` in internal/contractlint/fo_function_reference_invariant_test.go is raised to accommodate the measured surface with a small working margin, in its own commit whose message states what grew and why the estate's budget is being re-set rather than the change trimmed — one constant, nothing else in that file, no other Go file touched.
  Commit f2c5a40e, `1 file changed, 1 insertion(+), 1 deletion(-)`: 122,634 → 123,323. Margin 507 bytes of usable headroom under the strictly-below comparison — the number main carries today (122,634 − 1 − 122,126), so the ratchet re-arms at unchanged tension and funds this member's landing only, not a next train.
- DONE: The `«feedback.route»` edit in first-officer-shared-core.md is verified against z7's LANDED text (z7 merged as PR #540; the branch has been rebased onto it), not the pre-merge copy it was authored against — re-anchor if z7's final round moved that section.
  No re-anchor needed, and this is not a rebase-clean inference: `git show 45f54678 -- .../first-officer-shared-core.md` has hunks at `@@ -94,7` (`«gate.assemble-verdict»` block line) and `@@ -106,11` (State Management) that BRACKET the `«feedback.route»` section without touching it, so the two lines the edit rewrites are the same two lines z7 landed. z7 did not touch `feedback-rejection-flow/SKILL.md` at all. Semantically the edit sits with rather than against z7's landed "Prefer the cheapest check that can fail" ordering — this cut is prose plus a reader, with the command deferred, which is that ordering applied.
- DONE: `go test ./...` passes including TestFOFunctionPromptSurfaceShrinks, and the stage report records the final measured surface, the new baseline, the margin, and that the 0-Go self-check tripped and was resolved by an escalated captain decision.
  `go test ./...` exit 0, 15 packages ok, no FAIL. `TestFOFunctionPromptSurfaceShrinks` PASS; the same assertion was RED at 122,815 vs 122,634 before the constant moved, so the pass is not vacuous — it fails again the moment the surface exceeds 123,322.

### Final measurement

Measured FO prompt surface **122,815** (`FO_FUNCTION_METRICS addresses=0 bytes=122815`); main measures 122,126, so this member's net is **+689** (skill +547, shared-core +142). New baseline **123,323**; largest passing value 123,322; **margin 507**.

Final surface against the declared "~3 prose files, ~17 lines, 0 Go": `git diff --numstat main...HEAD` = **4 files, +18/−8** — `docs/dev/README.md` +2, `feedback-rejection-flow/SKILL.md` +13/−5, `first-officer-shared-core.md` +2/−2, and `fo_function_reference_invariant_test.go` +1/−1. Prose 3 files / 17 lines, exactly as declared; the fourth file and the eighteenth line are the baseline constant.

### The 0-Go self-check: tripped, escalated, approved

This member declared 0 Go LOC, and touching a Go file trips that self-check. It tripped, and that is the check doing its job — it exists to force a recorded decision before Go code enters a prose-only member, not to be waived by whoever is holding the keyboard. The sequence: the prior ensign hit the ratchet twice and STOPPED both times with numbers rather than shaving rule text (cycle 1, 2,903 bytes over; cycle 2, 206 over after the 02av re-scope recovered 94% of it), the FO escalated it as a call about the estate's budget, and the captain approved re-setting the baseline. Grounds recorded in f2c5a40e: the ratchet already collected its value this sprint — four members funded themselves, ~3,400 bytes of genuine redundancy harvested, one unaffordable payload killed, one delivery re-homed — and a duplication scan across all 13 measured files then returned only 110 bytes, because the seam is already harvested. The residual buys no leanness; closing it means cutting rule text. Governance, not gaming.

Nothing else moved: no new committed check, gate, or lint; no test logic; no minted terminology; 02av's parked material untouched. The one residue left deliberately, per the one-constant instruction: the failure message still calls the number the "post-#531 baseline", which is now a stale label.

### Summary

Raised the ratchet baseline to 123,323 in its own commit (f2c5a40e) with the grounds in the message, verified the `«feedback.route»` edit against z7's landed text by reading z7's actual hunks rather than trusting a clean rebase, and confirmed `go test ./...` green with the ratchet included. The four preserved properties survive untouched — the `- Cycle {N}:` leading form that satisfies both the convention and the shipped `feedback-3-cycle-escalation` assertion, the cycle-3 escalation clause, the dev one-liner staying out of the generic skill per AC-7, and the sentence pinning deviation to the approved estimate rather than the prior round. The `### Feedback Cycles` entry above records this round as a RECONFIRM, so the entity keeps practicing its own convention on the round that produced it.

### Roborev branch review (pre-validation, FO-triaged 2026-07-21)

`roborev review --branch --panel branch_final` (job 394, codex, roles correctness+product) returned **changes requested — two medium contract ambiguities**. FO triage under the sprint's posture (declared estimate first, recorded decline for correct-but-disproportionate, no over-building): **both DECLINED with grounds; bw is not re-dispatched.**

1. **"Undefined surface-percentage calculation" (`SKILL.md:28`) — DECLINED.** The finding is answered by line 31, which the reviewer's summary did not account for: the actuals come from "the one-liner the workflow documents for its own surface unit" and deviation is "the actuals over the estimate the ideation gate approved" — i.e. P = actuals ÷ approved-estimate, with the unit deliberately workflow-defined (bytes for the FO surface, net prose lines for a template member, net LOC for a code member). The reviewer's "additions vs deletions vs churn / how file and LOC ratios combine" imports a multi-dimensional metric the convention does not use — it takes one figure from each workflow's own surface one-liner, so within a workflow it is deterministic. Hardcoding the unit to LOC would break the convention for byte-measured surfaces. Acting on it would mint the aggregation machinery this sprint exists to price out.

2. **"No compatibility path for entities without estimates" (anchored `SKILL.md:14`) — DECLINED here, routed upstream.** Names a real gap but not bw's: bw ships the convention that CONSUMES an ideation estimate; the estimate SOURCE (every entity declaring one) is upstream (`js6` stakes-declaration-read-through / the ideation gate). An entity with no estimate produces no deviation figure — graceful degradation, not a failure. Any legacy/migration path is backward-compatibility, which requires an explicit captain decision by standing rule; not folded in unilaterally. Surfaced to the captain.

Both are the sprint's own thesis in miniature: a reviewer asking a deliberately prose-only, enforcement-deferred recording convention to become a fully-specified computed tolerance gate. bw's scope (command, refusal machinery, escalation backstop DEFERRED) is a captain decision; the review does not disturb it.

## Stage Report: validation (detached adversarial audit + AC-1 value proof, FO-driven under the conn 2026-07-21)

**Verdict: CLEAR FOR MERGE — refuted nothing material.** Two independent detached validators on throwaway checkouts (correctness+scope lens; live-lane+AC lens) both returned refuted-nothing-material; the only non-NOTHING findings were two CORRECT-BUT-DISPROPORTIONATE items, both captain-sanctioned (the AC-2/4/7 one-off greps are blessed by the 2026-07-20 prose-grep ruling; AC-1's value proof was a validation-stage deliverable, now run — below).

**Scope exactly as declared:** 4 files, one Go line (baseline constant only). No fifth file, no second Go change, no new lint/gate/test func, no minted terminology.

**Four preserved properties — all SURVIVED, each independently checked:**
- P1 (`- Cycle {N}:` satisfies the shipped `^- Cycle \d+:` live assertion): PROVEN, not asserted — a validator fed a bw-format escalation entity to the REAL `assertThirdCycleEscalation` + `feedbackCycleEntry` regex (internal/ensigncycle/shared_assertions_impl_test.go:43) in a throwaway in-package test → PASS, 3 entries + escalation matched, no fixture edited. The claude/codex/pi feedback-3-cycle-escalation lanes will not red on format grounds.
- P2 (cycle-3 escalation clause retained): SURVIVED; contractlint TestFOLocalOrderedProceduresPreserved anchors "cycle 3".
- P3 (dev numstat one-liner OUT of the generic skill, AC-7): SURVIVED; the one-liner lives only in docs/dev/README.md, not the skill.
- P4 (deviation pinned to approved estimate, not prior round): SURVIVED and shown load-bearing by the AC-1 negative case below.

**Ratchet:** re-arm proven by padding the file to a real RED (137215 > 123323) then reverting to clean; +689 baseline bump exactly matches +689 surface growth (identical 508 B slack to main). Disciplined, not loose. `go test ./...` green.

### AC-1 (VALUE) — PASS, on real e6j data

Run as a one-off exercise (output here, nothing committed — see the fixture decision below). e6j's real per-cycle history was git-reconstructed from `archive/2026-07-14/e6j-state-ready-rejected` (10 commits = 10 roborev cycles) and cross-checked against the forensics evidence and the live transcript (msg 203669 corroborates cycle-2 = 1,068 additions).

- Conforming entries, 2-file/40-LOC approved estimate:
  - `- Cycle 1: REJECTED — roborev; surface 6 files/630 LOC vs estimate 2 files/40 LOC (300% files / 1575% LOC); AC unchanged`
  - `- Cycle 2: REJECTED — roborev; surface 14 files/1068 LOC vs estimate 2 files/40 LOC (700% files / 2670% LOC); AC unchanged`
- **Round-2 deviation = 700% (files) / 2670% (LOC), both ≥ 200% → PASS.** Runaway visible at cycle 2; e6j's real history reframed only at cycle 10 — surfaced 8 cycles earlier, meeting AC-1's "by round 2" DoD target. (The ideation spike used a conservative synthetic 7-file/350%; real data is ~2× more extreme.)
- Falsification cases (all behaved as AC-1 requires): (a) within-tolerance 3-file/60-LOC = 150% → design-reset NOT tripped; (b) omitted actuals/estimate → check FAILS (uncomputable); (c) measured vs the PRIOR ROUND instead of the approved estimate → caps at 170% and the runaway hides every round, versus 2670% at round 2 against the fixed estimate — this is why "vs approved estimate, never the prior round" is load-bearing (matches the forensics finding that every observed runaway was contract-legal round by round).
- AC-5 drift half: narrowed entry carries the note, control does not, a mislabeled entry FAILS. Verified.

### Fixture-commit decision (recorded FO call, captain may override)

AC-1's text names "a checked-in fixture entity." The FO did NOT commit a fixture. Grounds: the captain struck AC-1's committed check at ideation (gate-attempt bw-ideation-3, "one-off script, no committed check"); bw is scoped prose-only, 0 committed tests/gates/lints/fixtures; and a committed entity whose sole purpose is to be read by a one-off script is the fabricated-rigor shape this sprint prices out. The value is proven falsifiably by the one-off exercise above with real data. If the captain wants a durable fixture, it is a small add — surfaced, not silently skipped.

### Required lanes before merge

bw touches `skills/**/references/**` (the host-neutral contract core), so ALL THREE host live lanes (claude-live, codex-live, pi-live) plus the deterministic lanes are required green — deterministic-only merge is not permitted. The P1 format-compat proof predicts the live lanes stay green; any red is a captain reconciliation decision, never a silent fixture edit.
