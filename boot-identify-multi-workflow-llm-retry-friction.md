---
title: Self-describing boot identify schema and contract hint to eliminate LLM duplicate CLI retry loop
status: ideation
score: 0.85
id: 32vshm0h2h04gs7hzcf315g0
---

## Problem

When `spacedock status --boot --identify --json` is executed at the root of a project containing multiple commissioned workflows (such as this repo with `docs/dev` and a test fixture workflow), the status binary's `resolveIdentifyBootDir` function returns a minimal JSON document:

```json
{"command":"boot","discovery":["/path/to/docs/dev","/path/to/fixture"]}
```

When an LLM agent operating under `$spacedock:first-officer` receives this short payload, it expects a full status summary with entity/stage lists. Because the response lacks explicit completion signals or descriptive status fields, the LLM hallucinates that the output was truncated or mixed with stderr. This triggers an immediate 8+ turn retry loop where the LLM repeatedly runs `spacedock status`, `jq`, `python3` subprocess wrappers, and `go run` before accepting the output. In a recorded Pi session, this duplicate retry loop bloated the context window by thousands of unnecessary tokens before the agent rendered its greeting.

## Context & Past Change Research

Research into repo history reveals that this behavior was intentionally introduced in **PR #480** (commit `0ba08c54`, July 2026: *"Boot identifies, engage converges — collapse the FO Startup recipe to ≤4 prose steps"*):

- `--boot --identify` was designed to be lightweight, side-effect-free, and fast.
- For multi-workflow roots, it returns only the discovered workflow paths (`discovery`) without deep-booting any single workflow or making network/git calls.
- The design intentionally defers workflow selection and convergence to `engage <workflow>`.

The lightweight multi-workflow startup behavior is a core design constraint and must be preserved. The friction is purely an interface legibility / contract communication issue: the JSON payload and FO contract do not make it clear to LLMs that this minimal output is complete and expected.

## Proposed Approach

1. **Self-describing CLI JSON envelope (`internal/status/native_runner.go`):**
   Update the multi-workflow JSON output of `resolveIdentifyBootDir` so it explicitly describes its state:
   ```json
   {
     "command": "boot",
     "status": "multi_workflow_discovered",
     "message": "Multiple workflows discovered. Select one via engage <workflow>.",
     "discovery": [ ... ]
   }
   ```
   Adding `status` and `message` guarantees that an LLM evaluator sees a complete, unambiguous response structure.

2. **First Officer Contract Clarification (`skills/first-officer/references/first-officer-shared-core.md`):**
   Update the `«state.boot»()` / Startup section to explicitly describe the multi-workflow boot payload shape so LLM agents recognize `multi_workflow_discovered` as valid and complete, prohibiting duplicate retry probes before the greet.

3. **CLI Test Coverage (`internal/status/boot_identify_test.go`):**
   Update `TestBootIdentifyUniformZeroOneMany` and related boot identify unit tests to assert the new self-describing fields.

## Acceptance Criteria

- **AC-1 (Self-describing CLI JSON):** Running `spacedock status --boot --identify --json` on a multi-workflow root outputs JSON with `command: "boot"`, `status: "multi_workflow_discovered"`, and the `discovery` array. *Verified by:* unit test in `internal/status/boot_identify_test.go`.
- **AC-2 (Contract note):** `skills/first-officer/references/first-officer-shared-core.md` explicitly documents the multi-workflow boot JSON shape and instructs First Officers to greet immediately without retrying CLI probes. *Verified by:* contract inspection / unit test.
- **AC-3 (Zero retry behavior in LLM boot):** An LLM running `$spacedock:first-officer` on a multi-workflow root executes `spacedock status --boot --identify --json` exactly once and greets immediately without subsequent CLI probes or wrappers. *Verified by:* live/simulated scenario trace.

## Stage Report: ideation

- DONE: Research past changes (`git log -S resolveIdentifyBootDir` / PR #480) to understand why multi-workflow boot returns a minimal discovery list.
  Confirmed PR #480 intentionally collapsed Startup to ≤4 steps and deferred convergence to `engage`.
- DONE: File local task entity under `docs/dev` detailing the LLM retry friction, root cause, past intent, and proposed solution.

