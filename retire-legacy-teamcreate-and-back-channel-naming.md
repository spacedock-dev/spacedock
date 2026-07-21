---
id: q4pw3xb4nf4cwfdjtwbn17mz
title: Retire legacy TeamCreate path and rename back-channel to inter-agent communication
status: ideation
source: captain (CL), 2026-07-20 session
started: 2026-07-21T16:05:13Z
completed:
verdict:
score:
worktree:
issue:
---

Retire the legacy TeamCreate machinery from the shipped FO contract and stop naming the worker/FO messaging surface "back-channel".

## Problem

Two contract-surface cleanups, captain-directed:

1. **Legacy TeamCreate path.** `claude-fo-dispatch.md` carries a legacy-override line that makes every session's first dispatch probe `ToolSearch(select:TeamCreate)`; when it matches, the FO loads `skills/using-legacy-claude-team`. Current runtimes no longer expose TeamCreate (this session's probe: no match), so the probe is dead weight on the first-dispatch path, and `using-legacy-claude-team/SKILL.md` is one of the 13 FO-surface measured files — retired, it shrinks the measured prompt surface. The claude runtime adapter and the terminal-teardown prose also carry legacy-override references that ride only on this trigger.

2. **"Back-channel" naming.** The contract names inter-agent messaging "back-channel" throughout (`## Worker Back-Channel` in claude-fo-dispatch.md, Codex "mailbox back-channel", Pi "contact_supervisor/intercom back-channel", «addressable-worker» prose). Captain direction: do not name it "back-channel" — name it for what it is, e.g. "inter-agent communication" (final term chosen at ideation).

## Degraded Mode retirement — split to `9q4`, but this task must not land blind to it

The captain asked (2026-07-21) whether this task should also retire Degraded Mode. A two-seat adversarial scope review says the retirement belongs in its own entity — `9q4` (dispatch-failure-retry-rung) — because the defect is the missing retry rung BELOW the trigger, not the trigger alone, and because this task's own `## Out of scope` promises no behavior change to dispatch or messaging. Two findings from that review bind THIS task regardless:

1. **This task deletes the repo's only retry rung.** `using-legacy-claude-team/SKILL.md:50` — «legacy-team.recover» rung 1, "Attempt one new TeamCreate with a fresh name" — is the sole try-once-before-degrading step anywhere in the contract. Retiring this file without `9q4` landing together or first leaves the contract with strictly ZERO retry surface, which is a regression, not a cleanup. Sequence them.

2. **VERIFY BARE-MODE REACHABILITY BEFORE DELETING THE LEGACY-OVERRIDE LINE.** A review seat claims the `e3z`-proven bare-mode entry (ToolSearch-no-match → bare) was inverted into the legacy-override at `claude-fo-dispatch.md:9`, which this task deletes — and that clause (3) of the Degraded Mode trigger is the only remaining route into `bare_mode: true`. That reading is CONTESTED: `claude-fo-dispatch.md:7` carries an independent `SendMessage`-availability probe that falls back to "fresh one-shot dispatch" without touching the legacy line. The two are not obviously the same state (`«addressable-worker» ABSENT` vs. the `bare_mode: true` build flag). **Resolve this at ideation with a live no-team drive, not by reading.** `e3z` proved bare mode reachable through the pre-retirement contract; nothing yet proves it reachable through the post-retirement one, and the coverage vacuum below means CI will not catch a mistake.

3. **`internal/claudeteam/claudeteam.go:68-79` (`BareModeAdvisory`)** instructs the FO to run `ToolSearch select:TeamCreate` on every bare dispatch lacking recent team evidence. After this retirement that names a tool which no longer exists, on the legitimate bare path — a required co-edit, and Go source, so it widens this task beyond prose.

4. **No oracle exists in either direction.** `TestLiveDegradedBareRecovery` is `//go:build live` and appears in ZERO CI `-run` filters; `.github/workflows/runtime-live-e2e.yml:111` sets `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`, so the teams-unavailable branch never runs in CI.

