---
id: ntnywe6wfk1g5sersjbe5yt7
title: Verify the Pi spawn skill name resolves - bare "ensign" versus "spacedock:ensign"
status: ideation
source: "Pi/GLM FO field report 2026-08-26: dispatched ensigns repeatedly produced broken DONE formatting, and the FO's diagnosis was that the dispatch artifact passes skill \"ensign\" while the Pi skill loader wants the exact name spacedock:ensign — if true, Pi ensigns spawn without the ensign contract at all"
started: 2026-08-27T00:12:01Z
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:ntnywe6wfk1g5sersjbe5yt7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ntnywe6wfk1g5sersjbe5yt7-backlog-1
              briefing:
                id: briefing:ntnywe6wfk1g5sersjbe5yt7:backlog:attempt-1:revision-1
                digest: sha256:184ce70a29b675fee2ea8cd40983596f1c2d5c7a73796c1d9221d726a9e1fe9e
                request-digest: sha256:0bcc9236a024db1353e8c8ca9f770965634d5634082625db1503064403d6bcfc
                room-ref: ./pi-spawn-skill-name-resolution/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ntnywe6wfk1g5sersjbe5yt7:backlog:1
                briefing: briefing:ntnywe6wfk1g5sersjbe5yt7:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T00:11:28.952704Z"
                decision: approve
                reason: 'Captain approve: enter ideation to flesh out the approach and test plan'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:ntnywe6wfk1g5sersjbe5yt7:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ntnywe6wfk1g5sersjbe5yt7-ideation-1
              briefing:
                id: briefing:ntnywe6wfk1g5sersjbe5yt7:ideation:attempt-1:revision-1
                digest: sha256:c1214eb165d4ad512de5fb703c46795cfaa710bc26f386e2223c5f31bb9b46ce
                request-digest: sha256:5909800fe51bd5c4d73aa7a3b644187ae7387adabc2cc5207e4f6a374a9d6433
                room-ref: ./pi-spawn-skill-name-resolution/review/ideation/briefing-1
---

`dispatch build --host pi` emits the bare skill name `ensign` by design (internal/dispatch/build.go piSpawnSkill), on a documented assumption: "pi-subagents resolves agents and skills by directory basename only." A field report says the loader needs `spacedock:ensign`. If the assumption rotted, every Pi ensign runs without its contract, and the observed broken stage-report formatting is a symptom, not a separate defect.

## Problem

The discriminator is one transcript read, and this ideation performed it. The field report claims the Pi skill loader wants `spacedock:ensign` (qualified) where the dispatch artifact emits bare `ensign`. Two independent lines of evidence refute that claim:

1. **Live transcript (this session, run 4c6e9c44).** This very ensign session was dispatched via `dispatch build --host pi` with `skill: "ensign"`. The system prompt carries an `<available_skills>` block listing the ensign skill at `/Users/clkao/git/spacedock-research/spacedock-v1/skills/ensign/SKILL.md` — the skill resolved and is resident before stage work begins. The ensign contract content is present in the inherited context. AC-1's falsifying condition does not hold.

2. **Source-level spike of the loader.** Read `resolveSkillPath` in `src/agents/skills.ts` on two versions: installed 0.37.2 and latest npm 0.57.0. The function is byte-identical across both: `skills.find((s) => s.name === skillName)` — an exact basename match, no qualified/namespace handling. Skill names are derived from directory basename (`path.basename(dirPath)` for `SKILL.md`) or file basename; a colon in the name would never match. `normalizeSkillInput` only trims whitespace and handles arrays — there is no colon-splitting or namespace stripping that would turn `spacedock:ensign` into `ensign`. Passing `spacedock:ensign` as the skill name would fail to resolve and land in `missing`, producing a contract-free boot — the exact failure the field report attributes to the bare name, but actually caused by the qualified name.

The field report's diagnosis is inverted: the bare `ensign` is the only form that resolves; `spacedock:ensign` is the form that would fail. The report likely confused the `subagent_type` field (host-neutral `"spacedock:ensign"`) with the `skill` field (bare `"ensign"`) — two different fields on the dispatch artifact. The existing test `TestBuildPiHostEmitsSpawnAgentAndSkill` (build_pi_host_test.go:427) already asserts `subagent_type = "spacedock:ensign"` AND `skill = "ensign"` simultaneously, confirming both fields coexist.

