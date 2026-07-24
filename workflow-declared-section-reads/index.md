---
title: Workflow-declared section reads assemble stage context by pointer
status: ideation
source: "Durable-decisions sprint audit of Subspace 198f762, Spacedock 6y, and held js6 stakes-read-through, 2026-07-24."
score: "0.8"
id: 0cj3qf6fefedfj7j9exq62jb
gates:
    version: 1
    current:
        gate: gate:docs-dev:0c:backlog
    records:
        - id: gate:docs-dev:0c:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0c-backlog-1
              briefing:
                id: briefing:docs-dev:0c:backlog:attempt-1:revision-1
                digest: sha256:116794807857c45a046d56cf9217b822e6c4fbe21077edb36f0951a02995f364
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:0c:backlog:1
                briefing: briefing:docs-dev:0c:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-24T03:59:33.82556Z"
                decision: approve
                reason: approve, no dispatch yet
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
started: 2026-07-24T08:41:39Z
---

A workflow can put a cross-stage policy in its README, but `dispatch show-stage-def` returns only the selected `### stage` subsection. Let workflow metadata point at additional relevant sections and let the dispatched ensign load those exact sections through the existing structured read path.

## Problem

Subspace commit `198f762` added a top-level `## Constraint Authority`, while its stage definitions mention the policy without carrying its three authority rules. The live dispatch packet fetches only the stage subsection, so the worker cannot reliably apply the top-level policy.

Copying the policy into every stage is ad hoc and creates divergent authorities. Injecting the whole README is noisy. Hardcoding `Constraint Authority`, `Stakes`, or another heading into the launcher would create one product mechanism per policy.

The existing substrate is close. Workflow stages already live in README frontmatter. `internal/status/section_read.go` already produces a fence-safe, level-aware heading map with exact offsets and section lengths, while every worker assignment already carries one `dispatch show-stage-def --workflow-dir … --stage …` fetch. What is missing is an ordered metadata pointer, builder-time validation, and assembly by that existing read.

## Proposed approach

### Declared contract

Add an optional ordered YAML sequence named `context-sections` to `stages.defaults` and to an individual `stages.states[]` item:

```yaml
stages:
  defaults:
    context-sections:
      - Constraint Authority
  states:
    - name: ideation
    - name: implementation
      context-sections:
        - Constraint Authority
        - Repository Safety
    - name: done
      context-sections: []
```

Each value is the exact parsed text of one ATX heading in the same workflow README. It has no product-specific meaning and may name a heading at any level. A stage-local sequence replaces the default sequence; an absent stage-local key inherits the default; an explicit empty sequence clears it. Replacement rather than merge makes the effective ordered list visible in one place and avoids implicit de-duplication.

The current `### {stage}` subsection remains mandatory and is always first. `dispatch show-stage-def` then appends each resolved section in declaration order, with one blank line between sections and one final newline. Each section spans its heading through the line before the next heading of equal or higher rank, exactly as the existing fence-safe reader defines ownership; trailing blank lines receive the same trim the current stage extractor applies. With no effective `context-sections`, stdout is byte-identical to today's command.

Repeated selector values are invalid rather than silently de-duplicated. A selector resolving to zero headings is missing; one resolving to more than one heading is ambiguous; selecting the already-selected stage heading is duplicate context. Selector comparison uses the existing parser's normalized heading text (whitespace and optional closing `#` run already stripped), not a new Markdown interpretation.

### Builder and worker ownership

Expose a narrow in-process selection operation from the existing `internal/status/section_read.go` parser. Both dispatch paths use it:

1. `dispatch build` resolves the stage's effective `context-sections` and validates the complete selection before writing a dispatch artifact. A failure exits nonzero and names the README path, stage, selector, and failure kind.
2. The emitted assignment keeps exactly one existing fetch command:
   `dispatch show-stage-def --workflow-dir <dir> --stage <stage>`.
3. When the ensign runs that command during its one-file bootstrap, `show-stage-def` reads the same metadata and returns the stage definition plus the selected context. It revalidates so a README edit between build and bootstrap cannot silently change the selection.

The First Officer therefore neither reads nor assembles policy sections. No `--include-section` flag, `status --read` projection, README snapshot, extra bootstrap read, or second heading parser is introduced.

### Mechanism-to-value trace

