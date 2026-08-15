---
id: v7a6xqh2rm3asjvj8qz1y4p0
title: Retire banned prose-grep contract tests and dedupe surviving pins
status: validation
source: "Captain review of the 0.27 stack + audit-r2 (2026-08-15); captain directive: file, dispatch off stack tip, PR as stack layer"
started: 2026-08-15T18:31:08Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-retire-prose-grep-contract-tests
issue:
gates:
    version: 1
    records:
        - id: gate:v7a6xqh2rm3asjvj8qz1y4p0:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:v7a6xqh2rm3asjvj8qz1y4p0-backlog-1
              briefing:
                id: briefing:v7a6xqh2rm3asjvj8qz1y4p0:backlog:attempt-1:revision-1
                digest: sha256:db3d3e2c2eabd9426e3eeec0af2179e2124792c069cecfb280ae50c060165f3e
                request-digest: sha256:799382d8dca6e2b600e6141353217b02e347f0288e2749291aedb9d75fb1c687
                room-ref: ./retire-prose-grep-contract-tests/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v7a6xqh2rm3asjvj8qz1y4p0:backlog:1
                briefing: briefing:v7a6xqh2rm3asjvj8qz1y4p0:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T18:15:14.081321Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: file, dispatch based off stack tip, PR on top of the stack'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:v7a6xqh2rm3asjvj8qz1y4p0:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:v7a6xqh2rm3asjvj8qz1y4p0-ideation-1
              briefing:
                id: briefing:v7a6xqh2rm3asjvj8qz1y4p0:ideation:attempt-1:revision-1
                digest: sha256:e05717deab88ae353e7aa7a81776a2f6fbff1be6aebe2276d93f73a0eae11c0f
                request-digest: sha256:98e8895855a32e8afeaf82cb2c6caea45aad62065aa2c912931a0fde213f1055
                room-ref: ./retire-prose-grep-contract-tests/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v7a6xqh2rm3asjvj8qz1y4p0:ideation:1
                briefing: briefing:v7a6xqh2rm3asjvj8qz1y4p0:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T18:30:32.23654Z"
                decision: approve
                reason: 'Captain approved 2026-08-15 with both queued rulings as recommended: install-gate goes checkless per the 2026-07-20 ruling with the sentinel behavior test filed separately; the narrow read-quarantine exemption is adopted'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:v7a6xqh2rm3asjvj8qz1y4p0:validation
          stage: validation
          attempts:
            - id: gate-attempt:v7a6xqh2rm3asjvj8qz1y4p0-validation-1
              briefing:
                id: briefing:v7a6xqh2rm3asjvj8qz1y4p0:validation:attempt-1:revision-1
                digest: sha256:77049f4058ecf77efc91881ccf4bf5eae38457e8335fe92847c2e0ccf8c80803
                request-digest: sha256:72357a0b11c6aee70fd2271fa57485d9eb4fd1821d8670cfa38ea738c219fba2
                room-ref: ./retire-prose-grep-contract-tests/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v7a6xqh2rm3asjvj8qz1y4p0:validation:1
                briefing: briefing:v7a6xqh2rm3asjvj8qz1y4p0:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T19:36:29.344097Z"
                decision: approve
                reason: Captain approved validation 2026-08-15; stacks as layer 10
              application:
                target-stage: done
                state: pending
---

Delete the committed prose-grep tests the Proof policy bans (paraphrase reds them, inversion passes them), dedupe the two double pins, and resolve the four gray cases. Base all work on the stack tip (branch stack27/09-trim-version-output); the deliverable becomes stack layer 10. This MUST land before make-shipped-contracts-self-contained (layer 11), whose prose rewrites would red several of these pins.

Audit inventory (verified at the ship tree):

