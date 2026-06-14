---
id: 87j19afq4tj5te1hjvgd6rs4
title: pr-merge mod hardcodes base branch `next` (pre-flip); refit to `main`/config-driven
status: validation
source: "0202 Commander drive (2026-06-13). The pr-merge mod (v0.12.1) opens PRs against `next`; the Commander overrode the base to `main` per the dispatch doc on every merge. Same post-flip stale-trunk class as dispatch reconcile."
group: cleanup
sprint: 0203-fo-efficiency
started: 2026-06-14T05:04:07Z
worktree: .worktrees/spacedock-ensign-pr-merge-mod-base-branch-post-flip
---

The `pr-merge` merge-hook mod (`docs/dev/_mods/pr-merge.md`, v0.12.1) opens code-branch PRs against base `next`, the pre-flip trunk. Post-flip the trunk is `main`, so every merge this drive required a manual base override.

## Problem

`_mods/pr-merge.md` is FO-read prose (the binary parses only its frontmatter via `internal/dispatch/mods.go`, never the body). The body resolves the integration trunk to a hardcoded literal `next` in **eight** places — not only the `gh pr create --base next`, but every trunk-relative reference: the header description, the split-root intro, the `Branch: {branch} -> next` draft line, the `git diff --stat origin/next...` change count, the rebase-target sentence, the branch-push remote note, and the local-merge fallback target (`merge onto next`, `merge-commit SHA on next`).

The 2026-06-08 flip made `main` the trunk (`next` is dev-only per `docs/releasing.md`); the mod was written pre-flip (v0.12.1) and never refit. The 0202 Commander manually overrode `--base main` on all five member PRs. A fresh FO following the mod literally would open PRs against the deprecated `next` branch and compute change counts against the wrong base.

This is the **same root cause** as the sibling `sr` task (`dispatch-reconcile-deconflate-repo-hygiene`): `reconcile.go` classD/classE hardcode `origin/next` in Go. The captain's constraint binds the pair: **settle ONE trunk-config source jointly; DoD: no helper/mod/ref/doc resolves the integration trunk to `next` post-flip.**

## Coupling with `sr` — ONE canonical trunk-config source (converged)

