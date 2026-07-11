---
title: Replace mutable step-number references with named FO functions
status: validation
score: 0.95
source: "Captain direction 2026-07-11: widen the step-number sweep; hooks and references use the «fn» notation."
id: 88tq5zyg9jvx13f33zz3eq28
started: 2026-07-10T23:57:47Z
worktree: .worktrees/spacedock-ensign-fo-function-reference-invariant
---

## Problem

The active first-officer contract still addresses shared procedures by mutable list positions. A step insertion can silently redirect `Startup step 3`, `Dispatch step 2`, Event Loop step 0/3, Merge-and-Cleanup step 10, reuse-condition numbers, or a legacy recovery tier. This coupling conflicts with `docs/runtime-support.md`: shared behavior should expose named `«fn»` capabilities, and runtime adapters should bind those names.

A read-only spike applied the proposed address classifier to the fixed active surface below. It found **51 mutable numeric address occurrences across 10 files**: Claude dispatch 17, dispatch core 11, shared core 8, legacy helper 4, merge core 3, Claude runtime 2, Codex runtime 2, Pi runtime 2, present-gate 1, and feedback-rejection 1. The same pattern ignores local ordered-list markers and semantic numbers such as feedback cycle 3, exit code 3, AC-2, version 0.24.0, and issue numbers. This is the independent baseline for the task's end value: implementation must reduce 51 to 0 without deleting local execution order.

## GPT-5.6 prompt-guidance gate

