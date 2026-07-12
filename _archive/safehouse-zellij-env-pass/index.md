---
title: Preserve Zellij targeting metadata through Safehouse
status: done
score: 0.55
source: "Captain request 2026-07-12."
id: 6jk4gverktmthbzkn4vas1kg
worktree: .worktrees/spacedock-ensign-safehouse-zellij-env-pass
started: 2026-07-12T04:17:30Z
mod-block:
pr: pr-merge:499
verdict: passed
completed: 2026-07-12T14:41:08Z
archived: 2026-07-12T14:41:08Z
---

Safehouse sanitizes the environment passed to a wrapped Spacedock host. Tooling
inside a Zellij-hosted sandbox needs the current Zellij identity in order to
address the correct pane and session.

## Chosen direction

The shared Safehouse wrapper owns one built-in terminal/session-targeting
default. When `ZELLIJ` is present, it adds one static
`--env-pass=ZELLIJ,ZELLIJ_PANE_ID,ZELLIJ_SESSION_NAME` argument; the variable
names are fixed, while Safehouse copies the inherited values. Safehouse's own
named-env parser composes that argument with the existing launcher allowance
and optional `SAFEHOUSE_ENV_PASS` operator configuration.

Claude, Codex, and Pi retain their ordinary front-door assembly and existing
`SPACEDOCK_BIN` handling. When the parent is not in Zellij, the wrapper adds
nothing. A future terminal multiplexer can extend this same wrapper helper
without introducing a configuration framework or a Safehouse argv parser.

Keep Safehouse's native `SAFEHOUSE_ENV_PASS` as the operator's global mechanism
for any additional, explicitly chosen variables. Do not add a competing
Spacedock config file, and do not put environment forwarding in the trusted
repository `.safehouse` profile.

## Acceptance criteria

**AC-1 (value): A Safehouse-wrapped host launched from Zellij receives the exact inherited values of `ZELLIJ`, `ZELLIJ_PANE_ID`, and `ZELLIJ_SESSION_NAME`, so in-sandbox tooling can identify and target the active pane/session.** Verified by: a wrapper-owned scrubbed-environment executable smoke whose host prints the three values, plus a control that omits them.

**AC-2 (scope): A non-Zellij parent neither manufactures targeting values nor adds a Safehouse argument, and Claude, Codex, and Pi retain their ordinary `SPACEDOCK_BIN` and launch-environment shapes.** Verified by: a wrapper argv control and an ambient-`ZELLIJ=0` repository test run.

**AC-3 (compatibility): The built-in Zellij argument uses Safehouse's native repeated named-env contract and composes with `SAFEHOUSE_ENV_PASS` for unrelated operator-selected variables, without a Spacedock config surface or broad `--env` inheritance.** Verified by: a focused Safehouse integration fixture and scoped code review.

## Test plan

1. Add one wrapper helper that conditionally constructs the static Zellij
   `--env-pass=` argument, leaving Safehouse to compose named-env inputs.
2. Extend the scrubbed Safehouse smoke to prove exact forwarding of all three
   Zellij values and prove the wrapper-owned stripped control cannot see them.
3. Prove the ordinary front-door shapes remain stable under ambient `ZELLIJ=0`
   and native `SAFEHOUSE_ENV_PASS` composes with the built-in argument.
4. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

## Out of scope

- Hardcoding `ZELLIJ=0`, a pane ID, or a session name.
- Passing the full host environment or secrets through Safehouse.
- Adding a Spacedock-specific persistent configuration format.

### Feedback Cycles

**Cycle 1 — captain ownership and extensibility correction (2026-07-12).**
The pass-through mechanism must be named for generic terminal/session targeting
metadata, not Zellij, so a later tmux set composes without renaming the
abstraction. Keep the current Zellij variables as today's concrete allowlist
data, not the helper/API name. Move environment-to-`--env-pass` composition to
the Safehouse wrapper/its immediate generic helper; do not rethread parent
environment through Claude, Codex, and Pi or add host-specific `WithEnv`
front-door variants merely to implement this shared Safehouse behavior. Retain
the existing `SPACEDOCK_BIN` and operator `SAFEHOUSE_ENV_PASS` composition,
avoid broad inheritance, and reduce the diff/tests to the wrapper-owned proof.

**Cycle 2 — taste and ownership correction (2026-07-12).** The one-entry
terminal-targeting registry, argv scanner, custom `--env-pass` merger, and
ambient-environment test choreography are speculative machinery. Safehouse
already owns named-env composition and exposes the operator-level generic
`SAFEHOUSE_ENV_PASS` configuration. Prefer that zero-code configuration when
it satisfies the desired default. If a Spacedock-built-in default remains
necessary, keep the wrapper-only change to one clean function that constructs
one static terminal-targeting `--env-pass` argument; let Safehouse combine it
with the existing launcher argument. Do not invent a set registry, parse or
merge Safehouse argv, or touch the three front doors. Update the task body and
ACs to express the outcome rather than prescribing a generic data structure;
keep each AC heading on one line so the gate extractor can audit it.

## Stage Report: ideation

