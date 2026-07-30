---
title: Make ensign dispatch self-contained across launcher drift
status: validation
sprint: durable-decisions
source: "Real source-build sprint dogfood, 2026-07-26"
score: "1.0"
id: kd7877nnbd19d528xnpwwaj4
gates:
    version: 1
    current:
        gate: gate:kd7877nnbd19d528xnpwwaj4:ideation
    records:
        - id: gate:kd7877nnbd19d528xnpwwaj4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:kd7877nnbd19d528xnpwwaj4-backlog-1
              briefing:
                id: briefing:kd7877nnbd19d528xnpwwaj4:backlog:attempt-1:revision-1
                digest: sha256:c22348f0cd78c4310e640953f1574d3c66a63fd1c5cad6669c5960234e767c6e
                digest-domain: canonical-bytes
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
        - id: gate:kd7877nnbd19d528xnpwwaj4:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:kd7877nnbd19d528xnpwwaj4-ideation-1
              briefing:
                id: briefing:kd7877nnbd19d528xnpwwaj4:ideation:attempt-1:revision-1
                digest: sha256:d372d2b50317aea3e30be46834f4b71d0cc7a8c4198734fbcf4221937b436dbc
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:kd7877nnbd19d528xnpwwaj4:ideation:1
                briefing: briefing:kd7877nnbd19d528xnpwwaj4:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T16:02:05.670336Z"
                decision: approve
                reason: 'Delegated sprint conn: exact two-binary spike and final independent staff re-review approve the self-contained artifact, pointer-only transport design; implementation must preserve legitimate ensign round recording through the builder-pinned launcher.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
started: 2026-07-26T15:08:14Z
worktree: .worktrees/spacedock-ensign-self-contained-ensign-dispatch
---

## Outcome

**Self-contained artifact; pointer-only transport.** A successful `dispatch build`
writes a dispatch file containing the exact stage/context snapshot, checklist,
standing-teammate context, scope/feedback instructions, and launcher binding selected
by that invocation. Here and throughout this design, “self-contained” describes that
file, never the outer transport. “Pointer-only transport” means the dispatch-file
pointer plus fixed transport/routing metadata. Fresh `output.prompt` retains its fixed
read instruction, and reuse-advance `prompt` retains the fixed-format
`Advancing to next stage: {stage}.` routing label. Neither may carry stage-definition,
declared-context, checklist, standing, scope, or feedback payload bytes.

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

1. Make **Self-contained artifact; pointer-only transport** the rendering invariant
   shared by fresh and `--advance` across Claude, Codex, and Pi. The dispatch file owns
   all assignment payload bytes; fresh `output.prompt` and reuse-advance `prompt` own
   only the file locator and fixed transport/routing metadata. Preserve the current
   fixed-format `Advancing to next stage: {stage}.` label; the stage identifier routes
   the message, but stage-definition, declared-context, checklist, standing, scope,
   and feedback payload bytes must never be interpolated into either prompt. Add that
   distinction in the source comment beside the two prompt renderers in
   `internal/dispatch/build.go`. This changes neither prompt/bootstrap format.
2. At CLI entry, reuse the launcher's existing `os.Executable` → absolute path →
   `EvalSymlinks` resolution and executable-file check. Pass that internally to
   `dispatch build`; do not accept a caller launcher flag. A successful artifact build
   fails closed if the running executable cannot be resolved.
3. Add a short `### Workflow launcher` block near the top of every fresh and
   `--advance` dispatch file. It contains the shell-quoted absolute A path and directs
   each ensign workflow-helper call to begin with that literal path. It explicitly
   forbids substituting inherited `SPACEDOCK_BIN` or PATH and separates a worktree
   product-test binary.
4. Append the already-computed `resolveStageContext(readmeData, stage)` result at the
   former fetch position. This snapshots the stage subsection and ordered inherited,
   replaced, or cleared declared context sections from one validated README buffer.
5. Render `EnumerateDeclaredStandingTeammates` through the existing standing renderer
   during build and append non-empty output directly. Preserve today's condition:
   legacy team dispatch with declared standing teammates gets the section; bare,
   merged/no-team, and empty-standing cases do not.
6. Emit no `### Fetch commands` block. Keep the existing schema-v2
   `fetch_commands` member as an empty array: it truthfully reports that the assignment
   has no bootstrap fetch and avoids inventing a replacement protocol or output field.
   Remove the ensign fetch-bootstrap and environment-fallback launcher prose.
