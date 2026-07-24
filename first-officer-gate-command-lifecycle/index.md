---
id: 6yyyyemkqwsett3g1c991w9f
title: Make First Officers operate the recorded gate lifecycle
status: implementation
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
---

Make the normal First Officer gate path bind the exact reviewed package, record the authorized decision, and durably consume it before ordinary workflow dispatch.

## Problem

The recorder and application commands exist, but the shipped First Officer contract still assembles and presents gates only in prose. PR #557 (`fa240a76`) changed no First Officer skill files despite budgeting approximately ten lines for that integration; PR #560 (`f06cce04`) likewise added eligibility and consume without an FO caller. `skills/first-officer/references/first-officer-shared-core.md` names judgment, `present-gate`, feedback routing, advancement, and dispatch, but literally invokes none of `spacedock gate record`, `gate validate`, `gate eligibility`, or `gate consume`. This sprint therefore used the commands only under manual Captain directives, and an ordinary FO can still approve and advance outside the durable lifecycle.

The missing work is not a new decision model. It is the smallest procedure that makes the sprint FO's authority durable: retain a canonical package, bind its exact `briefing.json`, make and present the evidence judgment, record the Captain or delegated decision, consume that one-use application, durably commit the consumed state, then return to the ordinary dispatch loop. Before h1 landed, the same dogfood stopped after recording the Resolution and manually changed workflow state. `gate consume` is the sole approval transition: it re-evaluates eligibility under the entity lock and atomically co-writes the next `status` and `application.state: consumed`. A separate `status --set` advance is a bypass and a contract violation.

The repository also exposed two CLI-boundary defects. First, the retained repo-root `./spacedock` identifies as `0.26.0+dev` yet lacks the gate surface. Version compatibility is therefore insufficient. One `spacedock gate --help` response already carries all command and semantic-form tokens needed to reject that stale launcher; four repeated subcommand probes do not add discrimination. Second, the CLI passes a relative retained-input path to a recorder that compares it with an absolute entity root, so `filepath.Rel` fails before mutation. The CLI must normalize `--briefing`, `--result`, and `--association` against its launch directory; callers should not manufacture absolute paths.

## Demonstrated baseline

The design adopts the landed behavior and historical dogfood; it does not reopen whether durable gates should exist.

1. The 3k validation room retained a canonical `review/validation/briefing-v1/briefing.json` plus a concise gate review and frozen references. State commit `2c616b7e` bound it; `77590ebd` recorded an `agent:first-officer` approval carrying the Captain decision as adoption provenance. The resulting record stayed at `status: validation`: 3k intentionally recorded the decision but did not apply it.
2. The remaining 3k closeout used manual lifecycle mutations (`mod-block`, PR state, and later terminal archive) because h1 was not yet available. Those commits are historical evidence for the exact integration gap, not a procedure to preserve.
3. h1 landed `gate eligibility` and `gate consume`. Its CLI fixture proves an approve closure yields `advance/pending`, eligibility prints `condition=approved-pending eligible=true`, and consume atomically changes `status` plus `application.state: consumed`; a second consume exits nonzero with `condition=consumed consumed=false` and byte-identical state.
4. The recorder's retained package boundary is settled: the manifest basename is exactly `briefing.json`; its canonical bytes independently define the complete artifact inventory. Existing artifact payloads remain URI + SHA references when the presentation resolver can reproduce those exact bytes. Mutable, cross-root, or otherwise unreproducible reviewed material is frozen as a room copy before publication. The recorder verifies but never copies artifact payloads. For folder-form entities, landed vn/PR #558 makes one `spacedock state commit <slug>` include the index and every non-ignored room artifact without sweeping sibling dirt.
5. Direct and delegated decisions remain distinct. A Captain who personally renders the chat decision is `person:captain`. An FO acting under delegated conn is `agent:first-officer`, supplies its evidence-bearing `--reason`, and stores the exact quoted grant in `--directive`. An exact provider Result uses `--result` plus its retained association and authorized actor; advisory adoption also names the authorizer.

The cycle-6 cardinality spike exercised the landed CLI/gate tests and inspected the write boundaries. Both semantic `gate record` forms validate the complete document and transition before atomic replacement and return the selected summary; a later close reads and validates the bound record again. `gate consume` reads and validates the record, re-evaluates reviewed-input currency, expected successor, blockers, and one-use state under the same entity lock, then atomically writes status plus consumed state. Therefore open/closed `gate validate` and pre-consume `gate eligibility` are diagnostic reads, not independent authority checks: deleting them leaves the three mutation guards intact, while deleting briefing record, decision/Result record, or consume destroys a required durable state transition. Focused `internal/gates` and `internal/ensigncycle` tests passed this baseline on 2026-07-24. Implementation turns this spike into the first three-command fixture and deletion controls.

## Operational lifecycle

The procedure belongs at the existing gated-stage branch in `first-officer-shared-core.md`, around `«gate.assemble-verdict»`; it does not move or duplicate `«gate.ac-cross-check»`, FO evidence judgment, or `present-gate` rendering.

### 0. Select a capable executable

Before the first gate mutation in an FO session, resolve `${SPACEDOCK_BIN:-spacedock}` once and invoke only:

```text
${SPACEDOCK_BIN:-spacedock} gate --help
```

