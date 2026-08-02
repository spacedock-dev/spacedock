---
title: status pagination and stage-then-score default sorting
status: validation
score: 0.70
id: rwpe45pdxffk2zfy24ejde6a
started: 2026-07-22T06:31:16Z
gates:
    version: 1
    current:
        gate: gate:rwpe45pdxffk2zfy24ejde6a:ideation
    records:
        - id: gate:rwpe45pdxffk2zfy24ejde6a:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:rwpe45pdxffk2zfy24ejde6a-ideation-1
              briefing:
                id: briefing:rwpe45pdxffk2zfy24ejde6a:ideation:attempt-1:revision-1
                digest: sha256:67ab918283c341a88f4b851da0a646f4de209766c5ef3e8a68485f3b10162fc0
                digest-domain: canonical-bytes
                request-digest: sha256:f602a1173a15ad93fcd313aecf4b898e97cdb6bdef2f39dc84bd4088b7d04f32
                room-ref: ./status-pagination-and-default-sorting/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rwpe45pdxffk2zfy24ejde6a:ideation:1
                briefing: briefing:rwpe45pdxffk2zfy24ejde6a:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T16:01:41.91137Z"
                decision: approve
                reason: 'Captain directed in chat: ''dispatch both'' (status-pagination-and-default-sorting + status-where-robust-and-discoverable). Ideation already complete with no spike needed; proceed to implementation.'
              application:
                target-stage: implementation
                state: consumed
worktree: .worktrees/spacedock-ensign-status-pagination-and-default-sorting
---

### Goal
Update `spacedock status` so output is paginated and sorted by stage (closer to the end first: e.g. validation, implementation, ideation, backlog) then score descending by default.

### Problem statement
Large workflows make the default `spacedock status` overview expensive to read and expensive to paste into a captain-facing turn. The current native table sorts earlier workflow stages first (`backlog` before `ideation` before `implementation` before `validation`) and then score descending, so the items closest to completion are often buried after the least-actionable backlog rows.

The new default should make the overview a triage surface: show later-stage active work first, order peers within the same stage by score descending, and bound the first page so a large workflow does not flood the terminal or an agent context. Existing focused queries (`--where`, `--next`, `--boot`, `--read`, `--resolve`, mutation flags) should retain their specialized behavior unless this task explicitly routes through the default status listing.

### Proposed approach
1. Change the default status listing sort key from `(stage order ascending, score descending)` to `(stage order descending, score descending, stable discovery tie-break)`. In a workflow whose README stages are `backlog`, `ideation`, `implementation`, `validation`, terminal/archive stages, this renders `validation` before `implementation` before `ideation` before `backlog`. Unknown stages remain grouped after known stages so malformed/unmodeled rows are visible but do not displace modeled active work.
2. Add explicit pagination flags for the default status listing: `--page N` and `--limit N`.
   - `--limit` is the maximum number of rows emitted for a page. `--limit 0` disables pagination and emits all matching rows.
   - `--page` is 1-based and defaults to `1`; with no explicit `--limit`, it uses the default limit of 25 and selects that page slice after filtering and sorting.
   - Invalid values (`--page 0`, negative or non-integer page/limit) and contradictory values (`--page N` with `--limit 0`) fail with the existing `Error: ...` / exit-1 shape.
3. Default pagination behavior: a bare default `spacedock status` and a default-listing `--where` use `page=1`, `limit=25`. This keeps the captain-facing overview bounded without changing the meaning of filters. Operators can pass `--limit 0` for the full legacy-size table.
4. Apply pagination only to the default status listing (`spacedock status`, `--where`, `--archived`, `--fields`, `--all-fields`, `--json` in the status-envelope path). Do not paginate `--next`, `--boot`, `--validate`, `--read`, `--resolve`, `--short-id`, `--next-id`, `--set`, or `--archive`.
5. Emit a short footer for paginated text tables when additional rows exist: `Showing 1-25 of 83 (page 1; use --page 2 or --limit 0 for all)`. In `--quiet`, omit the footer so data-only row output remains script-friendly. In `--json`, add a `pagination` object beside `entities` with string-valued fields (`page`, `limit`, `total`, `start`, `end`, `has_next`) to preserve the status JSON all-strings convention while making the slice auditable.
6. Update shell completions to include `--page` and `--limit` in status flags.

### Mechanism rationale
- Value AC served: AC-1 and AC-2. Explicit flags beat interactive terminal paging because current `status` is frequently consumed by agents and scripts where an interactive pager would block or hide bytes. The simplest alternative, piping through `less` or advising `head`, does not help JSON, does not expose totals, and leaves agents to implement inconsistent slicing.
- Value AC served: AC-1. Reversing the stage-order comparator in the shared status sort helper is smaller than adding a separate named sort mode because the task asks for a default behavior change, not a user-selectable sort matrix.
- Value AC served: AC-2 and AC-3. Pagination after filtering and sorting is simpler and more predictable than paginating raw scanned entities: page membership then matches exactly what the user sees, and tests can assert stable page windows.

