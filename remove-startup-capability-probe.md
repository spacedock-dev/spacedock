---
title: Remove boot-time gate withdraw capability probe
status: implementation
source: "Captain directive, 2026-08-15: pedantic compat check for unreleased 0.27.x is overkill; incomplete upgrades already surface a proper error at use time"
score: 0.9
group: tooling
id: dav9qnjhsbbg7k1a8x1260h6
gates:
    version: 1
    records:
        - id: gate:dav9qnjhsbbg7k1a8x1260h6:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:dav9qnjhsbbg7k1a8x1260h6-backlog-1
              briefing:
                id: briefing:dav9qnjhsbbg7k1a8x1260h6:backlog:attempt-1:revision-1
                digest: sha256:c5ca9ef1e81bd58a76d435a83a27ae7554485d193e7ecdfdbc432e13c7ccd1f4
                request-digest: sha256:fdc1e0b484e075d56c870855a84cc1484e02db4e2acc1144b3db517253a58d75
                room-ref: ./remove-startup-capability-probe/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:dav9qnjhsbbg7k1a8x1260h6:backlog:1
                briefing: briefing:dav9qnjhsbbg7k1a8x1260h6:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:13.080348Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:dav9qnjhsbbg7k1a8x1260h6:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:dav9qnjhsbbg7k1a8x1260h6-ideation-1
              briefing:
                id: briefing:dav9qnjhsbbg7k1a8x1260h6:ideation:attempt-1:revision-1
                digest: sha256:c74e56309a151da3019d51a56f33bca2bf25e4717e20519982c1d471e59d828d
                request-digest: sha256:2dcee93188db3170f960d3c77fe661952c0bda28b1938a4d3daa5fa003fa9e0d
                room-ref: ./remove-startup-capability-probe/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-15T03:50:34.621223Z"
                reason: 'Entity amended post-prepare (966a9d384): AC-4 reframed as differential against the HEAD baseline; evidence re-run in slug-named worktrees; re-preparing against current bytes'
            - id: gate-attempt:dav9qnjhsbbg7k1a8x1260h6-ideation-2
              briefing:
                id: briefing:dav9qnjhsbbg7k1a8x1260h6:ideation:attempt-2:revision-1
                digest: sha256:80aef53eafe46b7a615a57383c007d77f87b28744fd134bfcbd1606436c9a638
                request-digest: sha256:a3979025d6f14aa622803ea2d436f8f9ef6c30bdba45cc368e33563250118384
                room-ref: ./remove-startup-capability-probe/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:dav9qnjhsbbg7k1a8x1260h6:ideation:2
                briefing: briefing:dav9qnjhsbbg7k1a8x1260h6:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-15T03:56:16.467825Z"
                decision: approve
                reason: 'Captain ruling 2026-08-15 (approve all except x8): approved into implementation'
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-15T02:55:22Z
worktree: .worktrees/spacedock-ensign-remove-startup-capability-probe
---

Remove the same-minor capability probe added in commit b331baf4f ("fix: reject stale same-minor launchers", 2026-08-09). The binary version gate should rely solely on the minor version match; if a stale build lacks `gate withdraw --reason`, the CLI already returns a clear "unknown subcommand" error at the point of use.

## Problem

The probe pays a cost on every first-officer boot to prevent a failure that is
already loud, already actionable, and can happen at most once per stale install.

The cost is paid every boot: one extra `gate --help` subprocess, plus two bullets
in `first-officer-shared-core.md` — the boot-resident core, the most
context-expensive file in the skill set — plus a shell-mirror probe block, a
93-line Go fixture driver, and a troubleshooting paragraph on the published
install page.

The value it buys is close to zero:

1. The hazard window is the unreleased 0.27 dev line only. A stale *same-minor*
   0.27.x launcher predating `gate withdraw` can exist only before 0.27 ships.
2. When the hazard does occur, the binary already fails clearly at the point of
   use. Verified by exercising `spacedock@next/0.27.0-pre4`:
   `spacedock gate withdrawx foo --reason bar` → exit 2, stderr
   `spacedock gate: unknown subcommand (want: prepare|withdraw|record|validate|consume)`.
   A build lacking `withdraw` hits that same guard with `withdraw` absent from
   the allowed set.

So the value chain is broken at the consumer: the captain the probe protects is
already served by the at-use error, one gate later, with a better message.

## Scope, confirmed against HEAD ef8f55c83

First confirmed at 4d1912a69 and re-confirmed at ef8f55c83 after the repo
advanced mid-stage. The two intervening commits are workflow-doc changes;
`git diff 4d1912a69..HEAD` over the four sites is empty, so they are
byte-identical and the confirmation below carries forward unchanged.

