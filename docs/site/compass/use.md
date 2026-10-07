# Use Compass

State an intent for a session, read the board, and act when an agent drifts from it.

You need Spacedock installed with Compass enabled. See [Install Spacedock and enable Compass](../index.md#install).

Start Compass from your agent session:

```text
/cargento:cargento
```

> **[command may change with one install]**

In Codex, type `$cargento`. You can also ask your agent to "open cargento". Compass starts on your machine and gives you its address, `http://127.0.0.1:4553/` by default. Open it in a browser.

> **[Screenshot to take:** the Sessions screen right after Compass opens.**]**

## Set the goal

1. **Open the session.** On the Sessions screen, click the session's row. Its page opens with an **Intent** section at the top.

2. **Type the goal.** One line, up to 240 characters. If you typed nothing, Compass drafts one from your first prompt in the session, tinted and marked "from your prompt". The draft is not saved until you save it.

3. **Add the expected outcome.** Up to six lines, one check per line. Press **add a line** for each. This step is optional.

4. **Press Save intent.** Compass keeps the words you typed and shows them beside the session from now on. Nothing is sent to the agent.

> **[Screenshot to take:** the Intent section with a typed goal and three outcome lines, saved.**]**

## Read the board

Sessions puts active work first. Each active session answers four questions: where it is, what it is doing now, what it does next, and whether it is blocked. Your goal sits beside what the agent is doing, so a mismatch is visible without opening anything.

| Open | To see |
| --- | --- |
| **Sessions**, key <kbd>s</kbd> | Every session, active work first |
| **Projects**, key <kbd>p</kbd> | Sessions grouped by project, and each project's decisions and recent activity |
| **Attention**, key <kbd>a</kbd> | Sessions waiting on you, and sessions at risk |

Press <kbd>Escape</kbd> in a session to go back.

## Check a session for drift

1. **Press Analyze drift** on the session's page.

2. **Read what it sends, then press Allow and analyze.** The first press shows what Compass will send, and to which agent tool. Compass remembers your answer. **Turn off readings** on any session page takes it back.

3. **Read the result.** It says one of: "Departs from your intent", "Can't tell", "Not reached at this stop", or "Nothing found against what it read". A departure lists each line of your intent it contradicts, with the numbered session entries as evidence.

> **[Screenshot to take:** a "Departs from your intent" result, with the cited entries numbered beside it.**]**

A session with a departure on record gets a **Drift** mark on the Sessions screen and moves up the list.

Each analysis uses model capacity from your Claude Code or Codex account. Compass allows twelve in a rolling twenty-four hours and says when you can run the next one.

For a Claude Code session you can also turn on **Live monitor** on the session page. It re-estimates drift after every turn from the checks the session ran and the files it wrote, with no model call, and shows the result as None or low, Medium, High or Extreme. It never notifies you.

> **Draft note, not page copy:** confirm before 10-09 which harnesses Analyze drift reads on 10-15. At v0.28.0 the Drift section reads Claude Code and Pi sessions and says "Cargento can't read work from this harness" on the rest.

## Act on a departure

Compass never writes into your session. You choose what happens next.

1. **Press Steer back.** Compass composes a correction from your goal and the evidence, ready to copy. Steer back is available for Claude Code sessions.

2. **Paste it into the session.** Rows for Claude Code and Codex sessions carry a control that copies the command to re-enter the session.

If the agent was right and your goal was out of date, press **Add it to my intent** instead, and Compass updates the goal. If the reading is wrong, press **Not accurate?**. That reading stops counting until you clear the mark.

> **[Screenshot to take:** a departure with Steer back open, showing the composed correction.**]**

## Let Compass check while you are away

Compass can also check annotated sessions without being asked, and raise only the departures. It spends model capacity with nobody watching, so it is off until you turn it on:

```text
python3 "$SKILL/server.py" --daemon --unasked-readings
```

> **[command may change with one install]**

> **Draft note, not page copy:** decision pending on whether this ships on or off for 10-15, and how the setting is turned on. Confirm before 10-09. The command above is today's, and finding `$SKILL` is in the [reference](reference.md#start-and-stop).

## Get alerted when work reaches a step

With [Workflow](../workflow/index.md) running, Compass can tell you once when a work item reaches a step you choose.

1. **Open the project** from Projects, then its **Course** tab.

2. **Find Workflow stage conditions, and choose a step.** The list holds the steps your workflow declares.

3. **Press Save.** The next time an item enters that step, Compass sends one notification. Items already in the step do not trigger it. Press **Rearm** to be alerted again.

Notifications are native on macOS, and browser notifications elsewhere while a Compass tab is open.

## Next

- [Compass reference](reference.md): every command, option and key.
- **Workflow.** [Spacedock Workflow](../workflow/index.md) keeps work moving through steps you define, and Compass shows where each item is.
