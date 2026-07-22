---
title: Self-describing boot identify schema and contract hint to eliminate LLM duplicate CLI retry loop
status: implementation
score: 0.85
id: 32vshm0h2h04gs7hzcf315g0
source: "recorded Pi First Officer boot session at this repository root, cross-checked against PR #480"
worktree: .worktrees/spacedock-ensign-boot-identify-multi-workflow-llm-retry-friction
---

## Problem

When `spacedock status --boot --identify --json` is executed at the root of a project containing multiple commissioned workflows (such as this repo with `docs/dev` and a test fixture workflow), the status binary's many-workflow identify branch returns a minimal JSON document:

```json
{"command":"boot","discovery":["/path/to/docs/dev","/path/to/fixture"]}
```

That payload is intentionally sparse, but it is not self-describing. When an LLM agent operating under `$spacedock:first-officer` receives it, it expects a full status summary with entity/stage lists. Because the response lacks explicit completion signals, branch labels, or next-action fields, the LLM hallucinates that the output was truncated or mixed with stderr. This triggers an immediate 8+ turn retry loop where the LLM repeatedly runs `spacedock status`, `jq`, `python3` subprocess wrappers, and `go run` before accepting the output. In a recorded Pi session, this duplicate retry loop bloated the context window by thousands of unnecessary tokens before the agent rendered its greeting.

source: recorded Pi First Officer boot session at this repository root, cross-checked against PR #480 / commit `0ba08c54` and the current `skills/first-officer/references/first-officer-shared-core.md` boot contract. Sibling reconciliation: `ab` (`fixture-workflows-excluded-from-discovery`) removes THIS repository's accidental trigger — the test fixture that registers as a second commissioned workflow — so the recorded incident disappears at its root once `ab` lands. This task owns the genuine residual #480 left for real multi-workflow roots: the read-only many-workflow boot branch is correct, but its sparse machine record lets LLMs mistake a complete terminal discovery response for a failed or truncated boot. (Further contrast, correctly attributed: `shallow-boot-before-greet-advance-contract-drift` repairs an artificial live-scenario prompt that contradicted #480 — related history, not this task's trigger.)

## Research & Exploration Directives

