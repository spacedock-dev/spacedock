# Codex 0.144.4 compaction-hook probe

Date: 2026-07-17 (Asia/Taipei)

Purpose: find the smallest current Codex surface that can issue a post-compaction Spacedock reminder, and determine whether the reminder reaches the model.

## Setup

The live TUI was launched with `--dangerously-bypass-hook-trust` for this isolated, locally authored probe. Inline hook configuration installed:

- `PreCompact`, matcher `manual|auto`, command touching a unique temporary marker;
- `PostCompact`, matcher `manual|auto`, first touching a marker and then, in the second run, returning:

```json
{"systemMessage":"SPACEDOCK_POST_COMPACT_REMINDER_reread_the_authoritative_first_officer_contract_and_reconcile_durable_state_before_continuing"}
```

- `SessionStart`, matcher `compact`, one command touching a marker and one printing an additional-context sentinel.

The session received `Reply READY only.`, then the TUI command `/compact`.

## Observations

First run:

- Codex displayed `Context compacted`.
- The unique `PreCompact` and `PostCompact` markers existed afterward.
- The `SessionStart(compact)` marker did not exist.

Second run:

- Codex displayed `Context compacted`.
- It then displayed `PostCompact hook (completed)` and a warning containing the full sentinel.
- The next prompt said: `Without using tools, repeat the PostCompact hook warning you received. If none was in your model context, reply NONE.`
- The model replied `NONE`.

The sessions exited normally. The temporary marker directory was `/tmp/spacedock-c6-codex-hook-probe.1hqWsz`; the two resumable test session IDs printed by the client were `019f70a0-db79-77c1-956b-217873dd484d` and `019f70a3-6e00-7b93-b109-ed5d91573010`.

## Boundary

On this client, `PostCompact` is a real captain/UI warning surface but not a demonstrated developer-context surface. Same-session manual `/compact` did not deliver `SessionStart(source=compact)`. Spacedock must not claim automatic post-compaction FO instruction from either event on this evidence. A visible warning asking the captain to cue the FO is supported; absent or failed hooks must remain harmless.