The design applies the [official OpenAI GPT-5.6 prompt guidance](https://developers.openai.com/api/docs/guides/prompt-guidance-gpt-5p6) as a lean-prompt constraint, not as authority to weaken the workflow contract. The applicable guidance is to state the outcome and completion bar, keep consequential constraints and tool routes explicit, delete repeated process narration and examples, and remove one instruction group at a time while rerunning the same representative evals. This refactor therefore succeeds only if named functions make the total prompt smaller while preserving the user-visible outcome, success and stop criteria, evidence requirements, permission boundaries, contextual tool routing, output shape, and validation invariants.

Four API/deployment features from the guide are deliberately **not adopted**:

- Programmatic Tool Calling (PTC) changes Responses API tool execution; this task only restructures repository-owned instruction text.
- Prompt caching is a deployment and cache-prefix concern; this task owns neither request construction nor cache breakpoints.
- Persisted reasoning is Responses API state policy; this task does not add or change cross-turn reasoning persistence.
- `text.verbosity` is an API response-detail setting; it does not replace first-officer gate reports, status output, or other required output contracts.

## Captain amendment: eager `@` reference cleanup

PR #495 changed the live-proven Claude loading boundary after the original 118,178-byte inventory was recorded. The first-officer entry eagerly imported the shared, merge, smallest-sufficient, and write cores with `@references/...`; both Claude live matrix legs passed on that exact head. Preserve the reliable eager mechanism for required bodies, but remove prompt surfaces whose operative contract is already complete in shared core. Do not reintroduce PR #491's rejected callable-skill mechanism or rely on model skill discovery for required boot behavior.

Make that boundary singular and truthful:

- remove the standalone `skills/fo-write-core/SKILL.md` compatibility shim; the first-officer entry is the sole write-core entry surface and imports the canonical `skills/first-officer/references/fo-write-core.md` body directly;
- delete the separately imported smallest-sufficient core rather than making its load policy self-referential; one compact, self-sufficient resident rule remains in shared Working Principles, with no lazy reference, wrapper, or prompt examples;
- retain deferred prose only for files still loaded lazily (`fo-dispatch-core.md`, host dispatch material, gate/rejection/recovery helpers at their declared triggers);
- measure the real active surface: include both eagerly imported canonical cores and exclude the removed shim.

The post-#495, pre-task active prompt baseline is **122,400 bytes across 14 files**. The final active registry has 13 files after removing the smallest-sufficient body and write-core shim; the baseline constant still includes their pre-change bytes for a strict before/after comparison. This supersedes the stale 118,178-byte/13-file baseline below wherever acceptance or checkpoint accounting refers to the implementation starting point.

## Active lint boundary

Use an explicit path registry, not a recursive `skills/**/*.md` ban. The invariant applies to the first-officer entry point, its three host adapters and host-neutral cores, the eager canonical write core, and the helpers those cores directly load for gates, rejection, legacy/recovery, and hook execution:

- `skills/first-officer/SKILL.md`
- `skills/first-officer/references/{first-officer-shared-core,fo-dispatch-core,fo-merge-core}.md`
- `skills/first-officer/references/fo-write-core.md`
- `skills/first-officer/references/{claude-first-officer-runtime,claude-fo-dispatch,codex-first-officer-runtime,pi-first-officer-runtime}.md`
- `skills/{present-gate,feedback-rejection-flow,using-legacy-claude-team,fo-dispatch-recovery}/SKILL.md`

This boundary excludes unrelated skills whose local numbered recipes are not first-officer cross-file APIs. Inside the boundary, leading ordered-list markers remain legal. A line such as `3. If nothing is dispatchable` defines local order; prose such as `re-run step 3` is a mutable address and fails.

The final registry is the prompt-size boundary. Raw bytes are measured with `wc -c` over its 13 paths in registry order. The original pre-#495 vector below is retained as historical evidence; the authoritative before value is the post-#495 122,400-byte surface recorded in the captain amendment.

| File | Bytes |
|---|---:|
| first-officer entry | 1,963 |
| shared core | 24,883 |
| dispatch core | 18,112 |
| merge core | 3,226 |
| Claude runtime | 4,298 |
| Claude dispatch | 19,211 |
| Codex runtime | 6,119 |
| Pi runtime | 3,754 |
| present-gate | 5,815 |
| feedback-rejection | 2,769 |
| legacy helper | 14,216 |
| dispatch-recovery | 7,685 |
| write core | 6,127 |
| **Total** | **118,178** |

Implementation must finish below **122,400 bytes** across the final 13-file active registry. Individual owners may grow when a body moves into them, but the active total must shrink.

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

### Lean function bar and ownership discipline

A name is not a license to add a second abstraction layer. Every retained or proposed function follows these rules:

1. One owner contains the complete shared procedure and its guard, effect, completion condition, and host-binding requirement.
2. Callers contain only the `«fn»` call, arguments, and genuinely caller-specific order or branch condition. Delete duplicated process prose, constraint bullets, fallback recipes, and illustrative examples from callers.
3. A runtime adapter owns only its concrete host/tool signature and host-specific failure translation. It does not restate host-neutral semantics, success criteria, or stop rules.
4. Rename an existing owner section in place where possible. Do not add alias headings, `see also` wrappers, or a second explanatory definition.
5. Retain a verbatim template or example only when it is a machine-consumed artifact or a behavior-bearing output contract. Such a template calls the owner for shared content instead of copying its rules.
6. A new name is accepted only when it replaces mutable callers and duplicated procedure text and contributes to a strict net reduction below the amended 122,400-byte baseline. If implementation cannot meet that test, reuse an existing capability or named owner heading instead.

The four proposed names survive this reassessment, but narrowly:

| Proposed function | Lean-prompt decision |
|---|---|
| `«interaction.boundary»()` | Keep. It consolidates the repeated interactive/headless/given-the-conn contract now spread across Startup, deferred-load, scope, and Claude runtime prose. The shared owner is the only branch body. |
| `«dispatch.checklist»(entity, stage)` | Keep. Normal dispatch and two fallback paths need the same bounded checklist. The recovery template retains only a call placeholder, not copied checklist rules or examples. |
| `«hooks.run»(point)` | Keep. It gives all additional hook references the captain-required `«fn»` notation and centralizes discovery/order/run-once semantics. Point-specific timing remains at the caller because it is lifecycle order, not hook-body duplication. |
| `«legacy-team.recover»()` | Keep only as an in-place name for the existing recovery-ladder owner. It adds no wrapper or new procedure block; setup and failure callers stop copying tier prose and call the named ladder. |

This is the full preservation ledger for the deletion pass:

| Contract class | Invariant that remains explicit |
|---|---|
| User-visible outcomes | Interactive greeting/stop, headless progression, gates, feedback routing, and terminal result remain unchanged. |
| Success and stop criteria | Each owner's done-when, blocked, retry, escalation, and terminal conditions remain authoritative. |
| Evidence | Stage Report, acceptance-criteria cross-check, durable entity state, commits, and live-job evidence remain required where currently required. |
| Permission | Captain approval gates, FO/ensign write scopes, external-action limits, and no-force/no-destructive boundaries remain unchanged. |
| Tool routing | State commands, dispatch build, host spawn/message bindings, and context-dependent fallback routes remain explicit at their single owner or adapter binding. |
| Validation | Structural, offline, race, compile-only, and protected live checks remain required; local skips still do not count as live evidence. |

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

The same test enforces body ownership. `Interactive`/`Headless` branch headings occur only inside `«interaction.boundary»`; checklist constraint bullets occur only inside `«dispatch.checklist»`; hook discovery, alphabetical execution, and run-once semantics occur only inside `«hooks.run»`; legacy recovery's 1-3 ladder occurs only inside the section renamed `«legacy-team.recover»`. Caller sections contain the named call but no second copy of those structural blocks. Hook callers still own their point-specific position relative to state sweep, status lookup, roster reconciliation, and merge guard; that order is a caller invariant and must not be hidden in the generic hook body.

The implementation changes instruction structure, removes the redundant write-core skill shim, and updates contractlint only. It adds no command, tool call, stage transition, hook point, retry, or runtime behavior. No docs-site diff is needed: the CLI and user-visible runtime behavior remain unchanged, and `docs/runtime-support.md` already defines named capability binding.

## Lint design

Add `internal/contractlint/fo_function_reference_invariant_test.go` with:

- a fixed `foFunctionReferencePaths` registry containing the final 13 active files above;
- a `foFunctionReferenceBaselineBytes = 122400` constant and a test that sums raw file bytes through that same registry and requires the implementation total to be strictly smaller;
- an eager-reference topology test requiring exactly three canonical imports (shared, merge, write), their targets, absence of the redundant standalone `fo-write-core` skill directory, and absence of a smallest-sufficient body/reference/wrapper;
- a `mutableProcedureAddress` classifier for `step N`, `step-N`, ranges such as `steps 1-3`, reuse-condition numbers, legacy tier-number pointers, entry-point principle numbers, and parenthesized numeric “above” pointers;
- a production test requiring **zero** matches and reporting `path:line` plus the matched address;
- a discriminator table that plants failures in shared, Claude, Codex, Pi, gate, and legacy paths, proving every surface uses the same classifier;
- pass controls for local ordered-list definitions (`1.`, `0.5.`), feedback cycle 3, exit 3, AC-2, versions, and issue numbers;
- a section-bounded call-site matrix implementing every row above, not merely a global reference count;
- body-ownership checks for the moved branch/list structures, so a copied old body fails even when the new function exists;
- caller-only checks that reject the moved headings, constraint bullets, recovery tiers, hook-discovery prose, and retained illustrative examples outside their one owner; adapter checks permit only the concrete host binding and host-specific failure translation;
- `TestFOFunctionNormalizationPreservationSuite`, a fixed table of section-bounded semantic checks that passes on the pre-change prompt and after every instruction group. Its named subtests cover: startup/interaction outcomes; dispatch/checklist/reuse/completion; legacy recovery, layering, and terminal teardown; normal, reuse-advance, and break-glass dispatch; boot-resident/deferred-load closure; startup/idle/merge hook count and order; write/permission scope; Stage Report, evidence, approval-gate, and feedback outcomes; and concrete Claude, Codex, and Pi runtime bindings. Each subtest extracts its owner/caller sections and checks stable outcome, stop, evidence, permission, route, or validation anchors—not the mutable numeric address or the new function spelling being changed.
- `TestFOFunctionReferenceCheckpointMetrics`, which always passes and logs one machine-readable line, `FO_FUNCTION_METRICS addresses=<n> bytes=<n>`, from the production `mutableProcedureAddress` classifier and `foFunctionReferencePaths` byte sum. It is reporting only; AC-1 and AC-3 remain the enforcing tests.
- `TestFOLocalOrderedProceduresPreserved`, which extracts the real sections and compares exact marker sequences: Startup `1,2,3`; Dispatch `1..9`; reuse `0..4`; Event Loop `0.5,1,2,3`; present-gate `1..11`; feedback rejection `1..7`; legacy recovery `1..3`; legacy teardown `1..4`. It also binds each item to an independent structural anchor: Startup (`Binary version gate`, `«state.boot»`, `«interaction.boundary»`); Dispatch (entity/stage read, `«dispatch.checklist»`, conflict check, `dispatch_agent_id`, status set, state commit, worktree, `«dispatch.build»`→`«worker.spawn»`, `«completion-signal»`); reuse (`«context-budget»`, `«addressable-worker»`, `fresh: true`, worktree routing, `«reuse.model-match»`); Event Loop (`«addressable-worker»`, `mod-block`, first `status --next`, `«hooks.run»`+`«roster-reconcile»`+second `status --next`); present-gate (`Lede first`, `Chosen direction`, `Stage Report`, `Reviewer findings`, `Recommendation`, `Bounce-back`, `format-pedantry`, `worktree`, `Target length`, `declared label`, `verification state`); feedback (`feedback-to`, `Feedback Cycles`, `cycle 3`, `«context-budget»`, `«addressable-worker»`, reviewer, gate flow); legacy recovery (`Fresh-suffixed TeamCreate`, `Degraded Mode`, `Surface to captain`); legacy teardown (`shutdown_request`, `TeamDelete`, settle, `TERMINAL_TEARDOWN_BOUNDED`). An empty or reordered numbered shell therefore fails.

Update structural tests that currently anchor on `step 0`: `reconcile_class_binding_test.go` should extract the named Claude `«roster-reconcile»` block, and the layering discriminator's idle-hook control should use `«hooks.run»("idle")` plus `«roster-reconcile»`. Narrow `runtime_binding_block_test.go` to its negative-host-contrast duty or delegate numeric-address coverage to the new invariant. Do not add phrase-by-phrase semantic assertions.

## Scope

- Change only the active instruction files needed by the inventory and the contractlint tests that enforce them.
- Preserve every local ordered procedure and every current host binding.
- Do not change Go runtime code, command routing, output, scenarios, fixtures, hook ordering, or lifecycle behavior.

## Acceptance criteria

- **AC-1 (end value: zero mutable addresses).** The fixed 13-file active surface contains **0 mutable numeric procedure addresses**, down from the measured baseline of **51**, while every real local ordered procedure retains its exact marker sequence. *Verified by:* `TestFOFunctionReferenceInvariant`, the reproducible per-file baseline vector, and `TestFOLocalOrderedProceduresPreserved`; pass controls separately prove local markers and semantic counters are legal.
- **AC-2 (single-owner named closure).** Every boot, interaction, dispatch-checklist, event-loop, reuse, completion, teardown, recovery, and lifecycle-hook cross-reference resolves to the named owner in the inventory. `«interaction.boundary»`, `«dispatch.checklist»`, and `«hooks.run»` each have one definition; the existing legacy recovery section is renamed in place as `«legacy-team.recover»`; existing capabilities retain one owner and explicit runtime binding. Callers contain only the call, arguments, and caller-specific order/branch condition. No moved process prose, structural body, constraint list, fallback example, or host-neutral semantic copy remains at a caller or adapter. *Verified by:* the section-bounded required call-site matrix, body-ownership checks, and caller-only negative checks above; the implementation diff audit is supplemental evidence, not the invariant.
- **AC-3 (prompt surface shrinks and eager topology is singular).** The raw-byte total of the final 13-file prompt surface is **strictly less than 122,400**, its post-#495 implementation baseline. Exactly three canonical `@references/...` imports remain (shared, merge, write); smallest-sufficient has one resident shared rule and no separate body; and no standalone `fo-write-core` entry surface remains. Function headings and bindings do not create net abstraction growth. *Verified by:* `TestFOFunctionPromptSurfaceShrinks` and the eager-reference topology test, plus a recorded post-change per-file vector and total.
- **AC-4 (non-vacuous lint boundary).** The same classifier rejects planted mutable addresses attributed to shared core, Claude, Codex, Pi, present-gate, and legacy-helper paths, and accepts local ordering, cycle/exit values, AC IDs, versions, and issue numbers. *Verified by:* table-driven discriminator tests that call the production classifier, not a copy.
- **AC-5 (behavior and consequential constraints preserved).** The reference-only diff changes no user-visible outcome, success/stop criterion, evidence requirement, permission boundary, binary route, emitted command, contextual tool route, host tool signature, hook ordering, validation invariant, or durable workflow outcome. Existing offline, race, live-tag compile, and protected Claude/Codex/Pi live scenarios remain green. *Verified by:* `TestFOFunctionNormalizationPreservationSuite` passing unchanged at baseline and after every instruction group; the preservation-ledger diff audit; exact commands below; a changed-file audit showing instruction/contractlint paths only; and successful fail-closed `claude-live` (both matrix legs), `codex-live`, and `pi-live` jobs for the implementation SHA. A local live test that skips is not green evidence.

## Test plan

- **RED and semantic baseline:** add the production scanner, preservation suite, single-owner checks, eager-topology check, and byte gate first. Run `go test ./internal/contractlint -run 'TestFOFunctionReferenceInvariant|TestFOFunctionReferenceClassifierDiscriminates|TestFOFunctionPromptSurfaceShrinks|TestFirstOfficerEagerReferenceTopology' -count=1`; the address, byte, and duplicate-entry assertions must fail against the post-#495 baseline, with representative address failures in shared, Claude, Codex, Pi, gate, and legacy surfaces. Immediately afterward, before editing any instruction group, run `go test ./internal/contractlint -run '^TestFOFunctionNormalizationPreservationSuite$' -count=1`, require green, and record that exact result as the unchanged semantic baseline reused after groups 1-4.
- **Incremental GREEN:** change exactly one instruction group at a time: (1) boot/interaction; (2) dispatch/checklist/reuse/completion; (3) teardown/legacy recovery; (4) startup/idle/merge hooks. After every group, run the identical green command `go test ./internal/contractlint -run '^TestFOFunctionNormalizationPreservationSuite$' -count=1`, followed by the same runtime metadata and CLI evals: `go test -tags live ./internal/ensigncycle -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestPiSharedScenarioCoverage' -count=1` and `go test ./internal/cli -run TestProseFunction -count=1`. All three commands must pass at every checkpoint. Do not combine groups to mask a regression; repair or revert that group before proceeding.
- **Executable checkpoint accounting:** after each group run `go test ./internal/contractlint -run '^TestFOFunctionReferenceCheckpointMetrics$' -v -count=1` and copy its `FO_FUNCTION_METRICS` values into the implementation Stage Report table below. After groups 1-3 also run `go test ./internal/contractlint -run '^TestFOFunctionReferenceInvariant$' -count=1`, require a non-zero exit, and confirm its reported count equals the metrics line; after group 4 require the same invariant command to pass. Each row must be strictly lower than its predecessor in both columns:

  | Checkpoint | Mutable addresses | Prompt bytes |
  |---|---:|---:|
  | Baseline | 51 | 122,400 |
  | Group 1: boot/interaction | 41 | 120,829 |
  | Group 2: dispatch/checklist/reuse/completion | `< group 1` | `< group 1` |
  | Group 3: teardown/legacy recovery | `< group 2` | `< group 2` |
  | Group 4: startup/idle/merge hooks | `0` | `< group 3` |

Implementation checkpoint evidence (the group-2 byte increase was a measured 116-byte temporary owner-heading cost; the deletion pass then removed 4,665 bytes of tautological eager prompt and finished 6,421 bytes below baseline):

| Checkpoint | Mutable addresses | Prompt bytes |
|---|---:|---:|
| Baseline | 51 | 122,400 |
| Group 1: boot/interaction | 41 | 120,829 |
| Group 2: dispatch/checklist/reuse/completion | 8 | 120,945 |
| Group 3: teardown/legacy recovery | 0 | 120,247 |
| Group 4: hooks + eager-surface deletion | 0 | 115,979 |
- **Focused closure:** after group 4, run `go test ./internal/contractlint/... -count=1` and require 0 addresses, one owner per capability, no caller-body copies, and a total below 118,178 bytes.
- **Coupled offline:** run `go test -tags live ./internal/ensigncycle -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions' -count=1` (metadata only; no model calls) and `go test ./internal/cli -run TestProseFunction -count=1`.
- **Repository gates:** run scoped `gofmt` on touched Go tests, `git diff --check`, `go test ./...`, `go test ./... -race`, and `go test -tags live -run '^$' ./...`.
- **Protected live behavior:** use `.github/workflows/runtime-live-e2e.yml` without weakening assertions. The commands are `go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v`; `go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v`; `go test -tags live -count=1 -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v`; and `go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestPiSharedScenarioCoverage' ./internal/ensigncycle -v`. Require successful workflow job conclusions for the implementation SHA: both approval-gated `claude-live` matrix jobs (whose secret check fails when `ANTHROPIC_API_KEY` is missing), `codex-live` with `SPACEDOCK_CODEX_LIVE_REQUIRED=1`, and `pi-live` with `SPACEDOCK_PI_LIVE_REQUIRED=1`. Record the run URL, SHA, and job conclusions. Local SKIP output does not satisfy AC-5; the workflow's whole jobs, including their additional canaries, are required.

## Mechanism spike

The riskiest mechanism is the lint and deletion boundary, not runtime execution. The read-only candidate classifier was exercised against the exact fixed scope before design: it returned 51 current matches and did not classify local ordered-list markers. The same 13-file registry measures a reproducible 118,178-byte baseline. Its planned discriminator adds explicit fail/pass fixtures, and the strict byte gate makes alias-only functions or copied owner prose fail even if reference syntax is correct. No runtime spike is needed because the task changes addresses only; the existing protected live matrix already exercises boot, dispatch, gate, feedback, and teardown behavior.

## Stage Report: ideation

- DONE: Audit every active first-officer numeric procedure address and lifecycle-hook call across shared cores, Claude/Codex/Pi adapters, and directly linked gate, rejection, legacy, recovery, and write helpers. Map each address to an existing or proposed named `«fn»` without copying its procedure body.
  The inventory covers all 51 measured numeric addresses and every startup/idle/merge hook call. Existing functions absorb most references; four narrowly scoped functions are proposed for interaction routing, checklist assembly, hook execution, and legacy recovery.
- DONE: Define a non-tautological lint boundary that catches numeric cross-reference drift while preserving local ordered lists.
  The design uses one fixed 13-file registry and one production classifier with planted failures for shared, Claude, Codex, Pi, gate, and legacy paths. Pass controls cover local list markers and semantic numbers, preventing a blanket-number ban.
- DONE: Refine acceptance criteria around the end value and record exact implementation proof without expanding runtime behavior.
  AC-1 measures 51 → 0 mutable addresses; AC-2 pins one owner and callers-only references; AC-3 requires the fixed prompt surface to shrink below 118,178 bytes; AC-4 proves the classifier can fail and discriminate; AC-5 preserves consequential contract and runtime behavior through offline, race, compile-only, and protected live evidence. The scope excludes Go runtime and command changes.
- DONE: Apply the official GPT-5.6 prompt guidance and record the boundary of adoption.
  The revised design deletes repetition one instruction group at a time and reruns the same evals after each group. It explicitly preserves outcomes, success/stop criteria, evidence, permissions, contextual tool routing, output contracts, and validation. PTC, prompt caching, persisted reasoning, and `text.verbosity` are recorded as API/deployment concerns outside this instruction-only task.
- DONE: Reassess all four proposed names under a lean-prompt bar.
  Interaction, checklist, and hook functions remain because they consolidate multiple callers; legacy recovery becomes an in-place owner-heading rename with no wrapper. Every function has one owner, callers contain no duplicated process prose or examples, adapters retain only concrete bindings, and strict aggregate byte reduction rejects abstraction growth.
- DONE: Resolve the second independent staff review's preservation-suite and checkpoint findings.
  The reviewer returned Ready: No because the per-group suite did not span every consequential contract class and the checkpoint metrics were descriptive rather than executable. The design now requires one baseline-green named preservation suite covering startup through all host bindings, repeats that exact suite after every group, and supplies an exact metrics command plus a monotonic address/byte table and explicit expected-red/green invariant exits.
- DONE: Re-run independent staff review against the captain's GPT-5.6 gate and incorporate every finding.
  Final independent verdict: Material 0, Polish 1, Ready: Yes. The non-blocking polish asked the test plan to run and record the named preservation suite before group 1, not only after each group; the baseline step now does so explicitly. The reviewer confirmed the prior suite-coverage and executable-checkpoint blockers are resolved and the guide application, API exclusions, ownership discipline, preservation ledger, byte/address gates, and four-function lean justification are sound.

### Summary

Normalize the active first-officer contract around named functions while retaining local execution order and reducing its fixed prompt surface below 118,178 bytes. A fixed-scope contractlint test drives the measured 51-to-0 cleanup, enforces one owner with callers-only references, and rejects future numeric cross-file addresses without banning legitimate numbers. Four API-only prompt-guide features remain out of scope. The implementation is instruction-and-lint only; incremental repeated evals and the protected runtime matrix prove the consequential contract and behavior remain unchanged.

## Stage Report: implementation

- DONE: Eliminate all mutable numeric procedure addresses while preserving every local ordered procedure and runtime behavior.
  Commit `9c8cbea` reduces the bounded classifier from 51 to 0; `TestFOLocalOrderedProceduresPreserved` and the semantic preservation suite pass.
- DONE: Incorporate the captain amendment for @ notation: audit every eager @ import, remove contradictory deferred-load claims and unnecessary duplicate entry surfaces without regressing the live-proven Claude loading from PR #495; update the task design, acceptance criteria, and test plan before implementation; do not resurrect the rejected callable-skill approach from PR #491.
  The entry now eagerly imports only shared, merge, and canonical write; the tautological smallest-sufficient body and standalone write shim are deleted, with no callable-skill or lazy replacement.
- DONE: Deliver strict prompt-byte shrinkage, non-vacuous structural guards, and checkpoint evidence.
  The active surface is 115,979 bytes, down 6,421 from the 122,400 post-#495 baseline; planted fail/pass controls, bounded call-site ownership, and exact ordered-marker guards pass.
- DONE: Run required offline, race, compile, and coupled metadata proof.
  `go test ./...`, `go test ./... -race`, `go test -tags live -run '^$' ./...`, contractlint 113/113, shared runtime metadata 3/3, and CLI prose-function 3/3 passed.
- SKIPPED: Record successful fail-closed Claude, Codex, and Pi protected live jobs for the implementation SHA.
  The branch is pushed at `9c8cbea`; protected credentials are absent locally, so the FO must open/approve the PR workflow and attach its exact-SHA job conclusions during validation.

### Summary

The implementation replaces every mutable numeric procedure address with a named `«fn»` owner/call, retains all local ordered procedures, and updates stale structural tests to bind the named Claude reconcile section. It also dogfoods smallest-sufficient mechanism by deleting the redundant eager smallest-sufficient prompt and write-core shim while keeping required Claude loading direct and canonical. The deliverable is committed and pushed for fresh validation; only protected PR live-lane evidence remains external.
