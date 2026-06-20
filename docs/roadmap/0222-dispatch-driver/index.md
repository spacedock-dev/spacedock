# 0222 — dispatch-driver (0.22.2)

> **Renamed `0206-dispatch-driver` → `0222-dispatch-driver` (captain, 2026-06-19)** — realigns the slug to the version line (it encoded the original 0.20.6 target; now targets **0.22.2** after the version jumped to 0.22.0 at the m4 cut). Member `vp` (dispatch-build-request-file) re-stamped `sprint: 0222-dispatch-driver`; sibling `0221-layered-fo` was `0205`.

> **Seeded 2026-06-17 (captain):** the deterministic dispatch-driver track, descoped from 0221 to keep the layered-FO cut focused. Not yet scope-locked.

**Theme:** push the FO's mechanical dispatch loop into the binary — compute the next action deterministically, and let the FO assemble a worker dispatch from one structured request file instead of hand-built scratch files.

**Membership is the query, never a hard-coded list:** `spacedock status --workflow-dir docs/dev --where sprint=0222-dispatch-driver`.

## Seed members

- **`next-action-driver`** (uncarved) — `spacedock dispatch next-action`: the event-loop computation (reconcile sweep + PR-state check via live `gh pr view` + mod-block scan + dispatchable query + context-budget reuse probe) in ONE deterministic call, returning a fully-qualified `{action, slug?, stage?, team_action?, reason}` where `team_action` is a *resolved* instruction (`send_shutdown` to a named agent; `rebase` of a named worktree with `halt_on_conflict`), NOT a drift-class name the FO interprets. The FO becomes a dispatcher of returned actions. **Riskiest piece:** the return schema IS the FO-facing contract; the w4 spike's mode-2 (drift-class semantic interpretation) was NOT-EXERCISABLE and must be empirically validated here. **Depends on 2y** (`shared-merge-dispatch-contract`, MERGED #385) — unblocked.
- **`vp` `dispatch-build-request-file`** (filed, backlog) — `spacedock dispatch build --request-file`: one Markdown request file (frontmatter scalars + `## Checklist` / `## Scope Notes` / `## Feedback Context`) replacing the FO's 2–3 scratch files. The CONSUMER now; `next-action` becomes the PRODUCER later. **No 2y dependency.**

## The binding constraint

The request-file **schema is defined ONCE and shared** — `vp` consumes it, `next-action` produces it (the prefilled request file the FO lightly edits; the FO is left with the judgment — which signals belong in the checklist). Producer and consumer cannot drift.

**Recommended sequencing:** `vp` first (consumer + schema, buildable now), `next-action` second (producer + the event-loop compute, 2y-unblocked). The **round-trip** is the integration proof: `next-action` emits a request file → `dispatch build --request-file` consumes it → byte-identical to the equivalent flag-form dispatch output. Bundle the *schema*, sequence the *build* — don't fold `vp` into the heavier, riskier `next-action`.

## Why descoped from 0.20.5

0.20.5 (layered-FO) is the **judgment / safety** track — tier-delegation (`72`), gate-extract (`6re`), the state / merge-guard verbs, the prose-function restructure, and the `haiku-drive-validation` gate proof. `next-action` is **mechanical-dispatch mechanization**, and the w4 spike found `«next»` HOLDs (stays prose, Haiku-operable) — so 0.20.5's validation does NOT depend on `next-action` being binarized. Cleanly separable → 0.20.6.

## Definition of Done

To be scope-locked with the captain. (Likely: `vp` + `next-action` shipped with the shared request-file schema, the round-trip integration test green, and the `next-action` return schema validated against the w4 mode-2 risk.)
