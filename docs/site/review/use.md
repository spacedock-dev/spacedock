# Use Spacedock Review

Ask your agent for a review, comment on the lines that need it, and send back your decision.

```text
/spacedock:review plan.md
```

That is the command in Claude Code. In Codex, ask your agent to review the file with Spacedock Review. You need Spacedock installed; see [Install](../index.md#install).

## Choose where to review

1. **Run the command.** The first time in a project, your agent asks "Terminal (TUI) or Web?". Terminal keeps the file on your machine. Web uploads it to review.spacedock.md, where anyone with the link can open it.
2. **If you choose Web, agree to the upload.** Your agent asks once more: "Web uploads the file to review.spacedock.md. Anyone with the link can open it, and it expires automatically after seven days. Agree?"

Your answer is remembered for this project. To change it, run `/spacedock:review` with no file.

## Review in the terminal

1. **Read the file.** It opens beside your agent session.
2. **Select text and press <kbd>c</kbd> to comment**, or <kbd>s</kbd> to suggest the text you want instead.
3. **Press <kbd>q</kbd> to finish.** The finish screen shows "No decision".
4. **Press <kbd>Tab</kbd> to choose** Approve, Request changes or Hold decision. Type to add an overall note. Request changes and Hold decision need a note or at least one comment.
5. **Press <kbd>Enter</kbd> to send.** Your agent gets the decision and every comment, each attached to the text it is about.

> **[Screenshot to take:** the finish screen with Request changes chosen and two comments.**]**

Press <kbd>?</kbd> in the review for every key.

## Review in the browser

1. **Open your owner link.** Your agent shows you a share link and opens your owner link in your browser. The owner link is for you alone, works once, and expires after 15 minutes.
2. **Press "Sign in as owner".** The page says "Your agent gets every comment and your decision." Signing in also links your agent to your account.
3. **Select text and comment.**
4. **Send the share link** to anyone whose comments you want. Their comments appear on your page.
5. **Press "Finish review…" and choose** Approve, Request changes or Hold decision. Tick the comments, yours or others', that go with it.
6. **Send.** The review closes, and your agent, which has been waiting, picks up your decision and the comments you chose.

> **[Screenshot to take:** the Finish review sheet with another reviewer's comment ticked and Request changes chosen.**]**

## What a colleague sees

Send the share link to the person who knows the area best. They need nothing installed and no account.

1. **They open the link.** The page says "You're the reviewer", with a random name filled in.
2. **They continue with that name, type their own, or sign in** to use their account's name.
3. **They select text and comment.** "My feedback" lists what they sent; they do not see other reviewers' comments.
4. **Their comments reach your page**, and you choose which go to your agent with your decision.

A colleague can also hand the link to their own agent. It reads `review.spacedock.md/llms.txt` and comments the same way, under the name it declares.

## Next

- [Look up every argument, key and limit](reference.md).
- [Check what happens to a document you share](index.md#what-happens-to-a-document-you-share) before you send a link.
