---
title: FO contract path portability sweep and hardening
status: ideation
score: 0.9
source: captain request 2026-07-08 after PR #487 review found repo-specific paths in shipped FO contract
issue:
pr:
mod-block:
worktree:
started: 2026-07-08T12:45:36Z
completed:
verdict:
id: fafb8xb7ptn2dx5k2c423b50
---

## Problem

A shipped first-officer contract must be workflow-generic. A review of PR #487 found concrete Spacedock development workflow paths inside `skills/fo-write-core/SKILL.md`. A full shipped-instruction sweep found one additional older operational assumption in Claude reconcile recovery.

## Sweep log

### Finding A - #487 FO write-core classifier overfit to this repo

- Source PR: #487, `FO product-edit guard loads write-core before mutation`, merged 2026-07-08T07:32:59Z, merge commit `821e5b08742696b1bde6e6304024086bce03b5fb`.
- Original entity: `17` / `fo-write-core-product-edit-guard`, archived with `pr=pr-merge:487`, `verdict=passed`, completed 2026-07-08T07:33:40Z.
- Offending shipped lines: `skills/fo-write-core/SKILL.md` classifier block.
- Bad concrete paths / shapes:
  - `.spacedock-state/**` as an allowed-state rule, which assumes one state checkout spelling.
  - `docs/dev/README.md` as allowed-process, which assumes this repo's workflow directory.
  - `docs/dev/_mods/**` as blocked-product, which assumes this repo's mod directory.
  - Repo-shape product globs such as `docs/specs/**`, `docs/roadmap/**`, `docs/site/**`, and `fixtures/**`, which are Spacedock source-repo examples rather than a portable FO rule.
- Test leakage: `internal/contractlint/fo_write_core_mutation_gate_test.go` and `internal/ensigncycle/fo_product_edit_guard_test.go` encode the same concrete examples, so validation protected the leak.

### Finding B - pre-existing Claude reconcile recovery rebuilds arbitrary repos as Spacedock

- Source PR: #382, `Reconcile drift output carries descriptive class names instead of A-E letters`, merged 2026-06-15T15:44:54Z, merge commit `69bcd93fdbf063ab6f15551c376e287388135986`.
- Original entity: `pd` / `reconcile-drift-class-names`, archived as `reconcile-drift-class-names`, `pr=#382`, `verdict=PASSED`, completed 2026-06-15T15:45:51Z.
- Offending shipped line: `skills/first-officer/references/claude-fo-dispatch.md` local-main-drift action.
- Bad concrete command: `cd {repo} && go build -o spacedock ./cmd/spacedock` after fast-forwarding `{repo}`. That assumes the managed repo is the Spacedock source repo. A workflow embedded in another codebase should not rebuild `spacedock` from that repo's `./cmd/spacedock` path.

### Non-findings from this sweep

The following shipped-instruction hits were reviewed and are not the same class of bug:

- `skills/commission/SKILL.md` uses `docs/dev` and `.spacedock-state` as explicitly labeled examples/placeholders for scaffolding.
- `skills/first-officer/references/first-officer-shared-core.md` uses `go build -o spacedock ./cmd/spacedock` only as a Spacedock source-build install hint when the launcher binary is absent.
- `skills/first-officer/references/fo-dispatch-core.md` and `skills/ensign/references/ensign-shared-core.md` use `state: .spacedock-state` as an example of a README-declared split-root checkout and then defer to the dispatch-provided entity path.
- `skills/survey/SKILL.md` probes `.spacedock-state`, `docs/**/.spacedock-state`, and `_mods` as discovery signals, not as an operational write rule.
- `.worktrees/{worker_key}` / `spacedock-ensign` are workflow worker-key conventions, not this repo's `docs/dev` workflow path.

## Proposed approach

Replace concrete current-repo paths in shipped FO operational contracts with workflow-derived placeholders or binary-owned classification. The FO should reason from discovered workflow state, state checkout, entity path, workflow directory, and exact captain grants; unmatched write targets should default to blocked product work unless they are classified as state/process by discovered workflow metadata.

## Acceptance criteria

**AC-1 - Shipped FO write-core classifier is workflow-generic.**
Verified by: `skills/fo-write-core/SKILL.md` no longer names `docs/dev`, `docs/dev/_mods`, or `.spacedock-state/**` as universal rules. It uses `{workflow_dir}`, discovered state checkout/entity paths, or a binary-owned classifier surface. The contract explicitly says unmatched targets default to blocked-product.

**AC-2 - Contractlint catches repo-specific operational path leaks.**
Verified by: a focused contractlint test scanning shipped `SKILL.md` and `references/*.md` files. It must fail on operational mentions of this repo's dev workflow paths (`docs/dev`, `docs/dev/_mods`) or Spacedock source rebuild commands in generic managed-repo recovery, while allowing explicitly marked examples.

**AC-3 - The #487 tests stop protecting the bad examples.**
Verified by: `internal/contractlint/fo_write_core_mutation_gate_test.go` and `internal/ensigncycle/fo_product_edit_guard_test.go` use synthetic non-`docs/dev` workflow/state paths and include a red control proving `docs/dev/README.md` is not specially allowed unless it is the discovered `{workflow_dir}/README.md`.

**AC-4 - Claude local-main-drift recovery is repo-generic.**
Verified by: `skills/first-officer/references/claude-fo-dispatch.md` no longer instructs the FO to run `go build -o spacedock ./cmd/spacedock` inside arbitrary `{repo}` after drift. The new behavior either syncs only, uses the already resolved launcher invariant, or names a Spacedock-source-only rebuild condition explicitly.

**AC-5 - Provenance remains documented in the fix report.**
Verified by: implementation or validation report cites the two source PRs/entities: #487 / entity 17 and #382 / entity pd, and states which non-findings were intentionally left alone.

## Test plan

- Run focused contractlint for the new path-leak guard and the FO write-core mutation gate.
- Run focused ensigncycle product-edit guard tests after replacing concrete fixtures with synthetic workflow/state paths.
- Run `go test ./...` and `go test ./... -race`.
- Detached audit: temporarily reintroduce `docs/dev/README.md` or `docs/dev/_mods/**` into the shipped write-core classifier and verify the new guard fails.
- Detached audit: temporarily reintroduce `cd {repo} && go build -o spacedock ./cmd/spacedock` into the Claude local-main-drift action and verify the guard fails unless the line is explicitly marked as a Spacedock source-build install hint.
