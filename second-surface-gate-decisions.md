---
id: 301s961pr7b2cnywnkz88jec
title: A gate that a second surface can read and decide safely
status: validation
source: "MetaFO friction list 2026-11-24, sent on the word of CL, /tmp/spacedock-frictions-2026-11-24.md, group A items 1-5"
group: A-second-surface-gate
started: 2026-11-25T01:34:10Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-second-surface-gate-decisions
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
        - id: gate:301s961pr7b2cnywnkz88jec:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:301s961pr7b2cnywnkz88jec-ideation-1
              briefing:
                id: briefing:301s961pr7b2cnywnkz88jec:ideation:attempt-1:revision-1
                digest: sha256:cb1eca5184d7bacf397ec2473b0caf78d11e238b2c52d76d3ce2a0cf7fd3381e
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:301s961pr7b2cnywnkz88jec:ideation:1
                briefing: briefing:301s961pr7b2cnywnkz88jec:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-11-25T06:34:33.678138Z"
                decision: approve
                reason: |-
                    Approve: CL approved the group A ideation design for briefing:301s961pr7b2cnywnkz88jec:ideation:attempt-1:revision-1 (sha256:cb1eca5184d7bacf397ec2473b0caf78d11e238b2c52d76d3ce2a0cf7fd3381e) with the words 'yes to all'; the line CL saw was 'The Ready Room pane depends on it.' The design's listed decisions apply with the first officer's recommendations.
                    Via: MetaFO, relaying CL's answer of 2026-11-24 22:33 PT.
              application:
                target-stage: implementation
                state: consumed
---

A captain can decide a gate from a pane, a Spacedock Review page or another agent's message. The decision lands only on the Briefing content the captain read, the stored review is the text the captain read, and the record says how the decision arrived.

## Problem

A second surface is anything other than the first officer's own chat. Checked against origin/main `a85a26a35`:

1. A decision cannot name the Briefing it answers. `gate record` closes whatever attempt is current (`internal/gates/operation.go:25` `RecordInput` has no expected Briefing; `:248` `currentStageAttempt`). A late decision can land on a Briefing the captain did not read. The same fault was seen on the Review page: a decision went to version 1 while the agent waited on version 2. An id alone does not fix it: prepare can bind new content to an open attempt under the same id when the attempt's room directory is missing or its binding names a Briefing file directly (Spike 2, steps 1 and 2).
2. A surface cannot get the current decision coordinate and the review text from the binary. `gate read --json` returns only the stored `gates:` block (`internal/status/gate_history.go:13-52`). A pane parses the entity file and the room itself, and it does not verify the digest of the text it shows.
3. The review the captain reads is not stored with the gate. `--summary` has no length limit (`internal/gates/prepare.go:85-89`), but the first officer prepares with a short summary and composes the review afterwards, in chat only (`skills/fo-gate-lifecycle/SKILL.md:22-25`). The mail-triage room held a 147-character summary of a review of about 20 lines.
4. A captain Resolution has no place that says where it came from: `by`, `at`, `decision`, `reason` (`internal/gates/model.go:69-78`), and conn fields are refused for `person:captain` (`internal/gates/operation.go:246`).
5. When another agent starts a first officer's turn, the first officer wrote a short report for that agent and the gate review did not appear in its output. A manager-prompt line now corrects it in one project; the product rule is missing.

## Proposed approach

One rule ties the five items together: **a decision names the Briefing id and digest it answers, and the Briefing holds the review the captain read.** The binary gets two required record flags and one new read mode. Stored formats do not change. The first-officer skills get a new order (compose the review, prepare, commit, then present the verified stored text), a transition rule for open attempts prepared before this change, an origin line in the reason, a presentation rule for agent-started turns, and a stop rule for a binary that is older than the skills.

### 1. `gate record` names the Briefing id and digest it answers

Grammar (chat decisions; `--round` is unchanged):

```text
spacedock gate record <entity> --decision approve|revise|hold --actor ID --briefing-id ID --briefing-digest DIGEST
    [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]
```

- Both flags are required on every `--decision` record, for both actors. All record grammar checks (decision, actor, conn pairing, `--consume` with a non-approve decision, and the two binding flags) run before the entity reference is resolved. So a malformed command exits 2 even when the entity does not resolve.
- Under the entity lock, the recorder compares the pair with the current attempt's stored `briefing.id` and `briefing.digest`. If either differs, it exits 1 before the close: the entity bytes, the state commit, the sync and the application stay unchanged, and `--consume` consumes nothing.
- The refusal names the pair, the attempt of this gate that binds that pair (or says that no attempt binds it), and the current attempt's pair and state: `Error: decision names briefing:task:validation:attempt-1:revision-1 sha256:4a11…, and that is attempt gate-attempt:task-validation-1 (withdrawn); the current attempt is briefing:task:validation:attempt-2:revision-1 sha256:4910… (open); nothing was recorded`. The message prints both digests in full.
- A match then runs today's checks unchanged (open state, retained authority, verifiable digest, actor and conn rules).
- `--briefing-id` and `--briefing-digest` are refused with `--round` (exit 2), like the conn flags.

The digest is required because the id does not name one set of bytes. The freeze in prepare (`internal/gates/prepare.go:247`) applies only when the previous binding is classified as a prepared room. A request-free room whose directory is missing, and a legacy binding whose `room-ref` names a Briefing file, are not classified so. Prepare then rebinds the open attempt in place with the same id and a new digest (`prepare.go:270`), and an id-only check accepts a late decision for the earlier content (Spike 2, steps 1 and 2). With the digest compare, that in-place re-prepare stays available as the recovery for both shapes: it is safe because a decision for the earlier digest is refused (Spike 2, step 4).

The pair comes from the presentation the decision answers, never from a later read. For a decision in the first officer's chat or under conn, it is the binding the first officer presented (section 3). For a pane or Review page, it is the binding that surface displayed. For a relayed decision, it is the binding the relaying agent names for the presentation the captain read. A copied digest with a typing error is refused like any other mismatch, so the error costs one more presentation and cannot land a decision on other content.

### 2. `gate read --current --json`: the decision coordinate and the verified review

`gate read <ref> --json` keeps its history contract exactly: the stored `version` and `records`, no room access, readable for a missing or edited room and for archived entities. The verified presentation is a separate mode, so its failure does not take the history away:

```text
spacedock gate read <ref-or-path> --current --json [--workflow-dir DIR]
```

```json
{
  "stage": "validation",
  "attempt": {
    "id": "gate-attempt:task-validation-2",
    "state": "open",
    "briefing": { "id": "briefing:task:validation:attempt-2:revision-1", "digest": "sha256:..." },
    "room": "/abs/.../task/review/validation/briefing-2",
    "question": "...",
    "summary": "<the stored review>"
  }
}
```

