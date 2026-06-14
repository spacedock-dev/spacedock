---
id: e6aaveste2tm0nsyqt407k55
title: spacedock status read helper — entity/markdown FM + section-heading offsets for surgical reads
status: ideation
source: "captain (2026-06-14) — 0.20.4 backbone. Reading whole entity bodies / stage reports / the README is a recurring FO + ensign token sink (this session: 280-line bodies, ~315-line README, 143KB CI logs, the Read-then-status--set staleness echo). A `spacedock status` helper that returns FM + a section-heading→offset map lets callers read the one section they need (Read offset/limit) instead of the whole file. Helps the README work (rzp) and other report-reading areas."
started: 2026-06-14T21:14:19Z
completed:
verdict:
score: "0.40"
worktree:
issue:
sprint: 0204-structured-reads
---

A `spacedock status` helper that, given an entity (or any markdown file), returns its parsed frontmatter PLUS a map of section headings with line offsets — so a caller reads the one section it needs (`Read(file, offset, limit)`) instead of loading the whole file. The general capability behind 0.20.4's reading-cost reduction.

## Problem

- The FO and ensigns read whole files to reach one part: a 280-line entity to get the latest `## Stage Report`, a ~315-line README to check the `## Sprints` section, a stage report buried mid-file. Tokens scale with the whole file, not the part needed.
- On Claude Code, a `Read` followed by a `status --set` mutation of the same file triggers the staleness echo (the whole file re-emitted as cache-write tokens) — another whole-file tax.
- There is no structured way to ask "what sections does this doc have, and where do they start" — callers read the whole file or guess offsets.

## Proposed approach (seed — ideation designs the interface + output shape)

