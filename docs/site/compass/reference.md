# Compass reference

## Start and stop

| Command | What it does | Default |
| --- | --- | --- |
| `/spacedock:compass` | Starts Compass from Claude Code. | Port 4553 |
| `$spacedock:compass` | Starts Compass from Codex. | Port 4553 |
| `spacedock compass --daemon` | Starts Compass in the background. Without `--daemon` it runs in the foreground. | Foreground |
| `spacedock compass --status` | Says whether Compass is running on the port. | Port 4553 |
| `spacedock compass --stop` | Stops Compass on the port. | Port 4553 |
| `spacedock compass --diagnose` | Lists where Compass looks for each tool's sessions and whether it can read them. Starts nothing. `--json` prints it as JSON. | Off |
| `spacedock compass --forget` | Deletes observation history, recorded session ends, stored session summaries and copied-correction records. Pauses model calls and turns summaries off until you turn them back on in the page. Keeps the goals you typed. Refused while Compass is running on the port. | Off |
| `spacedock compass --version` | Prints the Compass version. | Off |
| `spacedock compass --sessions [waiting]` | Lists running and recently finished sessions; with `waiting`, lists only sessions waiting on you, excluding the calling Claude Code session. Needs a running Compass instance. | No waiting filter |
| `spacedock compass --footer` | Prints a Claude Code status-line segment from its input and Compass's saved state. Sends no request. May be used alone or with `--session`. | Silent when nothing can be reported |
| `spacedock compass --help` | Prints the installed Compass options, including three experimental terminal and frontend development options omitted here. | Off |

```bash
spacedock compass --port 4581 --daemon
spacedock compass --port 4581 --status
spacedock compass --port 4581 --stop
```

## Options

Server options apply to the run you start. Stop and restart Compass to change one. Spell each option in full: an abbreviation is refused.

| Option | What it does | Default |
| --- | --- | --- |
| `--port PORT` | Port to serve on. Also selects the instance for `--status`, `--stop` and `--forget`. | `4553` |
| `--window-hours HOURS` | Hides sessions with no activity in this window. | `24` |
| `--unasked-readings` | Checks eligible sessions with a typed goal for drift without being asked, and raises departures only. Claude Code is preferred for Pi and Antigravity sessions; a Codex reader needs permission for the destination shown in the Drift card. Spends model capacity: at most three checks per session and twelve a day. | Off |
| `--session HARNESS:SESSION_ID` | Selects the session for `--footer`, for example `claude:SESSION_ID`. Used only with `--footer`. | Status-line session ID, then `CLAUDE_CODE_SESSION_ID` |
| `--claude-reading-model MODEL` | Claude model that reads a Claude Code session when you press Analyze drift. Takes an explicit Sonnet or Opus ID of generation 5 or later. | `claude-sonnet-5-5` |
| `--claude-monitor-model MODEL` | Claude model that Live monitor asks about once a minute. Takes an explicit Haiku ID of generation 4 or later. | `claude-haiku-5-5` |
| `--no-observer-model`, `--no-harness-usage` | Refuse every model call for the run, including Analyze drift, Find goal, summaries, Live monitor and unasked checks. **Pause model calls** in the page does the same until you resume. | Off |
| `--no-annotations` | Does not read or save goals or outcome lines. Analyze drift stays visible and disabled. | Off |
| `--no-git` | Does not run the end-of-session git check in your repositories. | Off |
| `--no-ci-read` | Does not read pull-request or merge-request checks from GitHub or GitLab when you press Analyze drift. | Off |
| `--no-focus` | Does not raise a session's terminal or record terminal identities. | Off |
| `--no-dismiss` | Does not read or write the store of sessions marked handled. | Off |
| `--no-ask` | Stops sessions from asking you questions through Compass. | Off |
| `--no-events` | Stops the event coordinator, so state comes from scanning session files on a timer. | Off |
| `--no-irreversible` | Stops reporting of risky command shapes such as force push and hard reset. | Off |
| `--reach-url URL` | Webhook that receives counts-only nudges when sessions need input while you are away. | None |
| `--no-reach` | Refuses off-machine nudges whatever URL is configured. | Off |
| `--history-max-bytes BYTES` | Size cap on the history file, also its read cap. | `1048576` |
| `--history-days DAYS` | Days of observed history to keep. | `14` |
| `--no-history` | Does not read or write observed history. | Off |
| `--no-spacedock` | Does not read Workflow definitions, so stage strips are not shown. | Off |
| `--no-tripwires` | Does not read or save stage alerts. | Off |
| `--quiet-hours WINDOW` | Silences non-question notifications, including unasked-check notifications and stage alerts, in a local window such as `22:00-08:00`. Unasked model checks still run. | None |
| `--no-quiet-hours` | Ignores any configured quiet hours. | Off |
| `--no-usage` | Does not read Claude Code and Cursor quota from the vendor for the run. Quota is read otherwise, and sends no session content. **Turn off** under the quota windows on the Sessions screen stops it from the page. | Off |
| `--host ADDRESS` | `127.0.0.1`, or `0.0.0.0` to serve every network interface with no authentication. Prefer an SSH tunnel. | `127.0.0.1` |
| `--observer-model`, `--allow-tool-output` | Accepted but do nothing; will be removed in a later release. | Off |

