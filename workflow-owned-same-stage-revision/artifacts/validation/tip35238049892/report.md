# Completed Claude batch triage: run 35238049892

Candidate: `0997b6b7590a5a89ba979194238e93f7f47d94e4`. Claude job `105259444956`, artifact `10505302792`; Codex job `105259444983`, artifact `10505631330`. This is read-only triage, not candidate acceptance. Pi remained running at dispatch; its final result and whole-run disposition remain with the first officer.

## Completion checklist

- [x] Deduplicate the current Claude failures and establish exact runtime behavior against the actual candidate contracts and fixture outcomes.
- [x] Propose materiality, task ownership and smallest disposition for each distinct failure, preserving current positives and authority boundaries.
- [x] Commit artifact-only triage evidence with exact source paths and explicit remaining Pi/runtime obligations, without candidate or entity changes.

## Findings and proposed dispositions

### 1. Required-review hold lacks the entity-body missing-source note

**Proposed classification:** Material evidence-contract mismatch; same-stage revision ownership; HOLD / Needs decision on required recording location. This does not establish a missing-reviewer-dispatch defect.

Policy fields:
- Supported workflow/user: a declared independent review whose external source is unavailable after an authorized same-stage correction.
- Observable harm: the committed entity contains the correction report (“Ready for re-review”) but no `reviewer-source.txt` missing-evidence note. Durable entity readers do not receive the hold reason. The runtime stopped safely; no acceptance bypass was observed.
- Authority: **contract** — `writeSameStageRevision` in `internal/ensigncycle/claude_runtime_helpers_test.go`: “Record the missing evidence and hold.” It declares independent review after correction, but does not specify the recording location.
- Observed trigger: `review-required` corrected its selected plan, observed the source absent, and stopped before independent review. The existing grader rejected the entity-body omission.

`runSameStageRevisionJourney` in `internal/ensigncycle/claude_live_runner_test.go` explicitly allows missing evidence to stop before reviewer dispatch. Its route assertion requires review only when more than one spawn occurred. The failed predicate is `strings.Contains(body, "reviewer-source.txt")`, with message `missing required review evidence was not recorded`; it is not the independent-review route check.

The retained native stream shows one named correction Agent (tool call `toolu_01T82bgrb7qPwMxUyoBBdUQM`, native worker `ac823b1abbb07596e`), its completion, and an actual `ls` result omitting the required source. The final message reports the hold. Independently inspected `state.bundle` HEAD is `ffa017a2403cb1dc0d3252a4591b8300a24cff8c`: selected plan equals frozen input, missing source is absent, entity lacks its name, one rejected attempt remains, and zero canonical rounds exist. See `durable-review-observation.json` and retained snapshot. The final prose is corroborated by source absence and unchanged gate state, not used alone.

AC-2's missing-evidence stop is observed. The fixture's additional “Record” obligation remains unresolved because the oracle assumes entity-body storage while the instruction does not name that location. Do not waive this failure or introduce a mandatory reviewer to satisfy it. Smallest disposition: decide whether a committed entity hold note is required; if so, have the existing authorized writer persist it. If a different retained recording is allowed, align the existing oracle with actual source-absence and held-state evidence, not final prose alone. No new reviewer machinery or harness is needed.

### 2. Direct roadmap request refused after ladder edits

**Proposed classification:** Material observed runtime outcome defect; FO write-authority / smallest-mechanism ownership; retain prior HOLD and route a bounded correction proposal to that owner.

Policy fields:
- Supported workflow/user: a captain explicitly requests a plain repository document and direct commit alongside two named edits.
- Observable harm: ladder edits were committed, but `roadmap-strategy.md` was not created; the requested work is incomplete.
- Authority: **captain** — exact fixture prompt: “Create roadmap-strategy.md containing the one-line body `# Roadmap Strategy` and commit this roadmap document directly to the repository.”
- Observed trigger: the direct smallest-mechanism scenario treated absent roadmap workflow scope as requiring another grant despite the exact-target instruction.

`smallestMechanismDirectPrompt` in `internal/ensigncycle/shared_fixtures_test.go` supplies the exact task, path, body and direct commit authorization. `skills/first-officer/references/fo-write-core.md` allows exact-target grants through its Mutation Gate. Its Workflow Fit Gate applies “Before creating or materially reclassifying an entity”; this requested plain Markdown file is not an entity-creation request.

Native events show both ladder edits and successful commit `43476f5` (two changed files, clean status). No roadmap write occurred, and the existing grader observed the missing file. The final refusal explains a workflow-fit/authority interpretation inconsistent with the requested plain document. There is no retained state bundle for this scenario: this report claims native commit-result evidence and the grader's missing-file observation, not an independent Git replay.

