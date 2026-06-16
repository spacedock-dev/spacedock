---
title: Make the host-neutral dispatch core genuinely runtime-neutral
status: validation
source: "captain + FO follow-up after 2y merged as v0.20.3 (2026-06-16) — remaining concern: the extracted host-neutral dispatch core still carries Claude/team-only language, even though Codex and Pi now load it too. Verified against origin/main 9bd1f46a: fo-dispatch-core.md still says first team-mode dispatch, keeps spawn-standing-all --team {team_name} in the core, requires team_name as helper output, and makes reuse depend on not being in bare mode."
started: 2026-06-16T20:02:00Z
completed:
verdict:
score: 0.66
worktree: .worktrees/spacedock-ensign-runtime-neutral-dispatch-core-cleanup
issue:
id: ezfwkw33awtqgztgr6v7bb59
sprint: 0204-structured-reads
sprint-readiness: ready
---

The 2y merge fixed the important reachability gap by extracting `fo-dispatch-core.md` and `fo-merge-core.md`, making the merge and dispatch ceremony reachable to Codex and Pi. One follow-up remains: the new `fo-dispatch-core.md` is named as host-neutral, but parts of it still describe Claude team mode as if it were universal.

## Problem

On `origin/main` at `9bd1f46a` (`v0.20.3` plus manifest stamp), `skills/first-officer/references/fo-dispatch-core.md` still carries host-specific assumptions:

- The module says it is loaded at the first "team-mode dispatch" even though Codex/Pi do not necessarily have team mode.
- Standing-teammate injection is in the host-neutral core and calls `spacedock dispatch spawn-standing-all --team {team_name}`, which is Claude-team-shaped.
- Reuse condition 1 says `Not in bare mode (teams available)`, which makes Codex `send_input` reuse look invalid despite Codex having a reusable mailbox handle.
- The dispatch builder section says `team_name` is a key field that MUST come from helper output, while Codex explicitly has no `team_name` lifecycle.
- The core says bare mode dispatch blocks until subagent completion, which is a host behavior, not a host-neutral invariant.

This is not a blocker for the 2y merge because the extracted cores are reachable and tested, but it is a real contract-quality issue: future Codex/Pi FOs load a host-neutral core that still reads partly like Claude.

## Desired direction

Keep the 2y split. Do not collapse the cores back into per-host copies.

Make `fo-dispatch-core.md` describe only host-neutral dispatch lifecycle rules:

- load at first worker dispatch, not first team-mode dispatch;
- dispatch via the runtime adapter's spawn call;
- reuse when the adapter exposes a live reusable handle and all generic reuse conditions pass;
- helper output fields are forwarded when emitted/applicable, with host adapters mapping concrete fields;
- blocking/nonblocking dispatch and bare/team behavior live in runtime adapters.

Move or gate Claude-only pieces in `claude-fo-dispatch.md`:

- `spawn-standing-all --team {team_name}` and standing teammate injection;
- Claude `team_name` / TeamCreate sequencing;
- bare-mode blocking semantics;
- any rule whose source of truth is Claude team behavior.

Codex should remain explicit that `spawn_agent` is initial dispatch, `send_input` is reuse/feedback routing, `wait_agent` is only a foreground idle wait, there is no `team_name` lifecycle, and there is no reconcile sweep. Pi should retain its `subagent(...)` / `pi-agent-teams` substrate split.

## Acceptance criteria

**AC-1 — The dispatch core contains no unconditional Claude/team-only dispatch requirements.**
Verified by a structural contractlint check in `internal/contractlint/` (the instruction-read quarantine — the only package allowed to read an instruction file's content) that scans `fo-dispatch-core.md` for the LITERAL Claude-team command/flag tokens `spawn-standing-all` and `--team {team_name}` and FAILS on any line carrying one UNLESS that sentence's subject is the adapter realizing the call (phrase-level exemption: `Claude adapter`, `the adapter realizes`, `adapter maps`, `Claude realization`). A bare same-line `runtime adapter` mention is NOT a sufficient exemption — see the spike determination below for why.

The expected value comes from the rule (these are literal host command tokens that must not appear as core imperatives), not from the file's own prose, so a valid host-neutral paraphrase that drops the literal passes and an inverted/host-coupled imperative fails. This is structural-absence in the same family as `TestRetiredPluginPrivatePathsAbsent` (a literal token's presence is itself the defect), NOT a prose-grep behavior substitute. The check carries a positive discriminator control (a host-coupled fixture line that MUST flag, a host-neutral paraphrase that MUST pass) so it can never pass vacuously. The narrowly-paraphrasable phrase `team-mode` is deliberately NOT a literal token in the scanner — its fix is the prose move, not a token ban; banning the word would be the prose-grep tautology the policy forbids.

