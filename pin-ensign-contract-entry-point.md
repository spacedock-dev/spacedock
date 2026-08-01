---
title: "Dispatched workers load the ensign contract through a resolving entry point, never the first-officer core"
status: validation
priority: 2
sprint-readiness: ready
source: "Captain, 2026-08-01, after diagnosing the pi ensign misload: every pi-spawned ensign this session (8 workers, both Kimi and gpt-5.6-luna) booted on the first-officer shared core — sometimes from stale .claude/.gemini plugin caches — because ~/.pi/agent/agents/ensign.md declares skills: ['spacedock:ensign'] and the preload silently fails, leaving the model to file-search for its contract. Root-cause question is OPEN: pi-subagents may not route agent-def preloads through pi's package resolver at all (the session was not started by assigning an agent, which motivates verifying rather than inferring)."
id: mxaaqb96syv7pq7ekg5a5194
gates:
    version: 1
    current:
        gate: gate:mxaaqb96syv7pq7ekg5a5194:validation
    records:
        - id: gate:mxaaqb96syv7pq7ekg5a5194:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mxaaqb96syv7pq7ekg5a5194-backlog-1
              briefing:
                id: briefing:mxaaqb96syv7pq7ekg5a5194:backlog:attempt-1:revision-1
                digest: sha256:d861ca46643dc1cf9100a563ad8ae289a697d67c8232cbb68372d02def8c850a
                digest-domain: canonical-bytes
                request-digest: sha256:c0a9bc4b100816269503de91e29f50c196dafbcc80f4f603ade1e69529724a31
                room-ref: ./pin-ensign-contract-entry-point/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:mxaaqb96syv7pq7ekg5a5194:backlog:1
                briefing: briefing:mxaaqb96syv7pq7ekg5a5194:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T12:26:59.34092Z"
                decision: approve
                reason: 'Accepts the task with amended direction: generic agent + load the right skill at spawn (codex pattern); no shipped def, no .pi, no manifest agents key; acceptance = spawned child provably boots ensign contract, zero first-officer reads'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:mxaaqb96syv7pq7ekg5a5194:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:mxaaqb96syv7pq7ekg5a5194-ideation-1
              briefing:
                id: briefing:mxaaqb96syv7pq7ekg5a5194:ideation:attempt-1:revision-1
                digest: sha256:b29e25e8303e6175f9d23f4bb81eb9104d49a13d38671c36ae964f84aae263bf
                digest-domain: canonical-bytes
                request-digest: sha256:6835053f386bc35578c331282090258011558f9af8407520203c9d7aefe45138
                room-ref: ./pin-ensign-contract-entry-point/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:mxaaqb96syv7pq7ekg5a5194:ideation:1
                briefing: briefing:mxaaqb96syv7pq7ekg5a5194:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T12:42:14.985156Z"
                decision: approve
                reason: Accepts the build-owned pi ensign delivery design + child exemption as presented (verdict approve, chat)
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:mxaaqb96syv7pq7ekg5a5194:validation
          stage: validation
          attempts:
            - id: gate-attempt:mxaaqb96syv7pq7ekg5a5194-validation-1
              briefing:
                id: briefing:mxaaqb96syv7pq7ekg5a5194:validation:attempt-1:revision-1
                digest: sha256:d7c78539e27fca663b55ff78be5877635d106392b725e300a20faa255fd6da6d
                digest-domain: canonical-bytes
                request-digest: sha256:7bc284b9c5c5c55bd5d6135550148545f08c3c42a4b632ca5d90ea1be0ff059e
                room-ref: ./pin-ensign-contract-entry-point/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:mxaaqb96syv7pq7ekg5a5194:validation:1
                briefing: briefing:mxaaqb96syv7pq7ekg5a5194:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-01T13:37:28.13926Z"
                decision: revise
                reason: 'Revise into one narrow correction round (captain: proceed): apply the three validator-proven test-only patches — argv ''dispatch'' for the build exec; stdout/stderr separation with stdout-only JSON parse; seed piHome/settings.json with the repo as path package — then re-run the tagged AC-1 live leg on the candidate worktree to produce the gate-grade boot artifact'
              application:
                action: feedback
                target-stage: implementation
                state: pending
started: 2026-08-01T12:27:09Z
worktree: .worktrees/spacedock-ensign-pin-ensign-contract-entry-point
---

## Problem

The pi ensign agent definition's contract preload doesn't bind. Workers fall
back to their own filesystem search and land on `first-officer` material, so
the dispatcher's authority leaks into every worker boot. Strikes observed:
8/8 in one session across two models.

## End value

A `.pi`-dispatched worker provably boots on the ensign shared core (the exact
file, from the package) and demonstrably does NOT load anything named
first-officer — verified by reading its session log's first tool calls.

## Investigation directions for ideation (captain-seeded)

1. Search pi's issue tracker/discussions for `agent:` semantics in
   pi-subagents: do agent-definition `skills:` preloads resolve package
   namespaces (`spacedock:ensign`, package pi manifest), and does the loader
   differ for user-scope vs package-scope agent definitions?
