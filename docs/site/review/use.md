# Use Spacedock Review

Ask your agent for a review, comment on the lines that need it, and send back your decision.

You need:

- Spacedock installed; see [Install](../index.md#install). After a Homebrew install, run `spacedock install --host claude` (or `--host codex`) once, or the command does not exist in your agent.
- One Markdown file to review.
- For the terminal: Zellij, tmux, Herdr, CMUX, Ghostty or Apple Terminal. In any other terminal, review on the web.

```text
/spacedock:review plan.md
```

That is the command in Claude Code. In Codex, ask in plain words, for example "review plan.md with Spacedock Review".

## Choose where to review

1. Run the command. The first time in a project, your agent asks "Terminal (TUI) or Web?". Terminal keeps the file on your machine. Web uploads it to review.spacedock.md, where anyone with the link can open it.
2. If you choose Web, agree to the upload. Your agent asks once more: "Web uploads the file to review.spacedock.md. Anyone with the link can open it, and it expires automatically after seven days. Agree?" If you answer No, nothing is recorded; run `/spacedock:review tui` to use the terminal instead.

Spacedock Review remembers your answer for this project. To change it, run `/spacedock:review` with no file, or `/spacedock:review tui` or `/spacedock:review web`.

## Review in the terminal

1. Read the file. It opens in a split beside your agent session in Zellij, tmux and Herdr, in your pane in CMUX, and in a new tab or window in Ghostty and Apple Terminal.
2. Select text and press <kbd>c</kbd> to comment, or <kbd>s</kbd> to suggest the text you want instead.
3. Press <kbd>q</kbd> to finish. Your comments go back to your agent, each attached to its text.

For a review you start with `/spacedock:review`, <kbd>q</kbd> sends your comments and finishes. To send a decision too, review on the web.

### If a workflow asks for a decision

A [Spacedock Workflow](../workflow/index.md) gate can open a terminal review that asks for your decision. There, <kbd>q</kbd> opens a finish screen first:

1. Press <kbd>Tab</kbd> to choose Approve, Request changes or Hold decision. The screen starts at "No decision". Request changes and Hold decision need an overall note or at least one comment.
2. Type an overall note if you want one. <kbd>Enter</kbd> saves it; <kbd>Esc</kbd> discards it.
3. Press <kbd>Enter</kbd> to send. Your agent gets the decision, your comments each attached to their text, and your overall note.

> **[Screenshot to take:** the finish screen with Request changes chosen and two comments.**]**

Press <kbd>?</kbd> in the review for every key.

## Review in the browser

1. Open your owner link. Your agent shows you a share link and opens your owner link in your browser. The owner link is for you alone, works once, and expires after 15 minutes.
2. Press "Sign in as owner". Signing in also links your agent to your account.
3. Select text and comment.
4. Send the share link to anyone whose comments you want. Their comments appear on your page.
5. In "All feedback", tick Include beside each already-sent comment you want your agent to receive, including colleagues' comments. Your own unsent comments are included automatically.
6. Press "Finish review…". Choose Approve to continue the task, Request changes to have the agent revise, or Hold decision to stop without changes.
7. Press "Send decision". The review closes.

Your agent picks up your decision and the comments you chose, and acts on it as the [decision table](reference.md#decisions) says. After Request changes, it shares the edited file under a new link; colleagues on the old link do not move to the new version, so send them the new link.

Your agent waits up to 110 minutes in Claude Code and 12 hours in Codex. If it stopped waiting, ask it to "check my decision on plan.md". If your owner link was used or expired, ask it to "publish again" for a new one.

> **[Screenshot to take:** the Finish review sheet with another reviewer's comment ticked and Request changes chosen.**]**

## What a colleague sees

Send the share link to your reviewers. They need nothing installed and no account.

1. They open the link. The page says "You're the reviewer" and offers a random name.
2. They continue as that name, or sign in to use their account's name.
3. They select text, comment, and press "Submit". The button counts their comments, then reads "Sent". "My feedback" lists what they sent; they do not see other reviewers' comments.
4. Their comments appear on your owner page.

A colleague can also hand the link to their own agent. It reads `review.spacedock.md/llms.txt` and comments the same way, under the name it declares.

## Next

- [Look up every argument, key and limit](reference.md).
- [Check what happens to a document you share](index.md#what-happens-to-a-document-you-share) before you send a link.
