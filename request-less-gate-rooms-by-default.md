---
id: rx3daftacggfmw1pt2febw31
title: Gate rooms hold one canonical Briefing file; request.json is retired
status: validation
source: "Captain gate-format review 2026-08-25: variance scan found 499/499 rooms carry request.json while only ~8 July-era provider results ever consumed one; q0 (subspace-r-scaffolded-gate-room, spacedock-subspace, at validation) finalizes the provider journey the file exists for"
started: 2026-08-25T17:27:46Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-request-less-gate-rooms-by-default
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
            - id: gate-attempt:rx3daftacggfmw1pt2febw31-ideation-2
              briefing:
                id: briefing:rx3daftacggfmw1pt2febw31:ideation:attempt-2:revision-1
                digest: sha256:fff85ca3ae80452056db3eb6b5337b10639978c77400749234f7b893c6aa03c3
                request-digest: sha256:93fa9931addc02449e16c9b3402af2a5272465f66341f603f2f40679d5345aa9
                room-ref: ./request-less-gate-rooms-by-default/review/ideation/briefing-2
              withdrawal:
                by: agent:first-officer
                at: "2026-08-25T22:44:05.495002Z"
                reason: 'Body finalized after binding: the --room ambiguity recorded as resolved, out-of-scope and test-plan edits (7573cd4c0) postdate briefing-2''s frozen snapshot'
            - id: gate-attempt:rx3daftacggfmw1pt2febw31-ideation-3
              briefing:
                id: briefing:rx3daftacggfmw1pt2febw31:ideation:attempt-3:revision-1
                digest: sha256:09b51baaf0dc18b87b81aff1eea38667e211d30871ac86264fcda8bdc269a154
                request-digest: sha256:6ee5a290e31c631a680efc9d019816ecf85b9aa411760be3ca307b9b16a38789
                room-ref: ./request-less-gate-rooms-by-default/review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:rx3daftacggfmw1pt2febw31:ideation:3
                briefing: briefing:rx3daftacggfmw1pt2febw31:ideation:attempt-3:revision-1
                by: person:captain
                at: "2026-08-26T00:29:36.785316Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''approve rx'' — accepts the one-shape reset design on briefing-3; approved baselines: production +70 across 5 files (±25/±2) and proof +330 across 7 files (±80/±2), reported separately'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:rx3daftacggfmw1pt2febw31:validation
          stage: validation
          attempts:
            - id: gate-attempt:rx3daftacggfmw1pt2febw31-validation-1
              briefing:
                id: briefing:rx3daftacggfmw1pt2febw31:validation:attempt-1:revision-1
                digest: sha256:43879ee5ff13cf1f48b9fd23bd5602aa07fbe790c549a89fbe956ed30a1ebe79
                request-digest: sha256:efca06add38403fe2153271b0e79c5a3ebc1eeb916f8a5f113e4de75e1962172
                room-ref: ./request-less-gate-rooms-by-default/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rx3daftacggfmw1pt2febw31:validation:1
                briefing: briefing:rx3daftacggfmw1pt2febw31:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T04:36:54.389919Z"
                decision: approve
                reason: 'Captain chat 2026-08-26: ''approve both'' — accepts validation PASSED; deliver via PR to main; archived-shape finding dispositioned deferred with remedy filed as gate-record-prepared-room-guard'
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:rx3daftacggfmw1pt2febw31-validation-2
              briefing:
                id: briefing:rx3daftacggfmw1pt2febw31:validation:attempt-2:revision-1
                digest: sha256:396c0aa7ecf247738f04260fb783d13ed46e7b7ffd7a78d8f809dfbdba692ed5
                request-digest: sha256:2375f64b92b93dd64ca4b2280761952afcd76f25eb4844f79d722376acd8e793
                room-ref: ./request-less-gate-rooms-by-default/review/validation/briefing-2
---

Every gate room holds two files. Only the subspace provider journey reads `request.json`, and the entity frontmatter already holds every fact that file carries. Make the room one file. Name it `index.json`, and keep every legacy name readable forever.

## Problem

`gate prepare` publishes `gate-briefing.json` and `request.json`. A scan of the dev state checkout on 2026-08-25 counts 499 of 499 prepared rooms with `request.json`. About 8 July-era rooms have a provider result.

The request adds no fact the entity does not already hold. Its `gate`, `attempt`, and Briefing id and digest all repeat the frontmatter binding. Its `actor` and `approver` are the constant `person:captain`, and `operation.go:530` refuses any other value.

Removing the file is not the whole task. Eight code sites read `request-digest != ""` as the test for "this attempt has a prepared room":

- `prepare.go:245` — the frozen open binding cannot be rebound.
- `prepare.go:300` — the entity replay source.
- `prepare.go:351` — exact prepare replay.
- `operation.go:178` — `gate withdraw`.
- `operation.go:551`, `operation.go:562` — the bound Briefing read.
- `application.go:244` — the reviewed-input staleness check.
- `io.go:204` — retained-authority validation.
- `model.go:281` — the frontmatter rule that a withdrawn attempt must retain a `request-digest`.

Two of those guards were measured, not read. Today's `gate withdraw` exits 1 on a request-less room with `current attempt is not request-backed`. And `io.go:204` skips retained-authority validation for a request-less binding, so the shipped binary DURABLY RECORDS a captain approval over a corrupted Briefing: `state=closed ... decision=approve`. Both measurements are from cycle 0 and both fixes salvage forward.

Cycle 1 measured the cost of solving this with a presentation-channel fork: +740 net across 18 files, against an approved +200 across 14. The production column was close to estimate at +126; the fork's cost hid in the proof column, where a doubled test matrix reached +611. One room shape removes the fork, and with it the matrix.

## Proposed approach

One room shape for every presentation channel. There is no channel flag and no command-grammar change.

**1. The room is one file, named `index.json`.** `gate prepare` publishes the canonical Briefing alone and binds no `request-digest`. The name mirrors the entity convention `<slug>/index.md`.

- Value AC: AC-1.
- Simplest alternative: keep both files and let the chat path ignore the request. That is today's behavior, and it is what AC-1 measures against.
- Rejected alternative, with a measured cost: the cycle-0 `--provider` channel fork. It reached +740 across 18 files, and its proof matrix was 84% of the overage. One shape needs neither the flag nor the matrix.

