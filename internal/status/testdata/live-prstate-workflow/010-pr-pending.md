---
id: "010"
title: PR-pending non-terminal entity
status: implementation
pr: "#42"
score: "0.50"
source: roadmap
---
# PR-pending non-terminal entity

A non-terminal entity carrying a `pr:` field. At boot, `checkPRStates` runs
`gh pr view 42 --json state` LIVE and serializes the result into
`pr_state.entries[].state`. The stored `pr:` is `#42`; the live state is whatever
`gh` reports (MERGED in the pin test), proving the boot envelope carries live
merge state, not the stored field.
