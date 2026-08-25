# FO Install (deferred)

## Channel selection

Classify the channel from the retained `{first_officer_base}`. Substitute it for `B` and run:

    V=$(basename "$(dirname "$(dirname "$B")")"); case "$B" in */spacedock-edge/*) echo edge ;; *) case "$V" in [0-9]*-*) echo edge ;; *) echo stable ;; esac ;; esac

## Install commands

Detect the OS with `uname -s`, never `doctor`/`OS:`. Pick the row for the detected OS and the classified channel.

- **Linux, stable:** `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh`
- **Linux, edge:** `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | SPACEDOCK_CHANNEL=edge sh`
- **macOS:** `brew tap spacedock-dev/tap`, then stable `brew install spacedock-dev/tap/spacedock`, edge `brew install spacedock-dev/tap/spacedock@next`
- **Any other OS:** hint `go build -o spacedock ./cmd/spacedock`, ABORT — no package install.

On Linux edge the assignment binds to `sh`, not to `curl`: prefixing `curl` leaves the variable unset in the script and silently installs stable.

## Sandbox

In a sandbox, per the shared core's Startup check: do not offer or run the install; print the channel-correct command and tell the human to run it outside the sandbox.

## Install-and-resume offer

Outside a sandbox only.

1. **One attempt ever — sentinel `${TMPDIR:-/tmp}/spacedock-install-attempted`.** If it exists, skip the offer and hint-and-abort: print the channel-correct command to run manually, say the sentinel exists, and say `rm ${TMPDIR:-/tmp}/spacedock-install-attempted` re-enables the offer. If `SPACEDOCK_BIN` was set-but-stale, name its session-scoped replacement. Same message on the step-4 fallback.
2. **Sentinel absent — OFFER:** "Run the install for you and resume startup once the binary lands?" On decline: ABORT with the manual command; once `spacedock` is on PATH, start the first officer with the host's spacedock launcher command.
3. **On approval — `touch` the sentinel BEFORE the install runs**, so an attempt that dies mid-run still counts as the one attempt. Then run the command exactly once.
4. **Converge (session-scoped repoint).** `install.sh` prints `install.sh: installed spacedock <version> to <dir>/spacedock` to stderr and warns when the dir is off PATH. Parse the installed path from that stderr line; if absent, probe `$HOME/.local/bin/spacedock`. If a path resolves, set `SPACEDOCK_BIN` to it **for this session only — never persist it to a shell profile** — and re-check `${SPACEDOCK_BIN:-spacedock} --version`. If line 1 parses to a compatible version, that repointing IS the gate's one launcher resolution (blessed by the shared-core invariant): resume Startup. Next session: if the install dir is on PATH, bare `spacedock` resolves with no override; if not, the gate fails again and the fallback hint names the exact installed path and tells the human to add the dir to PATH (or launch with `SPACEDOCK_BIN=<path>`).
5. **Fall back, never loop.** If the re-check still fails, or no path resolves, fall back to the step-1 hint-and-abort — **no second install attempt, no proceeding without the re-check**.
