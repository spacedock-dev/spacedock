---
id: 9xnaq83nryb38fyt18mh0gbt
title: The completion-guard error names the failing sub-check
status: validation
source: "External field report 2026-08-26 (/tmp/spacedock-durability-guard-defect.md, WhisperLiveKit sandbox session): three human round-trips to diagnose a dirty entity file because one generic error covers four distinct failures; the reporter misread the guard as requiring a remote push — refuted by source read and a live no-remote repro on 0.28.0-pre0"
started: 2026-08-26T19:51:27Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-completion-guard-error-names-failing-subcheck
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:9xnaq83nryb38fyt18mh0gbt:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9xnaq83nryb38fyt18mh0gbt-backlog-1
              briefing:
                id: briefing:9xnaq83nryb38fyt18mh0gbt:backlog:attempt-1:revision-1
                digest: sha256:484cdf1ea7ca83bc9d95c2302368bb5ca70c64fe31e2220a3b6e37e4449115f5
                request-digest: sha256:65459f1836c9fb15d4ea01d5a55184a3581fd59098ffdf0833bbec2b0dbb7473
                room-ref: ./completion-guard-error-names-failing-subcheck/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9xnaq83nryb38fyt18mh0gbt:backlog:1
                briefing: briefing:9xnaq83nryb38fyt18mh0gbt:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T19:50:43.055728Z"
                decision: approve
                reason: The seed has a bounded diagnostic outcome, preserves the existing guard invariants, and names falsifiable ideation proof for each failure class.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:9xnaq83nryb38fyt18mh0gbt:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:9xnaq83nryb38fyt18mh0gbt-ideation-1
              briefing:
                id: briefing:9xnaq83nryb38fyt18mh0gbt:ideation:attempt-1:revision-1
                digest: sha256:7c3d06cf5e1fa2dcec86f81bff1a82d94a3598e493cff0eb7744366856bece47
                request-digest: sha256:bbeddb51f2622ef1af95fe8ae5c7d5f9b432710bee1bd7ae4076488e189b9722
                room-ref: ./completion-guard-error-names-failing-subcheck/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-26T20:08:11.020575Z"
                reason: Captain requested a decision-focused primary artifact with the entity retained only as supporting reference.
            - id: gate-attempt:9xnaq83nryb38fyt18mh0gbt-ideation-2
              briefing:
                id: briefing:9xnaq83nryb38fyt18mh0gbt:ideation:attempt-2:revision-1
                digest: sha256:da5664c8495a4e9a55f2d02b499fe4831a705594ef73f25fe70dfb23f6615c24
                request-digest: sha256:1422ffb53cec73320d3a8d2496b6d1baf7d303f09ddfcbcac83115703bf7d89e
                room-ref: ./completion-guard-error-names-failing-subcheck/review/ideation/briefing-2
              withdrawal:
                by: agent:first-officer
                at: "2026-08-26T20:25:52.132529Z"
                reason: Rebuild the prepared room with the local gate-room-v1 binary so Subspace can validate index.json.
            - id: gate-attempt:9xnaq83nryb38fyt18mh0gbt-ideation-3
              briefing:
                id: briefing:9xnaq83nryb38fyt18mh0gbt:ideation:attempt-3:revision-1
                digest: sha256:8b9808bcd11396d088220ead5047865fee7f7b15270bf8d9e8fefd2e4cbd5e06
                room-ref: ./completion-guard-error-names-failing-subcheck/review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:9xnaq83nryb38fyt18mh0gbt:ideation:3
                briefing: briefing:9xnaq83nryb38fyt18mh0gbt:ideation:attempt-3:revision-1
                by: person:captain
                at: "2026-08-26T20:30:45.290883Z"
                decision: approve
              application:
                target-stage: implementation
                state: consumed
---

The stage-advance guard emits one message — "cannot change status away from entered stage until a durable, complete ## Stage Report is committed" — for four distinct failures: the stage-report heading is not found (exact-token match), the checklist is incomplete (missing bullets, FAILED items, or no Summary), the entity file is untracked, or the entity file differs from local HEAD. An operator who cannot tell which failed diagnoses by round-trip; one field session needed three.

