---
id: yv3w8rhxrjyqywadmy666nph
title: Dispatch checklists must not delegate status advancement out of worker scope
status: validation
source: "claude-live smallest-sufficient-mechanism red, PR #762 run attempt 2, diagnosed from artifacts 2026-08-26: the live FO's checklist told both ensigns to 'advance status to done' — a frontmatter mutation the ensign contract forbids and documents no syntax for"
started: 2026-08-26T20:04:25Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-dispatch-checklist-out-of-scope-delegation
issue:
pr: "#779"
mod-block: merge:pr-merge
gates:
    version: 1
    records:
        - id: gate:yv3w8rhxrjyqywadmy666nph:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:yv3w8rhxrjyqywadmy666nph-backlog-1
              briefing:
                id: briefing:yv3w8rhxrjyqywadmy666nph:backlog:attempt-1:revision-1
                digest: sha256:fdeded6ab1ba8db3e1f57aa7a922ad2f54b7f62ea9904566b6bcf27a2b7c3f8f
                request-digest: sha256:a282e23565a3dd60c9d2b1c30a74446d283f395997d4be887147d801441421f0
                room-ref: ./dispatch-checklist-out-of-scope-delegation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:yv3w8rhxrjyqywadmy666nph:backlog:1
                briefing: briefing:yv3w8rhxrjyqywadmy666nph:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T20:04:00.456365Z"
                decision: approve
                reason: The seed defines a clear authority boundary, excludes unrelated repairs, and requires the smallest mechanism plus live proof of tool-mediated, correctly attributed transitions.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:yv3w8rhxrjyqywadmy666nph:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:yv3w8rhxrjyqywadmy666nph-ideation-1
              briefing:
                id: briefing:yv3w8rhxrjyqywadmy666nph:ideation:attempt-1:revision-1
                digest: sha256:a5ff08f926ac6b7664522093e17a167bd5a19011d54fc53bb5fc6503a871ab36
                room-ref: ./dispatch-checklist-out-of-scope-delegation/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:yv3w8rhxrjyqywadmy666nph:ideation:1
                briefing: briefing:yv3w8rhxrjyqywadmy666nph:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T20:34:45.146103Z"
                decision: approve
              application:
                target-stage: implementation
                state: consumed
        - id: gate:yv3w8rhxrjyqywadmy666nph:validation
          stage: validation
          attempts:
            - id: gate-attempt:yv3w8rhxrjyqywadmy666nph-validation-1
              briefing:
                id: briefing:yv3w8rhxrjyqywadmy666nph:validation:attempt-1:revision-1
                digest: sha256:61e8e11595f0dc249f2434f7511fb322626c536a645ac5a52dfe4920ffd8c3f6
                room-ref: ./dispatch-checklist-out-of-scope-delegation/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:yv3w8rhxrjyqywadmy666nph:validation:1
                briefing: briefing:yv3w8rhxrjyqywadmy666nph:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T21:57:43.28551Z"
                decision: approve
              application:
                target-stage: done
                state: pending
---

A live FO wrote "advance status to done" into two ensign checklists. One ensign complied by silently hand-editing frontmatter (the write-scope ban is unenforced for workers). The other honored the ban and improvised: five failed `--set` attempts, `strings` on the binary, and an FO-only skill loaded to find the syntax. The durable-journey grader then failed the run on the hand-edited transition. The FO's own state work was clean throughout.

## Problem

The First Officer turns a stage's free-form `Outputs:` and entity acceptance criteria into a worker completion checklist. The current composer rule says to choose outcome signals, but it does not exclude lifecycle work that belongs to the First Officer. In PR #762's second `claude-live` smallest-sufficient-mechanism run, the passive stage sentence “the entity advances to `done`” became “advance status to done” in both worker checklists.

That instruction is impossible to satisfy lawfully under the existing authority split. Ensigns may edit the entity body and append their stage report, but must not edit YAML frontmatter; they receive no supported status-mutation syntax. The First Officer owns frontmatter and must perform each transition with the applicable `spacedock` command after observing the worker completion signal and accepting the report. A bad checklist therefore turns an unambiguous state owner into conflicting instructions and produces host-dependent recovery rather than a reliable lifecycle.

## Proposed approach

Add one host-neutral authority rule immediately after `«dispatch.checklist»`'s existing assembly paragraph in `skills/first-officer/references/fo-dispatch-core.md`. Keep the current paragraph unchanged and append this exact text:

