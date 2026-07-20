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
          state: open
          briefing:
            id: briefing:bw-ideation-1
            digest: sha256:3dcd7a42c4ee7899f8c074e646a284baf77d15cfeb198e3bf1a7e8bf852d8afa
---

Turn the feedback-rejection correction loop from a prose-only cycle count into a measured, calibrated loop. At ideation the entity declares an **expected surface** as part of the captain-approved design; every correction round records its **actuals** into durable `### Feedback Cycles` state; the **deviation** of actuals from the captain-approved estimate is **narrated** at each re-dispatch decision point and gate presentation. The banked 3-cycle escalation survives as the backstop hard stop, now driven by the declared threshold rather than a blind count. A new `status --record-feedback-cycle` command owns the append + count + actuals + escalation marker as the measurement substrate — but **narration, not a dispatch-time wall, is the guard.**

**This supersedes the earlier framing** (hard dispatch-refusal on raw diff growth), per captain direction 2026-07-20. Hard refusal is demoted for two reasons: (1) it false-refuses a legitimately growing fix, spawning override ceremony; (2) a non-zero `dispatch build` exit is read by the FO contract as an infra failure that triggers Break-Glass **manual** dispatch (`fo-dispatch-core.md:145`, `claude-fo-dispatch.md:48`), so a build-side refusal would be routed *around* rather than honored. The frontmatter `title:` still names "diff-growth refusal" (frontmatter is not edited at this stage); the body below is authoritative for the design.

## Problem

- **The runaway loops were contract-legal round by round.** All four HIGH incidents in the 0260 forensics (`_evidence/0260-agent-derail-forensics/synthesis.md`) — e6j (2-defect fix → 10 roborev cycles, 26 files / +3,373, PR closed), dp (one-paragraph fix → 4-cycle ladder, discarded, ~38.5h), task-91 (16 roborev panels, own round-limit bypassed), 7h (harness repaired twice before park) — passed every round against the only baseline available: the prior round's accident. No single round screamed because each increment looked ordinary. The baseline has to be the entity's own captain-approved intent, not the last cycle's overrun.
- **In-stage rounds have no ledger at all.** `### Feedback Cycles` tracks only cross-stage gate bounces (dp, 7h). e6j's 10 roborev cycles and task-91's 16 panels were **in-stage** (roborev-at-end-of-implementation, never crossing a gate), so cycle tracking never saw them — the runaway ran with no durable record. Any calibration must capture both loop shapes into one ledger.
- **xa's determination stands** (`_archive/feedback-guarantee-binary-gate`): the cycle count IS durable on-disk state (section-scoped count is deterministic and tamper-evident — spike: ~25 lines, ignores a `Cycle N` line in a sibling section); a `--set status={feedback-to-target}` guard FALSE-FIRES (the disambiguating `is_feedback_reflow` lives on the dispatch-build input path, `build.go:299`, not as a `--set` field or durable state); the correct hook is the cycle-record WRITE, which is unambiguously a bounce event.
- The prose-only guarantee's ceiling is "the wording is present"; its drift mode is the infinite reject→re-implement→reject loop burning tokens.

## Proposed approach

### Generic principle (workflow-agnostic — lands in contract/skill prose)

1. **Declare expected surface at ideation.** The gated design states the surface it expects to touch, in the workflow's own unit — dev: "~x files, ~y LOC"; a research workflow: "~10 external docs". Captain-approved at the ideation gate → this is the baseline.
2. **Record per-round actuals.** Every correction round (whatever the loop shape) appends its actuals to `### Feedback Cycles`.
3. **Deviation vs the declared baseline** (not the prior cycle), beyond a **declared tolerance** (contract default, entity may override), triggers **reconfirmation-or-reframe**.
4. **Narration is the guard.** Round state — `round N/limit, findings k material / m deferred, surface at P% of estimate` — is rendered at the FO's re-dispatch decision point and in gate presentations. Awareness at the decision point, not a wall.
5. **Escalation is the backstop.** The banked 3-cycle hard stop survives as a last resort, now tripped by the declared threshold (tolerance crossed) as well as the bare cycle count.

### Dev-specific realization (lands in dev template + binary)

- **Surface unit:** files touched / LOC.
- **Actuals capture:** the record command reads the entity's `worktree:` and runs `git diff --numstat {merge-base}..HEAD` — cumulative surface of the entity's work vs its branch point, deterministic, no new round-store (roborev has none; rounds are Git-versioned — spike-proven below). A generic `--surface {files,loc}` override lets a non-dev caller pass actuals when there is no git surface.
- **roborev round hook:** the in-stage roborev loop invokes the record command when it re-panels after material findings, so in-stage rounds land in the **same** `### Feedback Cycles` ledger as cross-stage bounces.