`sr` LEADS the canonical-source design; this task CONSUMES it for the PR base. We converged over three SendMessage rounds; the pair has settled ONE source (`sr`'s firm design, latest commit `94d89dd6`). I bind to all of it:

- **Canonical source (settled):** a **top-level `trunk:` key** in the workflow README frontmatter (`docs/dev/README.md`), a sibling of `state:` — for this repo `trunk: main`. NOT nested under `stages.defaults`. The decisive rationale (`sr`'s, verified by my spike below): `status.ParseFrontmatter` renders nested mappings as empty string (`frontmatter.go:123-126`), so a `stages.defaults.trunk` key is **invisible** to the simple flat-map parser the `dispatch trunk` command and reconcile use — forcing a second, heavier parser (`ParseStagesWithDefaults`) and a second read path. Top-level `trunk:` keeps exactly ONE parser (`ParseFrontmatter`) behind ONE resolver. `internal/status` is already imported in reconcile (`reconcile.go:17`), so `resolveTrunk` adds NO new import — but its `ParseFrontmatter`-on-the-README is a NEW one-line call site (reconcile's existing `ParseFrontmatter` at `:296` reads per-ENTITY files; the README is currently read via `ParseStagesWithDefaults` at `:319`). Absent key → resolver returns `main`, never `next`.
- **One resolver (`sr`'s `resolveTrunk(workflowDir)`, unexported in `internal/dispatch`):** reads `status.ParseFrontmatter(README)["trunk"]`, returns `main` on empty/absent. Both callers live in `internal/dispatch` — `reconcile.go`'s `classD`/`classE` consume its return, and the `dispatch trunk` command (`dispatch.go`) prints its return — so one same-package symbol serves both, no export needed. `internal/status` is already imported (`reconcile.go:17`), so no new import; `resolveTrunk`'s `ParseFrontmatter`-on-the-README is a new one-line call site (corrected mis-citation: reconcile's `ParseFrontmatter` at `:296` reads per-entity files, not the README). This task CONSUMES the same `resolveTrunk` — it does NOT define a parallel resolver. The single symbol IS the one-source proof: reconcile and the command cannot diverge.

**Anti-conflation guard (carried from `sr`'s design):** the PR base MUST NOT route through `cli.devBranch` (`internal/cli/frontdoor.go`'s `var devBranch = "next"`). That variable is the *marketplace channel stamp* (`=main` → stable entry, `=next` → edge entry, per `docs/releasing.md`), a different axis where `next` is *correct*. The pr-merge base reads the top-level `trunk:` key only. Conflating the two would re-create the bundling bug.

**Division of labor (settled with `sr`, confirmed both ways):** `sr`'s **AC-5 ships the `dispatch trunk` command + `TestDispatchTrunkCommand` + the `resolveTrunk` resolver + the README `trunk: main` key + the dispatch usage-text registration.** That command-level test is the **shared oracle** this pair binds to — explicitly co-owned per `sr`'s AC-5. **This task's sole deliverable is the mod refit** (`docs/dev/_mods/pr-merge.md`). Because the mod is FO-read prose, this task cannot earn a *separate* non-grep AC over the mod (a substring sweep is the banned prose-grep; "the mod says run `dispatch trunk`" is exactly that tautology). So this task's behavioral proof IS `sr`'s AC-5: the base the refit mod uses is the stdout of `dispatch trunk`, and AC-5 proves that command returns the configured trunk (reds on a `next` regression). Command ownership confirmed with `sr` to stay on `sr`'s side (AC-5 as-written); `sr` reports no divergence to flag at the gate (its commit `94d89dd6`).

**The sole load-bearing contract between the two halves: `dispatch trunk` stdout is a BARE branch name.** `sr` pins AC-5 to **byte-exact stdout** — exactly the branch name plus a single trailing newline (`"main\n"`), all diagnostics to stderr; `sr`'s `TestDispatchTrunkCommand` asserts byte-exact stdout, not "contains". This matters because `$(...)` preserves *interior* newlines, so a stray log line on stdout would poison `$BASE`; the byte-exact pin forecloses that. This task's mod consumes it as `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` (the `$( )` strips the single trailing newline) and uses quoted `"$BASE"` at every site (`gh pr create --base "$BASE"`, the draft, the diff stat, the fallback merge target). The mod's capture is safe as written; the body assumes nothing else about the stdout shape — bare-name only. (Confirmed with team-lead at the staff review; byte-exact tightening confirmed with `sr`.)

## Proposed approach (firm)

1. **Refit the mod prose** (`docs/dev/_mods/pr-merge.md`, the split-root variant — NOT the repo-root `mods/pr-merge.md`, which is a separate non-split-root mod already on `main`) so the trunk is resolved by running the binary, not a literal:
   - The FO resolves the base once at the start of the merge hook: `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` (the configured trunk, default `main`; `dispatch trunk` emits a bare branch name, so `"$BASE"` is the clean branch ref).
   - All eight trunk-relative references become `"$BASE"` / `origin/"$BASE"` (rendered `origin/$BASE`) instead of the `next` literal — header description, split-root intro, `Branch: {branch} -> $BASE`, `git diff --stat origin/$BASE...{branch}`, the `gh pr create --base "$BASE"`, the rebase-target sentence, the push-remote note, and the local-merge fallback (`merge onto $BASE`, `merge-commit SHA on $BASE`).
   - Bump the mod `version:` (0.12.1 -> next patch) since the hook contract changed.
2. **Consume `sr`'s top-level `trunk: main` README declaration, `resolveTrunk` resolver, and `dispatch trunk` command.** `sr`'s task ships all three (its AC-5 + doc-diff); this task does not duplicate any of them. This task's mod points the FO at `dispatch trunk`.

The mod is no longer self-resolving: its base comes from a command that reads config, so a future trunk change is a one-line README edit, no mod edit and no per-merge override.

## Acceptance criteria

**AC-1 — The refit `pr-merge` mod's PR base resolves to the configured trunk (`main`), via `dispatch trunk`, never a hardcoded `next`.**
Verified by: `sr`'s **AC-5 command-level test** (`TestDispatchTrunkCommand`) — the **shared oracle** this pair binds to. The refit mod's base is `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)`; AC-5 drives that command against a fixture README:
  - fixture README with top-level `trunk: main` -> command prints `main`;
  - fixture README with **no** `trunk` key -> command prints `main` (default), proving absence does not fall through to `next`;
  - fixture README with top-level `trunk: ftrunk` (sentinel, neither `main` nor `next`) -> command prints `ftrunk`, proving the value is *sourced from the fixture config*, not a literal.
The expected value comes from the fixture README (an independent artifact the test writes), not from the mod prose or a code literal — the legitimate "relationship between independent values" the boundary guard permits, not a prose/code grep. **A `next`-defaulting or `next`-hardcoded resolver reds the absent-key case.** Because the mod's `gh pr create --base "$BASE"` takes that command's stdout verbatim, the command honoring the fixture trunk is the proof the FO opens the PR against the configured trunk. (The sentinel `ftrunk` matches `sr`'s reconcile fixture, so both consumers exercise the same `resolveTrunk` against the same sentinel.) This AC is co-owned: `sr` ships the command + test; this task's deliverable (the mod refit) consumes the command's output.

**AC-2 — `reconcile.go` classD/classE and the `pr-merge` base read the SAME resolver (one source, not two).**
Verified by: `resolveTrunk` is a single function reading the top-level `trunk:` key via `ParseFrontmatter`; `sr`'s reconcile test (fixture `trunk: ftrunk` drives classD/classE detection against `origin/ftrunk`) and `sr`'s AC-5 `dispatch trunk` test both exercise the *same* `resolveTrunk` symbol against a fixture README declaring the top-level key. A `grep -c` over the codebase is **not** the proof — the proof is that both tests pass against one symbol (if either consumer forked a second source, one test would read a value the other never wrote). The mod's base flows through that command, so the mod binds to the same single resolver as reconcile.

**Note on what is NOT an AC:** "the mod text no longer contains `next`" is real authoring work (it is in the doc diff and the stage-report checklist) but is **not** an acceptance criterion — a substring sweep over `pr-merge.md` is the banned prose-grep (`internal/contractlint/boundary_guard_test.go`). Equally, "the mod says run `dispatch trunk`" is a prose-grep and is NOT an AC. The behavior that matters — the FO opens the PR against the configured trunk — is proven only by `sr`'s AC-5 command-level test (the shared oracle), never by reading the mod. This task is a prose refit; its behavioral backing is the shared command, by design (a doc-only refit cannot carry an independent non-grep AC, and manufacturing one over the mod text would be the exact tautology the guard bans).

## Test plan

| What | Where | Owner | Cost |
|---|---|---|---|
| `dispatch trunk` resolves `main` / default / sentinel `ftrunk` (AC-1 shared oracle) | `sr`'s `TestDispatchTrunkCommand` (behavior fixture driving `dispatch trunk`) | `sr` (AC-5); this task consumes its output | low |
| `resolveTrunk`: top-level `trunk: X` -> X, absent -> `main` | `sr`'s resolver unit test (`sr` names it) | `sr` | low |
| One-source check (AC-2) | both `sr`'s reconcile test and `sr`'s `dispatch trunk` test pass against the single `resolveTrunk` symbol | shared | low (assertion: "both green, one symbol") |
| Mod refit (this task's deliverable) | `docs/dev/_mods/pr-merge.md` — 8 trunk-relative `next` literals -> `{base}` = `dispatch trunk` output | **this task** | low (prose diff; no separate test — proof is AC-1's shared command oracle) |

No live workflow / CLI smoke test needed: the claim is command-level (does the command print the configured trunk) and unit-level (does the resolver read config), both in `sr`'s tests. The end-to-end "FO opens PR against `main`" is a prose instruction whose only mechanical part — resolving the base — is the shared command; the `gh pr create` call itself is not exercised in CI (outward-facing, gated on captain approval). This task adds no new test of its own: a doc-only mod refit's behavioral backing is the shared `dispatch trunk` oracle by design; a test over the mod text would be the banned prose-grep.

## Spike result (riskiest unknown, exercised first)

The design rested on one unverified mechanism, exercised in two parts before the design firmed (throwaway tests, removed after):

**(1) Does `status.ParseFrontmatter` surface a top-level `trunk:` key, and does absence yield empty (so the resolver owns the `main` default)?** This is the parser `sr`'s `resolveTrunk` reads (on the README); `internal/status` is already imported in reconcile (`reconcile.go:17`), so the new resolver adds no import, only a new `ParseFrontmatter`-on-README call site.

```
=== RUN   TestSpikeTrunkFrontmatter
    SPIKE OK: trunk="main" state=".spacedock-state"; absent-trunk yields empty
--- PASS
```

**(2) Does the key-placement choice hold — is a NESTED `stages.defaults.trunk` invisible to `ParseFrontmatter` (the decisive reason for top-level)?**

```
=== RUN   TestSpikeTrunkPlacement
    SPIKE3 OK: top-level trunk visible to ParseFrontmatter; nested INVISIBLE (sr's claim confirmed)
--- PASS
```

Result: top-level `trunk: main` → `ParseFrontmatter["trunk"]=="main"`; absent → `""` (resolver owns the `main` default; fallback reachable, not dead code). A nested `stages.defaults.trunk` → `""` via `ParseFrontmatter` (nested mappings render empty, `frontmatter.go:123-126`) — **confirming `sr`'s decisive rationale** that a nested key would force a second parser. **Confirmed**: the canonical source is readable with zero parser changes via `ParseFrontmatter` (already imported in reconcile at `:17`; the resolver adds a new one-line call site, not a new import), and top-level placement is the correct choice. (An intermediate spike against `ParseStagesWithDefaults` for a nested key — during a brief mid-reconciliation detour — is superseded by this convergence.) The `dispatch trunk` command is a ~12-line `case "trunk":` on the existing `dispatch.Run` switch (`internal/dispatch/dispatch.go:80`), host-neutral (README parse only, no `~/.claude` coupling) — proven feasible by the existing `case "reconcile"` sibling that drives the same fixture-README harness.

## Doc diff (mod refit — applied at implementation)

Concrete before/after for `docs/dev/_mods/pr-merge.md`. The FO resolves the base once near the top of the merge hook; all literals become `{base}`/`origin/{base}`.

- **Frontmatter `version:`** `0.12.1` -> bump one patch (hook contract changed).
- **Frontmatter `description:`** "Open a code-branch PR to next at the merge boundary…" -> "Open a code-branch PR to the configured trunk at the merge boundary…"
- **Intro (line 9):** "base branch `next`" -> "base branch resolved from the workflow's `trunk:` config (default `main`)".
- **Merge hook, new opening step:** add "Resolve the PR base once: `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` — the workflow's configured integration trunk (default `main`). `dispatch trunk` emits exactly a **bare branch name** (e.g. `main`), so `$( )` yields `$BASE` clean (command substitution strips the single trailing newline). Always quote `"$BASE"` at use sites. Use `"$BASE"` for the draft, the diff stat, the `gh pr create --base`, and the fallback merge target below."
- **Draft `Branch:` line (45):** `{branch} -> next` -> `{branch} -> {base}`.
- **Draft `Changes:` line (46):** `git -C {worktree} diff --stat origin/next...{branch}` -> `... origin/{base}...{branch}`.
- **Approval rebase sentence (56):** "rebases the code branch onto `origin/next`" -> "onto `origin/{base}`".
- **Push note (58):** "(`origin` = the main repo, base branch `next`)" -> "(`origin` = the main repo, base branch `{base}`)"; "Do NOT push `next`" -> "Do NOT push the trunk".
- **`gh pr create` (62):** `--base next` -> `--base "$BASE"`.
- **Fallback merge target (100):** "a local `--no-ff` merge … onto `next`" -> "onto `{base}`".
- **Fallback sentinel (103):** "compute the merge-commit SHA on `next`" -> "on `{base}`"; "The SHA must already exist on `next`" -> "on `{base}`".

No other doc files reference the **split-root** pr-merge base. The dev `README.md` gains a top-level `trunk: main` key — that diff is `sr`'s canonical-source change (in `sr`'s doc-diff already); this task does not duplicate it. The repo-root `mods/pr-merge.md` is a separate non-split-root mod already on `main` and is out of this task's scope.

## Notes
Mod files are refit-or-worker territory, not FO direct edits. Coupled pair with `sr` (`dispatch-reconcile-deconflate-repo-hygiene`): `sr` OWNS the single source (top-level `trunk:` in `docs/dev/README.md`), the `resolveTrunk` resolver, and the `dispatch trunk` command (its AC-5); this task CONSUMES all three — its sole deliverable is the split-root mod refit pointing the FO at `dispatch trunk`. Anti-conflation guard carried: base reads the top-level `trunk:` key, never `cli.devBranch`. High-stakes (merge/release machinery) → expect a staff review before the gate. A parallel 5-surface trunk-config audit cross-checks the `next`-as-trunk inventory at the gate.

## Stage Report: ideation

- DONE: The design refits the `pr-merge` mod's PR base branch from hardcoded `next` to `main`, read from the SAME single canonical trunk-config source `sr` defines — config-driven, not a per-merge manual override.
  Firm design + doc diff refit all 8 trunk-relative `next` literals in the split-root `docs/dev/_mods/pr-merge.md` to `{base}`, resolved by `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` reading the settled top-level `trunk:` README key via `sr`'s `resolveTrunk` resolver. Converged with `sr` over three SendMessage rounds to ONE source (latest `sr` commit `94d89dd6`); both halves now agree on top-level `trunk:`.
- DONE: An acceptance criterion proven by a real check (a test/fixture that reds if the mod opens a PR against `next` or resolves the base to `next`) — never a prose/string grep over the mod text.
  AC-1 binds to `sr`'s AC-5 command-level test (`TestDispatchTrunkCommand`, the shared oracle): a fixture README's top-level `trunk:` drives `dispatch trunk` stdout; the absent-key case reds a `next`-defaulting resolver, the sentinel `ftrunk` case proves config-sourcing (independent oracle = fixture config, not mod prose). The refit mod's `--base "$BASE"` consumes that stdout, so the command's correctness IS the mod's base correctness. Explicit non-AC note records that a mod grep (incl. "the mod says run `dispatch trunk`") is the banned prose-grep (`boundary_guard_test.go`); a doc-only refit carries no independent non-grep AC by design.
- DONE: The single trunk-config source is consistent with `sr`'s design (the coupled pair settles ONE source) — exercised that the `pr-merge` mod actually reads it.
  One `resolveTrunk` (sr's) over the top-level `trunk:` key via `ParseFrontmatter`; AC-2 proves one-source by `sr`'s reconcile test and `sr`'s `dispatch trunk` test passing against the same `resolveTrunk` symbol on the shared `ftrunk` sentinel — the mod's base flows through that command. Two de-risk spikes (throwaway, removed): (1) `ParseFrontmatter` surfaces top-level `trunk: main` and yields empty when absent; (2) a nested `stages.defaults.trunk` is INVISIBLE to `ParseFrontmatter`, confirming `sr`'s decisive rationale for top-level placement — both exercised before the design firmed.

### Summary

Refit designed and converged with `sr` to ONE source: the split-root `pr-merge` mod stops hardcoding `next` and resolves its PR base by running `spacedock dispatch trunk`, a thin command wrapping `sr`'s `resolveTrunk`, which reads the single canonical top-level `trunk:` key in the workflow README (default `main`) — the same source+symbol `sr` consumes for `reconcile.go` classD/classE. Key-placement (top-level vs nested) was settled on `sr`'s decisive parser argument, which I verified by spike: a nested key is invisible to `ParseFrontmatter`, forcing a second parser, so top-level keeps ONE parser behind ONE resolver. `sr` owns the command + resolver + README key + their tests (its AC-5 is the shared oracle); this task's sole deliverable is the mod refit, whose behavioral proof is that shared oracle — a doc-only refit cannot carry an independent non-grep AC, which is surfaced honestly rather than papered over with a mod grep. Anti-conflation guard carried: the base never routes through `cli.devBranch`. Division of labor confirmed with `sr` both ways — command ownership stays with `sr` (AC-5); `sr` reports no divergence to flag at the gate. The pair has settled ONE source.

**Convergence settled (crossed-wires resolved):** an earlier message round briefly had me reconcile to `sr`'s superseded `stages.defaults.trunk`; both halves are now provably on **top-level `trunk:`** (the placement I first proposed, which `sr` adopted on the merits — altitude: trunk is a workflow-level repo property, sibling of `state:`, not a stage default). `sr` resolved my open question on the resolver's home: **`resolveTrunk` lives unexported in `internal/dispatch`** (both callers — `reconcile.go` and the `dispatch trunk` command — are in that package, so no export). I confirmed-back with a fresh top-level-`trunk:`-via-`ParseFrontmatter` spike (throwaway, removed): `main`/`ftrunk`/absent→empty, all green. I never wired any code against an `internal/status` home (ideation is non-worktree), so nothing to unwind there.

## Stage Report: implementation

- DONE: Refit `docs/dev/_mods/pr-merge.md` (split-root variant): replace all 8 trunk-relative `next` literals with `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` / `origin/$BASE` / `"$BASE"`; bump version one patch.
  Commit 6f2f97b6 on worktree branch. New opening resolve step + description + intro + Branch/Changes draft + rebase sentence + push note + `gh pr create --base "$BASE"` + fallback merge target + sentinel. Version 0.12.1 -> 0.12.2. `grep -n next pr-merge.md` returns nothing.
- DONE: Behavioral backing rides sr's AC-5 shared oracle: `go test ./internal/dispatch -run TestDispatchTrunkCommand` green on the worktree.
  `ok github.com/spacedock-dev/spacedock/internal/dispatch`. Also exercised the live command: `dispatch trunk --workflow-dir docs/dev` emits byte-exact `main\n` (xxd: `6d61 696e 0a`), so the mod's `BASE=$(...)` capture yields a clean bare branch name.

### Summary

Refit the split-root pr-merge mod so its PR base is resolved by running `spacedock dispatch trunk` (sr's shared command reading the canonical top-level `trunk:` README key, default `main`) instead of the hardcoded pre-flip `next`. All 8 trunk-relative literals now use `"$BASE"` / `origin/$BASE`. Base reads the top-level `trunk:` key only, never `cli.devBranch` (anti-conflation guard honored). Proof is the shared AC-5 oracle (green) plus a live byte-exact `main\n` check; no separate mod prose-grep AC by design.
