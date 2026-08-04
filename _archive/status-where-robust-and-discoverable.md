---
title: Make `spacedock status --where` robust, complete, and discoverable (so the FO queries, never greps)
status: done
sprint:
score: ""
source: "GitHub #314 (status --where silently matches-all on an unknown field) + an FO repeatedly falling back to find/grep over the state dir. Three gaps on one surface push the FO/Commander off the native query onto raw shell. FO session 2026-07-04: a live example of gap 3 — the FO pulled the full ~50-row board (no --where filter) to locate 2 rows (one in-flight, one gated) instead of a filtered query. Added as a motivating example for AC-3's discoverability case."
priority: medium
id: 3t9r36n9tbj116jp9g1k01tz
gates:
    version: 1
    records:
        - id: gate:3t9r36n9tbj116jp9g1k01tz:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3t9r36n9tbj116jp9g1k01tz-backlog-1
              briefing:
                id: briefing:3t9r36n9tbj116jp9g1k01tz:backlog:attempt-1:revision-1
                digest: sha256:7e024b6712a45735483d7b6546e0aef114b7a7bb385b049038147d585537c0d6
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
                target-stage: ideation
                state: consumed
        - id: gate:3t9r36n9tbj116jp9g1k01tz:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:3t9r36n9tbj116jp9g1k01tz-ideation-1
              briefing:
                id: briefing:3t9r36n9tbj116jp9g1k01tz:ideation:attempt-1:revision-1
                digest: sha256:25dc9938c0182e38f2158831de267b500e6a43326a3d140a4af8e67daf31b467
                request-digest: sha256:86055f2d29fbfab3d7f5fba21ad4ee473cc25c83a0eba9e4c1d32ea77ae8dd38
                room-ref: ./status-where-robust-and-discoverable/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3t9r36n9tbj116jp9g1k01tz:ideation:1
                briefing: briefing:3t9r36n9tbj116jp9g1k01tz:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T16:15:14.729845Z"
                decision: approve
                reason: 'Captain approved ideation in chat: AC-2 mechanism correctly deleted, AC-1 split into 2 guards, AC-3 status --help gap confirmed. Proceed to implementation.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:3t9r36n9tbj116jp9g1k01tz:validation
          stage: validation
          attempts:
            - id: gate-attempt:3t9r36n9tbj116jp9g1k01tz-validation-1
              briefing:
                id: briefing:3t9r36n9tbj116jp9g1k01tz:validation:attempt-1:revision-1
                digest: sha256:dfe84eb83bc263b11ff5f4e8e79e30e6ccc6050577dbed44c540fb6c8ff1e4cb
                request-digest: sha256:0adda560d2e04428a87e52ede91c96f948f4eea16ed59a7e144fd72cf94490fe
                room-ref: ./status-where-robust-and-discoverable/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3t9r36n9tbj116jp9g1k01tz:validation:1
                briefing: briefing:3t9r36n9tbj116jp9g1k01tz:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T03:30:20.857642Z"
                decision: approve
                reason: 'Captain approved in chat (''push status-where-robust-and-discoverable''): validation PASSED, all 4 ACs independently reproduced against baseline+HEAD binaries, 1 non-material deferred risk. Proceed to merge.'
              application:
                target-stage: done
                state: consumed
started: 2026-08-02T16:02:35Z
worktree: .worktrees/spacedock-ensign-status-where-robust-and-discoverable
mod-block:
pr: pr-merge:605
verdict: passed
completed: 2026-08-03T04:00:18Z
archived: 2026-08-03T04:00:19Z
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

## Measured baseline (ideation, 2026-08-03, `docs/dev`, 156 active / 495 with `--archived`)

Every number below was produced by running the binary built from this checkout.

| Invocation | Rows today | Correct answer |
| --- | --- | --- |
| (no `--where`) | 156 | — |
| `--where 'nosuchfield!=x'` | **156**, exit 0 | error |
| `--where 'sprint=A sprint-readiness!=defer'` | **156**, exit 0 | error |
| `--where 'sprint!=A sprint-readiness!=defer'` | **156**, exit 0 | error |
| `--where 'status!=backlog sprint=A'` | **156**, exit 0 | error |
| `--where 'sprint=A sprint-readiness=ready'` | 0, exit 0 | error (residual, see below) |
| `--where 'spint=0250-…'` (typo) | 0, exit 0 | error |
| `--where sprint=0202-survey-improvements` | 6 | 6 active |
| `--where sprint=0202-survey-improvements --archived` | **11** | 11 = 6 active + 5 archived |

