---
id: 0qdayfbj4kpb9myj7y2mcc4k
title: Make the status SOURCE column conditional — drop it from the default listing, render only when requested
status: validation
source: captain (2026-06-14) "make source field not default" — the `spacedock status` listing renders the SOURCE column UNCONDITIONALLY (internal/status/format.go:14 defaultStatusFields; :111-116 header, :136-137 rows). At boot the overview emitted ~7,200 tokens for ~30 entities, each carrying a multi-sentence SOURCE (boot-forensics, this session). A greet/overview rarely needs SOURCE. Separate from e6a/status-section-reader (which is file-BODY structural reads); this is the status TABLE column render. 0.20.4 read-cost theme.
started: 2026-06-15T05:35:19Z
completed:
verdict:
score: 0.35
worktree: .worktrees/spacedock-ensign-status-source-not-default
issue:
sprint: 0204-structured-reads
sprint-readiness: ready
mod-block: merge:pr-merge
pr: "#387"
---

The `spacedock status` listing table always renders the SOURCE column, the heaviest per-row field. Drop it from the default render so the common overview is cheap; surface it only when a caller explicitly asks. Reduces the recurring boot/overview read-cost the FO pays every session.

## Problem

`internal/status/format.go` renders SOURCE unconditionally: it is in `defaultStatusFields` (line 14) and emitted in the table header (lines 111-116) and every row (lines 136-137). SOURCE is multi-sentence provenance prose; across ~30 entities it dominated the boot overview at ~7,200 tokens (boot-forensics this session). The FO running a greet, an overview, or a "what's dispatchable" check almost never needs SOURCE — but pays for it every time. There is no current way to render the human table without it (`--fields` narrows only the `--json` read, not the human table).

## Proposed approach

### What the code does today (spike-verified, see Spike below)

SOURCE is **not** driven by `defaultStatusFields`. The human table (`printStatusTable`, format.go:102-138) hardcodes SOURCE as a sixth FIXED column in both its no-extras and extras render paths. `defaultStatusFields` (format.go:14, includes `source`) drives only two OTHER things:

1. the `--json` default key set, via `resolveJSONFields` (handlers.go:431);
2. the `displayedColumns` de-dupe set passed to `resolveExtraFields` (handlers.go:435), which is why `--fields source` today shows SOURCE exactly once (it is de-duped against the fixed column, not appended).

Consequences the spike pinned: `--fields id,slug,status,title,score` produces byte-identical output to the bare default (the fixed SOURCE column cannot be subtracted by `--fields`); `--fields source` still shows one SOURCE column (de-duped); `--all-fields` shows SOURCE in its fixed slot; `--json` default carries `source`; `--json --fields id,source` carries `source`. `--boot`'s DISPATCHABLE uses `printNextTable` (no SOURCE), so the ~7,200-token boot cost the task cites comes from the FO running a bare `spacedock status` overview, which is the surface this task targets. `--resolve` (`formatResolveLine` / `resolveJSON`) renders `workflow= scope= slug= id= path=` and has **never** carried SOURCE.

### Chosen direction

Make the human-table SOURCE column **conditional on `source` being in the resolved field set**, and reuse the existing field-projection surface as the opt-in rather than adding a new flag:

- Remove `source` from `defaultStatusFields` so SOURCE is no longer a fixed column, no longer in the de-dupe `displayedColumns`, and no longer in the `--json` default key set.
- In `printStatusTable`, render the SOURCE column only when `source` is among the columns to show. With `source` gone from the fixed set, `--fields source` (and `--all-fields`, whose scan now surfaces `source` as a non-default key) re-adds it as an extra column through the existing `resolveExtraFields` path — no new flag, no new render branch beyond the conditional.
- `--resolve` / single-entity views: **unchanged** — they never carried SOURCE, so there is nothing to drop or gate. The AC about them resolves to "no SOURCE before or after."

**Opt-in surface decision: `--fields source` (and `--all-fields`), not a new `--with-source` flag.** A dedicated flag would duplicate the field-projection mechanism that already exists and already reads naturally ("name the field you want"). YAGNI: the projection surface is the opt-in.