**AC-2 — Codex reuse remains expressible through the shared core.**
Verified by a fixture or existing ensigncycle test showing a Codex continuation/reuse path uses `send_input` against a live worker handle and is not rejected merely because there is no team mode or `team_name`.

**AC-3 — Claude behavior is preserved.**
Verified by existing focused gates: `go test ./internal/contractlint`, `go test ./internal/dispatch -run 'Standing|Build'`, and the Claude standing/teardown focused tests that cover `spawn-standing-all`, TeamCreate sequencing, and the terminal marker.

## Concrete prose move

Specific before/after wording for the five host-coupled leaks. Keep the 2y split — do NOT collapse the cores. Host-neutral language stays in `fo-dispatch-core.md`; the Claude realization moves to / stays in `claude-fo-dispatch.md`; Codex states its mapping in `codex-first-officer-runtime.md`; Pi in `pi-first-officer-runtime.md`.

**Leak 1 — Standing-teammate injection + `spawn-standing-all --team {team_name}`** (`fo-dispatch-core.md` line 9, `## Dispatch`).

- BEFORE: "**Standing-teammate injection.** Before the first team-mode dispatch, inject the workflow's declared standing teammates: run `spacedock dispatch spawn-standing-all --workflow-dir {wd} --team {team_name}` and forward each spawn spec in the returned JSON array to the runtime adapter's spawn call (verbatim, same discipline as `spacedock dispatch build` output). The call is idempotent … Standing teammates are team-scoped: they die with the team at teardown."
- AFTER (core, host-neutral): "**Standing-teammate injection.** Before the first worker dispatch, inject the workflow's declared standing teammates via the runtime adapter's standing-injection call. The adapter resolves the workflow's declared standing teammates and forwards each returned spawn spec to its spawn call (verbatim, same discipline as `spacedock dispatch build` output). Injection is idempotent — already-alive members are omitted, so re-running is safe — and is a no-op when no standing teammate is declared or the runtime has no shared-teammate surface. Standing-teammate lifetime is the adapter's (team-scoped where the runtime has teams). Read each teammate's routing usage from its mod, not from here."
- ADD to `claude-fo-dispatch.md` `## Spawn Call (Agent)` (new sub-bullet under the sequencing rule, or a short `## Standing-Teammate Injection (Claude)` section): "The Claude realization of the core's standing-injection call: before the first team-mode dispatch, run `spacedock dispatch spawn-standing-all --workflow-dir {wd} --team {team_name}` and forward each spawn spec in the returned JSON array to `Agent()`. The call is idempotent (already-alive members omitted) and emits `[]` in bare mode or when none is declared. Standing teammates are team-scoped: they die with the team at teardown." (The `## Spawn Call (Agent)` sequencing rule already references `spawn-standing-all` requiring a real `team_name` from a prior `TeamCreate` — that sentence stays.)
- Codex (`codex-first-officer-runtime.md`): one line — "Codex has no shared standing-teammate surface; the standing-injection call is a no-op." Pi (`pi-first-officer-runtime.md`): one line — "`pi-subagents` has no standing surface (no-op); `pi-agent-teams` MAY map injection to `member_spawn` per the adapter's lifecycle mapping."

**Leak 2 — `team_name` / TeamCreate sequencing as a universal dispatch fact** (`fo-dispatch-core.md` `## Dispatch Adapter`, the MANDATORY block).

- BEFORE: "The key fields that MUST come from helper output are `subagent_type`, `name`, `team_name`, `model`, and `prompt` (which contains the completion signal)."
- AFTER (core): "The key fields that MUST come from helper output are `subagent_type`, `name`, `model`, and `prompt` (which contains the completion signal), plus any host-scoped fields the adapter declares (e.g. Claude `team_name`). The adapter names which emitted fields map to its spawn call and which are absent on its host." (The `## Dispatch Adapter` step-4 sentence "Forward `output.model` … when present" stays; it is already host-neutral.)
- `claude-fo-dispatch.md` already owns Team Creation, the TeamCreate-first sequencing rule, and the `Agent(... team_name=output.team_name // omit if bare mode ...)` mapping — no move needed; it is already the source of truth. Codex `codex-first-officer-runtime.md` already states "no `team_name` lifecycle to create, recover, or tear down" — keep. Pi already states "no `team_name` is required" for the first slice — keep.

