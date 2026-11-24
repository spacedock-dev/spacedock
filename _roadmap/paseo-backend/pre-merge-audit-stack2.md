# Pre-Merge Whole-Stack Staff Review

**Verdict: HOLD — 3 ship-blockers, 2 non-blocking follow-ups, 1 polish item.**

The stack improves runtime selection and keeps worker recovery free of a second persistent registry. Two contract defects remain: replacement-spawn failures can re-enter the retry path, and explicit Paseo selection can fail boot on a marker-free harness. Required live evidence for the assembled tip is also incomplete.

## Scope and review discipline

Reviewed all 14 changed files in `origin/main..0f01ad2f1ba2535c5b8c86b3c7fa7666c659bcff`:

- Base: `3d134b255fbcc493bd5e5ccc77a16c496c8b3662`.
- Layer 1: PR #895, through `ebb9f7b1fc01fb644deaea9e81f2e5bddfa410b8`.
- Layer 2: PR #897, through `0f01ad2f1ba2535c5b8c86b3c7fa7666c659bcff`.
- Review date: November 23, 2026, Pacific time. CI observations: approximately `2026-11-24 03:53–03:58 UTC`.

The roadmap requires an independent review of the assembled diff before merge. I read its separate state checkout at `/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/_roadmap/README.md`. I also read the runtime-support and development proof policies. This report is an independent source and evidence review, not an adoption of the earlier validation verdicts.

No Git checkout or `.spacedock-state` was edited. The checkout was clean before and after review. Test execution excluded tests that create or modify Git/state fixtures. The only durable review file is this report. No CI jobs were approved, started, or changed.

## Ship-blockers (must fix before merge)

1. **A failed replacement spawn can bypass the one-retry limit.** Owner: Layer 1, PR #895.

   Exact citations: `skills/first-officer/references/paseo-fo-dispatch.md:68`, `:114`, `:117`, `:118`, `:119`, `:136`, `:139`.

   Consider this concrete sequence:

   - Worker `c1-a1` ends without a valid report.
   - Row 5 archives it and starts replacement `c1-a2`.
   - The `a2` launch exits 1 before creating an agent.
   - Verb 2 requires reconciliation. The `--all` query still returns archived `a1`.
   - There are zero running workers, no valid report, and no `a2` row. The nonempty roster takes Row 5 again.

   The routing rule at line 118 sends a failed spawn to Row 8 only when the entire filtered roster is empty. Row 8 instead describes absence for the specific `(stage, cycle, attempt)`. These predicates disagree whenever a prior attempt exists. Row 5 can therefore request `a2` repeatedly, despite Row 8 requiring a stop after an `a2` failure.

   This is a source-level counterexample to the executable prose, not a claim that a live agent was observed looping. An agent can infer the intended limit, but the shipped decision paths must not require that repair.

   **Required change:** After a failed create, distinguish the current attempt from earlier attempts. If the current attempt exists, adopt or recover it. If it does not exist, apply Row 8 with the current in-turn attempt before ordinary roster recovery. Keep the existing labels and the prohibition on lease/state sidecars.

   **Acceptance:** Exercise archived `a1` plus failed `a2` creation. There must be exactly one replacement-create attempt, an entity-local stop, and no third create. Also preserve adoption when a successful `a2` create loses its response.

