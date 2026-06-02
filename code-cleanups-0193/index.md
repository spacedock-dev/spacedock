---
id: 1xk9bz4fr7qefgcqz6stzpkk
title: 0.19.3 code-quality cleanups — release regex clobber, statecommit test phrase, StandingTeammate dedup, prose-marker brittleness, init→install rename drift, unknown-subcommand silent exit
status: validation
source: sprint-end antipattern reviews (2026-06-01) — 0.19.3 minor-findings bucket (all Minor)
started: 2026-06-02T04:57:30Z
completed:
verdict:
score: "0.22"
worktree: .worktrees/spacedock-ensign-code-cleanups-0193
issue:
---

A batch of low-risk code-quality cleanups surfaced by the sprint-end staff-SWE + AI-eng reviews. All Minor; grouped because each is too small to track alone.

- **internal/release** — collapse the twin version regexes into one, and fix the latent multi-`version`-key `ReplaceAll` clobber (a stamp file with more than one `version` key would over-rewrite). Add a multi-version fixture proving the targeted replace.
- **build_statecommit_test** — assert the "never a bare `git add -A`" guidance phrase is present in the emitted split-root state-commit guidance (regression-guards the concurrency-safety instruction).
- **stateCommitGuidance** (`internal/dispatch/build.go:459`) — fold a push-on-commit line into the guidance so state commits are reminded to push.
- **StandingTeammate dedup** — `dispatch.StandingTeammate` and `claudeteam.StandingTeammate` are duplicate structs + duplicate identity mappers; collapse to one.
- **prose-test brittleness** — de-brittle the single-word prose-test markers (match a phrase, not a lone word); trivial `min`/`itoa` cleanup in `prose_neutrality_test.go`.
- **init→install rename drift** (`internal/contract` + `internal/cli` + docs) — the `init`→`install` rename left stale `spacedock init` callouts: `internal/contract/contract.go:200` tells a binary-present *version-mismatch* user to `Upgrade it: spacedock init --host %s` (a command that no longer exists), `internal/cli/init.go` still prefixes its error messages `spacedock init:` (×4: lines 33/50/102/110), and the docs (`README.md` ×3 lines 24/48/51, `docs/install-journey.md:61`) still document `spacedock init --host claude`. Sweep all sites to `install`. The `internal/cli/init.go` filename is itself stale — renaming the file is optional polish, not required. Surfaced live by the captain (2026-06-02) testing the install path; the version-mismatch user currently gets a dead command.
- **unknown-subcommand silent exit** (`internal/cli` root) — `spacedock <unknown> --someflag` exits 2 with ZERO output (captain hit `spacedock init --host claude` → silent exit 2). Bare `spacedock init` *does* print the usage block, but adding a flag makes cobra swallow even that. Any unknown subcommand must print the usage block to stderr regardless of trailing flags/args, and exit non-zero. Behavioral fix + a test asserting non-empty usage output on an unknown command carrying a flag.
- **releasing.md missing bump-calendar step** (captain directive 2026-06-02, DOC-ONLY) — `marketplace.json` carries a CALENDAR version (`0.0.YYYYMMDDNN`, the `claude plugin update` re-pull key) stamped by `bump-calendar`, a SEPARATE release step from `stamp-version` (which stamps the plugin.json `version`). The documented release sequence in `docs/releasing.md` runs only `stamp-version`, so `bump-calendar` was missed for 0.19.3, leaving `marketplace.json` stale at `0.0.2026060105` vs `plugin.json` `0.19.3`. Add `go run ./cmd/spacedock-release bump-calendar .claude-plugin/marketplace.json` to the documented sequence (DOC ONLY — do not bump `marketplace.json` itself in this change).

## Root cause + minimal fix (per cleanup)

**1. release twin-regex + multi-`version` clobber** (`internal/release/release.go`).
Root cause: `topLevelVersionRe` (line 17) and `entryVersionRe` (line 41) are byte-identical (`("version"\s*:\s*")[^"]*(")`) — twin regexes. Both `StampVersion` (line 34) and `BumpCalendarVersion` (line 73) call `regexp.ReplaceAll`, which rewrites EVERY `"version":` match in the blob; the doc-comments claim "the first such member" / "the first match", but `ReplaceAll` does not honor first-only. A manifest with a second nested `version` key would be over-rewritten (latent clobber).
Minimal fix: one shared regex; replace `ReplaceAll` with a replace-FIRST-match-only helper (`FindSubmatchIndex` + splice, or a `ReplaceAllFunc` with a once-guard) used by both functions. No public-signature change. (No spike: Go `regexp` first-match-splice is a proven stdlib idiom.)

