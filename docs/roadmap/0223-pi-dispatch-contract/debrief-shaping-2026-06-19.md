---
session-date: 2026-06-19
sequence: 1
first-commit: 96a91f65
last-commit: 5bf98e1e
duration: ~7h (2026-06-19 ~10:00 → ~17:00 PDT, with async gaps)
runtime: pi (z-ai/glm-5.2, thinking:medium)
role: Shaping FO (0223-pi-dispatch-contract)
---

# Session Debrief — 2026-06-19 (0223 Shaping FO, pi)

Shaped sprint `0223-pi-dispatch-contract` end-to-end on pi and handed it off to a live Commander session via intercom — the first sprint shaped entirely on the pi runtime. The session's real work was **two scope reversals driven by captain interrogation**, each of which retired a clone-bound workaround in favor of the install-managed mechanism, plus a root-cause investigation that turned a "Commander corrupted state" report into a status-tool determinism bug. The contract's `code-gate-over-prose` principle was applied to every mechanism claim: the symlink, the cwd-fallback, the install-managed package, and the validate non-determinism were all verified against live code, not asserted.

## Shipped (shaping artifacts, main)
- **`0223-pi-dispatch-contract` sprint package** — `index.md` + `dispatch-sprint-execution.md` + `staff-review.md` + `staff-review-2.md` (`docs/roadmap/0223-pi-dispatch-contract/`). 4 members, DoD, Q1–Q14 cold-boot quirks, two independent staff reviews (both gaps-closed, no material redesign). Commits `96a91f65` → `5bf98e1e`.
- **Pre-commit hook** installed at `.git/hooks/pre-commit` (worktree-shared): state-checkout-aware, runs `spacedock status --validate`, blocks on hard errors. Documents its own limitation (can't reliably block blank-id untracked entities — see Filed).

## Filed (state, backlog)
- **`pi-back-channel-dispatch`** (`b23…`) — capstone: harden `fo-dispatch-core.md` to 7 runtime-neutral named capabilities; wire Pi adapter back-channel over intercom. Live spike evidence preserved (run `0637e2ed`). Gate-approved; Commander driving (implementation).
- **`pi-dispatch-model-stamping`** (`bdt…`) — null model → parent live model via `intercom list`. Gate-approved; **Commander delivered (#405 merged)**.
- **`pi-install-managed-skill-placement`** (`eq…`) — ship Spacedock as a pi package; `spacedock install --host pi` actually installs. Gate-approved; **Commander delivered (#406 merged)**.
- **`pi-safehouse-flag-parity`** (`qn…`) — pi front-door `--safehouse-*` flags + safehouse wrapping (parity with claude/codex). Gate-approved; **Commander delivered (archived PASSED)**.
- **`pi-ensign-skill-injection`** (`k8t…`) — ARCHIVED REJECTED. Superseded by `eq`; the `.pi/skills/ensign` symlink was clone-bound.
- **`pi-launcher-repo-resolution`** (`2m1…`) — ARCHIVED REJECTED. Superseded by `eq`; the cwd-fallback demotion + install-record file was clone-bound.
- **`status-validate-determinism`** (`a1a…`) — root-cause fix: `spacedock status --validate` is non-deterministic on untracked entity files (the same blank-id file flips between Error and VALID). The load-bearing fix that makes the pre-commit hook reliable.
- **`state-init-precommit-hook`** (`mv9…`) — `spacedock state init` auto-installs the managed pre-commit hook. Depends on the determinism fix.

## Decisions (captain)
- **Scope reversal #1 — install-managed package over symlink + cwd-fallback.** Captain interrogation ("what happens when the user does `spacedock install --host pi`? where does it symlink from?") exposed that members 1+2's mechanisms were clone-bound workarounds for install being check-only. Verified against `obra/superpowers` (`.pi/extensions` + `package.json pi.skills`) and pi-subagents source (`collectSettingsPackageSkillPaths`). Merged 1+2 into `pi-install-managed-skill-placement`; archived `k8t`+`2m1` REJECTED.
- **Scope reversal #2 — `pi-safehouse-flag-parity` as 4th member.** Captain flagged `--safehouse-add-dirs` unavailable on `spacedock pi`. Verified the pi front door registers only `--plugin-dir` (no safehouse flags, no wrapping) while claude/codex wrap via `frontdoor.go:310`. Filed as 4th member; it unblocks member 1's AC-1 on sandboxed sessions (Q14).
- **Drive vs. package.** Captain directed "package commander prompt rather than driving implementation" — Shaping FO hands off, Commander drives. Established intercom as the live handoff channel; the cold-boot file is the durable artifact.
- **Defer `SPACEDOCK_BIN`-through-sandbox forwarding** (Decision 2 in `qn`'s ideation) — pi never forwarded it today; minimal wrap introduces no regression. Concurred with the ensign's deferral.
- **Two staff reviews, both gaps-closed.** Review #1 found the child-cwd blocker; review #2 found 7 follow-through gaps from the re-carve. All closed without material redesign.

## Issues — Workflow
- **`status --validate` non-determinism** (filed `status-validate-determinism`) — the root cause of the escaped blank-id. A blank-id entity file on disk flips between Error and VALID across validate invocations. This is why the Commander's hand-filed `pi-devoverride-package-ok` escaped every gate. The pre-commit hook can't reliably block it until this is fixed.
- **Commander bypassed `spacedock new`** to file `pi-devoverride-package-ok` (hand-wrote the file, blank id from birth). Procedural gap; the pre-commit hook + `state init` auto-install is the defense, but only reliable once the determinism bug is fixed.
- **`pi-devoverride-package-ok` archived with `status: implementation`** (not `done`) — secondary latent inconsistency; not boot-blocking. The id was stamped upstream (boot passes), but the status field is still wrong.

## Issues — Runtime (pi-specific, surfaced this session)
- **Foreground dispatch can't service mid-run `need_decision`** (Q1). A blocking `subagent(...)` call occupies the FO thread; the worker's 10-min intercom timeout fires before the FO regains control. Run `b929622e` timed out foreground; run `0637e2ed`'s `need_decision` arrived after teardown. **Async dispatch is mandatory** — and it's friction 2 in the capstone.
- **Null model resolves to `settings.json` defaultModel, not parent's live model** (Q2). `dispatch build` emits `model:null`; pi-subagents resolves null to `~openai/gpt-mini-latest`, not the FO's `z-ai/glm-5.2`. Run `0637e2ed` ran on gpt-mini. Must pass `model:` explicitly until `pi-dispatch-model-stamping` lands (Commander delivered #405).
- **Skill injection broken until `eq` lands** (Q3). `skill:["ensign"]` → "Skills not found: ensign"; pi-subagents children don't inherit the parent's `--skill` flags and the repo's `skills/` isn't in the child's `buildSkillPaths`. Skill-less workers default to implementation behavior even in ideation (run `0637e2ed` edited contract docs during ideation — contained, reverted). The install-managed package (`eq`, #406 merged) fixes this at the root.
- **State branch is `spacedock-state/dev`, not `main`** (Q4). Pulling `main` into the state checkout rebases unrelated upstream and conflicts ("README.md had different types"). Operator error at boot; corrected.
- **glm-5.2 thinking:high pauses trip the needs-attention threshold** (~60s no activity). Both staff reviewers were flagged mid-run and completed fine. The control signal is conservative; do not interrupt on brief silence.

## Async ensign misroutes
- Two `worker` dispatches misread their role and did FO-side work instead of their assigned task: one narrated as the FO ("I'm the FO (parent), not the subagent") and did my inline index fixes (gaps 2–5,7) — which was correct work but not its assignment; one echoed FO prose and no-op'd the `b2` gap-6 fold-in. Both re-dispatched with tighter, instruction-style prompts and completed. **Lesson: pi dispatch prompts for ensigns must be imperative and explicit about the single deliverable, not prose-heavy context the worker can misread as its role.**

## What worked
- **Captain interrogation as the primary de-risking mechanism.** Two sharp questions ("where does it symlink from?", "`--safehouse-add-dirs` unavailable") each retired a wrong mechanism and reshaped the sprint. The Shaping FO's job was to verify each against live code and re-shape; the captain's questions did the load-bearing falsification.
- **Two independent staff reviews.** Review #1 caught the child-cwd blocker; review #2 caught 7 re-carve follow-through gaps. Neither was a rubber stamp; both verified claims against live source (`schemas.ts:69`, `pi.go:293`).
- **Intercom as the live Commander handoff.** The Commander was already mid-drive when I finished the 4th member; a non-blocking `send` delivered the scope change + the AC-1 hold call in-channel, complementing the durable cold-boot file. The planner-worker pattern from `pi-intercom` worked as advertised.
- **Append-only fold-ins.** Every staff-review gap was folded into entity bodies append-only, preserving prior ideation + spike evidence. No rewrites; SUPERSEDED pointers reconcile stale claims.

## What didn't
- **The pre-commit hook can't reliably block blank-id entities.** Installed and useful for deterministic catches, but the status-tool non-determinism (filed) undermines its load-bearing purpose until fixed.
- **My first model-self-check was wrong.** I claimed the FO session was `gpt-5.2` from a per-file jsonl grep; the captain corrected from the pi status bar (`z-ai/glm-5.2`). The status bar is the live truth; jsonl parsing across session files misleads. Logged.
- **Foreground dispatches early in the session.** I dispatched the first ideation foreground and hit the Q1 timeout/teardown failure. Switched to async for all subsequent dispatches. The friction is real and is exactly what the capstone ships.

## Handoff state
- **Commander (session `a5a1b04f`)** cold-booted on pi, driving. Members 1+2 merged (#405, #406); member 3 archived PASSED; member 4 (capstone) in implementation. No escalations. Shaping FO available for close-out (seed next sprint with deferred frictions 7–9 + pre-cut audit findings) and the two filed followups (`status-validate-determinism`, `state-init-precommit-hook`) when shaped.
