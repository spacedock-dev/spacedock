---
title: Defer FO Write Scope + ID Styles out of the boot-resident core (sibling to entity-status)
status: ideation
sprint: 0240-lean-contract
group: cleanup
id: k408d2ydgj7s3s9yg7csyw81
started: 2026-06-30T16:20:13Z
---
Sibling to `entity-status`. The boot-resident FO core (`skills/first-officer/references/first-officer-shared-core.md`) loads ~1,027 tokens of write/filing-phase reference at greet that isn't needed until the FO writes or files: **FO Write Scope** (~696 tok — what the FO may write on main + the `spacedock new` atomic-create procedure, needed at write-time) and **ID Styles** (~331 tok — sd-b32/sequential/slug minting detail, needed at new-entity filing). Defer both into a lazily-loaded reference, loaded on first write / `--set` / `spacedock new`, using the same pointer/defer pattern as the dispatch/merge modules and entity-status.

Measured (`wc -c`/4 on `origin/main`, this ideation): FO Write Scope `## ` lines 180–193 = 2,827 chars ≈ 707 tok; ID Styles `## ` lines 72–80 = 1,341 chars ≈ 335 tok; combined ≈ **1,042 tok** ≈ 14% of the ~7,295-token boot core — confirms the captain's ~1,027 estimate. After this + entity-status, the remaining deferrable surface (Probe & Ideation Discipline ~347, Single-Entity ~107, Mod Hook Convention ~147) is a smaller third sibling if wanted. Sequencing: coordinate with z4 #6 (the registry block) and entity-status — all three edit first-officer-shared-core.md.

## Acceptance criteria
- **AC-1 (occupancy — the reason this exists)** — first-officer-shared-core.md drops by the two deferred sections (~1,042 tok measured: FO Write Scope ~707 + ID Styles ~335), measured **net-NEGATIVE vs `origin/main`**; the content lives in a new lazily-loaded `references/fo-write-core.md`, off the boot path. *Test:* `wc -c`/token delta of first-officer-shared-core.md vs `origin/main` after the cut is NEGATIVE (the same file-delta-vs-origin/main methodology the 0240 DoD mandates, not a prose-grep). The pointer block's own size is counted, so a bloated pointer moves the number the wrong way — an independent baseline that can regress.
- **AC-2 (greet/boot independence — the load-bearing guard)** — a greet-and-stop boot completes WITHOUT loading fo-write-core.md; a DISPATCHING boot loads it before the first frontmatter write / `spacedock new`; and the boot's own `«state.sweep-merged»` write path does **not** depend on fo-write-core.md. *Test:* (a) **sweep-independence — SPIKED below**: `spacedock state sweep` (read-only detect) + `spacedock merge guard --verdict passed` (binary terminalize+archive+commit) finalize a merged-PR entity in a directory carrying ZERO Write-Scope prose; pinned by `internal/cli/state_sweep_test.go::TestStateSweepIsReadOnly` + `internal/status/merge_guard_test.go`. (b) structural — the only mention of fo-write-core.md in the boot-resident bodies is the deferred-module pointer (with a greet guard), mirroring Dispatch/Merge, so the greet path (Startup 1–8 + `«state.*»` + `status --boot`) names no fo-write-core.md content.
- **AC-3 (reachability — no guard lost)** — every FO Write-Scope rule and every id-style minting path resolves at its trigger via the boot-core pointer; **no write-permission guard is lost** in the move. *Test:* (a) the move is a **relocation, not a rewrite** — a `diff` of the moved sections against `origin/main` proves byte-preservation (no rule dropped); (b) `internal/contractlint/boot_resident_closure_test.go::TestHostNeutralCoresResolveAndCarryCeremony` extended with `fo-write-core.md → ["## FO Write Scope", "## ID Styles"]` enforces, via the filesystem oracle (not a prose-grep), that the shared core NAMES the reference AND the reference CARRIES both sections; `TestBootResidentDeferredLoadPointsResolve` already stats the new `references/fo-write-core.md` pointer; (c) the cross-ref inventory (Coordination, below) confirms every "see FO Write Scope" / id-style pointer resolves to fo-write-core.md with no dangling reference.

