---
title: Clean recorded-gate lifecycle proof hygiene before v1
status: done
source: "Durable-decisions exact-tip close-out audit at deac7f8a, corrected and committed as 4ff98d8c."
score: "0.9"
sprint: durable-decisions
id: fh3n4w4jg7tk015512tn1tsd
gates:
    version: 1
    current:
        gate: gate:docs-dev:fh3n:validation
    records:
        - id: gate:docs-dev:fh3n:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:fh3n-backlog-1
              briefing:
                id: briefing:docs-dev:fh3n:backlog:attempt-1:revision-1
                digest: sha256:c51c3611a1855ff5eab3eed3507559b676a6748faa133150ca681f7a02984de1
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fh3n:backlog:1
                briefing: briefing:docs-dev:fh3n:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T16:53:40.939705Z"
                decision: approve
                reason: The exact-tip audit reproduces both proof defects, existing executable controls preserve lifecycle value, and the accepted cleanup adds no product behavior or standing obligation.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:docs-dev:fh3n:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:fh3n-ideation-1
              briefing:
                id: briefing:docs-dev:fh3n:ideation:attempt-1:revision-1
                digest: sha256:fcf51be338b2fa7929c8b54ba6ea5577653f13866877bb0364e497f86b03fe64
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fh3n:ideation:1
                briefing: briefing:docs-dev:fh3n:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T17:08:33.585106Z"
                decision: approve
                reason: The corrected design removes the two proven proof-hygiene defects within an exact four-test-file, no-product boundary; executable lifecycle owners remain, and LOC is an estimate rather than gate authority.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:docs-dev:fh3n:validation
          stage: validation
          attempts:
            - id: gate-attempt:fh3n-validation-1
              briefing:
                id: briefing:docs-dev:fh3n:validation:attempt-1:revision-1
                digest: sha256:a21ab988482ef787427d6735d3911787e724924753fbb27bfe27e17973cb6c50
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fh3n:validation:1
                briefing: briefing:docs-dev:fh3n:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T18:24:19.911705Z"
                decision: approve
                reason: Independent validation proves the exact deletion-first four-test-file candidate satisfies all acceptance criteria, structural and archive controls, repository gates, and final Roborev; the credentialed journey red is byte-identical pre-existing canonical-open oracle debt, not a candidate regression.
              application:
                target-stage: done
                state: consumed
worktree: .worktrees/spacedock-ensign-recorded-gate-lifecycle-proof-hygiene
started: 2026-07-25T17:09:51Z
review-round:
    id: round:fh3n4w4jg7tk015512tn1tsd:implementation:2
    stage: implementation
    cycle: 2
    briefing:
        id: briefing:fh3n4w4jg7tk015512tn1tsd:implementation:round-2
        digest: sha256:b009e6c313cfdc0319281f1a02be8b788a698f6cf9cef90d51d4f1ebe4149eeb
        digest-domain: canonical-bytes
        room-ref: ./review/implementation/round-2
mod-block:
pr: pr-merge:567
verdict: passed
completed: 2026-07-25T19:04:46Z
archived: 2026-07-25T19:04:46Z
---

## Problem

At current `main` (`4ff98d8c`), the recorded-gate lifecycle has executable command and live-journey coverage, but two proof-hygiene defects obscure that evidence:

1. The prose-as-proof defect has two directly verified sites. `internal/ensigncycle/recorded_gate_lifecycle_test.go:568-598` reads the shipped `skills/fo-gate-lifecycle/SKILL.md` and deletes command superstrings before asking its own substring parser whether the corresponding substring disappeared. `procedureEvents` at `:650-667` has only the two call sites inside that test. The check runs no product behavior and survives meaning-inverting prose.
2. `internal/ensigncycle/codex_live_runner_test.go:210` reads the active entity before `:217-219` compares those bytes with `resolveRecordedGateEntity`. If the entity was archived, the read fails first; otherwise the resolver re-reads the same active file. The named archive assertion therefore cannot diagnose the failure it names.

