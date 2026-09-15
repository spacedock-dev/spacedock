---
title: Honor workflow-declared same-stage revision without mandatory reviewer machinery
status: ideation
source: Captain request after email-triage FO issue 792
issue: spacedock-dev/spacedock#792
score: 0.95
started: 2026-09-15T04:30:03Z
completed:
verdict:
worktree:
pr:
mod-block:
id: zz1yqc2w2katp28wpa8nghx2
gates:
    version: 1
    records:
        - id: gate:zz1yqc2w2katp28wpa8nghx2:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zz1yqc2w2katp28wpa8nghx2-backlog-1
              briefing:
                id: briefing:zz1yqc2w2katp28wpa8nghx2:backlog:attempt-1:revision-1
                digest: sha256:ea21af7990ca6e4a52dac540a434ad088b3bbecca35dcb03cbd1c70d39950ea6
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:zz1yqc2w2katp28wpa8nghx2:backlog:1
                briefing: briefing:zz1yqc2w2katp28wpa8nghx2:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T04:29:46.231251Z"
                decision: approve
                reason: Captain explicitly requested local filing and ideation dispatch to spike actual same-stage revision and inspect all related Roborev settings/history.
              application:
                target-stage: ideation
                state: consumed
---

Make same-stage gate revision follow the declared workflow while preserving independent review where it is required.

## Problem

Issue https://github.com/spacedock-dev/spacedock/issues/792 reports an email-triage workflow with feedback-to equal to its own gated stage, no independent reviewer stage, and no canonical review-round log. The shared feedback-rejection-flow nevertheless mandates both round publication and reviewer rerun. This can manufacture absent machinery or stall legitimate correction.

The captain suspects same-stage revision already shipped alongside Roborev and the development workflow's implementation-stage review. Do not assume new binary machinery is needed. The skill is unchanged from v0.27.2 to v0.28.0-pre3; the five-step rewrite at 884b55af0 retained earlier assumptions.

## Proposed approach

First trace all relevant Roborev configuration, hooks, workflow settings, same-stage gate/revise semantics, and introduction history. Inspect .roborev.toml, docs/dev/README.md, related local tasks, shipped skills and binary/tests at current origin/main. Distinguish shipped behavior from pending task designs and installed configuration. Related tasks: roborev-validation-hook, roborev-workflow-setup-skill, portable-state-checkout-roborev-followup, workflow-neutral-advisory-round-recorder, simplify-feedback-rejection-flow. The neutral-recorder task owns data interpretation; this task owns when the generic flow invokes that recorder/reviewer.

Spike same-stage revise using an actual on-disk workflow and real public CLI operations, not only prose or a state-machine mock. Use a disposable isolated copy of a real workflow definition and the installed/current binary; retain exact config/version, commands, exit codes, Git state, old/new attempts, reports, and approval state. Exercise both a triage-like self-feedback gate without a separate reviewer/log and the existing dev/Roborev same-stage path. Preserve the existing conventional producer-to-reviewer path as a control. If an agent drive is needed, coordinate the serialized local live slot with the FO; do not contend with the Codex stack runs. No Gmail writes, production workflow mutation, candidate implementation, or generic framework during this ideation.

Select the smallest fix from observed evidence. Likely scope is conditional shared-skill obligations based on the workflow's declared review/round contract, independently of whether a Feedback Cycles projection exists. Reuse existing same-stage mechanics if proven. Do not add flags, metadata, reviewer stages, or logs just to satisfy generic prose.

## Risk evidence

The issue provides a concrete same-stage triage incident on pre2+dev. Confirm exact current semantics and the Roborev history before asserting the cause or design. Record negative cases proving missing required review/round evidence still blocks the development path, and that stale or rejected attempt authority cannot approve the corrected plan.

## Out of scope

Codex release-test stack fixes, archive durability (#790), terminal output clarity (#791), AC scanner naming (#793), neutral-recorder schema/taxonomy changes, and live email execution.

## Expected surface and tolerance

Ideation only for now: this task body, its retained spike artifacts, and disposable workflow checkout. Propose exact product files/net LOC/tolerance after the spike; do not edit shipped code or skills yet.

## Acceptance criteria

**AC-1 — A self-feedback gate can correct and re-present without invented reviewer or round artifacts.**
Verified by: a real workflow drive records revise, dispatches/completes the same declared stage, commits corrected artifacts, and prepares exactly one fresh open gate. It preserves the frozen input and old attempt history, and performs no approval, consume, or external execution on the new attempt.

**AC-2 — Existing Roborev and development review obligations remain enforced.**
Verified by: the actual declared same-stage implementation/Roborev path and conventional separate-reviewer control complete required review/round steps; omitting required review or recorder evidence fails. Settings/history establish which behavior already ships versus what is only proposed.

**AC-3 — The fix changes only the proven generic-routing mismatch.**
Verified by: a scoped before/after workflow exercise shows the self-feedback path progresses with existing mechanics while required-review controls remain strict. The design names exact owning files and a falsifiable test for each changed behavior; no prose-grep stands in for workflow evidence.

## Test plan

Ideation must first reproduce the real same-stage workflow and inspect all relevant Roborev settings. Use existing CLI/gate fixtures and a disposable real workflow, retaining its artifacts and Git history. Distinguish binary workflow proof from live-agent proof. Before implementation approval, present the minimal tested design and exact future targeted regression/live commands; full repository tests follow implementation, not this ideation-only commission.

### Feedback Cycles

