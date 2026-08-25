---
id: rx3daftacggfmw1pt2febw31
title: Gate rooms are request-less by default; the provider handoff is opt-in at prepare
status: implementation
source: "Captain gate-format review 2026-08-25: variance scan found 499/499 rooms carry request.json while only ~8 July-era provider results ever consumed one; q0 (subspace-r-scaffolded-gate-room, spacedock-subspace, at validation) finalizes the provider journey the file exists for"
started: 2026-08-25T17:27:46Z
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
        - id: gate:rx3daftacggfmw1pt2febw31:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rx3daftacggfmw1pt2febw31-backlog-1
              briefing:
                id: briefing:rx3daftacggfmw1pt2febw31:backlog:attempt-1:revision-1
                digest: sha256:9a936b06ee55e932f0ef8997e0b3b79810718d78076fa69df98af75a87eae2e8
                request-digest: sha256:e9dbc6437c5a986912593c1ec9cbe6e016f561b89f7d53818221ad5ac556af86
                room-ref: ./request-less-gate-rooms-by-default/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rx3daftacggfmw1pt2febw31:backlog:1
                briefing: briefing:rx3daftacggfmw1pt2febw31:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T17:27:27.228339Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''dispatch rx'' — approved seeding into design; q0 coherence requirement baked into the seed'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:rx3daftacggfmw1pt2febw31:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:rx3daftacggfmw1pt2febw31-ideation-1
              briefing:
                id: briefing:rx3daftacggfmw1pt2febw31:ideation:attempt-1:revision-1
                digest: sha256:c9e0b6fd6463605536e72f5e3d98b7d9df4d94a9565f8f7419e81b2d04205f06
                request-digest: sha256:d09993fd1005af3bbce7d4b015b39f367b6fea405eaa1c8e8ccdeeeb8592db77
                room-ref: ./request-less-gate-rooms-by-default/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rx3daftacggfmw1pt2febw31:ideation:1
                briefing: briefing:rx3daftacggfmw1pt2febw31:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T21:02:13.592893Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''rx - approve'' — accepts the predicate-centered design as presented, field drops cut per q0 evidence; approved baseline net +200 across 14 files, tolerance ±60/±3'
              application:
                target-stage: implementation
                state: consumed
---

`gate prepare` mints `request.json` for every room. Only the subspace provider journey reads it. Chat-presented gates are 100% of current practice. Almost every room therefore pays for a handoff that no consumer reads. Make the chat room one file. Mint the provider handoff only when the caller selects the provider channel.

## Problem

Every `gate prepare` publishes two files. `gate-briefing.json` holds the canonical Briefing. `request.json` holds the provider handoff. A scan of the dev state checkout on 2026-08-25 counts 499 of 499 prepared rooms with `request.json`, and about 8 July-era rooms with a provider result.

The request adds no fact that the entity does not already hold. Its `gate`, `attempt`, and Briefing id and digest all repeat the frontmatter binding. Its `actor` and `approver` are the constant `person:captain`, and `operation.go:530` refuses any other value.

The spec already permits the omission: "Request-less and chat-only attempts may omit it". `validateGateRoomRequest` (`operation.go:512`) returns nil on an empty `request-digest`.

That permission is not sufficient, and this is the defect the seed did not see. Eight code sites read `request-digest != ""` as the test for "this attempt has a prepared room":

- `prepare.go:245` — the frozen open binding cannot be rebound.
- `prepare.go:300` — the entity replay source, which lets a re-prepare reuse the frozen entity bytes.
- `prepare.go:351` — exact prepare replay, which makes a repeat prepare a no-op.
- `operation.go:178` — `gate withdraw`.
- `operation.go:551`, `operation.go:562` — the bound Briefing read.
- `application.go:244` — the reviewed-input staleness check.
- `io.go:204` — retained-authority validation.
- `model.go:281` — the frontmatter rule that a withdrawn attempt must retain a `request-digest`.

Two of these are measured, not read. A chat room that this design prepares makes today's `gate withdraw` exit 1 with `Error: current attempt is not request-backed`. And `io.go:204` skips retained-authority validation for a request-less binding, so a corrupted chat Briefing passes unnoticed.

