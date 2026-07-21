---
id: q4pw3xb4nf4cwfdjtwbn17mz
title: Retire legacy TeamCreate path and rename back-channel to inter-agent communication
status: ideation
source: captain (CL), 2026-07-20 session
started: 2026-07-21T16:05:13Z
completed:
verdict:
score:
worktree:
issue:
---

Retire the legacy TeamCreate machinery from the shipped FO contract and stop naming the worker/FO messaging surface "back-channel".

## Problem

Two contract-surface cleanups, captain-directed:

1. **Legacy TeamCreate path.** `claude-fo-dispatch.md` carries a legacy-override line that makes every session's first dispatch probe `ToolSearch(select:TeamCreate)`; when it matches, the FO loads `skills/using-legacy-claude-team`. Current runtimes no longer expose TeamCreate (this session's probe: no match), so the probe is dead weight on the first-dispatch path, and `using-legacy-claude-team/SKILL.md` is one of the 13 FO-surface measured files — retired, it shrinks the measured prompt surface. The claude runtime adapter and the terminal-teardown prose also carry legacy-override references that ride only on this trigger.

2. **"Back-channel" naming.** The contract names inter-agent messaging "back-channel" throughout (`## Worker Back-Channel` in claude-fo-dispatch.md, Codex "mailbox back-channel", Pi "contact_supervisor/intercom back-channel", «addressable-worker» prose). Captain direction: do not name it "back-channel" — name it for what it is, e.g. "inter-agent communication" (final term chosen at ideation).

## Degraded Mode retirement — split to `9q4`, but this task must not land blind to it

The captain asked (2026-07-21) whether this task should also retire Degraded Mode. A two-seat adversarial scope review says the retirement belongs in its own entity — `9q4` (dispatch-failure-retry-rung) — because the defect is the missing retry rung BELOW the trigger, not the trigger alone, and because this task's own `## Out of scope` promises no behavior change to dispatch or messaging. Two findings from that review bind THIS task regardless:

1. **This task deletes the repo's only retry rung.** `using-legacy-claude-team/SKILL.md:50` — «legacy-team.recover» rung 1, "Attempt one new TeamCreate with a fresh name" — is the sole try-once-before-degrading step anywhere in the contract. Retiring this file without `9q4` landing together or first leaves the contract with strictly ZERO retry surface, which is a regression, not a cleanup. Sequence them.

2. **VERIFY BARE-MODE REACHABILITY BEFORE DELETING THE LEGACY-OVERRIDE LINE.** A review seat claims the `e3z`-proven bare-mode entry (ToolSearch-no-match → bare) was inverted into the legacy-override at `claude-fo-dispatch.md:9`, which this task deletes — and that clause (3) of the Degraded Mode trigger is the only remaining route into `bare_mode: true`. That reading is CONTESTED: `claude-fo-dispatch.md:7` carries an independent `SendMessage`-availability probe that falls back to "fresh one-shot dispatch" without touching the legacy line. The two are not obviously the same state (`«addressable-worker» ABSENT` vs. the `bare_mode: true` build flag). **Resolve this at ideation with a live no-team drive, not by reading.** `e3z` proved bare mode reachable through the pre-retirement contract; nothing yet proves it reachable through the post-retirement one, and the coverage vacuum below means CI will not catch a mistake.

3. **`internal/claudeteam/claudeteam.go:68-79` (`BareModeAdvisory`)** instructs the FO to run `ToolSearch select:TeamCreate` on every bare dispatch lacking recent team evidence. After this retirement that names a tool which no longer exists, on the legitimate bare path — a required co-edit, and Go source, so it widens this task beyond prose.

4. **No oracle exists in either direction.** `TestLiveDegradedBareRecovery` is `//go:build live` and appears in ZERO CI `-run` filters; `.github/workflows/runtime-live-e2e.yml:111` sets `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`, so the teams-unavailable branch never runs in CI.

## Proposed approach

{Ideation fills this in. Expected shape: delete the legacy-override line and the using-legacy-claude-team skill; update the FO-surface file list and reference-closure expectations in contractlint; rename the concept across claude-fo-dispatch.md, fo-dispatch-core.md capability prose, and the codex/pi adapter bullets. Both ratchets must account for the surface shrink honestly.}

## Out of scope

{Ideation fills this in. Likely: no behavior change to dispatch or messaging itself — naming and dead-path removal only.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - No TeamCreate probe or legacy-team path remains in the shipped contract; first dispatch has no legacy branch.**
Verified by: contractlint reference-closure/structural checks after the skill's removal; the FO-surface measured-file list no longer contains `using-legacy-claude-team/SKILL.md`; review-time grep cited as evidence.

**AC-2 (VALUE) - The FO prompt surface shrinks by at least the retired skill's byte size.**
Verified by: `TestFOFunctionPromptSurfaceShrinks` / checkpoint-metrics byte counts before vs after, an independent number that can move the wrong way.

**AC-3 - The shipped contract nowhere names inter-agent messaging "back-channel"; the concept carries a descriptive name.**
Verified by: review-time grep over the FO surface cited as evidence (no committed prose-grep, per proof policy).

**AC-4 - The FO-surface ratchet baseline is re-tightened to the post-retirement measured value (zero-slack convention).**
This task carries the F1 mitigation from `codex-post-compaction-contract-reload` (archived, validation cycle 7): the scratch-proven inert-heading escape passes every mechanical gate whenever slack exists, leaving the review-time grep as the only guard for inert-prose reintroduction. Re-tightening here closes it mechanically.
Verified by: `foFunctionReferenceBaselineBytes` equals measured+1 after the retirement lands; `TestFOFunctionPromptSurfaceShrinks` green at zero slack.

*Figures updated 2026-07-21 — the "204 B of slack at baseline 122634" this AC originally cited is stale.* Main measures **122126** against baseline **122634** (507 B of usable headroom under the strictly-below gate). Member `bw` re-baselined to **123323** on 2026-07-20 as a captain-approved governance decision after a duplication scan across all 13 measured files recovered only 110 B; that branch is complete but unmerged, so re-measure at implementation rather than trusting either number here. Retiring `using-legacy-claude-team/SKILL.md` alone removes **14065 B** from the measured surface.

## Test plan

{Ideation fills this in. Expected: contractlint structural suites + surface-metrics deltas; claude-live lane required (shipped contract change); no committed grep tests.}
