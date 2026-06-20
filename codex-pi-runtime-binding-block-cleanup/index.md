---
title: Convert Codex and Pi runtime adapters to capability binding blocks
status: validation
source: "Captain direction (2026-06-20): move toward per-host runtime files as bindings blocks keyed by core «fn» capability names, starting with Codex and Pi; recommend sequencing with ad/trim-dispatch-adapter-prose."
started: 2026-06-20T20:14:02Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-codex-pi-runtime-binding-block-cleanup
issue:
sprint: 0221-layered-fo
sprint-readiness:
id: t0gk2fatt18tj28xm6sr1xd1
mod-block: merge:pr-merge
pr: "#418"
---

Codex and Pi first-officer runtime references should become compact runtime implementation maps keyed by the shared core's `«fn»` capability names, rather than prose sections that re-narrate the lifecycle. The finished shape is a short runtime intro, a `## Runtime implementation` binding list, and only short residual sections for probe, wait, feedback, or harness notes that do not fit a capability bullet.

## Problem

`skills/first-officer/references/pi-first-officer-runtime.md` and `skills/first-officer/references/codex-first-officer-runtime.md` already contain most of the required binding facts, but they are spread across lifecycle prose sections (`Dispatch`, `Awaiting Completion`, `Follow-up and Reuse`, `Shutdown`, etc.). Pi still carries negative host contrast and step-number coupling. Codex is probe-based after PR #414 but still mixes binding facts with lifecycle narration. This duplicates responsibilities the shared core should own: when each capability is invoked.

There is also a core-shape prerequisite: `fo-dispatch-core.md` currently defines `«addressable-worker»`, `«async-dispatch»`, `«worker-identity»`, `«completion-signal»`, `«context-budget»`, and `«roster-reconcile»`, but it references `«worker.spawn»` and `«worker.shutdown»` without first-class capability headings. Codex and Pi cannot honestly key complete adapter binding blocks to shared `«fn»` names until those two capabilities are defined in the core.

## Proposed approach

Implement this as a focused Codex+Pi cleanup:

1. Add minimal first-class core headings in `skills/first-officer/references/fo-dispatch-core.md` for `## «worker.spawn»:` and `## «worker.shutdown»:`.
   - `«worker.spawn»` should bind the initial worker creation call that consumes `spacedock dispatch build` output. Keep this small: the existing mandatory build flow remains the detailed procedure.
   - `«worker.shutdown»` should bind cooperative terminal/supersede shutdown. Keep the core generic and put each host's concrete action in its `→` line.
   - Add only the minimum Claude segments needed for host coverage and existing `capability_binding` invariants. Do not trim or reorganize Claude adapter prose here; `ad` owns that.
2. Rewrite the Codex FO adapter to:
   - keep a short intro saying the shared core owns invocation timing
   - keep `## Live Tool Surface Probe` as a compact residual section
   - add `## Runtime implementation` with bullets for `«worker.spawn»`, `«addressable-worker»`, `«async-dispatch»`, `«worker-identity»`, `«completion-signal»`, `«worker.shutdown»`, `«context-budget»`, and `«roster-reconcile»`
   - move foreground-wait interruption/reinstallation and queued/autonomous mailbox classification into a short residual `## Codex wait notes`
   - move feedback reviewer reuse into the `«addressable-worker»` bullet or a short residual `## Feedback reviewer reuse` note if the bullet becomes too dense
   - remove `## Dispatch`, `## Reuse And Feedback Routing`, and `## Awaiting Completion` headings
