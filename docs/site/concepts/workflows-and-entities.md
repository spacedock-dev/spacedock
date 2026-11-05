# Workflows and entities

A workflow keeps its rules and work records as plain text in your repository. An **entity** is one work item.

## Workflow rules

The workflow README defines the stages, standards and [gate content](gates-and-decisions.md#gate-content).

Ask your agent to change these rules as you learn what the work requires. You can also edit the README directly.

## Work items

Each work item has a Markdown record with its goal, acceptance criteria, progress and stage reports. Related artifacts can live beside it.

Your agent maintains the records as work proceeds. You can read them to understand what happened or resume in another session.

## State storage

Workflow state can live in a [separate checkout](../advanced/split-root-state.md) to keep routine progress commits off your code branch.

<a id="each-work-item-is-one-file"></a>
<a id="keep-workflow-state-off-your-code-branch"></a>
<a id="the-readme-is-where-you-set-the-rules"></a>

**Related guides**

- [Workflow design](../running-workflows/commission.md)
- [Frontmatter contract](../reference/frontmatter-contract.md)
