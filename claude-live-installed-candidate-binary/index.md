---
title: Install the checked-out candidate into a neutral bin for Claude live E2E
status: ideation
source: "Captain direction 2026-07-13 after Opus treated the checkout-shaped SPACEDOCK_BIN path as a workflow root."
score: 0.85
started: 2026-07-13T04:20:06Z
completed:
verdict:
worktree:
issue:
id: 2h786rdtw1f26q7q98v5bp43
---

The Claude live job currently builds `./spacedock`, exports `SPACEDOCK_BIN=$(pwd)/spacedock`, and adds the checkout to PATH. The front door then deliberately re-exports its own executable path to the first officer. In a real Opus failure, that repository-shaped path became the apparent project root even though the host process started in a fixture.

## Problem

The live candidate must remain the exact checked-out release SHA, but its binary path should not visually identify the repository as a workflow root. Public package installation, Homebrew, and the release install script cannot prove the candidate commit. A symlink is also invalid because the launcher resolves it back to the checkout.

## Proposed approach

For the Claude live job only, install the local checkout with `GOBIN="$RUNNER_TEMP/spacedock-live-bin" go install ./cmd/spacedock`, then export that physical binary through `SPACEDOCK_BIN` and PATH. Retain the staged local `--plugin-dir` for current skill bytes. Record enough runtime evidence to show the binary came from the run's checked-out SHA. This is harness isolation, not a substitute for the separate launch-cwd-authority task.

## Out of scope

- Public/brew/release-script installation.
- Symlink-based relocation.
- Removing `SPACEDOCK_BIN` propagation or changing the launcher invariant.
- Treating this as a semantic guarantee that a model will never choose a wrong cwd.

## Acceptance criteria

- **AC-1 (VALUE):** Claude live E2E exercises a physical candidate binary outside the checkout while continuing to run the exact workflow-dispatch SHA and current local plugin bytes.
  - Verified by: a final-SHA Claude live job with recorded binary provenance and successful durable live scenarios.
- **AC-2:** The host-visible binary path no longer names the repository checkout, and no symlink resolves it back there.
  - Verified by: live-job environment/provenance evidence rather than instruction-text matching.
- **AC-3:** The change does not weaken the separate wrong-root detector or claim to solve launch-cwd authority.
  - Verified by: focused regression controls plus review of the bounded diff.

## Test plan

Use a focused workflow/config check only if it asserts an independently observable build/provenance value; do not add contractlint or a prose-grep test. Run the Claude live lane on the exact candidate SHA and inspect its recorded binary provenance, fixture durable state, and wrong-root controls. Retain the existing full Runtime Live E2E gate for release certification.
