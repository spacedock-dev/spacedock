---
title: "Collapse the FO Startup recipe to <=4 prose steps behind one shipped boot call"
status: ideation
source: "Captain, 0250 Commander session 2026-07-07, post k7 fo-boot-engage-split merge: 'i am allergic to the > 4 steps recipes. can we now further clean it up?' Startup in first-officer-shared-core.md is an 8-step prose recipe (version gate, project root, discovery, taxonomy read, state.boot, state.ensure-ready, state.sweep-merged, greet/headless). Steps 2-7 are deterministic orchestration the binary can own — the workflow's own prefer-code-gate-over-prose principle applied to its boot."
started: 2026-07-06T23:40:08Z
completed:
verdict:
score: 0.5
worktree:
issue:
id: 1y4ynffdxcgxn5eqcgw1mps3
sprint: 0250-fo-behavioral-discipline
---

# Collapse the FO Startup recipe to <=4 prose steps behind one shipped boot call

## End value

The FO Startup recipe reads as **3 numbered prose steps** (down from 8), and the boot-resident `first-officer-shared-core.md` is **fewer bytes** than before the change — measured against the file as it stands when this implementation opens (post-vcm-merge). It ships a real `spacedock boot` verb that runs the whole pre-greet sequence in one call, and **every per-class abort keeps its own distinct, remediation-carrying failure**. The recipe gets shorter without the failure UX getting blurrier.

## Problem

Startup spends six of its eight prose steps (steps 2-7) hand-orchestrating deterministic reads the binary already implements piecemeal:

- step 2 — project root (`git rev-parse --show-toplevel`)
- step 3 — workflow discovery (`status --discover`)
- step 4 — stage-taxonomy read (`status --read README --json`)
- step 5 — `«state.boot»()` (`status --boot --json`)
- step 6 — `«state.ensure-ready»()` (`state ready`)
- step 7 — `«state.sweep-merged»()` (`state sweep`)

One `git` call plus five `spacedock` calls the FO issues by hand every boot, each with its own resident prose paragraph and (for steps 5-7) its own `«state.*»` prose-function section further down the file. It is the workflow's own **prefer-a-code-gate-over-a-prose-rule** principle applied to its own boot: deterministic orchestration belongs in the binary, not in six paragraphs the model re-reads each session.

## Approach — one `spacedock boot` verb absorbs steps 2-7

A new `spacedock boot --json` runs the whole pre-greet sequence internally and emits **one** JSON boot record: workflow discovery + stage taxonomy + the existing `--boot` sections (MODS / ID_STYLE / NEXT_ID / MIN_PREFIX / ORPHANS / PR_STATE / DISPATCHABLE / TEAM_STATE / STATE_BACKEND) + the state-ready outcome + the sweep outcome. The version gate (step 1) stays a **distinct** step — it must run before any other verb (an absent/too-old binary cannot self-report through `boot`), which is why the checklist reads "absorbing steps 2-7," not 1-7.

Collapsed recipe (3 numbered steps):

1. **Binary version gate** — unchanged: resolve the launcher, `${SPACEDOCK_BIN:-spacedock} --version`, abort binary-absent / wrong-version each with its own remediation, do not proceed.
2. **Boot** — `${SPACEDOCK_BIN:-spacedock} boot --json`; consume the record as JSON. Abort by class on a non-zero exit, each remediation on its own stderr — do not blur them:
   - **zero discovery** (exit 4) — no managed workflow; report and STOP; do NOT broad-search the filesystem (still code-gated by `detectBroadSearchAtBoot`).
   - **split-root rebase-conflict** (exit 3) — `«halt.rebase-conflict»(paths)` per the stderr remediation.
   - **multiple workflows is NOT an abort** — the record lists them; the greet lists and proceeds (the captain picks one later via «engage»).
   The record carries the sweep outcome (real-empty vs `gh`-unavailable UNKNOWN, never collapsed) and the state-ready result.
3. **Interactive vs headless** — greet/drive from the record (current step 8, substantially as-is).

### The two guard channels the combined verb must not swallow (the crux)

The riskiest part: the two boot-time guards signal through **different channels**, and a naive combine loses one.

