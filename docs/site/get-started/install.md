# Installation

Spacedock works with Claude Code and Codex.

<a id="1-install-spacedock"></a>

## Install Spacedock

=== "macOS with Homebrew"

    <div class="sd-agent-prompt sd-command-prompt">
      <code>brew trust spacedock-dev/tap&#10;brew install --cask spacedock-dev/tap/spacedock</code>
      <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
    </div>

=== "Linux or macOS without Homebrew"

    <div class="sd-agent-prompt sd-command-prompt">
      <code>curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh&#10;export PATH="${SPACEDOCK_INSTALL_DIR:-$HOME/.local/bin}:$PATH"</code>
      <button type="button" class="sd-prompt-copy" aria-label="Copy command" aria-live="polite" hidden>copy</button>
    </div>

The installer looks up releases at `spacedock.md` with User-Agent `spacedock-install/<version|unknown> (<os>; <arch>[; ci][; internal])`.
The version identifies the installed Spacedock binary. If the version probe fails or takes more than two seconds, the installer uses `unknown`.
The installer uses GitHub directly if the route fails or takes more than five seconds.

PostHog counts each lookup without the caller's IP address. Plausible uses the caller's IP address to count the lookup but does not store it.
CI lookups are not counted. The installer adds `; ci` when `CI` is `1`, `true`, `yes`, or `on`, regardless of case or surrounding whitespace.
Set `SPACEDOCK_INTERNAL=1` to add `; internal`.

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

In Codex, run `spacedock doctor --host codex`. Follow the instructions in the doctor output.

## Keep it up to date

Run `spacedock update`. It updates Spacedock, Review, Compass and their plugins together. A Homebrew install prints the `brew upgrade` command to run instead.

[Your first review](../review/use.md).