### Open decision for the gate — `--json` default

Dropping `source` from `defaultStatusFields` is the single shared knob; it ALSO drops `source` from the `--json` default envelope (today `seq-default.json` carries `"source":"roadmap"`). Two ways to resolve, and this is a machine-contract change the captain should ratify, not the ensign:

- **(A, recommended) `source` is opt-in everywhere.** Human default and `--json` default both omit `source`; `--fields source` / `--all-fields` restore it in both. Simplest change (one knob), uniform rule, AC-3 still holds (`source` present when named). Cost: a breaking change to the `--json` default key set for machine consumers that read `source` without naming it.
- **(B) Human default drops SOURCE; `--json` default keeps `source`.** Preserves the machine contract but splits "default fields" into two different sets (human vs JSON), needing a separate human-default list distinct from `defaultStatusFields`. More code, divergent rules.

Recommendation: **(A)**. The task frames the win as the default LISTING read-cost, and a uniform "opt-in everywhere" rule is the cleaner contract; AC-3 only ever required `source` *when named*. Flagged here because (A) is a breaking change to the `--json` default envelope — gate-level ratification, not an ensign call.

### Parity-suite handling (mechanism note — divergence touches FIVE human cases + THREE JSON invariants)

Dropping SOURCE is a **deliberate divergence from the Python oracle**, mirroring the captain-approved `--fields` de-dupe divergence already in the tree (native_read_test.go:24-28: the divergent case was pulled OUT of oracle parity and locked by a native-only test, `TestFieldsDedupeNoDuplicateDefaultColumns`, instead). The divergence surface is wider than the bare default — it touches **every case that renders the human fixed table** (which carries the SOURCE column) and **every JSON envelope that carries the default `source` key**. Implementation follows the de-dupe precedent: pull each diverging case OUT of oracle parity into native-only golden fixtures regenerated with `-update`, plus native-only assertions on the new behavior. This is recorded so the implementer does not "fix" the now-failing oracle parity by reverting the column.

**Human fixed-table cases (`TestNativeReadMatchesOracle`, native_read_test.go) — all FIVE render SOURCE today and must move to native-only:**

- `default` — bare table.
- `archived` (`--archived`) — same fixed table, archived rows appended.
- `where-status` (`--where status=ideation`) — same fixed table, filtered.
- `fields` (`--fields worktree`) — the fixed table PLUS a worktree extra; the fixed SOURCE column is still rendered, so this golden also carries SOURCE (spike-confirmed: `native-read-fields.txt` HAS a SOURCE token).
- `all-fields` (`--all-fields`) — SOURCE moves from its fixed slot to an EXTRA column (the `--all-fields` scan surfaces `source` as a non-default key); the column persists but its position/width change, so this golden changes and must be native-only too.

The other four read cases (`next`, `validate`, `resolve`, `short-id`) render no SOURCE (spike-confirmed) and STAY in oracle parity unchanged. The `ind-*-default` / `ind-*` independent-parity fixtures in `zz_independent_parity_test.go` that render the fixed table change identically and are regenerated in the same sweep.

**JSON-side invariants (under direction A, where the `--json` default also drops `source`) — THREE places, not one:**

1. **`TestJSONReadGolden` goldens carrying the default key set:** `seq-default.json`, `seq-archived.json`, AND `seq-where.json` all carry `"source":…` today (the default envelope) and must be regenerated. `seq-next.json` / `seq-resolve.json` carry no `source` and are unchanged.
2. **`TestJSONStatusRoundTripsTableColumns`** cross-walks the `--json` default entities against the default TABLE columns, iterating `defaultStatusFields` and treating `source` as the unpadded row tail (json_read_test.go:88, 156, 167-169). With `source` gone from both the JSON default and the table, this test must be updated: `defaultStatusFields` shrinks to five, the new unpadded row tail becomes `score`, and `defaultColumnWidths` (json_read_test.go:88) is re-derived accordingly. Leaving it stale makes the round-trip walk index past the end of the row.
3. **`json_commands.go` `resolveJSONFields`** keys off `defaultStatusFields` for the default/`--all-fields` key set — no code change beyond the slice edit, but its behavior shift is what items 1-2 freeze.

