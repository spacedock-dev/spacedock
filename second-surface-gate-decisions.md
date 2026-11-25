---
id: 301s961pr7b2cnywnkz88jec
title: A gate that a second surface can read and decide safely
status: ideation
source: "MetaFO friction list 2026-11-24, sent on the word of CL, /tmp/spacedock-frictions-2026-11-24.md, group A items 1-5"
group: A-second-surface-gate
started: 2026-11-25T01:34:10Z
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
        - id: gate:301s961pr7b2cnywnkz88jec:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:301s961pr7b2cnywnkz88jec-backlog-1
              briefing:
                id: briefing:301s961pr7b2cnywnkz88jec:backlog:attempt-1:revision-1
                digest: sha256:52a3c58a278bc4f991b7cffcae827ada3f87f1b68dfcc51e81e09a86a3d605a7
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:301s961pr7b2cnywnkz88jec:backlog:1
                briefing: briefing:301s961pr7b2cnywnkz88jec:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-25T01:33:48.095174Z"
                decision: approve
                reason: 'MetaFO decision for CL: approve the backlog seed for briefing:301s961pr7b2cnywnkz88jec:backlog:attempt-1:revision-1; dispatch ideation, whose design returns to CL at the ideation gate.'
                conn:
                    quote: 'MetaFO decides your backlog gate for the task "A gate that a second surface can read and decide safely", under the approval rule that CL gave MetaFO for a gate that only moves a task to its next stage: approve, for the briefing briefing:301s961pr7b2cnywnkz88jec:backlog:attempt-1:revision-1.'
                    source: MetaFO message to the spacedock-v1 manager, 2026-11-24, under CL's approval rule for gates that only move a task to its next stage
              application:
                target-stage: ideation
                state: consumed
---

A captain can decide a gate from a pane, a Spacedock Review page or another agent's message, and the decision lands on the briefing the captain read, with its full text and its origin recorded.

## Problem

A second surface is anything other than the first officer's own chat. Checked against origin/main `a85a26a35`:

1. A chat decision cannot name the briefing it answers. `gate record` closes whatever attempt is current (`internal/gates/operation.go:25` `RecordInput` has no expected briefing or digest; `:248` `currentStageAttempt`). A late decision can land on a briefing the captain did not read. The same fault was seen on the Review page: a decision went to version 1 while the agent waited on version 2.
2. `gate read --json` returns only the `gates:` block (`internal/status/gate_history.go:13-52`). It does not give the entity path and status, each attempt's state (`attemptState`, `internal/gates/model.go:403`), which attempt is current, or the resolved room (`ResolveRoomRef`, `internal/gates/prepare.go:786`). A pane parses the entity file itself.
3. The briefing the captain reads is not stored with the gate. `--summary` has no length limit (`internal/gates/prepare.go:85-89`: valid UTF-8 and nonblank only), so a short room is a first-officer choice. But the rendered gate review (the `present-gate` text) exists only in chat; nothing stores it in the room.
4. A captain Resolution has no field for where it came from: `by`, `at`, `decision`, `reason` (`internal/gates/model.go:69-78`), and conn fields are refused for `person:captain` (`internal/gates/operation.go:246`).
5. When another agent starts a first officer's turn, the gate presentation did not reach the captain's chat. A manager-prompt line now corrects it in one project; the product question stays open.

## Proposed approach

One rule ties the five items together: **a decision names the Briefing it answers, and the Briefing holds the review the captain read.** The binary gets one required flag, one optional flag, one optional stored field and a wider `gate read` result. The first-officer skills get a new order (compose the review, then prepare) and one presentation rule.

### 1. `gate record` names the Briefing it answers

Grammar (chat decisions; `--round` is unchanged):

```text
spacedock gate record <entity> --decision approve|revise|hold --actor ID --briefing-id BRIEFING_ID
    [--reason TEXT] [--via TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]
```