The directly inspected adjacent site is `internal/contractlint/layering_restore_test.go:184-197`. Its three-file `path -> want prose` map is a shipped-prose grep forbidden by `internal/contractlint/doc_test.go:11-13`. Git history supplies the deletion-safe boundary: restore only the prior `c049d1eef` structural positive discriminator that verifies the allow-listed `mods/pr-merge.md` really contains `gh pr view`. Do not delete or weaken that structural owner.

## Current executable owners

Deleting prose checks does not delete lifecycle proof:

- `TestRecordedGateLifecycleRealCLIReplay` drives the real binary through briefing record, decision record, consume, commit ordering, and successor dispatch, with negative controls.
- `TestRecordedGateLifecycleMissingEventControls`, `TestRecordedGateLifecycleAC5RefusalMatrix`, and `TestRecordedGateLifecycleAC7ResumeMatrix` fail on missing, refused, duplicate, or resumed lifecycle events.
- `TestRecordedGateLifecycleTerminalConsumeHasNoDispatchableSuccessor` runs the terminal consume path and observes `dispatchable: []`.
- The Claude, Codex, and Pi recorded-gate live journeys grade command execution and durable state produced from the shipped skill.
- `TestMergeGuardFinalizesFromMergedSentinelNonArmed`, `TestMergeGuardFinalizesTerminalUnblockedEntityFromMergedSentinel`, `TestMergeGuardBlocksOnOpenPRNoModBlock`, `TestPRIndicatesMerged`, and the CLI merge-guard end-to-end tests execute the terminal merge behavior formerly represented by prose-map values.
- `TestNoUnexpectedPRViewScanIntroduced`, `TestPRViewAllowListIsLoadBearing`, and `TestPRViewAllowListConstrains` continue to own the structural `gh pr view` quarantine and its non-vacuity controls.

## Proposed approach

Treat the work as two proof-hygiene corrections, with the verified contractlint cleanup included in the prose-as-proof correction:

1. Delete `TestRecordedGateLifecycleCommandTextMutants` and `procedureEvents` without replacement. Keep `recordedGateRepoRoot`, `readFile`, all real-CLI tests, lifecycle observation/assertion helpers, and live journeys.
2. Restore `TestPRViewAllowListIsLoadBearing` to its prior structural-only form: read `mods/pr-merge.md` and require the allow-listed `gh pr view` token. Delete the three-file shipped-prose map. Add no prose token, behavior claim, detector, cap, or owner.
3. In `runCodexGateGuardrailScenario`, check the known folder-form archive path with `os.Stat` immediately after the runner returns and before reading `fixture.entity`. Delete the post-read `resolveRecordedGateEntity(fixture) != after` comparison. Keep `resolveRecordedGateEntity` because `recordedGateLiveObservation` uses its archive fallback at `recorded_gate_lifecycle_test.go:851`.
4. Add the same pre-read folder-archive absence assertion to `runClaudeGateGuardrailScenario`. That function owns both the shared `gate-guardrail` journey and `TestLiveDefaultHeadlessStopsAtGate`, so parity lands only where the same scenario already exists.

The archive assertion is against concrete on-disk state:

```go
if _, err := os.Stat(filepath.Join(fixture.stateRoot, "_archive", "recorded-gate-task", "index.md")); !os.IsNotExist(err) {
    t.Fatalf("recorded-gate-task was archived while waiting at the gate; stat err=%v", err)
}
```

No spike is needed: the repo already uses this explicit `os.Stat`/`os.IsNotExist` pattern in both live runners, the fixture has a deterministic folder-form archive path, the audit deletion experiment passed, and the prior structural-only contractlint implementation is in git history.

## Alternatives rejected

- Moving the command-text mutant into `internal/contractlint` would legalize the file read location but preserve the banned prose-grep tautology.
- Expanding the AST instruction-read detector would add a standing enforcement mechanism to remove one known orphan and would not prove lifecycle behavior.
- Comparing `resolveRecordedGateEntity` after reading the active path retains the unreachable diagnostic; reading only through the resolver would also allow archival instead of asserting its absence.
- Deleting all of `TestPRViewAllowListIsLoadBearing` would remove the positive discriminator for a real allow-list entry. Restoring its prior structural-only form deletes the prose map while preserving that executable owner.
- Adding prompt wording, a new behavior detector, a compatibility rule, or a component cap would create obligations unrelated to either defect.

