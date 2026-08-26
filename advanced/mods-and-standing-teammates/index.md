---
title: "Mods & standing teammates"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-26 17:14:19"
---

# Mods & standing teammates

As a workflow matures, you start wanting behavior every workflow wants: create a PR and act on it when it merges or CI fails, or keep a specialist whose judgment persists across the whole session (a prose polisher). You could model such steps as stages of your own; they are common enough across workflows that they are factored out as mods instead. A mod is the hook mechanism: standardized behavior that hooks into the workflow's run without changing your workflow definition. The stages and gates in your README stay as they are; a mod is one markdown file in the workflow's `_mods/` directory, and the first officer reads it and acts on it.

## Lifecycle hooks

A hook adds a step the first officer performs at a fixed point in the run: at `startup`, on an `idle` pass when nothing is ready to dispatch, or at the `merge` boundary when an entity reaches its final stage. Hook points are workflow-independent: any workflow can register the same hook to get the same behavior.

The canonical example is the [`pr-merge` mod](https://github.com/spacedock-dev/spacedock/blob/main/docs/dev/_mods/pr-merge.md): it opens the code-branch PR at merge, records the PR on the entity, and holds the terminal transition until the PR merges. The block is enforced; a half-merged entity cannot slip past the gate.

## Standing teammates

A standing teammate is a long-lived specialist agent declared by a mod. It stays available for the session and is addressed by name. Reach for one when the same specialist judgment recurs across entities and is worth a persistent agent rather than a fresh dispatch each time.

The canonical example is the [**comm-officer**](https://github.com/spacedock-dev/spacedock/blob/main/docs/dev/_mods/comm-officer.md), a prose-polisher the first officer routes deliberate drafts through (PR bodies, gate summaries, debriefs) before they reach you. Routing is best-effort: if the teammate is absent or slow, the work proceeds without it.

Ask the first officer to install a shipped mod or write a new one; the file format is its job.

## Sitemap

- [Multiple workflows](../multi-workflow/index.md)
- [Split-root state](../split-root-state/index.md)
- [Bridge an external tracker](../external-tracker/index.md)
- [Refit a workflow](../refit/index.md)