**Leak 3 — bare-mode blocking-until-completion as a host-neutral invariant** (`fo-dispatch-core.md` `## Dispatch Adapter`, the line after step 5).

- BEFORE: "In bare mode, dispatch blocks until the subagent completes — concurrent dispatch is not possible. Dispatch one entity at a time and process completions inline."
- AFTER (core): "Whether a dispatch blocks until worker completion or returns a handle to await later is the adapter's behavior, not a host-neutral invariant. When the adapter's dispatch is blocking, dispatch one entity at a time and process each completion inline before the next." (Removes the `bare mode`-as-universal framing; the concrete blocking behavior moves to the adapter.)
- ADD to `claude-fo-dispatch.md` `## Degraded Mode (spacedock seams)` (it already owns bare-mode emission): "In Claude bare mode the `Agent()` call blocks until the subagent completes — concurrent dispatch is not possible, so dispatch one entity at a time and process completions inline." Codex: already covered by `## Awaiting Completion` (async final-status notification; `wait_agent` is the explicit foreground wait) — no new prose. Pi: `pi-subagents` result IS the completion (already stated) — no new prose.

**Leak 4 — "first team-mode dispatch" → "first worker dispatch"** (`fo-dispatch-core.md` header line 3, `## Dispatch` step intro line 9 — covered by Leak 1 — and any other self-reference).

- BEFORE (line 3): "Lazily loaded at the first team-mode dispatch (named by the boot-resident core); a greet-and-stop boot never reads it."
- AFTER (line 3): "Lazily loaded at the first worker dispatch (named by the boot-resident core); a greet-and-stop boot never reads it."
- The Claude adapter (`claude-fo-dispatch.md` line 3, `claude-first-officer-runtime.md` line 7) legitimately keeps "first team-mode dispatch" — on Claude the first worker dispatch IS a team-mode dispatch, so the host-specific phrasing is correct there.
- SCOPE DECISION (flag at gate): `first-officer-shared-core.md` line 94 — the boot-resident core's pointer — also says "lazily loaded at the first team-mode dispatch." It is the boot-resident, host-neutral core, so by the same rule it should read "first worker dispatch." Recommend including this one-token fix in scope (it is the same defect in the same host-neutral layer). It is NOT covered by the AC-1 contractlint check (which scans only `fo-dispatch-core.md`); the check would need a second scanned file if we want it guarded. Recommendation: fix the prose in this task, leave the AC-1 check scoped to `fo-dispatch-core.md`, and note the shared-core line as prose-only.

**Leak 5 — reuse condition 1 ("Not in bare mode (teams available)")** (`fo-dispatch-core.md` `## Reuse and Fresh Dispatch`, condition 1).

- BEFORE: "1. Not in bare mode (teams available)."
- AFTER (core): "1. The runtime adapter exposes a live, reusable handle to the completed worker (its reuse-advance handle). When the adapter has no reusable-handle surface, this condition fails and the FO dispatches fresh." (Reframes from the Claude `bare mode` negative to the host-neutral positive: reuse needs a live handle. Codex `send_input` against a bound thread satisfies it; Claude bare mode — no kept-alive teammate — fails it; Pi fresh-by-default fails it unless a substrate exposes a resumable handle.)
- `claude-fo-dispatch.md` `## Feedback Rejection Flow (bare mode)` already states bare mode is sequential / no kept-alive reviewer — keep; that is the Claude realization of "no reusable handle in bare mode." Codex `## Reuse And Feedback Routing` already states `send_input` reuse against an existing worker — keep. Pi `## Follow-up and Reuse` already states fresh-redispatch default — keep.

## Spike determination

The riskiest unverified bit is the AC-1 contractlint check's feasibility as STRUCTURAL-and-NARROW (not a prose-grep tautology). Spiked it in `internal/contractlint/` (throwaway test, removed after; baseline `go test ./internal/contractlint` green before and after).

