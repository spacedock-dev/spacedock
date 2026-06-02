---
id: 3g6gbbn1bvk41a57tjbe50rv
title: Upload the live-e2e transcript as a CI artifact (diagnose failures past gh log truncation)
status: backlog
source: "FO (2026-06-02): 38's streamWatcher tees the FO transcript to t.Log, but on a failed run gh truncates the large job log AND the CI uploads only the spacedock binary, not the transcript — so the sonnet FO-stall on 38's PR was not fully diagnosable. Prerequisite for root-causing the headless FO-drive flake."
started:
completed:
verdict:
score: "0.30"
worktree:
issue:
---

`38` (live-e2e-per-stage-timeouts) made the live test STREAM the FO's stream-json to `t.Log` so a hang names the stalled step. But that diagnosability is defeated downstream: on 38's own PR CI-E2E (#261) the sonnet job failed (FO stalled headless), and the streamed transcript was NOT recoverable — `gh run view --log` truncates the multi-MB streamed log, AND the CI "Upload live artifacts" step (`.github/workflows/runtime-live-e2e.yml`) uploads only the `spacedock` binary, not the test transcript. So the captured-transcript half of 38's diagnosability never reaches a human on a real failure.

This is the PREREQUISITE for root-causing the headless FO-drive flake (the sibling follow-up `headless-fo-drive-flake`): you cannot diagnose WHY the FO stalls without the failing run's transcript.

## Scope
- Capture the live test's `-v` output (the streamed stream-json transcript) to a file in the CI step (e.g. tee `go test -tags live … -v` to a path under the artifact dir), and include that file in the "Upload live artifacts" upload (alongside the binary). The transcript MUST survive a FAILED/killed run (tee as it streams, not capture-on-success).
- Small `.github/workflows/runtime-live-e2e.yml` change; possibly a tiny test-side hook if a known transcript path is cleaner than a tee. NOT the status lane, NOT the FO runtime.

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — A failed live-e2e run uploads its transcript.** On a CI-E2E job that FAILS (e.g. a watcher stepTimeout), the uploaded artifact contains the streamed FO transcript (the stream-json lines + the labelled stepTimeout/stepFailure), not just the binary.
Verified by: inspect a failing run's artifact (or a forced-fail) and confirm the transcript file is present + non-empty + carries the failure tail.

**AC-2 — Success path + exit code unaffected.** A passing run still uploads cleanly; the tee/redirect does NOT swallow the test exit code (the job still goes red on failure — pipefail or equivalent).
Verified by: a passing run green, a failing run red; the workflow exit-code semantics preserved.

## Notes
- `.github/workflows/runtime-live-e2e.yml`. Merge this FIRST (captain) — prerequisite for the `headless-fo-drive-flake` investigation.
