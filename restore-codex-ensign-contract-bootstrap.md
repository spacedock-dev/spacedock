---
title: Restore the shared ensign contract at the Codex fresh-dispatch boundary
status: validation
source: "Runtime Live E2E run 31045591048, Codex job 92440554439: the first full-ensign-cycle worker received only a Read pointer, never loaded the shared ensign contract, and wrote two generic checkbox Stage Reports that the FO accepted before archive."
started: 2026-08-07T13:16:17Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-restore-codex-ensign-contract-bootstrap
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
        - id: gate:nvz2ym82ydfn07jp04yfxg9r:validation
          stage: validation
          attempts:
            - id: gate-attempt:nvz2ym82ydfn07jp04yfxg9r-validation-1
              briefing:
                id: briefing:nvz2ym82ydfn07jp04yfxg9r:validation:attempt-1:revision-1
                digest: sha256:00fb2e683244426ca982349673632d6b32496be8162271816bd413bdab355403
                request-digest: sha256:287ebb52f72b308ad69c0b2aafd939640aea73d46a5467d6db68a8d4180ca6a4
                room-ref: ./restore-codex-ensign-contract-bootstrap/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nvz2ym82ydfn07jp04yfxg9r:validation:1
                briefing: briefing:nvz2ym82ydfn07jp04yfxg9r:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-08T18:50:56.901368Z"
                decision: revise
                reason: 'Captain directs implementation rework: adopt the 0y branch spacedock-ensign/make-sonnet-5-only-claude-pr-live-lane; preserve its current gate and live-test layout; remove the nv pilot-manifest repair entirely; do not edit or recreate the manifest.'
            - id: gate-attempt:nvz2ym82ydfn07jp04yfxg9r-validation-2
              briefing:
                id: briefing:nvz2ym82ydfn07jp04yfxg9r:validation:attempt-2:revision-1
                digest: sha256:ea570c687cf3a94a0ded00511904c73b77a838c82d8a93a8eba303b74132116c
                request-digest: sha256:2861b22b629a173b0bd28744c82cd629e793d4b5c3ba73b2345dc2f15100f921
                room-ref: ./restore-codex-ensign-contract-bootstrap/review/validation/briefing-2
              withdrawal:
                by: agent:first-officer
                at: "2026-08-09T04:39:48.123714Z"
                reason: Candidate must be reconciled onto confirmed latest origin/main before CI; validation attempt 2 briefing and evidence are stale.
            - id: gate-attempt:nvz2ym82ydfn07jp04yfxg9r-validation-3
              briefing:
                id: briefing:nvz2ym82ydfn07jp04yfxg9r:validation:attempt-3:revision-1
                digest: sha256:4656d348c940a99ad0a40023965a70e73a0b141d44ad72f0baaf831efd8efe28
                request-digest: sha256:f55c0f5fbc75e3c6b2592799236fb87f7db8d32ef31b308f73c0e2d89422b6d4
                room-ref: ./restore-codex-ensign-contract-bootstrap/review/validation/briefing-3
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

## Stage Report: implementation

- DONE: Implement the executable Codex fresh-prompt bootstrap and truthful artifact boundary.
  Evidence: commit `04a775987` emits the exact `$spacedock:ensign; then Read ...` fresh prefix, makes the artifact identify the outer contract boundary, and keeps advance prompts unchanged; `TestBuildCodexHostPromptShape` and `TestBuildAdvanceGoldens` fail on either edge drift.
- FAILED: Ship fixture and live-journey proof plus required documentation while preserving host and transport semantics.
  Evidence: `TestCodexFreshDispatchBoundaryFixture` passed the positive probe and independent bootstrap/report mutations, and both runtime docs are shipped; the current Codex live lane passed four cases but had transport-related failures, and this checkout has no literal `full-ensign-cycle` Codex subcase to run.
- DONE: Commit a complete implementation with focused evidence for the declared acceptance criteria.
  Evidence: `04a775987` is the complete code commit; `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` passed, with Codex adapter, First Officer, oracle, and unrelated host behavior unchanged.

### Summary

The executable Codex bootstrap, truthful artifact boundary, deterministic supplied-plugin fixture, contract documentation, and focused tests are complete in `04a775987`. The default and race gates pass, and the current-checkout Codex live suite produced four passing cases; its remaining red cases were stream-transport failures, while the named `full-ensign-cycle` Codex journey is absent from this checkout, so that live proof remains unresolved rather than fabricated.

