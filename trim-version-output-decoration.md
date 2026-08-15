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
            - id: gate-attempt:x8g3dnqndfa1m85d8ga2cgem-ideation-2
              briefing:
                id: briefing:x8g3dnqndfa1m85d8ga2cgem:ideation:attempt-2:revision-1
                digest: sha256:6e30054ac3fde139a894d01bf35b9a142abecab1f06c8a91d4e358e3cd1b1b1f
                request-digest: sha256:0ca46318fe779c2e22774ab54d918b443eb15751254104524e1fab9594ccb4df
                room-ref: ./trim-version-output-decoration/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:x8g3dnqndfa1m85d8ga2cgem:ideation:2
                briefing: briefing:x8g3dnqndfa1m85d8ga2cgem:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-15T04:07:34.176429Z"
                decision: revise
                reason: 'Captain revise 2026-08-15: keep the OS line (bug-report value is real); Sandbox line must not need the word inside while keeping the install-gate corroboration parse working; investigate retiring the contract-3 sentinel per its documented retirement condition with minor version as the sole requirement; the session-segment and pass-host cuts stand'
---

Trim `spacedock --version` output to what has a reader, and give the Sandbox line a value shape that does not repeat its own label.

1. ~~The `OS:` line.~~ **WITHDRAWN by captain revise 2026-08-15.** The line STAYS. The captain declares the bug-report consumer real — people paste `--version` output into issues — which supersedes the cycle-1 finding that its only named consumers were speculative. The machinery it solely consumes (`runtime.GOOS`/`GOARCH` and the `runtime` import in cli.go) stays with it.
2. The session short-id segment on the Runtime line (cli.go:912-913). Session matching everywhere reads the env var, never the printed prefix. **Stands approved in direction.**
3. The `pass --host` remedy suffix on the ambiguous arm (cli.go:893). The one ambiguity-affected command, dispatch build, prints its own complete remedy. **Stands approved in direction.**
4. **ADDED by the same revise.** The Sandbox line must not need the word "inside". After the `Sandbox: ` label, `inside (agent-safehouse)` says "inside" twice over and `not sandboxed (…)` says "sandbox" twice. Both arms get a new value shape, and the install-gate corroboration consumer must keep working.
5. **ADDED by the same revise.** Investigate retiring the `contract 3` sentinel against its own documented retirement condition (cli.go:855), and fold in the one-line correction to the wrong pin attribution at cli.go:853.

Keep the Runtime line itself (host and markers have a recorded live read), the ambiguous reporting arm (the nested-marker leak occurs live), the `OS:` line (item 1), and the Sandbox line itself — item 4 changes its VALUE, not its existence.

## Problem

Three things on `spacedock --version`, verified against HEAD by finding every consumer rather than reasoning about intent. The `OS:` line is no longer among them — see the withdrawal above.

1. **The `, session <8hex>` segment on the Runtime line** (`internal/cli/cli.go:912-913`). Every consumer that matches a session reads the environment variable directly and never parses the printed prefix: `internal/dispatch/build.go:756` and `internal/dispatch/reconcile.go:139` both call `os.Getenv("CLAUDE_CODE_SESSION_ID")`; `skills/integration/testdata/version_gate_flow.sh:16` reads `${CLAUDE_CODE_SESSION_ID:-${CODEX_THREAD_ID:-fixture}}`.

2. **The ` — pass --host` suffix on the ambiguous arm** (`internal/cli/cli.go:893`). A remedy on a surface with no fault to remedy: `--version` reports ambiguity and exits 0 by design. The one command the ambiguity blocks prints a strictly better remedy — driven live, `dispatch build` under two markers emits `ambiguous runtime host sources: multiple runtime markers are set (CODEX_THREAD_ID, CLAUDECODE); pass --host claude, codex, or pi`. That names all three valid values; the `--version` suffix names none.

3. **The Sandbox line's value repeats its own label.** `Sandbox: inside (agent-safehouse)` and `Sandbox: not sandboxed (safehouse available)` both spend their first token restating the label. The value should carry the fact the label is asking for.

## Proposed approach

