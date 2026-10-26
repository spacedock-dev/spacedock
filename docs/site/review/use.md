# Your first review

[Install Spacedock](../get-started/install.md). In Claude Code, use the review skill. Replace `plan.md` with your file name:

<!-- Keyword invocation verification: https://github.com/spacedock-dev/spacedock/issues/852 -->

<div class="sd-agent-prompt">
  <code>/spacedock:review plan.md decision</code>
  <button type="button" class="sd-prompt-copy" aria-label="Copy agent prompt" aria-live="polite" hidden>copy</button>
</div>

Or ask your agent:

> Send plan.md to me in Spacedock Review for a decision.

No active plan right now? Try this:

<div class="sd-agent-prompt">
  <code>Based on my recent agent history, draft a concise recommendation on whether we should adopt Spacedock Review. Use /spacedock:review to send it to me for a decision.</code>
  <button type="button" class="sd-prompt-copy" aria-label="Copy agent prompt" aria-live="polite" hidden>copy</button>
</div>

<a id="choose-where-to-review"></a>

## Open the review

Choose **Terminal** to review on your machine, or **Web** to review in your browser. Your agent remembers the choice for this project.

For Web, agree to the upload, then press **Sign in as owner** to pair with your agent.

## Add a comment

Select the text you want changed.

- **Terminal:** press <kbd>c</kbd> to add a comment.
- **Browser:** type your comment in the box that opens.

## Finish the review

- **Approve** to let your agent continue.
- **Request changes** to send the work back with your feedback.

=== "Terminal"

    Press <kbd>q</kbd> to open the finish screen. Press <kbd>Tab</kbd> to choose your decision, then <kbd>Enter</kbd> to send.

=== "Browser"

    Press **Finish review**, choose your decision, then **Send decision**.

<a id="read-the-response"></a>
<a id="review-the-revision"></a>

Your agent receives your comments and decision automatically.

<a id="review-in-the-terminal"></a>
<a id="review-in-the-browser"></a>
<a id="what-a-colleague-sees"></a>

**Related guides**

- [Comments and decisions](decisions.md): what happens after your decision.
- [Terminal controls and approval requests](terminal.md)
- [Browser review](browser.md)
- [Sharing](sharing.md)