- **`state ready` halt → exit code 3.** Propagate by short-circuit: `boot` runs ready and, on exit 3, returns exit 3 *before* the sweep runs, so `«halt.rebase-conflict»` still fires. `runStateReady` already returns the int → in-process propagation is trivial. (Proven live: `TestStateReadyHaltsOnBootConflict`, `TestStateReadyHaltStderrCarriesRemediationAndPeerCommit`.)
- **`state sweep` gh-unavailable → JSON field, exit 0.** `Sweep` returns exit **0** for BOTH real-empty and gh-unavailable; the only distinguisher is the `gh: "unavailable"` / `reason` field. Exit-code inspection alone loses it — `boot` must **merge the sweep outcome struct** into the record. (Proven live: `TestSweepGhUnavailableReportsUnknown`, `TestSweepGhPartiallyAvailableStillReportsNormally`.)

**Concrete implementation constraint the spike surfaced:** `gatherBoot` already returns `*bootData` (mergeable) and `runStateReady` returns its int, but `Sweep` (`internal/dispatch/reconcile.go:686`) builds its `sweepResult` inline and *encodes it to stdout*, returning only the int. A `boot` wrapper calling `Sweep` would get exit 0 and have to re-parse stdout to recover the UNKNOWN — the swallow point. The implementation must expose the sweep outcome as a **returned struct** (extract `sweepData(workflowDir, gh) (sweepResult, int)`; `state sweep` serializes it, `boot` merges it). That is the smallest change that lets one call preserve both channels.

### Standalone verbs stay; the prose-functions collapse

`spacedock state ready` / `state sweep` / `status --boot` remain in the CLI unchanged (tests and the idle pr-merge mod path use them); only the FO's *boot recipe* stops calling them by hand. So the three shared-core prose-functions «state.boot» / «state.ensure-ready» / «state.sweep-merged» collapse into one `«boot»()` section, and the four in-file cross-references (deferred-load-points note line 45, write-scope carve-out line 48, State-Management note line 113, halt.rebase-conflict block line 139) repoint to «boot». No deferred module (`fo-dispatch-core.md` / `fo-merge-core.md`) references these — the collapse is contained to `first-officer-shared-core.md`.

### Multi-workflow resolution (design decision — recommended)

`boot` acts on the discovered set: zero → exit 4; one → full record (taxonomy + state + ready + sweep); **multiple → list the workflows with per-workflow dispatchable/gate counts and DEFER the ready/sweep convergence to the captain's first «engage»** on a chosen workflow. Converging N state checkouts eagerly at a greet that will act on one is exactly the expensive deferral k7 already pushes past the greet. Recommending defer-per-engage over eager-all-N; flag for the gate.

## Scope

**In:** the `spacedock boot` verb (discovery + taxonomy + boot sections + ready + sweep, classed aborts); the `sweepData` return-struct extraction; the Startup 8→3 rewrite; the «state.*»→«boot» prose-function collapse + cross-ref repoint; the `command-reference.md` entry.

**Out:** the `engage`/driver binary (roadmap 0222); multi-workflow eager convergence; any change to `state ready` / `state sweep` / `status --boot` external behavior (boot reuses them, so their own tests stay green unmodified).

## Acceptance criteria

**AC-1 (value — leaner recipe AND leaner file).** After the change the `## Startup` section carries **≤4 numbered prose steps** (from 8), AND `wc -c skills/first-officer/references/first-officer-shared-core.md` is **strictly less** than the same file measured immediately before this implementation's edits (the post-vcm-merge pre-change byte count, recorded in the implementation's first commit message). *Test:* a check counting top-level `^[0-9]+\.` items under `## Startup` (≤4) and asserting `post_bytes < recorded_pre_bytes`. Both halves must hold — a shorter recipe that grew the file fails.

**AC-2 (mechanism — one call, all sections present).** `spacedock boot --json` on a healthy split-root workflow exits 0 and emits one JSON object carrying the discovery result, the `stages` taxonomy, every existing `--boot` section, the state-ready result, and the sweep outcome. *Test:* a behavior fixture that drives the binary and asserts the record's key set. (Mechanism AC — serves AC-1's value via the AC cross-check.)

**AC-3 (version-gate classes unchanged).** binary-absent and wrong-version remain step-1 behavior (boot does not move these classes). *Test:* the existing version-gate tests stay green.

