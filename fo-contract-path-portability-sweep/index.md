---
title: FO contract path portability sweep and hardening
status: validation
score: 0.9
source: captain request 2026-07-08 after PR #487 review found repo-specific paths in shipped FO contract
issue:
pr:
mod-block:
worktree: .worktrees/spacedock-ensign-fo-contract-path-portability-sweep
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

Replace concrete current-repo paths in shipped FO operational contracts with workflow-derived inputs and a small binary-owned classifier surface. The FO should reason from discovered workflow state, state checkout, entity path, workflow directory, and exact captain grants; unmatched write targets should default to blocked product work unless they are classified as state/process by discovered workflow metadata.

The implementation should make the operational contract portable in two places:

1. `skills/fo-write-core/SKILL.md` should stop teaching a global table of repo paths. The classifier block should describe inputs and precedence, not Spacedock's own development layout.
2. `skills/first-officer/references/claude-fo-dispatch.md` should stop rebuilding `spacedock` from the managed repository during local-main-drift recovery. A managed repo may be any commissioned project, not the Spacedock source tree.

### Contract guidance

#### FO write-core classifier

Current bad shape:

```markdown
| allowed-state | `.spacedock-state/**`; `{workflow_dir}/_archive/**` | ... |
| allowed-process | `docs/dev/README.md`; `{workflow_dir}/README.md` | ... |
| blocked-product | `cmd/**`; `internal/**`; ... `docs/dev/_mods/**` | ... |
```

Required shape:

```markdown
Before any FO-authored file write, classify every target path with `«write.classify»(target, intent, workflow_context)`.

`workflow_context` is the resolved workflow metadata: `{workflow_dir}`, `{state_checkout}` when declared, active `{entity_path}`, archive roots, worktree path when set, registered mod paths, and the exact captain grant text for this turn.

| class | source | rule |
| --- | --- | --- |
| allowed-state | resolved state/entity/archive paths | Entity frontmatter, new entity creation, archive moves, state-transition commits, and `### Feedback Cycles` under the existing state/worktree rules. |
| allowed-process | `{workflow_dir}/README.md` only | The FO may edit the workflow README it operates because that file defines process, not the product being built. |
| blocked-product | registered mods plus every target not classified as state/process | Code, tests, product docs, fixtures, release/CI files, shipped skill/agent/reference scaffolding, plugin manifests, mods, and deliverable content go through a dispatched worker. |
| override | exact-target-grant | A blocked-product target is writable only when the captain explicitly grants direct-FO editing for this exact task and target path or exact path class. |
```

Do not replace the leaked `docs/dev` paths with another concrete state spelling such as `.spacedock-state/**`. The classifier may mention `.spacedock-state` only in an explicitly labeled example outside the operational classifier, such as "example: a workflow may declare `state: .spacedock-state`."

#### Claude local-main-drift recovery

Current bad shape:

```markdown
local-main-drift ... `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`
```

Required shape:

```markdown
local-main-drift ... behind only (`drift.behind > 0 && drift.ahead == 0`): `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk}`. Continue using the already resolved `${SPACEDOCK_BIN:-spacedock}` launcher for Spacedock helper calls. Rebuild the launcher only when `{repo}` is explicitly the Spacedock source checkout selected by the operator as the launcher source; never infer that from an arbitrary managed workflow repo.
```

This keeps drift reconciliation about synchronizing the managed repo. Launcher rebuilds remain a separate Spacedock-source install/update concern.

### Riskiest mechanism decision

The riskiest mechanism is the new contractlint guard: a naive grep for `docs/dev` or `.spacedock-state` would fail useful commission/survey examples, while an over-broad allowlist would miss the two operational leaks. The spike below exercises the scanner split against synthetic snippets before implementation touches shipped contracts.

The scanner should classify markdown by local context:

- Operational regions are first-officer contracts and deferred FO skills that tell the FO what to do: classifier tables, event-loop action bullets, recovery action bullets, and imperative command lines.
- Explicit examples are allowed only when the same nearby block labels them as an example, placeholder, install hint, discovery signal, or Spacedock source-build hint.
- Red controls must prove both sides: a planted operational `docs/dev/README.md` or `docs/dev/_mods/**` rule fails, while a planted commission example saying "example: `state: .spacedock-state`" passes.
- A planted operational `cd {repo} && go build -o spacedock ./cmd/spacedock` in local-main-drift fails, while a planted Spacedock-source install hint passes.

### Riskiest mechanism spike

Ran 2026-07-08 from the repo root with a throwaway `zsh` scanner over synthetic markdown snippets. No product or scaffolding files were edited. The scanner treated operational path/rebuild patterns as leaks unless the local snippet context explicitly labeled the mention as an example, placeholder, discovery signal, or source-build install hint; local-main-drift `{repo}` rebuilds stayed hard failures.

Observed result:

```text
ok - operational docs/dev README rule -> fail
ok - operational universal state checkout rule -> fail
ok - operational local-main-drift repo rebuild -> fail
ok - explicit state checkout example -> pass
ok - source-build install hint -> pass
```

Conclusion: the contextual guard is feasible and should become the first implementation test. Seed it with these five snippets, then expand the real-file walk only after the discriminator passes so the guard cannot collapse into a tautological "string absent everywhere" check.

## Acceptance criteria

**AC-1 - The shipped FO operational contract surface has zero repo-specific path leaks.**
Verified by: a focused `internal/contractlint` scan over shipped FO contract files returns zero operational leaks for `docs/dev`, `docs/dev/_mods`, `.spacedock-state/**` universal state rules, and generic managed-repo `go build -o spacedock ./cmd/spacedock` recovery commands. The scan must report file and line for each leak and must fail on both planted historical leak shapes from #487 and #382.

**AC-2 - The write-core classifier is workflow-derived, not repo-derived.**
Verified by: `internal/contractlint/fo_write_core_mutation_gate_test.go` parses the `FO-WRITE-CLASSIFIER` block and classifies synthetic workflow paths: `/tmp/acme-flow/README.md` or `workflows/acme/README.md` is `allowed-process` only when it equals the discovered `{workflow_dir}/README.md`; `/tmp/acme-flow/.state/task/index.md` or another synthetic resolved state checkout is `allowed-state`; `docs/dev/README.md` is `blocked-product` when `docs/dev` is not the discovered workflow dir; registered `_mods` paths are `blocked-product`. The unmatched-target default remains `blocked-product`.

**AC-3 - FO product-edit guard tests prove behavior with synthetic paths.**
Verified by: `internal/ensigncycle/fo_product_edit_guard_test.go` uses non-`docs/dev` workflow and state paths in the good state/process examples, then includes a red control where the FO claims `docs/dev/README.md -> allowed-process` while the discovered workflow is synthetic; that transcript must fail unless `docs/dev` is the resolved workflow dir.

**AC-4 - Claude local-main-drift recovery is repo-generic.**
Verified by: the shipped Claude dispatch contract no longer instructs the FO to run `cd {repo} && go build -o spacedock ./cmd/spacedock` after syncing local-main-drift. A red-control contractlint fixture that reintroduces that action under `local-main-drift` fails; a Spacedock-source-only install hint outside managed-repo recovery passes.

**AC-5 - Provenance remains documented in the fix report.**
Verified by: implementation or validation report cites the two source PRs/entities: #487 / entity 17 and #382 / entity pd, and states which non-findings were intentionally left alone.

## Test plan

- Contractlint, low cost: add `internal/contractlint/fo_contract_path_portability_test.go` or equivalent. It should walk the shipped FO operational surfaces (`skills/fo-write-core/SKILL.md`, `skills/first-officer/references/*.md`, and deferred FO skill `SKILL.md` files) and scan markdown context for operational leaks. Include discriminator tests with synthetic snippets for: operational `docs/dev/README.md` fail, operational `docs/dev/_mods/**` fail, operational `.spacedock-state/**` universal state rule fail, explicitly labeled `state: .spacedock-state` example pass, survey discovery-signal pass, operational `{repo}` rebuild fail, Spacedock-source install hint pass.
- Write-core classifier test, low cost: update `internal/contractlint/fo_write_core_mutation_gate_test.go` so classification takes a synthetic `workflow_context`. Use non-`docs/dev` paths for allowed state/process and include `docs/dev/README.md` as a blocked red control when it is not the discovered workflow README.
- Ensigncycle behavior fixture, medium cost: update `internal/ensigncycle/fo_product_edit_guard_test.go` to use the same synthetic workflow/state paths in Codex and Claude transcript fixtures. Add the red transcript where the FO labels `docs/dev/README.md` as allowed-process under a synthetic workflow context and verify the guard rejects it.
- Claude dispatch contract test, low cost: add a focused scanner case for the local-main-drift bullet that rejects `cd {repo} && go build -o spacedock ./cmd/spacedock` in that action. Keep an allowed fixture for source-build install guidance so the guard distinguishes managed-repo recovery from launcher installation.
- Repo gates: run `go test ./internal/contractlint ./internal/ensigncycle`, then `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.
- Detached red audits: temporarily reintroduce the #487 write-core table paths and verify the new contractlint/ensigncycle tests fail; temporarily reintroduce the #382 local-main-drift rebuild command and verify the new guard fails. Revert audit edits before final verification.

## Stage Report: ideation

- DONE: Turn the sweep log into a workflow-generic, behavior-first fix plan with concrete before/after guidance for the shipped FO contracts.
  Added contract guidance for `fo-write-core` and Claude local-main-drift with current bad shapes and required portable shapes.
- DONE: Strengthen acceptance criteria so they measure portability and fail on both known leaks: #487 write-core path overfit and #382 Claude local-main-drift source rebuild.
  AC-1 through AC-4 now require zero operational leaks, synthetic workflow-path behavior, and red controls for both historical leak classes.
- DONE: Record the riskiest-mechanism decision and a non-tautological test plan, including red controls for reintroducing docs/dev path rules and arbitrary-repo go-build recovery.
  Added a riskiest-mechanism section and a test plan centered on contextual contractlint plus behavior fixtures, not prose self-assertions.

### Summary

Ideation converted the existing sweep notes into an implementation-ready portability plan. The plan preserves provenance for #487/entity 17 and #382/entity pd, distinguishes operational contract leaks from labeled examples, and gives implementation concrete red controls for both known regressions.

## Stage Report: ideation (cycle 2)

- DONE: Do not edit product/scaffolding files. Entity-body updates are fine.
  The repair changed only this state-checkout entity file.
- DONE: Add a concise `Riskiest mechanism spike` or equivalent section to `docs/dev/.spacedock-state/fo-contract-path-portability-sweep/index.md`.
  Added `### Riskiest mechanism spike` after the mechanism decision.
- DONE: Prefer a lightweight throwaway spike against synthetic snippets: demonstrate that a contextual scanner can fail operational `docs/dev/README.md`, operational `.spacedock-state/**`, and local-main-drift `{repo}` rebuild snippets while passing an explicitly labeled `state: .spacedock-state` example and a source-build install hint. If you cannot run a spike, record a defensible no-spike rationale with concrete evidence from existing mechanisms.
  Ran the throwaway scanner; all five synthetic cases produced the expected pass/fail result.
- DONE: Update the ideation stage report with DONE/SKIPPED/FAILED accounting for this repair.
  Appended this cycle-2 ideation report.
- DONE: Commit the state-checkout update if required by the workflow.
  The state update is limited to this entity path; the repair commit records it path-scoped.

### Summary

The pre-gate repair exercised the riskiest contextual-scanner mechanism instead of deferring it. The recorded spike gives implementation a concrete first test that fails the known operational leaks while allowing explicitly labeled examples and install guidance.

## Stage Report: implementation

- DONE: Implement contextual contractlint guards with synthetic discriminator fixtures that fail the #487 write-core path leaks and #382 local-main-drift repo rebuild while allowing explicitly labeled examples/install hints.
  Code commit `23317ad` adds `fo_contract_path_portability_test.go`; detached red audits failed on #487 write-core paths and #382 local-main-drift rebuild before restoration.
- DONE: Update shipped FO operational contracts so write-core classification is workflow-derived and Claude local-main-drift recovery is repo-generic, with no operational docs/dev, universal .spacedock-state, or arbitrary-repo go-build leaks.
  `skills/fo-write-core/SKILL.md` now classifies from workflow context; `claude-fo-dispatch.md` syncs local-main-drift without rebuilding from `{repo}`.
- DONE: Update behavior fixtures/tests to use synthetic workflow/state paths, run focused plus repo gates, and include detached red-audit evidence before final implementation report.
  `go test ./internal/contractlint ./internal/ensigncycle`, `go test ./...`, and `go test ./... -race` passed after `gofmt -w ./cmd ./internal`.

### Summary

Implementation hardens the shipped FO contracts against the two logged provenance leaks: #487 / entity 17 and #382 / entity pd. The non-findings from the sweep were intentionally left alone, while the new guard preserves labeled examples, discovery signals, and source-build install hints.