Require one successful response containing `record`, `validate`, `eligibility`, `consume`, and the record forms/flags `--briefing`, `--result`, `--association`, `--decision`, `--actor`, and `--directive`. This is the minimum current capability fingerprint: all tokens come from the one gate help body, so four subcommand-help calls distinguish no additional stale binary.

Cache the pass only in the current FO session, keyed by the executable identity resolved immediately before use: canonical absolute target path plus a content digest of that target. Re-resolve and re-digest before reusing the cache. An unchanged identity reuses the pass; a same-path replacement, symlink retarget, different PATH resolution, changed `SPACEDOCK_BIN`, or unknown-command/unknown-flag result invalidates it and requires one new help probe. Launcher text alone is insufficient because the file at that text can change; path plus stat metadata is also insufficient because a replacement can preserve size/time. The content-key serves AC-5's pre-mutation stale-launcher value with no workflow-persistent cache.

If the fingerprint is incomplete, halt before mutation and prescribe refresh or a fresh checkout build selected through `SPACEDOCK_BIN`. Never hand-edit `gates:`. The CLI, not the FO, normalizes retained `--briefing`, `--result`, and `--association` arguments against the launch directory before recorder entry; absolute paths remain accepted.

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

Only after the open binding succeeds and is committed does the FO perform `«gate.ac-cross-check»`, make its evidence judgment, and invoke `present-gate`. The presenter remains an overridable channel; this task does not redesign it. Regardless of channel, the Captain sees the concise primary review rather than a raw entity dump, spec, room listing, or `briefing.json`. Those remain linked references.

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

The close call must exit 0, name the bound attempt/Briefing, report `state=closed`, and reproduce the intended decision. Its internal write path validates the existing record, authority fields, Resolution, derived application, transition, and rebuilt record before atomic replacement; a separate closure `gate validate` adds no authority. Every successful close is then committed through `state commit` before any route—approve, revise, hold, or a consume attempt. This close commit is a durability barrier: it preserves the exact Resolution even if consume later refuses or the route stops. A failed close or failed close commit halts without feedback, advancement, or dispatch.

### 4. Route the closed decision fail-closed

The Resolution is recorded before any route:

| Closed decision / condition | Required FO route |
| --- | --- |
| `approve` | After the close commit, invoke `gate consume` directly. Only exit 0 plus `condition=approved-pending eligible=true consumed=true` and the expected immediate successor is success. Consume itself evaluates the authoritative condition under lock, atomically advances `status`, and marks the application consumed. Commit that mutation, then hand the newly current stage to the existing ordinary reuse-or-fresh dispatch procedure. No separate `status --set` advance is allowed. |
| `revise` | After the close commit, never consume it as an advance. Invoke the existing `«feedback.route»` procedure; this task does not implement advisory correction rounds. `gate eligibility` remains available only when the condition needs diagnosis. |
| `hold` | After the close commit, leave the entity at the gate and surface the reason; never consume, advance, or dispatch. `gate eligibility` remains an optional diagnostic read. |
| approved but blocked, held, unknown, wrong-stage, or otherwise ineligible | The Resolution is already durable in the close commit. Halt advancement and dispatch, preserve status bytes, and name the exact reported condition and missing/current artifact or field. |
| `stale` | The consume attempt reports stale, exits nonzero, leaves status unchanged, and changes the pending application only to `superseded`. Commit that supersession, retain/bind a new Briefing, and re-present; never expose a separate rebind/supersede ceremony to the Captain. |
| already `consumed` | Treat the approval as spent. Do not record another decision or consume again; follow current `status` into the ordinary dispatch/recovery loop. A diagnostic repeat consume is nonzero and byte-clean. |

After successful consume, its second `state commit` is the approval-effect durability barrier: no successor `dispatch build` or host dispatch may occur until a descendant of the close commit contains `application.state: consumed` and successor status. After that barrier, the generic gate lifecycle ends. The ordinary dispatch contract owns build/reuse/fresh dispatch, and each host adapter owns how spawn or completion is observed; this procedure adds no “very next host event,” returned-handle, or common transport requirement.

### 5. Resume and retry without minting duplicates

Use `status --boot --identify --json`, entity state, and the previous command result to choose the next semantic operation. `gate validate` and `gate eligibility` remain on-demand diagnostics for malformed or ambiguous state; neither is a mandatory prelude:

- open attempt + same retained Briefing: repeating `record --briefing` is an idempotent bind; continue presentation;
- open attempt + changed package: `record --briefing` selects the changed binding under the same attempt; the FO speaks only of updating the presented Briefing, not recorder mechanics;
- closed + pending approval: do not re-record the decision; first ensure the exact closure has a durable state commit, then consume, which authoritatively accepts or refuses;
- closed + revise/hold: first ensure the exact closure has a durable state commit, then route or remain held; never call advance-consume;
- consumed + already advanced status: dispatch/recover the current stage without consuming again;
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
| Version-compatible repo-root `./spacedock` lacked the gate subcommands. | In scope: one `gate --help` fingerprint cached per resolved canonical target path + content digest plus fresh-build/refresh route before mutation. AC-5. Same-path replacement and PATH-target swap invalidate the cache; repeated subcommand probes remain removed. |
| Fresh `/tmp/spacedock-current gate record first-officer-gate-command-lifecycle --briefing docs/dev/.spacedock-state/first-officer-gate-command-lifecycle/review/ideation/briefing-1/briefing.json --workflow-dir docs/dev` failed before mutation with `resolve briefing room: Rel: can't make ... relative to ...`. | In scope CLI fix: normalize relative `--briefing`, `--result`, and `--association` paths against the launch directory before recorder entry. AC-5 makes the natural relative form succeed and keeps invalid/missing inputs byte-clean; absolute paths remain compatible. |
| The presentation helper could report failure before a late provider Result arrived, and earlier cleanup destroyed result bytes. | Boundary: xb/presentation transport owns exact Result and association retention; this procedure accepts only retained bytes. AC-4's provider arm uses retained inputs. No polling or transport enters this task. |
| Provider and canonical Briefing ids differed even when primary artifact bytes matched. | Boundary: landed recorder verifies complete association and normalizes identity only after exact bytes/authority. The FO supplies Result + association; it does not normalize by hand. AC-4. |
| `status --set` reserialized hand-authored `gates`, breaking anchors. | Prohibited workaround: binary owns all gate writes; `status --set` never substitutes for recorder or consume. AC-2/AC-5. |

## Acceptance criteria

**AC-1 (VALUE) — A normal FO-approved gated stage cannot dispatch without one durably consumed authorization.** In the real 3k validation-package replay, the FO emits the exact ordered authority trace (briefing record, delegated decision/Result record, consume), commits the retained package before presentation, commits every successful close before routing, and commits the consumed successor as a descendant before exactly one successful successor `dispatch build`. On-disk state ends at the expected successor with the same attempt's exact Briefing, Resolution, and `application.state: consumed`; the successor stage produces one new durable marker-bearing report/commit that was absent before dispatch. The measured count of advances or dispatches without a matching consumed application commit is **0**, versus the current prose contract's possible manual path of 1. *Verified by:* deterministic real-CLI replay plus one representative engaged live FO journey on each supported host, whose command log names the close/consume commit hashes and entity snapshots; the dispatch pre-HEAD must equal/descend from the consumed commit, and state Git history/marker cardinality prove one later durable effect. After that barrier, the ordinary dispatch contract owns host handoff.

**AC-2 — The minimum authority sequence is load-bearing and fails closed.** Omitting briefing record, decision/Result record, or consume makes the scenario grader refuse successor dispatch. Removing the former open validate, closed validate, and eligibility reads from the positive trace does not weaken the recorder/consumer guards: malformed open/close writes fail inside `record`, and stale/blocked/wrong-stage/spent approvals fail inside locked `consume`. Revise, hold, and ineligible controls produce zero advance/dispatch effects; revise invokes feedback routing and neither revise nor hold is consumed as advance. *Verified by:* three existing-fixture controls that each execute the remaining sequence through the freshly built real binary and logging wrapper (not a mutated expected slice), assert actual exits/state/commits, and offer an otherwise-valid dispatch observation that the oracle rejects; malformed-record and consume-refusal fixtures remain direct controls.

**AC-3 — Direct and delegated authority remain durably distinguishable and every successful decision close survives a stopped route.** A direct Captain decision records `by: person:captain`; a delegated FO decision records `by: agent:first-officer`, a nonblank evidence reason, and the exact quoted Captain conn as directive/adoption provenance. Approve, revise, and hold each produce a close commit before consume/feedback/stop; an approval whose consume is later refused retains that exact Resolution commit. Missing delegated fields fail before closure and leave the entity byte-identical; a procedure mutant that records the delegated live decision as `person:captain` fails the provenance grade. *Verified by:* public CLI fixtures, delegated live baseline, actor-swap mutant, and real-Git approve-refused/revise/hold controls asserting commit hash, closed snapshot, unchanged status, and no dispatch.

**AC-4 — The FO package is complete, reproducible, and durable without unnecessary copying.** The replay uses the retained 3k three-artifact validation Briefing: every URI resolves through the declared resolver and every SHA matches; reproducible existing artifacts are not duplicated, while the frozen reviewed snapshot stays in the room. One folder-form `state commit` includes the index and all new room files and excludes dirty sibling paths. Fixtures live only in testdata/temp locations and do not alter repo workflow discovery. *Verified by:* exact digest/path assertions, a real-Git state-commit assertion, and before/after workflow-discovery candidate equality.

**AC-5 — Readiness and retained-input paths work at the CLI boundary while failures preserve landed write boundaries.** One `gate --help` call rejects a version-compatible but capability-stale executable before room/entity mutation and names refresh/fresh-build remediation; the pass is reused only while canonical resolved target path and content digest remain identical in the current FO session. A same-path replacement and a PATH-target swap each force a new probe and reject a stale target. A repo-root-relative `--briefing` binds the same canonical bytes as its absolute form, and relative `--result`/`--association` inputs reach semantic validation rather than failing mixed-root normalization. Missing/noncanonical Briefing and invalid association/actor/provenance name the missing artifact/step, return nonzero, preserve entity bytes, and leave no lock residue where promised. Validate/eligibility diagnostics are byte-clean; hold/blocked/repeat-consume refusals are byte-clean; stale consume changes only pending → superseded and never status. *Verified by:* one-help identity-keyed counting cache, capable→stale same-path replacement, two-directory PATH resolution swap, CLI relative/absolute equivalence fixtures, and the existing whole-file/tree-hash refusal matrix.

