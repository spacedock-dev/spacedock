---
title: Increase multi_agent_v2 default async wait timeout
status: ideation
score: 0.5
source: "Captain filing request 2026-07-11."
completed:
verdict:
worktree:
issue:
id: 95we0fhydgx5rbay5fw3qy4q
started: 2026-07-11T11:51:34Z
---

`multi_agent_v2`’s short default wait causes repeated timeout churn while subagents perform normal verification or interactive work. Because the wait already returns immediately when an agent responds or new user input arrives, increase the default timeout to five minutes. This reduces unnecessary polling and tool traffic without delaying updates or reducing responsiveness.

## Problem

The current live Codex collaboration surface accepts an omitted `timeout_ms`, but
that omission produced a timeout after about 34 seconds of wall-clock time in
the ideation probe (the call itself has orchestration overhead; this establishes
the practical current default as roughly 30 seconds). A normal ensign that is
verifying, running a focused test, or waiting for an interactive clarification
therefore causes repeated timeout/re-wait churn even though it is healthy.

That churn is not needed for responsiveness. The live `wait_agent` contract
states that a mailbox update and new user input end the wait early. In the probe,
a `timeout_ms: 60000` wait returned `timed_out: false` after about 21 seconds
when a parent mailbox message arrived. The completion/user-input behavior is
therefore distinct from the timeout ceiling. A sibling agent's final status did
not wake this worker's five-minute wait, which is useful negative evidence that
the wait observes the caller's collaboration scope rather than acting as a
global completion bus.

## Evidence and ownership boundary

- **Live surface:** `wait_agent(timeout_ms?)` currently accepts 10,000 through
  3,600,000 milliseconds. Its surfaced behavior says it returns early for
  mailbox updates and steered user input. The no-argument probe ran from
  `2026-07-11T12:01:40.3Z` to `2026-07-11T12:02:14.3Z` and timed out; the
  60-second mailbox probe ran from `12:08:03.3Z` to `12:08:24.3Z` and completed
  early.
- **Actual owner:** the omitted-timeout default is owned by the Codex runtime,
  not this Go repository or the Spacedock plugin. Spacedock currently owns only
  its first-officer instructions and the `spacedock codex` argv it launches.
- **Supported Codex configuration:** the current Codex schema and both local
  executables expose
  `features.multi_agent_v2.default_wait_timeout_ms`, with a 0--3,600,000 ms
  range. On `codex-cli 0.144.1` and the bundled `0.144.0-alpha.4`,
  `--strict-config -c 'features.multi_agent_v2.default_wait_timeout_ms=300000'`
  passed configuration parsing (then correctly stopped only because stdin was
  not a terminal); an unknown nested key and `3600001` each failed validation.
  The initially plausible root key `multi_agent_v2.default_wait_timeout_ms` is
  invalid. Do not use it.
- **Repository boundary:** a user-level `~/.codex/config.toml` changes the
  operator's Codex default. A trusted project `.codex/config.toml` changes only
  that project. Neither makes an installed Spacedock plugin change an already
  open Codex session or every repository where the plugin is invoked. The
  current `runCodex` front door also emits no timeout override.

## Chosen direction

Bind the Spacedock Codex first officer's *foreground wait* explicitly to five
minutes: `wait_agent(timeout_ms: 300000)`. This is the smallest durable
in-repository behavior because it applies wherever the shipped first-officer
skill is used--including an existing Codex session that did not start through
`spacedock codex`--and it uses the already-probed live tool call shape.

Keep the live-surface probe generic (`wait_agent(timeout_ms)`) so a missing or
incompatible tool remains a concrete runtime blocker. In the normative Codex
wait notes, replace the currently unspecified call with the explicit
five-minute value and say that it is a Spacedock per-call policy, not a claim
that Codex's global default changed.

Implementation changes only these product surfaces:

1. `skills/first-officer/references/codex-first-officer-runtime.md` -- bind the
   normative idle wait to `timeout_ms: 300000` while preserving timeout and
   interruption lifecycle rules.
2. The focused `internal/contractlint` wait-shape expectations -- keep the
   structural command/cue invariant aligned with the rendered adapter.
3. A focused Codex live scenario under `internal/ensigncycle` -- prove a
   delayed owned worker crosses the old short default without a timeout/re-wait
   and still ends the five-minute wait as soon as it completes.

Do **not** add a global setting, modify the operator's `~/.codex/config.toml`,
add a repository `.codex/config.toml`, or add `-c` to `runCodex` in this
change.

## Alternatives considered

1. **Session configuration in `runCodex`:** inject
   `-c features.multi_agent_v2.default_wait_timeout_ms=300000`. This can change
   the default for sessions launched by `spacedock codex`, and the exact current
   key has been parser-proven. It does not cover a skill invoked in an existing
   session, changes a broader host-session default than the Spacedock wait, and
   makes the front door fail on older Codex builds that do not support the
   nested feature configuration. Defer unless a later task explicitly wants a
   launcher-session policy with version/override handling.
2. **Check in `.codex/config.toml`:** applies only when the Spacedock source
   checkout itself is the trusted workspace. It would not travel as active
   configuration with the cached plugin, so it is unsuitable for the requested
   operational behavior.
3. **Ask operators to edit global config:** this changes state outside the
   repository, is per-user rather than shipped behavior, and is out of scope.
