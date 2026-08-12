---
title: Use Codex subscription auth in live CI
status: ideation
source: captain-request
started: 2026-08-12T14:26:14Z
completed:
verdict:
worktree:
id: hrercm3ff4ww94rnqhqbqkyp
gates:
    version: 1
    records:
        - id: gate:hrercm3ff4ww94rnqhqbqkyp:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:hrercm3ff4ww94rnqhqbqkyp-backlog-1
              briefing:
                id: briefing:hrercm3ff4ww94rnqhqbqkyp:backlog:attempt-1:revision-1
                digest: sha256:4725aa1d24ef0107e1d244fd1e666415c73cded25f72011060f8069526344c87
                request-digest: sha256:2f934a126b11db2e1a45db3445cf1ed198de1992864f2c63f7b69a8c2b39fa7a
                room-ref: ./codex-subscription-live-auth/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:hrercm3ff4ww94rnqhqbqkyp:backlog:1
                briefing: briefing:hrercm3ff4ww94rnqhqbqkyp:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-12T14:24:42.50027Z"
                decision: approve
                reason: Captain approved the seed. A backlog stage report is not required for this gate.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:hrercm3ff4ww94rnqhqbqkyp:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:hrercm3ff4ww94rnqhqbqkyp-ideation-1
              briefing:
                id: briefing:hrercm3ff4ww94rnqhqbqkyp:ideation:attempt-1:revision-1
                digest: sha256:d0e95d39cb3defd5c5d9cc0a31d67d4f45a17f2861024517e2d773695a41bc0a
                request-digest: sha256:308aaa9d79291082cbab9cb6bccade2d21413dcef2236c5d6ea3db1c144a3434
                room-ref: ./codex-subscription-live-auth/review/ideation/briefing-1
---
# Use Codex subscription auth in live CI

Configure the Codex and Pi live CI lanes to prefer Codex subscription OAuth.

Keep `OPENAI_API_KEY` as a fallback during the migration. A lane must fail only when neither its OAuth secret nor `OPENAI_API_KEY` is available.

The implementation must preserve isolated homes and must not print OAuth secrets.

## Scope

Included:

- Codex CI authentication selection and isolated `CODEX_HOME` seeding.
- Pi CI authentication selection and isolated Pi-home seeding.
- GitHub Environment secret instructions for both OAuth file formats.
- Workflow and authentication tests.
- Runtime documentation for rotation and API-key fallback.

Excluded:

- Changes to the live model or thinking level.
- Changes to local developer authentication.
- Automatic persistence of refreshed OAuth tokens back to GitHub secrets.

## Problem statement and value

The current Codex live job treats `OPENAI_API_KEY` as mandatory and logs in with it. The current Pi job has the same preflight requirement; its isolated-home harness can copy an operator OAuth file locally, but CI has no subscription payload to seed and its configured model is not the `openai-codex` provider. This makes the approval-gated lanes spend an API key even when the maintainer has a Codex subscription, and it prevents a credential-free migration test of the subscription path.

The value is a live, approval-gated Codex/Pi lane that uses the subscription OAuth already paid for, while retaining a deterministic API-key fallback during rotation. The lane still runs the existing common and substrate journeys, so the value is measured by successful current-checkout workflow execution and durable artifacts rather than by a YAML or prose check. Credentials remain confined to isolated homes and setup/status output never contains their payloads.

## Proposed approach

### Secret names, formats, and precedence

Use two new secrets in the existing GitHub Environments; keep the names distinct because the hosts consume different file formats:

- `CI-E2E-CODEX` environment: `CODEX_AUTH_JSON`, the complete contents of the Codex `~/.codex/auth.json` file. The payload is the object with `auth_mode`, `tokens` (including the refresh token and account data), and `last_refresh`; no wrapper or base64 encoding is added.
- `CI-E2E-PI` environment: `PI_OPENAI_CODEX_AUTH_JSON`, only the value of Pi `auth.json`'s `openai-codex` entry. The record has `type: "oauth"`, `access`, `refresh`, `expires`, and `accountId`; the harness wraps this record under the `openai-codex` key when writing the isolated Pi `auth.json`.
- Both environments retain `OPENAI_API_KEY` as the migration fallback.

