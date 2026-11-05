# Spacedock Workflow

Your agents can do more work than you can review. Spacedock Workflow lets you define what needs your attention and what agents can handle themselves.

## How it works

A workflow gives recurring work a set of stages, such as design, implementation and review. You define the standards for each stage. Agents choose how to do the work within those standards.

You also choose where work pauses for your decision. These checkpoints are called [gates](../concepts/gates-and-decisions.md).

The stages and standards live in plain text in your repository. Your agent coordinates the work and records its progress and decisions.

## Gate content

Different stages need different context and evidence for your judgment. **Gate content** defines what the agent presents at each decision point.

For example, a feature workflow can pause before implementation. The agent brings you the proposed behavior, tradeoffs and a recommendation. You approve the plan, request changes or hold the work.

After approval, agents can implement and review the change against your standards. You can require another decision before release, with test results and unresolved risks as evidence.

[Spacedock Review](../review/index.md) presents the work for your feedback and decision. Workflow defines when that review happens, what it contains and what happens next.

## Standards and trust

You and your agent can refine the standards and gate content as you learn what each decision needs. Add missing context, remove distracting detail, or change the evidence required.

You can also add decision points or explicitly delegate decisions to an agent. You control how much attention the workflow requires.

<a id="next"></a>

**Related guides**

- [Your first workflow](../get-started/first-workflow.md): design a workflow and run its first work items.
- [Project survey](../get-started/survey.md): identify the process in your existing agent history.
- [Review standards](standards.md): define what good work looks like at each stage.
