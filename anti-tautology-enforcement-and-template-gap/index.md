---
id: azh879wdzm72ysxg16hbg39q
title: Falsifiable-test estate — contract evidence rule and the automated-gate decision, resolved once
status: ideation
source: "Two related findings from this session's audit work, not yet designed: (1) docs/dev's own Proof policy ('no prose-grep over instruction files', the detached adversarial audit, the AC template's 'Verified by: ... something outside this task body ... that can fail' clause) only catches tautology reactively, at high-stakes-surface merge time via human/reviewer judgment -- there is no standing, automatic check against the mirror-assertion / no-op-assertion patterns that four real tests in this repo turned out to have. (2) The commission-skill templates that scaffold NEW workflows do not carry equivalent discipline for GO TEST CODE tautology (as opposed to instruction-file/prose tautology, which is entity ey / proof-policy-shipped-scaffolding's separate, pre-existing scope -- see OVERLAP NOTE below): skills/commission/references/templates/development.md's AC template stub and skills/commission/SKILL.md's base AC template both lack docs/dev's own 'outside this task body'/'that can fail' clauses. Reference: github.com/kenn-io/middleman skills/testing-without-tautologies/SKILL.md. A design workflow + an independent fable review ran this session and both converged (mechanism choice, scope, sequencing all hold per fable's independent check) on: NO automated hard AST gate for mirror-assertions (an internal/testlint check for the OTHER pattern, assertion-free tests, is its own sibling entity: testlint-assertion-free-gate) -- instead extend the existing detached-adversarial-audit trigger to fire on AC PROVENANCE ('any AC whose expected value is derived from the same package's production functions or constants'), not the originally-drafted broader 'equality/byte-identity check' wording fable found would over-fire on nearly every unit test in this repo. Concrete diffs exist for docs/dev/README.md's Proof policy + pr-merge gate rule, development.md, and SKILL.md. CORRECTED accounting (fable caught this): internal/status/boot_probe_parity_test.go is NOT a stampID mirror-assertion instance -- it mirrors a different production CONSTANT (teamStateNeutralHint), not a function call; the confirmed mirror-assertion-via-shared-function count is 2 (native_new_test.go, zz_independent_parity_test.go), already covered under the separate tautological-test-fixes entity, don't double-count here. OVERLAP NOTE, UNRESOLVED -- flagging for the captain, not silently resolving: entity ey (proof-policy-shipped-scaffolding, filed 2026-06-04, pre-existing) targets the SAME file (skills/commission/references/templates/development.md) for a related-but-distinct concern -- porting the INSTRUCTION-FILE/prose tautology test (not code-test tautology) to shipped scaffolding, plus first-officer-shared-core.md and ensign-shared-core.md, with a heavier behavioral AC (a live scenario proving a validator REJECTS a presence-only proof). This entity's development.md diff and ey's development.md target could collide if ideated independently without coordination. Captain has not yet said how to reconcile (fold together, sequence, or keep fully separate with a coordination note in each) -- do not dispatch this entity's ideation until that's decided."
started: 2026-07-20T04:53:02Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0260-proportionality
group: test-cleanups
gates:
  version: 1
  current:
    gate: gate:docs-dev:az:ideation
    attempt: gate-attempt:az-ideation-1
  records:
    - id: gate:docs-dev:az:ideation
      stage: ideation
      current-attempt: gate-attempt:az-ideation-1
      attempts:
        - id: gate-attempt:az-ideation-1
          sequence: 1
          state: open
          briefing:
            id: briefing:docs-dev:az:ideation:attempt-1:revision-4
            digest: sha256:610fcfab5250d0d23eb7ed01f10eb702657b05b7be0692b194819715dff5bdc4
            room-ref: "./review/ideation/briefing-1"
          note: "Captain hold via float 2026-07-20 (resolution:actor-1784524247673759000): the FO gate summary was session-jargon dense and unreadable without full context. Attempt open; plain-language rewrite via the comm-officer before re-presentation."
---

