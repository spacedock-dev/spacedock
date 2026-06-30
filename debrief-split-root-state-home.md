---
title: Debrief skill writes to the definition dir, not the split-root state checkout
status: validation
group: cleanup
id: 7d47cgfj6h6z2xf5kk7ydbd9
sprint: 0240-lean-contract
started: 2026-06-30T16:20:13Z
worktree: .worktrees/spacedock-ensign-debrief-split-root-state-home
mod-block: merge:pr-merge
pr: pr-merge:449
---
The debrief skill (`skills/debrief/SKILL.md`) resolves `{dir}/_debriefs/` — read (session-boundary anchor), write (new file), and commit (current branch) — where `{dir}` = `spacedock status --discover` = the workflow DEFINITION dir (`docs/dev`, on `main`). For a split-root workflow (README declares `state:`), that is the wrong home: the established convention and the bulk of history live in the state checkout (`{state}/_debriefs/`), which auto-syncs via `spacedock state commit` and keeps session churn off the code branch — the same isolation the pr-merge mod enforces for `pr:`/`mod-block:`.

Symptoms observed (this repo, 2026-06-29): debriefs split across both homes; main-home debriefs left `main` ahead of `origin` needing a manual push outside the PR flow; and because the two homes numbered sequences independently, `2026-06-19-01` and `2026-06-21-01` exist in BOTH dirs as DIFFERENT sessions (a filename collision that a naive "de-duplicate" would silently destroy). Same split-root-unawareness family as the two deferred tooling gaps (`status --validate` near-dup ids; `status --resolve` archived-scope).

The documented intent already agrees: `docs/specs/state-behavior-extension.md` (V0 Layout) places `_debriefs/` *under* `.spacedock-state`, and `spacedock state commit` / the FO+ensign contract already isolate every other active-state write (entities, stage reports, archive) to the checkout. The debrief skill is the lone active-state writer that still resolves to the definition dir. The fix is skill-only — `skills/debrief/SKILL.md`, a different file from every other 0240 member — and ships no new binary behavior.

## Proposed approach

