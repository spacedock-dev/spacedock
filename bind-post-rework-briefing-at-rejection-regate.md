---
title: "Bind the post-rework Briefing when rejection flow re-enters the gate"
status: ideation
source: "se0 complete Opus live suite, 2026-07-28: after validation/1 rejection, successful advisory round recording, implementation rework, and cycle-2 PASS, the FO rebound the original round-1 Briefing as the ordinary final gate instead of the distinct post-rework Briefing."
score: 0.85
sprint: durable-decisions
group: gate
sprint-readiness: ready
issue:
id: zbcj98qfwtax61vxdzrf615e
gates:
    version: 1
    current:
        gate: gate:zbcj98qfwtax61vxdzrf615e:ideation
    records:
        - id: gate:zbcj98qfwtax61vxdzrf615e:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zbcj98qfwtax61vxdzrf615e-backlog-1
              briefing:
                id: briefing:zbcj98qfwtax61vxdzrf615e:backlog:attempt-1:revision-1
                digest: sha256:a8668228e65695fdea30226ee877edb1031da0356a36cca5b245d644c3434802
                digest-domain: canonical-bytes
                request-digest: sha256:70a3922bccbd3031f2b3e4b7f5921d1b081db5861350376d92d8f6be23b6cc35
                room-ref: ./bind-post-rework-briefing-at-rejection-regate/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zbcj98qfwtax61vxdzrf615e:backlog:1
                briefing: briefing:zbcj98qfwtax61vxdzrf615e:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T15:39:42.727381Z"
                decision: approve
                reason: 'Sprint conn: task protects decision integrity after rework, reuses existing machinery, and can be ideated independently of other critical lanes.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:zbcj98qfwtax61vxdzrf615e:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:zbcj98qfwtax61vxdzrf615e-ideation-1
              briefing:
                id: briefing:zbcj98qfwtax61vxdzrf615e:ideation:attempt-1:revision-1
                digest: sha256:83bca9143edd196d23ffa806c2bbf8a8439d6057e1ad80832a3f7696be5f86be
                digest-domain: canonical-bytes
                request-digest: sha256:e816cf506a7ace8da21bfc41db71abd78ccc9dc4cbd3f065bb2bb804651b0c33
                room-ref: ./bind-post-rework-briefing-at-rejection-regate/review/ideation/briefing-1
---

## Problem

A supported feedback-rejection journey can correctly retain the validation/1 advisory review round, complete rework and a second validation, then bind the stale pre-rework Briefing as the ordinary final gate. In the retained complete Opus run, `gate record --round validation/1` succeeded, cycle 2 passed, and the FO selected `rejection-task/inputs/briefing.json` again. PR run `30412397240` reproduced the same defect on Sonnet: two validations completed, but the final gate bound `briefing:rejection-task:validation:round-1`, so `assertReviewRoundRecorded` correctly rejected the advisory round being retained as a gate attempt.

A prior focused Opus run on the same candidate correctly created and bound a distinct `gate-validation-round2/briefing.json`. The behavior is therefore nondeterministic contract conduct, not a recorder failure or an oracle false positive. After feedback changes the candidate, final approval must spend a freshly bound post-rework Briefing; an advisory review-round package cannot silently become that gate.

The Opus and Sonnet executions of the shared rejection-flow journey are temporarily TODO under this task. The repository also currently carries a Codex skip for the same defect; remove it so the stated keep-active boundary is restored. Keep Pi, recorded-gate, keep-moving, and deterministic round/gate coverage active.

## Value invariant

An approval after correction is valid only when its ordinary gate Briefing is bound to the newest recorded correction episode and presents the reworked, revalidated candidate. The rejected snapshot remains advisory evidence. It can never receive, borrow, or satisfy the later approval.

This is decision integrity, not a sequencing preference: a stale approval is invalid even when an FO happened to perform the expected prose steps.

## Mechanical definition: newest post-rework episode

