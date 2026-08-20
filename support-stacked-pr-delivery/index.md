---
id: nxnbzzw9z60tx4bp3daeypkc
title: Make stacked-PR delivery a supported path in the pr-merge mod
status: implementation
source: "Captain, 2026-08-20, after a stacked delivery of three entities exposed a parallel layer, a hazardous trunk push, a contradicted conflict rule, and a mod whose installed copy lacked the whole section."
started: 2026-08-20T17:16:44Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-support-stacked-pr-delivery
issue:
---

Make a stacked delivery safe to drive from the `pr-merge` mod alone, and stop the mod from carrying tool knowledge that the official `gh stack` skill owns.

## Problem

Driving one real stack of three entities surfaced four defects in one session.

A layer was created from its parent before the parent's worker had committed, so both branches forked from the same base. The layer declared the parent as its PR base and contained none of the parent's work. Every available signal agreed it was healthy: the base was correct, `merge-tree` was clean, GitHub reported it mergeable, and the PR diff listed only the layer's own files. A CI run on that tip would have exercised the upper change alone and reported green. The captain caught it; nothing in the mod could.

The mod's front half pushes the trunk unconditionally on approval. In stacked mode `$BASE` is the parent layer, so that push sends the parent branch, including commits the captain has not approved.

The mod told the agent to follow the official skill for every `gh stack` mechanic and then hand-rolled the rebase and push, without recording why. An agent obeying both documents stalls, because the skill prescribes `gh stack rebase` and `gh stack push`.

The conflict rule contradicted `codify-conflict-owner-dispatch-handoff` (`D8`, merged as #645): the mod sent a moved base to a halt, while the shipped contract routes exactly that case to the worker recorded for the layer's registered branch and worktree.

Underneath all four: the installed `docs/dev/_mods/pr-merge.md` carries no stacked section at all while declaring the same `version: 0.27.0` as the repository copy. Version strings cannot distinguish them, so nothing detects the drift.

## Proposed approach

The mod already carries a corrected `### Stacked mode`, written under a captain direct-edit grant and reviewed twice by an independent engineer. This task delivers that text through the normal path and closes what the reviews left open.

- Sync the installed mod with the repository copy, and make same-version drift detectable rather than silent.
- Keep the mod deferring to `gh skill preview github/gh-stack gh-stack` for tool mechanics, and keep only this ceremony's own rules plus its three declared overrides.
- Test whether `gh stack rebase --downstack` runs inside a layer's own worktree, where that layer's branch is legitimately checked out. `gh stack link` creates no local tracking and `gh stack checkout` cannot claim a branch another worktree holds, which is why the hand-rolled rebase exists. If the worktree path works, the seven-step procedure deletes itself.
- Decide whether the ancestry check earns a mechanical guard. It is pure local git, needs no network or forge, and separates a parallel layer from a stacked one in one command. A guard is a new standing check and needs captain approval as its own decision.

## Out of scope

Do not restate the official skill. Do not change the gate ceremony, the merge guard, the PR body template, or the conflict-owner handoff contract that `D8` shipped. Do not add a resolver.

## Expected surface and tolerance

Ideation sets this. The delivered text is 45 lines in `mods/pr-merge.md`. The installed-copy sync and a drift check are the additions.

## Acceptance criteria

Ideation writes these. They must include the ancestry check distinguishing a parallel layer from a stacked one, and evidence that the installed copy no longer diverges from the repository copy at the same version string.

## Test plan

Ideation writes this. The ancestry cases are offline and cheap: build a parallel pair and a stacked pair and grade both.
