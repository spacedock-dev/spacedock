---
id: yv3w8rhxrjyqywadmy666nph
title: Dispatch checklists must not delegate status advancement out of worker scope
status: implementation
source: "claude-live smallest-sufficient-mechanism red, PR #762 run attempt 2, diagnosed from artifacts 2026-08-26: the live FO's checklist told both ensigns to 'advance status to done' — a frontmatter mutation the ensign contract forbids and documents no syntax for"
started: 2026-08-26T20:04:25Z
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
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
