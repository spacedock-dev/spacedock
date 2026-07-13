---
title: Claude first officer must retain its launch cwd for boot
status: backlog
source: "Runtime Live E2E run 29214183726, Opus job 86712945963: the fixture was launched at /tmp/TestLiveDefaultHeadlessStopsAtGate.../002, but the FO later cd'd to /home/runner/work/spacedock/spacedock before status --boot."
score: 0.95
started:
completed:
verdict:
worktree:
issue:
id: z93q87sh2sg5hppx95frdkf8
---

The Opus live lane performed a genuine wrong-root boot even though the process cwd was the isolated workflow fixture. It read the launcher invariant, treated the directory containing `SPACEDOCK_BIN` as the project root, and explicitly ran `cd /home/runner/work/spacedock/spacedock && ${SPACEDOCK_BIN:-spacedock} status --boot --identify --json`. It then operated on the repository's `docs/dev` workflow while the fixture task remained at `draft`.

## Problem

The first-officer contract says to use `SPACEDOCK_BIN` consistently and to stay at the project root, but it does not state that the launch cwd is the workflow root or forbid deriving a workflow root from executable, plugin, PATH, or CI paths. A capable model can therefore conflate the binary location with the project it must operate. The live test correctly caught this; it is not a detector false positive, a `SPACEDOCK_REPO_ROOT` leak, or the Codex resume change.

## Proposed approach

Make the launch cwd authoritative for first-officer boot. The contract should require preserving that cwd and running the boot command without `cd`; `SPACEDOCK_BIN`, plugin locations, PATH, and CI paths locate programs only and must never select a workflow root. Keep the existing isolated-fixture live scenario as the behavioral proof. Do not add a fixture-root prompt, parser, or broad environment-scrubbing workaround for an env var absent from this incident.

## Out of scope

- The independent keep-moving grader false positive from earlier Opus attempts.
- Changes to Codex resume classification or frontdoor argv handling.
- Weakening wrong-root detection or accepting a lucky rerun as proof.

## Acceptance criteria

- **AC-1 (VALUE):** A headless Claude first officer launched in an isolated workflow fixture boots and drives that fixture to its review gate rather than the repository workflow, with the existing live test observing `draft -> review` plus the authored gate review.
  - Verified by: `TestLiveDefaultHeadlessStopsAtGate` on the affected Opus lane and the final-SHA Runtime Live E2E run.
- **AC-2:** Launcher/plugin/executable paths remain program locators only; the implementation introduces no root inference from them and no test-only fixture-root instruction.
  - Verified by: focused implementation review plus the same live fixture behavior.
- **AC-3:** The exact wrong-root failure remains detectable if it recurs; no detector weakening or suppression is used to green the lane.
  - Verified by: focused wrong-root detector controls and the failing-run command retained as the incident baseline.

## Test plan

Use the existing `TestLiveDefaultHeadlessStopsAtGate` as the primary runtime proof because this is model behavior, not static text. Run its focused deterministic support tests and the required Claude live lane after implementation; then rerun Runtime Live E2E on the exact release SHA. Avoid contractlint or a prose substring test as evidence.
