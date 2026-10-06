# Seed review: the commission skill lets an initial gate claim to review work no worker produced

## What this seed is

A documentation change to the Spacedock plugin, prompted by an incident while commissioning an exploratory workflow on `0.28.0-pre4`.

## The incident

An `exploration` stage was declared both `initial: true` and `gate: true`. Its prose described the gate as approving the researched direction. But an initial gated stage reviews the committed seed, because no prior stage wrote a report — so approving it would have advanced straight to `prototype`, skipping the research worker and the requested reference scout.

The captain caught it before any worker launched. The workflow was corrected to `backlog (initial, gate) -> exploration (gate) -> prototype`.

## Why it is worth a fix

`internal/status/entered_stage.go` states the semantics outright: "A non-initial gated stage owes a complete committed stage report. An INITIAL gated stage had no prior stage to write one: the committed clean seed IS the artifact the captain reviews, so durability alone is the proof."

`initial: true` with `gate: true` is legitimate — this repo's own `docs/dev` workflow uses `backlog` that way. The defect is a gate whose prose promises review of output its stage has not produced, and the skill neither teaches the distinction nor checks it at generation time.

## Proposed change (documentation only)

1. `skills/commission/SKILL.md` — teach gate timing where stages are designed: what an initial gate reviews, and that permission to begin work and permission to build from results are separate decisions.
2. `skills/commission/SKILL.md` generation checklist — for every gated stage, name the reviewed artifact and the actual successor work.
3. `skills/present-gate/SKILL.md` — when the gate is initial, state plainly that it reviews the committed seed and that no worker has run. This catches the mismatch at the decision point rather than at authoring time.
4. Any worked example must pass `status --validate` before it lands.

No new standing lint, lane, or validator is proposed.

## Acceptance (behavioural, not text-matching)

1. Commission a seed gate, then a research stage, then a prototype stage.
2. Confirm the initial ready gate attaches to the seed stage and reviews only the seed.
3. Approve it, and observe the spawned worker assignment targets the research stage, not the prototype stage.
4. Confirm the research-result gate cannot be presented before its report exists.
5. After completed research, confirm approval advances to the prototype stage.

Falsifying outcome: an initial seed approval that skips research, or a research-result review presented without a completed report.