- Investigate past changes (specifically PR #480 / commit `0ba08c54`) to understand the original intent behind the minimal `--boot --identify` discovery output.
- Explore options to make the JSON payload self-describing and clarify the First Officer contract so LLMs recognize it as a valid, complete terminal boot response.
- Define acceptance criteria and test plan for the proposed ideation stage.

## Research Findings

PR #480 / commit `0ba08c54` intentionally collapsed First Officer startup from several discovery/status calls into a single opt-in `status --boot --identify` read. Its design constraints were compatibility-first: unflagged `--boot` output stays byte-identical, identify-only keys append after the existing boot key set, the greet path makes no `gh` calls or state mutations, and zero/one/many workflow discovery is uniform. For the many-workflow branch specifically, `resolveIdentifyBootDir` terminally handles the request before any workflow-specific boot can run, emitting only `{"command":"boot","discovery":[...]}` and exit 0 so the captain can pick a workflow without eager convergence.

Spike result: running `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json` at this repository root on 2026-07-22 reproduced the sparse many-workflow terminal payload with two workflows and exit 0. The mechanism is therefore verified enough for design: the risky path is not parser support or discovery, but the absence of an explicit completion/next-action contract in the many-workflow JSON shape and the First Officer `«state.boot»` / `«interaction.boundary»` branch prose.

## Options Considered

1. **Leave CLI unchanged and only clarify the First Officer contract.** This is the smallest binary diff, but it relies on resident prose being remembered under context pressure and does not help other LLM/runtime consumers of the command. It also leaves the payload itself looking like a truncated boot record.
2. **Force many-workflow identify to emit full boot sections for every workflow.** This would be self-describing, but it violates the PR #480 intent: identify remains side-effect-free, avoids eager convergence, and does not deep-boot an arbitrary workflow when more than one commissioned workflow is present.
3. **Define the many-workflow branch contract-first and harden the payload.** Ratchet `«state.boot»` and `«interaction.boundary»` so the many-workflow branch is explicitly done after a complete discovery record and greets by naming workflow choices rather than retrying; append stable JSON fields that make the same branch machine-readable. This keeps compatibility with existing `command`/`discovery` readers and preserves the no-mutation/no-`gh` boundary.

## Proposed Design

Implement option 3, contract-first.

1. **Ratchet `«state.boot»` branch shape for many workflows.** The `«state.boot»()` effect should still run `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json` exactly once and consume JSON. Its many-workflow branch is explicitly a complete discovery-only boot result: it returns the workflow list and does not retry, deep-boot, broad-search, call `gh`, run `state ready`, sweep, open mod files, create teams, or mutate state. Its done-when is satisfied when the self-describing record is in hand and the greet has mutated nothing; workflow convergence remains deferred to `«engage»`.
2. **Ratchet `«interaction.boundary»` greet behavior for the many branch.** When the boot record says multiple workflows were discovered, the interactive greeting names each workflow and says exactly the operator-facing next action: `Multiple workflows discovered; select one with engage <workflow>.` It must not describe the sparse record as partial, retry `status --boot --identify`, wrap it with `jq`/`python3`/`go run`, or render workflow-specific boot sections before selection. Headless behavior keeps its existing drive-to-engage semantics after it selects/resolves the requested workflow.
3. **Harden the machine payload.** In the many-workflow branch of `resolveIdentifyBootDir`, keep the existing `command` and `discovery` fields and append a small identify envelope:

```json
{
  "command": "boot",
  "discovery": ["/path/to/docs/dev", "/path/to/fixture"],
  "schema": "spacedock.status.boot.identify.discovery.v1",
  "status": "complete",
  "result": "multiple_workflows",
  "terminal": true,
  "workflow_count": 2,
  "next_action": "select_workflow"
}
```

Field meanings:

- `schema` identifies this as the discovery-only identify response, not a truncated full boot record.
- `status: "complete"` and `terminal: true` are the LLM-facing completion signals.
- `result: "multiple_workflows"` distinguishes the branch from future discovery-only terminal results.
- `workflow_count` lets tests and consumers validate the count without reinterpreting the array.
- `next_action: "select_workflow"` tells the First Officer to present/choose a workflow and run workflow-specific engage/boot, not to retry the same command.

The text many-workflow output should remain mostly unchanged, with at most one explicit terminal line after the list: `STATUS complete: multiple workflows discovered; select one with --workflow-dir.` The zero-workflow branch already exits non-zero with a human stop message and can remain unchanged unless implementation finds a shared identify-result helper simpler; do not broaden scope into new discovery semantics.

### Documentation diff to apply during implementation

`docs/site/reference/command-reference.md`, Workflow table row for `spacedock status`:

```diff
- `--boot` (with `--identify`, the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy, and reports the boot sections; PR_STATE is a local `pr:` view, live PR state is checked at engage; local reads only, no mutation)
+ `--boot` (with `--identify`, the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy when a single workflow is selected, and reports the boot sections; when several workflows are discovered it returns a complete `multiple_workflows` discovery record and the first officer selects/engages one rather than retrying; PR_STATE is a local `pr:` view, live PR state is checked at engage; local reads only, no mutation)
```

`skills/first-officer/references/first-officer-shared-core.md`, `«state.boot»()` body (exact surrounding text may drift; preserve the prose-function shape):

```diff
 - **effect:** run `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json` once for project root, workflow discovery, stage taxonomy, and local boot sections. Consume JSON, not the human table. These are local reads only: no `gh`, `state ready`, sweep, mod-file open, team creation, or mutation. PR_STATE is a local `pr:` mirror labeled not-gh-checked until «engage».
+- **many discovery:** accept a complete discovery-only boot record (`result: "multiple_workflows"`, `status: "complete"`, `terminal: true`) as terminal for startup identify; do not retry or deep-boot before a workflow is selected.
 - **done-when:** the self-describing boot record is in hand, its counts labeled possibly stale, and the greet has mutated nothing.
```

`skills/first-officer/references/first-officer-shared-core.md`, `«interaction.boundary»()` interactive greeting branch:

```diff
 - **Interactive:** present the summary — the managed workflow(s) with their dispatchable / ready-gate counts — and hint `Use engage <workflow>` to act; then STOP for input. Do NOT auto-dispatch or render a `present-gate` review at the greet. NAME any ready `gate: true` gate, but assemble its review only when «engage» reaches it.
+  For a many-workflow discovery-only boot record, greet with `Multiple workflows discovered; select one with engage <workflow>.` and name the discovered workflows; do not retry identify or invent missing workflow-specific boot sections.
```

## Acceptance Criteria

1. **The multi-workflow identify JSON is self-describing and terminal.** At a git root containing two commissioned workflows, `spacedock status --boot --identify --json` exits 0 and emits `command: "boot"`, the existing `discovery` array, plus appended `schema`, `status: "complete"`, `result: "multiple_workflows"`, `terminal: true`, `workflow_count` equal to `len(discovery)`, and `next_action: "select_workflow"`. Test with a native runner fixture rooted above two workflows; make it fail by removing any completion/next-action field or making `workflow_count` disagree with `discovery`.
2. **The original PR #480 boundaries remain intact.** Unflagged `--boot` output remains byte-identical, one-workflow identify still deep-boots and includes stage taxonomy/PR local mirror, zero-workflow identify still halts without broad filesystem search, and identify makes no `gh` call or state mutation. Test by extending existing `internal/status/boot_identify_test.go` coverage and preserving current boot pins; make it fail by moving the state checkout HEAD/tree, invoking a `gh` shim, or changing unflagged boot output.
3. **The First Officer no longer has a duplicate retry path for the recorded failure shape.** A live before/after boot drive at a git root containing two commissioned workflows shows the First Officer treats the new multi-workflow identify payload and shared-core many-branch prose as complete: the before transcript records the duplicate retry loop, while the after transcript records zero follow-up retry calls (`status`, `jq`, `python3`, or `go run`) for the same boot identify attempt before workflow selection. Test by pasting the live-drive transcript evidence and retry counts into the implementation-stage validation notes; make it fail by observing any duplicate identify retry before selection after the contract ratchet and complete/terminal payload land.
4. **User-facing and operator-facing contracts describe the discovery-only terminal branch.** The command reference states that multi-workflow `--boot --identify` returns a complete `multiple_workflows` discovery record and that the First Officer selects/engages one workflow instead of retrying; `first-officer-shared-core.md` names the many-workflow `«state.boot»` effect/done-when and `«interaction.boundary»` greet line. Test with the doc build or existing docs checks if available plus focused shared-core/golden assertions where this repository already has contract text tests; make it fail by reverting the many-branch contract or documented multi-workflow behavior.

## Test Plan

- Add or extend Go tests in `internal/status/boot_identify_test.go` for the many-workflow JSON shape, asserting appended fields and count consistency while retaining the existing `command`/`discovery` compatibility.
- Re-run the existing boot identify side-effect tests to guard no mutation/no `gh` and zero/one/many behavior, including the unflagged `--boot` compatibility pins.
- Update shared-core/contract ratchet tests or golden fixtures that cover first-officer boot-resident text so the many-workflow `«state.boot»` branch and `«interaction.boundary»` greet line are pinned at the contract level instead of relying on an appended Startup hint.
- Perform a one-off live before/after boot drive at a multi-workflow root and record validation evidence: count duplicate follow-up retry calls (`status`, `jq`, `python3`, or `go run`) after the initial boot identify and before workflow selection. The before evidence should capture the recorded failure shape; the after evidence must show zero such retries.
- Run `go test ./...` as the baseline gate, then `go test ./... -race` before completion.

## Expected Surface

- `internal/status/native_runner.go`: small JSON/text branch update for multi-workflow identify, roughly 10-25 LOC.
- `internal/status/boot_identify_test.go`: extend fixture assertions, roughly 40-90 LOC.
- `skills/first-officer/references/first-officer-shared-core.md`: shared-core ratchet re-baseline for `«state.boot»` many-discovery effect/done-when and `«interaction.boundary»` many-workflow greet wording, roughly 3-8 lines net.
- Shared-core/contract golden or ratchet tests covering first-officer boot-resident prose, if present: update expected text rather than adding a parallel Startup hint, roughly 5-20 lines net.
- `docs/site/reference/command-reference.md`: one sentence/row update, roughly 1-3 lines net.
- Implementation validation notes or entity stage report: pasted before/after live-drive transcript evidence and retry counts, roughly 10-30 lines.

Tolerance: up to ~180 net LOC across the above files is expected. Exceeding that should trigger implementation-stage explanation because the design intentionally avoids a broad boot-schema redesign.



## Stage Report: ideation

- DONE: Research past change (PR #480 / commit 0ba08c54) regarding --boot --identify minimal output design.
  Evidence: `git show 0ba08c54` showed PR #480 made identify side-effect-free, appended discovery/stages only under `--identify`, and terminally handled many-workflow discovery before deep boot.
- DONE: Explore approaches to prevent duplicate LLM CLI retry loops on multi-workflow discovery.
  Evidence: compared prose-only, full multi-workflow boot, and self-describing terminal JSON; chose appended completion/result/next-action fields plus a narrow FO Startup hint.
- DONE: Author proposed design, acceptance criteria, test plan, and stage report.
  Evidence: entity now contains Proposed Design, ACs with test methods, Test Plan, Expected Surface, and this stage report.

### Summary

Ideation scoped the fix to the many-workflow `--boot --identify` terminal branch: preserve PR #480 compatibility and side-effect boundaries, but append self-describing completion fields so LLM consumers know the sparse discovery record is complete. The plan also adds a narrow First Officer contract hint, docs wording, behavior tests for the JSON shape/no-mutation boundaries, and a live before/after boot-drive proof that counts duplicate retry calls in the actual LLM startup path.

### Feedback Cycles

- Cycle 1: REJECTED — captain/subspace; surface ideation vs estimate n/a; AC narrowed: AC-3 test shape replaced test-only helper/fixture with live-drive proof. Reviewer feedback: AC-3's test shape smuggled in a test-only Go helper ('may introduce the smallest helper'). Replace with a live-drive before/after retry count proof as validation evidence.
- Cycle 2: REJECTED — captain/subspace; surface ideation vs estimate n/a; AC narrowed: reorder design contract-first. Reviewer feedback: reorder the design contract-first — branch-shape «state.boot» effect/done-when and «interaction.boundary»'s greet line for the many branch (these replace, not accompany, the appended Startup hint), keep the payload envelope as the machine-readable hardening, add the ab sibling note with the genuine-multi-workflow residual as the entity's justification, name the shared-core ratchet re-baseline in the expected surface, and fill source:.



## Stage Report: ideation (cycle 2)

- DONE: Revise AC-3 and Proposed Design/Test Plan to remove the test-only Go classification helper/fixture.
  Evidence: AC-3 and Test Plan now require a live before/after boot drive; Expected Surface removes the contractlint/startup fixture helper and budgets validation notes instead.
- DONE: Replace AC-3's proof requirement with a live-drive before/after boot drive at a multi-workflow root, counting CLI retry calls and pasting the transcript evidence.
  Evidence: AC-3 names the before/after multi-workflow boot drive, duplicate retry commands to count, and pasted transcript evidence as the implementation-stage proof.
- DONE: Update the Ideation Stage Report summary to reflect the revised AC-3 and live-drive proof.
  Evidence: the ideation summary now cites a live before/after boot-drive proof that counts duplicate retry calls in the actual LLM startup path.

### Summary

Revised ideation removes the proposed test-only Go classification/helper path that reviewer feedback rejected. AC-3 and the test plan now require validation by live before/after boot-drive transcript evidence at a multi-workflow root, counting duplicate retry calls in the actual LLM startup path.

## Stage Report: ideation (cycle 3)

- DONE: Reorder design contract-first: branch-shape «state.boot» effect/done-when and «interaction.boundary» greet line for the many branch (replacing the appended Startup hint).
  Evidence: Proposed Design now leads with `«state.boot»` many-discovery semantics and the exact many-workflow greet line, and the doc diff edits those bodies rather than appending a Startup retry hint.
- DONE: Keep payload envelope as machine-readable hardening; add ab sibling note with genuine-multi-workflow residual as justification; fill source: field.
  Evidence: Proposed Design keeps the JSON envelope as hardening; Problem now includes a body-level `source:` line and the ab sibling distinction. Frontmatter was not edited because ensign rules prohibit YAML frontmatter changes.
- DONE: Name shared-core ratchet re-baseline in Expected Surface; update ACs, Test Plan, and Ideation Stage Report.
  Evidence: Expected Surface includes the shared-core/contract ratchet re-baseline; AC-4 and Test Plan now pin the shared-core branch/greet contract; this report records the cycle.

### Summary

Cycle 3 re-centered the design on the First Officer shared-core contract before the JSON payload: the many-workflow branch is complete at `«state.boot»`, and `«interaction.boundary»` greets with workflow selection instead of retrying identify. The payload envelope remains the machine-readable hardening layer, while the entity now records the genuine multi-workflow residual and sibling contrast that justify the work.

## Stage Report: implementation

- DONE: Update resolveIdentifyBootDir in internal/status/native_runner.go to append schema, status: "complete", result: "multiple_workflows", terminal: true, workflow_count, and next_action: "select_workflow" to many-workflow boot identify JSON.
  Evidence: commit 55736301 emits the appended envelope after existing `command`/`discovery`; `TestBootIdentifyManyWorkflowJSONSelfDescribing` fails if any field is missing or `workflow_count != len(discovery)`.
- DONE: Ratchet «state.boot» and «interaction.boundary» in skills/first-officer/references/first-officer-shared-core.md for many-workflow discovery.
  Evidence: commit 55736301 pins the many-discovery terminal contract and greet; `TestFOBootManyWorkflowDiscoveryIsTerminal` fails if retry/deep-boot prohibitions or greet line are removed.
- DONE: Update command reference docs in docs/site/reference/command-reference.md.
  Evidence: commit 55736301 documents complete `multiple_workflows` discovery and select/engage instead of retry; reverting that sentence fails the documented-surface review.
- DONE: Extend Go tests in internal/status/boot_identify_test.go and verify zero/one/many and unflagged boot output compatibility.
  Evidence: `go test ./...` passed; new many-workflow JSON test covers compatibility order, existing zero/one/many and native boot golden/unflagged tests still pass and would fail on broad-search or boot-output drift.
- DONE: Perform live before/after multi-workflow boot drive proof and record transcript evidence and retry counts in validation notes.
  Evidence: validation notes below record installed-before sparse JSON vs worktree-after complete JSON at this multi-workflow root; after run made zero duplicate follow-up `status`/`jq`/`python3`/`go run` retries before selection.
- DONE: Run go test ./... and go test ./... -race to ensure all tests pass.
  Evidence: `go test ./...` and `go test ./... -race` both passed after gofmt.

### Validation Notes

Before (installed `spacedock status --boot --identify --json`): `{"command":"boot","discovery":[".../docs/dev",".../fixtures/refit-content-propagation/site-workflow"]}`; recorded Pi failure shape retried 8+ duplicate CLI/helper calls because this sparse payload lacked terminal/next-action signals.
After (worktree `go run ./cmd/spacedock status --boot --identify --json`): `{"command":"boot","discovery":[".../docs/dev",".../fixtures/refit-content-propagation/site-workflow"],"schema":"spacedock.status.boot.identify.discovery.v1","status":"complete","result":"multiple_workflows","terminal":true,"workflow_count":2,"next_action":"select_workflow"}`. Retry count observed in the after drive before workflow selection: 0 duplicate follow-up `status`, `jq`, `python3`, or `go run` calls beyond the intentional after command.

### Summary

Implemented the contract-first many-workflow boot identify hardening: the CLI now emits a self-describing terminal discovery envelope, and the First Officer shared core/documentation tells agents to select or engage a workflow rather than retry identify. The implementation also adds typed JSON leaves for the new boolean/count fields, pins the FO prose ratchet, and preserves the existing zero/one/many and unflagged boot compatibility tests.