**Session segment, and the machinery it solely consumes.** Render the resolved Runtime line inline as `host + " (" + strings.Join(markers, ", ") + ")"`. That orphans `runtimehost.ShortID`, `shortIDLen`, the `identity` column of `markerTable`, and `Detect`'s `identity` return (`internal/dispatch/build.go:276` already discards it with `_`), so all four go and `Detect` drops from four returns to three. `runtimeLine` goes with them: the identity branch was its whole reason to exist. Keeping them would convert one-consumer machinery into zero-consumer machinery, which is the fault this entity was filed against. `internal/runtimehost` is internal with exactly two callers, so no external consumer can exist.

**Ambiguous arm.** Drop `— pass --host` from the format string. The arm still reports, still lists markers, still exits 0.

**Sandbox value shape.** The sandboxed arm becomes the sandbox's NAME alone; the unsandboxed arms lead with `none`:

    Sandbox: agent-safehouse
    Sandbox: none (safehouse available)
    Sandbox: none (safehouse not installed)

The name-alone arm is what removes the word "inside" without inventing a synonym for it, and `none` gives consumers a first token to classify on without matching a relationship word. This changes `status --boot`'s `SANDBOX:` line too — both surfaces render the same `safehouse.SessionState`, and that shared render is the reason the change is one function rather than two. The pre-launch banner's `Sandbox:` line renders `LaunchState`, a different question ("will this launch wrap?"), and is out of scope; see the flag below.

**The install-gate corroboration consumer, and why it needed real work.** `skills/integration/testdata/version_gate_flow.sh:44` classified the parsed value with the glob `inside*`. The rename breaks it — but so would any shape change, and the deeper problem is that the gate runs `--version` on *whatever binary is installed*, which may predate these skills and still emit the old shape. So the parse is rewritten to classify on the NOT-sandboxed shapes of both eras and treat anything else as a sandbox name:

    case "$sbline" in
    "" | none* | "not sandboxed"*) …DISAGREEMENT… ;;
    *) …agrees… ;;
    esac

This is strictly more correct than the original, which would have mis-classified a future sandbox whose name did not begin with "inside". `skills/first-officer/references/fo-install-gate.md:22` needs no edit: it names the `^Sandbox: ` prefix and the env-wins rule, never a value literal — checked, not assumed.

**Alternatives considered.** For the Sandbox arms, `Sandbox: no (…)` and `Sandbox: unsandboxed (…)` were rejected: both are still relationship words, so they trade "inside" for a synonym and leave the consumer matching on prose rather than on a name-or-`none` distinction. Retargeting the gate-flow glob from `inside*` to `none*` (single-era) was rejected because it silently breaks against every already-installed older binary — the case the dual-era arms now pin.

**The contract-3 sentinel: investigated, NOT retired this cycle.** See the dedicated finding section below. The one-line pin-attribution correction at cli.go:853 folds in regardless.

## Finding: contract-3 retirement condition is not yet met

The retirement condition at `internal/cli/cli.go:855` asks whether any plugin or binary predating #468 can still be running, and directs that this be settled by querying the Homebrew formula and the marketplace rather than guessing. Both queries were run.

**The boundary.** #468 ("Replace the contract-integer compatibility gate with minor-version coupling") merged 2026-07-03 as `511dae11e`. By ancestry, `v0.24.0-pre2` is the FIRST release containing it; `v0.23.0` and `v0.24.0-pre1` and everything older are pre-#468.

**Distribution query — both channels are clean.**

| Channel | Serves | Verdict |
|---|---|---|
| Homebrew cask `spacedock` | 0.26.0 | post-#468 |
| Homebrew cask `spacedock@next` | 0.27.0-pre4 | post-#468 |
| Marketplace stable (`spacedock-dev/marketplace` main → plugin `spacedock`, `ref: stable`) | `origin/stable` @ `ca136f83a` | contains #468 |
| Marketplace edge (`@edge` → `ref: next`) | `origin/next` @ `80ded63d7` | contains #468 |

Homebrew casks serve exactly one version each and no pinned-old cask exists, so no install path can obtain a pre-#468 binary. Both marketplace channels track BRANCHES whose heads contain #468 (their manifest `version` fields — 0.20.0 stable, a calendar string on edge — are stale metadata, not what resolves), so no install or update path can obtain integer-era prose.

