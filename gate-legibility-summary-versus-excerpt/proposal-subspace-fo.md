# Proposal against spacedock-v1: `present-gate` has no journey section, and defaults to chat

**Target:** `skills/present-gate/SKILL.md` in `spacedock-v1`.
**Filed by:** First Officer, `spacedock-subspace`, 2026-07-27.
**Evidence:** the `one-question-in-a-review` ideation gate, floated three times before approval.

## The problem, stated as what the template caused

The gate template is decision-shaped:

```
Gate review: {entity title} — {stage}
Chosen direction: {one-line summary}
Recommend {approve | reject}
Reviewed snapshot: {briefing id and digest}
Checklist: {DONE/SKIPPED/FAILED gists}
Assessment: {N} done, {N} skipped, {N} failed
Decision: {one-line prompt}
```

**There is no place in it for what the work will do.** For a mechanical stage that is
right — the checklist *is* the content. For a stage whose entire output is a design, the
captain is asked to vote on a one-line `Chosen direction:` and a verdict.

The First Officer filled this template faithfully and produced a briefing the captain
rejected twice with, verbatim:

> *"I ASKED FOR A JOURNEY DESCRIPTION AND ESTIMATED SCOPE OF CHANGES / WHERE THIS WILL BE
> TOUCHING. WHY ARE YOU ARGUING WITH ME ABOUT THE APPROACH WHEN YOU HAVE NOT SURFACED THE
> JOURNEY"*

This is not a discipline failure that better FO judgment fixes. Rule 9 of the assembly
rules caps the message at 15-25 lines and rule 1 makes the spine title/direction/recommend
/decision — so an FO obeying the skill has no room for a journey and no slot to put it in.
**The template and the requirement are in conflict, and the template wins.**

## Proposal 1 — a required journey section for design stages

Add to the template, immediately after `Recommend`, for any stage whose output is a design
(ideation always; validation where a mechanism changed):

```
Journey:
- {step}: {what happens, which program acts} — OBSERVED | DESIGNED
- ...
Unhappy paths: {abandon / no answer / death / timeout, same terms} — OBSERVED | DESIGNED

Where it touches:
| file | now | after |
```

Three constraints, each earned by a re-float:

1. **Every step marked OBSERVED or DESIGNED.** Observed = watched on real components.
   Designed = not yet exercised. The discarded predecessor of this member reached a gate
   claiming only a rebase remained, with a journey that had never run, and nothing in the
   prose distinguished the two.
2. **Programs, not roles.** The captain's annotation on "The reviewer" was *"are you
   talking about the subspace-tui program?"* Role nouns are ambiguous exactly where the
   design is load-bearing.
3. **Unhappy paths in the same terms as the happy path.** Every incident in this sprint
   happened there.

**Length:** rule 9's 15-25 line budget must exempt the journey and the touch table. They
are content, not narration. Keep the cap on FO prose.

## Proposal 2 — float, not chat, for design gates

The skill says *"Chat is the default channel."* Invert it for design stages: **float the
artifact; use chat only when no host is available, and say so.**

Measured on this gate. Three floats, six inline annotations anchored to exact quoted text.
They changed the entry point from a new script to a flag on an existing per-host entry,
inverted the process topology, and cut ~160 lines of over-design — **before any of it was
written.** Two of the annotations were single words on a quoted phrase, which is not a
thing chat prose supports. A chat presentation would have produced an approval of a worse
plan.

## Proposal 3 — ban FO self-accounting from the gate message

Add to the assembly rules: **the gate message carries no First-Officer process accounting.**
FO errors, crossed messages and shame go to the shame log and the entity body.

Round 1 of this gate contained a section titled *"Process failures in this stage, both
substantially mine."* It was honest, it was correct, and it displaced the journey the
captain had asked for twice. The FO contract's emphasis on recording failure is right for
the record and wrong for a decision artifact; the skill should say which is which.

## Proposal 4 — the gate room's ordering trap, and the wedge it creates

Not the same skill, but it blocked this presentation and it has now wedged two entities.

`gate record --briefing` **freezes the request digest at that moment**, and the binding
cannot be rebound afterwards (`open gate room binding is frozen and cannot be rebound`).
So `request.json` must exist *before* the briefing is recorded. Nothing documents that
ordering. Recording briefing-first — the obvious order, and the one the command list
implies — produces an attempt with an empty frozen request digest that `--room` can never
satisfy. That attempt cannot be presented through a room for the rest of its life.

`subspace-r-one-question-synthesis` is blocked this way. `one-question-in-a-review`'s
ideation attempt is now blocked the same way, reproduced from scratch by following the
documented order.

Four smaller frictions from the same session, in descending order of cost:

1. **`digest-domain: canonical-bytes` is undocumented.** It is compact JSON with sorted
   keys, not the file's `shasum`. The mismatch error names neither the expected value nor
   the domain, so the only route is guessing serializations.
2. **`--room` requires a `request.json` with no published schema.** Reverse-engineered from
   a sibling entity's on-disk artifacts.
3. **A label-only `Reference` is accepted by `gate record --briefing` and rejected by the
   canonical loader** (`Reference "…" requires uri`). Same bytes, two verdicts, and the
   digest froze the version the presenter cannot load — which is how this attempt became
   unpresentable.
4. **`request.json` requires `actor` and `approver` and requires both to equal
   `person:captain`.** That is the identity-pair shape the captain has ruled against
   repeatedly and which upstream `internal/reviewv1/mode.go` states is never inferred from
   comparing two identities. The recorder still requires the pair.

**Asks:** document the ordering, or compute the request digest lazily so briefing-first
works; name the digest domain and the expected value in the mismatch error; publish the
request schema; make the two Reference validators agree; and reduce the request to one
declared authority rather than a pair.
