---
title: Make FO dispatch of moving-target conflict owners explicit
status: validation
source: "Captain follow-up after the 2026-08-04 conflict-owner and shared-credential diagnosis."
started: 2026-08-06T15:54:14Z
completed:
verdict:
score: 0.93
worktree: .worktrees/spacedock-ensign-codify-conflict-owner-dispatch-handoff
issue:
sprint: durable-decisions
group: fo-contract
milestone: 0.27.0
id: d8qmey415fsb5q9h6q639ngf
gates:
    version: 1
    records:
        - id: gate:d8qmey415fsb5q9h6q639ngf:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:d8qmey415fsb5q9h6q639ngf-backlog-1
              briefing:
                id: briefing:d8qmey415fsb5q9h6q639ngf:backlog:attempt-1:revision-1
                digest: sha256:e94e3a3c61d1947a0a06e071414b7cb175f4b297cf6a87f0f0c09c092d8a98e6
                request-digest: sha256:0ce185eaddb7133c17ad91e086e0d59db9adcf0ca9bc9cb64d9ffb201ef69731
                room-ref: ./codify-conflict-owner-dispatch-handoff/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:d8qmey415fsb5q9h6q639ngf:backlog:1
                briefing: briefing:d8qmey415fsb5q9h6q639ngf:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-06T15:53:42.057665Z"
                decision: approve
                reason: 'Captain sprint conn authorizes gate approval. D8 has a distinct mechanical end value after G3 was cut: dispatch the recorded workflow owner without credential inference, a resolver, a parser, or new state.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:d8qmey415fsb5q9h6q639ngf:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:d8qmey415fsb5q9h6q639ngf-ideation-1
              briefing:
                id: briefing:d8qmey415fsb5q9h6q639ngf:ideation:attempt-1:revision-1
                digest: sha256:6b8aeea5d2d00e05f555a20fe84882387d18c285eed8bdc4a05f531968ee17fa
                request-digest: sha256:5762bfcf622c60cb58f298da590997df8744593937322a2ca2971a5e6aa2d7cb
                room-ref: ./codify-conflict-owner-dispatch-handoff/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:d8qmey415fsb5q9h6q639ngf:ideation:1
                briefing: briefing:d8qmey415fsb5q9h6q639ngf:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-07T04:43:57.873275Z"
                decision: approve
                reason: The cycle-2 design removes new grammar and state, uses existing owner identity and dispatch routes, and requires one real live owner handoff.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:d8qmey415fsb5q9h6q639ngf:validation
          stage: validation
          attempts:
            - id: gate-attempt:d8qmey415fsb5q9h6q639ngf-validation-1
              briefing:
                id: briefing:d8qmey415fsb5q9h6q639ngf:validation:attempt-1:revision-1
                digest: sha256:2b249456f83e2f7298d0e723004403582446274cba01de81a4e9b385ec58237f
                request-digest: sha256:40491a3fd0e3faec24de912bc5bc83e7e2d2834a5ad7f424da5e14ba29ece5f3
                room-ref: ./codify-conflict-owner-dispatch-handoff/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:d8qmey415fsb5q9h6q639ngf:validation:1
                briefing: briefing:d8qmey415fsb5q9h6q639ngf:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-07T05:49:48.980218Z"
                decision: revise
                reason: Validation proved that both live and fresh fixture arms seed the owner tuple instead of deriving it from the initial stamped dispatch. Fix only that provenance defect, then rerun focused, full, race, and live evidence.
            - id: gate-attempt:d8qmey415fsb5q9h6q639ngf-validation-2
              briefing:
                id: briefing:d8qmey415fsb5q9h6q639ngf:validation:attempt-2:revision-1
                digest: sha256:e956833f41e4790b85fc16f9f0ed75da35710cd6f93419a8551cf32bb48d0483
                request-digest: sha256:8f5acbc9380a1d224d433fedd7040ea4d5a37a112f9ded397e88ead7bbe524c6
                room-ref: ./codify-conflict-owner-dispatch-handoff/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:d8qmey415fsb5q9h6q639ngf:validation:2
                briefing: briefing:d8qmey415fsb5q9h6q639ngf:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-08T06:43:36.968806Z"
                decision: approve
                reason: Exact head 7057ad6f7 derives the owner tuple from the initial stamped dispatch, proves live and fresh routing from that tuple, satisfies all three ACs, and adds no grammar, state, parser, resolver, or authority change.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:d8qmey415fsb5q9h6q639ngf-validation-3
              briefing:
                id: briefing:d8qmey415fsb5q9h6q639ngf:validation:attempt-3:revision-1
                digest: sha256:b34a020053b58244301a47e25ddbdc6c57b33a977fa74314985e6027f74dc569
                request-digest: sha256:408425f89de6664fc518f72b950bfafc643b8a8803b723fb29220c75d2218709
                room-ref: ./codify-conflict-owner-dispatch-handoff/review/validation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:d8qmey415fsb5q9h6q639ngf:validation:3
                briefing: briefing:d8qmey415fsb5q9h6q639ngf:validation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-08-08T17:09:01.003807Z"
                decision: revise
                reason: 'Post-YS validation proves the live recipient and cardinality checks are self-referential. Authorize only the D8-owned typed-event correlation fix: exactly one spawn and one follow-up to the stamped handle, with text-only, wrong-handle, and extra-call negatives.'
            - id: gate-attempt:d8qmey415fsb5q9h6q639ngf-validation-4
              briefing:
                id: briefing:d8qmey415fsb5q9h6q639ngf:validation:attempt-4:revision-1
                digest: sha256:b53720064a79f90e21aaed54ca4e22f2a256b58996f58cd0819a9ee00183600c
                request-digest: sha256:0e1f8ab32a90a5df9ccd4d6201439e22d8b008215a1ed8cd410b330e0a5505aa
                room-ref: ./codify-conflict-owner-dispatch-handoff/review/validation/briefing-4
              withdrawal:
                by: agent:first-officer
                at: "2026-08-09T14:47:32.46304Z"
                reason: Latest validation report heading must exactly name the validation stage before gate decision.
