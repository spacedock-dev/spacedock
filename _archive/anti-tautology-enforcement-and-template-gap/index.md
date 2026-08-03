---
id: azh879wdzm72ysxg16hbg39q
title: Falsifiable-test estate — contract evidence rule and the automated-gate decision, resolved once
status: done
source: "Two related findings from this session's audit work, not yet designed: (1) docs/dev's own Proof policy ('no prose-grep over instruction files', the detached adversarial audit, the AC template's 'Verified by: ... something outside this task body ... that can fail' clause) only catches tautology reactively, at high-stakes-surface merge time via human/reviewer judgment -- there is no standing, automatic check against the mirror-assertion / no-op-assertion patterns that four real tests in this repo turned out to have. (2) The commission-skill templates that scaffold NEW workflows do not carry equivalent discipline for GO TEST CODE tautology (as opposed to instruction-file/prose tautology, which is entity ey / proof-policy-shipped-scaffolding's separate, pre-existing scope -- see OVERLAP NOTE below): skills/commission/references/templates/development.md's AC template stub and skills/commission/SKILL.md's base AC template both lack docs/dev's own 'outside this task body'/'that can fail' clauses. Reference: github.com/kenn-io/middleman skills/testing-without-tautologies/SKILL.md. A design workflow + an independent fable review ran this session and both converged (mechanism choice, scope, sequencing all hold per fable's independent check) on: NO automated hard AST gate for mirror-assertions (an internal/testlint check for the OTHER pattern, assertion-free tests, is its own sibling entity: testlint-assertion-free-gate) -- instead extend the existing detached-adversarial-audit trigger to fire on AC PROVENANCE ('any AC whose expected value is derived from the same package's production functions or constants'), not the originally-drafted broader 'equality/byte-identity check' wording fable found would over-fire on nearly every unit test in this repo. Concrete diffs exist for docs/dev/README.md's Proof policy + pr-merge gate rule, development.md, and SKILL.md. CORRECTED accounting (fable caught this): internal/status/boot_probe_parity_test.go is NOT a stampID mirror-assertion instance -- it mirrors a different production CONSTANT (teamStateNeutralHint), not a function call; the confirmed mirror-assertion-via-shared-function count is 2 (native_new_test.go, zz_independent_parity_test.go), already covered under the separate tautological-test-fixes entity, don't double-count here. OVERLAP NOTE, UNRESOLVED -- flagging for the captain, not silently resolving: entity ey (proof-policy-shipped-scaffolding, filed 2026-06-04, pre-existing) targets the SAME file (skills/commission/references/templates/development.md) for a related-but-distinct concern -- porting the INSTRUCTION-FILE/prose tautology test (not code-test tautology) to shipped scaffolding, plus first-officer-shared-core.md and ensign-shared-core.md, with a heavier behavioral AC (a live scenario proving a validator REJECTS a presence-only proof). This entity's development.md diff and ey's development.md target could collide if ideated independently without coordination. Captain has not yet said how to reconcile (fold together, sequence, or keep fully separate with a coordination note in each) -- do not dispatch this entity's ideation until that's decided."
started: 2026-07-20T04:53:02Z
completed: 2026-07-20T15:52:09Z
verdict: passed
score:
worktree: .worktrees/spacedock-ensign-anti-tautology-enforcement-and-template-gap
issue:
sprint: 0260-proportionality
group: test-cleanups
pr: pr-merge:536
archived: 2026-07-20T15:52:09Z
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

## Stage Report: implementation

- DONE: Edits A-D land as the approved decision + contract prose: A removes the "5/5 passed is sufficient" shortcut from skills/ensign/references/ensign-shared-core.md:85; B adds the falsifiable-evidence rule to docs/dev/README.md's Proof policy and sharpens the AC-template `Verified by` line; C appends the captain's prose-grep ruling verbatim; D appends the audit-trigger widening to the detached-adversarial-audit bullet. Edit D is CAPTAIN-APPROVED — see the entity's edit-d-resolution record.
  Commit `2245554f` on `spacedock-ensign/anti-tautology-enforcement-and-template-gap`. A at `ensign-shared-core.md:85` (sentence replaced, not added); B as a new Proof-policy bullet after the prose-grep bullet plus the `Verified by` clause at the AC template; C appended to the prose-grep bullet verbatim from `## Documentation changes`; D appended to the existing detached-audit bullet, wording as approved (AC-provenance form, not byte-identity).