**2. `preparedRoomBinding` replaces `request-digest != ""` at the six runtime sites.** SALVAGE from 312fa3f95; do not re-derive. A request-backed binding always has a prepared room. Otherwise the binding's `room-ref` must name a real directory that does not hold `briefing.json`.

- Value AC: AC-1. Without it, a one-file attempt loses replay, rebind protection, retained-authority validation, and withdraw.
- Simplest alternative: a new frontmatter field such as `channel: chat`. It costs a schema change, a validator rule, and a back-compat read for 600 archived bindings. The room shape already carries the fact.
- Why the test reads `briefing.json` and not `index.json`: a prepared room whose Briefing is deleted must stay a prepared room. An `index.json` test sends that room to the archived read path. There it takes the retained-authority skip, and `gate record` closes over a Briefing that no longer exists. SALVAGE: this is the deleted-Briefing inversion from cycle 0.

**3. Locator resolution reads `index.json`, then `gate-briefing.json`.** `boundBriefingPath` resolves all four binding shapes in one place. SALVAGE from 312fa3f95, extended by one name.

- Value AC: AC-1 and AC-4.
- Simplest alternative: add the name at each read site. That duplicates the order rule, and `application.go` and `operation.go` can then disagree.
- The order is fixed so that two names in one room resolve deterministically. Preparation never writes two, and no room in the corpus holds two.

**4. The withdraw guard moves from the frontmatter validator to the verb.** SALVAGE from 312fa3f95. `Withdraw` requires a prepared room. `model.go:281` drops the rule that a withdrawn attempt must retain a `request-digest`.

- Value AC: AC-4.
- Why the move is necessary: `Validate(doc)` reads frontmatter only. It cannot stat a room, so it cannot tell a one-file room from an archived opaque ref.
- Declared narrowing: the pure validator no longer asserts this property. The mutating verb asserts it.

**5. Retained-authority validation covers every prepared room.** SALVAGE from 312fa3f95. The request checks stay, and they run for a binding that carries a `request-digest`. The Briefing digest and the Git-source checks run for every prepared room.

- Value AC: AC-3.
- Why it is necessary: without it the one shape gets less validation than today's default. That inverts the value of the task.

**6. Room-entry validation follows the binding, not a flag.** A binding with a `request-digest` requires exactly its locator plus `request.json`. A binding without one requires exactly its locator.

- Value AC: AC-1.
- Simplest alternative: accept any one or two files. That drops the "no copied sources, no provider subtree" invariant the spec states.

**7. `request.json` is retired at preparation and readable forever.** Preparation never writes it again. The 499 existing rooms keep theirs, keep their bound `request-digest`, and keep the full request validation they have today. No room migrates.

## Cross-repo interface contract: the q0 preflight

This is a dependency, not a deliverable. The captain routes the q0 change in `spacedock-subspace`. This contract is what the two sides must both meet.

q0's `internal/gateroom.Prepare` amends to:

1. Read the canonical Briefing alone. Try `index.json` first, then `gate-briefing.json`. Refuse with a name list when neither is a contained regular file.
2. Derive identity from the Briefing id's positional components. For `briefing:<entity>:<stage>:attempt-<N>:revision-<M>`, the gate is `gate:<entity>:<stage>` and the attempt is `gate-attempt:<entity-tail>-<stage>-<N>`, where `<entity-tail>` is the last colon-separated segment of `<entity>`.
3. Take the content digest by recomputation: JCS-canonicalize the Briefing bytes, then SHA-256.
4. Set the actor and approver to the constant `person:captain`. That was always the value, and `operation.go:530` refuses any other.
5. Read artifact Git coordinates exactly as today, through the same `reviewv1` source resolver.
6. Drop the request-absence refusal.
7. KEEP the provider-absence refusal. A provider subtree in the room is still an error.

**q0 stops being an authority checkpoint, and this is the load-bearing consequence.** With no request, q0 recomputes the digest instead of comparing it, so it cannot detect a tampered Briefing. Risk evidence run 7 shows it accepting tampered bytes. Authority enforcement therefore sits entirely with the Spacedock recorder, which compares the room against the entity binding under the entity lock. That wall already exists as `validateRetainedAuthority` and it already runs at every mutating verb: `prepare.go:134`, `operation.go:170` and `operation.go:248`, `application.go:84` and `application.go:122`, and `delivery.go:173`. Run 8 shows it refusing the exact bytes q0 accepted.

RESOLVED, and recorded here so the gate does not re-open it. The scope notes for this reset sited the wall at `gate record --room`. That surface does not exist, and three places exclude it: the spec's "Explicitly outside v1" list, `docs/roadmap/durable-decisions/index.md:81`, and `TestGateRecordRejectsProviderRoomBeforeMutation`, which pins the exit-2 refusal. The First Officer confirmed on 2026-08-25 that the phrase named where the wall lives, and that it came from their own scope notes, not from the captain. There is no new command surface and no v1-exclusion reversal. `TestGateRecordRejectsProviderRoomBeforeMutation` stays green, and it now guards this design's boundary.

## Risk evidence

The riskiest claim is that q0 can preflight from the canonical Briefing alone. The spike exercised it first, against `spacedock-subspace` `internal/gateroom` at `75ef1a2`. The ref was read out with `git archive`, so that repo keeps its HEAD, its worktree list, and its tracked bytes. It holds 46 worktrees and other agents' in-flight work.

The amended preflight was written and run against real rooms minted by the real CLI on a split-root fixture with two Git roots:

1. **A new one-file `index.json` room is accepted.** `PREFLIGHT-ACCEPTED gate=gate:task:validation attempt=gate-attempt:task-validation-1 actor=person:captain`, exit 0.
2. **A legacy one-file `gate-briefing.json` room is accepted**, with identical derived identity.
3. **An existing two-file room is accepted**, with identical derived identity. The 499 existing rooms keep working under the amended preflight.
4. **The derived identity equals what today's request-backed path reads.** Running both preflights over the same room returns the same gate, attempt, actor, and Briefing id.
5. **Recomputation equals the bound digest.** The preflight recomputes `sha256:05022451384f66de...`, and the entity frontmatter binds the same value.
6. **The negatives still refuse.** A provider subtree gives `gate room: provider must be absent`. A room with no canonical Briefing gives `gate room: no canonical Briefing (index.json, gate-briefing.json)`. Both exit 1.
7. **A tampered Briefing is ACCEPTED by the amended preflight**, and returns a different digest, `sha256:c9def6143cf8...`. q0 is not the wall.
8. **The Spacedock recorder refuses those same tampered bytes.** `gate record --decision approve` prints `Error: bound canonical Briefing bytes do not match the frozen digest`, exits 1, and leaves entity bytes unchanged.

