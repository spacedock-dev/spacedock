---
id: 18963egcskzxaje6b5vnas3q
title: Simplify the feedback-rejection flow to five steps or fewer
status: implementation
source: "Captain principle 2026-08-16: anything more than 4 or 5 steps is a footgun. Live evidence: both rejection-flow failure mechanisms are step-count casualties - codex drops step 8's bundled tail (completion gap, ErrNoGateRecord); claude's oracle splits across the step-6/step-8 two-call publish (entries=2 vs 4, wrong round id)."
started:
completed:
verdict:
score: "0.85"
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:18963egcskzxaje6b5vnas3q:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:18963egcskzxaje6b5vnas3q-backlog-1
              briefing:
                id: briefing:18963egcskzxaje6b5vnas3q:backlog:attempt-1:revision-1
                digest: sha256:db36beb0817ac7684245ed81ade8d77f4e20fb9a3e8aefb91acfa479f41497e0
                request-digest: sha256:a8a79d493bc4736f5cae81387d4165fabc076cd3e6487a96d03ea5770a99db73
                room-ref: ./simplify-feedback-rejection-flow/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:18963egcskzxaje6b5vnas3q:backlog:1
                briefing: briefing:18963egcskzxaje6b5vnas3q:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T04:00:54.684905Z"
                decision: approve
                reason: 'Captain directive 2026-08-16: dispatch 18'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:18963egcskzxaje6b5vnas3q:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:18963egcskzxaje6b5vnas3q-ideation-1
              briefing:
                id: briefing:18963egcskzxaje6b5vnas3q:ideation:attempt-1:revision-1
                digest: sha256:9441cd69c131ec2835738e9962409051be1d92fefe5582925f46dafb4f31e231
                request-digest: sha256:ed65b1874c4c30fd003d34fb96d7811c0a85d90e09e0acb1c467bebc39b88c57
                room-ref: ./simplify-feedback-rejection-flow/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:18963egcskzxaje6b5vnas3q:ideation:1
                briefing: briefing:18963egcskzxaje6b5vnas3q:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T04:25:47.692481Z"
                decision: approve
                reason: 'Captain 2026-08-16: approved with all three rulings as recommended - the +107-word surface amendment accepted (the completion conditions are the fix), the three-line fixture-prose repair folded in as a declared scope addition, and the grounding item accepted on the FO''s cross-check attestation'
              application:
                target-stage: implementation
                state: consumed
---

Redesign skills/feedback-rejection-flow/SKILL.md from 8 steps to at most 5, with the single-publish shape the diagnostic evidence recommends: accumulate the round, publish ONCE under one round id, and make the gate re-entry a single unbundled step. Every step becomes one action with one completion condition; no step may bundle a tail an FO can drop.

Coordination, binding: hz (repair-codex-rejection-round-recording, in implementation) owns the binary side (inline state commit) and grading honesty; its reconciliation of the two diagnosis accounts determines whether the oracle repairs land there or the oracle re-anchors to this redesign's single-publish shape. Whichever lands second reconciles. The live rejection-flow journey and its oracles are the proof surface for both.

## Problem

Today's flow is eight steps, and the two live failures are **two different bugs**,
not one bug wearing two hats. Conflating them was the seed's error and this design's
own first draft: claude never entered the dirty-tree trap at all. Each mechanism
below is quoted from the run stream that produced it and re-exercised against a
binary built from `origin/main` (`0c6a2c32a`).

**Claude: the round is published before it is complete.** A round is the reviewer's
two entries plus the correction worker's two disposition entries — all four carrying
the same briefing id (`rejectionCompleteLog`, `shared_fixtures_test.go:109`). In CI
run 31922268382 the FO fired step 6's publish while the log still held only the
reviewer's half, then fired step 8's publish after the worker's entries landed:

    round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=2
    round=round:rejection-task:validation:2 stage=validation cycle=2 briefing=briefing:rejection-task:validation:round-1 entries=4

The stream check requires `validation:1` and `entries=4` **on one line**
(`rejectionRoundSuccess`, `shared_round_recording_test.go:17`). The first call has the
right id and the wrong count; the second has the right count and the wrong id. Neither
matches, `claudeRecordedRejectionRound` returns false, and the run fails as
`rejection-round-missing` — a label that is simply untrue, since the second call
recorded the round correctly and completely.

