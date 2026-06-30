---
title: Debrief skill writes to the definition dir, not the split-root state checkout
status: backlog
group: cleanup
id: 7d47cgfj6h6z2xf5kk7ydbd9
---
The debrief skill (`skills/debrief/SKILL.md`) resolves `{dir}/_debriefs/` — read (session-boundary anchor), write (new file), and commit (current branch) — where `{dir}` = `spacedock status --discover` = the workflow DEFINITION dir (`docs/dev`, on `main`). For a split-root workflow (README declares `state:`), that is the wrong home: the established convention and the bulk of history live in the state checkout (`{state}/_debriefs/`), which auto-syncs via `spacedock state commit` and keeps session churn off the code branch — the same isolation the pr-merge mod enforces for `pr:`/`mod-block:`.

Symptoms observed (this repo, 2026-06-29): debriefs split across both homes; main-home debriefs left `main` ahead of `origin` needing a manual push outside the PR flow; and because the two homes numbered sequences independently, `2026-06-19-01` and `2026-06-21-01` exist in BOTH dirs as DIFFERENT sessions (a filename collision that a naive "de-duplicate" would silently destroy). Same split-root-unawareness family as the two deferred tooling gaps (`status --validate` near-dup ids; `status --resolve` archived-scope).

## Acceptance criteria
- **AC-1** — for a split-root workflow (`state:` declared), the debrief skill reads prior debriefs from and writes the new debrief to `{state_checkout}/_debriefs/`, committing+pushing on the state branch — not `{definition_dir}/_debriefs/` on the trunk. Verified by a test running the debrief flow on a split-root fixture and asserting the file + commit land in the state checkout.
- **AC-2** — a single-root (non-split) workflow is unchanged: debriefs stay in `{dir}/_debriefs/` on the trunk. Verified on a single-root fixture.
- **AC-3** — the sequence number is computed from the state-checkout debriefs for split-root, so numbering is continuous with no cross-home collisions. Verified by a fixture seeded with existing state-checkout debriefs.
