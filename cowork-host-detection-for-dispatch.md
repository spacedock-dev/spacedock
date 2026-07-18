---
title: Binary-side host auto-detection for dispatch helpers on Claude Cowork
status: backlog
source: gate follow-up from cowork-bootstrap-shared-launcher-gate ideation (live Cowork dispatch 2026-07-18)
started:
completed:
verdict:
score:
worktree:
issue:
id: wdmqf8kwd26dvwfv6v7e9g11
---

`spacedock dispatch build` derives host from CLAUDECODE / CODEX_THREAD_ID env markers; the Claude Cowork VM sets neither, so a Cowork FO must pass `--host claude` manually (observed live 2026-07-18). Give the binary a positive Cowork host source — e.g. a documented env/marker the FO sets from its runtime detection, or a first-class `cowork` host value — so helper calls need no manual escape hatch. Scope: helper host resolution only; skill-side detection is owned by cowork-bootstrap-shared-launcher-gate.

## Acceptance criteria

**AC-1 - A dispatch build invoked from a Cowork-marked environment resolves its host without --host.**
Verified by: unit test over the host-resolution function with the Cowork marker set; the current missing-host error remains for genuinely unmarked environments.
