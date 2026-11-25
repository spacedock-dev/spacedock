# Compass reference

Look up Spacedock Compass commands, options and board labels.

## Start and stop

| Command | What it does | Default |
| --- | --- | --- |
| `/spacedock:compass` | Starts Compass from Claude Code. | Port 4553 |
| `$spacedock:compass` | Starts Compass from Codex. | Port 4553 |
| `spacedock compass --daemon` | Starts Compass in the background. | Foreground without the flag |
| `spacedock compass --status` | Reports whether Compass is running on the selected port. | Port 4553 |
| `spacedock compass --stop` | Stops Compass on the selected port. | Port 4553 |
| `spacedock compass --diagnose` | Reports where Compass looks for each agent's records and whether it can read them. Starts nothing. Add `--json` for JSON. | Text |
| `spacedock compass --forget` | Deletes observation history, recorded ends, launcher parent declarations, stored summaries and copied-correction records. Keeps the goals and group labels you typed and the current model-use budget. Pauses AI use and turns summaries off. Requires proof that all dashboards using that Compass home have stopped. | Port 4553 |
| `spacedock compass --version` | Prints the installed Compass version. | |
| `spacedock compass --sessions [waiting]` | Lists sessions from a running board, starting with **Right now** when that group has entries. Add `waiting` to list only sessions waiting on you. Excludes the calling session when Compass can identify it. | All listed sessions |
| `spacedock compass --drift [again]` | Prints a session's progress check. Reads a saved check for free. With no saved check or with `again`, asks the model and waits. | Session from `--session`. Otherwise `CLAUDE_CODE_SESSION_ID` or `CODEX_THREAD_ID`, whichever is set. Requires `--session` when both are set. |
| `spacedock compass --pause-model-calls` | Pauses model use on a running board. Use **Resume AI use** on the board. | Port 4553 |
| `spacedock compass --footer` | Prints a Claude Code status line segment from its input and Compass's saved state, including `AI use paused` when applicable. Sends no request. Accepts `--session`. | Silent when nothing can be reported |
| `spacedock compass --forward-usage` | Reads Claude Code's status line JSON on standard input and sends only the five-hour and seven-day usage percentages and reset times to a running local board. Reads no credential, prints nothing and starts nothing. | Off |
| `spacedock compass --help` | Prints the installed Compass options. | |

```bash
spacedock compass --port 4581 --daemon
spacedock compass --port 4581 --status
spacedock compass --port 4581 --stop
```

