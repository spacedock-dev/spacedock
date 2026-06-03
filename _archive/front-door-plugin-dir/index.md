---
id: jgc29m3pjb80efvmrc5bkc2n
title: spacedock claude --plugin-dir before `--` (captain dev workflow) + restore the live-e2e net
status: validation
source: "captain (2026-06-02) — CI-E2E live-runtime e2e fails: cobra `spacedock claude` rejects --plugin-dir before `--`; the captain's primary dev workflow and the e2e net both depend on it. 0.19.3 keystone."
started: 2026-06-02T05:06:55Z
completed: 2026-06-02T06:23:36Z
verdict: PASSED
score: "0.42"
worktree: 
issue:
mod-block: 
pr: "#260"
archived: 2026-06-02T06:23:36Z
---

The approval-gated live-runtime e2e (`TestLiveEnsignCycle`, `internal/ensigncycle/live_test.go:91-100`) has been silently broken since the cobra migration (#241): it invokes `spacedock claude --plugin-dir <repo> --skip-contract-check -p … --model … -- <task>` with the host flags BEFORE `--`, but the cobra-migrated front door rejects them (`unknown flag: --plugin-dir`) — host flags must follow `--`. The break stayed hidden because the env-approval gate kept every prior live run "waiting"; it surfaced the first time the job actually ran (this session, after CI approval). The live-e2e net has therefore been dead through the entire 0.19.2 sprint.

This also blocks the captain's PRIMARY dev workflow: nearly every real launch is `claude --plugin-dir /…/spacedock …`. For `spacedock claude` to truly REPLACE the hand-typed launcher (sprint goal #1), `spacedock claude --plugin-dir <dir>` (before `--`) must work — load the local plugin checkout and relax the contract gate — not require the `-- --plugin-dir` form.

## Scope (captain: "test + binary ergonomics")

- **Binary:** `spacedock claude` accepts `--plugin-dir <dir>` / `--plugin-dir=<dir>` BEFORE `--`, forwarding it into the host passthrough and relaxing the contract gate (the existing `hasPluginDir` gate-relax in `frontdoor.go` must see it). Multiple `--plugin-dir` allowed. Ideation decides whether `spacedock codex` gets the symmetric treatment and how broadly other host flags before `--` are forwarded (vs the documented `--` convention).
- **Test:** `TestLiveEnsignCycle` passes green against the fixed binary, restoring the live-e2e net.
- Lands on `next` as its own PR (whose own CI-E2E validates the fix end-to-end), then greens the in-flight #258 (5p) and #259 (#251) once they rebase onto the fixed `next`.

## Spike result (riskiest unknown, run FIRST)

**Question:** which cobra/pflag mechanism forwards host flags before `--` without swallowing spacedock's own flags (`--safehouse-*`, `--skip-contract-check`) or breaking the after-`--` passthrough.

**Answer: no generic flag-mode toggle works; the sound design is to register `--plugin-dir` as the ONE host flag spacedock teaches itself, and keep every OTHER host flag after `--`.** Spiked end-to-end against `github.com/spf13/pflag v1.0.9` (the pinned version) on the real `spacedock claude` arg shape:

- **`ParseErrorsWhitelist.UnknownFlags = true` SILENTLY EATS unknown flags.** Reading `pflag/flag.go` `parseArgs`/`parseLongArg`/`stripUnknownFlagValue` (v1.0.10 mirror of the v1.0.9 logic): an unknown long flag token is never appended to `f.args`, and `stripUnknownFlagValue` drops its following value too. Spike confirmed: `--plugin-dir /repo … -- task` → `passthrough=[task]`, the `--plugin-dir /repo` vanished. The worst outcome — data loss with no error. Rejected.
- **`SetInterspersed(false)` cannot distinguish spacedock flags from host flags.** It either errors on the first unknown flag (`unknown flag: --plugin-dir`) or, when the first token is a positional, stops parsing and dumps everything (including `--skip-contract-check`) into `Args()`. Rejected.
- **A manual generic pre-pass cannot pair value-taking host flags with their values** — it does not know each host flag's arity (`--verbose` takes none, `--model` takes one). Spike showed `--plugin-dir /repo -p Drive. --model sonnet` mis-split into `hostBefore=[--plugin-dir -p --model]` / `taskTokens=[/repo Drive. sonnet]`. This arity problem is *exactly* why the Option-2 grammar moved host flags after `--`; forwarding arbitrary host flags before `--` reintroduces it. Rejected.
- **WINNER — register `--plugin-dir` as a repeatable `StringArray` on the existing front-door pflag.FlagSet.** pflag then knows its arity (one value, repeatable), parses `--plugin-dir <dir>` / `--plugin-dir=<dir>` / repeated forms before `--` correctly, and the captured dirs re-inject into the FRONT of passthrough so they forward to the host and `hasPluginDir` sees them unchanged. Every other host flag (`-p`, `--model`, `--permission-mode`, `--output-format`, `--verbose`, …) stays after `--`. Spike end-to-end:
  - `--plugin-dir /repo --skip-contract-check -- -p Drive. --model sonnet` → inner argv `claude --agent spacedock:first-officer --plugin-dir /repo -p Drive. --model sonnet <prompt>`, `hasPluginDir=true`.
  - `--plugin-dir /repo "review the PRs"` (no `--`, captain dev workflow) → `--plugin-dir /repo` parsed, task `review the PRs`, gate relaxed.
  - `--plugin-dir /a --plugin-dir=/b --` → both captured (`/a`, `/b`), both forms.
  - A non-plugin-dir host flag left before `--` (e.g. `-p`) → loud `unknown shorthand flag: 'p'` error, NOT a silent mis-split — proving the live test MUST move its non-plugin-dir flags after `--`.

Spike harness: `/tmp/pflag-spike` (throwaway). The implementation's first failing test seeds from the WINNER row.

## Decisions (hardened at ideation)

- **D1 — `--plugin-dir` before `--` is ADDITIVE, not a replacement.** Both before-`--` and after-`--` `--plugin-dir` work; after-`--` stays verbatim-forwarded (unparsed positional) exactly as today. The help.go convention (`tokens after -- forward verbatim … and --plugin-dir`) is RETAINED and AUGMENTED: help gains a `--plugin-dir DIR` flag row (it is now a real spacedock-parsed flag) plus one example line `spacedock claude --plugin-dir ./checkout`. No documented behavior is removed.
- **D2 — only `--plugin-dir` is promoted; all OTHER host flags stay after `--`.** Per checklist #2's "move the non-plugin-dir ones after `--`": the live test moves `-p/--permission-mode/--output-format/--verbose/--model` to after `--`. Forwarding arbitrary host flags before `--` is rejected (arity problem, spike).
- **D3 — symmetric `spacedock codex` treatment.** `--plugin-dir` is a marketplace/plugin flag both hosts share and `hasPluginDir` already gates both `runClaude` and `runCodex`; promoting it in the shared `parseFrontDoorArgs` covers codex for free. No extra codex-only work.
- **D4 — placement: re-inject captured before-`--` `--plugin-dir <dir>` tokens at the FRONT of passthrough, before the after-`--` tokens.** Keeps the spacedock prompt the always-last assembled token (the Option-2 invariant `TestDanglingValueTakingHostFlagStillSwallows` pins) and keeps `hasPluginDir(passthrough)` the single gate-relax reader (no new field threaded through `runClaude`/`runCodex`).

## Acceptance criteria

**AC-1 — `spacedock claude --plugin-dir <dir>` BEFORE `--` is parsed, forwarded, and relaxes the gate.** The flag is accepted (no "unknown flag"), the captured dir(s) appear in the host passthrough as `--plugin-dir <dir>`, and `hasPluginDir(passthrough)` is true so the contract gate is relaxed (a failing manifest still launches). `--plugin-dir=<dir>` and repeated `--plugin-dir` both work; the captain's no-`--` form `spacedock claude --plugin-dir <dir> "task"` parses the dir as a flag and `task` as the launch prompt.
Verified by: a new `parseFrontDoorArgs` table case (before-`--` `--plugin-dir` space, equals, and repeated forms → captured into passthrough; a stray non-plugin-dir host flag before `--` → parse error) AND a `runClaude` seam test asserting the launch argv carries `--plugin-dir <dir>` and the gate is relaxed on a failing manifest, mirrored for `runCodex`.

**AC-2 — the live-e2e net is restored.** `go test -tags live -run TestLiveEnsignCycle ./internal/ensigncycle/` (the exact CI invocation) passes against the fixed binary. `live_test.go` keeps `--plugin-dir <repoRoot>` before `--` and moves `-p`, `--permission-mode`, `--output-format`, `--verbose`, `--model` to AFTER `--`, ahead of the fenced task.
Verified by: the live job green on this entity's own PR (CI-E2E), plus a local `-tags live` run recorded in the stage report (or SKIPPED-with-reason if no live credential on the implementer's machine, with the CI-E2E green as the binding proof).

