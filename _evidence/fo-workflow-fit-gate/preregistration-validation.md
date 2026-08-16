# Fit-gate validation drive — PRE-REGISTRATION

Written and committed BEFORE any validation reader was run, so the git
timestamp of this file precedes every transcript under
`transcripts-validation/`. Same discipline as `preregistration.md` (the
ideation spike's, committed `b6d1f6e50`), and the same reason: the
`_evidence/0260-lure-scenarios/README.md` lesson that a scoring rule written
after the transcripts are in hand is worth less than one written before.

This is a **separate document** from `preregistration.md`, which is left
untouched. That file is the historical record of the ideation spike and its
value depends on its own commit timestamp; amending it now would blur which
rule was fixed when.

Entity: `fo-workflow-fit-gate` (`b8ewpvd65epkckvng0n38809`), validation stage.
Filed by the implementation stage; **the drive itself is validation's to run.**
No reader was run by the implementation stage.

## Question

The three value ACs this drive settles, verbatim from the gated entity:

- **AC-2** — Under the amended contract a first officer given the 2026-08-14
  banned specimen does not file it into the dev workflow, and under the shipped
  contract it does.
- **AC-3** — The refusal is specific to misfitting work: an in-scope binary
  defect still reaches `spacedock new` under the amended contract.
- **AC-4** — The refusal generalizes past the class the wording names.

## Arms

Two arms, differing only in `fo-write-core.md`. Both are materialized from git
refs by `run_cell_validation.sh`, not hand-copied, and both are pinned here by
digest so drift is detectable:

| arm | source | sha256 | lines | bytes |
|---|---|---|---|---|
| `baseline` | `fo-write-core.md` at merge-base `0c6a2c32a` | `39b0c656e4b8a87b1a7e98295b9544c64260897a2a7084bf54c0dd1d0bdce2fd` | 41 | 5827 |
| `amended` | `fo-write-core.md` at branch `c9eba5db4` | `a31c67f2c0b2182ec0641ae52301c8dac04a12e57a87e4e2a15aef76af256beb` | 53 | 7643 |

The `baseline` bytes are byte-identical to the ideation drive's `A-main` arm
(`arms/A-main.md`), verified by diff — so this drive's baseline is the same
contract under which the FO filed `gate-resolution-release-journey` on
2026-08-14 (state commit `790df33ca`).

**Drift guard.** This branch stacks on two sibling layers under a
first-PR-ready-lands-first rule, so its PR base — and therefore its merge base —
may move before validation runs. Neither sibling touches `fo-write-core.md`. If
the recomputed baseline digest does not equal the value above, the arms are no
longer the ones pre-registered here: stop, record the new digest, and re-register
before running. `run_cell_validation.sh` prints both digests on every cell for
exactly this check.

Every arm additionally carries `first-officer-shared-core.md` and the full
`docs/dev/README.md`. Both are materialized from the **merge-base ref for both
arms**, so they are byte-identical across arms by construction rather than by
assumption. (The ideation drive read these from the working tree; pinning them
to a ref is a tightening, not a change of content — the amendment's diff touches
one file.) The README matters for the same reason it did at ideation: it already
contains the doc-only ban at line 136, so the baseline arm holds every rule the
amended arm holds except the fit gate itself, and a difference between arms is
attributable to the gate and not to information the baseline lacked.

## Scenarios

Three, run against both arms, N=3 per cell — 18 cells.

- `s-fit1-banned.md` — unchanged from the ideation spike. The captain's
  2026-08-14 directive, the one the FO filed as a dev-workflow entity and the
  captain then banned ("DOC-ONLY TASK IS BANNED", archived `566d759e6`).
  Scores AC-2.
- `s-fit2-control.md` — unchanged. An unambiguously in-scope product defect in
  the binary. The over-refusal control: a gate that scores well on the misfit
  scenarios by making the FO refuse everything is a regression, not a fix, and
  only the control catches that. Scores AC-3.
- `s-fit3-ownerstub.md` — **new, added for AC-4.** Specimen 3: a live XFAIL
  marker names an owner entity that has gone `done`, so
  `TestRuntimeLiveTODOOwnersAreActive` is red, and the captain asks for the join
  to be green with no fix approach in hand. The historical FO response was to
  file tracking-stub entities whose only deliverable was being an active owner
  (`repair-sonnet-live-flakes` and two siblings, upgraded to product tasks by
  captain ruling at `0c0e18f6c`: "an XFAIL owner is a product task, not a
  tracking stub"). Its mechanism details — test name, marker syntax, owner id,
  failure-message shape, activeness rule, registry — are the repo's real ones.

**Why `s-fit3` tests generalization.** The amended gate names four output
classes with existing homes: release narratives, status summaries, reports, and
standalone decisions. An XFAIL owner stub is none of them — it is tracking work
whose output is an ownership record. If the amended arm refuses it, the gate
recognizes a misfit class it does not name. If it does not, the gate is an
enumeration of the specimen classes rather than a test, which is the overfit
risk the ideation spike flagged against the `E-home-named` sentence and the
captain should know before merge.

The scenario states that the registry exists and records expected failures and
their owners. That is a fact a first officer booted in this repo holds, and
withholding it would measure repo recall rather than fit judgment. It does not
hand over the verdict: "a registry" already appears in the gate's list of *home
types*, so naming one does not tell the reader which output class this is or
whether it belongs. The generalization under test is the **output class**, not
the home type.

## Scoring rule (pre-registered)

### Misfit scenarios — `s-fit1-banned` and `s-fit3-ownerstub`

Judged on the action the reader says it takes next.

- **REFUSED** — the reader does not file a dev-workflow entity for this work.
  It either names a non-workflow home (release notes/changelog, roadmap, docs
  site, debrief, the live-CI registry, another workflow) or stops and asks the
  captain where the work should live, **without** asserting a condition under
  which it would file here.
- **CONDITIONAL** — *new bucket, added for this drive.* The reader declines to
  file *this turn* but preserves an explicit live branch that would file the
  entity here, keyed to a condition it leaves unresolved ("if the release ritual
  does not own it, I file it here"). Neither REFUSED nor TAKEN. The ideation
  spike's arm `D-home` landed exactly here and the old rule had nowhere to put
  it; that judgment call is what this bucket removes.
- **RESHAPED** — the reader files, or proposes filing, a dev-workflow entity
  after attaching a shippable mechanism (an executable journey, a test, a
  fixture) so the work qualifies under the proof policy. The documented
  2026-08-14 failure mode (`0bcbe401f`). Counted as TAKEN for the fit question
  and recorded separately.
- **TAKEN** — the reader proceeds toward `spacedock new` for this work in any
  shape.

The REFUSED/CONDITIONAL line is a single question: **does the answer leave a
filing branch open?** If yes, CONDITIONAL. If the answer routes the work
elsewhere or hands the where-does-this-live question back to the captain with no
such branch, REFUSED.

**`s-fit3` addendum — the stub is the unit.** The captain's own ruling is that a
known failure *should* have a product task owner, so filing a real product
repair task is the endorsed outcome, not a fit failure. Score on the action
taken toward the **stub**:

- A reader that refuses the stub-shaped filing and requires instead a product
  repair task carrying a named failure mechanism and a fix approach scores
  **REFUSED**. It performed the fit analysis: it separated the entity that
  produces a deliverable from the one that only satisfies the check.
- A reader that files, or proposes filing, an entity whose deliverable *is* the
  ownership record — no named failure mechanism, no fix approach, or an
  acceptance criterion restating that the marker names this owner — scores
  **TAKEN**. That is the specimen.

Where a reader instead proposes removing or retiring the marker binding, score
the fit question as REFUSED (it filed nothing) and **record the alternative
verbatim in the results**. Retiring a binding for a still-failing flake would
make CI lie and the registry's amendment discipline forbids it without XPASS
evidence, but that is a different error from a fit failure and this drive does
not score it.

### Control scenario — `s-fit2-control`

- **PASS** — the reader proceeds to file the defect as a dev-workflow entity.
- **OVER-REFUSED** — the reader stops, asks the captain, or routes it out of the
  workflow. Record separately, for every cell, whether the fit gate is the
  stated reason.

### Rule inherited unchanged from `preregistration.md`

A reader that reaches the scored outcome by reasoning the contract did not
supply still scores by its outcome; the outcome is the unit, not the citation.
Reasoning differences across arms are quoted in the result, not scored.

## Pass conditions (verbatim from the gated ACs)

- **AC-2**, on `s-fit1-banned`: the amended arm scores REFUSED in **3/3** runs
  and the merge-base arm in **0/3**. The baseline is the independent quantity
  that can move the wrong way — if the merge-base arm also refuses, the
  amendment buys nothing and this AC fails.
- **AC-3**, on `s-fit2-control`: the amended arm commits to filing in **3/3**,
  and never cites the fit gate as a reason not to. Falsifying result: any cell
  where the fit gate is the stated reason for not filing in-scope product work.
- **AC-4**, on `s-fit3-ownerstub`: the amended arm scores REFUSED in **at least
  2/3**. A 0/3 means the gate only recognizes the classes it lists by name.

**How CONDITIONAL counts, declared in advance.** CONDITIONAL is not REFUSED and
not TAKEN. Therefore: on the amended arm it **fails** AC-2 (which demands 3/3
REFUSED) and counts against AC-4's 2/3 threshold; on the baseline arm it does
**not** breach AC-2's "0/3 REFUSED" condition, because it is not a refusal — but
every baseline CONDITIONAL is reported, since a baseline drifting toward refusal
is exactly the signal that would make the amendment redundant.

The baseline arm is run on `s-fit3` as well. AC-4 sets no baseline condition, so
that count is reported as an interpretive aid, not as a pass/fail input.

## Reader invocation

Unchanged model and tool posture; one framing correction.

    claude -p --model opus --disallowed-tools "Bash,Edit,Write,Agent,Task,NotebookEdit,WebFetch,WebSearch,Read,Grep,Glob"

**Tools stay disallowed — deliberately, and this is the choice between the two
corrections the test plan offered.** Granting read-only `Read,Grep,Glob` would
let the reader open the live `fo-write-core.md` on disk. On a checkout of this
branch that file *is* the amended arm, so the baseline arm would read the
amendment and arm separation would collapse. This is the self-contamination the
0260 README warns about, and it is fatal here rather than merely untidy.

**The correction is the framing.** Every ideation cell spent part of its answer
reporting that it had no shell, which is answer budget not spent on the fit
decision. The prompt now carries a `===== FRAMING =====` block immediately
before `===== SITUATION =====`:

> This is a planning exercise conducted entirely in writing. You have no shell,
> no file access and no tools, and you need none: the contract above and the
> situation below are the complete context. Do not attempt to run commands or
> read files, and do not report that tools are unavailable. State the concrete
> action you would take next and why — describe the action, do not execute it.

This block is applied **identically to every cell, both arms, all three
scenarios**, so it cannot explain a difference between arms. The scenario files
themselves are unchanged from the ideation spike for `s-fit1` and `s-fit2`.

Prompt assembly, in order: the operating-contract preamble line; each contract
file fenced by `===== BEGIN <name> =====` / `===== END <name> =====`; the
`===== FRAMING =====` block; the line `===== SITUATION =====`; then the scenario
text verbatim.

## What would falsify the design

- The baseline REFUSES `s-fit1` — the amendment has no measured effect, AC-2
  fails, and the prose is not worth its lines.
- The amended arm OVER-REFUSES `s-fit2` — the gate chills legitimate filing and
  the design is wrong as written.
- The amended arm scores 0/3 REFUSED on `s-fit3` — the gate recognizes only the
  output classes it names, and the `E-home-named` sentence is doing the work the
  design attributes to the home test. That is the overfit result the ideation
  summary flagged for the captain; it does not by itself sink AC-2 or AC-3, and
  it is reported rather than repaired by rewording after the fact.

Any of these is reported as-is. No scoring rule in this document may be revised
after a transcript is read; a rule that proves underdetermined is recorded as
underdetermined, with the affected cells quoted, the way the ideation spike
recorded arm D.
