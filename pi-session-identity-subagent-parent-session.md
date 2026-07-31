---
id: ekw79nn8z9829d77dw7y9353
title: "Pi session identity via PI_SUBAGENT_PARENT_SESSION — runtimehost identity column + install-gate sentinel key"
status: backlog
source: "Live env evidence, 2026-07-31: both the captain's shell and the FO's own root-session tool shell carried PI_SUBAGENT_PARENT_SESSION equal to the running pi session's own id (019fb5d1-85af-73f6-bb07-20bfc04004db). The runtimehost marker table (internal/runtimehost/runtimehost.go:23-24) claims pi exposes no identity env var — the code is stale about pi's actual env surface."
started:
completed:
verdict:
score:
worktree:
issue:
---

The runtimehost marker table marks both pi rows with `identity: ""`, so every pi consumer of session identity (starting with the FO install-gate's one-attempt sentinel in `fo-install-gate.md`, D-3) must fall back to a project/cwd-hash scope. Live evidence from 2026-07-31 shows pi DOES inject a session-scoped identity variable: `PI_SUBAGENT_PARENT_SESSION`, set in the session's own tool shells (not just nested subagent shells) to the session's own id. Mapping it as pi's identity lets the install-gate sentinel be truly session-scoped on pi, removing the admitted over-reach where a failed install in session N suppresses the install offer in session N+1 of the same project until tmp cleanup. Follow-up to `fo-boot-install-hint-linux-direct-sandbox` (D-3).

## Problem

{Ideation fills this in, seeded: pi shells carry a session-scoped id variable that runtimehost's identity column doesn't read, forcing every consumer to a coarser scope (cwd hash / project). The consequence in the install-gate sentinel is a documented false-negative window across sessions on pi.}

## Proposed approach

{Ideation fills this in, seeded: add `{host: "pi", marker: "PI_SUBAGENT_PARENT_SESSION", identity: "PI_SUBAGENT_PARENT_SESSION"}`-style mapping to runtimehost's marker table (careful: marker vs identity semantics — the table's comments show identity is read ONLY for the first-set host row, so the pi identity mapping needs the row ordering handling, or a dedicated identity-separate lookup); then have the install-gate sentinel key use `runtimehost` identity with cwd-hash as the fallback when unset. Verify a ROOT pi session sets it identically (this session's own tool shell is one data point); check whether pi sets it in every context (bare terminal pi vs TUI vs subagent spawned).}

## Out of scope

{Ideation fills this in. Seeded exclusions: changing pi itself; changing claude/codex identity behavior; re-deriving any decision from fo-boot-install-hint-linux-direct-sandbox beyond the sentinel key.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - runtimehost detects pi session identity: when `PI_SUBAGENT_PARENT_SESSION` is set, the identity result for the pi host is its value; when unset, behavior is unchanged (identity omitted as today).**
Verified by: {ideation names the test — Go unit test driving Detect with a getenv map; independent baseline: today's identity column is a hard-coded empty, so the test fails if the mapping is absent.}

**AC-2 - The install-gate sentinel key on pi is session-scoped when an identity is available: a failed install in session N does NOT suppress the install offer in session N+1 of the same project (the cwd-hash over-reach is gone when PI_SUBAGENT_PARENT_SESSION is present), and the cwd-hash fallback still protects sessions where the var is absent.**
Verified by: {ideation names the test — behavior fixture simulating two sessions with different PI_SUBAGENT_PARENT_SESSION values sharing a project dir, asserting the second session's offer is NOT suppressed.}

**AC-3 - The claim that a root pi session sets PI_SUBAGENT_PARENT_SESSION in tool shells is verified against the live env rather than asserted.**
Verified by: {ideation names the verification — a live evidence capture recorded in the entity (e.g. `env | grep` in a root-session tool shell), and an honest "unknown for pi versions older than X" note if the matrix can't be established.}

## Test plan

{Ideation fills this in: runtimehost unit tests for the pi identity mapping rows (order/ambiguity unchanged for the two existing pi markers; CODEX/CLAUDE paths untouched), behavior fixture for AC-2's two-session simulation, and the env-surface dogfood for AC-3.}
