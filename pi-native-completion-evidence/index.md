---
title: Recognize native Pi worker completion evidence
status: implementation
source: Captain binding resolution:binding-1789576723761874000; tip CI 35058669297
started: 2026-09-16T17:00:12Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-pi-native-completion-evidence
issue:
pr:
mod-block:
id: cpx5q07b7wqmd6xv6wc7khvd
---

Recognize Pi native completion notifications so valid worker lifecycles are graded correctly without accepting another worker or hiding downstream failures.

## Problem

Tip CI 35058669297 reports incomplete implementation/validation workers although retained parent and child evidence proves their identity, committed reports, completion, and notification before advancement. The existing grader recognizes status/wait responses only. Auto-continue separately fails gate preparation on a mistyped reference path; that recorded whole journey must remain red.

## Proposed approach

Extend the existing Pi lifecycle evidence path. Join the parent dispatch result run ID to the notification-referenced child session identity and original assignment. Credit completion at its position in the parent timeline. Keep existing report/gate assertions and historical completion forms. Reuse existing parsers and retained sessions. Establish the smallest concrete surface and a failing captured replay before candidate edits.

Captain approved this bounded follow-up, independent validation, and a PR above the current stack through binding resolution:binding-1789576723761874000 at 2026-09-16T16:38:43.761876Z. Approved artifact /private/tmp/pi-native-notification-fix-review.md, sha256:fb4c087fbde4cad12db9080553c5a85a84a7ada2f504da760151d8c03fbcc3d5. This seed records the approved brief; no earlier stage report is manufactured.

## Risk evidence

Existing read-only triage: ../fo-continuation/artifacts/validation/pi-tip-ci-triage/report.md and retained-events.json/source-manifest.json. Downloaded artifacts: /tmp/spacedock-tip-ci/pi. Exact parent toolCallId/result run ID joins child session_info run identity, assignment, successful commit and terminal event. The native notification precedes advancement. Child directory UUID differs from run ID, so generic labels or filename guessing are insufficient.

## Out of scope

Product runtime changes, process controllers, new CI lanes, new public commands, naming changes, weakening gate checks, or fixing a one-off path typo by adding a mechanism. No live model runs locally. Expensive verification stays at the final tip.

## Expected surface and tolerance

Use existing lifecycle helper, Pi evidence adapter, shared native-stream bridge, and existing test owners with a compact captured fixture. Before editing, report concrete paths and size plus the failing replay. The approved boundary is one narrow correlation path; return for design review if it requires a new framework or materially broader scope. Do not invent an arbitrary numerical cap absent from the approved brief.

## Acceptance criteria

**AC-1 — Native Pi completion credits only the dispatched worker at the correct boundary.**
Verified by: captured parent/child replay from CI credits default implementation and auto validation, while wrong run/assignment, missing or unsuccessful child, and out-of-order notification fail. Removing identity correlation must make a negative case pass incorrectly and be caught.

**AC-2 — Correct completion observation does not hide outcome failures or change other hosts.**
Verified by: the original auto-continue whole journey still fails because its gate was never prepared; historical Pi status/wait and existing Claude/Codex lifecycle controls continue to pass. Bypassing gate assertions fails the retained negative case.

**AC-3 — The repaired stack supplies honest final CI evidence.**
Verified by: deterministic checks locally and required tip PR host lanes after independent validation. Retain concrete failures and expected-failure policy; never call skipped/cancelled checks green. No local model or duplicate middle-layer live run.

## Test plan

Start with the exact captured failure replay through existing lifecycle test owners. Add narrowly targeted correlation/ordering negatives, preserving durable state assertions. Run go test ./..., go test ./... -race, gofmt -w ./cmd ./internal. Known local installed-manifest resolver failure requires exact-log confirmation; it never waives new failures. Independent validation must test attribution and incomplete-outcome boundaries. Full expensive host checks run at the stack tip only.

