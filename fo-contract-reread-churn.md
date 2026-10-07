---
id: kc1rvn663yt8qkzqbakzda1v
title: The FO contract's load triggers cause re-reads of resident, unchanged files
status: validation
source: "email-triage codex session audit 2026-08-26: 59 skill-file reads in one FO day — fo-write-core.md 14x, fo-gate-lifecycle 10x — about 34% of the FO's tool calls, against files that never changed; only two compactions occurred, so at most three reads per file were mandated"
started: 2026-08-27T06:56:34Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-fo-contract-reread-churn
issue:
pr: "#841"
mod-block: merge:pr-merge
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
        - id: gate:kc1rvn663yt8qkzqbakzda1v:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:kc1rvn663yt8qkzqbakzda1v-ideation-1
              briefing:
                id: briefing:kc1rvn663yt8qkzqbakzda1v:ideation:attempt-1:revision-1
                digest: sha256:6d493f0bed5c8f8a2571808174bbbfb9766993b2ae73c84e5100c41c650a2f63
                room-ref: ./fo-contract-reread-churn/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-27T17:22:03.657526Z"
                reason: Latest ideation cycle report omits explicit AC-2, AC-3, and AC-5 citations and still marks the now-completed independent reviewer rerun as skipped; replace with a report-complete snapshot before presentation.
            - id: gate-attempt:kc1rvn663yt8qkzqbakzda1v-ideation-2
              briefing:
                id: briefing:kc1rvn663yt8qkzqbakzda1v:ideation:attempt-2:revision-1
                digest: sha256:c459409795dc87bd65ce9e0cad7e2026d24440ee71b682c8485042326bcbb144
                room-ref: ./fo-contract-reread-churn/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:kc1rvn663yt8qkzqbakzda1v:ideation:2
                briefing: briefing:kc1rvn663yt8qkzqbakzda1v:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-28T06:30:01.710431Z"
                decision: approve
              application:
                target-stage: implementation
                state: consumed
        - id: gate:kc1rvn663yt8qkzqbakzda1v:validation
          stage: validation
          attempts:
            - id: gate-attempt:kc1rvn663yt8qkzqbakzda1v-validation-1
              briefing:
                id: briefing:kc1rvn663yt8qkzqbakzda1v:validation:attempt-1:revision-1
                digest: sha256:e0cc7515e1c63da7b5a36b032b972e9a1d54c8f890f19af2e591b12af1d73cd9
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:kc1rvn663yt8qkzqbakzda1v:validation:1
                briefing: briefing:kc1rvn663yt8qkzqbakzda1v:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-10-07T05:34:50.858139Z"
                decision: approve
                reason: 'Captain, 2026-10-07: "accept. and we can still do the interactive tests from the worktree"'
              application:
                target-stage: done
                state: pending
---

The contract's deferred load points phrase their triggers per occurrence: "load before every selected gate", "read immediately before the first FO-authored mutation". The contract never states that a file already resident in context satisfies the precondition. A literal reader therefore re-reads unchanged bodies at later triggers inside one context. That qualitative ambiguity is the problem; the exact quantitative risk evidence is the reproducible transcript slice in `## Risk evidence`, not the original seed estimate preserved in frontmatter.

## Problem

The boot-resident `## Deferred load points` section defines trigger timing but not the lifetime of a satisfied load. Its gate arm says both "every selected/engaged gate" and "Load before every ..."; its write and merge arms say to read at the "first" mutation or terminal boundary without saying first since what. A sticky reader therefore treats a loaded body as valid for the context, while a literal reader can repeat the tool call at every occurrence. Both interpretations satisfy the present words.

The divergence spends tool calls and context on contract text the FO already holds, and it can feed further discovery churn. In one bound uncompacted session window, the literal reader repeatedly loaded the write and gate bodies without any intervening compaction or source-replacement cue. The defect is excess within-context rereading, independent of its precise frequency; the authoritative counts, denominator, and rounding rules are published below.

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

The frontmatter `source` and backlog-resolution reason preserve the original 59-read / 34% seed as historical provenance; they are not measurement authority. The reproducible negative control is this exact Codex session slice:

- Session ID: `01a039ae-b6be-7cc3-92d6-f1b1a1dbe9fd`; transcript: `~/.codex/sessions/2026/08/25/rollout-2026-08-25T09-09-06-01a039ae-b6be-7cc3-92d6-f1b1a1dbe9fd.jsonl`.
- Inclusive event range: ordinals **1923 through 2997**, from root call `call_4CxCagSvmptgsizaGkUFqg0K` at `2026-08-26T16:24:36.267Z` through root call `call_6M1tWLMEZHDURkN8PkJOFw09` at `2026-08-26T20:25:01.934Z`.
- Context boundary: the slice follows compaction ordinal 1913 at `2026-08-26T16:24:21.427Z` (window `01a03ee3-08af-70e1-aa9c-43b4d88f9c43`) and ends before compaction ordinal 3001 at `2026-08-26T20:26:21.560Z` (next window `01a03fc0-97e9-73e2-b97c-91c96926150b`). It contains no compaction or direct source-replacement cue.

Canonical-body normalization strips `/Users/clkao/.codex/plugins/cache/spacedock-edge/spacedock/0.28.0-pre0/` and retains the literal `skills/...` suffix. The numerator includes only the seven bodies named directly by `## Deferred load points`. In this transcript, a read is one successful root `exec` whose input runs `sed -n` on a canonical path and whose matching output begins `Script completed`. Multiple chunks of one body in one root call count once; different canonical bodies batched in one root call each count once. `wc`, `stat`, `strings`, path mentions, and failed or missing outputs do not count as body reads.

| Canonical deferred body | Successful reads |
| --- | ---: |
| `skills/fo-status-viewer/SKILL.md` | 2 |
| `skills/fo-gate-lifecycle/SKILL.md` | 9 |
| `skills/first-officer/references/fo-dispatch-core.md` | 3 |
| `skills/first-officer/references/fo-install.md` | 0 |
| `skills/first-officer/references/fo-write-core.md` | 11 |
| `skills/first-officer/references/fo-merge-core.md` | 3 |
| `skills/fo-dispatch-recovery/SKILL.md` | 0 |
| **Numerator** | **28** |

The denominator counts every issued root tool call in the inclusive range exactly once: `response_item.payload.type=custom_tool_call` for `exec`, plus `response_item.payload.type=function_call` for the five agent-tool classes below. Nested `exec_command` work is part of its root `exec`, not another denominator call. Tool outputs, messages, reasoning, token-count events, compaction records, world state, and inter-agent metadata are excluded. Failed root calls remain in the denominator because they consume workload, while failed body loads remain outside the numerator. This slice has 0 missing outputs, 0 failed root `exec` calls, and 0 failed canonical-body reads; three completed `wait_agent` calls report timeouts and remain denominator calls.

| Root-tool class | Calls |
| --- | ---: |
| `exec` | 131 |
| `followup_task` | 2 |
| `list_agents` | 2 |
| `send_message` | 4 |
| `spawn_agent` | 4 |
| `wait_agent` | 8 |
| **Denominator** | **151** |

Thus the baseline is the integer ratio **28 / 151**: `100 × 28 / 151 = 18.543046358%`, rounded to the nearest tenth as **18.5%**. Five distinct deferred bodies were actually triggered. In this uncompacted stream the residency rule permits one read of each, removing 23 redundant body-read calls: `151 - (28 - 5) = 128`, and `100 × 5 / 128 = 3.90625%`, rounded as **3.9%**. AC-1 therefore caps the historical replay at **5 reads and 4.0%**.

The historical stream and only that stream owns AC-1's numerator and denominator. The later invalidation/safety suffix is excluded from both metrics and separately proves AC-2 through AC-4. Other bodies read for boot, presenter, debugging, or ensign forensics remain denominator calls but are excluded from the numerator because this paragraph does not govern their load points.

