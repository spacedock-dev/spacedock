---
id: ca9gqaz4n97gsv2nd9p80xbf
title: The codex filing recognizer misses a create terminated by a newline
status: implementation
source: "Captain CL, 2026-08-18, from the live-lane inventory. Root-caused with a runnable repro: codexFilingCreateCount's terminator class (?:[ \\t';|&]|$) in shared_filing_test.go omits \\n, so a valid create followed by a verification line on the next line counts zero. Red in run 32105482382 codex-live as observed=[filing-command-not-observed]."
started: 2026-08-18T18:41:27Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-filing-recognizer-newline-terminator
issue:
gates:
    version: 1
    records:
        - id: gate:ca9gqaz4n97gsv2nd9p80xbf:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ca9gqaz4n97gsv2nd9p80xbf-ideation-1
              briefing:
                id: briefing:ca9gqaz4n97gsv2nd9p80xbf:ideation:attempt-1:revision-1
                digest: sha256:dc31438f394cc960e2f834884ec29dec12f4d71805437a7e6cf1aef0338c8eee
                request-digest: sha256:014e43209381e1604ed3925ab97a94aff26cd2b5cbd2fa2b7d6af2ddaa8e2620
                room-ref: ./filing-recognizer-newline-terminator/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ca9gqaz4n97gsv2nd9p80xbf:ideation:1
                briefing: briefing:ca9gqaz4n97gsv2nd9p80xbf:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T19:31:30.61822Z"
                decision: approve
                reason: 'Captain approved in chat: ''approve those 4, and have them be on a pr stack.'' Accepts the ideation direction — both terminator classes widened, exact-bytes fixture from run 32105482382, existing negatives held red.'
              application:
                target-stage: implementation
                state: consumed
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

Add `\n` to two terminator character classes in `internal/ensigncycle/shared_filing_test.go` — not one. The run-32105482382 command is gated before the seed's named class is ever consulted:

1. `commandFilesViaNew`'s direct-invocation terminator (line 53): `(?:[ \t']|$)` becomes `(?:[ \t'\n]|$)`. This is the actual first blocker: `codexFilingCreateCount` returns 0 at its `!commandFilesViaNew` guard, so the create-counting regex never runs. The same function serves `assertClaudeFilingViaNew` and `assertCodexFilingViaNew`, so the Claude-side recognizer carried the same latent newline miss for multi-line Bash commands.
2. `codexFilingCreateCount`'s create counter (line 183): `(?:[ \t';|&]|$)` becomes `(?:[ \t';|&\n]|$)`.

Fixture: add the exact public bytes of run 32105482382's create item (codex-exec.jsonl line 14, artifact 9313785785 — the atomic create, a literal `\n`, then the read-only `status --read` verification in the same `bash -lc` item) as a `codexRun32105482382PublicCommand` constant in `shared_filing_negative_test.go`, mirroring the PR679 constant, plus two rows in the seven-rung matrix: the path-localized exact bytes must pass, and the same bytes with the create verb corrupted (`\" new ` to `\" mew `) must stay red. One positive fixture pins both classes: with only line 53 fixed the count stays 0, and with only line 183 fixed the gate rejects — either regression fails the row with "0 atomic creates".

Spike (riskiest path, exercised — 0 today, 1 after): both filing test files copied into a scratch module and run in three configurations against the exact run bytes. Unpatched: `codexFilingCreateCount` = 0 and `assertCodexPublicFilingTransaction` fails with "0 atomic creates" — the red reproduced offline. Line-183-only fix: count still 0 — refutes the seed's single-class hypothesis and proves the gate shares the omission. Both classes fixed: count exactly 1, the full public transaction passes (receipt and entity checks included), the malformed-verb mutation stays red, and all 42 existing filing fixture subtests still pass, so the widening turns no existing negative green.

For the one new mechanism (the exact-bytes fixture constant plus matrix rows), the value AC it serves is AC-1. Simplest alternative considered: a synthetic minimal `new slug\nverify` command string. Insufficient because the live miss came from the real quoting seam (`| \""'${SPACEDOCK_BIN:-spacedock}" new …`), and PR679 set the family precedent that retained public bytes are what keep the recognizer honest against that seam.

Alternatives considered for the fix itself:

1. Fix only `codexFilingCreateCount`'s class (the seed's hypothesis): spike-refuted; the count stays 0 behind the gate.
2. `(?m)` multiline flag so `$` matches before `\n`: the flag would govern the whole concatenated pattern (including the `newInvocation` prefix in `commandFilesViaNew`), a wider semantic change than the two-character class addition that matches the existing `\t` idiom.
3. Newline-split segments before matching (as `capturedLauncherFilesViaNew` already does): restructures `commandFilesViaNew`'s launcher-to-verb pairing for no additional recognition value.

Sibling survey (the rest of the family, checked either way):

1. Shares the omission: line 53 (`commandFilesViaNew` direct terminator) and line 183 (`codexFilingCreateCount` create terminator) — both fixed here. Nothing else does.
2. Newline-safe by construction: `newInvocation` (line 30) uses `[^\n]*?` as a deliberate same-line launcher-to-verb pairing constraint, not a terminator; `launcherCapture` (line 37) bounds with `(?:[;\s]|$)` and `\s` includes `\n`; `capturedLauncherFilesViaNew` (lines 85-86) terminates with `(?:[ \t]|$)` only after splitting segments on `\r?\n|;|&&|\|\||\|`, so a newline is a segment boundary; the receipt regex (line 213) anchors with `(?m)^…$`.
3. Adjacent recognizer family: every boundary class in `shared_round_recording_test.go` is `\s`-based (`(?:\s|$)`, `(?:\s|[;&|]|$)`) — newline-safe. No other filing recognizer exists in the package.

No doc changes: the fix is test-side grader behavior; no user-visible command surface, output, or docs-site behavior changes.

## Test plan

Offline fixture tests only — the assertions are pure functions over transcript bytes; no CLI or live workflow runs needed. Cost: one `go test ./internal/ensigncycle` cycle.

1. New matrix row, exact run-32105482382 public bytes (path-localized like the PR679 row): transaction passes. Fails if either terminator class regresses.
2. New matrix row, same bytes with the create verb corrupted: stays red. Fails if the widening ever loosens verb matching in the newline shape.
3. Existing suite unchanged and green: `TestAssertClaudeFilingViaNew`, `TestAssertCodexFilingViaNew`, `TestCodexPR679ExactPublicCommandTransaction`, `TestCodexPublicFilingTransactionSevenRungMatrix` (spike-verified: all 42 subtests pass on the patched recognizers).

## Out of scope

The skipped-second-variant mode owned by `own-codex-filing-variant-miss`. Any change to what the FO is expected to run when filing.

## Expected surface and tolerance

Estimate net LOC change: +12 across 2 files (`internal/ensigncycle/shared_filing_test.go`: two lines modified; `internal/ensigncycle/shared_filing_negative_test.go`: fixture constant, comment, localized variable, two matrix rows). Insertions ≈ 14, deletions ≈ 2, reported separately; no gross tolerance declared. Tolerance: ±10 net, same 2 files. Semantics changed: none in the product; the grader accepts a newline-terminated create it previously rejected — a shape the filing contract already blesses.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A create terminated by a newline is recognized.**
This is the measuring AC: the recognizer's count over the run-32105482382 command must be 1, where it is 0 today (spike-measured both directions on the exact public bytes). Verified by a fixture carrying that exact command — create, a literal newline, then the verification line — passing the full public filing transaction. The single fixture fails if either terminator class regresses: a line-53 revert fails the gate and a line-183 revert zeroes the count, both surfacing as "0 atomic creates".

**AC-2 - No shape that should fail now passes.**
Verified by the existing negative fixtures still grading red — a run that never files, a malformed launcher capture, a wrong slug, a manual create, duplicate creates, and the `--next-id` manual flow — plus a new newline-shaped control: the run-32105482382 bytes with a corrupted create verb must stay red. Fails if widening the terminator turns the recognizer permissive, which would hide a real filing failure.

## Stage Report: ideation

- DONE: Fix the terminator class and prove it with the real run-32105482382 command as the fixture — 0 matches today, 1 after.
  Proven by spike on the exact 723-byte public event (codex-exec.jsonl line 14, artifact 9313785785): unpatched count 0 with "0 atomic creates"; both classes widened, count exactly 1 and the full transaction passes. The before/after diff is recorded in the body; implementation applies it (ideation runs without a worktree, per the stage contract).
- DONE: Check whether any sibling recognizer in the same family shares the newline omission, and say so either way.
  Yes, one does: `commandFilesViaNew`'s direct terminator (shared_filing_test.go:53) shares it AND gates the seed's named class — a 183-only fix still counts 0 (spike-proven). Everything else is newline-safe: `launcherCapture` uses `\s`, `capturedLauncherFilesViaNew` splits on `\r?\n` first, the receipt regex is `(?m)`-anchored, and `shared_round_recording_test.go`'s classes are all `\s`-based.
- DONE: Keep the negative fixtures red: a run that never files and a malformed create must both still fail, so widening the terminator does not make the grader permissive.
  All 42 existing filing fixture subtests pass on the patched copy (their red expectations hold: never-files, malformed captures, wrong slug, manual create, duplicates, --next-id), and a new control — the run bytes with a corrupted create verb — stays red.

### Summary

Root cause sharpened: the seed named one omitting terminator class, but the spike proved two — `commandFilesViaNew` (line 53) gates `codexFilingCreateCount` (line 183), so both must gain `\n` for the run to count 1. The proposed fix is a two-character recognizer diff plus an exact-bytes fixture (constant + two seven-rung-matrix rows) that pins both classes with one positive row and one malformed-verb negative. Spike exercised both directions on the exact run bytes with all existing negatives staying red; net estimate +12 across 2 files, no product semantics change.