2. Read the pi-subagents agent-loading code path for the actual resolution
   semantics instead of inferring them.
3. Worst case accepted by design: drop the custom agent def; dispatch pi's
   general-purpose agent and have the dispatch assignment itself load the
   ensign contract through a resolvable skill name (invoked by name, the way
   `$spacedock:first-officer` works in sessions).
4. The stale `.claude/plugins/cache/**` first-officer copies are part of the
   trap's surface; decide whether cleanup is in scope.

---

### Feedback Cycles

- Cycle 1 (2026-08-01, captain ruling in-session): Direction reset before first stage work — use the generic agent and load the right skill at spawn time ("same as codex"); no shipped agent-definition file, no .pi/agents path (.pi is local), no manifest agents registration. Acceptance stands: a spawned worker provably boots the ensign contract. FO diagnosis on record: basename-only skill resolution (probe-verified), the FO-bootstrap extension injecting into every child session (the active leak vector), and the current pi dispatch-build shape (file-pointer prompt, self-contained assignment, env-level SPACEDOCK_BIN wrapper, subagent_type spacedock:ensign, model=null with opus unsettable on pi).

## Ideation

### Root cause (verified this stage, not inferred)

1. **Basename-only pi-subagents skill resolution.** pi-subagents resolves every agent-def
   and spawn-param skill by directory basename: `resolveSkillPath` matches
   `s.name === skillName` where `name` is the `SKILL.md` parent dir
   (`src/agents/skills.ts`, pi-subagents 0.35.1 installed at
   `~/.pi/agent/npm/node_modules/pi-subagents`). The string `spacedock:ensign`
   names nothing, so `~/.pi/agent/agents/ensign.md`'s
   `skills: ["spacedock:ensign"]` preload silently resolves to zero skills and the
   worker boots contract-free. Claude resolves the same namespaced string in
   plugin agent frontmatter — the two hosts genuinely differ, which is why
   "just fix the declaration" is not a portable fix.
2. **The extension is the second, independent leak vector.** pi-subagents spawns
   children as full pi processes (`--no-extensions` only when the caller
   overrides; default inherits project+user extension discovery), so repo
   `.pi/extensions/spacedock.ts` loads in every child, `session_start` arms
   `injectBootstrap`, and the `context` hook inserts the FO bootstrap child-side.
   Live-confirmed: this worker's transcript (runId
   `146c60be-ccf8-4203-bfd4-50d67105d571`) contains the SPACEDOCK-FO-BOOTSTRAP
   text. Compliance with the ensign contract here was model courtesy, not
   structure — a child that obeys the injected claim reads
   `skills/first-officer/SKILL.md` and the strike is back.
