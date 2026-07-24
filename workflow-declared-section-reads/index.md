---
title: Workflow-declared section reads assemble stage context by pointer
status: validation
source: "Durable-decisions sprint audit of Subspace 198f762, Spacedock 6y, and held js6 stakes-read-through, 2026-07-24."
score: "0.8"
id: 0cj3qf6fefedfj7j9exq62jb
gates:
    version: 1
    current:
        gate: gate:docs-dev:0c:validation
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
        - id: gate:docs-dev:0c:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:0c-ideation-1
              briefing:
                id: briefing:docs-dev:0c:ideation:attempt-1:revision-1
                digest: sha256:b1e0cf2b6b7300464550d7af90d873802f7f51550a908d9eeb8cc62e0ea34400
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:0c:ideation:1
                briefing: briefing:docs-dev:0c:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-24T09:11:43.521688Z"
                decision: approve
                reason: Cycle 3 preserves the single-fetch boundary, proved the real-parser raw-span and mixed-newline mechanisms, stays within the declared surface, and independent staff re-review reports APPROVE with no material findings.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:docs-dev:0c:validation
          stage: validation
          attempts:
            - id: gate-attempt:0c-validation-1
              briefing:
                id: briefing:docs-dev:0c:validation:attempt-1:revision-1
                digest: sha256:fc9e60c5c84fb15d4aa4b5aaecb8f1001419c6c90cf36ffae3386403449ab14c
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:0c:validation:1
                briefing: briefing:docs-dev:0c:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-24T11:41:44.818812Z"
                decision: approve
                reason: 'Delegated sprint approval: exact candidate 396c92a3 satisfies all six ACs with 15/0/0 validation, Roborev 1990 clean, and all five live CI jobs green.'
                adoption-note: Captain granted the First Officer the conn to approve sprint gates, relevant CI lanes, PRs, and merge; land only exact candidate 396c92a302e1ac220d552d6f16600e9f08ccb622.
              application:
                action: advance
                target-stage: done
                state: pending
                blockers: []
started: 2026-07-24T08:41:39Z
worktree: .worktrees/spacedock-ensign-workflow-declared-section-reads
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

Each value is the exact parsed text of one ATX heading in the same workflow README. It has no product-specific meaning and may name a heading at any level. A stage-local sequence replaces the default sequence; an absent stage-local key inherits the default; an explicit empty sequence clears it. The typed metadata representation therefore retains both a presence bit and an ordered slice: absent is distinct from present-and-empty. Replacement rather than merge makes the effective ordered list visible in one place and avoids implicit de-duplication.

### Precise parser reconciliation

Do not reinterpret the current stage through the fence-safe selector scanner: that would change undeclared output for inputs where the legacy stage parser treats a heading inside a fence—or after a non-LF `splitlines` separator—as eligible. Instead, preserve both existing parsers and give their results one comparison coordinate:

- Extend the existing `extractStageSubsection` implementation to accept an in-memory README byte slice and return its existing rendered stage bytes plus the selected half-open **raw source byte span** `[startByte,endByte)`. Its decorated-token matching, fence treatment, `##`/`###` termination, malformed diagnostic, universal `splitTextLines` separators, trailing-blank trim, and output bytes remain unchanged.
- Extend the existing fence-safe `internal/status/section_read.go` scan to expose each selected heading's existing rank/ownership result plus its half-open raw source byte span. Its exact heading text, fence skipping, level-aware ownership, and `status --read` projection remain unchanged.
- Compare only raw byte intervals, never line numbers from the two incompatible line models. Both parsers receive the same immutable byte buffer, so parent/child/stage intersection is unambiguous even across CRLF or exotic logical separators.

This is a reconciliation of the two shipped parsers, not a third parser and not a claim that their heading semantics are identical. The legacy extractor remains the authority for the mandatory current-stage subsection and therefore preserves every no-declaration byte. The fence-safe scanner remains the authority for declared context selectors, so fenced lookalikes never satisfy a selector. Existing decorated, malformed, separator, CRLF, nested-heading, and no-declaration fixtures plus the cycle-2 pre-refactor fenced-stage characterization pin both sides without golden regeneration.

