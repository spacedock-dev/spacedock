---
name: fo-gate-lifecycle
description: Handle an engaged gate entry.
user-invocable: false
---

# First Officer Gate Lifecycle

## «gate.lifecycle»(slug, stage): bind, decide, apply one recorded authorization

Load before any engaged gate action. It grants no writes; read `fo-write-core.md` before FO mutation.

**Boot projection.** Use unresolved actionable `ready_gates` from `status --boot --identify --json`. Engage row `slug`, read its entity; never infer readiness from stage alone. `awaiting-captain` means an open current-stage Briefing; `approved-awaiting-merge`/`approved-awaiting-advance` are unblocked. For gated current `status`, any selected or retained gate record used to resume/present must have `stage` equal to that status. A prior-stage `gates.current` is history, not reusable authority. If no selected attempt (omitted `validating`) or it mismatches, run `gate prepare` for the current status and present only its emitted binding.

**Capability preflight.** Resolve `${SPACEDOCK_BIN:-spacedock}` and per lifecycle run exactly one fresh `gate --help`. Require `prepare`, `record`, `validate`, `eligibility`, `consume`; prepare flags `--question`, `--artifact`, `--summary`, `--reference`, `--workflow-dir`; and record flags `--briefing`, `--room`, `--decision`, `--actor`, `--reason`; reject retired `--directive` exposure. On absence or exposure, halt before mutation: commit no selected source, prepare no room, change no state; prescribe refresh or a fresh build via `SPACEDOCK_BIN`. Never hand-edit `gates:`.

**Prepare and bind.** Select one Markdown gate-review Artifact and References, author its exact concise summary, then after preflight commit new selections in their main/state histories. Supply only judgment and paths; never author JSON, ids, digests, Git-root locators, or room coordinates. Relative paths use launch cwd.

```text
${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary SUMMARY [--reference FILE ...] --workflow-dir WORKFLOW_DIR
```

Require the emitted `room`, `briefing`, `digest`, and `state=open` lines. The emitted clean absolute room is sole authority; never reconstruct it. Preparation binds two recorder-ready files without source copies. `«state.commit»(slug)` commits the folder entity or flat Markdown-plus-companion room unit. Before presentation, load `Skill(skill="spacedock:present-gate")`; an override replaces only chat display. Invoke `«gate.ac-cross-check»`, judge evidence, then `«gate.assemble-verdict»`. Chat presentation completes only after one root review names entity/stage, exact bound Briefing id/digest, recommendation, and decision ask. It follows the bind commit, must precede decision record; delegated conn does not waive it.

**Record and durably close.** Use exactly one semantic source:

```text
# Captain's chat decision
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor person:captain [--reason REASON] --workflow-dir WORKFLOW_DIR

# FO decision under explicit Captain conn
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve|revise|hold --actor agent:first-officer --reason EVIDENCE_JUDGMENT --workflow-dir WORKFLOW_DIR

# Exact retained provider Result
${SPACEDOCK_BIN:-spacedock} gate record ENTITY --room ROOM --workflow-dir WORKFLOW_DIR
```

No explicit Captain grant in the active conversation: bind/present; leave the gate open. A grant including one issued later in that conversation permits delegation. Record an FO-rendered decision as `agent:first-officer` with a nonblank reason, never `person:captain`; reserve it for the Captain's own decision. Recorder authenticates/retains no grant. `revise`/`hold` need reasons; provider mappings must be complete.

Map Captain calls before recording: `approve` maps to `approve` with an accepts-direction evidence reason; `redo with feedback` maps to `revise` with an accepts-direction reason; `reject` with `feedback-to` maps to `revise` with a rejects-direction reason; `reject` without `feedback-to` maps to `hold` with a pause reason; `hold` maps to `hold` with a pause reason; `not yet` maps to `hold` with a pause reason naming what remains. Routed redo/reject reasons include concrete asks and invoke `«feedback.route»` after the close commit; hold decisions commit and stop at the gate.

Require exit 0, the bound attempt/Briefing, `state=closed`, and the decision; record already validates the Resolution/application before atomic write. After every successful close, `«state.commit»(slug)` must commit that exact Resolution before approve, revise, hold, or any consume attempt. Close/commit failure halts.

**Route fail-closed.** For approve, run:

```text
${SPACEDOCK_BIN:-spacedock} gate consume ENTITY --workflow-dir WORKFLOW_DIR
```

Consume itself rechecks currency, successor, blockers, and one-use state under lock. Require exit 0 with `approved-pending`, `eligible=true`, `consumed=true`, and the expected successor. It atomically writes successor status plus consumed state; commit that descendant through `«state.commit»(slug)`. A nonterminal target then enters ordinary reuse-or-fresh dispatch; a terminal target enters the existing merge guard/hook and has no successor dispatch. Never use `status --set` to advance a gate.

- `revise`: after its close commit, never consume; invoke `«feedback.route»`.
- `hold`: after its close commit, remain at the gate and surface the reason.
- blocked/wrong-stage/unknown/ineligible approval: its close is already durable; halt and preserve status bytes.
- `stale`: consume exits nonzero, leaves status unchanged, changes only pending → superseded; commit it, bind a replacement Briefing, re-present.
- already `consumed`: the authorization is spent. A nonterminal current status resumes ordinary dispatch/recovery; a terminal current status resumes the existing merge ceremony. Do not re-record, consume, or dispatch a terminal successor. A diagnostic repeat consume must be nonzero and byte-clean.

**Resume.** Use boot/entity state and prior result; `gate validate`/`gate eligibility` are optional diagnostics, never positive-path requirements. Same-Briefing bind is idempotent. Changed open request-less Briefings update/re-present; request-backed ones are frozen—surface rebind refusal and stop. Require the exact Resolution commit before routing closed state. Pending approval → consume; revise/hold → route/stop; consumed → dispatch only if nonterminal, else merge; stale → supersede then replace. Surface nonzero command, exit, remedy; never repair frontmatter.
