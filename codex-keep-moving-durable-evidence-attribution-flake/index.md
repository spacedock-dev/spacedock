---
id: 8bnkrtq4rw46xkbez5zrbmmj
title: Codex keep-moving durable-evidence attribution false-red
status: ideation
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
        gate: gate:8bnkrtq4rw46xkbez5zrbmmj:ideation
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
---

The keep-moving live scenario must credit each completed task from its own durable workflow journey, without making a provider transcript dialect part of the product verdict.

## Problem

PR #513's `codex-live` lane failed `TestLiveCodexSharedScenarios/keep-moving-posture` even though `approved-gate`, `ready-one`, and `ready-two` had each completed the real workflow journey: dispatch entry, worker completion, terminalization, and archive. The observer instead tried to reconstruct that journey from Codex JSONL, report-heading cardinality, command text, and merge-loop output. Harmless changes in those representations produced a false red twice after the same exact-head behavior had passed offline, install, Claude/Opus, and Pi validation.

Candidate `1dd01b6db` replaced that observer with ordered Git history, but its first live diagnosis was destroyed by `t.TempDir` cleanup and could not justify changing the oracle. A second exact-tip run retained the completed workflow and narrowed the false red: the unmodified candidate credited `ready-one` and `ready-two` but reported `2/3` because `approved-gate`'s canonical archive commit also moved two gate-review files owned by that entity. Treating every path beyond `<slug>.md` as foreign conflates a same-entity atomic archive with cross-attribution.

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

## Proposed approach

Keep the candidate's host-neutral durable-task journey oracle. For each expected completed slug, follow that entity's Git history from its archive path and require, in ancestry order:

1. An entity-file-scoped `dispatch: {slug} entering {stage}` commit whose entity blob has the expected stage and non-empty `started`.
2. A later entity-file-scoped worker commit whose entity blob contains that stage's durable Stage Report. `started` or a generated dispatch artifact alone earns no worker-execution credit.
3. A later entity-file-scoped terminal blob with non-empty `completed` and `verdict`.
4. A final entity-owned archive commit: every changed destination must be exactly the active/archive entity Markdown path or lie beneath the active/archive `<slug>/` support subtree. The canonical archive Markdown must exist and the active Markdown must not.

Use path-component boundaries, not string-prefix aliases: `ready` does not own `ready-two.md` or `ready-two/...`. Preserve the stricter entity-file predicate for dispatch, report, and terminalization; only archive and the final durable `questioned` hold use the entity-owned namespace predicate. Any other slug or repository path makes the commit cross-attributed.

Grade the three independent tasks as a set of three separate journeys. `questioned` remains a negative guard: it must be durably re-shaped and nonterminal, and none of its commits may satisfy another slug. The live runner supplies only the workflow root and expected slugs/stages; it does not supply JSONL, final narration, commands, or provider events.

Update the deterministic real-Git positive fixture to mirror the retained archive: `approved-gate` carries same-entity gate-review sidecars, while the other two remain single-file archives. Add a foreign-path archive mutation and a slug-prefix collision control. Keep the existing missing dispatch/report/terminal/archive, stale order, and foreign worker-report mutations.

Give the Codex keep-moving live fixture a failure-preserving temp root: delete it on success, but on failure retain and print the ordinary Git repository path. This is diagnosis only; the oracle still grades the in-place repository. Reuse the same durable completion oracle for the commissioned-task completion fallback in the smallest-sufficient-mechanism scenario.

The rejected alternatives are: weakening dispatch-entry requirements (contradicted by retained history), allowing support subtrees at worker-report commits (unneeded), final archive state alone, `status --archived` alone, an instrumented wrapper log, and any transcript/event/command parser.

## Deletion inventory

- Delete all 433 lines of `internal/ensigncycle/shared_keep_moving_test.go`: provider event decoding, command token/regex inference, narration inference, host-neutral motion trace, and transcript correlation tests.
- Delete all 409 lines of `internal/ensigncycle/shared_keep_moving_negative_test.go`: fabricated Claude/Codex streams, replay dialects, and final-message grammar controls.
- Delete all 202 lines of `internal/ensigncycle/codex_dispatch_evidence_test.go`: `item.completed` decoding, dispatch-build result decoding, Stage Report cardinality, status text matching, merge-output matching, and shell-loop variable parsing.
- Delete all 350 lines of `internal/ensigncycle/codex_dispatch_evidence_regression_test.go`: JSONL constructors and observer compatibility cases. Move only the generic `codexCommandOutput` fixture constructor if its unrelated round-recording consumer still needs it.

