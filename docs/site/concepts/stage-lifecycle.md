# The stage lifecycle

Work moves through stages. Each stage declares what good looks like, produces a deliverable, and routes the entity to the next stage. The development workflow's lifecycle is the canonical example: backlog → ideation → implementation → validation → done.

## The stages

- **backlog** — a work item enters here when it's first proposed. It has a seed description but no design work yet.
- **ideation** — a pilot fleshes out the idea: clarify the problem, explore approaches, and produce a concrete definition of done with acceptance criteria and a test plan.
- **implementation** — produce the deliverable. Write the code, generate the fixtures, make the changes the work item describes. Implementation is complete when the deliverable exists and is ready for independent verification.
- **validation** — a *fresh* agent verifies the deliverable against the acceptance criteria. The validator checks what was produced; it does not produce the deliverable itself.
- **done** — validation is complete and the captain approves. The entity closes with a verdict of PASSED or REJECTED.

## Fresh context at validation

Validation runs with fresh context — a separate agent that has no access to the maker's reasoning. The agent doesn't get to judge its own work. This is what makes review push back on thin evidence instead of rubber-stamping the implementation it just wrote.

## Worktree vs. inline

Stages that touch shared state run in their own git worktree, so concurrent work doesn't collide. Lighter stages run inline, in the repo root. The workflow README declares which stages need a worktree.

## What a stage declares

Each stage in the README names its inputs, its outputs, and — most importantly — what *good* and *bad* look like for that stage. The agent works to that line on its own. When a standard turns out fuzzy in practice, the agent proposes an edit to the stage's written criteria for your approval, so the bar sharpens as you use it.

For how the transition between stages is decided, see [Gates & decisions](gates-and-decisions.md).