**Gap 2 does not exist as a capability gap.** `--archived` is already *inclusive*:
`activeAndArchivedEntities` (`internal/status/discover.go:326`) returns
`scanEntitiesActive(...) + archiveEntities(...)`. `--where sprint=X --archived`
already answers the full-sprint question (11 = 6 + 5, measured above), and the
merged list is already legible through canonical fields — measured:

```
--where sprint=0202-survey-improvements --archived --fields slug,status,verdict,archived
  survey-lens-honesty          backlog     v=-       arch=-
  …
  survey-output-redesign       validation  v=PASSED  arch=2026-06-13T21:12:39Z
```

So AC-2 needs **no new flag and no new binary mechanism**. `--include-archived`
/ `--all` would be a synonym for a shipped flag — pure command-surface sprawl
(cf. `dw command-surface-gate`). What actually failed is the FO contract, which
calls `--archived` an "**Archive view**" (`skills/fo-status-viewer/SKILL.md:41`)
— wording that reads as *archived-only*, i.e. as a scope swap. The FO believed
it had to choose between active and archived and reached for `find`. AC-2
collapses into AC-3 plus a regression test that pins the composition.

**Field-name validation alone does not close #314.** Shape `sprint!=A
sprint-readiness!=defer` cuts on the *first* `!=`, yielding the **known** field
`sprint` with the garbage value `A sprint-readiness!=defer`; no entity holds
that value, so `!=` matches all 156 — silently, with a valid field name. Three
of the four compound shapes silently match all; only a multi-operator *syntax*
check catches all of them. This is why the design needs two independent checks,
not one.

## Proposed approach

Two independent guards plus documentation. No new flags.

### Check A — one clause per `--where` (syntax; `parseWhereFilters`, `parse.go:172`)

Count operators in the argument: occurrences of `!=`, plus occurrences of `=`
not preceded by `!`. More than one ⇒ error. Needs no entity corpus, so it fires
even on an empty workflow, and it catches all four compound shapes.

```
Error: --where takes one clause per flag; 'sprint=A sprint-readiness!=defer' has 2 operators — repeat --where to AND clauses: --where 'sprint=A' --where 'sprint-readiness!=defer'
```

*Simplest alternative considered:* reject whitespace in the field name only.
Insufficient — it misses `sprint!=A sprint-readiness!=defer`, whose field name
is clean (measured: 156 rows). *Cost of the chosen rule:* a value legitimately
containing `=` is rejected. Measured blast radius: 3 of 495 entities have a
non-`source` value containing `=`, all in `title` — and `--where title=<prose>`
is not a query anyone issues (`--resolve`/`slug` is). Declared as accepted.

### Check B — known field names (semantic; new `validateWhereFields`, called from `handlers.go:495`)

`applyFilters` has exactly one call site (`handlers.go:495`), immediately after
the two materializers (`:493`, `:494`), so the check slots in there with the
full entity set in hand.

```
knownFields = ∪ keys(e.fields) over scanned entities        // sprint, sprint-readiness, group, milestone, …
            ∪ keys(loadEntitySchema().fields)               // canonical: verdict, completed, archived, mod-block, pr, …
            ∪ {gate-condition, gate-eligible, gate-readiness, next-suppressed-by}   // derived
```

```
Error: --where: unknown field 'spint' — known fields: archived, completed, gate-condition, …, worktree
```

- **The schema union is required, not belt-and-braces.** The schema is
  explicitly open (`permissive_additions: true`, `custom_fields.policy:
  preserve_unknown`), so `sprint` is *not* canonical and a schema-only list
  would reject the single most common query. Conversely `verdict`, `completed`
  and `archived` are canonical but **archive-only in the active corpus**
  (measured), so a corpus-only list would break `--where 'verdict!='` on an
  active read. Each half covers the other's hole.
- **The four derived names are static** rather than read from `e.fields` after
  materialization: `materializeGateEligibility` `continue`s on error
  (`discover.go:299-302`), so a workflow where eligibility fails for every
  entity would otherwise lose the name. Four constants remove all ordering
  dependence.
- **Guard:** zero scanned entities ⇒ skip Check B (the result is empty either
  way; without this, every `--where` on a fresh workflow would error).
- `--fields` has the same silent-empty-projection behavior. **Out of scope** —
  this entity is `--where`, and `--fields` feeds pinned JSON/golden surfaces
  (`section_read.go:192`) whose consumers may rely on the empty-key projection.

### Check C — none. `--archived` is left exactly as it is.

## Documentation changes

### `spacedock status --help` — a surface that does not exist today

Measured: `spacedock status --help` **runs the entity listing** (exit 0, 156
rows); `spacedock help status` prints the top-level menu. `newStatusCommand`
(`internal/cli/cli.go:564`) has `DisableFlagParsing: true` and, unlike
`newNewCommand` (`:592`), no `wantsHelp(args)` guard. AC-3 therefore *creates* a
help surface, it does not amend one. The mechanism is proven: `spacedock new
--help` renders correctly under the identical cobra configuration.

Wire `wantsHelp(args) → cmd.Help()` into `newStatusCommand` and add
`setStatusHelp` to `internal/cli/help.go` beside `setNewHelp`/`setMergeHelp`:

```
Show or update workflow state

