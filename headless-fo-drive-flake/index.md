---
id: hdr9dd8qywdc42cq700vefy4
title: Headless FO intermittently doesn't drive the cycle (-p / teams) — root-cause the stall
status: backlog
source: "FO (2026-06-02): live-e2e flakiness — under spacedock claude -p with EXPERIMENTAL_AGENT_TEAMS=1 the FO intermittently stalls (one model per run): on 38's PR sonnet booted, loaded the contract (Skill->Read), then went silent (~70s quiet trip) while opus drove the full cycle. The deferred FO-runtime-await question, now confirmed real. Investigation; gated on the transcript-artifact follow-up for data."
started:
completed:
verdict:
score: "0.34"
worktree:
issue:
---

The live-e2e flakiness root: under headless `spacedock claude -p` with teams enabled, the FO (first-officer) intermittently does NOT drive the cycle to done — one model per run stalls (failures alternated sonnet/opus across the 0.19.3 + 38 PRs). On 38's PR (#261) the new `streamWatcher` caught it cleanly: sonnet booted, started loading the operating contract (`Skill → Read`), then went silent (the ~60s quiet budget tripped at ~70s) while opus drove the full cycle green. So the binary + the watcher are sound; the FO MODEL sometimes stalls headless.

This is the FO-runtime-await / headless-drive question deferred from 38 — now confirmed real. It is an INVESTIGATION: root-cause WHY the FO stalls headless, then decide a fix.

## Dependency
GATED on `live-e2e-transcript-artifact` (3g — merge first): the investigation needs the FAILING run's transcript, which today is lost (gh truncation + binary-only artifact). Ideation may plan the approach now; data-gathering + root-cause execution follow the artifact landing + a captured stall.

## Hypotheses to test (ideation)
- The `-p` turn ends mid-await (the FO ended its turn waiting for a teammate completion the headless loop never re-wakes it for).
- The FO never reached single-entity/bare mode: the contract reserves bare blocking-Agent dispatch (which "prevents premature session termination in -p mode") for single-entity mode, but the live launch uses a GENERIC prompt + teams on → full team async-await under -p.
- Pure model nondeterminism (the model stalled loading the contract), independent of dispatch mode.

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — The stall is root-caused with evidence.** A captured failing-run transcript pinpoints WHERE the headless FO stalls (boot / contract-load / TeamCreate / dispatch-await) and WHY, distinguishing a runtime/harness cause from pure model nondeterminism.
Verified by: the captured transcript(s) + an analysis naming the stall point + cause.

**AC-2 — A fix or a justified disposition.** Either a runtime/harness change making the headless FO drive reliably (candidate: force bare/single-entity blocking dispatch under -p, per the contract's premature-termination guidance), OR a documented justified disposition (accept model nondeterminism + rely on 38's fast-localized detection + a bounded retry).
Verified by: the chosen direction with rationale + (if code) tests.

## Notes
- Touches the FO runtime (skills/first-officer/references/ + internal/cli/frontdoor.go) — high-stakes; investigation first. Gated on `3g`.
