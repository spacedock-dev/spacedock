---
id: pgdyphtaqfx1zn0h7ax31e5h
title: FO scheduler reads must not silently truncate at the page limit
status: validation
source: "trim-dispatch-core-stale-prose ideation and validation, 2026-08-15; captain 2026-08-15: necessary before 0.27"
started: 2026-08-15T23:04:16Z
completed:
verdict:
score: "0.95"
worktree: .worktrees/spacedock-ensign-fo-scheduler-read-pagination-truncation
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
        - id: gate:pgdyphtaqfx1zn0h7ax31e5h:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:pgdyphtaqfx1zn0h7ax31e5h-ideation-1
              briefing:
                id: briefing:pgdyphtaqfx1zn0h7ax31e5h:ideation:attempt-1:revision-1
                digest: sha256:8d0dff62a156a12a2c0242eb142ee068ce109e8fd6e8f757d47f6f1ee7f944f8
                request-digest: sha256:5fbd9ac55b422949c6147de131eb4f0989c5a91e438035b95e680c10b3a5e5fb
                room-ref: ./fo-scheduler-read-pagination-truncation/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pgdyphtaqfx1zn0h7ax31e5h:ideation:1
                briefing: briefing:pgdyphtaqfx1zn0h7ax31e5h:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T23:03:54.249344Z"
                decision: approve
                reason: 'Captain batch approval 2026-08-15 (approve all): into implementation as stack layers'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:pgdyphtaqfx1zn0h7ax31e5h:validation
          stage: validation
          attempts:
            - id: gate-attempt:pgdyphtaqfx1zn0h7ax31e5h-validation-1
              briefing:
                id: briefing:pgdyphtaqfx1zn0h7ax31e5h:validation:attempt-1:revision-1
                digest: sha256:e141176cf9bf0056cb8b542bbf1d377974b474b953e0c5770f7648ec58679c56
                request-digest: sha256:e5903caf79256e79e79be0da73f24fb51bb48409e6939d74b29709924b039a11
                room-ref: ./fo-scheduler-read-pagination-truncation/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:pgdyphtaqfx1zn0h7ax31e5h:validation:1
                briefing: briefing:pgdyphtaqfx1zn0h7ax31e5h:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T02:22:26.519366Z"
                decision: approve
                reason: 'Captain 2026-08-16 (approve all to the stack): validation PASSED; lands as a stack layer; captain reviews the PRs'
              application:
                target-stage: done
                state: pending
pr: "#712"
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

## Stage Report: implementation

- DONE: Execute the gated design exactly: FO-internal JSON reads (--where, bare status) return the complete set like --next; text-mode cap and explicit paging preserved; the 5 spiked test failures return green under the gated retargets
  Commit `f61f48c1d`, the `asJSON && !pageSet && !limitSet` guard applied verbatim at `native_runner.go:154`. The 5 spiked failures reproduced exactly (3 golden subtests + `TestStatusPaginationDefaultBounds` + `TestStatusPaginationFilteredComposition`) and returned green under the retargets and regenerated goldens. `--next --json` is byte-identical between the base and patched binaries on the live workflow and still emits no `pagination` key, so the two scheduler envelopes now agree.
- DONE: Value AC: the 30-entity fixture observes all rows (baseline 25); the live workflow read returns 172 of 172, recorded as one-off evidence
  Both measured by running the base binary built from `ab112b6d3` and the patched binary side by side. The step-2 read `--where "mod-block !=" --json --fields id,slug,mod-block` on a 30-entity mod-blocked fixture: 25 rows (`row-01`..`row-25`, `has_next:true`) before, 30 rows (`row-01`..`row-30`, `has_next:false`) after. Live `docs/dev` bare `status --json`: 25 of 172 before, 172 of 172 after. Text mode on the live workflow is byte-identical between the two binaries, footer included.
- DONE: Base is the stack top ab112b6d3; land as the next available layer per the ratified Stacked mode (gh pr create --base onto the current top, verify membership by gh stack view --json read-back, never the banner)
  PR #712, candidate `f61f48c1d`, created with `gh pr create --base spacedock-ensign/tune-dev-template-for-gated-stages` (the head of #711, itself at `ab112b6d3`). `gh pr view` returns the entity title, `isDraft:false`, and a body whose normalized sha256 equals the built `PR_BODY_FILE` (1-byte trailing-newline delta only). No `gh stack submit`, no `gh pr edit`.
- DONE: the `gh stack view --json` membership read-back required by the same item
  The join reproduced the failure mode `ab112b6d3` shipped the guard for: `gh stack link 699 ... 712` printed `⚠ failed to update base branch for PR #709 ... HTTP 422 PullRequest.base is invalid` then `✓ Updated stack to 13 PRs`, while the read-back enumerates 11 branches with #711 and #712 both absent; a narrower `gh stack link 707 712` separately claimed #712 is "already in stack #707". Per the shipped rule the banner was not trusted and #712 was not restacked: its base chain 712 -> 711 -> stack27/11 is coherent, which is the documented healthy end-state.

### Declared deviation — surface past tolerance

