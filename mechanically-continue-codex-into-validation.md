---
title: Mechanically continue Codex into validation
status: backlog
score: "0.90"
source: "v8 post-n28 pending-only validation receipt, 2026-08-10"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-product
id: s9hn38t0gwhzknnmr5w4m9d6
---
## Problem

After implementation completes, Codex can create a stamped validation dispatch envelope but never invoke the native validation worker. The n28 receipt stays pending, and the workflow can stop before validation runs.

## Value

After implementation completion, the supported Codex CLI host receives and acknowledges the validation dispatch envelope through n28 before any validation gate entry.

## Scope

- Use the exact v8 post-n28 pending-only receipt as baseline evidence.
- Name the missing supported-host primitive and the smallest bounded product surface.
- Prove validation pending → armed → consumed, a complete committed validation report, and later open gate.
- Do not implement product bytes during ideation.
- Do not add a transcript parser or Pi work.

## Acceptance criteria

- AC-1: Exact local Codex auto-continue reaches one validation pending → armed → consumed receipt chain.
- AC-2: The same validation worker produces a complete committed report before gate entry.
- AC-3: Missing native dispatch acknowledgment remains red and cannot enter the gate.
- AC-4: The design names exact files, gross/net estimate, and local proof plan before product work.
