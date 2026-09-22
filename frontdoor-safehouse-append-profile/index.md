---
title: Forward appended Safehouse profiles through the front door
status: ideation
source: Captain request 2026-09-22
started: 2026-09-22T15:38:57Z
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
id: 9cbmcq3yyfjynhms1yd12drh
---

Support `--safehouse-append-profile=file.sb` on Spacedock runtime front doors so users can add a Safehouse profile without bypassing the launcher. Captain explicitly requested filing and dispatch to ideation.

## Problem

The front door exposes existing Safehouse knobs but does not register an appended-profile option. Determine the smallest compatible pass-through across supported front doors.

## Proposed approach

Extend the existing Safehouse flag parsing and launch translation. Ideation must verify Safehouse’s actual option semantics and propose exact forwarding, path handling, repetitions/order, delimiter behavior, and docs before implementation.

## Risk evidence

Existing flag owners are internal/cli/frontdoor.go, pi.go, help.go and their parser/launch-parity tests. Probe the installed Safehouse CLI and existing translation path without altering global configuration or starting model sessions.

## Out of scope

New sandbox policy syntax, policy composition engines, auth changes, and unrelated front-door refactors.

## Expected surface and tolerance

Ideation must estimate net LOC change and file count from existing owners, with explicit tolerance and allowed semantic changes. Prefer the existing pass-through mechanism.

## Acceptance criteria

**AC-1 — A front-door invocation with --safehouse-append-profile=file.sb launches through Safehouse with the requested appended profile.**
Verified by: existing launch-seam fixtures observe the actual Safehouse argv and unchanged host command across supported front doors; omission or host-argument leakage fails the test.

**AC-2 — Profile paths and option boundaries retain their intended meaning.**
Verified by: argument tests cover equals-form, normal space-form consistency, paths containing spaces, relative paths, missing values, and tokens after --. Repetition/order and error handling follow the verified Safehouse contract; a malformed invocation never silently launches without its requested profile.

**AC-3 — Existing Safehouse and host argument behavior remains compatible.**
Verified by: existing front-door parser/launch tests preserve other knobs, sandbox selection and host passthrough; CLI help and user documentation show the supported example.

## Test plan

Ideation identifies existing primary proof owners, runs a bounded probe of Safehouse’s actual option and forwarding, and proposes focused failing tests before changes. Implementation runs applicable Go normal/race and formatting checks. No local model runs. Provide a concrete user-doc diff and expected surface for design review.

### Feedback Cycles

