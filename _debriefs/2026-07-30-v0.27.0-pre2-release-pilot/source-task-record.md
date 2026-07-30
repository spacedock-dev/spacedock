---
title: Cut and field-test the current gate feature as a prerelease
status: validation
score: 1.0
sprint: durable-decisions
source: "Captain direction on 2026-07-30 to release the current gate feature before sprint closure and use real installations to discover friction."
id: 0hympdejewzwkhe60ygqk15a
gates:
    version: 1
    current:
        gate: gate:0hympdejewzwkhe60ygqk15a:ideation
    records:
        - id: gate:0hympdejewzwkhe60ygqk15a:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0hympdejewzwkhe60ygqk15a-backlog-1
              briefing:
                id: briefing:0hympdejewzwkhe60ygqk15a:backlog:attempt-1:revision-1
                digest: sha256:8e001f16c030b80444488bc79a581b04da24ec4f9b9aee31bfa734a88525caa2
                digest-domain: canonical-bytes
                request-digest: sha256:519a9e44c1703eb2a751839e98755b9710a25e4b43781f37827045f8547ebcff
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0hympdejewzwkhe60ygqk15a:backlog:1
                briefing: briefing:0hympdejewzwkhe60ygqk15a:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T14:03:29.385661Z"
                decision: approve
                reason: Captain directed an instrumented prerelease before sprint closure; the bounded pilot preserves release integrity while deferring non-blocking friction.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:0hympdejewzwkhe60ygqk15a:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:0hympdejewzwkhe60ygqk15a-ideation-1
              briefing:
                id: briefing:0hympdejewzwkhe60ygqk15a:ideation:attempt-1:revision-1
                digest: sha256:9169db971542a870433966c1f345ebf1c543a02d5efbce7a6f59cc9b18d25ce4
                digest-domain: canonical-bytes
                request-digest: sha256:67bd56e1ed2ba2884d9a052e507475481d92cd3f11c6bd04952fcc181add922e
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0hympdejewzwkhe60ygqk15a:ideation:1
                briefing: briefing:0hympdejewzwkhe60ygqk15a:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T14:27:20.511401Z"
                decision: approve
                reason: Captain authorized v0.27.0-pre2; design proves the provider-free chat journey after gqs and keeps provider/follow-up work out of the cut.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
started: 2026-07-30T14:04:09Z
worktree: .worktrees/spacedock-ensign-gate-feature-prerelease-pilot
---

## End value

Publish `v0.27.0-pre2` from freshly fetched current `main`, then use the installed edge bundle for one real, provider-free gate journey. The pilot must show that an ordinary user and two successive First Officer sessions can prepare, understand, authorize, resume, and spend one decision without hand-authored authority, while preserving the known post-consume actionability gap as sprint-critical evidence.

## Problem

The gate commands and First Officer contract have strong fixture and live-lane coverage, but no user has yet installed one released bundle and driven the complete provider-free lifecycle. Waiting for sprint closure would hide install, presentation, cold-resume, and exactly-once friction until the design is harder to change.

## Boundary

This is an instrumented prerelease, not the durable-decisions sprint acceptance walk. It changes no gate grammar, stored format, authority rule, or runtime implementation. It ships the merged surface, runs one chat journey, and reports friction; it adds no release machinery and absorbs no follow-up fix.

Captain override on 2026-07-30 removed `gqs` (`dispatch-entered-stage-after-gate-consume`) as a pre2 prerequisite. Capture `BASE_SHA=$(git rev-parse origin/main)` only after a fresh fetch. The current post-consume actionability gap is an expected sprint-critical pilot finding to preserve and route, not a release blocker or a reason to widen this cut.

Captain authorized the release-time E2E waiver on 2026-07-30. Runtime Live E2E run `30543637878` was green at PR #575 head `a0391fdbdb957d703ee4f836a5dfd5e21cc70153`; it is waiver evidence only. It is neither an exact release-SHA run nor proof that its head is the post-`gqs` release base. The report must preserve that distinction.

## Proposed approach

### Cut and publish

