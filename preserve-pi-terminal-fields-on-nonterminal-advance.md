---
title: Preserve Pi terminal fields on nonterminal advance
status: implementation
score: "1.0"
source: Captain recovery directive; fh6 commit 4a98f40b4, 2026-08-11
sprint: test-behavior-completeness
sprint-readiness: ready
group: pi-product
id: kqdnfzjh921ryad7n6h82m1a
gates:
    version: 1
    records:
        - id: gate:kqdnfzjh921ryad7n6h82m1a:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:kqdnfzjh921ryad7n6h82m1a-backlog-1
              briefing:
                id: briefing:kqdnfzjh921ryad7n6h82m1a:backlog:attempt-1:revision-1
                digest: sha256:32e187c414902c2a7644a6db6ab005ed6909d8a31c409de008224bfa005637cf
                request-digest: sha256:592276c34d915145bb9c1d69da1e820b377379c81dd858f964a5d593e95e14b9
                room-ref: ./preserve-pi-terminal-fields-on-nonterminal-advance/review/backlog/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-11T00:41:41.736257Z"
                reason: seed lacks required backlog Stage Report
            - id: gate-attempt:kqdnfzjh921ryad7n6h82m1a-backlog-2
              briefing:
                id: briefing:kqdnfzjh921ryad7n6h82m1a:backlog:attempt-2:revision-1
                digest: sha256:7a96f5c08af39757fb5ce74b61f8c733a2ce4029204f66286cd12c0893e72f4e
                request-digest: sha256:f6e73e87b5350b673a1ff7f0a87aea5a506dd6e91ffdce3c7ce5aac6b28e2790
                room-ref: ./preserve-pi-terminal-fields-on-nonterminal-advance/review/backlog/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:kqdnfzjh921ryad7n6h82m1a:backlog:2
                briefing: briefing:kqdnfzjh921ryad7n6h82m1a:backlog:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-11T00:44:34.017994Z"
                decision: approve
                reason: The seed isolates one user-visible data-preservation defect, excludes mechanisms and live work, and defines falsifiable evidence for every acceptance criterion.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:kqdnfzjh921ryad7n6h82m1a:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:kqdnfzjh921ryad7n6h82m1a-ideation-1
              briefing:
                id: briefing:kqdnfzjh921ryad7n6h82m1a:ideation:attempt-1:revision-1
                digest: sha256:d155e03a82bca4007326625c56aa2f1b8af3233b7395d1f06ee7ea2dc3529470
                request-digest: sha256:74fbe0474f2374721cf53e6f128ee85c25dc01294dd87df9f986a27f70181b28
                room-ref: ./preserve-pi-terminal-fields-on-nonterminal-advance/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:kqdnfzjh921ryad7n6h82m1a:ideation:1
                briefing: briefing:kqdnfzjh921ryad7n6h82m1a:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-11T00:53:46.496189Z"
                decision: approve
                reason: The design reverses one destructive Pi instruction, uses an existing deterministic status boundary, defines a two-file +45/-1 baseline with bounded tolerance, and gives every AC falsifiable evidence without new mechanisms.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:kqdnfzjh921ryad7n6h82m1a:validation
          stage: validation
          attempts:
            - id: gate-attempt:kqdnfzjh921ryad7n6h82m1a-validation-1
              briefing:
                id: briefing:kqdnfzjh921ryad7n6h82m1a:validation:attempt-1:revision-1
                digest: sha256:81c4f32df3c595e5ee60a6613244aa0fbc16f365f3c2f96e312766a6c870fcd8
                request-digest: sha256:727bc33f674c42f52c630d958acc26de7538521a0fd3f7f3f4cf61f91abe7a9c
                room-ref: ./preserve-pi-terminal-fields-on-nonterminal-advance/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:kqdnfzjh921ryad7n6h82m1a:validation:1
                briefing: briefing:kqdnfzjh921ryad7n6h82m1a:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-11T01:41:41.198823Z"
                decision: approve
                reason: Independent validation passed. The candidate restores legitimate Pi workflow state without adding test-observability machinery or changing XFAIL ownership.
              application:
                target-stage: done
                state: superseded
