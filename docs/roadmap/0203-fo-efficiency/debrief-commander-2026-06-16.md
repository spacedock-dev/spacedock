---
session-date: 2026-06-16
sequence: 1
first-commit: 6c464b1a
last-commit: 9bd1f46a
duration: ~interactive Commander session
---

# Session Debrief — 2026-06-16 (Commander)

Headline: **2y merged (#385) and v0.20.3 CUT + published.** The 2y "live-suite flake" that drove five prior CI cycles was finally root-caused — and the handoff's premise was wrong. It was never a prose-fixable team-mode coin: an upstream Claude regression (2.1.178 dropped headless team tools, **anthropics/claude-code#68721**) layered on a structural SDK-headless limitation (**anthropics/claude-code-action#1124**), which 7e had already ruled "go bare under `-p`." Aligned the suite with that determination, shipped 2y, filed the deferred live-team-e2e gap, and cut the release.

## Shipped
- **2yf** `shared-merge-dispatch-contract` — [#385](https://github.com/spacedock-dev/spacedock/pull/385). The shared first-officer merge/dispatch ceremony was siloed in Claude-named references, so codex/pi FOs reaching terminalization had no contract and would break. Extracted two host-neutral cores (`fo-merge-core.md` + `fo-dispatch-core.md`) named once by the boot-resident core; each runtime adapter keeps only its host-specific part. This session also folded in the CI alignment that got it green: reverted the 2.1.177 pin (CI tests latest, headless bare), retired the two non-viable forced-team `-p` tests, and fixed a deterministic helper bug (`spawn-standing-all` emitting Agent specs without `description`). Verdict PASSED (validation cycle 4 confirmed the extraction byte-unchanged since the cycle-2 PASS; all 8 PR checks green; live run 27596254379 all four legs).

## Released
- **v0.20.3** cut + published (release run 27600604617 green: e2e-gate + goreleaser + journey-ledger). `main` stamped to 0.20.3 (`plugin.json`), `stable` branch advanced, homebrew cask bumped, 8 platform tarballs published. Notes were hand-corrected to lead with the 2y terminal-contract bullet — `spacedock-release notes` had dropped #385 (it missed the post-flip main merge).

## Filed (backlog)
- **m40** `live-team-mode-terminal-harness` (id `m40mphxan8phr3t3tp03gk89`) — build an interactive/pty live harness for team-mode e2e (comm-officer roster injection + bounded teardown marker). The deferred gap from go-bare: the entire live suite is `claude -p` (headless), where team mode is unsupportable; team mechanisms stay unit/fixture-covered, only live team e2e needs the harness.

## Decisions (captain)
- **Pin first, then pivot to go-bare.** Tactically pinned the claude-live CI to 2.1.177 (restored team tools) — residency went green, proving the pin worked, but the suite stayed red as mirror-flakes (each team test passed one leg, failed the other). On confirming the flakiness is structural (#1124) and that 7e already recorded "go bare if team mode too fragile under `-p`," switched strategy: **reverted the pin, ran latest+bare, retired the two forced-team tests.**
- **"We tackled this before" was the pivot.** Caught the FO before it built a headless team-await keep-alive fix — which would have cut against 7e's (#381) go-bare determination and #271's deferred-await finding. Aligned with the prior determination instead of re-litigating it.
- **Cut v0.20.3 with no fresh live-e2e** → the `e2e-gate` was waived via `SPACEDOCK_E2E_GATE_WAIVER` (auditable reason recorded; `6c464b1a` is the squash of the fully-green `f1be9f07`). One-shot: the waiver var was removed immediately after the gate cleared.
- **Bug 1 fixed regardless of mode** — the missing `description` broke comm-officer injection in interactive team mode too, so it ships independent of the go-bare disposition.

## Issues — Spacedock
- **Headless team mode is unsupportable on current Claude (root-caused, empirically).** Claude Code **2.1.178** removed the native `TeamCreate`/`TeamDelete` tools from headless `-p`/`sdk-cli` sessions (present in **2.1.177**) — anthropics/claude-code#68721 (filed). Confirmed by a controlled version probe holding every CI env var constant (flag=1, sdk-cli, child-session, fresh config): 2.1.178 → 28 tools no Team*; 2.1.177 → 30 tools with both. Even with tools present, the SDK headless session races to `end_turn` before teammates finish (anthropics/claude-code-action#1124). **Implication: headless `-p` FOs run BARE; team mode is interactive-only.** The three prior `forceTeamModeCue` prose escalations could never work — the tool was absent, not declined.
- **`spacedock dispatch spawn-standing-all` emitted Agent specs without `description`** — the Agent tool requires it, so forwarding the spec verbatim failed `InputValidationError` and the comm-officer never spawned (also broken in interactive team mode). `dispatch build` already includes one. **FIXED this session** (added `Description` to `spawnSpec` + a red→green test). Candidate to note: the two helpers should share spec construction so they can't drift again.
- **Bare-mode completion-signal noise:** every bare-mode ensign tried `SendMessage(to="team-lead")` for its completion, found no addressable team-lead, and fell back to its final message. Cosmetic (the FO reads the final message as the artifact) but noisy across four dispatches. Candidate: the ensign completion contract should detect bare mode and skip the SendMessage.

## Issues — Workflow / live suite
- **The five-cycle "flake" was a misdiagnosis.** The handoff framed it as the FO generating a fake team name / single-implicit-team, fixable by stronger prose. The CI transcripts show the FO behaved correctly: it searched for / called `TeamCreate`, got "No such tool available", and fell to documented bare mode. The two team-forced tests (residency #375, teardown #381) were born this sprint and never reliably green — they encoded a headless-team assumption the vendor doesn't support.
- Team-mode coverage that survives go-bare: `internal/dispatch/spawn_standing_all_test.go` + `standing_parity_test.go` (injection), `internal/ensigncycle/teardown_grade_watcher_test.go` + `testdata/sonnet_teamdelete_*.jsonl` (marker grading). Only the live e2e is deferred (→ m40).

## Observations
- **Local-first validation (captain steer):** when uncertain, run the relevant live tests locally on one model FIRST; reserve the full CI matrix (sonnet+opus+codex+pi) for confirmation. This session over-relied on CI cycles, deferring to the handoff's 429 caution instead of just trying locally first.
- **Spike discipline paid off:** treating the handoff's "verified" diagnosis as a hypothesis and re-proving on the actual target (the version probe) overturned three sessions of misdirected prose fixes. Docs-confidence (the web research said team mode is "structurally interactive-only") was contradicted by the empirical probe (2.1.177 exposes the tools headless) — only the empirical method found the real boundary (the 2.1.177→2.1.178 version line).
- The prior `re` follow-up (`team-mode-verdict-omission`) is likely subsumed by the headless=bare reality — verdict-omission was a team-mode-under-`-p` symptom; headless now runs bare.

## What's Next
- **Recommended:** **m40** `live-team-mode-terminal-harness` — the live team-mode e2e gap. Needs a pty harness; sizeable, its own project.
- **Dispatchable now** (implementation→validation): **0q** `status-source-not-default`, **48** `commission-templates-defer-to-contract`, **6rt** `ci-log-read-summary`.
- **Reconcile notes (report-only, peer-owned):** worktrees `agy-runtime-support` (behind origin/main by 79) and `status-section-reader` flagged stale-branch — not this session's to act on.
- v0.20.3 is OUT; the next cut re-requires a real green e2e-gate (waiver removed).
