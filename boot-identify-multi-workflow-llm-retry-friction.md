---
title: Self-describing boot identify schema and contract hint to eliminate LLM duplicate CLI retry loop
status: ideation
score: 0.85
id: 32vshm0h2h04gs7hzcf315g0
---

## Problem

When `spacedock status --boot --identify --json` is executed at the root of a project containing multiple commissioned workflows (such as this repo with `docs/dev` and a test fixture workflow), the status binary's `resolveIdentifyBootDir` function returns a minimal JSON document:

```json
{"command":"boot","discovery":["/path/to/docs/dev","/path/to/fixture"]}
```

When an LLM agent operating under `$spacedock:first-officer` receives this short payload, it expects a full status summary with entity/stage lists. Because the response lacks explicit completion signals or descriptive status fields, the LLM hallucinates that the output was truncated or mixed with stderr. This triggers an immediate 8+ turn retry loop where the LLM repeatedly runs `spacedock status`, `jq`, `python3` subprocess wrappers, and `go run` before accepting the output. In a recorded Pi session, this duplicate retry loop bloated the context window by thousands of unnecessary tokens before the agent rendered its greeting.

## Research & Exploration Directives

- Investigate past changes (specifically PR #480 / commit `0ba08c54`) to understand the original intent behind the minimal `--boot --identify` discovery output.
- Explore options to make the JSON payload self-describing (or clarify the First Officer contract) so LLMs recognize it as a valid, complete terminal boot response.
- Define acceptance criteria and test plan for the proposed ideation stage.

## Research Findings

PR #480 / commit `0ba08c54` intentionally collapsed First Officer startup from several discovery/status calls into a single opt-in `status --boot --identify` read. Its design constraints were compatibility-first: unflagged `--boot` output stays byte-identical, identify-only keys append after the existing boot key set, the greet path makes no `gh` calls or state mutations, and zero/one/many workflow discovery is uniform. For the many-workflow branch specifically, `resolveIdentifyBootDir` terminally handles the request before any workflow-specific boot can run, emitting only `{"command":"boot","discovery":[...]}` and exit 0 so the captain can pick a workflow without eager convergence.

Spike result: running `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json` at this repository root on 2026-07-22 reproduced the sparse many-workflow terminal payload with two workflows and exit 0. The mechanism is therefore verified enough for design: the risky path is not parser support or discovery, but the absence of an explicit completion/next-action contract in the many-workflow JSON shape and the First Officer startup prose.

## Options Considered

1. **Leave CLI unchanged and only clarify the First Officer skill.** This is the smallest diff, but it relies on resident prose being remembered under context pressure and does not help other LLM/runtime consumers of the command. It also leaves the payload itself looking like a truncated boot record.
2. **Force many-workflow identify to emit full boot sections for every workflow.** This would be self-describing, but it violates the PR #480 intent: identify remains side-effect-free, avoids eager convergence, and does not deep-boot an arbitrary workflow when more than one commissioned workflow is present.
3. **Make the many-workflow terminal payload self-describing and add a narrow First Officer contract hint.** Append stable JSON fields after `command` and `discovery` that say the response is complete, terminal for startup, and requires workflow selection; then teach Startup to treat that outcome as success rather than retry. This keeps compatibility with existing `command`/`discovery` readers and preserves the no-mutation/no-`gh` boundary.

## Proposed Design

Implement option 3. In the many-workflow branch of `resolveIdentifyBootDir`, keep the existing `command` and `discovery` fields and append a small identify envelope:

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

Update the First Officer Startup contract narrowly: after running `status --boot --identify --json`, if the JSON has `result: "multiple_workflows"`, `status: "complete"`, and `terminal: true`, accept it as the complete terminal identify result, do not retry or wrap the command, and proceed to workflow selection/engage. Keep the existing one-workflow deep boot path unchanged.

### Documentation diff to apply during implementation

`docs/site/reference/command-reference.md`, Workflow table row for `spacedock status`:

```diff
- `--boot` (with `--identify`, the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy, and reports the boot sections; PR_STATE is a local `pr:` view, live PR state is checked at engage; local reads only, no mutation)
+ `--boot` (with `--identify`, the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy when a single workflow is selected, and reports the boot sections; when several workflows are discovered it returns a complete `multiple_workflows` discovery record and the first officer selects/engages one rather than retrying; PR_STATE is a local `pr:` view, live PR state is checked at engage; local reads only, no mutation)
```

`skills/first-officer/references/first-officer-shared-core.md`, Startup step for local identify (exact surrounding text may drift; preserve the existing step count):

```diff
- Run `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json` to fold discovery, stage taxonomy, and local PR mirror into boot.
+ Run `${SPACEDOCK_BIN:-spacedock} status --boot --identify --json` to fold discovery, stage taxonomy, and local PR mirror into boot. If the JSON result is `multiple_workflows` with `status: "complete"` and `terminal: true`, treat that sparse discovery record as the complete boot-identify response, do not retry the command, and select/engage one workflow.
```

## Acceptance Criteria

1. **The multi-workflow identify JSON is self-describing and terminal.** At a git root containing two commissioned workflows, `spacedock status --boot --identify --json` exits 0 and emits `command: "boot"`, the existing `discovery` array, plus appended `schema`, `status: "complete"`, `result: "multiple_workflows"`, `terminal: true`, `workflow_count` equal to `len(discovery)`, and `next_action: "select_workflow"`. Test with a native runner fixture rooted above two workflows; make it fail by removing any completion/next-action field or making `workflow_count` disagree with `discovery`.
2. **The original PR #480 boundaries remain intact.** Unflagged `--boot` output remains byte-identical, one-workflow identify still deep-boots and includes stage taxonomy/PR local mirror, zero-workflow identify still halts without broad filesystem search, and identify makes no `gh` call or state mutation. Test by extending existing `internal/status/boot_identify_test.go` coverage and preserving current boot pins; make it fail by moving the state checkout HEAD/tree, invoking a `gh` shim, or changing unflagged boot output.
3. **The First Officer no longer has a duplicate retry path for the recorded failure shape.** A live before/after boot drive at a git root containing two commissioned workflows shows the First Officer treats the new multi-workflow identify payload as complete: the before transcript records the duplicate retry loop, while the after transcript records zero follow-up retry calls (`status`, `jq`, `python3`, or `go run`) for the same boot identify attempt before workflow selection. Test by pasting the live-drive transcript evidence and retry counts into the implementation-stage validation notes; make it fail by observing any duplicate identify retry before selection after the new complete/terminal hint lands.
4. **User-facing documentation describes the discovery-only terminal branch.** The command reference states that multi-workflow `--boot --identify` returns a complete `multiple_workflows` discovery record and that the First Officer selects/engages one workflow instead of retrying. Test with the doc build or existing docs checks if available, plus a focused assertion only if this repository already uses command-reference contract tests; make it fail by reverting the documented multi-workflow behavior.

## Test Plan

- Add or extend Go tests in `internal/status/boot_identify_test.go` for the many-workflow JSON shape, asserting appended fields and count consistency while retaining the existing `command`/`discovery` compatibility.
- Re-run the existing boot identify side-effect tests to guard no mutation/no `gh` and zero/one/many behavior.
- Perform a one-off live before/after boot drive at a multi-workflow root and record validation evidence: count duplicate follow-up retry calls (`status`, `jq`, `python3`, or `go run`) after the initial boot identify and before workflow selection. The before evidence should capture the recorded failure shape; the after evidence must show zero such retries.
- Run `go test ./...` as the baseline gate, then `go test ./... -race` before completion.

## Expected Surface

- `internal/status/native_runner.go`: small JSON/text branch update for multi-workflow identify, roughly 10-25 LOC.
- `internal/status/boot_identify_test.go`: extend fixture assertions, roughly 40-90 LOC.
- `skills/first-officer/references/first-officer-shared-core.md`: one narrow Startup hint, roughly 1-4 lines net.
- `docs/site/reference/command-reference.md`: one sentence/row update, roughly 1-3 lines net.
- Implementation validation notes or entity stage report: pasted before/after live-drive transcript evidence and retry counts, roughly 10-30 lines.

Tolerance: up to ~160 net LOC across the above files is expected. Exceeding that should trigger implementation-stage explanation because the design intentionally avoids a broad boot-schema redesign.



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


## Stage Report: ideation (cycle 2)

- DONE: Revise AC-3 and Proposed Design/Test Plan to remove the test-only Go classification helper/fixture.
  Evidence: AC-3 and Test Plan now require a live before/after boot drive; Expected Surface removes the contractlint/startup fixture helper and budgets validation notes instead.
- DONE: Replace AC-3's proof requirement with a live-drive before/after boot drive at a multi-workflow root, counting CLI retry calls and pasting the transcript evidence.
  Evidence: AC-3 names the before/after multi-workflow boot drive, duplicate retry commands to count, and pasted transcript evidence as the implementation-stage proof.
- DONE: Update the Ideation Stage Report summary to reflect the revised AC-3 and live-drive proof.
  Evidence: the ideation summary now cites a live before/after boot-drive proof that counts duplicate retry calls in the actual LLM startup path.

### Summary

Revised ideation removes the proposed test-only Go classification/helper path that reviewer feedback rejected. AC-3 and the test plan now require validation by live before/after boot-drive transcript evidence at a multi-workflow root, counting duplicate retry calls in the actual LLM startup path.
