# Durable-decisions Commander package — cold-boot execution dispatch

You are the Commander: a fresh FO session driving durable-decisions to its target train (0.27.0 — the movable line in the index). The four original ideation gates are closed-approved; `vn` joined by captain amendment on 2026-07-22 after live recorder dogfood reproduced a leaky folder-artifact commit boundary, and must complete its ideation gate before implementation. Your job: implementation → validation → done per member, merge to main, the contract landing pass, the pre-cut audit, the cut.

## Boot order

1. Load the FO contract (`spacedock:first-officer` skill) and engage the `docs/dev` workflow. NOTE: workflow discovery is currently ambiguous repo-wide (a landed test fixture registers as a second workflow) — pass `--workflow-dir docs/dev` explicitly on every auto-discovering command (`state commit`, `new`); the fix is seeded as backlog task `ab`.
2. Read `docs/roadmap/durable-decisions/index.md` IN FULL — the Constraints (inherited rulings + the recording-identity ruling), the Responsibility boundary table and its five rules, and the dogfooding change protocol bind this session's own conduct.
3. Read BOTH staff-review seats in this folder — `staff-review.md` (fable) and `staff-review-codex.md` (codex, NOT READY at the time) — plus the index Constraints recording the codex-seat rulings (authorization-only consumption, the named subspace release condition, digest domains, recording identity). Every codex material finding is closed in durable state: room evidence tracked, adoption provenance recorded, the contract reconciled at sha256 681b2348... via 3k gate attempt 9 (the first FO-identity closure under the recording-identity ruling).
4. The design authority for 3k, h1, xb, and 02av is ONE spec: `docs/dev/.spacedock-state/durable-gate-approval-pending-blockers/gate-resolution-frontmatter-contract.md`, owner-tagged per section. Amendments route to section owners per the change protocol; nobody forks the document. vn owns the orthogonal state-commit boundary and does not amend gate semantics.
5. Drivable set: `spacedock status --workflow-dir docs/dev --where sprint=durable-decisions` (returns five members; membership is the query, not this prose).

## Consuming the recorded approvals

Each original member's approval is in its entity `gates:` frontmatter: a closed attempt with `resolution.decision: approve` and a pending `advance` application IS the captain's approval — apply once via the normal transition path, mark `consumed`, never re-ask a closed gate. vn joined after that approval round and still requires its ideation gate. Every briefing here binds a byte-verifiable frozen artifact in the entity's review room (verify digests against the room, never by re-hashing the live entity). The drift waiver: approval-directed edits and captain-approved preflight folds recorded in the same gates record are not drift.

## Members, landing order, and binding conditions

1. **3k — the recorder — leads.** The clean unreleased-v1 recorder at validated commit `024a2c56` accepts only the canonical gates projection, records semantic Briefing/Result/chat inputs, verifies the complete association inventory from frozen canonical Briefing bytes, and preserves only unrelated entity bytes plus the explicit opaque application boundary. Prototype compatibility and replay machinery are deleted. The **contract landing pass** remains its final pre-merge step.
2. **vn — folder-form state commit — is the persistence companion.** It may implement in parallel after its ideation gate closes, but it must land before xb's end-to-end presentation journey or sprint close-out. One `state commit <folder-entity>` must durably include the entity index plus new/changed room reports and artifacts while excluding sibling dirt. The exact live regressions are state commits `d8e4180c`/`2c616b7e`: the package required a manual exact-path commit while the subsequent index mutation committed normally.
3. **h1 — the application layer — after 3k**, extending the same binary (never a second gates writer). Ships the one-use application lifecycle with the cross-attempt single-pending invariant; the blocker-evaluator half is a RECORDED DECLINE — do not build it; its promotion condition is in the entity.
4. **02av — triage on advisory records — parallel with h1** (prose + fixtures + one offline check, zero product LOC). Round records are hand-authored interim (the contract says so); the rounds-plumbing generalization is explicitly out of scope. Carries the release line: the seeded disproportionate finding produces a recorded decline and a zero-line diff in live replay at validation.
5. **xb — presentation as an overridable present-gate channel — after 3k's recorder verbs exist and before the final vn-backed journey.** It must present the canonical Briefing's question and complete artifact set, retain the exact Result, and produce an honest association from what was actually presented. A single-file Subspace float does not count as complete-package presentation merely because `briefing.json` names additional artifacts. The binary stays subspace-free. The override script + its committed CI drive suite remain the named cross-repo release condition.
6. **The sprint eats its own output:** the recorder owns remaining gates; once vn lands, the same `state commit` invocation must durably commit each folder-form gate package and its entity mutation without a manual Git fallback.

High-stakes surfaces (the recorder internals, the contract, present-gate skill changes): detached adversarial audit at validation before merge.

## Gate presentations

Present via the float ritual: probe with a throwaway first; packages and exact results are room-resident. Until xb lands, a `subspace:r` one-file float is advisory evidence only: it uses a generated single-file Briefing and generic question, so never mint a full-package association unless the reviewer actually received the canonical question and every mapped artifact. No Subspace → chat presentation per the present-gate template, recorded identically.

## Close-out

1. Contract landing pass verified done (3k's final step, above).
2. **Folder-package durability:** vn is landed and a live folder-form gate journey proves one `state commit` commits the index, Briefing, legible gate artifact, exact Result, and association while leaving sibling dirt untouched.
3. **The subspace release condition:** verify the pinned subspace revision carrying the hardened override script and its committed CI-run drive suite; the captain gates this before the tag — the tag does not fire without it or without the captain explicitly narrowing the condition.
4. **Pre-cut audit** — independent reviewer over the assembled sprint before the tag.
5. `go test ./...` and `go test ./... -race` green; `gofmt -w ./cmd ./internal` clean; final clean `git status`.
6. **Cut 0.27.0** per `docs/releasing.md` *(captain authorizes the tag)*.
7. Sprint-close mining is the Shaping FO's, not yours.
