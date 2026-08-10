---
title: Polish the v1 gate command and documentation surface
status: validation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: top-level help omits gate prepare and the prose still describes prototype lifecycle details that the semantic cuts will remove."
started: 2026-08-08T00:04:11Z
completed:
verdict:
score: "0.8"
worktree: .worktrees/spacedock-ensign-polish-v1-gate-command-surface
issue:
pr: "#667"
sprint: durable-decisions
id: f6cvn0s87ywbs158yy0b5q7k
gates:
    version: 1
    records:
        - id: gate:f6cvn0s87ywbs158yy0b5q7k:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:f6cvn0s87ywbs158yy0b5q7k-backlog-1
              briefing:
                id: briefing:f6cvn0s87ywbs158yy0b5q7k:backlog:attempt-1:revision-1
                digest: sha256:2372b3af80ff6dcfbc792369d495bde4f3beb4b74238a334109cb1c72a2e46c9
                request-digest: sha256:557bd043331c2db7efd24534289d48d21a537bd6cae11efe2db628b7e3f738f6
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:f6cvn0s87ywbs158yy0b5q7k:backlog:1
                briefing: briefing:f6cvn0s87ywbs158yy0b5q7k:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-07T23:59:41.992257Z"
                decision: approve
                reason: Captain directs F6C to kill prose-grep assertions introduced by durable-decisions PRs, including KD, 824, and D8, while retaining exact text checks only for published human output.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:f6cvn0s87ywbs158yy0b5q7k:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:f6cvn0s87ywbs158yy0b5q7k-ideation-1
              briefing:
                id: briefing:f6cvn0s87ywbs158yy0b5q7k:ideation:attempt-1:revision-1
                digest: sha256:a492a2061dd788377b7965037691eede103a410b69809755b50e5b8ccec5b42d
                request-digest: sha256:bbbaeddc4aa35f27528128af45f77b5b4713711ad1b420b55d34f0c06c667818
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:f6cvn0s87ywbs158yy0b5q7k:ideation:1
                briefing: briefing:f6cvn0s87ywbs158yy0b5q7k:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-08T00:29:17.546204Z"
                decision: approve
                reason: Captain directed F6C to kill sprint-added prose greps. The ideation report accounts for all three dispatched items, inventories 16 coupling groups, targets zero semantic prose-search dependencies, and preserves exact text only for published help.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:f6cvn0s87ywbs158yy0b5q7k:validation
          stage: validation
          attempts:
            - id: gate-attempt:f6cvn0s87ywbs158yy0b5q7k-validation-1
              briefing:
                id: briefing:f6cvn0s87ywbs158yy0b5q7k:validation:attempt-1:revision-1
                digest: sha256:57157057f0b460b56b1bb2bdc4ddb8c65716454726312e2b2474c94693afb781
                request-digest: sha256:4b5c1b80f5df0b812e816945087f5cc1c78e39754d85165a6b8e4c614bd1a69b
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:f6cvn0s87ywbs158yy0b5q7k:validation:1
                briefing: briefing:f6cvn0s87ywbs158yy0b5q7k:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T21:12:34.883305Z"
                decision: approve
                reason: Validation proves all 11 semantic prose-search dependencies in the finite 15-row inventory are removed or replaced, all five acceptance criteria have evidence, and the full required verification is green.
              application:
                target-stage: done
                state: pending
mod-block: merge:pr-merge
---

Make the stable help, command reference, specification, and First Officer instructions describe the final minimal gate lifecycle after the semantic cuts land. Remove sprint-added tests that infer behavior by searching free-form command, prompt, document, or transcript prose. Exact text remains asserted only for deliberately published human-facing help.

## Problem

Top-level `spacedock --help` currently lists only `gate record | validate | consume`; it omits the retained `prepare` and `withdraw` verbs even though `spacedock gate --help` exposes them. The canonical contract diagram still shows separate state commits after `record` and `consume`, despite those commands now committing and syncing their own split-root writes. The concept, command reference, roadmap journey, and First Officer instructions describe overlapping versions of the lifecycle.

The durable-decisions PRs also added tests that treat prose as a machine protocol. The worst cases search a generated assignment, Claude transcript, skill file, or specification for a phrase and then claim the corresponding dispatch, recovery, provider, or gate behavior occurred. Such tests can stay green when the behavior is broken and fail when equivalent wording changes.

