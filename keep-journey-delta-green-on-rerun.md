---
id: qz0ap96nt5k93tgbsphq9ahy
title: Keep journey-delta reporting green across failed-job reruns
status: validation
source: PR #639 Runtime Live E2E attempts 2 and 3 on 2026-08-08.
started: 2026-08-09T04:16:06Z
completed:
verdict:
score: 0.75
worktree: .worktrees/spacedock-ensign-keep-journey-delta-green-on-rerun
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:qz0ap96nt5k93tgbsphq9ahy:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:qz0ap96nt5k93tgbsphq9ahy-backlog-1
              briefing:
                id: briefing:qz0ap96nt5k93tgbsphq9ahy:backlog:attempt-1:revision-1
                digest: sha256:e5f04abe6ac0725eac6b831d33309ee789b6bb4ae8888518444ed971971013f3
                request-digest: sha256:b1f7d469585bc21a3d1b8180b8e49253a1abf6ca75347b60ab4874ea87e8de39
                room-ref: ./keep-journey-delta-green-on-rerun/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:qz0ap96nt5k93tgbsphq9ahy:backlog:1
                briefing: briefing:qz0ap96nt5k93tgbsphq9ahy:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T04:15:12.231791Z"
                decision: approve
                reason: Captain directed the First Officer to file and dispatch this journey-metrics fix.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:qz0ap96nt5k93tgbsphq9ahy:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:qz0ap96nt5k93tgbsphq9ahy-ideation-1
              briefing:
                id: briefing:qz0ap96nt5k93tgbsphq9ahy:ideation:attempt-1:revision-1
                digest: sha256:82b8d4b9e7f7b2692f1abb9bf1795c95b4a9af41930f058735fa913c0aecc647
                request-digest: sha256:4d55aafe4bca1976bf75b023791c0698e35787bbd4e241f90bf2696815e02496
                room-ref: ./keep-journey-delta-green-on-rerun/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:qz0ap96nt5k93tgbsphq9ahy:ideation:1
                briefing: briefing:qz0ap96nt5k93tgbsphq9ahy:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T04:29:01.214942Z"
                decision: approve
                reason: Captain approves the exact-artifact recovery design and its warning-and-skip fallback.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:qz0ap96nt5k93tgbsphq9ahy:validation
          stage: validation
          attempts:
            - id: gate-attempt:qz0ap96nt5k93tgbsphq9ahy-validation-1
              briefing:
                id: briefing:qz0ap96nt5k93tgbsphq9ahy:validation:attempt-1:revision-1
                digest: sha256:29bf02a0e2093c1aba011a80b234ad238e3e968137cf5e8edd3199581708dfc8
                request-digest: sha256:d881013b790978df4670360fed0074c52d32b600effd222e99b3c057f3051e4c
                room-ref: ./keep-journey-delta-green-on-rerun/review/validation/briefing-1
---

A failed-job rerun can reuse successful live jobs from an earlier attempt. The reporting job then reads artifacts from different attempts of one run.

All required test jobs can pass in this state. However, one optional metrics download can keep the complete workflow red.

## Problem

Runtime Live E2E run `31288504298` showed the fault. Attempt 2 and attempt 3 selected Codex artifact `9030747663`.

In each attempt, `actions/download-artifact@v5` found the artifact metadata. Then the ZIP download failed after five retries.

The `codex-live` job had passed in attempt 1. The reruns did not run that paid lane again.

The required test result was green. The optional `journey-delta-comment` job made the run red and prevented a clean merge.

## Spike result

The spike used the current run's artifact REST API. The API listed artifacts from successful jobs across earlier attempts.

The spike selected exact artifact IDs and downloaded each ZIP through `gh api`:

- Claude artifact `9031942337` contained 15 metrics JSON files.
- Codex artifact `9030747663` contained 10 metrics JSON files.
- `unzip -t` reported no error for either ZIP.

A name-pattern download also completed. However, the run had two Claude artifacts with the same name.

That path can merge two attempts into one directory. Therefore, the implementation must select one exact artifact ID for each producer name.

## Proposed approach

Replace the two required artifact actions with one best-effort shell step. Keep this change inside `journey-delta-comment`.

The step queries artifacts for `GITHUB_RUN_ID`. It selects the newest nonexpired artifact for each exact producer name.

The selection compares `created_at`. If two times are equal, it compares the artifact IDs.

The producer names are:

- `runtime-live-e2e-claude-live-claude-sonnet-5`
- `runtime-live-e2e-codex-live`

The step downloads each exact ID into a separate temporary directory. It checks each ZIP and requires at least one journey-metrics JSON file per producer.

