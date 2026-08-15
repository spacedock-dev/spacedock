---
id: x8g3dnqndfa1m85d8ga2cgem
title: Trim version output decoration without a reader
status: ideation
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started: 2026-08-15T02:55:38Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:x8g3dnqndfa1m85d8ga2cgem:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:x8g3dnqndfa1m85d8ga2cgem-backlog-1
              briefing:
                id: briefing:x8g3dnqndfa1m85d8ga2cgem:backlog:attempt-1:revision-1
                digest: sha256:29b040ac52ea6ac3b5ec088df13bd65c14cf47c40a23c4fbd9c953003d046676
                request-digest: sha256:f9aed9bec822660eeb6e0eaac5f049b2e1c6707f2fdcb404ba072ef674111704
                room-ref: ./trim-version-output-decoration/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:x8g3dnqndfa1m85d8ga2cgem:backlog:1
                briefing: briefing:x8g3dnqndfa1m85d8ga2cgem:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:43.914595Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:x8g3dnqndfa1m85d8ga2cgem:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:x8g3dnqndfa1m85d8ga2cgem-ideation-1
              briefing:
                id: briefing:x8g3dnqndfa1m85d8ga2cgem:ideation:attempt-1:revision-1
                digest: sha256:0a4dcfa4226fa4c8c01cdea9b1e77a5e737cbced12103a57c2ecec585350f207
                request-digest: sha256:bf06b0da29d4f2cb13688baa0fff04eb239cbaf21a3c339c31422c1b6da8ba48
                room-ref: ./trim-version-output-decoration/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-15T03:39:58.026643Z"
                reason: Entity gained evidence re-verification section post-prepare (33ccb5479); re-preparing against current bytes; all figures unchanged
---

Remove three decorations from `spacedock --version` output. Each has no reader.

1. The `OS:` line (internal/cli/cli.go:883). The shared contract orders the FO to use `uname -s`, not this line. Doctor computes GOOS itself.
2. The session short-id segment on the Runtime line (cli.go:893). Session matching everywhere reads the env var, never the printed prefix.
3. The `pass --host` remedy suffix on the ambiguous arm (cli.go:910-916). The one ambiguity-affected command, dispatch build, prints its own complete remedy.

Keep the Runtime line itself (host and marker have a recorded live read), the ambiguous reporting arm (the nested-marker leak occurs live), and the Sandbox line (fo-install-gate corroboration).

Captain override recorded here: the OS line was a captain-ratified annotation this cycle. Its ratification record names only speculative consumers. Approval of this entity reverses that annotation.

## Problem

`spacedock --version` prints five lines inside a session and two outside it. Three of those decorations have no reader. Each was verified against HEAD (4d1912a69) by finding every consumer, not by reasoning about intent.

1. **The `OS: <goos>/<goarch>` line** (`internal/cli/cli.go:883`). Its two ratified consumers were "the FO reads the platform" and "user issue reports carry the platform". Neither exists. The shared contract orders the opposite source: `skills/first-officer/references/first-officer-shared-core.md:10` says "Use `uname -s`, not `doctor`/`OS:`" — and `uname -s` is permanently load-bearing there, because the binary-absent class produces no `--version` output at all. `doctor` never prints an OS line; it computes `runtime.GOOS` itself to pick the brew-vs-curl remedy (`internal/contract/contract.go:115,133`). No issue template, CI lane, or skill quotes the line. `.github/workflows/install-e2e.yml:55` reads `--version` but globs `spacedock ?*` across the whole output, so it is indifferent to the line's presence.

