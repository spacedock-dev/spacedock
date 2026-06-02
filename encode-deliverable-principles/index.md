---
id: v1awnfhs996ykymv409anywh
title: Encode deliverable principles + template ergonomics into README, FO/ensign contract, and status guards
status: validation
source: FO triage (2026-06-01) — docs/dev/_proposals cleanup; consolidates the deliverable-principles study + the TDD/template-adoption ergonomics
started: 2026-06-02T04:57:30Z
completed:
verdict:
score: "0.30"
worktree: .worktrees/spacedock-ensign-encode-deliverable-principles
issue:
mod-block: 
pr: "#266"
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

## Stage Report: validation

- DONE: Reproduce AC-1/AC-2/AC-3 by MUTATING the real first-officer-shared-core.md (not re-reading)
  Ran 17 adversarial mutations on a DETACHED throwaway worktree (`/tmp/sd-audit-encode-deliverable`, removed after). AC-1: deleting each of the 5 cross-check clauses (external-proof, can-never-fail, terminal-PASSED, roadmap) + both backstop pointers (`status --validate`, `set guard`) FAILS TestACCrossCheckRequiresExternalProof. AC-2: removing the `## Detached Adversarial Audit` heading and the `read-only` / `adversarial edit` axes FAILS TestDetachedAdversarialAuditSection; `status` + `CI` trigger surfaces and the full `Feedback` routing axis are pinned. AC-3: deleting any of the 5 mechanism citations (`strings.Count`, `skip-on-zero`, `strings.Contains`, `disclaimer`, `#262`) FAILS TestAuditGroundingCitesRealCatch, and the cross-check against the archived `binary-absent-fo-bootstrap/index.md` (lines 131-132) runs (not skipped) and matches.