7. Update ensign instructions so `status --read`, including the report-append offset
   read, uses the pinned path. Update dev feedback-round wording so
   `gate record --round` does the same. Fresh and reused workers obtain the binding
   from the current dispatch file, so an advance built by a newer A supersedes the
   previous assignment's launcher path.
8. Extend the two existing transcript consumers in place. In
   `internal/journeymetrics/readadoption.go`, make the existing status-read command
   recognizer accept an absolute pinned executable token, including shell-quoted paths
   with spaces and names such as `/tmp/spacedock-s4-*`, while preserving its
   non-invocation controls. In
   `internal/ensigncycle/shared_round_recording_test.go`, thread the live fixture's
   already-known A path into the existing Claude/Codex transcript extractors and count
   the round invocation only when the observed executable is exactly A. Keep the
   durable round oracle's existing advisory Resolution, unchanged status, and absent
   gate/application assertions. This extends current observers; it adds no event,
   artifact field, observer, or launcher protocol.

The exact command identity is a resolved path, not a content digest, copied binary, or
lease. If that file disappears, a worker helper fails loudly and reports the path;
there is no fallback to B. Replacing bytes in place at the same path belongs to the
separate v21 source-build identity task. Break-glass manual dispatch has no successful
builder identity to pin and remains outside this change.

**Deferred risk:** a pinned path can later disappear or be replaced in place. Current
supported release and source-build launch paths are durable executable files for the
assignment lifetime, and the observed `/tmp/spacedock-s4-*` builder remains present
for its run. Do not add copying, hashing, shortening, a lease, or recovery here.
Promote this risk only if an ephemeral/mutable launcher path becomes supported or is
observed changing during a dispatched assignment.

## Acceptance criteria

**AC-1 — Self-contained artifact; pointer-only transport.** For fresh and
reuse-advance builds on every supported host, assignment-bearing bytes exist only in
the dispatch file. Fresh `output.prompt` and reuse-advance `prompt` contain only the
dispatch pointer plus fixed transport/routing metadata. The latter preserves
`Advancing to next stage: {stage}.`; neither prompt contains stage-definition,
declared-context, checklist, standing, scope, or feedback payload bytes.

Verified by paired relational fixtures for Claude, Codex, and Pi. Holding the output
path, stage name, and all fixed transport/routing metadata constant, independently
grow the selected stage definition, declared context, and checklist by N bytes with
distinct sentinels. In each pair the dispatch file contains the new sentinel and grows
with it, while the corresponding outer prompt is byte-identical to baseline and
contains none of the sentinels. Separate sentinels prove standing, scope, and feedback
stay file-only. Moving any payload byte into a prompt, or leaving any selected byte
out of the file, makes the fixture fail.

Both prompts retain the existing O(`dispatch_file_path`) shape. Explicit success
fixtures build a fresh max-legal 251-character dispatch stem and an advance
max-legal 243-character base stem plus the literal `-advance` suffix. Each exits 0,
keeps the complete expected filename unchanged, and remains distinct from a peer stem
that differs at its final character; truncation, shortening, or collision makes the
fixture fail. Those cases explicitly accept outer prompts longer than 300 bytes. The
≤300 values remain measurements of the supported dev-workflow fresh and advance
fixtures, never a universal name/path invariant or rejection threshold.

**AC-2 (VALUE)** In the same two-binary fresh-plus-advance fixture, the number of
worker-time stage/standing resolutions through B is **0 of 2**, down from the observed
baseline **2 of 2**, while both files contain A's exact resolved stage subsection and
ordered declared sections.

Verified by a command-level fixture that builds through A, makes both
`SPACEDOCK_BIN` and PATH resolve an instrumented B, reads each generated assignment,
and asserts B's invocation log stays empty. Removing the inline package or restoring a
fetch command makes it fail.

**AC-3** Every successful generated assignment carries A's shell-safe resolved
absolute executable path, and a retained ensign helper (`status --read` and the dev
round-record probe) reaches A even when inherited `SPACEDOCK_BIN` and PATH both select
B. An unresolvable builder executable refuses artifact creation; no fallback occurs.

Verified by extracting the launcher line from files built from a path containing a
space, executing the two representative helper probes under the B environment, and
asserting A-only output/on-disk state plus a zero-entry B log. Replacing the literal
path with `${SPACEDOCK_BIN:-spacedock}` or bare `spacedock` makes it fail.

