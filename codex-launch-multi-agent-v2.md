---
id: z5gwwz2748sg6vxr0g3kdsar
title: "Codex launcher guarantees multi-agent v2 surface"
status: implementation
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
    records:
        - id: gate:z5gwwz2748sg6vxr0g3kdsar:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-ideation-1
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:ideation:attempt-1:revision-1
                digest: sha256:988c0401c86e31c5e9a300df8dd7e3527e171d1b7dbbb56aa40262ffa6b0ec24
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
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-ideation-2
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:ideation:attempt-2:revision-1
                digest: sha256:0cc3451492227cd3b8b756842301674a4339ae6e6c1e42c9e45944980bd6c75b
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
                target-stage: implementation
                state: consumed
        - id: gate:z5gwwz2748sg6vxr0g3kdsar:validation
          stage: validation
          attempts:
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-1
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-1:revision-1
                digest: sha256:1065cc2ea5b9aca994efb80974074c6cf36c7c486c853f7a5e84031ac0252314
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
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-2
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-2:revision-1
                digest: sha256:ca917663cdcb933eae844888a3fe1661c3c6d4aef7aa133b1a0bc28c3873fed3
                request-digest: sha256:36f7b9dfe6811ad19739d9bbac0ac8e582884baece283ff45f5ec610ecbb3923
                room-ref: ./codex-launch-multi-agent-v2/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:z5gwwz2748sg6vxr0g3kdsar:validation:2
                briefing: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-01T17:57:10.685244Z"
                decision: revise
                reason: 'Captain-authorized direct send-back following Science Officer REVISE: fix M4 by parsing both ordered wait_agent argument records atomically and requiring each target to equal the spawned worker identity; retain the detached wrong-worker transcript as a negative and rerun validation. Prior M1-M3 corrections and product behavior remain accepted.'
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-3
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-3:revision-1
                digest: sha256:6e5d92302a28c91d4eb16f3c9e9a7d3a06d998db97e05a87545755f3ed859d31
                request-digest: sha256:26a53a4f9cdc8a4929c26bc574862f7b682e47c7bd0031ed199041ee03910175
                room-ref: ./codex-launch-multi-agent-v2/review/validation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:z5gwwz2748sg6vxr0g3kdsar:validation:3
                briefing: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-3:revision-1
                by: person:captain
                at: "2026-08-01T23:42:58.090537Z"
                decision: approve
                reason: Captain approved z5 after the cycle-3 gate review and Science Officer HOLD advisory; this approval accepts the current validation evidence and the bound implementation candidate.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-4
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-4:revision-1
                digest: sha256:09da751636b2cf4f71c3f5e4c4c56ef52a3ae295d95db472f51f4e562734556b
                request-digest: sha256:5517193c7a03a9d8a78eb927f27fa86333a6ffa7a9bf179198cefa3e43b5b1ef
                room-ref: ./codex-launch-multi-agent-v2/review/validation/briefing-4
              resolution:
                type: Resolution
                id: resolution:spacedock:z5gwwz2748sg6vxr0g3kdsar:validation:4
                briefing: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-4:revision-1
                by: person:captain
                at: "2026-08-02T01:45:32.207409Z"
                decision: approve
                reason: 'Captain approved validation cycle 4: AC-1 through AC-4 are independently evidenced, the exact multi_agent_v2 surface and M1-M4 behavior are green, the detached M5 spoof audit passes, and the 8-file/+379-26 candidate is within the approved scope.'
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-5
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-5:revision-1
                digest: sha256:d026b8eecbaa9a19a68bbb76a378c94e307e457f014e0ccab81ea76fe92fa736
                request-digest: sha256:1592a30ba68b38d177c7081463bd2f787012e21edceb73330cfa22a201da133c
                room-ref: ./codex-launch-multi-agent-v2/review/validation/briefing-5
              withdrawal:
                by: agent:first-officer
                at: "2026-08-02T09:17:35.468332Z"
                reason: 'Attempt 5 is stale: the candidate moved from dcab66435 on base 988163969 to 64f2f8773 on base 48a7ea0d9 after hq PR #598 merged. Cycle 7 now reports a reproducible recorded-gate-lifecycle failure; bind a fresh gate to that report.'
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-6
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-6:revision-1
                digest: sha256:c40931f5b3d340924fd670063b83071867b9466d645a52fef20f4a09db27f0ab
                request-digest: sha256:cc65d216d6aeba86493a8d1fa3a7a8693c7cef7bfad8d9f65b8606174f0aa427
                room-ref: ./codex-launch-multi-agent-v2/review/validation/briefing-6
              resolution:
                type: Resolution
                id: resolution:spacedock:z5gwwz2748sg6vxr0g3kdsar:validation:6
                briefing: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-6:revision-1
                by: agent:first-officer
                at: "2026-08-02T09:28:59.887424Z"
                decision: hold
                reason: 'Science classifies a material Codex First Officer runtime/state-commit outcome defect: the full and targeted recorded-gate-lifecycle runs have no structured worker.spawn event after consume, so the durable successor-dispatch predicate fails. Hold z5 unchanged; route the defect to local task codex-recorded-gate-successor-dispatch and require a fresh exact-head Codex lane.'
            - id: gate-attempt:z5gwwz2748sg6vxr0g3kdsar-validation-7
              briefing:
                id: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-7:revision-1
                digest: sha256:5c953491e929a08ac6809fa37921aeed7a557fbb5e792e33bbdcc492f6e058bf
                request-digest: sha256:0dbcaa3fee210617a69ab6291111df55886a6130fb616562adf28ce42ac70315
                room-ref: ./codex-launch-multi-agent-v2/review/validation/briefing-7
              resolution:
                type: Resolution
                id: resolution:spacedock:z5gwwz2748sg6vxr0g3kdsar:validation:7
                briefing: briefing:z5gwwz2748sg6vxr0g3kdsar:validation:attempt-7:revision-1
                by: person:captain
                at: "2026-08-05T00:18:50.95193Z"
                decision: approve
                reason: Captain approved the exact-head validation after reviewing the 9-file +380/-42 scope and current green evidence.
              application:
                target-stage: done
                state: superseded