3. **The receiving shape already works.** That same run was spawned as generic
   agent `worker` with pi-subagents `skill: "ensign"` delivery: the injected
   `available_skills` block carried only `ensign` with its absolute package
   location, and the transcript's first tool calls are read(dispatch file),
   read(skills/ensign/SKILL.md), read(entity), read(ensign-shared-core.md),
   read(pi-ensign-runtime.md) — zero first-officer file reads. That is exactly
   the captain's acceptance shape, achieved today only by FO hand-composition,
   which the adapter forbids ("must not be replaced by a locally composed
   assignment"). The fix is to make this build-owned and remove vector 2 so the
   outcome is deterministic.

### Proposed approach (implements captain ruling 2026-08-01: generic agent + load the right skill at spawn; no shipped def, no .pi/agents, no manifest agents key)

1. **`dispatch build --host pi` emits two pi-subagents spawn fields** in the
   JSON artifact: `agent: "worker"` (the pi-subagents generic write-capable
   agent) and `skill: "ensign"` (basename — the only form pi's resolver binds;
   the package path comes from the spacedock user-package `pi.skills` entry).
   `subagent_type` stays `spacedock:ensign` as the host-neutral role identity
   and worker-name/worktree key, so pi worker names, branches, and worktree
   conventions do not churn. A stage `agent:` override replaces `agent` and
   omits `skill` (an override agent owns its own contract). Claude and codex
   host output is byte-identical to today.
2. **Pi FO adapter + wrapper map the artifact fields.**
   `internal/piruntime.SubagentDispatch` gains `Agent`/`Skill`;
   `pi-first-officer-runtime.md` `«worker.spawn»` is updated to pass
   `{artifact.agent}` as the subagent agent and `{artifact.skill}` as the spawn
   `skill` param (alongside existing `context: "fresh"` and `cwd`). Exact
   before/after wording proposed at implementation time; the contractlint pin
   (test plan) fixes the two load-bearing sentences.
3. **Extension FO-bootstrap exemption — RECOMMENDED: exempt.**
   `.pi/extensions/spacedock.ts` skips FO-bootstrap injection when
   `process.env.PI_SUBAGENT_CHILD === "1"` (exported by pi-subagents
   `src/runs/shared/pi-args.ts:18`, set on every spawned child since at least
   0.35.1; presence runtime-verified in this child session). Root FO sessions
   never carry the var; children are delegated workers by definition and should
   never be commissioned as first officer. Edge: an FO run as a pi-subagents
   child would lose its bootstrap — accepted; FO commissioning is a
   root-session concern, and this is the pi-subagents installer contract.
   Alternative (skip the exemption, rely on skill-param-only context): rejected
   as nondeterministic per root cause 2.
4. **Claude and `agents/` untouched.** No changes under `agents/`,
   `.claude-plugin/`, or the claude/codex host branches of build.go —
   `agents/ensign.md`'s namespaced declaration is correct for Claude, and the
   pi fix does not need it.
5. **Stale-cache decision (investigation direction 4): out of shipping scope.**
   `~/.pi/agent/agents/ensign.md` and stale `.claude/plugins/cache/**` /
   `.gemini` first-officer copies are operator-local trap surface, not repo
   state — no hidden machine dependencies may ship. After approach 1–3, no pi
   dispatch path consults agent-def preloads or filesystem search at all, so
   the caches are unreachable on the normal path. Operator-side deletion of the
   local `~/.pi/agent/agents/ensign.md` is recommended hygiene and is recorded
   in the stage report, but no acceptance depends on it.

### Acceptance criteria (entity-level; the gate approves these as the baseline)

- **AC-1 (the value measure, captain-pinned):** A worker spawned through
  `spacedock dispatch build --host pi` + the pi-subagents spawn provably boots
  the ensign contract: the child's session transcript
  (`.pi-subagents/artifacts/*_transcript.jsonl`) has
  `<package>/skills/ensign/SKILL.md` among its first five read-type tool calls
  and contains zero reads of any path matching `first-officer`. Baseline before
  this task: 0/8 observed pi ensigns (all booted FO material). Pass bar: the
  live acceptance run leaves a graded transcript artifact meeting both
  conditions; the count is reported in the stage report.
  Test: transcript grader wired into the existing tagged live pi harness
  (`internal/ensigncycle` pi lane), extended from
  `TestLivePiSubagentEnsignSmoke`.
- **AC-2:** For a default (no `agent:` override) stage,
  `spacedock dispatch build --host pi` JSON carries `agent`=`"worker"` and
  `skill`=`"ensign"`; an `agent:` override replaces `agent` and causes `skill`
  to be omitted; `--host claude` build output is byte-identical to its
  pre-change golden. Test: Go fixture tests in `internal/dispatch`,
  including one byte-compare guard against the current claude golden.
- **AC-3:** With `PI_SUBAGENT_CHILD=1` in the environment, the Spacedock pi
  extension's `context` hook injects zero FO bootstrap across
  session_start/session_compact; without the var, every assertion of the
  existing behavior harness still passes unchanged (injection, de-dup,
  compaction placement, agent_end suppression). Test: the existing node
  behavior harness in `internal/piruntime/spacedock_extension_test.go`, run
  twice (env set / unset).
- **AC-4:** `skills/first-officer/references/pi-first-officer-runtime.md`
  binds `«worker.spawn»` to the artifact's `agent` and `skill` fields, and no
  pi-subagents agent named `spacedock:ensign` appears anywhere in the pi
  dispatch path (adapter text, wrapper struct, build pi branch).
  Test: contractlint pin asserting both sentences' presence and the banned
  string's absence in the pi path.

### Expected surface (baseline for later calibration)

Files (9) and estimated insertions/deletions, tolerance ±50% on insertion
count; any file outside this list or a semantic outside the declarations below
is a boundary breach:

- `internal/dispatch/build.go`: +20/-3 — pi-branch `Agent`/`Skill` emission;
  `buildOutput` gains two `json:",omitempty"` fields; `«dispatch.build»` schema
  doc comment updated.
- `internal/dispatch/build_pi_host_test.go`: +45 — AC-2 assertions including
  override omission and the claude byte-identical guard.
- `internal/piruntime/subagents.go`: +10 — `Agent`/`Skill` fields on
  `SubagentDispatch` + doc comment.
- `internal/piruntime/subagents_test.go`: +30 — wrapper round-trip/omitempty
  tests.
- `.pi/extensions/spacedock.ts`: +8/-2 — `PI_SUBAGENT_CHILD` exemption gate.
- `internal/piruntime/spacedock_extension_test.go`: +25 — harness run with the
  env var asserting zero injection; existing assertions re-run unset.
- `skills/first-officer/references/pi-first-officer-runtime.md`: ±8 lines —
  `«worker.spawn»` wording (before/after recorded by implementation).
- `internal/contractlint/` (adapter pin test): +20.
- `internal/ensigncycle/pi_live_runner_test.go` (+possible sibling helper):
  +60 — transcript grader (scan run artifacts, first-five-reads and
  zero-FO-read assertions).

Estimated total ≈ 190 insertions, 5 deletions.

**Observable semantics changed (declared):** (a) pi dispatch-build JSON output
gains two keys — additive output-shape change; claude/codex bytes unchanged;
(b) pi-subagents child sessions no longer receive FO bootstrap injection —
runtime behavior change. **Unchanged:** command grammar, stored/entity formats,
authority/write-scope rules, claude+codex dispatch bytes, `agents/` and
`.claude-plugin/` contents.

### Test plan

- Go fixture tests for AC-2 (default fields, override omission, claude
  byte-identity). Cost: low. No new fixtures beyond one golden reuse.
- Node behavior harness (existing pattern) for AC-3. Cost: low; runs under
  plain `go test`.
- Wrapper JSON tests (omitempty behavior) for the spawn fields.
- contractlint pin for AC-4 (prose assertions on the two load-bearing adapter
  sentences + banned-string absence in pi path). Skill-text change is covered
  by this smoke pin per repo skill-development rules.
- AC-1 is the live lane: extend the existing tagged harness with the
  transcript grader (reads `.pi-subagents/artifacts/`, asserts ordering and
  absence). Cost: one live dispatch per CI run; complexity low-medium. It runs
  only in the live CI lane; default `go test ./...` and `-race` are unaffected.

### Spike record

No spike needed — both load-bearing mechanisms were exercised live during this
ideation, in this session:

1. **Spawn-time basename skill delivery works.** This run
   (`146c60be-ccf8-4203-bfd4-50d67105d571`, agent `worker`, skill delivery
   `ensign`) injected `available_skills` with only the ensign skill's absolute
   package path; the transcript's first tool calls read the dispatch file then
   `skills/ensign/SKILL.md`, and contain zero `first-officer` reads — the
   captain's acceptance shape, before any code change.
2. **The exemption marker exists and binds.** pi-subagents 0.35.1
   `pi-args.ts:18` exports `SUBAGENT_CHILD_ENV = "PI_SUBAGENT_CHILD"` and sets
   it to `"1"` on every child spawn; runtime-verified in this child process
   (`PI_SUBAGENT_CHILD=1`).
3. The resolver mismatch itself is source-read proof (basename `===` match),
   corroborated by the entity's recorded 8/8 strikes.

The implementation's first tests are the AC-2 fixture and AC-3 harness
extensions derived directly from these observations. No remaining unverified
mechanism gates the design; the transcript grader is new code but reads a
JSONL schema this session's artifact demonstrates.

### Documentation impact

The behavior change is pi-host-internal. No docs-site page documents the pi
build JSON keys or extension injection behavior; `docs/runtime-support.md`'s
harness `--skill` flags are unaffected. No doc diff proposed; FO/ensign skill
instruction text change is limited to the `«worker.spawn»` binding line above
and ships with its contractlint pin.

- Cycle 2 (validation attempt-1, 2026-08-01): REJECTED — AC-1's committed tagged harness cannot produce the gate artifact: (1) build exec argv missing the `dispatch` subcommand (exit 2, unknown command: build); (2) `cmd.CombinedOutput()` merges the bare-mode stderr advisory into the JSON envelope parse; (3) hermetic piHome lacks settings.json package registration so the basename skill resolves to nothing (child meta `skills: []`, contract-free boot). All three findings Material/evidence-kind, owned by this task's diff. Fixes live-proven on a throwaway checkout (live leg passes in 124.7s; `pi-ensign-boot-grade.json` = ensign rank 1, zero first-officer reads). Routed to implementation: apply the three patches, re-run the tagged live leg on the candidate worktree, produce the gate-grade artifact. Deferred risk accepted: `piLiveEnv` lacks an ambient `PI_SUBAGENT_*` scrub (promote condition: lane run nested under pi-subagents).
- Cycle 3 (validation attempt-2, 2026-08-01): PASSED — validator recommendation on the captain-authorized correction candidate d9f38ca1e; surface 10 files / +643/-24 vs estimate 9 files / ~190/5 (338% insertions; deviation pre-adjudicated JUSTIFIED in attempt-1 — the delta is the authorized test-only correction +57/-5); AC unchanged: all four ACs freshly reproduced, AC-1 live leg re-graded rank 1 / zero FO reads / agent worker / skills [ensign].
## Stage Report: ideation

- DONE: Verify Claude-side plain-basename vs namespaced skill declaration semantics before touching agents/ (no claude regression).
  Verified: `agents/ensign.md` and `agents/first-officer.md` declare namespaced `skills: ["spacedock:..."]` (Claude plugin frontmatter resolves these; claude live lanes exercise them); pi-subagents `resolveSkillPath` matches basename only (`skills.ts`, `s.name === skillName`) — so the pi fix avoids `agents/` entirely (approach point 4) and pins claude output byte-identity (AC-2).
- DONE: Shape the pi generic-agent dispatch path: spawn-time skill delivery (subagent skill param / boot-read instruction), per captain ruling "same as codex, generic agent + load the right skill", plus the extension FO-bootstrap exemption decision with a recommendation.
  Approach points 1–3: build emits `agent: "worker"` + `skill: "ensign"` for `--host pi`; exemption decision = exempt (`PI_SUBAGENT_CHILD === "1"` gate), recommendation and edge case recorded; marker runtime-verified in this child session.
- DONE: Pin acceptance to the captain's: a spawned worker's session-log first tool calls show skills/ensign/SKILL.md read and zero first-officer reads.
  AC-1 states exactly this, graded over `.pi-subagents/artifacts/*_transcript.jsonl` with a live-run pass bar and the recorded 0/8 baseline.

### Summary

Root-caused both leak vectors with source and live evidence: pi-subagents' basename-only skill resolver makes the `agents/ensign.md` namespaced preload unresolvable, and the Spacedock pi extension injects the FO bootstrap into every pi-subagents child session. The design makes contract delivery build-owned (`dispatch build --host pi` emits `agent: worker` + `skill: ensign`, mapped by the adapter/wrapper) and exempts `PI_SUBAGENT_CHILD` sessions from FO injection; this very session's transcript already demonstrates the target boot shape (generic worker + ensign skill delivery, zero first-officer reads), recorded as the spike evidence. Stale local caches are declared operator hygiene, out of the shipping path.

## Stage Report: implementation

- DONE: AC-2/AC-4 surfaces: dispatch build --host pi emits agent="worker" + skill="ensign" (override omits skill, claude golden byte-identical); internal/piruntime wrapper maps Agent/Skill; contractlint pins the two «worker.spawn» binding sentences and bans the namespaced string in the pi path
  TestBuildPiHostEmitsSpawnAgentAndSkill (default fields + subagent_type unchanged), TestBuildPiHostAgentOverrideOmitsSkill (skill key absent), TestBuildClaudeHostGoldenByteIdentical (byte-compare reusing golden build-crossproduct-single+flat+nonworktree+bare), TestSubagentDispatchSpawnFieldsRoundTrip, contractlint TestPiAdapterWorkerSpawnBindsArtifactAgentAndSkill + TestPiDispatchPathBansNamespacedAgentName. Red base af42c77d5: `agent = <nil>, want "worker"` / `skill = <nil>, want "ensign"`; `wrapped.Agent undefined`; «worker.spawn» missing both binding sentences. Pi golden build-host-pi-model-ignored.txt regenerated — diff is exactly the two additive trailing keys `agent`/`skill`.
- DONE: AC-3: .pi/extensions/spacedock.ts exempts PI_SUBAGENT_CHILD=1 from FO bootstrap; extension harness green run twice (marker set and unset) with all existing assertions unchanged
  TestSpacedockPiExtensionBootstrapBehavior (existing assertions verbatim, env marked scrubbed of PI_SUBAGENT_CHILD) and new TestSpacedockPiExtensionChildExemption. Red base af42c77d5: `Error: PI_SUBAGENT_CHILD=1 session_start injects zero FO bootstrap` (pre-fix extension still injected under the marker).
- DONE: AC-1: transcript grader implemented in internal/ensigncycle pi lane (extends TestLivePiSubagentEnsignSmoke fixture helpers): asserts ensign SKILL.md within first five reads and zero first-officer reads over .pi-subagents artifacts
  Per FO ruling (AC-textual), the smoke now spawns through the real build path: runPiSmokeDispatchBuild execs `dispatch build --host pi` for the fixture entity; the FO prompt forwards artifact agent/skill/task verbatim; assertPiEnsignBootContract grades the child's *_meta.json (agent/skills/task match the artifact) and *_transcript.jsonl (ensign SKILL.md read rank <= 5, zero first-officer reads), writing a durable pi-ensign-boot-grade.json. go vet -tags live green; live execution deferred to the tagged live lane (needs pi CLI + auth + network) — boot-shape evidence today is the recorded spike transcript 146c60be (this host's own 1/1).

