---
id: 87j19afq4tj5te1hjvgd6rs4
title: pr-merge mod hardcodes base branch `next` (pre-flip); refit to `main`/config-driven
status: ideation
source: "0202 Commander drive (2026-06-13). The pr-merge mod (v0.12.1) opens PRs against `next`; the Commander overrode the base to `main` per the dispatch doc on every merge. Same post-flip stale-trunk class as dispatch reconcile."
group: cleanup
sprint: 0203-fo-efficiency
started: 2026-06-14T05:04:07Z
---

The `pr-merge` merge-hook mod (`docs/dev/_mods/pr-merge.md`, v0.12.1) opens code-branch PRs against base `next`, the pre-flip trunk. Post-flip the trunk is `main`, so every merge this drive required a manual base override.

## Problem

`_mods/pr-merge.md` is FO-read prose (the binary parses only its frontmatter via `internal/dispatch/mods.go`, never the body). The body resolves the integration trunk to a hardcoded literal `next` in **eight** places — not only the `gh pr create --base next`, but every trunk-relative reference: the header description, the split-root intro, the `Branch: {branch} -> next` draft line, the `git diff --stat origin/next...` change count, the rebase-target sentence, the branch-push remote note, and the local-merge fallback target (`merge onto next`, `merge-commit SHA on next`).

The 2026-06-08 flip made `main` the trunk (`next` is dev-only per `docs/releasing.md`); the mod was written pre-flip (v0.12.1) and never refit. The 0202 Commander manually overrode `--base main` on all five member PRs. A fresh FO following the mod literally would open PRs against the deprecated `next` branch and compute change counts against the wrong base.

This is the **same root cause** as the sibling `sr` task (`dispatch-reconcile-deconflate-repo-hygiene`): `reconcile.go` classD/classE hardcode `origin/next` in Go. The captain's constraint binds the pair: **settle ONE trunk-config source jointly; DoD: no helper/mod/ref/doc resolves the integration trunk to `next` post-flip.**

## Coupling with `sr` — ONE canonical trunk-config source (reconciled)

`sr` LEADS the canonical-source design; this task CONSUMES it for the PR base. I do **not** invent a second source. Reconciled against `sr`'s firm design (its entity, committed `b84e3b4a`); my earlier draft assumed a top-level `trunk:` key — `sr` settled a different, better source and I bind to it:

