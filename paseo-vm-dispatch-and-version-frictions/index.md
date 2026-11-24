---
id: kdpzm9pckwj4pm1djvaf01ce
title: "Paseo stack Layer 3: dev build SHA in --version, collision-free immutable dispatch files, and --dispatch-dir"
status: validation
source: "CL via MetaFO (2026-11-23, Studfleet VM dogfooding /tmp/paseo-frictions-studfleet-vm-2026-11-24.md)"
started: 2026-11-24T06:53:43Z
completed:
verdict:
score: 0.90
worktree: .worktrees/spacedock-ensign-paseo-vm-dispatch-and-version-frictions
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:kdpzm9pckwj4pm1djvaf01ce:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:kdpzm9pckwj4pm1djvaf01ce-backlog-1
              briefing:
                id: briefing:kdpzm9pckwj4pm1djvaf01ce:backlog:attempt-1:revision-1
                digest: sha256:f15aac9719b6ca4f73dc627f7e3a36a8e5c2eb235554cb19cd92c03cb495e42e
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:kdpzm9pckwj4pm1djvaf01ce:backlog:1
                briefing: briefing:kdpzm9pckwj4pm1djvaf01ce:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-24T06:52:56.224469Z"
                decision: approve
                reason: 'Continuing Paseo PR stack atop PR #897 (5210d36d0) per CL authorization; merge/release/tag remain CL gates.'
                conn:
                    quote: 'the driver is authorized to continue the paseo PR stack to fix those frictions. This word of CL covers the continuation of that stack as it was done until now. It does not cover a merge, a release or a tag: those stay gates of CL.'
                    source: CL via MetaFO message to driver 846d9a2 (2026-11-23 22:47 PT)
              application:
                target-stage: ideation
                state: consumed
        - id: gate:kdpzm9pckwj4pm1djvaf01ce:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:kdpzm9pckwj4pm1djvaf01ce-ideation-1
              briefing:
                id: briefing:kdpzm9pckwj4pm1djvaf01ce:ideation:attempt-1:revision-1
                digest: sha256:5ce69114377b181a60c57a86bb14a14127af7264c668d78f6d03de6b61bdc95e
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:kdpzm9pckwj4pm1djvaf01ce:ideation:1
                briefing: briefing:kdpzm9pckwj4pm1djvaf01ce:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-24T07:02:48.991783Z"
                decision: approve
                reason: 'Ideation design covers all Studfleet VM Spacedock frictions (#1, #4, #5, plus 3-line operational notes #3/#5/#6) with concrete falsifying ACs, atomic link publication, golden raw-digest verification, and untouched first-officer-shared-core.md.'
                conn:
                    quote: 'the driver is authorized to continue the paseo PR stack to fix those frictions. This word of CL covers the continuation of that stack as it was done until now. It does not cover a merge, a release or a tag: those stay gates of CL.'
                    source: CL via MetaFO message to driver 846d9a2 (2026-11-23 22:47 PT)
              application:
                target-stage: implementation
                state: consumed
        - id: gate:kdpzm9pckwj4pm1djvaf01ce:validation
          stage: validation
          attempts:
            - id: gate-attempt:kdpzm9pckwj4pm1djvaf01ce-validation-1
              briefing:
                id: briefing:kdpzm9pckwj4pm1djvaf01ce:validation:attempt-1:revision-1
                digest: sha256:16a5d6a893df55f59ff4adf510b8c5230943b2f9800af007167396a6d67c5270
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:kdpzm9pckwj4pm1djvaf01ce:validation:1
                briefing: briefing:kdpzm9pckwj4pm1djvaf01ce:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-24T07:58:01.867142Z"
                decision: approve
                reason: Cycle 2 validation PASSED on 03a796c48 across AC-1..AC-5 with 13/13 detached-audit mutations caught, full normal+race coverage, strict MkDocs, untouched 23,464 B shared core, and clean standalone-clone live Paseo provenance.
                conn:
                    quote: 'the driver is authorized to continue the paseo PR stack to fix those frictions. This word of CL covers the continuation of that stack as it was done until now. It does not cover a merge, a release or a tag: those stay gates of CL.'
                    source: CL via MetaFO message to driver 846d9a2 (2026-11-23 22:47 PT)
              application:
                target-stage: done
                state: pending
---

Identify the running source build and preserve each dispatch assignment. Let first officers select a directory that their workers can read.
Layer 3 targets PR #897 at `5210d36d092af463164b5f1c73ba24264ebe5354` (`origin/paseo-and-generic-harness-fo-support`).

## Problem

1. Friction #1: an unstamped binary reports `<manifest-version>+dev`. Operators cannot identify the installed stack commit from `--version`.
2. Friction #4: two briefs with the same parent, entity, and stage use one pathname. The second write replaces the first assignment.
3. Friction #5: the fixed `/tmp/spacedock-dispatch` path can name different storage in the parent and worker namespaces.
4. Paseo operators need three short notes about shared storage, workspace adoption, and the actual thinking setting.

## Proposed approach

### Dev version

`displayVersion()` retains the first branch: `Version != "dev"` returns the linker value unchanged.
An unstamped build reads the embedded manifest and Go VCS settings. A package seam with signature `func() (*debug.BuildInfo, bool)` defaults to `debug.ReadBuildInfo`.
Tests replace this seam and restore it with cleanup. Those tests do not run in parallel.

With a valid `vcs.revision`, append its first seven hexadecimal characters inside build metadata: `0.29.0-pre0+dev.5210d36`.
With `vcs.modified == "true"`, append `.dirty`: `0.29.0-pre0+dev.5210d36.dirty`.
Accept hexadecimal revisions with at least seven characters and convert them to lowercase.
Absent build information, absent revisions, short revisions, or malformed revisions retain `<manifest-version>+dev`. A dirty flag without a valid revision adds nothing.
A malformed embedded manifest retains the existing `dev` fallback. No Git subprocess runs during `--version`.

Line 1 retains `^spacedock \S+$`. Output retains three lines outside a session and six lines inside a session.
`contract.ParseMajorMinor` must return the same major and minor for the manifest, clean dev version, and dirty dev version.
The contract parser and linker-stamped release behavior stay unchanged.

### Dispatch identity and publication

Hash the exact final `dispatchBody` bytes with SHA-256. Append `-` plus the first eight lowercase hexadecimal digits to the existing filename stem.
The order is `[parent-]<derivedName>[-advance]-<hash8>.md`. All four hosts use this rule: Claude, Codex, Pi, and Paseo.
Keep the existing parent token rules and worker name. Keep schema version 2 and the existing envelope fields.
The hash excludes the output directory and the pointer prompt. Equal body bytes retain the same path within the same directory and parent context.

