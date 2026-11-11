---
title: "In Browser"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 06:33:53"
---

# Browser review

Review your agent’s work in your browser. Spacedock hosts the uploaded document; your agent stays in its existing environment.

## Open your review

In Claude Code, request a decision in the browser:

`/spacedock:review plan.md web decision` copy

In Codex, use `$spacedock:review` with the same arguments.

Or ask your agent:

> Send plan.md to me in Spacedock Review in the browser for a decision.

Agree to the upload. The first time, your agent runs `subspace-relay-client pair` for you to pair this computer with your account: sign in and approve in the browser page it opens. You pair once per computer.

Over SSH, the terminal shows a link and a code instead. Open the link in a browser on any device, sign in, enter the code and approve. Then confirm in the terminal: answer `y`. If your agent started the pairing and reports `paired: not yet`, run `subspace-relay-client pair --confirm` in your own terminal on that computer and answer `y` there.

Your agent then opens the review. Signed in, you are the owner and make the decision. Send the same share link to [invite others](../sharing/); they comment, and only you decide.

## Comment and decide

1. Select text and add a comment.
1. Press **Finish review**. Choose **Approve**, **Request changes** or **Hold decision**.
1. Press **Send decision**. Your agent receives the decision and included comments. The review closes.

Comments you have added are included automatically; add or cancel any text still in the comment box first. For comments already sent, select **Include** beside each one you want returned with your decision.

To send comments while keeping the review open, choose **Send comments, leave open**. You can decide later.

**Related guides**

- [Sharing](../sharing/)
- [Comments and decisions](../decisions/)
- [Link expiry and waiting limits](../reference/#limits)

## Sitemap

- [Overview](../index.md)
- [First review](../use/index.md)
- [In Terminal](../terminal/index.md)
- [Sharing](../sharing/index.md)
- [Comments and decisions](../decisions/index.md)
- [Reference](../reference/index.md)