4. **Leave the timeout omitted:** preserves the current approximately
   30-second churn and fails the operational goal.

## Documentation impact

The shipped Codex adapter is the user-facing contract surface. No CLI output,
command, or site documentation changes because no launcher or configuration
surface changes. The implementation must apply this adapter diff:

```diff
@@ ## Codex wait notes
-When there is an unresolved Codex worker and no other dispatchable, gate, or state work, the FO MUST call `wait_agent(timeout_ms)` before ending the turn or reporting idle/status.
+When there is an unresolved Codex worker and no other dispatchable, gate, or state work, the FO MUST call `wait_agent(timeout_ms: 300000)` before ending the turn or reporting idle/status. This five-minute ceiling is Spacedock's per-call foreground-wait policy; it does not change Codex's global default.
```

## Acceptance criteria

**AC-1 (value): A healthy owned worker that takes longer than the current short
wait completes without timeout churn.** In an isolated Codex live run, a
no-write/delayed owned worker taking at least 45 seconds causes one foreground
wait to remain active through the delay and reach the worker's completion; no
timeout/re-wait occurs before that completion. The measured delay exceeds the
approximately 30-second baseline, so an omitted/default wait can fail this
criterion.
Verified by: a focused `//go:build live` Codex scenario that captures the
collaboration stream plus durable workflow state, with an assertion that fails
on a timeout or repeated wait before the owned worker finishes.

**AC-2 (value): Five-minute waiting remains responsive to a real completion.**
The same delayed-worker run returns from its 300,000-ms foreground wait on the
worker completion rather than waiting out five minutes, then verifies the
worker's committed stage report before the workflow advances.
Verified by: the live scenario's timestamped wait/completion trace and durable
stage-report/git-state assertion; it must finish shortly after the delayed
worker completes, not at the timeout ceiling.

**AC-3 (scope and safety): Spacedock's Codex wait policy is explicit without
changing host-wide configuration or interruption semantics.** The shipped
adapter names the per-call 300,000-ms duration, keeps a timeout normal and
retryable, and preserves the explicit rule that an operator interruption does
not fail, close, or redispatch the worker. The implementation changes neither
the global Codex config nor `spacedock codex` argv.
Verified by: focused contractlint structural coverage for the adapter's required
call/cue shape, targeted front-door argv regression coverage showing no added
`-c` override, and the live AC-1/AC-2 run exercising the actual timeout
behavior rather than treating instruction-text matching as behavioral proof.

## Test plan

1. Add the live delayed-worker control first. Use the existing isolated
   `CODEX_HOME` Codex runner so no developer-local config can mask the baseline.
   Delay the owned no-write worker for 45 seconds--longer than the observed
   short default but below the existing 60-second stream-silence guard. Record
   the collaboration wait begin/completion events, timeout/retry count, process
   exit, entity body, and state-checkout git log. Estimated cost: one Codex run,
   roughly 1--2 minutes.
2. Add/adjust focused contractlint only for the stable adapter structure and
   safety cue. It is a guard against accidental wording/shape drift, not the
   proof of AC-1 or AC-2.
3. Run the required offline gates: `go test ./...`, `go test ./... -race`, and
   `gofmt -w ./cmd ./internal`. Because the changed adapter is loaded by the
   Codex live lane, run the focused live control and then
   `go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` before validation/merge.
4. If a later task elects the rejected launcher-config alternative, add separate
   parser/argv tests for the exact `features.multi_agent_v2...` key, an
   unsupported-version failure path, user-supplied override precedence, fresh
   vs resume behavior, and a fresh-process live proof. Those tests are not part
   of this per-call change.

## Out of scope

- Changing Codex's product-wide default or asking users to edit global config.
- Adding a user-facing Spacedock timeout flag, profile, or configuration file.
- Changing Claude or Pi wait behavior.
- Treating a sibling/cross-tree completion as the foreground wait's completion
  signal.
- Replacing the existing timeout/interruption lifecycle guarantees.

## Stage Report: ideation

- DONE: Identify the real configuration owner and supported surface without
  editing global settings. The actual key is
  `features.multi_agent_v2.default_wait_timeout_ms`; both local Codex binaries
  accepted 300,000 under strict configuration parsing, while wrong-key and
  out-of-range controls failed.
- DONE: Establish current wait and early-return behavior with the smallest safe
  live probe. An omitted wait timed out after about 34 seconds; a 60-second
  wait returned early on a mailbox message after about 21 seconds. The surfaced
  tool contract also explicitly guarantees early return on steered user input.
- DONE: Choose the smallest compatible repository behavior. Bind a five-minute
  value explicitly in the Codex first-officer adapter, rather than claiming a
  global default changed or injecting a config override that only applies to
  front-door launches.
- DONE: Specify a live delayed-worker control that measures the operational
  value against the observed short baseline and proves completion still wakes
  early.

### Summary

Codex owns the omitted wait default, but it exposes a validated session config
surface at `features.multi_agent_v2.default_wait_timeout_ms`. Spacedock should
not alter that global/session default for this task: the robust, compatibility-
first behavior is an explicit `wait_agent(timeout_ms: 300000)` in the shipped
Codex first-officer adapter. A 45-second owned-worker live control will prove
that the new five-minute ceiling removes the approximately 30-second timeout
churn without delaying completion handling.
