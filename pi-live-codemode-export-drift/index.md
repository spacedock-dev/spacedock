---
id: cvrwyn3xynbwg4wkgc5khyfq
title: pi-live fails because pi-codemode no longer exports MCP_TYPESCRIPT_PREAMBLE
status: validation
source: "Captain directive, 2026-11-10: the pi-live lane is red on main; pin-refresh already merged, so the repair is a new task."
started: 2026-11-10T18:37:11Z
completed:
verdict:
score: 0.9
worktree: .worktrees/spacedock-ensign-pi-live-codemode-export-drift
issue:
pr: "#874"
mod-block:
gates:
    version: 1
    records:
        - id: gate:cvrwyn3xynbwg4wkgc5khyfq:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:cvrwyn3xynbwg4wkgc5khyfq-backlog-1
              briefing:
                id: briefing:cvrwyn3xynbwg4wkgc5khyfq:backlog:attempt-1:revision-1
                digest: sha256:2544a2af2281e604befe554972fd7e5132abba4caacdc03bb0991794ec9bc943
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:cvrwyn3xynbwg4wkgc5khyfq:backlog:1
                briefing: briefing:cvrwyn3xynbwg4wkgc5khyfq:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-10T18:37:03.84826Z"
                decision: approve
                reason: Seed names the diagnosed cause with local evidence, the unproven journeys, and a bounded pin-bump-plus-drift-check outcome; ready for ideation.
                conn:
                    quote: dispatch pi drift one too and approve toward PR stack on top of the haiku one.
                    source: captain chat, 2026-11-10
              application:
                target-stage: ideation
                state: consumed
        - id: gate:cvrwyn3xynbwg4wkgc5khyfq:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:cvrwyn3xynbwg4wkgc5khyfq-ideation-1
              briefing:
                id: briefing:cvrwyn3xynbwg4wkgc5khyfq:ideation:attempt-1:revision-1
                digest: sha256:4121c52326b6aa8912a3ceabaf9299121f7906575e8cf7e6c8681a7f1f1f5035
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:cvrwyn3xynbwg4wkgc5khyfq:ideation:1
                briefing: briefing:cvrwyn3xynbwg4wkgc5khyfq:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-10T18:47:57.648811Z"
                decision: approve
                reason: Design pins the locally proven Pi family, guards the transitive float with a deterministic exact-pin check shown to flag the 1.0.2 tree and pass the 1.5.0 tree, measures the fix on a Linux pi-live run against the failing baseline, and keeps unrelated journey failures out of scope with named dispositions.
                conn:
                    quote: dispatch pi drift one too and approve toward PR stack on top of the haiku one.
                    source: captain chat, 2026-11-10
              application:
                target-stage: implementation
                state: consumed
---

The pi-live job of Runtime Live E2E has not passed since 2026-10-07. Worker launches die before a child session exists with a Pi package mismatch, so no commit on `main` can be greened on the Pi lane. This blocks the 0.28.0 release, which must green its exact release SHA.

## Problem

