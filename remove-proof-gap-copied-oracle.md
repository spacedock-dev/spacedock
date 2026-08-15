---
id: pph11xwsa2s3z73ts983wd0k
title: Remove the copied proof-gap oracle from the reconciliation lint
status: ideation
source: "Captain directive 2026-08-15 after ffacf584a reintroduced the pattern zvk9 deleted; zv ensign composition report flagged it"
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:pph11xwsa2s3z73ts983wd0k:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:pph11xwsa2s3z73ts983wd0k-backlog-1
              briefing:
                id: briefing:pph11xwsa2s3z73ts983wd0k:backlog:attempt-1:revision-1
                digest: sha256:f3a9c595b6653640f2bfa8e58198d69c5b65e5a73ddd9456b5d3fef6ccb89f8c
                request-digest: sha256:55cd1712c5b2a8150909aad6747819de2dcefeb3886ee9b2c7d17c5ddf231828
                room-ref: ./remove-proof-gap-copied-oracle/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pph11xwsa2s3z73ts983wd0k:backlog:1
                briefing: briefing:pph11xwsa2s3z73ts983wd0k:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T21:24:34.622076Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: dispatch all five onto the stack tip'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:pph11xwsa2s3z73ts983wd0k:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:pph11xwsa2s3z73ts983wd0k-ideation-1
              briefing:
                id: briefing:pph11xwsa2s3z73ts983wd0k:ideation:attempt-1:revision-1
                digest: sha256:d7125bc2bd6ee69f571c621e142be52a72b9fe9bf9fbf2efc919d10ea286960e
                request-digest: sha256:e7619ef9893b8bc4fff04e40204dda701f6eb84706fbb89a6dad2d4d98c54c42
                room-ref: ./remove-proof-gap-copied-oracle/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pph11xwsa2s3z73ts983wd0k:ideation:1
                briefing: briefing:pph11xwsa2s3z73ts983wd0k:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T23:04:10.82064Z"
                decision: approve
                reason: 'Captain batch approval 2026-08-15 (approve all): into implementation as stack layers'
              application:
                target-stage: implementation
                state: pending
---

ffacf584a added `wantProofGaps`, a hand-written map of proof-gap bindings, DeepEqual-checked against `readRuntimeProofGaps` - an AST reader that derives the same data from the live test sources three lines away. The copy's only failure mode is "you forgot to update the copy": it catches no real defect and reintroduces the two-file lockstep tax for XFAILs on live-proof tests, one abstraction above the wantGaps map the workflow just removed for the same reason.

Keep-boundary: `readRuntimeProofGaps` itself, the `parseLiveGap` shape validation it feeds, and the owner-join extension in TestRuntimeLiveTODOOwnersAreActive are load-bearing and stay.

## Problem

Confirmed against the **stack tip** `origin/stack27/11-self-contained-contracts` (c908ac657, file blob 02f347acb), not against `main`. All line numbers below are tip line numbers.

`internal/contractlint/live_registry_reconciliation_test.go:51-57` derives the live-proof gap bindings by AST walk and then compares them to a hand-written copy of the same data written three lines later:

```go
	proofGaps := readRuntimeProofGaps(t, repo, targets)
	wantProofGaps := map[string][]liveGapRow{
		"TestLiveBreakGlassShimRecovery": {{"xfail", "claude-sonnet", "060xp69y61yhrww23g3wvwqy"}},
	}
	if !reflect.DeepEqual(proofGaps, wantProofGaps) {
		t.Errorf("runtime proof gaps = %#v, want %#v", proofGaps, wantProofGaps)
	}
```

The copy's only failure mode is "you forgot to update the copy". Measured, not asserted: transferring the XFAIL owner in `internal/ensigncycle/dispatch_recovery_live_test.go:108` and changing no other file fails at `live_registry_reconciliation_test.go:56` with a DeepEqual mismatch — the derived value is correct and the frozen copy is stale. No defect is caught; a second file is taxed.

This is the same pattern one abstraction up from the `wantGaps` map that cb48188ca (`remove-redundant-lint-mirrors`, zvk9) deleted for exactly this reason. That deletion is on the stack, not on `main`; ffacf584a landed on `main` afterwards and reintroduced the pattern for runtime proofs.

It also re-falsifies the registry doc. `docs/runtime-live-ci-registry.md:350-351` states the reconciliation "derives findings from current source; it does not use a copied gap oracle." zvk9 made that sentence true for common journeys; ffacf584a made it false again for runtime proofs.

**Why the tip is the only correct base.** On `main` the file still carries `wantGaps`, so `reflect` survives this removal. On the tip zvk9 already deleted `wantGaps`, so the copied proof-gap oracle is the file's **only** remaining `reflect` user. A patch authored against `main` — which would correctly retain the import — fails on the tip: `vet: internal/contractlint/live_registry_reconciliation_test.go:11:2: "reflect" imported and not used` (measured). The base is load-bearing, not bookkeeping.

