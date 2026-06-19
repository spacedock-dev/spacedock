---
id: 3e2bb3r432vbgp8nnw54wnjk
title: Wrong-root boot guard mis-parses a compound `cd X; …` boot (strip the trailing `;`)
status: ideation
source: 'v0.22.0 cut (2026-06-19): the opus claude-live lane on PR #390 went red on TestLiveClaudeSharedScenarios/rejection-flow — a FALSE POSITIVE. The cut shipped on the sonnet proof with an e2e-gate waiver; this is the recorded follow-up to make the guard robust so the live lane is clean-green (no waiver) next time.'
started: 2026-06-19T07:45:59Z
completed:
verdict:
score: 0.4
worktree:
issue:
sprint: 0205-layered-fo
group: cleanup
sprint-readiness: ready
---

Harden the wrong-root boot guard so a compound `cd X; …` boot command is parsed correctly, eliminating a false-positive that latest-opus trips.

## Problem

`bootPathArgs` (`internal/ensigncycle/wrong_root_detect_impl_test.go:106`) does `strings.Fields(command)` and takes the token after `cd` as the boot target. When the FO issues a compound boot in one bash call — `cd /tmp/…/001; ls -la; echo "===README==="; sed -n '1,60p' README.md` — `strings.Fields` yields `["cd", "/tmp/…/001;", …]`, so the cd target comes out as `/tmp/…/001;` with the trailing `;` glued on (no space before it). `isUnder("/tmp/…/001;", "/tmp/…/001")` is then false → `detectWrongRootBoot` fires "FO booted the wrong root" though the FO booted the CORRECT root and merely chained exploration after it.

The parser's own comment concedes the assumption: *"It splits on whitespace (the boot commands are simple `cd …` forms; quoting is not exercised by the real boot stream)."* That held for pinned 2.1.177 and for sonnet, but **latest-opus composes its boot as a chained `cd X; ls; cat README` command**, which violates the assumption. The guard's "a CI env leak likely lured the FO off its cwd" message is a misleading built-in hypothesis — the real cause is the unstripped `;`.

Evidence: v0.22.0 release run — claude-live (opus, CI-E2E-OPUS) red on this; claude-live (sonnet) + the merged-team-mode lane + all deterministic lanes green. The 0.22.0 e2e-gate was waived for exactly this false-positive.

## Proposed approach (seed — ideation to flesh out)

Make `bootPathArgs` robust to compound commands: split the command on `;`, `&&`, `|` boundaries (or at minimum strip a trailing `;`/`&&` from each extracted path token) before the `cd`/`--workflow-dir` extraction, so the cd target is the clean path. Keep it a pure parsing fix in the test harness — no product code changes.

## Out of scope

- The broad-search overstep guard (`tq0`) — a sibling boot-guard hardening, distinct mechanism.
- Any product/FO-contract change — this is a test-harness parser fix.

## Acceptance criteria (proof = behavior, never prose-grep)

**AC-1 — a compound `cd X; …` boot is NOT flagged wrong-root when X is the fixture root.**
Verified by: a `TestDetectWrongRootBoot` case feeding a chained boot command (`cd <fixtureRoot>; ls -la; sed README`) and asserting `detectWrongRootBoot` returns nil; plus the inverted case (`cd <outside>; …`) still correctly flags.

**AC-2 — the live opus claude-live lane goes green on the shared scenarios without a waiver.**
Verified by: a Runtime Live E2E run (opus) with TestLiveClaudeSharedScenarios green (no wrong-root false positive), so the next cut needs no e2e-gate waiver for this cause.

## Test plan

Offline table test over synthetic boot commands (the AC-1 cases — chained cd, plain cd, cd-outside, no cd) is the cheap deterministic proof; the live opus lane green (AC-2) is the integration confirmation, naturally exercised by the next live-e2e run. Cost: minutes.

## Related

- `tq0` zero-discover-broad-search-hardening — sibling boot-guard flakiness (broad-search detector).
- v0.22.0 / PR #390 — where this surfaced; the e2e-gate waiver recorded the cause.