`dispatchFileNameMaxLen = 251` MUST stay enforced after the nine-character suffix.
Do not truncate the stem. An old stem of 242 bytes succeeds at 251 bytes. An old stem of 243 bytes fails at 252 bytes.
Keep the existing exit-1 filename diagnostic and emit no JSON or file on this failure. The `.md` extension adds three bytes.

Eight hexadecimal digits do not guarantee mathematical uniqueness. A collision must produce an error, never a replacement assignment.
Publish complete files with a temporary file in the destination directory, followed by `os.Link(temp, destination)`.
Write all bytes, apply mode `0644`, and close the temporary file before publication. Delete the temporary name on every return path.
A successful link publishes the complete bytes without replacing a destination. On `EEXIST`, require a regular destination with identical bytes.
An identical retry succeeds without rewriting the existing file. Different bytes, symlinks, directories, and unreadable destinations fail closed.
Keep the `dispatch_file_write_failed: <path>: <reason>` error prefix, exit 1, and empty stdout on publication failures.
A filesystem without hard-link support returns this error. Do not fall back to a truncating write or replacing rename.

### Dispatch directory

Add `DispatchDir string` and a flag-presence bit to `buildOptions`. Accept both `--dispatch-dir DIR` and `--dispatch-dir=DIR`.
Select the explicit flag first, then `os.LookupEnv("SPACEDOCK_DISPATCH_DIR")`, then `/tmp/spacedock-dispatch`.
Presence matters: an explicitly empty flag or environment value is an error, not permission to use the next source.
A valid explicit flag overrides an invalid environment value. Do not inspect unused lower-precedence values.

Reject empty values, whitespace-only values, NUL, CR, and LF. Preserve other bytes, including spaces within a path.
Resolve the selected value with `filepath.Abs`, relative to the caller's current directory. This also cleans `.` and `..` components.
Do not expand `~` or environment references. Do not canonicalize symlinks or interpret relative paths against `--workflow-dir`.
Create the directory with `os.MkdirAll(dir, 0755)`. Emit the resolved absolute spelling in both `dispatch_file_path` and the pointer prompt.

Missing flag operands and invalid selected values exit 2 with empty stdout. The invalid-value diagnostic is `error: invalid dispatch directory %q`.
An absolute-path resolution error or directory creation failure exits 1 with empty stdout and the underlying error.
Resolve, validate, and create the directory before `runStamp`. Invalid directories must not change entity state or Git history.
Later publication failures retain the existing post-stamp failure boundary. This change does not add stamp rollback.

In `parity_harness_test.go`, `TestMain` clears `SPACEDOCK_DISPATCH_DIR` alongside `CLAUDE_CODE_SESSION_ID`.
Directory tests use dedicated temporary directories. Concurrent calls use explicit hosts and a fixed environment, not helpers that mutate process environment.

### Golden stability

The raw body includes temporary fixture paths. Its hash changes between test runs even after existing path normalization.
Extend `normRun` in `golden_harness_test.go` to replace only a terminal `-[0-9a-f]{8}.md` suffix in dispatch locators with `-<HASH8>.md`.
Apply the same substitution to the pointer path. Do not normalize body text or remove the suffix entirely.
Extend `readDispatchBody` to compare the raw pathname suffix with an independent SHA-256 calculation over the file bytes before normalization.
A wrong hash must fail even when normalized goldens match. Production helper output must not supply this expected hash.

Update the six existing test files that assume unhashed paths. Negative cases must inspect their dedicated output directory for all artifacts.
A negative assertion against an obsolete unhashed filename can pass while a new hashed file escapes detection.

### Proposed documentation changes

These diffs use Layer 2 commit `5210d36d0`, not the local main copy.
`first-officer-shared-core.md` MUST stay untouched: Layer 2 is `23,464` B, within its `23,500` B cap.
The inspected local copy is `23,497` B. That difference belongs to another change and is not part of this task.

Add exactly three operational bullet lines before the existing frozen-label bullet in `skills/first-officer/references/paseo-fo-dispatch.md`:

```diff
+- **Shared dispatch directory:** If workers use a private `/tmp`, pass `--dispatch-dir DIR` or set `SPACEDOCK_DISPATCH_DIR`. Use a directory that both processes can read. The flag takes precedence. The default is `/tmp/spacedock-dispatch`.
+- **Workspace location:** Paseo creates worktrees under `~/.paseo/worktrees/` with `--isolation worktree`. Spacedock owns the task worktree. Adopt `code_cwd` with `--isolation local`.
+- **Thinking setting:** After spawn, run `paseo inspect "$agent_id" --json`. Verify `Thinking` against the requested setting. The `paseo models` metadata alone does not prove the worker setting.
 - **Frozen shared label contract:**
```

The command reference lists `spacedock dispatch`, but has no `dispatch build` section. Add a compact paragraph there and extend its version section:

```diff
 `spacedock --version` reports the version, platform, and release channel. Inside a supported or Paseo-hosted agent session, it also reports the runtime, sandbox, and contract level.
+
+An unstamped source build includes a seven-character commit ID when Go provides VCS information: `0.29.0-pre0+dev.5210d36`.
+A modified checkout adds `.dirty`. Without VCS information, the version ends in `+dev`. Linker-stamped versions stay unchanged.
@@
 | `spacedock completion` | Generate shell completions |

+`spacedock dispatch build --dispatch-dir DIR` selects the assignment directory. The flag takes precedence over `SPACEDOCK_DISPATCH_DIR`.
+The default is `/tmp/spacedock-dispatch`. Relative paths resolve from the current directory. Empty or invalid values cause an error.
+Assignment filenames include an eight-character content hash. Equal assignments reuse the file. Conflicting content causes an error without replacement.
+If workers use a private `/tmp`, select a directory that both processes can read.
+
 `spacedock gate read <ref-or-path> --json` returns the stored `gates` object, including every attempt and its application state.
```

Add this help line in `internal/dispatch/dispatch.go` immediately after the existing workflow-directory option:

```diff
+  --dispatch-dir DIR           Assignment directory. Overrides SPACEDOCK_DISPATCH_DIR. Default: /tmp/spacedock-dispatch. Relative to the current directory.
```

## Risk evidence

No spike needed: the inspected code already assembles `dispatchBody`, publishes a pointer, reads embedded versions, and accepts value flags.
The change uses standard Go VCS settings, SHA-256, path resolution, and atomic hard-link publication. No new host protocol is necessary.
The first implementation tests must exercise publication conflicts and overlapping calls before production changes.

Source inspection on 2026-11-23 found the same `build.go`, `dispatch.go`, `dev_version.go`, and parity harness bytes at local HEAD and Layer 2.
`version_session_test.go` differs between those revisions. Implementation must extend the Layer 2 version of its existing render tests.
The existing source-build test expects exactly `+dev`. It must derive its expected revision independently from the built binary's VCS information.
The golden harness currently normalizes root and home paths only. Without hash normalization, all successful build goldens become nondeterministic.
The state history records Layer 2 authorization and this ideation dispatch. These notes also serve as the session journal.