**AC-6 — The Captain receives a concise decision review, not recorder mechanics or raw artifacts.** The live output contains capability/change, test/evidence, reviewed snapshot, findings, one recommendation, and a concrete decision ask; entity/spec/Briefing are linked references. It does not lead with raw JSON/YAML, room listings, or `rebind`/`supersede` ceremony. *Verified by:* the live scenario's structured gate-review grader and retained output artifact; a raw-file-dump mutant fails the grade.

**AC-7 — Resume is idempotent and one-use without mandatory diagnostic ceremony.** Open same-Briefing retry creates no duplicate attempt; a closed pending approval is not re-recorded and is committed before consume if a prior session stopped between close and commit; consumed state never consumes twice; stale state becomes superseded without advancing and requires a replacement Briefing. Across three fresh-process resume passes, the count of close commits, successful consumes, and resulting stage transitions is exactly 1 each, whether or not an operator invokes optional `validate`/`eligibility` diagnostics. *Verified by:* fresh-process fixture passes, Git ancestry/hash assertions, and byte comparisons after each retry.

**AC-8 — The procedure is runtime-portable, fixture-backed, and transport-neutral.** The recorded-gate-lifecycle journey remains one host-neutral scenario with Claude, Codex, and Pi runners. Each host executes one representative engaged route and must satisfy the same observable oracle: three ordered successful gate mutations, a close commit, its descendant consumed-successor commit before exactly one `dispatch build`, and exactly one later durable successor effect. The logging shim records commit hash and entity snapshot after each state commit and immediately before dispatch, making the ordering falsifiable without transport forensics. After build, each host follows its ordinary dispatch adapter; the generic lifecycle requires no next-event ordering, returned handle, common spawn schema, or exact-child attribution. Static contract tests prove every gate-entry route funnels through the lifecycle; route-by-host live multiplication is not required. *Verified by:* the three live journeys, AC-2's three real-binary omitted-mutation controls, `TestFOGateLifecycleOwnsEveryEngagedEntry`, and existing Claude/Codex/Pi native dispatch fixtures.

## Minimum replay and test plan

The test package copies the real retained 3k validation Briefing and its declared artifacts into a temp folder-form workflow whose gated `validation` stage has a supported nonterminal successor. It preserves the package bytes and uses the historical delegated-approval shape: `--actor agent:first-officer`, `--reason "Captain directive: approved after reviewing the presented 3k validation gate."`, and `--directive "you have the conn toward the sprint goal; authorized to approve gates, PR, relevant CI lanes, and merge; use your judgement."`. The existing command-logging wrapper delegates to the freshly built real binary and records exit/stdout/stderr. Extend that same log with the state HEAD hash plus entity snapshot digest after each `state commit` and immediately before `dispatch build`; the existing dispatch stub is permitted only when the real trace and Git ancestry prove the close and consumed barriers.

1. **Baseline live journey (AC-1/AC-3/AC-6/AC-8; high):** launch the real FO, engage the held validation gate, grade its concise evidence review, supply/retain the delegated approval, and assert three ordered gate mutations. The bind commit snapshot is open; the close commit snapshot contains the exact Resolution/pending application; the consumed commit is its descendant and contains advanced+consumed state; dispatch's pre-HEAD equals or descends from consumed; one later commit adds the marker-bearing successor report. Run this representative route once on Claude, Codex, and Pi; ordinary host fixtures continue to own transport.
2. **Minimum-sequence proof (AC-2; medium):** for each omission, use the existing fixture/wrapper and freshly built real binary to invoke every remaining authority command in order, record its actual exit/state/commit trace, and offer a prospective dispatch observation. Omitted bind makes close/consume fail, omitted close makes consume fail ineligible, and omitted consume leaves pending state; every arm denies dispatch. Do not satisfy this with expected-slice deletion. Keep the positive free of mandatory validate/eligibility and retain direct malformed-record/consume-refusal fixtures.
3. **Close durability and non-approval controls (AC-2/AC-3/AC-5; medium):** approve-then-block/refuse, revise, and hold each record a real decision and commit its exact closed snapshot before the later route. Revise routes feedback with zero advance-consume; hold and refused approval remain at the gate; stale adds only a later supersession commit. Assert close hash/ancestry, status, Resolution bytes, and zero dispatch.
4. **Provenance, package, and path matrix (AC-3/AC-4/AC-5; medium):** direct versus delegated versus exact Result inputs; complete versus truncated association; canonical versus alternate basename; relative versus absolute `--briefing`, `--result`, and `--association`; reproducible URI+SHA versus frozen copy. Compare bytes and lock residue on every refusal, and prove relative/absolute forms bind or adopt identical bytes. Add the narrow normalization assertion in existing `internal/cli/gate_test.go`; do not change recorder path semantics.
5. **Resume passes (AC-7; low):** rerun fresh processes over open, closed-but-uncommitted, pending-committed, stale, and consumed snapshots without mandatory diagnostic reads; exactly one close commit, transition, and consume across three passes. Add an optional-diagnostic arm only to prove reads do not alter the result.
6. **Capability identity cache (AC-5; low):** use the existing stale shim/counting probe. An unchanged canonical path+content digest yields one help call across two gates. Replace the capable executable at the same path with a stale shim and assert digest change, reprobe, and pre-mutation refusal. Then place capable/stale same-named launchers in two PATH directories, swap resolution order, and assert canonical target change, reprobe, and refusal. `SPACEDOCK_BIN` identity change uses the same invalidation rule.
7. **Runtime lanes (AC-8; high):** keep the host-neutral scenario once with Claude, Codex, and Pi bindings. Each live lane retains command logs, gate review output, before/after entity, state Git history, and the one durable successor effect. Host-native fixtures separately prove each adapter's transport event shape; live lanes do not require public event uniformity or every route on every host.
8. **Repository gates:** `gofmt -w ./cmd ./internal`, focused scenario/skill tests, `go test ./...`, `go test ./... -race`, live-tag compilation, and the required live lanes.

