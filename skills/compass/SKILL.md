---
name: compass
description: "Open Spacedock Compass, the local dashboard of what your coding-agent sessions are doing: sessions, subagents, progress and drift. Use when the user wants to see, monitor or check on agent sessions, or asks to open Compass or Cargento."
---

# Open Compass

Compass runs a local dashboard on this machine. It reads agent session files
and uploads nothing.

1. Check it is there:

   ```bash
   command -v spacedock-compass
   ```

   If this prints nothing, Compass is not installed. Tell the user to reinstall
   spacedock, which installs Compass with it, and stop.

2. If a dashboard is already running, reuse it:

   ```bash
   spacedock compass --port 4553 --status
   ```

3. Otherwise start it in the background:

   ```bash
   spacedock compass --port 4553 --daemon
   ```

4. Give the user http://127.0.0.1:4553/.

If `spacedock compass` says it needs Python 3.11 or newer, pass that line on as
written: the user needs a newer `python3`, or `CARGENTO_PYTHON` pointing at one.
To stop the dashboard: `spacedock compass --port 4553 --stop`.
