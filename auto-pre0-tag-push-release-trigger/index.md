---
title: Ensure the automatic pre0 tag push triggers its release workflow
status: ideation
source: "v0.25.0 stable cut, 2026-07-15: edge-advance pushed annotated v0.26.0-pre0, but no release run appeared; manually replaying the identical tag triggered release.yml immediately."
started: 2026-07-21T15:58:31Z
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
id: 5aqczjeq6rq3mckbc5gyjqe3
---

A stable release must publish its automatically generated next-minor edge binary without requiring an operator to delete and replay the pre0 tag.

## Problem

The stable `edge-advance` job created and pushed the correct annotated `vX.(Y+1).0-pre0` tag, and the remote tag pointed at the intended green release commit. No `release.yml` run was created for that push. Deleting and manually re-pushing the identical annotated tag immediately created the expected run, which then published the edge release and cask successfully.

The current job can therefore report success while leaving the edge binary behind the newly advanced `next` skills.

## Reproduction

1. Push a stable `vX.Y.Z` tag and let the release workflow reach `edge-advance`.
2. Observe the job create and push annotated tag `vX.(Y+1).0-pre0` using its configured release token.
3. Confirm the remote tag exists on the stable release commit.
4. Observe that no release workflow run exists for the pre0 tag.
5. Delete and manually re-push the same local annotated tag; observe that the release workflow starts immediately.

## Acceptance criteria

**AC-1 (VALUE):** A stable cut automatically publishes the next-minor pre0 edge release and updates the edge cask without an operator replaying the tag.

Verified by: a release-workflow fixture or controlled repository exercise that observes the pre0 tag push produce a distinct `release.yml` run on the same commit and the edge artifact publication complete.

**AC-2:** The credential and event path used by `edge-advance` is capable of triggering tag-push workflows; a token shape that GitHub suppresses or cannot trigger fails before the stable run reports the edge handoff complete.

Verified by: an integration check over the configured push mechanism, plus a negative control using a workflow-suppressed credential or event.

**AC-3:** The pre0 tag remains annotated, non-empty, and bound to the stable release commit; recovery never changes its source tree or invents standalone changes.

Verified by: exact tag-object, peeled-commit, and release-body assertions before and after the automated handoff.

## Test plan

Ideation should first determine why the configured token push produced a remote tag without an Actions run. Prefer the smallest correction to the authentication or event path. Do not add a polling controller unless a direct trigger-capable push cannot satisfy AC-1. Exercise one controlled stable-to-pre0 handoff and retain the tag, run, release, and cask evidence.