The CI decision is OAuth secret first, API key second, and fatal only when both are absent while the corresponding `SPACEDOCK_*_LIVE_REQUIRED` guard is set. When both are present, the API key is omitted from the child environment and the OAuth file is authoritative. Local developer authentication keeps its existing behavior; this task changes only the CI secret path.

| CI inputs | Codex action | Pi action | expected result |
| --- | --- | --- | --- |
| OAuth secret only | seed isolated `CODEX_HOME/auth.json`, run `codex login status` | seed isolated Pi `auth.json`, use `openai-codex/gpt-5.6-luna:max` | subscription path runs without an API key |
| OAuth secret plus API key | same OAuth actions; do not pass the key to the child | same OAuth actions; do not pass the key to the child | OAuth preference is observable |
| API key only | `codex login --with-api-key` through stdin, then `codex login status` | pass the key and use `openai/gpt-5.6-luna:max` | fallback path runs |
| Neither | non-secret preflight/decision error | non-secret preflight/decision error | lane fails before a host launch |

A non-empty secret is written with `0600` permissions and is never interpolated into an argv, log line, artifact name, or final message. Host status is the authority that the payload is usable; malformed or rejected OAuth payloads produce a redacted auth failure rather than dumping the file. Refreshed tokens may be written only inside the throwaway isolated home and are never copied back to the GitHub secret.

### Codex lane

Extend the existing `codexLiveAuthDecision` and isolated-home helpers instead of introducing a new runner. The live runner reads `CODEX_AUTH_JSON` from the CI environment, writes it as the isolated home's `auth.json`, and invokes `codex login status` without stdin on the OAuth branch. The existing API-key login branch remains the fallback. The workflow preflight accepts either `CODEX_AUTH_JSON` or `OPENAI_API_KEY`; `SPACEDOCK_CODEX_LIVE_REQUIRED` continues to make the no-credential case fail in CI. The environment builder must pass `OPENAI_API_KEY` only for the API-key branch and must drop any inherited `CODEX_AUTH_JSON` before launching the child.

### Pi lane

Keep Pi's isolated `PI_CODING_AGENT_DIR`, session directory, clean `HOME`, and package discovery. Add a pure auth selector and a seeding helper next to the current Pi live harness. The OAuth branch parses the dedicated record, writes `{ "openai-codex": <record> }` as `auth.json` with `0600` permissions, drops `OPENAI_API_KEY`, and selects the exact `openai-codex/gpt-5.6-luna:max` model. The fallback branch writes no OAuth record, passes `OPENAI_API_KEY`, and selects the equivalent API provider model `openai/gpt-5.6-luna:max`. The model ID and `max` thinking level stay unchanged; only the provider qualifier follows the selected auth mechanism. Explicit local model overrides remain full provider-qualified values and are not silently rewritten.

### Isolation and no-secret boundary

Use the existing temp-home and artifact plumbing. Do not use the operator's home as `CODEX_HOME` or `PI_CODING_AGENT_DIR`, do not copy unrelated plugin/session/config state, and do not use shell tracing around secret-handling steps. Fake-host tests capture argv, environment, stdout, stderr, and seeded on-disk files; they assert sentinel credentials are absent from every captured output while still observing the selected branch.

## Expected file surface and tolerance

Expected implementation surface is one workflow, one runtime document, the existing Codex auth/runner tests, and the existing Pi auth/runner/control tests:

- `.github/workflows/runtime-live-e2e.yml` — CI secret bindings, either-credential preflights, and the Pi provider-qualified model defaults.
- `internal/ensigncycle/codex_liveenv.go`, `codex_liveenv_test.go`, and `codex_live_runner_test.go` — selection, isolated seeding, status invocation, and output-boundary fixtures.
- `internal/ensigncycle/pi_live_runner_test.go`, `pi_shared_live_runner_test.go`, and `pi_live_controls_test.go` — Pi selection/seeding, mode-dependent model and environment, and isolation fixtures. A small `pi_liveenv.go`/`pi_liveenv_test.go` split is allowed if the build-tag boundary otherwise prevents offline tests.
- `docs/runtime-live-ci.md` — secret formats, Environment setup, rotation, refresh behavior, fallback, and the exact model examples.

Estimate: 8 modified files and about 220–340 insertions, with a tolerance of plus or minus 2 files and 80 lines. At most two focused helper/test files may be added; no new dependency is expected. The implementation must not touch the launcher command grammar, production workflow-state code, skills, live model ID, thinking level, local developer auth behavior, or automatic secret persistence.

## Declared semantic changes

- **Command grammar:** none.
- **Stored formats:** two new GitHub Environment secret contracts only; repository entities and local auth files are unchanged.
- **Authority:** in CI, the dedicated host OAuth secret is authoritative when present; `OPENAI_API_KEY` is the explicit fallback. Local auth remains outside this change.
- **Runtime behavior:** Codex gets a seeded isolated `CODEX_HOME`; Pi gets a seeded isolated home and selects `openai-codex/...` for OAuth or `openai/...` for the key fallback. The model ID and thinking level do not change.
- **Security behavior:** neither OAuth payload nor API key is printed, and refreshed credentials never flow back to GitHub.

## Risks and mitigations

1. **Refresh-token expiry or revocation.** The stored OAuth record can be syntactically valid but rejected by the host. `codex login status`, the real Pi launch, and redacted diagnostics make this a normal auth failure; rotation instructions replace the Environment secret without persisting refreshed tokens automatically.
2. **Host format drift.** Codex and Pi do not consume the same JSON shape. Exact-format fixtures, the pinned Pi 0.80.10 probe, and JSON parsing before file writes catch accidental wrapping or field renames. The live Codex package floats to the workflow's selected version, so the OAuth-only CI leg remains required evidence.
3. **Provider/model mismatch.** Pi's `openai-codex` provider is OAuth-only while `openai` is the API-key provider. A mode-to-model table and exact model assertions prevent a key fallback from launching an OAuth-only model or silently changing the model ID.
4. **Secret leakage through logs or artifacts.** GitHub masking is not a sufficient product boundary. `printf`-to-file, `0600` files, dropped child-env secrets, and fake-host output scans prove the boundary independently of GitHub masking.
5. **Missing migration secret.** Preflights test the OR condition and leave `OPENAI_API_KEY` usable; only the neither-credential case fails. Documentation names the two Environment scopes and gives non-printing rotation commands.

## Simplest sufficient mechanism and alternatives rejected

The smallest mechanism is to reuse each existing live runner's auth decision, isolated-home seed, and host launch, adding only the dedicated secret input, the mode-dependent provider model, and deterministic fixtures. No shared auth package, token broker, API client, or refresh writer is needed.

| mechanism | value AC served | simplest alternative considered | why the alternative is insufficient |
| --- | --- | --- | --- |
| Write the dedicated payload into the existing isolated home | AC-1, AC-2, AC-4 | Pass the JSON as an environment variable to the host | Codex and Pi read files, not arbitrary secret variables; leaving the payload in the child env also widens the leak surface. |
| Select Pi's provider-qualified model from auth mode | AC-2, AC-4 | Keep one `openai/...` model for both branches | The OAuth record belongs to `openai-codex`; the API-key branch belongs to `openai`, and the provider is part of Pi's model identity. |
| Run host-native Codex status and the real Pi smoke | AC-1, AC-2, AC-4 | Check only that a GitHub secret variable is non-empty | Presence cannot detect an unusable payload, wrong file shape, wrong provider, or an isolated-home launch failure. |
| Capture fake-host output and scan it for sentinels | AC-1, AC-4 | Rely on GitHub's masking feature | Masking is platform behavior and does not prove local tests, subprocess env construction, or uploaded artifacts are safe. |
| Do not write refreshed OAuth credentials back | scope/security boundary | Persist the changed file to GitHub on every run | That is explicitly excluded, races with the secret owner, and would turn a CI run into a credential rotation authority. |