```bash
spacedock compass --port 4581 --window-hours 6 --no-observer-model
spacedock compass --port 4581 --daemon --unasked-readings --quiet-hours 22:00-08:00
spacedock compass --port 4581 --no-annotations --no-history
spacedock compass --port 4581 --daemon --reach-url https://ntfy.sh/my-topic
```

## Environment

| Variable | Effect |
| --- | --- |
| `COMPASS_HOME` | Moves Compass's own state, logs and stores. Default: `.compass` under your home directory. Use the same value for start, stop, status and forget. |
| `COMPASS_REACH_URL` | Reach webhook when `--reach-url` is absent. |
| `COMPASS_QUIET_HOURS` | Quiet-hours window when `--quiet-hours` is absent. |

## What works on each tool

Every session from the ten agent tools shows on the board. The rest depends on what each tool records.

| Feature | Claude Code | Codex | Pi | Antigravity | Other tools |
| --- | :-: | :-: | :-: | :-: | :-: |
| Session on the board | ✓ | ✓ | ✓ | ✓ | ✓ |
| Goal filled from `/goal` | ✓ | ✓ | | ✓ | |
| Goal drafted from your first prompt | ✓ | ✓ | | | |
| Outline of what the session has done | ✓ | ✓ | | | |
| Model-written summary | ✓ | | | | |
| Drift level | ✓ | ✓; checks need codex-cli 0.149+ | ✓ | | |
| Goal-met claim checked against the record | ✓ | ✓ | | | |
| Live monitor | ✓ | | | | |
| Steer back | ✓ | | | | |
| Unasked checks | ✓ | With the Drift card toggle | ✓; Codex fallback needs permission | ✓; directions only, Codex fallback needs permission | |

## Drift levels

Shown after Analyze drift, and by Live monitor on Claude Code sessions. Live monitor's level, labelled **Live monitor · Haiku**, is Haiku's read of where the latest turns are heading. Until its check returns, or while it is paused or capped, the page shows an estimate from checks and written files, which does not read what your intent says. The meanings below are that estimate's and the analysis's.

| Level | Meaning |
| --- | --- |
| None or low | No recorded signal against your intent. |
| Medium | A check passed and files were written after it, a failed check was followed by writes, or the agent made a claim the record contradicts. For the live estimate, some writes outside the folders your intent names also contribute Medium. |
| High | A check still fails after your last message with no later write. For the live estimate, most writes outside the folders your intent names also count as High. |
| Extreme | Both High conditions hold. An analysis never reads Extreme. |
| Not enough recorded yet | Too little in the record to estimate, or a goal too vague to read the work against. |

## Keys

| Key | Does |
| --- | --- |
| <kbd>s</kbd> | Opens Sessions |
| <kbd>p</kbd> | Opens Projects |
| <kbd>a</kbd> | Opens Attention |
| <kbd>/</kbd> | Focuses session search on the Sessions screen |
| <kbd>Escape</kbd> | Returns from a session to where you opened it, and from any other view to Sessions |
