---
title: Restore the shared ensign contract at the Codex fresh-dispatch boundary
status: implementation
source: "Runtime Live E2E run 31045591048, Codex job 92440554439: the first full-ensign-cycle worker received only a Read pointer, never loaded the shared ensign contract, and wrote two generic checkbox Stage Reports that the FO accepted before archive."
started: 2026-08-07T13:16:17Z
completed:
verdict:
score: 0.95
worktree:
issue:
sprint: test-behavior-completeness
id: nvz2ym82ydfn07jp04yfxg9r
gates:
    version: 1
    records:
        - id: gate:nvz2ym82ydfn07jp04yfxg9r:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:nvz2ym82ydfn07jp04yfxg9r-backlog-1
              briefing:
                id: briefing:nvz2ym82ydfn07jp04yfxg9r:backlog:attempt-1:revision-1
                digest: sha256:1bd6b055d6ff099786e2ed9dda6a9e4d74215a7e819a7994db71632af9f3db9e
                request-digest: sha256:2cb0362a63ad80bad420e625983ce003db17ec09ecfd5f0026590587b5be25d8
                room-ref: ./restore-codex-ensign-contract-bootstrap/review/backlog/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-07T05:03:19.891718Z"
                reason: Captain will dispatch NV separately outside the durable-decisions sprint.
            - id: gate-attempt:nvz2ym82ydfn07jp04yfxg9r-backlog-2
              briefing:
                id: briefing:nvz2ym82ydfn07jp04yfxg9r:backlog:attempt-2:revision-1
                digest: sha256:49b8a62bfcc38e9230f9f7a2db34680e1e586d551cf98acacd2882dddf2bc0a0
                request-digest: sha256:8820b1554093e427ca253fbff880ebe4d7a02158d808943f1597fc375887bf57
                room-ref: ./restore-codex-ensign-contract-bootstrap/review/backlog/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:nvz2ym82ydfn07jp04yfxg9r:backlog:2
                briefing: briefing:nvz2ym82ydfn07jp04yfxg9r:backlog:attempt-2:revision-1
                by: person:captain
                at: "2026-08-07T13:15:04.753368Z"
                decision: approve
                reason: Captain accepts the task direction and approves transition to ideation.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:nvz2ym82ydfn07jp04yfxg9r:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:nvz2ym82ydfn07jp04yfxg9r-ideation-1
              briefing:
                id: briefing:nvz2ym82ydfn07jp04yfxg9r:ideation:attempt-1:revision-1
                digest: sha256:0834a43f42a25251fb05daf7d946511a8ce887e149dfee57847f701b42237783
                request-digest: sha256:bdbf442f28214e48d96bf4eb730fd3ff5e230a7582ae9ebf5568eefe8d9c856c
                room-ref: ./restore-codex-ensign-contract-bootstrap/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nvz2ym82ydfn07jp04yfxg9r:ideation:1
                briefing: briefing:nvz2ym82ydfn07jp04yfxg9r:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-08T00:53:46.562114Z"
                decision: approve
                reason: Captain accepts the ideation design and approves implementation.
              application:
                target-stage: implementation
                state: consumed
---

## Outcome

Every fresh Codex stage worker receives the shared ensign operating contract before it executes the dispatch assignment. The worker therefore follows the exact Stage Report protocol, and the First Officer never treats a generic `## Stage Report` checkbox section as contract-complete.

## Problem

The Codex helper emits a pointer-only outer prompt that says only `Read /tmp/spacedock-dispatch/...`. Its generated file claims to contain the shared ensign discipline, but it contains only the host first-action note, stage/entity/checklist context, fetches, and completion signal. It does not contain `skills/ensign/references/ensign-shared-core.md` or another executable bootstrap to that contract.

Exact-head run `31045591048` reached this supported path on the first `full-ensign-cycle` journey. The isolated Codex child appended `## Stage Report` with markdown checkbox items. The First Officer sent one report-repair turn, accepted a second generic checkbox report, advanced through the workflow, and archived the entity. The live oracle correctly failed because no anchored `## Stage Report: {stage}` existed.

The archived Codex runtime adapter task (`r09jrf0k6qjv6c1sddhe1sh6`) deliberately removed `Skill(skill="spacedock:ensign")` in June 2026 because that Codex skill surface was then unavailable. That historical choice is not proof that the current host cannot load the skill. The current dispatch file's claim that it already carries the shared discipline is false.