Result — the mechanism is feasible, and the spike surfaced a load-bearing refinement: a same-line `runtime adapter` mention is NOT a sufficient exemption. The first spike pass exempted any line containing "runtime adapter" and produced a FALSE NEGATIVE on the real current leak — line 9 names the adapter only as the forward destination ("forward each spawn spec … to the runtime adapter's spawn call") while still issuing `run \`spacedock dispatch spawn-standing-all --team {team_name}\`` imperatively. The fix: the exemption must be PHRASE-LEVEL — the adapter is the sentence's subject realizing the call (`Claude adapter`, `the adapter realizes`, `adapter maps`, `Claude realization`), not merely a same-line forward-target word. With that refinement the spike (a) flagged the real current leak, (b) passed a host-neutral paraphrase, (c) passed a genuinely adapter-delegated sentence, (d) passed a clean meaning-inverting paraphrase. This refinement is baked into AC-1 above so implementation does not re-discover it.

Other mechanisms — no spike needed: the AC-2 Codex `send_input` reuse path is already proven by `internal/ensigncycle/shared_reviewer_reuse_test.go` (`assertCodexReviewerReuse` scans the `codex exec --json` transcript for a `send_input collab_tool_call` to the cycle-1 reviewer's bound thread, and explicitly handles the no-team-mode/no-`team_name` Codex case). The AC-3 Claude gates (`spawn-standing-all`, TeamCreate sequencing, terminal marker) are existing passing tests. The prose move itself is plain editing of markdown reference files — no new mechanism.

## Test plan

Cost/complexity: low. One new structural Go test (~60 lines) in the existing quarantine package, plus prose-only edits to four reference files. No fixtures, no live workflow, no CLI surface change.

1. **AC-1 (new, structural).** Add the scanner + discriminator control to `internal/contractlint/structural_checks_test.go` (the quarantine package). Write it FAILING first against the current core (the spike confirmed it flags line 9's `spawn-standing-all --team {team_name}` imperative), then make the prose move and confirm it passes. The positive control (host-coupled fixture flags, host-neutral paraphrase passes) keeps it from passing vacuously.
2. **AC-2 (existing).** `go test ./internal/ensigncycle -run 'SharedReviewerReuse'` — the Codex `send_input`-against-live-handle reuse path stays green after the reuse-condition-1 reframe.
3. **AC-3 (existing Claude gates).** `go test ./internal/contractlint`, `go test ./internal/dispatch -run 'Standing|Build'`, and the Claude standing/teardown focused tests that cover `spawn-standing-all`, TeamCreate sequencing, and the terminal marker.

Run:

```bash
go test ./internal/contractlint
go test ./internal/dispatch -run 'Codex|Pi|Standing|Build'
go test ./internal/ensigncycle -run 'Codex|SharedReviewerReuse|WaitWatchdog|GradeMarkerMatchesContract'
go test ./...
```

(`ForceTeamMode` was dropped from the original test-plan command — no test by that name exists in `internal/ensigncycle`; verified by grep. The remaining named tests all exist.)

## Stage Report: ideation

- DONE: AC-1's contractlint check must be STRUCTURAL and NARROW, never a prose-grep tautology
  AC-1 rewritten: scanner targets the LITERAL tokens `spawn-standing-all` / `--team {team_name}` (presence is the defect, family of `TestRetiredPluginPrivatePathsAbsent`), expected value from the rule not the file's words, with a positive discriminator control and an explicit note that the paraphrasable `team-mode` is deliberately excluded from the token set. Quarantine-package placement named.
- DONE: Pin the exact before/after prose move for the five host-coupled leaks, Codex + Pi mappings, keep the 2y split
  Added `## Concrete prose move` with verbatim BEFORE/AFTER for all five leaks (standing-injection + `spawn-standing-all --team`, `team_name`/TeamCreate sequencing, bare-mode-blocking, "first team-mode dispatch"→"first worker dispatch", reuse-condition-1) plus Codex (`spawn_agent`/`send_input`/`wait_agent`, no `team_name`, no reconcile) and Pi (`subagent`/`pi-agent-teams`) mappings; "do NOT collapse the cores" stated.
- DONE: Record the spike determination
  `## Spike determination` records the AC-1 check spike: feasible, AND surfaced a load-bearing refinement (same-line `runtime adapter` mention is NOT a sufficient exemption — it false-negatived the real leak; exemption must be phrase-level adapter-as-subject). AC-2/AC-3/prose-edit recorded as "no spike needed" with the proven mechanisms.

### Summary

Fleshed the ideation body into a concrete, checkable spec: sharpened AC-1 into a structural-absence token scan with a phrase-level adapter-subject exemption (not a prose-grep), pinned verbatim before/after prose for all five leaks with Codex/Pi adapter mappings, and recorded the spike result. The spike (throwaway, removed; contractlint baseline green) proved the check is feasible and caught that a bare same-line "runtime adapter" word false-negatives the real leak — the exemption must be phrase-level, now baked into AC-1. One scope decision flagged for the gate: `first-officer-shared-core.md:94` carries the same "first team-mode dispatch" leak in the boot-resident host-neutral layer; recommend fixing the prose in-task but leaving the AC-1 check scoped to `fo-dispatch-core.md`.

## Stage Report: implementation

- DONE: Apply the five-leak prose move EXACTLY as pinned in `## Concrete prose move` — host-neutral language stays in `fo-dispatch-core.md`, the Claude realization moves to / stays in `claude-fo-dispatch.md`, Codex + Pi state their one-line mappings; keep the 2y split; apply the gate-approved Leak-4 scope add fixing `first-officer-shared-core.md:94`.
  Code commit 90089789. fo-dispatch-core.md: standing-injection now defers to the adapter (Leak 1), `team_name` is adapter-declared not universal (Leak 2), bare-mode blocking reframed as adapter behavior (Leak 3), "first team-mode dispatch"→"first worker dispatch" on the header (Leak 4), reuse-cond-1 reframed to "adapter exposes a live reusable handle" (Leak 5). claude-fo-dispatch.md gained `## Standing-Teammate Injection (Claude)` + a bare-mode-blocking Degraded-Mode bullet. codex-/pi-first-officer-runtime.md each gained the one-line standing-injection mapping. first-officer-shared-core.md:94 fixed to "first worker dispatch" (cores NOT collapsed; Claude layer keeps "first team-mode dispatch" at claude-fo-dispatch.md:3/:7/:55 and claude-first-officer-runtime.md:7).
- DONE: Implement AC-1's structural contractlint check in `internal/contractlint/` exactly as designed — literal-token scan for `spawn-standing-all` and `--team {team_name}`, phrase-level adapter-as-subject exemption (not a bare same-line "runtime adapter" word), positive discriminator control, written FAILING first then greened by the prose move; do NOT add `team-mode` to the token set.
  Added `TestDispatchCoreHasNoClaudeTeamImperative` + `TestDispatchCoreClaudeTeamScannerDiscriminates` to structural_checks_test.go. Confirmed red against the current core (flagged line 9's `spawn-standing-all --team {team_name}` imperative), then green after the move. Discriminator proves: host-coupled flags, bare "runtime adapter" word still flags, host-neutral/adapter-as-subject/inverted paraphrases pass. `team-mode` deliberately excluded from the token set. Adversarial mutation re-confirmed the check is non-vacuous on the real file (re-introduced leak FAILs, revert PASSes).
- DONE: `go test ./internal/contractlint`, `go test ./internal/dispatch -run 'Codex|Pi|Standing|Build'`, `go test ./internal/ensigncycle -run 'Codex|SharedReviewerReuse|WaitWatchdog|GradeMarkerMatchesContract'`, and `go test ./...` all green — AC-2 (Codex `send_input` reuse stays valid after the reuse-condition-1 reframe) and AC-3 (Claude behavior preserved) hold.
  All four green: contractlint ok, dispatch ok (6.5s), ensigncycle ok (1.2s, SharedReviewerReuse confirms AC-2), full `go test ./...` ok across all 16 packages. AC-3 Claude gates (spawn-standing-all, TeamCreate sequencing, terminal marker via dispatch Standing/Build + ensigncycle) pass.

### Summary

Made the host-neutral dispatch core genuinely runtime-neutral via the pinned five-leak prose move (cores NOT collapsed): the core now defers Claude-team commands, `team_name`, bare-mode blocking, and the reuse-handle condition to the runtime adapter, with the Claude realizations re-homed in claude-fo-dispatch.md and one-line Codex/Pi mappings added. AC-1 is enforced by a new structural-absence contractlint check (literal-token scan, phrase-level adapter-as-subject exemption, discriminator control) written failing-first per TDD and greened by the move; an adversarial mutation confirmed it catches a re-introduced leak on the real file. All focused gates plus the full suite are green — Codex `send_input` reuse (AC-2) survives the reuse-condition-1 reframe and Claude behavior (AC-3) is preserved.
