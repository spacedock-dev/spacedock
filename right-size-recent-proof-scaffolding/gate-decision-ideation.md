# Ideation decision: right-size recent proof scaffolding

## Recommendation

APPROVE implementation. The design removes 560 net lines across 18 files while retaining one behavioral proof owner for every distinct failure mode.

## Selected approach

Remove duplicate Pi journeys, source-reading checks, cross-product coverage, and per-scenario installation scaffolding. Keep deterministic tests at the lowest layer that owns each behavior. Retain one Pi live journey for runtime behavior, and use a manual installation from the real release channel to prove published-package provenance.

Add two concise lines to the existing workflow guidance: require one primary proof owner, and require a distinct falsifying edit for every additional committed check. This adds no gate, lint, lane, fixture, table, or recurring step.

## Behavior proof before commit

Implementation will apply the README edit without committing it. Engram's behavior-diff will mint a base capsule and run three Codex trials before the edit and three after it. A neutral task will ask each trial to design proof for the Pi initial-dispatch bootstrap incident without revealing the expected one-owner result. The report must show that the rule changes test selection; identical before-and-after flows do not count as proof.

## Surface

Expected change: +65/-625 lines, a net reduction of 560 lines across 18 files. Tolerance: 100 net lines and two files. Relevant Pi journeys fall from two to one, and common live-suite package installations fall to zero. Product behavior, command grammar, stored formats, and authority remain unchanged.

## Decision

Approval starts implementation on a new stack layer above commit `37f50588aa37cfa571a88e4aa87b8f5c8f1b39e8`.
