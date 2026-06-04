---
title: Pi intercom runtime capability probe
status: implementation
source: captain (2026-06-04) — Codex idle-notification evidence pattern may generalize to prove Pi intercom contact_supervisor runtime capability
score: "0.31"
started: 2026-06-04T00:00:00Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-intercom-runtime-capability-probe
issue:
id: cq9kb7cdpp9y48tn8gwzmqzq
---

# Pi intercom runtime capability probe

Create a generalized runtime capability probe pattern, using Pi intercom supervisor talkback as the first target capability.

## Problem

`subagents-doctor` can report that the Pi intercom bridge is active. That is necessary setup evidence: it means the package path, bridge process, or extension discovery path is available enough for the host to advertise the bridge. It is not runtime capability evidence. Spacedock operations need a stronger property: a child subagent launched after bridge activation can actually reach its supervising session, send a non-blocking progress update, block on a decision request, receive the supervisor's reply, resume after that reply, and write durable post-reply evidence.

Without that distinction, a future implementation could pass `doctor --host pi` or a bridge-active preflight while still failing the workflow behavior Spacedock depends on: child-to-supervisor talkback. The failure modes are different and externally observable: the tool may be absent from the child prompt, progress messages may not arrive, `need_decision` may never wake the parent, the reply may not route back to the child, or the child may resume but not leave durable state proving it resumed after the reply.

The reference pattern is the Codex idle-notification probe in `docs/dev/codex-idle-notification-probe.md`, `docs/dev/_evidence/codex-idle-notification-probe/2026-06-03-dogfood.json`, and `skills/integration/codex_idle_notification_test.go`. That pattern separates runtime classifications, operator recipe, durable JSON evidence, and tests that reject over-claiming. Pi intercom should use the same style: not a transcript anecdote, but a reusable runtime capability probe with recipe + evidence JSON + tests.

Note: the requested Pi runtime docs (`docs/runtime-support.md` and `skills/ensign/references/pi-ensign-runtime.md`) are not present on this checkout's main tree at ideation time. I inspected the available Pi runtime support worktree copies at `.worktrees/spacedock-ensign-pi-runtime-support/docs/runtime-support.md` and `.worktrees/spacedock-ensign-pi-runtime-support/skills/ensign/references/pi-ensign-runtime.md`. Those notes already treat live runtime proof as distinct from setup and describe the proven Pi subagent live-smoke path; this task narrows the next proof to child-to-supervisor intercom talkback.

## Proposed approach

Add a reusable runtime capability probe pattern and instantiate it for Pi intercom supervisor talkback.

### Reusable probe pattern

A runtime capability probe should ship four linked artifacts:

1. **Probe recipe** under `docs/dev/<capability>-probe.md`.
   - Names the capability and the claim it proves.
   - Names setup preconditions that are necessary but insufficient.
   - Provides exact operator steps and exact child/worker prompt text.
   - Defines classification outcomes so failures and partial proofs are not mislabeled as success.
   - States where evidence JSON must be written.
2. **Evidence JSON** under `docs/dev/_evidence/<capability>-probe/*.json`.
   - Uses `schema_version: 1` and a stable `capability` / `host` / `classification` vocabulary.
   - Records setup observations separately from behavioral observations.
   - Records timestamps as RFC3339 UTC or `null` when unknown.
   - Records durable proof paths and marker content rather than relying on transcript prose alone.
3. **Focused integration tests**, analogous to `TestCodexIdleNotificationRecipeShape` and `TestCodexIdleNotificationEvidenceSchema`.
   - Recipe-shape tests fail when the exact smoke prompt, required sections, or interpretation rules disappear.
   - Evidence-schema tests parse every JSON record and reject missing required fields, invalid classifications, impossible `passed` records, or setup-only evidence labeled as capability success.
   - Runtime-contract/doc invariant tests enforce text claims at the correct level: setup wording may be checked statically, while talkback success must require live/manual evidence.
4. **Live/manual smoke path** for the behavior itself.
   - The live path is not required for every local `go test`, but it must be documented and record durable evidence when run.
   - The smoke should be cheap and bounded: one child, one progress update, one decision request, one parent reply, one durable marker write.