mod-block:
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

## Stage Report: implementation (cycle 3)

- DONE: inject supported controls plus the exact multi_agent_v2 inline TOML layer across all launch variants
  Commit `1a8a64cdb` preserves the accepted M1-M3 launcher layer and complete-argv oracle unchanged; the final candidate remains 8 files/+379 against `main`, within 10/+380.
- DONE: reject reserved or conflicting forwarded settings before plugin install or host launch
  The final focused unsupported-host check and full/race suites pass while the accepted attached-short, quoted-key, nested-field, and feature-toggle rejection coverage remains unchanged.
- DONE: prove complete argv, unsupported-host failure, and isolated-home lifecycle behavior within the approved surface
  `TestAdversarialLifecycleRejectsWrongWaitTargets` rejects both explicit `worker-b` waits and targetless waits whose child mailbox/session identity is misbound; the 42.35s live built-front-door run, disabled control, `go test ./...`, `go test ./... -race`, `git diff --check`, and formatting all pass.

### Summary

Cycle 3 corrects only M4 in the lifecycle oracle. Both ordered wait argument records are parsed atomically; an explicit wait target must equal the spawned task, while Codex 0.146's native targetless waits are accepted only as idle observations because exact one-spawn cardinality and the child `parent_thread_id`/`agent_path`, ordered READY/DONE mailbox output, and parent completion independently bind the lifecycle to that spawned identity.

## Stage Report: validation (cycle 3)

- DONE: reproduce each AC with applicable unit, race, CLI, and isolated-home lifecycle evidence
  At `1a8a64cdb`, focused argv/conflict/unsupported/vocabulary/wrong-worker tests, `go test ./...`, `go test ./... -race`, `git diff --check`, and `gofmt -d ./cmd ./internal` passed; the real built-front-door lifecycle and disabled zero-event control passed in 52.81s.
