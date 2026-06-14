---
id: sryzghzqazj9s9km6ebqkf5s
title: dispatch reconcile conflates team-management with repo-hygiene (and hardcodes pre-flip trunk `next`)
status: ideation
source: "0202 Commander drive (2026-06-13). Boot reconcile flagged Class-D/E drift against origin/next; Class-E remedy 'reset main->origin/next' would have reverted the entire post-flip trunk. Investigation (captain-prompted) found the deeper cause: reconcile bundles git-hygiene into a team-management helper, so it carries repo/trunk knowledge it shouldn't."
group: cleanup
sprint: 0203-fo-efficiency
started: 2026-06-14T05:04:07Z
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
- **One parser serves both consumers:** `status.ParseFrontmatter` surfaces top-level scalar keys as a flat `map[string]string`, so `fm["trunk"]` is free for the `dispatch trunk` command AND reconcile. Critically, `ParseFrontmatter` renders NESTED mappings as empty string (frontmatter.go:123-126) — so a `stages.defaults.trunk` key would be **invisible** to the simple parser the command uses, forcing the heavier `ParseStagesWithDefaults` and a second read path. Top-level keeps exactly one parser behind one resolver. reconcile already imports `status` and already calls `ParseFrontmatter` (reconcile.go:296), so no new import.

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

Conclusion: the source is a top-level `trunk:` key, read via `status.ParseFrontmatter` — a parser reconcile **already imports and calls** (reconcile.go:296) and that the `dispatch trunk` command can call identically. The design composes one proven mechanism (`ParseFrontmatter` returning a flat map) plus a one-line lookup behind a single `resolveTrunk` resolver. No on-disk format, no new handoff — the only unverified bit was the round-trip for a new key, now verified for both placements.

## Acceptance criteria

