---
id: 6yyyyemkqwsett3g1c991w9f
title: Make First Officers operate the recorded gate lifecycle
status: ideation
source: "Durable-decisions dogfood audit: PRs #557/#560 shipped gate commands without the planned FO operating contract, 2026-07-23"
started: 2026-07-23T02:01:56Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-first-officer-gate-command-lifecycle
issue:
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:docs-dev:6y:ideation
    records:
        - id: gate:docs-dev:6y:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6y-ideation-1
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-1:revision-1
                digest: sha256:39dada7e95453a8738f41ca886881deebfa31edf16ba677cf95a580596f7dbc6
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:1
                briefing: briefing:docs-dev:6y:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-23T02:23:39.812888Z"
                decision: approve
                reason: Ideation codifies the demonstrated 3k/h1 lifecycle, makes every transition command load-bearing, captures all observed friction, and preserves the no-recorder/no-production-Go boundary.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6y-ideation-2
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-2:revision-1
                digest: sha256:53b3cd4c9ba72ecbe375bb2a638cba5cd840c0e19481342b52da1cf8db5f11f7
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:2
                briefing: briefing:docs-dev:6y:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-23T06:38:37.347699Z"
                decision: approve
                reason: Cycle-2 ideation moves the complete lifecycle behind one deferred gate trigger, preserves the boot-core ceiling and strict spawn ACs, closes every gate-entry route including headless, and names all missing behavioral proof without changing product semantics.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6y-ideation-3
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-3:revision-2
                digest: sha256:b6ecdb249de0c91b3857b218cf2464ab7b8bafbdb1d2fce70c4f8526d0c827ce
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:3
                briefing: briefing:docs-dev:6y:ideation:attempt-3:revision-2
                by: agent:first-officer
                at: "2026-07-24T00:14:56.965625Z"
                decision: approve
                reason: Canonical AC repair now preserves product value while removing unsupported proof obligations and cutting at least 399 LOC.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6y-ideation-4
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-4:revision-1
                digest: sha256:ae9e79e14e9df46d29af40a9e570c2af47b0bb5dd049b9f30cf57b73a1607036
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-4
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:4
                briefing: briefing:docs-dev:6y:ideation:attempt-4:revision-1
                by: agent:first-officer
                at: "2026-07-24T09:08:13.257256Z"
                decision: approve
                reason: Cycle 7 preserves the authority and durability outcomes, removes unproven ceremony, passed the focused gate tests named in ideation, and independent staff re-review reports APPROVE with no material findings.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6y-ideation-5
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-5:revision-1
                digest: sha256:af93b06086234aa95c8ad1a98bf52bea7866c4ceeb86ce8ba3140881aa761ad3
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-5
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:5
                briefing: briefing:docs-dev:6y:ideation:attempt-5:revision-1
                by: agent:first-officer
                at: "2026-07-24T13:55:08.534375Z"
                decision: approve
                reason: Cycle 13 preserves the three durable authority mutations, fixes the observed terminal merge deadlock through the existing ceremony, removes the false cache proof, constrains Pi evidence to root-assistant output, and declares a net-deleting nine-file implementation boundary.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6y-ideation-6
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-6:revision-1
                digest: sha256:47bcff38eb3425ded2ee321c3639ad3a016a8db4f455a0e7029d65ad335af584
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-6
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:6
                briefing: briefing:docs-dev:6y:ideation:attempt-6:revision-1
                by: agent:first-officer
                at: "2026-07-24T18:24:50.60913Z"
                decision: approve
                reason: Cycle 18 closes both staff-review findings, restores supported rejection semantics without a new schema, makes live procedure ownership falsifiable, and limits implementation to seven existing files with a +95-addition hard stop; independent staff re-review approves with no material finding.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6y-ideation-7
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-7:revision-1
                digest: sha256:9c2501ebc0b72b6fe4ae64119265832a341d89713fc80b569a8c18d76c7a62bf
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-7
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:7
                briefing: briefing:docs-dev:6y:ideation:attempt-7:revision-1
                by: agent:first-officer
                at: "2026-07-24T19:14:28.321408Z"
                decision: approve
                reason: Independent staff review found no material issue; the repaired eight-file design gives gate presentation one owner, makes exact-one root visibility deletion-sensitive, and preserves the existing recorder, provider, and runtime boundaries.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6y-ideation-8
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-8:revision-1
                digest: sha256:c9a7c203e6d1db996be3325f5d55bcc3bc7ebe6d40d4b2bae9affba5d6f8bd10
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-8
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:8
                briefing: briefing:docs-dev:6y:ideation:attempt-8:revision-1
                by: agent:first-officer
                at: "2026-07-24T20:08:36.771738Z"
                decision: approve
                reason: Independent staff review found no material issue; the unreleased-v1 reset removes prototype marker and unbound-review assumptions, binds presentation to successful durable evidence, and keeps all supported-host and authority ACs falsifiable within an explicit ceiling.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:6y-ideation-9
              briefing:
                id: briefing:docs-dev:6y:ideation:attempt-9:revision-1
                digest: sha256:202522443343dc2cf1c18284f79459f4dc3b4bde62e30a9061c560df90822292
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-9
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:ideation:9
                briefing: briefing:docs-dev:6y:ideation:attempt-9:revision-1
                by: agent:first-officer
                at: "2026-07-24T21:18:39.847514Z"
                decision: approve
                reason: Independent staff review found no material issue; the repaired design restores the Captain-approved host proof split, makes delegated authority byte-exact, retains Pi failure evidence before cleanup, and requires measured deletion before any expansion.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:docs-dev:6y:implementation
          stage: implementation
          attempts:
            - id: gate-attempt:6y-implementation-1
              briefing:
                id: briefing:docs-dev:6y:implementation:attempt-1:revision-1
                digest: sha256:3b5dcf4d8d48d6d3991976efe622e9d577d465b67bb82e52e380310ac0da1334
                digest-domain: canonical-bytes
                room-ref: ./review/implementation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:implementation:1
                briefing: briefing:docs-dev:6y:implementation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-23T06:28:08.375697Z"
                decision: revise
                reason: The boot-resident lifecycle exceeds the hard shared-core ceiling by 5,534 bytes beyond available headroom, leaves a headless gate-entry bypass, and lacks required live-spawn and adversarial proof; ACs remain unchanged and require deferred gate-triggered topology re-ideation.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: feedback
                target-stage: implementation
                state: superseded
            - id: gate-attempt:6y-implementation-2
              briefing:
                id: briefing:docs-dev:6y:implementation:attempt-2:revision-1
                digest: sha256:cc82fcd0474089c65e415ba09545dc18c26fd50b47e646785a5dda6cb827f61a
                digest-domain: canonical-bytes
                room-ref: ./review/implementation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:implementation:2
                briefing: briefing:docs-dev:6y:implementation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-23T11:09:18.116708Z"
                decision: revise
                reason: 'The deferred lifecycle implementation preserves the end value, but its remaining failures are proof-boundary defects: a model-reconstructed command-only prompt mutant, an all-routes-by-all-hosts live matrix, and an unexposed Codex public-stream handle are not valid mandatory oracles; Pi still needs real async completion waiting and Claude needs a final corrected positive run.'
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement.
              application:
                action: feedback
                target-stage: implementation
                state: pending
        - id: gate:docs-dev:6y:validation
          stage: validation
          attempts:
            - id: gate-attempt:6y-validation-1
              briefing:
                id: briefing:docs-dev:6y:validation:attempt-1:revision-1
                digest: sha256:1bddbf112a8367cce288e428a4900715221e62d04ad35098e7e9dec50841baee
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:validation:1
                briefing: briefing:docs-dev:6y:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-24T17:13:28.512931Z"
                decision: approve
                reason: Validation reproduced AC-1 through AC-8 at exact tip b99f9c66; focused, full, race, docs, Codex, Claude, Pi, detached mutants, and final-tip Roborev are green with no candidate-scope material finding.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: done
                state: consumed
                blockers: []
            - id: gate-attempt:6y-validation-2
              briefing:
                id: briefing:docs-dev:6y:validation:attempt-2:revision-1
                digest: sha256:c5fc51c725cd43460b01b65af4558fcfed164eb141c9cd7b8ab67b87c606557a
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:validation:2
                briefing: briefing:docs-dev:6y:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-24T17:19:34.698945Z"
                decision: revise
                reason: 'PR #565 offline CI and the isolated clean-config reproduction fail at the first fixture-backed state commit because the temporary state repository lacks local Git author identity; this is a material evidence defect and the candidate cannot merge.'
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: feedback
                target-stage: implementation
                state: superseded
            - id: gate-attempt:6y-validation-3
              briefing:
                id: briefing:docs-dev:6y:validation:attempt-3:revision-1
                digest: sha256:0ae905b1aecc75768763f7d0c960af1f0894c7cf5d68c3f827e2b6ecfd6e957e
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:6y:validation:3
                briefing: briefing:docs-dev:6y:validation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-07-24T17:57:53.206098Z"
                decision: revise
                reason: Exact tip 13d70249 fixes clean-runner Git identity but detached controls prove four material defects across supported rejection routing, captain decision semantics, live prompt ownership, and root-visible Claude evidence; AC-1, AC-6, AC-7, and AC-8 fail and require a design reset.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: feedback
                target-stage: implementation
                state: pending
review-round:
    id: round:6yyyyemkqwsett3g1c991w9f:implementation:16
    stage: implementation
    cycle: 16
    briefing:
        id: briefing:first-officer-gate-command-lifecycle:implementation:round-16
        digest: sha256:5099e38f70c4c6de42a1a68a5abfab5e4666f34febd492fbf32b8cce0c8f12e2
        digest-domain: canonical-bytes
        room-ref: ./review/implementation/round-16
mod-block:
pr: "#565"
---

Make the normal First Officer gate path bind the exact reviewed package, record the authorized decision, and durably consume it before ordinary workflow dispatch.

## Problem

The recorder and application commands exist, but the shipped First Officer contract still assembles and presents gates only in prose. PR #557 (`fa240a76`) changed no First Officer skill files despite budgeting approximately ten lines for that integration; PR #560 (`f06cce04`) likewise added eligibility and consume without an FO caller. `skills/first-officer/references/first-officer-shared-core.md` names judgment, `present-gate`, feedback routing, advancement, and dispatch, but literally invokes none of `spacedock gate record`, `gate validate`, `gate eligibility`, or `gate consume`. This sprint therefore used the commands only under manual Captain directives, and an ordinary FO can still approve and advance outside the durable lifecycle.

The missing work is not a new decision model. It is the smallest procedure that makes the sprint FO's authority durable: retain a canonical package, bind its exact `briefing.json`, make and present the evidence judgment, record the Captain or delegated decision, consume that one-use application, durably commit the consumed state, then route the resulting status. A nonterminal target returns to ordinary dispatch; a terminal target enters the existing merge ceremony and has no dispatchable successor. Before h1 landed, the same dogfood stopped after recording the Resolution and manually changed workflow state. `gate consume` is the sole approval transition: it re-evaluates eligibility under the entity lock and atomically co-writes the next `status` and `application.state: consumed`. A separate `status --set` advance is a bypass and a contract violation.

The repository also exposed two CLI-boundary defects. First, the retained repo-root `./spacedock` identifies as `0.26.0+dev` yet lacks the gate surface. Version compatibility is therefore insufficient. One `spacedock gate --help` response already carries all command and semantic-form tokens needed to reject that stale launcher; four repeated subcommand probes do not add discrimination. Second, the CLI passes a relative retained-input path to a recorder that compares it with an absolute entity root, so `filepath.Rel` fails before mutation. The CLI must normalize `--briefing`, `--result`, and `--association` against its launch directory; callers should not manufacture absolute paths.

Validation Cycle 3 exposed four additional defects in the shipped-skill contract and its evidence. A feedback reviewer’s `REJECTED` recommendation no longer has an automatic correction route; ordinary Captain words `redo` and `reject` are not mapped to canonical recorder decisions; the host-neutral live prompt scripts the lifecycle it claims to prove; and the Claude extractor can mistake a child review for the root assistant review. The reset must correct those defects without adding recorder fields, a parallel lifecycle/controller, host-specific lifecycle copies, compatibility behavior, or transport forensics beyond Claude’s existing top-level parent marker.

Implementation Cycle 19 corrected those four defects but exposed a narrower presentation-ownership failure. Two unchanged goal-only Codex runs bound and committed the exact Briefing, formed the evidence judgment, loaded `present-gate`, recorded and committed delegated approval, consumed and committed it, and dispatched the successor, yet neither emitted the required six-field root-visible gate review. The retained first transcript records the decisive sequence: item 19 narrates the evidence and recommendation, item 20 reads `skills/present-gate/SKILL.md`, and item 21 immediately invokes `gate record ... --decision approve`; the review oracle fails at `capability/change: count=0`. After the lifecycle skill was amended to name all six labels, the second run loaded that corrected skill and failed identically after 190.29 seconds. This disproves a missing-goal or missing-label diagnosis.

The conflict is inside shipped ownership. `fo-gate-lifecycle` currently names a six-field concise review, while the presenter it invokes defines a different `Gate review` / `Chosen direction` / `Checklist` / `Assessment` / `Decision` template and says only “render ... then hand” the decision to the recorder. Loading the presenter therefore gives the runtime two render contracts and no observable completion barrier between presentation and mutation. Codex retained the judgment and optimized away the intermediate chat artifact. The correction must put the canonical review shape and the chat completion condition in `present-gate`, leave recorder/consumer authority unchanged, and keep the live oracle strict.

## Demonstrated baseline

The design adopts the landed behavior and historical dogfood; it does not reopen whether durable gates should exist.

1. The 3k validation room retained a canonical `review/validation/briefing-v1/briefing.json` plus a concise gate review and frozen references. State commit `2c616b7e` bound it; `77590ebd` recorded an `agent:first-officer` approval carrying the Captain decision as adoption provenance. The resulting record stayed at `status: validation`: 3k intentionally recorded the decision but did not apply it.
2. The remaining 3k closeout used manual lifecycle mutations (`mod-block`, PR state, and later terminal archive) because h1 was not yet available. Those commits are historical evidence for the exact integration gap, not a procedure to preserve.
3. h1 landed `gate eligibility` and `gate consume`. Its CLI fixture proves an approve closure yields `advance/pending`, eligibility prints `condition=approved-pending eligible=true`, and consume atomically changes `status` plus `application.state: consumed`; a second consume exits nonzero with `condition=consumed consumed=false` and byte-identical state.
4. The recorder's retained package boundary is settled: the manifest basename is exactly `briefing.json`; its canonical bytes independently define the complete artifact inventory. Existing artifact payloads remain URI + SHA references when the presentation resolver can reproduce those exact bytes. Mutable, cross-root, or otherwise unreproducible reviewed material is frozen as a room copy before publication. The recorder verifies but never copies artifact payloads. For folder-form entities, landed vn/PR #558 makes one `spacedock state commit <slug>` include the index and every non-ignored room artifact without sweeping sibling dirt.
5. Direct and delegated decisions remain distinct. A Captain who personally renders the chat decision is `person:captain`. An FO acting under delegated conn is `agent:first-officer`, supplies its evidence-bearing `--reason`, and stores the exact quoted grant in `--directive`. An exact provider Result uses `--result` plus its retained association and authorized actor; advisory adoption also names the authorizer.
6. The consumed terminal target is observable in the existing development workflow. `gate-review-presentation-command` is `status: done`, has `application.state: consumed` targeting `done`, and retains `mod-block: merge:pr-merge` plus open `pr: "#564"`. At this state, `status --next --json` reports an empty `dispatchable` array. The shared merge core can resume the recorded `mod-block`, but the shipped and development `pr-merge` startup/idle hooks currently scan only nonterminal entities, so they cannot observe this terminal in-flight PR and write the merge sentinel. The narrow repair is to include only terminal PR rows that still carry `mod-block: merge:pr-merge`; it does not create a second terminal transition or change consume.

The cycle-6 cardinality spike exercised the landed CLI/gate tests and inspected the write boundaries. Both semantic `gate record` forms validate the complete document and transition before atomic replacement and return the selected summary; a later close reads and validates the bound record again. `gate consume` reads and validates the record, re-evaluates reviewed-input currency, expected successor, blockers, and one-use state under the same entity lock, then atomically writes status plus consumed state. Therefore open/closed `gate validate` and pre-consume `gate eligibility` are diagnostic reads, not independent authority checks: deleting them leaves the three mutation guards intact, while deleting briefing record, decision/Result record, or consume destroys a required durable state transition. Focused `internal/gates` and `internal/ensigncycle` tests passed this baseline on 2026-07-24. Implementation turns this spike into the first three-command fixture and deletion controls.

## Operational lifecycle

The procedure belongs at the existing gated-stage branch in `first-officer-shared-core.md`, around `«gate.assemble-verdict»`; it does not move or duplicate `«gate.ac-cross-check»`, FO evidence judgment, or `present-gate` rendering. Before Captain presentation, a completed feedback-gate reviewer recommendation of `REJECTED` automatically invokes the existing `«feedback.route»` correction flow. That path does not mint a Captain Resolution, close or consume a gate, or dispatch a successor; its existing reviewer/worker correction round is the durable record. Every other outcome enters `fo-gate-lifecycle`.

### 0. Select a capable executable

Immediately before every gate lifecycle, freshly resolve `${SPACEDOCK_BIN:-spacedock}` and invoke exactly once:

```text
${SPACEDOCK_BIN:-spacedock} gate --help
```

Require one successful response containing `record`, `validate`, `eligibility`, `consume`, and the record forms/flags `--briefing`, `--result`, `--association`, `--decision`, `--actor`, and `--directive`. This is the minimum current capability fingerprint: all tokens come from the one gate help body, so four subcommand-help calls distinguish no additional stale binary.

Do not cache the result or compute executable path/content identity. The just-in-time probe naturally re-resolves a replacement, symlink/PATH retarget, or changed `SPACEDOCK_BIN` for the next lifecycle without cache invalidation machinery. Require the one response to contain `record`, `validate`, `eligibility`, `consume`, and the record forms/flags `--briefing`, `--result`, `--association`, `--decision`, `--actor`, and `--directive`. If the fingerprint is incomplete, halt before mutation and prescribe refresh or a fresh checkout build selected through `SPACEDOCK_BIN`. Never hand-edit `gates:`. The CLI, not the FO, normalizes retained `--briefing`, `--result`, and `--association` arguments against the launch directory before recorder entry; absolute paths remain accepted.

### 1. Retain and bind the presented package

The FO assembles the same judgment inputs it already owns, but retains their exact presentation package first:

- `ROOM/briefing.json`, with the required basename and the decision question;
- a concise primary gate review stating capability/change, test and evidence, exact reviewed snapshot, material/deferred/polish findings, one FO recommendation, and the concrete decision ask;
- the entity, spec, reports, and other raw material as reference artifacts, not as substitutes for the primary review;
- each artifact identified by URI and exact SHA. Reuse an existing immutable file only when the presentation resolver can reproduce its bytes; otherwise freeze a room copy. Do not duplicate reproducible payloads merely to satisfy the recorder.

After the complete package is retained, invoke:

```text
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --briefing ROOM/briefing.json --workflow-dir WORKFLOW_DIR
```

The record call must exit 0, identify the selected gate/attempt/Briefing, and report `state=open`. Its internal write path has already parsed the canonical package, validated the rebuilt record and transition, and replaced the entity atomically; a second `gate validate` would only reread the same bytes. Anything else halts presentation. In a split-root workflow, `${SPACEDOCK_BIN:-spacedock} state commit SLUG` then durably commits both the folder room and index mutation before the Captain sees the gate.

### 2. Judge and present

Only after the open binding succeeds and is committed does the FO perform `«gate.ac-cross-check»`, make its evidence judgment, and invoke `present-gate`. `present-gate` is the sole render owner. Its canonical concise review has exactly these ordered, nonblank fields: `Capability/change:`, `Test and evidence:`, `Reviewed snapshot:`, `Findings:`, `Recommendation:`, and `Decision ask:`. The reviewed-snapshot field names the bound Briefing identity and digest; the ask names `approve`, `revise`, and `hold`. Entity, spec, room listing, and `briefing.json` remain linked evidence rather than replacing the review with raw material.

For the default chat channel, presentation completes only when that assembled review has been emitted as one root-assistant message after the bind commit and before any chat decision record. The presenter must emit that message before the next gate-mutation tool call; a prior narration of judgment, a child message, a skill-read result, or a final summary after closure does not satisfy the barrier. Direct Captain operation then waits for the semantic call. Under an explicit delegated conn, the FO may proceed automatically only after the same review message is observable; delegation chooses who renders the decision and never waives presentation. The lifecycle cannot invoke `gate record --decision`, `gate consume`, feedback routing, or successor dispatch while chat presentation is incomplete.

A declared presentation override retains its existing blocking room/Result contract. There, provider exit, Result validation, and retention complete presentation before `gate record --room`; no new root-chat message or second record is required. This reset changes only the default-chat render owner and completion boundary, not provider transport, retention, or recorder authority.

### 3. Record the decision

After a decision arrives, record exactly one semantic source:

```text
# Captain personally rendered the chat decision
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold \
  --actor person:captain [--reason REASON] --workflow-dir WORKFLOW_DIR

# FO rendered the delegated decision under an explicit conn
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold \
  --actor agent:first-officer --reason EVIDENCE_JUDGMENT \
  --directive EXACT_QUOTED_CAPTAIN_GRANT --workflow-dir WORKFLOW_DIR

# Exact provider result
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --result RESULT --association ASSOCIATION \
  --actor AUTHORIZED_ACTOR [--adoption-note AUTHORIZER] \
  --workflow-dir WORKFLOW_DIR
```

