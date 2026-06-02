---
id: 1xk9bz4fr7qefgcqz6stzpkk
title: 0.19.3 code-quality cleanups — release regex clobber, statecommit test phrase, StandingTeammate dedup, prose-marker brittleness, init→install rename drift, unknown-subcommand silent exit
status: ideation
source: sprint-end antipattern reviews (2026-06-01) — 0.19.3 minor-findings bucket (all Minor)
started: 2026-06-02T04:57:30Z
completed:
verdict:
score: "0.22"
worktree:
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

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — release multi-version clobber fixed + guarded.** A version-stamp over a file containing more than one `version` key rewrites only the intended key.
Verified by: a multi-version fixture test in `internal/release` that fails against the current `ReplaceAll` and passes after the targeted replace.

**AC-2 — the remaining cleanups land with no behavior change.** StandingTeammate is one struct + one mapper; the statecommit phrase assertion and push-on-commit line are present; prose markers match phrases.
Verified by: `go test ./...` green; the dedup compiles against a single definition; the new assertions present.

**AC-3 — no `spacedock init` remains as a user-facing command callout.** Code and docs reference `spacedock install` (or the live subcommand) wherever they previously said `spacedock init`; the version-mismatch remedy emits a runnable command.
Verified by: a grep over `internal/`, `README.md`, `docs/` finds no `spacedock init --host` / `spacedock init:` user-facing string; a test on the contract version-mismatch message asserts it names `install`.

**AC-4 — unknown subcommand always prints usage.** `spacedock <unknown>` and `spacedock <unknown> --flag` both print the usage block to stderr and exit non-zero — never a silent exit 2.
Verified by: a CLI test that runs an unknown command carrying a flag and asserts non-empty usage output on stderr + non-zero exit.

## Notes
- Touches `internal/release` + `internal/dispatch` + `internal/cli` + `internal/contract` + docs + tests — NOT the `internal/status` serialized lane, so its implementation worktree can run alongside status-lane work. `internal/cli/cli.go` (root command) is a hot single-writer file — coordinate the implementation worktree with any other `cli.go` writer. The StandingTeammate dedup traces to the `zs` claude-runtime-segregation split (upstream-flow note).
- The init→install drift + unknown-subcommand-silent-exit were folded in by the captain (2026-06-02) from a live install-path test; the sibling FO-contract guidance gap (binary-*absent* startup routes to the missing-binary `doctor`) is tracked separately as `binary-absent-fo-bootstrap`.
