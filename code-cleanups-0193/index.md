---
id: 1xk9bz4fr7qefgcqz6stzpkk
title: 0.19.3 code-quality cleanups — release regex clobber, statecommit test phrase, StandingTeammate dedup, prose-marker brittleness
status: backlog
source: sprint-end antipattern reviews (2026-06-01) — 0.19.3 minor-findings bucket (all Minor)
started:
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

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — release multi-version clobber fixed + guarded.** A version-stamp over a file containing more than one `version` key rewrites only the intended key.
Verified by: a multi-version fixture test in `internal/release` that fails against the current `ReplaceAll` and passes after the targeted replace.

**AC-2 — the remaining cleanups land with no behavior change.** StandingTeammate is one struct + one mapper; the statecommit phrase assertion and push-on-commit line are present; prose markers match phrases.
Verified by: `go test ./...` green; the dedup compiles against a single definition; the new assertions present.

## Notes
- Touches `internal/release` + `internal/dispatch` + tests — NOT the `internal/status` serialized lane, so its implementation worktree can run alongside status-lane work. The StandingTeammate dedup traces to the `zs` claude-runtime-segregation split (upstream-flow note).
