---
title: "Installation"
description: "Your agents are waiting on you. Spacedock brings you the calls that need you, keeps your agents on the goal, and lets you decide in one place."
doc_version: "0.20.2"
last_updated: "2026-11-06 00:33:04"
---

# Installation

Spacedock works with Claude Code and Codex.

## Install Spacedock

`brew trust spacedock-dev/tap brew install --cask spacedock-dev/tap/spacedock@next` copy

`curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | SPACEDOCK_CHANNEL=edge sh export PATH="${SPACEDOCK_INSTALL_DIR:-$HOME/.local/bin}:$PATH"` copy

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
