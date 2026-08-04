---
title: Land the live-test-truth planning package
status: backlog
source: "Captain-directed topology prerequisite before 3d, 2026-08-04"
score: 1.0
sprint: live-test-truth
group: truthful-results
sprint-readiness: ready
id: nz6816mcf0qsdp1bsh3wzp1w
---

Land the already-authored sprint plan separately so implementation members retain their approved product surfaces.

## Outcome

The live-test-truth sprint index, staff review, execution package, desired registry, and durable-decision links are available on `main` before `3d` opens its implementation PR.

## Problem

The five planning and registry paths were authored and committed before sprint execution, but they remained on unpushed local ancestry. The `3d` branch inherited them, making its current diff 14 files instead of its approved nine-file implementation surface.

Merging those files as part of `3d` would misattribute 777 insertions and seven deletions to an implementation member whose approved expected surface does not own them. Dropping them would lose the sprint contract and desired-state registry.

## Proposed approach

Use the existing final planning commit `b0a2541f656e0b6d72d01d411160c1694910cbee` as one parent and merge current `origin/main` as the other in a dedicated branch. The preflight merge-tree is conflict-free and preserves both the authored planning history and all current-main work.

Open and merge a separate documentation PR. Then merge the resulting current `origin/main` into `3d` and require its diff to contain only the approved nine implementation files.

Do not rewrite local `main`, force-push, squash the planning history into `3d`, or change any 3d implementation file in this prerequisite.

## Expected surface

| File | Expected change | Estimate |
| --- | --- | ---: |
| `docs/roadmap/durable-decisions/index.md` | Link the sprint and remainder exit contract | +21 / -7 |
| `docs/roadmap/live-test-truth/dispatch-sprint-execution.md` | Add the Commander execution package | +140 / -0 |
| `docs/roadmap/live-test-truth/index.md` | Add the reviewed sprint plan and sequence | +192 / -0 |
| `docs/roadmap/live-test-truth/staff-review.md` | Add the preflight and DoD review record | +50 / -0 |
| `docs/runtime-live-ci-registry.md` | Add the desired live-test registry | +374 / -0 |

Expected total: exactly five files, 777 insertions, and seven deletions. No tolerance for additional paths; normal line-count drift is not expected because the candidate reuses the already-authored tree.

## Acceptance criteria

**AC-1 (VALUE) - The reviewed sprint plan and desired live-test registry are durable on `main` before 3d lands.**

Verified by: the candidate contains exactly the five expected files at the authored final contents, including the member sequence `3d` → `15e` → `ys` and the deferred `rm` boundary.

**AC-2 (VALUE) - The planning package no longer contaminates 3d's implementation PR surface.**

Verified by: after this candidate lands and current main is integrated, `git diff --name-status origin/main...<3d-candidate>` contains only 3d's approved nine paths.

**AC-3 - Current-main work and authored planning history are both preserved without destructive rewriting.**

Verified by: the candidate has the existing planning tip and current `origin/main` as merge ancestry, is zero commits behind current main, and no force/reset operation is used.

## Test plan

Verify the exact five-file `origin/main...candidate` name-status and +777/-7 numstat. Run `git diff --check`, then the repository-required checks:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

No live runtime lane is required because this prerequisite changes documentation only. Relevant PR CI is the secret-free offline, docs, and install surface.