**2. statecommit "never a bare git add -A" phrase assertion** (`internal/dispatch/build_statecommit_test.go`).
Root cause: the three existing tests assert the resolved `git -C` paths and brace-freedom, but NONE assert the concurrency-safety phrase "never a bare `git add -A`" that `stateCommitGuidance` (`build.go:459`) emits — so a future edit could silently drop the warning. Pure test-coverage addition.
Minimal fix: extend `TestStateCommitGuidanceResolvesPaths` (or a sibling) to assert the emitted body contains the literal "never a bare `git add -A`" substring for the split-root case. No production-code change.

**3. push-on-commit line in `stateCommitGuidance`** (`internal/dispatch/build.go:459`).
Root cause: the emitted guidance stops at "Retry on index.lock contention after a short wait." and never tells the worker to push the state branch — yet the ensign shared core (`ensign-shared-core.md` "Multi-writer sync") requires `git -C {state_checkout} push origin {state_branch}` then `pull --rebase` on rejection. The dispatched prose omits the sync step entirely.
Minimal fix: append a push-on-commit sentence to `stateCommitGuidance`. NOTE — this needs the state-branch name; `stateCommitGuidance(stateCheckout, entityPath)` does not currently receive it. The implementer must either (a) thread the resolved state-branch through, or (b) if the branch is not readily available at that call site, emit a branch-neutral "then push the state branch and `pull --rebase` on rejection" reminder. Decide at implementation; this is the one item whose blast radius (a new parameter) edges past trivial — flag if (a) proves non-trivial. Guard with a phrase assertion in the same statecommit test.

**4. StandingTeammate dedup** (`internal/dispatch/mods.go:93` + `internal/claudeteam/standing.go:16` + `toClaudeTeammates` at `internal/dispatch/standing.go:184`).
Root cause: two structurally identical structs (`{Name, Description, RoutingUsageBody string}`) plus an identity field-copy mapper, a leftover of the `zs` claude-runtime-segregation split.
Import graph (verified): `dispatch` imports `claudeteam`; `claudeteam` does NOT import `dispatch` (no cycle). So the single struct must live in `claudeteam` (the lower layer) for `dispatch` to use it without reversing the dependency.
Minimal fix: delete `dispatch.StandingTeammate` and `toClaudeTeammates`; have `EnumerateDeclaredStandingTeammates` return `[]claudeteam.StandingTeammate` and `RenderStandingTeammatesSection` take it directly. DESIGN NOTE: this makes the runtime-neutral `dispatch` enumerator return a `claudeteam`-package type, slightly nudging the segregation boundary the `zs` split drew. The alternative (struct stays in `dispatch`, `claudeteam` imports it) creates the cycle and is rejected. Implementer confirms the chosen direction compiles cycle-free.

