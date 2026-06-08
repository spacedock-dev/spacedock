---
id: "001"
title: Measure prompt-caching latency impact
status: review
source: fixture
started: 2026-06-08T00:00:00Z
---

# Measure prompt-caching latency impact

Quantify the end-to-end latency change from enabling prompt caching on the hot path.

## Proposed approach

Run a paired A/B over 500 representative requests with caching off then on; report
the median and p95 latency delta with a confidence interval.

## Acceptance criteria

- **AC-1 — the median latency delta is reported with a confidence interval.** The
  run yields a median delta and a 95% CI, not a point estimate. Verified by the
  analysis output containing both.
- **AC-2 — the A/B arms are paired on the same request set.** Both arms replay the
  identical 500 requests so the delta is not confounded by traffic mix. Verified by
  the harness logging matching request ids per arm.

## Stage Report: review

- DONE: Median latency delta reported with a 95% confidence interval
  Analysis output shows a -38ms median delta, CI [-44ms, -31ms] over 500 paired requests.
- DONE: A/B arms paired on the same request set
  Harness log confirms identical request ids replayed across both arms.

### Summary

Ran the paired A/B over 500 requests and measured a 38ms median latency improvement
from prompt caching, with a tight confidence interval. Both acceptance criteria are
satisfied; the reviewer recommends approval.
