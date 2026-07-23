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
---

Make the normal First Officer gate path use the landed 3k/h1 commands so a presented decision is durably recorded, validated, checked for eligibility, and consumed before workflow advancement or dispatch.

## Problem

The recorder and application commands exist, but the shipped First Officer contract still assembles and presents gates only in prose. PR #557 (`fa240a76`) changed no First Officer skill files despite budgeting approximately ten lines for that integration; PR #560 (`f06cce04`) likewise added eligibility and consume without an FO caller. `skills/first-officer/references/first-officer-shared-core.md` names judgment, `present-gate`, feedback routing, advancement, and dispatch, but literally invokes none of `spacedock gate record`, `gate validate`, `gate eligibility`, or `gate consume`. This sprint therefore used the commands only under manual Captain directives, and an ordinary FO can still approve and advance outside the durable lifecycle.

The missing work is not a new decision model. It is the procedure the sprint FO already dogfooded: retain a canonical package, bind its exact `briefing.json`, validate the open attempt, make and present the irreducible evidence judgment, record the Captain or delegated decision, validate the closure, observe eligibility, consume the one-use application, then continue through the existing advance/dispatch loop. Before h1 landed, the same dogfood stopped after recording the Resolution and manually changed workflow state. The minimal correction is to replace only that manual state transition with `gate eligibility` followed by successful `gate consume`; `consume` atomically co-writes the next `status` and `application.state: consumed`. A separate `status --set` advance is thereafter a bypass and a contract violation.

The repository also exposed a readiness trap during this ideation: the retained repo-root `./spacedock` identifies as `0.26.0+dev` and passes the skill version floor, yet `./spacedock gate record --help` (and the other gate subcommands) falls through to top-level usage. Version compatibility is therefore insufficient evidence that the selected executable has the four required capabilities. The FO must prove the command surface before opening a gate package and use a fresh checkout build or refreshed installation when it is stale.

## Demonstrated baseline

The design adopts the landed behavior and historical dogfood; it does not reopen whether durable gates should exist.

1. The 3k validation room retained a canonical `review/validation/briefing-v1/briefing.json` plus a concise gate review and frozen references. State commit `2c616b7e` bound it; `77590ebd` recorded an `agent:first-officer` approval carrying the Captain decision as adoption provenance. The resulting record stayed at `status: validation`: 3k intentionally recorded the decision but did not apply it.
2. The remaining 3k closeout used manual lifecycle mutations (`mod-block`, PR state, and later terminal archive) because h1 was not yet available. Those commits are historical evidence for the exact integration gap, not a procedure to preserve.
3. h1 landed `gate eligibility` and `gate consume`. Its CLI fixture proves an approve closure yields `advance/pending`, eligibility prints `condition=approved-pending eligible=true`, and consume atomically changes `status` plus `application.state: consumed`; a second consume exits nonzero with `condition=consumed consumed=false` and byte-identical state.
4. The recorder's retained package boundary is settled: the manifest basename is exactly `briefing.json`; its canonical bytes independently define the complete artifact inventory. Existing artifact payloads remain URI + SHA references when the presentation resolver can reproduce those exact bytes. Mutable, cross-root, or otherwise unreproducible reviewed material is frozen as a room copy before publication. The recorder verifies but never copies artifact payloads. For folder-form entities, landed vn/PR #558 makes one `spacedock state commit <slug>` include the index and every non-ignored room artifact without sweeping sibling dirt.
5. Direct and delegated decisions remain distinct. A Captain who personally renders the chat decision is `person:captain`. An FO acting under delegated conn is `agent:first-officer`, supplies its evidence-bearing `--reason`, and stores the exact quoted grant in `--directive`. An exact provider Result uses `--result` plus its retained association and authorized actor; advisory adoption also names the authorizer.