mod-block:
pr: "#645"
---

The First Officer contract names safety behavior for an open-PR moving-target
conflict, but it does not make the next owner dispatch mechanical. Git author
data is unusable as worker identity because workers share the Captain's Git
credential. This task is deliberately narrower than G3: it proves only the
owner-handoff transaction after the conflicted rebase has been aborted.

## Problem

`halt.rebase-conflict` says to abort and surface paths, but stops before a
dispatchable package and recipient are produced. The FO can therefore report
"the owner must reconcile" without contacting anyone, or can confuse the
Captain-authored PR with a Captain-owned worktree. The existing stage `agent:`,
entity status, registered branch/worktree, live worker identity, and fresh
dispatch builder already contain the routing identity. The missing piece is an
explicit same-stage owner-handoff contract that reuses those proven values
without inventing ownership state.

## Proposed approach

Use ordinary `spacedock dispatch build` with a scope-notes file containing the
entity, current stage, PR, registered branch/worktree, old base/head, moved
base, exact conflict paths, and the next owner action. The package is opaque
Markdown; the helper does not parse, classify, or resolve conflicts.

No new build invariant is needed for this event. A manual FO-held rebase is
permitted only after the reconcile/runtime ownership check identifies a live
worker tuple. The original `dispatch build --stamp` already required the
canonical entity, matching status/stage, and registered derived branch before
creating that worker. `git rebase --abort` restores the same worktree and branch,
and this handoff remains in the same stage, so its `agent:` and derived branch
cannot change. A cold or unowned worktree remains report-only and never reaches
this owner-handoff path.

The FO selects transport from identity it already holds:

- If the live worker identity matches entity, stage, branch, and
  worktree and `«addressable-worker»` is present, build ordinary `--advance`
  with the scope notes and send its prompt through the existing follow-up route.
- Otherwise run an ordinary fresh build with the same scope notes and forward
  its spawn envelope through the existing fresh-dispatch route.

Both modes preserve entity bytes and Git refs. The PR author and commit
credential are never inputs. This is a code-worktree per-entity hold; the FO
may continue its event loop after dispatch. A split-root state publication
conflict still invokes the existing workflow-wide state-safety halt and never
enters the owner-handoff route.

This contract serves AC-1. No handoff-specific mode or stronger ordinary build
guard is justified: the ownership and original-stamp preconditions already
prove the tuple before the rebase begins; ordinary scope notes carry the new
assignment, stage metadata selects the fresh worker, and live identity selects
reuse. If a supported handoff can be reached without those preconditions, D8
must be cut rather than growing a resolver or a second ownership system.

## Spike result

The existing transport seams are sufficient. On 2026-08-06,
`TestStampIdempotentReRunSkipsAlreadyStamped` proved an already-stamped stage
and registered worktree can emit another fresh envelope without a new state
commit; `TestBuildAdvanceGoldens/codex-host` proved the same-stage live-route
pointer envelope; and `TestCodexMultiAgentV2SpawnInputAlwaysIsolatesFreshDispatch`
proved the fresh Codex mapping forces isolated context. All passed.