**AC-3 — after-`--` host passthrough and spacedock-flag parsing are unchanged.** After-`--` `--plugin-dir /a --plugin-dir /b` still forwards verbatim in operator order; spacedock's own flags (`--safehouse`, `--safehouse-*`, `--skip-contract-check`) still parse before `--` in space and equals forms; the spacedock prompt is still the last assembled host-argv token.
Verified by: the existing front-door suite stays green unchanged — `TestParseFrontDoorArgs`, `TestPluginDirRelaxesGate`, `TestDanglingValueTakingHostFlagStillSwallows`, `TestDevLanePluginDirReachesLaunchSeam`, `TestClaudeFrontDoorLaunchesOnCompatible` — and `go test ./...` is green.

## Approach

In `internal/cli/frontdoor.go`:

- `bindFrontDoorFlags` gains one binding: `pluginDir: fs.StringArray("plugin-dir", nil, "Load a local plugin checkout (relaxes the contract gate); repeatable")`. Because `bindFrontDoorFlags` is the single source feeding both the parser and `declareFrontDoorHelpFlags`, the help row appears automatically with no drift (the comment on `frontDoorFlags` already promises this).
- `parseFrontDoorArgs` reads `*flags.pluginDir` after `Parse` and re-injects `--plugin-dir <dir>` pairs at the FRONT of `fd.passthrough`, then appends the post-`--` positionals after them. The after-`--` `--plugin-dir` path is untouched (those tokens are post-dash positionals, never seen by the flagset).
- `hasPluginDir`, `runClaude`, `runCodex` are UNCHANGED — they already read `fd.passthrough`.

