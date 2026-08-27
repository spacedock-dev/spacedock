---
id: kc1rvn663yt8qkzqbakzda1v
title: The FO contract's load triggers cause re-reads of resident, unchanged files
status: ideation
source: "email-triage codex session audit 2026-08-26: 59 skill-file reads in one FO day — fo-write-core.md 14x, fo-gate-lifecycle 10x — about 34% of the FO's tool calls, against files that never changed; only two compactions occurred, so at most three reads per file were mandated"
started: 2026-08-27T06:56:34Z
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:kc1rvn663yt8qkzqbakzda1v:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:kc1rvn663yt8qkzqbakzda1v-backlog-1
              briefing:
                id: briefing:kc1rvn663yt8qkzqbakzda1v:backlog:attempt-1:revision-1
                digest: sha256:a562429be9ee66eddf15442aa4bea797111178a583f49f23cd68bc9b8a45379d
                room-ref: ./fo-contract-reread-churn/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:kc1rvn663yt8qkzqbakzda1v:backlog:1
                briefing: briefing:kc1rvn663yt8qkzqbakzda1v:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T06:55:14.394209Z"
                decision: approve
                reason: 'Captain approved the bound Subspace backlog review: the 59-read baseline justifies ideation while preserving gate, write, merge, and post-compaction reload boundaries.'
              application:
                target-stage: ideation
                state: consumed
---

The contract's deferred load points phrase their triggers per occurrence: "load before every selected gate", "read immediately before the first FO-authored mutation". The contract never states that a file already resident in context satisfies the precondition. A literal reader therefore re-reads the full file at every trigger occurrence. One measured day: 14 reads of the same unchanged write-core, roughly one per mutation boundary.

## Problem

The boot-resident `## Deferred load points` section defines trigger timing but not the lifetime of a satisfied load. Its gate arm says both "every selected/engaged gate" and "Load before every ..."; its write and merge arms say to read at the "first" mutation or terminal boundary without saying first since what. A sticky reader therefore treats a loaded body as valid for the context, while a literal reader can repeat the tool call at every occurrence. Both interpretations satisfy the present words.

The divergence is expensive and can feed further discovery churn. The audited 2026-08-26 Codex FO day made 59 skill/contract reads (34% of tool calls), including 14 reads of `fo-write-core.md` and 10 of `fo-gate-lifecycle`, although those sources did not change. Only two compactions occurred, so an unchanged body needed at most one load in each of three context windows. Each repeated read brings roughly 200 lines back through the tool surface; in the same session, one failed command flag led to `--help` and `strings` fallback discovery because the contract offers no explicit, cheap answer to "is the body I already hold still valid?"

The correct safety boundary must remain: after compaction, and after known replacement of the loaded contract, the next relevant workflow effect must use a fresh body. Gate, write, and merge actions must also keep their existing prerequisite and ordering rules. The defect is only the absence of an explicit validity interval between those boundaries.

## Proposed approach

Add one host-neutral paragraph to `skills/first-officer/references/first-officer-shared-core.md`, directly after the opening paragraph of `## Deferred load points`. Do not add a registry, marker file, hash probe, runtime-adapter rule, or binary state. The paragraph makes every existing trigger an ensure-resident precondition and gives that precondition two explicit invalidators.

Concrete contract diff:

```diff
 A greet-and-stop boot loads NONE of these — it composes its summary from `«state.boot»()` and follows the interactive branch of `«interaction.boundary»()`. Each loads only at its trigger:
+**Residency and invalidation.** At every trigger below, the named body must be resident before the listed effect. One successful load satisfies later triggers for that body in the same uncompacted context; `load` and `read` below mean ensure resident, not repeat a tool call. Only a harness notice or captain cue of compaction, or direct evidence that the loaded source was replaced, invalidates it. After invalidation, reload at the next existing trigger—never eagerly—and preserve that trigger's ordering and own-host-event requirements. Do not probe the filesystem, version, or loader merely to look for replacement.
+
 **Combined-boundary order:** evaluate the write trigger before the merge trigger.
```

"Direct evidence" means evidence already delivered by the environment or the captain, such as a completed plugin/skill replacement or an explicit source-change notice. Absence of evidence preserves validity; the FO does not poll for change. Invalidation is per loaded body, while compaction invalidates all context-resident bodies. A body reloads once, lazily, at its next existing trigger; invalidation itself does not cause a read.

