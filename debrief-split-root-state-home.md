---
title: Debrief skill writes to the definition dir, not the split-root state checkout
status: ideation
group: cleanup
id: 7d47cgfj6h6z2xf5kk7ydbd9
sprint: 0240-lean-contract
started: 2026-06-30T16:20:13Z
---
The debrief skill (`skills/debrief/SKILL.md`) resolves `{dir}/_debriefs/` — read (session-boundary anchor), write (new file), and commit (current branch) — where `{dir}` = `spacedock status --discover` = the workflow DEFINITION dir (`docs/dev`, on `main`). For a split-root workflow (README declares `state:`), that is the wrong home: the established convention and the bulk of history live in the state checkout (`{state}/_debriefs/`), which auto-syncs via `spacedock state commit` and keeps session churn off the code branch — the same isolation the pr-merge mod enforces for `pr:`/`mod-block:`.

Symptoms observed (this repo, 2026-06-29): debriefs split across both homes; main-home debriefs left `main` ahead of `origin` needing a manual push outside the PR flow; and because the two homes numbered sequences independently, `2026-06-19-01` and `2026-06-21-01` exist in BOTH dirs as DIFFERENT sessions (a filename collision that a naive "de-duplicate" would silently destroy). Same split-root-unawareness family as the two deferred tooling gaps (`status --validate` near-dup ids; `status --resolve` archived-scope).

The documented intent already agrees: `docs/specs/state-behavior-extension.md` (V0 Layout) places `_debriefs/` *under* `.spacedock-state`, and `spacedock state commit` / the FO+ensign contract already isolate every other active-state write (entities, stage reports, archive) to the checkout. The debrief skill is the lone active-state writer that still resolves to the definition dir. The fix is skill-only — `skills/debrief/SKILL.md`, a different file from every other 0240 member — and ships no new binary behavior.

## Proposed approach

