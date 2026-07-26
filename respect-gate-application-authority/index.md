---
title: Apply a recorded gate only within the First Officer's assigned authority
status: ideation
source: "Durable-decisions real-sprint correction, 2026-07-26: gate record correctly preserved ideation, but a Shaping FO invoked Commander-owned gate consume because the shipped lifecycle treated approval as authority to apply."
started:
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: mnea9vq3pv1rz1x1hdjbvdg9
gates:
    version: 1
    current:
        gate: gate:mnea9vq3pv1rz1x1hdjbvdg9:backlog
    records:
        - id: gate:mnea9vq3pv1rz1x1hdjbvdg9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mnea9vq3pv1rz1x1hdjbvdg9-backlog-1
              briefing:
                id: briefing:mnea9vq3pv1rz1x1hdjbvdg9:backlog:attempt-1:revision-1
                digest: sha256:064cebf9ee6699261c5213d4b8b9ff42350c64e13060bcd07876a995c7a8e8e7
                digest-domain: canonical-bytes
                request-digest: sha256:4fb4b00f0b0d4dd7dfe9be7b26b4dd904d00a9ed3f7361c2e73dd24871a1fcef
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:mnea9vq3pv1rz1x1hdjbvdg9:backlog:1
                briefing: briefing:mnea9vq3pv1rz1x1hdjbvdg9:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T11:42:35.622923Z"
                decision: approve
                reason: 'The sprint''s real role error proves a narrow general contract gap: recording approval must not imply transition authority, while an explicit conn still permits consume and dispatch. Shape the contract and existing live proof without changing gate mechanics.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Recording a binding approval and applying it are separate operations. The shipped First Officer lifecycle must preserve that separation at the role boundary: presenting or recording an approval never enlarges the current session's assigned transition scope. A Shaping First Officer explicitly assigned to hold members at a gate records the exact Resolution and stops with approved-awaiting-advance; a Commander with the conn consumes the same frozen attempt, advances, and dispatches the entered stage.

This is a general authority rule, not a development-workflow stage name. It must compose with ordinary single-FO operation: when the current First Officer's explicit assignment or conn includes applying the transition, approval may proceed directly through consume and successor dispatch.

## Problem and observed replay

The durable-decisions walking-skeleton task already stated the role split in ordinary
language: approve a real review for a later Commander, leave it
`approved-awaiting-advance`, and let a cold Commander apply it. The shipped lifecycle
replay nevertheless produced this path in the state checkout:

1. `0376b646` filed
   `durable-decisions-release-walking-skeleton/index.md` at `backlog`; its body said
   the Shaping First Officer must not call the Commander's `gate consume` verb.
2. `684d6603` bound the backlog Briefing while status remained `backlog`.
3. `3c176a74` recorded one `agent:first-officer` approval. Status still read
   `backlog`, and its application was `advance/pending` to `ideation`.
4. `a4a60051`, 16 seconds later in the same shaping drive, changed status to
   `ideation` and the application to `consumed`.

The recorder and consumer each did exactly what their command contract promises.
The role error is upstream: `skills/fo-gate-lifecycle/SKILL.md` currently says
“For approve, run” `gate consume` immediately after the close commit.
`skills/first-officer/references/first-officer-shared-core.md` repeats that every
approval is consumed before routing. Neither instruction asks whether the active
assignment authorized application. Recording the Resolution therefore becomes an
accidental authority escalation.

## Host-neutral contract

Normative rule:

> Presenting or recording an approval does not expand the active assignment's
> transition scope. Resolve application authority from the active conversation:
> an express reservation or narrowing controls; otherwise an explicit application
> assignment or an unqualified broad grant—including exact
> `you have the conn`—includes consume and successor routing. Without application
> authority, the First Officer commits the Resolution and stops with
> `approved-awaiting-advance`; with it, the First Officer consumes and commits the
> one-use application, then enters ordinary successor routing.

This rule is independent of host, provider, workflow stage, or organizational role.
The Shaping/Commander incident and ideation/implementation transition are examples,
not words the normative contract depends on. A normal single First Officer holding
the unqualified broad grant `you have the conn` still owns record, consume, commit,
and successor dispatch unless the active prompt expressly reserves application.

Decision authority and transition authority remain distinguishable. A limited conn
may authorize the First Officer to render and record a binding decision while the
same prompt explicitly reserves application for another session. The reservation is
what narrows the otherwise broad conn; recording alone neither grants nor revokes
application authority. Conversely, a later First Officer may be assigned to apply an
already-recorded approval without reconstructing or replacing the decision.

## Proposed wording

Replace the unconditional approve route in
`skills/fo-gate-lifecycle/SKILL.md`.

Before:

> **Route fail-closed.** For approve, run:
>
> `gate consume ...`
>
> Consume itself rechecks currency, successor, blockers, and one-use state under
> lock. ... A nonterminal target then enters ordinary reuse-or-fresh dispatch; a
> terminal target enters the existing merge guard/hook and has no successor
> dispatch.

After:

> **Route fail-closed.** Presenting or recording approval never expands the active
> assignment's transition scope. An express reservation controls; otherwise an
> explicit application assignment or unqualified broad grant such as exact
> `you have the conn` authorizes application. Without that authority, stop after the
> close commit as `approved-awaiting-advance`; do not consume, mutate status, or
> dispatch the successor. When authorized, run:
>
> `gate consume ...`
>
> Consume rechecks currency, successor, blockers, and one-use state under lock. ...
> Commit the consumed descendant before routing: a nonterminal target enters ordinary
> reuse-or-fresh dispatch; a terminal target enters the existing merge guard/hook and
> has no successor dispatch.

Replace the routed summary in
`skills/first-officer/references/first-officer-shared-core.md`.

Before:

> It commits the bound package before presentation, every successful close before
> routing, and consumed approval before routing: nonterminal → ordinary dispatch,
> terminal → existing merge ceremony; revise invokes `«feedback.route»`, while
> hold/ineligibility stops.

After:

> It commits the bound package before presentation and every successful close before
> routing. Presenting or recording approval adds no transition scope. An express
> reservation narrows application; otherwise an application assignment or unqualified
> broad conn authorizes it. Stop pending without authority; when authorized, commit
> consume before nonterminal dispatch or terminal merge. Revise invokes
> `«feedback.route»`; hold/ineligibility stops.

Implementation should tighten adjacent sentences rather than raise the existing
6,600-byte gate-lifecycle or 26,754-byte shared-core caps.

The user-facing concept page currently says every authorized approval proceeds to
consume. Apply this concrete documentation diff to
`docs/site/concepts/gates-and-decisions.md`:

```diff
-After an authorized decision, it records and commits the Resolution before every route. Approval then uses `gate consume`, which rechecks eligibility and atomically writes the successor stage and consumed mark; the consumed descendant commit lands before ordinary successor dispatch.
+After an authorized decision, it records and commits the Resolution before every route. Presenting or recording approval adds no transition scope. An express reservation narrows application; otherwise an explicit application assignment or unqualified broad grant such as `you have the conn` authorizes it. Without application authority, approval remains `approved-awaiting-advance`; when authorized, `gate consume` rechecks eligibility and atomically writes the successor stage and consumed mark, and the consumed descendant commit lands before ordinary successor dispatch.
```

No CLI wording, schema, or command behavior changes.

## Proof design

Contractlint remains structural and gains no semantic phrase or order assertion for
this change. Existing reference closure, skill frontmatter, structural
absence/deduplication, and component byte caps continue to run. They prove only that
the shared references remain loadable and bounded. The two cold live sessions own all
claims about how the model resolves authority, stops pending, consumes, and
dispatches.

Use the existing shared `recorded-gate-lifecycle` live scenario and its current
Claude, Codex, and Pi runners. Do not add a scenario, harness, event format, or
provider. Split its existing journey into two cold launches over the same fixture:

Negative record-and-handoff control, exact prompt:

> Use `$spacedock:first-officer` for this whole run. Workflow directory:
> `{workflow-root}`. Engage only `recorded-gate-task`. You have the conn to judge and
> record this retained gate for the next First Officer, but application is expressly
> reserved for that next First Officer. Approve the retained validation package,
> commit the decision, report the approval awaiting advance, and stop without
> consuming, changing status, or dispatching the successor.

Its checkpoint requires exactly one Briefing record and one decision record, a
committed `advance/pending` application, unchanged `status: validation`,
`approved-awaiting-advance` on boot/readiness, and zero consume, `status --set`, or
dispatch events. A mutant that consumes, changes status, or builds a dispatch before
the checkpoint must fail.

Positive application control, exact prompt:

> Use `$spacedock:first-officer` for this whole run. Workflow directory:
> `{workflow-root}`. Engage only `recorded-gate-task`; you have the conn.

The positive phase is gated on sibling `gqs`
`dispatch-entered-stage-after-gate-consume`: its entity must be terminal `PASSED`,
its product commit must be merged into the exact candidate, and its focused owner
suite must pass before this cold launch. Immediately after consume, the landed gqs
projection—not prompt coaching or harness mutation—must expose the entered working
stage as both `current` and `next`; the ordinary First Officer loop then dispatches
it. The prompt must not name consume, commit, successor dispatch, target stage,
marker, status mutation, or command spelling. The harness must not inject a report,
rewrite status, or special-case the fixture to manufacture dispatchability.

