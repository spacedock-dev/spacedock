---
id: qbppyssy13pyh1gtkh2n8vp5
title: Fix the gate-lifecycle ls-tree fallback command
status: validation
source: Captain intake; recovered from deleted public issue spacedock-dev/spacedock#669
started: 2026-08-10T22:53:32Z
completed:
verdict:
score: 0.9
worktree: .worktrees/spacedock-ensign-fix-gate-lifecycle-ls-tree-fallback
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
        - id: gate:qbppyssy13pyh1gtkh2n8vp5:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:qbppyssy13pyh1gtkh2n8vp5-ideation-1
              briefing:
                id: briefing:qbppyssy13pyh1gtkh2n8vp5:ideation:attempt-1:revision-1
                digest: sha256:a5b43398f187ec19f74141482044f230b0fe6851b2f942e8949856b0276bd9e4
                request-digest: sha256:3df64bcb04e181fbb2fea036bda380a60028121b2062e963d352c2ef1d7ea145
                room-ref: ./fix-gate-lifecycle-ls-tree-fallback/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:qbppyssy13pyh1gtkh2n8vp5:ideation:1
                briefing: briefing:qbppyssy13pyh1gtkh2n8vp5:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T23:11:02.302967Z"
                decision: approve
                reason: Captain directed dispatch; independent review passed the one-skill, no-new-test, split-root design with falsifiable existing-journey proof.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:qbppyssy13pyh1gtkh2n8vp5:validation
          stage: validation
          attempts:
            - id: gate-attempt:qbppyssy13pyh1gtkh2n8vp5-validation-1
              briefing:
                id: briefing:qbppyssy13pyh1gtkh2n8vp5:validation:attempt-1:revision-1
                digest: sha256:bbf3c3b6872ec6c61acfd87bfdce5b9c85a5e75b17af1e51d9062c2439e9999b
                request-digest: sha256:0b9e6a2acc234fb999602ee392e53ac0910b92f68beb011c3713ffa8106387e4
                room-ref: ./fix-gate-lifecycle-ls-tree-fallback/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:qbppyssy13pyh1gtkh2n8vp5:validation:1
                briefing: briefing:qbppyssy13pyh1gtkh2n8vp5:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-11T03:23:31.885533Z"
                decision: revise
                reason: 'Accept validation rejection and route QBP to implementation: bind the state query to the known task directory and make proof detect whole-root discovery without adding a new test framework.'
            - id: gate-attempt:qbppyssy13pyh1gtkh2n8vp5-validation-2
              briefing:
                id: briefing:qbppyssy13pyh1gtkh2n8vp5:validation:attempt-2:revision-1
                digest: sha256:cdfb43872d56046662630a0f480af8b8893bfb9bfb8c2ff8f8fb251434aedd49
                request-digest: sha256:b41caf52cc6ab2a92fab08427908271c59d7e4dd07e3d35040349aa562733c85
                room-ref: ./fix-gate-lifecycle-ls-tree-fallback/review/validation/briefing-2
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

with this complete command shape after resolving and substituting concrete values, once for each distinct applicable retained root:

```sh
git -C "<resolved-git-root>" ls-tree -r --name-only HEAD -- "<root-relative-intended-path>" | awk 'tolower($0) ~ /\.(md|markdown)$/ { print }'
```

Before execution, the First Officer considers both retained absolute roots from boot: `definition_dir` and `entity_dir`. For each distinct root that can contain an intended gate source, it resolves the intended candidate location within that root and converts it to a root-relative path. In the existing split-root guardrail, that means the definition-root location containing `recorder-contract.md` and the known engaged entity task directory `recorded-gate-task` containing the review and snapshot. For the state root, `P` is that task directory and may not collapse to `.` while the task directory is known. It substitutes each shell-quoted root/path pair into one invocation; the angle-bracket terms are notation, not environment variables. If both retained roots are identical, the identical root/path pair is queried once. `HEAD` supplies the required committed tree-ish, `-r --name-only` yields recursive candidate paths, the value after `--` fences each root's discovery, and the `awk` filter admits the same case-insensitive `.md`/`.markdown` extensions that `gate prepare` accepts.

Smaller alternatives do not meet the value boundary:

- Adding only `HEAD` still yields tree entries rather than recursive Markdown candidates and does not fence the intended path.
- `ls-tree -r --name-only HEAD` is committed-only but scans the whole repository and includes non-Markdown files.
- Git `:(glob)` pathspecs would avoid the pipe, but the spike proved `ls-tree` rejects that pathspec magic (`fatal: ... pathspec magic not supported by this command: 'glob'`).
- `git ls-files` observes the index, so a staged but uncommitted Markdown path leaks into discovery. `find` or filesystem globbing likewise observes uncommitted worktree files.

The instruction remains a one-shot read-only discovery fallback with one query per applicable retained root, not repeated probing. Validation uses one exact local Codex `TestLiveCommonGateGuardrail` journey and inspects its observed command structure directly: the state-root invocation must use `P=recorded-gate-task`, and `P=.` fails AC-2. This is direct journey evidence paired with persisted package identities, not command-output comparison or standing test infrastructure.

### Risk spike

A temporary Git fixture contained committed `intended/review.md`, committed `intended/nested/reference.markdown`, committed `unrelated/noise.md`, committed `intended/note.txt`, and staged-only `intended/uncommitted.md`. The proposed command exited 0 and printed exactly the two intended committed Markdown paths. The same fixture's `git ls-files` output included the staged-only path, proving why the shorter index-based alternative is insufficient.

### Expected surface and semantic boundary

- `skills/fo-gate-lifecycle/SKILL.md`: bind state-root `P` to the engaged task directory while retaining the complete command.

Hard boundary: exactly 1 product file at +2/-2 lines versus `main`, with no enlarged tolerance. No Go, test, fixture, harness, stored format, authority, gate lifecycle, supplied-path, CLI grammar, or documentation-site change is permitted. The only runtime change remains absent-path discovery requiring the state-root query to use the known engaged task directory rather than `.`.

## Out of scope

- Changing gate authority, preparation, digest, or consume semantics.
- Changing the supplied Artifact/Reference path.
- Broad repository discovery or uncommitted-file discovery.
- Introducing another root/path policy beyond the retained `definition_dir`/`entity_dir` and their resolved intended candidate locations.
- Adding a reusable discovery API or production Go implementation.
- Adding or modifying Go tests, fixtures, live journeys, frameworks, or command-output oracles.

## Acceptance criteria

**AC-1 (VALUE) - A Codex First Officer with missing Artifact and Reference paths discovers the intended committed Markdown and prepares the gate without a Git usage failure.**
Verified by: existing `TestLiveCommonGateGuardrail` under `SPACEDOCK_LIVE_RUNTIME=codex`. Its existing prompt supplies no Artifact or Reference paths, loads the copied product skill, and must leave exactly one persisted prepared gate package open at the human decision boundary. Removing the explicit tree-ish prevents that package from being created and fails the existing durable assertion.

**AC-2 - The fallback command is complete, read-only, path-scoped, and excludes uncommitted Markdown.**
Verified by: one exact local Codex `TestLiveCommonGateGuardrail` journey. Its observed command structure must contain the state-root `ls-tree` invocation with `HEAD` and `P=recorded-gate-task`, not `P=.`, and its persisted package must retain the expected definition-root contract plus state-root review/snapshot identities. A state invocation using `-- .`, a missing retained root, a missing `HEAD`, or an incorrect persisted identity fails AC-2; no command output is compared.

## Test plan

Implementation does not run the model journey. It runs formatting, applicable focused existing gate tests, the full non-race suite, the full race suite, and `git diff --check` as deterministic repository-health checks. Validation then runs exactly one local `SPACEDOCK_LIVE_RUNTIME=codex` `TestLiveCommonGateGuardrail` journey and directly inspects its observed command structure for state `P=recorded-gate-task` rather than `.`, together with the persisted package identities required by AC-1 and AC-2. No command-output comparison, duplicate journey, new fixture, or standing assertion infrastructure is added.

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

## Stage Report: ideation (cycle 3)

- DONE: Define the exact before/after fallback command and why the smaller alternatives do not satisfy intended-path, committed-only discovery.
  The command remains complete and path-fenced, now invoked once for each distinct applicable retained `definition_dir`/`entity_dir` root with a concrete root-relative intended path.
- DONE: Declare expected files, insertion count, tolerance, and the only permitted runtime semantic change.
  The one-skill-file, about +5/-1 surface and +10/-2 same-file tolerance remain; only missing-path discovery expands from one retained root to all applicable retained roots.
