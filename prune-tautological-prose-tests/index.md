---
id: rx13p1svcchfs8rtdmwsxq1k
title: Prune/demote/replace the tautological prose-pin integration tests (the banned grep-over-prose antipattern)
status: implementation
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
