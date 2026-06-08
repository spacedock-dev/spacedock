# Sprint 0198 — staff readiness review

Same mechanism as `019x` (see [`../019x-pre-flip-cleanups/staff-review.md`](../019x-pre-flip-cleanups/staff-review.md)
for the full description): before the Commander drives, an INDEPENDENT reviewer does a
readiness gap-analysis over the ideated sprint — producibility, build-readiness,
grounded-risk, coherence, integration-test reality — output as BLOCKING / NON-BLOCKING
plus a go/no-go.

## Findings — run 1

> Status: **pending.** Run after the ideation wave settles (`qa`/`nzb7`/`69rk`/`jh` ideated;
> `4t` ideation pending; `kb`/`1p` fast-track; captain still filing install-path members).
> The review must weigh the items below.

Going in (to confirm or refute):

- **`jh` recommends a `CONTRACT_VERSION` bump (1→2) — a captain-level architectural fork, not
  a routine cleanup.** Its spike confirmed the skew and found the flag/file form shipped
  (`cfa3b671`) without the contract bump its own doc comment mandated. Bumping now corrects
  that — but it makes 0.19.8 a **breaking** contract change (a 0.19.7 contract-1 binary is
  rejected by a 0.19.8 plugin requiring `>=2,<3`), which contradicts the flip's
  "no contract change" stance and couples directly to `qa` (the incompatible-binary UX would
  then be exercised by real 0.19.7→0.19.8 upgraders). **The captain must rule on the contract
  direction before `jh` is driven**; the review judges whether the bump is proportionate vs.
  `jh`'s rejected stdin-fallback.
- **Coherence (survey group):** `69rk` found agentsview keys `project` by GIT-ROOT basename,
  not cwd basename — diverging from the current fixture/SKILL.md model. This touches `1p27`
  (scaffold fact) and `4t` (agentsview detection). Confirm the three survey members converge
  on one corrected agentsview model, not three divergent ones.
- **`qa` ↔ `jh` coupling:** both binary-ux; if `jh` bumps the contract, `qa`'s
  incompatible/stale messages must reflect the new versions. Sequence/align them.
- **Fast-track check:** `kb`/`1p` skip ideation (well-specified) — confirm build-ready from
  their bodies, or pull them back to ideation.
