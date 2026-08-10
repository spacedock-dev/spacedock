---
id: qbppyssy13pyh1gtkh2n8vp5
title: Fix the gate-lifecycle ls-tree fallback command
status: ideation
source: Captain intake; recovered from deleted public issue spacedock-dev/spacedock#669
started: 2026-08-10T22:53:32Z
completed:
verdict:
score: 0.9
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:qbppyssy13pyh1gtkh2n8vp5:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:qbppyssy13pyh1gtkh2n8vp5-backlog-1
              briefing:
                id: briefing:qbppyssy13pyh1gtkh2n8vp5:backlog:attempt-1:revision-1
                digest: sha256:7ca8a11685cfba977d897cea574e18061a7fe81a1cb8a565f4c272f5a09f7593
                request-digest: sha256:90c2970ab0e6992546e9f0ae718f393fc3bc37336d7a3a48669896d5b76cfb08
                room-ref: ./fix-gate-lifecycle-ls-tree-fallback/review/backlog/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-10T22:37:04.823517Z"
                reason: Backlog record has no Stage Report, so checklist and acceptance evidence cannot be assembled.
            - id: gate-attempt:qbppyssy13pyh1gtkh2n8vp5-backlog-2
              briefing:
                id: briefing:qbppyssy13pyh1gtkh2n8vp5:backlog:attempt-2:revision-1
                digest: sha256:819dd15e568e3a55efc09967edd4acfaf7485728be52ca08b89b892103d09a25
                request-digest: sha256:fc1a79f889bcfcc19d8666e6dae5f108ce568c2e56121a2999103e1f436e45e5
                room-ref: ./fix-gate-lifecycle-ls-tree-fallback/review/backlog/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:qbppyssy13pyh1gtkh2n8vp5:backlog:2
                briefing: briefing:qbppyssy13pyh1gtkh2n8vp5:backlog:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-10T22:53:09.147634Z"
                decision: approve
                reason: 'Captain directed dispatch of QBP; the recovered seed is bounded to the PR #659 fallback regression and defines falsifiable fixture and exact-Codex proof.'
              application:
                target-stage: ideation
                state: consumed
---

When gate Artifact or Reference paths are absent, the First Officer must discover committed Markdown with a complete, path-scoped Git command. The current instruction abbreviates this as `git -C ... ls-tree`, which deterministically exits 129 because it omits the required tree-ish.

## Problem

Codex can follow the fallback instruction literally during gate preparation. The command fails before artifact selection, even though offline checks and all supplied-path routes can pass. This makes a deterministic product instruction defect look intermittent.

The defect entered through PR #659, commit `bbfad5b4c7886dbdee797e66e34e67a348d05cfd`. Local reproduction:

```console
$ git -C . ls-tree
usage: git ls-tree [<options>] <tree-ish> [<path>...]
$ echo $?
129
```

The fallback needs one committed-tree query, not worktree or index discovery. It must stay under the intended candidate path already selected from the gate/workflow context; this task does not introduce a new search root or selection policy.

## Proposed approach

Replace the abbreviated fallback command:

```sh
git -C ... ls-tree
```

with this complete command shape after resolving and substituting concrete values:

```sh
git -C "<resolved-git-root>" ls-tree -r --name-only HEAD -- "<root-relative-intended-path>" | awk 'tolower($0) ~ /\.(md|markdown)$/ { print }'
```

Before execution, the First Officer resolves `<resolved-git-root>` to the retained absolute `entity_dir` from boot and resolves `<root-relative-intended-path>` to the path-resolved entity's containing directory relative to that root. It substitutes those two shell-quoted values into the command; they are notation, not environment variables. It does not broaden the intended path to `.` unless the resolved entity is flat at the entity root. `HEAD` supplies the required committed tree-ish, `-r --name-only` yields recursive candidate paths, the value after `--` fences discovery, and the `awk` filter admits the same case-insensitive `.md`/`.markdown` extensions that `gate prepare` accepts.

Smaller alternatives do not meet the value boundary:

- Adding only `HEAD` still yields tree entries rather than recursive Markdown candidates and does not fence the intended path.
- `ls-tree -r --name-only HEAD` is committed-only but scans the whole repository and includes non-Markdown files.
- Git `:(glob)` pathspecs would avoid the pipe, but the spike proved `ls-tree` rejects that pathspec magic (`fatal: ... pathspec magic not supported by this command: 'glob'`).
- `git ls-files` observes the index, so a staged but uncommitted Markdown path leaks into discovery. `find` or filesystem globbing likewise observes uncommitted worktree files.

The instruction remains a one-shot read-only discovery fallback. No test or harness mechanism is added: the existing `TestLiveCommonGateGuardrail` prompt already omits Artifact and Reference paths, forcing the installed skill to discover committed sources before `gate prepare`, and its durable assertions require exactly one prepared package held at the human decision boundary.

### Risk spike

A temporary Git fixture contained committed `intended/review.md`, committed `intended/nested/reference.markdown`, committed `unrelated/noise.md`, committed `intended/note.txt`, and staged-only `intended/uncommitted.md`. The proposed command exited 0 and printed exactly the two intended committed Markdown paths. The same fixture's `git ls-files` output included the staged-only path, proving why the shorter index-based alternative is insufficient.

### Expected surface and semantic boundary

- `skills/fo-gate-lifecycle/SKILL.md`: replace the incomplete fallback wording with the exact command; approximately +5/-1 lines.