In `internal/cli/help.go`: the `setFrontDoorHelp` prose keeps the after-`--` sentence and gains a one-line example `spacedock <host> --plugin-dir ./checkout`; the `--plugin-dir` flag row renders via the new binding. Update the ABOUTME/prose only where it states host flags ride "after `--`" to acknowledge `--plugin-dir` is also accepted before.

In `internal/ensigncycle/live_test.go`: move the five non-plugin-dir host flags to after `--`.

## Test plan

- **Unit (fixture-free, ~instant): `parseFrontDoorArgs` table** — add before-`--` `--plugin-dir` cases (space, equals, repeated; mixed with `--skip-contract-check` and a task; a non-plugin-dir host flag before `--` asserting a returned error). Cost: trivial. This is the failing-test-first artifact seeded by the spike WINNER row.
- **Unit: `runClaude`/`runCodex` seam** — extend the launch-parity suite with a before-`--` `--plugin-dir` case asserting the inner argv carries `--plugin-dir <dir>` and the gate relaxes on `tooOldBinaryManifest`. Cost: trivial, uses the existing `fakeHost`.
- **Regression: `go test ./internal/cli/...`** — the full existing suite (121 tests) must stay green, proving AC-3. Cost: ~1s.
- **Live (gated): `go test -tags live -run TestLiveEnsignCycle ./internal/ensigncycle/`** — the binary-change proof. Needs the built binary + a live credential (OAuth benchmark-token or `ANTHROPIC_API_KEY`); skips cleanly without one. The binding proof is the CI-E2E job green on this entity's PR. Cost: one live model run (~350s observed).

