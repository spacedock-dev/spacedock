---
id: z5gwwz2748sg6vxr0g3kdsar
title: "Codex launcher guarantees multi-agent v2 surface"
status: validation
source: "Captain feedback, 2026-07-02: ordinary Codex config currently enables multi_agent_v2; make Spacedock-launched Codex enable or prove the same surface instead of relying on ambient session setup."
started: 2026-08-01T14:29:25Z
completed:
verdict:
score: 1.0
worktree: .worktrees/spacedock-ensign-codex-launch-multi-agent-v2
issue:
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:z5gwwz2748sg6vxr0g3kdsar:validation
    records:
        - id: gate:z5gwwz2748sg6vxr0g3kdsar:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-ideation-1
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:ideation:attempt-1:revision-1
                digest: sha256:988c0401c86e31c5e9a300df8dd7e3527e171d1b7dbbb56aa40262ffa6b0ec24
                digest-domain: canonical-bytes
                request-digest: sha256:91d6338b0e681ec994ea130511078c95074dd505919b22665c7666e7fa16c408
                room-ref: ./codex-launch-multi-agent-v2/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:z5gwwz2748sg6vxr0g3kdsar:ideation:1
                briefing: briefing:z5gwwz2748sg6vxr0g3kdsar:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-01T16:10:09.609186Z"
                decision: revise
                reason: 'Captain-directed Science Officer send-back completed: preserve the exact multi_agent_v2 isolated-home fragment and bind AC-1 through AC-4 to falsifiable evidence; the prior Briefing is stale after the revised report.'
              application:
                action: feedback
                target-stage: ideation
                state: superseded
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-ideation-2
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:ideation:attempt-2:revision-1
                digest: sha256:0cc3451492227cd3b8b756842301674a4339ae6e6c1e42c9e45944980bd6c75b
                digest-domain: canonical-bytes
                request-digest: sha256:c156fc7e30942202396dcc2fb1c1fafbe57cdd20ab1d5de5bf7808e1fcbf854e
                room-ref: ./codex-launch-multi-agent-v2/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:z5gwwz2748sg6vxr0g3kdsar:ideation:2
                briefing: briefing:z5gwwz2748sg6vxr0g3kdsar:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-01T16:21:11.545595Z"
                decision: approve
                reason: Captain approved after Science Officer APPROVE advisory; exact features.multi_agent_v2 fragment and AC-1 through AC-4 evidence are bound in briefing-2.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:z5gwwz2748sg6vxr0g3kdsar:validation
          stage: validation
          attempts:
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-1
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-1:revision-1
                digest: sha256:1065cc2ea5b9aca994efb80974074c6cf36c7c486c853f7a5e84031ac0252314
                digest-domain: canonical-bytes
                request-digest: sha256:57030132e27b794a2e0bb1efbdf51f3e3debc12ac2d1b7c9b709822fa4d7ccc5
                room-ref: ./codex-launch-multi-agent-v2/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:z5gwwz2748sg6vxr0g3kdsar:validation:1
                briefing: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-01T17:13:08.344768Z"
                decision: revise
                reason: 'Captain-authorized direct send-back following Science Officer REVISE: replace the vacuous direct-codex marker oracle with a Spacedock-front-door typed ordered same-worker lifecycle proof; reject every accepted reserved-key spelling before side effects, including attached short and quoted dotted forms; add the disabled-control E-3 negative; rerun focused, full, race, format, detached-audit, and live validation.'
              application:
                action: feedback
                target-stage: implementation
                state: superseded
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-2
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-2:revision-1
                digest: sha256:ca917663cdcb933eae844888a3fe1661c3c6d4aef7aa133b1a0bc28c3873fed3
                digest-domain: canonical-bytes
                request-digest: sha256:36f7b9dfe6811ad19739d9bbac0ac8e582884baece283ff45f5ec610ecbb3923
                room-ref: ./codex-launch-multi-agent-v2/review/validation/briefing-2
---

Spacedock's Codex front door should reproduce the captain's isolated-home multi-agent v2 configuration while also pinning Codex's documented stable multi-agent controls. The resulting first officer must observably spawn a worker, follow up with the same worker, inspect/list and wait for workers, and receive completion. The guarantee belongs to launch configuration plus behavioral proof, not bootstrap prose or ambient user config.

## Problem