If both artifact sets are complete, the step writes `found=true`. The existing locate and comment steps then publish one complete comment.

If one query, download, ZIP check, or metrics check fails, the step deletes the partial metrics input. It writes `found=false` and exits successfully.

The warning names the unavailable artifact. For example:

`::warning::journey metrics artifact runtime-live-e2e-codex-live is unavailable; skipping journey-cost delta comment`

Gate the ledger, locate, and comment steps on `found == 'true'`. Thus, no later step can publish a partial comment.

This recovery serves AC-1, AC-2, AC-3, and AC-5. The live spike proves exact-ID recovery from an earlier successful attempt.

The simplest alternative was `continue-on-error` on each current download action. That option protects AC-1 but loses a recoverable complete comment.

Another alternative was a name-pattern `gh run download`. The spike found duplicate Claude names, so that option cannot identify one attempt.

The existing job graph serves AC-4. The task does not add a producer, a live lane, or a rerun command.

## Documentation change

Add this operator rule to `docs/runtime-live-ci.md`:

> The optional journey-delta job uses the newest metrics artifact for each live producer in the run.
>
> If one artifact is unavailable or incomplete, the job warns and skips the comment. The required test result does not change.

No command documentation changes. The `journey-delta` command and the comment layout stay unchanged.

## Expected surface

- `.github/workflows/runtime-live-e2e.yml`: replace two download steps and gate three consumer steps, with at most 65 inserted lines.
- `internal/release/journey_delta_workflow_test.go`: add workflow exercises and graph checks, with at most 170 inserted lines.
- `docs/runtime-live-ci.md`: add the operator rule, with at most 8 inserted lines.

The insertion ceiling is 243 lines across three files. A fourth file or more than 243 inserted lines requires a new design review.

The task changes CI runtime behavior and warning text. It does not change command grammar, stored formats, comment layout, or write authority.

## Out of scope

- Do not change the live journey assertions.
- Do not add a live lane.
- Do not rerun a passed live lane only to produce metrics.
- Do not change the journey-cost format or comment layout.
- Do not change the required status policy.

## Acceptance criteria

**AC-1 (VALUE) - An unavailable metrics artifact cannot make a pull-request workflow red after all required test jobs pass.**
Verified by: execute the real extracted download step with one failed artifact response. The step exits zero, and the comment steps do not run.

**AC-2 - Available artifacts from successful producer attempts result in one complete journey-delta comment.**
Verified by: give the extracted step two valid fixture ZIPs. Both metric sets reach the existing comment command, and exactly one comment call occurs.

**AC-3 - A rerun selects the newest artifact for each exact producer name.**
Verified by: give the API stub two Claude artifacts in reverse order. Only the artifact with the later `created_at` supplies Claude metrics.

**AC-4 - A rerun never publishes a partial journey-delta comment.**
Verified by: fail Claude once and Codex once in separate exercises. Each exercise names the artifact, removes partial input, and makes zero comment calls.

**AC-5 - The fix does not add or rerun a live test job.**
Verified by: parse the workflow graph. It has `offline`, one Claude producer, one Codex producer, and the existing reporting consumer.

**AC-6 - Operators can distinguish a test failure from an optional metrics failure.**
Verified by: the failed-artifact exercise exits successfully and emits one GitHub warning with the exact artifact name.

## Test plan

Add table-driven tests that execute the real extracted workflow script. Use a `gh` stub and small ZIP fixtures.

Cover complete artifacts, duplicate names, a missing artifact, a failed ZIP download, an invalid ZIP, and an empty metrics tree.

Make the stub record artifact IDs and comment calls. These records prove exact selection and the absence of a partial comment.

Keep the existing locate-path and comment-directory tests. Their failure protects the complete-comment path after the new artifact step.

Parse the workflow job graph. Make the test fail if a new live producer or a rerun command appears.

Run `go test ./internal/release`, `go test ./...`, and `go test ./... -race`. Then run `gofmt -w ./cmd ./internal`.

Validation must run a detached adversarial audit. The audit must fail one artifact response and observe a green job with no comment call.

## Stage Report: ideation

- DONE: Define the smallest rerun-safe artifact policy that preserves test truth and complete metrics comments.
  Exact-ID recovery publishes only two complete metric sets. An artifact-specific skip keeps the optional job green.
- DONE: Exercise the riskiest prior-attempt artifact path before selecting recovery or a nonfatal skip.
  The spike downloaded run `31288504298` artifacts `9031942337` and `9030747663`. Both ZIP checks passed.
