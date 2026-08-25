---
id: dt8j3pas83725fma2wbez5ss
title: Product README documents no edge channel path
status: ideation
source: "Captain CL work-machine report 2026-08-25: README at v0.27.0 documents only brew tap/install (stable); someone following it for edge parity lands on stable and loses parity silently"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:dt8j3pas83725fma2wbez5ss:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:dt8j3pas83725fma2wbez5ss-backlog-1
              briefing:
                id: briefing:dt8j3pas83725fma2wbez5ss:backlog:attempt-1:revision-1
                digest: sha256:c35c57a4aa2757547aae765122da7e1cea5396fe267ae9c5f52959b6799c9161
                request-digest: sha256:70b91d9acb281b900da42fd590b89cab6d58cb951ea020386a137234aa5709ab
                room-ref: ./product-readme-edge-channel-mention/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:dt8j3pas83725fma2wbez5ss:backlog:1
                briefing: briefing:dt8j3pas83725fma2wbez5ss:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T14:34:27.001703Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''file and dispatch'' — approved seeding and immediate dispatch of this task into design; README posture decision reserved to the captain at the ideation gate'
              application:
                target-stage: ideation
                state: consumed
---

The product README's install section documents only `brew tap spacedock-dev/homebrew-tap` + `brew install spacedock` — the stable channel. No edge cask, no `SPACEDOCK_CHANNEL=edge`, no pointer to the edge marketplace. A reader seeking edge parity follows the README, lands on stable, and loses parity with no signal. Decide the intended README posture and ship the README diff.

## Problem

README.md lines 60-64 (verified 2026-08-25) show only the stable brew steps. The archived task `marketplace-readme-channel-model-repair` (done 2026-08-24) fixed the MARKETPLACE repo README and docs/site install page; the product README was out of its scope. CLAUDE.md frames edge as a dev convenience, so whether the product README should advertise edge at all is a posture decision the captain holds. {Ideation refines.}

## Proposed approach

{Ideation fills this in. Seeded: two candidate postures — (a) document the edge install path in the README (edge cask / SPACEDOCK_CHANNEL=edge plus the marketplace@edge plugin add), or (b) state explicitly that the README covers stable only and link the docs-site install page for edge. Captain rules on the posture at the ideation gate; ideation records the concrete before/after doc diff per the workflow's doc-diff rule.}

## Risk evidence

{Backlog: the captain's silent-parity-loss report plus the verified README lines are the observation that decides design should start. Likely `no spike needed`: doc-only change over proven install commands — ideation confirms the commands it documents against the shipped install machinery.}

## Out of scope

The marketplace repo README (done: marketplace-readme-channel-model-repair). install.sh behavior. CLAUDE_CODE_PLUGIN_PREFER_HTTPS documentation (plugin-install-https-keyless-machines / rb0). The claude sibling-plugin cleanup (claude-install-sibling-channel-cleanup).

## Expected surface and tolerance

Estimate net LOC change: +15, across 2 files. {Backlog seed; ideation refines with tolerance.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A reader following the product README toward the edge channel reaches an edge install, or is told explicitly that the README path is stable-only and where edge lives — never silently landed on stable.**
Verified by: {ideation refines — seed: the shipped README diff naming the edge path or the explicit stable-only statement with link. Presence is the claim here, so a one-off existence check recorded in the validation report is legitimate evidence per the Proof policy. Falsifying edit: revert the section to the bare stable brew steps.}

## Test plan

{Doc-only; no code tests expected. Validation small-change fast path applies. If posture (a) is chosen, ideation states how the documented edge commands are verified against a real host once.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
