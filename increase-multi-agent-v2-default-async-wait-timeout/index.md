---
title: Increase multi_agent_v2 default async wait timeout
status: ideation
score: 0.5
source: "Captain filing request 2026-07-11."
completed:
verdict:
worktree:
issue:
id: 95we0fhydgx5rbay5fw3qy4q
started: 2026-07-11T11:51:34Z
---

`multi_agent_v2`’s short default wait causes repeated timeout churn while subagents perform normal verification or interactive work. Because the wait already returns immediately when an agent responds or new user input arrives, increase the default timeout to five minutes. This reduces unnecessary polling and tool traffic without delaying updates or reducing responsiveness.
