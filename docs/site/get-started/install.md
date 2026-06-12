# Install Spacedock

Spacedock plugs into a coding agent harness you already run: Claude Code, Codex,
or Pi. Install one of those first.

=== "macOS (Homebrew)"

    ```bash
    brew tap spacedock-dev/homebrew-tap
    brew install spacedock
    ```

=== "Binary (macOS / Linux)"

    ```bash
    curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/next/install.sh | sh
    ```

    Installs a checksum-verified binary to `~/.local/bin`.

## Launch

Point Spacedock at a project you already have and let it survey:

```bash
spacedock claude "/spacedock:survey"
```

Starts the first officer in Claude Code and runs the survey. Spacedock manages
the agents your harness loads: the first launch sets them up in Claude Code, so
no separate setup step is needed. To set them up ahead of time, or to refresh
them later, run `spacedock install --host claude`.

With Codex or Pi (experimental), add the agents with
`spacedock install --host codex` (or `--host pi`), then launch with the
matching subcommand: `spacedock codex "your task"` or `spacedock pi "your task"`.

Working on Spacedock itself? See [Build from source](https://github.com/spacedock-dev/spacedock/blob/next/docs/site/contributing/build-from-source.md).
It builds from the `next` branch and loads the agents from your checkout so
local edits are live.

## Sandboxing

When a `.safehouse` profile is present in the working directory, the launch
runs sandboxed through the `safehouse` command; `--safehouse` forces it.
Spacedock ships no sandbox of its own: when no `safehouse` binary is on your
`PATH`, the launch prints an install hint and proceeds **unsandboxed**.

## Keep things in sync

`spacedock doctor` checks compatibility. If it reports the agents your harness
loads are out of date, refresh with:

```bash
spacedock install --host claude
```

## Next

Run your [first launch](first-launch.md). The
[command reference](../reference/command-reference.md#launch)
covers the full launch grammar: the task argument, what rides after `--`, and
the sandbox flags.
