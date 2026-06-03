# Codex Ensign Runtime

How the shared ensign core executes on Codex.

## Agent Surface

The ensign is dispatched through Codex multi-agent dispatch. The dispatch prompt is authoritative for all assignment fields: entity, stage, stage definition, worktree path, and checklist.

Codex dispatch build prompts are file pointers. Read the named dispatch file directly and treat its content as the assignment; do not invoke a Claude `Skill(skill=...)` wrapper.

Codex declares none for the context-budget probe. The FO owns reuse decisions; the ensign follows the assignment.

## Dispatch

Initial dispatch comes from `spacedock dispatch build` with `host: "codex"`. The generated file carries fetch commands, worktree rules, split-root state commit guidance, checklist, and completion-signal wording. Do not reconstruct those fields manually.

## Awaiting Completion

The FO observes completion through the async final-status notification in the FO mailbox. After sending the completion signal, stop and let Codex deliver it; do not send follow-up chatter unless the FO routes more work through `send_input`.

## Clarification

If requirements are unclear, ask in the Codex worker thread. Describe what you understand and what is ambiguous so the FO can route an answer through `send_input`.

## Captain Communication

When dispatched for a stage that involves direct interaction with the captain, communicate via direct Codex conversation text. Keep operational completion signals concise so the FO mailbox notification stays easy to interpret.

## Completion Signal

When your work is done, send one minimal final message in the Codex worker
thread:

```text
Done: {entity title} completed {stage}. Report written to {entity_file_path}.
```

The entity file is the artifact. Do not include the checklist or summary in the message. Plain text only. Never send JSON; never emit a Claude `SendMessage(to="team-lead", ...)` call on Codex.

## Feedback Interaction

For feedback stages, the FO may route follow-up through `send_input`. Apply the feedback, update the stage report if needed, and send the same completion signal when finished.

## Shutdown Response Protocol

If the FO sends an explicit shutdown request through `send_input`, acknowledge in plain text and stop unless load-bearing in-flight work would be lost. In that case, briefly name what must be preserved before shutdown.