**AC-4** Existing journey-metrics and correction-round observers recognize arbitrary
artifact-pinned launcher paths without relaxing identity or invocation evidence:
`status --read` counts for unquoted `/tmp/spacedock-s4-*` and shell-quoted
space-bearing A paths, while a recorded round counts only when the transcript's
launcher equals the known fixture A.

Verified by positive/negative table cases in
`internal/journeymetrics/readadoption_test.go` and
`internal/ensigncycle/shared_round_recording_test.go`, plus the current durable round
oracle. A quoted A mis-tokenization, a B invocation, a no-invocation transcript, a
non-advisory Resolution, or any status/gate/application mutation makes the evidence
fail.

**AC-5** Workflow-control authority remains unchanged: workers can read entities,
append reports, commit their assigned paths, and publish advisory rounds; only the
First Officer applies gates, transitions status, runs state scheduling, or merges.
Product-under-test runs use an explicit worktree binary and cannot silently become the
workflow launcher.

Verified by the existing gate-round, dispatch, and host live scenarios plus a
two-path test that gives the worktree product binary a third identity C and observes A
for workflow helpers and C only for the explicit product test. The round oracle
continues to prove advisory-only output with no status/application side effects. Any
inferred ensign-side gate application or use of C for `status --read` fails the
scenario.

## Test plan

Add the exact two-binary fixture before production edits. Use real `dispatch build`
through A and an instrumented executable B, with fresh and `--advance` subtests,
declared-context sentinels, a standing-mod case, and launcher paths containing spaces.
Reuse the existing dispatch golden harness rather than adding a launcher protocol or
runtime wrapper. Estimated focused cost: under 10 seconds and moderate fixture work.

Add the AC-1 paired relational behavior test for fresh and advance on Claude, Codex,
and Pi. For each host, hold the dispatch locator, stage name, and fixed routing metadata
constant while independently adding N bytes and a unique sentinel to stage-definition,
declared-context, and checklist inputs. Assert the dispatch file contains and grows
with each addition, while `output.prompt` or the reuse-advance `prompt` stays
byte-identical and excludes every sentinel. Give standing, scope, and feedback
independent sentinels and assert they are file-only. Preserve the advance prompt's
fixed `Advancing to next stage: {stage}.` line in the baseline and grown cases. This
test, rather than the ≤300-byte dev-fixture measurement, protects the transport
boundary.

Add direct success cases for a fresh 251-character dispatch stem and an advance
243-character base stem that becomes 251 characters only after appending `-advance`.
For each, assert exit 0 and the exact unsimplified filename. Pair each with a second
legal stem differing at the final character and assert both full filenames exist and
do not collide. Place the dispatch directory deeply enough that the resulting outer
prompt exceeds 300 bytes and assert the build still succeeds. A universal prompt
length check, truncation, hash replacement, or collision must make these cases fail.

Update stage-discipline, advance-content, standing, hazard/quoting, host, and output
shape tests. Preserve direct `dispatch show-stage-def` and `show-standing` tests because
the commands remain useful inspection surfaces; only build-time bootstrap use changes.
Regenerate **exactly the 26 affected `build-*` goldens** and assert
`fetch_commands: []` with no fetch heading in the file. Preserve every meaningful
single/split-root, flat/folder, worktree/non-worktree, team/bare, host, model,
scope/feedback, quoting, and advance cross-product row; do not reduce the golden count
to avoid fixture churn.

Add positive and negative cases to the existing journey-metrics matcher for an
unquoted `/tmp/spacedock-s4-*`, a shell-quoted A path containing spaces, an absolute
non-invocation/other subcommand, and quoted text passed to another command. Exercise
both Claude and Codex transcript shapes so a correct direct command counts once and a
lookalike does not.

Change the existing shared round invocation extractors to take known fixture A and
compare the transcript's parsed executable to it. Positive cases cover direct,
multiline, and shell-quoted space-bearing A; negative cases run the otherwise-identical
round command through B and omit invocation entirely. Reuse
`assertRejectionRecordedRound` to prove the retained room is canonical, Resolutions
remain advisory, lifecycle sentinels survive, status stays unchanged, and no
gate/application state appears. Do not add a second transcript observer or event
format.

Run `go test ./internal/dispatch ./internal/cli`, `go test ./skills/integration`,
`go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`. Because the
host-neutral dispatch path and ensign contract change, all existing Claude, Codex, and
Pi live lanes are required green. The live checks must observe the pinned A helper,
not merely grep instruction prose.

