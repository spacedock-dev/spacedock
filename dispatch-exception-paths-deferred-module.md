---
title: "Exception paths (Degraded Mode, Break-Glass, dead-ensign detail) split out of claude-fo-dispatch.md into failure-triggered deferred modules"
status: ideation
group: tooling
source: "fable-token-trim-scout analysis 2026-07-02: Degraded Mode (claude-fo-dispatch.md:87-117, ~780 tok), the Break-Glass manual-dispatch template (:42-52, ~520), and Context-Budget/Dead-Ensign detail (:128-149, ~650) load at first dispatch but fire only on failure (second dispatch failure, helper non-zero exit, budget-fail) — ~1.6-1.9k tok on every dispatching session's happy path. #457 just shipped the same non-user-invocable-skill pattern for adapter-less deferred modules."
id: 41cfak9bgwtpa01m1z4qkprq
started: 2026-07-02T03:02:51Z
---

## Problem
The Claude dispatch module front-loads its exception machinery: every session that dispatches even one worker carries the full Degraded Mode contract, the Break-Glass template, and the dead-ensign/budget-failure detail, though the happy path never fires them.

Measured at ideation (2026-07-02, `main` @ 59251f18): `skills/first-officer/references/claude-fo-dispatch.md` is 22,992 B total; the three exception spans are the Break-Glass block (lines 42-52, 2,082 B), `## Degraded Mode` (lines 87-117, 3,131 B), and `## Context Budget and Dead Ensign Handling` (lines 128-149, 2,607 B) — 7,820 B (~1.9k tok) loaded at every first dispatch, fired only on failure.

## Proposed approach

Move the three exception bodies into ONE non-user-invocable skill, `skills/fo-dispatch-recovery/SKILL.md`, loaded at its failure triggers. This is the #457 pattern followed exactly — shipped precedents: `skills/fo-status-viewer/SKILL.md` and `skills/fo-write-core/SKILL.md` (commit 59251f18, #457), and the earlier `skills/present-gate/SKILL.md`, `skills/feedback-rejection-flow/SKILL.md`, `skills/using-legacy-claude-team/SKILL.md`. Frontmatter: `user-invocable: false`; the `description` carries the when-to-load (the three failure triggers), mirroring fo-status-viewer's.

**Why one skill, not three.** Every SKILL.md description is resident in the available-skills list of EVERY session running the plugin (FO, ensigns, unrelated sessions) — three skills triple that plugin-wide overhead for paths that are rare and correlated (a helper failure can cascade into degraded mode). One load at failure time brings in ~1.9k tok of deferred context exactly when it is needed; the finer granularity of three skills buys nothing on the happy path. Multi-anchor skills are the shipped shape (fo-status-viewer carries three section anchors).

