---
id: 77k7m0dmwm10mz6zrdq86tv3
title: Retire or re-justify the test-only gate summary reader trio
status: implementation
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
