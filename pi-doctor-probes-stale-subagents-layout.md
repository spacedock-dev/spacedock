---
title: Pi doctor probes a pi-subagents source layout the package no longer ships
status: validation
source: "Captain ran `spacedock doctor --host pi` against pi-subagents 0.74.0, 2026-10-01, and reported the two MISSING lines."
score:
started: 2026-10-03T02:42:26Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-doctor-probes-stale-subagents-layout
issue:
id: mc0ajnpb4wh5nhd2p4q0vx9n
gates:
    version: 1
    records:
        - id: gate:mc0ajnpb4wh5nhd2p4q0vx9n:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mc0ajnpb4wh5nhd2p4q0vx9n-backlog-1
              briefing:
                id: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:backlog:attempt-1:revision-1
                digest: sha256:1238b011f4548835cdc914d79e2b7ec094be63f774d7ed21e6cea30f2c69ce8f
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:mc0ajnpb4wh5nhd2p4q0vx9n:backlog:1
                briefing: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-10-03T02:42:14.267902Z"
                decision: approve
                reason: 'Captain directed dispatch of the seed: the defect, scope, and the ACs pinning the probe''s proof owner are stated.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:mc0ajnpb4wh5nhd2p4q0vx9n:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:mc0ajnpb4wh5nhd2p4q0vx9n-ideation-1
              briefing:
                id: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:ideation:attempt-1:revision-1
                digest: sha256:c8cb45aa0f2ab5d611cf05fcf867023d63b56105599a2cd42f42d2e901a685d3
                room-ref: '@review/ideation/briefing-1'
              withdrawal:
                by: agent:first-officer
                at: "2026-10-03T02:58:08.665278Z"
                reason: Bound artifact's AC-1 identifier wraps its bold span across two lines (lines 176-177), so the shipped status --read --ac-scan reader omits AC-1 entirely; proven by unwrapping it in a copy, where the scan returns AC-1, AC-2, AC-3. The gate AC cross-check cannot see the primary value AC. Repairing before presentation.
            - id: gate-attempt:mc0ajnpb4wh5nhd2p4q0vx9n-ideation-2
              briefing:
                id: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:ideation:attempt-2:revision-1
                digest: sha256:b77d1eb699105693bfbc421dfd2d921e1d5cddce9f28f1e4a1772e48e4d42c9c
                room-ref: '@review/ideation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:mc0ajnpb4wh5nhd2p4q0vx9n:ideation:2
                briefing: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-10-03T05:27:26.731446Z"
                decision: approve
                reason: Covers resolving both doctor probes from the installed package declarations, including the folded remedy acceptance; the CI run proves every Pi launch is refused until this lands.
                conn:
                    quote: yes do it
                    source: captain instruction, this session, 2026-10-03
              application:
                target-stage: implementation
                state: consumed
        - id: gate:mc0ajnpb4wh5nhd2p4q0vx9n:validation
          stage: validation
          attempts:
            - id: gate-attempt:mc0ajnpb4wh5nhd2p4q0vx9n-validation-1
              briefing:
                id: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:validation:attempt-1:revision-1
                digest: sha256:b1924916ac17b43cfe00e60da44784f3cc056e7d80616f06d06632230b03382e
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:mc0ajnpb4wh5nhd2p4q0vx9n:validation:1
                briefing: briefing:mc0ajnpb4wh5nhd2p4q0vx9n:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-10-03T15:43:18.139111Z"
                decision: approve
                reason: Captain approved the independent validation result in this session.
              application:
                target-stage: done
                state: pending
sprint: pi-ux
group: tooling
sprint-readiness: ready
pr: pr-merge:817
---

## Problem

`spacedock doctor --host pi` reports two prerequisites missing against a
correctly installed, up-to-date `pi-subagents`:

```
MISSING pi-subagents extension: .../pi-subagents/src/extension/index.ts
MISSING pi-subagents intercom bridge: .../pi-subagents/src/intercom/intercom-bridge.ts
```

Both probes test for TypeScript sources. The package ships compiled output.
The problem statement recorded pi-subagents 0.74.0 declaring
`"pi": { "extensions": ["./index.js"] }`; the installed 0.75.0 has the same
layout and pins the exact manifest resolution this task uses. Inspected on
2026-10-03 at `~/.pi/agent/npm/node_modules/pi-subagents` (version 0.75.0):

- `pi.extensions` = `["./index.js"]`; the root `index.js` exists.
- `src/extension/index.js` and `src/intercom/intercom-bridge.js` exist.
- No plain `.ts` source exists under `src/` (the only `.ts` files are `.d.ts`
  declarations); neither `src/extension/index.ts` nor
  `src/intercom/intercom-bridge.ts` exists.
- The intercom bridge is **not** declared in `pi.extensions`. It is reachable
  through the package `exports` map: `exports["./intercom-bridge"].default` =
  `./src/api/intercom-bridge.js`, which re-exports
  `./src/intercom/intercom-bridge.js`.

The printed remedy cannot clear the line it appears under: `pi install
npm:pi-subagents` reports `up to date`.

