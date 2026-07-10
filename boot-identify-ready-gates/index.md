---
title: Expose ready-gate entities in boot identify JSON
status: implementation
score: 0.9
source: PR #493 local-live m3 investigation on 2026-07-10. The default Claude greeting omitted the already-gated Gate Check entity because status --boot --identify --json returned only dispatchable entities; dispatchAnalysis intentionally suppresses current gate stages, so the authoritative boot record made the shipped greet requirement impossible to satisfy.
id: 8n55etrw9wj10jfejdq5f1s8
worktree: .worktrees/spacedock-ensign-boot-identify-ready-gates
---

## Problem

The interactive first-officer contract requires the greeting to name ready gates without rendering their reviews. `status --boot --identify --json` exposes only `dispatchable`, which intentionally excludes entities whose current stage is a gate. In the real PR #493 m3 shallow-boot record, `gate-check.md` existed at `review`, but the boot JSON named only `merged-pr` in `pr_state` and `dispatchable`. `discovery` named the workflow path, and `stages` described the taxonomy; neither field identified `gate-check`. The model therefore had no authoritative gate identity to name.

## Current behavior and ordering

`scanEntitiesActive` discovers active entities in lexical slug order and parses frontmatter once. `sortDefault` provides the status order: declared stage order, then score descending with empty scores last, then the stable discovery order. `dispatchAnalysis` instead uses score order for dispatch candidates and suppresses a known current stage in this precedence: terminal, gate, worktree-set, concurrency-full. A gate is therefore absent from `dispatchable` by design.

`gatherBoot` currently retains only the computed dispatchable rows. `bootJSON` emits those rows as fixed `id`, `slug`, `current`, `next`, and `worktree` objects; identify mode then appends `discovery` and `stages`. The existing identify tests pin this key order, local-only behavior, and zero/one/many discovery. The JSON structure tests pin ordinary boot keys and the equality of boot and `--next` dispatchable rows. No existing boot field can recover an active entity's identity from the stage taxonomy alone.

The current focused baseline is green: `go test -count=1 -v ./internal/status -run 'TestBoot(Identify|JSON)'` passed all 10 selected boot identify/JSON tests. This task changes their expected schema; it does not repair a failing status test.

## Approaches considered

1. **Recommended: identify-only `ready_gates`.** Append a compact array after `stages` only for `--boot --identify --json`. This gives the greet the missing identities while preserving ordinary boot bytes and `dispatchable` semantics.
2. **Emit `ready_gates` on every boot JSON record.** This gives non-identify consumers the same data, but the greet is the only stated consumer. It needlessly widens the stable ordinary-boot schema and its fixtures.
3. **Put gate rows in `dispatchable` or expose every active entity.** Gate rows would contradict `--next` and break the fixed dispatchable contract. A complete active-entity inventory would expose more frontmatter and schema than the greet needs.

## Proposed approach

Add an identify-only `ready_gates` array with three string fields per row:

```json
"ready_gates":[{"id":"gate-check","slug":"gate-check","current":"review"}]
```

Select active entities whose current status resolves to a declared stage with `gate: true` and `terminal: false`. Exclude terminal stages even if misconfigured with `gate: true`, unknown statuses, archived entities, and ordinary non-gate stages. A gate's worktree value does not affect readiness because `dispatchAnalysis` already gives gate suppression precedence over worktree suppression.

Order the selected rows with the existing status order: declared stage order first, score descending with empty scores last, and lexical slug order as the stable tie-break inherited from discovery. Each row emits keys in `id`, `slug`, `current` order. Emit `ready_gates: []` when identify mode finds no current gates.

Append `ready_gates` after `stages`. Every existing identify key retains its relative order. Ordinary `--boot --json` omits the field and remains byte-compatible; identify JSON gains one additive field. `dispatchable` remains byte-for-byte unchanged and continues to match `--next --json`.

No spike is needed. The design reuses four mechanisms already exercised by repository tests: active frontmatter scanning, ordered stage metadata, `sortDefault`, and the insertion-ordered JSON builder's non-nil empty arrays. It adds no entity-body read, filesystem walk, prompt cue, or state mutation.

## Implementation scope

- `internal/status/format.go`: add the current-gate selector and reuse `sortDefault` for row order.
- `internal/status/boot.go`: carry ready-gate entities in `bootData` and populate them only in identify mode.
- `internal/status/json_commands.go`: render the fixed three-field row and append `ready_gates` after `stages`.
- `internal/status/boot_identify_test.go`: add multiple/one/zero gate controls, exact row/key order, terminal/unknown exclusion, and the m3-shaped `gate-check` fixture.
- `internal/status/json_boot_test.go`: pin ordinary boot omission and the unchanged dispatchable raw value.
- `docs/site/reference/command-reference.md`: document the new identify field. Change “folds in the stage taxonomy, and reports the boot sections” to “folds in the stage taxonomy and active ready-gate identities (`id`, `slug`, `current`), and reports the boot sections.”

No skill, first-officer prompt, shared live fixture, or scenario assertion changes in this task. PR #493's m3 integration tree consumes the shipped field afterward and owns the live default-startup oracle update.

