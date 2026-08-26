---
id: 9n02rsw1s4tztqzgmwb07n1k
title: gate prepare resolves operator-supplied artifact paths without doubling
status: ideation
source: "email-triage codex audit 2026-08-26: three of six gate-prepare attempts across two days failed with a doubled state path (.../.spacedock-state/docs/triage/.spacedock-state/...); under the no-retry rule the third failure left a batch's gate unprepared for the rest of the window"
started: 2026-08-26T21:19:07Z
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:9n02rsw1s4tztqzgmwb07n1k:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9n02rsw1s4tztqzgmwb07n1k-backlog-1
              briefing:
                id: briefing:9n02rsw1s4tztqzgmwb07n1k:backlog:attempt-1:revision-1
                digest: sha256:a1db71cb8071b29305e8b027b71b0c14747ac5c4c62eaee6631f73b1de70ede5
                room-ref: ./gate-prepare-accepts-operator-paths/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9n02rsw1s4tztqzgmwb07n1k:backlog:1
                briefing: briefing:9n02rsw1s4tztqzgmwb07n1k:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T21:18:44.298895Z"
                decision: approve
              application:
                target-stage: ideation
                state: consumed
---

`gate prepare --artifact` and `--reference` resolve supplied paths relative to the state entity directory. An operator who supplies a project-relative path (`docs/triage/.spacedock-state/x.md`) gets it re-prefixed under the state root and the command fails: "selected source must be a readable non-symlink regular file" with the doubled path. Three of six live preparations hit this, and the gate skill's cwd-path wording is ambiguous against the binary's resolution rule.

## Problem

`gate prepare` receives selected-source paths at the CLI boundary but resolves every non-absolute spelling inside `internal/gates` from the entity root. In a split-root workflow, an operator launching from the repository root can legitimately pass `docs/triage/.spacedock-state/task.md`: the shell meaning is the committed state file, while the binary joins that spelling beneath the state checkout and seeks `.spacedock-state/docs/triage/.spacedock-state/task.md`. The same defect affects `--artifact` and `--reference`.

The live audit found this doubled path in three of six preparations over two days. Because the First Officer may invoke prepare only once, each exit 1 withholds the gate rather than costing a retry; one failure stopped the rest of a batch window. The current contracts disagree: `docs/site/reference/command-reference.md` says relative selected sources resolve from the launch working directory, `fo-gate-lifecycle` says to supply judgment/cwd paths, and the binary now treats relative inputs primarily as state-root-relative after an earlier split-root fix.

The fix must preserve both useful meanings. Operators already pass state-relative paths such as `task/index.md`, while transcripts contain launch-cwd-relative paths and callers also use absolute paths. Silently choosing one root is unsafe when the same relative spelling exists at two different committed files.

## Proposed approach

### One resolver with an explicit launch directory

Add an optional `LaunchDir` to `gates.PrepareInput`; the CLI sets it from its existing invocation `dir`, and direct package callers that omit it deterministically default to `WorkflowDir` rather than reading process-global cwd. `resolveSelectedSource` remains the only selected-source resolver, but receives launch dir, workflow dir/entity root, and flag kind.

Resolution is deterministic:

1. An absolute input is cleaned and selected as its sole candidate.
2. A relative input produces de-duplicated lexical candidates under the launch directory and entity root. Preserve the currently supported workflow-root candidate for a state-checkout-basename-prefixed Artifact and for a Reference at the workflow root.
3. Candidate discovery uses `Lstat` only to distinguish absent paths. Exactly one present candidate is selected. Multiple different present candidates refuse before mutation, name each resolved path, and say to use an absolute path. No present candidate refuses with the attempted paths and the three accepted forms.
4. Selection does not certify the file. The chosen path still flows through the existing `gitsource.Inspect`, which remains sole owner of readable regular-file, non-symlink, committed-object, and main/state repository-boundary guards. A dangling symlink, directory, unreadable file, or foreign-repository file therefore cannot be skipped in favor of another interpretation.

