---
title: Reliably preload the two small first-officer cores with @ references
status: implementation
score: 0.95
source: "Fresh local Claude shallow-boot evidence on 2026-07-11: exact engage ran an unscoped find / for fo-merge-core.md after the bare deferred reference failed to resolve from the fixture cwd. The captain previously rejected promotion to separately callable skills and chose the established @-reference loading pattern for two small cores; the larger fo-dispatch-core remains an explicit tradeoff."
completed:
verdict:
worktree:
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
