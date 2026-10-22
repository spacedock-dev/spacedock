# Your first review

Ask your agent for a review, comment on the lines that need it, and send back your decision.

You need:

- [Spacedock installed and connected to your agent](../get-started/install.md).
- One Markdown file to review.

```text
/spacedock:review plan.md
```

That is the command in Claude Code. In Codex, ask in plain words, for example "review plan.md with Spacedock Review".

## Choose where to review

1. Run the command. The first time in a project, your agent asks "Terminal (TUI) or Web?". Terminal keeps the file on your machine. Web uploads it to review.spacedock.md, where anyone with the link can open it.
2. If you choose Web, agree to the upload. Your agent asks once more: "Web uploads the file to review.spacedock.md. Anyone with the link can open it, and it expires automatically after seven days. Agree?" If you answer No, nothing is recorded; run `/spacedock:review tui` to use the terminal instead.

Spacedock Review remembers your answer for this project. To change it, run `/spacedock:review` with no file, or `/spacedock:review tui` or `/spacedock:review web`.

<a id="review-in-the-terminal"></a>

## Terminal

[Read and comment beside your agent](terminal.md).

<a id="review-in-the-browser"></a>

## Browser

[Review and decide in your browser](browser.md).

<a id="what-a-colleague-sees"></a>

## Sharing

[Shared reviews](sharing.md).
