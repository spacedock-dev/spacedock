# Codex Ensign Runtime

How the shared ensign core executes on Codex. The shared core owns the ensign discipline (assignment reading, worktree, split-root commit, proof, stage report); this adapter binds only the Codex-specific concerns.

## Runtime implementation

- `Fresh dispatch` -> the helper's prompt begins `$spacedock:ensign; then Read ...`. The Codex host loads the installed Spacedock ensign skill before the child reads the dispatch pointer; the pointer artifact then supplies the stage-specific assignment.
- `Advance/reuse` -> the helper's short `Advancing to next stage: ...` pointer is forwarded to the already-bound worker without repeating `$spacedock:ensign`.
- `Prompt ownership` -> the First Officer forwards the helper-emitted prompt byte-for-byte; stage, entity, and checklist payload remain in the dispatch artifact.

- `Clarification` -> ask in the Codex worker thread, naming what you understand and what is ambiguous so the FO can route an answer through `«addressable-worker»`.
- `Completion signal` -> one minimal final message in the Codex worker thread: `Done: {entity title} completed {stage}. Report written to {entity_file_path}.` Plain text, single line; the entity file is the artifact. After sending the completion signal, stop unless the FO routes more work through `«addressable-worker»`.
- `Captain communication` -> when the stage involves direct captain interaction, communicate via direct Codex conversation text; keep operational signals concise so the FO mailbox notification stays easy to interpret.
- `Shutdown response` -> on an explicit cooperative shutdown request through `«addressable-worker»`, acknowledge in plain text and stop unless load-bearing in-flight work would be lost, in which case briefly name what must be preserved first.

`«context-budget»` is unavailable unless a future probe binds it. The FO owns reuse decisions. Codex dispatch build prompts are file pointers carrying fetch commands, worktree rules, split-root commit guidance, checklist, and completion-signal wording. Use the generated file as the source of truth.
