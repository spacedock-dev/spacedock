---
title: Right-size recent proof scaffolding
status: backlog
source: "Captain review of PRs #767, #768, #776, #777, #780, and #781 on 2026-08-28."
started:
completed:
verdict:
score: 0.98
worktree:
issue:
pr:
mod-block:
milestone: 0.28.0
id: 3nm832m6pcnm8008n3wt7h9s
gates:
    version: 1
    records:
        - id: gate:3nm832m6pcnm8008n3wt7h9s:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3nm832m6pcnm8008n3wt7h9s-backlog-1
              briefing:
                id: briefing:3nm832m6pcnm8008n3wt7h9s:backlog:attempt-1:revision-1
                digest: sha256:110006f02ab3dbf144f9b908a67ff15f6eb208135613e3b6556d17ffb457c998
                room-ref: '@review/backlog/briefing-1'
---

Remove redundant tests and fake proof from the current 0.28 stack. Keep one authoritative observation for each distinct failure mode.

This task is the cleanup layer above `spacedock-ensign/flat-entity-gate-room-durability` at `37f50588aa37cfa571a88e4aa87b8f5c8f1b39e8`. The implementation worktree must start from that commit. The adjacent `retire-non-authoritative-test-oracles` task removes older dead oracles and stays separate.

## Problem

Recent fixes added more proof code than product code. Some checks inspect prose or source text. Some repeat the same invariant across package, command, fixture, and live layers. PR #781 also installs the current checkout with a stable-looking version in every common live scenario. That run does not prove the published stable package.

The result is slow CI, large fixtures, and false confidence. A validator can prove these claims with fewer committed tests and targeted manual execution.

## Required outcome

Delete or combine the redundant proof added by PRs #767, #768, #776, #777, #780, and #781. Preserve one primary proof owner for each distinct product failure mode.

Add a short workflow rule for test selection. The rule must not add a lint, gate, CI lane, required table, or recurring ceremony.

## Cleanup targets

- Replace the duplicate Pi live journey and instruction-source checks from #776 with one existing live frontdoor journey that observes the child transcript, report, and commit.
- Remove source and generated-prompt text checks that claim agent or runtime behavior.
- Keep one public `gate prepare` propagation case in #768. Leave symlink, foreign-root, and path grammar matrices at the owning resolver layer.
- Reduce #780 round tests to one flat publication boundary and one replay boundary. Keep canonical references, malformed reserved references, and frozen historical references.
- Narrow #767 terminal-delivery coverage when an existing terminal journey already owns the same failure mode.
- Share #777 setup while retaining the distinct resume aliases and one non-resume control.
- Stop #781 from presenting a copied checkout as a published stable package. Use candidate installation only for candidate behavior. Record a real stable-package installation as manual validation when that release claim matters.
- Review the older adjacent `pi_live_controls_test.go` and `pi_evidence_grade_impl_test.go` source checks. Remove them when they claim behavior from source text.

## Risk evidence

The audit already found the concrete duplication. #776 adds about 301 test lines around a behavior that the existing `TestLivePiFrontDoorSmoke` can observe. #781 copies the current checkout, stamps it with its manifest version, points installation at a temporary marketplace, and calls that local candidate stable. The behavior run is real, but the release-source claim is false.

## Out of scope

- Product behavior changes.
- New test frameworks, CI lanes, runtime controllers, or fixture protocols.
- The older dead-oracle cleanup owned by `retire-non-authoritative-test-oracles`.
- Release-process redesign.

## Expected surface and tolerance

Initial estimate net LOC change: -450, across 14 files. Ideation must refine this estimate. Tolerance: 100 net lines and 2 files.

Allowed semantic change: test selection, test setup, live-run count, and validation guidance. Command grammar, stored formats, authority, and supported runtime behavior must not change.

## Acceptance criteria

**AC-1 (VALUE): The stack keeps one primary proof for each distinct failure mode and removes at least 300 net test or harness lines.**
Verified by: compare `git diff --numstat` from the approved stack parent. Map each retained check to a distinct falsifying edit. A duplicate proof owner or a reduction below 300 net lines fails this criterion.

**AC-2: No committed test claims agent or runtime behavior from prose, generated prompt text, or Go source text.**
Verified by: inspect the changed tests and execute their retained behavioral owners. A retained assertion whose observed value comes only from the text under test fails this criterion.

**AC-3: Pi frontdoor behavior has one live journey that observes the real child transcript, ensign skill load, report, and clean commit.**
Verified by: run that targeted Pi journey once. A second journey, fixture, registry row, or CI lane for the same behavior fails this criterion.

**AC-4: Package-source claims name the real source. A copied checkout is never labeled as published stable.**
Verified by: run one candidate frontdoor install from the staged candidate and inspect its source. When published stable is relevant, manually install from the real release channel and record the observed version and behavior in validation. A temporary marketplace presented as published stable fails this criterion.

**AC-5: Resolver and round coverage remains at its owning layer, with at most one public-boundary propagation check for each invariant.**
Verified by: run focused resolver, `gate prepare`, round publication, and replay tests. Repeating one invariant across cross-product forms without a distinct failure mode fails this criterion.

**AC-6: The workflow tells authors to select one primary proof owner and to add another check only for a distinct failure mode.**
Verified by: inspect the workflow diff and one ideation report produced under it. A new lint, gate, CI lane, required table, or recurring validation step fails this criterion.

**AC-7: The full deterministic suite remains green, and required live evidence is smaller than the pre-cleanup stack.**
Verified by: run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. Run only the targeted live journeys justified by AC-3 and changed host boundaries. Record before-and-after live journey and package-install counts.

## Test plan

Reuse existing tests first. Delete or combine checks before writing any replacement.

Use deterministic Go tests for parser, command, path, and storage behavior. Use one live journey only for agent behavior that deterministic tests cannot observe. Use manual validation for published-package provenance, one-time wiring, and release-channel behavior.

Do not add test files, live scenarios, fixtures, registry rows, CI lanes, or harness code unless ideation identifies a distinct uncovered failure mode and the captain approves that expansion.

### Feedback Cycles