- **Canonical source (`sr`'s, adopted):** `stages.defaults.trunk` in the workflow README frontmatter (`docs/dev/README.md`); for this repo `stages.defaults.trunk: main`. Parsed by `status.ParseStagesWithDefaults`, which returns `(stages, defaults map[string]string)` — `defaults["trunk"]` is free, and `reconcile.go` **already imports and calls this parser** (`reconcile.go:319`). No new parser, no top-level key. Absent key → resolver defaults to `main`, never `next`.
- **One resolver (`sr`'s `resolveTrunk(workflowDir)`):** `sr` extracts `resolveTrunk` out of `classD`/`classE`; it reads `stages.defaults.trunk` and falls back to `main`. This task CONSUMES the same `resolveTrunk` — it does NOT define a parallel resolver. The single symbol IS the one-source proof (AC-3): both `sr`'s reconcile test and this task's command test exercise the same `resolveTrunk`.

**Anti-conflation guard (carried from `sr`'s design):** the PR base MUST NOT route through `cli.devBranch` (`internal/cli/frontdoor.go`'s `var devBranch = "next"`). That variable is the *marketplace channel stamp* (`=main` → stable entry, `=next` → edge entry, per `docs/releasing.md`), a different axis where `next` is *correct*. The pr-merge base reads `stages.defaults.trunk` only. Conflating the two would re-create the bundling bug.

**ALIGNMENT FLAG for the gate (open, defer to `sr`):** my mod is FO-read PROSE, so the FO can only become config-driven by *running* the binary to get the base — it cannot read a Go function. The consumer-side surface this task needs is a thin command, `spacedock dispatch trunk --workflow-dir docs/dev`, that calls `sr`'s `resolveTrunk` and prints the resolved trunk. This is a ~12-line `case "trunk":` on the existing `dispatch.Run` switch (`internal/dispatch/dispatch.go:80`), host-neutral (README parse only). Open question routed to `sr`: whether `resolveTrunk` is exported from `internal/dispatch` (so the `dispatch trunk` case calls it directly) or lifted to `internal/status` — I consume wherever it lands. If the command is rejected, this task either adds the command itself or cannot meet the non-grep AC (re-scope) — surfaced for the staff reviewer / gate.

## Proposed approach (firm)

1. **Refit the mod prose** (`docs/dev/_mods/pr-merge.md`, the split-root variant — NOT the repo-root `mods/pr-merge.md`, which is a separate non-split-root mod already on `main`) so the trunk is resolved by running the binary, not a literal:
   - The FO resolves the base once at the start of the merge hook: `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` (the standing trunk, default `main`).
   - All eight trunk-relative references become `{base}` / `origin/{base}` instead of the `next` literal — header description, split-root intro, `Branch: {branch} -> {base}`, `git diff --stat origin/{base}...{branch}`, the `gh pr create --base {base}`, the rebase-target sentence, the push-remote note, and the local-merge fallback (`merge onto {base}`, `merge-commit SHA on {base}`).
   - Bump the mod `version:` (0.12.1 -> next patch) since the hook contract changed.
2. **Consume `sr`'s `stages.defaults.trunk: main` README declaration and `resolveTrunk` resolver.** `sr`'s task adds the README key and the resolver; this task does not duplicate either. If `sr` defers the README-key edit, this task adds the one `stages.defaults.trunk: main` line — but the resolver stays `sr`'s single symbol.
3. **Expose `spacedock dispatch trunk --workflow-dir <dir>`** that prints `resolveTrunk(workflowDir)` to stdout. This is the FO-runnable oracle the refit mod points at, and the surface this task's non-grep AC tests.

The mod is no longer self-resolving: its base comes from a command that reads config, so a future trunk change is a one-line README edit, no mod edit and no per-merge override.

## Acceptance criteria

**AC-1 — `spacedock dispatch trunk` resolves the base from `stages.defaults.trunk`, defaulting to `main`, never `next`.**
Verified by: a Go command-level test in `internal/dispatch` (mirroring `reconcile_test.go`'s fixture+`Run`+stdout pattern) that writes a fixture README and drives the command:
  - fixture README with `stages.defaults.trunk: main` -> command prints `main`;
  - fixture README with **no** `trunk` key -> command prints `main` (default), proving absence does not fall through to `next`;
  - fixture README with `stages.defaults.trunk: ftrunk` (a sentinel that is neither `main` nor `next`) -> command prints `ftrunk`, proving the value is *sourced from the fixture config*, not a literal in the code under test.
The expected value comes from the fixture README (an independent artifact the test writes), not from the mod prose or a code literal — this is the legitimate "relationship between independent values" the boundary guard permits, not a prose/code grep. **A `next`-defaulting or `next`-hardcoded resolver reds the second case.** (Independent oracle: the fixture's configured trunk drives the output.) The sentinel `ftrunk` matches `sr`'s AC-1 fixture trunk so both halves test the same resolver against the same sentinel.

**AC-2 — the resolver is config-sourced, not a literal: a fixture trunk other than `next`/`main` flows end-to-end to the PR base the FO would use.**
Verified by: the AC-1 third case (`ftrunk`) is itself the proof for the consumer path — the command the refit mod invokes returns the fixture's trunk verbatim. Because the mod's `gh pr create --base "$BASE"` takes the command's stdout, a test proving the command honors the fixture trunk proves the FO opens the PR against that trunk. No grep over the mod text is involved; the oracle is the fixture config -> command stdout relationship.

**AC-3 — `reconcile.go` classD/classE and the `pr-merge` base read the SAME resolver (one source, not two).**
Verified by: `sr`'s `resolveTrunk` is a single function; `sr`'s reconcile test (fixture `stages.defaults.trunk: ftrunk` drives classD/classE detection against `origin/ftrunk`) and this task's `dispatch trunk` test both exercise it against a fixture README declaring the same key. A `grep -c` over the codebase is **not** the proof — the proof is that both tests pass against the *same* `resolveTrunk` symbol (if either half forked a second source, one test would be reading a value the other never wrote). Reconciled with `sr`: source settled as `stages.defaults.trunk` via `resolveTrunk`; if at the joint gate `sr`'s landed symbol moves (e.g. lifted to `internal/status`), this AC's test target follows it — there must be exactly one resolver.

**Note on what is NOT an AC:** "the mod text no longer contains `next`" is real authoring work (it is in the doc diff and the stage-report checklist) but is **not** an acceptance criterion — a substring sweep over `pr-merge.md` is the banned prose-grep (`internal/contractlint/boundary_guard_test.go`). The behavior that matters — the FO opens the PR against the configured trunk — is proven only by the command-level test (AC-1/AC-2), never by reading the mod.

## Test plan

| What | Where | Cost | Type |
|---|---|---|---|
| `dispatch trunk` resolves `main` / default / sentinel `ftrunk` | `internal/dispatch/trunk_test.go` (new), fixture README + `dispatch.Run([]string{"trunk","--workflow-dir",dir}, ...)` asserting stdout | low (~1 file, mirrors `reconcile_test.go:432` `TestReconcileIncludeScope` harness) | Go command-level fixture test |
| `resolveTrunk` (sr's): `stages.defaults.trunk: X` -> X, absent -> `main` | `sr`'s `TestResolveTrunk` (sr owns; this task does not duplicate) | low | Go unit test (sr's) |
| One-source check (AC-3) | both `sr`'s reconcile test and this command test pass against the single `resolveTrunk` symbol | low (no new test; assertion is "both green, one symbol") | reconciliation at gate |

No live workflow / CLI smoke test needed: the claim is command-level (does the command print the configured trunk) and unit-level (does the resolver read config). The end-to-end "FO opens PR against `main`" is a prose instruction whose only mechanical part — resolving the base — is the command under test; the `gh pr create` call itself is not exercised in CI (it is outward-facing and gated on captain approval).

## Spike result (riskiest unknown, exercised first)

The design rests on one unverified mechanism: **does `status.ParseStagesWithDefaults` surface `stages.defaults.trunk` from a README, and does an absent key yield empty (so the resolver owns the `main` default)?** This is the parser `sr`'s `resolveTrunk` reads and that `reconcile.go` already imports (`reconcile.go:319`). Exercised before committing the design (throwaway test, removed after):

```
=== RUN   TestSpikeTrunkStagesDefaults
    SPIKE2 OK: trunk=main; sentinel=ftrunk; absent=empty
--- PASS
```

Result: a README with `stages.defaults.trunk: main` → `defaults["trunk"]=="main"`; `stages.defaults.trunk: ftrunk` → `"ftrunk"` (arbitrary value flows through); absent key → `defaults["trunk"]==""` (so `resolveTrunk` owns the `main` default, and the fallback path is reachable, not dead code). **Confirmed**: the canonical source is readable with zero parser changes via the parser reconcile already uses. (My earlier spike used the wrong parser — `ParseFrontmatter` for a top-level `trunk:` key — and is superseded by this one once `sr`'s `stages.defaults.trunk` source was settled.) The `dispatch trunk` command is a ~12-line `case "trunk":` on the existing `dispatch.Run` switch (`internal/dispatch/dispatch.go:80`), host-neutral (README parse only, no `~/.claude` coupling) — proven feasible by the existing `case "reconcile"` sibling that drives the same fixture-README harness.

## Doc diff (mod refit — applied at implementation)

Concrete before/after for `docs/dev/_mods/pr-merge.md`. The FO resolves the base once near the top of the merge hook; all literals become `{base}`/`origin/{base}`.

- **Frontmatter `version:`** `0.12.1` -> bump one patch (hook contract changed).
- **Frontmatter `description:`** "Open a code-branch PR to next at the merge boundary…" -> "Open a code-branch PR to the configured trunk at the merge boundary…"
- **Intro (line 9):** "base branch `next`" -> "base branch resolved from the workflow's `stages.defaults.trunk` config (default `main`)".
- **Merge hook, new opening step:** add "Resolve the PR base once: `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` — the workflow's configured integration trunk (default `main`). Use `$BASE` for the draft, the diff stat, the `gh pr create --base`, and the fallback merge target below."
- **Draft `Branch:` line (45):** `{branch} -> next` -> `{branch} -> {base}`.
- **Draft `Changes:` line (46):** `git -C {worktree} diff --stat origin/next...{branch}` -> `... origin/{base}...{branch}`.
- **Approval rebase sentence (56):** "rebases the code branch onto `origin/next`" -> "onto `origin/{base}`".
- **Push note (58):** "(`origin` = the main repo, base branch `next`)" -> "(`origin` = the main repo, base branch `{base}`)"; "Do NOT push `next`" -> "Do NOT push the trunk".
- **`gh pr create` (62):** `--base next` -> `--base "$BASE"`.
- **Fallback merge target (100):** "a local `--no-ff` merge … onto `next`" -> "onto `{base}`".
- **Fallback sentinel (103):** "compute the merge-commit SHA on `next`" -> "on `{base}`"; "The SHA must already exist on `next`" -> "on `{base}`".

No other doc files reference the **split-root** pr-merge base. The dev `README.md` gains `stages.defaults.trunk: main` — that diff is `sr`'s canonical-source change (in `sr`'s doc-diff already); this task does not duplicate it. The repo-root `mods/pr-merge.md` is a separate non-split-root mod already on `main` and is out of this task's scope.

## Notes
Mod files are refit-or-worker territory, not FO direct edits. Coupled pair with `sr` (`dispatch-reconcile-deconflate-repo-hygiene`): `sr` OWNS the single source (`stages.defaults.trunk` in `docs/dev/README.md`) + the `resolveTrunk` resolver; this task CONSUMES both for the PR base via a thin `dispatch trunk` command. Anti-conflation guard carried: base reads `stages.defaults.trunk`, never `cli.devBranch`. High-stakes (merge/release machinery) → expect a staff review before the gate. A parallel 5-surface trunk-config audit cross-checks the `next`-as-trunk inventory at the gate.

## Stage Report: ideation

- DONE: The design refits the `pr-merge` mod's PR base branch from hardcoded `next` to `main`, read from the SAME single canonical trunk-config source `sr` defines — config-driven, not a per-merge manual override.
  Firm design + doc diff refit all 8 trunk-relative `next` literals in the split-root `docs/dev/_mods/pr-merge.md` to `{base}`, resolved by `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` reading `sr`'s `stages.defaults.trunk` source via `sr`'s `resolveTrunk`. Reconciled with `sr`'s firm design (commit `b84e3b4a`) over two SendMessage rounds; my initial top-level-`trunk:` assumption was corrected to `sr`'s `stages.defaults.trunk`.
- DONE: An acceptance criterion proven by a real check (a test/fixture that reds if the mod opens a PR against `next` or resolves the base to `next`) — never a prose/string grep over the mod text.
  AC-1/AC-2: a command-level Go test (mirrors `reconcile_test.go:432` fixture+`Run`+stdout) where a fixture README's `stages.defaults.trunk` drives `dispatch trunk` stdout; the absent-key case reds a `next`-defaulting resolver, the sentinel `ftrunk` case proves config-sourcing (independent oracle = fixture config, not mod prose). Explicit non-AC note records the banned prose-grep (`boundary_guard_test.go`).
- DONE: The single trunk-config source is consistent with `sr`'s design (the coupled pair settles ONE source) — exercised that the `pr-merge` mod actually reads it.
  One `resolveTrunk` (sr's) over `stages.defaults.trunk`; AC-3 proves one-source by both `sr`'s reconcile test and this command test passing against the same `resolveTrunk` symbol on the shared `ftrunk` sentinel. De-risk spike (throwaway, removed) against the REAL parser `ParseStagesWithDefaults` (the one reconcile already imports) confirmed `stages.defaults.trunk` surfaces `main`/`ftrunk` and yields empty when absent — exercised before committing the design.

### Summary

Refit designed and reconciled with `sr`: the split-root `pr-merge` mod stops hardcoding `next` and resolves its PR base by running `spacedock dispatch trunk`, a thin command wrapping `sr`'s `resolveTrunk`, which reads the single canonical `stages.defaults.trunk` key in the workflow README (default `main`) — the same source+symbol `sr` consumes for `reconcile.go` classD/classE. The riskiest unknown was spiked green against the actual parser before the design firmed. Anti-conflation guard carried: the base never routes through `cli.devBranch`. One open item flagged for the gate: whether the `dispatch trunk` command wrapper lands and where `resolveTrunk` is exported from — this task's non-grep AC depends on the command; defer the symbol's home to `sr`.
