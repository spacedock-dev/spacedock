---
id: 18963egcskzxaje6b5vnas3q
title: Simplify the feedback-rejection flow to five steps or fewer
status: ideation
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
---

Redesign skills/feedback-rejection-flow/SKILL.md from 8 steps to at most 5, with the single-publish shape the diagnostic evidence recommends: accumulate the round, publish ONCE under one round id, and make the gate re-entry a single unbundled step. Every step becomes one action with one completion condition; no step may bundle a tail an FO can drop.

Coordination, binding: hz (repair-codex-rejection-round-recording, in implementation) owns the binary side (inline state commit) and grading honesty; its reconciliation of the two diagnosis accounts determines whether the oracle repairs land there or the oracle re-anchors to this redesign's single-publish shape. Whichever lands second reconciles. The live rejection-flow journey and its oracles are the proof surface for both.

## Problem

Today's flow is eight steps. Both live failure mechanisms are located in the two
steps that carry more than one action, and both were reproduced deterministically
from a scripted fixture against a binary built from `origin/main` (`0c6a2c32a`) —
not inferred from the streams.

**One rejection cycle is published twice, under two round ids.** Step 6 publishes
`--round validation/1`; step 8 publishes `--round validation/2` after the reviewer
re-run. Driving both calls with the same briefing and log produces:

    round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4
    round=round:rejection-task:validation:2 stage=validation cycle=2 briefing=briefing:rejection-task:validation:round-1 entries=4

Round 2 is a byte-identical duplicate of round 1 under a different id: two rooms
(`review/validation/round-1/` and `round-2/`, same two files), and a `review-round`
pointer reading `id: round:rejection-task:validation:2` whose own briefing is still
`briefing:rejection-task:validation:round-1`. The end state is therefore ambiguous,
and the oracle absorbs the ambiguity: `assertRejectionRecordedRound`
(`shared_round_recording_test.go:184-189`) resolves `validation/1`, and on the
"pointer does not resolve" error retries as `validation/2`. That fallback exists
only because the flow publishes twice.

**The round can be published before it is complete.** A round is the reviewer's two
entries plus the correction worker's two disposition entries — all four carrying the
same briefing id (`rejectionCompleteLog`, `shared_fixtures_test.go:109`). Publishing
with only the reviewer's half exits 0 and prints `entries=2`. The oracle's success
regex pins `entries=4` (`shared_round_recording_test.go:17`), so a successful early
publish is graded as `rejection-round-missing`. Steps 4, 5 and 6 spend three steps
expressing one condition — "the worker's entries are in the log" — and the publish
happens in the third of them.

**The tail of step 8 is droppable, and dropping it is silent.** Step 8 bundles four
actions: append the Cycle line, publish `validation/2`, invoke
`spacedock:fo-gate-lifecycle`, run `«gate.lifecycle»`. An FO that stops after the
publish leaves the entity uncommitted, and the uncommitted entity is invisible:

| state | `status --next --json` | `gate prepare --artifact <entity>` |
|---|---|---|
| after `gate record --round`, uncommitted | `{"dispatchable":[],"ready_gates":[]}` | exit 1, `selected source differs from its committed Git object` |
| committed | row appears once the verdict is non-rejecting | exit 0, `state=open` |