The cycle-2 source audit found that ordinary `runBuildFields` alone checks only
a readable entity and existing worktree directory. That is sufficient here
only because the same-stage handoff inherits the stronger original `--stamp`
and owned-rebase proofs; it is not a claim that arbitrary build input is fully
validated. Strengthening all builds would touch at least 11 unrelated fixture
files and solve a broader problem with no demonstrated supported failure. Retain
D8 as contract plus live proof; cut it if the live journey cannot demonstrate
the inherited tuple without literals.

## Acceptance criteria

**AC-1 (VALUE) — A moving-target conflict reaches the workflow owner of the
existing registered checkout.** One real live-host fixture starts from an
owned rebase conflict and a Captain-authored commit with its worker handle still
addressable. It proves the durable owner-handoff outcome: the required marker is
committed on the fixture's pre-existing registered branch/worktree with the
expected Git author, the checkout is clean, and no replacement branch or
worktree exists. Writing no marker, leaving it uncommitted, or using a new
checkout makes the live grade fail.

**AC-2 — Live and fresh routing preserve the proven owner tuple.** The live
fixture records the tuple produced by the initial stamped dispatch, triggers an
owned rebase conflict, aborts it, and proves the entity authority bytes,
registered branch, and registered worktree still equal that tuple while the
committed owner result exists there and no replacement checkout or branch was
created. The existing focused `--advance` and isolated fresh-spawn tests remain
green; a second fixture arm removes the reusable handle and asserts the ordinary
build envelope derives the same stage agent, branch instruction, and worktree.
Changing any durable tuple member makes the corresponding arm fail.

**AC-3 — Handoff is credential-independent and authority-clean.** The live
fixture snapshots `status`, `pr`, `mod-block`, and the pending gate application
before routing and asserts those frontmatter values remain byte-identical after
the worker appends its required report. The Captain remains Git author while a
distinct runtime worker identity receives the work and writes the marker on the
existing code branch. The implementation contains no Git-author lookup, generic
resolver route, parser, tokenizer, simulated command language, new field, state
publication by the FO, or conflict-resolution action.

## Expected surface and semantic boundaries

Expected files are
`skills/first-officer/references/{first-officer-shared-core.md,fo-dispatch-core.md}`
plus one shared live-scenario fixture and its runner registration: 3–5 files and
about +90/-10 lines. Tolerance is two extra live-harness files and +80 insertions;
crossing either bound requires a new gate review.

The only allowed semantic change is FO runtime behavior: after aborting its own
owned code-worktree rebase conflict, it dispatches the recorded same-stage owner
through the existing live-or-fresh route and may continue unrelated entities.
Command grammar, command output, stored formats, authority, stage selection,
branch/worktree derivation, PR policy, Git credentials, split-root halt behavior,
and conflict-resolution behavior may not change. No docs-site diff is required;
the user-visible CLI is unchanged. The specific skill-text diff is:

- `first-officer-shared-core.md` before: manual code rebase conflict says abort,
  surface paths/peer, and stop. After: state-sync exit 3 remains a workflow-wide
  halt; an owned code-worktree conflict aborts, records the tuple/package, and
  enters the dispatch-core owner handoff while unrelated entities may continue.
- `fo-dispatch-core.md` before: reuse/fresh routing is defined only for stage
  advance. After: the same routes carry an opaque same-stage conflict package to
  the matching live worker or a fresh stage-agent worker, never the Git author.

## Test plan

Add one tagged live Codex or Claude scenario. Its setup creates a real conflicting
rebase in a worktree created by an initial stamped dispatch and retains the
initial worker identity. The FO must abort, send ordinary scope notes through
the real addressable route, and the worker writes and commits a unique marker in
that worktree. A second offline arm removes the reusable handle and validates the
ordinary fresh build envelope against the recorded tuple. Grade only durable
frontmatter authority, the registered Git branch/worktree, committed marker,
clean checkout, expected Git author, aborted rebase, and absence of a replacement
branch/worktree. The supported public stream is not evidence of the internal
recipient or call cardinality.

