# Ideation gate review — pin-ensign-contract-entry-point (mxaa)

- Stage: ideation (attempt 1)
- Reviewer: agent:first-officer (Pi host)
- Verdict: **approve** — accept AC-1..AC-4 as the baseline and enter implementation in a worktree.

## Root cause (verified, two independent vectors)

1. pi-subagents resolves skills by directory basename only; `spacedock:ensign` is unresolvable in any cwd, and preload failure is silent. Probe-verified live: plain `ensign` resolves, namespaced resolves nothing.
2. `.pi/extensions/spacedock.ts` injects the FO commission into every pi-subagents child (no exemption), re-arming on `session_start`/`session_compact`. The 8/8 FO-boot strikes were courtesy-sensitive: assignment self-sufficiency saved some workers, injection defeated others.

## Design (captain ruling 2026-08-01 implemented: generic agent + load the right skill at spawn)

1. `dispatch build --host pi` emits `agent: "worker"` + `skill: "ensign"` (basename — the only resolvable form); `subagent_type` stays `spacedock:ensign` as role identity/name key; claude/codex build bytes unchanged (golden-guarded).
2. Pi adapter + `internal/piruntime.SubagentDispatch` map the two fields to the subagent call (`context: "fresh"`, `cwd` as today).
3. Extension exempts `PI_SUBAGENT_CHILD === "1"` (pi-subagents-exported marker, source- and live-verified) from FO-bootstrap injection. Edge recorded: an FO nested as a subagents child would lose its bootstrap — accepted, FO commissioning is a root-session concern.
4. `agents/`, `.claude-plugin/`, claude/codex branches untouched — no claude regression surface.
5. Operator-local caches (`~/.pi/agent/agents/ensign.md`, stale `.claude/.gemini` copies): out of shipping scope per captain's removal + redesign making them unreachable on the normal path.

## AC map (mechanisms serve the value)

- **AC-1 (value):** spawned worker's transcript shows `skills/ensign/SKILL.md` within first five reads and zero first-officer reads; live-lane grader over `.pi-subagents/artifacts/*`; baseline 0/8 → graded pass count reported.
- **AC-2:** build JSON fields (default/override/claude-byte-identical).
- **AC-3:** extension harness double-run (marker set/unset), existing assertions unchanged.
- **AC-4:** adapter `«worker.spawn»` binding lines pinned by contractlint; banned namespaced string absent in pi path.

## Necessity and risk

Smallest mechanism: two additive JSON fields + one env check + one contract sentence; alternatives rejected inline. Spike judged unneeded on live evidence (spawn channel ran this session end-to-end; marker exists in the child's env). Staff review judged unnecessary: two small surfaces, twice-captain-challenged design. Declared changes: (a) additive pi build JSON keys; (b) no FO bootstrap in child sessions. All else byte-stable. Live cost: one pi-live dispatch per CI run, extending the existing tagged harness (`TestLivePiSubagentEnsignSmoke`), no new standing lane.

## Fan-out on approval

One implementation worker in a worktree (worktree stage), then fresh validation (stage `fresh: true`), tolerance ±1 for correction rounds; CI live lanes per path→lane (pi-live required; claude byte-identical guard keeps claude lanes routine).