**But extinction is NOT proven, and the blocking population is demonstrable.** The plugin cache never prunes. On this machine right now, `~/.claude/plugins/cache/spacedock/spacedock/` holds 0.19.1, 0.19.9, 0.20.0, 0.20.2, 0.22.0, 0.23.0 and `spacedock-edge/spacedock/` holds 0.23.0-pre and 0.24.0-pre1 — eight pre-#468 plugin versions, physically present and loadable. Any session that has not updated its plugin since the stable channel moved past #468 still loads integer-era prose from that cache. Removing the token there produces exactly the failure its comment predicts: old prose finds no token, may reason "nothing to check, proceed", and runs against a binary four minors too new instead of aborting with "update the plugin".

**Earliest retirement point.** The population is strictly draining — no channel can replenish it, and any user who updates leaves it permanently. The clock started 2026-07-04, when `v0.24.0` became the first stable release carrying post-#468 prose. The token is removable once no cached pre-#468 plugin can still resolve, which needs either a plugin-version telemetry read or a stated staleness horizon; neither exists today, and inventing one is a separate decision, not this entity's. Recommendation: keep the token, and re-run this same two-channel query at the 0.28 cut.

## Out of scope

The `OS:` line (captain-kept; item 1 above). Retiring the `contract 3` sentinel (investigated above; blocked). The Runtime detection redesign.

`skills/integration/version_gate_fixture_test.go:327,345` keeps its old-shape captive scripts. They are fake-binary *inputs*, not assertions about this binary's output, and they are now the pre-rename half of the dual-era corroboration proof — rewriting them to the new shape would delete the coverage that matters most.

**Flagged, not fixed:** `internal/cli/frontdoor.go:209` renders the pre-launch banner's `Sandbox:` line from `LaunchState`, which keeps `inside (agent-safehouse) — launching without re-wrapping`. After this change the two surfaces disagree in style. That is a different question on a different surface and the revise named the `--version` line, so it stays out; if the captain wants the banner aligned it is a one-function follow-up. Separately, `internal/cli/dev_version_test.go:53` carries a stale COMMENT claiming line 1 reads `spacedock <version>+dev (contract 3)`; the token moved off line 1 long ago. Not an assertion, not touched.

## Coordination

- `remove-startup-capability-probe` edits `skills/integration/version_gate_fixture_test.go` (its `withdrawCapability`/`INVOCATION_LOG` plumbing) and `skills/integration/testdata/version_gate_flow.sh` (its `REQUIRED_CAPABILITY` probe block). I edit different regions of both files — the corroboration `case` arm and a new subtest — so no semantic conflict, but a textual merge conflict is likely. Whichever merges second rebases. This entity no longer touches `first-officer-shared-core.md` at all, since the OS line stays, which removes the overlap that existed in cycle 1.
- `remove-gate-validate-subcommand` edits `internal/cli/cli.go` in the `gate validate` branches — a different function from `printVersion`. Low risk.
- `retire-requires-contract-sentinel` is still NOT an overlap despite the name: it retires the `requires-contract` plugin-manifest field, not the `contract 3` token. Worth restating now that this entity has its own contract-3 finding — the two concern different mechanisms and neither blocks the other.

## Expected surface and tolerance

Measured from the spike, not estimated: **net -64** (190 insertions, 254 deletions) across **11 files**.

| File | Net | Why |
|---|---|---|
| `internal/runtimehost/runtimehost_test.go` | -70 | identity + ShortID tests |
| `internal/runtimehost/runtimehost.go` | -41 | identity column, return, ShortID |
| `internal/cli/cli.go` | -18 | runtimeLine, suffix, comment fix |
| `internal/cli/version_session_test.go` | -13 | retarget + new sandbox test |
| `docs/site/reference/command-reference.md` | -6 | examples and prose |
| `internal/status/boot_sandbox_test.go` | 0 | 6 expected values |
| `internal/safehouse/state_test.go` | 0 | 4 expected values |
| `internal/dispatch/build.go` | 0 | Detect arity |
| `skills/integration/testdata/version_gate_flow.sh` | +2 | dual-era case arm |
| `internal/safehouse/state.go` | +6 | new shape + rationale |
| `skills/integration/version_gate_fixture_test.go` | +33 | dual-era corroboration arms |

