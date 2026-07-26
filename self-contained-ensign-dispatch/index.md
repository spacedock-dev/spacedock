---
title: Make ensign dispatch self-contained across launcher drift
status: ideation
sprint: durable-decisions
source: "Real source-build sprint dogfood, 2026-07-26"
score: "1.0"
id: kd7877nnbd19d528xnpwwaj4
gates:
    version: 1
    current:
        gate: gate:kd7877nnbd19d528xnpwwaj4:backlog
    records:
        - id: gate:kd7877nnbd19d528xnpwwaj4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:kd7877nnbd19d528xnpwwaj4-backlog-1
              briefing:
                id: briefing:kd7877nnbd19d528xnpwwaj4:backlog:attempt-1:revision-1
                digest: sha256:c22348f0cd78c4310e640953f1574d3c66a63fd1c5cad6669c5960234e767c6e
                digest-domain: canonical-bytes
                request-digest: sha256:cd3dab4ce2732d142d02b952b8f087a60c62ea13935e3ff24a000fc659b288eb
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:kd7877nnbd19d528xnpwwaj4:backlog:1
                briefing: briefing:kd7877nnbd19d528xnpwwaj4:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T15:01:27.186672Z"
                decision: approve
                reason: Real sprint dogfood proves fresh ensigns can resolve 0.26.0+dev while the FO uses 0.27.0-pre1; dispatch build already owns the exact stage subsection, so ideation should remove the redundant worker fetch and bound any remaining helper use.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
started: 2026-07-26T15:08:14Z
---

## Outcome

A successful `dispatch build` writes a self-contained assignment containing the exact
stage/context snapshot and standing-teammate context selected by that invocation. The
same file pins the builder executable's resolved absolute path for legitimate
ensign-owned workflow helpers. The spawn/reuse message remains only the existing tiny
dispatch-file pointer.

## Problem statement

The First Officer invoked `/tmp/spacedock-s4-73810d1e` (`0.27.0-pre1`), while the
Codex worker environment exposed
`SPACEDOCK_BIN=/Users/clkao/git/spacedock-research/spacedock-v1/spacedock`
(`0.26.0+dev`) and PATH exposed Homebrew `0.27.0-pre0`. `spawn_agent` has no
environment argument, and a shell export in one tool call does not persist into later
tool calls. Nevertheless, `dispatch build` discarded the stage/context bytes it had
already resolved and emitted a worker command that re-selected them through
`${SPACEDOCK_BIN:-spacedock}`. The worker could therefore obey its assignment and
still use a different launcher contract.

Removing all worker-side Spacedock use is not the goal. `status --read` is a legitimate
fence-safe entity read and supplies the exact append offset for a stage report. The dev
workflow also assigns `gate record --round` to the worker completing an advisory
Roborev/in-stage feedback round; that command publishes a round but neither applies a
gate nor advances status. Those calls need the builder's exact launcher path.

## Historical separation

Fetch-on-demand entity `claude-team-build-fetch-on-demand-dispatch-spec`
(`0x93enxe1hpmk95a25476zyn`, archived state commit `a46706e3`) removed stage and
standing prose from the full `Agent()` prompt because the same prose was then paid in
helper stdout, full tool arguments, and ensign context. Entity
`dispatch-prompt-as-file-pointer` (`rdt`, archived state commit `a17449cd`) later moved
the full PR #231-shaped body to `/tmp` and measured a 93.2% prompt reduction plus
roughly 80% lower First Officer per-dispatch context. It intentionally retained the
older fetch commands without revisiting them.

This task does not return prose to full spawn arguments. It preserves `rdt`'s
approximately 175-character file pointer and moves already-resolved bytes only into
the dispatch file that the pointer names.

## Spike: exact two-binary failure

A throwaway fixture built binary A from `50f8d1fb` and binary B from
`dd6bd114` (before declared `context-sections`). Its workflow selected an `ideation`
subsection plus `context-sections: [Authority]`.

