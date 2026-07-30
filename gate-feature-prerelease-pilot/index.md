---
title: Cut and field-test the current gate feature as a prerelease
status: implementation
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

After `gqs` merges, publish `v0.27.0-pre2` from the resulting releasable `main` gate surface, then use the installed edge bundle for one real, provider-free gate journey. The pilot must show that an ordinary user and two successive First Officer sessions can prepare, understand, authorize, resume, and spend one decision without hand-authored authority.

## Problem

The gate commands and First Officer contract have strong fixture and live-lane coverage, but no user has yet installed one released bundle and driven the complete provider-free lifecycle. Waiting for sprint closure would hide install, presentation, cold-resume, and exactly-once friction until the design is harder to change.

## Boundary

This is an instrumented prerelease, not the durable-decisions sprint acceptance walk. It changes no gate grammar, stored format, authority rule, or runtime implementation. It ships the merged surface, runs one chat journey, and reports friction; it adds no release machinery and absorbs no follow-up fix.

`gqs` (`dispatch-entered-stage-after-gate-consume`) is the sole pre2 prerequisite. Do not cut until it is terminal `PASSED`, its implementation is merged, and its merge commit is an ancestor of freshly fetched `origin/main`. Only then capture `BASE_SHA=$(git rev-parse origin/main)`.

Captain authorized the release-time E2E waiver on 2026-07-30. Runtime Live E2E run `30543637878` was green at PR #575 head `a0391fdbdb957d703ee4f836a5dfd5e21cc70153`; it is waiver evidence only. It is neither an exact release-SHA run nor proof that its head is the post-`gqs` release base. The report must preserve that distinction.

## Proposed approach

### Cut and publish

1. Require `gqs` terminal `PASSED` and merged. Fetch `origin/main` and tags, verify the `gqs` merge commit is an ancestor of `origin/main`, stop if `v0.27.0-pre2` already exists, then capture and retain `BASE_SHA=$(git rev-parse origin/main)`.
2. Create a release worktree from `BASE_SHA`. Run `gofmt -w ./cmd ./internal`, require a clean diff, then run `go test ./...` and `go test ./... -race`.
3. Run `spacedock-release stamp-version 0.27.0-pre2` against `.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, and `skills/first-officer/references/first-officer-shared-core.md`. Commit only those three files. Define `REL_SHA` as that commit; require its sole parent to equal `BASE_SHA` and its diff to contain only the three version substitutions.
4. Run `manifest-tag-gate v0.27.0-pre2` locally. Push the release branch to `main`, then generate and review public notes with `spacedock-release notes 0.27.0-pre2`. The resulting annotated tag must point to `REL_SHA` and carry a nonempty body.
5. Set `SPACEDOCK_E2E_GATE_WAIVER` to a reason naming the captain authorization, run `30543637878`, its green head, `BASE_SHA`, and `REL_SHA`, and explicitly state that the run is waiver evidence rather than release-SHA proof. Push only `v0.27.0-pre2`. Capture the release run and its waiver step summary. Delete the repository variable after the release run reads it, including on release failure, and verify that it is absent.
6. Require the existing release workflow to succeed. It must publish the prerelease assets and checksums, update `spacedock@next`, reconcile `next` because pre2 is newer than its current pre1 stamp, and bump the edge marketplace calendar key. Do not edit the workflow.

This path preserves one source identity: `v0.27.0-pre2 -> REL_SHA -> BASE_SHA`, where `BASE_SHA` is captured from `origin/main` only after `gqs` merges and the stamp diff is constrained to three files. Goreleaser supplies the binary version stamp; the existing edge-advance path supplies the installed pre2 plugin.

### Clean install and provider-free pilot

Use the published `spacedock@next` cask with a fresh host configuration. Do not use a checkout binary, `SPACEDOCK_BIN`, `--plugin-dir`, or an existing plugin cache.

1. Reinstall `spacedock-dev/tap/spacedock@next`, create a fresh `CODEX_HOME`, and run `spacedock install --host codex` plus `spacedock doctor --host codex`. Require the binary and installed plugin to report the `0.27` line, with the binary reporting `0.27.0-pre2`.
2. Add the exact local pilot path `.spacedock/pilots/v0.27.0-pre2/experiment-workflow` to the main checkout's untracked `.git/info/exclude` and verify it with `git check-ignore`. Copy the checked-in disposable fixture `skills/integration/testdata/entity-label-drive/experiment-workflow` there, initialize it as a standalone Git repository, reset `001-prompt-caching-latency.md` to the gated `backlog` stage, replace its report with a valid `## Stage Report: backlog`, and commit a baseline. This fixture's next stage is the nonterminal, nongated `implementation` stage. Do not create another sprint task, add product machinery, or track the pilot repository in `main` or the state checkout.
3. From a fresh First Officer session, use the cloned entity Markdown as the existing Artifact and run ordinary `gate prepare`. Require exit 0, an open room containing only `gate-briefing.json` and `request.json`, and a local-clone commit that binds the emitted room before presentation.
4. Present one legible chat review. It must name the experiment and `backlog` stage, chosen direction, recommendation, bound snapshot, checklist, and concrete decision effect. Use no provider override, caller-authored JSON, direct frontmatter edit after the baseline, or Reference.
5. Record the captain's exact semantic decision with `gate record --decision ... --actor person:captain`; commit the closed Resolution in the local clone before ending the session.
6. Start a cold First Officer session with the same installed bundle and fresh conversational context. Require boot to discover the experiment as approved and awaiting advance, then run `gate consume` once and commit the result.
7. Prove the consequence on disk: status is `implementation`, the application is `consumed`, and, because pre2 contains `gqs`, cold boot/status projects `current=implementation,next=implementation`. Require the installed binary to build an `implementation` dispatch package, but do not dispatch or complete that stage. Hash the entity, repeat `gate consume`, and require nonzero exit plus byte-identical state.

