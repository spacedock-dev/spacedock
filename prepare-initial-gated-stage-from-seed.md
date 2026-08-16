---
title: Gated initial stage deadlocks — no report means not gate-ready, gated means not dispatchable
status: implementation
source: "Field report from a pre5/0.27 consumer FO, relayed by the captain 2026-08-16 with focused-test confirmation; captain order: file and dispatch"
id: ra9qzfz94hzgsq938jz998mj
gates:
    version: 1
    records:
        - id: gate:ra9qzfz94hzgsq938jz998mj:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ra9qzfz94hzgsq938jz998mj-ideation-1
              briefing:
                id: briefing:ra9qzfz94hzgsq938jz998mj:ideation:attempt-1:revision-1
                digest: sha256:2e0b2dbee158b2af2f453bb5a4d29b1885f90d259a240e26f3d5b74584acaf28
                request-digest: sha256:4e95675088ccdd457cb555d4495b314fefc87f64c1a2f512c7af416a94129f04
                room-ref: ./prepare-initial-gated-stage-from-seed/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ra9qzfz94hzgsq938jz998mj:ideation:1
                briefing: briefing:ra9qzfz94hzgsq938jz998mj:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T17:50:45.210274Z"
                decision: approve
              application:
                target-stage: implementation
                state: consumed
worktree: .worktrees/spacedock-ensign-prepare-initial-gated-stage-from-seed
---

## Problem

A workflow whose INITIAL stage is gated deadlocks. Three layers disagree, verified below against this
repo at `df0bd50d9`:

1. Contract: `skills/fo-gate-lifecycle/SKILL.md:28` requires an exact-stage report before preparation
   and mandates `report-incomplete` otherwise; `SKILL.md:24` requires `--checklist` and `--ac-scan`,
   both of which exit nonzero without a Stage Report.
2. Scheduler: `internal/gates/model.go:148/153/158` labels a gated stage with no gate record
   `validating`; `internal/status/discover.go:253-256` promotes `validating` to `needs-preparation`
   only when `hasCompleteCommittedStageReport` succeeds; `internal/status/entered_stage.go:28` defines
   that as a committed exact-stage report with checklist evidence and a Summary;
   `internal/status/format.go:150-151` excludes `validating` from `ready_gates` and `format.go:296`
   excludes every current gated stage from ordinary dispatch.
3. Result: no report → not gate-ready; at a gate → not dispatchable. The entity sits forever, labeled
   `validating` — an untrue label for a stage that has never validated anything.

The decisive intent proof: `internal/gates/prepare.go:56` does NOT require a stage report.
`validatePreparedStage` (`prepare.go:636-654`) requires only a safe path segment, a stage defined in
the README, and `Gate && !Terminal`. Its own fixture (`internal/gates/prepare_test.go:1002`) declares
an `initial: true, gate: true` stage and an entity body of `# Task` with no report. The low-level gate
machinery was designed for this; the scheduler and contract layers above it contradict the machinery's
own test.

The only workaround is a fabricated "backlog completion report" restating the committed seed — fake
work of exactly the class `fo-workflow-fit-gate` (b8) names, except here the contract demands the
ceremony rather than the FO inventing it.

## Citation audit (2026-08-16, against `df0bd50d9`)

Every seed citation was checked at line level. **No behavioral claim failed.** Two line numbers drifted
(the seed's came from a consumer's 0.27.0 plugin cache); the behavior they name is real and located
above. The design direction stands unchanged.

- `SKILL.md:24`, `SKILL.md:28`, `entered_stage.go:28`, `format.go:150`, `format.go:296`,
  `prepare.go:56`, `prepare_test.go:1002` — exact, both line and behavior.
- `model.go:155` — **line wrong, behavior right.** Line 155 is `return "invalid"` for a malformed
  record. The `validating` labeling is at 148 (`doc == nil`), 153 (`no logical gate`), and 158 (zero
  attempts).
- `discover.go:261` — **line off by seven, behavior right.** Line 261 is the assignment
  `entity.fields["gate-readiness"] = readiness`; the report-conditioned promotion guard is at 253-256.