| Observation | Fresh | Reuse advance |
|---|---:|---:|
| A's selected contract included `A-only-declared-context` | yes | yes |
| Generated file embedded A's selected bytes | no | no |
| Generated `fetch_commands` count | 1 | 1 |
| Worker fetch with both `SPACEDOCK_BIN` and PATH resolving B included the declared context | no | no |

B returned the stage subsection successfully, so the failure was silent rather than a
missing-command error. Current baseline is therefore **2 of 2** generated assignments
re-resolving through B and **0 of 2** carrying A's complete selection. This exercise is
the seed for the implementation fixture.

## Ownership inventory

| Operation | Final owner and launcher rule |
|---|---|
| Select stage subsection plus declared `context-sections` | `dispatch build`/FO side; inline the single `resolveStageContext` result into the file |
| Render optional standing-teammate routing context | `dispatch build`/FO side; inline the existing renderer's result under its current legacy-team/non-empty condition |
| `dispatch show-stage-def` / `dispatch show-standing` bootstrap fetches | Delete from generated assignments; retain public inspection subcommands |
| Entity section reads and stage-report `total_lines` via `status --read` | Ensign-owned read; invoke the dispatch-pinned absolute launcher |
| Stage Report body append | Ensign-owned direct file edit at the dispatched entity path; no Spacedock helper |
| Split-root entity/report commit and push | Ensign-owned path-scoped Git sequence already rendered with resolved paths; no Spacedock helper |
| Advisory `gate record --round` for dev Roborev/in-stage feedback | Ensign-owned; invoke the dispatch-pinned launcher; it records no gate application/status transition |
| Gate decision/application/consume, status transitions, state ready/sweep, merge guard | First Officer-owned and absent from the ensign assignment |
| Worktree-built Spacedock binary as product under test | Ensign-owned explicit test target; use its worktree path, never as the workflow-control launcher |

Stage text may name other project-specific Spacedock product tests. The assignment-level
rule distinguishes those explicit test targets from workflow helpers; the builder does
not rewrite arbitrary stage prose.

## Proposed approach

1. At CLI entry, reuse the launcher's existing `os.Executable` → absolute path →
   `EvalSymlinks` resolution and executable-file check. Pass that internally to
   `dispatch build`; do not accept a caller launcher flag. A successful artifact build
   fails closed if the running executable cannot be resolved.
2. Add a short `### Workflow launcher` block near the top of every fresh and
   `--advance` dispatch file. It contains the shell-quoted absolute A path and directs
   each ensign workflow-helper call to begin with that literal path. It explicitly
   forbids substituting inherited `SPACEDOCK_BIN` or PATH and separates a worktree
   product-test binary.
3. Append the already-computed `resolveStageContext(readmeData, stage)` result at the
   former fetch position. This snapshots the stage subsection and ordered inherited,
   replaced, or cleared declared context sections from one validated README buffer.
4. Render `EnumerateDeclaredStandingTeammates` through the existing standing renderer
   during build and append non-empty output directly. Preserve today's condition:
   legacy team dispatch with declared standing teammates gets the section; bare,
   merged/no-team, and empty-standing cases do not.
5. Emit no `### Fetch commands` block. Keep the existing schema-v2
   `fetch_commands` member as an empty array: it truthfully reports that the assignment
   has no bootstrap fetch and avoids inventing a replacement protocol or output field.
   Remove the ensign fetch-bootstrap and environment-fallback launcher prose.
6. Update ensign instructions so `status --read`, including the report-append offset
   read, uses the pinned path. Update dev feedback-round wording so
   `gate record --round` does the same. Fresh and reused workers obtain the binding
   from the current dispatch file, so an advance built by a newer A supersedes the
   previous assignment's launcher path.

The exact command identity is a resolved path, not a content digest, copied binary, or
lease. If that file disappears, a worker helper fails loudly and reports the path;
there is no fallback to B. Replacing bytes in place at the same path belongs to the
separate v21 source-build identity task. Break-glass manual dispatch has no successful
builder identity to pin and remains outside this change.