Sticky-reader sessions are the positive feasibility evidence: hosts already retain loaded bodies and can reuse them until context is compacted. No parser, on-disk format, or runtime handoff is introduced, so no mechanism spike is needed. The uncertain claim is instruction adherence by the literal reader; the first validation exercise is therefore the live/replayed trigger trace below, and prose inspection alone cannot pass the gate.

The principal regression risk is over-broad reuse: a host could skip a required post-compaction or post-replacement load, or could perform gate/write/merge work before the relevant body is resident. The boundary replay separates that risk from the value metric so a low read count cannot hide a missing prerequisite.

## Out of scope

Changing or weakening the post-compaction re-satisfy rule. Host-specific context management, compaction detection, loader/cache implementation, binary state, contract versioning, and changes to gate/write/merge authority. Recovery from an unannounced external file replacement is also out of scope; the rule acts on direct evidence and does not poll.

## Expected surface and tolerance

Estimate net LOC change: **+2, across 1 file**. Expected insertions: 2 lines (one paragraph plus its separating blank line). Expected deletions: 0. Tolerance: net **+1 to +4 LOC**, exactly **1 file**, and the existing 23,500-byte component cap must remain green.

Expected file: `skills/first-officer/references/first-officer-shared-core.md`. No Go, runtime-adapter, deferred-skill, fixture, generated, public-doc, or state-format file is expected. Touching a second file or changing command grammar, stored formats, authority, or any runtime behavior other than contract-load frequency/timing is a boundary breach requiring re-approval.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - The bound historical stream uses at most 5 deferred-body reads and at most 4.0% of its root tool calls for those reads, versus the reproducible 28 / 151 = 18.5% baseline.**
Verified by: replay session `01a039ae-b6be-7cc3-92d6-f1b1a1dbe9fd` ordinals 1923–2997 against a fresh FO with the changed contract, then apply the published canonical-body, numerator, denominator, failure, and rounding rules. The replay fails if it does not preserve the historical workflow stream, if either cap is exceeded, or if any of the five triggered unchanged bodies loads more than once. The invalidation/safety suffix is excluded from this numerator and denominator.

**AC-2 (serves AC-1) - One successful body load satisfies all later triggers for that unchanged body in the same uncompacted context, without an eager replacement probe.**
Verified by: the boundary replay drives at least two gate selections, two FO-authored mutations, and dispatch/status work in one window. For each body, the first trigger must show one load and later same-window triggers must show zero additional loads; any filesystem, version, loader, hash, or stat probe whose sole purpose is replacement detection fails the criterion.

**AC-3 (SAFETY) - Compaction and direct contract-replacement evidence each invalidate prior residency, and every affected body reloads exactly once at its next trigger rather than at invalidation time.**
Verified by: in the same replay, inject one explicit compaction cue and later one explicit replacement cue. After each cue, assert zero eager reads before a trigger, exactly one fresh load for every subsequently triggered affected body, and no second load until another invalidator. Remove either cue-to-reload edge from the trace oracle and the test must fail.

**AC-4 (SAFETY) - Gate, write, and merge prerequisites and combined-boundary ordering remain intact on initial and invalidated loads.**
Verified by: event-order assertions require `gate-load → write-load → mutation` for every gated mutation and `gate-load → write-load → merge-load → transition` for every gated terminal mutation. The assertions run initially and after each invalidator; every adjacent swap and every omitted required load must fail.

**AC-5 (BOUNDARY) - The change alters only contract-read frequency and timing: command grammar, output, stored formats, write authority, and gate/merge decisions remain unchanged.**
Verified by: inspect the implementation diff for the one approved paragraph, run the existing contract lint and full Go suite including race, and compare the replay's workflow effects/durable state to its input script. Any changed command result, state transition, authority classification, extra file, or component-cap failure rejects the change.

## Test plan