- DONE: AC-4's negative deliverable holds and is proven by git state, not assertion: the shipped diff adds zero *_test.go, zero new gate/lint/CI lane, and touches only the two named instruction files. Edit D adds prose to an EXISTING audit's trigger — it must not become a new mechanism.
  `git diff --name-only` vs merge base returns exactly `docs/dev/README.md` and `skills/ensign/references/ensign-shared-core.md`; `--diff-filter=A` × `_test.go` = 0; files matching `\.go$|\.ya?ml$|contractlint|testlint` = 0. `--stat`: 2 files, 5 insertions, 4 deletions (net +1 line). Edit D touched no file other than the README bullet it extends.
- DONE: The decision record names a single decision, its losing side, and three grounds — and the mutation exercise that seeded the VALUE AC is reproducible from the record (mirror-assertion stays green under a production-function mutation while the independent-literal control goes red).
  Rerun from the record, throwaway package, not committed (see `### Reproducing the value exercise`): baseline both PASS; after mutating `Greeting()` `"hello"`→`"HELLO"`, `TestGreetingMirror` **PASS** (the hole) and `TestGreetingControl` **FAIL** `got "HELLO" want "hello"`. `## Decision: gate vs no gate` (body, unchanged this stage) carries one decision, names the losing AST/lint-or-required-roborev proposal, and gives the three grounds.

### Reproducing the value exercise

