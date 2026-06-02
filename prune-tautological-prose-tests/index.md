---
id: rx13p1svcchfs8rtdmwsxq1k
title: Prune/demote/replace the tautological prose-pin integration tests (the banned grep-over-prose antipattern)
status: validation
source: "FO classify-integration-tests workflow (2026-06-02) — 6 of 26 skills/integration funcs confirmed tautological-prose-pins (gut-it-keep-green proven by meaning-inversion mutation); captain: prune them"
started: 2026-06-02T17:28:09Z
completed:
verdict:
score: "0.28"
worktree: .worktrees/spacedock-ensign-prune-tautological-prose-tests
issue:
---

A 4-agent classification workflow (mutation-proved each suspect by inverting the asserted prose and confirming the test stayed GREEN) found **6 of 26** `skills/integration` test funcs are the banned **tautological-prose-pin** antipattern (P2: grep-over-prose asserting the contract "says" something — paraphrase-invertible, gives false confidence). Same class as the deleted `v1` grep tests. The other 20 are legitimate (10 absence/structural lints, 8 behavioral). This entity removes the 6.

## Disposition (per the workflow; behavioral ACs)

- **PRUNE** `working_principles_test.go:TestWorkflowGuideCarriesPrinciples` — all four principle clauses inverted in README, every grepped substring kept, stayed green. Prose-quality is doc review, not a Go test.
- **PRUNE** `working_principles_test.go:TestEnsignContractCarriesTestFirstRule` — test-first + no-hidden-deps rules reversed, tokens kept, green (a capitalized-`Hidden` variant only failed on case — a spelling artifact).
- **PRUNE** `ship_local_ceremony_test.go:TestPRMergeFallbackProseIsHonest` — fallback prose rewritten into guard-gaming instructions, stayed green; "honesty of prose" is intrinsically doc review. (If worth keeping any signal: fold its one banned-literal-absence check into the banned-token lint family.)
- **DEMOTE** `working_principles_test.go:TestFOContractCarriesWorkingPrinciplesAndPosture` → a structural lint: assert the `## Working Principles` section/heading EXISTS (a real structural anchor) and drop the four posture-marker substring greps.
- **REPLACE** `skill_text_test.go:TestCommissionStateBackendDecisionRule` → a BEHAVIORAL check: drive the real commission scaffolder/codepath for a standalone input and a PR-bearing-code-repo input, and assert the generated frontmatter actually contains vs omits `state: .spacedock-state`. (Fallback: if no scaffolder codepath is invocable from a test, DEMOTE to a structural decision-table assertion + record why REPLACE wasn't feasible.)
- **REPLACE or PRUNE** `ship_local_ceremony_test.go:TestShipLocalCeremonyBlockExists` — if the no-`--force` terminal-transition property is already covered by a behavioral gate over the launcher/merge codepath, PRUNE; else REPLACE with a test that drives the local-merge transition and asserts no `--force` is emitted.

KEEP (legitimate, do NOT touch): the 10 lint-invariants (`TestShippedInstructionsCarryNoInsiderJargon`, the no-plugin-private-path guards, `portability` host-dependency guards, frontmatter/closure/structure checks) and the 8 behavioral tests (dispatch.Run, launcher binary, real git, manifest version-bracket parses).

## Acceptance criteria (behavioral — no grep self-proof)

**AC-1 — the 4 PRUNE'd funcs are gone; suite stays green.** The four named funcs no longer exist; `go test ./skills/integration/` passes; no dangling reference.

**AC-2 — the DEMOTE'd func pins structure, not prose-meaning, and FAILS on a missing anchor.** `TestFOContractCarriesWorkingPrinciplesAndPosture` (or its replacement) asserts the `## Working Principles` heading exists; verified by mutation — deleting the heading makes it FAIL, while a meaning-inverting paraphrase that keeps the heading no longer needs to pass/fail on wording.

**AC-3 — the REPLACE'd decision-rule test exercises real behavior.** `TestCommissionStateBackendDecisionRule`'s replacement drives the real commission codepath (or is honestly DEMOTED with recorded rationale); verified by mutation — flipping the scaffolder's backend decision makes it FAIL. The ship-local one is PRUNED-or-REPLACED per the same rule.

**AC-4 — no new tautological prose-pin introduced.** The replacements are behavioral/structural, not new grep-over-prose; `go test ./...` green, gofmt/vet clean.

## Notes
- Touches `skills/integration/*_test.go` only (test files) — NOT the `internal/status` serialized lane; parallel-safe with 02/zj.
- Grounded in the v1 saga: a denylist of inversion phrasings cannot close a polarity hole against open-ended paraphrase (the v1 cycle-1 re-audit proved this 5× over), so these don't get "hardened" — they get pruned/demoted/replaced.

## Stage Report: implementation

- DONE: PRUNE (delete) the confirmed pure prose-pins: TestWorkflowGuideCarriesPrinciples, TestEnsignContractCarriesTestFirstRule, TestPRMergeFallbackProseIsHonest, and the no-`--force` half of TestShipLocalCeremonyBlockExists
  Three funcs deleted; no-force half PRUNED (its property is behaviorally covered — see below). Commit af070c5a. Suite green; no dangling refs (grep over repo: zero hits for all four deleted names); 10 lint-invariants + 8 behavioral tests untouched (confirmed TestShippedInstructionsCarryNoInsiderJargon, TestNoPluginStatusPathInVendoredSkills, TestNoPRMergeOrModBehaviorIntroduced + dispatch/launcher/manifest files present).
- DONE: the no-force half PRUNED because the no-`--force` terminal/local-merge transition property is already behaviorally covered
  `internal/status/merge_policy_guard_test.go` drives the real native+oracle status binary: TestMergeLocalNoSentinelTerminalSetSucceeds / ...ArchiveSucceeds / TestSentinelSatisfiesGuard{TerminalSet,Archive} / TestMergeLocalEntityActuallyAdvances all assert the terminal transition succeeds with NO --force (5/5 green). The block-exists + `merge: local` + `local-merge:` sentinel checks stayed as legitimate structural anchors.
- DONE: DEMOTE TestFOContractCarriesWorkingPrinciplesAndPosture to a structural lint: assert the `## Working Principles` heading EXISTS, drop the four posture-marker substring greps
  Renamed to TestFOContractCarriesWorkingPrinciplesSection; asserts only the heading. Commit af070c5a. Mutation-verified: renaming the heading → FAIL ("missing the `## Working Principles` section heading"); inverting the posture prose meaning while keeping the heading → STILL PASS (no longer prose-pinned). AC-2 satisfied.
- DONE: REPLACE TestCommissionStateBackendDecisionRule with a BEHAVIORAL check ... If no scaffolder codepath is invocable from a Go test, DEMOTE to a structural decision-table assertion and RECORD why REPLACE wasn't feasible
  DEMOTED (REPLACE infeasible — recorded). No Go scaffolder emits workflow frontmatter: commission is instruction-driven from SKILL.md's "Write the README with ..." steps; `claude-team` only *reads* an existing README's `state:` field (has_state_checkout), never decides/emits one. New test extracts the two labeled decision rows under the `**State backend` lead-in and asserts each binds its own outcome (split-root → `state: .spacedock-state`; inline → `$inline`/omit and NOT the split-root path). Commit 65be705b. Mutation-verified: swapping the binding → FAIL, dropping the Inline branch → FAIL, inverting condition prose with table kept → STILL PASS. AC-3 satisfied.
- DONE: No new grep-over-prose anywhere. Full `go test ./...` green + gofmt/vet clean.
  `go test ./...` → 732 passed in 12 packages; `gofmt -l` on all 3 changed files → clean; `go vet ./...` → no issues. The commission replacement is a per-row-bounded structural decision-table assertion (pins table shape + each row's bound outcome), not a free-floating substring grep. AC-4 satisfied.

### Summary

Executed the per-func disposition without re-litigating the classification: PRUNED 3 pure prose-pins + the no-force half of the ship-local test (that property is behaviorally owned by internal/status merge_policy_guard_test.go), DEMOTED the FO posture test to a `## Working Principles` heading lint, and DEMOTED (not REPLACED) the commission decision-rule test to a structural decision-table assertion — REPLACE was infeasible because the commission flow is LLM-instruction-driven from SKILL.md with no invocable Go scaffolder, which I recorded per the AC fallback. Every disposition is mutation-proven: the real structural/binding mutation FAILS while a meaning-inverting paraphrase that keeps the anchor stays green, closing the polarity hole by pinning structure instead of hardening greps (the v1-saga principle). Touched only `skills/integration/*_test.go` (3 commits on the worktree branch); full `go test ./...` is 732/732 green, gofmt/vet clean.