So the task is not "stop writing a file". The task is to give the reader a channel-independent test for "this attempt has a prepared room", and to keep every guard that `request-digest != ""` accidentally carried.

## Proposed approach

Select the presentation channel at prepare. The default is chat. The caller opts in to the provider channel.

**1. `gate prepare` gains `--provider`.** Without the flag, prepare publishes `gate-briefing.json` alone and binds no `request-digest`. With the flag, prepare mints `request.json` and binds `request-digest` exactly as today.

- Value AC: AC-1.
- Simplest alternative: always mint the request, and let the chat path ignore it. That is today's behavior, and it is what AC-1 measures against. It cannot deliver AC-1.
- Second alternative: a `--present chat|provider` value flag. A boolean is smaller, and two channels exist. Reject the value flag until a third channel exists.
- The First Officer's command text does not change. `skills/fo-gate-lifecycle/SKILL.md:22` runs the chat form, which is the new default.

**2. `preparedRoomBinding(entityPath, binding)` replaces `request-digest != ""` at the six runtime sites.** It reports true when the binding's `room-ref` names a real directory that holds `gate-briefing.json` as a regular file.

- Value AC: AC-1. Without it, a chat attempt loses replay, rebind protection, and withdraw.
- Simplest alternative: a new frontmatter field, such as `channel: chat`. It costs a schema change, a validator rule, and a back-compat read for 600 archived bindings. The room shape already carries the fact.
- Why the room shape is unambiguous: a new room holds `gate-briefing.json`. The 101 archived request-less rooms hold `briefing.json` under the old name, or an opaque ref such as `subspace-room:3k-gate-design`. No archived request-less binding can match.

**3. `boundBriefingPath(entityPath, binding)` resolves the three binding shapes in one place.** A provider binding resolves through the frozen request locator. A chat binding reads the reserved locator in its room. A legacy binding names the Briefing file itself.

- Value AC: AC-1 and AC-2.
- Simplest alternative: add the room case at each of the two read sites. That duplicates the shape test, and `application.go` and `operation.go` then can disagree. One resolver keeps them equal.

**4. The withdraw guard moves from the frontmatter validator to the verb.** `Withdraw` requires a prepared room instead of a request-backed attempt. `model.go:281` drops the rule that a withdrawn attempt must retain a `request-digest`.

- Value AC: AC-4.
- Why the move is necessary: `Validate(doc)` reads frontmatter only. It cannot stat a room, so it cannot tell a new chat binding from a legacy opaque ref. `Withdraw` already stats the room at `operation.go:185`.
- Simplest alternative: keep the frontmatter rule and let chat gates stay non-withdrawable. That is the measured regression in AC-4. It is not acceptable.
- Declared narrowing: the pure validator no longer asserts this property. The mutating verb asserts it. The gate must approve that trade.

**5. Retained-authority validation extends to chat rooms.** `io.go` keeps the request checks for a provider binding, and it now runs the Briefing digest and Git-source checks for every prepared room.

- Value AC: AC-3.
- Why it is necessary: without it, the default channel gets less validation than today's default. That inverts the value of the task.
- Simplest alternative: skip it, because a chat room has no provider authority to protect. Insufficient: the Briefing digest and the selected Git objects are the authority a chat gate presents. They need the same check.

**6. Room-entry validation follows the published shape.** `validatePreparedRoomEntries` requires the exact file set that the room's channel publishes. `publishPreparedRoom` treats the request as optional.

- Value AC: AC-1.
- Simplest alternative: accept one or two files without pinning the set. That drops the "no copied sources, no provider subtree" invariant the spec states at line 125.

### Two seeded cleanups are cut, not deferred

The seed proposed to drop the per-item `type` field and to consider dropping `mediaType`. The spike refutes both against the shipped q0 preflight. Neither belongs in this task.

