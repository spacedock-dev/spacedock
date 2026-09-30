# Use Review on its own

Review your agent's work, fix the lines that need fixing, and send your comments back, all without leaving your session.

You need Review installed. See [Install Review](../start.md).

## Comment on the work

Review works the same way in the terminal and the browser: you comment on the exact text, and your comments go back to the agent.

=== "In the terminal"

    1. **Ask your agent to open the file in Review**, or open it yourself:

        ```text
        /r plan.md
        ```

        In Codex, type `$r plan.md`. Review opens beside your session.

    2. **Select the text you want to talk about.** Press <kbd>c</kbd> to comment on it, or <kbd>s</kbd> to suggest a replacement.

        ![A comment attached to selected text in Review](https://raw.githubusercontent.com/spacedock-dev/subspace/main/assets/anchored-feedback.png)

    3. **Comment on as many lines as you need.** Leave the rest alone: lines you don't comment on stand as they are.

    4. **Press <kbd>q</kbd> to finish.** Your comments and suggested edits go back to the agent, each one attached to the text it's about.

    The agent's next turn starts from your comments. If its revision needs another round, open the file again and comment on what's left.

    > **[Screenshot to take:** a suggested replacement, showing the original text and the proposed text.**]**

    Press <kbd>?</kbd> in Review for every key.

=== "In the browser"

    > **Draft note, not page copy:** browser steps to come, once the browser view ships (depends on the relay reaching production). Question for Kent Chen below.

    > **[Screenshot to take:** a comment on one line in the browser view.**]**

## Ask the agent while you read

Sometimes you need to know why the agent wrote something before you can judge it. In the terminal, start the review with questions turned on:

```text
/r --allow-question plan.md
```

1. **Press <kbd>Q</kbd>** and type your question about the file.
2. **Keep reviewing.** The agent answers from the same file you're reading, and the answer appears in Review.

> **[Screenshot to take:** a question and the agent's answer, shown in the Review pane.**]**

> **Draft note, not page copy:** asking the agent works in the terminal only today. The browser viewer (subspace-web v0.2.6) handles comments, with nothing for questions. Question for Kent Chen below.

## Ask for a second opinion

Some calls need someone who knows the area better than you do. Send the review to a colleague, and their comments come back to your session next to your own. The call stays yours.

> **Draft note, not page copy:** steps to come. Sharing depends on the relay reaching production, and the terminal side is an open PR, so the real steps aren't known yet. **Questions for Kent Chen:** on 10-15, (1) how does a user open a review in the browser, and how do they comment and finish there? (2) what exact steps does a user take to share a review, from the terminal and from the browser, and what does the colleague see and do? (3) will asking the agent a question work in the browser on 10-15, or is it terminal-only at launch?

> **[Screenshot to take:** your session with a colleague's comments beside your own, each marked with who wrote it.**]**

## Next

- [What connects today](../start.md): decide a call [Workflow](../workflow/index.md) holds for you in Review, instead of in chat.
- [Review reference](reference.md): every command and key.
