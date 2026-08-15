---
id: 77k7m0dmwm10mz6zrdq86tv3
title: Retire or re-justify the test-only gate summary reader trio
status: validation
source: "Gate finding from remove-gate-validate-subcommand ideation (2026-08-15): SummaryFile/SummaryFileAt/SummaryFileDiagnosticsAt dead-end in tests after the subcommand removal; captain accepted the follow-up at the gate"
started: 2026-08-15T23:04:22Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-retire-test-only-gate-summary-readers
issue:
gates:
    version: 1
    records:
        - id: gate:77k7m0dmwm10mz6zrdq86tv3:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:77k7m0dmwm10mz6zrdq86tv3-backlog-1
              briefing:
                id: briefing:77k7m0dmwm10mz6zrdq86tv3:backlog:attempt-1:revision-1
                digest: sha256:8bd171d5a202a453728de4bce4cf9760d3044a3a17e4cf28615810d5d311c20d
                request-digest: sha256:db59d64f51fc38a4a61f0950e227498d48e75206c0eedae708aa20dc832c0ece
                room-ref: ./retire-test-only-gate-summary-readers/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:77k7m0dmwm10mz6zrdq86tv3:backlog:1
                briefing: briefing:77k7m0dmwm10mz6zrdq86tv3:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T21:24:39.885375Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: dispatch all five onto the stack tip'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:77k7m0dmwm10mz6zrdq86tv3:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:77k7m0dmwm10mz6zrdq86tv3-ideation-1
              briefing:
                id: briefing:77k7m0dmwm10mz6zrdq86tv3:ideation:attempt-1:revision-1
                digest: sha256:de0d38445f0639c82598673eb74e46ec931f10b0fb0f1901b16780d0ef066fff
                request-digest: sha256:186c33d4c6993b2f0d1f041170ab7e53824eb6c84d6532f5b595dc98952e9b08
                room-ref: ./retire-test-only-gate-summary-readers/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:77k7m0dmwm10mz6zrdq86tv3:ideation:1
                briefing: briefing:77k7m0dmwm10mz6zrdq86tv3:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T23:04:05.521463Z"
                decision: approve
                reason: 'Captain batch approval 2026-08-15 (approve all): into implementation as stack layers'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:77k7m0dmwm10mz6zrdq86tv3:validation
          stage: validation
          attempts:
            - id: gate-attempt:77k7m0dmwm10mz6zrdq86tv3-validation-1
              briefing:
                id: briefing:77k7m0dmwm10mz6zrdq86tv3:validation:attempt-1:revision-1
                digest: sha256:4ce9ac4405b24316252530a17b81ffc92563982ec32b842e6a3dc3407e2bc380
                request-digest: sha256:07833adf63c5bcaa787e1b1e3c17f9a40a49748b4fa1b6c8af0d6cc16a6d5275
                room-ref: ./retire-test-only-gate-summary-readers/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:77k7m0dmwm10mz6zrdq86tv3:validation:1
                briefing: briefing:77k7m0dmwm10mz6zrdq86tv3:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T23:54:37.554977Z"
                decision: approve
                reason: 'Captain 2026-08-15 (push them for the stack): validation PASSED approved; lands as a stack layer'
              application:
                target-stage: done
                state: pending
pr: "#714"
---

After remove-gate-validate-subcommand lands, the `SummaryFile` / `SummaryFileAt` / `SummaryFileDiagnosticsAt` trio has no production caller, and the warnings slice `SummaryFileDiagnosticsAt` returns is observable by no caller. Retire the trio, or record the one concrete justification that keeps it. Removing it requires re-pointing two internal/gates tests.

## Problem

Designed against the stack tip `stack27/11-self-contained-contracts` (abc1a569f), not `main`.

On `main` the trio is still live: `internal/cli/cli.go:337` calls `SummaryFileDiagnosticsAt` from the `gate validate` branch. `stack27/07-remove-gate-validate` deletes that branch, so from layer 07 onward the trio's only remaining callers are two tests inside `internal/gates`. At the tip:

- `SummaryFile` is called once, by `internal/gates/application_test.go:365`.
- `SummaryFileAt` is called once, by `internal/gates/prepare_test.go:956`, and once by `SummaryFile`.
- `SummaryFileDiagnosticsAt` is called once, by `SummaryFileAt`.

