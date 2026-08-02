---
title: Make `spacedock status --where` robust, complete, and discoverable (so the FO queries, never greps)
status: ideation
sprint:
score: ""
source: "GitHub #314 (status --where silently matches-all on an unknown field) + an FO repeatedly falling back to find/grep over the state dir. Three gaps on one surface push the FO/Commander off the native query onto raw shell. FO session 2026-07-04: a live example of gap 3 — the FO pulled the full ~50-row board (no --where filter) to locate 2 rows (one in-flight, one gated) instead of a filtered query. Added as a motivating example for AC-3's discoverability case."
priority: medium
id: 3t9r36n9tbj116jp9g1k01tz
gates:
    version: 1
    current:
        gate: gate:3t9r36n9tbj116jp9g1k01tz:backlog
    records:
        - id: gate:3t9r36n9tbj116jp9g1k01tz:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3t9r36n9tbj116jp9g1k01tz-backlog-1
              briefing:
                id: briefing:3t9r36n9tbj116jp9g1k01tz:backlog:attempt-1:revision-1
                digest: sha256:7e024b6712a45735483d7b6546e0aef114b7a7bb385b049038147d585537c0d6
                digest-domain: canonical-bytes
                request-digest: sha256:5079ac1476c67b3215a065329c37f6183a0eb29a278e83783cbd76026ae623b3
                room-ref: ./status-where-robust-and-discoverable/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3t9r36n9tbj116jp9g1k01tz:backlog:1
                briefing: briefing:3t9r36n9tbj116jp9g1k01tz:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T16:02:24.444201Z"
                decision: approve
                reason: 'Captain directed in chat: ''dispatch both''. Approve backlog->ideation for the --where robustness/discoverability/GH#314 gaps.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
started: 2026-08-02T16:02:35Z
---

## Problem

`spacedock status --where` is the entity query, but three gaps push the FO/Commander to raw `find`/`rg`/grep over `.spacedock-state/*.md` — slower, error-prone, and blind to archived state.

### 1. No field-name validation → silent wrong results (GitHub #314)
`--where` validates the operator but not the field. An unknown/misspelled field reads as the empty string in `applyFilters`, so:
- A typo (`--where spint=foo`) silently returns the wrong set.
- The intuitive compound-in-one-string (`--where 'sprint=A sprint-readiness!=defer'`) `Cut`s on the first operator → a garbage field name, and with `!=` it **silently matches every entity** (a confidently-wrong, unfiltered result).
The silent match-all is the dangerous part; a loud error would point straight at the fix (repeated `--where` flags). Root cause + repro in #314 — `internal/status/parse.go` `parseWhereFilters` (operator-only validation, ~L187) and `applyFilters` (unknown field → `""`, ~L215).

### 2. `--where` excludes done/archived → no full-sprint view
`--where sprint=X` returns only ACTIVE entities; done/archived members drop out, so "where is sprint X, completely (incl. done)" cannot be answered in one query — it forces a grep of `_archive/`. Observed live: querying a nearly-complete sprint surfaced only the in-flight members and hid the done ones, which is what drove the fall-back to `find`.

### 3. Under-documented → the FO greps instead
`--where <field>=<value>` is the canonical entity query, but it is not surfaced as THE query in the FO contract or `status --help` (the known-field list, the repeated-flag AND-semantics, the archived flag). So the FO reflexively writes `find`/grep over the state files — slower, blind to `_archive`, and prone to the confusing partial answers raw greps give.

## Acceptance criteria

- **AC-1 (value / #314)** — An unknown or misspelled `--where` field ERRORS with a clear message naming the known fields (à la `whereSyntaxHelp`), instead of silently matching. Independent baseline that moves the wrong way: today `--where 'nosuchfield!=x'` AND the compound `--where 'sprint=A sprint-readiness!=defer'` both silently return ALL entities (verified per #314); the fix must turn both into a loud, fixable error. Proven by a test exercising the unknown-field and compound-string cases — not a prose assertion.
- **AC-2 (full-sprint view)** — A single command answers "all members of sprint X across every stage, incl. done/archived" — e.g. `--where sprint=X --include-archived` (or `--all`), or a sprint-rollup view. Baseline: today `--where sprint=X` omits the done/archived members. Verified on a sprint that has both active and archived members (the done ones now appear).
- **AC-3 (discoverability)** — The FO contract and `spacedock status --help` document `--where <field>=<value>` as THE entity query: the known-field list, the AND-via-repeated-flags rule (the exact footgun #314 stems from), and the archived-inclusive flag. Cross-checked that a fresh FO can answer "list entities at status X" / "all of sprint Y incl. done" from the docs without reaching for shell.
- **AC-4 (validation)** — `go test ./internal/status/…` green; the field-validation and archived-inclusive behaviors are covered by tests (unknown-field error, compound-string error, archived-inclusive query returns the done members).

## Notes
- **GitHub issue #314** (`status --where: unknown field name silently matches-all`) — this task subsumes it; close #314 when AC-1 lands.
- The three gaps share one theme: `--where` should be **robust** (errors on misuse), **complete** (queryable across archived), and **discoverable** (documented as THE query) so the FO never falls back to raw shell — kin to the smallest-sufficient-mechanism gate (`fo-smallest-sufficient-mechanism`).
