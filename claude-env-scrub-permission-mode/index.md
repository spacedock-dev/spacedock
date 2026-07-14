---
id: d1sah62r0xckeysjet5j6zzk
title: Preserve Claude first-officer permissions with subprocess environment scrubbing
status: ideation
source: GitHub issue spacedock-dev/spacedock#504, reported by jesserobbins 2026-07-14
started: 2026-07-14T04:12:50Z
completed:
verdict:
score:
worktree:
issue: spacedock-dev/spacedock#504
---

`CLAUDE_CODE_SUBPROCESS_ENV_SCRUB=1` must compose with `spacedock claude` without silently changing a first officer from Claude Code's `auto` mode to `default`.

## Problem

[Issue #504](https://github.com/spacedock-dev/spacedock/issues/504) records this external observation on Spacedock 0.24.0: an unsandboxed `spacedock claude` launch showed the normal `Sandbox: available, not enabled` banner, then Claude Code warned that subprocess environment scrubbing forced the session to `default` because the parent had not declared `allowedTools`. Every consequential tool call then required approval.

This is a composition defect, not user error. The unsandboxed lane intentionally starts Claude with `--permission-mode auto`; Safehouse is optional. Auto mode applies Claude's model-classified permission policy. Safehouse supplies separate OS-level filesystem and network isolation. The reporter chose the supported policy-only lane and also enabled credential scrubbing.

Claude Code 2.1.208 exposes both `--permission-mode auto` and `--allowedTools`. Its [CLI reference](https://code.claude.com/docs/en/cli-usage) says `--allowedTools` pre-approves listed tools; it does not restrict unlisted tools. Its [permissions reference](https://code.claude.com/docs/en/permissions) distinguishes permission policy from sandbox enforcement, and its [environment reference](https://code.claude.com/docs/en/env-vars) says scrub mode removes Anthropic and cloud credentials from Bash, hook, and MCP child environments. Therefore Spacedock must not declare broad `Bash`, `Edit`, or `Write` access merely to silence the warning.

## Mechanism spike: external terminal required

The first captain-run probe is a negative control, not permission evidence. `/permissions` dismissed without a durable mode record; Claude correctly refused a request that would print `AWS_ACCESS_KEY_ID` into the transcript; and the launch prompt told the first officer to wait instead of running its normal boot. That transcript proves the old exercise was invalid. It says nothing about whether the candidate allow declaration preserves auto mode.

Use two independent, machine-written observations instead:

- Claude's `system/init` stream event records `permissionMode`; no dialog or model interpretation is involved.
- A `SessionStart` hook records only whether two fabricated AWS credential variables are present in the hook's child environment. It never reads or prints either value. Claude's environment reference explicitly includes hooks in subprocess scrubbing.

Run this probe from a normal macOS Terminal. It creates disposable project-local settings and does not alter `~/.claude/settings.json`:

```bash
PROBE="$(mktemp -d /tmp/spacedock-504.XXXXXX)"
BIN=/opt/homebrew/Caskroom/spacedock@next/0.25.0-pre2/spacedock
for lane in off on candidate safehouse; do
  mkdir -p "$PROBE/$lane/.claude"
  cat > "$PROBE/$lane/probe-scrub.sh" <<'SH'
#!/bin/sh
set -eu
if [ "${AWS_ACCESS_KEY_ID+x}" = x ] || [ "${AWS_SECRET_ACCESS_KEY+x}" = x ]; then
  state=visible
else
  state=scrubbed
fi
(umask 077; printf 'credential_vars=%s\n' "$state" > "$CLAUDE_PROJECT_DIR/scrub-state.txt")
SH
  chmod +x "$PROBE/$lane/probe-scrub.sh"
done
for spec in off:0 on:1 candidate:1 safehouse:1; do
  lane="${spec%%:*}"
  scrub="${spec##*:}"
  printf '%s\n' "{\"env\":{\"CLAUDE_CODE_SUBPROCESS_ENV_SCRUB\":\"$scrub\",\"AWS_ACCESS_KEY_ID\":\"FABRICATED_ID\",\"AWS_SECRET_ACCESS_KEY\":\"FABRICATED_SECRET\"},\"hooks\":{\"SessionStart\":[{\"matcher\":\"startup\",\"hooks\":[{\"type\":\"command\",\"command\":\"\${CLAUDE_PROJECT_DIR}/probe-scrub.sh\"}]}]}}" > "$PROBE/$lane/.claude/settings.local.json"
done
"$BIN" doctor --host claude
```

Run the normal headless first-officer bootstrap in all four lanes. Do not add a task prompt; Spacedock's ordinary bootstrap prompt must run.

```bash
(cd "$PROBE/off" && "$BIN" claude -- -p --output-format stream-json --verbose --model haiku > stream.jsonl 2> banner.txt); echo $? > "$PROBE/off/exit.txt"
(cd "$PROBE/on" && "$BIN" claude -- -p --output-format stream-json --verbose --model haiku > stream.jsonl 2> banner.txt); echo $? > "$PROBE/on/exit.txt"
(cd "$PROBE/candidate" && "$BIN" claude -- --allowedTools Read -p --output-format stream-json --verbose --model haiku > stream.jsonl 2> banner.txt); echo $? > "$PROBE/candidate/exit.txt"
(cd "$PROBE/safehouse" && "$BIN" claude --safehouse -- -p --output-format stream-json --verbose --model haiku > stream.jsonl 2> banner.txt); echo $? > "$PROBE/safehouse/exit.txt"
```

Extract only mode, boolean scrub state, exit, and relevant banner lines:

```bash
for lane in off on candidate safehouse; do
  printf '%s ' "$lane"
  jq -r 'select(.type=="system" and .subtype=="init") | "mode=" + .permissionMode' "$PROBE/$lane/stream.jsonl" | head -1
  cat "$PROBE/$lane/scrub-state.txt" "$PROBE/$lane/exit.txt"
  grep -E '^(Sandbox:|Use --safehouse to enable\.|.*Permission mode forced)' "$PROBE/$lane/banner.txt" || true
done
```

The expected matrix is independent of model behavior:

| lane | `system/init.permissionMode` | hook result | sandbox banner |
| --- | --- | --- | --- |
| scrub off, unsandboxed | `auto` | `credential_vars=visible` | available, not enabled |
| scrub on, unsandboxed baseline | `default` plus upstream warning | `credential_vars=scrubbed` | available, not enabled |
| scrub on, `--allowedTools Read` | `auto`, no forced-default warning | `credential_vars=scrubbed` | available, not enabled |
| scrub on, `--safehouse` | `bypassPermissions` | `credential_vars=scrubbed` | enabled, unchanged |

The candidate mechanism passes only when the third row matches exactly. The preimplementation binary does not yet print the proposed opt-in hint; AC-4 covers that later banner change. If the third row fails, retain the artifact, reject the allow-list approach, and require Safehouse for scrubbed autonomous operation. This Codex sandbox cannot replace the external unsandboxed rows.

## Proposed approach

Preserve scrub-on unsandboxed auto as the supported target. On a fresh, unsandboxed Claude launch where Spacedock injects `--permission-mode auto`, also inject the smallest externally proved explicit declaration: initially `--allowedTools Read`. `Read` is already a no-prompt, read-only built-in, so the declaration must not pre-approve Bash or writes. Gate the exact token on the external probe above; if `Read` does not preserve auto, do not widen the list by trial and error. Reframe that host combination as unsupported and require Safehouse for autonomous operation.

Apply the declaration at the existing `!wrap && !resume && no operator permission mode` branch in `internal/cli/frontdoor.go`. Preserve an operator's `--allowedTools` or `--allowed-tools` list without adding a duplicate. Preserve operator permission modes, resumes, and the Safehouse argv unchanged. Never clear, override, or recommend disabling `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB`.

Extend Claude doctor capability checking through the existing host-operation seam. A host whose help surface includes the proved permission mode and allow-list flags reports:

```text
OK: Claude Code supports scrub-compatible auto mode (--allowedTools).
```

A host that lacks the proved surface exits nonzero after the normal plugin report and prints:

```text
UNSUPPORTED: Claude Code cannot preserve scrubbed unsandboxed auto mode; update Claude Code or run `spacedock claude --safehouse`. Credential scrubbing remains enabled.
```

Doctor checks the host's command surface, not managed-policy state. Claude remains authoritative for a policy that disables auto or rejects a launch; Spacedock must preserve that host error.

Add the Safehouse opt-in hint to the shared fresh-launch banner. Print these adjacent lines exactly when all three conditions hold: platform is macOS, `safehouse` resolves on `PATH`, and this launch did not select Safehouse through a `.safehouse` profile, `--safehouse`, or another `--safehouse-*` flag.

```text
Sandbox: available, not enabled (no .safehouse profile)
Use --safehouse to enable.
```

Pass the platform into a small banner-rendering helper so Linux CI can cover both macOS and non-macOS cases. Fresh Claude, Codex, and Pi launches share the hint because they share `launchBanner`. Resumes remain banner-free. Unavailable Safehouse prints no hint. An enabled Safehouse launch keeps its existing line byte-for-byte and prints no hint.

## Documentation change

In `docs/site/reference/command-reference.md`, replace the unsandboxed-launch paragraph with:

> An unsandboxed bootstrap launch carries no Safehouse isolation. `spacedock claude` starts in `--permission-mode auto`, and Codex starts in `--ask-for-approval on-request`, unless you supply a mode. When Claude subprocess credential scrubbing is enabled, Spacedock adds a minimal read-only declaration so auto mode and credential scrubbing compose; it does not pre-approve Bash or writes. A sandboxed bootstrap instead bypasses prompts inside Safehouse's OS boundary. Run `spacedock doctor --host claude` if Claude rejects the scrub-compatible surface. Spacedock never disables credential scrubbing.

Add immediately after that paragraph:

> On macOS, when Safehouse is installed but the current launch has no `.safehouse` profile or Safehouse flag, the launch banner says `Use --safehouse to enable.` The hint advertises the existing opt-in; Spacedock does not enable the sandbox automatically.

Expand `docs/site/reference/sandbox.md` after the table:

> Permission mode is policy, not process isolation. Claude `auto` classifies tool calls but does not confine the process. Safehouse enforces filesystem and network boundaries around the host. Both launch lanes are supported; choose Safehouse when you need OS-level containment.

Then add the exact unselected macOS banner:

```text
Sandbox: available, not enabled (no .safehouse profile)
Use --safehouse to enable.
```

## Out of scope

- Disabling or clearing `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB`.
- Broadly pre-approving Bash, file mutation, network access, or MCP tools.
- Automatically enabling Safehouse for users who chose the unsandboxed lane.
- Redesigning Claude permission policy or changing Codex and Pi permission behavior. The shared macOS launch-banner hint intentionally appears for all three hosts.

## Acceptance criteria

**AC-1 - Scrubbed, unsandboxed Claude first officers retain auto policy.** On the supported current Claude host, the candidate's machine-written `system/init.permissionMode` equals `auto` and stderr contains zero forced-default warnings. The scrub-on launch without Spacedock's minimal declaration independently records `default` and the upstream warning. No dialog, model-selected tool call, or final answer supplies evidence.

**AC-2 - Credential scrubbing remains active in child processes.** A deterministic `SessionStart` hook receives fabricated AWS variables in the scrub-off control and records only their presence. With scrub enabled, the hook records `credential_vars=scrubbed` in both unsandboxed and Safehouse rows. No credential value enters stdout, stderr, the Claude transcript, or the evidence artifact.

**AC-3 - The compatibility declaration adds no unsandboxed mutation authority.** The launcher adds only the externally proved read-only rule, never `Bash`, `Edit`, `Write`, a bypass flag, or a second operator allow list. Operator modes, both allow-list spellings, resumes, and Safehouse retain their exact authority and order.

**AC-4 - The macOS launch banner makes available Safehouse discoverable without changing containment.** Every fresh, unselected launch with Safehouse available prints the exact two adjacent lines above once. Selected Safehouse keeps exactly `Sandbox: enabled (safehouse)` with no hint; unavailable Safehouse, non-macOS platforms, and resumes print no hint.

**AC-5 - Unsupported Claude command surfaces fail with a stable remedy that preserves hardening.** Doctor distinguishes the supported and missing-capability fixtures with the exact output above. The remedy offers a Claude update or Safehouse and never advises disabling scrub mode.

**AC-6 - The documentation distinguishes policy from isolation.** The command reference names the scrub-compatible launch behavior and macOS hint; the sandbox reference says auto is policy and Safehouse is OS containment. Documentation render and link checks pass.

## Test plan

- Add table-driven argv tests in `internal/cli` for unsandboxed/Safehouse, fresh/resume, explicit mode, `--allowedTools`, and `--allowed-tools`. Assert exact token identity, order, and cardinality.
- Add a fake-Claude launch journey that emits the upstream forced-default warning only for scrub + auto + absent allow declaration. Assert the corrected argv removes that failure without changing the supplied environment.
- Add a live-gated macOS journey that runs the normal Spacedock bootstrap with stream JSON and the presence-only SessionStart hook. Assert the four-row matrix from `system/init`, hook state, stderr, and exit codes. Its evidence reader must reject the captain's old refusal transcript because it lacks the required machine mode/scrub pair.
- Add exact full-banner tests over `{darwin, linux} × {selected, unselected} × {available, unavailable}`. Assert line adjacency and cardinality; cover Claude, Codex, and Pi through the shared renderer; retain the existing enabled line unchanged.
- Add doctor fixture tests for supported and missing host flags, including exact stdout, stderr, and exit codes.
- Protect the external four-row artifact as a release gate. Record the Claude version, `system/init` row, boolean hook file, banner, stderr warning count, and exit code. No model response or credential value is evidence.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` after focused tests.

### Feedback Cycles

**Cycle 1 (captain external-probe result).** The proposed probe is not
exercisable as written. In a real external Claude session, `/permissions`
dismissed without durable mode evidence and the model correctly refused to run
the requested command because it would print `AWS_ACCESS_KEY_ID` into the
transcript. The compatibility prompt also instructed Claude to wait, so the
normal first-officer boot path did not run.

Routed back to ideation:

- Replace the credential-value echo and model judgment with a safe,
  deterministic observation. Do not print any credential value or rely on a
  permissions dialog that produces no record. A presence-only check may be
  considered only if it actually distinguishes scrub behavior without exposing
  secret material; otherwise use host/process evidence outside model tool
  selection.
- Separate permission-mode evidence from credential-scrub evidence so a model
  refusal cannot make the matrix inconclusive. Preserve the real transcript as
  a negative test for the old probe.
- Add this exact macOS launcher UX requirement when Safehouse is available but
  no profile enables it:

  ```text
  Sandbox: available, not enabled (no .safehouse profile)
  Use --safehouse to enable.
  ```

  Cover its platform/availability conditions and exact output without changing
  the already-enabled Safehouse banner.

## Stage Report: ideation

- DONE: Define and, where externally runnable, exercise the scrub on/off × unsandboxed auto/Safehouse launch matrix, distinguishing Claude permission policy from process isolation and recording exact banners plus harmless tool behavior.
  Defined a four-row external macOS probe with banner, `/permissions`, prompt, and fake-cloud-credential evidence; issue #504 supplies the original failure, and the current-candidate run remains explicitly pending outside this sandbox.
- DONE: Decide whether env-scrub plus unsandboxed auto remains a supported combination or autonomous first-officer operation requires Safehouse, without disabling credential scrubbing or inventing host guarantees.
  Chose scrubbed unsandboxed auto as the target with one read-only declaration; a failed external mechanism probe triggers a Safehouse-required reframe instead of broader allow rules or disabled scrubbing.
- DONE: Produce an implementation-ready approach, end-state acceptance criteria, exact tests, and doctor/docs behavior for supported and unsupported Claude host configurations.
  Recorded the launcher branch, override rules, exact diagnostics, documentation text, four behavioral ACs, focused fixtures, and the mandatory external release artifact.

### Summary

Issue #504 exposes a conflict between two supported policies, not a failure to opt into Safehouse. The design preserves both credential scrubbing and unsandboxed auto with the smallest read-only declaration Claude proves sufficient; implementation must wait for the external terminal probe and must reframe to a Safehouse requirement if that proof fails.

## Stage Report: ideation (cycle 2)

- DONE: Replace the refused credential-value/model-mediated probe with a safe deterministic exercise that separately proves permission-mode behavior and subprocess credential scrubbing without printing secret values.
  Permission mode now comes from `system/init`; a SessionStart hook records only fabricated credential-variable presence, and the normal first-officer bootstrap runs without a wait prompt.
- DONE: Incorporate the real refused transcript as negative evidence and decide the supported auto-versus-Safehouse posture from observations that do not depend on model compliance or a non-recording permissions dialog.
  The refusal now fails the evidence contract; unsandboxed auto remains the target only if the machine-written candidate row passes, else scrubbed autonomous operation requires Safehouse.
- DONE: Add the exact macOS no-profile hint `Use --safehouse to enable.` beneath the existing sandbox banner, with implementation-ready conditions, acceptance criteria, and exact-output tests.
  The design gates the hint on macOS + available + unselected, shares it across fresh host banners, and preserves enabled, unavailable, non-macOS, and resume output.

### Summary

Cycle 2 removes both unsafe credential display and model judgment from the mechanism proof. The revised entity is implementation-ready after the captain records the external candidate row; the banner hint is fully specified and independent of that result.
