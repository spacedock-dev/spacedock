---
id: 5h0chdcad99dq0z50qzwf1za
title: Replace README-substring (prose-grep) test assertions — prove the seam, not the prose
status: ideation
source: "captain + nb (readme-main-flip-reconciliation) reconciliation 2026-06-07 — PR #315 edits tests/test_codex_plugin_packaging.py to assert README content by substring (assert \"spacedock codex\" in readme). Not on next today; when #315's content lands on main it must be replaced. Same proof-policy class as #309/4q and the survey signal-correction work."
started: 2026-06-15T01:55:42Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0203-fo-efficiency
group: proof-policy
sprint-readiness:
---

A test that asserts README content by substring (`assert "spacedock codex" in readme`) proves only that the text is present — the banned prose-grep tautology. When PR #315's content lands on `main`, replace that assertion with a behavioral one.

## Problem

PR #315's `tests/test_codex_plugin_packaging.py` greps the README for substrings to "test" Codex install/launch. A substring match over a doc the implementer wrote can't fail on a valid paraphrase and passes on an inverted clause — it tracks spelling, not behavior. This is exactly the proof-policy class the project fences (the 4q reframe, #309, the survey signal-correction `references/queries.sql` work).

## Proposed approach

When #315's content reaches `main` (or at the flip), replace the README-substring assertion: prove Codex install/launch by exercising the SEAM (the packaging/launch behavior — does `spacedock codex` resolve + launch), not by grepping README prose. If the only real invariant is "the README documents the command," that is not a behavioral test and should be dropped, not kept as a grep.

## Out of scope

- Rewriting #315 itself (it's a separate PR; this task handles the assertion once it lands).
- The general prose-grep ban already shipped for skills (#309/4q) — this is the README/packaging-test instance.

## Acceptance criteria

Ideation/implementation fills in. Sketch: no test asserts README content by substring/regex; Codex install/launch is proven by exercising the packaging/launch seam (output/exit/observable behavior), with a fixture or live check that can actually fail on broken behavior.

## Test plan

The replacement test goes RED on broken Codex install/launch behavior (not on a paraphrase of the README).

## FO note — concrete trigger (2026-06-14)

The dev-README slim (commit `a9e669ae`, FO-direct, unpushed local `main`) relocated the Runtime-Live-CI / Codex-watchdog content out of `docs/dev/README.md` into `docs/runtime-live-ci.md`, which EXPOSED two prose-grep doc-contract guards this task targets:

- `internal/ensigncycle/shared_scenarios_docs_test.go` — `TestSharedScenarioDocsContract`
- `internal/ensigncycle/codex_collab_wait_watchdog_test.go` — `TestCodexForegroundWaitWatchdogDocsContract`

Both assert literal clause strings are present in `docs/dev/README.md` (e.g. `### Shared runtime scenarios`, `Codex foreground-wait watchdog`, `codexScenarioRunners()`); they now fail (clause count 0) because the prose moved. Same README-substring anti-pattern this task removes — surfaced by the 0.20.4 `e6a` implementation, which found them red at its branch base.

Blocking relationship: this breakage blocks BOTH the README-slim push AND any 0.20.4 `e6a` merge (e6a's base includes `a9e669ae`). Resolution (this task's call): retarget the guards to `docs/runtime-live-ci.md`, or convert them to bind an independent source per the proof policy — not a relocated prose-grep.
