---
id: ca9gqaz4n97gsv2nd9p80xbf
title: The codex filing recognizer misses a create terminated by a newline
status: done
source: "Captain CL, 2026-08-18, from the live-lane inventory. Root-caused with a runnable repro: codexFilingCreateCount's terminator class (?:[ \\t';|&]|$) in shared_filing_test.go omits \\n, so a valid create followed by a verification line on the next line counts zero. Red in run 32105482382 codex-live as observed=[filing-command-not-observed]."
started: 2026-08-18T18:41:27Z
completed: 2026-08-18T23:09:20Z
verdict: PASSED
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
        - id: gate:ca9gqaz4n97gsv2nd9p80xbf:validation
          stage: validation
          attempts:
            - id: gate-attempt:ca9gqaz4n97gsv2nd9p80xbf-validation-1
              briefing:
                id: briefing:ca9gqaz4n97gsv2nd9p80xbf:validation:attempt-1:revision-1
                digest: sha256:1d2ab3adb11ead300603be9a068067168e296e0751563b6ff3afbdbac000f254
                request-digest: sha256:9f392cf299238332b7d5b2b6f5a891be7d95ef14474f5b791ca2fcf6d79d7edb
                room-ref: ./filing-recognizer-newline-terminator/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ca9gqaz4n97gsv2nd9p80xbf:validation:1
                briefing: briefing:ca9gqaz4n97gsv2nd9p80xbf:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T21:54:16.967412Z"
                decision: approve
                reason: 'Captain approved in chat: ''approve all PR and trigger ci on tip.'' Validation PASSED with one deferred risk, no material finding: AC-1 disproof reproduced against a self-fetched artifact, 27-shape adversarial matrix found no new permissiveness class, surface exactly on estimate at net +12/2 files.'
              application:
                target-stage: done
                state: consumed
mod-block:
pr: pr-merge:731
archived: 2026-08-18T23:09:20Z
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

## Stage Report: implementation

- DONE: Both terminator classes accept a newline — commandFilesViaNew (shared_filing_test.go:53) and codexFilingCreateCount (:183) — and the exact run-32105482382 bytes count 1 where they count 0 today.
  ffaee5d6e widens both classes to `(?:[ \t'\n]|$)` and `(?:[ \t';|&\n]|$)`. New matrix subtest `rung 2 exact run-32105482382 public bytes (create + newline + verify)` in TestCodexPublicFilingTransactionSevenRungMatrix carries the exact 722-byte public event (codex-exec.jsonl line 14, artifact 9313785785, fetched live via `gh api .../artifacts/9313785785/zip`) and passes; reverting line 53 alone or line 183 alone (checked independently, each reverted and reverted back) makes that same subtest fail with "0 atomic creates" — proving the fixture needs both classes, not just the seed's named one.
- DONE: All 42 existing filing negatives stay red, plus a new malformed-verb control row.
  `go test ./internal/ensigncycle/... -race` passes in full (234.4s) — every pre-existing subtest in TestAssertClaudeFilingViaNew, TestAssertCodexFilingViaNew, TestCodexPR679ExactPublicCommandTransaction, and TestCodexPublicFilingTransactionSevenRungMatrix (including the never-files, malformed-capture, wrong-slug, manual-create, duplicate, and --next-id red controls) still fails as before. Added negative: `rung 2 run-32105482382 malformed create verb` (the same run bytes with `" new ` corrupted to `" mew `) still returns "0 atomic creates" — verified failing (asserted via t.Fatalf on wantErr mismatch) rather than merely asserted.
- DONE: Report actual net surface against the approved estimate of +12 across 2 files.
  Actual: +14/-2 across `internal/ensigncycle/shared_filing_test.go` (+2/-2, the two terminator-class edits) and `internal/ensigncycle/shared_filing_negative_test.go` (+12/-0, the fixture constant plus two matrix rows) — net +12, exactly the estimate, 0% deviation.

### Summary

Widened both terminator classes that gated the newline miss (commandFilesViaNew's direct-invocation terminator and codexFilingCreateCount's create counter), pinned by a fixture carrying the exact public bytes of the failing run (32105482382, artifact 9313785785, downloaded live via `gh api`) rather than a synthetic approximation. Confirmed by disproof: reverting either class alone still fails the new positive row with "0 atomic creates," so the single fixture load-bears both fixes. Committed to `ffaee5d6e` on `spacedock-ensign/filing-recognizer-newline-terminator`; `go test ./internal/ensigncycle/... -race` is green and the broader `go test ./...` has one pre-existing, unrelated failure in `internal/cli` (a local-machine codex-plugin-install mismatch, reproduced identically with this change stashed out).

## Stage Report: validation

