# Terminal review

[Open your file in Review](use.md) and choose Terminal. You need Zellij, tmux, Herdr, CMUX, Ghostty or Apple Terminal.

## Read and comment

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

## Next

- [Comments and decisions](decisions.md).
- [Terminal keys](reference.md#terminal-keys).
