---
id: p48fjxe9e2rz23afaqmzj1pg
title: Encode the live-verification requirement — gate runtime-observable ACs on a live run + promote 8y's shared scenarios to a first-class live-scenario primitive
status: implementation
source: captain (2026-06-03, session 11) — after yy's fix passed offline + 3 detached-audit cycles but FAILED its live run, the captain noted the live-verification rule is prose-only (yy slipped) AND that agents don't intuitively reach for "descriptive runbook → run a real agent → grade durable outcomes" live tests
score: "0.40"
worktree: .worktrees/spacedock-ensign-live-verification-gate
started: 2026-06-03T15:28:23Z
completed:
verdict:
issue:
---

The requirement that a **runtime-observable AC** (one whose truth can only be decided by running the real producer) be confirmed by a **live run** is encoded only as advisory prose — the README validation stage and the merge-on-proof bar in handoffs. Nothing refuses to mark such an AC PASSED until a live run is on the record. `yy` (sonnet-live-ci-flake, #282) was marked `pending-live-run` and merged with its live AC unconfirmed; the post-merge live retrigger then revealed the fix doesn't actually close the flake. This violates the workflow's own Working Principle ("prefer a code gate over a prose-only rule; a prose-only rule has a ceiling of wording-present, which is not behavior").

## Problem

Two coupled gaps:

1. **The live-verification rule is prose-only and unenforced.** A runtime-observable AC can reach terminal/PASSED on offline proof alone (unit tests, a recorded-stream fixture, a contract-text check). But a recording proves the *consumer* (the streamwatcher localizes a hang in a frozen stream), never the *producer* (does the real model, reading the contract, actually exit?) — the fix produces a *different, future* stream you can only get by running it live. yy passed every offline check + 3 adversarial audit cycles and was still wrong.

2. **The live-scenario pattern is not reachable/intuitive for agents.** Agents default to deterministic offline proxies because "write a runbook → run a real agent → grade the durable outcomes" is **eval design**, not software testing — it requires defining observable success criteria for non-deterministic behaviour, a skill the contract/templates neither surface nor make cheap. `8y` built the substrate (host-neutral `sharedRuntimeScenarios()` + per-host runner adapters + durable state/output assertions) but it is CI-internal, not a primitive the FO/ensign reaches for. So agents pick the proxy and don't think to run the thing.

## Proposed approach (three layers, code-first)

1. **Runtime-observable classification + a citation gate, built ON 2a's extractor.** The in-flight `2a` (require-external-proof-guard, status: validation, PR #277) already SHIPPED an AC extractor + classifier: `ClassifyEntityACs` (`internal/status/external_proof.go:107`) walks `**AC-N` blocks, isolates the proof clause from a marker (`Verified by`/`Oracle:`/`Proof:`/`End state`), and flags self-referential proofs. It is wired into BOTH `runSet`'s terminal `--set` guard (`handlers.go:186`, behind `isTerminalUpdate()`) and `validateWorkflow`'s sub-check, gated by the README opt-in `require-external-proof: true` (`resolveExternalProofPolicy`, `external_proof.go:210`, default OFF). **p4 does NOT invent the extractor — it EXTENDS 2a's**: it adds a *runtime-observable* classification beside 2a's self-referential one, reusing the same block-walk / `isolateProofClause` / `stripFrontmatter` machinery. **p4 sequences AFTER 2a lands** (it depends on `ClassifyEntityACs` and the `external_proof.go` surface being present).

   - **Classification = explicit author convention, not keyword inference.** A runtime-observable AC declares itself by leading its proof clause with `live` — `Verified by: live <artifact-ref>`. A property of the declared text at the claim's own level, not a prose-keyword scan (which over/under-gates — see spike).
   - **Gate point = the same terminal `--set` guard 2a uses, NOT `runArchive`.** An entity reaches terminal status via `--set` (verdict/completed) BEFORE `runArchive` ever runs; 2a's guard already fires there under `isTerminalUpdate()`. p4's live-run check sits in the SAME guard block in `runSet`, right beside 2a's self-referential check — refuse the terminal `--set` when a runtime-observable AC cites no *resolvable* live-run artifact. Putting it at `runArchive` instead would be a SECOND AC-classifying guard at a LATER point with its own opt-in — two guards, two points, two policies that collide. One extractor, one gate point.
   - **Opt-in = ride 2a's `require-external-proof` policy (no new key, no non-dev over-gate).** The live-run check is INERT unless the workflow declares `require-external-proof: true` — exactly like 2a. A non-dev workflow whose proof is legitimately a published artifact or human review (the `development.md:113` caveat) never declared the key, so its `Verified by: live …` AC (if any) is untouched. This is the answer to the over-gate risk: the convention only bites in a workflow that opted into external-proof enforcement.
   - **Citation = resolvable artifact ref.** `ci-run:<id>` (the GitHub Actions run resolves — see spike) or `session:<path>` (JSONL present on disk). A placeholder like `<pending>` fails. This turns `pending-live-run` from a hopeful label into an enforced state.
   - **`--force` honored, with a louder warning than 2a's.** 2a lets `--force` bypass its guard with a `Warning: --force overriding require-external-proof`. p4's live-run bypass mirrors the mechanism but the warning names the specific risk being overridden — `Warning: --force overriding live-run requirement on entity X — a runtime-observable AC is being terminalized without a cited live run (this is exactly the yy slip)`. The point is not to make it impossible (an operator override must exist) but to make the escape LOUD and self-incriminating, so a default merge cannot silently skip it.
2. **Promote 8y's scenarios to a first-class live-scenario primitive.** 8y's substrate lives in `internal/ensigncycle/`. The **host-neutral triple already runs offline** — `sharedRuntimeScenario` / `sharedRuntimeScenarios()` (`shared_scenarios_test.go`) and the per-scenario fixture-writer + prompt + assertion carry NO `//go:build live` tag; only the two RUNNER ADAPTERS (`claude_live_runner_test.go`, `codex_live_runner_test.go`) are `//go:build live`-gated, because only the real-agent launch needs a credential. This STRENGTHENS the reachability argument: the authorable surface (define a scenario, assert durable outcomes) is already build-tag-free; what's buried is only that the whole thing sits in `_test.go` inside one internal package. Its shape already IS the target triple: the per-scenario *prompt* (`gatePrompt()`) is the runbook, the *fixture writer* (`writeGateWorkflow()`) is the setup, and the *assertion over (before, after, observedFinalMessage)* (`assertGateHeld()`) is the durable-outcome grade. Promotion = lift this triple into a primitive an FO/ensign authors **outside** the ensigncycle test package: a scenario = {descriptive runbook/prompt, setup/fixture, expected durable outcomes (entity state before→after + observed output, NOT transcript phrasing)}, runnable by a real Claude/Codex agent via the existing runner adapters and graded on those outcomes. The runner adapter (launch + observed-extract) and the host-neutral table stay; what changes is reachability — a scenario becomes authorable as easily as a Go test instead of buried behind the live build tag in one internal package.
3. **Surface the pattern in the dev template.** Add to "choose proof at the claim's altitude": *if the claim is what an agent/model DOES at runtime, the proof is a scripted live scenario graded on durable state/output — a recording proves the watcher not the producer; a contract-text check proves the words not the behaviour.*

## Spike (riskiest unknown — RUN at ideation)

**The genuinely-novel unknown is `ci-run:<id>` RESOLUTION**, not classification. Classification (lead-with-`live`) is a trivial string check, and 2a's `ClassifyEntityACs` already proves the block-walk + proof-clause isolation works. The unexercised mechanism this design rests on is: how does `ci-run:<id>` resolve to "this run exists", and what happens offline / unauthenticated?

**Spike A — `ci-run:<id>` resolution (the genuinely-novel bit).** Exercised `gh` (already on PATH; the binary already shells to it in `internal/status/*.go`) against the real repo:

| Input | Command | Exit | Signal |
|---|---|---|---|
| real run id `26891718562` | `gh api repos/.../actions/runs/26891718562` | 0 | resolves → CITED-AND-REAL |
| nonexistent id `1` | `gh api repos/.../actions/runs/1` | 1 | stderr `Not Found (HTTP 404)` → DEFINITIVELY ABSENT |
| unreachable host (offline/auth-broken) | same, `GH_HOST=invalid` | 1 | stderr `failed to determine base repo…` → INDETERMINATE |

**Key finding that reshapes the design:** a naive "exit ≠ 0 ⇒ invalid citation" is WRONG — it conflates "run genuinely absent" (gate SHOULD refuse) with "I can't reach GitHub right now" (a network blip must NOT refuse a valid entity). The resolver MUST distinguish the **404 → refuse** from a **connectivity/auth error → indeterminate**: on indeterminate, the gate does NOT mark the citation invalid; it surfaces a tooling error (`live-run citation could not be checked: <gh error>`) distinct from a gate refusal, so a transient never masquerades as a missing live run. The `gh api …/actions/runs/<id>` form gives the cleanest signal (definitive `HTTP 404` vs other error text).

**Spike B — classification + session citation (confirmatory; classification itself is 2a-proven).** A throwaway Go prototype over discriminating fixture AC bodies, adversarially broken to confirm the tests aren't tautological:

| Fixture AC body | Expected | Gate result |
|---|---|---|
| offline: `Verified by: a Go unit test TestRoundTrip…` | not refused | not refused (PASS) |
| runtime, no citation: `Verified by: live <no artifact yet>` | REFUSE | refused (PASS) |
| runtime, cited: `Verified by: live ci-run:12345` | not refused | not refused (PASS) |
| **yy-shape**: `Verified by: live <none — passed offline + 3 audit cycles, merged pending-live-run>` | REFUSE | refused (PASS) |
| `session:<present .jsonl>` vs `session:<absent path>` | pass / REFUSE | correct both ways |

Adversarial edits caught both failure directions: making the classifier never mark runtime-observable (the status-quo under-gate) **reds the yy-shape test**; making the gate refuse offline ACs too (over-gate) **reds the offline test**.

**Settled design (seeds the implementation's first test):**
- **Classification = explicit convention, not inference.** Runtime-observable iff the proof clause leads with `live`. Reuses 2a's block-walk; no prose-keyword scanning.
- **Citation = resolvable artifact ref.** `ci-run:<id>` via `gh api …/actions/runs/<id>` — **404 ⇒ refuse, connectivity/auth error ⇒ indeterminate tooling error (not a refusal)**; `session:<path>` ⇒ JSONL present on disk. Placeholder ⇒ refuse.
- **Gate point = 2a's terminal `--set` guard in `runSet`** (`handlers.go:186`), beside 2a's self-referential check — NOT `runArchive`. Inert unless `require-external-proof: true`. `--force` bypasses with a loud, risk-naming warning.
- **No further spike for the other layers:** layer 2 composes 8y's already-proven (build-tag-free) triple + runner adapters (reachability change, not a new mechanism); layer 3 is a text-presence change to the dev template (proof at the text's own level).
- **Two indeterminate cases the implementer must handle, named here:** an entity with **no ACs** ⇒ no flags ⇒ no-op (matches `classifyEntityFile` returning nil); an **already-archived** entity ⇒ NOT retro-checked (the gate fires only on the live terminal `--set`, never sweeps the `_archive`).

## Out of scope

- Re-running yy's fix — that's `at` (non-interactive-teardown-exit); this entity is the *meta*-fix that makes such slips gate-caught.
- A general eval framework beyond what spacedock workflows need.

## Acceptance criteria

**AC-1 — Under `require-external-proof: true`, a runtime-observable AC (declared by `Verified by: live …`) cannot reach terminal `--set` without a resolvable live-run citation.**
Verified by: a Go test driving the terminal `--set` guard in `runSet` (the SAME guard 2a extends) over a fixture entity — RED (exit 1, frontmatter untouched, diagnostic) when the runtime-observable AC's `live` clause cites no resolvable artifact, GREEN when it cites a real `ci-run:`/`session:` ref; an offline-checkable AC (any non-`live` proof clause) is never refused; and with `require-external-proof` absent/false the guard is inert (no over-gate). The classifier is the explicit `live`-lead convention reusing 2a's `ClassifyEntityACs` block-walk, not prose keyword inference. (Depends on 2a's `external_proof.go` surface — p4 sequences after 2a lands.)

**AC-2 — The live-scenario primitive is authorable + runnable independent of 8y's CI-internal plumbing.**
Verified by: a scenario authored as {runbook, setup, durable-outcome assertions}, run by a real agent and graded on outcomes, with at least one negative case (a broken durable outcome reds the grade). (Note the recursion: AC-2 is ITSELF a runtime-observable `live` AC — only a real agent run decides it — so AC-2's own proof is subject to AC-1's gate. Its `Verified by` must carry a `ci-run:`/`session:` citation at terminal, making p4 the first entity to eat its own gate.)

**AC-3 — The dev template carries live-scenario-for-runtime-claims guidance as a third opt-in recommended practice.**
Verified by: a Go presence check over `skills/commission/references/templates/development.md` (the same kind that guards the existing "External-proof acceptance criteria" / "Detached adversarial audit" recommended-practice blocks) confirming a new block names the live-scenario pattern and the recording-proves-the-watcher / text-check-proves-the-words distinction. (Presence check is legitimate here: the claim is about the text itself — proof at the claim's own level.)

**AC-4 — The yy slip would now be gate-caught.**
Verified by: a fixture reproducing yy's shape (a runtime-observable AC declared `Verified by: live …` but cited only by offline proof / a placeholder, at terminal `--set` under `require-external-proof: true`) is REFUSED by the AC-1 guard. Plus a `ci-run` resolution unit test asserting the spike's three-way split: real id → cited-and-real, 404 → refuse, connectivity/auth error → indeterminate tooling error (NOT a refusal), so a network blip cannot masquerade as a missing live run.

## Test plan

- **Spike DONE at ideation** (see Spike section): the `ci-run:<id>` resolution mechanism (the genuinely-novel bit) is exercised — `gh api …/actions/runs/<id>`, 404⇒refuse / connectivity-error⇒indeterminate. Classification + session citation confirmed. Implementation's first test seeds from the fixture rows + the three-way resolution split.
- **AC-1/AC-4 (Go, table-driven over fixture entity files):** drive the terminal `--set` guard in `runSet`; assert exit-1 + untouched frontmatter on an uncited runtime-observable AC, pass on a cited one and on offline ACs, and inert under `require-external-proof` absent/false. AC-4 adds the yy-shape fixture row + the `ci-run` three-way resolution unit test. Cost: low — fixture entities + 2a's existing `runSet`/`classifyEntityFile` harness; the resolution test stubs `gh` output (404 vs connectivity) so no network in unit tests. Adversarial guard: re-run the spike's two break-edits (never-classify-runtime; refuse-offline-too) against the real gate to confirm the suite reds, not just greens.
- **AC-3 (Go presence check):** assert the dev template's new recommended-practice block. Cost: trivial.
- **AC-2 (live-scenario exercise):** author one scenario via the promoted primitive {runbook, setup, durable-outcome assertions}, run it against a real agent through the runner adapter, grade on entity before→after + observed output; add a negative case where a deliberately broken durable outcome reds the grade. Cost: medium — needs a live credential + minutes of agent wallclock (live-gated, like 8y's scenarios). This is `at`'s live AC-1 in primitive form.
- **High-stakes → detached adversarial audit BEFORE merge.** Trigger: the gate machinery (a status/guard mutation path) **and** the shipped contract/template change — both are the high-stakes surfaces the dev template names. The audit runs read-only on a detached checkout of the merge result and tries to REFUTE: construct an entity edit that should be gate-refused and confirm it is, and confirm an offline AC survives an edit that would over-gate.

## Notes

- Motivated by yy's live-run failure (session 11): offline + 3 audit cycles passed, the live run failed — the strongest possible case for the gate. `8y` built the substrate; `at` will be the first entity that *needs* this primitive (its live AC is a scripted scenario).
- This **extends 2a's `ClassifyEntityACs`** (`external_proof.go`, shipped in PR #277, status validation) with a runtime-observable classification + live-run citation check, at 2a's terminal `--set` guard, under 2a's `require-external-proof` opt-in. p4 sequences AFTER 2a lands. It is NOT a new extractor and NOT a new gate point — one classifier surface, one gate, two AC families (self-referential + uncited-runtime).

## Stage Report: ideation

- DONE: SPIKE FIRST — the gate's classification + citation mechanism
  Ran a throwaway Go classifier+citation prototype over 4 discriminating fixture AC bodies (incl. yy-shape) + 2 adversarial break-edits; all discriminated correctly and both break-edits red the suite. Settled: explicit `Verified by: live <ref>` convention (not keyword inference) + resolvable `ci-run:`/`session:` citation, gated at archival. Recorded in the "Spike" section with a result table.
- DONE: Land the three-layer design with each AC backed by a CODE or TEST gate
  (a) AC-1 refusal gate anchored to `runArchive`'s existing merge-hook guard + the `--validate`/`validateWorkflow` AC-extractor (which does NOT exist yet — corrected the body's prior "extend the existing cross-check" claim); (b) live-scenario primitive grounded in 8y's actual `_test.go`-only `sharedRuntimeScenario` triple (prompt/fixture-writer/assertion) — promotion = reachability outside the ensigncycle test package; (c) AC-3 lands as a third opt-in recommended practice beside External-proof / Detached-audit.
- DONE: Test plan (AC-1/AC-4 gate refusal incl. yy fixture, AC-3 presence, AC-2 live-scenario + negative case; detached-audit trigger named)
  Test plan rewritten: AC-1/AC-4 table-driven Go over fixture entities with the spike's break-edits as an adversarial guard, AC-3 Go presence check, AC-2 live-scenario run + broken-outcome negative; detached-audit trigger named as gate-machinery + shipped contract/template (both high-stakes per the dev template).

### Summary

Fleshed out the ideation with the riskiest unknown exercised first. The spike settled the classifier as an explicit `Verified by: live <ref>` convention (rejecting prose keyword inference, which over/under-gates) with a resolvable `ci-run:`/`session:` citation, gated at the archival/terminal refusal modeled on the existing merge-hook guard in `runArchive`. Corrected two factual errors in the prior body: there is no pre-existing AC-evidence cross-check to extend (this builds the first AC-extractor in `validateWorkflow`), and 8y's scenario substrate is `_test.go`-only/live-gated so promotion is a reachability change, not a new mechanism. High-stakes (gate machinery + shipped contract/template) → detached adversarial audit before merge.

## Stage Report: ideation (cycle 2)

Reworked after staff REWORK (5 material findings, all verified against real code). Three sound pieces kept: explicit `Verified by: live <ref>` convention, the merge-hook-guard family resemblance, AC-1/AC-3/AC-4.

- DONE: Finding 1 — re-anchor on 2a's shipped extractor
  Read `external_proof.go` in 2a's worktree: `ClassifyEntityACs` (line 107) is shipped, wired at `runSet` terminal guard (handlers.go:186) + validateWorkflow. Body, ACs, Notes now say p4 EXTENDS 2a (not "first extractor"); recanted "opt-in prose, not binary code". p4 sequences after 2a lands.
- DONE: Finding 2 — gate point reconciled to 2a's terminal `--set`, not `runArchive`
  Verified an entity terminalizes via `--set` (verdict/completed) BEFORE `runArchive`; 2a fires under `isTerminalUpdate()` there. Moved p4's check into the same `runSet` guard block; explained why two points/two policies collide.
- DONE: Finding 3 — non-dev over-gate closed
  Live-run check rides 2a's `require-external-proof` opt-in (default OFF, `resolveExternalProofPolicy`); inert unless declared, so a non-dev workflow's `live` AC is untouched.
- DONE: Finding 4 — spiked `ci-run:<id>` resolution (the genuinely-novel bit)
  Ran `gh api …/actions/runs/<id>`: real id→exit 0; fake→exit 1 `HTTP 404`; offline→exit 1 different error. KEY: 404⇒refuse, connectivity/auth⇒indeterminate (NOT refusal) so a blip can't masquerade as a missing live run. Recorded as Spike A + AC-4 unit test.
- DONE: Finding 5 — `--force` specified
  Guard honors `--force` (an operator override must exist) but with a louder, risk-naming warning than 2a's, so a default merge can't silently escape.
- DONE: Polish — 8y triple is build-tag-free (only runner adapters are `//go:build live`), tightened; AC-2 recursion noted (it's a `live` AC its own gate holds); no-AC entity (no-op) + already-archived (not retro-checked) named in Spike.

### Summary (cycle 2)

p4 is now a runtime-observable EXTENSION of 2a's `ClassifyEntityACs` at the SAME terminal `--set` guard under the SAME `require-external-proof` opt-in — not a parallel gate at `runArchive`, and not a new extractor. The riskiest unexercised mechanism (`ci-run:<id>` resolution) is spiked: a 404 refuses, a connectivity/auth error is indeterminate (never a false refusal). `--force` stays as a loud, self-incriminating escape. p4 depends on 2a landing and sequences after it.

## Stage Report: implementation

- DONE: AC-1 + AC-4 (the gate) — terminal `--set` guard refuses an uncited runtime-observable `live` AC under `require-external-proof: true`
  `internal/status/live_proof.go` (`classifyLiveACs` reuses ep's block-walk + `isolateProofClause`; `resolveLiveCitation` three-way: ci-run:/session: resolve→pass, placeholder/404/absent→refuse, connectivity/auth→indeterminate tooling error not a refusal) + guard block in `handlers.go` beside ep's self-ref check (same opt-in, same `isTerminalUpdate`, inert when off). `--force` bypasses with a loud yy-naming warning. Commit b8f03087.
- DONE: AC-4 yy-shape fixture refused + ci-run 3-way resolution unit test
  `TestLiveRunGuardRefusesUncitedRuntimeAC` (yy-shape placeholder → exit 1, frontmatter untouched) + `TestResolveLiveCitationThreeWay` (real→citedAndReal, 404→definitivelyAbsent, connectivity→indeterminate via injected stub, offline). TDD: guard test written first, watched RED (exit 0, no guard), then GREEN.
- DONE: Adversarial guard — both spike break-edits red the suite
  never-classify-runtime (return nil) reds `TestLiveRunGuardRefusesUncitedRuntimeAC` + `TestClassifyLiveACs/yy-shape`; refuse-offline-too (over-gate) reds `TestLiveRunGuardNeverRefusesOfflineAC`. Restored, suite green.
- DONE: AC-2 (the primitive) — live-scenario authorable outside ensigncycle, graded on durable outcomes, with a negative case
  `internal/livescenario/scenario.go` (`Scenario{Runbook,Setup,Assert}` + `Run` + `Runner` seam). `TestScenarioGradesBrokenOutcomeNegative` proves a broken durable outcome (gate advanced to done, or no gate-review observed) reds the grade. `//go:build live` adapter `internal/ensigncycle/livescenario_adapter_live_test.go` drives the primitive through 8y's real Claude launch adapter — AC-2's runnable-by-a-real-agent half (vet-clean under `-tags live`).
- DONE: AC-3 — dev template carries live-scenario-for-runtime-claims as a third opt-in recommended practice
  New `### Live scenario for runtime claims` block in `skills/commission/references/templates/development.md` beside External-proof / Detached-audit (recording-proves-the-watcher / text-check-proves-the-words distinction). Locked by `internal/hostneutrality/live_scenario_practice_test.go` (presence check scoped to the Recommended-practices section).
- DONE: Full offline `go test ./...` green; `go vet` clean
  1020 passed in 13 packages; `go vet ./...` and `go vet -tags live ./internal/ensigncycle/` clean; `gofmt -l` clean.

### Summary

Re-anchored on ep's (#290) injected-vocabulary `ClassifyEntityACs` and layered a NEW runtime-observable dimension (`Verified by: live <ref>` lead) at the SAME terminal `--set` guard under the SAME `require-external-proof` opt-in. The live-run citation resolves three-way so a network blip is an indeterminate tooling error, never a false refusal. 8y's {runbook, setup, durable-outcome} triple is promoted into `internal/livescenario` — importable anywhere, graded on durable state/output with a real negative case, and runnable against the existing live launch adapter. AC-2's own live run (which produces its citation under AC-1's gate) is the validation-stage live exercise; the offline gradable core + the live adapter are both in place. Two notes for the FO: (1) AC-1/AC-4 guard machinery is the high-stakes status-mutation surface the detached audit should target; (2) AC-2's live half (`go test -tags live -run TestLivePrimitiveRunsAgainstClaudeAdapter`) needs a credential to produce the session/ci-run artifact that satisfies its own `Verified by: live …` citation.