This serves AC-1/AC-2. Simplest alternative considered: absolutize all relative values in `internal/cli` before calling `gates.Prepare`. That fixes launch-cwd paths but destroys state-relative inputs because the CLI cannot know which relative meaning the operator intended. Reading `os.Getwd` inside `gates` was also rejected: the in-process runner already carries the true launch dir, and a process-global cwd would make tests and embedded callers depend on hidden machine state.

### Refuse ambiguity instead of assigning silent root priority

There is no cwd-first or state-first winner when two different files exist. Refusal is the precedence rule after absolute inputs: unique existing relative interpretation succeeds; multiple interpretations fail closed; no interpretation reports all attempted locations. The error is attributed to `--artifact` or `--reference`, includes the original spelling and resolved candidates, and directs the next independently authorized attempt to an absolute spelling. The failed invocation leaves entity bytes, room tree, and gate binding unchanged.

This serves AC-3 and protects AC-1 from opening a gate over the wrong committed object. Simplest alternative considered: try the supplied cwd path first, then fall back to state root. It is fewer lines, but a coincidental cwd file silently changes the review authority and makes behavior depend on unrelated filesystem occupancy.

### Align the First Officer and operator documentation

Rewrite the existing path sentence without adding a command or a preflight. The binary owns resolution and refusal; the First Officer keeps supplied paths and uses an absolute spelling before its one prepare when its own evidence scan reveals two meanings. This serves AC-1 by making the successful forms and no-retry ambiguity response explicit. Simplest alternative considered: documentation-only refusal of cwd-relative paths. It contradicts the published cwd contract and leaves all three measured stalls intact.

The concrete `skills/fo-gate-lifecycle/SKILL.md` wording change is:

```diff
-**Prepare.** Resolve `${SPACEDOCK_BIN:-spacedock}`; keep supplied paths. Else, per distinct applicable retained absolute R=`definition_dir`/`entity_dir`, resolve intended root-relative P; state R uses the engaged task directory, never `.`. Run: `git -C "<R>" ls-tree -r --name-only HEAD -- "<P>" | awk 'tolower($0)~/\.(md|markdown)$/{print}'`. `<R>/<P>`: shell-quoted values. Read/use once; summarize. No harness logs/history/status/help probes. Supply judgment/cwd paths; no binary JSON/ids/digests/Git locators/room coords.
+**Prepare.** Resolve `${SPACEDOCK_BIN:-spacedock}` and keep supplied paths. For each retained absolute root R (`definition_dir`, `entity_dir`), derive intended relative P; state R is the engaged task directory, never `.`. Run `git -C "<R>" ls-tree -r --name-only HEAD -- "<P>" | awk 'tolower($0)~/\.(md|markdown)$/{print}'` with quoted R/P. Read once; summarize; no logs/history/status/help probes. Pass each source exactly as judged: absolute, launch-cwd-relative, or state-relative. If one relative spelling can name different cwd and state files, use absolute before the one prepare. No binary JSON/IDs/digests/Git locators/room coordinates.
```

The replacement must keep the skill under its existing 7,700-byte contract cap; this task does not raise the cap. The concrete `docs/site/reference/command-reference.md` change is:

```diff
-Relative selected-source and room inputs resolve from the launch working directory.
+Selected-source inputs may be absolute, launch-working-directory-relative, or state-checkout-relative. A relative spelling that names different existing files under the launch directory and state checkout is refused; use an absolute path. Relative room inputs continue to resolve from the launch working directory.
```

## Risk evidence

The riskiest claim was exercised in a throwaway focused Go spike against `main` at `43fc79a23` and then removed. A split-root fixture placed one committed file at `<main>/docs/dev/.state/selected/review.md` and supplied the launch-cwd-relative spelling `docs/dev/.state/selected/review.md`. Both Artifact and Reference cases reproduced the same result:

- input: `docs/dev/.state/selected/review.md`
- selected seek path: `<main>/docs/dev/.state/docs/dev/.state/selected/review.md`
- outcome: exit path returned `selected source must be a readable non-symlink regular file: ... no such file or directory`

