---
id: ca9gqaz4n97gsv2nd9p80xbf
title: The codex filing recognizer misses a create terminated by a newline
status: ideation
source: "Captain CL, 2026-08-18, from the live-lane inventory. Root-caused with a runnable repro: codexFilingCreateCount's terminator class (?:[ \\t';|&]|$) in shared_filing_test.go omits \\n, so a valid create followed by a verification line on the next line counts zero. Red in run 32105482382 codex-live as observed=[filing-command-not-observed]."
started: 2026-08-18T18:41:27Z
completed:
verdict:
score:
worktree:
issue:
---

The filing grader fails a compliant agent. A create followed by a newline counts as zero creates, so any run that verifies its own filing goes red.

## Problem

`codexFilingCreateCount` recognizes the blessed atomic create by regex. Its terminator class is `(?:[ \t';|&]|$)` — space, tab, quote, semicolon, pipe, ampersand, or end of string. It omits the newline.

In run 32105482382 the FO did exactly the right thing: ran `"${SPACEDOCK_BIN:-spacedock}" new wire-the-thing` piped from printf, `item.completed`, exit 0, receipt `created: … id=001` in the aggregated output. Then, in the same `bash -lc` item on a new line, it ran a read-only `status --read wire-the-thing --json` to verify.

The create counted zero. The journey graded `filing-command-not-observed`.

Reproduced both directions: the regex yields 0 matches on that command and 1 on the existing PR679 valid fixture in `shared_filing_negative_test.go`.

Two things make this worth fixing beyond the single red.

It is deterministic. Any codex run that appends a line after its create will red, forever, until the terminator class changes. Nothing about it is flaky.

It punishes the behavior we want. The FO verified its own filing before reporting — exactly the discipline the contract asks for — and the grader marked it as not having filed.

`own-codex-filing-variant-miss` does not cover this. That entity is the skipped-second-variant mode, and its out-of-scope line declares the filing grader honest, which this occurrence disproves for the newline shape.

## Proposed approach

{Ideation fills this in. The change is adding `\n` to the terminator class. Confirm no other recognizer in the same family shares the omission, and use the run-32105482382 command as the falsifying fixture.}

## Out of scope

The skipped-second-variant mode owned by `own-codex-filing-variant-miss`. Any change to what the FO is expected to run when filing.

## Expected surface and tolerance

Estimate net LOC change: +15 across 2 files. Report insertions and deletions separately. Do not declare a gross tolerance. Semantics changed: none in the product; the grader accepts a shape it previously rejected.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A create terminated by a newline is recognized.**
This is the measuring AC: the recognizer's count over the run-32105482382 command must be 1, where it is 0 today. Verified by a fixture carrying that exact command, create followed by a newline and a verification line. Fails if the terminator class regresses.

**AC-2 - No shape that should fail now passes.**
Verified by the existing negative fixtures still grading red — a run that never files, and a run whose create is malformed. Fails if widening the terminator turns the recognizer permissive, which would hide a real filing failure.
