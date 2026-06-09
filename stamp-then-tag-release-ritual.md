---
id: ezn308z0chwc2zvmyny9ry8w
title: Stamp-then-tag release ritual so the tagged commit's manifest matches its tag
status: backlog
source: "Post-flip release-model decision (roadmap 0201, captain 2026-06-09). Surfaced by the flip: v0.20.0 -> commit 6accd320 whose plugin.json says 0.19.9 (the 0.20.0 stamp landed 2 commits later on main). Harmless while serving branch HEAD; breaks tag-serving (Model B)."
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: release-model
sprint-readiness:
---

Invert the release ordering to **stamp-then-tag** so the commit a release tag points at carries the matching `plugin.json` version, and pin `docs/releasing.md` to the post-flip reality. See `docs/roadmap/0201-post-flip-release-model/index.md`.

## Problem

The flip placed the `v0.20.0` tag on the green tip (`6accd320`, whose `plugin.json` says **0.19.9**) and stamped `0.20.0` two commits **later** on `main` (`3a694742`). So `git show v0.20.0:.claude-plugin/plugin.json` reports `0.19.9`. This is harmless while the stable channel serves a branch HEAD, but **it breaks Model B's tag-serving**: a stable entry pinned to `ref: v0.20.0` would serve a manifest that says 0.19.9 wearing a 0.20.0 tag. Separately, `docs/releasing.md` (reconciled for the flip) still describes a main-integration cut and is not pinned to the recurring "stamp → tag → repoint" ritual.

## Proposed approach (ideation firms)

- Invert to **stamp-then-tag**: bump `plugin.json` version → commit → one live-e2e (the existing `e2e-gate`) → **tag that commit** `v0.X.Y` → repoint the stable marketplace entry's `ref` to the tag. The same tag drives the binary release (unchanged) and the stable plugin repoint — one tag, two consumers.
- Enforce the invariant in machinery, not just prose: a release-time guard/test that the tagged commit's `plugin.json` version equals the tag's semver.
- Rewrite `docs/releasing.md` to the recurring ritual (stamp-then-tag, two-consumer tag, the clean `next → main` fast-forward once the marketplace repo lands).
- Delete the now-dead `version` field on the `next` marketplace entry — `plugin.json` wins silently over it (https://code.claude.com/docs/en/plugins-reference), so it is a footgun, not a control.

## Out of scope

- The marketplace repo structure + channel entries — separate task `marketplace-repo-and-pinned-channels`.
- The install-refresh correctness + upgrade hint — existing task `tes`.

## Acceptance criteria

(Ideation firms. Verified by an independent value that can diverge from the file under test — not a prose-grep.)

**AC-1 (sketch) — the commit a release tag points at carries the matching manifest version.** Verified by: a release-machinery guard/test asserting the tagged commit's `plugin.json` version equals the tag's semver — it FAILS on the current backwards ordering (`v0.20.0` -> 0.19.9) and passes after the inversion. This compares two independent values (the git tag and the manifest), which can genuinely disagree, so it is a real invariant, not a tautology.

**AC-2 (sketch) — the dead `next` marketplace `version` field is removed.** Verified by: its absence in the marketplace manifest + a test/lint that the served version derives from `plugin.json` (resolved-version state), not the marketplace entry.

## Test plan

(Ideation/implementation firms.) A guard over the cut machinery (or a fixture release) asserting tag-commit manifest == tag semver; a check that the served/resolved version comes from `plugin.json`. `docs/releasing.md` is real authoring but is a deliverable, not an AC on its own — the enforceable invariant (AC-1) is the proof.
