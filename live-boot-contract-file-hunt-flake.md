---
id: f4qwzrd5etzrrc1e5emhpd6k
title: Live sonnet flake — boot-preamble contract-file broad-search hang (rejection-flow)
status: backlog
source: "PR#516 CI investigation 2026-07-17 (claude-live sonnet, rejection-flow, run 29419612048)."
started:
completed:
verdict:
score:
worktree:
issue:
---

The `claude-live (sonnet, CI-E2E)` lane intermittently reds on shared scenarios because the FO broad-searches the filesystem at boot instead of invoking a skill. Track and remediate this flake class.

## Problem

On PR#516 (run 29419612048), `TestLiveClaudeSharedScenarios/rejection-flow` failed on the sonnet lane only (opus passed). Evidence from the run's own `claude-shared-scenarios-detail.jsonl`:

- Failure: `claude_live_runner_test.go:130` — "made no stream progress within 1m0s (no-progress quiet budget) — a hung stage; killed the subprocess."
- Diagnostic: `claude_live_failure_diagnostic_impl_test.go:24` — "FO broad-searched the filesystem at boot."
- Root cause (transcript tail): during boot the FO ran `find / -iname "*feedback-rejection-flow*" 2>/dev/null` — a full-`/` scan to locate a contract/skill file — instead of invoking `Skill(skill="spacedock:feedback-rejection-flow")`. The scan emits no stream events, so the 60s no-progress watchdog killed the subprocess.

This is the **boot-preamble contract-file-hunt** broad-search class from PR #490, which added (a) a fast classifier and (b) a 3× `bootPreambleMaxAttempts` retry lever — explicitly a mitigation that "by design does not prevent" the behavior (PR #490 measured the class at ~33% on the sonnet lane). In this run the retry engaged (2 boot attempts, both hit the same `find /` hunt, ~114s total) and still exhausted to red.

It is a DISTINCT trigger from the FO-hunts-the-`spacedock`-binary case already solved by the "Install spacedock candidate" CI step (`runtime-live-e2e.yml`: builds the candidate, sets `SPACEDOCK_BIN` + `$GITHUB_PATH`). Installing the binary does not stop the FO from hunting a *skill/contract file*.

## Proposed approach

{Ideation fills this in. Candidate directions to weigh: strengthen the shared-core zero-discovery / no-broad-search wording so sonnet reliably invokes the skill rather than `find`-ing it; make the skill/contract path discoverable so no search is tempting; and/or verify the boot-preamble retry actually re-launches on the no-progress-kill path (2 attempts observed vs the 3 configured — confirm whether the kill path and the classify-and-retry path fully agree).}

## Out of scope

- PR#516's reviewer-reuse oracle change (test-only; unrelated — this flake is model boot behavior).
- The `pi-live` job stuck WAITING (separate CI/runner-capacity issue).
- The FO-hunts-the-binary trigger (already solved by the CI candidate-binary install).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - The sonnet live lane no longer reds on the boot-preamble contract-file broad-search class for the shared scenarios.**
Verified by: {ideation defines the measurable — e.g. a repeated live-lane pass rate, or an offline test proving the boot path invokes the skill rather than a filesystem search.}

## Test plan

{Ideation fills this in. Likely offline coverage in `internal/ensigncycle` (boot-preamble classify/retry) plus a live sonnet rejection-flow smoke to confirm the class no longer surfaces.}