## Proposed approach

1. Publish one retained verb list in top-level help: `prepare | withdraw | record | validate | consume`. Keep the full top-level help row and the existing `gate --help` usage block as the only new exact-text assertions because they are deliberately published human-facing command output. Behavioral tests continue to own exit codes, mutations, and state.
2. Replace every semantic prose-search group in the inventory below with the smallest existing evidence surface: decoded command JSON, a bounded Markdown section parser, canonical gate/frontmatter reads, Git state, invocation logs, or final filesystem behavior. Delete an assertion when stronger behavior evidence already proves the claim.
3. Reconcile the command reference, concept page, frontmatter reference, canonical contract, roadmap journey, merge mod, recovery guidance, and First Officer gate/dispatch instructions with the landed `0.27.0` behavior. The canonical contract remains the only lifecycle specification; other pages link to or summarize it.
4. Complete the re-homed contract landing pass: remove any residual owner/task aliases, use generic example identities, and render-check the Mermaid lifecycle. If the final candidate already has no residual scaffolding, record the zero-edit verification instead of creating churn.

No spike needed: current command fixtures already decode JSON envelopes and canonical gate state; `status --read` and existing section parsers prove bounded Markdown parsing; real Git fixtures already prove commits, branches, clean trees, and byte-clean refusal.

## Prose-coupling inventory and replacements

The finite inventory contains **15 rows**: **11 semantic-output couplings** require removal or replacement, and **four retained/non-semantic rows** remain as regression snapshots, published help, documentation-only reconciliation, or direct behavior. Implementation reduces semantic prose-search proof dependencies in this inventory to zero.

| Owner | Current coupling | Classification | Smallest replacement |
|---|---|---|---|
| KD `self-contained-ensign-dispatch` | Generated dispatch tests search raw bodies for stage/context sentinels, `### Fetch commands`, `### Workflow launcher`, standing-route headings, and a launcher-line regexp. | Semantic-output coupling. | Decode the `dispatch build` envelope, split the artifact into named sections, and compare each section with the source input; assert `fetch_commands` is empty and execute the parsed launcher command under A/B invocation logs. |
| KD | Twenty-six existing build snapshots changed with the assignment format. | Existing output snapshots, not behavioral proof. | Retain only as regression snapshots; no AC or behavior claim may cite their prose. F6 adds no new exact-text assertion to them. |
| KD | Error tests search stderr for the refusal sentence. | Semantic-output coupling unless wording is separately published. | Assert exit 1, no artifact path/file, and empty A/B invocation logs. Do not pin the sentence. |
| 824 `align-claude-break-glass-dispatch-oracle` | `assertBreakGlassObservables` searches report text and Agent prompt prose for `dispatch build`, ensign skill, stage definition, stage report, and completion signal. | Semantic-output coupling. | Keep typed tool-event fields only for bare/team mode and dispatch cardinality; prove completion through the committed marker, parsed Stage Report, clean entity path, and path-scoped Git commit. |
| 824 | Recovery completeness searches the whole entity for a marker, report heading, `DONE`, and `Summary`. | Semantic-output coupling. | Parse the final implementation Stage Report and require a DONE item plus nonempty Summary; verify the marker in the owned body span and the complete blob in a path-scoped commit. |
| 824 | Merged-dispatch tests search an artifact for skill-load and `SendMessage` phrases; unit tests search prompts for fixture slugs. | Semantic-output coupling. | Parse the packaged assignment sections/typed dispatch input, then use the durable worker result and selected-mode fields as the outcome oracle. |
| 824 | The retired degradation sentence is searched as a negative control. | Semantic-output coupling; the exact sentence is not published API. | Remove it. Selected bare mode plus absence of a recovery-skill tool event proves the retired route without freezing narration. |
| D8 `codify-conflict-owner-dispatch-handoff` | The live test searches raw Codex JSONL for `spawn_agent` and `followup_task`. | Semantic-output coupling to provider serialization. | Remove it. The marker on the stamped owner branch, expected Git author, unchanged authority bytes, aborted rebase, and committed follow-up are the behavior. |
| D8 | Fresh dispatch verification searches the generated body for tuple substrings. | Semantic-output coupling. | Decode the envelope and compare the bounded scope-notes section with the exact tuple input; keep worker name as a structured equality. |
| A7 `prove-or-cut-provider-backed-gate-closure` | `chat_gate_surface_test.go` searches two skills, the spec, and roadmap for required/forbidden phrases and treats matches as recorder convergence. | Semantic-output coupling. | Delete the prose oracle. Existing CLI tests must reject removed `--room`/provider forms byte-clean; decision recording and canonical state tests prove the sole semantic recorder. File-absence checks prove Result/inventory fixtures are cut. |
| BV `remove-standalone-gate-eligibility` | Status/next/field tests search JSON strings for application, target, readiness, and removed eligibility fields. | Semantic-output coupling to serialized JSON text. | Unmarshal the documented JSON objects and assert fields/absence structurally. Command-log argument inspection and removed-verb exit/state tests remain behavioral. |
| BV | Help and unknown-subcommand tests assert the retained verb list. | Published-text contract. | Consolidate into the exact top-level gate row and exact `gate --help` usage block; keep explicit comment that published help is the exception. |
| EP2 `avoid-unnecessary-pr-rebases` | No committed prose-grep test was introduced. | Documentation reconciliation only. | Keep the merge mod's behavioral boundary: clean mergeability preserves candidate SHA; conflict/unknown preserves pending authority and delegates owner selection without sprint aliases. |
| SK `gate-agent-ergonomics` | Readiness tests use JSON prefix/substring checks before also decoding the same output. | Semantic-output coupling to serialized JSON text. | Remove the prefix/substring checks and assert decoded `dispatchable`, `ready_gates`, and field projection values. The exact JSON golden remains machine-output compatibility evidence, not a prose semantic oracle. |
| SK | Cold-report readiness depends on a structurally valid committed Stage Report. | Direct behavior, but the fixture must not use loose whole-body phrase matching. | Continue through the production report parser and assert the emitted readiness plus dirty/malformed exclusions and durable Git state. |