The primary independent live oracle is durable state plus the command/dispatch log and state Git history. It resolves the logged close and consumed hashes, verifies their entity snapshots and ancestry, and compares dispatch's pre-HEAD; substring order alone cannot pass. Static command names in skill prose are inspection evidence only. The three real-binary omission controls prove authority writes are load-bearing; existing gate-unit fixtures prove diagnostic reads are unnecessary to the write guards.

Reuse the current `recorded_gate_lifecycle_test.go`, command wrapper, 3k fixture, dispatch oracle, state-Git oracle, topology tests, and one live binding per host. Change `recordedGateRequiredEvents` from six to three; add close/consumed commit hash+snapshot records to the wrapper/oracle; replace slice-only deletion with three real-binary subtests; remove positive validate/eligibility, repeated probes, and absolute-path workaround expectation. Extend existing `internal/cli/gate_test.go` for retained-input normalization. Do not add a second scenario, harness, live matrix, transport parser, or provider runtime.

## Obligation delta

| Obligation | Authority | Bearer | Proof burden |
| --- | --- | --- | --- |
| Three ordered authority mutations; one consumed approval; zero unauthorized advances/dispatches | Canonical AC-1/AC-2 and `fo-gate-lifecycle` | First Officer lifecycle contract plus landed recorder/consumer guards | Real CLI fixture, three real-binary omitted-mutation controls, malformed-write/consume-refusal fixtures, and each host's representative live command/state trace |
| Exact package committed before presentation; every successful close committed before routing; consumed successor committed before dispatch | Canonical AC-1/AC-3 | First Officer write contract | Logged commit hashes/snapshot digests and ordered state Git ancestry showing all durability barriers and no earlier dispatch build |
| Exactly one successor dispatch and one new durable successor effect | Canonical AC-1 | Existing ordinary dispatch loop after durable consume | One post-commit `dispatch build`, pre/post entity bytes, marker/report cardinality, and state Git history in each live journey |
| Every engaged route loads the lifecycle before gate action | Deferred topology contract | Shared First Officer core | `TestFOGateLifecycleOwnsEveryEngagedEntry` and native load-before-action fixtures; not every route live on every host |
| Host transport shape | Canonical AC-8 and each supported runtime adapter | Claude, Codex, and Pi adapter/fixture owners | Host-native fixtures, including Codex's supported build → wait → durable-report evidence when public spawn records are absent |

Captain rulings of 2026-07-24 remove five non-authority burdens from the generic lifecycle: mandatory open/closed validation reads, mandatory pre-consume eligibility, repeated subcommand help probes, absolute-path caller ritual, and runtime-specific next-event/handle language. Earlier proof-boundary cuts also remove uniform public exact-child forensics and the route × host live matrix. No authority outcome is removed: the FO still owes exact package binding, exact decision provenance, locked one-use consume, durable consumed state before ordinary dispatch, one dispatch, and one durable effect. This canonical AC block, test plan, and obligation delta are the sole active proof authority. Later historical reports do not compete with it.

## Expected surface and tolerance

- **FO contract:** rewrite the existing `skills/fo-gate-lifecycle/SKILL.md` in place and adjust its one shared-core pointer in `skills/first-officer/references/first-officer-shared-core.md`. Add the every-close commit barrier and resolved-path+content-digest cache identity while removing mandatory diagnostic reads, repeated probes, absolute-path ritual, and runtime handoff prose; do not add a skill, host adapter wording, or presentation template.
- **Narrow CLI fix:** `internal/cli/cli.go` plus `internal/cli/gate_test.go`, approximately 10–25 production lines and 25–60 test lines, normalizing only retained gate input flags against the injected launch directory. No recorder/schema change.
- **Codified/live proof:** reuse `internal/ensigncycle/recorded_gate_lifecycle_test.go`, its existing 3k fixture/wrapper/oracle, topology/contract tests, and current Claude/Codex/Pi bindings. Replace six expected events with three, add commit-hash/snapshot ordering, drive three omissions through the real binary, add two cache-identity controls, and delete obsolete positive diagnostics/probe repetitions/path workaround. No second harness. From preserved implementation commit `86bad049` (1,468 added LOC vs `origin/main`), additions are capped at 170 and must be offset by at least 50 deletions; total branch surface is expected at or below **1,588 added LOC**. This is a tolerance boundary, not a prose-compression target.
- **Docs:** revise the already-touched `docs/site/concepts/gates-and-decisions.md` and `docs/site/reference/command-reference.md` to describe three authority mutations, every-close/consumed commit barriers, identity-keyed help cache, CLI path normalization, and ordinary dispatch handoff; approximately 10–24 changed lines.
- **Component budgets and attribution:** shared first-officer core remains **≤26,754 bytes** and deferred lifecycle remains **≤6,600 bytes**. For every host, implementation records the current-main load manifest, then reports the byte delta attributable only to this task's touched load-set files; that attributable delta must be non-positive. Whole-host totals remain measurements, not fixed ceilings.
- **Rebaseline rule:** at every implementation commit selected as a test/review tip, and again after every rebase, measure each task-touched load file at the current upstream parent and at the candidate tip, then recompute the per-host attributable delta as `tip bytes - current-upstream bytes`. If upstream changed a task-touched file, first replay the task patch, then recompute against that new upstream version; never carry forward stale arithmetic. Informational whole-host baselines may absorb only bytes present in committed current upstream, including its changes to untouched or overlapping files before the task patch. Never absorb candidate/task-owned bytes, an uncommitted concurrent edit, or a component-cap breach. Any positive recomputed task delta or cap breach requires ideation reapproval; unrelated upstream growth does not trigger semantic compression here.
- **Tolerance:** 2× for changed files/LOC. Reconfirm if implementation needs more than the narrow CLI normalization, a schema/recorder/application change, another gate/state command, package generator, presenter UI, provider polling/transport, retry engine, or new test framework.