## Out of scope

State-checkout initialization, Paseo model metadata repairs, host registration, worker reuse, new launchers, and release operations are outside this task.
Do not edit the shared first-officer core, ensign instructions, workflow schema, worker labels, or the Layer 2 lifecycle rules.
Do not add locks, a registry, hash-collision retries, file cleanup policy, or a new dependency.
The selected directory must support hard links and be readable in both process namespaces. Directory selection does not mount or share storage.

## Expected surface and tolerance

Estimate net LOC change: **+574, across 43 code, test, fixture, and documentation files**.
Estimated insertions: 691. Estimated deletions: 117. Tolerance: net +400 to +750 and at most 45 files.
The state entity is one additional file in the state checkout. This ideation changes only that file.

| Files | Count | Insertions | Deletions | Net |
| --- | ---: | ---: | ---: | ---: |
| `internal/cli/{dev_version.go,dev_version_test.go,version_session_test.go}` | 3 | 100 | 15 | +85 |
| `internal/dispatch/{build.go,dispatch.go}` | 2 | 120 | 20 | +100 |
| `internal/dispatch/{parity_harness_test.go,golden_harness_test.go}` | 2 | 40 | 0 | +40 |
| New `internal/dispatch/build_artifact_test.go` and six existing test files listed next | 7 | 220 | 25 | +195 |
| `skills/integration/dispatch_test.go` | 1 | 45 | 5 | +40 |
| `skills/first-officer/references/paseo-fo-dispatch.md`, `docs/site/reference/command-reference.md` | 2 | 14 | 0 | +14 |
| Existing golden files listed next | 26 | 52 | 52 | 0 |

The six existing test files are `build_advance_test.go`, `build_errors_test.go`, `build_paseo_host_test.go`, `build_codex_host_test.go`, `build_input_mode_test.go`, and `self_contained_assignment_test.go`.
All six are under `internal/dispatch/`.

Exactly **26 of 79** golden files contain dispatch paths. Each needs two line replacements: `dispatch_file_path` and `prompt`.
The count matches both the inspected checkout and Layer 2. Paths below are relative to `internal/dispatch/testdata/golden/`:

```text
build-abs-worktree.txt
build-advance-codex-host.txt
build-advance-feedback-reflow.txt
build-advance-plain.txt
build-advance-split-root.txt
build-crossproduct-feedback+scope+reflow.txt
build-crossproduct-single+flat+nonworktree+bare.txt
build-crossproduct-single+flat+nonworktree+team.txt
build-crossproduct-single+flat+worktree+team.txt
build-crossproduct-single+folder+worktree+team.txt
build-crossproduct-split+flat+worktree+team.txt
build-crossproduct-split+folder+nonworktree+team.txt
build-crossproduct-split+folder+worktree+team.txt
build-host-codex-model-ignored.txt
build-host-pi-model-ignored.txt
build-html-escape.txt
build-model-defaults-fable.txt
build-model-defaults-haiku.txt
build-model-null.txt
build-model-stage-fable.txt
build-model-stage-wins-opus.txt
build-mods.txt
build-nonascii-title.txt
build-schemaversion-float-2.0.txt
build-schemaversion-int-2.txt
build-space-path.txt
```

Permitted semantic changes: dev version metadata, dispatch filename suffixes, immutable publication, the directory flag and environment variable, and the three operational notes.
JSON keys, worker identity, stage authority, and release versions retain their existing semantics. An additional semantic change requires a design revision.

## Mechanism-to-value justification

| Mechanism | Value criterion | Simplest alternative | Why the alternative is insufficient |
| --- | --- | --- | --- |
| VCS read seam | 1 | Keep `+dev` alone | Different installed commits remain indistinguishable. The seam provides deterministic missing and dirty cases. |
| Eight-digit body hash | 2 | Keep the parent prefix alone | One parent can issue different briefs for the same entity and stage. |
| Atomic link and byte comparison | 2 | Write directly to the hashed path | A retry can expose partial bytes. A short-hash collision can replace a different brief. |
| Directory flag and environment fallback | 3 | Keep the fixed `/tmp` path | Parent and worker namespaces can resolve that spelling to different storage. |
| Hash normalization and raw digest assertion | 2 | Record literal hashes in goldens | Temporary fixture paths change the hash on each run. Normalization alone can hide an incorrect digest. |
| Three adapter notes and one integration smoke extension | 3, 5 | Leave existing operational instructions | Operators lack the directory override and a direct check of the actual thinking setting. |

## Acceptance criteria

**AC-1 (VALUE) — Operators can distinguish two valid source revisions from `--version` alone.** Clean and dirty builds retain compatibility and the existing output shape.
   Verified by: injected revisions `5210d36d0` and `abcdef012` produce distinct literal tokens. A built-binary fixture compares the token with `debug/buildinfo.ReadFile` settings.
   Missing, short, malformed, and dirty settings cover fallback behavior. Stamped `0.29.0` stays exact and does not call the seam.
   Existing render tests require `^spacedock \S+$` and three/six lines. Literal major/minor assertions exercise `contract.ParseMajorMinor`.
   Falsifying edits: omit the revision, append a separate SHA line, change release output, or let metadata alter compatibility.

**AC-2 (VALUE) — Two overlapping briefs from one Claude or Paseo parent preserve two complete assignments.** Identical retries preserve the original file without replacement.
   Verified by: concurrent builds with independent checklist sentinels return distinct paths. Reading both paths after both calls returns each original sentinel and complete body.
   The same-body retry retains the path and file identity. A conflicting preexisting destination returns exit 1, empty stdout, and unchanged bytes.
   Falsifying edits: omit the content suffix, restore `os.WriteFile` publication, or accept conflicting bytes as an identical retry.

**AC-3 (VALUE) — A worker can read the complete assignment from the selected shared directory.** Selection is flag, then environment, then default.
   Verified by: CLI fixtures cover both flag forms, environment fallback, unset default, relative paths, spaces, and flag precedence over an invalid environment.
   Empty, whitespace-only, NUL, newline, missing-operand, and non-directory inputs return the specified failures with no JSON.
   Invalid-directory `--stamp` cases retain entity bytes and Git HEAD. A fresh reader process opens the exact emitted absolute path and compares checklist bytes.
   Falsifying edits: ignore the override, use the workflow directory as the relative base, fall back after an empty value, or stamp before rejection.