**That claude run committed and re-entered the gate.** Right after the second publish
it ran `git add rejection-task/index.md rejection-task/review/validation/round-2` and
committed. Claude executed all eight steps, tail included. Its only defect was *when*
inside step 6 it published. Steps 4, 5 and 6 spend three steps expressing the single
condition that governs that timing — "the worker's entries are in the log" — and the
publish sits in the third of them, with no step owning the condition.

**Codex: the tail of step 8 is dropped, and dropping it is silent.** Step 8 bundles
four actions: append the Cycle line, publish `validation/2`, invoke
`spacedock:fo-gate-lifecycle`, run `«gate.lifecycle»`. In run 31915540750 attempt 1
the FO did the first two and stopped. Its whole remaining command tail after the
publish was `status --next --json`, `status --read`, `status --where`: it queried the
scheduler, got an empty answer, and reported the journey complete. `gate prepare` was
never invoked, so the entity never gained a `gates:` block, so the codex-only clause
at `claude_live_runner_test.go:400-402` read `ErrNoGateRecord` and reported it under
the same `rejection-round-missing` code as claude's unrelated failure.

The empty answer is unambiguous once the mechanism is exercised:

| state | `status --next --json` | `gate prepare --artifact <entity>` |
|---|---|---|
| after `gate record --round`, uncommitted | `{"dispatchable":[],"ready_gates":[]}` | exit 1, `selected source differs from its committed Git object` |
| committed | row appears once the verdict is non-rejecting | exit 0, `state=open` |

An empty scheduler with no error reads as "the run is finished". Codex's *passing*
runs escape by improvising: after `state commit` they run `git status --short`,
notice nothing was committed, and issue a raw `git add && git commit` before
`gate prepare`. That improvisation appears nowhere in the skill. Passing is currently
a matter of whether the model happens to distrust an exit code.

**The two bugs share no cause, and the single-publish shape alone fixes only one of
them.** Publishing once, at the point the round is complete, removes claude's failure
outright. It does nothing for codex, whose trap is armed by any uncommitted write
before a gate surface read — one publish arms it exactly as well as two. That is why
the commit below is a step rather than a clause, and why this entity is not
self-sufficient: the verb that step names is a no-op on inline workflows until hz's
fix lands.

Two findings from the spike sharpen the picture beyond the two diagnosis accounts:

1. The `gate prepare` refusal is **artifact-scoped**, not tree-scoped. With the
   entity dirty but the artifact a separately committed file, prepare exits 0 with
   `state=open`. It refuses only when the selected artifact is the dirty entity.
   The unconditional half of the trap is the `status --next` drop.
2. A clean tree is necessary but **not sufficient** for the `needs-preparation` row.
   Measured on the same fixture, clean each time: one REJECTED validation report →
   `ready_gates":[]`; the same report with a `### Summary` → `[]`; two reports, both
   REJECTED → `[]`; a PASSED re-review report present → one `needs-preparation` row.
   The **verdict** gates the row, not report completeness or count. So the gate is
   unreachable until the re-review posts a non-rejecting verdict — the gate re-entry
   genuinely cannot precede the re-review, and no ordering choice can avoid that.

## Proposed approach

Five steps. Each is one action with one completion condition, stated as an explicit
`**Done when**` clause, and the ordering is forced by the measurements above rather
than chosen for symmetry:

- The round's log is complete only after the correction worker's entries land, so
  the single publication sits immediately after the correction — not after the
  re-review, which contributes nothing to this round's log.
- Every FO write for the round (the Cycle line, the round pointer) precedes one
  commit, so the commit makes the whole round durable at once.
- The commit is its own unbundled step because its omission is the only failure here
  that is **silent**: an empty scheduler and a success report. The two writes it
  follows are adjacent entity writes with no wait between them, and dropping either
  of those fails loudly through the oracle. The one action that must never be a tail
  is therefore the one action that gets its own step.
- The gate re-entry is last, with nothing after it, because the graded end state is
  exactly "one open gate presented, then stop".

