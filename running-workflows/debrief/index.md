---
title: "Debrief a session"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-25 04:51:37"
---

# Debrief a session

A workflow improves when friction is recorded and revisited instead of lost when the session ends. `/spacedock:debrief` captures the session into a chronological record of the work completed: what shipped, what you decided at the gates, and where it rubbed. Run it before you stop; the next session's first officer reads the latest debrief and starts with context instead of cold.

The friction splits two ways:

- **Workflow friction stays local.** Quirks in your pipeline, fixed by tightening the workflow's README.
- **Spacedock friction goes upstream.** For each framework issue, debrief offers to file an anonymized GitHub issue: the body carries the bug, repro steps, and scale, never your mission, entity titles, or domain. You approve, edit, or decline each one before anything is filed.

Debrief drafts everything from git and the workflow files, and pauses for you three times: confirm where the session starts, add the why behind decisions and observations, approve each upstream issue. It commits the record, and the next debrief picks up where this one ended.

## Sitemap

- [Commission a workflow](../commission/index.md)
- [Operate a workflow](../operating/index.md)