So the three functions form a closed chain whose only entry points are two test call sites. The `[]Warning` slice `SummaryFileDiagnosticsAt` returns is observable by nobody at all: its only renderer, `FormatWarning`, was removed by `stack27/08-trim-dead-gate-model`, and `status --validate` derives the same warning class independently through its own `ReadDiagnostics` call in `gateValidationDiagnostics` (`internal/status/validate.go:230`). The tip's own docstring already records the dead end: "No command prints that warnings slice today."

There are zero non-Go references to the trio anywhere in the repo, so no doc, contract, skill, or lint pin depends on the names.

## Proposed approach

Delete all three functions and re-point the two tests at readers production actually calls. No new mechanism is introduced — this is deletion plus two call-site substitutions — so there is no enabling mechanism to justify.

**`internal/gates/io.go` — delete the trio (-28).**

**Test A, `application_test.go` `TestResolutionSummaryDoesNotHashBriefing`.** Replace `SummaryFile(entity)` with `Read` + `entityStatus` + `CurrentSummary`, which is the exact composition `SummaryFileDiagnosticsAt` performed. The one step not carried over is `validateRetainedAuthority`, and a runtime probe proved it never engages for this fixture: `applicationWorkflow` writes a briefing mapping of `{id, digest, room-ref}` with no `request-digest`, so `validateRetainedAuthorityExcept` hits its `RequestDigest == ""` continue for every attempt. A `panic()` planted immediately after that continue did not fire and the test still passed. The substitution is therefore behavior-preserving, not coverage-losing.

```go
	summaryDoc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	stage, err := entityStatus(entity)
	if err != nil {
		t.Fatal(err)
	}
	if summary := CurrentSummary(summaryDoc, stage); summary.Decision != "approve" {
		t.Fatalf("resolution-only summary = %#v", summary)
	}
```

**Test B, `prepare_test.go` `TestPreparedAuthorityIsRecomputedDuringReadOnlyValidation`.** Replace `SummaryFileAt(entity, workflow)` with `EligibilityFileAt(entity, workflow)`. This is strictly stronger than the original: `EligibilityFileAt` runs the same `validateRetainedAuthority` on a read-only path *and* is production-called at `internal/status/merge.go:566`, so the claim "read-only validation recomputes prepared authority" moves from a test-only reader onto one the product exercises. Unlike Test A's fixture, this fixture goes through `Prepare`, so `request-digest` is set and the authority check genuinely runs.

```go
			if _, err := EligibilityFileAt(entity, workflow); err == nil {
				t.Fatal("read-only validation accepted drifted prepared authority")
			}
```

### Keep boundaries

Every adjacent surface this change must not touch, and why it stays:

- `Warning` struct and `ReadDiagnostics` — live; `status --validate` consumes them via `gateValidationDiagnostics`.
- `validateRetainedAuthority` / `validateRetainedAuthorityExcept` — live; 7 production call sites across `application.go`, `delivery.go`, `operation.go`, `prepare.go`.
- `CurrentSummary` — live; `RecordSemanticSummary`, `Withdraw`, and `internal/livescenario/ac2_reanchor.go:45`.
- `EligibilityFile` / `EligibilityFileAt` — live; `internal/status/merge.go:566`.
- `nearestWorkflowDir` (5 remaining callers) and `entityStatus` (9) — not orphaned by dropping the trio; verified after deletion.
- `newSummaryFile` in `cmd/spacedock-release/e2e_gate_test.go` — an unrelated local release-gate test helper that only shares the name stem. Out of scope; do not touch.

### Overlap coordination

- `stack27/07-remove-gate-validate` removes the trio's last production caller. This entity must land on top of it; against `main` alone the deletion would not compile.
- `stack27/08-trim-dead-gate-model` already removed `FormatWarning`, the `Diagnostic` alias, and `ReadWithWarnings`. Those are not this entity's surface and must not be re-claimed.
- No file overlap with layers 09, 10, or 11.

## Out of scope

`ReadDiagnostics` and its `status --validate` consumer. The read tolerance.

## Expected surface and tolerance

Measured from the spike against `stack27/11-self-contained-contracts` (abc1a569f): **+11 / -32, net -21, across 3 files.**

