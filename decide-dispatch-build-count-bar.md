---
id: 6htv3p97aq8sexrkjhrcdwt1
title: Decide whether one dispatch build per stage is the right bar for a live agent
status: implementation
source: "Captain CL, 2026-08-18, from the live-lane inventory reframe. The recorded-gate-lifecycle journey red in run 32105482382 on 'successor dispatch build attempts/successes = 2/2, want 1/1' — two SUCCESSFUL builds, not an error-then-retry. codex-live-dispatch-build-checklist-race already carries the open question and is codex-scoped; this occurrence is claude."
started: 2026-08-18T18:41:27Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-decide-dispatch-build-count-bar
issue:
gates:
    version: 1
    records:
        - id: gate:6htv3p97aq8sexrkjhrcdwt1:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6htv3p97aq8sexrkjhrcdwt1-ideation-1
              briefing:
                id: briefing:6htv3p97aq8sexrkjhrcdwt1:ideation:attempt-1:revision-1
                digest: sha256:915fe5c741a3e479b5062e7db78257c983a18bde293e13c406e6848eda0a2dcb
                request-digest: sha256:7e287fe25f888e8a8b50294bf23107dbdffa759fa2e4cd1756594805576628d7
                room-ref: ./decide-dispatch-build-count-bar/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6htv3p97aq8sexrkjhrcdwt1:ideation:1
                briefing: briefing:6htv3p97aq8sexrkjhrcdwt1:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T19:31:48.647812Z"
                decision: approve
                reason: 'Captain approved in chat: ''approve those 4, and have them be on a pr stack.'' Accepts the decision that 1/1 does not stand and the bar widens to one corrective rebuild, with identical rebuilds and three-plus attempts staying red.'
              application:
                target-stage: implementation
                state: consumed
---

A journey grades an FO red for building a dispatch envelope twice. Ideation's decision: 1/1 does not stand — the retained streams show the second successful build differed from the first (a corrected rebuild), so the bar widens to tolerate exactly one corrective rebuild while identical rebuilds, three-plus attempts, and final failures stay red.

## Problem

The recorded-gate-lifecycle journey asserts `dispatch.builds/successfulBuilds == 1/1` for a successor dispatch. Run 32105482382 produced `2/2` on claude — two builds that both succeeded, which rules out the benign error-then-retry sub-mode the sibling entity describes.

The question the bar begs has already been written down and left open. `codex-live-dispatch-build-checklist-race` asks whether `1/1` is the right bar for a live agent versus a scripted CLI replay. It is scoped to codex and describes `2/1` with a failed first attempt. Neither the host nor the sub-mode matches this occurrence.

`repair-opus-recorded-gate-lifecycle` is opus-scoped, names no mechanism, and is priority-held. It does not cover sonnet.

So the failure is real in the sense that the count exceeded the bar, and unproven in the sense that nobody has established the bar is right. A live FO that builds an envelope, reconsiders, and builds again with corrected inputs may be behaving well. Or it may be wasting a build and signalling confusion. The grader currently assumes the second without evidence.

Deciding this matters more than fixing this instance: the same bar now has two hosts and two sub-modes failing against it, and every future occurrence inherits the assumption.

## Proposed approach

**Decision: 1/1 does not stand. Widen the bar to tolerate exactly one corrective rebuild; keep everything else red.** A successor dispatch grades green when its build attempts are either one successful build, or two attempts whose last succeeded and which form a correction — the first attempt failed (codex's error-then-retry), or the second changed the command (claude's mode correction). Two identical successful builds stay red (pure waste, no correction signal), as do three-plus attempts (flailing) and a final failed attempt. One uniform rule for live and scripted replay: all four existing scripted controls (zero-build, failed-build, identical two-build, build-before-consume) stay red under it unchanged.