A throwaway Go module outside the worktree, three files, ~15 lines. `greet.go` exports `Greeting() string { return "hello" }`. `greet_test.go` holds the pair: `TestGreetingMirror` takes `want := Greeting()` then asserts `Greeting() == want` (expected value derived from the same package's production function — the AC-provenance shape Edit D now triggers the audit on); `TestGreetingControl` asserts `Greeting() == "hello"` against an independent literal. Run `go test ./... -v`, edit the return to `"HELLO"`, run again.

What this asserts and what would make it fail: it asserts the mutation divergence, so it fails if the mirror test goes red (the shape is not actually tautological) or the control stays green (the mutation did not reach the assertion). Both would refute the no-gate decision's premise; neither happened. Observed line cite was `greet_test.go:16`, not the `:17` recorded at ideation — same failure, different line because my file carries two comment lines in different positions. Line numbers are an artifact of the throwaway; the green/red divergence is the claim.

### Edit A's falsifiability anchor (report-distinguishability)

The two tests above, written as report evidence under each protocol:

- OLD ("`5/5 passed` is sufficient"): both render as `2/2 passed` — one string, no way for a gate to tell the tautology from the real test.
- NEW (this commit): control → "asserts `Greeting()` returns the literal `hello`; mutating the production function to return `HELLO` makes it fail." Mirror → "asserts `Greeting()` equals `Greeting()`; **no change to the production function makes it fail** — the mutation left it green."

Identical under the old protocol, distinguishable under the new one, and the mirror line cannot name a falsifying change at all. That inability is the tell Edit A hands the gate.

### Summary

Applied Edits A-D as contract prose across the two declared instruction files and nothing else: commit `2245554f`, 2 files, net +1 line, against a declared surface of ~7 net prose lines at 2× tolerance. The hard self-check did not trip — no Go source, no test, no gate, lint, or CI lane, no third instruction file — and Edit D stayed inside its narrow captain-approved scope, widening when the existing detached audit fires rather than creating a mechanism. Evidence is a rerun mutation exercise and git state, both one-off; the presence greps for the edited strings are recorded here as validation-time evidence and are deliberately not committed, per the ruling Edit C ships. `internal/contractlint` passes, including `TestFOFunctionPromptSurfaceShrinks` — confirming neither edited file sits under the first-officer byte ratchet.

### Review findings (roborev)

- MEDIUM, accepted and fixed — `docs/dev/README.md:76`, Edit C. The appended captain ruling called pasted grep output "legitimate external evidence for that run" without bounding what it may be evidence FOR; a reader landing on the ruling could take it as licence for a grep to satisfy a behavioral AC. Moving a self-derived check out of CI and onto the page changes its permanence, not the independence of its expected value.
  Fixed ADDITIVELY in `014fe501`. The captain's ruling is untouched — verified byte-for-byte, not by eye: extracted the shipped span and `diff`'d it against the approved Edit C text in this entity's `## Documentation changes`, 326 bytes each, identical. The bounding clause follows it as separate prose: the grep is structural or inventory evidence only (present or absent) and can neither satisfy a behavioral acceptance criterion nor serve as the independent source the falsifiable-evidence rule requires, because its expected value still comes from the file under test. It states a boundary the verbatim sentence does not, so no restatement was needed and no captain question arose.
  What would make this fix fail: an edit that makes the licence unbounded again, or any change to the 326-byte ruling span — the `diff` above goes non-empty. Surface unchanged at net +1 line (the clause extends the same line), still 2 files.

### Merge-base caveat for validation (read before running AC-4's check)

CORRECTED after a rebase (see the rebase note under Review findings — an earlier revision of this section named `972129ac`, which is no longer in this branch's ancestry; do not run that command).

Use `git diff $(git merge-base HEAD origin/main) --name-only`. The sprint-setup commits have since landed on `origin/main`, so `git merge-base HEAD origin/main` and `git merge-base HEAD main` now agree at `5dac2d6a` and the literal AC-4 command is correct as written — the earlier false-positive condition is gone.

Against `5dac2d6a`: **2 files** — `docs/dev/README.md`, `skills/ensign/references/ensign-shared-core.md` — 6 insertions, 5 deletions, net +1 line; `--diff-filter=A` × `_test.go` = 0; files matching `\.go$|\.ya?ml$|contractlint|testlint` = 0. AC-4 holds. Prefer the `$(git merge-base ...)` form over any hardcoded SHA: this branch has been rebased once already, and hardcoded SHAs silently rot into wrong answers rather than errors.

### Required CI lanes for this diff

This diff touches `skills/**/references/**` (`ensign-shared-core.md`), which under `docs/dev/README.md:78` makes the host live lanes REQUIRED green before merge — not optional, and a flake there is grounds to re-run to green, never to skip. Captain waiver in force: the **pi** lane's red is WAIVED for this train. That waiver covers **pi only** and must not be stretched over a failure in any other lane; a non-pi red is a real blocker. Flagging here so validation does not rediscover it late.

Local `go test ./...` after the fix: EXIT=0, 15 packages ok, 0 FAIL, including `TestFOFunctionPromptSurfaceShrinks` (the FO byte ratchet) — neither edited file is under that cap and both stay out of `first-officer/**`.

### Review findings (roborev, re-review)

Earlier medium (unbounded grep licence) closed — the fix drew no finding. Four new items: two accepted and fixed, two declined as out-of-surface with grounds and promote conditions.

- ACCEPT — audit-trigger scope conflict, `docs/dev/README.md`. Edit D's provenance trigger read ambiguously against the validation bullet's flat "Routine, low-blast-radius changes do not need it": a reader could not tell whether the provenance trigger fires on every change or only on the four high-stakes surfaces — the difference between "sometimes" and "every task".
  Fixed in `b7ec01c6`. Both lines now say the same thing: the two triggers are independent; the four-surface trigger runs the full audit on a throwaway checkout, the provenance trigger fires wherever such an AC appears — including on a change routine enough to skip the full audit — and covers that AC's adversarial-edit check alone, not the whole change. This is the widening the captain approved (WHEN the existing audit fires), conditioned on the provenance tell, so it is not an always-on audit for all changes and did not require a stop-and-ask.
  Deviation, flagged: this took TWO short clauses, not the one requested. Fixing only the Proof-policy bullet would have left the validation bullet's flat exemption still contradicting it, since that sentence asserts the exemption outright rather than deferring. Both are in this entity's own file; no third file, no widening.
  What would make it fail: any reading under which a routine change with a provenance AC skips the check, or under which a change without such an AC now owes a full audit. Neither is derivable from the two lines as they stand.
- ACCEPT — report-size precedence, `skills/ensign/references/ensign-shared-core.md`. The 30-50 line cap and the new per-cited-test assert-and-falsifier obligation are unsatisfiable together for an ensign citing many tests.
  Fixed in `b7ec01c6`, one clause: falsifiability wins, satisfied by grouping tests by the claim they prove and giving one such line per claim rather than dropping the falsifying change to fit. Precedence note, not a new protocol — it resolves the collision without adding steps.
- DECLINE (correct finding, wrong owner) — `skills/commission/SKILL.md` and `skills/commission/references/templates/development.md` still present grep as first-class proof and omit the falsifying-change clause, so ordinary commission runs scaffold ACs that contradict the policy this entity ships. The observation is CORRECT and it matters. It is not mine to fix: those two files are the declared surface of sibling `template-rigor-propagation` (2ae), which exists to carry this sprint's rigor into the commission templates and lands last precisely so it copies the settled wording landed here rather than an ideation-time draft. Touching them would be a third and fourth instruction file, tripping this entity's hard self-check, and would duplicate-deliver another member's surface. Routing note: carried to 2ae's dispatch by the FO; recorded here only so the trail is complete.
  Promote condition: if 2ae does not land this train, the contradiction between shipped policy and shipped scaffolding becomes live for anyone commissioning a new workflow, and needs its own entity.
- DECLINE (real, out of surface) — `docs/site/contributing/proof-policy.md:36` says "Prose-grep is banned", which now conflicts with the bounded one-off structural allowance shipped here. Real, contributor-facing, and worth fixing — but `docs/site/**` is a third instruction file and trips the hard self-check, so it is filed, not fixed. Next-train candidate.
  Promote to material if: a contributor acts on the stale absolute wording — either withholding legitimate one-off validation evidence, or reading the contradiction as licence to ignore the bullet entirely.
- DECLINE (outright, must not be implemented) — the suggested cross-file consistency test. A new committed check needs explicit captain approval and normally its own entity, and a test asserting that two instruction files agree in wording is prose-over-prose: the banned shape, and precisely what AC-4 promises this entity does not ship. Recorded so a later reader sees it was considered and rejected on grounds, not overlooked.

### Rebase note (affects every SHA cited above this line)

This branch was rebased between the first correction and this one. My commits were rewritten: `2245554f`→`b9197ac9` and `014fe501`→`3b3a8fa3`, and the branch point moved from `972129ac` to `5dac2d6a`. Any SHA cited in the earlier report sections is stale. Current branch: `b9197ac9`, `3b3a8fa3`, `b7ec01c6`. The Merge-base caveat section above has been corrected in place because it named a command (`git diff 972129ac`) that now returns a wrong answer rather than an error — actively misleading, so left uncorrected it would have sent a validator to a false conclusion.

`go test ./...` after both accepts: EXIT=0, 15 packages ok, 0 FAIL, ratchet `TestFOFunctionPromptSurfaceShrinks` included. Captain's Edit C ruling re-verified byte-for-byte after these edits: 326 bytes, `diff` against the approved text empty.

### Bounding clause revised on captain clarification (supersedes the first fix)

The bounding clause shipped in `3b3a8fa3` said a one-off grep "cannot satisfy a behavioral acceptance criterion." That wording is withdrawn. It came from the roborev finding via an FO instruction that was later corrected as too restrictive — recorded plainly here rather than smoothed over, because the trail should show the correction.

Two things were wrong with it. It banned legitimate evidence: a grep run at validation time IS sound evidence that something exists or is absent, and a validator can reproduce it. And it invited the wrong workaround — relabelling a claim "non-behavioral" to make a weak grep admissible, which is the failure the rule exists to prevent, reached by a different road.

Revised in `178b8c4d`. The boundary is now honesty of evidence rather than category: presence or absence is an existence fact and a grep establishes it soundly when that fact is itself the claim; when the claim is about what a program or an agent does, the words being present says nothing about that, so the grep is misleading and the claim must be re-expressed in a form that can be exercised rather than evidenced by a grep that cannot bear it. Committing any such grep as a test remains banned — that half is carried by the captain's verbatim ruling immediately preceding the clause, not restated.

Verified: the string "behavioral acceptance criterion" no longer appears anywhere in `docs/dev/README.md` (count 0), and the captain's Edit C ruling is still byte-identical at 326 bytes (`diff` against the approved text, empty). What would make this fix fail: a reading under which an existence-fact grep is inadmissible, or one under which a grep may stand as proof of runtime behavior. Neither survives the clause as written.

`go test ./...`: EXIT=0, 15 packages ok, 0 FAIL, ratchet `TestFOFunctionPromptSurfaceShrinks` included. Surface unchanged against `5dac2d6a`: 2 files, 6 insertions, 5 deletions, net +1 line. Branch: `b9197ac9`, `3b3a8fa3`, `b7ec01c6`, `178b8c4d`. The Edit D audit-scope reconciliation and both declines are unchanged by this revision.

## Stage Report: validation

- DONE: Each AC verified with evidence I REPRODUCE. AC-4's negative deliverable is a git-state check I ran myself, via `git diff $(git merge-base HEAD origin/main)` — not a hardcoded SHA.
  Merge base resolves to `bdf39f01` (NOT the `5dac2d6a` the implementation report names — the branch was rebased again; the `$(git merge-base ...)` form is what kept AC-4 answerable). Against it: exactly 2 files (`docs/dev/README.md`, `skills/ensign/references/ensign-shared-core.md`), 6 insertions, 5 deletions, net +1; `--diff-filter=A` = empty; `\.go$|\.ya?ml$|contractlint|testlint` = 0 matches. AC-1 re-run below; AC-2 read against an independent source (the sprint index, not this body); AC-3 by removal-check plus the falsifiability anchor.
- DONE: The VALUE exercise is re-run, not cited — throwaway mirror/control pair built from scratch, production function mutated, divergence observed.
  Throwaway module outside the repo. Baseline: both PASS. Mutation `"hello"`→`"HELLO"`: `TestGreetingMirror` **PASS** (the hole), `TestGreetingControl` **FAIL** `greet_test.go:14: got "HELLO" want "hello"`. What this asserts and what would flip it: it asserts the mutation reaches the control and not the mirror; it fails if the mirror reds or the control stays green — neither did, across a second value mutation (`""`) as well. Line cite differs from both prior runs (`:14` vs `:16`/`:17`) because the file is regenerated each time; the divergence, not the line, is the claim.
- DONE: Edit A's real payload — under the new report protocol the two tests are DISTINGUISHABLE where "5/5 passed" rendered them identically, and the mirror cannot name a falsifying change.
  Old protocol: both green tests render `2/2 passed` — one string, zero assertion content, the tautology invisible to the gate. New protocol: control → "asserts `Greeting()` returns the literal `hello`; mutating the production function's return value reds it" (verified RED twice); mirror → "asserts `Greeting()` equals `Greeting()`; no change to the returned VALUE reds it" (verified GREEN under both value mutations). The mirror's inability to name a value-level falsifier is what the gate now sees.
- DONE: The captain's Edit C ruling is byte-for-byte intact and the surrounding text does not contradict it; the four Proof-policy bullets and the validation stage-def are mutually coherent.
  Extracted the shipped span from `docs/dev/README.md:76` and compared it against the ideation-approved Edit C text in this entity's `## Documentation changes` — IDENTICAL. Byte accounting: 325 bytes for the ruling sentence, 326 including the trailing separator space before the bounding clause (the report's "326" counts the separator; the text is identical either way). The shipped boundary is honesty, not category: it says presence/absence is an existence fact a grep establishes soundly, and that where the claim is about what a program or agent DOES the claim must be re-expressed in exercisable form. It does NOT ban existence/absence evidence, and it keys on what the evidence can establish rather than on how the claim is labelled, so relabelling a claim buys nothing. Committed-grep ban intact via the verbatim ruling.
  Edit D scope: I read the two lines as INDEPENDENT triggers, not always-on. `README.md:78` says the provenance trigger "fires wherever such an AC appears ... and covers only that AC's adversarial-edit check, not the whole change"; the validation bullet (`:157`) carries the matching exception. "Such an AC" has a clear antecedent (same-package provenance), so no reading makes a change WITHOUT such an AC owe an audit — the widening stays inside the captain's approved scope.
- DONE: Semantic adversarial pass and finding classification.
  Adjacent-variant matrix on the mechanism: 3 mutations × 2 tests. Value mutations (`"HELLO"`, `""`) → invariant holds (mirror green, control red). The ONE mutation that reds the mirror is making `Greeting()` non-deterministic (stateful counter: `got "hello2" want "hello1"`) — not a value change, so the no-gate premise survives, but see polish 1. Fleet-wide sweep: "5/5 passed" survives nowhere in `skills/` (only the sprint ledger records its removal); `agents/ensign.md` carries no duplicate size guideline to diverge from. AC-2's three grounds trace to `docs/roadmap/0260-proportionality/index.md:17` and `:48` — an independent document, so the decision record is not self-referential. Both declines are factually accurate: commission templates still read `Verified by: {grep / test name / ...}` with 0 occurrences of a falsifying-change clause, and `docs/site/contributing/proof-policy.md:36` does carry the stale absolute. The declined cross-file consistency test was not implemented (AC-4's added-file check is empty).
