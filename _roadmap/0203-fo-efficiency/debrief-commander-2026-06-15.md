---
session-date: 2026-06-15
sequence: 1
first-commit: db414503
last-commit: f2b6f3dc
duration: ~long (full-day interactive Commander session)
---

# Session Debrief — 2026-06-15 (Commander)

Headline: **7e merged (v0.20.3 fold-in #1).** 2y (host-neutral merge/dispatch extraction, fold-in #2) drove through **five CI cycles**, surfacing and fixing the sprint's deep live-suite reliability bugs one by one — now in-flight on the last two (wrong-root-boot + force-team coin). One follow-up filed. The v0.20.3 cut is HELD for 2y.

## Shipped
- **7e** `headless-dispatch-mode-intent` — [#381](https://github.com/spacedock-dev/spacedock/pull/381). The FO's team-vs-bare dispatch-mode choice for headless `-p` was under-specified, so models coin-flipped and `TestLiveEnsignCycle` flaked. Determined + documented the two-mode behavior (interactive greets-and-stops / headless drives-to-gate, keyed on `-p`), dissolved the team-vs-bare coin (proven deterministic on sonnet + opus, `count=3` each), retargeted the live smoke to mode-invariant facts (status:done + path-scoped commit, not the `isTeamCreate` coin), moved teardown-marker coverage to a forced-team scenario, added the AC-a default-`-p`-stops-at-gate scenario, switched the live completion barrier to dispatch-OPEN + on-disk end-state (the team-mode close anchor is stale in Claude Code 2.1.177), refined the AC-1 >60s guard to a name-gated quiet-budget allowlist, and added `-timeout 40m` to the claude-live CI suite. Verdict PASSED (both-model live + detached adversarial audit clean). The team-mode `verdict:`-omission it surfaced was filed as a follow-up (below). Cycle-1 validation REJECTED (the coin was relocated to the teardown grade, not dissolved) → fixed in cycle 2.

## Filed (backlog)
- **re** `team-mode-verdict-omission` — a headless `-p` FO in **team mode** non-deterministically omits the `verdict:` frontmatter value when terminalizing (sets status:done + archives + emits the teardown marker, but ~1/3 of team runs skip the verdict); bare mode writes it reliably. Surfaced by 7e cycle-2. Same team-mode-under-`-p` fragility theme.

## In-flight — 2y (the headline; NOT yet merged)
- **2yf** `shared-merge-dispatch-contract` — at `validation`, `pr=#385` (CI failed), worktree `.worktrees/spacedock-ensign-shared-merge-dispatch-contract`, branch `spacedock-ensign/shared-merge-dispatch-contract`. The host-neutral merge/dispatch extraction: codex/pi FOs reaching terminalization had NO merge ceremony to follow (`shared-core` deferred to a per-host reference only Claude provided) → they'd break. Extracts two host-neutral cores (`fo-merge-core.md` + `fo-dispatch-core.md`) named once by the boot-resident core; each runtime adapter keeps only its host-specific part. **Journey (5 CI cycles):**
  1. First cut PASSED validation but GREW the contract (~+1.5k words; Claude FO loaded ~+965 more) — the host-neutral core-naming was duplicated in every adapter (AC-1's regex forced it) and AC-2 was a transient migration assurance. **Captain sent it back.**
  2. Lean rework: single core-naming, adapters stripped to host-specific parts, AC-2 deleted, AC-1 reframed to reachability, jargon dropped. Hit a floor at **+499 over the 9779 Claude baseline** — the worker proved (section-by-section) it's intrinsic core/adapter split-overhead, not duplication (mostly DEFERRED loads). **Captain: accept the documented overage.** (Worker correctly reverted an id_style cut on finding it was faithful prose, landing at +499 not +434.) PASSED re-validation.
  3. PR #383 → CI **suite-timeout** (7e's line-185 missing `-timeout 40m`) — was 7e's gap; fixed in 7e's PR before merge.
  4. PR #384 → CI failed `TestLiveDefaultHeadlessStopsAtGate` ("FO did not report gate status"). Root cause (a 7e-era gap, not 2y): the headless FO drove to the gate + reported a stop reason (step-9 compliant) but didn't author the full present-gate review the test scans for. **Captain: strengthen the contract** → step-9 headless bullet now MUSTs the full present-gate review at a gate-stop. PASSED re-validation; +125 (Claude total 10403).
  5. PR #385 → the **gate-stop MUST WORKED** (`TestLiveDefaultHeadlessStopsAtGate` passes both legs!) but the suite failed on TWO MORE pre-existing bugs: `TestLiveStandingResidencyInjectsCommOfficer` (the **team-vs-bare coin** — FO drove bare despite force-team, the literal original 0203 flake) and `TestLiveEnsignCycleTeamTeardown` on opus (a **wrong-root-boot** — the FO `cd`'d to the plugin path = the real CI checkout and drove the *real* `docs/dev` workflow for 27 min, never touching the fixture). **Captain: fix the team-mode fragility / make the tests run isolated.** Both fixes DISPATCHED (the worker is applying them).

## Decisions (captain)
- 7e: pre-approved + given the conn; merged on green (both-model live proof is the gate). v0.20.3 placement: 7e rides v0.20.3.
- 2y sent back TWICE: for token bloat (→ lean rework, control Claude delta; accepted +499 intrinsic overage once the duplication was gone) and for the headless gate-stop (→ strengthen the contract, NOT relax the test).
- AC-2 dropped as transient migration-theater (the durable guarantees are AC-1 reachability + AC-3a live).
- Live-suite reliability: chose to FIX the team-mode fragility + the wrong-root-boot (not de-gate, not accept). The wrong-root-boot fix = isolate the plugin/contract path so the real `docs/dev` is undiscoverable.

## Issues — Workflow / live suite (the sprint's core)
- The **claude-live e2e suite is deeply, MULTIPLY flaky** — flaked four distinct ways across 2y's runs: suite-timeout, the gate-stop gap, the wrong-root-boot, and the team-vs-bare residency coin. 2y was the first comprehensive run that exercised the FO live across all these paths; it became the vehicle hardening the live suite.
- **Wrong-root-boot (root-caused):** the live FO anchors project-root discovery on the **plugin path** (`--plugin-dir` = the real repo's `skills/`, where it reads the contract) instead of its isolated cwd; the real repo has a discoverable `docs/dev` sibling, so an FO that `cd`s there escapes into the real workflow. Fix in flight: isolate the plugin/contract path (temp copy w/o `docs/dev` sibling) so it's undiscoverable. Hardens EVERY claude-live test.
- **Force-team coin:** team-REQUIRING tests (residency, team-teardown) still flake because the FO goes bare under `-p` despite the force-team prompt (7e's determination sanctions bare). Fix in flight: a strong reliable force-team cue. (7e made the team-AGNOSTIC tests robust; the team-requiring ones need reliable forcing.)

## Issues — Spacedock
- **Same-name worker supersede race:** re-dispatching a same-name worker (a fresh validator) while the prior instance's shutdown is still resolving races the `teammate_terminated`. The roster `members[]` check is INSUFFICIENT (a pending-shutdown agent is already out of `members[]` but its shutdown still reaps the name). Happened twice this session; survived both by luck (the new instance took the slot). **Fix the FO contract / using-claude-team: wait for `teammate_terminated` before reusing a name, or always use a distinct `-cycleN` name on supersede.** (Candidate for a Spacedock issue.)

## Observations / process scars (for the next Commander)
- **Don't `cd` into the split-root state checkout** — use `git -C docs/dev/.spacedock-state`. A persistent Bash-cwd slip caused a confusing "state vanished" diagnostic detour (it hadn't; the cwd had moved).
- **Token + 429:** the `sk-ant-oat01` benchmark-token rotates server-side (~30–60 min); the **auto-refresh loop `b8h7bwdhj`** keeps `~/.claude/benchmark-token` synced to the keychain every 90s (re-mint cmd: `security find-generic-password -s "Claude Code-credentials" -w | python3 -c 'import sys,json;print(json.load(sys.stdin)["claudeAiOauth"]["accessToken"])' > ~/.claude/benchmark-token`). Separately, the account hit a **sustained 429 rate-limit** (account quota, NOT auth) that local refresh can't fix — it blocks LOCAL live runs all session. CI lanes use their own quota (codex/pi legs passed throughout), so the live behavior gates at the PR CI.
- **CI live-e2e approval:** the `Runtime Live E2E` run gates 3 environments behind deployment approval (CI-E2E / CI-E2E-OPUS / CI-E2E-CODEX). Approve with `gh api -X POST repos/spacedock-dev/spacedock/actions/runs/{LRID}/pending_deployments -F 'environment_ids[]=14103548752' -F 'environment_ids[]=14143978779' -F 'environment_ids[]=14104999836' -f state=approved`. The `gh pr checks --watch` background watcher sometimes exits 1 on a transient `api.github.com` drop — RE-FETCH `gh pr checks {N}` before trusting a "fail".

## What's Next
- **WATCH 2y (the cold-boot handoff `/tmp/0203-2y-watch-handoff.md` covers this).** Worker is applying the wrong-root-boot + force-team fixes → re-validation → re-PR (or push to #385) → CI. The decisive signals on BOTH claude legs: `TestLiveDefaultHeadlessStopsAtGate` PASS (already proven), `TestLiveStandingResidencyInjectsCommOfficer` PASS (force-team fix), `TestLiveEnsignCycleTeamTeardown` PASS (wrong-root-boot fix). Merge on green → then cut.
- **Cut v0.20.3** (7e + 2y). `release.yml` requires the TAGGED main commit to have its own green `Runtime Live E2E` run (squash-merge makes a fresh SHA with no run) — handle at the cut (trigger live-e2e on the final main commit or follow prior-release mechanics), and PAUSE before the tag for the captain.
- Backlog: `team-mode-verdict-omission` (filed); the team-mode-`-p` fragility broadly (the sprint's core).