- `context-sections` serves AC-1 by giving an authored policy a stage-scoped pointer. The simplest alternative, copying policy prose into stage definitions, is insufficient because it creates multiple authorities that can diverge.
- The shared fence-safe selector serves AC-1 and AC-2 by preserving authored section spans and generic order. A new dispatch-only parser is insufficient because fenced headings and ownership rules could diverge from the established reader.
- Builder validation plus the existing `show-stage-def` fetch serves AC-1 and AC-3 by failing before spawn while preserving the worker's single fetch. FO assembly or one fetch per section is insufficient because it reintroduces hand assembly or a bootstrap multi-read protocol.

### Spike basis

No new spike is needed. The held `js6` spike exercised the risky byte path: the shipped fence-safe reader located an arbitrary README section, its offset/length slice matched the authored section byte-for-byte after the current trailing-blank trim, and a fenced heading did not enter the map. `TestBuildStageDisciplineRidesFetchNotInlineAssignment` already proves that the assignment's `show-stage-def` fetch—not inline FO prose—delivers workflow stage context. YAML sequence order is preserved by the `yaml.v3` node representation already used by `ParseStagesWithDefaults`; implementation begins with a red metadata/assembly behavior fixture rather than another throwaway parser.

### Expected surface

Expected implementation surface is 8-10 files and about 380-560 changed lines: 120-180 production lines in `internal/status/stages.go`, `internal/status/section_read.go`, `internal/dispatch/build.go`, `internal/dispatch/showstagedef.go`, and dispatch help; 220-330 fixture/test lines under `internal/status` and `internal/dispatch`; and 20-35 documentation lines. The approved tolerance is at most 12 files or 700 changed lines. Crossing either bound, adding a package, or touching a skill/runtime adapter requires returning to the captain before expanding scope.

## Obligation delta

- **Bearer:** a dispatched ensign receives only the stage definition and sections explicitly selected for its current stage; the dispatch builder validates the same resolved selection before spawn.
- **Burden:** one existing fetch and a list bounded by the effective ordered declaration; no implicit whole-README read and no per-heading special case.
- **Authority:** the workflow's committed frontmatter plus captain direction from this audit.
- **No inferred obligation:** a referenced section supplies workflow context; it does not independently widen an entity's accepted ACs.

## Out of scope

- Hardcoding `Constraint Authority`, `Stakes`, development stage names, or any section's meaning.
- Injecting the entire workflow README or copying shared policies into every stage definition.
- Exposing selected bodies or selectors in `status --read`, adding a second section parser, or adding `--include-section`/similar CLI flags.
- Asking the First Officer to assemble sections, writing a README snapshot into the dispatch artifact, or making the ensign perform multiple bootstrap reads.
- A general artifact registry, remote URI fetcher, or dynamic include language.
- Changing gate/recorder semantics; the sibling FO follow-up owns gate-time provenance.

## Acceptance criteria

**AC-1 (VALUE) - A dispatched worker receives every byte of the workflow policy selected for its stage, without copying that policy into the stage subsection.**
Verified by: a fixture with a unique policy section measures the current undeclared baseline as 0 policy bytes, then drives `dispatch build`, executes its emitted one-file bootstrap fetch, and measures authored-policy-length bytes in the result with exact equality to `stage + separator + authored policy + final LF`. Removing the selector, changing one authored byte, or making the fetch stage-only makes the assertion fail.

**AC-2 - Workflow-owned section selection is generic, ordered, and bounded.**
Verified by: a custom-stage fixture selects two arbitrary heading names in reverse README order and asserts output follows declaration order, contains no unrelated section bytes, inherits a default, replaces it at stage scope, and clears it with `[]`. A fenced lookalike heading is not selectable. Swapping order, merging defaults, leaking an unrelated section, or recognizing the fenced lookalike fails the fixture.

**AC-3 - Invalid section declarations stop dispatch before a worker artifact can be consumed.**
Verified by: table-driven `dispatch build` cases for a missing selector, repeated selector, ambiguous matching headings, selection of the current stage heading, and non-sequence YAML all exit nonzero, write no usable dispatch body, and identify the README path, stage, selector/key, and failure kind. Correcting each independent fixture makes only its corresponding case pass.