- DONE: Required-lane mapping derived and stated.
  `skills/ensign/references/ensign-shared-core.md` matches `skills/**/references/**` and is loaded unconditionally by `skills/ensign/SKILL.md:8` for EVERY host, so it is the host-neutral ensign contract, not one adapter → under `README.md:79` every host lane is REQUIRED: `claude-live`, `codex-live`, `pi-live` (pi red WAIVED by the captain, pi only). Status: the branch is not pushed, there is no PR, and `gh run list` returns zero runs for it — the lanes are UNRUN, not green. This is a merge precondition, not a deliverable defect, and it must not be discharged by analogy with the sibling that merged on deterministic lanes.

### Summary

**PASSED — no material findings.** All four ACs verified against evidence I reproduced rather than cited: AC-1's mutation exercise rebuilt from scratch and extended to a 3-mutation matrix, AC-4's git state recomputed against a merge base that had moved again, AC-3's distinguishability grounded in the observed green/red split, AC-2's grounds traced to the sprint index. Edit C is byte-identical and the honesty framing is the one the captain corrected to — legitimate existence evidence preserved, committed greps still banned, no relabelling escape. Edit D reads as a conditioned widening, not an always-on audit. **Merge is blocked on one outstanding precondition:** the host live lanes are required green for this diff and have not run.