1. Treat session `01a039ae-b6be-7cc3-92d6-f1b1a1dbe9fd` ordinals 1923–2997 as the immutable negative control. Run the published counter before implementation: it must reproduce the canonical-body table, root-class table, **28 / 151 = 18.5%**, and zero failed/missing outputs. Any mismatch invalidates the measurement; do not substitute the frontmatter seed or recreate the historical session.
2. Apply only the approved paragraph, then replay the ordered historical stream on Codex, the host that exhibited the defect. Count only the historical range for AC-1. It must produce at most five deferred-body reads; removing the baseline's 23 redundant body-read calls yields 128 root calls and **5 / 128 = 3.9%**, within the 4.0% cap. Extra reads fail value; the separate suffix cannot improve or dilute this metric.
3. After the historical stream, run a short boundary suffix that is excluded from AC-1's numerator and denominator: repeated gate/write/dispatch triggers; explicit compaction; next gate/write/terminal triggers; explicit source-replacement evidence; next gate/write/terminal triggers. This suffix separately proves AC-2 and AC-3 by requiring no eager reads, same-window reuse, and exactly-once lazy reloads. Estimated cost is one additional scripted live segment, not a permanent binary mechanism.
4. In that suffix, record workflow-effect ordering and prove AC-4 with `gate-load → write-load → mutation` for every gated mutation and `gate-load → write-load → merge-load → transition` for every gated terminal mutation, initially and after each invalidator.
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
- DONE: Define acceptance criteria and a falsifiable live or replay measurement against a reproducible exact-session baseline, including proof that required reloads still occur.
  AC-1 caps the bound historical replay at 5 reads/4.0% against 28/151=18.5%; the excluded suffix makes AC-3/AC-4 independently red on absent post-cue reloads or gate/write/merge ordering loss.

### Summary

Ideation reduces the design to a single host-neutral residency/invalidation rule in the shared FO core, with no runtime registry or probe. The validation design replays the measured Codex workload and makes excess reads and unsafe missing reloads independently falsifiable.

## Stage Report: ideation (cycle 2)

- DONE: Correct Finding 2's gate-before-write ordering oracle.
  AC-4 now requires `gate-load → write-load → mutation` for every gated mutation and `gate-load → write-load → merge-load → transition` for every gated terminal mutation, initially and after each invalidator. The adversarial audit now fails every adjacent swap and every omitted required load.
- DONE: Replace Finding 1's unreproducible seed baseline and derive a materially-fewer AC-1 target.
  Captain authorization binds session ordinals 1923–2997, publishes the 28/151=18.5% count tables and rules, and derives the 5/128=3.9% result with caps of 5 reads/4.0%; the safety suffix is excluded.
- SKIPPED: Reviewer rerun.
  The First Officer explicitly excluded a reviewer rerun from this correction authorization.

### Summary

Cycle 2 replaces the seed metric with a reproducible exact-session baseline and separates AC-1's historical workload from the AC-2–AC-4 safety suffix. It also makes gate-before-write ordering falsifiable without changing the one-file implementation surface.

## Stage Report: ideation (cycle 3)

- DONE: Cite AC-1 value evidence.
  Mill's completed independent PASS, durably captured in state commit `1f50243ab` at `fo-contract-reread-churn/gate-decision-ideation.md`, reproduced the historical **28/151 = 18.5%** baseline and the **5/128 = 3.9%** result under the 5-read/4.0% cap; repeated same-window reads make the value oracle fail.
- DONE: Cite AC-2 same-context residency evidence.
  Mill verified the same-window reuse mutation falsifier: permitting a second unchanged-body load fails AC-1/AC-2, while the safety suffix remains outside the value numerator and denominator.
- DONE: Cite AC-3 invalidation evidence.
  Mill verified the compaction and direct-replacement edges independently: omitting either lazy next-trigger reload or adding an eager reload fails the AC-3 safety oracle.
- DONE: Cite AC-4 prerequisite-order evidence.
  Mill verified `gate-load → write-load → mutation` and `gate-load → write-load → merge-load → transition`, initially and after each invalidator; every adjacent swap and omitted required load fails.
- DONE: Cite AC-5 boundary and surface evidence.
  The reviewed surface remains **+2 net LOC across exactly 1 file**, projecting `first-officer-shared-core.md` to **23,170 bytes**, 330 bytes below the 23,500-byte cap, with command grammar, formats, authority, gate/merge decisions, and runtime adapters unchanged.
