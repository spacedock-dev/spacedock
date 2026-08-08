---
title: Polish the v1 gate command and documentation surface
status: ideation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: top-level help omits gate prepare and the prose still describes prototype lifecycle details that the semantic cuts will remove."
started: 2026-08-08T00:04:11Z
completed:
verdict:
score: "0.8"
worktree:
issue:
pr:
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

The baseline is **16 semantic-coupling groups** across KD, 824, D8, A7, BV, and SK. Implementation reduces that count to zero. A group is one repeated assertion strategy, not every individual `strings.Contains` call.

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

**AC-1 (VALUE) - Semantic behavior has zero prose-search proof dependencies, down from the 16-group baseline.**
Verified by: diff-audit every inventory row; mutate or remove the claimed structured/state/behavior outcome and observe its focused test fail while equivalent prose rewording stays green. The only exact-text assertions added or retained by this cleanup are explicitly labeled published top-level and `gate --help` contracts.

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

## Stage Report: ideation

- DONE: Inventory prose-grep assertions introduced by durable-decisions PRs, including KD, 824, and D8, and classify each as published-text contract or semantic-output coupling.
  The task body records 16 coupling groups across KD, 824, D8, A7, BV, EP2, and SK; only deliberately published top-level and gate help remain exact-text contracts.
- DONE: Design the smallest replacements for semantic prose greps using structured fields, parsed sections, durable state, or behavioral outcomes without adding product behavior.
  Each inventory row names its replacement and falsifying outcome; explicit cuts forbid new schemas, parsers, commands, provider adapters, live lanes, and lint machinery.
- DONE: Reconcile the final help and documentation surface with all landed sprint tasks, and define exact files, LOC tolerance, acceptance evidence, and explicit cuts.
  The final scenarios, concrete doc edits, 24-34-file/220-520-line estimate, ±6-file and +250/−150-line tolerance, five ACs, and seven-step test plan form the gated baseline.

### Summary

Ideation turns F6C into a finite proof-hygiene cleanup: semantic prose-search dependencies fall from 16 groups to zero, while exact published help remains stable. The design also reconciles the self-syncing gate lifecycle, provider and eligibility cuts, pinned dispatch/recovery/conflict ownership, mergeability policy, canonical contract landing pass, and final documentation without adding behavior.
