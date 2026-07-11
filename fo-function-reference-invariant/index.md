---
title: Replace mutable step-number references with named FO functions
status: ideation
score: 0.95
source: "Captain direction 2026-07-11: widen the step-number sweep; hooks and references use the «fn» notation."
id: 88tq5zyg9jvx13f33zz3eq28
started: 2026-07-10T23:57:47Z
---

## Problem

The active first-officer contract still addresses shared procedures by mutable list positions. A step insertion can silently redirect `Startup step 3`, `Dispatch step 2`, Event Loop step 0/3, Merge-and-Cleanup step 10, reuse-condition numbers, or a legacy recovery tier. This coupling conflicts with `docs/runtime-support.md`: shared behavior should expose named `«fn»` capabilities, and runtime adapters should bind those names.

A read-only spike applied the proposed address classifier to the fixed active surface below. It found **51 mutable numeric address occurrences across 10 files**: Claude dispatch 17, dispatch core 11, shared core 8, legacy helper 4, merge core 3, Claude runtime 2, Codex runtime 2, Pi runtime 2, present-gate 1, and feedback-rejection 1. The same pattern ignores local ordered-list markers and semantic numbers such as feedback cycle 3, exit code 3, AC-2, version 0.24.0, and issue numbers. This is the independent baseline for the task's end value: implementation must reduce 51 to 0 without deleting local execution order.

## Active lint boundary

Use an explicit path registry, not a recursive `skills/**/*.md` ban. The invariant applies to the first-officer entry point, its three host adapters and host-neutral cores, and the helpers those cores directly load for gates, rejection, legacy/recovery, and hook execution:

- `skills/first-officer/SKILL.md`
- `skills/first-officer/references/{first-officer-shared-core,fo-dispatch-core,fo-merge-core}.md`
- `skills/first-officer/references/{claude-first-officer-runtime,claude-fo-dispatch,codex-first-officer-runtime,pi-first-officer-runtime}.md`
- `skills/{present-gate,feedback-rejection-flow,using-legacy-claude-team,fo-dispatch-recovery,fo-write-core}/SKILL.md`

This boundary excludes unrelated skills whose local numbered recipes are not first-officer cross-file APIs. Inside the boundary, leading ordered-list markers remain legal. A line such as `3. If nothing is dispatchable` defines local order; prose such as `re-run step 3` is a mutable address and fails.

## Numeric-reference inventory

The table maps every current address family to one existing or proposed binding. Line numbers are the pre-change baseline.

| Mutable address | Active sites | Named replacement |
|---|---|---|
| `Startup step 2` | shared core 34, 130, 162; present-gate 45 | Existing `«state.boot»()`; the boot record owns discovery, labels, mods, and summary inputs. |
| `Startup step 3` / bare `step 3` for launch behavior | shared core 34, 44; Claude runtime 21 | New `«interaction.boundary»()`; it owns interactive greet-stop, headless drive, and given-the-conn behavior. |
| `Dispatch step 2` | Claude dispatch 46 | New `«dispatch.checklist»(entity, stage)`; move the ≤3-linchpin checklist rules into this single definition and call it from Dispatch. |
| Event Loop `step 0` / `step-0` | dispatch core 112, 147, 152; Claude dispatch 108, 110, 120, 124 | Existing `«roster-reconcile»()`; rename the Claude block as its binding and call the function at first dispatch, idle, and post-merge. |
| Event Loop `step 0.5`, `steps 1-3`, `step 1`, and idle `step 3` | dispatch core 82, 147, 157; Claude dispatch 110, 112, 120, 124 | Existing `«dispatch.next-action»()`; keep its local 0.5/1/2/3 ordered skeleton, but describe cross-line branches by action or function name. |
| Merge-and-Cleanup `step 10` / `Step 10` | merge core 3, 18; Claude runtime 13; Claude dispatch 112, 124 | Existing `«worker.shutdown»()`; remove the duplicate numbered teardown definition and bind Claude's terminal teardown under the capability name. |
| `reuse-condition-0` / `reuse condition 0` | dispatch core 107; Claude dispatch 98; Codex 25; Pi 13; feedback-rejection 18 | Existing `«context-budget»()`. |
| `reuse-condition-1` | dispatch core 82; Claude dispatch 36; Codex 20 | Existing `«addressable-worker»`. |
| `reuse-condition-4` | dispatch core 46, 91, 93; Claude dispatch 102; Pi 16 | Existing `«reuse.model-match»`. The local reuse-condition list may stay numbered. |
| completion signal `(1, 2, or 3 above)` | Claude dispatch 64 | Existing `«completion-signal»`; the local three-item signal list remains ordered. |
| legacy `tier-1`, `tiers 1 and 2` | using-legacy 49, 50, 52 | New `«legacy-team.recover»()`; it owns the ordered fresh-name → degraded-mode → captain ladder. Cross-references name the branches, not tier numbers. |
| legacy teardown `step 3` | using-legacy 63 | Existing `«worker.shutdown»` legacy override; refer to its interactive settle branch. |
| `entry-point principle 1/3` | shared core 159, 161 | Delete the numeric parentheticals. The surrounding bullets already state the complete, non-callable principles; inventing functions would duplicate semantics. |

