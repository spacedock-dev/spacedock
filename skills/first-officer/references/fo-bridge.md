# FO Bridge Egress — Permission Blocks and Alerts

Deferred reference for the FO's Bridge egress signals beyond the boot heartbeat. Loaded when the FO hits a permission block; not read at boot. The heartbeat + inbox-drain half lives in the boot-resident core (Startup step 2b) and the fleet routing in `references/fo-fleet.md`; the richer egress surface is documented in `docs/dev/bridge-egress-contract.md`.

## Permission Blocks and Bridge Alerts

When a host sandbox or permission boundary blocks a workflow action that would otherwise be valid to attempt, surface it as a top-level Bridge alert before parking the loop. Use the hidden helper from the repo root:

```
${SPACEDOCK_BIN:-spacedock} bridge alert permission --host <host> --workflow <workflow-slug> --entity <entity-slug> --reason "<one-line reason>" --command "<command summary>" [--prefix-rule "git,-C"]
```

The helper appends `_bridge/fo-alerts.jsonl`, which Bridge renders as an Approve/Deny alert. The helper's returned `id` is the alert join key; Bridge echoes it back as `request_id` on the typed `permission-decision` inbox record. `deny` leaves the action blocked, `approve-once` retries the exact blocked action once through the runtime's escalation path, and `approve-rule` retries with the proposed reusable prefix rule when one was present. This is an FO intent signal, not a bypass of host-native security; if the host still presents its own approval prompt, honor that prompt normally.
