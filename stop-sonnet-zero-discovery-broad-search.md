---
title: Stop Sonnet zero-discovery broad search
status: validation
score: "0.90"
source: "PR #663 Sonnet zero-discovery failure, 2026-08-10"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: 3rns0vh3svq49w43cfr0wdqd
gates:
    version: 1
    records:
        - id: gate:3rns0vh3svq49w43cfr0wdqd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3rns0vh3svq49w43cfr0wdqd-backlog-1
              briefing:
                id: briefing:3rns0vh3svq49w43cfr0wdqd:backlog:attempt-1:revision-1
                digest: sha256:065ca5a81f4f0358e6bf789fba0888899021fb8a1236a79d56a9cb8504ae8ab2
                request-digest: sha256:6478d6c2a4d6f48380b4bbed9541397c34fef4b40ff9afabb286da733bf8f069
                room-ref: ./stop-sonnet-zero-discovery-broad-search/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3rns0vh3svq49w43cfr0wdqd:backlog:1
                briefing: briefing:3rns0vh3svq49w43cfr0wdqd:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T18:38:26.692525Z"
                decision: approve
                reason: Captain directed immediate ideation for the exact Sonnet zero-discovery product repair.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:3rns0vh3svq49w43cfr0wdqd:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:3rns0vh3svq49w43cfr0wdqd-ideation-1
              briefing:
                id: briefing:3rns0vh3svq49w43cfr0wdqd:ideation:attempt-1:revision-1
                digest: sha256:9a4d63cd96d5b7951a57d9c433ba2978e7a636f927493ebc57678d78e4067ae0
                request-digest: sha256:4676d3306a8138d0c187cfe03be796f30b2e44c081be21b9793ae0bd1db0ae4e
                room-ref: ./stop-sonnet-zero-discovery-broad-search/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3rns0vh3svq49w43cfr0wdqd:ideation:1
                briefing: briefing:3rns0vh3svq49w43cfr0wdqd:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T18:42:37.269787Z"
                decision: approve
                reason: Captain directed this Sonnet product repair. The exact two-file design preserves all excluded surfaces and has a 14-line hard limit.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-10T18:38:46Z
worktree: .worktrees/spacedock-ensign-stop-sonnet-zero-discovery-broad-search
---
## Problem

The exact Sonnet zero-discovery journey still broad-searches the filesystem at boot. In PR #663 run 31415550991, job 93543929461, artifact 9074747236, the First Officer ran `find` across the installed First Officer references. The journey failed in 24.73 seconds.

## Value

A Sonnet CLI user can start Spacedock with no workflow and receive the declared local identification result. The First Officer does not search the project, plugin, skill, or reference filesystem.

## Scope

- Repair only the Sonnet zero-discovery boot behavior.
- Keep the change in the smallest product instruction or runtime surface plus focused proof.
- Do not change n28, Pi, Codex, shared XFAIL policy, or unrelated live journeys.
- Do not add a permanent XFAIL.
- Use local Sonnet subscription authentication before required PR CI.

## Acceptance criteria

- AC-1: The exact Sonnet `TestLiveCommonZeroDiscovery` target passes normally and retains artifacts.
- AC-2: Boot uses only the declared local identification path. No `find`, recursive `grep`, or equivalent filesystem or reference sweep occurs.
- AC-3: A focused negative control fails for the exact broad-search command from artifact 9074747236.
- AC-4: Existing full, race, format, registry, and active-owner checks pass.
- AC-5: Required exact PR lanes pass before merge. Pi remains skipped.

## Baseline evidence

- Released user and workflow: Sonnet CLI zero-discovery boot.
- Observable harm: the First Officer broad-searches reference files instead of stopping as declared.
- Value authority: the zero-discovery live journey and Captain direction for truthful Sonnet evidence.
- Trigger: run 31415550991, job 93543929461, artifact 9074747236, command `find /tmp/spacedock-live-plugin-3439009114/skills/first-officer/references -iname "*claude*"`.

## Ideation requirements

- Name exact files and gross/net estimate before product edits.
- Identify the smallest behavior boundary and one falsifying control.
- Keep the normal local Sonnet failing baseline, repaired normal PASS, validation, PR, and merge flow.

