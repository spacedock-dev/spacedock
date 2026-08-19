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
              resolution:
                type: Resolution
                id: resolution:spacedock:zf7rymtke3b6xp7r0337hjj4:ideation:1
                briefing: briefing:zf7rymtke3b6xp7r0337hjj4:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-19T19:11:26.744371Z"
                decision: revise
                reason: 'Captain rejected in chat: ''if it stays red it has no value. combine it with something that turns it green.'' The label analysis is accepted and stays — including the dead-code finding on the candidate.txt byte check. But a remedy whose entire outcome is a truer sentence on a still-red lane ships nothing: the lane keeps blocking, the FO keeps recording early, and the workflow''s own bar rejects a task whose only output is a decision. Combine the honest codes with the mechanism that makes the early record impossible, so a conforming run goes GREEN.'
            - id: gate-attempt:zf7rymtke3b6xp7r0337hjj4-ideation-2
              briefing:
                id: briefing:zf7rymtke3b6xp7r0337hjj4:ideation:attempt-2:revision-1
                digest: sha256:e7cf066ba0ad99ab0eb95501de1dbe11e8c29e2c0d124c5dba64e260414c0a90
                request-digest: sha256:0611668a809bd7fa50da6c63ed0b908c6091ae4ffe767b927a8caa5c1b989877
                room-ref: ./own-claude-early-rejection-round-record/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:zf7rymtke3b6xp7r0337hjj4:ideation:2
                briefing: briefing:zf7rymtke3b6xp7r0337hjj4:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-19T19:23:50.729175Z"
                decision: revise
                reason: 'Captain rejected in chat: ''it''s not about asking permission. it''s about not having that in the fricking AC.'' The live run is in the test plan as a residual, not in an acceptance criterion. A test-plan line can be deferred; an AC cannot, because validation must reproduce every AC''s cited evidence and the gate cross-checks it. Add an AC whose verification IS a targeted local live run of rejection-flow, naming the journey and the expected observed codes. The workflow README now requires this of any task changing a live grader.'
            - id: gate-attempt:zf7rymtke3b6xp7r0337hjj4-ideation-3
              briefing:
                id: briefing:zf7rymtke3b6xp7r0337hjj4:ideation:attempt-3:revision-1
                digest: sha256:e959bc47e87f862029273c0adb1d78d33bd8bf4d9dd6e52a77f098b904d9bb5e
                request-digest: sha256:2a4a6012a3d7fd202f4d88a032f7ebf0f87be2d820e1b85a7ff5ca850619eb39
                room-ref: ./own-claude-early-rejection-round-record/review/ideation/briefing-3
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

**Cycle-2 supersession note.** The gate REVISED this design: the analysis below stands (executed-verified chain, both pins, the dead-code finding, the honest codes as a component), but shipping only the relabel was rejected — captain: "if it stays red it has no value. combine it with something that turns it green." The Acceptance criteria, Test plan, Expected surface, and the "Red stays red" tip guidance below are superseded by "Ideation design, cycle 2" further down. In particular the cycle-1 tip guidance (expect an honest claude red) is now WRONG — see cycle 2's tip reading.

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

## Ideation design, cycle 2 (2026-08-19): the precondition that turns the lane green

Captain's direction: the honest codes ship as a component, combined with a precondition on the round recorder so the early record becomes impossible rather than better-labelled. Record-after-correction was prose (`feedback-rejection-flow`, ignored twice live); this release's evidence is 4:0 that a binary refusal holds where prose loses. The five questions from the gate, each settled:

### 1. The predicate: what makes a round closed

Derived from the recorder's own model, not from the observed 2-entry log. `parseReviewLog` already computes the round's verdict: `reviewLog.Reviewer` is the FIRST Resolution in the log (`internal/gates/review.go:82-84`), entries are ordered (`includes` may only reference earlier identities), and a log with no Resolution is already refused (`:90-92`). The model's own decision dichotomy (`internal/gates/model.go:350-357`): `approve` stands alone; `revise` and `hold` demand rationale — verdicts that demand a response.

**A round log is CLOSED iff its final entry is a Resolution, and either that Resolution is not the verdict (`Reviewer`), or the verdict's decision is `approve`.** Refused, therefore, when the log ends with an Annotation (a dangling finding or disposition with no closing resolution) or ends exactly at a non-approve verdict (a demanded response that was never logged). No entry count, no actor identity, no fixture ID, no `includes`-graph semantics — the exact over-fit this family kept producing is structurally absent.