## Documentation change proposal

After “The three calls” in `docs/site/concepts/gates-and-decisions.md`:

```diff
+Before the First Officer shows a gate, it records the exact retained Briefing and commits
+that package. After the decision, it records and commits the Resolution before any route,
+including revise, hold, or a refused approval. An approval advances only when `gate
+consume` rechecks and spends that authorization; consume writes the new stage and consumed
+mark together, and a descendant commit is durable before ordinary dispatch begins.
+`gate validate` and `gate eligibility` remain on-demand diagnostics, not approval steps.
+Revise routes feedback instead, and hold stays at the gate.
+The review itself remains concise: capability, evidence, reviewed snapshot, findings,
+recommendation, and decision ask. The entity, spec, and package are linked references.
```

After the setup compatibility paragraph in `docs/site/reference/command-reference.md`:

```diff
+For source-checkout or retained development launchers, version compatibility is not a
+command-capability check. Cache one successful `spacedock gate --help` per session only
+while the resolved canonical executable path and content digest are unchanged; same-path
+replacement or different PATH resolution forces a new probe. Help must list `record`,
+`validate`, `eligibility`, and `consume` plus the semantic record flags. Otherwise refresh
+or rebuild the launcher; do not hand-edit gates. Relative retained-input paths resolve
+from the launch directory, so callers need not manufacture absolute paths.
```

## Dependencies, collision, and non-goals

- 3k/PR #557 and h1/PR #560 are landed, binding inputs. vn/PR #558 supplies folder-room durability. xb owns presentation-channel transport/retention and can remain a reference dependency; this task keeps `present-gate` judgment/rendering intact.
- The concurrent advisory-round generalization explicitly owns `docs/dev/README.md`, `skills/feedback-rejection-flow/SKILL.md`, recorder/schema surfaces, and a round-only caller. It declares no `first-officer-shared-core.md` change, so the intended FO contract surfaces do not collide. Both tasks may edit `docs/site/reference/command-reference.md`; serialize that doc edit (land/rebase one before the other). If either implementation expands into the other's skill file or ordinary gate procedure, stop and serialize rather than resolving two concurrent meanings.
- No schema or recorder/application logic; no gate judgment duplication; no presenter redesign; no advisory correction-round implementation; no blocker evaluator or hold authoring; no arbitrary artifact registry; no provider launch/polling; no new gate verb; no hand-authored `gates:` repair; no stable-release work.

## Minimum value demonstration seed

In one fixture-backed live First Officer journey, package an exact validation Briefing, record and commit it, record and commit an evidence-bearing delegated approval, consume and commit the application, and only then enter ordinary dispatch. The log must bind each barrier to a Git hash and entity snapshot. Revise, hold, and refused approval retain their close commit without advancing. Three omission controls run the remaining calls through the real binary and fail dispatch; removing diagnostic `validate`/`eligibility` reads from the positive remains green because their safety checks live inside mutation commands.

## Boundary

This task owns the First Officer invocation contract, the narrow CLI retained-path normalization, and their behavioral proof. It reuses the landed 3k recorder and h1 consume semantics; `validate`/`eligibility` stay public diagnostics. It does not add a recorder schema, duplicate gate judgment, change presentation UI, alter host dispatch adapters, or implement advisory correction rounds.

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

### Feedback Cycles

- Cycle 2: CHANGES REQUESTED — Roborev job 700; surface 7 test/live files and 584 added LOC vs estimate 7 files and at most 615 added LOC (95%); AC unchanged

### Roborev follow-up triage (cycle 2, job 708)

