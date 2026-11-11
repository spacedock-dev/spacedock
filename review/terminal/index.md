---
title: "In Terminal"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 08:30:05"
---

# Terminal review

Review beside your agent in tmux, Zellij, Herdr or cmux, or in a Ghostty or Apple Terminal window. The file stays on your machine.

## Open the review

In Claude Code, request a decision in the terminal:

`/spacedock:review plan.md terminal decision` copy

In Codex, use `$spacedock:review` with the same arguments.

`terminal` is optional. Omit it to use your project's saved preference. To save Terminal for this project, run `/spacedock:review terminal`.

Or ask your agent:

> Open plan.md in Spacedock Review in the terminal and ask for my decision.

## Comment or suggest an edit

Select text, then press `c` to comment or `s` to suggest replacement text.

## Finish the review

Press `q` to open the finish screen. Press `Tab` to choose **Approve**, **Request changes** or **Hold decision**, then `Enter` to send.

Request changes and Hold decision need a comment, a suggested edit or an overall note. Your agent receives your decision and comments.

## Questions

In Claude Code, open the terminal review with questions enabled:

`/spacedock:review plan.md terminal questions` copy

Or ask your agent:

> Open plan.md in Spacedock Review in the terminal with questions enabled.

Answers appear beside the document and use the agent's original conversation context.

Press `Q` to ask a question, such as "What's next?"

Press `T`, enter an interval in minutes, then `Enter` to refresh the answer automatically. Press `T` again to stop.

**Related guides**

- [Comments and decisions](../decisions/)
- [Terminal keys](../reference/#terminal-keys)
- [Workflow gates](../../concepts/gates-and-decisions/)

## Sitemap

- [Overview](../index.md)
- [First review](../use/index.md)
- [In Browser](../browser/index.md)
- [Sharing](../sharing/index.md)
- [Comments and decisions](../decisions/index.md)
- [Reference](../reference/index.md)