Cycle 0 proved three more facts that this design keeps. The canonical Briefing is channel-independent. `status --validate` over the real 700-room state checkout is byte-identical before and after. The corpus separation is total. A classification of every archived request-less binding found 118 opaque refs, 33 refs that resolve to nothing, and about 70 directory-shaped rooms. Every one of those directories holds `briefing.json`, and none holds `index.json` or `gate-briefing.json`.

## Documentation changes

Before/after is stated against `main`, because the cycle-0 fork is deleted unbuilt and the net effect must land on `main`.

**`docs/specs/gate-resolution-frontmatter-contract.md` line 28.** Before: `ROOM[("Frozen gate room<br/>request.json and canonical Briefing")]`. After: `ROOM[("Frozen gate room<br/>the canonical Briefing, one file")]`.

**Same file, line 106.** Before: "A prepared provider-neutral room binds `request-digest`, the JCS digest of its `request.json`. Request-less and chat-only attempts may omit it." After: "A prepared room is one file: the canonical Briefing, named `index.json`. Its binding carries no `request-digest`. A binding whose `room-ref` names a directory that does not hold `briefing.json` is a prepared room. Rooms prepared before this change hold `gate-briefing.json` and `request.json`, and bind the JCS digest of that request. They stay readable, and they keep the full request validation. No room migrates."

**Same file, line 124.** Before: "It writes only `gate-briefing.json` and `request.json` at preparation time." After: "It writes only `index.json` at preparation time. Readers resolve the canonical Briefing by name: `index.json` first, then the earlier `gate-briefing.json`."

**Same file, after line 109 (new paragraph).** "The recorder is the authority wall. A presentation channel materializes the room and can recompute what it reads. Only the recorder compares the room against the entity binding, under the entity lock. Every mutating verb runs that comparison before it writes."

**Same file, line 174.** Before: "retires only the selected current-stage open request-backed attempt. Under the shared lock it validates all retained authority and requires the room to contain exactly `gate-briefing.json` and `request.json`." After: "retires only the selected current-stage open prepared attempt. Under the shared lock it validates all retained authority and requires the room to contain exactly the file set its binding implies."

**`docs/site/reference/command-reference.md` line 96.** Before: "Immediately after successful preparation the room contains exactly `gate-briefing.json` and `request.json`, with no copied sources or association." After: "Immediately after successful preparation the room contains exactly `index.json`, the canonical Briefing, with no copied sources or association."

**`docs/site/concepts/gates-and-decisions.md` line 41.** Before: "a two-file recorder-ready room". After: "a one-file recorder-ready room".

**`skills/present-gate/SKILL.md` line 13.** Before: "`request.json`, the canonical Briefing, and recorder validation retain the exact full digest authority." After: "The canonical Briefing and recorder validation retain the exact full digest authority."

## Out of scope

The q0-side code change (the captain routes it in `spacedock-subspace`). A `gate record --room` surface, which stays outside v1 and keeps its refusal test. Migrating or rewriting any existing room. The round recorder's `briefing.json` and `briefing.review.jsonl` shapes. Dropping the per-item `type` field and the artifact `mediaType` field, both refuted in cycle 0 against the shipped q0 loader.

## Expected surface and tolerance

Production and proof are stated separately. Cycle 1's lesson is that the fork's cost hid in the proof column, so a single combined figure hides the signal.

**Production: net +70 LOC, across 5 files. Tolerance ±25 net LOC and ±2 files.** The files are `internal/gates/prepare.go`, `operation.go`, `io.go`, `model.go`, and `application.go`. `internal/cli/cli.go` is NOT in the list: one shape adds no flag, so the command grammar does not change. Cycle 1 measured its own production column at +126 across 6 files with the fork. The flag, the channel plumbing, and the replay-decline check account for the difference.

**Proof: net +330 LOC, across 7 files. Tolerance ±80 net LOC and ±2 files.** Four files change: `internal/gates/prepare_test.go`, `gates_test.go`, `internal/cli/gate_test.go`, and `internal/cli/terminal_consume_test.go`. One more changes in `internal/ensigncycle`. Two are new test files, one for each package. This figure is calibrated against a measurement, not a guess: cycle 1's proof column measured +611 across 8 files with two channels asserted everywhere. One channel roughly halves the new-test volume and keeps the edit volume.

Combined for reference only: net +400 across 12 files. A correction round measures the two columns separately.

Declared semantic changes:

- **Command grammar:** unchanged. No flag is added or removed.
- **Stored format:** a new gate room holds one file named `index.json`, and its binding carries no `request-digest`. Existing rooms and bindings are unchanged.
- **Authority:** `gate withdraw` accepts any prepared room. The frontmatter validator no longer requires a withdrawn attempt to retain a `request-digest`.
- **Runtime behavior:** retained-authority validation now covers request-less prepared rooms. A corrupted Briefing that the shipped binary approves over will be refused. The `status --validate` run over the real state checkout shows no change there.
- **Cross-repo:** a room this version prepares is not readable by the unamended q0 preflight, which refuses it by name. The amended preflight reads both shapes.

## Acceptance criteria

**AC-1 (VALUE) - A gate room is one file, and the gate closes, consumes, and archives as it does today.**
Verified by: a behavior test that drives the built binary through prepare, record, and consume on a split-root fixture. It asserts the prepared room's entry set is exactly `{index.json}`, and that the binding has no `request-digest`. Independent baseline that fails today: the same journey publishes 2 files and binds a `request-digest`, in 499 of 499 rooms scanned. Falsifying edit: make `Prepare` mint the request again. The entry-set assertion reds.

**AC-2 - The canonical Briefing bytes are unchanged, so the amended q0 preflight derives the same identity and digest.**
Verified by: a golden test pinning the `index.json` bytes for fixed inputs, plus the recorded spike runs 1 through 5. Falsifying edit: change any manifest field, the key order, or the indent. The golden reds, and the recomputed digest stops matching the bound digest.

