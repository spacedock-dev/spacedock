---
title: Make FO dispatch of moving-target conflict owners explicit
status: ideation
source: "Captain follow-up after the 2026-08-04 conflict-owner and shared-credential diagnosis."
started: 2026-08-06T15:54:14Z
completed:
verdict:
score: 0.93
worktree:
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
dispatch builder already contain all identity needed; the missing piece is a
same-stage handoff mode that binds them without inventing ownership state.

## Proposed approach

Add a narrow `spacedock dispatch build --handoff` mode. It reuses the ordinary
dispatch builder and requires a scope-notes file containing the entity, current
stage, PR, registered branch/worktree, old base/head, moved base, exact conflict
paths, and the next owner action. It validates that the entity is still at the
named stage and that its recorded worktree is the registered checkout of the
stage agent's derived branch. The package is opaque Markdown; the helper does
not parse, classify, or resolve conflicts.

The FO selects transport from identity it already holds:

- If the live worker identity matches entity, stage, branch, and worktree and
  `«addressable-worker»` is present, build `--handoff --advance` and send the
  emitted prompt through the existing live follow-up route.
- Otherwise build `--handoff` and forward the emitted ordinary spawn envelope
  through the existing fresh-dispatch route. Handoff mode validates the already
  stamped stage/worktree itself and performs no state publication.

Both modes preserve entity bytes and Git refs. The PR author and commit
credential are never inputs. This is a code-worktree per-entity hold; the FO
may continue its event loop after dispatch. A split-root state publication
conflict still invokes the existing workflow-wide state-safety halt and never
enters `--handoff`.

This mechanism serves AC-1. The simpler alternative was raw `--advance` or
`--stamp` plus scope notes, but those envelopes describe a next-stage advance
or ordinary entry and do not validate a same-stage registered owner checkout.
No new owner resolver is justified: stage metadata plus registered Git state
already determines the fresh recipient, while live identity determines reuse.

## Spike result

The existing seams are sufficient. On 2026-08-06,
`TestStampIdempotentReRunSkipsAlreadyStamped` proved an already-stamped stage
and registered worktree can emit another fresh envelope without a new state
commit; `TestBuildAdvanceGoldens/codex-host` proved the same-stage live-route
pointer envelope; and `TestCodexMultiAgentV2SpawnInputAlwaysIsolatesFreshDispatch`
proved the fresh Codex mapping forces isolated context. All passed. The missing
semantic is only the explicit handoff mode and its owner/worktree validation;
no new stored state or resolver is required, so D8 should remain in the release.

## Acceptance criteria

**AC-1 (VALUE) — Every accepted moving-target handoff produces one dispatch to
the workflow owner of the existing branch/worktree.** A direct fixture starts
from a real aborted Git conflict and a Captain-authored commit, then exercises
both transports. With a matching addressable handle it observes one live
same-stage follow-up; with that handle absent it observes one fresh spawn whose
stage `agent:`, branch, and worktree equal the fixture's registered values.
Changing the stage agent or registered branch makes the assertion fail.

**AC-2 — Handoff is byte-clean outside its dispatch artifact.** The same fixture
snapshots entity bytes (including `status`, `pr`, `mod-block`, and gate record),
state HEAD, code branch, worktree HEAD, and porcelain before both transports and
asserts equality afterward. Any mutation, new state field, rebase, resolution,
commit, or force update makes the test fail.

**AC-3 — Credential identity cannot select the recipient.** A typed build-output
test configures the Captain as Git author and a distinct stage `agent:` owner,
then asserts the live/fresh targets derive from worker identity/stage metadata
and registered Git state. The implementation contains no Git-author lookup,
generic resolver route, parser, tokenizer, or simulated command language.

## Expected surface and semantic boundaries

Expected files are `internal/dispatch/{dispatch.go,build.go,build_handoff_test.go}`,
`skills/first-officer/references/{first-officer-shared-core.md,fo-dispatch-core.md}`,
one shared live-scenario fixture/runner registration, and
`docs/site/reference/command-reference.md`: 6–8 files and about +140/-10 lines.
Tolerance is two extra files and +100 insertions for host-neutral fixture wiring;
crossing either bound requires a new gate review.

Allowed semantic changes are one additive CLI grammar flag (`dispatch build
--handoff`) and runtime behavior that routes an aborted code-worktree conflict
to a matching live owner or fresh same-stage owner. Stored formats, authority,
stage selection, PR policy, Git credentials, split-root halt behavior, and
conflict-resolution behavior may not change.

Concrete documentation diff for `docs/site/reference/command-reference.md`:

- Before: the `spacedock dispatch` row ends after describing `--stamp` failure
  discrimination.
- After: append: "`dispatch build --handoff` builds a byte-clean same-stage
  conflict-owner package: pair it with `--advance` for a matching live handle or
  omit `--advance` for a fresh dispatch against the registered worktree."

## Test plan

Add one direct Go fixture around the existing builder and typed runtime adapter
output, not a prose oracle. It creates a real conflicting rebase, aborts it,
then supplies the observed argv-independent Git values as opaque scope notes.
It invokes `--handoff --advance` and fresh `--handoff`, captures the emitted
typed envelopes, and compares recipient/worktree plus the before/after bytes and
refs. A tagged live Codex or Claude scenario must exercise one emitted handoff
through the actual addressable or spawn tool and leave a worker-written marker
in the existing worktree; the offline test must red if that marker is attributed
to a Captain/Git-author identity.

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
- Implementation must ship the typed handoff mode and direct fixture before
  changing the FO contract; no lifecycle model or interpreter may substitute.
- Validation must run the direct fixture and one actual runtime handoff, and
  must reject if either recipient or registered worktree is only a literal.

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