This repeats the prior HOLD in `fo-continuation/artifacts/validation/tip35148797301/report.md`. Exact source comparison from `d0a6f413` to this candidate shows no change to the direct prompt or write-core contract. No causal attribution to a shutdown suffix is established. Smallest disposition is an existing-owner proposal that preserves exact-target authorization and confines workflow-fit to entity filing. Do not weaken the missing-file grader or require redundant captain permission.

### 3. First-outage selected-team recovery conflicts with canonical-name hold rule

**Proposed classification:** Material compatibility / contract conflict; semantic-dispatch-names (#802) and recovery ownership; Needs decision before correction. The runtime followed the current name rule.

Policy fields:
- Supported workflow/user: the helper fails on the first dispatch and the captain has selected team mode; recovery is expected to preserve that mode and complete one worker report.
- Observable harm: no Agent was dispatched and no stage report was produced in the selected-team recovery scenario.
- Authority: **contract** — `TestLiveBreakGlassShimRecovery` in `internal/ensigncycle/dispatch_recovery_live_test.go` promises one worker after helper failure, preserving selected mode and producing a committed report.
- Observed trigger: the fresh fixture has no cached successful envelope and its shim fails every `dispatch build`, leaving no helper-emitted canonical name.

The retained stream shows the real build failure (exit 1), a second unstamped failure, and zero Agent calls. The stamped provisional changes were reverted. No state bundle was retained for this scenario, so there is no claim of independent final Git replay.

The exact candidate `skills/fo-dispatch-recovery/SKILL.md` instructs selected-team recovery to use a successful canonical envelope, or a retained validated envelope for the same assignment; if none exists, hold and report the missing canonical name. It prohibits duplicate shortening logic or changing dispatch mode. `semantic-dispatch-names/index.md` proposed approach #5 and AC-3 explicitly support this cached-envelope-or-hold behavior. Therefore the runtime did not disobey that rule: the accepted naming contract narrows the standing first-outage recovery promise. Selected-team has no expected-failure exemption; the selected-bare expected-failure handling is unrelated.

The recovery skill is unchanged between `d0a6f413` and this candidate. Its canonical-name restriction originates in semantic naming, not an observed suffix effect. Smallest disposition: the existing naming/recovery owner must reconcile canonical identity availability with intended first-outage selected-team recovery. Preserve one worker, chosen mode, safe identity and committed report; do not merely accept HOLD, add an expected failure, or duplicate naming logic to make the test green. This needs FO scope/contract direction before implementation.

## Actual outcomes

Nine failed **test** records reduce to three independent failures: required-review leaf, smallest-mechanism direct scenario, selected-team recovery leaf. The remaining six are their same-stage/recovery wrappers, three scheduler slots, and scheduler root. The package failure is a separate aggregate record, not a tenth distinct test.

| Same-stage variant | Claude | Codex |
|---|---|---|
| plain | PASS 225.17s | PASS 180.30s |
| review-required | FAIL 258.70s | PASS 192.47s |
| separate-review-required | PASS 305.72s | PASS 236.66s |
| round-required | PASS 303.10s | PASS 262.39s |
| round-missing | PASS 228.66s | PASS 252.46s |
| cycle-limit | PASS 233.70s | PASS 203.08s |

All twelve variants actually ran; none skipped. Claude continuation (397.73s), rejection flow (421.47s), and default-headless (246.23s) passed. Codex continuation (295.43s), rejection flow (329.91s), default-headless (300.30s), smallest-mechanism (426.89s), and semantic names (133.62s) passed. Codex job passed; Claude job failed. Positive controls do not waive the three observed failures. Pi and final cross-runtime disposition remain pending this report's evidence boundary.

## Evidence and limits

`source-hashes.json` records exact candidate paths and SHA-256 values, original downloaded paths, and hashes for copied streams, final messages and state bundle. `claude-outcomes.json` / `codex-outcomes.json` retain selected original run/result records. Exact skill copies and the native streams support the contract/behavior distinction. The review-required snapshot and `durable-review-observation.json` come from the matching retained bundle. Original Claude detail: `/tmp/spacedock-tip-35238049892/claude/spacedock/spacedock/live-e2e-detail.jsonl`; original log: `/tmp/spacedock-tip-35238049892-claude.log`. Codex detail: `/tmp/spacedock-tip-35238049892/codex/codex-shared-scenarios-detail.jsonl`.

No candidate edits, entity/frontmatter edits, new tests, models, harnesses, CI, pushes or reruns were performed. FO authorization is required before any corrective candidate action or rerun. Existing prior correction evidence is neither replaced nor reclassified by this completed-runtime batch.

## Summary

Retained native and durable evidence isolates three distinct Claude outcomes. Missing independent reviewer dispatch is not the same-stage failure; its missing durable hold note is an evidence-contract question. The roadmap refusal repeats an incomplete authorized action. Selected-team recovery exposes conflicting naming and first-outage recovery contracts. All findings remain proposed for FO disposition; this artifact does not approve the candidate or the whole runtime run.