Usage:
  spacedock status --workflow-dir DIR [query flags]

Query flags:
  --where FIELD=VALUE  Filter entities — THE entity query; repeat the flag to AND clauses
  --archived           Include archived entities (active PLUS archived, not archived-only)
  --fields a,b,c       Project to these fields    --all-fields  Show every stored field
  --next               Dispatchable entities      --boot        Startup roll-up
  --resolve REF        Resolve slug/id/prefix     --next-id     Preview the next id
  --validate           Check workflow state       --json        Machine-readable output

Operators (one clause per flag):
  field=value   equals             field!=value  not equals
  field=        field is empty     field!=       field is non-empty

Repeat --where to AND clauses. Two clauses in one string is an error, not an AND.

Known fields are this workflow's entity frontmatter keys plus the canonical set
(id, slug, status, title, score, source, worktree, pr, started, completed,
verdict, mod-block, archived, issue). An unknown field is an error that lists them.

Examples:
  spacedock status --workflow-dir docs/dev --where status=ideation
  spacedock status --workflow-dir docs/dev --where sprint=X --where 'sprint-readiness!=defer'
  spacedock status --workflow-dir docs/dev --where sprint=X --archived
  spacedock status --workflow-dir docs/dev --where sprint=X --archived --fields slug,status,verdict,archived
```

### `skills/fo-status-viewer/SKILL.md` — the FO contract

Before (line 22): `- `--next` / `--where "pr !="` — targeted event-loop queries.`
Before (line 41): `- Archive view: `--archived`.`

After — replace line 22 with:

```markdown
- `--where <field>=<value>` — **THE entity query.** One clause per flag; repeat
  the flag to AND clauses (`--where sprint=X --where 'sprint-readiness!=defer'`).
  Two clauses in one string is an error, not an AND. `field!=` means non-empty,
  `field=` means empty. Unknown field names are a loud error listing the known
  fields. Known fields are this workflow's frontmatter keys plus the canonical
  set: `id slug status title score source worktree pr started completed verdict
  mod-block archived issue`. Never `find`/`grep` the state dir — query it.
- `--next` — dispatchable entities.
```

After — replace line 41 with:

```markdown
- Archived-inclusive view: `--archived` returns active **plus** archived, not
  archived-only. A full-sprint answer incl. done is one query:
  `--where sprint=X --archived --fields slug,status,verdict,archived`.