Expected total: exactly 1 product file, about +5/-1 lines. Tolerance: up to +10/-2 lines in that same skill file and no additional files. No Go, fixture, harness, stored format, authority, gate lifecycle, supplied-path, CLI grammar, or documentation-site change is permitted. The only permitted runtime semantic change is that the absent-path Codex fallback uses the resolved entity Git root and root-relative intended directory to select committed Markdown instead of exiting 129. The site command reference has no fallback-selection wording, so no site documentation diff is proposed.

## Out of scope

- Changing gate authority, preparation, digest, or consume semantics.
- Changing the supplied Artifact/Reference path.
- Broad repository discovery or uncommitted-file discovery.
- Introducing another root/path policy beyond the retained `entity_dir` and path-resolved entity directory.
- Adding a reusable discovery API or production Go implementation.
- Adding or modifying Go tests, fixtures, live journeys, or harness code.

## Acceptance criteria

**AC-1 (VALUE) - A Codex First Officer with missing Artifact and Reference paths discovers the intended committed Markdown and prepares the gate without a Git usage failure.**
Verified by: existing `TestLiveCommonGateGuardrail` under `SPACEDOCK_LIVE_RUNTIME=codex`. Its existing prompt supplies no Artifact or Reference paths, loads the copied product skill, and must leave exactly one persisted prepared gate package open at the human decision boundary. Removing the explicit tree-ish prevents that package from being created and fails the existing durable assertion.

**AC-2 - The fallback command is complete, read-only, path-scoped, and excludes uncommitted Markdown.**
Verified by: inspect the persisted package produced by that same existing journey and its selected source identities, not console output. Every selected source must resolve beneath the entity's intended directory at `HEAD`; the unrelated dirty sibling and any staged-only input must be absent. Removing `HEAD` or the intended-path fence either prevents preparation or changes the persisted selection and fails this package/state check.

## Test plan

Implementation changes only `skills/fo-gate-lifecycle/SKILL.md`; it adds no standing test or fixture. Validation runs the existing `TestLiveCommonGateGuardrail` once with `SPACEDOCK_LIVE_RUNTIME=codex`. AC-1 checks the resulting on-disk gate state for exactly one open prepared package. AC-2 opens that package and verifies its selected source identities are committed at `HEAD`, remain under the resolved intended directory, and exclude the fixture's unrelated dirty sibling and any staged-only input. No command-output comparison, prose substring check, duplicate Codex journey, full live matrix, or new test infrastructure is permitted.

The supplied-path behavior is outside the changed branch and receives no additional journey. The one-off Git spike remains risk evidence for the command mechanism, not acceptance proof and not a committed test.

## Stage Report: backlog

Seed outcome: When gate paths are absent, Codex discovers the intended committed Markdown and prepares the gate without a Git usage failure.

Included scope: Complete the read-only, path-scoped `git ls-tree` command with an explicit tree-ish. Keep discovery limited to committed Markdown. Prove the command in a Git fixture and in the Codex missing-path journey.

Excluded scope: Do not change gate authority, gate record semantics, supplied-path behavior, or discovery of uncommitted files.

Proof needed for ideation: Reproduce the PR #659 regression, run the proposed command against a controlled Git fixture, and confirm that the Codex journey reaches gate preparation. The evidence must show runtime behavior, not compare instruction text.

## Stage Report: ideation

- DONE: Define the exact before/after fallback command and why the smaller alternatives do not satisfy intended-path, committed-only discovery.
  The body specifies `ls-tree -r --name-only HEAD -- <intended-path> | awk ...`; the spike also disproved `:(glob)`, while `ls-files`, broad `ls-tree`, and filesystem discovery violate the boundary.
- DONE: Declare expected files, insertion count, tolerance, and the only permitted runtime semantic change.
  The design declares 3 files, about +115/-1 lines, one helper-file/+60-line tolerance, and only repairs absent-path committed Markdown discovery; production Go and all gate/supplied-path semantics remain unchanged.
- DONE: Tie each acceptance criterion to fixture-backed command proof and the exact missing-path Codex journey; exercise the riskiest Git command shape first.
  The temporary Git fixture produced exactly two intended committed Markdown paths and excluded unrelated, non-Markdown, and staged-only entries; AC-1 names the forced missing-path Codex journey and AC-2 names the exact-output fixture test.

### Summary

Ideation now defines the complete fallback command, a narrow semantic boundary, and falsifiable proof for both the Git behavior and the exact Codex journey. The risk spike found and avoided unsupported `ls-tree` glob magic, then proved the selected path-scoped committed-tree pipeline before any implementation work.

## Stage Report: ideation (cycle 2)

- DONE: Define the exact before/after fallback command and why the smaller alternatives do not satisfy intended-path, committed-only discovery.
  The revised design makes the Git root and root-relative intended directory concrete substitutions resolved from retained workflow state, while preserving the exercised `HEAD`/recursive/path-fenced command shape.
- DONE: Declare expected files, insertion count, tolerance, and the only permitted runtime semantic change.
  The authorized surface is exactly `skills/fo-gate-lifecycle/SKILL.md`, about +5/-1 lines with +10/-2 same-file tolerance; all Go, fixture, harness, and unrelated semantic changes are forbidden.
- DONE: Tie each acceptance criterion to fixture-backed command proof and the exact missing-path Codex journey; exercise the riskiest Git command shape first.
  Both ACs now reuse `TestLiveCommonGateGuardrail`: AC-1 requires exactly one persisted prepared package, while AC-2 inspects that package's committed, intended-path source identities and excludes unrelated or staged-only inputs without comparing console output.

### Summary

The authorized correction removes both proposed test mechanisms and narrows implementation to one product skill edit. Proof now comes from the existing missing-path Codex journey and its persisted package, with explicit root/path substitution and no command-output oracle.