Choose the **correction-cycle record** as the authoritative key. For a gate stage `S`, the newest post-rework episode key is the current immutable `review-round.id` when `review-round.stage == S`. The episode begins with the successful atomic publication of `gate record --round S/N` and remains current until a later same-stage round advances that pointer. Its ID already contains entity, stage, and cycle; the pointer already binds the rejected Briefing ID, digest, and immutable room; and the recorder already projects the matching `Cycle N` line. Feedback Cycles legibility, freshness, and binding therefore share one source.

An open current-stage gate is fresh after a correction round iff:

1. its attempt Briefing binding carries `correction-round: <current review-round.id>`;
2. its canonical prepared Briefing contains exactly one ordinary Context `Reference` to the current round Briefing, pinned through the existing Reference ID/URI/rev fields and therefore covered by the existing Briefing and request digests;
3. the referenced round still validates from its immutable room; and
4. the gate Briefing ID differs from the round Briefing ID.

Missing, older, mismatched, duplicate-Reference, or same-Briefing correction authority is stale. A later same-stage round mechanically invalidates an older open attempt. Closed historical attempts are preserved and are not retroactively required to bind a round recorded after they closed. A gate with no current same-stage correction round keeps today's behavior.

The cycle record is preferable to a rework commit because split-root workflows can place state evidence and deliverable commits in different Git histories; comparing their commit clocks would invent a false ordering. An explicit epoch would duplicate the round recorder. The cycle ID is the existing cross-root semantic boundary. The rework and second validation remain proven by their stage reports and live journey; the correction binding proves that the approval is spent on that episode rather than on the rejected snapshot.

## Proposed approach

Extend only the attempt's existing Briefing binding with the recorder-derived round identity:

```yaml
correction-round: round:rejection-task:validation:1
```

`gate prepare` reads and validates the current `review-round` pointer while it already holds the entity lock. When the pointer is for the gate's current stage, prepare writes that ID to `correction-round` and automatically appends the retained round's `briefing.json` as an ordinary canonical Briefing Context `Reference`. The Reference uses the existing deterministic Reference ID plus `git-root` URI and raw-byte `rev`; its bytes remain covered by the existing canonical Briefing digest and request digest. There is no new canonical Briefing field, flag, caller-supplied ID, metadata reconstruction, second recorder, or compatibility branch.

Entity-aware status validation compares the scalar `correction-round` with the current same-stage pointer without room or Git-object IO. Retained-authority validation locates exactly one deterministic correction Reference, compares its stored Source with `gitsource.Inspect(<round-room>/briefing.json)` through existing `SameLogicalRevision` (same logical root, repository path, and raw bytes even after unrelated commits), resolves the frozen object, parses its Briefing ID, recomputes its canonical digest, and compares both to the current pointer. It also validates the immutable round room normally. Thus Briefing ID and digest are derived, not duplicated in the gate binding or a new manifest field.

Consequently `gate validate`, `gate record`, eligibility/consume, prepare replay, and room-backed recording all fail closed through their existing reads. The diagnostic names attempt ID, expected round ID, and actual/missing round ID, and tells the FO to prepare a new post-rework gate Briefing. Reference path/bytes/identity errors use the existing retained-authority diagnostic chain. Rejection is byte-clean.

Status discovery must treat a stale open binding as `gate-readiness: invalid` and exclude it from `ready_gates`; it must never advertise that stale snapshot as `awaiting-captain`. The durable observability dependencies are the entity attempt's `correction-round`, the canonical Briefing's ordinary pinned correction Reference, the immutable round room, the existing `review-round` pointer/Feedback Cycles projection, and the status fail-closed projection. No success-line grammar changes are needed.

Update the FO text only to point at the mechanical behavior. The binary, not the instruction, owns the guarantee.

## Acceptance criteria

**AC-1 (VALUE) - Every final approval after rejection is spent on the reworked and revalidated candidate, never the rejected snapshot.**
Verified by: the shared rejection-flow oracle requires two implementation reports, two validation reports with cycle 2 passing, the retained validation/1 advisory round, and one open ordinary gate whose distinct canonical Briefing has the exact pinned round Reference while its attempt names the newest `correction-round`. Negative controls substitute the rejected Briefing, omit the scalar or Reference, or bind an older cycle and must fail. Focused Opus and Sonnet journeys pass three consecutive clean runs each.

