---
title: Retire the codex_foreground_wait_shape phrase checks — same prose-grep antipattern family as #394
status: done
source: "Validation of contractlint-antidrift-guard-hardening (#394, 2026-06-17) flagged codex_foreground_wait_shape_test.go (current shape from commit ab39d5d8, #378) carrying phrase checks at :118/:135 — enumerated hyphenation variants standing in for a lifecycle MEANING, the same prose-grep-as-proof family as the two doctrine guards retired in #394. PRE-EXISTING (ab39d5d8 is an ancestor of v0.20.3; the file is untouched in v0.20.3..HEAD), so it was OUT of the v0.20.4 cut surface and did NOT block the cut. Tracked here for a separate pass."
sprint:
id: b5qwmp2nc1yg6casrcb0dcc4
mod-block:
pr: pr-merge:500
verdict: passed
completed: 2026-07-12T14:38:25Z
archived: 2026-07-12T14:38:25Z
---

`codex_foreground_wait_shape_test.go` (#378) asserts a lifecycle meaning via enumerated literal hyphenation-variant phrases — the same prose-grep-as-proof antipattern #394 retired in the contractlint package (a literal-phrase match standing in for a meaning, which a paraphrase slips past). It is pre-existing and was outside the v0.20.4 cut surface, so it did not block the cut.

## What's needed (decide in ideation)
Assess the check against the same bar #394 used: is there a genuine token/structural property here (the token IS the defect), or is it a meaning-proxy that should be narrowed honestly / replaced by a behavior test / removed with the owed test recorded? Resolve consistent with #378's direction and the contractlint `doc_test.go:11` policy. Prove by running the behavior, not a prose-grep.

## Relates to
- `#394` (contractlint-antidrift-guard-hardening, archived PASSED) — the sibling retirement this follows.