**AC-3 - A prepared room's retained authority is validated whether or not it has a request.**
Verified by: a test that corrupts a one-file room's Briefing after binding, then asserts the next `gate prepare` and `gate record` refuse and leave entity bytes unchanged. Independent baseline that fails today: the shipped binary records `state=closed ... decision=approve` over the same corruption, because `io.go:204` skips a request-less binding. Falsifying edit: restore the `request-digest == ""` skip. The refusal assertion reds.

**AC-4 - `gate withdraw` retires an open one-file attempt, and the retired attempt validates.**
Verified by: a test that prepares a one-file room, withdraws it, and asserts exit 0, a `withdrawal` block, and a clean `Validate`. Independent baseline that fails today: the same withdraw exits 1 with `current attempt is not request-backed`. Falsifying edit: restore the request-backed precondition. The exit-0 assertion reds.

**AC-5 - Every archived and existing room still reads, under all four binding shapes.**
Verified by: `spacedock status --workflow-dir docs/dev --validate` over the real state checkout, whose output must equal the pre-change binary's output byte for byte. A retained fixture pins one `gate-briefing.json` room with its request, one legacy `briefing.json` directory, and one opaque `subspace-room:` ref. Falsifying edit: drop the `briefing.json` test from `preparedRoomBinding`. The legacy fixture reds.

**AC-6 (MEANS, serving AC-1) - The gate spec and the site docs state the new contract: the one-file room, the locator order, the request retirement, and the authority wall.**
Verified by: the doc diff above, reviewed at validation against ASD-STE100 per the workflow README.

## Test plan

`internal/gates` is a status-mutation and guard surface, so the detached adversarial audit applies before merge.

- **Unit, `internal/gates`:** the one-file prepare shape, the `index.json` golden bytes, and `preparedRoomBinding` over the four binding shapes. Also locator resolution order, replay and rebind protection, and withdraw on a one-file room. Also the retained-authority refusal on a corrupted Briefing, and the deleted-Briefing guard. Cost: low. These reuse `prepareFixture`.
- **Edit, not double:** the existing tests that assert the two-file shape state the one-file contract instead. There is one shape, so no test gains a channel axis. The tests that corrupt `request.json` to reach retained-authority validation corrupt the Briefing instead.
- **Retained guard, unchanged:** `TestGateRecordRejectsProviderRoomBeforeMutation` stays green. It refuses `gate record --room`, and that refusal is now this design's boundary. Do not relax it.
- **Behavior, `internal/cli`:** the CLI journey for AC-1. Cost: medium, because the existing gate CLI tests are slow.
- **Back-compat fixture:** one `gate-briefing.json` room with its request, one legacy `briefing.json` directory, and one opaque ref, for AC-5.
- **Cross-repo, recorded not committed:** the q0 spike runs above. They need a `spacedock-subspace` checkout, so they are recorded evidence, not a committed test.
- No live workflow run is needed. The change touches no dispatch path, host adapter, or journey grader.


### Feedback Cycles

- Cycle 1: DESIGN RESET — captain ruling 2026-08-25 ("find a simpler way" → one room shape, `index.json`): implementation 312fa3f95 landed at 18 files/+740 net vs approved 14±3/+200±60; overage was the channel fork's doubled test matrix. Reset to ideation: one file for every channel, no request.json anywhere, q0 amended to preflight from the canonical Briefing alone; predicate work and both hole fixes salvage forward; the fork and its matrix are deleted unbuilt. Spec update (gate-resolution-frontmatter-contract.md) required per captain.
- Cycle 2: one-shape implementation (0c1679864, rebuilt from origin/main at b59c208d8); surface production net -32 vs +70±25 (favorable: the request-minting path and four orphaned functions deleted) and proof net +490 across 10 files vs +330±80 across 7 (room shape hard-coded in four ensigncycle files where ideation counted one, plus the retained legacy-shape fixture); captain accepted both as amended baselines (2026-08-26, 'approve'); validation proceeds against them.
- Cycle 3: REWORK — required live lanes red on PR #762: shared_live_runner_test.go:162 asserts the retired two-file room while the candidate mints the approved one-file index.json room (claude-live and codex-live, identical line); a fifth hard-coded room-shape site behind the live build tag, invisible to go test ./...; fix routed to the implementation worker; lanes re-run green are the delta validation.

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

## Stage Report: implementation

- DONE: The deliverable lands the approved design exactly: the preparedRoomBinding predicate replacing the request-digest test at the six named sites; the one-file chat default; the --provider flag minting byte-identical request bytes (golden-pinned); the withdraw guard moved to the verb that can stat the room; retained-authority validation extended to chat rooms; back-compat reads for archived rooms and legacy shapes; the four approved doc diffs applied.
  Commit 312fa3f95 on `spacedock-ensign/request-less-gate-rooms-by-default`. Six sites: `prepare.go:256/311/366`, `operation.go:180/553`, `io.go:210`; `boundBriefingPath` (`operation.go:576`) serves `operation.go:600` and `application.go:248`; `model.go:281` rule dropped. REFINEMENT to mechanism 2, declared for the gate: the predicate tests for the archived `briefing.json` name, not for the reserved `gate-briefing.json` the design names. Rationale below.
- DONE: Proofs recorded red-then-green wherever the AC names a failing-today baseline (AC-1 room entry set, AC-3 corrupted-chat-room refusal, AC-4 withdraw exit-0); status --validate over the real state checkout byte-identical to the pre-change binary; go test ./... and -race green; gofmt clean; all comments and doc text strict ASD-STE100.
  All three reds were produced by the real pre-change binary built at `95f877cd6`, not by a simulated edit. AC-1: the new CLI journey test reds on the unmodified tree with `chat room entries=[gate-briefing.json request.json]`. AC-4: baseline `gate withdraw` on a chat room prints `Error: current attempt is not request-backed`, exit 1; patched exits 0 and writes the withdrawal block. AC-3: on a chat room with a tampered `gate-briefing.json`, the baseline binary DURABLY RECORDED the approval (`state=closed ... decision=approve`); patched refuses both `gate record` and `gate prepare` with `bound canonical Briefing bytes do not match the frozen digest`, entity byte-clean. AC-5: `status --validate` over the real state checkout (700 rooms) is byte-identical, sha256 `2527199b19f0...`, both exit 0. AC-2: `diff -r` of the baseline default room and the patched `--provider` room at the same fixture commits reports no difference, and the bound frontmatter matches. `go test ./...` and `go test ./... -race` both exit 0. gofmt clean on all 18 changed files; `internal/release/runtime_live_evidence_workflow_test.go` is unformatted but was already unformatted at `95f877cd6` and is untouched here. Added comments and doc prose scanned mechanically: no sentence over 25 words, no banned modal.
