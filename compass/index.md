---
title: "Overview"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 08:30:07"
---

# What is Spacedock Compass

You give an agent a goal and go do something else. Twenty turns later it is rewriting a file you never mentioned, or it has declared the work done and it isn't. With several sessions running, you find out when you read the result. **Compass keeps the goal you set beside what each session is doing, so the session that has left it stands out while it is still running.**

## How it works

Compass is a screen on your own machine that lists every agent session you have running, across Claude Code, Codex and eight other agent tools. For each active session it shows where the work is, what the agent is doing now, what it does next, and whether it is blocked.

*Each active session shows its goal beside what the agent is doing. Sessions with a recorded departure carry a Drift mark and sit above ordinary working sessions, below sessions waiting on you or holding a question.*

This is **Intent and drift**. You state your intent for a session as a goal. A `/goal` you set in Claude Code, Codex or Antigravity becomes the goal on its own. When you ask, Compass [checks the session against that intent](use/#check-a-session-for-drift) and says where the agent departed, citing the exact entries in the session record. When Claude Code or Codex says a goal is met, Compass displays the claim; Analyze drift checks it against what the session recorded. You decide what to do about a departure: tell the agent, adjust the intent, or leave it. Treat a drift level as an account of the evidence, not a measurement.

Compass reads your agents' session files and never writes to them. It never types into a session or chooses an answer for you; when a session asks through Compass, it returns the option you select. It keeps its own notes, your intent and a short history of what it observed, on your machine. With Compass on its default local address, session content leaves your machine only through model calls: [Analyze drift](use/#check-a-session-for-drift), session summaries, [Live monitor](use/#watch-a-session-live), and [checks Compass runs on its own](use/#let-compass-check-while-you-are-away). Calls run when you press a control or enable a feature. Opening a finished Claude Code session can start a missing summary, with summaries enabled and within its limits. The first of each kind tells you what it sends and where. **Pause model calls** stops all of them, and `--no-observer-model` refuses them for a whole run. Compass also reads your Claude Code and Cursor quota, which sends no session content; turn it off on the Sessions screen or with `--no-usage`.

## With Workflow

If you run work through [Spacedock Workflow](../workflow/), Compass also shows which step each work item is in and the decisions Workflow recorded for it. It can [alert you once](use/#get-alerted-when-work-reaches-a-step) when an item reaches a step you pick.

## When to use it

- You run more than one agent session and cannot watch them all.
- You step away during long runs and want to know, on return, which sessions need you.
- Your review comments keep saying "this isn't what we agreed".

## Next

- **Set your first intent.** [Install Spacedock](../#install), then [state an intent for a session and check it for drift](use/#set-the-goal).
- **Then Workflow.** If drift keeps coming from work that was never pinned down, [Spacedock Workflow](../workflow/) lets you write down once what good looks like at each step.
- [Compass reference](reference/): commands, options and keys.

## Sitemap

- [Use Compass](use/index.md)
- [Reference](reference/index.md)