Local list definitions stay numbered: Startup 1-3, Dispatch 1-9, reuse conditions 0-4, Event Loop 0.5/1-3, present-gate 1-11, feedback rejection 1-7, and legacy recovery/teardown lists. Their callers use names.

## Lifecycle-hook inventory

Introduce one host-neutral function, `«hooks.run»(point)`, in the shared core's current Mod Hook Convention. It owns discovery of registered mods for `startup`, `idle`, or `merge`, alphabetical execution, and the rule that boot only reads the registry. Bind it as `→ prose`; it changes no binary or runtime tool surface.

| Current hook reference | Replacement |
|---|---|
| shared-core engage's startup-hook advancement | `«hooks.run»("startup")` after the separate read-only `state sweep`, on first `«engage»` |
| dispatch-core idle branch | `«hooks.run»("idle")` inside `«dispatch.next-action»()` |
| merge-core merge hook and startup/idle detection prose | `«hooks.run»("merge")`, `«hooks.run»("startup")`, and `«hooks.run»("idle")`; retain `«roster-reconcile»` as the non-hook sweep path |
| merge-core armed/mod-block continuation | `«merge.guard»` names `«hooks.run»("merge")` as the next action |
| fo-write-core's statement that the FO runs hooks | call `«hooks.run»(point)`; the write-core still forbids direct mod edits |

### Hook invocation order

`«hooks.run»(point)` replaces the existing calls; it does not add a second invocation.

| Point | Exact preserved boundary |
|---|---|
| `startup` | On the named workflow's first `«engage»`, run `state ready`, then the separate read-only `state sweep`, then invoke `«hooks.run»("startup")` exactly once. The registered startup mod retains its documented behavior, including its own `gh` scan and advancement; no sweep result is passed into it. Boot/greet never invokes it. |
| `idle` | Inside the no-dispatch branch of `«dispatch.next-action»()`, after the first `status --next` returns empty, invoke `«hooks.run»("idle")` exactly once, then `«roster-reconcile»()` when present, then the second `status --next`. |
| `merge` | Invoke `«hooks.run»("merge")` only when `«merge.guard»` returns its existing armed result and names the registered merge mod. Then re-run `«merge.guard»`. Blocked and finalized results do not run the merge hook. |

`state sweep` remains read-only, and `«roster-reconcile»()` remains a non-hook convergence path. The refactor renames each current invocation at the same boundary and order.

## Approaches considered

1. **Recommended: fixed-scope syntax classifier plus named-function closure.** Scan the explicit active paths for numeric address syntax, compare all required `«fn»` calls with one owner definition/binding, and include planted fail/pass controls. This directly measures the defect, permits local ordering, and avoids interpreting prose meaning.
2. **Build a Markdown procedure graph.** Parse headings, ordered lists, and cross-file links, then infer whether each number points outside its list. This is more general, but it adds a parser and semantic inference for a one-time contract normalization. It can still misclassify prose references.
3. **Ban `step N` across all skills.** This is short but over-broad. It would reject unrelated local recipes, miss reuse conditions, legacy tiers, and parenthesized numeric pointers, and invite a growing allowlist.

## Recommended design

Keep numbered lists only where local ordering aids execution. Any cross-line, cross-section, cross-file, hook, or runtime-override reference targets a named `«fn»` capability.