- Runs 42848114247 (`kent/pi-parity` @ `8d7307c99`) and 42849777496 (`spacedock-ensign/state-durability` @ `02536b6c3`), both 2026-11-10 and both now merged into `main` @ `5785dbed7`, fail the steps "Run live Pi common journeys" and "Run live Pi front-door smoke".
- Both logs carry, in the `recorded-gate-lifecycle` journey: ``Error: `@earendil-works/pi-codemode/declarations` does not export `MCP_TYPESCRIPT_PREAMBLE`.``
- The same error appears in failing runs 39608837699 (2026-10-21) and 41857395650 (2026-11-04).
- Last green pi-live: `main` @ `077fb3bc1` on 2026-10-06, and a PR branch on 2026-10-07. Every pi-live job since has failed or been cancelled.
- The Pi pins live in `internal/pilive/pilive.go` (`pi-coding-agent@1.0.2`, `pi-subagents@0.75.0`, `pi-intercom@0.16.0`), refreshed by `pi-live-lane-pin-refresh` (PRs #816, #830, #833). `pi-codemode` is not among them.

## Diagnosis

`internal/pilive/pilive.go` pins `@earendil-works/pi-coding-agent@1.0.2`, which depends on `@earendil-works/pi-codemode` `^1.0.2`. The install has no lockfile, so npm resolves `pi-codemode` 1.5.0, whose `./declarations` no longer exports `MCP_TYPESCRIPT_PREAMBLE`. Pi 1.0.2's bundled `dist/extensions/codemode/tool.js` imports that symbol, so the extension fails to load. Pi 1.5.0 does not import it.

Local evidence, 2026-11-10 (macOS, model `openai-codex/gpt-6-luna`):

- Pi 1.0.2 installed the way CI installs it: `TestLivePiFrontDoorSmoke` fails in 102 s with the export error.
- Pi 1.5.0 with the pinned substrates (`pi-subagents` 0.75.0, `pi-intercom` 0.16.0): the smoke passes in 263 s.
- Pi 1.5.0 with `pi-subagents` 0.84.0 and `pi-intercom` 0.18.1: 16 of 18 tests pass, including every journey that failed in CI (`smoke`, `smallest-sufficient-mechanism`, `full-ensign-cycle`, `filing`, `gate-guardrail`). `recorded-gate-lifecycle` failed and `rejection-flow` was still running at the suite deadline; see **Journey dispositions**.
- In that run three `pi` xfail bindings passed: `owned-conflict-owner-handoff`, `default-headless-gate-stop`, and `auto-continue-after-implementation` split-root.

Facts checked during ideation (published tarballs and the scratch installs `/tmp/pi-repro`, Pi 1.0.2 installed as CI installs it, and `/tmp/pi-repro2`, Pi 1.5.0):

- `MCP_TYPESCRIPT_PREAMBLE` is exported by `pi-codemode/declarations` up to 1.1.0 and is absent from 1.2.0 (published 2026-10-12) onward. Pi 1.0.2's `^1.0.2` admits all of them. This task does not attribute pi-live failures dated before 2026-10-12.
- Pi declares its first-party dependencies as caret ranges up to 1.2.0 and as exact versions from 1.3.0 on (1.3.0: `"@earendil-works/pi-codemode": "1.3.0"`). Every dependency of Pi 1.5.0 is an exact version. Its first-party packages also pin one another exactly: `pi-agent-core` pins `pi-ai` 1.5.0, and `pi-ai` pins `pi-schema` 1.5.0.
- In the 1.0.2 tree, every nested `@earendil-works/*` package resolved to 1.5.0, not 1.0.2. Local npm has `min-release-age=7`; without that limit, CI resolves 1.5.2.
- Both substrates declare `@earendil-works/pi-coding-agent` as a peer, so the substrate install also places a second, unpinned Pi copy under `<agentDir>/npm`. That copy does not run. Pi's extension loader aliases `@earendil-works/pi-coding-agent`, `pi-ai`, `pi-tui` and `pi-agent-core` to the host Pi (`dist/core/extensions/loader.js`, `getAliases`), and `pi-subagents` spawns child Pi from `process.argv[1]`.

Latest published versions: Pi 1.5.2 (published 2026-11-10 12:24Z), `pi-subagents` 0.87.0, `pi-intercom` 0.18.1.

## Outcome

The pi-live pin moves to a current Pi family that loads cleanly, and the install can no longer drift through an unpinned transitive dependency without a deterministic check failing. pi-live passes in CI on a tip containing the fix. Any journey that still fails for a different cause is named with its evidence rather than repaired here.

## Proposed approach

### Pins

The pins move to the newest family that a local pi-live run proved:

| Package | Current pin | New pin | sha512 integrity (from `npm pack --json`) |
|---|---|---|---|
| `@earendil-works/pi-coding-agent` | 1.0.2 | **1.5.0** | `sha512-xNDZ5kqowO6Z03FS2gwAkPEGk1ssxVcYKQ6qZId84mFEvwACvvyc2l0vAMOS+c5zvxLRMLBE9r4y44IGYkMJaQ==` |
| `pi-subagents` | 0.75.0 | **0.84.0** | `sha512-IagE1hcBxLi82jTxTmcy/pIWckBuqhMH1GfAOe6/aIMVfnopMybCbcTMUS5K0CNy4tqSwHGleaXabaDmJS1QOg==` |
| `pi-intercom` | 0.16.0 | **0.18.1** | `sha512-hY3y+vF6CrkpV0N/55pEsl3DSVjGerlTCQDY2mt37o0jQj8SLW+GmQoKyAEgCyfcxrqJb+Tx4hp4eXqM2hLYgg==` |

This family passed 16 of 18 pi-live tests locally, including every journey that failed in CI. The 1.5.0/0.75.0/0.16.0 combination has only a smoke pass. `pi-intercom` 0.18.1 is also the latest release. The tarballs of the new pins contain every file that `install` and `verify-manifest` require: `bundle/index.js` (`pi.extensions`), `bundle/intercom-bridge.js` (`exports["./intercom-bridge"].default`), `skills/pi-subagents/SKILL.md` and `skills/pi-intercom/SKILL.md`.

The floors stay as they are. `NodeEngineFloor` 22.19.0 matches Pi 1.5.x `engines`. `PiCodingAgentFloor` 1.0.0 is also the launcher's ready-gate floor, so raising it is a product decision and not part of the lane repair. `PiSubagentsFloor` 0.53.0 also stays. The const-block comment changes its basis and date ("the newest family proven by a local pi-live run, 2026-11-10"). By the rule from `pi-live-lane-pin-refresh`, the comment contains no version string.

### Drift guard

After `install` verifies the installed Pi version, it runs one more check. The check reads only manifests in Pi's own install tree:

1. Pi's `package.json` must declare every dependency as an exact version (`x.y.z` or `x.y.z-pre`).
2. Each `@earendil-works/*` package under `<Pi>/node_modules/@earendil-works/` must declare its `@earendil-works/*` dependencies as exact versions. This closes the first-party edges below depth 1, for example `pi-ai` → `pi-schema`. Third-party dependencies of first-party packages are not checked.
3. For each edge that the check covers, `<Pi>/node_modules/<dep>/package.json` must exist and its version must equal the declared one. This reuses `verifyInstalled`.

If any edge fails, `install` exits non-zero. It names the declaring package, the dependency, and the declared range or the installed version, for example `@earendil-works/pi-coding-agent declares @earendil-works/pi-codemode "^1.0.2", not an exact version: the install can float`. Edges are checked in sorted order, so the error is the same on every run.

The check lives in `install` because `install` already owns "the installed package is the pin". It fails in the "Install Pi CLI and substrates" step, before any journey spends model calls, so the workflow does not change. The check would have refused 1.0.2 at the first install of the `pi-live-lane-pin-refresh` pin on 2026-10-04, when the float was possible but had not yet broken anything.

**Spike** (Python equivalent of the three rules, run over the scratch installs; throwaway script at `/tmp/pi-drift-spike-scripts/exactdeps2.py`). `/tmp/pi-repro` (Pi 1.0.2) has 29 edges checked and 6 violations, all Pi 1.0.2 caret declarations, including `@earendil-works/pi-codemode declared '^1.0.2'`. The script exits 1. `/tmp/pi-repro2` (Pi 1.5.0) has 26 edges checked and 0 violations. The script exits 0. Every Pi 1.5.0 dependency is nested under the Pi root in a global install, so rule 3 needs no hoisting lookup. The implementation's first unit test uses the 1.0.2 manifest shape.

### Alternatives considered

- **Pin the latest family (Pi 1.5.2, `pi-subagents` 0.87.0).** CI can install it, but no local run has exercised it. Local `min-release-age=7` blocks it, Pi 1.5.2 is hours old, and it swaps a direct dependency (`grok-mermaid` becomes `lovely-mermaid`). The 0.28.0 release gate would then rest on a family that no run has proven. An ordinary pin refresh can move to 1.5.2 after a local run proves it, and the drift guard holds for that refresh too.
- **Load check: `node` imports Pi's main entry `dist/index.js`,** the same way `guard` already loads `pi-ai/compat`. The spike shows that on the 1.0.2 tree it fails with `The requested module '@earendil-works/pi-codemode/declarations' does not provide an export named 'MCP_TYPESCRIPT_PREAMBLE'`. On the 1.5.0 tree it loads in 0.3 s and writes nothing under an empty `HOME`. Not adopted: it detects drift only after the drift breaks a static import. A compatible drift passes, and the lane silently measures a different tree. The front-door smoke already proves that Pi loads.
- **Lockfile, `npm overrides`, or a separate `pi-codemode` pin.** `npm install -g <tarball>` honours no project lockfile or `overrides`, and a separately installed `pi-codemode` does not replace the copy nested under Pi. Forcing the nested version would need a wrapper `package.json` for the global install, which changes how Pi is installed and found on `PATH`. Pi 1.5.0 already pins its own tree exactly, so a check that it does is enough.

## Stack and proof run

- The implementation branch starts from the implementation branch of task `claude-live-forbid-schedulewakeup`, which by workflow convention is `spacedock-ensign/claude-live-forbid-schedulewakeup`. That task is in ideation and its branch does not exist yet. That branch is itself stacked on `ci/go-126-workflows` (PR #849, which moves the workflows to Go 1.26). The pull request targets that branch, not `main`. If a lower branch moves, the FO rebases this branch onto it.
- Runtime Live E2E's `pull_request` trigger admits only PRs that target `main` (`branches: [main]`), so the stacked PR gets no automatic pi-live run. The proof run is a manual dispatch on the fix branch: `gh workflow run runtime-live-e2e.yml --ref <fix-branch> -f live_cadence=pi`. At `pi` cadence this runs only `offline` and `pi-live`, and it needs CI-E2E-PI approval. Run 42849777496 is a precedent. The run's `headSha` must equal the fix tip.
- A pi-live failure caused by a lower branch's change is named under the same rule as any other different-cause failure (see below). It is not fixed here.

## Journey dispositions

None of these is repaired in this task. The Linux pi-live run on the fix tip is the proof for all of them. If that run fails for any cause other than Pi package loading, the implementation report records the run id, the failing test, its finding codes and messages, and the relevant command or log line. AC-2 is then reported **not met** for that named cause, and the captain decides on a follow-up. This task does not widen to repair it.

1. **`recorded-gate-lifecycle`** has no `pi` binding, so a failure fails the lane. The local failure has a different cause from the one the seed assumed. The Pi FO ran `"${SPACEDOCK_BIN:-spacedock}" dispatch build recorded-gate-task --stamp --entity-path "/var/folders/…/index.md" --workflow-dir "/private/var/folders/…/004"`, with no `--stage` and no `--checklist-file`. The binary exited 2 with `error: flag/file input requires --entity-path, --stage, and --checklist-file` (`internal/dispatch/build.go:134`), from FO session `2026-11-10T17-50-48-758Z_01a1c695…jsonl` under `/tmp/pi-local/real-1.5.0/pi-common/recorded-gate-lifecycle/sessions/`. The `/var` and `/private/var` spellings appear in the same command, but they are not what the binary rejected. `fo-dispatch-core.md` step 5 requires `--stamp --stage {dispatch_stage}`. The cause is how the FO model followed the instruction, it does not depend on the platform, and it can recur on Linux. If the fix-tip run shows the same `successor dispatch build attempts` finding, the report quotes the FO's `dispatch build` command line from the run artifacts, and the follow-up is a new task that owns the Pi FO's `dispatch build --stamp` invocation.
2. **`rejection-flow`** has a `pi` XFAIL bound to `6h3teccccn3qh71yqcmjbjx4` (`own-pi-rejection-worker-topology`). A graded XFAIL or XPASS does not fail the lane. Locally, the 50-minute suite deadline (`go test -timeout 50m`) panicked while this journey had been running for 25m29s. Its own cap is 30 minutes (`SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES`). The risk on CI is the same suite-deadline panic, which fails the step on time and not on a graded result. If that happens, the report quotes the panic's running-tests list and elapsed times and routes them to `own-pi-rejection-worker-topology` and `pi-live-journey-flakiness`. Timeout and parallelism budgets stay out of scope.
3. **The three passing `pi` XFAIL bindings stay unchanged**: `owned-conflict-owner-handoff` (`psvqjf0w8xh2txp9604gsvmz`), `default-headless-gate-stop` (`penfp034pt9s3cgwp7wg3ykk`), and `auto-continue-after-implementation` (`s0gq9p69nztejw8xp3by4k7f`). An XPASS logs `XPASS ALERT` and does not fail the lane (`liveGradeFailsLane`). One local pass does not prove a repair. The owned-conflict binding is documented as flipping between XFAIL and XPASS. The auto-continue binding covers both fixture variants, and in the same run the single-root variant was an XFAIL (`validation-worker-not-dispatched`). The report records the fix-tip run's grade for each of the three as evidence for its owner task.

## Relation to `pi-live-lane-pin-refresh`

That task, now in validation, built the mechanism this task relies on: one pin source in `internal/pilive`, integrity checks before install, `verify-manifest` at both checkpoints, `guard`, and a front-door smoke with both package-root variables unset. It merged through PRs #816, #830 and #833. Its deterministic halves are met. Its lane halves (AC-1/2/3/6/9) wait for a green pi-live run on its 1.0.2 pin, and that run cannot happen because 1.0.2's caret ranges float. That is this task's diagnosis. This task changes only the pin values and adds the drift guard, and it does not revisit that task's mechanisms. The fix-tip run makes the same lane observations at the new pins: installed-version log, `verify-manifest` at both checkpoints, and a smoke with the overrides unset. The FO can cite that run when it closes `pi-live-lane-pin-refresh` and record that the 1.0.2 lane halves were superseded. This task does not change that entity.

## Expected surface and tolerance

Estimate net LOC change: +75, across 3 files (about +85 insertions / -10 deletions).

| File | Change | Insertions / deletions |
|---|---|---|
| `internal/pilive/pilive.go` | 3 versions and 3 integrities, comment basis and date, `verifyExactDependencies`, its call in `install` | about +35 / -8 |
| `internal/pilive/pilive_test.go` | one new install test for the drift guard (cases listed in the test plan) | about +45 / -0 |
| `docs/runtime-live-ci.md` | the doc diff below | +3 / -2 |

Tolerance: net +45 to +115, at most 4 files. The fourth file is allowed only for an install-step comment in `.github/workflows/runtime-live-e2e.yml`. Any change to a workflow step command, timeout, journey binding, floor, or skill text is outside the declared surface.

**Semantics this task may change.** Runtime behavior of `spacedock-release install`: a new refusal when the Pi dependency tree is not exact-pinned, or when an installed copy differs from its declared version. The Pi family that the pi-live lane measures also changes. No command grammar, stored format, authority, launcher behavior or ready-gate floor changes.

**Doc diff** (`docs/runtime-live-ci.md`, Pi local-run paragraph; `spacedock-release` is the developer and CI release tool, so this is the only doc it touches):

```diff
-Run the Pi live proofs locally with the pinned Pi family. The version list lives
-in exactly one place, `internal/pilive`; install the pinned family from those
-constants, then point the harness at the package roots:
+Run the Pi live proofs locally with the pinned Pi family. The version list lives
+in exactly one place, `internal/pilive`; install the pinned family from those
+constants, then point the harness at the package roots. `spacedock-release install`
+refuses a Pi release whose own dependency tree is not pinned to exact versions, so a
+new Pi pin must be such a release:
```

## Acceptance criteria

1. **The Pi lane no longer fails on Pi package loading, measured against the failing baseline.** On the Linux pi-live run dispatched at the fix tip, the install step logs `installed @earendil-works/pi-coding-agent@1.5.0, pi-subagents@0.84.0, pi-intercom@0.18.1`. No job log line contains `does not export` or `does not provide an export named`. `TestLivePiFrontDoorSmoke`, `smallest-sufficient-mechanism`, `full-ensign-cycle`, `filing` and `gate-guardrail` each grade pass. Baseline: runs 42848114247 and 42849777496, where these failed on the export error.
   *Test:* read that run's job log and the `pi-coverage-detail.jsonl` and `pi-front-door-smoke-detail.jsonl` artifacts. This is a live check.
2. **The pi-live job concludes `success` at the fix tip.** If it does not, every failing test is named with its finding and run evidence as **Journey dispositions** describes, none of them is a Pi package-load error, and this criterion is reported not met for those named causes.
   *Test:* the same run's job conclusion and per-test grades. This is a live check.
3. **`spacedock-release install` refuses a floating Pi tree and accepts an exact-pinned one.** It exits non-zero and names the dependency in three cases: Pi or one of its `@earendil-works/*` packages declares a non-exact first-party dependency, Pi declares any non-exact dependency, or an installed copy's version differs from the declared one. It succeeds when every covered edge is exact and installed as declared, and a third-party caret inside a first-party package does not trip it.
   *Test:* the unit test in `internal/pilive/pilive_test.go`, which is deterministic. A local real install into an isolated npm prefix exits non-zero naming `@earendil-works/pi-codemode` at the old 1.0.2 pin, and exits 0 at the new pins. That local install is a one-off check.
4. **The fix tip is clean under the repository checks.** `go test ./...`, `go test ./... -race` and `gofmt -l ./cmd ./internal` (empty output) all pass.
   *Test:* run locally on the fix tip. This is deterministic.

## Test plan

- **Primary proof owner (deterministic):** `internal/pilive/pilive_test.go`. It already owns `install` behavior through `fakeRunner` and fixture manifests (`TestInstallVerifiesAndRegistersThePinnedFamily`). The existing install test writes a Pi manifest with no dependencies, so it passes unchanged. A new test, `TestInstallRefusesAFloatingPiDependencyTree`, builds the global Pi root with nested `node_modules` and asserts four cases. Each case has its own falsifying edit:
  1. Pi declares `"@earendil-works/pi-codemode": "^1.0.2"` while 1.5.0 is installed (the 1.0.2 shape). The error names `pi-codemode` and `^1.0.2`. *Fails if* the exact-version test is removed.
  2. Pi declares `1.5.0` but `1.5.2` is installed. The error names the installed version. *Fails if* the `verifyInstalled` call per edge is dropped.
  3. A nested `pi-ai` declares `"@earendil-works/pi-schema": "^1.5.0"`. The error names `pi-ai` and `pi-schema`. *Fails if* the first-party walk is dropped and only Pi's direct dependencies are checked.
  4. The tree is fully exact, and a nested `pi-ai` declares `"openai": "^6.0.0"`. `install` succeeds. *Fails if* the first-party walk loses its `@earendil-works/` prefix filter.
  Cost: about 45 lines in the existing fixture style, sub-second run time.
- **One-off local real install:** a real `spacedock-release install` with `npm_config_prefix` and `PI_CODING_AGENT_DIR` pointed at fresh temp directories, once at the old pin and once at the new pins. The old pin must fail at the new check, before any journey runs, and the new pins must pass. All new pins are older than the local `min-release-age=7`, so this runs on the developer machine. It confirms the guard on a real npm tree and not only on fixtures.
- **Live (CI):** one `workflow_dispatch` pi-live run at the fix tip (see **Stack and proof run**). It owns AC-1 and AC-2, costs about 35–50 minutes of live model calls, and needs CI-E2E-PI approval. No new live test is added. The existing journeys and the front-door smoke are the instruments.
- **Repository checks:** AC-4.
- **Mechanism justification:** the drift guard serves the Outcome's "cannot drift without a deterministic check failing" and backs AC-1 for every future pin refresh. The simplest alternative, the main-entry load check, is weaker for the reason given under **Alternatives considered**. The pin bump alone serves AC-1 today, but the next pin refresh could pick a caret-ranged Pi again without the guard noticing.

## Out of scope

- The unpinned peer copy of Pi, and `pi-intercom`'s `tsx ^4.20.0`, under `<agentDir>/npm`. Neither runs today (see **Diagnosis**).
- Repairing `recorded-gate-lifecycle`, `rejection-flow` timing, or any XFAIL owner's fault. Changing journey bindings, timeouts or parallelism.
- Raising `PiCodingAgentFloor` or changing the launcher's Pi ready gate.

## Stage Report: ideation

- DONE: The design names the exact Pi family versions to pin and the deterministic check that fails when a transitive Pi dependency can float, with at least one acceptance criterion measured by a pi-live CI run on the fix tip
  Pins are Pi 1.5.0, `pi-subagents` 0.84.0 and `pi-intercom` 0.18.1, with integrities from `npm pack --json`. The guard is an exact-pin check in `install` covering Pi's dependencies and first-party edges. A spike over `/tmp/pi-repro` found 6 violations (exit 1); over `/tmp/pi-repro2` it found 0 (exit 0). AC-1 and AC-2 are measured by a `workflow_dispatch` pi-live run against the failing runs 42848114247 and 42849777496.
- DONE: The two unproven journeys (recorded-gate-lifecycle, rejection-flow) and the three passing pi xfail bindings each have a stated disposition backed by evidence, not repaired in scope
  `recorded-gate-lifecycle`: the local exit 2 was `flag/file input requires --entity-path, --stage, and --checklist-file`. The FO omitted `--stage` and `--checklist-file`, so path spelling is not the cause, and the failure can recur on Linux. `rejection-flow`: it is XFAIL-bound, and the suite-deadline risk is quoted from the panic. The three XPASS bindings stay, because `liveGradeFailsLane(xpass)` is false and auto-continue single-root was still an XFAIL. A rule covers any failure from a different cause.
- DONE: Expected surface (files, net LOC, tolerance) and the stacked-PR base (the claude-live-forbid-schedulewakeup branch) are declared
  Net +75 across 3 files, tolerance +45 to +115 and at most 4 files. The PR is based on `spacedock-ensign/claude-live-forbid-schedulewakeup`, which sits on `ci/go-126-workflows` (#849). The trigger is `branches: [main]`, so the proof run is `gh workflow run … --ref <fix-branch> -f live_cadence=pi`.

### Summary

The fix is a pin bump to the newest locally proven Pi family (1.5.0 / 0.84.0 / 0.18.1) plus one deterministic install-time check. The check refuses a Pi tree whose dependencies, or whose first-party packages' first-party dependencies, are not exact versions installed as declared. Pi has exact-pinned its tree since 1.3.0, so a manifest check is enough without a lockfile. A main-entry `node` import was spiked: it reproduces the exact export error. It was rejected because it misses drift that does not break loading. Ideation found three facts not in the seed. The export was dropped in `pi-codemode` 1.2.0, not 1.1.0. The `recorded-gate-lifecycle` failure was a missing flag in the FO's command, not the macOS `/var` spelling. A stacked PR gets no automatic pi-live run. `pi-live-lane-pin-refresh` owns the mechanism this task reuses; its 1.0.2 lane halves cannot be met, and this task's fix-tip run can be cited to close it.

## Stage Report: implementation

- DONE: internal/pilive pins Pi 1.5.0, pi-subagents 0.84.0, pi-intercom 0.18.1 with verified integrity values, and a real isolated install at the new pins exits 0
  Pins and the date-only comment are in 3399a5928. The three `npm pack --json` integrities match the design table byte for byte. A real `spacedock-release install` at tip a57652a87, run with an isolated `npm_config_prefix`, `PI_CODING_AGENT_DIR` and `TMPDIR` under `/tmp`, printed `installed @earendil-works/pi-coding-agent@1.5.0, pi-subagents@0.84.0, pi-intercom@0.18.1` and exited 0. `guard` then exited 0.
- DONE: spacedock-release install refuses a floating Pi tree naming the dependency, proven by the deterministic unit test and by a one-off real install at the old 1.0.2 pin exiting non-zero on pi-codemode
  `TestInstallRefusesAFloatingPiDependencyTree` runs `install` four times. It refuses the 1.0.2 caret on `pi-codemode`; it fails if the exact-version test is removed. It refuses an installed 1.5.2 against a declared 1.5.0; it fails if the per-edge `verifyInstalled` is dropped. It refuses a nested `pi-ai` caret on `pi-schema`; it fails if the first-party walk is dropped. It accepts an exact tree that has an `openai` caret inside `pi-ai`; it fails if the `@earendil-works/` filter is lost. All four edits were applied and each one failed the test. A real install at the old pin, using a scratch binary with only the pin constants reverted, exited 1 before the substrates were installed (the agent dir stayed empty). It listed 6 violations, including `@earendil-works/pi-coding-agent declares @earendil-works/pi-codemode "^1.0.2", not an exact version: the install can float`.
- DONE: go test ./..., go test ./... -race pass and the diff stays within the declared 3-file, +75 net LOC (+45 to +115, at most 4 files) surface
  Both suites exit 0 at a57652a87, with no ENOSPC lines. The diff against `origin/ci/go-126-workflows` is 3 files, +127/-16, net +111: `pilive.go` +67/-10, `pilive_test.go` +57/-5, `docs/runtime-live-ci.md` +3/-1. That is 48% over the +75 estimate and inside the tolerance. `gofmt -l ./cmd ./internal` lists only `internal/release/runtime_live_evidence_workflow_test.go`, which is unformatted on the base branch and outside this surface, so AC-4's gofmt clause is not clean at the tip for a reason unrelated to this task.

### Summary

The pins now name the Pi 1.5.0 / 0.84.0 / 0.18.1 family. `install` calls `verifyExactDependencies` right after it checks the Pi version. The check covers every dependency Pi declares and every `@earendil-works/*` dependency of the first-party packages nested under Pi. Each covered edge must be an exact version and must be installed under Pi at that version. The guard reports every violating edge in sorted order, joined, not only the first. The design says only that edges are checked in sorted order. The first sorted violation at 1.0.2 is `@earendil-works/chord`, so a first-error-only guard would not name `pi-codemode` as AC-3 requires. To share fixture setup, the existing install test now uses an `installFixture` helper, and its assertions are unchanged. During the run the disk filled (ENOSPC). The FO cleared the Go cache, and the suites were re-run clean. `/tmp/pi-repro*` were gone by then, so the one-off proof used fresh isolated installs, which were deleted afterwards. AC-1 and AC-2 remain for the dispatched pi-cadence CI run on the stack tip.

## Stage Report: validation

- DONE: AC-3 and AC-4 have reproduced evidence: the drift-guard unit test rerun with one falsifying edit, and the guard traced against adjacent variants (exact, caret, tilde, tag, missing installed copy, version mismatch, third-party caret inside a first-party package)
  At `c6fbeb292`, `TestInstallRefusesAFloatingPiDependencyTree` passes. It runs `install` four times and asserts the refusal text for a Pi caret, an installed-version mismatch and a nested first-party caret, and success for an exact tree with an `openai` caret inside `pi-ai`. Five edits, applied one at a time through `go test -overlay` so the worktree stayed unchanged, each failed it: exact-version test removed (case 1), per-edge `verifyInstalled` dropped (case 2), nested first-party walk dropped (case 3), `@earendil-works/` filter dropped (case 4), and the `verifyExactDependencies` call removed from `install` (case 1).
  A scratch matrix of 36 trees called `verifyExactDependencies` directly. Accepted: exact, exact prerelease, third-party caret or tilde inside a first-party package (installed or not), and a Pi with no dependencies. Refused, naming the declaring package and the dependency: `^`, `~`, `latest`, `*`, empty, `5.0.x`, `5.0`, `>=`, `=5.0.0`, `v5.0.0`, leading or trailing space, `||`, hyphen range, `npm:` alias, `github:`, `file:`, `+build`, an embedded newline, a missing installed copy (Pi edge and nested edge), a version mismatch, a wrong package name at the installed path, and an unparsable nested manifest. A tree with several violations rendered identically 20 times, sorted by declaring package and then dependency.
  Real installs, isolated `npm_config_prefix`, npm cache, `PI_CODING_AGENT_DIR` and `TMPDIR` under `/tmp`, deleted afterwards. A scratch binary with only the six pin constants reverted exited 1 with 6 lines, including `@earendil-works/pi-coding-agent declares @earendil-works/pi-codemode "^1.0.2", not an exact version: the install can float`; the agent directory stayed empty and the nested first-party packages had all resolved to 1.5.0. The tip binary printed `installed @earendil-works/pi-coding-agent@1.5.0, pi-subagents@0.84.0, pi-intercom@0.18.1` and exited 0; `guard` exited 0.
  AC-4: `go test ./...` exits 0 at `c6fbeb292` with no failing package. `go test ./internal/pilive/ -race` passes. The full `-race` suite was not rerun at the restacked tip (disk limit, by FO instruction); it passed at the pre-restack tip `a57652a87`. `gofmt -l ./cmd ./internal` lists only `internal/release/runtime_live_evidence_workflow_test.go`, which this task does not touch; the task's own files are formatted.
- SKIPPED: AC-1 and AC-2 are read from the pi-cadence CI run the first officer names: installed versions line, zero export-error lines, per-test grades; any failing test is named with its finding and evidence
  Pending, not graded. Run 42883041637 (`workflow_dispatch`, branch `pi-live-codemode-export-drift`, `headSha` `c6fbeb292`, equal to the fix tip) was checked once at about 19:45Z: `offline` succeeded, `pi-live` (job 130065024727, started 19:15:01Z) was still in progress. Checks that remain: the install step logs the 1.5.0 / 0.84.0 / 0.18.1 line; no job log line contains `does not export` or `does not provide an export named`; `TestLivePiFrontDoorSmoke`, `smallest-sufficient-mechanism`, `full-ensign-cycle`, `filing` and `gate-guardrail` grade pass in `pi-front-door-smoke-detail.jsonl` and `pi-coverage-detail.jsonl`; the `pi-live` job concludes `success`, or each failing test is named with its finding codes and evidence and none is a package-load error.
- DONE: The diff against the haiku branch stays within the approved surface (at most 4 files, +45 to +115 net LOC)
  `git diff --numstat claude-live-forbid-schedulewakeup...HEAD` (`ebf5dccfa`..`c6fbeb292`): `pilive.go` +67/-10, `pilive_test.go` +57/-5, `docs/runtime-live-ci.md` +3/-1. Three files, net +111. No workflow, floor, journey binding, timeout or skill text changed.

### Findings (validator proposals; the first officer owns the disposition)

1. The guard reads only `dependencies`. A caret in Pi's `optionalDependencies` or `peerDependencies` is accepted (matrix rows "pi optionalDependencies caret" and "pi peerDependencies caret"). Proposed class: deferred risk, evidence defect against AC-3's "any non-exact dependency". User and workflow: a maintainer refreshing the Pi pin. Harm: an optional package npm installs from a range could float unreported. Authority: `value-ac[AC-3]` Pi declares any non-exact dependency. Trigger evidence: none today; the installed Pi 1.5.0 and all seven nested first-party manifests declare empty `optionalDependencies` and `peerDependencies`. Promote to material when a pinned Pi release declares either field non-empty.
2. A non-directory entry in `<Pi>/node_modules/@earendil-works/` (for example `.DS_Store`) makes `install` fail with `open …/.DS_Store/package.json: not a directory`. Proposed class: polish. It fails closed, npm does not create such entries, and CI installs into a fresh tree.

### Summary

The deterministic half holds: AC-3 is reproduced by the unit test with five falsifying edits, a 36-row variant matrix, and real installs at the old and new pins; AC-4 passes with the full `-race` run carried from the pre-restack tip and the gofmt clause clean for this task's files. The surface is 3 files, net +111, inside the tolerance. No material finding. AC-1 and AC-2 are not graded because run 42883041637 was still in progress, so the recommendation is withheld until that run is read: PASSED if the five named tests pass with no export-error line and the job concludes `success`; otherwise AC-2 is reported not met for the named causes.
