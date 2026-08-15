---
id: pgdyphtaqfx1zn0h7ax31e5h
title: FO scheduler reads must not silently truncate at the page limit
status: ideation
source: "trim-dispatch-core-stale-prose ideation and validation, 2026-08-15; captain 2026-08-15: necessary before 0.27"
started:
completed:
verdict:
score: "0.95"
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:pgdyphtaqfx1zn0h7ax31e5h:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:pgdyphtaqfx1zn0h7ax31e5h-backlog-1
              briefing:
                id: briefing:pgdyphtaqfx1zn0h7ax31e5h:backlog:attempt-1:revision-1
                digest: sha256:326c4376af67c748c38eb03421f758bad071579466e0abb3a3a8fb68f78e638c
                request-digest: sha256:fac2a608f4d4ed2c05bdc0efdf39b72a2c0d897f54b237fc00d51a446afb428b
                room-ref: ./fo-scheduler-read-pagination-truncation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pgdyphtaqfx1zn0h7ax31e5h:backlog:1
                briefing: briefing:pgdyphtaqfx1zn0h7ax31e5h:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T21:24:29.265299Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: dispatch all five onto the stack tip'
              application:
                target-stage: ideation
                state: consumed
---

The FO event loop's `status --where` scheduler read returns paginated JSON with a default limit of 25, and the FO contract reads no pagination field. Past 25 matching entities the loop silently never sees the rest: work is dropped with no signal. Captain: needed before 0.27 stable.

Direction space for ideation: the binary could serve FO-internal reads unpaginated (or with an explicit --all), or the contract could instruct consuming has_next - prefer the binary owning it so no prose depends on paging mechanics. Ideation took the binary-owned direction; see Proposed approach.

## Problem

**Seed correction, confirmed against HEAD by exercising the binary.** The seed named two truncating reads. Only one truncates.

- `status --where … --json` (and bare `status --json`) run the default listing path, which calls `paginate()` at `internal/status/handlers.go:531` with `defaultPageLimit = 25` (`internal/status/format.go:91`). **These truncate.**
- `status --next --json` never paginates. `handlers.go:523-527` returns `nextJSON(entities, …)` over the whole set and emits no `pagination` object at all. **This does not truncate.**

Measured on a purpose-built 30-entity fixture: `--where "mod-block !=" --json` returned 25 of 30 rows; `--next --json` returned all 30 `dispatchable` rows with no `pagination` key. The defect is one read path, not two.

The consuming site is FO event-loop step 2 (`skills/first-officer/references/fo-dispatch-core.md:204` at the stack tip): `status --where "mod-block !=" --json --fields id,slug,mod-block`, whose instruction is "For each row, re-read the blocking mod and resume its pending action." With more than 25 mod-blocked entities, rows 26+ are never routed. No contract text anywhere in `skills/` reads `pagination` or `has_next`.

**Severity today is latent at step 2, live elsewhere.** The live `docs/dev` workflow currently has 1 mod-blocked entity, so step 2 is not dropping work right now. But the same capped path serves the general listing, where the live workflow has 172 active entities and bare `status --json` returns 25 — a 147-row silent loss for any machine reader.

**AC-2 as filed was a tautology and is replaced.** It asked that "the output names the remainder", which today's envelope already does (`has_next: "true"`, `total: "30"`). It would have passed on unmodified HEAD. The real gap is not a missing signal; it is that a machine caller gets a partial set it never asked for, and no consumer reads the signal that says so. The revised ACs below measure the delivered row set instead.

## Proposed approach

**Pagination is a human-table affordance; machine output is complete unless the caller asks for a page.** Row capping exists so the captain-facing overview stays readable — its own footer says so ("use --page 2 or --limit 0 for all"). A JSON envelope has no readability constraint. So `--json` serves the full set unless `--page`/`--limit` explicitly selects a window.

One guard in `internal/status/native_runner.go`, placed after the flags are parsed and before the existing `--page/--limit` compatibility block:

```go
// Row capping is a human-table affordance: the 25-row default keeps the
// captain-facing overview readable. A machine reader has no such constraint
// and a partial page it never asked for is silently dropped work, so --json
// serves the complete set unless the caller selects a page explicitly.
if asJSON && !pageSet && !limitSet {
    limit = 0
}
```

`limit = 0` is the already-shipped "no pagination" path (`paginate()` at `format.go:109-111`), so this reuses proven behavior rather than adding a mechanism.