1. Fetch `origin/main` and tags, stop if `v0.27.0-pre2` already exists, then capture and retain `BASE_SHA=$(git rev-parse origin/main)`.
2. Create a release worktree from `BASE_SHA`. Run `gofmt -w ./cmd ./internal`, require a clean diff, then run `go test ./...` and `go test ./... -race`.
3. Run `spacedock-release stamp-version 0.27.0-pre2` against `.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, and `skills/first-officer/references/first-officer-shared-core.md`. Commit only those three files. Define `REL_SHA` as that commit; require its sole parent to equal `BASE_SHA` and its diff to contain only the three version substitutions.
4. Run `manifest-tag-gate v0.27.0-pre2` locally. Push the release branch to `main`, then generate and review public notes with `spacedock-release notes 0.27.0-pre2`. The resulting annotated tag must point to `REL_SHA` and carry a nonempty body.
5. Set `SPACEDOCK_E2E_GATE_WAIVER` to a reason naming the captain authorization, run `30543637878`, its green head, `BASE_SHA`, and `REL_SHA`, and explicitly state that the run is waiver evidence rather than release-SHA proof. Push only `v0.27.0-pre2`. Capture the release run and its waiver step summary. Delete the repository variable after the release run reads it, including on release failure, and verify that it is absent.
6. Require the existing release workflow to succeed. It must publish the prerelease assets and checksums, update `spacedock@next`, reconcile `next` because pre2 is newer than its current pre1 stamp, and bump the edge marketplace calendar key. Do not edit the workflow.

This path preserves one source identity: `v0.27.0-pre2 -> REL_SHA -> BASE_SHA`, where `BASE_SHA` is captured from `origin/main` only after `gqs` merges and the stamp diff is constrained to three files. Goreleaser supplies the binary version stamp; the existing edge-advance path supplies the installed pre2 plugin.

### Clean install and provider-free pilot

Use the exact published arm64 edge archive with a fresh host configuration. The Homebrew `spacedock@next` cask is remotely verified at pre2, but clean cask installation is unprovable in this sandbox because Homebrew cannot open its host lock path; preserve that as one host/sandbox evidence friction and do not retry or fix host machinery. Do not use a checkout binary, `SPACEDOCK_BIN`, `--plugin-dir`, or an existing plugin cache.

1. Download `spacedock_0.27.0-pre2_darwin_arm64_edge.tar.gz` and `checksums.txt` from the published prerelease into the authorized disposable repository, verify its SHA-256, and extract the released binary there. Create a fresh `CODEX_HOME`, run that binary first on a clean PATH with `SPACEDOCK_BIN` unset, then run `spacedock install --host codex` plus `spacedock doctor --host codex`. Require the binary and installed plugin to report the `0.27` line, with the binary reporting `0.27.0-pre2`.
2. Create a task-specific standalone repository with `mktemp -d /tmp/spacedock-pre2-pilot.XXXXXX`, and retain its exact path through validation. Copy the checked-in disposable fixture `skills/integration/testdata/entity-label-drive/experiment-workflow` there, initialize it as a standalone Git repository, reset `001-prompt-caching-latency.md` to the gated `backlog` stage, replace its report with a valid `## Stage Report: backlog`, and commit a baseline. This fixture's next stage is the nonterminal, nongated `implementation` stage. Do not write the main checkout or common `.git/info/exclude`, create another sprint task, add product machinery, or track the pilot repository in the release worktree or state checkout.
3. From a fresh First Officer session, use the cloned entity Markdown as the existing Artifact and run ordinary `gate prepare`. Require exit 0, an open room containing only `gate-briefing.json` and `request.json`, and a local-clone commit that binds the emitted room before presentation.
4. Present one legible chat review. It must name the experiment and `backlog` stage, chosen direction, recommendation, bound snapshot, checklist, and concrete decision effect. Use no provider override, caller-authored JSON, direct frontmatter edit after the baseline, or Reference.
5. Record the captain's exact semantic decision with `gate record --decision ... --actor person:captain`; commit the closed Resolution in the local clone before ending the session.
6. Start a cold First Officer session with the same installed bundle and fresh conversational context. Require boot to discover the experiment as approved and awaiting advance, then run `gate consume` once and commit the result.
7. Prove the consequence on disk: status is `implementation` and the application is `consumed`. Record the cold boot/status projection and implementation dispatch-build result; the known inability to project/build the entered implementation stage is retained as a sprint-critical limitation rather than an AC blocker. Hash the entity, repeat `gate consume`, and require nonzero exit plus byte-identical state.

