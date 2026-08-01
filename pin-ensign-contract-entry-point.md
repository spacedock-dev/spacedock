---
title: "Dispatched workers load the ensign contract through a resolving entry point, never the first-officer core"
status: ideation
priority: 2
sprint-readiness: ready
source: "Captain, 2026-08-01, after diagnosing the pi ensign misload: every pi-spawned ensign this session (8 workers, both Kimi and gpt-5.6-luna) booted on the first-officer shared core — sometimes from stale .claude/.gemini plugin caches — because ~/.pi/agent/agents/ensign.md declares skills: ['spacedock:ensign'] and the preload silently fails, leaving the model to file-search for its contract. Root-cause question is OPEN: pi-subagents may not route agent-def preloads through pi's package resolver at all (the session was not started by assigning an agent, which motivates verifying rather than inferring)."
id: mxaaqb96syv7pq7ekg5a5194
gates:
    version: 1
    current:
        gate: gate:mxaaqb96syv7pq7ekg5a5194:ideation
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
started: 2026-08-01T12:27:09Z
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

## Stage Report: ideation

- DONE: Verify Claude-side plain-basename vs namespaced skill declaration semantics before touching agents/ (no claude regression).
  Verified: `agents/ensign.md` and `agents/first-officer.md` declare namespaced `skills: ["spacedock:..."]` (Claude plugin frontmatter resolves these; claude live lanes exercise them); pi-subagents `resolveSkillPath` matches basename only (`skills.ts`, `s.name === skillName`) — so the pi fix avoids `agents/` entirely (approach point 4) and pins claude output byte-identity (AC-2).
- DONE: Shape the pi generic-agent dispatch path: spawn-time skill delivery (subagent skill param / boot-read instruction), per captain ruling "same as codex, generic agent + load the right skill", plus the extension FO-bootstrap exemption decision with a recommendation.
  Approach points 1–3: build emits `agent: "worker"` + `skill: "ensign"` for `--host pi`; exemption decision = exempt (`PI_SUBAGENT_CHILD === "1"` gate), recommendation and edge case recorded; marker runtime-verified in this child session.
- DONE: Pin acceptance to the captain's: a spawned worker's session-log first tool calls show skills/ensign/SKILL.md read and zero first-officer reads.
  AC-1 states exactly this, graded over `.pi-subagents/artifacts/*_transcript.jsonl` with a live-run pass bar and the recorded 0/8 baseline.

### Summary

Root-caused both leak vectors with source and live evidence: pi-subagents' basename-only skill resolver makes the `agents/ensign.md` namespaced preload unresolvable, and the Spacedock pi extension injects the FO bootstrap into every pi-subagents child session. The design makes contract delivery build-owned (`dispatch build --host pi` emits `agent: worker` + `skill: ensign`, mapped by the adapter/wrapper) and exempts `PI_SUBAGENT_CHILD` sessions from FO injection; this very session's transcript already demonstrates the target boot shape (generic worker + ensign skill delivery, zero first-officer reads), recorded as the spike evidence. Stale local caches are declared operator hygiene, out of the shipping path.
