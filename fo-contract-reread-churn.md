---
id: kc1rvn663yt8qkzqbakzda1v
title: The FO contract's load triggers cause re-reads of resident, unchanged files
status: backlog
source: "email-triage codex session audit 2026-08-26: 59 skill-file reads in one FO day — fo-write-core.md 14x, fo-gate-lifecycle 10x — about 34% of the FO's tool calls, against files that never changed; only two compactions occurred, so at most three reads per file were mandated"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

The contract's deferred load points phrase their triggers per occurrence: "load before every selected gate", "read immediately before the first FO-authored mutation". The contract never states that a file already resident in context satisfies the precondition. A literal reader therefore re-reads the full file at every trigger occurrence. One measured day: 14 reads of the same unchanged write-core, roughly one per mutation boundary.

## Problem

{Ideation fills this in. Seeded facts: the post-compaction rule ("re-satisfy each load precondition at its existing trigger") is correct and must stay — two compactions justify two re-load rounds; the measured excess is within-window re-reads of resident files. The cost compounds: each read is a ~200-line file, and when a command flag failed, the same session fell back to --help dumps and `strings` over the binary to rediscover semantics — the load model gives a literal reader no cheap way to re-check what it already holds. Different hosts diverge: a sticky reader treats loads as once-per-context; a literal reader re-executes per occurrence; the contract supports both readings.}

## Proposed approach

{Ideation fills this in. Captain constraint (2026-08-26): the seed states the problem only — no solution is prescribed here. Post-compaction re-reads remain mandated whatever the design.}

## Risk evidence

{Backlog: the audited codex FO thread (59 reads, per-file counts, two compaction timestamps) decides design should start. The audit report is in the 2026-08-26 session record.}

## Out of scope

The post-compaction re-satisfy rule. Host-specific context management.

## Expected surface and tolerance

{Backlog seed; ideation estimates once a direction exists.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A full FO day over an unchanged contract spends materially fewer tool calls on contract reads than the measured baseline, with post-compaction re-satisfaction intact.**
Verified by: {ideation refines — seed: the measured baseline is 59 reads/day at 34% of FO calls; the falsifying comparison is a live or replayed FO day's read count against it.}

## Test plan

{Ideation fills this in.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
