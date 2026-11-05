# Browser review

Review your agent’s work in your browser. Spacedock hosts the uploaded document; your agent stays in its existing environment.

## Open your review

In Claude Code, request a decision in the browser:

<!-- Keyword invocation verification: https://github.com/spacedock-dev/spacedock/issues/852 -->

<div class="sd-agent-prompt">
  <code>/spacedock:review plan.md web decision</code>
  <button type="button" class="sd-prompt-copy" aria-label="Copy agent prompt" aria-live="polite" hidden>copy</button>
</div>

In Codex, use `$spacedock:review` with the same arguments.

Or ask your agent:

> Send plan.md to me in Spacedock Review in the browser for a decision.

Agree to the upload. The first time, your agent runs `rr pair` for you to pair this computer with your account: sign in and approve in the browser page it opens. You pair once per computer.

Over SSH, the terminal shows a link and a code instead. Open the link in a browser on any device, sign in, enter the code and approve. Then answer `y` in the terminal.

Your agent then opens the review. Signed in, you are the owner and make the decision. Send the same share link to [invite others](sharing.md); they comment, and only you decide.

<a id="comment-and-decide"></a>

## Comment and decide

1. Select text and add a comment.
2. Press **Finish review**. Choose **Approve**, **Request changes** or **Hold decision**.
3. Press **Send decision**. Your agent receives the decision and included comments. The review closes.

Your unsent comments are included automatically. For comments already sent, select **Include** beside each one you want returned with your decision.

![Browser review with an anchored comment included and Request changes selected on the finish screen](assets/browser-decision.png)

To send comments while keeping the review open, choose **Send comments, leave open**. You can decide later.

<a id="next"></a>

**Related guides**

- [Sharing](sharing.md)
- [Comments and decisions](decisions.md)
- [Link expiry and waiting limits](reference.md#limits)
