---
title: Document entity authoring for the first officer — contract "Filing an entity" stanza + `spacedock new --help` body template
status: backlog
score: ""
source: FO-efficiency friction (0203) — filing the live-test-stream-tee entity, the FO had to probe `spacedock new --help` plus new.go / state_new_test.go because neither the operating contract nor the help documents the entity-authoring contract
priority: medium
id: ah6j5r36gp23egbsr5ddrxc5
---

## Problem

When the first officer is asked to file or shape an entity, nothing tells it HOW to author one. The operating contract documents the dispatch loop, gates, and the `status` / `dispatch` / `boot` / `merge` verbs — but not entity creation. And `spacedock new --help` shows the invocation and flags, not the body contract. So the FO reverse-engineers it from `--help` + source.

Specifically the FO must discover on its own:

- the verb is `spacedock new [--folder] SLUG < body`;
- the **id is minted by the tool** — never hand-write `id` in the body;
- the stdin **body carries the frontmatter** (`title`, `status`, `source`, `priority`) and `new` injects the minted `id`;
- the body should follow the value-AC rule (`docs/dev/README.md`).

Observed while filing `live-test-stream-tee-ci-stdout-noise`: the FO read `new --help`, then `internal/cli/cli.go`, `internal/status/new.go`, and `internal/cli/state_new_test.go` before authoring — most of it unnecessary once `--help` was read. The doc gap plus the absence of a body template in `--help` is what drove the spelunking.

## Approach

Two small, independent fixes:

1. A terse **"Filing an entity"** stanza for the FO — ideally in a LAZILY-LOADED reference (loaded when the FO is shaping/filing, the way `present-gate` loads at the gate), NOT the always-resident shared core, to protect the boot token budget.
2. A **frontmatter/body template** in `spacedock new --help`.

## Acceptance criteria

- **AC-1 (contract stanza)** — The FO can author a correct entity from a "Filing an entity" stanza: `spacedock new [--folder] SLUG < body` mints the id (never hand-write `id`); the body carries `title` / `status` / `source` / `priority` frontmatter; follow the value-AC rule. The stanza is terse and lives in a reference loaded at shaping/filing time, not the always-on boot core (net boot-token delta ~0).
- **AC-2 (help template)** — `spacedock new --help` prints a minimal frontmatter + body skeleton and states that the id is minted, so a single `--help` read suffices to author a correct entity without reading source.
- **AC-3 (value)** — From the stanza alone, the authoring path requires ZERO source reads (no `new.go` / `state_new_test.go`). Baseline: filing the predecessor entity needed `--help` + 3 source files; the measure must flip that to one-shot. Verified by a from-cold authoring walkthrough.
- **AC-4 (validation)** — If `new`'s help is asserted in `internal/cli/new_help_test.go`, update it for the template; `go test ./internal/cli` green.

## Notes

- FO-efficiency friction program: `0203-fo-efficiency`.
- The stanza must respect the contract token budget that 0221 / 0223 have been trimming — favor a lazily-loaded reference over always-resident core text.
