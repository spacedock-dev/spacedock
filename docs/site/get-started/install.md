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
    curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
    ```

    Installs a checksum-verified binary to `~/.local/bin`.

## Launch

Point Spacedock at a project you already have and let it survey:

```bash
spacedock claude "/spacedock:survey"
```

Starts a Spacedock session in Claude Code and runs the survey. The first launch
sets up everything Claude Code needs; there is no separate setup step.

With Codex or Pi (experimental), run `spacedock install --host codex` (or
`--host pi`) once, then launch with the matching subcommand:
`spacedock codex "your task"` or `spacedock pi "your task"`.

## Sandboxing

When a `.safehouse` profile is present in the working directory, the launch
runs sandboxed through the `safehouse` command; `--safehouse` forces it.
Spacedock ships no sandbox of its own: when no `safehouse` binary is on your
`PATH`, the launch prints an install hint and proceeds **unsandboxed**.

## Troubleshooting

Run `spacedock doctor`. It checks the install and names anything missing or
out of date.

## Next

Run your [first launch](first-launch.md). The
[command reference](../reference/command-reference.md#launch)
covers the full launch grammar: the task argument, what rides after `--`, and
the sandbox flags.
