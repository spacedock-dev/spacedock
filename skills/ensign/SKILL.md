---
name: ensign
description: Execute workflow stage work as a dispatched worker.
---

## Operating contract

@references/ensign-shared-core.md

## Runtime adapter

Load ONLY the first matching runtime adapter (do not load a second inner-host adapter):
- Paseo or host-neutral dispatch (the dispatch prompt names `references/paseo-ensign-runtime.md`, OR `PASEO_AGENT_ID` is set with none of `CLAUDECODE`, `CODEX_THREAD_ID`, `PI_CODING_AGENT_DIR`): read `references/paseo-ensign-runtime.md`
- Claude Code (`CLAUDECODE` env var is set): read `references/claude-ensign-runtime.md`
- Codex (`CODEX_THREAD_ID` env var is set): read `references/codex-ensign-runtime.md`
- Pi (`PI_CODING_AGENT_DIR` is set, or this session is running under Pi without the Claude/Codex markers above): read `references/pi-ensign-runtime.md`
- Other harness (none of the above): read `references/paseo-ensign-runtime.md`

Then read your assignment and begin work.