**The decisive evidence: the second successful build differed from the first.** Claude, run 32105482382, retained `command.log` and `claude-stream.jsonl`: the FO ran `dispatch build ... --stamp` (exit 0; stamped state — commit 27c43f6 "dispatch: recorded-gate-task entering handoff"), reconsidered the transport shape, re-ran `dispatch build ... --bare-mode` without `--stamp` (exit 0, no second commit, different `dispatch_file_path` in the emitted envelope), verified in-stream "Good — one dispatch commit only", and dispatched the bare envelope via a blocking `Agent()` call. Exactly one durable stamp, ordering intact; the journey's only red was the count. Codex, run 30754109029 (sibling entity's verbatim quote; the retained artifact is the green re-run): the identical `--bare-mode --stamp` command twice, first exit 1 (checklist file not yet written), then exit 0. Both occurrences are self-corrections of an imperfect first command, not duplicate transactions.

**Why the exact count is an assumption, not a bar.** Each property the value transaction needs already has its own assertion in `assertRecordedGateLifecycle`: the envelope is tool-built (a successful build produced the dispatched prompt), exactly one durable successor effect (`durableEffects == 1`), build after consume (`ordered`), consumed state committed before dispatch (`committed`). The exact count adds one real protection — detecting waste and flailing — and the corrected-rebuild rule keeps that protection while admitting observed good conduct. The instructions already demand a single build in the adapter-selected shape (`fo-dispatch-core.md:17` "Run «dispatch.build» --stamp ... immediately"; `:163` "invoke that shape once"; `claude-fo-dispatch.md:9` single-entity mode uses bare-mode), and two hosts under two different sub-modes still slipped benignly — a grader that reddens benign self-correction measures instruction-adherence noise, not product behavior. Instruction polish stays with the sibling entity.

**Design** (all in `internal/ensigncycle/recorded_gate_lifecycle_test.go`):

