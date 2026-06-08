---
id: qak50v1d6pghavfc0ewg5hjd
title: Binary-missing journey — help the user install spacedock + show they can launch with `spacedock claude`
status: implementation
source: "captain (2026-06-05) — 0.19.6 readiness. When the spacedock binary is absent/non-executable, the user should get a helpful install-and-launch journey (a hint that helps them install the binary, then shows they can launch with `spacedock claude`), not a bare abort."
score: "0.32"
started: 2026-06-08T02:14:39Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-spacedock-binary-missing-install-journey
issue:
sprint: 0198-pre-flip-hardening
group: binary-ux
sprint-readiness: ready
---

This task is the user-facing binary/version/upgrade UX, covering THREE cases the
user actually hits: (a) the binary is MISSING/unusable, (b) the binary is present
but INCOMPATIBLE with the installed plugin, and (c) the plugin is present but the
binary needs installing/upgrading. The captain's directive (2026-06-05): show the
user VERSIONS they understand — skill/plugin version + binary version — not the
internal "contract N / requires >=N,<M" jargon, and make the missing-binary path
an install-and-launch on-ramp (ending in the `spacedock claude` payoff), not a
bare abort.

## Problem

The user-facing copy across these three cases is either internal jargon or
absent:

- **(a) Binary missing.** Two distinct sub-cases by *who* invokes `spacedock`
  (grounded below). When the **user** types `spacedock claude` with no binary on
  PATH, the SHELL emits `command not found: spacedock` (exit 127) — no spacedock
  surface runs at all, because the binary that would print a hint is the thing
  that is missing. When the **FO** mid-session runs `${SPACEDOCK_BIN:-spacedock}
  --version` and it fails, the FO-prose abort
  (`first-officer-shared-core.md` §Startup step 1) carries the install hint (PR
  #262 shipped this) but does NOT show the `spacedock claude` launch payoff.
- **(b) Incompatible.** `internal/contract/contract.go` shows contract integers
  and a half-open range — `Spacedock contract mismatch: binary is contract 1,
  plugin requires >=2,<3.` and the OK path `OK: binary contract 1 satisfies
  plugin range >=1,<2.` — jargon the user does not understand. The manifest
  already carries a display `version` the message discards.
- **(c) Binary stale.** The too-old-binary remedy points at `go install
  …@latest/@next`, not the `brew upgrade spacedock` the captain wants, and does
  not distinguish upgrading the BINARY from refreshing the PLUGIN
  (`spacedock install`).

## Grounding — what the user actually sees TODAY (captured this session)

Riskiest unknown first: ran each surface and captured literal output.

**(a) Binary missing — two sub-surfaces, no rewritable message in sub-case 1.**
- User invokes directly: `PATH=/usr/bin:/bin spacedock claude` → shell
  `command not found: spacedock` / `env: spacedock: No such file or directory`,
  exit 127. No spacedock code or FO prose runs. The only install-and-launch
  surface possible here is the README / `docs/install-journey.md` the user reads
  *before* having a binary. `docs/install-journey.md:72` already shows
  `spacedock claude -- "your task"` as the payoff; README §Install
  (`README.md:14-38`) already carries the brew + go-build lanes.
- FO mid-session: the gate at `first-officer-shared-core.md:10` already aborts
  with `brew install spacedock-dev/homebrew-tap/spacedock` + `go build -o
  spacedock ./cmd/spacedock`, and is explicitly told NOT to run `spacedock
  doctor` (same missing binary). Shipped by archived `binary-absent-fo-bootstrap`
  (PR #262, PASSED). What it lacks: the `spacedock claude` launch payoff.

**(b) Incompatible — the jargon to rewrite.** `spacedock doctor --host claude
--plugin-manifest <fixture>`:
- too-old-binary (plugin `>=2,<3`, binary contract 1) → stderr, exit 1:
  `Spacedock contract mismatch: binary is contract 1, plugin requires >=2,<3.`
  + `too-old-binary: your spacedock binary (contract 1) predates this plugin
  (needs >=2,<3). Rebuild/upgrade spacedock: go install …@next (or pull and 'go
  build').` + `Run \`spacedock doctor\` for details.`
- too-old-plugin (plugin `>=0,<1`) → `Spacedock contract mismatch: binary is
  contract 1, plugin requires >=0,<1.` + `too-old-plugin: … Update it: spacedock
  install --host claude (or 'claude plugin update spacedock').`
- compatible OK path → stdout, exit 0: `OK: binary contract 1 satisfies plugin
  range >=1,<2.`

**(c) Binary stale.** Same too-old-binary remedy line as above — `go install`,
not `brew upgrade`.

**Surface map (each case → its real surface):**
| Case | Real surface | File:line |
| --- | --- | --- |
| a / user-invokes | shell `command not found` (no spacedock surface) → docs | `README.md:14`, `docs/install-journey.md:72` |
| a / FO mid-session | FO-prose abort (model-ingested) | `first-officer-shared-core.md:10` |
| b incompatible (live launch) | `gateHost` stderr | `internal/cli/frontdoor.go:138` |
| b incompatible / c stale (doctor) | `mismatchMessage` + remedies + OK path | `internal/contract/contract.go:148,156,166,180` |

## Riskiest mechanism — spiked this session (the version rewrite is feasible)

The case-(b)/(c) redesign rests on one unverified mechanism: **can the manifest's
display `version` be surfaced in the message, and can the binary's version reach
the message functions?**

- Spike (throwaway Go, `/tmp/sd-spike`): unmarshalling `{version,
  requires-contract}` from the real `too-old-binary.json` fixture and the live
  `.claude-plugin/plugin.json` returns `version="0.13.0"` / `"0.19.6"` alongside
  `requires-contract`. `readRequiresContract` (`doctor.go:22`) today discards
  `version`; surfacing it is a one-field unmarshal change. **Mechanism proven.**
- Binary version: `Version` lives in `internal/cli` (`cli.go:28`); `contract`
  does NOT import `cli` and `cli` imports `contract` (verified — no cycle). So
  `contract`'s message functions cannot reach `cli.Version` without an import
  cycle — the binary version MUST be threaded as a parameter from the cli
  callers (`frontdoor.go:130`, `init.go:48,91,101`) into
  `ManifestVerdict`/`RunDoctor`/`compareWithManifest`. Bounded, no cycle.
- Constraint: the `(contract N)` token in `spacedock --version`
  (`cli.go:502`) is load-bearing — the FO/ensign gate parses it. The rewrite
  changes the **doctor/gate messages only**, NOT `--version`.

## Direction — where each surface lives + exact before/after copy

### (a) Binary missing — NO new abort message to "fix"; add the launch payoff where a surface exists

The user-invokes sub-case has no spacedock surface to rewrite (shell owns
`command not found`); the install-and-launch journey there already lives in
README + `docs/install-journey.md`, which already end in `spacedock claude`. The
only delta this entity ships for case (a) is on the FO mid-session surface:
append the launch payoff to the FO-prose abort so a user who installs from the
hint sees the next step.

`first-officer-shared-core.md` §Startup step 1, binary-absent class.
BEFORE (line 10, tail): `… Do NOT run \`spacedock doctor\` — same missing
binary.`
AFTER: `… Do NOT run \`spacedock doctor\` — same missing binary. Once spacedock
is on PATH, launch with \`spacedock claude\` to start your first officer.`

(Note: this is a text change to a model-ingested contract file. Per the workflow's
tautology rule, a substring check over this file is NOT a behavioral AC; it is
admitted only as the same legitimate-text exception PR #262 used — the binary is
absent, so the contract prose is the only artifact present at the failure. See
AC-3 framing. The behavioral weight of this entity is in (b)/(c) below.)

### (b) Incompatible — rewrite the contract-jargon messages to show VERSIONS

`internal/contract/contract.go`. Thread the plugin's manifest `version` and the
binary `Version` into the messages; drop the contract integers and the
`>=N,<M` range from the user-facing lines.

`mismatchMessage` (`contract.go:156`), too-old-binary case:
BEFORE: `Spacedock contract mismatch: binary is contract 1, plugin requires
>=2,<3.` + `too-old-binary: your spacedock binary (contract 1) predates this
plugin (needs >=2,<3). Rebuild/upgrade spacedock: go install …@next …` + `Run
\`spacedock doctor\` for details.`
AFTER (comm-officer polish):
```
Spacedock version mismatch: binary 0.19.4, plugin 0.21.0. The plugin needs a newer binary.
  Upgrade the binary: brew upgrade spacedock
  Or build from source: go build -o spacedock ./cmd/spacedock
Run spacedock doctor for details.
```

too-old-plugin case:
BEFORE: `Spacedock contract mismatch: binary is contract 1, plugin requires
>=0,<1.` + `too-old-plugin: … Update it: spacedock install --host claude …`
AFTER (comm-officer polish):
```
Spacedock version mismatch: binary 0.21.0, plugin 0.18.0. The binary needs a newer plugin.
  Update the plugin: spacedock install --host claude
Run spacedock doctor for details.
```

OK path (`contract.go:148`):
BEFORE: `OK: binary contract 1 satisfies plugin range >=1,<2.`
AFTER: `OK: spacedock binary 0.19.4 and plugin 0.19.6 are compatible.`

(The internal `contract` token stays in `--version` and in the malformed-range
packaging-bug message — that message names the manifest file for a packager, not
a version for an end user, so it is out of scope for the version rewrite. The
verdict String() kebab tokens stay — they are the test oracle, not user copy.)

### (c) Binary stale — brew upgrade, distinct from the plugin refresh

`tooOldBinaryRemedy` (`contract.go:166`). The too-old-binary remedy (folded into
(b)'s too-old-binary AFTER copy above) leads with `brew upgrade spacedock` and
keeps the source-build fallback, and names the distinction:
BEFORE: `Rebuild/upgrade spacedock: go install
github.com/spacedock-dev/spacedock/cmd/spacedock@next (or pull and 'go build').`
AFTER (comm-officer polish):
```
To upgrade the binary: brew upgrade spacedock (Homebrew), or go build -o spacedock ./cmd/spacedock (from a clone).
To refresh the plugin instead: spacedock install.
```

(The `@next`/`@latest` branch-pin suffix is dropped from the user-facing line:
`brew upgrade spacedock` carries no branch, and the source-build line is
branch-agnostic. The devBranch pin remains relevant only to the auto-install /
codex-add paths, which this entity does not touch.)

## Acceptance criteria

**AC-1 — an incompatible binary/plugin pair reports VERSIONS, not contract
jargon.** `spacedock doctor` (and the live launch gate) for a too-old-binary and
a too-old-plugin manifest emit a message that contains the plugin's display
version AND the binary's display version, and contains NEITHER the substring
`contract ` followed by an integer NOR a `>=N,<M` range, in the user-facing
mismatch line.
- Verified by: a behavior fixture test driving `RunDoctor(manifestPath, host,
  binaryVersion, …)` against a too-old-binary fixture (plugin `version` X,
  `requires-contract >=2,<3`) and a too-old-plugin fixture, asserting (i) exit 1,
  (ii) stderr contains both version strings (the fixture's `version` and the
  passed binary version), (iii) stderr does NOT match the regex
  `contract \d` and does NOT contain `>=` ... `,<`. The expected values (the two
  versions) come from the fixture + the test's passed binary version, NOT from
  the message file — independent oracle, not a tautology.

**AC-2 — the compatible OK path and the too-old-binary remedy show versions +
the brew upgrade path.** The OK path names both versions; the too-old-binary
remedy leads with `brew upgrade spacedock` and names the binary-vs-plugin
distinction.
- Verified by: a behavior fixture test asserting (i) the compatible fixture
  (`>=1,<2`) exits 0 and stdout contains both the plugin `version` and the passed
  binary version and does NOT match `contract \d`; (ii) the too-old-binary
  fixture's stderr contains `brew upgrade spacedock` and contains `spacedock
  install` (the plugin-refresh distinction). Oracle independent of the message
  file.

**AC-3 — the FO binary-absent abort points at the launch payoff.**
`first-officer-shared-core.md` §Startup step 1 binary-absent class names
`spacedock claude` as the next step after install, and still omits `spacedock
doctor` from that class.
- Verified by: a presence/absence text check over the real file (step-1
  binary-absent class contains `spacedock claude` AND does not contain `spacedock
  doctor` within that class). Framing: this is the SAME legitimate-text exception
  PR #262 used and the gate already accepted — the binary is absent at the
  failure, so the contract prose IS the only artifact present; there is no
  runnable binary to drive. It is NOT a stand-in for a behavioral guarantee a
  code gate could enforce. The behavioral weight of the entity is AC-1/AC-2,
  which ARE code-gated.

## Test plan

- **AC-1 / AC-2 (the behavioral core):** Go fixture tests in
  `internal/contract/doctor_test.go`, extending the existing table-driven
  `TestDoctorVerdicts` harness (already drives `RunDoctor` over `testdata/*.json`
  fixtures with exit-code + substring assertions). The existing fixtures already
  carry `version` fields (`0.13.0`, `0.10.0`, `0.12.1`) — reuse them; pass a
  fixed binary version (e.g. `"0.19.4"`) through the new param. Add negative
  assertions (regex `contract \d` absent; `>=`/`,<` absent) so a paraphrase that
  re-introduces jargon FAILS and an inverted message FAILS. Cheap, offline, no
  host CLI. The oracle is the fixture's `version` + the test's passed binary
  version — values that live OUTSIDE the message under test, so the check tracks
  behavior, not spelling.
- **gate parity:** one assertion that the live launch gate (`gateHost` in
  `frontdoor.go`) emits the same version-bearing message for a mismatch (the
  existing frontdoor test seam), so the rewrite is not doctor-only.
- **AC-3:** a presence/absence string check over
  `first-officer-shared-core.md` step 1 (the admitted text exception). Cheap.
- **No spike needed beyond the one already run:** the version-plumbing mechanism
  (manifest `version` readable; binary version threaded as a param, no import
  cycle) was exercised this session (`/tmp/sd-spike`). Estimated cost: small —
  one field added to the manifest reader, one param threaded through
  `ManifestVerdict`/`RunDoctor`/`compareWithManifest` and its 4 cli callers,
  message-string edits, ~6 fixture-test cases.

## Out of scope (separate tasks — do NOT fold in)

- **Flip release-mechanics** (separate task): the goreleaser `devBranch=next→main`
  stamp and confirming the contract-1 0.19.6→0.20.0 upgrade end-to-end (`brew
  upgrade` + `spacedock install` re-points to `main`). The 0.20.0 flip is NO
  contract change (stays contract 1), so a 0.19.6→0.20.0 upgrade is Compatible —
  the skew/jargon path is NOT triggered by the flip. This entity is UX hardening
  for future bumps + any incompatibility, not flip-blocking.
- `44 bundle-asset-distribution` (bundle the plugin into the binary — structural).
- `5w notarize-macos-release` (Gatekeeper).
- `s0cq install-marketplace-ref-refresh`.
- Renaming the `--skip-contract-check` FLAG (a separate question — the flag name
  stays; only user-facing MESSAGES change).
- The malformed-range packaging-bug message (names the manifest file for a
  packager, not a version for an end user) and the verdict String() kebab tokens
  (the test oracle).

## Stage Report: ideation

- DONE: Settle WHERE each surface lives and the EXACT before/after copy for the three user-facing cases (a/b/c)
  Surface map table + before/after copy in "Direction"; (a) FO-prose payoff line, (b) contract.go:148/156 version rewrite, (c) tooOldBinaryRemedy brew-upgrade rewrite.
- DONE: Ground the riskiest unknown FIRST — run each surface to capture what the user ACTUALLY sees today and map each case to its real surface
  Captured literal output: shell `command not found` exit 127 (user-invokes), FO-prose abort line 10, `spacedock doctor` mismatch/OK strings (too-old-binary/plugin/compatible). Reconciled with FO gate (already emits brew/go-build, told NOT to run doctor when binary missing). Discovered PR #262 already shipped case-(a) FO hint.
- DONE: Write entity-level ACs + test plan: behavior fixtures asserting emitted guidance shows VERSIONS (not contract) + correct install/upgrade/launch payoff (output/exit-code assertions)
  AC-1/AC-2 are code-gated fixture tests over `RunDoctor` + testdata (oracle = fixture `version` + passed binary version, independent of message file, with negative `contract \d` / `>=,<` assertions); AC-3 is the admitted text exception for the binary-absent FO prose.
- DONE: Record what is OUT of scope
  Out-of-scope section: flip release-mechanics (goreleaser devBranch stamp + contract-1 upgrade-confirm), 44 bundle, 5w notarize, s0cq, --skip-contract-check rename, malformed-range/kebab-token messages.

### Summary

Expanded the entity from the single binary-MISSING AC to the captain's full
three-case binary/version/upgrade UX. Grounded every surface against real output
first: the decisive finding is that the user-invokes missing-binary case has NO
spacedock surface (shell `command not found`), so case (a)'s journey already lives
in README/install-journey + the FO-prose abort (PR #262), and the only new (a)
delta is appending the `spacedock claude` payoff to the FO abort. The behavioral
weight moved to cases (b)/(c): rewrite the `contract.go` doctor/gate messages to
show plugin version + binary version (spiked feasible — manifest `version`
readable, binary version threaded as a param since `contract` cannot import
`cli`) and lead the stale-binary remedy with `brew upgrade spacedock`. ACs are
code-gated fixture tests with oracles independent of the message file; AC-3 is the
same legitimate-text exception PR #262 used. comm-officer polish arrived after the
first commit and was incorporated into the (b)/(c) AFTER copy (tighter: cut
filler, positive forms, scannable command lines); the polished copy still
satisfies the AC oracles (both versions present, no `contract N` / `>=,<`, brew
upgrade leads, `spacedock install` distinction kept), so the ACs were unchanged.

