---
id: 8bnkrtq4rw46xkbez5zrbmmj
title: Codex keep-moving durable-evidence attribution false-red
status: implementation
source: "PR #513 Runtime Live E2E run 29392675038, codex-live job 87279446937"
started: 2026-07-30T23:05:41Z
completed:
verdict:
score: 0.8
worktree: .worktrees/spacedock-ensign-codex-keep-moving-durable-evidence-attribution-flake
issue:
milestone: 0.25.0
gates:
    version: 1
    current:
        gate: gate:8bnkrtq4rw46xkbez5zrbmmj:validation
    records:
        - id: gate:8bnkrtq4rw46xkbez5zrbmmj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:8bnkrtq4rw46xkbez5zrbmmj-backlog-1
              briefing:
                id: briefing:8bnkrtq4rw46xkbez5zrbmmj:backlog:attempt-1:revision-1
                digest: sha256:f74c978f7bac349914d3380c445388cf82767c0647d3267768c425649fb719a8
                digest-domain: canonical-bytes
                request-digest: sha256:83ed6b8050afb6ea2e3737bf14a7c57f164095922ad8f37db84546e5cde6f84d
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:8bnkrtq4rw46xkbez5zrbmmj:backlog:1
                briefing: briefing:8bnkrtq4rw46xkbez5zrbmmj:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-30T22:50:16.874174Z"
                decision: approve
                reason: Captain explicitly directed dispatch of 8b after Q3 merge and ruled that transcript-grammar parsing is offrail; ideation must replace the dialect proposal with behavior and durable-state proof.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:8bnkrtq4rw46xkbez5zrbmmj:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:8bnkrtq4rw46xkbez5zrbmmj-ideation-1
              briefing:
                id: briefing:8bnkrtq4rw46xkbez5zrbmmj:ideation:attempt-1:revision-1
                digest: sha256:a9e29deee9056d26159a1c6772e5b213e73467dd03ce07cfe9186d9660bbdde2
                digest-domain: canonical-bytes
                request-digest: sha256:1e8eef615d1bd19e2011de8bac8ca2650c19caa690e13e9a85903a430c39ba30
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:8bnkrtq4rw46xkbez5zrbmmj:ideation:1
                briefing: briefing:8bnkrtq4rw46xkbez5zrbmmj:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T23:04:28.063854Z"
                decision: approve
                reason: 'Captain conn accepted the no-transcript, net-negative design: ordered per-task Git history directly tests the keep-moving journey, adversarial controls can falsify it, and the implementation boundary forbids observer dialects or product semantics.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:8bnkrtq4rw46xkbez5zrbmmj-ideation-2
              briefing:
                id: briefing:8bnkrtq4rw46xkbez5zrbmmj:ideation:attempt-2:revision-1
                digest: sha256:bd9ac766740cfa3c2c1c43e7a54076063521e8f124824f2116c3c7d6b14b46fe
                digest-domain: canonical-bytes
                request-digest: sha256:1ec520cb9a98aeb59797200800c6dc1fcb06a3a07ad4f84717b1c1cac8ec24b3
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:8bnkrtq4rw46xkbez5zrbmmj:ideation:2
                briefing: briefing:8bnkrtq4rw46xkbez5zrbmmj:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-31T00:05:30.189465Z"
                decision: approve
                reason: 'Captain conn approves the retained-history correction: exact live Git evidence proves the only false red is same-task review sidecars at archive/held boundaries, while dispatch/report/terminal scope and all foreign-path controls remain strict and the parser deletion stays net-negative.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:8bnkrtq4rw46xkbez5zrbmmj-ideation-3
              briefing:
                id: briefing:8bnkrtq4rw46xkbez5zrbmmj:ideation:attempt-3:revision-1
                digest: sha256:33afe927a33a75aecb584781c9baccd3a047443a10c672ae1eb9ed184f27e8dd
                digest-domain: canonical-bytes
                request-digest: sha256:7aeb5e296e4be785452ad21b8391f18693898fc8178c746b7966dfaf0d698af3
                room-ref: ./review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:8bnkrtq4rw46xkbez5zrbmmj:ideation:3
                briefing: briefing:8bnkrtq4rw46xkbez5zrbmmj:ideation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-07-31T00:39:19.114265Z"
                decision: approve
                reason: Captain conn approves the dual-episode design because retained supported histories prove both ordinary and atomic first-worker forms, while parent/child blob checks and foreign/stale controls preserve actual worker proof without transcript/provider observation.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:8bnkrtq4rw46xkbez5zrbmmj:validation
          stage: validation
          attempts:
            - id: gate-attempt:8bnkrtq4rw46xkbez5zrbmmj-validation-1
              briefing:
                id: briefing:8bnkrtq4rw46xkbez5zrbmmj:validation:attempt-1:revision-1
                digest: sha256:1288b9c4b1ef14bc523f1a3a2fbd05d9290dbf29ef90f04c739b02d455e162a5
                digest-domain: canonical-bytes
                request-digest: sha256:585b217773f5644bf9042deb977dc7c15a65b941250d86b2d2752d31237f970b
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:8bnkrtq4rw46xkbez5zrbmmj:validation:1
                briefing: briefing:8bnkrtq4rw46xkbez5zrbmmj:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-31T01:45:21.357752Z"
                decision: revise
                reason: 'Rejected under sprint conn: validation reproduced four Material task-owned false greens—serialized independent work, preexisting-room rebind, nonterminal canonical archive, and unrelated questioned-ticket mutation—so AC-1, AC-2, and AC-4 are not proven.'
              application:
                action: feedback
                target-stage: implementation
                state: pending