Resolve a `{debrief_root}` once, in Phase 1, by consuming the boot record the binary already computes (`status --boot --json`: its `state_backend` discriminator, the resolved absolute `entity_dir`, and `entity_dir_present`) rather than re-parsing the README `state:` field — so the skill inherits the binary's full `state:` resolution (including ClassifyState's absolute-path / `..` rejections) instead of re-implementing a partial copy of it. Set `{debrief_root} = entity_dir` in both modes, then thread it through the four `_debriefs/` touch-points. Single-root: `entity_dir == {dir}` (today's behavior, unchanged). Split-root: `entity_dir == {state_checkout}`, and the commit switches from a bare on-trunk `git commit` to a path-scoped commit + push in the checkout. Before either path touches `_debriefs/`, a halt-gate stops a split-root flow whose checkout is declared but not materialized (`entity_dir_present == false`).

Specific skill edits (before → after):

1. **Phase 1, new Step 1b — "Resolve the debrief home":** after `Store the confirmed path as {dir}.`, add:
   > Run `${SPACEDOCK_BIN:-spacedock} status --boot --json --workflow-dir {dir}` once and read three fields the binary already resolves (so the skill never re-parses `state:` itself): `state_backend` (`single-root` | `split-root`), `entity_dir` (the absolute resolved debrief home), and `entity_dir_present`. Set `{debrief_root} = entity_dir`.
   > - `state_backend == single-root` → `entity_dir` is `{dir}`; debriefs read/write/commit under `{debrief_root}/_debriefs/` on the current branch (unchanged).
   > - `state_backend == split-root` → set `{state_checkout} = entity_dir` (`= {debrief_root}`). **Halt-gate (mirrors the FO's «state.ensure-ready» at `skills/first-officer/references/first-officer-shared-core.md`):** if `entity_dir_present == false` the declared checkout is NOT materialized (orphan state branch without a linked worktree — fresh clone or removed worktree); HALT, report "state not initialized," and stop — do NOT fall back to writing into `{dir}` on the code branch. Tell the captain to run `spacedock state init` (manual fallback: `git fetch origin {state_branch} && git worktree add {state_checkout} {state_branch}`), then re-run. Otherwise all debrief reads, the new-file write, and the commit run in the state checkout on its own branch — never in `{dir}` on the code branch. Resolve the state branch from the checkout itself: `git -C {state_checkout} rev-parse --abbrev-ref HEAD`.
   >
   > Use `{debrief_root}/_debriefs/` everywhere `_debriefs/` appears below.
   >
   > The debrief skill is user-invocable independent of an FO boot, so it cannot assume the checkout has converged — hence the explicit `entity_dir_present` gate before any checkout git op.
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

- **Discriminator.** `docs/dev/README.md` carries `state: .spacedock-state`. `internal/status/state.go ClassifyState` is the single shipped interpreter (empty/`$inline` → inline; any other relative path → split-root with `filepath.Clean`, after rejecting absolute or `..`-escaping paths at state.go:49–55). The skill does NOT re-implement that test — it consumes the binary's already-resolved `state_backend` + `entity_dir` from `status --boot --json` (boot.go:206–217 sets `split-root` iff `entity_dir != definition_dir`), so it inherits the full resolution including the absolute/`..` rejections a hand-copied prose test would drop. Verified live (this repo, 2026-06-30): `status --boot --json` on `docs/dev` emits `state_backend: "split-root"`, `entity_dir: …/docs/dev/.spacedock-state`, `entity_dir_present: true`.
- **Absent-checkout signal.** `entity_dir_present` is `os.Stat(entity_dir).IsDir()` (boot.go:218–220), so a `state:`-declared workflow whose checkout is a bare orphan branch with no linked worktree (fresh clone / removed worktree) boots `entity_dir_present == false` — the exact silent failure the FO halt-gates at first-officer-shared-core.md:159. The debrief skill, user-invocable without an FO boot, gates on the same signal (AC-4) so it never runs `git -C {state_checkout}` against a missing dir.
- **Checkout path.** `{dir}/{state}` = `docs/dev/.spacedock-state`, confirmed a real **linked worktree** (`.git` → `gitdir: …/.git/worktrees/-spacedock-state`) on branch **`spacedock-state/dev`** with an `origin` remote. `git -C {checkout} rev-parse --abbrev-ref HEAD` returns `spacedock-state/dev`, matching `status.StateBranch` (basename rule, `state-branch:` override wins) — so asking the checkout is a robust substitute for re-deriving the convention.
- **Commit isolation.** The checkout is an orphan branch distinct from the code branch `main`; a commit there does not move `main`. The 3 orphan definition-dir debriefs (`2026-06-19-01`, `2026-06-19-02`, `2026-06-21-01`) are committed on `main` — the exact "left main ahead of origin" symptom this fix removes.
- **`spacedock state commit` is NOT reusable for debriefs.** Its handler (`internal/cli/state_sync.go runStateCommit`) resolves `<slug>` via `resolveEntityPath` (`{slug}.md` / `{slug}/index.md`); a `_debriefs/X.md` file is not an entity, so it exits 1 "no entity". The debrief commit therefore uses raw path-scoped git (above) — matching `commitEntityPathScoped`'s sequence — and needs no binary change.
- **Collision is real and must not be de-duped.** State checkout holds `2026-06-19-01.md` and `2026-06-21-01.md` as DIFFERENT sessions from the same-named definition-dir files. Computing the split-root sequence from the checkout only (AC-3) sidesteps this; migrating/de-duping the orphans is out of scope (would destroy a session).

This spike is the throwaway exercise that seeds the implementation's first test (AC-1's "file in checkout, trunk unmoved").

## Acceptance criteria

- **AC-1 (split-root routing — the end-value)** — for a `state:`-declared workflow, the debrief flow reads prior debriefs from, writes the new debrief to, and commits it in `{state_checkout}/_debriefs/` on the state branch, and the workflow's trunk does NOT move. *Measured on a split-root fixture:* the new `_debriefs/*.md` exists under the state checkout and not under `{dir}/_debriefs/`; the commit is on the state branch (`git -C {checkout} log -1 --format=%D`); and the trunk HEAD is byte-identical before and after the flow — the baseline the current bug moves the wrong way (it commits the debrief onto the code branch). Confirmed end-to-end by a live `/debrief` drive observing the same on-disk state.
- **AC-2 (single-root unchanged — the discriminator)** — a workflow with no `state:` (or empty/`$inline`) is unchanged: prior reads, the new write, and the commit all use `{dir}/_debriefs/` on the trunk. Verified on a single-root fixture as a negative control — the routing branch is taken only when `state:` resolves to a real checkout.
- **AC-3 (continuous numbering, no cross-home collision)** — for split-root the sequence number is computed from the state-checkout debriefs only, so numbering continues the established state-checkout history and a same-named debrief orphaned in `{dir}/_debriefs/` neither perturbs the count nor is overwritten/de-duplicated. Verified by a fixture whose state checkout is seeded with `…-01/02/03` and whose definition dir holds a colliding same-name orphan: the new file is `…-04` and the orphan is untouched.
- **AC-4 (declared-but-absent checkout halts, not silent fallback)** — for a `state:`-declared workflow whose checkout is NOT materialized (`entity_dir_present == false` — fresh clone or removed worktree), the debrief flow HALTS with "state not initialized" and writes/commits NOTHING; it does not silently fall back to writing into the definition dir on the code branch, nor emit opaque `git -C {missing}` errors. *Measured on a fixture* whose README declares `state: .spacedock-state` with no linked worktree at that path: the flow stops before any `_debriefs/` write — no new file appears under `{dir}/_debriefs/` or the absent checkout, and neither the trunk nor a state branch advances. This mirrors the FO's «state.ensure-ready» halt-gate; the skill is user-invocable independent of an FO boot, so it cannot assume convergence. The baseline this guards against: the unguarded split-root branch runs `git -C {state_checkout}` against a missing dir and produces opaque errors.

## Out of scope

- **Migrating or de-duping** the 3 orphan definition-dir debriefs (`2026-06-19-01`, `2026-06-19-02`, `2026-06-21-01`). They collide by name with different sessions in the checkout; reconciling them is a manual cleanup decision, not the routing fix. This task only stops *new* debriefs from landing there.
- **Any new binary behavior.** The resolution reuses the already-shipped `state:` semantics; the commit is raw git. Consistent with the 0240 DoD ("no new binary behavior").

## Test plan

The skill is agent-executed markdown; its claim ("the debrief lands in the checkout, not the trunk") is proven by exercising the flow and observing on-disk state — never by grepping the skill prose (proof policy bans prose-grep; the `done` stage requires a live drive for a skill change). Two layers:

1. **Deterministic fixture regression guard (cheap, ~minutes).** A bash/fixture harness that builds real-git fixtures and runs the prescribed resolution + write + commit sequence:
   - *Split-root (AC-1):* repo with `docs/dev/README.md` (`state: .spacedock-state`, `trunk: main`), a `.spacedock-state` linked worktree on `spacedock-state/dev`, `main` checked out. Assert: new file under the checkout's `_debriefs/`, not `docs/dev/_debriefs/`; commit on `spacedock-state/dev`; **`main` HEAD identical before/after**.
   - *Single-root (AC-2, negative control):* README with no `state:`; assert file + commit land in `{dir}/_debriefs/` on the trunk.
   - *Numbering (AC-3):* checkout seeded `…-01/02/03` plus a colliding `…-01` orphan in the definition dir; assert new file is `…-04`, orphan untouched.
   - *Absent-checkout (AC-4):* README declares `state: .spacedock-state` but no linked worktree exists at that path (`entity_dir_present == false`, simulating a fresh clone / removed worktree); assert the flow halts with "state not initialized," writes no `_debriefs/` file anywhere, and advances no branch. Fixture-decisive — the halt path needs no live drive.
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

## Stage Report: ideation (cycle 2)

- DONE: Add the declared-but-absent-state-checkout error path — Phase-1 halt-gate reading `entity_dir_present`/`entity_dir` from `status --boot --json`
  Step 1b (edit 1) now runs `status --boot --json` and HALTS with "state not initialized" when `state_backend == split-root && entity_dir_present == false`, mirroring «state.ensure-ready» (first-officer-shared-core.md:159); gate sits before any `git -C {state_checkout}` op.
- DONE: Add an AC + fixture for the declared-but-absent case
  Added AC-4 (halt, write/commit nothing, no silent fallback to definition dir, no branch advance) plus an Absent-checkout fixture in test-plan layer 1; baseline-that-moves-wrong-way named (opaque `git -C {missing}` errors).
- DONE: (Polish) Drop the "mirrors exactly that ClassifyState" claim — consume `entity_dir` from `status --boot --json` instead of re-parsing `state:`
  Resolution now sets `{debrief_root} = entity_dir` from the boot record (single + split alike); Discriminator spike bullet rewritten to "does NOT re-implement that test," noting state.go:49–55 absolute/`..` rejections the prose copy would drop. One fix closes both this and the absent-checkout gap.

### Summary

Folded staff-review M7 (declared-but-absent state checkout) and the 7d4 Polish into the gate-APPROVED design without redesigning. Both resolved by one mechanism change: Step 1b now consumes `state_backend` / `entity_dir` / `entity_dir_present` from `status --boot --json` (verified live: boot.go:206–220 emits all three) rather than re-parsing `state:`, so the skill inherits ClassifyState's full absolute/`..` resolution and gains a presence halt-gate. Routing structure, four touch-points, commit discipline, and existing ACs unchanged. No product files edited — ideation is design-only.

## Stage Report: implementation

- DONE: skills/debrief/SKILL.md Phase-1 Step 1b resolves {debrief_root}=entity_dir from status --boot --json and threads it through all four _debriefs/ touch-points; split-root commit becomes a path-scoped state-checkout commit+push with the rebase-conflict halt, single-root path unchanged
  Commit `84765fed`. New Step 1b reads `state_backend`/`entity_dir`/`entity_dir_present`, sets `{debrief_root}=entity_dir`; threaded through read (Phase 1 Step 2), sequence (Phase 4 Step 1), write (Phase 4 Step 3), commit (Phase 4 Step 4). Step 4 now branches: single-root keeps the bare on-trunk commit; split-root does `git -C {state_checkout} add/commit -- _debriefs/<file>` + `push origin {state_branch}`, with non-FF → `pull --rebase` + re-push, rebase-CONFLICT → `rebase --abort`+report+stop (never force/auto-resolve), no-origin → commit local + skip push.
- DONE: Absent-checkout halt-gate (entity_dir_present==false -> "state not initialized", write/commit NOTHING, no silent fallback to the definition dir) sits before any git -C {state_checkout} op
  Step 1b split-root branch HALTs on `entity_dir_present` false before any `_debriefs/` read/write/commit; routed to `spacedock state init` (manual fallback documented). Proven by AC-4: under the gate the flow writes nothing and no branch advances; with the gate removed the buggy flow wrote into the definition dir on the code branch (the regression the guard locks out).
- DONE: Deterministic fixture regression guard covers AC-1 (split-root: file in checkout, trunk HEAD byte-identical), AC-2 (single-root negative control), AC-3 (numbering from checkout only -> ...-04, orphan untouched), AC-4 (absent-checkout halt)
  `skills/integration/debrief_split_root_test.go` (commit `84765fed`), 4/4 PASS. Drives the real `status --boot --json` over real-git fixtures (linked orphan-branch state worktree per the commission recipe) and transcribes the skill's resolve+write+commit. Bug-catch verified: temporarily reverting the flow to pre-fix behavior fails AC-1/AC-3/AC-4 while AC-2 (negative control) stays green.

### Summary

Skill-only fix plus a Go regression guard; no binary change (reuses the already-shipped `status --boot --json` resolution). `skills/debrief/SKILL.md` now routes split-root debriefs to the state checkout on its own branch (path-scoped commit+push, rebase-conflict halt) and halt-gates a declared-but-absent checkout, while the single-root path is unchanged. The guard (skills/integration, 4/4 PASS) measures file-in-checkout + trunk-HEAD-byte-identical (AC-1), single-root negative control (AC-2), checkout-only numbering past a colliding orphan → -04 (AC-3), and the absent-checkout halt (AC-4); flipping the flow to the pre-fix behavior fails the three split-root ACs, proving the guard bites. The PASSED-gating live `/debrief` drive is the validation stage's job, per the test plan.

## Stage Report: validation

- DONE: MEASURE AC-1..AC-4 via the fixture guard ... and confirm it BITES by the documented revert
  All four PASS via `go test ./skills/integration/ -run Debrief` (AC-1 SplitRootRoutesToCheckout, AC-2 SingleRootUnchanged, AC-3 NumberingIgnoresOrphan, AC-4 AbsentCheckoutHalts). BITE measured: reverting `runDebriefFlow` to pre-fix (root=defDir always, no halt-gate, bare on-trunk commit) → AC-1/AC-3/AC-4 FAIL, AC-2 stays green; failure msgs show the bug (writes `docs/dev/_debriefs/2026-06-30-01.md`, numbers `-02` from the orphan, never halts). Reverted my edit, re-green (full pkg `-count=1` ok). NOTE: the checklist's literal `-run DebriefSplitRoot` substring-matches only 2 of 4 (SingleRoot/AbsentCheckout lack "SplitRoot") — used `-run Debrief` for the true four.
- DONE: LIVE /debrief drive (PASSED-gating for a skill change) on a split-root fixture via a host, observe durable on-disk state
  Drove the worktree's FIXED `skills/debrief/SKILL.md` (installed plugin copy is pre-fix) via a host subagent against a real-git split-root fixture (bare origin + linked `.spacedock-state` worktree on `spacedock-state/dev`, seeded `-01/-02/-03`), worktree binary (`contract 2`) as SPACEDOCK_BIN. Independently verified vs baseline: new `2026-06-30-04.md` under `{state_checkout}/_debriefs/`, ABSENT from `{dir}/_debriefs/`; commit `f4319cd` on `spacedock-state/dev` (`%D = HEAD -> spacedock-state/dev, origin/spacedock-state/dev`), path-scoped to the one file, NOT an ancestor of `main`; trunk `main` byte-identical `4463fcb3` before/after; push fast-forwarded origin; sequence continued history to `-04`.
- DONE: Confirm the single-root SKILL.md path is unchanged (negative control); reproduce each AC's "Verified by" clause; deliver a PASSED/REJECTED recommendation with evidence
  Single-root unchanged: AC-2 fixture green AND green under the pre-fix revert (routing branch fires only for a real checkout) AND the diff preserves the bare on-trunk single-root commit. Each AC's "Verified by" reproduced: AC-1 file-in-checkout/commit-on-state-branch/trunk-unmoved + live drive; AC-2 single-root negative control; AC-3 seeded -01/02/03 + colliding orphan → new -04, orphan byte-identical; AC-4 declared-but-absent checkout halts, writes/commits nothing, no branch advances. gofmt clean, `go vet` clean, full `skills/integration` pkg PASS.

### Summary

Recommendation: **PASSED.** The skill-only fix (`skills/debrief/SKILL.md`, commit `84765fed`) routes split-root debriefs to the state checkout on its own branch and leaves the trunk untouched, halt-gates a declared-but-absent checkout, and leaves the single-root path unchanged. The fixture guard measures all four ACs (4/4 PASS) and genuinely bites (pre-fix revert fails AC-1/3/4, AC-2 holds). The PASSED-gating live `/debrief` drive on a real split-root fixture, verified independently against a captured baseline, lands the new debrief as `…-04` in `{state_checkout}/_debriefs/` committed+pushed on `spacedock-state/dev` while `main` stays byte-identical — the exact end-value AC-1 specifies. One non-blocking note for the FO: the checklist's literal `-run DebriefSplitRoot` filter matches only 2 of the 4 test functions; the true four require `-run Debrief`.