**AC-1 — Class-D/E resolve the integration trunk from the configured top-level `trunk:` key, never a hardcoded `next`.**
Verified by: a reconcile fixture test (extending `reconcile_test.go`'s real-git-tree harness) where the fixture README declares top-level `trunk: ftrunk` (a sentinel name that is neither `next` nor `main`) and the fixture git graph builds `origin/ftrunk` as the trunk ref. The test asserts Class-D fires when a worktree branch is behind `origin/ftrunk` and Class-E fires when local `main` carries commits not on `origin/ftrunk`. **The expected ref name comes from the fixture's declared trunk, not from any literal in reconcile.go** — so a regression that re-hardcodes `next` (or `main`) reds the test (it would query `origin/next`, find no such ref / wrong count, and fail to detect against `origin/ftrunk`). This is an independent oracle: the fixture trunk and the code's resolved trunk can disagree, which is what makes it able to fail. NOT a prose/string grep.

**AC-2 — The trunk literal is absent from `classD`/`classE`; trunk knowledge is sourced through one `resolveTrunk` resolver.**
Verified by: a Go test asserting `resolveTrunk` returns the declared trunk for a README with top-level `trunk:` set and returns `main` (the fallback) for a README with no trunk key — expected values from the fixture READMEs, not from the function's own source. Combined with AC-1's behavioral test, this proves D/E read the resolved value rather than a literal. (De-conflation is checkable in code/tests: the literal `"origin/next"` no longer appears in the `classD`/`classE` git-ref construction — but the *binding* proof is AC-1's behavioral red-on-regression, not a grep for the absence of the string.)

**AC-3 — Class-E's drift item carries the resolved `trunk`, so the FO remedy is data-driven, not prose-hardcoded.**
Verified by: the AC-1 fixture test also asserts the emitted Class-E `driftItem` has `trunk == "ftrunk"` (from the fixture's declared trunk). Expected value from the fixture config, not the code. This is the behavioral proof that the FO can read the trunk from JSON instead of a hardcoded remedy string; the FO-runtime prose change is doc-only and is NOT itself an acceptance criterion (a prose-grep over the runtime file is banned).

**AC-4 — The marketplace channel stamp stays un-conflated with the trunk.**
Verified by: a test (or assertion within the existing frontdoor tests) that `resolveTrunk` does not read `cli.devBranch` and `devBranch`'s value is unchanged by trunk resolution — i.e. resolving the trunk on a workflow whose README declares `trunk: main` leaves `devBranch == "next"`. Expected values are the two independent declared sources. This guards against re-bundling the two axes. (Cheap: it is an invariant over two already-separate variables.)

**AC-5 — `spacedock dispatch trunk --workflow-dir DIR` prints the resolved trunk; this is the shared oracle `87`'s prose mod binds to.**
Verified by: a command-level Go test (behavior fixture driving the `dispatch trunk` subcommand) where a fixture README declares top-level `trunk: ftrunk` and the command's stdout is asserted to be `ftrunk`; a second fixture with no trunk key asserts stdout `main`. **Expected value from the fixture README, not from the command's source** — so the command and `87`'s pr-merge base branch share ONE oracle (both resolve through `resolveTrunk`). This is the command surface that lets `87`'s prose mod be config-driven with a real, fail-able oracle instead of a banned prose-grep. (This AC is co-owned: `sr` ships the command + this test; `87`'s own AC tests its mod consuming the command's output.)

## Test plan

| Test | Surface | What it proves | Cost |
|---|---|---|---|
| `TestReconcileTrunkFromConfig` (new, in `reconcile_test.go`) | Go fixture, real git tree | AC-1 + AC-3: D/E detect against `origin/{declared-trunk}`; E item carries `trunk` | Low — extends existing `newReconcileFixture`; add a top-level `trunk: ftrunk` line to the fixture README and an `origin/ftrunk` ref builder mirroring `repoSetOriginNext` |
| `TestResolveTrunk` (new) | Go unit | AC-2: declared top-level trunk returned; `main` fallback on absent key | Low — table test, no git |
| `TestDispatchTrunkCommand` (new) | Go behavior fixture (drives `dispatch trunk`) | AC-5: command prints the resolved trunk; the shared oracle `87` binds to | Low — mirrors existing `show-stage-def`/`show-standing` command tests |
| `TestTrunkNotConflatedWithChannel` (new or folded into frontdoor tests) | Go unit | AC-4: trunk resolution leaves `devBranch` untouched | Low |
| Existing `TestReconcileFiveClasses` / `TestReconcileFlipReclassifies` | Go fixture | Regression: the five-class sweep still passes once the fixture README declares `trunk` (these fixtures gain top-level `trunk: next` or `trunk: main` so their existing `origin/next`/`origin/main` graph still matches) | Adjust fixtures only |

No live workflow test needed: the claim is detector behavior + JSON shape + a command's stdout, fully exercised offline by the fixture/behavior harness. All tests are Go unit/fixture; estimated total complexity Low — the parser (`ParseFrontmatter`), the git-tree harness, reconcile's import of `status`, and the `dispatch` command-router pattern all already exist.

## Documentation changes (ideation proposes; implementation applies)

Two instruction-file surfaces name `next` as the trunk in D/E remedies and must follow the resolved trunk. These are doc diffs, NOT acceptance criteria (prose-grep is banned as proof; the behavioral proof is AC-1/AC-3):

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

## Notes

Surfaced by the 0202 drive. **Coupled pair:** this seed (`sr`) OWNS the single-source design — a **top-level `trunk:` key** in `docs/dev/README.md`, read by one `resolveTrunk` resolver, surfaced as `spacedock dispatch trunk`. The sibling `87` (`pr-merge` mod base branch) CONSUMES the same source by reading the **output of `spacedock dispatch trunk --workflow-dir docs/dev`** — NOT a second key, NOT a `_mods`-local config, NOT a prose literal. Key placement converged with `87` to top-level (a sibling of `state:`) so one `ParseFrontmatter` read serves both the command and reconcile; the `dispatch trunk` command was adopted from `87`'s ask because it is the only thing that gives the prose mod a real, fail-able oracle. The canonical `mods/pr-merge.md` (repo root) was ALREADY migrated to base `main`; the un-migrated copy is the workflow's `docs/dev/_mods/pr-merge.md` (still base `next`) — that is `87`'s target. A parallel 5-surface trunk-config audit cross-checks the completeness of the `next`-as-trunk inventory against this design.

## Stage Report: ideation

- DONE: The design de-conflates `dispatch reconcile`'s git-hygiene (Class-D/E branch/trunk staleness) from its team-management (Class-A/B/C roster drift) — trunk knowledge no longer lives inside the roster-management helper.
  Design extracts a `resolveTrunk(workflowDir)` resolver out of `classD`/`classE`; the trunk literal becomes a parameter sourced from one declaration; roster detectors (classA/B from ~/.claude/teams) stay untouched. See "Proposed approach (firm)" §2.
- DONE: Class-D/E resolve the integration trunk to `main` from ONE canonical trunk-config source (the source 87 will also consume), never a hardcoded `next`; the proof is a real test that reds if reconcile resolves the trunk to `next` or a fixture asserting it reads the canonical source — never a prose/string grep.
  Source settled: `stages.defaults.trunk` in `docs/dev/README.md`. AC-1 is a fixture test using a sentinel trunk `ftrunk` (neither next nor main) so a re-hardcoded `next` reds the test against an independent oracle. AC-4 keeps the `cli.devBranch` channel stamp un-conflated. 87 directed to the same key.
- DONE: The riskiest unknown is exercised before committing the design: WHERE the single trunk-config source lives and that `reconcile` actually reads it — demonstrated with the smallest end-to-end exercise, not asserted.
  Throwaway Go test against real `status.ParseStagesWithDefaults`: `stages.defaults.trunk: main` → `defaults["trunk"]=="main"` PASS; absent key → no panic, fallback reachable PASS. reconcile.go already imports+calls this parser (reconcile.go:319). See "Spike result".

### Summary
Settled the single canonical trunk-config source as `stages.defaults.trunk` in the workflow README — the source `87`'s pr-merge mod also consumes — and designed the de-conflation as extracting a `resolveTrunk` resolver out of the `classD`/`classE` detectors with a `main` fallback (hardcoded `next` deleted, not demoted). Key anti-conflation call: kept `cli.devBranch` (marketplace channel stamp, where `next` is correct) explicitly separate from the integration trunk (AC-4). Spiked the riskiest unknown first — the parser already returns the trunk key and reconcile already imports it — so the design composes one proven mechanism plus a map lookup; all ACs are Go fixture/unit tests on an independent oracle, none a prose-grep.

## Stage Report: ideation (cycle 1 — coordination convergence with `87`)

Revised the single-source design after coupled-pair coordination with `87` (pr-merge mod base branch). No checklist item changed verdict; the source key and consumption surface were refined to satisfy BOTH halves of the pair:

- DONE (refined): Class-D/E resolve the integration trunk from ONE canonical source, never a hardcoded `next`.
  Source moved from nested `stages.defaults.trunk` → **top-level `trunk:` key** (sibling of `state:`), read by one `resolveTrunk` resolver via `status.ParseFrontmatter`. Top-level chosen because `ParseFrontmatter` renders nested mappings empty, so a nested key would be invisible to the simple parser the shared command needs. Both placements spike-confirmed; rejected-placement spike justifies the choice.
- DONE (added): The resolver is surfaced as `spacedock dispatch trunk --workflow-dir DIR` so `87`'s PROSE mod gets a real non-grep oracle (new AC-5). reconcile and the command call the SAME `resolveTrunk`, so the two consumers cannot diverge. The `dispatch` router (dispatch.go:26) extends with a `case "trunk":` trivially — confirmed against the existing show-standing/show-stage-def routes.

### Summary
Converged with `87` on the single trunk-config source: a **top-level `trunk:` key** in `docs/dev/README.md`, one `resolveTrunk` resolver, surfaced as the new `spacedock dispatch trunk` command. reconcile (Go) calls the resolver directly; `87`'s pr-merge prose mod consumes the command's stdout — one oracle, no second key, no prose literal. Added AC-5 (command behavior) as the shared, co-owned oracle. Re-spiked the top-level placement end-to-end (`ParseFrontmatter(README)["trunk"]` round-trips; absent→`""`→fallback `main`) and verified the nested placement is invisible to `ParseFrontmatter`, which is what settled top-level over nested.