- SKIPPED: Surface within the approved net +200 across 14 files (tolerance ±60/±3) with only the declared semantic changes; deliverable committed on the worktree branch with a summary of what was produced and where.
  Superseded by the captain's Cycle-1 design reset; surface obligation not waived — the design it measured is retired. This item read FAILED when the surface was first reported, and the ruling in `### Feedback Cycles` Cycle 1 is what changed it. The measured numbers stand as recorded. The commit half is done (312fa3f95). Measured `git diff --numstat` against the merge-base: 18 files, +1127/-387, net +740. Files: 18 vs the 11-17 band, one over. Net LOC: +740 vs the 140-260 band, +480 over. Production is close to the estimate at 6 files / net +126 (spike measured +62); docs and skills are 4 files / net +3. The overage is tests: 8 files / net +611 against roughly +138 implied for 5 test files. Two causes. First, the ACs require every prepare/withdraw/authority test to run in both channels, which roughly doubles five existing test functions and adds five new ones. Second, the design's blast-radius scan ran `go test ./internal/gates/... ./internal/cli/...` and so never saw `internal/ensigncycle`, which holds three more tests asserting the two-file room; those are the two files above the file band.

### Summary

Chat is now the default gate channel and `--provider` is the opt-in. The predicate work was the real content, as ideation found: `request-digest != ""` was carrying four guards, and all six runtime sites now use `preparedRoomBinding` instead.

One refinement to mechanism 2 needs a ruling. The design defines the predicate as "a directory that holds `gate-briefing.json`". Implemented literally, a chat room whose Briefing is deleted stops being a prepared room, takes the skip archived rooms get, and `gate record` then closes the captain's approval over a Briefing that no longer exists - the exact inversion mechanism 5 exists to prevent. The predicate therefore tests for the archived `briefing.json` name instead. A classification of all 700 rooms in the real state checkout shows the separation is total: every archived directory-shaped room holds `briefing.json`, and none holds `gate-briefing.json`. AC-5's falsifying edit still falsifies, because dropping the `briefing.json` test is what makes the predicate match any directory.

Two holes the design did not name were found and closed. An exact replay was channel-blind, so `gate prepare --provider` over an open chat room silently replayed the chat room and returned success with no `request.json`; replay now declines on a channel change and the frozen open binding refuses the rebind. The deleted-chat-Briefing hole above is the second.

### Addendum: AC-2 cross-repo leg exercised

The first pass inherited q0 acceptance from the ideation spike and proved only byte identity. The q0 preflight is now run directly against this implementation's rooms, so AC-2 rests on exercise and not on inheritance.

`internal/gateroom.Prepare` was read out of `spacedock-subspace` at `75ef1a2` with `git archive` into a scratch tree, so that repo keeps its HEAD, its worktree list, and its tracked bytes. Three rooms were prepared from one split-root fixture at one pair of commits:

- Patched binary, `--provider`: `PREFLIGHT-ACCEPTED gate=gate:task:validation attempt=gate-attempt:task-validation-1 actor=person:captain briefing=briefing:task:validation:attempt-1:revision-1`, exit 0. The scratch tree holds the materialized Briefing plus both Git blobs, `git-root:/main/4c3c7b9d.../gate-review.md` and `git-root:/state/d34ff644.../reference.json`.
- Baseline binary at `95f877cd6`, today's default, as the control: identical accepted line and identical materialized set, exit 0. `diff -r` of the two rooms reports no difference.
- Patched binary, chat default: `PREFLIGHT-REFUSED: gate room: request.json: ... no such file or directory`, exit 1. q0 refuses a chat room by name, before any host command, and does not degrade silently.

That q0 accepts the patched provider room and refuses the chat room reproduces spike results 2 and 3 against the delivered code. The falsifying edit is unchanged: alter any bound field, the key order, or the indent, and the golden reds and the q0 digest check reds.
## Stage Report: ideation (cycle 2)

- DONE: The reshaped design lands the captain's one-shape direction exactly: index.json canonical for new rooms, no request.json in any channel, legacy names readable with the predicate discriminator intact, the q0 interface contract stated byte-precisely as an explicit cross-repo dependency, the salvage list from 312fa3f95 named, and the fork deleted unbuilt.
  Seven mechanisms in "Proposed approach", each naming its value AC, the simplest alternative, and why that alternative fails. Mechanisms 2, 3, 4, and 5 are marked SALVAGE from 312fa3f95. The `--provider` flag, the channel plumbing, the replay-decline check, and the doubled matrix are gone; "Declared semantic changes" now records command grammar as unchanged, and `internal/cli/cli.go` leaves the production file list. The predicate discriminator is unchanged and still tests for `briefing.json`, which is what keeps a deleted Briefing loud; the corpus classification in "Risk evidence" shows the separation is total across 700 rooms. ONE AMBIGUITY was raised instead of guessed, and it is now RESOLVED. The scope notes sited the authority wall at `gate record --room`, which does not exist and is excluded from v1 in the spec, in `docs/roadmap/durable-decisions/index.md:81`, and by `TestGateRecordRejectsProviderRoomBeforeMutation`. The First Officer confirmed the phrase named where the wall lives, and that the imprecision was theirs. The design adds no command surface, and the refusal test stays green as a boundary guard.
- DONE: The spec diff for docs/specs/gate-resolution-frontmatter-contract.md is concrete before/after (captain-required), plus the site docs; all prose strict ASD-STE100.
  Eight before/after diffs in "Documentation changes": five in the gate spec (the mermaid room node, the request-digest paragraph, the preparation-writes sentence, a NEW recorder-as-authority-wall paragraph, and the withdraw room-contents sentence), plus `docs/site/reference/command-reference.md`, `docs/site/concepts/gates-and-decisions.md`, and `skills/present-gate/SKILL.md`. Before/after is quoted against `main`, not against the retired fork branch, because the fork is deleted unbuilt and the net effect must land on `main`. New body prose was scanned mechanically for the 25-word cap and the banned modals; the only `may` left is inside a quoted "before" string.
