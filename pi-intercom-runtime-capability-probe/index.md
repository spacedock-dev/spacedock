---
title: Pi intercom runtime capability probe
status: ideation
source: captain (2026-06-04) — Codex idle-notification evidence pattern may generalize to prove Pi intercom contact_supervisor runtime capability
score: "0.31"
started: 2026-06-04T00:00:00Z
completed:
verdict:
worktree:
issue:
id: cq9kb7cdpp9y48tn8gwzmqzq
---

# Pi intercom runtime capability probe

Create a generalized runtime capability probe pattern, using Pi intercom supervisor talkback as the first target capability.

## Problem

`subagents-doctor` can report the Pi intercom bridge as active, but that only proves package discovery/setup. It does not prove the runtime capability that matters to Spacedock operations: a child subagent launched after bridge activation can call back to the parent, send non-blocking progress, block on a decision, receive a parent reply, resume, and write durable evidence after the reply.

Codex recently introduced a useful pattern for runtime behavior evidence: a probe recipe plus durable evidence JSON validated by integration tests for notification behavior. Pi intercom needs the same style of proof, generalized as a reusable runtime capability probe pattern rather than a one-off transcript claim.

## Proposed approach

Ideation should adapt the Codex idle-notification evidence pattern into a runtime capability probe surface:

- a probe recipe under `docs/dev/` with the exact live smoke prompt and operator steps;
- evidence JSON under `docs/dev/_evidence/` with a small schema and explicit classifications;
- integration tests that validate the recipe shape, evidence schema, and runtime-contract wording;
- Pi-specific live smoke instructions that prove `contact_supervisor` progress and decision round trips with a durable marker after resume.

The Pi intercom probe should distinguish setup evidence from capability evidence. `subagents-doctor` bridge-active output is necessary but not sufficient.

## Out of scope

- Reworking `pi-subagents` or `pi-intercom` internals.
- Replacing the existing Pi runtime/frontdoor live smokes.
- Generalizing every host runtime in the first implementation; Codex notification evidence can be the reference pattern, while Pi intercom is the first new capability probe.

## Acceptance criteria

**AC-1 - Runtime capability probes have a documented reusable evidence shape.**
Verified by: integration tests over a probe recipe/evidence schema that fail if required fields, classifications, or interpretation text are missing.

**AC-2 - Pi intercom supervisor talkback is represented as a concrete capability probe.**
Verified by: a recipe containing an exact child prompt for `contact_supervisor` progress_update, `need_decision`, parent reply `APPROVED`, and durable marker `PI-INTERCOM-SMOKE-APPROVED`.

**AC-3 - Evidence distinguishes setup from real talkback capability.**
Verified by: evidence-schema tests requiring both bridge/setup fields (doctor bridge active, pi-intercom path) and behavioral fields (child received tool, progress observed, decision resumed, durable marker path/content).

**AC-4 - Runtime docs or probe notes do not claim doctor bridge-active alone proves talkback.**
Verified by: instruction/probe text invariant tests that require wording such as “bridge active is necessary but not sufficient” and forbid wording that equates doctor success with supervisor talkback proof.

**AC-5 - A live Pi intercom smoke can be recorded as durable evidence.**
Verified by: a live/manual probe run writes evidence JSON with classification `passed` only when the parent observed progress, replied to the decision, and the child wrote the durable post-reply marker.

## Test plan

- Add integration tests analogous to the Codex idle-notification evidence tests for recipe shape, allowed classifications, and JSON evidence schema.
- Add or update runtime/probe documentation tests for the doctor-vs-capability distinction.
- Run `go test ./skills/integration -count=1` and `go test ./... -count=1`.
- Run the live Pi intercom probe only when Pi auth, `pi-subagents`, and `pi-intercom` are available; otherwise record exact missing setup without marking AC-5 passed.
