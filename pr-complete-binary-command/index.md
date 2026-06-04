---
id: p23sxe8ec3mmwgekvz9041a9
title: "spacedock pr complete — collapse the Merge-and-Cleanup ceremony into one binary command"
status: backlog
source: "captain (2026-06-04) — token-efficiency decomposition of first-officer-shared-core.md + the binary-simplification roadmap #3. Merge-and-Cleanup is MECHANICAL ceremony, so it routes to a binary command (not a lazy skill): collapsing it removes the prose outright. The FO ran this ceremony ~7x BY HAND this session (at/n3/2a/am/6b/ep/zd/p4 merges) — the heaviest, most error-prone sequence — the strongest single binary-command ROI."
score: "0.36"
worktree:
started:
completed:
verdict:
issue:
---

The post-merge half of Merge-and-Cleanup (+ the Ship-Local Ceremony) is pure mechanical ceremony: mod-block clear (separate commit per the standalone rule) → terminalize (`--set completed verdict= worktree=`) → archive → worktree-remove → local-branch-delete → state push, with the terminal team teardown. It is the heaviest block in `first-officer-shared-core.md` and the most error-prone to run by hand. Ceremony → BINARY command (the roadmap's lever 1), which removes the prose from the boot read entirely rather than deferring it to a skill.

## Proposed approach (ideation firms)

`spacedock pr complete {slug}` orchestrates the existing binary operations transactionally: clear mod-block (own commit) → terminalize → archive → remove worktree → delete local branch → push state. Surfaces failures the same shape the `pr-merge` mod uses. The FO prose shrinks to "on merged PR detected, run `spacedock pr complete {slug}`" + the team-teardown step (team-tool teardown stays FO/runtime-side — it is a Claude tool, not shell-able, per the using-claude-team boundary).

**Constraints:**
- **Idempotency** — re-running on an already-archived entity is a clean no-op.
- **Guard-honoring** — the separate-commit discipline for mod-block clear (the `status --set` guard already enforces it) must be preserved; the command must not need `--force` on the happy path.
- HIGH-STAKES (status mutation/guard + CI/release machinery) → detached audit at validation.

## Acceptance criteria (seed)

- **AC-1 (seed):** `spacedock pr complete {slug}` takes a merged-PR-state entity to archived (mod-block cleared in its own commit, terminalized verdict=PASSED, worktree removed, local branch deleted, state pushed) — verified by an end-to-end fixture from merged-PR state to archived passing through every guard with no `--force`.
- **AC-2 (seed):** Idempotent — re-running on an already-archived entity exits 0 as a no-op (fixture test).
- **AC-3 (seed):** The FO contract prose for the post-merge ceremony is replaced by the single command invocation (instruction-text check), shrinking that block of `first-officer-shared-core.md`.

## Out of scope

- The team-tool teardown (Claude tool, not shell-able — stays in the using-claude-team skill / runtime).
- #1 `state sync` and #2 `dispatch advance` (sibling roadmap binary commands; file separately if/when Phase 1 is greenlit).

## Notes

Roadmap #3 (binary-simplification-roadmap.md, refreshed 2026-06-04 — named the strongest single binary ROI; the FO ran it ~7x manually this session). The qs reconcile helper is the test-coverage template. Sibling lever-2 extractions: `gate-presentation-skill-extraction`, `feedback-rejection-flow-skill-extraction`.
