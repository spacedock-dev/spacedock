---
title: Make advisory round recording structurally strict and semantically workflow-neutral
status: backlog
source: "Captain boundary audit, 2026-07-30: the provider-neutral gate recorder currently parses the development workflow's Material/fixed/declined vocabulary and LOC/AC feedback-line grammar."
started:
completed:
verdict:
score: "0.88"
worktree:
issue:
pr:
sprint: durable-decisions
id: zrcge8f5xbpc7dfcfn685mn4
---

Keep `gate record --round` strict about durable evidence while removing development policy from the storage and validation layer.

The recorder should continue to verify Briefing/log identity, digests, authority, ordering, complete one-time coverage of reviewer findings by worker responses, final advisory Resolution coverage, immutability, atomic publication, and the distinction between no findings and findings addressed. It should not interpret Material, correct-but-disproportionate, fixed/declined/mixed, why-not-material, promotes-when, value ACs, or a files/LOC/estimate/percentage/AC projection grammar.

## Problem

`internal/gates` and the gate-resolution specification currently turn one development workflow's Roborev taxonomy into canonical v1 storage semantics. That makes non-development review rounds invalid even when their retained evidence graph is complete and authoritative. It also causes the generic binary to decide policy that belongs to the active workflow stage.

## Proposed approach

Ideation should define the minimum structural review/response graph the recorder can validate without reading domain-specific Annotation bodies. Keep the existing recorder rather than adding another command or schema version. Treat the optional feedback-cycle projection as workflow-authored retained text with only the minimum identity/encoding checks needed for safe placement; derive generic summaries such as no-findings versus responded rather than all-fixed/all-declines/mixed. Prove the boundary with one development round and one refinement-style round using different vocabularies, plus negative cases for missing, duplicate, unauthorized, reordered, or tampered evidence.

## Out of scope

Do not define the development finding taxonomy or its worker behavior; the sibling workflow-owned disposition task owns that. Do not add compatibility wrappers, migrations, or v2 formats for this unreleased v1 surface.
