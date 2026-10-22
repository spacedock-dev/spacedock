# Installation

Spacedock should work with most coding agents. This guide covers manual setup for Claude Code and Codex.

## 1. Install Spacedock

Choose your operating system. This installs the programs on your machine.

=== "macOS with Homebrew"

    <div class="sd-agent-prompt sd-command-prompt">
      <code>brew trust spacedock-dev/tap&#10;brew install --cask spacedock-dev/tap/spacedock@next</code>
      <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
    </div>

=== "Linux or macOS without Homebrew"

    <div class="sd-agent-prompt sd-command-prompt">
      <code>curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | SPACEDOCK_CHANNEL=edge sh&#10;export PATH="${SPACEDOCK_INSTALL_DIR:-$HOME/.local/bin}:$PATH"</code>
      <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
    </div>

## 2. Set up your agent

Choose the agent you use with Spacedock.

=== "Claude Code"

    <div class="sd-agent-prompt sd-command-prompt">
      <code>spacedock install --host claude</code>
      <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
    </div>

    In Claude Code, run `/reload-plugins` or start a new session.

=== "Codex"

    <div class="sd-agent-prompt sd-command-prompt">
      <code>spacedock install --host codex</code>
      <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
    </div>

    Start a new Codex session to load the installed skills.

## 3. Confirm setup

<div class="sd-agent-prompt sd-command-prompt">
  <code>spacedock doctor</code>
  <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
</div>

Check the report for missing programs or plugins. If it reports missing programs, repeat step 1. For missing plugins, repeat step 2.

Setup is complete when the required programs and plugins are available. To use Review, continue to [Your first review](../review/use.md).

## Troubleshooting

### Your agent cannot see the skills

In Claude Code, run `/reload-plugins`. In Codex, start a new session. If the skills are still missing, repeat step 2.