### Feedback Cycles


## Stage Report: implementation

- DONE: Credit the captured native Pi completions only after matching the dispatched worker, child evidence and parent event ordering; wrong or missing evidence fails.
  Commit `ba8e6eeb4e28e310b17897a0a8b12e151366b8f9`: captured default implementation and auto validation replay pass; 27 negatives reject run/agent/epoch/task/cwd, missing/error evidence, escaped/symlink locators and ordering/duplicate conflicts. Removing correlation admits wrong workers; ignoring conflicting epochs reproduced a red test before its authorized fix.
- DONE: Preserve downstream gate failure and historical Pi/Claude/Codex behavior through existing replay and negative controls.
  Full retained auto lifecycle now passes; captured absent-gate reconstruction still returns `entity has no gates record`. Existing Git check rejects an uncommitted report; bypassing either durable check breaks the replay. Historical Pi status/wait, Claude replay, Codex lifecycle and cross-host bypass controls pass (`artifacts/implementation/focused.log`).
- DONE: Deliver a minimal committed correction with required local checks and clear tip-CI deferral.
  Code committed locally on the registered branch; focused tests, live-tag compile-only check and formatting completed. Final normal/race ran sequentially; ensigncycle passed in 354.805s / 335.058s. Native host lanes remain explicitly deferred below.

### Summary

`go test ./...` and `go test ./... -race` did not pass. Both exit 1 only at the anticipated `TestCodexResolveManifestAgainstInstalledHost`, `codex_resolve_test.go:44`: `spacedock@spacedock not installed in codex, but resolver returned "/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json"`. All other packages pass; no race diagnostic. Exact final/initial logs are in `artifacts/implementation/`. The FO explicitly DECLINED the exact pre-existing installed-manifest resolver defect for this task, as anticipated in the assignment; it is not an unresolved correction finding. This disposition does not change either exit-1 result or waive new failures.

`go test -tags live ./internal/ensigncycle -run '^$'` passes; it selects no tests. Removing artifactDir propagation from either existing call site would leave that live journey unable to credit native notices.

Final tip PR/native host CI and independent validation remain deferred, not passed. FO owns these after implementation handback; no local native/model, network-auth, PR, push, CI or rebase was performed. The original auto gate-path typo is retained as an outcome failure, not fixed or declared green. AC-3 remains pending final tip CI.

One 150-line Pi helper correlates the parent dispatch result with the notification-referenced child’s exact assignment and run/epoch, and credits completion at the native parent position. Both existing callers pass optional retained-artifact roots; existing parsers and independent Git/report/gate checks remain authoritative, with no synthetic event, child-call concatenation or second commit parser.

Actual surface is 12 files, +424/-11, within the proposed helper/test estimate: seven Go harness/test files plus four projected JSONL fixtures and provenance (14,818 bytes). Exact paths/numstat are in `artifacts/implementation/surface.txt`; they are `internal/ensigncycle/{claude_runtime_helpers_test.go,pi_native_completion_test.go,pi_rejection_extractors_test.go,auto_continue_fixtures_test.go,claude_live_runner_test.go,shared_promoted_live_test.go,pi_auto_continue_double_dispatch_replay_test.go}` and `internal/ensigncycle/testdata/pi_native_completion/{default-parent.jsonl,default-child.jsonl,auto-parent.jsonl,auto-child.jsonl,provenance.json}`.

Both pre-edit reproductions, separate FO FIX authorizations (G1 and conflicting epoch), command details and final/initial logs are retained under `artifacts/implementation/`. The captured fixture records are field projections with source SHA-256/line provenance; the Git end state is explicitly reconstructed because no original Git repository was retained. `gofmt -w ./cmd ./internal` ran, final touched files were formatted again, and the unrelated baseline release-test spacing change was restored; code worktree is clean. State report synchronization is left to the FO under the dispatch’s explicit no-push direction.