## Proposed approach

One file, one deletion, no new mechanism.

`internal/contractlint/live_registry_reconciliation_test.go`:

1. Delete the `wantProofGaps` literal and its `reflect.DeepEqual` check (tip lines 52-57).
2. Keep the derivation call, demoted to a bare validating statement.
3. Delete the then-unused `"reflect"` import (tip line 11).

Before (tip lines 51-57) is the block quoted under Problem. After:

```go
	// Fatals on a malformed live-proof gap binding; the owner join re-derives these
	// under SPACEDOCK_LIVE_STATE_DIR.
	readRuntimeProofGaps(t, repo, targets)
```

**Why the call survives rather than going with the map.** `readRuntimeProofGaps` has exactly two call sites. The other is in `TestRuntimeLiveTODOOwnersAreActive`, which `t.Skip`s unless `SPACEDOCK_LIVE_STATE_DIR` is set — and nothing under `.github/` sets it; the only occurrences outside the test are `docs/runtime-live-ci.md:37` and a roadmap doc, both documenting a local invocation. So the reconciliation call site is the only thing that runs the live-proof shape validation in CI or in a plain `go test ./...`. Deleting it outright would leave the keep-boundary's `parseLiveGap` validation present but unreachable — the broken value chain this program exists to remove.

Proven by falsifier, both directions (see Spike record): with the bare call retained, a malformed `liveXFail("claude-sonnet", "bad-owner")` in the live-proof source fails at `live_registry_reconciliation_test.go:174` with `malformed liveXFail binding`; with the call deleted outright, the identical mutation **passes**. The surviving line is load-bearing, not decoration.

**Alternatives considered.**

1. Delete the call along with the map (the simplest thing). Insufficient — falsified above: it takes a named keep-boundary item dead in every automated run.
2. Replace the copy with a derived-but-real assertion, e.g. "no live proof repeats a gap target", mirroring the duplicate-target check the common-journey loop already applies. Rejected as scope creep on a removal entity: no such defect has occurred, and the bare validating call already discharges the keep-boundary. Recorded so the gate can overrule rather than discover it.
3. Keep the map but generate it. That is what `readRuntimeProofGaps` already is; the map adds nothing.

No new mechanism is introduced, so the "name the value AC, the simplest alternative, and why it is insufficient" clause applies only to the retained line, covered above.

**Documentation:** none required. Test-only; no CLI output, command grammar, stored format, authority, or host integration changes. As with zvk9, the deletion makes `docs/runtime-live-ci-registry.md:350-351` true again rather than needing an edit.

**Spike:** done during ideation, riskiest path first — the value probe and the keep-boundary falsifier were both exercised before this design was written. Results under Spike record.

## Out of scope

The runner-file liveXFail vocabulary and the registry doc.

**Residual, named honestly.** ffacf584a also added a prose "Expected failure" line to `docs/runtime-live-ci-registry.md:290-291` naming the same owner `060xp69y61yhrww23g3wvwqy`. That is a third copy of the binding, and nothing parses it — `git grep "Expected failure" -- '*.go'` at the tip returns nothing. It stays, because the registry doc is out of scope here. Flagged for a follow-up entity rather than silently absorbed.

## Expected surface and tolerance

Measured in the spike, not estimated: **1 file, -8 / +3, net -5 lines** (631 → 626).

Tolerance ±5 lines, ±0 files. The band absorbs rewording of the retained two-line comment at review; it does not absorb a second file.

**Base guard.** The measurement holds against the stack tip. Based on `main` instead, the same change measures -7 / +3 (net -4) and must *not* drop the `"reflect"` import. An implementation that reports the main-shaped numbers has used the wrong base and should be rejected at validation.

**Semantic changes: none** to command grammar, stored formats, authority, or runtime behavior. The one semantic that moves is deliberate and is the point of the entity: changing or retiring an XFAIL on a live-proof test stops being a two-file edit and becomes a one-file edit.

## Acceptance criteria

**AC-1 (value) - An XFAIL change on a live-proof test is a one-file edit that leaves contractlint green.**
Measured against a baseline that moves the wrong way: transferring the `TestLiveBreakGlassShimRecovery` XFAIL owner in `internal/ensigncycle/dispatch_recovery_live_test.go` and changing no other file leaves `go test ./internal/contractlint/` passing. The identical one-file edit fails at the tip at `live_registry_reconciliation_test.go:56`. Verified by probing both sides and reverting the probe. An incomplete removal keeps it failing.

**AC-2 (value) - A malformed live-proof gap binding still fails the default test run.**
Falsifiable by mutation, not by grep: `liveXFail("claude-sonnet", "bad-owner")` in the live-proof source fails `go test ./internal/contractlint/` with `malformed liveXFail binding`. The keep-boundary is measured where it is actually consumed, not asserted by the symbols still being present.