Four call sites hardcode the source paths in `internal/cli/pi.go`:
`piRuntimeConfigFromEnv` (`:617`), `checkPiRuntime` (`:644`, `:646`), and
`printPiDoctorReport` (`:773`, `:777`). The comment at `:303` repeats the
assumption ("package discovery loads `<pkg>/index.ts`").

## Risk evidence

The extension is live, not merely present: a session with pi-subagents loaded
carries the `subagent`, `subagent_supervisor`, and `bg_wait` tools. The defect
is in the probe, so the design question is narrow — how the probe should resolve
a package's entry instead of assuming one.

Independent baseline that can move the wrong way: with the current binary both
lines print `MISSING` on a machine where the package is correctly installed and
loadable. The fix must move those two lines to `OK` against that baseline.

## Approach

Resolve each probe target from the package's own `package.json` — the same
artifact pi reads — instead of hardcoding a source path:

- **Extension**: the first entry of `pi.extensions`, resolved relative to the
  package root. This is what pi's own package discovery loads.
- **Intercom bridge**: the package's `exports["./intercom-bridge"].default`
  subpath, resolved relative to the package root. This is the package's declared
  bridge entry.

One resolver, `resolvePiSubagentsEntries(packageRoot)`, reads
`<packageRoot>/package.json` once and returns both absolute paths (empty when
the manifest is missing or unparseable, so the probes fail closed).
`checkPiRuntime` stats the resolved paths; `piCheckResult` carries them so
`printPiDoctorReport` prints the resolved path; and `runPi`'s
unregistered-package fallback passes the resolved extension path to
`--extension` (it previously passed the nonexistent `.ts` path).

**Design decision (A), approved by the first officer and recorded here.** The
seed's AC-1/AC-2 originally implied both probes derive from one `pi.extensions`
entry. The real package contradicts that: pi-subagents declares only
`./index.js` in `pi.extensions` and does not declare the bridge there at all.
The bridge is therefore resolved from `exports["./intercom-bridge"]`. The
rejected alternative — a source-extension fallback at
`src/intercom/intercom-bridge.{ts,js,mjs}` — was refused because it
reintroduces the exact layout assumption this task removes and fails AC-3 (the
expected path would come from the probe's own guess, not the artifact).

Mechanism justification:

- The resolver is the end-value mechanism for AC-1/AC-2, not an enabling
  mechanism proved for its own internals.
- Simplest alternative considered: keep the hardcoded path and add a `.js`
  sibling candidate. Insufficient — it still assumes `src/extension/` and
  `src/intercom/`, so it re-breaks on the next package reorganization, which is
  precisely the defect.
- The `exports["./intercom-bridge"]` lookup is the minimal manifest-driven way
  to name the bridge; the package publishes no other declaration for it.

**No spike needed.** The only mechanism is standard `encoding/json` decoding of
the package's `package.json` plus `ops.Stat` on the resolved paths — both
already proven in this file (`readPackagePiSkills`, `isWorkspaceRoot`) and in
the live 0.75.0 artifact inspected above. No runtime handoff or on-disk format
is unverified.

## Out of scope

The doctor's other pi checks, the `--version` floor, the duplicate-registration
warning, and the CI lane's own `src/...index.ts` presence assertions
(`internal/release/workflow_exec_guard_test.go`, `.github/workflows/*`) — a
separate surface. This task changes only how the two pi-subagents probes resolve
their target paths.

## Expected surface and tolerance

Code surface: **2 files, net +60 to +90 LOC** (tolerance ±25%), reported as
insertions **~+115** and deletions **~-40** (the deletions cover the retired
hardcoded-path lines and the reshaped stat maps).

- `internal/cli/pi.go` — the `resolvePiSubagentsEntries` helper plus wiring at
  the four call sites, and the `:303` comment.
- `internal/cli/pi_frontdoor_test.go` — manifest-driven fixture plus reshaped
  stat maps.

Documentation surface (proposed diff below): **2 files, net 0** (four one-line
path edits).

Observable semantics this task **may** change:

- The two doctor lines' verdict for a correctly installed package (`MISSING` →
  `OK`) and the path each line prints — the resolved manifest entry
  (e.g. `.../pi-subagents/index.js` and
  `.../pi-subagents/src/api/intercom-bridge.js`).
- The `--extension` target `runPi` passes in its unregistered-package fallback
  (now the resolved `pi.extensions` entry).

Observable semantics this task **must not** change:

- The report's line labels (`pi-subagents extension`, `pi-subagents intercom
  bridge`), the `OK`/`MISSING` grammar, remedy text, exit codes, and every other
  doctor line. **Proposed M4 scope exception, pending captain approval:** only
  the bridge remedy text may change, exactly as proposed below; all other
  remedy text remains unchanged.
- `pi.skills` discovery, the `--version` floor, the Spacedock package checks,
  and the front-door launch shapes.

## Acceptance criteria

**AC-1 - The two pi-subagents probes report OK for a package that declares a compiled entry.**
Verified by: a Go test that writes a fixture package from a declarative
manifest — a `package.json` declaring `pi.extensions` and
`exports["./intercom-bridge"]`, plus those declared entry files and no plain
`.ts` anywhere — then asserts both the extension and the intercom bridge lines
print `OK` and name the manifest-resolved paths. Independent baseline that can
move the wrong way: on today's binary the same fixture prints `MISSING` for both
lines.
Falsifying edit: restore the hardcoded `src/extension/index.ts` probe path —
the compiled-entry fixture must turn the test RED.

**AC-2 - Each probe reports MISSING when its own declared target is absent.**
Restated from the seed's single-file phrasing under approved decision (A):
pi-subagents declares the extension in `pi.extensions` but the bridge in
`exports`, so one shared entry cannot drive both lines. The restatement is
reported at the ideation gate and is proposed, not pre-approved.
Verified by: the same manifest-driven fixture with the declared `pi.extensions`
entry file removed → the extension line prints `MISSING` (the bridge line stays
`OK`); and with the `exports["./intercom-bridge"]` target removed (or its
`exports` entry deleted) → the bridge line prints `MISSING`. Each absence turns
its own line `MISSING`.
Falsifying edit: make a probe accept any existing file under the package root —
removing the declared entry must turn the test RED.

**AC-3 - The probe expectation does not come from the probe implementation.**
Verified by: the fixture reads back the `package.json` it wrote and derives the
expected paths from that manifest, rather than naming the paths the production
code checks in a `statOK` map.
Falsifying edit: delete the manifest write from the fixture (or the written
`package.json`) — the test must fail because the expectation cannot be derived
from the artifact.

**AC-4 - Each affected prerequisite's printed remedy clears that same failing doctor line in a disposable home.**
**Proposed M4 acceptance addition — pending the captain's gate, not approved.**
Verified by: an integrated, one-off validation using the implementation candidate
and a real installed pi-subagents package in a disposable home. Independently
remove the manifest-declared extension target, then the manifest-declared bridge
target in a restored case. For each, capture `doctor --host pi` reporting that
line `MISSING`, execute the command printed beneath it without substituting an
unprinted repair, and rerun the same binary/environment: that same line must be
`OK` at its manifest-declared path, while the other probe remains `OK`.
Retain package/tool versions, exact commands and exits, and full before/after
doctor stdout/stderr in the entity's validation report (or a committed artifact
linked there). The primary proof owner is the implementation ensign, with the
validation ensign independently checking the retained transaction; the Go tests
for AC-1–3 cannot discharge AC-4. No global install mutation is authorized.
Falsifying edit: print a successful no-op command as the remedy — its exit 0
must not pass while the same prerequisite remains `MISSING`. Reintroducing a
stale `.ts` probe also fails the post-remedy assertion even when installation
succeeds. If the printed command does not clear its own line, acceptance fails:
return the needed text or scope change to the captain, not an install-exit-only
or file-existence substitute for the outcome.

## Test plan

Primary proof owner: `internal/cli/pi_frontdoor_test.go`. Its current
`statOKForPiResources` and `writePiSubagentsFixtures` fabricate
`src/extension/index.ts` and `src/intercom/intercom-bridge.ts` — the exact paths
the production code stats — so the tests pass by construction and cannot fail
when a real package changes layout. Reshape:

- `writePiSubagentsManifest(t, pkg, manifest)` writes `<pkg>/package.json` from
  a small struct (`pi.extensions`, `exports`) and creates the declared files.
- `readPiSubagentsManifest(t, pkg)` reads the written `package.json` back and
  returns the resolved expected paths — the expectation derives from the
  artifact, which is AC-3.
- `fakePiRuntimeOps.Stat` keys become the manifest-resolved paths
  (`<pkg>/index.js`, `<pkg>/src/api/intercom-bridge.js`), not the `.ts` paths.

**Proposed M4 proof-scope change, pending captain approval:** deterministic
Go tests remain the proof for path resolution (AC-1–3), supplemented by the
one-off, disposable-home remedy transaction for AC-4 described below. No live
agent dispatch/talkback lane is added; this exercises the package installer and
the doctor command, not a host session.

Per-check falsifying edits (each turns its own check RED):

- AC-1 extension: hardcode `src/extension/index.ts` again.
- AC-1 bridge: drop the `exports` lookup, restoring
  `src/intercom/intercom-bridge.ts`.
- AC-2: accept any file under the package root.
- AC-3: remove the manifest write, or stop deriving expectations from
  `package.json`.

Cost: low — one new resolver plus fixture helpers, no runtime changes.

## M4 remedy acceptance and scope proposal

This is a staff-review fold over commits `2fac0d774` and `38a5f040d`, not a new
design. AC-1–3 and the resolver approach are unchanged. AC-4, the integrated
proof scope, and the bridge-only wording exception are **proposed for the
captain's gate**, not authorized implementation changes. The earlier "No spike
needed" applies to manifest decoding/stat resolution only; the newly proposed
installer transaction was exercised separately below.

### Remedy-text resolution (proposed)

The existing extension command `pi install npm:pi-subagents` cannot fix the old
binary's false negative: installation does not manufacture the obsolete `.ts`
paths. Resolving the two probes from the manifest remains necessary. Against
Pi 1.0.0 / npm 11.8.0 / pi-subagents 0.75.0, the same install command **did**
restore each independently deleted declared target in the isolated spike. Thus
no extension-remedy change or forced uninstall/reinstall is justified by this
evidence. The seed's "up to date" observation does not prove that a genuinely
missing declared file will remain absent; nor does this spike prove that the
candidate doctor line clears (there is no implementation candidate yet).

The bridge's current "install/update" phrase does not print an executable
command. Proposed exact change in `printPiDoctorReport`, bridge line only:

```diff
-install/update pi-subagents or set PI_SUBAGENTS_PACKAGE_ROOT to a package root containing the intercom bridge
+run `pi install npm:pi-subagents` or set PI_SUBAGENTS_PACKAGE_ROOT to a package root containing the intercom bridge
```

This is the needed narrow exception to the previous remedy-text freeze and the
out-of-scope statement that only target resolution changes. It does not change
other remedies, exit codes, override resolution, labels, or line grammar. The
extension remedy stays verbatim. The existing four-line proposed documentation
diff below stays unchanged; it does not quote the bridge remedy.

The additional deterministic assertion belongs to the existing primary proof
owner, `internal/cli/pi_frontdoor_test.go`: a missing bridge prints the exact
executable remedy above; reverting to vague "install/update" makes it fail.
That assertion proves wording only, not AC-4. The existing two-code-file estimate
and ±25% tolerance remain the proposed budget: the one-line replacement and
approximately ten test lines fit within it. Live proof is recorded in the entity,
not a new CI lane or test framework. Simplest alternative: leave the bridge text
unchanged and have the validator invent a command. Rejected because it does not
exercise a printed command and could conceal a non-actionable remedy. A forced
remove/reinstall was considered and exercised but adds disruption with no
observed benefit over install-only, so is not proposed.

### Integrated proof procedure (proposed; AC-4 owner assignment)

1. Implementation ensign builds the candidate once and records its revision,
   binary version/hash, Pi, Node, npm and pi-subagents versions. Use a fresh
   `mktemp -d` root; keep HOME, `PI_CODING_AGENT_DIR=$HOME/.pi/agent`, npm cache,
   npm prefix, working directory and all install/settings files inside it. Use
   `env -i` with only the required PATH and those isolation variables; inherit
   no `PI_SUBAGENTS_PACKAGE_ROOT`, auth, project `.pi`, or global npm settings.
   The existing Pi/Node executables may be read, never updated. Install the real
   `npm:pi-subagents` package there and retain its manifest/version. Registry
   access is required; an unavailable registry means proof is blocked, not passed.
2. Record the candidate's healthy **two pi-subagents probe lines** before fault
   injection. Derive targets from the installed manifest, not production code.
   Remove only the extension target; record full doctor output and exit, execute
   the exact printed install command, record output and exit, then record full
   doctor output and exit again. Require extension `MISSING` → `OK` at the
   declared path and bridge `OK` throughout. Never create a substitute `.ts` file
   or manually restore the target after the failing observation.
3. With the repaired install, repeat for only the bridge target, requiring bridge
   `MISSING` → `OK` and extension `OK` throughout. Retain manifest/package versions
   before and after each transaction to expose any registry-driven update. Do
   not claim a tested version beyond those recorded.
4. Capture all command exits, including nonzero doctor exits. Unrelated auth,
   pi-intercom or Spacedock checks may keep the overall doctor exit at 1; acceptance
   is the same affected **line** clearing, not a fabricated overall healthy host.
   Keep those unrelated diagnostics in the full output; no credentials or global
   installs are needed. Validation ensign checks versions, exact printed-command
   fidelity, failure setup and line-by-line postconditions in the durable report.
5. If install succeeds but the line stays `MISSING`, retain that failure and return
   the exact needed remedy-text or scope change to the captain. Do not substitute
   the spike's file-restoration evidence, silently run remove/reinstall, or relax
   AC-4. If the captain rejects the bridge text exception, explicitly resolve the
   resulting acceptance gap at the gate; the unchanged freeze is not evidence.

Estimated cost: two small package installs/repairs plus command captures, minutes
with registry access; no live model/session or new persistent mechanism. A
Go-only simulated install is simpler but insufficient: it cannot establish what
the real printed package-manager command does to the declared artifact.

### Isolated remedy spike evidence (ideation, not AC-4 completion)

Executed 2026-10-03 with Spacedock **0.28.0-pre3** (unfixed), Pi **1.0.0**,
Node **v24.13.1**, npm **11.8.0**, pi-subagents **0.75.0** before and after every
repair. Disposable root: `/tmp/spacedock-pi-remedy.R9fC5T`. No global install,
credentials, source code or runtime configuration was changed. The following
setup describes the actual isolation (PATH supplied existing read-only tools):

```sh
ROOT=$(mktemp -d /tmp/spacedock-pi-remedy.XXXXXX)
mkdir -p "$ROOT/home" "$ROOT/work"
cd "$ROOT/work"
# Run every pi/doctor command through this environment:
env -i PATH="$PATH" HOME="$ROOT/home" \
  PI_CODING_AGENT_DIR="$ROOT/home/.pi/agent" \
  npm_config_cache="$ROOT/home/.npm" \
  npm_config_prefix="$ROOT/home/.npm-global" <command>