## Problem

`status --set` currently reduces four completion checks to one boolean and one
generic refusal: a matching current-stage report exists, its checklist and Summary
are complete, the entity is tracked, and its bytes equal local `HEAD`. The field
report therefore misdiagnosed a dirty entity as an unpushed commit and spent three
human round trips on remote pushes that the guard never inspects.

The current parser matches the first whitespace-delimited stage token, not the
whole heading. Thus `## Stage Report: implementation (cycle 2)` is valid, while
`implementation-notes` is not. The new message must describe that real boundary
instead of repeating the field report's refuted exact-whole-heading theory.

Two related vocabulary failures are folded in. A status-like checklist bullet such
as `- DONE (annotation):` is ignored because the colon is not immediately after
`DONE`; an otherwise valid bullet can then let an omitted obligation pass unseen.
The folded field report also records merge guard exposing `condition "ineligible"`
for an ungated terminal path without naming the direct terminal `status --set` form.
That exact source shape has since received a semantic repair, so ideation must not
reintroduce it merely to improve its wording.

## Proposed approach

Replace the completion predicate's internal boolean with a small structured result
(`kind`, line/item detail, Git root) and retain a boolean wrapper for scheduler and
gate-readiness callers. `status --set` renders the first failure in existing guard
order, so precedence and all current fences remain stable:

1. No heading whose first stage token equals the current stage: name the required
   token and show the canonical heading.
2. Checklist incomplete: name the first malformed/FAILED/blank/unevidenced item,
   no recognized items, or a missing/empty `### Summary`, including a line number
   when one exists.
3. Entity untracked: name the path and enclosing local Git root, and say to add and
   commit that path.
4. Entity dirty: name the path and Git root, say it differs from local `HEAD`, and
   explicitly state that this guard does not require a remote push.

Add a narrow near-miss scan beside the existing checklist parser. A bullet beginning
with `DONE`, `SKIPPED`, or `FAILED` but not the canonical immediate-colon form is a
checklist-incomplete failure, with the original line and corrected form in stderr.
This intentionally tightens one edge: a report containing a valid item plus an
ignored annotated near-miss no longer advances. Valid syntax, stage-token selection,
local-HEAD durability, sibling-dirt isolation, and `--force` behavior do not change.

For the folded terminal case, retain the repaired behavior rather than adding a dead
error branch. A consumed nonterminal approval is ordinary history, so both merge
guard and the direct ungated terminal write succeed without `--force`; neither may
emit raw `ineligible`. Add the missing merge-guard regression beside the existing
direct-write regression, and document the concrete direct form:
`spacedock status --workflow-dir DIR --set SLUG status=TERMINAL completed verdict=PASSED worktree=`.
An actually unreadable, stale, superseded, or pending authority remains fail-closed
and must not receive that hint: on those shapes the suggested direct write would
itself refuse, so printing it would be false guidance.

Mechanisms and necessity:

- The structured result serves AC-1/AC-4. The simplest alternative is a generic
  suffix listing all four remedies; it is insufficient because it preserves the
  guess-and-retry loop and can still recommend a push for a syntax failure.
- The near-miss scan serves AC-2. A warning-only alternative is insufficient because
  the command could still advance while silently omitting the annotated obligation.
- The ungated regression plus concrete documentation serves AC-3. Adding a formatter
  for residual `ineligible` errors is insufficient and unsafe because current source
  already succeeds for the reported shape, while genuine residual authority errors
  also block direct `status --set` and cannot truthfully recommend it.

## Risk evidence

The external report records three unnecessary push round trips. Source inspection
shows the guard runs only `git ls-files` and `git diff --quiet HEAD`; it has no remote
operation. `go test ./internal/status -run TestEnteredStage -count=1` passed on
2026-08-26 and exercises complete/incomplete/dirty real-Git states;
its committed-completion fixture has no origin and advances successfully. The focused
checklist tests also passed and prove that `(cycle 2)` suffixes are selected by the
leading stage token.