2. **The `, session <8hex>` segment on the Runtime line** (`internal/cli/cli.go:912-913` — the seed's line reference for this item and the next are transposed; the segment lives in `runtimeLine`, the ambiguous arm at 893). Every consumer that matches a session reads the environment variable directly and never parses the printed prefix: `internal/dispatch/build.go:756` and `internal/dispatch/reconcile.go:139` both call `os.Getenv("CLAUDE_CODE_SESSION_ID")`; `skills/integration/testdata/version_gate_flow.sh:16` reads `${CLAUDE_CODE_SESSION_ID:-${CODEX_THREAD_ID:-fixture}}`. The printed eight characters are a human convenience with no recorded human use.

3. **The ` — pass --host` suffix on the ambiguous arm** (`internal/cli/cli.go:893`). It is a remedy attached to a surface that has no fault to remedy: `--version` reports ambiguity and exits 0 by design. The one command the ambiguity actually blocks prints its own complete remedy, which the suffix does not match and cannot replace — driven live at HEAD, `dispatch build` under two markers emits `ambiguous runtime host sources: multiple runtime markers are set (CODEX_THREAD_ID, CLAUDECODE); pass --host claude, codex, or pi`. That names all three valid values; the `--version` suffix names none.

## Proposed approach

Delete the three decorations, then delete the machinery each one was the sole consumer of. Stopping at the printed characters would leave exactly the shape this program exists to remove — machinery with no consumer — in a package that previously had one.

**The three output changes.** Drop the `OS:` `Fprintf`; drop `— pass --host` from the ambiguous format string; render the resolved Runtime line inline as `host + " (" + strings.Join(markers, ", ") + ")"`.

**The orphaned machinery.** With the session segment gone, `runtimehost.ShortID` and `shortIDLen` have no caller and `Detect`'s `identity` return has no consumer (`internal/dispatch/build.go:276` already discards it with `_`). Delete `ShortID`, `shortIDLen`, the `identity` column from `markerTable`, and the `identity` return value, taking `Detect` from four returns to three. `runtimeLine` also goes: with the identity branch removed it is a single-expression helper with one caller, and the branch was its whole reason to exist. `internal/runtimehost` is an internal package with exactly two callers, so no external consumer can exist.

**Alternatives considered.** Removing only the printed characters and keeping `ShortID`/`identity` was rejected: it converts one-consumer machinery into zero-consumer machinery, which is the fault this entity is filed against. Keeping `runtimeLine` as a two-argument marker joiner was rejected: a single-expression helper called once is harder to read than the expression.

**The one addition.** `TestVersionSessionRender` gets a nine-line loop asserting that `\nOS: `, `, session `, and `pass --host` are absent, checked before the exact-match comparison. The simplest alternative is to rely on the exact-match `want` strings alone, which do catch a reintroduction — proven by mutation below. That alternative is insufficient against the specific regression path in this file: the `want` strings are hand-maintained literals, so a future editor who reintroduces a decoration and regenerates the literals to match gets a green suite. The named loop fails with the decoration's name and the reason ("it has no reader"), forcing the reintroduction to argue with the rationale rather than update a string. Both behaviours were exercised (see the spike record). This is the only insertion of substance in an otherwise pure deletion; cutting it would not weaken AC coverage, only the failure message.

**Doc diff (recorded here; implementation applies it).**

```diff
--- a/docs/site/reference/command-reference.md
+++ b/docs/site/reference/command-reference.md
@@ -1,22 +1,20 @@
 # Command reference
 
-The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version and the host OS/arch, and — inside an agent session — that session's runtime and sandbox state). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.
+The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version, and — inside an agent session — that session's runtime and sandbox state). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.
 
 ## --version
 
-`spacedock --version` reports the binary version and the host OS/arch, and — when it is running inside an agent session — that session's runtime and sandbox state. Outside any session it prints two lines:
+`spacedock --version` reports the binary version and — when it is running inside an agent session — that session's runtime and sandbox state. Outside any session it prints one line:
 
-    spacedock 0.26.0
-    OS: darwin/arm64
+    spacedock 0.27.0
 
-Inside a session it also names the host OS/arch, the runtime it detected, the marker that proved it, which session this is, and whether this process is running inside a sandbox:
+Inside a session it also names the runtime it detected, the markers that proved it, and whether this process is running inside a sandbox:
 
-    spacedock 0.26.0
-    OS: darwin/arm64
-    Runtime: claude (CLAUDECODE, session afd74765)
+    spacedock 0.27.0
+    Runtime: claude (CLAUDECODE)
     Sandbox: inside (agent-safehouse)
     contract 3
 
-The session identifier is the first eight characters of the host's own session id — the same prefix Claude Code uses to name `~/.claude/teams/session-afd74765` — so you can tell two concurrent sessions apart and match one against its state on disk. Hosts that do not expose a session id, such as pi, omit it:
+A host can set more than one marker; pi sets two, and both are listed:
 
     Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)
@@ -24,7 +22,9 @@
 When markers for more than one runtime are set — a nested session can leak them — it reports the ambiguity rather than guessing, and still exits 0:
 
-    Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE) — pass --host
+    Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE)
 
-Being outside every runtime is a normal state, not a fault — it means a human at a terminal. There is no `Runtime:` line at all in that case, because there is no session to report: the output is the two lines shown above (the version line plus the `OS:` line).
+That report is informational. The one command the ambiguity actually blocks is `spacedock dispatch build`, which refuses and prints its own remedy naming the three valid values: `pass --host claude, codex, or pi`.
+
+Being outside every runtime is a normal state, not a fault — it means a human at a terminal. There is no `Runtime:` line at all in that case, because there is no session to report: the output is the single version line shown above.
@@ -34,5 +34,5 @@
-`spacedock status --boot` reports the same three. The pre-launch banner answers …
+`spacedock status --boot` reports the sandbox state too, on its `SANDBOX:` line. The pre-launch banner answers …

--- a/skills/first-officer/references/first-officer-shared-core.md
+++ b/skills/first-officer/references/first-officer-shared-core.md
@@ -8,5 +8,5 @@
-   - **Binary absent:** … Use `uname -s`, not `doctor`/`OS:`. Linux: …
+   - **Binary absent:** … Use `uname -s`, not `doctor`. Linux: …
```

Two notes on that diff. The `status --boot` sentence is corrected rather than merely de-referenced: driving the binary at HEAD shows `--boot` emits only `SANDBOX: inside (agent-safehouse)` and no runtime or OS line, so "reports the same three" was already wrong before this change; the removal breaks its referent, which is what brings it in scope. The example version numbers move from a stale `0.26.0` to `0.27.0` because the example blocks are rewritten anyway.

**Spike record.** The riskiest unverified claim was that the identity plumbing can be removed without a compile or test consequence somewhere unsearched. The whole cut was exercised in a throwaway detached worktree at HEAD (not committed; ideation ships no branch):

- `go build ./...` clean; `gofmt -l ./cmd ./internal` clean.
- Live binary, in-session: 4 lines (`spacedock 0.27.0+dev` / `Runtime: claude (CLAUDECODE)` / `Sandbox: inside (agent-safehouse)` / `contract 3`); outside every runtime: 1 line. HEAD's binary on the same machine: 5 and 2.
- `dispatch build` driven through stdin JSON under two markers still prints its complete remedy verbatim — the keep-boundary that makes the `--version` suffix redundant rather than merely unread.
- `go test ./...`: `internal/cli` (439s on a loaded machine), `internal/runtimehost`, `internal/dispatch`, `internal/contractlint`, and `skills/integration` all pass except `TestCodexResolveManifestAgainstInstalledHost`, which fails identically on **unmodified HEAD** in this sandbox (`Failed to read config file /Users/clkao/.codex/config.toml: Operation not permitted`) — environmental, not caused by the cut.
- Mutation 1: reintroducing the `OS:` line turns all five `TestVersionSessionRender` cases red on the exact-match.
- Mutation 2: reintroducing ` — pass --host` **and** regenerating the `want` strings to match still fails, with `--version rendered a --host remedy suffix ("pass --host") — it has no reader`. This is what justifies the nine-line addition.
- Measured diff: 7 files, 91 insertions, 256 deletions (net **-165**).

## Out of scope

The Runtime detection redesign. The Sandbox line. The `contract 3` sentinel. Contract prose beyond the docs examples and the one shared-core clause whose referent this change breaks.

`skills/integration/version_gate_fixture_test.go:327,345` keeps its `echo "OS: linux/amd64"` lines. Those are captive fake-binary *inputs*, not assertions about this binary's output, and the intervening line is what proves the gate's `^Sandbox:` parse is prefix-anchored rather than position-dependent. Removing them would weaken the fixture.

## Coordination

- `remove-startup-capability-probe` edits the same two files from a different angle: the Startup section of `first-officer-shared-core.md` (its "Compatible minor" and "Missing capability" bullets) and `skills/integration/version_gate_fixture_test.go`. Different lines, so no semantic conflict, but a textual merge conflict is likely if both land unrebased. Whichever merges second should rebase.
- `remove-gate-validate-subcommand` edits `internal/cli/cli.go` in the `gate validate` branches and usage text — a different function from `printVersion`. Low risk.
- `retire-requires-contract-sentinel` is **not** an overlap despite the name collision: it retires the `requires-contract` plugin-manifest field, not the `contract 3` token `printVersion` emits. Nothing in this entity touches that token.

## Expected surface and tolerance

Estimate net LOC change: **-165** (91 insertions, 256 deletions), across **7 files** — measured from the spike, not guessed.

| File | Net |
|---|---|
| `internal/cli/version_session_test.go` | -70 |
| `internal/runtimehost/runtimehost_test.go` | -70 |
| `internal/runtimehost/runtimehost.go` | -41 |
| `internal/cli/cli.go` | -23 |
| `docs/site/reference/command-reference.md` | -2 |
| `internal/dispatch/build.go` | 0 |
| `skills/first-officer/references/first-officer-shared-core.md` | 0 |

Tolerance: ±30 lines, ±1 file. A file count above 8 or a net delta above -100 means the cut grew or shrank beyond what was gated; stop and re-present.

**Observable semantics changed.** Declared deliberately, since the diff is small and the semantics are not:

1. `--version` stdout shape changes in every case. Outside every runtime: 2 lines → 1. Inside a session: 5 → 4. The resolved Runtime line loses its `, session <8hex>` segment. The ambiguous arm loses its ` — pass --host` suffix.
2. Package-internal Go API: `runtimehost.Detect` returns 3 values instead of 4; `runtimehost.ShortID` and `shortIDLen` cease to exist. `internal/` forbids an external caller, and both in-repo callers are updated in the same change.
3. Unchanged: every exit code (including the ambiguous arm's 0), every command grammar, every stored format, every authority rule, and every other command's output.

## Acceptance criteria

**AC-1 — `spacedock --version` emits one line outside every runtime and four inside a session, down from two and five.**
This is the end-value the entity exists for, measured against a baseline that can move the wrong way: a partial cut, or a decoration added back during implementation, lands on different numbers.
Verified by: building the binary and counting output lines in both shapes, against the same counts taken from `origin/main`'s binary on the same machine. Baseline recorded above (5 and 2 at HEAD).

**AC-2 — No reader-less decoration is reachable in any `--version` shape.**
Verified by: `TestVersionSessionRender`'s five exact-match `want` strings plus the named absence checks for `\nOS: `, `, session `, and `pass --host`. Falsifying change: reintroducing any one of the three turns the test red — and does so even if the `want` literals are regenerated in step, which was exercised.

**AC-3 — The cumulative line delta against `origin/main` is negative.**
Verified by: `git diff --stat origin/main` on the merge branch. Falsifying change: net-additive machinery smuggled in alongside the deletion.

**AC-4 — Every named keep-boundary still holds.**
The Runtime line still names host and markers; the ambiguous arm still reports and still exits 0; the `Sandbox:` line is byte-identical; `dispatch build` still refuses on ambiguity with its own complete three-value remedy; the FO contract still sources the OS from `uname -s`.
Verified by: `TestVersionSessionRender` (Runtime and Sandbox lines), `TestVersionAmbiguousMarkersExitZero` (exit code through the real `Run`), `internal/dispatch` `build_host_test.go:38` (the unchanged remedy string), and `internal/contractlint` `TestVersionGateProseOSAwareHint` (the `uname -s` token). Falsifying change: any of these four surfaces drifting turns its own test red.

**AC-5 — The suite is green.**
Verified by: `go test ./...` and `go test ./... -race`, with `TestCodexResolveManifestAgainstInstalledHost` excepted only where it also fails on unmodified HEAD in the same environment; on CI, where `~/.codex/config.toml` is readable, no exception applies.

## Test plan

No new test files. The work is deletion plus retargeting existing assertions; cost is low and no live workflow run is needed, because every claim is either output bytes from the binary or a unit-level render.

- `internal/cli/version_session_test.go` — retarget the five render cases to the new shape, add the named absence loop, drop `TestVersionRuntimeLineDistinguishesConcurrentSessions` (it exists solely to prove the session segment distinguishes two sessions, which is the behaviour being removed), and drop the `claude-without-session-id` case (its distinguishing purpose was the omitted-identity path). Serves AC-2 and AC-4.
- `internal/runtimehost/runtimehost_test.go` — drop `TestDetectIdentity` and `TestShortID`, fix `Detect`'s arity in `TestDetectMarkerMatrix`. The marker-matrix and same-host-non-ambiguity coverage is untouched. Serves AC-4.
- Binary-level line counts for AC-1, taken from both the built binary and `origin/main`'s, in both output shapes.
- `internal/contractlint` and `skills/integration` are run unchanged: they are the guard that the FO prose edit and the retained gate fixtures did not drift. Both passed in the spike.
- `go test ./... -race` was not run in the spike (the non-race `internal/cli` pass alone took 439s under load); it runs at implementation for AC-5.

## Stage Report: ideation

- DONE: Keep-boundary confirmed: Runtime host and marker, the ambiguous reporting arm, and the Sandbox line stay
  All three verified live at HEAD and in the spike binary; `dispatch build` re-driven through stdin JSON to confirm it still prints the complete `pass --host claude, codex, or pi` remedy the `--version` suffix does not replace.
- DONE: Version output assertions and docs examples updated to the new shape
  Concrete diff for `command-reference.md` and the one `first-officer-shared-core.md` clause recorded in the body per the stage contract (ideation records, implementation applies); the same diff was applied and proven in the throwaway spike.
- DONE: Output-shape semantic change declared in the estimate
  Three observable semantics declared: `--version` shape (2->1 lines outside, 5->4 inside, two segments dropped), `runtimehost.Detect` 4 returns -> 3 with `ShortID`/`shortIDLen` deleted, and an explicit unchanged list (exit codes, grammar, stored formats, authority).

### Summary

Confirmed all three removals against HEAD by locating every consumer rather than reasoning about intent: the `OS:` line is contradicted by the FO contract's own `uname -s` order and unused by `doctor`, which computes GOOS itself; all three session-matching consumers read the env var directly; and the ambiguous arm's `pass --host` suffix is a remedy for a surface that exits 0, while the one blocked command prints a strictly better remedy. Scope grew past the seed's three items to the machinery each was the sole consumer of — `ShortID`, `shortIDLen`, the `markerTable` identity column, `Detect`'s `identity` return, and `runtimeLine` — because stopping at the printed characters would leave zero-consumer machinery behind, the exact fault this program removes; that expansion is declared for the gate to ratify, with the narrower alternative recorded. The whole cut was spiked in a throwaway detached worktree: builds clean, both output shapes exercised live, affected packages green (one pre-existing sandbox-environmental failure reproduced on unmodified HEAD), and both the exact-match guard and the new named absence checks falsified by mutation. Measured surface is 7 files and net -165 lines, replacing the seed's -NN/~3-files placeholder.

### Evidence re-verification (FO shared-scratchpad advisory, 2026-08-15)

The original spike worktree was entity-named, but two binaries it produced sat at bare shared-scratchpad paths (`scratchpad/sd` for the HEAD baseline, `/tmp/sd-spike` for the cut) where a sibling ensign could have overwritten them. Every number those binaries produced was re-derived from scratch in a fully slug-named worktree (`scratchpad/spike-trim-version-output-decoration`), with both binaries built inside it. All figures reproduce unchanged:

- Baseline `--version`: 5 lines in-session, 2 outside. Cut: 4 and 1. (AC-1 baseline intact.)
- Diffstat: 7 files, 91 insertions, 256 deletions, net -165. (Estimate intact.)
- `dispatch build` under two markers still prints `ambiguous runtime host sources: multiple runtime markers are set (CODEX_THREAD_ID, CLAUDECODE); pass --host claude, codex, or pi` — byte-identical. (Keep-boundary intact.)
- `internal/runtimehost`, `internal/dispatch`, `internal/contractlint`, `skills/integration`, and `internal/cli -run TestVersion` all pass.
- Mutation 1 (reintroduce the `OS:` line): fails with `--version rendered an OS/arch line ("\nOS: ") — it has no reader`. Mutation 2 (reintroduce ` — pass --host` and regenerate the `want` literals in step): fails with `--version rendered a --host remedy suffix ("pass --host") — it has no reader`. Both now fail on the named check rather than the exact match, because the absence loop runs first.

Baseline commit note: the body cites HEAD as 4d1912a69; HEAD has since advanced to ef8f55c83 (two workflow-doc commits). `git diff 4d1912a69..HEAD` over all seven target files is empty, so the cut and every cited number apply unchanged.
