# Commander dispatch — sprint 0198-pre-flip-hardening

You are the **Commander** for `0198-pre-flip-hardening`. Drive its members to 0.19.8,
approve execution gates and merge with judgment, run the integration test, see 0.19.8 cut,
report. On the escalation triggers below you escalate — you do not decide.

## Preflight gate (must hold before you drive)

A two-lens staff readiness review ran (see [`staff-review.md`](staff-review.md)). Verdict:
GO **except** one open blocker —

- **B2 (survey model contradiction) — OPEN, captain-decided.** Do NOT drive any survey
  member (`69`/`1p`/`4t`) until the captain resolves whether `SKILL.md:64` + the
  `queries.sql` `scoping`/`folded_keys` rationale are corrected to agentsview's real
  **git-root-basename** `project` keying (`69`'s spike disproved the per-cwd-basename premise
  `xn` shipped) or deferred with a tracked follow-up. The other groups are clear to drive.
- **B1 (DoD#4 live drive) — CLOSED:** DoD#4 is now a captain-run sprint-acceptance live
  drive (index.md), not a member AC.

## Cold boot

1. `cd /Users/clkao/git/spacedock-research/spacedock-v1`
2. `Skill(skill="spacedock:first-officer")` — load the contract + Claude runtime, run Startup.
3. `git fetch origin next && go build -o ./spacedock ./cmd/spacedock`. **Do NOT `git reset --hard` in a shared tree.**
4. `git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev`.
5. Create your OWN team. **If you were spawned as a subagent you cannot dispatch sub-agents — STOP and hand back; a Commander needs a separate top-level session** (proven in the 019x dry run). Report your mode first.

## Your sprint

Goal / DoD: [`index.md`](index.md). **Deliverable:** spacedock 0.19.8 cut on `next`.
Members (the query is the source of truth):
```bash
./spacedock status --workflow-dir docs/dev --where sprint=0198-pre-flip-hardening --where 'sprint-readiness != defer' --json --fields id,slug,status,group
```

## Drive plan (ordering + per-member, from the preflight)

**Binary-UX:**
- `qa` (HEADLINE) — captain-gated ideation gate; drive once approved. Mechanism is sound;
  DoD#4 is the captain's sprint-acceptance live drive.
- `z9` (codex auto-install) — **drive AFTER qa** (both edit `frontdoor.go`+`host_exec.go`,
  disjoint functions, sequence to avoid textual overlap). High-stakes front-door →
  **validation gets a detached adversarial audit.** Build note: **fix** the now-false
  comments/error-strings it builds around (`host_exec.go:271-273`, `:32-34`;
  `frontdoor.go:314-316`), don't add around them. Inverts `frontdoor_test.go:414`.

**Survey (BLOCKED on B2 — resolve first):** land `4t` (line 27, isolated) → resolve B2 →
`69` (queries.sql + step-2/step-4 hint) → `1p` (step-3 + step-4 SCAFFOLD). `69`+`1p` both
edit the step-4 report fence (~141-169): land sequentially or hand-merge. Build note for
`69`: agentsview keys `project` by git-root basename — don't assume per-cwd keys. `1p` is
proof-bar-light (live-drive only) and thinner than the kb bar — firm its exact SKILL.md edit
first (drop the taxonomy, keep the "recovered from behavior" fact).

**Independent:**
- `kb` (migration-check + orphaned-fixture delete) — drive-now, skip-ideation-ready. Build
  note: add a POSITIVE assertion that `.spacedock-state` is pruned, and verify the
  `checked == 0` Fatal guard stays non-vacuous after the prune. Land before the survey
  live-drives (orphaned `scaffolds/` gone). **This is the precondition for DoD#2.**
- `nzb` (e2e-gate) — drive-now, parallelizable; **merge before the captain-gated release.**
  Build note: its true cut-block e2e is the flip's; require its `--dry-run` vs real
  `gh run list` as the not-just-prose bar.

## Captain-gated (do NOT decide)
- `qa` ideation gate; the **B2 survey-model** resolution; the **0.19.8 release cut** (version
  bump + tag + push, outward-facing); DoD#4's captain live drive.

## DoD ownership
1. members done/merged — FO/captain process. 2. `go test ./...` green w/ `.spacedock-state` —
**`kb`**. 3. 0.19.8 cut — FO/captain. 4. qa live drive — **captain** (sprint acceptance).

## Validation budget
`69`/`4t`/`1p` all bottom out on a live drive — exercise all three observable changes
(Codex-presence hint, sandbox probe non-prompt, reshaped SCAFFOLD) in one live-drive pass.

## Off-limits
The 0200-flip (`pj`), the Codex peers (`27`/`z6`), `jh` + `5h0` (deferred).

## Authority & escalation
Approve execution gates + merge for YOUR members with judgment. Escalate on: a 3rd feedback
cycle, a budget blowout, an irrecoverable block, a scope fork, or the captain-gated items.

## Report
Per-member outcome, integration-test result, 0.19.8 cut status, friction
(`docs/dev/.spacedock-state/fo-friction-log.md`), escalations.
