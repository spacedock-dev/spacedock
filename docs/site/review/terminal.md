# Terminal review

Review beside your agent in tmux, Zellij, Herdr or cmux, or in a Ghostty or Apple Terminal window. The file stays on your machine.

## Open the review

In Claude Code, request a decision in the terminal:

<!-- Keyword invocation verification: https://github.com/spacedock-dev/spacedock/issues/852 -->

<div class="sd-agent-prompt">
  <code>/spacedock:review plan.md terminal decision</code>
  <button type="button" class="sd-prompt-copy" aria-label="Copy agent prompt" aria-live="polite" hidden>copy</button>
</div>

`terminal` is optional. Omit it to use your project's saved preference. To save Terminal for this project, run `/spacedock:review terminal`.

Or ask your agent:

> Open plan.md in Spacedock Review in the terminal and ask for my decision.

<a id="read-and-comment"></a>
<a id="comment-on-the-work"></a>

## Comment or suggest an edit

Select text, then press <kbd>c</kbd> to comment or <kbd>s</kbd> to suggest replacement text.

<a id="make-a-decision"></a>
<a id="if-a-workflow-asks-for-a-decision"></a>

## Finish the review

Press <kbd>q</kbd> to open the finish screen. Press <kbd>Tab</kbd> to choose **Approve**, **Request changes** or **Hold decision**, then <kbd>Enter</kbd> to send.

Request changes and Hold decision need a comment or an overall note. Your agent receives your decision and comments.

<a id="ask-about-the-document"></a>

## Questions

In Claude Code, open the terminal review with questions enabled:

<div class="sd-agent-prompt">
  <code>/spacedock:review plan.md terminal questions</code>
  <button type="button" class="sd-prompt-copy" aria-label="Copy agent prompt" aria-live="polite" hidden>copy</button>
</div>

Or ask your agent:

> Open plan.md in Spacedock Review in the terminal with questions enabled.

Answers appear beside the document and use the agent's original conversation context.

<a id="next"></a>

**Related guides**

- [Comments and decisions](decisions.md)
- [Terminal keys](reference.md#terminal-keys)
- [Workflow gates](../concepts/gates-and-decisions.md)
