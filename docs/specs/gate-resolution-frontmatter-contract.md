# Gate Resolution frontmatter contract

Status: implemented recorder and application contract
Date: 2026-07-22

## Outcome and ownership

The recorder makes a captain's decision durable before any workflow status change or
dispatch. It owns the logical gate, ordered attempts, immutable Briefing binding, and
portable Resolution. A closed approval carries the typed one-use `application`
subtree; revise and hold are complete Resolutions with no application. The same
recorder binary owns its guarded writes.

Presentation remains an overridable channel of the present-gate skill, not a recorder
verb. Chat decisions use `gate record --decision`. A selected override receives only
the opaque, committed prepared room; a room-backed outcome uses `gate record --room`.

## End-to-end gate lifecycle

Chat and an override are alternative presentation channels. Both begin with the same
prepared, durably committed gate room and converge on the same recorder-owned closed
attempt. The First Officer chooses review content; Spacedock prepares and records
authority. The generic contract treats an override as an opaque room handoff.

```mermaid
flowchart TD
    FO["First Officer<br/>selects Artifact and References<br/>authors question and summary"]
    PREP["spacedock gate prepare<br/>derives IDs, digests, Git locators,<br/>authority, room, and binding"]
    ROOM[("Frozen gate room<br/>request.json and canonical Briefing")]
    COMMIT_PREP["spacedock state commit<br/>publishes the prepared binding"]
    CHANNEL{"Presentation channel"}

    FO --> PREP --> ROOM --> COMMIT_PREP --> CHANNEL

    subgraph CHAT["Default chat"]
        CHAT_REVIEW["First Officer presents<br/>the canonical Briefing"]
        CHAT_DECISION["Captain decides"]
        CHAT_RECORD["spacedock gate record --decision<br/>--actor person:captain"]
        CHAT_REVIEW --> CHAT_DECISION --> CHAT_RECORD
    end

    subgraph OVERRIDE["Selected presentation override"]
        HANDOFF["Override receives only<br/>the opaque committed room"]
        ROOM_RECORD["spacedock gate record --room<br/>verifies room-backed authority"]
        HANDOFF --> ROOM_RECORD
    end

    CHANNEL --> CHAT_REVIEW
    CHANNEL --> HANDOFF

    CHAT_RECORD --> CLOSED["Recorder closes the gate attempt"]
    ROOM_RECORD --> CLOSED
    CLOSED --> COMMIT_CLOSE["spacedock state commit<br/>publishes the Resolution"]
    COMMIT_CLOSE --> CONSUME["spacedock gate consume"]
    CONSUME --> COMMIT_CONSUME["spacedock state commit<br/>publishes application"]
    COMMIT_CONSUME --> NEXT["Successor stage"]
```

The generic lifecycle ends at that opaque handoff. It defines no override execution,
transport, presentation, or evidence-production mechanics. If the room later becomes
recorder-ready, `gate record --room` validates its canonical v1 Result and complete
inventory under the recorder contract below.

## Canonical v1 schema

The binary accepts and emits one canonical `gates:` shape:

```yaml
gates:
  version: 1
  records:
    - id: gate:example:sample:validation
      stage: validation
      attempts:
        - id: gate-attempt:sample-validation-1
          briefing:
            id: briefing:sample-validation-1a
            digest: sha256:3333333333333333333333333333333333333333333333333333333333333333
            request-digest: sha256:4444444444444444444444444444444444444444444444444444444444444444
            room-ref: ./review/validation/briefing-1
          provider-evidence:
            result-digest: sha256:5555555555555555555555555555555555555555555555555555555555555555
            presented-inventory-digest: sha256:6666666666666666666666666666666666666666666666666666666666666666
          resolution:
            type: Resolution
            id: resolution:captain-sample-validation-1a
            briefing: briefing:sample-validation-1a
            by: person:captain
            at: 2026-07-22T09:00:00Z
            decision: approve
          application:
            target-stage: done
            state: pending
```

`records` and `attempts` are ordered. The last attempt in a record is current. An
attempt is open when both `withdrawal` and `resolution` are absent, withdrawn when only
`withdrawal` is present, and closed when only `resolution` is present. Withdrawn and
closed attempts are frozen. Entity `status` selects exactly one record by its `stage`;
duplicate stage records fail closed. The last ordered attempt in that record is current
and eligible for later application. These facts remove any need for separate attempt pointers, sequence
numbers, lineage pointers, or explicit lifecycle state.

