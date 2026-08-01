---
id: z5gwwz2748sg6vxr0g3kdsar
title: "Codex launcher guarantees multi-agent v2 surface"
status: ideation
source: "Captain feedback, 2026-07-02: ordinary Codex config currently enables multi_agent_v2; make Spacedock-launched Codex enable or prove the same surface instead of relying on ambient session setup."
started: 2026-08-01T14:29:25Z
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
---

Spacedock's Codex front door should guarantee the observable collaboration lifecycle its first officer needs: spawn a worker, follow up with the same worker, inspect/list and wait for workers, and receive completion. The guarantee belongs to the launch configuration and behavioral proof, not to bootstrap prose, an ambient user config, or the undocumented `multi_agent_v2` label.

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

Conclusion: the label is neither necessary nor sufficient evidence. The two supported settings are the smallest launch-time control, and only the observable lifecycle is acceptance evidence.

## Proposed approach

1. `spacedock codex` injects `-c agents.enabled=true -c features.multi_agent=true` immediately after the inner `codex` token. This position works for fresh launches and `resume`, survives the Safehouse wrap unchanged, overrides user/profile defaults, and does not create or merge a temporary config file.
2. The launcher reserves those two config keys. A forwarded `-c`/`--config` assignment targeting either key, or `--enable`/`--disable multi_agent`, fails before plugin installation or `ops.Launch` with an actionable message explaining that Spacedock owns the collaboration guarantee. Redundant overrides are rejected too, keeping one effective source of truth and preventing last-override-wins ambiguity.
3. A Codex version which does not accept either supported setting rejects config before creating the interactive first-officer turn; the native parser error is preserved and the host exits nonzero. Documentation states this Codex support boundary. No prompt fallback or downgraded dispatch is permitted.
4. Tests use one shared expected-argv helper so existing byte-exact launch oracles continue to assert the entire command, not merely search for substrings. A small live probe separately proves that the supported settings produce the lifecycle; it does not modify PR #585 or its existing live harness.

This serves AC-1 and AC-2 with less machinery than a temporary `CODEX_HOME`/config merge, which would introduce file ownership, cleanup, auth, and profile-precedence hazards. Preserving the ambient home alone is insufficient because the clean-home case is the defect. A per-launch `codex doctor` preflight was rejected because it adds latency and reports the feature flag but not the complete collaboration lifecycle; explicit overrides plus a live lifecycle proof give a smaller and stronger boundary.

## Out of scope

- Changing Codex itself.
- Redesigning Spacedock's dispatch or reuse contract.
- Making multi-agent v1 behavior equivalent to v2.
- Depending on a particular developer's local config as the only proof.
- Enabling or asserting the `multi_agent_v2` feature-list label.
- Modifying PR #585 or its live harness.
- Changing bootstrap/first-officer prose to simulate capability detection.

## Expected surface

- `internal/cli/frontdoor.go`: 30-50 inserted lines for the owned config tokens, conflict recognition, and pre-launch diagnostic.
- `internal/cli/codex_multi_agent_test.go` (new): 120-180 inserted lines for the variant/conflict matrix and a live-test entry point guarded by existing live prerequisites.
- Existing exact-argv launcher tests in `internal/cli/frontdoor_test.go`, `internal/cli/safehouse_frontdoor_test.go`, `internal/cli/launch_parity_test.go`, `internal/cli/codex_plugin_dir_test.go`, and `internal/cli/safehouse_knob_test.go`: 20-50 net inserted lines through a shared expected-argv helper and updated complete argv expectations.
- `docs/site/reference/command-reference.md` and `docs/runtime-live-ci.md`: 8-20 inserted lines describing the launcher-owned settings, reserved forwarded keys, support boundary, and isolated-home lifecycle command.

Expected total: 8-9 files, 178-300 inserted lines, and at most 30 deleted lines. Tolerance is at most 10 files and 340 inserted lines. Exceeding either bound requires design re-entry rather than quietly broadening the mechanism.

Declared semantic changes are limited to Codex launch runtime behavior and two forwarded-config grammar reservations. Every `spacedock codex` invocation carries the supported multi-agent settings, including plain, local-plugin, Safehouse, and resume launches; conflicting attempts fail before plugin mutation or host launch. Command names, stored formats, workflow state, plugin format, dispatch/reuse authority, Claude/Pi behavior, and bootstrap text do not change.

## Acceptance criteria

**AC-1 - A supported Spacedock-launched Codex can complete the native collaboration lifecycle from an isolated home.**
Verified by: a live probe with no user config observes, in structured JSONL and child state, one spawn, one same-worker follow-up, list/inspect, wait, `CHILD_READY`, `CHILD_DONE`, and `PARENT_DONE`; it also records `multi_agent_version: v2`. Removing either launcher-owned supported setting or disabling `agents.enabled` makes this proof fail.

