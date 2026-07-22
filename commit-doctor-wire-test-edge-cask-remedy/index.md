---
title: Commit a doctor wire-test for the @next edge-cask remedy
status: backlog
source: "vz validation (#556, upgrade-remedy @next trim) — the value fix is live-proven but the gateHost->remedy wire (runningEdgeCask() through gateHost/init.go) has no committed test; a future edit severing it would not red (validator Mutation C stayed green)."
id: c0gaq7ccwh2t6bg267q7rarr
---

Add a committed doctor test that stages a binary under a fake `Caskroom/spacedock@next/` path and asserts the emitted too-old-binary remedy names `brew upgrade spacedock@next`, so the gateHost->remedy wire is guarded (severing it reds a committed test). Promotion condition named by the #556 validator.
