# Live-test-truth Commander package — cold-boot execution dispatch

You are the Commander for the `live-test-truth` sprint. Shaping is complete.
The Captain approved the three ideation gates on 2026-08-03. Each approval has
a pending advance. Drive implementation, validation, merge, assembled review,
and release evidence. Do not re-present the ideation gates.

## Boot order

1. Load the `spacedock:first-officer` skill in a fresh top-level session.
2. Engage the `docs/dev` workflow. Pass `--workflow-dir docs/dev` explicitly.
3. Read this package, `index.md`, and `staff-review.md` in full.
4. Read `docs/runtime-live-ci.md` and `docs/runtime-live-ci-registry.md` in full.
5. Query the full and drivable sprint sets from the commands in `index.md`.
6. Read each current entity before you consume its approved gate.
7. Create your own implementation and review team. Do not act as an implementer.

The drivable set is `3d`, `15e`, and `ys`. Member `rm` remains deferred and
release-blocking. Archived design inputs do not resolve through `status`. Read
them by their paths under `docs/dev/.spacedock-state/_archive/` when a current
entity cites them.

## Captain-approved value

Each member must deliver an operator-visible test result. A pure component change
does not satisfy its gate.

1. `3d` makes wrong decisions fail and keeps active Codex runs alive.
2. `15e` makes every selected live minute buy named evidence.
3. `ys` gives Claude, Codex, and Pi one common journey identity.

Preserve this order: `3d` then `15e` then `ys`. Do not dispatch later members
for implementation before the earlier member merges to `main`.

## Wave 1 — truthful results (`3d`)

Consume the recorded ideation approval for `3d`. Then dispatch implementation.

The member owns two visible results:

- A wrong AC2 decision produces a red result through a durable two-way oracle.
- Meaningful Codex JSONL progress resets the quiet timer.

The product quiet budget is 60 seconds. The earlier 30-second value was only a
probe setting. Preserve the suite-wide timeout as a loose runaway backstop.

The plan expects about 9 files and `+331/-161` lines. Treat material growth beyond
its recorded tolerance as scope drift. Require focused negative controls, then
the repository test commands. Use detached adversarial review before merge because
this member changes live-result truth.

Merge `3d` before work starts on `15e`. Its old Codex selector is only pre-`ys`
evidence. Member `ys` owns the final selector migration.

## Wave 2 — named lane evidence (`15e`)

After `3d` merges, consume the approval for `15e`. Then dispatch implementation.

The member must leave a clear selector-to-claim map. It must also leave a clear
producer-to-consumer map for every retained control, metric, and artifact.

Required visible results include:

- one Pi smoke for the front door, child dispatch, durable output, and boot contract;
- supported Claude merged-agent and recovery proofs without the dead PTY lane;
- 12 deterministic tests moved out of live-only selection;
- no decorative effort input or unconsumed metrics path;
- lower duplicate live spend with the retained cost visible.

The plan expects about 22 files and `+350/-1750` lines. The deletion is intentional.
Treat new live selectors or retained dead surfaces as scope drift. Validate the
workflow and release guards on the post-`3d`, pre-`ys` selector shape.

Merge `15e` before work starts on `ys`.

## Wave 3 — portable common journeys (`ys`)

After `15e` merges, consume the approval for `ys`. Then dispatch implementation.

This member has the largest blast radius. Keep its one branch, but review and
commit its seven recorded landing steps as separate increments. The final selector
cut, guide update, first reconciliation SHA, and path guard land last and together.

Required visible results include:

- one `TestLiveSharedScenarios/<journey-id>` identity for all common journeys;
- the same `shallow-boot` selector on Claude, Codex, and Pi;
- all 16 desired journeys represented and executable;
- source bindings for stable fixture IDs;
- zero orphan fixtures and zero unclassified live entry points;
- runtime differences contained in adapters;
- one normative live-CI entry point that incorporates the desired registry.

The registry stores desired state. Do not add reconciliation history to its table.
Keep builder symbols and selector text in source bindings, not registry prose.
The source annotation grammar is in
`docs/dev/.spacedock-state/_archive/bind-live-registry-to-source.md`.

The plan expects about 34 files and `+2060/-1650` lines. Its full live-evidence
estimate is $2.46. The approved ceiling, including reserve, is $3.08. Stop and ask
the Captain before exceeding that ceiling or changing the seven-step design.

## Deferred member and release decision

Do not dispatch `rm` until all recorded upstream product dependencies land. Keep
its `sprint: live-test-truth` stamp and `sprint-readiness: defer` state visible.

Before the `v0.27.0` cut, promote and complete `rm`. Its restored journeys must
produce real, non-skipped live evidence. If an upstream repair is late, stop and
ask the Captain to move the train or change the release bar.

## Gates and review

- Ideation gates are approved. Consume each once, only when its wave starts.
- Ensigns implement and write stage reports. The FO owns entity state.
- Run detached validation for each member. Present validation gates to the Captain.
- Never self-approve a gate.
- Use the normal feedback route for rejected validation. Do not bypass it.
- After all members merge, dispatch an independent sprint-wide antipattern review.

The assembled review must reconcile the registry against source and CI. It must
prove that the recorded SHA is current and that the path guard fails on drift.
The review must also prove that no lane, fixture, selector, control, or metric is
unaccounted for.

## Completion evidence

Before sprint completion:

1. Verify every Definition of Done item in `index.md` with current evidence.
2. Run `gofmt -w ./cmd ./internal`.
3. Run `go test ./...`.
4. Run `go test ./... -race`.
5. Run every changed runtime lane required by the member plans.
6. Confirm a clean worktree and merged `main` state.
7. Present the assembled review and release evidence to the Captain.
8. Cut `v0.27.0` only after explicit Captain authorization.

Escalate a material scope fork, a third feedback cycle, a live-cost overrun, or
an upstream release block. Otherwise, keep driving in the approved order.