- `stage` is the entity's status. `attempt` is the last attempt of that stage's gate record: the attempt `gate record` acts on. It is `null` when that stage has no gate attempt (no record, a stage that is not gated, a terminal or archived entity). `state` is `open`, `withdrawn` or `closed`, by the existing rule in `attemptState` (`internal/gates/model.go:403`).
- For an open attempt only, `room` is the resolved `room-ref` (`ResolveRoomRef`), and `question` and `summary` come from the bound Briefing bytes after `boundBriefingBytes` verifies them against the stored digest. The binding and the text come from the same bytes. Legacy locators resolve as they do for the recorder. A Briefing with no primary summary returns no `summary` key.
- If the open attempt's Briefing does not resolve or does not match its digest, the mode exits 1 with no JSON and names the attempt and the failure. `gate read --json` on the same entity still exits 0.

Named consumers: the first officer presents from it on every gate (section 3) and on resume; a pane or Subspace presents `summary` and records with `attempt.briefing`; `room` is the coordinate that a room-aware presenter (Subspace) opens on resume, when no prepare output is available. The entity path and the parsed manifest are not returned, because no consumer needs them: a surface already holds the reference it passed.

### 3. The review the captain reads is stored in the room and presented from the room

The review goes into the field that exists for it: the primary artifact's `summary` in the canonical Briefing. So the Briefing digest covers the text the captain reads, and a decision bound to that digest answers exactly that text. Prepare already stores the summary byte-exact and refuses a different summary on the same prepared open attempt (Spike 1, step 3).

The first officer's order changes. It reads the checklist and AC evidence, composes the review, prepares with the review as `--summary`, commits, then reads `gate read --current --json` and presents the `Reviewed snapshot:` line from `attempt.briefing` followed by `attempt.summary` unchanged. Presentation always comes from the verified stored text after the commit, on a fresh gate and on resume, so the interrupted case (prepare wrote the binding, the commit did not run) presents what is stored once the resume commit runs. The first officer also requires `attempt.briefing` to equal prepare's `briefing=` and `digest=` output, and stops before presentation if they differ.

The stored review starts with the template's `Gate review:` line, keeps its `Decision:` line, and leaves out the `Reviewed snapshot:` line: that line names the digest, and the digest cannot cover a text that contains it. At an `initial: true` stage the review is composed from the committed seed, as today. A changed review after a withdrawal or a revise belongs to the successor attempt; at an open attempt that already holds a full review, the first officer never composes a new one. Exact prepare replay needs the same question, summary and logical source revisions, so a re-composed review on an open attempt is refused as divergent.

### 4. Open attempts prepared before this change

An open attempt prepared under the old order holds a short synopsis. On resume, the first officer reads `gate read --current --json`. If `summary` is missing, does not start with a `Gate review:` line, or has no `Decision:` line, the attempt predates stored reviews:

- A prepared room (request-free or request-backed): `gate withdraw` with the reason `stored summary is not the full gate review`, commit, then compose, prepare, commit and present the successor attempt. The successor has a new id.
- A legacy direct-file binding: `gate withdraw` refuses it (`current attempt has no prepared gate room`, Spike 2, step 3). The first officer runs the prepare sequence on the same attempt. Prepare rebinds the attempt in place under the same id with a new digest, so a decision for the old synopsis is refused (section 1).
- A request-free room whose directory is missing: `gate read --current` exits 1, and the first officer stops with that error. The captain chooses to restore the room from the state history or to have the first officer re-prepare in place. Either way, a decision for content the captain did not read is refused.

### 5. A relayed decision records its origin in the reason

No new field. When the captain decided on another surface, or an agent carried the captain's decision to the first officer, the first officer records `--actor person:captain` and ends the reason with one labeled line after the substantive reason:

```text
--reason $'Approve: the design covers all five items and the spike shows the refusal.\nVia: Ready Room pane, relayed by MetaFO.'
```

A reason without a `Via:` line means the decision was given in the first officer's own chat, or its origin is unspecified for records made before this change. A decision that another agent made under a rule the captain gave it stays `--actor agent:first-officer` with `--conn-quote`/`--conn-source` citing that rule, as today; the grant evidence stays in `conn`. The actor set does not change. The stored format does not change, so every binary reads the record (Spike 2, step 6). The binary does not verify the line: like `conn`, it attributes and authorizes nothing.

### 6. A gate presentation in a turn that another agent started

This is a rule in `present-gate`. When the selected channel is chat and another agent's message starts the first officer's turn, the presentation goes in the first officer's own output for that turn, once, before any decision write. The answer to that agent adds only the entity, the stage and the Briefing id and digest; it never replaces the presentation. When Subspace is the selected channel, the presentation stays there and is not repeated in chat.

This task owns that emission. It does not own delivery: the host session the captain watches (a Claude Code session, a Paseo agent view) shows the first officer's output, and an agent that started the turn owns what it shows the captain. Section 2 lets that agent read the same verified review.

### 7. Rollout and compatibility

The binding flags change the record protocol in both directions, so the binary and skills ship as one matched change:

- **One change, one release.** One PR changes the binary and both gate skills. Source builds with `--plugin-dir` are matched by construction. A stable tag carries both. The change ships in a minor release, not in a patch to a line whose skills lack it, so the minor-exact boot gate (`first-officer-shared-core.md`, "These skills require binary minor X.Y") keeps stable installs matched. The edge plugin resolves `main` on merge while the edge binary follows prerelease tags, so an edge prerelease is cut right after the merge (`docs/releasing.md`).
- **Skills older than the binary.** An old first officer's record commands carry no binding flags. The new binary exits 2 before it resolves the entity: `Error: gate record --decision requires --briefing-id and --briefing-digest of the Briefing the decision answers. The calling skills can be older than this binary: update the Spacedock plugin (see spacedock doctor). Nothing was recorded.` The old skill halts on a record failure ("Close/sync failure halts"); a headless conn run halts too, with the gate open and nothing consumed. Recovery: update the plugin. The new first officer resumes at the open attempt, which still has its presented binding because nothing was written. It presents again. It records a decision made before the update only with a binding that the presentation or the relay named; otherwise it asks again, or decides again under the grant. It never retries without the flags or takes a missing value from a later read.
- **Skills newer than the binary.** The first new command at any gate is `gate read --current`, before presentation. An older binary exits 2 with `unknown gate read flag: --current` (Spike 2, step 5). The new skill stops there with the update diagnostic (`spacedock doctor`) and changes nothing. If a record reaches an older binary, `unknown gate flag: --briefing-id` stops it the same way. A pinned `SPACEDOCK_BIN` that is older than the plugin is this case.
- **Advisory rounds.** `gate record --round` keeps its grammar on both binaries. The binding flags are refused with it. A round publishes once with its original command and commit, and gate re-entry prepares a fresh attempt with a new id.

### Related backlog