```

Installed manifest: `pi.extensions[0] = ./index.js`;
`exports["./intercom-bridge"].default = ./src/api/intercom-bridge.js`.
For each case the selected file was removed before the "before" doctor run.
The table retains commands/exits and the observed artifact state (not an inferred
`OK`). `doctor` below means
`/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock doctor --host pi`.

| Case | Command | Exit | Observation |
| --- | --- | --- | --- |
| Setup | `pi install npm:pi-subagents` | 0 | 0.75.0 installed, both targets present |
| Setup | `doctor` | 1 | Both stale `.ts` probes MISSING |
| Extension target removed | `doctor` (before) | 1 | Full output retained below |
| Extension | `pi install npm:pi-subagents` | 0 | `changed 5 packages`; declared extension restored |
| Extension | `doctor` (after install) | 1 | Byte-identical to before; stale probes still MISSING |
| Extension alternative | `pi remove npm:pi-subagents` then `pi install npm:pi-subagents` | 0, 0 | Declared target present; no extra benefit |
| Extension alternative | `doctor` | 1 | Byte-identical output |
| Bridge target removed | `doctor` (before) | 1 | Full output retained below |
| Bridge | `pi install npm:pi-subagents` | 0 | `changed 5 packages`; declared bridge restored |
| Bridge | `doctor` (after install) | 1 | Byte-identical to before; stale probes still MISSING |
| Bridge alternative | `pi remove npm:pi-subagents` then `pi install npm:pi-subagents` | 0, 0 | Declared target present; no extra benefit |
| Bridge alternative | `doctor` | 1 | Byte-identical output |

All seven doctor captures (setup; each case before, after install, after the
alternative) have SHA-256
`7ec79e91046a878505739d9f276e13b322662f715afc8addafb1ec7dbf6edecd`.
Their full combined stdout/stderr is identical and retained once below without
path normalization. This is durable negative evidence: the **old** printed
remedy cannot clear the **old** probe. Candidate before/after `MISSING` → `OK`
proof remains required and assigned, not claimed by this ideation report.

```text
Pi runtime check
OK pi CLI: /Users/clkao/.local/state/fnm_multishells/77762_1790878083686/bin/pi
MISSING Pi auth: /tmp/spacedock-pi-remedy.R9fC5T/home/.pi/agent/auth.json
  remedy: run `pi` login/auth flow; live tests copy this file into an isolated PI_CODING_AGENT_DIR