The existing final oracle continues to require the ordered
Briefing-record → decision-record → consume trace, exactly one successful successor
dispatch, one durable successor effect, `state: consumed`, and consumed-commit
ancestry before dispatch. Add the phase-specific negative checks: the second launch
must not bind or record a replacement decision, and omission of consume or dispatch
must fail. It also records the post-consume `current=<entered>,next=<entered>` gqs row
before the successful dispatch; a successor-skipping `status --set`, coached prompt,
or injected completion cannot qualify. Phase-specific artifact names/session
directories keep the two launches cold while reusing each current runner.

No new spike is needed: `TestRecordedGateLifecycleAC7ResumeMatrix` already proves with the
real binary that a committed pending approval survives a fresh process and consumes
once, while gqs owns and has already spiked the entered-stage scheduler boundary.
This task waits for terminal, merged gqs rather than copying or bypassing it. The
unproven model choice is the value under test, so the two-prompt live journey is
implementation's first behavioral proof rather than a throwaway mechanism spike.

## Expected surface and estimates

These file and line counts are advisory planning aids, not implementation authority
or automatic gate-reset thresholds. Ordinary refactoring among existing contract,
documentation, and recorded-gate proof files does not require reconfirmation.

| File | Estimated change | Purpose |
|---|---:|---|
| `skills/fo-gate-lifecycle/SKILL.md` | `+8/-8` | Guard approval application by active assignment scope while preserving consume semantics and byte cap. |
| `skills/first-officer/references/first-officer-shared-core.md` | `+4/-4` | Carry the same rule at the completion/gate routing seam within its byte cap. |
| `docs/site/concepts/gates-and-decisions.md` | `+2/-1` | Document pending approval versus application authority. |
| `internal/ensigncycle/shared_scenarios_test.go` | `+2/-2` | Restate the existing scenario's two-session, role-separated intent. |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | `+105/-25` | Exact prompts, pending/gqs checkpoints, combined oracle, and falsifying mutants. |
| `internal/ensigncycle/claude_live_runner_test.go` | `+18/-6` | Run the existing scenario's record and apply phases cold. |
| `internal/ensigncycle/codex_live_runner_test.go` | `+18/-6` | Same two phases through the existing Codex adapter. |
| `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` | `+24/-8` | Same two phases with separate existing Pi session directories. |

Expected total: about `+181/-60` lines across eight existing files, with zero product
Go changes. `gqs` is a terminal prerequisite, not a file this implementation edits.
Semantic reconfirmation is required only if implementation introduces new authority
state, command or schema mechanics, a new harness, a role registry, or a duplicated
scheduler. Those are design changes regardless of line count.

## Acceptance criteria

**AC-1 (VALUE)** A First Officer explicitly assigned to record an approval but leave
application to another role records the binding Result, commits it, reports
`approved-awaiting-advance`, and does not call `gate consume`, mutate status, or
dispatch the successor.

Verified by the negative live prompt's command log, post-session entity bytes, state
commit, and boot readiness. Any consume, status change, dispatch build, or missing
pending approval fails even if the final narration says it stopped.

**AC-2 (VALUE)** A cold First Officer explicitly assigned the application role
discovers that recorded approval, verifies eligibility, consumes it once, commits the
transition, and dispatches the entered working stage without reconstructing or
replacing the decision.

Verified only after gqs is terminal `PASSED`, merged into the candidate, and green on
its focused owner suite, using the exact entity-only positive prompt containing
unqualified `you have the conn`. The post-consume gqs row, positive live trace, and
existing final oracle must show entered-stage dispatch. A second decision record,
absent/duplicate consume, coached command/status workaround, uncommitted consumed
state, or zero/multiple successor effects fails.

**AC-3 (GENERALITY)** The normative contract is phrased only in terms of active
assignment, conn, transition scope, application, and routing. The existing broad-conn
positive control still records, consumes, and drives without an artificial stop.

Verified by the two cold live controls: an express application reservation narrows
the first conn, while exact unqualified `you have the conn` authorizes the second
session. Contractlint contributes no semantic phrase/order evidence. A static wording
pass without both durable live outcomes does not satisfy this criterion.

**AC-4 (PROOF)** Existing deterministic proof establishes structural closure and
phase-oracle discrimination only. The existing shared live journey proves both exact
natural-language assignments: record-and-handoff stops pending, then authorized
application consumes and dispatches. No new harness or compatibility layer exists.

Verified by existing structural contractlint checks, focused deterministic
ensigncycle phase-oracle tests, the live Claude/Codex/Pi scenario, and a diff showing
no new scenario ID, runner abstraction, schema, command, authority state, or
scheduler.

## Test plan

1. Before editing skill text, add deterministic recorded-gate phase-oracle mutants
   and confirm reservation violations, replacement decisions, coached/status
   workarounds, missing gqs projection, and absent/duplicate dispatch qualify as
   failures. Add no semantic contractlint assertion.