- Replace `recordedGateDispatchProof.builds/successfulBuilds` with `attempts []recordedGateBuildAttempt` (`{command string, ok bool}`). `recordedGateLiveObservation` already walks per-line argv in command.log (lines 1073–1092); retain each begin/exit pair instead of collapsing to two counters. Mechanism justification: this serves AC-2; the simplest alternative — keep the counters and allow `builds <= 2` — is insufficient because it grades the identical-successful-rebuild waste case green, destroying the existing two-build control. The other alternative — keep 1/1 and only fix FO instructions — is insufficient because the instructions already mandate one build and both hosts slipped anyway; the journey would stay flaky-red on good runs (AC-2's failure mode).
- `assertRecordedGateLifecycle` grades: one successful attempt, or two attempts with the last successful and (first failed or argv differs); everything else red. The rule's reasoning is recorded in a comment at the assertion site — AC-1's recorded reasoning lives where the next reader of the bar looks.
- Error message names the new bar (e.g. "successor dispatch build attempts = N (M ok), want one build or one corrective rebuild").
- A live-versus-scripted bar fork was considered and rejected: the uniform rule preserves every scripted control as-is, so a second grader semantics buys no protection.

**Spike (done in ideation).** The rule was replayed over the retained claude `command.log` bytes from run 32105482382 (attempt extraction mirroring the existing begin/exit pairing → `[--stamp ok, --bare-mode ok]` → GREEN) and over the codex pair reconstructed verbatim from the sibling entity's quoted evidence (fail-then-identical-retry → GREEN), plus single-clean-build (GREEN) and zero-build, failed-only, identical-two-build, three-build shapes (all RED) — 7/7 as intended. This table seeds the implementation's first unit test.

**Sibling fold.** This decision subsumes `codex-live-dispatch-build-checklist-race`'s open bar question ("is 1/1 the right bar for a live agent"): its occurrence grades green under the new bar. Its other open question — whether ensign instruction ordering invites the write-command-before-file race — remains its own, codex-scoped. `repair-opus-recorded-gate-lifecycle` (opus-scoped, priority-held) is untouched.

**No doc diff needed:** no user-facing doc describes the 1/1 bar (checked `docs/runtime-support.md` and `docs/specs/`); the change is grader-internal test semantics.

## Out of scope

Changing dispatch build itself. The error-then-retry sub-mode owned by the sibling entity, except where this decision subsumes it.

## Expected surface and tolerance

The evidence read decided it: the bar changes. Estimate net LOC change: +30 across 1 file (`internal/ensigncycle/recorded_gate_lifecycle_test.go`) — approximately 45 insertions, 15 deletions (proof-struct swap, observation loop retaining per-attempt records, assertion rule plus reasoning comment, new positive controls). Semantic changes: grader semantics of the recorded-gate-lifecycle journey and its scripted replay assertion only — no command grammar, no stored formats, no runtime authority, no user-visible behavior.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - The bar is justified by evidence rather than assumed.**
This is the measuring AC: the count of dispatch-build assertions in the live journeys whose expected value has no recorded reasoning must be ZERO. Verified by the decision recorded against both occurrences' actual build inputs, showing whether the second build differed from the first. Fails if the bar is restated without establishing what a second successful build means.

**AC-2 - A live FO that behaves acceptably grades green.**
Verified by replaying both recorded occurrences against the decided bar as unit-test rows: the claude 2/2 corrected rebuild (run 32105482382, `--stamp` then `--bare-mode`) and the codex 2/1 error-then-retry (run 30754109029, identical command, exit 1 then 0) must pass; the unacceptable shapes — an identical successful rebuild (the existing two-build replay control), three-plus attempts, and a final failed attempt — must fail. Fails if the outcome is a bar that no observed live run can satisfy, which would make the journey permanently red rather than informative, or one that admits the waste control.

## Test plan

Go unit tests in `internal/ensigncycle` (a grader change is proven by fixtures, not a live run; the next scheduled live lanes provide the confirming green). Cost: small, no new test infrastructure.

1. Table-driven test over attempt lists mirroring the spike's seven shapes, with the two observed occurrences' argv verbatim (cited to their run IDs). Falsified by re-tightening to 1/1 (both occurrence rows go red) or over-widening (waste/flail rows go green).
2. Existing scripted replay controls unchanged and still red: zero-build, failed-build, build-before-consume, missing-ancestry, and the identical second build (`recorded_gate_lifecycle_test.go:301` — it becomes the waste control under the new rule). Falsified by a tolerance that admits identical rebuilds.
3. New positive log-mutation controls in the replay: duplicate the valid log's build pair with one flag changed → lifecycle passes; mutate the first pair to exit=1 with an identical retry appended → lifecycle passes. Falsified by an assertion that still requires exactly one attempt.

Run `go test ./... && go test ./... -race` per repo policy.

## Stage Report: ideation

- DONE: Read both occurrences' actual build inputs FIRST and say whether the second successful build differed from the first — the answer decides whether 1/1 is a real bar or an assumption, and everything else follows from it.
  Claude run 32105482382 retained command.log: build 1 `--stamp` (stamped, commit 27c43f6), build 2 `--bare-mode` (no stamp) — inputs DIFFERED, a corrected rebuild; codex run 30754109029: identical command, exit 1 then 0 (the retained artifact is the green re-run; the failing pair survives verbatim in the sibling entity's source).
- DONE: Bring the proposal early: state the recommended bar and its reasoning before elaborating the design, because the captain wants the decision surfaced, not buried.
  The entity lede and the Proposed approach's first line both state the decision: tolerate exactly one corrective rebuild; identical rebuilds, three-plus attempts, and final failures stay red.
- DONE: If the honest outcome is that 1/1 stands and this host folds into the sibling entity, say so and ship zero code — that is a legitimate result for this task.
  The honest outcome is the opposite: 1/1 misgrades observed good conduct on two hosts under two sub-modes, so the bar changes (net +30, 1 file); this decision subsumes the sibling's bar question while its instruction-ordering diagnosis stays codex-scoped.

### Summary

Pulled the retained streams for both occurrences and answered the decisive question: the claude 2/2 was a corrected rebuild (`--stamp` then `--bare-mode`, one durable stamp, FO verified "one dispatch commit only" in-stream), the codex 2/1 an error-then-retry — both benign self-corrections, so the exact 1/1 count is an assumption that reddens good conduct. Spiked the proposed corrected-rebuild rule against the real claude command.log bytes plus six control shapes (7/7 correct: both occurrences green, zero-build/failed-build/identical-two-build/three-build red), preserving all existing scripted controls unchanged. Ideation records the widened bar, its design in assertRecordedGateLifecycle (per-attempt records replacing the two counters), a fixture-only test plan, and the surface: net +30 across 1 file.