## Final common scenarios

- **Nonterminal approval:** commit selected review inputs; `gate prepare`; commit the prepared binding; present it; `gate record --decision approve --consume` closes, syncs, consumes, and syncs; `dispatch build --stamp` enters the successor. The standalone `record` then `consume` path remains supported, but neither needs a separate state commit after a successful split-root write.
- **Terminal approval:** prepare, commit, present, then `gate record --decision approve --consume`. The record closes and syncs; consume returns `approved-awaiting-merge` without spending. `merge guard` spends and terminalizes only after delivery proof; `merge guard --rework` supersedes and routes through `feedback-to`.
- **Changed review before decision:** `gate withdraw --reason ...`, commit the withdrawal, then prepare and commit an ordinary successor. No Resolution or hold is fabricated.
- **Revise or hold:** `gate record` commits and syncs the Resolution. Revise follows the workflow feedback route; hold stops. Neither carries application metadata.
- **Presentation:** chat and Subspace present the same committed Briefing and return semantic decision/reason input. Only `gate record --decision` records it; stable v1 has no `--room`, Result, or inventory ingestion.
- **Conflict and recovery:** a moving-target conflict returns to the owner tuple stamped by the initial dispatch, preserving authority and requiring fresh evidence after reconciliation. Break-glass recovery preserves the selected bare/team mode and is successful only when its assigned durable result is committed.

## Out of scope

No new gate logic, schema field, compatibility decoder, provider integration, runtime event schema, transcript parser, linter, command, live lane, or second lifecycle specification belongs here. Do not redesign KD, 824, D8, A7, BV, EP2, or SK. If reconciliation exposes behavior that is absent rather than merely undocumented or weakly tested, route it to the semantic owner.

## Acceptance criteria

**AC-1 (VALUE) - Semantic behavior has zero prose-search proof dependencies in the finite inventory; exact text remains only for deliberately published help.**
Verified by: diff-audit all 15 inventory rows, including replacement or removal of all 11 semantic-output couplings; mutate or remove the claimed structured/state/behavior outcome and observe its focused test fail while equivalent prose rewording stays green. The only exact-text assertions added or retained by this cleanup are explicitly labeled published top-level and `gate --help` contracts.

**AC-2 - Every retained v1 gate verb is discoverable, and removed provider/eligibility verbs are absent.**
Verified by: byte-exact top-level and `gate --help` fixtures for `prepare | withdraw | record | validate | consume`, plus command-level unknown-verb and byte-clean state tests for removed `--room` and eligibility forms. Adding or removing a retained verb fails the published help fixture; accepting a removed form fails the behavior test.