- `--briefing-id` is required on every `--decision` record, for both actors. A missing flag is a usage error (exit 2) before any read or write.
- Under the entity lock, the recorder compares the value with the `briefing.id` of the current attempt (the attempt `currentStageAttempt` returns today). A mismatch exits 1 and writes nothing: the entity bytes, the state commit and the sync stay unchanged. With `--consume`, the refusal happens before the close, so nothing is consumed.
- The refusal names the Briefing the decision named and that attempt's state (or says no attempt of this gate has that id), and names the current attempt's Briefing and state: `Error: decision names briefing:task:validation:attempt-1:revision-1 (withdrawn); the current attempt is briefing:task:validation:attempt-2:revision-1 (open); nothing was recorded`.
- A match then runs today's checks unchanged (open state, verifiable digest, actor and conn rules).
- `--briefing-id` and `--via` are refused with `--round` (exit 2), like the conn flags.

The input is the id only. A `--digest` flag was considered and is not added, because it is not necessary. A prepared open binding is frozen (`internal/gates/prepare.go:247-250`; spike step 3 below), and every new attempt mints a new id (`attempt-N`). So in a prepared room one id names one set of bytes. The surfaces also carry the id: MetaFO's relayed decision for this task's backlog gate quoted the id and not the digest (see the `conn.quote` in the frontmatter). Agents copy a 71-character hex digest badly; `reliable-exact-digest-in-gate-review` records one such failure. The limit: an open attempt bound before one-file rooms can still be rebound under the same id, so an id check cannot see that rebind. No such open attempt is known.

### 2. `gate read --json` gives what a second surface needs, in one call

The output keeps `version` and `records` exactly as stored and adds two keys:

```json
{
  "version": 1,
  "records": [ "...stored gate records, unchanged..." ],
  "entity": { "path": "/abs/.../task.md", "status": "validation" },
  "current": {
    "attempt": "gate-attempt:task-validation-2",
    "state": "open",
    "briefing": { "id": "briefing:task:validation:attempt-2:revision-1", "digest": "sha256:..." },
    "room": "/abs/.../task/review/validation/briefing-2",
    "manifest": { "type": "Briefing", "id": "...", "question": "...", "artifacts": [ { "summary": "<the stored review>" } ] }
  }
}
```

- `current` is the attempt that `gate record` acts on, or `null` when the entity's status stage has no actionable gate attempt (no record, a stage that is not `gate: true`, a terminal or archived entity).
- `room` is the absolute path from `ResolveRoomRef`. `manifest` is the canonical Briefing object, read through `boundBriefingBytes` and verified against the bound digest. It is present only when `state` is `open`; otherwise it is `null`. If an open attempt's room does not resolve or verify, `gate read` exits 1 without JSON, which is the existing rule for invalid gate data.
- A surface reads `current.briefing.id`, shows `current.manifest`, and records with that id. It does not parse the entity file or recompute the canonical-bytes digest itself.

Fields considered and not added: a per-attempt state list (the rule is one line, already in `attemptState`: a withdrawal means withdrawn, a resolution means closed, neither means open; the command reference states it, and a decision needs only `current.state`), and `current.gate` (it is `records[].id` for the stage `entity.status`).

### 3. The review the captain reads is stored in the room, under the digest

There is no binary change. The review goes into the field that exists for it: the primary artifact's `summary` in the canonical Briefing (`--summary`, which has no length limit). So the Briefing digest covers the text the captain reads, and an id-matched decision answers exactly that text. Spike step 1 below stores a 21-line review verbatim. Step 3 shows that a different review cannot replace it on the same attempt: a new review needs a new attempt and a new id.

The change is the first officer's order. Today it prepares with a short summary first and composes the review afterwards, in chat only (`fo-gate-lifecycle` steps 1-4). The new order reads the checklist and AC evidence, composes the review, prepares with the review as `--summary`, commits, then presents. The stored review leaves out the `Reviewed snapshot:` line. That line names the Briefing digest, and a digest cannot cover a text that contains the digest. The presenter adds the line from prepare's `briefing=` and `digest=` output. When the first officer resumes at an open attempt, it presents the stored review from `gate read --json` and does not compose a new one. Prepare replays an identical review and refuses a different one.

