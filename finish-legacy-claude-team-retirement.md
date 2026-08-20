---
id: nr75fq7ha3nmvsegbd22cgqa
title: Finish the legacy Claude TeamCreate retirement — the contract retired it, the binary and a live proof did not
status: validation
source: "Captain CL, 2026-08-18, in chat ('file a proper legacy claude team retirement task'), after TestLiveBreakGlassShimRecovery/selected-team failed the claude-live lane on the ca9/6ht/j7j stack by demanding a team_name the shipped contract tells the FO to omit."
started: 2026-08-18T23:15:53Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-finish-legacy-claude-team-retirement
issue:
gates:
    version: 1
    records:
        - id: gate:nr75fq7ha3nmvsegbd22cgqa:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:nr75fq7ha3nmvsegbd22cgqa-backlog-1
              briefing:
                id: briefing:nr75fq7ha3nmvsegbd22cgqa:backlog:attempt-1:revision-1
                digest: sha256:0bb3e4f258b019130d5c0bd47aad6f12bf8064cf9b30f3e9080524bf89a00f14
                request-digest: sha256:0fd53f31d5636406dca65acd9596814d0da6fb9c15c45a6fb2e57882b3c8827d
                room-ref: ./finish-legacy-claude-team-retirement/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nr75fq7ha3nmvsegbd22cgqa:backlog:1
                briefing: briefing:nr75fq7ha3nmvsegbd22cgqa:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T23:15:14.692985Z"
                decision: approve
                reason: 'Captain approved in chat: ''ok dispatch 4c and nr.'' The captain raised the retirement question himself after the break-glass red, and asked for this filing; approving it into ideation closes a nondeterministic merge blocker.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:nr75fq7ha3nmvsegbd22cgqa:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:nr75fq7ha3nmvsegbd22cgqa-ideation-1
              briefing:
                id: briefing:nr75fq7ha3nmvsegbd22cgqa:ideation:attempt-1:revision-1
                digest: sha256:90c45d49442a8feffd99c50c7ae65b82211af4952c428df049dd718916edf4c1
                request-digest: sha256:e1e8c27fc3163e5742b36b5d0e8d6898b3955cec4fb43a2dd3f1e2ac27749f7d
                room-ref: ./finish-legacy-claude-team-retirement/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-18T23:31:25.020665Z"
                reason: 'Artifact changed after prepare: the ensign committed 49adbafaf pinning the edit surface and the claudeteam keep/retire split, invalidating the frozen briefing digest. Withdrawing and re-preparing on the corrected artifact per the entity-freeze rule.'
            - id: gate-attempt:nr75fq7ha3nmvsegbd22cgqa-ideation-2
              briefing:
                id: briefing:nr75fq7ha3nmvsegbd22cgqa:ideation:attempt-2:revision-1
                digest: sha256:006b5d54dc9623be3672fda8624db8a4631c040f542be3639f80f8e3fbff9967
                request-digest: sha256:24ae4c0040055876cc553f2028b24112d18d7669d3e163e045f70cf75058cac1
                room-ref: ./finish-legacy-claude-team-retirement/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:nr75fq7ha3nmvsegbd22cgqa:ideation:2
                briefing: briefing:nr75fq7ha3nmvsegbd22cgqa:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-18T23:40:22.69892Z"
                decision: revise
                reason: 'Captain rejected in chat: ''i don''t understand why nr would keep the xfail.'' The retirement scope, the keep-list, and the root-cause diagnosis are all accepted and stay. Decision 2 is not: it recommends binding the {dispatch_agent_id} slot AND registering selected-team as liveXFail, which cancel out — an xfail declares an expected failure with an owner, and there is nothing to own once the cause is fixed. Decision 3''s own finding (''the instability is one unbound template field'') removes the only justification for hedging. Drop the xfail; ship the binding alone. If it reds again after the fix, it earns a registration then, on fresh evidence.'
            - id: gate-attempt:nr75fq7ha3nmvsegbd22cgqa-ideation-3
              briefing:
                id: briefing:nr75fq7ha3nmvsegbd22cgqa:ideation:attempt-3:revision-1
                digest: sha256:2ba2b9e1dbbdee644a7d3a483a1c373fc6014f021039b7f18ab5e387aa118210
                request-digest: sha256:6262224a6da38b2486831f500144c5c63570614b39bf17cf4a85251bd441acd7
                room-ref: ./finish-legacy-claude-team-retirement/review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:nr75fq7ha3nmvsegbd22cgqa:ideation:3
                briefing: briefing:nr75fq7ha3nmvsegbd22cgqa:ideation:attempt-3:revision-1
                by: person:captain
                at: "2026-08-19T03:24:46.596981Z"
                decision: approve
                reason: 'Captain approved in chat: ''approve both.'' Accepts retiring the legacy TeamCreate surface and binding the {dispatch_agent_id} slot with no xfail, at net -500 across ~47 files. The variance was explained from both transcripts: the FO resolved {worker_key} and dropped subagent_type in the same call, so binding removes a real choice point. Both live FOs also fabricated team names off the advertised flag, making its removal an observed fix.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:nr75fq7ha3nmvsegbd22cgqa:validation
          stage: validation
          attempts:
            - id: gate-attempt:nr75fq7ha3nmvsegbd22cgqa-validation-1
              briefing:
                id: briefing:nr75fq7ha3nmvsegbd22cgqa:validation:attempt-1:revision-1
                digest: sha256:e52a98712637bd00d187e1f3d0713be2d6422d2e9daf92ca6529986919f1a2f9
                request-digest: sha256:f023742af404c81f0ddae1edb1d5cea40599cd2e8a03631aa72b3ecf96782dee
                room-ref: ./finish-legacy-claude-team-retirement/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nr75fq7ha3nmvsegbd22cgqa:validation:1
                briefing: briefing:nr75fq7ha3nmvsegbd22cgqa:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-19T05:06:33.291198Z"
                decision: revise
                reason: 'Captain rejected in chat: ''ok send it back.'' The retirement itself is accepted and validated — keep-list verified by blob hash, every removed symbol''s callers read, all four ACs falsifiable. Rejected on AC-5''s mechanism only: the fix leaves {dispatch_agent_id} as a placeholder and attaches a resolution rule in a comment. The observed failure was a DROPPED parameter, and a comment cannot stop a drop. Put spacedock:ensign literally in the value and move the override into the comment, so the common case is the pure copy the FO performed 11/11 times.'
            - id: gate-attempt:nr75fq7ha3nmvsegbd22cgqa-validation-2
              briefing:
                id: briefing:nr75fq7ha3nmvsegbd22cgqa:validation:attempt-2:revision-1
                digest: sha256:7f75a366baee9c28486b7849e3bbf757d15717350bae9c172d4b5bacacc49b28
                request-digest: sha256:aa090be10f9f1ac2a186ddcdf3fdcee304376cbd679dc973f5d676afaa9070a4
                room-ref: ./finish-legacy-claude-team-retirement/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:nr75fq7ha3nmvsegbd22cgqa:validation:2
                briefing: briefing:nr75fq7ha3nmvsegbd22cgqa:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-19T15:34:47.592951Z"
                decision: approve
                reason: 'Captain approved in chat: ''approve push and i''ll look from the PR.'' Validation PASSED with no material finding: the literal was confirmed against the running binary across absent, named and empty agent: READMEs; the fix proved stronger than claimed by removing the template''s dependence on a self-contradicting token; both polish fixes verified with falsifying changes run. Cumulative net -724 inside the declared band.'
              application:
                target-stage: done
                state: pending
mod-block: merge:pr-merge
pr: pr-merge:736
---