## Spike results

The riskiest question was whether the current fresh Codex child can explicitly load the installed Spacedock ensign skill at the dispatch boundary. The spike used the current Codex CLI (`codex-cli 0.146.0`) and the installed Spacedock plugin checkout. The decisive supported-host-shaped probe was:

```text
codex exec --ephemeral --skip-git-repo-check \
  --dangerously-bypass-approvals-and-sandbox \
  -C /Users/clkao/git/spacedock-research/spacedock-v1 --json \
  '$spacedock:ensign; then Read /tmp/spacedock-dispatch/codex-bootstrap-spike-spacedock-ensign-nvz2ym82yd-ideation.md and treat its content as your assignment. Bootstrap probe only: do not execute the assignment; after reading, report whether the shared Stage Report protocol was loaded.'
```

The child recognized `$spacedock:ensign`, read `skills/ensign/SKILL.md`, `ensign-shared-core.md`, and `codex-ensign-runtime.md`, then read the dispatch pointer. Its final answer confirmed that the shared Stage Report protocol was loaded, including the anchored `## Stage Report: {stage_name}` form. The child actually issued `sed` reads against those contract files; this was not only a marker assertion.

The control probe used the current generated pointer-only prompt. It produced no contract load and could not establish the Stage Report protocol. A separate read-only probe also exposed an environment boundary: a child whose workspace did not include `/tmp` could not follow the pointer. The supported live Codex lane already launches with the Spacedock plugin checkout and the required workspace shape, so that sandbox control is a harness concern, not a reason to inline the contract into the dispatch artifact.

The spike selects explicit skill invocation over self-contained contract bytes. The `$spacedock:ensign` marker must be in the outer fresh-child prompt, before the dispatch pointer; putting it inside the file would be too late because the child must load the skill before it reads the assignment. The current host can load the skill, so duplicating the shared core into every dispatch file would add drift and prompt size without solving a demonstrated problem.

## Boundary

This task owns Codex fresh-worker contract bootstrap across the dispatch builder, Codex runtime binding, ensign entry point, and their behavioral tests. It does not change the live oracle, weaken Stage Report validation, restore parent-turn inheritance, reconstruct assignment payload in the First Officer, or absorb `self-contained-ensign-dispatch` (`kd7877nnbd19d528xnpwwaj4`), whose boundary explicitly excludes prompt/bootstrap changes. Preserve `fork_turns: "none"`, helper-prompt forwarding, pointer-only assignment transport, and Claude/Pi behavior.

The sibling `self-contained-ensign-dispatch` work may change how stage/context bytes are represented in the generated artifact. This task must re-anchor its small Codex prompt change on that resulting builder rather than cherry-pick stale hunks or take ownership of the sibling's artifact-payload, launcher-pinning, or assignment-transport work. The Codex bootstrap remains a fresh-worker contract concern even if the sibling later makes the assignment file more self-contained.

## Acceptance criteria

**AC-1 (VALUE) — A fresh Codex ensign follows the shared Stage Report protocol.** A real fresh Codex stage worker emits an anchored `## Stage Report: {stage}` with DONE/SKIPPED/FAILED accounting, `### Summary`, and no checkbox bullets. The independent baseline is the exact-head run `31045591048`, which produced zero contract-compliant fresh Codex reports in the first `full-ensign-cycle` attempt and was archived only after generic checkbox reports were accepted. The target fixed result is one passing clean candidate in that common journey, while its generic no-colon checkbox negative remains rejected.

Verified by: the Codex `full-ensign-cycle` common live journey from a clean exact candidate; an offline report parser/fixture that presents the generic `## Stage Report` checkbox shape; and the existing oracle assertion that requires the anchored stage heading, accounting, and summary.

**AC-2 — The bootstrap is executable, not a false prose claim.** The emitted Codex fresh prompt begins exactly with the explicit current-host skill invocation `$spacedock:ensign; then Read ` followed by the generated pointer. The generated file says that the outer prompt supplies the shared contract instead of claiming that the file itself contains it. The skill load happens before assignment execution, and the implementation does not inline a second copy of `ensign-shared-core.md`.

Verified by: a fixture-backed dispatch-to-child probe that reads the real ensign contract only when the fixed prompt edge is present; a removal mutation that deletes `$spacedock:ensign` while retaining the old false prose and must fail; and the completed current-host fresh-child spike recorded above.

