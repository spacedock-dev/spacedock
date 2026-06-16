---
id: e6aaveste2tm0nsyqt407k55
title: spacedock status read helper — entity/markdown FM + section-heading offsets for surgical reads
status: done
source: "captain (2026-06-14) — 0.20.4 backbone. Reading whole entity bodies / stage reports / the README is a recurring FO + ensign token sink (this session: 280-line bodies, ~315-line README, 143KB CI logs, the Read-then-status--set staleness echo). A `spacedock status` helper that returns FM + a section-heading→offset map lets callers read the one section they need (Read offset/limit) instead of the whole file. Helps the README work (rzp) and other report-reading areas."
started: 2026-06-14T21:14:19Z
completed: 2026-06-16T14:55:29Z
verdict: PASSED
score: "0.40"
worktree:
issue:
sprint: 0204-structured-reads
sprint-readiness: in-progress
mod-block:
pr: "#386"
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
- **AC6 — Live adoption at a real injection point.** The adoption is proven by an FO (or ensign) ACTUALLY USING `status --read` at one of the wired contract sites — not by the contract text saying to. Concretely: at the Completion-and-Gates step (`first-officer-shared-core.md:105`), the FO reads an entity's latest `## Stage Report` by calling `status --read <ref> --json`, taking that heading's `offset`/`lines`, and issuing `Read(path, offset, limit)` — and the returned bytes are exactly that report section, with the rest of the (≥100-line) body never loaded. The proof is the behavioral trace (the two tool calls and the section-only result), never the presence of the instruction.
  *Verified by: live* FO section-read drive during validation, exercising the `:105` injection point. The instruction edit alone proves nothing (proof-policy: a contract saying "use --read" is not evidence; an agent using it is). The spike already proved the offset mechanism in-process; AC6 proves the wired contract actually drives the behavior.
- **AC7 — Contract edits are under-test scaffolding.** Each wired contract/skill file (the FO + ensign references and the dispatch ensign-prompt) carries the `--read` adoption wording, and the change does not break the scaffolding's own guards (the dispatch golden-prompt tests, the ban-readme-substring guard, any skill-text contract test).
  *Tested by:* the existing dispatch golden-prompt suite (`internal/dispatch/testdata/golden/*`, regenerated) passing, plus a re-run of the contract-lint / skill-text guards over the edited files. (Scaffolding guardrail: implementation makes these edits under the existing tests.)

## Test plan

