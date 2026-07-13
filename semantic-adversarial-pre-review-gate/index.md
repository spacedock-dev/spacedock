---
title: Semantic adversarial pass before implementation requests review
status: ideation
source: Captain learning from spacedock-subspace Roborev adoption, 2026-07-14
started: 2026-07-13T16:24:34Z
completed:
verdict:
score: 0.85
worktree:
issue:
id: 4hgf1a01qkpm7p6a5hbxsn3a
---

Generalize the subspace Roborev learning into an implementation-exit quality gate in both the active development workflow and the shipped development workflow template.

## Problem

Implementation can satisfy its immediate checklist and tests while still missing representation drift, lifecycle edge cases, atomic validation gaps, scaling hazards, or false-green assertions. External review then spends its first pass discovering defects the implementer could have exposed through one structured semantic adversarial pass.

## Proposed approach

Add one concise pre-review contract to the `implementation` stage in `docs/dev/README.md` and `skills/commission/references/templates/development.md`. Require the implementer to trace changed values and events across representations and lifecycle phases; exercise an adjacent-variant matrix; prefer canonical validators or atomic record validation; inspect hot paths for multiplicative work and size limits; and ask how tests could pass while observable behavior remains wrong. Keep this author-side preparation separate from Roborev reviewer calibration.

## Out of scope

- Installing or configuring Roborev.
- Adding component-specific rules that belong in repository instructions or task acceptance criteria.
- Changing validation-stage ownership, feedback routing, or human gates.
- Treating presence of the prose alone as proof that the behavior improved.

## Acceptance criteria

**AC-1 (VALUE) - An implementation-stage worker catches a seeded semantic or lifecycle defect before requesting review because the adversarial pass forces the relevant invariant and adjacent variants to be checked.**
Verified by: ideation defines a comparative fixture or live drive with a planted false-green implementation; the revised development stage turns the case red or records the defect before review, while a removed/reordered adversarial-pass mutant does not.

**AC-2 - The active development workflow and newly commissioned development workflows carry the same generalized pre-review semantic adversarial contract.**
Verified by: `spacedock dispatch show-stage-def` against `docs/dev` and an isolated workflow commissioned from the shipped development template expose equivalent implementation outputs; a template-only or local-only mutant fails.

**AC-3 - The contract covers representation/lifecycle tracing, adjacent variants, atomic validation, scaling boundaries, and false-green test analysis without embedding component-specific policy.**
Verified by: fixture-backed dispatch evidence maps each category to a concrete implementation-report signal, while a seeded omission for each category is detected; structural checks may enforce shape but do not claim behavioral sufficiency.

**AC-4 - The author-side adversarial pass remains distinct from reviewer configuration and fresh validation.**
Verified by: the changed diff is bounded to development workflow definitions and their behavior-backed fixtures/tests; `.roborev.toml`, validation ownership, feedback routing, and review tooling remain unchanged.

## Test plan

Ideation should choose the smallest behavior oracle that can refute a ceremonial wording-only change. Prefer an isolated commissioned dev workflow plus a deterministic fake implementation/review boundary; add one representative live drive only if LLM response to the new instruction is the claim. Include a negative mutant that removes or reorders the pass and a planted false-green defect spanning at least two lifecycle representations. Run focused contract/commission tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
