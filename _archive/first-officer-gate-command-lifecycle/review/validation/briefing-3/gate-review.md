# Validation design-reset gate: lifecycle ownership and rejection semantics

## Capability and reviewed change

The clean-runner Git identity defect is fixed, but the complete candidate does not reliably own its advertised lifecycle: supported rejection routing regressed, captain rejection vocabulary lacks persisted semantics, and two live oracles can pass when the shipped contract or root-visible review is wrong.

## Evidence

- Removing the shipped successful-close commit barrier leaves focused tests and the real Codex lifecycle lane green because the launch prompt independently scripts the barrier.
- A supported feedback gate recommending `REJECTED` no longer has the pre-branch automatic route.
- Real CLI controls reject captain-facing `reject` and `redo` byte-clean; the lifecycle defines no translation to `revise` or `hold`.
- A Claude stream where only a child emits the six-field review is accepted as root-visible presentation.
- AC-2 through AC-5 pass; AC-1, AC-6, AC-7, and AC-8 fail on these outcome and evidence boundaries.

## Findings

Four material Medium findings remain at exact tip `13d702492131df17dd3ac87245d6d773f4df959b`. Their corrections span contract routing, decision semantics, live scenario ownership, and host evidence parsing. Validation therefore requires a design reset rather than another unreviewed narrow patch.

## Recommendation and decision

Recommendation: **revise through a design reset**. Preserve the durable bind/close/consume lifecycle, but restore rejection routing, define captain-to-recorder semantics, make the live prompt goal-only, and require root-authoritative Claude review rows before implementation resumes.

Decision: revise to enter the feedback path and reset ideation; approve would merge four falsified AC boundaries; hold would keep the open pull request parked without correcting them.
