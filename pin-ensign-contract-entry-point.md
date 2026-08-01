---
title: "Dispatched workers load the ensign contract through a resolving entry point, never the first-officer core"
status: backlog
priority: 2
sprint-readiness: ready
source: "Captain, 2026-08-01, after diagnosing the pi ensign misload: every pi-spawned ensign this session (8 workers, both Kimi and gpt-5.6-luna) booted on the first-officer shared core — sometimes from stale .claude/.gemini plugin caches — because ~/.pi/agent/agents/ensign.md declares skills: ['spacedock:ensign'] and the preload silently fails, leaving the model to file-search for its contract. Root-cause question is OPEN: pi-subagents may not route agent-def preloads through pi's package resolver at all (the session was not started by assigning an agent, which motivates verifying rather than inferring)."
id: mxaaqb96syv7pq7ekg5a5194
gates:
    version: 1
    current:
        gate: gate:mxaaqb96syv7pq7ekg5a5194:backlog
    records:
        - id: gate:mxaaqb96syv7pq7ekg5a5194:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mxaaqb96syv7pq7ekg5a5194-backlog-1
              briefing:
                id: briefing:mxaaqb96syv7pq7ekg5a5194:backlog:attempt-1:revision-1
                digest: sha256:d861ca46643dc1cf9100a563ad8ae289a697d67c8232cbb68372d02def8c850a
                digest-domain: canonical-bytes
                request-digest: sha256:c0a9bc4b100816269503de91e29f50c196dafbcc80f4f603ade1e69529724a31
                room-ref: ./pin-ensign-contract-entry-point/review/backlog/briefing-1
---

## Problem

The pi ensign agent definition's contract preload doesn't bind. Workers fall
back to their own filesystem search and land on `first-officer` material, so
the dispatcher's authority leaks into every worker boot. Strikes observed:
8/8 in one session across two models.

## End value

A `.pi`-dispatched worker provably boots on the ensign shared core (the exact
file, from the package) and demonstrably does NOT load anything named
first-officer — verified by reading its session log's first tool calls.

## Investigation directions for ideation (captain-seeded)

1. Search pi's issue tracker/discussions for `agent:` semantics in
   pi-subagents: do agent-definition `skills:` preloads resolve package
   namespaces (`spacedock:ensign`, package pi manifest), and does the loader
   differ for user-scope vs package-scope agent definitions?
2. Read the pi-subagents agent-loading code path for the actual resolution
   semantics instead of inferring them.
3. Worst case accepted by design: drop the custom agent def; dispatch pi's
   general-purpose agent and have the dispatch assignment itself load the
   ensign contract through a resolvable skill name (invoked by name, the way
   `$spacedock:first-officer` works in sessions).
4. The stale `.claude/plugins/cache/**` first-officer copies are part of the
   trap's surface; decide whether cleanup is in scope.

---

### Feedback Cycles

- Cycle 1 (2026-08-01, captain ruling in-session): Direction reset before first stage work — use the generic agent and load the right skill at spawn time ("same as codex"); no shipped agent-definition file, no .pi/agents path (.pi is local), no manifest agents registration. Acceptance stands: a spawned worker provably boots the ensign contract. FO diagnosis on record: basename-only skill resolution (probe-verified), the FO-bootstrap extension injecting into every child session (the active leak vector), and the current pi dispatch-build shape (file-pointer prompt, self-contained assignment, env-level SPACEDOCK_BIN wrapper, subagent_type spacedock:ensign, model=null with opus unsettable on pi).