Banned, delete (9 functions, 4 files):
- internal/contractlint/feedback_rejection_publication_smoke_test.go:9 and :40 (both functions)
- internal/contractlint/initial_worker_spawn_guard_test.go:30, :58, :67 (nine frozen sentences plus a token count)
- internal/contractlint/pi_spawn_binding_test.go:21 only (the :35 structural-absence sibling is ALLOWED - keep)
- internal/contractlint/version_gate_smoke_test.go:35, :77, :158 (keep the :53 gate --help structural-absence guard and TestVersionGateSandboxRegistry / TestVersionGateDeferredTrigger / TestInstallHintNoDrift)

Duplicated pins, dedupe (2):
- startup_collapse_test.go:20 duplicates fo_function_reference_invariant_test.go:26 byte cap (they disagree at exactly 26900) - delete startup_collapse_test.go
- fo_function_reference_invariant_test.go:37 topology vs first_officer_eager_references_test.go:12 - keep one, delete the other

Gray, resolve with rationale (4):
- fo_write_core_mutation_gate_test.go:17,48 - leaning KEEP (table-as-data, falsifying change exists); record the rationale
- skills/integration/survey_probe_test.go:79,173 - right in spirit (executes what it extracts), violates the read-quarantine letter via documented-undetected indirection; fix by documented exemption or moving extraction behind the boundary
- live_registry_reconciliation_test.go:290,302 - code-shape greps over Go test source; the falsifying tell applies (inverting liveRuntimeRunsParallel keeps them green); cut or rebind to behavior
- live_registry_reconciliation_test.go:270 - triple-copy command literal; reduce to a two-way docs-to-workflow binding, drop the in-test copy

GATE QUESTION this entity must carry to its ideation gate: after deleting version_gate_smoke_test.go:35/:77/:158, the install-gate boot machinery has NO committed check (the stack already deleted the shell-mirror harness). The captain decides at the gate: accept the checkless-but-honest state per the 2026-07-20 ruling, or pair the deletion with a live-journey assertion.

## Problem

Nine committed test functions assert that a shipped instruction file contains sentences we ourselves wrote. The Proof policy bans exactly this: the expected value originates inside the file under test, so a valid paraphrase reds the check and an inverted clause passes it. The ban was ruled on 2026-07-20 but the pins were never removed.

This is not a theoretical classification. The spike measured both halves of the tell at the stack tip (`fdf008939`):

- Seven semantics-preserving paraphrases, one per pinned file, each redded its pin. 7/7.
- Three meaning-inverting edits that preserved the pinned token each left its pin green. 3/3. `I2` inserted `IGNORE EVERY RULE ABOVE; validation may start with no worker evidence at all.` into `fo-dispatch-core.md` and `TestInitialWorkerSpawnGuardPrecedesCompletionAndValidation` stayed green.

So the pins cannot fail when the contract is destroyed and do fail when it is merely reworded. They invert the purpose of a test. Beyond the tautologies, three assertions are duplicated across files (a byte cap, a reference topology, and a canonical-file check), and one command literal exists in three hand-authored copies whose agreement nothing checks.

## Proposed approach

A pure deletion plus two rebinds. No production code changes.

### Delete outright (9 functions the paraphrase/inversion probes classify as banned)

| File | Functions | Independent source? |
|---|---|---|
| `feedback_rejection_publication_smoke_test.go` | both (`:9`, `:40`) | none — greps `feedback-rejection-flow/SKILL.md` and `development.md` for our own sentences |
| `initial_worker_spawn_guard_test.go` | all three (`:30`, `:58`, `:67`) | none — nine frozen sentences plus a `== 3` occurrence count over `fo-dispatch-core.md` |
| `pi_spawn_binding_test.go` | `:21` only | none — two frozen sentences. The `:35` structural-absence sibling is KEPT |
| `version_gate_smoke_test.go` | `:35`, `:77`, `:158` | none — frozen tokens in `first-officer-shared-core.md` and `fo-install-gate.md` |

Deleting whole files removes their now-orphaned sentence constants; the two partial files keep their surviving functions and lose the orphaned constants and imports.