**AC-4 (zero vs multi discovery, distinct outcomes).** `boot` with no workflow exits **4**, names "no workflow" on stderr, and does NOT broad-search; `boot` with multiple lists them and exits **0** (not an abort). *Test:* two fixtures (empty root; two-workflow root) asserting exit code + stderr/record.

**AC-5 (split-root halt propagates as exit 3, sweep not reached).** a same-entity boot rebase conflict makes `boot` exit **3** with the halt remediation on stderr, and the sweep does NOT run past the halt. *Test:* a live rebase-conflict fixture asserting exit 3 + remediation (models `TestStateReadyHaltsOnBootConflict`).

**AC-6 (sweep UNKNOWN survives the merge, not swallowed).** with a PR-pending entity and gh unavailable, `boot` exits **0** and the record's sweep field reads `gh: "unavailable"` (UNKNOWN), NOT "0 merged / empty". *Test:* a fixture injecting an erroring gh probe asserting the record's sweep field (models `TestSweepGhUnavailableReportsUnknown`).

## Spike (riskiest mechanism) — spiked, de-risked to a known small refactor

Riskiest claim = one call propagating **both** guard channels without swallowing either. Exercised in this ideation (all green):

- exit-3 channel — `go test ./internal/cli/ -run 'TestStateReadyHaltsOnBootConflict|TestStateReadyHaltStderrCarriesRemediationAndPeerCommit'` → `ok`.
- exit-0 JSON-field channel — `go test ./internal/dispatch/ -run 'TestSweepGhUnavailableReportsUnknown|TestSweepGhPartiallyAvailableStillReportsNormally'` → `ok`.
- read the sources: `gatherBoot → *bootData` (mergeable), `runStateReady → int` (propagatable), `Sweep` encodes its struct to stdout and returns only int (**the swallow point**, requiring the `sweepData` return-struct extraction above).

Conclusion: the combine is sequencing + one struct extraction, not new orchestration invention. The spike seeds the implementation's first two tests — drive `boot` through a conflict fixture (assert exit 3 *before* sweep) and through an erroring-gh fixture (assert record UNKNOWN at exit 0).

## Test plan

- **Each abort class by a live scenario observing exit code + stderr (never prose-grep):**
  - binary absent → 127 / wrong version → version-gate exit (step 1, existing tests, AC-3).
  - zero discovery → boot exit 4 + "no workflow" stderr (empty-root fixture, AC-4).
  - multi discovery → boot exit 0 + workflow list in record (two-workflow fixture, AC-4).
  - split-root rebase-conflict → boot exit 3 + halt remediation, sweep NOT reached (conflict fixture, AC-5).
  - sweep gh-unavailable → boot exit 0 + record `gh: "unavailable"` (erroring-gh fixture, AC-6).
- **Value check** (AC-1): Startup numbered-step count ≤4 + byte delta negative vs the recorded pre-change baseline.
- **Cost/complexity:** Go unit + behavior-fixture tests (seconds); no live-workflow smoke needed for the verb — fixtures drive the binary directly. The `sweepData` extraction is a refactor with existing sweep tests as the regression net.
- **CI lanes:** the shared core is host-neutral → 0250 §Required CI lanes apply: `claude-live` + `codex-live` + `pi-live` green before merge; a flake re-runs to green, never skipped.
- **Detached adversarial audit** (0250 §3): required before merge on a throwaway checkout — the shipped FO contract is one of docs/dev's high-stakes surfaces.

## Sequencing (0250 strict-serial tail)

`1y` edits the `## Startup` section of `first-officer-shared-core.md` — the exact file Wave 1's strict-serial chain (`k7→z25→zm→vcm`) serializes on (`docs/roadmap/0250-fo-behavioral-discipline/dispatch-sprint-execution.md` §2/§5: "do not open overlapping worktrees on `first-officer-shared-core.md`"). **The implementation worktree must not open until vcm's merge lands.** The NEGATIVE-byte baseline is captured from that post-vcm-merge file (NOT the v0.24.0 21,663-byte leanness pin), recorded in the implementation's first commit.

## Doc diff

### `skills/first-officer/references/first-officer-shared-core.md` — Startup (current lines 5-34 → below)

The `**Launcher command invariant:**` preamble stays verbatim. The 8 numbered steps become 3:

```
1. **Binary version gate.** Before boot, run `${SPACEDOCK_BIN:-spacedock} --version` and parse line 1: `spacedock <version>`. These skills require binary minor 0.24 (patch/prerelease skew fine). Abort by class:
   - **Binary absent or non-executable** — [current step 1 bullet 1 verbatim: retry-once-with-bare, install/source-build hint, no `doctor`].
   - **Binary present but wrong version** — [current step 1 bullet 2 verbatim: too-old/too-new/dev-integer, ABORT + `${SPACEDOCK_BIN:-spacedock} doctor`].
   In every class, do NOT proceed to `boot`.
2. **Boot.** `${SPACEDOCK_BIN:-spacedock} boot --json` runs the whole pre-greet sequence — project-root + workflow discovery, stage-taxonomy read, the boot record (MODS/ID_STYLE/NEXT_ID/MIN_PREFIX/ORPHANS/PR_STATE/DISPATCHABLE/TEAM_STATE/STATE_BACKEND), split-root convergence, and the merged-PR sweep — and emits one JSON boot record on exit 0 (consume as JSON; the human table is not rendered for FO reasoning). Abort by class on a non-zero exit, each remediation on its own stderr — do not blur them:
   - **zero discovery (exit 4):** no managed workflow — report and STOP; do NOT broad-search the filesystem (code-gated by `detectBroadSearchAtBoot`).
   - **split-root rebase-conflict (exit 3):** `«halt.rebase-conflict»(paths)` per the stderr remediation.
   - **multiple workflows:** NOT an abort — the record lists them; NAME them in the greet and proceed (the captain acts on one later via «engage»). Single-entity mode fails with an ambiguity error.
   The record carries the sweep outcome — real-empty vs `gh`-unavailable UNKNOWN, never collapsed — and the state-ready result; the pre-greet boot is all reads plus the one convergence pull, no mod-file open, no team creation.
3. **Interactive vs headless.** [current step 8 verbatim — greet/drive from the record; interactive summarizes managed workflow(s) + dispatchable/ready-gate counts and hints `Use engage <workflow>`, no gate render at greet; headless drives the loop and authors full gate reviews at stops; headless+conn resolves gates.]
```

### `skills/first-officer/references/first-officer-shared-core.md` — «state.*» prose-functions (current lines 117-131 → below)

Replace the three sections `«state.boot»` / `«state.ensure-ready»` / `«state.sweep-merged»` with one:

```
## «boot»(): run the whole pre-greet sequence in one call

- **effect:** discover the workflow(s), read the stage taxonomy, yield the boot record (the Startup step 2 sections consumed as JSON), converge the split-root checkout (the halt-gate + the one pull-on-boot; single-root is a no-op), and sweep merged PRs to terminal — all reads plus the one convergence pull, no mod-file open, no team creation. Each sub-outcome is a record field: the state-ready result, and the sweep outcome distinguishing real-empty from `gh`-unavailable UNKNOWN (named per its registered startup-hook mod).
- **done-when:** the boot record is in hand for the greet.
- **block:** a split-root rebase-conflict short-circuits to exit 3 → `«halt.rebase-conflict»(paths)` BEFORE the sweep runs (the halt is not swallowed by continuing); zero discovery is exit 4.
- → **shipped**: `` `spacedock boot --json` ``.
```

Then repoint the four cross-refs: line 45 `«state.boot» JSON` → `the «boot» record`; line 48 `boot «state.sweep-merged»/pr-merge advancement` → `boot «boot» sweep/pr-merge advancement`; line 113 `Startup «state.ensure-ready»() step` → `Startup «boot» convergence`; line 139 `«state.ensure-ready»/«state.commit» exit 3` → `«boot»/«state.commit» exit 3`. (`«state.commit»` and `«halt.rebase-conflict»` stay.)

### `docs/site/reference/command-reference.md` — Workflow table

Before (row): `| `spacedock state` | Manage a [split-root workflow]…'s state checkout (`state init` resumes one on a fresh clone, `state new` births one) |`

After: add one row above it —
`| `spacedock boot` | Run the whole first-officer pre-greet sequence in one call (discover the workflow, read the stage taxonomy, the boot record, converge split-root state, sweep merged PRs) and emit one JSON record; classed non-zero exits for no-workflow / split-root halt. The first officer's Startup step 2. |`

