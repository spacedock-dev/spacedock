---
title: Codex resume passthrough must not append the default Spacedock prompt
status: backlog
score: 0.7
source: "Captain report 2026-07-10: `spacedock codex -- resume` should not invoke Codex with the default Spacedock prompt."
id: xvcz44jbmye15bpz1ekzxvkc
---

## Problem

The Spacedock Codex launcher appends its default first-officer bootstrap prompt for a fresh launch. When the operator passes the Codex `resume` subcommand after `--`, that prompt changes the resumed invocation's meaning instead of preserving Codex's resume flow.

## Proposed approach

Classify a passthrough argv whose Codex subcommand is `resume` and do not append the default Spacedock bootstrap prompt. Preserve the operator's passthrough arguments exactly. Keep the default prompt for normal fresh `spacedock codex` launches, and keep safehouse/plugin-dir routing behavior unchanged.

## Out of scope

- Changing Codex's own resume semantics or session storage.
- Removing the first-officer prompt from fresh launches.
- Broad passthrough-command redesign beyond the minimum subcommand classification.

## Acceptance criteria

- **AC-1:** `spacedock codex -- resume` reaches Codex as `codex resume` with no appended Spacedock bootstrap prompt.
- **AC-2:** `spacedock codex -- resume <session>` preserves `<session>` and all later operator arguments in order, with no injected prompt.
- **AC-3:** A fresh `spacedock codex` launch still receives the normal first-officer bootstrap prompt.
- **AC-4:** Safehouse and no-safehouse argv fixtures prove the same resume behavior without leaking launcher flags past `--`.

## Test plan

- Add front-door argv red/green cases for bare resume, resume with a session ID/options, and fresh-launch control.
- Run focused CLI/safehouse tests, full tests, race tests, and formatting checks.
