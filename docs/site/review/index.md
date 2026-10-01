# What is Spacedock Review

Your agent writes plans, specs and reports faster than you can read them, and each one waits for your call. **Review puts the work in front of you, and your comments go back to the agent attached to the exact lines they're about.**

## How it works

Your agent finishes a piece of work and asks you to review it. Review opens it beside your session.

![A plan open in Review beside the agent session](https://raw.githubusercontent.com/spacedock-dev/subspace/main/assets/review-one-file.gif)

Tables and diagrams render properly, so you read the work the way it was meant to look.

![A rendered diagram and table in Review](https://raw.githubusercontent.com/spacedock-dev/subspace/main/assets/mermaid-and-table.png)

You select text and [comment on it, or suggest a replacement](use.md#comment-on-the-work). When you finish, your comments go straight back to the agent that asked, each one attached to the text it's about. The agent picks up from those lines, not from your summary of them.

![A comment attached to selected text in Review](https://raw.githubusercontent.com/spacedock-dev/subspace/main/assets/anchored-feedback.png)

If you start the review [with questions turned on](use.md#ask-the-agent-while-you-read), you can also ask the agent about the file while you read, and see its answer in the same place.

> **[Screenshot to take:** a question asked mid-review, with the agent's answer shown in the Review pane.**]**

## In the terminal or the browser

Review works the same way in both, so you use whichever suits the moment. In the terminal, it opens beside your agent session in [the terminal you already use](../index.md#install): Zellij, tmux, CMUX, Herdr, Ghostty or Apple Terminal. It works in Claude Code and Codex.

> **[Screenshots to take, side by side:** the same review in the terminal and in the browser, with the same comment on the same line in both, so the reader sees it's one experience.**]**

> **Draft note, not page copy:** the browser view is part of the 10-15 release and depends on the relay reaching production (launch plan, "Where the repos actually are"). Confirm before 10-09.

## Ask for a second opinion

Some calls aren't yours alone. Your agent drafts a plan to migrate your database, and you own the product side, but the database details are outside what you know best. You [send the review to a teammate](use.md#ask-for-a-second-opinion) who works on the database every day. Their comments come back to the session you're working in, next to your own, so you see the work and the second opinion together, and the call stays yours.

You can ask for a second opinion from the terminal or the browser. The browser is the easier place for your colleague to read and comment.

> **[Screenshot to take:** your session with a colleague's comments beside your own, each marked with who wrote it.**]**

> **Draft note, not page copy:** sharing depends on the relay reaching production. Sharing from the terminal exists, but whether it ships in production for 10-15 is an open question. Confirm both before 10-09.

## When to use it

- Your agent has written something you have to approve: a plan, a spec, a report, a draft.
- The work is mostly right, and you want to fix the few lines that aren't without rejecting the rest.
- You want someone with the right expertise to look before you decide.

## Next

- [Install Review](../index.md#install), then [use it on its own](use.md).
- **Then Workflow.** When you find yourself reviewing the same kind of work again and again, at different steps along the way, there's a workflow under it. [Spacedock Workflow](../workflow/index.md) can [run that process](../get-started/survey.md), and bring you only the calls that need you. You can still [decide those calls in Review](../workflow/index.md#next).
- **Or Compass.** If your comments keep saying "this isn't what we agreed", you're catching drift after the work is done. [Spacedock Compass](../compass/index.md) catches it while the work is still running.
