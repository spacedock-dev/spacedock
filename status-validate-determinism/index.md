---
title: status --validate non-deterministic on untracked entity files (root cause of escaped blank-id)
status: implementation
source: "Shaping FO (2026-06-20): investigating why a Commander-filed entity (pi-devoverride-package-ok) shipped with a blank id through its whole lifecycle (create → pr-pending → mod-block → terminalize → archive) and only surfaced when a fresh boot ran --validate. A pre-commit hook running `spacedock status --validate` could not reliably block the same blank-id entity: the same file on disk returns `Error: missing required id` on some invocations and `VALID` on others, across both clean-shell and git-hook (GIT_DIR/GIT_INDEX_FILE set) environments. The non-determinism is in the status tool's entity scan, not the hook."
score:
started: 2026-06-20T02:39:48Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-status-validate-determinism
issue:
sprint:
sprint-readiness:
id: a1a9k6tqnpj7292z80xxwkf2
---

# status --validate non-deterministic on untracked entity files

## End value

`spacedock status --validate` never returns `VALID` for a misdirected or empty validation target, and when pointed at the real workflow definition dir it scans every on-disk entity file (tracked or not, staged or not, committed or not). A hook or human caller can trust validate to either report the corrupt blank-id entity or fail closed with a root-resolution error.

## Problem — root cause pinned by ideation spike