started: 2026-08-11T00:45:13Z
worktree: .worktrees/spacedock-ensign-preserve-pi-terminal-fields-on-nonterminal-advance
mod-block:
---

Pi must not erase legitimate `completed` or `verdict` fields during a nonterminal advance.

## Problem

Commit `4a98f40b4` added a Pi First Officer clause that clears `completed` and `verdict` before a nonterminal status advance. Those fields can contain legitimate durable workflow state, so the automatic cleanup can destroy valid information.

## Proposed approach

Replace the current `«completion-signal»` bullet in `skills/first-officer/references/pi-first-officer-runtime.md` with the exact after text below. The replacement removes only the read/clear/reread/stop rule and its `pi-agent-teams` application. It retains the current completion signal, advisory, file-verification, stage-report-read, and task/member-completion rules.

Before:

```markdown
- `«completion-signal»`: For `pi-subagents`, the PRIMARY completion signal is the child return/status (`status: completed`); an optional advisory is only a non-blocking heads-up via raw `intercom send` before return (`contact_supervisor` carries no completion reason). Before a nonterminal status advance, the FO reads `completed` and `verdict`. If either field is non-empty, the FO runs `status --set <slug> completed= verdict=` once, reads the fields again, and stops if either field remains non-empty. The FO does not use `--force`. File verification remains the completion gate: the FO reads the entity file and reads the stage report before advancing state. For `pi-agent-teams`, apply the same field rule after task/member completion.
```

After:

```markdown
- `«completion-signal»`: For `pi-subagents`, the PRIMARY completion signal is the child return/status (`status: completed`); an optional advisory is only a non-blocking heads-up via raw `intercom send` before return (`contact_supervisor` carries no completion reason). File verification remains the completion gate: the FO reads the entity file and reads the stage report before advancing state. For `pi-agent-teams`, task/member completion is likewise verified against the entity file.
```

Extend `TestEnteredStageMutationControls` in `internal/status/entered_stage_test.go` with one table-driven preservation ladder. Each row builds the existing committed implementation fixture, performs the supported `status=validation` advance, and compares the raw `completed:` and `verdict:` lines before and after. Use one legitimate timestamp and the schema-cased `PASSED` verdict across the complete adjacent-state matrix: both empty, completed only, verdict only, and both nonempty. Every row must still assert exit 0 and `status: implementation -> validation`; clearing either populated line must fail its row.

This test mechanism serves AC-1. The simpler existing single empty/empty transition cannot detect data loss; a Pi/live run is unnecessary and forbidden because the supported status mutation supplies a deterministic on-disk observation boundary. The exact instruction diff serves AC-2. Blob equality against the ideation baseline serves AC-3; no new oracle mechanism is introduced.

## Expected surface and estimate

- `skills/first-officer/references/pi-first-officer-runtime.md`: replace one line, 1 insertion and 1 deletion.
- `internal/status/entered_stage_test.go`: add the four-row preservation ladder beside the existing committed-completion control, estimated 44 insertions and 0 deletions.
- Total estimate: 2 existing files, 45 insertions, 1 deletion. Tolerance: at most 2 files, 57 insertions, and 6 deletions. A third file, more than +12 insertions or +5 deletions, or any semantic change outside the declaration below requires a design reset before another candidate.

## Semantic scope