Before merge, run the dev workflow's detached adversarial audit in a throwaway
checkout. Perturbations reintroduce a worker fetch, swap the pinned literal for the
environment fallback, drop the standing render, and route a workflow helper through
the explicit product-test C. Further perturbations restore whitespace-only tokenization
or let the round observer accept B in place of known A. Apply a universal 300-byte
prompt rejection and a filename-shortening mutation; the max-legal fresh/advance
fixtures must turn red. Each corresponding test must turn red.

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
> worker-owned workflow helpers. “Self-contained” refers to this artifact; the outer
> fresh or reuse-advance prompt carries only its locator plus fixed transport/routing
> metadata and never transports assignment payload. Reuse advance retains its
> fixed-format stage-routing label.

In `docs/site/reference/command-reference.md`, replace the dispatch row's description:

> Build the worker dispatch artifacts (`dispatch build`, `dispatch show-stage-def`)

with:

> Build self-contained worker dispatch files with pointer-only outer transport
> (`dispatch build`) or inspect a resolved stage/context package directly
> (`dispatch show-stage-def`).

In `skills/ensign/references/ensign-shared-core.md`, delete the launcher fallback
invariant and `## Fetch-on-Demand Bootstrap`, and add:

> The dispatch file pins one absolute workflow launcher. Begin every ensign-owned
> Spacedock workflow-helper command, including `status --read`, with that literal
> path. Do not re-resolve `SPACEDOCK_BIN` or PATH. An explicitly named worktree binary
> is the product under test, not the workflow launcher.

Change the report append rule to obtain `total_lines` with the pinned launcher's
`status --read`. In `skills/ensign/references/codex-ensign-runtime.md`, change the
Codex runtime summary from “file pointers carrying fetch commands” to “dispatch files
carrying resolved stage/context and the pinned workflow launcher, reached through a
pointer plus fixed transport/routing metadata.” In `docs/dev/README.md`, replace the
feedback-round
`${SPACEDOCK_BIN:-spacedock} gate record` prefix with “the dispatch-pinned absolute
workflow launcher” and retain the existing round arguments and no-application rule.

In `skills/first-officer/references/fo-dispatch-core.md`, replace the loose
“~175-char file-pointer” description with:

> **Self-contained artifact; pointer-only transport.** Forward fresh `output.prompt`
> and reuse-advance `prompt` unchanged as a pointer plus fixed transport/routing
> metadata. Preserve `Advancing to next stage: {stage}.`; stage-definition,
> declared-context, checklist, standing, scope, and feedback payload bytes MUST NEVER
> enter either outer prompt and belong in the dispatch file the ensign reads first.

Mirror this invariant in a short source comment beside the fresh and advance prompt
renderers in `internal/dispatch/build.go`; do not create an ADR, specification, output
field, or command.

In both `docs/runtime-support.md` and
`docs/site/contributing/adding-a-runtime.md`, replace the current propagation
paragraph:

> `spacedock claude` and `spacedock codex` attach `SPACEDOCK_BIN` to the process they
> exec. ... If a wrapper or runtime strips `SPACEDOCK_BIN`, the skill contract's
> `${SPACEDOCK_BIN:-spacedock}` convention degrades to PATH.

with:

> The First Officer continues to invoke Spacedock workflow control through
> `${SPACEDOCK_BIN:-spacedock}`; the front door propagates that value through supported
> wrappers. A successful `dispatch build` separately writes its own resolved absolute
> launcher into the dispatch artifact, and a dispatched ensign uses that literal path
> for ensign-owned helpers. Worker environments may still inherit or expose B through
> `SPACEDOCK_BIN` or PATH; dispatch does not rely on that environment to bind builder
> identity. An explicit worktree binary remains only the product under test.

## Expected surface

Expected hand-authored surface is about **27 files / 650-900 changed LOC**:
`internal/cli/cli.go`; `internal/dispatch/{dispatch.go,build.go}`; deletion of
`launcher_command.go`; 10-12 focused dispatch/CLI test and harness files; the two
journey-metrics files; the shared round oracle plus 2-3 existing Claude/Codex live
runner plumbing files that already know fixture A; and the eight instruction/doc files
named above. Expect **exactly 26 affected build golden files / 450-750 generated
changed lines** to replace fetch blocks with resolved context and the normalized
launcher.