A throwaway CLI spike built a real gate, approved and consumed it into an ordinary
stage, then ran merge guard through the ungated terminal transition. It exited 0 and
reported `finalized`; the spike was removed after the run. Commit `7a0de4018` and
`consumedNonterminalHistory` explain why: the reported `ineligible` failure has already
been repaired. No further mechanism spike is needed because the remaining design
reuses exercised parser, CLI, and local-Git mechanisms. Its first red tests are the
four-message table and annotated near-miss; the spike becomes the durable AC-3
regression.

## Out of scope

Relaxing the exact stage-token match, checklist evidence rules, or clean-vs-local-HEAD
requirement; adding remote/sync coupling; changing JSON/status table output; changing
gate authority; adding speculative merge-guard error branches; or broadly redesigning
merge guard eligibility. Unrelated Git command failures may get an honest "unable to
inspect" diagnostic but must remain fail-closed.

## Expected surface and tolerance

Estimate net LOC change: +145, across 6 files. Expected insertions: 168; deletions:
23. Tolerance: net LOC +/-40 and file count +/-1; crossing either bound requires a
correction round before implementation continues.

Expected files: `internal/status/entered_stage.go`, `gate_extract.go`, `handlers.go`,
`entered_stage_test.go`, `internal/cli/terminal_consume_test.go`, and
`docs/site/reference/command-reference.md`. Command grammar, stored formats, and gate
authority do not change. Observable stderr becomes failure-specific; annotated
status-token near-misses newly refuse; consumed nonterminal history remains a successful
ungated terminal path; all other guard pass/fail outcomes, exit codes, byte-clean
refusals, and local-only durability remain unchanged.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A blocked stage advance names the failing sub-check, and the four core failure fixtures produce four distinct remedies instead of one generic message.**
Verified by: a table test drives real `status --set` calls for missing current-stage
heading, incomplete checklist/Summary, complete-but-untracked entity, and
complete-but-dirty entity. It asserts exit 1, empty stdout, unchanged bytes, exact
class-specific stderr, and four unique messages. Collapsing any result to the old
generic string makes the exact and uniqueness assertions fail; baseline distinct
message count is 1 and required count is 4.

**AC-2 (VALUE) - No status-like checklist obligation is silently omitted: canonical bullets retain current behavior, while each named near-miss is refused with its line and repair.**
Verified by: a table covers canonical `DONE`/`SKIPPED`, `FAILED`, blank text, missing
evidence, missing/empty Summary, and `- DONE (annotation):`. The annotated case also
contains one valid DONE item, so deleting the near-miss scan incorrectly advances and
fails the test. A spaced heading suffix remains accepted; changing selection to an
exact whole-heading comparison fails that control.

**AC-3 (VALUE) - The reported ungated terminal journey has zero `ineligible` refusals, and users have a concrete direct terminal command, while a gated terminal approval still routes only through merge guard.**
Verified by: two real CLI fixtures modeled on `/tmp/spacedock-gate-merge-frictions.md`
each prepare, approve, and consume a gate into an ordinary stage. One runs merge guard;
the other runs the documented direct `status --set` form. Both must exit 0 and produce
terminal on-disk state, and neither channel may contain `ineligible`. The existing
pending-terminal-approval control continues to refuse direct status and route through
merge guard. Removing `consumedNonterminalHistory` makes both ordinary paths red;
weakening sole-consumer authority makes the gated control red.

**AC-4 - Existing completion and authority invariants remain intact outside the declared near-miss tightening.**
Verified by: existing scheduler/gate-readiness matrices plus a no-origin real-Git
fixture prove local committed completion passes, sibling dirt is ignored, entity-path
dirt blocks, latest stage cycle wins, and `--force` behavior is unchanged. Full normal
and race suites must pass; introducing a remote query or accepting dirty bytes makes
the focused fixtures fail or hang on a repo with no origin.

The concrete command-reference change is:

```diff
- ... `merge guard` discovers/arms the delivery mechanism when it acts.
+ ... `merge guard` discovers/arms the delivery mechanism when it acts. For an
+ ungated current-stage-to-terminal transition, finalize directly with
+ `spacedock status --workflow-dir DIR --set SLUG status=TERMINAL completed
+ verdict=PASSED worktree=`; do not use that route for a pending terminal-target
+ approval, whose sole consumer remains `merge guard`.
```

## Test plan

- Add the four-row behavior table in `entered_stage_test.go` using real temporary Git
  checkouts. Cost: small, under five seconds; CLI fixtures are required, no live host.
- Add parser/guard rows for checklist defects and annotated near-misses. Cost: small;
  unit plus CLI fixture. Assertions include the mutation that would make each red.
- Extend the existing consumed-nonterminal CLI fixture with a merge-guard subtest;
  retain the direct-write and pending-terminal controls. Inspect resulting frontmatter
  and archive state rather than only searching prose. Cost: medium; no network/live run.
- Run `go test ./internal/status -run 'TestEntered' -count=1` and
  `go test ./internal/cli -run 'TestConsumedNonterminal|TestTerminalDelivery' -count=1`, then
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- Review the command-reference diff against the exercised stderr/command. No golden
  prose-only test is proposed because the CLI behavior tests are the independent
  oracle and a static substring check would be tautological.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Stage Report: ideation

- DONE: Define a bounded, behavior-first diagnostic design for every named completion-guard failure while preserving the existing guard invariants and resolving the two folded error-vocabulary cases.
  The body specifies ordered diagnostics, the one declared near-miss tightening, an ungated regression plus concrete guidance, and unchanged authority/local-HEAD fences.
- DONE: Declare the expected net LOC change, insertions, deletions, files, tolerance, affected observable semantics, simplest alternative, and why it is insufficient.
  The expected surface is +145 net (168 insertions, 23 deletions) across six files with +/-40 net and +/-1 file tolerance; each mechanism has a rejected simpler alternative.
- DONE: Specify falsifiable proof for the four core failure classes, checklist near-misses, and ungated-terminal guidance, including spike/no-spike evidence and any concrete user-facing documentation change.
  AC-1 through AC-4 name behavior fixtures and falsifying edits; Risk evidence records the exercised no-origin baseline, and the command-reference diff is concrete.

### Summary

Ideation now defines a bounded structured diagnostic without relaxing the completion
or gate-authority fences. It corrects the heading-boundary misconception, makes
annotated checklist near-misses fail visibly, and turns the already-repaired ungated
journey into a durable regression plus concrete command guidance.

## Stage Report: implementation

- DONE: Implement the four failure-specific completion diagnostics, annotated checklist near-miss refusal, terminal regression, and command documentation while preserving the approved authority and local-HEAD boundaries.
  Commit `d24e90a8e` adds ordered structured failures, strict status-bullet scanning, both consumed-history terminal journeys, and the documented direct command without changing gate writers or Git pathspecs.
- DONE: Add falsifiable behavior coverage for AC-1 through AC-4: four distinct byte-clean refusals, checklist syntax controls, ungated merge/direct completion, and the pending-terminal sole-consumer control.
  `TestEnteredWorktreeStageFailureDiagnosticsAreDistinctAndByteClean` fails if messages collapse, bytes mutate, or the annotated item advances; heading controls fail if token selection changes, terminal tests fail without consumed history, and the pending control fails if direct completion spends authority.
- DONE: Keep the change within +145 net lines across six files with the approved tolerance, run gofmt plus focused, full, and race suites, and report exact surface and results.
  Surface is 6 files, 216 insertions/41 deletions, +175 net; gofmt and both focused commands passed, `go test ./...` passed on rerun after two transient ensigncycle fixtures passed in isolation, and `go test ./... -race` passed.

### Summary

The implementation now reports the first actual completion failure and explicitly
separates local HEAD durability from remote publication. Annotated checklist
near-misses refuse with a line and repair, while the existing authority boundary is
pinned by successful ungated direct/merge journeys and the pending-terminal refusal.

