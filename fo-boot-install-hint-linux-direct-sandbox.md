---
id: z3j0tsbr6t3mqd39rhs8bbvq
title: "First-officer boot install journey: Linux-aware hint, direct install/upgrade offer, sandbox detection"
status: backlog
source: "GitHub issue spacedock-dev/spacedock#581 (nomen429, 2026-07-30): the install journey the first-officer skill hits at the version gate in first-officer-shared-core.md Startup step 1."
started:
completed:
verdict:
score:
worktree:
issue: spacedock-dev/spacedock#581
---

The first-officer boot hits the binary version gate (Startup step 1) and, on either abort class — binary absent, or binary present but wrong minor — prints a Mac-only Homebrew install hint and stops, leaving the human to copy-paste a command and restart the session. This task improves that install journey along the three axes the issue names: make the hint OS-aware (include the documented Linux `curl|sh` path, not just Homebrew), offer to run the install/upgrade directly and resume startup once the binary lands (turn hint-and-abort into one approved action), and detect sandboxed execution so a sandboxed install does not silently no-op (tell the human to run the install command themselves outside the sandbox, naming the exact command).

## Problem

{Ideation fills this in, seeded from the issue: today the FO prints a Mac-only hint and aborts; a Linux VM with no `spacedock` on PATH sees a Homebrew command it cannot run; a wrong-minor binary aborts with no self-repair; a sandboxed session would install the binary somewhere the host cannot see, a silent no-op from the human's perspective.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

{Ideation fills this in.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded from the issue; ideation refines, re-anchors each to the end value it measures, and fills the `Verified by` with a concrete falsifiable test.

**AC-1 - The binary-absent install hint is OS-aware, offering a working install path for the host OS.**
Verified by: {ideation names the test — e.g. a Go unit/fixture test that drives the hint builder for a Mac host (expects the Homebrew command) and a Linux host (expects the documented `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh` command), asserting the wrong-OS command is absent; name the concrete change (e.g. removing the Linux branch) that would make it fail.}

**AC-2 - On either abort class (binary absent, or binary present but wrong minor), the FO offers to run the install/upgrade itself and resumes startup once the binary lands, instead of hint-and-abort.**
Verified by: {ideation names the test — e.g. a behavior fixture that drives the boot abort path with a captive install command and asserts the FO runs it and re-checks `--version` to convergence, vs. the current print-and-exit; name the concrete change (reverting to print-and-exit) that would make it fail.}

**AC-3 - When the FO is running inside a sandbox, it detects the sandbox and tells the human plainly to run the install command themselves outside the sandbox, naming the exact command.**
Verified by: {ideation names the test — e.g. a Go unit test that drives the sandbox-detection path with a sandbox marker present and asserts the human-facing message names the exact install command and the "outside the sandbox" instruction, vs. attempting the install; name the concrete change that would make it fail.}

## Test plan

{Ideation fills this in, seeded from the issue: Go unit tests for OS-aware hint selection and sandbox detection; a behavior fixture for the offer-install-and-resume path (the riskiest mechanism — a self-modifying boot that re-invokes `--version` after install, which must converge and not loop); confirm `docs/site/get-started/install.md`'s Linux path is the canonical command the hint cites, so the hint and the docs do not drift.}