### Summary

Implementation of the ideation's approved approach landed as commit 694d21db7 on spacedock-ensign/pin-ensign-contract-entry-point: the dispatch-build pi branch emits the spawn binding (`agent: worker`, `skill: ensign`; override replaces agent, omits skill), piruntime.SubagentDispatch maps the fields, the Spacedock pi extension skips FO-bootstrap injection for PI_SUBAGENT_CHILD=1 sessions, and the pi FO adapter «worker.spawn» bullet binds the artifact's agent/skill to subagent(...) — each with its pinned test (fixture, node harness run twice, contractlint pin) and TDD red output quoted above. Claude goldens are byte-identical (explicit re-run guard against build-crossproduct-single+flat+nonworktree+bare); the only byte change outside pi code is the additive pi golden regeneration. Surface re-measure vs the ±50% estimate: 591+/24- actual vs ≈190/5 baseline — attribution per FO ruling: ~200 of that is the ruling-B build-path smoke + transcript grader (line budget relaxed by the ruling), the rest is assertion verbosity in the same 9 declared files; no file outside the declared surface was touched. Residual risk: the AC-1 live lane has not been executed in this worktree (auth/network-gated); grader parsing is vet-verified only until the live CI lane runs it. Operator hygiene (unchanged from ideation): locally deleting ~/.pi/agent/agents/ensign.md is recommended but no acceptance depends on it.