Three findings from the audit that **do** change the design, all recorded in the approach below:
the deadlock is one missing edge rather than four; `gates.ReadinessStage` cannot carry the exception;
and the report demand lives in two shipped contract files, not one.

## Spike record

The riskiest unverified mechanisms were exercised against the real binary before designing. Fixture
scripts: `/Users/clkao/.claude/jobs/4e49247e/tmp/deadlock-fixture.sh` and `guard-fixture.sh`.

1. **Deadlock reproduced (red baseline).** Workflow with `backlog` as `initial: true, gate: true`;
   seed committed and clean; no stage report. Observed: `status --next --json` →
   `{"dispatchable":[],"ready_gates":[]}`; `gate-readiness` → `validating`; `--checklist` and
   `--ac-scan` both exit nonzero with `no ## Stage Report for stage "backlog"`. The seed carried a
   real `## Acceptance criteria` section with two ACs and `--ac-scan` still refused — the refusal is
   report-gated, not AC-gated.
2. **The machinery accepts the seed (live, not just the unit fixture).** `gate prepare` on that same
   no-report initial entity exited 0 and emitted `room`, `briefing`,
   `digest=sha256:0aa234a9…`, `state=open` — exactly the four values `SKILL.md:29` requires.
3. **Everything after prepare already works.** Post-prepare the entity reported `awaiting-captain` and
   appeared in `ready_gates` under both `status --next --json` and `status --boot --identify --json`.
   **The deadlock is a single missing edge — promotion to `needs-preparation` — not four broken layers.**
4. **Guard baseline.** A NON-initial gated stage with no report → `validating`, absent from
   `ready_gates`. The same stage with a complete committed report → `needs-preparation`, present.
   This is the behavior the fix must leave untouched.
5. **The keystone mechanism is real.** `status --boot --identify --json` emits per-stage
   `initial: "true"|"false"`, parsed into the same `[]Stage` slice `materializeGateReadiness` receives
   (`Stage.initial`, `stages.go:26`). The exception can key on `initial` at the promotion site with no
   new parsing.

## Proposed approach

Because spike 3 proved the gate lifecycle is already whole after `prepare`, the fix is one scheduler
edge plus contract prose. Two of the seed's five items dissolve; the reasoning is recorded rather than
applied silently.

**1. Scheduler promotion (the only necessary code change).** In `internal/status/entered_stage.go`,
add one predicate; in `discover.go:254`, call it instead of `hasCompleteCommittedStageReport`:

```go
// gatePreparable reports whether a gated stage's promotion proof is satisfied.
// A non-initial gated stage owes a complete committed stage report. An INITIAL
// gated stage had no prior stage to write one: the committed clean seed IS the
// artifact the captain reviews, so durability alone is the proof. The exception
// keys on stage.initial ONLY — never on "no report exists" — so no later stage
// gains a path around its completion proof.
func gatePreparable(path string, stage Stage) bool {
	if stage.initial {
		return entityPathCleanInHEAD(path)
	}
	return hasCompleteCommittedStageReport(path, stage.Name)
}
```

`entityPathCleanInHEAD` (`entered_stage.go:85`) already exists and already IS the clean-in-HEAD guard
the seed requires preserved; it is reused verbatim rather than reimplemented. The change is confined to
`discover.go:254` — the only other caller of `hasCompleteCommittedStageReport` is
`enteredStageAwaitingCompletion`, which already returns `false` for any initial, gate, or terminal
stage (`entered_stage.go:16`) and is therefore untouched.

**2. The clean-in-HEAD guard and digest pinning survive** — by reuse, not by re-assertion. Spike 2
observed `prepare` emitting a digest over the committed seed; the fix touches nothing on that path.

**3. No new label — the seed's item 4 dissolves into item 1.** Once promotion fires, `validating`
never appears for a committed clean initial seed; it reports `needs-preparation`, which is true. The
set of possible `gate-readiness` values is unchanged. Introducing a sixth vocabulary word would add a
term every FO reader must learn for a state that item 1 already eliminates.

