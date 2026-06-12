# Install Spacedock

Spacedock plugs into a coding agent harness you already run: Claude Code, Codex,
or Pi. Install one of those first. Every command below names the output to
expect, so you can check each step.

=== "macOS (Homebrew)"

    ```bash
    brew tap spacedock-dev/homebrew-tap
    brew install spacedock
    ```

=== "Linux / no Homebrew"

    The Homebrew cask is macOS-only. On Linux (or on macOS without Homebrew),
    the `curl | sh` script detects your OS and architecture, downloads the
    matching tarball from the latest GitHub Release, verifies it against the
    release `checksums.txt`, and installs the `spacedock` binary to
    `~/.local/bin`.

    ```bash
    curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/next/install.sh | sh
    ```

    If `~/.local/bin` is not on your `PATH`, the script prints a note; add it
    (`export PATH="$HOME/.local/bin:$PATH"`) so the `spacedock` command resolves.
    Set `SPACEDOCK_INSTALL_DIR` to install elsewhere (e.g.
    `SPACEDOCK_INSTALL_DIR=/usr/local/bin`, which may need `sudo`).

    **Sandboxing.** Safehouse behaves the same as on macOS: when a `.safehouse`
    profile is present in the working directory, Spacedock wraps the launch
    through the `safehouse` command. Spacedock ships no sandbox of its own. A
    run is sandboxed only when a Linux-capable `safehouse` binary is on your
    `PATH`; when it is absent, Spacedock prints an install hint and the launch
    proceeds **unsandboxed**. Gatekeeper and quarantine handling are macOS-only
    and do not apply here.

=== "Codex or Pi"

    Claude Code is the primary surface; Codex and Pi are supported but
    experimental. Install with Homebrew (the macOS tab) or the Linux script,
    then add Spacedock's agents to your host:

    ```bash
    spacedock install --host codex      # or: --host pi
    ```

    Codex installs from your shell rather than programmatically, so this
    prints the `codex plugin` commands to run.

## Launch

Point Spacedock at a project you already have and let it survey:

```bash
spacedock claude "/spacedock:survey"
```

Starts the first officer in Claude Code and runs the survey. Spacedock manages
the agents your harness loads: the first launch sets them up in Claude Code, so
no separate setup step is needed. To set them up ahead of time, or to refresh
them later, run `spacedock install --host claude`. When a `.safehouse` profile
is present in the working directory, the launch runs sandboxed.

With Codex or Pi, launch with the matching subcommand instead:
`spacedock codex "your task"` or `spacedock pi "your task"`.

Working on Spacedock itself? See [Build from source](https://github.com/spacedock-dev/spacedock/blob/next/docs/site/contributing/build-from-source.md).
It builds from the `next` branch and loads the agents from your checkout so
local edits are live.

## Keep things in sync

`spacedock doctor` checks compatibility. If it reports the agents your harness
loads are out of date, refresh with:

```bash
spacedock install --host claude
```

If the `spacedock` command itself is missing, install it with Homebrew first,
then run `spacedock install --host claude`.

## Next

Run your [first launch](first-launch.md). The
[command reference](../reference/command-reference.md#launch)
covers the full launch grammar: the task argument, what rides after `--`, and
the sandbox flags.
