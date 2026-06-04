---
id: fcfe02dk6nbwk9pw896d3eff
title: Front-door launcher doesn't propagate its own binary path ($0) to the launched claude/codex session — the in-session skill falls back to $PATH spacedock instead of the launching binary
status: ideation
source: captain (2026-06-03) — `/path/to/spacedock claude …` does not preserve `$0` for the skill to know which spacedock launched it; need a way to pass env from the spacedock CLI over (safehouse permitting) to claude/codex, and have the skill use the env or default to $PATH
score: "0.28"
worktree:
started: 2026-06-04T00:00:00Z
completed:
verdict:
issue:
---

When `spacedock claude` (or `spacedock codex`) launches the agent CLI with the FO/ensign skills, the operating contract inside that session shells out to `spacedock` (e.g. `spacedock status`, `spacedock dispatch build`, `spacedock dispatch reconcile`). Those invocations resolve `spacedock` via **`$PATH`** — not the binary that actually launched the session (`$0`). When the launching binary is a dev build (`/path/to/spacedock`, or `claude --plugin-dir` dev mode, or a freshly-built `./spacedock`), the in-session skill silently runs a **different** binary on `$PATH` (often the installed one, possibly an older contract).

## Problem

The dev workflow is "launch with the new binary to exercise it" — but the in-session skill doesn't know which binary launched it, so it can't reuse it. The launcher's own path (`os.Args[0]` / the resolved executable) is not propagated into the child session's environment, and the skill has no signal to prefer it over `$PATH`.

Concrete instance (session 11): this FO session booted with the `$PATH` `spacedock` at **0.19.1** while the repo build was newer; the FO had to invoke `./spacedock` explicitly throughout. A dispatched ensign has no such escape hatch — it follows the contract's bare `spacedock` invocations and gets whatever `$PATH` resolves, which may be a stale contract version. This is the binary-identity half of the already-noted "dispatched ensigns load the INSTALLED plugin contract, not the repo's vendored copy" friction.

## Proposed approach

1. **Launcher exports its own resolved path.** `spacedock claude` / `spacedock codex` set an env var (e.g. `SPACEDOCK_BIN`) to the launching binary's absolute, symlink-resolved path in the child process environment.
2. **Safehouse passthrough.** The wrap goes through safehouse, which sandboxes/strips env — so the var must be on safehouse's allowed-passthrough set (or the launcher adds it). The design must state how the var survives (or is re-injected after) the safehouse wrap, and degrade gracefully when safehouse forbids it.
3. **Skill/contract resolves env-then-PATH.** The operating contract (and/or the dispatch helpers it calls) use `$SPACEDOCK_BIN` when set and non-empty, defaulting to `spacedock` on `$PATH` otherwise. Prefer a code-level resolution (a helper that reads the env) over a prose instruction, so the guarantee is enforced rather than merely documented.

Ideation decides the altitude (env var name, whether resolution lives in the binary's own re-exec helpers vs the skill prose vs both) and exercises the riskiest unknown first: that the env var actually survives the safehouse wrap into the child claude/codex process.

## Out of scope

- Changing safehouse itself; this works within safehouse's passthrough policy.
- A general env-passthrough framework beyond what this one signal needs (unless the design naturally generalizes).
- The separate "in-session skill loads the installed plugin contract, not the repo's vendored copy" problem — related but distinct (contract source vs binary path); coordinate, do not merge.

## Acceptance criteria

**AC-1 — The launcher injects its own resolved binary path into the child environment.**
Verified by: a Go test on the launcher's child-env construction asserting the env contains `SPACEDOCK_BIN` set to the resolved absolute path of the launching executable, for both `spacedock claude` and `spacedock codex`.

**AC-2 — The skill/contract prefers the injected path and falls back to $PATH.**
Verified by: a test (or the resolving helper's unit test) showing `spacedock` invocations resolve to `$SPACEDOCK_BIN` when set and to `$PATH` `spacedock` when unset — proof at the resolution code's level, not a prose-only assertion.

**AC-3 — The signal survives (or degrades cleanly under) the safehouse wrap.**
Verified by: a check that the var reaches the child process through safehouse when passthrough is permitted, and that absence of the var falls back to `$PATH` without error when safehouse strips it. Runtime-observable — confirm against a real safehouse-wrapped launch (by-construction-pending-live until then).

## Test plan

- Go unit tests for the launcher env construction (AC-1) and the env-then-PATH resolver (AC-2). Cost: low.
- A live (or safehouse-fixture) check that the var survives the wrap (AC-3). Cost: medium; exercise this riskiest unknown first in ideation.
- High-stakes surface (front-door launcher) → detached adversarial audit before merge.

## Notes

- Front-door launcher entity, not test-themed — filed for the captain to slot (candidate for 0.19.5 or a launcher-focused sprint).
- Pairs conceptually with the launch-parity `--plugin-dir` dev-mode work in sprint notes (2026-05-31): both are about making the dev launch self-consistent (right plugin, right binary).