The disposable fixture supplies the smallest real post-release journey with a gated stage followed by a nonterminal working stage. Keeping its standalone Git repository through this task's validation preserves durable evidence without adding a permanent fixture, sprint entity, or product surface.

## Retained evidence

The implementation report records:

- `BASE_SHA`, `REL_SHA`, their parent relation and three-file diff, the annotated tag object, local check commands, and exit results;
- Runtime Live E2E run `30543637878`, the explicit waiver reason, release-run URL and step summary, waiver-variable deletion, release asset names, checksum verification, `next` ancestry/version/calendar evidence, and cask version;
- the remotely published pre2 cask metadata; all three classified Homebrew host/sandbox attempts as one friction; the exact published arm64 edge archive and checksum result; clean `CODEX_HOME`; install and doctor results; exact binary/plugin versions; and absence of checkout overrides;
- the retained local clone's absolute path, baseline/bind/decision/consume commit SHAs, `git log` and `git fsck` result, prepared-room path and digest, room inventory, chat session reference, cold-boot discovery, post-consume entity state, actionable implementation projection/build result, and repeat-consume exit plus before/after hash; and
- every friction's journey step, reproduction, observable impact, evidence, class, owner, and next action.

Do not paste full logs into the entity. Cite run URLs, commit SHAs, concise command results, and durable state paths.

## Finding classes

- **Release blocker:** a red local format/full/race check; a stamp or tag identity breach; release workflow failure; missing or bad checksums/assets; edge cask or clean plugin install failure; or a binary/plugin version mismatch. Stop before or during publication and report the exact failed precondition. This task takes no unrelated fix.
- **Sprint-critical:** the installed path cannot prepare and bind the room, present a comprehensible chat decision, record authentic authority, rediscover pending approval cold, or consume exactly once; any digest, actor, status, or application corruption also belongs here. The pre2 inability to expose the entered `implementation` stage as its own actionable dispatch target is a known finding to preserve and route to `gqs`, without patching it or blocking this prerelease.
- **Deferrable:** convenience, rendering polish that does not obscure the decision, provider integration, extra Artifact or Reference support, dashboards, compatibility, or generalization. Record it without widening pre2.

## Acceptance criteria

**AC-1 — Published candidate.** `BASE_SHA` is captured from freshly fetched `origin/main`. Annotated tag `v0.27.0-pre2` names `REL_SHA`, a three-file stamp commit whose parent is `BASE_SHA`; local format, full Go, and race checks pass; the release audit records the captain-authorized waiver against green run `30543637878` without treating it as release-SHA proof; the waiver variable is removed; and the release, checksums, `next` reconcile, edge calendar bump, and `spacedock@next` cask complete successfully.

Verified by the tag object and commit diff; named command exits; release step summary; absent repository variable; GitHub release assets; checksum comparison; `origin/next` ancestry and manifest/calendar bytes; and cask metadata. Changing the tag target, adding a fourth stamp file, retaining the waiver, corrupting an asset, or leaving `next` at pre1 makes this criterion fail.

**AC-2 — Installed provider-free journey.** The checksum-verified published `spacedock_0.27.0-pre2_darwin_arm64_edge.tar.gz` binary and clean Codex plugin install drive a local standalone clone of the checked-in disposable experiment fixture from its gated `backlog` stage through prepare, ordinary-chat presentation, exact captain decision recording, cold-session discovery, one successful consume, and entry into the nongated, nonterminal `implementation` stage without caller-authored JSON, post-baseline direct state editing, a provider, or a Reference. Homebrew clean-install proof is explicitly unproved in this sandbox, with remote pre2 cask metadata and the exact host-lock failure retained.

Verified by the published-archive checksum, installed version/doctor output, retained local-clone room and Git commits, captain-visible chat, cold boot projection, and final on-disk status. A same-session resume, hidden checkout override, missing presentation fact, tracked sprint/product pilot state, or non-implementation successor makes this criterion fail.

**AC-3 — Exactly-once authority with retained actionability evidence.** After the cold First Officer consumes the approval, the cloned entity has `status: implementation` and the application has `state: consumed`; a second consume exits nonzero and leaves the entity byte-identical. Cold boot/status projection and implementation dispatch-build results are retained, with the known post-consume actionability failure classified sprint-critical and routed to `gqs` rather than blocking pre2.

