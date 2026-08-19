---
title: Early rejection-round record grades as never-recorded
status: ideation
source: "Run 31996696789 claude lane, rejection-flow: round recorded at stream line 287 with entries=2, rework spawned line 309; captain approved filing owners for tolerated residual modes at the 0.27 composite-green ruling, 2026-08-17"
id: zf7rymtke3b6xp7r0337hjj4
gates:
    version: 1
    records:
        - id: gate:zf7rymtke3b6xp7r0337hjj4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zf7rymtke3b6xp7r0337hjj4-backlog-1
              briefing:
                id: briefing:zf7rymtke3b6xp7r0337hjj4:backlog:attempt-1:revision-1
                digest: sha256:67e727f265794634086d12c6cd19fc36f1ced497cb6cf325905a606cf4aa0ebc
                request-digest: sha256:a95473793cb5a3ac73d64382a04502e43a92229e3c2112dc29b626c91951f847
                room-ref: ./own-claude-early-rejection-round-record/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zf7rymtke3b6xp7r0337hjj4:backlog:1
                briefing: briefing:zf7rymtke3b6xp7r0337hjj4:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-19T18:16:02.74664Z"
                decision: approve
                reason: 'Captain approved in chat: ''if the fable opinion lands reasonable, dispatch those two on to the stack of 736 so we can run tip CI lane to verify.'' The fable review landed and confirmed the diagnosis while finding two real remedy gaps, both now recorded.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:zf7rymtke3b6xp7r0337hjj4:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:zf7rymtke3b6xp7r0337hjj4-ideation-1
              briefing:
                id: briefing:zf7rymtke3b6xp7r0337hjj4:ideation:attempt-1:revision-1
                digest: sha256:15c8b8b36c905b9847acda9746a90e5f6f0daed1bcb2ef1d8422c2007251c533
                request-digest: sha256:b65bdb1451c568bcb9ce669a3b6c41dd19dc40174eeef327e263651f24e782b0
                room-ref: ./own-claude-early-rejection-round-record/review/ideation/briefing-1
started: 2026-08-19T18:17:24Z
---

## Problem

Two coupled defects around one residual mode. The mode: the live claude FO invoked `gate record --round validation/1` immediately after the rejection — before routing the correction — so the immutable round room durably holds the reviewer's 2 entries instead of the complete 4 the contract's record-after-correction order produces. The red is correct in substance. The label is not: `claudeRecordedRejectionRound` accepts only a success line pinned to `entries=4` (`rejectionRoundSuccess`, shared_round_recording_test.go), so a successful early record is reported as "resolved launcher never invoked `gate record --round validation/1`" — the lying-label class this release spent itself killing. A second label defect in the same oracle: the workflow-README immutability boundary condition reports under `rejection-round-missing` ("README.md changed from its exact expected bytes", observed run 31991864922 attempt 2), which reads as a recording failure when the round WAS recorded.

## Recurrence confirmed (2026-08-19, FO)

The mode recurred in PR #736's live lane, run 32270990171, claude-live. Same shape, and the retained stream settles it beyond the earlier inference:

    tool_use:    ${SPACEDOCK_BIN:-spacedock} gate record rejection-task --round validation/1 \
                   --briefing .../briefing.json --log .../briefing.review.jsonl --workflow-dir ...
    tool_result: is_error = false
                 round=round:rejection-task:validation:1 stage=validation cycle=1
                 briefing=briefing:rejection-task:validation:round-1 entries=2
                 entry=annotation:rejection-task:missing-marker  type=Annotation
                 entry=resolution:rejection-task:reviewer        type=Resolution decision=revise

`entries=2`, exactly as this entity predicted: the reviewer's two entries, recorded before the correction was routed. `rejectionRoundSuccess` pins `entries=4`, so the oracle graded it "resolved launcher never invoked" — a successful, exit-0, correctly-formed invocation reported as absent.

One diagnosis to strike from the record, so it does not mislead the next reader: the FO initially attributed this red to `invokesRejectionRoundRecorder` failing to parse the contract-mandated `${SPACEDOCK_BIN:-spacedock}` form. That is WRONG. `directRoundLauncher` (`shared_round_recording_test.go:17`) accepts the default-expansion form explicitly via `(?::-[^}]*)?`, and it is checked first. The launcher parses fine; the entry count is the whole cause. The FO read the fallback branch's inner regex and never tested the branch that returns before it.

