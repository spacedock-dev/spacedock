---
id: nxnbzzw9z60tx4bp3daeypkc
title: Make stacked-PR delivery a supported path in the pr-merge mod
status: validation
source: "Captain, 2026-08-20, after a stacked delivery of three entities exposed a parallel layer, a hazardous trunk push, a contradicted conflict rule, and a mod whose installed copy lacked the whole section."
started: 2026-08-20T17:16:44Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-support-stacked-pr-delivery
issue:
pr: "#743"
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

## Stage Report: implementation

- DONE: Write the implementation stage report for the delivered stacked-mode section, from the committed diff
  Read from `git diff origin/main...HEAD -- mods/pr-merge.md` (+45/-6) across 7994539f2, 8a9b12cda, 8f29ad577; each element below was located in the diff text, not taken from the dispatch's list.
- DONE: `### Stacked mode` defers every `gh stack` mechanic to `gh skill preview github/gh-stack gh-stack` and declares three overrides
  7994539f2: no `gh stack submit` (auto-titles, so approved bytes never reach GitHub); `gh stack link` takes only PR numbers confirmed by `gh pr view {N} --json number`; no local stack tracking, so `rebase`/`sync`/`push`/`view` do not apply.
- DONE: three-condition test that a layer contains its parent — ancestry, heads-not-equal, parent-holds-own-work
  Exercised read-only on this entity's own live layer (PR #743 over parent `spacedock-ensign/force-boot-at-compaction-boundary`): `merge-base --is-ancestor` exit 0, heads differ exit 0, `rev-list --count "$PARENT_HEAD" --not "$TRUNK_SHA"` = 8. The heads-not-equal condition is the one a freshly branched parallel layer fails, since a commit is its own ancestor.
- DONE: the trunk push is skipped in stacked mode
  "Use the branch below the layer as `$BASE` … and skip the trunk push on approval. That push sends the parent layer, including commits the captain has not approved."
- DONE: rebase procedure with `rebase --onto`, a bare lease, and `CANDIDATE_SHA` re-recorded
  Steps 1-7: `OLD_PARENT` read before the fetch, `rebase --onto "origin/$PARENT_BRANCH" "$OLD_PARENT" "$BRANCH"`, step 5 replaces `CANDIDATE_SHA` with `$NEW_HEAD`, step 7 pushes `--force-with-lease --force-if-includes` with no value on the lease (8f29ad577).
- DONE: conflict rule routes to the recorded owner of the layer's registered branch and worktree
  Step 3 aborts, mutates no refs, and holds the entity with its pending approval and `mod-block` rather than taking `--rework`; the diff text itself does not cite #645 — that citation lives in this entity's Problem section and in 7994539f2's message.
- DONE: why a stack exists, where its checks are approved, and the check-latency rule
  8a9b12cda: one tip run proves every layer beneath it; a middle-layer run needs the captain plus a named falsification the tip cannot reach; an empty check list does not prove the repository runs none (wait 30s, stop after three empties, start no duplicate run).
- DONE: two front-half lines shipped outside the declared stacked surface
  7994539f2 also made "Always present the draft" explicit and added that a standing conn counts as push approval only when the captain's own words grant push or PR authority — a change to the unstacked approval path, not to `### Stacked mode`.
- SKIPPED: Sync the installed mod with the repository copy, and make same-version drift detectable rather than silent
  Did not ship: `docs/dev/_mods/pr-merge.md` still has no `### Stacked mode` heading, both copies still declare `version: 0.27.0`, and the two differ by 57 diff lines — the entity's second acceptance clause is unmet.
- SKIPPED: Test whether `gh stack rebase --downstack` runs inside a layer's own worktree
  Untested, so the hand-rolled seven-step procedure ships rather than deleting itself.
- SKIPPED: Decide whether the ancestry check earns a mechanical guard
  Undecided; a standing guard is its own captain decision and none was sought. Nothing mechanical covers this section: `grep -rn "Stacked mode\|force-if-includes" --include=*.go` matches no test.
- DONE: Three commits removed prose the new section does not replace
  The `gh api --method PATCH` title/body repair note (and its `gh pr edit` warning), the `gh stack link` HTTP 422 read-back paragraph, and the "back half needs no stacked case" statement are gone; only the link read-back survives, condensed to a base check.

### Summary

The deliverable is 45 lines of prose in `mods/pr-merge.md`, written by the first officer under a captain direct-edit grant, and there is no mechanical check behind any of it — no validator ran, no test asserts its text, and this entity carries no prior gate record (`verdict` and `score` are empty, and there are no Cycle lines). The dispatch cites two independent adversarial reviews; those left no artifact in the repo or the entity, so I cannot confirm them. What the git record does show is the three live exercises finding two real defects in the text: 7994539f2 names the ancestry hole an equal-heads test closes, and 8f29ad577 records that step 7 was rejected with "stale info" because a named lease blocks your own push. The installed-copy sync, the drift check, the `--downstack` experiment, and the mechanical-guard decision all remain open; the captain approved this entity to terminal with no further validation, so they are recorded here rather than fixed.