**5. prose-marker brittleness + min/itoa cleanup** (`internal/hostneutrality/prose_neutrality_test.go`).
Root cause (a): `hostQualifierMarkers = []string{"Codex", "codex"}` (line 41) treats a span as host-qualified if it contains the BARE word "Codex"/"codex" anywhere — but the intended qualifier is the `X on Codex, Y on Claude` PHRASE shape (per the function's own comment). A passing mention of "codex" falsely qualifies a span, letting an unqualified Claude token slip through.
Root cause (b): local `min` (lines 210-216) shadows the go1.22 builtin `min`; `itoa` (lines 219-236) reimplements `strconv.Itoa` with a stale "avoids importing strconv twice" rationale (`strconv` is imported zero times in the package).
Minimal fix: change markers to the phrase form (e.g. `" on Codex"` / `" on codex"`); delete local `min` (use builtin); delete `itoa`, import `strconv`, use `strconv.Itoa`. Add a fixture span proving a bare-word "codex" mention no longer qualifies while the real `… on Codex …` span still does.

**6. init→install rename drift** (`internal/contract` + `internal/cli` + docs + lockstep tests). Full site enumeration (grep over `internal/` + `README.md` + `docs/`, excluding `.spacedock-state/_archive` history and the entity's own body):
- `internal/contract/contract.go:200` — `pluginPredatesContractRemedy` emits the DEAD command `spacedock init --host %s` (the live defect the captain hit). Sibling remedies `tooOldPluginRemedy` (line 183) and `noPluginMessage` (line 210) were ALREADY fixed to `install` by the `cli-cobra-redesign` work — this one was missed.
- `internal/contract/contract.go:43`, `:189`, `internal/contract/doctor.go:21` — evergreen comments still say `spacedock init`; update to match the code.
- `internal/cli/init.go:33,50,102,110` — 4 error-prefix self-labels `spacedock init:`. (Filename `init.go` itself is optional polish, NOT required.) Also the ABOUTME line 1 names `init`.
- `README.md:24` (brew lane block), `:48` (upgrade lane block), `:51` (prose "`spacedock init` reinstalls…").
- `docs/install-journey.md:61` (step-3 install command).
- LOCKSTEP test assertions that currently PIN the dead command and will break on the fix: `internal/contract/contract_test.go:81` (`predates-contract-empty` case asserts `spacedock init --host claude`), `:136` and `:152` (`TestPluginPredatesContractRemedy` asserts the init one-liner). Flip these three to `install`. (Comment-only mentions at `contract_test.go:125`, `init_test.go:12,104`, `init_devbranch_test.go:1,14`, `frontdoor_test.go:142,145` update for accuracy.) `frontdoor_test.go:167` and `contract_test.go:177` are NEGATIVE assertions already forbidding `init` — they stay and will now also pass for the predates-contract path once contract.go:200 is fixed.

## Spike — cobra silent-exit-2 mechanism (riskiest unknown, exercised first)

Root cause, reproduced live (`go build -o /tmp/sd-test ./cmd/spacedock`): the root command (`cli.go:91`) has `SilenceErrors`/`SilenceUsage: true` and flag parsing ENABLED (no `DisableFlagParsing`). For `spacedock init --host claude`, cobra resolves no `init` subcommand and runs the root `RunE` — but FIRST parses `--host claude` against the root flagset, which has no `--host`, returning a silenced `unknown flag` error → exit 2 with ZERO output. The root `RunE`'s `unknownCommand(args[0])` path never runs. Bare `spacedock init` (no flag) DOES print the usage because there is no flag to error on. Reproduced:
- `init --host claude` → exit 2, no output (the captain's case)
- `bogus --someflag` → exit 2, no output
- `bogus` / `init` (no flag) → `unknown command: …` + help, exit 2 (correct)

Fix verified by a throwaway cobra spike (`FParseErrWhitelist{UnknownFlags: true}` on the root): unknown flags stop erroring during parse, fall through to `RunE`, and `args[0]` correctly carries the command token →
- `init --host claude` → `unknown command: init`
- `bogus --someflag` → `unknown command: bogus`
- `--version`, bare `spacedock`, `install`, `bogus` (no flag) → all UNCHANGED.

HONEST SCOPE LIMIT (recorded so it is not a silent regression): with the whitelist, a LONE unknown flag and no command token (`spacedock --bogus`) is swallowed entirely — `RunE` sees `args=[]` — so it goes from today's silent exit 2 to printing help + exit 0. No existing test pins the lone-flag case (`TestUnknownCommand` only covers bare `bogus`), so no contract breaks. This is within the captain's stated requirement ("unknown SUBCOMMAND … regardless of trailing flags/args") — the lone-flag-no-command case is out of that scope and the help+exit-0 outcome is acceptable. Minimal fix: add `FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true}` to the root `cobra.Command` (one struct field).

## Acceptance criteria

**AC-1 — release: one regex, replace-first-only, multi-version safe.** `StampVersion` and `BumpCalendarVersion` share one regex and each rewrite only the FIRST `"version"` match; a manifest with a second `"version"` key keeps that second key untouched.
Verified by: a multi-`version`-key fixture test in `internal/release` that FAILS against the current `ReplaceAll` (asserts the second key is unchanged) and PASSES after the replace-first fix. `go test ./internal/release/` green.

**AC-2 — statecommit guidance: concurrency phrase + push reminder pinned.** The split-root state-commit guidance the dispatch emits contains both the "never a bare `git add -A`" concurrency-safety phrase and a push-the-state-branch (+ `pull --rebase` on rejection) reminder; single-root dispatches emit neither.
Verified by: assertions in `internal/dispatch/build_statecommit_test.go` over the emitted dispatch body for the substring "never a bare `git add -A`" AND a push reminder substring (split-root case); the existing `TestSingleRootNoStateCommitGuidance` negative still passes. `go test ./internal/dispatch/` green.

**AC-3 — StandingTeammate is one struct, no identity mapper.** A single `StandingTeammate` type exists (in `claudeteam`); `dispatch.StandingTeammate` and `toClaudeTeammates` are gone; `EnumerateDeclaredStandingTeammates` returns `[]claudeteam.StandingTeammate` consumed directly by the render. No behavior change to the rendered standing-teammates section.
Verified by: `go build ./...` compiles against the single definition (a residual second definition would not compile / a duplicate-type grep finds one `type StandingTeammate struct`); the existing standing-teammate render tests stay green. `go test ./internal/dispatch/ ./internal/claudeteam/` green.

**AC-4 — prose oracle qualifies on the phrase, not a bare word; helpers trimmed.** `hostQualifierMarkers` matches the `… on Codex …` phrase shape, so a span that merely mentions "codex" in passing no longer counts as host-qualified; the local `min` and `itoa` are removed in favor of the builtin `min` and `strconv.Itoa`.
Verified by: a `prose_neutrality_test.go` fixture-span case asserting (a) a bare-word "codex" span carrying an unqualified Claude token now FAILS the oracle (was wrongly passing) and (b) a real `… on Codex, … on Claude …` span still passes; the package compiles with no local `min`/`itoa`. `go test ./internal/hostneutrality/` green.

**AC-5 — no dead `spacedock init` user-facing command remains; the predates-contract remedy is runnable.** Every user-facing `spacedock init` callout (the predates-contract remedy, the `init.go` error prefixes, README, install-journey) names `spacedock install`; the predates-contract verdict message names `spacedock install --host <host>`.
Verified by: a grep over `internal/`, `README.md`, `docs/` (excluding `.spacedock-state/_archive`) finds no user-facing `spacedock init --host` / `spacedock init:` string; the retargeted `contract_test.go` assertions (lines 81/136/152, now asserting `install`) plus the existing negative `TestCompareHostSubstitution`/`TestGateRemedyNamesLiveInstallCommand` pass. `go test ./internal/contract/ ./internal/cli/` green.

**AC-6 — unknown subcommand always prints usage + exits non-zero.** `spacedock <unknown>` and `spacedock <unknown> --flag` both print `unknown command: <unknown>` + the grouped help to stderr and exit 2 — never a silent exit 2. (Lone-flag-no-command `spacedock --bogus` is out of scope per the spike's recorded limit.)
Verified by: a CLI test in `internal/cli` driving `Run([]string{"init", "--host", "claude"}, …)` (and a `bogus --someflag` case) asserting non-empty stderr containing "unknown command:" and exit code 2; the existing `TestUnknownCommand` (bare `bogus`) stays green. `go test ./internal/cli/` green.

**AC-7 — releasing.md documents the bump-calendar step in sequence.** The `docs/releasing.md` release-run section names `bump-calendar` alongside `stamp-version`, so a release operator following the doc bumps both the plugin.json `version` AND the marketplace.json calendar key. DOC-ONLY — `marketplace.json` itself is NOT bumped in this change.
Verified by: the release-run section of `docs/releasing.md` (step 3) names `go run ./cmd/spacedock-release bump-calendar .claude-plugin/marketplace.json` in sequence with `stamp-version`; the invocation matches `cmd/spacedock-release` (one `<marketplace.json>` arg, confirmed against the binary's usage).

## Test plan

- All proofs are Go unit/behavioral tests at the claim's level — no live workflow run needed; the only runtime claim (cobra exit behavior) is exercised by `cli.Run` in-process and was additionally reproduced against a built binary in the spike above. Estimated cost: sub-second per package; `go test ./...` for the final green sweep.
- TDD order per item: write the failing assertion first (multi-version fixture asserting the second key survives; the `init --host claude` → "unknown command" CLI assertion; the bare-word-codex oracle span; the missing concurrency/push phrases; the `install` remedy). The cobra fix and the release replace-first are the two with a runtime/parser mechanism — both already spike-verified above, so the spike's observed outcomes seed the first test.
- Riskiest-first: the cobra mechanism (proven) gates AC-6; the release replace-first idiom (proven stdlib) gates AC-1. Everything else composes already-proven behavior (struct dedup, phrase substrings, comment/doc sweep).
- Final gate: full `go test ./...` (excluding the pre-existing env-gated `TestCodexResolveManifestAgainstInstalledHost`, unrelated to this diff) green + `gofmt`/`go vet` clean.

## Scope assessment (all Minor?)

Five of six items are genuinely Minor / no-behavior-change (release internal regex; statecommit test; StandingTeammate dedup; prose-test markers; rename sweep). TWO carry an intentional, scoped behavior change that the captain explicitly requested: AC-6 (unknown-subcommand now prints usage — the whole point) and AC-2's push reminder (new guidance text). The single item that could exceed Minor is AC-2's push-on-commit line IF threading the state-branch name into `stateCommitGuidance` requires a new parameter through the call chain — flagged for the implementer to confirm or fall back to a branch-neutral reminder. No item requires a CONTRACT_VERSION bump.

## Notes
- Touches `internal/release` + `internal/dispatch` + `internal/cli` + `internal/contract` + docs + tests — NOT the `internal/status` serialized lane, so its implementation worktree can run alongside status-lane work. `internal/cli/cli.go` (root command) is a hot single-writer file — coordinate the implementation worktree with any other `cli.go` writer. The StandingTeammate dedup traces to the `zs` claude-runtime-segregation split (upstream-flow note).
- The init→install drift + unknown-subcommand-silent-exit were folded in by the captain (2026-06-02) from a live install-path test; the sibling FO-contract guidance gap (binary-*absent* startup routes to the missing-binary `doctor`) is tracked separately as `binary-absent-fo-bootstrap`.

## Stage Report: ideation

- DONE: Each cleanup has a root cause + minimal fix recorded; the init->install sweep enumerates EVERY stale `spacedock init` site (grep over internal/ + README + docs) so implementation lands complete, not partial.
  Per-cleanup "Root cause + minimal fix" section (6 items) + AC-5's full site list: contract.go:200/43/189, doctor.go:21, init.go:33/50/102/110+L1, README.md:24/48/51, install-journey.md:61, plus the 3 lockstep test assertions (contract_test.go:81/136/152) that currently PIN the dead command. `_archive` history and the entity body excluded.
- DONE: The unknown-subcommand-silent-exit fix names the cobra mechanism causing the silent exit 2 and specifies a behavioral test.
  Spike section: root cause is the root command parsing `--host` against its own flagset (no `DisableFlagParsing`) + `SilenceErrors/SilenceUsage` → silenced "unknown flag" → exit 2. Reproduced live against a built binary; fix `FParseErrWhitelist{UnknownFlags: true}` verified by throwaway cobra spike. AC-6 test: `Run(["init","--host","claude"])` + `bogus --someflag` assert "unknown command:" on stderr + exit 2.
- DONE: ACs hardened to be gate-checkable (not prose-only); confirm each cleanup is genuinely no-behavior-change, and flag any item that turns out to be more than Minor.
  6 ACs each name a Go test + a `go test ./<pkg>/` green oracle; "Scope assessment" section flags AC-2's push-reminder (state-branch threading) as the one item that could edge past Minor, and records AC-6's intentional behavior change + its honest lone-flag scope limit.

### Summary

Hardened all six cleanups with root cause, minimal fix, and a gate-checkable AC apiece. Exercised the two riskiest unknowns first: reproduced the cobra silent-exit-2 live (root flagset eats the unknown flag before RunE under SilenceErrors) and verified the `FParseErrWhitelist{UnknownFlags:true}` fix via a throwaway spike — recording the honest limit that a lone unknown flag with no command token shifts from silent-exit-2 to help+exit-0 (no test pins it, within the captain's "unknown subcommand" scope). Found the rename drift is deeper than a doc sweep: contract.go:200's `pluginPredatesContractRemedy` still emits the dead `spacedock init` AND three test assertions actively pin it (siblings `tooOldPluginRemedy`/`noPluginMessage` were already fixed), so the fix is contract.go + 3 lockstep test flips + comments/docs. Flagged AC-2's push-on-commit line as the only item that could exceed Minor (needs the state-branch name threaded into `stateCommitGuidance`).

## Stage Report: implementation

- DONE: All 6 cleanups land with their named tests green, TDD failing-first: release replace-first multi-version fixture; statecommit phrase + push reminder; StandingTeammate single struct cycle-free; prose phrase-marker + min/itoa trim; init->install full sweep; cobra unknown-subcommand usage.
  Commits 5b699456 (release+dispatch), bed1f3b4 (hostneutrality), fce91a2f (cli+contract+docs) on spacedock-ensign/code-cleanups-0193. Each test written failing-first and confirmed failing for the right reason before the fix.
- DONE: AC-1 release: one shared `versionRe` + `replaceFirstVersion` (FindSubmatchIndex splice) replaces ReplaceAll in both StampVersion and BumpCalendarVersion.
  `TestStampVersionRewritesOnlyFirstVersionKey` + `TestBumpCalendarVersionRewritesOnlyFirstVersionKey` assert a second nested `version` key survives; `go test ./internal/release/` 6/6 green.
- DONE: AC-2/AC-3 statecommit: threaded the resolved state branch via `status.StateBranch(workflowDir)` into `stateCommitGuidance(stateCheckout, entityPath, stateBranch)` and appended a push + `pull --rebase` reminder; concurrency phrase pinned.
  Chose option (a) — threading was trivial since `status.StateBranch` already derives it; no call-chain breakage. `build_statecommit_test.go` asserts the phrase + resolved push command; single-root negative still passes. Updated the parity-harness strip regex to cover the new line ending. `go test ./internal/dispatch/` green.
- DONE: AC-3 StandingTeammate dedup: deleted `dispatch.StandingTeammate` + `toClaudeTeammates`; `EnumerateDeclaredStandingTeammates` now returns `[]claudeteam.StandingTeammate` consumed directly by the render.
  `grep "type StandingTeammate struct"` finds one (claudeteam); `go build ./...` cycle-free; `go test ./internal/dispatch/ ./internal/claudeteam/` green.
- DONE: AC-4 prose oracle qualifies on the contrast, not a bare word; min/itoa trimmed.
  Markers changed to require BOTH `Codex` AND `Claude` (capitalized) — see DEVIATION note below; `TestSpanHostQualifiedRequiresContrast` proves a bare-word codex span with no Claude contrast no longer qualifies and the real `… on Codex, … on Claude …` span still does; local `min`/`itoa` deleted (builtin `min` + `strconv.Itoa`). `go test ./internal/hostneutrality/` green.
- DONE: AC-5 init->install sweep COMPLETE: contract.go:200 (live dead command) + comments (43/189, doctor.go:21), init.go x4 prefixes + ABOUTME, README.md x3, install-journey.md:61, and the 3 lockstep test flips (contract_test.go:81/136/152) + comment-only mentions (init_test.go, init_devbranch_test.go).
  `grep "spacedock init --host"/"spacedock init:"` over internal/+README+docs (minus _archive) finds only comments and the negative-assertion test strings (frontdoor_test.go:167, contract_test.go:177) that must stay; no user-facing dead command. `go test ./internal/contract/ ./internal/cli/` green.
- DONE: AC-6 unknown-subcommand: added `FParseErrWhitelist{UnknownFlags: true}` to the root cobra command.
  `TestUnknownCommandWithFlag` (`init --host claude` + `bogus --someflag`) asserts "unknown command:" on stderr + the usage block + exit 2; existing `TestUnknownCommand` stays green. `go test ./internal/cli/` green.
- DONE: AC-7 releasing.md documents bump-calendar in sequence (cleanup #7, captain directive).
  `docs/releasing.md` step 3 now runs `go run ./cmd/spacedock-release bump-calendar .claude-plugin/marketplace.json` alongside stamp-version and commits marketplace.json; invocation confirmed against the binary's usage (one `<marketplace.json>` arg). DOC-ONLY — marketplace.json not bumped.
- DONE: AC-2's state-branch threading resolved (option a, threaded through); full `go test ./...` green + gofmt/vet clean.
  `go test ./...` all packages ok (incl. `TestCodexResolveManifestAgainstInstalledHost`, which passed in this env); `go vet ./...` no issues; gofmt clean on all touched files (pre-existing `internal/status/enum_scope_test.go` gofmt-unclean is on origin/next, untouched by this branch).

### Summary

All 6 cleanups + cleanup #7 (releasing.md bump-calendar doc) landed across three worktree commits, each test written failing-first. Two intentional, captain-requested behavior changes: AC-6 (unknown subcommand now prints usage) and AC-2's push reminder. DEVIATION on AC-4: the entity's literal marker prescription (`" on Codex"`) breaks the real-file test because first-officer-shared-core.md:145 host-qualifies `context-budget` via "Codex declares none; Claude supplies one" — a genuine contrast but not the `on Codex` phrase shape (reproduced empirically). Implemented the more robust intent — a span qualifies iff it names BOTH `Codex` and `Claude` (the contrast, capitalized so the lowercase `claude-` in `claude-team` does not self-satisfy it) — which satisfies AC-4's fixture requirement and keeps the real-file test green. Flagged to team-lead non-blocking. AC-2's state-branch threading was trivial (option a) — `status.StateBranch` already derives it.

## Stage Report: validation

- DONE: Reproduce every AC's 'Verified by'; named per-package tests + full `go test ./...` green + gofmt/vet clean; spot-check failing-first for AC-6 and AC-1.
  Per-package (release/dispatch/claudeteam/hostneutrality/contract/cli) 309 passed; full `go test ./...` 720 passed/12 pkgs; gofmt clean on all 19 touched files; `go vet` clean. AC-1 revert→ReplaceAll: nested `metadata.version`/`schema.version` clobbered, both new tests FAIL for the right reason. AC-6 revert→drop whitelist: stderr empty (silent exit 2), `TestUnknownCommandWithFlag` FAILS, bare `TestUnknownCommand` still passes.
- DONE: Scrutinize behavior changes + AC-4 deviation; confirm valid-subcommand unknown flag NOT swallowed; AC-4 contrast genuinely qualifies.
  Built binary 8-case matrix: unknown-sub+flag→`unknown command:`+usage+exit2; `install --bogusflag`/`doctor --nope`→`unknown argument` exit2 (NOT swallowed — whitelist is root-only, not inherited); lone-flag→help+exit0 (recorded in-scope limit); `--version`/bare unchanged. AC-4: reverting to literal `" on Codex"` breaks the real-file test (first-officer-shared-core.md:145 `context-budget` flagged); contrast rule passes. Capitalized markers verified — `codex exec`+`claude-team` (lowercase) do NOT self-qualify; real `on Codex`/`Codex…Claude` contrasts do. AC-4 old-bare-word revert→`bare-codex-mention-no-claude` case FAILS.
- DONE: AC-5 init->install sweep COMPLETE; AC-7 releasing.md names bump-calendar in sequence.
  Grep (internal/+README+docs, minus _archive/state) finds zero user-facing `spacedock init --host`/`spacedock init:`; only comments + the two negative `if Contains("spacedock init"){fail}` guards (frontdoor_test:167, contract_test:177) remain and pass. contract.go:200 live remedy → `install`; 3 lockstep flips (contract_test 81/136/152) assert produced `Compare().Message`. AC-7: doc names `bump-calendar .claude-plugin/marketplace.json` matching the binary usage (one `<marketplace.json>` arg); marketplace.json/plugin.json untouched (DOC-ONLY honored).

### Summary

PASSED. All 7 ACs reproduced by exercising real behavior, not prose grep: AC-1 (release replace-first) and AC-6 (cobra usage) confirmed failing-first by reverting each fix; AC-4 (prose contrast) confirmed both failing-first (old bare-word) and that the entity's literal `" on Codex"` prescription genuinely breaks the real-file test at first-officer-shared-core.md:145, justifying the recorded deviation. The capitalized Codex+Claude-contrast rule correctly prevents lowercase `claude-team`/`codex exec` from self-qualifying. The cobra whitelist is root-only — a real unknown flag on a valid subcommand (`install --bogusflag`) is still rejected with exit 2, not newly swallowed. Diff scope is exactly the declared files; no collateral behavior change in cli.go/contract.go/dispatch/build.go/release.go beyond the two intentional captain-requested changes (AC-6 usage output, AC-2 push reminder). Full suite 720 passed, gofmt/vet clean.