The live scenario is required because G3 proved that literal ownership values
inside a fixture do not establish FO lifecycle behavior. The supported evidence
boundary is the resulting durable checkout state; do not infer internal calls
from prose or private transcripts and do not add a synthetic event protocol,
runtime schema, shell parser, command log interpreter, or mutation matrix. Run
focused dispatch and live-scenario unit tests, the selected live host lane,
`go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and
`git diff --check`.

## Stage-specific test gates

- Ideation must prove the existing live/fresh dispatch seams can reuse the
  recorded stage/worktree identity and must cut D8 if they cannot.
- Implementation must make the real live journey fail before changing the FO
  contract, then pass it with only the two contract edits; no helper, handoff
  mode, lifecycle model, or interpreter may substitute.
- Validation must run the offline fresh-envelope arm and one actual runtime
  handoff, and must reject if the registered checkout outcome is only a literal
  rather than durable Git and filesystem state.

## Stage Report: ideation

- DONE: Prove whether a self-contained mechanical owner handoff can use existing dispatch identity, branch, worktree, and live-or-fresh routing; recommend cut if it requires new state or a resolver.
  Focused spikes passed for reuse of an already-stamped worktree, the Codex-host `--advance` artifact, and isolated fresh-spawn mapping; only an additive handoff mode is missing, so retain D8 with no state or resolver.
- DONE: Define the smallest end-value contract, acceptance evidence, expected files and insertions, tolerance, and allowed semantics; separate code-worktree holds from split-root state halts.
  The design is one byte-clean same-stage dispatch transaction, 6–8 files/+140/-10 with +2 files/+100 tolerance; code holds may dispatch and continue, while split-root conflicts retain the global halt.
- DONE: Exercise the riskiest live-owner or fresh-owner handoff path once without a parser or simulated command language, and provide a concrete implementation and test boundary.
  Direct existing tests exercised both build routes and the Codex fresh mapping; implementation adds typed `--handoff` fixture coverage plus one actual tagged live runtime handoff, with no command interpreter.

### Summary

D8 is viable only after removing G3's unprovable lifecycle claims. The approved
shape is a small explicit builder mode that binds existing stage identity and
registered Git state to the existing live-or-fresh runtime routes, proved at a
typed artifact boundary and once through a real runtime tool call.

## Stage Report: ideation (cycle 2)

- DONE: Remove the proposed dispatch build --handoff grammar because it adds no distinct transport or state transition; use existing fresh build and reuse --advance routes.
  The revised design uses opaque scope notes with ordinary fresh build or `--advance`; command grammar and stored state are unchanged.
- DONE: Decide whether existing dispatch invariants already validate canonical status, registered branch, and worktree; if not, propose only a general invariant that cannot break legitimate dispatches.
  Ordinary build alone lacks all three checks, but the supported handoff inherits them from the original `--stamp` plus the owned-rebase gate and stays in the same stage; no new invariant is warranted, and a cold/unowned checkout remains report-only.
- DONE: Define a smaller contract and one real live fixture with no parser, simulator, new state, or handoff-specific mode; recommend cutting D8 if general invariants cannot prove the end value.
  The smaller contract adds two FO paragraphs, one real live follow-up fixture whose worker marker must land on the existing branch, and an offline fresh-envelope arm; cut D8 if the journey needs literal ownership or can bypass the inherited tuple.

### Summary

Cycle 2 removes the unnecessary handoff grammar and reuses the routes already
shipped. D8 remains viable as two small contract edits plus one actual live
owner dispatch and an offline fresh fallback check; it adds no helper, resolver,
parser, or state.

## Stage Report: implementation

- DONE: Route one owned code-worktree conflict through existing live or fresh dispatch without new grammar or state.
  Commit `9b0d30f13` adds the same-stage FO contract and a live Codex fixture; removing the `followup_task` runtime call or fresh ordinary build makes the fixture fail.
- DONE: Prove the real worker writes on the pre-existing registered branch and worktree while authority bytes stay unchanged.
  `TestLiveCodexOwnedConflictReturnsToRegisteredWorker` passed in 145.51s; it fails on a missing worker marker, wrong registered branch, non-Captain Git author, changed entity bytes, or un-aborted rebase.
- DONE: Stay within the approved surface and add no parser, simulator, resolver, or conflict-resolution action.
  The 3-file `+164/-3` diff changes only the two approved FO references and one standalone live fixture; `git diff --check` and contract lint passed, with no exact overlap with task 824's disclosed `internal/ensigncycle/live_test.go`/registration surface.

### Evidence

- `go test ./internal/contractlint ./internal/ensigncycle` passed; the live test also grades the ordinary fresh envelope's recorded stage, branch, worktree, and opaque conflict notes.
- `TestLiveCodexOwnedConflictReturnsToRegisteredWorker` passed in 145.51s through the selected live Codex host lane.
- `gofmt -w ./cmd ./internal` and `git diff --check` passed.
- `go test ./...` failed only in `internal/gates.TestV1PilotManifestReadsAndValidates`: `manifest path is not present: stat /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/codex-launch-multi-agent-v2.md: no such file or directory` and `manifest path is not present: stat /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/gate-agent-ergonomics.md: no such file or directory`.
- `go test ./... -race` failed in the same two `internal/gates.TestV1PilotManifestReadsAndValidates` subtests with the same two `manifest path is not present` errors above; all D8 packages passed.

### Summary

The FO now aborts an owned code-worktree rebase conflict and mechanically returns an opaque same-stage package to the recorded live worker or the existing fresh-dispatch path, never to a Git author. A real Codex worker wrote and committed the marker on the pre-existing registered branch while all authority bytes remained unchanged, and the offline arm proved the fresh envelope preserves the same tuple.

## Review-finding disposition

- Finding: the live and fresh fixture arms seed the expected owner tuple with literals instead of deriving it from an initial stamped dispatch.
- Released user and normal workflow: an FO handling a supported owned moving-target code-worktree conflict must route reconciliation to the workflow owner of the existing registered checkout.
- Observable harm: the tests pass when the FO never establishes that owner; a wrong, newly invented, or absent owner tuple can therefore ship behind green live evidence.
- Authority: `value-ac[AC-1]` requires the conflict to reach the workflow owner of the existing registered checkout, and AC-2 requires both routes to preserve the tuple produced by the initial stamped dispatch.
- Trigger evidence: `TestLiveCodexOwnedConflictReturnsToRegisteredWorker` manually creates the checkout and tells Codex the literal entity/stage/branch/worktree tuple; `assertConflictOwnerFreshEnvelope` writes the same literals into scope notes and checks that they reappear.
- Proposed classification: Material evidence defect, task-owned. The exact failing boundary is tuple provenance in AC-1 and AC-2; the existing runtime marker, branch, author, authority-byte, and rebase-abort checks remain useful after that boundary is repaired.
- Proposed disposition: fix by producing and recording the tuple through the initial stamped dispatch, then drive both the addressable follow-up and no-handle fresh fallback from that recorded result; rerun focused, full, race, and live evidence at the new exact head.

## Stage Report: validation

- FAILED: Verify the real same-stage owner handoff and ordinary fresh fallback at exact head.
  Exact-head live Codex passed in 188.32s, but both arms use fixture-authored tuple literals; removing stamped-owner derivation does not make them fail, so AC-1 and AC-2 lack valid provenance evidence.
- DONE: Audit that the 3-file surface adds no grammar, state, parser, resolver, or conflict-resolution action.
  Commit `9b0d30f13` changes only two FO reference files and one live test (+164/-3); no command, stored state, parser, resolver, or Git conflict-resolution implementation changed, and `git diff --check` passed.
- DONE: Run focused, full, race, and applicable live evidence; classify every failure and any 824 overlap.
  Focused suites passed; full and race failed only on two pre-existing absent shared-state manifests; the live Codex lane passed; task 824 has no exact file or semantic overlap, only package adjacency in `internal/ensigncycle`.

### Summary

Validation recommends REJECTED at exact head `9b0d30f13`. The contract surface is narrow and authority-clean, and all candidate behavior exercised successfully, but the sole required owner-provenance proof is self-seeded and can pass without an initial stamped dispatch; this Material evidence defect blocks AC-1 and AC-2 while AC-3 remains satisfied.

### Feedback Cycles

- Cycle 1: REJECTED — detached validation / owner tuple provenance; surface 3 files/+164/-3 vs estimate 3–5 files/about +90/-10; AC unchanged
- Cycle 2: REJECTED — post-YS validation / self-referential recipient and cardinality evidence; surface 10 files/+265/-11 vs estimate 3–5 files/about +90/-10; AC unchanged. Captain-directed YS reconciliation explains the registry surface; FO authorizes only the D8-owned typed-event correlation fix.
- Cycle 3: REJECTED — the public Codex stream does not expose spawn/follow-up events; Captain authorizes a design reset that narrows proof to the durable owner-handoff outcome, preserves the post-YS live journey map, and forbids prose matching, private transcript parsing, or a new event protocol.

## Stage Report: implementation (cycle 1)

- DONE: Derive the owner tuple from the initial stamped dispatch and drive both the addressable and fresh fallback routes from that recorded result.
  Commit `7057ad6f7` replaces fixture-authored checkout and identity values with the `--stamp` envelope identity plus its stamped worktree and actual branch; both live reuse and ordinary fresh build consume that tuple.
- DONE: Prove the repaired test fails if stamped-owner derivation is removed, then rerun focused, full, race, and live evidence at the exact head.
  Removing `--stamp` made the live test fail in 0.22s because implementation had no registered worktree; restored exact head `7057ad6f7` passed focused packages, `go test ./...`, `go test ./... -race`, and the real Codex handoff in 235.71s.
- DONE: Keep the repair within D8's approved test boundary with no new grammar, state, parser, resolver, or conflict-resolution action.
  The cycle changes only `internal/ensigncycle/conflict_owner_handoff_live_test.go` (+83/-19); `gofmt -w ./cmd ./internal` and `git diff --check` passed, and production/contract surfaces remain unchanged.

### Summary

The repaired fixture now establishes ownership through the same initial stamped dispatch that creates the registered checkout, then carries that recorded owner tuple through both supported runtime routes. The exact-head live run proved the existing worker receives the follow-up and writes on the stamped checkout, while the fresh arm preserves the same identity without adding product behavior.

## Stage Report: validation (cycle 2)

- DONE: Verify at exact head that the initial stamped dispatch produces the owner tuple used by both the live same-stage handoff and ordinary fresh fallback.
  At `7057ad6f7`, `stampConflictOwner` starts from blank `started`/`worktree`, parses slug/stage/worker from the stamped envelope and branch/worktree from its checkout, and passes that one tuple to both arms; the recorded exact-head 235.71s live pass fails on any mismatched recipient, branch/worktree, or fresh worker name.
- DONE: Audit the final surface for no new grammar, state, parser, resolver, conflict-resolution action, or authority change.
  Relative to `e24e8234a`, only two FO reference files and the standalone live test change (`+228/-3`); cycle 1 changes only that test, entity bytes remain graded unchanged, and no command, stored field, parser, resolver, state publisher, credential lookup, or conflict resolver was added.
- DONE: Run focused, full, race, formatting, diff, adversarial provenance, and applicable live evidence; classify every failure and any 824 overlap.
  Focused packages passed; full/race failed only on the same two absent shared-state manifests; `gofmt -l`, `git diff --check`, and clean-head checks passed; removing stamp provenance is falsified by the blank initial entity and required stamped checkout; task 824 has no exact file/semantic overlap, only `internal/ensigncycle` adjacency.

### Infrastructure

Two exact-head attempts to reproduce a second completed live Codex handoff failed at the host boundary: one was killed after 520.09s of stream silence after a completed turn, and the retained retry hit Go's 10-minute timeout amid repeated model-manager child-process timeouts after deriving the stamped tuple and aborting the rebase but before `followup_task`; neither reached a candidate assertion or contradicts the completed implementation live proof.

### Summary

Validation recommends PASSED at exact head `7057ad6f7`: the prior Material provenance defect is repaired, all three ACs have exact-head evidence, and no Material candidate finding remains. The independent live reruns are classified as non-candidate evidence-infrastructure failures and are preserved honestly; `/tmp/spacedock-d8-validation-live` retains the second transcript, while the completed 235.71s exact-head implementation run remains the runtime handoff/fresh-fallback proof.

## Stage Report: implementation (cycle 2)

- DONE: Reconcile the candidate with current main and the YS live-test inventory.
  The owner branch was rebased without conflict from `7057ad6f7` onto canonical YS/main `d3e70e958`; exact new head `84f9830ec` preserves both FO contract edits and the stamped-owner proof.
- DONE: Register the owner-handoff live scenario in every required inventory and coverage guard.
  `owned-conflict-owner-handoff` is the 17th registry/common entry with fixture `conflict-owner/stamped-checkout`; the source annotation, count guard, live-CI guide, workflow count, scenario specification, and architecture inventory agree, while all three unchanged `^TestLiveCommon` selectors and existing artifact uploads include it.
- DONE: Preserve the self-contained owner tuple contract and prove focused, full, race, and live checks.
  At `84f9830ec`, reconciliation and focused packages passed, `go test ./...` and `go test ./... -race` passed, and canonical `SPACEDOCK_LIVE_RUNTIME=codex TestLiveCommonOwnedConflictOwnerHandoff` passed in 170.78s; Claude Sonnet/Opus and Pi retain exact D8-owned TODO cells.

### Summary

D8 now participates in YS's canonical portable live inventory instead of restoring its retired standalone test entry point. The existing Codex proof remains self-contained and green, current workflow selectors and artifact paths collect it automatically, and no product grammar, state, parser, resolver, authority, or conflict-resolution behavior changed.

## Review-finding disposition

- Finding: the canonical live grader does not observe the recipient or cardinality of the owner handoff; its `spawn_agent` / `followup_task` check is satisfied by literals in the FO contract text read into the stream.
- Released user and normal workflow: an FO handling a supported owned moving-target code-worktree conflict must route exactly one same-stage follow-up to the workflow owner of the existing registered checkout.
- Observable harm: a fresh or additional worker can write the marker on the registered branch while the stamped live owner receives no follow-up, and every current live assertion can still pass.
- Authority: `value-ac[AC-1]` requires exactly one same-stage live follow-up to the workflow owner of the existing registered checkout.
- Trigger evidence: exact-head `codex-exec.jsonl` contains all `spawn_agent` / `followup_task` matches only in `item_1`, the command output that reads the runtime contract; emitted collaboration events expose only `wait`, while `assertConflictOwnerHandoff` uses unscoped `strings.Contains` and grades no target handle or call count.
- Proposed classification: Material evidence defect, task-owned; AC-1 and the live half of AC-2 remain unproved, although the stamped tuple, checkout marker, fresh envelope, authority bytes, and abort behavior are independently valid.
- Proposed disposition: design reset for the live observation boundary, because valid repair requires runtime/harness evidence that identifies and counts the actual addressed worker call rather than another source-text assertion.

## Stage Report: validation (cycle 3)

- FAILED: Verify the exact reconciled head preserves D8 owner-handoff behavior and all acceptance criteria.
  Exact head `84f9830ec` passed the Codex journey in 149.26s, but AC-1 and live AC-2 fail provenance: the only tool-name matches come from read contract text, so a wrong or additional worker would not fail the grade.
- DONE: Verify the YS registry, common runner, selectors, counts, and artifact collection include D8 consistently.
  `TestRuntimeLiveRegistryReconciliation` and `TestRuntimeLiveCommonSuiteTimeouts` passed: registry/source are 17/17, D8 annotation and fixture agree, all three unchanged `^TestLiveCommon` selectors include it, and exact-head Codex detail/host/metric artifacts were emitted.
- DONE: Run focused, full, race, and the applicable canonical live scenario; report unsupported runtime cells as explicit TODOs.
  Focused suites and formatting/diff/clean-head checks passed; full/race failed only on the two absent shared-state manifests; canonical Codex passed, while Claude Sonnet, Claude Opus, and Pi are explicit D8-owned TODO cells.

### Summary

Validation recommends REJECTED at exact head `84f9830ec`. Reconciliation preserved the stamped tuple and YS inventory, the fresh envelope and registered checkout remain correct, and the supported Codex run completed, but the live recipient/cardinality assertion is self-referential and cannot establish that the stamped owner received exactly one handoff; this Material evidence defect requires a harness observation-boundary reset rather than candidate repair in validation.

## Stage Report: implementation (cycle 3)

- DONE: Replace contract-text matching with exact typed spawn and follow-up event correlation to the stamped owner handle.
  Commit `de9960197` adds a focused `codexCollabItem` correlation assertion: exactly one completed `spawn_agent` must expose one receiver handle and bind its prompt to the stamped worker and dispatch-file pointer; exactly one `followup_task` must target that same handle and carry the recorded entity/stage/branch/worktree tuple plus action marker.
- DONE: Reject text-only, zero, extra, wrong-handle, and mismatched-marker handoff evidence without changing product behavior.
  `TestCodexConflictOwnerHandoffCorrelatesTypedEvents` passes table cases for contract text only, absent or duplicate calls, wrong routing, wrong marker, and a changed tuple. The 2-file `+131/-4` correction changes tests only; it adds no product behavior, command grammar, parser, resolver, state, synthetic language, or broad harness.
- DONE: Reconcile onto the current integration tip and rerun focused, inventory, full, race, and canonical Codex live evidence.
  `origin/main` remains `d3e70e958`, already the candidate base, so 0Y has not landed. At exact head `de9960197`, focused correlation and inventory guards, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` passed. Two canonical Codex attempts failed the new assertion after otherwise completing the owner handoff: the installed public `codex exec --json` stream emitted only typed `wait` collaboration items and zero completed `spawn_agent` items. The retained second run is at `/tmp/d8-typed-live/codex-shared-scenarios/owned-conflict-owner-handoff`; its marker, branch, author, authority-byte, abort, and fresh-envelope checks passed before the typed-event grade failed.

