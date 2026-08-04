---
title: Land the live-test-truth planning package
status: implementation
source: "Captain-directed topology prerequisite before 3d, 2026-08-04"
score: 1.0
sprint: live-test-truth
group: truthful-results
sprint-readiness: ready
id: nz6816mcf0qsdp1bsh3wzp1w
gates:
    version: 1
    records:
        - id: gate:nz6816mcf0qsdp1bsh3wzp1w:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:nz6816mcf0qsdp1bsh3wzp1w-backlog-1
              briefing:
                id: briefing:nz6816mcf0qsdp1bsh3wzp1w:backlog:attempt-1:revision-1
                digest: sha256:c7c0ad74b08f72262c252d3cb0ac737f78564490e45e1ead7d2607caf4173754
                request-digest: sha256:e121dd798c26230dd3488085f80aaee58276ffbc06722f9c7228b5480b2c8ab6
                room-ref: ./land-live-test-truth-planning-package/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nz6816mcf0qsdp1bsh3wzp1w:backlog:1
                briefing: briefing:nz6816mcf0qsdp1bsh3wzp1w:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-04T15:00:27.945095Z"
                decision: approve
                reason: Captain/root explicitly directed a separate prerequisite PR for these five planning and registry paths so 3d remains within its approved nine-file surface; delegated conn authorizes this admission.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:nz6816mcf0qsdp1bsh3wzp1w:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:nz6816mcf0qsdp1bsh3wzp1w-ideation-1
              briefing:
                id: briefing:nz6816mcf0qsdp1bsh3wzp1w:ideation:attempt-1:revision-1
                digest: sha256:725dc09060b05ec0be4e455727fd67c0cb8b882fb4816fd39aeac6b9423a5c97
                request-digest: sha256:d2f6308f1920b378be6ec434e6ac29b6dfd402e2089b44af569279f69ffa3503
                room-ref: ./land-live-test-truth-planning-package/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nz6816mcf0qsdp1bsh3wzp1w:ideation:1
                briefing: briefing:nz6816mcf0qsdp1bsh3wzp1w:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-04T15:08:32.089823Z"
                decision: approve
                reason: Delegated Captain conn authorizes this topology prerequisite. Ideation completed 4 DONE, 0 SKIPPED, 0 FAILED; AC-1 through AC-3 each have direct diff/topology evidence, and the approved surface is exactly five files +776/-7 with no other drift.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-04T15:03:51Z
worktree: .worktrees/spacedock-ensign-land-live-test-truth-planning-package
---

Land the already-authored sprint plan separately so implementation members retain their approved product surfaces.

## Outcome

The live-test-truth sprint index, staff review, execution package, desired registry, and durable-decision links are available on `main` before `3d` opens its implementation PR.

## Problem

The five planning and registry paths were authored and committed before sprint execution, but they remained on unpushed local ancestry. The `3d` branch inherited them, making its current diff 14 files instead of its approved nine-file implementation surface.

Merging those files as part of `3d` would misattribute 777 insertions and seven deletions to an implementation member whose approved expected surface does not own them. Dropping them would lose the sprint contract and desired-state registry.

## Proposed approach

Use the existing final planning commit `b0a2541f656e0b6d72d01d411160c1694910cbee` as one parent and merge current `origin/main` as the other in a dedicated branch. The preflight merge-tree is conflict-free and preserves both the substantive authored planning history and all current-main work.

Remove only the trailing blank line at EOF in `docs/roadmap/live-test-truth/staff-review.md`. The authored tree otherwise remains byte-identical to the planning tip. This is the only authorized drift from that tree.

Open and merge a separate documentation PR. Then merge the resulting current `origin/main` into `3d` and require its diff to contain only the approved nine implementation files.

Do not rewrite local `main`, force-push, squash the planning history into `3d`, or change any 3d implementation file in this prerequisite.

## Expected surface

| File | Expected change | Estimate |
| --- | --- | ---: |
| `docs/roadmap/durable-decisions/index.md` | Link the sprint and remainder exit contract | +21 / -7 |
| `docs/roadmap/live-test-truth/dispatch-sprint-execution.md` | Add the Commander execution package | +140 / -0 |
| `docs/roadmap/live-test-truth/index.md` | Add the reviewed sprint plan and sequence | +192 / -0 |
| `docs/roadmap/live-test-truth/staff-review.md` | Add the preflight and DoD review record; remove its trailing EOF blank line | +49 / -0 |
| `docs/runtime-live-ci-registry.md` | Add the desired live-test registry | +374 / -0 |

Expected total: exactly five files, 776 insertions, and seven deletions. No tolerance for additional paths or any other byte drift from the already-authored tree.

## Acceptance criteria

**AC-1 (VALUE) - The reviewed sprint plan and desired live-test registry are durable on `main` before 3d lands.**

Verified by: the candidate contains exactly the five expected files at the authored final contents except for the authorized EOF blank-line removal, including the member sequence `3d` → `15e` → `ys` and the deferred `rm` boundary.

**AC-2 (VALUE) - The planning package no longer contaminates 3d's implementation PR surface.**