## Root cause

Artifact 9074747236 records the complete failing Sonnet boot. The skill loader supplied this exact base:

`/tmp/spacedock-live-plugin-3439009114/skills/first-officer`

Sonnet then read `first-officer-shared-core.md` by its exact path. Before boot, Sonnet ran this command to locate the Claude adapter:

`find /tmp/spacedock-live-plugin-3439009114/skills/first-officer/references -iname "*claude*"`

The command failed the journey before `status --boot --identify --json` returned the terminal no-workflow result. The later boot action was correct.

The entry skill names the Claude adapter with a relative path. The same entry already forbids filesystem search and retains `{first_officer_base}`. The remaining ambiguity is the relative adapter path.

## Proposed approach

Change only the Claude adapter line in `skills/first-officer/SKILL.md`.

```diff
- Claude Code (`CLAUDECODE` env var is set): read `references/claude-first-officer-runtime.md`
+ Claude Code (`CLAUDECODE` env var is set): read exactly `{first_officer_base}/references/claude-first-officer-runtime.md`; do not list or search the references directory.
```

Add the artifact command as a table case in `TestDetectBroadSearchAtBoot`. This falsifying control must report a `find` broad-search error.

The existing live journey remains the behavioral proof. A static instruction check is not sufficient for the value claim.

### Alternatives

- More zero-discovery text in the shared core is insufficient. Sonnet searched before it invoked boot, and the shared core already forbids search.
- A binary error change is insufficient. The binary already returned the required terminal result after the prohibited search.
- A permanent XFAIL is prohibited. It would hide the released Sonnet behavior instead of repairing it.

## Spike record

The retained artifact supplies the required mechanism spike. It proves these facts:

- The loader supplied the exact absolute skill base.
- Sonnet used that base for a direct shared-core `Read` action.
- The named Claude adapter existed under the same base.
- Sonnet used discovery for the relative adapter instruction.
- After the search, Sonnet read the adapter directly and completed the zero-workflow stop correctly.

No additional spike is necessary. The artifact ZIP remains at `/tmp/xp6-pr663-sonnet.Qu0eOU/artifact.zip` with SHA-256 `def7637b2ebbe97006c1d0c10e13e95d04e4d35a76453da9320d07f27bb8826e`.

## Expected surface and semantic boundaries

The exact expected surface is two files:

- `skills/first-officer/SKILL.md`: +1/-1 lines.
- `internal/ensigncycle/broad_search_detect_test.go`: +8/-0 lines.

The estimate is +9/-1 lines, 10 gross lines, and +8 net lines. The hard tolerance is two files and 14 gross lines.

The allowed runtime change is one direct Claude adapter read instead of adapter-directory discovery. The final no-workflow output remains unchanged.

Command grammar, stored formats, authority, Codex behavior, Pi behavior, XFAIL policy, and unrelated journeys must not change. No documentation file changes.

The live registry already states the intended behavior: startup stops without broad filesystem discovery. This task restores that documented behavior.

## Acceptance criteria and proof

- **AC-1 (VALUE):** The exact local Sonnet zero-discovery journey has zero broad-search commands and passes normally. Compare artifact 9074747236, with one prohibited command and a 24.73-second failure, to the repaired retained artifacts.
- **AC-2:** Sonnet reads the Claude adapter at the exact retained-base path. The repaired stream contains the direct adapter `Read` action and no adapter-directory search.
- **AC-3:** The exact artifact command remains a falsifying control. The focused detector case feeds that command to `detectBroadSearchAtBoot` and requires a `find` error.
- **AC-4:** Full, race, format, registry, and active-owner checks pass. Existing test commands supply this evidence.
- **AC-5:** The exact required Sonnet and Codex PR lanes pass for the candidate. Pi remains skipped.

## Test plan

