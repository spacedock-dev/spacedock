---
id: azh879wdzm72ysxg16hbg39q
title: Enforce non-tautological tests beyond the reactive Proof policy, and close the commission-template gap
status: backlog
source: "Two related findings from this session's audit work, not yet designed: (1) docs/dev's own Proof policy ('no prose-grep over instruction files', the detached adversarial audit, the AC template's 'Verified by: ... something outside this task body ... that can fail' clause) only catches tautology reactively, at high-stakes-surface merge time via human/reviewer judgment — there is no standing, automatic check against the mirror-assertion / no-op-assertion patterns that four real tests in this repo turned out to have. (2) The commission-skill templates that scaffold NEW workflows do not carry equivalent discipline: skills/commission/references/templates/development.md has the relevant bullets ('External-proof acceptance criteria', 'Detached adversarial audit') only as OPT-IN text copied in 'when commissioning', missing docs/dev's 'outside this task body'/'that can fail' clauses in its own AC template; experiment.md has only a one-line delegation and no AC-N construct at all; refinement.md has neither; skills/commission/SKILL.md's base entity template (~line 467) carries the same gapped 'Verified by' stub for every commissioned workflow. Reference material: github.com/kenn-io/middleman skills/testing-without-tautologies/SKILL.md (mutation-check discipline, mirror-assertion / mock-tests-subject / branch-double-reuse / upstream-functionality / blindingly-obvious-assertion checklist). A design pass (mutation-tooling landscape for Go, whether internal/contractlint's existing go/ast structural-check architecture can be extended to catch the observed patterns automatically, and concrete before/after diffs for docs/dev/README.md plus the commission templates) was run separately this session and its output should seed ideation directly rather than be re-derived."
started:
completed:
verdict:
score:
worktree:
issue:
---

Design and land a standing mechanism (not just reactive review) against tautological tests, and bring the commission-skill templates that scaffold new workflows up to docs/dev's own Proof-policy bar so new workflows don't inherit the gap.