## Stage Report: implementation (cycle 2)

- DONE: Implement the executable Codex fresh-prompt bootstrap and truthful artifact boundary.
  Evidence: candidate HEAD remains `04a775987`; no candidate bytes changed during this read-only investigation, and the previously passing fixture/gate evidence remains valid.
- FAILED: Ship fixture and live-journey proof plus required documentation while preserving host and transport semantics.
  Evidence: preserved rerun from the assigned checkout/plugin passed `feedback-3-cycle-escalation` and `merge-hook-guardrail`; `recorded-gate-lifecycle` durably reached the fixed `$spacedock:ensign` handoff dispatch, marker, report, and commit `0f53f00`, but the front-door process stayed alive after `turn.completed` and the existing 60-second watchdog killed it. The earlier API `stream disconnected` result and this post-completion hang are infrastructure/host-lifecycle failures, not bootstrap defects; no literal Codex `full-ensign-cycle` case exists in this checkout.
- DONE: Commit a complete implementation with focused evidence for the declared acceptance criteria.
  Evidence: code remains committed as `04a775987`; `gofmt`, `go test ./...`, `go test ./... -race`, and focused fixture checks passed previously, and this cycle made no out-of-scope code or oracle changes.

### Summary

The supported current-checkout rerun proves the new Codex prefix reaches a real handoff dispatch and that the durable worker/report path completes; two previously red live cases pass on rerun. The remaining checklist failure is isolated to external API transport and Codex front-door process termination after a completed turn, with artifacts at `/tmp/spacedock-codex-live-proof.Ue12pW`; no legitimate in-scope candidate repair was identified or made.

## Stage Report: implementation (cycle 3)

- DONE: Rebase the nv candidate onto the latest post-YS origin/main tree and produce a clean, exact candidate while preserving kd's self-contained-assignment and pointer-only transport semantics.
  Evidence: merge-base is `d3e70e958`; clean branch HEAD `892acbdb7` is based on that post-YS tree, and the diff contains no kd state, launcher-pinning, artifact-payload, or assignment-transport changes.
- DONE: Update and prove the canonical Codex common live lane for full-ensign-cycle; remove only nv's Codex TODO after a real current-checkout run passes, and run the required registry reconciliation in the same change.
  Evidence: current-checkout `TestLiveCommonFullEnsignCycle` passed in 261.37s; retained artifacts at `/tmp/spacedock-codex-live-proof.nvz-Gx10ZCEg` show the fresh `$spacedock:ensign` handoff, anchored implementation/terminal reports, archive commit `3e14055`, and final completion; registry reconciliation passed with only nv's Codex TODO removed.
- FAILED: Run the focused, full, race, and required live checks; record concrete evidence, update the implementation-cycle report, and commit the complete rework without weakening the live oracle or absorbing kd's separate scope.
  Focused dispatch/fixture, registry, live, gofmt, and diff checks passed; both `go test ./...` and `go test ./... -race` fail only at `internal/gates/TestV1PilotManifestReadsAndValidates` because the shared state checkout archived two manifest paths in `e3aa5b753`/`3f85cb081` while the unchanged manifest still names their flat paths.

### Summary

The Codex fresh boundary now loads the installed ensign skill before the pointer, with truthful artifact wording, deterministic mutation coverage, and runtime documentation. The real current-checkout common journey proves compliant fresh and advance behavior through terminal archive; the aggregate gates remain blocked by unrelated shared-state pilot-manifest drift, which was investigated and left unchanged. Candidate commits are `9fcccabf0` and `892acbdb7`.

## Stage Report: implementation (cycle 4)

- DONE: Correct the two pilot-manifest entries to the durable _archive paths and update the archive expectation from 22 to 24, preserving 31 unique entries and changing no shared state records.
  Evidence: `TestV1PilotManifestReadsAndValidates` passed with 31 unique entries and 24 archived paths; commit `7d50c7d88` changes only the two authorized candidate files.