This single rule leaves every trigger-specific sentence operative. A selected gate still requires `fo-gate-lifecycle` resident before capability, evidence, Git, presenter, decision, replay, write, or dispatch work. A first write after boot or invalidation still loads `fo-write-core.md` in its own completed host event, after gate lifecycle when both apply. A first terminal boundary after boot or invalidation still loads write before merge, each in its own host event, before the transition. Later same-window actions reuse those unchanged bodies. The existing `## Compaction continuity` sentence remains unchanged and now has a precise meaning: compaction invalidates prior satisfaction, and each body re-satisfies at its next trigger.

No public documentation diff is proposed. This is internal FO execution semantics: command grammar, output, stored formats, and authority do not change; only redundant contract-read tool calls are removed. The shipped contract file is the user-facing source of this behavior.

Cheaper alternatives are insufficient:

- Editing only the gate bullet to say "first gate" leaves dispatch, write, merge, status, and failure-triggered bodies ambiguous and still leaves "first since what" undefined.
- Adding "once per session" is smaller but unsafe after compaction or a known contract replacement.
- Adding reminders to individual host adapters duplicates policy, preserves host divergence, and misses host-neutral reference-file reads.
- Keeping an explicit loaded-set in the binary or on disk cannot observe what remains in an agent's context, creates stale state across compaction, and is machinery for a prose validity rule.
- Hashing, statting, or probing the loader at each trigger converts read churn into probe churn and can inspect a different source than the loader-bound body. Replacement evidence already in hand is sufficient.

## Risk evidence

The 2026-08-26 Codex FO audit is both the negative control and baseline: 59 reads, 34% read share, per-file repetition against unchanged bytes, and two compaction timestamps. Sticky-reader sessions are the positive feasibility evidence: hosts already retain loaded bodies and can reuse them until context is compacted. No parser, on-disk format, or runtime handoff is introduced, so no mechanism spike is needed. The uncertain claim is instruction adherence by the literal reader; the first validation exercise is therefore the live/replayed trigger trace below, and prose inspection alone cannot pass the gate.

The principal regression risk is over-broad reuse: a host could skip a required post-compaction or post-replacement load, or could perform gate/write/merge work before the relevant body is resident. The boundary replay separates that risk from the value metric so a low read count cannot hide a missing prerequisite.

## Out of scope

Changing or weakening the post-compaction re-satisfy rule. Host-specific context management, compaction detection, loader/cache implementation, binary state, contract versioning, and changes to gate/write/merge authority. Recovery from an unannounced external file replacement is also out of scope; the rule acts on direct evidence and does not poll.

## Expected surface and tolerance

Estimate net LOC change: **+2, across 1 file**. Expected insertions: 2 lines (one paragraph plus its separating blank line). Expected deletions: 0. Tolerance: net **+1 to +4 LOC**, exactly **1 file**, and the existing 23,500-byte component cap must remain green.

Expected file: `skills/first-officer/references/first-officer-shared-core.md`. No Go, runtime-adapter, deferred-skill, fixture, generated, public-doc, or state-format file is expected. Touching a second file or changing command grammar, stored formats, authority, or any runtime behavior other than contract-load frequency/timing is a boundary breach requiring re-approval.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - The audited FO-day workload over unchanged contract sources uses at most 29 contract reads and at most 21% of all tool calls for contract reads, versus the fixed 59-read / 34% baseline.**
Verified by: replay the audit's ordered workflow triggers and its two compaction cues against a fresh FO with the changed contract, then count successful loads by canonical body identity from the host transcript (a batched command counts once per body). The replay fails if either cap is exceeded, if any unchanged body loads more than once inside one of the three context windows, or if the original baseline artifact is rewritten instead of compared.

**AC-2 (serves AC-1) - One successful body load satisfies all later triggers for that unchanged body in the same uncompacted context, without an eager replacement probe.**
Verified by: the boundary replay drives at least two gate selections, two FO-authored mutations, and dispatch/status work in one window. For each body, the first trigger must show one load and later same-window triggers must show zero additional loads; any filesystem, version, loader, hash, or stat probe whose sole purpose is replacement detection fails the criterion.

**AC-3 (SAFETY) - Compaction and direct contract-replacement evidence each invalidate prior residency, and every affected body reloads exactly once at its next trigger rather than at invalidation time.**
Verified by: in the same replay, inject one explicit compaction cue and later one explicit replacement cue. After each cue, assert zero eager reads before a trigger, exactly one fresh load for every subsequently triggered affected body, and no second load until another invalidator. Remove either cue-to-reload edge from the trace oracle and the test must fail.

**AC-4 (SAFETY) - Gate, write, and merge prerequisites and combined-boundary ordering remain intact on initial and invalidated loads.**
Verified by: event-order assertions require `gate-load → write-load → mutation` for every gated mutation and `gate-load → write-load → merge-load → transition` for every gated terminal mutation. The assertions run initially and after each invalidator; every adjacent swap and every omitted required load must fail.

