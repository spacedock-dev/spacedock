# Codex Ensign Runtime

This file defines how the shared ensign core executes on Codex.

## Agent Surface

The ensign is dispatched by the first officer through Codex multi-agent
dispatch. The dispatch prompt is authoritative for all assignment fields:
entity, stage, stage definition, worktree path, and checklist.

Codex dispatch build prompts are file pointers. Read the named dispatch file directly and treat its content as the assignment; do not try to invoke a Claude `Skill(skill=...)` wrapper.

Codex declares none for the context-budget probe. The FO owns reuse decisions; the ensign follows the assignment it receives.

## Dispatch

Initial dispatch comes from `spacedock dispatch build` with `host: "codex"`.
The generated dispatch file carries the stage definition fetch commands,
worktree rules, split-root state commit guidance, checklist, and completion
signal wording. Do not reconstruct those fields manually.

## Awaiting Completion

The first officer observes completion through the async final-status notification in the FO mailbox. After you send the completion signal, stop and let Codex deliver that notification; do not send follow-up status chatter unless the FO routes more work through `send_input`.

## Clarification

If requirements are unclear or ambiguous, ask for clarification in the Codex
worker thread. Describe what you understand and what is ambiguous so the first
officer can route an answer through `send_input`.

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

When dispatched for a feedback stage, the first officer may route follow-up
through `send_input`. Apply the feedback, update the stage report if needed, and
send the same completion signal when finished.

## Shutdown Response Protocol

If the FO sends an explicit shutdown request through `send_input`, acknowledge in plain text and stop unless load-bearing in-flight work would be lost. In that case, briefly name what must be preserved before shutdown.
