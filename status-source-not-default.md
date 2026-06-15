---
id: 0qdayfbj4kpb9myj7y2mcc4k
title: Make the status SOURCE column conditional — drop it from the default listing, render only when requested
status: ideation
source: captain (2026-06-14) "make source field not default" — the `spacedock status` listing renders the SOURCE column UNCONDITIONALLY (internal/status/format.go:14 defaultStatusFields; :111-116 header, :136-137 rows). At boot the overview emitted ~7,200 tokens for ~30 entities, each carrying a multi-sentence SOURCE (boot-forensics, this session). A greet/overview rarely needs SOURCE. Separate from e6a/status-section-reader (which is file-BODY structural reads); this is the status TABLE column render. 0.20.4 read-cost theme.
started: 2026-06-15T05:35:19Z
completed:
verdict:
score: 0.35
worktree:
issue:
sprint: 0204-structured-reads
sprint-readiness: ready
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

### Parity-suite handling (mechanism note)

The default render is currently pinned native-vs-**oracle** (the Python reference) by `TestNativeReadMatchesOracle` case `"default"` and the `ind-*-default` / `seq-default` golden fixtures. Dropping SOURCE is a **deliberate divergence from the oracle**, mirroring the captain-approved `--fields` de-dupe divergence already in the tree (native_read_test.go:24-28: the divergent case was pulled OUT of oracle parity and locked by a native-only test, `TestFieldsDedupeNoDuplicateDefaultColumns`, instead). Implementation follows that precedent: the default-render cases move from oracle-parity into native-only golden fixtures regenerated with `-update`, and a native-only test asserts the new behavior. This is recorded so the implementer does not "fix" the now-failing oracle parity by reverting the column.

## Out of scope

- e6a / status-section-reader's file-BODY structural section reads (the separate structured-read helper) — a different surface.
- Any other column change, reorder, or width change to the status / next / boot tables.
- The `--boot` DISPATCHABLE render (already SOURCE-free via `printNextTable`).
- Re-deriving oracle parity for the default render: the oracle (Python reference) is frozen; native deliberately diverges here, as it already does for `--fields` de-dupe.

## Acceptance criteria

Each AC names a property of the finished behavior, not a stage action, and how it is verified.

**AC-1 — The default human listing omits SOURCE.** `spacedock status --workflow-dir <wf>` (no `--fields`/`--all-fields`) renders a table whose header is `ID SLUG STATUS TITLE SCORE` with no SOURCE column and no per-row provenance prose.
Verified by: a native-only golden fixture for the default render (regenerated `ind-seq-default.txt` / `native-read-default.txt` / `seq-default.txt`) whose header line contains no `SOURCE` token and whose rows carry no source value; plus an assertion-style native test (`countToken(header, "SOURCE") == 0`) so the property is checked, not only frozen.

**AC-2 — `--fields source` restores SOURCE in the human table.** `spacedock status --workflow-dir <wf> --fields source` renders a table that includes exactly one SOURCE column carrying each entity's source value.
Verified by: a native test asserting `countToken(header, "SOURCE") == 1` for `--fields source` (mirroring `fields_dedupe_test.go`'s header-token oracle), and a golden fixture freezing the restored column. A second assertion confirms `--all-fields` also surfaces a single SOURCE column (now as an extra, not the fixed slot).

**AC-3 — `--json --fields source` carries `source`; bare `--json` does not (under direction A).** A machine caller that names `source` receives `"source":…` in every entity object; a caller that does not name it receives an envelope with no `source` key.
Verified by: a native test over `--json --fields id,source` asserting each entity object contains `source`, and over bare `--json` asserting no entity object contains a `source` key (regenerated `seq-default.json`). (If the gate selects direction B, AC-3 is narrowed to only the `--fields source` half and bare `--json` retains `source`.)

## Test plan

- **Default-render goldens (AC-1):** regenerate the three default human goldens (`ind-seq-default.txt`, `native-read-default.txt`, `seq-default.txt`) with `go test ./internal/status -run '<default cases>' -update`, then add a native-only assertion test (no oracle) checking the header carries no `SOURCE` token. Pull the `"default"` case out of `TestNativeReadMatchesOracle`'s oracle-parity set into native-only coverage, exactly as the `--fields` de-dupe divergence was handled (native_read_test.go:24-28). Cost: low — fixture regen + one small test in the existing `internal/status` package. No live workflow run needed; the binary render is deterministic over a fixture workflow.
- **Opt-in render (AC-2):** extend `fields_dedupe_test.go`'s header-token style with `--fields source` and `--all-fields` cases asserting exactly one `SOURCE` column; freeze a golden for the restored render. Cost: low, fixture-level.
- **JSON envelope (AC-3):** regenerate `seq-default.json` (no `source` key under direction A) and add native assertions over bare `--json` (no `source` key) and `--json --fields id,source` (`source` present). Cost: low.
- **Regression sweep:** `go test ./internal/status/...` to confirm no other golden (`--all-fields`, `--where`, boot, resolve) silently shifted; any fixture that legitimately changes is regenerated and reviewed in the diff. The `--resolve` goldens MUST be unchanged (it never carried SOURCE) — a changed resolve golden is a red flag.

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