Verified by the recorded Resolution/application, first consume output and local-clone commit, status projection, dispatch-build result, and a before/after SHA-256 around the second consume. A second advance, zero exit, or any byte change makes this criterion fail; skipped implementation projection or failed package build remains a recorded known limitation.

**AC-4 — Useful feedback.** The pilot report classifies every observed friction as release blocker, sprint-critical, or deferrable and records its step, reproduction, impact, evidence, owner, and next action. Tracked repository state contains only the release stamp and this task's reports; pilot gate state remains solely in the retained `/tmp` standalone clone.

Verified by reconciling the command/session evidence inventory against the findings section and git diffs. An unclassified failure, missing reproduction, or product fix outside the three-file stamp makes this criterion fail.

## Test plan

- **Release identity and local health, low cost:** inspect the tag/commit graph and diff; run format, full Go, race, and manifest-tag checks. These serve AC-1.
- **Publication and install, medium cost:** observe the existing release workflow, verify downloaded checksums, `next`, calendar, remote cask metadata, the exact published archive, clean host install, and doctor output. These serve AC-1 and AC-2; sandbox-blocked Homebrew installation is retained separately rather than substituted for archive proof.
- **Live chat journey, high judgment cost:** two real First Officer sessions and the captain exercise the disposable fixture's backlog gate. Durable local-clone room/entity commits, not transcript phrasing alone, grade AC-2.
- **Exactly-once negative control and actionability capture, low cost:** record cold projection/build, hash, repeat consume, and re-hash. This serves AC-3; removing the consumed-state guard must turn it red, while the known actionability gap is evidence rather than a blocker.
- **Triage audit, low cost:** map every nonzero exit or user/agent friction from the evidence inventory to one finding row. This serves AC-4.

A release-helper spike on 2026-07-30 ran the focused manifest/edge-decision tests green and exercised `edge-advance-decision v0.27.0-pre2` against `origin/next`; it printed `advance` with target pre2 versus next pre1. This serves AC-1 and would fail if the existing workflow classified pre2 as equal or older. The checked-in experiment fixture proves the required gated-to-working topology and supplies the disposable pilot shape. Existing fixtures prove gate prepare/record/consume, install CI proves both supported operating systems, and Runtime Live E2E run `30543637878` supplies waiver evidence across Claude, Codex, and Pi. It does not prove `BASE_SHA` or `REL_SHA`. The remaining unproved value is the published clean-install user journey itself, so the first post-publication action is the clean install and real fixture `gate prepare` capability probe.

The cold second session is necessary: a same-session continuation cannot prove durable discovery. The repeat consume is necessary: one successful consume alone cannot prove one-use authority. Provider presentation is unnecessary for the narrowed provider-free value. The post-consume actionability gap is intentionally left visible for `gqs`.

## Expected surface

