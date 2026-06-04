# Pi Intercom Runtime Capability Probe

Use this recipe when checking whether a Pi child subagent can perform supervisor
talkback through `contact_supervisor`. The probe separates setup preflight from
runtime behavior: `subagents-doctor` bridge-active output is necessary but
insufficient and does not prove supervisor talkback by itself.

## Capability under test

The concrete capability is `pi-intercom-supervisor-talkback` on host `pi`:
a child subagent launched after Pi intercom setup can send a non-blocking
`contact_supervisor` `progress_update`, send a blocking `contact_supervisor`
`need_decision`, receive the supervisor reply, resume after that reply, and write
a durable marker after resuming.

A live pass requires evidence for setup and behavior. Bridge active alone proves
only that setup discovery worked; it must never be interpreted as proof that a
child can contact its supervisor.

## Setup preflight

Before spending live runtime, record setup observations in the evidence JSON:

1. Run or capture `subagents-doctor --host pi` or the local equivalent.
2. Record whether the bridge was reported active in
   `subagents_doctor_bridge_active`.
3. Record available package/version/path observations in
   `pi_cli_version`, `pi_subagents_version`, and `pi_intercom_version`; use
   `null` when a value is unknown.
4. Preserve a short `bridge_active_output_excerpt` and
   `bridge_active_observed_at_utc` when known.

This setup preflight cannot classify a run as `passed`. If only setup was
checked, classify the record as `setup_only` or `not_run`.

## Child prompt

Assign the child a marker file in a temporary workflow/evidence directory, not
in product/source files. Use this exact child prompt:

```text
You are a Pi intercom supervisor-talkback capability probe. Do not modify product/source files. Use contact_supervisor twice: first send reason progress_update with message "PI-INTERCOM-PROBE-PROGRESS"; then send reason need_decision with message "Reply exactly APPROVED to let the probe continue". After receiving the supervisor reply, create or update only the assigned probe marker file with the exact line "PI-INTERCOM-SMOKE-APPROVED" and return a concise completion message naming the marker file.
```

The marker file must be written only after the child receives the supervisor
reply.

## Parent actions

1. Launch one Pi child with the exact child prompt and the assigned marker file.
2. Observe a progress update with exactly `PI-INTERCOM-PROBE-PROGRESS`.
3. Observe a decision request with exactly
   `Reply exactly APPROVED to let the probe continue`.
4. Reply to the child with exactly `APPROVED`.
5. Wait for the child to resume and complete.
6. Verify that the marker file contains exactly `PI-INTERCOM-SMOKE-APPROVED`.
7. Write one evidence JSON record under
   `docs/dev/_evidence/pi-intercom-runtime-capability-probe/`.

## Interpretation rules

Allowed classifications:

- `passed`: setup evidence exists and the progress update, decision request,
  supervisor reply `APPROVED`, child resume, and durable marker were observed.
- `setup_only`: bridge/package setup was observed, but no child talkback behavior
  was exercised.
- `tool_unavailable`: the child launched but did not receive usable
  `contact_supervisor` talkback tooling.
- `progress_only`: the progress update was observed, but decision/reply/resume
  behavior was not proven.
- `decision_blocked`: the decision request arrived, but the child did not resume
  after the supervisor replied.
- `no_talkback_observed`: the child ran, but no progress or decision talkback
  was observed.
- `not_run`: prerequisites were missing or the operator deliberately skipped
  live spend.

Do not claim that `subagents-doctor` bridge-active alone proves supervisor
talkback. A record classified as `passed` must include setup observations and all
behavioral observations. A `setup_only` record must not claim child resume or a
post-reply marker.

## Evidence record

Store run records as JSON under
`docs/dev/_evidence/pi-intercom-runtime-capability-probe/`. Use
`schema_version: 1`, RFC3339 UTC timestamps when known, and `null` when an exact
time was not captured.

Required fields:

```json
{
  "schema_version": 1,
  "run_id": "2026-06-04-example",
  "host": "pi",
  "capability": "pi-intercom-supervisor-talkback",
  "classification": "not_run",
  "pi_cli_version": null,
  "pi_subagents_version": null,
  "pi_intercom_version": null,
  "subagents_doctor_bridge_active": false,
  "bridge_active_observed_at_utc": null,
  "bridge_active_output_excerpt": "not run; live Pi prerequisites were not exercised",
  "child_tool_available": null,
  "progress_update_observed": false,
  "progress_update_message": "",
  "decision_request_observed": false,
  "decision_request_message": "",
  "supervisor_reply": "",
  "child_resumed_after_reply": false,
  "marker_path": "",
  "marker_content": "",
  "session_started_at_utc": null,
  "child_spawned_at_utc": null,
  "progress_observed_at_utc": null,
  "decision_observed_at_utc": null,
  "reply_sent_at_utc": null,
  "marker_written_at_utc": null,
  "interpretation": "Not passed because the live Pi intercom talkback smoke was not run."
}
```

For `passed`, the marker content must be exactly
`PI-INTERCOM-SMOKE-APPROVED`, the supervisor reply must be exactly `APPROVED`,
and progress/decision/resume booleans must all be true.

## Live/manual smoke path

Run this path only when Pi auth, `pi-subagents`, `pi-intercom`, and
`subagents-doctor` are already available and safe to use. Otherwise write a
`not_run` or `setup_only` evidence record rather than claiming success.

1. Create a temporary probe directory and choose a marker path such as
   `$TMPDIR/pi-intercom-runtime-capability-probe/pi-intercom-smoke-marker.txt`.
2. Run the setup preflight and record outputs.
3. Launch one child with the exact child prompt above and the marker path.
4. Complete the parent actions above.
5. Write or update the evidence JSON record with the observed classification.