## Related

- Sits on k7 `fo-boot-engage-split` (the post-k7 light-greet / «engage» Startup this collapses).
- Serializes after `vcm` `fo-contract-keep-moving-posture` (shared-core strict-serial tail).
- Reuses `state ready` (`internal/cli/state_sync.go`) and `Sweep` (`internal/dispatch/reconcile.go`) unchanged externally; extracts `sweepData` from the latter.

## Stage Report: ideation

- DONE: Design the single boot verb absorbing Startup steps 2-7 with EVERY per-class abort preserved as a distinct, remediation-carrying failure; value AC measures BOTH the recipe dropping 8 -> <=4 numbered prose steps AND a NEGATIVE shared-core resident byte delta
  Approach + Doc diff design `spacedock boot --json` (steps 2-7) with the version gate kept as step 1; aborts mapped to distinct channels — binary-absent/wrong-version (step-1 version gate), zero-discovery (exit 4), multi (list, not an abort), split-root halt (exit 3), sweep gh-unavailable UNKNOWN (exit-0 JSON field). AC-1 pairs step-count ≤4 with `post_bytes < recorded_pre_bytes`; measured headroom: Startup section 7,202 B + 3 collapsible «state.*» sections 1,066 B collapse to ~3 steps + 1 «boot» section.
- DONE: Spike-or-record the riskiest mechanism first per proof policy — whether the combined verb can propagate the per-command guard semantics through one call without swallowing them
  Spiked green: exit-3 halt channel (`TestStateReadyHaltsOnBootConflict` + remediation test) and exit-0 JSON-field UNKNOWN channel (`TestSweepGhUnavailableReportsUnknown` + partial-availability test); source read found the swallow point — `Sweep` encodes its struct to stdout and returns only int, so the combine needs a `sweepData` return-struct extraction (`gatherBoot`/`runStateReady` already return usable values). De-risked to sequencing + one struct extraction.
- DONE: Test plan proves each abort class by a live scenario observing exit code + stderr (never prose-grep), includes the docs-site diff for any command-surface change, and pins the sequencing constraint
  Test plan: one fixture per abort class asserting exit code + stderr/record (AC-3..AC-6) plus the AC-1 value check; docs-site diff drafted for `command-reference.md` (new `spacedock boot` row); Sequencing section pins "implementation worktree must not open until vcm's merge lands" per 0250 dispatch-sprint-execution §2/§5, with the NEGATIVE-byte baseline captured post-vcm-merge.

### Summary

Designed `spacedock boot --json` to absorb Startup steps 2-7, collapsing the recipe 8→3 steps while keeping the version gate as step 1 and each abort class a distinct, remediation-carrying failure. The load-bearing finding, spiked live: the two boot-time guards signal through different channels — `state ready`'s halt via exit code 3, `state sweep`'s gh-unavailable UNKNOWN via a JSON field at exit 0 — so the combine must short-circuit on the exit code AND merge the sweep struct; `Sweep` currently returns only its int (encoding its struct to stdout), so the implementation must extract `sweepData` as a returned struct to avoid swallowing the UNKNOWN. One open design point flagged for the gate: multi-workflow boot defers ready/sweep convergence to the captain's first «engage» (recommended over eager-all-N). Sequencing is strict-serial after vcm's shared-core merge.

### Feedback Cycles

- Cycle 1 (captain design revision, 2026-07-07 — FO gate approval RESCINDED, routed back to ideation): re-slice the boot/engage boundary. `boot` = identify only — version-gate-adjacent discovery + per-workflow taxonomy/labels + counts from LOCAL reads, no `state ready` pull, no gh, no sweep, no mutation; semantics uniform across zero/one/many workflows (zero -> abort; else a list — kill the N==1 eager special case). `«engage»(workflow)` opens with that workflow's convergence (state ready + merged-PR sweep + live PR state) before the event loop — "advanced at boot or never" becomes "advanced at engage or never"; a greet-only session mutates nothing. Accepted trade: greet counts/PR fields are a possibly-stale local view, labeled as such, until first engage. Reconcile the pr-merge mod's startup-hook timing and the ORPHANS/PR_STATE/TEAM_STATE section placement with the new boundary.