### Summary

The authorized correction removes the self-referential string grade and makes recipient identity and call cardinality explicit at the existing typed collaboration-event boundary. Offline adversarial coverage is green, while the canonical live lane now truthfully rejects the current host stream because it does not publish the required spawn/follow-up events; this is the observation-boundary failure cycle-3 validation identified, not a fallback to contract-text evidence.

## Stage Report: implementation

- DONE: Durable evidence proves reconciliation returns to the stamped owner tuple on the existing branch and worktree, with unchanged authority and no replacement checkout.
  Commits `a87c5534c` and `2dc2953cb` remove the typed-event observer and make the live grade fail on an uncommitted marker, dirty checkout, changed branch/worktree, replacement branch/worktree, wrong author, changed authority bytes, or an active rebase; canonical Codex passed in 183.31s.
- DONE: The post-YS live journey registry, runner annotation, runtime TODO cells, and inventory reconciliation remain accurate.
  Rebased onto post-YS `origin/main` `90b6e0e61`; `TestRuntimeLiveRegistryReconciliation` and `TestRuntimeLiveCommonSuiteTimeouts` pass with the existing `owned-conflict-owner-handoff` entry point, fixture, and Claude/Pi TODO cells unchanged.
- DONE: Focused, full, race, formatting, and canonical Codex live checks report honest results without prose matching, private transcript parsing, or a new event protocol.
  Exact head `2dc2953cb` passes focused dispatch and inventory tests, live-tag compilation, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, `git diff --check`, and the canonical Codex D8 journey; deleting the marker commit or adding a checkout/branch makes its durable grade fail.
