---
title: Generalize the Cowork binary bootstrap from survey into the shared launcher gate
status: validation
source: captain request, live Cowork dogfood 2026-07-13
started: 2026-07-18T02:23:50Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-cowork-bootstrap-shared-launcher-gate
issue:
id: 8v5w1fk28m5bssm4147ad0ae
---

Make every spacedock skill front door work on a fresh Cowork project, not just survey. Today only survey (per survey-claude-cowork-runtime-detection, id eqn) owns the Cowork consent/install lifecycle for `.spacedock/bin/spacedock`; any other entry point that reaches the binary version gate on a Cowork VM aborts with host-wrong remediation.

## Problem

The FO/commission startup contract's binary version gate (`skills/first-officer/references/first-officer-shared-core.md`, step 1 "Binary version gate", the binary-absent class) aborts a missing binary with `brew install spacedock-dev/homebrew-tap/spacedock` and "launch with `spacedock claude`". Both are dead inside the Cowork VM. Commission does not carry its own gate — it instructs the agent to "Execute the first-officer startup procedure directly", so every front door funnels through this one abort. A fresh-project Cowork user who goes straight to `/spacedock:commission` (or launches the FO) designs or resumes a workflow whose first boot then dead-ends there.

The working bootstrap — positive Cowork capability probe → selected-folder two-surface binding → consent/Full-network relaunch → checksum-verified session install → exclusive-create of the exact mounted `.spacedock/bin/spacedock` → dual-surface verify — is live-proven (eqn run 0/1, 2026-07-12/13) but is specified only inside survey's entity body, reachable by no other front door.

**First-party evidence (this ideation ran inside a live Cowork FO→ensign dispatch; see `## Spike` below):** `spacedock` is not on PATH and there is no `brew` in the Cowork VM, so the brew/curl/`spacedock claude` remediation is not merely wrong, it is unactionable. The exact mounted `.spacedock/bin/spacedock` resolves and reports `spacedock 0.25.0`. All three host env markers (`CLAUDECODE`, `CODEX_THREAD_ID`, `PI_CODING_AGENT_DIR`) are unset — which is why `spacedock dispatch build` in this very session failed host auto-detection and needed explicit `--host claude`.

## Positive marker and host-detection scope decision

The Cowork arm must fire from a **positive capability probe**, never inferred from absence — eqn's proven principle (a missing PATH binary on a plain Linux box must not false-positive as Cowork). The reused positive marker is the same one eqn's survey branch uses: `ToolSearch(query="select:mcp__session_info__list_sessions,mcp__session_info__read_transcript", max_results=2)` resolving both exact schemas. This is a **tool-surface** probe available to the FO/commission agent at the gate moment regardless of shell state — confirmed live this session, and it keeps working when the shell VM is down (`COWORK_SHELL_UNAVAILABLE`), which an env/mount-shape probe could not.

