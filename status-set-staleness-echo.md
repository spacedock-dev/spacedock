---
id: j7tbbka0fpd9twbqgr7bhrqd
title: Eliminate the Read-then-`status --set` staleness echo — mutating a just-read file re-emits the whole file as cache-write tokens
status: backlog
source: FO + 0.20.4 scope survey (2026-06-14, this session) — on Claude Code, a `Read` of an entity file followed by a `status --set` mutation of the same file trips the harness file-staleness safety net, re-emitting the whole file as cache-write tokens. A recurring FO/ensign tax (validate-then-set, dispatch-then-set). Cited in e6a's Problem section and the sprint friction notes. Distinct from e6a (which avoids the whole-file READ); this attacks the whole-file WRITE-echo on mutation. 0.20.4 read-cost theme.
started:
completed:
verdict:
score: 0.33
worktree:
issue:
sprint: 0204-structured-reads
sprint-readiness: ready
---

A `Read` of an entity file followed by a `status --set` on the same file re-emits the whole file as cache-write tokens (the Claude Code file-staleness safety net). The FO contract works around it today (Grep-over-Read; trust `--set` stdout), but the echo still fires whenever a Read before a set is unavoidable. Kill the echo, or prove it is harness-inherent and document the avoidance.

## Problem

When an agent `Read`s a file and then mutates it via Bash (including `spacedock status --set`), Claude Code's file-staleness safety net re-emits the whole file as cache-write tokens on the next turn. For a 280-line entity that is a whole-file tax on every read-then-mutate. The FO contract (`first-officer-shared-core.md:218`) already names this and routes around it, but the workaround is avoidance, not a fix.

## Proposed approach

{Ideation fills. SPIKE THE MECHANISM FIRST (this is the riskiest unknown): replicate the exact `Read -> status --set -> next turn` pattern in a controlled exercise and measure the token flow to determine WHERE the echo originates — the Claude Code harness's file-tracking, or anything `status --set` controls. The outcome gates the task:
- If `status --set` (or an atomic-write / quiet-mode) can suppress the re-emit, candidate directions: a `--no-echo`/quiet `--set`, an atomic write that does not trip staleness, or an alternate mutation surface. Build it.
- If the echo is purely harness-inherent (the tool cannot influence the harness's cache-write decision), this task produces a DOCUMENTED avoidance pattern and ships as a roadmap decision, NOT a code deliverable (per the dev-workflow rule that a decision with nothing shipped belongs in the roadmap).}

## Out of scope

{Ideation fills. Likely: e6a's read-side section helper (separate); non-`status` mutations.}

## Acceptance criteria

Each AC names a property of the finished outcome, not a stage action, and how it is verified.

**AC-1 — A read-then-`status --set` on the same file does not re-emit the whole file as cache-write tokens (or the task is recorded as a roadmap decision with the measured reason it cannot).** 
Verified by: {a measured before/after token-flow trace in a controlled exercise (external oracle = token counts), never prose. If the spike proves harness-inherence, the recorded measurement + roadmap decision IS the satisfied outcome.}

**AC-2 — `status --set` still narrates its mutation (`field: old -> new`) so the FO needs no re-read.**
Verified by: {a test asserting `--set` stdout carries the mutation narration; ideation pins it.}

## Test plan

{Ideation fills; the spike seeds it. Note the bimodal outcome: a code fix proven by token-delta measurement, OR a roadmap decision proven by the measurement that shows harness-inherence.}
