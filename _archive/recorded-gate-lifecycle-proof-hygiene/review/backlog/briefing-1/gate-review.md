# Backlog gate: recorded-gate lifecycle proof hygiene

## Capability

Make the recorded-gate lifecycle's test evidence legible without changing its product
behavior, command surface, or accepted runtime outcomes.

## Evidence

- The exact-tip close-out audit proved that
  `TestRecordedGateLifecycleCommandTextMutants` reads shipped prose outside the
  quarantine and tests a parser whose needles are strict substrings of its mutations.
  Existing real-CLI, missing-event, refusal, resume, and live controls already own the
  lifecycle behavior.
- The Codex archival case already fails when the active entity disappears, but it fails
  at the initial file read. Its later named comparison is tautological. An explicit
  pre-read archive check improves attribution and can be shared with Claude.
- The captain accepted both dispositions and authorized the First Officer to continue
  the sprint close.

## Boundaries

- Delete the tautological prose test and orphaned helper; do not invent a replacement
  prose check.
- Improve Codex and Claude archive diagnostics without claiming new outcome coverage.
- Inspect the adjacent quarantine mapping once and include it only against direct
  evidence.
- Add no product behavior, prompt obligation, provider rule, compatibility layer,
  detector, or standing cap.

## Recommendation

Approve entry into ideation for the smallest current-main test-only cleanup.

## Decision

Approve to ideate this exact proof-hygiene boundary; revise to narrow it; or hold it
outside the pre-release close.