### Acceptance criteria
1. A default active `spacedock status` listing orders known stages closest-to-terminal first, with score descending inside each stage and stable slug/discovery order for exact ties.
   - Tested by a Go status fixture containing at least two entities in each of `backlog`, `ideation`, `implementation`, and `validation`, including score ties. The assertion fails if `backlog` appears before `validation`, if lower-score peers appear first within a stage, or if equal-key rows reorder non-deterministically.
2. Default status output is bounded to the first 25 sorted rows, and `--page` / `--limit` select deterministic windows after all `--where` filters and archive inclusion are applied.
   - Tested by a fixture with more than 25 active rows and a second filtered/archived case. Assertions check visible row IDs, omitted row IDs, footer totals in text mode, and the JSON `pagination` object.
3. `--limit 0` preserves the full default-listing row set while retaining the new stage-then-score ordering.
   - Tested by running the same >25-row fixture with `--limit 0` and asserting every matching row appears and no pagination footer/JSON truncation is applied; a companion `--page 2 --limit 0` case asserts the contradictory combination is rejected.
4. Pagination flags are rejected on non-listing status modes where slicing would change command semantics.
   - Tested with representative commands (`--next`, `--boot`, `--read`, `--resolve`, and `--set`) asserting exit 1 and a diagnostic naming `--page/--limit` incompatibility.
5. User-visible command reference and FO status-viewer instructions describe the bounded default and the full-list escape hatch.
   - Tested by a documentation diff review in implementation plus existing completion/help tests extended to include `--page` and `--limit`.

### Expected surface
- `internal/status/format.go`: sort comparator and table pagination helper, approximately 40-80 LOC.
- `internal/status/native_runner.go` / `internal/status/handlers.go`: parse and route `--page` / `--limit`, incompatibility checks, JSON metadata, approximately 60-120 LOC.
- `internal/status/json_commands.go`: status envelope pagination object, approximately 20-40 LOC.
- `internal/status/*_test.go` and `internal/status/testdata/`: fixture/golden updates for sorting and pagination, approximately 150-250 LOC including fixtures.
- `internal/cli/cli.go`: completion flag list update, approximately 2 LOC.
- `docs/site/reference/command-reference.md` and `skills/fo-status-viewer/SKILL.md`: user-facing wording, approximately 10-20 LOC.

Tolerance: if implementation can reuse existing helpers, the code delta may be smaller; if golden fixtures are table-heavy, testdata may exceed the estimate. Any change outside `internal/status`, `internal/cli`, `docs/site/reference`, or `skills/fo-status-viewer` should be called out in the implementation report.

### Spike determination
No spike needed: the design relies on mechanisms already present in the native runner — argv flag parsing (`parseSingleArg`, `contains`), workflow stage parsing (`parseStagesBlock` / `stageOrder`), stable sorting (`sort.SliceStable`), text table rendering, JSON status envelopes, and fixture-backed CLI tests. The only new logic is integer flag validation and deterministic slicing after an already-materialized entity slice.