Evidence honesty (per the rule this entity ships): EXERCISED — the mutation matrix, `go test ./...` (EXIT=0, 15 test packages, 0 FAIL, ratchet included), the git-state check, and the Edit C span `diff` (shipped file vs an approved text in a different file — two values that can diverge). EXISTENCE FACT ONLY — the string presence/absence sweeps and the two decline confirmations; they establish that text is or is not there and nothing about behavior. JUDGMENT, no machine check — the four-bullet coherence read and the Edit D scope reading. Grepping this deliverable for the strings we added would prove only that we wrote them, so no such grep is offered as behavioral proof.

**Deferred risks** (none block; each with trigger and promote condition):
1. Edit A ships fleet-wide but WITHOUT the independence qualifier that Proof-policy bullet 77 carries, so read alone it admits a degenerate falsifier ("delete the test"). Trigger: an ensign in a workflow other than `docs/dev`, whose README lacks bullet 77. Supported path holds — in `docs/dev` the two clauses compose and bullet 77 imports "an independent source that can diverge". Promote when: a gate accepts a degenerate falsifier and a tautological test ships. Owner already exists (`template-rigor-propagation`).
2. The honesty clause's existence-fact sentence, quoted in isolation, could be offered to defend an AC whose whole claim is "the file contains X". Mitigated two sentences earlier in the same bullet (the tautology ban) and by bullet 77. Promote when: such an AC lands and passes a gate.
3. The provenance trigger does not restate "on a throwaway checkout, never the implementation worktree"; the natural reading inherits it from the same bullet. Promote when: a provenance adversarial edit is run in an implementation worktree and leaves residue.
4. Both recorded declines confirmed accurate and left as filed, with their existing promote conditions unchanged.

**Polish** (no action required this train):
1. The implementation report's "no change to the production function makes it fail" is strictly overstated — a non-determinism change does red the mirror. Accurate form: no change to the VALUE returned reds it. Worth correcting because naming falsifiers precisely is this entity's own thesis.
2. Every SHA in the report is stale again after a further rebase: current commits are `a2d4791b`, `4bacd95c`, `48addee2`, `d39e4306` and the base is `bdf39f01`, not `5dac2d6a`.
3. The report cites `docs/dev/README.md:78` for the required-lanes rule; Edit B's inserted bullet pushed that rule to line 79 (line 78 is now the detached-audit bullet).
4. The `docs/site` decline slightly overstates its conflict: that page's "Prose-grep is banned:" is immediately scoped by "a test asserting that ...", which already matches Edit C's committed-test ban. Still worth filing; less urgent than recorded.