**AC-2 - Advisory round and binding gate roles are mechanically distinct.**
Verified by: unit and shared-fixture tests retain the original advisory reviewer/worker Resolutions under `review/validation/round-1`, while the later gate has a different Briefing ID/digest, names the round identity, and presents the retained round only as Context. Making the gate/round Briefing IDs equal or changing the Reference path, bytes, rev, or referenced Briefing identity/digest makes validation fail without writes.

**AC-3 - Freshness is derived from the newest recorded correction episode and stale binding is mechanically falsifiable.**
Verified by: table tests cover missing/older/malformed scalar identities plus absent, duplicate, wrong-path, wrong-rev, wrong-bytes, and same-Briefing correction References; advancing `review-round` from `validation/1` to `validation/2` makes an earlier open attempt invalid. `gate prepare`, `gate validate`, `gate record`, eligibility, and consume refuse stale state and preserve entity/room bytes.

**AC-4 - Operators and automation observe stale authority before presentation.**
Verified by: status fixtures project a fresh binding as `awaiting-captain`, project the equivalent stale binding as `invalid`, and omit the stale entity from `ready_gates`. Exact prepare replay remains a no-op, and initial gates without a same-stage correction round retain existing output and behavior.

**AC-5 - The quarantined cross-runtime rejection journeys are restored without weakening their oracle.**
Verified by: remove the Codex skip and Claude Opus/Sonnet TODO scope, retain the durable round/gate assertion, run the focused Codex journey once and Opus/Sonnet three consecutive times each, then run both complete Claude shared suites at the exact candidate tip. Pi and deterministic round/gate tests remain enabled throughout.

## Boundary

Do not weaken the semantic oracle and do not allow a review-round Briefing to double as a later binding gate. The observed nondeterminism is evidence against a prose-only correction: identical contract text admitted both correct and stale binding, so ideation must price the prose-only variant and reject it explicitly unless a falsifiable exercise disproves this finding.

Do not parse `### Feedback Cycles` prose for authority, compare timestamps, infer ordering across Git roots, add a canonical Briefing extension or duplicated correction digest, add a CLI-supplied epoch, or turn advisory Resolutions into gate Resolutions. The `review-round` pointer is authoritative; the body projection remains human-readable output of the same recorder. V1 is unreleased, so change the attempt binding directly and add no fallback for an open post-round gate lacking the new scalar/Reference composition.

The missing authority belongs in gate preparation and retained binding validation. Feedback routing only transports the completed package and should not record or calculate freshness. Keep the 3k recorder boundary unchanged: this task consumes its existing immutable round pointer/room and does not alter round publication, triage, projection, or advisory semantics. This task also has no concrete shared semantic boundary with q3vp: it neither changes worker reuse nor routing identity, so keep it independent.

## Alternatives and cost

- **Prose-only correction:** about 2 skill files and 8-15 inserted lines. Reject it. The same lifecycle prose produced both the correct distinct Briefing and the stale binding on retained Opus/Sonnet runs; wording cannot make bad on-disk state fail, cannot protect non-LLM callers, and cannot make status suppress stale authority.
- **Rework Git commit:** about 6-9 files and 180-260 inserted lines, plus cross-root ancestry rules. Reject it because state and deliverable evidence can live in unrelated repositories; a single commit order is not provider-neutral authority.
- **New explicit epoch recorder/flag:** about 10-14 files and 300-450 inserted lines plus a new command/storage surface. Reject it because `gate record --round` already creates the required monotonic stage/cycle identity atomically.
- **Duplicated correction object:** about 11-14 files and 300-420 inserted lines. Reject it after correction review: repeating round ID, Briefing ID, and digest in both the attempt and a new canonical manifest field creates synchronization states without adding evidence already carried by Reference URI/rev plus the pointer.
- **Chosen scalar + existing Reference composition:** about 11-13 files and 250-340 inserted lines. It adds one status-readable attempt scalar, reuses canonical Context/Reference and `SameLogicalRevision`, has no new command or manifest vocabulary, and makes stale state directly falsifiable.