3. Rewrite the Pi FO adapter to:
   - keep a short intro saying the shared core owns invocation timing and Pi owns host realization
   - add `## Runtime implementation` with the same capability bullets
   - fold model stamping and canonical Pi model space into `«worker-identity»`
   - fold file-verification gate semantics into `«completion-signal»`
   - fold `context: "fresh"`, `cwd: <resolved repo root>`, no `acceptance`, and harness-selected substrate into `«worker.spawn»` / `«async-dispatch»`
   - fold first-slice fresh redispatch/default no non-fresh resume into `«addressable-worker»`
   - fold `pi-subagents` complete/closed memory and `pi-agent-teams` member/team teardown into `«worker.shutdown»`
   - keep `## Live harness isolation` as the only expected residual section
   - remove `## Runtime Shape`, `## Dispatch`, `### Model Resolution`, `### Canonical Model Space`, `## Awaiting Completion`, `## Follow-up and Reuse`, and `## Shutdown` headings
4. Keep the boundary with `ad`: this task changes Codex+Pi adapters, the two minimal core capability definitions, and focused contractlint/test helpers only. `ad` remains responsible for Claude adapter trim, broad adapter prose reduction, and any final cross-runtime polish after this binding-block shape exists.

No spike needed: the implementation relies on existing Markdown parsing patterns in `internal/contractlint`, existing `capability_binding_test.go` extraction of `## «fn»:` headings plus per-host `→` lines, and existing Codex/Pi runtime facts already present in the adapter files. The riskiest piece is structural enforcement, and that should be handled first with failing contractlint cases before rewriting the adapters.

## Semantics to preserve

Codex must preserve:

- live probe-based binding, not runtime-version naming: `spawn_agent`, `send_message` + `followup_task` or compatibility `send_input`, `wait_agent`, `list_agents`, and unbound `interrupt_agent` unless a shutdown-specific probe exists
- helper-built dispatch prompt forwarding, sanitized task-name mapping, and retained mapping from helper-emitted identity to slug/stage/cycle
- final-status mailbox notification as the completion signal; foreground wait is an idle action only, not the signal itself
- foreground wait timeout/interruption is non-terminal, and the global wait is reinstalled only when waiting is again the next useful idle action
- queued/activity-driven notification classification versus autonomous idle wake-up classification
- validation reviewer reuse after a REJECTED report through the same turn-starting `«addressable-worker»` binding when the reviewer remains addressable

Pi must preserve:

- Pi-native dispatch substrate mapping, with `pi-subagents` default and `pi-agent-teams` optional adapter mapping
- `subagent(...)` dispatch with `context: "fresh"`, `cwd: <resolved repo root>`, no `acceptance`, and only additive transport metadata around the build artifact
- null-model stamping from the parent's live Pi model via `intercom({action:"list"})`; stage-declared model override preserved
- Pi canonical model space as Pi-native provider/model strings, not the Claude enum
- file-verification gate: subagent return or advisory alone never advances state
- same-turn continuation after a verified non-gated, non-terminal report
- fresh assignment cycles as the default first Pi slice; non-fresh resume only as an explicit manual/debug exception with durable metadata
- live harness isolation: isolated Pi config/session dirs, copied auth only, durable proof via exit code, entity state changes, git log, and stage report content

## Sequencing recommendation with `ad`

Recommended order: run this task before `ad`'s broad adapter trim.

Reason: this task is a narrow structural normalization for Codex+Pi only. It creates the binding-block shape that `ad` can then use as the target when trimming remaining per-adapter operational prose. Running `ad` first would force it to decide the same binding-placement question while also preserving Claude await/reuse/guardrails, increasing risk and review surface.

Decision: this task should add the minimal first-class `«worker.spawn»` and `«worker.shutdown»` headings to `fo-dispatch-core.md` before rewriting Codex/Pi adapters. Deferring that prerequisite to `ad` would leave Codex/Pi binding blocks keyed partly to names the shared core does not define, and would make `ad` solve a cross-runtime core-shape question while doing broad Claude prose trimming. This task should not otherwise restructure Claude.

## Out of scope