Resolve a `{debrief_root}` once, in Phase 1, from the README `state:` field (the discriminator the binary already interprets), then thread it through the four `_debriefs/` touch-points. Single-root resolves `{debrief_root} = {dir}` (today's behavior, unchanged). Split-root resolves `{debrief_root} = {state_checkout}` and switches the commit from a bare on-trunk `git commit` to a path-scoped commit + push in the checkout.

Specific skill edits (before → after):

1. **Phase 1, new Step 1b — "Resolve the debrief home":** after `Store the confirmed path as {dir}.`, add:
   > Read the `state:` field from `{dir}/README.md` frontmatter.
   > - Absent, empty, or the literal `$inline` → **single-root**: set `{debrief_root} = {dir}`; debriefs read/write/commit under `{dir}/_debriefs/` on the current branch (unchanged).
   > - Any other value → **split-root**: set `{state_checkout} = {dir}/{state}` and `{debrief_root} = {state_checkout}`. All debrief reads, the new-file write, and the commit run in the state checkout on its own branch — never in `{dir}` on the code branch. Resolve the state branch from the checkout itself: `git -C {state_checkout} rev-parse --abbrev-ref HEAD`.
   >
   > Use `{debrief_root}/_debriefs/` everywhere `_debriefs/` appears below.
2. **Phase 1, Step 2, item 1:** `Look for {dir}/_debriefs/*.md files.` → `Look for {debrief_root}/_debriefs/*.md files.`
3. **Phase 4, Step 1 (sequence number):** `{dir}/_debriefs/` → `{debrief_root}/_debriefs/`, plus a clause: *for split-root this counts the state-checkout debriefs only — never `{dir}/_debriefs/` — so numbering continues the established history and a same-named orphan in the definition dir cannot perturb the sequence.*
4. **Phase 4, Step 3 (write):** `mkdir -p {dir}/_debriefs` and the write path `{dir}/_debriefs/{date}-{sequence:02d}.md` → `{debrief_root}/_debriefs/...`.
5. **Phase 4, Step 4 (commit):** keep the bare `git add {dir}/_debriefs/... && git commit` for single-root; for split-root replace it with the path-scoped checkout commit + push:
   ```bash
   git -C {state_checkout} add -- _debriefs/{date}-{sequence:02d}.md
   git -C {state_checkout} commit -m "debrief: session {date} #{sequence} — {summary}" -- _debriefs/{date}-{sequence:02d}.md
   git -C {state_checkout} push origin {state_branch}
   ```
   On a non-fast-forward push rejection: `git -C {state_checkout} pull --rebase origin {state_branch}` then re-push (the single-file commit replays atop the peer's; disjoint paths → no conflict). If the rebase CONFLICTS: `git -C {state_checkout} rebase --abort`, report the conflicting path(s), and stop — never force-push or auto-resolve. No `origin` → commit locally, skip push. This is the same "Concurrency-safe state commits / Multi-writer sync / Rebase-conflict halt" discipline the ensign contract already prescribes for entity writes.
6. **Phase 4, final report line:** `Debrief written to {dir}/_debriefs/...` → `{debrief_root}/_debriefs/...`.

No separate docs-site page describes the debrief skill, so the skill text is the doc; the before/after above is the doc diff.

## Riskiest mechanism (spiked)

The riskiest path is the resolution + isolated-commit: does reading `state:` route reads/write/commit to the checkout on its own branch, leaving the trunk untouched — and is `spacedock state commit` reusable? Spiked against this live repo (2026-06-30):

- **Discriminator.** `docs/dev/README.md` carries `state: .spacedock-state`. `internal/status/state.go ClassifyState` is the single shipped interpreter: empty/`$inline` → inline; any other relative path → split-root with `filepath.Clean(state)`. The skill mirrors exactly that test — no new parsing.
- **Checkout path.** `{dir}/{state}` = `docs/dev/.spacedock-state`, confirmed a real **linked worktree** (`.git` → `gitdir: …/.git/worktrees/-spacedock-state`) on branch **`spacedock-state/dev`** with an `origin` remote. `git -C {checkout} rev-parse --abbrev-ref HEAD` returns `spacedock-state/dev`, matching `status.StateBranch` (basename rule, `state-branch:` override wins) — so asking the checkout is a robust substitute for re-deriving the convention.
- **Commit isolation.** The checkout is an orphan branch distinct from the code branch `main`; a commit there does not move `main`. The 3 orphan definition-dir debriefs (`2026-06-19-01`, `2026-06-19-02`, `2026-06-21-01`) are committed on `main` — the exact "left main ahead of origin" symptom this fix removes.
- **`spacedock state commit` is NOT reusable for debriefs.** Its handler (`internal/cli/state_sync.go runStateCommit`) resolves `<slug>` via `resolveEntityPath` (`{slug}.md` / `{slug}/index.md`); a `_debriefs/X.md` file is not an entity, so it exits 1 "no entity". The debrief commit therefore uses raw path-scoped git (above) — matching `commitEntityPathScoped`'s sequence — and needs no binary change.
- **Collision is real and must not be de-duped.** State checkout holds `2026-06-19-01.md` and `2026-06-21-01.md` as DIFFERENT sessions from the same-named definition-dir files. Computing the split-root sequence from the checkout only (AC-3) sidesteps this; migrating/de-duping the orphans is out of scope (would destroy a session).

This spike is the throwaway exercise that seeds the implementation's first test (AC-1's "file in checkout, trunk unmoved").

## Acceptance criteria

- **AC-1 (split-root routing — the end-value)** — for a `state:`-declared workflow, the debrief flow reads prior debriefs from, writes the new debrief to, and commits it in `{state_checkout}/_debriefs/` on the state branch, and the workflow's trunk does NOT move. *Measured on a split-root fixture:* the new `_debriefs/*.md` exists under the state checkout and not under `{dir}/_debriefs/`; the commit is on the state branch (`git -C {checkout} log -1 --format=%D`); and the trunk HEAD is byte-identical before and after the flow — the baseline the current bug moves the wrong way (it commits the debrief onto the code branch). Confirmed end-to-end by a live `/debrief` drive observing the same on-disk state.
- **AC-2 (single-root unchanged — the discriminator)** — a workflow with no `state:` (or empty/`$inline`) is unchanged: prior reads, the new write, and the commit all use `{dir}/_debriefs/` on the trunk. Verified on a single-root fixture as a negative control — the routing branch is taken only when `state:` resolves to a real checkout.
- **AC-3 (continuous numbering, no cross-home collision)** — for split-root the sequence number is computed from the state-checkout debriefs only, so numbering continues the established state-checkout history and a same-named debrief orphaned in `{dir}/_debriefs/` neither perturbs the count nor is overwritten/de-duplicated. Verified by a fixture whose state checkout is seeded with `…-01/02/03` and whose definition dir holds a colliding same-name orphan: the new file is `…-04` and the orphan is untouched.

## Out of scope

- **Migrating or de-duping** the 3 orphan definition-dir debriefs (`2026-06-19-01`, `2026-06-19-02`, `2026-06-21-01`). They collide by name with different sessions in the checkout; reconciling them is a manual cleanup decision, not the routing fix. This task only stops *new* debriefs from landing there.
- **Any new binary behavior.** The resolution reuses the already-shipped `state:` semantics; the commit is raw git. Consistent with the 0240 DoD ("no new binary behavior").

## Test plan

The skill is agent-executed markdown; its claim ("the debrief lands in the checkout, not the trunk") is proven by exercising the flow and observing on-disk state — never by grepping the skill prose (proof policy bans prose-grep; the `done` stage requires a live drive for a skill change). Two layers:

1. **Deterministic fixture regression guard (cheap, ~minutes).** A bash/fixture harness that builds real-git fixtures and runs the prescribed resolution + write + commit sequence:
   - *Split-root (AC-1):* repo with `docs/dev/README.md` (`state: .spacedock-state`, `trunk: main`), a `.spacedock-state` linked worktree on `spacedock-state/dev`, `main` checked out. Assert: new file under the checkout's `_debriefs/`, not `docs/dev/_debriefs/`; commit on `spacedock-state/dev`; **`main` HEAD identical before/after**.
   - *Single-root (AC-2, negative control):* README with no `state:`; assert file + commit land in `{dir}/_debriefs/` on the trunk.
   - *Numbering (AC-3):* checkout seeded `…-01/02/03` plus a colliding `…-01` orphan in the definition dir; assert new file is `…-04`, orphan untouched.
   Cheap and decisive for the mechanism, but it transcribes the skill's git commands, so it is a regression guard / live-drive seed — NOT a substitute for the live drive.
2. **Live skill drive (PASSED-gating, per the `done` proof policy).** Run the real `/debrief` skill via a host on a split-root fixture and assert the same durable on-disk state (file + commit in the checkout, trunk untouched, continuous sequence). Scope as a single-host on-disk-state drive; a full cross-host `internal/ensigncycle` shared scenario is available if cross-host parity is wanted, but is heavier than this skill needs.

Cost/complexity: fixture guard cheap (git plumbing); live drive the expensive-but-decisive proof. No Go changes. **Staff review recommended** before the ideation gate — this is both split-root behavior and skill integration, the two complexity triggers the README names.

## Stage Report: ideation

- DONE: The split-root routing mechanism (AC-1), spiked
  Spiked live (repo, 2026-06-30): `state: .spacedock-state` → checkout `docs/dev/.spacedock-state`, linked worktree on `spacedock-state/dev` w/ origin; commit there leaves `main` unmoved. Result recorded under `## Riskiest mechanism (spiked)`; AC-1 verifies file+commit in checkout and trunk-HEAD-unchanged on a split-root fixture + live `/debrief` drive.
- DONE: Single-root unchanged (AC-2, the discriminator)
  `ClassifyState` semantics (empty/`$inline` → inline) drive the branch; AC-2 verifies file+commit stay in `{dir}/_debriefs/` on the trunk via a single-root negative-control fixture.
- DONE: Continuous numbering (AC-3)
  Confirmed the real cross-home collision (`2026-06-19-01`, `2026-06-21-01` are different sessions in both homes); AC-3 computes the sequence from the checkout only and verifies a colliding orphan is neither counted nor de-duped. Skill-only edit to `skills/debrief/SKILL.md`.

### Summary

Refined the body in place: added Proposed approach (six before/after skill edits to `skills/debrief/SKILL.md`), a spiked Riskiest-mechanism section, sharpened ACs (AC-1 now measures trunk-HEAD-unchanged against the buggy baseline), an Out-of-scope (no orphan migration, no binary change), and a two-layer test plan (cheap fixture guard + PASSED-gating live drive). Key spike findings: `spacedock state commit` is entity-only so the debrief commit uses raw path-scoped git matching the ensign discipline; the resolution reuses the already-shipped `state:` semantics with no new binary behavior. Flagged staff review (split-root + skill integration). No product files were edited — ideation is design-only.