An empty scheduler with no error reads as "the run is finished". That is the
completion gap behind `ErrNoGateRecord`: `gates.Read` returns "entity has no gates
record" for an entity whose gate was never prepared, and the codex-only clause at
`claude_live_runner_test.go:400-402` reports it under `rejection-round-missing`.

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
- `internal/contractlint/feedback_rejection_publication_smoke_test.go`: retire
  `TestFeedbackRejectionPublishesCompleteRoundBeforeRegating`. It pins five ordered
  substrings of the skill body ("append its authorized line", "Re-enter the normal
  gate flow", …) and would need only rewording to pass — it proves nothing about
  behavior and would ratify any bundling that kept the tokens in order. Its actual
  guarantee (a complete round is published before the reviewer re-run and gate
  preparation) is carried behaviorally by the counter above plus the existing
  `entries=4` pin. `TestDevelopmentTemplatePublishesRejectedRoundBeforeRegating` in
  the same file stays: its phrase remains true, since publication is still step 2 and
  the re-run and gate are steps 4 and 5.
- A structural check replaces it: the skill has exactly five numbered steps and
  exactly one `**Done when**` clause per step. This counts the deliverable itself
  rather than proxying for behavior with wording, and it is the only check that can
  catch re-bundling.

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

Six files, roughly 135 insertions and 45 deletions. Tolerance ±40 lines and ±2 files.

| File | Change |
|---|---|
| `skills/feedback-rejection-flow/SKILL.md` | 8 steps → 5 (~+35/-20) |
| `internal/ensigncycle/shared_round_recording_test.go` | drop the round-id fallback, add the invocation counter (~45) |
| `internal/ensigncycle/claude_live_runner_test.go` | wire the fourth `durableSemantic` (~4) |
| `internal/contractlint/feedback_rejection_publication_smoke_test.go` | retire the ordered-token test (~-30) |
| `internal/contractlint/structural_checks_test.go` | step / completion-condition count (~25) |
| `internal/ensigncycle/shared_fixtures_test.go` | fixture Cycle-line prose (~3), only if the gate rules it in scope |

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
entity. Contract surface: the skill loses its ordered-token contractlint pin and
gains a structural step/condition count.

## Risk evidence and spike record

The riskiest unverified mechanism was whether a **single** publication reaches the
end state the oracle demands, since removing the second call could have been what
advances the pointer or prepares the room. Exercised before designing, on a fixture
rebuilt from `writeRejectionWorkflow` against an `origin/main` binary:

- One `gate record --round validation/1` with the complete 4-entry log emits exactly
  the line the oracle's success regex pins, `entries=4`, and creates exactly one
  room with the two canonical files.
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
Verified by: a structural check counting numbered steps and `**Done when**` clauses
in `skills/feedback-rejection-flow/SKILL.md` and requiring exactly five of each, one
per step. It fails if a sixth step appears, if a step gains a second completion
condition (a re-bundled tail), or if a step loses its only one. This asserts a
mechanism and counts only in service of AC-2 and AC-3.

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

**AC-4 — The retired prose-grep leaves no coverage hole, and the suite stays green.**
Verified by: `go test ./...` plain and `-race`; and the specific guarantee of the
retired `TestFeedbackRejectionPublishesCompleteRoundBeforeRegating` — a complete
round published before the reviewer re-run and gate preparation — demonstrated by
AC-2's live counter (exactly one publication, at `validation/1`, carrying the
complete `entries=4` summary) together with AC-1's step ordering. The falsifying
change is a flow that publishes after the re-review: it passes the retired token
order and fails the counter.

## Test plan

Three layers, cheapest first; no new standing harness and no new fixture directory.

**Structural (instant, standing).** The step / completion-condition count in
`internal/contractlint`. Cost: trivial. It is the only guard against silent
re-bundling, and it replaces a strictly weaker ordered-token check.

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
  All spikes ran against a binary built from `origin/main` (`0c6a2c32a`) in a namespaced worktree, because `internal/gates` differs between the working checkout and origin/main. The new oracle counts recorder invocations in the run stream; the existing ordered-token prose-grep (`TestFeedbackRejectionPublishesCompleteRoundBeforeRegating`) is retired rather than reworded, with its guarantee re-carried behaviorally.
- FAILED: Ground every choice in the *two* diagnosis reports' quoted streams
  Only one of the two was available. hz's report was read in full from its entity body. `diag-rejection-flow-bisect` was asked for its quoted streams, its codex step-8 tail observation, and its task-#17 confirmation table, and did not reply. Rather than cite a report I never read, I reproduced both failure mechanisms myself from a scripted fixture against the origin/main binary, and every claim in the body is either quoted from hz's report or from output I generated. Two claims in the seed remain unverified by me: the exact claude stream showing `entries=2`, and which step the codex FO stopped at.

### Summary

Both failure mechanisms reduce to the same cause and both are now reproduced deterministically rather than inferred: one rejection cycle is published twice under two round ids, and the publication's durability tail is droppable and its omission silent. The redesign publishes once, immediately after the correction completes the round's log, then gives the commit its own unbundled step — because a missing commit is the only failure here that reports success. Two spike findings sharpen the design beyond the diagnosis: the `gate prepare` refusal is artifact-scoped rather than tree-scoped, and the `needs-preparation` row requires a non-rejecting verdict as well as a clean tree, which is what forces the gate re-entry to follow the re-review. Two items need a captain ruling: the seed's "near zero or negative" surface estimate is amended upward to +107 words (the five completion conditions are the fix, so cutting them to hit the number would remove it), and a three-line fixture-prose repair sits just outside this entity's declared out-of-scope line.