A read helper (likely under `spacedock status`; the exact surface is ideation's call) that takes an entity ref or a markdown path and emits:
- the parsed frontmatter, and
- an ordered list of section headings, each with: heading text, level, start line-offset, and the section's line range (to the next heading) — so a caller can `Read(offset, limit)` exactly that section.

Use cases to keep in view: read an entity's latest stage report by heading; read the dev README's `## Sprints` section without the other ~300 lines (helps rzp); read any report-shaped markdown structurally.

## Concrete contract (ideation's call)

**Invocation surface — a new flag on `status`, not a new command.** `spacedock status --read <ref-or-path>` (sibling to `--resolve`/`--short-id`/`--validate`). Rationale: the discovery, `--workflow-dir`/`PIPELINE_DIR` resolution, entity-ref resolution (`resolveReferenceCandidates`), and `--json` envelope machinery all already live in the `status` dispatch (`internal/status/native_runner.go`); a new top-level `read` command would re-plumb all of it. `--read` accepts either an entity reference (resolved the same way `--resolve` resolves it, via `resolveReferenceOrExit`'s resolver) OR a plain filesystem path (so the README and any report-shaped markdown are addressable without being a tracked entity). It is mutually exclusive with the other action flags, the same incompatibility-list shape every existing action flag carries.

**Output shape — JSON envelope (the established all-strings ordered-object shape), text mirror for parity.** `--read` defaults to the table-style key=value text the other reads default to, and emits the structured envelope under `--json` (the agent-facing form callers will actually parse). The envelope:

```json
{"command":"read",
 "path":"<realpath of the file read>",
 "total_lines":"33",
 "frontmatter":{"id":"…","status":"…", …},
 "headings":[
   {"text":"Problem","level":"2","offset":"10","lines":"7"},
   {"text":"Stage Report: ideation","level":"2","offset":"30","lines":"4"}]}
```

- `frontmatter` reuses `ParseFrontmatter` verbatim (already the workflow's FM parser; top-level scalars only).
- `headings` is an ordered list of every ATX heading (`#`..`######` followed by a space) in the BODY (after the closing `---` fence), each with `text`, `level`, `offset`, `lines`. Headings inside fenced code blocks (` ``` ` / `~~~`) are skipped.
- `offset` is the **1-based line number** of the heading line — identical to Claude Code's `Read(offset, …)` semantics. `lines` is the section's line count, so the caller passes `Read(path, offset, lines)` to get exactly that section.
- **Section ownership:** a heading owns from its line to the line before the next heading of level ≤ its own; the final heading runs to EOF. So `## Proposed approach` includes a nested `### Sub-detail`, but stops at the next `##`. (Proven below.)
- Values are strings throughout (`level`/`offset`/`lines` are stringified ints) to match the existing all-strings `--json` contract (`internal/status/json.go`).
- **No byte offsets.** Claude Code's `Read` is line-addressed; byte offsets would be an unused second coordinate system. YAGNI — add only if a non-`Read` consumer appears.

**Caller usage pattern.** Two calls replace one whole-file read: (1) `spacedock status --workflow-dir <wf> --read <ref> --json` → parse `frontmatter` + find the target heading's `offset`/`lines`; (2) `Read(path, offset=<offset>, limit=<lines>)` → exactly that section. For the README's `## Sprints`: `--read docs/dev/README.md --json`, then `Read(offset, lines)` for that one heading instead of all 227 lines.

## Riskiest-mechanism spike (DONE — exercised, not asserted)

The riskiest path was the offset contract: does the helper's emitted `offset`/`lines` align with Claude Code's `Read(offset, limit)` 1-based-line semantics? Spiked with a throwaway Go parser (`/tmp/section-reader-spike/spike.go`) over a fixture with a KNOWN structure (`/tmp/section-reader-spike/fixture.md`; headings at lines 10, 17, 21, 25, 30 per `cat -n`):

- Parser emitted `## Problem` → offset 10/lines 7; `## Proposed approach` → 17/lines 8; `## Stage Report: ideation` → 30/lines 4.
- Fed each to a REAL `Read(offset, limit)`:
  - `Read(offset=10, limit=7)` returned lines 10–16 = exactly `## Problem` (heading through trailing blank, stopping before `## Proposed approach`).
  - `Read(offset=17, limit=8)` returned lines 17–24 = `## Proposed approach` INCLUDING its nested `### Sub-detail` (stopping before the next `##`) — proves the level-aware ownership rule.
  - `Read(offset=30, limit=4)` returned lines 30–33 = exactly `## Stage Report: ideation` to EOF.
- Independent oracle: the fixture's `cat -n` line numbers, not any substring match.

Mechanism confirmed: `Read.offset` == helper `offset` (1-based line), `Read.limit` == helper `lines`. The implementation's first test seeds directly from this exercise.

## Out of scope

- **The status-TABLE `source` column render** (`defaultStatusFields`/`SOURCE` in `internal/status/format.go`) — that is the separate `contract`/source-not-default concern (whether SOURCE renders by default in the listing table). THIS task is file-BODY structural reads (entity bodies, stage reports, README sections), not table-column projection. Kept distinct by design.
- **Non-markdown formats** (no JSON/YAML/code-structure reads) — markdown ATX headings only.
- **Byte offsets** (see contract above — line-addressed only).
- **Rewriting the README** — the FO slims `docs/dev/README.md` directly (FO realm; not a tracked task). This helper makes the slimmed README cheap to read structurally; the two compose.
- **A new top-level `read` command** — rejected in favor of the `--read` flag on `status` (rationale in the contract above).

## Relation

- **The dev README** — the FO slims it directly (FO realm over the workflow `README.md`; not a tracked task). This helper makes the slimmed README (and entity reports) cheap to read structurally; the two compose. The README's heading structure today (`## Sprints` at line 38, `## Stages` at 77, the `## Task Template` block at 164) is exactly the shape `--read` addresses.
- **`--resolve`** (`internal/status/resolve.go`, `resolveReferenceOrExit`) — `--read`'s entity-ref form reuses this resolver to turn a ref into a path.

## Acceptance criteria

All criteria are entity-level properties of the finished helper, each proven by EXERCISING real output against the fixture's known structure (the external oracle), never a prose/regex match.

- **AC1 — Frontmatter parity.** Given the fixture markdown (FM + several sections), `status --read <fixture> --json` emits a `frontmatter` object equal to `ParseFrontmatter` over the same file (every top-level scalar key/value present).
  *Tested by:* Go unit test asserting the emitted `frontmatter` map equals `ParseFrontmatter(fixture)`.
- **AC2 — Heading map locates the right bytes.** For every heading in the fixture, the emitted `offset`/`lines`, when used to slice the file's lines (`lines[offset-1 : offset-1+count]`), equal the section's known text — heading line through the line before the next heading of level ≤ its own (final heading to EOF).
  *Tested by:* Go unit test that, for each emitted heading, slices the fixture by `offset`/`lines` and asserts byte-equality against the fixture's known section text (table-driven over the fixture's headings). This is the in-process equivalent of the `Read(offset,limit)` exercise already done in the spike.
- **AC3 — Fenced-code headings are not mistaken for sections.** A `#`-prefixed line inside a ` ``` ` / `~~~` fence does NOT appear in `headings`.
  *Tested by:* Go unit test over a fixture whose code fence contains a `# not a heading` line; assert it is absent from the map.
- **AC4 — Entity-ref form resolves like `--resolve`.** `status --read <entity-ref>` (slug/id/prefix) reads the resolved entity's file; an unknown/ambiguous ref fails with the same resolver error shape `--resolve` produces.
  *Tested by:* Go behavior test driving the binary against a fixture workflow: a valid slug yields the map for that entity's file; an unknown ref exits 1 with the resolver's error.
- **AC5 — Flag incompatibility.** `--read` combined with `--set`/`--next`/`--boot`/`--validate`/`--resolve`/etc. exits 1 with the established "cannot be combined with …" message.
  *Tested by:* Go test asserting the incompatibility exit/message for each conflicting flag (mirrors the existing per-flag incompatibility tests).
- **AC6 — Live section read.** An FO or ensign, given a real entity (e.g. this one) and the helper, reads ONE target section (its latest `## Stage Report`) via `status --read … --json` then `Read(offset, lines)`, and the returned text is exactly that section — without loading the whole file.
  *Verified by: live* FO/ensign exercise during validation (a runtime handoff is the claim; the spike already proved the offset mechanism in-process).

## Test plan

- **Mechanism (already paid):** the spike above — throwaway parser + real `Read(offset,limit)` over a known-structure fixture. Seeds AC2's first test. Cost: done.
- **Unit (Go, `internal/status`):** AC1–AC3 + AC5 as table-driven tests over a committed fixture markdown (FM + the multi-level, fenced-code, stage-report sections the spike fixture exercised). Fixture lives under `internal/status/testdata/`. Cost: low — pure-function parser, no I/O beyond reading the fixture; minutes.
- **Behavior (Go, drives the binary):** AC4 — `status --read <ref>` against a fixture workflow dir (the `testdata/seq-workflow` pattern `golden_read_test.go` already uses), asserting the resolved-file map for a valid ref and the resolver error for an unknown one. Cost: low; reuses the existing native-runner test harness.
- **Live (validation stage):** AC6 — one FO/ensign section-read handoff against a real entity. Cost: one live exercise; runtime behavior is the only thing not provable in-process.
- **No golden-file capture needed** for the parser (the oracle is the fixture's structure, computed in-test), avoiding a brittle byte-golden that would re-break on every fixture edit.

## Documentation changes

`--read` is a new user-visible CLI surface, so ideation proposes the doc diff (implementation applies it):

- **`docs/dev/README.md`** — under the existing `## Runtime Live CI` / helper-reference area (or a short `### Reading sections` note), add:
  > **Read one section, not the whole file.** `spacedock status --workflow-dir <wf> --read <entity-ref-or-path> --json` returns the file's frontmatter plus an ordered heading map (`text`, `level`, `offset`, `lines`). Pass a heading's `offset`/`lines` to `Read(path, offset, limit)` to load just that section — e.g. an entity's latest `## Stage Report`, or this README's `## Sprints`, without the rest.
- **`mkdocs.yml` / docs site** — if the CLI reference page enumerates `status` flags, add a `--read <ref-or-path>` row with the one-line description above. (Concrete wording deferred to implementation against whichever doc file enumerates the flags; the before/after is the addition of that single flag row — no existing wording changes.)

## Stage Report: ideation

- DONE: Spike the riskiest mechanism FIRST: prove `spacedock status` can emit a section-heading->offset map for a real fixture markdown such that feeding an offset/limit to an actual Read returns exactly that section's text; record the exercise in the entity body.
  Throwaway parser `/tmp/section-reader-spike/spike.go` over known-structure fixture; three real `Read(offset,limit)` calls returned exactly the predicted sections (offset 10/lines 7, 17/8 incl. nested L3, 30/4 to EOF). Recorded in body "Riskiest-mechanism spike".
- DONE: Pin the helper's concrete contract: invocation surface, output shape, caller usage pattern with Read(offset,limit). Acceptance criteria behavior-first with the fixture's known structure as the external oracle, plus a live FO/ensign section-read exercise.
  Body "Concrete contract": `--read <ref-or-path>` flag on `status` (not a new command), all-strings JSON envelope (frontmatter + ordered headings text/level/offset/lines, no byte offsets), two-call caller pattern. AC1–AC6 each cite a fixture-structure-oracle test; AC6 is the live exercise.
- DONE: Hold the scope boundary: file-BODY structural reads, NOT the status-TABLE SOURCE-column render (the separate `contract`/source-not-default concern). Keep them distinct in the design.
  Body "Out of scope" names the `defaultStatusFields`/`SOURCE` table-column concern in `internal/status/format.go:14` as explicitly out of scope and distinct.

### Summary

Designed the 0.20.4 read-helper as a `--read <ref-or-path>` flag on `spacedock status` (reusing the existing discovery/resolver/`--json` machinery) that emits parsed frontmatter plus a 1-based-line heading map (text, level, offset, lines) so a caller reads one section via `Read(offset, limit)`. Spiked the riskiest path first — the offset contract — by running a throwaway parser over a known-structure fixture and confirming three real `Read(offset,limit)` calls returned exactly the predicted sections, including a level-aware "section owns its nested headings" case. Scope held distinct from the status-table SOURCE-column (`contract`/source-not-default) concern; ACs are oracle-based (fixture structure, never prose-match) with a live FO/ensign section-read as AC6, plus a proposed README doc diff.
