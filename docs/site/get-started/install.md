# Installation

Spacedock works with most coding agents.

<a id="1-install-spacedock"></a>

## Install Spacedock

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

<a id="2-set-up-your-agent"></a>

## Set up your agent

=== "Claude Code"

    <div class="sd-agent-prompt sd-command-prompt">
      <code>spacedock install --host claude</code>
      <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
    </div>

    Run `/reload-plugins` or start a new session.

=== "Codex"

    <div class="sd-agent-prompt sd-command-prompt">
      <code>spacedock install --host codex</code>
      <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
    </div>

    Start a new session.

<a id="3-confirm-setup"></a>

## Confirm setup

<div class="sd-agent-prompt sd-command-prompt">
  <code>spacedock doctor</code>
  <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
</div>

Follow the instructions in the doctor output.

[Your first review](../review/use.md).