1. Run the focused detector case for the exact artifact command. It must pass because the detector rejects the command.
2. Run the exact local Sonnet target before the edit. Retain the normal failing baseline.
3. Apply the one-line Claude instruction edit.
4. Run the exact local Sonnet target. It must pass normally and retain the stream, final message, and command evidence.
5. Run focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
6. Run the registry and active-owner checks.
7. Run only the required exact PR lanes. Sonnet and Codex must pass, and Pi remains skipped.

The focused detector case has low cost. The local Sonnet target has model cost. Full, race, and PR checks retain the existing project cost.

## Stage Report: ideation

- DONE: Read the complete task and exact baseline evidence.
  Artifact 9074747236 shows the exact `find` command, the 24.73-second failure, the direct shared-core read, and the later correct stop.
- DONE: Inspect the zero-discovery boot contract and current Sonnet instruction path.
  The shared core and binary already stop correctly. The entry skill leaves the Claude adapter path relative.
- DONE: Define the smallest product mechanism and one focused falsifying control.
  Replace one Claude adapter line. Add the exact artifact command to the existing detector table.
- DONE: Name every proposed file and gross/net line estimate before product edits.
  The design changes two files, +9/-1 lines, 10 gross lines, and +8 net lines. The hard tolerance is 14 gross lines.
- DONE: Preserve n28, Pi, Codex, XFAIL policy, and unrelated journeys.
  The semantic boundary permits only the direct Claude adapter read. All named surfaces remain unchanged.
- DONE: Write and push a Simplified-English ideation Stage Report with recommendation.
  This report recommends approval of the two-file design. The state commit and push make the report durable.

### Recommendation

APPROVE ideation. Implement the one-line Claude adapter path and the exact detector control within the 14-line hard tolerance.

### Summary

Sonnet searched for a file whose exact path was already available. The smallest repair gives Sonnet that exact retained-base path at the selection point.

The existing live journey proves the outcome. The new focused control preserves the exact failing command from artifact 9074747236.

## Stage Report: implementation

- DONE: Read the complete task and approved ideation report.
  The approved scope permits two files and 14 gross lines.
- DONE: Confirm the exact baseline artifact and run the focused detector proof before the product edit.
  Artifact 9074747236 retains the failing command and has the recorded SHA-256 digest.
  The full detector passed on base `9c28d633a`.
  Two pre-edit Sonnet runs also passed without a search.
  The retained pre-edit run used a direct adapter read despite the relative instruction.
- DONE: Change only the two approved files within 14 gross lines.
  Candidate `289ede00015c7fa29961714f4c727cac77423524` changes two files and 10 gross lines.
  The diff is +9/-1 lines and +8 net lines.
- DONE: Run the exact local Sonnet zero-discovery target and retain normal PASS artifacts.
  The repaired target passed in 26.87 seconds with local subscription authentication.
  The stream contains the exact adapter `Read` action and no broad-search tool call.
  The evidence remains under `.spacedock-evidence/3rns/post-edit-sonnet/` in the code worktree.
- DONE: Run focused, format, full, race, registry, and active-owner checks.
  The focused detector, `gofmt`, full suite, race suite, and registry reconciliation passed.
  All five remaining XFAIL owner IDs resolve to active workflow records.
- DONE: Keep n28, Pi, Codex, XFAIL policy, and unrelated journeys unchanged.
  The candidate changes only the Claude adapter instruction and the exact detector table case.
- DONE: Commit and push the exact candidate and a Simplified-English implementation Stage Report.
  The candidate branch contains commit `289ede00015c7fa29961714f4c727cac77423524` on `origin`.
  This report records the implementation evidence for independent validation.

### Acceptance evidence

- **AC-1:** The repaired Sonnet target passed normally and emitted zero broad-search commands.
- **AC-2:** The retained stream reads `{first_officer_base}/references/claude-first-officer-runtime.md` directly.
- **AC-3:** The exact artifact command passes as a focused detector control and reports `find` as the signature.
- **AC-4:** The full, race, format, registry, and active-owner checks passed.
- **AC-5:** Required pull-request lanes remain for the First Officer after validation.

### Summary

The entry skill now gives Sonnet the exact Claude adapter path. The focused control retains the prohibited command from the failing artifact.

The local Sonnet target passes and uses the direct adapter read. The candidate is ready for independent validation.
