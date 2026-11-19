# Roadmap

The Spacedock roadmap — strategy docs, sprint folders, staff reviews, and durable decisions — lives with the workflow state, not on `main`.

- **Branch:** `spacedock-state/dev`, folder `_roadmap/`.
- **Local path:** `docs/dev/.spacedock-state/_roadmap/`, after `spacedock state init --workflow-dir docs/dev` creates the state checkout.

Roadmap entries are planning state tied to tasks: sprint membership is a query over the state checkout (`spacedock status --workflow-dir docs/dev --where sprint=<slug>`), so the roadmap is versioned beside the entities it plans.

Links of the form `docs/roadmap/...` in task bodies and older docs refer to the same paths under `_roadmap/`.
