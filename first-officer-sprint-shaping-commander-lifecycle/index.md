---
title: First-officer-owned sprint shaping and Commander lifecycle
status: backlog
score: 0.85
source: "Captain design discussion 2026-07-11: make the First Officer the only user-facing surface and generalize sprint shaping/drive across project types."
completed:
verdict:
worktree:
issue:
id: ppw2wn21s9qb97z2yz697qt4
---

## Problem

Users should state a goal and retained authority to the First Officer, not invoke or coordinate shaping and Commander roles themselves. The FO needs a reusable way to decide whether to drive directly or shape one or more sprint-level deliverables, persist each approved drive, dispatch one Commander per sprint, and remain the captain-facing supervisor.

The design must work across different project types without hardcoding development stages, release mechanics, design reviews, publishing flows, or private workflow details.

## Current direction

- The user invokes only `spacedock:first-officer` with a goal, observable completion condition, and protected actions.
- Sprint shaping is an internal FO capability, not a user-invocable skill.
- The FO owns synthesis: outcome, scope, cohort, dependencies, proof, gates, and the durable drive package.
- One Commander owns one integrated sprint deliverable. Different workflows or workstreams do not justify separate Commanders by themselves.
- The FO drives small goals inline and delegates only when fan-out, isolation, or an independent sprint boundary justifies a Commander.
- Commander is an internal worker skill in the existing Spacedock plugin. The FO dispatches it from an approved, immutable package reference.
- The FO remains the sole user-facing status and gate surface while Commanders drive their packages and dispatch workers.
- The FO chooses the sprint's integration topology during shaping. A sprint may merge members independently to the target branch, or assemble them on one shared sprint branch and present the integrated result as one final PR.

## Shared sprint-branch pattern to consider

Some sprint deliverables are meaningful only as an assembled whole. For those drives, the approved package may declare a shared integration branch:

1. The FO pins the target branch, immutable starting revision, sprint branch name, integration proof, and final review gate while shaping.
2. The Commander creates or verifies the sprint branch from that revision. Each member still executes and validates through its own entity and isolated worktree.
3. After a member passes its entity gate, the Commander merges it locally with `--no-ff` into the sprint branch. The sprint branch is pushed after each accepted member so shared CI, preview, or integration evidence can run against the assembled result.
4. Sprint-wide review and integration proof run on the shared branch, not on disconnected member branches.
5. After ship-blockers are resolved, the Commander opens one final sprint PR from the shared branch to the target branch. That PR is the review and publication unit for the integrated deliverable.

Member branches do not need separate PRs merely to enter the sprint branch. A member-level PR remains available when that member has an independent external review, compliance, ownership, or publication gate. Otherwise, PR-per-member adds a second review ceremony without changing the sprint's actual delivery boundary.

This topology is optional. Independently shippable members may merge directly through their normal entity policy. The FO selects the topology from the deliverable's integration and review needs; the Commander must not change it mid-drive without returning a sprint-wide amendment to the FO.

The design must also keep workflow state independent of the sprint branch. Entity state and Commander obligations remain in their authoritative state checkout; changing the code integration branch must not redirect or strand state commits.

## Artifact lifecycle to design

The sprint artifact is a durable contract and audit record, not a second task tracker. Entity frontmatter and stage reports remain authoritative for executable work.

The ideation should define:

1. **Sprint contract:** goal, definition of delivered, scope, membership queries, workflow roles, dependency order, integration proof, retained human gates, terminal action, and escalation conditions.
2. **Membership stamps:** a minimal convention such as `sprint`, with optional grouping/readiness fields, without duplicating a roster in prose.
3. **Approval stamp:** how the FO records the approved package revision and captain-granted conn without a self-referential hash or ambiguous mutable path.
4. **Dispatch pointer:** how a Commander receives an immutable `{repository revision, package path, drive epoch}` reference and proves it loaded that exact contract.
5. **Drive lifecycle:** a small stage-neutral state model such as draft, approved, driving, reviewing, and done, while avoiding a second workflow engine.
6. **Amendments:** entity-local execution detail stays in the entity; sprint-wide scope, order, proof, gate, or terminal changes return to the FO, append an explicit amendment, advance the package revision, and require the appropriate captain approval before the Commander continues.
7. **Resume:** how a new FO reconstructs active Commander obligations from durable artifacts and live workflow state without relying on a conversation summary or worker handle alone.
8. **Close:** how integration review, deferred findings, debrief, and next-sprint seeds attach without turning the package into a historical dump.
9. **Integration topology:** how the approved package selects direct member merges or a shared sprint branch; pins its base and target; defines member-merge, push, preview, integration-review, and final-PR boundaries; and prevents workflow state from riding the integration branch.

## Privacy boundary

Generalize only the reusable operating pattern. Shipped instructions, examples, fixtures, tests, and reports must not name, quote, link to, or reproduce private projects, workflows, people, campaigns, content, paths, or business metrics. Use synthetic project profiles and invented artifacts for validation.

## Questions for ideation

- Should the sprint contract and Commander dispatch package be one versioned document or two linked artifacts?
- Which stamps remain convention-only, and which need binary validation or mutation guards?
- Does package approval use a Git commit, an explicit revision field, a content digest outside the document, or a binary-minted drive identifier?
- How does the FO amend an active drive without silently invalidating a Commander's loaded contract?
- What durable evidence distinguishes a proposed Commander from an active, completed, superseded, or abandoned drive?
- How are profile-specific proof and outward-action gates supplied without leaking them into the stage-neutral core?
- When should one user goal become multiple sprints and therefore multiple Commanders?
- When does an assembled deliverable justify a shared sprint branch and one final PR, and when does a member still require its own PR before integration?
- Which actor may create, update, and close the sprint branch, and what evidence lets a resumed FO prove that its head contains exactly the accepted member revisions?

## Acceptance criteria for ideation

- Define a user journey in which the user interacts only with the FO from goal through completion.
- Define the FO's direct-drive versus Commander-dispatch decision with explicit smallest-sufficient conditions.
- Define one-Commander-per-integrated-deliverable semantics, including cross-workflow drives.
- Define both integration topologies: independently merged members and locally assembled shared sprint branches with one final PR; specify the selection rule and the exception for member-level external gates.
- Specify the artifact schema, authoritative-state boundaries, approval stamp, immutable dispatch pointer, amendment protocol, resume behavior, and close behavior.
- Exercise the design against at least three synthetic profiles with different proof, gate, and terminal-action shapes.
- Determine whether shaping and Commander remain internal skills/references in the existing plugin and identify any binary surface justified by enforcement needs.
- Provide privacy-safe forward-test scenarios that contain no private project material.
