# Durable-decisions Commander package — cold-boot execution dispatch

You are the Commander: a fresh FO session driving durable-decisions from approved ideations to its target train (0.27.0 — the movable line in the index). All four ideation gates are closed-approved with pending advances. Your job: implementation → validation → done per member, merge to main, the contract landing pass, the pre-cut audit, the cut.

## Boot order

1. Load the FO contract (`spacedock:first-officer` skill) and engage the `docs/dev` workflow. NOTE: workflow discovery is currently ambiguous repo-wide (a landed test fixture registers as a second workflow) — pass `--workflow-dir docs/dev` explicitly on every auto-discovering command (`state commit`, `new`); the fix is seeded as backlog task `ab`.
2. Read `docs/roadmap/durable-decisions/index.md` IN FULL — the Constraints (inherited rulings + the recording-identity ruling), the Responsibility boundary table and its five rules, and the dogfooding change protocol bind this session's own conduct.
3. Read BOTH staff-review seats in this folder — `staff-review.md` (fable) and `staff-review-codex.md` (codex, NOT READY at the time) — plus the index Constraints recording the codex-seat rulings (authorization-only consumption, the named subspace release condition, digest domains, recording identity). Every codex material finding is closed in durable state: room evidence tracked, adoption provenance recorded, the contract reconciled at sha256 681b2348... via 3k gate attempt 9 (the first FO-identity closure under the recording-identity ruling).
4. The design authority for all four members is ONE spec: `docs/dev/.spacedock-state/durable-gate-approval-pending-blockers/gate-resolution-frontmatter-contract.md`, owner-tagged per section. Amendments route to section owners per the change protocol; nobody forks the document.
5. Drivable set: `spacedock status --workflow-dir docs/dev --where sprint=durable-decisions` (returns four members).

## Consuming the recorded approvals

Each member's approval is in its entity `gates:` frontmatter: a closed attempt with `resolution.decision: approve` and a pending `advance` application IS the captain's approval — apply once via the normal transition path, mark `consumed`, never re-ask a closed gate. Every briefing here binds a byte-verifiable frozen artifact in the entity's review room (verify digests against the room, never by re-hashing the live entity). The drift waiver: approval-directed edits and captain-approved preflight folds recorded in the same gates record are not drift.

## Members, landing order, and binding conditions

1. **3k — the recorder — leads.** ~400-650 production LOC + equal tests under `internal/` + 1-2 `spacedock gate` verbs; the eight-entity replay fixture (this sprint's own gate history, including the corrected double-pending incident) must stay green; the red fixtures are the real pointer conflict and a frozen-closure mutation. Two conditions land inside its implementation: the **recording-identity sentence** in the contract's lifecycle rules (a resolution is recorded under the identity that rendered it — captain ruling, index Constraints), and the **contract landing pass** as its final pre-merge step (strip owner tags, convert diagram prefixes to component words with a render re-check via a float, genericize example ids — the Drive checklist line).
2. **h1 — the application layer — after 3k**, extending the same binary (never a second gates writer). Ships the one-use application lifecycle with the cross-attempt single-pending invariant; the blocker-evaluator half is a RECORDED DECLINE — do not build it; its promotion condition is in the entity.
3. **02av — triage on advisory records — parallel with h1** (prose + fixtures + one offline check, zero product LOC). Round records are hand-authored interim (the contract says so); the rounds-plumbing generalization is explicitly out of scope. Carries the release line: the seeded disproportionate finding produces a recorded decline and a zero-line diff in live replay at validation.
4. **xb — presentation as an overridable present-gate channel — parallel after 3k's recorder verbs exist.** Skill prose only in this repo: the recorder-side validation and id-normalization verbs are now IN the contract (reconciled, attempt 9) and land with 3k/h1's binary work. The binary stays subspace-free by checkable criterion. The override script + its committed CI drive suite are the NAMED CROSS-REPO RELEASE CONDITION: the 0.27.0 pre-cut gates on a pinned subspace revision carrying that suite (captain ruling — see Close-out).
5. **The sprint eats its own output:** the moment the recorder can record, this sprint's own remaining gates and validations use it — hand-recording ends.

High-stakes surfaces (the recorder internals, the contract, present-gate skill changes): detached adversarial audit at validation before merge.

## Gate presentations

Present via the float ritual: the override-script lineage in the shaping debrief (probe with a throwaway FIRST; packages in each entity's review room; results are room-resident from the first byte). No subspace → chat presentation per the present-gate template, recorded identically.

## Close-out

1. Contract landing pass verified done (3k's final step, above).
2. **The subspace release condition:** verify the pinned subspace revision carrying the hardened override script and its committed CI-run drive suite; the captain gates this before the tag — the tag does not fire without it or without the captain explicitly narrowing the condition.
3. **Pre-cut audit** — independent reviewer over the assembled sprint before the tag.
4. `go test ./...` and `go test ./... -race` green; `gofmt -w ./cmd ./internal` clean; final clean `git status`.
5. **Cut 0.27.0** per `docs/releasing.md` *(captain authorizes the tag)*.
6. Sprint-close mining is the Shaping FO's, not yours.
