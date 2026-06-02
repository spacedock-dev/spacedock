---
id: r09jrf0k6qjv6c1sddhe1sh6
title: Ship a Codex runtime and install-readiness contract — codex ensign/first-officer adapters + Codex-shaped dispatch + wait/observe semantics
status: validation
source: "captain (2026-06-02) — live Codex FO session: 0.19.x ships Claude-only runtime adapters; Codex FO/ensign dispatch/wait/observe is improvised (wait_agent timed out 10s then async notifications), and dispatch build emits the Claude-ism Skill(skill=...) which the Codex skill list does not expose"
started: 2026-06-02T15:38:31Z
completed:
verdict:
score: "0.33"
worktree: .worktrees/spacedock-ensign-codex-runtime-adapter
issue:
milestone: post-0.19.4 (own track — Codex runtime parity, NOT folded into 0.19.4)
mod-block: 
pr: spacedock-dev/spacedock#269
---

0.19.x ships **runtime-agnostic shared cores with Codex-aware mentions** (`send_input` routing, "Codex declares none" budget probe, codex resume) but **no Codex runtime contract**. The result is the "model-improvised, not contract-guaranteed" antipattern at runtime scale — confirmed live by the captain's Codex FO session. Captain feedback on 2026-06-02 folds the adjacent Codex install-readiness gap into this entity, because a Codex runtime contract is not usable if `spacedock install --host codex` stays blind to an already-installed compatible plugin.

Confirmed gaps in the shipped plugin:
- Only `skills/ensign/references/claude-ensign-runtime.md` and `skills/first-officer/references/claude-first-officer-runtime.md` exist — **no `codex-*-runtime.md`**.
- Both `SKILL.md` files dispatch the runtime adapter on `CLAUDECODE` only (`ensign/SKILL.md:13`, `first-officer/SKILL.md:23`) — there is no `Codex → read codex-runtime.md` branch, so a Codex FO/ensign loads only the shared core and improvises completion/wait/observe.
- `spacedock dispatch build` hard-codes the Claude-ism `Skill(skill="spacedock:ensign")` in the emitted prompt (`internal/dispatch/build.go:300,431`); the Codex session's skill list did not expose `spacedock:ensign`, forcing a read-the-dispatch-file-directly workaround.
- The Codex FO observed: `wait_agent` timed out after 10s, then Codex delivered async subagent-completion notifications — that is **Codex host behavior, not a Spacedock FO contract**.
- `spacedock install --host codex` currently prints static manual install prose even when `spacedock doctor --host codex` can already resolve a compatible installed plugin. The captain hit the concrete failure mode: `codex plugin marketplace add spacedock-dev/spacedock --ref next` errors when the `spacedock` marketplace is already registered from another source, and after manual plugin install the Spacedock install command still prints the same stale instructions instead of reporting that the Codex plugin is installed.

## Design surface (for ideation to scope)

- **Codex first-officer runtime adapter** (`codex-first-officer-runtime.md`) — the Codex analog of the Claude adapter's `## Awaiting Completion` / dispatch / reuse sections: how a Codex FO observes a completion signal (the wait/notification model), how it routes reuse/feedback (`send_input` per the shared-core mentions), the team-vs-no-team model on Codex, and the budget-probe declaration (shared-core already says "Codex declares none").
- **Codex ensign runtime adapter** (`codex-ensign-runtime.md`) — the Codex analog of `claude-ensign-runtime.md`: completion-signal protocol, worktree ownership, polling.
- **SKILL.md runtime dispatch** — add the Codex branch to both `ensign/SKILL.md` and `first-officer/SKILL.md` so the right adapter loads by platform (Codex env detection).
- **`dispatch build` prompt shape** — emit a Codex-appropriate prompt (not the Claude `Skill(...)` call) when targeting Codex, and resolve the skill-exposure mismatch (either ensure `spacedock:ensign` is exposed on Codex, or emit the read-the-dispatch-file form as the contract rather than a workaround).
- **Codex install readiness** — make `spacedock install --host codex` observe the same installed-plugin resolver/doctor result that `spacedock doctor --host codex` already uses. If the Codex plugin is installed and compatible, report that state and do not print manual marketplace-add steps. If no compatible plugin is installed, keep the manual Codex command guidance, but frame it as the required next action rather than a generic message.
- **Reconcile** with the existing codex launcher work (`be`/codex-safehouse-launcher — `spacedock codex` LAUNCHES codex), the install/doctor path, and the `spacedock-dev/spacedock` migration. Name the boundary so launcher, install readiness, and runtime don't double-file.