No `--summary-file` flag: a 21-line argument with backticks, quotes and non-ASCII text reaches the room byte-exact through argv (spike step 1).

### 4. A Resolution can say how the decision arrived

New optional flag `--via TEXT` stores `via: TEXT` on the Resolution. It is valid with both actors, and must be nonblank, valid UTF-8 when given. When it is absent, the decision was given in the first officer's own chat. The binary does not verify `via`. It attributes the decision and authorizes nothing, like `conn`.

Recording rules, stated in `fo-gate-lifecycle`:
- The captain decided on a surface the captain operated (a pane, a Review page), and an agent carried the message: `--actor person:captain --via "Ready Room pane, relayed by MetaFO"`. So a reader can see that the first officer accepted another agent's report of the captain's decision.
- Another agent decided under a rule the captain gave it (this task's backlog gate): `--actor agent:first-officer` with `--conn-quote`/`--conn-source` citing that rule, as today. The actor set does not change.

### 5. A gate presentation reaches the captain whoever started the turn

This is a rule in `present-gate`, not a binary change. The first officer renders the review in its own captain-facing output every time. When another agent's message started the turn, the first officer also answers that agent with the entity, the stage and the Briefing id. The agent reads the stored review through `gate read --json`. A short reply to the agent never replaces the captain-facing presentation.

### Related backlog

- `gate-legibility-summary-versus-excerpt`: **separate.** It decides what a review contains (the FO's summary or the worker's own words, and how a workflow declares that). This task decides where the review lives and what binds it. Storing the review in the room lets its later excerpt rule reach second surfaces without more work.
- `subspace-gate-attempt-presentation-integration`: **separate.** It is Subspace transport (a blocking command, floats, `open`/EOF handling). It is a consumer of this task: it reads `current.manifest` and records with `--briefing-id`.
- `multi-artifact-gate-prepare`: **separate.** It extends prepare grammar. The review stays in the primary artifact's summary for any count of artifacts.
- `reliable-exact-digest-in-gate-review`: **separate.** Its AC is about how a live Sonnet journey displays the full digest. After this task, recording no longer depends on the agent copying the digest. Its conflict with `present-gate`'s "compact digest prefix" sentence stays with it.
- `gate-room-request-ordering-trap`: **separate, and its trap is not reachable on origin/main.** `gate record --room` is an unknown flag (exit 2), and `gate record --briefing` without `--round` exits 2. One-file rooms have no `request.json` (spike step 4). Its owner should check again and close it. It is not folded in here.

### Alternatives considered

- **Optional `--briefing-id`.** A surface that forgets it silently takes the unsafe path, which is the fault itself. Required costs mechanical updates of about 50 test call sites.
- **Origin in `--reason` text instead of a `via` field.** This has no stored-format change, so binaries of different versions do not conflict. But it mixes why with how the decision arrived, and a reader cannot find it mechanically. `via` is chosen, and the rollout cost is declared below.
- **A separate presentation file in the room, or a new Briefing field.** A file outside the manifest is not covered by the digest. A new field changes the Briefing v1 schema that Subspace and Review read. The existing `summary` needs neither.

## Risk evidence

Spike, 2026-11-24: `git archive origin/main` (a85a26a35) to `/tmp/spike-second-surface-gate-decisions-src`, built with `GOPROXY=off go build ./cmd/spacedock` (exit 0). Fixture: a split-root workflow (`validation` gate, then `implementation`, `done`), with the state on `spacedock-state/dev` and no origin.

1. **Today's late-decision fault (origin/main binary).** The commands were `gate prepare` (attempt 1, "plan A", exit 0), `state commit` (0), `gate withdraw` (0), `state commit` (0), `gate prepare` (attempt 2, "plan B", exit 0) and `state commit` (0). Then the captain's late approval of plan A ran: `gate record task --decision approve --actor person:captain --consume` gave exit 0, `recorded ... attempt=gate-attempt:task-validation-2 state=closed`, `consumed=true target-stage=implementation`. The entity moved to `status: implementation` on plan B, which the captain did not read.
2. **The proposed refusal (throwaway prototype, about 15 lines in the spike copy).** The same fixture ran with `--briefing-id briefing:task:validation:attempt-1:revision-1 --consume`. Result: exit 1 with the refusal message above, entity sha256 unchanged (`9aab9e73fbd1` before and after), no state commit added. With attempt 2's id: exit 0, closed and consumed. With no `--briefing-id`: exit 2. With `--via "Ready Room pane, relayed by MetaFO"`, the Resolution stores `via:`, and the prototype's `status --validate` gives `VALID`.
3. **The stored review and the digest (origin/main binary).** `gate prepare --summary "<21-line review>"` exit 0: `index.json` `artifacts[0].summary` equals the review byte for byte (21 lines). A replay with the identical review gave exit 0 and the same digest `sha256:a6285c80…`. The same open attempt with the review changed by one character gave exit 1, `open gate room binding is frozen and cannot be rebound`. So the digest covers the stored review, and a changed review needs a new attempt and a new id.
4. **The ordering trap (origin/main binary).** `gate record task --room …` exit 2 (`unknown gate flag: --room`). `gate record task --briefing …/index.json` exit 2 (`requires --round`).
5. **Stored-format skew (found by the spike).** The origin/main binary reads an entity whose Resolution has `via:` and gets `gate read` exit 1 and `status --validate` exit 1: `field via not found in type gates.Resolution`. The cause is `KnownFields(true)` at `internal/gates/io.go:79`. Older binaries fail closed and do not misread the entity. See Declared semantic changes.

The picard evidence paths (`/Users/clkao/git/picard/...`) were not readable from this sandbox (`Operation not permitted`). The design rests on the friction list `/tmp/spacedock-frictions-2026-11-24.md` and on the source.

## Out of scope

The pane itself; Subspace's Review page; an "act done" decision kind (task `gate-act-done-record`); what a review contains (`gate-legibility-summary-versus-excerpt`).

## Expected surface and tolerance

Estimate net LOC change: +360, across 16 files (about +395 insertions, -35 deletions). Tolerance: net +250 to +470, at most 20 files.

- Production, about +90: `internal/gates/operation.go` (+11: the input fields and the refusal), `internal/gates/model.go` (+4: `Via` and its validation), `internal/gates/current.go` or the same package (+35: the current-attempt view with the verified manifest), `internal/cli/cli.go` (+12: two flags, the required check, help text), `internal/status/gate_history.go` (+30: `entity` and `current`).
- Tests, about +250: `internal/cli/gate_test.go` (new refusal, `via` and `gate read` cases), `internal/status/testdata/gate-history/complete.json` (golden), the `--briefing-id` added at existing record call sites (`gate_test.go`, `gate_consume_sync_test.go`, `gate_ceremony_count_test.go`, `gate_one_file_room_test.go`, `terminal_consume_test.go`, `internal/ensigncycle/recorded_gate_lifecycle_test.go`, `shared_keep_moving_fixture_test.go`), and one ensigncycle oracle with transcripts for AC-6.
- Skills and docs, about +20: `skills/fo-gate-lifecycle/SKILL.md`, `skills/present-gate/SKILL.md`, `docs/site/reference/command-reference.md`.

### Declared semantic changes

- **Command grammar.** `gate record --decision` requires `--briefing-id` (exit 2 if missing) and refuses a non-current id (exit 1, no write). New optional `--via`. Both are refused with `--round`. `gate read --json` adds the top-level keys `entity` and `current`. The JSON of an entity with no gate history becomes `{"current":null,"entity":{...},"records":[],"version":1}`.
- **Stored format.** A Resolution can carry `via:`. Binaries without this change refuse every read of an entity that has one (spike step 5). So the release with this change must reach each binary that reads a workflow's state before a first officer records `--via` in it. The Briefing format does not change.
- **Authority.** None. The actors stay `person:captain` and `agent:first-officer`. `via` and `conn` attribute and are not verified.
- **Runtime behavior.** The first officer composes the review before `gate prepare` and stores it as `--summary`. On resume it presents the stored review. It presents to the captain whoever started the turn.

### Docs and skill text, before and after

`docs/site/reference/command-reference.md`, `gate read` paragraph:

Before:
> `spacedock gate read <ref-or-path> --json` returns the stored `gates` object, including every attempt and its application state.
> …
> An entity without gate history returns `{"records":[],"version":1}`.

After:
> `spacedock gate read <ref-or-path> --json` returns the stored `gates` object (`version`, `records`), including every attempt and its application state. An attempt with a `withdrawal` is withdrawn, one with a `resolution` is closed, and one with neither is open. It also returns `entity` (the resolved `path` and its `status`) and `current`: the attempt that `gate record` acts on, with its `state`, `briefing` id and digest, the resolved `room`, and, while the attempt is open, the verified Briefing `manifest`. `current` is `null` when the entity's stage has no actionable gate attempt.
> …
> An entity without gate history returns `records: []` and `current: null`.

Add after it:
> `spacedock gate record <entity> --decision … --briefing-id ID` records a decision for the Briefing `ID` only. If `ID` is not the current attempt's Briefing, the command exits 1 and changes nothing. `--via TEXT` records how a decision reached the first officer, for example a pane or another agent's message.

`gate --help`, record line. Before: `spacedock gate record <entity> --decision approve|revise|hold --actor ID [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]`. After: `spacedock gate record <entity> --decision approve|revise|hold --actor ID --briefing-id ID [--reason TEXT] [--via TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]`.

`skills/fo-gate-lifecycle/SKILL.md`, the prepare sequence. Before:
```text
1. ${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary SUMMARY [--reference FILE ...] --workflow-dir WORKFLOW_DIR
2. ${SPACEDOCK_BIN:-spacedock} state commit ENTITY --workflow-dir WORKFLOW_DIR
3. In one shell event, run `${SPACEDOCK_BIN:-spacedock} status --read ENTITY --checklist --json --workflow-dir WORKFLOW_DIR` and the same command with `--ac-scan` instead of `--checklist`.
4. Ensure spacedock:present-gate resident; present once from that fresh evidence.
```
After:
```text
1. In one shell event, run `${SPACEDOCK_BIN:-spacedock} status --read ENTITY --checklist --json --workflow-dir WORKFLOW_DIR` and the same command with `--ac-scan` instead of `--checklist`.
2. Ensure spacedock:present-gate resident; compose the gate review from that evidence, without its `Reviewed snapshot:` line.
3. ${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary GATE_REVIEW [--reference FILE ...] --workflow-dir WORKFLOW_DIR
4. ${SPACEDOCK_BIN:-spacedock} state commit ENTITY --workflow-dir WORKFLOW_DIR
5. Present once: the `Reviewed snapshot:` line from prepare's `briefing=` and `digest=`, then the stored review unchanged.
```
Initial-stage gate, before: "Skip step 3 — both structured reads exit nonzero with no report." After: "Skip step 1 — both structured reads exit nonzero with no report; compose the review from the seed."

Record block: each of the three record commands gains `--briefing-id BRIEFING` after `--actor …`, and `[--via SURFACE]` after `[--reason REASON]`. Add after the block:
> BRIEFING is the id of the Briefing the decision answers: prepare's `briefing=` for a decision given in your chat, or the id that the surface or relay names. Never substitute the current attempt for a different or missing id. On refusal, present the current attempt and wait. `--via` names the surface or agent that carried a captain's own decision to you (a pane, a Review page, another agent's message); omit it for a decision given in your chat. A decision that another agent made under a captain rule is not `person:captain`: record it as `agent:first-officer`, with the conn citing that rule.

Resume, before: "open → always `state commit`, then structured reads and presentation". After: "open → always `state commit`, then present the stored review from `gate read --json` `current.manifest` with its snapshot line; never compose a new one".

`skills/present-gate/SKILL.md`, Presentation channels. Before: "Delegated authority does not waive presentation. Only then hand the semantic decision and reason to `${SPACEDOCK_BIN:-spacedock} gate record <entity> --decision approve|revise|hold --actor ID [--reason TEXT]`." After: "Delegated authority does not waive presentation, and neither does who started the turn: when another agent's message starts the turn, render the review in your own captain-facing output, then answer that agent with the entity, stage and Briefing id, which it reads with `gate read --json`. Only then hand the semantic decision, reason and answered Briefing id to `${SPACEDOCK_BIN:-spacedock} gate record <entity> --decision approve|revise|hold --actor ID --briefing-id ID [--reason TEXT] [--via SURFACE]`." Add an assembly rule: "Compose the review before `gate prepare`; it is stored as `--summary`. Leave out the `Reviewed snapshot:` line, and add it at presentation from prepare's output."

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A late decision cannot land on a Briefing the captain did not read: a record that names a Briefing other than the current open attempt's is refused and changes nothing (value; baseline: origin/main closes and consumes the new attempt).**
Verified by: a CLI test that prepares attempt 1, withdraws it, prepares attempt 2, then records `--decision approve --consume --briefing-id <attempt-1 id>`. Expected: exit 1, the entity bytes and the state-commit count unchanged, status still the gate stage. The same record with attempt 2's id closes and consumes attempt 2. Baseline: on origin/main (spike step 1) the same late approval exits 0 and moves the entity to `implementation`. Falsifier: removing the id comparison makes the first record close attempt 2 and fails the test.

**AC-2 - Every chat decision names its Briefing.**
Verified by: a CLI test in which `gate record --decision` without `--briefing-id`, for each actor, exits 2 with entity bytes unchanged, and `--briefing-id` or `--via` with `--round` exits 2. Falsifier: making the flag optional lets the first case exit 0.

**AC-3 - One `gate read --json` call gives a second surface the entity, the current attempt and the verified review.**
Verified by: the golden `complete.json` and a CLI test on a two-attempt entity (attempt 1 withdrawn, attempt 2 open). They assert `entity.path`, `entity.status`, `current.attempt`, `current.state == "open"`, `current.briefing` (id and digest equal to the stored binding), `current.room`, and `current.manifest.artifacts[0].summary`. `current` is `null` after the gate is consumed and on an entity without gate history. An open room whose `index.json` has been edited makes the command exit 1. Falsifiers: reading `manifest` without the digest check lets the edited room print JSON; choosing the first attempt instead of the last returns attempt 1.

**AC-4 - The room holds the full review the captain was shown (value; baseline: the mail-triage room held a 147-character summary of a roughly 20-line review).**
Verified by: a deterministic test that `gate prepare` stores a multi-line review byte-exact as `artifacts[0].summary`, and that a changed review on the same open attempt is refused. Also a live Claude recorded-gate journey that asserts that the captain-facing review, minus its snapshot line, equals the stored summary of the attempt it presents. Falsifier: restoring the prepare-first order in `fo-gate-lifecycle` leaves the stored summary as the short prepare-time text, and the live check fails.

**AC-5 - A Resolution records how a captain decision arrived, without a change to who can decide.**
Verified by: a CLI test in which `--via "Ready Room pane, relayed by MetaFO"` with `person:captain` stores `via:` on the Resolution and `status --validate` passes; with no `--via`, no `via` key is stored; `--via ""` exits 2; `--actor agent:metafo` is still refused. Falsifier: dropping the assignment of `Via` leaves no `via:` in the stored Resolution.

**AC-6 - The gate review reaches the captain when another agent starts the first officer's turn.**
Verified by: a live Claude recorded-gate case whose turn-starting prompt is another agent's message. Its oracle requires that the first officer's own final output contains the review (entity, stage, Briefing id, decision ask), and that the reply to the agent names the Briefing id. Oracle unit tests with transcripts: a reply-only transcript fails, and a transcript with the review passes. Falsifier: removing the new present-gate sentence is expected to fail the live case, and the reply-only transcript fails the oracle deterministically.

## Test plan

- **AC-1, AC-2, AC-5 (deterministic, low cost).** The proof owner is `internal/cli/gate_test.go`, which already owns the gate record grammar and refusals and has the `gatePrepareCLIFixture`/`gateOneFileCLIFixture` substrate. The spike script (two attempts, then a late record) becomes the first test. Each refusal compares entity sha256 and the state-checkout commit count before and after.
- **AC-3 (deterministic, low cost).** The proof owner is the `gate read` tests in `internal/cli/gate_test.go` with the golden `internal/status/testdata/gate-history/complete.json`. The edited-room case is the distinct falsifier for the digest check.
- **AC-4 (deterministic and live, medium cost).** The binary half extends `internal/gates/prepare_test.go` (byte-exact multi-line summary; changed summary refused). The conduct half extends the Claude recorded-gate live journey in `internal/ensigncycle/recorded_gate_lifecycle_test.go`, run once at the candidate tip.
- **AC-6 (live, medium cost).** One Claude live case in `internal/ensigncycle` with a deterministic oracle unit test and its pass and fail transcripts. There is no existing proof owner for agent-started turns. That is why this case is new.
- **Existing suites.** The call sites updated with `--briefing-id` keep their current assertions. `go test ./...`, `go test ./... -race`, `gofmt`, and `internal/contractlint` (it reads the gate skills) all pass. Live runs of Codex and Pi are not owed: the skill change is host-neutral, and AC-6 proves the rule on one host.

### Feedback Cycles

## Stage Report: ideation

- DONE: One design for the five group A items, with exact grammar and semantics: the expected briefing id and digest input to `gate record` and its refusal on mismatch (entity bytes unchanged); the extended `gate read --json` envelope (entity path and status, every attempt with its state, the current attempt, the resolved room); where and how the presentation the captain reads is stored with the gate and how that interacts with the briefing digest; the origin field on a resolution (including a decision a pane or MetaFO relays); and the rule that a gate presentation reaches the captain whoever started the turn. Say for each related backlog task (gate-legibility-summary-versus-excerpt, subspace-gate-attempt-presentation-integration, multi-artifact-gate-prepare, reliable-exact-digest-in-gate-review, gate-room-request-ordering-trap) whether it folds in or stays separate, and why; challenge whether each new field or flag is necessary before adding it.
  Proposed approach sections 1-5. Challenged and not added: `--digest` (ids are frozen per prepared attempt), a per-attempt state list, `current.gate`, `--summary-file`, a new actor. All five related tasks stay separate; the ordering trap is not reachable on origin/main.
- DONE: The riskiest mechanism is exercised, not asserted: with a spacedock binary built from origin/main, a two-attempt fixture shows today's late-decision fault (a record after re-prepare closes the new attempt), and a throwaway prototype or reasoned trace shows the proposed refusal; also show what a stored presentation does to the briefing digest. Record commands and exit codes.
  Risk evidence steps 1-5: the origin/main late approval exits 0 and consumes attempt 2. The prototype refusal exits 1 with entity sha256 unchanged. A 21-line review is stored byte-exact and frozen under the digest. The spike also found that older binaries reject a stored `via:`.
- DONE: Acceptance criteria under the exact `## Acceptance criteria` heading in `**AC-N**` form, at least one measuring the end value against today's baseline, each with a falsifying edit; a test plan naming proof owners; expected surface (net LOC, files, tolerance); declared semantic changes (command grammar, stored formats, authority); before/after wording for any docs/site and first-officer skill text change.
  AC-1 to AC-6. AC-1 and AC-4 measure against the baseline. Net +360 across 16 files (tolerance +250 to +470, at most 20 files). Semantic changes and before/after wording for the command reference, `gate --help`, `fo-gate-lifecycle` and `present-gate` are in the body.

### Summary

The design binds a gate decision to the Briefing it answers: `gate record` requires `--briefing-id` and refuses a non-current id with no write. `gate read --json` adds `entity` and `current`, including the verified manifest. The first officer stores its full review as the Briefing summary, so the digest covers what the captain reads. An optional `via` records how a captain decision arrived, and a present-gate rule keeps the review in front of the captain whoever started the turn. Decision points for CL at the gate: id-only instead of id plus digest; a required flag; and the `via` stored-format change, which older binaries reject when they read the entity.
