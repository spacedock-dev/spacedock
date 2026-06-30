---
title: "status --ac-scan skips ACs annotated inside the bold (e.g. **AC-1 (VALUE)**), hiding the value AC from the gate cross-check"
status: ideation
score: 0.45
source: "v0.23.0 cut FO session, 2026-06-30. During the 3p validation gate, `status --read <ref> --stage validation --ac-scan` reported only AC-2 and AC-3 — it silently SKIPPED AC-1 because `**AC-1 (VALUE)**` carries the (VALUE) annotation INSIDE the bold markers, which the scan's `**AC-N**` matcher does not catch. The value AC (the most important one) was invisible to the deterministic gate AC cross-check; the FO confirmed AC-1 evidence manually, so 3p was unaffected, but the automated cross-check is weakened for ANY annotated AC."
id: 48gz5715kc4d2j687jbags7v
sprint: 0240-lean-contract
group: tooling
started: 2026-06-30T16:55:24Z
---

`spacedock status --read <ref> --stage <stage> --ac-scan` enumerates `**AC-N**` items and reports each one's evidence/unevidenced status, feeding the gate AC cross-check. Its matcher only recognizes a bare `**AC-N**` token, so an AC whose bold span carries an annotation — `**AC-1 (VALUE)**`, `**AC-2 (no-regression)**`, etc. — is NOT enumerated and is silently dropped from the scan.

## Problem
The README ideation policy explicitly encourages a `(VALUE)`-tagged AC ("At least one AC must MEASURE the end-value"), and the contract's AC cross-check re-anchors on exactly that value AC. So the convention the workflow recommends (`**AC-1 (VALUE)**`) is the convention `--ac-scan` cannot see — the deterministic extraction drops the single most important AC. Live evidence: on `3p`, `--ac-scan` returned only AC-2 and AC-3; AC-1 (VALUE) was absent, even though it was present and evidenced. The gate held only because the FO cross-checked AC-1 by hand.

## Proposed approach
Broaden the `--ac-scan` AC matcher to recognize `**AC-{id}` as a prefix within the bold span (allowing a trailing annotation before the closing `**`), so `**AC-1 (VALUE)**` and `**AC-1**` both enumerate as AC-1. Keep the id capture exact (AC-1, AC-2, ...) and treat anything after it inside the bold as a label.

## Acceptance criteria
- **AC-1 (the value)** — an AC written `**AC-1 (VALUE)**` (annotation inside the bold) is enumerated by `status --ac-scan` with id `AC-1` and its evidence/unevidenced status reported, exactly as a bare `**AC-1**` would be. Verified by: a fixture-driven test feeding an entity body whose AC section uses `**AC-1 (VALUE)**`, asserting the `--ac-scan` output lists `ac=AC-1` — RED on the current matcher, GREEN after the fix.
- **AC-2** — no over-match regression: a bare `**AC-2**` and an inline prose mention like "see AC-3 above" are still handled correctly (AC-2 enumerated once; the prose mention not treated as a new AC item).