**AC-5 (BOUNDARY) - The change alters only contract-read frequency and timing: command grammar, output, stored formats, write authority, and gate/merge decisions remain unchanged.**
Verified by: inspect the implementation diff for the one approved paragraph, run the existing contract lint and full Go suite including race, and compare the replay's workflow effects/durable state to its input script. Any changed command result, state transition, authority classification, extra file, or component-cap failure rejects the change.

## Test plan

1. Preserve the 2026-08-26 audit transcript/count table as the immutable 59-read / 34% negative control. Derive the ordered trigger stream and the two existing compaction cues from it; do not derive expected post-change counts from the candidate prose. The 29-read cap is a greater-than-50% reduction, and 21% is the corresponding conservative read-share ceiling after redundant read calls disappear from the denominator.
2. Exercise the counter on the original transcript before editing: it must reproduce 59 reads and approximately 34%, or the measurement is invalid. Do not rerun the old FO day merely to recreate evidence already preserved by the audit.
3. Apply only the approved paragraph, then replay on Codex, the host that exhibited the defect. Record canonical-body reads, total tool calls, window boundaries, replacement cue, and workflow-effect ordering. This one trace proves AC-1 through AC-4 and is falsifiable in both directions: extra reads fail value; missing or misordered reads fail safety.
4. Use a short boundary suffix after the historical stream: repeated gate/write/dispatch triggers; explicit compaction; next gate/write/terminal triggers; explicit source-replacement evidence; next gate/write/terminal triggers. Require no eager reads and exactly-once lazy reloads. Estimated cost is one additional scripted live segment, not a permanent binary mechanism.
5. Run focused `go test ./internal/contractlint ./internal/ensigncycle`, then `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` with a clean diff check. Existing topology, component-cap, write-scope, gate, and routing tests cover AC-5; no static wording-presence test is added because it would prove only that the means shipped.
6. Because shipped contract/scaffolding is a high-stakes surface, run the workflow's detached adversarial audit on a throwaway checkout. Challenge the evidence with three mutation classes: permit repeated same-window loads (AC-1/2 must fail), remove a compaction/replacement reload edge (AC-3 must fail), and exercise every adjacent swap and every omitted required load in both AC-4 sequences, initially and after each invalidator. A green replay under any mutation is a validation hole, not a pass.

Test complexity is moderate because one live Codex replay is required; repository tests are low-cost and unchanged. No Claude/Pi live matrix is required: the semantics are host-neutral and host context management is out of scope, while the Codex negative control is the measured failing host. If the Codex replay cannot distinguish a skill invocation from a body load, validation must count the underlying tool/file-open events directly rather than weakening the metric.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Stage Report: ideation

- DONE: Select the smallest resident-contract and invalidation rule that removes within-context rereads while preserving post-compaction, contract-replacement, gate, write, and merge preconditions.
  `## Proposed approach` selects one boot-resident ensure-resident paragraph; AC-2 through AC-4 fail on duplicate reads, missing reloads, or reordered prerequisites.
- DONE: Declare concrete before/after contract wording, semantic boundaries, expected net LOC and file surface with tolerance, and why cheaper alternatives are insufficient.
  The exact one-paragraph diff, five rejected alternatives, unchanged semantic classes, and +2 LOC/one-file estimate (+1..+4 tolerance) are recorded in the body.
- DONE: Define acceptance criteria and a falsifiable live or replay measurement against the 59-read / 34% baseline, including proof that required reloads still occur.
  AC-1 caps the replay at 29 reads/21%; AC-3/AC-4 independently red on absent post-cue reloads or gate/write/merge ordering loss.

### Summary

Ideation reduces the design to a single host-neutral residency/invalidation rule in the shared FO core, with no runtime registry or probe. The validation design replays the measured Codex workload and makes excess reads and unsafe missing reloads independently falsifiable.

## Stage Report: ideation (cycle 2)

- DONE: Correct Finding 2's gate-before-write ordering oracle.
  AC-4 now requires `gate-load → write-load → mutation` for every gated mutation and `gate-load → write-load → merge-load → transition` for every gated terminal mutation, initially and after each invalidator. The adversarial audit now fails every adjacent swap and every omitted required load.
- SKIPPED: Finding 1's baseline reproducibility and threshold changes.
  The First Officer held the finding for captain decision, so the 59-read / 34% baseline and 29-read / 21% target remain unchanged.
- SKIPPED: Reviewer rerun.
  The First Officer explicitly excluded a reviewer rerun from this correction authorization.

### Summary

Cycle 2 tightens the safety oracle to make gate-before-write ordering falsifiable without changing the one-file implementation surface or the held workload metric.