**Host detection for dispatch is OUT of scope, named explicitly (per the dispatching FO's request).** The observed `dispatch build --host` auto-detection failure is a *binary-side* concern: the Go launcher infers host from `CLAUDECODE`/`CODEX_THREAD_ID`, which are unset in Cowork. The skill version-gate arm here does not depend on that binary detection — it uses the tool-surface probe above. Teaching the binary (or the FO helper calls) a positive Cowork host marker is a separate change to `internal/` host detection with its own proof surface; it is filed as a follow-up candidate (`cowork-host-detection-for-dispatch`), not absorbed here. The three checklist items concern only the skill gate arm, the mounted-path resolution, and dedup.

## Proposed approach

Extract eqn's proven lifecycle into one shipped shared reference and point the gate at it:

1. **Create the canonical shared reference** `skills/first-officer/references/cowork-bootstrap.md`, seeded verbatim from eqn's live-proven lifecycle: positive probe, selected-folder two-surface binding, the `COWORK_INSTALL_ABSENT`/`COWORK_INSTALL_BROKEN` classes, the consent/Full-network relaunch prompt, checksum session install, parent-first validation, the single noclobber `( set -C; cat … )` exclusive create, dual-surface verify, and the `COWORK_SHELL_UNAVAILABLE` stop class. It lives in the FO reference dir because the version gate (its primary consumer) already lives there.
2. **Version-gate Cowork arm.** In the binary-absent class, before emitting the brew/`spacedock claude` hint, run the positive Cowork probe. On a positive probe, deferred-`Read` the shared reference (sibling-relative from the skill base) and follow it — consent-gated bootstrap, not the host hint. On a negative probe, the existing host remediation is emitted unchanged. The reference is read on-demand only in this rare class, not made boot-resident.
3. **`SPACEDOCK_BIN` resolution learns the mounted path.** On a positive Cowork probe, resolve `<selected-folder>/.spacedock/bin/spacedock` first, with eqn's existing-state classes exactly: silent-`--version` success reuses with zero network; non-executable/dangling stops without overwrite; absent enters the consent flow. Off Cowork, `${SPACEDOCK_BIN:-spacedock}` resolution is unchanged.
4. **Single-source, no drift.** Survey's Cowork branch (shipped by eqn's implementation) references the same file instead of embedding a copy; whichever of {this task, eqn's survey implementation} ships the embedded lifecycle second converts it to a pointer. A `internal/contractlint/` test enforces the invariant so it cannot regress.

### New mechanisms — justification (value AC served / simplest alternative / why insufficient)

- **Positive Cowork tool-probe at the gate** → serves AC-1. Alternative: infer Cowork from `spacedock` absent + Linux + `/sessions/*/mnt` layout. Insufficient: violates eqn's "positive probe, never infer from absence" — a non-Cowork Linux box with no `spacedock` would false-positive into a bogus consent flow; and an env/mount probe is unavailable during `COWORK_SHELL_UNAVAILABLE`. The tool probe resolves live (spike below).
- **Mounted-path `SPACEDOCK_BIN` resolution** → serves AC-2. Alternative: rely on PATH or `$HOME/.local/bin/spacedock`. Insufficient: eqn live-falsified those — Cowork HOME is per-session scratch that resets, and PATH can be shadowed; the mounted `.spacedock/bin/spacedock` is the only persistent execution identity.
- **Canonical reference file + `contractlint` duplication test** → serves AC-3. Alternative A: duplicate the lifecycle prose into each front door and hand-sync. Insufficient: prose drift between front doors is the exact failure this task exists to remove; a partial edit to one copy silently diverges. Alternative B: a shared `@`-included physical file. Insufficient: no cross-skill `@`-include mechanism exists today (verified — `@`-includes are skill-relative), so this would be an unproven skill-loader change; a cross-skill on-demand `Read` against the packaging-guaranteed sibling layout (precedent: commission reads `references/decomposed-snippets.md` by path) is smaller and testable. The lint measures the count so drift cannot creep back.

## Spike (first-party, executed this session — the riskiest unverified mechanism)

The one mechanism not already proven by eqn is **whether a positive Cowork capability is observable at the FO/commission version-gate moment** (eqn only proved the probe inside survey). Exercised live in this Cowork FO→ensign dispatch:

- `ToolSearch(select:mcp__session_info__list_sessions,mcp__session_info__read_transcript)` returned **both** exact schemas — the positive marker resolves at the agent surface an FO/commission boot shares. Absent-schema is not evidence of absence (lazy load); the ToolSearch call is the disambiguator.
- `CLAUDECODE`/`CODEX_THREAD_ID`/`PI_CODING_AGENT_DIR` all unset; `command -v spacedock` → not found; `command -v brew` → not found; `<mount>/.spacedock/bin/spacedock --version` → `spacedock 0.25.0`.

Result: the gate arm's detection is proven available at the gate point, and the brew/launcher hint is proven unactionable in the target runtime. No further spike needed for the remaining mechanisms — they are eqn-proven (bootstrap run 0/1 live 2026-07-12/13; run 2 reuse tracked by eqn) or checkable at fixture/lint level (routing prose, duplication count, sibling-path resolution).

## Out of scope

- Changing eqn's design or re-proving the bootstrap lifecycle (run 0/1 live-proven; run 2 reuse tracked by eqn — AC-2 rides eqn's run-2 rather than re-proving).
- Binary-side / dispatch host auto-detection on Cowork (the `--host claude` gap) — filed as follow-up `cowork-host-detection-for-dispatch`; see the scope decision above.
- Git credential/push handling from the VM and state-checkout gitdir portability (filed separately as state-checkout-portable-gitdir).
- Non-Cowork runtimes' version-gate behavior (the negative-probe path is unchanged and is the control, not a target).

## Acceptance criteria

**AC-1 — On a positive-Cowork binary-absent gate, the number of front doors that terminate at the brew/`spacedock claude` remediation is 0 (today: all of them); the negative-probe control still emits that remediation unchanged.**
This measures the reason-for-existing against a baseline that can regress: a leak of the arm into non-Cowork drops the control's remediation (control fails); a regression that re-hardcodes the host hint on Cowork re-raises the count above 0. Verified by: a fixture replay of the version gate's binary-absent class driving both probe outcomes — the positive arm's output contains the consent/Full-network prompt and contains neither `brew install` nor `spacedock claude`; the negative arm's output still contains both. Because commission executes the FO startup procedure, exercising the single FO gate covers both front doors; the replay asserts commission carries no independent copy of the remediation.