The spike command was `go test ./internal/gates -run TestSpikeCWDRelativeSelectedSourcesDoubleTheStatePath -count=1 -v`; both subtests passed by asserting the doubled-path baseline. This is smaller than a live prepare and directly falsifies the current resolver for both flags. No spike is needed for explicit launch-dir injection or the downstream guards: the CLI already passes its launch `dir` to path resolution elsewhere, and `internal/gitsource` already exercises non-symlink regular-file and repository ownership checks.

## Out of scope

Room layout, Briefing schema or stored Git locators, gate authority, exact-replay behavior, the no-retry rule, state-root discovery, and widening selected sources beyond the workflow main/state Git histories. No new flag or dependency is introduced.

## Expected surface and tolerance

Estimate net LOC change: +145, across 6 files. Expected insertions: 182; expected deletions: 37.

- `internal/gates/prepare.go`: `PrepareInput.LaunchDir`, candidate construction, de-duplication, ambiguity/no-match errors; about +55/−28.
- `internal/cli/cli.go`: pass the existing invocation directory; about +2/−0.
- `internal/gates/prepare_test.go`: equivalent-form, ambiguity, no-match, and guard-preservation tables; about +85/−8 (including replacing the now-obsolete wrong-root expectation).
- `internal/cli/gate_test.go`: behavior-first split-root regression/smoke over the real CLI runner for both flags; about +35/−0.
- `skills/fo-gate-lifecycle/SKILL.md`: the exact bounded rewrite above; about +1/−1 with no byte-cap increase.
- `docs/site/reference/command-reference.md`: the exact operator contract above; about +4/−0.

Tolerance: net +100 to +210 LOC across 6 files, or at most 7 files if the resolver earns one focused helper/test file. Exceeding +210 net LOC, 7 files, the existing 7,700-byte FO skill cap, or introducing a new command flag, dependency, path cache, filesystem search, or second selected-source inspection mechanism requires a return to ideation.

Observable semantics declared by this task: command grammar is unchanged; stored gate/Briefing formats and authority are unchanged; absolute and uniquely resolving state-relative behavior stay unchanged; a uniquely resolving launch-cwd-relative Artifact or Reference now succeeds; a relative spelling with multiple distinct existing candidates now refuses explicitly instead of silently choosing; a spelling with no candidate names all attempted roots instead of reporting only the state-prefixed seek path. Downstream readable-file, symlink, commit, and repository-boundary errors retain their current owners and safety effect.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) — The reproduced split-root gate that is stalled by a launch-cwd-relative Artifact or Reference reaches one open prepared binding without a doubled path.**
Verified by: a CLI fixture launched from the synthetic repository root passes `docs/dev/.state/selected/review.md` first as Artifact and then as Reference in fresh cases. Main-at-`43fc79a23` baseline returns the doubled-path error and creates no binding; the finished command exits 0, prints `state=open`, creates exactly one room/binding, and never creates a doubled directory. Falsified by removing launch-dir candidate construction, which restores exit 1 and zero binding.

**AC-2 — Absolute, launch-cwd-relative, and state-relative spellings of the same selected file are equivalent for both Artifact and Reference.**
Verified by: a table-driven `internal/gates` test prepares once with the absolute spelling and replays with the other two forms, asserting the same `git-root://` URI, raw revision, Briefing digest/id, room, and `state=open`; single-root and currently supported state-basename/workflow-root Reference cases remain green. Falsified by dropping any candidate root or resolving a form to a different repository path.

**AC-3 — Relative-path ambiguity and absence fail before mutation with actionable, flag-attributed evidence, while all selected-source safety guards remain in force.**
Verified by: fixtures commit different files at the launch-dir and state-root candidates for the same relative spelling and assert nonzero refusal names `--artifact`/`--reference`, original spelling, both resolved paths, and absolute disambiguation; zero-match cases name attempted roots and accepted forms. Entity bytes, gate records, room tree, and repository status remain unchanged. A guard table feeds an absolute, cwd-relative, and state-relative symlink/non-regular/unreadable or foreign-repository candidate through the resolver and asserts the existing `gitsource.Inspect` refusals rather than fallback to another file. Falsified by cwd-first selection, post-lock ambiguity detection, or bypassing `gitsource.Inspect`.