### Pi intercom first capability

The first new capability should be named `pi-intercom-supervisor-talkback` with evidence stored under `docs/dev/_evidence/pi-intercom-runtime-capability-probe/` or a similarly stable path chosen during implementation.

Recommended recipe sections:

- `## Capability under test`: child subagent can use Pi intercom supervisor talkback through `contact_supervisor`.
- `## Setup preflight`: run/record `subagents-doctor` bridge-active output and resolved `pi-intercom` / `pi-subagents` paths; explicitly state this preflight cannot classify the run as `passed`.
- `## Child prompt`: exact prompt for one Pi child subagent.
- `## Parent actions`: observe progress, answer the decision request with exactly `APPROVED`, then wait for the child result.
- `## Interpretation rules`: allowed classifications and what each means.
- `## Evidence record`: required JSON fields and durable marker requirements.

Exact smoke prompt to include in the recipe:

```text
You are a Pi intercom supervisor-talkback capability probe. Do not modify product/source files. Use contact_supervisor twice: first send reason progress_update with message "PI-INTERCOM-PROBE-PROGRESS"; then send reason need_decision with message "Reply exactly APPROVED to let the probe continue". After receiving the supervisor reply, create or update only the assigned probe marker file with the exact line "PI-INTERCOM-SMOKE-APPROVED" and return a concise completion message naming the marker file.
```

The assigned marker file should be in a temp workflow/evidence directory, not in product/source. The marker must be written only after the child receives the parent reply. A `passed` evidence record must include the marker path, marker content, and the parent-observed progress/decision events.

Suggested evidence schema fields:

```json
{
  "schema_version": 1,
  "run_id": "2026-06-04-example",
  "host": "pi",
  "capability": "pi-intercom-supervisor-talkback",
  "classification": "passed",
  "pi_cli_version": null,
  "pi_subagents_version": null,
  "pi_intercom_version": null,
  "subagents_doctor_bridge_active": true,
  "bridge_active_observed_at_utc": null,
  "bridge_active_output_excerpt": "...",
  "child_tool_available": true,
  "progress_update_observed": true,
  "progress_update_message": "PI-INTERCOM-PROBE-PROGRESS",
  "decision_request_observed": true,
  "decision_request_message": "Reply exactly APPROVED to let the probe continue",
  "supervisor_reply": "APPROVED",
  "child_resumed_after_reply": true,
  "marker_path": ".../pi-intercom-smoke-marker.txt",
  "marker_content": "PI-INTERCOM-SMOKE-APPROVED",
  "session_started_at_utc": null,
  "child_spawned_at_utc": null,
  "progress_observed_at_utc": null,
  "decision_observed_at_utc": null,
  "reply_sent_at_utc": null,
  "marker_written_at_utc": null,
  "interpretation": "Passed only because setup and the full progress/decision/resume/marker behavior were observed."
}
```

Allowed classifications should include at least:

- `passed`: bridge/setup evidence exists and the full progress update, decision request, supervisor reply, child resume, and durable marker were observed.
- `setup_only`: `subagents-doctor` or package discovery passed, but no child talkback behavior was exercised.
- `tool_unavailable`: child launched but did not receive usable `contact_supervisor` talkback tooling.
- `progress_only`: progress update was observed, but blocking decision/reply/resume was not proven.
- `decision_blocked`: decision request arrived but the child did not resume after a supervisor reply.
- `no_talkback_observed`: child ran but no progress or decision talkback was observed.
- `not_run`: operator recorded missing prerequisites or deliberately skipped live spend.

The implementation should reject any `passed` evidence when setup fields are missing, when behavioral booleans are false, when `supervisor_reply` is not exactly `APPROVED`, or when the marker content is absent/wrong. It should also reject any `setup_only` evidence that claims `child_resumed_after_reply: true`, because that mixes setup and behavior classifications.

## Spike / riskiest unknown plan

The riskiest unknown is not the JSON parsing or docs shape; those are already proven by the Codex idle-notification tests. The unverified mechanism is live Pi child-to-supervisor talkback after a bridge-active preflight: does a child launched through the Pi subagent substrate actually have `contact_supervisor`, can the parent observe both message types, and does `need_decision` resume after reply?

