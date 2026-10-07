# Compass reference

## Start and stop

Find the newest launcher your agent installed, in Claude Code or Codex, then reuse it for every command below:

```bash
SKILL="$(dirname "$(find ~/.claude/plugins/cache ~/.codex/plugins/cache -name server.py -path '*cargento*' 2>/dev/null | sort -V | tail -1)")"
```

> **[command may change with one install]**

| Command | What it does | Default |
| --- | --- | --- |
| `/cargento:cargento` | Starts Compass from Claude Code. In Codex, `$cargento`. | Port 4553 |
| `python3 "$SKILL/server.py" --daemon` | Starts Compass in the background. Without `--daemon` it runs in the foreground. | Foreground |
| `python3 "$SKILL/server.py" --status` | Says whether Compass is running on the port. | Port 4553 |
| `python3 "$SKILL/server.py" --stop` | Stops Compass on the port. | Port 4553 |
| `python3 "$SKILL/server.py" --diagnose` | Lists where Compass looks for each tool's sessions and whether it can read them. Starts nothing. `--json` prints it as JSON. | Off |
| `python3 "$SKILL/server.py" --forget` | Deletes observation history and recorded session ends, and clears the permission to analyze. Keeps the goals you typed. Refused while Compass is running. | Off |

```bash
python3 "$SKILL/server.py" --port 4581 --daemon
python3 "$SKILL/server.py" --port 4581 --status
python3 "$SKILL/server.py" --port 4581 --stop
```

## Options

Options apply to the run you start. Stop and restart Compass to change one.

| Option | What it does | Default |
| --- | --- | --- |
| `--port PORT` | Port to serve on. Also selects the instance for `--status`, `--stop` and `--forget`. | `4553` |
| `--window-hours HOURS` | Hides sessions with no activity in this window. | `24` |
| `--unasked-readings` | Checks annotated sessions for drift without being asked, and raises departures only. | Off |
| `--claude-reading-model MODEL` | Claude model that reads a Claude Code session when you press Analyze drift. Takes an explicit Sonnet or Opus ID of generation 5 or later. | `claude-sonnet-5-5` |
| `--no-observer-model` | Refuses every model call for the run, including Analyze drift, goal summaries and unasked checks. | Off |
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
| `--no-tripwires` | Does not read or save step alerts. | Off |
| `--quiet-hours WINDOW` | Silences non-urgent notifications and unasked checks in a local window, for example `22:00-08:00`. | None |
| `--no-quiet-hours` | Ignores any configured quiet hours. | Off |
| `--no-usage` | Does not fetch usage quota from the vendor. | Off |
| `--host ADDRESS` | `127.0.0.1`, or `0.0.0.0` to serve every network interface with no authentication. Prefer an SSH tunnel. | `127.0.0.1` |

```bash
python3 "$SKILL/server.py" --port 4581 --window-hours 6 --no-observer-model
python3 "$SKILL/server.py" --port 4581 --daemon --unasked-readings --quiet-hours 22:00-08:00
python3 "$SKILL/server.py" --port 4581 --no-annotations --no-history
python3 "$SKILL/server.py" --port 4581 --daemon --reach-url https://ntfy.sh/my-topic
```

`python3 "$SKILL/server.py" --help` prints the options for the copy you installed. Two experimental options for session terminals are left out.

## Drift levels

Shown after Analyze drift, and by Live monitor on Claude Code sessions.

| Level | Meaning |
| --- | --- |
| None or low | No recorded signal against your intent. |
| Medium | A claim the record contradicts, or a failed check followed by file writes. |
| High | A cited signal after your last message. |
| Extreme | `[Definition to confirm]` |
| Not enough recorded yet | Too little in the record to estimate. |

> **Draft note, not page copy:** the v0.28.0 skill text defines Medium only in part and does not define Extreme. Confirm the Medium, High and Extreme meanings with the Cargento team before 10-09.

## Keys

| Key | Does |
| --- | --- |
| <kbd>s</kbd> | Opens Sessions |
| <kbd>p</kbd> | Opens Projects |
| <kbd>a</kbd> | Opens Attention |
| <kbd>Escape</kbd> | Returns from a session to where you opened it |
