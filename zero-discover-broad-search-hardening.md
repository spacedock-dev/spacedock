---
title: Harden the zero-discover boot against the residual broad-search overstep (flaky guard)
status: backlog
source: "0204 Commander (2026-06-16 session #4), byte-verified across CI runs: TestLiveZeroDiscoverReportsAndStops flakily REDs on the claude-live lane when, after a zero `status --discover`, the boot FO runs a broad filesystem sweep (find/ls over the project root) instead of report-and-stop. Confirmed NOT a CC version cliff: the red occurs on 2.1.177 AND 2.1.179, and flips between sonnet/opus matrix legs of the SAME run at the SAME version — i.e. stochastic model-discipline, not a version step-function. The lean-boot-hardening task (#374, archived PASSED) already shipped the binary-OUTPUT fix (the no-workflow message reframed terminal) + the detectBroadSearchAtBoot stream-scanner guard, and explicitly concluded the contract PROSE was not the lever (the FO read 'report and stop' and swept anyway). So #374 reduced but did not eliminate the overstep; the guard now catches the residual. This is why the zero-discover red on a required CI lane is re-run-grounds (flaky guard), never a merge blocker."
sprint:
id: tq0jjq803c1a6hc757pjymj6
---

The lean-boot zero-discover guard (`detectBroadSearchAtBoot`, shipped by #374) flakily reds because the boot FO still, sometimes, broad-searches the filesystem after a zero `status --discover` instead of report-and-stop. #374 proved the binary-output reframe + a detector guard, and proved prose alone does not hold. What remains is a stochastic model-discipline gap with no deterministic fix yet.

## Problem
- The contract's Startup zero branch is terminal: zero `--discover` -> report no workflow found and STOP. #374 reframed the no-workflow CLI output to be self-evidently terminal AND added the `detectBroadSearchAtBoot` stream-scanner as the regression guard.
- The guard still reds intermittently: the FO issues `find … README.md | head && ls …` over the project root on a zero-discover boot. Byte-verified flaky: reds on 2.1.177 and 2.1.179; flips between matrix legs at one version. So it is model nondeterminism, not a CC version regression and not a 0204 read-cost regression.
- The prose lever is spent (#374's finding: the FO read the report-and-stop directive and swept anyway). The no-broad-search prohibition lives only in the test detector, not in the contract prose the FO acts on.

## What's needed (decide in ideation — do not pre-commit)
A durable reduction of the residual overstep, evaluated against the spend-the-cheapest-lever-first rule:
- an even-stronger terminal signal in the binary's zero-discover output (the lever #374 used, pushed further), and/or
- a contract-prose hardening that names the no-broad-sweep prohibition explicitly AND is backed by the existing detector as its code-gate (composes with fo-self-evidence-bar, which notes the same gap), and/or
- a flaky-lane policy: quarantine/auto-retry the zero-discover live step so a stochastic red is re-run, not merge-blocking, with the flakiness rate tracked.
Whichever is chosen must be proven by the existing live detector behavior over real boots, not a prose-grep (the #374 / fo-self-evidence-bar lesson).

## Relates to
- `lean-boot-hardening` (#374, archived) — the binary-output fix + the `detectBroadSearchAtBoot` guard this task hardens.
- `fo-self-evidence-bar` (z25…) — the broad-search prohibition absent from contract prose is the same gap fo-self-evidence-bar flags for the FO's own boot discipline.
