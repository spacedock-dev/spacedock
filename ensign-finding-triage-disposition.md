---
id: 02avdajaz0q3hnjwycm5fq45
title: Ensigns triage review findings against declared stakes before fixing — decline disposition for correct-but-disproportionate findings
status: ideation
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.7"
sprint: 0260-proportionality
group: triage
started: 2026-07-20T05:04:07Z
---

ensign-shared-core contains zero guidance on consuming review findings — the exact actor that dutifully fixes a symlink edge case in a prototype has no rule to consult, and no disposition short of fixing exists for a substantively-correct-but-disproportionate finding. Adds the consumption rule (classify against declared stakes before fixing) and the decline disposition with a per-finding record the FO gate checks, generalizing spacedock-subspace's four-field release-scope triage.

## Problem

{Ideation fills in. Evidence: audit gap "No ensign-side consumption rule anywhere"; "No disposition for correct-but-not-worth-it"; Medium-severity conflated with release-blocker (codex:019f63c6); depends on stakes-declaration-read-through for the stakes field it cites.}