The binary-owned model is closed: unsupported fields inside `gates` fail validation.
In particular, the pilot-only attempt selector, `current-attempt`, `sequence`,
`previous-attempt`, and explicit attempt `state` encodings are rejected. There is no
migration or compatibility rewrite. The `application` field is an approval-only
authority token with exactly `target-stage` and `state`, where state is `pending`,
`consumed`, or `superseded`. Revise and hold carry no application.

Every Briefing binding includes an id, canonical SHA-256 digest, and an exact file or room
reference. Version 1 digests are unconditionally RFC 8785/JCS canonical Briefing JSON
bytes. A request-backed room additionally freezes its
request digest; that request names the canonical Briefing with a clean room-relative
locator, id, and digest. No reader infers a canonical basename.

A prepared provider-neutral room binds `request-digest`, the JCS digest of its
`request.json`. Request-less and chat-only attempts may omit it. Changing the request,
located Briefing, or a selected Git object after the attempt binds therefore fails
before Result validation or entity mutation.

## Provider-neutral preparation

`gate prepare` is the single mechanical operation that turns First Officer judgment and
committed file selections into an open recorder-ready attempt:

```text
spacedock gate prepare ENTITY --question TEXT --artifact REVIEW.md --summary TEXT \
  [--reference FILE ...] [--workflow-dir DIR]
```

The caller supplies exactly one question, Markdown primary Artifact, and nonblank
valid-UTF-8 primary summary; References may repeat in caller order. Spacedock preserves
the summary string exactly, assigns deterministic ordinal item identities, and derives
the gate, attempt, Briefing, Captain authority, digests, and room. It writes only
`gate-briefing.json` and `request.json` at preparation time. It copies no selected
source, writes no association, and creates no provider subtree.

Folder and flat entities share the same companion-room layout:

```text
<state-root>/<slug>/review/<stage>/briefing-<attempt>/
```

For folder form, the entity is `<slug>/index.md`; for flat form it is `<slug>.md` and
`<slug>/` is its artifact companion. State commit and archive operations treat the flat
Markdown plus companion directory as one literal path-scoped unit, including tracked
deletions and rollback, without sweeping siblings.

Each selected source is a readable, committed, non-symlink regular file owned by the
workflow's `main` or distinct `state` Git history. Its closed identity is
`git-root://<main|state>/<full-commit>/<canonically-escaped-repository-path>` and its
`rev` is the full raw-byte SHA-256. Preparation compares the selected worktree bytes
with `<commit>:<path>`; later operations reopen only that local Git object. They never
fetch, deepen, hydrate, retain a ref, search another checkout, or fall back to current
worktree bytes. A clean detached or linked worktree is valid when it shares the
expected Git history and the object is local.

The primary Artifact alone carries the exact caller-supplied `summary`. References
carry none. Advisory-round Briefings remain summary-free. Exact prepare replay is a
no-op, divergent occupancy fails closed, and
handled validation, publish, or bind failure removes the new candidate and any
newly-created empty parents. Success prints exactly `room`, `briefing`, `digest`, and
`state=open` lines; the emitted absolute room is the only later handoff coordinate.

## Recorder lifecycle

`spacedock gate prepare` is the First Officer's normal lifecycle entry. With no record
for the current stage it opens the first attempt. Exact replay of the current open room
is a no-op; divergent open occupancy fails closed. After a withdrawn or closed attempt,
preparation appends a successor while earlier authority remains frozen; only a closed
attempt may have a pending application to supersede.

`spacedock gate withdraw ENTITY --reason TEXT` retires only the selected current-stage
open request-backed attempt. Under the shared lock it validates all retained authority
and requires the room to contain exactly `gate-briefing.json` and `request.json`. It
records only `withdrawal: {by: agent:first-officer, at: <UTC>, reason: <TEXT>}`: no
Resolution, provider evidence, application, status change, successor, or room write.

`spacedock gate record` also accepts either the prepared room's room-backed Result or
a semantic chat decision. Either closing source closes only the last open
attempt for the current stage. Approve derives one `application` with
`target-stage` and `state: pending`; revise and hold write no application.