**AC-4 — The First Officer and command reference describe exactly the accepted path forms and the no-retry ambiguity response without changing command grammar.**
Verified by: apply the recorded before/after diffs, run the focused CLI regression as the skill-command smoke, confirm the existing exact help fixture is byte-identical, and run the FO component-cap test. Falsified by prose that promises only cwd resolution, tells the FO to retry, adds a flag, or grows the skill beyond 7,700 bytes.

## Test plan

- Before production changes, add the AC-1 CLI regression in `internal/cli/gate_test.go`; it is the focused smoke for the First Officer's existing command text and must fail on the doubled path at current main. Cost: medium fixture extension, no live host.
- Add one table in `internal/gates/prepare_test.go` for AC-2 over `{absolute, launch-cwd-relative, state-relative} × {Artifact, Reference}`. Assert immutable source identity and exact replay output, not merely pass count. Cost: medium; reuse `prepareFixture`.
- Add AC-3 tables for `{unique launch, unique state, ambiguous, absent}` and for downstream safety guards. Assert errors and byte/tree/status non-mutation. Cost: medium; ordinary Go fixtures. The change that must make ambiguity tests fail is restoring any silent root priority.
- Keep and reframe the existing single-root, state-basename, workflow-root Reference, selected-source flag attribution, replay, and `internal/gitsource` guard tests. The existing wrong-root Artifact test changes because a file at the explicit launch root is no longer the wrong root; replace it with the two-present ambiguity case rather than deleting its safety claim.
- After implementation and the exact skill/doc rewrites, run `go test ./internal/gates ./internal/cli ./internal/gitsource ./internal/contractlint`, then repository-required `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. No golden fixture or live workflow test is needed: CLI output/on-disk state is the runtime claim, and the in-process CLI runner injects the real launch directory deterministically.

The test mechanism serves AC-1 through AC-3. Simplest alternative considered: one test containing only the audited Artifact string. It is insufficient because the transcript also failed on Reference, it cannot prove all accepted forms resolve to the same immutable object, and it leaves ambiguity and safety-guard regressions invisible.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Stage Report: ideation

- DONE: Define deterministic precedence and ambiguity handling for absolute, cwd-relative, and state-relative gate Artifact/Reference paths while preserving readable-regular-file, symlink, and repository-boundary guards; record the reproduced doubled-path baseline or a smaller falsifying spike.
  Commit `de8e35fef` records absolute/unique-relative/ambiguous/absent rules and a two-flag spike whose doubled-path assertion would fail if launch-cwd resolution already worked; selected paths still flow through `gitsource.Inspect`.
- DONE: Turn the operator outcome into entity-level acceptance criteria and a behavior-first test plan that proves equivalent accepted forms (or precise supported-form refusals) and pairs the mechanism with the stalled-gate value measure.
  AC-1 measures exit 1/no binding becoming exit 0/one open binding; AC-2/AC-3 require immutable-source equivalence and pre-mutation ambiguity refusal, falsified by removing launch-root candidates or restoring silent priority.
- DONE: Declare expected net LOC/files with tolerance, all observable semantic changes, the concrete first-officer wording/doc diff, and the simplest alternative considered for each new mechanism.
  The body declares +145 net LOC (+182/−37) across 6 files with bounded tolerance, unchanged grammar/storage/authority, exact FO/site diffs, and rejected simpler alternatives per resolver, ambiguity rule, prose alignment, and proof mechanism.

### Summary

Ideation defines one explicit-launch-directory resolver that accepts unique absolute, cwd-relative, and state-relative forms while refusing ambiguous spellings before mutation. A focused throwaway spike reproduced the doubled Artifact and Reference paths, and the acceptance plan ties the fix to reopening the gates that currently stall.