The disposable fixture supplies the smallest real post-release journey with a gated stage followed by a nonterminal working stage. Keeping its standalone Git repository through this task's validation preserves durable evidence without adding a permanent fixture, sprint entity, or product surface.

## Retained evidence

The implementation report records:

- the `gqs` merge SHA and ancestry proof, `BASE_SHA`, `REL_SHA`, their parent relation and three-file diff, the annotated tag object, local check commands, and exit results;
- Runtime Live E2E run `30543637878`, the explicit waiver reason, release-run URL and step summary, waiver-variable deletion, release asset names, checksum verification, `next` ancestry/version/calendar evidence, and cask version;
- clean `CODEX_HOME`, install and doctor results, exact binary/plugin versions, and absence of checkout overrides;
- the retained local clone's absolute path, baseline/bind/decision/consume commit SHAs, `git log` and `git fsck` result, prepared-room path and digest, room inventory, chat session reference, cold-boot discovery, post-consume entity state, actionable implementation projection/build result, and repeat-consume exit plus before/after hash; and
- every friction's journey step, reproduction, observable impact, evidence, class, owner, and next action.

Do not paste full logs into the entity. Cite run URLs, commit SHAs, concise command results, and durable state paths.

## Finding classes

- **Release blocker:** a red local format/full/race check; a stamp or tag identity breach; release workflow failure; missing or bad checksums/assets; edge cask or clean plugin install failure; or a binary/plugin version mismatch. Stop before or during publication and report the exact failed precondition. This task takes no unrelated fix.
- **Sprint-critical:** the installed path cannot prepare and bind the room, present a comprehensible chat decision, record authentic authority, rediscover pending approval cold, consume exactly once, or expose the entered `implementation` stage as its own actionable dispatch target; any digest, actor, status, or application corruption also belongs here. Preserve evidence and route the finding to the existing sprint owner without patching it in this task.
- **Deferrable:** convenience, rendering polish that does not obscure the decision, provider integration, extra Artifact or Reference support, dashboards, compatibility, or generalization. Record it without widening pre2.

## Acceptance criteria

**AC-1 — Published candidate.** After `gqs` is terminal `PASSED` and its merge is proven ancestral to fetched `origin/main`, `BASE_SHA` is captured from that ref. Annotated tag `v0.27.0-pre2` names `REL_SHA`, a three-file stamp commit whose parent is `BASE_SHA`; local format, full Go, and race checks pass; the release audit records the captain-authorized waiver against green run `30543637878` without treating it as release-SHA proof; the waiver variable is removed; and the release, checksums, `next` reconcile, edge calendar bump, and `spacedock@next` cask complete successfully.

Verified by the tag object and commit diff; named command exits; release step summary; absent repository variable; GitHub release assets; checksum comparison; `origin/next` ancestry and manifest/calendar bytes; and cask metadata. Changing the tag target, adding a fourth stamp file, retaining the waiver, corrupting an asset, or leaving `next` at pre1 makes this criterion fail.

