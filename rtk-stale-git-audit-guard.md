---
id: m1whe67x7nnh1gjwrnsksgak
title: Detached audits / validation must catch rtk-stale-git via an un-proxied SHA verify
status: backlog
source: "Commander (2026-06-08) qa validation: the validation ensign caught an rtk-stale-git discrepancy on its own — the worktree contract.go blob ≠ the expected e050996d blob — and only a raw, un-proxied blob-SHA compare exposed it, forcing a re-pin to 32ceb73e before the adversarial audit re-ran (verdict held PASSED). FO boot (this session) independently hit rtk mangling/blocking ls/find/ps. A SHA-pinned audit that trusts proxied git is a silent hole in exactly the discipline meant to catch silent holes."
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0199-pre-flip-mechanics
group: dev-quality
sprint-readiness: ready
---

A SHA-pinned detached adversarial audit (or validation) must not be fooled by a git proxy that serves stale state. The defense the qa validator improvised — verify the resolved commit/blob SHA via an un-proxied git read before trusting a pin — should be a standard, checkable guard, not a one-off save.

## Problem

A token-saving git proxy (rtk) can return stale/cached git output. The detached-adversarial-audit discipline pins to a specific merge-result SHA and asserts that the deliverable's tests would catch a broken edit. If the audit resolves that SHA (or reads blobs) through the proxy and the proxy is stale, the audit pins to the WRONG tree and validates something other than the merge result — and reports green. That is a silent correctness hole inside the discipline whose entire job is to catch silent holes. It already happened once (qa, contract.go — the contract gate) and was caught only because the validator fell back to a raw blob-SHA compare.

## Proposed approach

Ideation fills this in and decides the deliverable form (the proof-policy bars a pure-prose "the discipline says to verify" as the AC — it must be a checkable guard):

- **Likely:** a helper/guard the audit + validation flow calls to resolve the checked-out commit/blob SHA via an explicitly un-proxied git invocation, and fail LOUDLY on a proxied-vs-un-proxied mismatch. The discipline text then points at that guard.
- The guard is testable: simulate a stale-proxy read (or a deliberately divergent blob) and assert the guard goes red / re-pins, where a trusting audit would have stayed green.

## Out of scope

- rtk itself (the operator's global tooling) — this is the spacedock-side defense, not an rtk fix.
- General command-output mangling by rtk (`ls`/`find`/`ps`) — annoyance, not a correctness hole; separate if ever worth it.

## Acceptance criteria

Ideation/implementation fills in. Sketch:

- A SHA-pinned audit/validation detects an rtk-stale-git divergence and refuses to trust the stale pin (verified by a test that feeds a stale/divergent blob and asserts the guard goes red / re-pins — the expected value comes from the real un-proxied git object, an independent source that can diverge from the proxy, so it is not a tautology).

## Test plan

Ideation/implementation fills in. The riskiest unknown is reproducing rtk staleness deterministically — exercise it first: construct a case where proxied git and un-proxied git disagree on a blob SHA, and confirm the guard catches it.