Tolerance is ±5 non-golden files and +250/−150 hand-authored changed LOC. The affected
golden count is fixed at 26, not a tolerance band. Any new command, package, observer,
event/output field or schema version, caller flag, environment propagation,
binary copy/digest/lease, plugin-private wrapper, universal pointer-length guard,
shortening/rejection mechanism, or loss of a meaningful cross-product fixture requires
a design reset rather than consuming tolerance.

## Boundary

**Self-contained artifact; pointer-only transport** is the boundary: assignment bytes
live in the dispatch file, and fresh/reuse outer prompts remain the pointer plus fixed
transport/routing metadata regardless of assignment size. The fixed-format
`Advancing to next stage: {stage}.` label is routing metadata, not stage-definition
payload. The ≤300-byte check is only a supported dev-fixture measurement; legal
prompts may exceed 300 bytes and must not be shortened or rejected. Do not fold this
into s4 gate-room behavior, v21 source-build content identity, a new launcher
lifecycle, an ensign-without-Spacedock design, a new ADR/specification/command, or
prompt/bootstrap changes. The durable-decisions walking skeleton and pre-release wait
for this launcher-consistent artifact.

## Stage Report: ideation

- DONE: Spike the exact two-binary failure: builder binary A must select the stage contract while the worker environment and PATH resolve binary B, and record the observable current behavior.
  A=`50f8d1fb`, B=`dd6bd114`: fresh and advance both omitted A's declared context, emitted one fetch, and silently returned B's smaller contract.
- DONE: Design the smallest self-contained fresh and advance dispatch that embeds the already-resolved stage subsection plus declared context sections, and inventory every remaining ensign-side Spacedock call with explicit FO/worker ownership.
  The design preserves the tiny pointer, inlines resolved stage/declared/conditional-standing context, pins A in-file, and assigns retained reads/advisory rounds to the worker while FO keeps applications/transitions/merge.
- DONE: Specify falsifiable fixtures and user-visible documentation changes, with an expected files/LOC surface and tolerance; do not add environment transport, caller flags, wrappers, compatibility, or another launcher protocol.
  ACs name red-making perturbations, docs carry concrete before/after wording, and the estimate is 18 hand-authored files/450-650 LOC plus 17 goldens with explicit tolerance and reset triggers.

### Summary

Ideation reproduced the launcher-drift failure for both dispatch lifecycles and traced why fetch-on-demand outlived its original full-prompt cost model. The approved direction keeps file-pointer economics, snapshots builder-selected context, and binds legitimate worker helpers to the builder's resolved absolute executable without adding transport or launcher machinery.

## Stage Report: ideation (cycle 2)

- DONE: MATERIAL — add the existing read-adoption and shared round-recording consumers, with positive and negative tests, to the launcher-pinning design.
  AC-4 and the test plan now cover `/tmp/spacedock-s4-*`, quoted space-bearing A, lookalikes, B, and no-invocation controls without a new observer.
- DONE: NEEDS DECISION — keep pointer size O(path) and limit the ≤300-byte claim to supported dev-workflow fresh and advance fixtures.
  AC-1 now makes **Self-contained artifact; pointer-only transport** the governing invariant, proves it relationally across hosts, preserves legal long names, and explicitly forbids a universal ceiling, shortening, or rejection mechanism.
- DONE: Rebaseline all affected build goldens to exactly 26 without cutting meaningful cross-product coverage.
  The test plan names the retained dimensions and the expected surface makes 26 a fixed count rather than a tolerance band.
- DONE: Add runtime-support documentation separating First Officer and dispatched-ensign launcher rules.
  Concrete diffs for both runtime-support guides keep FO `${SPACEDOCK_BIN:-spacedock}` while ensigns use the artifact-pinned literal.
- DONE: Keep path-not-content launcher identity as a deferred risk with an explicit promotion trigger.
  The design defers copying/hashing/leases unless ephemeral or mutable paths become supported or are observed changing in-flight.

### Summary

Cycle 2 closes the two observer blind spots and binds the round proof to known fixture A while preserving advisory-only durable-state checks. It makes self-contained-file versus pointer-only-transport the main anti-regression invariant, narrows the pointer measurement, restores the full 26-golden baseline, adds both runtime docs, and recalculates the surface without broadening the launcher mechanism.

## Stage Report: ideation (cycle 3)

