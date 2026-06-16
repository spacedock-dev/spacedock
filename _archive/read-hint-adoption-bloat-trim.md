---
title: Trim redundant `status --read` adoption instruction sites once usage is measured
status: done
source: "captain (2026-06-16, 0204 sprint) — the --read adoption guidance reaches the ensign in three places it sees every dispatch (ensign-shared-core.md:18, :92, and the build.go:553 site-6 dispatch-prompt hint added to every prompt — visible as the repeated line across all 13 dispatch goldens). Site-6 reinforces the contract sites the ensign already loads via Skill. Heavier instruction is not stronger proof; it is per-dispatch prompt bloat."
score:
sprint: 0204-structured-reads
sprint-readiness: ready
issue:
id: 4xghzpa1wjqw2vh5s1h24d5z
verdict: superseded
completed: 2026-06-16T20:25:54Z
---

## Problem
The `status --read` adoption guidance is instruction-heavy: `ensign-shared-core.md:18` (read sections), `ensign-shared-core.md:92` (stage-report append point), AND the `internal/dispatch/build.go` site-6 prompt hint added to EVERY dispatch prompt (the line repeated across all 13 dispatch goldens). Site-6 largely reinforces sites 4/5, which the ensign already loads via its `Skill` call. That is prompt bloat on every dispatch.

## What's needed
Once the journeymetrics `--read` adoption metric can measure actual usage, trim the redundant site(s) — most likely the site-6 dispatch-prompt hint — and confirm `--read` usage stays stable (behavioral proof the trimmed site was not load-bearing). Regenerate the dispatch goldens.

## Dependency
Depends on `journeymetrics-read-adoption-metric` — measure first, then trim with evidence. Do not blind-trim.

## Acceptance criteria
- **AC-1** — After trimming the redundant site(s), measured `status --read` usage in a live FO/ensign drive is not lower than before the trim (behavioral, via the journeymetrics metric), and the dispatch golden suite is regenerated and green. No prose-grep over the instruction files.