Evidence preserved at `/tmp/9g-evidence` (both hosts' full artifacts, 30 MB) independent of CI artifact expiry.

Sibling reds in the same journey, same release: `rejection-topology-count-bar` (`12z`, filed today) covers the codex-side `rejection-worker-topology` exact-count bar. Both are the same family as the merged `filing-recognizer-newline-terminator` (#731) and `decide-dispatch-build-count-bar` (#732) — a pattern written against one observed shape, then treated as the definition of correct conduct.

## Proposed approach

Split honesty from strictness in the oracle. Match the recorder's success line generically (`entries=(\d+)`) and grade the count as its own condition with an honest code (`rejection-round-incomplete`, naming got-vs-want and that the record preceded the correction); keep `rejection-round-missing` for genuinely absent invocations. Move the README-immutability condition to its own code (`rejection-workflow-doc-mutated`). Both proven by falsifying edits: replay the run-31996696789 stream bytes through the recognizer (early record must grade incomplete, not missing) and a README-mutated fixture (must grade doc-mutated, not round-missing). Skill-side, the record-after-correction order is already explicit in step 6; the owner tracks the mode's recurrence rate with the metrics instrument rather than adding prose the model already had.

## Adversarial review of this remedy (2026-08-19, fable ensign, read-only)

Diagnosis CONFIRMED against the live stream: one `--round` call at line 286, result at 287 carrying `entries=2`, `is_error=false`; the rework spawn is at index 307, so the record demonstrably preceded the correction. Every link holds.

It also settled the question the FO could not: is an honest label a way to make a real defect easier to ignore? No. The conduct IS wrong — the round room is immutable, so an early record permanently truncates the durable record to 2 of 4 entries. Under this remedy the red STAYS red; only its name changes. The opposite risk is already proven: the lying label sent the FO chasing a launcher-regex ghost. And it matches the codebase's own pattern (`assertRejectionCycleLine`'s heading-drift diagnostic, `TestRejectionUnpreparedGateReportsItsOwnCode`): the grade does not soften, the diagnostic must say what happened.

**But the remedy as written SHIPS HALF-DONE. Two gaps, both verified by the FO:**

1. **The `entries=4` pin lives in TWO places, not one.** The stream regex at `shared_round_recording_test.go:18`, and a durable summary check at `:400-405` — `len(summary.Entries) != 4`, message `"retained round summary = %#v"`, still under `rejection-round-missing`. Genericize only the stream regex and the same run reds at :404 under the same lying code with a WORSE message: a raw struct dump. The count grade must move to the DURABLE site under the honest code. And this remedy's proposed proof — replay the stream bytes through the recognizer — exercises only the stream path, so it would GREEN the half-fix. A durable fixture is required: record a 2-entry log, assert the code is `rejection-round-incomplete`.
2. **The byte-exact block at `:425-442` covers FOUR files**, not one: `briefing.json`, `briefing.review.jsonl`, `candidate.txt`, `README.md` — all reporting "changed from its exact expected bytes" under `rejection-round-missing`. This remedy splits only README. `candidate.txt` carries the identical lying label. The split criterion should be **"is this file part of the round record"** — the two room files keep the round code, everything else does not — rather than "is it README".

**Related conduct hotspot, named by the review and owned by neither entity.** Both hosts deviated at the SAME point in the flow: the gap between rejection and routed correction. Claude recorded the round early; codex inserted an extra review round. If recurrence continues, the fix is mechanical rather than prose — for example a precondition the FO can check before `gate record --round` (the log tail contains a correction entry). This entity's "metrics over prose" call stands for now, at two occurrences.

## Ideation design (2026-08-19, ideation ensign)

Validated the adversarial review's conclusions by execution and source, not by re-reading: the reviewer's harness at `/tmp/9g-evidence/matcher-exec` re-run this session (all four recognizer matchers true on the live command bytes; `rejectionRoundSuccess` false on the real `entries=2` result, true on the one-token `entries=4` mutant), and the live stream re-extracted first-hand (`rejection-flow/claude-stream.jsonl`: one `--round` tool_use at line 285, result `is_error=false entries=2` at 286, rework spawn at 307). The chain holds; both remedy gaps are real. One reachability correction to gap 2, found in source and recorded below. Line references cite `main` @ 61dd8e435; the grader files are byte-identical on PR #736's branch (verified by diff), so they hold at the implementation base.

### The contract, named first

The journey's round-record contract, independent of any observed run: (1) the FO records `validation/1` exactly once, AFTER routing the correction, so the immutable room durably holds the complete 4-entry log; (2) the room is canonical and immutable — one `round-1` room, `briefing.json` + `briefing.review.jsonl` byte-exact; (3) the flow mutates only what it owns — the fix lands in `rejection-task/index.md`; `candidate.txt` (the reviewed artifact) and `README.md` (workflow policy) are untouched end-to-end. The observed run satisfies (2) and (3) and violates (1)'s completeness half. The red is correct; only its name lies.

### Remedy, settled across BOTH entries=4 pins

**Pin 1, stream (`shared_round_recording_test.go:18`):** genericize `entries=4` to `entries=\d+` in `rejectionRoundSuccess`, keeping every other token of the line pinned (round id, stage, cycle, briefing). The stream oracle's question is "did a successful, well-formed `validation/1` record happen", not "was it complete" — completeness is durable state and is graded at the durable site. Host parity confirms this is the contract, not a convenience: `codexRecordedRejectionRound` (`:218`) already accepts on exit-code-0 with NO entries pin, so today the identical early-record conduct grades "never invoked" on Claude and reds at the durable site on Codex — one conduct, two labels. After the fix both hosts converge on the durable count.

**Pin 2, durable (`:400-405`):** split `len(summary.Entries) != 4` out of the composite into its own condition returning `&gradedErr{code: "rejection-round-incomplete", ...}`. Message: got-vs-want count, the retained entry IDs (self-evidencing — the reader sees exactly the reviewer pair and no worker entries), and, when got < want only, the causal clause: the round was recorded before the correction was routed, and the room is immutable, so the truncation is durable. When got > want (unobserved) the message stays got-vs-want + IDs without that clause. The identity half of the composite (ID/Stage/Cycle/Briefing, `:400-403`) stays as-is under `rejection-round-missing`; it has never fired live and is not this entity's mode. Inner-code pass-through is proven mechanism: `durableSemantic` preserves an assertion's own `gradedErr` (`claude_runtime_helpers_test.go:547-556`, pinned by `TestFilingSemanticFailureUsesTargetXFail`), precedented by `assertRejectionCycleLine`. The code name `rejection-round-incomplete` already exists in the roadmap vocabulary (`docs/roadmap/test-behavior-completeness/staff-review.md` M2: Pi holds exact evidence; Sonnet needed "matching exact evidence" — run 32270990171 is that evidence).

### Split criterion for the four-file block (`:425-442`): CONFIRMED, with one reachability correction

The review's criterion — "is this file part of the round record" — is confirmed. The two room files ARE the durable round record; their byte checks stay under the round oracle unchanged. `candidate.txt` and `README.md` are workflow files outside the record; their mutation is not a round-recording failure and moves to one new code.

The correction: the `candidate.txt` half of the split is UNREACHABLE at `:431` as the block stands. `gates.ValidateRoundFile` → `verifyRoundArtifacts` (`internal/gates/round.go:182-210`) digest-checks the WORKSPACE `candidate.txt` against the briefing's pinned artifact rev, so a mutated candidate reds at `:397-399` ("validate retained round: artifact … raw digest does not match Briefing revision") under `rejection-round-missing` before the byte block ever runs — and digest-match implies byte-match, so `:431` can never fire on mutation. The implementation must therefore MOVE the two non-record byte checks to before the `ValidateRoundFile` call inside `assertRejectionRecordedRound`. Falsifying edit: moving them back after it re-reds the candidate-mutation fixture under the lying code.

Two settled sub-decisions. (a) Code name: `rejection-workflow-file-mutated`, one code for both files, message naming the path — refutes the filed `rejection-workflow-doc-mutated` narrowly (`candidate.txt` is not a doc; the class is "a workflow file the flow does not own changed"); two codes for a class with one observed occurrence (README, run 31991864922 a2) and zero (candidate) would be pattern-multiplying. (b) Mechanism: the split stays INSIDE `assertRejectionRecordedRound` via inner `gradedErr`, not a new runner-wired assertion — a separate assertion would double-report candidate mutation (its honest code PLUS the round oracle's digest-check red under the lying code), so the in-function reorder is required either way and is the smaller mechanism. Contrast `assertRejectionGatePrepared`, split out because two conditions held contradictory positions on one state; no contradiction exists here. A third code for round-room tampering was considered and declined: zero occurrences, and the room byte checks' existing message names the path and the byte drift.

### Red stays red; only the diagnostic moves

The conduct is genuinely wrong: `recordRoundLockedWith` publishes the room immutably and rejects re-publication ("round identity already has a pointer without its immutable room"; the durable oracle reds a second room), so an early record permanently truncates the durable record to 2 of 4 entries — no post-hoc repair path exists. Mechanically the composite stays red: `gradeLive` returns status "fail" for any non-empty code set without xfail (default branch, `claude_runtime_helpers_test.go:610-611`), and after the fix the early-record end state yields exactly one code, `rejection-round-incomplete`. AC-1 pins status AND codes, so a softening (pass) and a kept lie (`rejection-round-missing`) both fail its test. Tip-run reading restated: claude-live RED under `rejection-round-incomplete` is success; RED under `rejection-round-missing` means the fix did not land; GREEN proves nothing about this fix (that run's FO happened to record after the correction — nondeterministic conduct).

### New mechanisms, value AC served, simplest alternative

1. Genericized stream regex → AC-1/AC-2. Alternative: keep `entries=4`, add a second early-record stream pattern — insufficient: duplicates count truth at a site that cannot see durable state, and leaves the Claude/Codex label asymmetry.
2. In-function `gradedErr` count split → AC-1. Alternative: new runner-wired assertion — insufficient: double-reports, larger diff, does not remove the lying label from the digest path.
3. Reordered non-record byte checks + `rejection-workflow-file-mutated` → AC-3. Alternative: split only README (as filed) — insufficient: candidate keeps the identical lying label behind a dead byte check.
4. Durable 2-entry fixture test → AC-1. Alternative: stream replay alone — insufficient: exercises only pin 1 and greens the half-fix (review gap 1).

## Acceptance criteria

- **AC-1 (measures the end-value).** For the preserved early-record durable end state (a `validation/1` room recorded from the 2-entry reviewer log — run 32270990171's shape), the lane grade is status `fail` with codes exactly `[rejection-round-incomplete]` and a detail naming got=2 want=4 and the retained entry IDs. Independent baseline that can move the wrong way: the current grader on the same state yields `[rejection-round-missing]` / "resolved launcher never invoked" (Claude) — wrong-way movements are status `pass` (grade softened) or the code unchanged (lie kept). Test: new durable fixture test (Test plan 2).
- **AC-2 (means, paired with AC-1).** `claudeRecordedRejectionRound` returns true on the live run's command/result byte shape with `entries=2`, and still false for missing/failed results and wrong round/entity/artifacts. Test: extractor case + existing controls (Test plan 1).
- **AC-3.** Mutating `README.md` or `candidate.txt` after a complete 4-entry record grades `rejection-workflow-file-mutated` (detail naming the path), not any `rejection-round-*` code; tampering a round-room file still grades under the round oracle's code. Test: three durable cases (Test plan 3).
- **AC-4.** The conforming end state (complete 4-entry record, files intact) still grades pass, and the no-invocation control still reds — no new red or green for conforming conduct. Test: existing assertions in `TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl` stay green unmodified.

## Test plan

1. **Stream path** (extends `TestRejectionFlowRoundInvocationExtractors`): result line with `entries=2` must satisfy the extractor; falsified by re-pinning `entries=4`. Independent instrument: re-run `/tmp/9g-evidence/matcher-exec` (extracts the regex literals mechanically from source) — the real-result match must flip false→true after the fix.
2. **Durable half-fix killer** (new focused test, mirroring `TestRejectionUnpreparedGateReportsItsOwnCode`'s shape): `writeRejectionWorkflow`, record `validation/1` from `rejectionReviewerLog` (the existing 2-entry fixture constant), then assert `gradeLive(false, durableSemantic("rejection-round-missing", assertRejectionRecordedRound(...)))` yields status `fail`, codes exactly `[rejection-round-incomplete]`, detail containing got/want and entry IDs. A stream-only fix leaves the durable site returning the plain error under `rejection-round-missing` → this test reds. Recorder acceptance of a 2-entry log needs no spike: proven live (exit 0, `entries=2`), and the room copy is the same bytes re-validated by the same parser (`round.go:105-113`).
3. **Split boundary, both sides**: after a complete 4-entry record — (a) mutate `README.md` → `rejection-workflow-file-mutated`; (b) mutate `candidate.txt` → same code (this case is the reachability falsifier: ordered after `ValidateRoundFile` it reds under the lying code); (c) tamper the room `briefing.review.jsonl` bytes → still the round oracle's code.
4. Full `go test ./...` and `-race` per repo policy; live tip run on the stack is the journey-level reading (red under the honest code), not this fix's proof.

Cost: unit/fixture tests only, no live spend; all in one existing test file.

## Expected surface

Estimate net LOC change: +80, across 1 file (`internal/ensigncycle/shared_round_recording_test.go`): ~95 insertions, ~15 deletions. Breakdown — product: 0; tests: +80; docs: 0. Tolerance: ±40 net lines; hard boundary: no product (non-test) code, no files outside `internal/ensigncycle`, no fixture-text or recorder changes. Semantics declared: CI lane diagnostic codes change for two failure modes (early record: `rejection-round-missing` → `rejection-round-incomplete`; non-record file mutation: `rejection-round-missing` → `rejection-workflow-file-mutated`); no command grammar, stored format, authority, or runtime behavior changes; composite red/green is unchanged for every conforming run. Doc diff: none required — the lane codes appear on no user-facing docs page (grep across `docs/` and `skills/`; only state records and the roadmap staff review). Skill prose: unchanged — record-after-correction is already explicit in `feedback-rejection-flow` step 6; recurrence is tracked by the metrics instrument (standing call, review concurred at two occurrences).

## Spike determination

No spike needed: the design rests on (1) recorder acceptance of a 2-entry reviewer-only log — proven live, run 32270990171 stream lines 285-286, re-verified first-hand this session; (2) `gradedErr` code pass-through under `durableSemantic` — proven by `TestFilingSemanticFailureUsesTargetXFail` and the shipped `assertRejectionCycleLine` precedent; (3) the recognizer chain on the real bytes — proven by executing `/tmp/9g-evidence/matcher-exec` this session; (4) `verifyRoundArtifacts` digest-checking the workspace candidate — source-verified (`round.go:182-210`), and the design is robust to either resolution (Test plan 3b pins the honest code regardless of which site would have red).

## Out of scope

- The round recorder binary (neutral; its output already carries the true entry count).
- Rejection-flow fixture text (the Cycle-line target pin shipped with the 0.27 stack).
- The codex-side `rejection-worker-topology` count bar (`12z`'s entity, same stack).
- The identity half of the durable composite (`:400-403`) and its `%#v` message — never fired live.
- A dedicated code for round-room tampering (considered, declined above: zero occurrences).

## Stage Report: ideation

- DONE: Settle the remedy across BOTH entries=4 pins, not one: the stream regex at shared_round_recording_test.go:18 and the durable summary check at :400-405. A fix to only the first re-reds the same run under the same lying code with a raw %#v message.
  Settled in "Remedy, settled across BOTH entries=4 pins": pin 1 genericizes to `entries=\d+` (host-parity: codex already has no stream count pin); pin 2 splits the count into `rejection-round-incomplete` via in-function gradedErr, identity composite untouched.
- DONE: Decide the split criterion for the four-file byte-exact block at :425-442 (briefing.json, briefing.review.jsonl, candidate.txt, README.md). The review proposes "is this file part of the round record" rather than "is it README". Confirm or refute with reasoning.
  Confirmed, with a reachability correction the review missed: verifyRoundArtifacts (round.go:182-210) digest-checks candidate.txt inside ValidateRoundFile, so the candidate byte check at :431 is dead code — the two non-record checks must move BEFORE the ValidateRoundFile call. One code `rejection-workflow-file-mutated` for both non-record files (refines the filed doc-only name); room files keep the round oracle's code.
- DONE: Design the proof so it cannot green a half-fix: a stream replay alone exercises only the stream path. Name the durable fixture that proves the count grade moved to the honest code.
  Test plan 2: record validation/1 from the existing 2-entry `rejectionReviewerLog` constant, assert gradeLive yields status fail + codes exactly [rejection-round-incomplete]; a stream-only fix leaves the durable site under rejection-round-missing and reds this test. Recorder acceptance of the 2-entry log proven live (run 32270990171, exit 0, entries=2).
- DONE: Confirm the red STAYS red. The conduct is genuinely wrong — the round room is immutable, so an early record permanently truncates the durable record. State how the composite stays red and what changes is only the diagnostic.
  "Red stays red" section: gradeLive's default branch fails the lane on any non-empty code set; the early-record end state yields exactly [rejection-round-incomplete], and AC-1 pins status AND codes so both softening and the kept lie fail its test. Tip reading restated: claude RED under the honest code = success; GREEN proves nothing.

### Summary

Validated the adversarial review by execution (re-ran /tmp/9g-evidence/matcher-exec; re-extracted the live stream first-hand: one --round call, is_error=false, entries=2, rework spawned after) and confirmed both remedy gaps, adding one new finding: the candidate.txt byte check is unreachable behind ValidateRoundFile's artifact-digest check, so the file split requires reordering, which the design and its falsifying fixture case now pin. Produced the gated design: contract-first remedy across both pins, honest codes `rejection-round-incomplete` and `rejection-workflow-file-mutated`, four ACs with a measuring AC-1 against the current grader as baseline, surface +80 net across 1 test file (±40), no product/doc changes, no spike needed (four proven mechanisms named).