- DONE: MATERIAL — add explicit planned success fixtures for fresh max-legal dispatch stem 251 and advance base stem 243 + `-advance`.
  AC-1 and the test plan require exit 0, exact unchanged filenames, distinct peer names without collision, and accepted prompts over 300 bytes; shortening or a universal ceiling turns them red.
- DONE: NEEDS DECISION — define pointer-only as the pointer plus fixed transport/routing metadata under the captain's ruling.
  The invariant preserves `Advancing to next stage: {stage}.`, holds routing metadata fixed in relational tests, and excludes only stage-definition/context/checklist/standing/scope/feedback payload bytes.
- DONE: Fix the named instruction/document count to eight.
  Expected Surface now counts the explicitly named concepts, command reference, ensign core, Codex adapter, dev workflow, FO core, and two runtime-support documents.
- DONE: Replace the misleading worker-environment sentence.
  The runtime-doc diff now allows workers to inherit or expose B through `SPACEDOCK_BIN`/PATH while stating that dispatch never relies on that environment to bind A.

### Summary

This third review correction closes proof and wording gaps under the captain's pointer-only ruling. It preserves the semantic direction and mechanism, changes no prompt/bootstrap contract, and adds no new mechanism, so no design reset is required.

## Stage Report: implementation

- DONE: Deliver self-contained fresh and advance dispatch files while keeping every outer prompt pointer-only across Claude, Codex, and Pi, including the max-legal filename fixtures.
  Commit `706b5cb8` snapshots the resolved stage/context and conditional standing
  package into successful build artifacts, removes fetch commands, preserves fixed
  routing metadata, and keeps exact 251-character fresh/advance stems distinct.
- FAILED: Pin the builder's resolved executable as workflow launcher A so worker-time stage/context resolution through ambient B is 0 of 2 and retained helper probes reach A.
  The command-level fresh/advance fixture proves B receives 0 of 2 stage/context
  reads and synthetic status/round probes reach A, but the real correction-round
  publisher is the First Officer and invokes a variable launcher rather than the
  dispatch-pinned literal. The required live identity proof therefore fails.
- FAILED: Preserve worker/FO authority and existing observer behavior; run the focused two-binary, dispatch, round-recording, Go, race, formatting, and diff checks required by the entity.
  Offline focused, integration, full, race, formatting, and diff checks pass, but
  the strict known-A round observer false-reds the legitimate FO-owned publisher.
  The all-live-lanes criterion is also externally blocked by `se0`.

### Evidence and review

- Offline green: focused dispatch/CLI, skills integration, full, race, formatting, and diff checks.
- Final surface: 32 hand-authored files / 967 changed LOC against the reset cap 32 / 991; exactly 26 goldens / 401 generated changed LOC.
- Roborev job 3153 returned CHANGES_REQUESTED. Its High is reclassified MATERIAL:
  the allowed FO launcher variable cannot prove transcript equality to literal A.
- The three Medium findings were fixed before the stop: same-segment/successful
  round correlation, canonicalized live A, and fail-closed no-artifact evidence.
- Reproduction command: `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/kd-codex-rejection.vI3Y1j go test -tags live -count=1 -timeout 40m -run 'TestLiveCodexSharedScenarios/rejection-flow' ./internal/ensigncycle`.
- Preserved trace: `/tmp/kd-codex-rejection.vI3Y1j/codex-shared-scenarios/rejection-flow/codex-exec.jsonl`;
  item 41 records successful `gate record --round validation/1` through
  `launcher="${SPACEDOCK_BIN:-spacedock}"; "$launcher" ...`.
- The rerun also reproduced `se0`'s Codex rejection-flow failure before the strict
  oracle: one implementation report remained where the trajectory requires two.
- Current `se0` blockers are Claude Sonnet default-gate-stop/gate-guardrail; Claude Opus default-gate-stop/gate-guardrail/recorded-lifecycle-review;
  Codex gate-guardrail/rejection-flow/keep-moving; and Pi recorded-lifecycle-review.
- Current-main dogfood artifact `/tmp/spacedock-dispatch/spacedock-ensign-live-lanes-red-on-every-branch-ideation.md` retains the pre-change fetch fallback.
  This branch changes successful builds; break-glass/manual dispatch without builder identity remains outside scope.

### Required design reset

- Alternative A: move advisory `gate record --round` publication to the in-stage
  ensign so it can execute the artifact-pinned literal A.
