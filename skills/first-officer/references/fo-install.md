# FO Install (deferred)

Loaded by the first officer whenever the version gate (Startup step 1, shared core) lands in the **binary-absent class**, sandbox or not — reading this file is always safe; only *running* an install inside a sandbox is forbidden. Never loaded on the happy path. This file owns channel classification, the per-OS install commands, the sandbox arm, and the install-and-resume offer. The shared core keeps only the trigger and the never-install-inside-a-sandbox rule.

## Channel selection

Classify the channel from the retained `{first_officer_base}`, then hint that channel's command. Substitute the base for `B` and run:

    R="${B%/skills/first-officer}"; V="${R##*/}"; case "$R" in */spacedock-edge/*) echo edge;; *) case "$V" in [0-9]*-*) echo edge;; [0-9]*) echo stable;; *) echo local;; esac;; esac

- **Why two signals.** The marketplace segment is the documented channel name, but a dev build can be installed under the stable marketplace name while carrying a prerelease-suffixed version; the suffix catches that. Conversely an edge-marketplace install can carry an unsuffixed version; the marketplace segment catches that. Either signal alone misclassifies a real observed install shape.
- **`local`.** A base with no version-shaped directory is a `--plugin-dir` source checkout. Do not hint a package install: the repo is present, so hint `go build -o spacedock ./cmd/spacedock` and ABORT. Guessing a channel here would install over a tree the human is editing.
- **Never widen the match.** Test the version segment, never the whole path: a home directory whose name contains a prerelease-like token would otherwise force every install on that machine onto the edge channel.

## Install commands

Detect the OS with `uname -s`, never `doctor`/`OS:`. Pick the row for the detected OS and the classified channel.

- **Linux, stable:** `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh`
- **Linux, edge:** `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | SPACEDOCK_CHANNEL=edge sh`
- **macOS:** `brew tap spacedock-dev/tap`, then stable `brew install spacedock-dev/tap/spacedock`, edge `brew install spacedock-dev/tap/spacedock@next`
- **`local` channel, or any other OS:** hint `go build -o spacedock ./cmd/spacedock`, ABORT — no package install.

**Where the env var goes.** On Linux edge the assignment binds to `sh`, not to `curl`. A shell variable prefix applies to the first command of a pipeline, so prefixing `curl` leaves the variable unset in the script and silently installs stable — the exact channel skew this file exists to close.

## Install-and-resume offer

Outside a sandbox only — inside, `## Sandbox arm` stops at the printed command and this section never runs.

1. **Loop bound — durable sentinel file, one attempt.** Before OFFERING to run the install, check `test -f <sentinel>` where the sentinel is `${TMPDIR:-/tmp}/spacedock-install-attempted-<key>` and `<key>` is the host runtime's session-identity env var — `$CLAUDE_CODE_SESSION_ID` (claude) or `$CODEX_THREAD_ID` (codex). Pi exposes no session-identity env var, so on pi `<key>` falls back to a short hash of the gate's working directory: the sentinel is then **project-scoped and tmp-durable** — a failed install in one session suppresses the install offer in later sessions of the same project until tmp cleanup. That over-reach is admitted and safe: the blocked path is the same hint-and-abort carrying the exact OS-aware command. The sentinel lives only under `${TMPDIR:-/tmp}` — never in the project tree or state checkout.
   - **Sentinel present** → skip the offer; hint-and-abort with the exact OS-aware command. The message MUST tell the human the sentinel exists, print its path, and say `rm <sentinel-path>` re-enables the install offer.
   - **Sentinel absent** → OFFER: "Run the install for you and resume startup once the binary lands?" On decline: ABORT with the manual hint (the exact OS-aware command; once `spacedock` is on PATH, start the first officer with the host's spacedock launcher command).
2. **On approval — create-before-run, then run once.** Create the sentinel BEFORE the install runs (`touch <sentinel>`), so an attempt that dies mid-run still counts as the one attempt. Then run the OS-aware install command exactly once.
3. **Converge (session-scoped repoint).** `install.sh` prints `install.sh: installed spacedock <version> to <dir>/spacedock` to stderr and warns when the dir is off PATH. Parse the installed path from that stderr line; if absent, probe `$HOME/.local/bin/spacedock`. If a path resolves, set `SPACEDOCK_BIN` to it **for this session only — never persist it to a shell profile** — and re-check `${SPACEDOCK_BIN:-spacedock} --version`. If line 1 parses to a compatible version, that repointing IS the gate's one launcher resolution (blessed by the shared-core invariant): resume Startup. Next session: if the install dir is on PATH, bare `spacedock` resolves with no override; if not, the gate fails again and the fallback hint names the exact installed path and tells the human to add the dir to PATH (or launch with `SPACEDOCK_BIN=<path>`).
4. **Fall back, never loop.** If the re-check still fails, or no path resolves, fall back to hint-and-abort — **no second install attempt, no proceeding without the re-check**.

## Fallback-message grammar (both the post-failure fallback and every sentinel-blocked re-entry)

- The exact OS-aware install command for the human to run manually (Linux `curl|sh` lead, macOS brew lead, source-build fallback — the channel-correct row from `## Install commands`).
- The sentinel exists, its printed path, and `rm <sentinel-path>` to re-enable the install offer.
- If `SPACEDOCK_BIN` was set-but-stale (the cause of the gate failure), name its session-scoped replacement.

## Sandbox arm (which branch runs)

This file loads in both cases; the sandbox verdict selects the branch.

The shared core's Startup step 1 performs the sandbox check in EVERY gate class via the session-env registry read — Bash name+VALUE match, e.g. `[ "${APP_SANDBOX_CONTAINER_ID:-}" = "agent-safehouse" ]`, extended to every name+value pair the binary treats as a sandbox marker. Matching is on the VALUE, not mere presence: `APP_SANDBOX_CONTAINER_ID` is a generic macOS app-sandbox variable other containers also set, so presence alone would claim any of them.

- **Inside a sandbox.** Classify the channel and select the OS row exactly as above, then STOP: do NOT offer, and do NOT run, any install — it would land where the host cannot see it, a silent no-op. Skip `## Install-and-resume offer` entirely; touch no sentinel. Message body: "You're running inside a sandbox. Run this command yourself, outside the sandbox: `<exact channel-correct install command>`." The command must still be channel-correct — an in-sandbox human who runs a stable command against an edge plugin lands the same rejected binary.
- **Outside a sandbox.** Proceed to `## Install-and-resume offer` with the selected command.

In binary-present classes the `--version` `^Sandbox: ` line (prefix-anchored) corroborates the env verdict only; on disagreement the env check wins and the message names the disagreement.
