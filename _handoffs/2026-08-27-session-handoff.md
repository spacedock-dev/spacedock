# Session handoff — 2026-08-27 (pre-compaction)

## In flight

### PR #776 — Pi firstActionBlock loads the ensign skill (THE FIX)
- Branch: `spacedock-ensign/pi-spawn-skill-name-resolution` → base `main`
- Entity: `embed-stage-report-protocol-in-dispatch` (t4) + `pi-spawn-skill-name-resolution` (ntn), both folded, both at `done`/`approved-awaiting-merge`, `pr: #776`
- Content: Pi firstActionBlock invokes `/skill:ensign` before the dispatch file (mirrors Claude Skill() and Codex $spacedock:ensign); embed (stageReportFormatBlock) removed; live lane (TestLivePiNonSelfDescribingDispatch) + CI wiring kept; offline guard TestPiFirstActionInvokesEnsignSkill
- **Locally verified GREEN on lunaroute**: TestLivePiNonSelfDescribingDispatch passes (160s) — a bare-checklist Pi worker loads the skill and writes a complete stage report
- Merge guards armed, blocked on open PR #776
- Push CI: build/install/offline pass; claude-live/codex-live pending; pi-live skipping (manual dispatch only)
- Pi live lane dispatched: run 33121640433, `in_progress`, needs CI-E2E-PI environment approval (CL said "i'll approve")

### PR #777 — Gate the Pi front door bootstrap on resume (stacked on #776)
- Branch: `spacedock-ensign/gate-pi-frontdoor-bootstrap-on-resume` → base #776's branch
- Entity: `gate-pi-frontdoor-bootstrap-on-resume` (4av), at `done`/`approved-awaiting-merge`, `pr: #777`
- Content: pi.go containsResume guard wrapping piBootstrapPrompt append (mirrors frontdoor.go:553); TestPiResumeSuppressesBootstrapPrompt 7/7 green
- Merge guard armed, blocked on open PR #777
- Checks not yet registered (just created)

### Merge order
1. #776 first (the fix) — when pi-live goes green and CL merges
2. #777 second (resume-gate) — rebase to main after #776 merges

## Terminalized this session
- finish-pi-rejection-flow (#750) → done, archived
- pi-subagents-duplicate-extension-load (#747) → done, archived
- align-pi-compaction-with-force-boot (#754) → BLOCKED (worktree on exedev host, not local; merge guard can't verify gate authority)

## Closed PRs (superseded)
- #764 (embed approach, superseded by #776), #773 (stack-locked, recreated as #777), #774 (folded into #776), #770/#771 (stack tool mishap)

## Key lessons this session
- The bug: Pi firstActionBlock said "treat this dispatch file as your contract" — never told the worker to load the ensign skill. Fix: invoke /skill:ensign (Pi's skill-invoke per docs/skills.md) before the dispatch file, like Claude Skill() and Codex $spacedock:ensign.
- The embed approach (#764 original) was WRONG — duplicated the contract in the artifact body. Removed in the rework.
- ntn (pi-spawn-skill-name-resolution) was originally scoped wrong (comment-pin); reworked to BE the firstActionBlock fix.
- glm-5.2-vision-background stalls on the acceptance-report contract (template-echo, 0 tool calls) — needs a steer. Happened 4× this session.
- The completion guard (hasCompleteCommittedStageReport) caught a real stage-report format defect (inline evidence vs separate line) — the embed fix (t4) makes that guard reliable on Pi.
- gh stack link reorganizes PR bases into its own chain shape; close+recreate is the reliable stack path.

## After compaction
- Re-satisfy: fo-gate-lifecycle (loaded), fo-write-core (loaded), fo-merge-core (loaded), fo-dispatch-core (loaded), present-gate (loaded)
- Check: pi ci run 33121640433 status; #776 + #777 check status
- When #776 pi-live green + CL merges: record pr-merge:#776 on both entities, re-run merge guard, finalize
- When #777 green + CL merges: record pr-merge:#777, finalize