> **Authority boundary.** Checklist items describe only worker-owned outcomes: never assign `status`/frontmatter mutation or stage advancement to a worker. After the completion signal and report gate, the First Officer alone advances entity state through the applicable `spacedock` command; passive stage prose such as “the entity advances” is lifecycle context, not a worker checklist item.

This is guidance at the decision point that created the defect. The dispatch builder should continue transporting the First Officer's structured checklist byte-for-byte; it cannot reliably infer authority from arbitrary natural-language checklist items. Keyword lint would both miss paraphrases and reject legitimate verification items such as “confirm status remains validation.” A typed checklist-owner schema or binary mutation protocol would add grammar, compatibility, and recovery surface for one observed composition error. Changing only the fixture's passive sentence would hide the reproducer while leaving every other commissioned workflow vulnerable. Changing the ensign contract is also insufficient: it already forbids frontmatter writes and can only fail after the bad assignment has shipped.

Guidance is sufficient only if the unchanged live fixture demonstrates the intended behavior. If the candidate live journey still delegates a transition, ideation's mechanism fails and must return for a stronger design; implementation must not patch the fixture, grader, or worker recovery to manufacture a pass.

## Risk evidence

PR #762 run attempt 2 is the exercised negative baseline. The First Officer authored two dispatch artifacts, and both contained the out-of-scope “advance status to done” item before either worker began. Worker 1 made no supported state-command attempt and then hand-edited frontmatter; worker 2 recorded four failed command attempts before eventually reaching the requested state. The common upstream checklist, not either recovery path, is therefore the smallest causal surface.

The incident analysis also established the attribution rule for proof: Claude stream events must be partitioned on `parent_tool_use_id`. Root events belong to the First Officer; descendant events belong to dispatched workers. Treating the combined stream as one actor previously blamed the First Officer for child mutations and cannot prove this task.

No additional spike is needed. Existing behavior already proves that `dispatch build` preserves the structured checklist, the ensign contract forbids frontmatter mutation, and First Officer lifecycle commands own status changes. The only unverified mechanism is whether the added composer rule changes the model's live choice, which is exactly the required candidate journey below.

## Out of scope

The smallest-sufficient-mechanism fixture and grader, including new extraction or attribution code. The per-host h30 and bf repairs, Pi/Codex behavior, retroactive transcript tooling, binary checklist lint/schema, ensign mutation commands, and unrelated host recovery. No command reference or docs-site change is proposed because command grammar and public stored behavior do not change.

## Expected surface and tolerance

Estimate net LOC change: **+1, across 1 file**.

- Production: 1 insertion, 0 deletions in `skills/first-officer/references/fo-dispatch-core.md` (the exact authority paragraph above).
- Proof: 0 insertions, 0 deletions; exercise the existing live fixture and retain its artifacts rather than altering its grader.
- Tolerance: net 0 to +3 LOC, exactly 1 file. Line wrapping or replacing part of the adjacent paragraph may move the net within that range. A second file, a changed fixture/grader, or more than one deleted semantic line requires a design reset.
- Command grammar: unchanged. Stored formats: unchanged. Checklist transport: unchanged.
- Authority: clarified, not transferred—workers own body/report outcomes; the First Officer owns all status/frontmatter transitions through `spacedock`.
- Runtime: a passive or explicit lifecycle phrase is excluded from worker checklists; after the completion signal and accepted report, the First Officer performs the applicable tool-mediated transition. The rule applies through the shared First Officer contract on every host, while this observed Claude defect is proved only in `claude-live`.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - The unchanged smallest-sufficient-mechanism journey delegates zero entity-state transitions to workers, and every durable status transition is performed by the First Officer through `spacedock` after worker completion.**
Verified by: compare the retained PR #762 attempt-2 baseline (2 of 2 worker checklists delegated `done`) with one candidate `claude-live` journey (target 0 of 2). For each entity, retain the dispatch body, completion event, public stream, ordered path history, and final archive; derive every `status` delta from Git, then require a successful root-FO `spacedock` call after the completion signal for each delta and zero descendant worker state commands or frontmatter edits. Partition root and child events by `parent_tool_use_id` before attribution. Removing the new authority rule or again emitting “advance status to done” makes the checklist count nonzero or leaves a status delta without its root-FO command and fails this AC.

**AC-2 (BOUNDARY) - Free-form stage lifecycle prose cannot become a worker obligation, while exact checklist transport, CLI grammar, stored state, the fixture, and the grader remain unchanged.**
Verified by: keep the fixture's passive “the entity advances to `done`” sentence byte-identical, observe that each resulting worker checklist contains only the worker-owned stage-report outcome, and use `git diff --numstat`/`--name-only` to require the approved one-file surface. Repository tests establish non-regression only; a presence grep over the new prose does not satisfy AC-1 or AC-2.

