# Independent ideation review: issue 813

Recommendation: **revise** for R1. The first officer classified R1 as owned Material and directed a bounded correction through the existing Safehouse allowlist.
The selected approach remains necessary. This review changes neither the design nor the five approved value criteria.

## Reviewed evidence

- Design and risk evidence: state commit `6b2c3ecc28276d06828ddd88774160c2a46d195d`, [task](index.md).
- Approved seed: state commit `3ae93ac36`, including its five value criteria.
- Product source: `661c607ae34620caa9579f2f65b53e4a5d242344`, especially the launcher, Safehouse wrapper, Pi extension, and existing proof owners.
- [Issue 813](https://github.com/spacedock-dev/spacedock/issues/813), the workflow README, repository instructions, and `docs/runtime-support.md`.
- The recorded native protocol probes and focused launcher regressions supplied existing evidence. This review did not repeat them.

## Necessity before soundness

| Mechanism | Removal challenge and result |
|---|---|
| Explicit mode | Protocol-token detection needs knowledge of three native grammars and cannot identify every caller-controlled session. One opt-in flag is sufficient. |
| Suppression of identity, prompts, permissions, and plugin work | Retention changes the caller's transaction or adds persistent setup. Local guards preserve the existing supervisor and wrapper. |
| Suppression of Codex collaboration injection and override rejection | Retention replaces caller settings and rejects valid caller choices. Native settings remain the caller's choice in this mode. |
| Suppression of Pi explicit resources and both inherited markers | Argv suppression alone leaves the installed extension's bootstrap and compaction gates active. Native resource discovery remains an explicit declared boundary. |
| Environment fallback with explicit-flag presence | Fixed provider argv cannot express session-specific capability choices. The shared helper avoids three precedence implementations. |
| Rejection of pre-fence tasks and plugin selection | Silent acceptance loses a task or triggers persistent setup. Native arguments already have the post-fence route. |
| Pi compatibility alias | AC-4 requires common acceptance. Preserving ordinary Pi readiness avoids an unrelated bypass mechanism. |

No separate dispatcher, protocol parser, process controller, dependency, or standing enforcement mechanism is necessary.
Automatic native-resource isolation is unnecessary because the design explicitly retains native discovery and caller flags.

## R1: wrapped launches lose explicit caller settings

**Classification:** Material, as directed by the first officer. **Ownership:** the task's design owner, then its implementation and validation owners.
The review initially proposed Needs decision because automatic forwarding and a caller obligation represented different boundaries.
The first officer resolved that ambiguity within the seed promise and authorized the bounded forwarding route.

- **Released user and normal workflow:** An orchestrator uses the provider command with Safehouse and an explicit environment selector for native settings.
  Codex supports `CODEX_HOME` for its settings and state. The proposed mode promises caller control in wrapped launches.
- **Observable harm:** The retained wrapper selected the wrong Codex model and changed the effective approval policy from `on-request` to `never`.
  The wrong policy came from an isolated fallback home. No model turn or shared host settings participated.
- **Affected authority:** `value-ac[AC-2]` The orchestrator retains its identity, prompting, and permission settings in wrapped passthrough launches.
- **Trigger evidence:** Real Codex 0.166.0 and Agent Safehouse 0.12.0 exercised the existing wrapper shape and native `config/read`.
  Only `SAFEHOUSE_ENV_PASS=CODEX_HOME` differed between the two observations.

### Independent observations

The fixture used an isolated `HOME`, an isolated working directory, and two literal settings files:

| File | Model | Approval policy |
|---|---|---|
| `<fixture>/home/.codex/config.toml` | `review-default` | `never` |
| `<fixture>/selected-codex/config.toml` | `review-selected` | `on-request` |

Both observations set `HOME=<fixture>/home` and `CODEX_HOME=<fixture>/selected-codex`.
Both used `safehouse --trust-workdir-config --add-dirs=<fixture> -- codex app-server` from the isolated working directory.
The fixture grant kept filesystem access identical. It is a probe input, not a proposed product grant.
Stdin remained open for `initialize`, `initialized`, and `config/read` with `{"includeLayers":false}`.
The observation recorded only the returned model and approval policy, then closed stdin.

| Observation | Effective model | Effective approval policy |
|---|---|---|
| Existing wrapper, no named allowance | `review-default` | `never` |
| Existing wrapper plus `SAFEHOUSE_ENV_PASS=CODEX_HOME` | `review-selected` | `on-request` |

A separate sentinel observation exercised the four promised environment values through the same real wrapper:

| Explicit parent value | Default wrapper child | Child with the existing named allowance |
|---|---|---|
| `CODEX_HOME=SPACEDOCK_REVIEW_CODEX_HOME` | absent | exact parent value |
| `CLAUDE_CONFIG_DIR=SPACEDOCK_REVIEW_CLAUDE_CONFIG_DIR` | absent | exact parent value |
| `PI_CODING_AGENT_DIR=SPACEDOCK_REVIEW_PI_AGENT_DIR` | absent | exact parent value |
| `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=0` | absent | `0` |

Both sentinel observations returned `APP_SANDBOX_CONTAINER_ID=agent-safehouse`.
The sentinel observation establishes environment transport. The Codex observation independently establishes a changed effective permission setting.
The observations do not establish effective Claude or Pi settings. Those hosts need the planned candidate evidence after implementation.

Safehouse documents its environment controls in [Environment modes](https://agent-safehouse.dev/docs/usage#environment-modes).
OpenAI documents the selector in [Config and state locations](https://learn.chatgpt.com/docs/config-file/config-advanced#config-and-state-locations).
OpenAI also documents the native [settings RPC](https://learn.chatgpt.com/docs/app-server#config-rpc-examples-for-app-settings).

### Authorized correction boundary

The first officer directed automatic preservation of the selected host's explicitly supplied canonical settings through the existing Safehouse allowlist:

- Claude: `CLAUDE_CONFIG_DIR` and the explicit `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` value, including `0`.
- Codex: `CODEX_HOME`.
- Pi: `PI_CODING_AGENT_DIR`.

Presence must select an allowance. The mode must preserve the supplied value without synthesizing absent values.
The correction adds no filesystem grant, general environment forwarding, caller setup obligation, or ordinary-mode behavior change.
The design owner must record this boundary and correct the proposed documentation before the gate.
The existing environment fixtures can own the deterministic proof. The planned candidate startup can include the Codex settings observation.
An edit that omits the selected selector's allowance must fail that proof even with byte-exact inner argv.

## Remaining proof and surface assessment

AC-1 has independent literal-argv and candidate protocol owners. Native-only startup remains baseline evidence until the candidate repeats it.
AC-2 needs R1 in those existing owners. Pi marker removal also closes the extension's bootstrap and compaction gates in the reviewed source.
AC-3 has a complete precedence table, including explicit false, capability replacement, profiles, missing Safehouse, and stable Pi paths.
AC-4 preserves ordinary behavior through existing regression owners. Its full-suite baseline remains unsuccessful, as the task records.
AC-5 proposes concrete documentation changes. R1 requires the caller-settings promise to match the corrected behavior.

The estimate is consistent: 470 insertions minus 70 deletions gives +400 net LOC across 13 product files.
The tolerance is +280 through +520 net LOC and at most 15 product files. It counts net growth, not gross churn.
The bounded correction fits the listed launcher helper, host functions, environment fixtures, and documentation responsibilities.
Implementation must count the actual result against that baseline.
The inherited detached audit and applicable live CI lanes remain required. This review proposes no additional standing check or live journey.

All three review checklist items are complete. R1 requires revision before approval, with the original value criteria unchanged.
