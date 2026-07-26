---
title: Source builds use checkout compatibility identity, not Git-tag provenance
status: backlog
source: Captain ruling after repeated post-release auto-pre0 source-build drift, 2026-07-26
started:
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: v2183mw7c09a10pw185p33cw
---

Prevent an ordinary source build from impersonating the automatic next-minor `pre0` release merely because that tag is the nearest Git ancestor.

## Problem

The automatic post-release `vX.(Y+1).0-pre0` tag is intentionally placed on the green stable release commit. A local build path stamped `internal/cli.Version` from `git describe`, so a later source checkout with manifest `0.26.0` reported `v0.27.0-pre0-89-g...`. The First Officer's minor-version gate then treated the source binary and adjacent source skills as incompatible. An ordinary unstamped `go build` already reports the correct checkout identity, `0.26.0+dev`; the false drift comes from allowing Git provenance to masquerade as compatibility identity.

## Approved direction

Keep the automatic `pre0` release tag and strict First Officer compatibility gate. Make source versus release build intent mechanical: source builds default to the embedded checkout manifest plus `+dev`; only an explicitly marked release build may trust a linker-stamped release version. Git revision, dirty state, and nearest tag are provenance only and must not participate in the compatibility comparison. Remove the misleading source-build `git describe` stamping example and retain a plain canonical `go build` path.

Do not solve this by tolerating a one-minor mismatch, weakening the First Officer gate, removing the auto-tag, or adding compatibility behavior for unreleased source-build mistakes.

## Acceptance criteria seed

**AC-1 (VALUE) - A source build remains compatible with its adjacent checkout even when Git history contains an automatic future-minor `pre0` tag.**
Verified by: a real or fixture-backed source build whose embedded manifest is minor Y and whose Git/tag-derived candidate is minor Y+1 reports `Y+dev`, and the existing compatibility gate accepts the adjacent Y plugin. The proof must fail if the source build begins trusting the tag-derived candidate.

**AC-2 - Only the explicit release pipeline can make the binary claim a release tag version.**
Verified by: release-profile tests and release wiring that produce the exact tag version only with the explicit release marker; a version ldflag without that marker remains a source build. The proof must fail if an ambient or copied `git describe` stamp is sufficient.

**AC-3 - Git provenance cannot change compatibility decisions.**
Verified by: cases varying revision/tag/dirty provenance while holding the embedded manifest and build profile fixed produce the same compatibility identity and gate result. The proof must fail if provenance is parsed as the binary compatibility version.

**AC-4 - Source-build guidance has one canonical unstamped build path.**
Verified by: executable source-build behavior plus documentation review against the working command; no committed prose-grep may stand in for behavior. The proof must fail if the documented command again injects a Git-derived version.

## Test plan seed

Ideation must choose the smallest explicit build-profile marker and prove the riskiest case first: a source build carrying a misleading Git-derived version input still reports the embedded manifest minor. Add focused unit/behavior fixtures before implementation, retain the existing real `go build` source-build test, cover release configuration wiring, and run full/race/format checks. Because the change touches the front-door compatibility and release surfaces, include the workflow-required detached adversarial audit and map any live lane requirement from the actual diff rather than assumed relatedness.