- Alternative B: keep the First Officer as publisher and revise AC-4 plus the
  observer to prove the FO variable's controlled resolution instead of demanding
  literal A in transcript text.
- No alternative was selected or implemented. Validation, merge, and the
  all-live-lanes criterion remain blocked pending a captain-approved design and
  `se0` restoration.

### Summary

The core self-contained dispatch implementation and offline proof landed in `706b5cb8`, but live dogfood falsified the approved round-identity assumption. Implementation stops FAILED rather than weakening the oracle or expanding authority without a new design ruling.

## Stage Report: implementation (cycle 2)

- DONE: Implement the captain-selected ensign-owned advisory-round publication so the pinned launcher identity remains end to end.
  Commits `a1ba68aad` and `1d7dfc0b7` assign publication to the in-stage ensign and prove its retained Codex worker transcript invokes exact A; a FO/root or B invocation cannot satisfy the extractor.
- DONE: Preserve the self-contained pointer-only dispatch artifact and fresh/advance max-filename guarantees without adding compatibility behavior.
  Rebased commit `19d422af7` retains zero-fetch self-contained artifacts, exact launcher literals, pointer-only prompts, and the max-legal fresh/advance fixtures; restoring fetches, payload-bearing prompts, shortening, or collisions makes the focused dispatch tests fail.
- DONE: Reproduce the relevant focused live identity journey, classify findings before mutation, and keep the correction within the approved kd semantics and expected surface.
  Candidate-stamped `TestLiveCodexSharedScenarios/rejection-flow` passed in 486.765s with exact-A worker publication and durable advisory-only state; artifacts are under `/tmp/spacedock-kd-cycle2-live.poBn8k/artifacts-harness/codex-shared-scenarios/rejection-flow`.

### Summary

Captain-selected Alternative A is implemented: the ensign publishes completed in-stage advisory rounds through the launcher literal embedded by the builder, while the First Officer retains routing and authorization authority. Focused, full, and race suites pass; the final surface is 35 hand-authored files / 1074 changed lines plus exactly 26 build goldens / 401 generated changed lines, within the explicitly authorized correction expansion.

## Stage Report: implementation (cycle 3)

- DONE: Close the post-rebase ownership-proof gaps: emitted-launcher AC-3 execution, exact ensign-A publication cardinality, any-launcher FO duplicate rejection, current Codex retained-event correlation, and the real rejected-attempt/superseded-feedback/open-regate durable shape.
- DONE: Preserve the product boundary: no launcher, transport, contract, schema, or compatibility change followed the captain-selected Alternative A; all cycle-3 commits are test-only.
- DONE: Validate and review the exact final code tip `9148e557255529ecdca74cb7a988b0c7b22989aa`; focused rejection replay, `go test ./...`, `go test ./... -race`, full formatting, and diff checks pass.
- Exact c4 evidence: `/tmp/spacedock-kd-round-ownership-c4.XXNJXy/artifacts-exact/codex-shared-scenarios/rejection-flow`; the pinned ensign published, root did not, the round was durable, and the retained two-attempt gate shape superseded feedback then left attempt 2 open.
- Roborev 392 and 408 findings were MATERIAL/task-owned and fixed; job 411 findings 1/3 were MATERIAL and fixed while inert-string finding 2 was deferred-risk/task-adjacent DECLINED; job 414 nested-terminal-status was Polish/task-adjacent DECLINED because composed durable validation protects AC-4.
- Bootstrap exception: the pre-feature kd dispatch lacked `### Workflow launcher` and routed round arguments, so this one pre-merge cycle omitted canonical round publication rather than reconstructing authority.
- Final surface is 35 hand-authored files / 1282 changed lines versus the reset estimate 32 / 991, expanded only by authorized proof corrections; generated surface remains exactly 26 goldens / 401 changed lines.

### Summary

Alternative A is merge-ready at the exact code tip: self-contained pointer transport pins the worker launcher, advisory publication is exclusively ensign-owned under the composed transcript/durable oracle, and every material review finding is fixed; the two remaining hypothetical proof refinements are explicitly declined with promotion triggers.

## Stage Report: validation

- FAILED: Verify AC-1 through AC-5 against exact tip 9148e557255529ecdca74cb7a988b0c7b22989aa, reproducing each cited mechanism/value proof and rejecting stale or self-referential evidence.
  AC-1, AC-2, AC-3, and AC-5 reproduce; AC-4 fails the current-Codex, shell-quoted-space-bearing-A observation boundary below.
