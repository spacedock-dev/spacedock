# Use Compass

State an intent for a session, read the board, and act when an agent drifts from it.

You need Spacedock installed with Compass enabled. See [Install Spacedock and enable Compass](../index.md#install).

Start Compass from your agent session:

```text
/compass:compass
```

> **[command may change with one install]**

In Codex, type `$compass`. You can also ask your agent to "open compass". Compass starts on your machine and gives you its address, `http://127.0.0.1:4553/` by default. It does not open a browser for you, so open the address yourself.

> **[Screenshot to take:** the Sessions screen right after Compass opens.**]**

New to Intent and drift? Open **See a sample session (view only)** above the Sessions list. It is an anonymized session that drifted, with its reading and the correction Compass composed. Nothing in it is saved, sent or run.

## Set the goal

1. **Open the session.** On the Sessions screen, click the session's row. Its page opens with an **Intent** section at the top.

2. **Type the goal.** One line, up to 240 characters. If you typed nothing, Compass drafts one from your first prompt in a Claude Code or Codex session, tinted and marked "from your prompt". **Use your prompt** lists up to five of your prompts to choose from. A draft is not saved until you press Save intent or press Analyze drift, which adopts it.

3. **Add the expected outcome.** Up to six lines of up to 240 characters, one check per line. Press **Add a line** for each. This step is optional, but a drift level needs at least one line. Until then the level reads "Not enough recorded yet".

4. **Press Save intent.** Compass keeps the words you typed and shows them beside the session from now on. Nothing is sent to the agent.

> **[Screenshot to take:** the Intent section with a typed goal and three outcome lines, saved.**]**

## Read the board

Sessions puts active work first. Each active session answers four questions: where it is, what it is doing now, what it does next, and whether it is blocked. Your goal sits beside what the agent is doing, so a mismatch is visible without opening anything.

| Open | To see |
| --- | --- |
| **Sessions**, key <kbd>s</kbd> | Every session, active work first |
| **Projects**, key <kbd>p</kbd> | Sessions grouped by project. Each project opens to Now, Course, Decisions and Console |
| **Attention**, key <kbd>a</kbd> | What needs you: sessions waiting on a person, and sessions at risk |
| **Intent log** | Every session that has an intent, with what Compass raised about it. A link in the navigation, with no key |

Press <kbd>Escape</kbd> in a session to go back.

## Check a session for drift

1. **Press Analyze drift** on the session's page.

2. **Read what it sends, then press Allow and analyze.** The first press shows what Compass will send, and to which agent tool. A Claude Code session is read by Claude Code, or by Codex when `claude` is not installed. Compass remembers your answer. **Turn off readings** on any session page takes it back.

3. **Read the result.** It says one of: "Departs from your intent", "Can't tell", "Not reached at this stop", or "Nothing found against what it read". A departure lists each line of your intent it contradicts, with the numbered session entries as evidence. **What the analysis read** shows the model, your intent, the part of the session it read, and what it did not read.

> **[Screenshot to take:** a "Departs from your intent" result, with the cited entries numbered beside it.**]**

A session with a departure on record gets a **Drift** mark on the Sessions screen and moves up the list, below any session blocked on you.

For a Claude Code session, the analysis also sends the checks the session ran and the files it wrote, but only after you allow tool output for the destination the page names.

Analyze drift reads the work of Claude Code and Pi sessions. On the other tools, the Drift section says "Compass can't read work or the agent's replies from this harness".

Each analysis uses model capacity from your Claude Code or Codex account. Compass allows 100 in a rolling twenty-four hours, shared across sessions and tools, and failed attempts count. When the limit is spent, a banner at the top of the session page says when you can analyze again.

For a Claude Code session you can also turn on **Live monitor** in **More** on the session page. It is off by default. It re-estimates drift after every turn from the checks the session ran and the files it wrote, with no model call, and shows the result as None or low, Medium, High or Extreme. It never notifies you. The [levels](reference.md#drift-levels) are in the reference.

## Act on a departure

Compass never writes into your session. You choose what happens next.

1. **Press Steer back.** Compass composes a correction from your goal and the evidence, ready to copy. Steer back is available for Claude Code sessions.

2. **Paste it into the session.** Attention and the project view offer a command to copy that re-enters the session.

If you gave the agent a later direction, Compass asks about it before you analyze again. Press **Add it to my intent** to make it a line of your intent, or **Keep my intent and analyze** to leave your intent as it was.

If a reading is wrong, press **Mark as inaccurate**. The analysis stops setting the drift level and Steer back, and nothing is sent or re-run. Press **Undo** to bring it back.

> **[Screenshot to take:** a departure with Steer back open, showing the composed correction.**]**

## Let Compass check while you are away

Compass can also check annotated sessions without being asked, and raise only the departures. It spends model capacity with nobody watching, so it is off until you turn it on. It needs your remembered Allow for Codex, reads your goal and not your outcome lines, and stops at a per-session and a per-day limit of its own (twelve a day). Start Compass with the option:

```text
python3 "$SKILL/server.py" --daemon --unasked-readings
```

> **[command may change with one install]**

Finding `$SKILL` is in the [reference](reference.md#start-and-stop). What Compass raised shows on the session page under **Raised while you were away**, and in the Intent log.

## Get alerted when work reaches a step

With [Workflow](../workflow/index.md) running, Compass can tell you once when a work item reaches a step you choose.

1. **Open the project** from Projects, then its **Course** tab.

2. **Find Stage alerts, and choose a step.** The list holds the steps your workflow declares.

3. **Press Save.** The next time an item enters that step, Compass sends one notification. Items already in the step do not trigger it. Press **Rearm** to be alerted again.

Notifications are native on macOS, and browser notifications elsewhere while a Compass tab is open.

## Next

- [Compass reference](reference.md): every command, option and key.
- **Workflow.** [Spacedock Workflow](../workflow/index.md) keeps work moving through steps you define, and Compass shows where each item is.