## Expected surface

Baseline: `main` at `4ff98d8c`.

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | `+0/-49` | Delete the command-text mutant and orphaned parser only. |
| `internal/ensigncycle/codex_live_runner_test.go` | `+3/-3` | Replace the post-read comparison with a pre-read archive absence assertion. |
| `internal/ensigncycle/claude_live_runner_test.go` | `+3/-0` | Add the same pre-read assertion to the existing scenario runner. |
| `internal/contractlint/layering_restore_test.go` | `+8/-12` | Separately restore the prior structural positive discriminator and delete the three-file prose map. |
| **Total** | **`+14/-64`, net `-50` LOC** | Four test files; no product or instruction change. |

Estimate/checkpoint: the table's `+14/-64` (net `-50`) comes from exact current-main line ownership; it is not an acceptance band or an automatic design-reset trigger. The four declared test files and zero product/instruction changes remain the binding surface. At implementation review, report the actual numstat against this estimate; explain and triage any implementation-surface drift for materiality against the deletion-only task intent. LOC variance alone does not fail the design. The fourth file is the pre-authorized direct-evidence branch.

## Acceptance criteria

1. **The finished diff touches exactly the four declared test files, removes the named proof-hygiene constructs, and changes zero product, command, skill, prompt, provider, compatibility, or lifecycle-AC files.**
   Test: `git diff --numstat "$(git merge-base main HEAD)"..HEAD` and `git diff --name-only` show the declared surface and actual LOC beside the `+14/-64` estimate; exact-name searches show the mutant and orphaned parser are gone, and any numeric variance is explained and triaged for materiality rather than failed by threshold.
2. **No shipped-prose mapping remains in `TestPRViewAllowListIsLoadBearing`, while its structural `mods/pr-merge.md` / `gh pr view` positive discriminator and negative control remain executable.**
   Test: the focused contractlint command below passes; deleting `gh pr view` from the allowed mod makes the positive discriminator fail, while planting it outside the allow-list makes the negative control fail.
3. **Both Codex and Claude guardrail journeys test archive absence before reading the active folder entity, and an archived fixture produces the named archive diagnostic rather than a generic active-file read error.**
   Test: run the three targeted live commands below; a temporary control that moves the entity to `_archive/recorded-gate-task/index.md` after the runner returns must fail at the named pre-read assertion.
4. **Every executable lifecycle and supported runtime outcome that passed before the cleanup still passes after it.**
   Test: the focused real-CLI/refusal/resume suite, live-tag compile, targeted live journeys, full suite, and race suite below pass. Removing or reordering record/consume events, allowing terminal successor dispatch, or archiving at the held gate makes an owner test fail.

## Test plan

Run focused executable owners before the edit and again after it:

```bash
go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle(RealCLIReplay|TerminalConsumeHasNoDispatchableSuccessor|AC5RefusalMatrix|AC7ResumeMatrix|ProvenanceAndPresentationMutants|MissingEventControls)$' -count=1
go test ./internal/contractlint -run 'Test(NoInstructionReadsOutsideQuarantine|NoUnexpectedPRViewScanIntroduced|PRViewAllowListIsLoadBearing|PRViewAllowListConstrains)$' -count=1
go test -tags live ./internal/ensigncycle -run '^$' -count=1
```

Exercise the existing host journeys when their credentials are available (CI remains the required host-backed lane otherwise):

```bash
go test -tags live ./internal/ensigncycle -run '^TestLiveCodexSharedScenarios$/^gate-guardrail$' -count=1
go test -tags live ./internal/ensigncycle -run '^TestLiveClaudeSharedScenarios$/^gate-guardrail$' -count=1
go test -tags live ./internal/ensigncycle -run '^TestLiveDefaultHeadlessStopsAtGate$' -count=1
```

Then run repository gates in the required order:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