### Test plan
- Add focused unit tests around a pagination helper: zero/one-based boundaries, `--limit 0`, page past the end (valid empty page with total metadata), and invalid values.
- Add a default listing fixture under `internal/status/testdata/` with staged rows and scores that prove the comparator: `validation 0.1` must precede `implementation 0.9`; `implementation 0.9` must precede `implementation 0.4`; tied implementation rows retain discovery order.
- Add CLI/native runner tests for text output and JSON output. These should parse rows/JSON instead of relying only on broad substrings, and each assertion should name the row order that would fail if sorting/pagination is applied before filtering or in the old stage direction.
- Update relevant golden outputs whose default status order or row count changes. Keep `--next` goldens unchanged except for explicit incompatibility coverage when pagination flags are present.
- Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` before implementation completion.

### Documentation diff to apply during implementation

`docs/site/reference/command-reference.md` status row should change from:

```markdown
| `spacedock status` | Read or mutate the state: the entity table (omits the SOURCE column by default; `--fields source` or `--all-fields` restores it), `--next`, `--where`, `--set`, `--validate`, `--boot` (with `--identify`, the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy, and reports the boot sections; PR_STATE is a local `pr:` view, live PR state is checked at engage; local reads only, no mutation), `--read <ref-or-path>` (a file's structured body/frontmatter read), `--resolve REF`, and related helpers. |
```

to:

```markdown
| `spacedock status` | Read or mutate the state: the default entity table is sorted by later workflow stage first then score descending, paginated to 25 rows by default (`--page N`, `--limit N`, `--limit 0` for all), omits the SOURCE column by default (`--fields source` or `--all-fields` restores it), `--next`, `--where`, `--set`, `--validate`, `--boot` (with `--identify`, the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy, and reports the boot sections; PR_STATE is a local `pr:` view, live PR state is checked at engage; local reads only, no mutation), `--read <ref-or-path>` (a file's structured body/frontmatter read), `--resolve REF`, and related helpers. |
```

`skills/fo-status-viewer/SKILL.md` invocation should change from:

```markdown
${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} [--next-id|--next|--archived|--where ...|--boot|--validate|--resolve REF]
```

to:

```markdown
${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} [--page N|--limit N|--next-id|--next|--archived|--where ...|--boot|--validate|--resolve REF]
```

and the mode list should gain:

```markdown
- Overview: no extra flags shows the first 25 rows, sorted by later stage first then score descending; use `--page N` for more or `--limit 0` for the full table.
```

## Stage Report: ideation

- DONE: Stage-then-score default sorting implementation plan/design in entity body.
  See Proposed approach items 1 and 4 plus Acceptance criteria 1 for the comparator, tie-break, and scope.
- DONE: CLI pagination flags (--page, --limit) or default pagination behavior defined in design.
  See Proposed approach items 2, 3, and 5 plus Acceptance criteria 2-4 for defaults, flag semantics, and incompatibilities.
- DONE: Tests and fixture expectations identified for stage-then-score and pagination.
  See Test plan and Acceptance criteria; assertions name row-order/window changes that would fail under old ordering or pre-filter slicing.

### Summary
Expanded the seed into an implementation-ready design for bounded default status output with later-stage-first sorting and explicit pagination controls. The plan scopes pagination to the default listing path, identifies tests/fixtures/docs, and records that no spike is needed because existing native runner mechanisms cover the risky plumbing.

## Stage Report: implementation

- DONE: Reverse the default sort to (stage order descending, score descending, discovery tie-break); unknown stages grouped after known.
  internal/status/format.go sortDefault (dfaf1519c); internal/status/sort_default_test.go TestSortDefaultStageThenScore proves validation-lo(0.10) < impl-hi(0.90) < tied impl-lo-a/impl-lo-b (discovery order) < ideation < backlog.
- DONE: Add --page N / --limit N to the default listing (page 1-based default 1, limit default 25, --limit 0 disables); reject invalid/contradictory values.
  internal/status/parse.go parsePageLimitArgs; TestParsePageLimitArgs covers page 0/negative/non-integer, limit negative/non-integer, missing args, and --page-with-limit-0 (rejected even at --page 1, since an explicit page makes no sense once --limit 0 asks for everything on one page).
- DONE: Apply pagination only to the default listing path; reject --page/--limit on the nine non-listing modes with a named diagnostic.
  internal/status/native_runner.go dispatch() incompatibility block; native_usage_test.go usageCases page-with-{next,boot,resolve,short-id,next-id,archive}/limit-with-{validate,read,set} plus page-with-limit-zero, each with a captured golden (e.g. testdata/golden/native-usage-page-with-next.txt: "Error: --page/--limit cannot be combined with --next").
- DONE: Text-mode "Showing X-Y of Z (page N; use --page N+1 or --limit 0 for all)" footer (omitted in --quiet); JSON pagination object (page/limit/total/start/end/has_next, all strings) beside entities.
  internal/status/format.go paginationFooter/printPaginationFooter; internal/status/json_commands.go paginationJSONObj; TestStatusPaginationDefaultBounds/Page2/PageOutOfRange/LimitZero in pagination_test.go assert the exact footer text, JSON field values, and the valid-empty-page-past-the-end case.
- DONE: Update shell completions to include --page and --limit in status flags.
  internal/cli/cli.go bashCompletion/zshCompletion status_flags; internal/cli/verbs_test.go TestCompletionShells extended to assert both flags appear.
- DONE: Apply the exact documentation diff to command-reference.md and fo-status-viewer/SKILL.md.
  docs/site/reference/command-reference.md status row and skills/fo-status-viewer/SKILL.md invocation line + Overview bullet, matching the entity's diff verbatim (the command-reference "from" text had since gained an unrelated --boot/--read parenthetical upstream of ideation's capture; only the targeted substring was replaced, preserving that text).
- DONE: Tests per the Test plan (pagination-helper unit tests, stage-then-score fixture, CLI/JSON row-order tests, golden updates, --next goldens unchanged except incompatibility).
  TestPaginate/TestPaginationFooter/TestParsePageLimitArgs (unit); TestSortDefaultStageThenScore (comparator fixture); TestStatusPagination* (CLI/JSON, parses rows via splitTableAndFooter and JSON via encoding/json, not substrings); 30 golden files regenerated via -update, reviewed diff-by-diff (git diff) to confirm each changed line is only the expected reorder or the added pagination object; seq-next.txt/seq-next.json/boot-structural.txt untouched.
- DONE: go test ./..., go test ./... -race, gofmt -w ./cmd ./internal all pass.
  Full suite green (internal/status 22.2s plain / 49.6s -race, whole repo including skills/integration); gofmt -l ./cmd ./internal empty after -w.

### Summary
Reversed the default status sort to later-stage-first/score-descending with unknown stages grouped last, and added --page/--limit pagination (default page=1 limit=25, --limit 0 for all) to the default listing only, refusing the flag on every other read/mutation mode and on the --page-with-limit-0 contradiction. Text and JSON both surface the pagination window (footer / pagination object); completions and the two named docs were updated per the entity's diff. All work is on branch spacedock-ensign/status-pagination-and-default-sorting at dfaf1519c, scoped to internal/status, internal/cli, docs/site/reference, and skills/fo-status-viewer as estimated — no surface deviation to report.