Do not change runtime behavior, live runners, launch/install UX, or Claude guardrails. Do not remove Codex wait-interruption semantics or feedback-reviewer reuse semantics; express them as binding/probe notes if they remain load-bearing. Do not touch ensign runtime adapters unless an existing contractlint test requires wording alignment after the FO adapter headings move.

## Acceptance criteria

**AC-1 - Codex and Pi FO runtime adapters use a binding-block structure.**
Verified by: contractlint showing each file has exactly one `## Runtime implementation` section with bullets keyed to `«worker.spawn»`, `«addressable-worker»`, `«async-dispatch»`, `«worker-identity»`, `«completion-signal»`, `«worker.shutdown»`, `«context-budget»`, and `«roster-reconcile»`, and no longer uses lifecycle narration headings for dispatch/await/reuse/shutdown.

**AC-2 - Runtime support principles are enforced for Codex and Pi.**
Verified by: contractlint rejecting negative host contrast, mutable shared-procedure step-number coupling, and concrete host tool names outside `## Runtime implementation` / probe / harness sections in the Codex/Pi FO runtime adapters.

**AC-3 - Load-bearing host-specific semantics are preserved.**
Verified by: focused contractlint tests or table-driven structural checks showing Codex probe-based tool binding, wait reinstallation/interruption semantics, feedback reviewer reuse, Pi model stamping, Pi file-verification gate, and Pi harness isolation notes remain present in compact form.

**AC-4 - The task composes cleanly with `ad`.**
Verified by: implementation report naming what remains for `ad`, with no Claude adapter rewrite beyond unavoidable shared/core/test edits.

**AC-5 - Adapter prose is measurably smaller while preserving binding coverage.**
Verified by: a line-count or byte-count check in the implementation report showing the combined Codex+Pi FO adapter body size is lower than the pre-change `main` baseline, while AC-1 and AC-3 pass. The expected target is at least a 20% combined line-count reduction because lifecycle narration headings are removed rather than renamed.

## Test plan

Add or update focused tests under `internal/contractlint`:

- Extend capability binding coverage to include the new `«worker.spawn»` and `«worker.shutdown»` core definitions; the existing definition/call/host-coverage check should go red until the headings exist.
- Add a Codex/Pi runtime binding-block test that parses each adapter, extracts `## Runtime implementation`, and compares its bullet capability set to the expected ordered set.
- Add a lifecycle-heading absence test for Codex/Pi adapters rejecting `## Dispatch`, `## Awaiting Completion`, `## Follow-up and Reuse`, `## Reuse And Feedback Routing`, `## Shutdown`, `### Model Resolution`, and `### Canonical Model Space`.
- Replace or update `TestCodexToolNamesStayInRuntimeBindingSection` so Codex concrete tool names are allowed only in `## Runtime implementation` and `## Live Tool Surface Probe`; keep the existing negative-control intent around `interrupt_agent`.
- Add the same style of Pi guard for `subagent(`, `intercom(`, `member_spawn`, `delegate`, `message_dm`, `member_shutdown`, `team_done`, `context: "fresh"`, `cwd: <resolved repo root>`, and `acceptance`, allowing them only in runtime binding or harness sections.
- Add tests that reject mutable step-number coupling in Codex/Pi adapters (`Dispatch step`, `Event Loop step`, `Merge-and-Cleanup step`, bare `step 10`) and reject known negative host-contrast phrases unless a deliberately named compatibility note section is added.
- Keep `codex_foreground_wait_shape_test.go` or its successor pointed at the new Codex wait residual section so the non-terminal wait interruption wording remains guarded.

Run:

```bash
go test ./internal/contractlint
go test ./...
```

Use `go test ./... -race` only if implementation changes beyond Markdown and contractlint tests touch runtime code.

## Implementation notes for the next worker

