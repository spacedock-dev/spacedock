---
title: Boot identifies, engage converges — collapse the FO Startup recipe to ≤4 prose steps on the existing verbs
status: validation
source: "Captain, 0250 Commander session 2026-07-07, post k7 fo-boot-engage-split merge: 'i am allergic to the > 4 steps recipes. can we now further clean it up?' Startup in first-officer-shared-core.md is an 8-step prose recipe (version gate, project root, discovery, taxonomy read, state.boot, state.ensure-ready, state.sweep-merged, greet/headless). Steps 2-7 are deterministic orchestration the binary can own — the workflow's own prefer-code-gate-over-prose principle applied to its boot."
started: 2026-07-06T23:40:08Z
completed:
verdict:
score: 0.5
worktree: .worktrees/spacedock-ensign-startup-single-boot-call
issue:
id: 1y4ynffdxcgxn5eqcgw1mps3
sprint: 0250-fo-behavioral-discipline
mod-block: merge:pr-merge
pr: pr-merge:480
---

# Boot identifies, engage converges — collapse the FO Startup recipe to ≤4 prose steps on the existing verbs

## End value

The FO Startup recipe reads as **3 numbered prose steps** (down from 8), and the boot-resident `first-officer-shared-core.md` is **fewer bytes** than before the change — measured against the file as it stands when this implementation opens (post-vcm-merge). It ships **no new verb**: boot reuses the existing `status --boot` (modestly extended — discovery + taxonomy folded into the record, PR_STATE rendered as a local view), run as **local reads only** — no `gh`, no `state ready` pull, no sweep, no mutation — so a greet-only session changes nothing on disk. The per-workflow convergence (state ready + merged-PR sweep + live PR state) moves to the captain's first `«engage»(workflow)`, which invokes the existing `state ready` / `state sweep` **directly**, each per-class guard keeping its own distinct, remediation-carrying failure. The recipe gets shorter and the greet becomes side-effect-free, without the failure UX getting blurrier or a new binary surface to maintain.

## Problem

Startup is 8 prose steps; six of them (steps 2-7) hand-orchestrate work the binary already implements piecemeal. But those six split into two DIFFERENT classes, and the current recipe treats them the same:

**Steps 2-5 — local identify reads** (safe, deterministic, side-effect-free):
- step 2 — project root (`git rev-parse --show-toplevel`)
- step 3 — workflow discovery (`status --discover`)
- step 4 — stage-taxonomy read (`status --read README --json`)
- step 5 — `«state.boot»()` (`status --boot --json`)

**Steps 6-7 — per-workflow convergence** (network + mutation):
- step 6 — `«state.ensure-ready»()` (`state ready` — a `git pull --rebase`; resumes an absent checkout)
- step 7 — `«state.sweep-merged»()` (`state sweep` + the pr-merge advancement — `gh pr view` per PR-pending entity, then `status --set` / `archive` to terminalize a merged entity)

Two problems, not one. First, steps 2-5 are hand-issued reads with their own resident prose paragraphs — deterministic local identify the existing `status --boot` should serve in one record. Second, steps 6-7 run **eagerly at the greet, for every discovered workflow**, before the captain has acted on any of them — a `gh` round-trip and a state mutation spent on workflows a greet-only session will never touch. The convergence belongs where the captain acts (engage), not at the greet; the greet should be all reads and mutate nothing.

## Approach — extend the existing `status --boot` for identify; `«engage»(workflow)` owns convergence

Per the captain's Cycle-1 addendum: **no new `spacedock boot` verb.** The recipe collapses by folding the local reads into the EXISTING `status --boot` record and moving the convergence into `«engage»`, which invokes the EXISTING `state ready` / `state sweep` directly. The binary delta approaches zero — a modest `status --boot` extension; the deliverable is chiefly the contract follow-up. The version gate (step 1) stays a **distinct** step — an absent/too-old binary cannot self-report through any `status` verb.

**Ownership at a glance — boot identifies, engage converges:**

| Responsibility | Boot (`status --boot`, local reads, every session) | `«engage»(workflow)` (existing verbs, on first act) |
|---|---|---|
| Discover the managed workflow(s) | yes — lists zero / one / many, uniform | — |
| Stage taxonomy + entity labels | yes — folded into the record | — |
| MODS map / ID_STYLE / NEXT_ID / MIN_PREFIX / STATE_BACKEND | yes — local | — |
| ORPHANS (git/fs), TEAM_STATE (`~/.claude` probe) | yes — local reads | — |
| Dispatchable / ready-gate counts | yes — local counts, labeled possibly-stale | (refreshed as it acts) |
| PR state | local `pr:` mirror, labeled not-gh-checked | **live** — `gh pr view` |
| State-checkout convergence (pull / resume) | never | **yes** — `state ready` (exit 3 → halt, sweep not reached) |
| Merged-PR sweep + advance to terminal | never | **yes** — `state sweep` + pr-merge startup-hook |
| Network call / disk mutation | **none — provably side-effect-free** | yes, where convergence needs it |

Collapsed recipe (3 numbered steps):