- DONE: Reproduce AC-1's cited evidence independently: the exact run-32105482382 bytes count 1 atomic create, and reverting line 53 alone or line 183 alone each returns it to 0 — confirm the disproof rather than trusting the report.
  Re-fetched artifact 9313785785 myself (`gh api repos/spacedock-dev/spacedock/actions/artifacts/9313785785/zip`, not expired); `filing/codex-exec.jsonl` line 14 is 722 bytes and compares byte-equal to `codexRun32105482382PublicCommand` under full string equality (not a substring check), with 0x0A immediately after the slug. In a throwaway checkout at `ffaee5d6e`: the new row passes; reverting only line 53, and separately only line 183, each fails that same row with "has 0 atomic creates". Disproof reproduced, not taken on trust.
- DONE: Semantic adversarial pass on the widened terminator: prove that accepting a newline does not make the recognizer permissive on adjacent shapes — CR, CRLF, whitespace-only, an embedded newline inside a quoted argument, and a create verb split across lines must NOT count as an atomic create.
  27-shape matrix, every shape graded patched vs unpatched `main`. Red both ways: CR, CRLF, bare CR, vertical tab, form feed, NBSP, U+2028, U+0085 (0x0A is the only character admitted); whitespace-only command; verb/slug split across lines, verb split mid-token, launcher line then verb line, backslash continuation (`newInvocation`'s `[^\n]*?` still bars line crossing); slug prefix collision `wire-the-thing-2`. Embedded-newline-in-a-quoted-argument shapes (heredoc body, quoted multiline echo, unexecuted var assignment, shell comment) DO newly count — but a control run refutes it as a new permissiveness class: on unpatched `main` the same seven never-executed mentions terminated by space/tab/quote/EOS already grade GREEN, so the leak is the pre-existing "bounded simple-command recognizer, not a shell parser" boundary and `\n` only joins an already-leaky class. The widening also CLOSES a hole: two creates on two lines counted 1 (false GREEN) on `main`, now 2 (correctly red).
- DONE: Run go test ./internal/ensigncycle/... -race and confirm all pre-existing negatives stay red; separately confirm the reported internal/cli failure is pre-existing and unrelated by reproducing it with this branch's change stashed out.
  `-race -count=1` — the plain run returns `(cached)` and proves nothing — is ok at 292.9s. 44 filing subtests pass: the 42 pre-existing plus the 2 new rows, including 22 red-expectation matrix rows; a permissive recognizer fails them because the harness asserts `(err != nil) != wantErr`, so they cannot rubber-stamp. `internal/cli`'s `TestCodexResolveManifestAgainstInstalledHost` fails identically with `main`'s copies of both changed files checked out — stale local `~/.codex/plugins/cache/spacedock-local/…/0.27.0-pre7` state, no dependency on `internal/ensigncycle`.

### Deferred risk — newline widening can double-count one real create

The mirror of the fixed bug. `assertCodexPublicFilingTransaction` requires `creates == 1`, so one genuine atomic create plus any other `new <slug>` occurrence terminated by a newline now counts 2 and grades red where `main` graded green. Measured: `spacedock new wire-the-thing\necho ran spacedock new wire-the-thing` counts 2/red patched, 1/GREEN unpatched.

- Released user and normal workflow: the codex-live `filing` journey grading a real FO run.
- Observable harm: an FO that files atomically and also names the command on another line of the same `bash -lc` item grades `filing-command-not-observed` — the same false red this entity removes.
- Affected value AC or non-negotiable boundary: `none:` — AC-1 holds (the run bytes count 1) and AC-2 holds, since this is a false RED, not a permissive pass. No value AC promises "no shape that should pass now fails" beyond AC-1's single fixture.
- Trigger evidence: scanned all 191 retained `codex-exec.jsonl` files; 8 command items name the slug; exactly 1 changes count under the fix (0 to 1 — the intended repair) and 0 count above 1. The retained Claude PR399 stream shows 0 flips, and the Claude lane's `assertFilingCommands` has no cardinality check, so line 53 is monotone-permissive there. Never observed.

Promote to material if any live codex or Claude filing transcript shows one command item carrying a real create plus a second newline-terminated `new <slug>` occurrence, or if the FO contract starts asking the FO to echo the command it ran.

### Summary

Recommendation: **PASSED**, with one deferred risk and no material finding. AC-1's disproof reproduced independently against a self-fetched artifact whose bytes are byte-identical to the fixture; AC-2 holds under a 27-shape adversarial matrix where the only newly-accepted non-create shapes are proven pre-existing by a same-run control on unpatched `main`. The specific risk I was sent to attack — a newline making the grader permissive — did not land: no adjacent shape counts as an atomic create that did not already count with a space, tab, quote, or end-of-string terminator, and the change additionally repairs a duplicate-create undercount. The one real defect found runs the other direction (double-counting a create alongside a mention), is unobserved across the full retained corpus, and fails no value AC.