- **Runtime behavior changed:** Pi First Officers no longer erase nonempty `completed` or `verdict` values as a prerequisite to a nonterminal status advance. Existing status mutation behavior remains the authority: fields not named by the command remain unchanged.
- **Command grammar unchanged:** no CLI verb, option, argument order, exit code, or output contract changes.
- **Stored formats unchanged:** no frontmatter field, value convention, schema, fixture, gate record, or dispatch envelope changes.
- **Authority unchanged:** the First Officer remains the entity-state writer; the status binary retains its mutation guards; gate and captain authority do not move.
- **Other runtime behavior unchanged:** completion signaling, worker identity/lifecycle, file verification, stage-report verification, gate preparation/hold behavior, and Claude/Codex behavior remain unchanged.
- **Oracle and XFAIL surfaces unchanged:** preserve the current blobs of `internal/ensigncycle/claude_live_runner_test.go`, `claude_runtime_helpers_test.go`, `gate_assert_impl_test.go`, `gate_assert_test.go`, `pi_shared_live_runner_test.go`, and `shared_live_runner_test.go`, plus `internal/contractlint/live_registry_reconciliation_test.go`. This retains all fh6 assertion improvements and the current post-fh6 binding/owner state.

## Out of scope

No test-only product mechanism, hook, protocol, state store, parser loop, XFAIL mutation, live or Pi run, CI change, or unrelated fh6 change.

## Acceptance criteria

**AC-1 (VALUE) - A nonterminal Pi advance preserves legitimate `completed` and `verdict` fields byte-for-byte.**
Verified by: the four-row ladder drives the supported implementation-to-validation mutation and compares both raw lines before and after. It covers empty/empty, set/empty, empty/set, and set/set in one batch. Clearing the timestamp or `PASSED` line must make the corresponding row fail.

**AC-2 - The fh6 terminal-field-clearing instruction is absent while unrelated Pi runtime instructions remain unchanged.**
Verified by: the implementation diff must equal the exact one-line before/after replacement declared above, aside from the separately declared test file. Restoring any clear/reread/stop sentence or changing another Pi instruction fails the diff comparison. This is a one-off existence/scope check, not a committed prose-grep test.

**AC-3 - Useful fh6 oracle improvements and all XFAIL bindings, assertions, reconciliation rows, and owners remain unchanged.**
Verified by: the seven named oracle/registry files must retain their blobs from ideation baseline `ff9bb4506be73787a684e5fd80b7b772ea7473a5`; focused oracle, registry-reconciliation, and active-owner tests must pass. Any blob or ownership change fails the comparison.

**AC-4 - Repository behavior remains green after the narrow reversal.**
Verified by: focused tests, `go test ./...`, `go test ./... -race`, gofmt, and `git diff --check` on one immutable candidate.

## Spike determination

No spike needed. The design relies only on already-proven mechanisms: `runSet` applies named status updates without clearing unspecified fields; `TestEnteredStageMutationControls` already drives the committed implementation-to-validation transition through `runNative`; and the split-root fixture exposes the resulting entity bytes. The smallest falsifier is to make that status transition blank either unspecified terminal line: the new completed-only, verdict-only, or both-set row turns red. No parser round-trip, runtime handoff, new format, or tool flag is assumed.

## Test plan

1. Before the instruction edit, add the preservation ladder to `internal/status/entered_stage_test.go`. Run its exact focused test and demonstrate the falsifier by temporarily blanking an unspecified terminal line in the status mutation path; at least the applicable populated rows must turn red, then restore the mutation path.
2. Apply the exact instruction replacement. Run the focused entered-stage test and the relevant status package tests. No instruction-file grep is accepted as behavioral proof.
3. Compare the seven protected blobs to baseline `ff9bb4506be73787a684e5fd80b7b772ea7473a5`; run the focused gate-hold oracle, live-registry reconciliation, and active-owner checks to confirm their current semantics and ownership remain intact.
4. Run `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` on one immutable candidate. Gofmt must not leave unrelated changes.
5. Validation batches the complete four-row adjacent-state ladder and one adversarial falsifier before issuing one verdict. One candidate-owned rejection may receive one FO-authorized correction. A second candidate-owned rejection, any surface/tolerance breach, or any evidence fix that adds a mechanism triggers a design reset or HOLD; cycle 3 escalates to the captain. Preserve existing oracle and XFAIL ownership throughout both rounds.

## Stage Report: backlog