---

The keep-moving live scenario must credit each completed task from its own durable workflow journey, without making a provider transcript dialect part of the product verdict.

## Problem

PR #513's `codex-live` lane failed `TestLiveCodexSharedScenarios/keep-moving-posture` even though `approved-gate`, `ready-one`, and `ready-two` had each completed the real workflow journey: dispatch entry, worker completion, terminalization, and archive. The observer instead tried to reconstruct that journey from Codex JSONL, report-heading cardinality, command text, and merge-loop output. Harmless changes in those representations produced a false red twice after the same exact-head behavior had passed offline, install, Claude/Opus, and Pi validation.

Candidate `1dd01b6db` replaced that observer with ordered Git history, but its first live diagnosis was destroyed by `t.TempDir` cleanup and could not justify changing the oracle. A second exact-tip run retained the completed workflow and narrowed the false red: the unmodified candidate credited `ready-one` and `ready-two` but reported `2/3` because `approved-gate`'s canonical archive commit also moved two gate-review files owned by that entity. Treating every path beyond `<slug>.md` as foreign conflates a same-entity atomic archive with cross-attribution.

Candidate `ae6cf1ee9` fixed that scope distinction and retained failing roots, then exposed a second supported durable shape. Its one exact-head live run used prescribed break-glass dispatch after the dispatch helper rejected empty checklists. Each task's first entity-file-only worker commit atomically added non-empty `started` plus the new current-stage Stage Report; a later entity-owned archive added terminal fields. With no separate `dispatch:` commit, the ordinary-only oracle returned `0/3` despite concrete worker output.

The architectural constraint remains: adding a generic-heading exception, recognizing a shell variable, normalizing `item.completed`, or replaying a smaller JSONL sample would retain a second protocol whose grammar changes independently of Spacedock behavior. The correction must stay inside existing Git paths, blobs, and ancestry.

