# FO Install Gate (deferred)

Loaded by the first officer ONLY when the version gate (Startup step 1, shared core) lands in the **binary-absent class** and the session is NOT inside a sandbox. Never loaded on the happy path. The shared core owns the OS-aware hint text and the sandbox check; this file owns the direct install-offer machinery.

## Install-and-resume offer

1. **Loop bound — durable sentinel file, one attempt.** Before OFFERING to run the install, check `test -f <sentinel>` where the sentinel is `${TMPDIR:-/tmp}/spacedock-install-attempted-<key>` and `<key>` is the host runtime's session-identity env var — `$CLAUDE_CODE_SESSION_ID` (claude) or `$CODEX_THREAD_ID` (codex). Pi exposes no session-identity env var, so on pi `<key>` falls back to a short hash of the gate's working directory: the sentinel is then **project-scoped and tmp-durable** — a failed install in one session suppresses the install offer in later sessions of the same project until tmp cleanup. That over-reach is admitted and safe: the blocked path is the same hint-and-abort carrying the exact OS-aware command. The sentinel lives only under `${TMPDIR:-/tmp}` — never in the project tree or state checkout.
   - **Sentinel present** → skip the offer; hint-and-abort with the exact OS-aware command. The message MUST tell the human the sentinel exists, print its path, and say `rm <sentinel-path>` re-enables the install offer.
   - **Sentinel absent** → OFFER: "Run the install for you and resume startup once the binary lands?" On decline: ABORT with the manual hint (the exact OS-aware command; once `spacedock` is on PATH, start the first officer with the host's spacedock launcher command).
2. **On approval — create-before-run, then run once.** Create the sentinel BEFORE the install runs (`touch <sentinel>`), so an attempt that dies mid-run still counts as the one attempt. Then run the OS-aware install command exactly once.
3. **Converge (session-scoped repoint).** `install.sh` prints `install.sh: installed spacedock <version> to <dir>/spacedock` to stderr and warns when the dir is off PATH. Parse the installed path from that stderr line; if absent, probe `$HOME/.local/bin/spacedock`. If a path resolves, set `SPACEDOCK_BIN` to it **for this session only — never persist it to a shell profile** — and re-check `${SPACEDOCK_BIN:-spacedock} --version`. If line 1 parses to a compatible version, that repointing IS the gate's one launcher resolution (blessed by the shared-core invariant): resume Startup. Next session: if the install dir is on PATH, bare `spacedock` resolves with no override; if not, the gate fails again and the fallback hint names the exact installed path and tells the human to add the dir to PATH (or launch with `SPACEDOCK_BIN=<path>`).
4. **Fall back, never loop.** If the re-check still fails, or no path resolves, fall back to hint-and-abort — **no second install attempt, no proceeding without the re-check**.

## Channel selection

The shared core's Binary-absent bullet carries the classifier itself, because the sandbox arm must name a channel-correct command without loading this file. This section is the rationale it defers.

- **Why two signals.** The marketplace segment is the documented channel name, but a `next`-branch dev build can be installed under the stable marketplace name (`cache/spacedock/spacedock/0.27.0-pre7-dev`); a prerelease suffix catches that. Conversely an edge install can carry an unsuffixed version (`cache/spacedock-edge/spacedock/0.27.0`); the marketplace segment catches that. Either signal alone misclassifies a real observed install.
- **`local`.** A base with no version-shaped directory is a `--plugin-dir` source checkout. Do not hint a package install: the repo is present, so hint `go build -o spacedock ./cmd/spacedock` and ABORT. Guessing a channel here would install over a tree the human is editing.
- **Never widen the match.** Test the version segment, never the whole path: a home directory such as `/Users/pre-release-tester` would otherwise force every install on that machine onto the edge channel.
- **Where the env var goes.** On Linux edge the assignment binds to `sh`, not to `curl`: `curl … | SPACEDOCK_CHANNEL=edge sh`. A shell variable prefix applies to the first command of a pipeline, so prefixing `curl` leaves the variable unset in the script and silently installs stable — the exact channel skew this gate exists to close.

## Fallback-message grammar (both the post-failure fallback and every sentinel-blocked re-entry)

- The exact OS-aware install command for the human to run manually (Linux `curl|sh` lead, macOS brew lead, source-build fallback — per the shared core's hint).
- The sentinel exists, its printed path, and `rm <sentinel-path>` to re-enable the install offer.
- If `SPACEDOCK_BIN` was set-but-stale (the cause of the gate failure), name its session-scoped replacement.

## Sandbox arm (message body)

The shared core's Startup step 1 performs the sandbox check in EVERY gate class via the session-env registry read — Bash name+VALUE match, e.g. `[ "${APP_SANDBOX_CONTAINER_ID:-}" = "agent-safehouse" ]`, extended to every name+value pair the binary treats as a sandbox marker. Matching is on the VALUE, not mere presence: `APP_SANDBOX_CONTAINER_ID` is a generic macOS app-sandbox variable other containers also set, so presence alone would claim any of them. When the registry check says inside: do NOT offer or run any install (it would land where the host cannot see it — a silent no-op). Message body: "You're running inside a sandbox. Run this command yourself, outside the sandbox: `<exact OS-aware install command>`." In binary-present classes the `--version` `^Sandbox: ` line (prefix-anchored) corroborates the env verdict only; on disagreement the env check wins and the message names the disagreement.
