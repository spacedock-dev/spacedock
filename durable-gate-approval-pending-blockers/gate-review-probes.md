# Ask once, verify the next revision

Status: companion proposal for `3k` ideation review
Date: 2026-07-18

## Product promise

> Ask once. The next revision shows whether your concern was addressed.

The first version should reach that payoff before it offers configuration. A user asks
one useful question, sends the work back, and sees the same question answered against
the next revision. The product shows the old answer, the new answer, and the evidence
that changed.

## First-use experience

The ensign finishes the work and publishes the decision surface. The user sees:

```text
Ready for your decision

What do you need to know before deciding?
[Does this change preserve approval while dispatch is blocked?          ]

[Ask]
```

The field accepts the suggested question or free-form text. The first screen offers
one question, not a review setup form.

After the user selects **Ask**, the surface says `Checking this revision…`. Subspace
starts a fresh responder and returns a short attributed answer with citations:

```text
Yes. The change records approval before it checks dispatch blockers.

Evidence
- gate-resolution-frontmatter-contract.md:296
- gate-resolution-frontmatter-contract.md:306

Answered by Codex · GPT-5.5

[Approve]  [Send back with this concern]
```

If the evidence cannot answer the question, the result says `Insufficient evidence`
and names what it could not verify. It never fills the gap with an unsupported answer.

When the user selects **Send back with this concern**, the surface carries the question
and answer into a concise concern:

```text
Send back with this concern

[The approval is recorded, but the evidence does not prove that dispatch stays blocked.]

[Send back]
```

The ensign revises the work and publishes it again. The surface automatically asks the
same question against the new revision and shows:

```text
Your concern was checked again

Old answer
The evidence does not prove that dispatch stays blocked.

New answer
The scheduler now checks the committed blocker before it consumes approval.

What changed
The new contract adds the blocker check and a test that fails if dispatch occurs early.

Evidence
- gate-resolution-frontmatter-contract.md:304
- index.md:301

[Approve]  [Send back with this concern]
```

The user selects **Approve**. Spacedock first records the approval and shows `Approval
recorded. Spacedock will advance when safe.` The First Officer then consumes that
recorded approval and advances, routes, or dispatches as the workflow allows.

Only after the user sees the old/new comparison does the surface offer:

```text
Use this question in future reviews?  [Save]
```

The first-run UI never exposes Probe, Briefing, attempt, room, or Resolution vocabulary.

## Minimal internal model

A **probe** is exactly one reusable question. It is never a multi-question lens. The
first version stores probes as overlays in a durable Git-backed Subspace room. It does
not add them to portable Review & Gate `Context`. A portable probe format may follow
after the room model proves useful.

A room-local probe has a stable id, exact question text, revision, creation source, and
carry policy. Changing the question creates a new probe revision. The room keeps prior
revisions; Spacedock entity frontmatter does not copy them.

A probe result binds all of these values:

- `briefing`: the exact Review & Gate Briefing id and Spacedock JCS digest;
- `probe`: the exact probe id, question text, and probe revision;
- `requested-by`: the authenticated principal, harness, and model;
- `answered-by`: the responder principal, fresh harness run, and model version;
- `outcome`: a short answer or `insufficient-evidence`;
- `citations`: verified references to the reviewed artifacts; and
- `id` and `recorded-at`: the immutable result identity and time.

The room derives an old/new delta from two immutable results for the same probe lineage.
The delta names both result ids and presents the old answer, new answer, changed claim,
and supporting citations. It never substitutes for either source result.

## Gate and revision lifecycle

A logical Spacedock gate contains adjudication attempts. An open attempt points to one
current immutable Review & Gate Briefing. Enriching the presentation or revising its
inputs advances that pointer inside the same open attempt. State Git preserves earlier
pointer and digest values.

A binding portable Resolution closes the attempt. When the user sends work back, the
authenticated `revise` Resolution closes that attempt; the First Officer later consumes
its feedback application. Publishing the gate after rework creates a new attempt with a
new current Briefing.

The room automatically carries every previously answered probe into the next attempt
and runs it against the new current Briefing. A preset probe that merely appears
applicable stays visible and lazy until someone asks it. The first slice does not expose
preset matching or applicability language.

The current Spacedock entity stores the stable room reference, logical gates, distinct
attempts, each attempt's current or resolved Briefing id and digest, the exact binding
Resolution, and application state. It stores no full Briefing, log, probe, result, or
delta history. State Git preserves entity pointer revisions; the room preserves review
history.

## Ownership and authority

The ensign publishes and presents the gate. It may transport an authenticated
Resolution through a narrowly guarded recorder. The recorder checks the entity, gate,
attempt, exact current Briefing id and digest, provider attribution, and externally
authorized approver identity. The ensign cannot assert captain authority, create an
attributed captain decision, consume the application, or transition workflow state.

Subspace resolves and verifies artifacts, invokes a fresh responder for the first
version, records responder attribution, and persists the Git-backed room. It also
authenticates the acting user and stamps the binding Resolution. Workflow tooling still
supplies the authorized approver identity; Subspace does not choose who may bind the
gate.

The First Officer validates the committed Resolution against the current entity and
workflow authority. Only the First Officer consumes its application to transition,
route, or dispatch.

Review & Gate v1 remains the portable authority. The complete contract is
`../spacedock-subspace/docs/review-and-gate.md` at `spacedock-subspace` commit
`bd17bdb23318f815d17a1d10ea2a6d39ab449520`, blob
`14f3eb91ec85bfcc08bb3330c21b94cc77f4529f`. Each immutable Briefing owns one ordered
log. The first Resolution attributed to the externally authorized approver is binding,
and no later entry may follow it. These room-local probes and results do not alter that
contract.

## Definition storage after the first slice

Project or team probe definitions may later live in shared Git configuration. Personal
definitions may live in a user registry. Contextual preset matching, sharing, scope
selection, and promotion into a portable format remain deferred.

Saving a question after its first successful comparison can copy its definition to one
of those stores. The room results, citations, responder identities, and deltas remain
attached to the reviewed room and exact Briefings.

## Behavioral acceptance scenario

1. An ensign publishes a gate. The user sees `Ready for your decision` and one editable
   question.
2. The user asks the question. A fresh attributed responder returns a short result with
   verified citations bound to the exact current Briefing and question revision.
3. The user selects `Send back with this concern`. The authenticated binding decision
   closes the attempt, but recording it does not change workflow state.
4. The First Officer consumes the feedback application and routes the work for rework.
5. The ensign revises the work and publishes a new gate attempt with a new Briefing.
6. The room automatically re-runs the same question. The UI shows the old answer, new
   answer, and an evidence-backed delta. The result cites only the new Briefing's
   verified artifacts; the comparison links both immutable results.
7. The user approves. The recorder commits the exact authenticated approval before any
   workflow transition or dispatch.
8. A fresh read proves the approval is durable. The First Officer validates and
   consumes it exactly once, then advances the workflow.

The scenario fails if the responder reuses an unattributed prior run, answers without
evidence, evaluates the old Briefing, changes the question silently, hides the old/new
comparison, lets the ensign assert captain authority, or advances before the approval
commit.

## First-version cuts

The first version has no exposed probe management, scopes, applicability language,
lens collections, multi-probe synthesis, or portable probe format. It also omits
contextual preset matching and sharing. One question, one cited answer, one rework
comparison, and one durable decision deliver the first useful loop.
