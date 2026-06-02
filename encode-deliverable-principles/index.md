---
id: v1awnfhs996ykymv409anywh
title: Encode deliverable principles + template ergonomics into README, FO/ensign contract, and status guards
status: implementation
source: FO triage (2026-06-01) — docs/dev/_proposals cleanup; consolidates the deliverable-principles study + the TDD/template-adoption ergonomics
started: 2026-06-02T04:57:30Z
completed:
verdict:
score: "0.30"
worktree: .worktrees/spacedock-ensign-encode-deliverable-principles
issue:
---

Encode the deliverable principles + template-ergonomics snippets into the workflow's own
contract surfaces so future dev work cannot drift into the antipatterns they forbid, AND
formalize the detached-adversarial-audit discipline that has repeatedly caught what
validation passed. Seeded from the parked design study in `proposal.md` (formerly
`docs/dev/_proposals/encoding-deliverable-principles.md`) — which carries the original
before/after wording, the dogfood/falsifiability analysis, and the placement map.

The four principles, as the proposal stated them: (1) **no doc-only deliverable** — every AC's
oracle is EXTERNAL to the entity body; (2) **proof is behavioral, not grep** — exercise the
behavior and observe the outcome (substring present/absent is not proof; invariant-over-real-values
is); (3) **enforce in code, not prose** — a prose-only contract has ceiling "wording present" and
can't alone satisfy an AC; (4) **spike the riskiest unknown in ideation** — smallest end-to-end
exercise of the riskiest path first, recording behavioral evidence or `no spike needed: {mechanisms}`.

## Problem