The observed broken DONE formatting is a symptom of the stage-report protocol surface, not a skill-resolution defect. That belongs to embed-stage-report-protocol-in-dispatch (t4, in flight) — out of scope here.

## Proposed approach

**Verify-first posture (executed).** The riskiest unverified mechanism was the newer-loader resolver behavior. Spiked it before deciding:

- Installed a throwaway copy of pi-subagents@0.57.0 (latest npm) in `/tmp/pi-subagents-spike`. Read `resolveSkillPath` and `normalizeSkillInput` in `src/agents/skills.ts`. Both are byte-identical to the installed 0.37.2. No namespace/qualified-name handling exists in either version. CHANGELOG review across 0.35–0.57 confirms no skill-name-qualification change (line 69 mentions "namespaced extension bindings" — extension RPC contracts, not skill name resolution; line 964 mentions nested grouped skills but the derived name is still the leaf basename).

- Confirmed the live dispatch path: this session's system prompt carries the resolved ensign skill in `<available_skills>`, proving the bare name resolves end-to-end on the current loader.

**Outcome: verified-working, no code change to piSpawnSkill.** The bare `ensign` resolves on both 0.37.2 (installed) and 0.57.0 (latest npm). No loader-version-aware logic is warranted — the resolver behavior has not changed across the version range. The only production change is a comment refinement: pin the verified loader versions in the `piSpawnSkill` comment block (build.go:56–63) so the documented assumption carries the evidence of re-verification rather than the stale unverified claim.

No new mechanism is introduced. The simplest alternative — changing `piSpawnSkill` to `"spacedock:ensign"` — is actively harmful: it would break resolution on every version, producing exactly the contract-free boot the field report fears. It is insufficient because it solves a problem that does not exist and creates the one it claims to fix.

## Risk evidence

**Spike result (riskiest mechanism, exercised first).** The newer-loader resolver behavior was the only unverified mechanism. It is now verified: `resolveSkillPath` is byte-identical on installed 0.37.2 and latest 0.57.0, both using exact basename match. The live session transcript (4c6e9c44) confirms end-to-end resolution of the bare name. No residual unverified mechanism remains.

Backlog baseline: the field report plus the source comment's unverified assumption decided design should start. The repeated DONE-format failures give the symptom baseline — now attributed to the stage-report protocol surface (t4), not skill resolution.

## Out of scope

The stage-report grammar and its messaging (9x) and the protocol embed (t4).

## Expected surface and tolerance

**Verify-first outcome: verified-working, no code change.**

- Production: net +0 across 1 file (internal/dispatch/build.go). The only edit is a comment refinement — replacing the stale unverified assumption text with the pinned loader versions and re-verification date. Insertions ~3, deletions ~3, net ~0. Tolerance: ±3 lines.
- Proof: net +0. No new tests needed — the existing `TestBuildPiHostEmitsSpawnAgentAndSkill` (build_pi_host_test.go:427) already asserts `skill = "ensign"` (the basename form) and `subagent_type = "spacedock:ensign"` (the host-neutral form). The spike evidence is recorded in this entity body, not in a new test file.
- Observable semantics: none changed. The dispatch artifact's `skill` field remains `"ensign"`; no CLI grammar change; no stored-format change; no runtime behavior change.

This task may close as verified-working (no code change beyond the comment pin) because the bare name still resolves on current pi-subagents. No loader-version-aware piSpawnSkill is required.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A Pi-dispatched ensign demonstrably loads the ensign contract: a live or recorded Pi dispatch transcript shows the skill content resident before stage work begins.**
Verified by: a live Pi dispatch transcript (this session, run 4c6e9c44) whose system prompt carries an `<available_skills>` block with the ensign skill entry pointing at `skills/ensign/SKILL.md`. The falsifying condition — the skill name fails to resolve and the transcript shows no contract load — does not hold: the ensign skill content is resident in the inherited context. Corroborated by source-level spike: `resolveSkillPath` on installed pi-subagents 0.37.2 and latest npm 0.57.0 both use exact basename match (`s.name === skillName`), confirming the bare `"ensign"` is the resolving form. The existing Go test `TestBuildPiHostEmitsSpawnAgentAndSkill` (build_pi_host_test.go:427) asserts the artifact emits `skill = "ensign"`, locking the verified-correct value.

