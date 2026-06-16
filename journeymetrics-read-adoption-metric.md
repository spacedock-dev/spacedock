---
title: Measure actual `status --read` adoption in FO/ensign journeys via journeymetrics
status: backlog
source: "captain (2026-06-16, 0204 sprint) — e6a's status --read adoption was proven WORKING once (AC6 single live drive). No ongoing verification that real FO/ensign sessions call --read; the six contract sites are wording, and wording-present is not behavior. journeymetrics already parses every journey transcript but records ToolCallsByName by tool NAME only — `status --read` is a Bash subcommand, invisible today as a generic Bash call."
score:
sprint: 0204-structured-reads
sprint-readiness: ready
issue:
id: hf4jmbksapyg2d9s0zj85wca
---

## Problem
e6a's `status --read` adoption is proven to work exactly once (AC6's single live FO drive during e6a validation). There is no standing, behavioral check that real FO/ensign sessions actually invoke `--read`. The instruction lives in six contract sites; that is wording, not observed behavior.

`internal/journeymetrics` parses every journey transcript and records `ToolCallsByName` — but keyed by tool NAME only (`toolCallsByName[block.Name]++`). `status --read` is a Bash subcommand, so it is counted as a generic `Bash` call and is invisible.

## What's needed
Extend `internal/journeymetrics` (claude.go + codex.go) so a journey record surfaces `--read` adoption: detect `status --read` in the Bash `tool_use` command string, and/or count scoped `Read` calls carrying `offset`/`lines`. Surface a per-journey count in the ledger so adoption is measurable over time.

## Why
The durable behavioral proof that the e6a adoption sticks — and the evidence that should drive the redundant-instruction-site trim (see read-hint-adoption-bloat-trim). Measure before trimming.

## Acceptance criteria
- **AC-1** — A journey record surfaces a count of `status --read` invocations (Bash-arg detection) and/or scoped-`Read` calls. Verified by a journeymetrics unit test over a fixture transcript containing a `status --read` Bash call + a scoped `Read`, asserting the count is captured, and zero when absent (non-tautological by perturbation).
