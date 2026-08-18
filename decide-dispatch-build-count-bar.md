---
id: 6htv3p97aq8sexrkjhrcdwt1
title: Decide whether one dispatch build per stage is the right bar for a live agent
status: validation
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
        - id: gate:6htv3p97aq8sexrkjhrcdwt1:validation
          stage: validation
          attempts:
            - id: gate-attempt:6htv3p97aq8sexrkjhrcdwt1-validation-1
              briefing:
                id: briefing:6htv3p97aq8sexrkjhrcdwt1:validation:attempt-1:revision-1
                digest: sha256:2c2f2daaee0b592d2c39d7c32539423f13f6669115aa7b7b4ee3e41261d695e7
                request-digest: sha256:409a41753ab8c84933c943f5854d8960f1eb06249c2c5dfbf63d8d06c755b6d3
                room-ref: ./decide-dispatch-build-count-bar/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6htv3p97aq8sexrkjhrcdwt1:validation:1
                briefing: briefing:6htv3p97aq8sexrkjhrcdwt1:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T21:55:08.075992Z"
                decision: approve
                reason: 'Captain approved in chat: ''approve all PR and trigger ci on tip.'' Validation PASSED with two deferred risks, no material finding: 12-shape matrix confirms the widened bar still reds flailing, waste and final failure; falsifiability reproduced under a strict-1/1 revert; +94-vs-+30 overrun traced line by line to approved design and test-plan items, one file, zero production code.'
              application:
                target-stage: done
                state: pending
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

## Stage Report: implementation

- DONE: The build-count bar tolerates exactly one corrective rebuild; identical rebuilds, three-plus attempts, and final failures still grade red.
  `recordedGateBuildAttemptsAcceptable` (internal/ensigncycle/recorded_gate_lifecycle_test.go:86-108) replaces `builds != 1 || successfulBuilds != 1`; wired into `assertRecordedGateLifecycle`. `recordedGateDispatchProof.attempts []recordedGateBuildAttempt{command, ok}` replaces the two counters; `recordedGateLiveObservation`'s per-line loop now retains each begin/exit pair instead of collapsing them.
- DONE: Both observed occurrences grade green (claude 32105482382 --stamp then --bare-mode; codex 30754109029 error-then-retry) and the six control shapes grade as spiked, 7/7.
  `TestRecordedGateBuildCountBar` table-drives both real occurrences (argv verbatim, cited to run IDs) plus single-clean-build/zero-build/failed-build/identical-two-build-waste/three-build-flailing — 7/7 as spiked. Falsifiable: temporarily reverting `recordedGateBuildAttemptsAcceptable` to strict `len==1 && attempts[0].ok` fails both real-occurrence rows and the new corrected-rebuild real-CLI control (verified locally, then reverted). Also added two new positive log-mutation controls inside `TestRecordedGateLifecycleRealCLIReplay` (corrected-rebuild, error-then-retry), and confirmed the existing scripted controls — zero-build, failed-build, build-before-consume, missing-ancestry, identical-two-build (now the waste control) — stay red unchanged.
- DONE: Report actual net surface against the approved estimate of +30 across 1 file.
  Actual: 1 file (`internal/ensigncycle/recorded_gate_lifecycle_test.go`), +110/-16 (net +94) — 3x the ideation estimate. The estimate's own breakdown (proof-struct swap, loop rewrite, assertion+comment, "new positive controls") under-budgeted two items the design and test plan both required for AC-2's falsifiability: the 7-row table-driven unit test citing both real run IDs (test plan item 1, ~35 lines) and the two new positive log-mutation controls inside the real-CLI-replay integration test (test plan item 3, ~18 lines). No scope beyond the approved design was added; the gap is estimation, not scope creep.

### Summary

Widened `assertRecordedGateLifecycle`'s build-count check per ideation's decision: exactly one successful build, or exactly one corrective rebuild (last attempt succeeded, and either the prior attempt failed or its command differed). `go test ./...` and `go test ./... -race` both pass for `internal/ensigncycle` (and every other package except one pre-existing, unrelated failure — `TestCodexResolveManifestAgainstInstalledHost` in `internal/cli`, reproduced identically against unmodified `main`, caused by local codex plugin-cache state on this machine, not this diff). `gofmt -w ./cmd ./internal` is clean. Committed to `spacedock-ensign/decide-dispatch-build-count-bar` at `cce97b8fd`, a clean diff against `main` touching only `internal/ensigncycle/recorded_gate_lifecycle_test.go`.

## Review-finding disposition

- **Reviewer (validation, 2026-08-18) — Deferred risk: the waste control is raw-string equality, so a cosmetic argv difference evades it.** Released user and normal workflow: an FO dispatching a successor in the live recorded-gate-lifecycle journey. Observable harm: a duplicate build carrying no correction grades green when its argv differs only cosmetically — the same flags reordered, a doubled space, or a repeated boolean flag (audit probes E/F/G). Affected authority: `value-ac[AC-2]` — not violated; AC-2 promises the *byte-identical* rebuild stays red, and it does (probe A and the live two-build control at `recorded_gate_lifecycle_test.go:357`). Trigger evidence: none observed — the FO's build argv comes from the fixed-order template at `skills/first-officer/references/fo-dispatch-core.md:151`, and both recorded occurrences (32105482382, 30754109029) show stable flag order. Promotes to material on any recorded occurrence whose two build argvs differ only by token order, whitespace, or a duplicate flag; the narrow fix is to compare normalized argv fields rather than the raw string.
- **Reviewer (validation, 2026-08-18) — Deferred risk: the rule cannot separate a corrected rebuild from two distinct dispatches.** Released user and normal workflow: the same journey. Observable harm: a second successful build naming a different `--stage` or `--entity-path` grades green (audit probe H), where the old 1/1 bar reddened it. Affected authority: `value-ac[AC-2]` — not violated; AC-2's unacceptable list is identical rebuild, three-plus attempts, and final failure, and a distinct-target double dispatch is not on it. Trigger evidence: unreachable in this journey — the scenario drives one entity through one gate to exactly one successor stage (`claude_live_runner_test.go:123-146`, `expectedNext: "handoff"`), and the sibling `durableEffects == 1` assertion independently reds any second stamped build or second dispatched worker. Promotes to material on a live occurrence with two builds naming different `--stage`/`--entity-path`; the narrow fix is to require both attempts to agree on those two flags.