2. **Explicit Paseo selection fails the binary gate on a marker-free harness.** Owner: Layer 2, PR #897.

   Exact citations: `skills/first-officer/SKILL.md:35`, `:36`, `:42`, `:43`; `internal/cli/cli.go:914`, `:915`, `:921`; `skills/first-officer/references/first-officer-shared-core.md:9`.

   Reproduction conditions: no Claude, Codex, Pi, or Paseo environment markers; the captain explicitly selects Paseo; the daemon is reachable; the binary has the required minor version.

   Step 1 selects the Paseo runtime and skips the unsupported-harness branch. The binary still emits only its three terminal lines, because selecting a transport does not set `PASEO_AGENT_ID`. The shared startup gate then aborts because `contract 3` is absent. The exception at line 43 only belongs to branch 3, whose condition is false when Paseo is reachable.

   I executed the real CLI with every host marker removed and `Version=0.29.0`. It returned exactly:

   ```text
   spacedock 0.29.0
   OS: darwin/arm64
   Channel: edge (spacedock@spacedock-edge)
   ```

   Thus this failure is independent of a wrong-version installation. The CLI output is executable evidence; the boot failure follows from the two conflicting instruction branches. I did not launch an agent to mutate a workflow.

   **Required change:** Make the marker-free version rule apply independently of worker-transport selection, including explicit reachable Paseo. Preserve the required-minor check and rejection of an explicitly incompatible contract token. Do not manufacture an agent ID or add Paseo to native-host ambiguity detection.

   **Acceptance:** A marker-free harness with explicit reachable Paseo passes boot and selects only Paseo. Marker-free native fallback also passes. A wrong binary minor still stops. Keep the shared-core byte cap unchanged.

