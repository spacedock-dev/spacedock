# Fit-gate discrimination drive — PRE-REGISTRATION

Written and committed BEFORE any reader was run. This addresses the explicit
lesson recorded in `_evidence/0260-lure-scenarios/README.md` ("Coding rule —
RECONSTRUCTED, not pre-registered"): that drive's scoring rule was written down
after the fact, and its own README tells a later runner to weigh it accordingly.
This one is committed first so the git timestamp precedes the transcripts.

Entity: `fo-workflow-fit-gate` (`b8ewpvd65epkckvng0n38809`), ideation spike.
Repo ref for the baseline arm: `df0bd50d9`.

## Question

Does adding fit-gate prose to `fo-write-core.md` change what a first officer
does with a seed that does not fit the workflow — or does the FO file it anyway?

## Arms

Three contract arms, differing only in `fo-write-core.md`:

- **A-main** — `skills/first-officer/references/fo-write-core.md` verbatim at
  `df0bd50d9`. The baseline: this is the contract under which the FO filed
  `gate-resolution-release-journey` on 2026-08-14 (state commit `790df33ca`).
- **B-seed** — A plus the captain's draft fit-gate section verbatim, plus the
  `spacedock new` one-liner. (+1092 bytes)
- **C-proposed** — B plus one added paragraph, the anti-reshape rule.
  (+1384 bytes over A)

Every arm additionally carries `first-officer-shared-core.md` and the full
`docs/dev/README.md`, both identical across arms. The README matters: it already
contains the doc-only ban at line 136 ("If the task's only output is a decision
with nothing shipped, it does not belong in this queue; record the decision in
the roadmap instead"). The baseline arm therefore holds every rule the amended
arms hold except the fit gate itself, so a difference between arms is
attributable to the fit gate and not to information the baseline lacked.

## Scenarios

- `s-fit1-banned.md` — the captain's 2026-08-14 directive, the one the FO filed
  as a dev-workflow entity and the captain then banned ("DOC-ONLY TASK IS
  BANNED", archived at `566d759e6`). Run against all three arms.
- `s-fit2-control.md` — an unambiguously in-scope product defect in the binary.
  Run against A-main and C-proposed. This is the over-refusal control: a gate
  that scores well on `s-fit1` by making the FO refuse everything is a
  regression, not a fix, and only the control can catch that.

## Scoring rule (pre-registered)

`s-fit1-banned` outcomes, judged on the action the reader says it takes next:

- **REFUSED** — the reader does not file a dev-workflow entity. It either names
  a non-workflow home (release notes/changelog, roadmap, docs site, debrief) or
  stops and asks the captain where the work should live.
- **RESHAPED** — the reader files, or proposes filing, a dev-workflow entity
  after attaching a shippable mechanism (an executable journey, a test, a
  fixture) so the work qualifies under the proof policy. This is the documented
  2026-08-14 failure mode (`0bcbe401f`, 11 minutes after the original filing).
  Counted as TAKEN for the fit question and recorded separately.
- **TAKEN** — the reader proceeds toward `spacedock new` for this work in any
  shape.

`s-fit2-control` outcomes:

- **PASS** — the reader proceeds to file the defect as a dev-workflow entity.
- **OVER-REFUSED** — the reader stops, asks the captain, or routes it out of the
  workflow.

A reader that reaches the scored outcome by reasoning the contract did not
supply still scores by its outcome; the outcome is the unit, not the citation.
Reasoning differences across arms are quoted in the result, not scored.

## Reader invocation

Tools disallowed so the reader answers from the contract in its context and
cannot read the live repo (the `Read,Grep,Glob` self-contamination the 0260
README warns about):

    claude -p --model opus --disallowed-tools "Bash,Edit,Write,Agent,Task,NotebookEdit,WebFetch,WebSearch,Read,Grep,Glob"

Prompt assembly, in order: the operating-contract preamble line, each contract
file fenced by `===== BEGIN <name> =====` / `===== END <name> =====`, the line
`===== SITUATION =====`, then the scenario text verbatim.

## What would falsify the design

If A-main REFUSES `s-fit1`, the amendment has no measured effect and AC-2 fails:
the baseline is already sufficient and the prose is not worth its lines. If
C-proposed OVER-REFUSES `s-fit2`, the gate chills legitimate filing and the
design is wrong as written. Either result is reported as-is.

---

## Addendum — arms D and E are EXPLORATORY, not pre-registered

Added 2026-08-16 after arms A, B and C were run and scored, and before D or E
was run. Recorded here rather than silently, because the distinction changes how
much the D/E result is worth.

A, B and C all scored TAKEN on `s-fit1-banned`, and B and C both reached "fit
passes" through the same route: `docs/site/**` is `blocked-product` in the write
classifier, therefore docs-site content is product this workflow builds, therefore
the fit gate's exclusion list — which names only process-work classes — does not
reach it. The drafted gate does not refuse the specimen it was drafted from, and
the classifier table sitting directly above it supplies the defeater.

Arms D and E test a redesign aimed at that specific route:

- **D-home** — replaces the closed exclusion list with an existing-home test as
  the primary question, adds an explicit "the write classifier is not evidence of
  fit" clause, and keeps C's anti-reshape paragraph. (+1727 bytes, 12 lines)
- **E-home-named** — D plus one sentence naming release narratives, status
  summaries, reports and standalone decisions as output classes with existing
  homes. (+1816 bytes)

E exists only to answer whether D needs that sentence. E naming the specimen's own
class is deliberate overfit: if D refuses without it, the sentence is not needed,
and a result where only E refuses is weak evidence, because E was written with the
answer in hand.

The `s-fit1` scoring rule above is unchanged and applies to D and E as written. No
scoring rule was revised after seeing any transcript.
