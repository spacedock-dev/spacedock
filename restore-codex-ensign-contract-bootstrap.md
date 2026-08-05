---
title: Restore the shared ensign contract at the Codex fresh-dispatch boundary
status: backlog
source: "Runtime Live E2E run 31045591048, Codex job 92440554439: the first full-ensign-cycle worker received only a Read pointer, never loaded the shared ensign contract, and wrote two generic checkbox Stage Reports that the FO accepted before archive."
started:
completed:
verdict:
score: 0.95
worktree:
issue:
sprint:
id: nvz2ym82ydfn07jp04yfxg9r
---

## Outcome

Every fresh Codex stage worker receives the shared ensign operating contract before it executes the dispatch assignment. The worker therefore follows the exact Stage Report protocol and the First Officer never treats a generic `## Stage Report` checkbox section as contract-complete.

## Problem

The Codex helper emits a pointer-only outer prompt that says only `Read /tmp/spacedock-dispatch/...`. Its generated file claims to contain the shared ensign discipline, but it contains only the host first-action note, stage/entity/checklist context, fetches, and completion signal. It does not contain `skills/ensign/references/ensign-shared-core.md` or another executable bootstrap to that contract.

Exact-head run `31045591048` reached this supported path on the first `full-ensign-cycle` journey. The isolated Codex child appended `## Stage Report` with markdown checkbox items. The First Officer sent one report-repair turn, accepted a second generic checkbox report, advanced through the workflow, and archived the entity. The live oracle correctly failed because no anchored `## Stage Report: {stage}` existed.

The archived Codex runtime adapter task (`r09jrf0k6qjv6c1sddhe1sh6`) deliberately removed `Skill(skill="spacedock:ensign")` in June 2026 because that Codex skill surface was then unavailable. That historical choice is not proof that the current host cannot load the skill. The current dispatch file's claim that it already carries the shared discipline is false.

## Boundary

This task owns Codex fresh-worker contract bootstrap across the dispatch builder, Codex runtime binding, ensign entry point, and their behavioral tests. It does not change ys's live oracle, weaken Stage Report validation, restore parent-turn inheritance, reconstruct assignment payload in the First Officer, or absorb `self-contained-ensign-dispatch` (`kd7877nnbd19d528xnpwwaj4`), whose boundary explicitly excludes prompt/bootstrap changes. Preserve `fork_turns: "none"`, helper-prompt forwarding, pointer-only assignment transport, and Claude/Pi behavior.

## Acceptance criteria

**AC-1 (VALUE) — A fresh Codex ensign follows the shared Stage Report protocol.** A real fresh Codex stage worker emits an anchored `## Stage Report: {stage}` with DONE/SKIPPED/FAILED accounting, `### Summary`, and no checkbox bullets.

Verified by: the Codex `full-ensign-cycle` common live journey passes from a clean exact candidate; a negative fixture with the generic no-colon checkbox shape still fails.

**AC-2 — The bootstrap is executable, not a false prose claim.** The emitted Codex prompt/file has one current-host mechanism that actually loads or carries `ensign-shared-core.md` before assignment execution. A first implementation spike must exercise the current fresh-child skill surface before choosing skill invocation versus self-contained contract bytes.

Verified by: a fixture-backed dispatch-to-child probe fails when the bootstrap edge is removed even if the generated file retains the sentence that it "contains the shared ensign discipline."

**AC-3 — Fresh-context and prompt ownership remain intact.** Codex spawn still uses `fork_turns: "none"`; the First Officer forwards the helper-emitted prompt byte-for-byte; stage/entity/checklist payload remains in the dispatch artifact rather than being reconstructed by the runtime adapter.

Verified by: the existing exact spawn-map and pointer-only relational tests plus a removal mutation for the new fixed bootstrap.

**AC-4 — Other hosts do not regress.** Claude and Pi retain their supported bootstrap and completion-signal behavior.

Verified by: focused cross-host dispatch tests, full/race gates, and the live lanes required by the final changed-path set.

## Test plan

Start with the current-host spike required by AC-2. Add the smallest failing offline behavioral fixture at the dispatch/bootstrap boundary, including false-claim and generic-report negatives. Preserve the existing Codex isolation and prompt-identity controls. Then run focused dispatch/contract/ensigncycle tests, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and the changed-path-required live lanes.
