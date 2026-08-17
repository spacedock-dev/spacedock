---
id: 18963egcskzxaje6b5vnas3q
title: Simplify the feedback-rejection flow to five steps or fewer
status: validation
source: "Captain principle 2026-08-16: anything more than 4 or 5 steps is a footgun. Live evidence: both rejection-flow failure mechanisms are step-count casualties - codex drops step 8's bundled tail (completion gap, ErrNoGateRecord); claude's oracle splits across the step-6/step-8 two-call publish (entries=2 vs 4, wrong round id)."
started: 2026-08-16T04:26:21Z
completed:
verdict:
score: "0.85"
worktree: .worktrees/spacedock-ensign-simplify-feedback-rejection-flow
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
        - id: gate:18963egcskzxaje6b5vnas3q:validation
          stage: validation
          attempts:
            - id: gate-attempt:18963egcskzxaje6b5vnas3q-validation-1
              briefing:
                id: briefing:18963egcskzxaje6b5vnas3q:validation:attempt-1:revision-1
                digest: sha256:93cdb73ed599f9bc19492ddbb251e05804fc978ee3bce2ff5df2a981b0f8fa35
                request-digest: sha256:0ae7d44d9a5ee2f7373139e07e7180bb8ccf2dfc8682da3cf5cc7949c8354b9b
                room-ref: ./simplify-feedback-rejection-flow/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:18963egcskzxaje6b5vnas3q:validation:1
                briefing: briefing:18963egcskzxaje6b5vnas3q:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-17T02:56:32.971274Z"
                decision: approve
              application:
                target-stage: done
                state: pending
pr: pr-merge:718
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
| committed | row appears once the stage report is structurally complete | exit 0, `state=open` |

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
   The sufficient condition is a *structurally* complete stage report for the gate
   stage: `hasCompleteStageReport` (`internal/status/entered_stage.go:36-52`) requires
   every checklist bullet's status token to be literally `DONE` or `SKIPPED`, each
   with evidence, plus a `### Summary`.

   **Corrected — an earlier draft of this entity got this wrong.** It reported that
   the *verdict* gates the row, from a four-way comparison that varied the bullet
   token and the verdict together and attributed the effect to the wrong one. The
   diagnostic agent challenged it against the parser, and the decisive isolation
   settles it, clean tree in both cases:

   | report | row |
   |---|---|
   | `- DONE:` bullet, "Recommendation: REJECTED" in prose | `needs-preparation` |
   | `- FAILED:` bullet, "PASSED" in prose | none |

   The scheduler reads no review semantics whatsoever. `extractChecklist` parses any
   `- WORD:` bullet and the completeness loop then rejects any token outside
   `{DONE, SKIPPED}`; `entered_stage_test.go`'s `"failed item"` case pins exactly this
   with a `- FAILED:` bullet. Earlier no-row results came from `- FAILED:` bullets,
   not from rejection.

   **This inverts the consequence, and the new one is a hazard the flow must handle.**
   A REJECTED validation report written the ordinary way — `- DONE:` bullets with the
   verdict in prose, which is what `writeRejectionWorkflow`'s README implies, since it
   mandates no bullet vocabulary — is structurally complete. So after step 3 commits,
   a `needs-preparation` row can appear **while the entity still carries only the
   rejected report**. The gate is not mechanically withheld until the re-review; the
   FO can walk straight into preparing a gate over un-re-reviewed work. Step 5's line
   that a `needs-preparation` row is work to perform, not a stopping condition, is
   therefore only true *at step 5* — see the ordering note below.

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
- The re-review precedes the gate for **contract** reasons, not mechanical ones.
  An earlier draft claimed the scheduler withheld the row until a non-rejecting
  verdict, so the ordering was forced for free; finding 2 above retracts that. The
  ordering stands on what the gate is *for*: it presents corrected work for a
  decision, and `assertRejectionRoundGateBoundary` requires the prepared attempt to
  bind `briefing:rejection-task:validation:attempt-1:revision-1` — the briefing
  prepared over the re-review's report, not the rejected one. The conclusion did not
  move; its justification did, and the difference matters because the mechanical
  version implied a guard that does not exist.