- **Material — AC-7 resume compares only entity bytes.** Released user/workflow: an FO replaying any open/closed/consumed gate after restart; harm: the replay may create or mutate sibling state while entity-only equality remains green; AC/boundary: AC-7 exact pass cardinality and byte preservation; trigger: same/open, revise, hold, and consumed arms compare only `fixture.entity`. Fix every idempotent arm against a complete state-root snapshot.
- **Material — topology proof remains structural, not runtime-observed.** Released user/workflow: every engaged gate route on Claude, Codex, and Pi; harm: a host may fail to load the deferred skill before gate commands while the Markdown-derived trace remains green; AC/boundary: AC-1/AC-2 and the cycle-2 load-topology stop condition; trigger: `TestRecordedGateLifecycleLoadTopologyMatrix` reads only shipped Markdown. Retain it as structural proof, add host-stream load/command extractors with all route fixtures and red controls, and make each live recorded-gate runner grade its observed load order.
- **Material — successor controls bypass the real stream parsers.** Released user/workflow: post-consume successor dispatch on every host; harm: a parser regression could accept narration, a blank handle, parent-written output, or an empty wait; AC/boundary: AC-1 strict spawn/handle/correlated-output oracle; trigger: `TestRecordedGateLifecycleSuccessorOracleControls` constructs `recordedGateDispatchProof` directly. Fix by routing adversarial Claude, Codex, and Pi streams through their production test extractors.
- **Polish — prompt-load growth note.** Released user/workflow: none; harm: none, because exact baselines and hard ceilings are already measured and the boot core shrank; AC/boundary: reporting clarity only; trigger: the rebaseline comment does not explicitly contrast full-lifetime growth with boot residency. Decline this round; promote if a reviewer or release note conflates the two metrics.
- **Polish — byte-clean helper name includes stale caller.** Released user/workflow: none; harm: none, because the stale caller separately asserts the exact allowed mutation; AC/boundary: test naming only; trigger: `assertRecordedGateByteCleanFailure` checks nonzero/output/lock, not bytes itself. Decline; promote if a caller relies on the helper name without its own byte assertion.

### Feedback Cycles

- Cycle 3: CHANGES REQUESTED — Roborev job 708; surface 7 test/live files and 708 added LOC vs estimate 7 files and at most 615 added LOC (115%, below 2× tolerance); AC unchanged

### Intended-change amendment 4

- `internal/ensigncycle/claude_live_runner_test.go`, `internal/ensigncycle/codex_live_runner_test.go`, and `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go`: add one runtime-load-order assertion per positive recorded-gate runner, using the shared host-stream extractors added to `recorded_gate_lifecycle_test.go`.

This adds at most 12 live LOC to make the previously declared all-host topology proof load-bearing. It changes no prompt, scenario, host-specific lifecycle instruction, or acceptance criterion.

### Roborev final triage (cycle 2, job 711)

- **Material, fixed — runtime action detection could accept echoed skill prose.** Released user/workflow: every live recorded-gate route; harm: an echoed `gate --help` inside a skill read could false-green load order; AC/boundary: AC-1/AC-2 observed load-before-action proof; trigger: the extractor matched any JSON line containing the text. Fixed by requiring each host's structured Bash/command-execution event and recognizing the first actual `spacedock gate` command, not only help.
- **Needs decision — route labels do not execute distinct runtime paths.** Released user/workflow: headless conn/no-conn, engage, worker completion, and every resume state; harm: a host-specific route can bypass the load while synthetic route fixtures remain green; AC/boundary: the binding cycle-2 topology matrix and stop condition; trigger: every route subtest currently reuses one normalized load/action fixture. The required actual-host matrix cannot be made green in this cycle: Claude is 401-blocked, Pi cannot load its extension dependency, and Codex omits the required structured spawn event. Do not substitute fixtures or exceed the 650-LOC proof cap; return for captain/FO reset or repaired live hosts.
- **Needs decision — copied-skill command deletion lacks green runtime execution.** Released user/workflow: a live FO consuming the shipped deferred skill; harm: a host could ignore a missing command while text-derived checks pass; AC/boundary: AC-2 and the copied-plugin live-mutant stop condition; trigger: deterministic deletion uses `procedureEvents`, while the actual Claude copied-plugin mutant cannot start because of 401 credentials. The original six-event baseline and all six deterministic deletions are now explicit, but the required runtime mutant remains blocked pending credential repair; do not claim validation readiness.
- **Material, fixed — deterministic missing-event discriminator was removed.** Released user/workflow: the common successor grader; harm: an event-completeness regression could escape the more specialized source/runtime checks; AC/boundary: AC-2 six-event load-bearing proof; trigger: job 711 observed the removed deterministic control. Restored a compact six-arm control that requires the grader's event-completeness error.

### Feedback Cycles

- Cycle 4: NEEDS DECISION — Roborev job 711; surface 7 test/live files and 741 added LOC vs estimate 7 files and at most 615 added LOC (120%, below 2× tolerance); AC unchanged, actual-route/live-mutant proof blocked

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

### Feedback Cycles

- Cycle 5: CHANGES REQUESTED — Roborev job 744; resumed-cycle surface 4 files and 31 added LOC against the bounded 650-line behavioral and 59-line live caps; AC unchanged; 8n projection is an explicit integration dependency

### Roborev final resumed-cycle triage (job 775)

