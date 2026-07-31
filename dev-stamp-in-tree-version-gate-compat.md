---
title: "Dev-stamp compat: in-tree builds must pass the FO version gate against in-tree skills, for dev usage and CI"
status: ideation
source: "Captain directive 2026-08-01, in the PR #586 merge decision: file the dev-stamp problem separately from z3, 'tackling actual dev usage and ci'. Same-class bites in one day (2026-07-31): (1) FO boot self-aborted on a stale 0.26.0+dev env binary against an in-tree 0.27 plugin (captain override needed); (2) a gpt-5.6-sol staff reviewer self-blocked on the same prose mid-review; (3) codex-live keep-moving-posture + gate-guardrail on PR #586 failed twice, the live FO declaring the PR candidate's 0.27.0-pre2+dev stamp 'not a non-development 0.27 build' — while the same scenarios stayed green on main and on claude/pi lanes."
issue: spacedock-dev/spacedock#581
id: zexbrjhartgykvhm012f527w
gates:
    version: 1
    current:
        gate: gate:zexbrjhartgykvhm012f527w:backlog
    records:
        - id: gate:zexbrjhartgykvhm012f527w:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zexbrjhartgykvhm012f527w-backlog-1
              briefing:
                id: briefing:zexbrjhartgykvhm012f527w:backlog:attempt-1:revision-1
                digest: sha256:0e8e9d466ebe2768f6c01a89a402a427d03ad1ed03ef5080a4b111e2ac0ec2b7
                digest-domain: canonical-bytes
                request-digest: sha256:cd777a48b4365012f90310e3cda8ae4e29c3d20f2439dc16bd00db91141bd4d0
                room-ref: ./dev-stamp-in-tree-version-gate-compat/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zexbrjhartgykvhm012f527w:backlog:1
                briefing: briefing:zexbrjhartgykvhm012f527w:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-31T16:32:37.288799Z"
                decision: approve
                reason: Under the captain's explicit conn to proceed with the six critical lanes, the bound backlog Briefing demonstrates a repeated live-agent failure, constrains ideation to the smallest sufficient correction, and requires behavioral proof without compatibility machinery.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

The FO version gate's dev-stamp abort class does not survive contact with in-tree builds: source builds and CI candidates stamp `<latest-tag>+dev`, which IS major.minor-compatible with the in-tree skills, but the gate prose names the dev-stamp class only by its extreme (`dev` — an integer-era build carrying no major.minor at all) and never states the converse, so strict live-agent readers classify compatible dev-suffixed tokens as the abort class. Evidence is deterministic: codex's runtime red twice on identical reads where claude/pi read green; main's lanes with the old prose were green 07-13/07-19.

Ideation should firm the shape, but the seed classes are: (a) prose disambiguation for agent readers — a token that parses a major.minor is never the dev-stamp class; build metadata and prerelease suffixes do not make it one (this is the smallest codex-lane fix; z3's PR #586 shipped the rewritten step-1 prose without it); (b) an enforced invariant that an in-tree build of any tree state yields a --version token the co-located skills accept — machinery in the spirit of stamp-then-tag-release-ritual's 'enforce the invariant in machinery, not just prose', covering both actual dev usage (stale env binaries: detect-and-rebuild guidance, not a hard fault) and CI (a guard/lane assertion, with the codex-live keep-moving-posture scenario as the standing live regression). Related: stamp-then-tag-release-ritual, next-independent-release-line, fo-boot-install-hint-linux-direct-sandbox (z3), the #468 minor-coupling gate. Proof future stages must provide: a live-agent reading of every realistic in-tree stamp shape passes the gate (fixture or live lane, not prose-grep), plus the compat invariant exercised mechanically.