All four sites from the seed exist at HEAD and are the only tracked occurrences.
`git grep -nE "REQUIRED_CAPABILITY|Compatible minor|withdrawCapability|Missing capability"`
outside the state checkout returns exactly nine lines, all inside these files.
Three corrections to the seed's site descriptions:

1. `skills/first-officer/references/first-officer-shared-core.md` — remove the
   "Compatible minor" probe bullet and the "Missing capability" abort bullet
   (lines 12-13). Unchanged from the seed.
2. `skills/integration/testdata/version_gate_flow.sh` — remove
   `REQUIRED_CAPABILITY` and the probe block. The seed says "restore the simple
   `exit 0` path"; that is right, and it is a wider revert than it sounds.
   b331baf4f also replaced `exit 0` with
   `"$launcher" status --boot --identify --json; exit $?`. That boot call exists
   only to give the probe an ordering to be asserted against ("ProbesThen**Boots**"),
   and its sole test is one of the two being deleted, so it reverts with the probe.
3. `skills/integration/version_gate_fixture_test.go` — the seed's list
   understates the deletion. `writeGateLauncher` becomes entirely dead, not just
   its `gate --help` branch: its only two callers are the two tests being
   deleted. `captiveInstall`'s `status --boot --identify --json` branch also dies
   with the shell revert in site 2. Full removal: the `withdrawCapability` const,
   `writeGateLauncher`, `invocationLog`, the `INVOCATION_LOG` plumbing in
   `captiveInstall` and `runGateFlow`, the two dead `captiveInstall` branches, and
   both tests. The second test is named
   `TestGateFlowCompatibleSameMinorProbesThenBootsOnce` at HEAD, not
   `...ProbesThenBoots`.
4. `docs/site/get-started/install.md` — remove the "If startup says the installed
   launcher is missing a required command" paragraph. Unchanged from the seed.

## Coordination with remove-redundant-lint-mirrors (zvk9)

zvk9 deletes `skills/integration/testdata/version_gate_flow.sh` and
`skills/integration/version_gate_fixture_test.go` outright. Those are exactly the
files this entity's sites 2 and 3 edit, so **zvk9 supersedes sites 2-3**. Sites 1
and 4 are untouched by zvk9 and are this entity's durable contribution.

Delivery is order-independent, and both orders were exercised in the spike:

- **zvk9 lands first (expected).** This entity ships sites 1 and 4 only, and the
  stage report records sites 2-3 as already done, naming zvk9's commit.
- **This entity lands first.** It ships all four sites; zvk9 then deletes the two
  files it edited.

zvk9's whole-file delete is self-contained, which was checked rather than
assumed: every symbol defined in `version_gate_fixture_test.go`
(`curlInstallToken`, `gateFixtureDir`, `writeExe`, `captiveInstall`,
`runGateFlow`, `installRunCount`, `sentinelPath`, `fixtureSessionID`) is used
only inside that file, so no sibling in `skills/integration` loses a symbol when
it goes. `repoRoot()` lives in `plugin_manifest_test.go` and is untouched by
either entity. A third entity, retire-requires-contract-sentinel, makes
comment-only edits to `plugin_manifest_test.go` and `marketplace_manifest_test.go`
in this same package — disjoint files, no symbols added or removed, so no merge
order conflicts.

A finding that supports zvk9's premise fell out of the spike: with sites 1 and 4
applied and sites 2-3 left alone, `go test ./skills/integration/` stayed green.
The harness that documents itself as a "deterministic shell mirror of the FO
version gate's Startup step 1 flow" does not notice when the prose it mirrors
loses two of its classes. Nothing binds mirror to prose.

## Keep-boundaries

Neighbours that match a probe-shaped search and must survive:

1. `first-officer-shared-core.md` lines 32 and 46 — "capability read" and "gate
   capability probe" name the `spacedock:fo-gate-lifecycle` deferred-load
   trigger. That is a different, per-gate probe. Keep.
2. `internal/ensigncycle/recorded_gate_lifecycle_test.go:281` — asserts
   `gate --help` is *optional* in the recorded-gate lifecycle log. Different
   probe; the assertion also proves the two are not coupled. Keep.
3. `skills/fo-gate-lifecycle/SKILL.md:34` — the `gate withdraw` command itself.
   The command stays; only the boot-time probe *for* it goes.
4. The binary gate's other classes — sandbox check, `--version` line-1 parse,
   binary-absent, wrong-version, minor 0.27. Unchanged (AC-3).