| File | Delta |
| --- | --- |
| `internal/gates/io.go` | -28 |
| `internal/gates/application_test.go` | +10 / -3 |
| `internal/gates/prepare_test.go` | +1 / -1 |

Tolerance: ±10 lines and ±1 file. Exceeding either means the deletion pulled in a keep-boundary and the implementation should stop and report rather than absorb it.

**Semantic changes: none.** No command grammar, stored format, authority rule, or runtime behavior changes. The deleted functions sit behind no CLI surface, and the repo has zero non-Go references to them, so no documentation diff is required for this entity.

**Spike: done.** The full change was exercised in a throwaway worktree at the stack tip before this design was written; results are recorded under Test plan below. The worktree has been removed. Nothing in this design rests on an unexercised mechanism.

## Acceptance criteria

**AC-1 - No production code path retains a reader whose output no caller can observe.**
Verified by: the tree builds and the suite passes with `SummaryFile`, `SummaryFileAt`, and `SummaryFileDiagnosticsAt` absent from the source. The Go compiler is the caller-oracle here, not a substring search: any surviving caller of a deleted function fails the build, so a green `go build ./...` plus a green suite is direct evidence that nothing referenced them.

**AC-2 - The change removes more lines than it adds.**
Verified by: `git diff --numstat` against the stack parent `stack27/11-self-contained-contracts` sums to a negative net. Baseline expectation -21; this number can move the wrong way if the implementation adds scaffolding instead of deleting, which is the failure this AC exists to catch.

**AC-3 - The suite stays green.**
Verified by: `go test ./...` and `go test ./... -race` pass.

**AC-4 - The two re-pointed tests still fail when the behavior they claim is broken.**
Verified by: each test is run against a named breaking change and observed to fail (see Test plan). A re-point that quietly turned either test into a tautology is the specific risk this entity carries, since both tests are being moved off the reader they were written against.

## Test plan

No new tests. The change is a deletion plus two call-site re-points, so the existing suite is the net and the work is proving the re-points stay honest. All four results below were obtained in the spike at tip abc1a569f.

**Existing suite.** `go build ./...` clean, `gofmt -l ./cmd ./internal` clean, and `go test ./...` exited 0 with no failing package. Implementation must additionally run `go test ./... -race` per AC-3.

**Test A falsifiability.** `TestResolutionSummaryDoesNotHashBriefing` asserts the summary projection reports `Decision: "approve"` after `briefing.json` is deleted. Falsifying change: make the decision projection in `CurrentSummary` conditional on the briefing being readable. Observed under that change — `resolution-only summary = ... Decision:""` and the test fails. Unmodified, it passes.

**Test B falsifiability.** `TestPreparedAuthorityIsRecomputedDuringReadOnlyValidation` asserts a read-only read refuses drifted prepared authority across three mutations (request drift, briefing summary drift, missing local source object). Falsifying change: delete the `validateRetainedAuthority` call from `eligibilityFileAt`. Observed under that change — all three subtests fail with "read-only validation accepted drifted prepared authority". Unmodified, all three pass.

**No-op confirmation for Test A's dropped check.** A `panic()` planted immediately after the `RequestDigest == ""` continue in `validateRetainedAuthorityExcept` did not fire while Test A ran, and Test A passed. This is what licenses dropping `validateRetainedAuthority` from Test A's path: it was never executing there.

**Orphan check.** After the deletion, `nearestWorkflowDir` retains 5 production callers and `entityStatus` retains 9, so neither helper is left dangling.

## Stage Report: ideation

- DONE: Retire-or-justify the SummaryFile trio with the two internal/gates tests re-pointed; deletion preferred per the gated finding
  Deletion selected and fully exercised; no justification survives — at tip abc1a569f the trio is a closed 3-function chain whose only entry points are `application_test.go:365` and `prepare_test.go:956`, and its `[]Warning` return has had no renderer since `stack27/08` removed `FormatWarning`.
- DONE: Value AC: no production code retains a reader whose output no caller can observe; negative LOC
  AC-1 re-proofed to use the compiler as caller-oracle (green build + suite with the functions absent) rather than a grep; AC-2 baselined at -21 net against the stack parent, measured by `git diff --numstat` in the spike.
- DONE: Design against the stack tip
  Designed against `stack27/11-self-contained-contracts` (abc1a569f), not `main` — on `main` the trio still has a production caller at `internal/cli/cli.go:337` and the deletion would not compile.