- Start with tests. The first red should be the missing core headings and the missing `## Runtime implementation` sections.
- Use the existing `extractMarkdownSection` / `markdownSubsection` helpers if they remain suitable; otherwise add a small Markdown section extractor local to contractlint.
- When rewriting adapters, prefer one capability bullet per host fact. Use residual sections only when a fact spans multiple capabilities or documents a probe/harness concern.
- Keep docs/runtime-support.md unchanged unless the implementation discovers the preferred shape is incomplete; it already states the target contract.

## Peer feedback from Pi commander

PR #417 (`pi-fo-runtime-runtime-support-compliance`) is a subset of this Pi rewrite, not a conflict. This task supersedes #417's Pi FO adapter shape by removing the lifecycle sections #417 touched and relocating their positive bindings into `## Runtime implementation` bullets. Preserve #417's substance and update its guard tests rather than deleting their intent: `TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast` should still ban the smell phrases, but its positive-binding assertions should target the new `→` binding bullets (model-space provider/model strings, `«worker.shutdown»` realization), not old prose sentences. Keep `TestPiEnsignRuntimeAvoidsNegativeHostContrast` intact unless an ensign-side change is explicitly needed.

The runtime-support guide update is now an explicit prerequisite/scope input: merged `origin/main` commit `4403e095` adds `### Runtime binding-block shape` to `docs/runtime-support.md`. Implementation should leave it alone unless the binding-block work proves the guide incomplete.

The core-heading prerequisite remains correct. Add `«worker.spawn»` and `«worker.shutdown»` first with failing contractlint cases before adapter rewrites. Once `«worker.shutdown»` is first-class, the Pi `→` line must carry the #417 shutdown substance: for `pi-subagents`, a completed child invocation needs no mailbox shutdown and the FO marks the worker complete/closed in memory; for `pi-agent-teams`, the adapter lifecycle maps to `member_shutdown`.

Pi-specific substance to preserve from #417 / 0223:

- null-model stamping via `intercom({action:"list"})`; omit-on-null is wrong on Pi because it resolves to `settings.json` `defaultModel`
- Pi canonical model space is provider/model strings, not the Claude enum
- every `subagent(...)` call uses `cwd: <resolved repo root>` as a working-directory concern, not a skill-discovery concern
- file verification remains the completion gate; neither subagent return nor advisory alone advances state
- fresh redispatch remains the default first Pi slice; reuse-advance is friction 9 and deferred

Token expectation: #417 already cut roughly 67 tokens from the Pi runtime refs through positive-binding cleanup. This structural rewrite removes most of the current Pi FO adapter lifecycle narration, so implementation should expect a meaningful size reduction, not just a heading reshuffle.

## Stage Report: ideation

- DONE: Turn the seed task into a concrete implementation plan for converting Codex and Pi FO runtime adapters to `## Runtime implementation` binding blocks keyed by core `«fn»` capability names.
  The body now names exact file edits, target sections, capability bullets, residual sections, and first-test strategy.
- DONE: Decide whether this task must add first-class `«worker.spawn»` and `«worker.shutdown»` headings to `fo-dispatch-core.md`, or whether it should defer that prerequisite to `ad`; record the sequencing recommendation clearly.
  Decision recorded: this task should add minimal core headings first; `ad` keeps broad Claude adapter trim.
- DONE: Specify exactly which Codex and Pi semantics must be preserved in compact form: Codex probe-based binding, wait reinstallation/interruption, feedback reviewer reuse; Pi model stamping, file-verification gate, harness isolation.
  Added `## Semantics to preserve` with explicit Codex and Pi bullets.
- DONE: Define contractlint or focused test changes that will enforce the binding-block shape and prevent lifecycle narration/negative contrast/step-number coupling from returning.
  Added focused contractlint test plan covering binding sections, lifecycle-heading absence, tool-name placement, negative contrast, and step-number coupling.
- DONE: Keep the boundary with `ad` explicit: this task is Codex+Pi only; `ad` owns Claude and broad adapter trim.
  Boundary appears in proposed approach, sequencing decision, and AC-4.
