---
title: Fold the dispatched-ensign sub-agent transcript into the journeymetrics --read adoption metric
status: backlog
source: "hf validation (2026-06-16) — the journeymetrics --read adoption metric (status_read_calls/scoped_read_calls) parses the FO front-door `claude -p` stream only. The site-6 dispatch-prompt hint hf trimmed targets the DISPATCHED ENSIGN's reads, and the ensign runs as a separate team-agent session whose transcript the metric does not capture — four real FO captures all measure 0/0. The standing monitor does not yet observe the agent it was built to watch."
sprint: 0204-structured-reads
sprint-readiness: ready
issue:
id: f53zr0ehhzekzbgydpybaq5g
---

## Problem
The journeymetrics `--read` adoption metric (`status_read_calls`/`scoped_read_calls`, landed in hf) parses the FO front-door `spacedock claude -p` stream. But `status --read` adoption is principally an ENSIGN behavior — the contract sites that steer repeat-reads (and the site-6 hint hf trimmed) target the dispatched ensign, which runs as a separate team-agent session (`subagents/agent-*.jsonl`) whose transcript the metric does not parse. Empirically, four real FO front-door captures all register `status_read_calls=0`/`scoped_read_calls=0`. So the standing monitor does not observe the agent whose adoption it was built to measure, and hf's AC-7 before/after could only compare 0==0.

## What's needed
Extend the journeymetrics parse to fold the dispatched-ensign sub-agent transcript (`subagents/agent-*.jsonl` under the run's config dir) into the per-journey `--read` adoption counts, so the metric measures the ensign's `status --read` + scoped-`Read` behavior — the surface the contract sites and the trimmed site-6 hint actually steer. Then hf's adoption-not-regress check (pre-trim vs post-trim) becomes a meaningful before/after.

## Acceptance criteria
- **AC-1** — A journey record's `status_read_calls`/`scoped_read_calls` include the dispatched ensign sub-agent transcript's `status --read` Bash invocations and scoped `Read` calls, not just the FO front-door stream. Verified by a journeymetrics unit test over a fixture journey that includes an ensign sub-agent transcript carrying a `status --read` Bash call + a scoped `Read`, asserting the counts capture them, and zero when the ensign transcript has neither. Non-tautological by perturbation (removing the ensign-transcript fold drops the count to the FO-only value and FAILs).