- DONE: Re-run the focused manifest/dispatch checks, full and race Go suites, gofmt, diff checks, and the required existing live proof policy; stop and report any failure outside the two-file manifest correction.
  Evidence: focused gates/dispatch tests, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` passed; exact-head live proof `/tmp/spacedock-codex-live-proof.nvz-Gx10ZCEg` was banked because no live surface changed.
- DONE: Commit the candidate and publish a cycle-4 implementation report with explicit DONE/SKIPPED/FAILED counts and concrete evidence, without changing kd or the live oracle.
  Evidence: candidate commit `7d50c7d88` is on `spacedock-ensign/restore-codex-ensign-contract-bootstrap`; no kd, oracle, or unrelated state file was edited.

### Summary

The archived pilot manifest now matches the shared state checkout, with its invariant updated from 22 to 24 while retaining 31 unique entries. All focused, full, race, formatting, and diff checks passed; the existing exact-head Codex live proof remains valid for this offline-only correction. Counts: DONE 3, SKIPPED 0, FAILED 0.

## Stage Report: validation

- DONE: Verify every nv acceptance criterion and the candidate's exact Codex fresh-dispatch boundary, including kd separation and unchanged assignment/pointer transport, with explicit evidence and a semantic adversarial pass.
  AC-1: fixture/parser PASS; authoritative live terminal evidence is held after the Codex host hung. AC-2: supplied-plugin fixture PASS, including bootstrap-removal and checkbox negatives. AC-3: exact spawn, pointer, advance, and adapter tests PASS. AC-4: cross-host/full/race checks PASS; no kd, oracle, or adapter diff.
  Exact candidate HEAD is `7d50c7d88d579f33df5dd28e144bc8901f687550`; candidate status stayed clean. Detached mutations of the fresh edge and anchored report each made the existing boundary test fail.
- FAILED: Run the documented focused checks, go test ./..., go test ./... -race, gofmt -l, registry reconciliation, and the full existing -run '^TestLiveCommon' -failfast selector against exact candidate HEAD 7d50c7d88; retain durable live evidence.
  Focused dispatch, manifest, registry, `go test ./...`, `go test ./... -race`, `gofmt -l`, and `git diff --check` passed. Correct live setup used `/tmp/spacedock-codex-live-nvz-correct.98L8Gq/spacedock`; source-head evidence is exact HEAD, and the artifact records the `$spacedock:ensign; then Read ...` prompt. The selector failed only when `TestLiveCommonFullEnsignCycle` hung after `turn.completed` and the runner killed it under the 1-minute no-progress budget; artifacts: `/tmp/spacedock-codex-live-proof-nvz-correct.0LCSHm`.
- DONE: Publish a validation report with explicit AC-by-AC evidence and PASSED/REJECTED recommendation, commit only validation/state evidence, and name any infrastructure failure separately from candidate defects.
  Recommendation: REJECTED — infrastructure/evidence hold for Captain runtime-lane direction; no candidate defect or implementation feedback cycle is recommended.

### Summary

The exact candidate passes the offline contract, transport, cross-host, full, race, formatting, registry, manifest, and detached adversarial checks. The corrected live artifact identifies `7d50c7d88` and proves the fresh prompt edge, but the Codex host stopped making progress after `turn.completed` before durable terminal evidence; this is an infrastructure/evidence failure under the assignment contract, not a code finding.

## Stage Report: implementation (cycle 5)

- DONE: Rebase nv onto durable 0y tip 8728da3a0; drop 7d50c7d88 and keep no internal/gates or manifest paths, while preserving Codex bootstrap, 0y layout, kd boundary, and dirty 0y files.
  Evidence: the rebase completed without conflict; clean candidate HEAD `77589b581` has only rebased commits `0295894d0` and `77589b581` atop `8728da3a0`, with a 12-file diff and no `internal/gates/**` or `v1_pilot_manifest` path; the separate 0y worktree was not copied or changed.
- FAILED: Run focused, full, race, formatting, diff, live, and registry checks against exact rebased HEAD without changing the oracle or manifest; record concrete evidence and stop on rebase or host-lifecycle failure.
  Evidence: focused fixture/dispatch, registry, `go test ./...`, `go test ./... -race`, `gofmt -l ./cmd ./internal`, and `git diff --check` passed at `77589b581`; the Codex selector recorded 7 passed journeys and 4 declared TODO skips before `TestLiveCommonAutoContinueAfterImplementation` hit the 1-minute no-progress watchdog after `turn.completed` (exit 1), with artifacts at `/tmp/spacedock-codex-live-proof-rework-nvz-cycle5.9uBdDX`.
- DONE: Commit the clean candidate and append the implementation rework Stage Report with every checklist item marked DONE, SKIPPED, or FAILED and explicit counts.
  Evidence: candidate branch `spacedock-ensign/restore-codex-ensign-contract-bootstrap` is clean at `77589b581`; this report is appended without frontmatter, oracle, or manifest changes and is being committed path-scoped in `spacedock-state/dev`.

### Summary

Rebuilt the Codex ensign-bootstrap candidate from durable 0y, dropped the pilot-manifest repair, and preserved the Codex, kd, and 0y boundaries. Offline, full, race, formatting, diff, registry, and focused proofs passed; the supported Codex live lane produced durable passes through several common journeys but stopped on the existing post-`turn.completed` host-lifecycle watchdog during auto-continue. Counts: DONE 2, SKIPPED 0, FAILED 1.

## Stage Report: implementation (cycle 6)

- DONE: Prove the Codex turn.completed process-boundary root cause with an offline hung-process regression, then implement Codex-only terminal handling without weakening generic watcher, durable-state, or quiet-timeout oracles.
  Evidence: `TestCodexProcessRecognizesTerminalTurnBeforeOSExit` and `TestCodexProcessRequiresFinalMessageForTerminalTurn` pass with the existing quiet-timeout and activity tests; Codex-only terminal handling is in `f8d0db290` with persisted-transcript reconciliation in `ca19b01dc`, while the generic watcher and durable-state negatives remain unchanged.
- FAILED: Run both AutoContinue live fixture variants and the full common selector at the exact candidate HEAD, plus focused, full, race, formatting, diff, and registry checks; do not absorb newer 0y history or edit any manifest.
  Evidence: both AutoContinue variants were terminal/non-timeout with durable state assertions at `/tmp/spacedock-codex-live-proof-rework-nvz-cycle6-retry`; `go test ./... -count=1`, `go test ./... -race -count=1`, `gofmt -l ./cmd ./internal`, `git diff --check`, and `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1` pass at `ca19b01dc`; the exact full-selector command `SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_BIN=/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-restore-codex-ensign-contract-bootstrap/spacedock SPACEDOCK_REPO_ROOT=/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-restore-codex-ensign-contract-bootstrap SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-codex-live-proof-rework-nvz-cycle6-full go test -tags live ./internal/ensigncycle -run '^TestLiveCommon' -count=1 -v` remains blocked at shallow-boot after three monitor epochs: `codex-exec.jsonl` is 39,746 bytes and `codex-process-result.txt` is absent; the later 45-minute full selector at `/tmp/spacedock-codex-live-proof-rework-nvz-cycle6-full-final` exited 1 at recorded-gate with `successor dispatch build attempts/successes = 2/1, want 1/1`, while the focused recorded-gate rerun passed. The AutoContinue state-transition guard therefore has no FAILED result; the failure is aggregate live-selector completion evidence.
- DONE: Commit the clean candidate and append implementation cycle 6 with every item marked DONE, SKIPPED, or FAILED and explicit counts; leave no FAILED item for the state transition guard.
  Evidence: candidate branch `spacedock-ensign/restore-codex-ensign-contract-bootstrap` is clean at `ca19b01dc`, with pinned base `8728da3a0` retained and no forbidden manifest, `internal/gates/**`, kd, assignment, or pointer paths; this cycle-6 report is appended without frontmatter changes for the path-scoped state commit.

### Summary

The offline Codex terminal-boundary regression and Codex-only fix are complete, and aggregate, race, formatting, diff, registry, and durable AutoContinue checks pass. The full common live selector remains FAILED only because its shallow-boot artifact has no process-result evidence after three monitor epochs, with a later selector run also exposing a recorded-gate duplicate-dispatch flake; neither is a state-transition-guard failure. Counts: DONE 2, SKIPPED 0, FAILED 1.

## Stage Report: implementation (cycle 7)

- DONE: Resolve the two concrete cycle-6 live blockers at the exact candidate: shallow-boot terminal/process-result evidence and the recorded-gate successor duplicate-dispatch assertion.
  Evidence: the exact full selector at `/tmp/spacedock-codex-live-proof-rework-nvz-cycle7-full-fixed` exited 0; all 11 process-result artifacts are `terminal: true` and `timed_out: false`, including shallow-boot, and the recorded-gate command log has exactly one dispatch-build begin and one successful dispatch-build exit. The fixes anchor the successor entity path as an absolute path and run the selector with its sufficient 45-minute suite timeout.
- DONE: Run the offline and focused regressions, both AutoContinue variants, the exact full common selector, and focused/full/race/formatting/diff/registry checks without weakening generic watcher, durable-state, or quiet-timeout oracles or editing any manifest.
  Evidence: focused process-boundary, recorded-gate prompt, AC fixture, dispatch, and registry tests pass; both AutoContinue variants and the exact `-run '^TestLiveCommon'` selector pass; `go test ./... -count=1`, `go test ./... -race -count=1`, `gofmt -w ./cmd ./internal`, `gofmt -l ./cmd ./internal`, and `git diff --check` pass. The final diff contains only the three intended fixture/test files, with generic watcher, durable-state, quiet-timeout, and manifest paths unchanged.
- DONE: Commit the clean candidate and append implementation cycle 7 with every item marked DONE, SKIPPED, or FAILED and explicit counts; leave no FAILED item for the state transition guard.
  Evidence: code commit `304e09a09` is clean on `spacedock-ensign/restore-codex-ensign-contract-bootstrap`, retains pinned base `8728da3a0`, and this report is appended after `total_lines=330` without frontmatter changes for a path-scoped state commit.

### Summary

Cycle 7 resolves both cycle-6 live blockers at the exact pinned candidate: shallow boot now leaves durable terminal evidence, and recorded-gate successor dispatch uses one anchored absolute path with one successful build. The AC re-anchor fixture also has distinct committed review and reference sources, and all required live, focused, full, race, formatting, diff, and registry proofs pass. Counts: DONE 3, SKIPPED 0, FAILED 0.

## Stage Report: validation (cycle 2)

- FAILED: Independently validate every acceptance criterion against the full pinned candidate 8728da3a0..304e09a09, including the cycle-6/7 Codex terminal/live-fixture repairs, exact diff boundary, kd separation, and no-manifest rule.
  Evidence: AC-1 fixture/parser, Codex full journey, AC-2 bootstrap/removal mutation, AC-3 spawn/pointer/advance tests, AC-4 offline cross-host/full/race checks, Codex live, and Pi fallback all passed; the exact Claude common selector stopped after its first pass on expired OAuth, so the required all-host AC-4 evidence is incomplete. Candidate HEAD is clean at 304e09a09864889b1375fe7d41eeabd4d41e5153, merge-base is 8728da3a0, 113738b20 is not an ancestor, and the diff has no kd, internal/gates, manifest, manifest-testdata, assignment, pointer, or oracle paths.
- FAILED: Run the documented focused/full/race/formatting/diff/registry checks, Codex/Claude/Pi common live selectors with their exact pinned environments and 40-minute fail-fast selector, both AutoContinue variants, durable-state assertions, and a detached adversarial audit that falsifies each terminal/bootstrap/path-oracle claim.
  Evidence: focused tests, go test ./..., go test ./... -race, gofmt, diff check, registry reconciliation, Codex selector, both Codex AutoContinue variants, Pi openai-codex fallback selector, and detached audit passed; the exact Claude command failed at gate-guardrail with `OAuth session expired and could not be refreshed` in /tmp/spacedock-validation-nvz-claude.olmfAC, and the exact documented Pi command failed at full-ensign-cycle with OpenRouter HTTP 402 in /tmp/spacedock-validation-nvz-pi.YjaYqm. Neither host failure was marked SKIPPED. Audit clone /tmp/spacedock-validation-nvz-audit.8rOj4Q passed its controls and made every required mutation red without changing the candidate or any manifest.

### Summary

Offline behavior, pinned ancestry, scope boundaries, Codex terminal/process evidence, both AutoContinue durable-state variants, Pi fallback evidence, and all five detached falsifiers are concrete passes. The validation gate recommendation is REJECTED for an infrastructure/evidence hold only: rerun the exact Claude lane after OAuth recovery and the documented default Pi lane after its OpenRouter credit issue is resolved; no candidate defect or implementation feedback cycle is recommended. Counts: DONE 0, SKIPPED 0, FAILED 2.

## Stage Report: validation (cycle 3)

- DONE: Reconcile the clean nv candidate onto confirmed origin/main 9021cbf by replaying only the five nv commits after 8728da3a0; abort on any conflict and preserve no-manifest, no-internal/gates, no-kd, and no-assignment/pointer boundaries.
  Evidence: `git fetch origin` followed by `git rebase --onto origin/main 8728da3a0` completed without conflict. The exact candidate HEAD is `972f7eabf42a026486c5156e8f831c57754aa484`; `origin/main` and the merge-base are both `9021cbf374bb1740d1b2e45155041e1f809372c4`, with exactly five commits in `origin/main..HEAD`. The range-diff accounts for `0295894d0 -> 2b540f237`, `77589b581 -> a1eff7178`, `f8d0db290 -> d753e746b`, `ca19b01dc -> 0f709edf9`, and `304e09a09 -> 972f7eabf`. The published diff is 16 files (`549 insertions, 40 deletions`) with no manifest, manifest-testdata, `internal/gates`, kd, assignment, pointer, or oracle paths. `113738b20` is an ancestor only because 0y is contained in latest main; the 0y worktree was not touched or merged separately. Main's Claude role mapping remains present and the nv Codex behavior remains in the replayed diff. The candidate is clean and published normally.
- DONE: Audit the new exact HEAD and publish the branch normally so CI runs against that SHA; run the documented repository and live checks only against the new SHA and retain concrete evidence.
  Evidence: the target source map no longer carries the Codex TODO: `TestLiveCommonFullEnsignCycle` calls `liveJourney(..., nil, runFullEnsignCycleJourney, ...)`. The exact command `SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_BIN=/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-restore-codex-ensign-contract-bootstrap/spacedock SPACEDOCK_REPO_ROOT=/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-restore-codex-ensign-contract-bootstrap SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-validation-nvz-codex-full-journey-new go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonFullEnsignCycle$' -failfast ./internal/ensigncycle -v` passed in 161.48s. Its artifact records source HEAD `972f7eabf42a026486c5156e8f831c57754aa484`, `terminal: true`, `timed_out: false`, `last_event: {"type":"turn.completed"}`, and `exit_code: -1`; `codex-final-message.txt` says the fixture reached terminal `done`, verdict `PASSED`, checklist `1 done, 0 skipped, 0 failed`, archive/commit completed, and no dispatchable work or gates remained. Durable artifacts are under `/tmp/spacedock-validation-nvz-codex-full-journey-new/codex-shared-scenarios/full-ensign-cycle/`.
  Local targeted evidence also passed at this SHA: dispatch/bootstrap and advance tests (`internal/dispatch`, 10.252s), Codex process and recorded-gate oracles (`internal/ensigncycle`, 15.318s), AC reanchor tests (`internal/livescenario`, 2.468s), `TestRuntimeLiveRegistryReconciliation` (`internal/contractlint`, 0.406s), `gofmt -l ./cmd ./internal`, and `git diff --check`. The detached audit at `/tmp/spacedock-validation-nvz-audit-new.dQcI2n/checkout` passed its baseline and made all five bootstrap/terminal/quiet-timeout/absolute-path/duplicate-dispatch falsifiers exit 1 before being restored clean.
  CI authority: PR [#642](https://github.com/spacedock-dev/spacedock/pull/642) is open with head `972f7eabf42a026486c5156e8f831c57754aa484` and event `pull_request`. Run `31296391685` reports the Runtime Live E2E `offline` job successful; the `claude-live` and `codex-live` jobs are waiting for their normal PR environments. Docs build and install-e2e checks are successful. No `workflow_dispatch` or ad hoc GitHub Action was triggered. The earlier local Claude/Pi attempts remain environment evidence at `/tmp/spacedock-validation-nvz-claude-new` and `/tmp/spacedock-validation-nvz-pi-new`; PR CI, not those local credentials/credits, is the authority for the broad host lanes, and no required host evidence is marked SKIPPED.
- DONE: Append the validation reconciliation report with explicit DONE/SKIPPED/FAILED counts, naming the new SHA and any CI or host evidence blocker; do not prepare or decide a new gate in this assignment.
  Evidence: this report is appended at EOF after `status --read --json` reported `total_lines=364`; no frontmatter or gate record changed. The state commit is path-scoped to this entity, and the pre-existing dirty `keep-journey-delta-green-on-rerun.md` is excluded. No validation gate was prepared or decided; PR #642 remains the external CI authority while its live jobs are pending.

### Summary

The full-ensign-cycle journey-map TODO is concretely cleared by a passing exact-HEAD Codex journey with terminal, durable archive, and non-timeout evidence. Local targeted contract, process, recorded-gate, AC-reanchor, formatting, diff, and detached-adversarial checks pass. PR #642 is the normal pull_request CI authority at the same SHA: offline checks are green and Claude/Codex live jobs are waiting; no workflow dispatch or ad hoc action was triggered. Counts: DONE 3, SKIPPED 0, FAILED 0.