## Test plan

Since the verify-first outcome is verified-working with no behavior change, the test plan is minimal:

1. **Existing Go test (no change needed).** `TestBuildPiHostEmitsSpawnAgentAndSkill` (build_pi_host_test.go:427) already asserts `*out.Skill == "ensign"` and `*out.Agent == "worker"` for a default-path Pi dispatch, with the comment "basename, the only form pi's resolver binds." This test locks the verified-correct spawn constant. Cost: zero (already passing). Run: `go test ./internal/dispatch/ -run TestBuildPiHostEmitsSpawnAgentAndSkill`.

2. **Comment refinement (production change, if applied).** The only production edit is updating build.go:56–63 to pin the verified loader versions: replace "pi-subagents resolves agents and skills by directory basename only" with text noting re-verification on 0.37.2 (installed) and 0.57.0 (latest npm), both using exact basename match, dated 2026-08-27. No test exercises a comment, but `go build ./...` and `go test ./internal/dispatch/` confirm no accidental code change. Cost: low.

3. **Pi-live lane.** Per the path-to-lane rule, the pi-live lane is NOT triggered because the spawn constants do not change — `piSpawnSkill` remains `"ensign"`, no loader-version-aware branch is added. The live transcript evidence (run 4c6e9c44) already satisfies the lane's purpose for this task.

No new mechanism is introduced, so there is no new mechanism to name a value AC for. The comment refinement serves AC-1 by making the re-verification auditable in-source rather than leaving the stale unverified claim. The simplest alternative — no comment change at all — is insufficient because it leaves the field report's diagnosis unchallenged in-source, inviting a future contributor to "fix" the bare name to the qualified form and break resolution.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

### Summary

Spiked the riskiest unverified mechanism (newer-loader resolver behavior) before deciding. Read `resolveSkillPath` and `normalizeSkillInput` in pi-subagents `src/agents/skills.ts` on installed 0.37.2 and latest npm 0.57.0 — byte-identical, both use exact basename match (`s.name === skillName`), no namespace/qualified-name handling. Confirmed end-to-end with a live dispatch transcript (this session, run 4c6e9c44): the system prompt carries the resolved ensign skill in `<available_skills>`, proving the bare `"ensign"` resolves. The field report's diagnosis is inverted — the bare name is the resolving form; `spacedock:ensign` is the form that would fail (likely confused the `subagent_type` field with the `skill` field). Outcome: verified-working, no code change to `piSpawnSkill`. The only production change is a comment refinement pinning the verified loader versions in build.go. The broken DONE formatting belongs to the stage-report protocol embed (t4), out of scope here.

## Stage Report: ideation

- DONE: Four-section design completion (Problem/Proposed approach/AC/Test plan)
  Entity body carries all four sections; Problem refutes the field report via live transcript + source spike, Proposed approach records the verify-first outcome, AC-1 names the falsifying condition and the live transcript that falsifies it, Test plan cites the existing Go test locking `skill = "ensign"`.
- DONE: Verify-first spike of the newer-loader resolver (executed — 0.37.2 and 0.57.0 byte-identical basename match)
  Read `resolveSkillPath` and `normalizeSkillInput` in `src/agents/skills.ts` on installed 0.37.2 and latest npm 0.57.0 — byte-identical, both use exact basename match (`s.name === skillName`); no namespace/qualified-name handling. Live dispatch transcript (run 4c6e9c44) confirms end-to-end resolution of bare `"ensign"`.
- DONE: Expected surface declaration (net ~0 LOC, comment pin only)
  Expected surface section declares net +0 across 1 file (internal/dispatch/build.go), comment refinement only (pin verified loader versions 0.37.2/0.57.0, dated 2026-08-27), tolerance ±3 lines. No observable semantic change; no new test needed (existing `TestBuildPiHostEmitsSpawnAgentAndSkill` locks the value).

### Summary

Ideation design was already committed (9ff6e54ca); this report appends the missing stage-report section. All three ideation checklist items are complete: the four-section design is in the entity body, the verify-first spike of the newer-loader resolver was executed (0.37.2 and 0.57.0 byte-identical basename match, live transcript confirms end-to-end resolution), and the expected surface declares net ~0 LOC (comment pin only). Outcome: verified-working, no code change to `piSpawnSkill`.