## Out of scope

- Rendering gate reviews during boot.
- Changing dispatchability, stage transitions, or gate ownership.
- Adding recursive entity reads, title/body projection, or model-facing prompt cues.
- Repairing or weakening the separate m3 scenario oracle in this change.

## Acceptance criteria

- **AC-1 (value — every ready gate is visible at boot).** For an offline identify fixture with three active current-gate entities across two declared gate stages, `ready_gates` contains exactly 3/3 rows and 0 non-gate rows, each with only `id`, `slug`, and `current`; rows follow status order. The m3-shaped case contains `{"id":"gate-check","slug":"gate-check","current":"review"}`. The current baseline exposes 0/3 because the field does not exist. *Verified by:* a table-driven CLI test in `boot_identify_test.go` that runs the real native status path and asserts the exact decoded rows plus raw row-key order.
- **AC-2 (dispatchability remains authoritative and separate).** Adding any number of current gates changes neither the raw `dispatchable` array nor the `--next --json` result; current gates remain excluded. Terminal, terminal-plus-gate, unknown-stage, archived, and ordinary non-gate entities never enter `ready_gates`. *Verified by:* the same mixed fixture compares the exact raw dispatchable value with an independently specified expected value and with `--next --json`, then asserts every excluded slug is absent from `ready_gates`.
- **AC-3 (schema and ordering are backward-compatible).** Identify mode appends `ready_gates` after `stages`, emits `ready_gates: []` for zero gates, and preserves every prior key's relative order. Ordinary `--boot --json` omits `ready_gates` and remains byte-compatible. *Verified by:* focused key-order and zero-gate assertions in `boot_identify_test.go`, a negative ordinary-boot assertion in `json_boot_test.go`, and the existing JSON structure/dispatchable tests.
- **AC-4 (m3 greet receives authoritative gate identity).** On the coupled PR #493 integration tree, the real default-startup Claude shallow-boot run's single boot record contains `gate-check` at `review`, and the final light greeting names canonical `gate-check` or “Gate Check” without a post-boot entity-body read, broad search, custom prompt, or rendered `Gate review:` block. The current m3 boot record contains no `gate-check` identity and required later filesystem/status recovery. *Verified by:* re-running m3's live scenario and inspecting the boot tool result, pre-greet tool sequence, and final message; this is the end-value handoff, not a claim that this status-only task fixes the separate oracle by itself.

## Test plan

- Start with the CLI-level mixed fixture in `internal/status/boot_identify_test.go`. It costs milliseconds and exercises active discovery, stage resolution, ordering, boot gathering, and JSON rendering through the real native command path. Make it red first on the missing `ready_gates` field.
- Add one-gate m3-shape and zero-gate cases, plus the ordinary-boot omission pin in `internal/status/json_boot_test.go`. Keep the existing `TestBootJSONDispatchableMirrorsNext` green and compare the mixed fixture's raw dispatchable bytes against an independent literal.
- Run `go test ./internal/status -run 'TestBoot(Identify|JSON)'`, `go test ./internal/status ./internal/cli`, `go test ./...`, `go test ./... -race`, and `go test -tags live -run '^$' ./internal/ensigncycle/...`.
- After merging the status change into the coupled PR #493 integration tree, run the real Claude default-startup shallow-boot scenario once for AC-4. This is the only model-costly proof; the status schema and compatibility evidence remain offline.

## Stage Report: ideation

- DONE: Audit the current boot JSON schema, dispatchAnalysis gate suppression, status stage/entity ordering, and existing boot identify tests. Confirm the observed m3 record truly cannot expose `gate-check` through any current field.
  The code audit found gate suppression before worktree/concurrency, status ordering by declared stage then score then slug, and no boot inventory beyond dispatchable; the real m3 record names only `merged-pr`, while `gate-check.md` independently records `status: review`.
- DONE: Choose the smallest stable schema that lets the greet name ready gates without changing dispatchability. Specify exact fields, deterministic ordering, zero-gate behavior, compatibility impact, and exact implementation/test files; avoid entity-body reads or prompt changes.
  The recommended identify-only `ready_gates` field carries ordered `id`/`slug`/`current` rows after `stages`, emits `[]` at zero, leaves ordinary boot and dispatchable bytes unchanged, and scopes implementation to five Go files plus the command reference.
- DONE: Refine AC-1 through AC-4 with independent offline evidence and the m3 live end-value handoff. Append a complete ideation Stage Report and recommendation; do not change product code or claim unrun tests passed.
  AC-1 measures 3/3 ready gates against the current 0/3 baseline; AC-2/AC-3 pin dispatch and schema compatibility offline; AC-4 hands the authoritative row to m3's real live greet without claiming its separate oracle is fixed here. Recommend approval for implementation.

### Summary

The design adds one identify-only, append-only `ready_gates` array and reuses existing active-entity parsing and status ordering. It supplies the missing m3 gate identity while preserving ordinary boot output, dispatchability, prompts, and entity-body read boundaries.