The one residual: an initial gated seed that is *dirty or uncommitted* still reads `validating`. Left
alone deliberately. It is transient and self-correcting (commit the seed), and the identical label
already covers the non-initial no-report case, which is equally not-yet-validated. Fixing a label for a
one-commit-away transient state is not worth a vocabulary expansion.

**4. No seed-AC-scan mechanism — the seed's item 3 is rejected on necessity.** The seed proposed
replacing report-based checklist/AC extraction with a scan of the seed's own AC section. Three reasons
not to build it:

- At a seed gate every AC is `unevidenced: true` by construction — there is no report whose checklist
  ranges could cite one. Emitting a constant is ceremony, not evidence.
- `present-gate` already handles absent evidence: it instructs the FO to "omit missing evidence, empty
  result classes… and placeholders such as `None` or `N/A`", and requires only entity+stage, Briefing
  identity and digest, one recommendation, and one concrete decision effect. None of those need a
  checklist read. `present-gate` needs **no change**.
- Many legitimate seeds carry no `## Acceptance criteria` at all — including this entity's own seed,
  which had only Problem / Proposed approach / Out of scope. Requiring one so an `--ac-scan` can
  succeed would invent precisely the ceremony b8 exists to stop.

Simplest alternative considered: a `--ac-scan` fallback that skips the `selectStageReport` error for an
initial stage and emits ACs with empty citations (~5 lines in `gate_extract.go`). Rejected: it buys the
FO a structured list it can already read from the digest-pinned artifact, and it still fails loudly on
the AC-less seeds that are the common case. Instead the SKILL exception omits step 3 entirely — zero
binary surface.

**5. Contract prose in TWO shipped files, not one.** The report demand is stated twice. Amending only
`fo-gate-lifecycle/SKILL.md` leaves the FO reading `fo-dispatch-core.md:193` and deadlocking anyway.

`skills/fo-gate-lifecycle/SKILL.md` — insert before line 28, and scope line 28's heading:

> **Initial-stage gate.** At a stage the workflow marks `initial: true`, the committed seed IS the
> reviewed artifact: no prior stage existed to write a report. Prepare directly from the seed with
> `--artifact` naming the entity; never author, request, or wait for a stage report, and never emit
> `report-incomplete`. Skip step 3 — both structured reads exit nonzero with no report. Present
> outcome, scope, exclusions, and the ideation proof the next stage owes.

Line 28 before: `**Cold report candidate.** Structurally review…`
Line 28 after: `**Cold report candidate (non-initial stages).** Structurally review…`

Line 15 before: `` `needs-preparation`: review report; ``
Line 15 after: `` `needs-preparation`: review report, or at an `initial: true` stage the committed seed itself; ``

`skills/first-officer/references/fo-dispatch-core.md:193` before:

> For `needs-preparation`, load `spacedock:fo-gate-lifecycle`, re-read the entity, latest exact-stage
> report/checklist, and path-scoped clean commit, then decide semantic completeness.

after:

> For `needs-preparation`, load `spacedock:fo-gate-lifecycle`, re-read the entity, latest exact-stage
> report/checklist, and path-scoped clean commit, then decide semantic completeness. At an
> `initial: true` stage there is no prior report: the committed clean seed is the reviewed artifact,
> the structured checklist/AC reads are skipped, and `report-incomplete` never applies.

**6. Documentation diff.** `docs/site/concepts/gates-and-decisions.md:73-79` describes the transition
this change alters. Before:

> After completion verification, a gate with no current-stage authority remains `validating` until the
> mechanical report checks pass. It then appears as `needs-preparation` on boot and every machine
> scheduler read.

After:

> After completion verification, a gate with no current-stage authority remains `validating` until the
> mechanical report checks pass. It then appears as `needs-preparation` on boot and every machine
> scheduler read. A gated stage the workflow marks `initial: true` has no prior stage to have written a
> report: there the committed, clean-in-HEAD seed is itself the proof, and the entity appears as
> `needs-preparation` directly. Engage prepares from the seed and skips the structured checklist/AC read.

`docs/site/reference/command-reference.md:106-112` needs no change — it describes the `ready_gates`
envelope shape, which is unchanged.