**AC-3 - One canonical lifecycle and its summaries describe the final executable `0.27.0` behavior.**
Verified by: map every command transition in the canonical spec and the six scenarios above to focused command tests for prepare, withdraw, record/sync, consume/sync, terminal merge/rework, readiness, dispatch stamping, conflict ownership, and recovery. A documented transition without executable or durable-state evidence fails review.

**AC-4 - The shipped contract contains no sprint-owner scaffolding or prototype surface.**
Verified by: render the Mermaid lifecycle; inspect the rendered flow; check the spec for sprint aliases/owner tags and prototype commands; exercise the documented commands. The check is release evidence, not a new committed prose-grep test.

**AC-5 (SCOPE) - The cleanup changes no command grammar, stored format, authority, or runtime behavior.**
Verified by: production diff classification, unchanged canonical fixtures, and focused/full/race behavior suites. Any non-help production Go change, gates/frontmatter fixture change, new event/parser/schema, or changed exit/state outcome requires routing to its semantic owner.

## Test plan

Before editing, run the focused tests named by each inventory row to preserve a behavioral baseline. During implementation:

1. Update the exact published help fixture first; prove bare, `-h`, `help`, and `--help` remain byte-identical and the `gate` row contains all five retained verbs.
2. Replace KD artifact substring checks with decoded envelope plus bounded-section comparisons and A/B invocation logs. Reuse existing parsing helpers; add no generic Markdown package.
3. Replace 824 prompt/report prose checks with typed dispatch fields, parsed Stage Report, marker, clean path, and path-scoped commit evidence. Run its offline matrix; the existing selected live lane remains the only live proof and no new lane is added.
4. Remove D8 provider-JSONL event-name checks and retain its existing live durable outcome. Parse the fresh scope-notes section rather than searching the body.
5. Delete A7's contract prose oracle; run removed-command, provider-file-absence, canonical decision-recording, and byte-clean refusal tests. Decode BV/SK JSON structurally.
6. Apply the concrete documentation changes below, render the contract Mermaid, and map every remaining scenario to an executable test.
7. Run `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check`.

## Concrete documentation changes

- Top-level help: replace `gate        record | validate | consume` with `gate        prepare | withdraw | record | validate | consume`; keep the continuation line concise enough for the grouped layout.
- Canonical lifecycle diagram: replace `RECORD -> state commit -> CONSUME -> state commit` with the actual self-syncing record/consume boundary; retain the one explicit state commit after preparation.
- Concept and command reference: say that successful close/consume writes self-commit and self-sync in split-root workflows; describe `record --consume` as the shortest approval path; remove any separate post-close/post-consume commit instruction.
- Frontmatter reference: replace the stale room-backed Result sentence with semantic chat/Subspace presentation through `gate record --decision`; retain workflow-neutral correction-round recording only where the final command surface still exposes it.
- Roadmap common journeys: use the six scenarios above and remove superseded prototype/application/advisory phrasing. Do not create another normative lifecycle.
- First Officer gate/dispatch text: consume structured `sync=`/`phase=` and readiness fields; never scrape human prose. Keep KD's pinned launcher, 824's selected dispatch mode, D8's stamped owner tuple, BV's collapsed approval path, EP2's preserve-candidate mergeability decision, and SK's `needs-preparation` projection.
- Merge mod: remove sprint-local G3/D8 names from shipped prose. On conflict or unknown mergeability, preserve pending authority and hand ownership to the consuming workflow without rebase, automatic resolution, or force.
- Contract landing pass: remove residual owner labels or sprint ids, genericize example ids, and record the Mermaid render result. If none remain on the implementation base, state that explicitly.

## Expected surface

Expected changed surface is **24-34 files / 220-520 changed lines**, net negative in semantic-oracle code. The likely files are:

- CLI/help: `internal/cli/help.go`, `internal/cli/cli_test.go`, `internal/cli/gate_test.go`, and the BV JSON/help tests in `gate_ceremony_count_test.go`, `terminal_consume_test.go`, and `verbs_test.go`.
- Semantic-oracle cleanup: KD dispatch tests (including `self_contained_assignment_test.go`), 824 recovery assertion/unit/live tests, D8 `conflict_owner_handoff_live_test.go`, A7 `chat_gate_surface_test.go`, and SK `gate_readiness_needs_preparation_test.go`.
- Canonical/docs: `docs/specs/gate-resolution-frontmatter-contract.md`, `docs/site/concepts/gates-and-decisions.md`, `docs/site/reference/{command-reference,frontmatter-contract}.md`, `docs/roadmap/durable-decisions/index.md`, and `docs/runtime-live-ci.md`.
- Operational text: `skills/{present-gate,fo-gate-lifecycle,fo-dispatch-recovery}/SKILL.md`, `skills/first-officer/references/{first-officer-shared-core,fo-dispatch-core}.md`, and `mods/pr-merge.md`.

Tolerance is **±6 files and +250/−150 changed lines** because KD/824/D8 may land with small test-file movement before F6 starts. The reset trigger is semantic, not numeric: any new production behavior, command/schema/event surface, compatibility path, provider integration, or second lifecycle spec returns to ideation regardless of LOC.

## Explicit cuts

- No committed repository-wide prose linter or grep-ban test; this is a finite sprint cleanup inventory.
- No transcript grammar, provider-event dialect adapter, actor-attribution observer, or model-output parser.
- No replacement for A7's docs/skill phrase oracle; behavior and removed-surface tests are sufficient.
- No new JSON field merely to make a test easier; decode fields that already exist.
- No new live scenario. Reuse 824 and D8 lanes only where live runtime behavior is already the acceptance claim.
- No cleanup of unrelated historical prose assertions outside the durable-decisions diff unless a touched test directly depends on them.

## Captain ruling: inventory accounting

The historical “16” starting count is unsupported bookkeeping, not a value or validation boundary. The reconstructible inventory is 15 rows: 11 semantic-output couplings requiring removal or replacement and four retained/non-semantic rows. This ruling supersedes contrary count claims in the historical stage reports below without rewriting those reports.

## Stage Report: ideation

- DONE: Inventory prose-grep assertions introduced by durable-decisions PRs, including KD, 824, and D8, and classify each as published-text contract or semantic-output coupling.
  The task body records 16 coupling groups across KD, 824, D8, A7, BV, EP2, and SK; only deliberately published top-level and gate help remain exact-text contracts.
- DONE: Design the smallest replacements for semantic prose greps using structured fields, parsed sections, durable state, or behavioral outcomes without adding product behavior.
  Each inventory row names its replacement and falsifying outcome; explicit cuts forbid new schemas, parsers, commands, provider adapters, live lanes, and lint machinery.
- DONE: Reconcile the final help and documentation surface with all landed sprint tasks, and define exact files, LOC tolerance, acceptance evidence, and explicit cuts.
  The final scenarios, concrete doc edits, 24-34-file/220-520-line estimate, ±6-file and +250/−150-line tolerance, five ACs, and seven-step test plan form the gated baseline.

### Summary

Ideation turns F6C into a finite proof-hygiene cleanup: semantic prose-search dependencies fall from 16 groups to zero, while exact published help remains stable. The design also reconciles the self-syncing gate lifecycle, provider and eligibility cuts, pinned dispatch/recovery/conflict ownership, mergeability policy, canonical contract landing pass, and final documentation without adding behavior.

## Stage Report: implementation

- DONE: Reduce the 16 inventoried semantic prose-search proof groups to zero, retaining exact text only for deliberately published top-level and gate help.
  Commit e481864e4 replaces KD body sentinels with decoded envelopes and bounded sections; 824/D8 prompt and provider-event prose with typed dispatch, parsed report, exact scope input, and Git outcomes; A7's prose oracle is deleted; BV/SK JSON is decoded structurally. Removing any owned fetch/checklist/scope field, selected dispatch shape, standalone worker marker, report entry, clean path-scoped commit, or readiness value now fails its focused test, while equivalent narration changes do not.
- DONE: Reconcile help, canonical lifecycle, docs, and First Officer guidance with the final landed 0.27.0 behavior without changing command grammar, stored state, authority, or runtime behavior.
  Top-level and gate help publish all five retained verbs; the canonical Mermaid, concept/reference pages, roadmap journeys, and runtime proof registry describe prepare's one explicit commit and record/consume self-sync. The contract contains no sprint aliases or owner labels; `mods/pr-merge.md`, recovery mode selection, stamped owner handoff, ready-gate ordering, and split-root sync guidance already matched the final surface and required zero churn.
