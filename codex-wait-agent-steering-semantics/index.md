---
title: Codex wait_agent steering semantics describe captain input as active-loop resumption
status: backlog
source: "Captain request 2026-07-23: replace misleading wait-interruption language and use the corrected behavior in-session"
started:
completed:
verdict:
score: 0.9
worktree: ""
issue:
id: 6gkz4z2qweheyj17ck5tythn
---

## Problem

The Codex First Officer runtime currently imports the harness label “Wait interrupted by new input” into its operating language. That label does not describe the behavior: `wait_agent` is asynchronous monitoring, captain input resumes the FO's active loop, and unresolved workers continue unchanged. Repeating an interruption disclaimer before each wait epoch makes normal steering sound destructive and adds noise.

## Required behavior

The Codex runtime contract must express these semantics directly:

- Captain input resumes the FO's active loop while workers continue unchanged.
- When the FO becomes idle again, it resumes monitoring unresolved workers.
- The contract does not require a repetitive interruption disclaimer before every wait epoch.
- A wait return or captain message never becomes worker-completion evidence; durable reports and the final-status signal remain authoritative.

Scope is the Codex runtime contract and its behavioral proof. Do not change `wait_agent`, invent cancellation/restart state, or broaden this into a generic scheduler redesign.

## Acceptance criteria

**AC-1 (VALUE)** In a Codex drive where one worker remains unresolved across captain input, the FO handles the input, continues useful active-loop work, and later monitors the same worker without failing, closing, or redispatching it. Verified by a runtime drive or fixture that correlates one worker handle before and after steering and observes its eventual durable completion.

**AC-2** The Codex runtime instruction models captain input as active-loop resumption and idle monitoring as conditional on becoming idle again, while preserving the rule that monitoring output alone is not completion. Verified by a behavioral contract scenario whose wrong-semantics variants—worker cancellation, redispatch, or treating wait return as completion—fail.

**AC-3** Two or more monitoring epochs do not require repeated captain-facing interruption disclaimers. Verified by an observed-output scenario that permits ordinary progress/idle communication but fails if the old mandatory disclaimer is emitted before each epoch.

**AC-4** The change stays Codex-specific unless ideation demonstrates a shared-core contradiction. Claude and Pi wait semantics remain unchanged, and no runtime tool/API behavior is fabricated. Verified by scoped diff review plus existing runtime contract and live-tag test gates appropriate to the touched paths.

## Stage test gates

- Ideation identifies every current clause and behavioral fixture affected, proposes an exact wording delta, and spikes the cheapest scenario that distinguishes steering from worker interruption.
- Implementation changes only the approved contract/test surfaces, runs focused and full/race gates, and requests Roborev because this is shipped runtime scaffolding.
- Validation performs the required detached adversarial audit and a Codex behavioral drive capable of catching cancellation, redispatch, stale completion, and repetitive-disclaimer regressions.