- DONE: Append a Stage Report: ideation and commit only this entity path in the state checkout.
  Report appended here; commit follows as a path-scoped state checkout commit.

### Summary

Ideation converted the seed into an implementation-ready cleanup plan. The key decision is to add minimal `«worker.spawn»` and `«worker.shutdown»` core capability headings in this task, then rewrite only Codex and Pi FO adapters into binding blocks while leaving broad Claude trimming for `ad`.

## Stage Report: implementation

- DONE: Implement the approved t0g plan from the entity body: add minimal first-class `«worker.spawn»` and `«worker.shutdown»` headings to `fo-dispatch-core.md`, then convert only Codex and Pi FO runtime adapters to `## Runtime implementation` binding blocks.
  Code commit `19dde2b8` adds the core headings and rewrites only `codex-first-officer-runtime.md` / `pi-first-officer-runtime.md` into ordered capability binding blocks.
- DONE: Base all work on current `origin/main`, which includes #417 (`bae7cba0`) and runtime-support docs commit `4403e095`; preserve #417 Pi substance and retarget its contractlint guard assertions to the new binding bullets.
  Rebased onto `origin/main` at `54580acc`; history includes `bae7cba0` and `4403e095`, and Pi guard assertions now target binding-block wording.
- DONE: Preserve Codex semantics: probe-based binding, sanitized task mapping, wait timeout/interruption/reinstallation semantics, queued-vs-autonomous notification classification, and feedback reviewer reuse through turn-starting `«addressable-worker»`.
  Guarded by `TestCodexCurrentMultiAgentRuntimeReferencesUseLiveToolSurfaceProbe`, `TestCodexForegroundWaitSectionCarriesOperatorInterruptionShape`, and `TestCodexAndPiFirstOfficerRuntimeSemanticsPreserved`.
- DONE: Preserve Pi semantics: null-model stamping via `intercom({action:"list"})`, provider/model canonical model space, `cwd:<resolved repo root>` on subagent calls, file verification as the completion gate, fresh redispatch default, and harness isolation notes.
  Guarded by `TestCodexAndPiFirstOfficerRuntimeSemanticsPreserved`, Pi tool-placement checks, and the retained Pi negative-contrast tests.
- DONE: Add/update contractlint tests for core capability headings, Codex/Pi binding-block sections, lifecycle-heading absence, negative contrast, step-number coupling, and host tool-name placement.
  Added `runtime_binding_block_test.go`; updated capability extraction and Codex tool placement guards.
- DONE: Keep boundary with `ad`: do not rewrite Claude adapter prose beyond minimal core `→` coverage needed by capability tests, and do not touch live runners or behavior code unless tests require it.
  Claude adapter and live runners were untouched; changes are scoped to FO core docs, Codex/Pi FO adapters, and contractlint tests.
- DONE: Run `go test ./internal/contractlint` and `go test ./...`; run `gofmt -w ./cmd ./internal` if Go tests are edited; include size/line-count impact for Codex+Pi FO adapters.
  `go test ./internal/contractlint`: 59 passed; `go test ./...`: 1712 passed in 17 packages; `go test ./... -race`: 1712 passed in 17 packages; `gofmt -w ./cmd ./internal` run. Codex+Pi FO adapter lines: 164 -> 74, down 90 lines (54.9%).
- DONE: Append a Stage Report: implementation with changed files, token/line impact, verification evidence, residual risk, and commit all worktree changes plus the implementation stage report path-scoped in state.
  Changed files: `fo-dispatch-core.md`, Codex/Pi FO runtime refs, and contractlint tests. Residual risk: compact prose still relies on validation to confirm no reviewer-facing nuance was over-trimmed; `ad` still owns broad Claude adapter trim and final cross-runtime polish.
- DONE: Push the state branch so peers see the entity/report.
  Root resolved the worker's SSH push failure by pushing the code branch and `spacedock-state/dev` over HTTPS successfully.

