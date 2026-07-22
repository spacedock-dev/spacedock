# Remember resolved concerns across a changing spec

Status: companion proposal for `3k` ideation review
Date: 2026-07-18

## The problem

The captain described the hard case:

> “but what's more tricky is for a particular instance of spec that is evolving through conversation. I'd ask a bunch of questions to confirm capability or details, but sometimes a resolved concern drifts with later changes, and right now it is still easier to re-read the spec to confirm.”

The core value is concern memory. A person asks one question about one revision, gets an
evidence-backed answer, and sees when a later revision changes or weakens that answer.
The person should not need to re-read the whole spec.

> Ask once. The next revision shows whether your concern was addressed.

A Probe may cause an author to improve the spec. Automatic incorporation is useful but
optional. The product succeeds when it preserves and rechecks the resolved concern.

## Before presentation: publication checks

General project or team questions run before presentation. The publisher also searches
the candidate for contradictions that the current change introduced.

If a check reveals a contradiction or an obvious self-revision, the publisher returns
the finding and citations to the author. It does not present avoidable inconsistency to
the reviewer. If the evidence reveals a genuine human choice, the publisher presents
that choice explicitly and leaves it unresolved.

Personal common questions may come from a skill or profile. A personal question remains
personal unless its owner deliberately shares it. These checks improve a candidate
before review; they do not replace the instance-specific concern memory below.

## During and after presentation: one persistent question

An instance-specific **Probe** is exactly one question attached to a Subspace review
room and spec lineage. It persists across later Review & Gate Briefings in that lineage.
It is not a multi-question lens and is not portable Review & Gate `Context`.

The first useful interaction stays small:

```text
Ready for review

What do you need to confirm?
[Using the concrete YAML, trace 3k from a revised ideation gate through a second
 ideation attempt and then two validation attempts, and show where each resolved
 Briefing, Resolution, and workflow application remains durable.]

[Ask]
```

Subspace answers against the exact presented Briefing. A later Briefing in the same
lineage automatically triggers a new result for every previously answered Probe.
Questions that merely match a general or personal preset stay lazy until someone asks
them for this lineage.

## Probe and ProbeResult

A Probe record has a stable id, spec-lineage id, exact question, revision, creator, and
time. Editing the question appends a new Probe revision; it never rewrites the prior
question.

Each immutable **ProbeResult** binds:

- the exact Probe id, revision, and question;
- the exact Briefing id and digest;
- `requested-by` principal, harness, and model;
- `answered-by` responder, fresh harness run, and model version;
- an answer or `insufficient-evidence` outcome;
- verified citations or evidence;
- explicit limitations; and
- an immutable result id and recorded time.

A ProbeResult reports what the evidence supports. It contains no recommendation,
decision, binding flag, advisory Resolution, or gate verdict. It never claims that the
reviewer should approve, reject, or revise.

## Provider-owned serialization

The review provider owns Probes and ProbeResults. A room may serialize them as an
append-only `probes.jsonl`:

```json
{"type":"Probe","id":"probe:3k-durability-trace","revision":1,"spec-lineage":"spec:3k-gate-design","question":"Using the concrete YAML, trace 3k from a revised ideation gate through a second ideation attempt and then two validation attempts, and show where each resolved Briefing, Resolution, and workflow application remains durable.","created-by":"person:captain","at":"2026-07-18T12:00:00Z"}
{"type":"ProbeResult","id":"probe-result:3k-durability-trace:b1","probe":{"id":"probe:3k-durability-trace","revision":1},"question":"Using the concrete YAML, trace 3k from a revised ideation gate through a second ideation attempt and then two validation attempts, and show where each resolved Briefing, Resolution, and workflow application remains durable.","briefing":{"id":"briefing:3k-design-b1","digest":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"outcome":"answered","answer":"The YAML preserves four closed attempts and keeps each resolved Briefing and Resolution beside its workflow application.","citations":["gate-resolution-frontmatter-contract.md#physical-representation"],"limitations":["Design trace only; no binary behavior was executed."],"requested-by":{"principal":"person:captain","harness":"subspace-web","model":"human"},"answered-by":{"principal":"agent:reviewer-1","harness":"codex","model":"gpt-5.5","run":"run:fresh-1"},"at":"2026-07-18T12:01:00Z"}
```

An equivalent provider store may use another durable format if it preserves append-only
record identity, immutable revisions, exact bindings, and replay. Git is one backend;
Git commits and paths are not Probe semantics.

Comparisons are derived from immutable ProbeResults. A provider may cache a comparison,
but it can always rebuild it and must not treat the cache as authority.

## Concrete ProbeResult

For the question above, a short narrative result could read:

```text
Answer

The YAML preserves four closed attempts under two logical gates.

1. Ideation attempt 1 freezes Briefing 1a beside a revise Resolution. Its consumed
   feedback application points to backlog.
2. Ideation attempt 2 freezes Briefing 2b beside an approve Resolution. Its advance
   application is consumed and points to implementation.
3. Validation attempt 1 freezes Briefing 1a beside a revise Resolution. Its consumed
   feedback application points back to implementation.
4. Validation attempt 2 freezes Briefing 2b beside an approve Resolution. Its advance
   application points to done but remains pending under an unsatisfied blocker and an
   active execution hold.

Each attempt keeps its resolved Briefing reference and exact Resolution next to its
workflow application. The current entity keeps those bindings; state history keeps
earlier open-pointer revisions; the review provider keeps the full Briefings and logs.

Evidence
- gate-resolution-frontmatter-contract.md, “Physical representation”
- gate-resolution-frontmatter-contract.md, “Recording, closing, and consuming are separate operations”

Limitations
- This result traces the proposed YAML. It does not prove parser, scheduler, or dispatch behavior.

Answered by Codex · GPT-5.5 · fresh run
```