Verified by: after this candidate lands and current main is integrated, `git diff --name-status origin/main...<3d-candidate>` contains only 3d's approved nine paths.

**AC-3 - Current-main work and substantive authored planning history are both preserved without destructive rewriting.**

Verified by: the candidate has the existing planning tip and current `origin/main` as merge ancestry, is zero commits behind current main, differs from the merged planning tree only by the authorized EOF blank-line removal, and no force/reset operation is used.

## Test plan

Verify the exact five-file `origin/main...candidate` name-status and +776/-7 numstat. Compare the five candidate blobs with `b0a2541f656e0b6d72d01d411160c1694910cbee`: four must match, and `staff-review.md` must differ only by the removed trailing EOF blank line. Run `git diff --check`, then the repository-required checks:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

No live runtime lane is required because this prerequisite changes documentation only. Relevant PR CI is the secret-free offline, docs, and install surface.

No spike is needed. `git merge-tree --write-tree origin/main b0a2541f656e0b6d72d01d411160c1694910cbee` already exercised the only topology mechanism without mutating a branch: it exited zero, produced the five declared paths, and preserved all five authored blobs. Its one concrete hygiene finding is independently reproducible: `git diff --check origin/main...b0a2541f656e0b6d72d01d411160c1694910cbee` rejects the trailing blank line that this plan requires the implementation to remove.

## Stage Report: ideation

- DONE: Audit whether a separate planning-package prerequisite is necessary to keep 3d within its approved nine-file surface; compare the simpler alternatives of dropping the files or carrying them in 3d and record why they are insufficient.
  `3d` declares nine implementation paths; dropping the package loses the sprint contract and desired registry, while carrying it adds five unrelated paths and misattributes 777 insertions and seven deletions to `3d`.
- DONE: Verify the existing final planning tip b0a2541f656e0b6d72d01d411160c1694910cbee and current origin/main merge conflict-free into exactly the five declared paths with +777/-7; record no-spike-needed grounds or a concrete finding.
  After `git fetch origin main`, `origin/main` was `1947aacb0d7c3481c18a846f3566645fd2cb89ee`; `git merge-tree --write-tree` exited zero with tree `635cabd99999f8f703526450eb3c28eae54fd05e`, and name-status/numstat proved five paths at +777/-7 with all five blobs matching the planning tip.
  `git diff --check` exited 2 on `docs/roadmap/live-test-truth/staff-review.md:50: new blank line at EOF`; the tightened plan authorizes only that removal and requires the post-trim green control at +776/-7.
- DONE: Review and, if needed, tighten Outcome, Problem, Proposed approach, Expected surface, AC-1 through AC-3, and Test plan so every criterion has falsifiable evidence and no 3d implementation path enters scope.
  The body now distinguishes substantive history from incidental whitespace, fixes the expected count, requires four exact blob matches plus one single-line exception, and keeps all nine `3d` paths outside nz6.
- DONE: Append a complete ideation stage report with exact commands/results and a concise summary; commit only land-live-test-truth-planning-package.md in the split-root state checkout.
  This report records the topology proof, necessity alternatives, exact whitespace finding, and the path-scoped state-commit requirement.

### Summary

The separate prerequisite is necessary to keep `3d` truthful to its approved nine-file surface while retaining the sprint plan and desired registry. The merge mechanism is already proven conflict-free; ideation found and bounded one trailing-blank-line defect, so implementation must land five paths at +776/-7 with no other authored-byte drift.

## Stage Report: implementation

- DONE: Integrate current origin/main into the dedicated nz6 branch while preserving the existing planning tip ancestry; remove only the approved trailing EOF blank line from staff-review.md and commit the resulting candidate without force/reset.
  Merge commit `603594374d034cfcef0190a5e15b71094badf1cb` preserves planning tip `b0a2541f656e0b6d72d01d411160c1694910cbee` and `origin/main` `1947aacb0d7c3481c18a846f3566645fd2cb89ee`; candidate `7461beda42382391a80ca9b14d42fafebaebaf55` removes only the final blank line.
- DONE: Prove the immutable candidate is 0 behind current main with exactly the five approved paths at +776/-7, four blobs byte-identical to b0a2541f and staff-review.md differing only by this deletion; run git diff --check.
  A post-test fetch left `origin/main` at `1947aacb`; left/right count was `0 12`, shortstat was exactly five files with 776 insertions and seven deletions, four blob hashes matched, the trimmed fifth matched after deleting its last blank line, and `git diff --check origin/main...HEAD` exited 0; any added path, byte drift, or whitespace defect fails these assertions.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race; append a complete implementation report with exact SHA/results and commit only the nz6 entity path in split-root state.
  `gofmt -w ./cmd ./internal` left the worktree clean at `7461beda4`; `go test ./...` and `go test ./... -race` both exited 0, exercising all packages including CLI, gates, status, live-scenario fixtures, and skill integration, where a compile, behavior, fixture, or race regression would fail the respective suite.

### Summary

Current main and the substantive planning history are now joined without rewriting either side. The immutable candidate contains exactly the approved five-file planning package, with only the authorized EOF cleanup, and all repository-required checks pass; no live runtime lane was required for this documentation-only change.