Before either ordinary close, the recorder resolves authoritative current status in the
workflow taxonomy and requires a nonterminal `gate: true` stage. The bound Briefing must use the canonical v1 stage-qualified identity and name that same stage. Malformed identity, mismatch, or non-actionable stage fails before Resolution construction and leaves entity bytes unchanged.

Cross-logical-gate re-entry is ordinary: workflow stage selects the target record even
when another stage's closed gate is retained in history. The successful write selects the
target record but does not modify either record's earlier closures.

## Room-backed Result association

The room-backed form consumes one prepared gate room. Its frozen `request.json` binds the
logical gate, attempt, canonical Briefing id and digest, and captain actor/approver authority.
Delegated chat decisions record `agent:first-officer` with a nonblank evidence reason and no directive or `adoption-note`;
the attempt's `request-digest` rejects post-binding request changes.
The fixed recorder inputs are `provider/result.json` and
`provider/presented-inventory.json`; callers supply neither path nor invocation detail.

The recorder validates `request.json`, resolves its exact frozen Briefing locator,
recomputes its JCS digest, and
derives the canonical inventory from every Artifact and recursively reached Reference.
It derives a private `spacedock-result-association` v1 by matching each presented id and
revision to that inventory and binding the raw Result digest. The mapping must cover the
whole inventory exactly once, including the Result's primary Artifact.

A direct binding Result uses Review v1's minimal envelope: authority comes from nested
`Resolution.by`, and redundant `status`, `binding`, `actor`, `approver`, or
`resolutionId` fields are absent. The recorder requires `Resolution.by` to equal the
request authority. Advisory output remains retained evidence; no adoption note can
promote it into a binding Resolution. Artifact payloads may remain external URI and
SHA references.

On room-backed close, the recorder stores only the raw-byte digests of
`provider/result.json` and `provider/presented-inventory.json` as `provider-evidence`.
They are part of the frozen attempt. `gate validate` recomputes both from the fixed room
files and fails if either is missing or changed. Chat-closed and open attempts carry no
provider evidence; the derived association remains ephemeral.

Request, located Briefing, Result, and presented inventory all pass through one
recursive token-stream duplicate-member check before typed decoding or
canonicalization. Conflicting members at any object depth fail closed; Go's
last-member-wins JSON behavior is never authority. Binding, room recording,
validation, eligibility, and consumption all resolve and recheck the same frozen
request/Briefing/source authority. Room-backed evidence is valid only on a prepared,
request-digest-bound attempt.

## Round records (workflow-neutral correction evidence)

A correction round reuses the recorder vocabulary without becoming a gate: its reviewed
snapshot is an immutable, digest-bound Briefing; same-Briefing Annotations and
Resolutions are retained in an ordered review log. Round entries carry no application,
do not select a logical gate, and cannot advance workflow status. Finding labels,
materiality, disposition, ownership, and any workflow prose are opaque to the generic
recorder and remain the responsibility of the active workflow and First Officer.

`spacedock gate record <entity> --round STAGE/CYCLE --briefing PATH/briefing.json
--log PATH/briefing.review.jsonl` publishes the canonical two-file room at
`review/<stage>/round-<cycle>`, then atomically writes the exact `review-round` pointer.
The producer does not parse or write a `### Feedback Cycles` section. A workflow may
append its authorized Cycle line before invoking the producer; the recorder preserves
that body byte-for-byte. `spacedock gate validate <entity> --round STAGE/CYCLE`
replays the pointer and reports every Resolution as advisory structural evidence.

Round recording requires a folder-form entity at `<slug>/index.md`, so its accumulating
`review/` artifacts are scoped beside that entity. Flat entities refuse before locking
or writing; the recorder does not alter the approved derived room path to compensate.
`STAGE` must name a stage in the workflow definition, but need not equal current
`status`: explicit historical backfill remains supported.

The room is immutable: exact whole-room replay is a whole-tree no-op; any different
Briefing, log, room shape, or pointer fails closed. New-room publication rolls back if
the full-entity compare-and-swap or atomic pointer replacement fails.

Narrowing a value AC to make a finding pass is not a round disposition. It opens a real
gate attempt whose binding Resolution is captain-owned; the correction loop cannot
self-approve that design reset.

## Write boundary and invariants

