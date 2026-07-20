---
id: bwr6j6edkmfx5sbz73cr2952
title: Feedback-cycle record and design-reset gate — binary-owned count, diff-growth refusal, reframe routing
status: ideation
source: "captain (2026-06-04) — forked from xa (feedback-guarantee-binary-gate) per the roadmap-the-decision + separate-build-task call. xa's ideation determined Candidate 1 (3-cycle escalation) is mechanizable via a dedicated cycle-record command (a spike disproved a --set status guard) and Candidate 2 (budget-probe) is not. This task SHIPS the Candidate-1 guard; xa closed as a roadmap decision."
score: "0.30"
started: 2026-07-20T03:29:33Z
completed:
verdict:
worktree:
issue:
sprint: 0260-proportionality
group: reframe
gates:
  version: 1
  current:
    gate: gate:docs-dev:bw:ideation
    attempt: gate-attempt:bw-ideation-1
  records:
    - id: gate:docs-dev:bw:ideation
      stage: ideation
      current-attempt: gate-attempt:bw-ideation-1
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
          application:
            action: feedback
            target-stage: ideation
            state: pending
          note: "Subspace advisory float; four captain annotations included by id in the resolution. Annotation 1 is FO-owned (briefing packaging), 2-4 routed to the worker; next attempt opens at re-presentation."
---

Turn the feedback-rejection correction loop from a prose-only cycle count into a measured, calibrated loop. At ideation the entity declares an **expected surface** as part of the captain-approved design; every correction round records its **actuals** into durable `### Feedback Cycles` state; the **deviation** of actuals from the captain-approved estimate is **narrated** at each re-dispatch decision point. A new `status --record-feedback-cycle` command owns the append + count + actuals + deviation + the stdout round-state line.

**Narration is the guard, not a wall.** Per the ideation-gate bounce (2026-07-20), the first cut ships only cheap detection: the record command's measurement and the two skill-prose swaps that make the FO read the round state before re-dispatching. The hard backstop (a durable escalation marker that refuses a further auto-bounce) ships **only if live drives show the narration being ignored** — it is deferred, not designed out. Hard dispatch-refusal on raw diff growth is demoted for good: it false-refuses a legitimately growing fix, and a non-zero `dispatch build` exit is read by the FO contract as an infra failure that triggers Break-Glass **manual** dispatch (`fo-dispatch-core.md:145`), so a build-side refusal would be routed around rather than honored. The frontmatter `title:` still names "diff-growth refusal"; the body is authoritative.

## First cut vs deferred

