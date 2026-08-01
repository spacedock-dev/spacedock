---
title: "merge guard + terminal --set over-refuse entities holding consumed nonterminal-target approvals (fail-closed classification exceeds spec)"
status: backlog
source: "FO session 2026-08-01 during PR #590 CI: validator transcript forensics plus local differential replay; captain follow-up directive"
id: fvzjk6w59gc6syy5kk0rta69
---

During the 1w62 CI ceremony, live runs on the codex AND claude-opus lanes (ungraded fixture approved-gate) plus an independent validator surfaced that PR #590's fail-closed classification refuses to terminalize any entity whose gates record is readable but carries no PENDING terminal-target approval — e.g. a consumed nonterminal-target approval from a gated earlier stage followed by an ungated terminal advance.

Replay-verified locally (minimal single-root fixture, gated backlog approved+consumed to implementation, ungated validation, then terminal ceremony): main@8d978b638 `--set status=done` and `merge guard --verdict passed` both exit 0; candidate@b54b79043 refuses BOTH exit 1 ("entity carries no binding pending terminal-target approval (condition \"consumed\")"), byte-clean, --force works.

Mechanism: internal/status/merge.go:553 pendingTerminalApproval flattens only gates.ErrNoGateRecord; any readable record short of pending-terminal-target errors, and both finalize (merge.go:427/493) and terminal --set (handlers.go:259-275) key off it. The diff's own spec (docs/specs/gate-resolution-frontmatter-contract.md:301) scopes the refusal to "while a pending terminal-target application is in force" — this shape carries none, so the refusal exceeds the spec. No test pins the shape (TestTerminalDeliveryRefusalsByteClean covers gate-less/superseded/digest-stale only).

Scope: durable-decisions unaffected (validation is gated; terminal transitions always carry pending authority). Affected: workflows mixing gated nonterminal stages with ungated terminal advances; hand --set terminal on gated-record entities.

Suggested direction (from the surface, not prescribed): pendingTerminalApproval should return legacy/no-authority (not error) for readable-but-not-pending-terminal records, keeping the sole-consumer refusal strictly for pending-terminal-target approvals; add one fixture test pinning the mixed-gating shape through both merge guard and terminal --set. Evidence: run 30684502027 (codex items 49/51, claude-opus lane); validator artifacts under /tmp/ze-validation-artifacts.

## Acceptance criteria

- **AC-1** An entity with a consumed nonterminal-target approval reaching an ungated terminal stage finalizes via merge guard and terminal --set exactly as a gate-less entity.
- **AC-2** A pending terminal-target approval still blocks terminal --set and routes through merge guard as sole consumer (existing refusal tests stay green).
- **AC-3** One new fixture/test pins the mixed-gating shape; full suite green.