## Stage Report: validation

- DONE: Test the overrun claim: the surface came in at net +94 against an approved +30, and the implementation asserts the gap is estimation error with no scope beyond the approved design. Trace every added line to the approved ideation design or test plan; recommend REJECTED if scope beyond the approved design landed under cover of the estimate miss.
  Overrun confirmed and benign. `git diff --numstat $(git merge-base main HEAD)..HEAD` at cce97b8fd: `110 16 internal/ensigncycle/recorded_gate_lifecycle_test.go` — 1 file (estimate held), net +94 vs +30 (3.1x), and zero production code (the sole changed file is a `_test.go`). Every insertion traces: 34 lines of `recordedGateBuildAttempt` + `recordedGateBuildAttemptsAcceptable` + its reasoning comment (design bullets 1-2, incl. "the rule's reasoning is recorded in a comment"); 8 for the assertion swap and new error message (bullet 3); 7/-6 for the observation loop retaining begin/exit pairs (bullet 1); 3 mechanical `builds`→`len(attempts)` call-site updates forced by the struct swap; 21 for `buildLine` plus the two positive log-mutation controls (test plan item 3, verbatim); 40 for `TestRecordedGateBuildCountBar` (test plan item 1, verbatim). Nothing outside the approved design. Deletions landed at 16 against an estimated 15 — the accurate deletion budget corroborates that the miss is entirely unbudgeted new tests (~61 of 110 insertions are test-plan items 1 and 3), not a wider rewrite.
- DONE: Semantic adversarial pass on the widened bar — this is the real risk. The new rule accepts "last attempt succeeded, and either the prior attempt failed or its command differed." Prove it still reds genuinely bad conduct: two different failing commands then a success; a success followed by a redundant differing rebuild; a differing rebuild that changes nothing material. If any of those now grade green, that is a material finding.
  Ran a 12-shape matrix through the real `recordedGateLiveObservation` + `assertRecordedGateLifecycle` on the real-CLI replay fixture, in a throwaway checkout (since removed). RED: byte-identical successful rebuild (waste), three successes (flailing), two DIFFERENT failing commands then a success, success then a failed differing rebuild, success then a build that started and never finished — plus the five pre-existing scripted controls. The bar reds nothing-becomes-green: the named "two different failures then a success" shape is 3 attempts and hard-red. GREEN: cosmetic-difference rebuilds (flag reorder / doubled space / repeated boolean flag) and a rebuild naming a different `--stage`. Both are recorded above as deferred risks, not material — AC-2's promised waste control is the byte-identical one and it still reds, and neither green shape has an observed or reachable trigger. A second green shape (a build line containing " --help" is excluded from attempts) is pre-existing: the identical `!strings.Contains(line, " --help")` guard predates this diff and is load-bearing for the accepted capability probe at `claude_runtime_helpers_test.go:504`.
- DONE: Independently reproduce the falsifiability claim: reverting recordedGateBuildAttemptsAcceptable to strict len==1 && attempts[0].ok must fail both real-occurrence rows and the corrected-rebuild control. Run go test ./internal/ensigncycle/... -race and confirm the pre-existing scripted controls stay red.
  Reproduced in the throwaway checkout: with the rule reverted to `len(attempts) == 1 && attempts[0].ok`, `TestRecordedGateBuildCountBar/claude-32105482382-corrected-rebuild` and `/codex-30754109029-error-then-retry` both FAIL, and `TestRecordedGateLifecycleRealCLIReplay` fails at "corrected-rebuild control failed to qualify" — so the two occurrence rows and the replay control each fail on the exact change that re-tightens the bar. The five control rows (zero-build, single-clean-build, failed-build, identical-two-build-waste, three-build-flailing) still PASS under the revert, which is correct: they must be red under both bars, so they cannot be the rows carrying the widening. Unreverted, `go test ./... && go test ./... -race` pass for `internal/ensigncycle` and every package except `TestCodexResolveManifestAgainstInstalledHost` in `internal/cli`, which I reproduced identically on unmodified `main` in a separate checkout (it depends on this machine having no codex-installed spacedock); `gofmt -l ./cmd ./internal` is empty.

### Summary

PASSED, with two deferred risks. The +94-vs-+30 overrun is an honest estimation miss: the file count held, no production code changed, and every added line maps to an approved design bullet or a verbatim test-plan item — the estimate simply never budgeted test-plan items 1 and 3. Both value ACs verify with reproduced evidence: AC-1's zero-unreasoned-assertions count holds (the only exact-count dispatch-build assertion in the live journeys now carries recorded reasoning citing both run IDs; the sibling hold journey asserts ordering, not counts, and already tolerates a repeated envelope), and AC-2's occurrence rows are genuinely falsifiable under a strict-1/1 revert while every shape AC-2 promises to redden still reddens. The widened bar keeps real teeth — flailing, waste, and final failure all stay red — but its correction test is raw-string inequality, so cosmetic argv differences and distinct-target rebuilds slip through; neither has an observed or reachable trigger today, so both are recorded as deferred risks with concrete promote conditions rather than blocking findings.
