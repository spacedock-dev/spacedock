---
id: "001"
title: Add login rate limiting
status: review
source: fixture
started: 2026-06-08T00:00:00Z
---

# Add login rate limiting

Throttle repeated failed login attempts per account to blunt credential-stuffing.

## Proposed approach

Add a sliding-window counter keyed by account id; reject the sixth failed attempt
within five minutes with a 429 and a retry-after hint.

## Acceptance criteria

- **AC-1 — the sixth failed login within the window is rejected.** Five failures
  pass through; the sixth returns 429 with a retry-after header. Verified by an
  integration test driving six attempts.
- **AC-2 — a successful login resets the window.** A success between failures clears
  the counter so legitimate users are not locked out. Verified by a test interleaving
  a success.

## Stage Report: review

- DONE: Sliding-window counter rejects the sixth failed attempt within five minutes
  Integration test drives six attempts; the sixth returns 429 with retry-after.
- DONE: Successful login resets the window
  Test interleaves a success between failures; counter clears, no lockout.

### Summary

Implemented the per-account sliding-window throttle and proved both acceptance
criteria with integration tests. The reviewer found no material issues and
recommends approval.