No fixture or golden changes needed. No schema/on-disk-format change.

## Stage Report: ideation

- DONE: AC pins the captain-workflow fix: `spacedock claude --plugin-dir <dir>` (BEFORE `--`) accepted, forwarded, relaxes the gate; `=`/repeated forms; before-vs-after-`--` decided ADDITIVE; frontdoor unit test named.
  AC-1 + D1 + D4; named unit test: a new `parseFrontDoorArgs` table case + a `runClaude`/`runCodex` launch-parity seam case.
- DONE: AC restores the live-e2e net: `go test -tags live -run TestLiveEnsignCycle ./internal/ensigncycle/` passes against the fixed binary; test-change-vs-binary-change decided BOTH (binary parses --plugin-dir before `--`, test moves the other five host flags after `--`).
  AC-2 + D2; the five non-plugin-dir flags (-p/--permission-mode/--output-format/--verbose/--model) move after `--`.
- DONE: Spike the riskiest unknown FIRST: which cobra mechanism forwards host flags before `--` without swallowing spacedock's own flags or breaking after-`--`; exercised end-to-end against the real `spacedock claude` arg shape.
  Spike result section; ran 4 mechanisms in /tmp/pflag-spike against pflag v1.0.9 — UnknownFlags silently eats, SetInterspersed(false) and generic pre-pass both unsound; WINNER = register --plugin-dir as StringArray. Verified the real live-test shape end-to-end.

### Summary

Spiked the riskiest unknown first and it overturned the obvious approach: pflag's `ParseErrorsWhitelist.UnknownFlags` SILENTLY DROPS unknown host flags (confirmed by reading pflag source + running it), and no generic flag-mode toggle can forward arbitrary host flags before `--` because spacedock cannot know each host flag's arity — which is precisely why the Option-2 grammar moved them after `--`. The sound design promotes ONLY `--plugin-dir` (the one host flag whose arity spacedock already reasons about via `hasPluginDir`) to a real repeatable StringArray on the front-door flagset; every other host flag stays after `--`, so the live test moves its five non-plugin-dir flags there. Decisions recorded: before-`--` is ADDITIVE (after-`--` retained), codex gets it for free via the shared `parseFrontDoorArgs`, and captured dirs re-inject at the front of passthrough so `hasPluginDir`/`runClaude`/`runCodex` stay unchanged and the prompt-always-last invariant holds. Baseline: all 121 internal/cli tests green before any change.

## Stage Report: implementation

- DONE: RED-FIRST: add the parseFrontDoorArgs table case (before-`--` --plugin-dir space/equals/repeated + the no-`--` captain form `--plugin-dir <dir> "task"`; a stray non-plugin-dir host flag before `--` → parse error) and watch it fail before the binding exists. Then bind `--plugin-dir` as a repeatable StringArray (D-WINNER), re-inject captured dirs at the FRONT of fd.passthrough; hasPluginDir/runClaude/runCodex unchanged.
  RED first: 6 new `--plugin-dir` table cases failed `err = unknown flag: --plugin-dir`; added `TestStrayHostFlagBeforeDashErrors` (stray `-p`/`--model` before `--` → parse error). Bound `pluginDir StringArray` in bindFrontDoorFlags; parseFrontDoorArgs re-injects `--plugin-dir <dir>` pairs at the FRONT of passthrough. Seam: extended `TestPluginDirRelaxesGate` with before-`--` claude/codex + captain no-`--`+task cases asserting argv carries `--plugin-dir <dir>` and the gate relaxes on tooOldBinaryManifest. Commit b20bed1f.
