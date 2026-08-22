# FO Gate Review — `align-pi-compaction-with-force-boot` implementation

## Question
Is the implementation correct: does the compaction boundary fire a boot read (not re-inject contract), with all 5 ACs satisfied and no scope breach?

## Verdict: APPROVE → validation

## Code verified (d24fbaada, off 753 tip 185b53477)
- `injectBootRecord` flag added, set by `session_compact` ONLY (which now clears `injectBootstrap`) — spacedock.ts. `session_start` still sets `injectBootstrap` (unchanged).
- `context` handler is now `async`; boot-record path awaits `pi.exec("spacedock", ["status","--boot","--identify","--json"])`, parses stdout, injects user message with `[SPACEDOCK-FO-BOOT-v2]` marker + directive + boot JSON. `FO_BOOTSTRAP_TEXT` fully removed from compaction path (retained for session_start).
- `pi.exec` failure → try/catch → directive-only fallback (never blocks session). Recorded in stage report.
- `hasBootRecord` dedup (AC-5); `agent_end` clears both flags.
- Insert-after-leading-compaction-summary logic shared by both paths (clean refactor).

## ACs (all PASS via focused test — real .ts through Node harness)
- AC-1 (value): compaction injects boot record (`"command":"boot"`), 0 bootstraps, no marker/contract pointer — PASS
- AC-2 (mechanism): `pi.exec` called once with `["status","--boot","--identify","--json"]` — PASS
- AC-3 (scope): session_start unchanged (FO_BOOTSTRAP_TEXT + "Load the skill") — PASS
- AC-4 (child exemption): PI_SUBAGENT_CHILD=1 → undefined on both paths — PASS
- AC-5 (dedup): existing boot record → undefined — PASS

## Scope
- Expected surface: ideated +15-20 net LOC; actual +91 net (2 files). Over because the two-path async handler + shared insert logic + test harness (pi.exec mock, execCalls, await) is larger than the single-path original. The stage report names this honestly. No scope breach: observable-semantics declaration holds — only compaction-boundary context content changes; no FO bootstrap content, gate/dispatch/state mechanics, command grammar, stored formats, or authority touched.
- gofmt: Go test file clean. (The .ts file was mistakenly gofmt'd in the checklist — gofmt is for .go; the .ts is idiomatic TS, `git diff --check` clean.)
- Pre-existing failures: 4 internal/cli tests fail on 753's tip too (env-driven: PI_CODING_AGENT set + /tmp helper-binaries warning) — NOT caused by this change.

## Validation plan
Validation = the focused test (the ideation specified "fixture test, no live run needed" — the behavior is the hook, which the fixture exercises against the real .ts). FO runs the focused test as the validation proof. Then open PR on 753 tip, add to stack #748, trigger the pi-live lane (the consolidated stack signal the captain wants).