**AC-3 - The tree gets smaller and the keep-boundary survives.**
Line delta against the stack tip is negative and within the tolerance band above (measured -5, 1 file), and `readRuntimeProofGaps`, the `parseLiveGap` shape validation, and the owner-join extension in `TestRuntimeLiveTODOOwnersAreActive` are all still present and still reachable.

**AC-4 - The suite is no worse than its baseline.**
`go test ./internal/contractlint/` plain and `-race` pass, `gofmt -l` is empty and `go vet` is clean for the package, and `go test ./...` reports no failure that is not also present at the tip on the same machine.

## Test plan

No new tests. The change is a deletion, and its proof is the two probes above rather than an added assertion — adding a test to watch a deleted copy would rebuild the thing being removed.

- **AC-1**, both sides: apply the one-file owner transfer at the tip (expect FAIL at `:56`), apply the deletion, re-apply the same one-file edit (expect PASS), revert.
- **AC-2**, both variants: with the deletion applied, mutate the live-proof binding to a malformed owner (expect FAIL at `:174`); then delete the retained derivation call as well and re-run the same mutation (expect PASS, which is the failure the retained line prevents). Restore the retained line.
- **AC-3**: `diff` against the tip blob for the line delta; the keep-boundary symbols are covered by AC-2 firing, not by grep.
- **AC-4**: `gofmt -l ./internal`, `go vet ./internal/contractlint/`, `go test ./internal/contractlint/` plain and `-race`, then `go test ./...`.

Cost: minutes, no fixture and no live workflow test. contractlint reads the live sources by AST and never executes them, so no credential, runtime host, or `-tags live` build is involved — the probes are ordinary `go test` runs.

## Spike record

Run during ideation against tip c908ac657, extracted with `git archive` into a scratch directory. (A `git worktree` was used first and was destroyed mid-run by a sibling agent sharing this session's job directory; the archive extraction is immune to that, and every result below is from the clean re-run.)

| Step | Expected | Observed |
| --- | --- | --- |
| Baseline at tip | pass | `ok` |
| One-file XFAIL owner transfer, oracle intact | fail | `FAIL live_registry_reconciliation_test.go:56` DeepEqual mismatch |
| Deletion applied, clean tree | pass | `gofmt` empty, `go vet` clean, `ok` |
| Same one-file edit, deletion applied | pass | `ok` — **AC-1** |
| Malformed binding, retained call present | fail | `FAIL :174 TestLiveBreakGlassShimRecovery: malformed liveXFail binding` — **AC-2** |
| Malformed binding, retained call deleted too | pass | `ok` — falsifier: proves the retained line is load-bearing |
| `-race` on the patched tree | pass | `ok` (20.7s) |
| `go test ./...` on the patched tree | pass | exit 0, 20 packages ok, 0 FAIL |
| Main-shaped patch (import retained) on the tip | fail | `vet: :11:2: "reflect" imported and not used` |

The full suite is green on the patched tree, so AC-4 can be stated absolutely rather than against a degraded baseline — unlike zvk9, which had to carve out a local sandbox failure.

Nothing in the design now rests on an unverified mechanism: the value claim, the keep-boundary claim, and the base-selection claim were each exercised and observed.

## Stage Report: ideation

- DONE: One-file deletion design: wantProofGaps and its DeepEqual go, readRuntimeProofGaps, parseLiveGap validation, and the owner-join stay
  Proposed approach specifies the exact before/after for `internal/contractlint/live_registry_reconciliation_test.go` (tip lines 11 and 51-57); the derivation call is retained as a bare validating statement because it is the only automated-run consumer of the shape validation.
- DONE: Value AC: an XFAIL change on a live-proof test is a one-file edit leaving contractlint green, probe both sides
  AC-1; probed both sides at tip c908ac657 — the one-file owner transfer FAILS at `:56` with the oracle intact and PASSES after the deletion. Probe reverted.
- DONE: Design against the stack tip
  Based on `origin/stack27/11-self-contained-contracts` (c908ac657, blob 02f347acb), where cb48188ca already removed `wantGaps`. Proven load-bearing: the main-shaped patch fails on the tip with `"reflect" imported and not used`.

### Summary

The removal is a single-file deletion of the `wantProofGaps` literal, its `reflect.DeepEqual` check, and the then-unused `"reflect"` import: measured -8/+3, net -5 lines, no second file. The one non-obvious decision is keeping `readRuntimeProofGaps(t, repo, targets)` as a bare validating call — its only other call site skips unless `SPACEDOCK_LIVE_STATE_DIR` is set, which no workflow under `.github/` does, so deleting it would take the declared `parseLiveGap` keep-boundary dead in CI. A falsifier probe proves that line is load-bearing: a malformed `liveXFail` binding fails at `:174` with the line present and passes with it removed. Two items for the gate: the base must be the stack tip, not `main` (the patches differ and the main-shaped one does not compile on the tip), and the registry doc's prose "Expected failure" line is a third, unparsed copy of the same binding that this entity deliberately leaves alone as out of scope.
