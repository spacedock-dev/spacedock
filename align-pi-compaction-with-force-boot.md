---
title: "Align the Pi extension's compaction hook with force-boot (re-read state, not re-inject contract)"
status: backlog
source: "Captain (2026-08-21): the Pi extension's session_compact hook re-injects FO_BOOTSTRAP_TEXT (a contract pointer), but PR #738 (force-boot-at-compaction-boundary, merged) established the opposite mechanism — re-read durable state via one «state.boot»(), do NOT re-inject the contract. The Pi extension did not follow #738."
started:
completed:
verdict:
score: 0.8
worktree:
issue:
id: h9nn5brc1dp0m82x5en21d56
---

The Pi extension's `session_compact` hook (`.pi/extensions/spacedock.ts`) re-injects `FO_BOOTSTRAP_TEXT` at the compaction boundary — a contract pointer telling the FO to re-satisfy load preconditions and re-read durable state. PR #738 (`force-boot-at-compaction-boundary`, merged) established the opposite for Claude/Codex: fire one `«state.boot»()` (re-read durable state); the contract does not need re-injecting. The Pi extension is misaligned — it does the thing #738 rejected (re-inject the contract) and only points at re-reading state rather than doing it.

Direction: change the Pi `session_compact` hook to fire a `«state.boot»()` (the `spacedock status --boot --identify --json` read) at the compaction boundary, aligning with #738's "re-read durable state, don't re-inject contract." The boot record carries mods/ready_gates/dispatchable/pr_state/team_state/state_backend — the same answers #738 provides Claude/Codex. Keep the bootstrap text re-injection only if a follow-on load actually requires the contract pointer (it likely does not — #738's conclusion is the contract survives compaction; only the state goes stale).

Acceptance sketch: value — a compacted Pi FO re-reads durable state on resume (the boot record), matching Claude/Codex; the bootstrap-text re-injection is removed or reduced to what #738's mechanism doesn't cover. mechanism — a behavior test asserting the compaction boundary triggers a boot read, not a contract re-injection. Expected surface: `.pi/extensions/spacedock.ts` + test; small. Stacks on PR 753 (same Pi-extension file).
