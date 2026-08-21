---
title: "spacedock pi dedupes the pi-subagents extension when it is installed as a package"
status: backlog
source: "GitHub issue spacedock-dev/spacedock#746 — spacedock pi fails at startup when pi-subagents is installed as a package: duplicate extension load ('Tool \"subagent\" conflicts with ...')"
started:
completed:
verdict:
score: 0.8
worktree:
issue: spacedock-dev/spacedock#746
id: 5xwwj9c2w50921t16s840p49
---

`spacedock pi` loads the `pi-subagents` extension twice under two different specifiers when the package is also registered in `~/.pi/agent/settings.json` `packages` — the exact setup `spacedock pi --check` recommends. Package discovery loads `<pkg>/index.ts` (re-exporting `./src/extension/index.ts`), and `spacedock pi` additionally passes `--extension <pkg>/src/extension/index.ts --skill <pkg>/skills/pi-subagents` (internal/cli/pi.go, argv built at the `--extension`/`--skill` block; extensionPath is `filepath.Join(pkg, "src", "extension", "index.ts")`). Pi keys extension identity by resolved specifier, not module identity, so the second registration of `subagent`/`subagent_wait` collides and startup fails with `Tool "subagent" conflicts with ...`. Plain `pi` works; `spacedock pi -- --ne` works only by silencing all discovered extensions (lossy — also drops pi-intercom and host extensions).

Direction (from issue #746): the explicit flags exist for the `--plugin-dir`/`SPACEDOCK_REPO_ROOT` dev-override case (cfg.repoRoot != "", where the extension is NOT registered), per the code comments. Gate them on the package NOT already being registered — skip the explicit `--extension`/`--skill` when pi-subagents is in `~/.pi/agent/settings.json` `packages` (or when extensionPath resolves inside `~/.pi/agent/npm/node_modules`) — preserving the dev-override purpose. Alternative if the explicit flag must stay unconditional: point it at the package's declared entry (`<pkg>/index.ts`) so both load paths produce one specifier (more fragile, depends on pi-subagents' internal layout).

Acceptance sketch: value — with `npm:pi-subagents` in `settings.json` `packages`, `spacedock pi` starts and the subagent tools register exactly once (no `Tool "subagent" conflicts`), and the dev-override path still loads the extension explicitly; mechanism — a behavior test driving the flag-selection path against a registered-vs-unregistered package config asserts one extension load in the registered case and the explicit flag in the dev-override case. Expected surface: `internal/cli/pi.go` plus a test; small diff.