- DONE: A Pi nonterminal advance preserves legitimate `completed` and `verdict` bytes.
- DONE: The scope reverts only the terminal-field-clearing clause from fh6 commit `4a98f40b4`. It retains the oracle improvements and adds the smallest behavioral test.
- DONE: The exclusions are mechanisms, hooks, protocols, state stores, parser loops, XFAIL changes, live runs, Pi runs, and CI changes.
- DONE: AC-1 uses a focused behavioral test that proves byte-for-byte preservation of both terminal fields.
- DONE: AC-2 uses the exact source diff and focused contract checks to prove the narrow instruction reversal.
- DONE: AC-3 uses exact diff inspection and focused ownership checks to prove that the oracle improvements and XFAIL records are unchanged.
- DONE: AC-4 uses focused tests, full tests, race tests, gofmt, and `git diff --check` on one immutable candidate.

## Stage Report: ideation

- DONE: Declare the exact before/after instruction change, exact files, insertion/deletion estimate, tolerance, and unchanged semantic surfaces.
  Commit `d4fbb505c` specifies a two-file 45-insertion/1-deletion baseline, bounded tolerance, exact Pi bullet replacement, and unchanged semantic surfaces.
  AC-2 evidence: compare the candidate to `ff9bb4506be73787a684e5fd80b7b772ea7473a5`; only the exact Pi instruction replacement and `internal/status/entered_stage_test.go` may differ.
- DONE: Define the smallest behavioral test that fails when a legitimate completed or verdict value is cleared and covers the complete adjacent-state matrix in one ladder.
  One `entered_stage_test.go` ladder advances four empty/set combinations and compares raw terminal-field lines; blanking either populated value falsifies AC-1.
- DONE: Record no-spike-needed evidence or the smallest falsifier, preserve oracle/XFAIL ownership, and apply the two-round reset rule.
  Existing `runSet` and split-root fixtures provide the observation boundary; seven protected blobs stay fixed, and a second candidate-owned rejection requires reset or HOLD before cycle-3 escalation.
  AC-3 evidence: compare the seven named protected oracle/registry blobs and their active registry owners byte-for-byte to baseline `ff9bb4506be73787a684e5fd80b7b772ea7473a5`.
  AC-4 evidence: run the focused entered-stage and ownership checks, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` on one immutable candidate.

### Summary

Ideation narrows the recovery to one Pi instruction-line replacement and one deterministic status preservation ladder. It protects every fh6 oracle and current XFAIL/registry owner byte while making destructive clearing falsifiable across the full adjacent-state matrix.

## Stage Report: implementation

- DONE: Apply only the exact Pi completion-signal line replacement and the four-row entered-stage preservation ladder within the approved 2-file +57/-6 cap.
  Commit `578ed943a` changes two files with 49 insertions and 1 deletion. The Pi replacement matches the approved text.
  AC-2 evidence: The runtime diff is the approved one-line replacement. Only the runtime reference and entered-stage test changed.
- DONE: Prove the ladder can fail when an unspecified completed or verdict line is cleared, then keep all fh6 oracle blobs and XFAIL registry/owner bytes unchanged.
  A temporary `completed` clear failed the completed-only and both-nonempty rows. All seven protected blobs match baseline `ff9bb4506`.
  AC-1 evidence: Four rows preserve raw `completed:` and `verdict:` lines across empty/empty, completed-only, verdict-only, and both-nonempty advances.
  The temporary completed-clear fault failed two populated rows. The verdict-only and both-nonempty comparisons also detect a verdict clear.
  AC-3 evidence: The protected blobs for six `internal/ensigncycle` oracle files and `live_registry_reconciliation_test.go` match the baseline.
  `TestRuntimeLiveRegistryReconciliation` and `TestRuntimeLiveTODOOwnersAreActive` passed with the unchanged registry and active owners.
- DONE: Produce one immutable candidate with focused, full, race, gofmt, registry, active-owner, and diff checks green; run no live, Pi, or CI checks.
  The focused, race, gofmt, registry, active-owner, and diff checks passed on `578ed943a`. No live runtime, Pi session, or CI job ran.
  Two full runs found external suite load failures. Both isolated tests passed unchanged. The captain stopped additional verification.
  AC-4 evidence: `TestEnteredStageMutationControls`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` passed.
  The first `go test ./...` had a quiet timeout in `TestCodexProcessRequiresFinalMessageForTerminalTurn`. Its exact isolated run passed.
  The second `go test ./...` reached the 10-minute package timeout during `TestDurableKeepMovingDelayedPersistence`. Its exact isolated run passed.
  The Captain prohibited a third full run and timeout changes. No standalone output artifact file was recorded.