## Stage Report: validation

- DONE: Reproduce the pinned evidence per AC: rerun TestBuildPiHost*/TestBuildClaudeHostGoldenByteIdentical, TestSubagentDispatchSpawnFieldsRoundTrip, the contractlint pins, and the extension harness BOTH env states on the worktree; quote exits.
  All green at 694d21db7: dispatch pi-host suite (incl. TestBuildPiHostEmitsSpawnAgentAndSkill, TestBuildPiHostAgentOverrideOmitsSkill, TestBuildClaudeHostGoldenByteIdentical) EXIT=0; TestSubagentDispatchSpawnFieldsRoundTrip EXIT=0; contractlint TestPiAdapterWorkerSpawnBindsArtifactAgentAndSkill + TestPiDispatchPathBansNamespacedAgentName EXIT=0; TestSpacedockPiExtensionBootstrapBehavior + TestSpacedockPiExtensionChildExemption EXIT=0, also EXIT=0 with PI_SUBAGENT_CHILD=1 exported (harnessEnv hermeticity works). Full `go test ./...` and `-race` EXIT=0 (ambient runtime markers scrubbed); gofmt clean.
- DONE: Replay red-first on base af42c77d5 for at least one leg per AC-2/3/4 (apply tests to base or cite the quoted red) and report red->green per leg.
  Throwaway worktree at af42c77d5 with the new test files applied: AC-2 red `agent = <nil>, want "worker"` / `skill = <nil>, want "ensign"` (exit 1); AC-3 red `Error: PI_SUBAGENT_CHILD=1 session_start injects zero FO bootstrap` (exit 1, node harness); AC-4 red `«worker.spawn» binding missing sentence "Pass the artifact's `agent` field as the subagent `agent` parameter verbatim."` (exit 1). TestPiDispatchPathBansNamespacedAgentName passes on base (absence pin, noted). All legs green on the candidate branch (exit 0). TestSubagentDispatchSpawnFieldsRoundTrip red on base as build failure `back.Agent undefined`.
