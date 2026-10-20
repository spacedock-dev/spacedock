# Compass reference

## Start and stop

Find the newest launcher your agent installed, in Claude Code or Codex, then reuse it for every command below:

```bash
SKILL="$(dirname "$(find ~/.claude/plugins/cache ~/.codex/plugins/cache -name server.py -path '*compass*' 2>/dev/null | sort -V | tail -1)")"
```

> **[command may change with one install]**

| Command | What it does | Default |
| --- | --- | --- |
| `/compass:compass` | Starts Compass from Claude Code. In Codex, `$compass`. | Port 4553 |
| `python3 "$SKILL/server.py" --daemon` | Starts Compass in the background. Without `--daemon` it runs in the foreground. | Foreground |
| `python3 "$SKILL/server.py" --status` | Says whether Compass is running on the port. | Port 4553 |
| `python3 "$SKILL/server.py" --stop` | Stops Compass on the port. | Port 4553 |
| `python3 "$SKILL/server.py" --diagnose` | Lists where Compass looks for each tool's sessions and whether it can read them. Starts nothing. `--json` prints it as JSON. | Off |
| `python3 "$SKILL/server.py" --forget` | Deletes observation history, recorded session ends and copied-correction records, and clears the permission to analyze. Keeps the goals you typed. Refused while Compass is running on the port. | Off |

```bash
python3 "$SKILL/server.py" --port 4581 --daemon
python3 "$SKILL/server.py" --port 4581 --status
python3 "$SKILL/server.py" --port 4581 --stop
```

## Options

Options apply to the run you start. Stop and restart Compass to change one. Spell each option in full: an abbreviation is refused.

| Option | What it does | Default |
| --- | --- | --- |
| `--port PORT` | Port to serve on. Also selects the instance for `--status`, `--stop` and `--forget`. | `4553` |
| `--window-hours HOURS` | Hides sessions with no activity in this window. | `24` |
| `--unasked-readings` | Checks annotated sessions for drift without being asked, and raises departures only. Needs your remembered Allow for Codex, and spends model capacity. | Off |
| `--claude-reading-model MODEL` | Claude model that reads a Claude Code session when you press Analyze drift. Takes an explicit Sonnet or Opus ID of generation 5 or later. | `claude-sonnet-5-5` |
| `--allow-tool-output` | Allows, at launch, sending a Claude Code reading's checks and the paths it wrote, for the destination Compass names. The first press still asks for the words. `--forget` and **Turn off readings** clear it. | Off |
| `--no-observer-model`, `--no-harness-usage` | Refuse every model call for the run, including Analyze drift, goal summaries and unasked checks. | Off |
| `--no-annotations` | Does not read or save goals or outcome lines. Analyze drift stays visible and disabled. | Off |
| `--observer-model` | Offers model-written goal summaries in Console. Asks for its own consent in the browser. Does not turn on unasked checks. | Off |
| `--no-git` | Does not run the end-of-session git check in your repositories. | Off |
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
| `--quiet-hours WINDOW` | Silences non-urgent notifications, unasked checks and stage alerts in a local window, for example `22:00-08:00`. | None |
| `--no-quiet-hours` | Ignores any configured quiet hours. | Off |
| `--no-usage` | Does not fetch usage quota from the vendor. | Off |
| `--host ADDRESS` | `127.0.0.1`, or `0.0.0.0` to serve every network interface with no authentication. Prefer an SSH tunnel. | `127.0.0.1` |

```bash
python3 "$SKILL/server.py" --port 4581 --window-hours 6 --no-observer-model
python3 "$SKILL/server.py" --port 4581 --daemon --unasked-readings --quiet-hours 22:00-08:00
python3 "$SKILL/server.py" --port 4581 --no-annotations --no-history
python3 "$SKILL/server.py" --port 4581 --daemon --reach-url https://ntfy.sh/my-topic
```

`python3 "$SKILL/server.py" --help` prints the options for the copy you installed. Three options for experimental session terminals and contributor frontend work are left out.

## Environment

| Variable | Effect |
| --- | --- |
| `COMPASS_HOME` | Moves Compass's own state, logs and stores. Default: `.compass` under your home directory. Use the same value for start, stop, status and forget. |
| `COMPASS_ALLOW_TOOL_OUTPUT` | Set to `1` to start as `--allow-tool-output` does. |
| `COMPASS_REACH_URL` | Reach webhook when `--reach-url` is absent. |
| `COMPASS_QUIET_HOURS` | Quiet-hours window when `--quiet-hours` is absent. |

## Drift levels

Shown by Live monitor on Claude Code sessions, and after Analyze drift. Live monitor reads checks and written files, not what your intent says.

| Level | Meaning |
| --- | --- |
| None or low | No recorded signal against your intent. |
| Medium | A check passed and files were written after it, a failed check was followed by writes, some writes fall outside the folders your intent names, or the agent made a claim the record contradicts. |
| High | A check failed after your last message, or most writes fall outside the folders your intent names. |
| Extreme | Both High conditions hold. Live monitor only: an analysis never reads Extreme. |
| Not enough recorded yet | Too little in the record to estimate. Normal until you add an outcome line. |

## Keys

| Key | Does |
| --- | --- |
| <kbd>s</kbd> | Opens Sessions |
| <kbd>p</kbd> | Opens Projects |
| <kbd>a</kbd> | Opens Attention |
| <kbd>Escape</kbd> | Returns from a session to where you opened it, and from any other view to Sessions |