**AC-2 - The launcher cannot silently downgrade through ambient or forwarded configuration.**
Verified by: table tests start from ambient false settings and show the owned true overrides in complete argv, while every forwarded spelling targeting the reserved settings exits nonzero with the pinned diagnostic and records zero plugin installs and zero host launches. An instrumented unsupported host rejects the owned settings before opening a session and its nonzero exit/config diagnostic propagate; a live `agents.enabled=false` control produces no collaboration calls.

**AC-3 - Plain, local-plugin, Safehouse, and resume launches carry one identical supported enablement layer.**
Verified by: complete-argv launcher tests assert the two tokens occur exactly once, immediately after the inner `codex` token, across all four variants; deleting either token or placing one outside the Safehouse inner argv fails the matrix.

**AC-4 - Runtime proof and operator documentation use the supported settings and behavioral oracle, not the v2 label or personal config.**
Verified by: the isolated-home live test copies authentication only, applies the same two settings, and grades structured lifecycle events; documentation names the supported keys, the rejection behavior, and the minimum host-support boundary. A guard rejects `multi_agent_v2` in this task's test/helper surface.

## Test plan

Implementation starts with failing focused tests, then changes launcher behavior:

1. Add the complete-argv matrix for plain, explicit local-plugin, Safehouse, and exact `resume`. The matrix asserts adjacency, uniqueness, and the inner side of the wrapper boundary. Cost: low, fixture-backed, no model spend.
2. Add a conflict matrix for `-c`, `--config`, equals forms accepted by Codex, and `--enable`/`--disable multi_agent`. Each case asserts exit 1, byte-stable actionable stderr, no install, and no launch. Include nearby unrelated config as a green control. Cost: low, fixture-backed.
3. Add an instrumented unsupported-host case that rejects the two settings before starting a session. Assert the native actionable config diagnostic and nonzero exit propagate, and that its session/dispatch marker remains absent. Cost: low, process fixture.
4. Add the supported settings in `frontdoor.go`, before passthrough and after the inner program token, and reject reserved conflicts immediately after front-door parsing. Update all affected complete-argv oracles through one test helper. Cost: low-to-medium because exact argv is intentionally broad regression coverage.
5. Run the focused package tests, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.
6. Run the isolated-home live proof with copied auth only. Parse JSONL rather than final prose: require the ordered parent tool set, one stable child identity across follow-up, both child terminal messages, parent completion, and `multi_agent_version: v2`. Run the disabled control and require zero collaboration calls. Cost: one short positive and one short negative Codex run; live and credential-gated.
7. Repeat the instrumented frontdoor capture for plain, local-plugin, and Safehouse launch shapes and retain concise artifacts. The live lifecycle is run once because the same inner tokens are proven byte-identical across variants; three paid lifecycle runs would add cost without testing a distinct Codex mechanism.

## Consultation record

Standing Science Officer consultation (signed consultation supplied with the 2026-08-01 ideation dispatch): **REVISE**. Reference: current official Codex `Subagents` and `Configuration Reference` documentation plus the isolated Codex 0.146.0 probe. Reason to preserve in the gate: `agents.enabled` and stable `features.multi_agent` are the supported current controls; `multi_agent_v2=false` is only a local feature-label observation, so approval must rest on observable spawn, same-worker follow-up, list/wait, and completion across launcher variants, with a disabled/unsupported path that fails before worker dispatch.

## Stage Report: ideation

- DONE: Spike the real Codex frontdoor argv/config handoff from an isolated Codex home and record the falsifying result.
  Built launcher plus instrumented host captures proved no current override in plain/local/Safehouse argv; isolated 0.146.0 and live JSONL proved the stable-settings lifecycle while `agents.enabled=false` removed it.
- DONE: Declare the expected files/LOC/semantic surface and select the smallest supported enablement or fail-closed preflight.
  The body caps the change at 10 files/340 insertions and selects launcher-owned `agents.enabled=true` plus stable `features.multi_agent=true`, with conflicting forwarded controls rejected before install/launch.
- DONE: Add falsifiable launcher-test and live-proof plans covering plain, local-plugin, and sandbox launches.
  Complete-argv matrices fail on missing/duplicate/misplaced tokens; reserved-key cases require zero install/launch; isolated live JSONL requires spawn, same-worker follow-up, list/wait, child and parent completion.

### Summary

Reframed the task from enabling an unproven `multi_agent_v2` label to guaranteeing the observable Codex collaboration lifecycle through supported settings. The spike demonstrated that current Spacedock does not own that configuration, while a clean-home live run with the supported settings completed the full lifecycle and the disabled control did not expose collaboration; the signed Science Officer REVISE rationale is retained for the gate.