**Why this and not the alternatives.** The seed's own preference (binary owns it, no prose depends on paging mechanics) rules out having the contract loop on `has_next`; that makes every future machine reader re-implement paging and leaves the defect class open. Considered and rejected:

- **`--where` implies unpaginated.** Narrower but wrong seam. Bare `status --json` is equally a machine read and would still truncate, and `--where` in text mode is still a human table that should stay capped. Splitting on the filter flag rather than the output mode fixes one instance, not the class.
- **Add an `--all` flag the contract passes.** Requires the prose the seed wants to avoid, and fails open: any future machine read that forgets `--all` silently truncates again.
- **Drop `defaultPageLimit` to 0 outright.** One-line simpler, but dumps 172 rows into the captain-facing overview. Captain-facing table rendering is out of scope, so the `asJSON` condition is what keeps the change inside the boundary.

**Consistency argument.** `--next --json` already returns its complete set. After this change `--where --json` matches it. The inconsistency between the two scheduler envelopes is exactly the bug.

**No contract prose changes.** The `pagination` object stays in the envelope (reporting `limit: "0"`, `has_next: "false"`), so the envelope shape documented at `fo-dispatch-core.md:197` remains accurate byte-for-byte, and step 2 at line 204 needs no edit. Explicit `--page`/`--limit` keep working for JSON callers, so no capability is lost.

**Documentation diff** (`docs/site/reference/command-reference.md:91`, the `spacedock status` row). Before:

> the default entity table is sorted by later workflow stage first then score descending, paginated to 25 rows by default (`--page N`, `--limit N`, `--limit 0` for all), omits the SOURCE column by default

After:

> the default entity table is sorted by later workflow stage first then score descending, paginated to 25 rows by default in the human table (`--page N`, `--limit N`, `--limit 0` for all) while `--json` returns every row unless `--page`/`--limit` selects a window, omits the SOURCE column by default

`skills/fo-status-viewer/SKILL.md:46` needs no edit: it documents the captain-facing overview ("no extra flags shows the first 25 rows … use `--page N` for more"), which is text mode and stays true.

## Out of scope

Captain-facing table rendering.

## Expected surface and tolerance

6 files, roughly +45 / -10 lines. Tolerance: +/- 2 files and +/- 25 lines.

| File | Change |
|---|---|
| `internal/status/native_runner.go` | +8 (the guard and its comment) |
| `internal/status/pagination_test.go` | ~+35 / -7 (two JSON halves retargeted, one new explicit-paging test) |
| `internal/status/testdata/golden/seq-default.json` | `"limit":"25"` -> `"limit":"0"` |
| `internal/status/testdata/golden/seq-where.json` | `"limit":"25"` -> `"limit":"0"` |
| `internal/status/testdata/golden/seq-archived.json` | `"limit":"25"` -> `"limit":"0"` |
| `docs/site/reference/command-reference.md` | 1 line |

**Observable semantics changed.** Exactly one, declared deliberately: for the default listing with `--json` and no explicit `--page`/`--limit`, the `entities` array now carries every matching row instead of the first 25, and the `pagination` object reports `limit: "0"`, `end: "{total}"`, `has_next: "false"`. Unchanged: command grammar (no flag added or removed), envelope key set, authority, text-mode rendering, `--next`/`--boot`/`--validate`/`--read` behavior, and JSON output when `--page`/`--limit` is explicit.

**Stack position.** Designed against stack tip `stack27/11-self-contained-contracts`. `internal/status/` and both target docs are untouched between `main` and the tip, and `fo-dispatch-core.md:204` is byte-identical there, so the change applies clean. Implementation lands as the next layer, `stack27/12-*`.

## Acceptance criteria

**AC-1 - A `--where … --json` read over 30 matching entities delivers all 30 rows.**
Verified by: a 30-entity fixture workflow driven through the built binary, asserting `len(entities) == 30` with boundary slugs `row-01`..`row-30`. Measured baseline on unmodified HEAD: 25. Fails if the cap survives, and fails the other way if a worker removes pagination wholesale (AC-2/AC-3 catch that).

**AC-2 - The captain-facing text table still caps at 25 with its footer intact.**
Verified by: the same fixture in text mode asserting exactly rows `row-01`..`row-25` and the literal footer `Showing 1-25 of 30 (page 1; use --page 2 or --limit 0 for all)`. This is the out-of-scope boundary made checkable; it fails if the fix is applied to the shared path instead of the JSON path.