At implementation end, after the branch diff and tests are final, run `roborev review --branch --panel branch_final`. Triage every finding against the declared surface and acceptance criteria before making edits; record a feedback cycle if a finding triggers another pass. Roborev does not authorize a new standing obligation or surface expansion.

## Stage Report: ideation

- DONE: The design removes only the two proven proof-hygiene defects while preserving every executable lifecycle oracle and accepted runtime outcome.
  Current-main call sites were verified at `4ff98d8c`; focused lifecycle owners passed, and the adjacent prose map was admitted only through its pre-authorized direct-evidence branch.
- DONE: Each proposed assertion is falsifiable against on-disk state or command execution, with no shipped-prose behavioral grep and no new standing enforcement.
  Archive absence uses the fixture's real `_archive/recorded-gate-task/index.md`; real CLI, merge-guard, allow-list mutation controls, and host journeys own behavior.
- DONE: The exact current-main files/LOC, focused tests, full/race gates, and Roborev boundary are small enough for direct implementation review.
  Expected surface is four test files, `+14/-64` (net `-50`, tolerance `-44..-56`), with the contractlint `+8/-12` delta recorded separately.

### Summary

Ideation pinned a deletion-first, test-only cleanup against current `main`: remove the tautological lifecycle prose oracle, restore the adjacent contractlint test to structural-only proof, and make archive diagnostics direct and pre-read in both existing host runners. No product behavior or new obligation is proposed; executable lifecycle, merge, and live-journey owners remain the acceptance evidence.

### Feedback Cycles

- Cycle 1: REJECTED — first officer; surface 0 code files/0 implementation LOC vs estimate 4 code files/+14/-64 (ideation correction); AC unchanged in task intent: removed the invented LOC band while preserving the exact four-file/no-product boundary.

## Stage Report: ideation (cycle 2)

- DONE: The design removes only the two proven proof-hygiene defects while preserving every executable lifecycle oracle and accepted runtime outcome.
  The accepted four-file direction is unchanged; no product, instruction, runtime outcome, or executable owner was added or removed in this correction.
- DONE: Each proposed assertion is falsifiable against on-disk state or command execution, with no shipped-prose behavioral grep and no new standing enforcement.
  The archive and command oracles remain executable; LOC is now observed and triaged, not enforced as a new standing gate.
- DONE: The exact current-main files/LOC, focused tests, full/race gates, and Roborev boundary are small enough for direct implementation review.
  `+14/-64` remains the review checkpoint, while actual variance requires explanation and materiality triage rather than automatic rejection.

### Summary

Cycle 2 preserves the exact four-file, test-only design and removes the unauthorized `-44..-56` acceptance band. This report supersedes only the prior report's numeric tolerance claim; file scope, executable evidence, full/race gates, and the Roborev boundary remain unchanged.

## Stage Report: implementation

- DONE: The finished diff touches exactly the four declared test files, removes the named proof-hygiene constructs, and changes no product or instruction files.
  Commit `1fa8ead5` is `+13/-63` (net `-50`): the one-line-per-side variance from `+14/-64` is immaterial line accounting, with identical four-file scope and deletion-first net change; exact-name searches find neither removed construct.
- DONE: The structural `mods/pr-merge.md` / `gh pr view` positive discriminator remains executable while the shipped-prose map is gone.
  Focused contractlint passed; deleting the allowed token fails the retained positive discriminator, while its existing planted-path negative control still rejects an outside occurrence.
- DONE: Codex and Claude guardrail journeys check archive absence before reading the active entity and produce the named archive diagnostic.
  A temporary post-run archive control failed both paths with `recorded-gate-task was archived while waiting at the gate; stat err=<nil>` before any active-file read; the control file was removed before commit.
- DONE: Focused lifecycle and contractlint tests, live-tag compilation, gofmt, full tests, race tests, and a final Roborev panel pass or have explicit evidence-backed dispositions.
  Focused suites, live-tag compilation, `go test ./...`, and `go test ./... -race` passed; Roborev job 2331 passed 2/2 with no findings and is recorded as advisory round `implementation/2`.