Design and land a standing mechanism (not just reactive review) against tautological tests, and bring the commission-skill templates that scaffold new workflows up to docs/dev's own Proof-policy bar so new workflows don't inherit the gap.

## Scope trim (0260 re-lock 2026-07-20)

This entity owns the CONTRACT half only: the falsifiable-evidence rule (AC evidence must be
able to fail; "show the change that makes it fail") and removal of "5/5 passed is
sufficient" so gates read assertion content. The commission-template half of the original
scope moves to the template group (proof-policy-shipped-scaffolding and
template-rigor-propagation) — one owner per surface, no duplicate delivery.

## Merged scope (adopted cross-review re-lock, 2026-07-20)

Absorbs `anti-tautology-enforcement-output-grep-shape` (the output/prose-grep third shape). The cross-review caught a latent contradiction: this entity's banked direction leans against an automated prose-gate while the absorbed entity proposes extending one — the merged ideation must resolve gate-vs-no-gate DELIBERATELY, once, with the losing side recorded. Scope: the contract falsifiable-evidence rule ("show the change that makes it fail"; remove "5/5 passed is sufficient" so gates read assertion content) + the third-shape enforcement decision. The 8 concrete test fixes ship independently (their ideation is complete and valid under either outcome here).

## Additional evidence: the 0.25.1 AC-narrowing incident (2026-07-20)

See the synthesis addendum (`_evidence/0260-agent-derail-forensics/synthesis.md`). Design input for this entity's merged ideation: the falsifiable-evidence rule must cover not only proof that cannot fail but VALUE CLAIMS that quietly shrink — an AC edit weakening the value claim after a validation rejection is a design-reset event (captain-visible), never a task-internal edit; and evidence placement matters: proof runs at the exact place the failure occurs (0.25.1 proved the adapter, the live invocation stayed unproven, and the unproven place is exactly where the failure recurred).

## Problem

The opening line and frontmatter title predate the re-lock; the sections above and below are authoritative. Two things are settled: the scope is the CONTRACT half plus one gate-vs-no-gate decision, and the 8 concrete test fixes ship elsewhere (`fix-tautological-output-grep-tests`).