5. `curlInstallToken` in the fixture test — still used by two surviving tests.
6. install.md's `=== "Binary (macOS / Linux)"` and `=== "macOS (Homebrew)"` tabs
   — `TestInstallHintNoDrift` extracts its token-equality source from those tab
   sections. Only the `## Troubleshooting` paragraph goes.
7. `brew upgrade spacedock` in `internal/contract` — a different remedy surface
   (too-old binary), explicitly excluded from the drift lint. Keep.

## Doc diff

`docs/site/get-started/install.md`, under `## Troubleshooting`:

```diff
 Run `spacedock doctor`.

-If startup says the installed launcher is missing a required command, upgrade
-Spacedock and relaunch. On macOS, run `brew upgrade spacedock`. On Linux, rerun
-the checksum-verified binary installer shown above.
-
 ## Next
```

`## Troubleshooting` keeps `Run `spacedock doctor`.` and nothing is added: with
the probe gone there is no startup message for the paragraph to explain.

## Expected surface and tolerance

- **Primary (zvk9 first):** 2 files, -6 lines, +0. Tolerance ±2 lines, no extra files.
- **Contingent (this first):** 4 files, -117/+1. Tolerance ±10 lines, no extra files.

Measured from the spike worktree, not estimated.

## Semantic changes

- **Runtime behavior:** the FO boot no longer runs `spacedock gate --help`. One
  fewer subprocess per boot. A stale same-minor launcher now boots and fails at
  the first `gate withdraw` with exit 2 and the unknown-subcommand message,
  instead of aborting before boot. This is the intended trade.
- **Documentation:** install.md loses one troubleshooting paragraph.
- **Contingent only:** the shell mirror's gate-pass path stops invoking
  `status --boot --identify --json`, and the same-minor gate-pass branch of the
  mirror loses its only test. That branch was equally untested before b331baf4f,
  so this restores the pre-probe coverage state rather than regressing past it —
  and zvk9 deletes the whole mirror regardless.
- **Command grammar, stored formats, authority:** unchanged. No CLI surface moves.

## Spike determination

**No spike needed.** This is pure deletion over proven mechanisms: markdown
prose, a shell fixture, Go tests, and a docs page, all already exercised by the
existing suite. There is no parser round-trip, runtime handoff, on-disk format,
or unproven flag in the change.

The one non-obvious risk was not the mechanism but the *blast radius*: the two
shared-core bullets sit in a file that seven Go packages read, and
`internal/contractlint` both pins shared-core tokens and cross-checks install.md.
Deleting a bullet that a lint pins would strand the build. That was resolved by
exercise, not by reading — in a throwaway worktree at HEAD 4d1912a69 the deletion
was applied and the tests run.

- Sites 1+4 only: `go test ./internal/contractlint/ ./skills/integration/` → both ok, -6/+0.
- All four sites: `gofmt -l` clean, diff -117/+1, and every package that reads a
  changed file passes — `cmd/spacedock-release`, `internal/claudeteam`,
  `internal/contractlint`, `internal/ensigncycle` (463s run alone),
  `internal/release`, `skills/integration`, plus `internal/cli`'s three
  `TestProseFunction*` tests, which are that package's shared-core readers.
- CI's deterministic live-harness control subset
  (`.github/workflows/runtime-live-e2e.yml:78`) → ok in 0.6s.

Every measurement above was re-run from scratch in worktrees named
`spike-remove-startup-capability-probe` and
`control-remove-startup-capability-probe`. The first pass used bare names
(`probe-spike`, `control-clean`) in the shared session scratchpad, which the FO
flagged as a collision hazard after one contaminated a sibling's evidence. The
re-run reproduced every number exactly, so nothing here rests on a shared path.

Three test failures appeared along the way and all three were shown NOT to be
caused by the change, by re-running them in a clean control worktree at the same
HEAD. `TestCodexResolveManifestAgainstInstalledHost` fails identically at clean
HEAD — the sandbox denies reading `~/.codex/config.toml`. The `internal/cli` and
`internal/ensigncycle` package timeouts are load artifacts: clean-HEAD
`internal/cli` also takes 503s, and `internal/ensigncycle` passes when it is not
competing with other test runs. `internal/ensigncycle`'s two shared-core mentions
are inert string literals inside fixture transcripts; neither reads the file's
contents, so the package cannot be sensitive to this change. See AC-4 for why a
bare `go test ./...` is the wrong gate here.

