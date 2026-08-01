---
title: "Dev-stamp compat: in-tree builds must pass the FO version gate against in-tree skills, for dev usage and CI"
status: done
source: "Captain directive 2026-08-01, in the PR #586 merge decision: file the dev-stamp problem separately from z3, 'tackling actual dev usage and ci'. Same-class bites in one day (2026-07-31): (1) FO boot self-aborted on a stale 0.26.0+dev env binary against an in-tree 0.27 plugin (captain override needed); (2) a gpt-5.6-sol staff reviewer self-blocked on the same prose mid-review; (3) codex-live keep-moving-posture + gate-guardrail on PR #586 failed twice, the live FO declaring the PR candidate's 0.27.0-pre2+dev stamp 'not a non-development 0.27 build' — while the same scenarios stayed green on main and on claude/pi lanes."
issue: spacedock-dev/spacedock#581
id: zexbrjhartgykvhm012f527w
gates:
    version: 1
    current:
        gate: gate:zexbrjhartgykvhm012f527w:ideation
    records:
        - id: gate:zexbrjhartgykvhm012f527w:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zexbrjhartgykvhm012f527w-backlog-1
              briefing:
                id: briefing:zexbrjhartgykvhm012f527w:backlog:attempt-1:revision-1
                digest: sha256:0e8e9d466ebe2768f6c01a89a402a427d03ad1ed03ef5080a4b111e2ac0ec2b7
                digest-domain: canonical-bytes
                request-digest: sha256:cd777a48b4365012f90310e3cda8ae4e29c3d20f2439dc16bd00db91141bd4d0
                room-ref: ./dev-stamp-in-tree-version-gate-compat/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zexbrjhartgykvhm012f527w:backlog:1
                briefing: briefing:zexbrjhartgykvhm012f527w:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-31T16:32:37.288799Z"
                decision: approve
                reason: Under the captain's explicit conn to proceed with the six critical lanes, the bound backlog Briefing demonstrates a repeated live-agent failure, constrains ideation to the smallest sufficient correction, and requires behavioral proof without compatibility machinery.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:zexbrjhartgykvhm012f527w:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:zexbrjhartgykvhm012f527w-ideation-1
              briefing:
                id: briefing:zexbrjhartgykvhm012f527w:ideation:attempt-1:revision-1
                digest: sha256:14a4ffd0a63ed60b4622b5759c0fa191bd13fa72fec3f81698bd6349204bbcc8
                digest-domain: canonical-bytes
                request-digest: sha256:0e78927b5089450943d1f789e4c44d2a1cec5ee38cceeb0b6a67845213600d64
                room-ref: ./dev-stamp-in-tree-version-gate-compat/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zexbrjhartgykvhm012f527w:ideation:1
                briefing: briefing:zexbrjhartgykvhm012f527w:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-01T01:26:00.274289Z"
                decision: approve
                reason: 'Delegated sprint conn: the bound ideation proves the two-line bare-dev distinction is the smallest sufficient correction, preserves wrong-minor refusal, and requires one exact-tip Codex journey without compatibility or standing machinery.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
started: 2026-07-31T16:34:00Z
worktree: .worktrees/spacedock-ensign-dev-stamp-in-tree-version-gate-compat
mod-block:
pr: pr-merge:591
verdict: passed
completed: 2026-08-01T06:44:42Z
archived: 2026-08-01T06:44:42Z
---

## Problem statement

An unstamped in-tree build reports the co-located manifest version plus `+dev`.
For this tree that is `spacedock 0.27.0-pre2+dev`, and the compatibility parser
correctly extracts major.minor `0.27`. PR #586 compressed the first-officer gate
from “the version token carries no major.minor at all (`dev`)” to
“major.minor below/above, or `dev`.” That compression erased the distinction
between the bare integer-era sentinel and a compatible build-metadata suffix.

The failure is in the agent-facing classification, not in version production or
comparison. PR #586's Codex live run reached no workflow work in both
`gate-guardrail` and `keep-moving-posture`; the latter ended:

> Blocked by the mandatory Spacedock version gate: installed version is
> `0.27.0-pre2+dev`, but the workflow requires a non-development 0.27 build.

The same-minor binary was valid. The old, explicit no-major.minor wording had
been green on the corresponding main-branch lanes, while Claude and Pi also
accepted the candidate stamp. Existing focused tests independently prove the
mechanical boundary:

- `TestUnstampedSourceBuildReportsCheckoutVersionPlusDev` builds the real
  binary and observes `<checkout-version>+dev` on line 1.
- `TestParseMajorMinor` accepts `0.24.0-pre1+dev` and rejects bare `dev`.
- `TestCompare` accepts a same-minor dev-suffixed token and rejects bare `dev`.

These tests passed during ideation with no product changes. No spike is needed:
the production mechanism is already exercised, and the live before/after branch
evidence isolates the removed wording distinction. The exact edited wording
still requires one fresh Codex journey before acceptance because agent behavior,
not prose presence, is the value under test.