Under direction B (JSON default keeps `source`) only the FIVE human cases diverge; the JSON goldens and the round-trip test are unchanged, at the cost of a separate human-default list. This is the concrete extra cost of B over A.

## Out of scope

- e6a / status-section-reader's file-BODY structural section reads (the separate structured-read helper) — a different surface.
- Any other column change, reorder, or width change to the status / next / boot tables.
- The `--boot` DISPATCHABLE render (already SOURCE-free via `printNextTable`).
- Re-deriving oracle parity for the default render: the oracle (Python reference) is frozen; native deliberately diverges here, as it already does for `--fields` de-dupe.

## Acceptance criteria

Each AC names a property of the finished behavior, not a stage action, and how it is verified.

**AC-1 — Every default human fixed-table render omits SOURCE.** The fixed status table — bare default, `--archived`, and `--where …` — renders a header of `ID SLUG STATUS TITLE SCORE` with no SOURCE column and no per-row provenance prose. `--fields <extra>` likewise no longer carries the fixed SOURCE column (the extra is appended without it).
Verified by: native-only golden fixtures for the five diverging renders (regenerated `native-read-{default,archived,where-status,fields,all-fields}.txt` plus the matching `ind-*` independent-parity fixtures), each pulled out of `TestNativeReadMatchesOracle`'s oracle set; plus an assertion-style native test (`countToken(header, "SOURCE") == 0`) over the bare default, `--archived`, and `--where` headers so the property is checked, not only frozen. The `--all-fields` case asserts SOURCE appears as an EXTRA (one token), not the fixed slot.