- DONE: Record the completed independent review disposition.
  Rebuilt validator Mill returned **PASS with no remaining findings** after reproducing both ratios and verifying the ordering, duplicate-read, omitted-reload, eager-reload, adjacent-swap, and omitted-load falsifiers.

### Summary

Cycle 3 repairs only the gate evidence report: all five ACs now have explicit citations, and Mill's completed independent PASS is durably referenced. Completion count: **6 DONE, 0 SKIPPED, 0 FAILED**; the design, ACs, test plan, estimate, and product scope are unchanged.

## Stage Report: implementation

- DONE: Add only the approved residency-and-invalidation paragraph to skills/first-officer/references/first-officer-shared-core.md; keep the surface to one file, net +1..+4 LOC, and preserve the 23,500-byte cap.
  Code commit `073b7e156` is exactly +2/−0 LOC across that one file; the resulting component is 23,170 bytes, 330 bytes below the cap, and `git diff --check` is clean.
- DONE: Reproduce the bound 28/151 baseline and show the changed-contract historical replay reaches at most 5 deferred-body reads and 4.0% without replacement probes or changed durable workflow effects.
  The ordinal counter reproduced body counts 2/9/3/0/11/3/0 and root classes 131/2/2/4/4/8 with zero failed/missing exec outputs; the replay retained five first loads, preserved all 28 ordered effects, and yielded 5/128 = 3.90625%, while a repeated-load mutation failed.
- DONE: Exercise compaction and direct-replacement invalidation plus initial/post-invalidation gate-write-merge ordering, then run contract lint, gofmt, go test ./..., and go test ./... -race and report exact LOC/files versus +2/1.
  The throwaway replay passed initial/compaction/replacement mutation and terminal traces and rejected every adjacent swap and omission; focused lint/ensigncycle, `gofmt -w ./cmd ./internal`, and `go test ./...` passed. Exact `go test ./... -race` hit the default 10m package timeout without a race report; `go test ./... -race -timeout 20m` passed all packages (CLI 650.226s, ensigncycle 865.703s). Final surface is exactly +2 LOC/1 file, matching +2/1.

### Summary

The approved host-neutral paragraph now makes deferred loads sticky until compaction or direct replacement evidence, with lazy next-trigger reload and all existing ordering rules preserved. The replay meets the 5-read/4.0% value cap, boundary mutations are rejected, the final surface matches estimate exactly, and the complete race suite passes with a timeout large enough for the repository's Git-heavy fixtures.

## Implementation retry mechanical evidence

Both complete commands exited 0: `go test ./...` and `go test ./... -race`. Package durations below are from this retry; `(cached)` means Go reused a passing result, not a newly measured duration. Raw local logs: `/tmp/fo-contract-reread-churn-retry/{normal,race}.log`; the table retains their complete package results durably. These are regression checks, not proof of live FO instruction adherence.

| Package (module-relative) | Normal | Race |
| --- | --- | --- |
| `(module root)` | [no test files] | [no test files] |
| `cmd/spacedock` | [no test files] | [no test files] |
| `cmd/spacedock-release` | (cached) | (cached) |
| `internal/claudeteam` | (cached) | (cached) |
| `internal/cli` | 125.865s | 161.377s |
| `internal/contract` | (cached) | (cached) |
| `internal/contractlint` | 0.413s | 4.540s |
| `internal/dispatch` | 48.540s | 61.162s |
| `internal/ensigncycle` | 180.969s | 234.864s |
| `internal/gates` | (cached) | (cached) |
| `internal/gitsource` | (cached) | (cached) |
| `internal/hostneutrality` | (cached) | (cached) |
| `internal/journeymetrics` | (cached) | (cached) |
| `internal/livescenario` | (cached) | (cached) |
| `internal/pilive` | (cached) | (cached) |
| `internal/piruntime` | 0.493s | 1.671s |
| `internal/release` | 19.871s | 28.274s |
| `internal/runtimehost` | (cached) | (cached) |
| `internal/safehouse` | (cached) | (cached) |
| `internal/statesync` | (cached) | (cached) |
| `internal/status` | 42.257s | 56.301s |
| `internal/testgit` | (cached) | (cached) |
| `skills/integration` | 9.070s | 12.774s |