## Riskiest-unknown spike — Codex multi-agent wait/observe contract (DONE, grounded)

The riskiest claim is AC-3: how a Codex FO knows a dispatched worker finished. Per the spike rule it was exercised first, against the real `codex` binary on PATH (`codex-cli 0.136.0`) — not assumed. `codex features list` shows `multi_agent  stable  true`; the multi-agent toolset and its descriptions were read directly from the binary. Authoritative findings (verbatim from the binary's tool descriptions, evidence command in the stage report):

- **Tool surface (`multi_agent_v1`, "Tools for spawning and managing sub-agents"):** `spawn_agent`, `send_input`, `send_message`, `wait_agent`, `close_agent`, `resume_agent`, plus an `agent_status` field on results. There is **no `TeamCreate` analog** — sub-agents are organized by a canonical task-path tree (`spawn_agent` with `task_name "task_3"` under parent `/root/task1` → `/root/task1/task_3`). Codex has no team-registry lifecycle; the "team-vs-no-team" axis on Claude has no Codex counterpart.
- **`wait_agent` (the captain's observed 10s-timeout mechanism), verbatim:** *"Wait for agents to reach a final status. Completed statuses may include the agent's final message. **Returns empty status when timed out. Once the agent reaches a final status, a notification message will be received containing the same completed status.**"* And: *"Wait for a mailbox update from any live agent… Does not return the content; returns either a summary of which agents have updates (if any), or a timeout summary if no mailbox update arrives before the deadline."* The result carries a boolean *"Whether the wait call returned due to timeout before any agent reached a final status."*
- **Guidance, verbatim:** *"Call wait_agent very sparingly. Only call wait_agent when you need the result immediately for the next critical-path step and you are blocked until it returns."*
- **Timeouts:** per-worker default runtime `1800` seconds ("config may set a different default"); wait has config keys `default_wait_timeout_ms` / `max_wait_timeout_ms` and a per-call `timeout_ms`.
- **`send_input`, verbatim:** *"Send a message to an existing agent. Use interrupt=true to redirect work immediately. You should reuse the agent by send_input if you believe your assigned task is highly dependent on the context of a previous task."* (This is exactly the shared-core's reuse/feedback routing claim.)
- **`close_agent`:** *"Close an agent and any open descendants when they are no longer needed… Don't keep agents open for too long if they are not needed."* (the explicit-shutdown analog.)

**What this resolves.** The captain saw `wait_agent` time out at ~10s then async completion notifications arrive — that was NOT a malfunction and NOT improvisation needing fixing; it is the documented design. The completion signal on Codex is **the async final-status notification delivered to the FO's mailbox**, not the `wait_agent` return value. `wait_agent` is an *optional* blocking accelerator the FO uses sparingly when it is blocked on a critical-path result; a timeout return is normal and means "no mailbox update yet," not "worker failed." This is the Codex realization of Claude's `## Awaiting Completion` (where the completion signal is a `task_notification` system entry, and the FO ends the turn empty between dispatch and signal). The parity wording writes itself from this: **end-the-turn-and-await-the-notification is the contract; `wait_agent` is the sparing accelerator, never the polling loop.**

No further spike is needed for AC-1/AC-2 mechanism *soundness* — they compose already-proven Spacedock behavior (SKILL.md `@`-include + env-branch; `dispatch build` JSON envelope assembly) — but each still ships a code-checkable artifact (see test plans). AC-4's install-readiness mechanism is grounded in the already-shipped Codex resolver: in the captain's environment `spacedock doctor --host codex` resolves `spacedock@spacedock` and reports a compatible contract, while `spacedock install --host codex` ignores that resolver and prints static prose. One sub-mechanism remains genuinely unverified and is called out as an open question below: the **Codex session env-detection signal** (the analog of `CLAUDECODE=1`).

## Open question to resolve in implementation (not a blocker for ideation)

**How does an agent detect it is running under Codex?** Claude branches on `CLAUDECODE`. Binary strings show Codex sets `CODEX_SANDBOX` (and `CODEX_SANDBOX_NETWORK_DISABLED`) inside its sandbox — but the `spacedock codex` launcher (`internal/cli/frontdoor.go`) runs `codex --dangerously-bypass-approvals-and-sandbox`, which may suppress the sandbox-only var. The implementation MUST exercise the actual var present in a launched `spacedock codex` session (a one-line `env | grep -i codex` in a live session) before committing the branch condition in both SKILL.md files. If no reliable env var exists, the fallback is a negative branch: `CLAUDECODE` unset → Codex adapter (with a comment naming the dependency). This is mechanism-level and belongs to the implementation's first test, per the spike discipline.

## Acceptance criteria

Each AC names a property of the finished entity and a check outside the entity body that can fail.

**AC-1 — A Codex FO and a Codex ensign each load a dedicated runtime adapter, not just the shared core.** Two new files exist — `skills/first-officer/references/codex-first-officer-runtime.md` and `skills/ensign/references/codex-ensign-runtime.md` — and both `SKILL.md` files carry a Codex branch (alongside the existing `CLAUDECODE` branch) that points at the Codex adapter. Each adapter has a section mapping to every load-bearing Claude-adapter section: Awaiting Completion (the wait/notification model from the spike), Dispatch, reuse/feedback routing (`send_input`), and ensign completion-signal. Each adapter declares the budget-probe posture (Codex declares none — already the shared-core wording).
- *Verified by:* a Go presence-test in `internal/hostneutrality/` (the package that already polices prose parity) that (a) both adapter files exist, (b) each `SKILL.md` body contains a Codex branch line referencing its `codex-*-runtime.md`, and (c) each Codex adapter contains the required section headings (`Awaiting Completion`, `Dispatch`, the `send_input` reuse term, and a completion-signal section). This is proof-at-the-claim's-level: the claim is "the contract text carries these clauses," and a presence/structure check over the real instruction files is legitimate per the ideation rule (a property of the text when the text is the claim). The wait/observe *behavior* claim is AC-3's, proven separately.

**AC-2 — `dispatch build` emits a Codex-target prompt that does not depend on a `Skill(skill=...)` call.** The dispatch envelope gains a host discriminator (a `host` input field, default `claude` for back-compat-preserving current callers, plus `codex`). On `host: "codex"`, the emitted `prompt` is the read-the-dispatch-file form WITHOUT a leading `Skill(skill="spacedock:ensign")` call (the spike-confirmed risk: the Codex skill list may not expose `spacedock:ensign`). The dispatch-file body's first-action block is correspondingly Codex-shaped (read-this-file as the operating contract entry, not "call the Skill tool"). The completion-signal block emitted for Codex names the Codex mailbox/notification contract, not `SendMessage(to="team-lead")`.
- *Verified by:* a Go table-test in `internal/dispatch/` (sibling to `build_advisory_probe_test.go`, which already proves the Claude-vs-nil envelope divergence). Two arms over one fixture: `host` absent/`claude` → current envelope unchanged (the existing parity tests are the regression floor); `host: "codex"` → the emitted `prompt` contains no `Skill(skill=` substring, the dispatch-file first-action block matches the Codex shape, and the generated dispatch-file completion section names the Codex mailbox/notification signal while omitting `SendMessage(to="team-lead")`. The skill-exposure-mismatch resolution is the *decision* (emit read-the-dispatch-file form as the contract, do not rely on skill exposure); the binary test is the proof that the emitted bytes honor it. The "verified against the real Codex skill-exposure behavior" half is the spike already on record (Codex session did not expose `spacedock:ensign`), cited not re-run.

**AC-3 — The Codex wait/observe contract is explicit in the adapter, grounded in the spike, not host-incidental.** `codex-first-officer-runtime.md`'s Awaiting Completion section states: (a) the completion signal is the async final-status notification in the FO mailbox; (b) `wait_agent` is a sparing, optional accelerator for critical-path blocking, and a `wait_agent` timeout return is normal (not a failure / not a teardown trigger); (c) between dispatch and notification the FO does not poll, re-dispatch, or close the worker — the Codex analog of Claude's end-turn-empty rule. `codex-ensign-runtime.md` states the ensign's completion-signal emission (the message form the FO's notification observes).
- *Verified by:* the same `internal/hostneutrality/` presence-test as AC-1, extended with anchored assertions that the Codex FO adapter's Awaiting-Completion section contains the load-bearing phrases (the notification-is-the-signal clause, the `wait_agent`-timeout-is-normal clause, and the do-not-teardown-on-timeout clause). Proof at the claim's level (the claim is that the contract text carries these specific clauses). The underlying *behavioral truth* of these clauses is the binary-grounded spike recorded above — the test pins that the adopted wording survives future edits.

**AC-4 — `spacedock install --host codex` is aware of an already-installed compatible Codex plugin.** The Codex install path first consults the installed-plugin resolver/doctor result. When `spacedock@spacedock` is installed and its manifest satisfies the binary contract, `spacedock install --host codex` exits 0 with the compatible/installed report and does not print `codex plugin marketplace add` or `codex plugin add` instructions. When no installed plugin is resolvable, the command keeps the current manual Codex install guidance. When the installed plugin is present but incompatible, the command surfaces the doctor verdict/remedy instead of implying a fresh marketplace add is always the next step.
- *Verified by:* Go unit tests in `internal/cli/` over the existing `hostOps` seam: compatible Codex manifest -> stdout contains the doctor-compatible report and lacks both manual Codex add commands; no resolved manifest -> stdout contains the manual command pair; incompatible manifest -> stderr/stdout carries the doctor mismatch/remedy and does not claim the plugin is absent. The existing live `TestCodexResolveManifestAgainstInstalledHost` remains the resolver smoke test; AC-4's command behavior is proven hermetically through `runInit`.

## Boundary — launcher, install readiness, and runtime (so they don't double-file)

- **`spacedock codex` launcher (`internal/cli/frontdoor.go`, shipped):** owns *getting into* a Codex session — version-gate, safehouse-wrap, `--dangerously-bypass-approvals-and-sandbox`, the bootstrap prompt that assumes `$spacedock:first-officer`. It stops at "a Codex session is running with the FO skill assumed."
- **Codex install readiness (`internal/cli/init.go` + `internal/cli/host_exec.go`):** owns making `spacedock install --host codex` truthful about whether the plugin is already installed and compatible before the user enters a Codex session. This is folded into this entity by captain direction because it gates usable Codex runtime onboarding.
- **This entity (Codex *runtime* contract):** owns FO/ensign *behavior once launched* — adapter files, SKILL.md branch, dispatch-build prompt shape. It begins where the launcher ends.
- **The one shared seam:** the env-detection signal (open question above) is the handshake between them — the launcher's session is where the implementation must read the actual env var. That read is implementation work, not a launcher change.
- **`spacedock-dev/spacedock` migration:** orthogonal — a repo/module-path concern (`internal/release`, manifests), not a runtime-behavior concern. Adapter files live under `skills/*/references/` regardless of module path; no migration coupling.

## Test plan summary

- **Fixture/CLI cost:** AC-1 and AC-3 share one `internal/hostneutrality/` Go presence+structure test (cheap, no live host). AC-2 is one `internal/dispatch/` Go table-test over a fixture (the existing build-test harness pattern), now including the generated completion-signal block. AC-4 is a small `internal/cli/` unit-test set over `runInit` and the fake `hostOps` resolver. No live workflow test is required for the gate — runtime behavior is grounded by the binary-read spike, and install-command behavior is proven through the existing command seam.
- **Estimated complexity:** Moderate. The risk was the wait/observe mechanism, now retired. Remaining runtime work is contract authoring (two adapter files + two SKILL.md branches) plus one `host` field threaded through `dispatch build`'s prompt-assembly branches (sections 0, 4, and 10 in `runBuild`) and the `buildOutput` input parse. The install-readiness fold-in is a narrower CLI change in `runInit`: consult the Codex resolver/doctor result before printing manual install prose. The host-neutrality prose oracle already exists and constrains where Codex-specific tokens may live.
- **Live verification deferred to implementation (not a gate AC):** the env-detection `env | grep -i codex` read in a real `spacedock codex` session — a one-line mechanism check the implementer pays before wiring the SKILL.md branch condition, per the spike discipline.

## Notes
- OWN TRACK (captain decision 2026-06-02): file + ideate now, but do NOT fold into 0.19.4 (already large). Likely a 0.20 "Codex runtime parity" milestone.
- Fold-in (captain decision 2026-06-02): include `spacedock install --host codex` installed-plugin awareness here, despite the original launcher/runtime boundary, because it is the Codex onboarding prerequisite for the runtime contract.
- This is scaffolding (`skills/*/references/` + SKILL.md) plus an internal dispatch change; implementation goes through a worktree.

## Stage Report: ideation

- DONE: The Codex completion/wait/observe contract is named concretely and GROUNDED in a real Codex-session observation, not assumed.
  Read the multi-agent toolset + verbatim descriptions from the real binary (`codex-cli 0.136.0`, `multi_agent stable true`): `wait_agent` returns empty on timeout, then "a notification message will be received containing the same completed status" — confirming the captain's 10s-timeout-then-async-notification as documented design, not malfunction. Recorded in the "Riskiest-unknown spike" section with the evidence command below.
- DONE: Each new adapter section maps to its Claude counterpart; the SKILL.md Codex-branch (both files) and the dispatch-build Codex-prompt-shape change are specified, including how the spacedock:ensign skill-exposure mismatch is resolved.
  AC-1 names the two adapter files + the required section parity (Awaiting Completion / Dispatch / `send_input` reuse / completion-signal); AC-2 specifies the `host` discriminator on `dispatch build` and resolves the mismatch by emitting the read-the-dispatch-file form (no leading `Skill(skill=...)`) as the contract.
- DONE: The boundary vs the codex launcher and the spacedock-dev/spacedock migration is named; AC-2/AC-3 are framed as runtime-behavior proofs against real Codex, not prose-only.
  "Boundary — launcher vs runtime" section draws the launcher-ends/runtime-begins line and marks the migration orthogonal; AC-3's behavioral truth is the binary-grounded spike (cited), with the prose-presence test pinning the adopted wording. AC-2's skill-exposure half cites the live Codex session that did not expose `spacedock:ensign`.

### Summary

Hardened the ideation by retiring the riskiest unknown first: read the real `codex` binary's `multi_agent` tool contract (spawn_agent/send_input/wait_agent/close_agent/resume_agent, no TeamCreate analog) and confirmed the completion signal is the async final-status mailbox notification — `wait_agent` is a sparing accelerator, a timeout return is normal. ACs now ship code-checkable artifacts: a `internal/hostneutrality/` presence+structure test for the two adapter files and SKILL.md branches (AC-1/AC-3) and a `internal/dispatch/` table-test that the `host: "codex"` envelope emits no `Skill(skill=...)` call (AC-2). One genuine open question — the Codex session env-detection signal (`CODEX_SANDBOX` may be suppressed by the launcher's sandbox-bypass) — is flagged as implementation-first mechanism work, not a gate blocker.

### Spike evidence command

    f=$(strings /opt/homebrew/bin/codex | grep -c "wait_agentWait for agents to reach a final status"); echo "$f wait_agent description(s) in codex-cli 0.136.0"

The full verbatim descriptions were extracted via `strings /opt/homebrew/bin/codex | grep -E "wait_agent|send_input|spawn_agent|close_agent|resume_agent"` and `codex features list | grep multi_agent`.

## Stage Report: ideation feedback repair

- DONE: AC-2 now verifies the generated Codex completion-signal block, not only the outer prompt.
  The AC-2 "Verified by" clause now requires the `internal/dispatch/` table-test to assert that the generated dispatch-file completion section names the Codex mailbox/notification signal and omits `SendMessage(to="team-lead")`.
- DONE: The captain's Codex install-awareness report is folded into this entity with a code-checkable AC.
  AC-4 now covers `spacedock install --host codex`: compatible installed plugin -> compatible/installed report with no manual add commands; no resolved plugin -> manual command pair remains; incompatible installed plugin -> doctor mismatch/remedy. The proof is a hermetic `internal/cli/` test over `runInit` and `hostOps`.
- DONE: The review polish item was repaired.
  The AC-2 host-default wording now says "back-compat-preserving current callers."

### Summary

Repaired the ideation gate finding and folded the captain's install-readiness report into the task. The entity now owns four checkable outcomes: Codex adapter loading, Codex-shaped dispatch output including completion-signal wording, Codex wait/observe contract text grounded in the spike, and `install --host codex` awareness of an already-installed compatible plugin.

## Stage Report: implementation

- DONE: Codex install path reports installed or compatible state before manual add guidance
  Code commit `2c4cd877` changes `spacedock init --host codex` to resolve the installed manifest first; `TestInitCodexInstallReadiness` covers compatible, missing, and incompatible states.
- DONE: Codex dispatch build host=codex output omits Skill(...) and uses Codex mailbox completion semantics
  Code commit `2c4cd877` adds the optional dispatch `host` field and `TestBuildCodexHostPromptShape`, asserting the Codex prompt/body omit `Skill(...)` and name the FO mailbox notification.
- DONE: Codex FO and ensign runtime adapters load via a verified Codex branch
  Code commit `2c4cd877` adds both Codex adapters plus SKILL.md `CODEX_THREAD_ID` branches; live env spike observed `CODEX_CI=1` and `CODEX_THREAD_ID=019e8949-0679-70a3-aa9a-1f626004898d`.

### Summary

Implemented all four ACs in the code worktree: Codex install readiness, Codex-shaped dispatch build output, dedicated FO/ensign Codex adapters, and adapter contract tests for the wait/observe clauses. Verification run from the code worktree: `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`; `gofmt` wanted one unrelated existing comment change in `internal/status/enum_scope_test.go`, which was not committed.

## Stage Report: validation

- DONE: Each AC-1 through AC-4 is independently verified with reproducible evidence
  AC-1/AC-3: `go test ./internal/hostneutrality` passed and verified both Codex adapter files, SKILL.md `CODEX_THREAD_ID` branches, required sections, and anchored mailbox/wait clauses; AC-2: `go test ./internal/dispatch` passed and verified host `codex` prompt/body omit `Skill(...)` and Claude `SendMessage`; AC-4: `go test ./internal/cli` passed and verified compatible, missing, and incompatible Codex install states.
- DONE: Implementation test gates and command behavior are run from the code worktree
  Ran from `.worktrees/spacedock-ensign-codex-runtime-adapter`: `go test ./internal/hostneutrality ./internal/dispatch ./internal/cli` PASS; `go test ./...` PASS; `go test ./... -race` PASS; `git status --short` clean before and after tests.
- DONE: Validation recommends PASSED or REJECTED and flags stale or weak proof
  Recommendation: PASSED. Proof caveats: AC-2's new test directly covers only the Codex arm while existing dispatch parity tests carry the Claude/default back-compat floor; `gofmt -l ./cmd ./internal` still reports unrelated pre-existing `internal/status/enum_scope_test.go`.

### Summary

PASSED: implementation commit `2c4cd877` satisfies AC-1 through AC-4 with focused tests and full Go gates passing from the code worktree. The validation made no deliverable code changes. The only weak proof is structural rather than behavioral for some runtime prose, which matches the AC's stated proof level, and the only format caveat is an unrelated existing file outside the implementation commit.

## Stage Report: merge prep

- DONE: Rebased the code branch onto `origin/next` with only the r0 implementation commit retained
  Branch `spacedock-ensign/codex-runtime-adapter` now points at `13fde1b6`, one commit ahead of `origin/next`, with the expected nine-file diff (`405 insertions, 56 deletions`).
- DONE: Post-rebase verification reran from the code worktree
  `go test ./internal/cli ./internal/dispatch ./internal/hostneutrality` passed, `go test ./...` passed, and `go test ./... -race` passed. Touched Go files were formatted with `gofmt -w`; whole-tree `gofmt -l ./cmd ./internal` still reports unrelated pre-existing `internal/status/enum_scope_test.go`.

### Summary

Prepared r0 for the PR merge hook after validation approval: the branch is rebased onto `origin/next`, the branch pointer is attached to the rebased commit, and post-rebase tests passed.

## Stage Report: PR refresh

- DONE: Rebasing r0 incorporated the live-E2E fixes now on `origin/next`
  Rebased branch `spacedock-ensign/codex-runtime-adapter` from old PR head `13fde1b6` onto current `origin/next`, producing new head `528ffa17`.
- DONE: Required local gates passed after the rebase
  From `.worktrees/spacedock-ensign-codex-runtime-adapter`, `gofmt -w ./cmd ./internal` completed, `go test ./...` passed 738 tests in 12 packages, and `go test ./... -race` passed 738 tests in 12 packages.
- DONE: Updated the existing PR branch for rerun
  Force-with-lease pushed `spacedock-ensign/codex-runtime-adapter` so PR #269 reruns on the fixed base.

### Summary

r0 is refreshed on top of current `origin/next`; the previous Opus failure was from the stale pre-fix PR head and the updated PR branch is now waiting on fresh CI.