The gap this closes: today the workflow catches tautological tests only reactively, at high-stakes-surface merge time, by human/reviewer judgment — and that judgment was applied leniently four times in one session (synthesis incidents #6, #7; the output-grep shape passed FOUR Claude validation reviews). Two contract clauses make the miss structural, not accidental:

- `ensign-shared-core.md:85` tells the ensign report protocol to compress test evidence to a bare pass count ("Do not paste full test output — `5/5 passed` is sufficient"). The FO and captain then judge test evidence by count. A tautological test that always passes is fully satisfying evidence under this rule — the gate never sees what the test asserts (synthesis, clause-level driver: "hides test content from gates by design").
- The proof discipline verifies only the PASSING direction. "Here is the run that proves it" is a green run; nothing requires observing the FAILING case. A test that cannot fail satisfies the rule.

The merged-in `anti-tautology-enforcement-output-grep-shape` proposed to close this by EXTENDING automated prose-gating (a structural AST/lint gate over test literals, or a required committed roborev codex-lens check). This entity's banked direction leans the other way (no new standing automated gate; route the class to review-time discipline). Those cannot both hold. The entity exists to resolve that once and ship the contract prose the winning side needs.

## Decision: gate vs no gate (resolved once)

**Chosen: no new standing automated prose-gate.** The mirror-assertion shape and the absorbed output/prose-grep third shape both route to review-time discipline, realized by contract prose and the detached adversarial audit that already exists — not a new committed machine check. The mechanically-decidable shape (assertion-free tests) keeps its own separate gate entity (`testlint-assertion-free-gate`), where a genuine `contractlint`/`testlint` gate is justified because "a test contains zero assertions" is structurally decidable without semantic judgment.

**What the winning side ships (the contract half):** the falsifiable-evidence rule (AC evidence must name the change that makes it fail); removal of "`5/5 passed` is sufficient" so the report — and therefore the gate — reads what a test asserts and how it would fail, not a count; the captain's prose-grep ruling carried verbatim (one-off validation greps legal, committed prose-greps banned). The concrete before/after is in `## Documentation changes`.

**Losing side, recorded with why.** The absorbed proposal — a standing automated AST/lint gate whose tell is "a test asserts a string literal that appears verbatim in a repo file it read," or a required committed roborev codex-lens check for this shape — is rejected. Three grounds, each from a captain ruling that binds this sprint (index `## Constraints`, 2026-07-20):

1. **The cheapest-check-that-can-fail ordering.** The order is: a shipped system guard, an existing mechanical check, a falsifiable exercise, captain judgment — and new machinery only as a consent-gated last resort. A new AST/lint gate is that last resort. The cheaper levels reach the value and are shown to (see the spike): the detached adversarial audit is an EXISTING mechanism that already "constructs an adversarial edit the deliverable's own tests should catch, then confirms they do"; the falsifiable-evidence rule turns that into an author-side one-off exercise ("name the change that makes it fail"). New machinery is unjustified while the cheaper levels reach the value.
2. **It cannot compute a real value against an independent source that can diverge.** The structural tell — "the test's expected value is a literal that appears in a file the test read" — is itself prose-over-prose: the value the gate flags comes from the same file, the exact tautology one level up. It over-fires (fable's independent review this session found the broader equality/byte-identity wording flags nearly every ordinary unit test in the repo) and under-catches (the absorbed entity's own documented blind spots: the execute-a-static-passthrough-then-grep-a-constant form, and semantically-tautological tests that look structurally fine). A gate that flags legit tests and misses two of the three sub-shapes is worse than the review exercise it would replace. Whether a machine consumer PARSES the asserted string and whether a real behavior change flips it is a semantic judgment; a committed gate would encode false precision over it.
3. **New enforcement needs explicit captain approval and normally its own entity.** A standing automated gate is exactly that; folding it into a "resolved once, contract prose" entity violates the rule. The one shape that IS cleanly machine-decidable already has its own gate entity; the mirror/output shapes are not, so they route to review, not a gate.

The connective mechanism for the winning side is the detached adversarial audit firing on AC PROVENANCE ("any AC whose expected value is derived from the same package's production functions or constants") rather than only on the four high-stakes surfaces — the narrowed wording fable validated, not the byte-identity form that over-fires. This is a prose sharpening of an EXISTING audit's trigger, not new machinery. Because it widens when enforcement fires, it is flagged for explicit captain consent at the gate (Edit D); the core decision and the value stand without it, since the falsifiable-evidence rule already makes an author unable to name a falsifying change for a provenance-derived AC.

## Out of scope

- The 8 concrete tautological output-grep test fixes (`fix-tautological-output-grep-tests`) and the 2 confirmed mirror-assertion instances `native_new_test.go` / `zz_independent_parity_test.go` (`tautological-test-fixes`). This entity ships the CONTRACT and the DECISION, not the fixes; `boot_probe_parity_test.go` is correctly excluded from the mirror count (it mirrors a constant, not a same-package function — corrected accounting, not re-litigated here).
- The assertion-free gate (`testlint-assertion-free-gate`) — a real committed gate for the shape that warrants one.
- Commission-template propagation (`proof-policy-shipped-scaffolding`, `template-rigor-propagation`) — the template group carries the falsifiable-evidence rule into new-workflow scaffolding; this entity ships only the fleet-wide report-protocol fix (`ensign-shared-core.md`) and the dev-workflow proof surface (`docs/dev/README.md`).
- Any new committed test, gate, lint rule, or CI lane. The negative deliverable is the point (AC-4).

## Expected surface (contract prose only) + tolerance

**Two instruction files, ~7 net lines of contract prose, 0 Go source, 0 product LOC, 0 new committed tests/gates:**
- `skills/ensign/references/ensign-shared-core.md` — line 85, one sentence replaced (Edit A). Fleet-wide (the report protocol is fleet-wide, and the string physically lives here).
- `docs/dev/README.md` — Proof policy: +1 bullet (Edit B, falsifiable-evidence rule), append to the prose-grep bullet (Edit C, verbatim ruling), append to the detached-audit bullet (Edit D, consent-gated); plus the AC-template `Verified by` line sharpened (Edit B).

**Declared tolerance: 2×** (≤ ~14 net contract-prose lines). **Hard self-check:** any Go/test/product code, any new committed check or gate, or a third instruction file beyond these two trips a reconfirm — because "resolved once = decision + contract prose, not new infrastructure" IS the point.

## Riskiest-mechanism spike (done first, per ideation policy)

**Claim under test:** does the no-gate mechanism the decision routes to actually catch this class without over-firing — i.e., does a real behavior change flip a legit test but NOT a mirror-assertion test? Exercised end-to-end in a throwaway Go package (results recorded, not asserted):

- Two tests over a production function `Greeting() == "hello"`: a mirror-assertion test (`want := Greeting()`; expected derived from the same-package function — AC provenance) and a control (expected = the independent literal `"hello"`). Baseline: both pass.
- Applied the audit's claim-breaking edit — mutate the production function to return `"HELLO"`. Result: the mirror-assertion test **stayed GREEN** (PASS — the hole: a real behavior change did not flip it), the control went **RED** (`greet_test.go:17: got "HELLO" want "hello"`).

**Determination — no further spike needed.** The green/red divergence against the mutated production function (an independent source that can diverge) is the empirical foundation of the no-gate decision: the cheap one-off exercise catches the tautology and does not falsely flag the control, so a standing gate buys nothing the exercise does not. The throwaway seeds AC-1 directly. Proven mechanisms this design rides: (a) the detached adversarial audit already exists with recorded real catches (#262 two test-strength holes; `1x` and `external-tracker-checkpoint` AC-6 self-referential ACs; `7h` AC-3); (b) the AC-provenance predicate was already spiked by fable's independent review this session (the narrowed same-package form does not over-fire where byte-identity did; confirmed mirror count 2). No unverified parser round-trip, runtime handoff, on-disk format, or tool-flag dependency remains.

## Acceptance criteria

Each AC names a property of the finished entity and how it is verified. Per the decision, verifications are one-off validation exercises or git-state checks, never committed prose-greps.

**AC-1 (VALUE) — The no-gate mechanism catches the tautology shape and does not flag a legit control.** Constructing the representative pair (a mirror-assertion test whose expected value derives from a same-package production function; a control asserting the same output against an independent literal) and applying the audit's claim-breaking edit to that function, the mirror-assertion test stays GREEN (caught as a hole) while the control goes RED. The independent baseline that can move the wrong way: the four lenient reviews that passed this shape (synthesis incident #6), and the control itself — a mis-scoped mechanism would either turn the tautology red or the control green.
Verified by: the ideation spike formalized as a one-off exercise recorded in the validation report — a throwaway Go package, the mutation applied, the green-stays-green vs control-goes-red outcome observed (the ideation run: mirror PASS, control `greet_test.go:17 FAIL` under `hello`→`HELLO`). This is a real value against the mutated function, not a prose-grep, and not a committed test (consistent with the no-gate decision). This is the outcome the entity exists for.

**AC-2 — The gate-vs-no-gate contradiction is resolved once, with the losing side recorded.** The entity body carries one decision (no new standing automated gate), names the losing proposal (the absorbed AST/lint or required-roborev committed gate), and gives the three grounds (cheapest-check ordering; cannot compute an independent value / over-fires + under-catches; new-enforcement consent + own-entity), each tied to a binding captain ruling.
Verified by: a one-off validation read of `## Decision: gate vs no gate` (output pasted into the report, not a committed test) confirming the single decision, the named losing side, and the three grounds are present and mutually consistent with the shipped diff (which adds no gate). Per the captain's prose-grep ruling this is validation-time evidence.

**AC-3 — The three contract edits land as specified, and the report-protocol edit demonstrably makes a tautology and a legit test distinguishable in a stage report.** "`5/5 passed` is sufficient" is gone from `ensign-shared-core.md`, replaced by the name-the-assertion-and-its-falsifying-change instruction; the falsifiable-evidence rule and the `Verified by` "name the change that would make it fail" sharpening are in `docs/dev/README.md`; the captain's prose-grep ruling is carried verbatim.
Verified by: a one-off validation grep for the removed/added strings (output in the report, not committed); PLUS the falsifiability anchor — write the AC-1 pair's two report lines under both protocols and observe that the OLD form renders both as an identical "2/2 passed" while the NEW form yields distinguishable lines (the tautology's line cannot name a change that flips it). Distinguishable-vs-identical is a real divergence, not a spelling check over the edit.

**AC-4 — The negative deliverable holds: nothing new is committed as enforcement.** The shipped diff adds zero `*_test.go`, zero new gate/lint rule, zero CI lane, and touches only the two named instruction files (Edit D excepted only if the captain approves it, and it adds prose, not a test).
Verified by: `git diff --stat` / `git diff --name-only` against the merge base showing only `skills/ensign/references/ensign-shared-core.md` and `docs/dev/README.md` changed, no added `*_test.go`, no new files under `internal/contractlint` or `internal/testlint`. This is on-disk state that fails the moment someone reintroduces a committed gate — the decision with teeth.

## Test plan

- **Value exercise (throwaway Go mutation) → AC-1:** the ideation spike formalized — a throwaway package, one production-function mutation, two `go test` runs; passes only if the mirror test stays green and the control goes red. One-off, recorded in the validation report, not committed. Cost: minutes.
- **Decision read + edit-presence grep (one-off, NOT committed) → AC-2 / AC-3:** validation-time greps/reads whose output is pasted into the report; a grep over prose the implementer wrote cannot fail independently, so per the captain's ruling it is external evidence for that run, never a checked-in test.
- **Report-distinguishability check → AC-3:** author both report lines under old and new protocol; assert identical-under-old, distinguishable-under-new. A real divergence, cheap.
- **Git-state check → AC-4:** `git diff --name-only` / `--stat` against the merge base; asserts the exact file set and the absence of any added test/gate. Real state, falsifiable.
- **No Go unit tests, no product code, no committed prose-grep in this deliverable.** High-stakes note: `ensign-shared-core.md` is shipped contract, so the detached adversarial audit applies at merge — here it refutes nothing material (the change removes a count shortcut and adds falsifiability prose), and that clean-audit outcome is recorded at the gate.

## Documentation changes (concrete before/after — ideation proposes, implementation applies)

**Edit A — `skills/ensign/references/ensign-shared-core.md`, line 85 (report protocol size guideline):**

> - Before: `Size guideline: 30-50 lines max. One-line evidence per checklist item. Do not paste before/after diffs — the git log is the diff; cite commit SHAs. Do not paste full test output — `5/5 passed` is sufficient.`
> - After: `Size guideline: 30-50 lines max. One-line evidence per checklist item. Do not paste before/after diffs — the git log is the diff; cite commit SHAs. Do not paste full test output; for each test cited, name in one line what it asserts and the change that would make it fail — a bare pass count hides a tautology from the gate.`

**Edit B — `docs/dev/README.md`, Proof policy (add one bullet after the "No prose-grep over instruction files" bullet):**

> `- **Evidence must be able to fail.** Each AC's cited evidence names the concrete change that would flip it — the falsifying edit. An author who cannot name what would make the evidence fail has not shown it can fail, and the criterion does not count. This is the author-side obligation behind the validation bar below (a static check counts only when it tests a real value against an independent source that can diverge); the gate reads that falsifying change, not a pass count.`

**Edit B (cont.) — `docs/dev/README.md`, AC template `Verified by` line (~224), sharpen the existing clause:**

> - Before: `… something outside this task body that a future reader can reproduce and that can fail.}`
> - After: `… something outside this task body that a future reader can reproduce and that can fail; name the concrete change that would make it fail.}`

**Edit C — `docs/dev/README.md`, append the captain's ruling verbatim to the "No prose-grep over instruction files" bullet:**

> `Captain ruling (2026-07-20, verbatim): prose-greps are one-off validation evidence, never committed tests. A grep whose output is pasted into the validation report is legitimate external evidence for that run; the same grep committed as a test is banned — it re-asserts that the file contains what we wrote and cannot fail.`

**Edit D (consent-gated — captain approves at the gate) — `docs/dev/README.md`, append to the detached-adversarial-audit bullet:**

> `The audit also fires on AC provenance: when an AC's expected value is derived from the same package's production functions or constants, run the adversarial-edit check on it — that provenance is the tautology tell. Scope it to that provenance form; the broader equality/byte-identity form over-fires on ordinary unit tests.`

If the captain reads Edit D as new enforcement rather than an existing-mechanism trigger sharpening, drop it: the falsifiable-evidence rule (Edit B) already makes a provenance-derived AC unfalsifiable at authoring time, and the audit's existing "construct a claim-breaking edit" method covers it. The decision and AC-1's value stand either way.

## Boundary and notes

- **Leanness (0250/0260):** ~7 net contract-prose lines across two instruction files; the ensign report line is REPLACED not added; the dev-README additions ride an already-lazy proof surface (not boot-resident). Net contract-byte delta is small and measured at validation.
- **Practices its own thesis:** the entity ships a decision + prose and proves its value with a one-off falsifiable exercise and a git-state check, not a committed prose-grep — the exact discipline it writes into the contract.
- **Siblings/coordination:** `fix-tautological-output-grep-tests` (the 8 fixes) and `tautological-test-fixes` (the 2 mirror instances) ship the remediation; `testlint-assertion-free-gate` owns the one shape that warrants a real gate; the template group propagates the falsifiable-evidence rule to new-workflow scaffolding. This entity ships neither the fixes nor the templates.
- No minted identifier schemes or coined abstractions; bare ordinals throughout.

## Stage Report: ideation

- DONE: The gate-vs-no-gate contradiction is resolved ONCE, deliberately: the absorbed output-grep-shape proposal to extend automated prose-gating vs this entity's banked lean against it — one decision, the losing side recorded with why.
  `## Decision: gate vs no gate` records one decision (no new standing automated gate), names the losing proposal (the absorbed AST/lint or required-roborev committed gate), and gives three grounds each tied to a binding captain ruling.
- DONE: Contract half drafted as concrete before/after under the leanness constraint: the falsifiable-evidence rule, removal of "5/5 passed is sufficient" so gates read assertion content, and the captain's prose-grep ruling carried verbatim.
  `## Documentation changes` Edits A-C: A removes "`5/5 passed` is sufficient" from `ensign-shared-core.md:85`; B adds the falsifiable-evidence rule + `Verified by` sharpening to `docs/dev/README.md`; C carries the captain's ruling verbatim. Edit D (audit provenance trigger) flagged consent-gated.
- DONE: Written expected surface + tolerance declared; riskiest mechanism spiked first or "no spike needed" recorded with the proven mechanisms.
  `## Expected surface`: 2 instruction files, ~7 net prose lines, tolerance 2×, hard self-check. `## Riskiest-mechanism spike`: throwaway Go mutation observed mirror-test GREEN / control RED under `hello`→`HELLO`; determination "no further spike needed" with the proven mechanisms cited.

### Summary

Resolved the gate-vs-no-gate contradiction once: no new standing automated prose-gate; the mirror/output-grep shapes route to review-time discipline (the existing detached adversarial audit) plus three contract-prose edits, with the losing automated-gate proposal recorded and rejected on the captain's cheapest-check ordering, the cannot-compute-an-independent-value / over-fire+under-catch grounds, and the new-enforcement-needs-its-own-entity rule. The load-bearing mechanism was spiked live (mirror-assertion stays green under a production-function mutation while an independent-literal control goes red), seeding the VALUE AC. Deliverable is decision + ~7 lines of contract prose across two instruction files, 0 product code and 0 new committed tests/gates (AC-4 makes that negative deliverable a git-state check); Edit D (widening the audit trigger to AC provenance) is the one edit flagged for explicit captain consent at the gate.