- `gate-legibility-summary-versus-excerpt`: **separate.** It decides what a review contains. This task decides where the review lives and what binds it. Its later excerpt rule reaches second surfaces through the stored summary with no more work.
- `subspace-gate-attempt-presentation-integration`: **separate.** It is Subspace transport (a blocking command, floats, `open`/EOF handling). It consumes this task: it reads `gate read --current` and records with the displayed binding. The deferred rendering risk below is checked there.
- `multi-artifact-gate-prepare`: **separate.** It extends prepare grammar. The review stays in the primary artifact's summary for any count of artifacts.
- `reliable-exact-digest-in-gate-review`: **separate.** Its AC is about how a live journey displays the full digest. After this task, the recorder no longer trusts an id alone, and the first officer takes the digest from `gate read --current`. Its conflict with `present-gate`'s "short digest prefix" sentence stays with it.
- `gate-room-request-ordering-trap`: **separate, and its trap is not reachable on origin/main.** `gate record --room` exits 2 (`unknown gate flag: --room`), and `gate record --briefing` without `--round` exits 2 (Spike 1, step 4). Its owner should check again and close it.

### Alternatives considered

- **Immutable content per id, with an id-only record flag.** Prepare would refuse every in-place rebind of an open attempt. That needs the freeze to cover a missing request-free room and a legacy direct-file binding, and then those attempts have no recovery: `gate withdraw` refuses both shapes (Spike 2, step 3), so withdraw would also change. The digest compare needs no change to prepare or withdraw and keeps in-place re-prepare as a safe recovery.
- **Digest only.** The digest of a prepared Briefing covers its id, but a captain, a relay and a refusal message read the id. Requiring both keeps the human-readable half and the content half together at the cost of one flag.
- **Optional binding flags.** A surface or skill that forgets them takes the unsafe path silently, which is the fault itself. Required flags cost mechanical updates of about 50 test call sites.
- **A structured `via:` field on the Resolution.** A reader can filter it, but no consumer filters by origin, and any binary without the field refuses every read of an entity that has one, even after the gate closes (Spike 1, step 5). The labeled line gives the same attribution with no format change.
- **A wider `gate read --json` with the verified manifest.** It makes history fail when a room is missing or edited, which is when an operator needs the history. A separate mode keeps the two failures apart.
- **A separate presentation file in the room, or a new Briefing field.** A file outside the canonical Briefing is not covered by the digest. A new field changes the Briefing v1 schema that Subspace and Review read.

## Risk evidence

**Spike 2 (this revision), 2026-11-24.** `git archive origin/main` (a85a26a35) to `/tmp/spike-second-surface-gate-decisions-r2-src`, built with `GOPROXY=off go build ./cmd/spacedock` (exit 0). Fixture: a split-root workflow (`validation` gate, then `implementation`, `done`), state on `spacedock-state/dev`, no origin. Removed after the run.

1. **Missing modern room, same id, new digest (origin/main binary).** `gate prepare --summary "plan A"` exit 0: `briefing=briefing:task:validation:attempt-1:revision-1`, `digest=sha256:bc64ca7b…`; `state commit` exit 0 (the room is in the state history). `rm -rf` the room directory. `status --validate` exit 0, `VALID`, with `Warning: retained gate room does not resolve`. `gate prepare --summary "plan B"` exit 0 with the **same id** and `digest=sha256:4cb4204d…`; `state commit` exit 0. `gate record --decision approve --actor person:captain --consume` exit 0, `consumed=true target-stage=implementation`, stored digest `4cb4204d…`. A captain who read plan A under that id has approved plan B.
2. **Legacy direct-file binding, same id, new digest (origin/main binary).** After prepare, the Briefing moved to `task/inputs/briefing.json` and `room-ref` set to that file (digest `acea2e7b…`); `status --validate` exit 0 `VALID`. `gate prepare --summary "plan B"` exit 0: same id, `digest=sha256:bca9362e…`, `room-ref` rewritten to `@review/validation/briefing-1`.
3. **Withdraw cannot retire either shape (origin/main binary).** `gate withdraw --reason stale` exit 1, `current attempt has no prepared gate room`, on the missing room and on the direct-file binding.
4. **The proposed refusal (throwaway prototype, about 25 lines: two `RecordInput` fields, the compare in `recordChatLocked` before the open-state check, a CLI check before path resolution).** Missing-room fixture after step 1's re-prepare: record with the id and plan A's digest, `--consume`: exit 1 with the refusal; entity sha256 `62e2d73c03da` and state commit count 3 before and after; no resolution written. The id without the digest: exit 2. No binding flags and an unknown entity `nosuch`: exit 2. The current id and digest: exit 0, consumed. Direct-file fixture: the old digest with `--consume` exit 1, sha256 `1a1dd31b7ba5` and commit count 3 unchanged; the delegated shape (`agent:first-officer`, reason, both conn flags) exit 1, unchanged. Successor fixture (attempt 1 withdrawn, attempt 2 open): attempt 1's pair exit 1, `that is attempt gate-attempt:task-validation-1 (withdrawn); the current attempt is briefing:task:validation:attempt-2:revision-1 … (open)`.
5. **Skills newer than the binary (origin/main binary).** `gate read task --current --json` exit 2, `unknown gate read flag: --current`. `gate record … --briefing-id x --briefing-digest y` exit 2, `unknown gate flag: --briefing-id`, also with an unknown entity. Today's grammar runs after path resolution: `gate record nosuch --decision approve` (no actor) exits 1, `unknown reference: nosuch`.
6. **Origin line and history (origin/main binary).** A captain record whose reason is `Approve: plan B covers the fixture.` + newline + `Via: Ready Room pane, relayed by MetaFO.` (written by the prototype): origin/main `status --validate` exit 0 `VALID`; `gate read --json` exit 0 and returns `by: person:captain` and the reason byte-exact. `gate read --json` also exits 0 with the room directory missing and with an edited `index.json`.

**Spike 1 (first ideation cycle), origin/main binary, same fixture shape.** (1) A late `gate record --decision approve --consume` after withdraw and re-prepare exits 0 and moves the entity to `implementation` on attempt 2. (3) `gate prepare --summary "<21-line review>"` stores the review byte-exact as `artifacts[0].summary`; an identical replay exits 0 with the same digest; a one-character change on the same open attempt exits 1, `open gate room binding is frozen and cannot be rebound`. (4) `gate record --room` exits 2 (`unknown gate flag: --room`); `gate record --briefing …` without `--round` exits 2. (5) A Resolution with a `via:` key makes `gate read` and `status --validate` exit 1, `field via not found in type gates.Resolution` (`KnownFields(true)`, `internal/gates/io.go:79`).

No spike for `gate read --current`: it composes `ResolveRoomRef`, `boundBriefingBytes` and the Briefing parser, which the record, withdraw and replay paths already run on the same bindings.

The picard evidence paths (`/Users/clkao/git/picard/...`) were not readable from this sandbox. The design rests on the friction list `/tmp/spacedock-frictions-2026-11-24.md` and on the source.

## Deferred risk