Acceptance note: AC-1 and the live half of AC-2 are narrowed to the Captain-approved durable outcome boundary. The entity now names committed owner state, unchanged authority, aborted rebase, original checkout inventory, and fresh fallback tuple as proof, while explicitly excluding internal recipient/cardinality claims from the public stream.

Publication note: GitHub rejected the lease-protected update because the OAuth App lacks `workflow` scope for the existing `.github/workflows/runtime-live-e2e.yml` change; the complete candidate remains committed locally at `2dc2953cb` for validation.

### Summary

The final candidate preserves the FO same-stage live-or-fresh owner handoff and the post-YS 17-journey map while removing the unsupported typed collaboration-event correlation test. The live grade now relies only on durable Git and filesystem outcomes, and the acceptance/test text forbids prose inference, private transcript parsing, synthetic event protocols, and new runtime schemas.

## Stage Report: validation

- DONE: Verify every revised D8 acceptance criterion with independent code, command, or durable-state evidence, including adversarial negatives that can falsify the owner-handoff outcome.
  Exact head `2dc2953cb` passed canonical Codex in 190.50s; a detached exact-head probe accepted the valid outcome and rejected missing/uncommitted marker, dirty checkout, wrong branch/author, changed authority, active rebase, replacement branch, and replacement worktree variants.
