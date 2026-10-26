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

Agree to the upload. Your agent opens your owner link. Press **Sign in as owner** so you can make the decision.

Keep the owner link for yourself. The separate share link is for [inviting others](sharing.md).

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