**AC-4 — All host paths retain the content suffix and filename limit.** Fresh, bare where supported, and advance modes retain their envelope contracts.
   Verified by: raw digest assertions cover Claude, Codex, Pi, and Paseo. Direct boundary cases accept a 242-byte old stem and reject 243 bytes.
   Goldens retain `-<HASH8>.md` in both locator fields. `TestMain` environment isolation keeps a deliberately set external dispatch directory from changing fixtures.
   Falsifying edits: hash only one host, measure length before the suffix, change a worker name, or omit environment clearing.

**AC-5 — The operational instructions expose shared storage and actual Paseo settings without expanding the shared core.**
   Verified by: a one-off Paseo smoke adopts the existing `code_cwd`, uses a shared dispatch directory, and records `inspect.Cwd` and `inspect.Thinking`.
   The worker reads the emitted pointer and commits its assigned result. Compare shared-core bytes with Layer 2 and require exactly 23,464 bytes.
   Falsifying changes: create a second worktree, use an unreadable dispatch directory, accept an unexpected thinking value, or change the shared core.

## Test plan

1. Add focused tests before each production change. `dev_version_test.go` owns the version branches and the real source-build fixture.
   `version_session_test.go` owns output shape. Reuse those owners instead of a second version harness. Cost: seconds for units, one Go build for the fixture.
2. `build_artifact_test.go` owns concurrent publication, conflicting destinations, deterministic retry, directory selection, and filename boundaries.
   Pin environment before concurrency. Use `RunWithLauncher` with explicit hosts and separate buffers. Do not call environment-mutating helpers from concurrent goroutines.
   Use a barrier to start Claude and Paseo pairs. Inspect files after both returns. No daemon is necessary for these deterministic CLI claims.
   Cost: small Git fixtures and milliseconds of file I/O per case. An explicit wrong destination proves conflict handling without a brute-force hash search.
3. Extend the existing advance and error owners. Replace obsolete negative filename assertions with dedicated-directory snapshots in the six identified files.
   `readDispatchBody` owns raw hash verification. The golden harness owns stable snapshots only. Regenerate the 26 files, then rerun without `-update`.
   Distinct falsifiers are a wrong digest, a missing suffix, a leaked artifact on failure, and a wrong advance envelope.
4. Extend `skills/integration/dispatch_test.go` before the adapter edit. Drive a Paseo dispatch into a shared directory and read it from a fresh process.
   This smoke proves pointer accessibility and split-root state paths. A wrong directory or rewritten state path fails it. Do not assert prose substrings.