**Ordering note, and the one clause implementation must add.** Because a structurally
complete *rejected* report on a clean tree already yields a `needs-preparation` row,
step 5's "a `needs-preparation` row is work to perform, not a stopping condition"
must not read as licence to prepare a gate the moment such a row appears. Step 5 is
reached only after step 4's completion condition holds. Implementation should bind
that explicitly — the cheapest form is one clause in step 5 along the lines of *this
step is reached only from step 4's verdict; a `needs-preparation` row observed before
that is step 3's commit landing, not a gate to prepare* — so the flow does not hand
the captain a gate over work no reviewer has seen since the rejection. This is a
correction to an approved design, not new scope: it stays inside
`skills/feedback-rejection-flow/SKILL.md`, which the approved surface already owns,
and it costs roughly one line of the +107 already granted.

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
  structurally complete stage report — every bullet `DONE`/`SKIPPED` with evidence,
  plus a Summary. The four-way comparison first recorded here was confounded (it moved
  the bullet token and the verdict together); the corrected two-case isolation and the
  parser locus are under Problem, finding 2.

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

The two failure mechanisms do NOT reduce to the same cause, which is the correction the second diagnosis forced and the most consequential fact in this design: claude published its first round against a half-written log and then published again when the log was complete, so neither call carried the id and the count the oracle needs together — but it committed and re-entered the gate correctly. Codex published once, never called `gate prepare` at all, read an empty scheduler off its own uncommitted write, and reported success. The single-publish shape fixes claude outright and does nothing for codex, whose trap one publish arms as well as two. The redesign publishes once, immediately after the correction completes the round's log, then gives the commit its own unbundled step — because a missing commit is the only failure here that reports success. Two spike findings sharpen the design beyond the diagnosis: the `gate prepare` refusal is artifact-scoped rather than tree-scoped, and the `needs-preparation` row needs a structurally complete stage report as well as a clean tree. The second of those was first recorded WRONG here — as "a non-rejecting verdict" — from a confounded comparison that moved the bullet token and the verdict together; the diagnostic agent challenged it against the parser after the gate closed, and the corrected isolation reverses its design consequence, since an ordinarily-written rejected report is structurally complete and so DOES surface a row before any re-review. The five-step ordering is unaffected but its justification is now contract-based rather than mechanical, and implementation carries one added clause so step 5 cannot fire on that early row. Designing against origin/main rather than the working checkout was load-bearing rather than procedural: it removed two files from the surface and killed a check I had proposed and would otherwise have shipped as a banned prose-grep. Two items need a captain ruling: the seed's "near zero or negative" surface estimate is amended upward to +107 words (the five completion conditions are the fix, so cutting them to hit the number would remove it), and a three-line fixture-prose repair sits just outside this entity's declared out-of-scope line.

## Stage Report: implementation

- DONE: Execute the gated design exactly: the 5-step single-publish flow with the unbundled commit step, plus the captain-approved three-line fixture-prose repair; every step one action, one completion condition
  `skills/feedback-rejection-flow/SKILL.md` at 884b55af0: 5 numbered steps, 5 `**Done when**` clauses, one per step (counted on the shipped file, not asserted). The gated blockquote transcribed to 34 lines / 679 words / 4802 bytes — the ideation's declared baseline to the byte — then amended by the four authorized additions below. Fixture prose repaired in `shared_fixtures_test.go` (+3/-1): the FO writes one Cycle line per rejection round; the unprojectable "recorder projects Cycle 1" instruction is gone.
