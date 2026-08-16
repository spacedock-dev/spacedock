---
title: "Exempt `spacedock new` seed-task filing from the fo-write-core first-load (classify is always allowed-state)"
status: backlog
source: "FO session 2026-07-10 (captain relay from a Codex-runtime FO): the Codex FO cautiously loaded fo-write-core to check write scope before filing a task with `spacedock new`; captain — 'this should be exempted … so we don't need to load the write scope in this case.' Codex FO's own reasoning: 'first-officer-shared-core.md:37 says to load fo-write-core at the first FO-authored file-write intent or state mutation … Calling the operation \"protected\" was the misleading part; after loading, task creation is explicitly allowed.' Verified by the science-officer against this repo's source copy — first-officer-shared-core.md:37, fo-write-core SKILL.md:16,29, internal/status/new.go:24-25,80-91 — the classify outcome is unconditionally allowed-state, so the first-load is pure ceremony for this trigger."
started:
completed:
verdict:
score: 0.25
worktree:
issue:
id: afptght8yjz6s277x8246z7h
---

The FO's boot contract requires loading `fo-write-core` at the first "FO-authored file-write intent or state mutation" (`skills/first-officer/references/first-officer-shared-core.md`, `## Deferred load points`, line 37). Filing with `spacedock new` matches that trigger (writes an entity file, mutates workflow state), so the FO correctly loads the skill and runs `«write.classify»(target, intent)`. But the classify is a guaranteed yes: `runNew` writes only to the discovered workflow's state dir for both flat and folder forms (`internal/status/new.go:24-25,80-91`); the only flags are `--folder` and `--workflow-dir` (both keep the write inside a workflow state dir); and the classifier's `allowed-state` row explicitly lists `spacedock new` for `.spacedock-state/**` (`skills/fo-write-core/SKILL.md:16`). No `new` invocation lands in a `blocked-product` path. The mandatory first-load never changes the outcome — pure latency, and it invites mislabeling an always-allowed filing as "protected."

Fix (single site, mirroring the already-carved engage-sweep exemption in the same section): extend line 37's exemption clause to add `spacedock new` seed-task filing —

> NOT «engage»'s sweep / pr-merge advancement, whose `status --set`/`archive` are pre-authorized at engage and need no write-scope load; and NOT `spacedock new` seed-task filing, whose atomic-create write always lands in the workflow state dir (`allowed-state`), so the classify is a guaranteed yes and it needs no write-scope load either.

Consistent with `fo-write-core`'s FO Write Scope, which already lists `spacedock new` as an allowed write (`skills/fo-write-core/SKILL.md:29`); the carve-out only says the FO need not LOAD the skill to reach that yes. No change to `fo-write-core` itself — the engage-sweep precedent lives only in Deferred load points.

Scope: `skills/**` is `blocked-product`, so this doc edit goes through a dispatched worker in a worktree, not a direct FO edit. Touches the same `## Deferred load points` section as `dpwp415…` (fo-deferred-load-point-hunt, at validation, PR #491) — independent change; rebase on that merge to avoid a textual conflict. Doc-only, low blast radius, no detached adversarial audit owed; validation just needs the reworded clause to parse and read cleanly.

## FO note (2026-08-16, from fo-workflow-fit-gate ideation)

This entity's rationale ("the classify outcome is unconditionally allowed-state, so the
first-load is pure ceremony for the `new` trigger") is invalidated if the Workflow Fit
Gate amendment lands: after it, fo-write-core.md carries the one question the classifier
cannot answer, and `spacedock new` becomes the single most important trigger to load it
at. Re-evaluate against the amended write core before any dispatch; shipping this as-is
would silently remove the fit gate's read trigger.
