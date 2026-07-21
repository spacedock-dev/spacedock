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

## Research & Exploration Directives

- Investigate past changes (specifically PR #480 / commit `0ba08c54`) to understand the original intent behind the minimal `--boot --identify` discovery output.
- Explore options to make the JSON payload self-describing (or clarify the First Officer contract) so LLMs recognize it as a valid, complete terminal boot response.
- Define acceptance criteria and test plan for the proposed ideation stage.