**AC-4 - Undeclared workflows preserve the stable dispatch and read contract byte-for-byte.**
Verified by: every existing dispatch golden remains unchanged; existing `show-stage-def` parity goldens remain unchanged; an explicit regression fixture compares the no-declaration build body and fetch stdout to their checked-in bytes and asserts exactly one existing `show-stage-def` fetch command. Any extra flag, fetch, packet block, separator, or status projection fails the fixed baseline.

**AC-5 - Declared context is delivered through the same host-neutral bootstrap on every runtime.**
Verified by: the shared dispatch/integration fixture executes the emitted fetch command for Claude, Codex, and Pi envelopes and observes the same exact assembled bytes, while the existing approval-gated `claude-live`, `codex-live`, and `pi-live` lanes remain green because the host-neutral dispatch path changed. Removing the fetch from any host envelope or changing one host's quoting fails the fixture; an unapproved live lane is not counted as a pass.

## Test plan

Implementation starts with the AC-1 red behavior fixture. It authors a README with a custom stage, one unique top-level policy, one unrelated section, and no copied policy in the stage subsection; the test builds the assignment, parses and executes its emitted fetch, and compares observed bytes to independently-authored fixture strings. Estimated cost: medium, about 80-120 test lines.

Add focused tests around the existing parser and stage metadata:

- `internal/status`: ordered YAML-sequence decoding, default inheritance, stage replacement, explicit empty clear, exact section spans, fenced-heading exclusion, and ambiguous-heading reporting. These exercise the existing parser; they do not add JSON fields or a second parser. Estimated cost: 90-130 test lines.
- `internal/dispatch`: build-time validation cases and exact `show-stage-def` assembly for two arbitrary headings in declaration order. Extend `TestBuildStageDisciplineRidesFetchNotInlineAssignment` or its fixture rather than creating a parallel bootstrap harness. Estimated cost: 120-170 test lines.
- Compatibility: leave every existing dispatch and `show-stage-def` golden untouched, then run the golden suites as the no-declaration byte baseline. Add only declared-context fixtures; do not regenerate unrelated goldens.
- Runtime coverage: use the existing shared host-envelope fixture to run the same fetch for Claude, Codex, and Pi. Because this changes the host-neutral dispatch path, the existing approval-gated `claude-live`, `codex-live`, and `pi-live` lanes are required green; no new LLM scenario is justified unless the shared fixture proves unable to execute the fetch.

Focused verification is `go test ./internal/status ./internal/dispatch ./skills/integration`, followed by `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

### Documentation diff for implementation

Update `docs/site/concepts/workflows-and-entities.md` after the frontmatter example:

```diff
+A stage can also select ordered context from other README sections:
+
+    stages:
+      defaults:
+        context-sections:
+          - Constraint Authority
+
+`context-sections` contains exact README heading text. A stage-specific list
+replaces the default, and `[]` clears it. `dispatch build` validates every
+selection before spawn; the worker's existing `show-stage-def` read returns
+the stage definition followed by those sections in declaration order.
+Workflows without the field keep the existing output unchanged.
```

Update `dispatch show-stage-def --help` from “Print the workflow README section for a stage.” to “Print a stage’s workflow README section followed by its declared context sections.”

## Stage Report: ideation

- DONE: Define the builder-owned section-selection contract: workflow metadata declares ordered sections, dispatch build validates them, and the existing stage-definition read returns the selected context for the ensign's one-file bootstrap.
  Defined `context-sections`, inheritance/replacement/clear semantics, exact assembly order, shared build/read validation, diagnostics, and the unchanged single fetch command.
- DONE: Reuse the existing fence-safe section parser and existing dispatch/live fixtures; add no status projection, second parser, bootstrap multi-read protocol, or product-specific heading semantics.
  The design exposes only a narrow in-process selector from `internal/status/section_read.go`, extends the existing fetch fixture, and explicitly excludes every prohibited surface.
- DONE: Measure that declared policy reaches the worker exactly while undeclared workflows remain byte-identical, with a bounded expected surface and falsifiable test plan.
  AC-1 measures 0-to-authored-length exact policy bytes; AC-4 pins existing goldens; the plan bounds work to 8-10 files/380-560 changed lines with a 12-file/700-line tolerance.

### Summary

The ideation now specifies one workflow-owned ordered selector contract whose validation belongs to `dispatch build` and whose content reaches the ensign through the existing `show-stage-def` fetch. It preserves no-declaration bytes, reuses the fence-safe parser and current fixture/live lanes, and records exact failure, proof, documentation, and surface bounds.
