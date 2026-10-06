---
title: The commission skill lets an initial seed gate claim to review work no worker produced
status: backlog
source: "Captain, 2026-10-06: incident while commissioning an exploratory workflow on 0.28.0-pre4. An `exploration` stage was declared both initial and gated, and its prose implied the gate reviewed research no worker had run. Approving it would have advanced straight past the research worker."
id: 6dhxxd047s3c32n21h8ppf9y
gates:
    version: 1
    records:
        - id: gate:6dhxxd047s3c32n21h8ppf9y:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:6dhxxd047s3c32n21h8ppf9y-backlog-1
              briefing:
                id: briefing:6dhxxd047s3c32n21h8ppf9y:backlog:attempt-1:revision-1
                digest: sha256:a3c37779e5cbc8beb54e51b618936ef0718fe5fe474b989d565b69b2bbaffb1c
                room-ref: '@review/backlog/briefing-1'
---

A commissioned workflow can declare an initial stage gated, then describe that gate as reviewing the stage's own worker output. The binary's semantics are the opposite: an initial gated stage reviews the committed seed, because no prior stage wrote a report. Approving it advances to the successor without the stage's worker ever running.

Verified in `internal/status/entered_stage.go`: "A non-initial gated stage owes a complete committed stage report. An INITIAL gated stage had no prior stage to write one: the committed clean seed IS the artifact the captain reviews, so durability alone is the proof."

Note that `initial: true` with `gate: true` is legitimate — this repo's own `docs/dev` workflow uses `backlog` that way. The defect is not the combination; it is a gate whose prose promises review of output the stage has not produced.

## Proposed change

Documentation only. No new standing lint, lane, or validator.

1. `skills/commission/SKILL.md` — explain gate timing where stages are designed: what an initial gate reviews, that the reviewed artifact is the seed, and that permission to begin work and permission to build from its results are separate decisions.
2. `skills/commission/SKILL.md` generation checklist — for every gated stage, name the reviewed artifact and the actual successor work, and confirm no stage is skipped when a seed approval advances.
3. `skills/present-gate/SKILL.md` — when the gate is an initial one, state plainly that it reviews the committed seed and that no worker has run. This surfaces the mismatch at the decision point instead of at authoring time.
4. Any worked example must pass `status --validate` before it lands.

## Acceptance

Behavioural, in a disposable commissioned workflow. Presence-grepping the skill text is not evidence.

1. Commission a seed gate followed by a research stage and a prototype stage.
2. Confirm the initial ready gate is attached to the seed stage and reviews only the seed.
3. Approve it and observe that the spawned worker assignment targets the research stage, not the prototype stage.
4. Confirm the research-result gate cannot be presented before the worker's report exists.
5. After completed research, confirm approval advances to the prototype stage.

Falsifying outcome: an initial seed approval that skips research, or a research-result review presented without a completed report.