### Both loop shapes, one record (trigger paths named)

- **Cross-stage feedback bounce** (dp, 7h): triggered by the FO's `feedback-rejection-flow` when a gate recommends REJECTED and routes to the `feedback-to` target — the FO invokes the record command. Generic; lands in the skill.
- **In-stage review round** (e6j's 10 roborev cycles, task-91's 16 panels): triggered by the implementation-stage roborev loop re-paneling after material findings within one stage — the loop driver invokes the record command per re-panel. Dev realization; lands in the implementation-stage roborev-loop prose.

Both write the same section: `### Feedback Cycles` becomes the single durable ledger of correction rounds regardless of loop shape. Making in-stage roborev rounds record here is precisely what would have surfaced e6j at all.

### The record command

`spacedock status --record-feedback-cycle {slug}` (composes proven machinery; see Test plan):
1. Appends a timestamped `- Cycle N:` entry to `### Feedback Cycles` (creating the section if absent), stamping the round's actuals, the declared estimate, and the computed deviation%.
2. Computes N section-scoped (a `Cycle N` line in a sibling section does not inflate it — xa spike).
3. On the threshold — deviation beyond the declared tolerance **or** the banked 3rd cycle — stamps a queryable `feedback-escalate: cycle-N` frontmatter marker (the `mod-block` idiom, `merge.go:arm`) **and** refuses a further auto-bounce (exit 1, `--force` overrides with a warning, the `handlers.go:156-171` idiom).

### Narration render

