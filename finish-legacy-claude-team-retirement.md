---
id: nr75fq7ha3nmvsegbd22cgqa
title: Finish the legacy Claude TeamCreate retirement — the contract retired it, the binary and a live proof did not
status: implementation
source: "Captain CL, 2026-08-18, in chat ('file a proper legacy claude team retirement task'), after TestLiveBreakGlassShimRecovery/selected-team failed the claude-live lane on the ca9/6ht/j7j stack by demanding a team_name the shipped contract tells the FO to omit."
started: 2026-08-18T23:15:53Z
completed:
verdict:
score:
worktree:
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