## Expected surface and tolerance

Seven files, net roughly +120 lines, of which about 80 are tests. Tolerance ±40 lines and ±1 file.

| File | Change |
| --- | --- |
| `internal/status/entered_stage.go` | +10, new `gatePreparable` |
| `internal/status/discover.go` | 1 line, call site at 254 |
| `internal/gates/model.go` | +3, rename the `reportComplete` parameter to `promotionProven` and update its doc comment; no behavior change |
| `skills/fo-gate-lifecycle/SKILL.md` | +6 / 2 edited |
| `skills/first-officer/references/fo-dispatch-core.md` | +3 |
| `docs/site/concepts/gates-and-decisions.md` | +4 |
| `internal/status/*_test.go` | +80 |

The `model.go` rename is honesty at the seam: `discover.go` will pass `true` for a case where no report
exists, and a parameter named `reportComplete` would then be a lie at the call boundary. It is a
positional argument, so no caller breaks.

This is wider than the dispatch's hinted "binary + one skill file", for two audited reasons: the report
demand is duplicated in `fo-dispatch-core.md`, and a shipped doc describes the transition. Conversely
the hinted `internal/gates` *label* is not needed — `gates.ReadinessStage` (`model.go:130-134`) carries
only Name/Gate/Terminal, so keying the exception there would force widening that struct and every
caller, while `internal/status` already holds `Stage.initial` for free.

## Semantics changed

- **Scheduler behavior.** A committed, clean entity at a gated `initial: true` stage now reports
  `gate-readiness: needs-preparation` and appears in `ready_gates` for `status --next --json` and
  `status --boot --identify --json`. Previously `validating`, absent from both.
- **Status vocabulary — explicitly UNCHANGED.** No `gate-readiness` value is added, removed, or
  retired.
- **Contract authority.** `report-incomplete` becomes inapplicable at an initial stage; the FO gains a
  documented path to prepare without a report. No other stage's completion proof weakens.
- **Command grammar, stored formats, `Prepare` semantics, digest pinning, and the dispatch exclusion
  at `format.go:296` — all unchanged.** A gated stage remains permanently excluded from ordinary
  dispatch; this task does not touch that.

## Acceptance criteria

**AC-1 (VALUE)** — On the deadlock fixture (gated initial stage, committed clean seed, zero
`## Stage Report` sections in the entity at every step), the entity reaches a captain-presentable gate:
`status --next --json` lists it in `ready_gates` with `readiness: needs-preparation`, `gate prepare`
then yields `awaiting-captain`, and the entity file's Stage Report count is still zero afterward.
Measured against the recorded red baseline `{"dispatchable":[],"ready_gates":[]}` from spike 1.
Verified by: a Go test in `internal/status` driving that fixture. Falsifying change: any fix that
requires a report makes the Stage Report count nonzero, and reverting `gatePreparable` returns
`ready_gates` to empty.

**AC-2 (GUARD)** — A NON-initial gated stage with no complete stage report still refuses promotion:
it reports `validating` and is absent from `ready_gates`, while the same stage with a complete
committed report reports `needs-preparation` and is present. Verified by: a two-entity Go test
mirroring the spike-4 guard fixture. Falsifying change: keying the exception on "no report exists"
instead of `stage.initial` promotes the no-report non-initial entity and fails this.

**AC-3 (GUARD)** — An entity at a gated initial stage whose seed is untracked or dirty against HEAD is
NOT promoted; it reports `validating`. Verified by: a Go test that writes the seed, asserts absence
before commit, commits, and asserts presence after. Falsifying change: dropping the
`entityPathCleanInHEAD` call promotes the dirty seed.

**AC-4** — The initial-stage gate briefing pins the exact committed seed bytes: preparing an unchanged
seed twice replays one identical digest, and a seed modified between prepares yields a different one.
Verified by: a Go test in `internal/gates` asserting digest equality then inequality. Falsifying
change: unpinning the artifact makes the two digests match across a modified seed.