MISSING pi-subagents extension: /tmp/spacedock-pi-remedy.R9fC5T/home/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
  remedy: run `pi install npm:pi-subagents` or set PI_SUBAGENTS_PACKAGE_ROOT
OK pi-subagents skill: /tmp/spacedock-pi-remedy.R9fC5T/home/.pi/agent/npm/node_modules/pi-subagents/skills/pi-subagents
INFO Pi auth/session dirs: auth=/tmp/spacedock-pi-remedy.R9fC5T/home/.pi/agent/auth.json session=/tmp/spacedock-pi-remedy.R9fC5T/home/.pi/agent/sessions
Supervisor-talkback setup prerequisites
MISSING pi-subagents intercom bridge: /tmp/spacedock-pi-remedy.R9fC5T/home/.pi/agent/npm/node_modules/pi-subagents/src/intercom/intercom-bridge.ts
  remedy: install/update pi-subagents or set PI_SUBAGENTS_PACKAGE_ROOT to a package root containing the intercom bridge
MISSING pi-intercom package root: /tmp/spacedock-pi-remedy.R9fC5T/home/.pi/agent/npm/node_modules/pi-intercom
  remedy: set PI_INTERCOM_PACKAGE_ROOT to the installed pi-intercom package root
MISSING pi-intercom skill: /tmp/spacedock-pi-remedy.R9fC5T/home/.pi/agent/npm/node_modules/pi-intercom/skills/pi-intercom
  remedy: install pi-intercom or set PI_INTERCOM_PACKAGE_ROOT to a package root containing skills/pi-intercom/SKILL.md
MISSING Spacedock package
  remedy: run `spacedock install --host pi` to install the Spacedock package (or `spacedock install --host pi --plugin-dir <checkout>` for a dev override)
MISSING Spacedock extension: .pi/extensions/spacedock.ts
  remedy: run `spacedock install --host pi` to reinstall the Spacedock package (its .pi/extensions/spacedock.ts delivers the FO contract)
MISSING Spacedock first-officer skill (package discovery): skills/first-officer
  remedy: run `spacedock install --host pi` to reinstall the Spacedock package