## Acceptance criteria

**AC-1 (VALUE)** In the same two-binary fresh-plus-advance fixture, the number of
worker-time stage/standing resolutions through B is **0 of 2**, down from the observed
baseline **2 of 2**, while both files contain A's exact resolved stage subsection and
ordered declared sections.

Verified by a command-level fixture that builds through A, makes both
`SPACEDOCK_BIN` and PATH resolve an instrumented B, reads each generated assignment,
and asserts B's invocation log stays empty. Removing the inline package or restoring a
fetch command makes it fail.

**AC-2** Fresh and reuse-advance dispatch files are self-contained for all context the
builder already owns: stage subsection, inherited/replaced/cleared declared sections,
and conditional standing-teammate routing. Their outer spawn/reuse prompts remain
file pointers no longer than the existing 300-byte ceiling.

Verified by fresh/advance golden and behavior fixtures with independent sentinels in
each selected section and standing mod. Reordering/dropping a section, broadening the
standing condition, or inlining the body into `prompt` makes them fail.

**AC-3** Every successful generated assignment carries A's shell-safe resolved
absolute executable path, and a retained ensign helper (`status --read` and the dev
round-record probe) reaches A even when inherited `SPACEDOCK_BIN` and PATH both select
B. An unresolvable builder executable refuses artifact creation; no fallback occurs.

Verified by extracting the launcher line from files built from a path containing a
space, executing the two representative helper probes under the B environment, and
asserting A-only output/on-disk state plus a zero-entry B log. Replacing the literal
path with `${SPACEDOCK_BIN:-spacedock}` or bare `spacedock` makes it fail.

**AC-4** Workflow-control authority remains unchanged: workers can read entities,
append reports, commit their assigned paths, and publish advisory rounds; only the
First Officer applies gates, transitions status, runs state scheduling, or merges.
Product-under-test runs use an explicit worktree binary and cannot silently become the
workflow launcher.

Verified by the existing gate-round, dispatch, and host live scenarios plus a
two-path test that gives the worktree product binary a third identity C and observes A
for workflow helpers and C only for the explicit product test. Any inferred
ensign-side gate application or use of C for `status --read` fails the scenario.

## Test plan

Add the exact two-binary fixture before production edits. Use real `dispatch build`
through A and an instrumented executable B, with fresh and `--advance` subtests,
declared-context sentinels, a standing-mod case, and launcher paths containing spaces.
Reuse the existing dispatch golden harness rather than adding a launcher protocol or
runtime wrapper. Estimated focused cost: under 10 seconds and moderate fixture work.

Update stage-discipline, advance-content, standing, hazard/quoting, host, and output
shape tests. Preserve direct `dispatch show-stage-def` and `show-standing` tests because
the commands remain useful inspection surfaces; only build-time bootstrap use changes.
Regenerate the affected build goldens and assert `fetch_commands: []` with no fetch
heading in the file.

Run `go test ./internal/dispatch ./internal/cli`, `go test ./skills/integration`,
`go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`. Because the
host-neutral dispatch path and ensign contract change, all existing Claude, Codex, and
Pi live lanes are required green. The live checks must observe the pinned A helper,
not merely grep instruction prose.

Before merge, run the dev workflow's detached adversarial audit in a throwaway
checkout. Perturbations reintroduce a worker fetch, swap the pinned literal for the
environment fallback, drop the standing render, and route a workflow helper through
the explicit product-test C; each corresponding test must turn red.

## User-visible documentation diff

In `docs/site/concepts/workflows-and-entities.md`, replace:

> `dispatch build` validates every selection before spawn; the worker's existing
> `show-stage-def` read returns the stage definition followed by those sections in
> declaration order. ... The read uses the then-current valid README.

with:

> `dispatch build` resolves and validates every selection from one README snapshot,
> then writes the stage definition followed by those sections into the dispatch file
> in declaration order. Later README edits affect later dispatches, not an assignment
> already built. The file also pins the builder's resolved absolute launcher for
> worker-owned workflow helpers.

In `docs/site/reference/command-reference.md`, replace the dispatch row's description:

> Build the worker dispatch artifacts (`dispatch build`, `dispatch show-stage-def`)

with:

> Build self-contained worker dispatch artifacts (`dispatch build`) or inspect a
> resolved stage/context package directly (`dispatch show-stage-def`).

In `skills/ensign/references/ensign-shared-core.md`, delete the launcher fallback
invariant and `## Fetch-on-Demand Bootstrap`, and add:

> The dispatch file pins one absolute workflow launcher. Begin every ensign-owned
> Spacedock workflow-helper command, including `status --read`, with that literal
> path. Do not re-resolve `SPACEDOCK_BIN` or PATH. An explicitly named worktree binary
> is the product under test, not the workflow launcher.

Change the report append rule to obtain `total_lines` with the pinned launcher's
`status --read`. Change the Codex runtime summary from “file pointers carrying fetch
commands” to “file pointers carrying resolved stage/context and the pinned workflow
launcher.” In `docs/dev/README.md`, replace the feedback-round
`${SPACEDOCK_BIN:-spacedock} gate record` prefix with “the dispatch-pinned absolute
workflow launcher” and retain the existing round arguments and no-application rule.

## Expected surface

Expected hand-authored surface is about **18 files / 450-650 changed LOC**:
`internal/cli/cli.go`; `internal/dispatch/{dispatch.go,build.go}`; deletion of
`launcher_command.go`; 8-10 focused dispatch/CLI test files (including rewriting the
old launcher-drift tests); and the five instruction/doc files named above. Expect
**17 build golden files / 250-450 generated changed lines** to replace fetch blocks
with resolved context and the normalized launcher.

Tolerance is ±5 non-golden files, ±200 hand-authored changed LOC, and ±5 golden
fixtures. Any new command, package, JSON field/schema version, caller flag,
environment propagation, binary copy/digest/lease, plugin-private wrapper, or outer
prompt above 300 bytes requires a design reset rather than consuming tolerance.

## Boundary

Do not fold this into s4 gate-room behavior, v21 source-build content identity, a new
launcher lifecycle, or an ensign-without-Spacedock design. The durable-decisions
walking skeleton and pre-release wait for this launcher-consistent artifact.

## Stage Report: ideation

- DONE: Spike the exact two-binary failure: builder binary A must select the stage contract while the worker environment and PATH resolve binary B, and record the observable current behavior.
  A=`50f8d1fb`, B=`dd6bd114`: fresh and advance both omitted A's declared context, emitted one fetch, and silently returned B's smaller contract.
- DONE: Design the smallest self-contained fresh and advance dispatch that embeds the already-resolved stage subsection plus declared context sections, and inventory every remaining ensign-side Spacedock call with explicit FO/worker ownership.
  The design preserves the tiny pointer, inlines resolved stage/declared/conditional-standing context, pins A in-file, and assigns retained reads/advisory rounds to the worker while FO keeps applications/transitions/merge.
- DONE: Specify falsifiable fixtures and user-visible documentation changes, with an expected files/LOC surface and tolerance; do not add environment transport, caller flags, wrappers, compatibility, or another launcher protocol.
  ACs name red-making perturbations, docs carry concrete before/after wording, and the estimate is 18 hand-authored files/450-650 LOC plus 17 goldens with explicit tolerance and reset triggers.

### Summary

Ideation reproduced the launcher-drift failure for both dispatch lifecycles and traced why fetch-on-demand outlived its original full-prompt cost model. The approved direction keeps file-pointer economics, snapshots builder-selected context, and binds legitimate worker helpers to the builder's resolved absolute executable without adding transport or launcher machinery.