- Consolidate Startup's boot semantics under existing `«state.boot»()` and move launch routing under new `«interaction.boundary»()`. Startup's local list calls those functions instead of duplicating their bodies.
- Move Dispatch's checklist rules into new `«dispatch.checklist»(entity, stage)`; both normal and break-glass dispatch call it.
- Keep Event Loop's local ordered skeleton, but address it only as `«dispatch.next-action»()` and its host pre-pass only as `«roster-reconcile»()`.
- Replace every terminal teardown pointer with existing `«worker.shutdown»()` and name each host binding or legacy override explicitly.
- Replace reuse-condition numbers with their existing capability names: `«context-budget»`, `«addressable-worker»`, and `«reuse.model-match»`.
- Define `«hooks.run»(point)` once and use it for startup, idle, and merge lifecycle calls.
- Define `«legacy-team.recover»()` once for the legacy recovery ladder. Keep the ladder ordered locally, but name its branches in cross-references.

Representative wording changes are literal:

| Before | After |
|---|---|
| `boot record (Startup step 2)` | `boot record from «state.boot»()` |
| `Startup step 3's headless rule` | `the headless branch of «interaction.boundary»()` |
| `[CHECKLIST — assemble from Dispatch step 2]` | `[CHECKLIST — assemble via «dispatch.checklist»(entity, stage)]` |
| `re-run the host's step-0 reconcile sweep` | `invoke «roster-reconcile»()` |
| `Merge-and-Cleanup step 10` | `«worker.shutdown»()` |
| `reuse-condition-0/1/4` | `«context-budget»` / `«addressable-worker»` / `«reuse.model-match»` |
| `Fire idle hooks` | `invoke «hooks.run»("idle")` |

### Required call-site matrix

The durable structural test checks each token in the bounded section shown; a token elsewhere cannot satisfy the row.

