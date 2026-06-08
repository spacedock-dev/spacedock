# Commander dispatch — sprint 0198-pre-flip-hardening

You are the **Commander** for `0198-pre-flip-hardening`. Drive its members to 0.19.8,
approve execution gates and merge with judgment, run the integration test, see 0.19.8 cut,
and report. On the escalation triggers below you escalate — you do not decide.

## Cold boot

1. `cd /Users/clkao/git/spacedock-research/spacedock-v1`
2. `Skill(skill="spacedock:first-officer")` — load the contract + Claude runtime, run Startup.
3. `git fetch origin next && go build -o ./spacedock ./cmd/spacedock`. **Do NOT `git reset --hard` in a shared working tree** — use fetch + build, or your own checkout/worktree.
4. `git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev`.
5. Create your OWN team. **If you were spawned as a subagent you cannot dispatch sub-agents — STOP and hand back; a Commander needs a separate top-level session** (proven during the 019x dry run). Report your mode first.

## Your sprint

Goal / DoD / deliverable: see [`index.md`](index.md). **Deliverable:** spacedock 0.19.8 cut on `next`.

Members (the query is the source of truth):
```bash
./spacedock status --workflow-dir docs/dev --where sprint=0198-pre-flip-hardening --where 'sprint-readiness != defer' --json --fields id,slug,status,group
```

## Drive split

- **You drive (through merge):** `nzb7` (e2e-gate), `69rk` (codex-cwd), `1p27` (scaffold-fact),
  `kb` (migration-check), `4t` (agentsview-detect). `kb`/`1p27` are well-specified →
  straight to implementation (the `78` precedent). Run ideation for `4t` if undone.
  **Coherence:** the survey members (`69rk`/`1p27`/`4t`) must share ONE corrected agentsview
  model — `69rk` found agentsview keys `project` by git-root basename, not cwd basename.
- **Captain-gated — do NOT decide, run the ceremony once ruled:**
  - `qa` (headline) — its design sits at the captain's ideation gate; drive once approved. Its
    proof requires a **live drive** (the binary/version UX observed), not prose.
  - `jh` — **BLOCKED on a captain architectural decision.** Its ideation recommends a
    `CONTRACT_VERSION` bump (1→2). Do NOT implement a contract bump without explicit captain
    approval — it is breaking and couples to the flip. Hold `jh` until the captain rules.
  - The **0.19.8 release cut** (version bump + tag + push) — outward-facing; the captain confirms.
- **Off-limits:** the 0200-flip (`pj`), the Codex peers (`27`/`z6`), `5h0` (blocked on #315).

## Definition of Done

1. Every `ready` member `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green with `.spacedock-state` + `.claude/worktrees` present.
3. `qa` proven by a live drive.
4. spacedock 0.19.8 stamped + cut on `next`.

## Authority & escalation

Approve execution gates (ideation / implementation / validation) and merge for YOUR sprint
members with judgment. Escalate to the captain / Shaping FO ONLY on: a 3rd feedback cycle,
a budget blowout, an irrecoverable block, a genuine scope fork, or the captain-gated items
above (the `jh` contract decision, the `qa` gate, the 0.19.8 release cut).

## Report

When the DoD holds (or you are blocked): per-member outcome, the integration-test result,
the 0.19.8 cut status, friction (log to `docs/dev/.spacedock-state/fo-friction-log.md`), and
any deferred / escalated items.