Ideation did not run a live intercom smoke. That is intentional: this stage is design, and a live run can spend model/runtime resources. The cheap evidence gathered without live spend is:

- Main checkout contains the Codex recipe, one evidence JSON record, and integration tests that validate recipe shape and evidence schema.
- Main checkout does not currently contain the requested Pi runtime docs; available worktree copies record a proven Pi subagent live-smoke mechanism and reinforce that runtime support requires live proof of durable state.
- No `pi-intercom`, `subagents-doctor`, or `contact_supervisor` documentation was found in the main checkout or the inspected Pi worktree notes, so the first implementation should start by documenting the recipe/schema rather than assuming an existing contract file covers talkback.

Implementation should pay the riskiest bill first after adding the static recipe/evidence tests: run a manual/live Pi intercom smoke only when prerequisites are available, record `not_run` or the exact failing classification if they are not, and never classify bridge-active alone as `passed`.

## Out of scope

- Reworking `pi-subagents` or `pi-intercom` internals.
- Replacing the existing Pi runtime/frontdoor live smokes.
- Adding a full Spacedock `install --host pi` path.
- Generalizing every host runtime in the first implementation; Codex notification evidence remains the reference, while Pi intercom is the first new capability probe.
- Treating transcript wording as proof of talkback success without durable JSON and marker evidence.

## Acceptance criteria

Each criterion is an end-state property of the implemented task and names proof outside this task body that can fail.

**AC-1 - Runtime capability probes have a documented reusable recipe and evidence contract.**
Verified by: a focused integration test that reads the new runtime capability probe recipe and evidence JSON files, parses them, and fails if required sections, required evidence fields, allowed classifications, or interpretation rules are missing.

**AC-2 - Pi intercom supervisor talkback is represented as a concrete capability probe.**
Verified by: a recipe-shape test that fails unless the Pi probe contains the exact child prompt requiring `contact_supervisor` `progress_update`, `contact_supervisor` `need_decision`, supervisor reply `APPROVED`, and durable marker content `PI-INTERCOM-SMOKE-APPROVED`.

**AC-3 - Setup evidence cannot be mistaken for talkback capability evidence.**
Verified by: evidence-schema tests that fail any `passed` record lacking both setup fields (`subagents_doctor_bridge_active`, resolved package/version/path data) and behavioral fields (`child_tool_available`, progress observed, decision observed, reply content, child resumed, marker path/content), and fail any `setup_only` record that asserts full behavioral success.

**AC-4 - Runtime/probe docs preserve the doctor-vs-capability distinction.**
Verified by: docs/probe invariant tests over the relevant recipe/runtime instruction files that require wording equivalent to “bridge active is necessary but insufficient” and reject wording that says or implies `subagents-doctor` bridge-active alone proves supervisor talkback.

**AC-5 - A live/manual Pi intercom smoke path can produce durable pass/fail evidence.**
Verified by: a documented live/manual command path that writes one JSON evidence record and a marker file; the evidence-schema test accepts `passed` only when the marker content is `PI-INTERCOM-SMOKE-APPROVED`, the parent observed progress, the parent replied `APPROVED` to the decision request, and the child resumed after that reply.

**AC-6 - Probe evidence remains reproducible and reviewable without live runtime access.**
Verified by: `go test ./skills/integration -run 'PiIntercom|RuntimeCapability' -count=1` or the implementation's equivalent focused package command passing against checked-in recipe/evidence fixtures, including a non-passed fixture such as `not_run` or `setup_only` when no live smoke has been recorded.

## Test plan