- DONE: AC set and test plan rebuilt single-channel with failing-today baselines; the estimate states production and proof surface separately, each with tolerance.
  Six ACs, no channel axis anywhere. AC-1 measures the room entry set against the 499/499 baseline. AC-3 and AC-4 each carry their own measured failing-today baseline, both re-measured this cycle against the real pre-change binary. AC-5 pins all four binding shapes. The estimate is split as the captain required: production net +70 across 5 files (±25/±2), proof net +330 across 7 files (±80/±2). The proof figure is calibrated against cycle 1's MEASURED +611 across 8 files, not guessed.

### Summary

The riskiest claim was that q0 can preflight from the canonical Briefing alone, so it was exercised first. The amended preflight was written against `spacedock-subspace` `internal/gateroom` at `75ef1a2` and run on real rooms minted by the real CLI. It accepts all three live shapes — a new `index.json` room, a legacy `gate-briefing.json` room, and an existing two-file room — and derives a gate, attempt, actor, and Briefing id identical to what today's request-backed path reads out of `request.json`. Its recomputed digest equals the digest the entity binds. The provider-absence and missing-Briefing refusals still fire by name. The ref was read out with `git archive`, so that repo keeps its HEAD, its worktree list, and its tracked bytes.

Run 7 is the finding that shaped the design. A tampered Briefing is ACCEPTED by the amended preflight, which returns a different digest, because recomputation cannot detect tampering. Run 8 shows the Spacedock recorder refusing those same bytes with `bound canonical Briefing bytes do not match the frozen digest`, exit 1, entity byte-clean. Removing `request.json` therefore moves a real check off the presentation side, and the recorder-as-authority-wall statement is load-bearing rather than decorative. The design states it, the spec diff adds a paragraph for it, and AC-3 tests it.

Falsifiability of the spike, in one line per claim: the identity derivation reds if the Briefing id grammar changes, because gate and attempt are parsed from its positional components; the digest claim reds if any manifest field, key order, or indent changes, because the recomputed value stops matching the bound value; the authority-wall claim reds if the `request-digest == ""` skip returns to `io.go`, because run 8 then records the approval instead of refusing it.

## Stage Report: implementation (cycle 2)

- DONE: FIRST ACTION: reset your worktree branch to origin/main's tip — the retired fork commit 312fa3f95 leaves the branch history (it was never pushed; reuse its hunks as raw material where the salvage list applies). Then land the approved one-shape design exactly: index.json canonical for new rooms, the preparedRoomBinding predicate at the six sites, retained-authority validation extended to one-file rooms (keep the tampered-Briefing red), the deleted-Briefing fix, legacy names readable with the discriminator intact, no --provider flag, no new command surface, and TestGateRecordRejectsProviderRoomBeforeMutation untouched as the design's boundary guard.
  Branch reset to `origin/main` at `b59c208d8`; the fork was confirmed unpushed first, and its patch was saved to scratch as raw material rather than kept as history. Deliverable is 0c1679864. Six predicate sites: `prepare.go:239/294/349`, `operation.go:180/546`, `io.go:210`. `boundBriefingPath` (`operation.go:565`) serves `operation.go:592` and `application.go:248`. No flag and no command-grammar change; `internal/cli/cli.go` is untouched. `TestGateRecordRejectsProviderRoomBeforeMutation` is untouched and green — the diff of that test is zero lines. Removing the request path orphaned `validateGateRoomRequest` and its three request-locating helpers, which are deleted; no test referenced them.
- DONE: Proofs per the approved plan, red-then-green on each failing-today baseline (room entry set, corrupted-room refusal, withdraw exit-0), the index.json golden, all four binding shapes, and status --validate byte-identity over the real state checkout; go test ./... and -race green; gofmt clean; all prose strict ASD-STE100; the eight approved doc diffs applied including the five spec blocks.
  All three reds were produced by the real pre-change binary, not a simulated edit. Entry set: pre-change publishes `gate-briefing.json request.json` and binds a request-digest; this build publishes `index.json` and binds none. Withdraw: pre-change prints `current attempt is not request-backed`, exit 1; this build exits 0 with the withdrawal block. Corrupted room: the pre-change binary DURABLY RECORDED the approval (`state=closed ... decision=approve`, 1 resolution written); this build refuses with `bound canonical Briefing bytes do not match the frozen digest`, exit 1, entity byte-clean. Golden: `TestPreparedBriefingBytesArePinned` pins every byte except the two Git coordinates, so a changed field, key order, or indent reds. Four binding shapes plus three negatives: `TestBindingShapesClassifyAndResolve`, which reds if the `briefing.json` test is dropped from the predicate. AC-5: `status --validate` over the real 703-room state checkout is byte-identical, sha256 `2527199b19f0...`, both exit 0. `go test ./...` exit 0 and `go test ./... -race` exit 0. gofmt clean on all 19 changed files. Prose scanned mechanically: no sentence over 25 words, no banned modal. Eight doc diffs applied, the five spec blocks at lines 28, 106, 114 (the new authority-wall paragraph), 133, and 184.
- DONE: Surface reported in BOTH approved columns separately via numstat: production against +70 across 5 files (±25/±2) and proof against +330 across 7 files (±80/±2); deliverable committed on the worktree branch with a summary of what was produced and where.
  Reported as required. This item read FAILED when both columns first measured outside their approved bands; the captain accepted both figures as amended baselines, recorded in `### Feedback Cycles` Cycle 2 (2026-08-26), and that ruling is what changed it. The measured columns and the original-band comparisons stand as recorded. The commit half is done (0c1679864). PRODUCTION: 5 files (+176/-208), net -32, against +70 ±25. Files land exactly on the approved 5; net is 102 below estimate and 77 under the band floor. The miss is favorable and one-directional: the estimate priced the predicate's additions but not the deletions the one-shape design forces, namely the whole request-minting block and the four now-orphaned request functions, about 130 deleted lines. PROOF: 10 files (+621/-131), net +490, against +330 ±80. Files are one over the 5-9 band, net is 80 over the 410 ceiling. Cause, measured not guessed: the room shape is hard-coded in four `internal/ensigncycle` files, and the ideation estimate named one; and AC-5's retained two-file fixture is new work this cycle that no earlier cycle had to build. DOCS/SKILLS: 4 files, net +8. COMBINED: 19 files, net +466.

### Summary