- **Per-item `type` stays.** The seed states "q0 does not read it". That is false. With `"type": "Reference"` removed from the context items, q0 refuses the room: `briefing: context 0 requires id and type`. In the accepted control run, q0 materializes the reference blob; in the typeless run it does not. Spacedock's own `referenceInventory` (`operation.go:758`) uses the same field to find References, so dropping it would also silence Git-source validation for every reference.
- **Artifact `mediaType` stays.** With `mediaType` removed from the primary artifact, q0 refuses the room: `briefing: artifact 0 requires id, uri, mediaType, and rev`. This is a resolved cross-repo requirement of the shipped `reviewv1` loader, not an open question and not a deferral.

## Risk evidence

The riskiest claim is cross-repo byte compatibility. The spike exercised it first, against the q0 preflight in `spacedock-subspace` at `spacedock-ensign/subspace-r-scaffolded-gate-room` (75ef1a2), through `internal/gateroom.Prepare`.

The spike built the design in a throwaway checkout of `95f877cd6` and drove the real CLI against a split-root fixture with two Git roots. Results:

1. **The provider opt-in reproduces today's room byte-for-byte.** In one fixture, at one pair of commits, the baseline binary and the patched binary with `--provider` produce identical `request.json` and identical `gate-briefing.json`. `diff -r` reports no difference.
2. **q0 accepts the patched provider room.** `PREFLIGHT-ACCEPTED gate=gate:task:validation attempt=gate-attempt:task-validation-1 actor=person:captain`. The scratch tree holds the materialized artifact blob and the materialized reference blob.
3. **q0 fails closed on a chat room, before any host command.** `PREFLIGHT-REFUSED: gate room: request.json: ... no such file or directory`. This is the correct outcome. A chat room is not a provider room, and q0 refuses it by name rather than by silent degradation.
4. **The canonical Briefing is channel-independent.** The chat room's `gate-briefing.json` is byte-identical to the provider room's, at the same commits.
5. **The withdraw regression is real, and the fix is load-bearing.** Against one chat room, the baseline binary prints `Error: current attempt is not request-backed` and exits 1. The patched binary prints `withdrawn ... state=withdrawn` and exits 0.
6. **The chat lifecycle closes and consumes.** `gate record --decision approve --consume` reports `state=closed`, `sync=local-only`, and `route=approved-awaiting-merge`.
7. **Archived rooms still read.** `status --workflow-dir docs/dev --validate` produces byte-identical output under the baseline binary and the patched binary, over the real state checkout of 600 rooms in every historical shape. Both exit 0.
8. **The blast radius is bounded and named.** `go test ./internal/gates/... ./internal/cli/...` fails 13 cases in 11 test functions. Every one asserts the two-file shape or the request-backed precondition. `TestWithdrawRefusalsLeaveEntityRoomAndLockBytesClean/chat-only_request-less_attempt` asserts the exact behavior this design inverts, so it must be rewritten, not repaired.

The spike code is throwaway. Its runs seed the implementation's first tests, named in the test plan.

## Documentation changes

**`docs/specs/gate-resolution-frontmatter-contract.md`, line 125.** Before: "It writes only `gate-briefing.json` and `request.json` at preparation time." After: "The caller selects the presentation channel. By default preparation writes only `gate-briefing.json`, and the binding carries no `request-digest`. With `--provider` it also writes `request.json` and binds its digest, for a provider that reads the room. Preparation copies no selected source, writes no association, and creates no provider subtree."

**Same file, line 107.** Before: "Request-less and chat-only attempts may omit it." After: "A chat attempt omits it. A binding whose `room-ref` names a directory that holds `gate-briefing.json` is a prepared room in either channel."

**Same file, line 176.** Before: "requires the room to contain exactly `gate-briefing.json` and `request.json`". After: "requires the room to contain exactly the file set its channel publishes".

**`docs/site/reference/command-reference.md`, line 96.** Before: "Immediately after successful preparation the room contains exactly `gate-briefing.json` and `request.json`, with no copied sources or association." After: "Immediately after successful preparation a chat room contains exactly `gate-briefing.json`, with no copied sources or association. Add `--provider` to also mint the `request.json` handoff that a room-reading provider needs." The command cell gains `[--provider]`.

**`docs/site/concepts/gates-and-decisions.md`, line 41.** Before: "two-file recorder-ready room". After: "recorder-ready room".

Implementation applies these diffs. Validation reviews them against ASD-STE100.