### Summary

Commit `578ed943a` removes the Pi field-clear instruction and adds the four-state preservation ladder. The race suite and all focused checks passed.
The full runs retained two external load artifacts: one quiet timeout and one 10-minute package timeout. Both isolated tests passed unchanged.

## Stage Report: validation

- DONE: Batch the complete four-row preservation matrix and one adversarial clearing falsifier before issuing one verdict; do not add cases one at a time.
  `TestEnteredStageMutationControls/successor_preserves_unspecified_terminal_fields` passed all four rows on `578ed943a`; one detached fault that cleared both unspecified fields failed completed-only, verdict-only, and both-nonempty while empty/empty stayed green.
- DONE: Verify the exact two-file +49/-1 candidate, approved instruction bytes, and all seven protected oracle/registry blobs plus XFAIL owners against baseline.
  `578ed943a` is exactly two files at +49/-1; the runtime line and approved after-text share SHA-256 `7d56f4b...a23a528`, all seven blobs equal `ff9bb4506`, and the reconciliation/active-owner tests passed.
- DONE: Independently assess AC-1 through AC-4 using the accepted focused/race/full-load evidence, run no redundant full/race/live/Pi/CI, and return one PASSED or REJECTED report under the two-round rule.
  PASSED in round 1: AC-1 through AC-4 have valid evidence, no material finding remains, and validation ran no full, race, live, Pi, or CI rerun.
  AC-1 evidence: The complete four-row preservation matrix passed, while the detached clearing falsifier failed every populated adjacent state.
  AC-2 evidence: The approved instruction replacement is byte-identical and the candidate has the exact two-file +49/-1 diff.
  AC-3 evidence: All seven protected blobs are baseline-identical, and the registry-reconciliation and active-owner checks are green.
  AC-4 evidence: The accepted focused and race evidence passed; both full-load failures passed unchanged in isolation under the captain's no-third-run ruling.

### Summary

PASSED. AC-1 is proven by the complete preservation matrix and a falsifier that turns every populated adjacent state red when clearing occurs; AC-2 is proven by the byte-identical approved instruction replacement and exact two-file diff.
AC-3 is proven by seven baseline-identical blobs plus green registry and active-owner checks. AC-4 uses the accepted focused/race/full-load record: focused and race evidence passed, and both load-time full-suite failures passed unchanged in isolation under the captain's no-third-run ruling; reviewer findings are none, with no deferred risk or polish item.

## Stage Report: implementation (cycle 2)

- DONE: Restore `internal/status/entered_stage_test.go` exactly to recovery base `c5d2e79c4fa220bb859c13495d180c9450c937e8`.
  Candidate `9c2cda2aa` and the recovery base share test-file blob `54583521989afbfa382cc64ccc88d65fd9eb6c47`.
- DONE: Preserve only the approved Pi completion-signal instruction replacement.
  The candidate diff against the recovery base is exactly one file with 1 insertion and 1 deletion.
- DONE: Add no replacement test or evidence mechanism.
  The Captain rejected the 48-line test expansion. This correction removes it without replacement.
- DONE: Run only the permitted formatting and diff checks.
  `gofmt -w ./cmd ./internal`, `git diff --check`, exact blob equality, and exact diff-scope checks passed. No test ran.

### Summary

Cycle 2 removes the rejected preservation ladder and retains only the approved Pi instruction change. Earlier implementation and validation reports describe the rejected snapshot.
Candidate `9c2cda2aa` is ready for independent validation against recovery base `c5d2e79c4fa220bb859c13495d180c9450c937e8`.