Exact evidence: [PR #513](https://github.com/spacedock-dev/spacedock/pull/513), [Runtime Live E2E run 29392675038 / codex-live job 87279446937](https://github.com/spacedock-dev/spacedock/actions/runs/29392675038/job/87279446937).

## Spike result

An isolated harness-only copy of exact candidate `1dd01b6db` replaced only the keep-moving `t.TempDir` allocation with a caller-named retained directory. A freshly built `0.27.0-pre2+dev` binary passed `doctor --host codex --plugin-manifest .codex-plugin/plugin.json`; the supported local-auth Codex run then completed and left the real repository at `/tmp/spacedock-km-retained.SHtfgD` for inspection before cleanup.

The retained histories are:

- `approved-gate`: `1149540 dispatch` → `622afb9` worker report → `bb4cf57` terminalize → `8c1a40a` archive. The archive commit renamed `_archive/approved-gate.md` plus two `_archive/approved-gate/review/...` briefing files and no foreign entity path.
- `ready-one`: `0ec6813 dispatch` → `75f8385` worker report → `dec565a` terminalize → `50420a3` single-file archive.
- `ready-two`: `0cf96e0 dispatch` → `aeb0f4c` worker report → `66fecd9` terminalize → `6899d54` single-file archive.
- `questioned`: a scoped ideation report was followed by review and `f9c71be` binding the corrected review; that final commit changed `questioned.md` plus only `questioned/review/...` briefing files, and the entity remained active/nonterminal at review.

All three dispatch blobs have `status: implementation` and the same non-empty `started: 2026-07-30T23:53:21Z`; all three worker commits contain the later implementation Stage Report; all three terminal blobs have `completed: 2026-07-30T23:56:20Z` and `verdict: passed`. Thus the prior destroyed run's `0/3 missing dispatch` diagnosis was not reproduced and is not a design premise.

A throwaway two-scope spike kept dispatch, worker-report, and terminal commits strictly entity-file-scoped, but allowed archive/held-state commits to change only paths inside that slug's active or archived namespace. It credited the retained journey `3/3`, preserved all existing missing/reordered/stale/cross-report controls, and kept a new cross-attributed archive control red:

`SPACEDOCK_KEEP_MOVING_RETAIN_ROOT=/tmp/spacedock-km-retained.SHtfgD go test ./internal/ensigncycle -run '^(TestDurableTaskJourneys|TestRetainedKeepMovingJourney)$' -count=1 -v`

The spike passed in 6.776s. Final archived state alone remains insufficient; the smallest correction is a path-ownership distinction at the two commits that legitimately carry same-entity support files.

Cycle 3 retained the clean failed repository at `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/spacedock-keep-moving-1137023568`. Its completed histories are:

- `approved-gate`: `053a50e init` → `9d710f1` first worker → `363cbe0` terminal archive with same-task review sidecars.
- `ready-one`: `053a50e init` → `0c7c6b8` first worker → `8fa4c7b` terminal archive.
- `ready-two`: `053a50e init` → `b4d7fdf` first worker → `08924c9` terminal archive.

In every first-worker commit, the immediate parent lacks both `started` and the implementation Stage Report; the entity-file-only child changes to the expected stage and adds non-empty `started` plus that new report. A throwaway exact-`ae6cf1ee9` spike accepted both this retained `3/3` history and the ordinary two-commit fixture, while rejecting preexisting started, preexisting report, report-only, pre-dispatch/stale report, foreign worker path, and slug-prefix worker path controls:

`SPACEDOCK_KEEP_MOVING_RETAIN_ROOT=/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/spacedock-keep-moving-1137023568 go test ./internal/ensigncycle -run '^(TestDurableTaskJourneys|TestRetainedAtomicWorkerJourney)$' -count=1 -v`

The spike passed 16 deterministic subtests plus the retained live journey in 11.229s. It added 57/deleted 7 lines only in a throwaway copy of `shared_keep_moving_durable_test.go`; frozen candidate `ae6cf1ee9` remained untouched.

## Proposed approach

Keep the candidate's host-neutral Git-history oracle, but recognize exactly two worker-proof episodes for the expected stage:

1. **Ordinary:** an entity-file-only `dispatch: {slug} entering {stage}` commit whose blob has the expected stage and non-empty `started`, followed by a later entity-file-only commit whose blob contains the current-stage Stage Report.
2. **Atomic break-glass:** one entity-file-only first-worker child whose immediate parent has empty `started` and no current-stage Stage Report, and whose child has the expected stage and atomically adds both non-empty `started` and the new current-stage Stage Report. That commit satisfies dispatch plus worker completion because the new worker report is present; a generated artifact, started-only transition, or report-only child never earns credit.

After either episode, require a later terminal/archive boundary. It may be a separate entity-file-only terminal commit followed by the final entity-owned archive, or the final entity-owned archive may atomically add non-empty `completed` and `verdict`. In both forms the canonical archive Markdown must exist and the active Markdown must not.

Use path-component boundaries, not string-prefix aliases: `ready` does not own `ready-two.md` or `ready-two/...`. Dispatch and worker proof always remain entity-file-only. Only the final archive and final durable `questioned` hold may include paths inside that exact slug's active/archive support subtree; any other path makes the commit cross-attributed.

Grade the three independent tasks as a set of three separate journeys. `questioned` remains a negative guard: it must be durably re-shaped and nonterminal, and none of its commits may satisfy another slug. The live runner supplies only the workflow root and expected slugs/stages; it does not supply JSONL, final narration, commands, or provider events.

Extend the deterministic real-Git table with an atomic three-task positive that mirrors the retained first-worker plus terminal-archive shape. Add atomic controls for preexisting started, preexisting report, report without started, a foreign entity path, and a slug-prefix support path. Keep the ordinary positive and its missing dispatch/report/terminal/archive, stale/reordered report, foreign report/archive, and slug-prefix archive mutations.

Give the Codex keep-moving live fixture a failure-preserving temp root: delete it on success, but on failure retain and print the ordinary Git repository path. This is diagnosis only; the oracle still grades the in-place repository. Reuse the same durable completion oracle for the commissioned-task completion fallback in the smallest-sufficient-mechanism scenario.

The rejected alternatives are: changing the shared fixture to force ordinary dispatch (hides a supported break-glass shape), accepting `started` or archive state without a newly added worker report, allowing support subtrees at worker commits, final archive state alone, `status --archived` alone, an instrumented wrapper log, and any transcript/event/command parser.

## Deletion inventory

- Delete all 433 lines of `internal/ensigncycle/shared_keep_moving_test.go`: provider event decoding, command token/regex inference, narration inference, host-neutral motion trace, and transcript correlation tests.
- Delete all 409 lines of `internal/ensigncycle/shared_keep_moving_negative_test.go`: fabricated Claude/Codex streams, replay dialects, and final-message grammar controls.
- Delete all 202 lines of `internal/ensigncycle/codex_dispatch_evidence_test.go`: `item.completed` decoding, dispatch-build result decoding, Stage Report cardinality, status text matching, merge-output matching, and shell-loop variable parsing.
- Delete all 350 lines of `internal/ensigncycle/codex_dispatch_evidence_regression_test.go`: JSONL constructors and observer compatibility cases. Move only the generic `codexCommandOutput` fixture constructor if its unrelated round-recording consumer still needs it.

No compatibility is retained for these internal observer formats.

## Expected surface and semantic boundary

Expected cumulative files remain the same 10-file candidate surface. Cycle-3 implementation may modify only existing `internal/ensigncycle/shared_keep_moving_durable_test.go` and `docs/runtime-live-ci.md`; the other eight candidate files, including the shared fixture and failure-retaining Codex runner, remain byte-identical to `ae6cf1ee9`.

Incremental budget from `ae6cf1ee9`: at most 65 inserted lines across those two files, no new/helper file. Cumulative hard bounds versus `origin/main`: exactly 10 files, at most 400 insertions, at least 1,449 deletions, and at least 1,049 net deletions. There is no tolerance; stop before crossing any bound.

Observable semantics changed: test-oracle runtime behavior only. The oracle credits either a separate ordinary dispatch/report episode or an atomic first-worker episode proven by the parent/child entity blobs. Terminal fields may arrive in the final entity-owned archive. Same-task archive/held ownership and failure retention from `ae6cf1ee9` remain. Command grammar, CLI output, stored formats, mutation authority, runtime dispatch behavior, fixture behavior, retry policy, and provider adapters do not change.

Concrete documentation diff in `docs/runtime-live-ci.md`: replace the single “entity-file-scoped dispatch with started, worker Stage Report” sequence with “worker proof is either an entity-file-only dispatch followed by a later report, or one entity-file-only child that adds started and the new report over a parent with neither.” Retain the existing same-slug archive/held and failed-root sentences.

## Out of scope

- The permission-narration oracle removal in #514.
- The separate Codex rejection-flow worker-reuse flake.
- Retry policy or the task-15 wait-watchdog replacement.
- Changes to the keep-moving behavioral contract or its requirement to commission each ready entity.
- Runtime host adapters, dispatch tools, workflow frontmatter, and Stage Report format.

## Acceptance criteria

**AC-1 - Every completed independent task is credited from its own durable journey.**
Verified by: ordinary and atomic deterministic real-Git positives each report `3/3`; removing any one worker-proof episode reports `2/3`. The retained break-glass repository moves from `ae6cf1ee9`'s independent `0/3` baseline to `3/3`.

**AC-2 - Missing, stale, reordered, or cross-attributed durable steps remain red per task.**
Verified by: the ordinary controls remain; atomic controls independently prepopulate started, prepopulate the current-stage report, add report without started, add started without a new report, touch another slug, and touch a prefix-collision subtree. Each rejects only the affected slug; state-only and generated-dispatch evidence have no positive path.

**AC-3 - The observer surface is smaller and provider-independent.**
Verified by: `git diff --numstat` meets the revised deletion/net-negative floors; focused tests accept identical durable journeys with empty/arbitrary transcript and final-message bytes, and no keep-moving completion code reads JSONL, shell/JavaScript text, provider event types, model narration, or a generated dispatch artifact.

**AC-4 - Repository and live confirmation gates are green.**
Verified by: focused dual-episode durable tests, retained smallest-mechanism consumers, `gofmt -l` empty for changed Go files, `go test ./...`, `go test ./... -race`, then exactly one compatible exact-head Codex keep-moving run reports `3/3`. A failure still retains/prints the Git root; success removes it.

## Test plan

Cheapest first:

1. Extend the real-Git table with an atomic first-worker/terminal-archive positive. Before implementation, require red controls for preexisting started/report, report-only, started-only/stale-report, foreign path, and slug-prefix path; keep every ordinary negative. Estimated cost: ~12 seconds, fixture-only.
2. Run the dual-episode oracle directly against the retained clean live repository; this is a spike/diagnostic check only, never a machine dependency in the committed suite. Require `3/3` and the corrected nonterminal `questioned` hold.
3. Run focused durable journey, smallest-mechanism, and live-runner helper/compilation tests. Feed empty/arbitrary transcript and final-message bytes to the existing provider-independence controls.
4. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; verify the revised direct cumulative numstat bounds.
5. After all deterministic proof passes, run exactly one compatible exact-head Codex `keep-moving-posture` journey. Require `3/3`, corrected nonterminal `questioned`, clean workflow, and success cleanup; if red, inspect the printed retained Git root before any further action.

No model transcript, JSONL dialect, final narration, shell/JavaScript grammar, provider event, correlated output, wrapper log, observer schema, receipt, controller, retry loop, daemon, or compatibility layer participates in grading.

## Stage Report: ideation

- DONE: Replace the seed’s transcript-dialect replay/parser direction with a behavior-first design that observes actual dispatch, completion, terminalization, and durable per-task state; deletion must be the primary mechanism metric.
  The design uses existing per-entity state plus path-scoped Git ancestry and deletes 1,394 parser/replay lines, with a net-negative floor of 700 lines.
- DONE: Spike the cheapest real mechanism that can distinguish completed keep-moving behavior from false-red without parsing model prose, shell commands, JavaScript wrappers, or provider event dialects; record the falsifying result before finalizing the design.
  Focused durable helper tests passed; the skipped-transition control falsified final-state-only grading while ordered Git history rejected it.
- DONE: Specify value ACs, semantic boundary, exact files/LOC and tolerance, negative controls, and one live confirmation only after deterministic behavior proof; forbid new observer schema, retry controller, transcript grammar, and unrelated runtime changes.
  ACs, a 10-file/+1 tolerance, insertion/deletion floors, per-task controls, semantic exclusions, and the single post-deterministic Codex run are recorded above.

### Summary

Ideation now removes the rejected transcript observer instead of extending its dialect grammar. The replacement proves each completed task from existing dispatch-entry, worker-report, terminal, archive, and Git-history facts, with a durable falsifier and strict net-negative implementation budget.

## Stage Report: implementation

- DONE: Replace transcript/replay observers with the bounded, net-negative per-task Git-history journey oracle.
  Code commit `1dd01b6db` deletes the four observer files and lands 273 insertions/1,447 deletions across exactly 10 files; restoring a deleted parser consumer or exceeding the declared surface makes the budget/provenance review fail.
- DONE: Prove three independent journeys and exact missing, stale, reordered, and cross-attributed negative controls.
  `TestDurableTaskJourneys` reports 3/3 for ordered real-Git journeys and exactly 2/3 for each affected-slug control; restoring a missing step, accepting a pre-dispatch report, or crediting a multi-path report commit makes the named control fail.
- FAILED: Run focused, format, full, race, then exactly one exact-head Codex live confirmation after deterministic proof is green.
  Focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed, but the authorized corrected exact-head live run completed all three journeys and the oracle still reported 0/3.

### Summary

Candidate `1dd01b6db` is frozen with deterministic and race proof green. The corrected live preflight built `0.27.0-pre2+dev` from that exact tip and `doctor --host codex --plugin-manifest .codex-plugin/plugin.json` confirmed compatibility; Codex then reported `approved-gate`, `ready-one`, and `ready-two` dispatched, worker-reported, passed, and archived, while `questioned` was re-shaped and held at its corrected review gate.

The oracle rejected every completed slug with the exact reason `missing path-scoped dispatch entry with stage and started`. Go test cleanup removed the temporary workflow before its path history could be inspected, so weakening the dispatch rule would be a guess; route back to ideation to redesign the runtime-neutral durable boundary and retain failing state, with transcript/provider parsing still forbidden.

## Stage Report: ideation (cycle 2)

- DONE: Retain and inspect the exact real-Git history shape from a supported completed Codex keep-moving journey before test cleanup.
  Exact candidate `1dd01b6db` with compatible `0.27.0-pre2+dev` retained `/tmp/spacedock-km-retained.SHtfgD`; the unmodified oracle reported `2/3`, and the recorded commit/blob sequence proves all dispatch, worker-report, terminal, and archive steps.
- DONE: Redesign the smallest runtime-neutral durable oracle that credits the real three-task journey and keeps missing, reordered, stale, and cross-attributed controls red.
  A harness-only two-scope spike credited the retained journey `3/3` while all original controls plus a foreign-path archive control stayed red; dispatch/report/terminal remain entity-file-scoped and only archive/held commits admit same-slug support paths.
- DONE: Revise ACs, semantic/file/LOC boundary, deletion inventory, and cheapest-first test plan without transcript/provider dialect parsing or guessed compatibility.
  The body records the retained hashes, revised 10-file +360/-1,400/net-negative-1,000 boundary, unchanged deletion inventory, concrete documentation wording, adversarial controls, and one post-deterministic exact-head live confirmation.

### Summary

The retained journey replaced the destroyed run's unproven `0/3 missing dispatch` diagnosis with a concrete `2/3` archive-attribution defect. The reset preserves strict proof of actual worker execution and the net-negative parser deletion, while narrowly recognizing existing gate sidecars owned by the same entity and retaining future failing Git roots for diagnosis.

## Stage Report: implementation (cycle 2)

- DONE: Apply component-safe same-task support-path ownership only at archive/held boundaries and retain failing live Git roots.
  Code commit `ae6cf1ee9` makes only archive/held commits entity-owned, keeps stricter worker-step scope, and retained the failed clean repo at `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/spacedock-keep-moving-1137023568`.
- DONE: Keep dispatch/report/terminal entity-file-only and prove missing, reordered, stale, foreign-path, and slug-prefix controls remain red.
  Focused real-Git tests passed the sidecar-bearing 3/3 fixture while missing dispatch/report/terminal/archive, pre-dispatch report, foreign report/archive, and `ready-one-other` prefix mutations each rejected the affected task.
- FAILED: Stay within cumulative 10-file +360/-1400/net-negative-1000 bounds; run focused/full/race then one compatible exact-head 3/3 live confirmation.
  Surface is 10 files, +331/-1,449, net -1,118; format, focused, `go test ./...`, and `go test ./... -race` passed, but the single compatible exact-head live run returned 0/3.

### Summary

The retention correction exposed a second supported Git shape rather than losing it: each task went `init → entity-file-only first-worker commit adding started plus the new Stage Report → entity-owned terminal archive`, with no separate `dispatch:` commit. Captain disposition is Material, task-owned, Needs decision and selects a return to ideation: support either the ordinary separate dispatch/report episode or this atomic first-worker episode, while keeping stale, reordered, report-only, foreign-path, and transcript/provider evidence red.

## Stage Report: ideation (cycle 3)

- DONE: Formalize ordinary two-commit and break-glass atomic first-worker durable episodes from the retained exact live history.
  The design records both worker-proof forms, their parent/child blob requirements, separate or combined terminal/archive boundaries, and the exact retained histories `053a50e→9d710f1→363cbe0`, `053a50e→0c7c6b8→8fa4c7b`, and `053a50e→b4d7fdf→08924c9`.
- DONE: Spike controls that accept both real worker proofs while rejecting preexisting, report-only, stale, reordered, foreign-path, and slug-prefix attribution.
  A throwaway `ae6cf1ee9` spike passed 16 deterministic cases plus retained-live `3/3` in 11.229s; parent prepopulation, report without started, started without a new report, pre-dispatch report, foreign worker path, and prefix-collision path all stayed red.
- DONE: Revise ACs, incremental/cumulative surface, and one-final-live test plan without fixture normalization or transcript/provider evidence.
  ACs now measure ordinary and atomic `3/3`; cycle-3 is limited to two existing files/+65 with cumulative 10 files/+400/-1,449/net-negative-1,049, and the plan keeps exactly one post-deterministic live run.

### Summary

Cycle 3 treats the retained break-glass Git history as authoritative without weakening worker proof: only a file-only child that newly adds both `started` and the current-stage Stage Report can replace the ordinary dispatch-plus-report pair. Frozen candidate `ae6cf1ee9` remains untouched; fixture normalization and every transcript/provider-derived alternative remain excluded.

### Feedback Cycles

- Cycle 1: REJECTED — validation/Roborev; surface 10 files/+463/-1,449 vs accepted cycle-4 reset 10 files/+463/-1,449 (100%); AC unchanged. Four Material task-owned false greens accept serialized independent work, a replaced preexisting room binding, a nonterminal canonical archive, and unrelated questioned-ticket mutation. Route the exact retained counterexamples to implementation without restoring transcript or provider observers.

## Review-finding disposition (implementation cycle 4)

### Approved-gate dispatch room binding

- Reviewer observation: the single exact-head `d6a7f462d` live run completed all three task journeys but the durable oracle reported `2/3`, rejecting only `approved-gate` as missing a path-scoped dispatch.
- Released user and normal workflow: the prescribed Codex keep-moving path consumes a captain-approved gate, advances that entity, dispatches its worker, verifies the worker report, terminalizes it, and archives it.
- Observable harm: the supported `approved-gate` journey false-reds even though its own dispatch, later worker report, terminal fields, and canonical archive are durable.
- Authority: `value-ac[AC-1]` requires every completed independent task to be credited from its own durable journey.
- Trigger evidence: retained root `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/spacedock-keep-moving-1344046602` is clean; dispatch `3b699e8` changed `approved-gate.md` plus two newly bound files below `approved-gate/review/review/briefing-1`, the new entity revision's `room-ref`; file-only commit `14b8ddb` later added the implementation Stage Report, followed by terminal `a7addc3` and entity-owned archive `bbf19b0`.
- Worker proposal: Material and Needs decision because accepting any non-file-only worker-proof commit would exceed the cycle-3 design boundary.
- First Officer classification: Material and Needs decision; hold candidate mutation and rerun until a design ruling.
- Captain ruling: preserve the ordinary two-step worker proof, but permit its First Officer dispatch commit to change the ticket file plus only gate-room files newly bound by that same ticket revision, with `room-ref` pointing to their exact subtree. The later worker report remains strictly ticket-file-only.
- Resolution: amend the implementation design narrowly; do not generalize dispatch scope to arbitrary ticket-owned paths or generated artifacts. Atomic first-worker proof, terminal/archive ownership, exact-slug boundaries, and every transcript/provider exclusion remain unchanged.

### Cycle 4 design amendment

An ordinary dispatch is attributable when it changes only the entity file, or when the same entity revision newly introduces a normalized relative `room-ref` below the exact slug and every additional changed path is newly added below that exact room subtree. The later current-stage report must still be a separate entity-file-only commit. Negative controls must reject a preexisting room binding, a room outside the exact slug, a sibling path outside the room, a modified rather than newly added room file, and a slug-prefix collision.

## Stage Report: implementation (cycle 3)

- DONE: Implement ordinary and atomic first-worker episode recognition in the existing durable oracle only.
  Commit `d6a7f462d` added the atomic parent/child proof and combined terminal/archive boundary; removing either worker episode makes its deterministic positive lose one independently credited task.
- DONE: Add parent-prepopulation, report-only, started-only, stale, foreign-path, and prefix-collision controls while preserving prior archive-sidecar controls.
  The 15-case cycle-3 table rejects only `ready-one` for each named mutation while ordinary `3/3`, atomic `3/3`, and exact-slug archive sidecars stay green.
- DONE: Apply the authorized cycle-4 gate-room dispatch amendment without widening worker attribution.
  Commit `970dec18d` accepts `3b699e8→14b8ddb` only when the dispatch revision newly binds added files strictly below its exact-slug `room-ref`; prebound, other-slug, sibling, modified-file, and prefix controls each remain red.
- DONE: Stay within the authorized two-file cycle-4 reset and preserve the provider-independent semantic boundary.
  Incremental cycle 4 is exactly two existing files at `+90/-8`; cumulative direct surface is exactly 10 files at `+463/-1,449` (net `-986`), with no fixture, parser, observer, transcript, provider-event, command, storage, authority, or runtime change.
- DONE: Run focused/full/race/format checks and exactly one fresh exact-head live `3/3` confirmation.
  Both retained roots grade `3/3`; 21 durable controls, smallest-mechanism consumers, live helper/parity guards, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` pass, followed by one `970dec18d` Codex keep-moving run passing in 484 seconds.

### Summary

The durable oracle now recognizes ordinary dispatch/report, atomic first-worker, separate terminal/archive, combined terminal archive, and the narrowly ruled same-revision gate-room dispatch shape while preserving exact worker-report attribution. Exact-head `970dec18d` is clean, provider-independent, deterministic/full/race green, and live-confirmed `3/3`.

## Stage Report: validation

- DONE: Reproduce the 21 durable controls and exact-head Codex keep-moving 3/3 proof, including the newly bound same-slug gate-room dispatch and strict later ticket-file-only worker report.
  Exact `970dec18d`: all 21 controls passed in 15.32s; clean retained roots `1137023568` and `1344046602` each graded `3/3`; one fresh Codex run passed `3/3` in 409.88s and removed its success root.
- DONE: Verify every acceptance criterion and the cumulative 10-file +463/-1449 provider-independent surface against the captain-approved cycle-4 design reset.
  Direct `origin/main..HEAD` numstat is exactly 10 files, `+463/-1,449`; `gofmt -l` was empty, `git diff --check`, focused consumers, `go test ./...`, and `go test ./... -race` passed.
- DONE: Perform the required semantic adversarial pass and classified Roborev review; route every finding through materiality and task ownership before recommending PASSED or REJECTED.
  Detached counterexamples and Roborev `branch_final` produced four findings; the First Officer classified every one Material, task-owned, and `REJECTED / route to implementation`, with no candidate mutation or rerun.
- FAILED: AC-1 - Every completed independent task is credited from its own durable journey.
  Ordinary/atomic and retained/live `3/3` pass, but the serialized positive also passes after each slug reaches archive before the next dispatch, losing the scenario's parallel-dispatch value and invalidating the AC proof.
- FAILED: AC-2 - Missing, stale, reordered, or cross-attributed durable steps remain red per task.
  Named controls pass, but a changed preexisting room binding and a final archive that clears terminal fields each remain green `3/3`, so adjacent missing/attribution variants are not red.
- DONE: AC-3 - The observer surface is smaller and provider-independent.
  The exact 10-file `+463/-1,449` surface deletes the four provider observers; keep-moving grading consumes only Git paths, blobs, ancestry, and workflow root, and arbitrary/empty provider bytes are outside its call boundary.
- FAILED: AC-4 - Repository and live confirmation gates are green.
  Format/full/race and fresh live `3/3` are green, but the live assertion accepts an unrelated edit as the required questioned re-shape, so the gate cannot validly establish the current intended behavior.
- FAILED: Preexisting-room rebind false green.
  Released workflow: captain-approved gate-room dispatch; harm: rebind is credited as a newly bound room; authority: `value-ac[AC-2]` requires preexisting attribution variants red; trigger: old room A replaced by room B plus newly added B file graded `3/3`.
  Proposal/classification: Material; ownership: task-owned; disposition: `REJECTED / route to implementation`.
- FAILED: Serialized independent journeys false green (Roborev).
  Released workflow: keep-moving parallel dispatch; harm: fully serialized motion passes; authority: `value-ac[AC-1]` requires independent durable journeys; trigger: the existing positive's slug loop dispatches through archive one task at a time and passes.
  Proposal/classification: Material; ownership: task-owned; disposition: `REJECTED / route to implementation`.
- FAILED: Nonterminal canonical archive false green (Roborev).
  Released workflow: durable terminal/archive proof; harm: persisted archive may be nonterminal but credited; authority: `value-ac[AC-2]` requires missing terminal evidence red; trigger: prior terminal commit followed by archive clearing status/completed/verdict graded `3/3`.
  Proposal/classification: Material; ownership: task-owned; disposition: `REJECTED / route to implementation`.
- FAILED: Unrelated questioned mutation proves re-shape (Roborev).
  Released workflow: corrected entity is re-shaped and held while independent work moves; harm: silent non-correction passes; authority: `value-ac[AC-4]` requires corrected nonterminal questioned state; trigger: an unrelated-note-only child made `assertDurableKeepMoving` return success.
  Proposal/classification: Material; ownership: task-owned; disposition: `REJECTED / route to implementation`.

### Summary

Exact-head deterministic, full, race, retained-root, and one fresh Codex live run are green, and the cumulative deletion-led surface is exact. Validation nevertheless recommends **REJECTED** because four reproducible Material false greens show that those gates do not prove the promised parallel motion, newly bound room attribution, terminal archive state, or questioned re-shape.

## Review-finding disposition (validation cycle 1 worker triage)

- Serialized independent journeys: confirmed Material and task-owned. The released keep-moving workflow promises independent ready-ticket motion; a dispatch/report/terminal/archive sequence completed for one slug before the next engagement passes `970dec18d`, violating `value-ac[AC-1]`. Proposed fix: require every expected ticket's recognized engagement commit to precede the first recognized terminal commit.
- Replaced preexisting room binding: confirmed Material and task-owned. In the supported approved-gate dispatch, replacing parent room A with room B plus a newly added B file passes, violating `value-ac[AC-2]`. Proposed fix: the narrow room exception requires the immediate parent to have no `room-ref`, not merely a different one.
- Nonterminal canonical archive: confirmed Material and task-owned. A prior valid terminal commit followed by a canonical archive clearing `status`, `completed`, and `verdict` passes, violating `value-ac[AC-2]`. Proposed fix: require those terminal fields in the final canonical archive blob as well as in the recognized terminal boundary.
- Unrelated questioned mutation: confirmed Material and task-owned. An unrelated-note-only child passes as the correction, violating `value-ac[AC-4]`. Proposed fix: require an entity-file-only commit that adds a durable Stage Report and a history transition away from the initial questioned stage, while the final entity remains active, nonterminal, and entity-owned.

The validator's retained counterexamples and Feedback Cycle 1 are unchanged. A throwaway exact-`970dec18d` spike at `/tmp/spacedock-km-validation-fixes.XPeJNw/repo` makes only `internal/ensigncycle/shared_keep_moving_durable_test.go` `+111/-21`; all 24 per-ticket controls, parallel/serialized ordering pair, questioned meaningful/unrelated pair, and both retained real roots pass. Candidate `970dec18d` remains clean pending the First Officer's distinct fix authorization.

## Review-finding disposition (implementation batch live evidence)

- Reviewer observation: the single exact-head `ccd7bee28` live run completed the scenario but returned `0/3`; retained clean root `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/spacedock-keep-moving-2756575595` records batch engagement `4b909b3`, three separate file-only reports, batch terminal `d8aae97`, and three individual terminal archives.
- Worker proposal: Material and Needs decision. The released Codex keep-moving workflow completed genuinely parallel independent work, but `value-ac[AC-1]` false-reds because per-ticket scope rejects the scenario-wide batch commits.
- First Officer classification: Material and Needs decision; hold candidate mutation and rerun pending a design ruling.
- Design ruling: a scenario-wide dispatch or terminal commit may represent all expected tickets only when it changes the complete expected-ticket set and no foreign ticket. Dispatch may additionally carry exact newly bound gate-room subtrees for scenario tickets; each worker report remains ticket-file-only, each canonical archive remains individually terminal, and the questioned re-shape remains independently proven.
- Resolution: add no generic lifecycle engine. The deterministic proof must accept the retained full-set batch history and reject adjacent partial-set and foreign-path batches before another live run.