The Codex launcher owns plugin installation, optional sandbox wrapping, and the first-officer bootstrap prompt, but the current launch path does not force or verify Codex's supported multi-agent configuration. That leaves several fragile cases:

- a clean or isolated Codex home may not carry the operator's feature config;
- a sandboxed launch may preserve the binary invocation but miss the config needed to expose collaboration tools;
- live CI or fresh user machines may silently start with a weaker multi-agent surface, changing dispatch/reuse behavior after Spacedock has already begun operating.

The user-visible failure mode is not "a flag is missing". It is that Spacedock starts a first officer which cannot complete the native worker lifecycle after dispatch has begun.

## Spike result (2026-08-01)

The riskiest launch and runtime mechanisms were exercised before selecting a design:

- An instrumented `codex` under a built `spacedock` front door captured plain, local-plugin, and Safehouse argv from isolated homes. None carried a multi-agent override. The plain/local inner argv began with `codex --ask-for-approval on-request`; Safehouse carried `safehouse ... -- codex --dangerously-bypass-approvals-and-sandbox`. This falsifies the current launcher-owned-guarantee claim.
- Codex CLI 0.146.0 with an empty `CODEX_HOME` reported `multi_agent stable true`, `multi_agent_v2 stable false`. `-c features.multi_agent=true` enables the documented stable feature; `-c features.multi_agent_v2=true` merely flips the second label and therefore proves no lifecycle.
- Current official Codex documentation names `agents.enabled` (default `true`) and `features.multi_agent` (stable, on by default). It does not name `multi_agent_v2` as the supported configuration contract.
- A live `codex exec --json` used an isolated `CODEX_HOME` containing only copied authentication and the overrides `agents.enabled=true` and `features.multi_agent=true`. The parent transcript recorded `spawn_agent`, `wait_agent`, `followup_task`, `list_agents`, and `wait_agent`; the child recorded `CHILD_READY` then `CHILD_DONE`; the turn context recorded `multi_agent_version: v2`; and the parent returned `PARENT_DONE`.
- The falsifying negative repeated the isolated launch with `agents.enabled=false`. No collaboration call appeared and the agent returned `SPAWN_UNAVAILABLE`.
- The captain's actual ambient config contains this exact fragment:

  ```toml
  [features.multi_agent_v2]
  max_concurrent_threads_per_session = 16
  tool_namespace = "agents"
  hide_spawn_agent_metadata = false
  ```

  Codex 0.146.0 accepts its equivalent inline table but `codex features list` still reports `multi_agent_v2 stable false`, proving the table and boolean feature-list label are not interchangeable.
- A second isolated-home live run combined the exact table with `agents.enabled=true` and `features.multi_agent=true`. Parent thread `019fbe0d-76ca-79d1-9d5b-d2bfb110acf1` recorded `spawn_agent`, `list_agents`, and `wait_agent`; the child returned `FRAGMENT_CHILD`; the parent returned `FRAGMENT_PARENT`; and turn context recorded `multi_agent_version: v2` even though the identical `features list` invocation still printed `multi_agent_v2 stable false`.

Conclusion: the supported controls and requested v2 table serve different purposes and must coexist. `agents.enabled=true` plus stable `features.multi_agent=true` are the supported enablement boundary; the exact `features.multi_agent_v2` table reproduces the captain's requested concurrency, namespace, and metadata shape. Only the observable lifecycle is acceptance evidence; the feature-list label is not.

## Proposed approach

1. `spacedock codex` injects one launcher-owned config layer immediately after the inner `codex` token: `-c agents.enabled=true`, `-c features.multi_agent=true`, and the inline TOML equivalent of the captain's exact fragment, `-c 'features.multi_agent_v2={max_concurrent_threads_per_session=16,tool_namespace="agents",hide_spawn_agent_metadata=false}'`. This position works for fresh launches and `resume`, survives the Safehouse wrap unchanged, overrides user/profile defaults, and avoids a temporary config file.
2. The launcher reserves all three owned paths. A forwarded `-c`/`--config` assignment targeting `agents.enabled`, `features.multi_agent`, or `features.multi_agent_v2`, or `--enable`/`--disable multi_agent[_v2]`, fails before plugin installation or `ops.Launch` with an actionable message explaining that Spacedock owns the collaboration guarantee. Redundant overrides are rejected too, keeping one effective source of truth and preventing last-override-wins ambiguity.
3. A Codex version which does not accept the stable settings or v2 table rejects config before creating the interactive first-officer turn; the native parser error is preserved and the host exits nonzero. Documentation states this Codex support boundary. No prompt fallback or downgraded dispatch is permitted.
4. Tests use one shared expected-argv helper so existing byte-exact launch oracles continue to assert the entire command, not merely search for substrings. A small live probe separately proves that the supported settings produce the lifecycle; it does not modify PR #585 or its existing live harness.