## Test plan

1. Run `go test ./...` and `go test ./... -race` after the one-file skill change. These are low-cost regression checks; they do not prove model behavior. Run `gofmt -w ./cmd ./internal` as the repository completion command, expecting no relevant diff.
2. Run the existing targeted journey with a retained artifact directory: `SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=sonnet SPACEDOCK_LIVE_ARTIFACT_DIR=<dir> go test -tags live -count=1 -timeout 30m -run '^TestLiveCommonSmallestSufficientMechanism$' ./internal/ensigncycle/`. This is the required, minutes-scale model-backed check because the change is to live First Officer choice; no new fixture, CLI test, or grader is needed.
3. Audit the candidate artifacts against AC-1 rather than trusting the current grader alone. For both entities, read the emitted completion checklist, split root First Officer events from descendant worker events using `parent_tool_use_id`, enumerate Git `status` deltas, and pair every delta with the successful root `spacedock` transition after the worker completion/report. Confirm zero worker state commands/frontmatter edits and terminal durable archives.
4. Compare the candidate's 0-of-2 delegation count with the retained PR #762 attempt-2 2-of-2 baseline. Provider errors, missing artifacts, or a run that does not complete both journeys are infrastructure failures, not passes. Any remaining delegated item rejects the guidance mechanism and returns the task to ideation; it is not repaired by changing the grader or a host-specific worker path.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Stage Report: ideation

- DONE: Select the smallest mechanism that prevents dispatch checklists from assigning first-officer-owned status mutations to workers, and state the actor boundary unambiguously.
  Selected one First Officer checklist-composer rule; workers own body/report outcomes, while the First Officer alone performs tool-mediated state transitions after completion.
- DONE: Define the expected net LOC, insertions, deletions, files, tolerance, observable authority/runtime semantics, and why guidance alone or stronger enforcement is sufficient or insufficient.
  Expected surface is +1 net (1 insertion, 0 deletions) in exactly 1 file, tolerance net 0..+3; the approach section rejects fixture-only, ensign-side, binary-lint, and typed-schema alternatives.
- DONE: Specify falsifiable contract and live-journey proof that every status transition is tool-mediated by the first officer and correctly attributed, while excluding grader and unrelated host repairs.
  AC-1 requires 2/2 bad baseline checklists to become 0/2 and pairs every Git status delta with a successful root-FO command after partitioning child events by `parent_tool_use_id`; fixture, grader, and h30/bf repairs stay excluded.

### Summary

Ideation narrows the repair to one host-neutral composer rule at the point where the bad checklist was created. The unchanged passive fixture and retained failing run form the negative control, while one targeted Claude live journey must prove zero worker delegation and correctly attributed First Officer transitions before the guidance is accepted.

## Stage Report: implementation

- DONE: Add only the approved host-neutral authority paragraph at the dispatch-checklist composition boundary; do not change the fixture, grader, CLI, storage, or worker contract.
  Commit `6af283ec0` adds the approved paragraph after `«dispatch.checklist»` assembly in `skills/first-officer/references/fo-dispatch-core.md`; removing it makes lifecycle prose eligible for checklist composition again.
- DONE: Keep the candidate within net 0..+3 lines in exactly one file and run gofmt, go test ./..., and go test ./... -race.
  Diff is 2 insertions/0 deletions in one file; gofmt ran with its unrelated formatter-only delta excluded, and exact `go test ./...` plus `go test ./... -race` pass (the former also reproduced its initial timeout target alone before a clean retry).
- DONE: Run the unchanged targeted Claude live journey and audit retained artifacts for 0-of-2 delegated transitions, zero worker state mutations, and every durable status delta paired with a successful post-completion root-FO spacedock command.
  Both retained runs emit 0-of-2 delegated checklists (4/4 contain only the ready report), descendant events contain no state/frontmatter command, and all four terminal `ready -> done` deltas follow completion signals and successful root `merge guard --verdict passed` calls; artifacts are under `/tmp/spacedock-live-yv3w8rhxrj-candidate-6af283ec0{,-rerun}`.
- FAILED: The unchanged targeted Claude live test exits successfully.
  Run 1 failed because ready-two's worker report was not separately committed; run 2 failed because ready-two lacked the path-scoped dispatch-entry commit. Both completed the two journeys and preserved the mechanism audit above, but these out-of-scope durability findings kept the unchanged grader red.

### Summary