- DONE: Verify the post-YS 17-journey registry, runner annotation, runtime TODO ownership, selectors, and artifact inventory remain consistent.
  Focused reconciliation/timeout tests pass at 17/17; `owned-conflict-owner-handoff` retains fixture `conflict-owner/stamped-checkout`, all three canonical selectors collect it, and Claude Sonnet, Claude Opus, and Pi remain D8-owned TODO cells while Codex emitted exact-head detail artifacts.
- DONE: Run the applicable focused, full, race, formatting, and canonical Codex live checks at the exact candidate head; recommend PASSED or REJECTED without using prose or private-event observation.
  Focused packages, live-tag compilation, `go test ./...`, isolated `go test ./... -race`, `gofmt -w ./cmd ./internal`, `git diff --check`, and the canonical Codex journey pass at clean head `2dc2953cb`; recommendation is PASSED solely on code, command, Git, filesystem, and dispatch-envelope evidence.

### Summary

Validation recommends PASSED at exact local candidate `2dc2953cb`. AC-1 is proved by the committed marker on the stamped branch/worktree plus exact inventory and cleanup checks; AC-2 by the stamped tuple carried through the durable live result and ordinary fresh envelope; AC-3 by byte-identical authority, Captain Git attribution, aborted rebase, and the absence of product grammar, state, parser, resolver, credential lookup, or conflict-resolution changes. Internal `spawn_agent`/`followup_task` recipients or cardinality and host prose were not used as acceptance evidence.