### Summary

Implementation normalized Codex and Pi FO adapters into binding-block maps keyed to the shared capability names while preserving the load-bearing runtime semantics from #417 and the entity body. The code work is committed as `19dde2b8`; this state report records the verification evidence and line-count reduction.

### Feedback Cycles

- Cycle 1 fix: validation rejected the initial shared-core shape because the new `«worker.spawn»` and `«worker.shutdown»` headings carried concrete Claude/Codex/Pi `→` bindings. Commit `e60979c` keeps those headings first-class but host-neutral in `fo-dispatch-core.md`, and updates contractlint so concrete lifecycle bindings for Codex/Pi are asserted from runtime adapter `## Runtime implementation` blocks instead of new shared-core `→` clauses.
- Verification after cycle 1: `go test ./internal/contractlint` passed with 59 tests; `go test ./...` passed with 1712 tests in 17 packages; `go test ./... -race` passed with 1712 tests in 17 packages. `gofmt -w ./cmd ./internal` was run; unrelated formatting-only diffs outside the task scope were restored before commit.
- Cycle 2 fix: follow-up requested hard Codex idle-stop foreground wait semantics. Commit `892a734` updates the Codex FO runtime so an unresolved Codex worker plus no dispatchable/gate/state work requires `wait_agent(timeout_ms)` before ending the turn or reporting idle/status, and requires reinstalling foreground wait after interrupted waits while preserving final-status mailbox notification as the completion signal.
- Verification after cycle 2: focused Codex wait tests passed with 4 tests; `go test ./internal/contractlint` passed with 60 tests. `gofmt -w ./cmd ./internal` was run; unrelated formatting-only diffs outside the task scope were restored before commit.
- Cycle 3 fix: CI codex-live rejection-flow fresh-dispatched cycle-2 validation because the adapter implied `send_message` plus `followup_task` were both required. Commit `7da3c22` makes `followup_task(target,message)` the current turn-starting reuse/advance route by itself, keeps `send_message` non-triggering only, and requires reusing the completed-but-still-addressable validation reviewer before fresh dispatch; it also states PASSED re-review must re-enter gate flow and advance or terminalize from durable state.
- Verification after cycle 3: focused Codex reviewer-reuse/tool-placement checks passed with 4 tests; `go test ./internal/contractlint` passed with 61 tests; `go test ./...` passed with 1714 tests in 17 packages. `gofmt -w ./cmd ./internal` was run; unrelated formatting-only diffs outside the task scope were restored before commit.
- Cycle 4 fix: fresh CI for PR #418 at `7da3c22a` still failed codex-live rejection-flow because the public Codex live surface exposed only `spawn_agent` and `wait_agent`, with no `followup_task` / `send_message` reuse route available. Commit `aba3010` makes Codex `«addressable-worker»` explicitly probe-based: PRESENT only when a turn-starting route such as `followup_task(target,message)` or a proven equivalent is live, and ABSENT when only spawn/wait exists. The v2 binding remains the PRESENT branch; the ABSENT branch is now characterized as fresh-dispatching a separate cycle-2 validation reviewer, with live-harness assertions requiring two validation spawns and rejecting contradictory reuse tools. The fixture prompt explicitly says "Do not advance the entity to done," so no terminalization change was made for this scenario.
- Verification after cycle 4: `go test ./internal/contractlint` passed with 61 tests; `go test ./internal/ensigncycle` passed with 223 tests; `go test ./...` passed with 1717 tests in 17 packages; `go test ./... -race` passed with 1717 tests in 17 packages. `gofmt -w ./cmd ./internal` was run; unrelated formatting-only diffs outside the task scope were restored before commit.

## Stage Report: validation

- DONE: Review code commit `19dde2b8` on branch `spacedock-ensign/codex-pi-runtime-binding-block-cleanup` in worktree `.worktrees/spacedock-ensign-codex-pi-runtime-binding-block-cleanup` against the t0g entity body, Pi commander feedback, and acceptance criteria.
  Reviewed commit `19dde2b8`; recommendation: PASSED.