The reset design landed as approved. The room is one file named `index.json`, nothing migrates, and the salvaged predicate work carried over unchanged in substance: `preparedRoomBinding` at six sites, the withdraw guard on the verb, retained-authority validation over request-less rooms, and the deleted-Briefing fix. Cutting the channel fork also cut the production column below its estimate, because the request-minting path and four functions that existed only to serve it are gone rather than made conditional.

The cross-repo contract was re-exercised against a real room this build minted, not the renamed fixture the ideation spike used. The amended q0 preflight accepts the `index.json` room and recomputes `sha256:b1add0b89266...`, exactly the digest the entity binds. The shipped, unamended preflight refuses the same room by name, which is the declared cross-repo semantic change. The recorder remains the authority wall, and that is load-bearing rather than decorative: the amended preflight accepts a tampered Briefing, and the recorder refuses those same bytes.

One process note, not a defect. An early `go test ./... -race` run reported `panic: test timed out after 10m0s` on `internal/ensigncycle`. That was contention from my own concurrent test runs, not the code. Run alone at the default timeout the package passes `-race` in 202s, and this change adds no test functions there.

Falsifiability, one line per claim: the entry-set claim reds if `Prepare` mints a request again; the authority claim reds if the `request-digest == ""` skip returns to `io.go`, because the recorder then records the approval instead of refusing it; the back-compat claim reds if the `briefing.json` test leaves `preparedRoomBinding`, because every archived room directory then has to supply a Briefing name it does not hold.

## Stage Report: validation

- DONE: Every AC (AC-1 through AC-6) reproduced first-hand on candidate 0c1679864: the room entry-set red-then-green, the index.json golden, the corrupted-room refusal, withdraw exit-0, all four binding shapes, and status --validate byte-identity over the real state checkout; go test ./... and -race fresh (-count=1, read the run's own exit code, never through a pipe); gofmt clean; STE spot-check over added prose; the eight approved doc diffs applied including the five gate-spec blocks.
  All reds came from the real baseline binary built at merge-base b59c208d8 on my own split-root fixtures, not from the implementer's recordings. AC-1: baseline publishes `gate-briefing.json request.json` and binds a request-digest; candidate publishes exactly `index.json`, binds none, and record --consume closes, consumes, advances to implementation, and leaves the archived room one file with `status --validate` exit 0. AC-3: baseline exit 0 and a durable `decision: approve` over a tampered index.json; candidate refuses record AND successor prepare with `bound canonical Briefing bytes do not match the frozen digest`, exit 1, entity bytes unchanged. AC-4: baseline `current attempt is not request-backed` exit 1; candidate exit 0 with a withdrawal block and clean validate. AC-2: golden red proven by mutation (a one-space indent change reds TestPreparedBriefingBytesArePinned) and the preflight recomputation below. AC-5: baseline and candidate `status --validate` over the real state checkout are byte-identical (sha256 ec1a9ff194b7..., both exit 0, VALID, on today's corpus — the corpus moved since implementation, so the hash differs from the recorded 2527199b19f0... while the identity claim holds); all four binding shapes pinned by TestBindingShapesClassifyAndResolve, and the candidate closed/consumed a baseline-minted two-file room and refused its corrupted retained request. AC-6: all eight doc diffs verified against the approved quotes (the SKILL.md line also swaps banned `may` for `can`; the spec paragraph keeps one adapted fails-before-mutation sentence). Fresh `go test ./... -count=1` PLAIN_EXIT=0 and `go test ./... -race -count=1` RACE_EXIT=0, exit codes read from the runs' own `$?` into log files. `gofmt -l ./cmd ./internal` names only `internal/release/runtime_live_evidence_workflow_test.go`, untouched here and already unformatted at b59c208d8. Mechanical STE scan of added prose: no banned modal, no sentence over 25 words.