OK pi version: 1.0.0 (floor 0.83.0)
NOTE: These checks verify necessary supervisor-talkback setup prerequisites only; they are insufficient to prove live child talkback.
NOTE: Live proof still requires the cq-style progress -> decision -> supervisor reply -> child resume -> durable marker probe for pi-intercom-supervisor-talkback.
```

## Proposed documentation diff

The manual `--extension` examples name the stale path this task retires. Same
edit in two docs (four lines):

`docs/runtime-support.md:192`

```diff
---extension ~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
+--extension ~/.pi/agent/npm/node_modules/pi-subagents/index.js
```

`docs/runtime-support.md:247`

```diff
-~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
+~/.pi/agent/npm/node_modules/pi-subagents/index.js
```

`docs/site/contributing/adding-a-runtime.md:113`

```diff
---extension ~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
+--extension ~/.pi/agent/npm/node_modules/pi-subagents/index.js
```

`docs/site/contributing/adding-a-runtime.md:194`

```diff
-~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
+~/.pi/agent/npm/node_modules/pi-subagents/index.js
```

## Stage Report: ideation

- DONE: Task body states the problem, the chosen approach, and criteria under the exact `## Acceptance criteria` heading, with AC-1 through AC-3 each paired to a named falsifying edit that turns it RED.
  Task body has `## Problem`, `## Approach`, and `## Acceptance criteria` (scanner-exact heading); AC-1 falsifier = re-hardcode `src/extension/index.ts`, AC-2 = accept any file under the package root, AC-3 = remove the manifest write.
- DONE: Expected surface declares net LOC change and file count with tolerance, reports insertions and deletions separately, and names the observable semantics the task may change (doctor per-probe OK/MISSING lines) versus those it must not.
  `## Expected surface and tolerance`: 2 code files net +60..+90 (±25%), insertions ~+115 / deletions ~-40; may-change = the two doctor lines' verdict + printed path and the `--extension` fallback target; must-not-change = labels, OK/MISSING grammar, remedies, exit codes, other lines.
- DONE: Test plan names `internal/cli/pi_frontdoor_test.go` as primary proof owner, reshaped to a manifest-driven fixture whose expected paths derive from the written `package.json` rather than a stat map, and records the spike result (or "no spike needed" with the proven mechanisms).
  `## Test plan` names `internal/cli/pi_frontdoor_test.go` and the `writePiSubagentsManifest`/`readPiSubagentsManifest` reshape; `## Approach` records "No spike needed" with the proven `encoding/json` + `ops.Stat` mechanisms.

### Summary

Authored the ideation baseline for the two stale pi-subagents doctor probes.
Inspected the real installed package (pi-subagents 0.75.0) and found the seed's
AC-1/AC-2 implied one shared `pi.extensions` entry, but the bridge is declared
only via the package `exports` map; escalated and got decision (A) approved
(extension from `pi.extensions`, bridge from `exports["./intercom-bridge"]`).
Output: a fleshed task body with the resolver approach, a restated per-probe
AC-2 (reported at the gate as proposed, not pre-approved), and a manifest-driven
test plan.

### Same-stage repair: AC-1 bold-span line wrap

Repaired the one defect from the gate withdrawal: AC-1's bold identifier
wrapped across two lines (176-177 at authoring), so the line-based `--ac-scan`
reader could not see AC-1. Joined the AC-1 bold span onto a single line; no
other paragraph was reflowed.

Evidence — `spacedock status --read docs/dev/.spacedock-state/pi-doctor-probes-stale-subagents-layout.md --ac-scan --json --workflow-dir docs/dev`:

```json
{"command":"read","stage":"ideation","acs":[{"id":"AC-1","line":"180","unevidenced":"false"},{"id":"AC-2","line":"191","unevidenced":"false"},{"id":"AC-3","line":"204","unevidenced":"false"}]}
```

Checklist count summary:

- DONE (1/3): AC-1's bold identifier closes on a single line; `--ac-scan` now
  lists AC-1, AC-2, and AC-3 (previously AC-2 and AC-3 only).
- DONE (2/3): every other byte of the acceptance-criteria section, including the
  `## Acceptance criteria` heading and each AC's 'Verified by' and falsifying
  edit, is unchanged apart from that line wrap.
- DONE (3/3): the repair is committed path-scoped in the state checkout with no
  other file touched and no staged residue.


## Stage Report: ideation (cycle 2)

- DONE: Assign the integrated remedy acceptance: in a disposable home, demonstrate the applicable failing prerequisite, execute its printed remedy, and show the same line clear.
  Proposed AC-4 assigns implementation ownership and independent validation of both per-target MISSING → printed command → same-line OK transactions; AC-1, AC-2 and AC-3 remain byte-identical to the prior body.
- DONE: Retain package versions, command exits, and before/after doctor output as that proof; no global install mutation is authorized.
  AC-4 requires durable full transaction evidence; the isolated ideation spike retains Pi 1.0.0/npm 11.8.0/pi-subagents 0.75.0, exits and all identical old-doctor captures above, explicitly not a candidate pass. Every install mutation stayed inside the disposable home.
- DONE: Resolve the remedy-text scope question explicitly: if the existing text cannot clear its own line, return the needed text or scope change to the gate rather than weakening the outcome.
  Old doctor stays MISSING after successful repair; manifest resolution is still required. The extension command restores real missing declared files; proposed bridge-only “install/update” → executable `pi install npm:pi-subagents` wording and AC-4 are pending captain approval, with no relaxed postcondition.