Shapes checked against the predicate: the live 2-entry reviewer log → refused (the defect); the complete 4-entry rejection log → accepted; the no-findings single-`approve`-Resolution round → accepted (an existing TESTED recorder shape, `round_test.go:189-206`, which a naive "two resolutions" rule would have broken); the default advisory testdata log, a decline-triage close ("0 material fixed; 2 declined") → accepted, so answering ≠ fixing; a hold-verdict tail → refused by the model's dichotomy (zero observed occurrences; recorded as model-derived, not observation-fitted).

### 2. Where it lives

`recordRoundLockedWith`, immediately after `loadValidateRound` (`internal/gates/round.go:92-95`), before pointer construction and `publishRound` — the same preflight-refusal band as the artifact/room checks, inheriting the byte-clean refusal contract the suite already pins (`TestRoundNoFindingsAndPreflightRefusals` asserts refusals leave the tree digest unchanged). Explicitly NOT in `loadValidateRound` itself: that helper is shared with `ValidateRoundFile` (`:126`), and read-time closure enforcement would make legacy early-recorded rooms unreadable — the durable oracle would red them under "validate retained round:" → `rejection-round-missing`, resurrecting the lying label for historical state. Record-time admission control; read-time stays permissive; the cycle-1 durable count check (`rejection-round-incomplete`) remains as the read-side honest backstop for legacy or bypassed state.

CLI: zero changes. The round branch already routes the error to stderr with exit 1 (`internal/cli/cli.go:356-359`). The refusal message must name what the tail is, that nothing was written, and the recovery in contract order — e.g. "round validation/1 is not closed: the review log ends at the reviewer's revise resolution; route the correction, append its triage entries (disposition annotations and a closing resolution), then record the round". Stable test token: "not closed".

### 3. Refusal, not record-and-mark-open

Record-open loses on the room's own model, not on taste: `publishRound` writes an immutable room exactly once, and the durable oracle plus the single-publication counter pin exactly one `round-1` room and one successful publication. An "open" record needs a second finalize write — either room immutability breaks or a new amend verb and lifecycle appear — and if the finalize never comes, the truncated record still lands durably: the defect tolerated with a marker, which the gate ruled out ("impossible rather than tolerated"). Refusal's cost is one exit-1 mid-cycle, and that cost is the mechanism working: the message carries the recovery, and the FO retries after routing the correction — the same move as j7j's gate attribution. Third alternative (a skill-prose precondition the FO checks itself) is the mechanism that already failed twice.

### 4. What stays red — falsifying cases

- Genuinely absent record — FO never invokes, or treats the refusal as a terminal hold and stops: no room, no pointer → `rejection-round-missing`, now honest (the round truly is missing).
- Durably incomplete room that exists anyway — legacy pre-precondition state, or a bypassed recorder: durable count check → `rejection-round-incomplete` (cycle-1 component, kept as backstop).
- Mutated `README.md`/`candidate.txt` → `rejection-workflow-file-mutated`; room-file tamper stays under the round oracle (cycle-1 split, unchanged).
- Open-log record attempts — reviewer-only tail, dangling annotation, hold tail: refused at record time, byte-clean (gates unit tests).
- Conforming or refusal-recovered run: GREEN. A failed attempt is not a publication (`claudeRejectionRoundPublications` counts only non-error results; codex counts exit-0 only), so refuse-then-recover still grades exactly one publication.

### 5. The tip reading FLIPS — correcting my cycle-1 guidance

Cycle 1 told the implementation to read claude-live RED under `rejection-round-incomplete` as success. That guidance is superseded and wrong under this design. With the precondition shipped the FO cannot produce the truncated room: an early attempt exits 1 with the recovery, so the run either records after the correction (GREEN) or never records (`rejection-round-missing`, honest red). **Expected claude tip result: GREEN.** Codex: green if `12z`'s tolerance lands, unchanged by this entity. Honest caveats: a green now also measures the FO following the refusal's recovery — a run that converts exit 1 into a terminal hold reds honestly as round-missing, a model-adherence residual the metrics instrument tracks; and `rejection-round-incomplete` should never appear in a post-fix live lane — if it does, the precondition was bypassed, which is itself a finding. The lane color is the journey's evidence; the fixture tests below are this fix's proof.