## Acceptance criteria

**AC-1** The Codex approval-gated lane reaches the existing current-checkout live suite with OAuth-only inputs, and a key-only fixture reaches the same runner through API-key login; with neither input the required lane stops before host launch. The isolated `CODEX_HOME` contains only the minimal config plus the selected auth file, `codex login status` is exercised on the OAuth branch, and neither credential appears in captured output.

- **Test plan:** deterministic Go table tests cover OAuth preference, API-key fallback, and the required missing-input decision; file tests cover exact auth copy, `0600` mode, and no unrelated home state; a fake Codex shim records argv/env/status without printing values; the approval-gated `codex-live` job supplies the OAuth-only proof against the real host.

**AC-2** The Pi approval-gated lane reaches the existing front-door/common smoke with an isolated home seeded from the dedicated `openai-codex` record and uses the exact `openai-codex/gpt-5.6-luna:max` identity; a key-only fixture uses `openai/gpt-5.6-luna:max`; both branches preserve isolated homes and neither credential appears in captured output.

- **Test plan:** offline Go fixtures cover record wrapping, exact mode-dependent model selection, key dropping on OAuth, key retention on fallback, path isolation, and file permissions; the pinned Pi 0.80.10 model-list probe checks both provider/model identities; the approval-gated `pi-live` front-door/common run supplies the real OAuth evidence.

**AC-3** `docs/runtime-live-ci.md` contains the exact two secret names, both JSON formats, the Environment scopes, non-printing setup/rotation commands, refresh-token handling, and the API-key fallback/expiry guidance; it states that refreshed tokens remain isolated and are not written back to GitHub.

- **Test plan:** documentation review checks each required field and command against the workflow and harness names; a focused text-shape test may assert the secret names, provider-qualified model strings, and no-secret command shape, while behavioral auth claims remain covered by AC-1/AC-2 rather than static prose alone.

**AC-4** Workflow and authentication fixtures observe OAuth preference, API-key fallback, provider/model selection, isolated-home contents, login/status exit behavior, and no-secret-output boundaries; the end-value is a successful OAuth-only current-checkout Codex/Pi run with the existing durable artifacts, measured against the key-only fallback fixture and the neither-credential failure case.

- **Test plan:** run the focused default Go suite and fake-host workflow controls for all branches, then run the required offline gate and the approval-gated Codex/Pi live commands; review exit codes, model/session records, artifacts, and on-disk home paths, not transcript assertions alone. Cost/complexity is low-medium: fixture work is deterministic, while one OAuth-only Codex suite and one Pi front-door/common suite consume the live budget.

## Provider-auth spike (2026-08-12)

The riskiest unverified mechanism was whether the two installed hosts accept copied subscription records in isolated homes and whether the pinned Pi package exposes the exact provider-qualified model. A no-model, isolated-home spike was run with output redirected to temporary files; only exit codes, byte counts, model names, and a credential-pattern scan were reported.