No spike is needed for the command semantics: 3k/h1's landed CLI and fixture tests prove binding, validation, closure, eligibility, atomic one-use consumption, staleness, and retries, while the retained 3k state commits prove the real package/decision inputs. The unproven mechanism is whether the shipped FO actually performs the full sequence. That is this task's value claim, so implementation starts with the live, fixture-backed FO replay below rather than another command-model spike.

## Operational lifecycle

The procedure belongs at the existing gated-stage branch in `first-officer-shared-core.md`, around `«gate.assemble-verdict»`; it does not move or duplicate `«gate.ac-cross-check»`, FO evidence judgment, or `present-gate` rendering.

### 0. Select a capable executable

Before writing a room or gate record, resolve `${SPACEDOCK_BIN:-spacedock}` and inspect its gate help. Success means the output names all four landed operations and their semantic forms: `record` with `--briefing`, `--result`, and `--decision`; `validate`; `eligibility`; and `consume`. A compatible version string alone does not satisfy this check. If the surface is absent, halt before mutation and name the remedy: refresh the launcher, or from a source checkout build a fresh temporary binary with `go build -o <temp>/spacedock ./cmd/spacedock`, set `SPACEDOCK_BIN` to that exact path, and rerun the capability check. Resolve every retained Briefing, Result, and association argument to an absolute path before invoking the current CLI; a repo-root-relative Briefing parent is currently combined with the absolute entity root and fails before mutation. The skill records either problem as CLI friction; it does not silently fall back to hand-edited frontmatter.

### 1. Retain and bind the presented package

The FO assembles the same judgment inputs it already owns, but retains their exact presentation package first:

- `ROOM/briefing.json`, with the required basename and the decision question;
- a concise primary gate review stating capability/change, test and evidence, exact reviewed snapshot, material/deferred/polish findings, one FO recommendation, and the concrete decision ask;
- the entity, spec, reports, and other raw material as reference artifacts, not as substitutes for the primary review;
- each artifact identified by URI and exact SHA. Reuse an existing immutable file only when the presentation resolver can reproduce its bytes; otherwise freeze a room copy. Do not duplicate reproducible payloads merely to satisfy the recorder.

After the complete package is retained, invoke:

```text
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --briefing ROOM/briefing.json --workflow-dir WORKFLOW_DIR
${SPACEDOCK_BIN:-spacedock} gate validate ENTITY --workflow-dir WORKFLOW_DIR
```

The record call must exit 0 and identify the selected gate/attempt/Briefing. The immediate validation must exit 0, name the same identities, and report `state=open` with no decision. Anything else halts presentation. In a split-root workflow, `${SPACEDOCK_BIN:-spacedock} state commit SLUG` then durably commits both the folder room and index mutation before the Captain sees the gate.

### 2. Judge and present

Only after the open binding validates does the FO perform `«gate.ac-cross-check»`, make its evidence judgment, and invoke `present-gate`. The presenter remains an overridable channel; this task does not redesign it. Regardless of channel, the Captain sees the concise primary review rather than a raw entity dump, spec, room listing, or `briefing.json`. Those remain linked references.

### 3. Record and validate the decision

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

The close call must exit 0. Immediately invoke `gate validate` again; it must name the same attempt and Briefing, report `state=closed`, and reproduce the intended decision. Commit the entity mutation with `state commit` only after closure validation succeeds. A failed close or failed validation halts without feedback, advancement, or dispatch.

### 4. Route the closed decision fail-closed

The Resolution is durable before any route:

| Closed decision / condition | Required FO route |
| --- | --- |
| `approve` | Run `gate eligibility`. Only exit 0 plus `condition=approved-pending eligible=true` and the expected immediate successor authorizes `gate consume`. Consume must exit 0 with `consumed=true`, atomically advance `status`, and mark the application consumed. Commit that mutation, then run the existing reuse-or-fresh dispatch procedure for the newly current stage before yielding. No separate `status --set` advance is allowed. |
| `revise` | Optionally read eligibility for diagnostics; it must be ineligible and identify `feedback/pending`. Never consume it as an advance. Invoke the existing `«feedback.route»` procedure; this task does not implement advisory correction rounds. |
| `hold` | Read eligibility; it must report `not-applicable`/ineligible. Leave the entity at the gate, do not call consume, advance, or dispatch, and surface the reason. |
| approved but blocked, held, unknown, wrong-stage, or otherwise ineligible | Halt advancement and dispatch. Preserve status bytes. Name the exact reported condition and missing/current artifact or field. |
| `stale` | A read-only eligibility check first reports stale. Invoke `gate consume` only as the landed staleness materialization path: it must exit nonzero, leave status unchanged, and change the pending application to `superseded`. Commit that supersession, retain/bind a new Briefing, and re-present; never expose a separate rebind/supersede ceremony to the Captain. |
| already `consumed` | Treat the approval as spent. Do not record another decision or consume again; follow current `status` into the ordinary dispatch/recovery loop. A diagnostic repeat consume is nonzero and byte-clean. |

### 5. Resume and retry without minting duplicates

On resume, run `gate validate` before replaying any write:

- open attempt + same retained Briefing: repeating `record --briefing` is an idempotent bind; validate and continue presentation;
- open attempt + changed package: `record --briefing` selects the changed binding under the same attempt; the FO speaks only of updating the presented Briefing, not recorder mechanics;
- closed + pending approval: do not re-record the decision; run eligibility and, if authorized, consume;
- closed + revise/hold: route or remain held; never call advance-consume;
- consumed + already advanced status: dispatch/recover the current stage without consuming again;
- stale: materialize supersession through the one landed consume failure path, bind the replacement Briefing, and re-present.

Every nonzero command becomes explicit FO-visible friction with the command, exit, and actionable missing artifact/field/step. The FO never repairs `gates:` by hand or substitutes `status --set`.

## Friction inventory and disposition

