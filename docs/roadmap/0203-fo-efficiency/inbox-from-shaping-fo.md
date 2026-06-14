# 0203 — Inbox: shaping FO → Commander

> Best-effort FO↔Commander notifications. The live channel is the `xp` gap, so this file is a durable inbox the Commander reads at boot/pull. The **source of truth for membership is always** `spacedock status --workflow-dir docs/dev --where sprint=0203-fo-efficiency` — not this note, not the index.

## 2026-06-14 — two tasks folded into 0203 (captain)

The captain folded two fresh backlog seeds into `sprint: 0203-fo-efficiency`. **Both are unshaped — they need ideation; neither is ready to execute.**

- **`lean-boot-hardening`** (backlog) — Startup discovery must report-and-stop when `spacedock status --discover` returns zero, never fall back to a broad `find`/`grep` filesystem sweep (a contract + lean-boot violation, observed live). Extends j9's shallow-boot ethos. Its proof must be behavioral/code (a live drive of the report-and-stop, or a code gate), never a contract prose-grep.
- **`ci-log-hygiene`** (backlog) — the live shared-scenario runner (`internal/ensigncycle/claude_live_runner_test.go:365`) `t.Logf`s the full host jsonl stream to CI stdout, bloating logs (~143KB on one failed step) and burying the failure. The jsonl is already a per-scenario artifact, so make it artifact-only / dump-on-failure. Surfaced while debugging #368's opus `gate-guardrail` no-progress flake.

Both now appear in the sprint query above. Pick them up with the next wave at your discretion; no readiness flag is set.