`revise` and `hold` require a reason (or the provider's included same-Briefing Annotation). A delegated FO approval always supplies both a reason and exact directive even though a direct portable approval may be reasonless. This preserves the identity that actually rendered the decision while distinguishing direct Captain authority from delegated conn.

Ordinary Captain language maps to the existing schema as follows:

| Captain call | Canonical decision | Durable meaning and route |
| --- | --- | --- |
| `approve` | `approve` | Record approval, commit, then consume. |
| `redo with feedback` | `revise` | The reason says the Captain accepts the direction and lists concrete asks; commit, then invoke `«feedback.route»`. |
| `reject` when the stage has `feedback-to` | `revise` | The reason says the Captain rejects the direction and lists concrete asks; commit, then invoke `«feedback.route»`. |
| `reject` without `feedback-to` | `hold` | The reason says the direction was rejected but no correction owner is configured; commit and stop. |
| `hold` or an equivalent pause | `hold` | Record why work is paused; commit and stop. |
| `not yet` | `hold` | Require a nonblank reason naming what remains; commit and stop in the same held state. |

`redo` and routed `reject` deliberately share canonical `revise`, while `Resolution.reason` durably preserves whether the Captain accepted or rejected the direction. No recorder enum or schema field is added.

The close call must exit 0, name the bound attempt/Briefing, report `state=closed`, and reproduce the intended decision. Its internal write path validates the existing record, authority fields, Resolution, derived application, transition, and rebuilt record before atomic replacement; a separate closure `gate validate` adds no authority. Every successful close is then committed through `state commit` before any route—approve, revise, hold, or a consume attempt. This close commit is a durability barrier: it preserves the exact Resolution even if consume later refuses or the route stops. A failed close or failed close commit halts without feedback, advancement, or dispatch.

### 4. Route the closed decision fail-closed

The Resolution is recorded before any route:

| Closed decision / condition | Required FO route |
| --- | --- |
| `approve` | After the close commit, invoke `gate consume` directly. Only exit 0 plus `condition=approved-pending eligible=true consumed=true` and the expected immediate successor is success. Consume itself evaluates the authoritative condition under lock, atomically advances `status`, and marks the application consumed. Commit that mutation. If the target is nonterminal, hand the current stage to ordinary reuse-or-fresh dispatch. If the target is terminal, enter the existing merge ceremony and do not build or dispatch a successor. No separate `status --set` advance is allowed. |
| `revise` | After the close commit, never consume it as an advance. Invoke the existing `«feedback.route»` procedure; this task does not implement advisory correction rounds. `gate eligibility` remains available only when the condition needs diagnosis. |
| `hold` | After the close commit, leave the entity at the gate and surface the reason; never consume, advance, or dispatch. `gate eligibility` remains an optional diagnostic read. |
| approved but blocked, held, unknown, wrong-stage, or otherwise ineligible | The Resolution is already durable in the close commit. Halt advancement and dispatch, preserve status bytes, and name the exact reported condition and missing/current artifact or field. |
| `stale` | The consume attempt reports stale, exits nonzero, leaves status unchanged, and changes the pending application only to `superseded`. Commit that supersession, retain/bind a new Briefing, and re-present; never expose a separate rebind/supersede ceremony to the Captain. |
| already `consumed` | Treat the approval as spent. Do not record another decision or consume again. A nonterminal current status resumes ordinary dispatch/recovery; a terminal current status resumes the existing merge ceremony, including any recorded merge `mod-block`, and has no successor dispatch. A diagnostic repeat consume is nonzero and byte-clean. |

After successful consume, its second `state commit` is the approval-effect durability barrier. No nonterminal successor `dispatch build` or host dispatch may occur until a descendant of the close commit contains `application.state: consumed` and successor status. A terminal successor never enters dispatch: after the same barrier it follows `«merge.guard»` and the registered merge hook. The generic gate lifecycle ends at this fork. Ordinary dispatch owns nonterminal build/reuse/fresh dispatch, merge owns terminal finalization, and each host adapter owns how spawn or completion is observed; this procedure adds no “very next host event,” returned-handle, or common transport requirement.

### 5. Resume and retry without minting duplicates

Use `status --boot --identify --json`, entity state, and the previous command result to choose the next semantic operation. `gate validate` and `gate eligibility` remain on-demand diagnostics for malformed or ambiguous state; neither is a mandatory prelude:

- open attempt + same retained Briefing: repeating `record --briefing` is an idempotent bind; continue presentation;
- open attempt + changed package: `record --briefing` selects the changed binding under the same attempt; the FO speaks only of updating the presented Briefing, not recorder mechanics;
- closed + pending approval: do not re-record the decision; first ensure the exact closure has a durable state commit, then consume, which authoritatively accepts or refuses;
- closed + revise/hold: first ensure the exact closure has a durable state commit, then route or remain held; never call advance-consume;
- consumed + nonterminal status: dispatch/recover the current stage without consuming again;
- consumed + terminal status: resume `«merge.guard»` and the recorded merge hook without consuming or dispatching again;
- stale: the consume failure materializes supersession; bind the replacement Briefing and re-present.

Every nonzero command becomes explicit FO-visible friction with the command, exit, and actionable missing artifact/field/step. The FO never repairs `gates:` by hand or substitutes `status --set`.

## Friction inventory and disposition

| Observed sprint friction | Disposition and coverage |
| --- | --- |
| The FO skill contains no literal `record`/`validate`/`eligibility`/`consume` calls. | In scope: lifecycle procedure at the gated-stage branch. AC-1/AC-2 live and mutation proof. |
| Briefing, association, and Result JSON were manually crafted and error-prone. | In scope only as exact input guidance and retained real fixtures. The recorder continues to derive lifecycle/ids; no new JSON builder. AC-5 checks actionable failures. Better package authoring remains presentation-provider work. |
| It was unclear which entity/spec/reference artifacts to package and whether to use URI+SHA or copied bytes. | Settled rule above: reproducible existing bytes stay URI+SHA; mutable/unresolvable snapshots are frozen in-room; the recorder never copies. AC-4 replays the three-artifact 3k package and hashes every declared revision. |
| `briefing.json` linkage/discoverability was unclear; an accepted alternate basename later dead-ended Result resolution. | In scope guidance requires the canonical basename and a room whose relative references resolve. AC-4/AC-5 include the landed alternate-basename no-mutation control. |
| Bind → decision → consume was not presented as one legible authority lifecycle. | In scope: three mutation phases above and exact trace oracle. `validate`/`eligibility` remain diagnostics rather than mandatory ceremony. AC-1/AC-2. |
| Raw room/entity/spec files were mistaken for validation-gate evidence. | In scope: primary review must state capability, test/evidence, snapshot, findings, recommendation, and ask; raw artifacts are references. AC-6 grades the rendered gate. |
| Delegated Captain conn left actor/reason/directive provenance ambiguous. | In scope: direct `person:captain` versus delegated `agent:first-officer` plus evidence reason and exact quoted directive. AC-3. |
| A successful close could be followed immediately by a refused consume or stopping route without a durable Resolution commit. | In scope: every successful close has a state-commit barrier before approve/revise/hold routing; successful consume adds a descendant commit before dispatch. AC-1/AC-3. |
| Before h1, the FO manually changed status after recording a Resolution. | In scope: replace that step with locked `consume` and prohibit `status --set` bypass. Optional eligibility preview is not authority. AC-1/AC-2 measure zero unconsumed advances. |
| Retries, stale input, and duplicate application behavior lacked an agent path. | In scope resume matrix maps each landed condition to idempotent bind, no re-close, consume once, stale supersession, or dispatch recovery. AC-2/AC-7. No retry daemon or new ids. |
| Failures did not reliably tell the FO which artifact/step was missing. | In scope as halt-and-surface contract plus output assertions. AC-5; parser/record failures are byte-clean, validate/eligibility are read-only, held/blocked/repeat-consume refusals are byte-clean, and stale consume's sole permitted mutation is explicit supersession. |
| Fixtures placed as workflows polluted repo-wide workflow discovery. | In scope test layout: all fixtures live below existing ignored/test-owned `skills/integration/testdata/` or temp dirs, never as discoverable repo workflows. AC-4/AC-8. |
| A real runtime replay, including Pi, was absent. | In scope: one host-neutral shared scenario, codified fixture, Claude/Codex adapters, and Pi live-capable coverage. AC-1/AC-8; host-neutral core changes require the relevant live lanes. |
| Extra folder-form room artifacts were omitted by `state commit` and required manual exact-path commits. | Resolved dependency, not reimplemented: vn/PR #558 makes folder-form `state commit` the durable unit. AC-4 exercises it; no new artifact registry or state verb. |
| Historical gates held canonical decisions, but not every one came through the literal CLI. | In scope: use the real retained 3k validation package and delegated decision as the replay baseline, then require the actual command trace. AC-1 closes the history/procedure gap rather than treating hand-authored YAML as proof. |
| Gate presentation became long and context-dependent. | In scope: one concise primary review, entity/spec as linked references. AC-6. Presentation UI/channel design remains xb-owned. |
| `rebind`/`supersede` vocabulary exposed recorder mechanics to humans. | In scope vocabulary rule: the FO operates a Briefing/gate with existing commands; it speaks of updated, stale, or spent decisions. No new command/subcommand. AC-6. |
| The FO silently worked around command/CLI friction during the sprint. | In scope: every command failure is surfaced and recorded; no manual `gates` edit. AC-5. Product CLI UX defects discovered by the replay are reported as findings, not patched inside skill work. |
| Version-compatible repo-root `./spacedock` lacked the gate subcommands. | In scope: exactly one `gate --help` probe immediately before each lifecycle, plus a fresh-build/refresh route before mutation. AC-5. There is no session cache, executable digest, or launcher-swap laboratory; fresh resolution naturally observes replacement or PATH changes on the next lifecycle. |
| Consumed state was routed unconditionally to dispatch even when its target was terminal. | In scope: terminal targets enter the existing merge ceremony; only nonterminal targets dispatch. The existing `gate-review-presentation-command` state proves terminal+consumed is absent from `status --next`. Its still-open PR also exposes the smallest merge-hook selector gap: startup/idle must include terminal rows only while `mod-block: merge:pr-merge` records an in-flight ceremony. AC-1/AC-7. |
| Pi review extraction concatenated every session and accepted any message role. | In scope evidence repair: inspect the single flat root Pi session and accept the six-field review only from a root `assistant` message between bind and decision. Nested subagent sessions and user/tool text cannot satisfy AC-6/AC-8. |
| A feedback reviewer can recommend `REJECTED` while the First Officer waits for an ordinary Captain gate decision. | In scope: restore the automatic pre-lifecycle `«feedback.route»` branch. Its deterministic structural deletion mutant is the authoritative ownership proof; the goal-only rejection journey is positive integration evidence and need not red under the mutant because a model may infer the bounce. |
| Captain `redo` and `reject` do not name recorder decisions. | In scope: the lifecycle skill owns one explicit translation table into `approve|revise|hold`; `Resolution.reason` preserves whether the Captain accepted or rejected the direction. |
| The live prompt can make a broken skill look correct by naming every gate command and commit barrier. | In scope evidence repair: prompts retain only fixture, authority, goal, and stop-marker constraints. Static controls reject lifecycle commands, ordered procedure, and review labels in prompt text. |
| Two goal-only Codex runs load the lifecycle and presenter, then mutate the decision without a root-visible review. | In scope ownership repair: `present-gate` owns the sole six-field chat template and its emit-before-record completion barrier; the lifecycle waits for that completion. The existing root transcript interval and unchanged live oracle prove the outcome. No prompt label, second harness, or recorder field is added. |
| Claude child output can satisfy the root-review oracle. | In scope evidence repair: accept only Claude assistant rows with an empty top-level `parent_tool_use_id`; a focused child-rejected/root-accepted stream proves the boundary. |
| Fresh `/tmp/spacedock-current gate record first-officer-gate-command-lifecycle --briefing docs/dev/.spacedock-state/first-officer-gate-command-lifecycle/review/ideation/briefing-1/briefing.json --workflow-dir docs/dev` failed before mutation with `resolve briefing room: Rel: can't make ... relative to ...`. | In scope CLI fix: normalize relative `--briefing`, `--result`, and `--association` paths against the launch directory before recorder entry. AC-5 makes the natural relative form succeed and keeps invalid/missing inputs byte-clean; absolute paths remain compatible. |
| The presentation helper could report failure before a late provider Result arrived, and earlier cleanup destroyed result bytes. | Boundary: xb/presentation transport owns exact Result and association retention; this procedure accepts only retained bytes. AC-4's provider arm uses retained inputs. No polling or transport enters this task. |
| Provider and canonical Briefing ids differed even when primary artifact bytes matched. | Boundary: landed recorder verifies complete association and normalizes identity only after exact bytes/authority. The FO supplies Result + association; it does not normalize by hand. AC-4. |
| `status --set` reserialized hand-authored `gates`, breaking anchors. | Prohibited workaround: binary owns all gate writes; `status --set` never substitutes for recorder or consume. AC-2/AC-5. |

## Acceptance criteria

**AC-1 (VALUE) — Every supported gate outcome reaches its durable owner without an unauthorized route.** For an approved gate, the real 3k replay emits the exact ordered authority trace (briefing record, delegated decision/Result record, consume), commits the package before presentation, commits the close before routing, and commits the consumed successor as a descendant. A nonterminal target has exactly one successor `dispatch build` and one durable marker-bearing effect; a terminal target has no dispatchable successor and enters the existing merge ceremony. For a completed feedback gate whose reviewer recommends `REJECTED`, the FO invokes the existing correction flow before Captain presentation, with zero Captain Resolution, consume, or successor dispatch; the existing advisory correction round is its durable record. The measured count of advances or dispatches without the required consumed-application commit is **0**. *Verified by:* deterministic real-CLI nonterminal and terminal controls, an authoritative shared-core branch-deletion mutant, a positive goal-only rejection journey, and one goal-only representative approved journey on each supported host. The structural mutant—not a live mutant—is the ownership proof for automatic rejection because the model can infer the fixture’s bounce trajectory.

**AC-2 — The minimum authority sequence is load-bearing and fails closed.** Omitting briefing record, decision/Result record, or consume makes the scenario grader refuse successor dispatch. Removing the former open validate, closed validate, and eligibility reads from the positive trace does not weaken the recorder/consumer guards: malformed open/close writes fail inside `record`, and stale/blocked/wrong-stage/spent approvals fail inside locked `consume`. Revise, hold, and ineligible controls produce zero advance/dispatch effects; revise invokes feedback routing and neither revise nor hold is consumed as advance. *Verified by:* three existing-fixture controls that each execute the remaining sequence through the freshly built real binary and logging wrapper (not a mutated expected slice), assert actual exits/state/commits, and offer an otherwise-valid dispatch observation that the oracle rejects; malformed-record and consume-refusal fixtures remain direct controls.

**AC-3 — Direct and delegated authority remain durably distinguishable and every successful decision close survives a stopped route.** A direct Captain decision records `by: person:captain`; a delegated FO decision records `by: agent:first-officer`, a nonblank evidence reason, and the exact quoted Captain conn as directive/adoption provenance. Approve, revise, and hold each produce a close commit before consume/feedback/stop; an approval whose consume is later refused retains that exact Resolution commit. Missing delegated fields fail before closure and leave the entity byte-identical; a procedure mutant that records the delegated live decision as `person:captain` fails the provenance grade. *Verified by:* public CLI fixtures, delegated live baseline, actor-swap mutant, and real-Git approve-refused/revise/hold controls asserting commit hash, closed snapshot, unchanged status, and no dispatch.

**AC-4 — The FO package is complete, reproducible, and durable without unnecessary copying.** The replay uses the retained 3k three-artifact validation Briefing: every URI resolves through the declared resolver and every SHA matches; reproducible existing artifacts are not duplicated, while the frozen reviewed snapshot stays in the room. One folder-form `state commit` includes the index and all new room files and excludes dirty sibling paths. Fixtures live only in testdata/temp locations and do not alter repo workflow discovery. *Verified by:* exact digest/path assertions, a real-Git state-commit assertion, and before/after workflow-discovery candidate equality.

**AC-5 — Readiness and retained-input paths work at the CLI boundary while failures preserve landed write boundaries.** Exactly one freshly resolved `gate --help` call immediately before every lifecycle rejects a capability-stale executable before room/entity mutation and names refresh/fresh-build remediation. There is no session cache, target digest, same-path replacement experiment, or PATH-swap experiment: the next lifecycle performs its own single probe against the then-selected launcher. A repo-root-relative `--briefing` binds the same canonical bytes as its absolute form, and relative `--result`/`--association` inputs reach semantic validation rather than failing mixed-root normalization. Missing/noncanonical Briefing and invalid association/actor/provenance name the missing artifact/step, return nonzero, preserve entity bytes, and leave no lock residue where promised. Validate/eligibility diagnostics are byte-clean; hold/blocked/repeat-consume refusals are byte-clean; stale consume changes only pending → superseded and never status. *Verified by:* live and fixture logs that count exactly one successful help probe before each lifecycle's first mutation, stale-launcher pre-mutation refusal, CLI relative/absolute equivalence fixtures, and the existing whole-file/tree-hash refusal matrix.

**AC-6 — The Captain receives one concise root-visible review before any chat decision mutation, and ordinary decisions have one durable meaning.** After the bound-package commit and before `gate record --decision`, the default chat path emits exactly one root-assistant review containing the ordered, nonblank `Capability/change:`, `Test and evidence:`, `Reviewed snapshot:`, `Findings:`, `Recommendation:`, and `Decision ask:` fields. It names the exact bound Briefing identity/digest and all three calls (`approve`, `revise`, `hold`), links entity/spec/Briefing, and does not lead with raw files or recorder ceremony. Prior judgment narration, child output, skill-read output, and post-close summary do not count. Extraction preserves multiplicity: two qualifying root reviews in the bind→decision interval fail grading rather than silently selecting the last. `approve` records `approve`; `redo with feedback` records `revise` with a reason that accepts the direction; `reject` with `feedback-to` records `revise` with a reason that rejects the direction; and `reject` without `feedback-to` or `hold` records `hold`. `not yet` independently maps to `hold`, requires a nonblank reason, and stops. Concrete asks remain in the reason for every correction route. Pi accepts only a root-session assistant review; Claude accepts only assistant rows with empty top-level `parent_tool_use_id`. *Verified by:* the existing live review extractor bracketed by bind commit and decision record, Codex before/after/non-agent controls, a deletion-sensitive bind-commit→two-qualifying-root-reviews→decision-record multiplicity control, presenter-template/barrier deletion mutants paired with the unchanged live outcome, mapping delete/swap mutants, the `not yet` deletion control, routed/held controls, Pi nested/nonassistant controls, and a Claude child-rejected/root-accepted stream.

**AC-7 — Resume is idempotent, one-use, and faithful to the recorded Captain meaning.** Open same-Briefing retry creates no duplicate attempt; a closed pending approval is not re-recorded and is committed before consume; consumed state never consumes twice; stale becomes superseded without advancing. Across three fresh-process approval passes, close commits, successful consumes, and transitions each total exactly 1. A recorded `revise` resumes `«feedback.route»`, and its reason still distinguishes accepted-direction redo from rejected-direction correction; a `hold`, including rejection without `feedback-to` or `not yet`, resumes as held and never consumes or dispatches. The `not yet` snapshot retains its nonblank reason. Consumed nonterminal status resumes dispatch; consumed terminal status resumes the existing merge ceremony. *Verified by:* the existing resume/merge controls plus fresh-process routed-redo, routed-reject, held-reject, and not-yet snapshots asserting exact reason, route, zero duplicate close, and zero unauthorized consume/dispatch.

**AC-8 — The shipped skills, not the live prompt, own one runtime-portable lifecycle.** The recorded-gate and rejection prompts name only the workflow/entity fixture, retained Briefing or delegated connection authority, the outcome goal, and the durable stop marker. They contain no gate subcommands, state-commit sequence, dispatch-build instruction, route invocation, or six-field review labels. `present-gate` owns the sole canonical chat template and emit-before-record completion barrier; `fo-gate-lifecycle` owns the bind/close/consume route and treats presenter completion as its precondition instead of duplicating the template. Claude, Codex, and Pi must nevertheless satisfy the same approved-route oracle: the root-visible review between bind and decision, three ordered mutations, close commit, descendant consumed-successor commit before exactly one dispatch build, and exactly one later durable effect. The deterministic terminal and rejection variants provide zero-dispatch complementary proof without multiplying the host matrix. The logging shim remains transport-neutral; Pi and Claude apply their root-row rules while native adapters retain their own spawn shape. *Verified by:* a prompt-procedure exclusion control, single-owner structural/deletion controls, the three unchanged goal-only live journeys, shared terminal and rejection controls, and existing host-native dispatch fixtures.

## Minimum replay and test plan

The test package copies the real retained 3k validation Briefing and its declared artifacts into the existing temp folder-form workflow. Its ordinary form keeps the gated `validation` stage's supported nonterminal successor; a table arm changes only that same fixture's immediate target to the declared terminal stage. It preserves the package bytes and historical delegated authority. The live prompt supplies only the workflow/entity, retained Briefing, exact connection grant, outcome goal, and stop marker; the shipped First Officer skills must discover and execute the procedure. The existing command-logging wrapper delegates to the freshly built real binary and records exit/stdout/stderr, state HEAD hashes, and entity snapshot digests after each `state commit` and immediately before `dispatch build`.

1. **Goal-only baseline and terminal fork (AC-1/AC-3/AC-6/AC-8; high):** launch the real FO with only fixture, delegated authority, goal, and stop-marker constraints. Grade exactly one concise root-visible review after the bind commit and before the decision-record command, then assert three ordered gate mutations. The bind commit snapshot is open; the close commit contains the exact Resolution/pending application; the consumed commit is its descendant and contains advanced+consumed state; dispatch's pre-HEAD equals or descends from consumed; one later commit adds the marker-bearing successor report. Run the corrected unchanged host-neutral journey on Codex first, then Claude and Pi. The deterministic terminal arm retains zero dispatch and the existing merge route. Static controls fail if either live prompt contains gate commands, state-commit ordering, dispatch construction, feedback-route instructions, or the six review-field labels. Table-driven Codex transcript controls reject a qualifying review before bind, after decision, in a non-agent row, or twice as root agent messages between bind and decision. The duplicate fixture expects an explicit count-two failure; deleting either input yields the exact-one positive, while deleting the multiplicity rejection makes the duplicate control red rather than silently accepting the last message.
2. **Minimum-sequence proof (AC-2; medium):** for each omission, use the existing fixture/wrapper and freshly built real binary to invoke every remaining authority command in order, record its actual exit/state/commit trace, and offer a prospective dispatch observation. Omitted bind makes close/consume fail, omitted close makes consume fail ineligible, and omitted consume leaves pending state; every arm denies dispatch. Do not satisfy this with expected-slice deletion. Keep the positive free of mandatory validate/eligibility and retain direct malformed-record/consume-refusal fixtures.
3. **Rejection branch, decision translation, and close durability (AC-1/AC-2/AC-3/AC-6/AC-7; high):** run the existing rejection fixture with a goal-only prompt as positive integration proof. A completed reviewer `REJECTED` recommendation must enter `«feedback.route»` before Captain presentation with zero gate close, consume, or successor dispatch. The existing structural contract test deletes the shared-core pre-branch and is the authoritative red control; the live journey is not a mutation oracle and may still bounce by model inference. Table-test Captain approve, redo, reject-with-target, reject-without-target, hold, and `not yet`; deleting or swapping each mapping must fail, including deletion of only the `not yet` alias. Real routed-redo/routed-reject/held-reject/not-yet snapshots assert canonical decision, nonblank accepted/rejected/pause reason, close commit, held resume state where applicable, and zero unauthorized consume/dispatch.
4. **Provenance, package, and path matrix (AC-3/AC-4/AC-5; medium):** direct versus delegated versus exact Result inputs; complete versus truncated association; canonical versus alternate basename; relative versus absolute `--briefing`, `--result`, and `--association`; reproducible URI+SHA versus frozen copy. Compare bytes and lock residue on every refusal, and prove relative/absolute forms bind or adopt identical bytes. Add the narrow normalization assertion in existing `internal/cli/gate_test.go`; do not change recorder path semantics.
5. **Resume and merge passes (AC-7; medium):** rerun fresh processes over open, closed-but-uncommitted, pending-committed, stale, consumed-nonterminal, and consumed-terminal snapshots without mandatory diagnostic reads; exactly one close commit, transition, and consume across three passes. The terminal arm derives from the existing lifecycle fixture and proves zero dispatchable successors. Extend the existing merge fixture so terminal+`mod-block`+merged sentinel finalizes idempotently, and pin both `pr-merge` templates to scan terminal PR rows only when that exact in-flight block remains. Add no hook runner or second terminal lifecycle.
6. **Per-lifecycle capability probe (AC-5; low):** delete the identity-cache test and helper. The existing command log must show exactly one successful `gate --help` immediately before the first mutation in each fixture/live lifecycle. Keep the stale shim only as the direct pre-mutation failure control. Do not add same-path replacement, PATH-swap, digest, or cache-reuse experiments.
7. **Runtime lanes and root authority (AC-6/AC-8; high):** keep the goal-only host-neutral scenario once with Claude, Codex, and Pi bindings. Each live lane retains command logs, one root-visible gate review between bind commit and decision record, before/after entity, state Git history, and one durable successor effect. Extractors return all qualifying interval reviews or an explicit count/error; none may overwrite an earlier match with the last. Pi accepts only a root assistant row in that interval. Claude decodes the existing top-level `parent_tool_use_id` and ignores every nonempty-parent row; a focused stream with a qualifying child review and qualifying root review must select only the root. Codex accepts only `agent_message` rows in the same interval, and its two-root control must fail exact-one grading. Host-native fixtures separately prove transport shape; no exact-child tracing or common spawn schema is added.
8. **Repository gates:** `gofmt -w ./cmd ./internal`, focused scenario/skill tests, `go test ./...`, `go test ./... -race`, live-tag compilation, and the required live lanes.

The primary independent live oracle is durable state plus the command/dispatch log and state Git history. It resolves the logged close and consumed hashes, verifies their entity snapshots and ancestry, and compares dispatch's pre-HEAD; substring order alone cannot pass. Deterministic skill-contract mutants independently prove the successful-close barrier, automatic reviewer-rejection branch, Captain-language mapping, prompt ownership, and Claude root boundary are load-bearing.

Reuse the current `recorded_gate_lifecycle_test.go`, `shared_fixtures_test.go`, command wrapper, 3k fixture, rejection fixture, dispatch/state-Git oracles, and `TestFOGateLifecycleOwnsEveryEngagedEntry`. Preserve the repo-local fixture Git identity fix. Narrow the existing prompts rather than adding a scenario, and decode Claude’s already-present parent marker in the existing extractor rather than adding transport machinery. Do not add a harness, live matrix, runtime protocol, compatibility path, exact-child/artifact attribution, recorder behavior, schema field, or product-side decision mapper.

**Hard reset trigger:** deterministic controls must be green before live spend. Because Codex has already omitted this exact artifact twice, the corrected presenter owner gets one unchanged goal-only Codex run first. If it again reaches chat decision record without the qualifying root review, stop immediately; do not add prompt coaching, another prose synonym, a second harness/controller, or an oracle exception. If Codex passes, each remaining supported host gets one diagnosed run and at most one unchanged rerun for a different missing obligation. A goal-only rejection run that still bounces after the shared-core branch is deliberately deleted does not trigger reset or demand another prompt/controller/harness; the structural mutant already owns that attribution.

## Obligation delta

| Obligation | Authority | Bearer | Proof burden |
| --- | --- | --- | --- |
| Three ordered authority mutations; one consumed approval; zero unauthorized advances/dispatches | Canonical AC-1/AC-2 and `fo-gate-lifecycle` | First Officer lifecycle contract plus landed recorder/consumer guards | Real CLI fixture, three real-binary omitted-mutation controls, malformed-write/consume-refusal fixtures, and each host's representative live command/state trace |
| Exact package committed before presentation; every successful close and consumed successor committed before routing | Canonical AC-1/AC-3 | First Officer write contract | Logged commit hashes/snapshot digests and ordered state Git ancestry showing all durability barriers; no earlier nonterminal dispatch or terminal merge continuation |
| Reviewer `REJECTED` automatically enters the existing correction owner before Captain lifecycle | Canonical AC-1 | Shared First Officer core and existing feedback-rejection flow | Authoritative structural pre-branch deletion mutant; positive goal-only rejection journey with zero Resolution, consume, or successor dispatch is integration evidence, not a live mutant oracle |
| Captain words have one canonical durable meaning | Canonical AC-6/AC-7 | Existing deferred gate lifecycle | Mapping delete/swap controls, including deletion of only `not yet`, and real routed-redo, routed-reject, held-reject, and not-yet snapshots preserving reason and route |
| Route by consumed target kind | Canonical AC-1/AC-7 | Existing ordinary dispatch loop for nonterminal; existing merge ceremony for terminal | Nonterminal has one post-commit `dispatch build` and durable effect; terminal has empty dispatchability, zero builds/effects, and merge-guard/hook continuation |
| Every engaged route loads the lifecycle before gate action | Deferred topology contract | Shared First Officer core | `TestFOGateLifecycleOwnsEveryEngagedEntry` and native load-before-action fixtures; not every route live on every host |
| Host transport shape | Canonical AC-8 and each supported runtime adapter | Claude, Codex, and Pi adapter/fixture owners | Host-native fixtures, including Codex's supported build → wait → durable-report evidence when public spawn records are absent |
| Gate capability check | Canonical AC-5 | First Officer lifecycle contract | Exactly one freshly resolved `gate --help` immediately before each lifecycle; no cache/digest/swap machinery |
| Exactly one root-authoritative review before chat mutation | Canonical AC-6/AC-8 | `present-gate` for render/completion; existing host transcript extractors for observation | One six-field root message after bind commit and before decision record; bind→two-root-reviews→decision must expose count 2 and fail, with deletion of multiplicity rejection caught; Codex interval controls, Pi root-session/assistant controls, Claude empty-parent control, and unchanged live oracle |
| Shipped-skill procedure ownership | Canonical AC-1/AC-8 | Shared First Officer core plus deferred lifecycle for route; `present-gate` for the sole review template and chat completion barrier | Goal-only prompts, prompt-procedure exclusion control, single-owner/barrier deletion controls, and unchanged three-host oracle |
| Terminal open-PR resumption | Existing merge ceremony plus AC-7 | Shipped and development `pr-merge` startup/idle hooks | Selector includes terminal PR-bearing rows only with `mod-block: merge:pr-merge`; existing sentinel and merge guard finalize, while OPEN/CLOSED behavior stays unchanged |

Captain rulings of 2026-07-24 still remove mandatory diagnostic reads, repeated probes, executable caching, absolute-path caller ritual, runtime-specific next-event/handle language, exact-child/built-artifact forensics, and a route × host matrix. Cycle 17 added four obligations proven missing by validation: automatic reviewer rejection routing, Captain-language translation through the existing schema, goal-only prompt ownership, and Claude root-row filtering. Cycle 20 adds only the presenter ownership/completion obligation proven missing twice by Codex. It removes no proven authority mutation, provenance distinction, one-use guard, routing fork, retained room, or fixture identity. This amended canonical AC block, test plan, obligation delta, surface, and documentation diff are the sole active proof authority; historical reports do not compete with them.

## Expected surface and tolerance

The failed Cycle-19 tip `37d6980b` is intentionally local and off PR #565. Against the approved baseline `13d702492131df17dd3ac87245d6d773f4df959b`, it changes exactly seven files at **+95/-14**; shared core is 26,737 bytes and `fo-gate-lifecycle` is 6,369 bytes. Cycle 20 preserves that evidence tip and resets the implementation boundary to one incremental ownership repair in an eighth existing file:

| File | Existing Cycle-19 delta | Intended Cycle-20 incremental delta |
| --- | ---: | ---: |
| `skills/first-officer/references/first-officer-shared-core.md` | +1 / -1 | 0 / 0 |
| `skills/fo-gate-lifecycle/SKILL.md` | +3 / -1 | +3 / -2 |
| `skills/present-gate/SKILL.md` | 0 / 0 | +16 / -14 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +60 / -5 | +22 / -0 |
| `internal/ensigncycle/shared_fixtures_test.go` | +3 / -4 | 0 / 0 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +17 / -0 | +16 / -0 |
| `docs/site/concepts/gates-and-decisions.md` | +9 / -2 | +8 / -15 |
| `docs/specs/scenario-testing-principles.md` | +2 / -1 | 0 / 0 |

The incremental intent is **+65/-31** from failed tip, for a total intended **+160/-45** against `13d70249`. Additions have a hard stop at **+175** from that baseline. The correction should be net-small in the presenter/docs template because it replaces the conflicting spine rather than adding a second one. No ninth file is permitted without re-ideation.

- **Presenter contract:** `skills/present-gate/SKILL.md` becomes the sole owner of the exact six-field canonical chat template and the root-message-before-mutation completion barrier. Its override/provider rules remain unchanged.
- **Lifecycle contract:** `skills/fo-gate-lifecycle/SKILL.md` removes its duplicate field list, invokes presenter completion as a precondition, and preserves every bind/close/consume/route barrier.
- **Evidence:** the existing contract test makes presenter ownership and lifecycle ordering deletion-sensitive; the existing Codex extractor gains before/after/non-agent and duplicate-root interval controls, preserving all matches or exposing multiplicity instead of returning the last. The goal-only prompts and live oracle are unchanged.
- **Documentation:** replace the conflicting example spine with the same concise six-field review and state that delegated conn does not suppress chat presentation.
- **Component budgets:** shared first-officer core remains **≤26,754 bytes**; deferred lifecycle remains **≤6,600 bytes**; presenter remains **≤8,900 bytes** (currently 8,616).
- **Stop conditions:** any schema/recorder/application/CLI production change, product-side decision mapper, new fixture/harness/runtime protocol, parallel gate lifecycle, per-host lifecycle copy, compatibility layer, exact-child transport forensics, ninth file, addition cap breach, or repeated Codex omission returns to ideation before editing.

### Cycle-20 mechanism evidence

No additional throwaway runtime spike is needed. The retained first Codex JSONL proves all host mechanisms the correction relies on: root `agent_message` events occur between tool calls; the bind and state commit complete at items 14-15; a root judgment narration occurs at item 19; the shipped presenter is read at item 20; and the decision mutation begins at item 21. The existing extractor already isolates root `agent_message` content between bind and decision. The second unchanged run proves that adding labels to the caller does not resolve ownership. What remains is the shipped presenter correction itself, so implementation exercises the riskiest claim first with focused structural/interval controls followed by the single unchanged Codex live run before any other live host or advisory spend.

The following durable excerpt preserves the machine-independent event sequence from that retained JSONL; paths and runtime-generated ids are intentionally elided:

```text
14 command: spacedock gate record recorded-gate-task --briefing <briefing.json> --workflow-dir <workflow>
15 command: spacedock state commit recorded-gate-task
19 root agent_message: "The structured scan confirms ... no material findings. I recommend approval ... and am applying the gate presentation rules before recording it."
20 command: read <installed-plugin>/skills/present-gate/SKILL.md
21 command: spacedock gate record recorded-gate-task --decision approve --actor agent:first-officer ...
root agent_message with all six fields between events 15 and 21: 0
oracle: gate review field "capability/change:" count=0, want 1
```

This excerpt is evidence of the missing artifact, not an alternate fixture or replay source; the retained JSONL remains the full audit record.

## Documentation change proposal

Replace the “What you see at a gate” prose and example in `docs/site/concepts/gates-and-decisions.md`:

```diff
-A gate review has a fixed spine: the first three lines and the last line carry the decision; everything between is supporting evidence. If you stop reading after line three, you can still vote.
+A chat gate review has one concise evidence spine. The First Officer emits it before
+recording either your decision or a decision made under delegated conn:
@@
 ```text
-Gate review: Fix the flaky login test — review
-Chosen direction: replace sleep-based waits with event polling
-Recommend reject: the AC-2 retry scenario has no covering test.
-
-Checklist (from ## Stage Report in docs/ship-features/fix-the-flaky-login-test.md):
-- DONE: login test stable across 50 consecutive runs
-- FAILED: retry scenario unproven — no test exercises it
-
-Reviewer findings
-  Material: AC-2 cites a test file that does not exist
-  Polish:   stage report wording drifts from the template
-
-Assessment: 1 done, 0 skipped, 1 failed.
-
-Decision: approve to close; reject to bounce back to implementation.
+Capability/change: replace sleep-based waits with event polling.
+Test and evidence: login is stable across 50 runs; AC-2 retry has no covering test.
+Reviewed snapshot: Briefing `...` at digest `sha256:...`.
+Findings: material — AC-2 cites a test file that does not exist.
+Recommendation: revise to add the retry scenario.
+Decision ask: approve to close, revise to bounce back, or hold at review.
 ```
@@
-Material findings are the ones that should move your vote; Polish never blocks. The Decision line tells you concretely what your vote does. Every acceptance criterion is cross-checked before the review reaches you; a criterion without cited evidence is named rather than passed over.
+Material findings are the ones that should move your vote; polish never blocks. The
+decision ask tells you concretely what each call does. Every acceptance criterion is
+cross-checked before the review reaches you; delegated authority does not hide the review.
```

In “The three calls” in `docs/site/concepts/gates-and-decisions.md`:

```diff
- **Redo with feedback.** You accept the direction but send concrete fixes back.
+ **Redo with feedback.** You accept the direction but send concrete fixes back. The
+ recorded decision is `revise`, and its reason says the direction is accepted.
- **Reject.** The work bounces back to the stage that owns the fix, carrying your findings.
+ **Reject.** With a configured feedback target, the recorded decision is `revise`, its
+ reason says the direction is rejected, and the work bounces to that owner. Without a
+ feedback target, the decision is `hold` and the First Officer stops for routing
+ help.
 Redo and reject differ only in whether you accept the direction; both carry concrete asks.
+Captain words translate into the existing `approve`, `revise`, and `hold` record; durable
+reasons distinguish redo from routed reject. Automatic bounce applies only when a reviewer
+recommends `REJECTED` at a configured feedback gate; Captain reject without `feedback-to` records `hold` and stops.
```

In the `recorded-gate-lifecycle` seed entry in `docs/specs/scenario-testing-principles.md`:

```diff
- `recorded-gate-lifecycle` — the FO performs three authority mutations...
+ `recorded-gate-lifecycle` — from a prompt naming only fixture, delegated authority, goal,
+ and stop condition, the FO discovers the shipped three-mutation lifecycle before dispatch.
```

## Dependencies, collision, and non-goals

- The landed recorder/consumer, folder-room durability, merge ceremony, and advisory feedback route remain binding inputs. This task changes only how the shared First Officer selects those existing owners.
- `skills/feedback-rejection-flow/SKILL.md`, recorder/schema surfaces, host adapters, and merge hooks are unchanged. The shared core retains the automatic rejection trigger; the deferred lifecycle owns Captain-word translation and durable routing; `present-gate` owns chat rendering and its completion boundary.
- No schema or recorder/application logic; no product decision mapper; no provider transport redesign; no advisory-round implementation; no new gate verb; no second lifecycle/controller; no per-host lifecycle copies; no compatibility behavior; no exact-child transport forensics.

## Minimum value demonstration seed

From a goal-only prompt, one fixture-backed First Officer journey must package and commit an exact validation Briefing, emit the six-field root chat review, record and commit an evidence-bearing delegated approval, consume and commit the application, and only then enter ordinary dispatch. The same existing fixture controls terminal zero-dispatch, automatic reviewer rejection, routed redo, routed reject, held reject, and `not yet` hold/resume. The log binds every durability barrier to a Git hash and entity snapshot, while the transcript orders the review between bind commit and decision record. Deleting the shipped presenter barrier, close barrier, automatic rejection branch, decision translation, prompt exclusion, or Claude root filter must make its focused control red.

## Boundary

This task owns the shared First Officer routing trigger, the deferred lifecycle’s Captain-word translation and durable route, the presenter’s canonical chat rendering/completion boundary, and their behavioral proof. It reuses the landed recorder/consumer, retained package, advisory correction flow, merge ceremony, host adapters, and repo-local fixture identity. It does not add a recorder schema or enum, product decision mapper, duplicate gate judgment, provider change, host-specific lifecycle, compatibility behavior, new harness, or advisory correction-round implementation.

## Stage Report: ideation

- DONE: Map the current FO gate path to the demonstrated 3k/h1 lifecycle without redesigning judgment, presentation, recorder, or application semantics.
  The body fixes the exact bind/open-validate/present/record/closed-validate/eligibility/consume order, commit boundaries, direct/delegated/provider inputs, fail-closed routes, resume behavior, and the pre-h1 manual-transition replacement.
- DONE: Inventory all observed sprint friction and give each item an acceptance-covered correction or named owner/boundary.
  The table covers missing FO calls, manual package/result work, artifact linkage/copy rules, presentation legibility, provenance, stale/retry/duplicates, byte-clean failures, fixture discovery, Pi/live replay, folder artifacts, historical non-command records, vocabulary, CLI friction, provider retention/id mapping, status-set reserialization, the capability-stale 0.26 launcher, and the fresh-binary relative-Briefing-path failure. Exact replay evidence: `/tmp/spacedock-current ... --briefing docs/dev/.spacedock-state/.../briefing.json --workflow-dir docs/dev` failed before mutation with `resolve briefing room: Rel: can't make ... relative to .../first-officer-gate-command-lifecycle`; the FO contract now requires absolute retained-input paths, while CLI normalization/remediation remains a named product follow-up.
- DONE: Define falsifiable acceptance criteria and a minimum real-package live proof in which every recorder call is load-bearing.
  AC-1 measures zero advances/dispatches without a consumed authorization; the real retained 3k validation package drives the six-event live trace; skipped-step, revise, hold, ineligible, stale, duplicate, raw-dump, and capability-stale mutants make failure observable in event logs and durable state.
- DONE: Declare implementation surface, tolerance, documentation diff, dependency boundaries, and concurrent advisory-round collision handling.
  The design changes one FO core branch plus test/live/docs surfaces, no production Go; it serializes the shared command-reference edit and rejects any overlap that would put round semantics or ordinary gate integration into the wrong task.

### Summary

Codified the successful 3k/h1 dogfood procedure as the proposed normal FO gate lifecycle: capability-check the selected binary, retain and bind canonical `briefing.json`, validate open, judge/present, record and validate the exact direct/delegated/provider decision, then use eligibility and one-use consume as the sole approval advance before ordinary dispatch. The ideation turns the sprint's full friction record into acceptance coverage or explicit sibling boundaries and defines a real 3k-package live replay plus skipped-step mutants that fail whenever an FO bypasses a required lifecycle call.

## Historical implementation intended-change declaration (cycle 1)

Declared before code edits on 2026-07-23. Production Go estimate: **0 LOC**. First-Officer contract estimate: **70-100 lines** (hard stop above 120). Test/live estimate: **650-900 LOC**. Documentation estimate: **18-30 lines**. This remains within the approved 2x tolerance and introduces no new mechanism.

- `skills/first-officer/references/first-officer-shared-core.md` — 70-100 contract lines at the existing gated-stage branch for capability preflight, canonical binding, open/closed validation, provenance, eligibility/consume routing, fail-closed behavior, and resume/idempotency.
- `internal/ensigncycle/recorded_gate_lifecycle_test.go` — 500-650 test LOC for the retained 3k package fixture copied into temp roots, real CLI six-event replay, command/dispatch oracle, skipped-step/provenance/raw-review mutants, failure matrix, folder-form state commit, resume, and discovery controls.
- `internal/ensigncycle/shared_scenarios_test.go` — 6-10 test LOC for the host-neutral scenario definition.
- `internal/ensigncycle/claude_live_runner_test.go` — 20-35 live-test LOC for the Claude binding and retained evidence capture.
- `internal/ensigncycle/codex_live_runner_test.go` — 20-35 live-test LOC for the Codex binding and retained evidence capture.
- `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` — 80-120 live-test LOC for explicit Pi live-capable execution using the same fixture, prompt, and durable oracle.
- `internal/ensigncycle/pi_shared_coverage_test.go` — 4-8 live-test LOC marking the shared scenario's concrete Pi lane.
- `docs/specs/scenario-testing-principles.md` — 1-3 docs lines to keep the machine-locked shared scenario list synchronized.
- `docs/site/concepts/gates-and-decisions.md` — 8-12 docs lines explaining the ordinary FO lifecycle and concise review.
- `docs/site/reference/command-reference.md` — 6-10 docs lines for capability readiness; collision check was clear at declaration time (clean worktree, latest file commit `f06cce04`) and will be repeated immediately before editing.

If implementation needs production Go, schema/command changes, more than 120 FO-contract lines, provider transport/polling, a presenter redesign, advisory-round logic, or any file outside this declaration, implementation stops and reports drift for a design reset.

### Intended-change amendment 1

The first shared-scenario compile correctly failed because `internal/ensigncycle/shared_scenarios_meta_test.go` pins the exact scenario list. Add that existing test lock at **1 test LOC**. This is a mechanical parity update within the declared host-neutral scenario proof: production/schema/command LOC remain zero, FO contract remains 58 added lines, and no lifecycle mechanism or AC changes. No design-reset trigger fired.

### Roborev triage (job 649)

- **Dispatch proof — Material plus Needs decision.** Released workflow: an approved gate must start exactly one successor worker. User-visible harm: a First Officer could write a marker or narrate a spawn while no worker exists. Acceptance/boundary: AC-1, AC-7, and AC-8 require observed dispatch and a returned handle after consume. Trigger evidence: every completed Codex gate replay reached `dispatch build`, then emitted a wait with no receiver and no observed spawn. Tighten the oracle to require a runtime spawn event, non-empty handle, and worker-attributed durable output; escalate the repeated host behavior rather than weaken the AC.
- **Exact record/review oracle — Material.** Released workflow: the gate must bind the retained package, preserve the exact delegated directive/provenance, and present a concise evidence review. User-visible harm: wrong briefing, attempt, actor, directive, or an unpresented review could pass. Acceptance/boundary: AC-3, AC-4, and AC-6. Trigger evidence: the common grader currently checks generic substrings and receives a static review fixture. Parse exact durable identities/cardinality and grade the review actually observed from the runtime.
- **Claude successful-event attribution — Material.** Released workflow: Claude is a supported host for the shared six-event journey. User-visible harm: echoed narration or a failed tool call could be counted as a successful lifecycle event. Acceptance/boundary: AC-1 and AC-8. Trigger evidence: the current stream regex is unrestricted by tool-result success. Pair Bash tool-use IDs with successful results before admitting an event.
- **Help probes in event order — Material.** Released workflow: capability preflight precedes all mutation. User-visible harm: required `gate ... --help` probes could be miscounted as lifecycle operations and mask a skipped recorder call. Acceptance/boundary: AC-1, AC-5, and AC-8. Trigger evidence: the command-log parser currently accepts any line containing a gate event. Exclude help invocations and begin the lifecycle trace at the first successful briefing record.
- **Stale-launcher no-mutation proof — Material.** Released workflow: an installed version that lacks the full gate surface must halt before mutation. User-visible harm: a partially compatible launcher could mutate workflow state before the FO discovers a missing subcommand/form. Acceptance/boundary: AC-5. Trigger evidence: the current probe checks aggregate help through test code and does not prove every required semantic form or workflow-byte preservation. Exercise every required help form through a logging shim and assert only read-only probes occurred and workflow bytes remain identical.

### Roborev follow-up triage (job 663)

- **Worker attribution — Material plus Needs decision.** Released workflow: one consumed approval starts one successor worker. Harm: the parent or an unrelated writer could append the marker after a no-op spawn. AC-1/AC-7/AC-8 require worker-attributed durable output. Trigger: the live proof currently correlates spawn/handle/wait with the final marker but cannot tag the filesystem write to a worker thread. Add pre-dispatch absence and the strongest structured result correlation available; retain the host-evidence escalation rather than claim full attribution.
- **Pi ordering — Material plus Needs decision.** Released workflow: Pi must consume before spawn. Harm: independently valid command and session traces could hide an early dispatch. AC-1/AC-8 require one ordered observation. Trigger: the Pi command log and session are separate clocks and Pi auth blocks a fresh trace. Require consume evidence in the parent session before accepting `subagent`; the lane remains unverified until authenticated replay.
- **Skipped-step mutant — Material.** Released workflow: every one of six calls is load-bearing. Harm: deleting dispatch proof as well as the event makes the mutant tautological. AC-2. Trigger: every table arm zeroes `dispatch`. Retain valid dispatch proof and remove only the selected event.
- **Actual stale-FO entry — Material plus Needs decision.** Released workflow: the real FO must halt against a capability-stale launcher. Harm: declarative prose could omit the probe while the helper test remains green. AC-5. Trigger: the FO has no deterministic executable entry point; all actual entries are model-hosted and the available Claude/Pi auth is invalid. Preserve the byte-clean full-form helper proof and contract mutant, but report the unavailable model-host proof rather than substitute a tautological interpreter.
- **Hardcoded 0.26 — Material.** Released workflow: compatible source and later launchers use capability, not a frozen version. Harm: a routine version bump falsely blocks valid gate work. AC-5/documented boundary. Trigger: every later binary. Remove the version predicate and test only required command semantics.
- **Relative-path wording — Polish.** Released workflow still fails byte-clean and gives a remedy; harm is only test brittleness if Go wording improves. AC-4/AC-5 are protected by exit, byte, and lock assertions. Trigger requires an upstream error-message improvement. Do not spend the material-fix budget on this advisory.
- **Idempotent-repeat oracle — False positive.** The exact-six oracle grades the fresh baseline only; AC-7 repeats are graded separately by the three-pass resume test. No released workflow rejects a repeat. Trigger would require feeding a resume trace to the baseline grader, which no runner does. No change.
- **Empty open decision — Material.** Released workflow must not present an already-decided attempt as open. Harm: `decision=approve` satisfies the current substring. AC-1/AC-3. Trigger: any malformed open response with a nonempty decision. Require an exact output field line.

## Stage Report: implementation

- DONE: Publish the intended-change declaration before editing and remain inside the no-production-Go boundary.
  State commits `2a1ffd49` and `5c45f508` name every touched surface and the one mechanical amendment; production/schema/command LOC remain zero.
  Collision checks kept the shared command-reference edit serialized at base commit `f06cce04`.
- DONE: Add behavioral coverage first for the six-event lifecycle and a load-bearing skipped-step mutant.
  The focused suite grades exact ordered successful events, keeps valid dispatch proof while deleting each event, and reds the shipped-skill eligibility deletion.
  Initial RED was the shipped FO contract missing the briefing-record integration point; the focused suite is green on the checkpoint.
- FAILED: Land the ordinary gated-stage FO contract on an architecture that satisfies repository prompt-load invariants.
  The contract implements capability/absolute-path/package/provenance/validation/eligibility/consume/resume routing, but makes the boot-resident shared core 32,289 bytes versus the hard <26,755-byte invariant. Base is 26,092 bytes, leaving only 663 bytes; the lifecycle adds 6,197 bytes. A gate-triggered deferred reference or deliberate budget reset is necessary, and either is outside the declared file/boundary.
  Repository precedent keeps substantially larger procedure cores deferred; this is root-cause evidence for a topology reset, not permission to raise a ratchet.
- FAILED: Complete every deterministic control in the approved matrix.
  Fresh real-CLI bind/validate/decision/consume, folder commits, exact durable identity/digest/provenance, relative refusal/absolute success, revise/hold, resume-one-consume, capability-stale full-form probes, concise presentation, and structured host parsers pass. An actual model-host stale-launcher entry remains unavailable; worker-attributed Codex output and the full stale/ineligible live controls remain unproven.
  Claude counts only successful tool-result-paired Bash events; Pi requires successful consume before subagent and a marker-bearing worker result; Codex requires spawn and correlated wait handles.
- FAILED: Prove the host-neutral scenario across Claude, Codex, and Pi live lanes.
  Claude failed before workflow work on invalid 401 auth. Pi failed before workflow work on reused-refresh-token/no-API-key auth after an isolated compatible pi-subagents package was supplied. Repeated Codex runs completed all six gate commands and durable consume, then emitted wait with no spawn handle or worker output; strict grading correctly fails.
  Claude artifacts: `/tmp/spacedock-fo-gate-live-claude/claude-shared-scenarios/recorded-gate-lifecycle`.
  Pi artifacts: `/tmp/spacedock-fo-gate-live-pi/pi-recorded-gate-lifecycle/run`.
  Codex artifacts: `/tmp/spacedock-fo-gate-live`, `-rerun`, `-final`, `-spawn`, and `-sequenced` roots.
- DONE: Preserve the approved product boundary.
  No production Go, recorder/schema/application, command, provider transport, package generator, presenter, advisory-round, hand-authored-gates, or status-set advancement changes exist.
- FAILED: Pass the final repository and runtime gates.
  `gofmt -w ./cmd ./internal`, focused lifecycle tests, workflow-discovery controls, and `go test -tags live ./internal/ensigncycle -run '^$' -count=1` pass. Final `go test ./...` and `go test ./... -race` fail only `TestFOHostPromptLoadRatchet` (Claude 102278>96081, Codex 81493>75296, Pi 77623>71426) and `TestStartupRecipeCollapsedAndLeaner` (32289>=26755). Runtime lanes fail as recorded above.
- DONE: Request Roborev and triage every finding before edits.
  Jobs 649 and 663 completed; both triage records above carry workflow, harm, AC/boundary, trigger evidence, and classification. Material oracle defects were fixed; worker attribution, actual stale-FO entry, and host behavior are explicit Needs-decision items; polish/false-positive findings were declined with triggers.
- FAILED: Commit a self-contained implementation and declare validation readiness.
  WIP counterexample commit `cabdef33` preserves the complete checkpoint without claiming it is mergeable: its resident-core layout violates a binding design-reset trigger. Actual checkpoint delta: FO contract +58/-3 lines; test/live +1,157 lines; docs +7 lines; production Go 0. Do not advance to validation; the First Officer disposition is REVISE for bounded topology re-ideation.

### Independent audit additions

- **Recommended topology:** move the lifecycle behind a deferred, non-user-invocable `fo-gate-lifecycle` skill and retain only its trigger in the shared core. The boot-core ceiling remains unchanged; do not rebaseline or compress around it.
- **Material bypass:** the headless/no-conn path currently invokes `present-gate` directly, bypassing the recorded lifecycle. Re-ideation must route this path through the same lifecycle before presentation or decision handling.
- **Incomplete behavioral proof:** the checkpoint lacks a shipped-skill live mutant, actor-swap and raw-dump mutants, the full AC-5 refusal matrix, the full AC-7 resume matrix, and exact before/after workflow-discovery equality. Existing focused coverage must not be presented as satisfying those missing arms.
- **Runtime disposition:** Codex repeatedly consumes then waits without an observed spawn, so no-spawn is an implementation blocker. Claude and Pi fail before workflow work on credentials, so their auth failures are external validation conditions, not implementation failures.
- **Bounded next step:** topology re-ideation owns only the deferred skill boundary, the headless bypass, and the named proof gaps. This WIP commit is evidence/counterexample input; no implementation continuation or contract compression is authorized from this stage.

### Summary

Implemented and hardened the six-event integration and its deterministic/live oracles, then stopped at the binding design reset exposed by the final prompt-load gates. The behavioral work is preserved in the assigned worktree, all advisory findings are triaged, and no acceptance criterion or spawn evidence was weakened to manufacture a pass.

### Feedback Cycles

- Cycle 1: REVISE — independent topology audit; surface shared core +6,197 bytes vs available headroom 663 bytes (935%); AC unchanged
- Cycle 2: CHANGES REQUESTED — Roborev job 700; surface 7 test/live files and 584 added LOC vs estimate 7 files and at most 615 added LOC (95%); AC unchanged
- Cycle 3: CHANGES REQUESTED — Roborev job 708; surface 7 test/live files and 708 added LOC vs estimate 7 files and at most 615 added LOC (115%, below 2× tolerance); AC unchanged
- Cycle 4: NEEDS DECISION — Roborev job 711; surface 7 test/live files and 741 added LOC vs estimate 7 files and at most 615 added LOC (120%, below 2× tolerance); AC unchanged, actual-route/live-mutant proof blocked
- Cycle 5: CHANGES REQUESTED — Roborev job 744; resumed-cycle surface 4 files and 31 added LOC against the bounded 650-line behavioral and 59-line live caps; AC unchanged; 8n projection is an explicit integration dependency
- Cycle 6: NEEDS DECISION — Roborev job 775; behavioral proof remains capped at 650 added LOC and live proof at 59; AC unchanged; command-only copied-skill mutant has a credentialed blocking counterexample
- Cycle 7: REVISE — delegated proof-boundary ruling; surface 650 behavioral plus 59 live LOC vs caps 650 plus 59 (100%); AC unchanged, impossible/disproportionate proof mechanisms corrected
- Cycle 8: REVISE — Roborev job 837; two Material successor-oracle findings accepted; declared caps remain 650 behavioral and 59 live LOC; ACs and production Go unchanged
- Cycle 9: DESIGN RESET — Captain-approved proof-boundary ruling, 2026-07-24. Preserve the six-command fail-closed lifecycle, one successor dispatch, and its durable effect. Preserve deterministic command-deletion proof and one representative live journey per supported host. Remove mandatory public transport-event uniformity, exact-child forensic attribution where the supported host surface does not expose it, and every-route-by-every-host live execution. Host-native fixtures own transport details; live journeys own the observable lifecycle and durable outcome. Return to ideation for the smallest AC-1/AC-8 and test-surface delta before further implementation.
- Cycle 10: DESIGN RESET — Captain audit ruling, 2026-07-24. Preserve authorized record/consume-before-dispatch, durable successor state, and representative host journeys. Remove repeated four-subcommand help probes, the absolute-path workaround as a permanent FO ritual, runtime-specific “very next host event” wording, and zero-headroom absolute total-host ceilings. Spike whether every separate validation command adds safety before retaining exact six-command cardinality; keep resident/deferred component caps and define an attributable rebaseline policy. AC end values unchanged.
- Cycle 11: REVISE — independent ideation staff review; surface 1 design file with implementation not restarted vs expected implementation at or below 1,518 added LOC; AC end values unchanged. Preserve a state commit after every successful close, key the cached capability result to resolved executable identity, recompute task-attributable load at every tip/rebase, and prove consumed-commit ordering plus each omitted authority mutation through real traces.
- Cycle 12: REJECTED — validation resume matrix; surface 17 files/1,461 added LOC vs estimate 1,588 added LOC (92%); AC unchanged
- Cycle 13: DESIGN RESET — FO ruling after Roborev 2083. Preserve the three authority mutations and durability barriers; route a consumed terminal successor through the existing merge ceremony and only a nonterminal successor through ordinary dispatch. Replace the session identity/digest cache with one capability probe immediately before each gate lifecycle, deleting same-path/PATH cache machinery rather than adding a launcher-swap live laboratory. Pi review evidence must come from the root assistant message. Reuse the existing lifecycle fixture and host journey; no new harness, transport contract, or exact-child attribution.

- Cycle 14: REJECTED — Roborev jobs 2145–2158; surface 9 files/88 LOC vs estimate 9 files/75 LOC (117%); AC unchanged
- Cycle 15: REVISE — GitHub Actions offline CI run 30112268591; surface 21 files/1,470 added LOC vs 1,588-addition branch cap (93%); AC unchanged

- Cycle 16: REJECTED — Roborev job 2167 reviewing Cycle 15 CI feedback; surface 1 test file/2 added LOC vs estimate 1 test file/2 added LOC (100%); AC unchanged
- Cycle 17: DESIGN RESET — detached validation audit at `13d70249`; the 2-line identity repair is valid, but supported rejection/decision outcomes and AC-1/AC-6/AC-8 live evidence cross product and harness-controller boundaries; AC unchanged
- Cycle 18: DESIGN RESET — goal-only Codex lifecycle at `37d6980b`; surface 7 files/95 added LOC vs +95 hard stop (100%); AC unchanged; two unchanged runs omitted root-visible gate review after loading the corrected skill, triggering the approved skill-discovery/runtime-ownership reset
- Cycle 19: DESIGN RESET — Roborev job 2170 and supported-host live evidence at `08675f02`; surface 8 files/174 added LOC vs +175 hard stop (99%); AC unchanged; Codex and deterministic presenter proof pass, but legacy gate presentation, successful-commit attribution, onboarding scope, distinct Claude obligations, Pi runtime availability, and unresolved-material round disposition require design authority
- Cycle 20: DESIGN RESET — supported Pi lifecycle at `ce436505`; surface 17 files/497 added LOC vs +510 hard stop (97%); AC unchanged; Claude/Codex and offline gates pass, but Pi’s successful native successor is not attributable through the common oracle, its directive retains prompt delimiters, and the runner cleaned command/state evidence required to adjudicate either claim
- Cycle 21: DESIGN RESET — independent checkpoint review and correction prototype at `3c535105`; surface 3 correction files/+191 additions vs +95 repair cap (201%) before seven deleted v1 invariants are restored; AC authority must be re-anchored because expanding the forensic proof surface would reward the same overdesign

## Topology re-ideation delta (cycle 2; topology authority)

This delta remains authoritative only for deferred placement, byte budgets, and gate-entry ownership. The canonical acceptance criteria, test plan, and Obligation delta above own the proof boundary. The Problem, six-event lifecycle, provenance rules, fail-closed routing, concise presentation, capability readiness, exact package behavior, documentation wording, and no-production-Go boundary remain binding. Commit `cabdef33` remains an untouched counterexample/checkpoint.

### Deferred ownership and byte budget

The detailed `«gate.lifecycle»` body moves verbatim in substance from the boot-resident core to a new adapter-less `skills/fo-gate-lifecycle/SKILL.md`. Its frontmatter declares `name: fo-gate-lifecycle` and `user-invocable: false`; its description names the first engaged-gate trigger. It owns capability preflight, retained package rules, resume validation, open/closed binding, direct/delegated/provider provenance, eligibility/consume routing, stale/consumed behavior, and the handoff to existing feedback/dispatch procedures. It does not own judgment, `present-gate` rendering, write authority, feedback-round semantics, or runtime spawning.

The exact implementation budgets are:

- `skills/first-officer/references/first-officer-shared-core.md`: extract the 6,222-byte lifecycle section. The checkpoint is 26,067 bytes after extraction and before added route/pointer text; all gate-load triggers together have at most 687 bytes, and the finished file must be at most **26,754 bytes**. The hard `<26,755` guard is unchanged.
- `skills/fo-gate-lifecycle/SKILL.md`: new, **6,600 bytes maximum**, expected 62–75 lines. The 6,222-byte procedure may gain only frontmatter, an ownership heading, and load/write-order preconditions; lifecycle meaning may not be compressed away.
- Worst-case host accounting: register the new skill in `foFunctionReferencePaths` and `foSharedLoadPaths`, because every host can reach a gate. The implementation ceiling is current baselines plus the 6,600-byte skill plus at most 662 resident bytes versus the 26,092-byte parent: Claude **103,343**, Codex **82,558**, Pi **78,688**. `foHostLoadBaselineBytes` must be set to each candidate's exact measured bytes, with no unused allowance; the ceilings are stop limits, not ratchet values. The set-equality discriminator must fail if the file disappears from either address lint or worst-case accounting.
- Contract/topology guards: 40–80 test LOC across `internal/contractlint/boot_resident_closure_test.go`, `fo_function_reference_invariant_test.go`, and `structural_checks_test.go`. Register the skill in `lazyLoadSkills`/`deferredSkillCores`, require its lifecycle anchors, prove the resident pointer resolves, and prove it is absent from the user-invocable discovery set.
- Behavioral proof: follow the canonical test-plan retain/simplify/delete map and the current branch ceiling of 1,480 added LOC. The three checkpoint documentation edits remain as written and gain no new behavior.

No other production or instruction file is expected. Any production Go, schema/command change, host-specific lifecycle copy, `present-gate` redesign, feedback-round change, or file outside this list is a design-reset trigger.

### One gate-entry trigger and load order

The boot core owns only detection and the load precondition. Every engaged route funnels through one conceptual `gate.enter(slug, stage, route)` boundary:

```text
if route == interactive-greet:
    summarize the ready gate and STOP                 # no lifecycle load
if current/next stage is not gate:true:
    continue the ordinary event loop                  # no lifecycle load

complete Skill(skill="spacedock:fo-gate-lifecycle")  # its own host event
run the selected launcher's full gate capability probe
if resuming an existing attempt:
    run lifecycle resume validation before any replayed write
if the next action creates/changes room, entity, gate, or application state:
    complete the existing fo-write-core read          # separate prerequisite/event
retain/bind package -> validate open
AC cross-check/judgment -> Skill(spacedock:present-gate) -> present
record decision -> validate closed -> eligibility/route/consume
only after durable consumed state: dispatch build -> host-native dispatch -> durable successor effect
```

The required route table is:

| Entry route | Required behavior |
| --- | --- |
| Interactive startup with an already-gated entity | Greet names the gate and stops without loading the lifecycle or presenter. A later `engage` is the engaged entry below. |
| Headless/no-conn startup, including already-gated state | Load lifecycle before capability/package work, validation, or `present-gate`; bind and validate open, present, then stop without deciding. The checkpoint's direct presenter call is deleted. |
| Headless with delegated conn | Same entry and ordering; the grant changes only decision authority, never the pre-presentation sequence. |
| Interactive `engage` | Stay lazy through convergence/non-gated work; at the first gate, load before any gate action. |
| Worker completion whose current/next stage is gated | After report/checklist detection and before package mutation, validation, presentation, feedback routing, or advancement, load the lifecycle. |
| Resume/recovery of open, closed-pending, revise, hold, stale, or consumed records | Load before the first `gate validate` and before any state-specific routing. Consumed resumes dispatch without another consume; it still enters through the module so the one-use decision is observed. |

Loading `fo-gate-lifecycle` never satisfies write authority. The existing `fo-write-core` read remains a separate completed event immediately before the first FO-authored room/state mutation. Conversely, a write-core read cannot substitute for the lifecycle load.

Repository precedent already proves the adapter-less, non-user-invocable skill mechanism structurally (`fo-status-viewer` and `fo-dispatch-recovery`); no new loader or spike product is needed. Static contract tests own the complete route matrix, while one native load-before-action fixture per host and one representative live journey per host prove runtime portability.

### Checkpoint reuse and falsifiable proof

Retain the checkpoint's fresh real-binary six-command replay, exact attempt/Briefing/digest/Resolution/provenance assertions, exact empty open-decision field, help-event exclusion, relative/absolute path pair, folder-form commit, revise/hold controls, shared-scenario parity, and command/dispatch evidence artifacts. Host-native fixtures own Claude, Codex, and Pi transport details; the shared live oracle owns only the ordered commands, consumed state, one dispatch build, and one durable successor effect. Change contract-source reads from the shared core to the new skill and keep the resident trigger/topology assertions separate.

Replace two insufficient checkpoint checks:

- Replace the static one-off eligibility check with deterministic source-derived deletion of each of the six lifecycle events. Every mutant retains an otherwise-valid prospective dispatch observation, and the common grader must go red because the required trace is incomplete; model reconstruction of edited prose is not the oracle.
- Replace generic review substring acceptance with grading of the review actually emitted by the runtime: required structured fields, one recommendation/decision ask, retained snapshot identity, and no leading raw entity/Briefing/room dump.

Add these exact matrices:

| Matrix | Required arms and red control |
| --- | --- |
| Load topology | Interactive gated greet (no load); headless no-conn; headless conn; engage; worker-completion gate; open/pending/revise/hold/stale/consumed resume. Static shipped-contract extraction and route-deletion controls prove every funnel; one native load-before-action fixture per host proves transport parsing. The matrix is not multiplied into live route × host executions. |
| Six load-bearing events | Baseline plus deterministic deletion of briefing-record, open-validate, decision-record, closed-validate, eligibility, and consume. Each mutant keeps a valid prospective dispatch observation; the grader, not removal of dispatch evidence, refuses it. |
| Provenance/presentation | Delegated baseline; actor swapped to `person:captain`; blank reason; missing/exactly altered directive; raw entity dump; raw `briefing.json` dump. Actor/provenance mutants fail exact durable cardinality; raw-dump mutants fail the observed-review grader. |
| AC-5 refusals | Capability-stale executable with every help form logged; missing and alternate-basename Briefing; invalid association, actor, reason, and directive; relative Briefing refusal then absolute success; forced close-validation mismatch; validate/eligibility reads; hold, blocked, repeat-consume, and stale consume. Assert nonzero/actionable output, whole-tree/entity byte identity and no lock residue on every promised byte-clean arm; stale may change only pending→superseded and never status. The actual FO stale-launcher entry remains a required live arm, not replaceable by the helper probe. |
| AC-7 resume | Fresh processes over open/same package, open/changed package, closed approval-pending, closed revise, closed hold, stale, and consumed. Assert attempt/Resolution cardinality and bytes after each pass; across three passes exactly one successful consume and one transition. Stale must supersede, bind a replacement Briefing, and re-present without advancing. |
| Discovery | Capture the real workflow-discovery candidate list, create/run all fixtures, capture it again, and require exact sorted equality. A planted discoverable workflow proves the comparison can red. |
| Runtime successor | After consumed state, require exactly one successful `dispatch build`, pre-dispatch absence of the marker/report, and exactly one new marker-bearing successor report/commit. Zero/two builds, missing/duplicate effects, wrong ordering, or no durable commit fail. Host-native fixtures separately prove transport shapes; the common live oracle does not require public exact-child events. |

Claude 401 and Pi reused-refresh-token/no-key failures remain historical external validation conditions. Final validation requires one green identical host-neutral journey per supported host. Codex's supported public evidence is its native `dispatch build` → completed wait → durable report shape; the absent public `spawn_agent` event is not a product blocker or a reason to add a runtime interface.

### Stop conditions and excluded reset defects

Stop and return to ideation if the resident core reaches 26,755 bytes; the deferred file exceeds 6,600 bytes; any host worst-case file set omits the deferred skill or exceeds its ceiling; any engaged entry acts before the module load; any deterministic six-command deletion remains green; the full AC-5/AC-7/discovery matrices are incomplete; any live host lacks the six-command/consume/one-dispatch/one-durable-effect outcome; or implementation needs a surface excluded above. Full `go test ./...`, `go test ./... -race`, focused/live compilation, and one representative live lane per supported host remain final gates.

Two observed workflow-reset defects are recorded but explicitly not absorbed: revise currently derives feedback to implementation, and clearing a worktree on backward routing falsely trips the merge guard. This entity changes neither reset target derivation nor worktree/merge-guard behavior; those require a separately scoped correction.

## Stage Report: ideation (cycle 2)

- DONE: Define the exact deferred topology for the recorded gate lifecycle: a non-user-invocable gate-triggered module, a boot-core trigger that keeps first-officer-shared-core below 26,755 bytes, honest worst-case host prompt accounting, and load order before every mutation/validation/presentation/decision path.
  The authoritative delta assigns the 6,222-byte procedure to `skills/fo-gate-lifecycle/SKILL.md`, caps the resident core at 26,754 and the module at 6,600 bytes, charges it to all three host load sets, and specifies lifecycle-load then separate write-core ordering.
- DONE: Enumerate and close every gate-entry route—including already-gated startup, headless/no-conn, engage, worker completion, and resume—while preserving greet-and-stop/non-gated laziness; provide exact file/trigger ownership and a falsifiable load-topology test plan.
  The route table covers interactive and headless startup, conn/no-conn, engage, worker completion, and all resume states; route-deletion mutants must red while interactive greet and non-gated paths prove no lifecycle load.
- DONE: Replace the failed implementation plan with a bounded checkpoint-reuse plan that retains useful tests but adds the shipped-skill live mutant, actor-swap/raw-dump mutants, full AC-5/AC-7 matrices, workflow-discovery equality, and strict observed successor-spawn evidence; keep Claude/Pi auth as external validation conditions and do not narrow ACs.
  The retain/replace/add matrix names every missing proof arm and its red control; Codex consume-without-spawn remains blocking, Claude/Pi require later credentialed green runs, and AC-1 through AC-8 are unchanged.

### Summary

Re-ideated the rejected resident implementation as one gate-triggered, adapter-less `fo-gate-lifecycle` skill with a sub-26,755-byte boot core and honest all-host worst-case accounting. The revised plan closes the headless bypass and every other gate-entry route, preserves the detailed lifecycle and WIP checkpoint, and makes live mutation, complete refusal/resume matrices, discovery equality, and an actually observed successor spawn mandatory before validation.

## Historical implementation intended-change declaration (cycle 2)

- `skills/first-officer/references/first-officer-shared-core.md`: extract the checkpoint's 6,222-byte lifecycle body, add only the common gate-entry/load-order trigger, and finish at no more than 26,754 bytes; this is the boot-resident detector and pointer, not the procedure owner.
- `skills/fo-gate-lifecycle/SKILL.md`: add the sole adapter-less, non-user-invocable lifecycle owner at no more than 6,600 bytes; it is required for capability, package, provenance, closure, refusal, resume, consume, feedback, and dispatch handoff semantics.
- `internal/contractlint/boot_resident_closure_test.go`: add 24 test LOC for lazy/deferred registration, lifecycle anchors, and resident-pointer ownership.
- `internal/contractlint/fo_function_reference_invariant_test.go`: add 24 test LOC for address/worst-case set equality and exact per-host measured baselines that include the deferred skill.
- `internal/contractlint/structural_checks_test.go`: add 20 test LOC proving the skill is non-user-invocable, adapter-less, and absent from user discovery.
- `internal/ensigncycle/recorded_gate_lifecycle_test.go`: historical cycle-2 allocation; superseded by the canonical retain/simplify/delete map, deterministic six-event deletions, and common command/state/Git oracle above.
- The live/shared-scenario files retain one host-neutral scenario and one representative journey per host; host-native fixtures own transport details.
- `docs/dev/.spacedock-state/first-officer-gate-command-lifecycle/index.md`: append this declaration and the cycle-2 implementation report only; it owns workflow evidence, not product behavior.

Planned cycle-2 test addition is exactly 568 LOC (68 contract/topology plus 500 behavioral) and planned cycle-2 live-runner addition is exactly 0 LOC. Production Go, schemas, commands, presenter behavior, advisory-round surfaces, host-specific lifecycle copies, and all other instruction/documentation files remain untouched.

### Intended-change amendment 2

- `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go`: add exactly 1 live LOC registering the newly deferred adapter-less skill with Pi; without this explicit skill root the identical Pi scenario cannot consume the shipped lifecycle.
- `internal/ensigncycle/claude_live_runner_test.go`: add up to 45 live LOC for the required copied-plugin missing-command runtime mutant; this is the existing copied-plugin host seam and proves the runtime consumed the changed shipped tree while the common grader rejects the missing event.

The corrected plan is 568 test LOC plus at most 46 live LOC. These two additions are proof wiring only; they do not create host-specific lifecycle instructions or alter the shared scenario.

### Intended-change amendment 3

- `internal/ensigncycle/codex_live_runner_test.go`: change exactly 1 live LOC so the positive recorded-gate scenario's logging shim delegates to a fresh current-checkout binary, after retaining the separate observed stale-launcher halt artifact. The checkpoint's repo-root binary is intentionally stale and cannot be the positive lane.

The live-proof budget is now at most 47 LOC. The stale executable remains a required red live arm; this amendment only prevents that red control from replacing the fresh-binary positive scenario.

### Roborev triage (cycle 2, job 700)

- **Material — empty observed-review fields can pass.** Released user/workflow: a Captain receiving the ordinary FO gate review; harm: labels without evidence or a decision ask can be graded as a usable review; AC/boundary: AC-6 and the observed-review matrix; trigger: `assertConciseRecordedGateReview` accepts each required label with a blank value. Fix by parsing nonblank field values and pinning the retained snapshot and decision ask.
- **Material — stale consume can corrupt unrelated bytes and pass.** Released user/workflow: an FO resuming an ordinary stale approval; harm: another record, field, or file could be changed while the test notices only `status`/`superseded`; AC/boundary: AC-5 byte-clean exception and AC-7 stale resume; trigger: the current assertion checks only three substrings. Fix with a whole-tree before/after comparison whose independent expected tree permits only `pending` → `superseded`.
- **Material — changed-package resume does not pin replacement identity/digest.** Released user/workflow: an FO updating an open gate package; harm: the old binding could remain or unrelated fields could change while the test passes; AC/boundary: AC-7 open/changed-package resume and exact provenance; trigger: the test checks only byte inequality and attempt count. Fix by independently deriving the replacement canonical digest, asserting exact identity/digest, and comparing all other bytes.
- **Material — discovery equality queries a different root.** Released user/workflow: repository startup discovery after lifecycle fixtures; harm: a planted fixture workflow in the actual queried repository could pollute discovery undetected; AC/boundary: AC-8 exact workflow-discovery equality; trigger: `writeRecordedGateFixture` creates an unrelated temporary root. Fix by placing and exercising the fixture below the queried repository's ignored/test-owned directory, then proving the planted discoverable control changes that same root.
- **Material — load topology is self-authored rather than runtime-derived.** Released user/workflow: every engaged gate entry on every host; harm: deleting or reordering a real funnel load can remain green because the test's event arrays do not come from the shipped contract/runtime trace; AC/boundary: AC-1/AC-2 engaged-entry load order and the cycle-2 topology stop condition; trigger: the matrix supplies its own events and deletes only event zero. Fix by extracting route events from the copied shipped core plus deferred skill, grading each route, and adding route deletion and interior-order mutants through that extractor.

### Roborev follow-up triage (cycle 2, job 708)

- **Material — AC-7 resume compares only entity bytes.** Released user/workflow: an FO replaying any open/closed/consumed gate after restart; harm: the replay may create or mutate sibling state while entity-only equality remains green; AC/boundary: AC-7 exact pass cardinality and byte preservation; trigger: same/open, revise, hold, and consumed arms compare only `fixture.entity`. Fix every idempotent arm against a complete state-root snapshot.
- **Material — topology proof remains structural, not runtime-observed.** Released user/workflow: every engaged gate route on Claude, Codex, and Pi; harm: a host may fail to load the deferred skill before gate commands while the Markdown-derived trace remains green; AC/boundary: AC-1/AC-2 and the cycle-2 load-topology stop condition; trigger: `TestRecordedGateLifecycleLoadTopologyMatrix` reads only shipped Markdown. Retain it as structural proof, add host-stream load/command extractors with all route fixtures and red controls, and make each live recorded-gate runner grade its observed load order.
- **Material — successor controls bypass the real stream parsers.** Released user/workflow: post-consume successor dispatch on every host; harm: a parser regression could accept narration, a blank handle, parent-written output, or an empty wait; AC/boundary: AC-1 strict spawn/handle/correlated-output oracle; trigger: `TestRecordedGateLifecycleSuccessorOracleControls` constructs `recordedGateDispatchProof` directly. Fix by routing adversarial Claude, Codex, and Pi streams through their production test extractors.
- **Polish — prompt-load growth note.** Released user/workflow: none; harm: none, because exact baselines and hard ceilings are already measured and the boot core shrank; AC/boundary: reporting clarity only; trigger: the rebaseline comment does not explicitly contrast full-lifetime growth with boot residency. Decline this round; promote if a reviewer or release note conflates the two metrics.
- **Polish — byte-clean helper name includes stale caller.** Released user/workflow: none; harm: none, because the stale caller separately asserts the exact allowed mutation; AC/boundary: test naming only; trigger: `assertRecordedGateByteCleanFailure` checks nonzero/output/lock, not bytes itself. Decline; promote if a caller relies on the helper name without its own byte assertion.

### Intended-change amendment 4

- `internal/ensigncycle/claude_live_runner_test.go`, `internal/ensigncycle/codex_live_runner_test.go`, and `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go`: add one runtime-load-order assertion per positive recorded-gate runner, using the shared host-stream extractors added to `recorded_gate_lifecycle_test.go`.

This adds at most 12 live LOC to make the previously declared all-host topology proof load-bearing. It changes no prompt, scenario, host-specific lifecycle instruction, or acceptance criterion.

### Roborev final triage (cycle 2, job 711)

- **Material, fixed — runtime action detection could accept echoed skill prose.** Released user/workflow: every live recorded-gate route; harm: an echoed `gate --help` inside a skill read could false-green load order; AC/boundary: AC-1/AC-2 observed load-before-action proof; trigger: the extractor matched any JSON line containing the text. Fixed by requiring each host's structured Bash/command-execution event and recognizing the first actual `spacedock gate` command, not only help.
- **Needs decision — route labels do not execute distinct runtime paths.** Released user/workflow: headless conn/no-conn, engage, worker completion, and every resume state; harm: a host-specific route can bypass the load while synthetic route fixtures remain green; AC/boundary: the binding cycle-2 topology matrix and stop condition; trigger: every route subtest currently reuses one normalized load/action fixture. The required actual-host matrix cannot be made green in this cycle: Claude is 401-blocked, Pi cannot load its extension dependency, and Codex omits the required structured spawn event. Do not substitute fixtures or exceed the 650-LOC proof cap; return for captain/FO reset or repaired live hosts.
- **Needs decision — copied-skill command deletion lacks green runtime execution.** Released user/workflow: a live FO consuming the shipped deferred skill; harm: a host could ignore a missing command while text-derived checks pass; AC/boundary: AC-2 and the copied-plugin live-mutant stop condition; trigger: deterministic deletion uses `procedureEvents`, while the actual Claude copied-plugin mutant cannot start because of 401 credentials. The original six-event baseline and all six deterministic deletions are now explicit, but the required runtime mutant remains blocked pending credential repair; do not claim validation readiness.
- **Material, fixed — deterministic missing-event discriminator was removed.** Released user/workflow: the common successor grader; harm: an event-completeness regression could escape the more specialized source/runtime checks; AC/boundary: AC-2 six-event load-bearing proof; trigger: job 711 observed the removed deterministic control. Restored a compact six-arm control that requires the grader's event-completeness error.

## Stage Report: implementation (cycle 2)

- DONE: Move the detailed lifecycle into the approved non-user `fo-gate-lifecycle` skill and funnel every engaged gate-entry route through one deferred load before gate action; shared core remains ≤26,754 bytes, deferred skill ≤6,600 bytes, all host worst-case sets account for it exactly, and non-gated greet remains load-free.
  Commit `9a3cc782` leaves the boot core at 26,526 bytes and the adapter-less deferred skill at 6,599 bytes; exact host loads are Claude 103,114, Codex 82,329, and Pi 78,459, each below its unchanged ceiling. Contract/topology tests fail if address/load sets diverge, the deferred file disappears, the skill becomes user-invocable/adapter-backed, the headless route bypasses lifecycle ownership, or interactive greet loads it.
- FAILED: Retain and strengthen the checkpoint proof with route-deletion/load-order controls, a real copied-skill missing-command live mutant, exact observed-review and provenance mutants, complete AC-5 refusal and AC-7 resume matrices, discovery equality, and a strict post-consume successor oracle requiring one spawn, nonblank handle, and correlated worker output.
  Fresh real-binary replay, all six deterministic command deletions, structured host load/action matrices, exact review/provenance controls, complete AC-5/AC-7 command matrices, whole-tree byte checks, discovery equality/red control, and real-parser successor negatives pass. Roborev job 711 correctly blocks completion because route names still reuse normalized runtime fixtures and the actual copied-plugin mutant cannot execute under Claude's 401; weakening those requirements or replacing them with fixture narration is prohibited.
- FAILED: Preserve zero production Go and unchanged ACs, pass focused/full/race/prompt/load/discovery gates and required credentialed live lanes, request Roborev and triage findings before edits, then commit the self-contained deliverable and report validation readiness only if Codex spawn evidence and later green Claude/Pi evidence exist.
  Production Go remains unchanged; `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, focused lifecycle/contract tests, and live-tag compilation pass. Roborev jobs 700, 708, and 711 were triaged before edits; all bounded material fixes landed in `9a3cc782`. Live readiness is red: stale Codex correctly halts before mutation; fresh Codex consumes and produces the durable worker marker but emits no qualifying structured spawn event/handle; Claude fails 401 before FO work; Pi fails loading `@earendil-works/pi-coding-agent`.

### Summary

Cycle 2 successfully moved the recorded gate lifecycle behind the approved deferred skill while preserving byte budgets, honest host accounting, zero production Go, and the complete deterministic command/refusal/resume/discovery checkpoint. The strict live proof remains intentionally fail-closed: current host evidence cannot establish every route's load order, the copied-plugin runtime mutant, or Codex's required structured spawn/handle/result correlation. The committed checkpoint is therefore not validation-ready and the entity remains in implementation for a captain/First Officer decision or repaired live-host conditions.

### Roborev resumed-cycle triage (job 744)

- **Needs decision/dependency — 8n readiness projection is not in this branch.** Released user/workflow: an FO selecting an actionable recorded gate from boot; observable harm: this branch alone would omit every projected row and the new instruction forbids stage-name inference; affected value/boundary: AC-1 and the captain-routed no-production-Go dependency on 8n; trigger: job 744 compared the skill to this worktree's pre-8n binary. Retain the captain-provided exact contract, make 8n's projection a merge prerequisite, and do not duplicate its production Go here. Promote to Material/block merge if 8n is absent from the integration base.
- **Material, accepted — aggregate Claude Bash success can mask an earlier command failure.** Released user/workflow: the six-event Claude lifecycle; observable harm: a failing recorder call followed by a successful validator could false-green; affected value/boundary: AC-1/AC-2 successful ordered events; trigger: job 744 observed that physical-line splitting inherited one overall exit status. Require fail-fast `&&` between multiple gate commands before distributing the successful tool result, and add the unsafe-separator red control.
- **Material, accepted — Pi load was not correlated to a successful read result.** Released user/workflow: Pi gate entry before action; observable harm: configuration text or a failed `read` could satisfy the load-order assertion; affected value/boundary: AC-1/AC-2 and the runtime-observed topology stop condition; trigger: job 744 observed path-only matching. Correlate the structured Pi `read` tool-call id to its non-error tool result before accepting the load.
- **Polish, declined — split Pi provider/model flags onto separate source lines.** Released user/workflow: none; observable harm: none; affected value/boundary: readability only, while the approved live-test budget is exactly 59 added lines; trigger: job 744's style note. Promote only if the line becomes functionally ambiguous or the live budget is explicitly reset.

### Roborev final resumed-cycle triage (job 775)

- **Material, accepted — reason validation must reject canonical empty scalars.** Released user/workflow: delegated FO approval; observable harm: whitespace or quoted-empty reason could pass as evidence; affected value/boundary: AC-3 nonblank evidence judgment; trigger: job 775 observed the relaxed substring check. Anchor the unique canonical Resolution `reason:` field and trim whitespace plus YAML quote characters before accepting it. Do not restore one historical sentence: the contract deliberately permits the FO's evidence judgment while exact actor and Captain directive remain pinned.
- **Needs decision/blocker — a two-part mutation would not prove the command alone load-bearing.** Released user/workflow: the required copied-skill missing-command red control; observable harm: removing explanatory authorization prose with the command manufactures a red without isolating command deletion; affected value/boundary: AC-2 and the explicit command-only live-mutant stop condition; trigger: job 775, plus serial artifact `/tmp/spacedock-fo-gate-live-claude-mutant-serial`, where Claude loaded the copied command-only mutant, reconstructed eligibility from remaining instruction, emitted all six events, consumed, and dispatched. Revert to command-only deletion and retain this as a real blocker; do not weaken the criterion or claim green.
- **Low, accepted narrowly — prove the copied skill body as well as its directory.** Released user/workflow: mutant provenance; observable harm: another file in the copied directory could theoretically be consumed; affected value/boundary: AC-2 runtime mutant; trigger: job 775. Claude's native load event reports the base directory rather than `/SKILL.md`, so require that exact directory and absence of the removed placeholder command in the emitted loaded body.
- **Low, declined — semicolon-joined multiple gate commands are conservatively discarded.** Released user/workflow: none in the controlled route; observable harm: false negative only, never false green; affected value/boundary: AC-1 successful-event attribution; trigger: a host emits multiple gate calls joined only by `;`. Promote if a supported route intentionally emits per-command exit capture for that shape.
- **Low, declined — detect alteration from one historical reason sentence.** Released user/workflow: delegated approvals legitimately carry the FO's current evidence judgment; observable harm: pinning the historical sentence would reject correct live decisions; affected value/boundary: AC-3 requires nonblank evidence plus exact actor/directive, not exact prose; trigger: a binding specification changes to require a canonical reason string.

## Stage Report: implementation (cycle 3)

- DONE: Move the detailed lifecycle into the approved non-user `fo-gate-lifecycle` skill and funnel every engaged gate-entry route through one deferred load before gate action; shared core remains ≤26,754 bytes, deferred skill ≤6,600 bytes, all host worst-case sets account for it exactly, and non-gated greet remains load-free.
  Commits `9a3cc782` and `2680ce97` retain the approved adapter-less topology. The resident core is 26,526 bytes, the deferred skill is 6,594 bytes, and the exact host prompt baselines remain below their ceilings. The skill consumes 8n's exact `ready_gates` projection (`id`, `slug`, `current`, `readiness`; three actionable readiness values and fail-closed omissions) without inferring readiness from status/stage or adding production Go. 8n landing that projection is an explicit integration prerequisite.
- FAILED: Retain and strengthen the checkpoint proof with route-deletion/load-order controls, a real copied-skill missing-command live mutant, exact observed-review and provenance mutants, complete AC-5 refusal and AC-7 resume matrices, discovery equality, and a strict post-consume successor oracle requiring one spawn, nonblank handle, and correlated worker output.
  Deterministic route deletion, structured host load/action controls, exact review/provenance, complete AC-5/AC-7 matrices, discovery equality/red control, aggregate-exit rejection, Pi read/result correlation, and real Claude spawn/result parsing pass at the hard 650-added-line proof cap. The required command-only copied-skill live mutant does not red: credentialed serial artifact `/tmp/spacedock-fo-gate-live-claude-mutant-serial` proves Claude loaded the copied directory, reconstructed the removed eligibility command from remaining semantics, emitted all six events, consumed, and dispatched. Removing both command and authorization prose would manufacture a red and was rejected by Roborev job 775. Actual distinct runtime executions for every named gate-entry route also remain absent; normalized fixture labels cannot substitute.
- FAILED: Preserve zero production Go and unchanged ACs, pass focused/full/race/prompt/load/discovery gates and required credentialed live lanes, request Roborev and triage findings before edits, then commit the self-contained deliverable and report validation readiness only if Codex spawn evidence and later green Claude/Pi evidence exist.
  Production Go and ACs remain unchanged. Current-tree `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, focused lifecycle/contract/discovery controls, and `go test -tags live ./internal/ensigncycle -run '^$' -count=1` pass. Roborev jobs 744 and 775 are fully triaged above; accepted material fixes are in `2680ce97`.
  Live readiness remains red. Claude's positive credentialed runs performed the six commands and real Agent spawn with correlated durable output, but the last run stopped on the now-corrected historical-reason overconstraint rather than a green test result. The honest command-only mutant remains green as described above. Pi performed the full six-command lifecycle and returned native async handle `dc46d149-43f3-4032-830b-7c96e42e69dc`, but `pi --print` exited before the child result and the fixture was graded/cleaned too early; artifact `/tmp/spacedock-fo-gate-live-pi-fresh`. Codex's public `codex exec --json` stream still exposes no spawn event or nonblank receiver handle although private rollout files prove an internal child; accepting that private archive would weaken the declared public-stream oracle. No validation-readiness claim is made.

### Summary

Hardened the deferred recorded-gate lifecycle and its host parsers to the exact approved byte/LOC limits, folded in 8n's precise readiness contract, repaired Claude stored-login and fresh-launcher execution, and obtained real Claude/Pi lifecycle and worker activity. The strict proof correctly remains fail-closed on three unresolved boundaries: a command-only copied-skill mutant that the model semantically reconstructs, Codex's unobservable public spawn handle/result, and Pi print-mode async completion correlation. The bounded implementation is committed, but the entity must remain in implementation and must not advance to validation.

## Stage Report: implementation (cycle 4)

- DONE: Move the detailed lifecycle into the approved non-user `fo-gate-lifecycle` skill and funnel every engaged gate-entry route through one deferred load before gate action; shared core remains ≤26,754 bytes, deferred skill ≤6,600 bytes, all host worst-case sets account for it exactly, and non-gated greet remains load-free.
  Commit `3cc6225b` consumes merged 8n `12e7bccb`; core is 26,526 bytes and satellite is 6,594 bytes.
- DONE: Retain and strengthen the checkpoint proof with route-deletion/load-order controls, a real copied-skill missing-command live mutant, exact observed-review and provenance mutants, complete AC-5 refusal and AC-7 resume matrices, discovery equality, and a strict post-consume successor oracle requiring one spawn, nonblank handle, and correlated worker output.
  Per the binding ruling, copied-body provenance plus six deterministic event-deletion controls replaces the reconstructible command-only mutant. Behavioral proof is exactly 650 added lines; live proof is 54.
- DONE: Preserve zero production Go and unchanged ACs, pass focused/full/race/prompt/load/discovery gates and required credentialed live lanes, request Roborev and triage findings before edits, then commit the self-contained deliverable and report validation readiness only if Codex spawn evidence and later green Claude/Pi evidence exist.
  Zero production Go changed. Format, full, race, and diff checks pass. Claude `/tmp/spacedock-fo-gate-green-claude`, Codex `/tmp/spacedock-fo-gate-sequential-codex`, and async Pi `/tmp/spacedock-fo-gate-green-final3-pi` pass. Roborev job 837 was requested.

### Summary

All binding correction boundaries are closed without duplicating 8n. Three credentialed hosts prove the lifecycle through correlated durable successor output; implementation is ready for validation.

## Stage Report: implementation (cycle 5)

- FAILED: Require a correlated public Codex `spawn_agent` call, child completion, and exact-child durable output; parent writes, narration, waits without a receiver, and synthetic commit handles must fail.
  Credentialed artifact `/tmp/spacedock-fo-gate-sequential-codex/codex-shared-scenarios/recorded-gate-lifecycle/codex-exec.jsonl` has no public spawn or child-completion event. After `item_25` dispatch-build, `item_26` only narrates `/root/recorded_gate_task_handoff`; `item_27` is `tool:"wait"` with `receiver_thread_ids:[]`; `item_28` is a parent-side state read. This cannot be correlated without synthesizing the missing public handle.
- SKIPPED: Parse Pi child assistant/tool-result events and require child-produced write evidence correlated with the completed async handle and resulting durable state; add parent-writes/no-op-child red controls.
  Stopped at the explicit Codex public-evidence boundary before expanding the zero-free 650-line behavioral budget; Pi remains a material accepted finding for the next correction after a valid Codex public surface is available.
- SKIPPED: Rerun focused/full/race and all live lanes, request final-tip Roborev, wait, and triage every finding.
  No code change was made because the required public Codex observation is absent; rerunning cannot manufacture the missing event, and weakening the oracle is forbidden.

### Summary

Roborev 837 is fully triaged and both Medium findings are Material. The implementation remains blocked in implementation on a missing public Codex spawn/completion correlation; no validation-readiness claim is made.

## Stage Report: ideation (cycle 3)

- DONE: Rewrite only the AC-1/AC-8 proof boundary and dependent test plan so the six-command lifecycle, one successor dispatch, durable effect, deterministic command-deletion controls, and one representative live journey per supported host remain mandatory without requiring unavailable public transport events or every-route-by-every-host live execution.
  The authoritative reset gives all three live lanes one command/state/Git oracle while static route coverage and six source-derived deletion mutants remain mandatory.
- DONE: Identify the exact existing tests and LOC to retain, simplify, or delete; show that host-native fixtures own transport details while live journeys own observable lifecycle outcomes, and quantify the expected net reduction from the current 1,879-line branch.
  Current-tip line ranges assign 479 deletions and at most 80 replacement lines, capping the branch at 1,480 added LOC (at least 399 LOC/21.2% smaller).
- DONE: Produce a concise implementation delta and validation plan that preserves every product/value promise and requires no new recorder, runtime, or public host interface; make no product changes during ideation.
  Only the shared test oracle and live bindings change; AC-2–AC-7, skill/docs/topology, zero-production-Go scope, and full/race/three-host validation remain intact.

### Summary

Reset the proof mechanism without narrowing the product: every host must still execute the six-command lifecycle, consume once, dispatch once, and leave one durable successor effect. Transport details now stay in host-native fixtures, eliminating unsupported cross-host forensic requirements and at least 399 lines from the current implementation plan.

## Stage Report: ideation (cycle 4 repair)

- DONE: Edit the canonical `## Acceptance criteria` so AC-1 and AC-8 express the captain-approved proof boundary.
  `spacedock status --read first-officer-gate-command-lifecycle --workflow-dir docs/dev --ac-scan --json` anchors AC-1 at line 241 and AC-8 at line 255; both now require six commands, one dispatch, durable effect, and one live journey per host without public exact-child forensics.
- DONE: Update the canonical test plan and relevant design text rather than relying on a later override.
  The canonical replay plan, expected surface, topology handoff, checkpoint matrices, and stop conditions now use the command/state/Git oracle and deterministic six-event deletion.
- DONE: Remove or condense duplicate/competing authority so there is one legible specification.
  The duplicate later AC/test-plan override is removed, cycle-2 is labeled topology-only, and prior intended-change declarations are labeled historical.
- DONE: Add a concise `## Obligation delta` identifying preserved obligations and authority, removed invented obligations, and bearer/burden.
  The table assigns lifecycle, dispatch/effect, route, and host-transport duties; the Captain ruling removes only uniform public exact-child evidence and route × host live multiplication.
- DONE: Prove `status --read ... --ac-scan` returns the repaired ACs.
  Scanner JSON returns canonical line anchors 241/255; resolving those lines through the returned entity path shows the repaired AC-1/AC-8 text, and `git diff --check` is clean.
- DONE: Make no product-code changes and report exact files, evidence, and commit.
  Only `docs/dev/.spacedock-state/first-officer-gate-command-lifecycle/index.md` changed; the state commit is path-scoped and the code worktree remains clean at `3cc6225b`.

### Summary

Repaired the authoritative body rather than stacking another override: scanners, readers, and implementers now see one proof boundary at the canonical AC/test-plan location. All product obligations remain; only unsupported proof inventions are removed.

## Stage Report: implementation (cycle 6)

- DONE: Declare the exact reduction surface before editing and remain within the approved cap.
  The declared surface was four test/live-binding files with at most 80 additions and at least 479 deletions. Commit `86bad049` changes exactly those four files by **+71/-482**: Claude `+22/-10`, Codex `+2/-12`, Pi `+4/-37`, and the common lifecycle oracle `+43/-423`. The branch is **1,468 added LOC**, below the 1,480 ceiling. The estimate variance is nine fewer replacement lines and three additional deletions; no production Go, schema, command, runtime interface, skill, or documentation behavior changed.
- DONE: Implement the canonical proof-boundary reduction while preserving AC-2 through AC-7.
  One shared command/state/Git oracle now grades six ordered successful lifecycle calls, consumed successor state, exactly one successful post-consume `dispatch build`, and exactly one marker-bearing successor commit. Cross-host exact-child parsers, the route-by-host runtime matrix, Pi async polling, and their controls are deleted. Claude, Codex, and Pi retain host-native transport fixtures while their representative journey bindings use the common observable oracle.
- DONE: Run the bounded pre-integration verification and classify live evidence without weakening the oracle.
  Focused lifecycle/contract tests, `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, `go test ./... -race`, and live-tag compilation pass at `86bad049`. Representative live journeys pass on Codex (`/tmp/spacedock-fo-gate-reduction-live/codex`), Claude (`/tmp/spacedock-fo-gate-reduction-live/claude-zdot`), and Pi (`/tmp/spacedock-fo-gate-reduction-live/pi-sonnet46-final`), each with the six commands, one post-consume dispatch, and one durable marker commit. Earlier Pi peer-dependency/OAuth/default-child-auth failures were runtime harness friction; an extra-open-validation run was a genuine AC-2 red and was not used to relax the grader.
- BLOCKED: Reconcile the committed reduction with current `origin/main`, rerun final-tip evidence, and request final-tip Roborev.
  Rebase onto `origin/main` `dd6bd114` stops at the exact prompt-load ceiling. The merged advisory-round recorder grows the shared-load `skills/feedback-rejection-flow/SKILL.md` from 4,032 to 4,386 bytes. Arithmetic: original exact loads + fr's required **354 bytes** - prior **229 bytes** headroom = **125 bytes over** each ceiling. Applying 6y's exact deferred lifecycle produces measured totals Claude **103,468**, Codex **82,683**, and Pi **78,813**, above the canonical hard ceilings 103,343 / 82,558 / 78,688. Raising the ratchets would defeat the stop condition; trimming fr's skill would cross its concurrent-owner boundary. The partial rebase was aborted, preserving clean commit `86bad049`; no Roborev request or validation-ready claim is made for an unreconciled tip.

### Summary

The approved reduction is complete, committed, within its LOC cap, and green across deterministic, full/race, and all three representative live-host journeys. Integration is honestly stopped on a newly landed shared prompt-load collision. Return to a bounded 6y design reset that finds at least 125 bytes in 6y's own deferred-lifecycle/resident-pointer surface while preserving semantics, the 6,600/26,754 caps, and all proof; do not increase the baselines or edit fr's file. Rebase, final-tip replay, and Roborev follow only after that ruling.

## Historical rejected byte-budget delta (cycle 5; superseded by cycle 6)

This rejected delta records the prior 125-byte proposal for audit only. Captain cycle-6 feedback supersedes its six-command assumption, exact total-host ceilings, absolute-path ritual, and zero-margin compression; the canonical lifecycle, ACs, test plan, expected surface, and attributable rebaseline policy above govern implementation.

### Exact reduction

Apply two line-preserving rewrites in `skills/first-officer/references/first-officer-shared-core.md` at `86bad049`:

1. Line 46, the deferred-load pointer, changes from:
   ```text
   - `Skill(skill="spacedock:fo-gate-lifecycle")` — every engaged gate entry: headless with or without conn, `engage`, gated worker completion, and open/pending/revise/hold/stale/consumed resume. Complete this load before every capability probe, gate read/write/validation, presenter load, decision route, replay, or dispatch; interactive gated greet only names the gate and stops load-free.
   ```
   to:
   ```text
   - `Skill(skill="spacedock:fo-gate-lifecycle")` — load for headless with or without conn, `engage`, gated worker completion, and open/pending/revise/hold/stale/consumed resume; interactive gated greet only names the gate and stops load-free. Complete before gate probe/read/write/validation, presenter, decision route, replay, or dispatch.
   ```
   The exact UTF-8 line size falls from 391 to 341 bytes: **-50 bytes**.
2. Line 81, the gated-completion pointer, changes from:
   ```text
   If the stage is gated, first complete `Skill(skill="spacedock:fo-gate-lifecycle")`, then `«gate.lifecycle»(slug, stage)`. The deferred procedure binds and validates the package, presents, records and validates the decision, and authorizes approval through one-use consume before dispatch; revise routes through `«feedback.route»`, while hold and every ineligible condition stop.
   ```
   to:
   ```text
   If gated, complete `Skill(skill="spacedock:fo-gate-lifecycle")`, then `«gate.lifecycle»(slug, stage)`. It binds and validates the package, presents, records and validates the decision; approval permits dispatch only after one-use consume, revise invokes `«feedback.route»`, and hold/ineligibility stops.
   ```
   The exact UTF-8 line size falls from 383 to 308 bytes: **-75 bytes**.

The total is exactly **-125 bytes**, the mathematical minimum that satisfies “at least 125.” The first rewrite retains every engaged route, the greet exception, and load-before-action ordering; only repeated determiners and “load” nouns disappear. The second retains load-before-call, bind/open and close validation, presentation, consume-before-approved-dispatch, revise routing, and fail-closed hold/ineligibility. The detailed 6,594-byte lifecycle remains the semantic owner, so neither pointer needs to restate its full ceremony.

### Composition and cap proof

| Surface | Rebased before | Reduction | Proposed | Hard cap | Margin |
| --- | ---: | ---: | ---: | ---: | ---: |
| Claude worst-case load | 103,468 | -125 | **103,343** | 103,343 | 0 |
| Codex worst-case load | 82,683 | -125 | **82,558** | 82,558 | 0 |
| Pi worst-case load | 78,813 | -125 | **78,688** | 78,688 | 0 |
| Shared core | 26,526 | -125 | **26,401** | 26,754 | 353 |
| Deferred lifecycle skill | 6,594 | 0 | **6,594** | 6,600 | 6 |

An in-memory UTF-8 measurement of both replacements produced `load=391->341 completion=383->308 saved=125`; a token-presence control retained every literal route discriminator required by `TestFOGateLifecycleOwnsEveryEngagedEntry`. No fr, ratchet, host adapter, command, schema, runtime interface, test oracle, or AC changes.

### Reapply and final-tip verification

Replay the 6y implementation through `3cc6225b` onto `dd6bd114`, apply only the two pointer rewrites, then replay `86bad049` unchanged. Confirm its four approved test/live files remain **+71/-482** and the branch remains 1,468 added LOC; the line-preserving prose rewrite changes bytes, not LOC.

Run the focused topology/load and lifecycle tests first: `TestFOGateLifecycleOwnsEveryEngagedEntry`, deferred-skill reachability/discovery/prose-pointer controls, `TestFOHostPromptLoadRatchet`, load-set equality, function-reference invariants/metrics, and the recorded lifecycle source-deletion/common-oracle suite. The metrics log must print exactly 103343/82558/78688 and `wc -c` must print shared core 26401 and lifecycle skill 6594.

Then run `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, `go test ./... -race`, live-tag compilation, and the representative Claude/Codex/Pi journeys. Each live journey still owes six ordered commands, consumed successor state, one post-consume dispatch build, and one durable successor effect. Only after all final-tip evidence is green request Roborev and triage every finding; do not claim validation readiness before that review.

## Stage Report: ideation (cycle 5)

- DONE: Define the smallest semantics-preserving reduction of at least 125 bytes confined to 6y's deferred lifecycle and resident-pointer surface; do not edit fr or raise any host ceiling.
  Two line-preserving resident-pointer rewrites save exactly 125 UTF-8 bytes while retaining every route, load-order, lifecycle, routing, and fail-closed obligation; fr and all ratchets remain untouched.
- DONE: Prove the proposed rebased composition fits Claude 103343, Codex 82558, Pi 78688 while preserving the 26754 shared-core and 6600 lifecycle-skill caps.
  Direct byte arithmetic and in-memory replacement measurement yield 103343/82558/78688; shared core becomes 26401 and the unchanged deferred skill remains 6594.
- DONE: Name the exact files/lines, semantic equivalence argument, and focused verification needed to reapply preserved reduction commit 86bad049 and reach final-tip review without expanding product scope.
  The cycle-5 delta names `first-officer-shared-core.md:46` and `:81`, exact before/after text, retained meanings, replay order, focused/full/race/live gates, and the final-tip Roborev boundary.

### Summary

The exact integration repair is a 125-byte resident-pointer compression with zero semantic or product-surface change. It composes the preserved `86bad049` reduction with `dd6bd114` exactly at all three canonical host ceilings while leaving both instruction-file caps green.

## Stage Report: ideation (cycle 6)

- DONE: Produce the smallest lifecycle design that preserves authorized record/one-use consume before dispatch and a durable successor effect while removing provider/path/runtime mechanics already owned elsewhere.
  Canonical AC-1/AC-8 now require three authority mutations, package/consumed state commits as durability barriers, and ordinary post-consume dispatch; CLI normalization owns retained paths and host adapters own transport observation.
- DONE: Spike whether every separate validate/help command contributes independent safety; recommend the minimum falsifiable capability and lifecycle sequence rather than preserving six commands by assumption.
  Code-path inspection plus `go test ./internal/gates ./internal/ensigncycle -run 'TestRecordedGateLifecycle(CapabilityStaleLauncherHaltsBeforeMutation|ReviseAndHoldDoNotApprovalDispatch|HappyPath|Skipped|Resume)|Test(Consume|Application|Gate)' -count=1` passed; removing either `record` or `consume` internal guards would red these fixtures, while the design's new three-event deletion controls will red each omitted authority mutation.
- DONE: Replace zero-margin total-host byte ceilings with component budgets plus an attributable rebaseline policy, then update ACs, expected surface, and existing fixture/live verification without adding a new test framework.
  Shared core remains capped at 26,754 bytes, deferred lifecycle at 6,600; host totals are informational, task-owned load delta must be non-positive, and rebaseline may absorb only committed upstream files outside this task's touched surface.

### Summary

Cycle 6 reduces the happy path from six gate commands to three load-bearing mutations and the readiness check from five help calls to one session-cached fingerprint. It moves relative-path repair into the CLI, ends the lifecycle at a durable consumed commit, and reuses the existing scenario, refusal fixtures, live bindings, and host-native dispatch proofs.

## Stage Report: ideation (cycle 7)

- DONE: Retain a durable state commit after every successful close, including revise, hold, and approval that later fails consume, while keeping the three gate mutations as the authority sequence.
  Canonical lifecycle and AC-1/AC-3 now require an exact close commit before every route; successful approval adds a descendant consumed commit before dispatch, while stopped/refused routes retain the Resolution.
- DONE: Key the session capability cache to resolved executable identity and add same-path replacement plus PATH-target-swap controls; recompute task-attributable load at every tip/rebase while rebaselining only upstream-owned bytes.
  AC-5 keys reuse to canonical target path plus content digest and names both invalidation fixtures; expected-surface policy recomputes tip-minus-current-upstream bytes after every candidate/rebase and forbids task-owned bytes in the baseline.
- DONE: Upgrade existing proof so live traces establish close/consumed commit ordering before dispatch and each omitted authority mutation is exercised through the real binary rather than only an expected-slice mutation.
  The existing wrapper/oracle now owes logged commit hashes and entity snapshot digests at both barriers and dispatch pre-HEAD; three existing-fixture omission subtests invoke the remaining real commands. Current baseline tests `TestRecordedGateLifecycleRealCLIReplay`, capability/missing-event controls, and focused gate CLI tests pass; implementation must replace the slice-only control rather than claim it as final proof.

### Summary

Cycle 7 closes the staff-review gaps without expanding the authority sequence: every decision is durable, cache reuse follows executable identity, load attribution is recomputed, and proof follows real command/Git history. Mandatory diagnostics, repeated probes, caller path workarounds, total-host ceilings, and generic host-event requirements remain excluded.

## Implementation intended-change declaration (cycle 7)

Rebased preserved commits onto current `origin/main` `dd6bd114`; clean tip `2f7a6414` remains 1,468 added / 11 deleted LOC across the approved existing branch surface. Cycle-7 work is limited to these eight files:

| File | Intended cycle-7 delta |
| --- | ---: |
| `internal/cli/cli.go` | +6–10 / -0 |
| `internal/cli/gate_test.go` | +18–25 / -0–5 |
| `skills/fo-gate-lifecycle/SKILL.md` | +15–25 / -25–40 |
| `skills/first-officer/references/first-officer-shared-core.md` | +1–3 / -2–4 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +45–65 / -35–55 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +12–18 / -25–45 |
| `docs/site/concepts/gates-and-decisions.md` | +3–6 / -3–6 |
| `docs/site/reference/command-reference.md` | +3–4 / -1–2 |

Total intended amendment: at most 156 additions and at least 66 deletions, within the approved +170/-50 correction envelope and 1,588-LOC branch cap. No other file, harness, host adapter, schema, recorder/application package, or runtime transport surface is intended; any need outside this table stops for drift review before editing it.

## Stage Report: implementation (cycle 7)

- DONE: Implement the three-mutation recorded-gate lifecycle and narrow retained-input CLI normalization within the declared eight-file surface.
  Tip `624a9c2e` records Briefing, delegated decision/Result, and one-use consume as the only authority mutations; commits the retained package before presentation, every successful close before routing, and the consumed descendant before ordinary dispatch. Relative `--briefing`, `--result`, and `--association` paths now normalize against launch cwd. Legacy unbound human gates retain a presentation-only, non-mutating compatibility route. The exact cycle-7 delta is **+162/-221**, within the approved +170/-50 correction envelope; the branch remains **1,469 added LOC**, below 1,588.
- DONE: Make capability identity, presentation timing, and Git durability falsifiable without expanding generic host transport.
  The counting control caches one successful `gate --help` only for canonical target plus content digest, reprobes a same-path replacement, and reprobes a PATH retarget before rejecting stale capability without workflow mutation. Every live lane requires exactly one successful help call before mutation. The transcript grader accepts the six-field review only after the successful open-package commit and before the decision command. Git history requires that bound no-Resolution snapshot as the decision commit's parent, then closed and consumed descendants before the single dispatch build and marker-bearing durable effect. The structural command mutant test is named and scoped as structural evidence; three omitted-mutation controls execute the remaining commands through the freshly built binary.
- DONE: Preserve component/load and documentation boundaries.
  Shared core is **26,540 bytes** (cap 26,754); `fo-gate-lifecycle` is **5,837 bytes** (cap 6,600). Task-attributable instruction changes remain non-positive against the rebased parent. Strict MkDocs succeeds. Concepts/reference docs describe the three authority writes, package/close/consume commit barriers, executable-identity cache, launch-cwd retained paths, and ordinary dispatch handoff.
- DONE: Pass deterministic, repository, compatibility, and representative live-host verification.
  `gofmt -w ./cmd ./internal`, `git diff --check`, focused CLI/contract/lifecycle suites, `go test ./...`, `go test ./... -race`, and strict docs all pass on the final tip. Legacy unbound gate-guardrail journeys pass on Claude (**61.711s**) and Codex (**92.973s**). Final recorded-lifecycle journeys pass with post-bind/pre-decision review ordering and exact-one capability probe on Claude (**270.713s**, retained at `/tmp/spacedock-claude-order.XpJIha`), Codex (**241.589s**, `/tmp/spacedock-codex-order.QwQmCW`), and Pi (**196.332s**, `/tmp/spacedock-pi-order.WssMG5`).
- DONE: Exercise runtime-support recovery and reject false-green live runs.
  The installed default Pi `pi-subagents` package lacked its optional coding-agent peer. Runtime-support tracing selected the compatible FNM Pi package root `/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/lib/node_modules/pi-subagents`; the spent OpenAI OAuth path was replaced by the already-configured OpenRouter provider, with `anthropic/claude-sonnet-4.5` parent and `openrouter/anthropic/claude-haiku-4.5` worker. Runs that recomputed the wrong Briefing hash, skipped capability preflight, emitted two successful dispatch builds, presented before the durable bind, or changed the marker in two commits were retained as failures and never used to relax the oracle.
- DONE: Triage Roborev branch-final jobs 2011, 2031, 2036, and 2042 against the canonical captain rulings.
  Material findings fixed: behavioral executable-identity caching and exact-one live probes; legacy unbound-gate compatibility; open-package commit ancestry; review-after-bind/before-decision ordering; state-only Resolution counting; and the shadowed retained-path variable. Declined as correct-but-disproportionate or already-decided: generic exact-child/spawn attribution, because AC-8 assigns transport proof to host-native fixtures and explicitly requires no common spawn schema/exact-child forensics; `present-gate` redesign, because the canonical expected surface makes it xb-owned and a design-reset trigger; restored whole-host prompt ratchets, because cycle-6 captain rulings replaced them with component caps and attributable deltas; and advisory-round path normalization, which belongs to the concurrent round owner. Job 2042's Claude benchmark-token finding concerns pre-existing host-runner setup outside the declared eight-file correction and is deferred to that runtime owner; final Claude live evidence passed without changing it.

### Summary

The final implementation operates the recorded gate lifecycle with three load-bearing authority mutations, durable package/Resolution/consume barriers, one identity-keyed capability probe, and a concise review ordered between bind commit and decision. It preserves legacy human-gate holds and the captain's transport/presentation boundaries, stays under every approved LOC and byte cap, and is green under full/race/docs plus Claude, Codex, and Pi live evidence.

## Stage Report: validation

- FAILED: Reproduce all eight acceptance criteria at exact tip 624a9c2e, including three authority mutations, Git ancestry, review ordering, executable-identity cache invalidation, cwd-relative inputs, and ordinary dispatch handoff.
  AC-1 through AC-6 and AC-8 reproduce, but AC-7's cited resume proof omits the promised closed-but-uncommitted and pending-committed approval states; `TestRecordedGateLifecycleAC7ResumeMatrix` never proves one close commit precedes one consume across three resumed processes.
- DONE: Adversarially audit the 17-file surface and every approved cap/exclusion against omitted writes, PATH/same-path replacement, stale or reordered commits, duplicate dispatch, compatibility inventions, and host-transport scope drift.
  The three real-binary omission controls, capability identity controls, state-Git oracle, review extraction, legacy hold route, and exact-one successful dispatch checks red the approved mutants without adding exact-child, presenter, whole-host-ratchet, advisory-round, or second-harness obligations.
- DONE: Verify focused/full/race/docs, retained exact-tip Claude/Codex/Pi transcripts, and Roborev triage; issue a fresh PASSED or REJECTED verdict with every finding classified.
  Focused suites, `go test ./...`, `go test ./... -race`, live-tag compilation, `gofmt -d`, `git diff --check`, and `uvx ... mkdocs build --strict` pass; retained post-624a9c2e host transcripts reproduce bind-commit/review/close-commit/consume-commit/dispatch ordering. Verdict: **REJECTED**.

### Reviewer findings

- MATERIAL / EVIDENCE DEFECT / Medium — AC-7's required supported resume boundary is not reproduced.
  Trigger: a session stops after a successful approval close, either before or after the close commit. Harm: the suite can stay green if resume re-records the decision, consumes before durably committing the exact Resolution, or produces more than one close/consume/transition. Narrow correction: add the existing-fixture fresh-process controls for closed-uncommitted and pending-committed snapshots, then assert one close commit, one successful consume, one transition, byte stability, and ancestry across three passes.
- POLISH / OUTCOME DRIFT / Low — selected-state boot prose says “then `gate validate`” while the active resume rule makes diagnostics optional.
  The three authority mutations and value boundary remain intact, but `skills/fo-gate-lifecycle/SKILL.md:13` should stop implying mandatory diagnostic ceremony already excluded by Cycle 6/7.
- POLISH / DOCUMENTATION DRIFT / Low — the host-neutral scenario descriptions still narrate the superseded six-step lifecycle.
  `shared_scenarios_test.go:31` and `scenario-testing-principles.md:59` should describe three authority mutations and optional diagnostics; the executable oracle itself already does.

### Summary

Exact-tip implementation behavior, Git barriers, CLI normalization, identity invalidation controls, host handoff, caps, repository tests, strict docs, retained live evidence, and the scope-conscious Roborev dispositions all reproduce. Validation rejects only the missing AC-7 resume evidence; no ruled-out presenter redesign, generic child attribution, whole-host ratchet, advisory-round normalization, compatibility expansion, or second live harness is requested.

## Implementation intended-change declaration (cycle 8)

Validation rejected only AC-7 evidence at product tip `624a9c2e`; the branch currently has 1,461 added LOC against its merge base. The correction is limited to these four existing files:

| File | Intended cycle-8 delta |
| --- | ---: |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +65 / -27 |
| `skills/fo-gate-lifecycle/SKILL.md` | +1 / -1 |
| `internal/ensigncycle/shared_scenarios_test.go` | +1 / -1 |
| `docs/specs/scenario-testing-principles.md` | +1 / -1 |

The exact intended amendment is **+68/-30**. It replaces the existing consumed-resume subtest with a three-fresh-process approval resume proof: close and stop before commit; resume to commit the exact pending Resolution and stop; resume from committed pending state to consume and transition. The same fixture will assert one Resolution, one close commit, one successful consume, one status transition, byte-clean duplicate/repeated failures, and close → consumed → dispatch ancestry. The three prose substitutions make diagnostics optional and name the three authority mutations without adding commands or changing semantics. The resulting branch ceiling is 1,529 added LOC, 59 below the approved 1,588 cap. No new harness, presenter, host adapter/live lane, generic transport, exact-child attribution, total-host ratchet, advisory-round behavior, compatibility path, CLI, schema, or runtime surface is intended; any need outside this table stops for drift review before editing.

## Implementation intended-change declaration (cycle 13)

At exact product tip `d55ddca3`, the branch is 1,500 additions / 63 deletions against its merge base. The canonical nine-file table under “Expected surface and tolerance” is the implementation declaration: **+75/-99** intended, **+88** additions hard maximum, net deletion preferred, and no tenth file. It deletes the 41-line capability-cache experiment and 32-line helper; changes no CLI, schema, recorder, application, host runner, or runtime protocol; and reuses the existing lifecycle, Pi, merge-guard, and contract fixtures.

## Stage Report: ideation (cycle 13)

- DONE: Rewrite the canonical lifecycle, ACs, and test plan around the consumed target-kind fork.
  Three authority mutations and all package/close/consume commit barriers remain unchanged. A consumed nonterminal target enters ordinary dispatch; a consumed terminal target enters `«merge.guard»` and the existing registered hook with zero successor dispatch.
- DONE: Exercise the riskiest terminal route and identify the smallest complete recovery path.
  Existing real lifecycle/resume tests and merge-sentinel tests pass at `d55ddca3`. The live recorded sample `gate-review-presentation-command` is `done` with a consumed application targeting `done`, `mod-block: merge:pr-merge`, and PR `#564`; `status --next --json` returns `dispatchable: []`. The merge core can resume the block, but both `pr-merge` startup/idle selectors exclude terminal rows, so neither can observe the open PR or later write the existing merge sentinel. The design changes only those selectors to admit terminal rows carrying that exact in-flight block, then reuses the existing sentinel → merge-guard finalize path.
- DONE: Replace cache and ambiguous Pi evidence with smaller falsifiable obligations.
  AC-5 now requires one freshly resolved `gate --help` immediately before every lifecycle and explicitly deletes session path/digest caching plus replacement/PATH-swap experiments. AC-6/AC-8 accept Pi review evidence only from an assistant message in the single flat root session; nested subagent and nonassistant rows cannot grade green.
- DONE: Bound implementation to a precise net-deleting surface.
  Nine existing files total **+75/-99**, with an **+88** hard addition stop preserving the 1,588-LOC branch cap. The plan reuses all representative Claude/Codex/Pi journeys, derives the terminal arm from the existing lifecycle fixture, and adds no harness, transport contract, compatibility layer, terminal lifecycle, exact-child/built-artifact attribution, or recorder machinery.

### Summary

Cycle 13 resolves the dispatch/merge contradiction without changing gate authority: probe once per lifecycle, consume once, commit once, then dispatch only nonterminal targets and merge terminal targets. It also closes the actual terminal-PR resume gap and Pi evidence ambiguity while deleting the cache laboratory and remaining inside the existing fixture and host boundaries.

## Stage Report: implementation (cycle 14)

- DONE: Implement the approved recorded-gate operating contract inside the exact nine-file boundary.
  Final code tip `ac36c5576a0b0881dac9e1cf1e97728ce225506f` changes exactly **9 files, +88/-120** from `d55ddca3`, meeting the +88 hard stop and adding no tenth file. The First Officer now resolves and probes one fresh `gate --help` immediately before each lifecycle with no session executable cache; preserves bind, decision, and consume plus every required state-commit barrier; sends consumed nonterminal targets through the existing dispatch path; and sends consumed terminal targets into the existing merge ceremony with zero successor dispatch.
- DONE: Make terminal merge recovery restart-safe within the owned selector and guard boundary.
  The shipped and split-root PR-merge templates select terminal rows carrying either canonical merged sentinel (`pr-merge:` or `local-merge:`), bypass `gh`, and resume `merge guard` directly. The shipped mod advances from `0.12.2` to `0.12.3`, while the local `0.12.4` template records `reconciled-from-shipped: 0.12.3`. The existing restart control constructs terminal status, each canonical sentinel, and an empty `mod-block` before proving guard archival. Template parity pins both selectors and the reconciled rollout version.
- DONE: Tighten representative lifecycle evidence without adding a harness or transport contract.
  The shared oracle proves terminal consume has no dispatchable successor, excludes `dispatch build --help` from successful-build cardinality, and selects only state-commit transitions for Git ancestry. Pi evidence comes from exactly one flat root session and accepts the six-field review only from a root assistant message between the durable bind and decision; nested/nonassistant text cannot satisfy it. Claude, Codex, and Pi retain their existing host-native dispatch transports.
- DONE: Pass deterministic, repository, race, documentation, and representative live verification.
  `gofmt -w ./cmd ./internal`, `git diff --check`, focused status/contract/lifecycle suites, `go test ./...`, `go test ./... -race`, and `uv run --with-requirements docs/requirements.txt mkdocs build --strict` pass. Recorded-lifecycle journeys pass on Codex in **247.815s** (`/tmp/spacedock-af8-final-codex`) and Claude in **322.933s** (`/tmp/spacedock-af8-final-claude`); final tip differs from those runs only by the test-only parity/name correction. Pi passes sequentially at exact final tip in **231.682s** (`/tmp/spacedock-ac36-final-pi-sequential`). A prior concurrent Pi attempt at `/tmp/spacedock-af8-final-pi` is excluded: Codex overwrote the shared `/tmp/spacedock-dispatch/...` fixture artifact, Pi detected the stale path, and the run timed out rather than false-green.
- DONE: Request final-tip Roborev, apply the First Officer rulings, and record the complete advisory round through the command surface.
  Jobs 2145–2158 drove sentinel-first finalize prose, cleared-sentinel selection, truthful fixture preconditions, both canonical sentinels, mod-version rollout, template parity, and the corrected test name. Exact-tip job **2158** identified no correctness regression and repeated the request for a new startup/idle execution harness. The First Officer declined that request as correct-but-disproportionate under the binding no-new-harness boundary; promotion requires an approved existing hook-execution harness or live discovery/`gh`-bypass failure. The split-root final-archive remote synchronization finding is material to the prerelease but remains owned by the separate durable-decisions synchronization task, outside this nine-file entity. After path-scoped normalization preserved all 13 historical cycle entries under one heading, `gate record --round implementation/14` and `gate validate --round implementation/14` recorded the **mixed, 18-entry** triage at `review/implementation/round-14`.
- DONE: Push the code branch and leave an exact handoff.
  Commits `e9415a17`, `60adfc1f`, `6f9108d1`, `af8d73be`, and `ac36c557` are pushed to `origin/spacedock-ensign/first-officer-gate-command-lifecycle`. The code worktree is clean. Validation should retain the separate remote-archive synchronization blocker as an integration prerequisite rather than expanding this entity.

### Summary

First Officers now operate the recorded gate lifecycle with three durable authority mutations, one fresh capability probe, root-assistant-only Pi review evidence, and an explicit consumed-target fork: ordinary dispatch for nonterminal successors and the existing merge ceremony for terminal successors. Terminal merged-sentinel recovery is restart-safe for both canonical sentinel types, rollout is versioned, the implementation is exactly at its nine-file/+88 boundary, all required deterministic/full/race/docs/live evidence is green, and the complete Roborev triage is durably recorded.

## Stage Report: implementation (cycle 15)

- DONE: Verify the exact final endpoint after rebasing onto current `origin/main`.
  Final code endpoint `b99f9c664912b18e729b639e737ef203c00cacbe` is based on `origin/main` `cc51e518a3420b01fd4b455e9710d38803dc6d3e`. The focused CLI/contract/lifecycle/status suite, `gofmt -w ./cmd ./internal`, `git diff --check origin/main..HEAD`, `go test ./...`, `go test ./... -race`, and `uv run --with-requirements docs/requirements.txt mkdocs build --strict` all pass. Exact-tip recorded-gate journeys pass sequentially on Codex in **236.87s** (`/tmp/spacedock-b99-final-codex`), Claude in **291.24s** (`/tmp/spacedock-b99-final-claude-rerun`), and Pi in **207.23s** (`/tmp/spacedock-b99-final-pi-rerun`). The first Claude attempt is excluded because the host repeated `gate --help` only to print its exit code; the first Pi attempt is excluded because its visible review omitted the recorded Briefing digest. Both failures were retained, diagnosed against the strict oracle, and rerun unchanged rather than weakening evidence.
- DONE: Preserve the closed Cycle-13 boundary while identifying the final integration endpoint separately.
  The closed approved Cycle-13 range is `d7466f0c..0acfcbc6`: exactly **9 files, +88/-120**. The final endpoint range from the same predecessor is `d7466f0c..b99f9c66`: **10 files, +93/-141**. The explicit post-range integration reconciliation commit `b99f9c66` is **3 files, +6/-22**; it removes stale legacy Result/Association composition, restores the provider-neutral `--room` contract, and corrects the rebased CLI reference without rewriting the already closed range.
- DONE: Complete and record final-tip advisory review.
  Roborev job **2161** reviewed exact endpoint `b99f9c66` with correctness and product reviewers and reported **no issues found**. `gate record --round implementation/15` and `gate validate --round implementation/15` record the exact Git artifact and diff digest as a one-entry, no-findings approving round under `review/implementation/round-15`.
- DONE: Publish the rebased endpoint without rewriting the divergent historical remote branch.
  A non-force push created successor branch `origin/spacedock-ensign/first-officer-gate-command-lifecycle-rebased-cc51e518`; remote verification resolves it exactly to `b99f9c664912b18e729b639e737ef203c00cacbe`. The old remote branch was not changed.
- SKIPPED: Advance workflow state or merge the successor branch.
  The dispatch explicitly reserves advancement, validation routing, and merge to the independent First Officer. This handoff leaves the entity at implementation and returns the exact successor endpoint for validation.

### Summary

All assigned implementation outcomes account as **4 DONE, 1 SKIPPED, 0 FAILED**. The approved Cycle-13 range remains closed and separately auditable; the final integration endpoint is fully formatted, deterministic/full/race/docs green, live-green on Codex/Claude/Pi, no-findings under exact-tip Roborev, and available on a new non-force successor branch. No workflow advancement or merge was performed.

## Stage Report: validation (cycle 2)

- DONE: Reproduce AC-1 through AC-8 at exact candidate b99f9c664912b18e729b639e737ef203c00cacbe, including durable lifecycle ordering, fail-closed controls, authority provenance, retained package scope, idempotent resume, and terminal-versus-nonterminal routing.
  Fresh real-CLI, refusal, mutant, resume, terminal, merge-guard, provider-room, and three-host live controls reproduce every cited boundary at the exact remote tip.
- DONE: Perform the required detached semantic adversarial audit and check the implementation, tests, task body, and post-rebase reconciliation against the latest Captain rulings; classify every finding by defect kind and release scope.
  A detached b99f9c66 checkout red-tested user-role Pi review acceptance and loss of `local-merge:` recovery; the intentional `b99f9c66` 3-file +6/-22 reconciliation is the only post-range drift.
- DONE: Run every applicable focused/full/race/docs and existing Claude, Codex, and Pi recorded-gate live lane at the exact candidate tip, inspect the retained final-tip Roborev evidence, and issue a PASSED or REJECTED recommendation with concrete evidence.
  Focused/full/race/live-tag/docs pass; fresh Codex 224.59s, Claude 249.03s, and Pi 226.18s pass; round-15 Roborev retains exact b99f9c66 and an approving Resolution with no annotations. Recommendation: **PASSED**.

### Acceptance-criterion evidence

- AC-1: `TestRecordedGateLifecycleRealCLIReplay`, terminal consume, both-sentinel merge guard, and all three live lanes fail on missing ancestry, nonempty terminal dispatchability, duplicate build/effect, or dispatch before consume.
- AC-2: `TestRecordedGateLifecycleMissingEventControls` executes the remaining real commands for each omitted bind/decision/consume; the refusal matrix proves malformed, stale, blocked, revise, hold, and repeat-consume routes fail closed.
- AC-3: the real replay and three-process resume pin the exact close commit; actor-swap, blank-reason, altered-directive, direct-chat, and provider-room authority controls reject indistinguishable or unauthorized provenance.
- AC-4: the replay verifies canonical retained Briefing bytes, room artifacts, one folder-form state commit surface, exclusion of `dirty-sibling.md`, and unchanged workflow discovery.
- AC-5: relative/absolute Briefing and relative prepared-room CLI tests pass; whole-tree refusal controls stay byte-clean; each live command log requires exactly one successful `gate --help` before mutation.
- AC-6: the structured six-field review grader rejects raw dumps; the Pi parser accepts only one flat root-session assistant review between durable bind and decision, and its user-role mutant reds.
- AC-7: the resume matrix proves open retry, replacement, revise, hold, stale, closed-uncommitted, pending-committed, and consumed one-use behavior across three fresh processes; terminal state routes only to merge recovery.
- AC-8: one host-neutral scenario passes unchanged through fresh Claude, Codex, and Pi runners; native adapters retain their own transport while the shared oracle proves exact mutations, commits, dispatch cardinality, and durable effect.

### Reviewer findings

- No material outcome or evidence defect exists in the candidate scope.
- DEFERRED RISK / EVIDENCE DEFECT / Low — startup/idle hook selection and `gh` bypass are structurally pinned rather than executed end to end.
  Trigger: an actual hook runner diverges from the shipped selector prose. This is outside the approved no-new-harness proof boundary; direct both-sentinel guard behavior and paired template parity pass. Promote on an approved existing hook harness or observed discovery/`gh`-bypass failure.
- EXTERNAL INTEGRATION PREREQUISITE / OUTCOME DEFECT / material to prerelease — split-root merge finalization still needs the separately owned conflict-safe remote archive synchronization.
  This endpoint promises and proves the local path-scoped commit, not remote publication; it does not block this entity gate, but remains a release blocker if the synchronization owner does not land before prerelease.

### Summary

Exact endpoint `b99f9c664912b18e729b639e737ef203c00cacbe` preserves the closed 9-file Cycle-13 boundary and adds only the Captain-approved provider-room reconciliation. All ACs have independently reproduced behavioral evidence, the detached mutants red at the intended boundaries, repository and live lanes pass, and validation recommends **PASSED** with no candidate-scope material finding.

## Stage Report: implementation (cycle 16)

- DONE: Make the existing recorded-gate fixture self-contained by persisting repository-local test Git identity in its temporary state checkout, with no product, CI, global-config, or unrelated-fixture change.
  Commit `13d702492131df17dd3ac87245d6d773f4df959b` adds only `git config user.email t@t` and `git config user.name t` after the fixture initializes its state checkout: **1 test file, +2/-0** from `b99f9c66`; removing either identity leaves product-owned commits dependent on the host.
- DONE: Prove the exact isolated clean-config reproduction changes from red to green, then run the affected package plus required full/race/format/diff verification and diagnose any failure from this run.
  With `[user] useConfigOnly=true` and no identity, `TestRecordedGateLifecycleRealCLIReplay` first failed at the initial `spacedock state commit` with `Author identity unknown`; after the fix, all CI-failing replay/terminal/refusal/resume cases pass, the clean-config `internal/ensigncycle` package passes in **14.421s**, and `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` pass.
- DONE: Request final-tip Roborev, record exact correction surface and DONE/SKIPPED/FAILED evidence in the implementation Stage Report, and fast-forward the existing PR successor branch without advancing workflow state.
  Branch-final panel job **2167** reviewed exact range `cc51e518..13d70249`: product passed with no issues; correctness repeated four branch-wide medium findings outside `b99f9c66..13d70249`. Rejection routing and captain-decision mapping require forbidden lifecycle/schema expansion; goal-only prompting and Claude root-row attribution require proof-surface redesign outside this setup correction. They remain worker-authored declines pending independent validation, with one-to-one promotion conditions in the **10-entry, all-declines** advisory package recorded and validated at `review/implementation/round-16`; no First Officer approval is claimed. The successor branch was fast-forwarded and remotely resolves exactly to `13d702492131df17dd3ac87245d6d773f4df959b`.
- SKIPPED: Rebuild strict documentation.
  The final correction changes only Go test fixture setup; no documentation, generated site, command text, or user-visible contract changed.
- SKIPPED: Advance workflow state, mutate gate records or PR metadata, or merge.
  The routed assignment reserves fresh validation and replacement-gate handling to the First Officer; this implementation writes only its body report.

### Summary

The CI-only red was an undeclared host-identity dependency in test setup, not a product lifecycle failure. The one-file/two-line correction persists the fixture helper's existing identity for independent product Git subprocesses, is clean-config/full/race green, is published at `13d70249`, and now has a durable implementation/16 advisory round carrying the Cycle 15 correction context; checklist accounting remains **3 DONE, 2 SKIPPED, 0 FAILED**.

## Stage Report: validation (cycle 3)

- DONE: At exact tip 13d702492131df17dd3ac87245d6d773f4df959b, reproduce the identity-free Git failure turning red-to-green, verify the correction is exactly one test file +2/-0 with no product drift, and run the applicable focused/full/race/format evidence.
  With isolated `user.useConfigOnly=true`, b99f9c66 fails its first fixture commit with `Author identity unknown` and 13d70249 passes; the exact delta is one test file +2/-0, and focused/full/race/format/diff/live-tag compilation pass.
- DONE: Independently assess every Roborev 2167 finding preserved in advisory round implementation/16 against the current task ACs and supported First Officer workflow; reproduce the cheapest falsifiable control for each disputed materiality claim and do not accept “predates the two-line fix” as a release-scope argument.
  Detached controls exercised goal-only rejection, unmapped reject/redo decisions, a deleted shipped close-commit barrier, and child-only Claude review evidence; all four findings were assessed on the complete candidate.
- FAILED: Re-anchor AC-1 through AC-8 on the final tip and issue PASSED or REJECTED with each finding classified by defect kind and release scope; do not repair code, narrow ACs, mutate gate/PR state, or merge.
  AC-1, AC-6, AC-7, and AC-8 retain material outcome/evidence defects below; AC-2 through AC-5 reproduce. Recommendation: **REJECTED — DESIGN RESET**.

### Acceptance-criterion re-anchor

- AC-1: REJECTED evidence — deleting the shipped close-commit barrier leaves focused tests green and the Codex recorded-gate live lane passes in 211.37s because its prompt independently scripts every barrier.
- AC-2: PASSED — real-binary omitted bind/decision/consume and the malformed/stale/blocked/revise/hold/repeat-consume refusal matrix remain green.
- AC-3: PASSED — direct/delegated/provider authority, actor-swap, exact directive, close durability, and stopped-route controls remain green.
- AC-4: PASSED — canonical retained package digests, folder commit scope, dirty-sibling exclusion, and workflow-discovery equality remain green.
- AC-5: PASSED — isolated identity repair, one fresh capability probe, relative retained inputs, byte-clean failures, and no lock residue reproduce.
- AC-6: REJECTED outcome/evidence — captain-facing `reject`/`redo` have no recorder mapping, and a Claude child-only six-field review is accepted as if the root presented it.
- AC-7: REJECTED outcome — existing approval resume is idempotent, but an ordinary captain rejection cannot enter a defined persisted decision/resume route without model-invented translation.
- AC-8: REJECTED evidence — the host-neutral live prompt owns the procedure under test, and Claude extraction cannot establish root captain visibility; the three-host claim is therefore not validly proven.

### Reviewer findings

- MATERIAL / OUTCOME DEFECT / Medium — automatic reviewer-REJECTED routing was removed from the shared core.
  Trigger: a supported feedback gate recommends REJECTED. The pre-branch contract routed it automatically; the candidate leaves routing to inference. A goal-only Codex run happened to load feedback routing but then skipped the initial implementation report while claiming success, so the supported journey is not contractually reliable.
- MATERIAL / OUTCOME DEFECT / Medium — captain decision vocabulary has no durable mapping.
  Trigger: the captain selects presented `reject` or `redo`; both real CLI controls refuse byte-clean because the recorder accepts only approve/revise/hold. This blocks a common supported gate route and leaves rejection versus redo semantics undefined.
- MATERIAL / EVIDENCE DEFECT / Medium — the live prompt duplicates lifecycle procedure.
  Adversarial edit: remove the shipped successful-close commit barrier while retaining command tokens. Focused lifecycle tests and the real Codex live grader both remain green, so AC-1/AC-8 can pass when the shipped contract regresses.
- MATERIAL / EVIDENCE DEFECT / Medium — Claude review extraction accepts subagent-only output.
  Adversarial stream: root binds, only a child row emits the six-field review, then root decides. The extractor returns the child review because row `parent_tool_use_id` is discarded, invalidating AC-6/AC-8 visibility evidence.
- RULING: design reset, not another automatic narrow correction.
  Restore the supported rejection branch and define captain-to-recorder semantics, then rebuild the live lane around a goal-only prompt and root-aware Claude rows. This crosses outcome and evidence layers/controllers, meeting the validation policy’s reset condition.

### Summary

The two-line Git-identity correction is correct, minimal, and green, but it only exposes a clean harness for a branch whose broader supported behavior and proof boundary remain unsound. Validation rejects exact tip `13d702492131df17dd3ac87245d6d773f4df959b` and recommends a design reset covering rejection routing, decision semantics, goal-only live ownership, and root-authoritative Claude evidence.

## Stage Report: ideation (cycle 17)

- DONE: Define the smallest coherent contract that restores automatic reviewer-`REJECTED` routing and maps ordinary Captain gate language to the existing recorder schema.
  A completed feedback-gate reviewer `REJECTED` recommendation enters the existing feedback-rejection flow before Captain presentation and creates no Captain Resolution, consume, or successor dispatch. Captain `approve` maps to `approve`; `redo with feedback` maps to `revise` with an accepts-direction reason; `reject` with `feedback-to` maps to `revise` with a rejects-direction reason; and `reject` without a correction owner, `hold`, or `not yet` maps to `hold`. The durable reason preserves redo versus reject without a new enum, field, product mapper, or parallel lifecycle.
- DONE: Redesign the host-neutral evidence so shipped skills own procedure and each validation counterexample fails at its real boundary.
  The existing recorded-gate and rejection prompts retain only fixture, authority, goal, and stop-marker constraints. Focused deletion/swap controls cover the shared-core rejection branch, successful-close commit barrier, Captain-word mapping, and absence of procedure from prompt text. The existing Claude extractor consumes the already-present top-level `parent_tool_use_id` and rejects child rows; Pi retains its root-session assistant rule. Approved live behavior remains the same exact three-mutation/Git-ancestry/dispatch oracle on Claude, Codex, and Pi.
- DONE: Update the canonical task body, acceptance criteria, test plan, documentation diff, and implementation surface together, preserving every previously proven boundary.
  From current tip `13d70249`, implementation is seven existing files, intended **+75/-16**, with a **+95** additions hard stop and no eighth file. The shared core stays at or below 26,754 bytes, the lifecycle skill at or below 6,600 bytes, and the repo-local fixture identity fix remains. Recorder/schema/application/CLI production code, provider distinction, retained rooms, one-use consume, nonterminal dispatch, terminal merge recovery, and host adapters are unchanged.

### Risk spike and reset trigger

No implementation spike is needed. The existing shared-core history already contains the compact automatic feedback-rejection pre-branch, `fo-gate-lifecycle` already owns canonical `approve|revise|hold` recording and fail-closed routing, the rejection fixture already supplies the two-cycle correction journey, and Claude JSONL already exposes a top-level parent marker. The validation mutants proved the current evidence hole rather than a missing product primitive.

Deterministic controls must pass before live spend. After one diagnosed goal-only failure, each host gets at most one unchanged rerun. The same missing shipped lifecycle obligation twice on any supported host is a hard design-reset trigger: stop rather than add prompt procedure, weaken the oracle, or exclude the host.

### Summary

Cycle 17 keeps one recorder vocabulary and one First Officer lifecycle. It restores the existing automatic correction owner, translates Captain language durably through `Resolution.reason`, moves procedure ownership out of test prompts and back into shipped skills, and makes Claude root visibility falsifiable. The design is ready for implementation but this stage makes no code, PR, gate, workflow-status, or merge change.

## Stage Report: ideation (cycle 18)

- DONE: Make the deterministic structural deletion mutant authoritative for automatic reviewer-`REJECTED` branch ownership and retain the goal-only rejection journey only as positive integration proof.
  The canonical AC, friction inventory, test plan, obligation table, and hard-reset rule now state that model-inferred live bounce neither proves ownership nor triggers another prompt/controller/harness when the structural mutant already reds.
- DONE: Add an independently falsifiable `not yet` to canonical `hold` mapping.
  The mapping table, AC-6, AC-7, and test plan require deletion of only the alias to fail; the positive snapshot must retain a nonblank pause reason, resume held, and produce zero consume or dispatch.
- DONE: Qualify the user documentation and correct the seven-file arithmetic.
  The concrete concepts diff limits automatic bounce to reviewer `REJECTED` at a configured feedback gate and maps ordinary Captain reject without `feedback-to` to hold. Its exact amendment is +9/-2; the seven-file total is **+77/-16**, intended branch additions are 1,549, and the +95 hard stop preserves exactly 21 lines below the 1,588 cap.
- SKIPPED: Modify code, PR #565, gate records, workflow status, or merge state.
  The independent Cycle-17 staff review requested ideation-only corrections; the code worktree remains at exact tip `13d70249`.

### Summary

Independent staff review disposition **CHANGES REQUESTED** is fully accepted: both material corrections and the documentation/surface notes are incorporated without weakening Claude root-row scope or any no-schema/controller/harness boundary. Checklist accounting is **3 DONE, 1 SKIPPED, 0 FAILED**; this cycle changes only the canonical entity design and report.

## Stage Report: implementation (cycle 19)

- DONE: Implement the approved rejection and Captain-language contract in the existing shared core and lifecycle skill, preserving one recorder schema and all durable bind/close/consume barriers.
  Commits `c4afd8f7` and `37d6980b` change exactly seven approved files, **+95/-14** from `13d70249`: automatic reviewer rejection has one shared-core owner, Captain calls map through `approve|revise|hold`, prompts are goal-only, Claude filters non-root rows, docs are qualified, and the identity fix remains.
- FAILED: Make the existing prompts goal-only, structural/mapping mutants load-bearing, and Claude review extraction top-level-only; prove all four validation counterexamples red at their real boundaries without a new harness, controller, adapter, or procedure hint.
  Focused controls red before implementation and pass afterward, including close-barrier deletion, reviewer-route deletion, every mapping plus `not yet` deletion/held resume, prompt-procedure exclusion, and child-only Claude rejection. But two unchanged goal-only Codex runs completed bind/close/consume/dispatch while omitting the root-visible six-field review; the second loaded the corrected label-bearing skill and failed the same oracle in 190.29s, firing the approved hard reset.
- FAILED: Stay within the exact seven existing files and +95-addition hard stop, apply the qualified documentation diff, run all required focused/full/race/docs and Claude/Codex/Pi live lanes, request final-tip Roborev, and durably record its advisory round.
  Surface and docs are exact; affected/full/race/strict-docs/live-compile passed at `c4afd8f7`, and final `37d6980b` passed focused/format/diff checks. The first Codex rejection run reached corrected Cycle 2 validation before the combined test hit its 10-minute ceiling, but the repeated recorded-gate failure required stopping before Claude/Pi, final-tip full/race/docs, Roborev, or an advisory round.

### Summary

Checklist accounting is **1 DONE, 0 SKIPPED, 2 FAILED**. The design’s hard-reset trigger worked as intended: skill text plus a goal-only prompt does not reliably produce the required root-visible review on Codex, so implementation stops without prompt coaching, oracle weakening, a new controller/harness, PR #565 changes, gate/status mutation, or merge; the First Officer must re-ideate skill-discovery/runtime ownership.

## Stage Report: ideation (cycle 20)

- DONE: Identify why two goal-only Codex runs omitted the root-visible review after loading the corrected lifecycle.
  The retained first artifact at `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/TestLiveCodexSharedScenarios3309482735/001/recorded-gate-lifecycle/codex-exec.jsonl` proves the FO understood the outcome and retained authority sequence: items 14-15 bind/commit, item 19 judges and recommends approval, item 20 loads `present-gate`, and item 21 immediately records the delegated decision. No root `agent_message` contains `Capability/change:`; the oracle reports count 0. The second recorded-only run failed the identical field/count after 190.29 seconds at the corrected `37d6980b` tip. Its temporary directory was cleaned after test exit, but the exact command, failure, duration, and path are preserved in Cycle 19. The cause is not missing natural-language end value or missing skill discovery: `fo-gate-lifecycle` and `present-gate` were both loaded. It is conflicting render ownership plus no observable presenter-completion boundary.
- DONE: Design the smallest sound ownership boundary that preserves the strict oracle and existing lifecycle.
  `present-gate` becomes the sole owner of the six ordered review fields. On default chat it completes only after one root-assistant review is emitted after bind commit and before decision record; delegated conn may automate the decision only after that same artifact. `fo-gate-lifecycle` owns bind/close/consume and waits for presentation completion instead of duplicating the template. Provider override completion remains its existing retained Result boundary. This adds no command, schema, recorder, controller, harness, compatibility path, host copy, or prompt coaching.
- DONE: Specify the revised surface and falsifiable proof plan for independent staff review.
  The failed seven-file tip remains local and off PR #565. The reset adds only existing `skills/present-gate/SKILL.md`, yielding eight files with an incremental **+61/-31** intent and a total **+165** additions hard stop from `13d70249`. Existing contract and transcript tests make sole ownership, exact template, and bind-review-decision order deletion-sensitive; the unchanged goal-only Codex lane runs first and any third identical omission stops immediately. Only after it passes do Claude, Pi, full/race/docs, and advisory review proceed. The retained JSONL proves Codex can emit root messages between tools, so no new runtime mechanism needs a throwaway spike.
- SKIPPED: Modify code, PR #565, gate records, workflow status, or merge state.
  This hard-reset stage is design-only. The local failed implementation commits `c4afd8f7` and `37d6980b` remain unpushed for the next implementation worker to repair or rework after independent staff review.

### Summary

Checklist accounting is **3 DONE, 1 SKIPPED, 0 FAILED**. Cycle 20 relocates the missing Captain-visible artifact to the skill that actually renders it, removes the conflicting dual template, and makes presentation completion observable before decision mutation while preserving every recorder/consumer authority barrier and the unchanged root-visible-review oracle.

## Stage Report: ideation (cycle 20 repair)

- DONE: Make exact-one root review cardinality deletion-sensitive without changing the harness or oracle.
  AC-6, the baseline/runtime test plan, and the obligation table now require bind commit → two qualifying six-field root `agent_message` rows → decision record to expose count 2 and fail. Extractors must preserve all qualifying interval matches or return an explicit multiplicity error; they may not overwrite with the last. Deleting either fixture message produces the exact-one positive, while deleting the multiplicity rejection makes the duplicate control fail. The control stays in `internal/ensigncycle/recorded_gate_lifecycle_test.go`.
- DONE: Preserve the decisive Codex failure evidence independently of its temporary artifact path.
  The canonical mechanism-evidence section now contains a minimal ordered excerpt: Briefing bind, state commit, root judgment/recommendation, presenter skill load, immediate delegated decision record, zero qualifying six-field root renders in the interval, and the exact count-0 oracle failure. It explicitly remains evidence rather than a new fixture or replay source.
- DONE: Reconcile the expected implementation budget with the staff correction.
  The same eight existing files remain the complete surface. The table-driven transcript controls increase the recorded-gate test estimate by four lines, making the incremental intent **+65/-31**, the total intent **+160/-45** from `13d70249`, and the total additions hard stop **+175**. No ninth file, new harness, prompt procedure, compatibility path, or oracle weakening is allowed.
- SKIPPED: Modify implementation commits, PR #565, gate records, workflow status, or merge state.
  This repair changes only the durable ideation design and report; the failed code tip remains local and unpushed.

### Summary

Independent staff **CHANGES REQUESTED** is fully incorporated. Checklist accounting is **3 DONE, 1 SKIPPED, 0 FAILED**: exact-one review cardinality is now falsifiable, the decisive runtime evidence survives temp cleanup, and the corrected budget preserves the same ownership boundary and existing test surface.

## Implementation intended-change declaration (cycle 21)

Declared before code edits at local failed tip `37d6980b`. This cycle repairs that tip in place and keeps it off PR #565 until validation. The complete surface is the approved eight existing files: the prior seven plus `skills/present-gate/SKILL.md`.

| File | Intended incremental delta |
| --- | ---: |
| `skills/first-officer/references/first-officer-shared-core.md` | 0 / 0 |
| `skills/fo-gate-lifecycle/SKILL.md` | +3 / -2 |
| `skills/present-gate/SKILL.md` | +16 / -14 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +22 / -0 |
| `internal/ensigncycle/shared_fixtures_test.go` | 0 / 0 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +16 / -0 |
| `docs/site/concepts/gates-and-decisions.md` | +8 / -15 |
| `docs/specs/scenario-testing-principles.md` | 0 / 0 |

Intended incremental total is **+65/-31**; total additions from baseline `13d70249` have a hard stop at **+175**. `present-gate` will be the sole six-field root-visible render owner; `fo-gate-lifecycle` will preserve bind/close/consume and require completed presentation before chat decision mutation. Existing host extractors will expose exact-one multiplicity without a new harness. Provider override behavior, prompts, recorder/schema/application/CLI production code, host adapters, compatibility behavior, and every ninth file remain unchanged.

## Stage Report: implementation (cycle 21)

- DONE: Implement the approved single-owner presentation boundary and exact-one root-review controls.
  Commits `7419e61a` and `08675f02` repair the local failed tip without pushing PR #565. `present-gate` solely owns the six ordered fields and requires one root-assistant chat render before decision mutation; `fo-gate-lifecycle` waits for presentation completion while retaining bind/close/consume and provider override behavior. Claude, Codex, and Pi extractors now preserve every qualifying interval review so duplicates fail exact-one grading. Structural and three-host duplicate controls failed before implementation and passed afterward. Exact tip `08675f02ff278bb16965ce583f00b15bd724e81f` changes the approved eight files by **+174/-64** from `13d70249` (**+81/-52** from failed tip `37d6980b`), one line below the +175 additions hard stop. Shared core is 26,737 bytes, lifecycle is 6,384 bytes, and presenter is 6,779 bytes, all within cap.
- DONE: Run deterministic, repository, documentation, compile, and decisive Codex proof.
  Focused ownership/multiplicity controls, full `internal/contractlint` plus `internal/ensigncycle`, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, strict MkDocs through an isolated environment built from `docs/requirements.txt`, and `go test -tags=live ./... -run '^$'` all pass. The unchanged goal-only Codex recorded-gate journey passes in 190.91 seconds, proving the prior repeated root-review omission is corrected without prompt labels, procedure coaching, a new harness, or oracle weakening.
- FAILED: Complete green Claude and Pi live evidence on the unchanged supported journeys.
  Claude run 1 emitted the six-field review and completed bind/close/consume/dispatch, but failed after 216.03 seconds because the review omitted the exact retained Briefing digest (`gate review does not name the exact retained snapshot identity and digest`). Its one unchanged retry again presented and completed the lifecycle, then failed after 234.38 seconds because durable state omitted the exact delegated directive. These are distinct from each other and from the corrected Codex omission; the design permits no third Claude run. Pi failed before scenario launch in 0.98 seconds because installed `pi-subagents` could not resolve `@earendil-works/pi-coding-agent`; the documented isolated-home and explicit-extension wiring was already active, and no installed copy of the missing dependency exists.
- FAILED: Clear and canonically record the final-tip Roborev advisory round.
  `branch_final` job 2170 reviewed the exact four-commit range and returned six medium findings. Four are material or need design authority: supported gate-guardrail assertions still require retired `Gate review:` / `Decision:` markers while Cycle 20 forbids compatibility; unbound legacy gates have no truthful bound-Briefing snapshot value; transcript ordering trusts narrated `Committed` text rather than successful state-commit events at a 174/175 cap; and the primary onboarding example contradicts the new format but is a forbidden ninth file. Reviewer-`REJECTED` lifecycle-load ordering is a deferred contract risk. The requested Captain-phrase runtime matrix is declined as correct-but-disproportionate because deletion-sensitive mapping and durable route controls already own it, with promotion on an observed wrong goal-only mapping. The round recorder refused the honest unresolved triage with `worker triage contains an invalid disposition`: it accepts only material-fixed or canonical non-material decline dispositions. The complete Briefing/log and refusal are retained under `review/implementation/roborev-2170/`; no round pointer or entity mutation was invented.
- SKIPPED: Fix findings across prohibited boundaries or mutate external workflow state.
  Under the Captain’s sprint conn, the First Officer ruled design reset. No compatibility marker, ninth file, new parser/harness, global Pi package mutation, PR push, gate record, workflow status change, or merge action was made. The four code commits remain local and unpushed for re-ideation.

### Summary

Checklist accounting is **2 DONE, 1 SKIPPED, 2 FAILED**. The cycle-20 presenter correction succeeds on the decisive unchanged Codex journey and all deterministic checks, but supported-route compatibility, legacy presentation, commit-event attribution, onboarding scope, and non-Codex live evidence prevent fresh validation. Implementation stops at exact tip `08675f02` and returns the ticket to design-reset ideation with complete retained evidence.

## Canonical v1 design reset (cycle 22; Feedback Cycle 19 authority)

This section is the sole active delta over the canonical task body above. It replaces the prototype presentation assumptions in the Problem, Operational lifecycle, AC-1, AC-3, AC-6, AC-8, Minimum replay and test plan, Obligation delta, Expected surface, Documentation change proposal, Dependencies, Minimum value seed, and Boundary. AC-2, AC-4, AC-5, and AC-7 remain active except where this section strengthens their evidence boundary. Historical reports remain evidence, not competing authority.

### Problem and evidence boundary

Roborev job 2170 exposed a v1 contradiction rather than a compatibility requirement. The shared `gate-guardrail` journey and its Claude/Codex assertions still require the prototype-only `Gate review:` and `Decision:` markers and an entity that is presented without a recorded Briefing. The v1 presenter instead requires an exact bound Briefing identity and digest. Because v1 is unreleased, retaining both shapes would manufacture a compatibility layer and make the snapshot field false on the unbound path.

The truthful v1 invariant is **bind before present**. A gate-stage entity with no selected attempt is `validating`, not presentation-ready. The FO retains a canonical Briefing, records the open binding, commits the package, and only then presents the six-field review. If no Captain decision or delegated conn exists, it stops with that open attempt; it does not close, consume, advance, or dispatch. The existing `gate-guardrail`, headless gate-stop, and live-scenario primitive journeys are refitted to this bind-and-hold outcome. Their old markers and byte-identical-entity expectation are deleted, not accepted as aliases.

Cycle 21 retained only the exact Claude failure statements, not the full streams. The first Claude run completed the lifecycle but its review omitted the exact bound digest. The unchanged retry completed the lifecycle but durable state omitted the exact delegated directive. No evidence identifies the omitted command argument, actor spelling, or intermediate prose, so this design makes no claim about them. The supported diagnosis is two independent data-flow losses: the bound snapshot tuple did not survive into the rendered review, and the delegated authority token did not survive into the durable Resolution. The next Claude run must use a stable `SPACEDOCK_LIVE_ARTIFACT_DIR` and retain stream, final message, command log, before/after entity, and state Git history on pass or fail.

### Proposed v1 contract

1. **One bind-before-present path.** `first-officer-shared-core.md` and `fo-gate-lifecycle` remove “unbound legacy gate” behavior. The gate lifecycle owns `validating -> open/committed -> presented`; absence of a supplied package means the FO assembles and retains one from current stage evidence, not that it presents without binding.
2. **Immutable presentation and authority inputs.** After the successful bind commit, the FO rereads the selected record and captures the exact `(Briefing.id, Briefing.digest)` pair. `present-gate` renders both values byte-for-byte in `Reviewed snapshot:` and blocks decision mutation if either is absent. At engage, an explicit conn is captured once as an opaque directive token; delegated close passes that exact token to `--directive`, never a summary, reconstructed grant, review text, or reason. These are shipped-skill obligations; the goal-only live prompt still carries only fixture, authority input, goal, and stop marker.
3. **Successful-commit attribution, not narration.** The existing logging wrapper remains the only command substrate. It already records `begin`, `exit=0`, and `state-head` for successful `state commit`. Each host extractor pairs the structured state-commit command occurrence in its existing stream/session with the corresponding successful wrapper-log completion and state HEAD. Only a later root review and an eventual structured decision command form the review interval. Assistant narration containing `Committed`, a failed commit, text before the commit tool call, child output, and tool/skill output cannot open the interval. Codex additionally uses its existing `item.completed` command event; Claude correlates the existing Bash `tool_use` with the later successful wrapper entry; Pi correlates its root-session `toolCall` with the same wrapper entry. No new harness, recorder field, timestamp protocol, or host lifecycle is introduced.
4. **Feedback rejection stays outside the Captain gate lifecycle.** A completed feedback reviewer `REJECTED` result is owned by `feedback-rejection-flow`, before Captain presentation. It loads that owner and routes correction without binding, Resolution, consume, or successor dispatch. `fo-gate-lifecycle` is intentionally not a precondition because no Captain gate is entered. The shared-core contract test must pin this branch before the ordinary lifecycle load, resolving job 2170's ordering question without duplicating lifecycle rules.
5. **Provider and decision behavior stay settled.** Presentation overrides retain their blocking room/Result boundary. Captain phrase mappings remain deterministic and deletion-sensitive; no route-by-host live matrix is added absent an observed wrong goal-only mapping. Recorder, schema, application, provider, consume, resume, terminal merge, and one-use semantics do not change.

The immutable tuple/token mechanism serves AC-3 and AC-6. The simpler alternative—repeat “exact” in the live prompt—is insufficient because it coaches the oracle and leaves shipped skills unproven. The structured-command/log pairing serves AC-1, AC-6, and AC-8. The simpler alternative—keep splitting on `Committed`—accepts narration; a second harness or recorder event is unnecessary because the existing stream and wrapper already expose command identity, exit, and committed HEAD. Refitting the guardrail journey serves AC-1 and AC-6. Keeping an unbound fallback is insufficient because no truthful digest exists; deleting all no-conn proof would lose the human-owned stop outcome.

### Acceptance-criterion delta

**AC-1 (VALUE, replacement) — Every v1 presentation is backed by one durable binding and every outcome reaches only its owner.** Across the refitted no-conn guardrail, delegated approved lifecycle, and terminal/rejection controls, the measured count of root reviews before a successful open-binding commit is **0**, decisions without a prior successful binding commit is **0**, advances/dispatches without a consumed descendant commit is **0**, and no-conn closes/consumes/dispatches is **0**. A no-attempt gate becomes open and committed before its one review, then stays at the gated stage. Reviewer `REJECTED` instead enters its feedback owner with zero Captain lifecycle mutations. *Verified by:* real-CLI bind-and-hold fixtures, structured command/log ordering controls, unchanged approved-route durable-state oracle, terminal and rejection controls, the refitted Claude/Codex guardrail journey, and the headless Claude gate-stop journey.

**AC-3 (replacement) — Direct and delegated authority are durably distinguishable and delegated conn is lossless.** Direct Captain close records `person:captain`. Delegated close records `agent:first-officer`, a nonblank evidence reason, and exactly one byte-for-byte occurrence of the conn captured at engage. Any missing, changed, duplicated, or substituted directive refuses the provenance grade even if bind, review, consume, and dispatch otherwise succeed. Every successful close still has its commit barrier. *Verified by:* public CLI fixtures, exact-token keep/drop/mutate/substitute controls over the lifecycle command contract, actor-swap control, real-Git stopped-route controls, and each supported-host approved journey's durable entity.

**AC-6 (replacement) — The Captain receives exactly one truthful six-field root review after durable bind and before chat decision.** The sole review uses the ordered nonblank fields already defined by `present-gate`; `Reviewed snapshot:` contains the exact captured Briefing id and digest. Prototype `Gate review:` and `Decision:` markers have no contractual meaning. The refitted no-conn journey binds and commits, emits the same review, and stops open. Two qualifying root reviews, a review with either tuple member changed/omitted, or a review bracketed only by narration/failed commit fails. *Verified by:* presenter ownership tests, id/digest keep/drop/mutate controls, host transcript controls using successful wrapper-log completion, exact-one multiplicity controls, and unchanged goal-only Codex/Claude/Pi approved journeys.

**AC-8 (replacement) — Shipped skills own one runtime-portable v1 lifecycle and supported-host evidence remains required.** Goal-only prompts contain no gate command sequence or review labels. Codex, Claude, and Pi each prove the approved lifecycle with the same durable oracle; Codex's Cycle-21 success is preserved. Claude must pass both exact tuple and exact directive obligations with full retained artifacts. Pi must pass after selecting a package root whose `pi-subagents` imports resolve. A dependency failure before the Pi scenario is runner provisioning, not product behavior and not a waiver of the Pi lane. *Verified by:* prompt exclusion, structural owner/load-order controls, deterministic three-host stream/session controls, and one retained live approved journey per supported host.

AC-2's minimum three mutations, AC-4's package durability, AC-5's capability/path boundary, and AC-7's retry/meaning rules are unchanged. Captain-language runtime expansion remains correct-but-disproportionate: existing delete/swap controls plus durable routed/held snapshots own it, promoting only if a goal-only Captain phrase records or routes incorrectly.

### Test plan and reset conditions

1. **Deterministic bind-and-hold (medium; AC-1/AC-6):** refit the existing guardrail fixture and assertions, do not add a scenario. Execute the real binary through the existing wrapper. Assert one successful help, bind, and state commit; the committed HEAD contains the exact open binding/package; one six-field review follows; decision, consume, status advance, dispatch, verdict, completion, and archive counts remain zero. Delete the bind, successful commit, or review and the control must red.
2. **Commit-attribution controls (medium; AC-1/AC-6/AC-8):** table-test the existing Claude, Codex, and Pi stream/session shapes paired with wrapper logs. Narrated `Committed` plus pre-commit review fails; failed commit plus review fails; successful matching commit then review passes; duplicate root reviews fail. Change the logged exit or `state-head`, structured command, root role/parent, tuple value, or decision boundary and the grade must fail.
3. **Immutable tuple/directive controls (low; AC-3/AC-6):** contract tests delete or mutate each captured Briefing id/digest and exact directive handoff independently. Run the real CLI with the resulting exact arguments and assert the durable entity, not prompt substrings. The tuple control fails before close; directive mutation fails provenance even when the rest of the lifecycle is valid.
4. **No-conn and rejection ownership (medium; AC-1/AC-8):** the refitted shared guardrail and headless Claude drive prove bind-and-hold. A structural control pins feedback reviewer `REJECTED -> feedback-rejection-flow` before ordinary `gate.lifecycle` and proves zero Captain mutations; the existing positive rejection journey remains integration evidence.
5. **Live lanes (high; AC-1/AC-3/AC-6/AC-8):** run deterministic/full/race/docs/live-compile first. Preserve the unchanged Codex pass. Run Claude once with stable retained artifacts; one unchanged retry is allowed only for a distinct obligation. For Pi, first set `PI_SUBAGENTS_PACKAGE_ROOT` to a read-only verified compatible installation (the current FNM root contains the older package and its `@mariozechner/pi-coding-agent` peer); keep isolated home/auth and explicit extension/skill wiring, then run the unchanged journey. No global install is part of 6y.
6. **Repository and docs gates (low):** run `gofmt -w ./cmd ./internal`, focused contract/lifecycle tests, `go test ./...`, `go test ./... -race`, live-tag compilation, and strict MkDocs. Run final-tip Roborev only after all required host evidence is retained.

No throwaway code spike is needed. The retained Codex stream proves structured completed command events with exit code; the existing Claude stream parser already correlates Bash tool-use ids with tool results; the existing Pi root-session parser consumes assistant `toolCall` blocks; and the wrapper already records successful exits plus state HEADs. Implementation first combines those proven inputs in deterministic controls before live spend.

Hard reset before further editing if successful-commit attribution requires a new recorder field, second wrapper/harness, timestamp protocol, or host-specific lifecycle; if a gate must be presented without a canonical bound Briefing; if goal-only prompts gain procedure/labels; if exact tuple/directive grading is weakened; if supported Pi is dropped rather than provisioned; if any eighteenth file or the LOC cap below is needed; or if the same shipped obligation fails twice unchanged on a supported host.

### Expected surface and tolerance

Implementation starts from local unpushed tip `08675f02ff278bb16965ce583f00b15bd724e81f`. It may refit the existing guardrail tests beyond the former eight-file boundary, but the complete surface is these **17 existing files**:

| File | Intended incremental delta |
| --- | ---: |
| `skills/first-officer/references/first-officer-shared-core.md` | +3 / -3 |
| `skills/fo-gate-lifecycle/SKILL.md` | +12 / -6 |
| `skills/present-gate/SKILL.md` | +12 / -5 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +90 / -25 |
| `internal/ensigncycle/shared_fixtures_test.go` | +20 / -25 |
| `internal/ensigncycle/gate_assert_impl_test.go` | +22 / -15 |
| `internal/ensigncycle/gate_assert_test.go` | +12 / -8 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +10 / -6 |
| `internal/ensigncycle/shared_scenarios_test.go` | +2 / -2 |
| `internal/ensigncycle/claude_live_runner_test.go` | +10 / -5 |
| `internal/ensigncycle/codex_live_runner_test.go` | +10 / -5 |
| `internal/ensigncycle/live_gate_stop_test.go` | +35 / -20 |
| `internal/ensigncycle/livescenario_adapter_live_test.go` | +15 / -12 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +22 / -4 |
| `docs/site/concepts/gates-and-decisions.md` | +8 / -4 |
| `docs/site/get-started/first-workflow.md` | +10 / -13 |
| `docs/specs/scenario-testing-principles.md` | +4 / -2 |

Intended incremental total is **+297/-160**. The full branch intent versus baseline `13d70249` is **+471/-224**; additions have a hard stop at **+510** versus that baseline. Shared core remains at or below **26,754 bytes**, lifecycle at or below **7,000 bytes**, and presenter at or below **7,500 bytes**. Test refactoring may move lines among the listed `internal/ensigncycle` files, but neither an eighteenth file nor the full-addition cap may be crossed without re-ideation.

`recorded_gate_lifecycle_test.go` owns the common structured-command/log interval and three-host deterministic controls. The existing host runners only supply their native stream and wrapper log; they do not gain lifecycle procedure. The guardrail fixture remains one shared scenario but changes its durable outcome from byte-identical unbound presentation to committed-open bind-and-hold.

### Documentation delta

In `docs/site/concepts/gates-and-decisions.md`, replace the ambiguous gate-entry paragraph:

```diff
-After completion verification, the first officer binds the retained Briefing
-before presenting the gate.
+After completion verification, the first officer retains a canonical Briefing,
+binds it to the current attempt, and commits that package before presenting the
+gate. A gate with no selected attempt is still validating, not ready for an
+unbound review. Without decision authority, the First Officer presents the bound
+review and stops with the attempt open.
```

In `docs/site/get-started/first-workflow.md`, replace the prototype example:

```diff
-Gate review: Add rate limiting to the API — review
-Chosen direction: token-bucket limiter at the API middleware layer
-Recommend approve.
-
-Checklist (from ## Stage Report in docs/ship-features/add-rate-limiting-to-the-api.md):
-- DONE: limiter implemented with per-client buckets
-- DONE: tests cover burst and refill behavior
-
-Assessment: 2 done, 0 skipped, 0 failed.
-
-Decision: approve to close; reject to bounce back to implementation.
+Capability/change: add a token-bucket limiter at the API middleware layer.
+Test and evidence: 2 done, 0 skipped, 0 failed; burst and refill tests pass.
+Reviewed snapshot: Briefing `...` at digest `sha256:...`.
+Findings: none.
+Recommendation: approve because the acceptance checks pass.
+Decision ask: approve to close, revise with concrete feedback, or hold at review.
```

Follow it with: “The First Officer binds and commits that reviewed Briefing before this message. You approve, revise with feedback, or hold; delegated conn changes who makes the decision, not whether the review appears.”

In `docs/specs/scenario-testing-principles.md`, replace `gate-guardrail`'s “without mutation” seed with: “the FO durably binds the retained Briefing, presents exactly one current six-field review, and halts without Resolution, consume, advance, dispatch, or archival when no decision authority is present.”

### Classifications and non-goals

The Pi failure is **runner provisioning friction** under `docs/runtime-support.md`: the explicitly loaded `pi-subagents` 0.28.0 package imports optional `@earendil-works/pi-coding-agent`, which is absent, so the scenario never launched. A compatible installed 0.24.0 package and `@mariozechner/pi-coding-agent` peer exist under the FNM root, and the runner already accepts `PI_SUBAGENTS_PACKAGE_ROOT`. Validation must select and verify that root in isolation. This is not product evidence, not permission to mutate global packages, and not grounds to remove Pi from AC-8; package-health UX, if desired, is a separate runtime-support task.

The round recorder's refusal is **correct fail-closed behavior**. A published advisory round represents resolved worker triage, so unresolved material cannot truthfully use `disposition: fixed` or a non-material decline. Retaining the complete pre-publication Briefing/log/refusal and returning to design is the correct path. Draft/unresolved-round support may be proposed separately, but no recorder-contract change belongs in 6y.

Provider overrides, canonical package verification, Captain-language mapping, recorder/schema/application code, CLI paths, one-use consume, resume, terminal merge, host adapters, prompts, retained Roborev evidence, and PR #565 remain outside this reset. No code, gate, status, PR, merge, global package, or review-round mutation occurs during ideation.

## Stage Report: ideation (cycle 22)

- DONE: Resolve the v1 contract contradictions exposed by Roborev 2170: retire prototype-only marker/legacy obligations rather than adding compatibility, define truthful bind-before-present behavior, and re-scope the onboarding/documentation surface explicitly.
  The canonical reset deletes marker semantics and unbound presentation, refits the existing guardrail to committed-open bind-and-hold, and records concrete concepts/onboarding/spec diffs within 17 named existing files.
- DONE: Design falsifiable fixes for successful state-commit attribution and the two distinct Claude misses, and classify the Pi dependency and unresolved-material round-recorder refusal without weakening the supported-host or durable-authority ACs.
  Existing structured host events pair with wrapper `exit=0`/`state-head`; narration and failed commits red. Exact bound tuple and opaque conn capture red independently on omission/mutation. Pi is required runner provisioning; unresolved round refusal is correct fail-closed behavior.
- DONE: Produce one coherent canonical design delta, expected file/LOC boundary, reset conditions, and test plan for independent staff review before implementation resumes.
  The complete boundary is 17 existing files, intended +297/-160 from `08675f02`, full +471/-224 from `13d70249`, hard stop +510 additions, with per-file estimates, component caps, deterministic/live sequencing, and explicit reset triggers above.
- SKIPPED: Modify code, PR #565, gate records, workflow status, global Pi packages, review-round state, or merge state.
  This design-reset stage changes only the entity body; the four implementation commits remain local and unpushed.

### Summary

Cycle 22 establishes one unreleased-v1 gate path: commit the exact binding, present one truthful review, then record only exact authority. It preserves Codex success and every recorder/consumer/provider boundary, turns both Claude misses into lossless data-flow obligations, keeps Pi required after runner provisioning, and treats unresolved advisory publication as correctly fail-closed. Checklist accounting is **3 DONE, 1 SKIPPED, 0 FAILED**; the design is ready for independent staff review before implementation resumes.

## Implementation intended-change declaration (cycle 23)

Declared before code edits at local unpushed tip `08675f02ff278bb16965ce583f00b15bd724e81f`. This cycle implements the approved unreleased-v1 reset in exactly these 17 existing files:

| File | Intended incremental delta |
| --- | ---: |
| `skills/first-officer/references/first-officer-shared-core.md` | +3 / -3 |
| `skills/fo-gate-lifecycle/SKILL.md` | +12 / -6 |
| `skills/present-gate/SKILL.md` | +12 / -5 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +90 / -25 |
| `internal/ensigncycle/shared_fixtures_test.go` | +20 / -25 |
| `internal/ensigncycle/gate_assert_impl_test.go` | +22 / -15 |
| `internal/ensigncycle/gate_assert_test.go` | +12 / -8 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +10 / -6 |
| `internal/ensigncycle/shared_scenarios_test.go` | +2 / -2 |
| `internal/ensigncycle/claude_live_runner_test.go` | +10 / -5 |
| `internal/ensigncycle/codex_live_runner_test.go` | +10 / -5 |
| `internal/ensigncycle/live_gate_stop_test.go` | +35 / -20 |
| `internal/ensigncycle/livescenario_adapter_live_test.go` | +15 / -12 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +22 / -4 |
| `docs/site/concepts/gates-and-decisions.md` | +8 / -4 |
| `docs/site/get-started/first-workflow.md` | +10 / -13 |
| `docs/specs/scenario-testing-principles.md` | +4 / -2 |

Intended incremental total is **+297/-160**. The full branch intent versus `13d70249` is **+471/-224**; additions have a hard stop at **+510** versus that baseline. Shared core remains at or below **26,754 bytes**, lifecycle at or below **7,000 bytes**, and presenter at or below **7,500 bytes**. Any eighteenth file, cap breach, new recorder field/harness/protocol/host lifecycle, unbound presentation, prompt coaching, weakened tuple/directive grade, or repeated unchanged supported-host obligation returns to ideation before further edits.

## Stage Report: implementation (cycle 23)

- DONE: Implement and commit the approved v1 reset as a bounded counterexample checkpoint.
  Commit `ce4365053158ef80c1a4dc34c63256fd83da62d0` on `spacedock-ensign/first-officer-gate-command-lifecycle-rebased-cc51e518` changes exactly the declared 17 files by **+497/-460** from `13d70249`, within the +510 hard stop. It removes the unbound legacy path, makes no-attempt gates validating, requires bind/commit before presentation, captures the exact tuple and opaque directive, refits the no-authority path to committed-open hold, and structurally keeps reviewer `REJECTED` outside the Captain lifecycle. Component sizes are shared core **26,521**, lifecycle **6,505**, and presenter **6,906** bytes.

  | File | Additions / deletions from `13d70249` |
  | --- | ---: |
  | `docs/site/concepts/gates-and-decisions.md` | +24 / -22 |
  | `docs/site/get-started/first-workflow.md` | +6 / -11 |
  | `docs/specs/scenario-testing-principles.md` | +3 / -2 |
  | `internal/contractlint/fo_function_reference_invariant_test.go` | +70 / -1 |
  | `internal/ensigncycle/claude_live_runner_test.go` | +28 / -6 |
  | `internal/ensigncycle/codex_live_runner_test.go` | +12 / -7 |
  | `internal/ensigncycle/gate_assert_impl_test.go` | +29 / -17 |
  | `internal/ensigncycle/gate_assert_test.go` | +20 / -40 |
  | `internal/ensigncycle/live_gate_stop_test.go` | +14 / -191 |
  | `internal/ensigncycle/livescenario_adapter_live_test.go` | +10 / -30 |
  | `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +230 / -36 |
  | `internal/ensigncycle/shared_fixtures_test.go` | +10 / -49 |
  | `internal/ensigncycle/shared_scenarios_negative_test.go` | +17 / -18 |
  | `internal/ensigncycle/shared_scenarios_test.go` | +1 / -1 |
  | `skills/first-officer/references/first-officer-shared-core.md` | +2 / -2 |
  | `skills/fo-gate-lifecycle/SKILL.md` | +5 / -3 |
  | `skills/present-gate/SKILL.md` | +16 / -24 |

- DONE: Complete every offline, repository, race, documentation, and compile gate at the checkpoint.
  `gofmt -w ./cmd ./internal`, focused lifecycle/contract controls, `go test ./...`, `go test ./... -race`, `go test -tags live ./internal/ensigncycle -run '^$' -count=1`, `git diff --check`, and strict MkDocs through `uv run --with-requirements docs/requirements.txt mkdocs build --strict` all pass. The real-CLI no-authority control proves one bound open attempt with the exact Briefing tuple and no Resolution, consume, advance, or dispatch.

- FAILED: Retain complete pass/fail Pi evidence without mutating global packages.
  Claude passes `gate-guardrail` in 108.49s and `recorded-gate-lifecycle` in 261.48s; Codex passes them in 99.83s and 219.01s. The package preserves all three Pi root sessions, surviving process output, and round 3 successor subagent metadata/output/session. It does **not** preserve the required Pi wrapper command logs, entity before/after snapshots, or state Git histories: those remained in temporary fixture repositories that were cleaned when the tests returned, and neither the state package nor the surviving `/tmp/spacedock-cycle23-gate-lifecycle-evidence` tree contains them. The retained sessions cannot recreate those facts. The package is therefore an incomplete Pi counterexample record, not complete host pass/fail evidence.

- FAILED: Produce a green supported-Pi lifecycle without crossing the declared reset boundary.
  Pi round 3 again reports `successor dispatch was not observed after consume` after 271.08s, while its retained subagent artifacts show a fresh `delegate` exited 0 and reported a path-scoped commit of the requested marker. Without the missing command log, before/after entity pair, and state Git history, that durable effect cannot be independently reconstructed after fixture cleanup. The same root session also transformed the opaque conn into a `BEGIN_CONN`/`END_CONN` block before recorder input. Reconciling this split requires a new or expanded harness/host-lifecycle attribution boundary, weakening exact directive grading, or another unchanged supported-host retry. Each is an explicit cycle-23 hard-reset trigger, so implementation stops with the surviving counterexample artifacts intact.

- SKIPPED: Final-tip Roborev, advisory round publication, PR #565 update/push, gate or round recording, workflow-status mutation, merge, and any further Pi retry or fix.
  The named hard reset fired before those actions. The checkpoint remains local and off PR #565; no gate, status, review-round, merge, or global package state changed.

### Summary

Checklist accounting is **2 DONE, 1 SKIPPED, 2 FAILED**. Cycle 23 proves the reset across deterministic controls, the repository, Claude, and Codex, and preserves the surviving artifacts from all three Pi rounds, but fails the complete Pi retention checklist. The repeated Pi successor boundary, exact-token mutation, and missing durable-state evidence require a new ideation decision about existing harness/host attribution; the implementation worker stopped at that first named boundary.

## Canonical design reset (cycle 24)

### Problem and evidence boundary

The cycle-23 checkpoint `ce4365053158ef80c1a4dc34c63256fd83da62d0` is the only code starting point. Its retained Pi round-3 root session contains one real native `subagent` call with `context: "fresh"`, the fixture cwd, explicit model `openrouter/openai/gpt-4.1-mini`, and the emitted handoff assignment path. The result with the same tool-call ID reports `isError: false`, exit 0, that model, one child session/artifact, and one completed run; the child records the marker and reports a path-scoped commit. Those are valid structured inputs to a future oracle.

They are not a retrospective pass. The cleaned fixture took away `command.log`, the before/after entity pair, and the state repository history, so the surviving marker and sessions cannot prove that the consumed commit preceded that spawn or that the child commit is the retained state descendant. The root session also contains two `gate --help` invocations. Cycle 23's observation folded `count(help)==1` into `ordered` and then reported the combined failure as “successor dispatch was not observed”; that is a design hypothesis for the misleading classification, not a regrade of the missing wrapper log.

The retained Codex lifecycle stream likewise cannot be promoted to stronger proof: it shows successful `dispatch build`, root narration that a worker was dispatched, an empty collaboration wait, and a later marker commit, but no structured native spawn receipt. Narration is not dispatch. Claude's retained structured Task call/result remains useful, but all three hosts must satisfy the same new contract on their next evidence-producing run.

### One host-neutral successor observation

Replace the current boolean/count bundle with one normalized `successorObservation`. Host adapters only extract native event shapes; they do not own lifecycle rules. The common observation contains:

- `consume`: the successful consume command, its consumed state commit OID, and its event position;
- `build`: exactly one successful `dispatch build`, returned assignment path, assignment digest captured by the existing wrapper at build return, wrapper `dispatch-head`, and event position;
- `start`: exactly one native worker-start receipt, native call/run identity, assignment path consumed, fresh-context value, cwd, stamped model, and event position;
- `completion`: exactly one native result joined to the start's call/run identity, successful terminal state/exit, reported model, and event position;
- `effect`: final entity bytes, exact marker/report count, the introducing commit OID, and retained state-history evidence.

The common evaluator accepts only the chain `consume commit ancestor-or-equal to build dispatch-head` → successful build result → native start consuming that exact returned path → matched native completion → one marker-introducing commit descended from dispatch-head. It also requires final nonterminal `handoff` state, one build/start/completion/effect, the expected cwd and freshness, and the expected Pi model when the scenario declares one. A build without a native start is not dispatch; a native result without the matching call identity is not completion; a marker without its introducing descendant commit is not a worker effect.

The assignment join is path plus build-time SHA-256, not path alone. The existing wrapper records the returned assignment path and digest immediately after the successful build, before another process can replace the shared `/tmp` pathname. Each native start must reference that exact path in its structured argument; the state ancestry and same-session event order complete the join. This extends the existing logging shim and host stream parsers, not a second harness or a host-specific lifecycle.

Pi extraction pairs the retained `subagent` tool call/result by `toolCallId`, checks `agent: delegate`, `context: fresh`, cwd, task assignment path, result count/status/exit, run ID, and attempted/reported model. Claude extraction pairs Task tool-use/tool-result IDs and assignment pointer. Codex extraction pairs collaboration spawn/terminal events and assignment pointer. Each emits the same structure. If a host stream does not expose a native start—as in the retained Codex cycle-23 stream—the stronger lifecycle scenario fails honestly while its independent gate guardrail remains unaffected.

Capability preflight is a separate observation and failure class. It still requires exactly one successful `gate --help` immediately before lifecycle mutation. An extra or missing help call fails `capability preflight cardinality/order`; it can no longer erase or masquerade as native successor evidence. This separates diagnostics without weakening either invariant.

### Lossless delegated authority

Replace sentinel framing with one natural-language data sentence whose payload is produced by `json.Marshal`, for example `Delegated authority data (JSON object; the decoded conn member is the exact grant): {"conn":"..."}`. JSON escaping is the transport; after decoding, only the member value is authority. The goal-only prompt still contains no shell command, review label, decision, actor, reason, lifecycle procedure, or recorder instruction.

The durable grade parses the single `adoption-note` scalar and compares its decoded bytes to `recordedGateDirective`; substring presence is insufficient. It separately requires one directive-bearing scalar and zero `BEGIN_CONN`/`END_CONN` sentinel bytes anywhere in the prompt or durable value. Table controls prove exact keep passes and drop, byte mutation, substitution, duplication, prefix/suffix wrapping, sentinel wrapping, and JSON-object leakage all fail. Real CLI argument capture continues to require the exact directive bytes.

### Fail-safe Pi artifact retention

`recorded_gate_lifecycle_pi_live_test.go` becomes the justified eighteenth cumulative file because it alone owns all three values needed to beat cleanup: the permanent artifact directory, the captured pre-run entity bytes, and the temporary fixture/state-repository lifetime. After every temp directory has registered its cleanup and immediately before the fatal-capable Pi command, it registers a later `t.Cleanup`; Go's last-added-first-run cleanup order makes retention execute before `t.TempDir` removal even after `t.Fatalf`.

That cleanup writes a manifest and independently attempts every artifact: wrapper `command.log`, exact before entity, resolved final entity, state HEAD/status, and `git bundle create state.bundle --all` for inspectable refs/ancestry. Existing permanent root/child Pi sessions stay where they are. Missing or unreadable required input is recorded in the manifest, does not prevent remaining copies, and fails the test via cleanup-safe error reporting. Thus even a process failure or oracle fatal leaves either the required evidence or an explicit per-artifact retention failure.

Deterministic controls run before any live host:

- a Pi-shaped valid call/result, build digest, consume ancestry, and child effect passes the common evaluator;
- drop/duplicate/mismatched tool ID, wrong model, non-fresh context, wrong cwd/path, completion-before-start, effect-before-consume, non-descendant commit, root-authored marker, and build-only evidence each fail their named invariant;
- a valid successor with two help calls still passes successor evaluation but fails the separate preflight evaluation;
- retention under simulated fatal preserves all five required fixture artifacts and an inspectable bundle; simulated missing log preserves the other four plus manifest and returns failure;
- delegated-conn controls independently falsify keep/drop/mutate/duplicate/wrapper leakage as specified above.

Only after those controls, the focused offline suites, Claude lifecycle, Codex lifecycle, repository tests/race, live compile, and docs checks are green may implementation run Pi once with the unchanged supported-host scenario and model. Pass requires both observations and the retained package. Any same preflight or successor failure, missing retention artifact, auth/provisioning failure, directive inequality, or host-native proof gap stops immediately with no prompt adjustment, retry, fallback model, global package change, or second Pi run.

### Deletion-first surface and implementation tolerance

The 497 added lines at cycle 23 are not a new floor. Before adding the Pi owner file or successor fields, implementation must consolidate/delete at least 107 checkpoint-added lines and show measured cumulative additions at or below **390** versus `13d70249`. Table-drive repeated host-review and directive fixtures in `recorded_gate_lifecycle_test.go`; remove synthetic gate-assert and shared-negative cases duplicated by the real-CLI scenario; fold the two contract-lint lifecycle variants into their existing ownership table. No behavioral invariant above may be deleted to hit the number.

The complete revised implementation surface is seven existing files plus the one newly admitted Pi live owner:

| File | Intended delta from `ce436505` |
| --- | ---: |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +55 / -120 |
| `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` | +24 / -0 |
| `internal/ensigncycle/claude_live_runner_test.go` | +8 / -12 |
| `internal/ensigncycle/codex_live_runner_test.go` | +8 / -12 |
| `internal/ensigncycle/gate_assert_test.go` | +0 / -16 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +0 / -12 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +0 / -20 |

The intended cycle-24 delta is **+95/-192**. The final measured branch should be at or below **+465 additions** versus `13d70249`; **+480** is the hard stop, strictly below the cycle-23 checkpoint. Failure to reach the deletion checkpoint, a nineteenth cumulative file, more than 95 new lines, a second harness, new production Go, a recorder/provider/schema field, a host-owned lifecycle, prompt coaching, or weakened exact-once/ancestry/directive grade returns to ideation before live work.

No skill or user documentation changes are designed: cycle 23 already states the user-visible v1 lifecycle. This reset changes test serialization, observation, attribution, and failure retention only. Existing bind-and-hold, provider, recorder, consumer, resume, terminal merge, failure-retention, and Claude/Codex gate-guardrail behavior remain unchanged.

### Acceptance and independent-review plan

- **AC-1 / AC-8 successor value:** each supported host yields one normalized native start/completion and one descendant worker effect after a consumed commit. The negative matrix proves no build, narration, root write, duplicate, wrong model, or reordered event can substitute.
- **Exact authority value:** the only durable delegated grant equals the Captain bytes after JSON decoding; all loss/mutation/wrapper counterexamples fail independently, and the prompt adds no operating procedure.
- **Failure-retention value:** every Pi outcome leaves independently inspectable command, entity, and Git ancestry evidence before cleanup, or a manifest that names the missing source while retaining the rest.
- **Lean-contract value:** deletion precedes expansion, final additions fall below cycle 23, and the only eighteenth file is the fixture-lifetime owner.

Independent staff review occurs after the deterministic controls and deletion checkpoint, before any live host. It inspects the normalized schema, native host joins, Git bundle, directive falsifiers, and measured diff. Implementation then follows one unchanged Claude/Codex proof sequence and at most one unchanged Pi run; the first named failure returns the retained package and resets rather than being coached around.

### Classifications and non-goals

This design does not infer the missing cycle-23 wrapper/state facts, regrade Pi or Codex, change runtime support, add a driver, alter the gate command surface, or mutate global Pi packages. It does not modify code, run tests, record a gate/round, change status, push PR #565, merge, or alter the checkpoint during ideation.

## Stage Report: ideation (cycle 24)

- DONE: Define one host-neutral successor-dispatch observation that accepts Pi's native completion and durable effect without weakening order, exactly-once, model, commit, or other-host proof.
  The normalized chain joins consume commit, build-time assignment identity, native start, matched completion, and descendant effect.
  Pi, Claude, and Codex adapters extract structured native events only; narration and `dispatch build` alone are explicitly insufficient.
  Capability-preflight cardinality remains exact-one but is graded separately, preventing its failure from being mislabeled as missing dispatch.
  The retained incomplete Pi and Codex records are inputs and counterexamples, not retroactive passes.

- DONE: Define lossless delegated-conn serialization, exact durable grading, and cleanup-safe live artifact retention.
  `json.Marshal` carries the conn as a data member in one natural-language sentence with no commands, labels, or procedure.
  The decoded `adoption-note` must equal the exact directive bytes; keep/drop/mutate/substitute/duplicate/prefix/suffix/sentinel/JSON leakage are independent controls.
  Pi's owning live test registers retention after temp setup and before the fatal-capable run so it executes before fixture cleanup.
  It preserves command log, before/after entity, state HEAD/status, and a Git bundle, while a manifest records missing inputs without suppressing other copies.

- DONE: Put deterministic Pi-shaped counterexamples and a one-run supported-host stop rule ahead of another live attempt.
  The matrix falsifies identity, model, freshness, cwd/path, ordering, ancestry, authorship, duplicate, build-only, preflight, directive, and retention failures.
  Offline, Claude, and Codex proofs must be green first.
  Pi then runs once unchanged; the first repeated or new named failure stops with the retained package and no coaching, retry, or global mutation.

- DONE: Produce a deletion-first, independently reviewable implementation boundary that challenges the 497-line checkpoint.
  At least 107 checkpoint-added lines must disappear before expansion, reaching at most +390 cumulative additions.
  Seven existing files plus the Pi fixture-lifetime owner target +95/-192 from `ce436505`, at most +465 additions versus baseline, with +480 a hard stop.
  The eighteenth file is admitted only because it owns artifact destination, pre-run bytes, and cleanup ordering; a nineteenth file resets design.

- SKIPPED: Implement, format, compile, test, run a live host, mutate a gate/round/status/PR, change global packages, or merge.
  Ideation changed only this entity body on the shared state branch.

### Summary

Cycle 24 replaces an ambiguous build-and-marker oracle with one native-event and Git-ancestry contract shared by all hosts, transports authority losslessly as data, and makes Pi failure evidence survive fixture cleanup. It cuts before expanding, keeps every product boundary fixed, and permits exactly one unchanged Pi proof only after deterministic and other-host evidence is independently reviewable. Checklist accounting is **4 DONE, 1 SKIPPED, 0 FAILED**.

## Canonical design repair (cycle 24 staff-review correction)

Independent staff review returned **CHANGES REQUESTED**. This section supersedes cycle 24's native-start/common-SHA proof, automatic artifact-destination assumption, and +390 deletion checkpoint. The retained cycle-24 Stage Report remains historical accounting, not the approved design.

### Restored Cycle-9 proof split

The common live lifecycle contract is exactly:

1. the approved gate is consumed and that state is committed;
2. one successor `dispatch build` succeeds after the consumed commit;
3. the active host's already-supported dispatch observation occurs after that build;
4. one durable successor effect is introduced by a state commit descended from the consumed/build snapshot.

The common evaluator owns those four order/ancestry facts, exactly-one build/effect, final nonterminal `handoff` state, exact marker/report counts, and the separate exact-one capability preflight. It does not prescribe a universal transport event. A build alone remains insufficient because it still needs the host-supported observation and descendant effect.

Codex uses the approved multi-agent-v2 evidence already owned by `codex_dispatch_evidence_test.go`: a successful entity-targeted `dispatch build`, a completed foreground `wait` ordered after that build, and a later durable stage-report observation. `codexDispatchCompletionEvidenceFromJSONL` and its regression controls remain the transport-shape authority. The recorded lifecycle runner consumes that result; it does not require, synthesize, or infer `spawn_agent`, native-start, run-ID, or native-completion events that the Codex stream may omit.

Claude may enrich its host observation by pairing the Task call/result IDs it exposes. Pi may enrich its host observation by pairing `subagent` call/result IDs and checking fresh context, cwd, successful completion, and the explicitly required child model. Those richer facts are host-owned evidence, not fields every host must manufacture. Pi's retained round-3 call/result is therefore a valid deterministic fixture input, but the missing wrapper/entity/history still prevents a retrospective lifecycle verdict.

Capability preflight remains independently exact-one and ordered before mutation. Two help calls can fail preflight while an otherwise valid successor chain remains accurately classified; neither result overwrites the other.

### No authoritative assignment-byte join

Remove assignment SHA-256 from the common observation and from acceptance. A returned assignment path, displayed prompt, or optional diagnostic digest may help explain a failure, but no path/digest combination proves the bytes a worker consumed. The common value is the observed host dispatch followed by descendant durable state, not a synthetic assignment-byte custody claim.

Pi/Claude structured calls may show that their native task argument references the emitted path because those hosts expose it. That is additional transport evidence only. Codex retains its successful-build/completed-wait ordering observation. No host needs a second harness, assignment interception protocol, or new recorder field.

### Delegated-conn serialization and exact-byte controls

The sentinel replacement remains one natural-language data sentence containing a `json.Marshal` result, with the decoded `conn` member defined as the exact grant. The prompt adds no command, review label, decision, actor, reason, recorder call, or lifecycle procedure. Durable grading decodes the single `adoption-note` scalar and compares its bytes for equality; substring, prefix, suffix, duplication, or sentinel wrapping cannot pass.

The deterministic table now includes JSON values containing a double quote, a backslash, an actual newline, and one value combining all three. For each, marshal → prompt extraction → JSON unmarshal must reproduce the original bytes exactly. The durable directive matrix still proves exact keep passes while drop, mutation, substitution, duplicate, prefix/suffix, `BEGIN_CONN`/`END_CONN`, and encoded-object leakage fail independently.

### Persistent Pi retention is a launch precondition

The admitted `recorded_gate_lifecycle_pi_live_test.go` owner must resolve `SPACEDOCK_LIVE_ARTIFACT_DIR` before `runPiLiveCommand`. Empty, relative, uncreatable, or cleanup-owned destinations fail the test before Pi launches. The owner canonicalizes the configured directory, creates a run-specific child, and rejects any destination equal to or contained by the fixture root, Pi home, session directory, clean home, or other test cleanup root it created. There is no `t.TempDir` fallback for this scenario.

Only after that preflight succeeds does the owner register cleanup retention and launch Pi. Last-added-first-run cleanup copies the wrapper log and exact before/after entity, writes state HEAD/status, creates `state.bundle` with all refs, and emits a per-artifact manifest into the persistent run directory. One copy failure records the failure, preserves every other available artifact, and fails the test safely.

Deterministic owner controls exercise three paths without launching Pi:

- missing `SPACEDOCK_LIVE_ARTIFACT_DIR` fails with a destination-precondition diagnostic;
- a destination inside a supplied cleanup root fails before any launch callback can run;
- an absolute, writable destination outside all supplied cleanup roots is accepted, receives the manifest/snapshots/bundle under the simulated-fatal cleanup path, and is explicitly removed by the control after inspection rather than registered as a launch cleanup root.

The launch callback is a counter in these controls, so both invalid cases prove zero launches. These controls and the Pi-shaped host observation run before the single unchanged supported Pi run.

### Corrected deletion-first surface and arithmetic

Cycle 23 starts at **+497 additions** versus `13d70249`. With up to **+95** implementation additions, final additions can remain at or below **+465** only if deletion/consolidation first reaches **+370** or fewer cumulative additions. Therefore implementation must remove or consolidate at least **127 checkpoint-added lines before any expansion**. The prior +390 checkpoint is void.

At the deletion checkpoint, record both:

- checkpoint-relative `git diff --numstat ce436505 -- <declared surface>` totals, showing what cycle 24 removed; and
- cumulative `git diff --numstat 13d70249 -- <full branch surface>` totals, whose additions must be **≤370**.

Record both tables again at final tip. Final checkpoint-relative additions must be **≤95**; cumulative additions target **≤455**, must be **≤465** under the declared tolerance, and retain the absolute **+480 hard stop**. If deletion changes baseline-owned lines, only measured cumulative numstat—not subtraction by intent—decides whether expansion may start.

The corrected implementation surface remains seven existing files plus the one admitted Pi owner:

| File | Expected delta from `ce436505` |
| --- | ---: |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +45 / -120 |
| `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` | +28 / -0 |
| `internal/ensigncycle/claude_live_runner_test.go` | +6 / -8 |
| `internal/ensigncycle/codex_live_runner_test.go` | +6 / -8 |
| `internal/ensigncycle/gate_assert_test.go` | +0 / -16 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +0 / -12 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +0 / -20 |

Expected cycle-24 delta is **+85/-184**, with **+10 additions tolerance** and no deletion credit assumed until numstat measures it. `codex_dispatch_evidence_test.go` is an existing, unchanged proof owner, not a ninth repair file. The Pi file remains justified because it alone owns the configured persistent destination, before bytes, fixture lifetime, and pre-launch/cleanup ordering.

The same reset triggers remain: failure to reach +370 before expansion, more than +95 checkpoint-relative additions, cumulative additions above +465 (or any +480 hard-stop breach), a nineteenth cumulative file, a second harness, production Go, recorder/provider/schema changes, host lifecycle duplication, prompt coaching, or weakened exact-once/ancestry/directive behavior.

### Repaired acceptance criteria and test plan

- **Supported-host lifecycle value:** Claude, Codex, and Pi each prove consumed commit → successful build → their supported dispatch observation → one descendant durable successor effect. Common deterministic negatives remove/reorder/duplicate each link. Codex regression fixtures specifically fail if wait precedes build, build fails, wait is incomplete, or durable report evidence precedes wait; no spawn event is required.
- **Exact-authority value:** JSON data transport round-trips current, quoted, backslashed, newline, and combined grants byte-for-byte, and only one exactly equal durable scalar passes. Every drop/mutation/wrapper case fails.
- **Failure-retention value:** Pi cannot launch without a verified persistent external destination. Missing and cleanup-owned controls observe zero launches; the accepted-destination simulated failure leaves command/entity/history artifacts and an inspectable bundle.
- **Lean-contract value:** measured cumulative additions are ≤370 before expansion and ≤465 at final tip, with checkpoint-relative and cumulative numstat published at both points. This is the independent baseline that can move the wrong way.

The simplest proof considered was the cycle-23 build-plus-marker count; it cannot distinguish host dispatch ordering or root-authored effects. A universal native-call schema was also considered and rejected because Codex intentionally exposes a different approved transport shape. The restored host observation plus common Git ancestry is the smallest mechanism that serves the supported-host lifecycle AC.

The simplest retention alternative was copying after the live call; `t.Fatalf` and temp cleanup can erase the sources first. A cleanup hook without a persistent destination merely moves the same loss. Pre-launch destination validation plus last-added cleanup is the smallest mechanism that makes failure evidence survive.

No spike runs during this repair. The mechanisms are already fixture-backed: Cycle-9 Codex dispatch evidence tests own build/wait/report ordering, Go JSON owns lossless string round-trip, Go cleanup ordering is established, and Git bundles preserve refs. Implementation's first work is the deterministic deletion checkpoint and the newly named missing/persistent-destination plus JSON edge controls. Then run focused offline controls, full/race/live compile/docs, unchanged Claude and Codex lifecycle proofs, independent staff review, and at most one unchanged Pi run. The first named failure stops with retained evidence.

No skill text, user documentation, product command, provider/recorder/consumer boundary, bind-and-hold behavior, gate guardrail, global package, PR, gate, round, status, or merge state changes in this repair.

## Stage Report: ideation (cycle 25)

- DONE: Restore the binding Cycle-9 proof split without requiring or synthesizing Codex native start/completion events.
  The common chain is consumed commit, successful build, host-supported dispatch observation, and descendant durable effect.
  Codex reuses its existing successful-build/completed-wait/later-report evidence owner; Claude and Pi may retain richer matched events only where exposed.
  Assignment paths and optional digests are diagnostic and cannot claim worker-consumed bytes.
  Existing Codex temporal negatives remain the falsifier for failed, incomplete, or misordered waits.

- DONE: Make persistent Pi retention a pre-launch requirement with deterministic invalid/valid destination controls.
  Missing, relative, cleanup-owned, or uncreatable destinations fail before the launch callback.
  A persistent external destination accepts cleanup retention and preserves command, before/after entity, state metadata, manifest, and Git bundle on simulated failure.
  The admitted Pi live owner is still the only file that holds destination, fixture lifetime, and cleanup ordering together.
  Invalid-destination controls assert a zero launch counter rather than relying on error text alone.

- DONE: Complete lossless delegated-conn coverage and preserve exact durable authority grading.
  JSON quote, backslash, actual-newline, and combined values round-trip to identical decoded bytes.
  Exact keep is the sole pass; drop, mutate, substitute, duplicate, prefix/suffix, sentinel, and encoded-wrapper leakage remain independent failures.

- DONE: Correct deletion arithmetic and publish an auditable checkpoint/final measurement plan.
  Pre-expansion cumulative additions are now ≤370, requiring at least 127 measured checkpoint-added lines removed before up to +95 additions.
  Both checkpoint-relative and cumulative numstat are recorded after deletion and at final tip.
  Expected cycle delta is +85/-184; final target is ≤455, tolerance is ≤465, and +480 remains the hard stop.
  Measured cumulative numstat, not intended deletions, decides whether implementation may expand.

- SKIPPED: Implement, format, compile, test, launch a host, or mutate code, PR, gate, round, status, package, or merge state.
  The staff-review repair changes only this entity body on the shared state branch.

### Summary

The repair restores Codex's approved host evidence instead of imposing an unavailable native-event schema, removes the false assignment-byte proof, and makes Pi evidence persistence a condition of launch. It also adds adversarial JSON round trips and corrects the deletion checkpoint so the promised final surface is arithmetically possible. Checklist accounting is **4 DONE, 1 SKIPPED, 0 FAILED**.

## Implementation intended-change declaration (cycle 26)

Declared before code edits at local unpushed checkpoint `ce4365053158ef80c1a4dc34c63256fd83da62d0`. This cycle changes exactly these eight files:

| File | Expected delta from `ce436505` |
| --- | ---: |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +45 / -120 |
| `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` | +28 / -0 |
| `internal/ensigncycle/claude_live_runner_test.go` | +6 / -8 |
| `internal/ensigncycle/codex_live_runner_test.go` | +6 / -8 |
| `internal/ensigncycle/gate_assert_test.go` | +0 / -16 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +0 / -12 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +0 / -20 |

Expected checkpoint-relative delta is **+85/-184**, with at most **+95 additions**. Deletion/consolidation happens before expansion and must first produce measured cumulative additions **≤370** versus `13d70249`, removing at least 127 checkpoint-added lines without inferred credit. Checkpoint-relative and cumulative numstat are recorded after deletion and at final tip. Final cumulative target is **≤455**, tolerance is **≤465**, and **+480** is the absolute hard stop.

Implementation stops before every live host after deterministic controls are green so independent staff can review the actual code and measured diff. The reset triggers and unchanged PR/workflow boundaries are those in the approved cycle-24 staff-review repair.

## Cycle-26 deletion checkpoint

Deletion/consolidation completed before feature expansion. Existing deterministic lifecycle, no-authority, gate-guardrail, and FO contract tests remain green:

- `go test ./internal/ensigncycle -run '^(TestRecordedGate|TestAssertGateHeld|TestGateGuardrailNegative)' -count=1`
- `go test ./internal/contractlint -run '^TestFO' -count=1`

Checkpoint-relative numstat from `ce436505`:

| File | Additions / deletions |
| --- | ---: |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +0 / -64 |
| `internal/ensigncycle/gate_assert_test.go` | +5 / -33 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +12 / -55 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +0 / -15 |
| **Total** | **+17 / -167** |

Cumulative numstat from `13d70249`:

| File | Additions / deletions |
| --- | ---: |
| `docs/site/concepts/gates-and-decisions.md` | +24 / -22 |
| `docs/site/get-started/first-workflow.md` | +6 / -11 |
| `docs/specs/scenario-testing-principles.md` | +3 / -2 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +6 / -1 |
| `internal/ensigncycle/claude_live_runner_test.go` | +28 / -6 |
| `internal/ensigncycle/codex_live_runner_test.go` | +12 / -7 |
| `internal/ensigncycle/gate_assert_impl_test.go` | +29 / -17 |
| `internal/ensigncycle/gate_assert_test.go` | +6 / -54 |
| `internal/ensigncycle/live_gate_stop_test.go` | +14 / -191 |
| `internal/ensigncycle/livescenario_adapter_live_test.go` | +10 / -30 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +188 / -37 |
| `internal/ensigncycle/shared_fixtures_test.go` | +10 / -49 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +7 / -23 |
| `internal/ensigncycle/shared_scenarios_test.go` | +1 / -1 |
| `skills/first-officer/references/first-officer-shared-core.md` | +2 / -2 |
| `skills/fo-gate-lifecycle/SKILL.md` | +5 / -3 |
| `skills/present-gate/SKILL.md` | +16 / -24 |
| **Total** | **+367 / -480** |

The measured additions fell by **130** from cycle 23's +497, satisfying both the at-least-127 removal requirement and the pre-expansion **≤370** gate. No deletion credit is inferred. Expansion may now begin within the declared eight-file surface and +95 checkpoint-relative additions cap.

## Cycle-26 deterministic implementation checkpoint

Local code commit `45c1a9a3ea0e8e29a61a39036de05de1125eb0ee` is ready for the required independent staff review. No Claude, Codex, or Pi live host journey ran.

The checkpoint adds a host-neutral dispatch-observation requirement, reuses the existing Codex build/wait/report evidence owner unchanged, observes Claude's structured Agent call, and observes Pi's structured subagent call/result shape. Delegated authority now travels as marshaled JSON and is graded as one exact durable adoption-note line; quote, backslash, actual-newline, combined, sentinel-wrapper, and JSON-object controls are deterministic. The Pi owner rejects empty, relative, uncreatable, and symlink-resolved cleanup destinations before launch, creates a run-specific external artifact directory, and retains command log, before/after entity, HEAD/status metadata, per-artifact manifest, and an inspectable all-refs Git bundle from cleanup.

Checkpoint-relative numstat from `ce436505`:

| File | Additions / deletions |
| --- | ---: |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +0 / -64 |
| `internal/ensigncycle/claude_live_runner_test.go` | +1 / -1 |
| `internal/ensigncycle/codex_live_runner_test.go` | +1 / -1 |
| `internal/ensigncycle/gate_assert_test.go` | +0 / -42 |
| `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` | +56 / -2 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +37 / -97 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +0 / -15 |
| **Total** | **+95 / -222** |

Cumulative numstat from `13d70249` is **+416 / -506**, within the ≤465 tolerance and below the +480 hard stop. The repair changes only the seven named paths above; `internal/ensigncycle/codex_dispatch_evidence_test.go` is byte-unchanged and remains the Codex transport authority.

Green deterministic/repository evidence:

- `go test ./internal/ensigncycle -run '^(TestRecordedGate|TestAssertGateHeld|TestGateGuardrailNegative)' -count=1`
- `go test ./internal/contractlint -run '^TestFO' -count=1`
- `go test -tags live ./internal/ensigncycle -run '^TestRecordedGatePersistentArtifactControls$' -count=1`
- `gofmt -w ./cmd ./internal`
- `go test ./...`
- `go test ./... -race`
- `go test -tags live ./internal/ensigncycle -run '^$' -count=1`
- `uv run --with-requirements docs/requirements.txt mkdocs build --strict`
- `git diff --check`

**STOP FOR INDEPENDENT STAFF REVIEW.** Review the actual commit/diff, especially the host-observation joins, exact directive falsifiers, symlink containment, cleanup ordering, manifest failure behavior, and bundle inspectability. Do not run a live host until that review explicitly approves this checkpoint.

## Stage Report: implementation (cycle 26)

- DONE: Replace Claude/Pi substring booleans with structured host-specific dispatch observations.
  Counterexample commit `3c53510529c095878f5156d490d498e6742fb31c` joins successful build results to later root Claude Agent or Pi delegate calls and rejects pre-build, wrong-task/description/context, list-only, unrelated-result, and missing-successor shapes; deleting any join or field constraint makes `TestRecordedGateStructuredDispatchObservations` fail.
- DONE: Parse frontmatter and the single Resolution adoption-note YAML scalar as decoded bytes.
  The same commit parses complete YAML frontmatter, separates Resolution cardinality from exact decoded directive equality, reports JSON decode errors, and passes quoted, backslash, newline, combined, sentinel, and object controls; reverting to serialized-line comparison makes `TestRecordedGateDurableDirectiveYAMLScalars` fail.
- DONE: Prove deterministic Pi artifact retention after fatal cleanup and a missing command log.
  `TestRecordedGatePersistentArtifactControls` now forks a child that fatals after registering cleanup; removing cleanup retention, manifest independence, or bundle creation loses a required artifact or makes bundle verification fail.
- FAILED: Restore/consolidate every deletion-sensitive v1 lifecycle invariant.
  The reviewer-required unbound/byte-identical hold refusal, Resolution and noncanonical-review negatives, reviewer-REJECTED ordering, Captain mappings, presenter sole ownership, exact tuple/directive, and root-review-before-mutation barrier remain absent because the hard addition stop fired before restoration.
- FAILED: Keep the correction within the checkpoint-relative +95 additions and cumulative +465 ceiling.
  Exact `45c1a9a3..3c535105` numstat is +191/-22 across three authorized files: Claude +1/-1, Pi +41/-10, and recorded lifecycle +149/-11; cumulative `13d70249..3c535105` is +592/-513, beyond both +465 tolerance and +480 hard stop.
- SKIPPED: Run full, race, live-compile, and documentation gates.
  The approved reset rule requires an immediate design stop once the measured cap cannot hold; only the focused deterministic controls and `git diff --check` ran green.
- SKIPPED: Request same-staff re-review or run any live host.
  The prototype is deliberately local and unpushed; no Claude, Codex, or Pi live journey, PR push, gate, round, status, or Resolution mutation occurred.

### Summary

Staff review found four material proof defects: substring host evidence, serialized rather than decoded authority comparison, unproven fatal cleanup retention, and deletion of v1 guardrails. The local counterexample implements and verifies the first three, but already costs +191 additions before the fourth, proving the approved +95 correction surface cannot truthfully satisfy the unchanged acceptance criteria. Checklist accounting is **3 DONE, 2 SKIPPED, 2 FAILED**; implementation stops for a new design rather than weakening invariants or exceeding the ceiling.

## Canonical obligation-authority reset (cycle 27)

This section replaces the canonical task body and every cycle-22 through cycle-26 correction. Historical reports and local commits remain diagnostic evidence, not implementation authority. The implementation baseline is the accepted remote tip `13d702492131df17dd3ac87245d6d773f4df959b`.

### Reset decision

Abandon all seven local commits after `13d70249` (`c4afd8f7`, `37d6980b`, `7419e61a`, `08675f02`, `ce436505`, `45c1a9a3`, and `3c535105`) and rebuild from the accepted tip. Preserve their SHAs as counterexamples; do not cherry-pick or incrementally salvage them.

This is safer than salvage because those commits mix three different things in one dependency stack: two useful skill-contract corrections, prototype-only unbound-review compatibility, and progressively stronger transcript/retention laboratories. The stack reached +592/-513 while still omitting required contract controls. Starting from `13d70249` keeps the already-proven real CLI, runner, recorder, consumer, and host-native fixtures, then lets implementation add only the two product corrections and delete the obsolete compatibility/proof layers.

There is no released-v1 compatibility authority: v1 is unreleased. Landed recorder and consumer behavior is a reusable implementation input, not a reason to preserve prototype markers, an unbound review, or a common live transport schema.

### Canonical authority audit

| Current obligation | Source classification | Cycle-27 disposition |
| --- | --- | --- |
| Bind one exact Briefing and commit it before review | Explicit Captain/product value; landed recorder contract | **KEEP.** Direct real-CLI proof reads the selected `(Briefing.id, digest)` from durable state. |
| Show one legible root-visible review, including under delegated conn | Explicit Captain/product value | **KEEP.** `present-gate` owns the existing review template; `fo-gate-lifecycle` waits for one root emission before decision mutation. Live proof checks legibility and multiplicity, not six exact labels. |
| Record direct Captain and delegated FO authority distinctly | Explicit Captain/product value; landed actor/directive contract | **KEEP.** Direct CLI cases prove `person:captain` versus `agent:first-officer`, nonblank reason, and exact decoded directive bytes. |
| Commit close, consume once, commit consumed state, then route | Explicit Captain/product value; landed consume contract | **KEEP.** Zero dispatch/advance before the consumed commit is the primary measured outcome. |
| Hold, revise, stale, blocked, and repeat-consume fail closed | Landed command semantics serving the product barrier | **KEEP, LEAN.** Existing lower-level tests remain authoritative; one direct lifecycle table covers only cross-command state/commit behavior. |
| Terminal approval has no successor dispatch | Landed status/consume/merge behavior | **KEEP as deterministic complement.** It does not need a route-by-host matrix. |
| Exact three-artifact 3k package, relative-path equivalence, folder commit, discovery equality | Validator/implementor expansion over already-landed CLI/state contracts | **CUT from 6y ACs.** Existing CLI, gates, and state tests own these mechanics; the representative fixture may reuse them without re-proving them. |
| Exactly one `gate --help`, cache absence, path-swap and stale-launcher experiments | Implementation convenience prompted by one stale launcher incident | **CUT from end-value ACs.** Keep the shipped fail-closed capability sentence; do not add cardinality or identity laboratories. |
| Reviewer `REJECTED` automatic routing before the Captain lifecycle | Existing feedback-flow behavior, not this Captain gate end value | **CUT from 6y.** `feedback-rejection-flow` and its own scenarios remain the owner; 6y neither changes nor re-proves it. |
| Captain aliases `redo`, routed/unrouted `reject`, and `not yet` across every host | Prototype/validator attempt to prove model-language interpretation | **CUT.** Canonical `approve|revise|hold` and their durable routes remain; no phrase-by-host matrix or model-causality claim. |
| Exact six-field labels, exact transcript interval, child/root row filters in the common grader | Prototype validator overlay | **CUT from common/live grading.** The shipped presenter template remains independently contract-tested; live journeys require one legible root review before durable decision. |
| A universal native start/completion schema or exact child author | Explicitly removed by the Captain’s Cycle-9 ruling | **CUT.** Existing Claude, Codex, and Pi native fixture owners retain their own transport assertions. |
| Assignment-path/byte custody, child model/cwd/context equality | Prototype forensics and implementor convenience | **CUT.** Successful durable successor outcome is the value; no test claims which bytes caused a model action. |
| Common build/wait/call/result joins as lifecycle proof | Prototype validator overlay | **CUT from 6y live journeys.** Host-native dispatch fixtures own transport shape; the lifecycle journey observes the command/state outcome naturally exposed by its runner. |
| JSON prompt encoding for arbitrary conn bytes and a new YAML frontmatter parser | Counterexample implementation mechanism | **CUT.** Direct CLI tests pass arbitrary directives and inspect the resulting `gates.Document` through the production parser. |
| Pi persistent-destination preflight, fatal-child retention, manifest, and Git bundle | Failure-forensics expansion beyond the existing runner | **CUT.** The existing Pi runner’s normal artifact directory is sufficient; a missing artifact is a runner failure, not a new 6y subsystem. |
| Legacy `Gate review:`/`Decision:` markers and byte-identical unbound gate hold | Prototype compatibility behavior | **DELETE.** A presentation-ready gate is bound; no alias, marker bridge, or unbound-review path remains. |
| Goal-only prompts and shipped-skill ownership | Evidence-integrity requirement serving the product value | **KEEP.** Prompts name fixture, authority, desired outcome, and stop condition, never the command procedure or review text. |

### Proposed minimum v1 behavior

1. On an engaged gate with a retained package, `fo-gate-lifecycle` records the Briefing and commits the folder state.
2. It rereads the selected Briefing identity/digest, performs the existing evidence judgment, and invokes `present-gate`. Chat presentation completes only after one legible root review is emitted. Delegated conn changes who may decide, not whether the Captain sees the review.
3. It records exactly one canonical `approve`, `revise`, or `hold` decision. Direct Captain and delegated First Officer actors remain distinct; delegated authority carries the exact directive value and a nonblank reason. The close is committed before any route.
4. Approval consumes once and commits the resulting state before nonterminal dispatch or terminal merge. Revise routes through the existing feedback owner after the close commit. Hold stops. Stale, blocked, wrong-stage, and spent approvals fail closed under existing command semantics.
5. A gate without a selected Briefing is validating, not presentation-ready. The old unbound/byte-identical review path and its marker-based live scenarios are removed.

### End-value acceptance criteria

**AC-1 (VALUE) — The normal First Officer gate path has zero unauthorized routes.** In the real-CLI fixture, the durable sequence is one bound Briefing commit, one decision close commit, and for approval one consumed successor commit. The count of decision records before a successful bind commit is **0**, and the count of advances or dispatches before a consumed descendant commit is **0**. Revise, hold, stale, blocked, and repeat-consume cases have **0** unauthorized consumes/dispatches. *Verified by:* real CLI plus real Git snapshots; omission/reordering controls make the state or ancestry assertion fail.

**AC-2 — Direct and delegated authority are lossless and fail closed.** Direct close records `person:captain`; delegated close records `agent:first-officer`, a nonblank reason, and one directive whose decoded value exactly equals the supplied bytes. Quote, backslash, newline, and combined directives round-trip through the real CLI and production `gates.Document` parser. Missing actor/reason/directive or a changed actor refuses or fails the durable grade. *Verified by:* table-driven public-CLI fixtures; no prompt serializer or test-only YAML parser.

**AC-3 (VALUE) — Each supported host demonstrates the user-visible lifecycle without a common transport fiction.** One unchanged goal-only approved journey on Claude, Codex, and Pi emits exactly one legible root review under delegated conn and finishes with a durable consumed nonterminal state plus the successor’s durable report/marker. The review precedes the durable decision in the runner evidence naturally available on that host. Native dispatch call, wait, result, child, model, cwd, and session-shape assertions remain exclusively in existing host-native fixtures. *Verified by:* the three existing runners with a reduced outcome grader; the terminal deterministic control supplies the zero-dispatch complement.

**AC-4 — The shipped contract is singular and the implementation is materially smaller than the abandoned stack.** `present-gate` is the sole review renderer; `fo-gate-lifecycle` owns bind/close/consume routing and waits for presentation. Goal-only prompts contain no procedure. No unbound-review compatibility, common native-event parser, assignment-custody oracle, or new retention subsystem exists. Against `13d70249`, additions remain at most 80 and the branch deletes at least 300 lines; the independent value measure is net source/test removal while AC-1 through AC-3 remain green.

### Lean test and implementation plan

1. **Reset and deletion checkpoint (no runtime cost; AC-4):** move the code branch to `13d70249`, preserve `3c535105` by SHA, remove the legacy gate-held assertion/live scenario and every registration/fixture/coverage entry that exists only for unbound marker compatibility, then record exact numstat before adding behavior.
2. **Real CLI lifecycle (low/medium; AC-1/AC-2):** retain the existing positive bind/close/consume and terminal cases. Consolidate omission, revise/hold, stale/blocked/spent, tuple, actor, and commit controls around actual binary calls and Git snapshots. Add arbitrary directive cases by reading the durable entity with the production gates parser. Simplest alternative—serialized substring checks—cannot distinguish YAML encoding from decoded authority.
3. **Skill ownership (low; AC-3/AC-4):** add one structural test that the presenter owns chat rendering/delegation visibility and the lifecycle waits before decision mutation; reject duplicate ownership and legacy unbound language. This serves root-visible review and has no model-causality claim.
4. **Supported-host journeys (high; AC-3):** after deterministic/full/race/docs/live-compile and independent staff approval, run the existing goal-only approved journey once on Claude, Codex, and Pi. Grade root-visible review and durable outcome only. Do not parse a common native dispatch event or add a failure-retention harness; inspect the existing runner artifacts on failure.
5. **Repository verification:** `gofmt -w ./cmd ./internal`, focused lifecycle/contract controls, `go test ./...`, `go test ./... -race`, live compile, strict MkDocs, and `git diff --check`. Host-native fixture suites remain unchanged and green.

No new spike is needed. `13d70249` already proves the real CLI mutations, Git commits, production parser, and three supported runners exist; retained Codex runs prove the missing root-visible review, while `3c535105` proves that common event/YAML/retention machinery is expensive and unnecessary. Implementation starts with deletion and deterministic tests, not live spend.

### Expected surface and evidence-based tolerance

The clean implementation touches **17 existing paths**: four skill/test paths gain implementation or evidence, two documentation paths receive small wording edits, and eleven are deletion-only or remove legacy registrations. The estimate derives from measured baseline blocks: `live_gate_stop_test.go` is 199 lines, `gate_assert_{impl,test}.go` total 104 lines, the legacy adapter is 97 lines, and the abandoned stack demonstrated that the skill barrier itself needs only a few lines.

| Path | Expected delta from `13d70249` |
| --- | ---: |
| `skills/fo-gate-lifecycle/SKILL.md` | +4 / -3 |
| `skills/present-gate/SKILL.md` | +4 / -1 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +14 / -0 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +35 / -160 |
| `internal/ensigncycle/gate_assert_impl_test.go` | +0 / -42 |
| `internal/ensigncycle/gate_assert_test.go` | +0 / -62 |
| `internal/ensigncycle/live_gate_stop_test.go` | +0 / -199 |
| `internal/ensigncycle/livescenario_adapter_live_test.go` | +0 / -45 |
| `internal/ensigncycle/shared_fixtures_test.go` | +0 / -30 |
| `internal/ensigncycle/shared_scenarios_test.go` | +0 / -5 |
| `internal/ensigncycle/shared_scenarios_meta_test.go` | +0 / -1 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | +0 / -25 |
| `internal/ensigncycle/pi_shared_coverage_test.go` | +0 / -4 |
| `internal/ensigncycle/claude_live_runner_test.go` | +0 / -25 |
| `internal/ensigncycle/codex_live_runner_test.go` | +0 / -20 |
| `docs/site/concepts/gates-and-decisions.md` | +6 / -4 |
| `docs/specs/scenario-testing-principles.md` | +2 / -2 |
| **Expected total** | **+65 / -628** |

Tolerance is **+15/-120** around that estimate: at most **+80 additions**, at least **300 measured deletions**, no new file, and no production Go. The wide deletion tolerance reflects shared helpers that may still serve unrelated scenarios; the narrow addition tolerance reflects the measured 4-line skill barrier, one compact contract test, and one production-parser directive table. Numstat—not intended credit—decides the checkpoint.

### Hard reset triggers

Return to ideation before live work if implementation requires production Go/schema/provider/recorder changes, a new file or harness, a common native-event parser, exact-child or assignment-custody evidence, prompt coaching, an unbound-review compatibility path, new Pi artifact-retention machinery, more than one approved journey per existing host, more than +80 additions, fewer than 300 deletions, or any path outside the declared surface. A host journey failure is reported from its existing artifacts; it does not authorize coaching, retry laboratories, or cross-host schema growth.

No code, test, live host, PR, gate, review-round, status, or merge mutation occurs during this ideation reset. Independent staff review must assess the authority cuts, AC sufficiency, and measured surface before implementation begins.

## Stage Report: ideation (cycle 27)

- DONE: Audit every current 6y obligation by authority and propose rigorous cuts.
  The canonical table distinguishes explicit product value and landed command inputs from prototype, validator, and implementor inventions; unreleased v1 carries no compatibility authority.
- DONE: Re-anchor the minimum product value and proof split.
  AC-1/AC-2 put exact bytes, tuple binding, one-use state, stale/hold/revise, and commit/consume guards in direct real-CLI controls; AC-3 leaves live journeys with one root-visible review and durable outcome while host-native fixtures own transport.
- DONE: Produce a clean reimplementation design from the accepted remote baseline.
  The design discards all seven local commits, preserves their SHAs as counterexamples, limits additions to four skill/test and two documentation paths, makes eleven paths deletion-only, estimates +65/-628 with evidence-based +80/-300 bounds, and names hard reset triggers.
- SKIPPED: Recover additional historical citations from AgentsView.
  `agentsview session search` could not initialize because its configured `~/.agentsview` path is inaccessible as a directory; the ticket’s durable Captain rulings and repository commit history supplied the authoritative record.
- SKIPPED: Edit/reset code, run tests or live hosts, push PR, or mutate gate/round/status.
  Ideation changed only this entity body in the shared state checkout.

### Summary

Cycle 27 removes proof inflation instead of compressing it. The minimum v1 contract is bind and commit, show one root review even under conn, record exact direct/delegated authority, consume and commit, then route; real CLI tests own exact state semantics, supported-host journeys own only visible review plus durable outcome, and native fixtures own transport. Checklist accounting is **3 DONE, 2 SKIPPED, 0 FAILED**; the design is ready for independent staff review.

## Canonical authority repair (cycle 27 staff-review correction)

Independent staff review returned **CHANGES REQUESTED**. This section supersedes Cycle 27’s cuts to dispatch observation, rejection/mapping ownership, semantic review ordering, no-authority coverage, folder durability, and its LOC-based AC/floor. The clean-baseline decision and the exclusions on common native transport, exact-child causality, assignment custody, prompt coaching, and new failure-forensics machinery remain binding.

### Restored product obligations

1. **Durable successor route.** A successful nonterminal approval has exactly one successful `dispatch build`, ordered after the descendant commit containing consumed application plus successor status, and exactly one later durable successor effect. The shared grader may use the logging wrapper and Git ancestry for these command/state facts. It does not require a common native call/result, child identity, model, cwd, session, or completion event; existing host-native fixtures own those facts.
2. **Correction and Captain-language ownership.** No other integration owner currently connects a completed reviewer `REJECTED` result to correction before Captain presentation, so shared core retains that automatic pre-branch. The deferred lifecycle retains deterministic mappings: `approve → approve`; redo-with-feedback → `revise` with an accepts-direction reason; reject with `feedback-to` → `revise` with a rejects-direction reason; reject without `feedback-to`, explicit hold, and `not yet` → `hold`, with `not yet` requiring a nonblank reason. Structural/mapping deletion controls own this contract; there is no route × host live matrix.
3. **Semantic root review.** Default chat emits exactly one qualifying root review after the selected-Briefing folder commit and before the decision invocation. A qualifying review semantically names entity/title and stage, the exact bound Briefing id and digest, one recommendation, and a decision ask. It need not use six exact labels. Claude, Codex, and Pi use their existing host-specific root-message and order filters; no common transcript/native-event schema is introduced.
4. **Bound-open no-authority behavior.** `gate-guardrail`, `live_gate_stop`, `livescenario_adapter`, their registrations, and negative controls are refitted rather than deleted. A no-conn run binds the selected Briefing, commits the folder, presents one qualifying review, and stops open at the gated stage with no Resolution, consume, advance, dispatch, or archival. The standalone headless/no-conn path still proves drive-to-gate-and-stop; the primitive still proves self-approval refusal independently.
5. **Package and capability boundary.** The positive path proves the selected Briefing folder commit includes the entity index and retained room files without sibling sweep. Capability prose stays fail-closed: an incomplete `gate --help` surface halts before room/entity mutation and names refresh/build remediation. Exact probe cardinality, executable identity/cache, and swap laboratories remain non-obligations.
6. **Directive domain.** Exact delegated directive coverage is every argv-capable UTF-8 string except NUL. Public CLI cases include quotes, backslashes, actual newlines, and combined values; the durable value is read through production `gates.Document`, never a serialized-line or test-only YAML parser. Real Git snapshots prove bind, close, and consume barriers around the same record.

Reviewer `REJECTED`, Captain mappings, exact-one semantic review, and bound-open no-authority behavior are therefore product/integration obligations, not compatibility aliases. Prototype `Gate review:`/`Decision:` marker matching and byte-identical unbound state remain deleted.

### Repaired acceptance criteria

**AC-1 (VALUE) — Every normal approved route is durably authorized before dispatch.** The real-CLI positive records and commits one selected Briefing, records and commits one approval, consumes and commits the successor, then performs exactly one successful `dispatch build` and produces exactly one later durable successor report/marker. Decision-before-bind count is **0**; build-before-consumed-descendant count is **0**; unauthorized advance/dispatch count across revise, hold, stale, blocked, and spent cases is **0**. The terminal complement has zero builds/effects and enters the existing merge path. *Verified by:* real CLI, logging wrapper, Git snapshots/ancestry, omission/reorder/duplicate controls, and one goal-only approved journey per supported host. No common native transport event is required.

**AC-2 — Authority, correction ownership, and Captain meanings are durable and deterministic.** Direct close records `person:captain`; delegated close records `agent:first-officer`, a nonblank evidence reason, and exactly the supplied argv-capable UTF-8 directive bytes excluding NUL. Reviewer `REJECTED` enters `feedback-rejection-flow` before the Captain lifecycle with zero Resolution/consume/dispatch. Approve, redo, routed reject, unrouted reject, hold, and `not yet` map and route as stated above; reasons preserve accepted direction, rejected direction, or pause condition. *Verified by:* public CLI plus production `gates.Document`, real Git bind/close/consume snapshots, shared-core branch-deletion control, mapping delete/swap controls, and deterministic routed/held snapshots. Live host matrices do not own language mapping.

**AC-3 (VALUE) — The Captain sees exactly one truthful root review before decision invocation.** On approved and no-authority paths, exactly one host-qualified root review occurs after the selected-Briefing folder commit and before `gate record --decision`. It semantically names the entity/title and stage, exact Briefing id and digest, recommendation, and decision ask. Child/tool/user output, pre-commit narration, post-decision summary, omission, or a second qualifying root review fails. Delegated conn never waives presentation. *Verified by:* presenter/lifecycle ownership controls, semantic keep/drop/mutate/duplicate cases, and existing host-specific Claude empty-parent, Codex root-agent interval, and Pi root-session assistant filters.

**AC-4 (VALUE) — A no-authority First Officer binds, presents, and stops without self-approval.** The guarded entity changes from no selected attempt to exactly one committed open attempt before presentation, then remains at the gated stage with no Resolution, consume, advance, dispatch, completion, or archive. The headless scenario independently proves drive to the gate under no conn; the primitive independently proves an already-engaged gate cannot self-approve. *Verified by:* refitted `assertGateHeld`, its unbound/advanced/Resolution/self-approved/noncanonical-review negatives, refitted Claude/Codex guardrail runners, standalone `live_gate_stop`, and host-neutral primitive coverage.

**AC-5 — The selected package and capability boundary fail closed.** The selected Briefing id/digest in review and Resolution equals durable state, one folder commit includes the index plus retained room files and excludes dirty siblings, and missing/incomplete capability or retained input halts before decision/consume/dispatch without manual frontmatter repair. *Verified by:* real-Git folder commit inspection, exact tuple mutations, existing CLI/gates refusal tests, and a structural capability-prose deletion control. Capability probe counts and binary-identity experiments are not acceptance requirements.

Source size and numstat are implementation guardrails only. They are intentionally absent from the value ACs.

### Repaired lean test plan

1. **Clean branch setup, before edits:** create the explicit local archival ref `spacedock-archive/first-officer-gate-command-lifecycle-3c535105` at `3c53510529c095878f5156d490d498e6742fb31c`, then `git switch -c spacedock-ensign/first-officer-gate-command-lifecycle-cycle27 13d702492131df17dd3ac87245d6d773f4df959b`. Do not reset, force, delete, or rewrite either ref. If either name already resolves elsewhere, stop for First Officer direction. No ref operation occurs during ideation.
2. **Direct CLI and Git (AC-1/AC-2/AC-5; medium):** retain bind/close/consume, omission, stale/blocked/spent, revise/hold, and terminal cases. Inspect the selected folder commit and exact tuple. Table-drive argv-capable UTF-8 directives excluding NUL through the public CLI; read `gates.Document` and assert bind, close, and consumed Git snapshots. The production parser replaces every test-only YAML/serialized-line mechanism.
3. **Shared contracts (AC-2/AC-3; low):** pin reviewer-`REJECTED` before ordinary Captain lifecycle, all six Captain call classes, presenter sole ownership, lifecycle wait-before-decision, semantic review minimum, and capability fail-closed prose. Deleting one branch, mapping, semantic member, multiplicity rejection, or barrier must make its focused control red.
4. **No-authority paths (AC-3/AC-4; medium/high):** refit the gate fixture with a retained Briefing and grade bound-open state. `assertGateHeld` requires one open binding and semantic review while forbidding Resolution/consume/advance/dispatch. Preserve separate shared guardrail, standalone headless no-conn, and primitive scenarios; do not accept legacy markers or unchanged bytes.
5. **Approved supported hosts (AC-1/AC-3; high):** after deterministic/full/race/docs/live-compile and independent staff approval, run one unchanged goal-only approved journey on Claude, Codex, and Pi. Common proof is selected commit → one semantic root review → decision/consume commits → one later build → one durable effect. Each runner’s existing host-specific extractor owns root/order filtering; separate host-native fixture suites own dispatch transport.
6. **Repository gates:** focused lifecycle/contract/no-authority controls, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, live compile, strict MkDocs, and `git diff --check`. No route-language live matrix, exact-child assertion, assignment-custody proof, prompt serializer, or new Pi retention harness.

No new spike is needed: `13d70249` already contains the real CLI, production parser, Git fixture, guardrail/primitive runners, and host-specific extractors. Implementation first refits deterministic and no-authority controls, records measured per-file deltas, and stops for staff review before live spend.

### Obligation-based surface ranges

These are planning ranges from the actual `13d70249` file sizes, not subtraction credit or acceptance value. A path outside the table, a new file, or a per-file delta outside its range pauses implementation for design review. Shared helpers may reduce deletion safely; no global deletion floor applies.

| Path | `13d70249` size | Expected implementation range |
| --- | ---: | ---: |
| `skills/first-officer/references/first-officer-shared-core.md` | 200 lines / 26,602 bytes | +1–3 / -1–3 |
| `skills/fo-gate-lifecycle/SKILL.md` | 56 / 5,788 | +4–8 / -1–5 |
| `skills/present-gate/SKILL.md` | 62 / 8,615 | +4–12 / -0–12 |
| `internal/contractlint/fo_function_reference_invariant_test.go` | 436 / 20,889 | +20–45 / -0–10 |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | 1,146 / 53,715 | +70–120 / -100–190 |
| `internal/ensigncycle/gate_assert_impl_test.go` | 39 / 1,380 | +15–30 / -10–25 |
| `internal/ensigncycle/gate_assert_test.go` | 62 / 1,741 | +10–25 / -20–45 |
| `internal/ensigncycle/live_gate_stop_test.go` | 199 / 9,805 | +12–24 / -130–180 |
| `internal/ensigncycle/livescenario_adapter_live_test.go` | 97 / 4,056 | +8–18 / -20–45 |
| `internal/ensigncycle/shared_fixtures_test.go` | 809 / 41,480 | +5–15 / -25–55 |
| `internal/ensigncycle/shared_scenarios_test.go` | 88 / 5,534 | +1–3 / -1–3 |
| `internal/ensigncycle/shared_scenarios_negative_test.go` | 477 / 27,457 | +5–15 / -10–30 |
| `internal/ensigncycle/claude_live_runner_test.go` | 615 / 29,564 | +5–15 / -5–25 |
| `internal/ensigncycle/codex_live_runner_test.go` | 483 / 20,943 | +5–15 / -5–20 |
| `docs/site/concepts/gates-and-decisions.md` | 80 / 6,393 | +8–18 / -8–22 |
| `docs/site/get-started/first-workflow.md` | 107 / 3,626 | +3–8 / -5–12 |
| `docs/specs/scenario-testing-principles.md` | 98 / 10,032 | +1–4 / -1–4 |

The measured center is approximately **+250/-489** across these 17 existing files; the abandoned cycle-23 implementation measured +497 additions across the same product area. This comparison is a review guardrail, not a mandate to delete useful proof. Production Go, host-native dispatch fixture files, `pi_shared_coverage_test.go`, and `shared_scenarios_meta_test.go` remain unchanged because the scenario names and transport owners remain intact.

### Reset triggers and non-goals

Stop before live work for any production Go/schema/provider/recorder change, new file/harness, common native dispatch parser, exact-child/model/cwd/session requirement, assignment custody, prompt coaching, route × host matrix, legacy marker/unbound compatibility, new Pi retention subsystem, ref rewrite/reset/force operation, or per-file range breach. A test that can pass on an unbound review, a build before consumed commit, two qualifying root reviews, or a changed directive is a design failure, not a reason to weaken the AC.

No code, test, live host, branch/ref, PR, gate, review-round, status, or merge mutation occurs during this repair. Independent staff re-review must approve the authority and surface corrections before the clean branch is created.

## Stage Report: ideation (cycle 27 repair)

- DONE: Restore the durable successor, correction/mapping, semantic review, bound-open guardrail, folder commit, capability, and directive-domain obligations requested by staff.
  The repaired ACs retain exactly one post-consume build/effect without native causality, deterministic rejection/mappings without host matrices, exact-one semantic root review with host-specific filters, and independent no-authority/self-approval coverage.
- DONE: Replace deletion-as-value with obligation-based implementation guardrails.
  Seventeen existing paths list actual `13d70249` line/byte sizes and per-file delta ranges; no global deletion floor or numstat AC remains.
- DONE: Preserve a non-destructive clean-branch implementation plan.
  Implementation first archives `3c535105` under an explicit local ref, then creates a fresh branch at `13d70249`; reset, force, ref deletion/rewriting, and all ideation-time ref changes are prohibited.
- SKIPPED: Edit/reset code, create refs/branches, run tests or live hosts, push PR, or mutate gate/round/status.
  The repair changes only this entity body in the shared state checkout.

### Summary

Cycle 27 repair restores product obligations that had no other integration owner while preserving the clean proof split. Direct CLI/Git tests own exact state, authority, mappings, and folder durability; supported-host journeys own one semantic root review plus durable route; host-native fixtures alone own transport. Checklist accounting is **3 DONE, 1 SKIPPED, 0 FAILED**; the repaired design is ready for independent staff re-review.