- **Focused integration tests (low cost, no live spend):** add tests analogous to `TestCodexIdleNotificationRecipeShape` and `TestCodexIdleNotificationEvidenceSchema` for the Pi intercom probe. They should parse the new recipe, require the exact smoke prompt and sections, parse every evidence JSON record, validate required fields by classification, validate RFC3339-or-null timestamps, and enforce the allowed classification enum.
- **Evidence schema negative coverage (low cost, no live spend):** add table tests or fixture mutations showing that `passed` fails without bridge-active setup, without progress, without decision, without exact `APPROVED` reply, without resume, or without the marker content; show that setup-only records cannot assert behavioral success.
- **Docs/probe invariant tests (low cost, no live spend):** add tests over the new recipe and any Pi runtime docs touched by implementation. These tests should enforce the doctor-vs-capability distinction as a text claim and forbid obvious over-claims that equate bridge-active with talkback proof.
- **Command gates (moderate cost, no live spend):** run `go test ./skills/integration -count=1` for the focused skill integration surface and `go test ./... -count=1` before handing off implementation. Run `gofmt -w ./cmd ./internal` only if Go files under those paths are edited.
- **Live/manual Pi intercom smoke (bounded live spend, run only with prerequisites):** run the documented Pi recipe with `pi-subagents`, `pi-intercom`, Pi auth, and `subagents-doctor` available. Launch one child, observe `PI-INTERCOM-PROBE-PROGRESS`, reply exactly `APPROVED` to the `need_decision` request, verify the marker file contains `PI-INTERCOM-SMOKE-APPROVED`, then write evidence JSON. If prerequisites are missing, write a `not_run` or `setup_only` evidence record rather than claiming `passed`.

## Stage Report: ideation

- DONE: Clarified why `subagents-doctor` bridge-active is necessary setup proof but insufficient for the child-to-supervisor talkback capability Spacedock needs.
- DONE: Proposed a reusable runtime capability probe pattern based on recipe, evidence JSON, integration tests, docs/probe invariants, and a live/manual smoke path.
- DONE: Specialized the pattern for Pi intercom with exact child prompt, evidence schema fields, allowed classifications, and pass/fail interpretation rules.
- DONE: Recorded the riskiest unknown and why no live intercom smoke was run during ideation.
- DONE: Rewrote acceptance criteria as end-state properties with external, fail-able proof.
- DONE: Added a test plan covering focused integration tests, evidence schema tests, docs/probe invariant tests, command gates, and a live/manual Pi intercom smoke path.
- SKIPPED: Live Pi intercom smoke; ideation scope is design and the live path can spend runtime/model resources.
- FAILED: Requested Pi runtime docs were not present in the main checkout; inspected available Pi runtime support worktree copies instead and recorded that limitation above.

### Summary

The task is ready for an implementation gate decision as a concrete design: implement a Codex-style runtime capability probe pattern, instantiate it for Pi intercom supervisor talkback, prove setup and behavior separately, and require durable evidence before any `passed` classification.

## Stage Report: implementation

- Product worktree: `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-intercom-runtime-capability-probe`
- Product branch: `spacedock-ensign/pi-intercom-runtime-capability-probe`
- Product commit: `81a11d066fa369b1e70e25c9e7d46cc2693824b7` (`Add pi intercom runtime capability probe`)
- Changed product files:
  - `docs/dev/pi-intercom-runtime-capability-probe.md`
  - `docs/dev/_evidence/pi-intercom-runtime-capability-probe/2026-06-04-not-run.json`
  - `skills/integration/pi_intercom_runtime_capability_test.go`
- Tests run from implementation worktree:
  - `go test ./skills/integration -run 'PiIntercom|RuntimeCapability' -count=1` — passed
  - `go test ./skills/integration -count=1` — passed
  - `go test ./... -count=1` — passed
- Formatting:
  - `gofmt -w skills/integration/pi_intercom_runtime_capability_test.go` — run for the new integration test
  - `gofmt -w ./cmd ./internal` — not needed; no `cmd/` or `internal/` Go files changed
- Implementation notes:
  - Added a Codex-style runtime capability probe recipe for `pi-intercom-supervisor-talkback`.
  - Added a non-passed `not_run` evidence fixture because no live Pi intercom smoke was run in this static implementation stage.
  - Added focused recipe/evidence/invariant tests that distinguish bridge-active setup evidence from behavioral talkback evidence and reject `passed` unless setup, progress, decision, exact `APPROVED` reply, resume, marker path, and `PI-INTERCOM-SMOKE-APPROVED` marker content are all present.
- Residual risks:
  - Live Pi intercom child-to-supervisor talkback remains unproven until an operator runs the documented manual smoke with safe prerequisites and records a non-`not_run` evidence record.
