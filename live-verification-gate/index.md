---
id: p48fjxe9e2rz23afaqmzj1pg
title: Encode the live-verification requirement — gate runtime-observable ACs on a live run + promote 8y's shared scenarios to a first-class live-scenario primitive
status: ideation
source: captain (2026-06-03, session 11) — after yy's fix passed offline + 3 detached-audit cycles but FAILED its live run, the captain noted the live-verification rule is prose-only (yy slipped) AND that agents don't intuitively reach for "descriptive runbook → run a real agent → grade durable outcomes" live tests
score: "0.40"
worktree:
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

1. **AC classification + a real gate.** At ideation, classify each AC as *offline-checkable* (a unit test / fixture / text check decides it) or *runtime-observable* (only a live run decides). Classification is an **explicit author convention, not keyword inference**: a runtime-observable AC declares itself by leading its Verify clause with `live` — `Verified by: live <artifact-ref>`. This keeps the classifier a property of the declared text at the claim's own level, not a substring search guessing intent from prose (which over- and under-gates — see the spike). The terminal gate then **refuses** archival of an entity carrying a runtime-observable AC unless that AC cites a *resolvable* live-run artifact (`ci-run:<id>` whose run exists, or `session:<path>` whose JSONL is present); a placeholder like `<pending>` fails the check. This turns `pending-live-run` from a hopeful label into an enforced state.

   The gate's home is the **terminal/archival refusal in `runArchive`** (`internal/status/mutate.go`), modeled on the merge-hook guard already there (archival is refused unless `pr` is set / `mod-block` in flight / `--force` / `merge: local`). The live-run guard is its sibling: refuse archival of an entity with an uncited runtime-observable AC. The AC extraction + classification itself rides the `--validate` flow (`validateWorkflow` in `internal/status/validate.go`), which today is purely structural (id rules, flat/folder conflicts, stage-name regex) and carries **no AC-extractor yet** — this builds the first one. There is no pre-existing "AC-evidence cross-check" to extend; the external-proof guard the dev template mentions is opt-in prose, not binary code.
2. **Promote 8y's scenarios to a first-class live-scenario primitive.** 8y's substrate (`internal/ensigncycle/shared_scenarios_test.go` + the per-host `claude_live_runner_test.go` / `codex_live_runner_test.go`) is today **`_test.go`-only and `//go:build live`-gated** — `sharedRuntimeScenario` and `sharedRuntimeScenarios()` appear in no non-test source. Its shape already IS the target triple: the per-scenario *prompt* (`gatePrompt()`) is the runbook, the *fixture writer* (`writeGateWorkflow()`) is the setup, and the *assertion over (before, after, observedFinalMessage)* (`assertGateHeld()`) is the durable-outcome grade. Promotion = lift this triple into a primitive an FO/ensign authors **outside** the ensigncycle test package: a scenario = {descriptive runbook/prompt, setup/fixture, expected durable outcomes (entity state before→after + observed output, NOT transcript phrasing)}, runnable by a real Claude/Codex agent via the existing runner adapters and graded on those outcomes. The runner adapter (launch + observed-extract) and the host-neutral table stay; what changes is reachability — a scenario becomes authorable as easily as a Go test instead of buried behind the live build tag in one internal package.
3. **Surface the pattern in the dev template.** Add to "choose proof at the claim's altitude": *if the claim is what an agent/model DOES at runtime, the proof is a scripted live scenario graded on durable state/output — a recording proves the watcher not the producer; a contract-text check proves the words not the behaviour.*

## Spike (riskiest unknown — RUN at ideation)

**Question:** the gate's classification + citation mechanism. How does the binary KNOW an AC is runtime-observable, and how does it verify a cited live-run artifact is real, without over-gating (blocking offline-checkable ACs) or under-gating (letting runtime ACs through — the status quo that let yy slip)?

**Exercise:** a throwaway Go test prototyped the classifier + citation check over four discriminating fixture AC bodies and proved discrimination, then was adversarially broken to confirm the tests aren't tautological. Results:

| Fixture AC body | Expected | Gate result |
|---|---|---|
| offline: `Verified by: a Go unit test TestRoundTrip…` | not refused | not refused (PASS) |
| runtime, no citation: `Verified by: live <no artifact yet>` | REFUSE | refused (PASS) |
| runtime, cited: `Verified by: live ci-run:12345` | not refused | not refused (PASS) |
| **yy-shape**: `Verified by: live <none — passed offline + 3 audit cycles, merged pending-live-run>` | REFUSE | refused (PASS) |
| `session:<present .jsonl>` vs `session:<absent path>` | pass / REFUSE | correct both ways |

Adversarial edits confirmed the tests catch both failure directions: making the classifier never mark runtime-observable (the status-quo under-gate) **reds the yy-shape test**; making the gate refuse offline ACs too (over-gate) **reds the offline test**.

**Settled design (seeds the implementation's first test):**
- **Classification = explicit convention, not inference.** An AC is runtime-observable iff its Verify clause leads with `live`: `Verified by: live <artifact-ref>`. No prose-keyword scanning (that is a substring search standing in for a behavioral claim — banned by the proof-altitude rule and the source of over/under-gating).
- **Citation = resolvable artifact ref.** `ci-run:<id>` (run id exists) or `session:<path>` (JSONL present on disk). A placeholder fails.
- **Gate point = the terminal/archival refusal**, sibling to the existing merge-hook guard in `runArchive`.
- **No spike needed beyond this** for the other layers: layer 2 composes 8y's already-proven runner/assertion substrate (reachability change, not a new mechanism); layer 3 is a text-presence change to the dev template (proof at the text's own level).

## Out of scope

- Re-running yy's fix — that's `at` (non-interactive-teardown-exit); this entity is the *meta*-fix that makes such slips gate-caught.
- A general eval framework beyond what spacedock workflows need.

## Acceptance criteria

**AC-1 — A runtime-observable AC (declared by `Verified by: live …`) cannot reach terminal/archival without a resolvable live-run citation.**
Verified by: a Go test driving the archival/terminal gate over a fixture entity — RED (archival refused, non-zero exit + diagnostic) when the runtime-observable AC's `live` clause cites no resolvable artifact, GREEN when it cites a real `ci-run:`/`session:` ref; an offline-checkable AC (any non-`live` Verify clause) is never refused. The classifier is the explicit `live`-lead convention, not prose keyword inference.

**AC-2 — The live-scenario primitive is authorable + runnable independent of 8y's CI-internal plumbing.**
Verified by: a scenario authored as {runbook, setup, durable-outcome assertions}, run by a real agent and graded on outcomes, with at least one negative case (a broken durable outcome reds the grade).