**Rejected: teaching the skill to distrust `state commit`.** The diagnosis proposed,
as a skill-level mitigation, that the step tell the FO to check `git status` after
`state commit` and issue a raw `git add && git commit` when it finds the tree still
dirty — which is precisely what codex's passing runs improvise today. That is
declined, for the reason hz's entity already gives: it hands a mechanical guard back
to model discipline at exactly the point model discipline is proven to fail, and it
would have to be repeated in `fo-gate-lifecycle` and `fo-dispatch-core` as well. It
also contradicts the project's own priority that the binary owns mutation guards.
Step 3 instead states the required end state — the tree is clean for that entity —
and names no escape hatch. An FO cannot satisfy that condition by trusting an exit
code, and once hz's fix lands the contract verb genuinely achieves it. If hz's fix
does not land, this entity's step 3 is unsatisfiable on inline workflows and the
right response is to say so, not to write the raw-git workaround into the contract.

The replacement `## Feedback Rejection Flow` section reads:

> When a feedback stage recommends REJECTED, run these five steps in order. Each is
> one action with one completion condition, and is unfinished until that condition
> holds: never treat a step's first command as its completion.
>
> 1. **Deliver the authorized correction.** Read the rejected stage's `feedback-to`
>    target — the stage that receives the fix request, not the reviewer — and the
>    already-authorized workflow package: rejected snapshot, finding evidence,
>    existing workflow classifications, FO-authorized dispositions, and concrete
>    revise assignment. If the distinct authorization or assignment is missing, hold
>    at the active workflow's review-finding checkpoint; routing is ineligible. Route
>    the package unchanged to the target stage in the same worktree, carrying the
>    concrete assignment, not an acknowledgment request or a new classification
>    request. Reuse the existing handle through `«addressable-worker»` only when it
>    is addressable, reuse conditions pass, and `«context-budget»()` reports it under
>    budget; otherwise shut down and fresh-dispatch. If no probe is declared, proceed
>    to reuse.
>    **Done when** the correction is complete in durable workflow state: the target
>    worker's own entries closing this round's review log where the workflow keeps
>    one, otherwise its `«completion-signal»` attributed by mailbox content, task
>    path, or durable state. The immediate routing response is never completion.
>
> 2. **Record this round, once.** Append the authorized `### Feedback Cycles` line
>    for this round when the active workflow declares that projection, then invoke
>    the neutral recorder exactly once for the whole rejection cycle:
>    `${SPACEDOCK_BIN:-spacedock} gate record --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl`.
>    CYCLE is this rejection cycle's number and the only round id this cycle
>    publishes; there is no second publication after the reviewer re-run. Do not
>    define, normalize, or interpret the Cycle line's category labels, fields,
>    tolerance, estimate, or drift grammar. The recorder retains the canonical
>    two-file room and advances `review-round`, without receiving or interpreting the
>    Cycle line.
>    **Done when** the recorder exits successfully reporting the complete round
>    summary, counting every entry this round accumulated. If it fails, produces no
>    result, or reports an incomplete round, hold the flow. Do not claim that the
>    round was recorded, re-run the reviewer, or prepare the next gate.
>
> 3. **Commit the recorded round.** Invoke `«state.commit»(slug)`.
>    **Done when** the recorded round is durably committed and the tree is clean for
>    that entity. Until it is, the entity is absent from `status --next` and refused
>    by `gate prepare`: an empty scheduler here means this step is unfinished, never
>    that the run is over.
>
> 4. **Re-run the reviewer on the corrected state.** Ask that reviewer to re-review
>    the updated entity state, never to validate its own fix work. Re-run the
>    kept-alive reviewer through the `«addressable-worker»` capability step 1 used,
>    when it remains addressable and reuse conditions pass; fresh-dispatch a reviewer
>    only when it is not. On cycle 3 run no reviewer: escalate to the human instead
>    of a fourth round.
>    **Done when** this cycle has a reviewer verdict — or, on cycle 3, the escalation
>    to the human is recorded and the flow has stopped here.
>
> 5. **Re-enter the gate.** Invoke `Skill(skill="spacedock:fo-gate-lifecycle")` and
>    run `«gate.lifecycle»` for the updated stage. A `needs-preparation` row is work
>    to perform, not a stopping condition.
>    **Done when** exactly one fresh open gate has been prepared and presented, and
>    the flow stops — without resolving or applying it, changing terminal state, or
>    dispatching a successor.