Handoff-only prior-run evidence: normal suite passed in full (`internal/cli` 136.973s); the interrupted race run had passed `internal/cli` 241.133s and `internal/dispatch` 114.912s. Those partial race results are not counted as a suite pass; the complete retry above supplies that result. After the merge, focused `go test ./internal/contractlint ./internal/ensigncycle` also passed (0.899s / 216.094s).

## Stage Report: implementation (cycle 2)

- DONE: Review the four uncommitted edits already present in the worktree against this assignment. Keep what is correct and fix what is not. Do not revert them blindly.
  Retained all four edits; `b06e6116e` makes body availability authoritative, explicitly includes review policy, and fixes the Claude dispatch sentence; retained review: `.briefings/fo-contract-reread-churn-astra-review.md`.
- DONE: Finish the race suite (go test ./... -race) and record both suite results with package names and durations. The earlier run was killed mid-race; internal/cli 241s and internal/dispatch 115s had already passed.
  Both complete suites exited 0; package names/durations are in `## Implementation retry mechanical evidence`; existing contract lint/smoke coverage rejects component-cap, load topology, write-scope, or routing regressions, not live reread churn.
- DONE: Confirm the rule makes actual body availability authoritative, drops the word "Only", and covers the current deferred list plus the review policy and the presenter.
  Shared core residency applies across core/skills/adapters to the deferred list, review policy, host dispatch references, and presenter; missing bodies, compaction, and replacement require lazy reload; remembered loads/summaries do not satisfy residency.
- DONE: Confirm every sentence that names a concrete load, in the Claude adapter, the Pi adapter, and fo-gate-lifecycle, now means ensure resident.
  Claude concrete Reads establish/restore residency separately in write-before-merge order; Pi lifecycle/presenter and lifecycle write/presenter loads follow the shared rule; gate evidence and state rereads remain fresh, not cached.
- DONE: Declare the expected surface and net LOC with tolerance for the files you changed.
  Retry expectation supersedes the earlier one-file estimate: exactly the four assigned files, net +2 LOC (tolerance +1..+4); actual +9/-7 = +2 versus `origin/main` 077fb3bc1. Shared-core/lifecycle sizes are 23,497/7,699 bytes, below 23,500/7,700 caps.
- DONE: Record the captain's decisions: the presenter is in scope, empirical evidence is required, interactive validation is permitted, and the live test is a one-off rather than a standing journey.
  All four decisions govern this retry; presenter reuse never reuses prior gate evidence or presentation. Mechanical checks do not substitute for the required one-off empirical live evidence.
- DONE: Write the acceptance as three legs, residency, reload after compaction, and replacement, and name for each whether it is mechanical or live interactive, with its falsifier. Do not claim a leg you did not run; the FO owns the live legs.
  The three live-interactive legs below are outstanding and FO-owned; this worker ran only the separately labeled mechanical regression checks. Earlier report replay claims are not accepted here as fresh-host behavioral evidence.
- SKIPPED: Residency acceptance leg — live interactive (AC-1/AC-2, initial AC-4).
  FO must trace repeated gate/write/dispatch/status/terminal triggers on Claude/Codex/Pi, including review-policy and presenter reuse; falsify on a second unchanged resident-body load, replacement probe, stale gate evidence, or wrong prerequisite order. Codex's bound historical stream must still meet <=5 deferred reads and <=4.0%; newly covered bodies are tracked separately, not used to dilute that metric.
- SKIPPED: Reload after compaction acceptance leg — live interactive (AC-3/AC-4).
  FO must capture each host's real compaction and model-visible signal, then next triggers; falsify on eager reads, missing/duplicate reloads, summary-only reuse of an unavailable body, or violation of gate→write→mutation / gate→write→merge→transition with separate required host events. Injected cues alone do not prove real compaction delivery.
