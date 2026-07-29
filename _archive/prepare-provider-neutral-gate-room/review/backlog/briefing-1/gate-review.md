# BACKLOG GATE — Provider-neutral gate-room preparation (`s4`)

Recommendation: **APPROVE and dispatch ideation.**

## Capability and value

Give First Officers one mechanical way to turn a concise gate review and selected existing artifacts into a complete recorder-ready room. This removes hand-authored request ids, locators, digests, and basename assumptions while keeping Spacedock provider-neutral.

## Binding boundaries

- Align directly with Subspace `em` after it lands; v1 needs no compatibility layer.
- Bind the canonical Briefing by locator, id, and digest, with any readable filename.
- Derive the durable association from request, Briefing, Result, and inventory; do not create `association.json`.
- Reject duplicate JSON members independently at every authority-bearing recorder boundary.
- Keep Subspace transport and `/subspace:r gate <room>` in Subspace q0.
- Leave broader lifecycle guidance, advisory rounds, and readiness projection in `gate-agent-ergonomics`.

## Proof direction

Ideation must first reproduce xb’s basename failure with an otherwise valid arbitrary-name Briefing. The final behavior needs command-level preparation and recording fixtures, byte-clean adversarial duplicate-member failures, the existing skill integration exercise, full Go/race gates, and a detached high-stakes audit.

## Decision ask

Approve this narrow release dependency for ideation, or revise/hold with a concrete boundary.