- FAILED: reproduce AC-1 and AC-4's targetless-wait same-worker identity as typed evidence
  Detached `TestDetachedTargetlessWaitRejectsTypedMisboundChildWithSpoofedIdentityText` failed: an actual top-level `other-parent/worker-b` child passed when an unrelated JSON record carried expected `parent_thread_id/agent_path`, because the grader searches raw session bytes instead of parsing the child identity record atomically.
- DONE: reproduce AC-2's exact table and fail-closed downgrade boundaries
  Complete argv carried the owned stable controls and exact v2 table once; every covered forwarded reserved spelling failed before install/launch, the unrelated-config control passed, and the unsupported-host fixture propagated exit 78 with its native config diagnostic.
- DONE: reproduce AC-3 across plain, local-plugin, Safehouse, and resume
  `TestCodexCollaborationLayerCompleteArgv` atomically verified exact adjacency, order, uniqueness, table bytes, and the inner Safehouse boundary in all four variants.
- FAILED: reproduce the complete E-1 through E-5 ledger as acceptance-quality evidence
  E-1 argv, E-3 disabled zero-event control, E-4 real v2 context/completion, and E-5's supported-controls distinction were reproduced; E-2/E-4 cannot certify targetless same-worker binding because the child mailbox identity check accepts the detached misbound transcript above.
- DONE: perform semantic adversarial audit of launcher composition, ordering, cardinality, same-worker identity, and unsupported/conflict boundaries
  Throwaway checkout `/tmp/spacedock-codex-audit-cycle3.58z2ib/repo` at exact commit `1a8a64cdb` retained the explicit wrong-target and vocabulary negatives but exposed the targetless typed-identity falsifier without changing candidate bytes.
- DONE: report exact evidence, docs/support-boundary checks, regressions, and a validation verdict without changing the deliverable
  Candidate HEAD stayed clean at 8 files/+379 within 10/+380; docs correctly distinguish supported controls from the v2 label, but the claim that typed parent/child records grade a same-worker lifecycle still overstates the raw-substring identity check.

### Summary

The launcher behavior, full/race suites, live front-door lifecycle, disabled control, explicit wrong-target negative, and requested v2 context all pass. Validation nevertheless recommends REJECTED because the targetless-wait proof still does not parse the child mailbox identity atomically, so AC-1 and the shared AC-4 oracle remain falsifiable.

### Recommendation

REJECTED. M5 is an evidence defect and Material: released user/normal workflow — this oracle gates the supported launcher lifecycle; observable harm — a typed misbound child transcript passes; affected value — `value-ac[AC-1]` requires one same-worker isolated-home lifecycle; trigger evidence — the detached test above fails at `1a8a64cdb`. Cycle 3 therefore escalates under validation policy; do not mutate the candidate without a new disposition.

## Stage Report: implementation (cycle 4)

- DONE: Fix targetless-wait child identity binding by parsing the typed child record atomically and matching the spawned worker.
  Commit `155747c2fc` decodes each `session_meta.payload.source.subagent.thread_spawn` record and requires its `parent_thread_id` and `agent_path` together; changing either field makes the lifecycle grader reject the child.
- DONE: Retain a detached spoofed-identity negative and rerun focused lifecycle, full, race, format, and audit checks.
  `TestDetachedTargetlessWaitRejectsTypedMisboundChildWithSpoofedIdentityText` rejects an actual `other-parent/worker-b` record despite unrelated expected identity text; focused tests, `go test ./...`, `go test ./... -race`, `git diff --check`, and `gofmt -d ./cmd ./internal` pass, as does the 47.13s live lifecycle plus disabled control.
- DONE: Preserve the exact multi_agent_v2 launch/config surface and accepted M1-M4 behavior within the existing file/LOC boundary.
  The complete-argv, reserved-conflict, unsupported-host, vocabulary, and explicit wrong-wait negatives pass; the candidate remains exactly 8 files/+379 against `main`, and the correction changes only `internal/cli/codex_multi_agent_test.go`.

### Summary

Cycle 4 corrects only M5 by binding the child identity fields from one typed session metadata record instead of independent raw-session substrings. The real isolated-home lifecycle, disabled control, full/race suites, formatting, and exact approved surface all pass at `155747c2fc`.