**External rendering of the full summary.** The producer stores the summary byte-exact and replays it exactly. That does not show whether a Review page or Subspace truncates Markdown, hides findings, or returns the binding it displayed when a newer attempt exists. Trigger: Subspace or Spacedock Review is named as a supported surface for second-surface decisions. Promote to Material before that support claim if a representative multi-line review (template headings, a findings list, backticks, non-ASCII text) cannot be read in full on that surface, or if the decision it returns does not carry the id and digest it displayed.

## Out of scope

The pane itself; Subspace's Review page; delivery of the first officer's output to the captain (section 6); an "act done" decision kind (task `gate-act-done-record`); what a review contains (`gate-legibility-summary-versus-excerpt`).

## Expected surface and tolerance

Estimate net LOC change: +675, across 19 files (about +760 insertions, -85 deletions). Tolerance: net +520 to +830, at most 23 files.

- Production, about +150 / -15, 5 files: `internal/gates/operation.go` (+25: the two input fields, the compare and the refusal text), `internal/gates/current.go` (new, +60: the current presentation from the verified bytes), `internal/cli/cli.go` (+35 / -15: two flags, record grammar ahead of path resolution, help text), `internal/cli/gate_read.go` (+8: `--current`), `internal/status/gate_history.go` (+25: the `--current` JSON writer).
- Tests, about +565 / -50, 10 files: `internal/cli/gate_test.go` (+260: the refusal matrix, the grammar and mixed-pair table, the `--round` refusals, the history and `--current` state matrices), two new goldens under `internal/status/testdata/gate-history/` (+25), the binding flags at existing record call sites (`gate_test.go`, `gate_consume_sync_test.go`, `gate_ceremony_count_test.go`, `gate_one_file_room_test.go`, `terminal_consume_test.go`, `internal/ensigncycle/recorded_gate_lifecycle_test.go`, `shared_keep_moving_fixture_test.go`; +60 / -50, including the resume-matrix additions), and one new `internal/ensigncycle` file with the two-phase live scenario, its oracle and the oracle's unit tests (+220).
- Skills and docs, about +45 / -20, 4 files: `skills/fo-gate-lifecycle/SKILL.md`, `skills/present-gate/SKILL.md`, `docs/site/reference/command-reference.md`, `docs/specs/gate-resolution-frontmatter-contract.md`.

### Declared semantic changes

- **Command grammar.** `gate record --decision` requires `--briefing-id` and `--briefing-digest` (exit 2 if either is missing) and refuses a pair that is not the current attempt's binding (exit 1, no write, no sync, no consume). Both flags are refused with `--round`. Record grammar errors now exit 2 before the entity is resolved; for example `gate record nosuch --decision approve` without `--actor` changes from exit 1 (`unknown reference`) to exit 2. New read mode `gate read --current --json`. `gate read --json` output does not change.
- **Stored format.** None. The Resolution and Briefing schemas do not change. A `Via:` line is a convention inside the existing `reason` text.
- **Authority.** None. The actors stay `person:captain` and `agent:first-officer`. The `Via:` line and `conn` attribute and are not verified.
- **Runtime behavior.** The first officer composes the review before `gate prepare`, stores it as `--summary`, and presents the verified stored text from `gate read --current` after the commit, on a fresh gate and on resume. It records with the binding of the presentation the decision answers. It withdraws or re-prepares an open attempt whose stored summary is not a full review. It adds a `Via:` line for relayed captain decisions. It presents in its own output when another agent starts the turn. It stops with an update diagnostic when the binary is older than the skills.

### Docs and skill text, before and after

`docs/site/reference/command-reference.md`, `gate read` paragraph. The existing paragraph does not change. Add after it:
> `spacedock gate read <ref-or-path> --current --json` returns the presentation that a decision answers now. `stage` is the entity status. `attempt` is the last attempt of the gate for that stage, with its `id`, its `state` (`open`, `withdrawn` or `closed`) and its `briefing` `id` and `digest`. For an open attempt, it also returns the resolved `room`, and the `question` and `summary` of the bound Briefing after the command verifies its digest. `attempt` is `null` when the stage has no gate attempt. If the open attempt's Briefing is missing or does not match its digest, the command exits 1 without JSON. `gate read --json` still returns the history.
>
> `spacedock gate record <entity> --decision … --briefing-id ID --briefing-digest DIGEST` records a decision only when `ID` and `DIGEST` are the binding of the current attempt. Otherwise the command exits 1 and changes nothing. A decision without both flags is a usage error (exit 2).

`gate --help`. Record line, before: `spacedock gate record <entity> --decision approve|revise|hold --actor ID [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]`. After: `spacedock gate record <entity> --decision approve|revise|hold --actor ID --briefing-id ID --briefing-digest DIGEST [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]`. Read line, before: `spacedock gate read <ref-or-path> --json [--workflow-dir DIR]`. After: `spacedock gate read <ref-or-path> --json [--current] [--workflow-dir DIR]`. Add to the help paragraph: "`--briefing-id` and `--briefing-digest` name the Briefing the decision answers; a pair that is not the current attempt's binding is refused and nothing is recorded."