- DONE: The captain approved a built-in conditional allowlist for the three
  Zellij targeting variables, retaining Safehouse's native global
  `SAFEHOUSE_ENV_PASS` as the escape hatch for additional variables.
- DONE: The smallest mechanism is the existing Safehouse named environment
  forwarding surface; a new Spacedock config format would duplicate it and make
  security-sensitive allowlisting harder to audit.

## Stage Report: implementation

- DONE: Implement a conditional, shared Safehouse named-env allowlist for ZELLIJ, ZELLIJ_PANE_ID, and ZELLIJ_SESSION_NAME without hardcoding values.
  Product commit `53848e52` forwards only the three names when `ZELLIJ` is present and retains the resolved `SPACEDOCK_BIN` allowance.
- DONE: Add behavior-level scrub/pass-through and non-Zellij argv controls for every wrapped Spacedock host using the shared helper.
  Explicit-parent fixtures cover Claude, Codex, and Pi; the scrubbed executable smoke proves exact `ZELLIJ=0`, pane, and session values plus the stripped control.
- DONE: Preserve Safehouse's native global SAFEHOUSE_ENV_PASS route for arbitrary extras; run focused tests plus gofmt and repository baseline gates.
  `go test ./... -count=1` and `go test ./... -race -count=1` each passed 2168 tests in 17 packages; `ZELLIJ=0 go test ./internal/cli -count=1` passed, and only pre-existing `internal/release/journeydelta.go` remained from `gofmt -l ./cmd ./internal`.

### Summary

Safehouse-wrapped Claude, Codex, and Pi launches now conditionally retain the
current Zellij targeting metadata without manufacturing values or widening the
environment. Existing non-Zellij argv fixtures are terminal-independent, while
the native global Safehouse allowlist remains available for operator extras.

## Stage Report: validation

- DONE: Independently exercised the scrubbed Safehouse executable smoke. With
  `ZELLIJ=0`, `ZELLIJ_PANE_ID=51`, and
  `ZELLIJ_SESSION_NAME=excellent-pheasant`, the child printed those exact three
  values; the otherwise identical stripped control printed all three as unset.
  The same smoke proved the native `SAFEHOUSE_ENV_PASS=EXTRA_TARGET` operator
  extension composes with the built-in trio.
- DONE: Independently ran the conditional argv checks across Claude, Codex, and
  Pi. A Zellij parent emits the single named allowlist
  `SPACEDOCK_BIN,ZELLIJ,ZELLIJ_PANE_ID,ZELLIJ_SESSION_NAME`; a parent without
  `ZELLIJ` retains only the established `SPACEDOCK_BIN` allowance. The helper is
  shared by all three wrapped front doors, and no new Spacedock config surface
  or broad environment pass-through appears in the product diff.
- DONE: Validation gates passed: focused Zellij/Safehouse tests (17 test units),
  `go test ./... -count=1`, and `go test ./... -race -count=1`. `git diff
  --check origin/main...HEAD` was clean. The changed Go files were gofmt-clean;
  repository-wide `gofmt -l ./cmd ./internal` reports only the unrelated,
  pre-existing `internal/release/journeydelta.go`.

### Validation summary

PASSED. The implementation preserves exact inherited Zellij targeting metadata
through the Safehouse scrub boundary only when `ZELLIJ` is present, leaves
non-Zellij launches stable, and retains `SAFEHOUSE_ENV_PASS` as the global
operator extension mechanism.

## Stage Report: implementation (cycle 1)

- DONE: Centralize conditional terminal/session-targeting environment pass-through in the Safehouse wrapper or its immediate generic helper. Use a generic abstraction name so a future tmux allowlist composes; retain the current Zellij trio only as today's data.
  Product commit `5f61c16b` moves composition into `safehouse.Wrap` via generic terminal-targeting sets; Zellij is only the first data set.
- DONE: Restore the ordinary Claude, Codex, and Pi front-door call shapes: do not rethread parent environment or add per-host WithEnv variants solely for Safehouse metadata. Preserve SPACEDOCK_BIN and native SAFEHOUSE_ENV_PASS composition without broad inheritance.
  The combined diff against `origin/main` has no Claude/Codex/Pi product changes; the wrapper merges targeting names with the existing `SPACEDOCK_BIN` allowlist, and the scrubbed smoke preserves `SAFEHOUSE_ENV_PASS=EXTRA_TARGET`.
- DONE: Shrink the diff and tests to wrapper-owned proof, update the task design/acceptance/report for this captain correction, commit the existing worktree, and run proportionate verification.
  Independent review found no actionable issue; `ZELLIJ=0 go test ./... -count=1` and `ZELLIJ=0 go test ./... -race -count=1` each passed 2168 tests in 17 packages, while `gofmt -l ./cmd ./internal` reports only pre-existing `internal/release/journeydelta.go`.

### Summary

The correction makes Safehouse—not each host front door—the sole owner of
terminal/session metadata pass-through. The wrapper combines current Zellij
metadata with the existing least-privilege allowlist and leaves a clean data
extension point for future terminal multiplexers.

## Stage Report: validation (cycle 1)