6 files (estimate 6, on target), +81/-16 = net +65 against the declared net +35, tolerance +/- 25 lines. Overage is +30 net and confined to `internal/status/pagination_test.go` (+69/-12 against its ~+35/-7 estimate); the other five files match the estimate line for line. Cause: the ideation costed the two JSON retargets as a `limit`-field swap, but AC-1 as written also requires `len(entities) == 30` plus boundary-slug assertions in each retargeted half, and the file's convention is a falsifiability doc comment per test (14 of the 69 added lines). The new AC-3 test was rewritten table-driven, in the file's own idiom, to cut ~20 lines without dropping an assertion. No AC was narrowed. Flagging for the gate rather than trimming assertions to fit the number.

### Summary

The guard landed as designed and the spike's prediction held exactly: 5 failures in 3 functions, all green after the two retargets and three regenerated goldens whose diff is the `limit` field alone. `go test ./...` and `go test ./... -race` are both exit 0 across 20 packages.

Every AC was proved by a demonstrated falsifying edit rather than an asserted one. Running the retargeted tests against unpatched base source reds `TestStatusPaginationDefaultBounds`, `TestStatusPaginationFilteredComposition`, and the 3 goldens, so they observe the change instead of restating it. Two wrong-guard variants were built and run to prove the boundary tests are load-bearing: `if asJSON { limit = 0 }` (ignoring explicit paging) reds `TestStatusPaginationExplicitJSONWindow`, `PageOutOfRange`, and `ArchivedComposition`, guarding AC-3; an unconditional `limit = 0` (the fix applied to the shared path) additionally reds the text halves of `DefaultBounds` and `FilteredComposition`, guarding AC-2's out-of-scope boundary.

Two items for the gate. The surface deviation above is past its declared tolerance and is the ideation's under-estimate of AC-1's own assertions, not added scope. And PR #712 is ungrouped from the stack object for the reason `ab112b6d3` documents; the FO owns the `pr:` field.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against the pg worktree commit, never by reading the report: --where and bare status JSON reads return the complete set (30-entity fixture, baseline 25); --next unchanged; text-mode cap and explicit paging preserved; the 5 retargeted tests green
  Base (ab112b6d3) and patched (f61f48c1d) binaries built and driven against a validator-built 30-entity mod-blocked fixture: step-2 read `--where "mod-block !=" --json` 25/30 base vs 30/30 patched (row-01..row-30, limit "0", has_next "false"); bare `status --json` 25 vs 30; `--next --json` byte-identical between binaries (stages fixture and live rerun) with no `pagination` key; text mode byte-identical, 25 rows, literal footer; `--json --limit 10` -> 10 rows has_next "true", `--json --page 2` -> rows 26-30, `--page 2 --limit 10` -> rows 11-20, explicit `--limit 25` still windows; the 5 retargeted tests plus TestStatusPaginationExplicitJSONWindow green in the worktree.
- DONE: One-off live evidence: the full workflow read returns every row (baseline was 25 of 172)
  Live docs/dev (171 active now after peer state movement): base bare `status --json` returns 25 of 171, patched returns 171 of 171 with count == total, base's 25 rows exactly the patched first 25, ids unique; live text mode byte-identical with footer "Showing 1-25 of 171".
- DONE: Suites green plain and -race; verdict PASSED or REJECTED with per-AC citations
  In the worktree at f61f48c1d: `go test ./...` exit 0 and `go test ./... -race` exit 0 across 20 packages; `gofmt -l ./cmd ./internal` clean. Per-AC citations under Verdict below.

### Falsifiability (independently reproduced)

Transplanting the retargeted tests onto unpatched base source reds exactly the 5 predicted failures with the truncation signature (count 25, want 30). Wrong-guard `if asJSON { limit = 0 }` reds ExplicitJSONWindow, PageOutOfRange, and ArchivedComposition (AC-3 test is load-bearing); unconditional `limit = 0` reds the text halves of DefaultBounds and FilteredComposition (AC-2 boundary test is load-bearing).

### Verdict: PASSED

- AC-1: PASS. Fixture delivers 30/30 vs baseline 25, boundary slugs row-01..row-30, on both `--where` (with and without `--fields`) and bare `status --json`.
- AC-2: PASS. Text mode byte-identical to base on the 30-row fixture and the live workflow; 25 rows with the literal footer `Showing 1-25 of 30 (page 1; use --page 2 or --limit 0 for all)`.
- AC-3: PASS. Explicit `--limit`/`--page` all window, including combined `--page 2 --limit 10` and value-equal `--limit 25`; the guard keys on flag presence (parse.go:380), so an explicit operator request always wins.
- AC-4: PASS. Full suite green plain and -race; the three golden diffs in ab112b6d3..f61f48c1d are confined to the pagination `limit` field.

No material findings. Adversarial probes beyond the ACs: empty-workflow envelope differs only in the `limit` field; `--boot`/`--read`/`--validate` `--json` byte-identical between binaries; the first live `--next` diff was concurrent peer state movement, byte-identical on rerun. Deferred risks: none new. The implementation's declared surface deviation (+30 net, confined to pagination_test.go, no AC narrowed) is confirmed against the diff and stands as a gate-visible declaration, not a validation failure.

### Summary

Re-exercised every AC against f61f48c1d with independently built base and patched binaries and a validator-built fixture, taking no evidence from the implementation report. All four ACs hold, the falsifying edits were independently reproduced (5 base-source reds plus two wrong-guard variants caught by the boundary tests), and both suites are green with gofmt clean. Recommend PASSED; the only gate-visible item is the already-declared test-file surface overage.