**AC-3 — Fresh-context and prompt ownership remain intact.** Codex spawn still uses `fork_turns: "none"`; the First Officer forwards the helper-emitted prompt byte-for-byte; and stage/entity/checklist payload remains in the dispatch artifact rather than being reconstructed by the runtime adapter. The fixed bootstrap applies only to a fresh worker. A reuse/advance pointer remains the existing short prompt and does not invoke the skill again.

Verified by: the existing exact spawn-map, Codex adapter, pointer-identity, and advance-prompt-size tests; a new assertion that the fresh prompt has the fixed prefix while the advance prompt does not; and a removal mutation for the bootstrap edge.

**AC-4 — Other hosts do not regress.** Claude retains `Skill(skill="spacedock:ensign"); then Read ...`, Pi retains its supported pointer behavior, and all hosts retain their existing completion-signal and worktree rules. No First Officer validation rule or live oracle is weakened.

Verified by: focused cross-host dispatch tests, the full and race gates, and the live lanes required by the final changed-path set. The Codex live lane uses the existing `spacedock codex --plugin-dir <checkout>` setup so the skill invocation is exercised in the same plugin boundary that production dispatch uses.

## Proposed approach

1. Change only the Codex branch of the fresh pointer prompt in `internal/dispatch/build.go` to emit:

   ```text
   $spacedock:ensign; then Read /tmp/spacedock-dispatch/{name}.md and treat its content as your assignment.
   ```

   Keep the helper-emitted prompt as the value forwarded by `codex_v2_adapter.go`, keep `fork_turns: "none"`, and do not add assignment fields to the runtime adapter. The fixed prefix is bootstrap metadata, not a reconstruction of the assignment.

2. Replace the Codex artifact's false first-action prose with a truthful boundary statement: the outer prompt invokes `$spacedock:ensign`; the skill supplies the shared operating contract; this file supplies the stage-specific assignment. Continue to prohibit Claude `Skill(...)` and `SendMessage` syntax in the Codex artifact.

3. Make the shared contract and Codex runtime reference describe both fresh and reuse paths. The dispatch bootstrap section should show Claude's `Skill(skill="spacedock:ensign"); then Read ...` and Codex's `$spacedock:ensign; then Read ...` forms, and should state that the Codex marker belongs in the outer prompt. The Codex runtime reference should state that fresh prompts load the skill before the pointer and later stage advances use the already-bound worker.

4. Add a small dispatch-boundary fixture. Build a real Codex dispatch from a fixture request, then pass the emitted prompt and artifact to a deterministic child probe with a plugin root copied from the current checkout. The probe must match the exact prefix at byte zero before it reads `skills/ensign/SKILL.md`, `ensign-shared-core.md`, and `codex-ensign-runtime.md`; only then may it read the dispatch artifact and emit a report that the existing Stage Report parser/oracle accepts. It must use only paths supplied by the fixture, not `~/.codex`, network auth, or an ambient `/tmp` file. Include two independent negative mutations: remove the bootstrap edge while retaining the old false sentence, and replace the anchored report with the generic checkbox shape. Keep the existing exact prompt, isolation, and cross-host fixtures as controls; the live spike remains the proof of the real Codex loader, while this fixture makes the dispatch boundary deterministic.

5. On the branch that owns the common live journey, remove only a Codex-specific `full-ensign-cycle` bootstrap quarantine if that exact entry is present, and only after the clean live run passes. Do not create a second journey, remove the current host-neutral `liveDurableJourneyTODO` product-defect skips for `smallest-sufficient-mechanism` or `keep-moving-posture`, touch the Claude-only `TestLiveEnsignCycle`, or alter the shared oracle. Capture the child transcript/report, entity transition, and terminal archive evidence required by the runtime registry; if no Codex-specific quarantine exists at implementation time, make no quarantine change.

## Mechanism choices and alternatives

**Explicit Codex skill invocation — selected.** Value: it uses the current host's progressive-disclosure skill surface to load one authoritative shared core before the fresh child reads the pointer. Simple alternative: keep the pointer-only prompt and rely on the artifact's prose claim. That alternative is the exact failing incident and cannot execute a contract. More elaborate alternative: inline the shared core into every dispatch file. The spike showed the installed skill loads, so inline bytes would duplicate a versioned contract, increase artifact size, and create drift between Claude/Pi/Codex surfaces.

