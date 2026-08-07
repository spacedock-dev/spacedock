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
addressable. It observes exactly one same-stage live follow-up and a
worker-written marker on the fixture's pre-existing registered branch/worktree.
Sending the work to the Captain identity, a new checkout, or no worker makes the
live grade fail.

**AC-2 — Live and fresh routing preserve the proven owner tuple.** The live
fixture records the tuple produced by the initial stamped dispatch, triggers an
owned rebase conflict, aborts it, and observes a same-stage live follow-up whose
entity, agent, branch, and worktree equal that tuple. The existing focused
`--advance` and isolated fresh-spawn tests remain green; a second fixture arm
removes the reusable handle and asserts the ordinary build envelope derives the
same stage agent, branch instruction, and worktree. Changing any tuple member
makes the corresponding arm fail.

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
the real addressable route, and the worker writes a unique marker in that
worktree. A second offline arm removes the reusable handle and validates the
ordinary fresh build envelope against the recorded tuple. Grade only durable
frontmatter authority, registered Git branch/worktree, marker branch, and actual
runtime worker attribution; Git author remains the Captain in both arms.

The live scenario is required because G3 proved that literal ownership values
inside a fixture do not establish FO lifecycle behavior. The simplest
alternative—only inspecting the envelope—is insufficient to prove the runtime
call occurred. No shell parser, command log interpreter, or mutation matrix is
permitted. Run focused dispatch and live-scenario unit tests, the selected live
host lane, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`,
and `git diff --check`.

## Stage-specific test gates

- Ideation must prove the existing live/fresh dispatch seams can reuse the
  recorded stage/worktree identity and must cut D8 if they cannot.
- Implementation must make the real live journey fail before changing the FO
  contract, then pass it with only the two contract edits; no helper, handoff
  mode, lifecycle model, or interpreter may substitute.
- Validation must run the offline fresh-envelope arm and one actual runtime
  handoff, and must reject if either recipient or registered worktree is only a
  literal.

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