- **Material, accepted — reason validation must reject canonical empty scalars.** Released user/workflow: delegated FO approval; observable harm: whitespace or quoted-empty reason could pass as evidence; affected value/boundary: AC-3 nonblank evidence judgment; trigger: job 775 observed the relaxed substring check. Anchor the unique canonical Resolution `reason:` field and trim whitespace plus YAML quote characters before accepting it. Do not restore one historical sentence: the contract deliberately permits the FO's evidence judgment while exact actor and Captain directive remain pinned.
- **Needs decision/blocker — a two-part mutation would not prove the command alone load-bearing.** Released user/workflow: the required copied-skill missing-command red control; observable harm: removing explanatory authorization prose with the command manufactures a red without isolating command deletion; affected value/boundary: AC-2 and the explicit command-only live-mutant stop condition; trigger: job 775, plus serial artifact `/tmp/spacedock-fo-gate-live-claude-mutant-serial`, where Claude loaded the copied command-only mutant, reconstructed eligibility from remaining instruction, emitted all six events, consumed, and dispatched. Revert to command-only deletion and retain this as a real blocker; do not weaken the criterion or claim green.
- **Low, accepted narrowly — prove the copied skill body as well as its directory.** Released user/workflow: mutant provenance; observable harm: another file in the copied directory could theoretically be consumed; affected value/boundary: AC-2 runtime mutant; trigger: job 775. Claude's native load event reports the base directory rather than `/SKILL.md`, so require that exact directory and absence of the removed placeholder command in the emitted loaded body.
- **Low, declined — semicolon-joined multiple gate commands are conservatively discarded.** Released user/workflow: none in the controlled route; observable harm: false negative only, never false green; affected value/boundary: AC-1 successful-event attribution; trigger: a host emits multiple gate calls joined only by `;`. Promote if a supported route intentionally emits per-command exit capture for that shape.
- **Low, declined — detect alteration from one historical reason sentence.** Released user/workflow: delegated approvals legitimately carry the FO's current evidence judgment; observable harm: pinning the historical sentence would reject correct live decisions; affected value/boundary: AC-3 requires nonblank evidence plus exact actor/directive, not exact prose; trigger: a binding specification changes to require a canonical reason string.

### Feedback Cycles

- Cycle 6: NEEDS DECISION — Roborev job 775; behavioral proof remains capped at 650 added LOC and live proof at 59; AC unchanged; command-only copied-skill mutant has a credentialed blocking counterexample
- Cycle 7: REVISE — delegated proof-boundary ruling; surface 650 behavioral plus 59 live LOC vs caps 650 plus 59 (100%); AC unchanged, impossible/disproportionate proof mechanisms corrected

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

### Feedback Cycles

- Cycle 8: REVISE — Roborev job 837; two Material successor-oracle findings accepted; declared caps remain 650 behavioral and 59 live LOC; ACs and production Go unchanged

## Stage Report: implementation (cycle 5)

- FAILED: Require a correlated public Codex `spawn_agent` call, child completion, and exact-child durable output; parent writes, narration, waits without a receiver, and synthetic commit handles must fail.
  Credentialed artifact `/tmp/spacedock-fo-gate-sequential-codex/codex-shared-scenarios/recorded-gate-lifecycle/codex-exec.jsonl` has no public spawn or child-completion event. After `item_25` dispatch-build, `item_26` only narrates `/root/recorded_gate_task_handoff`; `item_27` is `tool:"wait"` with `receiver_thread_ids:[]`; `item_28` is a parent-side state read. This cannot be correlated without synthesizing the missing public handle.
- SKIPPED: Parse Pi child assistant/tool-result events and require child-produced write evidence correlated with the completed async handle and resulting durable state; add parent-writes/no-op-child red controls.
  Stopped at the explicit Codex public-evidence boundary before expanding the zero-free 650-line behavioral budget; Pi remains a material accepted finding for the next correction after a valid Codex public surface is available.
- SKIPPED: Rerun focused/full/race and all live lanes, request final-tip Roborev, wait, and triage every finding.
  No code change was made because the required public Codex observation is absent; rerunning cannot manufacture the missing event, and weakening the oracle is forbidden.

### Summary

Roborev 837 is fully triaged and both Medium findings are Material. The implementation remains blocked in implementation on a missing public Codex spawn/completion correlation; no validation-readiness claim is made.

### Feedback Cycles

- Cycle 9: DESIGN RESET — Captain-approved proof-boundary ruling, 2026-07-24. Preserve the six-command fail-closed lifecycle, one successor dispatch, and its durable effect. Preserve deterministic command-deletion proof and one representative live journey per supported host. Remove mandatory public transport-event uniformity, exact-child forensic attribution where the supported host surface does not expose it, and every-route-by-every-host live execution. Host-native fixtures own transport details; live journeys own the observable lifecycle and durable outcome. Return to ideation for the smallest AC-1/AC-8 and test-surface delta before further implementation.
- Cycle 10: DESIGN RESET — Captain audit ruling, 2026-07-24. Preserve authorized record/consume-before-dispatch, durable successor state, and representative host journeys. Remove repeated four-subcommand help probes, the absolute-path workaround as a permanent FO ritual, runtime-specific “very next host event” wording, and zero-headroom absolute total-host ceilings. Spike whether every separate validation command adds safety before retaining exact six-command cardinality; keep resident/deferred component caps and define an attributable rebaseline policy. AC end values unchanged.
- Cycle 11: REVISE — independent ideation staff review; surface 1 design file with implementation not restarted vs expected implementation at or below 1,518 added LOC; AC end values unchanged. Preserve a state commit after every successful close, key the cached capability result to resolved executable identity, recompute task-attributable load at every tip/rebase, and prove consumed-commit ordering plus each omitted authority mutation through real traces.
- Cycle 12: REJECTED — validation resume matrix; surface 17 files/1,461 added LOC vs estimate 1,588 added LOC (92%); AC unchanged

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
