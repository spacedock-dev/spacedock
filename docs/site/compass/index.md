# What is Spacedock Compass

You give an agent a goal and go do something else. Twenty turns later it is rewriting a file you never mentioned, or it has declared the work done and it isn't. With several sessions running, you find out when you read the result. **Compass keeps the goal you set beside what each session is doing, so the session that has left it stands out while it is still running.**

## How it works

Compass is a screen on your own machine that lists every agent session you have running, in Claude Code, Codex and the other agent tools you use. For each active session it shows where the work is, what the agent is doing now, what it does next, and whether it is blocked.

> **[Screenshot to take:** the Sessions screen with several active sessions, each showing its goal beside what the agent is doing now, and one marked as drifting.**]**

You type the goal for a session, and up to six lines describing the outcome you expect. When you ask, Compass [checks the session against them](use.md#check-a-session-for-drift) and says where the agent departed, citing the exact entries in the session record. You decide what to do about it: tell the agent, adjust the goal, or leave it.

Compass only reads. It never types into a session, and it never answers a prompt for you. Session content stays on your machine unless you [allow a check](use.md#check-a-session-for-drift), and the check asks first.

> **Draft note, not page copy:** as of Cargento v0.28.0, Compass checks for drift only when you press the button. Checking on its own is a setting that is off by default. Confirm before 10-09 whether that stays true for 10-15, then update [Use it](use.md#let-compass-check-while-you-are-away).

## With Workflow

If you run work through [Spacedock Workflow](../workflow/index.md), Compass also shows which step each work item is in and the decisions Workflow recorded for it. It can [alert you once](use.md#get-alerted-when-work-reaches-a-step) when an item reaches a step you pick.

## When to use it

- You run more than one agent session and cannot watch them all.
- You step away during long runs and want to know, on return, which sessions need you.
- Your review comments keep saying "this isn't what we agreed".

## Next

- [Install Compass](../index.md#install), then [set a goal and read the board](use.md).
- **Then Review.** When Compass shows a session has drifted and you want to mark up what it produced, [Spacedock Review](../review/index.md) puts the work beside your session and sends your comments back to the exact lines.
- [Compass reference](reference.md): every command and option.