No compatibility is retained for these internal observer formats.

## Expected surface and semantic boundary

Expected files: delete the four files above; add `internal/ensigncycle/shared_keep_moving_durable_test.go`; adjust `shared_smallest_mechanism_test.go`, `shared_fixtures_test.go`, `codex_live_runner_test.go`, `claude_live_runner_test.go`, and `docs/runtime-live-ci.md`. A tiny shared fixture-helper move is allowed if compilation requires it.

Budget: at most 10 files plus one helper-only file; at most 360 inserted lines; at least 1,400 deleted lines; cumulative diff must remain at least 1,000 lines net negative. Candidate `1dd01b6db` is already 10 files, +273/-1,447; the reset permits at most 87 further inserted lines in those files. Tolerance is +1 helper-only file and +20 insertions only if compilation requires a fixture move; the net-negative floor is not waived.

Observable semantics changed: test-oracle runtime behavior only. Dispatch, worker completion, and terminalization remain entity-file-scoped; canonical archive and corrected held-state commits may also contain existing support paths owned by the same slug. Keep-moving failures retain their ordinary temp Git root for diagnosis. Command grammar, CLI output, stored formats, mutation authority, runtime dispatch behavior, retry policy, and provider adapters do not change. There is no new observer schema, receipt, retry controller, transcript grammar, or runtime normalization.

Concrete documentation diff in `docs/runtime-live-ci.md`: replace “ordered, path-scoped Git journey” with “entity-file-scoped dispatch/report/terminal commits followed by an entity-owned canonical archive; same-slug review sidecars are allowed only at the archive/held boundary and any foreign slug rejects.” Add one sentence that a failed Codex keep-moving run prints and retains its workflow Git root, while a passing run removes it.

## Out of scope

- The permission-narration oracle removal in #514.
- The separate Codex rejection-flow worker-reuse flake.
- Retry policy or the task-15 wait-watchdog replacement.
- Changes to the keep-moving behavioral contract or its requirement to commission each ready entity.
- Runtime host adapters, dispatch tools, workflow frontmatter, and Stage Report format.

## Acceptance criteria

**AC-1 - Every completed independent task is credited from its own durable journey.**
Verified by: the deterministic real-Git fixture reproduces the retained sidecar-bearing `approved-gate` archive plus two single-file archives and reports `3/3`; removing any one journey reports `2/3`. The independent retained live baseline moves from the unmodified candidate's `2/3` to `3/3`.

**AC-2 - Missing, stale, reordered, or cross-attributed durable steps remain red per task.**
Verified by: table-driven real-Git controls independently remove dispatch, report, terminal fields, or archive; place a report before dispatch; cross-attribute a worker report or archive to another slug; and use a slug-prefix collision. Each control names and rejects only the affected slug, while a same-slug gate sidecar archive stays green.

**AC-3 - The observer surface is smaller and provider-independent.**
Verified by: `git diff --numstat` meets the revised deletion/net-negative floors; focused tests accept identical durable journeys with empty/arbitrary transcript and final-message bytes, and no keep-moving completion code reads JSONL, shell/JavaScript text, provider event types, model narration, or a generated dispatch artifact.

**AC-4 - Repository and live confirmation gates are green.**
Verified by: focused durable-journey and retain/delete helper tests, `gofmt -l` empty for changed Go files, `go test ./...`, `go test ./... -race`, then exactly one compatible exact-head Codex keep-moving run reports `3/3`. A forced harness failure retains and prints the Git root; a successful helper control removes it.

## Test plan

Cheapest first:

1. Extend the deterministic real-Git fixture with the retained `approved-gate` sidecar archive. Assert `3/3`; then independently remove dispatch/report/terminal/archive, reorder the report, add a foreign report/archive path, and use `ready` versus `ready-two` to prove component-safe ownership. Estimated cost: seconds, fixture-only.
2. Exercise the failure-retained temp-root helper with success-removes and failure-retains controls. It must preserve native Git state, not generate a receipt or observer format. Estimated cost: milliseconds.
3. Run focused durable journey, smallest-mechanism, and Codex/Claude runner compilation tests. Feed empty/arbitrary transcript and final-message bytes to the existing provider-independence controls.
4. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; verify the declared numstat floors.
5. After all deterministic proof passes, run exactly one compatible exact-head Codex `keep-moving-posture` journey. Require `3/3`, the corrected nonterminal `questioned` hold, and a clean workflow; if red, inspect the printed retained Git root before any cleanup.

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