- FAILED: Confirm the tests carry NO skip-on-zero / bare-substring holes (the #262 lesson applied to v1's OWN tests)
  AC-1's banned-self-claim assertion genuinely rejects an injected `this cross-check enforces` phrase (M6 — not a spell-check). No skip-on-zero (`strings.Count > 0`) holes anywhere. BUT TestDetachedAdversarialAuditSection has a #262-M2-class bare-substring leak: the `front-door` and `contract` trigger markers also appear elsewhere in the section ("front-door change" in the What-it-produces bullet; "contract_gate_test.go" in the M1/M2 grounding). M12 proved an adversarial edit can gut the trigger-surface line down to only `status` + `CI` (dropping both `front-door` and `contract/scaffolding` surfaces) and the test STAYS GREEN. Full `go test ./skills/integration/` = 29/29 green; TestShippedInstructionsCarryNoInsiderJargon green; no stray `oracle` token in the FO contract.
- DONE: PASSED/REJECTED; confirm no README edit, 2a's code guard pointed-at not re-filed, change confined
  No README edit (`git diff --name-only main...HEAD` carries none). This entity's only commit (d77f7553) touches exactly `first-officer-shared-core.md` + `deliverable_principles_test.go`; the `contract_gate_test.go` diff vs main is from #262's own commits (b5e5ced7/4268dc1a) in the branch ancestry, not this entity. 2a's guard is pointed-at only — the banned-self-claim half (M6) confirms no re-filing. Verdict: REJECTED on the AC-2 material finding below.

### Summary

Independent detached audit (read-only, separate worktree). AC-1 and AC-3 are robust: every required clause and citation fails on deletion, the banned-self-claim half actively rejects an injected self-claim, and AC-3's archived-#262 cross-check runs against the real prior artifact. The full suite is 29/29 green and there is no stray `oracle` token. REJECTED on ONE material finding: TestDetachedAdversarialAuditSection's `front-door` and `contract` trigger-surface markers are bare-substring checks that leak into "front-door change" and "contract_gate_test.go" elsewhere in the same section, so an adversarial edit (M12) can remove two of the four named trigger surfaces from the trigger line and the test stays green. This is precisely the #262-M2 bare-`strings.Contains` leakage class the entity dogfoods against — two of four trigger assertions spell-check rather than pin their clause. Fix: scope the trigger-surface assertions to the opening trigger-surface sentence/list (as `acCrossCheckParagraph` and `startupStep1` already do for their clauses), or assert the surfaces with phrasing unique to the trigger line (e.g. `front-door launcher`, `shipped contract/scaffolding`), so removing a surface from the WHEN axis fails. Routes to implementation via Feedback Cycles.

## Feedback Cycles

**Cycle 1 (FO, 2026-06-02) — validation REJECTED + detached audit MATERIAL (4 findings, superset); routed to implementation.** The prose is correct and well-grounded (audit verified: scope surgical +13/-1, no README/SKILL edit, jargon rename complete, 2a pointed-at not re-filed, the #262 citation matches the archived record). The defect is entirely in the TESTS: `deliverable_principles_test.go` pins TOKEN PRESENCE via bare `strings.Contains`, never the rule's MEANING — so a meaning-inverting rewrite that keeps the token passes. This is the exact #262 M1/M2 class the entity exists to teach against, reproduced in its own tests. The validator's 17 mutations were deletions (all caught); the audit added inversion/skip tests the validator missed. Each proven by an adversarial edit that left the suite green:

- **M-1 (AC-1 self-proof rule — M2 disclaimer).** `:60` checks `Contains(para, "can never fail")`; rewriting the refusal to its opposite ("…can never fail, but that concern no longer applies; self-review is acceptable") stays green. Fix: reject the negated/disclaimer form (assert the rule's polarity — ban "no longer applies"/"is acceptable"/"historically"-type inversions near the clause), not just the token.
- **M-2 (AC-1 banned-self-claim is a 3-string denylist).** `:86-90`; an unlisted self-sufficiency phrasing ("this very cross-check is itself the binding guarantee … no external code gate is required") passes. Fix: make it a relationship check — the paragraph must point OUT to 2a/external proof and must not claim self-sufficiency; document any residual limit honestly.
- **M-3 (AC-2 section is token-salad-satisfiable).** `:111-152`; replacing the whole section with "intentionally left blank; do not perform any audit" + a token-salad line stays green. Fix: assert the section STRUCTURE/instruction (load-bearing bullets as distinct items with their verbs; reject a "left blank"/"do not perform" inversion), not tokens-anywhere. Subsumes the validator's narrower `front-door`/`contract` trigger-marker leak — scope those to the trigger line or use trigger-line-unique phrasing.
- **M-4 (AC-3 citation — M2 inversion AND M1 skip-on-zero).** `:190-201`; inverting the citation ("unrelated to #262, which found nothing; ignore any mention of strings.Count…") stays green. WORSE: the archived-record cross-check at `:207-211` is `if !ok { return }` skip-on-ABSENCE, and the record is ABSENT in a code-only worktree / CI — so the one protective arm NEVER runs in CI (the validator saw it run only because it had the state checkout). Fix: (a) reject the citation inversion; (b) make the record cross-check NOT skip-on-absence in CI — vendor the relevant #262 fragment into the test as a testdata fixture so it always runs, or hard-fail when the record is expected and absent.

**Unifying requirement:** pin the rule's RELATIONSHIP/polarity, not token presence; and specifically close the two named #262 holes — no skip-on-absence (M1) and reject the negated/disclaimer form (M2) — since those are the exact failure modes this entity teaches against. Where a prose test genuinely cannot pin full meaning (the doc-as-deliverable ceiling this entity itself documents), that is acceptable ONLY for the part whose behavioral teeth live in 2a — record that boundary explicitly; it does not excuse the M1/M2 holes above. Re-run failing-first: each adversarial edit (inversion / left-blank / citation-flip / absent-record) must now FAIL. **Polish:** align the test docstrings still saying "self-oracle" (`:34,43,45,60`) with the shipped "self-proof" rename.

## Stage Report: implementation (cycle 1)

- DONE: M-1 (AC-1 self-proof polarity) — reject the negated/disclaimer form, not just require the token.
  TestACCrossCheckRequiresExternalProof now bans inversion phrases near each clause ("no longer applies", "is acceptable", "historically", "but that concern") AND asserts the relationship: "can never fail" must be FOLLOWED by "does not satisfy" in the same clause. Mutation (`can never fail, but that concern no longer applies; self-review is acceptable`) → RED at :102 and :115. Commit 0cd1f742.
- DONE: M-2 (AC-1 banned-self-claim) — make it a RELATIONSHIP check; document residual limit honestly.
  Replaced the 3-string denylist with: (a) required backstop pointers `status --validate` + `set guard` + `backstop` (points OUT to 2a), and (b) a self-sufficiency-pattern ban ("no external code gate", "itself the binding guarantee", "self-sufficient", …). The unlisted phrasing from the finding ("this very cross-check is itself the binding guarantee … no external code gate is required") → RED at :102 ("no external") and :159 (two patterns). Residual prose-ceiling limit recorded in the test docstring: a wholly novel self-sufficiency phrasing could still slip a prose check — that ceiling is the doc-as-deliverable limit; the behavioral teeth are 2a's `status --validate` guard, out of scope here.
- DONE: M-3 (AC-2 section structure) — assert STRUCTURE, reject inversion, scope trigger surfaces to the WHEN line.
  TestDetachedAdversarialAuditSection now asserts three DISTINCT load-bearing bullets (`- **Detached and read-only.**`, `- **What it produces.**`, `- **How it is recorded and routed.**`) each carrying its instruction, plus the `DETACHED checkout`/`read-only`/`REFUTE` requirements, bans "left blank"/"do not perform"/"no audit required" inversions, and scopes the four trigger surfaces to the opening WHEN sentence with trigger-line-unique phrasing (`front-door launcher`, `` `status` mutation/guard ``, `shipped contract/scaffolding`, `CI/release machinery`) — subsuming the validator's narrower front-door/contract bare-substring leak. Mutation (whole section → "left blank; do not perform any audit" + token-salad) → RED at :231 (×2 trigger), :263 (×3 bullets), :285 (×2 requirements). Commit 0cd1f742.
- DONE: M-4 (AC-3 citation — reject inversion AND kill skip-on-absence).
  TestAuditGroundingCitesRealCatch bans citation-inversion phrases ("unrelated to #262", "found nothing", "ignore strings.Count…") alongside the affirmative M1 (`strings.Count`/`skip-on-zero`) + M2 (`strings.Contains`/`disclaimer`) + `#262` citations. The archived-record cross-check now reads a VENDORED in-repo testdata fixture (`skills/integration/testdata/binary-absent-fo-bootstrap-feedback-cycles.md`, content-identical to the archived `_archive/binary-absent-fo-bootstrap/index.md` Feedback Cycles fragment) and HARD-FAILS when absent — no more `if !ok { return }` skip-on-absence, so the cross-check runs in CI / a code-only worktree. Mutation (citation inversion) → RED at :336; deleting the fixture → RED (hard-fail). Commit 0cd1f742.
- DONE: Polish — align test docstrings "self-oracle" → "self-proof"; full suite green.
  Rewrote all docstrings to "self-proof"; `grep -i oracle` over the test file, the fixture, and the FO contract returns nothing. Also scoped the sibling TestStartupGateGuidanceHasSingleProseSource to skip `testdata/` (the vendored fixture carries "spacedock doctor" and would otherwise be flagged as a second gate-prose mirror) — mutation-verified that a real mirror outside testdata still fails. `go test ./skills/integration/` 29/29 green, `go vet` clean, `gofmt -l` clean.

### Summary

Closed all four feedback-cycle-1 findings: the tests now pin RELATIONSHIP/POLARITY (inversion bans, ordered-consequence assertion, distinct-bullet structure, scoped trigger sentence) rather than token presence, and the AC-3 archived-record cross-check runs ALWAYS via a vendored in-repo #262 fixture that hard-fails when absent — closing both named #262 holes (no skip-on-absence M1; reject the disclaimer/inversion form M2) in this entity's own tests. Each adversarial edit (self-proof inversion / self-sufficiency claim / left-blank section / citation flip / absent fixture) was mutation-verified RED, then the real file restores to 29/29 green. The shipped contract PROSE is unchanged from cycle 0 (zero diff); only the test surface changed, plus a one-bullet testdata-skip to the sibling single-source test that my fixture would otherwise have tripped. Honestly recorded the residual prose-ceiling limit on M-2 (a novel self-sufficiency phrasing could slip; behavioral teeth are 2a's). Commit 0cd1f742 on branch spacedock-ensign/encode-deliverable-principles.

## Stage Report: validation (cycle 1)

- DONE: M-1 — self-proof refusal polarity inversion now FAILS TestACCrossCheckRequiresExternalProof
  Rewrote the clause to "can never fail, but that concern no longer applies; self-review is acceptable" on a DETACHED worktree (0cd1f742) → FAILS via the inversion denylist (`self-review is acceptable`, `no longer applies`, `but that concern`) AND the consequence-follows check (deliverable_principles_test.go:115 — `can never fail` not followed by `does not satisfy`). Was green in cycle 0.
- DONE: M-2 — unlisted self-sufficiency injection now FAILS
  Injected "this very cross-check is itself the binding guarantee … no external code gate required" → FAILS on both the `no external code gate` and `is itself the binding guarantee` self-sufficiency patterns (deliverable_principles_test.go:159).
- DONE: M-3 — token-salad section replacement AND dropped trigger surfaces now FAIL TestDetachedAdversarialAuditSection
  M-3a (replace whole section with "intentionally left blank; do not perform any audit" + token-salad) → FAILS on no-op inversion guard + missing distinct bullets + missing WHEN surfaces. M-3b (drop front-door + contract from the WHEN line, the cycle-0 material hole) → FAILS: trigger surfaces are now scoped to the WHEN sentence with trigger-line-unique phrasing (`front-door launcher`, `shipped contract/scaffolding`), so the leak into "front-door change"/"contract_gate_test.go" no longer green-lights it. Hole closed.
- DONE: M-4 — inverted #262 citation FAILS and fixture deletion HARD-FAILS (not skip)
  M-4a (negate grounding: "unrelated to #262, which found nothing; ignore strings.Count…") → FAILS on the citation-inversion guard. M-4b (delete skills/integration/testdata/binary-absent-fo-bootstrap-feedback-cycles.md) → HARD-FAILS via t.Fatalf at deliverable_principles_test.go:359, never skips — the #262-M1 skip-on-absence hole in this entity's own test is closed. Fixture verified byte-for-byte verbatim against the real archived _archive/binary-absent-fo-bootstrap/index.md M1/M2 bullets (lines 131-132).
- DONE: Regression — single-source invariant NOT weakened by the testdata-skip
  R1: a real gate-prose mirror (markers Contract version gate / per-class remedy / spacedock doctor) placed OUTSIDE testdata/ (skills/first-officer/references/) STILL FAILS TestStartupGateGuidanceHasSingleProseSource (reports two sources). R2: the fixture carries 2 gate markers, so the SkipDir is load-bearing and narrowly scoped to testdata/ — it excludes only inert vendored fixtures, not real prose mirrors.
- DONE: Hygiene — full suite green, vet/gofmt clean, no stray "oracle"
  `go test ./skills/integration/ -count=1` = 29/29 passed; `go vet ./skills/integration/` clean; `gofmt -l` on both changed .go files lists nothing; grep for "oracle" in first-officer-shared-core.md + deliverable_principles_test.go is empty.

### Summary

PASSED. Re-validated all four cycle-0 findings by MUTATING the real files on a detached throwaway worktree (0cd1f742, removed after) — every adversarial edit that was green in cycle 0 now FAILS, and the fixture-deletion case HARD-FAILS rather than skips. The cycle-0 material finding (AC-2 trigger-surface bare-substring leak into "front-door change"/"contract_gate_test.go") is closed by scoping the WHEN-axis assertions to the opening trigger sentence with trigger-line-unique phrasing. The new tests pin relationship/polarity (inversion denylists + consequence-follows + distinct-bullet structure), not token-presence, dogfooding the #262 M1/M2 lesson on this entity's own tests. The testdata-skip on the sibling single-source test is narrowly scoped (verified: a real mirror outside testdata/ still fails). Full suite 29/29, vet/gofmt clean, no stray "oracle". Contract prose unchanged from cycle 0; change confined to deliverable_principles_test.go + the testdata fixture + a one-bullet contract_gate_test.go testdata-skip.

## Feedback Cycles

**Cycle 2 (FO + captain, 2026-06-02) — RESHAPE: wrong home + banned proof. Re-home the disciplines to the dev workflow + template; drop the grep tests.** Captain review of the held PR raised two points, both decisive:

1. **Wrong home.** The `## Detached Adversarial Audit` discipline AND the AC-coverage-cross-check hardening were placed in the *universal* `first-officer-shared-core.md`. But that contract governs EVERY commissioned workflow, dev and non-dev. The deliverable principles (external-oracle proof / behavioral-not-grep / enforce-in-code / spike-the-risk) are **dev-centric** — a non-development workflow (research, ops, writing/triage) has ACs whose proof is a published artifact, a metric, or a human review, NOT a test/command/exit-code. Imposing the external-proof constraint universally is a category error. So these disciplines are **dev-workflow policy**, not universal FO mechanism.
2. **Banned proof.** `deliverable_principles_test.go` is the tautological grep-over-prose antipattern this entity itself forbids (P2). A contract/prose edit's proof is review + the *code* guard (2a), never a grep that spell-checks the prose. The entire cycle-1 M-1..M-4 saga was hardening a construct we ban.

**Reshape (this cycle):**
- **REVERT both edits** to `skills/first-officer/references/first-officer-shared-core.md`: remove the AC-cross-check hardening and the `## Detached Adversarial Audit` section. The universal FO contract returns to its prior generic state (it keeps only "each AC needs an evidence citation"; no external/behavioral qualifier).
- **ADD** the dev-specific guidance to **`docs/dev/README.md`** (this workflow's own contract): the AC-external-proof hardening as validation-stage / gate guidance, and the detached-adversarial-audit discipline as a validation/review practice (when it triggers, what it produces, how it's recorded — grounded in the #262 / 1x / 7h catches).
- **ADD** both as **recommended (opt-in)** practices in the commission **dev-template** (the workflow-README template under `skills/commission/`), so new dev workflows can adopt them without the universal FO contract mandating them.
- **DELETE** `skills/integration/deliverable_principles_test.go` AND its testdata fixture; **REVERT** the `contract_gate_test.go` testdata-skip (it only existed to accommodate that fixture). No grep self-proof.
- **Design note for `2a`** (not this entity): the `status --validate` external-proof guard must be **workflow-opt-in** (a workflow declares it wants the guard), not a universal binary behavior — same dev-specificity reason. Record this so 2a's ideation picks it up.

**v1's proof after the reshape:** doc edits (the dev README + the commission dev-template), reviewed at the gate; the behavioral enforcement is delegated to `2a` (opt-in). This is an honest documentation/policy deliverable — no self-proof grep, no universal-contract overreach. Record the boundary explicitly in the entity body.

## Design note for `2a` (workflow-opt-in guard)

Recorded here so `2a` (`require-external-proof-guard`) picks it up in its own ideation; this is NOT this entity's deliverable.

The `status --validate` external-proof guard (flagging an AC whose proof clause is self-referential / circular) and the terminal-PASSED `--set` guard must be **workflow-opt-in**, not universal binary behavior. The reason is the same dev-specificity that drove this entity's cycle-2 reshape: the "an AC's proof must be a test / command / exit-code / on-disk state" constraint is a **development-workflow** policy. A non-development commissioned workflow (research, ops, writing/triage) legitimately has ACs whose proof is a published artifact, a metric, or a human review — none of which a code/command external-proof guard should reject. So a workflow must **declare** that it wants the guard (e.g. a workflow-README flag the binary reads), and `status --validate` applies the self-referential-proof check only for workflows that opted in. Imposing it on every workflow would be the same category error this entity just corrected at the contract level.

The dev README and the commission dev-template both phrase the external-proof rule as the dev-workflow policy and point at this opt-in guard as the behavioral backstop "when a workflow wants it" — so the policy text and 2a's guard stay consistent: prose states the dev policy, the opt-in guard enforces it for workflows that declare it.

## Stage Report: implementation (cycle 2)

- DONE: REVERT both edits to skills/first-officer/references/first-officer-shared-core.md (remove AC-cross-check hardening AND the Detached Adversarial Audit section; return the universal FO contract to its prior generic state).
  Reverted via `git checkout d77f7553^ -- first-officer-shared-core.md`; working-tree diff vs the pre-edit base d77f7553^ is EMPTY (byte-identical). The AC-coverage cross-check reads only "confirm each AC has at least one evidence citation" again — no external/behavioral qualifier; no `## Detached Adversarial Audit` section; `grep -i` for "detached adversarial"/"OUTSIDE the entity body"/"self-proof lint" returns 0. Commit 413340e9.
- DONE: ADD the dev-specific guidance to docs/dev/README.md — AC-external-proof hardening as validation-gate guidance, and the detached-adversarial-audit discipline as a validation/review practice (when/what/how), grounded in #262/1x/7h.
  Added an external-proof requirement bullet to the validation stage Outputs (evidence must be a test/command/exit-code/file/on-disk-state outside the task body; reject self-referential ACs; pure-decision → roadmap not PASSED; explicitly marked dev-workflow policy). Added a `**Detached adversarial audit:**` validation bullet with the four axes (when it triggers / what it produces / how recorded / why), grounded in #262 (strings.Count skip-on-zero + bare strings.Contains disclaimer), 1x/external-tracker AC-6 (self-referential proof), and 7h release-notes AC-3 (tag-cut folded notes into subject not body). Commit 413340e9.
- DONE: ADD both as RECOMMENDED (opt-in) practices in the commission dev-template (skills/commission/references/templates/development.md).
  Added a `## Recommended practices (opt-in)` section with two subsections (external-proof acceptance criteria; detached adversarial audit), framed as recommended-not-mandatory and explicitly noting the universal FO contract does not impose them because a non-dev workflow's proof may be an artifact/metric/review. Points at the workflow-opt-in `status --validate` guard as the behavioral backstop.
- DONE: DELETE skills/integration/deliverable_principles_test.go and its testdata fixture; REVERT the contract_gate_test.go testdata-skip.
  `git rm` removed the grep test and `testdata/binary-absent-fo-bootstrap-feedback-cycles.md` (the now-empty testdata dir went with it). contract_gate_test.go's testdata-skip reverted to the bare `if info.IsDir() || filepath.Ext(p) != ".md"` walk — working-tree diff vs its pre-cycle-1 state (4268dc1a) is EMPTY. No dangling references to the deleted helpers (`grep` across skills/internal/cmd is clean). `go test ./...` 712 passed, `go vet` clean, `gofmt -l` clean.
- DONE: In the entity body, add a design note for 2a — its status --validate external-proof guard must be workflow-OPT-IN, not universal — same dev-specificity reason.
  Added the `## Design note for 2a (workflow-opt-in guard)` section above: the guard applies only to workflows that declare they want it, because the external-proof constraint is dev-workflow policy and a non-dev workflow's AC proof may be an artifact/metric/review. Recorded so 2a's ideation picks it up; explicitly not this entity's deliverable.

### Summary

Reshaped per captain review: the prose was correct but lived in the wrong home (the universal FO contract) and was proven the wrong way (a grep-over-prose self-proof — the exact P2 antipattern this entity teaches against). Reverted both first-officer-shared-core.md edits (working tree now byte-identical to the pre-edit base), deleted the grep test + fixture, and reverted the testdata-skip it required. Re-homed the external-proof and detached-adversarial-audit disciplines as DEV-specific guidance in docs/dev/README.md (validation-gate policy) and as opt-in recommended practices in the commission development template, grounded in the #262/1x/7h catches. The honest proof boundary is recorded explicitly: this is a documentation/policy deliverable reviewed at the gate, with behavioral enforcement delegated to 2a as a workflow-opt-in guard (design note recorded for 2a's ideation). go test ./... 712 passed, go vet clean, gofmt clean. Commit 413340e9 on branch spacedock-ensign/encode-deliverable-principles.

## Stage Report: validation (cycle 2 — the reshape)

- DONE: REVERT complete — FO contract byte-identical to pre-edit base
  `git diff d77f7553^ HEAD -- skills/first-officer/references/first-officer-shared-core.md` is EMPTY (0 bytes). Residual hardening phrases: "OUTSIDE the entity body"=0, "self-proof"=0, "Detached Adversarial Audit"=0. "can never fail"=1 — but that lone hit is the pre-existing `**Prefer a code gate over a prose-only rule.**` paragraph in `## Working Principles` (line 383), byte-identical in the base d77f7553^ (verified: present in `git show d77f7553^:…`), i.e. se #248 text, NOT this entity's hardening residual. Empty full-file diff proves the revert is complete.
- DONE: RE-HOME landed (not lost) in the two dev-specific homes
  docs/dev/README.md `### validation` stage (lines 110, 117) carries the external-proof AC bullet ("must come from a check OUTSIDE the task body") AND the `**Detached adversarial audit:**` practice bullet, both marked dev-workflow policy with the explicit why-not-universal framing ("This is dev-workflow policy… A non-development workflow's AC proof may legitimately be a published artifact, a metric, or a human review"). skills/commission/references/templates/development.md carries `## Recommended practices (opt-in)` with both disciplines, marked "recommended, not mandatory" and explaining the universal FO contract does not impose them. Directly addresses the captain's wrong-home complaint.
- DONE: DELETION clean — grep test + fixture gone, testdata-skip reverted, no dangling refs, build/test/lint green
  deliverable_principles_test.go and testdata/binary-absent-fo-bootstrap-feedback-cycles.md both gone (testdata/ dir removed). contract_gate_test.go byte-identical to 4268dc1a (`git diff 4268dc1a HEAD` = 0 bytes — testdata-skip fully reverted). Zero dangling references to any deleted helper/test (grepped skills/ internal/ cmd/ for deliverable_principles, archived262, the fixture name, acCrossCheckParagraph, auditTriggerSentence, auditSection, and all three deleted Test names — all 0). `go test ./...` = 10 packages ok, 0 FAIL; `go vet ./...` clean. gofmt: the entity-touched contract_gate_test.go is clean; the one gofmt-dirty file (internal/status/enum_scope_test.go) is pre-existing base drift — not in this entity's commit, 0 diff vs d77f7553^, and dirty in the base too.

### Summary

PASSED. The reshape correctly addresses both captain rejections of cycle 1. WRONG HOME: the FO-contract edits are reverted byte-for-byte (empty diff vs d77f7553^; the lone "can never fail" hit is pre-existing se #248 prose, not a residual), and the disciplines are re-homed to docs/dev/README.md (dev-workflow policy) + skills/commission/.../development.md (opt-in recommended practices) — both with explicit why-not-universal framing. WRONG PROOF: the tautological grep self-proof test and its fixture are deleted, the contract_gate_test.go testdata-skip that existed only for that fixture is reverted to 4268dc1a, and there is intentionally NO replacement grep test (verified absent, not expected). No dangling references; `go test ./...` 10/10 ok, vet clean, the touched Go file gofmt-clean (the one dirty file is pre-existing base drift outside this entity). Entity scope is exactly the 6 expected files. This entity stays HELD for the captain's own review; no merge recommended here.