- Codex CLI `0.147.0`: copied the existing operator `~/.codex/auth.json` to a fresh `CODEX_HOME`, cleared `OPENAI_API_KEY`, and ran `codex login status`. Result: exit `0`, zero stdout, and the combined output secret scan was clean.
- Pi coding agent `0.80.10` (the workflow pin): copied the existing operator Pi `auth.json` and `models.json` to a fresh `PI_CODING_AGENT_DIR`, cleared `OPENAI_API_KEY`, and ran `pi --list-models`. Result: exit `0`, `openai-codex/gpt-5.6-luna` was listed, and the output secret scan was clean.
- Pi coding agent `0.80.10` API-key fixture: used an empty auth file plus a non-printing sentinel `OPENAI_API_KEY` and ran `pi --list-models`. Result: exit `0`, `openai/gpt-5.6-luna` was listed, and neither the sentinel nor OAuth-shaped fields appeared in output.

No model request was made and no credential value was printed. The spike supports the proposed file shapes and provider split; implementation still owes fake-host branch tests and one approval-gated OAuth-only live run because local status/model listing does not prove a complete workflow journey.

## Documentation diff to apply

Insert under `### GitHub setup` in `docs/runtime-live-ci.md`:

```diff
+#### Subscription OAuth secrets
+
+The `CI-E2E-CODEX` Environment accepts `CODEX_AUTH_JSON`, the complete
+`~/.codex/auth.json` payload. The `CI-E2E-PI` Environment accepts
+`PI_OPENAI_CODEX_AUTH_JSON`, only the `openai-codex` object from Pi's
+`~/.pi/agent/auth.json`. Keep `OPENAI_API_KEY` in each Environment during the
+migration; the lane uses it only when its OAuth secret is absent.
+
+Create or rotate the secrets from a trusted workstation without printing them:
+
+    gh secret set CODEX_AUTH_JSON --env CI-E2E-CODEX < "$HOME/.codex/auth.json"
+    jq -c '."openai-codex"' "$HOME/.pi/agent/auth.json" \
+      | gh secret set PI_OPENAI_CODEX_AUTH_JSON --env CI-E2E-PI
+
+The Codex value must remain the full file object (`auth_mode`, `tokens`, and
+`last_refresh`). The Pi value must remain one OAuth record (`type`, `access`,
+`refresh`, `expires`, and `accountId`); the runner supplies the outer provider
+key. Never echo either value or include it in an artifact. A refresh may update
+only the isolated run home. Re-run the host login flow and replace the relevant
+Environment secret when the refresh token is revoked or the host reports an
+expired credential; CI never writes refreshed tokens back to GitHub. If the
+OAuth secret is absent, the existing `OPENAI_API_KEY` is used; a lane fails
+before launch only when both credentials are absent.
+
+For the OAuth branch Pi's exact model is
+`openai-codex/gpt-5.6-luna:max`; the API-key fallback is
+`openai/gpt-5.6-luna:max`. The model ID and `max` thinking level are unchanged.
```

This is a documentation-only addition; no local-auth instructions are removed and no live model/thinking setting is changed.

## Stage Report: ideation

- DONE: Produce a problem statement and value framing for Codex and Pi CI subscription authentication with API-key fallback.
  The `Problem statement and value` section ties the current mandatory-key behavior to the measurable value: OAuth-only live workflow execution with the existing durable artifacts and a key-only fallback.
- DONE: Record one concrete approach, expected file surface and tolerance, semantic changes, risks, and the simplest sufficient mechanism.
  `Proposed approach`, `Expected file surface and tolerance`, `Declared semantic changes`, `Risks and mitigations`, and the alternatives table define the bounded implementation and its necessity.
- DONE: Update the acceptance criteria with test plans, a provider-auth spike result or no-spike record, and the documentation diff needed for implementation.
  AC-1 through AC-4 each name observable tests; the pinned-host spike records exit/model/secret-scan evidence; `Documentation diff to apply` gives the exact insertion and rotation commands.

### Summary

The design keeps the existing isolated live runners and adds only two Environment secret contracts, mode-dependent auth/model selection, and deterministic no-secret fixtures. The provider spike proved Codex status and both pinned Pi provider/model identities without a model request or credential output; implementation's remaining proof is the fake-host tests plus approval-gated OAuth-only lanes.
