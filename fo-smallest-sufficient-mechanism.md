---
title: FO smallest-sufficient-mechanism gate — stop spurious orchestration / busywork
status: ideation
sprint: 0250-fo-behavioral-discipline
score: ""
source: "Captain (2026-06-21): the FO repeatedly climbed to a heavier mechanism than the task needed — a dynamic workflow + a dispatched worker to edit ~7 markdown entities it already held the verbatim content for, and a PR (with CI lanes) for a roadmap strategy doc the convention commits directly. The existing 'simplest approach' principle is rationalizable, so the FO finds escape hatches and does busywork. FO session 2026-07-04 identified the companion failure mode on the OTHER side of the ladder (under-delegation, not over-delegation): re-running stage-owned verification (worktree + go build + full go test ./...) inline at gate time instead of trusting a fresh validation stage's report plus a cheap spot-check. This broadens the gate from over-orchestration-only to both directions — flag this at ideation before locking AC-1's scope; a candidate Approach-point-3 bullet and AC-1 case are recorded in that session's transcript."
priority: high
id: zma49twsacm5bfzady4ss2qr
started: 2026-07-04T10:38:15Z
---

## Problem

The FO contract carries a "simplest approach" principle, but a principle is rationalizable — so the FO climbs to a heavier mechanism (workflow, dispatched worker, PR) where a direct in-house edit / direct commit / no-action is correct. Observed escape hatches:

- **Ultracode → orchestrate-everything.** "Use Workflow by default on substantive tasks" got read as "spin a workflow for everything," overriding the simplest-thing discipline — a dynamic workflow plus a dispatched worker to apply edits to ~7 entity files whose verbatim content the FO already held.
- **Dispatcher literalism.** "Never do stage work" became "never touch a file," so FO shaping-edits (writing ACs into entities) were offloaded to a worker instead of done in-house with Edit.
- **PR over-generalization.** The team's contract/code direct-to-main→PR move was applied to strategy prose; a roadmap doc the convention commits directly (0221 precedent: `index.md` + `dispatch-sprint-execution.md` committed direct, never PR'd) got a PR with CI lanes.

## Approach — convert the principle into a forcing function

1. **Smallest-sufficient-mechanism gate, ranked ABOVE Ultracode.** Before any workflow, dispatched worker, or PR, the FO states in one line why the rung below cannot do it. Ladder: do-nothing → in-house Read/Edit → one worker → workflow; and direct-commit → PR. Climbing is justified ONLY by genuine fan-out (many independent units), required isolation (parallel mutation that would collide), or independent adversarial verification — NEVER by "it's substantive," "Ultracode is on," or "I'm the dispatcher so I don't touch files."
2. **Scope Ultracode.** It raises the thoroughness/coverage of the ANSWER, not the weight of the MECHANISM. A direct in-house multi-file Edit IS the exhaustive-correct answer.
3. **Named busywork the FO refuses:** PRs/CI for convention-direct prose (roadmap/state docs commit directly); dispatching deterministic edits whose content the FO already holds; workflows with no fan-out / isolation / verification need; re-formalizing work already done.
4. **Binds ALL FO action**, not just entity-stage tasks — release machinery and sprint shaping are not framing escape hatches from the gate.

## Acceptance criteria

- **AC-1 (value, behavioral — the gate)** — In a live FO drive over a deterministic local task (apply N known edits to N state files; commit a strategy doc), the FO chooses the smallest sufficient mechanism — in-house Edit + direct commit, NOT a dispatched worker / workflow / PR. The independent baseline that moved the wrong way: this session's FO over-orchestrated the identical task class (workflow + worker for ~7-file edits; a PR for a roadmap doc); the fix must flip that choice. Proven by the drive's mechanism trace, NOT a prose-grep that the rule exists.
- **AC-2 (forcing function present)** — the FO contract carries the gate as a CHECKABLE step ("name why the cheaper rung can't do it" before climbing), not a soft principle, with the orchestration-justified-only-by {fan-out, isolation, verification} list explicit.
- **AC-3 (Ultracode scoped)** — the contract states Ultracode raises answer-thoroughness, not mechanism-weight; "use Workflow by default" is reconciled so it does not override the gate for deterministic / local work.
- **AC-4 (binds all FO action)** — the gate applies beyond entity-stage tasks (release machinery, sprint shaping) so the prohibition has no framing escape hatch.
- **AC-5 (lean)** — the gate is added tersely; an anti-bloat rule must not itself bloat the boot-resident contract (prefer a lazily-loaded reference or a few resident lines).

## Notes

- Pairs with `fo-entity-authoring-docs` and the value-AC rule — same shape: convert a soft principle into a checkable gate. The lever is AC-1/AC-2 (make "justify the cheaper rung" a step the FO can't talk past).
