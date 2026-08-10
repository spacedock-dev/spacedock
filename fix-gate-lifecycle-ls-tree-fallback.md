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

with this complete, directly executable command shape:

```sh
git -C "$GIT_ROOT" ls-tree -r --name-only HEAD -- "$INTENDED_PATH" | awk 'tolower($0) ~ /\.(md|markdown)$/ { print }'
```

`GIT_ROOT` is the Git root that contains the already-intended candidate location. `INTENDED_PATH` is that location expressed relative to the Git root. It is not broadened to `.` unless the intended location is the root. `HEAD` supplies the required committed tree-ish, `-r --name-only` yields recursive candidate paths, `-- "$INTENDED_PATH"` fences discovery, and the `awk` filter admits the same case-insensitive `.md`/`.markdown` extensions that `gate prepare` accepts.

Smaller alternatives do not meet the value boundary:

- Adding only `HEAD` still yields tree entries rather than recursive Markdown candidates and does not fence the intended path.
- `ls-tree -r --name-only HEAD` is committed-only but scans the whole repository and includes non-Markdown files.
- Git `:(glob)` pathspecs would avoid the pipe, but the spike proved `ls-tree` rejects that pathspec magic (`fatal: ... pathspec magic not supported by this command: 'glob'`).
- `git ls-files` observes the index, so a staged but uncommitted Markdown path leaks into discovery. `find` or filesystem globbing likewise observes uncommitted worktree files.

The instruction remains a one-shot read-only discovery fallback. The implementation adds a fixture-backed smoke test that executes the instruction's command shape and a Codex live journey whose prompt omits Artifact and Reference paths, forcing the installed skill to discover the committed package before `gate prepare`.

### Risk spike

A temporary Git fixture contained committed `intended/review.md`, committed `intended/nested/reference.markdown`, committed `unrelated/noise.md`, committed `intended/note.txt`, and staged-only `intended/uncommitted.md`. The proposed command exited 0 and printed exactly the two intended committed Markdown paths. The same fixture's `git ls-files` output included the staged-only path, proving why the shorter index-based alternative is insufficient.

### Expected surface and semantic boundary

- `skills/fo-gate-lifecycle/SKILL.md`: replace the incomplete fallback wording with the exact command; approximately +5/-1 lines.
- `internal/contractlint/fo_gate_lifecycle_fallback_test.go` (new): execute the exact instruction command against a controlled Git fixture; approximately +75 lines.
- `internal/ensigncycle/codex_live_runner_test.go`: add the exact missing-path Codex journey and assertions that discovery precedes successful preparation; approximately +35 lines.

Expected total: 3 files, about +115/-1 lines. Tolerance: one additional existing `internal/ensigncycle` helper/fixture file and up to +60 insertions if needed to reuse the recorded-gate harness without duplicating it. No production Go, stored format, authority, gate lifecycle, supplied-path, CLI grammar, or documentation-site change is permitted. The only permitted runtime semantic change is that the absent-path Codex fallback now returns committed Markdown candidates beneath the already-intended path instead of exiting 129. The site command reference has no fallback-selection wording, so no site documentation diff is proposed.

## Out of scope

- Changing gate authority, preparation, digest, or consume semantics.
- Changing the supplied Artifact/Reference path.
- Broad repository discovery or uncommitted-file discovery.
- Changing how the intended candidate path or Git root is selected.
- Adding a reusable discovery API or production Go implementation.

## Acceptance criteria

**AC-1 (VALUE) - A Codex First Officer with missing Artifact and Reference paths discovers the intended committed Markdown and prepares the gate without a Git usage failure.**
Verified by: an exact local Codex recorded-gate journey whose prompt omits both supplied paths and whose copied `fo-gate-lifecycle` skill is the implementation under test. The test requires a successful discovery command before exactly one successful `gate prepare`, and no Git usage exit; removing `HEAD` makes the journey stop at exit 129.

**AC-2 - The fallback command is complete, read-only, path-scoped, and excludes uncommitted Markdown.**
Verified by: a fixture-backed command test with committed intended `.md` and nested `.markdown`, committed unrelated Markdown, committed non-Markdown, and staged-only intended Markdown. It executes the exact documented command and asserts exit 0 plus the two-path output exactly; removing the intended path, recursive flag, tree-ish, or committed-tree source changes the observed output or exit and fails.

## Test plan

Implementation starts by preserving the spike as the focused fixture-backed smoke test, then changes the skill wording. Run that test directly and confirm the exact output and exit code. Next run the dedicated local Codex missing-path recorded-gate journey and inspect its command trace/on-disk room state: it must use the fallback, prepare exactly once, and produce the expected open gate package without supplied paths. The supplied-path recorded-gate journey remains the regression control and must not invoke discovery.

Finally run the applicable focused `internal/contractlint` and `internal/ensigncycle` packages, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. No generic prose substring check or full live matrix is accepted as proof for AC-1 or AC-2.

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