**Skill sections (bodies move verbatim, per #457's "section bodies verbatim"):**
- `## Degraded Mode` — the Triggers rationale, Effects invariants, the verbatim Captain Report Template, and the Cooperative Shutdown Sweep (including the feedback-cycle exemption).
- `## Break-Glass Manual Dispatch` — the full manual `Agent()` template block plus the "concrete Claude form of fo-dispatch-core.md's Break-Glass template" note.
- `## Context Budget Failure and Dead Ensign Handling` — the three budget-unavailable stderr conditions, the Recovery clause text, the four dead-ensign bullets, and the forward-family context-window note.

**What stays resident in claude-fo-dispatch.md** (the named reliability risk: exception time is when the FO is least reliable, so each resident trigger carries the FIRST ACTION, never a bare pointer; the pre-invocation guards stay resident):
- `## Awaiting Completion` — untouched (guards the common wait path).
- The context-budget probe invocation and the fail-safe rule — reuse-condition-0 runs on the happy path.
- The Model-to-context mapping paragraph — reuse-condition-4's comparator and the canonical model enum (values not quoted here: the `wc` entity amends the enum before this task merges; take the landed text verbatim); `fo-dispatch-core.md:93` names "the Claude enum in `claude-fo-dispatch.md`" and reuse-condition-4 runs on every happy-path reuse decision, so the enum cannot move to a failure-time skill. Post-split the enum appears once (today it is duplicated at lines 52 and 140).
- `## Terminal Worker Teardown`, `## Feedback Rejection Flow (bare mode)`, `## Event Loop`, the Spawn Call, and the Worker Back-Channel — untouched.

### Resident trigger lines: before/after diffs

**Diff 1 — Break-Glass (replaces lines 42-52, 2,082 B → 444 B measured).** Before: the full bolded intro + 8-line `Agent()` template + closing note. Note: this span directly overlaps the `jh` entity's break-glass demotion, which merges first — re-derive this diff against jh's landed text (see the merge-order dependency below). After:

> **Break-Glass Manual Dispatch (ONLY when `spacedock dispatch build` exits non-zero or is unavailable):** never hand-assemble a dispatch while the helper works. First action: report the helper failure (command, exit code, stderr) to the captain. Then `Skill(skill="spacedock:fo-dispatch-recovery")` and fill its `## Break-Glass Manual Dispatch` template; the conditional `model=` slot draws from the canonical enum in `## Context Budget` below.

**Diff 2 — Degraded Mode (replaces lines 87-117, 3,131 B → ~715 B with the pre-trip-names clause; 660 B measured without it).** Before: `## Degraded Mode` + `### Triggers` + `### Effects` + `### Captain Report Template` + `### Cooperative Shutdown Sweep`. After:

> ## Degraded Mode (trigger)
>
> Any ONE of these trips Degraded Mode — a session-wide, irreversible fallback to sequential bare dispatch: (1) any SECOND dispatch failure within the session (no time window, no counter — deliberate: the FO cannot reliably track failure timestamps); (2) the captain command `/spacedock bare`; (3) `Agent` or `SendMessage` themselves unavailable. On trip, FIRST ACTION: stop all named/background dispatch for the remainder of the session and do not route `SendMessage` to any pre-trip worker name. Then `Skill(skill="spacedock:fo-dispatch-recovery")` and follow its `## Degraded Mode` section: the bare-dispatch invariants, the verbatim captain report, and the cooperative shutdown sweep.

The "do not route `SendMessage` to any pre-trip worker name" clause is a deliberate ~55 B spend against the AC-1 margin: the no-reuse-of-pre-degrade-names prohibition otherwise lives only in the moved Effects, and a failed skill load must not leave the FO free to advance work to a presumed-zombified name.

**Diff 3 — Context Budget (replaces lines 128-149, 2,607 B → ~1,350 B measured).** Before: probe check + three stderr bullets + mapping paragraph + recovery clause + four dead-ensign bullets. After:

> ## Context Budget
>
> This is the Claude realization of reuse-condition-0's budget probe (and the feedback-rejection budget check); Codex declares none.
>
> **Context budget check:** Run `${SPACEDOCK_BIN:-spacedock} dispatch context-budget --name {ensign-name}`. Parse the JSON output. `reuse_ok: true` → reuse may proceed. `reuse_ok: false`, or ANY non-zero exit with no reading, is fail-safe — never silent-reuse: log to captain and fresh-dispatch. Before the replacement dispatch (budget-fail, zombie, or dead ensign), `Skill(skill="spacedock:fo-dispatch-recovery")` — its `## Context Budget Failure and Dead Ensign Handling` section carries the recovery clause for the prior worktree and the dead-ensign rules.
>
> **Model-to-context mapping:** [the line-140 paragraph retained verbatim as landed at merge time — reuse-condition-4's comparator, the canonical enum (do NOT re-quote its values here: `wc` amends the enum and merges first), the fallback re-stamp rule]

### Why each load point is safe at failure time

- **Degraded Mode:** its triggers are team-substrate failures (`Agent`/`SendMessage` down, repeated dispatch failure) — but `Skill()` is a harness-native local file read with no dependency on the team substrate, so the load path is disjoint from the failure that fired it. The safety-critical invariants (stop background dispatch, session-wide, irreversible, no messaging pre-trip worker names) are on the resident trigger line itself, and the host-neutral degraded effects are ALSO resident via `fo-dispatch-core.md`'s `## Dispatch Adapter` (the "no `team_name` on subsequent dispatch" effect and the blocking-dispatch clause, which claude-fo-dispatch.md:103-104 cites today): if the skill load somehow failed, the FO has already halted team dispatch and degrades to halt-and-report, never to continued wrong dispatch or to routing an advance to a presumed-zombified name.
- **Break-Glass:** the trigger is a helper-binary failure (`spacedock dispatch build` non-zero), orthogonal to skill loading. The first action (report to the captain) needs no loaded body; a failed load leaves the FO already talking to the captain with the failure surfaced.
- **Budget-fail/dead-ensign:** the trigger is a probe-binary failure or reading, orthogonal to skill loading. The fail-safe decision (fresh-dispatch, never silent-reuse) is resident (and restated host-neutrally at `fo-dispatch-core.md:34`); the skill supplies only the recovery-clause text and dead-ensign bookkeeping, so a failed load cannot cause silent-reuse.

### Registry, re-points, and doc diff

**`## Deferred load points` registry** (`skills/first-officer/references/first-officer-shared-core.md:35-42`) gains a fifth row (additive — the fold test's four tokens and greet-guard clause survive):

> - `Skill(skill="spacedock:fo-dispatch-recovery")` — dispatch failure recovery (Degraded Mode, break-glass manual dispatch, budget-fail/dead-ensign handling); named at its triggers inside the Claude dispatch module — no boot and no happy-path dispatch loads it.

This row is also what keeps `TestDeferredSkillCoresResolveAndCarryCeremony` (which hardcodes the shared core as sole namer) working unmodified once the skill joins `deferredSkillCores`.

**Re-points off the moved sections:**
- `skills/using-legacy-claude-team/SKILL.md:25` — before: "…before `claude-fo-dispatch.md`'s `## Degraded Mode` fires." → after: "…before the Degraded Mode trigger in `claude-fo-dispatch.md` fires (body in `Skill(skill="spacedock:fo-dispatch-recovery")`)."
- `skills/using-legacy-claude-team/SKILL.md:49` — before: "**Fall back to Degraded Mode** per `## Degraded Mode` in `claude-fo-dispatch.md`." → after: "**Fall back to Degraded Mode** per its trigger in `claude-fo-dispatch.md` — first action there, body via `Skill(skill="spacedock:fo-dispatch-recovery")`."
- `skills/first-officer/references/claude-first-officer-runtime.md:7` — before: "…the Awaiting-Completion idle guardrail, Degraded Mode, the Context-Budget probe, and the Event-Loop reconcile sweep + Backstop — live in `references/claude-fo-dispatch.md`…" → after: "…the Awaiting-Completion idle guardrail, the Degraded-Mode/break-glass/budget-failure trigger lines, the Context-Budget probe, and the Event-Loop reconcile sweep + Backstop — live in `references/claude-fo-dispatch.md` … (the exception bodies behind those triggers load at failure time via `Skill(skill="spacedock:fo-dispatch-recovery")`)."
- `claude-fo-dispatch.md:3` intro — "the idle and degraded-mode guardrails" → "the idle guardrail and the failure-recovery trigger lines".
- `fo-dispatch-core.md:133` ("use the adapter's Break-Glass Manual Dispatch template") — deliberately NOT edited: the host-neutral core must not name a Claude-only skill, and the pointer resolves transitively — the adapter's resident trigger line IS the adapter's break-glass entry (first action + skill load), so "the adapter's template" still reaches the template through it.
- `skills/using-legacy-claude-team/SKILL.md:50` and `:52` — checked and valid as-is: both mention Degraded Mode as a state ("cannot be entered" / "is entered"), not as a pointer at the moved section body; no re-point needed.

Blast radius checked clean at review: beyond the sites listed above, no other contract or test file references the moved section bodies.

### Merge-order dependency (recorded for the FO, which owns enforcement)

Merge order is `wc` → `jh` → this entity (41cfak9bgw). Consequences honored in this body: Diff 1 overlaps jh's break-glass demotion, so the implementation re-derives Diff 1 (and the moved break-glass body) against jh's landed text; Diff 3's bracket note does not quote the canonical enum's values because wc amends the enum; the pre-split baseline captures (equivalence protocol below) and the AC-1 measurement run AFTER jh lands — jh grows the moved span, which improves AC-1's margin.

**Doc diff (`docs/runtime-support.md:23`, user-visible on the docs site)** — two edits to the discriminator paragraph. First, the lead clause accommodates the third case — before: "A deferred module realizes as one of two shapes keyed on WHERE its host binding lives, not on whether it loads an adapter section alongside at the trigger." → after: "A deferred module realizes as a core reference file or a non-user-invocable skill, keyed on WHERE its host binding lives, not on whether it loads an adapter section alongside at the trigger." Second, append:

> A failure-triggered exception body whose trigger lives inside a host adapter's deferred reference realizes as a non-user-invocable skill too: the adapter keeps a resident trigger line that carries the first action and names the skill (`spacedock:fo-dispatch-recovery` — Degraded Mode, break-glass manual dispatch, budget-fail/dead-ensign handling), so the body loads only when its failure actually fires; the trigger line, not the boot core, is the namer.

### Test wiring

- **contractlint** (`internal/contractlint/boot_resident_closure_test.go`): add `skills/fo-dispatch-recovery/SKILL.md` to `deferredSkillCores` with anchors `## Degraded Mode`, `## Break-Glass Manual Dispatch`, `## Context Budget Failure and Dead Ensign Handling`; add `fo-dispatch-recovery` to `lazyLoadSkills`; add a walked deferred-module-bodies list containing `claude-fo-dispatch.md` so in-module `spacedock:` trigger tokens are stat-checked by the same extraction+os.Stat oracle (this retroactively closure-checks the existing `using-legacy-claude-team` token there, uncovered today). Extend `TestDeferredSkillProsePointersResolve`'s walk to that same deferred-module-bodies list — today it walks only `deferredSkillCores` files, so without the extension the trigger lines in `claude-fo-dispatch.md` would go unscanned. The extension is sound: trigger lines resolve because they carry the `Skill(skill="spacedock:fo-dispatch-recovery")` token on the same line, and the resident `## Degraded Mode (trigger)` heading is not an exact match for the watched `Degraded Mode` anchor, so it does not false-define the moved section.
- **ensigncycle greet guard** (`internal/ensigncycle/shallow_boot_measure_test.go:16`): `deferredFOSkillNames` += `"fo-dispatch-recovery"` so a pre-greet load reds; the existing RED-control machinery already proves the oracle can fail.
- **Live scenarios** (Claude-only — NOT `sharedRuntimeScenarios`, which fan out to the Codex runner; these flows are Claude-host-specific):
  - `degraded-bare` (transport decision: HEADLESS): the `/spacedock bare` trigger rides in the initial `-p` prompt — a mid-run captain turn is NOT injectable through the headless runner, so the mid-run variant would need the `ptyLiveDriver` (tmux send-keys nudge precedent, `pty_live_driver_test.go:426-482`) plus new detect-completion-then-send orchestration; that is recorded as the upgrade path, not shipped here. Headless observables: (i) `Skill` invocation with arg `spacedock:fo-dispatch-recovery` after the trigger, (ii) the verbatim captain report sentence ("Falling back to bare mode for the remainder of this session due to infrastructure failure. …"), (iii) every post-trigger `Agent()` lacks `name` and `run_in_background`. The cooperative-shutdown-sweep observable (a `shutdown_request` to a pre-trigger worker) is dropped with this transport — it needs the pty mid-run variant, since a boot-time trigger precedes any worker.
  - `break-glass-shim`: a PATH-shim `spacedock` wrapper that delegates to the real binary except `dispatch build` → stderr + exit 1 (mechanism precedent: the shallow-boot scenario's stub `gh` via `withStubPATH` — interface `claude_live_runner_test.go:82-88`, impl `:201-208`, use `:327`). Stream oracle: (i) a captain-facing helper-failure report before any `Agent()`, (ii) `Skill` arg `spacedock:fo-dispatch-recovery`, (iii) an `Agent()` call with `run_in_background=true`, a `{worker_key}-{slug}-{stage}` name, and a prompt carrying `Skill(skill="spacedock:ensign")` plus an inline `### Stage definition`.
  - **Equivalence protocol:** run both scenarios once against pre-split `main` — AFTER `jh` lands (see the merge-order dependency) — to capture baseline streams and prove the harness (shim + oracle) against the resident-body contract; check the streams in as fixtures under `internal/ensigncycle/testdata/`; offline unit oracles drive the fixtures with RED controls (a mutated fixture that keeps named/background dispatch after the bare trigger; one that hand-assembles a dispatch with no captain report). Post-split live runs must satisfy the same observables PLUS the skill-load assertion. This is the mechanism-validation-first ordering: the riskiest claim — a live FO at failure time actually follows a trigger-line skill load — is exercised end-to-end before/immediately after the contract edit, not asserted.

### Spike determination

No pre-design spike needed — every mechanism the design rests on is shipped and proven: (1) non-user-invocable `Skill()` loading — five shipped skills, live-validated in #457's greet oracle work (commit 59251f18); (2) PATH-stub binary injection into the live FO subprocess — the shallow-boot stub `gh` (`claude_live_runner_test.go:82-88, :201-208, :327`); (3) stream skill-argument extraction — `journeymetrics` `ClaudeTurn.SkillNames` (#457); (4) the contractlint closure/anchor walkers being extended. The one new harness piece (the failing-`dispatch build` shim) is validated first by the pre-split baseline run in the equivalence protocol above — the smallest end-to-end exercise, ordered before the contract edits.

## Acceptance criteria

- **AC-1 (value-measuring, the end the entity exists for).** `claude-fo-dispatch.md`'s resident size vs `origin/main` drops by ≥1,200 tokens measured (expected ~1,300), byte proxy ≤ −5,000 B (expected ~ −5,310: moved spans 7,820 B minus measured resident replacements 444 + ~715 + ~1,350 B — thin ~310 B headroom over the floor as measured today, but the baselines are taken after `jh` lands, which grows the moved span and improves the margin; see the merge-order dependency). The baseline is independent and can move the wrong way (the file can grow), cf. `trim-dispatch-adapter-prose` AC-1. **Tested:** `git show origin/main:skills/first-officer/references/claude-fo-dispatch.md | wc -c` vs `wc -c` on the merged file (reproducible on-disk evidence, commands + output in the stage report); token delta measured once via the Anthropic count-tokens API on the pre/post file contents and recorded alongside. The honest plugin-wide sum (resident file deltas + the new SKILL.md size + the registry row) is reported with it, so the resident-vs-deferred trade is on the record.
- **AC-2 (behavioral equivalence, degraded mode).** Post-split, the `/spacedock bare` trigger produces the same observables the pre-split contract mandates — the verbatim captain report sentence and bare-shape `Agent()` calls (no `name`, no `run_in_background`) for every subsequent dispatch — plus the `spacedock:fo-dispatch-recovery` skill load. (The cooperative-shutdown-sweep observable is out of this AC's oracle: the shipped headless transport triggers at boot, before any worker exists; the pty mid-run variant that could exercise it is the recorded upgrade path.) **Tested:** the `degraded-bare` live scenario (claude-live lane; the merge gate, since this touches `skills/**`) and its checked-in fixture oracle with RED control; the pre-split baseline capture is the equivalence reference. No prose-grep: the oracle keys on stream tool-call structure and the emitted report text, not on contract wording.
- **AC-3 (behavioral equivalence, break-glass).** Post-split, a non-zero `spacedock dispatch build` produces: captain-facing helper-failure report first, the skill load, then a break-glass-shaped `Agent()` (named background worker, `Skill(skill="spacedock:ensign")` first action, inline stage definition). **Tested:** the `break-glass-shim` live scenario + fixture oracle with RED control, baseline-captured pre-split as above.
- **AC-4 (structural closure).** Every deferred load point named by the boot-resident bodies AND by `claude-fo-dispatch.md` resolves on disk (including `spacedock:fo-dispatch-recovery`); the skill carries its three section anchors; no surviving contract file points at a moved section in its old home (the `using-legacy-claude-team` re-points land); a pre-greet load of the new skill reds the greet guard. **Tested:** `go test ./internal/contractlint ./internal/ensigncycle` offline — `TestBootResidentDeferredLoadPointsResolve` (+ the new deferred-module-bodies walk), `TestDeferredSkillCoresResolveAndCarryCeremony`, `TestDeferredSkillProsePointersResolve`, the fold test, and the extended `deferredFOSkillNames` guard, each with its existing/added RED control.
- **AC-5 (mechanism, serves AC-1).** The `docs/runtime-support.md` discriminator covers the failure-triggered case per the doc diff above, and the `## Deferred load points` registry row ships. **Tested:** the applied diffs are on-disk state AC-4's closure and fold tests bind to (the registry row is what makes the sole-namer test pass); counts only paired with AC-1.

## Test plan

1. **Offline (no model spend, runs per-PR):** the contractlint additions and greet-guard extension in AC-4, plus fixture-driven unit oracles for the two scenarios (RED controls included). Cost: small Go test additions to existing walkers; `go build ./... && go test ./...` green.
2. **Live (claude-live lane, maintainer-approved environment; gates the merge since `skills/**` changes):** the two Claude-only scenarios in AC-2/AC-3. Cost: two multi-minute live runs per matrix model; the pre-split baseline capture is a one-time cost during implementation.
3. **Measurement (AC-1):** one-off git/wc + count-tokens comparison recorded in the implementation stage report.
4. **Order:** build shim + scenarios and capture pre-split baselines FIRST (mechanism validation + equivalence reference), then the contract split, then contractlint wiring, then post-split live runs and the AC-1 measurement.

## Stage Report: ideation

- DONE: The split honors the named reliability risk: each resident trigger line carries the first action (not a bare pointer), pre-invocation guards stay resident (context-budget probe invocation, Awaiting Completion), and the design states per exception body why its load point is safe at failure time.
  Diffs 1-3 each carry the first action on the resident line; `## Awaiting Completion` and the probe/fail-safe/model-enum lines stay resident (reuse-conditions 0 and 4 run on the happy path); per-body safety rationale in "Why each load point is safe at failure time".
- DONE: ACs include a measured token delta for claude-fo-dispatch.md's resident size and a behavioral-equivalence check for degraded-mode and break-glass flows post-split, each paired with how it is tested (contractlint closure, live or fixture evidence — no prose-grep as proof).
  AC-1 (≥1,200 tok / ≤ −5,000 B vs origin/main, git+count-tokens evidence); AC-2/AC-3 (pre-split-baselined live scenarios + checked-in fixture oracles with RED controls, keyed on tool-call structure, not wording); AC-4 (contractlint closure + greet guard).
- DONE: The #457 non-user-invocable-skill pattern is followed exactly (name the shipped precedent files) and the deferred-load-points registry in the shared core is updated consistently, with concrete before/after prose diffs in the task body.
  Precedents named: skills/fo-status-viewer/SKILL.md, skills/fo-write-core/SKILL.md (59251f18/#457), plus present-gate, feedback-rejection-flow, using-legacy-claude-team; registry gains an additive fifth row (before/after quoted); Diffs 1-3 plus four re-point diffs and the runtime-support.md doc diff are in the body verbatim.

### Summary

Ideation designs the split as ONE non-user-invocable skill (`fo-dispatch-recovery`, three verbatim section anchors) rather than three, because every SKILL.md description is resident in every session's skill list and the three failure paths share one trigger family. Measured the spans (7,820 B, ~1.9k tok; expected resident drop ~5.6 KB / ~1.5k tok), resolved the model-enum constraint (fo-dispatch-core.md:93 pins the enum resident for happy-path reuse-condition-4), and recorded a "no spike needed" determination on four shipped mechanisms, with the one new harness piece (the failing-`dispatch build` shim) validated first via pre-split baseline captures that double as the equivalence reference.

## Stage Report: ideation (cycle 2)

Staff review returned six material findings; all applied as body reconciliations (no architecture change).

- DONE: M1 — AC-1 arithmetic corrected to measured replacements (444 / ~715 / ~1,350 B), expected net ~ −5,310 B / ~1,300 tok, floors unchanged with thin headroom noted.
  AC-1 text and Diffs 1-3 byte figures updated; margin note points at the merge-order dependency.
- DONE: M2 — `fo-dispatch-core.md:133` added to the re-point list as a deliberate non-edit with the transitive-resolution rationale (host-neutral core must not name a Claude-only skill; the adapter's trigger line is the entry).
- DONE: M3 — the "do not route `SendMessage` to any pre-trip worker name" clause added to the Diff 2 trigger line (recorded ~55 B spend against AC-1), and the safety rationale now cites fo-dispatch-core.md's host-neutral degraded effects as the resident backstop.
- DONE: M4 — degraded-bare transport decided: HEADLESS (trigger in the initial `-p` prompt), sweep observable dropped from AC-2's oracle, pty mid-run variant (`pty_live_driver_test.go:426-482` send-keys precedent) recorded as the upgrade path; `withStubPATH` citations fixed to `:82-88, :201-208, :327`.
- DONE: M5 — merge-order dependency section added (wc → jh → this entity): Diff 1 re-derived against jh's landed text, enum values un-quoted everywhere (wc amends the enum), baselines and AC-1 measurement after jh lands.
- DONE: M6 — the vacuous scanner sentence replaced with an explicit extension of `TestDeferredSkillProsePointersResolve`'s walk to the deferred-module-bodies list, with the soundness argument (same-line Skill token; `## Degraded Mode (trigger)` does not false-define).
- DONE: Polish — using-legacy-claude-team `:50`/`:52` recorded as checked-and-valid, the runtime-support.md lead clause diff now accommodates the third shape, and the clean-blast-radius record is in the body.

### Summary

All six staff findings applied in the task body: corrected AC-1 measurement arithmetic, the fo-dispatch-core transitive-pointer rationale, the pre-trip-names trigger clause with its byte cost on the record, a headless transport decision for degraded-bare with the pty upgrade path named, the wc → jh → 41 merge-order dependency, and a non-vacuous prose-pointer-walk extension. Scope held to task-body edits only.