- DONE: Detached adversarial audit on a THROWAWAY checkout (internal/gates is the status-guard high-stakes surface): start with the implementer's own named targets — the preparedRoomBinding discriminator (it reads the archived briefing.json name; a deleted Briefing must stay LOUD) and the recorder-as-authority-wall (the amended q0 preflight recomputes rather than compares, so replay the tampered-Briefing sequence and confirm the recorder refuses what the preflight accepts) — then the AC-provenance trigger; findings enter the Review-finding disposition checkpoint.
  Audit ran on a `git archive 0c1679864` throwaway with its own binary; candidate worktree bytes and HEAD untouched. Deleted Briefing stays LOUD on record, withdraw, and successor prepare (all exit 1, entity byte-clean). Tampered-Briefing replay: my independently written amended preflight (from the entity's contract, not the spike code) ACCEPTS the tampered room with a recomputed different digest; the recorder refuses the same bytes, exit 1 — the wall sits where the design says. Falsifying-edit mutations on the throwaway prove the guards can fail: restoring the io.go request-digest skip reds both corrupt-Briefing refusal tests; dropping the briefing.json test reds the archived-room case; restored tree clean. AC-provenance: no self-referential evidence in any AC. One real finding (below): planting a `briefing.json` file into a prepared room, or swapping the room dir for a symlink, flips the discriminator to archived, and `gate record` then closes with retained-authority SKIPPED — combined with a tampered index.json it durably approves over tampered bytes (probes 8/9, exit 0). Entered in Review-finding disposition below as a proposed deferred risk.
- DONE: Cross-repo and surface: re-verify the amended q0 preflight contract against a room this build minted (reproduce the recorded runs, do not trust them); semantic adversarial pass scaled to the diff; surface confirmed against the captain-amended baselines (production net -32, proof net +490 across 10 files); recommend PASSED or REJECTED with material findings, deferred risks, and polish listed separately.
  Amended preflight rewritten from the contract text against `spacedock-subspace` `internal/gateroom` at 75ef1a2 (read out with `git archive`; that repo untouched). Over a room this build minted: ACCEPTED with gate/attempt/actor/briefing identical to the entity binding and recomputed digest exactly equal to the bound `sha256:0937ab3f...`; the shipped preflight refuses the same room by name (the declared cross-repo change); a legacy `gate-briefing.json` room accepts with identical identity; a baseline-minted two-file room gives identical identity under shipped and amended; provider-subtree and no-Briefing rooms refuse with the named errors. Semantic pass covered identity, cardinality (entry sets, two-locator refusal), order (index.json wins deterministically over a divergent second name), bytes, attribution, authority, and terminal states; no multiplicative work added (two Lstat calls per attempt). Surface verified by numstat against merge-base: production 5 files +176/-208 net -32, proof 10 files +621/-131 net +490, docs/skills 4 files net +8 — exactly the captain-amended baselines. RECOMMEND PASSED; no material finding; one deferred risk and two polish items below.

### Summary

All six ACs reproduced first-hand with red sides driven by the real merge-base binary and the cross-repo leg re-run through an independently rewritten amended preflight; both test suites fresh-pass and the surface matches the amended baselines exactly. The detached audit confirms the two named targets and adds one finding: the archived-classification skip can be reached adversarially from a prepared room (plant `briefing.json` or symlink the room dir), silently restoring today's approve-over-corruption hole; the trigger needs hand mutation of the state checkout, is absent from the 700-room corpus, and is strictly narrower than the shipped binary's behavior, so it is proposed as a deferred risk, not material. Recommendation: PASSED.

### Review-finding disposition entries (validation, cycle 2)

- Finding (reviewer entry, proposed Deferred risk): an archived-shape perturbation of a prepared one-file room drops the recorder's digest wall. Evidence: audit probes on throwaway checkout — plant `briefing.json` beside index.json, or replace the room dir with a symlink, then tamper index.json; `gate record --decision approve` exits 0 and durably writes the approval. Four fields: (1) released user and normal workflow — FO gate ceremony over tool-managed rooms; no supported command writes into `review/<stage>/briefing-N/` (rounds retain to `round-<cycle>`, round.go:49; prepare writes exactly one file); (2) observable harm — durable captain approval over tampered Briefing bytes, the AC-3 inversion; (3) authority — contract[docs/specs/gate-resolution-frontmatter-contract.md] "Every mutating verb runs that comparison before it writes": the archived-classified skip breaks that sentence's promise when the archived shape was minted by tampering; (4) trigger evidence — two hand mutations inside the state checkout, zero occurrences in the corpus classification, and an actor with that write access can already rewrite the entity binding itself. Promote-to-material condition: any feature that writes additional files into `briefing-N/` rooms, or an observed prepared room holding `briefing.json`. Possible narrow remedy for FO consideration (not performed): `gate record` could require `preparedRoomBinding` for the open attempt it closes, as `Withdraw` already does.
- Polish: the new spec sentence "Every mutating verb runs that comparison before it writes" overstates — archived-classified bindings (opaque refs, legacy rooms) are skipped, as today. Suggest scoping the sentence to prepared and request-backed rooms.
- Polish (pre-existing scope, unchanged by this diff): `status --validate` reports VALID over a prepared room whose Briefing file is deleted while the room dir resolves; the #739 warning covers unresolvable refs only, identically in the baseline. Mutating verbs stay loud, so nothing silently proceeds.

## Stage Report: implementation (cycle 3)

- DONE: Every live-tagged room-shape assertion matches the approved one-file contract — shared_live_runner_test.go:162 fixed and the full live-tag sweep recorded with the grep evidence (all sites listed, fixed in one commit).
  CORRECTION TO THE ASSIGNMENT'S PREMISE, stated so the gate reads the right file: the assertion is NOT in `shared_live_runner_test.go`. Line 162 there is the `TestLiveCommonWithdrawnGateRecovery` registration. The assertion lives at `claude_live_runner_test.go:334`, and CI attributed it to the shared file because `liveJourney` (`shared_live_runner_test.go:28`) calls `t.Helper()`, which walks the reported line up to the test function. It now requires exactly one regular entry named `index.json`, and asserts the successor binding carries no `request-digest`. SWEEP, all sites listed. Build-constrained files repo-wide: 19 total, 15 with `//go:build live` (all in `internal/ensigncycle`) and 4 with unix constraints in `internal/cli` that hold no room reference. Of the 15 live files, exactly one asserted a gate-room shape, and it is the one fixed. The FO's `grep -rn "gate-briefing.json\|request.json" internal/ensigncycle/` now returns only `testdata/*.stream.jsonl` hits, which are frozen August live-session recordings and must not be edited. A filename grep alone would have missed this site: it counted entries (`len(entries) != 2`) instead of naming a file, so the sweep also ran `grep -rn "os.ReadDir\|len(entries)" internal/ensigncycle/`. That second grep surfaced `shared_round_recording_test.go:472-476`, which is correctly unchanged: it pins the round recorder's `round-1` room (`briefing.json` plus `briefing.review.jsonl`), which the design puts out of scope. Remaining repo-wide `request.json` references were each classified and left: retained-request production paths (`io.go`, `operation.go`, `prepare.go`), the AC-5 retained two-file test, hand-built legacy fixtures in `merge_guard_test.go`, `native_mutation_test.go`, `state_commit_test.go`, `gate_test.go`, and `gate_application_warning_test.go`, the unrelated `dispatch build < request.json` stdin file, and my own legacy-compat doc sentences.
- DONE: The live tag compiles clean (go vet -tags live; build-only test invocation) with no scenario executed locally.
  `go vet -tags live ./internal/ensigncycle/` is silent. `go test -tags live -run TestNothing ./internal/ensigncycle/` reports `ok ... [no tests to run]` in 0.46s, so no scenario ran. `go vet ./...` on the default build is also silent, and `gofmt -l` on the changed file is empty.
- DONE: One commit atop 0c1679864, unpushed; stage report round appended; state committed.
  `acd06d4ce`, one file, +5/-2, parent `0c1679864`. `git branch -r --contains acd06d4ce` is empty, so it is unpushed; the CI re-run of the required lanes is the delta validation.

### Summary

The finding is correct and the miss is mine. My cycle-2 sweep grepped for room filenames and ran `go test ./...`, and neither can see a `//go:build live` file: the tag excludes it from the default build, so the stale assertion compiled and ran only in CI. My cycle-2 report named four `internal/ensigncycle` files as the map; this was a fifth, and it sat behind the tag.

Two lessons worth carrying, because both are general and neither is specific to this task. A room-shape assertion does not have to name a file — this one counted directory entries — so a filename grep is not a sufficient sweep for a stored-format change. And build tags hide whole files from `go test ./...`, so any format change needs an explicit tagged-build compile check; `go vet -tags live` plus a `-run TestNothing` invocation costs under a second and would have caught this in cycle 2.

Falsifiability: revert `Prepare` to mint a request and the new assertion reds on the entry name and count; drop the `request-digest` check and a room that regains a bound request passes the entry check while still violating the one-file binding contract.
