---
id: 8e2053706c2c77xdm696dfr6
title: dispatch build emits Agent name >64 chars for long slugs — Agent() dispatch fails
status: backlog
source: github#366 (captain intake 2026-06-14)
started:
completed:
verdict:
score: "0.40"
worktree:
issue: "#366"
---

`spacedock dispatch build` constructs the worker `name` as `{worker_key}-{slug}-{stage}` with no length cap, but Claude Code's `Agent` tool enforces `name` matching `^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$` (max 64 chars). For a long slug the emitted `name` exceeds 64 chars, so forwarding the helper output verbatim to `Agent()` — the contract's MANDATORY dispatch path — fails with `InputValidationError` on `name`, and the dispatch cannot proceed.

## Problem

- Name construction (`internal/dispatch/build.go`, `derivedName := fmt.Sprintf("%s-%s-%s", workerKey, slug, stage)`) has no 64-char cap.
- Observed live: `pr-merge-mod-base-branch-post-flip` → name 60 chars, dispatched fine; `dispatch-reconcile-deconflate-repo-hygiene` → name 68 chars, `Agent(name=...)` rejected with `InputValidationError … (max 64 chars)`.
- **Decomposition is load-bearing.** A hand-shortened or naively-truncated `name` no longer decomposes cleanly to `(worker_key, slug, stage)`. That decomposition drives supersede-shutdown, terminal-team-teardown cohort derivation, and reconcile A/B (lingering/superseded) classification. A name that does not map back to its entity risks a false Class-A "lingering" flag (a spurious shutdown against a live worker) or a missed cohort teardown.

## Notes (suggested-fix direction from github#366 — ideation owns the actual design)

- Cap the generated `name` at 64 chars deterministically inside `dispatch build`, preserving a stable, decomposable identity — e.g. truncate the slug component to fit while keeping the `worker_key` prefix + `stage` suffix, or append a short stable hash of the full slug so distinct long slugs do not collide.
- Whatever the form, the reconcile/teardown name-decomposition must round-trip back to the entity (e.g. resolve via a stored short-id rather than string-splitting the full slug).
- `name` and `dispatch_file_path` need not match — only `name` needs capping.

## Proposed approach

{Ideation fills this in.}

## Out of scope

{Ideation fills this in — e.g. whether to introduce stored short-id resolution vs pure deterministic truncation+hash.}

## Acceptance criteria

{Ideation fills this in — each AC proven by a Go test over real `dispatch build` output (name ≤64 chars for a long slug AND the name resolves back to the correct entity/cohort), never a prose grep over an instruction file.}

## Test plan

{Ideation fills this in.}
