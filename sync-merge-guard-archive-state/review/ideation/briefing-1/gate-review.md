# IDEATION GATE — Durable split-root merge finalization (`rd`)

Recommendation: **APPROVE the corrected design; keep implementation pending until 6y lands.**

## Chosen direction

Extract the existing push/rebase/HALT discipline into one standard-library internal publisher. `merge guard` publishes its path-scoped archive commit before reporting final success; `state commit <slug>` becomes the restart path for a clean archived slug and publishes existing ahead history without staging or committing archived state.

## Evidence

- A real two-clone spike left origin at active `pr-merge:99` while host A was clean and one archive commit ahead; no existing supported command could publish it.
- The design adds no public verb and reuses the current two-host CLI fixtures.
- Staff review forced archived scope to remain publish-only, duplicate shapes to fail closed, and interrupted rebase preflight to abort before entity resolution or mutation.
- ACs bind remote refs, HEAD/index/worktree state, exact one-value JSON plus EOF, path-scoped sibling isolation, and truthful inline/no-origin outcomes.
- Corrected expected surface is 10 files and about 560 changed LOC, with a 13-file/900-line reset ceiling.
- Independent re-review returned APPROVE with no material finding.

## Boundary

No autostash, raw-Git FO workaround, new public command, gate/provider/PR-host change, or 6y lifecycle expansion. Sibling dirt plus non-fast-forward fails recoverably and resumes after the dirt settles.

## Decision ask

Approve the design now; retain the application pending until 6y lands, then enter implementation in `.worktrees/spacedock-ensign-sync-merge-guard-archive-state`.