`docs/specs/gate-resolution-frontmatter-contract.md`, command surface (line 274). Before: `spacedock gate record ENTITY --decision approve|revise|hold --actor ID [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]`. After: `spacedock gate record ENTITY --decision approve|revise|hold --actor ID --briefing-id ID --briefing-digest DIGEST [--reason TEXT] [--conn-quote TEXT --conn-source TEXT] [--consume] [--workflow-dir DIR]`, and add `spacedock gate read REF --json [--current] [--workflow-dir DIR]`. Recorder lifecycle, add after the withdraw paragraph: "A `--decision` record names the Briefing id and digest it answers. Under the lock, the recorder refuses any pair that is not the current attempt's binding, with no write, sync or consume. This includes an open attempt that prepare rebound in place under the same id." Behavioral proof, add: "7. A decision whose Briefing id and digest are not the current attempt's binding changes the entity, the state commit or the sync."

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
5. ${SPACEDOCK_BIN:-spacedock} gate read ENTITY --current --json --workflow-dir WORKFLOW_DIR
6. Present once: the `Reviewed snapshot:` line from step 5's `attempt.briefing`, then step 5's `attempt.summary` unchanged. That id and digest are the binding the decision answers.
```
Initial-stage gate, before: "Skip step 3 — both structured reads exit nonzero with no report." After: "Skip step 1 — both structured reads exit nonzero with no report; compose the review from the seed."

After "Require prepare to emit `room`, `briefing`, `digest`, and `state=open`; never reconstruct the room.", add: "Require step 5's `attempt.briefing` to equal prepare's `briefing=` and `digest=`; a difference stops before presentation."

New paragraph after the prepare paragraphs: "**Binary older than these skills.** If `gate read --current` exits 2 with `unknown gate read flag`, or `gate record` exits 2 with `unknown gate flag: --briefing-id`, the binary predates these skills. Stop before presentation or recording, tell the Captain to update the binary (`spacedock doctor`), and change nothing. Never retry without the flags."

Record block, after: each of the three record commands gains `--briefing-id ID --briefing-digest DIGEST` after its `--actor …`. Add after the block:
> ID and DIGEST are the binding of the presentation the decision answers: your step-6 presentation for a decision in your chat or under conn, or the id and digest that a pane, a Review page or a relaying agent names for the presentation the Captain read. If a relayed decision names no digest, do not record: present the current attempt and wait. Never take a missing value from a later read. On a refusal (`nothing was recorded`), present the current attempt and wait. When the Captain decided on another surface, or an agent carried the Captain's decision to you, record `--actor person:captain` and end REASON with one line after the substantive reason: `Via: SURFACE, relayed by AGENT.` A decision that another agent made under a Captain rule is not `person:captain`: record it as `agent:first-officer`, with the conn citing that rule.

Resume, before: "open → always `state commit`, then structured reads and presentation". After: "open → always `state commit`, then steps 5 and 6; a failure stops. If `attempt.summary` is missing, does not start with a `Gate review:` line, or has no `Decision:` line, it predates stored reviews: withdraw with reason `stored summary is not the full gate review`, commit, then run steps 1 to 6 for the successor. If withdraw refuses with `current attempt has no prepared gate room`, run steps 1 to 6 on the same attempt. Never compose a new review for an open attempt that holds a full one."

`skills/present-gate/SKILL.md`, Presentation channels. Before: "After the prepared room and binding commit, render the gate template exactly once in the selected channel before the next decision-mutation tool call." After: "After the prepared room and binding commit, read the stored review with `${SPACEDOCK_BIN:-spacedock} gate read <entity> --current --json` and render it exactly once in the selected channel before the next decision-mutation tool call: the `Reviewed snapshot:` line from `attempt.briefing`, then `attempt.summary` unchanged."

Before: "Delegated authority does not waive presentation. Only then hand the semantic decision and reason to `${SPACEDOCK_BIN:-spacedock} gate record <entity> --decision approve|revise|hold --actor ID [--reason TEXT]`." After: "Delegated authority does not waive presentation, and neither does who started the turn. When the selected channel is chat and another agent's message starts your turn, the presentation goes in your own output for that turn; your answer to that agent adds only the entity, stage, Briefing id and digest, and never replaces the presentation. Only then hand the semantic decision, reason and answered binding to `${SPACEDOCK_BIN:-spacedock} gate record <entity> --decision approve|revise|hold --actor ID --briefing-id ID --briefing-digest DIGEST [--reason TEXT]`."

Captain-facing assembly rules, add: "Compose the review before `gate prepare`; it is stored as `--summary`. It starts with the template's `Gate review:` line, keeps its `Decision:` line, and leaves out the `Reviewed snapshot:` line, which presentation adds from `gate read --current`."

## Acceptance criteria

Each AC names a property of the finished entity and how it is verified.

**AC-1 - A decision lands only on the Briefing content it names, including when prepare rebound an open attempt under the same id (value; baseline: origin/main closes and consumes whatever attempt is current, Spike 1 step 1 and Spike 2 step 1).**
Verified by: CLI tests in `internal/cli/gate_test.go`. Each records `--decision approve --consume` and compares the entity sha256 and the state commit count before and after; case (a) also runs on the `singleHostOrigin` fixture of `internal/cli/gate_consume_sync_test.go` and compares the pushed state ref. (a) Successor: attempt 1 withdrawn, attempt 2 open, attempt 1's pair: exit 1, the refusal names attempt 1 (withdrawn) and attempt 2 (open). (b) Missing room: attempt 1 prepared and committed with summary A, the room directory removed, prepared again with summary B; the id and A's digest: exit 1. (c) Direct-file binding prepared again; the id and the old digest: exit 1. (d) Case (b) with the delegated shape: exit 1. In each refusal nothing changes and status stays the gate stage. The current pair closes and consumes. Falsifiers: comparing only the id makes (b), (c) and (d) exit 0 and consume; removing the compare makes all four exit 0.

**AC-2 - Every chat decision carries its binding, and a skill/binary mismatch stops before any write with a diagnostic.**
Verified by: (a) a table test in `internal/cli/gate_test.go`: the origin/main record shapes the old skills emit (captain approve with `--consume`, captain revise with a reason, delegated approve with a reason, `--conn-quote` and `--conn-source`), without binding flags and with only one of the two, exit 2 with the update diagnostic, also when the entity reference does not resolve; the entity bytes do not change. The same shapes with both flags close the attempt and store the reason and both grant fields. (b) The binding flags with `--round` exit 2; the existing round test publishes once with its unchanged command, and the existing re-entry test prepares a fresh attempt. (c) One-off at validation, the newer-skills pair: the new skills' `gate read --current` and `gate record --briefing-id` against the released `0.29.0-pre0` binary exit 2 with the unknown-flag error and leave the entity unchanged; one live Claude run with that binary on `PATH` stops with the update diagnostic, and its command log has no `gate record` without the binding flags. Falsifiers: making the flags optional lets (a) exit 0; checking after path resolution makes the unresolved-reference case exit 1; removing the stop rule lets the one-off run retry without the flags.

**AC-3 - Gate history stays readable whatever state a room is in.**
Verified by: the golden `internal/status/testdata/gate-history/complete.json` does not change, and a CLI test: `gate read --json` exits 0 with the stored records for an open attempt whose room directory is missing, an open attempt whose `index.json` was edited, a withdrawn attempt, a closed attempt with a pending application, an archived entity, and an entity with no gate history (`{"records":[],"version":1}`). Baseline: origin/main gives the same results (Spike 2, step 6). Falsifier: routing history through the verified presentation makes the missing-room and edited-room cases exit 1.

**AC-4 - `gate read --current` returns the decision binding and the verified review from the same bytes, and its failure leaves history readable.**
Verified by: a golden and a CLI test. A two-attempt entity (attempt 1 withdrawn, attempt 2 open) returns `stage`, attempt 2's id, `state: open`, a `briefing` equal to the stored binding, `room`, and `question` and `summary` equal to the stored Briefing. A direct-file binding returns its verified summary. Withdrawn returns `state: withdrawn` with no `summary`; closed with pending application returns `state: closed`; a consumed gate at the next stage, no record, and an archived entity return `attempt: null`. A missing room and an edited `index.json` exit 1 with no stdout, while `gate read --json` on the same entity exits 0. Falsifiers: skipping the digest check prints the edited summary; selecting the first attempt returns attempt 1.

**AC-5 - The room holds the full review the captain was shown, also after an open attempt from before this change (value; baseline: the mail-triage room held a 147-character summary of a review of about 20 lines).**
Verified by: phase 1 of a two-phase live Claude scenario. The fixture holds an open attempt with a short summary, prepared and committed by the binary, and a stage report that contains a seeded, decision-relevant finding with a unique phrase. The oracle requires that the short-summary attempt is withdrawn; that the successor's stored summary starts with `Gate review:`, has a `Decision:` line and contains the seeded finding; and that the first officer's presentation, with its `Reviewed snapshot:` line removed, equals the stored summary line for line (trailing spaces ignored). Deterministic support in `TestRecordedGateLifecycleAC7ResumeMatrix`: after an interrupted prepare, the resume commit and `gate read --current` return the summary and binding that prepare wrote; a published full review returns the same bytes on resume; a re-composed review on the open attempt is refused. Falsifiers: restoring the prepare-first order stores a short summary without the seeded finding; an oracle transcript that drops the finding from the presentation or from the stored summary fails; an oracle input whose presentation differs from the stored text fails.

**AC-6 - A relayed captain decision records where it came from, in a format that every binary reads.**
Verified by: phase 2 of the live scenario. A fresh first-officer process gets a turn-starting message from another agent that relays the captain's approval from a named pane, with phase 1's id and digest. The oracle requires a Resolution `by: person:captain` whose `briefing` is the relayed id, whose reason has a substantive line and ends with `Via: <pane>, relayed by <agent>.`, and a consumed application. One-off at validation: the released `0.29.0-pre0` binary runs `gate read --json` and `status --validate` on the phase-2 entity with exit 0. Falsifiers: an oracle input whose reason has no `Via:` line, or whose actor is `agent:first-officer`, fails; a structured `via` key makes the old-binary check exit 1 (Spike 1, step 5).

**AC-7 - When another agent's message starts the first officer's turn, the first officer's own output carries the gate review once, before any decision write.**
Scope: emission by the first officer. Delivery of that output to the captain belongs to the host session (section 6).
Verified by: the phase 1 and phase 2 oracles over the stream and the command log: exactly one assistant text event holds the review (entity, stage, Briefing id, recommendation, decision ask, seeded finding), and it comes before every `gate record` and `gate consume` in the command log. Oracle unit tests: the passing transcript with that event changed into a relay-only reply (the Briefing id and one line) fails; a second full copy fails; the event moved after the record command fails. Falsifier: changing the delivered event into a relay-only event.

## Test plan

- **AC-1 to AC-4 (deterministic, medium cost).** The proof owner is `internal/cli/gate_test.go`, which owns the gate record grammar, refusals and `gate read`, and has the `gatePrepareCLIFixture`/`gateOneFileCLIFixture` substrate. The Spike 2 fixtures (missing room, direct-file binding, successor) become the first refusal tests. Each refusal compares entity sha256 and the state commit count; the origin-backed case compares the pushed ref. The goldens own the `--current` shape; the edited-room case is the distinct falsifier for the digest check, and the missing-room history case is the distinct falsifier for keeping the two reads apart.
- **AC-5 deterministic support (low cost).** Extend `TestRecordedGateLifecycleAC7ResumeMatrix` in `internal/ensigncycle/recorded_gate_lifecycle_test.go` with the interrupted-prepare, published-review and re-composed-review cases. Existing prepare tests keep ownership of byte-exact storage and replay.
- **AC-5, AC-6, AC-7 (live, medium-high cost).** One new two-phase Claude scenario in `internal/ensigncycle`, built on the recorded-gate runner, its logging shim and command log, with oracle unit tests on synthetic transcripts for each falsifier. Phase 1 is an agent-started turn at a pre-upgrade open attempt, with no conn; phase 2 is a fresh process that receives the relayed decision. It runs once at the candidate tip. No existing proof owner covers agent-started turns or relayed decisions.
- **One-off at validation (low cost).** AC-2 (c) and the AC-6 old-binary read use the released `0.29.0-pre0` binary, an external artifact that a committed test cannot reproduce truthfully. The validation report records the commands and exit codes.
- **Existing suites and required lanes.** The record call sites gain the binding from prepare's output and keep their assertions. `go test ./...`, `go test ./... -race`, `gofmt` and `internal/contractlint` pass. The gate skills are host-neutral and every host's gate journeys load them, so under `docs/dev/README.md` the Claude, Codex and Pi live lanes are required green before merge: the no-conn stop, the headless conn continuation and the shared live journeys (`internal/ensigncycle/shared_live_runner_test.go`). The new two-phase scenario runs on Claude only.

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

## Stage Report: ideation (cycle 2)

- DONE: Revise the committed ideation design so each Material finding of the staff review at /tmp/spacedock-dispatch/staff-review-second-surface-gate-decisions.md is resolved in the body (FO disposition: fix findings 1 to 6 and 8; record finding 7 as a Deferred risk with its trigger and promote condition).
  Dispositions, one per finding:
  1. Same id can authorize different content — fixed: records require id and digest (Proposed approach 1); immutability rejected in Alternatives considered (withdraw refuses both shapes); AC-1 (b)-(d).
  2. Required flag changes the protocol both ways — fixed: Proposed approach 7 (one PR, minor release, edge prerelease after merge, both mixed pairs with diagnostic and recovery, `--round` unchanged); grammar ahead of path resolution; AC-2 (a)-(c); required Claude, Codex and Pi lanes in Test plan.
  3. Optional `via:` field breaks old readers — fixed: `Via:` line in the existing reason, no format change (Proposed approach 5); AC-6 asserts actor, preserved reason and old-binary read.
  4. Resume can present an old synopsis — fixed: present from verified stored text after commit (section 3); transition for prepared rooms, direct-file bindings and missing rooms (section 4); AC-5 phase 1 and resume-matrix cases.
  5. History would depend on a room file — fixed: `gate read --json` unchanged; `--current` is a separate mode that fails alone, with no entity path or full manifest (section 2); AC-3 and AC-4 split.
  6. AC-6 observed output, claimed delivery — fixed: narrowed to emission, delivery owner stated (section 6); now AC-7 with the relay-only falsifier.
  7. External rendering of the full summary — deferred: `## Deferred risk`, trigger and promote condition; AC-5 uses its seeded-finding oracle.
  8. Stage report and references — fixed: this report states the final contract; spike references corrected; the contract spec grammar is in the edit list.
