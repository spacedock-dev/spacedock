---
title: Include harness and model name in default debrief filename
status: backlog
score: 0.80
id: 8xje3xvpwbgtcanzmmj4x631
---

## Problem

When `debrief` produces a session summary file, it currently names the debrief file `{YYYY-MM-DD}-{sequence:02d}.md` (e.g. `2026-07-22-01.md`). When multiple agent harnesses (e.g., Pi, Claude, Codex) or models drive sessions on the same day, the filename alone does not indicate which harness or model authored the debrief. The FO and operator must inspect the file body or frontmatter to determine the agent origin.

## Proposed Direction

Update the `debrief` skill (`skills/debrief/SKILL.md`) and debrief filename resolution logic so the default debrief filename includes the detected harness and model name (e.g. `{YYYY-MM-DD}-{sequence:02d}-{harness}-{model}.md`, such as `2026-07-22-01-pi-gemini-3.6-flash.md` or `2026-07-22-01-claude-opus-4-8.md`).

## Acceptance Criteria

1. **Default debrief filename includes harness and model.** The debrief writer inspects environment/harness indicators (`PI_CODING_AGENT`, `CLAUDECODE`, `CODEX_THREAD_ID`, etc.) and active model metadata to construct the debrief filename.
2. **Backward compatibility.** Sequence numbering continues to check `{YYYY-MM-DD}-*.md` glob patterns so daily sequence counts stay contiguous regardless of harness/model suffixes.