- DONE: Preserve the previous design and make every acceptance or scope change proposed for the captain's gate.
  AC-1–3, frontmatter and historical reports were byte-compared to HEAD before append; only AC-4, its one-off proof and the bridge-text exception are added as explicit proposals. No production/test/doc-source edits are shipped in ideation.
- DONE: Check the body and existing doctor regression tests.
  `status --read --ac-scan --json` discovers AC-1–4 with report citations; `git diff --check` passes. `go test ./internal/cli -run '^TestPiDoctor' -count=1` passes normally and with `-race`; these existing tests check missing/healthy diagnostics, not the new AC-4 transaction.
- SKIPPED: Complete repository-wide regression results for this body-only fold.
  `go test ./...` exceeded the 120s command budget; `go test ./... -race` was stopped after its 120s budget without a full result. `gofmt -w ./cmd ./internal` ran; its unrelated two-line formatting diff was reversed, leaving no code changes.

### Summary

Folded M4 into the existing ideation body, retaining the manifest-driven design and
adding a proposed integrated remedy acceptance with named proof owners. The real
isolated installer restored each deleted declared target, but the unfixed doctor
still reported stale `.ts` paths; the report does not confuse that negative spike
with completed AC-4, and returns the exact bridge remedy-text exception to the gate.

## Stage Report: implementation

- DONE: The extension probe and the intercom bridge probe resolve from the installed package's own declarations, proven by a fixture that ships compiled output and reports OK, and by the same fixture with the declared entry removed reporting MISSING.
  `resolvePiSubagentsEntries` reads `<pkg>/package.json`; `TestPiSubagentsProbesResolveFromPackageManifest` writes a compiled fixture (`pi.extensions=["./index.js"]`, `exports["./intercom-bridge"].default="./src/api/intercom-bridge.js"`) and asserts `OK` for both lines at the manifest-resolved paths; removing `index.js` → extension `MISSING`, bridge `OK`; removing the bridge target or deleting the `exports` entry → bridge `MISSING`.
- DONE: The integrated remedy acceptance holds: in a disposable home the applicable failing prerequisite appears, its printed remedy runs, and the same line then clears.
  Disposable `mktemp` home (`env -i`), candidate `0.28.0-pre3+dev`, pi 1.0.0/npm 11.8.0/pi-subagents 0.75.0: extension removed → `MISSING` → printed `pi install npm:pi-subagents` (exit 0) → `OK` at declared path, bridge `OK` throughout; bridge removed → `MISSING` → same remedy (exit 0) → `OK`, extension `OK` throughout. Full before/after captured in `_evidence/pi-doctor-probes-stale-subagents-layout/remedy-transactions.txt`.
- DONE: No launcher behaviour changes beyond probe resolution; the shared runtime, the XFAIL bindings, and every other doctor line are untouched.
  `pi.go` diff is one resolver plus wiring at the four call sites and the stale comment; runtime/XFAIL/CI workflow/lane assertions untouched. The only remedy-text change is the in-scope intercom-bridge line, `install/update ...` → `run \`pi install npm:pi-subagents\` ...` (the folded M4 exception the gate approved), asserted by `TestPiSubagentsProbesResolveFromPackageManifest`.

### Summary

Implemented manifest-driven resolution of the two pi-subagents doctor probes:
the extension from `pi.extensions[0]`, the bridge from
`exports["./intercom-bridge"]`, both relative to the package root, failing closed
on a missing/unparseable manifest. Reshaped `pi_frontdoor_test.go` so the fixture
ships compiled output and derives its expected paths from the written
`package.json` (AC-3). Verified the readiness gate passes against the fixture and
that both probes move to `MISSING` when their own declared target is removed.
Docs example paths updated. Focused tests green; noted a pre-existing,
environment-dependent `TestCodexResolveManifestAgainstInstalledHost` failure in
the full suite.


## Review-finding disposition

Validation recommendation: **PASSED (AC-1–4)**. No material outcome or evidence defect found in the promised probe-resolution/remedy behavior; no deferred task defect identified. No candidate repair is proposed or authorized by this report.

- Baseline/environment caveats, not task regressions: `TestCodexResolveManifestAgainstInstalledHost`, `TestVersionAmbiguousMarkersExitZero`, and `TestSurveyCodexPresenceThroughSync` fail unchanged on `dc7d9a0cd` as well as the candidate. The installed Codex cache, ambient `PI_CODING_AGENT`, and agentsview cwd behavior respectively trigger them; none supplies evidence of a violated task value AC. Full-suite results are not claimed green.
- Scope-accounting advisory for the FO: merged code is two files, +287/-20 (net +267), versus the stated +60..+90 ±25% estimate; two documentation files are +4/-4. This exceeds the LOC tolerance but introduces no additional semantic surface. This validator does not retroactively approve the variance or initiate a candidate feedback cycle.

## Stage Report: validation

