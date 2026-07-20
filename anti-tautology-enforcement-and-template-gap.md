---
id: azh879wdzm72ysxg16hbg39q
title: Enforce non-tautological tests beyond the reactive Proof policy, and close the commission-template gap
status: backlog
source: "Two related findings from this session's audit work, not yet designed: (1) docs/dev's own Proof policy ('no prose-grep over instruction files', the detached adversarial audit, the AC template's 'Verified by: ... something outside this task body ... that can fail' clause) only catches tautology reactively, at high-stakes-surface merge time via human/reviewer judgment -- there is no standing, automatic check against the mirror-assertion / no-op-assertion patterns that four real tests in this repo turned out to have. (2) The commission-skill templates that scaffold NEW workflows do not carry equivalent discipline for GO TEST CODE tautology (as opposed to instruction-file/prose tautology, which is entity ey / proof-policy-shipped-scaffolding's separate, pre-existing scope -- see OVERLAP NOTE below): skills/commission/references/templates/development.md's AC template stub and skills/commission/SKILL.md's base AC template both lack docs/dev's own 'outside this task body'/'that can fail' clauses. Reference: github.com/kenn-io/middleman skills/testing-without-tautologies/SKILL.md. A design workflow + an independent fable review ran this session and both converged (mechanism choice, scope, sequencing all hold per fable's independent check) on: NO automated hard AST gate for mirror-assertions (an internal/testlint check for the OTHER pattern, assertion-free tests, is its own sibling entity: testlint-assertion-free-gate) -- instead extend the existing detached-adversarial-audit trigger to fire on AC PROVENANCE ('any AC whose expected value is derived from the same package's production functions or constants'), not the originally-drafted broader 'equality/byte-identity check' wording fable found would over-fire on nearly every unit test in this repo. Concrete diffs exist for docs/dev/README.md's Proof policy + pr-merge gate rule, development.md, and SKILL.md. CORRECTED accounting (fable caught this): internal/status/boot_probe_parity_test.go is NOT a stampID mirror-assertion instance -- it mirrors a different production CONSTANT (teamStateNeutralHint), not a function call; the confirmed mirror-assertion-via-shared-function count is 2 (native_new_test.go, zz_independent_parity_test.go), already covered under the separate tautological-test-fixes entity, don't double-count here. OVERLAP NOTE, UNRESOLVED -- flagging for the captain, not silently resolving: entity ey (proof-policy-shipped-scaffolding, filed 2026-06-04, pre-existing) targets the SAME file (skills/commission/references/templates/development.md) for a related-but-distinct concern -- porting the INSTRUCTION-FILE/prose tautology test (not code-test tautology) to shipped scaffolding, plus first-officer-shared-core.md and ensign-shared-core.md, with a heavier behavioral AC (a live scenario proving a validator REJECTS a presence-only proof). This entity's development.md diff and ey's development.md target could collide if ideated independently without coordination. Captain has not yet said how to reconcile (fold together, sequence, or keep fully separate with a coordination note in each) -- do not dispatch this entity's ideation until that's decided."
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0260-proportionality
group: test-cleanups
---

Design and land a standing mechanism (not just reactive review) against tautological tests, and bring the commission-skill templates that scaffold new workflows up to docs/dev's own Proof-policy bar so new workflows don't inherit the gap.

## Scope trim (0260 re-lock 2026-07-20)

This entity owns the CONTRACT half only: the falsifiable-evidence rule (AC evidence must be
able to fail; "show the change that makes it fail") and removal of "5/5 passed is
sufficient" so gates read assertion content. The commission-template half of the original
scope moves to the template group (proof-policy-shipped-scaffolding and
template-rigor-propagation) — one owner per surface, no duplicate delivery.
