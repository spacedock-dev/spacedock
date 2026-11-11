---
title: "Gates and decisions"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 06:33:53"
---

# Gates and decisions

A gate pauses work for a decision. You choose which stages need gates and what evidence each gate presents.

## Gate content

Gate content defines what you need to judge the work at that stage.

| Decision                 | Useful context and evidence                    |
| ------------------------ | ---------------------------------------------- |
| Approve a plan           | Proposed behavior, alternatives and tradeoffs  |
| Accept an implementation | Changes, test results and unresolved findings  |
| Release a change         | User impact, remaining risks and rollback plan |

Ask your agent to adjust the gate content as you learn what helps you decide. The instructions live in the workflow README.

## Your decision

The agent presents the work, supporting evidence and a recommendation. You can respond in the conversation or use [Spacedock Review](../../review/).

| Decision            | What happens next                                                       |
| ------------------- | ----------------------------------------------------------------------- |
| **Approve**         | The agent continues with the approved work                              |
| **Request changes** | The agent routes your feedback to the stage responsible for the changes |
| **Hold**            | The agent pauses the work                                               |

The workflow records your decision and its reason.

## Rejections

A review stage can send failed work back for correction before it reaches you. The reviewer then checks the revised work.

The automatic correction loop stops after three cycles and asks you how to proceed.

## Delegation

You can explicitly authorize an agent to make specified decisions. State the scope and conditions of that authority in the conversation.

The agent still presents the evidence and records the decision under its own identity.

**Related guides**

- [Workflow operation](../../running-workflows/operating/)
- [Review standards](../../workflow/standards/)

## Sitemap

- [Overview](../../workflow/index.md)
- [Workflow design](../../running-workflows/commission/index.md)
- [Workflow operation](../../running-workflows/operating/index.md)
- [Review standards](../../workflow/standards/index.md)