**AC-5 (VALUE, prose)** — A live First Officer drive over the deadlock fixture reaches a presented gate
without any stage report being written and without a `report-incomplete` stop: the driven run's
resulting on-disk entity contains zero `## Stage Report` sections and one `state=closed`-eligible open
attempt. Verified by: a one-off live drive recorded in the validation report. Falsifying change: a
drive against the unamended contract stops at `report-incomplete` or fabricates a report.

**AC-6** — The suite stays green: `go test ./...` plain and `-race`.

## Test plan

| AC | Vehicle | Cost |
| --- | --- | --- |
| AC-1, AC-2, AC-3 | Go tests in `internal/status`, building README + entity fixtures in a temp git repo and asserting `gate-readiness` and `ready_gates`. The existing `internal/status` fixture helpers already build exactly this shape. | Low. ~80 lines, all three share one builder. |
| AC-4 | Go test in `internal/gates` reusing `prepareFixture`, which already declares an `initial: true, gate: true` stage with no report. | Low. ~20 lines, no new fixture. |
| AC-5 | One-off live FO drive at validation, recorded in the report. | Medium — one live run. |
| AC-6 | `go test ./...` and `-race`. | Existing lanes. |

**No committed prose-grep.** AC-5 is deliberately a one-off exercise, not a standing test. A committed
check that the SKILL text contains the exception would assert only that the file contains what the
implementer wrote (dev README line 81; the same ruling b8 records for its own prose amendment).

**Required live lanes.** The diff touches `skills/first-officer/references/fo-dispatch-core.md`, the
host-neutral dispatch core, so by dev README line 84 **every** host lane (claude, codex, pi) is
required green before merge, not just one. This is the dominant cost of the task and is a direct
consequence of finding 5 above.

**A standing live journey is NOT proposed.** Adding a common journey for the initial-seed gate would
cost roughly one new journey across three lanes on every run (the registry's measured baseline is
127s Sonnet / 144s Opus per journey) to guard model obedience to prose. Recommended instead: the
one-off drive. If the captain wants standing regression protection for the FO behavior specifically,
that is the price, and it is the gate's call rather than ideation's.

**Implementation base.** Implementation branches from the current stack tip at its dispatch time; the
base is deliberately not pinned here.

## Cross-reference: `fo-workflow-fit-gate` (b8, in ideation)

Complementary, not overlapping — b8's own line 64 already records this boundary from its side
("Reconcile framing, do not absorb that fix"). Both target the same failure class: ceremony that
produces no value. b8 addresses the FO *inventing* ceremony and remedies it with an admissibility gate
in `fo-write-core.md`; this entity addresses the *contract demanding* ceremony and remedies it with a
scheduler edge plus gate-lifecycle prose. The shipped files are disjoint (b8 owns `fo-write-core.md`;
this owns `fo-gate-lifecycle/SKILL.md`, `fo-dispatch-core.md`, and `internal/status`), and neither
blocks the other in either order. Both adopt the same proof discipline: prose amendments proven by a
one-off exercise, never a standing prose check.

## Out of scope

- The workflow-fit-gate write-core amendment (b8 — complementary entity, cross-referenced above).
- The consumer workflow's own state (8e re-prepare is their field validation, referenced only).
- A new `gate-readiness` vocabulary value, including for the dirty-initial-seed transient.
- A seed-AC-scan mechanism in the binary (rejected on necessity; reasoning recorded above).
- A standing live journey for the initial-seed gate (costed above; the gate may overrule).

## Stage Report: ideation

- DONE: The seed's citations come from a consumer's 0.27.0 plugin cache — verify every one against THIS repo's source at line level first (model.go labeling, discover.go promotion, entered_stage predicate, format.go ready_gates + dispatch exclusion, SKILL.md:24/28, prepare.go:56, prepare_test.go:1002), and run the focused tests yourself. Any citation that does not hold changes the design; say so rather than adapting silently.
  All nine checked at `df0bd50d9`; zero behavioral misses, two line drifts stated openly in Citation audit (`model.go:155` is `return "invalid"`, not the `validating` labeling at 148/153/158; `discover.go:261` is the assignment, guard at 253-256). `go test ./internal/gates/ ./internal/status/` green (22.0s / 43.8s).