- DONE: Restore the live-e2e: in internal/ensigncycle/live_test.go keep `--plugin-dir <repoRoot>` BEFORE `--` and move `-p`, `--permission-mode`, `--output-format`, `--verbose`, `--model` to AFTER `--` (ahead of the fenced task). The exact CI invocation `go test -tags live -run TestLiveEnsignCycle ./internal/ensigncycle/` must pass against the rebuilt binary — run it locally if a live credential exists, else SKIPPED-with-reason and note CI-E2E is the binding proof.
  live_test.go rewired exactly as specified; compiles clean under `go vet -tags live`. Ran the exact CI invocation locally against a /tmp-built binary: the front-door fix is PROVEN end-to-end — the binary parsed `--plugin-dir <worktree>` before `--`, forwarded it, and claude loaded the local plugin (`"plugins":[{"name":"spacedock","path":".../.worktrees/...","source":"spacedock@inline"}]`), the exact arg-shape the cobra migration rejected. The cycle then 401'd on auth: the local `~/.claude/benchmark-token` (108B, `sk-ant-o…`) is expired/invalid (`apiKeySource:none`, `authentication_failed`) — a stale machine credential, not a code defect. Per AC-2 this is the SKIPPED-with-reason path: CI-E2E green on this entity's PR is the binding proof.
- DONE: AC-3 regression: the named existing front-door tests stay green (TestParseFrontDoorArgs, TestPluginDirRelaxesGate, TestDanglingValueTakingHostFlagStillSwallows, TestDevLanePluginDirReachesLaunchSeam, TestClaudeFrontDoorLaunchesOnCompatible), full `go test ./...` green, `go vet` clean; the spacedock-prompt-always-last invariant holds; help.go gains the --plugin-dir row + one example line.
  All 5 named tests PASS; full `go test ./...` green (10 packages); `go vet ./...` exit 0. Prompt-always-last invariant held: TestDanglingValueTakingHostFlagStillSwallows + the new before-`--` seam cases assert the spacedock prompt is the last argv token. help.go renders the `--plugin-dir stringArray` row (auto via bindFrontDoorFlags), before/after-`--` prose, and the `spacedock <host> --plugin-dir ./checkout` example — TestFrontDoorHelpCarriesDetail green for both hosts.

### Summary

Implemented the spike WINNER as settled: `--plugin-dir` is now a repeatable StringArray bound on the front-door pflag.FlagSet, so it parses BEFORE `--` (space/equals/repeated + the captain no-`--` form) and re-injects at the FRONT of passthrough — runClaude/runCodex/hasPluginDir untouched, the prompt-always-last invariant preserved (D4). Codex got it for free via the shared parseFrontDoorArgs (D3); before-`--` is additive to the after-`--` verbatim path (D1). The live-e2e net is restored in live_test.go (five non-plugin-dir host flags moved after `--`); a local run PROVED the front-door fix end-to-end (the binary forwarded `--plugin-dir` before `--` and claude loaded the worktree plugin) but the cycle 401'd on a stale local OAuth benchmark-token — the AC-2 SKIPPED-with-reason path, with CI-E2E on this entity's PR as the binding proof. Offline suite fully green, `go vet` clean. Code committed to the worktree branch (b20bed1f).

## Stage Report: validation

- DONE: Reproduce AC-1 — parseFrontDoorArgs table cases + launch-parity seam cases pass; independently built the binary and confirmed `--plugin-dir /x "task"` (no `--`), `--plugin-dir=/x`, repeated `--plugin-dir` all parse, forward as `--plugin-dir <dir>` into host argv, relax the gate on a failing manifest; a stray non-plugin-dir host flag before `--` errors loudly (no silent drop).
  Ran the behavior: TestParseFrontDoorArgs (6 before-`--` plugin-dir subcases) + TestStrayHostFlagBeforeDashErrors (3) + TestPluginDirRelaxesGate (claude/codex before-`--`, captain no-`--`+task, no-plugin-dir-fails-fast) all PASS. Binary built to /tmp/spacedock-val-binary; with a fake `claude` recording argv: space/equals/repeated forwarded `--plugin-dir <dir>`, gate relaxed (launched with no installed plugin), `-p`→`unknown shorthand flag: 'p'` exit 1 and `--model`→`unknown flag: --model` exit 1, claude NOT invoked. Without `--plugin-dir` the gate ran and denied launch (exit 1).
