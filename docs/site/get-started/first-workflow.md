# Your first workflow

[Install Spacedock](install.md), then describe a repeatable piece of work and decide where it needs your judgment.

## Describe the work

```bash
spacedock claude "/spacedock:commission ship features through design, implementation, and review"
```

The agent asks about your work and proposes a workflow. For example:

> Each task moves through design, implementation and review.
>
> You approve the design before implementation. A separate agent reviews the result before your final decision.
>
> Requested changes return to implementation.

## Set the standards

Adjust the proposed stages, decision points and quality standards. Type `review stages` to walk through the standards with your agent.

For each decision point, specify the context and evidence you need. This is the workflow's [gate content](../concepts/gates-and-decisions.md#gate-content).

## Run the first tasks

Choose the first work items with your agent. After you accept the design and standards, the agent starts the work.

At a gate, the agent presents the evidence and asks for your decision. Approve, request changes or hold the work.

## Resume later

Run [`/spacedock:debrief`](../running-workflows/debrief.md) before you end the session.

To resume:

```bash
spacedock claude
```

The agent reads the workflow records and picks up the work.

<a id="commission-a-workflow"></a>
<a id="operate-the-workflow"></a>
<a id="the-pilot-run"></a>
<a id="what-gets-generated"></a>

**Related guides**

- [Workflow operation](../running-workflows/operating.md)
- [Project survey](survey.md): start from your existing agent history.