- DONE: Independently prove the behavior is owned by Safehouse alone: compare against origin/main to confirm Claude, Codex, and Pi production front doors retain their prior shapes, while the wrapper carries generic terminal-targeting metadata and today’s Zellij data.
  `git diff origin/main...5f61c16b` is empty for production Claude/Codex/Pi front doors; only `internal/safehouse` owns the generic `terminalTargetingEnvSet` data and allowlist composition.
- DONE: Exercise the wrapper-owned scrubbed smoke and composition controls; verify exact inherited Zellij values, no manufacture without the targeting parent, and native SAFEHOUSE_ENV_PASS composition without broad inheritance.
  Fresh ambient-Zellij focused controls passed; the scrubbed child observed `0`, `51`, and `excellent-pheasant`, the stripped control observed all unset, and `SAFEHOUSE_ENV_PASS=EXTRA_TARGET` composed without broad inheritance.
- DONE: Perform a smallest-sufficient-mechanism review of the code and tests. Reject any residual host-front-door rewiring, Zellij-shaped API, or disproportionate test/fixture scaffolding; run focused, full, and race evidence before recommending PASSED.
  The production change is a single wrapper helper plus generic data; the three focused controls cover join, public wrap, and real scrubbed-child behavior. Full and race repository gates exited 0; a detached mutation that suppressed forwarding made all three controls fail.

### Summary

PASSED. The cycle-1 correction removes the prior front-door rethreading and
keeps the test mechanism proportionate: one existing scrubbed Safehouse fixture
extended for all three acceptance criteria, plus small wrapper unit controls.

## Stage Report: implementation (cycle 2)

- DONE: Replace one-entry registry and custom argv merger with clean wrapper-only terminal-targeting argument construction; do not touch front doors.
  Product commit `d57d82d` replaces the speculative machinery with the single
  `terminalTargetingEnvArgs` helper in `safehouse.Wrap`. The combined production
  diff still contains no Claude, Codex, or Pi front-door change.
- DONE: Use Safehouse's actual named-env composition contract; preserve built-in Zellij default while leaving SAFEHOUSE_ENV_PASS generic operator extension.
  The helper conditionally supplies one static
  `--env-pass=ZELLIJ,ZELLIJ_PANE_ID,ZELLIJ_SESSION_NAME` argument. Review of
  Safehouse's parser confirmed it deduplicates repeated `--env-pass` names and
  composes them with `SAFEHOUSE_ENV_PASS`; the scrubbed-child smoke proves the
  exact three values, the stripped control, and `EXTRA_TARGET` composition.
- DONE: Update task body and one-line ACs, keep meaningful forwarding/scrub proof, run focused and required tests, write report, commit, and do not push.
  The task body and one-line ACs describe the static wrapper behavior.
  `ZELLIJ=0 go test ./... -count=1` and
  `ZELLIJ=0 go test ./... -race -count=1` each passed 2168 tests in 17
  packages. `gofmt -l ./cmd ./internal` reports only pre-existing unrelated
  `internal/release/journeydelta.go`; the local product commit was intentionally
  not pushed.

### Summary

The cycle-2 correction relies on Safehouse's native named-environment
composition instead of reimplementing it. Independent review found no
actionable issues: the result is one conditional wrapper helper, exact
least-privilege forwarding, and no front-door or configuration-surface change.

## Stage Report: validation (cycle 2)

- DONE: Independently verify the wrapper has only clean static terminal-targeting argument construction: no registry, argv parser/merger, or Claude/Codex/Pi front-door change.
  `d57d82d` replaces the stale registry/merger with wrapper-only `terminalTargetingEnvArgs()`, which conditionally returns one static `--env-pass=ZELLIJ,ZELLIJ_PANE_ID,ZELLIJ_SESSION_NAME`; `git diff origin/main...HEAD` contains no production Claude, Codex, or Pi change.
- DONE: Prove exact forwarding and scrubbed absence against Safehouse named-env composition, with a meaningful adversarial control; reproduce relevant normal/race evidence.
  Focused ambient-Zellij controls passed; a detached audit changed that helper to return `nil`, making both the scrubbed-child exact-forwarding smoke and wrapper contract test fail. `ZELLIJ=0 go test ./... -count=1` and `ZELLIJ=0 go test ./... -race -count=1` completed cleanly; changed Go files are gofmt-clean and `git diff --check origin/main...HEAD` is clean.
- DONE: Cross-check every AC with evidence, including the now-extractable one-line AC headings; report PASSED/REJECTED and do not push or merge.
  AC-1 is the exact-three-values smoke plus stripped control; AC-2 is the absent-parent wrapper control and no-front-door diff; AC-3 is the composition smoke with `SAFEHOUSE_ENV_PASS=EXTRA_TARGET`. The remote PR head remains the older `5f61c16b`; this validation covers local `d57d82d` and does not push or merge it.

### Summary

PASSED for `d57d82d`. The prior `withTerminalTargetingEnvPass` shape was
unnecessarily indirect; the candidate now has the smallest legible form: a
single conditional function that constructs the one Safehouse argument and
leaves composition to Safehouse. PR #499 must be fast-forwarded to this
validated commit before it can be presented as the reviewed change.