- DONE: Verify AC-1: Codex and Pi FO runtime adapters have exactly one `## Runtime implementation` binding block keyed to the expected capability set and no lifecycle narration headings for dispatch/await/reuse/shutdown remain.
  `TestCodexAndPiFirstOfficerRuntimeBindingBlocks` and lifecycle-heading checks passed; manual heading scan found one binding block in each adapter with the expected ordered capability set.
- DONE: Verify AC-2: runtime-support principles are enforced for Codex and Pi, including negative host contrast, mutable step-number coupling, and concrete host tool placement limited to binding/probe/harness sections.
  Contractlint guard tests passed for negative contrast, step coupling, Codex tool placement, and Pi tool placement; host tool names are limited to Codex probe/runtime and Pi runtime/harness sections.
- DONE: Verify AC-3: load-bearing Codex and Pi semantics from the entity body and #417 feedback are preserved in compact form, including Codex wait/reuse semantics and Pi model/cwd/file-verification/fresh-redispatch semantics.
  Guarded by `TestCodexCurrentMultiAgentRuntimeReferencesUseLiveToolSurfaceProbe`, `TestCodexForegroundWaitSectionCarriesOperatorInterruptionShape`, and `TestCodexAndPiFirstOfficerRuntimeSemanticsPreserved`; manual review confirmed #417 Pi shutdown/model/file-verification substance remains.
- DONE: Verify AC-4: the change composes with `ad` by leaving Claude adapter prose and live runners untouched except minimal shared/core test coverage.
  Commit `19dde2b8` touches only FO core docs, Codex/Pi FO runtime refs, and contractlint tests; no Claude adapter rewrite or live runner changes found.
- DONE: Verify AC-5: reproduce or confirm the reported Codex+Pi FO adapter size reduction and ensure binding coverage remains complete.
  Reproduced line counts against `origin/main`: Codex 92 -> 52, Pi 72 -> 22, combined 164 -> 74 lines, down 90 lines (54.9%); binding coverage tests passed.
- DONE: Run `go test ./internal/contractlint` and a broader appropriate suite; report exact commands/results. If you do not run `go test ./... -race`, say why.
  `go test ./internal/contractlint`: 59 passed; `go test ./...`: 1712 passed in 17 packages; `go test ./... -race`: 1712 passed in 17 packages. `gofmt -w ./cmd ./internal` completed but produced unrelated formatting/comment rewrites outside commit `19dde2b8`; validation restored those local changes and committed no code.
- DONE: Append a Stage Report: validation with PASSED or REJECTED recommendation, evidence for each AC, and commit only this entity path in state.
  This report records the PASSED recommendation and AC evidence; state commit follows path-scoped for this entity only.

### Summary

Validation recommends PASSED. The implementation meets the binding-block structure, runtime-support guardrails, preserved Codex/Pi semantics, `ad` boundary, and size-reduction acceptance criteria, with contractlint, full test, and race suites passing.

## Stage Report: validation (cycle 2)

- DONE: Review updated code commit `e60979c9` on branch `spacedock-ensign/codex-pi-runtime-binding-block-cleanup` after feedback rejected concrete host bindings in the shared core.
  Reviewed `e60979c9`; recommendation: PASSED.
- DONE: Verify the new `«worker.spawn»` and `«worker.shutdown»` headings in `fo-dispatch-core.md` are host-neutral and do not add concrete Claude/Codex/Pi tool bindings; existing older transitional `→` lines may remain for `ad` scope.
  `TestDispatchCoreDefinesWorkerLifecycleCapabilities` passed; manual review found no `→` lines or `**Claude:**` / `**Codex:**` / `**Pi:**` bindings inside those two new core blocks.
