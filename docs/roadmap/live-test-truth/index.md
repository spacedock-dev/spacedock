# Sprint: live-test-truth

**Sprint:** the tasks matching `sprint: live-test-truth`:

```bash
spacedock status --workflow-dir docs/dev --where sprint=live-test-truth
```

The drivable set excludes deliberately deferred candidates:

```bash
spacedock status --workflow-dir docs/dev \
  --where sprint=live-test-truth \
  --where 'sprint-readiness != defer'
```

**Target train:** `v0.27.0` stable.

## Deliverable

Runtime live CI has one truthful desired-state registry and a lean execution
shape. Common behavioral journeys have one runtime-neutral identity and are
required on every supported runtime by default. Host-specific live work proves
only the host boundary that cannot be shared. Every live test and fixture is
either bound to that desired state, deliberately non-gating with a reason, moved
to offline coverage, or removed.

[`docs/runtime-live-ci.md`](../../runtime-live-ci.md) remains the single normative
entry point for live CI. It incorporates the desired-state registry at
[`docs/runtime-live-ci-registry.md`](../../runtime-live-ci-registry.md) and names
the executable reconciliation check.

## Locked semantics

- A **common journey** is required on every supported runtime target by default.
  Only genuine non-applicability earns an exception. Missing support, cost,
  quarantine, and an unwired selector do not.
- A common journey has one canonical exported `TestLiveCommon...` function.
  Its adjacent `liveJourney(...)` call binds the stable journey ID, fixture,
  TODO owners, runtime-neutral exercise, and durable assertion. Authentication,
  launch, output, and liveness differences stay behind runtime adapters.
- The registry can name a desired journey before code exists. Missing code or CI
  invocation is a reconciliation result, not a reason to weaken desired state.
- Each journey references one or more stable fixture IDs and their semantic
  contracts. Source annotations own fixture-ID-to-builder bindings. Multiple
  fixtures can prove the same intention, and multiple journeys can share a
  fixture. An annotated fixture referenced by no registered journey or runtime
  proof is orphaned.
- Runtime-specific live proofs are separate and lean. They prove only a unique
  host substrate boundary and do not repeat common workflow semantics.
- An intentionally unselected live test is not release evidence and carries an
  explicit reason. Otherwise it is promoted, moved to offline coverage, or
  deleted.
- The standing reconciliation guard parses the registry, Go declarations and
  calls, adjacent annotations, and the live workflow. It checks their semantic
  join directly. It does not use a recorded SHA. Fixture resumption and
  unwired-count ratchets remain out.

## Scope

Carved work is discovered from workflow state, never enumerated here. The sprint
groups work in these outcome areas:

- Bind stable desired registry IDs to tests, shared scenarios, and fixture
  builders through source annotations. Keep builder symbols and selector text out
  of registry prose.
- Converge common live behavior on one canonical entry point and promote existing
  standalone common journeys into it.
- Add Pi as a first-class common-journey adapter while keeping its unique
  front-door proof lean.
- Restore required common evidence after the product defects behind current TODO
  quarantines land. This integration work remains deferred until those upstream
  owners complete.
- Repair the AC value-reanchor proof so its oracle can fail when the promised end
  value regresses.
- Move deterministic tests out from behind the live build tag.
- Retire the unreachable legacy Claude pty/tmux lane and select the supported
  Claude-specific recovery proofs.
- Replace Codex's fixed false-red timeout with progress-aware liveness.
- Remove or connect decorative live-lane inputs and unconsumed metrics paths.

One earlier candidate proposed a second Pi smoke and remains deferred. The locked
direction puts its additional boot-contract evidence in the selected Pi
front-door smoke. This direction avoids a duplicate live run.

## Definition of Done

1. `docs/runtime-live-ci.md` is the single normative entry point. It incorporates
   `docs/runtime-live-ci-registry.md` as its desired-state component and names the
   executable reconciliation command.
2. Every common journey has one canonical executable identity. Each target is
   either passing live evidence or an exact target-scoped `TODO(owner)` backed by
   a reproduced product defect. An unverified runnable cell, an unowned gap, or
   a broad quarantine blocks this sprint. The product repairs and later evidence
   restoration belong to `test-behavior-completeness` before the release cut.
3. Every registered fixture ID resolves through a source annotation to one
   concrete builder. Every annotated live fixture is consumed by a registered
   journey or runtime-specific proof. No orphan fixture or unclassified live
   entry point remains.
4. Deterministic tests run under default build tags. Intentional live experiments
   remain unselected with explicit reasons.
5. Runtime-specific lanes contain only unique substrate proofs. These proofs cover
   Claude merged-agent and recovery behavior, Codex front-door behavior, and one
   Pi front-door/subagent smoke. The common runner already exercises the Codex
   front door. The unreachable pty/tmux lane is gone.
6. Retained live-lane inputs affect a real runtime setting, and retained metrics
   paths have a producer and consumer. Codex liveness resets on meaningful
   progress while the suite-wide deadline remains a loose runaway backstop.
