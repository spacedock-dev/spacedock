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

## Coupling with `sr` — ONE canonical trunk-config source

`sr` LEADS the canonical-source design; this task CONSUMES it for the PR base. I do **not** invent a second source.

Settled jointly (my proposal, sent to `sr`; `sr`'s AC-1 sketch independently names "workflow/README config" as the trunk source, so we are aligned on the location):

- **Canonical source:** a `trunk:` key in the workflow README frontmatter (`docs/dev/README.md`), parsed by the binary exactly as `state:` is — `status.ParseFrontmatter` returns `map[string]string`, so `fm["trunk"]` is free. Default `main` when the key is absent (old READMEs).
- **One resolver in code:** a `TrunkBranch(workflowDir)` resolver in `internal/status` (or `internal/dispatch`) that reads the README's `trunk:` and applies the `main` default. `sr`'s `reconcile.go` classD/classE call it instead of the `origin/next` literal; this task's command surface (below) calls the same resolver.

**ALIGNMENT FLAG for the gate (open, defer to `sr`):** whether the resolver is exposed as a thin command (`spacedock dispatch trunk --workflow-dir docs/dev`, prints the resolved trunk) in addition to the in-code function. This task's real (non-grep) AC *requires* the command: the mod is prose, so the FO can only become config-driven by *running* the binary to get the base, not by reading a literal. `sr` needs the resolver function regardless (for the reconcile test); the command is a ~12-line wrapper on the existing `dispatch.Run` switch (`internal/dispatch/dispatch.go:26`). If `sr` keeps it function-only, this task either (a) adds the `dispatch trunk` command itself as the consumer-side surface, or (b) cannot meet the non-grep AC and must be re-scoped — surfaced here for the staff reviewer / gate to decide. Proceeding on the assumption the command lands (cheapest path; benefits both halves).

## Proposed approach (firm)

1. **Refit the mod prose** (`docs/dev/_mods/pr-merge.md`) so the trunk is resolved by running the binary, not a literal:
   - The FO resolves the base once at the start of the merge hook: `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` (the standing trunk, default `main`).
   - All eight trunk-relative references become `{base}` / `origin/{base}` instead of the `next` literal — header description, split-root intro, `Branch: {branch} -> {base}`, `git diff --stat origin/{base}...{branch}`, the `gh pr create --base {base}`, the rebase-target sentence, the push-remote note, and the local-merge fallback (`merge onto {base}`, `merge-commit SHA on {base}`).
   - Bump the mod `version:` (0.12.1 -> next patch) since the hook contract changed.
2. **Add (or confirm `sr` adds) the `trunk:` README key** = `main` and the `TrunkBranch` resolver. This is the single source both halves consume.
3. **Expose `spacedock dispatch trunk --workflow-dir <dir>`** that prints the resolved trunk to stdout (calls `TrunkBranch`). This is the FO-runnable oracle the refit mod points at.

The mod is no longer self-resolving: its base comes from a command that reads config, so a future trunk change is a one-line README edit, no mod edit and no per-merge override.

## Acceptance criteria

**AC-1 — `spacedock dispatch trunk` resolves the base from the README `trunk:` key, defaulting to `main`, never `next`.**
Verified by: a Go command-level test in `internal/dispatch` (mirroring `reconcile_test.go`'s fixture+`Run`+stdout pattern) that writes a fixture README and drives the command:
  - fixture README with `trunk: main` -> command prints `main`;
  - fixture README with **no** `trunk:` key -> command prints `main` (default), proving absence does not fall through to `next`;
  - fixture README with `trunk: release-3.0` (a value that is neither `main` nor `next`) -> command prints `release-3.0`, proving the value is *sourced from the fixture config*, not a literal in the code under test.
The expected value comes from the fixture README (an independent artifact the test writes), not from the mod prose or a code literal — this is the legitimate "relationship between independent values" the boundary guard permits, not a prose/code grep. **A `next`-defaulting or `next`-hardcoded resolver reds the second case.** (Independent oracle: the fixture's configured trunk drives the output.)

**AC-2 — the resolver is config-sourced, not a literal: a fixture trunk other than `next`/`main` flows end-to-end to the PR base the FO would use.**
Verified by: the AC-1 third case (`trunk: release-3.0`) is itself the proof for the consumer path — the command the refit mod invokes returns the fixture's trunk verbatim. Because the mod's `gh pr create --base "$BASE"` takes the command's stdout, a test proving the command honors the fixture trunk proves the FO opens the PR against that trunk. No grep over the mod text is involved; the oracle is the fixture config -> command stdout relationship.

**AC-3 — `reconcile.go` classD/classE and the `pr-merge` base read the SAME resolver (one source, not two).**
Verified by: the `TrunkBranch` resolver is a single exported function; `sr`'s reconcile test (fixture trunk drives classD/classE) and this task's `dispatch trunk` test both exercise it against a fixture README. A `grep -c` over the codebase is **not** the proof — the proof is that both tests pass against the *same* resolver symbol (if either half forked a second source, one test would be reading a value the other never wrote). Reconciled with `sr` at the joint gate: if `sr`'s landed source differs from README `trunk:`, this AC's test target moves to whatever single symbol `sr` shipped — there must be exactly one.

**Note on what is NOT an AC:** "the mod text no longer contains `next`" is real authoring work (it is in the doc diff and the stage-report checklist) but is **not** an acceptance criterion — a substring sweep over `pr-merge.md` is the banned prose-grep (`internal/contractlint/boundary_guard_test.go`). The behavior that matters — the FO opens the PR against the configured trunk — is proven only by the command-level test (AC-1/AC-2), never by reading the mod.

## Test plan

| What | Where | Cost | Type |
|---|---|---|---|
| `dispatch trunk` resolves `main` / default / arbitrary fixture trunk | `internal/dispatch/trunk_test.go` (new), fixture README + `dispatch.Run([]string{"trunk","--workflow-dir",dir}, ...)` asserting stdout | low (~1 file, mirrors `reconcile_test.go:432` `TestReconcileIncludeScope` harness) | Go command-level fixture test |
| `TrunkBranch` resolver: `trunk: X` -> X, absent -> `main` | `internal/status` (or wherever `sr` lands it) unit test; seeded by the throwaway spike below | low | Go unit test |
| One-source check (AC-3) | both `sr`'s reconcile test and this command test pass against the single `TrunkBranch` symbol | low (no new test; assertion is "both green, one symbol") | reconciliation at gate |

No live workflow / CLI smoke test needed: the claim is command-level (does the command print the configured trunk) and unit-level (does the resolver read config). The end-to-end "FO opens PR against `main`" is a prose instruction whose only mechanical part — resolving the base — is the command under test; the `gh pr create` call itself is not exercised in CI (it is outward-facing and gated on captain approval).

## Spike result (riskiest unknown, exercised first)

The design rests on one unverified mechanism: **does `status.ParseFrontmatter` surface a `trunk:` key from a README the way it surfaces `state:`, and does an absent key yield empty (so the resolver must default)?** Exercised before committing the design (throwaway test, removed after):

```
=== RUN   TestSpikeTrunkFrontmatter
    SPIKE OK: trunk="main" state=".spacedock-state"; absent-trunk yields empty
--- PASS
```

Result: `ParseFrontmatter` returns `fm["trunk"]="main"` from a fixture README carrying `trunk: main`, and `fm["trunk"]==""` when the key is absent. **Confirmed**: the canonical source is readable with zero parser changes, and the resolver owns the `main` default. This seeds the implementation's first test (the resolver default case). The `dispatch trunk` command is a ~12-line `case "trunk":` on the existing `dispatch.Run` switch (`internal/dispatch/dispatch.go:80`), host-neutral (README parse only, no `~/.claude` coupling) — proven feasible by the existing `case "reconcile"` sibling that drives the same fixture-README harness.

## Doc diff (mod refit — applied at implementation)

Concrete before/after for `docs/dev/_mods/pr-merge.md`. The FO resolves the base once near the top of the merge hook; all literals become `{base}`/`origin/{base}`.

- **Frontmatter `version:`** `0.12.1` -> bump one patch (hook contract changed).
- **Frontmatter `description:`** "Open a code-branch PR to next at the merge boundary…" -> "Open a code-branch PR to the configured trunk at the merge boundary…"
- **Intro (line 9):** "base branch `next`" -> "base branch resolved from the workflow's `trunk:` config (default `main`)".
- **Merge hook, new opening step:** add "Resolve the PR base once: `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` — the workflow's configured integration trunk (default `main`). Use `$BASE` for the draft, the diff stat, the `gh pr create --base`, and the fallback merge target below."
- **Draft `Branch:` line (45):** `{branch} -> next` -> `{branch} -> {base}`.
- **Draft `Changes:` line (46):** `git -C {worktree} diff --stat origin/next...{branch}` -> `... origin/{base}...{branch}`.
- **Approval rebase sentence (56):** "rebases the code branch onto `origin/next`" -> "onto `origin/{base}`".
- **Push note (58):** "(`origin` = the main repo, base branch `next`)" -> "(`origin` = the main repo, base branch `{base}`)"; "Do NOT push `next`" -> "Do NOT push the trunk".
- **`gh pr create` (62):** `--base next` -> `--base "$BASE"`.
- **Fallback merge target (100):** "a local `--no-ff` merge … onto `next`" -> "onto `{base}`".
- **Fallback sentinel (103):** "compute the merge-commit SHA on `next`" -> "on `{base}`"; "The SHA must already exist on `next`" -> "on `{base}`".

No other doc files reference the pr-merge base. (The dev `README.md` gains the `trunk: main` key — that diff is `sr`'s canonical-source change or this task's if `sr` defers; tracked under AC-3's one-source reconciliation.)

## Notes
Mod files are refit-or-worker territory, not FO direct edits. Coupled pair with `sr` (`dispatch-reconcile-deconflate-repo-hygiene`) — settle ONE source. High-stakes (merge/release machinery) → expect a staff review before the gate.

## Stage Report: ideation

- DONE: The design refits the `pr-merge` mod's PR base branch from hardcoded `next` to `main`, read from the SAME single canonical trunk-config source `sr` defines — config-driven, not a per-merge manual override.
  Firm design + doc diff in the task body refit all 8 trunk-relative `next` literals in `docs/dev/_mods/pr-merge.md` to `{base}`, resolved by `spacedock dispatch trunk` reading the README `trunk:` key (the source `sr`'s AC-1 sketch independently names); aligned with `sr` via SendMessage.
- DONE: An acceptance criterion proven by a real check (a test/fixture that reds if the mod opens a PR against `next` or resolves the base to `next`) — never a prose/string grep over the mod text.
  AC-1/AC-2: a command-level Go test (mirrors `reconcile_test.go:432` fixture+`Run`+stdout) where a fixture README's `trunk:` drives `dispatch trunk` stdout; the absent-key case reds a `next`-defaulting resolver, the `trunk: release-3.0` case proves config-sourcing. Independent oracle = fixture config, not mod prose. Explicit non-AC note records the banned prose-grep (`boundary_guard_test.go`).
- DONE: The single trunk-config source is consistent with `sr`'s design (the coupled pair settles ONE source) — exercised that the `pr-merge` mod actually reads it.
  One `TrunkBranch` resolver over README `trunk:`; AC-3 proves one-source by both `sr`'s reconcile test and this command test passing against the same symbol. De-risk spike (throwaway, removed) confirmed `ParseFrontmatter` surfaces `trunk: main` and yields empty when absent — the mechanism the mod's base-resolution rides on, exercised before committing the design.

### Summary

Refit designed: the `pr-merge` mod stops hardcoding `next` and instead resolves its PR base by running `spacedock dispatch trunk`, which reads a single canonical `trunk:` key in the workflow README (default `main`) — the same source `sr` consumes for `reconcile.go` classD/classE. The riskiest unknown (can the binary read a `trunk:` README key) was spiked green before the design firmed. One open item flagged for the gate: whether `sr` exposes the resolver as the `dispatch trunk` command this task's non-grep AC depends on, or whether this task adds that ~12-line command itself — defer to `sr`'s canonical-source choice at the joint gate.