Implementation clarifies the First Officer's checklist-composition authority at the approved one-file, two-line surface and passes both repository regression suites. Two retained Claude journeys eliminate the original 2-of-2 delegation defect and correctly attribute every terminal state transition to a successful post-completion root-FO command, while the unchanged grader remains red on two distinct lifecycle durability failures outside the approved candidate surface.

### Review-finding proposal 1: ready-two terminal transition lacks a preceding durable worker report

- Exact observation/evidence: in run 1, the child edited `ready-two.md` at 21:13:01Z and sent completion to root at 21:13:11Z but issued no Git commit; root read the report at 21:13:14Z and `merge guard --verdict passed` finalized at 21:13:19Z with exit 0, while history jumps from `0c1bed1 dispatch: ready-two entering ready` to `b0be70f archive ready-two (merge guard)`.
- Released user and normal workflow: supported Claude First Officer-to-ensign standing dispatch for a commissioned ready entity, the normal workflow exercised by the unchanged live journey.
- Observable harm: the archive commit absorbs both report and terminal fields, so durable Git history cannot prove that worker evidence existed before the First Officer's terminal state mutation even though the public stream orders completion, report read, and root command correctly.
- Affected authority: `value-ac[AC-1]` — its required ordered path history must establish the report/completion gate before each durable terminal transition.
- Trigger evidence: a worker appends the required report, sends completion, and omits its mandatory pre-signal commit; terminalization then commits the dirty report with the archive transition.
- Proposed release scope: Needs decision; the harm is material to the protected durability proof, but enforcement belongs to worker completion or merge-guard behavior outside this ticket's approved one-paragraph surface.
- Task ownership: outside this checklist-composer ticket; route to the existing `merge-guard-requires-preceding-report` / ensign commit-before-signal work rather than changing this candidate.
- Proposed disposition: route for decision and hold the candidate unchanged; do not rerun the reviewer until the captain selects or explicitly waives the separate durability repair.
- AC effect: full AC-1 is not established for run 1 because durable report ordering is missing; AC-2 is unaffected, and the 0-of-2 delegation, zero child state-mutation, and successful post-completion root-FO transition evidence remain valid.

### Review-finding proposal 2: ready-two dispatch entry is absent from the grader's qualifying path history

- Exact observation/evidence: in run 2, root `status --set ready-two status=ready started` reports `ready -> ready`, sets `started`, and exits 0 at 21:17:22Z; root `dispatch build --stamp` exits 0 at 21:17:25Z; the worker reads those fields and final global history lists `d343851 dispatch: ready-two entering ready`, yet the unchanged grader reports no path-scoped dispatch entry with stage and started.
- Released user and normal workflow: the same supported Claude standing-dispatch path, using the shipped status and stamped dispatch commands before a normal ensign spawn.
- Observable harm: retained archive-following history cannot certify the entity's durable dispatch boundary despite successful commands and a visible global commit, leaving the commissioned journey at 1-of-2 durable proofs.
- Affected authority: `value-ac[AC-1]` — its verifier requires an ordered per-entity Git history that identifies every durable lifecycle boundary before terminalization.
- Trigger evidence: ready-two is stamped and later archived among two highly similar entities; after cleanup, the grader's archive-following path history does not qualify the globally visible dispatch commit. Retained artifacts lack the commit's changed-path listing, so whether the fault is stamp scope, rename/history identity, or grading remains unresolved.
- Proposed release scope: Needs decision; the proof gap is material to the protected durable journey, while diagnosis or repair would cross into the explicitly excluded CLI, storage, fixture, or grader surfaces.
- Task ownership: outside this checklist-composer ticket; captain must assign dispatch-stamp/archive-history ownership after preserving the current artifacts.
- Proposed disposition: route for decision and hold the candidate unchanged; no reviewer rerun or speculative fix until the owning surface is selected.
- AC effect: full AC-1 is not established for run 2 because its path history is not qualifying; AC-2 is unaffected, and the 0-of-2 delegation, zero child state-mutation, and successful post-completion root-FO transition evidence remain valid.

## Stage Report: implementation evidence run

- DONE: Run one additional unchanged targeted Claude live journey exactly as the entity test plan specifies, with a fresh retained artifact directory and the existing candidate commit `6af283ec0`.
  Exact command passed with exit 0 in 168.841s; raw artifacts and audit manifest are retained at `/tmp/spacedock-live-yv3w8rhxrj-candidate-6af283ec0-proof3/claude-shared-scenarios/smallest-sufficient-mechanism/`.