For the remaining terminal actions, see the [Compass command manual](https://github.com/spacedock-dev/compass/blob/main/MANUAL.md#options).

## Options

Options are optional and must be spelled in full. Server options apply to the run you start. Stop and restart Compass to change them.

| Option | What it does | Default |
| --- | --- | --- |
| `--port PORT` | Port to serve on. Also selects the instance for `--status`, `--stop` and `--forget`. | `4553` |
| `--window-hours HOURS` | Hides sessions with no activity in this window. | `24` |
| `--unasked-readings` | Checks eligible sessions with a goal you wrote and raises departures only. Claude Code and Codex sessions use their own provider. Other eligible sessions prefer Claude Code. If the preferred CLI is unavailable, the other can check. Codex needs permission for the destination shown in **Is it on track?** At most three checks per session and twelve per day. | Off |
| `--session HARNESS:SESSION_ID` | Selects a session for `--footer`, `--drift`, `--goal`, `--project`, `--summary`, `--monitor` or `--set-group`. Example value: `claude:SESSION_ID`. | For `--footer`, stdin session ID then `CLAUDE_CODE_SESSION_ID`. Other session reads and monitoring use `CLAUDE_CODE_SESSION_ID`. `--drift` also uses `CODEX_THREAD_ID` when it is the only one set, and requires `--session` when both are set. Required for `--set-group`. |
| `--claude-reading-model MODEL` | Model Claude Code uses for **Check progress**. Takes an explicit Sonnet or Opus ID of generation 5 or later. | `claude-sonnet-5-5` |
| `--claude-monitor-model MODEL` | Model Claude Code uses for **Live monitor**. Takes an explicit Haiku ID of generation 4 or later. | `claude-haiku-5-5` |
| `--no-observer-model`, `--no-harness-usage` | Refuse all model calls for the run, including **Check progress**, **Find goal**, summaries, **Live monitor** and unasked checks. | Off |
| `--no-annotations` | Skips saved goals, outcome lines, readings and group labels. **Check progress** remains visible and disabled. | Off |
| `--no-git` | Skips the end-of-session git check and the read of the project's committed check scripts. | Off |
| `--no-ci-read` | Skips pull request and merge request checks from GitHub and GitLab when you press **Check progress**. | Off |
| `--no-focus` | Disables raising a session's terminal and recording its terminal identity. | Off |
| `--no-dismiss` | Skips the store of sessions marked handled and disables marking them. | Off |
| `--no-ask` | Disables questions that sessions can ask through Compass. | Off |
| `--no-events` | Stops collecting pushed session events. Compass scans session files on a timer. | Off |
| `--no-irreversible` | Disables reports of command shapes such as force push and hard reset. | Off |
| `--reach-url URL` | HTTPS webhook for counts-only nudges when sessions need input or end unread. A URL starting with `http://` is refused. | Unset |
| `--no-reach` | Disables off-machine nudges for the run. | Off |
| `--history-max-bytes BYTES` | Maximum size of the history file, also its read limit. | `1048576` |
| `--history-days DAYS` | Days of observed history to keep. | `14` |
| `--no-history` | Skips reading and writing observed history. | Off |
| `--no-spacedock` | Skips Spacedock Workflow state and hides its stage strips. Saved stage alerts are suspended. | Off |
| `--no-tripwires` | Skips reading and saving stage alerts. | Off |
| `--quiet-hours WINDOW` | Suppresses notifications and nudges other than questions in a local window such as `22:00-08:00`. Unasked checks still run. | Unset |
| `--no-quiet-hours` | Ignores any configured quiet hours. | Off |
| `--no-usage` | Stops vendor quota fetches and ignores quota forwarded from a status line. Quota read from an agent's own records can still appear. | Off |
| `--no-update-check` | Skips the daily release lookup. The board also has **Look for new versions once a day** in **Settings**. | Off |
| `--host ADDRESS` | Listen on `127.0.0.1` or expose the board on every IPv4 interface with `0.0.0.0`. Compass provides no authentication on that bind. | `127.0.0.1` |
| `--allow-host NAME` | Trusts an exact DNS hostname for a board behind an authenticated proxy with direct port access blocked. Repeat for more names. Applies with `--host 0.0.0.0`. | Unset |
| `--observer-model`, `--allow-tool-output` | Accepted for compatibility and have no effect. | Off |

```bash
spacedock compass --port 4581 --window-hours 6 --no-observer-model
spacedock compass --port 4581 --daemon --unasked-readings --quiet-hours 22:00-08:00
spacedock compass --port 4581 --daemon --reach-url https://ntfy.sh/my-topic
```

## Environment

| Variable | Effect |
| --- | --- |
| `COMPASS_HOME` | Moves Compass's state, logs and stores. Defaults to `.compass` under your home directory. Use the same value for start, stop, status and forget. |
| `COMPASS_REACH_URL` | HTTPS webhook when `--reach-url` is absent. |
| `COMPASS_QUIET_HOURS` | Quiet-hours window when `--quiet-hours` is absent. |

## What works on each tool

The board lists sessions from the ten supported agent tools when it can find their records.

| Feature | Claude Code | Codex | Pi | Antigravity | Other tools |
| --- | :-: | :-: | :-: | :-: | :-: |
| Session on the board | ✓ | ✓ | ✓ | ✓ | ✓ |
| Goal filled from `/goal` | ✓ | ✓ | | ✓ | |
| Goal drafted from a prompt | ✓ | ✓ | | | |
| Session outline | ✓ | ✓, your messages only | | | |
| Model-written session summary | ✓ | | | | |
| Check progress | ✓ | ✓, findings without a level. Evidence from tests and builds requires codex-cli 0.149 or later. | ✓ | | |
| Compare reported goal completion with the record | ✓ | ✓ | | | |
| Live monitor | ✓ | | | | |
| Steer back | ✓ | | | | |
| Unasked checks | ✓ | ✓, with permission | ✓ | ✓, your prompts only | |

## Progress levels

| Level | Meaning |
| --- | --- |
| **On track** | No departure found in what Compass read. |
| **Slightly off course** | Some work wandered from the goal. |
| **Off course** | The work no longer serves the goal. |
| **Way off course** | The work has gone far from the goal. Both **Live estimate** and **Live monitor** can show this level. A progress check you start shows at most **Off course**. |

A saved progress check shows a level only while the evidence it used still matches. Lost history, an expired direction, a truncated transcript or a changed project key can remove the level. Checks saved before v3.0.0 need **Check progress again** before they can show a level.

Sessions rows can show **Slightly off course** or **Off course** from that same saved evidence. A row says **Not checked yet** only if Compass can check the goal and has no record of a check or attempt. Codex leaves the progress cell blank and shows any findings on the session page.

| Other progress text | Meaning |
| --- | --- |
| **Too early to tell** | The session has recorded too little to check its progress. |
| **Can't tell** | Compass could not establish progress. The card gives the reason. |
| **Goal too vague to check** | The goal needs a clearer end state. |
| **Out of date** | The previous check no longer describes the evidence Compass can use. |

## Keys

| Key | Does |
| --- | --- |
| <kbd>s</kbd> | Opens Sessions. |
| <kbd>p</kbd> | Opens Projects. |
| <kbd>n</kbd> | Opens Needs you. |
| <kbd>/</kbd> | Focuses search on Sessions. |
| <kbd>Escape</kbd> | Returns from a session to where you opened it, and from other pages to Sessions. |

## Upgrading to v3.0.0

| Area | Contract |
| --- | --- |
| Observation routes | Require an admitted browser `Origin` or same-origin fetch metadata. |
| Compass state files | State files covered by the strict state policy must be regular files you own. Symlinks are refused. |
| Scripts matching output | Match `AI use paused` from `--footer`, `Right now` from `--sessions` and `Check progress` in progress messages. The `--drift` flag keeps its name. |

See the [v3.0.0 release notes](https://github.com/spacedock-dev/compass/releases/tag/v3.0.0#why-this-is-a-major-release) and the [Compass security policy](https://github.com/spacedock-dev/compass/blob/main/SECURITY.md#scope) for the full contracts.

**Related guides**

- [Understand what Compass reads and shows](index.md).
- [Check a session's progress on the board](use.md#check-a-sessions-progress).