## Proposed approach

Restore only the missing distinction in Startup step 1. Keep the required-minor
stamp, sandbox branch, absent-binary branch, abort action, and doctor remedy
unchanged. Apply this exact diff:

```diff
-1. **Binary version gate.** ... Then run `${SPACEDOCK_BIN:-spacedock} --version`, parse line 1: `spacedock <version>`. These skills require binary minor 0.27. Classes:
-   - **Wrong version** (major.minor below/above, or `dev`): ABORT with the mismatch; run `${SPACEDOCK_BIN:-spacedock} doctor`.
+1. **Binary version gate.** ... Run `${SPACEDOCK_BIN:-spacedock} --version`, parse line 1: `spacedock <version>`. These skills require binary minor 0.27.
+   - **Wrong version**: major.minor below/above/absent (bare `dev`, not `+dev`). ABORT with the mismatch; run `${SPACEDOCK_BIN:-spacedock} doctor`.
```

The ellipsis above stands only for the byte-identical sandbox preamble on the
same line; implementation changes the displayed `Then`/`Classes:` fragments and
the wrong-version bullet exactly as shown. The two replacements add seven bytes
in total, leaving the boot core at 26,753 bytes under its existing 26,754-byte
component cap.

Observable semantics are deliberately asymmetric:

- Same required major.minor with patch, prerelease, or build suffixes, including
  `0.27.0-pre2+dev`, passes and continues into boot.
- A below/above major.minor still aborts, including stale `0.26.0+dev` against
  0.27 skills. The suffix is not a compatibility bypass.
- A token with no major.minor, specifically bare `dev`, still aborts before
  discovery or boot and retains the doctor remedy.
- Binary-absent/install, sandbox, command grammar, stored formats, authority,
  and binary runtime behavior do not change.

## Alternatives rejected

- Do not make all `+dev` tokens pass. That would admit a stale checkout from the
  wrong minor and weaken the existing minor-coupling contract.
- Do not change dev stamping, add a compatibility layer, or introduce a broader
  version model. The current producer/parser/compare path already yields the
  correct verdict for every case in scope.
- Do not add a committed prose oracle, compatibility guard, new CI lane, or
  standing live regression. Existing mechanical tests cover the independent
  version values, and one live Codex validation covers the changed agent
  behavior. If that proof is insufficient, stop and consult the captain before
  adding machinery.
- Do not add stale-environment detection or rebuild guidance here. A stale
  wrong-minor binary is correctly rejected; improving its remedy is a separate
  usability change and is not necessary to fix compatible in-tree builds.

## Expected surface and semantics

- `skills/first-officer/references/first-officer-shared-core.md`: replace two
  physical lines; 2 insertions, 2 deletions, net +7 bytes.

Tolerance: this one file only, at most 3 inserted and 3 deleted lines, and the
existing 26,754-byte component cap remains unchanged. No Go, shell fixture,
manifest, release tool, runtime harness, workflow, site, or CI file changes.
No command grammar, stored format, authority, or binary runtime semantic changes.
The only changed runtime behavior is a first officer accepting a parseable
same-minor `+dev` token instead of inventing a non-development-build requirement.

No site-documentation diff is needed. `docs/releasing.md` already states that a
proxy/source build can report `X.Y.Z+dev`, gates correctly under minor coupling,
and that the suffix is not a compatibility issue. The exact affected
agent-facing documentation diff is the two-line contract diff above.

## Acceptance criteria

**AC-1 — A compatible in-tree dev build crosses the live first-officer version
gate.** A binary built from the exact candidate reports
`spacedock 0.27.0-pre2+dev`; the real Codex `gate-guardrail` journey reaches the
workflow and creates exactly one prepared gate request instead of stopping on a
version mismatch. Test with one exact-tip live invocation and compare the
prepared-request count (`1`) with PR #586's independent failure baseline (`0`).

**AC-2 — The dev suffix does not weaken minor coupling.** Same-minor
`0.27.0-pre2+dev` is compatible, wrong-minor `0.26.0+dev` is incompatible, and
bare `dev` remains incompatible because it has no major.minor. Test with the
existing real-build, parser, compare, and startup-gate tests; each must retain
its current verdict and discovery/boot boundary.

**AC-3 — The correction remains prose-only and bounded.** The delivered diff is
limited to the two declared lines, stays under the existing shared-core byte
cap, and introduces no compatibility machinery, committed check, standing lane,
or wider version semantics. Test by diff/byte inspection followed by the full
offline and race suites.

## Test plan

