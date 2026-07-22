---
title: status pagination and stage-then-score default sorting
status: ideation
score: 0.70
id: rwpe45pdxffk2zfy24ejde6a
started: 2026-07-22T06:31:16Z
---

### Goal
Update `spacedock status` so output is paginated and sorted by stage (closer to the end first: e.g. validation, implementation, ideation, backlog) then score descending by default.

### Scope
- Implement stage-then-score default sorting for status listing.
- Add pagination flags/options (e.g. `--page`, `--limit`, or interactive terminal pagination where appropriate).
- Add automated unit and CLI fixture tests.