**Expected surface (this entity's own ideation estimate, in the form the convention it ships asks for):** ~2 Go files in `internal/status` (1 new handler `feedback_cycle.go` + the `--record-feedback-cycle` flag branch in `native_runner.go`) + 1 new test file, ~150 LOC source / ~200 LOC test; 2 skill-prose files (`feedback-rejection-flow`, `first-officer-shared-core`), ~6 lines swapped; 1 dev-README line. Declared tolerance: 2× — reconfirm or reframe if implementation crosses ~4 Go files / ~400 changed lines, or if the skill-prose file count grows.

**First cut (this entity):**
- The record command core: append a `### Feedback Cycles` entry, section-scoped count, per-round actuals + deviation vs the declared estimate, and the stdout round-state line.
- The two skill-prose swaps: `feedback-rejection-flow` steps (invoke the command, read the narration, weigh a design-reset decision) and the `first-officer-shared-core` `«feedback.route»` lines.
- One line in the dev README ideation stage-def so the declared-surface baseline exists (its template propagation is the `template` member's job, not this entity's).
- **Recommended (my call on the gate's ask):** the AC-drift digest — one entry field + one narration clause, riding the same command (see `### AC-drift detection`).

**Deferred to a later cut (ships only if live drives show narration ignored):**
- The escalation marker + refuse-further-auto-bounce + `--force` machinery, and its `feedback-escalate` schema field.
- The `present-gate` surface-deviation evidence line.
- The template ideation-scaffold `Outputs` propagation (owned by the `template` member).
- The CLI docs page for the command.
- The enforcement half of AC-drift: treating an AC-weakening edit as a design-reset event requiring captain sign-off (the AC-narrowing backstop; see `### AC-drift detection`).

## Problem

- **The runaway loops were contract-legal round by round.** All four HIGH incidents in the 0260 forensics (`_evidence/0260-agent-derail-forensics/synthesis.md`) — e6j (2-defect fix → 10 roborev cycles, 26 files / +3,373, PR closed), dp (one-paragraph fix → 4-cycle ladder, discarded, ~38.5h), task-91 (16 roborev panels, own round-limit bypassed), 7h (harness repaired twice before park) — passed every round against the only baseline available: the prior round's accident. No single round screamed. The baseline has to be the entity's own captain-approved intent, not the last cycle's overrun.
- **AC-narrowing under validation pressure** (synthesis addendum, 2026-07-20, the 0.25.1 release; the addendum calls it repair-forward's paperwork twin). When validation correctly found a value claim unproven, the task narrowed its AC until a weaker claim passed — a real rejection converted into a paperwork pass; the failure then reproduced live. The gate cross-check compares against the CURRENT AC text, so silent narrowing defeats it by construction. A calibrated loop must make AC-drift across rounds machine-visible, the same way it makes surface-growth visible.
- **In-stage rounds have no durable record.** `### Feedback Cycles` tracks only cross-stage gate bounces (dp, 7h). e6j's 10 roborev cycles and task-91's 16 panels were in-stage (roborev-at-end-of-implementation, never crossing a gate), so cycle tracking never saw them. Both loop shapes must land in one section.
- **xa's determination stands** (`_archive/feedback-guarantee-binary-gate`): the cycle count IS durable on-disk state (section-scoped count deterministic and tamper-evident — spike: ~25 lines, ignores a `Cycle N` line in a sibling section); a `--set status={feedback-to-target}` guard FALSE-FIRES (the disambiguating `is_feedback_reflow` lives on the dispatch-build input path, `build.go:299`, not as a `--set` field or durable state); the correct hook is the cycle-record WRITE, unambiguously a bounce event.
- The prose-only guarantee's ceiling is "the wording is present"; its drift mode is the infinite reject→re-implement→reject loop.

## Proposed approach

### Generic principle (workflow-agnostic — lands in contract/skill prose)

1. **Declare expected surface at ideation.** The gated design states the surface it expects to touch, in the workflow's own unit — dev: "~x files, ~y LOC"; a research workflow: "~10 external docs". Captain-approved at the ideation gate → this is the baseline.
2. **Record per-round actuals.** Every correction round (whatever the loop shape) appends its actuals to `### Feedback Cycles`.
3. **Deviation vs the declared baseline** (not the prior cycle), beyond a **declared tolerance** (contract default, entity may override), is **narrated** and prompts a **reconfirm-or-reframe** judgment at the decision point.
4. **Narration is the guard.** Round state — `round N, findings k material / m deferred, surface at P% of estimate` — is rendered where the FO decides whether to re-dispatch. Awareness at the decision point, not a wall.
5. **Escalation is the backstop — DEFERRED.** A hard stop that refuses a further auto-bounce ships only if live drives show narration being ignored. It is not in the first cut.

### Dev-specific realization (lands in dev template + binary)

- **Surface unit:** files touched / LOC.
- **Actuals capture:** the record command reads the entity's `worktree:` and runs `git diff --numstat {merge-base}..HEAD` — cumulative surface of the entity's work vs its branch point, deterministic, no new round-store (roborev has none; rounds are Git-versioned — spike-proven below). A generic `--surface {files,loc}` override lets a non-dev caller pass actuals.
- **roborev round hook:** the in-stage roborev loop invokes the record command when it re-panels after material findings, so in-stage rounds land in the same `### Feedback Cycles` section as cross-stage bounces.

### Both loop shapes, one record (trigger paths named)

- **Cross-stage feedback bounce** (dp, 7h): triggered by the FO's `feedback-rejection-flow` when a gate recommends REJECTED and routes to the `feedback-to` target — the FO invokes the record command. Generic; lands in the skill.
- **In-stage review round** (e6j's 10 roborev cycles, task-91's 16 panels): triggered by the implementation-stage roborev loop re-paneling after material findings within one stage — the loop driver invokes the record command per re-panel. Dev realization; lands in the implementation-stage roborev-loop prose.

Both write the same section: `### Feedback Cycles` becomes the single durable record of correction rounds. Making in-stage roborev rounds record here is precisely what would have surfaced e6j at all.

### The record command (first cut)

`spacedock status --record-feedback-cycle {slug}`:
1. Appends a timestamped `- Cycle N:` entry to `### Feedback Cycles` (creating the section if absent), stamping the round's actuals, the declared estimate, the computed deviation%, and the AC-digest (below).
2. Computes N section-scoped (a `Cycle N` line in a sibling section does not inflate it — xa spike).
3. Emits the round-state line on stdout.

It does NOT, in the first cut, stamp an escalation marker or refuse — those are deferred. It composes proven machinery (`section_read.go` body-parse, `mutate.go` `atomicWrite`, `merge.go arm` for the later marker); the only net-new byte-writes are appending the `- Cycle N:` body line and the git-diff capture (`runGitCmd`).

### Narration render (first cut)

- The record command emits the round-state line on stdout (byte-observable, unit-tested).
- `feedback-rejection-flow` renders/reads it at the re-dispatch decision point (the FO's own awareness before routing another repair — not `output.prompt`, which is forwarded verbatim).
- The `present-gate` evidence line is **deferred**.

### Design-reset routing (narration-driven in the first cut)

When the narrated deviation crosses tolerance, the FO weighs a **design-reset decision** — reconfirm the estimate / re-scope / park / escalate — instead of automatically routing the next repair. In the first cut this is a judgment prompted by narration (the skill step), not an enforced halt; a recorded reframe re-baselines (the new estimate becomes the baseline). The enforced backstop that forces the halt when narration is ignored is deferred.

### AC-drift detection (recommended for the first cut — the gate's ask #4)

The record command already reads the entity body and writes a per-round entry. Riding that, each entry also stamps a **digest of the `## Acceptance criteria` section** (normalized text → short hash), and the round-state line flags **`AC CHANGED since cycle 1`** when the digest differs from the first recorded round's. This makes AC-narrowing under validation pressure (synthesis addendum) machine-visible from two durable entries — exactly the pattern surface-deviation uses.

**Recommendation: include the DETECTION+NARRATION in the first cut; defer the ENFORCEMENT.** Reasoning:
- It is the narration half of the fix, and the gate's own logic ships cheap detection first, enforcement only if ignored. The two derail directions (grow-diff, narrow-AC) then get symmetric detection from the same command.
- Cost is genuinely marginal: it rides the existing `readSections` call (the AC heading span comes back for free) and the existing entry-write — adding a hash, one entry field, and one narration clause. No new command, gate, frontmatter field, or schema change.
- It covers a live-reproduced failure the fresh addendum names as repair-forward's twin.
- **Honest limits, recorded:** the baseline is the first recorded round, so it catches narrowing DURING the correction loop (the 0.25.1 shape) but not narrowing during initial implementation before any rejection — that needs an ideation-gate AC snapshot, deferred. The ENFORCEMENT half — lesson (a): an AC-weakening edit after a rejection is a design-reset event requiring captain sign-off — is prose+backstop and defers with the escalation machinery.

If the gate prefers the tightest possible core, this AC (AC-6) is severable: defer it whole without disturbing AC-1..AC-5.

### Per-mechanism justification (value AC served / simplest alternative / why insufficient)

- **Record command (append + count + actuals + deviation):** serves AC-1. Alt: FO prose-tracking (status quo). Insufficient — a prose count carries no actuals, so deviation is uncomputable; ceiling is "wording present"; in-stage rounds go unrecorded.
- **Deviation vs the declared estimate (not the prior cycle):** serves AC-1. Alt: diff-growth vs the prior cycle (the superseded design). Insufficient — every e6j round passed round-by-round; the prior-cycle baseline never screams (spike counterfactual: +5 files reads as ordinary growth).
- **AC-digest (AC-6):** serves AC-1's calibration by covering the narrow-AC direction. Alt: rely on the gate's live cross-check of current AC text. Insufficient — the addendum proves that cross-check is defeated by construction when the AC is silently narrowed; the digest makes the edit durable and comparable.
- **Escalation marker + refuse (DEFERRED):** would serve a backstop. Alt A: `--set status` guard — xa spike: false-fires. Alt B: `dispatch build` refusal — Break-Glass routes it around + false-refuses growing fixes. If it ships later, the record-command write is the only correct hook (bounce-unambiguous AND an intentional decision-gate exit).
- **Narration as guard (not a wall):** serves AC-1's awareness-at-the-decision-point. Alt: a hard wall. Insufficient in the first cut — override ceremony on legitimate growth (captain demotion); narration is tried first.

## Landing-spot map (generic principle vs dev realization; first-cut vs deferred)

| Piece | Kind | Cut | Lands in |
|---|---|---|---|
| Declare expected surface at ideation (one line) | generic | first | dev README ideation stage-def |
| Template propagation of that line | generic | deferred | `template` member |
| files/LOC as the surface unit | dev | first | record command |
| per-round actuals; deviation vs baseline; reconfirm-or-reframe; narration render | generic | first | `feedback-rejection-flow/SKILL.md`; `first-officer-shared-core.md` effect line |
| git-diff actuals capture (`--numstat {merge-base}..HEAD`) | dev | first | `internal/status` record command (`runGitCmd`) |
| roborev round hook (in-stage rounds record) | dev | first | dev README implementation-stage roborev-loop prose |
| section-scoped count + append; stdout round-state line | dev/binary | first | `internal/status` (`section_read.go`, `mutate.go` `atomicWrite`) |
| AC-section digest + drift narration | dev/binary | first (recommended) | `internal/status` record command |
| present-gate surface-deviation evidence line | generic | deferred | `present-gate/SKILL.md` |
| escalation marker + refuse + `--force` + `feedback-escalate` schema field | dev/binary | deferred | `internal/status` (`merge.go arm`, `handlers.go`, embedded schema) |
| AC-weakening = design-reset requiring captain sign-off | generic | deferred | `feedback-rejection-flow/SKILL.md` + backstop |

## Riskiest-mechanism spike (done first, per ideation policy)

**Claim under test:** can per-round actuals be captured from durable evidence, and deviation computed vs the declared estimate (not the prior cycle)? Exercised end-to-end against real git state (throwaway repo, results recorded — not asserted):

- **Cumulative actuals from durable git** via `git diff --numstat {BASE}..HEAD`: round 1 = **2 files**, round 2 = **7 files** (base = the pre-fix commit; entity carries `surface-estimate: 2 files, 40 LOC`).
- **Deviation vs the declared estimate** (2 files): round 1 = **100%**, round 2 = **350%** → crosses a 200% tolerance **at round 2** (the e6j "scream by round 2" target from the DoD).
- **Counterfactual vs the prior cycle** (the accident baseline e6j effectively used): round 2 − round 1 = **+5 files**, reads as ordinary growth and never screams — proving the baseline must be the estimate, not the accident.
- **Deviation recomputes offline from two durable `### Feedback Cycles` text entries**, no live process — confirming the record is a reader/measurement, not a runtime bouncer.

**Determination:** the estimate/actual semantics compose deterministic git capture + section-scoped counting (xa-proven) + text digest. No further spike needed; the throwaway seeds the implementation's first unit test.

## Acceptance criteria

**AC-1 (value) — Fed the archived e6j per-round surface shape (2-file / 40-LOC estimate; 10 cycles ending 26 files / +3,373), the computed cumulative-surface deviation crosses the declared tolerance no later than round 2, and the emitted round-state line reads ≥ 200% of estimate at round 2.**
Verified by: a Go unit test in `internal/status` feeding the e6j fixture cycle boundaries; asserts the computed deviation% and the round-2 narration string byte-for-byte. Independent baseline that moved the wrong way: e6j's real history surfaced nothing until round 10 (archived, 26 files / +3,373). This is the outcome the entity exists for — the runaway made visible at round 2.

**AC-2 — `status --record-feedback-cycle {slug}` owns the `### Feedback Cycles` append, a section-scoped count, and the per-round actuals + deviation stamp (dev: files/LOC from git; generic: `--surface` override), emitting the round-state line on stdout.**
Verified by: a Go unit test driving the command against a temp git-backed entity — first invoke appends the section + a cycle-1 entry carrying actuals + estimate + deviation%; Nth invoke appends cycle-N; a `Cycle N` line in a sibling section does not inflate the count; stdout carries the round-state line. Serves AC-1 (the durable actuals are the measurement substrate).

**AC-3 (deferred-scope note) — The hard backstop (escalation marker + refuse-further-auto-bounce + `--force`, and its `feedback-escalate` schema field) is NOT in this cut.** It ships only if live drives show the narration being ignored. Recorded here so the boundary is explicit; when it ships, the record-command write is the hook (not `--set status`, not `dispatch build`), for the reasons in the per-mechanism justification.

**AC-4 — The record command emits the round-state line on stdout, and `feedback-rejection-flow` reads it at the FO's re-dispatch decision point.**
Verified by: the stdout line is byte-observable (unit-tested under AC-2); a section-scoped presence oracle over `feedback-rejection-flow` asserts the re-dispatch step invokes the command and reads its line. The `present-gate` render is deferred. The FO ACTING on the narration is gq's live scenario, not this AC.

**AC-5 — `feedback-rejection-flow` and `first-officer-shared-core` invoke the record command and narrate instead of prose-tracking; the dev README ideation stage-def gains the one-line expected-surface declaration.**
Verified by: a section-scoped presence oracle over the skill prose (`skill_text_test.go` `sectionAfter` pattern) asserting the two swaps; the dev README line is checkable by reading the stage-def. Template propagation of the line is the `template` member's, not this AC.

**AC-6 (recommended for the first cut; severable) — Each round entry stamps a digest of the `## Acceptance criteria` section, and the round-state line flags AC-drift when the digest differs from the first recorded round.**
Verified by: a Go unit test — record a round, narrow the AC text, record another round, assert the second entry's digest differs and the stdout line carries the drift flag; a negative control (unchanged AC → no flag). Serves AC-1's calibration by covering the narrow-AC direction. Enforcement (captain sign-off on AC-weakening) is deferred.

## Test plan

- **Unit (Go, `internal/status`):** drive the command against temp git-backed entities; assert append / count / actuals / deviation / stdout line / section-scoping / AC-digest + drift flag. Byte-observable over the on-disk file + stdout — same altitude as `archive_guard_test.go`. Composes proven machinery (`section_read.go`, `mutate.go` `atomicWrite`); net-new: the body-append and the git-diff capture (`runGitCmd`). Cost: low.
- **Value replay (Go fixture):** the e6j cycle-boundary fixture → AC-1. The DoD's "unit test fed that fixture shape" retargets from the demoted dispatch-refusal guard to the record command's deviation computation.
- **AC-drift control (Go):** narrow-AC-across-rounds → digest differs + flag; unchanged → no flag → AC-6.
- **Presence oracle (offline):** AC-4 / AC-5 skill-prose checks ride the `skill_text_test.go` `sectionAfter` pattern (currently only in a worktree copy, `.worktrees/audit-bq-release-gate/skills/integration/skill_text_test.go:277`; porting to the main tree is part of implementation).
- **Deferred (not in this cut's test plan):** the escalation-marker refuse + `--force` + negative control, and the `present-gate` render — they ship with the backstop.
- **High-stakes → detached adversarial audit before merge:** this edits the `status` mutation/guard paths (a named high-stakes surface).

## Doc diffs (first cut only; ideation proposes, implementation applies)

**`skills/feedback-rejection-flow/SKILL.md`** — steps 2-3:

> - Before (step 2): `2. Track cycles in `### Feedback Cycles` in the entity body.`
> - After: `2. Record the round with `${SPACEDOCK_BIN:-spacedock} status --record-feedback-cycle {slug}` — it appends the `### Feedback Cycles` entry (actuals + estimate + deviation% + AC-digest) and returns the round-state line.`
> - Before (step 3): `3. On cycle 3, escalate to the human instead of another round.`
> - After: `3. Read the round-state line (`round N, surface at P% of estimate`; `AC CHANGED` if the AC drifted). When deviation is beyond tolerance or the AC narrowed, weigh a design-reset decision (reconfirm / re-scope / park / escalate) instead of an automatic next round.`

**`skills/first-officer/references/first-officer-shared-core.md:102-103`** — the `«feedback.route»` effect/done-when:

> - Before: `… track `### Feedback Cycles`, escalate on cycle 3, …` / `… (or escalated at cycle 3).`
> - After: `… record each round with `status --record-feedback-cycle` (which owns the count and stamps actuals + deviation), read the narrated round state, …` / `… (or a recorded design-reset decision when deviation is beyond tolerance).`

**Dev README ideation stage-def** — add to the ideation `Outputs` list (one line; template propagation is the `template` member's):

> `The task body declares an expected surface — the files/LOC (or the workflow's own unit) it expects to touch — as part of the gated design; the ideation gate approves it as the baseline the correction loop calibrates against.`

Deferred doc changes (not proposed here): the `present-gate` evidence line, the template ideation-scaffold `Outputs` propagation, and the CLI docs page.

## Boundary and notes

- **gq (`feedback-nonhappy-live-coverage`) owns the live half:** that the FO acts on the narration / obeys a (later) refusal is FO-LLM behavior with no in-process Go seam — its scenario, not an AC here.
- **Out of scope:** the budget-probe fail-safe (xa Candidate 2, non-mechanizable); any `dispatch build` refusal (demoted — the Break-Glass routing makes it the wrong hook).
- **Leanness (0250/0260):** the count prose moves OUT of the always-on FO contract INTO the binary (boot-resident prose shrinks); the narration and declare-surface additions are lazy-loaded (`feedback-rejection-flow` loads at the rejection point; the one dev-README line rides the stage-def, not boot). Deferring the backstop, present-gate line, template propagation, and CLI page keeps this cut's net byte delta small; reported at implementation.
- **Lane calibration:** conventions settled here (the `### Feedback Cycles` entry shape carrying actuals+estimate+deviation+AC-digest; the round-state line format) are what ek and w0 build against after this gate; ht runs parallel.
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