1. Before editing, run the focused existing tests:
   `go test ./internal/cli -run
   'Test(UnstampedSourceBuildReportsCheckoutVersionPlusDev|DisplayVersionFallsBackToEmbedOnlyWhenUnstamped)$'
   -count=1` and `go test ./internal/contract -run
   'Test(ParseMajorMinor|Compare|StartupGateAbortsBeforeDiscover)$' -count=1`.
   They go red if the real build returns bare `dev`, a suffix changes parsed
   major.minor, bare `dev` becomes compatible, or a wrong version discovers.
2. Apply only the exact two-line diff. Run the focused tests plus
   `go test ./internal/contractlint -run
   'Test(SharedCoreRemainsBelowPreChangeByteCap|FOInstructionComponentCaps|ProseMinorMatchesVendoredManifestMinor)$'
   -count=1`. The cap and manifest-stamp tests fail on size or required-minor
   drift; they are existing independent checks, not a new prose oracle.
3. Build an exact-tip temporary binary, assert its first version line is
   `spacedock 0.27.0-pre2+dev`, then run exactly one real Codex journey:
   `go test -tags live ./internal/ensigncycle -run
   '^TestLiveCodexSharedScenarios$/^gate-guardrail$' -count=1 -timeout 20m -v`
   with `SPACEDOCK_BIN` set to that binary and `SPACEDOCK_REPO_ROOT` set to the
   candidate checkout. Require one prepared request, no mismatch final message,
   and no workflow mutation beyond the scenario's ordinary gate preparation.
   This is one-off validation evidence, not a new standing check.
4. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and
   `go test ./... -race`. Confirm `git diff --check`, the one-file surface, the
   two-line count, and the 26,754-byte cap.

## Stage Report: ideation

- DONE: Define the minimum observable semantics: parseable in-tree +dev passes; bare dev still fails.
  AC-1/AC-2 distinguish same-minor `+dev`, wrong-minor `+dev`, and bare `dev`, with live and mechanical evidence named for each.
- DONE: Bound the expected surface and reject compatibility layers or standing machinery unless the value proves them necessary.
  The approved baseline is one skill file, two line replacements, net +7 bytes; wider machinery is explicitly rejected pending captain consultation.
- DONE: Specify focused behavioral proof, including the live Codex journey and exact user-facing wording.
  The body records the exact diff and a one-off exact-tip Codex `gate-guardrail` journey whose prepared-request count must move from 0 to 1.

### Summary

Ideation isolates PR #586's lost bare-`dev` qualifier as the root cause and
specifies the smallest cap-preserving correction. Existing tests retain the
mechanical boundary; one live Codex journey proves the changed agent behavior
without adding compatibility code or a standing check. Repository formatting,
`go test ./...`, and `go test ./... -race` are green at the ideation baseline.

## Stage Report: implementation

- DONE: Declare the intended file and line replacement before editing.
  Only `skills/first-officer/references/first-officer-shared-core.md` lines 9 and 11 may change: line 9 removes only `Then ` and trailing ` Classes:`; line 11 becomes `   - **Wrong version**: major.minor below/above/absent (bare \`dev\`, not \`+dev\`). ABORT with the mismatch; run \`${SPACEDOCK_BIN:-spacedock} doctor\`.`
- DONE: Apply only the two declared line replacements in skills/first-officer/references/first-officer-shared-core.md; preserve every other file and semantic boundary.
  Commit `4ec2ebeba` is one file, 2 insertions/2 deletions, net +7 bytes; rebasing onto `origin/main` preserved that exact diff and no other candidate file changed.
- DONE: Prove same-minor 0.27.0-pre2+dev passes while wrong-minor 0.26.0+dev and bare dev still fail, staying within the existing 26,754-byte cap.
  Exact-tip binary/doctor returned compatible for `0.27.0-pre2+dev`; a stamped `0.26.0+dev` returned exit 1 with mismatch; `TestCompare` retained bare-`dev` rejection; core size is 26,753 bytes.
- DONE: Run the focused mechanical tests, one exact-tip live Codex gate-guardrail journey, then gofmt/full/race/diff checks; request Roborev and classify every finding before any further edit.
  Focused CLI/contract/contractlint, `go test ./...`, `go test ./... -race`, gofmt, and diff checks passed; live Codex `gate-guardrail` passed in 131.14s at source head `4ec2ebeba`; Roborev panel jobs 623/624 and synthesis 625 passed.
- DONE: Record Roborev observations and authorized dispositions without mutating the reviewed tip.
  FO authorized DECLINE unchanged for both explicitly non-blocking Polish notes: optional fixture arms are a forbidden standing-check expansion, and one-byte cap headroom causes no current harm; no rerun or candidate edit followed.

### Summary

Implementation ships the approved two-line clarification at `4ec2ebeba`, with
same-minor `+dev` accepted and wrong-minor/bare `dev` rejected. Exact-tip live,
mechanical, full, race, cap, and two-reviewer Roborev proof are green; both Polish
review notes were declined unchanged under FO authorization, so the reviewed tip
and one-file boundary remain intact.
