# Terminal review

Review beside your agent in tmux, Zellij, Herdr or cmux, or in a Ghostty or Apple Terminal window. The file stays on your machine.

## Comment on the work

1. Select a passage and press <kbd>c</kbd> to comment, or <kbd>s</kbd> to suggest replacement text.
2. Press <kbd>q</kbd> to finish and return your feedback to the agent.

In Claude Code, a plain `/spacedock:review plan.md` asks for comments. If your agent requested a decision, <kbd>q</kbd> opens the finish screen instead.

<a id="read-and-comment"></a>
<a id="if-a-workflow-asks-for-a-decision"></a>

## Make a decision

To request an approval review, tell your agent:

<div class="sd-agent-prompt">
  <code>Open plan.md in Spacedock Review and ask for my decision before continuing.</code>
  <button type="button" class="sd-prompt-copy" aria-label="Copy agent prompt" aria-live="polite" hidden>copy</button>
</div>

On the finish screen, press <kbd>Tab</kbd> to choose **Approve**, **Request changes** or **Hold decision**, then <kbd>Enter</kbd> to send. Request changes and Hold decision need a comment or an overall note.

[Workflow gates](../concepts/gates-and-decisions.md) can also request a decision here.

## Ask about the document

Ask your agent to enable questions when it opens the review. You can ask about a passage and read the agent’s answer in the review, with its original conversation context.

Questions currently work in comments-only terminal reviews.

<a id="next"></a>

**Related guides**

- [Comments and decisions](decisions.md)
- [Terminal keys](reference.md#terminal-keys)
