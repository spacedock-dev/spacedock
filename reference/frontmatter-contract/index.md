---
title: "Frontmatter contract"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-08-25 05:01:49"
---

# Frontmatter contract

Spacedock reads YAML frontmatter as the machine-readable state of a workflow and its entities. The first officer and dispatched workers write it; you operate through them, not by hand-editing. The field-level contract (names, types, patterns, defaults, invariants) is two machine-checkable schemas.

## Workflow README

The workflow `README.md` frontmatter declares the entity type, the id style, and the stages with their per-stage defaults and gates. The contract is [`workflow-readme.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/workflow-readme.mdschema.yml), which also specifies the required per-stage body subsections.

## Entity

Each entity's frontmatter carries its id, current stage, outcome, and worktree state. The contract is [`entity.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml), which defines the fields, the custom-field policy, the recognized body headings, and the invariants.

A gate decision is recorded in the entity's versioned `gates` collection before any stage transition or worker dispatch. `spacedock gate prepare` derives a provider-neutral room, canonical Briefing, request authority, and open binding from committed selected files. Chat or Subspace presents that committed Briefing and returns semantic decision and reason input; `spacedock gate record --decision` closes canonical-v1 gate/attempt/Briefing/Resolution state through the sole recorder. `spacedock gate withdraw` retires a stale request-backed open attempt without a Resolution or application. Open means neither withdrawal nor Resolution; withdrawn means withdrawal alone; closed means Resolution alone. Withdrawn and closed attempts are frozen. A withdrawal has fixed First Officer attribution, UTC time, nonblank reason, and retained request authority. No Briefing basename is canonical: prepared rooms freeze a clean relative locator, id, digest, request digest, and exact primary Artifact summary. Advisory-round Briefings remain summary-free. Unknown or prototype fields inside binary-owned `gates` fail closed except unknown keys inside an exact `records[*].attempts[*].application` mapping; those keys produce warnings only on explicit `status --validate` over active entities, are ignored for authority, and are never written. The known `application` subtree carries the typed one-use lifecycle on the same writer surface. Canonical gates writes preserve top-level frontmatter outside `gates` and the entity body byte-for-byte. Successful split-root close and consume writes commit and synchronize themselves. `spacedock status --fields gate-readiness` surfaces the current gate readiness without changing the default status table.

`review-round` is the current pointer to one immutable correction-review room. It repeats only the round, stage, cycle, and canonical Briefing binding; Resolutions remain in `briefing.review.jsonl`. Round recording shares the gate recorder's lock, full-entity compare-and-swap validation, and atomic entity replacement, but carries no gate selection, application, or status effect. The neutral producer retains canonical Briefing/log bytes and does not classify findings or write workflow body projections; a workflow-owned Cycle line, if present, remains untouched.

The known `application` subtree is approval-only. Canonical writes contain exactly `target-stage` and `state`; read-only compatibility may observe arbitrary unknown application keys as warnings. Unknown application extensions are ignored and are never authority. Revise and hold carry no application.

A recorded approval carries a one-use application. Status surfaces it as `approved-pending`, `consumed`, or `superseded`. Hold and revise carry no application; their Resolution remains the durable decision. The application is never consumed twice.

## Sitemap

- [Command reference](../command-reference/index.md)
- [Glossary](../glossary/index.md)
- [Supported sandboxes](../sandbox/index.md)
