# Workflow operation

Talk to your agent to add work, check progress and make decisions.

## Start a session

```bash
spacedock claude
```

The agent reads the workflow records and continues work that is ready. It coordinates workers and brings you decisions or blockers.

## Check progress

Ask questions such as:

- What needs my decision?
- What is blocked?
- What is ready next?

You can add work or change priorities in the same conversation.

## Decide at a gate

The agent presents the evidence specified by the workflow's [gate content](../concepts/gates-and-decisions.md#gate-content). You can approve, request changes or hold the work.

Your decision and feedback become part of the work item's record. The agent continues or routes the requested changes according to the workflow.

## End a session

Run [`/spacedock:debrief`](debrief.md) to record the completed work, decisions and follow-up items.

The next session reads the workflow records to resume the work.

<a id="delegating-the-loop"></a>
<a id="ending-a-session"></a>
<a id="the-session-loop"></a>
<a id="when-a-gate-arrives"></a>