## Expected surface and semantic budget

Baseline: **12 files, about 295 insertions and 35 deletions; tolerance ±3 files and ±100 insertions.**

- `internal/gates/model.go`, `prepare.go`, `io.go`: about 80 insertions for `correction-round`, automatic round Reference derivation, retained-authority comparison through existing gitsource helpers, and actionable stale diagnostics. `operation.go` needs no new manifest field.
- `internal/gates/prepare_test.go`, `round_test.go`: about 125 insertions for fresh/stale/Reference/replay/history/immutability controls.
- `internal/status/boot_identify_test.go`: about 35 insertions for fresh versus stale readiness.
- `internal/ensigncycle/shared_round_recording_test.go`: about 50 insertions to strengthen the durable oracle and negative controls.
- `internal/ensigncycle/codex_live_runner_test.go`, `claude_live_runner_test.go`: about 20 deletions to remove this task's quarantine and TODO scope.
- `docs/specs/gate-resolution-frontmatter-contract.md`: about 25 insertions documenting the scalar/Reference binding and stale refusal.
- `skills/feedback-rejection-flow/SKILL.md`, `skills/fo-gate-lifecycle/SKILL.md`: about 15 insertions/changes pointing operators at the binary-owned guard.

Declared semantic changes:

- **Command grammar:** none; no new flags or success fields.
- **Stored formats:** request-backed gate attempt Briefing bindings gain `correction-round` when a current same-stage round exists; canonical Briefing v1 gains only an ordinary existing Context `Reference`, not a new field.
- **Authority:** the existing latest `review-round` pointer becomes the sole correction-episode key for a later open gate; callers cannot supply it.
- **Runtime:** stale post-round open attempts fail retained-authority operations, status marks them invalid, and `ready_gates` omits them. Initial gates and closed history retain current behavior.

Any new command input, a second epoch store, canonical Briefing correction field, duplicated Briefing/digest authority, parsing of Feedback Cycles prose, cross-root commit ordering, or change to round advisory semantics is outside tolerance regardless of LOC.

## Test plan

Write the focused gate tests first, then implement the schema/validator, status projection, shared oracle, docs/skills, and finally live re-enablement.

1. Gate unit/fixture tests create and commit `validation/1`, prepare a later gate, and inspect resulting bytes. They fail if `correction-round` or the deterministic Context Reference is absent/mismatched, if `SameLogicalRevision` no longer tolerates an unrelated state commit, if the ordinary and advisory Briefing IDs coincide, or if retained round validation is skipped. Estimated 125 LOC, fixture/unit only.
2. A stale table drives prepare replay, validate, record, eligibility, and consume against missing/older scalar identities and absent/duplicate/wrong-path/wrong-rev/wrong-byte References, asserting nonzero diagnostics and exact before/after tree equality. It fails if any semantic seam accepts stale authority. Estimated 70 LOC, unit/CLI behavior.
3. Status fixtures compare otherwise-identical fresh and stale entities. They fail if stale state is `awaiting-captain` or enters `ready_gates`. Estimated 35 LOC, native status fixture.
4. The shared rejection oracle continues to prove real two-cycle state and adds exact correction linkage. Its controls replace the gate Briefing with the round Briefing and advance the round pointer. Estimated 60 LOC, host-neutral behavior fixture.
5. Live proof removes the skips, runs Codex once, Opus three times, Sonnet three times, then the full Opus and Sonnet shared suites at one commit. These runs fail on either stale binding or failure to complete rework/revalidation; transcript wording alone cannot pass them.

