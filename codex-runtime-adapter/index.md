---
id: r09jrf0k6qjv6c1sddhe1sh6
title: Ship a Codex runtime contract — codex ensign/first-officer adapters + Codex-shaped dispatch + wait/observe semantics
status: ideation
source: "captain (2026-06-02) — live Codex FO session: 0.19.x ships Claude-only runtime adapters; Codex FO/ensign dispatch/wait/observe is improvised (wait_agent timed out 10s then async notifications), and dispatch build emits the Claude-ism Skill(skill=...) which the Codex skill list does not expose"
started: 2026-06-02T15:38:31Z
completed:
verdict:
score: "0.33"
worktree:
issue:
milestone: post-0.19.4 (own track — Codex runtime parity, NOT folded into 0.19.4)
---

0.19.x ships **runtime-agnostic shared cores with Codex-aware mentions** (`send_input` routing, "Codex declares none" budget probe, codex resume) but **no Codex runtime contract**. The result is the "model-improvised, not contract-guaranteed" antipattern at runtime scale — confirmed live by the captain's Codex FO session.

Confirmed gaps in the shipped plugin:
- Only `skills/ensign/references/claude-ensign-runtime.md` and `skills/first-officer/references/claude-first-officer-runtime.md` exist — **no `codex-*-runtime.md`**.
- Both `SKILL.md` files dispatch the runtime adapter on `CLAUDECODE` only (`ensign/SKILL.md:13`, `first-officer/SKILL.md:23`) — there is no `Codex → read codex-runtime.md` branch, so a Codex FO/ensign loads only the shared core and improvises completion/wait/observe.
- `spacedock dispatch build` hard-codes the Claude-ism `Skill(skill="spacedock:ensign")` in the emitted prompt (`internal/dispatch/build.go:300,431`); the Codex session's skill list did not expose `spacedock:ensign`, forcing a read-the-dispatch-file-directly workaround.
- The Codex FO observed: `wait_agent` timed out after 10s, then Codex delivered async subagent-completion notifications — that is **Codex host behavior, not a Spacedock FO contract**.

## Design surface (for ideation to scope)

- **Codex first-officer runtime adapter** (`codex-first-officer-runtime.md`) — the Codex analog of the Claude adapter's `## Awaiting Completion` / dispatch / reuse sections: how a Codex FO observes a completion signal (the wait/notification model), how it routes reuse/feedback (`send_input` per the shared-core mentions), the team-vs-no-team model on Codex, and the budget-probe declaration (shared-core already says "Codex declares none").
- **Codex ensign runtime adapter** (`codex-ensign-runtime.md`) — the Codex analog of `claude-ensign-runtime.md`: completion-signal protocol, worktree ownership, polling.
- **SKILL.md runtime dispatch** — add the Codex branch to both `ensign/SKILL.md` and `first-officer/SKILL.md` so the right adapter loads by platform (Codex env detection).
- **`dispatch build` prompt shape** — emit a Codex-appropriate prompt (not the Claude `Skill(...)` call) when targeting Codex, and resolve the skill-exposure mismatch (either ensure `spacedock:ensign` is exposed on Codex, or emit the read-the-dispatch-file form as the contract rather than a workaround).
- **Reconcile** with the existing codex launcher work (`be`/codex-safehouse-launcher — `spacedock codex` LAUNCHES codex; this entity is about the FO/ensign RUNTIME once launched) and the `spacedock-dev/spacedock` migration. Name the boundary so launcher vs runtime don't double-file.

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — a Codex FO/ensign loads a dedicated runtime adapter, not just the shared core.** Both `SKILL.md` files branch to a Codex adapter on Codex; the adapters exist and define completion/wait/observe + dispatch/reuse for Codex.

**AC-2 — `dispatch build` emits a prompt a Codex agent can act on directly.** The Codex-target prompt does not depend on a `Skill(skill=...)` call the Codex skill list may not expose; verified against the real Codex skill-exposure behavior.

**AC-3 — the wait/observe contract is explicit, not host-incidental.** The Codex completion-signal observation is specified in the adapter (so a Codex FO does not fall back to a raw `wait_agent` timeout + improvised async handling).

## Notes
- OWN TRACK (captain decision 2026-06-02): file + ideate now, but do NOT fold into 0.19.4 (already large). Likely a 0.20 "Codex runtime parity" milestone.
- Ideation-first; this is a design entity (the completion/wait/observe contract is the riskiest unknown — exercise it against a real Codex session before committing the adapter wording).
