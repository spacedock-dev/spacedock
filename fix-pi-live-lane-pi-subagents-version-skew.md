---
id: 2686bggef0qz5hsrft2aks0t
title: Restore pi-live lane green by resolving the pi-subagents/pi-coding-agent version skew
status: validation
source: c6 validation cycle-7 live-lane triage, 2026-07-20
started: 2026-07-20T13:56:42Z
completed: 2026-07-20T14:41:04Z
verdict: PASSED
score: 1
worktree: .worktrees/spacedock-ensign-fix-pi-live-lane-pi-subagents-version-skew
issue:
mod-block: merge:pr-merge
---

The pi-live CI lane is deterministically red for every commit, including main, due to an upstream extension version skew. Until it is fixed, no PR can satisfy the every-host-lane merge requirement on pi-live.

## Problem

`TestLivePiFrontDoorSmoke` fails at Pi extension load: `pi-subagents@0.35.1` (published 2026-07-18) requires the `@earendil-works/pi-ai` `/compat` export path, which `pi-coding-agent@0.74.2`'s vendored copy lacks. Evidence: identical failure on main's run 29691271679 (2026-07-19, pre-dating unrelated PRs' commits) and on PR #531's run 29709011051; last green pi-live run (2026-07-17) resolved `pi-subagents@0.34.0`. npm still serves 0.35.1, so reruns red identically — this is a deterministic break, not a flake, and re-run-to-green policy cannot apply.

## Proposed approach

Bump only the pi-live job to Node 24, install exact pinned Pi substrate tarballs (`@earendil-works/pi-coding-agent@0.80.10`, `pi-subagents@0.35.1`, `pi-intercom@0.6.0`) after checking their recorded npm `dist.integrity` sha512 values, and add a fast compatibility guard for the pi-ai `/compat` export required by pi-subagents 0.35.x.

## Out of scope

Pi runtime behavior changes and unrelated live lanes' Node versions.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - The pi-live lane is green on main.**
Verified by: a main-branch pi-live run passing after the change — an independent external result that can move the wrong way.

**AC-2 - The version constraint is explicit and evergreen, not incidental.**
Verified by: the lane setup declares the pinned/bumped versions with a comment stating the compatibility requirement, so the next skew is diagnosable from the file.

## Test plan

- `go test ./...`
- `go test ./... -race`
- Local npm smoke of the pinned tarball/integrity flow and pi-ai `/compat` guard using temporary install prefixes.
- Primary live evidence remains the PR pi-live run, then a main pi-live run after merge.

## Stage Report

Stage: implementation

### Summary

Implemented the pi-live workflow fix: Node 24 for the Pi lane, exact pinned Pi substrate versions with recorded sha512 integrity checks before install, installed tarballs instead of bare latest selectors, and a deterministic compatibility guard that fails fast if pi-subagents cannot use the pi-ai `/compat` export. Updated release workflow guard tests to enforce the new pinned/integrity-based install contract.

### Dispatch checklist

- DONE — Engines-safe Node runtime for the pinned pi-coding-agent: the pi-live job now uses `actions/setup-node` with `node-version: "24"`, satisfying `@earendil-works/pi-coding-agent@0.80.10`'s `>=22.19.0` engine requirement and preventing Node 20 legacy fallback.
- DONE — Exact-version pins with sha512 `dist.integrity` verification, no bare-latest selectors: the install step pins `@earendil-works/pi-coding-agent@0.80.10`, `pi-subagents@0.35.1`, and `pi-intercom@0.6.0`, records their sha512 integrities, verifies `npm pack` output before install, and installs only the verified tarballs.
- DONE — Deterministic fail-fast compat guard for the pi-ai `/compat` export: the new guard checks Node version, pinned package versions, the `@earendil-works/pi-ai` `./compat` export declaration/file, and dynamically imports the compat module before running the live smoke.

### Validation

- DONE — `gofmt -w ./cmd ./internal`: completed with no formatting diff outside the intended Go test updates.
- DONE — `go test ./...`: passed.
- DONE — `go test ./... -race`: passed.
- DONE — Local npm smoke: passed for the pinned tarball/integrity flow and pi-ai `/compat` guard using temporary install prefixes.

Stage: validation

### Summary

Validated implementation commit `158927c13fe3db3398bfbca8f0b1d970b6166f3d` for the pi-live version-skew fix. Recommendation: **PASSED** for the offline validation gate. AC-2 is satisfied by exact pinned versions, recorded sha512 integrities, install-step compatibility comments, and a fail-fast `/compat` guard before the live smoke. AC-1 remains intentionally reserved for the FO's post-merge main-branch pi-live run.

### Acceptance criteria evidence

- **AC-1 (VALUE) - The pi-live lane is green on main.** Not independently verified in this validation stage; dispatch explicitly assigns post-merge main evidence to the FO. This does not block this offline validation gate.
- **AC-2 - The version constraint is explicit and evergreen, not incidental.** PASS. Static review of `.github/workflows/runtime-live-e2e.yml` confirms the pi-live job uses Node 24, declares exact pins for `@earendil-works/pi-coding-agent@0.80.10`, `pi-subagents@0.35.1`, and `pi-intercom@0.6.0`, records sha512 integrity strings in-file, verifies `npm pack` integrity before installing tarballs, contains the compatibility comment naming `@earendil-works/pi-ai/compat` and the Node engine requirement, and runs `Guard Pi substrate compatibility` before `Run live Pi front-door smoke`.

### Validation

- PASS — `gofmt -w ./cmd ./internal && git diff --check && git status --short`: no output; worktree remained clean.
- PASS — `go test ./...`: all packages passed on this run.
- PASS — `go test ./... -race`: all packages passed on this run.
- PASS — Detached adversarial audit on hidden throwaway checkout `.adversarial-audit`: corrupted `PI_SUBAGENTS_INTEGRITY`, then ran `go test ./internal/release -run TestWorkflowsPreserveAndPublishJourneyCosts -count=1`; the guard failed with `runtime-live-e2e.yml Pi live job does not declare exact Pi substrate version and integrity pins`, proving the deliverable's guard catches a material integrity-pin regression. Throwaway checkout was removed after the audit.

### Reviewer findings

No material findings. The offline guard is static by design for the workflow contract, but the adversarial audit confirms it is not a self-referential spelling-only acceptance of the implemented file.

### Deferred risks

- AC-1 live evidence is pending the FO-owned post-merge main pi-live run. Promotion condition: if the PR or main pi-live run still fails at Pi extension load or substrate compatibility after merge, reopen as material against AC-1.
