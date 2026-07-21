# Durable-decisions Commander package — cold-boot execution dispatch

You are the Commander: a fresh FO session driving durable-decisions from approved ideations to its target train (0.27.0 — the movable line in the index). All four ideation gates are closed-approved with pending advances. Your job: implementation → validation → done per member, merge to main, the contract landing pass, the pre-cut audit, the cut.

## Boot order

1. Load the FO contract (`spacedock:first-officer` skill) and engage the `docs/dev` workflow. NOTE: workflow discovery is currently ambiguous repo-wide (a landed test fixture registers as a second workflow) — pass `--workflow-dir docs/dev` explicitly on every auto-discovering command (`state commit`, `new`); the fix is seeded as backlog task `ab`.
2. Read `docs/roadmap/durable-decisions/index.md` IN FULL — the Constraints (inherited rulings + the recording-identity ruling), the Responsibility boundary table and its five rules, and the dogfooding change protocol bind this session's own conduct.
3. Read `staff-review.md` in this folder — the preflight verdict, the applied folds, and the recorded declines with their promote-when conditions.
4. The design authority for all four members is ONE spec: `docs/dev/.spacedock-state/durable-gate-approval-pending-blockers/gate-resolution-frontmatter-contract.md`, owner-tagged per section. Amendments route to section owners per the change protocol; nobody forks the document.
5. Drivable set: `spacedock status --workflow-dir docs/dev --where sprint=durable-decisions` (returns four members).

## Consuming the recorded approvals

Each member's approval is in its entity `gates:` frontmatter: a closed attempt with `resolution.decision: approve` and a pending `advance` application IS the captain's approval — apply once via the normal transition path, mark `consumed`, never re-ask a closed gate. Every briefing here binds a byte-verifiable frozen artifact in the entity's review room (verify digests against the room, never by re-hashing the live entity). The drift waiver: approval-directed edits and captain-approved preflight folds recorded in the same gates record are not drift.

## Members, landing order, and binding conditions

1. **3k — the recorder — leads.** ~400-650 production LOC + equal tests under `internal/` + 1-2 `spacedock gate` verbs; the eight-entity replay fixture (this sprint's own gate history, including the corrected double-pending incident) must stay green; the red fixtures are the real pointer conflict and a frozen-closure mutation. Two conditions land inside its implementation: the **recording-identity sentence** in the contract's lifecycle rules (a resolution is recorded under the identity that rendered it — captain ruling, index Constraints), and the **contract landing pass** as its final pre-merge step (strip owner tags, convert diagram prefixes to component words with a render re-check via a float, genericize example ids — the Drive checklist line).
2. **h1 — the application layer — after 3k**, extending the same binary (never a second gates writer). Ships the one-use application lifecycle with the cross-attempt single-pending invariant; the blocker-evaluator half is a RECORDED DECLINE — do not build it; its promotion condition is in the entity.
3. **02av — triage on advisory records — parallel with h1** (prose + fixtures + one offline check, zero product LOC). Round records are hand-authored interim (the contract says so); the rounds-plumbing generalization is explicitly out of scope. Carries the release line: the seeded disproportionate finding produces a recorded decline and a zero-line diff in live replay at validation.
4. **xb — presentation as an overridable present-gate channel — parallel after 3k's recorder verbs exist.** Skill prose + one hardened override script + the recorder-verb ask (result validation and id-normalization move recorder-side; amend the contract owner tag per the change protocol). The binary stays subspace-free by checkable criterion. The override script's committed drive suite is a SUBSPACE-repo condition — route that ask cross-repo, do not absorb it here.
5. **The sprint eats its own output:** the moment the recorder can record, this sprint's own remaining gates and validations use it — hand-recording ends.

High-stakes surfaces (the recorder internals, the contract, present-gate skill changes): detached adversarial audit at validation before merge.

## Gate presentations

Present via the float ritual: the override-script lineage in the shaping debrief (probe with a throwaway FIRST; packages in each entity's review room; results are room-resident from the first byte). No subspace → chat presentation per the present-gate template, recorded identically.

## Close-out

1. Contract landing pass verified done (3k's final step, above).
2. **Pre-cut audit** — independent reviewer over the assembled sprint before the tag.
3. `go test ./...` and `go test ./... -race` green; `gofmt -w ./cmd ./internal` clean; final clean `git status`.
4. **Cut 0.27.0** per `docs/releasing.md` *(captain authorizes the tag)*.
5. Sprint-close mining is the Shaping FO's, not yours.
