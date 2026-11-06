---
title: "Frontmatter contract"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-06 02:11:11"
---

# Frontmatter contract

YAML frontmatter stores workflow settings and work item state. Your agent maintains these fields.

| Record          | Contents                             | Schema                                                                                                           |
| --------------- | ------------------------------------ | ---------------------------------------------------------------------------------------------------------------- |
| Workflow README | Stages, defaults and decision points | [Workflow schema](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/workflow-readme.mdschema.yml) |
| Work item       | Identity, current stage and outcome  | [Entity schema](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml)            |

The schemas define field names, types and required values. [Workflows and entities](../../concepts/workflows-and-entities/) explains the files you see.

## Sitemap

- [Command reference](../command-reference/index.md)
- [Glossary](../glossary/index.md)
- [Supported sandboxes](../sandbox/index.md)
