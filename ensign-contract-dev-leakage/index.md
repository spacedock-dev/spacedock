---
id: ep0ra3zjf4hhkhx5rrkwsxbb
title: Universal ensign contract has absorbed dev-workflow assumptions (TDD, code-only deliverables, "CODE only" worktree) — re-home dev policy out of the shared core
status: ideation
source: session-10 detached audit (Task waymqcmru, 2026-06-03) — overall verdict MATERIAL-PRESENT; validated the captain's grilling instinct ("why am I seeing dev-workflow specific ones like TDD in the universal contract?")
score: "0.19"
worktree:
started: 2026-06-03T07:09:59Z
completed:
verdict:
issue:
---

The shared ensign contract (`skills/ensign/references/ensign-shared-core.md`) is loaded by **every** dispatched ensign across all three commission templates (development, experiment, refinement). A detached adversarial audit confirmed it has absorbed dev-workflow-specific discipline — TDD authoring, a code-shaped closed deliverable list, and a "CODE only" worktree assumption — and asserts them as universal. The FO contract was already corrected to deliverable-shape-agnostic phrasing (`first-officer-shared-core.md` lines 342-353); the ensign contract is the **uncorrected twin**.

**Scope (this entity owns TWO re-homings of the same dev-leakage class).** A 2026-06-03 captain decision on entity 2a (`require-external-proof-guard`, PR #277) merges 2a as-is and folds the generic-vs-dev split of its external-proof guard into THIS entity. So ep now owns:

- **(a) Prose layer — re-home the ensign shared core's dev discipline.** Lift the `## Working Practices` block's two dev-specific principles (TDD authoring, code-shaped deliverable list) out of the universal core into the dev layer; neutralize the "CODE only" worktree noun and the "worktree path" universal-field enumeration. (The original ep scope.)
- **(b) Code layer — split 2a's external-proof guard along the same generic/dev seam.** 2a's `internal/status/external_proof.go` conflates a GENERIC self-reference kernel (`ClassifyEntityACs`: an AC's proof cannot be re-reading its own write-up) with a DEV-specific external-token VOCABULARY (`externalTokenRe`: external == `.go`/`gofmt`/`goreleaser`/`--cask`/`commit`/CI/…). The generic kernel is currently **only satisfiable through** the dev vocabulary, fenced behind the `require-external-proof: true` opt-in. ep designs the clean split: the generic kernel stays workflow-agnostic; the dev vocabulary is supplied by a dev-specific layer, so a non-dev workflow could supply its own proof vocabulary without inheriting Go-toolchain tokens.

**Governing principle (captain, this session).** *"An AC's proof must be external and able to fail"* is GENERIC — true for any workflow (research: a metric; content: a published artifact; ops: a runbook run). *"External means code/command/state/CI/release"* is DEV-specific. 2a's guard conflates them: the generic self-reference kernel is reachable only through a dev-flavored `externalTokenRe`. The project's PROSE layer already tags this distinction — `docs/dev/README.md:110` ends with "(This is dev-workflow policy: an AC's proof here is code/command/state. A non-development workflow's AC proof may legitimately be a published artifact, a metric, or a human review.)" — and the generic self-reference kernel itself already lives in the universal FO core (`first-officer-shared-core.md:346`: "The gate's AC cross-check refuses a criterion whose only proof is review of the entity's own prose."). The CODE layer (`external_proof.go`) is the lone place the two are still welded together. ep's job is to make the code seam match the seam the prose already draws.

## Problem

`ensign-shared-core.md` lines 22-27 carry a `## Working Practices` block that duplicates dev-workflow policy verbatim into the portable contract. The project's own docs already tag this same policy as dev-only opt-in:
- `docs/dev/README.md` line ~110: "This is dev-workflow policy: an AC's proof here is code/command/state. A non-development workflow's AC proof may legitimately be a published artifact, a metric, or a human review."
- `skills/commission/references/templates/development.md` lines ~111-117: a "Recommended practices (opt-in)" carve-out stating "the universal first-officer contract does not impose them."

The shared ensign core strips that qualifier and presents the rules as universal ensign discipline, breaking at least two non-dev template classes: refinement entities whose deliverable **is** prose (PRD, outreach reply, content — `refinement.md` says they never touch the repo), and experiment entities whose deliverable is a hypothesis/analysis verdict against pre-registered success criteria.

**Material findings from the audit (Task waymqcmru):**
- **L2-01 (umbrella):** The entire `## Working Practices` section duplicates dev policy into the universal contract. Lift it out; re-home each principle.
- **L1-F1 / L2-02:** TDD bullet (line 24, "Write the failing test first… The test is what the gate judges") asserts test-first as the universal authoring contract — false for refinement (gate judges the artifact body) and experiment (gate judges evidence vs. pre-registered criteria) workflows.
- **L2-03:** "Every task produces a real, checkable change" + its escape clause under-serve prose-only deliverables (refinement/experiment bodies that are checkable via pre-fixed success criteria, not external files).
- **L2-05:** The Split-Root State Contract (line ~35) hard-codes "the worktree isolates **CODE only**" — excludes non-code work product (experiment evidence files, refinement attachments). Same fix needed in **both** `ensign-shared-core.md` and `first-officer-shared-core.md`.
- **L3-1:** Structural twin of the FO's already-corrected four-principle block — the ensign block is the uncorrected version.
- **L4-P2 / runtime adapters (Polish):** `claude-ensign-runtime.md` and `codex-ensign-runtime.md` enumerate "worktree path" as a top-level always-present assignment field, while the shared core already uses the substrate-neutral "workflow location" for the same slot.

## Proposed approach

Mirror the FO-contract correction (`first-officer-shared-core.md` Working Principles, lines 342-353) onto the ensign contract. The exact edits are pinned below against verified current line numbers (ground-truthed during ideation, 2026-06-03).

### Re-homing design — what lifts, and where each principle lands

The `## Working Practices` block is `ensign-shared-core.md:22-27` (four bullets). Re-home, do not delete:

1. **TDD authoring bullet (line 24, "Write the failing test first… The test is what the gate judges").** This is the only principle with **no current dev home**. Lift it into `development.md`'s existing `## Recommended practices (opt-in)` section (line 111) as a new `### Test-first authoring` subsection, phrased identically to how that section already frames the validation-stage disciplines ("recommended, not mandatory — the universal contract does not impose them"). `docs/dev/README.md`'s ideation stage already carries the code-gate-over-prose corollary (line 85); the TDD authoring discipline itself is what's missing and must land in `development.md`.

2. **"Every task produces a real, checkable change" bullet (line 25).** Already homed: `docs/dev/README.md:80` carries this verbatim-equivalent as a dev ideation-stage rule. The lift is a **deletion from the shared core** (its dev home already exists) — no new home to author. Confirm during implementation that the README clause still covers the prose-only escape ("belongs in the roadmap").

3. **"Prove by exercising, not by re-reading" + "No hidden machine dependencies" bullets (lines 26-27).** These are **genuinely universal** (they hold for experiment evidence and refinement artifacts too) — they are not dev-leakage. They stay, but move out of a block named "Working Practices" (a dev-flavored heading) into the existing universal flow. Fold them into the `## Working` numbered list or a neutral `## Proving your work` heading. The implementer decides placement; the invariant is: these two survive in the universal core, the two dev-specific ones do not.

4. **Optional universal restatement.** The shared core MAY retain one discipline-neutral sentence so a non-dev ensign still knows it owes proof: *"The stage definition states the proof your stage owes — a test, a metric, a published artifact, a human review. Satisfy that, not a generic test-first ritual."* This covers TDD (dev), pre-registered criteria (experiment), and acceptance criteria (refinement) without asserting any one tradition. This mirrors the FO core's already-neutral "satisfy the stage definition's proof requirements" framing — adopt the same shape, do not re-invent the four-bullet block.

`experiment.md` / `refinement.md` need **no new homes** — confirmed during ideation: `experiment.md` already encodes "success criteria fixed before evidence is gathered" (hypothesis stage), and `refinement.md` carries draft/review/polish. The re-homing targets are `development.md` (new TDD subsection) and the shared core's own structure (delete two, keep two, optionally add one neutral line).

### "CODE only" → substrate-neutral (AC-3)

- `ensign-shared-core.md:35` — "the worktree isolates **CODE only**" → "the worktree isolates **the deliverable work product only**".
- `first-officer-shared-core.md:270` — "a worktree stage isolates **CODE only**" and "The worktree still owns code: … apply to **code changes** only" → "isolates **the deliverable work product only**" / "apply to **deliverable-artifact changes** only". Keep the rest of the clause (pr-mirror exception, state-checkout path) intact — only the substrate noun changes.

### "worktree path" → neutral location vocabulary (AC-4)

- `claude-ensign-runtime.md:7` and `codex-ensign-runtime.md:7` enumerate "entity, stage, stage definition, **worktree path**, and checklist" as universal assignment fields. Rename the slot to the shared core's neutral term: **"workflow location"** (the term `ensign-shared-core.md:11` already uses). A non-worktree stage (ideation, backlog) has no worktree path but always has a workflow location.
- **In-scope nuance found during ideation:** `ensign-shared-core.md:17` ALSO says "If you were given a **worktree path**…". This usage is **conditional** ("if you were given"), not an asserted-universal field, so it is legitimate and does NOT need changing. The AC-4 oracle must scope to the *universal-field-enumeration* position (the runtime adapters' "fields are: …" lists), not ban the literal "worktree path" everywhere — otherwise it false-fails on the legitimate conditional at line 17.

The shared ensign core should retain only universals: read the assignment, follow the stage definition's proof requirements, commit the entity to the state checkout, signal completion.

### Code layer — the 2a generic/dev split (AC-6)

This is the NEW scope folded in per the captain's 2a decision. Map of `internal/status/external_proof.go` (2a's worktree `.worktrees/spacedock-ensign-require-external-proof-guard`, PR #277, HEAD `128041a9`) along the generic/dev seam:

**GENERIC (workflow-agnostic — stays in the universal kernel):**
- `ClassifyEntityACs(body string) []ACFlag` — the 5-step algorithm (lines 107-165): extract `**AC-N` blocks → isolate proof clause → strip quoted spans → match a self-phrase → **require external-token absence**. Steps 1-4 carry NO dev vocabulary.
- `selfPhraseRes` (lines 47-52) — the self-reference phrases (`this entity's`, `review of … section`, `the entity's own prose/decision/section`, `re-reading … body`). These describe the SELF-REFERENCE antipattern, not a deliverable substrate — generic.
- `acHeaderRe`, `proofMarkerRe`, `quotedSpanRe`, `isolateProofClause`, `matchSelfPhrase`, `stripFrontmatter` — pure structural parsing, no domain vocabulary.

**DEV-specific (the proof VOCABULARY — must move to a dev-supplied layer):**
- `externalTokenRe` (lines 61-71) — the entire alternation is Go/CI/release/commit vocabulary: `\.go`, `gofmt`, `\bvet\b`, `goreleaser`, `\bcask\b`, `spctl`, `notariz`, `frontmatter`, `stdout`, `\bCI\b`, `\bPR\b`, `commit`, `\bGitHub\b`, the 7+ hex SHA shape. This IS "external means code/command/state/CI/release" — the dev reading of "external". A research workflow's external-proof vocabulary (a DOI, a p-value, a dataset hash) shares none of these tokens.

**The seam.** Step 5 ("require external-token absence") is where generic meets dev: the KERNEL owns the *structural rule* ("the proof clause must cite a recognized external token, else flag"); the dev layer owns *which tokens count as external*. The clean split makes `externalTokenRe` an **injected parameter** of the kernel rather than a hard-coded package global:

    // generic kernel — no dev vocabulary
    func ClassifyEntityACs(body string, externalTokens *regexp.Regexp) []ACFlag

The dev caller (handlers.go `runSet`, validate.go `validateWorkflow`) passes the dev vocabulary; `externalTokenRe` moves to a `dev`-named home (a `devExternalTokens` var, or — cleaner — the same dev layer the README opt-in already gates). A non-dev workflow that one day opts in supplies its own token set without inheriting `goreleaser`/`gofmt`.

**Why this is a re-homing, not a behavior change.** The dev workflow is the ONLY current opt-in (`require-external-proof:` defaults OFF, byte-identical to absent). Passing `devExternalTokens` from the dev callers reproduces today's behavior token-for-token — `TestClassifierPrecisionRecallOnLiveCorpus` must stay green at the same flagged set (`{external-tracker-checkpoint/index.md AC-6}`). The split is a parameterization that leaves the dev behavior identical while making the kernel reusable. This honors the captain's "merge 2a as-is" — 2a ships welded; ep performs the seam cut on top.

**Re-conflation guard.** The lock-in is a test that goes RED if the dev vocabulary leaks back into the generic kernel: assert `external_proof.go`'s `ClassifyEntityACs` signature takes the token regex as a parameter (the kernel does not name a package-global `externalTokenRe`), AND assert no Go-toolchain literal (`gofmt`, `goreleaser`, `\.go\b`) appears in the kernel function body. If a future edit re-hard-codes the dev tokens inside the kernel, the guard fails. (AC-6 below pins the exact shape.)

## Out of scope

- Rewriting the dev-workflow disciplines themselves — they are correct *for dev*; this is a re-homing, not a policy change.
- Changing the FO contract's already-corrected four-principle block (only its "CODE only" line, shared with the ensign core, is in scope — finding L2-05).
- The experiment/refinement template bodies beyond confirming they already carry the re-homed analogues.
- **Changing the dev workflow's external-proof BEHAVIOR.** The 2a generic/dev split (AC-6) is a re-homing/parameterization: `TestClassifierPrecisionRecallOnLiveCorpus` must stay green at the identical flagged set. ep does NOT add, remove, or re-tune any `externalTokenRe` vocabulary — that is 2a's settled spec. ep only relocates WHERE the dev vocabulary is supplied (caller-injected) vs. WHERE the generic kernel lives.
- **Re-litigating 2a's opt-in design.** `require-external-proof:` README key, default-OFF, terminal-only gating, the dev-scoped opt-in — all settled in 2a and merged as-is. ep inherits them.
- **Building a non-dev external-proof vocabulary.** ep makes the kernel *accept* an injected token set; it does NOT author a research/content token set (YAGNI — no non-dev workflow opts in today). The split is proven sufficient by the dev caller reproducing today's behavior, not by shipping a second vocabulary.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action. Proof is presence/absence over the instruction files (legitimate per the README: "a presence check over instruction files proving they carry a required clause or stay free of a banned token is proof at the claim's own level") plus the lock-in oracle. Ideation refines.

**AC-1 — The universal ensign shared core no longer asserts dev-only authoring discipline.**
Verified by: a Go oracle test (extending the proven `internal/hostneutrality/prose_inflator_locks_test.go` pattern) asserting `ensign-shared-core.md` is free of the dev-only test-first vocabulary — the exact banned literals are "failing test", "feature or bugfix", and "the test is what the gate judges". The test fails if the Working Practices TDD bullet is reintroduced. Negative-proof shape: the implementer inserts the banned phrase, watches the test go red, removes it, watches green (the same lock-in-is-real demonstration `TestNoAuditTrailExposition` uses).

**AC-2 — The dev disciplines survive in their dev-specific homes.**
Verified by: the same (or paired) oracle asserting the *positive* presence of the re-homed guidance — `development.md` carries a "test-first" clause in its `## Recommended practices (opt-in)` section, AND `docs/dev/README.md` retains the "real, checkable change" / deliverable-proof policy. This proves the lift *relocated* rather than *deleted* the guidance; the test fails if a future edit strips the dev home.

**AC-3 — The worktree-isolation boundary is substrate-neutral in both shared cores.**
Verified by: the oracle confirming neither `ensign-shared-core.md` nor `first-officer-shared-core.md` contains the literal "CODE only" (case-sensitive, as it appears today), and that both still carry a worktree-isolation clause (presence of the isolation concept, absence of the code-substrate noun — so the fix neutralizes rather than deletes the boundary).

**AC-4 — Both runtime adapters use the shared core's neutral location vocabulary.**
Verified by: the oracle confirming the `## Agent Surface` / `## Dispatch` assignment-field enumeration in `claude-ensign-runtime.md` and `codex-ensign-runtime.md` lists "workflow location" and NOT "worktree path". Scoped to the field-enumeration sentence ("authoritative for all assignment fields: …"), NOT a blanket ban on the literal "worktree path" — the conditional usage at `ensign-shared-core.md:17` ("if you were given a worktree path") is legitimate and must remain passing.

**AC-5 — The contract still drives a real workflow correctly.**
Verified by: a live ensign dispatch on the swept ensign contract completing a stage cleanly (the strongest signal that load-bearing meaning was preserved — same bar Phase 0.A met by driving on both sonnet and opus).

**AC-6 — The external-proof guard's generic self-reference kernel is separated from the dev-specific proof vocabulary.**
Verified by TWO Go tests in `internal/status/external_proof_test.go` (the package 2a's classifier ships in):
- **(a) Structural-seam test `TestClassifierKernelTakesInjectedTokens`** — asserts `ClassifyEntityACs`'s signature takes the external-token regex as a parameter (the kernel is callable with a *non-dev* token set), and that the dev callers (`runSet`, `validateWorkflow`) supply `devExternalTokens`. This is a real relationship over parsed values — the kernel function accepts the injected matcher — not a substring search. The dev behavior is unchanged: a separate assertion drives the kernel with `devExternalTokens` over the same fixtures 2a's `TestExternalTokensClearSelfPhrase` uses and gets the identical accept/flag verdicts.
- **(b) Re-conflation guard `TestKernelFreeOfDevVocabulary`** — asserts the kernel function body (`ClassifyEntityACs` and the generic helpers it calls: `isolateProofClause`, `matchSelfPhrase`, `selfPhraseRes`) contains NO Go-toolchain literal (`gofmt`, `goreleaser`, `\bvet\b`, `\.go\b`, `\bcask\b`). Negative-proof shape: the implementer re-hard-codes a dev token inside the kernel, watches this test go RED, removes it, watches green — proving it is a real lock, not a tautology. The dev vocabulary lives ONLY in the `devExternalTokens` var, which this test confirms is the sole site naming those literals.
- **Behavior-preservation backstop:** 2a's `TestClassifierPrecisionRecallOnLiveCorpus` must stay green at the identical flagged set (`{external-tracker-checkpoint/index.md AC-6}`) after the split — the live-corpus invariant proves the parameterization changed structure, not behavior. If the split altered any verdict, this pre-existing test (not authored by ep) goes red on the real corpus.

## Lock-in oracle (the negative-proof mechanism)

The proven pattern lives at `internal/hostneutrality/prose_inflator_locks_test.go` and ALREADY enumerates `ensign-shared-core.md` and both runtime adapters in its `contractProseFiles` / `runtimeAdapterPaths` tables. The lock-in oracle for this entity extends that file (same package, same table-driven shape) rather than introducing a new test harness:

- **Negative half (AC-1, AC-3, AC-4) — dev vocabulary absent from the universal core.** A `devLeakageLiterals` table of phrases that must NOT appear in `ensign-shared-core.md` / `first-officer-shared-core.md`: `"failing test"`, `"feature or bugfix"`, `"the test is what the gate judges"` (AC-1, ensign core only), `"CODE only"` (AC-3, both cores). Mirrors `auditTrailLiterals` exactly — a `strings.Contains` loop that `t.Errorf`s on any hit. The AC-4 worktree-field check is a scoped regex over the runtime adapters' field-enumeration sentence (anchor on "authoritative for all assignment fields:"), asserting "workflow location" present and "worktree path" absent **within that sentence** — NOT a file-wide ban (the `ensign-shared-core.md:17` conditional must stay green).

- **Positive half (AC-2) — dev vocabulary present in the dev homes.** A `devHomePresence` table mapping each dev home to a required clause: `development.md` must contain a test-first clause inside `## Recommended practices (opt-in)`; `docs/dev/README.md` must contain "real, checkable change". A `strings.Contains` that `t.Errorf`s on *absence* — the inverse polarity, proving the lift relocated rather than deleted.

This is a presence/absence property of the text where the text IS the claim (legitimate per the README: "a presence check over instruction files proving they carry a required clause or stay free of a banned token is proof at the claim's own level"). The negative-proof discipline (insert banned phrase → red, remove → green) is what makes it a real lock rather than a substring spelling check.

**Code-layer oracle (AC-6) — a DIFFERENT package, a DIFFERENT proof level.** AC-1..AC-4 are prose invariants in `internal/hostneutrality/` (the text IS the claim → presence check is proof at the claim's own level). AC-6 is a CODE-structure claim, so its proof is stronger than presence: the seam test (a) DRIVES the parameterized kernel with a non-dev token set and observes the verdicts — exercising behavior, not asserting spelling; the re-conflation guard (b) is a presence/absence check over the *kernel function body* (legitimate because "the kernel is free of dev vocabulary" is itself a property of that text), backstopped by 2a's live-corpus behavioral test which would go red if the structure change altered any verdict. AC-6's tests live in `internal/status/external_proof_test.go` (2a's package), NOT the `hostneutrality` prose oracle — the two layers stay in their own packages.

## Spike determination

**No spike needed.** This change composes only already-proven mechanisms:
- The table-driven prose-invariant oracle is proven in-repo (`prose_inflator_locks_test.go`, shipped 0.19.4 / Phase 0.A AC-4).
- The live-ensign-dispatch path (AC-5) is the same dispatch mechanism 0.19.4 already exercised on both sonnet and opus; re-running it on the swept contract is a regression check, not an unverified handoff.
- No new on-disk format, parser round-trip, or runtime handoff is introduced — the prose edits are re-homing across files the oracle already reads.
- **AC-6 (code split):** parameterizing a pure Go function (`ClassifyEntityACs` gains a `*regexp.Regexp` argument; callers pass `devExternalTokens`) is a textbook refactor, not an unverified mechanism. The proof it composes — a Go test driving the function + a live-corpus invariant — is the same `external_proof_test.go` pattern 2a already proved across two feedback cycles (820→853 tests green). The only judgement call is the seam location (step 5: kernel owns "external-token absence is a flag", dev owns "which tokens count"), recorded in the design above and settled at the gate by review.

**Hard dependency: AC-6 cannot land until 2a merges.** `internal/status/external_proof.go` does NOT exist on `main` — it lives only in 2a's worktree (`.worktrees/spacedock-ensign-require-external-proof-guard`, PR #277). ep's AC-6 edits and tests target that file. ep's implementation stage MUST be sequenced AFTER 2a merges to `main`; until then there is nothing to split. (AC-1..AC-5 — the prose layer — have no such dependency; they edit files already on `main`.) This is a real sequencing constraint, not a spike: the mechanism is proven, the FILE is just not landed yet.

## Test plan

- Go oracle tests for AC-1..AC-4 (extending `internal/hostneutrality/prose_inflator_locks_test.go` — the proven Phase 0.A negative-proof + presence pattern). Cost: low — text invariants over instruction files, same package, no new harness.
- One live ensign dispatch for AC-5. Cost: medium (one live cycle). Drive a real stage on the swept contract and confirm clean completion — the strongest signal load-bearing meaning survived the lift.
- Go tests for AC-6 in `internal/status/external_proof_test.go` (2a's package): the structural-seam test (kernel takes injected tokens; dev callers supply `devExternalTokens`; dev verdicts unchanged over 2a's existing fixtures) + the re-conflation guard (kernel body free of Go-toolchain literals). Cost: low — Go unit tests over the parameterized function, same pattern 2a proved. Plus the behavior-preservation backstop: 2a's `TestClassifierPrecisionRecallOnLiveCorpus` must stay green at the identical flagged set, and full `go test ./...` stays green.
- High-stakes surface (shipped contract/scaffolding + the `status` guard path) → a detached adversarial audit is required before merge per the dev README's validation stage. The auditor should probe (1) the AC-4 scoping (try removing "if you were given" from `ensign-shared-core.md:17` and confirm the field-scoped oracle stays correctly green — no over-reach into the conditional usage), and (2) the AC-6 split (try re-hard-coding a dev token inside the kernel and confirm `TestKernelFreeOfDevVocabulary` goes red; try a kernel verdict-altering edit and confirm 2a's live-corpus test catches it).

### Sequencing and coordination

- **Depends on 2a (PR #277) merging to `main`.** AC-1..AC-5 (prose layer) can start anytime — they edit `main` files. AC-6 (code split) MUST wait for 2a to land, because `external_proof.go` does not exist on `main`. The cleanest order: prose layer first (independent), code split after 2a merges. ep's implementation stage gate should not open the AC-6 work until `git log main -- internal/status/external_proof.go` is non-empty.
- **Shared-file coordination with `at` (`non-interactive-teardown-exit`, id `atwf2w6p…`).** Both ep and at edit `skills/first-officer/references/first-officer-shared-core.md` — but DIFFERENT sections. ep touches ONLY the `## Split-Root Worktree Contract` "CODE only" clause (line ~270). at touches the teardown clauses (`## ` line 148 Supersede-shutdown + line 226 "Teardown agents at terminal"). The sections are disjoint, but the file is shared: ep MUST NOT run concurrent-implementation with at on this file. Sequence the two FO-core edits serially (whoever lands first; the second rebases its one-line clause edit on top — disjoint sections → no conflict, but a concurrent worktree-vs-worktree edit of the same file risks a merge headache). Flag to the FO at dispatch so the two implementation stages are not opened simultaneously on the shared file.

## Notes

- The full audit synthesis (6 material + 2 polish findings, MATERIAL-PRESENT) is Task `waymqcmru` from session 10. Findings are quoted above; re-run the audit if deeper detail is needed.
- This is the ensign-side completion of the contract-simplification arc shipped in 0.19.4 (Phase 0.A swept the FO + ensign cores for *length*; this sweeps the ensign core for *dev-leakage* — a different axis the captain's grilling surfaced).
- Sequence: 0.19.5. **Depends on 2a (`require-external-proof-guard`, PR #277) merging** — the AC-6 code split targets `external_proof.go`, which 2a introduces. The prose layer (AC-1..AC-5) is independent of 2a and of the sonnet-live-CI-flake entity.
- **Captain decision (2026-06-03, this session):** merge 2a as-is + fold the generic-self-ref/dev-vocab split into ep (this re-ideation's AC-6). The generic self-reference kernel is universal; the dev external-token vocabulary is dev-specific; the code seam must match the seam the prose already draws (`docs/dev/README.md:110` + `first-officer-shared-core.md:346`).

## Stage Report: ideation

- DONE: Design the re-homing: what lifts out of ensign-shared-core.md Working Practices and where each principle lands (development.md / docs/dev/README.md), mirroring the already-corrected FO block (first-officer-shared-core.md 342-353).
  Added `### Re-homing design` to Proposed approach: per-bullet disposition (TDD → new development.md "Recommended practices (opt-in)" subsection; "real, checkable change" already homed at docs/dev/README.md:80 → delete-only; "prove by exercising" + "no hidden machine deps" are universal → keep, de-dev the heading; optional one neutral universal line mirroring the FO core). Pinned exact before/after wording for AC-3 ("CODE only", both cores) and AC-4 ("worktree path" → "workflow location", both adapters) against verified current line numbers.
- DONE: Specify the lock-in oracle (dev-only vocabulary absent from the universal core, present in the dev homes) — the proven Phase 0.A negative-proof pattern.
  Added `## Lock-in oracle` + `## Spike determination` sections: the oracle EXTENDS the in-repo `internal/hostneutrality/prose_inflator_locks_test.go` (which already lists these files) with a negative-half banned-literal table (mirrors `auditTrailLiterals`) and a positive-half dev-home-presence table. Recorded "no spike needed" with the proven mechanisms (table-driven prose oracle + already-exercised live-dispatch path).

### Summary

Refined the existing spec rather than rewriting it: ground-truthed every cited line number, then sharpened the two checklist items into an implementable design. Key ideation finding: only the TDD authoring bullet genuinely needs a NEW home (development.md) — the deliverable-shape principle is already at docs/dev/README.md:80 (delete-only), and two of the four Working-Practices bullets ("prove by exercising", "no hidden machine deps") are genuinely universal and must STAY. Flagged an AC-4 over-reach trap: the conditional "if you were given a worktree path" at ensign-shared-core.md:17 is legitimate, so the oracle must scope to the runtime adapters' field-enumeration sentence, not ban the literal file-wide. Recorded "no spike needed" — composes only the proven prose-oracle and live-dispatch mechanisms.

## Stage Report: ideation (re-ideation)

- DONE: Re-anchor ep's scope on the CURRENT state — ep now owns (a) prose re-homing AND (b) the 2a generic/dev factoring; map exactly which lines/sections are dev-leaky vs universal.
  Added a `**Scope (this entity owns TWO re-homings…)**` block + `**Governing principle**` to the Problem statement. Mapped `external_proof.go` (2a worktree, HEAD 128041a9): GENERIC = `ClassifyEntityACs` 5-step kernel (lines 107-165), `selfPhraseRes` (47-52), structural helpers; DEV = `externalTokenRe` (61-71, the Go/CI/release/commit alternation). Ground-truthed that the prose layer ALREADY draws the seam (`docs/dev/README.md:110` dev-policy parenthetical; the generic self-reference kernel already in `first-officer-shared-core.md:346`) — the CODE is the lone welded place.
- DONE: Each AC a REAL gate at the text's own level, PLUS the 2a-factoring AC with a check that fails if generic/dev re-conflate.
  Added `### Code layer — the 2a generic/dev split (AC-6)` design (inject `externalTokens *regexp.Regexp` into the kernel; `externalTokenRe` → caller-supplied `devExternalTokens`) + AC-6 with two Go tests in `internal/status/external_proof_test.go`: (a) structural-seam test (kernel takes injected tokens, dev callers supply them, dev verdicts unchanged) and (b) re-conflation guard `TestKernelFreeOfDevVocabulary` (kernel body free of Go-toolchain literals, negative-proof: insert → red, remove → green), backstopped by 2a's pre-existing live-corpus invariant.
- DONE: Test plan + sequencing — state the 2a dependency, the no-spike determination, and the no-concurrent-implementation-with-at constraint on the shared FO core.
  Added `### Sequencing and coordination`: AC-6 hard-depends on 2a (PR #277) merging — `external_proof.go` is NOT on `main` (verified: file absent on main, present only in `.worktrees/spacedock-ensign-require-external-proof-guard`). Recorded "no spike needed" (parameterizing a pure Go function is a textbook refactor on 2a's proven test pattern). Flagged the `at` (`atwf2w6p…`) shared-file coordination: ep edits ONLY `## Split-Root Worktree Contract` "CODE only" (FO core ~270); at edits the teardown clauses (148, 226) — disjoint sections, but ep MUST NOT run concurrent-implementation with at on the shared file.

### Summary (re-ideation)

Folded the captain's 2a decision into ep: ep now owns BOTH the original prose re-homing (AC-1..AC-5) and the NEW code-layer generic/dev split of 2a's external-proof guard (AC-6). The governing principle — "an AC's proof must be external and able to fail" is GENERIC; "external means code/command/state/CI/release" is DEV — maps cleanly onto `external_proof.go`: the `ClassifyEntityACs` self-reference kernel is generic and stays workflow-agnostic; the `externalTokenRe` Go/CI/release vocabulary is dev-specific and becomes a caller-injected `devExternalTokens`. The seam is step 5 (kernel owns "external-token absence is a flag"; dev owns "which tokens count"). The split is a behavior-preserving parameterization — 2a's `TestClassifierPrecisionRecallOnLiveCorpus` must stay green at the identical flagged set — verified by a structural-seam test plus a re-conflation guard with negative-proof discipline. AC-6 hard-depends on 2a (PR #277) landing because `external_proof.go` is not yet on `main`; the prose layer is independent. Flagged the shared-FO-core coordination with `at` (disjoint sections, no concurrent implementation). No spike needed: a function parameterization on 2a's proven test pattern.