| Observed sprint friction | Disposition and coverage |
| --- | --- |
| The FO skill contains no literal `record`/`validate`/`eligibility`/`consume` calls. | In scope: lifecycle procedure at the gated-stage branch. AC-1/AC-2 live and mutation proof. |
| Briefing, association, and Result JSON were manually crafted and error-prone. | In scope only as exact input guidance and retained real fixtures. The recorder continues to derive lifecycle/ids; no new JSON builder. AC-5 checks actionable failures. Better package authoring remains presentation-provider work. |
| It was unclear which entity/spec/reference artifacts to package and whether to use URI+SHA or copied bytes. | Settled rule above: reproducible existing bytes stay URI+SHA; mutable/unresolvable snapshots are frozen in-room; the recorder never copies. AC-4 replays the three-artifact 3k package and hashes every declared revision. |
| `briefing.json` linkage/discoverability was unclear; an accepted alternate basename later dead-ended Result resolution. | In scope guidance requires the canonical basename and a room whose relative references resolve. AC-4/AC-5 include the landed alternate-basename no-mutation control. |
| Bind → open validate → decision → closure validate was not presented as one legible lifecycle. | In scope: named phases above and exact trace oracle. AC-1/AC-2. |
| Raw room/entity/spec files were mistaken for validation-gate evidence. | In scope: primary review must state capability, test/evidence, snapshot, findings, recommendation, and ask; raw artifacts are references. AC-6 grades the rendered gate. |
| Delegated Captain conn left actor/reason/directive provenance ambiguous. | In scope: direct `person:captain` versus delegated `agent:first-officer` plus evidence reason and exact quoted directive. AC-3. |
| Before h1, the FO manually changed status after recording a Resolution. | In scope: replace only that step with eligibility + consume and prohibit `status --set` bypass. AC-1/AC-2 measure zero unconsumed advances. |
| Retries, stale input, and duplicate application behavior lacked an agent path. | In scope resume matrix maps each landed condition to idempotent bind, no re-close, consume once, stale supersession, or dispatch recovery. AC-2/AC-7. No retry daemon or new ids. |
| Failures did not reliably tell the FO which artifact/step was missing. | In scope as halt-and-surface contract plus output assertions. AC-5; parser/record failures are byte-clean, validate/eligibility are read-only, held/blocked/repeat-consume refusals are byte-clean, and stale consume's sole permitted mutation is explicit supersession. |
| Fixtures placed as workflows polluted repo-wide workflow discovery. | In scope test layout: all fixtures live below existing ignored/test-owned `skills/integration/testdata/` or temp dirs, never as discoverable repo workflows. AC-4/AC-8. |
| A real runtime replay, including Pi, was absent. | In scope: one host-neutral shared scenario, codified fixture, Claude/Codex adapters, and Pi live-capable coverage. AC-1/AC-8; host-neutral core changes require the relevant live lanes. |
| Extra folder-form room artifacts were omitted by `state commit` and required manual exact-path commits. | Resolved dependency, not reimplemented: vn/PR #558 makes folder-form `state commit` the durable unit. AC-4 exercises it; no new artifact registry or state verb. |
| Historical gates held canonical decisions, but not every one came through the literal CLI. | In scope: use the real retained 3k validation package and delegated decision as the replay baseline, then require the actual command trace. AC-1 closes the history/procedure gap rather than treating hand-authored YAML as proof. |
| Gate presentation became long and context-dependent. | In scope: one concise primary review, entity/spec as linked references. AC-6. Presentation UI/channel design remains xb-owned. |
| `rebind`/`supersede` vocabulary exposed recorder mechanics to humans. | In scope vocabulary rule: the FO operates a Briefing/gate with existing commands; it speaks of updated, stale, or spent decisions. No new command/subcommand. AC-6. |
| The FO silently worked around command/CLI friction during the sprint. | In scope: every command failure is surfaced and recorded; no manual `gates` edit. AC-5. Product CLI UX defects discovered by the replay are reported as findings, not patched inside skill work. |
| Version-compatible repo-root `./spacedock` lacked the gate subcommands. | In scope capability readiness probe plus fresh-build/refresh route before mutation. AC-5. Version/doctor-only readiness is explicitly insufficient. |
| Fresh `/tmp/spacedock-current gate record first-officer-gate-command-lifecycle --briefing docs/dev/.spacedock-state/first-officer-gate-command-lifecycle/review/ideation/briefing-1/briefing.json --workflow-dir docs/dev` failed before mutation with `resolve briefing room: Rel: can't make docs/dev/.spacedock-state/first-officer-gate-command-lifecycle/review/ideation/briefing-1 relative to /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/first-officer-gate-command-lifecycle`. | In-scope FO correction: pass absolute retained-input paths and surface that exact remedy. AC-5 replays relative failure byte-clean, then absolute-path success. Product normalization/actionable-error behavior is a separately named CLI follow-up; no recorder schema or lifecycle change enters this task. |
| The presentation helper could report failure before a late provider Result arrived, and earlier cleanup destroyed result bytes. | Boundary: xb/presentation transport owns exact Result and association retention; this procedure accepts only retained bytes. AC-4's provider arm uses retained inputs. No polling or transport enters this task. |
| Provider and canonical Briefing ids differed even when primary artifact bytes matched. | Boundary: landed recorder verifies complete association and normalizes identity only after exact bytes/authority. The FO supplies Result + association; it does not normalize by hand. AC-4. |
| `status --set` reserialized hand-authored `gates`, breaking anchors. | Prohibited workaround: binary owns all gate writes; `status --set` never substitutes for recorder or consume. AC-2/AC-5. |

## Acceptance criteria