- DONE: Oracle changes ride in step with the new shape, never as prose-greps
  `assertRejectionRecordedRound` drops the `validation/1`→`validation/2` fallback, pins `validation/1`, and requires exactly one `round-N` room (filtered on the `round-` prefix — `review/validation/` also holds the gate's own `briefing-N` rooms, which the first draft of this check did not know and the test caught). Falsification exercised, not asserted: re-introducing the fallback makes the new republication control return `<nil>` and the test fail with `republication control diagnostic = <nil>`. The same fixture state that origin/main's test asserted the oracle ACCEPTED is now asserted to fail. No contractlint file touched; no test reads the skill text.
- DONE: Publication counter wired as the fourth `durableSemantic`
  `assertSingleRejectionRoundPublication` requires exactly one successful `gate record --round`, at `validation/1`, read from the run stream on both hosts (`claude`/`codexRejectionRoundPublications`), wired as `rejection-round-publication-count`. `TestRejectionRoundPublicationCounter` runs 6 cases × 2 hosts; 5 of the 6 assert failure — two calls, zero calls, a lone `validation/2`, a failed call, and both calls chained into one shell command. Deleting the "want exactly 1" branch reds all five.
- DONE: SEQUENCING IS REAL: write the skill and fixture now, hold the live loop and the stack PR until hz's layer lands
  Skill, fixture and oracle are committed on `spacedock-ensign/simplify-feedback-rejection-flow` (884b55af0, 0e8dd3c49), based on `origin/main` 0c6a2c32a. No live run started, no PR opened, no rebase onto hz attempted — held for the FO's signal, per the coupling. hz is still in its mechanism-confirmation loop (task #14).
- SKIPPED: Value proof: targeted local rejection-flow runs on BOTH runtimes on the composed tree, then the full matrix on the stack PR
  Blocked by the same sequencing, not declined. Step 3's completion condition is unsatisfiable on today's binary — `«state.commit»` is a no-op on inline workflows until hz lands — so a live run now would measure the pre-fix tree and prove nothing. AC-3 is unmet and this entity cannot reach PASSED until these runs execute on the composed tree.
- DONE: `go test ./...` plain and `-race`
  Both exit 0 with no failures and no data races, on the tree at 884b55af0 (`PLAIN_EXIT=0`, `RACE_EXIT=0`). `TestLiveCommonFeedbackThreeCycleEscalation`'s offline surfaces are unaffected; its live half rides the same held matrix.

### Deviations from the gated design — four, all needing captain ratification

1. **Entity operand in step 2's documented command.** The gated text documents `gate record --round …` with no entity operand. Exercised against a binary built from this tree: exit 2, `Error: unknown gate flag: validation/1`. The real signature is `gate record <entity> --round …` (`internal/cli/cli.go:181`; operand dropped from the skill by ea6723acb). A by-the-book FO never recorded a round at all. Root cause found by hz's ensign, routed to this entity as binding by the FO (handoff 5835d833a).
2. **Usage error is not a hold.** Step 2's gated failure language turned that exit-2 into a terminal hold. It now distinguishes a malformed invocation (correct it and re-run — the same single publication) from a recorder refusal (hold). The counter agrees by construction: it counts only calls that exited successfully, so the retry cannot read as a second publication.
3. **Step 5's precondition clause, in its corrected form.** Written first from the entity's original finding 2, then rewritten after that finding's retraction (00a86e53a). Verified independently in `internal/status/entered_stage.go`: the completeness loop accepts only `DONE`/`SKIPPED` bullet tokens with evidence plus a Summary, and reads no review semantics — so a structurally complete REJECTED report does surface a row before any re-review. The shipped clause is the entity's prescribed wording.
4. **Surface overrun, concentrated in one file.** 4 files as declared (`±1` honoured), but +273/-78 (net +195) against the declared ~85/~25 (net ~+60) at ±30 lines. All of it is `shared_round_recording_test.go`: +249/-67 against its ~45 estimate. It splits roughly as the counter helpers and both hosts' stream readers (~120), the counter's table test with its five inverted cases (~74), the room check (~13), and the republication counterexample (~20). One trim already applied: a per-segment command splitter was replaced by reusing the existing launcher matcher, ~25 lines lighter. The rest is the counter as AC-2 specifies it; cutting the inverted cases to reach the number would leave the counter unfalsifiable, which is the failure the estimate was protecting against.

The skill also grew past the +107 words the gate approved: 572 → 784 words (4158 → 5419 bytes), 28 → 34 lines. Deviations 1-3 are the whole of the extra +105.

### Summary

The five-step single-publish flow, its oracle, and the fixture prose repair are committed and the full suite is green, but the value proof is not in: AC-3's live runs are held behind hz's inline `state commit`, without which step 3 cannot complete on an inline workflow. The most consequential thing found during implementation was not in the design at all — the gated step 2 documents a `gate record` invocation that exits 2, so an FO following the skill by the book never recorded a round. That is hz's ensign's discovery, routed here as binding, and it means the five-step shape alone would not have fixed the journey. Two other clauses moved too: a usage error no longer becomes a terminal hold, and step 5's precondition was written twice — once from the entity's original verdict-gating finding, then again after that finding was retracted, which I re-verified against the parser rather than taking either version on trust. The oracle work is honestly larger than estimated, entirely in the counter and its inverted controls; I trimmed what could go without weakening falsifiability and am reporting the rest rather than trimming into the proof.

### Addendum: recorder-command provenance and evidence status (FO-required)

Provenance of the two defects this rewrite closes, each verified by reading the shipped
skill text at the commit rather than by trusting the citation:

- `c355fbe44` (2026-07-23) documented the recorder correctly, with the operand:
  `${SPACEDOCK_BIN:-spacedock} gate record <entity> --round <stage>/<cycle> --briefing <path>/briefing.json --log <path>/briefing.review.jsonl`.
- `ea6723acb` (2026-08-03) rewrote that step line to `gate record --round STAGE/CYCLE
  --briefing PATH/briefing.json --log PATH/briefing.review.jsonl`, dropping the entity
  operand from the contract-facing command. The FO's provenance is confirmed on this
  point. I first read it as contradicted, because the same file still carries fixture
  examples that DO name the entity (`gate record rejection-task …`) and a naive grep
  returns those first; the step line at `ea6723acb:21` is the one the FO follows.
- `d05730564` (2026-08-09) added the invoke-once / hold-on-failure language, which
  turns the resulting usage error into a terminal hold. `d865e8b2a` (2026-08-12)
  carried both forward to the text this entity replaced.
- One citation correction: `00f8c203f` (2026-07-30) carries no recorder command line at
  all in this file, so it is not a prior-correct witness. `c355fbe44` is.

Discovery attribution: the root cause was found by the repair-codex-rejection-round-recording
ensign's live-loop analysis and FO-verified against `internal/cli/cli.go:181`. This
entity's contribution was to exercise the defect directly — the documented command
exits 2 with `Error: unknown gate flag: validation/1` — and to close both defects in
the rewrite.

Layer relationship, stated plainly for the gate: the repair layer's inline
`state commit` fix is **necessary but not sufficient**. A by-the-book FO never reached
the dirty-tree trap at all, because the documented command failed first. The journey
cannot go green on the full matrix until this rewrite lands, so this layer is
load-bearing for the journey AC and the composed-tree loop at the stack tip is where
the journey green gets proven.

Evidence status, recorded so it is not silently reused: the live run in which the FO
recorded both rounds and then declared the journey complete without ever invoking
`gate prepare` was driven under bare codex defaults, not the CI PATH-shim model
configuration. It is withdrawn as an asserted property of the journey and is NOT cited
as oracle evidence anywhere in this report. Step 5's clause stands on the
needs-preparation mechanism resolution and the parser reading, which predate that run.

FO dispositions received on the deviations above: the handoff-channel clauses are
authoritative (identical to the direct relay); the gate must state plainly that the
five-step shape alone would not have fixed the journey; the step-5 retraction handling
is correct; and the surface overrun is a KEEP ruling — the publication counter and its
five inverted cases stay, since falsifiability outranks the line estimate, and the
+195 net rides as a declared deviation.

### Addendum 2: final surface figures, and the plainness the gate requires

**Required plain statement (FO disposition, not optional).** The gated five-step shape
alone would NOT have fixed the journey. The gated design's step 2 documented
`gate record --round STAGE/CYCLE …` with no entity operand, which exits 2 without
recording anything, so an FO following the shipped five-step flow verbatim would have
failed at step 2 on every attempt — before ever reaching the commit step the shape was
designed around. **The ENTITY operand clause is load-bearing for the journey AC**, on
equal footing with the single-publish shape itself: neither alone carries the journey.
The repair layer's inline `state commit` fix is likewise necessary but not sufficient,
because a by-the-book FO never reached the dirty-tree trap it fixes.

**Final surface, superseding the figures in deviation 4 above.** Those were measured at
884b55af0; `ba91b333f` then added the quote-run tolerance and its table case. At
`ba91b333f`, against base `origin/main` 0c6a2c32a:

| file | change |
|---|---|
| `internal/ensigncycle/shared_round_recording_test.go` | +261 / -67 (net +194) vs a ~45 estimate |
| `skills/feedback-rejection-flow/SKILL.md` | +16 / -10 |
| `internal/ensigncycle/claude_live_runner_test.go` | +5 / -0 |
| `internal/ensigncycle/shared_fixtures_test.go` | +3 / -1 |
| **total** | **4 files, +285 / -78, net +207** vs declared ~85 / ~25 (net ~+60), tolerance ±30 |

File count is exactly as declared; the line overrun is entirely the one oracle file.
FO disposition: KEEP — the publication counter and its inverted table cases stay,
because a counter that cannot be falsified is worse than extra test lines, and no
further trimming is to be attempted. The ~25-line trim already applied (a per-segment
command splitter replaced by the existing launcher matcher) is the whole of it.

**Skill growth attribution.** 572 → 784 words, 28 → 34 lines, 4158 → 5419 bytes. The
captain approved +107 words at the ideation gate, which the transcribed design hit to
the byte (679 words). The further +105 words are wholly the three FO-authorized
clauses — the ENTITY operand sentence, the usage-error≠hold wording, and step 5's
precondition clause — and nothing else.

### Addendum 3: stack PR opened as the first layer

Captain order 2026-08-16 reversed the layer order: this entity opens the stack's FIRST
PR and the repair layer rebases onto it, so this branch was NOT rebased onto anything.
It remains based on `origin/main` 0c6a2c32a, which is correct as-is.

- **PR: #718** — https://github.com/spacedock-dev/spacedock/pull/718
- Branch `spacedock-ensign/simplify-feedback-rejection-flow` -> base `main`, state OPEN,
  mergeable, not a draft. Carries exactly the three commits 884b55af0, 0e8dd3c49,
  ba91b333f. Read back with `gh pr view 718 --json`, not trusted from the create banner.
- Pushed over the SSH remote `git@github.com:spacedock-dev/spacedock.git`.

**Stack reference: none yet, and that is the correct state.** `gh stack link 718` fails
with `requires at least 2 arg(s), only received 1` — a stack cannot exist with a single
member, so `gh stack view` reports this branch is not part of a stack. The stack forms
when the repair layer's PR exists; the command to run then, bottom-to-top, is
`gh stack link 718 {repair-PR}`, which makes #718 the bottom and sets the repair PR's
base to this branch. Worth noting for whoever runs it: the failure above still exited 0
through a pipeline, which is the same print-success-on-failure hazard already recorded
for this command — read the stack back with `gh stack view --json` either way.

The PR body carries the layer-status plainness required for the gate: this layer is
necessary but not sufficient, it fixes the documented-command and hold-on-usage-error
defects that stopped a by-the-book FO before the dirty-tree trap, and full-matrix
journey green is expected only at the stack tip once the repair layer lands on top. No
live loop was run from this layer; that proof moves to the composed tree at the tip.

`pr:` frontmatter is deliberately NOT set here — entity frontmatter is the first
officer's write scope, not an ensign's. The FO sets `pr=#718` per the pr-merge mod.

## Stage Report: validation

- DONE: Validate on the stack tip 52b94f9aa; entity stage report and amendments as source of truth
  Throwaway detached checkout of 52b94f9aa; ba91b333f verified as its ancestor; PR #718 OPEN on base main with exactly 884b55af0, 0e8dd3c49, ba91b333f (`gh pr view --json`).
- DONE: Verify the gated skill rewrite on the tip's bytes; byte-compare against the approved design's wording where the entity pins it
  No commit above ba91b333f touches the skill; tip `wc` = 34 lines / 784 words / 5419 bytes exactly as Addendum 2; 5 numbered steps, 5 `**Done when**`. Word-diff of the entity's gated blockquote vs the tip section: byte-identical except exactly the three FO-authorized insertions (ENTITY operand + sentence, usage-error≠hold wording, step-5 clause binding step 5 to step 4's verdict). Operand defect re-exercised on a tip-built binary: operand-less `gate record --round …` exits 2 `unknown gate flag: validation/1`.