`ecffcedef` (#549) retired the legacy TeamCreate dispatch path and deleted `skills/using-legacy-claude-team/`. The retirement stopped at the skill boundary. The binary still carries the flag and the field, and a live proof still fails a First Officer for obeying the current contract.

## Problem

The shipped Claude adapter now states the post-retirement shape plainly (`skills/first-officer/references/claude-fo-dispatch.md:7`):

> a worker is `Agent(name=…, run_in_background=true)` **(no `team_name`)** … `spacedock dispatch build` emits this shape (`name` present, `team_name` absent, `run_in_background` true)

and at `:26` it marks the field `// absent in the normal shape; map it verbatim if the build emits it`.

Three things still contradict that:

1. **The binary keeps the legacy surface.** `internal/dispatch/build.go` still declares `TeamName *string \`json:"team_name,omitempty"\`` and `teamNamePattern`, and `dispatch build --help` still advertises `--team-name NAME  Select the legacy TeamCreate-registry dispatch shape`.
2. **A live proof demands the retired shape.** `TestLiveBreakGlassShimRecovery/selected-team` (`internal/ensigncycle/dispatch_recovery_live_test.go:116`) fails with *"the only Agent() call did not preserve selected team mode"* when the FO emits `team_name-present=false` — which is exactly what `claude-fo-dispatch.md` instructs. Its sibling `selected-bare` is already `liveXFail("claude-sonnet", …)` against `repair-sonnet-live-flakes`; `selected-team` carries no registration and hard-fails through `t.Fatalf`.
3. **It blocks merges nondeterministically.** That test reddened the `claude-live` lane on the `ca9`/`6ht`/`j7j` stack (run 32189720211) while the same commit passed the same test locally, and `main` passed it both in CI and locally. The captain merged the stack over it.

## Why this matters beyond one flaky test

This is the third live-lane test in one day that reddened a **compliant** agent. `ca9` fixed a filing recognizer that missed a newline-terminated command. `6ht` fixed a build-count bar that punished a corrective rebuild. This one fails an FO for emitting the shape its own shipped contract prescribes.

A grader that punishes the documented behavior is worse than no grader: it trains the reader to treat live reds as noise, which is exactly the habit that lets a real red through.

## Ideation must settle

1. **Whether legacy team mode is retired in fact or only in the skills.** Decide one way. If retired, `--team-name`, `TeamName`, `teamNamePattern`, and the `selected-team` live case go together, and `dispatch build --help` stops advertising a mode we do not ship. If it is a supported fallback, say what supports it and where that is written, because `claude-fo-dispatch.md` currently says the opposite.
2. **What `selected-team` should become.** Delete with the mode, re-anchor to assert the post-retirement shape (break-glass preserves `name` and `run_in_background` and emits no `team_name`), or register it under `repair-sonnet-live-flakes` beside its sibling. Deleting a live proof needs the captain's explicit approval; recommend, do not assume.
3. **Whether the break-glass path has a real defect underneath.** The failure is nondeterministic across otherwise identical runs, so before retiring the assertion, establish whether the FO's shape choice under break-glass is stable at all. A retirement that hides a real instability is not a retirement.
4. **What else `#549` left behind.** Sweep for other survivors of the same retirement — help text, fixtures, mod prose, adapter references, `dispatch reconcile --team-name`.

## Out of scope

`repair-sonnet-live-flakes` and its existing xfail registrations. The auto-team dispatch shape itself. Any change to break-glass recovery behavior beyond what settling item 3 requires.

## Value

A live lane is only worth its cost if a red means something. Today one of its tests fails agents for following the contract we ship, and it has already forced a merge decision over a red. Either the mode exists and the contract is wrong, or the mode is gone and the test is. Both are cheap; the ambiguity is not.

## Ideation findings (2026-08-18)

### The red was misdiagnosed — the proof already asserts the post-retirement shape

The filing's premise that `selected-team` "demands a team_name the shipped contract tells the FO to omit" is wrong. The team-mode oracle (`internal/ensigncycle/dispatch_recovery_assert_test.go:182`, landed `e481864e4` 2026-08-10) requires `!hasTeamName` — team_name ABSENT — plus a shaped `name`, `run_in_background=true`, and `commonShape` (`subagent_type == "spacedock:ensign"` and a non-empty description). Pulled from the CI log of run 32189720211, the red's only failing conjunct was `subagent_type=""`:

    name="spacedock-ensign-widget-task-implementation" name-present=true
    team_name-present=false run_in_background=true run-present=true
    subagent_type="" description-present=true

Every mode predicate passed. The FO under break-glass hand-fills the manual `Agent()` template (`skills/fo-dispatch-recovery/SKILL.md:28`), whose `subagent_type="{dispatch_agent_id}"` slot has NO binding instruction — and the live fixture's README declares no `agent:` field, so the FO had nothing to resolve the slot from. `dispatch build` itself defaults it (`internal/dispatch/build.go:551`: `spacedock:ensign` when the stage declares no agent); the template never says so. The run's final message reports the durable side (marker, Stage Report, path-scoped commit `9c91c5b`) completed.

### Decision 1: legacy team mode is retired in fact — finish it in the binary

Retired. What supports it: the shipped contract mandates omission (`claude-fo-dispatch.md:7`); #549 (`ecffcedef`) deleted the mode's skill and `contractlint` keeps it deleted; and the legacy envelope (team_name present, run_in_background absent) requires `TeamCreate`, which no supported host exposes — the CI floor is the merged host (unpin #390, "claude ≥2.1.178, TeamCreate gone" per `build.go:338` and `merged_team_mode_test.go`), and the current Agent tool marks `team_name` "Deprecated; ignored". No host in the support matrix can consume what `--team-name` emits. What still says otherwise — `dispatch build --help` advertising the mode, the "sunsetting" advisory, `spawn-standing`'s required TeamCreate name, and commission Step 3 — is the surface this task removes.

Live evidence that the surviving flag actively misroutes (found in cycle 2's stream pulls): in BOTH break-glass transcripts — red 32189720211 and green main 32040773989 — the sonnet FO, told only "run in team mode", passed `--team-name` to `dispatch build` with a fabricated name (`team-lead`, `team1`) before the shim failed it. Had the build succeeded, it would have emitted the LEGACY envelope (team_name present, `run_in_background` ABSENT) plus a stderr advisory, and the contract's verbatim mapping (`claude-fo-dispatch.md:26`) would have produced a blocking call with a host-ignored field instead of the named background teammate the run mandated. The advertised flag is an attractive nuisance on the exact word "team"; removal is not just cleanup, it closes a live mis-dispatch path.

### Decision 2 (revised per gate): `selected-team` — keep the proof byte-unchanged, bind the slot, NO xfail

The test is not a proof of the legacy mode (its oracle asserts team_name absence), so "delete with the mode" is a category error. It becomes: exactly what it is. No deletion, no re-anchor, no owner registration — the test stays hard-gating with its oracle bytes unchanged, and the fix is the root cause: bind the template's `{dispatch_agent_id}` slot. The cycle-1 recommendation to ALSO register it under `repair-sonnet-live-flakes` was rejected at the gate as self-cancelling — an xfail declares an expected failure someone owns fixing, and binding the slot in the same change leaves nothing to expect or own. If it reds again after the binding, that red is fresh evidence of a different cause and earns a registration then.

Concrete binding, at the deviation's entry point: in `skills/fo-dispatch-recovery/SKILL.md`, both template arms' `subagent_type="{dispatch_agent_id}"` line gains an inline comment in the template's existing comment style (the `name=` line already carries one): `// the stage's agent: field from the workflow README; spacedock:ensign when the README names none (the build helper's default, build.go rule 6)`. The value sits AT the fill point, where the omission happened.

### Decision 3: break-glass shape choice is stable; the instability is one unbound template field

Across all observed runs — CI red 32189720211, two same-commit local passes (captain-reported), main green in CI and locally — the FO's MODE choice never wavered: single Agent() call, named background shape, no team_name, durable result committed, in the red run too. The nondeterminism is sonnet's template fidelity on exactly the one slot the template leaves unbound. Prior-era failures (PR #680 run 31640122346: two-call cardinality, missing bare marker) were fixed by the `bbe3d7a05`/`43fd2e79d`/`e481864e4` oracle rewrite and are not this failure.

**The variance, explained from both streams.** Cycle 1 left open why an unbound slot passes on main yet failed on the stack tip. Both transcripts are now in evidence — red run 32189720211 and green main run 32040773989 (both claude-sonnet-5, identical plugin bytes; each stream quotes the same recovery-skill section before the call):

- Green emitted the full six-parameter call, `subagent_type: "spacedock:ensign"` included.
- Red emitted five of six: the `subagent_type` line is absent from the call entirely — not mis-valued, dropped. Yet the SAME red call resolves `{worker_key}` to `spacedock-ensign` in the name, and its transcribed prompt carries the literal `Skill(skill="spacedock:ensign")`. The red FO demonstrably possessed the ensign identity in-context and still dropped the one parameter whose slot is an undefined token.

So what varies is not environment (CI passed on main, failed on the stack tip, passed twice locally on the tip — same bytes every time) and not the FO's knowledge (both FOs had the value three ways: the template prompt's own `Skill(skill="spacedock:ensign")` literal, the worker-key rule derived from the agent id, and the plugin's registered agent-type listing on the Agent tool). What varies is per-sample handling of an underdetermined instruction: transcribing a template into a tool call, the model reaches a placeholder with no binding rule, and "resolve it from context" vs "drop the parameter" are both plausible continuations — sampling picks one per run. The split in the evidence lands exactly on the bound/unbound line: across both calls, 11 of 11 explicitly-valued template fields were emitted correctly; the sole unbound slot went 1 for 2.

Why the binding removes this variance: it deletes the choice point. With the value stated at the fill point, emitting `subagent_type` is the same copy operation the FO performed correctly on every explicit field in both runs, not an inference. Stated plainly rather than overclaimed: sampling can still violate an explicit instruction — the claim is that the identified degree of freedom (an unbound slot forcing a resolve-or-drop choice) is removed, not that sonnet becomes infallible. A post-binding red would be a different defect on fresh evidence, handled then.

### Decision 4: survivor sweep (what #549 left behind)

- `internal/dispatch/build.go` — `TeamName` opt + `--team-name` wiring, `teamNamePattern`, legacy envelope branch, team-prefixed dispatch filename, `LegacyTeamNameAdvisory` call, input-schema `team_name` entry.
- `internal/dispatch/dispatch.go` — flag registration, value-flag list, `build --help` advertising the mode, input-JSON key doc, `spawn-standing` routing.
- `internal/dispatch/standing.go` — `runSpawnStanding` (singular; hard-REQUIRES a TeamCreate name: "Call TeamCreate first", line 81) and the legacy `--team` branch of `spawn-standing-all` with its `MemberExists` config dedup. The shipped contract only calls `spawn-standing-all` with no `--team` (`claude-fo-dispatch.md:50`) — the singular form is contract-orphaned, same status as the TeamDelete machinery #549's captain ruling retired.
- `internal/claudeteam/claudeteam.go:89` — `LegacyTeamNameAdvisory`; `internal/claudeteam/standing.go:71` — `MemberExists` (orphaned once the legacy dedup goes).
- **`skills/commission/SKILL.md:672-678` (Step 3 — Team Probe)** — still MANDATES `TeamCreate(...)` and forwarding `team_name` into every dispatch input, calls skipping it "the failure mode", and cites a "Team Creation section" of an adapter #549 deleted. Directly contradicts the shipped FO contract; the worst survivor.
- Skill prose: `claude-fo-dispatch.md:26` (map `team_name` verbatim), `fo-dispatch-core.md:158` (`--team-name` in the build synopsis), `codex-first-officer-runtime.md:20`, `fo-gate-lifecycle/SKILL.md:67`, `fo-dispatch-recovery/SKILL.md:15,25` ("omit team_name").
- Tests carrying the legacy shape as their vehicle: `build_team_name_advisory_test.go`, `build_teamname_path_test.go` (whole files); `"team_name": "fixture-team"` as the parity crossproduct's BASE stdin (six `+team` goldens embed the legacy envelope and the WARN advisory); `gate_ceremony_count_test.go` drives `--team-name` at 3 call sites; team-name cap/pattern cases in `build_errors_test.go`; legacy `spawn_standing_all_test.go`; `skills/integration/dispatch_test.go` (2 `"team_name": "fixture-team"` stdin sites); scattered stdin lines in ~15 more files.
- Scope trap, measured: the grep universe is much larger than the edit surface. 52 test files match a legacy string but many are keeps (absence assertions in `merged_team_mode_*`, oracle RED controls, pi ban lists); `docs/dev/.spacedock-state/_archive|_evidence`, `_reviews`, `_debriefs`, and `docs/roadmap/**` are historical records — never edited. Live non-test Go surface EDITED here is 6 files: `internal/dispatch/{build,dispatch,standing}.go`, `internal/claudeteam/{claudeteam,standing}.go`, plus the live-test arm; `internal/dispatch/reconcile.go`, `internal/claudeteam/reconcile.go`, `internal/cli/frontdoor.go` (comment), `shallow_boot_window_record.go`, `journeymetrics/claude.go`, `piruntime/teams.go` all match but are keeps per the list above.
- `internal/claudeteam` scope, explicit: the package is the SURVIVING auto-team seam, not legacy — `reconcile.go` (leadSessionId auto-discovery), `contextbudget.go`, `pyjson.go`, `BareModeAdvisory`, and `RenderStandingTeammatesSection` all serve the shipped merged-mode contract and stay. Exactly two members retire with #549's mode: `LegacyTeamNameAdvisory` (`claudeteam.go:89`, only trigger is the removed `--team-name`) and `MemberExists` (`standing.go:71`, sole caller is the legacy `spawn-standing-all --team` dedup branch at `internal/dispatch/standing.go:216`).
- NOT survivors (keep): `dispatch reconcile --team-name` (names an auto-team `config.json` in the still-real teams registry — identity selection, not the TeamCreate envelope); `journeymetrics`/`shallow_boot_window_record` TeamCreate tolerance (measurement over historical transcripts, incl. `eager-team-boot.stream.jsonl`); pi/codex banned-tool lists (they BAN the legacy tools); `merged_team_mode_*` and `boot_probe_parity` absence assertions; roadmap docs (history).
- Related open filing: `align-claude-break-glass-agent-proof` (`s2atxdv146qknetdjvx0xer6`, backlog, PR #680 era) — its two ACs are embodied by the current oracle; recommend the captain close or fold it.

### Proposed approach

1. Land the root-cause fix first: the `{dispatch_agent_id}` slot binding in both break-glass template arms (inline comment at the fill point, Decision 2). No xfail, no live-test or registry change — `selected-team` stays hard-gating and byte-unchanged.
2. Retire the binary surface: remove `--team-name`/`TeamName`/`teamNamePattern`/legacy envelope/team-prefix filename/advisory from `dispatch build`; remove `spawn-standing` (singular) and the `--team` legacy branch of `spawn-standing-all` (+ orphaned `MemberExists`); keep `reconcile --team-name`.
3. Re-anchor the test fleet: delete the two legacy-only test files; drop `team_name` from base stdins; regenerate goldens (crossproduct `+team` rows become the merged shape — names kept, "team" now denotes the auto-team); drop `--team-name` from `gate_ceremony_count_test.go`; port mode-neutral validation cases from the legacy standing test into the merged one.
4. Extend the retirement invariant: `contractlint` also asserts no shipped skill text instructs `TeamCreate`; a focused CLI test asserts `dispatch build --team-name` now exits 2 (unknown flag).
5. Rewrite commission Step 3 from the TeamCreate probe to the shipped boot probe (SendMessage availability per `claude-fo-dispatch.md:7`).

New-mechanism justification: contractlint-over-skills serves AC-4 (the exact regression — a skill re-teaching TeamCreate — already survived #549's sweep in commission; review alone missed it). The slot binding serves AC-5 (alternatives: loosening the oracle by dropping `subagent_type` stops grading a field the template mandates and hides the defect instead of fixing it; an owner-bound xfail was rejected at the gate as self-cancelling alongside the fix).

### Expected surface

Estimate net LOC change: -500, across ~47 files. Insertions ≈ +115, deletions ≈ +615. Breakout: product (6 Go files + 6 shipped skill files) net ≈ -180 (ins +45 / del +225); tests + fixtures (~27 test files incl. `skills/integration/dispatch_test.go`, ~10 goldens) net ≈ -320 (ins +70 / del +390 — test-side additions already budgeted at ~2x first instinct per today's calibration); docs net ≈ ±0 (help-text diff counted in product). Tolerance: net within [-250, -750]; files ≤ 60. Declared in net — this task's purpose is removal; gross would count the deletions as growth. Historical trees (`_archive`, `_evidence`, `_reviews`, `_debriefs`, `docs/roadmap`) are excluded from the surface by classification, not by accident. Revision delta from the cycle-1 estimate: -2 files (`dispatch_recovery_live_test.go` and `runtime-live-ci-registry.md` are no longer touched — the xfail was dropped); net unchanged within noise.

### Declared semantics

- Command grammar: `dispatch build --team-name` becomes an unknown-flag usage error (exit 2). `dispatch spawn-standing` (singular) becomes an unknown subcommand. `spawn-standing-all` loses `--team`. A stdin `team_name` key degrades to the ignore-unknown-keys path (spiked, see below) — same input with and without it emits identical envelopes.
- Stored formats: the build envelope can no longer carry `team_name`; dispatch file names lose the `{team}-` prefix shape (merged session-token prefix unchanged).
- Runtime behavior: none on supported hosts (the removed shapes are unreachable via the shipped contract); `selected-team`'s grading is unchanged (hard-gating), and the break-glass template text gains the slot binding.
- Authority: none.

### Acceptance criteria

- **AC-1 (measures the end-value):** The retirement lands as net removal — cumulative net LOC vs origin/main ≤ -250 (insertions minus deletions), measured by `git diff --shortstat origin/main...HEAD` on the landed branch. Independent baseline that can move the wrong way: compat shims or re-advertising would push it positive.
- **AC-2:** The binary refuses the retired selector: `--team-name` exits 2 with a usage error, and a stdin `team_name` key changes nothing (byte-identical envelope with and without it). Tested by a focused CLI test (exit code + stderr) and an envelope byte-equality pair.
- **AC-3:** No shipped dispatch path emits the legacy envelope: every non-bare host=claude row of the parity crossproduct and merged-mode suites emits `name` + `run_in_background:true` and never `team_name` (regenerated goldens are the byte-compare proof; re-introducing an emission fails them).
- **AC-4:** The shipped instruction surface contains no TeamCreate imperative: commission's pilot-run probe instructs the SendMessage-availability probe, and the extended `contractlint` invariant fails on any shipped skill text instructing `TeamCreate` (paired with AC-1/AC-3 as their instruction-side enforcement).
- **AC-5:** The break-glass template's `{dispatch_agent_id}` slot is bound in both arms (inline comment at the fill point: stage `agent:` field; `spacedock:ensign` when the README names none), removing the resolve-or-drop choice point run 32189720211 red on — while `selected-team` remains hard-gating and byte-unchanged: no `liveXFail`, no registry Expected-failure line, oracle and its RED controls (two-call, team_name-present, bare/team cross) untouched and still failing wrong topologies. Tested by: the RED-control unit suite unchanged-green, `live_registry_reconciliation` green with no selected-team gap row, and the template diff shipping in the same change. No live run required to land; a post-binding live red would be fresh evidence of a different cause and earns its own filing then.

### Test plan

Offline only; no live run needed to land. Focused Go CLI tests (flag refusal, stdin byte-equality), regenerated goldens (crossproduct, nonascii-title, namecap, ceremony fixtures), contractlint extension, existing oracle RED-control unit tests untouched, `live_registry_reconciliation` staying green with no selected-team gap row. Full `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`. Cost: mechanical-wide, logic-shallow; the risk is missed-survivor, mitigated by the grep-complete sweep above.

### Spike record

Riskiest unverified mechanism — "an unknown stdin key is ignored by `dispatch build`" (the post-removal fate of `team_name` input) — exercised against the 0.27.0-pre8 binary: stdin carrying `"bogus_unknown_key":"zzz"` exited 0 with a clean envelope. Remaining mechanisms are proven in-repo: `liveXFail`/`finishLiveScenario` (sibling arm), registry reconciliation lint (existing test), golden regen harness (existing).

### Proposed doc diff (user-visible surfaces)

`dispatch build --help` (internal/dispatch/dispatch.go):

    -  --team-name NAME              Select the legacy TeamCreate-registry dispatch shape. On host=claude, auto-team is the default — omit this unless you mean legacy team mode.
    -  --bare-mode                   Emit the bare sequential shape (no name, no team_name, no run_in_background); unsupported on host=codex.
    +  --bare-mode                   Emit the bare sequential shape (no name, no run_in_background); unsupported on host=codex.

with `--team-name` also dropped from the flag summary line and the optional input-JSON key list (`team_name` removed). `docs/runtime-live-ci-registry.md` is NOT touched (revised: no Expected-failure line — `selected-team` stays hard-gating).

### Open questions for the gate

1. `spawn-standing` (singular) + `spawn-standing-all --team` retirement widens beyond the literal `--team-name` sweep (contract-orphaned standing legacy, per the #549 TeamDelete precedent) — default IN, trim if the gate says so.
2. The AC-5 template slot binding is a break-glass instruction change; judged within "what settling item 3 requires" (it is the identified root cause) — flag if the gate reads the carve-out narrower.
3. `align-claude-break-glass-agent-proof` (`s2a…`) looks superseded by the current oracle — recommend closing it.

(Resolved at the cycle-1 gate: `selected-team` is kept, not deleted; the xfail half was rejected — bind the slot only, no registration.)

## Stage Report: ideation

- DONE: Settle whether legacy team mode is retired in fact or only in the skills. One answer, and name what supports it — claude-fo-dispatch.md says the normal shape omits team_name, while dispatch build --help still advertises --team-name as the legacy TeamCreate shape.
  Retired in fact (Decision 1): the legacy envelope requires TeamCreate, absent from every supported host (merged CI floor, #390; build.go:338); the binary/help/commission surfaces are the leftovers to remove.
- DONE: Say what TestLiveBreakGlassShimRecovery/selected-team should become: delete it with the mode, re-anchor it to assert the post-retirement shape, or register it under repair-sonnet-live-flakes beside its already-xfailed sibling. Deleting a live proof requires the captain's explicit approval — recommend, do not assume.
  It ALREADY asserts the post-retirement shape (oracle requires team_name absent, dispatch_recovery_assert_test.go:182); recommend keep + register under owner 060xp69y61yhrww23g3wvwqy + bind the template's {dispatch_agent_id} slot. No deletion requested.
- DONE: Establish whether the FO's dispatch-shape choice under break-glass is stable at all before retiring the assertion. It failed in CI and passed twice locally on the same commit. A retirement that hides a real instability is not a retirement.
  Stable in mode across all observed runs incl. the red (CI log, run 32189720211: only failing conjunct was subagent_type="" — verified from the failure detail line); the instability is one unbound template slot, an identified product gap the task fixes, not an unknown.
- DONE: Declare net LOC change and file count with a realistic test budget, and sweep for other survivors of #549 (help text, fixtures, mod prose, adapter references, dispatch reconcile --team-name).
  Estimate net -500 across ~48 files (ins ≈ +120 / del ≈ +620; product -180 / tests -320 / docs ±0; tolerance net [-250,-750], files ≤ 60); sweep enumerated in Decision 4 — worst survivor is commission Step 3 still mandating TeamCreate; reconcile --team-name classified keep (auto-team identity, not the legacy envelope).

### Summary

Corrected the filing's central premise: the live proof does not demand a team_name — the CI red's only failing predicate was an empty subagent_type, traced to an unbound {dispatch_agent_id} slot in the break-glass manual template. Settled retirement as fact (no supported host can consume the legacy envelope), scoped the binary/tests/skills removal with a keep-list (reconcile --team-name, historical-transcript tolerances), and spiked the one unverified mechanism (unknown stdin keys are ignored, exercised against the pre8 binary). Recommendation to the gate: keep selected-team, xfail-bind it under the existing sonnet-flake owner, and ship the slot binding as the root-cause fix.

## Stage Report: ideation (cycle 2)

- DONE: Revise Decision 2 to bind-the-slot-only — no xfail registration.
  Decision 2 rewritten (test stays hard-gating, byte-unchanged; binding is an inline comment at the fill point in both template arms); AC-5 merged accordingly, registry/live-test edits removed from approach, doc diff, semantics, and estimate (-2 files, now ~47).
- DONE: Explain the variance the cycle-1 body left open — why an unbound slot passes on main yet failed on the stack tip.
  Pulled both transcripts (red 32189720211, green main 32040773989): green emitted all six template fields; red emitted five, dropping only subagent_type while resolving {worker_key} to spacedock-ensign in the same call — the FO possessed the value and dropped the unbound parameter. Variance is per-sample resolve-or-drop handling of an undefined placeholder, not environment (same bytes across CI-main pass / CI-tip fail / local passes) and not missing knowledge (the value was in-context three ways). Across both calls: 11/11 explicit fields correct, 1/2 on the sole unbound slot — the split lands exactly on the bound/unbound line. Binding removes the choice point, turning the fill into the copy operation both runs performed correctly on every explicit field; stated with its honest bound (sampling can still defy an explicit literal — a post-binding red is fresh evidence, filed then).
- DONE: Report any delta from net -500 / ~49 files.
  Net -500 unchanged within noise; ~47 files (dispatch_recovery_live_test.go and runtime-live-ci-registry.md no longer touched).

### Summary

Dropped the xfail per the gate and answered the variance puzzle with direct stream evidence from the red and a green main run: the nondeterminism is the model's per-sample handling of the one template slot with no binding rule — it either resolves it from ambient context or drops the parameter — and binding the slot at the fill point deletes that choice. Bonus finding folded into Decision 1: both live FOs, told only "team mode", reached for the advertised --team-name with fabricated names, so the flag's removal closes an observed mis-dispatch path, not a theoretical one.

### Feedback Cycles

- Cycle 1: SURFACE OVERRUN DISCLOSED (no rejection) — implementation self-report; surface 77 files/net -751 vs estimate ~47 files (cap ≤60)/net -500, tolerance [-250,-750]. LOC lands one line past the tolerance FLOOR (-751 vs -750) — the correct direction for a removal — while the file count exceeds the ≤60 cap by 17, concentrated entirely in tests and goldens (66 of 77 vs ~37 estimated). Cause, counted directly: `team_name: "fixture-team"` decoration appeared as inert base-stdin filler in ~24 test files beyond the crossproduct family ideation's grep sweep anticipated, and each such test captures its `dispatch_file_path` byte-for-byte in its own golden, so making the field a no-op forced far more golden regeneration than the estimated ~10. Keep-list held: `_archive`/`_evidence`/`_reviews`/`_debriefs`/`docs/roadmap`, both `reconcile.go` files, the live test, the oracle, and the registry are all byte-untouched (FO-verified independently against `main`). Two items validation must weigh rather than accept: (1) the implementation added a package-level `TestMain` in `internal/dispatch` clearing `$CLAUDE_CODE_SESSION_ID`, because merged-mode dispatch keys the golden filename on it. CORRECTED BY VALIDATION (2026-08-19): this was recorded here and in the implementation report as "a pre-existing bug fixed in passing." It is not. `main` with the var exported is 0 failures — the env-keyed branch pre-existed but no golden reached it, because every fixture's `team_name` diverted to the legacy prefix branch. Making the key inert is what made the path reachable, so this is a hermeticity gap this change INTRODUCED and fixed in the same commit. The fix is correct and minimal either way; the characterization was wrong, and the FO wrote it into this line as unrelated goodwill. (2) It DEFERRED two decorative, now-inert `team_name` stdin keys in `internal/ensigncycle` (`cycle_test.go`, `feedback_test.go`), disclosed rather than silently left, to avoid growing the file count.

- Cycle 2: REJECTED — captain at the validation gate; surface 77 files/net -751 (unchanged this round); AC unchanged. Validation recommended PASSED with no material finding, and the retirement itself is accepted: keep-list verified by blob hash against f64c733ef, every removed symbol's former callers read rather than grepped, the one non-obvious deletion verified differentially against main, and all four ACs carrying named falsifying changes that were actually run. Rejected on AC-5's MECHANISM only. The fix leaves `subagent_type="{dispatch_agent_id}"` as a placeholder and attaches a resolution rule in a comment. The observed failure was a DROPPED parameter, not a wrongly computed one, and a comment cannot stop a drop — the FO resolved `{worker_key}` correctly in the same call it omitted `subagent_type` (11/11 on explicit fields, 1/2 on the sole placeholder). Validation independently reached the same conclusion and named the stronger form. FO context established at the gate: `{dispatch_agent_id}` predates this task (introduced by #466, present identically on main), the `subagent_type` field itself is sound in both the Agent tool and the binary, and break-glass is the ONLY path where a placeholder reaches a tool call — every normal dispatch has `dispatch build` fill the field via rule 6 (build.go:551). So the template is asking the FO to re-derive a default the binary already owns. Recorded cost of the fix: the literal duplicates rule 6, and the oracle hardcodes it too, making three copies of one default.

## Stage Report: implementation

- DONE: Remove the legacy TeamCreate surface per the settled sweep: --team-name, TeamName, teamNamePattern, LegacyTeamNameAdvisory, MemberExists, spawn-standing singular and --team, and the help text that advertises the mode. Respect the keep-list exactly.
  `internal/dispatch/{build,dispatch,standing}.go` and `internal/claudeteam/{claudeteam,standing}.go` (commit 21b7609f7); `dispatch build --team-name` now exits 2 (usage error), a stdin `team_name` key is ignored, `spawn-standing` (singular) and `spawn-standing-all --team` are gone. Keep-list verified byte-for-byte untouched: `internal/dispatch/reconcile.go`, `internal/claudeteam/reconcile.go`, `internal/cli/frontdoor.go`, `internal/piruntime/teams.go`, `internal/journeymetrics/claude.go`, `internal/ensigncycle/shallow_boot_window_record.go`, `merged_team_mode_*`, `boot_probe_parity_test.go` — none show in `git diff HEAD~1 HEAD`.
- DONE: Bind the {dispatch_agent_id} slot at the fill point in BOTH template arms, so the FO copies a value rather than choosing whether to resolve one. selected-team stays hard-gating and byte-unchanged — no xfail, no registry edit, no live-test edit.
  `skills/fo-dispatch-recovery/SKILL.md` both `Agent()` arms gain `// the stage's agent: field from the workflow README; spacedock:ensign when the README names none (the build helper's default, build.go rule 6)` at the exact `subagent_type="{dispatch_agent_id}",` line. `internal/ensigncycle/dispatch_recovery_assert_test.go`, `dispatch_recovery_live_test.go`, and `docs/runtime-live-ci-registry.md` show empty diffs against `HEAD~1` — confirmed untouched.
- DONE: Prove the removal does not break the surviving auto-team seam: internal/claudeteam reconcile, leadSessionId discovery, contextbudget, BareModeAdvisory and standing-teammate rendering all keep working, and reconcile --team-name still resolves auto-team identity.
  `go test ./...` green except one pre-existing, environment-local failure unrelated to this change (`TestCodexResolveManifestAgainstInstalledHost`, fails identically on `main` — a local Codex plugin-cache probe). `internal/claudeteam` (reconcile, contextbudget, standing render), `internal/dispatch` (reconcile_test.go, reconcile_session_test.go, reconcile_namecap_test.go — `reconcile --team-name` resolution paths), and `internal/ensigncycle` (326s, incl. `merged_team_mode_*`) all pass unchanged. `go test ./internal/dispatch/... ./internal/claudeteam/... ./internal/cli/... ./internal/contractlint/... -race` clean (same pre-existing failure only). Fixed a real bug this removal would otherwise have exposed as test flakiness: merged-mode dispatch keys the filename on `$CLAUDE_CODE_SESSION_ID` when set, and with legacy `team_name` gone every non-bare claude fixture now takes that path — a developer shell exporting the var (mine did) would leak a live session id into golden fixtures. Added a package-level `TestMain` in `internal/dispatch` clearing it; verified against `internal/ensigncycle`, which has no such guard but already only checks body/field content, not exact dispatch filenames.
- DONE: Report actual net surface against the approved net -500 across ~47 files, broken out product / tests / docs. This is a removal — report NET, never gross.
  Net -751 across 77 files (commit 21b7609f7; `git show --stat HEAD`): product 5 Go files net -188 (ins 100/del 288, matches the -180 estimate closely); skills docs 6 files net -1 (ins 11/del 12, matches the ±0 estimate); tests 30 `_test.go` files net -455 (ins 320/del 775); golden fixtures 36 `.txt` files net -107 (ins 73/del 180). LOC lands essentially at the declared tolerance floor (-751 vs -750). File count (77) exceeds the declared ≤60 cap by 17, concentrated entirely in tests+fixtures (66 of 77 files there vs. the ~37 estimated). Cause, verified by direct count: `team_name: "fixture-team"`-style decoration appeared as harmless base-stdin filler in ~24 test files beyond the crossproduct family ideation's grep sweep anticipated (build_advance/build_hazards/build_pi_host/codex_bootstrap/build_state_no_origin/build_namecap/build_json_ergonomics/self_contained_assignment/reconcile_namecap/cycle2_parity, etc.), and each such test's `dispatch_file_path` is captured byte-for-byte in its own golden, so removing the field's effect (now a no-op per AC-2) forced far more golden regeneration than the cycle-1 estimate's ~10. Not padding: every touched file needed either the dead key removed (declared semantics require `team_name` to have zero effect) or its golden regenerated to match (AC-3's byte-compare proof) — trimming any of it back would leave stale legacy-shaped fixtures or an un-regenerated golden. Deliberately deferred, disclosed rather than silently fixed: two `internal/ensigncycle` fixture helpers (`cycle_test.go`, `feedback_test.go`) still carry a decorative, now-inert `team_name` stdin key — already harmless no-ops (tests pass), left alone to avoid growing the file count further for zero functional gain. Flagging the tolerance/cap deviation for the gate rather than trimming correct, necessary changes to hit the number.

### Summary

Removed the legacy TeamCreate-registry dispatch surface from the binary (`--team-name`/`TeamName`/`teamNamePattern`/`LegacyTeamNameAdvisory`/`MemberExists`/`spawn-standing` singular/`spawn-standing-all --team`) and bound the break-glass template's `{dispatch_agent_id}` slot at its fill point in both arms — the identified root cause of the `selected-team` live red, landed with the oracle byte-unchanged and no xfail. Rewrote commission's Step 3 from a TeamCreate probe to the shipped SendMessage-availability boot probe and dropped stale `--team-name` mentions from four other skill files; added a contractlint invariant (AC-4) that fails on any shipped skill instructing `TeamCreate(...)`, verified falsifiable against the pre-fix commission text. Ported the mode-neutral spawn-standing validation cases (fable model, non-ASCII prompt escaping, missing-model/bad-enum/trailing-heading) onto `spawn-standing-all` before deleting the singular-only legacy test files, and added dedicated AC-2 coverage (flag refusal, stdin byte-equality). `go test ./...` and `-race` are green (one pre-existing, unrelated, environment-local failure). Net surface landed at -751/77 files against the approved -500/~47 estimate: LOC is essentially at the tolerance floor, but the file count exceeds the ≤60 cap by 17, entirely in the test+golden surface — real, necessary work the ideation's grep-based estimate undercounted, not scope creep; reported for the gate to weigh rather than trimmed to fit.

## Stage Report: validation

- DONE: Verify the retirement removed only what was authorized and nothing that was not. Check the keep-list byte-for-byte against main yourself: both reconcile.go files, frontdoor.go, piruntime/teams.go, journeymetrics/claude.go, shallow_boot_window_record.go, merged_team_mode_*, boot_probe_parity_test.go, the oracle (dispatch_recovery_assert_test.go), the live test, and runtime-live-ci-registry.md. Also confirm no historical tree (_archive, _evidence, _reviews, _debriefs, docs/roadmap) was edited.
  All 12 keep-list paths blob-hash-identical to `f64c733ef` (`git rev-parse BASE:path` vs `HEAD:path`, not a read); diff touches only `internal/` and `skills/` — zero `_archive|_evidence|_reviews|_debriefs|docs/roadmap|.spacedock-state` paths.
- DONE: Verify the surviving auto-team seam still works, not just that its tests pass: reconcile --team-name resolves auto-team identity, leadSessionId discovery, contextbudget, BareModeAdvisory, and standing-teammate rendering. Removing LegacyTeamNameAdvisory and MemberExists must not have taken a live caller with them — find every former caller and confirm each is genuinely dead, not merely untested.
  Callers read, not grepped: `MemberExists`'s sole production caller was the removed legacy dedup at `dispatch/standing.go:216`; `LegacyTeamNameAdvisory`'s was the removed `teamName != ""` branch at `build.go:799`; `isFile` went with it while `readTeamConfig`/`hasMember` stay live in `contextbudget.go` (no dead code, `go vet` clean). `reconcile --team-name` exercised on both binaries: nonexistent team → identical exit 1 + identical stderr; auto-discovery resolved the real `session-fedfe9c3` identically. The non-obvious deletion — build.go point 9's show-standing injection — verified differentially, not by reading its claim: a merged build against `docs/dev` (which HAS a standing mod) is byte-identical base↔HEAD apart from the launcher path and neither emits show-standing, while a BASE build with `team_name` emitted 2 fetch commands including it. Falsifying change: `EnumerateDeclaredStandingTeammates` returning non-nil at `teamName == ""` would have made those two envelopes differ.
- FAILED: Test the TestMain bug fix as a necessity claim, and hunt for what it might mask. The implementation added a package-level TestMain in internal/dispatch clearing CLAUDE_CODE_SESSION_ID because merged-mode dispatch keys the golden filename on it. Confirm the leak is real, that clearing it is the minimal fix rather than hiding a behavior change, and CRITICALLY that no golden now passes only because the variable is cleared — a golden that would differ under a real session is a false green.
  Fix is correct and necessary; the DISCLOSURE's characterization is not, which is what fails here. Leak real: TestMain neutered + var exported → 28 failures; var unset → green. But `main` with the var exported is 0 failures — the env-keyed branch pre-existed yet no golden reached it, because every fixture's `team_name` diverted to the legacy prefix branch. Making the key inert is what made the path reachable, so this is not "a pre-existing bug fixed in passing" (Cycle 1 + implementation report) but a hermeticity gap this change introduced and fixed in the same commit. Minimal: clears one var, matching the harness's existing CLAUDECODE/CODEX_THREAD_ID/PI_* normalization; product behavior untouched (`else if mergedMode` → `if mergedMode`, same body). Masking: none. Mutation-proved — disabling the session prefix entirely leaves ALL golden tests green, so they never covered the disambiguator, while `TestBuildMergedModeDispatchFileDisambiguator` reds on that exact mutation via its own `t.Setenv` (asserts distinct paths + filename embeds `sessionaaa`). Retired team-keyed collision coverage was replaced, not dropped; its merged-floor equivalent pre-existed on main.
- DONE: Check the deferred pair and the 36 regenerated goldens. Two inert team_name keys remain in internal/ensigncycle (cycle_test.go, feedback_test.go): confirm they are genuinely no-ops post-change, not latent. Then spot-check regenerated goldens against the declared semantics — a stdin team_name key must have ZERO effect, so a golden whose dispatch_file_path still carries a team-shaped prefix is a real finding.
  Deferred keys (3, not 2 — feedback_test.go carries two) removed in a throwaway: ensigncycle stays green, so genuine no-ops; reinforced by the full suite passing with the session var exported (ensigncycle has no TestMain, so any effect there would have surfaced). Goldens swept exhaustively, not spot-checked: zero contain `team_name`, zero carry a team-shaped dispatch-path prefix, zero retain the WARN advisory; every non-bare non-error claude golden has `name` + `run_in_background: true`. All 36 diffs confined to retired-shape lines (path, prompt pointer, team_name, WARN, run_in_background) plus one JSON trailing comma and `build-mods.txt`'s show-standing line, both explained above.

### Acceptance criteria

- AC-1 PASS: `git diff --shortstat origin/main...HEAD` = 77 files, +504/-1255, net **-751** ≤ -250.
- AC-2 PASS: all four argv forms (`--team-name=v`, `--team-name v`, bare trailing, flag-first) exit 2, empty stdout, stderr naming the flag; stdin key byte-identical incl. adversarial values (null/number/object/array/`../../etc/passwd`/300-char) — all inert, no panic. The tombstone is load-bearing: generic unknown flags exit 0 silently, so without it `--team-name v` would have emitted a merged envelope with no warning. Falsifying change: removing the tombstone reds `TestBuildTeamNameFlagRefused` (exit=0, want 2) — verified.
- AC-3 PASS: goldens are the byte proof and they bite. Falsifying change: suppressing `run_in_background` on the merged shape reds all 7 non-bare crossproduct rows + `TestBuildMergedModeEmission` — verified.
- AC-4 PASS: commission Step 3 now the SendMessage-availability probe. Falsifying change: injecting `TeamCreate(` into a shipped skill reds `TestNoShippedSkillInstructsTeamCreate`; run against **main's** skills tree it flags `skills/commission/SKILL.md` — it catches the real historical regression, not just its own text. Exactly one bare-word `TeamCreate` remains in all shipped skills, legitimate prose about omitting `reconcile --team-name`.
- AC-5 PASS on its stated terms: binding present in both `Agent()` arms; oracle, live test, and registry blob-identical; no `selected-team` xfail (only sibling `selected-bare`); RED-control unit suite green incl. `selected_bare_rejects_team_name_present`; `TestRuntimeLiveRegistryReconciliation` green at `xfail/claude-sonnet=1` with no selected-team gap row. See the prose assessment below — the AC waives live proof, so this is not a rejection ground.
- Tests: `go test ./...` (var unset), `go test ./...` (var exported), `go test ./... -race` — each exactly one failure, `TestCodexResolveManifestAgainstInstalledHost`, which fails identically on `main` with the same message (local `~/.codex` plugin-cache probe); 0 data races. `gofmt -l ./cmd ./internal` clean.

### The AC-5 root-cause fix is prose — assessed as prose

The comment supplies a resolution rule where the token had none, so it is a real improvement. It does **not** remove the choice point; it replaces "resolve from nothing" with a two-branch lookup (read the stage's `agent:`, else default). The observed failure was *dropping the parameter*, and a comment cannot prevent a drop. A strictly stronger form was available and not taken: put the literal in the value — `subagent_type="spacedock:ensign"` with the override in the comment — making the common case the pure copy the FO performed 11/11 times, with no lookup. As shipped the mechanism is prose, and this repo's own record (three live reds in one day from prose-shaped contracts) says prose is weaker than a binary refusal. Ideation states this bound honestly and the captain approved "no live run required to land", so it stands as a deferred risk, promoted by a post-binding `selected-team` red.

### Review-finding disposition

No material finding. All below are Polish or deferred risk; no value AC fails and the delivered code is correct.

1. Polish — "a pre-existing bug fixed in passing" is inaccurate (evidence in checklist item 3). Matters only because the gate is asked to weigh it as unrelated goodwill when it is a consequence of this change. Recommend the Cycle line and gate wording be corrected; no code change.
2. Polish — `spawn-standing-all` silently ignores a legacy `--team` (exit 0, merged output) while `build` loudly refuses `--team-name` (exit 2). Inconsistent migration story for two flags retired in one commit; declared semantics only said `spawn-standing-all` "loses `--team`", so no AC fails.
3. Polish — ideation Decision 4 listed `fo-dispatch-recovery/SKILL.md:15,25` ("omit team_name") in the sweep; both remain, undisclosed. On the merits, keep them: they reinforce the oracle's `!hasTeamName` requirement and counteract the FO's observed habit of inventing team names. Record as a deliberate keep, not an open item.
4. Polish — `EnumerateDeclaredStandingTeammates` keeps a vestigial `teamName` parameter whose only caller passes the sentinel `"_show_standing_"` purely to defeat the `teamName == ""` early return; its doc comment still says "Returns an empty slice for bare mode (empty teamName)" — stale prose for a retired mode.
5. Deferred risk — the oracle hardcodes `subagentType == "spacedock:ensign"` while the new comment tells the FO to prefer the stage's `agent:` field. They agree only because the break-glass fixture's README declares none. Promote if a live fixture ever declares `agent:` — a compliant FO would then be graded red.
6. Deferred risk — `buildSpawnSpec`'s `if !isFile(modPath)` branch is now unreachable (the only caller enumerates existing `_mods` files) and its golden was deleted with the singular subcommand. Promote if any caller passes an arbitrary path again.
7. Deferred risk (pre-existing, unchanged) — on codex/pi the dispatch filename now carries no disambiguator; main's team prefix was host-independent. The contract never passes `team_name` there, so the supported path is unchanged.

Worth recording positively: removing the team prefix deletes a user-controlled path component from a `/tmp` filename outright (main only pattern-guarded it).

### Summary

Recommendation: **PASSED**, with the file-count breach carried to the captain as a scope note rather than a rejection. I verified the keep-list by blob hash rather than by reading the implementation's account, and the one non-obvious deletion (build.go's show-standing injection) differentially against `main` — it fired only in the deleted legacy branch, so merged behavior is byte-unchanged. Every removed symbol's former callers were read and each is genuinely dead, not merely untested. The 77-vs-≤60 overage is real but every file in it is necessary: exactly 24 test files each deleted at least one inert `team_name` key, which the declared "zero effect" semantics require, and all 36 golden diffs are confined to retired-shape lines — I found nothing trimmable. LOC landed one line past the tolerance floor in the correct direction for a removal, and no AC declares the file cap. One item fails as recorded: the TestMain change is the right fix but is misdescribed as pre-existing — `main` is green with `CLAUDE_CODE_SESSION_ID` exported, so this change made the path reachable and then fixed it. The "false green" hunt came back clean by mutation: the goldens never covered the session disambiguator, and the dedicated test that does still reds when it breaks.

## Stage Report: implementation (cycle 2)

- DONE: Replace AC-5's mechanism. Both break-glass `Agent()` arms in `skills/fo-dispatch-recovery/SKILL.md` now read `subagent_type="spacedock:ensign"` (literal) with an inline override comment, not `subagent_type="{dispatch_agent_id}"` plus an explanatory comment. Oracle, live test, and registry stay untouched — no xfail, no registry edit, no live-test edit.
  Commit `8b8d009ea`. `git diff HEAD~1 HEAD -- internal/ensigncycle/dispatch_recovery_assert_test.go internal/ensigncycle/dispatch_recovery_live_test.go docs/runtime-live-ci-registry.md` is empty. `go test ./internal/dispatch/... ./internal/ensigncycle/...` and `-race` green (unchanged pre-existing `TestCodexResolveManifestAgainstInstalledHost` failure only).
- DONE: Record the maintenance cost honestly — no sync mechanism added.
  The literal now duplicates rule 6 (`build.go:536`, `subagentType := "spacedock:ensign"`) and the oracle hardcodes the same string a third time (`dispatch_recovery_assert_test.go:176`, unchanged, explicitly out of scope). One default now lives in three places; that trade was weighed at the gate and the literal won, deliberately, over adding a mechanism to keep them in sync.
- DONE: Two cheap polish items from validation's list.
  `EnumerateDeclaredStandingTeammates(workflowDir, teamName)` → `EnumerateDeclaredStandingTeammates(workflowDir)`: dropped the vestigial `teamName` parameter and the `teamName == ""` early return (the sole caller, `runShowStanding`, passed the sentinel `"_show_standing_"` only to defeat it), and the stale "empty teamName" doc-comment language. `spawn-standing-all` now refuses a legacy `--team` flag (exit 2, `TestSpawnStandingAllTeamFlagRefused`) instead of silently ignoring it — consistent with `dispatch build --team-name`'s refusal, closing the inconsistency validation named.
- DONE: Accept the "pre-existing bug" correction without re-litigating it.
  Validation is right: `main` with `CLAUDE_CODE_SESSION_ID` exported is 0 failures, so the env-keyed branch pre-existed but no golden reached it — every fixture's `team_name` diverted to the legacy prefix branch instead. Making the key inert is what made the path reachable; this cycle's report (and any future summary) should say "hermeticity gap this change introduced and fixed in the same commit," not "pre-existing." No code change follows from this — the `TestMain` fix itself was already validated correct and minimal.

### Cycle-2 surface (separate from cumulative)

Commit `8b8d009ea`: 5 files changed, +39/-12, net **+27** (not a removal — this cycle replaces a comment-based mechanism with a literal-based one, so the net is small and positive: one new dedicated test plus a slightly longer per-line literal-and-comment pair). Cumulative across both cycles vs. the approved -500/~47 estimate: net **-724** across **78** files (`git diff --shortstat f64c733ef...HEAD`) — LOC is within the declared [-250,-750] tolerance band (cycle 2's small positive delta moved it back inside the floor it had been one line past); file count is 78 vs the ≤60 cap, one more than cycle 1's 77 (the new `spawn_standing_all_merged_test.go` growth stayed in an already-touched file; no new file was added this cycle).

### Summary

Replaced AC-5's comment-based mechanism with a literal: both break-glass `Agent()` arms hardcode `subagent_type="spacedock:ensign"` with the override rule in the comment, matching build.go rule 6 and closing the exact failure mode observed live (a dropped, not miscomputed, parameter). Recorded the resulting three-copies-of-one-default maintenance cost per the gate's instruction, with no sync mechanism added. Landed both cheap polish items: removed `EnumerateDeclaredStandingTeammates`'s vestigial `teamName` parameter, and made `spawn-standing-all` refuse a legacy `--team` (exit 2) instead of silently ignoring it. Accepted validation's correction on the `TestMain`/session-leak characterization without re-litigating — it is a hermeticity gap this retirement introduced and fixed in the same commit, not a pre-existing bug. `go test ./...` and `-race` remain green (same one pre-existing, unrelated, environment-local failure). Cycle-2 surface is +27/5 files; cumulative is -724/78 files, back within the declared LOC tolerance with file count still one over cycle 1's count.

## Stage Report: validation (cycle 2)

Fresh validator; scope is the cycle-2 delta (`8b8d009ea` on `21b7609f7`) per the FO. Cycle 1's findings stand.

- SKIPPED: Verify the retirement removed only what was authorized and nothing that was not. Check the keep-list byte-for-byte against main yourself ... Also confirm no historical tree (_archive, _evidence, _reviews, _debriefs, docs/roadmap) was edited.
  Not reopened per the FO. Re-checked for the delta only: `git diff --name-only HEAD~1..HEAD` carries no `_archive|_evidence|_reviews|_debriefs|docs/roadmap|.spacedock-state` path, and the oracle, live test, and registry do not appear in the commit at all.
- DONE: Verify the surviving auto-team seam still works, not just that its tests pass ... Removing LegacyTeamNameAdvisory and MemberExists must not have taken a live caller with them — find every former caller and confirm each is genuinely dead, not merely untested.
  The delta's own version of this question — the `teamName` parameter removal — verified by callers, not greps: on main `EnumerateDeclaredStandingTeammates` had two (`build.go:693` legacy injection, `standing.go:62` sentinel); cycle 1 deleted the first, so no caller could pass `""` and the `teamName == ""` early return was already dead before this cycle removed it. `show-standing` exercised on both binaries: byte-identical stdout (1709 b, non-empty — the enumeration fires, so not a both-empty false pass) and stderr. `reconcile --team-name` byte-identical (2500 b, exit 0). Falsifying change run: making the enumerator return nil reds `TestShowStandingParity` (`standing_parity_test.go:82`, golden vs `""`).
- SKIPPED: Test the TestMain bug fix as a necessity claim, and hunt for what it might mask ... CRITICALLY that no golden now passes only because the variable is cleared.
  Not reopened. Cycle 1 proved the fix correct, minimal, and non-masking by mutation; cycle 2 accepted the "hermeticity gap this change introduced" recharacterization with no code change. The delta touches no golden.
- SKIPPED: Check the deferred pair and the 36 regenerated goldens ... a golden whose dispatch_file_path still carries a team-shaped prefix is a real finding.
  Not reopened. The delta regenerates no golden and touches no `internal/ensigncycle` file, so the deferred pair is unchanged; the full suite green (including the golden-parity tests) is the evidence no golden needed regeneration for this delta.

### Cycle-2 verification

1. **The literal mirrors rule 6, and the override direction matches.** Rule 6 (`build.go:535-539`) is `subagentType := "spacedock:ensign"` then override from `stageMeta.Agent()`; both template arms now carry that same structure — literal as base, stage `agent:` as override. Probed against the binary rather than read: no `agent:` emits `spacedock:ensign`; `agent: custom:worker` emits `custom:worker`; `agent:` present-but-empty is refused exit 1 at rule 7 pointing at `status --validate`, so the "names one" vs "key present" wording gap emits no dispatch at all. The break-glass fixture README (`dispatch_recovery_fixtures_test.go:12`) declares no `agent:`, so the literal is the correct value for the test that red. Gain the implementation did not claim: cycle 1's `{dispatch_agent_id}` is defined at `fo-dispatch-core.md:16` as defaulting to `ensign` — not `spacedock:ensign`, which is what `:93` says — so an FO resolving the slot from that step would have filled a value the oracle rejects. The literal removes the template's dependence on that divergent line. Both lines byte-identical to main.
2. **Deferred risk 5 is narrower on the axis this cycle touched, with a wider blast radius than cycle 1 recorded.** The trigger is unchanged: a live fixture declaring `agent:` still makes a compliant FO emit a value the oracle rejects. Narrower because the template's primary instruction now agrees with the oracle by construction — an FO that ignores the override clause emits the passing value, where failing the old lookup emitted nothing (the observed red) or `ensign`. Wider in blast radius, found this cycle but not caused by it: a second oracle hardcodes the literal, `mergedEnsignDispatches` (`merged_team_mode_test.go:63`), so the same trigger reds merged mode too. **Maintenance cost plainly: this cycle added zero copies.** Cycle 1's comment already carried the literal; cycle 2 moved it from the comment into the value. The string is already an unnamed de-facto constant — `build.go:536`, `stamp.go:94` (a verbatim duplicate of rule 6), `dispatch/reconcile.go:294`, `claudeteam/reconcile.go:75`, both oracles, the template's two arms, `fo-dispatch-core.md:93` and `:16` — and the last two already disagree. No drift guard exists; a string match over the template would be the check this workflow's proof policy forbids, so the only legitimate guard compares the template's literal against the binary's emitted default.
3. **The two polish fixes did not break the seam.** `--team` is refused at exit 2 on all four argv forms (`--team=v`, `--team v`, bare trailing, flag-first), stdout empty, stderr naming the flag; base `21b7609f7` exits 0 with the full 6865-byte array on all four, so the behavior change is real. Shape matches `build --team-name`: same argv scan, same stderr form, same exit 2. Not over-broad — `--team-roster` and `--teamname` still exit 0. No live caller passed `--team`: the shipped contract (`claude-fo-dispatch.md:49`) instructs it without `--team` and is byte-identical to main, so there is no old-skill/new-binary skew. Falsifying change run: deleting the tombstone reds `TestSpawnStandingAllTeamFlagRefused` (`exit=0, want 2`, both forms). Show-standing injection unchanged from cycle 1's differential — a merged build against `docs/dev` (which has a standing mod) is byte-identical base↔HEAD apart from the launcher path, neither emits show-standing, both carry `run_in_background: true` and no `team_name`.
4. **Nothing else rode along.** 5 files, +39/-12 — exact. Cumulative `f64c733ef...HEAD` is 78 files, +541/-1265, net -724, inside the declared [-250,-750] band. The 78th file is `internal/dispatch/mods.go`, which cycle 1 did not touch — not the test file the cycle-2 report's parenthetical points at. `go test ./...`: one failure, `TestCodexResolveManifestAgainstInstalledHost`, the environment-local `~/.codex` probe, which also reproduces in a clean throwaway checkout at HEAD. `-race` on dispatch+claudeteam clean, 0 races; `gofmt -l ./cmd ./internal` and `go vet ./internal/dispatch/` clean.

### Does the literal remove the choice point?

Yes, for the field that failed. `subagent_type` is a literal in the value: no token to resolve, no lookup, and byte-identical to what rule 6 actually emits for this fixture — verified by running the build, not by reading the template. That degree of freedom is gone. Prose remains load-bearing elsewhere in the template, but not for that field: `name="{worker_key}-{slug}-{stage}"` still carries an unbound `{worker_key}` whose definition lives in another file (`fo-dispatch-core.md:94`), and the team-mode oracle's `nameShaped` check reds `selected-team` for the same class of reason. It went 2/2 in the observed runs, so it is not the observed defect, but it is now the template's highest-risk remaining slot. The override clause is prose and stays prose; what changed is its failure mode — ignoring it now yields the correct value in the common case instead of a dropped parameter.

### Review-finding disposition

No material finding. AC-5's mechanism now matches both the captain's instruction and the binary's rule.

1. Polish — the new comment drops cycle 1's `(the build helper's default, build.go rule 6)` cross-reference, so the template no longer points a maintainer at the copy it must track. One-parenthetical fix; this is the only maintenance cost cycle 2 actually added.
2. Deferred risk (revises cycle-1 risk 5) — same trigger, now reddening two oracles rather than one. Promote when any live fixture declares `agent:`.
3. Deferred risk (pre-existing, out of scope) — the two oracles disagree on a dropped `subagent_type`: merged mode tolerates omission when the transport is merged-shaped, break-glass requires the literal. The exact red this task fixed would have passed merged mode. This argues for the fix chosen over loosening break-glass. Promote if a future red is triaged against the wrong oracle.
4. Deferred risk — `{worker_key}` is the template's remaining unbound slot; the analogous partial binding is `name="spacedock-ensign-{slug}-{stage}"`. Outside AC-5, which names the `{dispatch_agent_id}` slot only.
5. Polish (pre-existing, unchanged from main) — `fo-dispatch-core.md:16` says `dispatch_agent_id` defaults to `ensign` while `:93` says `spacedock:ensign`. Nothing shipped now reads the divergent line for break-glass, but it is a one-word fix.
6. Polish — `stamp.go:94` duplicates rule 6's four lines verbatim; a shared helper would make the production default single-source.
7. Polish — the cycle-2 report attributes the 77→78 file count to the test file's growth; the added file is `mods.go`.

Carried forward unchanged and already accepted: the 78-file surface against the ≤60 cap.

### Summary

Recommendation: **PASSED**. The literal does what the captain asked: `subagent_type` is now a value copy, structurally identical to rule 6, and I confirmed the correspondence by running the binary across absent, named, and empty `agent:` READMEs rather than reading the code — the empty case is refused upstream, so the comment's "names one" wording has no reachable gap. The fix is also stronger than the implementation claimed, because it removes the template's dependence on `{dispatch_agent_id}`, a token whose own definition contradicts itself between two lines of `fo-dispatch-core.md`. On the maintenance question the gate asked me to price: this cycle added no copy at all — cycle 1's comment already carried the literal — so the real cost is the dropped `build.go rule 6` cross-reference, and the underlying duplication is a pre-existing unnamed constant with copies in four production sites and two oracles. Both polish fixes hold: the `--team` refusal is exit 2 across all four argv forms and not over-broad, no live caller ever passed it, and the parameter removal was provably dead code after cycle 1 — with `show-standing` byte-identical on a non-empty output and mutation-covered by `TestShowStandingParity`. Delta is exactly 5 files/+39/-12 with no historical tree touched. My one honest reservation is recorded as a deferred risk rather than a rejection: `{worker_key}` is now the template's weakest slot for the same reason `subagent_type` was, and it is outside AC-5's scope.