- DONE: Every acceptance criterion is independently verified against the merged candidate with evidence that can fail, including the manifest-driven resolution and the integrated remedy acceptance.
  Validated `1f41f289f`, tree-identical to merged `origin/main` `12b695f26`; commands, complete captures, executable harnesses, and observations are retained in [_evidence/pi-doctor-probes-stale-subagents-layout/validation-independent.txt](_evidence/pi-doctor-probes-stale-subagents-layout/validation-independent.txt).
- DONE: AC-1 — compiled package declarations drive both OK lines and the actual launcher gate.
  `go test ./internal/cli -run '^TestPiSubagentsProbesResolveFromPackageManifest$' -count=1 -v` exits 0, including `piRuntimeLaunchReady(check)`; independent `spacedock doctor --host pi --plugin-dir <fixture>` and `spacedock pi validation-marker --plugin-dir <fixture>` exit 0 against real installed pi-subagents 0.75.0, with argv-capturing Pi stub and other prerequisites supplied as filesystem fixtures. The pre-change binary refuses that same package (exit 1); the candidate launches with the manifest's `index.js` argument, not a `.ts` source.
- DONE: AC-2 — each declared target's absence fails only its own probe.
  The named Go test passes both removal subtests and the removed-exports subtest; the independent real-package transaction separately removes `pi.extensions[0]` then `exports["./intercom-bridge"].default`, observing own line MISSING / other line OK before each repair.
- DONE: AC-3 — expected paths derive from the written artifact, independently of production probes.
  The Go fixture reads its written `package.json`; removing that write with `go test -overlay <AC3-manifest.json> ./internal/cli -run '^TestPiSubagentsProbesResolveFromPackageManifest/compiled' -count=1 -v` exits 1 at `pi_frontdoor_test.go:710` (fixture manifest missing), not at a compiler error. Independent CLI fixtures relocate declarations to Unicode/space-containing `.mjs`/`.cjs` paths and still launch at the declared extension.
- DONE: AC-4 — the printed remedies reproduce in a disposable home.
  `python3 <scratch>/live.py` exits 0: fresh `env -i` home `/tmp/spacedock-validation-k_wx7s8m`, Pi 1.0.0, Node v24.13.1, npm 11.8.0, pi-subagents 0.75.0 throughout; each exact printed `pi install npm:pi-subagents` exits 0 and moves its line MISSING → OK at its declared path while the other stays OK. Doctor exits remain 1 only for unrelated missing auth/intercom/Spacedock setup. Full stdout/stderr, versions, manifests, commands, exits and candidate SHA-256 are retained; no global mutation.
- DONE: The named falsifiers hold: re-hardcoding either TypeScript source path turns the manifest subtest RED, and removing a probe's own declared target turns that line MISSING while the other stays OK.
  `python3 <scratch>/mutations.py`: both independent stale-.ts overlays make the compiled subtest exit 1 at the corresponding line assertion; accepting unrelated `package.json` instead of each own target makes its removal subtest exit 1; deleting the manifest write exits 1. `live.py` also sees AC-4's exact postcondition go RED after each affected line prints `true` and it exits 0, and after either stale-path binary follows successful installation. Candidate files are never mutated.
- DONE: No launcher behaviour beyond probe resolution changed; the shared runtime, the XFAIL bindings, and every other doctor line are untouched.
  `git diff HEAD^ HEAD --name-only` names only `pi.go`, `pi_frontdoor_test.go`, and the two documentation examples; protected-surface diff is empty. Pre-change/candidate real-home doctor captures compare identically after excluding only the two affected lines and their remedies (including identical stderr/exits). The approved bridge executable-remedy exception is the only wording change.
- DONE: Semantic adversarial matrix and focused regressions.
  `python3 <scratch>/matrix.py` exits 0 across relocated Unicode/spaces, bare-string bridge export, first-of-two extensions, empty extensions/default, absent export, malformed JSON and missing manifest; exact manifest extension reaches the child only in ready cases, all invalid cases refuse without launching. Focused frontdoor/doctor/readiness tests pass normally and with `-race`; resolution is one linear manifest read, not a multiplicative hot path.
- FAILED: Repository-wide regression commands are not wholly green in this environment.
  `go test ./...` exits 1 on the three baseline-reproduced tests above; `go test ./... -race` also exits 1 with the same three failures and no reported data race. No failure is hidden or repaired. Changed Go files pass `gofmt -d`; `gofmt -w ./cmd ./internal` runs only in a candidate archive, exposing an unrelated pre-existing release-test formatting difference without touching the canonical worktree.
- SKIPPED: Real model session / supervisor-talkback validation.
  Explicitly out of scope: the installed-package repairs are real; the ready-gate child is a capture stub and does not claim loaded extensions, credentials, or live talkback.

### Summary

**PASSED:** independently reproduced AC-1–4 and made every named falsifier fail at its intended behavioral boundary. The remedy acceptance holds in a disposable home, the launcher's actual ready gate accepts compiled output, and no other doctor line or shared runtime changed; candidate bytes/HEAD remain untouched, with only this state report and its evidence committed.