Ordinary recorder operations rebuild only the canonical `gates:` subtree. Consumption of a
non-terminal target co-writes `status` and `application.state: consumed` in one atomic
replacement. Consumption of a terminal-target approval is mechanism-agnostic routing,
not a spend: it writes nothing, leaves the application `pending` and the status at the
gated stage, and returns the `approved-awaiting-merge` route; a repeated consume is
idempotent re-routing. The terminal merge ceremony (`spacedock merge guard`) is the
sole terminal consumer: with delivery proof it clears the `mod-block` in its
own step, then writes, in one locked replacement, `application.state:
pending→consumed` plus the terminal status, `verdict`, and `completed` — the
`pr` merge sentinel is retained through archive as durable delivery proof —
and a non-forced `status --set` to a terminal stage is refused while a
pending terminal-target application is in force. `merge guard --rework` writes
`application.state: pending→superseded` with `status :=` the record stage's declared
`feedback-to` and delivery state cleared — through the same guarded application
mutation; it refuses when no pending terminal-target application exists or when the
declared `feedback-to` is missing, undefined, or terminal. Before replacement it
validates the rebuilt full entity and compares the locked source subtree with the one it
read, so stale or invalid writes fail without replacing the file. All frontmatter fields
outside `gates` and the Markdown body are preserved byte-for-byte. The per-entity lock
rejects concurrent recorder writers; there is no retry, lease, daemon, or recovery
protocol.

The model enforces unique gate, stage, attempt, Briefing, and Resolution ids; a unique
status-matched logical gate when one exists; non-empty attempt histories; exact Resolution-to-Briefing binding;
and portable `approve`, `revise`, or `hold` decisions. `revise` and `hold` require a
reason or an included same-Briefing Annotation. Withdrawals require fixed
`agent:first-officer` attribution, a UTC timestamp, a nonblank reason, and a valid
request digest. Open attempts cannot carry provider evidence or application data;
withdrawn attempts can carry neither Resolution, provider evidence, nor application.

## Command surface

```text
spacedock gate prepare ENTITY --question TEXT --artifact REVIEW.md --summary TEXT [--reference FILE ...] [--workflow-dir DIR]
spacedock gate withdraw ENTITY --reason TEXT [--workflow-dir DIR]
spacedock gate record ENTITY --room PATH [--workflow-dir DIR]
spacedock gate record ENTITY --decision approve|revise|hold --actor ID [--reason TEXT] [--workflow-dir DIR]
spacedock gate validate ENTITY [--workflow-dir DIR]
spacedock gate eligibility ENTITY [--workflow-dir DIR]
spacedock gate consume ENTITY [--workflow-dir DIR]
```

Exactly one semantic source is required. Supported chat actor IDs are `person:captain` and `agent:first-officer`. The binary derives operation, ids, stage target,
and compare-and-swap state; callers cannot submit an operation envelope or candidate
identities. `gate validate` is read-only and reports the selected record's last attempt.

New delegated chat resolutions use `by: agent:first-officer`, require a nonblank
evidence reason, and reject `adoption-note` as an unknown prototype field. The recorder
constructs the portable Resolution under the asserted identity that rendered the
decision; it does not authenticate chat or apply the result.

## Explicitly outside v1

- Prototype-format compatibility, migration, and arbitrary
  unknown-field preservation inside `gates`.
- Presentation-channel execution, transport, UI, and evidence production beyond the
  opaque prepared-room handoff.
- Remote Git-object acquisition, retention refs, copied selected-source payloads, or
  generic URI/root registries.
- Blocker-satisfaction evaluation, execution-hold authoring, dispatch identities, or effect receipts.
- A second schema version or provider operation envelope.

## Behavioral proof

The release tests must fail if any of these outcomes regress:

1. Open/rebind/close/successor behavior changes, or a successor mutates a frozen closure.
2. Cross-gate re-entry follows current workflow stage rather than historical records.
3. A prototype field or arbitrary unknown binary-owned field becomes readable or writable.
4. A stale, invalid, or lock-contended write changes the entity.
5. A canonical write changes bytes outside `gates` or alters opaque application data.
6. Removing canonical artifacts and matching presentation entries together is accepted.
7. A room Result closes before exact bytes, request authority, bound Briefing digest,
   complete Artifact/Reference inventory, and full presentation mapping are verified.