## Stage Report: validation (cycle 4)

- DONE: Independently verify AC-1 through AC-4 against the cycle-4 implementation and current captain scope.
  At `155747c2fc`, AC-1/AC-4 passed the real built-front-door isolated-home lifecycle plus disabled zero-event control; AC-2 passed exact-layer, conflict, and unsupported-host tests; AC-3 passed complete argv for plain, local-plugin, Safehouse, and resume.
- DONE: Run focused lifecycle, full, race, format, detached adversarial, and audit checks; record exact green/red evidence.
  Focused lifecycle/config tests passed; `go test ./...` and `go test ./... -race` exited 0; `gofmt -d ./cmd ./internal`, `git diff --check`, and the candidate worktree were clean. No requested check was red.
- DONE: Confirm the isolated-home launch contains the exact multi_agent_v2 TOML and no M1-M4 regression.
  The live test passed in 49.99s through the built `spacedock codex` front door with the exact 16-thread/`agents`/visible-metadata table, ordered typed lifecycle and v2 context; reserved-spelling, vocabulary-only, wrong-wait, unsupported-host, and disabled-control proofs retained M1-M4.
- DONE: Perform the required detached adversarial audit and report any material/evidence finding with all four fields.
  Throwaway checkout `/tmp/spacedock-codex-audit-cycle4.1V41Lh/repo` at exact commit `155747c2fc7bc0f4e8d7538db126521e164e364a` rejected the typed misbound child despite spoofed expected identity text; no new material or evidence finding arose, so no four-field disposition entry is required.

### Summary

Cycle-4 validation independently reproduces AC-1 through AC-4 and recommends PASSED. The candidate remains exactly 8 files, +379/-26 against `origin/main`, preserves the approved launcher-owned multi-agent surface and M1-M4 behavior, and closes M5 by requiring the parent and worker identity fields to coexist in one typed child session record.

### Recommendation

PASSED. All promised ACs have behavioral evidence, the required live and detached adversarial proofs are green, and no material finding or deferred risk remains.

## Stage Report: implementation (cycle 5)

- DONE: Remove the obsolete forwarded --enable multi_agent_v2 pair from the shared Codex live runner and keep launcher-owned config as the sole source.
  Commit `60a75ba658` removes the reserved pair and rewrites the runner assertion so either forwarded `--enable` or `--disable` makes the build-tag check fail.
- DONE: Correct docs/runtime-live-ci.md to describe launcher-owned enablement and retain the exact v2 table/support boundary.
  The launch row now names both stable controls, the exact 16-thread/`agents`/visible-metadata table, forwarded-override rejection, and fail-closed host support; the setup paragraph retains fence, bypass, isolated-home, final-message, and lifecycle-oracle details.
