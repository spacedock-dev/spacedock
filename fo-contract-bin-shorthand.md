---
title: 'FO contract: replace bare `spacedock` command examples with `${SPACEDOCK_BIN:-spacedock}`'
status: ideation
source: 'boot-forensics (2026-06-16) — FO used homebrew cask (0.20.2) instead of $SPACEDOCK_BIN (dev) throughout a session because contract command blocks model the bare-spacedock shorthand. SPACEDOCK_BIN was set and correct but never used; cap divergence cost ~16k tok in failed --read calls. Fix: find-and-replace executable command positions in first-officer-shared-core.md and claude-first-officer-runtime.md.'
score: 0.5
id: 13f8b12x9f7ba25ywm5wt2x7
started: 2026-07-03T03:03:12Z
---

Contract command examples in `first-officer-shared-core.md` and `claude-first-officer-runtime.md` use bare `spacedock` in all executable command positions. The preamble explains the `${SPACEDOCK_BIN:-spacedock}` invariant but the examples undermine it by modeling the shorthand form. When the FO writes Bash calls by copying those examples, it silently uses the wrong binary when `SPACEDOCK_BIN` differs from the PATH `spacedock`.

Surfaced by boot-forensics on 2026-06-16: `$SPACEDOCK_BIN` was the local dev build (has `status --read`); PATH `spacedock` was the homebrew 0.20.2 cask (no `--read`). Both satisfy `contract 1`, so the version gate passed silently. Every `status --read` call fell through to the full entity dump (~8k tok each, two hits = ~16k tok wasted).

Fix: find-and-replace bare `spacedock` in the executable command lines of both files with `${SPACEDOCK_BIN:-spacedock}`. The introductory note already states the rule; the examples need to model it.
