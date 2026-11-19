# Sprint 019x — pre-flip cleanups

**Goal:** land the small ready cleanups and the two nearly-done pre-flip deliverables
(the README reconciliation and the survey fix) so **0.19.7** ships clean ahead of the
0.20.0 marketplace flip.

**Deliverable:** spacedock **0.19.7** cut on `next` with every `ready` member merged.

**Why a sprint, not loose tasks:** these are the things that must be true *before* the
flip. Bundling them gives one definition-of-done and one integration test instead of
six independently-tracked merges, and a single conn-to-drive handoff.

## Members

Membership is the query, not this table (the table is just the readable view):

```bash
spacedock status --workflow-dir docs/dev --where sprint=019x-pre-flip-cleanups --where 'sprint-readiness != defer'
```

| Entity | group | readiness | Drive owner | Remaining work |
|--------|-------|-----------|-------------|----------------|
| `nb` readme-main-flip-reconciliation | readme | ready | **captain** | merge PR #322 → ceremony + close #213/#220/#315 |
| `xn` survey-signal-correction | survey | ready | **captain** | AC-3 live drive → approve validation gate → merge |
| `78` brew-cask-message-cleanup | release-hygiene | ready | **captain** | approve validation gate (PASSED) → merge |
| `qy` internal-status-test-hygiene | release-hygiene | ready | **Commander** | full cycle: fix root `go test` + gofmt-dirty files |
| `e30` dispatch-path-teamname-sanitize | dispatch-hygiene | ready | **Commander** | full cycle: sanitize `team_name` path + fix comment |
| `jh` dispatch-build-flag-form-version-skew | dispatch-hygiene | ready | **Commander** | full cycle: runtime-doc / accepted-binary compat |
| `5h0` ban-readme-substring-assertions | proof-policy | **defer** | — | blocked on PR #315 landing on `main` |

**Drive split:** `nb` / `xn` / `78` are done or nearly done and sit at **captain gates**
(merge / approve) — the captain resolves those; the Commander runs any post-gate
ceremony. `qy` / `e30` / `jh` are backlog cleanups the **Commander drives** end-to-end.
`5h0` is carried but deferred (recorded, not silently dropped).

## Definition of Done

Checkable, and able to fail today (not a tautology):

1. Every `ready` member is `done` / PASSED and merged to `next`.
2. **`go test ./...` from the repo root is green with `.spacedock-state` present.**
   (RED today — `qy`'s bug. This is the sprint-wide integration test.)
3. `gofmt -l ./...` is empty (`qy`).
4. `goreleaser check` exits 0 (`78`).
5. spacedock **0.19.7** is stamped and cut on `next` (captain-gated release ceremony).
6. The deferred member (`5h0`) is explicitly recorded as carried, with its blocker.

## Integration test (sprint-wide proof, above per-entity validation)

```bash
go test ./...        # from repo root, with docs/dev/.spacedock-state present — must be green
gofmt -l ./...       # must be empty
goreleaser check     # must exit 0
spacedock --version  # must report 0.19.7 after the release cut
```

## Out of scope

The 0.20.0 flip itself (`pj`, sprint 2). Distribution (`v3` / `5w` / `44`) and the docs
site (`wv` / `e6`) — post-flip. The `spacedock:roadmap` skill build — only if this dry
run proves the construct (graduation, per the proposal).

## Outcome (delivered 2026-06-08)

All five `ready` members shipped to `next` and **spacedock 0.19.7 was cut** (commit
`a35c7002`, annotated tag `v0.19.7`):

- `nb` readme → #322 · `78` brew-cleanup → #323 · `xn` survey → #324 ·
  `e30` dispatch-path → #325 · `qy` test-hygiene → #326

DoD met: `go test ./...` from the repo root is green with `.spacedock-state` + the live
`.claude/worktrees` scratch present (the integration test that started RED — `qy`'s bug),
`gofmt -l` empty, `goreleaser check` exit 0, 0.19.7 stamped + tagged. Deferred and
recorded: `jh` (same-contract command-form skew — needs a gating ideation + a spike
against a real old binary) and `5h0` (blocked on PR #315 landing on `main`).

This sprint was the **convention-only dry run** of the construct
(`docs/dev/_proposals/sprint-roadmap-construct.md`): prose + frontmatter + the native
`--where` query, zero new binary code. The staff-review mechanism caught real readiness
blockers before driving (and refuted an over-cautious coherence flag); the two-tier
FO↔Commander handoff proved realizable only as a separate top-level session — a subagent
cannot originate dispatch.
