# Install Spacedock

Spacedock plugs into a coding agent harness you already run: Claude Code, Codex,
or Pi. Install one of those first. Every command below names the output you
should see, so you can check each step against the stated result.

Spacedock itself is two pieces that install separately:

1. **The `spacedock` launcher.** The command you run to start a session.
2. **The host plugin.** The first-officer and ensign agents, loaded by your
   harness (Claude Code, Codex, or Pi).

The recommended setup installs the launcher with Homebrew, then adds the plugin.
Pick your platform, then confirm and launch. Those last two steps are the same
everywhere. A from-source build is available for development.

## Install the launcher

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
    `PATH`; when the binary is absent, Spacedock prints an install hint and the
    launch proceeds **unsandboxed**. The macOS-only Gatekeeper/quarantine
    handling does not apply on Linux and is not needed there.

=== "Codex or Pi"

    Claude Code is the primary surface; Codex and Pi are supported but
    experimental. Install the launcher with Homebrew (the macOS tab), then add
    the plugin for your host:

    ```bash
    spacedock install --host codex      # or: --host pi
    ```

    Codex installs plugins from your shell rather than programmatically, so this
    prints the `codex plugin` commands to run.

## Confirm it

```bash
spacedock --version
```

Prints the installed version, e.g. `spacedock 0.20.0`.

## Launch

Point Spacedock at a project you already have and let it survey:

```bash
spacedock claude "/spacedock:survey"
```

Starts the first officer in Claude Code and runs the survey. The first launch
sets up the plugin for you, so this single command is enough. When a
`.safehouse` profile is present in the working directory, the launch runs
sandboxed. To set up the plugin ahead of time, or to refresh it later, run
`spacedock install --host claude`.

With Codex or Pi, launch with the matching subcommand instead:
`spacedock codex "your task"` or `spacedock pi "your task"`.

Working on Spacedock itself? See [Build from source](https://github.com/spacedock-dev/spacedock/blob/next/docs/site/contributing/build-from-source.md).
It builds the launcher from the `next` branch and loads the plugin from your
checkout so local edits are live.

## Keep things in sync

`spacedock doctor` is the compatibility check. If it reports your installed
plugin is out of date, refresh it:

```bash
spacedock install --host claude
```

If the `spacedock` command itself is missing, install the launcher with Homebrew
first, then run `spacedock install --host claude`.

## Next

Run your [first launch](first-launch.md). The
[command reference](../reference/command-reference.md#launch)
covers the full launch grammar: the task argument, what rides after `--`, and
the sandbox flags.