## Proposed approach

Two coordinated retirements on the shipped FO contract, landing in ONE commit with `9q4` (never q4-alone-first — see `### Joint landing with 9q4`).

**1. Retire the legacy TeamCreate path.**
- Delete `skills/using-legacy-claude-team/SKILL.md` (14065 B) and its directory.
- Delete the `**Legacy override (…)**` line at `claude-fo-dispatch.md:9`, and the two legacy-load-point sentences in `claude-first-officer-runtime.md` (:7 "the sole legacy load point" parenthetical, :13 "its further bounded teardown is the legacy override").
- The go-forward team-mode entry is unchanged: `claude-fo-dispatch.md:7`'s `SendMessage`-availability probe → named-background-`Agent` (no `TeamCreate`). THIS session is direct evidence the model works with SendMessage-but-no-TeamCreate — it is running under it. Bare-mode reachability is the SAME :7 probe's ABSENT branch, independent of the deleted :9 line; 9q4's captain ruling #1 dissolves the Degraded-Mode trigger (3) to exactly "teams-unavailable → bare at the dispatch site", and 9q4's AC-2 live no-team drive is the shared proof (so this task does NOT re-litigate bare-mode reachability).
- Rewrite the two Go advisories that name `TeamCreate` as the team-mode bootstrap so they name the go-forward entry instead: `internal/claudeteam/claudeteam.go` `BareModeAdvisory` (9q4 co-edit #6, q4-owned) AND `PresentFalseHint` (same defect class — the boot `present:false` hint tells the FO to "run TeamCreate", which the current auto-team contract already contradicts; retirement makes it actively wrong). Both bindings read the Go symbol, not a literal (`build_advisory_probe_test.go`, `boot_probe_parity_test.go`), so rewriting the strings keeps their tests green.

**2. Rename "back-channel" → "inter-agent communication"** across the shipped contract/adapters (captain's chosen term). Occurrences: `claude-fo-dispatch.md` (the heading `## Worker Back-Channel` + ~7 prose/comment uses), `claude-first-officer-runtime.md:7`, `pi-first-officer-runtime.md:8`, `pi-ensign-runtime.md:8`. The Codex adapter uses "mailbox" (already descriptive) — no codex rename. Exact before/after in `### Shipped-contract prose diff`.

**3. Re-tighten the ratchet (AC-4), correcting a stale mechanism reference.** The real gate is the per-host `foHostLoadBaselineBytes` map in `internal/contractlint/fo_function_reference_invariant_test.go` — NOT the `foFunctionReferenceBaselineBytes`/`TestFOFunctionPromptSurfaceShrinks` named in this task's old AC-4 (grep-confirmed: neither exists in the tree). It is a strictly-ABOVE gate (`load > baseline` reds), so zero-slack = baseline == measured (NOT measured+1). All three hosts sit at EXACTLY baseline today (measured: claude 111183, codex 74608, pi 70725). The joint change drops the claude load by ≥14 KB (using-legacy dominates) → re-baseline `["claude"]` DOWN to the new measured value. The rename adds ~13 B/occurrence to the pi load → re-baseline `["pi"]` UP (a justified naming-improvement growth). `["codex"]` unchanged.

## Out of scope

- No behavior change to dispatch or messaging itself — this is dead-path removal + naming. The go-forward auto-team model and bare mode are untouched (bare-mode reachability is 9q4's AC-2, not re-proven here).
- The `TeamCreate`/`TeamDelete` branches in `internal/dispatch/build.go` + `internal/dispatch/standing.go`, and `claudeteam.go`'s `LegacyTeamNameAdvisory` (which describes the `--team-name` build SHAPE — a still-existing CLI path with its own test `build_team_name_advisory_test.go`). The using-legacy skill itself declared these "dead code, removable at leisure"; removing them widens the change into dispatch-build machinery for no contract-surface gain. Follow-up, not required by AC-1.
- `skills/commission/SKILL.md`'s TeamCreate step (:667-671, "failure mode #201") — a separate skill, not on the FO dispatch surface. Same staleness class, but not named in this task; follow-up.
- Retiring bare mode or Degraded Mode (Degraded Mode is `9q4`).
- The installer-stable removal-trigger governance (see `### Gate flag`) — q4 SUPERSEDES the AC-5 "retain until both conditions" decision by captain direction; whether a live drive on a TeamCreate-capable runtime is required before landing is a gate decision, resolved at the gate, not here.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - The shipped FO contract carries no legacy-TeamCreate path: no override line, no `using-legacy-claude-team` skill, no `select:TeamCreate` probe on the dispatch surface.**
Verified by: the repurposed `legacy_teamcreate_layering_test.go` retirement invariant (`TestLegacyTeamCreatePathFullyRetired`) — the skill dir is absent, the `spacedock:using-legacy-claude-team` token is absent from `claude-fo-dispatch.md`, and no `select:TeamCreate` appears on the FO dispatch surface. Falsifying change: re-adding the override line, the skill file, or the probe reds this test (and re-adding the token without the file reds `boot_resident_closure_test.go`'s dangling-target check; re-adding the 14 KB skill reds the AC-4 claude ratchet).

**AC-2 (VALUE) - The claude FO session load shrinks by at least the retired skill's size (≈14 KB) net of the rename.**
Verified by: post-joint `foHostLoadBytes(claude)` is ≥14000 B below the pre-change 111183, measured by `TestFOFunctionReferenceCheckpointMetrics`' `FO_FUNCTION_METRICS claude_bytes=` line before vs after — an independent number that moves the wrong way if the rename or 9q4's retry rung lands net-positive. Falsifying change: leaving `using-legacy-claude-team/SKILL.md` in `foHostLoadPaths["claude"]`, or a mis-scoped verbose retry rung eating the shrink below the ≥14 KB threshold.

**AC-3 - The shipped contract/adapters name inter-agent messaging descriptively ("inter-agent communication"), never "back-channel".**
Verified by: review-time grep over the FO surface returns no "back-channel" (no committed prose-grep, per proof policy); the renamed `## Inter-Agent Communication` heading resolves the `claude-binding` preservation anchor. Falsifying change: any residual "back-channel"; the `claude-binding` anchor still pinning `## Worker Back-Channel` reds `TestFOFunctionNormalizationPreservationSuite/claude-binding` (spike-observed).

**AC-4 - Every host ratchet the joint change touches is re-tightened to its post-change measured load (zero-slack).**
CORRECTS this task's original AC-4: the gate is the per-host `foHostLoadBaselineBytes` map (strictly-ABOVE; zero-slack = baseline == measured, NOT measured+1). The `foFunctionReferenceBaselineBytes`/`TestFOFunctionPromptSurfaceShrinks` this AC named do not exist; the "122126/122634/123323" figures are a different, un-gated 13-file-union sum. Carries the F1 mitigation from `codex-post-compaction-contract-reload` (archived, cycle 7): with all three hosts already at zero slack, any inert-prose reintroduction reds the ratchet, so re-tightening closes the inert-heading escape mechanically. `["claude"]` → post-joint measured (down); `["pi"]` → post-rename measured (up ~13 B, justified by the rename); `["codex"]` unchanged (no back-channel there).
Verified by: `TestFOHostPromptLoadRatchet` green with each edited baseline == its measured load; `TestFOHostPromptLoadRatchetDiscriminates` proves a +1 regression reds exactly that host. Falsifying change: leaving a baseline loose (above measured), or leaving `["pi"]` un-raised for the rename reds the pi ratchet (spike-observed: pi was at exactly baseline).

**AC-5 - No dispatch/boot advisory instructs running `TeamCreate`; the go-forward team-mode entry (SendMessage-probe → named-`Agent`) is named instead.**
Verified by: `internal/dispatch/build_advisory_probe_test.go` green against the rewritten `BareModeAdvisory` and `internal/status/boot_probe_parity_test.go` green against the rewritten `PresentFalseHint` (both bind the Go symbol/behavior, not the literal wording). Falsifying change: an advisory still emitting "run TeamCreate".

## Test plan

- **Contractlint:** `go test ./internal/contractlint/...` green after all co-edits. The ideation spike enumerated the exact red set (6 test files) the joint change must resolve — see `### Ideation spike result`.
- **Full + race:** `go test ./...` and `go test -race ./...` green; `go build ./...` (spike: exit 0 with deletions applied).
- **Joint-only green (never q4-alone-first):** empirically re-grounded below — q4-alone reds `TestUsingLegacyClaudeTeamDegradedModePointersNameRealAnchor` + the 3 layering retention tests; 9q4-alone reds the same DegradedModePointers test via its `## Degraded Mode` anchor drop; only the one commit (q4 deletes that test's subject, 9q4 drops the anchor) is green.
- **Live lane:** the bare no-team drive is `9q4`'s AC-2 (shared). The installer-stable / TeamCreate-capable-runtime team-mode drive is FLAGGED as gate-required-but-unverifiable on this session (no TeamCreate here) — see `### Gate flag`.
- **No committed prose-grep** (proof policy): the rename absence is a review-time grep; the retirement is guarded mechanically by the repurposed retirement-invariant (AC-1) + the AC-4 ratchet + the boot-resident dangling-target check.

### Shipped-contract prose diff

Exact before/after for each shipped-contract edit (implementation applies verbatim; the gate reviews this diff).

**`skills/using-legacy-claude-team/SKILL.md`** — DELETE the file and its directory (14065 B).

**`skills/first-officer/references/claude-fo-dispatch.md`:**
- heading — before `## Worker Back-Channel` → after `## Inter-Agent Communication`
- :3 — before "…the runtime adapter): the worker back-channel, the `Agent()` spawn call…" → after "…): inter-agent communication, the `Agent()` spawn call…"
- :7 — before "Claude PROVIDES the worker back-channel (fo-dispatch-core.md `## Dispatch Adapter`, the organizing capability) via…" → after "Claude PROVIDES inter-agent communication (fo-dispatch-core.md `## Dispatch Adapter`, the organizing capability) via…"
- :9 — DELETE the entire `**Legacy override (delete this line to sunset legacy mode):** …otherwise this file applies as written.` line.
- :11 — before "skip the background back-channel and use bare-mode dispatch" → after "skip background inter-agent communication and use bare-mode dispatch"
- :27 — before "// the lead→worker back-channel; omit if bare mode (field absent)" → after "// the lead→worker channel; omit if bare mode (field absent)"
- :29 — before "// the worker→lead back-channel; omit when field absent" → after "// the worker→lead channel; omit when field absent"
- :106 — before "With the background back-channel, the fix agent and reviewer can interact via messaging." → after "With background inter-agent communication, the fix agent and reviewer can interact via messaging."

**`skills/first-officer/references/claude-first-officer-runtime.md`:**
- :7 — before "The Claude dispatch parts — the worker back-channel, the ID/next-id read, …" → after "The Claude dispatch parts — inter-agent communication, the ID/next-id read, …"; AND delete the trailing parenthetical "(`claude-fo-dispatch.md`'s one legacy-override line handles a runtime that still exposes `TeamCreate`; it is the sole legacy load point.)"
- :13 — before "…already loaded at first dispatch). When the runtime still exposes `TeamCreate`, its further bounded teardown is the legacy override." → after "…already loaded at first dispatch)." (delete the trailing legacy sentence).

**`skills/first-officer/references/pi-first-officer-runtime.md:8`** — before "PRESENT — the `contact_supervisor`/`intercom` back-channel." → after "PRESENT — the `contact_supervisor`/`intercom` inter-agent communication."

**`skills/ensign/references/pi-ensign-runtime.md:8`** — before "the FO routes follow-up through the back-channel." → after "the FO routes follow-up through inter-agent communication." (not ratcheted; consistency with the captain directive.)

**`internal/claudeteam/claudeteam.go`** (behavior-bound, final wording at implementation):
- `BareModeAdvisory` (:73-79) — before names "run ToolSearch select:TeamCreate and TeamCreate first"; after names the go-forward entry: "If you intend teams mode, ensure SendMessage is available (`ToolSearch select:SendMessage`) and dispatch a named background Agent — no TeamCreate."
- `PresentFalseHint` (:22) — before "run TeamCreate before first team-mode dispatch (claude runtime supports it)"; after "dispatch a named background Agent for team-mode messaging (no TeamCreate needed)".

### Co-edit map (contractlint + Go)

**q4-owned** (each keyed to a spike-observed red or a known binding):
1. `fo_function_reference_invariant_test.go` — remove `skills/using-legacy-claude-team/SKILL.md` from `foFunctionReferencePaths` (:31) AND `foHostLoadPaths["claude"]` (:55) together (`TestFOHostLoadSetsCoverAddressLintUnion` keeps them in sync); remove the using-legacy entries at :325 / :340 / :361 / :404-405; rename the `claude-binding` anchor :365 `## Worker Back-Channel` → `## Inter-Agent Communication`; AC-4 re-baseline `["claude"]` (down) and `["pi"]` (up), `["codex"]` unchanged. [reds: TestFOFunctionReferenceInvariant, TestFOHostPromptLoadRatchet, TestFOFunctionReferenceCheckpointMetrics, TestFOFunctionRequiredCallSites, TestFOFunctionNormalizationPreservationSuite/{legacy-recovery,claude-binding}, TestFOLocalOrderedProceduresPreserved]
2. `legacy_teamcreate_layering_test.go` — REPURPOSE the 3 retention tests into a retirement invariant (not pure delete): keep `TestNormalPathContractInlinesNoLegacyMachinery`'s contract-absence half (drop the now-missing legacy-skill read at :77/:88); invert `TestNormalPathContractNamesLegacySkill` → assert the token is ABSENT; rewrite `TestLegacyConsumerRetiredButPathLives` → `TestLegacyTeamCreatePathFullyRetired` (skill dir gone + token gone + no `select:TeamCreate` on the FO dispatch surface). This is AC-1's mechanical guard. [reds: all 3, resolved by the inversion]
3. `dispatch_recovery_value_binding_test.go` — delete `TestUsingLegacyClaudeTeamDegradedModePointersNameRealAnchor` (:72-136) + its `usingLegacyClaudeTeamDegradedModePointerMarkers` var (:62-70); its subject is deleted. [reds: that test — also 9q4's binding coupling]
4. `boot_resident_closure_test.go` — remove `"using-legacy-claude-team": true` from `lazyLoadSkills` (:94); swap the control fixture's real load point (:186 `spacedock:using-legacy-claude-team` → `spacedock:fo-dispatch-recovery`) so `sawReal` still resolves. [reds: TestBootResidentDeferredLoadPointGuardFailsOnDanglingTarget]
5. `structural_checks_test.go` — remove `"using-legacy-claude-team"` from `userSkills` (:23); drop the `using-legacy-claude-team` clause of the `~/.claude/teams` machine-dependency carve-out (:347), keep `survey`. [reds: TestUserSkillsParseWithFrontmatter, TestUserSkillReferenceClosureResolves]
6. `launcher_invariant_test.go` — remove `"using-legacy-claude-team"` from `deferredSkillPaths` (:65) and the discriminator want (:337); update the :58 comment. [reds: TestLauncherSurfaceUsesResolvedLauncher, TestDeferredSkillLauncherScopeDiscriminates]
7. `internal/claudeteam/claudeteam.go` — rewrite `BareModeAdvisory` + `PresentFalseHint` (AC-5). [no red — behavior-bound tests stay green]

**9q4-owned** (joint commit, listed to avoid duplication): `fo-dispatch-recovery/SKILL.md` (remove `## Degraded Mode`); `claude-fo-dispatch.md` (:83-85 trigger → retry rung + :68 carve-out + `### Dispatch Retries` note); `first-officer-shared-core.md:49` (recovery one-liner reword); `boot_resident_closure_test.go:53` (drop the `## Degraded Mode` anchor from `deferredSkillCores`); `dispatch_recovery_value_binding_test.go` (delete `TestDegradedModeCaptainReportPrefixBindsSkillBlockquote`); `dispatch_recovery_assert*`/`live`/`fixtures` (oracle rewrite + retarget); + the net-new AC-1 bounded-retry offline oracle.

**Overlap files (disjoint lines/functions → clean joint commit):** `claude-fo-dispatch.md` (q4 :3/5/7/9/11/27/29/106; 9q4 :68/83-85), `boot_resident_closure_test.go` (q4 :94/:186; 9q4 :53), `dispatch_recovery_value_binding_test.go` (q4 deletes one test; 9q4 deletes another).

### Joint landing with 9q4 (green-tree ordering)

Both couplings (retry-surface + binding) force ONE commit. The binding coupling is now empirically re-grounded: q4-alone reds `TestUsingLegacyClaudeTeamDegradedModePointersNameRealAnchor` (spike: RED) AND the 3 layering retention tests; 9q4-alone reds the SAME DegradedModePointers test via the `## Degraded Mode` anchor drop it makes (that test reads `deferredSkillCores[fo-dispatch-recovery]` at :94 and matches the pointer against it). Only the joint change is green — q4 deletes that test's subject and 9q4 drops the anchor it reads, both resolved by q4 deleting the test. Ordering within the one commit: apply all deletions + the layering-test inversion + the re-baseline together; run `go test ./...` before committing. Never q4-alone-first.

### Ideation spike result

Ran, in a throwaway worktree off HEAD: the deletions-only (delete `using-legacy` skill + the :9 override line) + the rename (heading + pi adapter). `go build ./...` exit 0. `go test ./internal/contractlint/...` red set was EXACTLY the mapped co-edit surface — 6 files: `fo_function_reference_invariant_test.go`, `boot_resident_closure_test.go`, `dispatch_recovery_value_binding_test.go`, `launcher_invariant_test.go`, `legacy_teamcreate_layering_test.go`, `structural_checks_test.go` — with NO surprise binding in an unmapped file. The rename reds `TestFOFunctionNormalizationPreservationSuite/claude-binding` ("none of the preservation headings found: [## Worker Back-Channel]") and the pi ratchet (pi at exactly zero slack). Conclusion: binding map complete, compile clean, joint-landing necessity re-grounded (the DegradedModePointers test reds q4-alone). Worktree removed; main tree unchanged.

### Gate flag: high-stakes surface + installer-stable governance

This is the shipped FO contract. Two items the gate must weigh:
1. **Detached adversarial audit + live drive** (stage-def staff-review for skill-integration changes). The go-forward auto-team-without-TeamCreate model is PROVEN on this runtime (this session runs it). The residual risk population is runtimes that expose TeamCreate but a different team model (installer-stable Claude 2.1.170): after retirement they hit the :7 `SendMessage` probe → auto-team if SendMessage is present, else bare. That branch CANNOT be exercised on this session (no TeamCreate here). 9q4's AC-2 bare live-drive covers bare reachability; it does NOT cover installer-stable team-mode.
2. **q4 supersedes an AC-5 governance decision.** `legacy_teamcreate_layering_test.go`'s `TestLegacyConsumerRetiredButPathLives` and its AC-5 comment encode a prior captain-approved decision to RETAIN the legacy path until BOTH (1) no live lane drives it AND (2) no targeted runtime exposes TeamCreate — and state condition (2) is UNMET ("a user on installer-stable (2.1.170, still TeamCreate-capable) still hits it; deletion is a separate trigger"). q4's captain directive reverses this. The gate/captain must consciously ratify: either accept that installer-stable users drop to the auto-team-or-bare path (treating the auto-team model as the go-forward floor), OR require a live drive on a TeamCreate-capable runtime before landing. Recommendation: the former IF the captain confirms installer-stable Claude is no longer a supported target; otherwise the latter, and this ideation should not be gate-approved for implementation until that live drive is scheduled.

### Expected surface

- Prose (measured FO surface): DELETE using-legacy SKILL.md (−14065 B); `claude-fo-dispatch.md` (−override line ~−330 B, +rename ~+90 B); `claude-first-officer-runtime.md` (−2 legacy sentences + rename ≈ −250 B); `pi-first-officer-runtime.md` (+~13 B); `pi-ensign-runtime.md` (+~13 B, not ratcheted). Net measured prose ≈ **−14.5 KB**, dominated by using-legacy.
- Go co-edits (q4): `claudeteam.go` (2 advisory rewrites, ~net-neutral) + 6 contractlint test files (list removals + 1 anchor rename + 3-test inversion + re-baseline).
- **Tolerance:** q4-side ~11 files (5 prose incl. 1 delete, 6 Go). Joint with 9q4 adds ~6-8 Go + 2-3 prose. Breach to watch: the pi/codex ratchet — the rename must raise ONLY pi (codex has no "back-channel"); leaving `["pi"]` un-raised reds `TestFOHostPromptLoadRatchet`. And a mis-scoped verbose 9q4 rung landing net-positive would eat the claude shrink (AC-2's ≥14 KB threshold catches it).

## Stage Report: ideation

- DONE: Read 9q4's completed ideation first; q4 and 9q4 MUST land as one coordinated change and your design must dovetail with it.
  Read `dispatch-failure-retry-rung.md` in full; adopted its co-edit resolutions #1-#6, the durable `### Dispatch Retries` ledger, and the "land as ONE commit, never q4-alone-first" verdict. q4's design owns co-edit #6 and defers the Degraded-Mode co-edits to 9q4 (no duplication).
- DONE: Flesh out the q4 task body: Problem, Proposed approach, Out of scope, Acceptance criteria, Test plan (replace the seed).
  All four seed placeholders replaced; Problem kept (still accurate).
- DONE: Design retiring the legacy TeamCreate machinery and renaming "back-channel" → inter-agent communication; specify exact before/after prose for each shipped-contract edit.
  `### Shipped-contract prose diff` gives verbatim before/after for the skill delete, 8 `claude-fo-dispatch.md` edits, 2 `claude-first-officer-runtime.md` edits, pi FO + pi ensign renames, and the 2 Go advisories. Codex uses "mailbox" — no codex rename (a correction to the seed's "codex mailbox back-channel").
- DONE: Own the co-edits 9q4 assigned to q4: delete using-legacy SKILL.md; rewrite claudeteam.go BareModeAdvisory (co-edit #6); AC-4 byte re-baseline.
  All three owned. CORRECTION recorded: AC-4's cited mechanism (`foFunctionReferenceBaselineBytes`/`TestFOFunctionPromptSurfaceShrinks`, "measured+1", "123236/123323") is stale — grep-confirmed neither symbol exists. Real gate is per-host `foHostLoadBaselineBytes` (strictly-above; zero-slack = baseline == measured). All three hosts measured at exactly baseline today. Also extended the advisory rewrite to `PresentFalseHint` (same TeamCreate-bootstrap defect), flagged for the gate.
- DONE: Guarantee no contractlint test is transiently red under either-alone; name each test and show the joint ordering.
  `### Co-edit map` + `### Joint landing with 9q4` name every touched test (incl. `TestUsingLegacyClaudeTeamDegradedModePointersNameRealAnchor` and the `deferredSkillCores`/`## Degraded Mode` anchor) and the single-commit ordering that keeps the tree green.
- DONE: Write ACs with a value-measuring AC + AC-4 re-baseline; each AC names its falsifying change.
  AC-1..AC-5; AC-2 is the VALUE AC (claude load shrinks ≥14 KB vs the pre-change 111183, independent number); AC-4 is the corrected re-baseline. Each names its falsifier.
- DONE: Write the Test plan (contractlint + go test ./... + -race, joint green never q4-alone-first); declare expected surface + tolerance.
  Test plan + `### Expected surface` written; joint-only-green re-grounded empirically.
- DONE: Spike the riskiest unverified mechanism.
  Spiked (`### Ideation spike result`): deletions + rename in a throwaway worktree → `go build` exit 0, contractlint red set == the mapped 6-file co-edit surface with no surprise binding; rename + pi-zero-slack ratchet bindings confirmed. Proven mechanisms relied on: FO body-write pattern (9q4), the SendMessage-availability probe as the go-forward entry (this session runs it), the symbol-bound advisory tests.
- DONE: Flag for the gate: high-stakes shipped FO contract — adversarial audit + live drive; propose the concrete diff.
  `### Gate flag` raises the installer-stable/TeamCreate-capable live-drive gap (unverifiable on this session) AND that q4 SUPERSEDES the AC-5 retain-until-both-conditions governance; the concrete prose diff is in `### Shipped-contract prose diff` for gate review.

### Summary

Designed q4 to dovetail with 9q4's completed ideation and land as one commit. Two shipped-contract retirements: delete the legacy TeamCreate override + skill (14065 B), and rename "back-channel" → "inter-agent communication". A throwaway-worktree spike proved the binding map complete (contractlint red set == the mapped 6 files, `go build` clean) and re-grounded the joint-landing necessity. Three corrections to the seed surfaced and recorded: (1) AC-4 named a stale test/constant — the real gate is per-host `foHostLoadBaselineBytes`, strictly-above, all three hosts at exactly zero slack today (so the rename reds the pi ratchet and needs a pi re-baseline); (2) the seed's "codex mailbox back-channel" — codex uses "mailbox", no rename needed; (3) the advisory defect extends to `PresentFalseHint`, not just `BareModeAdvisory`. Two gate flags: the installer-stable / TeamCreate-capable-runtime team-mode live drive is unverifiable on this session, and q4 SUPERSEDES the AC-5 governance decision (`TestLegacyConsumerRetiredButPathLives`) that RETAINS the legacy path — the gate/captain must ratify.

## Gate: ideation — APPROVED (FO, captain-ratified)

- **Verdict:** approved for JOINT implementation with `9q4` (one worktree, one commit, never q4-alone-first).
- **Governance ruling (captain CL, 2026-07-22): RETIRE NOW.** The no-TeamCreate auto-team model is the go-forward floor; installer-stable Claude (2.1.170) is no longer a blocking supported target for this path. This ratifies q4 superseding `TestLegacyConsumerRetiredButPathLives`' retain-until-both-conditions decision. **No TeamCreate-capable live drive is required before landing.**
- **Validation:** detached adversarial audit (shipped FO contract, high-stakes) + `9q4`'s bare no-team live drive (runnable on this session). The installer-stable team-mode drive is NOT required per the ruling.
- **Corrections banked from ideation:** AC-4 gate is the per-host `foHostLoadBaselineBytes` map (claude ↓ ≥14 KB, pi ↑ ~13 B, codex unchanged); codex uses "mailbox" (no rename); `PresentFalseHint` rewritten alongside `BareModeAdvisory`.
- **Base:** worktree off `origin/main` (`ca136f83`), not local `main`.