- DONE: Verify Codex and Pi `## Runtime implementation` blocks carry the concrete bindings for `«worker.spawn»`, `«worker.shutdown»`, and the expected capability set.
  `TestCodexAndPiFirstOfficerRuntimeBindingBlocks` passed; Codex/Pi runtime blocks carry the ordered eight-capability set including concrete spawn and shutdown bindings.
- DONE: Re-verify AC-1 through AC-5 from the entity body, including Codex/Pi preserved semantics, runtime-support guardrails, `ad` boundary, and adapter size reduction.
  AC-1 through AC-5 pass: contractlint guards binding shape, lifecycle-heading removal, host-tool placement, preserved semantics, and line reduction; changed files remain scoped to core docs, Codex/Pi adapters, and contractlint tests.
- DONE: Run `go test ./internal/contractlint` and `go test ./...`; run `go test ./... -race` unless there is a clear reason not to, and report exact commands/results.
  `go test ./internal/contractlint`: 59 passed in 1 package; `go test ./...`: 1712 passed in 17 packages; `go test ./... -race`: 1712 passed in 17 packages.
- DONE: Append a fresh Stage Report: validation cycle 2 with PASSED or REJECTED recommendation, evidence for every AC, and commit only this entity path in state.
  This cycle 2 report records PASSED and will be committed path-scoped to this entity in the state checkout.

### Summary

Validation cycle 2 recommends PASSED. The feedback fix in `e60979c9` keeps new shared-core lifecycle capability headings host-neutral while preserving concrete Codex/Pi bindings in runtime adapter `## Runtime implementation` blocks, and all required test gates passed.

## Stage Report: validation (cycle 3)

- DONE: Review updated code commit `892a734b` after the Codex wait hardening follow-up.
  Reviewed `892a734b`; recommendation: PASSED.
- DONE: Verify Codex wait notes now make unresolved-worker idle waiting a MUST: unresolved Codex worker plus no other dispatchable/gate/state work requires `wait_agent(timeout_ms)` before ending turn or reporting idle/status, and interrupted waits must be reinstalled on the next idle action if unresolved.
  `codex-first-officer-runtime.md` now states the FO MUST call `wait_agent(timeout_ms)` before ending/reporting idle and MUST reinstall foreground wait on the next idle action after interruption when the worker remains unresolved.
- DONE: Verify contractlint guards this hard invariant and does not allow weakening it into optional guidance.
  `TestCodexWaitNotesRequireIdleStopForegroundWait` requires the hard idle-stop phrases and rejects optional forms such as "may use foreground wait" or "use foreground wait only when".
- DONE: Re-verify the previous cycle-2 feedback remains fixed: new `«worker.spawn»` and `«worker.shutdown»` shared-core headings are host-neutral, with concrete Codex/Pi bindings in runtime blocks.
  Manual review confirmed the new shared-core lifecycle headings contain no host `→` lines; `TestCodexAndPiFirstOfficerRuntimeBindingBlocks` passed with concrete Codex/Pi runtime bindings.
- DONE: Run `go test ./internal/contractlint` and a justified broader suite or cite the implementation's broader tests if not rerun; report exact commands/results.
  `go test ./internal/contractlint`: 60 passed in 1 package; `go test ./...`: 1713 passed in 17 packages; `go test ./... -race`: 1713 passed in 17 packages; `gofmt -w ./cmd ./internal` was run and unrelated formatting side effects were restored.
- DONE: Append Stage Report: validation cycle 3 with PASSED or REJECTED recommendation, evidence, and commit only this entity path in state.
  This cycle 3 report records PASSED and will be committed path-scoped to this entity in the state checkout.

### Summary

Validation cycle 3 recommends PASSED. The Codex wait follow-up hardens idle-stop foreground waiting from optional guidance into a guarded MUST, while the cycle-2 host-neutral shared-core lifecycle headings and Codex/Pi runtime binding blocks remain intact.