**AC-1 (VALUE) — A normal FO-approved gated stage cannot advance or dispatch without one validated, consumed authorization.** In the real 3k validation-package replay, the FO emits the exact ordered lifecycle trace (briefing record, open validate, delegated decision record, closed validate, eligibility, consume), on-disk state ends at the expected successor with the same attempt's `application.state: consumed`, and the successor receives one observed dispatch. The measured count of advances/dispatches with no matching consumed application is **0**, versus the current contract's possible manual path of 1. *Verified by:* a fixture-backed live FO journey graded on command events, entity before/after bytes, and dispatch observation—not transcript phrasing.

**AC-2 — Every required lifecycle call is load-bearing and advancement fails closed.** Deleting any one of briefing record, open validation, decision record, closure validation, eligibility, or consume from the expected trace makes the scenario grader refuse the successor dispatch; the entity either remains at the gate or the trace is incomplete and cannot grade PASS. Revise, hold, and ineligible controls produce zero advance/dispatch effects; revise invokes feedback routing and neither revise nor hold is consumed as advance. *Verified by:* table-driven skipped-step trace mutants plus durable-state controls, with at least one mutated shipped-skill live replay; the test must demonstrate its mutant goes red.

**AC-3 — Direct and delegated authority remain durably distinguishable.** A direct Captain decision records `by: person:captain`; a delegated FO decision records `by: agent:first-officer`, a nonblank evidence reason, and the exact quoted Captain conn as directive/adoption provenance. Missing delegated fields fail before closure and leave the entity byte-identical; a procedure mutant that records the delegated live decision as `person:captain` fails the provenance grade even though the generic recorder cannot know who spoke in chat. *Verified by:* public CLI fixtures, the delegated live baseline, and the actor-swap scenario mutant.

**AC-4 — The FO package is complete, reproducible, and durable without unnecessary copying.** The replay uses the retained 3k three-artifact validation Briefing: every URI resolves through the declared resolver and every SHA matches; reproducible existing artifacts are not duplicated, while the frozen reviewed snapshot stays in the room. One folder-form `state commit` includes the index and all new room files and excludes dirty sibling paths. Fixtures live only in testdata/temp locations and do not alter repo workflow discovery. *Verified by:* exact digest/path assertions, a real-Git state-commit assertion, and before/after workflow-discovery candidate equality.

**AC-5 — Readiness and failures are actionable and preserve the landed write boundaries.** A version-compatible but capability-stale executable is rejected before any room/entity mutation and the output names refresh/fresh-build remediation. Missing/noncanonical Briefing, invalid association/actor/provenance, relative-path room resolution, and close validation failures name the missing artifact/step or an absolute-path remedy, return nonzero, preserve entity bytes, and leave no lock residue where the landed commands promise it. The exact replay command with a repo-root-relative Briefing path must fail byte-clean; the same retained bytes passed by absolute path must bind successfully. Validate/eligibility reads are byte-clean; hold/blocked/repeat-consume refusals are byte-clean; stale consume changes only pending → superseded and never status. *Verified by:* CLI failure matrix with whole-file/tree hashes, relative/absolute path pair, and a stale-binary shim whose version passes but gate help lacks the surface.

**AC-6 — The Captain receives a concise decision review, not recorder mechanics or raw artifacts.** The live output contains capability/change, test/evidence, reviewed snapshot, findings, one recommendation, and a concrete decision ask; entity/spec/Briefing are linked references. It does not lead with raw JSON/YAML, room listings, or `rebind`/`supersede` ceremony. *Verified by:* the live scenario's structured gate-review grader and retained output artifact; a raw-file-dump mutant fails the grade.

**AC-7 — Resume is idempotent and one-use.** Open same-Briefing retry creates no duplicate attempt; a closed pending approval is not re-recorded; consumed state never consumes twice; stale state becomes superseded without advancing and requires a replacement Briefing. Across three resume passes, the count of successful consumes and resulting stage transitions is exactly 1. *Verified by:* fresh-process fixture passes and byte comparisons after each retry.