No new integration mechanism needs a throwaway spike. The focused baseline run
`go test ./internal/gates ./internal/status ./internal/ensigncycle -run 'TestRejectionFlowRoundRecordingDurableOracleAndNoInvocationControl|TestBootReadyGatesFailClosedLifecycleControls|TestPrepare' -count=1`
passed on 2026-07-30. `TestSameLogicalRevisionIgnoresUnrelatedCommitButNotPathOrBytes` also proves the existing Reference Source comparison distinguishes path/bytes while ignoring unrelated commit movement. Together they exercise the exact primitives this design composes: atomic round pointer/projection publication, canonical References, a later prepared gate, strict retained room/digest validation, and status fail-closed classification. The missing scalar/Reference cross-binding is the test-first change.

## Proposed documentation diff

```diff
--- a/docs/specs/gate-resolution-frontmatter-contract.md
+++ b/docs/specs/gate-resolution-frontmatter-contract.md
@@ Recorder lifecycle
+When the entity carries a current review-round for the gate stage, prepare derives a
+`correction-round` identity for the attempt and includes the round Briefing as a pinned
+ordinary Context Reference in the new canonical Briefing. A later open attempt is valid
+only while its scalar and Reference resolve the newest same-stage round. Missing or
+older correction authority fails closed; the advisory round Briefing itself can never
+serve as the binding gate Briefing.
```

```diff
--- a/skills/feedback-rejection-flow/SKILL.md
+++ b/skills/feedback-rejection-flow/SKILL.md
@@
-7. Re-enter the normal gate flow with the updated result.
+7. Re-enter the normal gate flow with the updated result. Normal `gate prepare`
+   derives and binds the latest recorded correction round; never reuse or present
+   the advisory round Briefing as the later gate.
--- a/skills/fo-gate-lifecycle/SKILL.md
+++ b/skills/fo-gate-lifecycle/SKILL.md
@@ Prepare and bind
+After a correction round, rely on `gate prepare` to derive its authority. A stale-
+correction diagnostic halts presentation; never reconstruct or hand-edit the binding.
```

## Stage Report: ideation

- DONE: Define the value invariant that final approval is always spent on the reworked and revalidated candidate, never the rejected snapshot.
  The Value invariant and AC-1 make approval-to-candidate integrity the outcome, with durable two-cycle and stale-snapshot negative proof.
- DONE: Choose one mechanical definition of the newest post-rework episode and prove gate preparation or binding refuses a stale pre-rework Briefing without reconstructing metadata.
  The existing latest same-stage correction-round ID is the key; prepare derives it and retained-authority validation rejects missing, older, or mismatched bindings.
- DONE: Reuse existing gate and round recording machinery, declare expected files/LOC/semantics, and price and reject a prose-only solution.
  The plan extends current round/prepare/validation seams across 13 files/~360 insertions, declares four semantic axes, and rejects the 8-15-line prose option as non-falsifiable.

### Summary

Ideation now binds post-rework approval to the existing correction-cycle record, including exact freshness rules, observability, negative controls, surface budget, and live proof. It adds no recorder or caller-supplied metadata, rejects prose-only and cross-root commit ordering, and keeps q3vp independent.

## Stage Report: ideation (cycle 2)

- DONE: Preserve the value invariant and correction-cycle key while pricing the smallest existing-primitives composition.
  The revised design retains approval-to-reworked-candidate integrity and the latest same-stage `review-round.id`, but replaces duplicated authority with one scalar plus an existing Context Reference.
- DONE: Either revise to the leaner contract or provide a falsifiable reason existing Reference/binding primitives cannot carry the guarantees.
  Existing `git-root` Source, canonical Reference, request/Briefing digest, `SameLogicalRevision`, and immutable round validation carry presentation, path/byte pinning, and retained authority; no blocking gap remains.
- DONE: Re-estimate expected files, LOC, and semantics without changing the 3k recorder or q3vp boundaries.
  The baseline falls to 12 files/~295 insertions; there is no manifest extension, new recorder, compatibility path, 3k coupling, or q3vp coupling.

### Summary

Correction review produced a smaller contract: status compares only `correction-round`, while retained validation derives the round Briefing identity and digest from a normal pinned Context Reference. This removes duplicated authority and about 65 estimated insertions without weakening stale refusal, presentation evidence, or the original value invariant.