## Out of scope

Subspace-side changes (q0 owns its preflight and Result contract). The canonical five-key Result format. Retiring the provider machinery (the opt-in exists to keep it). The round recorder's `briefing.json` and `briefing.review.jsonl` shapes. Dropping the per-item `type` field and the artifact `mediaType` field, both refuted above. Converting existing archived rooms; they stay unchanged and readable.

## Expected surface and tolerance

Estimate net LOC change: +200, across 14 files (insertions ~+420, deletions ~-220). Tolerance: ±60 net LOC and ±3 files.

The production part is measured, not estimated. The spike's non-test diff is +144 / -82, net +62, across 6 files: `internal/gates/prepare.go`, `operation.go`, `application.go`, `io.go`, `model.go`, and `internal/cli/cli.go`. The remainder is 5 test files and 3 documentation files.

Declared semantic changes:

- **Command grammar:** `gate prepare` gains the `--provider` flag. No existing invocation changes meaning.
- **Stored format:** a chat gate room holds one file, and its binding carries no `request-digest`. Archived rooms and bindings are unchanged.
- **Authority:** `gate withdraw` accepts any prepared room, not only a request-backed one. The frontmatter validator no longer requires a withdrawn attempt to retain a `request-digest`.
- **Runtime behavior:** retained-authority validation now covers prepared chat rooms. A corrupted chat Briefing that passes today will be refused. The `status --validate` run over 600 archived rooms shows no change there.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A chat gate journey publishes a one-file room, and the room closes, consumes, and archives as it does today.**
Verified by: a behavior test that drives the built binary through prepare, record, and consume on a split-root fixture, and asserts the prepared room's entry set is exactly `{gate-briefing.json}` and the binding has no `request-digest`. Independent baseline that fails today: the same journey publishes 2 files and binds a `request-digest`, in 499 of 499 rooms scanned. Falsifying edit: make `Prepare` mint the request unconditionally. The entry-set assertion reds.

**AC-2 - A `--provider` prepare publishes the exact room bytes today's prepare publishes, and the q0 preflight accepts it.**
Verified by: a golden test that pins the `request.json` and `gate-briefing.json` bytes for fixed inputs, plus the recorded cross-repo run of `spacedock-subspace` `internal/gateroom.Prepare` against a `--provider` room. Falsifying edit: change any bound field, the key order, or the indent. The golden reds, and the q0 digest check reds.

**AC-3 - A prepared room's retained authority is validated in both channels.**
Verified by: a test that corrupts a chat room's `gate-briefing.json` after binding, then asserts the next `gate prepare` and `gate record` refuse and leave entity bytes unchanged. Independent baseline that fails today: the same corruption passes, because `io.go:204` skips a request-less binding. Falsifying edit: restore the `request-digest == ""` skip. The refusal assertion reds.

**AC-4 - `gate withdraw` retires an open chat attempt, and the retired attempt validates.**
Verified by: a test that prepares a chat room, withdraws it, and asserts exit 0, a `withdrawal` block, and a clean `Validate`. Independent baseline that fails today: the same withdraw exits 1 with `current attempt is not request-backed`. Falsifying edit: restore the request-backed precondition. The exit-0 assertion reds.

**AC-5 - Archived rooms of every historical shape still read.**
Verified by: `spacedock status --workflow-dir docs/dev --validate` over the real state checkout, whose output must equal the pre-change binary's output byte for byte. A retained fixture pins one legacy `briefing.json` room and one opaque `subspace-room:` ref. Falsifying edit: make `preparedRoomBinding` match any directory. The legacy fixture reds.

**AC-6 (MEANS, serving AC-1) - The gate spec and the command reference document the chat default and the provider opt-in.**
Verified by: the doc diff above, reviewed at validation against ASD-STE100 per the workflow README.

## Test plan

`internal/gates` is a status-mutation and guard surface, so the detached adversarial audit applies before merge.