### Declared-context rendering

The current `### {stage}` subsection remains mandatory and is always first. Its rendering is byte-for-byte the existing legacy result. Each selected raw source span then receives one explicit canonical rendering:

1. Split with the existing `splitTextLines` semantics, so LF, CRLF, lone CR, VT, FF, FS, GS, RS, NEL, LS, and PS are logical line boundaries and every boundary is represented as LF in output.
2. Remove only trailing logical lines whose `strings.TrimSpace` value is empty. Preserve the heading, every nonblank byte (including UTF-8 and trailing spaces on nonblank lines), and all leading or internal blank lines.
3. Join the remaining logical lines with `\n`; the rendered component itself has no trailing terminator.
4. Assemble the rendered stage followed by rendered context components in declaration order with exactly `\n\n` between components, then append exactly one final `\n` to command stdout.

Thus a raw declared section `## Pølicy\r\nα\rβ\u0085γ\r\n \t\r\n` canonically renders as `## Pølicy\nα\nβ\nγ`; when appended after a stage it is preceded by two LF bytes and the entire command ends in one LF byte. This is the exact-byte contract: equality is against the canonical rendering, not the source's newline encoding. With no effective `context-sections`, assembly is the existing rendered stage plus its existing final LF and remains byte-identical to today's command.

Repeated selector values are invalid rather than silently de-duplicated. A selector resolving to zero headings is missing; one resolving to more than one heading is ambiguous. After resolving unique headings, validate all half-open spans pairwise: every selected span must be disjoint from the current-stage span and from every other selected span. Intersection rejects all three material overlap shapes—selected parent containing the stage, selected child inside the stage, and selected parent containing another selected child—as well as exact duplicate spans. Sibling spans that only meet at an endpoint remain valid. Diagnostics name both headings and both `[start,end)` spans.

Selector comparison uses the existing parser's normalized heading text (whitespace and optional closing `#` run already stripped), not a new Markdown interpretation. Span intersection is structural only and assigns no meaning to a heading name or rank.

### Builder and worker ownership

Expose both existing extractors over an in-memory README byte slice, and add a strict in-memory stage-metadata parse that returns an error for malformed YAML, a wrong-kind `context-sections`, or a non-scalar sequence item. The existing path wrappers may retain their compatibility behavior, but build/read resolution uses one shared `resolveStageContext(readmeBytes, stage)` operation and one immutable byte buffer per invocation. No metadata/body reread can mix two README versions inside one command.

Both dispatch paths use that resolver:

1. `dispatch build` reads the README once, resolves the stage's effective `context-sections`, validates metadata, headings, and non-intersecting spans, and only then writes a dispatch artifact. A failure exits nonzero and names the README path, stage, selector, and failure kind.
2. The emitted assignment keeps exactly one existing fetch command:
   `dispatch show-stage-def --workflow-dir <dir> --stage <stage>`.
3. When the ensign runs that command during its one-file bootstrap, `show-stage-def` reads the then-current README once and runs the same resolver. Valid current metadata/content is intentionally adopted; invalid current metadata, missing/ambiguous selectors, or newly intersecting spans make the fetch exit nonzero before stage work.

This is a live-read contract, not a version pin. Build preflights the version it sees before spawn; the pointer intentionally does not carry a digest or snapshot, so build cannot reject a later valid edit. Bootstrap is authoritative for the bytes the worker receives and fails only when the current version violates the selection contract.

The First Officer therefore neither reads nor assembles policy sections. No `--include-section` flag, `status --read` projection, README snapshot, extra bootstrap read, or second heading parser is introduced.

### Mechanism-to-value trace

