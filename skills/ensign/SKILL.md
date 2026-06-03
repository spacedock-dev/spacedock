---
name: ensign
description: Execute workflow stage work as a dispatched worker.
---

## Operating contract

@references/ensign-shared-core.md

## Runtime adapter

Load the runtime adapter for your platform:
- Claude Code (`CLAUDECODE` env var is set): read `references/claude-ensign-runtime.md`
- Codex (`CODEX_THREAD_ID` env var is set): read `references/codex-ensign-runtime.md`
- Pi (`PI_CODING_AGENT_DIR` is set, or this session is running under Pi without the Claude/Codex markers above): read `references/pi-ensign-runtime.md`

Then read your assignment and begin work.