```

## Acceptance criteria

- **AC-1 (value / #314) — a wrong `--where` stops returning a confident wrong answer.**
  Measured baseline that can move the wrong way: the five error-worthy rows in
  the table above return **156 of 156** rows at exit 0 today (four match-all,
  one silently-empty). End state: each exits 1 with a message naming either the
  repeated-flag fix or the known fields. Counter-baseline against over-rejection:
  `--where status=ideation`, `--where sprint=0202-survey-improvements` (6),
  `--where 'pr!='`, `--where 'mod-block!='` and `--where 'next-suppressed-by = concurrency-full'`
  keep their current row counts and exit 0. Proven by tests driving the binary,
  not by prose.
- **AC-2 (full-sprint view) — one query answers a sprint completely, and stays that way.**
  `--where sprint=0202-survey-improvements --archived` returns 11 (6 active + 5
  archived) with the archived members distinguishable by `verdict`/`archived`.
  This already holds; the deliverable is the regression test that pins it — it
  must fail if `--archived` ever becomes archived-only or stops composing with
  `--where`. **No new flag ships.**
- **AC-3 (discoverability) — the two documents an FO actually reads state the query.**
  Serves AC-1 and AC-2. `spacedock status --help` exits 0 and prints help
  (today it prints the 156-row listing), containing: `--where` named as the
  entity query, the one-clause-per-flag AND rule, the canonical known-field
  list, and `--archived` described as active-plus-archived. `fo-status-viewer`
  carries the same four facts and no longer calls `--archived` an "Archive view".
- **AC-4 (validation)** — `go test ./...` and `go test ./... -race` green,
  including the pinned oracle-parity goldens; `gofmt` clean.

## Test plan

Go tests driving the binary through the existing `runNative` / `Run` harnesses.
Estimated cost: low — one new test file per package, no fixtures beyond the
existing workflow fixtures.

1. `internal/status/where_validate_test.go` (new, ~120 lines)
   - **Compound rejection, table-driven over all four shapes** (`a=1 b!=2`,
     `a!=1 b!=2`, `a=1 b=2`, `a!=1 b=2`): exit 1, stderr names the repeated-flag
     fix. *Fails if* Check A is dropped or only catches the `=`-first shape —
     the shapes return 156 rows at exit 0 today.
   - **Unknown field**: `--where 'nosuchfield!=x'` and `--where 'spint=v'` exit
     1 and stderr lists known fields. *Fails if* Check B is dropped (both exit 0
     today, returning 156 and 0 rows).
   - **No over-rejection**: `--where 'verdict!='` on an *active-only* read exits
     0. *Fails if* the schema half of the union is dropped — `verdict` is absent
     from the active corpus.
   - **Derived names accepted**: `--where 'next-suppressed-by = concurrency-full'`
     and `--where gate-eligible=true` exit 0. *Fails if* the static derived set
     is dropped.
   - **Empty workflow**: `--where sprint=X` against a zero-entity workflow exits
     0. *Fails if* the zero-entity guard is dropped.
2. `internal/status/where_archived_test.go` (new, ~40 lines) — AC-2 regression
   pin: on a fixture with both active and archived members of one sprint,
   `--where sprint=S --archived` returns active+archived and `--where sprint=S`
   returns active only. *Fails if* `--archived` becomes archived-only or stops
   composing with `--where`.
3. `internal/cli/status_help_test.go` (new, ~40 lines) — `status --help` exits 0,
   emits no entity rows, and contains the four AC-3 facts. *Fails if* the
   `wantsHelp` guard is missing (today the listing prints instead).
4. Existing suites are the over-rejection backstop: `go test ./...` plus the
   pinned parity goldens in `zz_independent_parity_test.go`.

No live workflow test is needed: every claim is command-level and observable
from exit code plus stdout/stderr bytes.

## Expected surface

| File | Change | ~LOC |
| --- | --- | --- |
| `internal/status/parse.go` | Check A + help const | +25 |
| `internal/status/parse.go` | `validateWhereFields` | +30 |
| `internal/status/handlers.go` | one call site | +5 |
| `internal/cli/help.go` | `setStatusHelp` | +40 |
| `internal/cli/cli.go` | `wantsHelp` guard + wire | +5 |
| `skills/fo-status-viewer/SKILL.md` | two stanzas | ~+12/-2 |
| 3 new test files | above | +200 |

**7 files, ~315 insertions, tolerance ±30%.**

**Declared observable-semantics changes** (the part no line count catches):

1. **Command grammar — `--where` gains two error classes.** Invocations that
   exit 0 today exit 1 after. Intended and captain-approved; nothing can depend
   on it, since the outputs being removed are the wrong answers.
2. **Command grammar — `status --help` stops listing entities and prints help.**
   A genuine behavior change to a currently-working invocation, and the one
   place a script could break.
3. **No change** to stored formats, frontmatter, authority, dispatch, gates, or
   any default/`--json`/golden output on a valid query.

## Spike determination

**No spike needed.** Every mechanism the design rests on was exercised on this
checkout during ideation rather than assumed:

1. The #314 repro reproduces here — 156/156 match-all at exit 0, both cases (and
   two further shapes not in the issue).
2. `--archived` is inclusive and composes with `--where` — 11 = 6 + 5, measured.
3. The merged view is already legible via `--fields slug,status,verdict,archived`
   — measured; no `scope` field needs exposing.
4. `wantsHelp` + `cmd.Help()` works under `DisableFlagParsing: true` — proven by
   running `spacedock new --help`, which uses the identical configuration.
5. `applyFilters` has exactly one call site, after both materializers
   (`handlers.go:493-495`) — the insertion point is unambiguous.
6. Blast radius on the existing suite is nil: an exhaustive grep of every
   `--where` in `*_test.go` yields only `status`, `next-suppressed-by`, and the
   deliberate bad inputs `statusideation` / missing-arg. The first is always in
   the corpus, the second is in the static derived set, and the last two hit
   error paths this change does not touch. The parity goldens pin no unknown-field
   case.

## Notes

- **GitHub issue #314** (`status --where: unknown field name silently matches-all`)
  — subsumed here; close when AC-1 lands. Ideation found the issue *understates*
  the bug: three compound shapes silently match-all, and one of them survives
  field-name validation.
- The three gaps share one theme: `--where` should be **robust** (errors on
  misuse), **complete** (queryable across archived), and **discoverable**
  (documented as THE query) so the FO never falls back to raw shell — kin to the
  smallest-sufficient-mechanism gate (`fo-smallest-sufficient-mechanism`). The
  smallest-sufficient reading is what deleted AC-2's flag.
- **Two unrelated data bugs found while measuring — filed as observations, not fixed here:**
  1. `docs/dev/.spacedock-state/multi-artifact-gate-prepare/index.md` has a
     frontmatter key `tatus:` (typo for `status:`), so `tatus` is a real key in
     the active corpus and would be accepted as a "known field".
  2. 12 entities under `_archive/` carry no `archived:` stamp, although
     `runArchive` (`mutate.go:444`) always writes one — suggesting another path
     archives without stamping. Worth its own entity.

## Stage Report: ideation

- DONE: Confirm and tighten AC-1: design the field-name validation fix for --where (internal/status/parse.go parseWhereFilters ~L187, applyFilters ~L215) -- exact known-field list, exact error message shape, and how the compound-in-one-string footgun (--where 'sprint=A sprint-readiness!=defer') is rejected rather than silently matching all.
  Split into two independent guards; measurement forced the split — `--where 'sprint!=A sprint-readiness!=defer'` yields the KNOWN field `sprint` and still matches all 156, so field-name validation alone cannot close #314. Known-field set, both error strings, and the single insertion point (`handlers.go:495`) are specified under "Proposed approach".
- DONE: Design AC-2: the archived-inclusive full-sprint query mechanism (e.g. --where sprint=X --include-archived or --all) -- exact flag name/shape, and how it interacts with existing --archived semantics.
  Designed to ZERO new mechanism: `--archived` is already active+archived (`discover.go:326`), measured at 11 = 6 active + 5 archived, with the archived members already distinguishable via `--fields slug,status,verdict,archived`. A new flag would be a synonym for a shipped one; AC-2 became a doc fix plus a regression pin.
- DONE: Design AC-3: exact wording to add to spacedock status --help and the FO contract documenting --where <field>=<value> as THE entity query -- the known-field list, the AND-via-repeated-flags rule, and the archived-inclusive flag.
  Full replacement text for both surfaces is in "Documentation changes". Material finding: `spacedock status --help` does not exist — it runs the 156-row listing, because `newStatusCommand` lacks the `wantsHelp` guard `newNewCommand` has. AC-3 creates the surface rather than amending it.
- DONE: Confirm no spike needed, or spike the one thing that might be non-obvious (e.g. whether reusing the existing --archived scan path composes cleanly with --where filtering, or needs a new merged-scan mechanism).
  Recorded as "no spike needed" with six mechanisms exercised on this checkout, including the named risk: `--archived` composes with `--where` with no merged-scan work (11 rows, measured). `wantsHelp` under `DisableFlagParsing: true` was proven by running `spacedock new --help`.
- DONE: Refine Acceptance criteria and Test plan only if needed; otherwise confirm as-is. Note this closes GitHub #314 -- confirm the exact repro (--where 'nosuchfield!=x' silently matching all) still reproduces on this checkout before designing the fix.
  Repro confirmed live: `--where 'nosuchfield!=x'` returns 156/156 at exit 0, as does the compound. ACs rewritten against measured baselines (AC-2 no longer asks for a flag; AC-1 gained an over-rejection counter-baseline). Test plan rewritten with the falsifying change named per test.

### Summary

The measured baseline changed two of the three ACs. AC-2's premise was wrong — `--archived` has always meant active-plus-archived, so the full-sprint query already works and the proposed `--include-archived` flag would have been command-surface sprawl; what failed was the FO contract calling it an "Archive view", which reads as archived-only. AC-1 turned out to need two guards rather than one, because one compound shape (`a!=1 b!=2`) parses to a legitimate field name and still silently matches all 156 rows — invisible without running it. AC-3 is larger than it reads: `spacedock status --help` is not an existing surface to amend, it currently prints the entity listing.

## Stage Report: implementation

- DONE: Check A rejects all four compound-in-one-string shapes with exit 1 naming the repeated-flag fix.
  `countWhereOperators` (internal/status/parse.go) counts `!=`/`=` operators via `regexp.MustCompile("!=|=")`; >1 rejects before the field/value split. Live-verified all four shapes against `docs/dev` (previously 156/156 at exit 0, now exit 1); pinned by `TestWhereCompoundShapesRejected` (4 subtests, internal/status/where_validate_test.go).
- DONE: Check B (validateWhereFields) rejects unknown/misspelled field names against the canonical known-field list with exit 1 listing known fields.
  `validateWhereFields`/`knownWhereFields` (internal/status/parse.go) union scanned-entity keys, `loadEntitySchema().fields`, and the static `derivedWhereFields` (gate-condition, gate-eligible, gate-readiness, next-suppressed-by). Live-verified `--where 'nosuchfield!=x'` and `--where 'spint=...'` against `docs/dev` (previously 156 and 0 rows at exit 0, now exit 1 each); pinned by `TestWhereUnknownFieldRejected`.
- DONE: Insert at the single call site after both materializers (handlers.go:493-495 per ideation).
  `internal/status/handlers.go`: `validateWhereFields(entities, whereFilters)` inserted between `materializeSuppressedBy` and `applyFilters`, returning `errExit` on error — 3-line diff, no other call site touched.
- DONE: Do not reject the counter-baseline cases named in AC-1.
  Live-verified against `docs/dev`: `--where status=ideation` (22 rows), `--where sprint=0202-survey-improvements` (8), `--where 'pr!='` (8), `--where 'mod-block!='` (8), `--where 'next-suppressed-by = concurrency-full'` (3) — all exit 0, unchanged. Pinned by `TestWhereNoOverRejection` and `TestWhereDerivedNamesAccepted`.
- DONE: Implement AC-2 as a doc-only change plus a regression pin; fix the FO contract's "Archive view" mislabel.
  `skills/fo-status-viewer/SKILL.md`: replaced the `--next` / `--where "pr !="` line with `--where` as THE entity query (AND rule, known-field set) and reworded "Archive view: `--archived`" to "Archived-inclusive view ... active plus archived". No flag added. Regression pin: `TestWhereArchivedComposesAsActivePlusArchived` (internal/status/where_archived_test.go) reuses the existing `enum-scope-workflow` fixture (one active + one archived entity sharing `status=backlog`, differing only by placement) — `--where status=backlog` returns active-only; `--where status=backlog --archived` returns both. Fails if `--archived` stops composing with `--where` or becomes archived-only.
- DONE: Implement AC-3 — wantsHelp guard on newStatusCommand, setStatusHelp in help.go.
  `internal/cli/cli.go` newStatusCommand gained the same `if wantsHelp(args) { return cmd.Help() }` guard `newNewCommand` uses. `internal/cli/help.go` setStatusHelp renders the synopsis (--where as THE entity query, one-clause-per-flag AND rule, canonical known-field list, --archived as active-plus-archived). Live-verified `spacedock status --help` and `-h` exit 0 and print help instead of the 159-row listing; `spacedock new --help` and top-level `--help` unaffected. Pinned by `TestStatusHelpRendersQuerySynopsisNotEntityListing` and `TestStatusShortHelpFlagAlsoRendersHelp`.
- DONE: Write the 3 new test files per the ideation's Test plan.
  `internal/status/where_validate_test.go` (143 lines): 4 compound shapes, 2 unknown-field cases, no-over-rejection (`verdict!=`), 2 derived-name cases, empty-workflow guard. `internal/status/where_archived_test.go` (51 lines): AC-2 pin. `internal/cli/status_help_test.go` (57 lines): AC-3 pin plus `-h` spelling.
- DONE: Run go test ./..., go test ./... -race, and the pinned parity goldens; gofmt -w ./cmd ./internal.
  All green: `go test ./...` (all 19 packages ok, internal/cli 118.9s, internal/status 39.2s), `go test ./... -race` (all ok, internal/status 81.4s), `go test ./internal/status -run TestInd -v` (all 14 parity tests pass, e.g. `TestIndReadFlagsSeq`, `TestIndUsageErrorsExitDomain`). `gofmt -l ./cmd ./internal` empty both before and after `-w`.
- DONE: Confirm the 2 declared observable-semantics changes land as scoped and nothing else.
  `git diff` on `internal/cli/cli.go` and `internal/status/handlers.go` shows exactly the wantsHelp guard (5-line net) and the one validateWhereFields call (3 lines) — no other command's help handling or call site touched. Full surface: 8 files, 395 insertions / 3 deletions (declared estimate 315 ± 30% = 220-410; within tolerance).
- SKIPPED: Fix the two incidental data bugs (tatus: typo, missing archived: stamps).
  Out of scope per assignment; confirmed still present and unaffected — the `tatus` key surfaces in the live known-fields error message as expected corpus noise, not a regression.

### Summary

Two independent guards close #314: an operator-count syntax check (Check A) catches all four compound-in-one-string shapes regardless of field validity, and a corpus+schema+derived known-field union (Check B) catches misspelled/unknown field names — together closing the gap a single check could not (one compound shape parses to a legitimate field name). AC-2 needed no new mechanism, only a doc fix plus a regression pin, since `--archived` was already active-plus-archived. AC-3 created a help surface that did not exist (`status --help` previously ran the entity listing). All changes were live-verified against the real `docs/dev` corpus before and after, and the full suite (including race and the pinned oracle-parity goldens) is green.
Commit: 59274d492 (branch spacedock-ensign/status-where-robust-and-discoverable)

## Review-finding disposition

- Reviewer (validation, cycle 1): a corpus-only custom field that exists
  exclusively on archived entities is falsely rejected as "unknown field" by
  Check B when queried without `--archived`, on the real `docs/dev` corpus.
  Live-verified: `docs/dev/.spacedock-state`'s own `superseded-by` field is
  written only under `_archive/` (`grep -rl superseded-by
  docs/dev/.spacedock-state` — hits only `_archive/*.md`).
  `spacedock status --workflow-dir docs/dev --where 'superseded-by!=x'`
  (no `--archived`) exits 1, `unknown field "superseded-by"`; the identical
  query with `--archived` appended succeeds. Pre-fix, the same query on the
  same corpus silently matched every active entity (the #314 pattern) rather
  than erroring, so this is not a new wrong-answer class — it swaps one wrong
  behavior (silent match-all) for a rough edge (a real field called
  "unknown" instead of "archive-scoped").
  - Released user and normal workflow: yes — `fo-status-viewer/SKILL.md`
    directs the FO to query `docs/dev` directly rather than grep it, and
    `superseded-by` is a real field the archive path writes today.
  - Observable harm: a query naming a real, currently-used field is rejected
    with a message implying a typo, with no hint to add `--archived`.
  - Affected value AC or boundary: none: AC-1's over-rejection guarantee
    enumerates five specific counter-baseline cases; an archive-only custom
    field is not among them, and no AC promises corpus-derived field names
    stay known regardless of `--archived` scope. The ideation doc's "each
    half covers the other's hole" claim holds for the *canonical* archive-only
    fields (verdict/completed/archived, covered by the schema half) but not
    for a *workflow-custom* field that is archive-only by nature.
  - Trigger evidence: reproduced live against `docs/dev` above.
  - Tentative classification: **Deferred risk.** Trigger is real but narrow
    (a custom field whose only writer is the archive path); promotes to
    Material if a documented persona is found relying on unscoped queries
    over such a field, or if a captain ruling extends AC-1's over-rejection
    guarantee beyond its five named cases.

## Stage Report: validation

- DONE: Verify AC-1: reproduce TestWhereCompoundShapesRejected (all 4 subtests) and TestWhereUnknownFieldRejected yourself against the live docs/dev corpus; independently confirm the previously-reported baseline (156/156 or similar match-all) actually flips to exit 1 with the claimed row counts, not just that the test passes in isolation.
  Built binaries from HEAD (59274d492) and from merge-base (48a7ea0d9). Baseline binary on live `docs/dev` (now 157 active rows, corpus has drifted since ideation's 156): `nosuchfield!=x`, both compound shapes, and `status!=backlog sprint=A` all return the full 159-line listing (157 rows + 2 header lines) at exit 0 — reproduces the #314 match-all; `sprint=A sprint-readiness=ready` returns 0 rows at exit 0 (silent empty). HEAD binary on the same corpus: all six shapes now exit 1 with `unknown field` or `takes one clause per flag` in stderr, stdout empty. `go test ./internal/status -run 'TestWhereCompoundShapesRejected|TestWhereUnknownFieldRejected' -v`: 6/6 subtests PASS.
- DONE: Verify AC-1's no-over-rejection guarantee: reproduce TestWhereNoOverRejection and TestWhereDerivedNamesAccepted; confirm the 5 counter-baseline cases (--where status=ideation, --where sprint=X, --where 'pr!=', --where 'mod-block!=', --where 'next-suppressed-by = concurrency-full') still exit 0 with unchanged row counts.
  Ran all 5 cases through both the baseline and HEAD binaries against live `docs/dev` and diffed line counts: status=ideation (22/22), sprint=0202-survey-improvements (8/8), pr!= (8/8), mod-block!= (8/8), next-suppressed-by = concurrency-full (3/3) — identical, both exit 0. `go test ./internal/status -run 'TestWhereNoOverRejection|TestWhereDerivedNamesAccepted|TestWhereEmptyWorkflowSkipsFieldValidation' -v`: all PASS.
- DONE: Verify AC-2: reproduce TestWhereArchivedComposesAsActivePlusArchived; confirm --archived still means active-plus-archived and composes with --where, and that the FO contract doc no longer mislabels it "Archive view".
  `go test -run TestWhereArchivedComposesAsActivePlusArchived -v`: both subtests (active-only, active-plus-archived) PASS. Live on `docs/dev`: `--where sprint=0202-survey-improvements` returns 6 data rows, adding `--archived` returns 11 — active is a strict subset, confirming composition, not a scope swap. `git diff` on `skills/fo-status-viewer/SKILL.md` confirmed the "Archive view" wording is gone, replaced by "Archived-inclusive view ... active plus archived".
- DONE: Verify AC-3: reproduce TestStatusHelpRendersQuerySynopsisNotEntityListing and TestStatusShortHelpFlagAlsoRendersHelp; confirm spacedock status --help and -h both exit 0 and print the query synopsis (not the entity listing), and that spacedock new --help and top-level --help are genuinely unaffected.
  Both tests PASS in isolation. Live: HEAD's `status --help` and `status -h` print the identical query synopsis (byte-diffed, `diff` empty) at exit 0; baseline binary's `status --help` instead ran the entity listing (confirmed: table header + rows). `new --help` and top-level `--help` on HEAD print their own unrelated synopses at exit 0, unchanged in shape from baseline.
- DONE: Confirm the surface-deviation note (8 files, 395 insertions/3 deletions vs declared 315+-30%=220-410) is within tolerance and legitimately explained -- spot check the diff yourself, don't take the count on faith.
  `git diff --numstat $(git merge-base main HEAD)..HEAD`: 8 files, exactly 395 insertions / 3 deletions, matching the report's claim. 315 ± 30% = 220.5–409.5; 395 falls inside. Read every hunk (`cli.go`, `handlers.go`, `parse.go`, `help.go`, `where_validate_test.go`, `where_archived_test.go`, `status_help_test.go`, `SKILL.md`): all lines are the declared mechanism (Check A/B, wantsHelp guard, setStatusHelp, the two doc stanzas) plus doc comments and 3 new test files — no unrelated or padding changes found.
- DONE: Confirm the SKIPPED item (the two incidental data bugs) is correctly out of scope and unaffected, not silently masking a real regression.
  Confirmed live: `docs/dev/.spacedock-state/multi-artifact-gate-prepare/index.md` still carries both `tatus: backlog` (typo) and `status: backlog`; the `tatus` key surfaces in the live known-fields error list as expected corpus noise, not a new regression. `_archive/` still contains entities without an `archived:` stamp (pre-existing, unrelated to this diff). Neither is touched by this change's files.
- DONE: Run go test ./... -count=1, go test ./... -race -count=1, and go test ./internal/status -run TestInd -v (the pinned parity goldens) yourself from scratch; confirm all green.
  All three run to completion from a clean state: `go test ./... -count=1` — 19/19 packages ok (internal/cli 113.0s, internal/status 45.2s). `go test ./... -race -count=1` — 19/19 packages ok (internal/cli 120.1s, internal/status 55.7s). `go test ./internal/status -run TestInd -v` — 8 top-level parity tests (incl. subtests) all PASS. `gofmt -l ./cmd ./internal` empty.
- DONE: Recommend PASSED or REJECTED with evidence for each AC; list deferred risks separately from material findings.
  See Summary and Review-finding disposition above.

### Summary

**PASSED.** All four ACs verified with independent, live evidence against the real `docs/dev` corpus (not just re-running the implementer's tests in isolation): AC-1's compound and unknown-field rejections flip a reproduced #314 match-all to a loud exit 1, its five named counter-baseline cases are byte-for-byte unchanged, AC-2's archive composition and contract wording are confirmed, and AC-3's help surface is proven against both the new and pre-fix binaries. The surface-deviation note checks out on inspection, not just by numstat. One new finding: an archive-only *custom* field (`superseded-by`, real and in current use in `docs/dev`) is over-rejected by Check B on an active-only query — classified as a **deferred risk**, not material, since no AC promises coverage beyond the five named counter-baseline cases and the pre-fix behavior for this exact case was the worse silent-match-all bug, not a correct result. `go test ./...`, `-race`, and the pinned parity goldens are all green from a clean build.
