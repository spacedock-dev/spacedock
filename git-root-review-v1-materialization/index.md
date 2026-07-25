---
title: Materialize Git-root Review v1 sources for provider presentation
status: backlog
source: "s4 cycle-6 staff rejection: recorder-valid git-root:// sources are not renderable by current Subspace package mode, 2026-07-25"
started:
completed:
verdict:
score: "1.0"
worktree:
issue:
sprint: durable-decisions
id: rqh46ey33aqq4rt72b4w1m2q
gates:
    version: 1
    current:
        gate: gate:docs-dev:rqh4:backlog
    records:
        - id: gate:docs-dev:rqh4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rqh4-backlog-1
              briefing:
                id: briefing:docs-dev:rqh4:backlog:attempt-1:revision-1
                digest: sha256:d620934ee0af1b72c38e80fdb640f6ea07bd95da9fd08729c38e9b9d04a4fce2
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:rqh4:backlog:1
                briefing: briefing:docs-dev:rqh4:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T06:28:30.836574Z"
                decision: approve
                reason: The task isolates the actual missing cross-repository consumer boundary, forbids durable source duplication, and requires a real moved-root Subspace proof before implementation.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

Bridge recorder-ready Git-addressed Briefings to actual provider presentation without
duplicating selected sources in durable gate state.

## Problem

s4's smallest source identity is
`git-root://<main|state>/<full-commit>/<repository-relative-path>` plus Review v1's raw
SHA-256 `rev`. Spacedock can reopen and verify that blob through the named split-root
Git object database after checkouts move independently. Current Subspace package mode
cannot present it: Artifact URIs are opened as filesystem paths and Reference URIs
containing `://` are rejected.

The missing owner is not terminal transport. It is the provider-neutral boundary that
turns a root/commit/path/raw-SHA reference into the exact bytes a Review v1 presenter
renders, while keeping the durable source singular in Git.

## Boundary to decide in ideation

Compare the smallest real options at the consumer boundary:

- a provider-neutral resolved-source manifest and ephemeral presentation materialization;
- a narrow Git-root resolver capability implemented by the provider;
- another mechanism that supplies exact verified bytes without rewriting durable gate
  history or making Subspace understand Spacedock workflow layout implicitly.

Name which repository owns each half, the stable API between them, retention and cleanup
of any ephemeral bytes, local-object/shallow/pruned-history failure behavior, and how
the presented inventory maps back to the canonical Git-root Artifact/Reference ids and
revisions. Do not treat q0 transport alone as resolution.

## Acceptance criteria

**AC-1 (VALUE) — A real provider presents every canonical Git-root source after main
and state checkouts move independently.** Starting with one s4-generated two-file room,
the original checkout nesting becomes unavailable; the supported integration resolves
the named committed objects, verifies their raw SHA-256 revisions, and Subspace displays
the complete Artifact/Reference inventory. Missing root, object, path, or digest mismatch
fails before a binding Result.

**AC-2 — Durable gate state contains no copied selected source payload.** The retained
room remains request plus canonical Briefing; any resolved bytes are provider-owned
ephemeral material or native object reads with explicit lifecycle and cleanup. Git
history remains the sole durable payload owner.

**AC-3 — Result association remains canonical and provider-neutral.** Presented
inventory maps every displayed item exactly once to the original Briefing id, source id,
Git-root URI, and raw revision. Spacedock still derives and verifies the association and
stores no `association.json`.

**AC-4 — The cross-repository contract is exercised, not inferred.** A committed
end-to-end test uses the actual Spacedock-generated room and actual Subspace package
mode/capability from independently located split roots. A local `git cat-file` helper
test alone does not satisfy this criterion.

## Initial test gate

Ideation must first run the smallest throwaway room-to-provider spike through current
Subspace and preserve its concrete failure/success boundary. It must declare expected
files and LOC in each repository before implementation. No compatibility wrapper,
remote fetch, generic URI framework, or durable source-copy cache is implied.