- DONE: Stay within the approved finite cleanup surface and explicit cuts; run focused behavior controls, formatting, diff checks, full tests, and race tests with exact files/LOC reported.
  The final surface is 23 files / 770 changed lines (275 additions, 495 deletions), exactly the declared upper LOC tolerance and within the 18-40-file tolerance; only `internal/cli/help.go` changes production Go. `gofmt -w ./cmd ./internal`, focused package controls, `go test ./...`, `go test ./... -race`, `git diff --check`, exact help execution, and a rendered 28,651-byte Mermaid SVG all pass.

### Summary

Implementation removes every inventoried prose-dependent semantic oracle and keeps only the deliberately published help text exact. Commit e481864e4 reconciles the self-syncing v1 lifecycle and operational guidance without changing command grammar, stored authority, or runtime behavior; full and race suites are green.

## Review-finding disposition

### V1 — AC-1 baseline count is not reconstructible

- Reviewer observation: the inventory contains 15 data rows, not 16 groups. Eleven rows are explicitly classified as semantic-output coupling; four are retained or non-semantic (KD snapshots, BV published help, EP2 documentation-only, and SK direct cold-report behavior).
- Released user and normal workflow: validation uses this finite inventory to certify that the durable-decisions proof estate has no remaining prose-search dependency.
- Observable harm: the gate can claim “all 16 to zero” without identifying all 16 inputs, so neither inventory completeness nor the zero endpoint can be independently reproduced.
- Affected authority: value-ac[AC-1] Semantic behavior has zero prose-search proof dependencies, down from the 16-group baseline.
- Trigger evidence: a bounded table parse reports `inventory_rows=15 semantic_rows=11 retained_or_nonsemantic_rows=4`; state history shows ideation commit `bf4e56336` introduced the same unsupported 16-group sentence and 15-row table together.
- Validator proposal: Evidence defect; Material; ownership is ideation/captain because changing the baseline or AC changes approved proof scope; disposition is Needs decision and a scope/design reset, not candidate mutation or an automatic implementation feedback cycle.

## Stage Report: validation

- FAILED: Audit all 16 inventory groups and prove semantic claims now use structured, parsed, durable-state, or behavioral evidence; allow exact text only for published help.
  Candidate commit `e481864e4` replaces every explicitly inventoried semantic row with decoded envelopes, bounded sections, typed events, parsed reports, Git outcomes, or command behavior, but the promised baseline is unreproducible: 15 rows, 11 semantic and 4 retained/non-semantic, with no splitting rule that yields 16.
- DONE: Verify help and lifecycle documentation match executable 0.27.0 behavior while production changes remain limited to published help and no command/state/authority/runtime semantics change.
  Executed top-level and `gate --help` show all five retained verbs; CLI lifecycle tests exercise removed-form byte-clean refusal, prepare/withdraw, record/consume sync, terminal merge/rework, readiness, stamping, conflict ownership, and recovery; only `internal/cli/help.go` is non-test production Go.
- DONE: Reproduce focused controls, formatting, full, race, diff, help, and Mermaid rendering evidence; verify the 23-file/770-line boundary and recommend PASSED or REJECTED.
  Focused CLI/dispatch/ensigncycle/status/contractlint packages, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` pass; the clean worktree remains unchanged, the diff is 23 files/770 lines, and Kroki renders the lifecycle to a 28,651-byte SVG. Recommendation: REJECTED.
- FAILED: AC-1 — Semantic behavior has zero prose-search proof dependencies, down from the 16-group baseline.
  Mutation-sensitive replacements are present for all 11 explicitly semantic rows, but the unenumerated 16-group baseline makes the completeness and zero claims non-reproducible.
- DONE: AC-2 — Every retained v1 gate verb is discoverable, and removed provider/eligibility verbs are absent.
  `TestTopLevelHelpFormsAreIdentical`, `TestGatePrepareCLIPrintsExactRoomBindingAndCurrentV1HelpSurface`, and removed-provider/unknown-verb controls fail if published verbs drift or removed grammar writes state.
- DONE: AC-3 — One canonical lifecycle and its summaries describe the final executable `0.27.0` behavior.
  Focused command tests cover preparation, withdrawal, self-synced close/consume, collapsed approval, terminal merge/rework, readiness, dispatch stamping, owner handoff, and durable recovery; changing any transition or terminal state fails its outcome assertion.
- DONE: AC-4 — The shipped contract contains no sprint-owner scaffolding or prototype surface.
  The Mermaid service parsed and rendered the full lifecycle, its SVG carries the recorder-close and successor-sync nodes, and the canonical spec contains no listed sprint aliases or owner tags.
- DONE: AC-5 — The cleanup changes no command grammar, stored format, authority, or runtime behavior.
  Production classification finds only the published help row changed; canonical behavior suites, the full suite, and the race suite remain green with no worktree mutation.

### Summary

The candidate is behaviorally clean, within its exact 23-file/770-line boundary, and all executable, formatting, help, lifecycle, and render checks pass. Validation nevertheless recommends REJECTED because AC-1's claimed 16-group baseline cannot be reconstructed from its 15-row inventory; resolving that evidence defect changes approved proof scope and requires an ideation/captain design reset rather than candidate edits.

### V1 disposition update — resolved by captain ruling

The `## Captain ruling: inventory accounting` section resolves V1 by declaring the reconstructible finite boundary: 15 rows, 11 semantic-output couplings, and four retained/non-semantic rows. The historical 16 count is explicitly non-authoritative bookkeeping, so V1 no longer identifies a failed value AC or material release finding; candidate `e481864e4` remains unchanged.