- DONE: AC-1 — Self-contained artifact; pointer-only transport.
  Fresh/advance relational tests for Claude, Codex, and Pi kept outer prompts byte-identical as stage/context/checklist/standing/scope/feedback grew; exact distinct 251-character stems and >300-byte prompts passed.
- DONE: AC-2 (VALUE) — ambient B performs 0 of 2 worker-time stage/standing resolutions.
  `TestDispatchBuildTwoBinaryCommandFixture` built fresh and advance through A under B-selected PATH/SPACEDOCK_BIN, found A-only stage/context in both files, zero fetches, and an empty B log.
- DONE: AC-3 — every successful assignment pins shell-safe resolved A and fails closed without it.
  `TestDispatchBuildPinsResolvedLauncherAndFailsClosed` executed status/round probes through space-bearing A, kept B unused and explicit product C separate, and proved unresolved A writes no artifact.
- FAILED: AC-4 — observers recognize arbitrary artifact-pinned paths without relaxing identity or invocation evidence.
  Existing tables pass, but adding the supported current-Codex custom-tool case `codexCustomRoundSession(quoted, result, 0)` with space-bearing A fails `TestRejectionFlowRoundInvocationExtractors`.
- DONE: AC-5 — workflow-control authority remains unchanged.
  Deterministic ownership/cardinality and durable-oracle tests pass; retained c4 shows one exact-A ensign publication, no root publication, advisory Resolution, superseded attempt 1, and open attempt 2.
- DONE: Confirm pointer-only fresh/advance dispatch, builder-pinned A under ambient B, ensign-only advisory publication, and unchanged FO application authority across deterministic and retained c4 evidence; do not rerun the expensive live drive unless its report evidence is absent or invalid.
  Retained c4 completed after the final production commit; the two later commits are test-only, its process exited 0, and transcript plus durable bytes remain valid, so no expensive rerun was made.
- DONE: Run the semantic adversarial pass over the exact changed surface, classify every finding under Review-finding disposition without mutating the candidate, and recommend PASSED or REJECTED with deferred risks separate.
  The exact candidate stayed clean at `9148e557`; empty/failed/wrong/B/duplicate/out-of-order/open/closed variants held, while the added current-Codex quoted-A variant exposed the Material defect below.
- DONE: Run required focused, integration, full, race, formatting, and diff checks.
  Focused AC tests, `go test ./...`, `go test ./... -race`, and non-mutating gofmt inspection pass; the code worktree remains clean at the exact tip.

### Review-finding disposition

- MATERIAL / task-owned / evidence defect / FO-authorized disposition REJECTED and route to implementation: current Codex custom-tool round evidence false-negatives shell-quoted pinned A when its path contains spaces.
  Released/normal workflow: successful dispatch explicitly supports shell-safe space-bearing A and current Codex emits custom exec events.
  Observable harm: a valid ensign advisory publication becomes unobservable and false-reds the required live identity proof.
  Authority: `value-ac[AC-3]` requires every successful assignment's shell-safe resolved A to carry representative helpers end to end.
  Trigger evidence: the disposable exact-tip probe fails with `codex custom quoted A extractor case failed`; `customInputRecordsRejectionRound` searches raw input for unquoted `knownLauncher + " gate..."`.

### Deferred risks

- Retain the approved path-lifecycle risk: promote only if ephemeral/mutable launcher paths become supported or are observed disappearing or changing during an assignment.
- Retain the previously declined inert-string false-positive risk: promote if a supported Codex custom event can carry a complete inert round command plus success-shaped output without executing it.

### Polish

- Retain the previously declined nested-terminal-status oracle refinement; composed durable parsing protects the supported AC-4 path, and no current value loss is established.

### Summary

Validation recommends REJECTED at exact tip `9148e557255529ecdca74cb7a988b0c7b22989aa`. Product dispatch, launcher pinning, two-binary isolation, ensign-only publication, and durable FO authority proofs are green, but AC-4 cannot validly observe the explicitly supported combination of current Codex custom-tool events and a quoted launcher path containing spaces; the FO classified and routed this Material evidence defect without validator mutation.

### Feedback Cycles

- Cycle 1: REJECTED — fresh validation / quoted-space Codex observer; surface 35 hand files/1282 LOC + 26 goldens/401 LOC vs estimate 27 hand files/650-900 LOC + 26 goldens/450-750 LOC (142%); AC unchanged