## Acceptance criteria (cycle 2 — replaces cycle 1's list)

- **AC-1 (measures the end-value).** `gate record --round` with an open log — the exact reviewer-only shape run 32270990171 recorded — exits non-zero, writes nothing (tree digest unchanged), and its error names the recovery; the same inputs with the correction's entries appended then record successfully. Independent baseline that can move the wrong way: today's recorder exits 0 on those bytes and durably truncates the room (proven live); wrong-way movements are acceptance (defect returns) or over-refusal (the closed log, the no-findings `approve` round, or the decline-triage close refused). Test: gates refusal + recovery + acceptance cases.
- **AC-2.** A durably incomplete round room grades status `fail` with codes exactly `[rejection-round-incomplete]`, detail naming got/want and the retained entry IDs. Baseline: the current grader yields `[rejection-round-missing]` / "resolved launcher never invoked" on the same state. Test: ensigncycle durable fixture (built by room-log overwrite — see Test plan 3).
- **AC-3 (means, paired with AC-2).** `claudeRecordedRejectionRound` recognizes the live run's command/result byte shape (`entries=2`) as a successful invocation; missing/failed results and wrong round/entity/artifact controls stay refused. Test: extractor case + existing controls; matcher-exec harness re-run flips false→true.
- **AC-4.** Mutating `README.md` or `candidate.txt` after a complete record grades `rejection-workflow-file-mutated`, not any `rejection-round-*` code; room-file tamper stays under the round oracle. Test: the three split-boundary cases (cycle-1 design, with the reordering ahead of `ValidateRoundFile`).
- **AC-5.** Conforming shapes stay accepted and the controls stay red: complete 4-entry record grades pass; no-findings `approve` round and decline-triage close still record; the no-invocation control still reds. Test: existing suites green unmodified plus the acceptance cases in AC-1.
- **AC-6 (the live journey, per the workflow README's Testing Resources requirement).** The rejection-flow journey, driven live against the built change, observes the expected codes — not "passes". Verification IS one targeted local run at implementation or validation: `SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=1 -timeout 30m -run '^TestLiveCommonRejectionFlow$' ./internal/ensigncycle/` with `SPACEDOCK_BIN` resolving to the built change and `SPACEDOCK_LIVE_ARTIFACT_DIR` set to retain the stream and topology digest. Expected observation: grade status `pass` with NO `rejection-round-missing`, NO `rejection-round-incomplete`, and NO `rejection-workflow-file-mutated` — the precondition makes the early record impossible, so a conforming run carries no rejection-round code at all. Non-green outcomes are graded findings, each reportable under its own meaning: `rejection-round-missing` observed → the FO produced no record — never attempted, or refused and did not take the recovery (model-adherence red, honest; the refusal-recovery residual measured); `rejection-round-incomplete` observed → a truncated room landed DESPITE the precondition — the refusal was bypassed or the fix did not land, a defect in this entity's own change (note: the gate feedback glossed this code as "refusal fired, FO did not take recovery"; under this design that state has no room and grades round-missing — the distinction matters because the two codes indict different things). Unavailability is named, never silent: a 429/quota or connection drop is a provider error, a stale `SPACEDOCK_BIN` is a setup failure — neither is a journey result and neither satisfies nor fails this AC; the run is retried or the blocker (sandbox, quota, unreachable host) is NAMED in the stage report along with what stays unproven — the FO's live adherence to the refusal recovery and the journey-level green, with the grader/recorder behavior then proven only offline. A named unavailability is a finding for the gate to weigh; silence is a gate rejection. The codex lane is out of this AC's scope (its isolated `CODEX_HOME` cannot be created in this sandbox — named here per the README rather than silently skipped; codex conduct is `12z`'s journey half).

## Test plan (cycle 2 — replaces cycle 1's)

1. **Gates refusal band** (extends `TestRoundNoFindingsAndPreflightRefusals`'s table): reviewer-only log → error contains "not closed", tree byte-clean; dangling-annotation tail → refused; hold-verdict tail → refused. Falsifying edit: deleting the precondition records a 2-entry room and reds the table.
2. **Gates recovery**: after the refusal, append the closing entries to the same log and re-record → success (proves the refusal leaves a recordable state).
3. **Durable honest code** (ensigncycle): the 2-entry room can no longer be built through `RecordSemantic` — that refusal is itself asserted by (1) — so build the legacy/bypass state by recording the COMPLETE log then overwriting the room's `briefing.review.jsonl` with the reviewer-only bytes; assert grade status `fail`, codes exactly `[rejection-round-incomplete]`. The count check at `:400-405` fires before the room byte block, so the state is unambiguous. A stream-only fix leaves this under `rejection-round-missing` → red.
4. **Stream + split cases**: cycle-1 items unchanged — `entries=2` extractor case, README/candidate/room-tamper boundary cases with the reorder falsifier.
5. `go test ./...` and `-race`; the targeted live rejection-flow run is AC-6's verification (not a deferrable plan line), executed at implementation or validation against the built change; the tip CI run remains the stack-level journey reading per the corrected reading above.

Cost: unit/fixture items are offline; AC-6 adds one targeted live run — minutes and one agent (the README's costing), against the 90-minute matrix behind a deployment approval. Live adherence to the refusal's recovery is the one thing offline tests cannot prove; AC-6 owns proving it.

## Expected surface (cycle 2 — replaces cycle 1's)

Estimate net LOC change: +150, across 3 files (~170 insertions, ~20 deletions). Breakdown — product: +20 (`internal/gates/round.go`: predicate function, call site, refusal message); tests: +130 (`internal/gates/round_test.go` ~+45: three refusal cases, recovery, acceptance controls; `internal/ensigncycle/shared_round_recording_test.go` ~+85: cycle-1 grading half with the durable fixture built per Test plan 3); docs: 0. Tolerance: ±60 net. Hard boundaries: product change confined to the round-record admission path in `internal/gates` — no CLI grammar, no new flags or verbs, no gate-ceremony or chat-decision path changes, no read-time (`ValidateRoundFile`) enforcement; no fixture-text changes; no skill prose changes (the refusal message is the binding surface).

Semantics declared: (1) runtime behavior — `gate record --round` now refuses an open round log (exit 1, zero writes) where it previously recorded it; command grammar, stored formats, and authority unchanged; recorder output format unchanged. (2) CI lane diagnostic codes change per cycle 1. (3) The early-record conduct's lane outcome changes from red to refused-then-recovered green — the defect becomes impossible, not tolerated.

## Spike determination (cycle 2)

No spike needed. Cycle-1's four proven mechanisms stand. The new mechanism rests on: `parseReviewLog`'s existing `Reviewer` = first-Resolution model and ordering guarantees (source-verified, `review.go:82-92`); the preflight-refusal band's byte-clean contract (tested, `round_test.go:255-262`); and the CLI's existing error plumbing (`cli.go:356-359`). Every existing round test exercises the same admission path the moment the check lands, so regressions surface immediately offline. The one unverifiable-offline element — the live FO following the refusal's recovery — is not spikeable (it is model adherence, not mechanism); it is owned by AC-6's targeted live run, not deferred to CI.

## Out of scope

- Rejection-flow fixture text (the Cycle-line target pin shipped with the 0.27 stack).
- The codex-side `rejection-worker-topology` count bar (`12z`'s entity, same stack).
- The identity half of the durable composite (`:400-403`) and its `%#v` message — never fired live.
- A dedicated code for round-room tampering (considered, declined in cycle 1: zero occurrences).
- The recorder's OUTPUT format (`round=… entries=N`) — unchanged; only admission of open logs changes. (Cycle 1 declared the whole recorder out of scope as "neutral"; the captain's direction supersedes that line — admission control is now the core of this entity.)
- New CLI flags/verbs, read-time closure enforcement, and any skill-prose additions (prose is the mechanism that already failed).

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

## Stage Report: ideation (cycle 2)

- DONE: Settle the remedy across BOTH entries=4 pins, not one: the stream regex at shared_round_recording_test.go:18 and the durable summary check at :400-405. A fix to only the first re-reds the same run under the same lying code with a raw %#v message.
  Cycle-1 settlement stands and now ships as the honest-code COMPONENT of a combined fix; the durable count check is additionally re-purposed as the read-side backstop for legacy/bypassed state the new precondition cannot reach.
- DONE: Decide the split criterion for the four-file byte-exact block at :425-442 (briefing.json, briefing.review.jsonl, candidate.txt, README.md). The review proposes "is this file part of the round record" rather than "is it README". Confirm or refute with reasoning.
  Unchanged from cycle 1: confirmed with the reachability correction (non-record checks must precede ValidateRoundFile), one code `rejection-workflow-file-mutated`.
- DONE: Design the proof so it cannot green a half-fix: a stream replay alone exercises only the stream path. Name the durable fixture that proves the count grade moved to the honest code.
  Revised: the precondition refuses recording a 2-entry log, so the durable fixture is now built by recording the COMPLETE log and overwriting the room's briefing.review.jsonl with the reviewer-only bytes (Test plan 3); the refusal of the direct construction is itself asserted (Test plan 1). Stream-only fix still reds Test plan 3.
- DONE: Confirm the red STAYS red. The conduct is genuinely wrong — the round room is immutable, so an early record permanently truncates the durable record. State how the composite stays red and what changes is only the diagnostic.
  Answered under the revised target: the early record is now REFUSED before any write, so it can no longer truncate anything; red is preserved for every genuinely bad END STATE (absent record → rejection-round-missing; legacy/bypassed incomplete room → rejection-round-incomplete; mutated workflow file → rejection-workflow-file-mutated), and the expected claude tip flips to GREEN — cycle 1's expect-honest-red guidance is explicitly superseded (design section 5).

### Summary

Redesigned per the gate's rejection: the relabel alone shipped no value, so the honest codes now combine with a record-time precondition on the round recorder — a round log is closed iff its final entry is a Resolution and it is either past the verdict or an approve verdict — derived from parseReviewLog's own Reviewer model and the model's approve/revise-hold dichotomy, count-free and actor-free. Settled all five gate questions: predicate, placement (recordRoundLockedWith's preflight band, record-time only so legacy rooms stay readable), refusal over record-and-mark-open (room immutability and single-publication invariants; prose already failed twice), the falsifying red cases, and the flipped tip reading (claude expected GREEN; my cycle-1 guidance corrected). Surface re-estimated honestly: +150 net across 3 files, product +20 in internal/gates/round.go, tolerance ±60.

## Stage Report: ideation (cycle 3)

- DONE: Settle the remedy across BOTH entries=4 pins, not one: the stream regex at shared_round_recording_test.go:18 and the durable summary check at :400-405. A fix to only the first re-reds the same run under the same lying code with a raw %#v message.
  Unchanged from cycle 2 (design accepted; this cycle's revision touched only the live-run AC placement).
- DONE: Decide the split criterion for the four-file byte-exact block at :425-442 (briefing.json, briefing.review.jsonl, candidate.txt, README.md). The review proposes "is this file part of the round record" rather than "is it README". Confirm or refute with reasoning.
  Unchanged from cycle 2.
- DONE: Design the proof so it cannot green a half-fix: a stream replay alone exercises only the stream path. Name the durable fixture that proves the count grade moved to the honest code.
  Unchanged from cycle 2, with the proof surface strengthened per this cycle's gate feedback: the targeted live rejection-flow run moved from Test plan item 5 ("declared as residual") into AC-6, whose verification IS the run with the README's exact command, named expected codes, and named-unavailability semantics.
- DONE: Confirm the red STAYS red. The conduct is genuinely wrong — the round room is immutable, so an early record permanently truncates the durable record. State how the composite stays red and what changes is only the diagnostic.
  Unchanged from cycle 2 (refusal makes the truncation impossible; every genuinely bad end state keeps its honest red; expected claude tip GREEN). AC-6 now also enumerates what each live-observed code would indict, including one accuracy correction to the gate feedback's gloss: a refused-but-unrecovered run grades `rejection-round-missing` (no room exists), while `rejection-round-incomplete` appearing live means the precondition was bypassed — a defect in this entity's own change.

### Summary

One structural addition per the gate: AC-6 makes the targeted local live rejection-flow run an acceptance criterion — naming the journey, carrying the README Testing Resources command verbatim (claude lane, sonnet, `^TestLiveCommonRejectionFlow$`), naming the expected observation as codes (pass with no rejection-round code at all) rather than "passes", mapping each non-green code to what it indicts, and defining unavailability semantics (429/quota and stale SPACEDOCK_BIN are never grades; a named blocker plus what-stays-unproven is a finding, silence is a rejection; codex lane named as sandbox-blocked and out of this AC's scope). Test plan 5, the cost line, and the spike determination re-pointed from "residual" to AC-6 ownership. Surface unchanged: the AC adds a run, not code.
