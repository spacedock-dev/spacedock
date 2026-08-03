---
title: "Frontmatter contract"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-03 03:53:10"
---

# Frontmatter contract

Spacedock reads YAML frontmatter as the machine-readable state of a workflow and its entities. The first officer and dispatched workers write it; you operate through them, not by hand-editing. The field-level contract (names, types, patterns, defaults, invariants) is two machine-checkable schemas.

## Workflow README

The workflow `README.md` frontmatter declares the entity type, the id style, and the stages with their per-stage defaults and gates. The contract is [`workflow-readme.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/workflow-readme.mdschema.yml), which also specifies the required per-stage body subsections.

## Entity

Each entity's frontmatter carries its id, current stage, outcome, and worktree state. The contract is [`entity.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml), which defines the fields, the custom-field policy, the recognized body headings, and the invariants.

A gate decision is recorded in the entity's versioned `gates` collection before any stage transition or worker dispatch. `spacedock gate prepare` derives a provider-neutral room, canonical Briefing, request authority, and open binding from committed selected files; `spacedock gate record` closes canonical-v1 gate/attempt/Briefing/Resolution state from a chat decision or the prepared room's room-backed Result. `spacedock gate withdraw` retires a stale request-backed open attempt without a Resolution or application. Open means neither withdrawal nor Resolution; withdrawn means withdrawal alone; closed means Resolution alone. Withdrawn and closed attempts are frozen. A withdrawal has fixed First Officer attribution, UTC time, nonblank reason, and retained request authority. No Briefing basename is canonical: prepared rooms freeze a clean relative locator, id, digest, request digest, and exact primary Artifact summary. Advisory-round Briefings remain summary-free. Unknown or prototype fields inside binary-owned `gates` fail closed. The known `application` subtree carries the typed one-use lifecycle on the same writer surface. Canonical gates writes preserve top-level frontmatter outside `gates` and the entity body byte-for-byte. `spacedock status --fields gate-state,gate-decision,gate-resolution` surfaces the current state without changing the default status table.

`review-round` is the current pointer to one immutable advisory Review & Gate room. It repeats only the round, stage, cycle, and canonical Briefing binding; Resolutions remain in `briefing.review.jsonl`. Round recording shares the gate recorder's lock, full-entity compare-and-swap validation, and atomic entity replacement, but carries no gate selection, application, or status effect.

A recorded approval carries a one-use application. Status surfaces it as `approved-pending` until it is consumed, `consumed` once the approval has advanced exactly once through the ordinary transition path, `superseded` if the reviewed input changed before it was consumed, and `not-applicable` for a reviewer hold. A present dependency blocker or an active execution hold keeps the application ineligible; the application is never consumed twice.

## Sitemap

- [Command reference](../command-reference/index.md)
- [Glossary](../glossary/index.md)
- [Supported sandboxes](../sandbox/index.md)