- `context-sections` serves AC-1 by giving an authored policy a stage-scoped pointer. The simplest alternative, copying policy prose into stage definitions, is insufficient because it creates multiple authorities that can diverge.
- Raw-byte span reconciliation between the unchanged legacy stage extractor and unchanged fence-safe selector scanner serves AC-1, AC-2, AC-3, and AC-4. Forcing either parser's heading semantics onto the other is insufficient because it changes an existing stage byte contract or weakens fence-safe context selection; a third parser would add another divergence.
- The canonical selected-span renderer serves AC-1 and AC-4 by making mixed newline input deterministic while leaving undeclared stage output untouched. Raw emission is insufficient because it combines a newline-normalized legacy stage with editor-dependent context terminators and leaves inter-section/final termination ambiguous.
- Strict, presence-aware YAML decoding serves AC-2 and AC-3 by preserving absent/inherited, explicit-empty, replacement, and wrong-kind states. Reusing the scalar map is insufficient because it collapses sequences to `""` and has no parse-error channel.
- Builder preflight plus live `show-stage-def` resolution serves AC-1, AC-3, and AC-5 by failing before spawn and again before work while preserving the worker's single fetch. FO assembly, snapshots, or one fetch per section are insufficient because they reintroduce hand assembly, stale copies, or a bootstrap multi-read protocol.

### Spike basis

The held `js6` spike already exercised the byte path: the shipped fence-safe reader located an arbitrary README section, its offset/length slice matched the authored section byte-for-byte after the current trailing-blank trim, and a fenced heading did not enter the map. `TestBuildStageDisciplineRidesFetchNotInlineAssignment` proves that the assignment's `show-stage-def` fetch—not inline FO prose—delivers workflow stage context.

Cycle 2 exercised the previously unverified metadata and temporal path with a throwaway Go test in `internal/status`, using the real `stagesNodeFromFrontmatter`, `mappingValue`, `scanHeadings`, and `splitLines` primitives. The first run failed on explicit `[]`: the prototype preserved the key's presence but returned a nil slice, while the oracle expected present-and-empty. Initializing an empty slice only when the YAML sequence node is present fixed that distinction. The rerun command `go test ./internal/status -run TestIdeationSpikeContextSectionsSingleBufferTriStateAndLiveRead -v` passed and observed:

- absent stage metadata inherited `[Authority Safety]`;
- explicit `[]` resolved to a present empty list;
- stage replacement preserved `[Safety Authority]` order;
- scalar `context-sections: Authority` returned a `want sequence` error;
- one-buffer build resolution saw `authority-v1`, a subsequent valid README edit made one-buffer fetch resolution adopt `authority-v2`, and a subsequent `[Missing]` edit failed with `selector "Missing" matches 0 headings`.

The throwaway file was removed after the run. Its single-buffer resolver shape and the explicit-empty red case seed the first permanent implementation tests.

A second throwaway characterization ran `go test ./internal/dispatch -run TestIdeationSpikeCharacterizeLegacyFencedStage -v` against the real `extractStageSubsection` and passed. With a fenced `### ideation` before a real one, today's extractor returns the fenced span exactly as `"### ideation\n\nfenced-body\n```"`. That measured compatibility behavior is why cycle 2 reconciles raw byte spans instead of moving stage identity onto the fence-safe selector scanner. The throwaway file was removed; implementation first commits this result as a permanent pre-refactor characterization.

Cycle 3 exercised the load-bearing raw-span mapping with one buffer containing UTF-8 (`Wørkflow`, `Pärent`, `α`-`δ`), CRLF, lone CR, VT, a decorated stage, its level-4 child, a level-2 parent, and a disjoint level-2 sibling. A temporary bridge invoked the real fence-safe `scanHeadings(splitLines(...))`; the dispatch test invoked the real `extractStageSubsection` and independently mapped both shipped parsers' coordinates to source byte offsets.

The first run failed with a slice-bounds panic: a naïve rune-range mapper advanced past a CRLF pair at CR, then visited its LF byte again and produced reversed bounds. The mapper was corrected to decode manually and advance over CRLF as one separator, exactly like `splitTextLines`. The rerun command `go test ./internal/dispatch -run TestIdeationSpikeMixedRawSpanReconciliation -v` passed and observed:

- legacy rendered stage exactly ``### `build` *(captain)*\nstage β\ncontinuation\n\n#### Child\nchild γ``;
- its raw interval sliced exactly ``### `build` *(captain)*\rstage β\vcontinuation\r\n\r\n#### Child\r\nchild γ\r\r``;
- the fence-safe parent, child, and disjoint intervals each sliced their independently authored raw UTF-8/mixed-newline substrings exactly;
- stage/parent, stage/child, and selected-parent/selected-child intersections were true;
- stage/disjoint and selected-parent/disjoint intersections were false.

