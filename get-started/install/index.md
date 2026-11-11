---
title: "Installation"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-11 08:30:05"
---

# Installation

Spacedock works with Claude Code and Codex.

## Install Spacedock

`brew trust spacedock-dev/tap brew install --cask spacedock-dev/tap/spacedock` copy

`curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh export PATH="${SPACEDOCK_INSTALL_DIR:-$HOME/.local/bin}:$PATH"` copy

The installer looks up releases at `spacedock.md` with User-Agent `spacedock-install/<version|unknown> (<os>; <arch>[; ci][; internal])`. The version identifies the installed Spacedock binary. If the version probe fails or takes more than two seconds, the installer uses `unknown`. The installer uses GitHub directly if the route fails or takes more than five seconds.

PostHog counts each lookup without the caller's IP address. Plausible uses the caller's IP address to count the lookup but does not store it. CI lookups are not counted. The installer adds `; ci` when `CI` is `1`, `true`, `yes`, or `on`, regardless of case or surrounding whitespace. Set `SPACEDOCK_INTERNAL=1` to add `; internal`.

## Set up your agent

`spacedock install --host claude` copy

Run `/reload-plugins` or start a new session.

`spacedock install --host codex` copy

Start a new session.

## Confirm setup

`spacedock doctor` copy

In Codex, run `spacedock doctor --host codex`. Follow the instructions in the doctor output.

## Keep it up to date

Run `spacedock update`. It updates Spacedock, Review, Compass and their plugins together. A Homebrew install prints the `brew upgrade` command to run instead.

[Your first review](../../review/use/).

## Sitemap

- [Documentation home](/docs/)