- DONE: Exercise the new riskiest path: with a binary built from origin/main, show (or refute) a same-id replacement with a different digest (missing modern room, or legacy binding), and show the proposed refusal or immutability rule on it; record commands and exit codes.
  Risk evidence, Spike 2. Both replacements shown on origin/main: a missing room re-prepares to the same id with a new digest and a late approve consumes it (exit 0); a direct-file binding also rebinds (exit 0). Withdraw refuses both (exit 1). The id-and-digest prototype refuses both, including the delegated shape (exit 1, entity sha256 and commit count unchanged). A missing flag gives exit 2 before path resolution. The origin/main binary gives exit 2 on the new flags and reads a `Via:` reason. Spike binaries and fixtures removed.
- DONE: Update Acceptance criteria (exact heading, `**AC-N**` form, each with a falsifying edit, value ACs against baseline), test plan, expected surface and tolerance, declared semantic changes, and before/after wording; append a new stage report that lists each staff-review finding with its disposition and where the body resolves it.
  AC-1 to AC-7 under `## Acceptance criteria`; each names a falsifier. AC-1 and AC-5 measure value against the baseline. Net +675 across 19 files (tolerance +520 to +830, at most 23 files). Before/after wording covers the command reference, `gate --help`, the gate contract spec, `fo-gate-lifecycle` and `present-gate`.