- DONE: Tie each acceptance criterion to fixture-backed command proof and the exact missing-path Codex journey; exercise the riskiest Git command shape first.
  AC-1 still requires exactly one persisted package from `TestLiveCommonGateGuardrail`; AC-2 now requires its inventory to retain expected committed definition-root and entity-root identities while excluding unrelated and uncommitted inputs.

### Summary

The split-root correction preserves the approved one-file implementation and existing-journey proof. It prevents a state-only query and incomplete package from passing by making both applicable roots and their persisted selected identities part of the design boundary.

## Stage Report: implementation

- DONE: The shipped skill gives one complete path-fenced committed-tree command per applicable retained root and defines every substituted root/path value.
  Commit `81e1dbdfe` defines R as each distinct applicable retained absolute `definition_dir`/`entity_dir`, P as its resolved root-relative intended location, and supplies the complete `HEAD -- "<P>"` command.
- DONE: The change preserves both-root intended source selection while excluding unrelated and uncommitted inputs, without changing supplied paths or gate authority.
  A temporary Git fixture returned only committed `intended/review.md` and `intended/nested/reference.markdown`; removing `HEAD` or the `intended` fence would admit the staged-only or unrelated controls.
- DONE: The candidate stays within the one-file +10/-2 tolerance and passes applicable focused, full, race, formatting, and diff checks; no new test infrastructure is added.
  Formatting, non-race, and diff checks passed unchanged; earlier capacity-constrained parallel and serialized race attempts timed out, but after cache-only capacity remediation the exact `go test ./... -race` passed (`internal/cli` 192.562s; `internal/ensigncycle` 246.836s) with no test or infrastructure change.

### Summary

The gate-lifecycle fallback now supplies one committed, recursive, Markdown-filtered and intended-path-fenced tree query for each applicable retained root while leaving supplied paths unchanged. Product commit `81e1dbdfe` remains byte-clean, and the required full race suite passed unchanged after external capacity remediation; the report retains the prior constrained-host failures as history.

## Review-finding disposition

### V-1 — Live state-root fallback used the whole repository

- Reviewer observation: Codex item_11 invoked both `ls-tree` commands with `P='.'` including the state checkout, rather than the designed state task directory. The persisted-package oracle cannot distinguish this whole-root scan because the fixture has no unrelated committed state Markdown.
- Released user and normal workflow: a Codex First Officer preparing a missing-path gate in the supported split-root `TestLiveCommonGateGuardrail` workflow.
- Observable harm: the fallback exposes every committed Markdown path in the state checkout to discovery instead of fencing selection to `recorded-gate-task`; the correct final package does not erase that broader observation boundary.
- Affected authority: `contract[skills/fo-gate-lifecycle/SKILL.md#Prepare]` requires each retained root query to use its resolved root-relative intended location `P`, while the task boundary forbids broad repository discovery.
- Trigger evidence: the sole authorized live run recorded completed Codex command item `item_11` with exit 0 and `-- .` for both roots; the state root and known slug/task directory were already available.
- Defect kind: outcome defect (the state query is whole-root) and evidence defect (the persisted inventory still passes, so the cited oracle cannot detect the wrong fence).
- Release scope: Material; the trigger is the supported missing-path split-root journey, and it violates the explicit discovery-scope contract on that normal path.
- Task ownership: the current task owns the fallback and its AC proof, but the authorized one-file candidate cannot repair both the runtime instruction result and the observation boundary.
- Correction shape: mechanism/design reset, not a bounded correction; re-anchor `P` to the known entity task directory and make the proof distinguish intended-path discovery from whole-root discovery.
- Proposed disposition: route for decision.

## Stage Report: validation

- DONE: Reproduce AC-1 with the existing exact local Codex TestLiveCommonGateGuardrail and verify exactly one prepared package without the Git usage failure.
  The single `SPACEDOCK_LIVE_RUNTIME=codex` run passed in 121.04s; one `request.json` was held open, `gate prepare` exited 0 exactly once, and the retained stream had no exit-129/`git ls-tree` usage failure.
- FAILED: Reproduce AC-2 from persisted package/source identities across definition_dir and entity_dir, excluding unrelated and uncommitted inputs without command-output comparison.
  The canonical Briefing had exactly state `gate-review.md`, main `recorder-contract.md`, and state `entity-snapshot.md` with independently matching SHA-256 revisions and no dirty sibling, but completed command item `item_11` used state `P='.'`; the package-only oracle therefore passes the prohibited whole-root scan.
