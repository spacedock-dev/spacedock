---
title: Reliably preload the two small first-officer cores with @ references
status: implementation
score: 0.95
source: "Fresh local Claude shallow-boot evidence on 2026-07-11: exact engage ran an unscoped find / for fo-merge-core.md after the bare deferred reference failed to resolve from the fixture cwd. The captain previously rejected promotion to separately callable skills and chose the established @-reference loading pattern for two small cores; the larger fo-dispatch-core remains an explicit tradeoff."
completed:
verdict:
worktree: .worktrees/spacedock-ensign-fo-small-core-at-reference-imports
issue:
id: m1y5k6w8any3gachwxtxqjfk
---

The 0.25 recovery stack proves the detector and fixture fixes are insufficient by themselves: a real Claude first officer still searched the host filesystem for the known merge-core document. Preload the two small cores through the existing top-level `@references/...` import mechanism so they are reliably addressed without adding callable capabilities. Keep the 2,386-word dispatch core deferred for this release, retain its anti-hunt stop rule, and do not merge the rejected skill-promotion portion of PR #491.

## Acceptance criteria

**AC-1.** A fresh-cwd Claude first officer reaches the merge boundary without `find`, recursive grep, or path reconstruction for `fo-merge-core.md`; the shallow-boot durable greeting/engage assertions pass.

**AC-2.** `fo-merge-core.md` and `fo-smallest-sufficient-mechanism.md` are loaded through top-level `@references/...` imports, while `fo-dispatch-core.md` remains deferred and no new skill/capability directory is added.

**AC-3.** Contract-lint, full, race, live-tag compile, and the coupled live recovery scenarios remain green.

## Proposed approach

Add the two small `@references/...` imports to the first-officer skill entry point and update structural tests to pin the chosen split: reliable eager imports for the 472-word merge core and 731-word smallest-sufficient core; lazy deferred loading for the 2,386-word dispatch core. Use the existing saved Claude transcript as the red behavioral baseline and the rerun as the green proof.

## Stage Report: implementation

- DONE: Work in a new isolated worktree based directly on the `release/0.25.0-live-recovery` integration tip, not the rejected promoted-skill branch.
  Commit `6baeed70192676a5fa6ae70e1de5dd1af7e3ef49` has parent and merge-base `89eee07e2fb2fe7845d9a7759c8def69175794ee`; the code worktree is clean.
- DONE: Preload exactly the two small first-officer cores through the established top-level import mechanism, keep dispatch core deferred, and add no callable skill/capability.
  `skills/first-officer/SKILL.md` now imports `@references/fo-merge-core.md` and `@references/fo-smallest-sufficient-mechanism.md` after the existing shared-core import. The shared core no longer lists merge core as deferred or exposes either eager core through a bare `references/...` cue; the Claude adapter describes merge as preloaded and applied at terminal. `fo-dispatch-core.md` remains solely in the deferred dispatch path. No `skills/fo-merge-core/` or `skills/fo-smallest-sufficient-mechanism/` directory exists.
- DONE: Add/update structural contract tests using TDD to pin the eager/deferred split and real target closure.
  RED: `TestFirstOfficerEagerReferencesKeepDispatchCoreDeferred` first failed with only the shared-core import, then failed again while merge remained in `## Deferred load points`, and finally failed on the smallest-sufficient bare load cue. GREEN: the test now asserts the exact import sequence, reads both eager targets as non-empty files, rejects either bare shared-core cue, rejects promoted skill directories, and relies on the updated deferred/host-neutral closure tests to keep dispatch lazy and merge eager. `go test ./internal/contractlint -count=1` passed 91 tests.
- DONE: Preserve the saved behavioral RED baseline and make no unproven behavioral-green claim.
  `/tmp/spacedock-m3-claude-current2/pty-team-mode/shallow-boot/pty-engage-stream.jsonl` records unscoped `find /` calls for `fo-merge-core.md` at events 77 and 83 from the fresh fixture cwd. This implementation supplies and structurally verifies the preload; AC-1 behavioral green remains for a fresh live validation rerun.
- DONE: Run focused contractlint, scoped formatting, coupled shallow-boot tests, full tests, race tests, live-tag compile, and diff checks; commit the scoped result.
  `go test ./internal/contractlint -count=1` passed 91 tests; `go test ./internal/ensigncycle -run '^TestShallowBoot' -count=1` passed 26 tests; `go test ./...` and `go test ./... -race` each passed 2,141 tests across 17 packages; `go test -tags live -run '^$' ./...` exited 0 without executing live tests. Scoped `gofmt -d` and `git diff --check` emitted no findings. The six-file commit contains only the first-officer import/prose alignment and structural contract tests.

### Summary

The first-officer entry point now eagerly loads the merge and smallest-sufficient cores as real `@references` imports, while the substantially larger dispatch core remains deferred. Structural tests enforce the exact split, target readability, absence of promoted skills, and absence of cwd-relative reload cues. All offline, race, and live-tag compile gates pass; the saved filesystem hunt remains the behavioral RED baseline until validation performs a fresh live run.