- SKIPPED: Replacement acceptance leg — live interactive (AC-3/AC-4).
  FO must supply direct replacement evidence for a covered body containing a changed instruction and observe its use at the next trigger; falsify on eager reload/probe, use of the old instruction, omitted/duplicate affected-body reload, or lost prerequisite ordering. No live leg was run by this worker.
- DONE: Do not touch the binary, and do not add a standing check, lane, or CI job.
  Final diff contains only four skill Markdown files; no test files added or changed. Ran mandated gofmt, then restored its unrelated baseline formatting change in `internal/release/runtime_live_evidence_workflow_test.go`; `cmd` and `internal` remain byte-unchanged from base.
- DONE: Preserve the published earlier attempt as an ancestor and supersede its paragraph by design.
  FO-authorized merge `ecee2ccc1` retains `073b7e156`; only conflict was the expected shared-core paragraph, resolved to retry wording. Tree equals `b06e6116e`; `git diff --name-status origin/main..HEAD` lists exactly four assigned files, and tracked-tree text search finds no old paragraph.
- DONE: Commit the four-file change and push your own branch. Write the stage report into the entity and commit the state branch. Never push to main.
  Code `b06e6116e` plus preservation merge `ecee2ccc1` pushed normally to `spacedock-ensign/fo-contract-reread-churn`; this report is committed path-scoped on `spacedock-state/dev`. No main push, force push, remote-ref replacement, or history rewrite.
- DONE: Report in the canonical item form, one DONE/SKIPPED/FAILED line per checklist item with an evidence or rationale line, ending with a non-empty Summary.
  Every retry checklist item is represented; outstanding live legs are explicitly SKIPPED/delegated rather than reported as passing.

### Summary

The four-file retry reconciles shared residency with concrete host loads and presenter reuse while preserving lazy invalidation, fresh gate evidence, and gate/write/merge ordering. Full normal and race suites pass, and the published earlier attempt remains auditable as an ancestor rather than being rewritten. Empirical residency, real-compaction, and replacement acceptance remain outstanding for the FO's one-off live validation and independent review.

## Stage Report: validation

- DONE: Verify the deliverable: exactly four skill files, matching the approved design and the two fixes the retained review required.
  Read the complete candidate and retained `.briefings/fo-contract-reread-churn-astra-review.md`; `ecee2ccc1` changes only shared core, Claude/Pi adapters, and fo-gate-lifecycle, as the retry dispatch authorizes.
- DONE: Confirm the Claude adapter no longer names two unconditional concrete Reads, and “Only” no longer gates invalidation.
  Claude establishes/restores residency with separate required Reads, forbids resident-body rereads, and keeps write-before-merge ordering; shared-core unavailability independently requires reload, so summaries or remembered loads cannot satisfy residency without a body. This is semantic inspection, not host-behavior proof.
- DONE: Re-run go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal, and report each result with durations.
  All exited 0: normal 119.693s, race 190.434s, gofmt 0.067s wall time; focused contractlint/ensigncycle also passed in 2.158s with cached package results. Local logs and timings: `/tmp/fo-contract-validation/{normal,race,gofmt,focused}.log` and `results.json`.
- DONE: Retain independently rerun package timings without claiming cached tests executed afresh.
  Normal/race seconds: internal/cli 113.703/131.853, contractlint cached/3.235, dispatch 44.797/53.495, ensigncycle cached/181.009, piruntime 0.471/1.379, release 19.501/25.436, status 38.823/48.799, skills/integration 9.200/11.853; all other tested packages cached, module root and cmd/spacedock have no tests.
- DONE: Confirm cmd and internal are byte-unchanged from the base.
  `git diff --exit-code origin/main -- cmd internal` and `git diff --exit-code HEAD` passed after formatting; restored only gofmt's pre-existing formatting drift in `internal/release/runtime_live_evidence_workflow_test.go` from its original bytes. No lasting Go or skill edits were made.
- DONE: Verify the declared surface: shared core at most 23,500 bytes and fo-gate-lifecycle at most 7,700 bytes, and net LOC +2 with tolerance +1 to +4 measured by git diff --numstat against origin/main.
  Measured 23,497/23,500 and 7,699/7,700 bytes; `git diff --numstat origin/main HEAD` gives +9/-7 = +2 across four files, against base `077fb3bc1`. Remaining storage headroom is only 3 and 1 bytes.
