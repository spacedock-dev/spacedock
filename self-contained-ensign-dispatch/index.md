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