### Summary

The revision binds each decision to the Briefing id and digest. The spike showed that an id alone authorizes replaced content on origin/main, for a missing room and for a direct-file binding. The digest check makes the existing in-place re-prepare safe, so prepare and withdraw do not change. History stays as it is, and the verified review moves to a separate `gate read --current` mode. Origin becomes a `Via:` line in the reason, so no stored format changes. The rollout section covers both mixed pairs. Decision points for CL: two required flags; the exit-code change for malformed records (exit 1 becomes exit 2); and the two-phase live scenario plus one-off old-binary checks, which raise the estimate from +360 to +675.

## Stage Report: implementation

- DONE: The approved cycle-2 design is implemented on this task's worktree branch, as written: `gate record --decision` requires `--briefing-id` and `--briefing-digest` and refuses a mismatch with no write (grammar checks before path resolution, exit 2; mismatch exit 1); `gate read --current --json` as specified while `gate read --json` stays byte-identical; the fo-gate-lifecycle and present-gate skill text, the gate contract spec and the command reference changed with the before/after wording in the task body; every existing record call site passes the binding. Skill smoke/contract tests are added or adjusted before the skill text changes.
  Binary, help, command reference and contract spec on `second-surface-gate-decisions` (690418b9e binary and AC-1 to AC-4 tests, 58400bbda resume matrix, 1a94d3161 docs, 42c356d43 live scenario and oracle, 0482563e1); `complete.json` golden unchanged. Every test record call site passes the binding (helpers `bound`, `withRecordedBinding`). The skill text is NOT on the task branch: fo-gate-lifecycle is 7699 bytes under a 7700-byte captain-approved cap (`TestFOInstructionComponentCaps`); the verbatim wording measures 9901, the pointer version 8873. On the FO's ruling, the pointer version (relay/Via/refusal and older-binary paragraphs moved to present-gate), the pi-runtime sequence line, the new contract test `TestGateSkillsRecordTheAnsweredBindingFromTheStoredReview` (fails before the skill edit, passes after) and a cap raise 7700 -> 9100 are one commit, 8ac8f459e, on side branch `second-surface-gate-decisions-skill-cap`. That cap change awaits CL. The library compares the pair only when one is passed; the CLI, the only production caller, requires both.
- DONE: AC-1 to AC-4 and the AC-5 deterministic support are committed tests that fail on origin/main b7ef2a885 and pass on the branch (show each fail-before/pass-after); the two-phase live Claude scenario for AC-5 to AC-7 and its oracle unit tests are committed, the oracle tests pass, and one local live run of the scenario is attempted and its grade or named unavailability recorded.
  Fail-before: the branch tests copied onto a b7ef2a885 worktree (plus a throwaway JSON type shim so they compile) all fail: AC-1 `TestGateRecordRefusesAPairThatIsNotTheCurrentBinding` (exit 2 `unknown gate flag: --briefing-id`; on main the stale pair cannot be named, and Spike 2 shows it consumes); AC-2 `TestGateRecordRequiresTheAnsweredBindingBeforeResolvingTheEntity` (bare shapes are accepted, grammar errors come after resolution); AC-3 `TestGateReadHistoryStaysReadableInEveryRoomState` fails only in the closed and archived setups that record with the binding; its four history cases pass on main, which is the declared unchanged baseline; AC-4 `TestGateReadCurrentReturnsTheAnsweredBindingAndVerifiedReview` (exit 2 `unknown gate read flag: --current`); AC-5 support, three new `TestRecordedGateLifecycleAC7ResumeMatrix` cases. All pass on the branch. Falsifiers: an id-only compare passes missing-room, direct-file and delegated stale pairs; resolving the path first turns the `nosuch` cases into exit 1; skipping the digest check prints the edited summary; selecting the first attempt breaks the `current-open.json` golden.
  Live: `TestLiveSecondSurfaceGateDecision` (live tag, 42c356d43) and `TestSecondSurfaceGateOracle` (relay-only reply, second copy, presented after record, dropped finding, text differs from stored, sub-agent-only text, no Via line, first-officer actor, other Briefing, not consumed: each fails the oracle; the passing transcript passes). One run on the side branch (90bbf81ca, the same skill bytes as 8ac8f459e): UNAVAILABLE, `authentication_failed: Not logged in · Please run /login` at the first event of phase 1, before any FO work; not repeated. Until the skill commit lands, AC-5 to AC-7 and AC-2 (c) cannot pass on the task branch alone.
- DONE: `go test ./... -timeout 30m`, `go test ./... -race -timeout 45m`, `gofmt -l ./cmd ./internal` and `go test ./internal/contractlint/` results are recorded, and actual net LOC and file count are recorded against the approved tolerance (net +520 to +830, at most 23 files).
  Task branch 0482563e1: `go test ./...` exit 0 (22 packages ok); `-race` exit 0 (22 ok); contractlint ok; `gofmt -l` lists only `internal/release/runtime_live_evidence_workflow_test.go`, unchanged from b7ef2a885. Side branch: contractlint, ensigncycle, journeymetrics, release ok. Size against b7ef2a885: task branch 18 files, +1267 -109, net +1158; with the side commit 23 files, +1348 -122, net +1226. Net is past the tolerance's +830 (+1158 is 172% of the +675 estimate); the file count is at the 23 limit. Production is close to the estimate (about +196 net vs +135); tests are the overage (about +955 net vs +515): the AC-1 to AC-4 matrix and its fixture helpers (+527) and the oracle with its unit tests (+225). A captain-visible design-reset decision on the size is needed before another pass.

### Summary

The binary binds each gate decision to the Briefing id and digest it answers and adds `gate read --current --json`, with tests for AC-1 to AC-4 and the resume support, the live two-phase scenario and its oracle; full and race suites pass. The skill text does not fit under fo-gate-lifecycle's byte cap, so on the FO's ruling it is a single commit on `second-surface-gate-decisions-skill-cap` (8ac8f459e) with a cap raise to 9100 that awaits CL. The one live run was unavailable (not logged in). The net size, +1158, is past the approved tolerance and needs a design-reset decision.

## Stage Report: validation