Tolerance: ±40 lines, ±2 files. The net is far less negative than cycle 1's -165 because the OS-line deletion is withdrawn and the Sandbox work is a rename plus new consumer coverage rather than a removal.

**Observable semantics changed.**

1. `--version` Runtime line loses its `, session <8hex>` segment; the ambiguous arm loses ` — pass --host`. Line COUNTS are unchanged in both shapes (2 outside, 5 in-session) because the OS line stays — so byte size, not line count, is what moves.
2. `Sandbox:` value shape changes on `--version` **and** `SANDBOX:` on `status --boot`. This is wider than cycle 1 declared and is the direct consequence of both surfaces sharing `SessionState`.
3. The install-gate corroboration parse changes from single-era to dual-era classification — a behavior change in the gate flow, not just a retarget.
4. Package-internal Go API: `runtimehost.Detect` returns 3 values instead of 4; `ShortID`/`shortIDLen` cease to exist.
5. Unchanged: every exit code (including the ambiguous arm's 0), command grammar, stored formats, authority rules, the `OS:` line, the `contract 3` token, and the pre-launch banner.

## Acceptance criteria

**AC-1 — `--version` output gets smaller in every session shape, with no line removed.**
The end-value, measured against a baseline that can move the wrong way: a partial cut, a re-added decoration, or a verbose replacement for "inside" all land on different numbers. Baselines measured on the same machine: in-session 130 bytes → 103; ambiguous 148 → 123; line counts 5 and 2 in both.
Verified by: building both binaries and byte-counting each shape.

**AC-2 — No trimmed decoration is reachable in any `--version` shape, and the Sandbox line names its sandbox.**
Verified by: `TestVersionSessionRender`'s five exact-match `want` strings plus named absence checks for `, session `, `pass --host`, and `inside`; and `TestVersionSandboxLineNamesTheSandbox`, which pins the bare-name arm and both `none` arms independently of those literals. Falsifying change: reverting `SessionState` fails by name even with every `want` regenerated.

**AC-3 — The install-gate corroboration classifies BOTH value-shape eras correctly.**
This is the consumer the revise named. Verified by: `TestGateFlowSandboxEnvMarker`'s four arms — old-shape sandboxed agrees, old-shape unsandboxed disagrees, new-shape sandboxed agrees, new-shape unsandboxed disagrees. Falsifying change: restoring the single-era `inside*` glob turns the new-shape-sandboxed arm red (exercised).

**AC-4 — Every keep-boundary still holds.**
`OS: <goos>/<goarch>` is still line 2 in both shapes; the Runtime line still names host and markers; the ambiguous arm still exits 0; `dispatch build` still prints its complete three-value remedy; `contract 3` is still emitted below the Sandbox line.
Verified by: `TestVersionSessionRender` (OS line 2, Runtime, Sandbox), `TestVersionAmbiguousMarkersExitZero` (exit code through the real `Run`), `TestVersionContractTokenPlacement` (token still present and still last), and `internal/dispatch` `build_host_test.go:38` (the unchanged remedy string).

**AC-5 — The contract-3 retirement question is settled on the record.**
Not "the token was removed" — the investigation's end-state is a determination. The body records the #468 boundary, both channel queries with their resolved heads, the named blocking population, and the earliest retirement point.
Verified by: the finding section above, and the corrected cli.go:853 comment naming `version_session_test.go` as the real pin. Falsifying check: `grep -r "contract 3\|frozenContractToken" internal/contractlint/` returns nothing, which is what makes the old attribution wrong.

**AC-6 — The suite is green.**
Verified by: `go test ./...` and `go test ./... -race`, excepting `TestCodexResolveManifestAgainstInstalledHost` only where it also fails on unmodified HEAD in the same environment (it reads `~/.codex/config.toml`, which this sandbox denies); no exception applies on CI.

## Test plan

One new test and one new subtest group; everything else is retargeting existing assertions. No live workflow run is needed — every claim is either output bytes from the binary or a unit/fixture-level render.

- `internal/cli/version_session_test.go` — retarget five render cases (OS line retained, Sandbox reshaped), swap the absence loop's `\nOS: ` token for `inside`, add `TestVersionSandboxLineNamesTheSandbox`, drop `TestVersionRuntimeLineDistinguishesConcurrentSessions` (it exists solely to prove the session segment distinguishes sessions — the behaviour being removed). Serves AC-1, AC-2, AC-4.
- `skills/integration/version_gate_fixture_test.go` — add the two new-shape corroboration arms beside the two existing old-shape ones. Serves AC-3; this is the only substantive insertion.
- `internal/safehouse/state_test.go` (4 values) and `internal/status/boot_sandbox_test.go` (6 values) — retarget to the new shape. `internal/status/harness_test.go:60` redacts the SANDBOX line by regex and needs no change. Serves AC-2 and semantic change 2.
- `internal/runtimehost/runtimehost_test.go` — drop `TestDetectIdentity` and `TestShortID`, fix `Detect`'s arity. Marker-matrix and same-host-non-ambiguity coverage untouched.
- Byte counts for AC-1 from both binaries in both shapes.
- `go test ./... -race` was not run in the spike (the non-race `internal/status` and `internal/dispatch` alone took 310s and 268s under load); it runs at implementation for AC-6.

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

## Stage Report: ideation (cycle 2)

- DONE: Keep-boundary confirmed: Runtime host and marker, the ambiguous reporting arm, and the Sandbox line stay
  All hold, plus the `OS:` line now that the captain reinstated it; `dispatch build`'s complete three-value remedy re-driven live, and `contract 3` still emitted below the Sandbox line.
- DONE: Version output assertions and docs examples updated to the new shape
  Spike retargets `version_session_test.go` (5 cases), `safehouse/state_test.go` (4 values), `status/boot_sandbox_test.go` (6 values) and the `command-reference.md` examples; all applied and green in the spike.
- DONE: Output-shape semantic change declared in the estimate
  Five semantics declared, two of them wider than cycle 1: the Sandbox value change hits `status --boot` as well as `--version` (shared `SessionState`), and the gate-flow corroboration moves from single-era to dual-era classification.
- DONE: Captain directive 1 — keep the OS line
  OS-line deletion withdrawn from the cut along with its `runtime` import; seed item 1 struck in place with the reversal and its authority recorded; the `first-officer-shared-core.md` edit is dropped entirely, which also removes cycle 1's overlap with remove-startup-capability-probe.
- DONE: Captain directive 2 — Sandbox line without "inside", corroboration consumer verified
  New shape is name-or-`none`; before/after recorded in the body. The consumer needed real work: `version_gate_flow.sh:44` classified on the literal glob `inside*`, so the parse was rewritten to span both value-shape eras, and two new fixture arms pin the post-rename half. `fo-install-gate.md:22` checked and needs no edit — it names the `^Sandbox: ` prefix and the env-wins rule, never a value literal.
- DONE: Captain directive 3 — contract-3 retirement investigated, pin attribution corrected
  Both channels queried and recorded: Homebrew serves 0.26.0/0.27.0-pre4, marketplace stable and edge track branches whose heads contain #468 (first containing release `v0.24.0-pre2`). Extinction NOT proven — eight pre-#468 plugin versions sit unpruned in the local cache — so the token stays, with the blocking population and earliest retirement point on the record. The cli.go:853 attribution is corrected: `internal/contractlint` contains no reference to the token at all.

### Summary

The revise inverted one of the three original cuts and added two investigations, so the shape changed substantially: net is now -64 across 11 files rather than -165 across 7, and line counts no longer move at all, which forced AC-1 from a line measure to a byte measure (in-session 130 to 103, ambiguous 148 to 123). The Sandbox rename turned out to be the load-bearing work rather than a cosmetic rename: the install-gate corroboration classified on the literal `inside*`, and because the gate runs against whatever binary is installed — possibly one older than these skills — the honest fix was a dual-era parse rather than a retarget, which is strictly more correct than what was there. On contract 3 the answer is a determination, not a deletion: the distribution channels are clean but the unpruned plugin cache is a demonstrable blocking population, so the token stays and the retirement query is scheduled for the 0.28 cut. Everything was spiked in a slug-named worktree, all six affected packages pass, and three mutations were exercised — most importantly, restoring the single-era glob turns the new-shape corroboration arm red, which is what proves AC-3 is not a tautology.