**AC-2 — Installed provider-free journey.** A fresh `spacedock@next` binary and clean Codex plugin install drive a local standalone clone of the checked-in disposable experiment fixture from its gated `backlog` stage through prepare, ordinary-chat presentation, exact captain decision recording, cold-session discovery, one successful consume, and entry into the nongated, nonterminal `implementation` stage without caller-authored JSON, post-baseline direct state editing, a provider, or a Reference.

Verified by the installed version/doctor output, retained local-clone room and Git commits, captain-visible chat, cold boot projection, and final on-disk status. A same-session resume, hidden checkout override, missing presentation fact, tracked sprint/product pilot state, or non-implementation successor makes this criterion fail.

**AC-3 — Exactly-once authority and actionable entry.** After the cold First Officer consumes the approval, the cloned entity has `status: implementation` and the application has `state: consumed`; cold boot/status names `implementation` as both current and next, an implementation dispatch package builds, and a second consume exits nonzero and leaves the entity byte-identical.

Verified by the recorded Resolution/application, first consume output and local-clone commit, exact status projection, dispatch-build output, and a before/after SHA-256 around the second consume. A successor projection that skips implementation, failed package build, second advance, zero exit, or any byte change makes this criterion fail.

**AC-4 — Useful feedback.** The pilot report classifies every observed friction as release blocker, sprint-critical, or deferrable and records its step, reproduction, impact, evidence, owner, and next action. Tracked repository state contains only the release stamp and this task's reports; pilot gate state remains solely in the ignored standalone clone.

Verified by reconciling the command/session evidence inventory against the findings section and git diffs. An unclassified failure, missing reproduction, or product fix outside the three-file stamp makes this criterion fail.

## Test plan

- **Release identity and local health, low cost:** inspect the tag/commit graph and diff; run format, full Go, race, and manifest-tag checks. These serve AC-1.
- **Publication and install, medium cost:** observe the existing release workflow, verify downloaded checksums, `next`, calendar, cask, clean host install, and doctor output. These serve AC-1 and AC-2.
- **Live chat journey, high judgment cost:** two real First Officer sessions and the captain exercise the disposable fixture's backlog gate. Durable local-clone room/entity commits, not transcript phrasing alone, grade AC-2.
- **Exactly-once negative control, low cost:** hash, repeat consume, and re-hash. This serves AC-3; removing the consumed-state guard must turn it red.
- **Triage audit, low cost:** map every nonzero exit or user/agent friction from the evidence inventory to one finding row. This serves AC-4.

A release-helper spike on 2026-07-30 ran the focused manifest/edge-decision tests green and exercised `edge-advance-decision v0.27.0-pre2` against `origin/next`; it printed `advance` with target pre2 versus next pre1. This serves AC-1 and would fail if the existing workflow classified pre2 as equal or older. The checked-in experiment fixture proves the required gated-to-working topology and supplies the disposable pilot shape. Existing fixtures prove gate prepare/record/consume, install CI proves both supported operating systems, and Runtime Live E2E run `30543637878` supplies waiver evidence across Claude, Codex, and Pi. It does not prove `BASE_SHA` or `REL_SHA`. The remaining unproved value is the published clean-install user journey itself, so the first post-publication action is the clean install and real fixture `gate prepare` capability probe.

The cold second session is necessary: a same-session continuation cannot prove durable discovery. The repeat consume is necessary: one successful consume alone cannot prove one-use authority. Provider presentation is unnecessary for the narrowed provider-free value.

## Expected surface

- `main`: three files, three version substitutions, no other tracked diff — `.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, and `skills/first-officer/references/first-officer-shared-core.md`.
- state checkout: this `index.md` and its implementation report only; no pilot room or pilot entity.
- local ignored pilot repository: `.spacedock/pilots/v0.27.0-pre2/experiment-workflow`, excluded through the main checkout's `.git/info/exclude` and containing the copied fixture, locally reset entity/report, generated backlog room, and baseline/bind/decision/consume Git history retained through this task's validation.
- `next`: the existing release workflow's merge of `REL_SHA` plus one calendar-key substitution; no hand edit.
- external: one annotated tag, one GitHub prerelease, release archives/checksums, the edge cask update, and host-native plugin cache.

Any extra tracked `main` or state file, command grammar, stored-format field, authority change, runtime behavior change, permanent fixture, sprint task, or product fix breaches the approved surface and requires a design reset. Expected copied and generated files inside the ignored standalone pilot clone are not tracked surface.

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
