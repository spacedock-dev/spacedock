# The operating model

Spacedock runs on three roles and one division of labor: you shape the work and make the calls; the agents drive each item through its stages and bring decisions back to you with evidence.

## Three roles

| Role | Who | What they own |
|------|-----|---------------|
| **Captain** | You. | The mission, and the call at every approval gate unless you delegate it. |
| **First Officer** | The orchestrator agent. | Running the workflow: dispatch, gate presentation, advancing entity state. |
| **Ensign** | The worker agent. | Moving one entity (one work item) forward through one stage. |

There is one captain and one first officer per session. The number of ensigns tracks the dispatchable work: the first officer dispatches one per entity per stage.

The first officer reads the workflow README, runs `spacedock status --next` to find entities ready to advance, and dispatches an ensign for each. An ensign reads its assignment, does the stage's work, commits, writes a stage report, and signals done. The first officer reviews that report against the checklist it dispatched. If the stage is gated, it pauses and presents the report to you. If not, it advances the entity and dispatches the next stage itself. A completed non-gated, non-terminal stage is not a stopping point.

## Shaping versus driving

The captain shapes; the agents drive. These are different jobs, and the split is what keeps you out of the per-step loop.

**Shaping is defining what good looks like before the work runs.** You set the mission, the stages, and the bar each stage must clear, all declared in the workflow README rather than negotiated mid-task. You commission a workflow with [`/spacedock:commission`](../running-workflows/commission.md), and you make the calls at gates: approve, redo with feedback, or reject. Some gates you answer yourself; others resolve through a delegated agent review. That is the whole of the captain's standing job.

**Driving is moving entities through the declared stages.** The first officer schedules and dispatches; the ensign does the stage work and proves it. The first officer is allowed to take obvious reversible steps without asking: a dispatch the workflow already permits, a status read, a routine state transition. It asks you only when requirements are materially ambiguous, a design choice would change the output meaningfully, or scope is too unclear to turn into concrete criteria. Everything else happens without a prompt to you.

The line holds because of one rule: the maker does not judge its own work. Review runs as a separate stage with fresh context and no access to the ensign's reasoning (see [Gates and decisions](gates-and-decisions.md)). The first officer never self-approves a gated stage.

## Batched, evidenced decisions

Decisions reach you batched and backed by evidence, not as a stream of interruptions. This is the point of the model: your attention is the bottleneck, so the agents queue work and surface only the calls that need a human.

**Batch the work; decide as it flows back.** Queue many entities at once. Ensigns advance each through its stages in parallel. You handle gates as they surface, not one session at a time, and not on the agent's schedule. While one entity waits on a clarification, the first officer keeps dispatching the others.

**Every gate carries evidence.** When the first officer presents a gate, it does not hand you the transcript. It renders a fixed gate-review format: the chosen direction in one line, a single clear recommendation you can approve with one "yes", a gist roll-up of the stage report's `DONE`/`SKIPPED`/`FAILED` items cited by file path and line range, and any reviewer findings split into `Material` and `Polish` tiers. The target is 15-25 lines. You decide on the evidence and the bar, not on a wall of output.

**The decision leaves a trail.** Each gate records the verdict and its reason alongside the stage report in the entity file. The record outlives the reviewer, so a bad result traces back to the call that caused it. When you end a session, [`/spacedock:debrief`](../running-workflows/debrief-and-refit.md) captures what happened (commits, state changes, decisions, open issues), and the next session picks up from it.

## Where to go next

- [Workflows and entities](workflows-and-entities.md) covers the directory and files the roles operate on.
- [Stage lifecycle](stage-lifecycle.md) walks an entity backlog → ideation → implementation → validation → done.
- [Gates and decisions](gates-and-decisions.md) lays out what a gate review carries, the three calls you make, and the detached adversarial audit.
