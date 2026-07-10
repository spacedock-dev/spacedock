---
title: Align inline workflow state commits with the first-officer durability contract
status: backlog
score: 0.35
source: "Fresh m3 shallow-boot evidence on 2026-07-11: canonical pr-merge startup advancement terminalized and archived an inline entity, then the prescribed spacedock state commit returned a documented no-op, leaving the main checkout dirty. Reproduced twice after stub gh and local Git identity were proven healthy."
completed:
verdict:
worktree:
issue:
id: a03ta7jrnm4termzdjdxzsse
---

The shared first-officer contract says state changes are committed at dispatch and merge boundaries and binds every state write to `spacedock state commit <slug>`. The binary intentionally commits only split-root state and returns a no-op for inline workflows, while the canonical pr-merge startup hook terminalizes/archives inline entities without a raw Git commit instruction. Decide one coherent owner: teach the binary to make a path-scoped inline commit, teach the canonical hook/contract the correct raw commit boundary, or explicitly relax inline durability. Do not hide the mismatch in unrelated live-test assertions.

## Acceptance criteria

**AC-1.** One documented command path durably records an inline entity transition without sweeping sibling changes.

**AC-2.** First-officer contract, canonical pr-merge mod, CLI output, and tests agree on whether `state commit` owns inline workflows.

**AC-3.** A live inline merged-PR advancement ends with terminal archive and the documented Git cleanliness/durability outcome.
