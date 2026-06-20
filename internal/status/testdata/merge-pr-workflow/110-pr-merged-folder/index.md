---
id: "110"
title: Folder-form entity whose PR merged (non-armed)
status: implementation
score: "0.50"
source: roadmap
pr: pr-merge:111
---
# Folder-form entity whose PR merged (non-armed)

A FOLDER-FORM entity ({slug}/index.md) whose PR landed, recorded as a
`pr-merge:{number}` sentinel, with an EMPTY mod-block. `merge guard` must FINALIZE
it — terminalize + archive the whole folder to `_archive/{slug}/` — and commit the
move PATH-SCOPED. This pins the bug where commitArchiveMove hardcoded the flat
`{slug}.md` paths and exit-128'd the git add on a folder-form entity, stranding it.