**AC-3 — The dev template carries live-scenario-for-runtime-claims guidance as a third opt-in recommended practice.**
Verified by: a Go presence check over `skills/commission/references/templates/development.md` (the same kind that guards the existing "External-proof acceptance criteria" / "Detached adversarial audit" recommended-practice blocks) confirming a new block names the live-scenario pattern and the recording-proves-the-watcher / text-check-proves-the-words distinction. (Presence check is legitimate here: the claim is about the text itself — proof at the claim's own level.)

**AC-4 — The yy slip would now be gate-caught.**
Verified by: a fixture reproducing yy's shape (a runtime-observable AC marked PASSED with only offline proof) is REFUSED by the gate from AC-1.

## Test plan

- **Spike DONE at ideation** (see Spike section): the AC-classification + citation-check mechanism is settled — explicit `live`-lead convention + resolvable `ci-run:`/`session:` ref, gated at archival. Implementation's first test seeds from the four fixture rows there.
- **AC-1/AC-4 (Go, table-driven over fixture entity files):** drive the archival/terminal gate; assert refusal (non-zero exit + diagnostic) on an uncited runtime-observable AC, pass on a cited one and on offline ACs. AC-4 adds the yy-shape fixture row verbatim. Cost: low — fixture entities + the existing `runArchive`/`--validate` harness; no live run. Adversarial guard: re-run the spike's two break-edits (never-classify-runtime; refuse-offline-too) against the real gate to confirm the suite reds, not just greens.
- **AC-3 (Go presence check):** assert the dev template's new recommended-practice block. Cost: trivial.
- **AC-2 (live-scenario exercise):** author one scenario via the promoted primitive {runbook, setup, durable-outcome assertions}, run it against a real agent through the runner adapter, grade on entity before→after + observed output; add a negative case where a deliberately broken durable outcome reds the grade. Cost: medium — needs a live credential + minutes of agent wallclock (live-gated, like 8y's scenarios). This is `at`'s live AC-1 in primitive form.
- **High-stakes → detached adversarial audit BEFORE merge.** Trigger: the gate machinery (a status/guard mutation path) **and** the shipped contract/template change — both are the high-stakes surfaces the dev template names. The audit runs read-only on a detached checkout of the merge result and tries to REFUTE: construct an entity edit that should be gate-refused and confirm it is, and confirm an offline AC survives an edit that would over-gate.

## Notes

- Motivated by yy's live-run failure (session 11): offline + 3 audit cycles passed, the live run failed — the strongest possible case for the gate. `8y` built the substrate; `at` will be the first entity that *needs* this primitive (its live AC is a scripted scenario).
- This builds the **first AC-extractor in `validateWorkflow`** (today purely structural). A future general `gate verify` / external-proof binary guard can reuse the same extractor; this entity ships the runtime-AC-aware slice of it, not the whole thing.

## Stage Report: ideation

- DONE: SPIKE FIRST — the gate's classification + citation mechanism
  Ran a throwaway Go classifier+citation prototype over 4 discriminating fixture AC bodies (incl. yy-shape) + 2 adversarial break-edits; all discriminated correctly and both break-edits red the suite. Settled: explicit `Verified by: live <ref>` convention (not keyword inference) + resolvable `ci-run:`/`session:` citation, gated at archival. Recorded in the "Spike" section with a result table.
- DONE: Land the three-layer design with each AC backed by a CODE or TEST gate
  (a) AC-1 refusal gate anchored to `runArchive`'s existing merge-hook guard + the `--validate`/`validateWorkflow` AC-extractor (which does NOT exist yet — corrected the body's prior "extend the existing cross-check" claim); (b) live-scenario primitive grounded in 8y's actual `_test.go`-only `sharedRuntimeScenario` triple (prompt/fixture-writer/assertion) — promotion = reachability outside the ensigncycle test package; (c) AC-3 lands as a third opt-in recommended practice beside External-proof / Detached-audit.
- DONE: Test plan (AC-1/AC-4 gate refusal incl. yy fixture, AC-3 presence, AC-2 live-scenario + negative case; detached-audit trigger named)
  Test plan rewritten: AC-1/AC-4 table-driven Go over fixture entities with the spike's break-edits as an adversarial guard, AC-3 Go presence check, AC-2 live-scenario run + broken-outcome negative; detached-audit trigger named as gate-machinery + shipped contract/template (both high-stakes per the dev template).

### Summary

Fleshed out the ideation with the riskiest unknown exercised first. The spike settled the classifier as an explicit `Verified by: live <ref>` convention (rejecting prose keyword inference, which over/under-gates) with a resolvable `ci-run:`/`session:` citation, gated at the archival/terminal refusal modeled on the existing merge-hook guard in `runArchive`. Corrected two factual errors in the prior body: there is no pre-existing AC-evidence cross-check to extend (this builds the first AC-extractor in `validateWorkflow`), and 8y's scenario substrate is `_test.go`-only/live-gated so promotion is a reachability change, not a new mechanism. High-stakes (gate machinery + shipped contract/template) → detached adversarial audit before merge.
