---
title: Forward appended Safehouse profiles through the front door
status: implementation
source: Captain request 2026-09-22
started: 2026-09-22T15:38:57Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-frontdoor-safehouse-append-profile
issue:
pr:
mod-block:
id: 9cbmcq3yyfjynhms1yd12drh
gates:
    version: 1
    records:
        - id: gate:9cbmcq3yyfjynhms1yd12drh:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:9cbmcq3yyfjynhms1yd12drh-ideation-1
              briefing:
                id: briefing:9cbmcq3yyfjynhms1yd12drh:ideation:attempt-1:revision-1
                digest: sha256:c3bf30c9d2d181a9bf249de314ac8d2e9010eff21d1ba935866a9cdaedce2414
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:9cbmcq3yyfjynhms1yd12drh:ideation:1
                briefing: briefing:9cbmcq3yyfjynhms1yd12drh:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-22T16:43:11.951568Z"
                decision: approve
                reason: Captain approved the prepared ideation design in binding Subspace review resolution:binding-1790095375131026000.
              application:
                target-stage: implementation
                state: consumed
---

Support `--safehouse-append-profile=file.sb` on `spacedock claude`, `spacedock codex`, and `spacedock pi`.

## Problem

The front doors reject this option before Safehouse can load an additional policy file. Users must bypass Spacedock to supply the equivalent Safehouse option.

## Proposed approach

Add one repeatable StringArray flag to the existing parsers and one translation case. Each value becomes one literal `--append-profile=VALUE` argument before Safehouse's `--`. The option alone selects the existing sandbox launch path.

| Owner | Bounded change |
|---|---|
| `internal/cli/frontdoor.go` | Add the flag to `frontDoorFlags` and `bindFrontDoorFlags`; collect values in `parseFrontDoorArgs`. Claude and Codex share these owners. |
| `internal/cli/pi.go` | Add the same flag and collection loop to `parsePiFrontDoorArgs`. |
| `internal/cli/help.go` | Add the same StringArray to Pi's separate help registration. Claude/Codex help uses the shared binding. |
| `internal/safehouse/safehouse.go` | Translate `append-profile=VALUE` to `--append-profile=VALUE`. |

These additions serve AC-1 and AC-2. The simpler alternative, host passthrough after `--`, cannot work: Safehouse never receives those host arguments. A generic pass-through parser or policy engine adds no required behavior. Existing launch functions and `Wrap` already carry the translated arguments.

### Path, order, and error contract

- Both `--safehouse-append-profile=file.sb` and `--safehouse-append-profile file.sb` work before the host delimiter.
- Each occurrence is one path. Preserve spaces, commas, colons, equals signs, duplicates, and profile occurrence order. Use StringArray, not a comma-splitting StringSlice.
- Spacedock performs no expansion, splitting, shell evaluation, existence check, or path normalization. The child inherits the caller's current directory. Relative CLI paths therefore resolve from the launch directory, not the workflow directory, repository root, or `.safehouse` location.
- Safehouse owns tilde expansion, normalization, file validation, and policy parsing. Quoted shell-like text remains argument data in Spacedock. The shell can expand unquoted input before Spacedock receives it.
- Keep the current grouping of different knobs: enable, add-dirs, add-dirs-ro, then append-profile. Profile occurrences retain their own order even when other knobs appear between them. Do not add an order-preserving parser for unrelated knobs.
- The existing `--trust-workdir-config` remains. Safehouse renders generated grants, project-config profiles, then CLI profiles. Safehouse's final profile/config write protections remain last. Pi's existing automatic directory grant still uses this generated-grant phase.
- A standalone `--` that is not a flag value ends Spacedock parsing. Subsequent tokens retain existing host passthrough behavior and do not select the sandbox.
- Preserve pflag value rules. A bare option at end fails parsing. A following token beginning with `-` is still a space-form value; use equals form to make such paths clear. Do not introduce a new missing-value heuristic.
- Explicit empty values reach Safehouse as `--append-profile=` and fail there. Missing, nonregular, unreadable, and invalid policy files remain Safehouse errors. Spacedock propagates failure and never retries without the profile. A missing Safehouse executable follows the existing availability error.
- Already inside Safehouse: retain the current wrap decision. `Inside` currently affects banner/status text, not launch selection. The new flag still selects wrapping and remains in outer argv. A nested sandbox attempt can fail; Spacedock must neither discard the option nor claim it changes the parent sandbox. No new nested-sandbox bypass or escalation.

## Risk evidence