2. Apply the byte-conscious shared wording and documentation diff. Run
   `go test ./internal/contractlint ./internal/ensigncycle` to prove structural
   closure/byte caps, existing command replay, resume, and phase-oracle
   discrimination.
3. Stop until gqs is terminal `PASSED`, its product commit is merged into the exact
   candidate, and its focused owner suite is green. Do not weaken the prompt, inject
   a status/report, or special-case the fixture to proceed earlier.
4. Run the existing `recorded-gate-lifecycle` live scenario through its current
   Claude, Codex, and Pi lanes. The two exact prompts, landed gqs projection, and
   resulting durable state—not substring presence or a green command count—are the
   behavioral evidence.
5. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and
   `go test ./... -race`. Formatting must leave the intended diff unchanged.

## Scope

Shape the narrow change to the shared gate lifecycle/First Officer contract and its existing fixtures/live lane. Do not alter gate recorder, eligibility, consume semantics, gate schema, or provider integration.

Also out of scope: a transition-authority field, role registry, dispatch receipt,
new CLI flag, host-specific wording, compatibility path, standing harness, or changes
to presenter/recorder authority. The active assignment and explicit conn remain
model-read contract inputs; this task only prevents a recorded approval from silently
enlarging them. New authority state, command/schema mechanics, a harness, role
registry, or duplicated scheduler requires semantic reconfirmation; advisory file or
line-count drift alone does not.

## Stage Report: ideation

- DONE: Replay the exact record-only Shaping First Officer versus application-authorized Commander journey using the shipped lifecycle and identify the narrow contract seam.
  State commits `0376b646` → `684d6603` → `3c176a74` → `a4a60051` prove record-pending then unintended same-drive consume; the unconditional approve route is the seam.
- DONE: Define a host-neutral authority rule: recording or presenting approval never expands assigned transition scope, while explicit conn or assignment permits consume and successor dispatch.
  The normative rule uses only assignment, conn, transition scope, application, and routing; sprint roles and stage names remain examples.
- DONE: Choose the smallest shared First Officer/gate-lifecycle wording and existing deterministic/live proof changes without altering recorder, eligibility, consume, schema, or provider mechanics.
  Exact before/after wording stays under existing component caps; contractlint stays structural and the existing shared recorded-gate live scenario becomes a two-cold-session proof.
- DONE: Declare the intended files and non-authoritative planning estimates, including exact natural-language positive and negative live controls.
  Nine existing files total about +189/-62 lines with ±80-line tolerance; both prompts and falsifying outcomes are recorded above.
- DONE: Append a complete ideation Stage Report and commit the state path only.
  This report is appended to the split-root entity; the path-scoped commit and push follow verification.

### Summary

Ideation isolates a prose-level authority escalation: the lifecycle correctly records
and consumes, but it consumes every approval without checking the active assignment's
transition scope. The proposed byte-conscious shared rule, concept-doc correction,
structural checks, and two-session reuse of the existing live scenario prove
record-and-handoff versus authorized application without changing any gate mechanics.

## Stage Report: ideation (cycle 2)

- DONE: State precedence—an existing unqualified broad grant such as exact `you have the conn` includes application authority unless the prompt expressly narrows/reserves it; use that exact broad grant in the cold positive phase, scoped only by entity.
  The normative rule now gives express reservation precedence; the negative prompt reserves application, while the positive prompt is only workflow, entity, and exact `you have the conn`.
- DONE: Remove semantic presence/order assertions from contractlint; keep only reference closure/frontmatter/structural absence-dedup/byte caps, and let two cold live sessions own behavior.
  Contractlint gains no semantic assertion or intended file edit; recorded-gate phase oracles and the current live runners own every authority and routing claim.
- DONE: Make terminal gqs prerequisite for the positive consume→dispatch phase; prohibit prompt/status workaround.
  The positive phase waits for terminal `PASSED`, merged, owner-suite-green gqs and requires its post-consume current=current projection with no coached command, injected report/status, or fixture branch.
- DONE: Make counts advisory; semantic reset only for new authority state, command/schema mechanics, harness, role registry, or duplicated scheduler.
  Eight-file +181/-60 estimates remain planning context only; the five listed mechanism changes—not count drift—require reconfirmation.
- DONE: Append corrected ideation report and commit path-scoped.
  Cycle 2 is appended without replacing the prior report; path-scoped commit and push follow structural verification.

### Summary

Cycle 2 makes authority precedence explicit and moves all semantic proof out of
contractlint into the two cold live sessions. The uncoached broad-conn positive phase
cannot run until landed gqs supplies entered-stage dispatch, and implementation size
is advisory unless it crosses one of the named mechanism boundaries.
