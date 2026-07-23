---
name: fo-gate-lifecycle
description: Load at any engaged gate entry before action.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, and apply one recorded gate authorization

Complete this skill load in one host event before probes, writes, validation, presenter load, routing, replay, or dispatch. It grants no write authority; read `fo-write-core.md` separately immediately before the first FO mutation.

**Capability preflight.** Before creating a room or gate record, inspect the ONE launcher selected at startup. Run `${SPACEDOCK_BIN:-spacedock} gate --help` and the four subcommand help forms. The surface is capable only when it exposes `record` with `--briefing`, `--result`, and `--decision`, plus `validate`, `eligibility`, and `consume`. A compatible version alone is insufficient. On a missing form, halt before mutation and say to refresh the launcher or build the current source checkout with `go build -o <temp>/spacedock ./cmd/spacedock`, set `SPACEDOCK_BIN` to that executable, and retry the probe. Never hand-edit `gates:` as fallback.

**Retain the package.** Assemble `ROOM/briefing.json` (that exact basename) before presentation. Its primary review is concise and names capability/change, test and evidence, exact reviewed snapshot, material/deferred/polish findings, one FO recommendation, and the concrete decision ask. The entity, spec, reports, and other raw inputs are linked references. Keep an existing payload as URI + SHA when its exact bytes remain reproducible through the presentation resolver; otherwise freeze a room copy. Resolve `BRIEFING`, every Result, and every association to absolute paths before invoking the current CLI.

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

After a successful close, immediately prove closure:

```text
${SPACEDOCK_BIN:-spacedock} gate validate ENTITY --workflow-dir WORKFLOW_DIR
```

It must name the same attempt and Briefing, report `state=closed`, and reproduce the intended decision. Only then `«state.commit»(slug)`. Any close/validation failure halts without feedback, advancement, or dispatch.

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

**Resume before writing.** Run `gate validate` first. An open attempt with the same Briefing accepts idempotent bind + validate; a changed package updates the Briefing under that attempt and is re-presented. A closed pending approval resumes at eligibility/consume without re-recording. Closed revise/hold routes or stays held. Consumed status resumes dispatch without another consume. Stale materializes supersession through the one consume-failure path, then binds a replacement. Every nonzero command is FO-visible friction with command, exit, and actionable remedy; never repair frontmatter by hand.