- DONE: Attempt the AC-1 tagged live leg on THIS host (pi CLI + auth present): go test -tags live -run TestLivePiSubagentEnsignSmoke ./internal/ensigncycle in the worktree; quote the boot-grade artifact (pi-ensign-boot-grade.json: ensign SKILL.md read rank <=5, zero first-officer reads) or the exact refusal.
  REJECTED as committed — three evidence defects, all in live-tagged test-only code, each reproduced live, each narrow: (1) pi_live_runner_test.go runPiSmokeDispatchBuild execs `spacedock build` — exit 2 `unknown command: build` (CLI surface is `dispatch build`, internal/cli/cli.go newDispatchCommand); (2) `cmd.CombinedOutput()` merges the bare-mode stderr advisory ("WARN: bare_mode dispatch with no recent TeamCreate evidence", build.go:347) into the stdout JSON envelope parse — `invalid character 'W' looking for beginning of value`; (3) the hermetic fixture seeds piHome with auth.json only, so pi-subagents' settings-package skill discovery (skills.ts collectSettingsPackageSkillPaths over agentDir/settings.json) has no package: verified live that the FO passed `skill: "ensign"` verbatim (session jsonl, subagent args agent=worker skill=ensign) yet the child meta recorded `skills: []` and the child booted contract-free. Throwaway investigation (694d21db7 + 3 patches: argv "dispatch", stdout/stderr separation, piHome/settings.json `{"packages":["file:"+repo]}`): PASS 124.7s, pi-ensign-boot-grade.json {"ensign_skill_read_rank":1,"first_officer_reads":0,"read_calls":5,"spawn_agent":"worker","spawn_skills":["ensign"],"verdict":"pass"} — the end-to-end product behavior (build artifact -> pi-subagents spawn -> ensign boot, zero FO reads) is live-proven.
- DONE: Necessity-audit the surface: 591+/24- over the declared ~190/5 profile across the same 9 files with the FD ruling attribution (~200 smoke/grader) - rule whether the assertion bulk is justified or collapsible, with one cited example either way.
  Ruling: JUSTIFIED, small collapsible margin. pi_live_runner_test.go 199+/19- matches the FO ruling's ~200 smoke/grader attribution exactly. Production/adapter bytes (~52/5) are within tolerance. The remaining test bulk (~340 vs ~125 declared) is per-test fixture scaffolding in package convention, not assertion verbosity: justified example — harnessEnv (spacedock_extension_test.go, ~25 lines) pins marker state for both harness runs and closes a real hermeticity hole (this very session carries PI_SUBAGENT_CHILD=1; without the scrub the existing bootstrap test would silently inherit the marker); collapsible example — the ~20-line fixture prelude in TestBuildPiHostEmitsSpawnAgentAndSkill (build_pi_host_test.go:434-453) repeats near-identically across the branch's three new tests; a shared helper would save ~40 lines but the package has no such helper convention.

Findings (review-finding disposition, recommendation authority only): three Material evidence defects on AC-1 as above (defect kind: evidence; release scope: material — the approval-gated pi-live CI lane runs exactly this command; AC-1 is the captain-pinned value AC). Narrow fixes demonstrated green on throwaway; recommend one FO-authorized fix cycle applying the three patches, then re-run the live leg. Deferred risk: piLiveEnv does not scrub ambient PI_SUBAGENT_* markers (this run needed manual scrubbing; CI lane env is clean — promote-to-material if the live lane ever runs nested under pi-subagents). Pre-existing (on base af42c77d5, unrelated to this diff): TestVersionAmbiguousMarkersExitZero fails when ambient PI_CODING_AGENT is set — internal/cli env-hermeticity gap, outside this task's surface.

