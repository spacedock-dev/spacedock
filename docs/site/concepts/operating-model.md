# The operating model

Spacedock runs work through three roles. You set the mission and make the calls; agents do the work and bring decisions back to you.

## Three roles

| Role | Who |
|------|-----|
| **Captain** | You. You define the mission and make the calls at approval gates unless delegated. |
| **First Officer** | The orchestrator agent that runs the workflow and reports to you at gates. |
| **Ensign** | The worker agent that moves one item forward through one stage. |

The first officer reads the workflow README, checks which items are ready to advance, and dispatches ensigns. Stages that need isolation run in their own git worktree; lightweight stages run inline. At a gate, the first officer pauses and presents the stage report for a decision: approve, redo with feedback, or reject. Some gates wait on you; others resolve through a delegated agent review.

## The shaping/driving split

At larger scale the orchestration splits into two operating roles:

- **The shaping first officer** owns strategy and shape — the roadmap, sprint definition (deliverable plus definition-of-done), the gating ideation of each sprint's entities, the staff readiness review, and packaging each sprint as a dispatch. It stays high-level and does not hand-drive stage execution.
- **The Commander** takes one packaged sprint and drives it to its deliverable: dispatches each stage, approves execution gates and merges with judgment, runs the sprint-wide integration test, and produces the report. It escalates only on a third feedback cycle, a budget blowout, an irrecoverable block, or a genuine scope fork.

See [Sprints & roadmap](../advanced/sprints-and-roadmap.md) for the strategy layer this split serves.

## How the model serves you

The point of the model is that decisions are batched and evidenced, not scattered. You queue work; agents advance each item through its stages; you handle gates as they surface, not one session at a time. Each decision leaves a trail — a stage report with findings, verdicts, artifacts, and anomalies — so you decide on evidence, not the transcript.