## Stage Report: validation

- DONE: Reproduce AC-1 and AC-2 from real status commands: prove four byte-clean failure classes have exact distinct remedies, canonical checklist forms retain behavior, and the annotated near-miss refuses with its line and repair instead of silently advancing.
  `TestEnteredWorktreeStageFailureDiagnosticsAreDistinctAndByteClean` ran real `status --set` calls and would fail on a changed byte, nonempty stdout, wrong exit, collapsed remedy, missing line/repair, or an advancing annotated obligation; canonical DONE/SKIPPED and spaced-cycle controls also passed.
- DONE: Verify AC-3 and AC-4 end to end: both consumed-nonterminal terminal journeys succeed with zero ineligible refusals, pending terminal approval remains merge-guard-only, local-HEAD/no-origin and dirty-path boundaries remain intact, and applicable focused, full, race, and formatting checks actually run green.
  Direct and merge-guard consumed-history fixtures reached terminal state without `ineligible`, the full documented direct form stayed byte-clean with pending authority, focused/full/race suites passed, and gofmt left every candidate-changed Go file unchanged.
- DONE: Perform the declared detached semantic adversarial audit from a throwaway checkout, trace diagnostic identity/order/bytes and authority state across adjacent variants, challenge how the tests could pass with wrong observable behavior, and report PASSED or REJECTED with every finding classified under the workflow policy.
  Detached audit proved report/checklist/Git precedence, exact Unicode near-miss identity and repair, spaced Unicode cycle acceptance, full pending-terminal refusal bytes and authority, and found no material, deferred-risk, or polish candidate finding; recommendation: PASSED.
- DONE: AC-1 evidence is external and falsifiable.
  The four-row CLI table asserts exact stderr and four unique values plus exit 1, empty stdout, and unchanged entity bytes; restoring the generic diagnostic or reordering the guard makes it red.
- DONE: AC-2 evidence is external and falsifiable.
  Parser and CLI controls cover canonical DONE/SKIPPED, FAILED, blank text/evidence, absent/empty Summary, exact stage-token selection, and a valid item followed by the annotated near-miss; deleting the scan incorrectly advances.
- DONE: AC-3 evidence is external and falsifiable.
  Real prepare/approve/consume journeys exercise direct status and merge guard, inspect terminal fields and consumed authority, while the pending-terminal control refuses the documented full write and preserves pending authority byte-for-byte.
- DONE: AC-4 evidence is external and falsifiable.
  `go test ./internal/status -run 'TestEntered' -count=1`, focused terminal tests, `go test ./...`, and `go test ./... -race` passed; no-origin, sibling-dirt, entity-dirt, latest-cycle, and force controls remain covered.
- DONE: Candidate surface remains within the approved ideation tolerance.
  Commit `d24e90a8e` changes the declared six files by 216 insertions and 41 deletions (+175 net), within the +145 +/-40 net and six +/-1 file bounds.
- DONE: Changed hot paths have no new scaling or blocking-I/O risk.
  Completion still reads one entity and runs the same two path-scoped Git commands; parsing adds bounded linear passes over those bytes, so no over-limit test is warranted.
- DONE: The command-reference guidance matches the exercised authority boundary.
  It gives the successful ungated full status form and excludes pending terminal-target approval, whose real full-form control still routes only through merge guard.
- DONE: Reviewer-finding disposition and delivery recommendation.
  Outcome defects: none. Evidence defects: none. Material findings: none. Deferred risks: none. Polish findings: none. The detached checkout exposed unrelated pre-existing gofmt drift in `internal/release/runtime_live_evidence_workflow_test.go`; it is outside commit `d24e90a8e` and does not affect this candidate recommendation.

### Summary

Validation recommends PASSED for commit `d24e90a8e`: AC-1 through AC-4 all have
behavioral, state-based evidence, and the detached audit found no candidate defect.
The candidate preserves diagnostic order, local-HEAD durability, and terminal gate
authority while making the named failure and repair observable.
