# Installation

Spacedock should work with most coding agents. This guide covers manual setup for Claude Code and Codex.

## 1. Install Spacedock

Choose your operating system. This installs the programs on your machine.

=== "macOS with Homebrew"

    ```sh
    brew trust spacedock-dev/tap
    brew install --cask spacedock-dev/tap/spacedock@next
    ```

=== "Linux or macOS without Homebrew"

    ```sh
    curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | SPACEDOCK_CHANNEL=edge sh
    export PATH="${SPACEDOCK_INSTALL_DIR:-$HOME/.local/bin}:$PATH"
    ```

## 2. Set up your agent

Choose the agent you use with Spacedock.

=== "Claude Code"

    ```sh
    spacedock install --host claude
    ```

    In Claude Code, run `/reload-plugins` or start a new session.

=== "Codex"

    ```sh
    spacedock install --host codex
    ```

    Start a new Codex session to load the installed skills.

## 3. Confirm setup

```sh
spacedock doctor
```

Check the report for missing programs or plugins. If it reports missing programs, repeat step 1. For missing plugins, repeat step 2.

Setup is complete when the required programs and plugins are available. To use Review, continue to [Your first review](../review/use.md).

## Troubleshooting

### Your agent cannot see the skills

In Claude Code, run `/reload-plugins`. In Codex, start a new session. If the skills are still missing, repeat step 2.