The Commander filed `pi-devoverride-package-ok` by hand-writing the file (bypassing `spacedock new`'s id-minting), so it shipped with `id:` blank from birth (state commit `74af536f`) through archive (`2784813b`) before later id backfills (`210182fe`, `882f9db9`). The blank id should have been blocked by validation before any lifecycle transition or commit.

The apparent "same file sometimes errors, sometimes VALID" behavior is deterministic root misresolution, not an `os.ReadDir` race. In a throwaway split-root fixture with `docs/dev/README.md` declaring `state: .spacedock-state` and an untracked blank-id folder entity at `.spacedock-state/zprobe-hook-test/index.md`:

- From the code repo root, `spacedock status --workflow-dir docs/dev --validate` returns `Error: missing required id` in both clean and git-hook envs.
- From the state checkout root, the same relative argv resolves `docs/dev` under `.spacedock-state/`, scans an empty/non-workflow target, and returns `VALID` in both clean and git-hook envs.
- From `docs/dev`, the same relative argv resolves `docs/dev/docs/dev`, scans an empty/non-workflow target, and returns `VALID`.
- With an absolute workflow definition dir, both `spacedock 0.22.0` and the current dev binary return `missing required id` 20/20 runs for clean env and 20/20 runs for hook-shaped `GIT_DIR`/`GIT_INDEX_FILE`.

The fail-open path is `dispatch`/`resolveRoots`: an explicit `--workflow-dir` is used verbatim; if it resolves to a directory with no commissioned README, `resolveRoots` falls back to inline state and `validateWorkflow` scans zero entities, then prints `VALID`. `discoverEntityFiles` itself is filesystem-only and deterministic once it receives the correct entity root.

## Approach

- Make explicit `--validate` fail closed before `validateWorkflow` when `roots.definitionDir` is not a commissioned workflow definition dir, or when the README declares split-root state and `roots.entityDir` is absent/not a directory. Keep the existing direct-state-checkout diagnostic when `stateCheckoutParent` recognizes it.
- Scope the behavior change to validation/guard surfaces, not the default table read: the existing explicit-empty-dir table compatibility can stay unchanged, but validation must never certify an empty non-workflow target as `VALID`.
- Preserve entity enumeration as pure on-disk discovery. Do not introduce git tracked/staged filters; flat `{slug}.md` and folder `{slug}/index.md` entities remain eligible based on files, opening fence, reserved-dir rules, and flat/folder conflict handling.
- Add regression tests around the command surface, not only helper functions: run `NativeRunner` repeatedly against a split-root fixture with untracked blank-id flat and folder entities, under clean env and hook-shaped env, from code-root, state-root, and definition-dir CWDs.

## Acceptance criteria

**AC-1 — `--validate` fails closed on a misresolved workflow root.**  
Verified by: a Go command-level test that runs `status --workflow-dir docs/dev --validate` from inside the state checkout and from inside `docs/dev`; both return non-zero with a named root-resolution error, never `VALID`.

**AC-2 — `--validate` deterministically errors on blank-id on-disk entities when pointed at the workflow definition dir.**  
Verified by: a Go command-level test that writes untracked blank-id flat and folder-form entities to the split-root state checkout, runs validation N times with an absolute definition dir, and asserts every run reports `missing required id` for the corresponding slug.

**AC-3 — Clean-shell and git-hook environments behave identically.**  
Verified by: the AC-1 and AC-2 tests run each command once with `GIT_DIR`, `GIT_INDEX_FILE`, and related hook env unset, and once with `GIT_DIR`/`GIT_INDEX_FILE` set to the linked state worktree's git paths.

**AC-4 — Tracked, staged, and committed entities keep existing behavior.**  
Verified by: existing validate suites stay green; an added fixture validates a good tracked entity, a good staged entity, and a good committed entity, while a bad untracked sibling still fails validation.

## Out of scope

- The pre-commit hook auto-install at `spacedock state init` — separate task (`state-init-precommit-hook`); depends on this task landing first (the hook is only reliable once validate is deterministic).
- The `status: implementation` (not `done`) latent inconsistency on archived `pi-devoverride-package-ok` — separate data-quality cleanup; not the root cause.
- Changing `validate`'s error/warning classification — only the entity-enumeration determinism.

## Test plan

- Go tests in `internal/status/`: split-root fixture with linked state worktree, blank-id flat + folder entities, absolute definition dir, wrong relative `docs/dev` from state/definition CWD, clean env, and hook-shaped env.
- Regression guard for the fail-open: wrong relative root must exit non-zero and must not print `VALID`.
- Regression guard for entity scanning: absolute definition dir must report both blank-id files every run, regardless of git tracking state.
- Existing focused status tests plus `go test ./...`; this changes a high-stakes status guard path, so validation needs a detached adversarial audit. The audit should temporarily remove the new fail-closed root check and confirm the new misresolved-root test turns red, then restore.
- This is a `status` mutation/guard-path change — high-stakes surface per the dev-workflow proof policy. Detached adversarial audit at validation.

## Documentation diff for implementation

The implementation changes user-visible `--validate` behavior for invalid workflow roots. Apply a small docs update alongside the code:

```diff
diff --git a/docs/dev/README.md b/docs/dev/README.md
@@
-| Status validator | `spacedock status --workflow-dir docs/dev --validate` | Spacedock entity-contract validation |
+| Status validator | `spacedock status --workflow-dir docs/dev --validate` from the repo root, or pass an absolute workflow definition dir | Spacedock entity-contract validation; fails closed if `--workflow-dir` does not resolve to a commissioned workflow |
```

## Related

- `pi-devoverride-package-ok` (archived, `id` now `18s156v0txff12fv61yqptk2`) — the escaped-blank-id instance that surfaced this.
- `internal/status/validate.go` (`validateWorkflow`, `validateSDB32`), `internal/status/discover.go` (`scanEntitiesActive`, `discoverEntityFiles`, `loadActiveEntityFields`) — the scan path.
- `state-init-precommit-hook` (sibling followup) — the pre-commit hook auto-install; depends on this task.
- The pre-commit hook currently installed at `.git/hooks/pre-commit` (worktree-shared) — documents this limitation in its header.

## Stage Report: ideation

- DONE: Pin the root cause hypothesis and minimal deterministic reproduction for blank-id entity scan nondeterminism, including clean shell and git-hook env cases.
  Evidence: disposable split-root fixture showed relative `docs/dev` validates the correct root from repo cwd but fails open to `VALID` from state/definition cwd; absolute definition dir reported missing-id 20/20 clean and 20/20 hook-env runs.
- DONE: Produce a concrete implementation approach that keeps validate scanning every on-disk flat or folder-form entity without regressing tracked/staged/committed entity handling.
  Evidence: approach now scopes the fix to fail-closed validation root checks while preserving filesystem-only entity enumeration and adding tracked/staged/committed regression coverage.
- DONE: Tighten acceptance criteria and test plan so AC-1 through AC-3 are proven by behavior, including the required high-stakes validation/audit implications.
  Evidence: ACs now bind to command-level Go tests for wrong-root fail-closed behavior, blank-id flat/folder scanning, clean-vs-hook env parity, and high-stakes adversarial audit.

### Summary

Ideation corrected the seed hypothesis: the observed flip is deterministic root misresolution for a relative `--workflow-dir`, not nondeterministic `discoverEntityFiles` scanning. The task is now shaped around fail-closed `--validate` root handling plus command-level regression tests that prove untracked on-disk entities are still scanned when the workflow definition dir is correct.