7. Every drivable sprint task reaches `done` with `PASSED`. A reproduced,
   target-scoped product gap with an exact owner does not block this sprint and
   does not count as green. The owner moves to `test-behavior-completeness`.
   Missing ownership or missing evidence for a runnable cell blocks completion.
8. `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` pass.
   Changed live surfaces also pass their exact required runtime lanes.

## Out of scope

- Rewriting the shared scenario/executor semantic model.
- Implementing fixture-session resumption, cross-journey batching, or live-run
  optimization before the registry and bindings establish a trustworthy baseline.
- A heavyweight AST-diffing contract lint, ratcheted unwired-count ceiling, or any
  standing enforcement beyond the approved path-scoped reconciliation-SHA guard.
- Product repair owned by another sprint, including `durable-decisions` members.
  This sprint can restore live evidence after such work lands but does not duplicate
  or investigate that work.
- Changes to `internal/gates` fields owned by the separate `durable-decisions`
  sprint.

## Delivered sequence

The implementation sequence was `3d`, then `15e`, then `ys`. Those members are
complete. They delivered truthful assertions and liveness, named lane evidence,
the portable common-journey entry points, the registry, and semantic
reconciliation.

The broad deferred member `rm` is archived as superseded. Exact target-scoped
owners in `test-behavior-completeness` own product repairs and later TODO removal.

## Ideation and staff-review focus

- The three outcome members are `ys`, `3d`, and `15e`. Their ideation is complete.
- This shaping pass used Codex Sol ensigns and `$simple-english` in pragmatic mode.
  The workflow-wide `model: opus` setting does not map to a Codex spawn model.
  The dispatch builder emitted `model: null`. The Sol override from the Captain governed.
- `ys` owns the portable common-journey surface. It reuses the annotation join,
  guard fail-then-pass, canonical entry-point, journey-promotion, and Pi-feasibility
  evidence from its absorbed design inputs.
- `ys` uses one candidate branch with seven reviewable landing steps. The staff
  review must audit its collision map, step boundaries, rollback points, and final
  selector cut as the seam with the highest blast radius.
- `3d` owns truthful live results. It reuses the AC2 false-green proof and adds the
  Codex progress-aware liveness proof.
- `15e` owns named lane evidence. It must map each selector to one evidence claim
  and each retained control or metric to one producer and consumer.
- The staff review must examine the dependency on `durable-decisions`. It must
  state the effect on the `v0.27.0` train if the upstream repair is late.
- The staff review must not count the planned SSOT link, SHA, or guard as current
  evidence before they land.

## Sprint lifecycle

### Shape — Shaping FO

- [x] **Scope-lock** — membership, semantics, the two-document SSOT relationship,
  the single SHA guard, deferred-upstream policy, and `v0.27.0` target are locked.
- [x] **Carve** — sprint fields are stamped, deferred candidates are separated from
  the drivable query, and this index records the real goal, scope, and DoD.
- [x] **Ideate each gated member** — exercise the riskiest mechanism first and do
  not re-ideate an already banked design.
- [x] **Preflight staff review** — Fable checked sprint-wide
  DoD ownership, sequencing, collisions, blast radius, missing scope, and Commander
  cold-boot readiness. Both rounds and their folds are in `staff-review.md`.
- [x] **Present ideation gates** — the Captain approved `3d`, `15e`, and `ys` on
  2026-08-03. Each gate has a recorded pending advance for the Commander to consume.
- [x] **Package** — `dispatch-sprint-execution.md` is the cold-boot Commander package.

### Drive — Commander

- [x] Drive implementation, detached validation where required, and merge for
  the live-test-truth members.

### Close — Shaping FO

- [ ] Fold deferred/non-blocking findings into the next sprint and check the
  released result.

## Delivered state and pre-close record

The desired registry now maps 16 common journeys across four runtime roles. On
2026-08-09, reconciliation derived 48 runnable cells and 16 exact TODO cells:
Claude Opus 1, Claude Sonnet 4, Codex 6, and Pi 5.

- `26n` remains the exact Sonnet, Codex, and Pi owner for
  `default-headless-gate-stop` on current `main`. PR #583 is the Sonnet and Codex
  promotion; Pi remains TODO until separate Pi evidence exists.
- `zh` (`zhcb4bcz1qgcn7ajx2ctxpxk`) replaces archived rejected `zbc` as the
  Sonnet, Codex, and Pi owner for `rejection-flow`.
- PR #643 merged as `90b6e0e61` and repaired optional journey-metrics reporting.
- `0a` (`0aqnm6v8ajns6cpsknxn9wf2`) owns optional manual Pi CI in
  `test-behavior-completeness`.
- Product repairs and later evidence restoration remain outside this sprint.

The delivery work is complete. Source reconciliation rejects unclassified live
entry points. The close procedure joins TODO IDs to the mutable state checkout
and rejects inactive owners without making stable code CI depend on later state.