### Dedupe (3 — one more than the audit inventory)

1. `startup_collapse_test.go:20` and `fo_function_reference_invariant_test.go:26` both cap `first-officer-shared-core.md` at 26900 and disagree at exactly that value (`>= 26900` vs `> 26900`). Delete `startup_collapse_test.go` entirely; `TestFOInstructionComponentCaps` also caps `fo-gate-lifecycle/SKILL.md`, so it is the superset.
2. `fo_function_reference_invariant_test.go:37` (`TestFirstOfficerReferenceTopology`) and `first_officer_eager_references_test.go:12` both parse `@references/` lines from `first-officer/SKILL.md` and assert the identical want-list. Keep the `first_officer_eager_references_test.go` copy — it additionally asserts each deferred load-point appears exactly once in the inventory block — and delete `TestFirstOfficerReferenceTopology`. Its one non-duplicated assertion (absence of `first-officer/references/fo-smallest-sufficient-mechanism.md`) is dropped rather than folded in: the surviving exact-import assertion already reds if such a file were eagerly imported, and an unimported stray file is dead weight, not a contract breach.
3. **Scope delta, not in the audit inventory.** `first_officer_eager_references_test.go:61` (`TestFirstOfficerDeferredWriteCoreHasSingleCanonicalFile`) is a *third* copy: both its assertions (`fo-write-core.md` non-empty; `skills/fo-write-core` absent) are already made at `:46-51` and `:53-58` of the function above it. AC-2 says topology is asserted exactly once, so delete it too.

### Gray resolutions (4)

