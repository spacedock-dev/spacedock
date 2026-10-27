# The stage lifecycle

Each work item moves through the stages in your workflow. Each stage defines the work to do and the evidence required.

## A development workflow

```mermaid
flowchart LR
  design --> implementation --> review --> done
  review -. changes .-> implementation
  classDef gate stroke:#e3b04b,stroke-width:2.5px
  class design,review gate
```

In this example, the outlined stages pause for your decision: before implementation and after review.

Requested changes return to implementation. The reviewer checks the revised work before it returns to you.

## Stage settings

You can ask your agent to configure a stage to:

- Pause at a [gate](gates-and-decisions.md).
- Use a separate reviewer with fresh context.
- Return findings to the stage responsible for corrections.
- Limit concurrent work or use a specialist agent.

## Isolated work

Code stages can use separate Git worktrees so work items have their own branches and files.

<a id="a-typical-dev-workflow"></a>
<a id="fresh-context-at-validation"></a>
<a id="isolated-worktrees"></a>
<a id="what-a-stage-declares"></a>
<a id="where-to-go-next"></a>

**Related guides**

- [Workflow design](../running-workflows/commission.md)
- [Gates and decisions](gates-and-decisions.md)
