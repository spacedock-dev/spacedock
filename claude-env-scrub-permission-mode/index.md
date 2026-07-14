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

This Codex worker already runs inside a sandbox, so it cannot prove the external unsandboxed boundary. Issue #504 supplies the original external scrub-on failure. A current-candidate result remains pending the following captain-run probe from a normal macOS Terminal.

Set up two disposable project settings without changing `~/.claude/settings.json`:

```bash
PROBE="$(mktemp -d /tmp/spacedock-504.XXXXXX)"
mkdir -p "$PROBE/off/.claude" "$PROBE/on/.claude"
printf '%s\n' '{"env":{"CLAUDE_CODE_SUBPROCESS_ENV_SCRUB":"0","AWS_ACCESS_KEY_ID":"SPACEDOCK_504_SENTINEL"}}' > "$PROBE/off/.claude/settings.local.json"
printf '%s\n' '{"env":{"CLAUDE_CODE_SUBPROCESS_ENV_SCRUB":"1","AWS_ACCESS_KEY_ID":"SPACEDOCK_504_SENTINEL"}}' > "$PROBE/on/.claude/settings.local.json"
BIN=/opt/homebrew/Caskroom/spacedock@next/0.25.0-pre2/spacedock
"$BIN" doctor --host claude
```

Run each lane separately. In each session, save the launch banner, open `/permissions`, then ask: `Use Bash once to run exactly: printf 'SPACEDOCK_504_ACTION_OK scrub=%s\n' "${AWS_ACCESS_KEY_ID-unset}". Do nothing else.` Record the displayed mode, whether Claude prompts, and the printed `scrub=` value.

```bash
(cd "$PROBE/off" && "$BIN" claude "Compatibility probe only: wait for my next message.")
(cd "$PROBE/on" && "$BIN" claude "Compatibility probe only: wait for my next message.")
(cd "$PROBE/on" && "$BIN" claude "Compatibility probe only: wait for my next message." -- --allowedTools Read)
(cd "$PROBE/on" && "$BIN" claude "Compatibility probe only: wait for my next message." --safehouse)
```

The four rows are scrub-off unsandboxed baseline, scrub-on unsandboxed reproduction, scrub-on unsandboxed candidate mechanism, and scrub-on Safehouse control. The candidate passes only if it keeps `auto`, emits no forced-default warning, runs the harmless action under auto policy, and prints `scrub=unset`. The Safehouse control must show `Sandbox: enabled (safehouse)`, retain its bypass-with-containment posture, and also print `scrub=unset`. A fixture or an in-sandbox run cannot replace these observations.

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

## Documentation change

In `docs/site/reference/command-reference.md`, replace the unsandboxed-launch paragraph with:

> An unsandboxed bootstrap launch carries no Safehouse isolation. `spacedock claude` starts in `--permission-mode auto`, and Codex starts in `--ask-for-approval on-request`, unless you supply a mode. When Claude subprocess credential scrubbing is enabled, Spacedock adds a minimal read-only declaration so auto mode and credential scrubbing compose; it does not pre-approve Bash or writes. A sandboxed bootstrap instead bypasses prompts inside Safehouse's OS boundary. Run `spacedock doctor --host claude` if Claude rejects the scrub-compatible surface. Spacedock never disables credential scrubbing.

Expand `docs/site/reference/sandbox.md` after the table:

> Permission mode is policy, not process isolation. Claude `auto` classifies tool calls but does not confine the process. Safehouse enforces filesystem and network boundaries around the host. Both launch lanes are supported; choose Safehouse when you need OS-level containment.

## Out of scope

- Disabling or clearing `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB`.
- Broadly pre-approving Bash, file mutation, network access, or MCP tools.
- Automatically enabling Safehouse for users who chose the unsandboxed lane.
- Redesigning Claude permission policy or changing Codex and Pi launch behavior.

## Acceptance criteria

**AC-1 - Scrubbed, unsandboxed Claude first officers retain auto policy and credential isolation.** On the supported current Claude host, the candidate launch shows `auto`, emits no forced-default warning, completes the harmless Bash action without a manual permission round trip, and exposes no `AWS_ACCESS_KEY_ID` in the Bash child. The same launch without Spacedock's minimal declaration is the independent failing baseline.

**AC-2 - The compatibility declaration adds no unsandboxed mutation authority.** The launcher adds only the externally proved read-only rule, never `Bash`, `Edit`, `Write`, a bypass flag, or a second operator allow list. Operator modes, both allow-list spellings, resumes, and Safehouse retain their exact authority and order.

**AC-3 - Unsupported Claude command surfaces fail with a stable remedy that preserves hardening.** Doctor distinguishes the supported and missing-capability fixtures with the exact output above. The remedy offers a Claude update or Safehouse and never advises disabling scrub mode.

**AC-4 - The documentation distinguishes policy from isolation.** The command reference names the scrub-compatible launch behavior; the sandbox reference says auto is policy and Safehouse is OS containment. Documentation render and link checks pass.

## Test plan

- Add table-driven argv tests in `internal/cli` for scrub off/on representation, unsandboxed/Safehouse, fresh/resume, explicit mode, `--allowedTools`, and `--allowed-tools`. Assert exact token identity, order, and cardinality.
- Add a fake-Claude launch journey that emits the upstream forced-default warning only for scrub + auto + absent allow declaration. Assert the corrected argv removes that failure without changing the supplied environment.
- Add doctor fixture tests for supported and missing host flags, including exact stdout, stderr, and exit codes.
- Protect the external four-row macOS journey above as a manual/live release artifact. It must record the Claude version, banner, permissions view, prompt behavior, and child `scrub=` output. No PR may claim AC-1 from fixture evidence alone.
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
