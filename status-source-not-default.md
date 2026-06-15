---
id: 0qdayfbj4kpb9myj7y2mcc4k
title: Make the status SOURCE column conditional — drop it from the default listing, render only when requested
status: backlog
source: captain (2026-06-14) "make source field not default" — the `spacedock status` listing renders the SOURCE column UNCONDITIONALLY (internal/status/format.go:14 defaultStatusFields; :111-116 header, :136-137 rows). At boot the overview emitted ~7,200 tokens for ~30 entities, each carrying a multi-sentence SOURCE (boot-forensics, this session). A greet/overview rarely needs SOURCE. Separate from e6a/status-section-reader (which is file-BODY structural reads); this is the status TABLE column render. 0.20.4 read-cost theme.
started:
completed:
verdict:
score: 0.35
worktree:
issue:
sprint: 0204-structured-reads
---

The `spacedock status` listing table always renders the SOURCE column, the heaviest per-row field. Drop it from the default render so the common overview is cheap; surface it only when a caller explicitly asks. Reduces the recurring boot/overview read-cost the FO pays every session.

## Problem

`internal/status/format.go` renders SOURCE unconditionally: it is in `defaultStatusFields` (line 14) and emitted in the table header (lines 111-116) and every row (lines 136-137). SOURCE is multi-sentence provenance prose; across ~30 entities it dominated the boot overview at ~7,200 tokens (boot-forensics this session). The FO running a greet, an overview, or a "what's dispatchable" check almost never needs SOURCE — but pays for it every time. There is no current way to render the human table without it (`--fields` narrows only the `--json` read, not the human table).

## Proposed approach

{Ideation fills. Candidate direction to evaluate, not prescribe: remove SOURCE from `defaultStatusFields` so the default human table omits it; add an explicit opt-in (e.g. a `--with-source` / `--all-fields` flag, or `--fields` honored for the human table too) that materializes it when named; keep the `--json` envelope carrying `source` when requested. Decide whether `--resolve`/single-entity views still show SOURCE (where it is cheap and useful).}

## Out of scope

{Ideation fills. Likely: e6a's file-BODY section reads (the separate structured-read helper); other column changes; reordering the table.}

## Acceptance criteria

Each AC names a property of the finished behavior, not a stage action, and how it is verified.

**AC-1 — The default listing omits SOURCE.** `spacedock status --workflow-dir <wf>` renders the table without the SOURCE column.
Verified by: {a golden-fixture test over the default human render asserting no SOURCE column; ideation pins it.}

**AC-2 — SOURCE renders when explicitly requested.** The opt-in surface (flag chosen in ideation) brings SOURCE back in the human table.
Verified by: {a golden/behavior test exercising the opt-in render; ideation pins it.}

**AC-3 — The `--json` envelope still carries `source` when requested.** Machine callers that name `source` still receive it.
Verified by: {a test over the `--json --fields …` path asserting `source` is present when named; ideation pins it.}

## Test plan

{Ideation fills. Likely golden-fixture tests over `internal/status/format.go` for the default vs opt-in human renders, plus the existing `--json --fields` coverage.}