- DONE: Reproduce the deadlock in a fixture before designing: a workflow whose initial stage is `gate: true` with a committed clean seed — show the scheduler emits no needs-preparation row and ordinary dispatch excludes it. This fixture becomes the red-before/green-after proof.
  Spike 1: `status --next --json` → `{"dispatchable":[],"ready_gates":[]}` with `gate-readiness: validating`. Script `deadlock-fixture.sh`; the red baseline AC-1 measures against.
- DONE: Design the fix at each layer, exception keyed on `initial: true` only: scheduler promotion for committed clean seeds, the honest state label for gated-no-record initial stages, the seed-AC-scan replacement for report-based checklist/AC extraction, and the SKILL.md initial-stage exception text (specific before/after wording, skill-change rule). The clean-in-HEAD guard and briefing digest pinning must survive unchanged.
  Designed, but two sub-items resolved as NOT-BUILD on necessity and the reasoning is on the record, not applied silently: the honest label falls out of promotion (no new vocabulary), and the seed-AC-scan is rejected because every seed AC is `unevidenced` by construction, `present-gate` already omits absent evidence, and AC-less seeds are the common case. Before/after wording given for both contract files plus the doc diff; guard and digest pinning survive by reuse of `entityPathCleanInHEAD` and an untouched prepare path (spike 2 observed prepare emitting a digest over the committed seed — the mechanism AC-4 binds).
- DONE: The guard-preservation direction is mandatory: a NON-initial gated stage without a complete report still refuses promotion — proven by test, not asserted.
  Spike 4 measured the live baseline both ways (no-report → `validating`, absent; with-report → `needs-preparation`, present) and AC-2 binds the committed test that reproduces it. The test itself lands at implementation; ideation owns the baseline and the falsifying change.
- DONE: Value AC measures the outcome: seed → prepared gate → captain-presentable with zero fabricated reports, on the deadlock fixture; paired with the guard-preservation AC.
  AC-1 measures `ready_gates` promotion plus a zero Stage Report count against the recorded red baseline; AC-5 measures the same outcome under a live FO drive. Both paired with AC-2/AC-3 guards.
- DONE: Cross-reference fo-workflow-fit-gate (b8, in ideation now): complementary — there the FO invents ceremony, here the contract demands it. Reconcile framing; do not absorb its scope.
  Reconciled in its own section; b8's line 64 already records the boundary from its side. Shipped files are disjoint and neither blocks the other.
- DONE: Declare expected surface (binary internal/status + internal/gates label + one skill file) and semantics changed (scheduler behavior, status vocabulary, contract prose). Test plan per AC with costs.
  Surface declared as seven files / ~+120 lines / ±40 tolerance, and it departs from the hint in both directions: the `internal/gates` label is NOT needed (`ReadinessStage` carries no `Initial`), while `fo-dispatch-core.md` and a docs-site page ARE, which makes every host live lane required by dev README line 84. Semantics table states status vocabulary explicitly UNCHANGED.
- DONE: Implementation will branch from the current stack tip at its dispatch time; note it, do not pick the base now. Write the ideation stage report and stop for the gate.
  Noted under Test plan; no base pinned. Stopping for the gate.

### Summary

Verified all nine seed citations at line level against `df0bd50d9`: every behavioral claim holds, two line numbers drifted from the consumer's 0.27.0 cache, so the design direction stands. Five live spikes against the real binary reproduced the deadlock, then found the decisive simplification — `gate prepare` already succeeds on a no-report initial seed and the entity immediately becomes `awaiting-captain` in `ready_gates`, so this is one missing scheduler edge, not four broken layers. That collapses the seed's five-item approach to a ~10-line predicate keyed on `stage.initial` plus contract prose, and I rejected two proposed mechanisms on necessity grounds rather than adapting silently: no new status vocabulary, and no seed-AC-scan in the binary. Two findings widen the surface beyond the dispatch's hint — the report demand is duplicated in `fo-dispatch-core.md`, which makes all three host live lanes required — and one narrows it, since `gates.ReadinessStage` carries no `Initial` field and is best left alone.