- FAILED: Each acceptance criterion's cited evidence is reproduced independently on throwaway exports, not taken from the report: AC-1 to AC-4 and the AC-5 resume support fail on origin/main b7ef2a885 and pass on the task branch head 0482563e1, and at least one named falsifying edit per AC makes its test fail; AC-5 to AC-7 and AC-2 (c) are checked on the side branch `second-surface-gate-decisions-skill-cap` (8ac8f459e) as far as possible without a live login, and each is marked proven, unproven, or blocked with the reason.
  Only the fail-before half ran. Free disk fell to 1.7 GiB, then 1535 MiB, and on the coordinator's order all builds and tests stopped. Fail-before: a detached b7ef2a885 worktree, the task branch's test files and goldens checked out over it, plus a type-only shim (`CurrentBinding`, `CurrentAttempt`, `CurrentPresentation` with their json tags, no functions) so the tests compile. Results:
  AC-1 `TestGateRecordRefusesAPairThatIsNotTheCurrentBinding`: all 5 subtests FAIL, exit 2 `unknown gate flag: --briefing-id`. This failure is grammar only. The proof of the value baseline (a late approve consumes the other content) is still Spike 2 step 1 from ideation. I did not reproduce it.
  AC-2 `TestGateRecordRequiresTheAnsweredBindingBeforeResolvingTheEntity`: all 4 subtests FAIL. The bare old-skill shapes exit 0 and close the gate, and the grammar-before-resolution case gets the unknown-flag error. This failure is substantive.
  AC-3 `TestGateReadHistoryStaysReadableInEveryRoomState`: 4 PASS (missing room, edited index.json, withdrawn, no history) and 2 FAIL (closed, archived), only because their setups record with the new flags. This is an invariance AC whose baseline is "unchanged", so fail-before does not apply to it.
  AC-4 `TestGateReadCurrentReturnsTheAnsweredBindingAndVerifiedReview`: all 7 subtests FAIL, exit 2 `unknown gate read flag: --current`.
  AC-5 support `TestRecordedGateLifecycleAC7ResumeMatrix`: the 3 new cases FAIL, exit 2 on `--current`. 3 existing cases also FAIL because their records now carry the flags. `open-same-package` PASS.
  Not run: pass-after on 0482563e1, all falsifying edits, and the side-branch tests (`TestGateSkillsRecordTheAnsweredBindingFromTheStoredReview`, `TestFOInstructionComponentCaps`, `TestSecondSurfaceGateOracle`).
  AC status: AC-1 to AC-4 and AC-5 support are unproven (fail-before only). AC-5, AC-6 and AC-7 live parts are blocked: `newClaudeLiveRunner` needs a benchmark OAuth token or `ANTHROPIC_API_KEY` against an empty HOME, and neither is available. The live run was not retried. The AC-2 (c) live Claude run is blocked for the same reason.
- FAILED: The one-off checks against the released 0.29.0-pre0 binary in AC-2 (c) and AC-6 are run where possible and recorded with commands and exit codes; `go test ./...`, `go test ./... -race`, `gofmt -l` and contractlint results are recorded for the task branch.
  Not run because of the disk stop. The binary is still installed at /opt/homebrew/bin/spacedock (`--version`: 0.29.0-pre0). There is no independent `go test ./...`, `-race`, gofmt or contractlint run. The only results are the implementation report's: full exit 0 (22 ok), race exit 0, contractlint ok, and `gofmt -l` listing only `internal/release/runtime_live_evidence_workflow_test.go`, the same as on b7ef2a885. These are cited, not reproduced. My own full-suite start never launched (nohup `nice` refused), so it used no disk.
- DONE: PASSED or REJECTED recommendation with findings classified per Review-finding disposition (material vs deferred), stating separately: the net size against the approved tolerance (+1158 vs +520..+830), the cap change awaiting CL, the unavailable live run, and the required Claude, Codex and Pi CI lanes that cannot run without a push.
  Recommendation: REJECTED-pending-reruns. The validation is incomplete. No candidate defect is shown in the binary half. Findings (proposals, for FO authorization):
  1. Pi `«gate.lifecycle»` binding is stale on 8ac8f459e (skill half). `skills/first-officer/references/pi-first-officer-runtime.md:15` still reads "After `gate prepare` returns `state=open`, run `state commit` once. Then run `status --read` ... emit one captain-facing review block from this gate's fresh evidence". The side commit updated only line 25 of that file. Four fields: (1) a Pi first officer at any gate after the skill half lands; (2) it can prepare before composing the review, store a short synopsis and present text the digest does not cover; (3) value-ac[AC-5]: the room holds the full review the captain was shown, which the Pi lane must also carry; (4) trigger evidence is the line itself, which Pi reads at every gate. It was not observed live. Proposed: Material, owned by this task, fix it on the skill half by rewriting line 15 to the steps 1-6 order. `fo-dispatch-core.md:221` (needs-preparation: "summary ... commit its binding, and re-read the same envelope") could carry the same pointer as polish.
  2. Polish: on 8ac8f459e, the fo-gate-lifecycle Resume paragraph puts the new sentences inside the route list ("Never compose a new review ...; `needs-preparation` → review; ..."). The later routes then read as part of that sentence.
  Size: net +1158 across 18 files on 0482563e1 (prod +248/-51, tests +1010/-57, docs +9/-1, numstat). The side commit adds 5 files, for net +1226 across 23 files. Both exceed the approved +520..+830, and 23 files is at the limit. This needs CL's design-reset decision. Validation cannot waive it.
  Cap: the fo-gate-lifecycle cap raise from 7700 to 9100 on 8ac8f459e awaits CL. It is not approved, and it is not merged here.
  Live: unavailable (no login, as above). AC-5 to AC-7 and AC-2 (c) have no live evidence.
  CI: the required Claude, Codex and Pi live lanes (no-conn stop, headless conn continuation, shared live journeys) cannot run without a push. That is forbidden here, so these lanes are not run.

### Summary

On b7ef2a885, the task-branch tests fail where they should for AC-1, AC-2, AC-4 and the resume cases. AC-3's history cases pass, which matches its unchanged baseline. A disk emergency (free space down to 1535 MiB) stopped all testing before the pass-after runs, the falsifiers, the side-branch tests and the old-binary one-offs, so those ACs are unproven or blocked. The recommendation is REJECTED-pending-reruns. A read-only review of 8ac8f459e found a stale Pi `«gate.lifecycle»` binding, proposed as Material for the skill half. Scratch was removed: the base worktree with `git worktree remove`, plus the private GOCACHE and the binaries.

Validation-findings note (FO disposition: fix both, side branch only): commit 2913ce0b1 on `second-surface-gate-decisions-skill-cap` (parent 8ac8f459e). It rewrites `pi-first-officer-runtime.md` line 15 to the steps 1-6 order, and it moves the pre-change resume sentences out of fo-gate-lifecycle's route list. fo-gate-lifecycle now measures 8880 bytes, which is under the proposed 9100 cap; that cap still awaits CL. `go test ./internal/contractlint/ -run 'TestGateSkillsRecordTheAnsweredBindingFromTheStoredReview|TestFOInstructionComponentCaps' -count=1`: ok, with 3.3 GiB free before the run. The task branch is unchanged at 0482563e1.
