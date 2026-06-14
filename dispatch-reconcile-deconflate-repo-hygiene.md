---
id: sryzghzqazj9s9km6ebqkf5s
title: dispatch reconcile conflates team-management with repo-hygiene (and hardcodes pre-flip trunk `next`)
status: validation
source: "0202 Commander drive (2026-06-13). Boot reconcile flagged Class-D/E drift against origin/next; Class-E remedy 'reset main->origin/next' would have reverted the entire post-flip trunk. Investigation (captain-prompted) found the deeper cause: reconcile bundles git-hygiene into a team-management helper, so it carries repo/trunk knowledge it shouldn't."
group: cleanup
sprint: 0203-fo-efficiency
started: 2026-06-14T05:04:07Z
worktree: .worktrees/spacedock-ensign-dispatch-reconcile-deconflate-repo-hygiene
mod-block: merge:pr-merge
---

`dispatch reconcile` is two helpers in one coat: roster/team reconciliation AND repo git-hygiene. The repo half hardcodes the pre-flip trunk `next`, which the 2026-06-08 flip silently invalidated.

## Problem

`internal/dispatch/reconcile.go` emits five drift classes:
- **A** (lingering agent), **B** (superseded agent) — genuine team management, sourced from the `~/.claude/teams` roster.
- **C** (un-advanced PR) — entity/PR state.
- **D** (stale branch), **E** (stale local main) — pure git hygiene, hardcoding `origin/next` (`reconcile.go:582`, `:605`, `:616`).