- DONE: Specify expected files, insertion ceiling, semantic changes, acceptance evidence, and detached audit.
  The task body sets three files, a 243-line ceiling, six tested criteria, and one adversarial audit.

### Summary

The design recovers the newest exact artifact for each producer across rerun attempts. It skips the optional comment with a named warning if recovery cannot produce both complete metric sets.

## Stage Report: implementation

- DONE: Implement exact-ID artifact recovery and the complete-pair warning-and-skip fallback within the approved three-file, 243-insertion ceiling.
  Commit `ce01c4e69` selects exact artifact IDs by creation time and ID. The two-file change has 54 insertions.
- SKIPPED: Add behavioral exercises for complete, duplicate, missing, invalid, and empty artifacts without adding or rerunning a live lane.
  The captain removed the simulator and prohibited a replacement harness. The implementation adds no live lane or rerun command.
- DONE: Update the operator rule and run focused, full, race, formatting, and workflow checks.
  The focused release and registry tests passed. The full and race suites passed after formatting.

### Summary

The reporting job now downloads the newest exact artifact for each required producer. It warns and exits successfully if either metric set is unavailable.

The operator rule now states that an optional metrics failure does not change the required test result.

## Stage Report: validation

- DONE: Verify exact-ID recovery chooses the newest nonexpired Claude and Codex artifacts and skips the optional comment cleanly when either artifact is unavailable or incomplete.
  A detached extracted-script audit selected Claude ID 12 over older ID 11 while excluding expired ID 99 and near-name ID 100; missing, failed-download, invalid-ZIP, and empty-metrics cases exited zero with `found=false` and deleted partial input.
- DONE: Confirm the final change is lean: no artifact simulator, generated ZIP fixtures, job-graph test, live-lane addition, or rerun command.
  Commit `ce01c4e69` changes only the workflow and operator guide, with 54 insertions/13 deletions; no simulator, committed ZIP, test file, new job, live lane, or rerun command exists.
- DONE: Run applicable existing workflow, release, registry, full, race, formatting, and detached-adversarial checks; verify every acceptance criterion and report PASSED or REJECTED.
  Focused release and registry checks passed; `go test ./...` and `go test ./... -race` passed; `gofmt -w ./cmd ./internal`, `gofmt -l`, and `git diff --check` were clean; verdict is PASSED.
- DONE: AC-1 (VALUE) - An unavailable metrics artifact cannot make a pull-request workflow red after all required test jobs pass.
  The detached audit ran the real recovery step with missing and failed artifact responses: each exited zero, wrote only `found=false`, removed partial input, and therefore left all three consumers gated off.
- DONE: AC-2 - Available artifacts from successful producer attempts result in one complete journey-delta comment.
  Two valid exact-ID ZIPs produced `found=true`; the real locate step delivered two JSON inputs, and `TestJourneyDeltaCommandCreatesNewCommentWhenNoneExists` proved exactly one post invocation (it would fail on zero or multiple calls).
- DONE: AC-3 - A rerun selects the newest artifact for each exact producer name.
  Reverse-ordered API data downloaded only Claude ID 12 (later `created_at`) plus Codex ID 21; the older, expired, and nonexact-name Claude artifacts were not downloaded.
- DONE: AC-4 - A rerun never publishes a partial journey-delta comment.
  Separate Claude-missing and Codex failure/incomplete exercises wrote `found=false`, removed recovered Claude input when applicable, and the independently parsed consumer predicates require `found == 'true'`.
- DONE: AC-5 - The fix does not add or rerun a live test job.
  Parsed active job headers remain exactly `offline`, `claude-live`, `codex-live`, and `journey-delta-comment`; no rerun invocation exists and the diff adds no job header.
- DONE: AC-6 - Operators can distinguish a test failure from an optional metrics failure.
  Every unavailable case exited zero and emitted exactly one GitHub warning containing the exact failing producer name; changing the warning count or name fails the detached audit.
- SKIPPED: Add durable regression coverage for the new recovery branch.
  Latest captain feedback explicitly removed the simulator/harness and required no replacement; detached mutation showed committed release tests stay green if `warn_and_skip` changes to exit 1, while the independent audit fails, so revisit if this branch is edited again.

### Summary

PASSED. The candidate satisfies all six acceptance criteria, stays inside the two-file lean scope, and completes focused, full, race, formatting, and detached behavioral validation.

One deferred evidence risk remains by captain-approved scope: the recovery branch has no committed behavioral regression test. Promote it to material if a later change touches artifact recovery without repeating the extracted-script audit or if CI exhibits another optional-comment red result.