### Summary

Confirmed the removal scope against the stack tip rather than HEAD, which changed the answer: the trio is live on `main` and dead only from `stack27/07-remove-gate-validate` onward, so this entity must land as a stack layer. Spiked the entire change in a throwaway worktree — deletion plus both re-points — and got a clean build, clean gofmt, and `go test ./...` exit 0 at +11/-32 (net -21) across 3 files. Because both tests were being moved off the reader they were written against, I separately proved each re-point still fails under a named breaking change, and proved by runtime probe that the `validateRetainedAuthority` step dropped from Test A never executed in that fixture; those results are recorded in the test plan and added as AC-4. Test B's re-point onto `EligibilityFileAt` is a net strengthening, since that reader is production-called from `internal/status/merge.go:566` while the old one was not called at all.

## Stage Report: implementation

- DONE: Execute the gated design exactly: delete the trio, both test re-points as designed, each proven to fail under its named breaking change
  `bc776e8e8` — the three functions deleted from `internal/gates/io.go`, Test A composed from `Read` + `entityStatus` + `CurrentSummary`, Test B moved to `EligibilityFileAt`; the design's code blocks landed verbatim, no deviation.
  Test A, `TestResolutionSummaryDoesNotHashBriefing`, asserts the summary projection still reports `Decision: "approve"` after `briefing.json` is deleted. Named breaking change: make the decision projection conditional on the briefing being readable (probe in `Read` cleared the attempt's resolution when `inspectReviewedInput` was not current). Observed FAIL — `resolution-only summary = ... Decision:""` — and PASS again once reverted.
  Test B, `TestPreparedAuthorityIsRecomputedDuringReadOnlyValidation`, asserts a read-only read refuses drifted prepared authority across request drift, briefing summary drift, and a missing local source object. Named breaking change: delete the `validateRetainedAuthority` call from `eligibilityFileAt`. Observed FAIL on all three subtests — "read-only validation accepted drifted prepared authority" at `prepare_test.go:957` — and PASS again once reverted.
  Re-confirmed the no-op that licenses Test A's dropped step: with `panic()` planted immediately past the `RequestDigest == ""` continue in `validateRetainedAuthorityExcept`, the *original* Test A at base `ab112b6d3` still passed and the panic never fired, so `validateRetainedAuthority` was never executing for that fixture. Run in a throwaway worktree namespaced to this entity, since the original call site no longer exists on this branch; worktree removed.
- DONE: Value AC: no production reader with unobservable output survives; net negative (spiked -21 over 3 files)
  AC-1: `go build ./...` clean with `SummaryFile`, `SummaryFileAt`, and `SummaryFileDiagnosticsAt` absent — the compiler is the caller-oracle, and a repo-wide search now returns zero references of any kind. The unobservable `[]Warning` return died with `SummaryFileDiagnosticsAt`; `status --validate` still derives that class through its own `ReadDiagnostics` call at `internal/status/validate.go:230`.
  AC-2: `git diff --numstat ab112b6d3..HEAD` = +11 / -32, **net -21 across 3 files** — exactly the spiked baseline, zero deviation against the ±10 line / ±1 file tolerance.
  AC-3: `go test ./...` exit 0 (20 packages ok) and `go test ./... -race` exit 0 (20 packages ok, no data race). `gofmt -l ./cmd ./internal` clean; `go vet ./internal/gates/` clean.
  Keep boundaries all held: the diff touches only the 3 designed files. `Warning`/`ReadDiagnostics`, `validateRetainedAuthority(Except)`, `CurrentSummary`, and `EligibilityFile`/`EligibilityFileAt` are untouched and still production-called; `newSummaryFile` in `cmd/spacedock-release/e2e_gate_test.go` was not touched (0 lines of diff in that tree). Orphan check re-run after deletion: `nearestWorkflowDir` retains 5 callers, `entityStatus` retains 9.
- DONE: Base is the stack top ab112b6d3; land as the next available layer per the ratified Stacked mode
  Branch `spacedock-ensign/retire-test-only-gate-summary-readers` is a single commit `bc776e8e8` whose parent is `ab112b6d3`, so it is the next layer on the tip with no rebase owed. The trio's last production caller (`internal/cli/cli.go` `gate validate`) is already gone below this base, which is what makes the deletion compile — against `main` alone it would not. PR creation and `gh stack link` are the FO's ceremony, not done here.

### Summary

The gated design landed exactly as written in `bc776e8e8`: the three-function chain deleted, Test A re-composed from the readers `SummaryFileDiagnosticsAt` itself used, Test B moved onto `EligibilityFileAt`, at +11/-32 (net -21) across the 3 predicted files — the spiked baseline hit on the nose with no tolerance consumed. Both suites are green, with and without `-race`.

The work that mattered was AC-4 rather than the deletion. Each re-pointed test was run against its named breaking change and observed to fail with the predicted message, then observed to pass again on revert, so neither re-point is a tautology; and the one step the substitution drops was re-proven inert by a panic probe against the original test at the base commit, not merely asserted from the ideation spike. Nothing here is left for validation to finish.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against worktree commit bc776e8e8, never by reading the report: trio gone (grep clean over internal/ cmd/ skills/), 3 files +11/-32 net -21, both re-pointed tests present and green
  AC-1: word-boundary grep for the three names over internal/ cmd/ skills/ returns nothing, and a repo-wide sweep (docs and contracts included) is also clean; only the declared keep-boundary `newSummaryFile` survives, untouched. `go build ./...` green with the functions absent — compiler as caller-oracle. AC-2: `git diff --numstat ab112b6d3..bc776e8e8` = +11/-32 net -21 across exactly the 3 designed files, zero tolerance consumed; the io.go hunk deletes the three functions and nothing else. Both re-pointed tests present at bc776e8e8 and green in a detached audit checkout, matching the gated design's code blocks verbatim.
- DONE: Each re-point proven falsifiable: apply the named breaking change per the gated design, expect red, revert; the EligibilityFileAt re-point exercises a production-called reader
  All probes run in a throwaway detached checkout at bc776e8e8, never the implementation worktree. Test A under "decision projection conditional on briefing readability" (attempt resolution cleared in Read when inspectReviewedInput is not current): FAIL with `resolution-only summary = ... Decision:""` at application_test.go:374; green again on revert. Test B under "delete validateRetainedAuthority from eligibilityFileAt": all three subtests FAIL "read-only validation accepted drifted prepared authority" at prepare_test.go:957; green on revert. The no-op license for Test A's dropped step was reproduced, not trusted: panic planted past the `RequestDigest == ""` guard at base ab112b6d3, original SummaryFile-based Test A passed with the panic silent, and as a positive control the same probe fires once under Test B's digest-bearing fixture — the probe placement is live, so its silence is meaningful. `EligibilityFileAt` is production-called at internal/status/merge.go:566.
- DONE: internal/gates green plain and -race; verdict PASSED or REJECTED with per-AC citations; the branch bases on ab112b6d3 whose 42c content is under rework - your verdict covers this diff, the restack is FO-owned
  internal/gates: plain pass in the audit checkout and a fresh (uncached) `-race` pass, exit 0. AC-3 full form: `go test ./...` exit 0 (20 packages ok) and `go test ./... -race` exit 0 in the candidate worktree; `gofmt -l ./cmd ./internal` clean. Keep boundaries verified live by direct grep: ReadDiagnostics consumed at internal/status/validate.go:230, validateRetainedAuthority at 6 production sites, CurrentSummary at 3, nearestWorkflowDir 5 callers, entityStatus 8 production callers. Verdict: PASSED — AC-1, AC-2, AC-3, AC-4 all verified by commands run here. No material findings, no deferred risks, no polish findings. The verdict covers the ab112b6d3..bc776e8e8 diff only; the restack over the 42c rework is FO-owned.

### Summary

PASSED. Every AC was re-exercised from scratch against bc776e8e8: the trio is gone with the compiler and a repo-wide sweep as independent oracles, the diff is exactly the spiked +11/-32 over 3 files, and both suites are green plain and under -race. The falsifiability work — the entity's real risk — held up under independent reproduction: both re-points go red under their design-named breaking changes with the exact predicted messages and green on revert, and the panic probe licensing Test A's dropped authority step was re-run at the base commit with a positive control proving the probe itself fires when the check engages. No findings of any class; recommend the gate accept this layer as-is, restack owed to the FO.