- **Mechanism (already paid):** the spike above — throwaway parser + real `Read(offset,limit)` over a known-structure fixture. Seeds AC2's first test. Cost: done.
- **Unit (Go, `internal/status`):** AC1–AC3 + AC5 as table-driven tests over a committed fixture markdown (FM + the multi-level, fenced-code, stage-report sections the spike fixture exercised). Fixture lives under `internal/status/testdata/`. Cost: low — pure-function parser, no I/O beyond reading the fixture; minutes.
- **Behavior (Go, drives the binary):** AC4 — `status --read <ref>` against a fixture workflow dir (the `testdata/seq-workflow` pattern `golden_read_test.go` already uses), asserting the resolved-file map for a valid ref and the resolver error for an unknown one. Cost: low; reuses the existing native-runner test harness.
- **Live (validation stage):** AC6 — one FO/ensign section-read handoff against a real entity. Cost: one live exercise; runtime behavior is the only thing not provable in-process.
- **No golden-file capture needed** for the parser (the oracle is the fixture's structure, computed in-test), avoiding a brittle byte-golden that would re-break on every fixture edit.

## Contract adoption (where the FO/ensigns actually use `--read`)

The tool only captures the recurring token savings if the live read-instruction sites adopt it. Six verified injection points (file:line confirmed cycle 2). The concrete before/after wording for each lives in `## Documentation changes` below; this section is the adoption MAP and the per-site rationale.

| # | Site | Today | Adoption |
|---|------|-------|----------|
| 1 | FO `first-officer-shared-core.md:105` (Completion & Gates step 1) | "Read the entity file's last `## Stage Report` section." | `status --read <ref> --json` → that heading's offset/lines → `Read(offset, lines)`. **The prime FO use case** (the FO just read 127 lines for one report); this is the AC6 injection point. |
| 2 | FO `first-officer-shared-core.md:217` (Probe & Ideation Discipline, Grep-over-Read) | "Prefer Grep over Read… Anchor on heading… Read only when you need the full text." | Add `status --read` as the STRUCTURED upgrade to heading-anchoring: when you want a whole named section (not a grep hit), `--read` gives the exact offset/lines so the follow-up `Read` is section-scoped, not whole-file. |
| 3 | FO `claude-first-officer-runtime.md:39` (Entity-Body Inspection) | Defers to the shared-core Grep-over-Read rule; names the staleness echo. | Re-point the defer to also cover the new `--read` upgrade in the shared core — no independent wording, just keep the runtime note consistent (it already says "avoid a full-file Read for targeted section lookups"). |
| 4 | Ensign `ensign-shared-core.md:18` (Working step 1) | "Read the entity file before making changes." | Qualify: read the SECTIONS you need via `status --read` (FM + the stage def's relevant section) rather than the whole body, when you don't need full text. The ensign's dispatch already hands it the stage def; `--read` covers re-reads of the body. |
| 5 | Ensign `ensign-shared-core.md:92` (Stage Report Protocol) | "Append the report at the end of the entity file — do not read the entire entity body to find an insertion point." | `--read`'s `total_lines` IS the exact append offset — turn the "do not read to find an insertion point" rule into a positive instruction: get `total_lines` from `status --read … --json`, append after it. |
| 6 | Dispatch ensign-prompt `internal/dispatch/build.go:541` & `:545` ("Read the entity file at {path} …") | Whole-spec read instruction baked into every dispatch prompt. | **Decision: do NOT swap this to `--read`.** The dispatched ensign genuinely needs the full spec to do the stage work — a section-only read here would starve it. Instead add a one-line hint AFTER the entity-read line: re-reads of the body during work (e.g. to find the stage-report append point, site #5) should use `status --read`. Keeps the first read whole, makes subsequent reads cheap. |

**Why a hint, not a swap, at site 6:** the entity-read line is the ensign's single source of the spec; the savings is in REPEAT reads, not the first one. The golden-prompt tests (`internal/dispatch/testdata/golden/*`) pin this line, so the hint is an additive line the goldens regenerate to include (AC7).

## Documentation changes

Two classes: (a) the user-facing CLI doc for the new `--read` surface, and (b) the contract/skill-text adoption diffs (specific before/after wording per the ideation rule; implementation applies them under the scaffolding guards).

### (a) User-facing CLI docs

- **`docs/dev/README.md`** — under the existing `## Runtime Live CI` / helper-reference area (or a short `### Reading sections` note), add:
  > **Read one section, not the whole file.** `spacedock status --workflow-dir <wf> --read <entity-ref-or-path> --json` returns the file's frontmatter plus an ordered heading map (`text`, `level`, `offset`, `lines`). Pass a heading's `offset`/`lines` to `Read(path, offset, limit)` to load just that section — e.g. an entity's latest `## Stage Report`, or this README's `## Sprints`, without the rest.
- **`mkdocs.yml` / docs site** — if the CLI reference page enumerates `status` flags, add a `--read <ref-or-path>` row with the one-line description above. (Concrete wording deferred to implementation against whichever doc file enumerates the flags; the before/after is the addition of that single flag row — no existing wording changes.)

### (b) Contract / skill-text adoption diffs

Site 1 — `skills/first-officer/references/first-officer-shared-core.md:105`:
> **Before:** `1. Read the entity file's last `## Stage Report` section.`
> **After:** `1. Read the entity file's last `## Stage Report` section — `status --read <ref> --json`, take the last `## Stage Report` heading's `offset`/`lines`, then `Read(offset, limit)` that range, instead of loading the whole body.`

Site 2 — `first-officer-shared-core.md:217`:
> **Before:** `- Prefer Grep over Read for targeted entity-body inspection. Anchor on heading or field name (`## Stage Report`, `### Feedback Cycles`, a specific frontmatter field). Read only when you need the full text.`
> **After:** append one sentence: `… Read only when you need the full text. To pull a whole named section, `status --read <ref> --json` returns its `offset`/`lines` so the follow-up `Read(offset, limit)` is section-scoped, not whole-file.`

Site 3 — `claude-first-officer-runtime.md:39`:
> **Before:** `… avoid a full-file Read for targeted section lookups and trust `status --set` stdout (`field: old -> new`) for mutation narration.`
> **After:** `… avoid a full-file Read for targeted section lookups (use the shared core's `status --read` section-read upgrade) and trust `status --set` stdout (`field: old -> new`) for mutation narration.`

Site 4 — `skills/ensign/references/ensign-shared-core.md:18`:
> **Before:** `1. Read the entity file before making changes.`
> **After:** `1. Read the entity file before making changes — when you need only specific sections (the relevant stage-def section, your prior report), `status --read <entity-path> --json` returns each section's `offset`/`lines` for a scoped `Read`, rather than the whole body.`

Site 5 — `ensign-shared-core.md:92`:
> **Before:** `- Append the report at the end of the entity file — do not read the entire entity body to find an insertion point.`
> **After:** `- Append the report at the end of the entity file — get the file's `total_lines` from `status --read <entity-path> --json` and append after it; do not read the entire entity body to find an insertion point.`

Site 6 — `internal/dispatch/build.go` (the ensign-prompt assembly, after the `Read the entity file at {path}` line at `:541`/`:545`):
> **Add** (a new line in the assembled prompt, regenerating the `testdata/golden/*` prompts): `For later re-reads of the body during work — e.g. to find the append point for your stage report — use `${SPACEDOCK_BIN:-spacedock} status --read {path} --json` to fetch a section's offset/lines rather than re-reading the whole file.`
> (The first entity-read stays a whole-spec read; only repeat reads are redirected. The golden-prompt suite pins this output, so the added line is reflected in regenerated goldens — AC7.)

## Stage Report: ideation

- DONE: Spike the riskiest mechanism FIRST: prove `spacedock status` can emit a section-heading->offset map for a real fixture markdown such that feeding an offset/limit to an actual Read returns exactly that section's text; record the exercise in the entity body.
  Throwaway parser `/tmp/section-reader-spike/spike.go` over known-structure fixture; three real `Read(offset,limit)` calls returned exactly the predicted sections (offset 10/lines 7, 17/8 incl. nested L3, 30/4 to EOF). Recorded in body "Riskiest-mechanism spike".
- DONE: Pin the helper's concrete contract: invocation surface, output shape, caller usage pattern with Read(offset,limit). Acceptance criteria behavior-first with the fixture's known structure as the external oracle, plus a live FO/ensign section-read exercise.
  Body "Concrete contract": `--read <ref-or-path>` flag on `status` (not a new command), all-strings JSON envelope (frontmatter + ordered headings text/level/offset/lines, no byte offsets), two-call caller pattern. AC1–AC6 each cite a fixture-structure-oracle test; AC6 is the live exercise.
- DONE: Hold the scope boundary: file-BODY structural reads, NOT the status-TABLE SOURCE-column render (the separate `contract`/source-not-default concern). Keep them distinct in the design.
  Body "Out of scope" names the `defaultStatusFields`/`SOURCE` table-column concern in `internal/status/format.go:14` as explicitly out of scope and distinct.

### Summary

Designed the 0.20.4 read-helper as a `--read <ref-or-path>` flag on `spacedock status` (reusing the existing discovery/resolver/`--json` machinery) that emits parsed frontmatter plus a 1-based-line heading map (text, level, offset, lines) so a caller reads one section via `Read(offset, limit)`. Spiked the riskiest path first — the offset contract — by running a throwaway parser over a known-structure fixture and confirming three real `Read(offset,limit)` calls returned exactly the predicted sections, including a level-aware "section owns its nested headings" case. Scope held distinct from the status-table SOURCE-column (`contract`/source-not-default) concern; ACs are oracle-based (fixture structure, never prose-match) with a live FO/ensign section-read as AC6, plus a proposed README doc diff.

## Stage Report: ideation (cycle 2)

Cycle-2 expands scope to include the contract ADOPTION wiring (captain bounce-back: the tool without adoption ships no recurring savings). All cycle-1 design kept; this is additive.

- DONE: Design the contract ADOPTION at each verified live read-instruction site.
  New body section "Contract adoption" maps all six sites (file:line re-verified cycle 2: FO shared-core 105 & 217, FO runtime 39, ensign shared-core 18 & 92, dispatch build.go 541/545) to a specific instruction change, with per-site rationale. Site 6 (dispatch ensign-prompt) is a deliberate HINT-not-swap: the first entity read must stay whole-spec or it starves the ensign; only repeat reads redirect.
- DONE: Expand `## Documentation changes` with CONCRETE before/after instruction diffs for the contract files (not just the README CLI ref).
  Body "Documentation changes" now split (a) user-facing CLI docs and (b) six per-site before/after blocks with verbatim current wording and the proposed replacement, per the ideation rule (template/skill text → specific wording).
- DONE: Strengthen AC6 so the live drive proves a REAL injection point in use; proof is behavioral, never instruction prose.
  AC6 rewritten to require the FO actually calling `status --read` at the `:105` Completion-and-Gates site (the two-tool-call trace: `--read --json` → `Read(offset,limit)` → section-only bytes), explicitly stating the instruction edit alone proves nothing. Added AC7 covering the contract edits as under-test scaffolding (dispatch goldens + skill-text guards).
- DONE: Keep scope distinct from the status-table SOURCE render (source-not-default).
  Unchanged from cycle 1 — "Out of scope" still names `defaultStatusFields`/`SOURCE` in format.go as the separate concern; adoption work is all file-BODY reads.

### Summary

Folded in the adoption wiring: a six-site contract-adoption map (all file:line re-verified) plus concrete before/after skill-text diffs for the FO references, the ensign references, and the dispatch ensign-prompt — with the dispatch site deliberately a re-read HINT rather than a first-read swap (a section-only first read would starve the ensign of the spec). AC6 now demands a behavioral live drive at the `:105` FO completion-gate site (proof is the agent using `--read`, never the contract saying to), and AC7 covers the contract edits as under-test scaffolding (dispatch goldens regenerate to include the added prompt line). Scope still held distinct from the status-table SOURCE-column concern.

## Stage Report: implementation

- DONE: Ship the `--read <ref-or-path>` feature per the pinned contract (flag on `status`; all-strings JSON envelope of frontmatter + ordered headings text/level/offset/lines; 1-based offsets matching Read(offset,limit); level-aware section ownership; fenced-code headings skipped; no byte offsets), with AC1-AC5 tests green where the oracle is the fixture's known structure.
  `internal/status/section_read.go` (parser + runner) + `json_commands.go:readJSON` + `native_runner.go` wiring; fixture `internal/status/testdata/section-reader/fixture.md`; tests `internal/status/section_read_test.go` — AC1-AC5 PASS (commit ebbd9ba3). Mutation-checked: flipping `<=` to `<` in the ownership rule fails AC2; disabling the fence-skip fails AC3 — the oracle is real, not a tautology.
- DONE: Apply the six contract-adoption edits using the EXACT before/after wording from `## Documentation changes` (b), and regenerate the dispatch golden prompts so AC7's scaffolding guards pass.
  FO shared-core 105 (completion-gate read) & 217 (Grep-over-Read upgrade), FO claude-runtime 39 (defer), ensign shared-core 18 (entity-read) & 92 (stage-report append point), dispatch `build.go` site 6 as a re-read HINT after the whole-spec entity-read line (commit e75fede2). 13 dispatch goldens regenerated to include the additive hint; `go test ./internal/dispatch/` green (AC7). Hint uses the inline `${SPACEDOCK_BIN:-spacedock}` idiom (matches the skill-prose convention; the full launcher-function form is reserved for the executable `### Fetch commands` block).
- DONE: Apply the user-facing CLI doc diff (README `### Reading sections` note + the command-reference flag row), commit the deliverable.
  `docs/dev/README.md` `### Reading sections` note + `docs/site/reference/command-reference.md` status-flag enumeration (commit e75fede2). mkdocs had no per-flag table, so `--read` was added to the inline `spacedock status` flag enumeration (addition only, no rewording).

### Summary

Shipped `status --read <ref-or-path>` (commit ebbd9ba3): a flag on `status` returning a file's frontmatter plus an ordered ATX-heading map (text/level/1-based offset/section line count) so a caller reads one section via `Read(offset, limit)` instead of the whole file. Entity-ref form resolves like `--resolve`; a plain path reads any markdown directly. AC1-AC5 are oracle-based Go tests (fixture structure, mutation-verified non-tautological); the offset round-trip was also exercised live against the real 231-line dev README (`## Sprints` at offset 38/lines 8 → a real `Read(38, 8)` returned exactly that section). Wired the six adoption sites + regenerated the 13 dispatch goldens (AC7) and added the user-facing docs (commit e75fede2). Dogfooded `--read` to locate this report's append point (site #5's own use case).

**Flagged for validation:**
1. **AC6 owes the live FO `--read` drive at the `:105` completion-gate site** — the behavioral two-tool-call trace (`--read --json` → `Read(offset, limit)` → section-only bytes, rest of a ≥100-line body never loaded). The instruction edit alone proves nothing; the mechanism is proven in-process + live-against-README, but the WIRED contract driving an FO is validation's call.
2. **Detached adversarial audit owed** — this ships a contract change touching every dispatch prompt and the FO/ensign read discipline (high-stakes surface).
3. **PRE-EXISTING, OUT-OF-SCOPE test failures (not mine):** `TestSharedScenarioDocsContract` and `TestCodexForegroundWaitWatchdogDocsContract` in `internal/ensigncycle` fail because the README-slim commit `a9e669ae` relocated the Runtime-Live-CI / Codex-watchdog content out of `docs/dev/README.md` but left these two doc-contract guards pointing at it. Confirmed absent at my branch base (clause count 0 at `a9e669ae`); `docs/dev/README.md`'s relocated sections untouched by me. Fix is a judgment call (restore clauses vs. retarget guards to `docs/runtime-live-ci.md`) that intersects the FO's deliberate README-slimming — surfaced rather than guessed.

## Stage Report: implementation (cycle 2)

Cycle-2 is a rebase re-dispatch to de-risk a 25-commit merge gap before validation. No feature changes — the AC1-AC5 + adoption + docs work (commits ebbd9ba3 + e75fede2) stays. The branch was rebased onto current origin/main with the redundant README-slim commit dropped.

- DONE: Rebase the branch onto origin/main by DROPPING the redundant slim commit a9e669ae (`git rebase --onto origin/main a9e669ae`) so only the two feature commits replay onto current main.
  `git rebase --onto origin/main a9e669ae` succeeded clean (no conflicts); branch now sits directly on origin/main 49aff7ab with the two feature commits replayed as a5e6ff92 (status --read) + 63c3756a (adoption wiring). Verified main already had docs/runtime-live-ci.md + a 229-line slimmed README (vs a9e669ae's 227) — a9e669ae confirmed redundant before dropping.
- DONE: Resolve the replay conflicts: re-land the `### Reading sections` note into main's already-slimmed README, re-apply the six contract-adoption edits at their current sites, and regenerate the dispatch goldens.
  Git's 3-way merge landed all edits by content, not line number. Verified by diff vs origin/main: README `### Reading sections` note sits cleanly after the `--next` block / before `## Runtime Live CI`; all six adoption edits match the entity's `## Documentation changes (b)` verbatim wording; the 13 dispatch goldens carry the additive site-6 hint. `go test ./internal/dispatch -run TestBuild -update` produced ZERO drift — goldens already current against main.
- DONE: Full `go build ./... && go test ./...` green on the rebased branch — the two previously-flagged guard failures are GONE (not suppressed); AC1-AC5 + AC7 still pass.
  `go build ./...` exit 0; `go test ./...` ALL packages ok (exit 0), including internal/ensigncycle. The two guards are absent from the codebase (grep exit 1, not skipped) — main dropped both functions (#378) and dropping a9e669ae removed the slim that broke them. AC1-AC5 pass by name (TestReadFrontmatterParity / TestReadHeadingMapLocatesBytes / TestReadFencedHeadingsSkipped / TestReadEntityRefResolvesLikeResolve / TestReadFlagIncompatibility); AC7 via the green dispatch goldens.

### Summary

Rebased `spacedock-ensign/status-section-reader` onto current origin/main (49aff7ab) by dropping the redundant README-slim commit a9e669ae via `git rebase --onto origin/main a9e669ae`; only the two feature commits replayed (a5e6ff92 status --read, 63c3756a adoption wiring). The replay was conflict-free — git's 3-way merge re-landed the README note, the six contract-adoption edits, and the dispatch goldens by content; a `-update` regen confirmed the goldens carry no drift against main. Full `go build ./... && go test ./...` is green (exit 0 across all packages): the two previously-flagged guards (`TestSharedScenarioDocsContract`, `TestCodexForegroundWaitWatchdogDocsContract`) dissolved as predicted — they are absent from the codebase (main dropped them per #378, and dropping a9e669ae removed the breaking slim), not suppressed. AC1-AC5 + AC7 still pass. Per dispatch: not pushed, no PR — the Commander handles merge. AC6 (the live FO `--read` drive) remains validation's job, not attempted here.

## Stage Report: validation

- DONE: AC6 — live FO `--read` drive at the `first-officer-shared-core.md:105` completion-gate site, behavioral two-tool-call trace against a real ≥100-line entity.
  Built `/tmp/e6a-sd` from this branch (`go build -o /tmp/e6a-sd ./cmd/spacedock`, exit 0). Target: THIS entity (222 lines). **Call 1:** `e6a-sd status --workflow-dir docs/dev --read status-section-reader --json` → last `## Stage Report` heading (level 2) = `{"text":"Stage Report: implementation (cycle 2)","level":"2","offset":"209","lines":"14"}`. **Call 2:** `Read(status-section-reader.md, offset=209, limit=14)` → returned exactly lines 209-222 (the full `## Stage Report: implementation (cycle 2)` section, heading through its nested `### Summary` to EOF) — and ONLY those 14 lines. Body spans 18-222 (~205 lines); 14 loaded, the other ~191 (all four prior reports + contract + ACs) never read. Independent oracle: `awk END` = EOF line 222; offset+lines-1 = 209+14-1 = 222 = EOF. The level-2 section correctly owned its nested L3 `### Summary` (live, not just unit test). This is the agent USING `--read`, not the contract saying to.
- DONE: Detached adversarial audit of the high-stakes surface (the `--read` parser + the six contract-adoption edits), read-only on a throwaway checkout, never the implementation worktree.
  Clean `git clone --branch` to `/tmp/e6a-audit` + separate binary `/tmp/e6a-audit-sd`. Attacked the parser with adversarial fixtures, each checked against a `cat -n` / hand-computed oracle (never a prose match): CRLF (`\r\n` collapsed first → offsets match `wc -l`, round-tripped through a real `Read`), closed-ATX `## x ##`, trailing-whitespace heading, `~~~` fence, ```` ```go ```` info-string fence, empty section, deep nesting `## → #### → ### → ##` (the `<=` ownership boundary, both directions), final-section-to-EOF, unclosed-fence-swallows-to-EOF, 7-hash reject, no-space reject, empty file (`total_lines:0`, exit 0), no-frontmatter file. Resolver error PARITY confirmed: `--read <unknown>` and `--resolve <unknown>` emit byte-identical errors (`unknown reference: …`; `sd-b32 address prefix … too short`); entity-ref success returns the same `path` as `--resolve`. Flag-incompatibility (AC5) confirmed with well-formed args for `--set`/`--next`/`--validate`/`--boot`/`--resolve` (`--read cannot be combined with …`, exit 1).
- DONE: Adoption-edit wording-fidelity audit + scaffolding-guard re-run.
  All six edits match `## Documentation changes (b)` VERBATIM (sites 1-5 prose diffed line-for-line; site 6 build.go hint matches, and additionally points the re-read hint at the worktree-rewritten `readPath` so it isn't a stale main path). Guards green on the throwaway checkout: `internal/contractlint` (skill-text/contract guards), `internal/dispatch` (golden-prompt suite, AC7), `internal/ensigncycle` (doc-contract guards). `go test ./internal/dispatch -run TestBuild -update` produced ZERO golden drift (0 changed files) — the adoption edits are consistently under-test.
- DONE: AC cross-check AC1-AC7 each carry evidence after this stage.
  On the rebased branch (HEAD 63c3756a on 49aff7ab): `go build ./...` exit 0; `go test ./...` exit 0, 0 FAIL/panic across all 15 tested packages (the two prior-flagged guards confirmed gone, not suppressed). AC1-AC5 by name green (TestReadFrontmatterParity / TestReadHeadingMapLocatesBytes / TestReadFencedHeadingsSkipped / TestReadEntityRefResolvesLikeResolve / TestReadFlagIncompatibility, verbose run all PASS); AC7 via green `internal/dispatch` goldens + the contractlint/ensigncycle guards; AC6 by the live drive above.

### Summary

PASSED. The two proofs implementation left open are now satisfied with behavior, not assertion. AC6: a live two-tool-call trace (build the branch binary → `status --read … --json` → take the last `## Stage Report` offset 209/lines 14 → `Read(209, 14)`) returned exactly that 14-line section out of a ~205-line body, with the rest never loaded — the wired contract actually drove the section-scoped read. The detached adversarial audit (throwaway checkout, ~16 edge-case fixtures each checked against an external oracle) refuted nothing material: CRLF, the level-aware `<=` ownership boundary, fenced-code skipping (both ` ``` ` and `~~~`, info strings, unclosed-to-EOF), empty/no-FM files, resolver error parity, and flag-incompatibility all behave per contract. Two NON-material observations recorded for honesty, neither touching the load-bearing offset/lines contract and neither warranting route-back: (1) a heading whose text legitimately ends in a literal `#` (e.g. `## … hash \#`) has its trailing `#` stripped by the `TrimRight(body, "#")` rendered-text rule, leaving a dangling `\` in the `text` field — cosmetic, offset/lines unaffected, and trailing-`#`-in-heading-text is vanishingly rare in entity bodies; (2) a blockquoted `> ## x` is NOT treated as a section heading (diverges from strict CommonMark) — but this is the SAFER behavior for a section-locator (blockquoted headings are not document section boundaries) and is consistent. All six adoption edits match the entity's spec verbatim and break no scaffolding guard (contractlint + dispatch goldens + ensigncycle all green; `-update` zero drift). AC1-AC7 each carry reproduced evidence. Recommendation: PASSED — ready for the Commander's post-gate merge.