- `main`: three files, three version substitutions, no other tracked diff — `.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, and `skills/first-officer/references/first-officer-shared-core.md`.
- state checkout: this `index.md` and its implementation report only; no pilot room or pilot entity.
- local disposable pilot repository: the exact path returned by `mktemp -d /tmp/spacedock-pre2-pilot.XXXXXX`, containing the copied fixture, locally reset entity/report, generated backlog room, and baseline/bind/decision/consume Git history retained through this task's validation. It is the only task-specific write allowed outside the assigned worktree and state entity.
- `next`: the existing release workflow's merge of `REL_SHA` plus one calendar-key substitution; no hand edit.
- external: one annotated tag, one GitHub prerelease, release archives/checksums, the edge cask update, and host-native plugin cache.

Any extra tracked `main` or state file, command grammar, stored-format field, authority change, runtime behavior change, permanent fixture, sprint task, or product fix breaches the approved surface and requires a design reset. Expected copied and generated files inside the disposable `/tmp` standalone pilot clone are not tracked surface.

## Documentation

No documentation diff is proposed. This task publishes and exercises behavior already documented in the gate command reference and concepts pages; it changes no user-visible command or contract text beyond the required release-version stamp and public tag notes.

## Explicit non-prerequisites

Subspace, provider evidence, cross-provider parity, multiple Artifacts, Reference rendering, the final `ph` walking skeleton, generalized advisory-round recording, dashboards, prototype compatibility, and host-specific convenience fixes are not prerequisites for pre2.

## Stage Report: ideation

- DONE: Define the exact v0.27.0-pre2 cut, identity, and green-run proof.
  The design makes merged `gqs` the sole prerequisite, captures `BASE_SHA` from post-merge `origin/main`, binds the three-file `REL_SHA` stamp to that base, and records run `30543637878` only as captain-authorized waiver evidence. The helper spike plus full/race suites passed.
- SKIPPED: Define clean-install chat and Subspace gate journeys with retained evidence.
  Captain narrowed pre2 to one provider-free chat journey; the design retains release, install, room, decision, cold-resume, and consume evidence and defers Subspace.
- DONE: Declare release blockers versus sprint-critical and deferrable findings, plus expected surface.
  The body defines three classes, a three-file `main` stamp, an ignored standalone pilot clone with retained local Git evidence, `next` effects, and the no-product-change boundary.

### Summary

After merged `gqs`, the design captures the actual release base, cuts `v0.27.0-pre2` through existing release machinery and an auditable captain waiver, then pilots one clean-installed chat gate across warm and cold First Officer sessions in a disposable local clone. It keeps provider work and every discovered follow-up outside the cut while retaining enough release and local Git evidence to judge the real user journey.

## Pilot Findings

- **Deferrable — Homebrew clean-install proof was blocked by host/sandbox permissions.** Step/reproduction: three normal `brew reinstall --cask spacedock-dev/tap/spacedock@next` attempts progressed through the existing cache `EPERM`, isolated cache plus inaccessible `~/.homebrew`, then isolated XDG trust to `/opt/homebrew/var/homebrew/locks/...incomplete.download.lock: Operation not permitted`; the last host view still resolved stale pre0. Impact: this sandbox could not prove the cask install path, although remote cask metadata was pre2 and release checksums were green. Evidence: exact exits retained in session; owner: host/sandbox evidence environment; next action: re-run on a normal writable Homebrew host, with no product or host-machinery fix in this task.
- **Deferrable — a fully empty `CODEX_HOME` requires authentication bootstrap.** Step/reproduction: the first installed-bundle Codex launch opened the sign-in selector. Impact: the pilot could not start until only the existing `auth.json` was seeded into the otherwise fresh home. Evidence: first launch session plus `codex login status`; owner: host setup/docs; next action: state that clean-home pilots may reuse auth while excluding config, sessions, and plugin cache.
- **Deferrable — automatic workflow discovery did not commission the disposable fixture root.** Step/reproduction: first session `019fb37b-b1e9-77d2-be49-b7756d71427f` ran boot from the fixture root and received `no commissioned Spacedock workflow`; explicit `--workflow-dir /tmp/spacedock-pre2-pilot.YHcWXZ` succeeded. Impact: one cold-start clarification was required before gate preparation. Evidence: the same session's stopped first turn and successful resume; owner: workflow discovery/pilot ergonomics; next action: pass the exact disposable workflow directory explicitly.
- **Deferrable — `state commit` is a no-op for the inline fixture.** Step/reproduction: after prepare, record, and consume it printed `Inline workflow — entities live beside the README; nothing to commit to a state checkout.` Impact: the First Officer had to make the required standalone-repository commits directly. Evidence: commits `cfa0318`, `67d743b`, and `7584d27`; owner: inline-workflow commit ergonomics; next action: preserve the explicit local Git commit instruction for disposable pilots.
- **Sprint-critical — consumed implementation is not actionable in pre2.** Step/reproduction: cold boot after consume projected `current=implementation,next=review,worktree=no`; explicit implementation package build failed `worktree stage 'implementation' but entity has no worktree path`. Impact: the entered implementation stage cannot be dispatched without forbidden state repair. Evidence: cold session `019fb37f-199d-7fd1-bb8e-8d374179ac4e`, entity commit `7584d27`; owner: `gqs` (`dispatch-entered-stage-after-gate-consume`); next action: finish and validate that sprint task, without patching pre2.

## Retained Pilot Evidence

- Release identity: `BASE_SHA=4eeb94e9b1f7d2e407961e28c941e422c28749fc`; `REL_SHA=5e7f1ffa08721c062fa9ae82636549a635983e95`; sole parent equals base; the release diff is exactly the two plugin manifests plus First Officer shared core, one substitution each. Annotated tag object `7b678576bcb1a289bd54cdb22cd4caad8a60fe00` targets `REL_SHA`.
- Local health: `gofmt -w ./cmd ./internal` left a clean diff; `go test ./...`, `go test ./... -race`, and `manifest-tag-gate v0.27.0-pre2` exited 0. Each full suite would fail on a package regression; the manifest gate would fail on a fourth/mismatched stamp.
- Publication: release run [30552325976](https://github.com/spacedock-dev/spacedock/actions/runs/30552325976) completed all four jobs green. Its E2E step records run `30543637878` at `a0391fdbdb957d703ee4f836a5dfd5e21cc70153` only as waiver evidence with base/release SHAs; the repository waiver variable is absent.
- Assets/channel: the prerelease has checksums plus eight platform/channel archives and the journey ledger; every archive passed `shasum -a 256 -c`. `origin/next=26057e74b59d5037eea93f959ac645f33aa97f33` contains `REL_SHA`, both manifests are pre2, calendar key is `0.0.2026073001`, and remote `spacedock@next` cask metadata is `0.27.0-pre2`.
- Installed archive: `spacedock_0.27.0-pre2_darwin_arm64_edge.tar.gz` matched SHA-256 `1345444854df436035fb8b059f4e9dc543e7934ca522d56b79b5eed268df2367`. With it first on PATH and `SPACEDOCK_BIN` absent, binary and clean-home plugin both reported `0.27.0-pre2`; install and doctor exited 0.
- Pilot repository: `/tmp/spacedock-pre2-pilot.YHcWXZ`; commits `bb6bf6d` baseline, `cfa0318` binding, `67d743b` exact captain decision, `7584d27` consume; `git fsck --full` and status were clean. Room `001-prompt-caching-latency/review/backlog/briefing-1` contains only `gate-briefing.json` and `request.json`; Briefing digest is `sha256:f98e3db5f9d031cf3b9f227cc4a7d2565df3ccd3f5985ebdaefe8002abe780e8`.
- Journey: first session `019fb37b-b1e9-77d2-be49-b7756d71427f` presented and recorded the exact `person:captain` decision; cold session `019fb37f-199d-7fd1-bb8e-8d374179ac4e` discovered `approved-awaiting-advance`, consumed once into `implementation`, and retained the actionability failure. Repeat consume exited 1 and entity SHA-256 stayed `72ef8a9e8002a0952bbbb2fded3426fbc26631ba58e1fbe5fb495344180e14f2`.

## Stage Report: implementation

- DONE: Capture `BASE_SHA`, run required local checks, stamp and publish `v0.27.0-pre2`, and remove the E2E waiver.
  Captain removed `gqs` as a prerequisite; fresh `origin/main` supplied the retained base, release commit `5e7f1ffa0` contains only the three allowed stamps, run `30552325976` is green, and the waiver is absent.
- DONE: Verify release assets, checksums, next reconciliation, edge calendar/cask, and clean installed 0.27 bundle.
  All release archives verified, `next` and remote cask report pre2, and the checksum-verified arm64 edge archive installed plugin pre2 into a fresh authenticated `CODEX_HOME`; Homebrew proof remains explicitly unproved by sandbox permissions.
- DONE: Drive the provider-free chat fixture through cold discovery, exactly-once consume, actionable implementation, and classified friction.
  The retained repository has baseline/bind/decision/consume commits, two fresh session IDs, exact captain authority, a nonzero byte-clean second consume, and the expected sprint-critical actionability failure routed to `gqs`.
- DONE: Update the task body and acceptance criteria for the captain override and archive-install proof.
  The body removes the hard `gqs` dependency, treats post-consume actionability as retained evidence, and distinguishes unproved Homebrew installation from the successful exact-published-archive proof.

### Summary

Published `v0.27.0-pre2` from a constrained three-file release commit, verified the complete release/channel output, and removed the one-shot waiver immediately after the gate consumed it. A clean installed bundle then completed the real provider-free gate journey across warm and cold First Officer sessions, proving durable exactly-once captain authority while exposing the known `gqs` actionability gap without patching it.