- DONE: Require 0-of-2 delegated lifecycle checklist items and zero descendant worker state/frontmatter mutations.
  Retained `dispatches/ready-{one,two}.md` contain only the report item; `claude-stream.jsonl` partitions children by `parent_tool_use_id` (`toolu_01Mz6G9rCRGCsuMFNcqJyq9q`, `toolu_01VtJ9wy45hEgPbRvqcFEHaZ`) and shows only body edits plus path-scoped report commits.
- DONE: For every durable status delta, prove a successful root-FO `spacedock` command after the attributed completion signal.
  Ready One signals at 08:51:24/28Z precede ROOT merge tool `toolu_01GC7GmFdXr4E6h7kA9WfyF9` at 08:51:48Z; Ready Two signals at 08:52:25/28Z precede ROOT `toolu_019ng6iBGH2U28xWM199vRGL` at 08:52:32Z; both successful guards produce the only `ready -> done` deltas.
- DONE: Prove the worker Stage Report was durably committed before the terminal transition.
  Qualifying history is `9ea2239` dispatch -> `aea1532` report -> `ee6f5b7` archive for Ready One and `540a941` dispatch -> `370d66a` report -> `db893a6` archive for Ready Two; worker results record one-file commits before completion.
- DONE: Prove the dispatch-entry commit is in the entity's qualifying path history with stage and nonempty started.
  The passing unchanged grader qualifies path-scoped `9ea2239` (`ready-one.md`, `status=ready`, `started=2026-08-27T08:51:00Z`) and `540a941` (`ready-two.md`, `status=ready`, `started=2026-08-27T08:52:01Z`); removing path scope, stage, or started makes the grader fail.
- DONE: Preserve command exits, changed-path lists, ordered commits, parent_tool_use_id partition, and artifact paths.
  `ac1-audit.md` beside the raw stream records the successful tool-result flags, entity-owned changed paths, ordered commit chains, actor IDs, exact dispatch files, and terminal artifact locations.

### Summary

The captain-directed unchanged proof run passes and supplies the previously missing AC-1 durability evidence for both entities. Candidate `6af283ec0`, fixture, grader, CLI, storage, and acceptance criteria remain unchanged; the result is 0-of-2 delegated lifecycle items with every report and state transition correctly committed, ordered, and attributed.

## Stage Report: validation

- DONE: Independently validate AC-1 from the retained proof3 artifacts: 0-of-2 delegated lifecycle items, zero child state mutations, report commits before terminalization, qualifying dispatch commits, and every status delta after completion through a successful root-FO command.
  PASSED — an independent JSONL audit found one report-only item in each dispatch, 2 child body edits/2 report commits/0 state mutations, and ordered commit→signal→root-guard events for both entities; changing an actor ID, order, command result, or chain hash fails the audit.
- DONE: Validate AC-2 and the approved surface: exactly one shared-core file, net 0..+3 LOC, unchanged fixture/grader/CLI/storage/worker contract, and passive lifecycle prose excluded from worker obligations.
  PASSED — `git diff --numstat main...HEAD` is `2 0 skills/first-officer/references/fo-dispatch-core.md`; protected surfaces are byte-identical to `main`, while both emitted checklists exclude the unchanged passive `the entity advances to done` lifecycle context.
- DONE: Perform detached adversarial review plus required focused/full/race/gofmt checks, compare actual LOC/files to estimate, and classify every finding with PASSED or REJECTED evidence.
  PASSED — detached artifact/order assertions and focused durable/mechanism tests pass; actual +2/-0 in 1 file is +1 LOC (+100%) versus the +1/1 estimate but remains inside the approved net 0..+3, exactly-1-file tolerance.

### Reviewer findings

- PASSED — Material findings: none; AC-1 and AC-2 have command/artifact evidence and the candidate stays inside the captain-approved surface.
- PASSED — Deferred risks: none; the two earlier durability gaps are absent in proof3, whose unchanged grader exits 0 and whose raw hashes/order agree with the audit manifest.
- PASSED — Polish observation: `gofmt -w ./cmd ./internal` changes only pre-existing `internal/release/runtime_live_evidence_workflow_test.go`, a file byte-identical between `main` and candidate; the candidate changes no Go file.
- PASSED — Infrastructure observation: fresh exact full/race attempts exhausted the temporary volume or hit the 10-minute package timeout under host pressure, not semantic assertions; implementation's exact full/race runs on `6af283ec0` are green, and the independent focused suite passes.

### Summary

Validation recommends PASSED: the retained candidate journey removes the baseline 2-of-2 delegation defect while preserving First Officer-only state authority and every approved compatibility boundary. No material, deferred-risk, or polish finding belongs to the candidate; fresh full-suite evidence is limited only by the recorded host resource failures.