Both temporary spike files were removed. The CRLF double-advance failure and the five interval outcomes seed permanent reconciliation tests; implementation must not replace them with line-number comparisons.

### Expected surface

The reviewed reconciliation raises the expected implementation surface to 9-12 files and about 520-760 changed lines: 170-250 production lines in `internal/status/stages.go`, `internal/status/section_read.go`, `internal/dispatch/build.go`, `internal/dispatch/showstagedef.go`, and dispatch help; 320-460 fixture/test lines under `internal/status`, `internal/dispatch`, and existing integration fixtures; and 20-35 documentation lines. The cycle-2 tolerance is at most 14 files or 900 changed lines. Crossing either bound, adding a package, or touching a skill/runtime adapter requires returning to the captain before expanding scope.

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

**AC-1 (VALUE) - A dispatched worker receives the exact canonical bytes of the workflow policy selected for its stage, without copying that policy into the stage subsection.**
Verified by: a fixture with a unique raw policy `## Pølicy\r\nα\rβ\u0085γ\r\n \t\r\n` measures the current undeclared baseline as 0 policy bytes, then drives `dispatch build`, executes its emitted one-file bootstrap fetch, and asserts literal equality to `legacy stage bytes + "\n\n## Pølicy\nα\nβ\nγ\n"`. Removing the selector, emitting raw CR/CRLF/NEL, retaining the trailing whitespace-only line, changing a UTF-8 byte, using any separator other than two LF, or making the fetch stage-only fails the assertion.

**AC-2 - Workflow-owned section selection is generic, ordered, and bounded.**
Verified by: a custom-stage fixture selects two arbitrary heading names in reverse README order and asserts output follows declaration order, contains no unrelated section bytes, inherits a default when absent, replaces it at stage scope, and preserves explicit `[]` as a clear. A fenced lookalike heading is not selectable. Swapping order, merging defaults, collapsing absent with empty, leaking an unrelated section, or recognizing the fenced lookalike fails the fixture.

**AC-3 - Invalid section declarations stop dispatch before a worker artifact can be consumed.**
Verified by: table-driven `dispatch build` cases for malformed/wrong-kind YAML, a missing selector, repeated selector, ambiguous matching headings, selected parent containing the stage, selected child inside the stage, and selected parent containing another selected child all exit nonzero, write no usable dispatch body, and identify the README path, stage, selector/key, failure kind, and both spans for overlap. Disjoint siblings pass. Correcting each independent fixture makes only its corresponding case pass.

**AC-4 - Undeclared workflows preserve the stable dispatch and read contract byte-for-byte.**
Verified by: every existing dispatch golden remains unchanged; existing `show-stage-def` decorated-stage, malformed-heading, separator-set, CRLF, nested-heading, and no-declaration parity fixtures remain byte-identical. A new pre-refactor characterization pins the measured fenced-stage output `"### ideation\n\nfenced-body\n```"`, while a reconciliation fixture proves the fence-safe selector scanner does not resolve that fenced context lookalike. An explicit regression fixture asserts exactly one existing `show-stage-def` fetch command. Any extra flag, fetch, packet block, changed legacy byte, or status projection fails the fixed baseline.

**AC-5 - The pointer adopts valid current workflow context and refuses invalid current context at bootstrap.**
Verified by: one fixture runs build against policy v1, rewrites the README to a valid policy v2 and observes the unchanged emitted fetch return v2 exactly, then rewrites to a missing selector, wrong-kind sequence, and newly overlapping span and observes nonzero fetch results before any worker action. Returning v1 after the valid edit, mixing v1 metadata with v2 body, or accepting any invalid-current edit fails the test.