## Stage Report: validation (cycle 2)

- DONE: Regrade the existing validation evidence against the Captain-corrected active F6C state without rerunning green checks or mutating the candidate.
  HEAD remains `e481864e4`; the prior focused, formatting, full, race, diff, help, and 28,651-byte Mermaid render evidence is reused unchanged, and the captain ruling removes the sole prior evidence defect.
- DONE: AC-1 — Semantic behavior has zero prose-search proof dependencies in the finite inventory; exact text remains only for deliberately published help.
  Diff audit covers all 15 rows: KD uses decoded envelopes, bounded sections, exact source comparison, invocation logs, and byte-clean refusal; 824 uses typed events, parsed Stage Reports, durable markers, clean paths, and Git commits; D8 uses durable owner outcomes and exact bounded scope input; A7's prose oracle is deleted in favor of CLI/state/file-absence behavior; BV and SK decode JSON structurally. The four retained rows are snapshots not cited as behavior, deliberately published help, documentation-only EP2 reconciliation, and production-parsed cold-report behavior.
- DONE: AC-2 — Every retained v1 gate verb is discoverable, and removed provider/eligibility verbs are absent.
  `TestTopLevelHelpFormsAreIdentical` and `TestGatePrepareCLIPrintsExactRoomBindingAndCurrentV1HelpSurface` pin the five published verbs; `TestGateRecordRejectsProviderRoomBeforeMutation` and unknown-subcommand controls fail if removed provider or eligibility grammar is accepted or writes state.
- DONE: AC-3 — One canonical lifecycle and its summaries describe the final executable `0.27.0` behavior.
  Existing focused evidence maps prepare/withdraw, record/consume sync, collapsed approval, terminal merge/rework, readiness, `dispatch build --stamp`, stamped conflict ownership, and durable recovery to exact exit, state, Git, or worker outcomes; a transition or terminal-state change breaks those assertions.
- DONE: AC-4 — The shipped contract contains no sprint-owner scaffolding or prototype surface.
  Prior release evidence rendered the canonical Mermaid to a 28,651-byte SVG and inspected the recorder-close, consume, and successor-sync flow; the spec audit found no listed sprint aliases, owner tags, `--room`, Result, inventory, or eligibility prototype surface.
- DONE: AC-5 — The cleanup changes no command grammar, stored format, authority, or runtime behavior.
  The 23-file/770-line diff changes only `internal/cli/help.go` in non-test production Go; unchanged canonical behavior fixtures plus focused, full, and race suites passed, and the candidate worktree stayed byte-clean.
- DONE: Recommend PASSED or REJECTED against the corrected boundary.
  Recommendation: PASSED. All five corrected ACs have reproducible evidence, the earlier V1 count finding is resolved by the captain's authoritative accounting ruling, and no material or deferred finding remains.

### Summary

Re-review recommends PASSED for unchanged candidate `e481864e4`. The captain-corrected 15/11/4 inventory makes AC-1 reconstructible, and the existing executable, structured, durable-state, scope, and rendering evidence satisfies AC-1 through AC-5 without another test run or candidate mutation.