**AC-8 — The procedure is runtime-portable and fixture-backed.** The recorded-gate-lifecycle journey is a host-neutral shared scenario with Claude and Codex runners and explicit Pi live-capable coverage, alongside the deterministic CLI/trace fixture. A host adapter may differ only in launch/transport; the command order and durable oracle are identical. *Verified by:* shared-scenario parity/meta tests, the Pi live replay (or required CI artifact), focused codified fixture, and all live lanes required by the host-neutral skill diff.

## Minimum replay and test plan

The test package copies the real retained 3k validation Briefing and its declared artifacts into a temp folder-form workflow whose gated `validation` stage has a supported nonterminal successor. It preserves the package bytes and uses the historical delegated-approval shape: `--actor agent:first-officer`, `--reason "Captain directive: approved after reviewing the presented 3k validation gate."`, and `--directive "you have the conn toward the sprint goal; authorized to approve gates, PR, relevant CI lanes, and merge; use your judgement."`. A command-logging wrapper delegates each gate operation to the freshly built real binary, records exit/stdout/stderr, and permits the dispatch stub only after the complete successful trace.

1. **Baseline live journey (AC-1/AC-3/AC-6/AC-8; high):** launch the real FO, engage the held validation gate, validate its concise evidence review, supply/retain the delegated approval, and assert six ordered gate events, one consume, atomic advanced+consumed state, one state commit per mutation boundary, and one successor dispatch.
2. **Skipped-step mutants (AC-2; medium):** table-delete each of the six required events from the trace accepted by the dispatch stub; every arm must deny dispatch. Run at least one copied-skill live mutant (remove one recorder call at the actual integration point) and prove the live grade reds.
3. **Revise/hold/ineligible controls (AC-2/AC-5; medium):** use the same Briefing with each decision/condition. Revise records/validates and routes feedback with zero advance-consume; hold and blocked approval remain at the gate; stale materializes only supersession.
4. **Provenance and package matrix (AC-3/AC-4/AC-5; medium):** direct versus delegated versus exact Result inputs; complete versus truncated association; canonical versus alternate basename; repo-root-relative versus absolute retained-input paths; reproducible URI+SHA versus frozen copy. Compare bytes and lock residue on every refusal, and prove the absolute-path retry binds the same bytes.
5. **Resume passes (AC-7; low):** rerun fresh processes over open, pending, stale, and consumed snapshots; exactly one stage transition and consume across three passes.
6. **Capability-stale launcher (AC-5; low):** shim a passing `--version` and top-level-only gate help; assert pre-mutation halt and fresh-build remedy. Then point `SPACEDOCK_BIN` at the newly built binary and run the baseline green.
7. **Runtime lanes (AC-8; high):** add the host-neutral scenario once, runner bindings for Claude/Codex, and a Pi live replay using isolated auth/home and explicit local skills. Required CI artifacts retain command logs, gate review output, before/after entity, and dispatch evidence.
8. **Repository gates:** `gofmt -w ./cmd ./internal`, focused scenario/skill tests, `go test ./...`, `go test ./... -race`, live-tag compilation, and the required live lanes.

The primary independent oracle is durable state plus the command/dispatch event log. Static presence of command names in skill prose is inspection evidence only and cannot satisfy any behavioral AC.

## Expected surface and tolerance

- **FO contract:** `skills/first-officer/references/first-officer-shared-core.md`, approximately 35–60 lines at the existing gate branch. No host adapter wording, new skill, or new presentation template.
- **Codified/live proof:** a common recorded-gate fixture/scenario under `internal/ensigncycle/`, the shared scenario table/meta/docs lock, narrow Claude/Codex runner bindings, Pi live-capable runner/coverage, and retained fixtures under `skills/integration/testdata/`; approximately 500–850 Go test LOC plus 150–300 fixture lines across 9–14 files. No production Go LOC.
- **Docs:** `docs/site/concepts/gates-and-decisions.md` and a capability-readiness note in `docs/site/reference/command-reference.md`, approximately 15–30 lines.
- **Tolerance:** 2×. Reconfirm if implementation needs production Go, a schema change, another gate/state command, a package/association generator, presenter UI, provider polling/transport, a new retry/idempotency engine, or more than 120 skill-contract lines. Those are new products, not integration variance.