**AC-6 - Declared context is delivered through the same host-neutral bootstrap on every runtime.**
Verified by: the shared dispatch/integration fixture executes the emitted fetch command for Claude, Codex, and Pi envelopes and observes the same exact assembled bytes, while the existing approval-gated `claude-live`, `codex-live`, and `pi-live` lanes remain green because the host-neutral dispatch path changed. Removing the fetch from any host envelope or changing one host's quoting fails the fixture; an unapproved live lane is not counted as a pass.

## Test plan

Implementation starts with the AC-1 red behavior fixture. It authors a README with a custom stage, one unique top-level policy, one unrelated section, and no copied policy in the stage subsection. The policy mixes UTF-8, CRLF, lone CR, NEL, internal blanks, and trailing whitespace-only lines; the test builds the assignment, parses and executes its emitted fetch, and compares stdout to one independently authored LF-canonical literal including the two-LF component separator and one-LF final terminator. Estimated cost: medium, about 90-130 test lines.

Add focused tests around the existing parser and stage metadata:

- `internal/status`: strict in-memory YAML decoding with an error channel; absent/inherited, explicit-empty, replacement, ordering, and wrong-kind cases; exact half-open section spans; fenced-heading exclusion; and ambiguous-heading reporting. The permanent version of the cycle-2 spike must preserve present-empty as a non-nil empty slice or an equivalent explicit state. These exercise the existing parser; they do not add JSON fields or a second parser. Estimated cost: 130-180 test lines.
- `internal/dispatch`: build-time validation and exact `show-stage-def` assembly for two arbitrary headings in declaration order; a pairwise interval matrix covering stage/parent, stage/child, selected parent/child, exact duplicate, and disjoint siblings; and a single-buffer resolver check. Extend `TestBuildStageDisciplineRidesFetchNotInlineAssignment` or its fixture rather than creating a parallel bootstrap harness. Estimated cost: 170-240 test lines.
- Reconciliation: before refactoring, commit the cycle-2 fenced-stage characterization and cycle-3 mixed-buffer spike as permanent tests. Drive the in-memory legacy stage extractor through every existing decorated-stage, malformed-heading, VT/FF/FS/GS/RS/NEL/LS/PS separator, CRLF/lone-CR, nested-level-4, and no-declaration case without regenerating its golden. In the same raw source, prove a fenced context lookalike remains absent from the fence-safe selector table. Assert raw-byte spans from the two existing parsers, not their incompatible line numbers, drive overlap decisions.
- Rendering: table-drive selected raw sections using LF, CRLF, lone CR, and one exotic separator; assert UTF-8 and nonblank spaces survive, internal blanks survive, trailing whitespace-only logical lines do not, components have exactly two LF between them, and stdout has exactly one final LF. Expected strings are literals, not output from the production canonicalizer.
- Temporal behavior: build against v1; mutate the same README to valid v2 and execute the emitted fetch to prove live adoption; then mutate independently to wrong-kind, missing, and overlapping current states and require fetch failure. Each resolver invocation reads the README once so the test can detect version mixing.
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
+Selected sections normalize newline boundaries to LF, drop trailing blank
+lines, use one blank line between sections, and end stdout with one LF.
+The read uses the then-current valid README; an invalid current selection
+fails before stage work. Workflows without the field keep existing output
+unchanged.
```

Update `dispatch show-stage-def --help` from “Print the workflow README section for a stage.” to “Print a stage’s workflow README section followed by its declared context sections.”

### Feedback Cycles

- Cycle 1: REVISE — independent ideation staff review; surface 1 design file with implementation not started vs estimate 8-10 files/380-560 changed lines; AC unchanged. Resolve the existing stage-reader/fence-safe-reader semantic split, reject intersecting stage/selected spans, exercise the actual ordered tri-state YAML round-trip before approval, and state honest live-read semantics for valid edits between build and bootstrap.
- Cycle 2: REVISE — independent ideation staff re-review; surface 1 design file with implementation not started vs estimate 8-10 files/380-560 changed lines; AC unchanged. Exercise the proposed raw-byte-span reconciliation across UTF-8, CRLF/lone-CR, exotic separators, and overlap cases before approval; define and test exact selected-context trailing-blank and newline rendering for non-LF README bytes.

## Stage Report: ideation

- DONE: Define the builder-owned section-selection contract: workflow metadata declares ordered sections, dispatch build validates them, and the existing stage-definition read returns the selected context for the ensign's one-file bootstrap.
  Defined `context-sections`, inheritance/replacement/clear semantics, exact assembly order, shared build/read validation, diagnostics, and the unchanged single fetch command.
- DONE: Reuse the existing fence-safe section parser and existing dispatch/live fixtures; add no status projection, second parser, bootstrap multi-read protocol, or product-specific heading semantics.
  The design exposes only a narrow in-process selector from `internal/status/section_read.go`, extends the existing fetch fixture, and explicitly excludes every prohibited surface.
- DONE: Measure that declared policy reaches the worker exactly while undeclared workflows remain byte-identical, with a bounded expected surface and falsifiable test plan.
  AC-1 measures 0-to-authored-length exact policy bytes; AC-4 pins existing goldens; the plan bounds work to 8-10 files/380-560 changed lines with a 12-file/700-line tolerance.

### Summary

The ideation now specifies one workflow-owned ordered selector contract whose validation belongs to `dispatch build` and whose content reaches the ensign through the existing `show-stage-def` fetch. It preserves no-declaration bytes, reuses the fence-safe parser and current fixture/live lanes, and records exact failure, proof, documentation, and surface bounds.

## Stage Report: ideation (cycle 2)

- DONE: Unify or explicitly reconcile stage extraction and fence-safe section spans while preserving every existing decorated-stage, separator, CRLF, nested-heading, and no-declaration byte contract.
  Precisely reconciled the unchanged legacy stage extractor and fence-safe selector scanner through raw-byte spans from one buffer; AC-4 pins every named existing fixture without golden regeneration.
- DONE: Define and test fail-before-spawn overlap rules for stage/parent/child selected spans, and exercise the real absent/inherit, explicit-empty, replacement, ordering, and wrong-kind YAML round-trip as the ideation spike.
  Pairwise half-open-span rejection covers all requested intersections; the throwaway one-buffer Go spike red-marked nil-vs-empty, then passed inheritance, explicit clear, replacement order, and wrong-kind rejection.
- DONE: Correct the temporal contract to live-read semantics: build preflights one version, bootstrap validates and adopts the then-current valid content, and invalid current metadata fails.
  AC-5 now exercises v1 build, valid-v2 adoption, and independent wrong-kind/missing/overlap failures from one immutable README buffer per invocation.

### Summary

Cycle 2 reconciles the incompatible shipped parsers without changing either one's public semantics: raw-byte spans provide one overlap coordinate while the stage and selector authorities stay distinct. It adds structural overlap rejection, records the real metadata spike, and makes the pointer's live-read semantics explicit and falsifiable.

## Stage Report: ideation (cycle 3)

- DONE: Run and record a bounded spike using both shipped parsers over one raw buffer with UTF-8, CRLF/lone-CR, an exotic logical separator, and intersecting/disjoint headings; prove every mapped interval slices the intended bytes and yields the intended overlap result.
  The mixed-buffer spike caught and fixed a CRLF double-advance panic, then passed exact legacy-stage/parent/child/disjoint raw slices and the three-true/two-false overlap matrix.
- DONE: Define exact declared-context rendering for trailing blanks and LF/CRLF/lone-CR inputs, including the separator between assembled sections and final terminator.
  Selected spans use legacy split-lines normalization, trim only trailing whitespace-only logical lines, preserve other bytes/blanks, join components with two LF, and end stdout with one LF.
- DONE: Extend ACs and the existing fixture plan so the exact-byte value claim fails on a non-LF declared context without changing undeclared legacy output or adding a third parser.
  AC-1 now pins a literal UTF-8 CRLF/lone-CR/NEL payload; AC-4 retains every undeclared golden; the fixture plan promotes both bounded spikes into the existing status/dispatch lanes.

### Summary

Cycle 3 proves the raw-byte reconciliation on the real parsers and records the mapping bug the spike exposed. The design now defines one exact canonical output for mixed-newline context and a literal non-LF value fixture, while leaving the legacy no-declaration path untouched.

## Stage Report: implementation

- DONE: Before editing, declare intended files and LOC; implement the approved shared one-buffer resolver, ordered context-sections metadata, structural overlap rejection, live-read behavior, and exact selected-context rendering while preserving undeclared bytes.
  Commit `396c92a3` delivers the resolver and docs in 10 files/560 changed lines; exact-byte, ordered, tri-state, overlap, fenced-stage, and live-edit tests fail if those contracts regress.
- DONE: Promote the tri-state and mixed-buffer spikes into existing status/dispatch tests, extend existing host-envelope/golden fixtures, and make malformed, ambiguous, repeated, intersecting, and current-version invalid selections fail before worker action.
  `go test ./internal/status ./internal/dispatch ./skills/integration` passed; the tests use literal outputs and invalid mutations, while every undeclared dispatch golden remained unchanged.
- FAILED: Stay within 12 files/700 changed lines, add no status projection/snapshot/FO assembly/bootstrap multi-read/third parser/new harness, and run focused/full/race/docs plus required existing live lanes and final-tip Roborev.
  Scope stayed at 10 files/560 lines with no forbidden surface; full and race suites, `mkdocs build --strict`, and offline live coverage passed, but the approval-gated Claude/Pi runs stopped before FO work on revoked/refresh OAuth and still require authorized CI secrets before merge.
- DONE: Request Roborev at the final implementation tip and triage every finding before reporting completion.
  Branch-final job `1990` reviewed `396c92a3`; its one medium and three low non-blocking findings were closed with advisory Deferred-risk/Polish dispositions, four-field evidence, and promotion conditions, with AC unchanged.
- DONE: Exercise the existing runtime lanes against the committed candidate and preserve exact evidence.
  Codex local-auth ran all 9 shared scenarios green in 1299.00s; Claude artifacts are under `/tmp/spacedock-live-workflow-declared-section-reads/artifacts/claude`, Codex under `artifacts/codex-local-auth`, and pinned Pi under `artifacts/pi-pinned`.

### Summary

Commit `396c92a3` adds generic ordered workflow context pointers, reconciles the two shipped parsers through raw half-open spans from one README buffer, and preserves legacy output when no context is declared. Focused, full, race, docs, golden, integration, and Codex live proofs are green; Claude and Pi reached their real hosts but failed before workflow work on external credential state, so their approved CI-secret lanes remain required before merge.

## Stage Report: implementation (CI repair)

- DONE: Stay within 12 files/700 changed lines, add no status projection/snapshot/FO assembly/bootstrap multi-read/third parser/new harness, and run focused/full/race/docs plus required existing live lanes and final-tip Roborev.
  Authorized Runtime Live E2E [run 30087598672](https://github.com/spacedock-dev/spacedock/actions/runs/30087598672) completed SUCCESS at exact head `396c92a302e1ac220d552d6f16600e9f08ccb622`: offline, pi-live, codex-live, claude-live sonnet, and claude-live claude-opus-4-8 were green; journey-delta-comment skipped by workflow design.

### Summary

The authorized CI run repairs the implementation report's sole FAILED checklist item without changing product code or the reviewed implementation tip. The earlier local Claude revoked-token and Pi OAuth-refresh failures remain valid historical harness evidence, while the approved secret-bearing lanes now prove the candidate itself green on every required runtime.

## Stage Report: validation

- DONE: Independently reproduce every canonical AC at exact candidate 396c92a3: exact mixed-newline selected bytes, ordered tri-state metadata, overlap/malformed refusals, live-read mutation behavior, undeclared byte identity, and host-neutral one-fetch envelopes.
  Committed AC tests passed; throwaway literal-byte probes also covered all splitlines separators, UTF-8/nonblank spaces, EOF and endpoint siblings, no-artifact failure, and actual execution of each Claude/Codex/Pi envelope fetch.
- DONE: Audit the 10-file/560-line surface against the approved exclusions and adversarially attack parser-coordinate identity, EOF/separator boundaries, duplicate/intersecting spans, current-version changes, allocation/work scaling, and tests that could pass on the wrong bytes.
  Exact tip `396c92a3` is clean at 10 files/560 changed lines; no excluded surface appears, raw half-open coordinates and pairwise overlap held, and a 1,000-selector probe completed in 7.48 ms on Apple M4.
- DONE: Verify focused/full/race/docs/Roborev evidence and exact-head Runtime Live E2E run 30087598672; issue a fresh PASSED or REJECTED recommendation with every finding classified by defect kind and release scope.
  Focused, full, race, gofmt, and strict MkDocs gates passed; Roborev 1990 and CI run 30087598672 were inspected at exact SHA, with all five required CI jobs green. Recommendation: PASSED.
- DONE: AC-1 (VALUE) - A dispatched worker receives the exact canonical bytes of the workflow policy selected for its stage, without copying that policy into the stage subsection.
  `TestDeclaredContextBuildAndHostNeutralFetch` plus the independent literal probe fail on any UTF-8/newline/trim/separator/final-LF byte change; the dispatch artifact contains only the pointer.
- DONE: AC-2 - Workflow-owned section selection is generic, ordered, and bounded.
  Tri-state/reverse-order fixtures proved inheritance, replacement, explicit clear, unrelated exclusion, arbitrary heading names, fence safety, and declaration order.
- DONE: AC-3 - Invalid section declarations stop dispatch before a worker artifact can be consumed.
  Malformed, wrong-kind, missing, repeated, ambiguous, parent/stage, child/stage, and selected parent/child cases failed loudly; a fresh missing-selector probe observed empty stdout and zero matching dispatch artifacts.
- DONE: AC-4 - Undeclared workflows preserve the stable dispatch and read contract byte-for-byte.
  Existing build goldens and decorated/malformed/separator/CRLF/nested parity tests passed unchanged; independent decorated, fenced-stage, mixed-separator, and plain fixtures equaled the legacy extractor exactly.
- DONE: AC-5 - The pointer adopts valid current workflow context and refuses invalid current context at bootstrap.
  The temporal fixture adopted v2 after a v1 build and returned empty stdout/nonzero for current wrong-kind, missing, and newly overlapping edits.
- DONE: AC-6 - Declared context is delivered through the same host-neutral bootstrap on every runtime.
  A throwaway validator built the exact candidate and executed each emitted Claude/Codex/Pi shell command, obtaining identical canonical bytes; CI run 30087598672 independently kept all live host lanes green.
- DONE: Material findings — none; no outcome defect or evidence defect affects a supported AC.
  Exact-byte and lifecycle behavior are valid, and the fresh command-level proofs close the observable boundaries that substring-only assertions would miss.
- DONE: Deferred risk — outcome risk, non-material: a column-zero YAML comment shaped as `## Heading` can collide with a selected body heading.
  The crafted trigger is absent from templates/live evidence; supported body headings pass. Promote on a supported template, live workflow, user report, or explicit frontmatter-transparency promise.