- **Unit, `internal/gates`:** both prepare shapes; the `--provider` golden bytes; `preparedRoomBinding` over the three binding shapes; replay and rebind protection on a chat room; withdraw on a chat room; retained-authority refusal on a corrupted chat Briefing. Cost: low. These reuse `prepareFixture`.
- **Rewrite, not repair:** the 11 named failing test functions. `TestWithdrawRefusalsLeaveEntityRoomAndLockBytesClean/chat-only_request-less_attempt` and `TestPrepareCreatesOneTwoFileRecorderRoomForFolderAndFlatEntities` assert the inverted behavior, so they must state the new contract.
- **Behavior, `internal/cli`:** the CLI journey for AC-1, and the `--provider` usage surface. Cost: medium, because the existing gate CLI tests are slow.
- **Back-compat fixture:** one retained legacy `briefing.json` room and one opaque ref, for AC-5.
- **Cross-repo, recorded not committed:** the q0 preflight run for AC-2. It needs a `spacedock-subspace` checkout, so it is validation evidence, not a committed test.
- No live workflow run is needed. The change touches no dispatch path, host adapter, or journey grader.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Stage Report: ideation

- DONE: Design lands the presentation-channel selection at prepare: the one-file chat default on the existing empty-request-digest path, the provider opt-in that mints today's exact request bytes, and back-compat reads for every archived room — each mechanism naming the value AC it serves, the simplest alternative, and why it is insufficient.
  Six mechanisms in "Proposed approach", each with value AC, simplest alternative, and why it fails. DEVIATION from the seed: the chat default does NOT ride the existing empty-request-digest path. That path means "room-ref names a Briefing file", so a request-less room loses replay, rebind protection, retained-authority validation, and withdraw. Mechanism 2 (`preparedRoomBinding`) replaces the predicate at six sites.
- DONE: Riskiest cross-repo mechanism exercised first and recorded: q0's subspace-tui preflight (spacedock-subspace repo, subspace-r-scaffolded-gate-room journey) run against rooms prepared in BOTH shapes by the modified design's fixture, or an auditable `no spike needed` line with the proven mechanisms; the mediaType question resolved as a recorded cross-repo agreement with the q0 Result contract or an explicit recorded deferral.
  Eight recorded runs in "Risk evidence". q0's `internal/gateroom.Prepare` (75ef1a2) accepts the patched `--provider` room and refuses the chat room by name. mediaType is RESOLVED, not deferred: removing it makes q0 refuse with `briefing: artifact 0 requires id, uri, mediaType, and rev`.
- DONE: Entity-level AC set finalized with a value-measuring AC against the failing-today baseline (every room carries request.json); net-LOC estimate with tolerance; internal/gates named as the status-guard high-stakes surface (detached adversarial audit applies); all comments and doc text follow ASD-STE100 per the workflow README's Prose style section.
  AC-1 measures the room's file set against the 499/499 baseline; AC-3 and AC-4 each carry their own failing-today baseline. Estimate is +200 net across 14 files, tolerance ±60 and ±3, with the production part measured at +62 from the spike diff. Test plan names the audit. Body written to ASD-STE100.
- DONE: Documentation changes proposed as a concrete doc diff.
  Four before/after wordings for `docs/specs/gate-resolution-frontmatter-contract.md` (lines 107, 125, 176), `docs/site/reference/command-reference.md` (line 96), and `docs/site/concepts/gates-and-decisions.md` (line 41).

### Summary

The seed framed this as removing a file. The spike found it is really about a predicate: `request-digest != ""` is what eight code sites use to mean "this attempt has a prepared room", so removing the file silently removes four guards with it. Measured, not read: today's `gate withdraw` exits 1 on a chat room, and `io.go:204` skips retained-authority validation for a request-less binding. The design therefore adds `preparedRoomBinding` as a channel-independent test, moves the withdraw guard from the pure frontmatter validator to the verb that can stat the room, and extends retained-authority validation to chat rooms so the new default is not less-checked than today's.

Two seeded cleanups are cut rather than deferred. The q0 preflight refuses a manifest whose context items lack `type`, and refuses an artifact that lacks `mediaType`, both proven by exercise with a passing control run. Byte compatibility holds: at one pair of fixture commits the baseline binary and the patched `--provider` binary produce identical room bytes, and `status --validate` over the real 600-room state checkout is byte-identical before and after.