1. **`fo_write_core_mutation_gate_test.go:17,48` — KEEP.** The file supplies DATA (a fenced `FO-WRITE-CLASSIFIER` pattern table); the test supplies fifteen independent path expectations. Moving `docs/dev/README.md` from `allowed-process` to `blocked-product` in the table reds it. Recorded bound: `classifyFOWriteTarget`/`pathPatternMatches` are a Go reimplementation of gate semantics the FO actually executes by reading prose, so this proves the table classifies as intended, not that the FO obeys it. That is honest and worth keeping; it is not a behavioral proof of the gate.
2. **`skills/integration/survey_probe_test.go:79,173` — KEEP, and amend the policy.** These are the anti-prose-grep pattern: they extract the shipped bash probe from `survey/SKILL.md` and EXECUTE it against fixture conditions, with the oracle being on-disk state. They violate only the letter of the read quarantine ("tests do not read instruction files except in `internal/contractlint`"). Moving extraction behind a boundary was rejected under trace-mechanism-to-value: it adds a mechanism and changes nothing about what is asserted. Moving the tests into `contractlint` is self-contradictory, since that package is quarantined to *structural* checks. The fix is a documented exemption. **Requires captain approval — see the gate section.**
3. **`live_registry_reconciliation_test.go:290,302` — CUT both.** `TestRuntimeLiveCodexParallelCapacityAndIsolation` greps `shared_live_runner_test.go` for the literal call-site text `liveRuntimeRunsParallel(os.Getenv("SPACEDOCK_LIVE_RUNTIME"))`; inverting that function's body leaves the call site unchanged and the test green. `TestGateStopRunnerDoesNotShortCircuitBoundAssertion` greps a function body for `if err := assert(before, after, expected);`. Both assert code shape over Go *test* source, which is the prose-grep tell in a different file extension. No replacement unit test: `liveRuntimeRunsParallel` sits behind `//go:build live` in `internal/ensigncycle`, and the property it guards (Codex admitted to `t.Parallel`) is already observable as the lane's own 40m wall-clock. Residual recorded: the gate-stop short-circuit property loses its only committed check, and a live journey that silently stops asserting would not red — noted for the captain, not patched here.
4. **`live_registry_reconciliation_test.go:270` — REBIND to a two-way binding.** Today the common-suite command exists as three hand-authored copies (workflow YAML, `docs/runtime-live-ci.md`, and the test's own literal) with the test asserting each copy matches its own literal, so the three sources are never compared to each other. Replace with: extract the command from the workflow per runtime, extract the docs command per runtime, and compare the *run shape* they must share (timeout, parallelism, fail-fast), ignoring the `gotestsum`-vs-`go test` reporting flags that legitimately differ. No expected value is written in the test.

**This rebind found a real defect on its first run.** `docs/runtime-live-ci.md:72` tells a human to run the Claude common suite with `-failfast`, while the workflow deliberately omits it and `TestRuntimeLiveCommonFailFastPolicy` asserts Claude "common journeys must all run before the job reports failure". The docs contradict the policy. The one-line docs fix (drop `-failfast` from the Claude command) is included in this entity; without it the rebind cannot be green. Codex and Pi already agree.

### AC-1 residual: one KEPT test still fails a paraphrase

The audit's keep-list is not sufficient for AC-1 as written. `TestVersionGateDeferredTrigger` (`version_gate_smoke_test.go:53`, kept for its `gate --help` structural-absence guard) also asserts the core contains ``read `references/fo-install-gate.md` ``. Paraphrasing `read` to `load` reds it — measured. Narrow it by deleting that one assertion; the surviving deferred-load-points inventory entry already establishes the reference, and the function keeps three structural checks (inventory entry present, `gate --help` absent, `fo-install-gate.md` resolves at >= 500 bytes). Measured after narrowing: the `read`->`load` paraphrase is green, and renaming the inventory entry still reds it.

### Documentation changes

Two, both one-line.

`docs/runtime-live-ci.md:72` — before:

    SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveCommon' -failfast -parallel 3 ./internal/ensigncycle -v

after:

    SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 90m -run '^TestLiveCommon' -parallel 3 ./internal/ensigncycle -v

`docs/dev/README.md`, Proof policy, read-quarantine bullet (gray 2; **captain decision**) — before:

    - **Instruction-file read quarantine.** Tests do not read prompt or instruction files except in `internal/contractlint`, and there only for structural checks: reference closure, frontmatter validity, structural absence, and dedup. Prose-grep and prose-to-code consistency checks never substitute for running the behavior.

after:

    - **Instruction-file read quarantine.** Tests do not read prompt or instruction files except in two cases. First, in `internal/contractlint`, and there only for structural checks: reference closure, frontmatter validity, structural absence, and dedup. Second, to extract a shipped runnable block and execute it against independent fixture conditions, where the oracle is the fixture's on-disk state or the block's observed output and never the file's wording — `skills/integration/survey_probe_test.go` is the reference shape. Prose-grep and prose-to-code consistency checks never substitute for running the behavior.

### GATE QUESTION: install-gate coverage after the deletion

Deleting `version_gate_smoke_test.go:35/:77/:158` leaves these committed checks on the install-gate boot machinery, all of which name a source outside the file under test:

- `TestVersionGateDeferredTrigger` (narrowed) — `fo-install-gate.md` exists at >= 500 bytes, is named in the deferred inventory, and no inline `gate --help` probe. Source: the filesystem.
- `TestVersionGateSandboxRegistry` — every `insideRegistry` row's env-var name AND value, source-grepped from `internal/safehouse/state.go`, must appear in both files. Source: the binary's own table.
- `TestInstallHintNoDrift` — the curl|sh command and brew tap/formula must equal `docs/site/get-started/install.md`. Source: the published install doc.

What genuinely loses all committed coverage: the sentinel one-attempt loop bound (`test -f <sentinel>`, create-before-run, the per-runtime identity keys, the `rm` recovery message) and the `uname -s` OS-detection and source-build escape wording. There is no live journey covering either — confirmed by search: nothing in `internal/ensigncycle/` or `docs/runtime-live-ci.md` names an install-gate or version-gate journey. The gate fires only when the binary is absent, which never happens in the live lanes, so this is genuinely unobserved rather than incidentally covered.

**Recommendation: accept the checkless-but-honest state, and file the sentinel-bound behavior test as a separate entity.** Three reasons. First, the 2026-07-20 ruling is unambiguous, and keeping `:77` to preserve a coverage number preserves a check that the inversion probes show cannot detect the failure it names — that is worse than no check, because it reads as coverage. Second, pairing a live-journey assertion onto this entity means a new live scenario that boots a runtime with no `spacedock` on PATH and performs a network install; that is a substantial new mechanism riding on a deletion, and it would make layer 10 risky to land ahead of layer 11's rewrites, which is the whole reason this entity is sequenced first. Third, the cheap honest replacement already exists as a proven pattern in this repo: the sentinel bound is a bash mechanism, and `survey_probe_test.go` shows how to extract a shipped bash block and execute it against fixture conditions. That deserves its own design, not a rider.

## Out of scope

The prose files themselves (layer 11 owns rewrites). The allowed tests named above. The replacement sentinel-bound behavior test (recommended as a follow-up entity). Any production (non-test) Go code.

## Expected surface and tolerance

Measured in the spike, not estimated: **9 files, +38 / -321, net -283.**

- 4 test files deleted outright: `feedback_rejection_publication_smoke_test.go`, `initial_worker_spawn_guard_test.go`, `startup_collapse_test.go` (-150)
- 4 test files edited: `version_gate_smoke_test.go` (-76), `live_registry_reconciliation_test.go` (rebind, +36/-34), `pi_spawn_binding_test.go` (-17), `fo_function_reference_invariant_test.go` (-32), `first_officer_eager_references_test.go` (-12)
- 1 docs file: `docs/runtime-live-ci.md` (1 line)
- Plus `docs/dev/README.md` (1 bullet) if the captain approves the gray-2 exemption

Tolerance: +/-2 files and +/-60 lines. The net MUST stay below -250; a smaller reduction means a deletion was skipped.

Semantic changes declared: none to command grammar, stored formats, authority, or binary runtime behavior — no non-test Go file is touched. Two observable non-binary changes: the documented human-facing Claude live command loses `-failfast`, and (pending captain approval) the workflow Proof policy gains a read-quarantine exemption clause. One test-surface semantic change: the install-gate sentinel bound and OS-hint wording lose all committed coverage, per the gate question above.

## Acceptance criteria

**AC-1 (value) - A semantics-preserving paraphrase of any shipped contract sentence leaves the whole suite green, where 7 such paraphrases red it today.**
Verified by: apply all 7 measured paraphrases simultaneously (the P1-P7 set: `fo-dispatch-core.md` empty-wait guard, `feedback-rejection-flow/SKILL.md` "hold the flow", shared-core launcher-resolution and unsupported-OS phrasings, `pi-first-officer-runtime.md` agent-binding sentence, `fo-install-gate.md` "no second install attempt", `development.md` publication-ordering phrase), run `go test ./internal/contractlint/ ./skills/integration/`, require green, revert. The independent baseline that can move the wrong way is the count of paraphrases that red: 7 at `fdf008939`, and the criterion demands 0.

**AC-2 - The byte cap and the reference topology are each asserted exactly once, and every surviving pin in a touched file names a source that can diverge from the file under test.**
Verified by: the three deduped functions absent from the tree; plus one demonstrated falsifying edit per surviving pin class, each shown to red — a second `@references/` import for the topology pin, a renamed deferred-inventory entry for `TestVersionGateDeferredTrigger`, a workflow-side timeout bump for the rebound common-suite binding, and an `insideRegistry` value change for the sandbox pin.

**AC-3 (value) - The change removes more lines than it adds, and the suite is green plain and under -race.**
Verified by: `git diff --stat` net below -250 against the stack-tip base `fdf008939`; `go test ./internal/contractlint/ ./skills/integration/` and the same with `-race`, both green; `gofmt -l ./cmd ./internal` empty.

## Test plan

No new test infrastructure. The plan is the deletion plus the two rebinds, with the existing allowed lints as the regression floor.

- **AC-1 probe harness:** a throwaway script that applies the 7 paraphrases, runs the two packages, and reverts. Already written and exercised in the spike; implementation reuses it and pastes the before/after verdicts. Cost: minutes. Not committed — per the 2026-07-20 ruling, a paraphrase probe is one-off validation evidence, never a committed test.
- **AC-2 falsifying edits:** four one-line edits, each applied, run, and reverted. All four verified RED in the spike.
- **AC-3:** `go test` plain and `-race` over `./internal/contractlint/` and `./skills/integration/`, plus `gofmt`. The spike measured both green after the full deletion set.
- **CI lanes:** the diff touches only `_test.go` files, `docs/runtime-live-ci.md`, and `docs/dev/README.md`. It touches no shipped contract text, no host adapter, and no dispatch/launch path, so the deterministic lanes suffice under the required-CI-lanes rule. The rebound binding reads `runtime-live-e2e.yml` but does not change it.
- **Spike status:** no unverified mechanism remains. Every claim in this design was exercised at `fdf008939` in a detached worktree — the ban classification (7 paraphrase reds, 3 inversion greens), the post-deletion suite (green plain and `-race`), the AC-1 combined proof (green), the AC-1 residual and its narrowing (red, then green, still falsifiable), the rebind (red on the real drift, green after the docs fix, red on a workflow timeout bump), and the topology dedupe (surviving pin red on a second eager import).

## Stage Report: ideation

- DONE: Design confirms every deletion target and gray resolution against the stack tip; the gate question (install-gate coverage after deletion: checkless-per-policy or paired live-journey assertion) is presented with a recommendation
  All 9 banned functions and both dedupe pairs verified present at `fdf008939` with the audited line numbers exact; 4 grays resolved (KEEP fo-write-core mutation gate; KEEP survey probe + policy exemption; CUT both live code-shape greps; REBIND the triple-copy literal); gate question answered with a recommendation to accept checkless and file the sentinel test separately.
- DONE: Value AC: paraphrase probes red at tip and green after; every surviving pin names its independent source
  At tip: 7/7 semantics-preserving paraphrases RED their pins and 3/3 meaning-inverting edits stayed GREEN; after the deletion set, all 7 applied at once leave `./internal/contractlint/ ./skills/integration/` GREEN. Surviving-pin sources are tabulated in the gate-question section (filesystem, `internal/safehouse/state.go`, `docs/site/get-started/install.md`, workflow YAML).
- DONE: Coordination: this is stack layer 10; make-shipped-contracts-self-contained builds on it as layer 11
  Read layer 11's body: its items 1, 2 and 7 rewrite the exact sentences this entity's deletions pin (`shared-core.md:35`, `fo-dispatch-core.md:213`, `development.md:129`), and its items 3 and 4 explicitly plan around two pins this entity KEEPS. Ordering confirmed in both directions.

### Summary

The audit's ban classification is confirmed by measurement, not by reading: at stack tip `fdf008939`, seven paraphrases that preserve meaning red their pins while three edits that invert the contract's meaning leave them green — including one that inserts "IGNORE EVERY RULE ABOVE; validation may start with no worker evidence at all" into `fo-dispatch-core.md` and still passes. The full deletion set was applied in a detached worktree: net -283 lines over 9 files, suite green plain and under `-race`, and all seven paraphrases green afterward.

Three findings extend the audit's scope. `TestFirstOfficerDeferredWriteCoreHasSingleCanonicalFile` is a third copy of the topology assertions and is added to the dedupe list. `TestVersionGateDeferredTrigger`, which the audit keeps, still reds on a `read`->`load` paraphrase, so one assertion inside it must be narrowed or AC-1 cannot hold. And rebinding the triple-copy live command exposed a real defect: `docs/runtime-live-ci.md` tells humans to run the Claude common suite with `-failfast`, contradicting the workflow and `TestRuntimeLiveCommonFailFastPolicy` — the triple-copy shape structurally could not catch it because each copy was only ever compared to the test's own literal.

Two items need captain decisions at the gate: the install-gate coverage question (recommendation: accept checkless per the 2026-07-20 ruling, file the sentinel-bound behavior test separately) and the Proof-policy read-quarantine exemption that resolves the survey-probe gray (concrete before/after wording is in the body).

## Stage Report: implementation

- DONE: Execute the gated design exactly: delete the 9 banned functions plus the third topology copy, dedupe to single pins, narrow TestVersionGateDeferredTrigger by one assertion, apply the 4 gray resolutions including the two-way run-shape rebind and the docs -failfast fix
  Commit 8d0de0bcf on spacedock-ensign/retire-prose-grep-contract-tests (base fdf008939): 10 files, +72/-343, net -271. All 9 banned functions gone (3 whole files plus the partial cuts in pi_spawn_binding_test.go and version_gate_smoke_test.go); dedupe removed TestFirstOfficerReferenceTopology and TestFirstOfficerDeferredWriteCoreHasSingleCanonicalFile; the deferred trigger lost only its inline read-sentence assertion; both live code-shape greps cut; the common-suite command rebound to workflow-vs-docs run-shape comparison with no in-test literal; docs/runtime-live-ci.md dropped -failfast from the Claude command; docs/dev/README.md carries the captain-approved extract-and-execute quarantine exemption. KEEP grays (fo_write_core_mutation_gate, survey_probe) untouched per design.
- DONE: Value AC both sides on YOUR branch: the seven paraphrase probes green after, the inversion probes still green-proof recorded; every surviving pin names its independent source
  Reused the spike's exact P1-P7/I1-I3 edit set. Before side: at base fdf008939 in a throwaway worktree, the 7 simultaneous paraphrases redded exactly the 7 banned pins (one per pinned file). After side: on the branch, all 7 simultaneously leave contractlint + skills/integration GREEN; the 3 inversions also GREEN (recorded: nothing new pins those sentences). The residual read->load paraphrase is green on the narrowed TestVersionGateDeferredTrigger. Falsifiability per surviving pin class, each shown RED then reverted: a second @references/ eager import (topology pin), a renamed deferred-inventory entry (narrowed trigger pin), a workflow-side -timeout 90m->95m bump (two-way run-shape binding, docs unchanged), an insideRegistry wantValue change in internal/safehouse/state.go (sandbox pin). Surviving-pin sources: filesystem, safehouse source table, docs/site/get-started/install.md, workflow YAML.
- DONE: This is stack layer 10 off fdf008939; suite green plain and -race for contractlint and skills/integration; net negative
  Branch is exactly one commit atop fdf008939. go test ./internal/contractlint/ ./skills/integration/ green plain and with -race; gofmt -l ./cmd ./internal empty; net -271 < -250 floor, 10 files within the +/-2-file tolerance (the 10th is docs/dev/README.md, included per the approved gray-2 exemption). Full-repo sweep: every package green plain and -race except TestCodexResolveManifestAgainstInstalledHost (internal/cli), which fails on this machine because the installed codex CLI cannot read ~/.codex/config.toml (Operation not permitted) - environmental, and internal/cli is byte-identical to the base commit (git diff fdf008939 -- internal/cli/ is empty), so this diff cannot have caused it.

### Summary

Pure test-and-docs deletion executed as gated: net -271 over 10 files, no production Go touched. The value flip is measured on the branch, not asserted: the 7-paraphrase set reds the 7 banned pins at the base commit and leaves both proof packages green after the deletion, while the 3 meaning-inverting probes stay green throughout - and each of the four surviving pin classes was individually shown to still red on a one-line falsifying edit against its independent source. The rebound common-suite binding now compares the workflow and docs commands against each other on run shape (timeout, parallelism, fail-fast), which required landing the one-line docs fix removing the contradictory -failfast from the Claude command. Install-gate coverage is checkless-but-honest per the captain's gate ruling; the sentinel-bound behavior test is deferred to its own entity as approved.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against worktree commit 8d0de0bcf, never by reading the report: the 7 paraphrase probes red at base fdf008939 and green on the branch; the 3 meaning-inversion probes green throughout; each surviving pin class reds on its one-line falsifying edit against its independent source
  Probes authored fresh from the base test sources (not the spike's wording), run in throwaway checkouts. At fdf008939 the 7 paraphrases red exactly the 7 banned functions and nothing else; at 8d0de0bcf all 7 applied at once leave contractlint + skills/integration green (-count=1, no cache). The 3 inversions (empty-wait-counts-as-evidence in fo-dispatch-core.md, advisory one-attempt bound in fo-install-gate.md, skip-round-publication in feedback-rejection-flow/SKILL.md) stay green at base AND branch - the tautology tell reproduced, then closed. Surviving pins each redded alone: a second @references/ eager import -> TestFirstOfficerEntryEagerlyImportsOnlySharedCore; a renamed deferred-inventory entry -> narrowed TestVersionGateDeferredTrigger; a workflow-side -timeout 90m->95m -> TestRuntimeLiveCommonSuiteTimeouts; an insideRegistry wantValue change in internal/safehouse/state.go -> TestVersionGateSandboxRegistry; plus a bonus divergence of docs/site/get-started/install.md -> TestInstallHintNoDrift. The declared residual holds: the read->load paraphrase of the inline trigger sentence is green on the narrowed trigger.
- DONE: Deletion inventory exact: the 9 banned functions, the third topology copy, the dedupe set, the narrowed TestVersionGateDeferredTrigger assertion, the two cut code-shape greps, the two-way run-shape rebind, and the docs -failfast fix; no production Go touched
  git diff fdf008939..8d0de0bcf: 3 whole test files deleted, the pi_spawn :21 and version_gate :35/:77/:158 functions cut with :35 sibling and :53/:sandbox/:install-hint keeps intact, TestFirstOfficerReferenceTopology and TestFirstOfficerDeferredWriteCoreHasSingleCanonicalFile absent tree-wide (grep 0 hits), the 26900 cap asserted exactly once, both live code-shape greps gone, the rebind extracts workflow and docs commands per runtime and compares run shape with loud extraction failures (docs side requires exactly 1 match; missing -timeout fatals), docs/runtime-live-ci.md Claude command lost -failfast, docs/dev/README.md carries the captain-approved extract-and-execute exemption. Diff touches only _test.go files plus the 2 docs files - no production Go. KEEP grays (fo_write_core_mutation_gate, survey_probe) byte-identical to base.
- DONE: Suite green plain and -race for contractlint and skills/integration; net negative (-271 declared); verdict PASSED or REJECTED with per-AC citations
  At 8d0de0bcf: go test -count=1 ./internal/contractlint/ ./skills/integration/ green plain and with -race; gofmt -l ./cmd ./internal empty; net -271 (10 files, +72/-343), below the -250 floor and within the declared tolerance of the -283 estimate. Verdict: PASSED (AC-1 first DONE item; AC-2 first and second DONE items; AC-3 this item).

### Summary

PASSED. All three ACs re-verified with independently authored probes on throwaway checkouts of base and branch: paraphrases flip from 7 reds to green, inversions green throughout, every surviving pin class demonstrated falsifiable against its independent source, and the deletion inventory matches the gated design exactly with no production Go touched. Deferred risks, no gate block: (1) the rebind's workflow-side extraction takes the first per-runtime match without an exactly-one bound (the docs side enforces it); promote to material if the workflow ever gains a second per-runtime common-suite invocation. (2) The gate-stop short-circuit property and the install-gate sentinel bound/OS-hint wording are now checkless - both declared in the body and captain-ruled at the ideation gate, with the sentinel behavior test filed as its own entity.