- DONE: Verify the oracle reshape claims by falsifying edits where the report declares them; verify the three clauses account for exactly +105 words over the approved +107
  On tip bytes: neutralizing the counter's count guard reds the 4 count-guarded inverted cases; neutralizing its round-id guard reds the other 2 — all 6 inverted cases load-bearing; restored → green. Re-introducing the validation/2 fallback reds the durable test with `republication control diagnostic = <nil>`, verbatim as the implementation report declared. Words: file 679→784, section 542→647, both +105; baseline 572→679 = +107; the word-diff contains nothing but the three clauses.
- DONE: Provenance: c355fbe44 alone as the prior-correct witness; entity records it and the naive-grep caveat
  c355fbe44 carries `gate record <entity> --round …`; ea6723acb:21 drops the operand; 00f8c203f carries zero recorder lines. Entity records all three plus the caveat. Polish: the caveat's "the same file" is imprecise — the operand-carrying examples live in `shared_round_recording_test.go` at that commit, not SKILL.md; the provenance conclusion is unaffected.
- DONE: Composed journey evidence: cite ledgers and tip CI; verify ledger artifacts exist where claimed; no live re-runs
  No live loop run. Ledgers committed (state 08242fe75) at `_evidence/zq-team-mode-ac1/`: codex 3/3 consecutive conforming greens, claude 3 consecutive (runs 3-5; run-2 red log retained), on composed tree 571017df3, which contains this layer. Tip CI 31986034831 (sha 52b94f9aa) concluded RED after dispatch: offline success; codex-live 1/17 failed = ACValueReanchor only (rejection and escalation green); claude-live 1/17 failed = RejectionFlow, observed=[rejection-cycle-line rejection-round-missing] — finding 1 below.
