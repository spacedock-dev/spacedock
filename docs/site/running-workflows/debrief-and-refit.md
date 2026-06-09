# Debrief & refit

Two skills keep a workflow durable across sessions and across Spacedock releases: `debrief` captures what a session did so the next one resumes from it, and `refit` upgrades your workflow scaffolding when a new Spacedock release lands.

## Debrief: capture a session

When you end a session, run:

```bash
spacedock claude "/spacedock:debrief"
```

It captures what happened — commits, state changes, decisions, and open issues — into a record the next session picks up. Because the record lives with the workflow, the next session starts from where you left off instead of from a cold transcript.

## Refit: upgrade your scaffolding

When a new Spacedock release is out, run:

```bash
spacedock claude "/spacedock:refit"
```

It upgrades your workflow's scaffolding files to the current version while keeping your local modifications. Use it to bring an existing workflow's structure up to date without losing the project-specific edits you've made.