Recommendation: REJECTED — AC-1's gate evidence cannot be produced by the committed harness (three material evidence defects, fixes proven); AC-2/AC-3/AC-4 pass with reproduced red->green evidence; product behavior live-proven correct.

### Summary

All pinned tests reproduce green on 694d21db7 and each AC-2/3/4 leg replays red on base af42c77d5 with the quoted failure strings. The AC-1 live leg uncovers three blocking harness defects (wrong argv, stderr contamination, missing settings-package skill surface), none of which touch product code; with three narrow throwaway patches the same live leg passes and grades a real spawned worker at ensign-SKILL.md-read rank 1 with zero first-officer reads, so the captain's value shape is proven and the correction is confined to live-tagged test code. The 591+/24- surface is ruled justified (ruling-B attribution + package-convention scaffolding) with only a ~40-line polish margin.

## Stage Report: implementation (cycle 2)

- DONE: Apply the three proven harness patches (argv dispatch; stdout/stderr split; seed piHome/settings.json path-package) — show the diff
  Commit d9f38ca1e on spacedock-ensign/pin-ensign-contract-entry-point (pushed to origin): internal/ensigncycle/pi_live_runner_test.go only, +57/-5 — `exec.Command(binary, "dispatch", "build", ...)`, `bytes.Buffer` stdout/stderr split with stdout-only `json.Unmarshal`, `writeFile(piHome/settings.json, {"packages":["file:"+repo]})`. Test-only; no candidate semantics touched.
- DONE: Run the tagged live leg on the candidate worktree on this host; quote pi-ensign-boot-grade.json fields as the gate-grade artifact
  `SPACEDOCK_BIN=/tmp/spacedock-candidate SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-pi-live-artifacts go test -tags live -run TestLivePiSubagentEnsignSmoke ./internal/ensigncycle` — PASS 135.4s on d9f38ca1e. Artifact: {"verdict":"pass","worker_transcripts_graded":1,"read_calls":5,"ensign_skill_read_rank":1,"first_officer_reads":0,"spawn_agent":"worker","spawn_skills":["ensign"]} — meets the AC-1 bar (rank <= 5, zero first-officer reads, agent=worker, skills=[ensign]).
- DONE: Full go test ./... and -race (markers scrubbed) green, gofmt clean; AC-2/3/4 pins re-verified untouched
  `go test ./...` EXIT=0 and `go test ./... -race` EXIT=0 (env scrubbed of PI_CODING_AGENT/PI_CODING_AGENT_DIR/CLAUDECODE/CODEX_THREAD_ID/SPACEDOCK_BIN); `gofmt -l ./cmd ./internal` empty. AC-2 pins TestBuildPiHostEmitsSpawnAgentAndSkill/TestBuildPiHostAgentOverrideOmitsSkill/TestBuildClaudeHostGoldenByteIdentical PASS; AC-3/AC-4 internal/piruntime and internal/contractlint suites PASS — the cycle diff touches only the live-tagged harness file.

### SKIPPED

- SKIPPED: none — every checklist item completed.

### Summary

Applied the validator-proven three-patch correction package verbatim (argv `dispatch build`, stdout/stderr separation with stdout-only envelope parse, piHome/settings.json path-package seed) plus the package's optional-adjacent ambient `PI_SUBAGENT_*` scrub in `piLiveEnv` — justified here because this run executed nested inside a pi-subagents session (the recorded promote-to-material condition; without the scrub the forked FO pi process would inherit `PI_SUBAGENT_CHILD=1` and exempt itself from its own bootstrap). One operational note for replay: this host's ambient `SPACEDOCK_BIN` points at the main checkout's stale binary, which lacks the pi agent/skill fields — the live leg must run with `SPACEDOCK_BIN` unset (harness falls back to `go build` of the worktree) or pointed at a freshly built candidate binary, mirroring the CI convention. The AC-1 gate artifact is now produced by the committed harness on the candidate branch: live-spawned worker boots the ensign contract at read rank 1 with zero first-officer reads.

## Stage Report: validation (cycle 2)

