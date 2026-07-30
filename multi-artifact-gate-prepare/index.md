---
title: Support multiple Artifacts in gate prepare
tatus: backlog
source: "q0 gate-room dogfood: captain requested a separately viewable git diff --stat Artifact, 2026-07-30."
score: 0.7
sprint: durable-decisions
id: vpf2143105f72t5m52dqmwyk
status: backlog
---

## Problem

`spacedock gate prepare` accepts exactly one `--artifact` and demotes every additional selected source to a Reference. A gate operator therefore cannot construct a canonical Briefing with multiple independently summarized Artifacts through the public command, even though Review v1 and Subspace can present them.

## Boundary

Extend mechanical room preparation without adding caller-authored JSON, copied payloads, provider-specific arguments, or compatibility behavior. Repeated Artifacts remain committed Git sources pinned by URI and SHA-256; References retain their existing meaning.

## Acceptance criteria

**AC-1 - An agent can prepare a gate room with two or more ordered Artifacts using the public CLI.**
Verified by: a focused CLI test whose canonical Briefing contains both ordered Artifact records with distinct summaries, Git-root locators, and byte digests.

**AC-2 - Multiple Artifacts add no source-payload duplication to the room.**
Verified by: a room-layout test that finds only the two authoritative metadata files immediately after preparation and resolves every Artifact from its pinned Git object.

**AC-3 - Existing one-Artifact and Reference behavior remains mechanically unambiguous.**
Verified by: focused positive and malformed-grammar CLI tests covering repeated Artifact arguments, per-Artifact summaries, ordering, duplicates, and References.

## Out of scope

Provider presentation changes, hand-authored Briefings, copied source bundles, or changes to gate decision and consumption semantics.

## Test plan

Ideation should choose the smallest unambiguous public grammar for pairing each Artifact with its summary, then extend the existing gate prepare CLI and room-preparation tests.