**AC-3 - Explicit paging still wins in JSON.**
Verified by: `--json --limit 10` returns 10 rows with `has_next: "true"`, and `--json --page 2` returns rows 26-30. Fails if the guard ignores an explicit operator request.

**AC-4 - The suite stays green.**
Verified by: `go test ./internal/status/ ./internal/contractlint/` and `go test ./... -race` on the changed packages, with the three golden files regenerated and their diff confined to the `limit` field.

## Test plan

All of it is Go unit/fixture tests in `internal/status`; no live workflow run is needed, because the claim is command output shape, not runtime behavior. Cost: low, the fixture builder already exists.

- `buildPaginationFixture(t, 30)` already builds the 30-row fixture. AC-1 and AC-2 retarget the JSON half of `TestStatusPaginationDefaultBounds` and `TestStatusPaginationFilteredComposition`; both already assert the text half, which stays as the AC-2 guard.
- AC-3 adds one focused test for `--json --limit 10`. `--json --page 2` is already covered by `TestStatusPaginationArchivedComposition`, and `--limit 0 --json` by `TestStatusPaginationLimitZero`; both pass unmodified under this change.
- AC-4 regenerates the three goldens with `-update` and confirms the diff is the `limit` field only.

**Spike result (riskiest path exercised first, not asserted).** The full change was built and run in a throwaway copy of HEAD before this design was written. Baseline `go test ./internal/status/ ./internal/ensigncycle/` green; with the guard applied, exactly 5 failures in 3 functions — `TestJSONReadGolden/{default,archived,where}` (golden `limit` field), `TestStatusPaginationDefaultBounds` (`pagination_test.go:198`, JSON count 30 vs 25), `TestStatusPaginationFilteredComposition` (`pagination_test.go:334`, `end` 30 vs 25). `internal/ensigncycle` and `internal/contractlint` stayed green. Applying the two test retargets and regenerating the goldens took the package back to green. End-to-end on the patched binary: the step-2 read returned 30/30 on the fixture, and the live `docs/dev` workflow returned 172/172. Every failure the implementation will meet is therefore known and bounded in advance.

## Stage Report: ideation

- DONE: Design the no-truncation mechanism with the binary owning it if the trade study holds; 30-entity fixture value AC baseline measured (today observes 25)
  Trade study held: binary owns it via one `asJSON && !pageSet && !limitSet` guard in `native_runner.go`; three alternatives rejected in Proposed approach. Baseline measured on a built binary against a 30-entity fixture: `--where "mod-block !=" --json` returns 25/30 on HEAD, 30/30 patched.
- DONE: Declare observable semantics for the scheduler read envelopes
  Exactly one semantic changes, declared under Expected surface and tolerance: default-listing `--json` without explicit `--page`/`--limit` now returns every row, `pagination` reports `limit:"0"`/`has_next:"false"`. Envelope key set, command grammar, text mode, and `--next` are unchanged, so `fo-dispatch-core.md:197` stays accurate and needs no edit.
- DONE: Design against the stack tip; implementation lands as a stack layer
  Diffed `main..origin/stack27/11-self-contained-contracts`: `internal/status/`, `command-reference.md`, and `fo-status-viewer/SKILL.md` are untouched there, and `fo-dispatch-core.md:204` (the step-2 read) is byte-identical. Lands as `stack27/12-*`.

### Summary

Confirmed the removal scope against HEAD by exercising the binary rather than reading prose, and corrected the seed on the way: only `--where`/bare `status` truncate, while `--next --json` already returns its complete set and emits no `pagination` object. That reframed the fix as making the two scheduler envelopes consistent, and it made the filed AC-2 a tautology — today's output already names the remainder via `has_next`/`total`, so it would have passed on unmodified HEAD; the ACs now measure the delivered row set, the preserved text-mode cap, and preserved explicit paging.

The riskiest path was spiked end-to-end before the design was written: the guard was applied in a throwaway copy of HEAD, the suite run to a green baseline first, and the resulting blast radius measured at exactly 5 failures in 3 functions plus 3 golden `limit` fields, all of which returned to green after the proposed test retargets. Severity is stated honestly: step 2 is latent today at 1 mod-blocked entity, but the same capped path returns 25 of 172 rows on the live workflow.

One process note for the FO: the first spike tree was deleted mid-run by a sibling ensign sharing the job tmp directory, which produced fabricated-looking cross-package failures (fixtures "missing" that exist in the repo). That run was discarded and redone in a private `mktemp -d` tree; the reported baseline and blast radius come only from the clean tree.
