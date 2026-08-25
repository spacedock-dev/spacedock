---
id: d0g21c517b5nvga1ybwckapk
title: Patch-release line support - gate release.yml's main-stamp on line-latestness, fix the stable-branch advance, automate the preversion bump
status: backlog
source: "Captain CL, 2026-08-25, reconciling gr and tw after the v0.27.0 cut: 'reconcile tw and gr and recommend best approach' - supersedes next-independent-release-line (twq68r4y8qg0wetztajtmmzz), whose body described the retired next-branch model; the live incidents of 2026-08-25 are the spec"
started:
completed:
verdict:
score:
worktree:
issue:
related: "next-post-release-preversion-bump (closed delivered); stamp-then-tag-release-ritual"
gates:
    version: 1
    records:
        - id: gate:d0g21c517b5nvga1ybwckapk:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:d0g21c517b5nvga1ybwckapk-backlog-1
              briefing:
                id: briefing:d0g21c517b5nvga1ybwckapk:backlog:attempt-1:revision-1
                digest: sha256:6217483748c63c0663f4d6a6b96d8b90277d19e216e08534a10f5fca95168c2f
                request-digest: sha256:7c1dce1d1eb8b5e2720ce244ea77f986847eec757358eed954c21f83a61208e1
                room-ref: ./patch-release-line-support/review/backlog/briefing-1
---

One owner for release.yml's stable-tag conditional. Today the machinery cannot ship an old-line patch (a v0.27.1 while main is the 0.28 line): two step-level defects, verified by reading the shipped workflow at the v0.27.0 cut, plus one automation the cut proved missing.

## Problem

{Ideation fills this in. The three items, from the 2026-08-25 analysis:
1. The "Stamp plugin manifests to the release version" step gates only on "stable tag" (no hyphen in the ref), not on line-latestness. On a v0.27.1 tag it switches to main and stamps main DOWN to 0.27.1 - rewriting the 0.28.0-pre0 manifests and the FO minor pin, re-breaking every edge user backwards. The latest-line decision already exists in the sibling edge-advance job (highest-known-edge-version + edge-advance-decision) but this step does not consult it.
2. The stable-branch advance pushes the tagged commit to refs/heads/stable with the comment "the tagged commit is on main's history, so this fast-forwards" - false for a release-branch commit, so the push dies non-FF and stable-channel plugin users never receive the patch even though goreleaser published its binaries.
3. Nothing automated stamps main past the released minor after a latest-line stable cut. The auto-pre0 job tags vX.(Y+1).0-pre0 but leaves main's manifests and FO pin at the released version; every edge install aborts at the FO version gate until a human runs the bump (observed live: the v0.27.0 cut, the 05:01Z install at 17f5cd591, the mismatch loop, the manual b8346ffc9 stamp). docs/releasing.md step 9 (b04b3effd) is the interim; this item retires it.}

## Proposed approach

{Ideation fills this in. Seeded: reuse the existing latest-line decision for all three - latest-line stable tag: stamp main to X.(Y+1).0-pre0 (not the release version, which the ritual's step-3 commit already carries), advance stable to the tagged commit (FF by construction); old-line stable tag: skip the main stamp entirely, advance stable ONLY if the tagged version exceeds the stable branch's current version (mechanism for a non-FF-safe advance is an ideation decision - force-with-lease against a verified ancestor tag, or refuse-with-instructions). Prove the riskiest claim first: that the decision helper correctly ranks an old-line patch against a newer prerelease (its comment says it was built for exactly this).}

## Out of scope

The e2e-gate/waiver mechanics (work for branch SHAs today). The marketplace display-version fields (cosmetic; tidied 2026-08-25, automation optional here). The release notes ritual.

## Expected surface and tolerance

Estimate net LOC change: ~+60 across 3 files (.github/workflows/release.yml, cmd/spacedock-release decision surface if extended, docs/releasing.md step-9 retirement). Ideation refines with tolerance.

## Acceptance criteria

Seeded; ideation refines and re-anchors.

**AC-1 (value) - An old-line patch tag publishes without touching main: main's manifests and FO pin are byte-identical before and after a vX.Y.(Z+1) release run on a non-latest line.**
Verified by: a workflow-level test or act/fixture replay of the stamp step's conditional with an old-line ref, asserting no main mutation; fails if the step stamps main on a non-latest tag.

**AC-2 (value) - A patch on the latest stable line reaches stable-channel plugin users: the stable branch advances to the tagged commit even when that commit is not on main's history.**
Verified by: a fixture repo replay of the advance step from a release-branch commit asserting the branch moves (or the refusal path prints its instruction, per the ideation decision); fails on the current non-FF die.

**AC-3 (value) - After a latest-line stable cut, an edge install passes the FO version gate with no human step.**
Verified by: the release run's main tip carrying the X.(Y+1).0-pre0 stamp (manifests + FO pin) at run completion, asserted against a replay or the next live cut; fails if the gap window needs a hand stamp. Retires docs/releasing.md step 9 in the same change.

## Test plan

{Ideation fills this in. Seeded: the decision-helper ranking test (old-line vs newer prerelease); workflow-step replays over fixture repos rather than live cuts; one live observation rides the next real release.}
