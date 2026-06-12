# Install Spacedock

Spacedock works with a coding agent you already have: Claude Code, Codex, or
Pi. Install one of those first.

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

Run `spacedock doctor`. It checks the install and names anything missing or
out of date.

## Next

[Your first launch](first-launch.md) covers what survey reports and the offer
it ends with.
