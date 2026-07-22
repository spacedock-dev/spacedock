# Local Spacedock commit evidence

- Repository: `/Users/clkao/git/spacedock-research/spacedock-v1`
- Worktree: `.worktrees/spacedock-ensign-gate-review-presentation-command`
- Base: `fa240a76`
- Candidate: `612b72fca1ef98a0dde97cf0b1cecdf2355a7b16`
- Parent: `cf6008fd85bfb371a83ae747100f06a18d91ac21`
- Subject: `Clarify complete presentation association`
- Author date: `2026-07-22T23:17:50+08:00`
- Candidate worktree porcelain at package assembly: empty

Exact base-to-candidate name-status:

```text
M	docs/site/concepts/gates-and-decisions.md
M	internal/cli/gate_test.go
M	internal/contractlint/fo_function_reference_invariant_test.go
M	skills/present-gate/SKILL.md
```

Exact base-to-candidate diff stat:

```text
 docs/site/concepts/gates-and-decisions.md           |  8 ++++++++
 internal/cli/gate_test.go                           | 21 +++++++++++++++++++++
 .../fo_function_reference_invariant_test.go         |  8 +++++---
 skills/present-gate/SKILL.md                        | 15 +++++++++++++++
 4 files changed, 49 insertions(+), 3 deletions(-)
```

The exact commands and outcomes for normal, race, documentation, formatting, detached adversarial, and Roborev verification are retained unchanged in `validator-report.md`.

