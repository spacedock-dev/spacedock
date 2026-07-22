---
title: Cold Commander package cannot consume shaping-era gate approvals after strict-v1 recorder lands
status: backlog
source: Durable-decisions cold-boot dogfood, 2026-07-22
score: "0.75"
id: h49jd1at66bkc6w131we7d87
---

The durable-decisions execution package directs a cold Commander to consume closed shaping-era approvals, but those entities retain pilot gate fields that the landed strict-v1 recorder must reject. The session required an exact captain-authorized state override before implementation dispatch could proceed.

## Problem

After the canonical-v1 recorder landed, the package still described `h1`, `02av`, and `xb` as ordinary pending applications to consume. Reproduction from repo root:

```text
go run ./cmd/spacedock gate validate gate-blockers-and-eligibility --workflow-dir docs/dev
```

The command exits 1 and reports unsupported pilot fields including `current.attempt`, `current-attempt`, `sequence`, `state`, and shaping notes. Durable state evidence is the shaping history in `gate-blockers-and-eligibility/index.md`, `ensign-finding-triage-disposition/index.md`, and `gate-review-presentation-command/index.md`; each held a closed approval plus `application.state: pending` while lifecycle status remained `ideation`.

Expected behavior: a packaged cold-boot execution plan must make the transition boundary executable without asking the strict recorder to accept or migrate pilot encodings, without leaving a false pending authorization after dispatch, and without requiring an improvised nested-frontmatter edit outside the FO write contract.

## Ownership boundary

This is execution-package and state-transition preparation friction. It does not belong in the recorder's canonical-v1 parser, and it does not authorize compatibility or migration code. A solution may pre-close shaping applications before strict-v1 becomes authoritative or package a named one-time transition mechanism with exact state scope and provenance. It must preserve the strict recorder boundary.

## Acceptance criteria

**AC-1 (VALUE)** A clean cold Commander can follow a packaged sprint from recorded shaping approval to implementation dispatch with no unsupported-schema error, no manual nested-frontmatter edit, and no pending authorization left behind.
Verified by: a fixture-backed cold-boot replay using a strict-v1 recorder plus shaping-era approved state; it fails if the recorder is asked to accept pilot fields or if post-dispatch state still reports the application pending.

**AC-2** The transition adds no pilot compatibility, migration branch, or unknown-field acceptance to the canonical recorder.
Verified by: the recorder's existing prototype-rejection tests remain green and the transition replay succeeds outside the recorder parse/write path.

**AC-3** The package records exact transition authority and preserves the historical approval bytes or an immutable reference to them.
Verified by: resulting state plus committed room/package evidence reconstruct who approved, which Briefing was reviewed, and how the one-time application was consumed.

## Promotion condition

Promote when another sprint packages approvals before a strict schema cut, or when the durable-decisions cold-boot replay is made repeatable for release verification. Until then this remains orthogonal backlog friction and does not expand the current sprint.
