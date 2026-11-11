---
title: "Command reference"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 05:18:33"
---

# Command reference

Run `spacedock <command> --help` for the full options. Run `spacedock` for the command list.

## Version

`spacedock --version` reports the version, platform and release channel. Inside an agent session, it also reports the runtime and sandbox.

## Launch

| Command            | Agent       |
| ------------------ | ----------- |
| `spacedock claude` | Claude Code |
| `spacedock codex`  | Codex       |
| `spacedock pi`     | Pi          |

Claude Code is the primary Workflow host. Codex and Pi support is experimental.

Launch from your project directory. You can include a task:

```
spacedock claude "continue the workflow"
```

Put agent-specific options after `--`. See [Supported sandboxes](../sandbox/) for sandbox options and [Local plugin development](../../advanced/plugin-development/) for local checkouts.

## Setup

| Command             | Purpose                                                                                            |
| ------------------- | -------------------------------------------------------------------------------------------------- |
| `spacedock install` | Install or refresh the agent integration                                                           |
| `spacedock doctor`  | Check installation and compatibility                                                               |
| `spacedock update`  | Update Spacedock, Review, Compass and their plugins to the newest release on this binary's channel |

`install` and `doctor` accept `--host claude`, `--host codex` or `--host pi`. The default is Claude Code. Follow the doctor output to resolve setup problems.

See [Installation](../../get-started/install/) for the manual setup path.

## Workflow

Your agent uses these commands to operate the workflow.

| Command                 | Purpose                                                 |
| ----------------------- | ------------------------------------------------------- |
| `spacedock status`      | Read or update work item state                          |
| `spacedock new`         | Create a work item                                      |
| `spacedock dispatch`    | Prepare worker assignments                              |
| `spacedock gate`        | Prepare reviews, record decisions and read gate history |
| `spacedock merge guard` | Complete delivery and record the outcome                |
| `spacedock state`       | Manage the workflow state checkout                      |
| `spacedock completion`  | Generate shell completions                              |

`spacedock gate read <ref-or-path> --json` returns the stored `gates` object, including every attempt and its application state. References and paths use the same resolution as `status --read`, including archived entities. The command changes no files. Records and attempts retain their stored order. Scalar values retain their stored values. An entity without gate history returns `{"records":[],"version":1}`. Invalid gate data exits nonzero without JSON output.

```
spacedock gate read gate-history-json-reader --json --workflow-dir docs/dev
spacedock gate read archive:finished-task --json --workflow-dir docs/dev
```

See [Workflow operation](../../running-workflows/operating/) for the conversation you use to direct the work.

## Sitemap

- [Frontmatter contract](../frontmatter-contract/index.md)
- [Glossary](../glossary/index.md)
- [Supported sandboxes](../sandbox/index.md)