## Design: the cut

Same name-a-pointer / defer-the-body pattern as the `## Dispatch (deferred module)` and `## Merge and Cleanup (deferred module)` sections (first-officer-shared-core.md:90–94, 134–138). **Locate sections by heading, NOT line number** — sibling `84` removes `## Status Viewer` (lines 35–70) from the same file and shifts everything below it; the implementation re-anchors on `## ID Styles` and `## FO Write Scope`.

**Edit 1 — new file `skills/first-officer/references/fo-write-core.md`** (host-neutral, deferred). Header + the two sections moved **verbatim** (byte-preserved, per AC-3a):

```markdown
# First Officer Write Core (host-neutral)

The FO's main-branch write-authority boundary, the `spacedock new` atomic-create
procedure, and new-entity id-style minting. Lazily loaded at the first write to main
(named by the boot-resident core's `## FO Write Scope and ID Styles (deferred module)`
pointer); a greet-and-stop boot never reads it. The host `spacedock new` invocation form
is the runtime adapter's new-entity binding, read alongside this file.

## FO Write Scope
{verbatim move of shared-core's `## FO Write Scope` section}

## ID Styles
{verbatim move of shared-core's `## ID Styles` section}
```

(Order: Write Scope first — the file is named for it. ID Styles' existing intra-section "(see FO Write Scope)" pointer at old line 80 then resolves intra-file.)

**Edit 2 — remove `## ID Styles` (shared-core 72–80) and `## FO Write Scope` (shared-core 180–193)** from first-officer-shared-core.md.

**Edit 3 — add the boot-resident pointer** where `## FO Write Scope` was (the write-phase part of the file, after the `«state.*»` functions):

```markdown
## FO Write Scope and ID Styles (deferred module)

- → **runtime-binding**: `references/fo-write-core.md` (host-neutral) — the FO's main-branch write-authority boundary, the `spacedock new` atomic-create procedure, and id-style minting — loaded at the FIRST write to main (first `status --set`, `spacedock new`, archive move, or `### Feedback Cycles` write). The host `spacedock new` invocation form is the runtime adapter's already-resident new-entity binding.
- **done-when:** the FO is about to write entity state or file a new entity on main — the write-authority boundary and the id-style minting detail are resident.
- **guard:** a greet-and-stop boot does not read it; the boot's own `«state.sweep-merged»` write is binary-executed (`spacedock state sweep` detect + `spacedock merge guard` finalize) and depends on no part of this reference.
```

**Edit 4 — fix the boot-resident cross-refs that named the moved sections** (before → after):
- shared-core State Management `- The FO owns YAML frontmatter on the main branch (see FO Write Scope below).` → `… on the main branch (full write-authority scope in the deferred write reference; see the FO Write Scope and ID Styles pointer below).`
- shared-core Status Viewer `--next-id` bullet `… use `spacedock new` (see FO Write Scope), …` → `… use `spacedock new` (see `references/fo-write-core.md`), …` *(this bullet is inside `## Status Viewer`, which `84` defers — coordination below)*.
- `claude-first-officer-runtime.md:35` `(see `## FO Write Scope` in the shared core for the full contract)` → `(see `references/fo-write-core.md` for the full contract)`.
- `skills/feedback-rejection-flow/SKILL.md:23` `Routing follows FO Write Scope:` → `Routing follows the FO write reference (`references/fo-write-core.md`):` (the worktree-vs-main rule is restated inline there, so this is a name-accuracy fix, not a load-bearing dependency).

**Edit 5 — register the reference in the contractlint structural guard** (the AC-3b gate): add to `internal/contractlint/boot_resident_closure_test.go`'s `foReferenceCores` map: `fo-write-core.md → {"## FO Write Scope", "## ID Styles"}`. Optionally add fo-write-core.md to `deferred_tier_absence_test.go`'s `foContractCores` (hygiene; not load-bearing).

## Riskiest-mechanism spike (AC-2, load-bearing) — DONE

**Claim under test:** deferring `## FO Write Scope` out of the boot core must NOT break the boot's own `«state.sweep-merged»` write path (terminalize verdict=PASSED / archive / remove-worktree at boot when a merged PR exists). The risk: if that path read the Write-Scope prose to authorize itself, deferral would silently break boot.

**Result — the sweep write is binary-self-authorizing, not prose-authorized.** I built the contract-2 binary (`spacedock dev (contract 2)`) and exercised the two boot-sweep commands in a fixture directory containing **zero** Write-Scope prose (only the workflow fixture; no `skills/`):
- `spacedock merge guard 080-pr-merged --verdict passed` (the entity carries `pr: pr-merge:99`, empty `mod-block`, status `implementation` — the exact boot-sweep precondition) → `finalized: 080-pr-merged -> done (verdict passed), archived.` On-disk after: `status: done`, `verdict: passed`, `completed:` stamped, file moved to `_archive/`, and the binary made its own commit `archive 080-pr-merged (merge guard)`. **The terminalize+archive+commit ran with no FO Write Scope section anywhere on disk.**
- `spacedock state sweep --json` on a fresh fixture copy → JSON envelope `"swept": []`, HEAD unchanged, clean working tree. **Detection is read-only; the actual write is `merge guard`.**

**Trace confirming the chain stays resident/binary after the cut:** Startup step 7 → `«state.sweep-merged»()` (boot-resident, first-officer-shared-core.md:165–171, NOT in either deferred section) → reads `_mods/pr-merge.md` → its startup hook delegates the finalize to `spacedock merge guard {slug} --verdict passed` (mod line 20–23). `«state.sweep-merged»` binds `→ shipped: spacedock state sweep`. The shipped command (`internal/cli/state_sync.go::runStateSweep`) is read-only ("makes no commit, push, or mutation"); the write is `MergeGuard` (`internal/status/merge.go`). **No node in this chain names `## FO Write Scope`.** The Write-Scope section declares what the FO may *hand-write* on main; the boot sweep's writes are binary-executed and need no such declaration loaded.

**No spike needed for the greet-no-load pattern itself:** it reuses the already-proven Dispatch/Merge deferred-module guard (a greet-and-stop boot traverses Startup 1–8 + `«state.*»` only, which name no deferred-module load point). The only NEW risk was the boot's own write path — spiked above.

## Test plan

| AC | What verifies it | Kind / cost |
|---|---|---|
| AC-1 | `git diff --stat origin/main -- first-officer-shared-core.md` shows net-NEGATIVE char/token delta (≈ −1,042 tok); fo-write-core.md is off the boot path so it doesn't count against boot occupancy. | measurement / trivial |
| AC-2 | (a) `TestStateSweepIsReadOnly` + `merge_guard_test.go` (existing, green) pin the binary self-authorizing sweep — re-run + the reproducible bare-dir spike above. (b) structural: grep confirms the sole boot-resident mention of `fo-write-core.md` is the deferred-module pointer. | Go unit (existing) + grep / low |
| AC-3 | (a) `diff` moved sections vs `origin/main` = byte-identical relocation. (b) `TestHostNeutralCoresResolveAndCarryCeremony` (extended) + `TestBootResidentDeferredLoadPointsResolve` go green = shared core names the reference, it carries both anchors, the pointer resolves on disk. (c) cross-ref inventory: 0 dangling "see FO Write Scope". | Go contractlint + diff / low |

Implementation is a contract-file + Go-test edit → a dispatched worker under contractlint in a worktree. **No live workflow lane required**: behavior is unchanged (the moved bodies say exactly what the prose said; the boot sweep is already binary and already tested). `go build ./...` + `go test ./internal/contractlint/ ./internal/cli/ ./internal/status/` is the gate.

## Coordination

- **Sibling `84` (entity-status) — SAME file, different sections.** `84` defers `## Status Viewer` (35–70) + `## Issue Filing`; I defer `## ID Styles` (72–80) + `## FO Write Scope` (180–193). Disjoint sections, but `84` shifts my line numbers — re-anchor on headings. Wave 2 order is **`84` then `k4`** (roadmap Sequencing), so `84` lands first. One shared edge: the `## Status Viewer` `--next-id` bullet's "(see FO Write Scope)" pointer (old line 47) — when `84` moves it into the entity-status reference, the pointer must end up as "(see `references/fo-write-core.md`)". Whichever of us lands the dependent edit, the gate's cross-ref inventory (AC-3c) catches a dangling pointer.
- **`z4` #6 (fn-binding-refinements) — Wave 3, follows me.** `z4` replaces the two pointer prose sections with one `«fn»`-registry block that "indexes the deferred modules (Dispatch, Merge, and the new Status-Viewer / Write-Scope references)." My `## FO Write Scope and ID Styles (deferred module)` section is one of the entries `z4` will fold into that registry. No conflict — `z4` runs last and consumes my pointer; I add the pointer in the existing prose-section shape and `z4` restructures it.
- **Contractlint test files** (`boot_resident_closure_test.go`) are also edited by `84` (to register its entity-status reference). Same-file, different map entries, sequenced — `84` adds its entry first, I add `fo-write-core.md` after.

## Stage Report: ideation

- DONE: The LOAD-BEARING AC-2 check, spiked — trace the boot's `«state.sweep-merged»` write path and PROVE it does NOT depend on the deferred FO Write Scope reference; exercise greet (no load) vs first-dispatch (load-before-first-write) on fixtures.
  Built contract-2 binary; ran `merge guard 080-pr-merged --verdict passed` in a dir with ZERO Write-Scope prose → terminalized+archived+committed; `state sweep --json` → read-only (HEAD unchanged). Sweep is self-authorizing via shipped `state sweep` (detect) + `merge guard` (write) + boot-resident Startup step 7 / `«state.sweep-merged»` (core lines 165–171), NOT the deferred prose. Spike section in body; pinned by `state_sweep_test.go` + `merge_guard_test.go`.
- DONE: The deferred content named exactly and moved to a lazily-loaded reference reached by a boot-core pointer at first write / `--set` / `spacedock new`; occupancy delta measured net-NEGATIVE.
  FO Write Scope (~707 tok, core 180–193) + ID Styles (~335 tok, core 72–80) → new `references/fo-write-core.md`; boot-resident `## FO Write Scope and ID Styles (deferred module)` pointer (`→ runtime-binding`, greet guard). AC-1 measures char/token delta vs `origin/main` (Design Edits 1–3).
- DONE: Reachability / no guard lost — every Write-Scope rule + ID-style minting path resolves at its trigger via the pointer, no write-permission guard lost; record coordination with `84` and `z4`.
  6-site cross-ref inventory in Coordination + Design Edit 4; AC-3b extends `TestHostNeutralCoresResolveAndCarryCeremony` (`fo-write-core.md → ["## FO Write Scope","## ID Styles"]`) as the filesystem-oracle guard; move is byte-preserved relocation (AC-3a). `84` (same file, different sections) + `z4` (registry) coordination recorded.

### Summary

Refined the ideation design in place (kept Problem/approach/AC-1..3, added Design doc-diff, the load-bearing spike, test plan, coordination). The riskiest path — the boot's own `«state.sweep-merged»` write — was SPIKED and proven binary-self-authorizing (`state sweep` + `merge guard` finalize a merged-PR entity with zero Write-Scope prose on disk), so deferring FO Write Scope cannot break boot. Design follows the Dispatch/Merge deferred-module pattern: a new host-neutral `fo-write-core.md` + a boot-core pointer; AC-3 reachability is gated by an extended `boot_resident_closure_test.go` (filesystem oracle), not a prose-grep. Coordination noted with sibling `84` (same file, disjoint sections, lands first) and `z4` (registry indexes the new pointer, runs last).