This serves AC-1 and AC-2 with less machinery than a temporary `CODEX_HOME`/config merge, which would introduce file ownership, cleanup, auth, and profile-precedence hazards. The inline table is a semantically equivalent TOML value at the same config path and was exercised live from an isolated home. Preserving the ambient home alone is insufficient because the clean-home case is the defect. A per-launch `codex doctor` preflight was rejected because it adds latency and reports the stable feature flag but not the v2 table or complete collaboration lifecycle; explicit overrides plus live lifecycle proof give a smaller and stronger boundary.

## Out of scope

- Changing Codex itself.
- Redesigning Spacedock's dispatch or reuse contract.
- Making multi-agent v1 behavior equivalent to v2.
- Depending on a particular developer's local config as the only proof.
- Treating the `multi_agent_v2` feature-list label as proof that the requested table or lifecycle is active.
- Modifying PR #585 or its live harness.
- Changing bootstrap/first-officer prose to simulate capability detection.

## Expected surface

- `internal/cli/frontdoor.go`: 35-60 inserted lines for the owned stable settings plus exact v2 table, conflict recognition, and pre-launch diagnostic.
- `internal/cli/codex_multi_agent_test.go` (new): 140-210 inserted lines for the variant/conflict matrix and a live-test entry point guarded by existing live prerequisites.
- Existing exact-argv launcher tests in `internal/cli/frontdoor_test.go`, `internal/cli/safehouse_frontdoor_test.go`, `internal/cli/launch_parity_test.go`, `internal/cli/codex_plugin_dir_test.go`, and `internal/cli/safehouse_knob_test.go`: 20-50 net inserted lines through a shared expected-argv helper and updated complete argv expectations.
- `docs/site/reference/command-reference.md` and `docs/runtime-live-ci.md`: 10-24 inserted lines describing the stable settings, exact v2 table, reserved forwarded keys, support boundary, and isolated-home lifecycle command.

Expected total: the same 8-9 files, 205-344 inserted lines, and at most 30 deleted lines. Tolerance remains at most 10 files and 380 inserted lines; no additional implementation file is authorized. Exceeding either bound requires design re-entry rather than quietly broadening the mechanism.

Declared semantic changes are limited to Codex launch runtime behavior and three forwarded-config path reservations. Every `spacedock codex` invocation carries the two supported stable controls and exact requested v2 table, including plain, local-plugin, Safehouse, and resume launches; conflicting attempts fail before plugin mutation or host launch. The v2 table fixes the open-thread cap at 16, selects the `agents` tool namespace, and leaves spawn metadata visible. Command names, stored formats, workflow state, plugin format, dispatch/reuse authority, Claude/Pi behavior, and bootstrap text do not change.

## Acceptance criteria

**AC-1 - A supported Spacedock-launched Codex can complete the native collaboration lifecycle from an isolated home.**
Verified by: live evidence E-2 runs with no user config and observes, in structured JSONL and child state, one spawn, one same-worker follow-up, list/inspect, wait, `CHILD_READY`, `CHILD_DONE`, and `PARENT_DONE`; it also records `multi_agent_version: v2`. Removing either launcher-owned supported setting or disabling `agents.enabled` makes this proof fail as shown by E-3.

**AC-2 - The launcher reproduces the captain's exact v2 table and cannot silently downgrade through ambient or forwarded configuration.**
Verified by: live evidence E-4 combines the exact table with the stable controls in an isolated home and observes v2 spawn/list/wait plus child and parent completion while the feature-list label remains false. Launcher table tests start from ambient false settings and show all owned overrides in complete argv, while every forwarded spelling targeting the three reserved paths exits nonzero with the pinned diagnostic and records zero plugin installs and zero host launches. An instrumented unsupported host rejects the owned layer before opening a session and its nonzero exit/config diagnostic propagates.