## Documentation change proposal

After “The three calls” in `docs/site/concepts/gates-and-decisions.md`:

```diff
+Before the First Officer shows a gate, it binds the exact retained Briefing and validates
+the open attempt. After the decision, it records and validates the Resolution before any
+effect. An approval advances only when `gate eligibility` reports the expected current
+successor and `gate consume` spends that authorization; consume writes the new stage and
+the consumed mark together. Revise routes feedback instead, and hold stays at the gate.
+The review itself remains concise: capability, evidence, reviewed snapshot, findings,
+recommendation, and decision ask. The entity, spec, and package are linked references.
```

After the setup compatibility paragraph in `docs/site/reference/command-reference.md`:

```diff
+For source-checkout or retained development launchers, version compatibility is not a
+command-capability check. Before gate work, `spacedock gate --help` must list `record`,
+`validate`, `eligibility`, and `consume` with the semantic record forms. If it does not,
+refresh the installed launcher or build the current checkout and select that binary with
+`SPACEDOCK_BIN`; do not hand-edit gate frontmatter as a fallback.
```

## Dependencies, collision, and non-goals

- 3k/PR #557 and h1/PR #560 are landed, binding inputs. vn/PR #558 supplies folder-room durability. xb owns presentation-channel transport/retention and can remain a reference dependency; this task keeps `present-gate` judgment/rendering intact.
- The concurrent advisory-round generalization explicitly owns `docs/dev/README.md`, `skills/feedback-rejection-flow/SKILL.md`, recorder/schema surfaces, and a round-only caller. It declares no `first-officer-shared-core.md` change, so the intended FO contract surfaces do not collide. Both tasks may edit `docs/site/reference/command-reference.md`; serialize that doc edit (land/rebase one before the other). If either implementation expands into the other's skill file or ordinary gate procedure, stop and serialize rather than resolving two concurrent meanings.
- No schema or recorder/application logic; no gate judgment duplication; no presenter redesign; no advisory correction-round implementation; no blocker evaluator or hold authoring; no arbitrary artifact registry; no provider launch/polling; no new gate verb; no hand-authored `gates:` repair; no stable-release work.

## Minimum value demonstration seed

In one fixture-backed live First Officer journey, package an exact validation Briefing, record and validate it, record an evidence-bearing delegated approval, observe eligibility, consume the application, and only then advance and dispatch. Two controls must fail closed without status mutation: a revise decision routes through feedback rather than consume, and a hold or ineligible approval does not advance. Removing any recorder command from the FO procedure must make the journey fail.

## Boundary

This task owns the First Officer invocation contract and its behavioral proof. It reuses the landed 3k recorder and h1 eligibility/consume commands; it does not add a recorder schema, duplicate gate judgment, change presentation UI, or implement advisory correction rounds.

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

## Implementation intended-change declaration

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
  The code worktree intentionally remains an uncommitted checkpoint because its resident-core layout violates a binding design-reset trigger. Actual checkpoint delta: FO contract +58/-3 lines; test/live +1,157 lines; docs +7 lines; production Go 0. Do not advance to validation until the Captain/FO authorizes the deferred gate-core topology, a justified prompt-budget reset, or a materially narrower contract.

### Summary

Implemented and hardened the six-event integration and its deterministic/live oracles, then stopped at the binding design reset exposed by the final prompt-load gates. The behavioral work is preserved in the assigned worktree, all advisory findings are triaged, and no acceptance criterion or spawn evidence was weakened to manufacture a pass.