5. During validation, run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` in the implementation worktree.
   Ideation performs no code formatting or product test runs because its authorized write scope is the state document only.
6. Before merge, run the required live lanes for every host touched by the dispatch core, including Paseo. Apply `docs/runtime-support.md` for setup friction.
   Record durable worker output and Git state. Use the existing runtime substrate for the one-off workspace and thinking observations in criterion 5.
   Live cost depends on host availability and model latency. A deterministic fixture does not replace these runtime proofs.
7. The shipped adapter change requires the workflow's detached adversarial audit before merge. Use a throwaway checkout and the stated falsifying edits.
   No live smoke, code test, or audit result is claimed by this ideation report.

### Feedback Cycles

- Cycle 1: REJECTED — detached adversarial audit (`0f5e8a8c`) + validation (`7dd87755`); 1 Material test-strength hole at `internal/dispatch/build_artifact_test.go:153-157` (relative `--dispatch-dir` subtest collapsed `cwd` and `--workflow-dir` bases) + Commander-owned linked-worktree binary provenance (resolved by standalone-clone build); surface 44/+489 vs estimate 43/+574 (-14.8%); AC unchanged

## Stage Report: ideation

- DONE: Flesh out concrete design in docs/dev/.spacedock-state/paseo-vm-dispatch-and-version-frictions/index.md for Friction #1 (dev_version.go VCS SHA via injectable debug.ReadBuildInfo seam preserving ^spacedock \S+$ and 3/6-line output), Friction #4 (deterministic 8-hex SHA-256 content suffix -<hash8> on dispatchBody across all hosts respecting 251-char cap + golden updates + concurrent-brief tests on Claude & Paseo), Friction #5 (--dispatch-dir flag > SPACEDOCK_DISPATCH_DIR env > /tmp/spacedock-dispatch resolved absolute, cleared in parity_harness_test.go TestMain), and exact before/after 3-line operational notes in paseo-fo-dispatch.md (checked against origin/paseo-and-generic-harness-fo-support 5210d36d0; plus command-reference.md if applicable)
  The five criteria define the design and falsifiers. This is design coverage, not evidence that the future implementation passes.
- DONE: Verify no spike needed (or inspect existing dev_version.go, build.go, dispatch.go, version_session_test.go, parity_harness_test.go, and golden fixtures directly), declare net LOC Expected surface & tolerance, Out of scope, and Mechanism-to-value justification
  Inspected the named source owners and both golden inventories. Exactly 26 of 79 goldens change. Estimated net +574 LOC across 43 files.
- DONE: Write scanner-exact ## Acceptance criteria (with at least one VALUE AC and falsifying Verified by clauses) and ## Test plan, run /tmp/spacedock-main status --validate / --ac-scan / --checklist, and commit+push path-scoped ## Stage Report: ideation in docs/dev/.spacedock-state
  Validation returned VALID. The scanners returned five criteria and three DONE items. Product evidence remains pending at ideation.

### Summary

The design specifies VCS version metadata, immutable dispatch publication, directory selection, and exact operational doc additions.
It preserves the filename cap and records the golden normalization requirement. Validation reported one unrelated warning for live-lane-gap-inventory.
The first scanner calls required this stage report. Both scanners succeeded after report creation. Product tests and live runtime proofs remain implementation and validation work.

## Implementation work log

1. Complete: dispatch publication and directory fixtures plus production behavior. Concurrent Claude/Paseo briefs, immutable retries, conflict rejection, pre-stamp refusal, and all-host digest assertions pass.
2. Complete: VCS version fixtures/implementation and approved operational documentation. Complete: full validation and code commit `1fdb1e42a`; state sync accompanies this report.
3. Runtime observation for AC-5 is assigned to Commander under the F8 ensign sandbox boundary. No journal tool is available; this log records implementation decisions and evidence.

4. Investigation notes: parent tokens are capped at 36 bytes, so filename limits use the approved direct boundary test. OS environments reject NUL before Go can read it; NUL is covered through both flag forms. One additional Pi pathname assertion required migration, within the 45-file cap.

5. Validation: all packages passed normal tests with the required CLI/default PATH split; full race validation also passed. Vet, explicit contractlint/integration race checks, strict docs build, external-directory isolation, and workflow validation passed. The workflow validator retains the unrelated live-lane-gap-inventory warning.
6. Docs tooling required an explicit Python interpreter to avoid ambient UV virtualenv discovery; the isolated pinned-dependency build passed. A 90-second diagnostic CLI run timed out during existing state fixtures; the standard-timeout normal and race CLI runs passed (253s/267s).

## Stage Report: implementation

- DONE: Implement Friction #1 in internal/cli/dev_version.go (+ dev_version_test.go, version_session_test.go): injectable debug.ReadBuildInfo seam appending .<sha7>[.dirty] inside +dev build metadata when Version=="dev", preserving ^spacedock \S+$, 3/6-line output, contract.ParseMajorMinor, and linker-stamped releases (AC-1)
  AC-1: `TestDevVersionVCSMetadata`, `TestUnstampedSourceBuildReportsCheckoutVersionPlusDev`, and `TestVersionSessionRender` prove revision/dirty tokens, independent binary VCS metadata, compatibility, and 3/6-line output; removing SHA, adding a line, or changing stamped output makes them fail.
- DONE: Implement Frictions #4 & #5 in internal/dispatch/build.go and internal/dispatch/dispatch.go (+ build_artifact_test.go, parity_harness_test.go TestMain, golden_harness_test.go raw-digest check + -<HASH8>.md norm, 6 updated test files, and 26 updated golden fixtures): deterministic 8-hex SHA-256 dispatchBody suffix across all 4 hosts enforcing 251-char cap after suffix, atomic temp+os.Link immutable publication with fail-closed EEXIST byte check, and --dispatch-dir > SPACEDOCK_DISPATCH_DIR > /tmp/spacedock-dispatch validated+created before runStamp (AC-2, AC-3, AC-4)
  AC-2/3/4: `build_artifact_test.go` exercises concurrent briefs, immutable retries/conflicts, directory selection, pre-stamp failures, every host/mode, and 242/243-byte boundaries; raw-body SHA assertions precede golden normalization.
- DONE: Apply exact 3-line operational notes in skills/first-officer/references/paseo-fo-dispatch.md (keeping first-officer-shared-core.md strictly untouched at 23,464 B), update docs/site/reference/command-reference.md + skills/integration/dispatch_test.go, run gofmt -w ./cmd ./internal, go test ./... + go test ./... -race, contractlint, and spacedock status --validate, commit cleanly on paseo-vm-dispatch-and-version-frictions atop 5210d36d0, and commit+push path-scoped ## Stage Report: implementation (AC-5)
  AC-5: exact three adapter bullets and command reference committed in `1fdb1e42a`; shared core is byte-identical to `5210d36d0` at 23,464 bytes. Code surface is 44 files, +489 net LOC (630 additions/141 deletions), within tolerance.
- DONE: AC-2 preserves both overlapping assignments and identical retry identity.
  `TestArtifactOverlappingBriefsAndRetry` and `TestArtifactPublicationConflict` inspect complete sentinel bodies, inode/mtime, conflicting bytes, symlinks, and directories; dropping the suffix, restoring truncating writes, or accepting conflicts makes them fail.
- DONE: AC-3 proves shared-directory selection and rejection before stamping.
  `TestArtifactDirectorySelection`, `TestArtifactInvalidDirectoryBeforeStamp`, and `TestSplitRootFolderWorktreeDispatch` compare selected paths, entity bytes/Git HEAD, and a fresh reader's complete body; ignoring precedence, falling back from empty values, or stamping before rejection makes them fail.
- DONE: AC-4 retains content hashes across all hosts and the filename cap.
  `TestArtifactHostsModesAndFilenameBoundary` plus `readDispatchBody` compare raw SHA-256 and exact limits; a wrong digest or measuring before the suffix fails. Dispatch suite passed with an invalid external `SPACEDOCK_DISPATCH_DIR`; all 26 locator goldens retain `-<HASH8>.md`.
- DONE: Required validation and documentation build.
  All packages passed normal and race tests with CLI PATH `/opt/homebrew/bin:/usr/bin:/bin` and default PATH elsewhere; `go vet ./...`, uncached race contractlint/integration, strict MkDocs, and workflow validation passed. `gofmt -w ./cmd ./internal` ran; touched Go files and `git diff --check` are clean.
- SKIPPED: Live Paseo runtime observation in AC-5.
  F8 ensign sandbox boundary: Commander will directly record inspect.Cwd, inspect.Thinking, pointer readability, and the worker's result commit. Fixture accessibility is proven; live observations are not claimed.

### Summary

Commit `1fdb1e42a` adds VCS dev metadata, immutable content-hashed assignments, configurable shared storage, and the approved operational notes. Workflow validation is VALID with the existing unrelated live-lane-gap-inventory warning; the code worktree is clean and the state report is committed path-scoped. Commander owns the AC-5 live observation; live host lanes and the detached adversarial audit remain pre-merge validation gates.

## Validation work log

1. Candidate HEAD is `1fdb1e42ade7f42e026a40c2d2cfd70ce032cb2e`. The worktree remains clean. The diff has 44 files and +489 net lines.
2. The shared first-officer core matches `5210d36d0` byte-for-byte at 23,464 bytes. All changed Go files pass `gofmt -l`.
3. Repository-wide `gofmt -l` lists `internal/release/runtime_live_evidence_workflow_test.go`. The same formatting difference exists at baseline `5210d36d0`.
4. The worker read the shared assignment pointer and independently computed SHA-256 prefix `608ad412`. Commander reports adopted workspace `wks_508514ca958b5a19` and `inspect.Thinking == "high"`.
5. The current `/tmp/spacedock-layer3` reports `0.29.0-pre0+dev.3d134b2.dirty`. `go version -m` confirms revision `3d134b255fbcc493bd5e5ccc77a16c496c8b3662` and `vcs.modified=true`. This conflicts with the stated candidate provenance. The worker observes `Runtime: codex (CODEX_THREAD_ID)` and contract 3. Commander reports Paseo in the parent environment.
6. `go vet ./...`, `git diff --check 5210d36d0..HEAD`, and workflow validation pass. The existing `live-lane-gap-inventory` warning remains. Normal and race suites are in progress.

No journal tool is available. This log retains the validation evidence.

## Review-finding disposition

1. Live candidate provenance remains unresolved.
   - Observation: Commander identifies `/tmp/spacedock-layer3` as a build from `1fdb1e42a`. The worker observes `+dev.3d134b2.dirty`. Embedded metadata independently confirms `vcs.revision=3d134b255fbcc493bd5e5ccc77a16c496c8b3662` and `vcs.modified=true`.
   - Released user and normal workflow: operators dispatch a worker with the candidate binary and use the resulting commit as runtime evidence.
   - Observable harm: the successful pointer read cannot establish which candidate produced that pointer. The binary can contain unrecorded source changes.
   - Authority: value-ac[AC-5] The live runtime evidence must establish the candidate's shared-directory handoff and adopted workspace behavior.
   - Trigger evidence: the binary metadata disagrees with the supplied candidate revision. The worker directly read pointer `608ad412` and independently verified its digest.
   - Worker proposal: evidence defect, Material. Commander owns the runtime evidence. Hold the gate until Commander reconciles provenance or repeats the observation with a uniquely named candidate binary.
   - Proposed narrow correction: retain the binary digest and embedded VCS metadata, `inspect.Cwd`, `inspect.Thinking`, emitted pointer, and worker result commit. No product change or new harness is proposed. FO authorization remains pending.
2. Baseline formatting remains outside this candidate.
   - Observation: `gofmt -l ./cmd ./internal` lists `internal/release/runtime_live_evidence_workflow_test.go`. Its two field-alignment differences also exist at `5210d36d0`.
   - Released user and normal workflow: contributors run the repository formatting command. Observable harm: that check remains nonempty.
   - Authority: none: the unchanged whitespace does not violate a value AC or a protected runtime boundary.
   - Trigger evidence: formatting the baseline blob produces the same diff. Every changed Go file passes the formatting check.
   - Worker proposal: Polish, pre-existing and outside this task. No candidate edit is proposed. FO disposition remains pending.

## Stage Report: validation

Recommendation: REJECTED

- DONE: Verify commit 1fdb1e42a atop 5210d36d0 matches approved ideation design (44 files <= 45, +489 net LOC within +400..+750, first-officer-shared-core.md byte-identical at 23,464 B, exact 3-line operational notes in paseo-fo-dispatch.md, command-reference.md + dispatch.go help text in STE)
  The candidate has 630 additions and 141 deletions. Independent baseline comparison confirms the shared-core bytes and exact three approved bullets. Changed documentation passes the structural STE review.
- DONE: Independently reproduce AC-1 through AC-5 tests: gofmt -l ./cmd ./internal, go vet ./..., PATH=/opt/homebrew/bin:/usr/bin:/bin go test -race -count=1 ./internal/cli, go test -race -count=1 ./internal/dispatch ./internal/contractlint ./skills/integration, git diff --check 5210d36d0..HEAD, and /tmp/spacedock-layer3 status --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --validate; perform semantic adversarial pass over VCS version parsing, atomic os.Link conflict handling, --dispatch-dir precedence & pre-stamp validation, and raw SHA-256 golden verification
  Every listed command ran. Vet, race suites, diff checks, and workflow validation pass. Repository formatting lists one unchanged baseline file. AC-5 remains unresolved as recorded next.
- DONE: Record ## Stage Report: validation (Recommendation: PASSED or REJECTED) citing AC-1 through AC-5 with falsifying edits, Review-finding disposition, and live dogfood evidence (you are reading a content-hashed pointer emitted by /tmp/spacedock-layer3 --dispatch-dir /tmp/spacedock-layer3-shared-dispatch) alongside F8 sandbox boundary
  This report and the preceding disposition section retain the evidence, scope, proposed classification, and Commander handoff.
- DONE: AC-1 — revision tokens, release compatibility, and output shape.
  `TestDevVersionVCSMetadata`, `TestUnstampedSourceBuildReportsCheckoutVersionPlusDev`, and `TestVersionSessionRender` pass. Removing the SHA, changing stamped output, or adding a line violates their exact assertions.
- DONE: AC-2 — complete overlapping briefs and immutable retries.
  `TestArtifactOverlappingBriefsAndRetry` and `TestArtifactPublicationConflict` pass for Claude/Paseo sentinels, identity/mtime, conflicting bytes, symlinks, and directories. Removing the suffix or restoring truncating publication fails these assertions.
- DONE: AC-3 — selected shared storage and refusal before stamping.
  `TestArtifactDirectorySelection`, `TestArtifactInvalidDirectoryBeforeStamp`, and `TestSplitRootFolderWorktreeDispatch` pass. Ignoring precedence, changing the relative base, accepting empty input, or stamping before rejection fails path/state assertions.
- DONE: AC-4 — raw hashes, envelopes, and filename limits.
  `TestArtifactHostsModesAndFilenameBoundary` and the golden suite pass across four hosts. `readDispatchBody` independently hashes raw bytes. A wrong digest or a pre-suffix length check fails. The suite also passes with an invalid ambient `SPACEDOCK_DISPATCH_DIR`.
- FAILED: AC-5 — establish candidate provenance for the live runtime observation.
  Commander reports the adopted `code_cwd` workspace `wks_508514ca958b5a19`, `isolation: worktree`, `background: true`, `inspect.Thinking == "high"`, and parent `Runtime: paseo (PASEO_AGENT_ID)`, contract 3. The worker directly read the shared pointer ending `-608ad412.md` and independently reproduced its digest.
  The current binary reports `0.29.0-pre0+dev.3d134b2.dirty`. `go version -m` confirms `3d134b255fbcc493bd5e5ccc77a16c496c8b3662` with modified sources. This conflicts with the claimed `1fdb1e42a` provenance. Binary SHA-256: dcf7714216add7b35233ae2e1ca0d8aa760ad578fcb8eeb73ffb47b270210195.
  F8 ensign sandbox boundary: the worker has no Paseo daemon grant. Parent inspection remains Commander evidence. Worker output reports `Runtime: codex (CODEX_THREAD_ID)`, `Sandbox: agent-safehouse`, and contract 3. Changing shared-core bytes, losing pointer access, adopting a second worktree, or accepting unexpected thinking falsifies AC-5.
- DONE: Semantic adversarial pass over adjacent inputs and lifecycle boundaries.
  The inspected matrix covers missing/malformed/uppercase/dirty/stamped VCS values, both directory flag forms and environment precedence, invalid paths before stamp, conflicting destination types, retries, all host modes, and 242/243-byte stems. Existing fixtures exercise exact output, file identity, cleanup, entity bytes, and Git HEAD. Hash checks precede normalization. No product defect emerged.
- DONE: Full normal and race package coverage with the prescribed PATH split.
  All packages pass uncached normal and race runs. CLI uses `/opt/homebrew/bin:/usr/bin:/bin`; other packages use the default PATH. CLI race took 434.669s. Dispatch race took 144.976s. The final race package group also passed.
- SKIPPED: Mutating formatting and candidate edits.
  The assignment requires read-only validation. `gofmt -l` replaces `gofmt -w`. Every changed Go file is formatted. The baseline formatting issue remains a proposed Polish finding.
- SKIPPED: Independent parent inspection, required host live lanes, and detached mutation audit.
  F8 reserves parent inspection for Commander. Host live lanes and the throwaway-checkout adversarial audit remain required before merge. This read-only validation does not claim those results or executed falsifying edits.

### Summary

Candidate tests pass, but unresolved live-binary provenance prevents acceptance of AC-5. The proposed Material finding is an evidence defect owned by Commander, with no product fix or design reset proposed. Code remains unchanged at `1fdb1e42a`; this path-scoped state report is the worker result for the shared-pointer handoff.

## Implementation work log (cycle 2)

1. Complete: Go 1.22-compatible fixture committed as `03a796c48` after explicit FO authorization. Targeted normal test, vet, formatting, mutation verification, and full dispatch race passed; all packages have passing normal/race coverage after the documented Pi runtime PATH adjustment.
2. Root cause: the old relative path climbs above both candidate bases, collapsing to the same absolute temp directory. The corrected fixture uses a child path under a distinct process cwd and derives its expected spelling from `os.Getwd()` for macOS symlink parity.
3. Scope: FO authorized only `internal/dispatch/build_artifact_test.go`; Commander owns standalone-clone build provenance. Existing formatting and post-stamp boundaries remain declined. No journal tool is available; this log retains the evidence.
4. FO authorized `os.Chdir(cwdDir)` with checked `t.Cleanup` restoration in place of `t.Chdir`. The resulting wrong-base overlay fails with `wrong locator` at line 178, while unmodified production passes; `build.go` stayed byte-identical and overlay files were removed.

## Stage Report: implementation (cycle 2)

- DONE: Rewrite `TestArtifactDirectorySelection/relative` to use a distinct temporary cwd, a relative child path, and an expected directory derived from `os.Getwd()` after changing cwd.
  Commit `03a796c48` uses FO-authorized `os.Chdir(cwdDir)` with checked `t.Cleanup` restoration for Go 1.22 compatibility. The caller directory is separate from `--workflow-dir`; the argument is `filepath.Join("shared assignments", "tmp", "..", "nested")`.
- DONE: Prove via throwaway Go `-overlay` mutation that resolving relative `--dispatch-dir` against `opts.WorkflowDir` fails `TestArtifactDirectorySelection/relative`, while unmodified `build.go` passes.
  The overlay inserts `if !filepath.IsAbs(dir) { dir = filepath.Join(opts.WorkflowDir, dir) }` before `filepath.Abs(dir)`. `go test -count=1 -overlay <overlay.json> ./internal/dispatch -run '^TestArtifactDirectorySelection/relative$' -v` fails with `wrong locator` at line 178 (exit 1); the same command without the overlay passes (exit 0). Production bytes remained unchanged; overlay files were removed.
- DONE: Run `gofmt -l ./internal/dispatch/build_artifact_test.go`, `go vet ./internal/dispatch`, and `go test -race -count=1 ./internal/dispatch`.
  Formatting output is empty and vet passes. Uncached artifact race tests pass (7.196s), and the complete uncached dispatch race suite passes (47.466s). The relative test checks the emitted absolute locator and complete assignment body; resolving against the workflow directory falsifies its locator assertion.
- DONE: Verify the baseline diff stays at 44 files and `skills/first-officer/references/first-officer-shared-core.md` stays byte-identical at 23,464 B; commit only the authorized test file on `paseo-vm-dispatch-and-version-frictions`.
  Post-commit checks against `5210d36d0` confirm both boundaries. Commit `03a796c48` changes one file (+13/-3); the code worktree is clean and `git diff --check` passes.
- DONE: Replace the Cycle 2 blocker report with the completed implementation report and commit/push the entity path in the split-root state checkout.
  This report replaces the blocked report per explicit FO instruction; the earlier blocker and authorization context remain in Git history and the implementation work log.
- DONE: Run repository-wide normal and race tests.
  `PATH=/opt/homebrew/bin:/usr/bin:/bin go test ./...` and `go test ./... -race` passed every package except `internal/piruntime`, where Homebrew Node could not load `libllhttp.9.3.dylib`. Rerunning `go test ./internal/piruntime` and `go test -race ./internal/piruntime` under the default PATH passed (0.286s/1.294s), giving passing normal and race coverage for every package.
- SKIPPED: Repository-wide formatting writes and Commander-owned binary provenance work.
  The authorized scope permits only the test-file correction and explicitly declines baseline formatting changes. Commander owns rebuilding the candidate from a standalone clone; this report makes no new AC-5 provenance claim.

### Summary

Commit `03a796c48` makes the relative dispatch-directory fixture distinguish the caller cwd from the workflow directory, with checked cleanup and Go 1.22-compatible APIs. The wrong-base mutation now fails, all required checks pass, and normal/race package coverage passes with the documented PATH adjustment. The cumulative surface remains 44 files with the first-officer shared core unchanged.


## Validation work log (cycle 2)

1. Candidate `03a796c4809c7e926040be9e34b17354ea789994` is clean. Cycle 2 changes only `internal/dispatch/build_artifact_test.go` (+13/-3); cumulative surface is 44 files, +640/-141 (net +499), within the approved 45-file and +400..+750 limits. Shared-core bytes match `5210d36d0` exactly at 23,464 B.
2. Independently reproduced the authorized wrong-base mutation using a temporary Go overlay: inserting workflow-relative resolution before `filepath.Abs` makes `TestArtifactDirectorySelection/relative` fail at line 178 with `wrong locator` (exit 1); unmodified code passes (exit 0). Overlay files were removed and production bytes stayed unchanged.
3. `go vet ./...`, changed-test formatting, diff checks, workflow validation, and focused AC-1 CLI tests pass. Targeted race suites pass: dispatch 38.035s, contractlint 3.528s, integration 6.508s. The directory-selection suite also passes with an invalid ambient override. Workflow validation retains only the existing unrelated `live-lane-gap-inventory` warning.
4. Standalone `/tmp/spacedock-layer3-clean` reports `0.29.0-pre0+dev.03a796c`; embedded revision is `03a796c4809c7e926040be9e34b17354ea789994`, `vcs.modified=false`, SHA-256 `3cb990664114cc9b9b6e797263ef46be4011817b546c4069c30ba23221df8890`. Assignment SHA-256 independently equals `6f21932667c3dd26570f2e4ddd394e69ae81435bb365b545d7112ad710af9d0d`, matching its filename suffix.
5. Commander records this pointer as emitted by that clean binary with `--host paseo --dispatch-dir /tmp/spacedock-layer3-shared-dispatch`, adopted workspace `wks_508514ca958b5a19`, `inspect.Cwd` at the assigned code worktree, `Worktree: null`, `inspect.Thinking == "high"`, and parent `Runtime: paseo (PASEO_AGENT_ID)`, contract 3. Worker directly observes the assigned cwd and pointer bytes; its inner runtime reports Codex and contract 3. Parent inspection is Commander evidence under the F8 boundary.

No journal tool is available; this work log retains the validation evidence and remaining checks.


## Review-finding disposition (cycle 2)

1. Material evidence defect, AC-3 relative-directory test: FO authorized the fixture-only correction in `03a796c48`. Independent overlay reproduction now fails on workflow-based resolution and passes on caller-cwd resolution. The detached audit's only Material finding is closed; no production change was needed.
2. Material evidence defect, AC-5 provenance: FO-owned correction supplies the standalone clean build and this fresh content-hashed dispatch. Independent binary version, embedded revision, clean flag, and digest match the assignment. Commander supplies parent inspection; this worker supplies pointer-read, cwd, and committed report evidence. The previous provenance blocker is resolved.
3. Deferred risk, accepted post-stamp boundary: an existing `0555` dispatch directory can pass `MkdirAll` and fail publication after stamping. The approved contract requires writable, hard-link-capable shared storage and explicitly retains post-stamp failure without rollback; supported writable-directory, immutable publication, and invalid-directory pre-stamp fixtures pass. Promote to Material if preflight rejection of unwritable directories or rollback becomes promised, or writable supported storage exhibits this failure. FO declined candidate work at this boundary.
4. Polish, pre-existing formatting: `internal/release/runtime_live_evidence_workflow_test.go` is unchanged from the baseline. FO declined unrelated formatting edits; the changed test is formatted. No new findings require disposition.

## Stage Report: validation (cycle 2)

Recommendation: PASSED

- DONE: Verify Cycle 2 commit `03a796c4809c7e926040be9e34b17354ea789994`, one-file correction, cumulative surface, and shared-core byte identity.
  `git diff --stat 1fdb1e42a..HEAD` contains only `internal/dispatch/build_artifact_test.go` (+13/-3); `5210d36d0..HEAD` is 44 files, +640/-141 (net +499), within tolerance. Shared core equals baseline byte-for-byte at 23,464 B.
- DONE: Re-verify the detached audit's ONLY Material finding directly with a throwaway Go `-overlay`, keeping the code worktree clean.
  Inserting workflow-relative resolution before `filepath.Abs(dir)` makes `TestArtifactDirectorySelection/relative` fail with `wrong locator` at line 178 (exit 1); unmodified code passes (exit 0). Overlay removed; production bytes unchanged.
- DONE: Run targeted formatting, vet, race, diff, and workflow verification.
  `gofmt -l ./internal/dispatch/build_artifact_test.go` is empty; `go vet ./...`, `go test -race -count=1 ./internal/dispatch ./internal/contractlint ./skills/integration` (38.035s/3.528s/6.508s), and `git diff --check 5210d36d0..HEAD` pass. `/tmp/spacedock-layer3-clean status --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --validate` is VALID with the existing unrelated `live-lane-gap-inventory` warning.
- DONE: Verify AC-5 live provenance cleanly and independently compute the assignment SHA-256 prefix.
  `/tmp/spacedock-layer3-clean --version` reports `0.29.0-pre0+dev.03a796c`; `go version -m` gives revision `03a796c4809c7e926040be9e34b17354ea789994`, modified=false; binary SHA-256 is `3cb990664114cc9b9b6e797263ef46be4011817b546c4069c30ba23221df8890`. Read pointer `3b905937-e237-4060-af8c-87b279dd96d9-paseo-vm-dispatch-and-version-frictions-validation-6f219326.md` from the selected shared directory; independently computed prefix is `6f219326`.
- DONE: AC-1 — Operators can distinguish two valid source revisions from `--version` alone.
  Focused uncached CLI tests pass: `TestDevVersionVCSMetadata`, `TestUnstampedSourceBuildReportsCheckoutVersionPlusDev`, `TestVersionSessionRender`, and stamped/fallback/contract-placement tests assert literal tokens, independently read build metadata, 3/6 lines, and compatibility; omitting SHA, changing stamped output, or adding a line falsifies them.
- DONE: AC-2 — Two overlapping briefs from one Claude or Paseo parent preserve two complete assignments.
  Race-tested `TestArtifactOverlappingBriefsAndRetry` and `TestArtifactPublicationConflict` verify distinct complete sentinels after both builds, retry inode/mtime, conflict bytes, symlinks/directories, mode, and cleanup; removing the content suffix, replacing publication with truncation, or accepting conflicting bytes falsifies these assertions.
- DONE: AC-3 — A worker can read the complete assignment from the selected shared directory.
  `TestArtifactDirectorySelection`, `TestArtifactInvalidDirectoryBeforeStamp`, and `TestSplitRootFolderWorktreeDispatch` pass flag/env/default precedence, both flag forms, spaces, invalid values, entity bytes/HEAD, and fresh-process body reads. The executed wrong-base mutation fails; ignoring overrides, falling back from empty input, or stamping before rejection also violates assertions.
- DONE: AC-4 — All host paths retain the content suffix and filename limit.
  `TestArtifactHostsModesAndFilenameBoundary` and goldens pass four-host fresh/advance/supported-bare cases, independently hashed raw bytes, and 242/243-byte stems. Directory selection passes with an invalid external `SPACEDOCK_DISPATCH_DIR`; wrong digests, pre-suffix length checks, or removed environment clearing falsify these checks.
- DONE: AC-5 — Operational instructions expose shared storage and actual Paseo settings without expanding the shared core.
  Commander records clean-binary emission, adopted `wks_508514ca958b5a19`, `inspect.Cwd` at this assigned code worktree, `Worktree: null`, `inspect.Thinking == "high"`, and parent Paseo/contract 3. Worker directly verifies cwd, pointer bytes/digest, shared-core identity, and supplies this committed result. A second worktree, unreadable pointer, wrong thinking, or changed core falsifies the claim; parent inspection remains attributed to Commander under F8.
- DONE: Semantic adversarial pass and FO-disposed findings.
  Rechecked version variants, path selection/rejection before stamp, immutable publication/conflict/cleanup, host envelopes, and raw digest-before-normalization invariants; the Cycle 2 mutation closes the exact evidence gap. No new Material finding remains. Accepted Deferred and baseline Polish findings are recorded above.
- SKIPPED: Repeating repository-wide normal/race tests and mutating formatting.
  Explicit Cycle 2 scope reuses passing full-package coverage on `03a796c` and Cycle 1 validation; unchanged CLI does not require its 435-second race suite again. Focused AC-1 tests were rerun. Read-only validation preserves the FO-declined baseline formatting.
- DONE: Record the Cycle 2 validation report and commit/push the assigned entity path in the split-root state checkout.
  Only this entity body is changed; frontmatter and candidate code remain unchanged. The state commit containing this report is the durable worker result for the clean-binary shared-pointer handoff.

### Summary

Candidate `03a796c48` satisfies AC-1 through AC-5 with the independent mutation check, passing targeted tests, and reconciled clean-build provenance. Both Material evidence defects are resolved; the accepted post-stamp Deferred risk and pre-existing formatting Polish remain as authorized. Validation recommends PASSED; Commander retains ownership of parent inspection and any separate pre-merge runtime gates.