- DONE: Assess every declared deviation against the captain rulings recorded in the entity
  Layer diff verified: 4 files +285/-78, oracle file +261/-67, exactly Addendum 2's KEEP-ruling figures; no further trim attempted. Deviations 1-3 are the three verified clauses; the falsification above is the falsifiability the KEEP ruling protects.
- DONE: Validation stage report with per-AC verdicts, recommendation, path-scoped commit, push, FO signal; no gate preparation or resolution
  This report; no gate action taken.

### Per-AC verdicts

- AC-1 PASS — tip counts: `grep -cE '^[0-9]+\. \*\*'` = 5, `grep -c '\*\*Done when\*\*'` = 5; the behavioral half rides AC-2/AC-3 as designed.
- AC-2 PASS — durable fixture drive green on tip and falsified both directions; counter green (7 cases × 2 hosts, 6 inverted) and wired as a durableSemantic on both hosts. Live corroboration: the red claude CI run's observed codes exclude rejection-round-publication-count — one conforming publication at validation/1 even in a failing run.
- AC-3 PASS on the designated evidence — no validation/2 retry anywhere in the tip oracle (grep absent; republication control enforces); both-host journey greens are the committed conforming-green ledgers. Tip CI attempt: codex rejection journey green; claude rejection red, mechanism outside this layer's surface (finding 1).
- AC-4 PASS — `go test ./...` plain and `-race` green on tip, all packages ok (cli/ensigncycle need >600s on a loaded machine, so the 10m default trips environmentally: plain re-run 627s/709s ok, race `-timeout 30m` exit 0; CI offline job green on the same sha corroborates). Escalation journey passed live on both hosts in the tip CI run.