The tokens the lints pin (`curl -fsSL …install.sh | sh`, `brew tap …`,
`brew install …`, `uname -s`, the sandbox registry rows, the launcher-invariant
sentence) all live in the surviving binary-absent bullet and the step-1 sentence,
and `TestInstallHintNoDrift` reads install.md's tab sections, never the
Troubleshooting paragraph. The throwaway worktrees are deleted; the implementation
reruns these as its own tests.

## Acceptance criteria

- **AC-1 (value): the change is a net deletion that leaves no boot-probe
  reference behind.** Cumulative line delta against `origin/main` is negative,
  and no tracked file outside `docs/dev/.spacedock-state/` matches
  `REQUIRED_CAPABILITY|withdrawCapability|Missing capability|Compatible minor`.
  Verified by: `git diff --shortstat origin/main...HEAD` reports more deletions
  than insertions, and the `git grep -nE` above exits non-zero (no match). Both
  halves can move the wrong way — a rewrite-instead-of-delete makes the delta
  positive, and a missed site keeps the grep matching.
- **AC-2: the first officer boots without invoking `gate --help`.** The Startup
  binary gate reaches `«state.boot»()` from a same-minor launcher with no
  capability probe between the version check and boot.
  Verified by: `internal/contractlint` passes with a check asserting the Startup
  section carries no `gate --help`; in the contingent case, additionally by the
  shell mirror reaching `gate passed` with `gate --help` absent from the
  invocation log.
- **AC-3: the surviving abort classes are unchanged.** The binary gate still
  enforces minor 0.27 and still aborts on binary-absent and wrong-version, with
  the same OS-aware hints and sandbox handling.
  Verified by: `go test ./internal/contractlint/` — `TestVersionGateProseOSAwareHint`,
  `TestVersionGateDeferredTrigger`, `TestVersionGateSandboxRegistry`,
  `TestVersionGateProseLauncherInvariantAmendment`, and `TestInstallHintNoDrift`
  pass unchanged. Deleting a surviving bullet or an install.md tab fails these.
- **AC-4: the change introduces no new test failure.** Every package that reads a
  changed file is green, and any package that fails does so identically at the
  same HEAD without the change.
  Verified by: `go test -count=1` green on `cmd/spacedock-release`,
  `internal/claudeteam`, `internal/contractlint`, `internal/release`,
  `skills/integration`, and `internal/ensigncycle` (run alone, not inside
  `./...`), plus `internal/cli -run TestProseFunction`; and CI's deterministic
  live-harness control subset (`.github/workflows/runtime-live-e2e.yml:78`)
  green. Any other failure must be reproduced in a clean control worktree at the
  same HEAD before it is attributed to this change.

  **Do not write this AC as "`go test ./...` passes" — it is not satisfiable on a
  dev box.** `internal/ensigncycle` is the live-agent harness; its live legs
  self-skip when no host credential is present, which is why CI's `offline` job
  (`runtime-live-e2e.yml:75`) runs a bare `go test ./...` green — it runs without
  secrets. Inside a credentialed Claude session those legs actually execute and
  blow the default 10-minute package timeout. Confirmed independently by the
  retire-requires-contract-sentinel ideation (600.5s timeout at HEAD before any
  edit) and by this task's own runs. `go test ./... -list '.*'` is safe — it
  compiles without running, so a test-inventory diff is unaffected.

## Test plan

Deletion plus the existing suite; the surviving contractlint prose pins are the
regression floor. No new test fixtures.

- `go test ./internal/contractlint/` — the five tests named in AC-3, plus the
  AC-2 no-`gate --help` assertion. Cost: seconds. Already exercised in the spike.
- `go test ./skills/integration/` — contingent case only, and only if zvk9 has
  not yet deleted the package's version-gate files. Cost: ~15s.
- The AC-4 package set plus CI's deterministic control subset — not a bare
  `go test ./...`. Run `internal/ensigncycle` on its own; inside `./...` it both
  executes its live legs and competes for CPU, and times out for reasons that
  have nothing to do with this change.
- The AC-2 contractlint assertion is the only line of new test code. It is
  needed because every other check in the file asserts a token is *present*; none
  would fail if the probe were reintroduced. The simplest alternative — relying
  on AC-1's grep — is insufficient because the grep lives in a stage report, not
  in the suite, so it cannot stop a regression after this entity closes.
- No live workflow test. The claim is boot prose plus deletion, not runtime
  orchestration; the FO boot path is already covered by the contractlint pins.

## Stage Report: ideation

- DONE: Design confirms the four-site removal scope against HEAD and reconciles the overlap with remove-redundant-lint-mirrors, which supersedes scope items 2-3
  All four sites exist at HEAD 4d1912a69 and are the only tracked matches; body records three seed corrections and routes sites 2-3 to zvk9 with both landing orders exercised.