The proposal predates `se` (#248, archived PASSED). `se` ALREADY shipped almost everything the
proposal proposed. A diff of the live contract surfaces (this ideation's spike — see "No spike
needed", below) shows the proposal's placement map is now mostly redundant. Re-filing it verbatim
would re-write text that is already present, which the workflow forbids (do-not-double-file). So
the real problem this entity closes is narrow and twofold:

1. **One P1 hardening is still missing.** `se` shipped the README's P1 prose, the ensign-contract
   habits, P3 code-gate-preference, the FO posture, and the P4 spike rule. It did NOT harden the FO
   operating contract's **AC coverage cross-check** (`## Completion and Gates`): that clause still
   reads "confirm each AC has at least one evidence citation" — it does not yet say the evidence must
   come from an oracle EXTERNAL to the entity body, that a self-oracle AC can never satisfy the
   cross-check, or that a pure-decision entity must not advance to terminal PASSED. The proposal's
   P1 "after" wording for that clause is the surviving gap.

2. **The detached-audit discipline is nowhere in the shipped contract.** The captain folded the
   detached-audit formalization into this entity (debrief 2026-06-02 #1: "Audit-formalization folds
   into `v1`"). Today the read-only adversarial audit on a detached checkout is an undocumented FO
   habit — it earned its keep again on #262 (M1/M2, below) and on 5p/jg/37 — but no contract surface
   names when it triggers, what it produces, or how it is recorded. It is the live proof that
   validation-PASSED is not the last word for high-stakes surfaces, yet a clean-room FO has no
   instruction to run it.

## Proposed approach

### A. Reconcile against shipped state — land only the gap (do not double-file)

| Piece | State | This entity's action |
|---|---|---|
| README P1 (ideation Outputs + done Bad) | shipped by `se` (README:80, :124) | none |
| README P2 ("choose proof" bullet + validation Bad) | shipped by `se` (README:84, :114) | none |
| README P3 (ideation Outputs bullet) | shipped by `se` (README:85) | none |
| README P4 (ideation Outputs + Staff review) | shipped by `se` (README:81, :90) | none |
| README template ergonomics (`## Out of scope`, `## Problem`/`## Proposed approach`/`## Test plan` headings, `--next` doc) | shipped by `se` (README:143, :164-185) | none |
| Ensign contract: TDD test-first, no-hidden-deps, real-checkable-change, prove-by-exercising | shipped by `se` (ensign-shared-core:22-29) | none |
| FO contract P3 code-gate-preference + FO posture | shipped by `se` (FO core:379-389) | none |
| FO contract P4 spike rule (`## Probe and Ideation Discipline`) | shipped by `se` (FO core:393) | none |
| **FO contract P1 AC-cross-check hardening** | **NOT shipped** (FO core:134 is the un-hardened original) | **ship — this entity** |
| **Detached-audit-discipline formalization** | **NOT shipped** (absent everywhere) | **ship — this entity** |
| P1's `status --validate` self-oracle lint + terminal-PASSED `--set` guard | **filed as `2a` (require-external-proof-guard)** | **do not re-file; coordinate** |

The two surviving edits both land in the VENDORED `skills/first-officer/references/first-officer-shared-core.md`
(confirmed leading edge: canonical `~/git/spacedock` lacks `## Working Principles` AND the detached-audit
section — verified by grep this ideation, see spike note). Author vendored-first, flow upstream during
self-hosting, exactly as `se`'s edits did.

### B. Edit 1 — harden the FO AC coverage cross-check (P1 wording; the code teeth are 2a)

`first-officer-shared-core.md` `## Completion and Gates`, the `**AC coverage cross-check.**` paragraph.

BEFORE (FO core:134):

```
**AC coverage cross-check.** Additionally, at every gate, scan the entity body's `## Acceptance criteria` section and confirm each `**AC-N**` item has at least one evidence citation from this stage's report or a prior stage report. Name any AC without evidence; REJECT if this stage was the natural place to address it. This cross-check is independent of checklist DONE/SKIPPED/FAILED accounting — checklist items are dispatch signals, AC items are entity properties.
```

AFTER:

```
**AC coverage cross-check.** Additionally, at every gate, scan the entity body's `## Acceptance criteria` section and confirm each `**AC-N**` item has at least one evidence citation from this stage's report or a prior stage report. The evidence must come from a check OUTSIDE the entity body — a test, a command's output or exit code, or resulting on-disk state. An AC whose only cited proof is review of the entity's own prose ("verified by review of this entity's decision section") proves only that the prose exists; it can never fail, so it does not satisfy this cross-check. Name any AC without external evidence; REJECT if this stage was the natural place to address it. When an entity's only deliverable is a decision with nothing shipped, do not advance it to a terminal PASSED verdict — surface to the captain that the decision belongs in the roadmap. (Where a behavioral guarantee can be enforced by the binary or a failing test, the real assurance is that gate, per `## Working Principles`; the `spacedock status --validate` self-oracle lint and the terminal-PASSED set guard are that backstop when present.) This cross-check is independent of checklist DONE/SKIPPED/FAILED accounting — checklist items are dispatch signals, AC items are entity properties.
```

This is **contract prose** — its honest ceiling is wording-present, which the workflow's own rule
caps at "cannot stand alone." The BEHAVIORAL guarantee (a self-oracle AC mechanically cannot reach
terminal PASSED) is `2a`'s code guard, NOT this entity's. The parenthetical points at 2a's gate
rather than re-asserting it; this edit ships the FO-facing wording 2a's guard backs. The proof for
THIS edit is a presence/banned-token test over the real file (legitimate doc-as-deliverable — see AC-1
proof note), in the established `skills/integration/working_principles_test.go` / `contract_gate_test.go`
style, written failing-first.

### C. Edit 2 — formalize the detached-adversarial-audit discipline (new contract section)

A new `## Detached Adversarial Audit` subsection in `first-officer-shared-core.md`, placed after
`## Gate Presentation` (the audit feeds the gate, so it sits adjacent to gate handling). It must be
concrete on three axes the captain named: WHEN it triggers, WHAT it produces, HOW it is recorded.

Proposed wording (implementation may tune prose; the load-bearing invariants are the trigger-surface
list, the read-only/detached-checkout requirement, and the recorded-finding-routing rule):

```
## Detached Adversarial Audit

For HIGH-STAKES surfaces — the front-door launcher (`spacedock claude`/`codex`/`doctor`), the `status`
mutation/guard paths, the shipped contract/scaffolding, and the CI/release machinery — a validation
PASSED is not sufficient on its own. Before merging such a change, the FO runs (or dispatches) a
read-only adversarial audit on a DETACHED checkout of the merge result:

- **Detached and read-only.** The auditor works on a separate checkout (a throwaway worktree of the
  merge candidate), never the implementation worktree, and never mutates the deliverable. It tries to
  REFUTE the validation: construct an adversarial edit that the deliverable's own tests should catch,
  and confirm they do. A test that stays green under an edit that breaks the claim is a hole.
- **What it produces.** A short finding list in two tiers — `Material:` (a real correctness or
  test-strength hole, e.g. an assertion that green-lights a regression) and `Polish:` (non-blocking).
  "Refuted nothing material" is itself a valid, recorded outcome — the audit confirms correctness
  independently, which is the point on a front-door change.
- **How it is recorded and routed.** Material findings route through the normal Feedback Rejection
  Flow into the prior stage (a `### Feedback Cycles` entry naming the audit and its adversarial edit);
  the gate is not presented as clean until they are closed. A clean audit is noted in the gate
  presentation's `Reviewer findings` block (or a one-line "detached audit: no material findings").

This is why a validation report is necessary but not final for these surfaces: the audit catches the
class of hole where the TEST passes but would also pass on a broken future edit — which validation,
trusting its own green suite, does not. It is read-only refutation, not a second implementation pass.
```

This is also **contract prose**; its ceiling is wording-present and it is marked as such. There is no
code gate that can decide "did a sound adversarial audit happen" — that is FO judgment on a
judgment-bound surface. The enforceable floor is "the discipline is written down where a clean-room FO
reads it, with concrete triggers and a recording path," which is the legitimate doc-as-deliverable case
(an agent-read contract where the prose IS the executable instruction). Proof is a presence test over
the real file asserting the section names its trigger surfaces, the detached/read-only requirement, and
the Feedback-Cycles routing — same proof tier as `se`'s `working_principles_test.go`.

### Grounding the audit formalization in real catches (AC-3)

The formalization is not abstract — it codifies catches already on the record:

- **#262 (`binary-absent-fo-bootstrap`, archived PASSED) M1 + M2.** Validation recommended PASSED on
  correct shipped prose. The detached audit then found TWO test-strength holes in
  `contract_gate_test.go`: **M1** — the Class-A no-doctor check was `if n := strings.Count(...); n > 0`,
  so deleting the prohibition entirely (zero mentions) SKIPPED the check and passed; **M2** — the
  Class-B doctor assertion was a bare `strings.Contains`, satisfied by a negated disclaimer mention. Both
  were verified by the auditor constructing the adversarial edit (delete the prohibition; replace the
  route with "Historically we suggested... but no longer") and watching the test stay green. Routed
  through Feedback Cycles, fixed, re-validated. This is the canonical "validation passed, audit refuted"
  catch — and it is a test-STRENGTH hole, exactly the class validation's own green suite cannot see.
- **`1x` (`code-cleanups-0193`) AC-6 / `external-tracker-checkpoint` AC-6.** The dispatch cites a
  `1x AC-6` "validation passed but the audit refuted" case in the same lineage; `external-tracker-checkpoint`
  AC-6 is the archived self-oracle ("verified by static prose review of this entity's decision section")
  that motivated P1 in the first place. These are cited as the dogfood lineage, not re-litigated here.
- The proposal's own warning applies: this entity's two contract edits MUST themselves arrive at the
  gate with behavioral proof (the failing-first presence/banned-token tests over the real files), then
  pass a detached audit — dogfooding the very disciplines they encode.

## Out of scope

- **The code guard** (`status --validate` self-oracle lint + terminal-PASSED `--set` guard) — owned by
  `2a` (`require-external-proof-guard`). This entity ships only the FO-facing WORDING that guard backs,
  and points at it. Coordinate sequencing with 2a; do not re-file the guard.
- **Re-shipping anything `se` (#248) already landed** — the README principles/ergonomics, the ensign
  habits, the FO P3/posture/P4. Verified present this ideation; touching them is double-filing.
- **The portability/clean-room TEST** (the proposal's part C lineage) — `se` already encoded the
  PORTABILITY correction (TDD lives in the shipped ensign contract, no global CLAUDE.md dependency);
  no further portability work is in this entity.
- **README edits** — none survive; everything the README needed shipped with `se`. This entity touches
  the FO contract only.
- **Canonical upstream fan-out** — the two vendored edits flow to `~/git/spacedock` during the
  self-hosting reconciliation (same path `se`'s edits take), not as part of this entity's merge.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Both ACs
below are LEGITIMATE doc-as-deliverable cases: the deliverable is the FO's own loaded contract text, so
a presence/banned-token check over the real file is proof at the claim's own level (the workflow's named
exception — "a presence check over instruction files proving they carry a required clause or stay free of
a banned token is proof at the claim's own level"). The BEHAVIORAL teeth for the self-oracle guarantee
live in `2a`, explicitly out of scope here; these ACs do not claim a behavioral guarantee a code gate
should enforce.

**AC-1 — The FO AC-coverage cross-check requires an EXTERNAL oracle and refuses a self-oracle / pure-decision-to-terminal-PASSED, in the FO contract.**
Verified by: a Go test in `skills/integration/` (extend `working_principles_test.go` or `contract_gate_test.go`)
that reads `first-officer-shared-core.md`, isolates the `## Completion and Gates` AC-coverage-cross-check
paragraph, and asserts it carries the external-oracle requirement, the self-oracle "can never fail" refusal,
and the pure-decision-to-roadmap clause — AND that it points at the 2a backstop without re-asserting the
guard as this entity's deliverable. Failing-first: the assertion fails against the current FO core:134
(un-hardened original). The banned-token half asserts the entity does NOT re-file the guard as its own code.

**AC-2 — The shipped FO contract carries a concrete `## Detached Adversarial Audit` section: WHEN it triggers, WHAT it produces, HOW it is recorded.**
Verified by: a Go test in `skills/integration/` that reads `first-officer-shared-core.md`, isolates the
new section, and asserts it names (a) the high-stakes trigger surfaces (front-door / status / contract /
CI), (b) the detached-and-read-only requirement, and (c) the Feedback-Cycles recording/routing path. The
test fails against the current file (no such section exists). This is a property-of-the-text check at the
claim's own level — the deliverable is the contract an FO loads, not a runtime behavior.

**AC-3 — The audit formalization is grounded in a real catch, not invented.**
Verified by: the section (or the entity body's grounding note) cites #262's M1/M2 by their concrete
mechanism (the `strings.Count(...) > 0` skip-on-zero hole and the bare-`strings.Contains` disclaimer hole),
which are reproducible from the archived `binary-absent-fo-bootstrap/index.md` Feedback Cycles section.
The oracle is the archived entity's recorded text (external to THIS entity), so the citation can be checked
against a real prior artifact — it can fail if the cited mechanism does not match #262's record.

## Test plan

- **Test type:** Go integration tests in `skills/integration/`, extending the existing
  `working_principles_test.go` / `contract_gate_test.go` real-file contract-text harness (`foSharedCore`,
  `sectionAfter`, `repoRoot` helpers). Fixture cost: LOW — file reads + section isolation + substring /
  banned-token relationships over the real contract file. No binary build, no live workflow.
- **AC-1:** isolate the AC-coverage-cross-check paragraph; assert the three hardening clauses present;
  assert the 2a-backstop pointer present and the guard-code NOT re-declared. Confirm failing-first by
  running against the pre-edit file (the current un-hardened paragraph fails the external-oracle assertion).
- **AC-2:** isolate the `## Detached Adversarial Audit` section; assert the three axes (trigger surfaces,
  detached/read-only, Feedback-Cycles routing) present. Failing-first is trivial (the section does not exist).
- **AC-3:** assert the grounding note cites #262's M1 (`strings.Count > 0` skip-on-zero) and M2 (bare
  `strings.Contains` disclaimer) mechanisms; cross-checkable against the archived #262 Feedback Cycles text.
- **Detached audit (dogfood):** because the FO contract IS a high-stakes scaffolding surface, this entity's
  own merge must pass a detached adversarial audit per the very section it adds — the auditor confirms the
  presence tests would FAIL on an adversarial edit that deletes a required clause (no skip-on-zero / no
  bare-substring holes, the #262 lesson applied to this entity's own tests).
- **No spike needed:** the design composes only already-proven mechanisms — static contract-prose edits the
  FO loads at startup, proven by the `skills/integration/` real-file presence-test harness that `se` (#248)
  and #262 already shipped and that is green in-repo today. There is no parser round-trip, no new on-disk
  format, no runtime handoff. The one mechanism-level question — "is the FO contract already shipping these
  two edits (would this be double-filing)?" — was the riskiest unknown and was exercised FIRST this ideation:
  a grep of the live vendored `first-officer-shared-core.md` (cross-check un-hardened at :134; no
  `## Detached Adversarial Audit` section anywhere) and the canonical `~/git/spacedock` copy (lacks
  `## Working Principles` and the audit section), confirming the gap is real and the authoring direction is
  vendored-first. That diff is recorded in the reconciliation table above.

## Scaffolding-guardrail note

`skills/first-officer/references/` is protected scaffolding — both edits go through this tracked entity
dispatched to a worker in a worktree, never an FO hand-edit. Author in the VENDORED copy first (this project
is the leading edge; canonical lacks both sections), then flow upstream to canonical `~/git/spacedock` during
self-hosting reconciliation. The code guard lands only via `2a`. 0.19.4-class — folded by the captain
(2026-06-02), off the 0.19.2/0.19.3 critical path.

## Stage Report: ideation

- DONE: The encoding plan names exactly WHERE each of the 4 deliverable principles + the detached-audit discipline lands, and reconciles with se #248 + 2a so nothing is double-filed
  Reconciliation table (## Proposed approach A) maps every piece to shipped/this-entity/2a; verified by grep that se (#248) already landed README P1-P4 + ergonomics, ensign habits, and FO P3/posture/P4 — so the ONLY surviving FO-contract gaps are the AC-cross-check P1 hardening (FO core:134) and the new Detached-Audit section. README needs no edit; the code guard is 2a's, pointed-at not re-filed.
- DONE: Each principle claiming a BEHAVIORAL guarantee is backed by a code gate, not prose-only; genuine contract-prose principles are marked with the doc-as-deliverable justification
  The self-oracle BEHAVIORAL guarantee (a self-oracle AC mechanically can't reach terminal PASSED) is explicitly 2a's code guard, kept out of scope and only pointed at. The two edits this entity ships are marked contract-prose (ceiling wording-present), justified as the legitimate doc-as-deliverable case (agent-read contract; presence/banned-token test over the real file is proof at the claim's own level — the workflow's named exception, the #262 AC-1 precedent), proven failing-first in the working_principles_test.go / contract_gate_test.go harness.
- DONE: The detached-audit formalization is concrete — when it triggers, what it produces, how it is enforced/recorded — grounded in real catches
  ## Proposed approach C gives the full proposed section: trigger surfaces (front-door/status/contract/CI), read-only detached checkout, two-tier finding output, Feedback-Cycles routing. Grounded in #262 M1 (strings.Count>0 skip-on-zero) + M2 (bare strings.Contains disclaimer) — both validation-PASSED but audit-refuted, reproducible from the archived binary-absent-fo-bootstrap Feedback Cycles — plus the 1x/external-tracker AC-6 self-oracle lineage. AC-3 makes the grounding citation itself checkable against the archived #262 record.

### Summary

Reconciled the predates-se proposal against live contract state: se (#248) already shipped the README principles + ergonomics, the ensign habits, and the FO P3/posture/P4, so the only surviving gaps are (1) hardening the FO AC-coverage cross-check with the external-oracle / self-oracle-refusal / pure-decision-to-roadmap clause, and (2) a new ## Detached Adversarial Audit FO-contract section. Both land vendored-first in first-officer-shared-core.md (canonical lacks both — verified); the behavioral self-oracle guarantee stays with 2a's code guard (pointed-at, not re-filed). The two edits are honestly marked contract-prose with the doc-as-deliverable justification, proved by failing-first presence/banned-token tests in the existing real-file harness, and the audit section is grounded in #262's concrete M1/M2 test-strength holes. Riskiest unknown (am I double-filing?) was spiked first via a grep diff of vendored + canonical contracts, recorded in the reconciliation table.

## Stage Report: implementation

- DONE: Both FO-contract edits land in references/first-officer-shared-core.md (vendored copy): (a) harden the ## Completion and Gates AC-coverage cross-check with the external-oracle requirement + the self-oracle-can-never-fail refusal + the pure-decision-to-roadmap clause; (b) a new ## Detached Adversarial Audit section naming trigger surfaces (front-door/status/contract/CI), detached+read-only, and Feedback-Cycles recording/routing.
  Both edits in commit d77f7553 (worktree branch spacedock-ensign/encode-deliverable-principles). Edit (a): cross-check now carries "evidence must come from a check OUTSIDE the entity body", the can-never-fail self-proof refusal, and "do not advance it to a terminal PASSED ... the decision belongs in the roadmap". Edit (b): new `## Detached Adversarial Audit` section after `## Gate Presentation` with all three axes. Avoided the banned "oracle" token (renamed "self-oracle lint" → "self-proof lint") so the existing TestShippedInstructionsCarryNoInsiderJargon stays green.
- DONE: Go presence/banned-token tests in skills/integration/ enforce both edits over the REAL file, failing-first; AC-1 also asserts the 2a code guard is POINTED-AT, not re-declared here (no double-file).
  skills/integration/deliverable_principles_test.go: TestACCrossCheckRequiresExternalProof (AC-1) + TestDetachedAdversarialAuditSection (AC-2). Confirmed failing-first against the pre-edit file (5 missing-clause failures for AC-1; "no section" for AC-2), then green after the edits. AC-1 asserts the `status --validate` / `set guard` backstop pointers present AND bans self-claim phrases ("this cross-check enforces", "this edit implements the guard") so the 2a guard is pointed-at not re-filed. 29/29 integration pass.
- DONE: AC-3: the audit section (or grounding note) cites #262's M1 (strings.Count>0 skip-on-zero) + M2 (bare strings.Contains disclaimer) by mechanism, checkable against the archived binary-absent-fo-bootstrap Feedback Cycles. Do NOT edit README and do NOT re-file 2a's code guard.
  Grounding paragraph in the audit section cites M1 (`strings.Count(...); n > 0` skip-on-zero) + M2 (bare `strings.Contains` disclaimer), with `#262`. TestAuditGroundingCitesRealCatch asserts the mechanism strings present in the section AND cross-checks them against the archived `_archive/binary-absent-fo-bootstrap/index.md` record when the state checkout is reachable (skips gracefully in a code-only worktree — no hidden machine dependency). The archived record carries `strings.Count`/`strings.Contains` at lines 131-132, confirming the citation matches the real prior artifact. No README edit; 2a's guard not re-filed.

### Summary

Landed the two surviving FO-contract edits in the vendored first-officer-shared-core.md: hardened the AC-coverage cross-check (external-proof / can-never-fail self-proof / pure-decision-to-roadmap, pointing at 2a's backstop) and added a concrete `## Detached Adversarial Audit` section grounded in #262's M1/M2. Both are contract prose proved by failing-first presence/banned-token tests in skills/integration/deliverable_principles_test.go, dogfooding the very disciplines they encode — confirmed robust against clause-deletion adversarial edits (no skip-on-zero hole; deletion fails rather than skips). Key decision: renamed the backstop pointer "self-oracle lint" → "self-proof lint" to avoid the banned insider-jargon token "oracle" that the existing jargon test guards. Code committed on the worktree branch (d77f7553); this entity's own merge should get the detached audit it adds.