3. **The assembled stack lacks the required completed live proof.** Owner: stack integration / PR #897.

   Exact citations: `docs/dev/README.md:85`; changed runtime surfaces at `skills/first-officer/SKILL.md:35`, `skills/ensign/SKILL.md:12`, and `skills/first-officer/references/fo-dispatch-core.md:188`; adoption proof at `internal/ensigncycle/paseo_adoption_live_test.go:22`.

   The declared policy requires every host lane when the host-neutral dispatch core changes. At review time, the exact-tip runs showed:

   | Lane | Observed result | Evidence |
   |---|---|---|
   | Claude | Waiting | [Runtime run 45322592611](https://github.com/spacedock-dev/spacedock/actions/runs/45322592611) |
   | Codex | Waiting | [Runtime run 45322592611](https://github.com/spacedock-dev/spacedock/actions/runs/45322592611) |
   | Pi | Waiting | [Runtime run 45322595978](https://github.com/spacedock-dev/spacedock/actions/runs/45322595978) |
   | Paseo dispatch | Waiting | [Paseo run 45322589888](https://github.com/spacedock-dev/spacedock/actions/runs/45322589888) |

   Offline jobs succeeded. Their success does not establish the live behavior. The recorded Layer 1 adoption result is on `ebb9f7b1f`; it cannot establish compatibility with Layer 2's changed boot selection. The existing Paseo CI command selects only `TestLivePaseoDispatchSmoke` (`.github/workflows/paseo-live-e2e.yml:111`). It will not prove adoption even when that job passes.

   **Required completion:** Obtain green required native lanes on the final assembled revision. Run both Paseo smoke tests on that revision. Preserve FO traces that show the reachability probe and the absence of inner-host FO adapter reads. Include native fallback evidence with ambient `PASEO_AGENT_ID` and an unreachable daemon. Run the required full normal/race checks after the fixes, or provide equivalent existing results for that exact revision.

   These are existing project requirements. No new release gate or standing test system is proposed. This review did not change environment approvals.

## Non-blocking follow-ups (seed after merge)

1. **Make the existing adoption test detect regressions in the new boot seam.**

   Exact citations: `internal/ensigncycle/paseo_adoption_live_test.go:43`, `:67`, `:73`, `:200`, `:218`, `:221`; `skills/first-officer/SKILL.md:35`, `:36`.

   The test correctly checks real worker identity, duplicate absence, archival, commits, terminal state, and cleanliness. Its trace predicate only requires an inspect/wait call for the original worker. It does not require the FO's `paseo status --json` probe or reject an inner-host runtime-adapter read. Both FO prompts also explicitly request Paseo, so they do not test selection from the ambient marker alone.

   A regression that restores an unnecessary Codex adapter load can therefore leave this smoke green. Reuse the existing trace observation to cover that load/probe seam, and use an ambient-only prompt in an appropriate case. Keep durable state as the completion oracle. This can follow merge if blocker 3 supplies direct trace evidence for the candidate.

2. **Separate addressability from context-budget availability in the generic fallback.**

   Exact citation: `skills/first-officer/SKILL.md:45`.

   The fallback groups `«addressable-worker»` and `«context-budget»` under the existence of a live follow-up handle. A handle establishes a route to the worker. It does not establish a context-limit/usage probe.

   Bind addressability from the handle and context budget from an actual supported probe. Keep context budget absent otherwise. This is a contract-precision follow-up: the sentence does not require a fabricated budget, and no incorrect budget decision was observed.

## Polish

1. **Keep the no-sidecar rule in one owning section.**

   Exact citations: `skills/first-officer/references/paseo-fo-dispatch.md:5`, `:145`, `:152`.

   The opening invariant fully states the three sources and the forbidden writes. The spawn and roster bindings repeat the prohibition and source list. They load with the same file, so that duplication adds no load-boundary protection. Keep the full rule at line 5 and use short section references in the bindings. Preserve all actual reconciliation predicates.

## Composition and boundary audit

### Paseo reachability, restart, and sandbox fallback

The normal marked-session seam is coherent. `SKILL.md:35` requires both exit 0 and `connectedDaemon == "reachable"`. A stopped daemon's successful status command alone cannot select Paseo. `SKILL.md:36` and `paseo-fo-dispatch.md:3` both require the adapter before resume reconciliation, next-action, spawn, or terminal teardown. The shared loop starts with roster reconciliation (`fo-dispatch-core.md:231`). A fresh FO resuming an active ungated worker therefore has a stated load boundary before roster use.

An ambient marker with an unreachable daemon falls through to the native adapter. Explicit captain selection instead stops. This distinction supports the documented workspace-only sandbox. The Paseo module's opening condition also requires active/reachable Paseo, so it no longer independently activates from the ambient marker after fallback.

No missed restart load boundary was found in this path. The marker-free explicit-selection defect is blocker 2. Recovery is not contradiction-free because of blocker 1. Other matrix paths preserve the important distinctions: timeout waits, permissions hold, one live worker is adopted, duplicate live workers hold, and recorded `a2` exhaustion stops.

### Worker-state ownership and prohibited dependencies

The Paseo module contains no instruction to create a lease, append `## Paseo Worker`, or append `### Dispatch Retries`. Those strings appear only in prohibitions. There are no ensign `references/*.md` citations in the FO module. The module retains required dispatch/status frontmatter mutations; the zero-write promise concerns the worker registry in the entity body.

The state inputs are the entity/history, code Git state, and Paseo labels/inspect. The epoch and path/time filtering prevent a simple name match from being treated as ownership. The report check requires committed new evidence and clean relevant paths. The cascade guard remains intact. No second persistent roster, lease subsystem, or production Go recovery implementation was added.

### Adoption smoke under the new SKILL.md

The harness installs the candidate skills and uses `$spacedock:first-officer` for both FOs. The isolated daemon helper requires `connectedDaemon == "reachable"` before returning (`internal/ensigncycle/paseo_live_helpers_test.go:123`). It passes the isolated Paseo environment onward (`internal/ensigncycle/paseo_live_controls_test.go:27`). The Codex path uses full-access mode.

Those conditions support the new probe and skipped Codex-adapter path. The test stops FO-1 without archiving it, observes W1 alive, starts FO-2, and releases the report only after an inspect/wait observation. It then checks one worker, original epoch, archival, committed output, a post-baseline report commit, terminal verdict, and clean trees.

**Assessment: statically compatible and compilation passed; live success on the assembled tip is not established.** Probe execution and skipped adapter reads need the trace evidence in blocker 3. Follow-up 1 addresses durable coverage of those claims.

### Native Claude, Codex, and Pi behavior

| Situation | Expected FO route | Expected ensign route | Review evidence |
|---|---|---|---|
| Paseo marker unset; no explicit Paseo request | Skip Paseo command; native adapter | Native adapter | Source trace; native detection/version unit tests pass |
| Ambient Paseo marker; daemon unreachable | Native adapter | Native adapter when prompt does not explicitly select Paseo | Source trace; native host still wins binary detection |
| Reachable Paseo selected | Paseo; skip inner FO adapter | Paseo when emitted prompt names its adapter | Source trace; exact-tip live proof pending |
| No markers and no selected Paseo | Generic fallback | Generic/Paseo completion adapter | Marker-free CLI result confirmed; artifact test structure reviewed |

The ensign's explicit dispatch-prompt match takes precedence over native markers. Ambient `PASEO_AGENT_ID` alone selects its Paseo adapter only when all four native markers are absent. This avoids treating a native child of a sandboxed FO as a Paseo worker merely because it inherits the outer marker.

`runtimehost.Detect` is unchanged and still recognizes only native hosts. `printVersion` applies the Paseo fallback only after unambiguous native detection returns no host. Explicit `--host paseo` controls dispatch shape separately. This preserves native ambiguity handling.

**I cannot certify “100% regression-free.”** No native detection regression was found, but unit output tests do not execute model-driven skill routing. The required live lanes are pending.

### Shared core, contract lint, and documentation

- Shared core: **23,464 bytes**, against a **23,500-byte cap**; 36 bytes remain. It shrinks by 33 bytes across the stack.
- The modified shared core, dispatch core, and feedback skill contain no `Skill(skill=...)` calls. The generic load wording keeps the same load targets and triggers.
- The dispatch core now permits an adapter-required explicit `--host`, which resolves the prior conflict with Paseo's build command.
- Conditional worker retention at a gate composes with Paseo's unaddressable, one-shot workers and with native addressable workers.
- Contract lint passes, normal and race, except the intentionally unexecuted Git-writing break-glass test. Structural lint does not prove recovery decisions or runtime routing.
- The changed `docs/` prose passes the supplied ASD-STE100 structural checks: sentence limits, active/simple forms, permitted modals, and command preservation. Both pages describe current behavior without history or discarded alternatives. This is a review of changed prose, not certification of the entire documentation tree against the official vocabulary dictionary.

## Prose-section and code-behavior necessity inventory

Every changed file is accounted for below. Citations refer to the reviewed tip.

| File / anchor | Added or modified content | Necessity and disposition |
|---|---|---|
| `docs/runtime-live-ci.md:105` | Combined smoke command, adoption description, FO provider override | Necessary discoverability for the new smoke. Existing CI limitation is accurately stated. No unnecessary section. |
| `docs/site/reference/command-reference.md:7` | Version output description | Necessary user-visible Paseo/contract behavior. Short and structurally STE-clean. |
| `internal/cli/cli.go:914` | Paseo-only version fallback | Small, sufficient display-layer change. No host command execution or dead branch. Keep detection policy outside `runtimehost.Detect`. |
| `internal/cli/version_session_test.go:125` | Paseo render cases, contract placement, CLI execution witness | Necessary output and privacy coverage. Native-marker scrubbing makes the witness test meaningful. Table cases are inexpensive; no separate mechanism to remove. |
| `internal/dispatch/build_host_test.go:29` | Native-plus-Paseo and Paseo-only host cases | Necessary proof that display fallback does not change dispatch detection. No duplicate production logic. |
| `internal/dispatch/build_paseo_host_test.go:20` | Reusable assertion body, scrubbed-marker case, artifact cleanup | Reuses the existing shape oracle instead of adding a harness. The larger indentation diff is mechanical. Does not prove agent routing, which needs live evidence. |
| `internal/ensigncycle/paseo_adoption_live_test.go:22` | FO interruption/adoption orchestration and durable assertions | Necessary restart proof. Uses shared daemon helpers and real Paseo operations. The release file is an external test barrier, not a production lease. Helpers are used; no dead code found. Follow-up 1 covers the new boot seam. |
| `skills/ensign/SKILL.md:12` | First-match selection and generic fallback | Necessary precedence rule. Explicit prompt selection and the full native-marker exclusion are load-bearing. No redundant second adapter. |
| `skills/ensign/references/ensign-shared-core.md:112` | Paseo pointer example; expands “either” to “any” | Necessary bootstrap addition alongside existing host forms. Preserves first assignment read. |
| `skills/feedback-rejection-flow/SKILL.md:27` | Host-neutral skill load | Necessary removal of a concrete host tool from a shared path. Correction sequencing remains unchanged. |
| `skills/first-officer/SKILL.md:33` | Reachability selection, deferred Paseo load, native fallback, generic fallback | Necessary central selection owner. Blocker 2 covers the missing version-rule composition; follow-up 2 covers capability precision. No new framework needed. |
| `skills/first-officer/references/first-officer-shared-core.md:45` | Host-neutral loads and conditional retention at `:76` | Necessary shared cleanup. Triggers remain explicit and byte cap passes. No reason to restore host-specific call syntax. |
| `skills/first-officer/references/fo-dispatch-core.md:61` | Generic gate load and explicit-host rule at `:188` | Necessary composition fix. Keeps one dispatch shape and the existing refusal policy. |
| `skills/first-officer/references/paseo-fo-dispatch.md:3` | Resume load boundary, epoch derivation, reconciliation, eight-row matrix, bindings | Necessary recovery owner. Consolidating wait/error outcomes is preferable to competing paths. Blocker 1 identifies an actual competing predicate; polish 1 identifies repeated invariants. |

## Validation record

Executed successfully without Git/state fixture writes:

1. Focused normal tests in `internal/runtimehost`, `internal/cli`, `internal/dispatch`, and `internal/ensigncycle`:
   - `TestDetectMarkerMatrix`
   - `TestVersionSessionRender`, `TestVersionContractTokenPlacement`, `TestVersionSandboxLineNamesTheSandbox`, `TestVersionExecutesNoHostCLI`, `TestVersionAmbiguousMarkersExitZero`
   - `TestResolveBuildHostFromMarkers`, `TestResolveBuildHostExplicitWins`
   - `TestPaseoLiveConfigEnvAndRunArgs`, `TestPaseoRunJSONAndPreconditions`
2. The same focused checks with `-race`, plus `TestFOInstructionComponentCaps`.
3. `go test -count=1 ./internal/contractlint -skip '^TestBreakGlassCommitLandsInEntityRepo$'` and its `-race` equivalent.
4. `go vet ./...`.
5. `go test -tags live -count=1 -run '^$' ./internal/ensigncycle` — compilation only, **not a live pass**.
6. Real marker-free CLI version execution with the required minor, as shown in blocker 2.
7. `git diff --check origin/main..HEAD` and `gofmt -l` on every changed Go file: clean.
8. Whole-tree `gofmt -l ./cmd ./internal`: one existing issue in `internal/release/runtime_live_evidence_workflow_test.go`. Reading that file from `origin/main` produces the same formatting diff. It is not introduced by either PR.

The exact-tip Runtime CI offline job also completed its `go test ./...` step successfully. That result was read from GitHub, not rerun locally. I did not execute local full-suite, full-suite race, live smoke, or dispatch-fixture tests because they write Git/state fixtures. I did not run `gofmt -w` because it writes the checkout. These omissions preserve the explicit read-only restriction; they are not reported as passing checks.

Review tracking is complete: all changed files inspected, composition paths traced, allowed checks run, CI evidence inspected, findings classified, and report delivered. No follow-up entities or source fixes were created.