- DONE: Reproduce the shared live Codex smoke through the Spacedock front door, then rerun focused, full, race, format, diff, and detached checks.
  `TestLiveCodexSharedScenarios` passed in 765.06s (seven enabled journeys passed; three existing TODO journeys skipped); focused CLI/live-tag tests, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` passed, as did the detached spoofed-identity test at exact commit `60a75ba658`.

### Review Finding Disposition

- Released user and normal workflow: the required shared Codex live lane invokes each supported journey through `spacedock codex` with its isolated home and local checkout.
- Observable harm: the runner forwarded a launcher-reserved `--enable multi_agent_v2` pair, so the shared lane exited 1 before host launch with `collaboration settings are managed by Spacedock; remove the forwarded override`.
- Affected authority: `value-ac[AC-2]` — launcher-owned collaboration settings cannot be silently downgraded or shadowed by forwarded configuration.
- Trigger evidence: pre-correction `internal/ensigncycle/codex_live_runner_test.go:454` appended the reserved pair; `go run ./cmd/spacedock codex -- --enable multi_agent_v2` reproduced the pinned rejection before launch.
- Classification, ownership, disposition: Material integration/test defect owned by the current candidate; the First Officer authorized fixing runner/docs only while preserving the fail-closed product guard.
- Exact scope delta: final candidate is 9 files, +380/-42 against `main`; file and insertion tolerance remain satisfied. Deletions exceed the ideation estimate of 30 by 12 because the obsolete live-runner assertion and stale 11-line operator narrative were replaced by a compact launcher-owned source-of-truth description; product semantics are unchanged.

### Summary

Cycle 5 restores the required shared Codex live lane without weakening launcher ownership or its reserved-override rejection. The real front-door suite, offline/race suites, formatting, diff, and detached adversarial evidence all pass at `60a75ba658`, with the correction held to 9 files and exactly +380 lines overall.

## Stage Report: validation (cycle 5)

- DONE: Independently verify AC-1 through AC-4 against the corrected launcher-owned collaboration surface and the authorized runner/docs disposition.
  At exact commit `60a75ba658`, AC-1/AC-4 passed the real built-front-door isolated-home lifecycle in 39.91s with ordered typed same-worker events, v2 context, and a disabled zero-event control; AC-2/AC-3 passed exact-layer, reserved-conflict, unsupported-host, and four-variant complete-argv tests. The refreshed official Codex manual documents `agents.enabled`, `agents.max_concurrent_threads_per_session`, and stable `features.multi_agent`, but not `multi_agent_v2` or its three private fields, matching the docs' supported-boundary distinction.
- DONE: Reproduce the build-tag shared Codex live lane, focused/full/race/format/diff/detached checks, and record exact green or skipped evidence.
  `TestLiveCodexSharedScenarios` passed in 792.62s: seven enabled journeys passed (gate-guardrail, recorded-gate-lifecycle, feedback-3-cycle-escalation, merge-hook-guardrail, filing, shallow-boot, self-evidence-merge-triage) and three existing TODO journeys skipped rather than counting as evidence. Focused CLI/live-runner tests, `go test ./...`, `go test ./... -race`, zero-byte `gofmt -d ./cmd ./internal`, and both `git diff --check` forms exited 0.
- DONE: Confirm the corrected candidate is 9 files, +380/-42 against main, preserves the fail-closed reserved-override guard, and has no unrelated scope expansion.
  `git diff --numstat main...HEAD` reports exactly 9 files/+380/-42; cycle 5 changes only `internal/ensigncycle/codex_live_runner_test.go` and `docs/runtime-live-ci.md`. Detached checkout `/tmp/spacedock-codex-audit-cycle5.OYNdOw/repo` at the exact commit passed reserved-spelling and lifecycle-spoof negatives, runner/front-door tests, and rejected the removed `--enable multi_agent_v2` form with exit 1 and the pinned pre-launch diagnostic.
- SKIPPED: Treat the three existing shared live TODO journeys as passes.
  `rejection-flow` remains `TODO(zbcj98qfwtax61vxdzrf615e)`; `smallest-sufficient-mechanism` and `keep-moving-posture` remain `TODO(9adv48yhye5s2vkhwd7ge52d)`. Their zero-second skips provide no capability evidence.
- DONE: Verify the deletion excess over the ideation narrative is explicitly justified and does not change product semantics.
  The 12 deletions beyond the estimate replace the obsolete forwarded-enable assertion and stale 11-line operator narrative with a compact launcher-owned source-of-truth description; the runner still preserves the front-door fence, bypass posture, isolated home, final-message path, and all seven enabled journeys.

### Summary

Validation cycle 5 recommends PASSED at `60a75ba658`. The shared live lane now reaches Codex through Spacedock without forwarding a reserved override, launcher ownership remains fail-closed, all requested offline/race/live/adversarial evidence is green, and no material finding, deferred risk, unrelated runner edit, compatibility fallback, or scope expansion remains.

### Recommendation

PASSED. The candidate satisfies AC-1 through AC-4 and the authorized 9-file/+380 tolerance; the three named TODO skips remain explicitly non-evidence.

## Stage Report: validation (cycle 6, stale base)

- DONE: Validate the rebased Codex multi-agent-v2 candidate at exact head dcab66435 against origin/main 988163969, preserving the approved nine-file +380/-42 scope and five-commit patch equivalence.
  Before the base advanced, `git range-diff` mapped all five pre-rebase commits exactly (`=`), and `git diff --numstat 988163969...dcab66435` reported exactly 9 files, +380/-42 with a clean candidate worktree.
- DONE: Run focused multi-agent-v2 launcher/config/lifecycle tests, full and race suites, format/diff checks, detached evidence, and the required Codex live lane; triage only this rewritten head.
  Focused launcher/config/oracle and live-runner boundary tests passed; `go test ./...`, `go test ./... -race`, `gofmt -d ./cmd ./internal`, and both diff checks exited 0; the isolated-home lifecycle plus disabled control passed in 52.24s.
- DONE: Perform the detached adversarial audit at the exact rewritten head.
  Detached checkout `/tmp/spacedock-codex-audit-rebase.cJTlOM/repo` at `dcab66435` passed complete-argv, reserved-conflict, unsupported-host, vocabulary-only, typed spoofed-identity, and live-runner boundary tests; the removed forwarded enable form exited 1 with the pinned pre-launch diagnostic.
- FAILED: Complete the required shared Codex live lane as final evidence.
  Four journeys completed (`gate-guardrail`, `recorded-gate-lifecycle`, `feedback-3-cycle-escalation`, `merge-hook-guardrail`) and `rejection-flow` retained its known TODO skip; the run was intentionally interrupted at the next boundary after `origin/main` advanced, so `filing` and later journeys are not evidence.
- DONE: Reproduce the AC-1 through AC-4 evidence against the dispatched candidate before staleness was reported.
  AC-1/AC-4 passed the isolated-home structured lifecycle and disabled zero-event negative; AC-2 passed exact layer, conflict, unsupported-host, and fail-closed runner checks; AC-3 passed byte-exact plain/local-plugin/Safehouse/resume argv identity.
- SKIPPED: Record a fresh final validation recommendation for the new base.
  `origin/main` advanced from dispatched base `988163969` to `48a7ea0d9` while the shared live lane was running; the First Officer requested a stop and fresh rebase, so this stale cycle cannot authorize PASSED or REJECTED for the successor candidate.

### Summary

The exact dispatched rewrite was green across focused, full, race, format, diff, detached, and short live evidence, with no new finding observed. The shared live lane was partially green but deliberately stopped after the base advanced; all evidence in this report is stale for gate authority and must be rerun on the next rebased head.

### Recommendation

WITHHELD — rebase onto `48a7ea0d9` or later and perform one fresh validation cycle.

## Stage Report: validation (cycle 7, final base)

- DONE: Validate the exact candidate and approved scope against the dispatched final base.
  `HEAD` remained `64f2f8773a3ba4d848c15f7c9f71bbc45e55c395` and `origin/main` remained `48a7ea0d97042f0e7aaac258e1b77f16157c5281`; the worktree was clean, `git diff --numstat` reported exactly 9 files/+380/-42, and `git range-diff 988163969..dcab66435 48a7ea0d9..64f2f8773` mapped all five commits with `=`.
- DONE: Reproduce focused launcher/config/lifecycle evidence, full and race suites, formatting/diff integrity, and detached adversarial evidence.
  Complete argv for plain/local-plugin/Safehouse/resume, every reserved forwarded spelling, unsupported-host propagation, vocabulary-only rejection, typed spoofed-child rejection, live-runner front-door boundaries, `go test ./...`, `go test ./... -race`, zero-byte `gofmt -d ./cmd ./internal`, and both diff checks passed. Detached checkout `/tmp/spacedock-codex-audit-final.QGQOmW/repo` at exact head retained those negatives and the removed `--enable multi_agent_v2` form produced the pinned exit-1 pre-launch diagnostic.
- DONE: Reproduce AC-1 and AC-4's isolated-home behavioral evidence on the final head.
  AC-1 and AC-4 evidence: the built-front-door lifecycle plus disabled zero-event control passed in 37.87s; the typed oracle observed ordered spawn/wait/follow-up/list/wait, same-worker binding, child and parent terminal messages, and `multi_agent_version: v2` without using the feature-list label or personal configuration as proof.
- DONE: Reproduce AC-2's exact table and fail-closed downgrade boundaries on the final head.
  AC-2 evidence: `TestCodexCollaborationLayerCompleteArgv`, `TestCodexRejectsOwnedCollaborationOverridesBeforeSideEffects`, and `TestCodexUnsupportedHostFailsBeforeSession` passed; the detached exact-head audit also rejected the removed forwarded `--enable multi_agent_v2` form with the pinned exit-1 pre-launch diagnostic.
- DONE: Reproduce AC-3's identical launcher-owned layer across all required launch variants on the final head.
  AC-3 evidence: `TestCodexCollaborationLayerCompleteArgv` atomically passed exact adjacency, order, uniqueness, and table bytes for plain, explicit local-plugin, Safehouse-inner, and resume argv.
- FAILED: Complete the required shared Codex live lane as green final evidence.
  `TestLiveCodexSharedScenarios` failed after 745.38s: six enabled journeys passed (`gate-guardrail`, `feedback-3-cycle-escalation`, `merge-hook-guardrail`, `filing`, `shallow-boot`, and `self-evidence-merge-triage`), the three existing TODO journeys skipped, and `recorded-gate-lifecycle` failed because the durable oracle reported `successor dispatch was not observed after consume`. A targeted retained-artifact rerun reproduced the same failure after 241.33s. Its JSONL shows prepare, bind commit, decision record, close commit, consume, consumed commit, a successful dispatch build, narration that one worker was being dispatched, no observed structured `worker.spawn`/spawn collaboration event, one structured `wait` event with an empty receiver set and agent-state map, and the committed marker; the shipped ordering/ancestry predicate still rejects the run, and neither narration nor the model's final-message claim is dispatch evidence.
- DONE: Preserve the reproducible live failure artifacts and classify the gate impact without changing candidate code.
  Retained artifacts are under `/tmp/spacedock-codex-final-live-artifacts/codex-shared-scenarios/recorded-gate-lifecycle`. Material outcome classification — released user/normal workflow: a Codex first officer driving a consumed recorded gate into successor handoff; observable harm: it claims durable successor completion after dispatch-build narration and a targetless wait, but no structured spawn event establishes dispatch and the required durable oracle cannot establish post-consume dispatch ordering/ancestry; affected authority: the mandatory shared Codex live lane required by this validation; trigger evidence: the full lane and one isolated rerun fail the same `assertRecordedGateLifecycle` predicate at exact head `64f2f8773`.

### Summary

Final-base validation independently establishes the launcher-owned multi-agent-v2 surface and AC-1 through AC-4 at the exact nine-file candidate. The authoritative report extraction now counts 6 DONE, 0 SKIPPED, and 1 FAILED checklist items, and the AC scan reports two evidence citations for each of AC-1 through AC-4 with none unevidenced. Approval is nevertheless blocked because the required shared Codex lane is reproducibly red on `recorded-gate-lifecycle`; six enabled journeys pass and three known TODOs remain non-evidence, but a claiming final message cannot override the failed durable oracle.

### Recommendation

REJECTED. Route the reproducible `recorded-gate-lifecycle` failure through the validation feedback gate, preserving the retained JSONL and exact predicate failure; do not treat the green isolated collaboration lifecycle or six other shared journeys as a substitute for the mandatory red lane.

## Stage Report: validation (cycle 8, post-PR #582)

- DONE: Run focused multi-agent-v2 launcher/config/lifecycle checks, full and race suites, format/diff checks, detached audit, and the Codex/shared-live lanes on exact HEAD 37abaa8b3.
  Focused argv/config/oracle and live-runner boundary tests, `go test ./...`, `go test ./... -race`, zero-byte `gofmt -d ./cmd ./internal`, and both diff checks passed; detached checkout `/tmp/spacedock-codex-audit-z5.HrdUPq/repo` passed the same high-risk negatives and a built front door returned exit 1 with the exact reserved-override diagnostic. The isolated-home live test passed in 52.42s, and `TestLiveCodexSharedScenarios` passed in 671.25s with all seven enabled journeys green and the three existing TODO journeys explicitly skipped.
- DONE: Verify AC-1 through AC-4 with independent front-door behavior and exact argv/state evidence; do not reuse cycle-7 evidence from the old base.
  AC-1: the fresh built-front-door isolated-home run observed typed ordered spawn/wait/follow-up/list/wait events, one stable child identity, both child terminal messages, parent completion, `multi_agent_version: v2`, and a disabled zero-event control. AC-2: exact-table, reserved-spelling, unrelated-control, unsupported-host, and built-binary pre-launch rejection checks passed. AC-3: complete argv passed exact adjacency, order, uniqueness, and bytes across plain, local-plugin, Safehouse-inner, and resume. AC-4: the fresh behavioral oracle passed while the same current host still reported `multi_agent_v2 stable false`; vocabulary-only and spoofed-identity negatives remained red.
- DONE: Perform the semantic adversarial matrix, classify the recorded-gate lifecycle finding, and recommend PASSED only when the required live evidence is current and green.
  The matrix covered separate/attached/equals/quoted config spellings, every reserved root and nested v2 path, unrelated config, four launch routes, unsupported host, disabled collaboration, vocabulary-only output, spoofed child identity, and targetless wait; the launcher hot path remains one bounded argv scan with no blocking I/O or implicit size allocation. The prior recorded-gate finding was a Material outcome failure on the superseded head, but it is not present on this candidate: the fresh durable oracle passed `recorded-gate-lifecycle` in 226.21s after prepare/bind/record/close/consume and observed successor handoff, so no current material or deferred finding remains.

### Summary

Fresh validation recommends PASSED at exact HEAD `37abaa8b35e7044d24dc06e009c047d875558bff`. The five-commit patch maps unchanged by `git range-diff`, scope remains exactly 9 files/+380/-42 against dispatched merge base `db7f1e84aef5df2daf20fb02deac440df4ae1af1`, and the worktree is clean; `origin/main` is now `655aa679f588523388b808d8686809522ab6e8ec`, a known distinction reserved for later PR/mergeability reconciliation rather than this commissioned exact-head validation.

### Recommendation

PASSED. AC-1 through AC-4 and the mandatory shared-live lane are current and green; the three named TODO skips remain non-evidence.

## Stage Report: implementation (cycle 6)

- DONE: Commit a minimal current-main reconciliation that preserves Z5 Codex multi-agent launch behavior.
  Merge commit `b0b9ef319f2d267c7da3a7df70dee83b6446ba86` incorporates `origin/main` `be0e8453e6075960075e3a34d033d822f2009f4c`; the clean candidate remains exactly 9 files/+380/-42 against that main head.
- DONE: Resolve only the proven codex_live_runner_test.go content conflict without expanding runtime, scheduler, or authority semantics.
  Current main relocated the argv helper into `team_capability_test.go`; the resolution retains main's quiet-budget runner and migrates Z5's removal/rejection assertion there, so reintroducing forwarded `--enable` or `--disable multi_agent_v2` fails the offline boundary test.
- FAILED: Run focused checks plus gofmt, full, and race suites; record the exact new candidate head and residual merge status.
  Focused CLI, offline runner, and live-tag runner checks passed; `gofmt -w ./cmd ./internal`, zero-byte `gofmt -d`, and diff checks passed. Both `go test ./...` and `go test ./... -race` passed every package except unrelated `internal/gates/TestV1PilotManifestReadsAndValidates`, whose current-main manifest names active `gate-agent-ergonomics.md` while the shared state checkout contains `_archive/gate-agent-ergonomics.md`; no Z5 or unrelated state bytes were changed. Exact head is `b0b9ef319f2d267c7da3a7df70dee83b6446ba86`; `origin/main` is its ancestor, `git merge-tree --write-tree HEAD origin/main` exits 0, and the code worktree is clean.

### Summary

Z5 is reconciled onto current main with its launcher-owned multi-agent-v2 semantics intact and no product-scope expansion. The candidate is clean and mergeable at `b0b9ef319f2d267c7da3a7df70dee83b6446ba86`; only the pre-existing shared-state manifest drift prevents repository-wide full/race commands from being entirely green.