- DONE: Value AC measures negative LOC with no surviving capability-probe reference
  AC-1 pairs a negative `git diff --shortstat origin/main...HEAD` with an empty `git grep -nE 'REQUIRED_CAPABILITY|withdrawCapability|Missing capability|Compatible minor'`; spike measured -117/+1 (all four sites) and -6/+0 (sites 1+4).
- DONE: No-spike determination recorded: pure deletion over proven mechanisms
  Recorded under "Spike determination" with the blast-radius risk resolved by exercise in a throwaway worktree rather than by assertion.

### Summary

The removal scope holds at HEAD, with three corrections to the seed: the second
deleted test is `TestGateFlowCompatibleSameMinorProbesThenBootsOnce`;
`writeGateLauncher` dies whole rather than losing one branch, because its only two
callers are the deleted tests; and restoring `exit 0` in the shell mirror also
reverts the `status --boot --identify --json` call that b331baf4f added to give
the probe an ordering to be asserted against. The overlap with
remove-redundant-lint-mirrors is settled by scope rather than by sequencing: zvk9
deletes both files this entity's sites 2-3 edit, so sites 1 and 4 are the durable
contribution and sites 2-3 are contingent on landing first. Both orders were
exercised green.

The design rests on two things proved by exercising rather than reading. The
premise for removal — that a stale launcher already fails clearly at the point of
use — was confirmed against `spacedock@next/0.27.0-pre4`: an unknown gate
subcommand exits 2 with
`spacedock gate: unknown subcommand (want: prepare|withdraw|record|validate|consume)`.
The blast-radius risk was confirmed harmless by applying the deletion in a
throwaway worktree at HEAD and running every package that reads a changed file.

Tests cited, and what would falsify each: `TestInstallHintNoDrift` asserts the FO
prose's curl and brew tokens equal install.md's tab-section forms — deleting the
surviving binary-absent bullet, or an install.md install tab, fails it, while
deleting the Troubleshooting paragraph does not, because it reads only the tab
sections. `TestVersionGateProseOSAwareHint`, `TestVersionGateDeferredTrigger`,
`TestVersionGateSandboxRegistry`, and
`TestVersionGateProseLauncherInvariantAmendment` assert the surviving Startup
classes keep their OS hints, deferred-read trigger, sandbox registry rows, and
launcher-invariant sentence — removing a class other than the two probe bullets
fails them. `internal/cli`'s three `TestProseFunction*` tests bind shared-core
prose notation to routing; deleting a `«function»` line fails them.

One gap is declared rather than closed: no existing check fails if the probe is
reintroduced, since every prose lint asserts presence. AC-2 therefore adds the
single line of new test code in this task — a contractlint assertion that the
Startup section carries no `gate --help`.

### Addendum (post-report corrections, 2026-08-15)

Two corrections landed after the report above was written. Both change the
entity; neither changes the removal scope or the verdict.

1. **AC-4 was unsatisfiable and is rewritten.** It said "`go test ./...` and
   `go test ./... -race` pass". The retire-requires-contract-sentinel ideation
   reported that a bare `go test ./...` does not finish on this box at HEAD
   *before any edit* — `internal/ensigncycle` hits the 10-minute timeout at
   600.5s because its live legs execute when an ambient host credential is
   present. I verified the explanation rather than taking it: CI's `offline` job
   (`.github/workflows/runtime-live-e2e.yml:75`) does run a bare `go test ./...`,
   and it is green there because that job runs without secrets, so the live legs
   self-skip. AC-4 is now a differential — green on every package that reads a
   changed file, plus CI's deterministic control subset, with any other failure
   reproduced in a clean control worktree before being attributed here. Left as
   written, it would have sent the implementer chasing a failure in a package
   this task does not touch.
2. **All evidence was re-run in slug-named worktrees.** The first pass used
   `probe-spike` and `control-clean` in the shared session scratchpad; the FO
   flagged bare names as a collision hazard after one contaminated a sibling's
   evidence. Re-run in `spike-remove-startup-capability-probe` and
   `control-remove-startup-capability-probe`, every number reproduced exactly:
   -6/+0 for sites 1+4, -117/+1 and gofmt-clean for all four, the same package
   set green (`internal/ensigncycle` 463s alone), and the codex control still
   failing identically at clean HEAD.

HEAD also advanced mid-stage from 4d1912a69 to ef8f55c83. The two new commits are
workflow-doc changes and the four sites are byte-identical across them, so the
scope confirmation carries forward; the body is re-anchored to ef8f55c83.
