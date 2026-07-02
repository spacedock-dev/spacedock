# Progress

## 2026-06-19 — 0223-pi-dispatch-contract sprint (Shaping FO session on z-ai/glm-5.2)

### Current state
Sprint `0223-pi-dispatch-contract` — 3 members (re-carved from 4):
- `eq` `pi-install-managed-skill-placement` (ideation) — merged task, supersedes archived `k8t`+`2m1`
- `bdt` `pi-dispatch-model-stamping` (ideation)
- `b2` `pi-back-channel-dispatch` (capstone, ideation)

### Staff reviews
- Review #1 (`3adf00ee`): found child-cwd blocker + 2 minor gaps → all closed by re-carve
- Review #2 (`efff49c9`): found 7 gaps (1 spec, 6 doc) → closing in progress

### Gap closure status (staff review #2)
- [x] Gap 1 (`eq` repoRoot source) — closed via fold-in commit `e862e42e`; resolution (a): retire repo-path Stat checks, add `spacedockPackageOK`, remove cwd fallback; AC-3 revised
- [ ] Gap 6 (`b2` fold-in AC-2 sub-bullet reconciliation) — parallel fold-in dispatched (run in flight)
- [ ] Gaps 2–5, 7 (index: Q11/Q12 rewrite, Sequencing section, DoD Proven-by items, Q3/Q9/Q13 cosmetics) — Shaping FO to fix inline

### Recent commits
- `e862e42e` pi-install-managed-skill-placement: staff review #2 fold-in — close gap 1
- `cb078d10` (peer) pi-back-channel-dispatch: gap-1 re-check — install-managed placement removes child-cwd seam
- `75469bb5` pi-install-managed-skill-placement: ideation — mechanism confirmed, spike PASSED
- `b2f2cbf3` (main) sprint: 0223 re-carve — 3 members