1. **Binary version gate** — unchanged: resolve the launcher, `${SPACEDOCK_BIN:-spacedock} --version`, abort binary-absent / wrong-version each with its own remediation, do not proceed.
2. **Boot (local identify)** — `${SPACEDOCK_BIN:-spacedock} status --boot --json`, the existing verb modestly extended to **fold discovery + the stage taxonomy into the record** and render **PR_STATE as a local `pr:` mirror** (no `gh` at boot). Every part is a **local read** — filesystem, git-read, entity frontmatter, the `~/.claude` host probe. **No `gh`, no `state ready` pull, no sweep, no mutation**; a greet-only session writes nothing. Semantics are uniform across the discovered set (the N==1 eager special case is gone):
   - **zero discovery** — no managed workflow; report and STOP; do NOT broad-search the filesystem (code-gated by `detectBroadSearchAtBoot`).
   - **one or many** — a LIST of the discovered workflow(s) with per-workflow local counts. One workflow is a list of length 1; it gets **no** eager convergence.
   The greet's counts and PR fields are a **possibly-stale local view, labeled as such**, until the captain's first «engage».
3. **Interactive vs headless** — greet/drive from the record (current step 8, substantially as-is).

### `«engage»(workflow)` opens with that workflow's convergence — via the existing verbs, directly

Before the event loop, engage converges the NAMED workflow with the existing verbs invoked as **separate calls** — each per-class guard stays on its own call, unchanged:

- **`state ready`** — the split-root pull/resume (single-root is a no-op). On **exit 3** (a same-entity rebase conflict) → `«halt.rebase-conflict»(paths)`, and the sweep is NOT reached.
- **`state sweep`** — detect merged-but-not-terminalized entities; the pr-merge startup-hook advancement (a mutation) fires HERE, at first engage. Its **exit-0 `gh: "unavailable"` field** distinguishes real-empty from gh-unavailable UNKNOWN, never collapsed.
- **live PR state** — the pr-merge hook's `gh pr view`, filling in the OPEN / MERGED / CLOSED the greet's local mirror could not.

"Advanced at **boot** or never" becomes "advanced at **engage** or never"; a greet-only session (boot then stop) mutates nothing.

### Why the two guard channels stay on separate calls (the cycle-1 swallow problem dissolves)

The two convergence guards signal through **different channels**:
- `state ready` halt → **exit code 3** (`runStateReady`, `internal/cli/state_sync.go:138`, returns the int; proven live: `TestStateReadyHaltsOnBootConflict`, `TestStateReadyHaltStderrCarriesRemediationAndPeerCommit`).
- `state sweep` gh-unavailable → **exit-0 JSON field** `gh: "unavailable"` (`Sweep`, `internal/dispatch/reconcile.go:686`; proven live: `TestSweepGhUnavailableReportsUnknown`, `TestSweepGhPartiallyAvailableStillReportsNormally`).

Cycle 1 combined ready+sweep into ONE boot call, which forced a merge of the two channels and surfaced a swallow point (`Sweep` encodes its struct to stdout and returns only int, so a wrapper would have to re-parse stdout to recover the UNKNOWN — requiring a `sweepData` return-struct extraction). **The addendum's boundary dissolves that entirely:** because engage calls `state ready` and `state sweep` as two separate calls, each channel is read straight off its own call's exit code / stdout. There is no combine, so there is nothing to swallow — **no `sweepData` extraction, no new orchestration.** The two-channel analysis survives only as the reason the calls stay separate.

### Reconciling ORPHANS / PR_STATE / TEAM_STATE across the new boundary

- **ORPHANS** — worktree fields cross-referenced against filesystem + git state (`scanOrphans`). Local git-read, no network, no mutation → **stays in the boot record**, per discovered workflow.
- **TEAM_STATE** — the `~/.claude` host probe: is a team already present. Local read, no mutation → **stays in the boot record**.
- **PR_STATE** — the one that splits. Today `status --boot` calls `gh pr view` per PR-pending entity (`checkPRStates`, `internal/status/boot.go:86`) — a network call, and `live_prstate_pin_test.go` deliberately PINS that `pr_state.entries[].state` reflects live gh (or explicit unknown when gh is absent). So the FO's identify-boot must be gh-free WITHOUT breaking that pin: render the local mirror behind an **opt-in identify mode** (a new `--boot` flag the FO passes), leaving the un-flagged `status --boot` — and its live-PR pin — unchanged. In identify mode, PR_STATE is a **local mirror**: entities carrying a non-empty `pr:` field (read from frontmatter), listed by number, **labeled "local view — not gh-checked."** The live OPEN/MERGED/CLOSED state is filled in at engage's convergence. This is the change that makes the FO's boot record network-free while keeping every existing `--boot` consumer green (see the blast-radius check in the test plan).

### The prose-functions consolidate

`spacedock status --boot` / `state ready` / `state sweep` stay in the CLI (`--boot` gains the modest identify extension; `state ready` / `state sweep` are unchanged externally — tests and the idle pr-merge path use them). Only the FO's *recipe* changes:
- `«state.boot»` stays — its description folds in discovery + taxonomy and the local PR mirror; still shipped by `status --boot --json`.
- `«state.ensure-ready»` and `«state.sweep-merged»` **consolidate into `«engage»`'s convergence** — their content (the exit-3 halt gate, the sweep's UNKNOWN, "advanced at engage or never") moves into the `«engage»` effect and its `→ shipped` lines, rather than standing as two boot-resident sections.
- The four in-file cross-references (deferred-load-points note, write-scope carve-out, State-Management note, halt.rebase-conflict block) repoint from the boot-time `«state.*»` framing to `«engage»`'s convergence / the shipped verbs.

