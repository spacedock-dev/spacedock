---
title: "Frontmatter contract"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-07-25 14:22:43"
---

# Frontmatter contract

Spacedock reads YAML frontmatter as the machine-readable state of a workflow and its entities. The first officer and dispatched workers write it; you operate through them, not by hand-editing. The field-level contract (names, types, patterns, defaults, invariants) is two machine-checkable schemas.

## Workflow README

The workflow `README.md` frontmatter declares the entity type, the id style, and the stages with their per-stage defaults and gates. The contract is [`workflow-readme.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/workflow-readme.mdschema.yml), which also specifies the required per-stage body subsections.

## Entity

Each entity's frontmatter carries its id, current stage, outcome, and worktree state. The contract is [`entity.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml), which defines the fields, the custom-field policy, the recognized body headings, and the invariants.

A gate decision is recorded in the entity's versioned `gates` collection before any stage transition or worker dispatch. The `spacedock gate record` verb is the sole writer for the canonical-v1 gate/attempt/Briefing/Resolution collection; callers provide a complete retained package manifest named `briefing.json`, a fully associated exact Result, or a chat decision, while the binary derives lifecycle and ids. Unknown or prototype fields inside binary-owned `gates` fail closed. The known `application` subtree carries the typed one-use lifecycle on the same writer surface. Canonical gates writes preserve top-level frontmatter outside `gates` and the entity body byte-for-byte. Each attempt stores one Briefing reference and digest: replaceable while Resolution is absent and frozen when present. `spacedock status --fields gate-state,gate-decision,gate-resolution` surfaces the current recorded result without changing the default status table.

`review-round` is the current pointer to one immutable advisory Review & Gate room. It repeats only the round, stage, cycle, and canonical Briefing binding; Resolutions remain in `briefing.review.jsonl`. Round recording shares the gate recorder's lock, full-entity compare-and-swap validation, and atomic entity replacement, but carries no gate selection, application, or status effect.

A recorded approval carries a one-use application. Status surfaces it as `approved-pending` until it is consumed, `consumed` once the approval has advanced exactly once through the ordinary transition path, `superseded` if the reviewed input changed before it was consumed, and `not-applicable` for a reviewer hold. A present dependency blocker or an active execution hold keeps the application ineligible; the application is never consumed twice.

## Sitemap

- [Command reference](../command-reference/index.md)
- [Glossary](../glossary/index.md)
- [Supported sandboxes](../sandbox/index.md)
