---
id: qz0ap96nt5k93tgbsphq9ahy
title: Keep journey-delta reporting green across failed-job reruns
status: implementation
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