No deferred module (`fo-dispatch-core.md` / `fo-merge-core.md`) references these — the consolidation is contained to `first-officer-shared-core.md`.

## Scope

**In:** a modest `status --boot` extension (fold discovery + stage taxonomy into the JSON record; render PR_STATE as a local `pr:` mirror with NO `gh` call); the Startup 8→3 rewrite; the `«state.ensure-ready»`/`«state.sweep-merged»` → `«engage»`-convergence consolidation + cross-ref repoint; the pr-merge startup-hook timing note (fires at engage); the `command-reference.md` update.

**Out:** any NEW verb (`spacedock boot` is NOT introduced — the addendum reuses `status --boot`); the `sweepData` return-struct extraction (dropped — engage calls `state ready` / `state sweep` separately, so no combine, no swallow point); the `engage`/driver event-loop binary (roadmap 0222 — engage stays a prose interaction verb wrapping `«dispatch.next-action»()`); any change to `state ready` / `state sweep` external behavior (engage invokes them unchanged, so their own tests stay green).

## Acceptance criteria

**AC-1 (value — leaner recipe AND leaner file).** After the change the `## Startup` section carries **≤4 numbered prose steps** (from 8), AND `wc -c skills/first-officer/references/first-officer-shared-core.md` is **strictly less** than the same file measured immediately before this implementation's edits (the post-vcm-merge pre-change byte count, recorded in the implementation's first commit message). *Test:* a check counting top-level `^[0-9]+\.` items under `## Startup` (≤4) and asserting `post_bytes < recorded_pre_bytes`. Both halves must hold — a shorter recipe that grew the file fails.