- The record command emits the round-state line on stdout.
- `feedback-rejection-flow` renders it at the re-dispatch decision point (the FO reads it before routing another repair — not `output.prompt`, which is forwarded verbatim, but the FO's own decision-point awareness).
- `present-gate` renders it as a surface-deviation **evidence** line (rule 11: "evidence, not label"), sourced from this run's `### Feedback Cycles`, never an inherited label.

### Design-reset routing

When deviation crosses tolerance or the backstop fires, the sanctioned exit is a **recorded design-reset decision** — reconfirm the estimate / re-scope / park / escalate — written into the entity, not another repair dispatch. A recorded reframe re-baselines (the new estimate becomes the baseline; the prior overrun does not).

### Per-mechanism justification (value AC served / simplest alternative / why insufficient)

- **Record command (append + count + actuals):** serves AC-1. Alt: FO prose-tracking (status quo). Insufficient — a prose count carries no actuals, so deviation is uncomputable; ceiling is "wording present"; in-stage rounds go unrecorded.
- **Deviation vs the declared estimate (not the prior cycle):** serves AC-1. Alt: diff-growth vs the prior cycle (the superseded merged-scope design). Insufficient — every e6j round passed round-by-round; the prior-cycle baseline never screams (spike counterfactual: +5 files reads as ordinary growth).
- **Escalation marker + refuse, on the record-command path:** serves AC-3 (backstop). Alt A: a `--set status` guard — xa spike proved it false-fires (no bounce signal in `--set`). Alt B: a `dispatch build` refusal — Break-Glass routes it around as infra failure AND it false-refuses growing fixes (captain demotion). The record-command write is the only hook that is both bounce-unambiguous AND an intentional decision-gate exit the FO reads (not an infra-failure exit).
- **Narration as guard:** serves AC-1's awareness-at-the-decision-point. Alt: a hard wall. Insufficient — override ceremony on legitimate growth (captain demotion).

## Landing-spot map (generic principle vs dev realization)

| Piece | Kind | Lands in |
|---|---|---|
| Declare expected surface at ideation | generic | ideation stage-def in the dev template (refit propagates); the dev README instance below |
| files/LOC as the surface unit | dev | dev template ideation scaffold + record command |
| per-round actuals recorded; deviation vs declared baseline; reconfirm-or-reframe | generic | `feedback-rejection-flow/SKILL.md`; `first-officer-shared-core.md` effect line |
| git-diff actuals capture (`--numstat {merge-base}..HEAD`) | dev | `internal/status` record command (built on `runGitCmd`, `handlers.go:611`) |
| roborev round hook (in-stage rounds record) | dev | dev README implementation-stage roborev-loop prose |
| narration render at re-dispatch decision point | generic | `feedback-rejection-flow/SKILL.md` |
| narration render in gate presentation | generic | `present-gate/SKILL.md` template |
| section-scoped count + append; escalation marker + refuse + `--force` | dev/binary | `internal/status` (`section_read.go`, `mutate.go`/`merge.go arm`, `handlers.go`) |
| escalation as backstop (banked 3-cycle, threshold-driven) | generic | `feedback-rejection-flow/SKILL.md`; marker enforced in binary |

## Riskiest-mechanism spike (done first, per ideation policy)

**Claim under test:** can per-round actuals be captured from durable evidence, and deviation computed vs the declared estimate (not the prior cycle)? Exercised end-to-end against real git state (throwaway repo, results recorded — not asserted):

- **Cumulative actuals from durable git** via `git diff --numstat {BASE}..HEAD`: round 1 = **2 files**, round 2 = **7 files** (base = the pre-fix commit; entity carries `surface-estimate: 2 files, 40 LOC`).
- **Deviation vs the declared estimate** (2 files): round 1 = **100%**, round 2 = **350%** → crosses a 200% tolerance **at round 2** (the e6j "scream by round 2" target from the DoD).
- **Counterfactual vs the prior cycle** (the accident baseline e6j effectively used): round 2 − round 1 = **+5 files**, reads as ordinary growth and never screams — proving the baseline must be the estimate, not the accident.
- **Deviation recomputes offline from two durable `### Feedback Cycles` text entries**, no live process — confirming the record is a reader/measurement, not a runtime bouncer.

**Determination:** the estimate/actual semantics compose deterministic git capture + section-scoped counting (xa-proven, ~25 lines) + frontmatter stamp (`mod-block`-proven). No further spike needed; the throwaway seeds the implementation's first unit test. The one net-new byte-write — appending a `- Cycle N:` line to the body — has no existing helper and rides `atomicWrite` (`mutate.go:244`).

## Acceptance criteria

**AC-1 (value) — Fed the archived e6j per-round surface shape (2-file / 40-LOC estimate; 10 cycles ending 26 files / +3,373), the computed cumulative-surface deviation crosses the declared tolerance no later than round 2, and the emitted round-state line reads ≥ 200% of estimate at round 2.**
Verified by: a Go unit test in `internal/status` feeding the e6j fixture cycle boundaries; asserts the computed deviation% and the round-2 narration string byte-for-byte. Independent baseline that moved the wrong way: e6j's real history surfaced nothing until round 10 (archived, 26 files / +3,373). This is the outcome the entity exists for — the runaway made visible at round 2.

**AC-2 — `status --record-feedback-cycle {slug}` owns the `### Feedback Cycles` append, a section-scoped count, and the per-round actuals stamp (dev: files/LOC from git; generic: `--surface` override).**
Verified by: a Go unit test driving the command against a temp git-backed entity — first invoke appends the section + a cycle-1 entry carrying actuals + estimate + deviation%; Nth invoke appends cycle-N; a `Cycle N` line in a sibling section does not inflate the count. Serves AC-1 (the durable actuals are the measurement substrate).

**AC-3 — On the threshold (deviation beyond the declared tolerance OR the banked 3rd cycle) the command stamps a queryable `feedback-escalate` frontmatter marker and refuses a further auto-bounce.**
Verified by: a Go unit test asserting the threshold invoke stamps the marker AND the next auto-bounce attempt exits non-zero (with `--force` overriding, emitting `Warning: --force overriding`), plus a negative control (strip the marker write in production → the refusal-on-threshold assertion goes red), mirroring `feedback_test.go` NEG-A. Serves the backstop. Note: `feedback-escalate` needs an embedded-schema field entry to be queryable/validated like `mod-block`.

**AC-4 — The round-state line (`round N/limit`, surface at P% of estimate, finding tiers) renders at the FO's re-dispatch decision point and in the captain gate presentation.**
Verified by: the record command emits the line (byte-observable stdout, unit-tested under AC-2); a section-scoped presence oracle over `feedback-rejection-flow` + `present-gate` asserts the render step invokes it (text claim, proven at its own level). The FO ACTING on the narration is gq's live scenario, not this AC.

**AC-5 — `feedback-rejection-flow` and `first-officer-shared-core` invoke the record command and render narration instead of prose-tracking / prose-escalation; the generic principle and dev realization land in the spots named in the Landing-spot map.**
Verified by: a section-scoped presence oracle over the skill prose (the existing `skill_text_test.go` `sectionAfter` pattern) asserting the cycle-record/escalation steps invoke `status --record-feedback-cycle` and reference the narration line; the landing-spot map is checkable by reading the named files. The behavioral half (the FO acts on a refusal) is gq's live scenario.

## Test plan

- **Unit (Go, `internal/status`):** drive the command against temp git-backed entities; assert append / count / actuals / deviation / threshold-marker / refuse / `--force` / section-scoping. Byte-observable over the resulting on-disk file + stdout — same altitude as `archive_guard_test.go`. Composes proven machinery (`section_read.go` body-parse, `mutate.go`/`merge.go arm` frontmatter stamp, `handlers.go` terminal guard); **net-new:** the `### Feedback Cycles` body-append (`atomicWrite`) and the git-diff capture (`runGitCmd`). Cost: low (no network, no live runtime).
- **Value replay (Go fixture):** the e6j cycle-boundary fixture → AC-1. The DoD's "unit test fed that fixture shape" retargets from the demoted dispatch-refusal guard to the record command's deviation computation.
- **Negative control:** strip the escalation-marker write and prove the refusal assertion goes red (mutation-proves-the-test).
- **Presence oracle (offline):** AC-4 / AC-5 skill-prose checks ride the `skill_text_test.go` `sectionAfter` pattern. Note: that helper currently lives only in a worktree copy (`.worktrees/audit-bq-release-gate/skills/integration/skill_text_test.go:277`); porting it to the main tree is part of implementation.
- **No further spike needed** beyond the one recorded: the estimate/actual semantics are exercised; the rest composes already-proven body-parse + frontmatter-mutate + terminal-guard machinery.
- **High-stakes → detached adversarial audit before merge:** this edits the `status` mutation/guard paths (a named high-stakes surface, Proof policy).

## Doc diffs (ideation proposes; implementation applies)

**`skills/feedback-rejection-flow/SKILL.md`** — steps 2-3, and a new narration step:

> - Before (step 2): `2. Track cycles in `### Feedback Cycles` in the entity body.`
> - After: `2. Record the round with `${SPACEDOCK_BIN:-spacedock} status --record-feedback-cycle {slug}` — it appends the `### Feedback Cycles` entry (actuals + estimate + deviation%) and returns the round-state line.`
> - Before (step 3): `3. On cycle 3, escalate to the human instead of another round.`
> - After: `3. Read the round-state line the command returned (`round N/limit, surface at P% of estimate`). On its refusal (non-zero exit / `feedback-escalate` marker) — deviation beyond tolerance or the banked 3rd cycle — halt at a recorded design-reset decision (reconfirm / re-scope / park / escalate), not another repair round.`

**`skills/first-officer/references/first-officer-shared-core.md:102-103`** — the `«feedback.route»` effect/done-when:

> - Before: `… track `### Feedback Cycles`, escalate on cycle 3, …` / `… (or escalated at cycle 3).`
> - After: `… record each round with `status --record-feedback-cycle` (which owns the count and stamps the escalation marker on threshold), narrate the round state, …` / `… (or escalated at the declared threshold).`

**`skills/present-gate/SKILL.md`** — add one evidence line near `Assessment:` (respecting the 15-25-line budget, rule 9; evidence-not-label, rule 11):

> After `Assessment: {N} done, {N} skipped, {N} failed.` add, when the entity has `### Feedback Cycles` entries:
> `Surface: round {N}, {P}% of the ~{estimate} declared at ideation{ — beyond tolerance if crossed}.`

**Ideation stage-def (dev README + template)** — the surface declaration:

> Add to the ideation `Outputs` list: `The task body declares an expected surface — the files/LOC (or the workflow's own unit) it expects to touch — as part of the gated design; the ideation gate approves it as the baseline the correction loop calibrates against.`

**CLI docs (`docs/site/contributing/development-workflow.md`)** — one entry documenting `status --record-feedback-cycle {slug}` (appends a cycle, stamps deviation, escalates on threshold) alongside the existing `status` surfaces.

## Boundary and notes

- **gq (`feedback-nonhappy-live-coverage`) owns the live half:** that the FO OBEYS a refusal / acts on the narration by escalating is FO-LLM behavior with no in-process Go seam — its `feedback-3-cycle-escalation` scenario, not an AC here. This guard makes the count and the deviation tamper-evident; gq proves the response.
- **Out of scope:** the budget-probe fail-safe (xa Candidate 2 — non-mechanizable, stays prose + gq's coverage); whether a `dispatch build` refusal should ever be honored (demoted — the Break-Glass routing makes it the wrong hook).
- **Leanness (0250/0260 constraint):** the count + escalation prose moves OUT of the always-on FO contract INTO the binary (boot-resident prose shrinks); the narration and declare-surface additions are lazy-loaded (`feedback-rejection-flow` loads at the rejection point; the ideation-scaffold line rides the template, not boot). Net contract-byte delta is reported at implementation.
- **Lane calibration:** conventions settled here (the `### Feedback Cycles` entry shape carrying actuals+estimate+deviation; the `feedback-escalate` marker; the round-state line format) are what ek and w0 build against after this gate; ht runs parallel.
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