reconcile landed 2026-06-02 (#273), six days BEFORE the 2026-06-08 flip — when `next` was genuinely the integration trunk and `main` tracked it, so D/E were correct then. The flip inverted the model (`main` is now the trunk; `next` is dev-only per `docs/releasing.md`), but the helper was never refit. Post-flip, Class-E's remedy ("reset main->origin/next") would throw away the real trunk and revert to the stale dev branch — the dangerous drift the 0202 Commander refused to act on at boot.

Root cause is the conflation: a helper whose name + roster-loading say "team management" should not carry repo/trunk knowledge. The hardcoded `next` is the symptom; the bundling is the cause.

## Proposed approach (firm)

Two moves, both required — they answer the captain's twin DoD ("de-conflate" AND "settle ONE trunk-config source"):

### 1. The single canonical trunk-config source (this seed OWNS the design; `87` consumes it)

The integration trunk is declared **once**, as a **top-level `trunk:` key** in the workflow README frontmatter (`docs/dev/README.md`). For this repo: `trunk: main`. Top-level — a sibling of `state:` — NOT nested under `stages.defaults` (key-placement settled with `87`; rationale below).

One resolver, two consumers, one command surface:
- **`resolveTrunk(workflowDir)`** is the ONE resolver, living in **`internal/dispatch`** (where both callers already are — `reconcile.go` and the `dispatch trunk` command — so it needs no export and there is exactly one symbol). It reads `status.ParseFrontmatter(README)["trunk"]` and returns `main` when the key is empty/absent.
- **reconcile** (Go) calls `resolveTrunk` directly — `classD`/`classE` consume its return.
- **`87`'s pr-merge mod** is PROSE the FO reads; it cannot call a Go function and cannot earn a non-grep AC if it embeds a base-branch literal. So the resolver is also surfaced as a thin command — **`spacedock dispatch trunk --workflow-dir docs/dev`** prints the resolved trunk to stdout. `87`'s mod becomes "PR base = output of `dispatch trunk`", and `87`'s AC is a behavior test on that COMMAND (a fixture README's declared trunk drives the printed value), not a grep over the mod. The command and reconcile call the SAME `resolveTrunk` resolver, so the two consumers cannot diverge. (Adopted from `87`'s coordination ask — the command is the only thing that gives the prose half a real oracle.)

Why top-level `trunk:` (not nested under `stages.defaults`):
- **Altitude:** the trunk is a workflow-level repo property, a sibling of `state:`, not a stage default like `worktree:`/`concurrency:`.
- **One parser serves both consumers:** `status.ParseFrontmatter` surfaces top-level scalar keys as a flat `map[string]string`, so `fm["trunk"]` is free for the `dispatch trunk` command AND reconcile. Critically, `ParseFrontmatter` renders NESTED mappings as empty string (frontmatter.go:123-126) — so a `stages.defaults.trunk` key would be **invisible** to the simple parser the command uses, forcing the heavier `ParseStagesWithDefaults` and a second read path. Top-level keeps exactly one parser behind one resolver. `internal/dispatch` already imports `internal/status` (reconcile.go:17) and already calls `ParseFrontmatter` elsewhere (reconcile.go:296, on per-ENTITY files), so no new import — but note `resolveTrunk`'s `ParseFrontmatter`-on-the-README call site is a NEW one-line call (the README is currently read only via `ParseStagesWithDefaults` at reconcile.go:319, not via `ParseFrontmatter`).

Fallback when the key is absent: `resolveTrunk` returns `main` (the current post-flip trunk), NOT `next`. The hardcoded `next` is deleted outright, not demoted to a fallback. Absence-of-key is a reachable path (spike-confirmed) so the fallback is real, not dead code.

**NOT this source — the explicit anti-conflation guard:** `internal/cli/frontdoor.go`'s `var devBranch = "next"` is the *marketplace channel stamp* (`cli.devBranch=main` → stable entry, `=next` → edge entry, per `docs/releasing.md:13`). It is semantically a different axis from the integration trunk and `next` is *correct* there (the edge channel really does track `next`). This design MUST NOT route the trunk through `devBranch` or vice-versa. Conflating the channel stamp with the trunk would re-create the exact bundling bug one layer over.

### 2. De-conflate: trunk knowledge leaves the roster-management helper's body

`dispatch reconcile` keeps emitting all five classes (the FO's single boot/idle sweep stays one call — no behavior split for the operator), but the **trunk-resolution logic is extracted out of the class detectors**:

- The `resolveTrunk(workflowDir)` resolver (above) reads the top-level `trunk:` key, returning `main` on absent/empty. It is resolved **once** in `Reconcile()` and passed into `classD`/`classE` as a parameter.
- `classD`/`classE` take the resolved trunk as an argument instead of embedding `origin/next` literals. The git refs become `HEAD..origin/{trunk}` and `origin/{trunk}..main`; the `Reason` strings interpolate `{trunk}`.
- The roster-derived detectors (`classA`/`classB` from `~/.claude/teams`) are untouched — they already carry zero repo/trunk knowledge. The conflation was D/E embedding a trunk literal; extracting the literal into a parameter sourced outside the detector is the de-conflation. Trunk knowledge now lives in ONE resolver reading ONE declared source, not scattered as literals inside roster-helper detectors.

**The FO-side remedy follows the data, not prose:** `classE`'s drift item gains a `trunk` field (the resolved trunk). The FO runtime's E-remedy reads `origin/{drift.trunk}` from the JSON rather than a hardcoded `git reset --hard origin/next`. This keeps the trunk out of the instruction-file prose too — the remedy becomes data-driven. (The doc-diff for the FO runtime is below; the *behavioral* proof is the JSON field, not the prose.)

### Spike result (riskiest unknown — exercised first)

**Unknown:** where the single trunk source lives AND that one parser can serve both the `dispatch trunk` command and reconcile. **Exercised** with throwaway tests against the real parsers:
- **Chosen placement** — top-level `trunk: main` → `status.ParseFrontmatter(README)["trunk"] == "main"`. PASS. Absent key → `fm["trunk"] == ""` (resolver defaults to `main`; fallback reachable, no panic). PASS.
- **Rejected placement (also exercised, to justify the choice)** — nested `stages.defaults.trunk: main` is reachable only via `ParseStagesWithDefaults` (`defaults["trunk"] == "main"`); it is INVISIBLE to `ParseFrontmatter` (nested mappings render empty). Confirming the heavier parser would be forced for a nested key is what settled the top-level choice.

Conclusion: the source is a top-level `trunk:` key, read via `status.ParseFrontmatter` — a parser `internal/dispatch` **already imports** (reconcile.go:17) and already calls on entity files (reconcile.go:296), and that the `dispatch trunk` command can call on the README identically. (`resolveTrunk`'s call of `ParseFrontmatter` on the README is a new one-line call site; the README itself is currently read only via `ParseStagesWithDefaults` at reconcile.go:319.) The design composes one proven mechanism (`ParseFrontmatter` returning a flat map) plus a one-line lookup behind a single `resolveTrunk` resolver. No on-disk format, no new handoff — the only unverified bit was the round-trip for a new key, now verified for both placements.

## Acceptance criteria

**AC-1 — Class-D/E resolve the integration trunk from the configured top-level `trunk:` key, never a hardcoded `next`.**
Verified by: a reconcile fixture test (extending `reconcile_test.go`'s real-git-tree harness) where the fixture README declares top-level `trunk: ftrunk` (a sentinel name that is neither `next` nor `main`) and the fixture git graph builds `origin/ftrunk` as the trunk ref. The test asserts Class-D fires when a worktree branch is behind `origin/ftrunk` and Class-E fires when local `main` carries commits not on `origin/ftrunk`. **The expected ref name comes from the fixture's declared trunk, not from any literal in reconcile.go** — so a regression that re-hardcodes `next` (or `main`) reds the test (it would query `origin/next`, find no such ref / wrong count, and fail to detect against `origin/ftrunk`). This is an independent oracle: the fixture trunk and the code's resolved trunk can disagree, which is what makes it able to fail. NOT a prose/string grep.

**AC-2 — The trunk literal is absent from `classD`/`classE`; trunk knowledge is sourced through one `resolveTrunk` resolver.**
Verified by: a Go test asserting `resolveTrunk` returns the declared trunk for a README with top-level `trunk:` set and returns `main` (the fallback) for a README with no trunk key — expected values from the fixture READMEs, not from the function's own source. Combined with AC-1's behavioral test, this proves D/E read the resolved value rather than a literal: the binding proof is AC-1's behavioral red-on-regression (a re-hardcoded `next` queries `origin/next` and fails to detect against the fixture's `origin/ftrunk`).

**AC-3 — Class-E's drift item carries the resolved `trunk`, so the FO remedy is data-driven, not prose-hardcoded.**
Verified by: the AC-1 fixture test also asserts the emitted Class-E `driftItem` has `trunk == "ftrunk"` (from the fixture's declared trunk). Expected value from the fixture config, not the code. This is the behavioral proof that the FO can read the trunk from JSON instead of a hardcoded remedy string; the FO-runtime prose change is doc-only and is NOT itself an acceptance criterion (a prose-grep over the runtime file is banned).

**AC-4 — The marketplace channel stamp stays un-conflated with the trunk.**
Verified by: a test (or assertion within the existing frontdoor tests) that `resolveTrunk` does not read `cli.devBranch` and `devBranch`'s value is unchanged by trunk resolution — i.e. resolving the trunk on a workflow whose README declares `trunk: main` leaves `devBranch == "next"`. Expected values are the two independent declared sources. This guards against re-bundling the two axes. (Cheap: it is an invariant over two already-separate variables.)

**AC-5 — `spacedock dispatch trunk --workflow-dir DIR` prints the resolved trunk as EXACT stdout (a bare branch name, single trailing newline, nothing else); this is the shared oracle the prose consumers bind to.**
Verified by: a command-level Go test (behavior fixture driving the `dispatch trunk` subcommand) where a fixture README declares top-level `trunk: ftrunk` and the command's stdout is asserted **byte-exact** — `"ftrunk\n"`, not "contains ftrunk"; a second fixture with no trunk key asserts stdout byte-exact `"main\n"`. **Expected value from the fixture README, not from the command's source.** Exact stdout is the SOLE load-bearing contract between this seed and the prose consumers: both `87`'s pr-merge mod and the shared-core ship-local ceremony (AC-7) capture the base with `BASE=$(spacedock dispatch trunk --workflow-dir docs/dev)` and pass `--base "$BASE"` / merge onto `$BASE` verbatim. A stray log line or a second output line would poison `$BASE` (shell `$(...)` keeps interior newlines), so the test asserts the command writes ONLY the bare trunk to stdout (any diagnostics go to stderr). This is the command surface that lets the prose consumers be config-driven with a real, fail-able oracle instead of a banned prose-grep. (Co-owned: `sr` ships the command + this test; the prose consumers' own ACs test their text consuming the command's output.)

**AC-6 — No `.github/workflows/*.yml` resolves the integration trunk to `next` post-flip.**
Verified by: a Go structural lint (new test in `internal/contractlint`, alongside `structural_checks_test.go`) that checks the **two specific trunk-resolving positions** the audit found, NOT a token sweep:
  (a) **`on.pull_request.branches`** — YAML-parse each workflow and assert the PR-base-filter sequence does not list `next` (spike note: handle the YAML-1.1 `on:`→`true` bare-key quirk).
  (b) **`gh run list … --branch next`** — a precise regex (`gh run list[^\n]*--branch\s+next\b`) over each file's body asserts no producer-run query targets `next`.
Checking these two positions (rather than "any `next` token minus a blocklist") is what makes the oracle precise: it does NOT false-positive on the legitimately-`next` uses — `next-publish.yml`'s edge-channel publish, the `awk … next` keyword in release.yml, an `@next` action pin, `--ref next`, or English in comments — because none of those occupy position (a) or (b). CI YAML is parsed by GitHub Actions, NOT ingested by the model as instruction prose, so a structural assertion over it is a legitimate non-grep oracle (outside "files the model reads" — the same exemption `structural_checks_test.go` relies on). It is able to fail: today it reds on exactly the 3 pre-flip offenders (spike-confirmed), and reds again on any future re-introduction. STANDING guard, not a one-time flip.

**Spike (AC-6 mechanism — exercised before committing):** ran the two-position lint against the live `.github/workflows/`. Result: exactly 3 offenders flagged (`install-e2e.yml` + `runtime-live-e2e.yml` `pull_request.branches`; `release.yml` `gh run list --branch next`), and ZERO false-positives on `next-publish.yml`/comments/awk-`next`/grep-guard. The spike also surfaced the YAML-1.1 `on:`→`true` key remap (a naive parser misses `on.pull_request` entirely), now folded into the design. The lint composes one proven mechanism (yaml.v3 mapping walk) plus a precise command regex.

The release.yml `--branch next` fix is also a latent-bug fix, not only staleness: that query degrades a no-match to `found=false; exit 0` (per its own `|| true` comment), so post-flip producer runs landing on `main` make the journey ledger silently empty forever. Flipping to `--branch main` restores the ledger.

**AC-7 — The shared-core ship-local ceremony's local-merge target resolves to the configured trunk, never a hardcoded `next`.**
Verified by: AC-5's command oracle (the SAME shared `dispatch trunk` test) — `skills/first-officer/references/first-officer-shared-core.md`'s no-PR-host ceremony is FO-read PROSE and a trunk CONSUMER, structurally identical to `87`'s pr-merge mod. The fix points it at `BASE=$(spacedock dispatch trunk --workflow-dir {workflow_dir})` and renders the merge target / sentinel as `{base}` instead of the `next` literal. A doc-only prose consumer cannot carry a separate non-grep AC (a substring sweep over shared-core is the banned prose-grep, and "the prose says run `dispatch trunk`" is the same tautology), so its behavioral proof IS AC-5: the base the ceremony uses is the command's stdout, and AC-5 proves that stdout is the configured trunk (reds on a `next` regression). The prose change itself is the committed doc-diff below; AC-5 is the behavior backing.

## Test plan

| Test | Surface | What it proves | Cost |
|---|---|---|---|
| `TestReconcileTrunkFromConfig` (new, in `reconcile_test.go`) | Go fixture, real git tree | AC-1 + AC-3: D/E detect against `origin/{declared-trunk}`; E item carries `trunk` | Low — extends existing `newReconcileFixture`; add a top-level `trunk: ftrunk` line to the fixture README and an `origin/ftrunk` ref builder mirroring `repoSetOriginNext` |
| `TestResolveTrunk` (new) | Go unit | AC-2: declared top-level trunk returned; `main` fallback on absent key | Low — table test, no git |
| `TestDispatchTrunkCommand` (new) | Go behavior fixture (drives `dispatch trunk`) | AC-5: command prints the resolved trunk as **byte-exact** stdout (`"ftrunk\n"` / `"main\n"`); the shared oracle the prose consumers bind to | Low — mirrors existing `show-stage-def`/`show-standing` command tests |
| `TestTrunkNotConflatedWithChannel` (new or folded into frontdoor tests) | Go unit | AC-4: trunk resolution leaves `devBranch` untouched | Low |
| `TestNoWorkflowResolvesTrunkToNext` (new, `internal/contractlint`) | Go structural lint over `.github/workflows/*.yml` | AC-6: no `on.pull_request.branches` lists `next` and no `gh run list --branch next` query remains; standing guard, two-position check (no blocklist) | Low — yaml.v3 mapping walk + command regex; spike-confirmed it flags exactly the 3 offenders, zero false-positives; CI YAML is non-model-read so the structural assertion is a legitimate oracle |
| Existing `TestReconcileFiveClasses` / `TestReconcileFlipReclassifies` | Go fixture | Regression: the five-class sweep still passes once the fixture README declares `trunk` (these fixtures gain top-level `trunk: next` or `trunk: main` so their existing `origin/next`/`origin/main` graph still matches) | Adjust fixtures only |

No live workflow test needed for the code surfaces: the claim is detector behavior + JSON shape + a command's exact stdout, fully exercised offline by the fixture/behavior harness. All code tests are Go unit/fixture/structural; estimated total complexity Low — the parser (`ParseFrontmatter`), the git-tree harness, the `internal/contractlint` structural-check pattern, and the `dispatch` command-router pattern all already exist. The prose consumers (AC-7 shared-core, `87`'s mod) add NO test of their own — their behavioral backing is AC-5's command oracle by design. The pure doc corrections (roadmap templates) are proven by their committed diffs, not a test.

### Completeness accounting — the cluster DoD ("no helper/mod/ref/doc resolves the integration trunk to `next` post-flip")

Every audit-confirmed trunk reference, mapped to its proof. 11 surfaces total: 5 already covered + the 6 expansion gaps. Per-surface proof strategy is matched to the surface's nature (behavioral test for code; the AC-5 command oracle for config-reading prose consumers; a CI structural lint for workflow YAML; committed diffs for pure go-forward doc corrections).

| # | Surface | Kind | Proof |
|---|---|---|---|
| 1 | `reconcile.go` classD/classE git refs | code | AC-1 (sentinel-trunk fixture, reds on `next`) |
| 2 | `reconcile.go` classE `Reason` string | code | AC-1/AC-3 (drift item carries resolved trunk) |
| 3 | `resolveTrunk` resolver + source key | code | AC-2 (resolver unit) + AC-5 |
| 4 | `claude-first-officer-runtime.md` D/E remedy | prose (data consumer) | AC-3 — reads `{drift.trunk}` from JSON; doc-diff below |
| 5 | `cli.devBranch` channel stamp (un-conflation guard) | code | AC-4 (invariant: trunk resolution leaves devBranch) |
| 6 | `first-officer-shared-core.md` ship-local merge target (188/191) | prose (config consumer) | **AC-7 → AC-5 command oracle**; doc-diff below |
| 7 | `.github/workflows/runtime-live-e2e.yml` PR-base filter | CI YAML | **AC-6 structural lint**; doc-diff below |
| 8 | `.github/workflows/install-e2e.yml` PR-base filter | CI YAML | **AC-6 structural lint**; doc-diff below |
| 9 | `.github/workflows/release.yml` `--branch next` query | CI YAML (+latent bug) | **AC-6 structural lint**; doc-diff below |
| 10 | `docs/roadmap/README.md` reusable sprint checklist (49/51) | doc (go-forward template) | committed diff; doc-diff below |
| 11 | `docs/roadmap/0203-fo-efficiency/dispatch-sprint-execution.md:9` | doc | committed diff; doc-diff below |

Out of this seed's scope (owned elsewhere, do not touch): `87`'s `docs/dev/_mods/pr-merge.md` refit (its base also rides AC-5). Verified non-trunk `next` uses excluded by AC-6's enumeration: the edge marketplace channel (`cli.devBranch=next`), `spacedock-state/*` branches, `@next` action pins, `--ref next` checkouts.

## Documentation & config changes (ideation proposes; implementation applies)

These are doc/config diffs. Where a behavioral AC backs the change, it is named; the diffs themselves are NOT acceptance criteria (a prose-grep is banned as proof). Pure go-forward doc corrections (the roadmap templates) are proven by their committed diffs.

**`docs/dev/README.md` frontmatter — add the canonical top-level declaration (sibling of `state:`):**
```diff
 id-style: sd-b32
 state: .spacedock-state
+trunk: main
 stages:
   defaults:
     worktree: false
```

**`skills/first-officer/references/claude-first-officer-runtime.md` D/E remedies — read the trunk from drift data, not a literal:**
```diff
-   - **D (stale branch)** → `git -C {worktree} pull --rebase origin next`; halt on conflict per the rebase-conflict halt rule.
-   - **E (stale local main)** → `git -C {repo} fetch origin next && git -C {repo} reset --hard origin/next && cd {repo} && go build -o spacedock ./cmd/spacedock`.
+   - **D (stale branch)** → `git -C {worktree} pull --rebase origin {drift.trunk}`; halt on conflict per the rebase-conflict halt rule.
+   - **E (stale local main)** → `git -C {repo} fetch origin {drift.trunk} && git -C {repo} reset --hard origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`.
```
(`{drift.trunk}` is the `trunk` field reconcile now emits on the D/E drift items — AC-3.)

**`dispatch` usage text (`internal/dispatch/dispatch.go` `printUsage`) — register the new subcommand:**
```diff
   spacedock dispatch show-stage-def --workflow-dir DIR --stage STAGE
+  spacedock dispatch trunk --workflow-dir DIR
   spacedock dispatch reconcile --workflow-dir DIR [--team-name NAME] [--repo-root DIR] [--include A,B,C,D,E]
```
(This is the code-surface change that backs AC-5; the usage text is a side effect of adding the `case "trunk":` route, not itself an AC.)

**`skills/first-officer/references/first-officer-shared-core.md` ship-local ceremony (lines ~188-191) — resolve the local-merge target via the configured trunk (gap #6; behavior backed by AC-7→AC-5):**
```diff
-2. Invoke the merge hook (local `--no-ff` merge of `{branch}` onto `next`).
+2. Resolve the integration trunk `BASE=$(spacedock dispatch trunk --workflow-dir {workflow_dir})` (configured trunk, default `main`), then invoke the merge hook (local `--no-ff` merge of `{branch}` onto `{BASE}`).
 3. Record the merge so the terminal guard is satisfied without `--force`:
    - If `merge: local`, the policy exempts the pr-requirement — skip to step 4.
-   - Otherwise set the post-merge sentinel `spacedock status --workflow-dir {workflow_dir} --set {slug} pr=local-merge:{short-sha}` (the merge commit on `next`; set ONLY after merge has landed; commit path-scoped). The status table renders as `{short-sha} (local)`.
+   - Otherwise set the post-merge sentinel `spacedock status --workflow-dir {workflow_dir} --set {slug} pr=local-merge:{short-sha}` (the merge commit on `{BASE}`; set ONLY after merge has landed; commit path-scoped). The status table renders as `{short-sha} (local)`.
```
(The exact-stdout contract of AC-5 is load-bearing here: `$(...)` captures the bare trunk so `{BASE}` is a clean branch name. Resolving the base inline at step 2 avoids renumbering the 1-5 ceremony.)

**`.github/workflows/runtime-live-e2e.yml` (gap #7; backed by AC-6) — drop `next` from the PR-base filter:**
```diff
   pull_request:
-    branches: [next, main]
+    branches: [main]
```

**`.github/workflows/install-e2e.yml` (gap #8; backed by AC-6) — same:**
```diff
   pull_request:
-    branches: [next, main]
+    branches: [main]
```

**`.github/workflows/release.yml` journey-ledger query (gap #9; backed by AC-6; also a latent-bug fix) — query the producer run on the post-flip trunk:**
```diff
-          run_id="$(gh run list --workflow "Runtime Live E2E" --branch next --status success --limit 1 --json databaseId --jq '.[0].databaseId' || true)"
+          run_id="$(gh run list --workflow "Runtime Live E2E" --branch main --status success --limit 1 --json databaseId --jq '.[0].databaseId' || true)"
           if [ -z "$run_id" ] || [ "$run_id" = "null" ]; then
-            echo "::warning::no successful Runtime Live E2E run found on next; skipping journey ledger" >&2
+            echo "::warning::no successful Runtime Live E2E run found on main; skipping journey ledger" >&2
```

**`docs/roadmap/README.md` reusable sprint checklist (gap #10; pure go-forward doc correction — this template SEEDS future sprints via line 33 "Copy into a sprint's `index.md`", so leaving it propagates the bug):**
```diff
-- [ ] **Merge** each to `next` (PR-merge); keep state commits concurrency-safe
+- [ ] **Merge** each to `main` (PR-merge); keep state commits concurrency-safe
...
-- [ ] **⚠️ Pre-cut antipattern audit** — with all members merged to `next` and the tag **not yet fired**, dispatch an *independent* reviewer …
+- [ ] **⚠️ Pre-cut antipattern audit** — with all members merged to `main` and the tag **not yet fired**, dispatch an *independent* reviewer …
```

**`docs/roadmap/0203-fo-efficiency/dispatch-sprint-execution.md:9` (gap #11; pure doc correction; the sibling 0203 index.md DoD line was already fixed by the FO):**
```diff
-**0.20.3** = the FO-efficiency restructure + the context-budget probe fix. Done when, merged to `next` (then `main` at the cut) — see `index.md` Definition of Done.
+**0.20.3** = the FO-efficiency restructure + the context-budget probe fix. Done when merged to `main` — see `index.md` Definition of Done.
```

## Notes

Surfaced by the 0202 drive. **Coupled pair:** this seed (`sr`) OWNS the single-source design AND the cluster-wide completeness — a **top-level `trunk:` key** in `docs/dev/README.md`, read by one `resolveTrunk` resolver, surfaced as `spacedock dispatch trunk`, with every audit-confirmed `next`-as-trunk reference accounted for (see the Completeness accounting table: 11 surfaces, each with a per-surface proof). The sibling `87` (`pr-merge` mod base branch) CONSUMES the same source by reading the **output of `spacedock dispatch trunk --workflow-dir docs/dev`** — NOT a second key, NOT a `_mods`-local config, NOT a prose literal — and owns ONLY the `docs/dev/_mods/pr-merge.md` refit (out of this seed's scope). Key placement converged with `87` to top-level (a sibling of `state:`) so one `ParseFrontmatter` read serves both the command and reconcile; the `dispatch trunk` command was adopted from `87`'s ask because it is the only thing that gives the prose consumers a real, fail-able oracle.

**Completeness scope (post-staff-review expansion):** the cluster DoD is "no helper/mod/ref/doc resolves the integration trunk to `next` post-flip." Beyond reconcile + the FO-runtime D/E remedy (cycle 1), the audit found 6 more operative references: the shared-core ship-local merge target (AC-7→AC-5), three CI workflows (AC-6 structural lint, including a latent journey-ledger bug in release.yml), and two go-forward roadmap doc templates (committed diffs). **Shared exact-stdout contract:** the prose consumers capture `dispatch trunk` via `$(...)`, so AC-5 now requires byte-exact stdout (bare branch name) — this is the sole load-bearing wire between this seed and both prose consumers (`87`'s mod and shared-core). **Mis-citation corrected** (shared with `87`): `reconcile.go:296` calls `ParseFrontmatter` on per-ENTITY files, not the README; the README is read via `ParseStagesWithDefaults` at `:319`. The "no new import" conclusion holds (`internal/status` imported at `:17`), but `resolveTrunk`'s `ParseFrontmatter`-on-README call is a NEW one-line site.

## Stage Report: ideation (cycle 0 — SUPERSEDED, moved to top-level `trunk:` in cycle 1)

> Superseded: this cycle-0 report names the nested `stages.defaults.trunk` source, which cycle 1 replaced with a top-level `trunk:` key. Retained for provenance; the source-key references below are stale.

- DONE: The design de-conflates `dispatch reconcile`'s git-hygiene (Class-D/E branch/trunk staleness) from its team-management (Class-A/B/C roster drift) — trunk knowledge no longer lives inside the roster-management helper.
  Design extracts a `resolveTrunk(workflowDir)` resolver out of `classD`/`classE`; the trunk literal becomes a parameter sourced from one declaration; roster detectors (classA/B from ~/.claude/teams) stay untouched. See "Proposed approach (firm)" §2.
- DONE: Class-D/E resolve the integration trunk to `main` from ONE canonical trunk-config source (the source 87 will also consume), never a hardcoded `next`; the proof is a real test that reds if reconcile resolves the trunk to `next` or a fixture asserting it reads the canonical source — never a prose/string grep.
  Source settled: `stages.defaults.trunk` in `docs/dev/README.md` *(superseded — moved to top-level `trunk:` in cycle 1)*. AC-1 is a fixture test using a sentinel trunk `ftrunk` (neither next nor main) so a re-hardcoded `next` reds the test against an independent oracle. AC-4 keeps the `cli.devBranch` channel stamp un-conflated. 87 directed to the same key.
- DONE: The riskiest unknown is exercised before committing the design: WHERE the single trunk-config source lives and that `reconcile` actually reads it — demonstrated with the smallest end-to-end exercise, not asserted.
  Throwaway Go test against real `status.ParseStagesWithDefaults`: `stages.defaults.trunk: main` → `defaults["trunk"]=="main"` PASS; absent key → no panic, fallback reachable PASS *(superseded — moved to top-level `trunk:` in cycle 1)*. See "Spike result".

### Summary
*(Superseded — moved to top-level `trunk:` in cycle 1.)* Settled the single canonical trunk-config source as `stages.defaults.trunk` in the workflow README — the source `87`'s pr-merge mod also consumes — and designed the de-conflation as extracting a `resolveTrunk` resolver out of the `classD`/`classE` detectors with a `main` fallback (hardcoded `next` deleted, not demoted). Key anti-conflation call: kept `cli.devBranch` (marketplace channel stamp, where `next` is correct) explicitly separate from the integration trunk (AC-4). Spiked the riskiest unknown first.

## Stage Report: ideation (cycle 1 — coordination convergence with `87`)

Revised the single-source design after coupled-pair coordination with `87` (pr-merge mod base branch). No checklist item changed verdict; the source key and consumption surface were refined to satisfy BOTH halves of the pair:

- DONE (refined): Class-D/E resolve the integration trunk from ONE canonical source, never a hardcoded `next`.
  Source moved from nested `stages.defaults.trunk` → **top-level `trunk:` key** (sibling of `state:`), read by one `resolveTrunk` resolver via `status.ParseFrontmatter`. Top-level chosen because `ParseFrontmatter` renders nested mappings empty, so a nested key would be invisible to the simple parser the shared command needs. Both placements spike-confirmed; rejected-placement spike justifies the choice.
- DONE (added): The resolver is surfaced as `spacedock dispatch trunk --workflow-dir DIR` so `87`'s PROSE mod gets a real non-grep oracle (new AC-5). reconcile and the command call the SAME `resolveTrunk`, so the two consumers cannot diverge. The `dispatch` router (dispatch.go:26) extends with a `case "trunk":` trivially — confirmed against the existing show-standing/show-stage-def routes.

### Summary
Converged with `87` on the single trunk-config source: a **top-level `trunk:` key** in `docs/dev/README.md`, one `resolveTrunk` resolver, surfaced as the new `spacedock dispatch trunk` command. reconcile (Go) calls the resolver directly; `87`'s pr-merge prose mod consumes the command's stdout — one oracle, no second key, no prose literal. Added AC-5 (command behavior) as the shared, co-owned oracle. Re-spiked the top-level placement end-to-end (`ParseFrontmatter(README)["trunk"]` round-trips; absent→`""`→fallback `main`) and verified the nested placement is invisible to `ParseFrontmatter`, which is what settled top-level over nested.

## Stage Report: ideation (cycle 2 — completeness expansion after staff review + trunk-config audit)

The single-source DESIGN (top-level `trunk:` + `resolveTrunk` + `dispatch trunk`) was accepted; the joint gate was held to expand coverage to the full cluster DoD ("no helper/mod/ref/doc resolves the integration trunk to `next` post-flip"). Expanded from 5 covered surfaces to all 11 (5 + 6 audit gaps), each with a per-surface proof matched to its nature.

- DONE: Completeness — all 6 audit-confirmed gaps folded in, each verified against the live file.
  Gap #6 shared-core ship-local merge target (188/191, confirmed) → resolve via `dispatch trunk`, proof rides new AC-7→AC-5. Gaps #7/#8 CI PR-base filters `branches: [next, main]` (confirmed runtime-live-e2e.yml:46, install-e2e.yml:15) → `[main]`. Gap #9 release.yml `--branch next` (confirmed:56) → `--branch main`; verified the latent bug (its own `|| true` comment degrades no-match to found=false/exit 0, so post-flip producer runs on main are silently skipped). Gaps #10/#11 roadmap go-forward templates (confirmed README.md:49/51 reusable checklist seeded via :33; dispatch-sprint-execution.md:9) → committed diffs.
- DONE: Per-surface proof strategy stated — behavioral ACs for code (AC-1/2/3/4), AC-5 command oracle for config-reading prose consumers (AC-7 shared-core; 87's mod), a CI structural lint for the workflows (new AC-6 in internal/contractlint — CI YAML is non-model-read so a structural assertion is a legitimate non-grep oracle), committed diffs for pure doc corrections. Completeness accounting table maps all 11 surfaces.
- DONE: AC-6's new lint mechanism spiked before committing (riskiest-unknown discipline). A two-position lint (yaml.v3 walk of on.pull_request.branches + regex for `gh run list --branch next`) run against the live workflows flags EXACTLY the 3 real offenders with ZERO false-positives on next-publish.yml/comments/awk-`next`/the grep guard — proving the two-structural-position approach beats a token-sweep-minus-blocklist. Spike also surfaced the YAML-1.1 `on:`→`true` bare-key quirk, now folded in. Enumerated all `.github/workflows` `next` uses to confirm the exclusion reasoning.
- DONE: Soundness polish folded.
  AC-5 now requires BYTE-EXACT stdout (`"ftrunk\n"`/`"main\n"`) — the sole load-bearing contract for the `$(...)`-capturing prose consumers. Mis-citation corrected in body AND flagged to 87 (reconcile.go:296 reads entity files not README; resolveTrunk's ParseFrontmatter-on-README is a new call). Cycle-0 stage report marked superseded. AC-2's "literal origin/next no longer appears" grep aside dropped.

### Summary
Accepted-design completeness pass: kept the cycle-1 single-source design intact and expanded coverage to the whole audit inventory — 11 surfaces, each with a proof matched to its kind (behavioral test for code, the AC-5 command oracle for config-reading prose consumers, a non-model-read CI structural lint for the workflow YAML including a real journey-ledger latent-bug fix, committed diffs for go-forward doc templates). Added AC-6 (CI lint) and AC-7 (shared-core ship-local consumer); tightened AC-5 to byte-exact stdout as the load-bearing wire to both prose consumers; corrected the shared mis-citation; marked cycle 0 superseded. 87's pr-merge-mod scope is untouched (frozen). The completeness DoD is now provable per-surface, not asserted.

## Stage Report: implementation

- DONE: resolveTrunk resolver + top-level `trunk: main` README key + `spacedock dispatch trunk` command shipped; reconcile classD/classE consume resolveTrunk (no `next`/`main` literal); TestResolveTrunk, TestDispatchTrunkCommand (byte-exact `"ftrunk\n"`/`"main\n"` stdout, AC-5), and TestReconcileTrunkFromConfig (sentinel `ftrunk` fixture; classE driftItem carries trunk) all green.
  internal/dispatch/trunk.go (resolveTrunk + runDispatchTrunk), reconcile.go (trunk resolved once in Reconcile, passed to classD/classE as param; driftItem.Trunk added; refs are HEAD..origin/{trunk} & origin/{trunk}..main), dispatch.go (`case "trunk":` + usage line), docs/dev/README.md `trunk: main`. Tests green (commit f3729ae2). Real-config exercise: `dispatch trunk --workflow-dir docs/dev` → byte-exact `main\n`, stderr 0 bytes, `BASE=$(...)` captures `main`; reconcile D/E against real docs/dev emits Class-E with `trunk: main`.
- DONE: AC-6 lint TestNoWorkflowResolvesTrunkToNext in internal/contractlint flags exactly the 3 offenders (install-e2e.yml + runtime-live-e2e.yml `pull_request.branches`; release.yml `gh run list --branch next`) with zero false-positives; and the 3 workflow YAML fixes are applied — including release.yml `--branch next`→`main` (the v0.20.3 journey-ledger cut dependency).
  internal/contractlint/workflow_trunk_test.go: two-position lint (yaml.Node walk of on.pull_request.branches + `gh run list…--branch next\b` regex) + a discriminator control proving precision (flags both offending positions, ignores edge-publish/@next/--ref next/awk next/comment English). Red on the 3 live offenders before the fix, green after (commit f7de4312). yaml.v3 keeps `on` a string key (no YAML-1.1 true remap) — walk tolerates both defensively. Stale `next` comments in runtime-live-e2e.yml corrected to `main`.
- DONE: All 11 Completeness-table surfaces resolved: prose doc-diffs applied (#4 D/E remedy reads `{drift.trunk}`; #6 ship-local merge target via `dispatch trunk`; #10/#11 roadmap templates) + AC-4 un-conflation guard (trunk resolution leaves `cli.devBranch=="next"`); full `go test ./...` green; 87's `docs/dev/_mods/pr-merge.md` left untouched (frozen).
  #1/#2/#3 code (commit f3729ae2); #4 claude-fo-dispatch.md:228-229 + #6 claude-fo-merge.md:31,34 (RELOCATED from the spec's pre-split paths by #367 — same diff content; flagged to team-lead); #5 AC-4 internal/cli/trunk_channel_unconflated_test.go; #7/#8/#9 CI YAML (commit f7de4312); #10 docs/roadmap/README.md, #11 dispatch-sprint-execution.md:9 (commit c6de781c). `go test ./...` all packages ok; gofmt + go vet clean; `git diff f87107b1 -- docs/dev/_mods/pr-merge.md` empty.

### Summary
Shipped the single canonical trunk-config source: a top-level `trunk:` README key, one `resolveTrunk(workflowDir)` resolver (default `main`), surfaced as `spacedock dispatch trunk` (byte-exact stdout, AC-5). reconcile's classD/classE consume the resolver via a parameter — the `next` literal is deleted, the de-conflation is the trunk leaving the roster-helper detectors — and classE's driftItem carries the resolved trunk so the FO remedy is data-driven (AC-3). All 7 ACs proven by independent oracles (sentinel `ftrunk` fixture, byte-exact command test, the un-conflation invariant, the two-position CI lint flagging exactly 3 offenders). NOTE for validation: the entity's doc-diffs named pre-split FO files; #367 (in my base) relocated the D/E remedy to claude-fo-dispatch.md and the ship-local ceremony to claude-fo-merge.md — I applied the identical diff content at the new locations (flagged to team-lead). pr-merge.md (task 87, frozen) untouched.

## Stage Report: validation

- DONE: AC-1 — Class-D/E resolve trunk from configured `trunk:` key (sentinel `ftrunk`).
  PASSED. `TestReconcileTrunkFromConfig` green from worktree; audit edit re-hardcoding `next` in classD/E reds it ("expected 1 D entry against origin/ftrunk; got 0").
- DONE: AC-2 — trunk literal absent from classD/E; sourced via one `resolveTrunk`.
  PASSED. `TestResolveTrunk` 3/3 subtests green (declared-sentinel, declared-main, absent→main fallback); reconcile.go classD/E take `trunk` param, ref strings are `origin/`+trunk.
- DONE: AC-3 — Class-E driftItem carries resolved `trunk`.
  PASSED. `TestReconcileTrunkFromConfig` asserts E item `trunk=="ftrunk"`; driftItem.Trunk at reconcile.go:149/627.
- DONE: AC-4 — channel stamp un-conflated with trunk.
  PASSED. `TestTrunkNotConflatedWithChannel` green (internal/cli/trunk_channel_unconflated_test.go).
- DONE: AC-5 — `dispatch trunk` prints byte-exact stdout.
  PASSED. `TestDispatchTrunkCommand` green (`"ftrunk\n"`/`"main\n"`); audit edit emitting a stray second stdout line reds it ("ftrunk\nextra\n").
- DONE: AC-6 — no workflow resolves trunk to `next` post-flip.
  PASSED. `TestNoWorkflowResolvesTrunkToNext` green; flags exactly 3 offenders (audit re-adding `next` to install-e2e pull_request.branches reds it); zero false-positives (next-publish.yml present, lint stays green).
- DONE: AC-7 — ship-local merge target resolves to configured trunk.
  PASSED via AC-5 oracle. claude-fo-merge.md step 2 captures `BASE=$(spacedock dispatch trunk ...)`; behavioral backing is AC-5's byte-exact stdout.
- DONE: Detached adversarial audit (MANDATORY).
  Refuted nothing material. Separate throwaway checkout; all 3 claim-breaking edits red the intended tests (re-hardcode next→AC-1 red; next in PR branches→AC-6 red; stray stdout line→AC-5 byte-exact red). AC-6 false-positive guard confirmed (next-publish.yml present, lint green).
- DONE: #367 relocation sound.
  Gap #4 (`{drift.trunk}` D/E remedy) landed in claude-fo-dispatch.md; gap #6 (`BASE=$(dispatch trunk)` ship-local) in claude-fo-merge.md — diff content matches spec. Pre-split files (claude-first-officer-runtime.md, first-officer-shared-core.md) no longer carry the remedy/merge content (grep empty). AC-3/AC-5 are code-behavioral, unaffected by the doc move.
- DONE: 87's pr-merge.md untouched.
  `git diff f87107b1 -- docs/dev/_mods/pr-merge.md` empty.
- DONE: full `go test ./...` green from worktree root.
  Exit 0, no failures.

### Summary
PASSED. All 7 ACs verified by reproducing their independent oracles from the worktree (not trusting the report): named tests green, full `go test ./...` exit 0. The mandatory detached adversarial audit on a separate throwaway checkout confirmed all three claim-breaking edits red the tests that should catch them, and the AC-6 lint does not false-positive on legitimate `next` uses — refuted nothing material. The #367 doc relocation is sound (gap #4/#6 content at the post-split paths, pre-split files cleared, behavioral ACs unaffected) and 87's frozen pr-merge.md is untouched.