Source baseline: `origin/main` = `9a6765fa8aed1141e6f5c1d38868b78c4975060a`; checkout HEAD = `438053493838dc70c9478b3d991309d566783e85`. The named parser, translator, launch, and sandbox-doc owners match that baseline. The root checkout has unrelated untracked files; this task changes only its shared-state directory.

Read `docs/runtime-support.md` before the probe. The installed `/Users/clkao/.local/bin/safehouse` denied access, including `ls`; this is an environment limitation, not evidence that the option is unsupported.

The independent source is [Safehouse commit e376993](https://github.com/eugene1g/agent-safehouse/tree/e376993ee8e15c4e4b3aa3a2ee282f15f6e3c680), version 0.12.0. Its real standalone script ran `--help` and `--stdout`; no model or sandboxed command ran. See [probe.py](probe.py) for the repeatable probe and [probe-results.txt](probe-results.txt) for observed results.

The policy output contains `PROBE_PROJECT`, `PROBE_A`, `PROBE_B`, `PROBE_B` in that order. The CLI files exist only in the caller directory; `--workdir` points elsewhere. This exercises relative-path base, spaces/punctuation, both flag forms, repetitions, duplicates, and composition with project config. Empty, missing, directory, and absent values each exit 1 with a specific error.

Upstream owners: `bin/lib/cli/parse.sh` collects each path; `bin/lib/policy/request.sh:policy_request_resolve_append_profile_paths` validates files. `bin/lib/policy/render.sh:policy_render_emit_dynamic_sections` orders project profiles before CLI profiles and final protections. This proves policy composition, not nested Seatbelt execution or installed-version compatibility.

## Out of scope

Generic sandbox interfaces, policy engines, profile-content changes, new environment options, dependency upgrades, install behavior, auth, permission defaults, and nested-sandbox redesign. Ideation includes no product edits, model sessions, CI, global config changes, or pushes.

## Expected surface and tolerance

Estimate net LOC change: **+195, across 9 files**. Estimate **205 insertions and 10 deletions**. Tolerance: net **+120 through +270 LOC**, **8 through 10 files**. This measures the implementation diff against its starting `origin/main`; shared-state ideation artifacts are separate.

Four production files: `frontdoor.go`, `pi.go`, `help.go`, `internal/safehouse/safehouse.go` (about +25 net). Four proof files: `safehouse_knob_test.go`, `pi_frontdoor_test.go`, `help_test.go`, `internal/safehouse/safehouse_test.go` (about +145 net). One documentation file: `docs/site/reference/sandbox.md` (about +25 net). An extra focused parser-test file is within tolerance only if the existing proof owners cannot keep the cases readable.

Allowed semantic changes: one new repeatable pre-delimiter CLI option; recognition now selects the existing sandbox path and appends the requested policy. Help and sandbox documentation describe it. Stored formats, FO/ensign authority, host arguments, permission defaults, existing knobs, and other runtime behavior stay unchanged.

## Acceptance criteria

**AC-1 — All three front doors deliver each requested profile to Safehouse without losing host arguments.**
The independent baseline rejects the option. After the change, a profile-only invocation without `.safehouse` records the expected Safehouse argv for Claude, Codex, and Pi. Each profile appears exactly once per occurrence before Safehouse's delimiter. Existing host command, prompt, environment forwarding, and host permission flags match the equivalent existing sandbox launch. Verified by deterministic launch-seam fixtures; removing flag registration, translation, or sandbox selection makes this fail.

**AC-2 — Profile values retain their path and composition meaning.**
Both forms, relative and absolute paths, spaces, punctuation, duplicates, and interleaved repetitions reach Safehouse as literal ordered values. The committed independent probe demonstrates caller-relative CLI resolution and project-before-CLI composition. Verified by parser-to-translator tables, launch fixtures, and the real Safehouse policy probe; comma splitting, path rebasing, shell evaluation, or deduplication makes a relevant assertion fail.

**AC-3 — Malformed input fails without an unprofiled fallback, and argument boundaries stay compatible.**
A missing terminal value fails before Launch. An empty value remains present and Safehouse rejects it. Safehouse file/policy errors propagate; no retry launches the host without the profile. Post-delimiter tokens remain host arguments. Existing inside-Safehouse signals do not suppress the requested wrapper or profile. Verified by missing-value/no-Launch fixtures, exact empty-value forwarding, propagated nonzero exit fixtures, delimiter fixtures, and inside-marker fixtures. A silent drop, fallback, or delimiter leak fails these assertions.

**AC-4 — Each front door advertises the option and its repeatability.**
Rendered CLI help includes the new flag for all three hosts. The sandbox reference supplies the exact example and path/order contract. Verified by existing help output tests and documentation review against the implementation and upstream probe. Omitting Pi's separate help registration fails the help assertion.

## Test plan

Write focused failing cases before implementation. Extend existing proof owners rather than introducing a separate parser framework.

| Proof owner | Cases and falsifying edit | Cost |
|---|---|---|
| `safehouse_knob_test.go` / `TestSafehouseKnobFormsEquivalent` | Both forms, special characters, interleaved repeats, duplicates, empty value. Splitting or rewriting a value fails exact argv comparison. | Deterministic, small table extension. |
| `safehouse_knob_test.go` launch fixtures | Claude/Codex profile-only selection, no leakage, delimiter, missing terminal value, inside marker, missing binary, nonzero exit. A removed wrapper or ignored failure fails observed Launch/return assertions. | Deterministic, reuse fakeHost; modest. |
| `pi_frontdoor_test.go` / existing Safehouse parser and launch tests | Same profile contract through Pi's separate parser; keep Pi resource flags, markers, and existing directory grants. Missing Pi registration or collection fails observed argv. | Deterministic, reuse Pi fake; modest. |
| `internal/safehouse/safehouse_test.go` / `TestTranslateFlags` | Append-profile translation, literal empty and punctuation values, order; retain unknown-key error. Missing translator case fails. | Deterministic, small table extension. |
| `help_test.go` and `pi_frontdoor_test.go` | Actual help output exposes the repeatable flag for all hosts. Missing help binding fails. | Deterministic text claim; small. |
| `probe.py` against pinned upstream script | Policy markers prove current-directory resolution and project/A/B/B order; invalid file cases fail. Moving the CLI base to workdir or reordering composition fails. | One-off independent policy probe, about 1 second; no model. |

Implementation runs `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` from its isolated worktree. No live model run is needed: host argument forwarding is the claim. A harmless real policy render supplies the external mechanism proof; do not substitute a fake Safehouse for that claim.

## Proposed documentation diff

In `docs/site/reference/sandbox.md`, replace the existing trigger cell:

```text
A `.safehouse` profile in the working directory, or the `--safehouse` flag
```

with:

```text
A `.safehouse` profile in the working directory, `--safehouse`, or any `--safehouse-*` option
```

Append this exact text:

````markdown
## Additional profiles

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `--safehouse-append-profile PATH` | Optional, repeatable path | Append a safehouse policy file for this launch. | None |

```bash
spacedock claude --safehouse-append-profile=file.sb
spacedock codex --safehouse-append-profile="profiles/local rules.sb"
spacedock pi --safehouse-append-profile=first.sb --safehouse-append-profile=second.sb
```

Both `--safehouse-append-profile=PATH` and `--safehouse-append-profile PATH` select safehouse.
Place the option before `--`; tokens after `--` go to the coding agent.
Relative paths start in the directory where you launch Spacedock.
Each occurrence supplies one path; repeats retain their order.
Safehouse loads project-config profiles before these profiles, then applies its final write protections.
Safehouse reports invalid profile paths or content, and Spacedock returns the failure.
An existing sandbox stays active; this option still requests a safehouse launch and cannot relax the parent sandbox.
````

New help description in the shared binding and both Pi registrations:

```text
Append a safehouse policy file; repeatable; relative paths use the launch directory
```

### Feedback Cycles

## Stage Report: ideation

- DONE: Prove the actual Safehouse append-profile option and map minimal forwarding across existing supported front doors.
  Upstream 0.12.0 `--help` and real `--stdout` probe passed; `probe.py` asserts project/A/B/B output, cwd-relative CLI paths, punctuation, and four error exits.
- DONE: Produce a bounded design with explicit path/order/error semantics, testable acceptance criteria, concrete doc diff, and net LOC/file estimate with tolerance.
  Task body specifies the four production owners, four proof owners, one doc page, +195 net LOC across 9 files, and exact semantic boundaries.
- DONE: Exercise existing forwarding proof owners.
  Focused `TestTranslateFlags`, `TestSafehouseKnobFormsEquivalent`, and `TestPiFrontDoorAcceptsSafehouseFlags` passed; removing current knob forwarding fails their exact argv assertions.
- FAILED: Repository-wide normal and race checks.
  Both encountered undefined symbols in pre-existing untracked `tmp/dispatch-027/stash-recovery-copy/untracked-files`; race also reported installed-Codex manifest mismatch. Remaining runs were interrupted after these unrelated failures.
- SKIPPED: `gofmt -w ./cmd ./internal`.
  Read-only `gofmt -l` found existing formatting in `internal/release/runtime_live_evidence_workflow_test.go`; changing product files violates this ideation assignment.
- SKIPPED: Installed Safehouse execution and nested-sandbox execution.
  Installed path access was denied; pinned upstream policy-only execution proves the required option without a model or sandbox session.
- SKIPPED: State push.
  Captain's explicit no-push instruction overrides generic split-root sync guidance; task paths are committed locally.

### Summary

The smallest design extends existing StringArray parsing and Safehouse translation for Claude, Codex, and Pi. A real Safehouse policy probe establishes literal path forwarding, caller-relative paths, repetition order, and composition with project profiles. The task includes acceptance criteria, concrete documentation wording, and an implementation estimate; no product files or frontmatter changed.


## Stage Report: implementation

- DONE: Implement the approved literal repeatable append-profile option across Claude, Codex and Pi using existing parsing and translation owners, within approved surface.
  Code commit `83b356c1c`, based on clean rebase to `9a6765fa8`; 9 files, 172 insertions/9 deletions (+163 net), within +120..+270 and 8..10 tolerance.
- DONE: Prove argv, path/order, delimiter and failure behavior with existing tests written red before implementation; apply approved help and doc changes.
  Focused tests first failed for unregistered flag, missing translation, and absent help; the same cases passed after implementation, normally and under `-race`.
- DONE: Commit the finished implementation and canonical report with required normal/race/format evidence and explicit limitations.
  Product committed on `spacedock-ensign/frontdoor-safehouse-append-profile`; this report is committed path-scoped in the shared state checkout; full-check failures are recorded below.
- DONE: AC-1/AC-3 launch proof for all three hosts.
  `TestAppendProfileLaunchContract` compares complete argv/env with the existing sandbox launch and counts Launch calls; dropping profiles, changing host flags/prompt/env, suppressing inside wrapping, or retrying after exit 23 fails it.
- DONE: AC-2 literal and composition proof.
  `TestAppendProfileLiteralParsing` and `TestTranslateFlags` assert both forms, relative/absolute and shell-like paths, punctuation, empty/dash values, duplicates, and grouped knob order; splitting, rebasing, evaluation, or deduplication fails exact comparisons.
- DONE: AC-3 argument-boundary and error proof.
  Launch fixtures assert missing terminal values and unavailable Safehouse produce no Launch, empty values remain present, and post-delimiter values go to the unwrapped host; swallowing a value, fallback, or delimiter leakage fails them.
- DONE: AC-4 advertised option and documentation.
  `TestFrontDoorHelpCarriesDetail` and `TestPiHelpCarriesSafehouseDetail` require the flag and repeatability/path description; missing Pi's separate help registration fails; sandbox reference uses the approved example and contract text.
- DONE: Independent Safehouse policy probe.
  Re-ran committed `probe.py` against upstream `e376993ee8e15c4e4b3aa3a2ee282f15f6e3c680` standalone script: project/A/B/B order and caller-relative punctuation paths passed; empty, missing, directory, and absent values exited 1 as expected.
- DONE: Run and record required normal/race checks (`go test ./...` and `go test ./... -race`).
  Both completed exit 1 solely at existing `TestCodexResolveManifestAgainstInstalledHost` (`codex_resolve_test.go:44`): stable ID absent while resolver returns local `spacedock-local/spacedock/0.28.0-pre0` manifest; all other packages passed, with no race report.
- DONE: `gofmt -w ./cmd ./internal` and owned-file format verification.
  Required formatter ran; it touched two pre-existing unrelated struct-field spacing lines in `internal/release/runtime_live_evidence_workflow_test.go`; FO explicitly declined cleanup, so only those formatter-induced bytes were restored. All eight owned Go files pass `gofmt -l`; `git diff --check` is clean.
- SKIPPED: Unrelated installed-host resolver repair.
  FO explicitly declined: previously observed mismatch, outside append-profile ownership; deferred for this task, promote when stable/local install resolution is assigned; no append-profile AC harm established.
- SKIPPED: Installed Safehouse, nested Seatbelt execution, local model sessions, CI, and pushes.
  Installed Safehouse was denied in ideation; pinned policy rendering and deterministic launch seams prove the assigned claim, not installed-version compatibility or nested execution. Captain's no-push constraint applies to code and state.

### Summary

The existing StringArray parsers and translator now forward each profile literally before Safehouse's delimiter and select the existing sandbox path. The bounded implementation, help, and documentation are committed; focused normal/race and independent policy evidence passed. Both required full suites completed with the known installed-Codex resolver mismatch, which the first officer explicitly declined to repair in this task.