This result answers the question without making a gate recommendation.

## Later-revision comparison

For each later Briefing, the provider runs the same Probe revision and derives one of
four relationships to the earlier result:

- `still-holds`: the later evidence supports the same answer;
- `changed`: the supported answer differs;
- `no-longer-supported`: the later evidence cannot support the earlier answer; or
- `not-affected`: the changed material does not bear on the question.

The review surface shows `changed` and `no-longer-supported` first. It keeps
`still-holds` and `not-affected` quiet unless the user asks to see unchanged checks.

If a later Briefing changes validation attempt 2 from pending to consumed, the derived
comparison could read:

```text
Changed

Old answer
Validation attempt 2 keeps its advance application pending under an unsatisfied blocker
and an active execution hold.

New answer
Validation attempt 2 still preserves Briefing 2b and its approve Resolution, but its
advance application is now consumed after the blocker cleared and the hold was released.

Evidence-backed delta
- application.state: pending → consumed
- execution-hold.state: active → released
- consumed-at: absent → 2026-07-18T13:10:00Z

Still holds (available on demand)
- Both ideation attempts and validation attempt 1 remain durable and unchanged.
```

The comparison links the two source ProbeResults. It does not edit either result or
automatically rewrite the spec.

## Ordinary Review & Gate flow

This design works with an ordinary Review & Gate Briefing. It needs no Spacedock repo,
gate, stage, attempt, entity, or First Officer.

1. A publisher runs publication checks and sends contradictions back to the author.
2. The publisher presents a clean immutable Briefing in a provider room.
3. A person asks one instance-specific question; the provider appends its Probe and
   ProbeResult.
4. Conversation produces another immutable Briefing in the same spec lineage.
5. The provider re-runs answered Probes and surfaces changed or unsupported answers.
6. Any human decision remains a separate Review & Gate Resolution.

Review & Gate v1 remains the portable authority. The complete contract is
`../spacedock-subspace/docs/review-and-gate.md` at `spacedock-subspace` commit
`bd17bdb23318f815d17a1d10ea2a6d39ab449520`, blob
`14f3eb91ec85bfcc08bb3330c21b94cc77f4529f`. Probe and ProbeResult are provider-owned
room records, not portable Briefings, Context, log entries, Annotations, or Resolutions.

## Optional Spacedock adapter

Spacedock may publish a Review & Gate Briefing and retain an opaque provider resume
reference. It does not need to understand Probe storage, reruns, or comparison.

If the review produces a binding Resolution, Spacedock may adopt the exact authenticated
object into entity state. Recording remains separate from consumption; the First
Officer later consumes the workflow application to advance, route, or dispatch.

The durable encoding in
[`gate-resolution-frontmatter-contract.md`](gate-resolution-frontmatter-contract.md)
already keeps full Probe and ProbeResult history out of entity frontmatter. One
contradiction remains: its concrete YAML uses a different `room-ref` for each attempt,
but this proposal requires one review-room/spec-lineage identity across Briefings and
across any Spacedock attempts. A follow-up must either reuse one opaque lineage-level
`room-ref` across those attempts or add a separate lineage reference. This companion
does not choose or edit that Spacedock encoding.

## Behavioral acceptance scenario

1. Run shared questions and dynamic contradiction checks against candidate Briefing A.
   Return an obvious contradiction to the author; present a genuine human choice.
2. Publish corrected Briefing B in a provider room with no Spacedock metadata.
3. Ask the concrete 3k durability question once. Append one Probe and one ProbeResult
   bound to the exact question revision and Briefing B id/digest.
4. Publish later Briefing C in the same spec lineage. Append a fresh attributed
   ProbeResult for the same Probe revision.
5. Derive and show the old answer, new answer, citations, limitations, and
   evidence-backed relationship. Put `changed` or `no-longer-supported` first; keep
   unchanged results on demand.
6. Restart from the provider store and reproduce both immutable results and the same
   comparison without re-reading the whole spec or consulting transcript prose.
7. If a human records a binding Resolution, keep it separate from ProbeResult. With the
   optional Spacedock adapter, commit that Resolution before the First Officer consumes
   it.

The scenario fails if a result lacks exact question/Briefing binding, responder
attribution, citations, or limitations; if it contains a recommendation or verdict; if
the provider overwrites history; if a later Briefing silently drops an answered Probe;
or if the flow requires Spacedock.

## First-version cuts

The first version does not require automatic spec incorporation, exposed Probe
management, scope controls, applicability language, multi-question lenses,
multi-Probe synthesis, contextual preset matching, sharing UI, or a portable Probe
format. Shared definitions may later live in project or team configuration. Personal
definitions may later come from a skill or profile. The first slice preserves and
rechecks one previously resolved concern.