**AC-2 (mechanism — the boot record folds in discovery + taxonomy + local PR mirror).** `status --boot --json` in identify mode on a healthy split-root workflow exits 0 and emits one JSON object carrying the discovery result, the `stages` taxonomy, every existing local boot section (MODS/ID_STYLE/NEXT_ID/MIN_PREFIX/ORPHANS/DISPATCHABLE/TEAM_STATE/STATE_BACKEND), and the **local PR mirror** (PR-pending entities by `pr:` number, labeled local/not-gh-checked). New keys are **appended** after the existing key set, preserving current order (per `json_boot_test.go`'s pinned key order). *Test:* a behavior fixture driving the binary and asserting the record's key set + order. (Mechanism AC — serves AC-1's value via the AC cross-check.)

**AC-3 (the identify-mode boot record is provably side-effect-free — the captain's "greet mutates nothing").** `status --boot --json` in identify mode makes **no network call** (no `gh` invocation) and performs **no mutation** (no `state ready` pull, no file write or commit to the state checkout). *Test:* a fixture running the identify-mode `--boot` with a recording/erroring `gh` probe and asserting it was never invoked, plus a before/after check that the state checkout's git HEAD and working tree are unchanged. This is the core boundary guarantee; it MUST have its own proof.

**AC-4 (version-gate classes unchanged).** binary-absent and wrong-version remain step-1 behavior. *Test:* the existing version-gate tests stay green.

**AC-5 (uniform zero/one/many — no N==1 eager convergence).** the boot record with no workflow reports "no workflow" and does NOT broad-search; with one workflow AND with multiple it emits a list, and in **neither** case does any convergence run (no `state ready`, no sweep — asserted the same way as AC-3). *Test:* three fixtures (empty root; one-workflow root → list of 1, no convergence side effect; two-workflow root → list of 2, no convergence side effect).

**AC-6 (engage converges — exit-3 halt propagates, sweep not reached).** the captain's first `«engage»(workflow)` runs `state ready` then `state sweep` before the loop; a same-entity rebase conflict makes `state ready` exit **3** with the halt remediation on stderr, and `state sweep` does NOT run past the halt. *Test:* the existing live rebase-conflict test (`TestStateReadyHaltsOnBootConflict`) stays green, and the `«engage»` contract prose orders ready-before-sweep with the exit-3 short-circuit.

**AC-7 (engage sweep UNKNOWN survives, not swallowed).** with a PR-pending entity and `gh` unavailable, `state sweep` at engage reports `gh: "unavailable"` (UNKNOWN), NOT "0 merged / empty". *Test:* the existing `TestSweepGhUnavailableReportsUnknown` stays green; because engage reads `state sweep`'s own exit-0 JSON field directly (no combine), the UNKNOWN cannot be swallowed.

## Spike (riskiest mechanism) — spiked, de-risked by source read

The addendum's boundary REMOVES cycle-1's riskiest mechanism (a one-call combine of two guard channels) rather than relocating it. What remains risky is a smaller claim: **can the boot record be made side-effect-free without losing useful greet counts?** Source read (`internal/status/boot.go`): `gatherBoot`'s ONLY network call is `checkPRStates` → `gh pr view` (line 86, called at line 184); every other section is a local read (`scanOrphans` git/fs, `computeDispatchable`, the `~/.claude` team probe, `os.Stat` for the backend). Dropping that `gh` call and rendering PR_STATE from the local `pr:` field is the single change that makes the record network-free. De-risked to "drop one gh call + relabel PR_STATE local + fold discovery/taxonomy into the record."

The two convergence guards are already proven live and stay UNCHANGED at engage — each on its own call, so no swallow point arises:

- exit-3 channel — `go test ./internal/cli/ -run 'TestStateReadyHaltsOnBootConflict|TestStateReadyHaltStderrCarriesRemediationAndPeerCommit'` → `ok`.
- exit-0 JSON-field channel — `go test ./internal/dispatch/ -run 'TestSweepGhUnavailableReportsUnknown|TestSweepGhPartiallyAvailableStillReportsNormally'` → `ok`.

Conclusion: no new orchestration and no struct extraction — the deliverable is the `status --boot` identify extension plus the contract follow-up. The spike seeds the implementation's first test: assert the extended `--boot` invokes no `gh` and writes nothing (AC-3), reusing the two green convergence suites above as the engage-side regression net.

## Test plan

- **The boot record is local-only, side-effect-free** (AC-2, AC-3): a behavior fixture drives `status --boot --json` against a healthy split-root workflow and asserts (a) the record's key set incl. discovery + taxonomy + the local PR mirror, (b) `gh` was never invoked (recording probe), (c) the state checkout's git HEAD + working tree are byte-identical before/after.
- **Uniform discovery, no N==1 eager convergence** (AC-5): three fixtures — empty root → "no workflow" + no broad-search; one-workflow → list of 1 + no convergence side effect; two-workflow → list of 2 + no convergence side effect.
- **Engage convergence guards stay green on their own calls** (AC-6, AC-7): the existing `TestStateReadyHaltsOnBootConflict` (exit-3 halt) and `TestSweepGhUnavailableReportsUnknown` (exit-0 UNKNOWN) stay green — engage invokes these verbs directly, so no new test is needed beyond confirming the `«engage»` contract orders ready-before-sweep with the exit-3 short-circuit.
- **Version gate** (AC-4): existing version-gate tests stay green (binary absent → 127 / wrong version → version-gate exit).
- **Value check** (AC-1): Startup numbered-step count ≤4 + byte delta strictly negative vs the recorded pre-change baseline.
- **Blast-radius check — existing `--boot` consumers provably unaffected** (the extension's diff surface is bounded): the un-flagged `status --boot` keeps its live-gh PR_STATE, so `internal/status/live_prstate_pin_test.go` stays green; new record keys (discovery/taxonomy) are **appended**, so `internal/status/json_boot_test.go`'s pinned key order + the native/text parity suites (`boot_probe_parity_test.go`, `zz_independent_parity_test.go`) stay green; the pr-merge mod does its OWN `gh pr view` over entity files and never parses `--boot` output, so its startup/idle/merge paths are untouched. The check: run the full `go test ./internal/status/... ./internal/dispatch/...` boot+sweep suites green with the extension, confirming no existing key, order, or live-PR assertion moved. Only the new identify-mode path adds tests (AC-2/AC-3/AC-5).
- **Cost/complexity:** Go unit + behavior-fixture tests (seconds); no live-workflow smoke needed — fixtures drive the binary directly. The one new code path is the `--boot` identify extension; the convergence verbs are unchanged, so their existing suites are the regression net.
- **CI lanes:** the shared core is host-neutral → 0250 §Required CI lanes apply: `claude-live` + `codex-live` + `pi-live` green before merge; a flake re-runs to green, never skipped.
- **Detached adversarial audit** (0250 §3): required before merge on a throwaway checkout — the shipped FO contract is one of docs/dev's high-stakes surfaces.

## Sequencing (0250 strict-serial tail)

`1y` edits the `## Startup` section of `first-officer-shared-core.md` — the exact file Wave 1's strict-serial chain (`k7→z25→zm→vcm`) serializes on (`docs/roadmap/0250-fo-behavioral-discipline/dispatch-sprint-execution.md` §2/§5: "do not open overlapping worktrees on `first-officer-shared-core.md`"). **The implementation worktree must not open until vcm's merge lands.** The NEGATIVE-byte baseline is captured from that post-vcm-merge file (NOT the v0.24.0 21,663-byte leanness pin), recorded in the implementation's first commit.

## Doc diff

### `skills/first-officer/references/first-officer-shared-core.md` — Startup (current 8 steps → 3)

The `**Launcher command invariant:**` preamble stays verbatim. The 8 numbered steps become 3:

```
1. **Binary version gate.** [current step 1 verbatim — resolve the launcher, `${SPACEDOCK_BIN:-spacedock} --version`, parse line 1, require minor 0.24; abort binary-absent and wrong-version each with its own remediation; do NOT proceed to boot.]
2. **Boot (local identify).** `${SPACEDOCK_BIN:-spacedock} status --boot --json` runs the whole pre-greet identify — project-root + workflow discovery, the stage taxonomy, and the local boot sections (MODS/ID_STYLE/NEXT_ID/MIN_PREFIX/ORPHANS/DISPATCHABLE/TEAM_STATE/STATE_BACKEND) — folded into one JSON record, as **local reads only** (no `gh`, no `state ready` pull, no sweep, no mutation). PR_STATE is a **local `pr:` mirror, labeled not-gh-checked**; live PR state is filled in at «engage». Semantics are uniform across the discovered set:
   - **zero discovery:** no managed workflow — report and STOP; do NOT broad-search the filesystem (code-gated by `detectBroadSearchAtBoot`).
   - **one or many:** a LIST of the discovered workflow(s) with their local dispatchable / ready-gate counts. One workflow is a list of length 1 — no eager convergence. NAME them in the greet; the captain converges and acts on one via «engage». Single-entity mode fails with an ambiguity error when many.
   The greet's counts and PR fields are a possibly-stale local view, labeled as such, until the first «engage»; boot writes nothing and touches no network.
3. **Interactive vs headless.** [current step 8 verbatim — greet/drive from the record; interactive summarizes managed workflow(s) + dispatchable/ready-gate counts and hints `Use engage <workflow>`, no gate render at greet; headless drives the loop, converging each workflow at its first engage, and authors full gate reviews at stops; headless+conn resolves gates.]
```

### `skills/first-officer/references/first-officer-shared-core.md` — «engage» gains the convergence; «state.boot» folds in identify; drop «state.ensure-ready» + «state.sweep-merged»

The `«engage»(workflow)` **effect** opens with that workflow's convergence before running `«dispatch.next-action»()`, invoking the existing verbs directly (add to the effect bullet):

```
- **effect — converge, then drive:** for the named `workflow` (default: the current / only managed workflow), FIRST converge its state with the existing verbs, each on its own call: `state ready` (the split-root pull/resume; single-root a no-op; on exit 3 → `«halt.rebase-conflict»(paths)` BEFORE the sweep), then `state sweep` (advance merged PRs to terminal — the pr-merge startup-hook fires HERE, at first engage, not the greet; "advanced at engage or never"; its exit-0 `gh: "unavailable"` field distinguishes real-empty from UNKNOWN, never collapsed) plus the live PR state (`gh`). THEN run `«dispatch.next-action»()` to its stopping condition.
- → **shipped** (converge): `` `spacedock state ready` `` then `` `spacedock state sweep` `` — two calls, each guard on its own; → **prose** (drive): the event loop wraps `«dispatch.next-action»()` (driver binary descoped to roadmap 0222).
```

`«state.boot»` stays; extend its effect to fold in discovery + taxonomy and the local PR mirror:

```
## «state.boot»(): read all local startup identify in one call

- **effect:** discover the workflow(s), read each stage taxonomy, and yield the boot record (MODS/ID_STYLE/NEXT_ID/MIN_PREFIX/ORPHANS/DISPATCHABLE/TEAM_STATE/STATE_BACKEND + the local `pr:` PR mirror) — **all local reads**, no `gh`, no `state ready` pull, no sweep, no mod-file open, no team creation, no mutation. Uniform across the discovered set: zero → no-workflow; one or many → a list. PR fields are the local `pr:` mirror, labeled not-gh-checked until «engage».
- **done-when:** the boot record is in hand for the greet; the greet has mutated nothing.
- → **shipped**: `` `spacedock status --boot --json` `` (extended to fold in discovery/taxonomy and render PR_STATE local).
```

DELETE the standalone `«state.ensure-ready»` and `«state.sweep-merged»` sections — their content now lives in `«engage»`'s convergence effect above. Repoint the four cross-refs off the boot-time `«state.*»` framing: the deferred-load-points note (`Startup step 8` → `step 3`; write-scope carve-out: the sweep/pr-merge advancement is pre-authorized **at engage**, not boot), the State-Management note (`Startup «state.ensure-ready»() fires before any dispatch` → `«engage»'s `state ready` fires before that workflow's loop`), and the `«halt.rebase-conflict»` block (`«state.ensure-ready»/«state.commit» exit 3` → `«engage»'s `state ready`/«state.commit» exit 3`). (`«state.commit»` and `«halt.rebase-conflict»` stay.)

### `docs/site/reference/command-reference.md` — Workflow table

No new verb row. Update the existing `status` entry (or add a `--boot` clause) so `status --boot` reads: "the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy, and reports the boot sections; PR_STATE is a local `pr:` view (live PR state is checked at engage). Local reads only, no mutation."

## Related

- Sits on k7 `fo-boot-engage-split` (the post-k7 light-greet / «engage» split this pushes convergence into).
- Serializes after `vcm` `fo-contract-keep-moving-posture` (shared-core strict-serial tail).
- Shares the "smallest sufficient mechanism" spirit with `zm` `fo-smallest-sufficient-mechanism` — the addendum picks no-new-verb, reuse-existing-verbs, precisely on that principle.
- Invokes `status --boot` (extended, `internal/status/boot.go`), `state ready` (`internal/cli/state_sync.go`) and `state sweep` / `Sweep` (`internal/dispatch/reconcile.go`) — the last two unchanged externally.

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
- Cycle 1 addendum (captain refinement, 2026-07-07): NOT a new `spacedock boot` verb. Properly follow k7 instead: boot = version gate + the EXISTING `status --boot` surface (at most modestly extended — discovery/taxonomy folded into the record; its live-gh PR_STATE polling deferred to engage or explicitly labeled a local view) + greet. `«engage»`'s convergence preamble invokes the EXISTING `state ready` / `state sweep` verbs directly, so each guard channel (exit-3 halt; gh-unavailable UNKNOWN) stays on its own call unchanged — the combined-verb swallow problem and the sweepData extraction drop out entirely. Binary delta approaches zero; the deliverable is chiefly the contract follow-up (Startup collapse + engage convergence + prose-function consolidation).

## Stage Report: ideation (cycle 2)

- DONE: Re-slice the Approach per Feedback Cycle 1's ownership boundary: boot = discover/identify only (no pull, no gh, no mutation; uniform zero/one/many; local-read counts labeled as local view); engage(workflow) opens with that workflow's convergence (ready + sweep + live PR state) before the loop; the N==1 eager special case is gone
  Approach rewritten: recipe collapses to 3 steps — version gate (step 1) + `status --boot --json` (extended local identify, PR_STATE a labeled local `pr:` mirror, no gh/no pull/no sweep/no mutation, uniform zero/one/many, N==1 eager case removed) + greet. Convergence (ready + sweep + live PR) moved into a new `«engage»(workflow)` convergence preamble. Honors the Cycle-1 ADDENDUM: no new `spacedock boot` verb — extend the existing `status --boot`; engage invokes existing `state ready` / `state sweep` directly.
- DONE: Reconcile every moved piece's contract: merged-PR sweep becomes advanced-at-engage-or-never; pr-merge mod startup-hook timing; ORPHANS/PR_STATE/TEAM_STATE placement; the «boot»/«engage» prose-function designs and all cross-refs updated coherently
  Reconciled: sweep + pr-merge startup-hook advancement fire at first «engage» (not the greet) — the mod file needs no timing edit (its `## Hook: startup` text is greet/engage-agnostic; the shared-core prose owns when startup hooks fire, consistent with the existing "startup hooks run deferred" framing). ORPHANS (local git/fs) + TEAM_STATE (local `~/.claude` probe) STAY in the boot record; PR_STATE splits — local `pr:` mirror at boot (drops the `checkPRStates` gh call, boot.go:86), live state at engage. «state.boot» stays (extended); «state.ensure-ready» + «state.sweep-merged» consolidate into «engage»; four cross-refs repointed in the Doc diff.
- DONE: Re-derive ACs + test plan for the new boundary: boot provably side-effect-free (fixture asserts no network/no state writes), engage-side convergence keeps every per-class guard (exit-3 halt, gh-unavailable UNKNOWN) with live-scenario proofs; AC-1's two-sided value stays (Startup <=4 numbered steps AND strictly-negative shared-core byte delta)
  AC-3 is now the side-effect-free proof (fixture asserts no gh invoked + state-checkout HEAD/tree unchanged). AC-6/AC-7 keep the exit-3 halt and exit-0 UNKNOWN guards via the existing green suites (`TestStateReadyHaltsOnBootConflict`, `TestSweepGhUnavailableReportsUnknown`) — engage reads each on its own call, so no combine and nothing to swallow (sweepData extraction dropped per the addendum). AC-1 unchanged (≤4 steps AND `post_bytes < recorded_pre_bytes`). AC-5 proves uniform zero/one/many with no N==1 convergence.

### Summary

Reworked the design to the captain's Cycle-1 ADDENDUM: no new verb — the recipe collapses to 3 steps by extending the existing `status --boot` (fold in discovery + taxonomy; PR_STATE a labeled local `pr:` mirror; drop the boot-time gh call so the greet is provably side-effect-free), and the per-workflow convergence (ready + sweep + live PR) moves to `«engage»(workflow)`, which invokes `state ready` then `state sweep` DIRECTLY. Because engage keeps the two verbs on separate calls, cycle-1's combined-verb swallow point dissolves: the exit-3 halt and exit-0 gh-unavailable UNKNOWN each read straight off their own call, so the `sweepData` return-struct extraction is dropped, not relocated. Binary delta approaches zero; the deliverable is chiefly the contract follow-up. Two flags for the gate: (1) this diverges from the team-lead dispatch framing that expected the sweepData extraction to survive/relocate — I followed the addendum as the captain's latest canonical word and messaged team-lead; (2) the `status --boot` PR_STATE change (live-gh → local) touches an existing verb's output — recommend a non-breaking path (a boot-mode gh-skip vs a wholesale drop) for the implementer to pick, since boot must be gh-free.

### Refinements (team-lead follow-up, addressed)

Three additional notes carried from the team-lead's confirmation, all applied:
- **Blast-radius check named in the test plan.** The PR_STATE change is now an **opt-in identify mode** (a new `--boot` flag), NOT a wholesale change — grounded by reading the consumers: `live_prstate_pin_test.go` deliberately pins un-flagged `--boot`'s live-gh PR_STATE, and `json_boot_test.go` pins the record key order/presence. So new keys append (order preserved) and the un-flagged `--boot` keeps live-gh; the pr-merge mod runs its own `gh pr view` over entity files and never parses `--boot`, so it is untouched. The test plan names the check (full `internal/status` + `internal/dispatch` suites green).
- **Title rename** — already applied to frontmatter by the captain ("Boot identifies, engage converges — …"); synced the H1 body heading to match.
- **Captain-legible ownership table** added at the top of Approach (boot identifies vs engage converges, row per responsibility), so the gate reads the boundary at a glance.
- Cycle 2 addendum (captain emphasis, 2026-07-07, pre-implementation): the boot/engage split is the sprint's principal PROSE-REMOVER — maximize the negative resident delta, do not merely satisfy AC-1's "strictly less." Specifically: (a) Startup step 5's nine boot-section sub-bullets must NOT survive as resident prose — the --boot JSON record is self-describing; any gloss worth keeping rides a lazy reference; (b) the «engage» convergence preamble stays terse — the shipped verbs' own stderr carries remediation, the contract does not restate it; (c) the implementation's stage report cites the achieved delta prominently against the pre-change measurement.
- Cycle 2 addendum 2 (captain-approved scope extension, 2026-07-07): 1y's implementation additionally carries a FROZEN, meaning-preserving dedup pass over the 0250 members' merged resident additions (FO editorial review, ~1,000-1,100 bytes; validated by 1y's own lanes at zero marginal cost, and counting toward its negative-delta AC). The frozen list — dedup ONLY, no guard clause's meaning changes: (1) «engage» block: state "not a binary / wraps existing logic" ONCE (keep the → line as canonical; trim the trigger restatement and the effect's fo-dispatch-core parenthetical that duplicates the Deferred-load-points list); replace the sprint-internal "a 0250 'Out of scope' extension" with "a named future extension" (~350B). (2) The no-render-at-greet rule is stated three times across step 8 and the deferred-load note — state once in step 8, minimal pointer elsewhere (~150B). (3) z25 S1 bullet 1: it states its rule forwards then backwards — keep the stronger back-half ("'unrelated' is a claim the change must substantiate"), cut the mirror (~120B). (4) zm blockquote: relocate the named-excuses list (including both "Ultracode" mentions — host-specific in a runtime-neutral core) to references/fo-smallest-sufficient-mechanism.md where AC-5's split promised it; the resident text keeps one short "never by a named excuse (see the reference)" clause (~260B). (5) present-gate rule 11: replace the not-a-pass enumeration duplicated verbatim from the shared-core S1 bullet with a pointer to it; keep rule 11's unique halves (the captain votes on which checks are green; a red is read from this run's evidence) (~200B).

## Stage Report: implementation

- DONE: The captain-approved cycle-2 design lands: status --boot gains the OPT-IN identify mode (un-flagged --boot behavior pinned by the existing live_prstate_pin/json_boot tests — keys additive, order preserved; pr-merge mod untouched), Startup collapses 8 -> 3 steps, «engage» opens with convergence via state ready THEN state sweep as SEPARATE calls (exit-3 short-circuit before sweep; gh-UNKNOWN read off sweep's own output), «state.ensure-ready»/«state.sweep-merged» consolidate into «engage», all four cross-refs repointed; AC-3's side-effect-free proof ships (no gh invocation, state-checkout HEAD/tree byte-identical)
  `--identify` flag threads dispatch→runRead→gatherBoot→checkPRStates/bootJSON (commit a847e472); un-flagged `--boot` byte-identical (live_prstate_pin/json_boot/native+text parity green UNMODIFIED). Startup 8→3, «engage» convergence + «state.*» consolidation + 4 cross-ref repoints (commit 132fbeb8). AC-3 green (boot_identify_test.go) + exercised on real docs/dev: identify boot exit 0, state checkout HEAD 3a4e4586 unchanged before/after.
- DONE: MAXIMIZE the negative resident delta (captain emphasis, Feedback Cycles addendum): Startup step 5's nine boot-section sub-bullets do NOT survive as resident prose (the JSON record self-describes; any keepable gloss rides a lazy ref); the convergence preamble stays terse; AND the frozen five-item dedup pass (addendum 2) executes exactly as listed, meaning-preserving — report achieved wc -c prominently vs pre-change 26,755 + 5,833
  Nine sub-bullets dropped; «engage» preamble terse (shipped verbs' own stderr carries remediation, not restated). Frozen 5-item dedup executed as listed (item 4 relocated to references/fo-smallest-sufficient-mechanism.md, which already carried the list — pure resident cut). Achieved: shared-core 26,755 → 24,722 (-2,033); present-gate 5,833 → 5,815 (-18); resident total 32,588 → 30,537 (**-2,051**). AC-1 green (Startup = 3 numbered steps AND strictly < 26,755).
- DONE: Every abort class keeps its distinct exit + stderr proof by fixture (zero discovery no-broad-search, multi lists-and-proceeds, exit-3 halt with sweep not reached, gh-unavailable UNKNOWN at engage); existing --boot consumer tests stay green UNMODIFIED as the blast-radius proof; full offline suite green; committed on the branch
  Zero (exit 1, "do NOT search the filesystem") + multi (lists 2, exit 0) proven by AC-5 fixtures (boot_identify_test.go) AND real-binary exercise. Exit-3 halt-sweep-not-reached + gh-unavailable UNKNOWN proven by the existing green suites UNCHANGED (engage keeps state ready / state sweep on separate calls, no combine): TestStateReadyHaltsOnBootConflict + TestSweepGhUnavailableReportsUnknown pass. `go test ./... -count=1` all green; committed on spacedock-ensign/startup-single-boot-call.

### Summary

Shipped the cycle-2 design: `status --boot` gains an opt-in `--identify` local-identify mode that folds workflow discovery + the stage taxonomy into the boot record, renders PR_STATE as a local `pr:` mirror (no gh), and makes no network call and no mutation — new keys append after the existing set so every un-flagged `--boot` consumer stays byte-identical (the blast-radius proof). The contract follow-up collapsed Startup 8→3 steps, moved the per-workflow convergence (`state ready` then `state sweep`, each on its own call) into «engage»'s preamble, consolidated «state.ensure-ready»/«state.sweep-merged» into «engage», and repointed all four cross-refs. The frozen 5-item dedup pass ran meaning-preserving. Resident delta is **-2,051 bytes** (32,588 → 30,537). Binary delta stayed modest as the captain required — one flag plus the fold; the exit-3 and gh-UNKNOWN guards were reused unchanged (the cycle-1 swallow point dissolved because engage keeps the two verbs on separate calls). Proven by exercising the real binary against docs/dev, not just fixtures.

Two implementer refinements to the ideation Doc-diff, flagged for the gate/validator: (1) the recipe + «state.boot» shipped line name `status --boot --identify --json` (the Doc-diff wrote `--boot --json` loosely; the AC and team-lead refinement pin the opt-in FLAG, so I made it precise); (2) «state.boot»'s effect does NOT re-enumerate the nine boot sections (the Doc-diff listed them, but cycle-2 addendum (a) says the JSON record self-describes and the nine sub-bullets must not survive as resident prose — I honored the addendum over the Doc-diff enumeration). The detached adversarial audit (0250 §3) and fresh validation remain the FO's pre-merge gate steps.

## Stage Report: validation

- DONE: Reproduce every AC independently (Startup steps <=4; resident delta exact; AC-2 key set; AC-3 side-effect-free; AC-4..7 abort classes; blast-radius pins UNMODIFIED — confirmed via diff)
  AC-1: Startup = 3 top-level numbered steps AND shared-core 24,722 < 26,755 — the 26,755 baseline confirmed INDEPENDENTLY via `git show aa01ae53:…|wc -c` (not the implementer's self-attested constant); present-gate 5,833→5,815; total 32,588→30,537 = **-2,051 exact**. AC-2/3/5 green (boot_identify_test.go); AC-4 version-gate suite green; AC-6 (TestStateReadyHaltsOnBootConflict) + AC-7 (TestSweepGhUnavailableReportsUnknown) green on UNTOUCHED files (internal/dispatch entirely unchanged). Blast-radius: live_prstate_pin_test.go + json_boot_test.go byte-untouched (diff-stat blank) and green; only edited status test files are boot_identify_test.go (new) + comment renumbering (no assertion weakened). Full offline `go test ./...` — every package ok, 0 failures.
- DONE: AC-3 proven by EXERCISING the real binary, not fixtures alone
  Built the binary; ran `status --boot --identify --json` against the real docs/dev workflow with a recording `gh` shim on PATH → exit 0, shim NEVER fired, state-checkout HEAD (a5e580d5) + working tree byte-identical before/after. Un-flagged `--boot` key-set diff adds ONLY discovery+stages, nothing reordered — the additive-key blast-radius proof, observed live.
- DONE: Dedup + design-fidelity audit — five frozen edits meaning-preserving; two implementer refinements assessed
  Enumerated each addendum-2 edit vs the diff: (1) «engage» "not a binary" stated once + trigger restatement & fo-dispatch-core parenthetical trimmed + "0250 Out of scope"→"named future extension"; (2) no-render-at-greet once in step 3 + minimal pointer in the deferred-load note; (3) z25 bullet keeps the stronger back-half, cuts the forward mirror; (4) zm named-excuses relocated — VERIFIED fo-smallest-sufficient-mechanism.md (untouched by 1y) already carries the full three-excuse list + Ultracode point at line 20, so it is a pure resident cut with NO information lost; (5) present-gate rule 11 pointerizes the not-a-pass enumeration into the shared-core self-evidence bar, keeps its two unique halves. Both implementer refinements (--identify flag spelling; «state.boot» dropping the nine-section enumeration) are faithful to the LATER captain word (team-lead opt-in-flag refinement; cycle-2 addendum (a)) — WITHIN the approved design, no route-back.
- DONE: Detached adversarial audit on a THROWAWAY checkout (high-stakes shipped contract + status surface)
  Detached worktree at HEAD 132fbeb8 (never the impl worktree; removed after). 4 claim-breaking edits, each caught RED by its owning AC, green on revert: A resurrect gh in identify → AC-3 + AC-2 FAIL; B zero-halt drops "do NOT search the filesystem" → AC-5 FAIL; C Startup 5 numbered steps → AC-1 FAIL ("has 5 … want <=4"); D halt exit 3→0 → AC-6 FAIL. Recorded as prose-only / not-test-gated BY DESIGN (surfaced, not blocking): the «engage» ready-before-sweep ORDERING is a prose interaction verb (driver binary descoped to roadmap 0222) — no test parses the effect; only the verb-level guards it composes (exit-3 via edit D; gh-UNKNOWN via AC-7) are gated. Dedup meaning-preservation is FO editorial (AC-1's byte floor catches only a large reversion) — hence the by-hand enumeration above.

### Summary

PASSED. All seven ACs reproduced with independent evidence: AC-1's byte baseline confirmed from git rather than the implementer's constant; AC-3 proven by driving the REAL binary against docs/dev (recording gh shim never fired, state HEAD/tree byte-identical); AC-2/4/5/6/7 green on fixtures and the existing untouched guard suites; the live_prstate_pin/json_boot blast-radius pins are byte-untouched and green. The five frozen dedup edits are meaning-preserving (item 4's relocation confirmed by the reference independently carrying the content), and both implementer Doc-diff refinements sit within the captain-approved cycle-2 design. The adversarial audit's four edits each went red under the owning AC and green on revert — the suite has teeth. Per the dispatch: claude-live/codex-live/pi-live remain PR-CI-at-merge and are NOT certified here (local benchmark-token expired 401; every validated AC is fixture/offline by design). One residual note for the gate: the «engage» ready-before-sweep ordering rests on prose, not a mechanical gate — acceptable given engage has no binary, but worth the captain's eye.
