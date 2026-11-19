# 0240 lean-contract — Commander dispatch (cold-boot execution package)

> **Self-contained cold-boot package.** A Commander session boots `spacedock:first-officer` on Claude, reads this file + the sprint `index.md` + `staff-review.md`, and drives the **8 members** to the 0240 Definition of Done, then cuts **0.24.0**. The Shaping-FO phase is done (ideate → independent multi-lens staff review → fold all 10 Material findings → ideation gates approved); this is the handoff to the drive phase. Approved by the captain 2026-06-30 after the preflight staff review (SOUND-WITH-FIXES) and the fold.

## 0. Boot

1. Guard PATH before any git-shelling spacedock call: `export PATH="/usr/bin:/bin:/usr/local/bin:/opt/homebrew/bin:$PATH"`.
2. Resolve the launcher at the version gate: `SPACEDOCK_BIN` (set/executable) else `spacedock` on PATH. The contract-2 binary is required (`spacedock --version` → `contract 2`; `origin/main` carries #443). Build from source if needed: `go build -o spacedock ./cmd/spacedock`.
3. Workflow dir: `docs/dev` (split-root; entity state on `spacedock-state/dev`, checked out at `docs/dev/.spacedock-state`). `«state.ensure-ready»` (pull-on-boot) before any dispatch.
4. Members query: `${SPACEDOCK_BIN:-spacedock} status --workflow-dir docs/dev --where sprint=0240-lean-contract` → the 8 members (all `status: ideation`, gate-approved, implementation-ready).

## 1. Gate status — do NOT re-present ideation gates

**Ideation gates are APPROVED** (captain, 2026-06-30) after the independent multi-lens preflight staff review (`staff-review.md`, verdict SOUND-WITH-FIXES) and the fold of all 10 Material findings into the member bodies + the carve. Each member's body carries the gate-approved design (problem / approach / measured ACs / test plan / spike). The Commander does **NOT** re-present ideation gates; it drives `implementation → validation → done` per member.

Every member carries a **measured** end-value AC (file-delta vs `origin/main`, a behavioral count, or on-disk state) — validation must MEASURE it, never assert it.

## 2. Members, waves, and the strict-sequence constraint

The contract-file members edit the **same 2–3 files**, so **implementation dispatches run in strict sequence, never in parallel** (parallel = merge conflicts in the same paragraphs). Implementation order:

- **Predecessor (parallel-safe, different surface): `f5` journeymetrics-ensign-read-adoption** — folds the dispatched-ensign sub-agent transcript into the journeymetrics `--read` metric AND normalizes launcher recognition (counts the canonical quoted `${SPACEDOCK_BIN:-spacedock}`). Edits `internal/journeymetrics` (Go) — no contract-file collision. **MUST land (impl + validation) before `82k`'s AC-1 validation** (82k's adoption before/after is vacuous until f5 + its launcher widening land).
- **Wave 1 — trim the cores:** `82k` read-guidance-redundant-with-grep, then `scr` ensign-contract-dev-leakage. Both edit `skills/ensign/references/ensign-shared-core.md` (disjoint sections; coordinate). NOTE: `82k` also edits `first-officer-shared-core.md` (its Site C `:229` Probe/Read bullet + Site D `:100` gate-verdict read) — it is the **first toucher of the FO shared core**, and Site D (`:100`) is adjacent to z4 #4's target — z4 must re-anchor on bullet text, not line numbers.
- **Wave 2 — defer reference out of the boot core:** `84` entity-status, then `k4` defer-write-scope-id-styles. Both edit `first-officer-shared-core.md` (disjoint sections) AND `internal/contractlint/boot_resident_closure_test.go`. **`84` lands first** (registers `fo-status-viewer.md` in `foReferenceCores` — first edit to that test file), then `k4` (registers `fo-write-core.md` + adds the prose-pointer oracle `TestDeferredReferenceProsePointersResolve`).
- **Wave 3 — last:** `z4` fn-binding-refinements. Its #6 four-entry registry `os.Stat`s the `references/*.md` files `84`+`k4` create — naming a Wave-2 file before it exists REDS the closure walk, so z4 **must** implement after Wave 2.
- **Independent (any time, different file): `7d4` debrief-split-root-state-home** — edits `skills/debrief/SKILL.md`. No collision.
- **Parallel-safe (any time, different file): `48g` ac-scan-value-annotation-skip** — one-line matcher broadening in `internal/status/gate_extract.go`. No collision; fixes the `--ac-scan` value-annotated-AC bug (so the deterministic gate cross-check sees the value AC).