- DONE: Verify the one-file +2/-2 candidate, applicable focused/full/race/formatting checks, and unchanged supplied-path and gate-authority behavior; recommend PASSED or REJECTED.
  Commit `81e1dbdfe` is one file at +2/-2 and `git diff --check` clean; focused tests, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` passed, while the live log stopped after one prepare/commit with no decision, consume, dispatch, or archive and the supplied-path branch remained unchanged.
- FAILED: Recommend PASSED or REJECTED with findings classified under workflow policy.
  REJECTED: V-1 is a Material outcome-plus-evidence defect in the chosen mechanism/proof boundary; route the design reset for decision without mutating or rerunning candidate `81e1dbdfe`.

### Summary

AC-1 and the exact persisted inventory passed in the sole authorized Codex journey, and all offline, race, formatting, surface, supplied-path, and gate-authority checks are green. Validation nevertheless recommends REJECTED because the live state query used the whole root and the specified persisted-package proof cannot detect that path-scope failure; V-1 requires a mechanism/design reset rather than a silent bounded correction.

## Stage Report: implementation (cycle 2)

- DONE: The fallback resolves the state-root path to the known task directory and cannot collapse it to `.` when the task directory is known.
  Final correction commit `69139fb1a` leaves the product skill requiring state-root `P` to use the engaged task directory and explicitly forbidding `.`, while removing every test/harness delta.
- DONE: Existing fixture or semantic proof makes whole-root discovery fail while preserving both-root intended selection, without a new test framework or command-output comparison.
  The authorized proof is one exact local Codex gate-guardrail journey in validation: its observed state-root command must use `P=recorded-gate-task` rather than `.`, and its persisted package must retain both-root intended identities; this direct evidence compares no command output and adds no infrastructure.
- DONE: Task design, acceptance evidence, test plan, implementation surface, and deterministic verification are updated consistently; the spent Codex journey is not rerun during correction.
  The final candidate is exactly 1 file at +2/-2 against `main`, with no enlarged tolerance. `gofmt -w ./cmd ./internal`, focused existing gate tests, `go test ./...`, `go test ./... -race`, and `git diff --check` passed (`internal/cli` 217.646s/109.244s and `internal/ensigncycle` 287.604s/183.767s for non-race/race); no model journey ran during implementation.

### Summary

The corrected fallback now fences state discovery to the known engaged task directory in a one-file +2/-2 candidate. Candidate `69139fb1a` passes deterministic full and race verification; the exact local Codex journey remains reserved for validation, where state `P=.` must fail the direct observed-command and persisted-package proof.

## Stage Report: validation (cycle 2)

- DONE: Reproduce AC-1 with the existing exact local Codex TestLiveCommonGateGuardrail and verify exactly one prepared package without the Git usage failure.
  The sole authorized cycle-2 Codex journey passed in 94.43s; completed command item `item_8` exited 0, exactly one gate package was prepared and held, and the retained evidence contains no exit-129 or `git ls-tree` usage failure.
- DONE: Reproduce AC-2 from persisted package/source identities across definition_dir and entity_dir, excluding unrelated and uncommitted inputs without command-output comparison.
  The canonical Briefing contains exactly state `gate-review.md`, main `recorder-contract.md`, and state `entity-snapshot.md` with matching immutable revisions and no dirty sibling; direct grading of item `item_8` shows main `P=.` and state `P=recorded-gate-task`, resolving V-1 without comparing command output.
- DONE: Verify the one-file +2/-2 candidate, applicable focused/full/race/formatting checks, and unchanged supplied-path and gate-authority behavior; recommend PASSED or REJECTED.
  Candidate `69139fb1a` is clean at one file +2/-2; focused tests, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` passed, while the live log has one prepare/commit and no decision, consume, successor dispatch, or archive; supplied paths remain outside fallback discovery.
- DONE: Recommend PASSED or REJECTED with material, deferred-risk, and polish findings separated.
  PASSED: V-1 is resolved by the observed state task-directory fence; no material, deferred-risk, or polish findings remain.

### Summary

Both acceptance criteria pass in the single authorized cycle-2 Codex journey, including direct observation of state `P=recorded-gate-task` and the exact persisted three-source inventory. The candidate stays within its one-file +2/-2 boundary, all deterministic checks pass, authority remains held at the human gate, and validation recommends PASSED.
