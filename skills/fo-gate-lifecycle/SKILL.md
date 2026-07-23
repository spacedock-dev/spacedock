---
name: fo-gate-lifecycle
description: Load at any engaged gate entry before action.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, and apply one recorded gate authorization

Load this skill in one host event before probe, write, validation, presentation, route, replay, or dispatch. No write authority; read `fo-write-core.md` before FO mutation.

**Boot projection.** Use only unresolved actionable `ready_gates` rows from `status --boot --identify --json`, fixed keys `id`, `slug`, `current`, `readiness`: `awaiting-captain` = selected current-stage open Briefing; `approved-awaiting-merge` = unblocked approve + advance/pending to terminal; `approved-awaiting-advance` = nonterminal target. Gate-stage/no selected attempt is omitted `validating`; malformed/stale selection, blocked/held, feedback, consumed/superseded/not-applicable are omitted. Opt-in human/JSON `gate-readiness` summarizes; `gate-*` retains diagnostics. Engage row `slug`, read entity, then `gate validate` for full Briefing/Resolution/application. Never infer readiness from status/stage.

**Capability preflight.** Probe the ONE startup launcher with `gate --help` and four subcommand help forms. Require `record` flags `--briefing`, `--result`, `--decision`, plus `validate`, `eligibility`, `consume`; version is insufficient. If missing, halt before mutation; prescribe refresh or `go build -o <temp>/spacedock ./cmd/spacedock`, set `SPACEDOCK_BIN`, retry. Never hand-edit `gates:`.

**Retain the package.** Before presentation assemble `ROOM/briefing.json`: concise capability/change, tests/evidence, reviewed snapshot, material/deferred/polish findings, one recommendation, decision ask; link raw entity/spec/reports. Keep reproducible payloads as URI + SHA, else freeze a room copy. Make `BRIEFING`, Results, associations absolute before the CLI.

Bind and prove the open attempt, in this order:

```text
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --briefing BRIEFING --workflow-dir WORKFLOW_DIR
${SPACEDOCK_BIN:-spacedock} gate validate ENTITY --workflow-dir WORKFLOW_DIR
```

Both calls must exit 0 and name the same gate, attempt, and Briefing; validation must report `state=open` and no decision. Otherwise halt before presentation and surface the command, exit, and missing artifact/field/step. In a split-root workflow, `«state.commit»(slug)` durably commits the folder-form room and index before presentation.

Only after open validation, invoke `«gate.ac-cross-check»`, make the evidence judgment, and invoke `«gate.assemble-verdict»`. The Captain sees the concise primary review, with the entity/spec/package as references—not a raw JSON/YAML dump or room listing.

**Record exactly who rendered the decision.** Use one semantic source:

```text
# Captain personally rendered the chat decision
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor person:captain [--reason REASON] --workflow-dir WORKFLOW_DIR

# FO rendered it under an explicit Captain conn
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor agent:first-officer --reason EVIDENCE_JUDGMENT --directive EXACT_QUOTED_CAPTAIN_GRANT --workflow-dir WORKFLOW_DIR

# Exact retained provider Result
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --result RESULT --association ASSOCIATION --actor AUTHORIZED_ACTOR [--adoption-note AUTHORIZER] --workflow-dir WORKFLOW_DIR
```

`revise` and `hold` require a reason (or the provider's included same-Briefing Annotation). Delegated FO approval always carries both its nonblank evidence reason and the exact quoted grant; never relabel it `person:captain`. A provider Result requires its complete retained association and authorized actor.

After closing, immediately prove closure:

```text
${SPACEDOCK_BIN:-spacedock} gate validate ENTITY --workflow-dir WORKFLOW_DIR
```

It must name the same attempt/Briefing, report `state=closed`, and reproduce the decision. Only then `«state.commit»(slug)`. Close/validation failure halts without feedback, advance, or dispatch.

**Route fail-closed.** For `approve`, run:

```text
${SPACEDOCK_BIN:-spacedock} gate eligibility ENTITY --workflow-dir WORKFLOW_DIR
${SPACEDOCK_BIN:-spacedock} gate consume ENTITY --workflow-dir WORKFLOW_DIR
```

Consume is authorized only when eligibility exits 0 with `condition=approved-pending eligible=true` and the expected immediate successor. Consume must exit 0 with `consumed=true`, atomically write that successor status and `application.state: consumed`, and then be committed through `«state.commit»(slug)`. Only after that durable consumed authorization may the deferred dispatch module run its ordinary reuse-or-fresh procedure for the newly current stage. `dispatch build` is still only an artifact: after building it, the very next host event is the runtime-bound `«worker.spawn»` tool call—no narration or wait may intervene. Require its returned live worker handle, and never enter wait or claim successor dispatch before that spawn is observed. Never use `status --set` to advance a gated stage.

- `revise`: eligibility, when read diagnostically, is ineligible `feedback/pending`; never consume it. Invoke `«feedback.route»` (including feedback-gate `REJECTED` and captain rejection at a `feedback-to` stage).
- `hold`: eligibility must be `not-applicable`/ineligible. Leave the entity at the gate and surface the reason; never consume, advance, or dispatch.
- approved but blocked, wrong-stage, unknown, or otherwise ineligible: halt with the reported condition and missing/current artifact or field; preserve status bytes.
- `stale`: first observe stale through read-only eligibility. Invoke consume only to materialize the landed failure: it exits nonzero, leaves status unchanged, and changes only pending → superseded. Commit, retain/bind a replacement Briefing, and re-present.
- already `consumed`: the authorization is spent. Follow the current status into ordinary dispatch/recovery; do not re-record or consume it. A diagnostic repeat consume must be nonzero and byte-clean.

**Resume before writing.** Run `gate validate` first. Same-Briefing `gate record --briefing` must correct stale `gates.current` without a new attempt; changed package updates the open attempt and is re-presented. Closed pending approval resumes eligibility/consume; revise/hold routes or stays held. Consumed resumes dispatch. Stale materializes supersession by consume failure, then binds a replacement. Surface nonzero command, exit, remedy; never repair frontmatter.