**AC-2 — After bootstrap, a second front-door invocation resolves the exact mounted `.spacedock/bin/spacedock` with 0 network/install events.**
Verified by: a second replay invocation observing reuse — the mounted-path `--version` succeeds and the network/install event count is 0, matching eqn's working existing-state class; broken (non-executable/dangling) still stops without overwrite; absent still enters consent. Reuses eqn's committed Cowork event fixtures; live reuse rides eqn's run-2 session.

**AC-3 — The Cowork bootstrap lifecycle body appears in exactly one shipped file; every front door carries a pointer, not a copy.**
This is the value metric for "specified once": the count of shipped skill/reference files containing the load-bearing lifecycle tokens (`COWORK_INSTALL_ABSENT`, the noclobber `set -C` exclusive-create line, `COWORK_SHELL_UNAVAILABLE`, the selected-folder two-surface binding) is 1 — a baseline that moves the wrong way the moment any front door pastes a second copy. Verified by: a new `internal/contractlint/` test asserting (a) exactly one shipped file contains those tokens, (b) the FO version-gate binary-absent class and survey's Cowork branch each reference that file's path rather than restating the lifecycle, and (c) the referenced sibling path resolves to an existing file.

## Test plan

- **Gate-arm routing replay (cheap, primary):** fixture replay of the binary-absent class over both probe outcomes for AC-1 (positive → consent prompt, no brew/launcher; negative → host hint unchanged), extending the existing gate/version-gate fixtures. No live run needed for routing.
- **Reuse replay (cheap):** second-invocation replay for AC-2 over eqn's committed Cowork event fixtures — reuse with 0 network/install events; broken/absent controls.
- **Dedup contractlint (cheap):** new `internal/contractlint/` test for AC-3 — token-count == 1, front-door pointers present, sibling path resolves. Runs in the existing `go test ./internal/contractlint/...` lane.
- **Live Cowork smoke (rides eqn's run-2):** one FO/commission boot on a fresh Cowork project observes the binary-absent gate route to consent (not brew) and, post-relaunch, reuse of the mounted binary with 0 network events. No separate live lane beyond eqn's run-2 session.
- **Docs render (low):** build the docs site and inspect the changed `install.md` Cowork note.

## Documentation change proposed for implementation

`docs/site/get-started/install.md` currently sends every reader to brew/curl install and `spacedock claude` launch, with no Cowork path — a fresh Cowork user has neither brew nor a persistent PATH. Add a Cowork subsection after the `## Launch` block (eqn already owns the parallel note in `survey.md`):

Before — after the "Replace `claude` with `codex` or `pi`…" line, no Cowork guidance.

After — insert:

> ### Claude Cowork
>
> In Claude Cowork there is no Homebrew and no persistent shell PATH. Instead, set **Full network access** in Cowork Settings, relaunch the session, and run any Spacedock front door (`/spacedock:survey`, `/spacedock:commission`, or the first officer). On first use Spacedock bootstraps a checksum-verified helper into `.spacedock/bin/spacedock` in your selected project folder; later sessions reuse it with no download. You do not run `brew install` or `spacedock claude` inside Cowork.

Implementation applies this diff; the ideation gate reviews it.

## Stage Report: ideation

- DONE: Version-gate binary-absent class gains a Cowork arm: positive capability probe routes to the shared consent/install flow; negative probe keeps today's brew/launcher remediation unchanged.
  Proposed approach step 2 + AC-1: the binary-absent class runs the positive `ToolSearch(session_info)` probe; positive → deferred-Read the shared reference and follow the consent flow, negative → existing host hint verbatim. AC-1 measures dead-ending front doors == 0 with the negative-probe control as the baseline that can regress.
- DONE: SPACEDOCK_BIN resolution names the exact mounted .spacedock/bin/spacedock first on Cowork, preserving eqn's existing-state classes (working reuses zero-network, broken stops without overwrite, absent enters consent).
  Approach step 3 + AC-2: on a positive probe resolve `<selected-folder>/.spacedock/bin/spacedock` first with eqn's working/broken/absent classes; off-Cowork `${SPACEDOCK_BIN:-spacedock}` unchanged. Live-confirmed the mounted path reports `spacedock 0.25.0`.
- DONE: The bootstrap lifecycle is specified exactly once: survey and the launcher gate reference one shared flow (dedup check named), including COWORK_SHELL_UNAVAILABLE and the two-surface binding.
  Approach step 1/4 + AC-3: a new canonical `skills/first-officer/references/cowork-bootstrap.md` seeded from eqn's proven lifecycle; a new `internal/contractlint/` test asserts the lifecycle tokens (incl. `COWORK_SHELL_UNAVAILABLE`, the two-surface binding) appear in exactly one shipped file and both front doors carry a pointer.
- DONE: Spike the riskiest unverified mechanism.
  Ran live in this Cowork FO→ensign dispatch: the positive `session_info` probe resolves at the agent surface the gate shares; host env markers all unset; `spacedock`/`brew` absent; mounted binary works. Recorded in `## Spike`. Remaining mechanisms are eqn-proven or fixture/lint-checkable.
- DONE: Propose the documentation change (user-visible behavior).
  Concrete before/after doc diff for `docs/site/get-started/install.md` (new `### Claude Cowork` subsection) recorded in the task body for implementation to apply.

### Summary

Sharpened the ideation into a gate-ready design grounded in first-party evidence gathered live inside this Cowork FO→ensign dispatch. The version gate is already single-sourced in the FO shared core (commission executes the FO startup procedure), so the Cowork arm is one edit reaching every front door; the canonical bootstrap lifecycle becomes a new shared reference seeded from eqn's proven design, with a contractlint duplication test enforcing "specified once". Two scope decisions were made explicit at the dispatching FO's request: the positive Cowork marker reuses eqn's tool-surface `session_info` probe (available even during `COWORK_SHELL_UNAVAILABLE`), and binary-side dispatch host auto-detection (the `--host claude` gap) is named out of scope as a separate follow-up. Open coordination for the FO/gate: AC-3 assumes eqn's in-flight survey Cowork branch points at the same shared file rather than embedding a copy — whichever ships the embedded lifecycle second converts it to a pointer, and the lint holds the invariant either way.

## Stage Report: implementation

- DONE: Positive-Cowork gate-arm replay green: consent/Full-network prompt present, zero brew/`spacedock claude` tokens; negative-probe control emits today's remediation unchanged (AC-1).
  `TestCoworkGateArmRoutingReplay` (internal/contractlint/cowork_bootstrap_gate_test.go) drives both probe outcomes over the shipped binary-absent class: order pinned probe→pointer→negative→remediation; positive arm renders the canonical consent prompt (Full network access present, both remediation tokens absent, including in the class's positive segment); negative arm still carries `brew install spacedock-dev/homebrew-tap/spacedock` and `spacedock claude`; commission verified to carry no independent remediation copy. PASS.
- DONE: Canonical skills/first-officer/references/cowork-bootstrap.md ships seeded from eqn's lifecycle; contractlint token-count==1 test green with front-door pointers resolving (AC-3).
  New canonical reference seeded from eqn's live-proven lifecycle (positive probe, selected-folder two-surface binding, COWORK_INSTALL_ABSENT/BROKEN classes, consent/relaunch, checksum session install, noclobber exclusive create, dual-surface verify, COWORK_SHELL_UNAVAILABLE). `TestCoworkLifecycleShipsExactlyOnce` (5 tokens, count==1) and `TestCoworkFrontDoorsPointAtCanonicalLifecycle` (gate pointer `{first_officer_base}/references/cowork-bootstrap.md`, survey pointer `../first-officer/references/cowork-bootstrap.md`, both resolving non-empty) PASS. Survey's pointer shipped here (this task ships the lifecycle first); the existing reference-closure lint learned the cross-skill `../<skill>/references/` form so the sibling pointer closure-checks like a local one.
- DONE: Reuse replay green over committed Cowork event fixtures (mounted-leaf reuse, 0 network/install; broken/absent controls) and the install.md Cowork doc diff applied (AC-2, docs).
  `TestCoworkBootstrapReuseReplay` (skills/integration) EXECUTES the shipped probe block extracted from the canonical file over committed sanitized fixtures (skills/integration/testdata/cowork/bootstrap-events.jsonl, event classes mirroring eqn's live 2026-07-12/13 runs): mounted-leaf reuse silent across two invocations with 0 instrumented network/install events and leaf `--version` resolving; nonexecutable + dangling-symlink → COWORK_INSTALL_BROKEN with byte/mode-identical leaf (no overwrite); absent → COWORK_INSTALL_ABSENT with 0 network events. 4/4 PASS. install.md gained the `### Claude Cowork` subsection per the ideation diff; mkdocs build green, subsection renders on the built page.

### Summary

Delivered in worktree commit d22acf0b: the FO shared-core binary-absent class now runs the positive session_info ToolSearch probe before any remediation — positive routes to the new on-demand canonical `cowork-bootstrap.md` (mounted `.spacedock/bin/spacedock` resolved first as SPACEDOCK_BIN, else consent-gated bootstrap), negative emits today's brew/launcher hint verbatim; commission funnels through the same gate so one edit covers every front door. Environment notes: the Cowork VM shipped no Go toolchain — installed go1.22.12 (linux/arm64) into session scratch; full `go test ./...` green (15 packages), `go vet` clean, docs built with pip-installed mkdocs-material. Live run-2 Cowork smoke deliberately not run here — it rides eqn's run-2 session per the test plan. Follow-up unchanged: binary-side dispatch host detection stays out of scope (cowork-host-detection-for-dispatch).

## Stage Report: validation

- DONE: Independently re-run the shipped proof: contractlint gate/dedup tests and the integration reuse replay green from the worktree, not taken from the maker's report.
  Fresh runs with `-count=1` at d22acf0b: TestCoworkLifecycleShipsExactlyOnce, TestCoworkFrontDoorsPointAtCanonicalLifecycle, TestCoworkGateArmRoutingReplay PASS; TestCoworkBootstrapReuseReplay 4/4 scenarios PASS; full suite all 15 packages ok (incl. the modified reference-closure regex lane); mkdocs build green and the `Claude Cowork` subsection renders on the built install page.
- DONE: AC-1/AC-2/AC-3 verified against the actual diff (shared-core gate arm, canonical reference, survey pointer, install.md) — evidence cited per AC, any unmet AC named.
  AC-1: diff to `first-officer-shared-core.md` runs the positive session_info ToolSearch probe before any remediation; positive arm routes to `{first_officer_base}/references/cowork-bootstrap.md` with no brew/`spacedock claude` token, negative arm keeps both verbatim; commission (SKILL.md:657) executes the FO startup procedure and carries no independent remediation copy — routing replay re-run green. AC-2: canonical reference resolves the exact mounted `.spacedock/bin/spacedock` first as SPACEDOCK_BIN with eqn's working/broken/absent classes; the reuse replay EXECUTES the shipped probe block over committed fixtures — 2-invocation reuse with 0 instrumented network/install events and leaf `--version` resolving; nonexecutable/dangling → COWORK_INSTALL_BROKEN with byte/mode-identical leaf; absent → COWORK_INSTALL_ABSENT, 0 events. AC-3: 5 load-bearing tokens ship in exactly one file (the canonical reference); gate and survey pointers both resolve; the extended cross-skill ref-closure lint passes repo-wide. install.md diff matches the ideation before/after exactly. No unmet AC.
- DONE: Adversarial pass on the dedup invariant: attempt a second-copy/token-drift mutation and confirm the lint fails it (the code gate, not the prose, holds the guarantee).
  Four mutations, all caught: (a) second copy of COWORK_INSTALL_ABSENT pasted into commission/SKILL.md → dedup fails "ships in 2 file(s)"; (b) noclobber exclusive-create line drifted to `set +C; cp` → dedup fails "ships in 0 file(s)"; (c) `brew install` re-hardcoded into the positive gate arm → routing replay fails the probe→pointer→negative→remediation order; (d) semantic swap of BROKEN/ABSENT echoes inside the shipped probe block (pinned first line intact) → reuse replay fails behaviorally in 3 scenarios, proving the replay executes the shipped block rather than pattern-matching it. Probe-block first-line drift also fails loud at extraction. All mutations reverted; worktree clean at d22acf0b.

### Summary

PASSED. Verdict grounds: every proof re-run fresh from the worktree (contractlint, integration replay, full 15-package suite, docs build), each AC traced to the actual diff, and the dedup/routing invariants shown to be held by failing code gates, not prose. Deferred risks, none material: (1) the live Cowork run-2 smoke deliberately rides eqn's session per the approved test plan — promote to material if eqn's run-2 never executes or fails; (2) AC-1's gate-arm proof is structural-over-prose, inherent to skill contracts — mitigated by the live spike, the negative-probe control, and the two independent files that can diverge; (3) eqn's in-flight survey branch must land as a pointer — the dedup lint will fail an embedded copy at their merge, which is the intended enforcement, but the FO should expect that red test in eqn's lane if they ship a copy. Environment note for reproducers: the VM has no system Go; toolchain at `/sessions/brave-compassionate-ramanujan/go-root/go/bin/go` (go1.22.12 linux/arm64) was reused.