**AC-2 — `--fields source` restores SOURCE in the human table.** `spacedock status --workflow-dir <wf> --fields source` renders a table that includes exactly one SOURCE column carrying each entity's source value.
Verified by: a native test asserting `countToken(header, "SOURCE") == 1` for `--fields source` (mirroring `fields_dedupe_test.go`'s header-token oracle), and a golden fixture freezing the restored column. A second assertion confirms `--all-fields` also surfaces a single SOURCE column (now as an extra, not the fixed slot).

**AC-3 — `--json --fields source` carries `source`; bare `--json` does not (under direction A).** A machine caller that names `source` receives `"source":…` in every entity object; a caller that does not name it receives an envelope with no `source` key.
Verified by: a native test over `--json --fields id,source` asserting each entity object contains `source`, and over bare `--json` asserting no entity object contains a `source` key (regenerated `seq-default.json`). (If the gate selects direction B, AC-3 is narrowed to only the `--fields source` half and bare `--json` retains `source`.)

## Test plan

- **Default-render goldens (AC-1) — FIVE cases, not one:** pull the `default`, `archived`, `where-status`, `fields`, and `all-fields` cases out of `TestNativeReadMatchesOracle`'s oracle-parity set into native-only coverage (exactly as the `--fields` de-dupe divergence was handled, native_read_test.go:24-28), regenerate their goldens (`native-read-{default,archived,where-status,fields,all-fields}.txt`) plus the fixed-table `ind-*` independent-parity fixtures with `-update`, then add a native-only assertion test (no oracle) checking the bare-default / `--archived` / `--where` headers carry no `SOURCE` token and `--all-fields` carries SOURCE as an extra. Cost: low-moderate — fixture regen across five cases + their `ind-*` siblings + one assertion test. The render is deterministic over a fixture workflow; no live run.
- **Opt-in render (AC-2):** extend `fields_dedupe_test.go`'s header-token style with `--fields source` and `--all-fields` cases asserting exactly one `SOURCE` column; freeze a golden for the restored render. Cost: low, fixture-level.
- **JSON envelope (AC-3) — widen beyond `seq-default.json`:** under direction A, regenerate the THREE JSON goldens that carry the default `source` key — `seq-default.json`, `seq-archived.json`, `seq-where.json` — and add native assertions over bare `--json` (no `source` key), `--json --archived`, `--json --where` (no `source` key), and `--json --fields id,source` (`source` present). ALSO update `TestJSONStatusRoundTripsTableColumns`: shrink `defaultStatusFields` handling to five columns, re-derive `defaultColumnWidths` (json_read_test.go:88) with `score` as the new unpadded row tail, so the cursor walk does not index past the row. `seq-next.json` / `seq-resolve.json` are unchanged (no `source`). Cost: low-moderate — three golden regens + the round-trip-test edit.
- **Regression sweep:** `go test ./internal/status/...` to confirm no OTHER golden (boot, next, resolve, short-id) silently shifted; any fixture that legitimately changes is regenerated and reviewed in the diff. The `--resolve`, `--next`, and `--validate` goldens MUST be unchanged (none carried SOURCE) — a changed one of these is a red flag, not a regen target.

No spike beyond the one below is needed: the render is a pure function over a fixture workflow exercised through the real binary; all three ACs are fixture/CLI tests, no live workflow.

## Spike (riskiest mechanism — exercised before design)

The riskiest assumption was *where* SOURCE comes from and whether `--fields` already gates it. Built the binary (`go build ./cmd/spacedock`) and ran it against `docs/dev`:

- `--fields id,slug,status,title,score` → byte-identical (19207 bytes) to the bare default → **proves `--fields` cannot subtract the fixed SOURCE column**; SOURCE is hardcoded in `printStatusTable`, not field-driven.
- `--fields id,source` → header shows SOURCE once → **proves the de-dupe against `displayedColumns`**; so `source` must leave `displayedColumns` (i.e. leave `defaultStatusFields`) for `--fields source` to re-add it as an extra.
- `--all-fields` → SOURCE present in its fixed slot today; after removing it from `defaultStatusFields` the `--all-fields` scan surfaces `source` as a non-default key, so it still renders (as an extra).
- bare `--json` → carries `"source":…`; `--json --fields id,source` → carries `"source":…` → both pinned, framing the AC-3 / direction-A-vs-B decision.
- `--boot` DISPATCHABLE uses `printNextTable` (no SOURCE) → the boot read-cost is the bare `status` overview, the surface this task targets.

Result: the design rests on removing `source` from `defaultStatusFields` + a single conditional in `printStatusTable`, with the opt-in carried by the existing `--fields` / `--all-fields` projection — all exercised, no unverified mechanism remains.

## Documentation diff

The binary emits no detailed `spacedock status --help` flag list (`newStatusCommand` has `DisableFlagParsing: true` and only a one-line `Short`); the command-reference defers to `--help`. No user-facing doc currently describes the SOURCE column or `--fields` behavior in prose (`grep` over `docs/site/` finds no SOURCE/`--fields`/`--all-fields` description outside this state file). So the only user-facing doc surface is the `spacedock status` row in `docs/site/reference/command-reference.md`, which lists the read modes.

Before (`docs/site/reference/command-reference.md`, the `spacedock status` row):

> | `spacedock status` | Read or mutate the state: the entity table, `--next`, `--where`, `--set`, `--validate`, `--boot` |

After:

> | `spacedock status` | Read or mutate the state: the entity table (omits the SOURCE column by default; `--fields source` or `--all-fields` restores it), `--next`, `--where`, `--set`, `--validate`, `--boot` |

No other doc file changes. If the gate selects direction B (JSON default keeps `source`), the parenthetical narrows to "the human table omits the SOURCE column by default" to avoid implying the `--json` default also drops it.

## Stage Report: ideation

- DONE: Decide the opt-in surface (the flag that brings SOURCE back) and whether `--resolve`/single-entity views still carry SOURCE; pin behavior-first ACs over the default-vs-opt-in human render with golden fixtures (default omits SOURCE, opt-in restores it, `--json --fields` still carries source when named).
  Opt-in = existing `--fields source` / `--all-fields` (no new flag — YAGNI; spike proved `source` must leave `defaultStatusFields` so it re-adds as an extra). `--resolve`/single-entity never carried SOURCE → unchanged. AC-1/2/3 rewritten behavior-first with native header-token assertions + golden fixtures (Acceptance criteria section).
- DONE: Record the user-facing doc-diff: the default listing drops the SOURCE column plus the opt-in flag that restores it (before/after wording for the CLI docs).
  Documentation diff section: before/after for the `spacedock status` row in `docs/site/reference/command-reference.md`; verified by grep that no other doc describes SOURCE/`--fields` in prose and the binary emits no detailed status `--help`.

### Summary

Reverse-engineered the actual render path and exercised the binary: SOURCE is a HARDCODED fixed column in `printStatusTable`, not field-driven; `defaultStatusFields` only governs the `--json` default key set and the `--fields` de-dupe set. The design removes `source` from `defaultStatusFields` plus one conditional in `printStatusTable`, reusing `--fields source`/`--all-fields` as the opt-in. One genuine fork is surfaced for the gate (open decision, not an ensign call): whether the `--json` DEFAULT also drops `source` (direction A, recommended, uniform but a machine-contract change) or keeps it (direction B, more code). Also flagged the parity-suite handling: dropping SOURCE is a deliberate native-vs-oracle divergence following the captain-approved `--fields` de-dupe precedent.

## Stage Report: ideation (cycle 2)

- DONE: Fold the staff-review Material item — oracle-parity pull-out was under-scoped (default only) and the JSON test plan was ~7x too narrow (`seq-default.json` only).
  Widened to the verified divergence surface. Spike confirmed which goldens carry a SOURCE token: the human fixed-table cases are FIVE, not four — `default`, `archived`, `where-status`, `fields` (`--fields worktree` still renders the fixed SOURCE column), and `all-fields` (SOURCE shifts to an extra). `next`/`validate`/`resolve`/`short-id` carry no SOURCE and stay in oracle parity.
- DONE: Widen the JSON test plan beyond `seq-default.json`.
  Three JSON goldens carry the default `source` key and must be regenerated under direction A — `seq-default.json`, `seq-archived.json`, `seq-where.json` (`seq-next.json`/`seq-resolve.json` are clean). Also caught a THIRD JSON invariant the review did not name: `TestJSONStatusRoundTripsTableColumns` iterates `defaultStatusFields` and treats `source` as the unpadded row tail (json_read_test.go:88,156,167-169) — it must be updated so `score` becomes the new tail and `defaultColumnWidths` is re-derived, else the cursor walk indexes past the row.

### Summary

Folded the staff-review parity-divergence finding and corrected the scope two ways: the human-render divergence is five oracle-parity cases (the review named four; the `--fields worktree` case also renders the fixed SOURCE column), and the JSON divergence is three goldens plus the round-trip test (the review named the goldens; I added the round-trip-test/`defaultColumnWidths` edit it implies). Body's "Parity-suite handling" section, AC-1, and the test plan now enumerate all five human cases and all three JSON invariants, with the unchanged cases (`next`/`validate`/`resolve`/`short-id`) named as red-flag-if-changed. The minimal one-line drop + `--fields source` opt-in and the direction-A `--json` gate fork are unchanged — both still stand.

## Stage Report: implementation

- DONE: Drop `source` from defaultStatusFields so the human fixed table (default, --archived, --where) AND the bare --json envelope both omit source, while --fields source / --all-fields restore it through the existing field-projection path (no new flag). Direction A is gate-ratified: the --json default drops source.
  format.go:14 `defaultStatusFields` now `[id slug status title score]`; `printStatusTable` SOURCE fixed column removed from both render paths (commit 6b9cd91e). `resolveJSONFields`/`resolveExtraFields` key off the shrunk slice with no further code change — `--fields source` and `--all-fields` re-add SOURCE as an extra. Verified live: bare default/`--json` omit SOURCE/`source`; `--fields source` and `--all-fields` show one SOURCE column.
- DONE: Follow the --fields de-dupe divergence precedent: pull each of the five diverging human-table cases and the three JSON invariants OUT of oracle parity into native-only goldens (regenerated with -update), plus assertion-style native header-token tests.
  Goldens regenerated via `go test ./internal/status/ -update`: the five human cases (`native-read-{default,archived,where-status,fields,all-fields}` — all-fields byte-identical since SOURCE was already the last column at the same offset) + their `ind-*` siblings + the three JSON goldens (`seq-{default,archived,where}.json` lose `source`) + collateral status-table goldens (next-suppressed-*, enum-scope-*, worktree-overlay-*, comment-roundtrip-fields-read, argv-unknown-toplevel, merge-sentinel-displays-local, native-folder-report-append, native-slug-default). Red-flag check: `seq-next/resolve/validate/short-id` and boot goldens UNCHANGED. New `source_conditional_test.go` asserts countToken(header,"SOURCE")==0 for default/--archived/--where, ==1 for --fields source and --all-fields, and the bare-`--json`-no-`source` / `--json --fields id,source`-has-`source` JSON property. `TestJSONStatusRoundTripsTableColumns` `defaultColumnWidths` re-derived to `[6 30 20 30]` (score is the new unpadded row tail).
- DONE: AC-1, AC-2, and AC-3 green per the entity's Acceptance criteria under Direction A.
  `go test ./internal/status/` green (8.9s); `go test ./...` green (OVERALL_EXIT=0, 15 packages, no FAIL); `go vet ./...` clean. AC-1 = TestSourceOmittedFromDefaultTables + regenerated goldens; AC-2 = TestFieldsSourceRestoresColumn + TestAllFieldsSurfacesSourceAsExtra; AC-3 = TestBareJSONOmitsSource.

### Summary

Removed `source` from `defaultStatusFields` and the hardcoded SOURCE fixed column from `printStatusTable` — the single-knob change the gate ratified (Direction A). The default human table (bare/`--archived`/`--where`) and the bare `--json` envelope now omit SOURCE; `--fields source` and `--all-fields` restore it through the existing field-projection path (no new flag). Followed the `--fields` de-dupe precedent: regenerated the diverging goldens and locked the new behavior with assertion-style native tests rather than reverting the column. One behavioral note for validation: opt-in SOURCE now renders through the extras path, which truncates values >20 runes with an ellipsis (e.g. a `#`-bearing source shows `consolidates #223, …` in the `--fields source` table) — the old fixed column rendered the full value. On-disk and `--fields`-read round-trip of the raw value are unaffected (TestCommentValueRoundTripParity still green); this is inherent to the YAGNI decision to reuse the extras render path as the opt-in surface.

## Stage Report: validation

- DONE: Live behavioral confirmation of Direction A: build the binary from THIS branch (go build -o /tmp/0q-sd ./cmd/spacedock), then drive `status` across default / --archived / --where / --fields source / --all-fields / bare --json / --json --fields source against a real workflow and confirm SOURCE/`source` presence matches AC-1/AC-2/AC-3 — absent by default everywhere, present only when named. Capture the actual headers/envelopes.
  Built `/tmp/0q-sd` from branch HEAD 6b9cd91e; drove against the real `docs/dev` workflow. Default/`--archived`/`--where` header = `ID SLUG STATUS TITLE SCORE` (countToken SOURCE=0). `--fields id,slug,status,title,score` and `--fields worktree` also carry no SOURCE (AC-1 holds: `--fields` can no longer subtract a fixed column). `--fields source` appends exactly one SOURCE extra after SCORE carrying real values; `--all-fields` surfaces SOURCE once as an extra (alphabetical among GROUP/MOD-BLOCK/PR/SOURCE/…), not a fixed slot (AC-2). Bare `--json` first-entity keys `[id score slug status title]`, no `source`; `--json --fields id,source` carries `source` with the FULL 447-rune value, no ellipsis (AC-3). `--next`/`--boot`/`--resolve`(full+short-id) carry zero SOURCE — unchanged.
- DONE: Detached adversarial audit (status paths = high-stakes surface), read-only on a throwaway checkout: confirm the oracle-parity divergence is correct — the four cases kept in parity (next/resolve/validate/short-id) and the boot goldens are truly UNCHANGED and still match the Python oracle, and the TestJSONStatusRoundTripsTableColumns defaultColumnWidths re-derivation ([6 30 20 30], score as the new row tail) has no off-by-one; hunt for any status path that still leaks SOURCE or wrongly drops it. Adjudicate the opt-in SOURCE >20-rune truncation: confirm it is consistent with how ALL --fields extras render (not a source-specific regression) and that the full value remains available via `--json --fields source` (no truncation) — if the human opt-in is materially useless beyond the extras norm, route back. "Refuted nothing material" is valid.
  Throwaway detached worktree at 6b9cd91e (`/tmp/0q-audit`, removed after). Refuted nothing material. Findings: (1) The Python oracle is RETIRED — `TestNativeReadMatchesOracle` compares against frozen goldens via `assertEnvelopeGolden`, not a live oracle (harness comments: "the retired Python oracle"). So "still matches the oracle" = "golden byte-unchanged + suite green against it"; I do NOT claim live-oracle parity. (2) The commit touched NONE of the red-flag oracle/JSON/boot goldens (seq-next/seq-resolve/seq-validate/native-read-{next,resolve,validate,short-id}/boot-structural/native-boot-sdb32/ind-*-{next,resolve,validate,boot,short-id}); the four `next-suppressed-*.txt` it DID change are the status-table (printStatusTable) render, correctly dropping SOURCE — not the DISPATCHABLE printNextTable. (3) `go test ./internal/status -update` on the throwaway produced ZERO golden diffs → no stale fixture; committed goldens == current binary. (4) Off-by-one: perturbing `defaultColumnWidths[3]` 30→29 REDS TestJSONStatusRoundTripsTableColumns with "score not at col 4" — the walk genuinely pins each value to its boundary; `[6 30 20 30]` is correct (score is the unpadded tail per format.go:107). (5) Leak hunt: only residual `source` is `discover.go:145 defaultEntityKeys`, which is the parse-time setdefault backfill (every entity gets `source:""`) — correct and required so `--fields source`/`--all-fields` find the key; NOT a render path. AC tests non-tautological: re-adding source to defaultStatusFields REDS TestBareJSONOmitsSource; re-adding the hardcoded SOURCE column REDS TestSourceOmittedFromDefaultTables (all 3 subcases). Truncation adjudication: ACCEPTABLE YAGNI — `formatExtraCell` truncates EVERY extra to 20 runes with U+2026 (uniform, position-independent, mirrors the oracle's format_extra_cell; `--all-fields` emits 61 ellipsis glyphs across all extra columns), so SOURCE is not special-cased; the full untruncated value is available via `--json --fields source` (447 runes, no ellipsis). Not a regression beyond the extras norm → no route-back.
- DONE: AC cross-check: AC-1, AC-2, AC-3 each carry reproduced evidence after this stage.
  AC-1: default/`--archived`/`--where`/`--fields worktree`/`--fields id,…,score` headers all SOURCE-free (live drive, countToken=0); regenerated native-only goldens + TestSourceOmittedFromDefaultTables (perturbation-verified non-tautological). AC-2: `--fields source` and `--all-fields` each render exactly one SOURCE extra with real values (live drive + TestFieldsSourceRestoresColumn + TestAllFieldsSurfacesSourceAsExtra). AC-3: bare `--json` no `source` key, `--json --fields id,source` carries full `source` (live drive + TestBareJSONOmitsSource, perturbation-verified). `go test ./internal/status/` green (9.0s); `go test ./...` green (RAW_EXIT=0, 15 packages, 0 FAIL).

### Summary

PASSED. Direction A is live-confirmed against the real `docs/dev` workflow with this branch's binary: SOURCE/`source` is absent from every default human table and the bare `--json` envelope, present only when named via `--fields source` or `--all-fields` (human) / `--json --fields source` (machine). The detached adversarial audit refuted nothing material: red-flag goldens (next/resolve/validate/short-id/boot) are byte-unchanged, `-update` shows zero stale fixtures, the round-trip `defaultColumnWidths` re-derivation has no off-by-one (perturbation reds it), and the AC tests are non-tautological. One framing correction for the record: the Python oracle is retired, so parity is proven against frozen goldens + a green suite, not a live oracle — the report does not claim otherwise. The implementer's flagged >20-rune truncation of opt-in SOURCE is an acceptable YAGNI tradeoff, not a regression: it is the uniform `formatExtraCell` behavior for ALL extras, and the full value remains available untruncated via `--json --fields source`. No route-back.
