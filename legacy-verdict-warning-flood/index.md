---
title: "Legacy verdict tokens flood every state-checkout validation commit"
status: implementation
source: "Durable-decisions Commander dogfood, 2026-07-22: a normal `spacedock state commit` ran the state-checkout pre-commit validator and printed 117 pre-existing verdict-enum warnings before succeeding, burying any new warning attributable to the entity being committed. Dedupe found the schema validator and pre-commit-hook tasks, but no entity owns legacy-token migration or bounded warning output."
sprint:
id: mr9k7c0g35jhrrdv4zqyjctw
started: 2026-08-08T00:16:07Z
worktree: .worktrees/spacedock-ensign-legacy-verdict-warning-flood
---

One path-scoped state mutation currently emits 117 `Warning: field 'verdict' ... is not one of [PASSED REJECTED]` lines from the checkout-wide pre-commit validation: 104 lowercase `passed`, 9 lowercase `rejected`, and four other legacy/superseded values. The command exits successfully, but the historical backlog swamps the warning signal for the current change.

## Problem

The schema validator correctly treats conventional-field violations as warnings, and the pre-commit hook correctly runs checkout-wide validation. The development state checkout predates the canonical verdict enum, however, so every state commit repeats a large fixed corpus of legacy warnings. A newly introduced warning is difficult to distinguish in the flood, and routine command output becomes disproportionately noisy.

## Boundary for ideation

Determine the smallest safe ownership boundary between a one-time legacy data migration and bounded/delta-aware pre-commit diagnostics. Preserve schema-driven validation, warning severity, historical decision meaning, path-scoped state commits, and hard-error blocking. Do not silently suppress novel warnings or rewrite archived prose bodies; only frontmatter tokens whose semantics can be proven equivalent are candidates for migration.

## Acceptance sketch

- A clean state mutation produces bounded, actionable validation output instead of replaying the current 117-line legacy corpus.
- A newly introduced invalid verdict still appears and remains distinguishable; structural validation errors still block the commit.
- Any migrated `passed`/`rejected` tokens preserve their semantic value as canonical `PASSED`/`REJECTED`, with an auditable count and no unrelated entity-body changes.

## Evidence

`spacedock status --workflow-dir docs/dev --validate` exits 0 and currently emits 117 stderr lines, all verdict-enum warnings: 104 contain `value "passed"`, 9 contain `value "rejected"`, and the remainder are other legacy/superseded tokens. The installed pre-commit hook captures combined validation output and echoes it wholesale on every state-checkout commit.

## Out of scope

- Changing the durable-decisions sprint or its release criteria.
- Weakening warning-tier schema conformance or pre-commit hard-error enforcement.
- Bulk rewriting entity bodies, reports, review rooms, or unrelated frontmatter.