**Truthful artifact boundary — selected.** Value: it makes a missing bootstrap visible and prevents a future reader or model from treating assignment context as the shared protocol. Simple alternative: leave the existing sentence unchanged while fixing only the outer prompt. That would preserve a false diagnostic surface and would let a bootstrap-removal mutation pass if it inspected only the file. The first-action text must therefore identify the outer prompt as the contract edge.

**Fresh-only bootstrap prefix — selected.** Value: it preserves the fresh-context guarantee while keeping reuse prompts small and stable. Simple alternative: prepend `$spacedock:ensign` to every stage advance. That would reload a skill unnecessarily, alter the measured advance prompt, and blur the distinction between fresh worker initialization and same-worker stage transition. The existing advance prompt remains unchanged.

**Dispatch-boundary behavioral fixture plus live journey — selected.** Value: the fixture catches prompt/artifact wiring regressions deterministically, and the live journey proves that a real fresh Codex child follows the loaded contract through First Officer validation and terminalization. Simple alternative: assert that the generated string contains `$spacedock:ensign`. That would not prove skill loading or report behavior, so it is insufficient as the only test.

## Semantic scope

- Command surface: no new CLI command or flag. `dispatch build` changes only the Codex fresh prompt prefix.
- Stored formats: no JSON envelope, state/entity schema, or dispatch-artifact payload fields change.
- Prompt grammar: Codex fresh output gains the fixed `$spacedock:ensign; then Read ` prefix. Claude's prompt, Pi's prompt, and the existing advance prompt remain stable.
- Authority and ownership: the helper still owns the generated prompt; the First Officer still forwards it; the dispatch artifact still owns stage/entity/checklist payload; the Codex runtime adapter still maps it to a named spawn with `fork_turns: "none"`. No First Officer or oracle authority changes.
- Runtime behavior: a fresh Codex child loads `skills/ensign/references/ensign-shared-core.md` through the installed Spacedock plugin before reading the assignment. A reused worker continues from the already-loaded contract.
- Validation: the Stage Report parser, live oracle, report-repair gate, and archive rules are unchanged. Generic checkbox reports remain invalid.
- Host parity: Claude and Pi receive no new Codex marker and keep their current bootstrap/completion behavior.

## Expected surface and estimate

The expected implementation is approximately 6–9 files and 120–220 changed lines, excluding generated dispatch artifacts and live transcript output. A two-times increase requires rechecking scope before proceeding.

Expected production/contract surface:

- `internal/dispatch/build.go`: one host-specific fresh-pointer branch and truthful Codex first-action text, approximately 15–30 lines.
- `internal/dispatch/build_codex_host_test.go` and `internal/dispatch/build_json_ergonomics_test.go`: replace the old no-`Skill(...)` expectation with the exact Codex `$spacedock:ensign` prefix while retaining Claude-syntax and completion-signal negatives.
- `internal/dispatch/codex_bootstrap_test.go` (or the equivalent existing dispatch-boundary test file): the supplied-plugin-root child probe and independent bootstrap/report removal mutations, approximately 80–150 test lines.
- `skills/ensign/references/ensign-shared-core.md`: document the host-specific fresh bootstrap forms, approximately 10–20 lines.
- `skills/ensign/references/codex-ensign-runtime.md`: document the fresh-versus-advance boundary, approximately 5–10 lines.
- The existing common-journey files on the live-journey branch, likely `internal/ensigncycle/shared_scenarios_test.go` and its focused unit test: clear only a Codex-specific `full-ensign-cycle` bootstrap quarantine if the implementation branch contains one and evidence is available. Do not modify the host-neutral product-defect quarantine or add a second journey. Update the runtime registry only if that branch's existing entry requires an evidence pointer; do not expand the registry for this task.

No changes are expected in `cmd/spacedock/`, `internal/dispatch/codex_v2_adapter.go`, the First Officer's assignment reconstruction, the Stage Report oracle, or the sibling task's launcher/artifact payload work. If the sibling lands first and moves a helper, preserve these semantics and re-anchor the test at its new call site.

## Required documentation diff

This is a user-visible host integration behavior, so the implementation must update the runtime support documentation. The intended `docs/runtime-support.md` change is:

```diff
- Codex fresh dispatch is always named, so `host=codex` rejects `bare_mode`.
+ Codex fresh dispatch is always named, so `host=codex` rejects `bare_mode`. Its
+ fresh prompt begins `$spacedock:ensign; then Read ...`, which loads the
+ installed Spacedock ensign skill before the child reads the dispatch pointer.
+ Reuse/advance prompts do not repeat this prefix because they target the worker
+ that already loaded the shared contract.
```

The canonical runtime-host guide carries the same sentence and gets the matching concrete change:

```diff
- Codex fresh dispatch is always named, so `host=codex` rejects `bare_mode`. Keep entity and worktree paths explicit, especially for split-root workflows (`state: .spacedock-state`). Test the positive shape and the banned-tool negative case.
+ Codex fresh dispatch is always named, so `host=codex` rejects `bare_mode`. Its fresh
+ prompt begins `$spacedock:ensign; then Read ...`, which loads the installed
+ Spacedock ensign skill before the child reads the dispatch pointer. Reuse/advance
+ prompts do not repeat this prefix because they target the worker that already
+ loaded the shared contract. Keep entity and worktree paths explicit, especially
+ for split-root workflows (`state: .spacedock-state`). Test the positive shape and
+ the banned-tool negative case.
```

This applies to `docs/site/contributing/adding-a-runtime.md`; no live-CI wording change is required because `docs/runtime-live-ci.md` already documents `spacedock codex --plugin-dir <checkout>` as the current-checkout plugin boundary. If that launch instruction changes during implementation, update that paragraph explicitly. The documentation must not promise that `/tmp` is readable in arbitrary sandboxed sessions; it describes the supported launcher/live-lane boundary.

## Test plan

1. Re-run the focused dispatch tests after the prompt change, including exact Codex prefix, truthful artifact wording, no Claude syntax in Codex artifacts, prompt identity, named spawn, `fork_turns: "none"`, and unchanged advance size.
2. Run the new fixture with the supplied checkout plugin root. The fixed prompt must make the probe read the real core and produce a compliant anchored report; removing only the prefix must fail before assignment execution; replacing only the report with the generic checkbox shape must fail the report oracle independently.
3. Run the existing Codex `full-ensign-cycle` live lane from a clean exact candidate with the supported plugin checkout. Preserve durable evidence for skill load, compliant report, validation, state transition, and archive. Do not treat the unrelated host-neutral product-defect quarantines as evidence for or against this bootstrap.
4. Run focused ensign-contract and cross-host tests, then `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
5. Review the final diff for unchanged frontmatter, unchanged First Officer/oracle semantics, no sibling-task absorption, and the required runtime documentation diff. Commit only this entity path in the state checkout and push `origin spacedock-state/dev`.

## Stage Report: ideation

- DONE: Define an executable Codex bootstrap approach and exercise its riskiest fresh-child path.
  Evidence: AC-2 — the supported Codex fresh-child probe (Codex thread `019fdea2-490f-7632-a599-f380366947b3`, CLI `0.146.0`) exited 0 after reading the real ensign skill/core/adapter and the pointer, then confirmed the anchored Stage Report protocol; the result selects that executable edge over duplicated contract bytes.
- DONE: Make every acceptance criterion independently checkable while preserving host and transport boundaries.
  Evidence: AC-3 transport-boundary evidence: AC-1 through AC-4 each name an independent live, fixture, mutation, or relational check; the repaired fixture contract supplies its plugin root and avoids ambient auth/`/tmp`; the scope preserves `fork_turns: "none"`, byte-for-byte prompt forwarding, pointer-owned assignment payload, Claude/Pi behavior, and the existing oracle.
- DONE: Declare expected surface, estimate, semantic scope, and any required documentation diff.
  Evidence: the expected files, 6–9-file/120–220-line estimate, explicit command/stored-format/authority/prompt/runtime/validation scope, sibling-task boundary, and concrete diffs for both runtime-support docs are recorded above; `runtime-live-ci.md` is intentionally unchanged because its plugin-boundary wording already matches.

### Summary

The current Codex host can load the installed Spacedock ensign skill at the fresh-dispatch boundary, as confirmed by the fresh-child probe. Implement the small outer-prompt prefix, make the artifact wording truthful, add the supplied-root dispatch-to-child negative fixture, and prove the result with the existing Codex live journey; do not inline the shared contract, remove unrelated product-defect quarantines, or change First Officer validation.
