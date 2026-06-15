---
id: r0n903prs88q29vjkmcmtw2x
title: FO contract documents the from-root `spacedock new` invocation (--workflow-dir) so filing an entity needs no trial-and-error
status: ideation
source: "captain (2026-06-14, this session) — filing an entity from the project root hit no-Spacedock-workflow-here / pass --workflow-dir. FO Write Scope (first-officer-shared-core.md:171) and Filing New Entities (claude-first-officer-runtime.md:43) both document `spacedock new <slug> [--folder] [--id-seed S --id-actor A]` WITHOUT --workflow-dir, but the Working Directory rule keeps the FO at the project root where `new` cannot auto-discover the workflow. `new --help` returns the general menu, not per-command usage, so the FO must trial-and-error the flags. Same fo-efficiency / friction-reduction class as the 0.20.4 structured-reads and 0203 work."
started: 2026-06-15T00:32:39Z
completed:
verdict:
score: 0.30
worktree:
issue:
sprint: 0203-fo-efficiency
---

The documented `spacedock new` invocation is incomplete for the FO's standing position. An FO at the project root (per the Working Directory rule) must pass `--workflow-dir`, but the contract's `new` examples omit it, so the first filing attempt fails. Close the gap so a fresh FO files correctly on the first try.

## Problem

Two contract sites document filing:
- `first-officer-shared-core.md:171` (FO Write Scope): "seed task creation via `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < stub`".
- `claude-first-officer-runtime.md:43` (Filing New Entities): "Use `spacedock new <slug> [--folder] [--id-seed S --id-actor A]` via Bash".

Neither shows `--workflow-dir {workflow_dir}`. The Working Directory rule ("Stay at the project root") means the FO runs `new` from root, where the binary errors `no Spacedock workflow here — pass --workflow-dir`. `new --help` prints the top-level command menu rather than per-command usage, so the only way to learn the required flag is to fail and read stderr.

## Proposed approach

{Ideation fills. Keep it simple. Candidate direction to evaluate: (a) update both contract sites to show the from-root invocation including `--workflow-dir {workflow_dir}`; (b) make `spacedock new --help` (and other subcommand `--help`) print per-command usage with the full flag surface.}

## Out of scope

{Ideation fills. Likely: generalizing to every subcommand's contract documentation — scope to `new` unless ideation finds the same gap is cheap to close broadly.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 — `spacedock new --help` prints the command's own usage and full flag surface.** Running it shows the required/optional flags (`--workflow-dir`, `--folder`, `--id-seed`, `--id-actor`) rather than the top-level menu.
Verified by: {a CLI test asserting `new --help` stdout contains the per-command usage and the flag names; ideation pins it.}

**AC-2 — The contract's documented `new` invocation matches what the binary actually requires from the FO's standing directory.** The from-root example includes `--workflow-dir`.
Verified by: {a check binding the contract-documented flags to the binary's actual accepted-flag set (two independent values that can diverge — legitimate, not prose-grep); ideation pins the form.}

## Test plan

{Ideation fills.}
