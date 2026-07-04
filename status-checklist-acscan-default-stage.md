---
id: tvstbznw83y5vwgc0jhr6ss8
title: "status --read --checklist/--ac-scan default --stage to the entity's current stage (drop the required-flag round-trip)"
status: backlog
sprint: 0203-fo-efficiency
source: "FO session 2026-07-04, boot-friction #3: status --read <entity> --checklist and --ac-scan both error 'requires --stage <stage>' when --stage is omitted (reproduced against entity 72, current stage validation), forcing a round-trip to learn a stage the entity's own status field already names. Distinct code path from 3t (--where robustness) and fk (--read frontmatter projection). Drafted by the session's science-officer teammate; sibling in the 0203-fo-efficiency sprint."
started:
completed:
verdict:
score: 0.3
worktree:
issue:
---

status --read <entity> --checklist and --read <entity> --ac-scan hard-require an explicit --stage even though the entity's own status frontmatter field already names its current stage, forcing a two-call sequence (read status, then re-issue with --stage) at gate assembly — one of the hottest FO read paths. Proposed direction: default --stage to the entity's current status when omitted, for both flags, while keeping explicit --stage <name> working for reading a non-current stage's report/ACs. Read-only projection change in internal/status; no mutation, no guard-path change.