The two trailing paragraphs (the `«write.classify»` ownership sentence and the
`## Workflow-defined correction-round projection` section) are unchanged.

**The flow's one branch is cycle 3**, placed at step 4 — the step it replaces —
rather than split across two steps. Because the publication and commit already
precede it, "publish the rejected round, then escalate" (today's step 6 tail) falls
out of the ordering instead of needing a special case.

**Oracle alignment, in step with the shape, not as prose-greps.**

- `shared_round_recording_test.go`: delete the `validation/1`→`validation/2` fallback
  (:184-189) and pin `validation/1`, so a second publication makes the run fail
  instead of being absorbed. Add an invocation **counter** — the current
  `commandRecordsRejectionRound` pins `--round validation/1`, so it cannot see a
  `validation/2` call at all; the counter must recognize any `gate record --round
  <stage>/<n>` for the entity, require exactly one, and require its round to be
  `validation/1`.
- `claude_live_runner_test.go`: wire that counter in as one more `durableSemantic`
  alongside the existing three. It reads the run stream, so it grades observed FO
  behavior, not skill wording.
- **No contractlint change, and no new structural check.** Two corrections found by
  building and testing against `origin/main` rather than the working checkout, which
  is behind it:
  1. `internal/contractlint/feedback_rejection_publication_smoke_test.go` no longer
     exists. `723028f01` ("Retire banned prose-grep contract pins") already deleted
     both of its pins as tests the Proof policy forbids. Running the contractlint
     package against this draft in an isolated `origin/main` worktree passes with
     `no tests to run` for that name. So nothing in contractlint constrains this
     rewrite, and there is no retirement for this entity to perform.
  2. The structural step / completion-condition count this design first proposed is
     itself banned. The policy's test is whether the expected value comes from
     outside the file under test: counting `**Done when**` does not qualify, because
     a semantics-preserving paraphrase ("Complete when") reds it while a step that
     quietly grows a second action keeps exactly one clause and passes it. It is
     dropped rather than defended.

  The consequence, stated plainly: after this change **no committed test guards the
  publication ordering by reading the skill**, and that is the correct end state
  under the policy — "A contract or skill change is PASSED only when a live drive
  observed the behavior it claims." The guard is the stream counter above plus the
  live journey. The one-off step count is validation evidence for the run, which the
  policy explicitly permits for an existence fact, and is never committed as a test.

**Ruling requested at the gate — one fixture-prose conflict.** The rejection
fixture's README (`shared_fixtures_test.go:148`) says "The first rejected cycle is
projected by the round recorder from `rejection-task/inputs/feedback-cycle.txt`; do
not hand-write Cycle 1. Record `- Cycle 2: PASSED` after the re-validation passes."
The recorder cannot do this: `gate record ... --feedback-cycle <path>` exits 2 with
`unknown gate flag: --feedback-cycle`, pinned by `internal/cli/gate_test.go:250-251`,
and nothing in `internal/gates` reads that file. So Cycle 1 has no writer today, and
the "Cycle 2: PASSED" slot assumes a post-re-review append the new shape does not
have. Nothing in the rejection journey grades either line — `assertRejectionFlow`
and `assertRejectionRecordedRound` never look at them; only the *escalation* journey
grades Cycle lines, and it hand-writes one `- Cycle N: REJECTED` per round, which
the new shape serves exactly. The minimal repair is three lines of fixture prose:
have the FO write one Cycle line per rejection round and drop the unprojectable
instruction. This entity declares the workflow-defined cycle-line grammar out of
scope, so it is surfaced here for an explicit ruling rather than absorbed.

## Coordination with hz (repair-codex-rejection-round-recording)

Boundary proposed to `spacedock-ensign-hz2ankag6f-implementation` and awaiting its
confirmation. hz keeps `internal/cli/state_sync.go`, `state_commit_test.go`, the
codex-only clause and new `rejection-gate-not-prepared` code at
`claude_live_runner_test.go:397-403`, `claude_runtime_helpers_test.go`, and the
`state commit` command-reference diff. This entity keeps
`skills/feedback-rejection-flow/SKILL.md` and the round-recording oracle. The two
touch `claude_live_runner_test.go` in different places — hz changes the codex clause,
this adds a fourth `durableSemantic` — so whichever lands second reconciles that one
file.

**This design depends on hz's fix.** Step 3's completion condition is only
satisfiable through `«state.commit»` once the inline change lands; on today's binary
that verb prints `Inline workflow — entities live beside the README; nothing to
commit to a state checkout.`, exits 0, and leaves the tree dirty. hz therefore lands
first. Its AC-2 requires five consecutive local greens driven through today's
eight-step skill, so this entity absorbs the cost of re-measuring the journey under
the new shape rather than invalidating that evidence — which is what AC-3 below
already buys.

## Out of scope

The binary state-commit fix and grading-honesty repairs (hz owns them). The
workflow-defined cycle-line grammar — with the single fixture-prose exception raised
for a ruling above. Widening the prepared-gate assertion beyond codex.

## Expected surface and tolerance

Four files, roughly 85 insertions and 25 deletions. Tolerance ±30 lines and ±1 file.

| File | Change |
|---|---|
| `skills/feedback-rejection-flow/SKILL.md` | 8 steps → 5 (~+35/-20) |
| `internal/ensigncycle/shared_round_recording_test.go` | drop the round-id fallback, add the invocation counter (~45) |
| `internal/ensigncycle/claude_live_runner_test.go` | wire the fourth `durableSemantic` (~4) |
| `internal/ensigncycle/shared_fixtures_test.go` | fixture Cycle-line prose (~3), only if the gate rules it in scope |

No contractlint file is touched, and none needs to be: the two pins that once read
this skill were deleted upstream in `723028f01`, and the replacement structural check
this design first proposed was dropped as banned. Verified by running the
contractlint package against the new skill text in an isolated `origin/main`
worktree.

**Amendment to the seed estimate, stated plainly.** The seed declared "near zero or
negative on the skill". The measured draft is the opposite: 572 → 679 words, 4158 →
4802 bytes, 28 → 34 lines. The entire increase is the five explicit `**Done when**`
clauses, which are AC-1's deliverable — the eight-step version has no completion
conditions at all, which is why tails were droppable. Cutting them to meet the seed's
number would remove the fix. Step *count* falls 8 → 5; word count rises. The gate is
asked to approve the word-count baseline as +107, not ~0.

**Declared semantic changes.** Command grammar: none. Stored formats: a completed
rejection cycle now retains one round room under one round id; `validation/2` is no
longer produced by this flow, so entities completing a cycle after this change have
one `review/<stage>/round-N/` directory where they previously had two. Authority:
none — no new write target, no new decision rights. Runtime behavior: the FO commits
as an explicit step after recording the round; the FO no longer publishes a round
after the reviewer re-run; the workflow's Cycle line is appended once per rejection
round before publication instead of twice. Live-lane vocabulary: unchanged by this
entity. Contract surface: unchanged — no contractlint pin is added or removed, since
the two that once read this skill were already deleted upstream.

## Risk evidence and spike record

The riskiest unverified mechanism was whether a **single** publication reaches the
end state the oracle demands, since removing the second call could have been what
advances the pointer or prepares the room. Exercised before designing, on a fixture
rebuilt from `writeRejectionWorkflow` against an `origin/main` binary:

- One `gate record --round validation/1` with the complete 4-entry log emits exactly
  the line the oracle's success regex pins, `entries=4`, and creates exactly one
  room with the two canonical files, with the `review-round` pointer's id, briefing,
  and room-ref all agreeing on round 1.
- Driving both of today's calls establishes AC-2's baseline: two rooms
  (`review/validation/round-1/` and `round-2/`) and a pointer reading
  `id: round:rejection-task:validation:2` whose own briefing is still
  `briefing:rejection-task:validation:round-1`. Note this spike passed the same
  complete log to both calls, so its two rooms hold identical bytes; the real claude
  run's did not, because its first call fired against a half-written log. The room
  and invocation counts are the baseline, not the bytes.
- Committing then running `gate prepare --artifact <entity>` exits 0 with
  `state=open` and briefing `briefing:rejection-task:validation:attempt-1:revision-1`
  — the exact id `assertRejectionRoundGateBoundary` requires. The second publication
  contributes nothing to reaching that state.
- The same prepare on the uncommitted entity exits 1 with `selected source differs
  from its committed Git object; commit the exact source before preparation`.
- `status --next --json` returns `{"dispatchable":[],"ready_gates":[]}` while dirty,
  and yields the `needs-preparation` row only once committed **and** carrying a
  non-rejecting verdict (four-way isolation recorded under Problem).

No further spike is needed: the mechanisms this design rests on — single-publish
round semantics, the committed-object prepare gate, the scheduler's readiness
conditions, and the recorder's lack of a `--feedback-cycle` capability — were each
exercised directly, the last against its pinning test at `gate_test.go:250-251`.

## Acceptance criteria

**AC-1 — The shipped flow is five steps, each one action with one completion
condition.**
Verified by: a one-off count of numbered steps and `**Done when**` clauses in the
shipped skill, pasted into the validation report as run evidence — never committed as
a test. The Proof policy permits a grep exactly here and no further: the step count is
an existence fact about the file, which is itself the claim. It does not and cannot
prove the *behavioral* half — that each step is genuinely one action — and no reading
of the file can, since a second action smuggled into a step's prose keeps the count at
five. That half is proven only by AC-2 and AC-3's live drives. This AC asserts a
mechanism and counts only in service of those two.

**AC-2 — One publication per rejection cycle: a completed cycle leaves exactly one
round room under one round id, and the recorder was invoked exactly once.**
Verified by: (a) a scripted fixture drive asserting that after the flow's single
publication the entity has exactly one `review/<stage>/round-N/` directory and a
`review-round` pointer whose id and briefing agree; (b) the live oracle counting
successful `gate record --round` invocations in the run stream and requiring exactly
one, at round `validation/1`. The independent baseline that can move the wrong way is
the room count and the invocation count: both are 2 today, reproduced above, and must
reach 1; an FO that publishes twice fails (b), and a design that publishes three
times moves it further wrong.

**AC-3 — The live rejection-flow journey passes on both hosts under the new flow,
with the oracle's round-id fallback removed.**
Verified by: `assertRejectionRecordedRound` resolving `validation/1` with no
`validation/2` retry, and `TestLiveCommonRejectionFlow` passing on claude and codex —
targeted local runs on each, then the full-matrix stack run. Removing the fallback is
what makes this falsifiable: with it, a double publication still passes; without it,
the run fails. Because hz's inline `state commit` fix is a precondition of step 3,
these runs are measured on top of it.

**AC-4 — The suite stays green, and the escalation journey is not collateral damage.**
Verified by: `go test ./...` plain and `-race`; and `TestLiveCommonFeedbackThreeCycleEscalation`
still passing on both hosts, since it drives this same skill through a fixture that
supplies no briefing or log and grades hand-written `- Cycle N: REJECTED` lines plus
the escalation marker. It is the falsifying surface for step 2's conditional
projection and step 4's cycle-3 branch: if the rewrite makes publication
unconditional, or drops the third-cycle escalation, that journey reds while the
rejection journey stays green.

## Test plan

Two committed layers plus one-off run evidence; no new standing harness, no new
fixture directory, and no new contractlint check.

**Run evidence (not a test).** The step and completion-condition count, pasted into
the validation report. Cost: one command. It settles AC-1's existence half and
nothing more.

**Offline behavioral (deterministic, standing).** Extend the existing
`TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl` rather than
adding a test: it already drives `gates.RecordSemantic` twice (`validation/1` then
`validation/2`) and then prepares the gate, so it is the exact place to assert the
single-publication end state — one room, agreeing pointer — and to keep the
inverted no-invocation control that stops the oracle becoming a tautology. Add table
cases for the invocation counter: one call passes; two calls fail; zero fails; a call
at `validation/2` fails.

**Live (measurement, throwaway).**
`SPACEDOCK_LIVE_RUNTIME={claude,codex} go test -tags live -run TestLiveCommonRejectionFlow ./internal/ensigncycle/`
as the repair loop, one journey per run, on top of hz's landed branch. Each green is
confirmed by reading its stream for the single `gate record --round validation/1` and
the `state=open` prepare — not by the pass alone, since the pre-fix journey is
stochastic (measured 1 pass in 3 at one sha) and a bare pass count cannot separate
the fix from luck. The full-matrix stack run is AC-3's terminal proof. No live test
joins the standing suite.

The escalation journey (`runClaudeFeedback3CycleEscalationScenario`) is a regression
surface for this change, not a target: it drives the same skill, its fixture supplies
no briefing or log, and it grades `- Cycle N: REJECTED` lines plus the escalation
marker. Step 4's cycle-3 branch and step 2's conditional projection must leave it
passing; it runs in the same matrix.

## Stage Report: ideation

- DONE: Design the flow at 5 steps or fewer, each ONE action with ONE completion condition, single round publish under one stable round id; ground every choice in the two diagnosis reports' quoted streams and the adversarially-confirmed trap
  Five steps, five `**Done when**` clauses, one per step (counted, not asserted). Single publish verified: one `--round validation/1` call leaves exactly 1 room and a pointer whose id, briefing, and room-ref all agree on round 1; the current two-call flow leaves 2 byte-identical rooms and a pointer reading cycle 2 whose briefing is still round-1. Grounding is partial and named below.
- DONE: Coordination is binding: hz owns the binary state-commit fix and grading honesty and is IN IMPLEMENTATION now - message its ensign (spacedock-ensign-hz2ankag6f-implementation) to agree the oracle boundary before you finalize; do not absorb its scope; whichever lands second reconciles
  Three messages sent proposing a file-level split, the shared-file reconcile rule, and hz-lands-first; no reply received, so the boundary is proposed and UNCONFIRMED — the gate should treat hz's yes as an open item. No hz-owned file appears in this entity's expected surface. Also sent hz a spike finding that its AC-1 may not hold as written (see below).
- DONE: Design against origin/main (post-stack); implementation lands as a stack PR; oracle-side changes ride in step with the new shape, never as prose-greps
  All spikes ran against a binary built from `origin/main` (`0c6a2c32a`) in a namespaced worktree, because `internal/gates` differs between the working checkout and origin/main. Doing so caught two errors in an earlier draft of this design: the ordered-token contractlint pin it planned to retire was already deleted upstream by `723028f01`, and the structural step-count check it planned as the replacement is itself banned by the Proof policy (a paraphrase reds it, a smuggled second action passes it). Both were dropped; contractlint now goes untouched, verified by running the package against the new skill text in an isolated origin/main worktree.
- DONE: Ground every choice in the *two* diagnosis reports' quoted streams
  The second report arrived after the first commit and corrected the design's framing, so the body was revised and recommitted. Both reports are now cited: hz's from its entity body, `diag-rejection-flow-bisect`'s quoted streams for both hosts. Its confirmation table reproduces all four trap sub-claims independently and matches my own scripted spike byte-for-byte on every one.

### Summary

The two failure mechanisms do NOT reduce to the same cause, which is the correction the second diagnosis forced and the most consequential fact in this design: claude published its first round against a half-written log and then published again when the log was complete, so neither call carried the id and the count the oracle needs together — but it committed and re-entered the gate correctly. Codex published once, never called `gate prepare` at all, read an empty scheduler off its own uncommitted write, and reported success. The single-publish shape fixes claude outright and does nothing for codex, whose trap one publish arms as well as two. The redesign publishes once, immediately after the correction completes the round's log, then gives the commit its own unbundled step — because a missing commit is the only failure here that reports success. Two spike findings sharpen the design beyond the diagnosis: the `gate prepare` refusal is artifact-scoped rather than tree-scoped, and the `needs-preparation` row requires a non-rejecting verdict as well as a clean tree, which is what forces the gate re-entry to follow the re-review. Designing against origin/main rather than the working checkout was load-bearing rather than procedural: it removed two files from the surface and killed a check I had proposed and would otherwise have shipped as a banned prose-grep. Two items need a captain ruling: the seed's "near zero or negative" surface estimate is amended upward to +107 words (the five completion conditions are the fix, so cutting them to hit the number would remove it), and a three-line fixture-prose repair sits just outside this entity's declared out-of-scope line.
