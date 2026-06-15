# FO layered architecture + model-aware delegation (later goal)

Strategy note, not yet a sprint. Captured 2026-06-15 (captain). Files a direction for the
first-officer system to revisit once the nearer read-cost and contract-structure work lands.

## The idea

Split the first-officer system into three layers with clean seams:

1. **Automation** — the mechanical, deterministic substrate: the `spacedock` binary
   (`status` / `dispatch` / `new` / `reconcile`), git plumbing, the scripted state
   transitions and guards. No judgment; pure mechanism. Already mostly real.
2. **Driving and dispatch** — an agent runs the mechanical flow by the book: reads state,
   builds checklists, dispatches workers, advances stages, runs gates and merges per the
   contract. This is today's FO loop executing the contract deterministically.
3. **High function** — judgment that can **override** the mechanical automation: design
   discussion, architectural decisions, gate adjudication beyond the default, scoping, the
   calls that need reasoning rather than rule-following.

## Self-aware operating level

An FO session should be **self-aware of which level it is fit to operate at — 2 or 3.**

- A capable model operates at level 3: it drives (level 2) *and* exercises high-function
  judgment, overriding automation when warranted.
- A weaker model self-identifies as **level-2-only**: it runs the mechanical driving and
  files work, but does **not** attempt level-3 judgment. Instead, **all discussion beyond
  filing is routed to a separate team member running on a higher model.**

The routing rule: a weak-model FO is the level-2 hands (mechanical dispatch + filing) and
delegates every judgment or discussion call to a stronger teammate, the level-3 brain.

## Why it matters

- **Capability/cost matching** — a cheap model can run the deterministic loop; the
  expensive model is spent only on the calls that actually need reasoning.
- **Safety by construction** — a weak model that knows it is weak does not make
  architectural or gate calls it should not; it escalates structurally rather than by luck.
- **Composability** — the three layers mirror the seams the system already has (mechanical
  binary vs FO contract vs captain judgment), so the split is a sharpening of the existing
  structure, not a new paradigm.

## Open questions (for ideation when this becomes a sprint)

- **Self-identification** — how does an FO know its level? Declared model tier, a self-probe,
  or a flag the launcher sets from the resolved session model?
- **Hand-off protocol** — what exactly routes to the higher-model teammate (gate
  adjudication? design discussion? scope decisions?) versus stays mechanical (filing, state
  transitions, dispatch)? The boundary is the whole game.
- **The level-3 teammate** — a standing member (like `comm-officer`), spawned on demand for
  judgment calls, or a different shape?
- **Relationship to existing roles** — how this composes with the shaping-FO / Commander
  split and the captain. The captain is already the ultimate level-3 override; this would
  insert a model-aware level-3 *agent* between the mechanical FO and the human.