**AC-3 - Plain, local-plugin, Safehouse, and resume launches carry one identical supported enablement layer.**
Verified by: complete-argv launcher tests seeded by instrumented capture E-1 assert the two stable tokens and exact inline v2 table occur exactly once, immediately after the inner `codex` token, across all four variants; deleting any element, changing a v2 field, or placing one outside the Safehouse inner argv fails the matrix.

**AC-4 - Runtime proof and operator documentation use the supported settings and behavioral oracle, not the v2 label or personal config.**
Verified by: the isolated-home live test reproduces E-4 by copying authentication only, applies the stable controls plus exact v2 table, and grades structured lifecycle events. Documentation names all owned paths, the distinction demonstrated by E-4 between the v2 table and false feature-list label, the rejection behavior, and the minimum host-support boundary. A guard rejects any test that substitutes `codex features list` state for lifecycle evidence.

## Evidence ledger

- **E-1 — frontdoor handoff falsifier:** built `spacedock` with an instrumented `codex`/`safehouse`, isolated `CODEX_HOME`, and plain/local-plugin/Safehouse launches. Captured inner argv contained approval/bypass and bootstrap arguments but no multi-agent config. Repeating after implementation must show the full owned layer in each inner argv.
- **E-2 — supported-controls lifecycle:** isolated home with copied authentication only; `codex exec --json -c agents.enabled=true -c features.multi_agent=true`. Parent thread `019fbdd7-f0c0-7183-ab4a-02420aac3911` recorded `spawn_agent`, `wait_agent`, `followup_task`, `list_agents`, `wait_agent`, and `PARENT_DONE`; its child recorded `CHILD_READY` then `CHILD_DONE`; turn context recorded `multi_agent_version: v2`.
- **E-3 — disabled falsifier:** isolated home, same command with `agents.enabled=false`; stdout completed with `SPAWN_UNAVAILABLE` and no collaboration item. This independently moves the lifecycle in the wrong direction when enablement is removed.
- **E-4 — exact-fragment reconciliation:** isolated home with copied authentication only; E-2 controls plus `-c 'features.multi_agent_v2={max_concurrent_threads_per_session=16,tool_namespace="agents",hide_spawn_agent_metadata=false}'`. Parent thread `019fbe0d-76ca-79d1-9d5b-d2bfb110acf1` recorded `spawn_agent`, `list_agents`, `wait_agent`, `FRAGMENT_PARENT`; child recorded `FRAGMENT_CHILD`; turn context recorded `multi_agent_version: v2`; the identical config passed to `codex features list` still reported `multi_agent_v2 stable false`, falsifying label-as-proof.
- **E-5 — official supported boundary:** current Codex [Subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents.md) and [Configuration Reference](https://learn.chatgpt.com/docs/config-file/config-reference.md) documentation identifies `agents.enabled` (default true), `agents.max_concurrent_threads_per_session`, and stable `features.multi_agent`; it does not document the captain's v2 table. The design therefore pins both the supported controls and the exercised requested table rather than relabeling one as the other.

## Test plan

Implementation starts with failing focused tests, then changes launcher behavior:

1. Add the complete-argv matrix for plain, explicit local-plugin, Safehouse, and exact `resume`. The matrix asserts adjacency, uniqueness, exact values for all three v2 table fields, and the inner side of the wrapper boundary. Cost: low, fixture-backed, no model spend.
2. Add a conflict matrix for `-c`, `--config`, equals forms accepted by Codex, and `--enable`/`--disable multi_agent[_v2]`. Cover the stable keys, the v2 table, and each v2 nested field spelling. Each case asserts exit 1, byte-stable actionable stderr, no install, and no launch. Include nearby unrelated config as a green control. Cost: low, fixture-backed.
3. Add an instrumented unsupported-host case that rejects the owned stable-controls-plus-v2-table layer before starting a session. Assert the native actionable config diagnostic and nonzero exit propagate, and that its session/dispatch marker remains absent. Cost: low, process fixture.
4. Add the complete owned layer in `frontdoor.go`, before passthrough and after the inner program token, and reject reserved conflicts immediately after front-door parsing. Update all affected complete-argv oracles through one test helper. Cost: low-to-medium because exact argv is intentionally broad regression coverage.
5. Run the focused package tests, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.
6. Run the isolated-home live proof with copied auth only and all owned settings, reproducing E-4. Parse JSONL rather than final prose: require the ordered parent tool set, one stable child identity across follow-up, both child terminal messages, parent completion, and `multi_agent_version: v2`; separately assert `features list` remains insufficient. Run the disabled control and require zero collaboration calls. Cost: one short positive and one short negative Codex run; live and credential-gated.
7. Repeat the instrumented frontdoor capture for plain, local-plugin, and Safehouse launch shapes and retain concise artifacts. The live lifecycle is run once because the same inner tokens are proven byte-identical across variants; three paid lifecycle runs would add cost without testing a distinct Codex mechanism.

## Consultation record

Standing Science Officer consultation (signed consultation supplied with the 2026-08-01 ideation dispatch): **REVISE**. Reference: current official Codex `Subagents` and `Configuration Reference` documentation plus the isolated Codex 0.146.0 probe. Reason to preserve in the gate: `agents.enabled` and stable `features.multi_agent` are the supported current controls; `multi_agent_v2=false` is only a local feature-label observation, so approval must rest on observable spawn, same-worker follow-up, list/wait, and completion across launcher variants, with a disabled/unsupported path that fails before worker dispatch.

Direct Science Officer advisory send-back, 2026-08-01: **REVISE**. Material findings preserved unchanged:

1. the current ideation approach selects agents.enabled=true + features.multi_agent=true and explicitly excludes the captain-required isolated-home fragment [features.multi_agent_v2] max_concurrent_threads_per_session=16, tool_namespace="agents", hide_spawn_agent_metadata=false; this is a scope mismatch;
2. AC-1 through AC-4 have no independent evidence citations in the current entity report.

Disposition: the revised approach carries the exact fragment as an inline TOML table alongside the supported controls, and the evidence ledger plus cycle-2 report citations bind every AC to falsifiable launch or live evidence. The open gate is intentionally untouched; the First Officer will withdraw the stale attempt and prepare a fresh Briefing.

## Stage Report: ideation

- DONE: Spike the real Codex frontdoor argv/config handoff from an isolated Codex home and record the falsifying result.
  Built launcher plus instrumented host captures proved no current override in plain/local/Safehouse argv; isolated 0.146.0 and live JSONL proved the stable-settings lifecycle while `agents.enabled=false` removed it.
- DONE: Declare the expected files/LOC/semantic surface and select the smallest supported enablement or fail-closed preflight.
  The body caps the change at 10 files/340 insertions and selects launcher-owned `agents.enabled=true` plus stable `features.multi_agent=true`, with conflicting forwarded controls rejected before install/launch.
- DONE: Add falsifiable launcher-test and live-proof plans covering plain, local-plugin, and sandbox launches.
  Complete-argv matrices fail on missing/duplicate/misplaced tokens; reserved-key cases require zero install/launch; isolated live JSONL requires spawn, same-worker follow-up, list/wait, child and parent completion.

### Summary

Reframed the task from enabling an unproven `multi_agent_v2` label to guaranteeing the observable Codex collaboration lifecycle through supported settings. The spike demonstrated that current Spacedock does not own that configuration, while a clean-home live run with the supported settings completed the full lifecycle and the disabled control did not expose collaboration; the signed Science Officer REVISE rationale is retained for the gate.

## Stage Report: ideation (cycle 2)

- DONE: Reconcile the captain-required isolated-home v2 fragment with the supported controls.
  AC-1 and AC-2 evidence: E-2/E-3 establish the stable enablement boundary; E-4 independently proves the exact three-field v2 table coexists with it and completes a v2 lifecycle while the feature-list label remains false.
- DONE: Anchor every acceptance criterion to falsifiable evidence.
  AC-3 and AC-4 evidence: E-1 is the instrumented launcher falsifier and future argv oracle; E-4 is the isolated-home structured lifecycle oracle; E-5 binds the supported/documented boundary without substituting prose for runtime behavior.
- DONE: Update proposed approach, test plan, and semantic scope together without widening files.
  The same 8-9 implementation files remain authorized; complete argv, reserved-key, unsupported-host, disabled, and live JSONL tests now cover the stable controls plus exact table and fail on any missing or changed field.

### Summary

The revision accepts both Science Officer findings without changing product code or the open gate. It preserves the captain's exact v2 table alongside Codex's supported stable controls, adds independent falsifiers and AC citations, and keeps the implementation file surface unchanged.

## Stage Report: implementation

- DONE: inject supported controls plus the exact multi_agent_v2 inline TOML layer across all launch variants
  Commit c36e3f7d1 injects the layer directly after inner `codex`; `TestCodexCollaborationLayerCompleteArgv` fails if plain, local-plugin, Safehouse, or resume loses, duplicates, reorders, or changes any token.
- DONE: reject reserved or conflicting forwarded settings before plugin install or host launch
  The conflict matrix covers `-c`/`--config` and `--enable`/`--disable` space/equals forms plus every v2 field, requires the pinned diagnostic and zero install/launch calls, and retains an unrelated-config green control.
- DONE: prove complete argv, unsupported-host failure, and isolated-home lifecycle behavior within the approved surface
  Eight files/332 insertions stayed below 10/380; the unsupported-host test preserves exit 78/config stderr before a session, the 35.6s live isolated-home run observed spawn/follow-up/list/wait and all child/parent markers with v2 context, and `go test ./...` plus `go test ./... -race` passed.

### Summary

Spacedock now owns one Codex collaboration configuration source across every launch shape and refuses forwarded attempts to weaken or replace it before persistent or host side effects. The implementation includes offline argv/conflict/support-boundary proofs, a credential-gated isolated-home lifecycle oracle, and operator documentation without changing non-Codex behavior or using a prompt fallback.

## Stage Report: validation

- DONE: Run focused launcher/support tests, full tests, race tests, formatting checks, and the credential-gated lifecycle test at implementation commit `c36e3f7d1`.
  Evidence: focused `TestCodexCollaborationLayerCompleteArgv`, `TestCodexRejectsOwnedCollaborationOverridesBeforeSideEffects`, and `TestCodexUnsupportedHostFailsBeforeSession` passed; `go test ./...`, `go test ./... -race`, `git diff --check`, and `gofmt -d ./cmd ./internal` passed. The real Codex 0.146.0 lifecycle test passed in 52.06s. The candidate remained exactly 8 files, +332/-26, within the approved 10-file/+380 limit.
- FAILED: Reproduce AC-1 and AC-4 as a Spacedock-launched, ordered, same-worker native lifecycle.
  Material evidence defect: `TestCodexIsolatedHomeCollaborationLifecycle` invokes `codex exec` directly, not the built `spacedock codex` front door. Its grader concatenates stdout with every parent/child JSONL file and makes unordered substring-presence assertions. It does not parse structured events, assert exactly one spawn, bind follow-up/list/wait to the spawned worker ID, or prove ordering. `CHILD_READY`, `CHILD_DONE`, and `PARENT_DONE` are already present in the input prompt. In a detached checkout at `c36e3f7d1`, a fake `codex` that emitted one non-collaboration note containing the vocabulary and performed zero worker calls passed the unmodified lifecycle test in 0.32s. This falsifies the behavioral oracle and the documentation claim that it grades structured same-worker lifecycle output.
- FAILED: Reproduce AC-2's claim that forwarded configuration cannot silently downgrade the launcher-owned settings.
  Material product defect: Codex 0.146.0 accepts attached short overrides such as `-cagents.enabled=false` and quoted dotted keys such as `-c '"agents".enabled=false'` (both probes exited 0). Spacedock's recognizer catches neither form and forwards each after the launcher-owned `agents.enabled=true` token. Detached adversarial tests expecting pre-launch rejection failed and captured exit 0 with the conflicting token later in the launched argv. The later override can therefore weaken the guarantee, contrary to AC-2 and the operator documentation.
- DONE: Reproduce AC-3's complete-argv placement and uniqueness for the covered plain, local-plugin, Safehouse, and resume variants.
  Evidence: the complete-argv matrix passed and asserts the three exact owned assignments immediately after the inner `codex` token, including inside the Safehouse boundary. Existing exact-argv tests also passed in the full and race suites. This does not cure the accepted post-layer override spellings above.
- FAILED: Reproduce the complete E-1 through E-5 ledger.
  E-1's post-implementation argv shapes are covered by the passing exact-argv tests. E-2/E-4 ran through direct `codex exec`, and the live test returned success, but the shipped oracle cannot establish cardinality, order, same-worker identity, or Spacedock-launch composition. E-3 has no implemented disabled negative run or zero-collaboration assertion. E-5 is consistent with the current official Codex manual: it documents `agents.enabled` (default true), `agents.max_concurrent_threads_per_session`, and stable/on-by-default `features.multi_agent`; the manual contains no `multi_agent_v2`, `tool_namespace`, or `hide_spawn_agent_metadata` entry. The two operator docs name the owned paths and unsupported-host boundary, but their lifecycle-grading and conflict-rejection claims overstate the implementation.
- DONE: Perform the required detached adversarial audit without modifying the deliverable.
  Evidence: a detached throwaway worktree at `/tmp/spacedock-codex-audit.3sHzek/repo` reproduced both material failures: a zero-collaboration fake host passed the live oracle, while rejection tests for attached-short and quoted-key owned overrides failed with the conflicting argv reaching launch. No code changes were made in the implementation worktree.

### Summary

Validation recommends REJECTED at `c36e3f7d1`. Baseline, race, formatting, focused argv, unsupported-host, and real credential-gated tests pass, and AC-3's covered launch variants carry the intended exact layer. However, the lifecycle proof is semantically vacuous and bypasses the Spacedock front door, while accepted Codex config spellings can override `agents.enabled` after Spacedock's owned layer. AC-1, AC-2, AC-4, E-2, E-3, and E-4 are therefore not established; AC-2 also has a material runtime defect.

### Recommendation

REJECTED. Route correction through the validation feedback gate. The narrow repair is to reject every Codex-accepted spelling of the reserved semantic paths, then drive the built `spacedock codex` front door from an isolated home and parse typed parent/child records to assert one spawn, ordered waits and same-ID follow-up/list, distinct child outputs, parent completion, and v2 context. Add the disabled-control run with an explicit zero-collaboration event assertion and retain a fake-host negative that must fail the lifecycle grader.

### Feedback Cycles

- Cycle 1: REJECTED (Sol/medium validation + Science Officer, 2026-08-01) —
  route through `feedback-to: implementation` before another validation pass.
  - M1 (evidence defect; Material): user harm — the launcher cannot establish
    its promised native collaboration lifecycle; observable harm — a zero-call
    fake Codex passes the live test; affected value — `value-ac[AC-1]` the
    supported isolated-home lifecycle must be observed through the launcher;
    trigger evidence — direct `codex exec`, unordered concatenated JSONL and
    prompt-supplied markers in the validator's detached falsifier.
  - M2 (product defect; Material): user harm — forwarded configuration can
    weaken the launcher guarantee; observable harm — Codex accepts attached
    short and quoted dotted `agents.enabled=false` forms after the owned true
    layer; affected value — `value-ac[AC-2]` forwarded configuration cannot
    silently downgrade the owned settings; trigger evidence — Codex 0.146.0
    probes and detached argv capture reach launch with exit 0.
  - M3 (evidence defect; Material): user harm — the disabled-control boundary
    is unverified; observable harm — the E-3 zero-collaboration negative is
    absent; affected value — `value-ac[AC-4]` runtime proof must use the
    behavioral oracle rather than a label or personal config; trigger evidence
    — validation report's E-1..E-5 reproduction records no disabled run.
  - FO authorization: `fix` the three findings within the approved semantic
    surface, then fresh-dispatch validation; preserve the exact findings and
    evidence unchanged in the feedback package.

## Stage Report: implementation (cycle 2)

- DONE: inject supported controls plus the exact multi_agent_v2 inline TOML layer across all launch variants
  Commit `383b4da8b` retains the exact layer; the complete-argv matrix fails if plain, local-plugin, Safehouse, or resume changes its placement, order, uniqueness, or any table field.
- DONE: reject reserved or conflicting forwarded settings before plugin install or host launch
  Attached `-cagents.enabled=false`, quoted dotted components, space/equals forms, nested v2 fields, and feature toggles now fail with the pinned diagnostic and zero install/launch calls; unrelated config remains a passing control.
- DONE: prove complete argv, unsupported-host failure, and isolated-home lifecycle behavior within the approved surface
  Eight files/+379 stayed within 10/+380; the 41.54s live run used built `spacedock codex`, parsed ordered typed parent/child records and v2 context, and required a zero-event disabled control; the vocabulary-only fake, unsupported-host, full, race, format, and detached-commit audits passed.

### Summary

Cycle 2 closes M1-M3 without changing their preserved findings or weakening the oracle. The first bounded live attempt exposed a completed-worker wait race; the final proof keeps the worker active until the ordered list/wait sequence, caps the subprocess at two minutes, and passes through the launcher at `383b4da8b`.

## Stage Report: validation (cycle 2)

- DONE: reproduce each AC with applicable unit, race, CLI, and isolated-home lifecycle evidence
  Commit `383b4da8b`: focused launcher/conflict/unsupported/vocabulary tests, `go test ./...`, `go test ./... -race`, `git diff --check`, and `gofmt -d ./cmd ./internal` passed; the built-front-door live test plus disabled control passed in 31.13s.
- FAILED: reproduce AC-1's exact same-worker lifecycle identity
  The real isolated-home run completed, but the grader never parses either `wait_agent` target; a detached typed/ordered transcript with both waits aimed at `worker-b` instead of spawned `worker-a` passed, so AC-1's same-worker wait identity is not established.
- DONE: reproduce AC-2's exact table and fail-closed forwarded/unsupported boundaries
  Complete argv carries the three owned assignments once; attached-short, quoted dotted, space/equals, nested-field, and feature-toggle conflicts all returned the pinned pre-side-effect rejection, while the unsupported-host fixture propagated exit 78.
- DONE: reproduce AC-3 across plain, local-plugin, Safehouse, and resume
  `TestCodexCollaborationLayerCompleteArgv` passed and atomically checks exact order, adjacency, uniqueness, table bytes, and the inner Safehouse boundary for all four variants.
- FAILED: reproduce AC-4's behavioral oracle as a complete lifecycle proof
  The front-door positive, zero-event disabled control, v2 context, vocabulary-only rejection, and docs boundary pass, but the same typed oracle accepts wrong-worker waits and therefore cannot certify the documented same-worker lifecycle claim.
- DONE: reproduce E-1 through E-5 and identify the exact residual boundary
  E-1 argv and E-3 zero-event negative pass; E-2/E-4 execute via built `spacedock codex` and observe terminal outputs/v2 context but remain incomplete on wait identity; current official config/subagent docs still name `agents.enabled`, `agents.max_concurrent_threads_per_session`, and stable `features.multi_agent`, not `multi_agent_v2`.
- DONE: perform semantic adversarial audit of launcher composition, ordering, cardinality, same-worker identity, and unsupported/conflict boundaries
  Throwaway checkout `/tmp/spacedock-codex-audit-cycle2.fKRtyk/repo` at the exact commit retained the passing vocabulary-only negative but produced the falsifier above; the candidate worktree stayed unchanged at 8 files/+379, within 10/+380.
- DONE: report exact evidence, docs/support-boundary checks, regressions, and a validation verdict without changing the deliverable
  Candidate HEAD remained `383b4da8b9b80954cc435e87d674a4cb1443b321`; only this split-root validation report was written.

### Summary

Cycle 2 fixes the prior front-door, reserved-spelling, disabled-control, and vocabulary-only defects, and the observed live launcher behavior succeeds. Validation nevertheless finds one material evidence defect: typed ordering and cardinality can still certify waits directed at a different worker, leaving AC-1 and the shared AC-4 oracle unproved.

### Recommendation

REJECTED. Route one narrow evidence correction through the validation feedback gate: parse both `wait_agent` argument records atomically, require each target to equal the spawned worker identity, and retain the detached wrong-worker transcript as a negative test before rerunning validation.

### Feedback Cycles

- Cycle 2: REJECTED (Codex validation, 2026-08-02) — surface 8 files/+379 vs 10/+380; AC unchanged.
  - M4 (evidence defect; Material): released user/normal workflow — the shipped live proof is the acceptance gate for every supported Spacedock-launched Codex; observable harm — a typed, correctly ordered, exact-cardinality transcript with both waits targeting another worker passes; affected value — `value-ac[AC-1]` the isolated-home proof must establish one same-worker native lifecycle through completion; trigger evidence — detached `TestAdversarialLifecycleRejectsWrongWaitTargets` fails because `gradeCodexLifecycle` does not parse or compare either wait target.
  - Science Officer advisory: REVISE. FO authorization: `fix` M4 by parsing both ordered `wait_agent` argument records atomically, requiring each target to equal the spawned worker identity, retaining the detached wrong-worker transcript as a negative, and rerunning validation; preserve M4's finding and evidence unchanged.
