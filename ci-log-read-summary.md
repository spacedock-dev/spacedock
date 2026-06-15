---
id: 6rtpj5avcp733tb15dfjcbbb
title: Summarize CI artifacts for FO/ensign reads — replace whole-log (143KB) reads with a triage summary
status: backlog
source: FO + 0.20.4 scope survey (2026-06-14, this session) — CI logs (143KB cited in e6a's source as a recurring read sink) are read whole into FO/ensign context for validation and triage, when the agent needs only pass/fail + the failing lines. Tokens scale with the whole log. 0.20.4 read-cost theme; lower-frequency than the entity/README reads e6a covers.
started:
completed:
verdict:
score: 0.30
worktree:
issue:
sprint: 0204-structured-reads
---

CI logs are read whole into context for triage when the agent needs only pass/fail, the exit code, and the failing tests/lines. Give the FO/ensign a triage summary they can read instead of the 143KB log, while keeping the full artifact reachable when the summary is not enough.

## Problem

Validation and triage read large CI logs whole (143KB cited). The signal an agent needs is small: did it pass, what failed, where. The rest is noise that fills context. There is no summary surface, so callers load the whole artifact or grep blind.

## Proposed approach

{Ideation fills. SPIKE FIRST: take a real large CI log and produce a 10-20 line summary (pass/fail, exit code, top-N failures with file:line), then confirm an FO/ensign can triage correctly from the summary alone. Candidate directions to evaluate: a `spacedock`-side log-summarize helper, or a CI-emitted summary `.md` artifact (status bullets + links) the agent reads instead of the log. Decide where the summary is produced (CI side vs read-time) and how the full log stays reachable.}

## Out of scope

{Ideation fills. Likely: e6a's entity/README section reads (separate); non-CI logs; changing what CI itself runs.}

## Acceptance criteria

Each AC names a property of the finished outcome, not a stage action, and how it is verified.

**AC-1 — A summary surface yields pass/fail + the failing tests/lines for a real large CI log in a small read, and an FO/ensign triages correctly from it without the whole log.**
Verified by: {a live triage drive against a real failing-log fixture (external oracle = the actual failures present in the log), behavioral — not a prose match.}

**AC-2 — The full log stays reachable when the summary is insufficient (no information loss).**
Verified by: {a test/exercise confirming the summary points to or preserves the full artifact; ideation pins it.}

## Test plan

{Ideation fills; the spike seeds it.}