- DONE: Deferred risk — outcome risk, non-material: strict parsing rejects the requested stage when another stage has wrong-kind `context-sections`.
  This requires an already-invalid workflow and matches the approved strict parse/AC-3; promote if partial operation with invalid unrelated metadata becomes supported.
- DONE: Deferred risk — outcome risk, non-material: inline state-branch fallback could emit `spacedock-state/.` for an underivable workflow basename.
  Resolved absolute workflow paths make the trigger unreachable in supported fixtures; promote if a supported root/dot path reaches dispatch build.
- DONE: Deferred risk — outcome risk, non-material: selector lookup and overlap checks are quadratic at very large declaration counts.
  The author-controlled unsupported stress shape remained 7.48 ms at 1,000 selectors; promote if supported workflows reach thousands of selected sections or resolver latency becomes observable.
- DONE: Polish finding — neither outcome nor evidence defect, non-material: `extractStageSubsection` remains production-unused.
  The wrapper supports an existing compatibility measurement and has no observable harm; remove only in later cleanup with that test migrated.

### Summary

Fresh validation recommends PASSED for exact candidate `396c92a302e1ac220d552d6f16600e9f08ccb622`. All six ACs have behavior-level evidence, required local and authorized CI gates are green, no material finding remains, and the advisory risks retain concrete promotion conditions without expanding this release.