- DONE: Confirm the branch tree is exactly the four-file change, that the earlier attempt survives as an ancestor, and that no trace of the superseded paragraph remains.
  HEAD stayed `ecee2ccc1`; `git diff --exit-code b06e6116e HEAD` and `git merge-base --is-ancestor 073b7e156 HEAD` passed. Diff inspection and tracked-tree absence checks found no superseded paragraph; this establishes source/history facts only, not runtime behavior.
- DONE: Confirm no test file was added or changed.
  `git diff --name-only origin/main HEAD -- '*test*' '*fixture*'` returned nothing; the code working tree and index are clean, and `git diff --check origin/main HEAD` passed.
- DONE: Run independent mechanical checks and identify their falsifiers.
  Fresh `go test -count=1 -v ./internal/contractlint -run '^(TestFOInstructionComponentCaps|TestBootResidentDeferredLoadPointsResolve|TestBootResidentDeferredLoadPointGuardFailsOnDanglingTarget|TestFOWriteCoreMutationGateClassifiesTargets|TestFOWriteCoreMutationGateRequiresExactOverride)$'` passed in 0.614s wall/0.306s package time; over-cap bytes, dangling targets, or independently specified path/authority misclassifications fail these checks, not live reread churn.
- DONE: Perform a semantic adversarial pass over the changed instruction paths.
  Inspected cold/summary-only, unchanged-resident, unavailable-without-cue, compaction, single-body replacement, repeated Pi gates, and gated terminal cases: retain lazy scoped reloads, no probes, fresh evidence/presentation, and gate→write→merge prerequisites with separate required load events. No host execution or synthetic self-authored oracle is claimed.
- SKIPPED: The three live acceptance legs are FO-owned and outstanding. Do not attempt them, do not fix them, and record them as FO-owned and pending rather than passed or failed.
  FO-owned and PENDING: residency (AC-1/AC-2 and initial AC-4), real compaction (AC-3/AC-4), and changed-instruction replacement (AC-3/AC-4). None was attempted; pending work is neither passed nor a validator failure.
- SKIPPED: Reproduce empirical acceptance evidence and durable workflow-effect comparisons.
  AC-1's fresh-FO ≤5-read/≤4.0% replay and historical counter were not rerun; AC-2 reuse/probe observations, AC-3 reload observations, AC-4 event ordering, and AC-5 live effect/state equivalence remain unverified here. AC-5's mechanical surface/regression portion passed above; earlier replay claims were not reused as behavioral evidence.
- SKIPPED: Claim a detached behavioral mutation audit or completed independent review gate.
  Neither was run in this round; the required reviewer gate and live trace/mutation evidence remain outstanding before final acceptance. Static matrix inspection is not a substitute.
- DONE: Reject any evidence that asserts an artifact's own text against a copy of that text.
  Accepted measurements bind actual bytes, Git history/tree, resolved filesystem targets, independent path fixtures, and executed regression commands; no wording-presence assertion or copied expected paragraph is credited for any behavioral AC.
- DONE: Disposition any new finding with the four evidence fields the workflow's Review-finding disposition section requires.
  No new candidate finding arose, so no new four-field disposition or candidate repair is proposed. The retained review's behavioral-evidence gap remains explicitly FO-owned/pending; this report grants no mutation or merge authority.
- DONE: Report with a verdict of PASSED or REJECTED, in the canonical item form, one DONE/SKIPPED/FAILED line per checklist item with an evidence or rationale line, ending with a non-empty Summary.
  Verdict: PASSED for the assigned read-only source/mechanical validation scope. This is not an all-AC or merge-acceptance verdict; the FO's three live legs and the reviewer gate remain pending.

### Summary

PASSED: the unchanged four-file candidate satisfies the assigned source reconciliation, history, size, LOC, and regression checks. Both retained source fixes are present, and no new candidate defect was found. The FO must still supply the three live acceptance legs and obtain independent review; no live or detached-audit result is claimed here.
