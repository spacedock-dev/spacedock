# Install Spacedock

Spacedock works with a coding agent you already have: Claude Code, Codex, or
Pi. Install one of those first.

=== "macOS (Homebrew)"

    ```bash
    brew tap spacedock-dev/homebrew-tap
    brew install spacedock
    ```

    `brew install` also pulls in `agentsview` (it powers `/spacedock:survey`).
    The optional sandbox, safehouse, is installed separately — see
    [Sandboxing](../reference/sandbox.md).

=== "Binary (macOS / Linux)"

    ```bash
    curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
    ```

    Installs a checksum-verified binary to `~/.local/bin`.

## Launch

In a project you already have:

```bash
spacedock claude "/spacedock:survey"
```

Or launch directly:

```bash
spacedock claude "what can spacedock do for me in this project"
```

Replace `claude` with `codex` or `pi` for the respective coding agents.

## Skills

Spacedock installs the relevant skills on launch. To install them manually:

```bash
claude plugin marketplace add spacedock-dev/spacedock
claude plugin install spacedock@spacedock
```

## Sandboxing

See [supported sandboxes](../reference/sandbox.md).

## Troubleshooting

Run `spacedock doctor`.

## Next

Read about the [survey report](survey.md) to understand your usage pattern
with coding agents, or start with [your first workflow](first-workflow.md).