- DONE: Re-run the AC-1 tagged live leg on the candidate worktree AT d9f38ca1e on this host (pi + auth): reproduce pi-ensign-boot-grade.json and quote every field as the verdict evidence (ensign_skill_read_rank <= 5, first_officer_reads == 0, spawn_agent worker, spawn_skills [ensign]).
  `env -u PI_CODING_AGENT{,_DIR} -u CLAUDECODE -u CODEX_THREAD_ID -u PI_SUBAGENT_{CHILD,RUN_ID,DEPTH} -u SPACEDOCK_BIN SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-pi-live-artifacts-validate go test -tags live -run TestLivePiSubagentEnsignSmoke ./internal/ensigncycle` — PASS 172.5s, exit 0. Reproduced artifact `/tmp/spacedock-pi-live-artifacts-validate/pi-subagent-ensign-smoke/run/pi-ensign-boot-grade.json`: {"verdict":"pass","worker_transcripts_graded":1,"read_calls":5,"ensign_skill_read_rank":1,"first_officer_reads":0,"spawn_agent":"worker","spawn_skills":["ensign"],"transcript":".../d65c3c7e_worker_0_transcript.jsonl"} — every AC-1 bar met (rank 1 <= 5, zero first-officer reads, agent worker, skills [ensign]). Falsifiability: a stale binary (no agent/skill envelope keys), a contract-free child (no ensign read in first 5), or any FO bootstrap leak (first_officer_reads > 0) each fail this test.
- DONE: Verify the correction was test-only: diff 694d21db7..d9f38ca1e shows only internal/ensigncycle/pi_live_runner_test.go (+57/-5), the three named findings' lines now correct (argv "dispatch", stdout/stderr separated with stdout-only parse, piHome/settings.json seeded), and TestPiLiveEnvScrubsAmbientPiSubagentMarkers exists and passes.
  `git diff --numstat 694d21db7..d9f38ca1e` = `57 5 internal/ensigncycle/pi_live_runner_test.go` (single file). Full diff inspected: argv now `dispatch build` (matches internal/cli surface), stdout/stderr split via bytes.Buffer with stdout-only json.Unmarshal (stderr cannot contaminate the envelope), `piHome/settings.json` seeded `{"packages":["file:"+repo]}` (settings-package skill discovery now resolves basename `ensign`). TestPiLiveEnvScrubsAmbientPiSubagentMarkers exists (pi_live_runner_test.go:367, live-tagged file) and PASSES under `-tags live` alongside TestPiLiveEnvDropsForeignRuntimeMarkers; it fails if piLiveEnv stops dropping the PI_SUBAGENT_ family.
- DONE: Re-verify AC-2/3/4 pin suites + full go test ./... and `go test ./... -race` (runtime markers scrubbed) exit 0 and gofmt clean on the worktree.
  AC-2 pins PASS: TestBuildPiHostEmitsSpawnAgentAndSkill (default agent=worker/skill=ensign; fails if build stops emitting either field), TestBuildPiHostAgentOverrideOmitsSkill (fails if an override re-adds the skill key), TestBuildClaudeHostGoldenByteIdentical (fails on any claude-output byte drift). AC-3/AC-4: `go test ./internal/piruntime ./internal/contractlint` PASS (extension exemption + spawn-binding pins). Full `go test ./...` EXIT=0 (19 ok), `go test ./... -race` EXIT=0 (19 ok), zero FAIL lines, both with runtime markers and SPACEDOCK_BIN scrubbed. `gofmt -l ./cmd ./internal` empty.

### Semantic adversarial pass

Changed behavior is test-only, so the pass targets the evidence boundary: (1) identity — verified `piSpacedockBinary` with SPACEDOCK_BIN unset runs `go build ./cmd/spacedock` from the worktree at d9f38ca1e, so the graded run exercised candidate code, and the envelope assertion (agent=worker/skill=ensign) cannot pass against the stale ambient binary noted as a hazard in cycle 2; (2) authority/attribution — the grade's zero first-officer reads binds the PI_SUBAGENT_CHILD extension exemption end-to-end under a real pi + pi-subagents spawn, not a fixture; (3) cardinality — `worker_transcripts_graded:1` prevents a vacuous zero-transcript pass.

Findings: none new. Deferred risk carried from attempt-1, unchanged: pre-existing TestVersionAmbiguousMarkersExitZero fails under ambient PI_CODING_AGENT (on base af42c77d5, internal/cli env-hermeticity gap, outside this task's surface) — promote condition: any supported workflow that must run with ambient runtime markers set. Operational replay note (non-blocking): this host's ambient SPACEDOCK_BIN is a stale main-checkout binary; the live leg must run with SPACEDOCK_BIN unset (worktree build fallback) or pointed at a fresh candidate build.

Recommendation: PASSED — all four ACs carry fresh, independently reproduced evidence on d9f38ca1e; the captain-authorized correction is confined to live-tagged test code; no material findings.

### Summary

Independently re-ran every verdict leg on the candidate at d9f38ca1e: the tagged AC-1 live leg re-passed in 172.5s and reproduced pi-ensign-boot-grade.json with ensign SKILL.md at read rank 1, zero first-officer reads, agent=worker, skills=[ensign]; the cycle-2 diff is exactly one live-tagged test file (+57/-5) with all three attempt-1 findings corrected and the new ambient-PI_SUBAGENT_ scrub test passing; AC-2/3/4 pin suites, full go test ./... and -race (markers scrubbed), and gofmt are all green. The rejected attempt-1 gap — a committed harness unable to produce the gate artifact — is closed by evidence, and the harness now self-scrubs the nested pi-subagents marker family that required manual scrubbing before. No candidate changes were made by this validation pass.
