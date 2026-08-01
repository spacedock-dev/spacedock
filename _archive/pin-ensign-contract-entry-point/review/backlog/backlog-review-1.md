# Backlog gate review — pin-ensign-contract-entry-point (mxaa)

- Entity: `pin-ensign-contract-entry-point` (`mxaaqb96syv7pq7ekg5a5194`), priority 2, sprint-readiness: ready
- Stage: backlog (first formal gate binding)
- Reviewer: agent:first-officer (Pi host), 2026-08-01
- Verdict: **approve** — accept into ideation with the captain's packaging rulings as direction.

## Root cause (verified this session, code + live probe)

1. pi-subagents names discovered skills by directory basename only (`collectFilesystemSkills`); the namespaced string `spacedock:ensign` is unresolvable in any cwd, and `resolveSkillsWithFallback` continues silently. Probe child at repo root: plain `ensign` resolved to `skills/ensign/SKILL.md` (via the repo `package.json pi.skills` project-package path); `spacedock:ensign` resolved nothing.
2. `.pi/extensions/spacedock.ts` injects the FO-commission bootstrap on every `session_start`/`session_compact` with no worker exemption; the probe child received the injection. This is the active vector steering contract-less workers onto the first-officer core.
3. The repo already ships `agents/ensign.md` + `agents/first-officer.md` with `skills: ["spacedock:ensign"]` (valid for Claude's plugin idiom, broken on pi), and the pi manifest registers `extensions` + `skills` but no `agents` — hence the hand-created user def the captain removed.
4. Pi extension API `resources_discover` (docs) cannot register agents (skill/prompt/theme paths only); agent defs reach pi-subagents only via user dir, project `.pi/agents`/`.agents`, or package manifests `pi.subagents.agents` / `pi-subagents.agents` (agents.ts:339-360).

## Captain rulings (2026-08-01, in-session)

1. No `.pi/` project-local agent dir — that dir is local, not shipped.
2. Ship the agent def inside the extension bundle (package manifest route).
3. Acceptance: **skills loaded correctly** — a spawned worker provably boots on the ensign contract (session-log first tool calls show the ensign SKILL.md read and no first-officer reads).

## Structural choices sent to ideation

1. Per-host agent-def split vs shared def: the shipped `agents/ensign.md` declares `skills: ["spacedock:ensign"]`, valid on Claude (namespace real), broken on pi (basename-only). Options: one plain-basename def verified against BOTH hosts' discovery, or a pi-specific shipped def dir (`agents/pi/`) referenced by the manifest's `pi.subagents.agents`. Ideation verifies the Claude-side semantics of the plain-basename form before choosing; no claude regression.
2. Manifest key: `pi: { subagents: { agents: [...] } }` vs top-level `"pi-subagents".agents` — ideation pins against the loader code.
3. Boot-sequence fallback line in the def — same namespace problem on pi; per-host text or neutral wording.
4. Extension exemption: WITHOUT it, a correctly-loaded ensign still receives the "You are the first officer" injection and the two contracts race — the captain's acceptance test (first tool calls ensign-only) cannot pass reliably. Included in this entity's scope for captain ratification at the ideation gate rather than silently split out.

## Notes

- The deleted user def (`~/.pi/agent/agents/ensign.md`) will not be recreated; this entity ships the proper home.
- Full dispatch on pi remains degraded (adhoc workers + assignment-carried contracts) until this lands.