| Function | Required owner/caller sections |
|---|---|
| `«state.boot»()` | owner: shared `## «state.boot»`; callers: shared `## Startup`, `## Deferred load points`, `## Mod Hook Convention`, and `## Working Principles`; present-gate `### Captain-facing assembly rules` |
| `«interaction.boundary»()` | owner: shared named heading; callers: shared `## Startup`, `## Deferred load points`, and `## Single-Entity Scope`; Claude runtime `## Captain Interaction` |
| `«dispatch.checklist»(entity, stage)` | owner and normal caller: dispatch core `## Dispatch`; fallback callers: Claude dispatch reuse-advance template and fo-dispatch-recovery `## Break-Glass Manual Dispatch` (replace `{numbered checklist}` with the named function's output) |
| `«dispatch.next-action»()` | owner/callers: dispatch core `## Event Loop`; Claude `«roster-reconcile»` binding and Backstop use it for branch timing instead of step numbers |
| `«roster-reconcile»()` | owner: dispatch core named heading; callers: dispatch-core `«dispatch.next-action»`; Claude named binding and Backstop |
| `«worker.shutdown»()` | owner: dispatch core named heading; callers: merge core `«merge.guard»`; Claude runtime `## Terminal teardown` pointer; concrete Claude dispatch section renamed/marked `Claude binding: «worker.shutdown»`; Claude reconcile/Backstop; legacy section renamed/marked `Legacy override: «worker.shutdown»` |
| `«context-budget»()` | owner: dispatch core named heading; callers: dispatch reuse list, Claude/Codex/Pi runtime binding sections, and feedback-rejection flow |
| `«addressable-worker»` | owner: dispatch core named heading; callers: dispatch reuse list, Claude Worker Back-Channel, and Codex runtime binding |
| `«reuse.model-match»` | owner: dispatch core named heading; callers: dispatch reuse list, Claude Context Budget, and Pi runtime binding |
| `«completion-signal»` | owner: dispatch core named heading; caller: Claude `## Awaiting Completion` decision procedure |
| `«hooks.run»(point)` | owner: shared `## «hooks.run»`; callers: shared `«engage»`, dispatch `«dispatch.next-action»`, merge `«merge.guard»`/Mod-Block, and fo-write-core mutation scope |
| `«legacy-team.recover»()` | owner: legacy named heading; callers: the legacy setup block and failure branches |

The same test enforces body ownership. `Interactive`/`Headless` branch headings occur only inside `«interaction.boundary»`; checklist constraint bullets occur only inside `«dispatch.checklist»`; lifecycle ordering occurs only inside `«hooks.run»`; legacy recovery's 1-3 ladder occurs only inside `«legacy-team.recover»`. Caller sections contain the named call but no second copy of those structural blocks.

The implementation changes instruction structure and contractlint only. It adds no command, tool call, stage transition, hook point, retry, or runtime behavior. No docs-site diff is needed: the CLI and user-visible runtime behavior remain unchanged, and `docs/runtime-support.md` already defines named capability binding.

## Lint design

Add `internal/contractlint/fo_function_reference_invariant_test.go` with:

- a fixed `foFunctionReferencePaths` registry containing the 13 files above;
- a `mutableProcedureAddress` classifier for `step N`, `step-N`, ranges such as `steps 1-3`, reuse-condition numbers, legacy tier-number pointers, entry-point principle numbers, and parenthesized numeric “above” pointers;
- a production test requiring **zero** matches and reporting `path:line` plus the matched address;
- a discriminator table that plants failures in shared, Claude, Codex, Pi, gate, and legacy paths, proving every surface uses the same classifier;
- pass controls for local ordered-list definitions (`1.`, `0.5.`), feedback cycle 3, exit 3, AC-2, versions, and issue numbers;
- a section-bounded call-site matrix implementing every row above, not merely a global reference count;
- body-ownership checks for the moved branch/list structures, so a copied old body fails even when the new function exists;
- `TestFOLocalOrderedProceduresPreserved`, which extracts the real sections and compares exact marker sequences: Startup `1,2,3`; Dispatch `1..9`; reuse `0..4`; Event Loop `0.5,1,2,3`; present-gate `1..11`; feedback rejection `1..7`; legacy recovery `1..3`; legacy teardown `1..4`. It also binds each item to an independent structural anchor: Startup (`Binary version gate`, `«state.boot»`, `«interaction.boundary»`); Dispatch (entity/stage read, `«dispatch.checklist»`, conflict check, `dispatch_agent_id`, status set, state commit, worktree, `«dispatch.build»`→`«worker.spawn»`, `«completion-signal»`); reuse (`«context-budget»`, `«addressable-worker»`, `fresh: true`, worktree routing, `«reuse.model-match»`); Event Loop (`«addressable-worker»`, `mod-block`, first `status --next`, `«hooks.run»`+`«roster-reconcile»`+second `status --next`); present-gate (`Lede first`, `Chosen direction`, `Stage Report`, `Reviewer findings`, `Recommendation`, `Bounce-back`, `format-pedantry`, `worktree`, `Target length`, `declared label`, `verification state`); feedback (`feedback-to`, `Feedback Cycles`, `cycle 3`, `«context-budget»`, `«addressable-worker»`, reviewer, gate flow); legacy recovery (`Fresh-suffixed TeamCreate`, `Degraded Mode`, `Surface to captain`); legacy teardown (`shutdown_request`, `TeamDelete`, settle, `TERMINAL_TEARDOWN_BOUNDED`). An empty or reordered numbered shell therefore fails.

Update structural tests that currently anchor on `step 0`: `reconcile_class_binding_test.go` should extract the named Claude `«roster-reconcile»` block, and the layering discriminator's idle-hook control should use `«hooks.run»("idle")` plus `«roster-reconcile»`. Narrow `runtime_binding_block_test.go` to its negative-host-contrast duty or delegate numeric-address coverage to the new invariant. Do not add phrase-by-phrase semantic assertions.

## Scope

- Change only the active instruction files needed by the inventory and the contractlint tests that enforce them.
- Preserve every local ordered procedure and every current host binding.
- Do not change Go runtime code, command routing, output, scenarios, fixtures, hook ordering, or lifecycle behavior.

## Acceptance criteria

- **AC-1 (end value: zero mutable addresses).** The fixed 13-file active surface contains **0 mutable numeric procedure addresses**, down from the measured baseline of **51**, while every real local ordered procedure retains its exact marker sequence. *Verified by:* `TestFOFunctionReferenceInvariant`, the reproducible per-file baseline vector, and `TestFOLocalOrderedProceduresPreserved`; pass controls separately prove local markers and semantic counters are legal.
- **AC-2 (named closure, no duplicated procedure).** Every boot, interaction, dispatch-checklist, event-loop, reuse, completion, teardown, recovery, and lifecycle-hook cross-reference resolves to the named owner in the inventory. `«interaction.boundary»`, `«dispatch.checklist»`, `«hooks.run»`, and `«legacy-team.recover»` each have one definition; existing capabilities retain one owner and explicit runtime binding. No moved procedure body remains duplicated at its old numbered site. *Verified by:* the section-bounded required call-site matrix and body-ownership checks above; the implementation diff audit is supplemental evidence, not the invariant.
- **AC-3 (non-vacuous lint boundary).** The same classifier rejects planted mutable addresses attributed to shared core, Claude, Codex, Pi, present-gate, and legacy-helper paths, and accepts local ordering, cycle/exit values, AC IDs, versions, and issue numbers. *Verified by:* table-driven discriminator tests that call the production classifier, not a copy.
- **AC-4 (behavior preserved).** The reference-only diff changes no binary route, emitted command, host tool signature, hook ordering, or durable workflow outcome. Existing offline, race, live-tag compile, and protected Claude/Codex/Pi live scenarios remain green. *Verified by:* the exact commands below, a changed-file audit showing instruction/contractlint paths only, and successful fail-closed `claude-live` (both matrix legs), `codex-live`, and `pi-live` jobs for the implementation SHA. A local live test that skips is not green evidence.

## Test plan

- **RED:** add the production scanner first. `go test ./internal/contractlint -run 'TestFOFunctionReferenceInvariant|TestFOFunctionReferenceClassifierDiscriminates' -count=1` must report the current 51 addresses, with representative failures in shared, Claude, Codex, Pi, gate, and legacy surfaces.
- **Focused GREEN:** run `go test ./internal/contractlint -run 'TestFOFunctionReference|TestFOLocalOrderedProceduresPreserved|TestReconcileClassBinding|TestCodexAndPiFirstOfficerRuntime|TestCapabilityBinding|TestStartupRecipe' -count=1`, then `go test ./internal/contractlint/... -count=1`.
- **Coupled offline:** run `go test -tags live ./internal/ensigncycle -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions' -count=1` (metadata only; no model calls) and `go test ./internal/cli -run TestProseFunction -count=1`.
- **Repository gates:** run scoped `gofmt` on touched Go tests, `git diff --check`, `go test ./...`, `go test ./... -race`, and `go test -tags live -run '^$' ./...`.
- **Protected live behavior:** use `.github/workflows/runtime-live-e2e.yml` without weakening assertions. The commands are `go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v`; `go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v`; `go test -tags live -count=1 -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v`; and `go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestPiSharedScenarioCoverage' ./internal/ensigncycle -v`. Require successful workflow job conclusions for the implementation SHA: both approval-gated `claude-live` matrix jobs (whose secret check fails when `ANTHROPIC_API_KEY` is missing), `codex-live` with `SPACEDOCK_CODEX_LIVE_REQUIRED=1`, and `pi-live` with `SPACEDOCK_PI_LIVE_REQUIRED=1`. Record the run URL, SHA, and job conclusions. Local SKIP output does not satisfy AC-4; the workflow's whole jobs, including their additional canaries, are required.

## Mechanism spike

The riskiest mechanism is the lint boundary, not runtime execution. The read-only candidate classifier was exercised against the exact fixed scope before design: it returned 51 current matches and did not classify local ordered-list markers. Its planned discriminator adds explicit fail and pass fixtures. No runtime spike is needed because the task changes addresses only; the existing protected live matrix already exercises boot, dispatch, gate, feedback, and teardown behavior.

## Stage Report: ideation

- DONE: Audit every active first-officer numeric procedure address and lifecycle-hook call across shared cores, Claude/Codex/Pi adapters, and directly linked gate, rejection, legacy, recovery, and write helpers. Map each address to an existing or proposed named `«fn»` without copying its procedure body.
  The inventory covers all 51 measured numeric addresses and every startup/idle/merge hook call. Existing functions absorb most references; four narrowly scoped functions are proposed for interaction routing, checklist assembly, hook execution, and legacy recovery.
- DONE: Define a non-tautological lint boundary that catches numeric cross-reference drift while preserving local ordered lists.
  The design uses one fixed 13-file registry and one production classifier with planted failures for shared, Claude, Codex, Pi, gate, and legacy paths. Pass controls cover local list markers and semantic numbers, preventing a blanket-number ban.
- DONE: Refine acceptance criteria around the end value and record exact implementation proof without expanding runtime behavior.
  AC-1 measures 51 → 0 mutable addresses; AC-2 pins named ownership and no duplicated body; AC-3 proves the classifier can fail and discriminate; AC-4 requires unchanged offline, race, compile-only, and protected live behavior. The scope excludes Go runtime and command changes. Independent staff review challenged ordered-procedure preservation, call-site ownership, hook timing, and fail-closed live evidence; the revised design resolved all material findings and received Ready: Yes.

### Summary

Normalize the active first-officer contract around named functions while retaining local execution order. A fixed-scope contractlint test drives the measured 51-to-0 cleanup, checks every new/existing owner binding, and rejects future numeric cross-file addresses without banning legitimate numbers. The implementation is instruction-and-lint only; the protected runtime matrix proves behavior remains unchanged.