## 3. Per-member build notes (folded staff-review corrections — these are already in the bodies; restated here as drive cues)

- **`f5`** — AC-2 recognizes launchers by NORMALIZING (strip quotes; canonicalize the `SPACEDOCK_BIN` family quoted-or-not), not a literal-token allowlist; add the quoted positive to `TestCommandInvokesStatusRead`. `$SPACEDOCK` is legacy, NOT recognized.
- **`82k`** — 4 trim sites incl. FO `:100` / Site D (confirmed in scope). AC-1 (adoption-not-regress) is a **live FO+ensign before/after** and is **gated on `f5`** landing; the dispatch goldens stay byte-identical (negative control — the dispatch hint was already removed in #392). Pre-agreed degradation if the small-count metric is inconclusive at the gate: ship the trim, treat the ensign-aware metric as a forward-looking non-regression check (captain ratified Wave 1 stays).
- **`scr`** — scope is the −49-char stage-name-parenthetical removal + a NEW `contractlint` structural-absence guard (with a discriminator control). #290 prior art (PASSED/archived) — NOT a 0240 deliverable. **Before implementation: re-title** (current title is a misnomer → "Stage-neutralize the ensign core + add a regression guard") **and rename the slug** to clear the `--resolve` collision with archived `ep0ra3z` (the `--where sprint=` membership query is unaffected).
- **`84`** — AC-3 is content-carriage: register `fo-status-viewer.md → ["## Status Viewer","### Captain-Facing State Display","## Issue Filing"]` in `foReferenceCores` (the `## Issue Filing` anchor proves the issue-filing approval gate isn't lost). Arrow-taxonomy: `84`+`k4`'s adapter-less refs use one label — `z4` classifies it (`→ reference` proposed vs `→ runtime-binding`) and owns the `runtime-support.md` sync.
- **`k4`** — boot-sweep is binary-self-authorizing via `state sweep` (read-only detect) + the pr-merge hook's `status --set`/`status --archive` finalize (pinned to `verdict_guard_test.go`/`archive_guard_test.go`, NOT `merge_guard_test.go`). Edit 4 converts the moved `(see FO Write Scope)` pointer (now inside `fo-status-viewer.md`) to the path form; Edit 6 adds `TestDeferredReferenceProsePointersResolve` (spiked RED-then-GREEN) over both new reference files.
- **`z4`** — AC-8 is the **four-entry** registry (Dispatch, Merge, Status-Viewer, Write-Scope); prove all four load-points + four greet-guards survive via an explicit structural grep (the closure walk only checks filename resolution). AC-4 is two-part: z4-own per-file delta ≤ 0 vs `origin/main`; absolute ≤ v0.22.0 only for files z4 materially restructures (`claude-fo-dispatch.md` EXEMPT — already +508 B over from prior work). AC-6 needs **one live re-run** of `TestLiveReanchorGateRejectsMeansOnlyRegressed` (`-tags live`); the rest is offline. Classify the arrow label + sync `runtime-support.md`.
- **`7d4`** — consume `entity_dir`/`entity_dir_present` from `status --boot --json` (not re-parsing `state:`); halt-gate on `entity_dir_present == false` (AC-4 + fixture). Route reads/write/commit to `{state_checkout}/_debriefs/` via path-scoped git (NOT `spacedock state commit`, which is entity-only). Live `/debrief` drive gates the `done` stage (skill change).
- **`48g`** — one regex at `gate_extract.go:52`: `\*\*(AC-[0-9A-Za-z]+)\*\*` → `\*\*(AC-[0-9A-Za-z]+)[^*]*\*\*`. RED-then-GREEN fixture + the real-entity 7→8 baseline (z4's `**AC-4 (value guardrail)**` surfaces). + the `docs/dev/README.md:131` doc diff.

## 4. Drive procedure (per member)

1. **Implementation** (worktree stage): the deliverable is **`skills/` scaffolding / Go code** built by the dispatched ensign under contractlint — do NOT hand-edit. Commit in the worktree; entity body + stage report to the split-root state checkout.
2. **Validation** (`fresh: true` — a fresh validator): MEASURE every AC's end-value (file delta vs `origin/main`, the behavioral count, on-disk state); reproduce each "Verified by" clause.
3. **Detached adversarial audit** — REQUIRED before merging the contract/scaffolding changes (the shipped FO/ensign contract + scaffolding is one of the four high-stakes surfaces). Run on a throwaway checkout, not the impl worktree. `48g`/`f5` (Go, low-blast-radius) and `7d4` (skill) take the routine path unless validation flags risk.
4. **Required CI lanes are a function of the diff:** a change to the shipped FO/ensign contract or a host adapter (`skills/**/references/**`) REQUIRES the matching live lane green before merge; a flake is re-run to green, never skipped. `z4`/`48g`/`f5` are offline-provable except `z4` AC-6's one live re-run and `82k`/`7d4`'s live drives.
5. **Merge** to `main` (PR-merge), then advance to `done`. Keep state commits concurrency-safe (path-scoped; rebase-conflict halt).

## 5. Release cut — 0.24.0

After all 8 members are merged to `main`:
1. `go test ./...` green from the repo root.
2. **Pre-cut antipattern audit** (independent reviewer over the assembled sprint) before the tag fires; ship-blockers fixed pre-cut, non-blockers recorded for the next sprint.
3. **⚠️ e2e-gate trap (`q3`, not yet fixed):** the release e2e-gate's `gh run list --status success` only matches *first-attempt* successes, so a Live-E2E run **re-run to green** (run_attempt ≥ 2) is invisible and blocks the cut — it bit the v0.23.0 cut. **Either land a clean first-attempt-green Live-E2E run for the tagged SHA, or fix `q3` first** (drop `--status success` from `ghRunListForCommit`).
4. Confirm contract-2 consistency (the tagged commit's manifest brackets `CONTRACT_VERSION == 2`).
5. Stamp the dev/plugin manifests → tag so the tagged commit's manifest matches its tag → publish → bump the Homebrew cask → advance the `stable` branch. Authoritative procedure: `docs/releasing.md`; consult the runbook entities `steady-state-stable-release-runbook`, `stamp-then-tag-release-ritual`, `next-independent-release-line`. *(Captain authorizes the cut.)*

## 6. Out of cohort — `0qt` (separate usability mustfix, NOT a 0240 member)

`0qt` dispatch-build-team-name-advisory was driven separately by the Shaping FO. Status: **implementation complete + verified clean**, committed in worktree `.worktrees/spacedock-ensign-dispatch-build-team-name-advisory` (branch `spacedock-ensign/dispatch-build-team-name-advisory`, off `origin/main` `e3f85ec3`, commit `d833cd5a`). It adds a stderr advisory when `--team-name` is passed on host=claude (success-path placement — fires only when the legacy envelope is emitted, verified by a two-arm test) + documents `--host`/`--team-name`/`--bare-mode` in `dispatch build --help`. Ready for a **fresh validator** (validation is `fresh: true`) → merge. Independent of the cohort; can ship before, with, or after 0.24.0.

## 7. Close (Shaping FO)

After the cut: fold the pre-cut audit's deferred findings into the next sprint's backlog; a light post-cut release verification (some release-machinery issues only manifest when the tag fires). 0.25.0 candidates held: `fo-opus-behavioral-robustness`, `minor-version-compat-coupling`.