- DONE: AC-2 — confirm live_test.go matches the EXACT CI invocation; build the binary and run the front-door arg-parse end-to-end to plugin-load; confirm the local 401 is a credential issue not a code defect; record CI-E2E green as the binding proof.
  live_test.go:95-105 keeps `--plugin-dir <repoRoot>` before `--`, moves `-p`/`--permission-mode`/`--output-format`/`--verbose`/`--model` + task after `--` — the sole place that builds the live argv; runtime-live-e2e.yml:148 runs `go test -tags live -run TestLiveEnsignCycle`, so the test IS the CI invocation. Ran it locally against the built binary: claude's session-init loaded the worktree plugin (`"plugins":[{"name":"spacedock","path":".../.worktrees/spacedock-ensign-front-door-plugin-dir","source":"spacedock@inline"}]`, model resolved `claude-sonnet-4-6`), THEN two api_retry → terminal `api_error_status:401` / `authentication_failed` / `apiKeySource:none`. Local ~/.claude/benchmark-token is a 108B stale OAuth token (`sk-ant-oat01…`), no ANTHROPIC_API_KEY — the 401 is downstream of a fully successful parse+load, a stale machine credential not a code defect. CI-E2E green on the entity's PR (ANTHROPIC_API_KEY from secrets) is the binding live proof.
- DONE: AC-3 — the 5 named front-door tests + full `go test ./...` + `go vet` clean; the prompt-always-last invariant holds. Then PASSED/REJECTED.
  All 5 named tests PASS (TestParseFrontDoorArgs, TestPluginDirRelaxesGate, TestDanglingValueTakingHostFlagStillSwallows, TestDevLanePluginDirReachesLaunchSeam, TestClaudeFrontDoorLaunchesOnCompatible); `go test ./...` 680 passed / 12 packages; `go build ./...` ok; `go vet ./...` and `go vet -tags live ./internal/ensigncycle/` clean. Prompt-always-last held: TestDanglingValueTakingHostFlagStillSwallows asserts the spacedock prompt is the last argv token with the dangling flag immediately before it, and every before-`--` seam case asserts the same. Binary-level: after-`--` `--plugin-dir /a --plugin-dir /b` forwarded verbatim in order; help renders the `--plugin-dir stringArray` row + `spacedock claude --plugin-dir ./checkout` example.

### Summary

PASSED. Independently reproduced every "Verified by" by running the behavior, not re-reading: built the binary and exercised the real front door (space/equals/repeated `--plugin-dir` before `--` forward and relax the gate; stray host flags before `--` error loudly with claude not invoked), and ran the exact CI live invocation locally, which loaded the worktree plugin inline (`source: spacedock@inline`) the way impl proved before 401'ing on the stale local OAuth token — the AC-2 SKIPPED-with-reason path; CI-E2E green on the PR is the binding live proof. Full offline suite 680/680, `go vet` clean (incl. `-tags live`), prompt-always-last invariant holds, after-`--` passthrough and spacedock-flag parsing unchanged. No competing live-invocation shape exists. No PR open yet (the captain/FO opens it post-validation, per scope). Recommendation: PASSED.

## Stage Report: implementation (post-validation addendum)

- DONE: audit-surfaced test-only addition — `TestBareValuelessPluginDirErrors` in internal/cli/frontdoor_parse_test.go pins that a bare before-`--` `--plugin-dir` with no value is a LOUD parse error naming the flag (`flag needs an argument: --plugin-dir`, runClaude/runCodex return 1 before Launch), not the silent drop UnknownFlags would have produced — the load-bearing reason StringArray beat UnknownFlags. No logic change; `go test ./internal/cli/` green. Commit 18240751 (worktree branch, not pushed — FO pushes the full branch at the merge boundary).