### Findings (recommendations for Review-finding disposition; candidate untouched)

1. Tip CI claude-live rejection red: the FO prepared a premature round-1 gate then withdrew it (selected gate holds ≠1 attempt) and left 0 durable Cycle lines while its final message claims one. This layer's own oracles graded green in that run; the codes are gate-attempt cardinality and the cycle line — the sonnet FO-adherence residual class the captain's 2026-08-16 ruling folded into run-rejection-journey-in-team-mode (whose ledger red is the same class, different sub-mechanism: unrecognized recorder invocation). Proposed: material for the stack-tip merge (required lane red), ownership zq/live-ci flake binding, route for decision at the gate.
2. Tip CI codex-live red: TestLiveCommonACValueReanchor only; no file overlap with this layer; cited for the stack picture, ownership with the reanchor-grading layer.
3. Polish: the Addendum's "the same file still carries fixture examples" location imprecision (see provenance item).
4. Deferred risk: the publication counter's entity recognizer requires a whitespace-delimited bare `rejection-task`; a path-qualified operand would under-count and false-red (safe direction — never a false pass). Unobserved across 8 ledger runs and 2 CI attempts. Promote if a live red reports "made 0 successful" while the stream carries a path-qualified publication.

### Summary

Recommend PASSED. Every mechanism claim was re-verified on the stack tip's bytes and falsified where the implementation declared falsification: the shipped skill is the gated design plus exactly the three FO-authorized clauses (+105 words to the word), the publication counter's six inverted cases all red when their guards are removed, and re-introducing the retired fallback reproduces the declared `<nil>` control failure verbatim. The journey-level evidence stands on the committed conforming-green ledgers on the composed tree; the tip CI attempt concluded red after dispatch on two codes outside this layer's surface — routed as finding 1 for the gate, since a red required lane is the stack's merge problem even though this layer's oracles graded green inside that very run.