- DONE: Exercise the credentialed Codex, Claude, and default-headless guardrail journeys when credentials are available.
  All three reached and presented the recorded gate, then failed the unchanged hold-shape oracle because `state: open` count was not one; each executed the new pre-read archive assertion first without firing, so this host-output drift is non-material to the approved proof-hygiene cleanup.

### Summary

The committed four-test-file cleanup removes the tautological command-text oracle and shipped-prose map, preserves the structural PR-view discriminator, and makes both guardrail runners diagnose archival before active-file reads. Deterministic gates are green; the credentialed host runs exposed an unrelated exact-state-shape drift after clearing the new assertion, and the final two-reviewer Roborev panel found no issues.

## Stage Report: validation

- DONE: Reproduce the exact four-test-file deletion-first diff and verify the removed prose-as-proof/parser constructs, retained structural PR-view discriminator, and pre-read archive diagnostic controls against acceptance criteria 1–3.
  `4ff98d8c..1fa8ead5` is exactly the declared four test files at `+13/-63`; both removed names are absent, and the candidate changes no product or instruction file.
- DONE: Reproduce focused lifecycle/contractlint tests, live-tag compilation, full tests, race tests, and the relevant host-journey evidence; determine from this candidate’s evidence whether the reported hold-shape drift is material or a separately triggered deferred risk.
  Focused lifecycle and contractlint, live-tag compile, `go test ./...`, `go test ./... -race`, and gofmt passed; Roborev job 2331 is 2/2 with no findings.
- DONE: Perform the semantic adversarial pass for all four acceptance criteria, classify every finding by defect kind and release scope, and give a PASSED/REJECTED recommendation without inventing product or workflow obligations.
  File/prose, allow-list present/absent/outside, active/archive, repeated/out-of-order/terminal lifecycle, and baseline/candidate variants were exercised; recommendation is PASSED.
- DONE: Acceptance criterion 1 has direct scope and deletion evidence.
  Numstat/name-status and exact-name searches prove the four-file boundary, zero product/instruction drift, and removal of the tautological mutant/parser.
- DONE: Acceptance criterion 2 has load-bearing positive and negative controls.
  Focused tests passed; deleting both allowed-mod `gh pr view` tokens failed the positive control, while stubbing the shared scanner failed the planted-path negative control.
- DONE: Acceptance criterion 3 has direct on-disk mutation evidence in both runner paths.
  Post-run moves to `_archive/recorded-gate-task/index.md` made Codex and Claude fail first with the named archive diagnostic and `stat err=<nil>`, never a generic active-read error.
- DONE: Acceptance criterion 4 has no-regression evidence against the exact base.
  Executable owners and repository gates pass; the canonical writer, fixtures, held-state oracle, CLI, lifecycle skill, and default-headless test are byte-identical at base and candidate.
- FAILED: The three credentialed journey commands exit zero after exercising the new assertion.
  Codex and Claude failed later at the unchanged explicit-`state: open` count; default-headless twice failed the unchanged held-boundary oracle after binding/presenting and clearing archive absence.
- DONE: Finding classification — pre-existing canonical-open hold oracle.
  Defect kind: evidence defect; release scope: deferred risk for this cleanup, not a value regression. Trigger: the supported canonical writer encodes open by Resolution absence while the old live oracle demands forbidden explicit attempt state.
- DONE: Deferred-risk boundary and promotion condition are recorded.
  The trigger is outside this cleanup's archive/proof-deletion promise and exists identically at `4ff98d8c`; promote it to material here only if this candidate changes gate state semantics or fails to diagnose actual archival before the active read.
- DONE: Recommendation — PASSED with no material finding.
  The credentialed red remains actionable follow-up for the live harness, but assigning obsolete canonical-state semantics to this four-file cleanup would invent a value regression absent from the diff and baseline.

### Summary

Validation reproduced the deletion-first diff, all deterministic gates, Roborev, both structural mutation controls, and named archive failures through both live runners. Credentialed journeys cleared the new archive assertion and exposed a byte-identical baseline oracle that contradicts canonical-v1 serialization; it is recorded as a separately triggered evidence risk, not a candidate regression. The candidate satisfies all four acceptance criteria and is recommended PASSED.
