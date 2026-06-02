---
id: qs87q0ca1wa3bhzfkj07mgwb
title: Harden FO reconciliation — teardown agents at terminal + supersede-shutdown + a state-keyed reconcile sweep (no lingering ensigns / stale branches / un-advanced PRs)
status: validation
source: "captain (2026-06-02) — observed the FO leak 3 ensigns (the 7h implementer post-merge + the old 3g cycle agents post-rework) and reconcile in-flight branches only reactively at merge (02 hit #263 late); all keyed off drift-prone FO session memory"
completed:
verdict:
score: "0.30"
worktree: .worktrees/spacedock-ensign-ensign-lifecycle-reconcile
issue:
started: 2026-06-02T21:14:43Z
mod-block: merge:pr-merge
---

The FO leaks agents and reconciles reactively because every lifecycle rule keys off **FO session memory**, which drifts under context pressure (this session: the FO tracked "5 alive" when 7 were — the 7h implementer lingered after its entity merged+archived, and the superseded old-3g `-implementation`/`-validation` lingered after 3g re-ideated). Same root cause as the prose-pin tests: a "remember to shut down / remember to rebase" rule can't self-check. The fix is to **compute the reconciliation set from state** that already exists on disk (team `config.json` members[] + entity frontmatter + git refs), and to let a code gate flag drift the FO would otherwise miss.

## Where it fails today (contract gaps)

1. **No teardown at terminal.** `first-officer-shared-core.md` `## Merge and Cleanup` step 9 removes the worktree + local branch but has NO step to shut down the entity's agents (implementer, validator, any `-N`/`-cycleN` variants, the detached auditor). → the 7h-class leak.
2. **No supersede-shutdown.** When an entity re-enters an earlier stage and a fresh / `-N`-suffixed agent supersedes a prior-cycle one (re-ideation, rework), nothing shuts down the superseded agent. → the old-3g-class leak.
3. **Reactive, memory-based reconciliation.** Merged-PR advancement, in-flight-branch freshness vs `next`, and agent liveness are all tracked in FO memory and reconciled only at merge time (02 hit the #263 conflict late — the pre-merge rebase caught it, but reactively). Local `main` accumulates merge-bubbles and the `./spacedock` binary runs stale, but neither is checked.

## Mechanism — what state is authoritative

Three disk-resident sources, all already present, all read-only:

1. **Team roster** — `~/.claude/teams/{team_name}/config.json` `members[]` (the same file `internal/claudeteam/contextbudget.go` already opens). Each member carries `name`, `agentType`, and `model`. `agentType == "spacedock:ensign"` is the ensign discriminator; anything else (e.g. `team-lead`, `comm-officer`) is exempt — standing teammates and the FO itself never get swept.
2. **Entity frontmatter** — the workflow's active + archived entity directories on the state branch (`{workflow_dir}/{state}/*/index.md` and `{workflow_dir}/{state}/_archive/*/index.md`). The fields read are `id`, `slug`, `status`, `worktree`, `pr`, `completed`. The status viewer's `--where`/`--json` interface already exposes everything we need; the helper consumes it, never re-parses YAML.
3. **Git refs** — `git -C {worktree} rev-list --count origin/next..HEAD` for branch-distance, `git -C {repo} rev-parse main origin/next` for local-main-staleness, `gh pr view {N} --json state` for PR-state confirmation when `pr:` is set.

The reconcile helper does NOT write to any of these; it emits a structured drift report. The FO acts on the report.

## Mechanism — the sweep algorithm

The helper computes four disjoint drift classes from the three sources above. Pseudocode-level spec:

```
sweep:
  ensigns          = [m for m in team.members if m.agentType == "spacedock:ensign"]
  active_entities  = status --where "status!=" --json   # non-archived
  archived         = status --archived --json
  in_flight        = [e for e in active_entities if e.worktree != ""]

  # Class A — LINGERING (terminal entity, live agent)
  for m in ensigns:
    slug, stage = decompose(m.name)        # strip "spacedock-ensign-" prefix; trailing "-{known_stage}"
    if archived.has(slug) or active.find(slug).status == "done":
      emit A: {name: m.name, slug, reason: "entity archived/terminal"}

  # Class B — SUPERSEDED (current cycle's ensign for the same slug+stage is a different name)
  for m in ensigns:
    slug, stage = decompose(m.name)
    cohort = [n for n in ensigns if decompose(n.name) == (slug, stage)]
    if len(cohort) > 1:
      winner = highest "-cycleN" (or unsuffixed=cycle1)
      losers = cohort - {winner}
      emit B for each loser: {name, slug, stage, reason: "superseded by {winner}"}

  # Class C — UN-ADVANCED PR (entity's PR merged on origin but status not terminal)
  for e in active_entities where e.pr != "":
    if gh.pr_state(e.pr) == "MERGED" and e.status != "done":
      emit C: {slug: e.slug, pr: e.pr, reason: "PR merged but status={e.status}"}

  # Class D — STALE BRANCH (worktree branch behind origin/next by ≥1 commit)
  for e in in_flight:
    behind = git -C {e.worktree} rev-list --count HEAD..origin/next
    if behind > 0:
      emit D: {slug: e.slug, worktree: e.worktree, behind, reason: "branch behind origin/next"}

  # Class E — STALE LOCAL MAIN (project's main has merge-bubbles vs origin/next)
  ahead = git rev-list --count origin/next..main
  if ahead > 0:
    emit E: {ahead, reason: "local main carries {ahead} commits not on origin/next; reset main->origin/next"}
```

Output is one JSON envelope `{"command":"reconcile","drift":[{"class":"A|B|C|D|E", ...}, ...]}`. Empty `drift[]` is a green sweep.

### Name decomposition

`decompose(name)` strips the leading `worker_key-` (default `spacedock-ensign-`), then strips a trailing `-{stage}` where stage is one of the workflow's README-declared stages (or any of the five-canonical ensign stages: `backlog | ideation | implementation | validation | done`). What remains is `{slug}` with an optional trailing `-cycleN` or `-N`. The cycle suffix lifts off cleanly because stage names never contain `-cycleN`. Stage names are read from the workflow README at sweep time, NOT hard-coded — so a workflow with custom stages works without code changes.

This decomposition is a string-shape contract on the dispatch helper's existing `name = workerKey-slug-stage` rule (`internal/dispatch/build.go`). The reconcile helper inherits that contract; if the dispatch ever changes the shape, both ends move together.

## Mechanism — actions per drift class

The helper emits; the FO acts. The action table is:

| Class | FO action | When |
|-------|-----------|------|
| A (lingering) | `SendMessage(to=name, {"type":"shutdown_request"})`; remove from session memory | At idle + after each merge |
| B (superseded) | Same as A on the loser(s) | At supersede moment + idle sweep |
| C (un-advanced PR) | Run the normal merge-completion flow (advance to `done`, archive, worktree-remove); see `## Merge and Cleanup` | At idle |
| D (stale branch) | `git -C {worktree} pull --rebase origin next` from worktree; if conflict, halt + surface | At idle |
| E (stale local main) | `git fetch origin next && git reset --hard origin/next` on local main + `cd {repo} && go build -o spacedock ./cmd/spacedock` | At boot + after each merge |

Class A/B shutdowns are **cooperative best-effort**, matching the Degraded-Mode sweep semantics already in the runtime adapter — fire-and-forget, ignore failures. Class C/D/E are deterministic git/gh operations whose exit codes are the proof.

## Riskiest unverified mechanism — the spike

The design's load-bearing unknown is **`SendMessage` to a zombie or absent name**. Three behaviors are possible:

- (i) returns success silently (no-op against in-memory registry, FO never sees it failed) — A/B classes would still leak,
- (ii) returns a typed error the FO can observe (`"member not found"`) — A/B classes can be retried or escalated,
- (iii) blocks until timeout — A/B sweep stalls the FO event loop.

Implementation MUST run a one-shot spike before the rest of the helper lands: spawn an ensign, immediately SendMessage shutdown_request, then SendMessage shutdown_request a second time to the now-dead name. Record which of (i)/(ii)/(iii) holds. If (i), augment the helper to confirm via post-sweep config.json re-read (the member's `name` should be gone from `members[]`); if (ii), surface the error to the captain; if (iii), wrap the shutdown loop in a goroutine with a hard timeout. The result goes into the implementation entity body before the rest of class A/B lands.

**Spike result (recorded at implementation):** behavior (i) holds. On Claude Code, `SendMessage(to="spacedock-ensign-this-name-does-not-exist-on-purpose-for-spike", message={"type":"shutdown_request", ...})` returned `{"success":true, "request_id":"shutdown-...@<name>", "target":"<name>"}` synchronously, with no error, no block, no crash. The runtime accepts the message against any name shape and returns a request_id regardless of whether the target is a live, dead, or never-existed member. Consequences for the design:
- The FO's teardown loop (Class A/B actions in the action table) is safe to call fire-and-forget on Claude Code — no exceptions to swallow, no timeout to bound, no synchronous wait that could stall the event loop. The existing prose at `claude-first-officer-runtime.md` line 185 (`SendMessage(shutdown_request) is cooperative — do NOT send to dead or unresponsive ensigns`) describes the same observation from the FO-author's perspective.
- Because SendMessage gives no signal that the shutdown took effect, the helper's Class A on the next idle sweep is the only authoritative confirmation that a name was actually torn down (the post-sweep `config.json` re-read the entity body anticipated). Lingering names that "should" have died will show up again as Class A, and the FO will re-emit shutdown_request — convergent, not exception-driven.
- No goroutine timeout wrapper is needed (rules out (iii)); no captain-escalation glue is needed for typed errors (rules out (ii)). The shutdown loop is the simplest possible code.

The other mechanism candidates do NOT need a spike: (b) `members[]` schema is already proven by `internal/claudeteam/contextbudget.go` reading the same fields today, and a `jq` against the live team config in this session confirmed `agentType` distinguishes ensigns from `team-lead`; (c) `gh pr view --json state` is the same call the FO already uses to advance merged PRs; (d) `git rev-list --count` is git's documented contract; (e) `git reset --hard origin/next` + `go build` are bog-standard. **No spike needed for (b)/(c)/(d)/(e); proven mechanisms: team `config.json` members[] shape (read by contextbudget.go), `gh pr view --json state` (already in the event loop), `git rev-list --count` and `git reset --hard` (documented git), `go build -o`.**

## Integration surface — one new CLI command + targeted contract steps

The choice is **both**: a new code helper bears the detection logic; small contract additions wire it into the FO event loop. Rationale: the detection MUST be a code gate (a prose-only rule keys off memory — exactly the bug we are fixing), AND the FO must know WHEN to call it (boot/idle/post-merge) and HOW to act on each drift class. Neither half can be omitted.

### Helper — `spacedock dispatch reconcile`

Lives in `internal/dispatch/reconcile.go`, registered in `internal/dispatch/dispatch.go` alongside `context-budget` and `list-standing`. Same `Run` router, same `--workflow-dir` flag convention.

**Argv:**
```
spacedock dispatch reconcile --workflow-dir DIR [--team-name NAME] [--repo-root DIR] [--include classes]
```
- `--workflow-dir` (required) — absolute path to the workflow directory.
- `--team-name` (optional) — overrides the default discovery (which scans `~/.claude/teams/*/config.json` for one whose `members[]` contains a `spacedock:ensign` AND whose `leadSessionId` matches the current session). Tests pass this explicitly.
- `--repo-root` (optional) — defaults to `git rev-parse --show-toplevel` from `--workflow-dir`. Tests pass a fixture path.
- `--include` (optional, comma-separated subset of `A,B,C,D,E`; default all) — lets the FO scope a sweep (e.g. `--include A,B` at supersede time, full sweep at idle).

**Output shape (stdout):**
```json
{
  "command": "reconcile",
  "team_name": "spacedock-v1-dev-...",
  "drift": [
    {"class":"A","name":"spacedock-ensign-7h-implementation","slug":"release-notes-local-summary","reason":"entity archived"},
    {"class":"D","slug":"yaml-parser-migration","worktree":".worktrees/spacedock-ensign-yaml-parser-migration","behind":3,"reason":"branch behind origin/next by 3"}
  ]
}
```

**Exit codes:**
- `0` — sweep ran; `drift[]` is empty (clean) OR populated (FO acts on it). Either way the sweep itself succeeded.
- `1` — sweep could not run (team config missing, workflow dir invalid, git not in repo). Stderr names the cause. The FO surfaces to the captain.
- `2` — usage error (missing required flag, unknown `--include` class).

Drift-present is NOT a non-zero exit — the FO consumes `drift[]` and decides per class. Mixing "sweep failed" and "drift found" into one exit code would force the FO to parse stderr to disambiguate, which is exactly the kind of brittle coupling we want to avoid.

The helper does NO writes — no `status --set`, no `git reset`, no `SendMessage`. It is a pure reader emitting a structured report.

### FO contract additions

Three small edits to `skills/first-officer/references/first-officer-shared-core.md` and the Claude runtime adapter:

1. **`## Merge and Cleanup` — step 9.5 (teardown agents).** Between worktree-remove (current step 9) and the next merge boundary, insert: "Derive the entity's agent cohort from the live team roster (every `spacedock:ensign` whose name matches `decompose(name).slug == {this entity's slug}`). SendMessage `shutdown_request` to each. Drop them from session-memory tracking." This codifies the teardown rule; the reconcile helper's Class A is the behavioral backstop that catches a missed step.
2. **`## Feedback Rejection Flow` (and `## Completion and Gates` reuse-or-fresh) — supersede-shutdown.** When the FO decides "fresh dispatch" because of cycle increment OR feedback rework, shut down the prior cohort BEFORE dispatching the new one (in its own message, per the existing sequencing rule). Class B is the backstop.
3. **`## Event Loop` (Claude runtime) — boot + idle + post-merge sweep.** Add a new step 0 to startup (after split-root pull-on-boot, before first dispatch): "Run `spacedock dispatch reconcile --workflow-dir {wd}`; act on each drift class per the action table; report a one-line summary to the captain." Add the same call to the idle-hook fan-in (step 4) and to post-merge cleanup. Class E is what catches the "stale local main + stale binary" boot scenario.

The class-A/B prose steps are not the proof — they are how the FO USES the helper. The helper's existence is the gate. AC-1 verifies the helper detects each class on a fixture; AC-2 verifies the prose steps are present AND the helper would flag a contract violation if the prose were skipped.

## Acceptance criteria

**AC-1 — `spacedock dispatch reconcile` detects all five drift classes from a synthetic fixture.** Given (a) a fake team `config.json` with members[] containing a `team-lead` (exempt), an ensign for a live entity (clean), an ensign for an archived entity (A), two ensigns for the same slug+stage (B winner+loser), (b) a workflow tree with active + archived entity dirs each containing a minimal index.md whose frontmatter encodes the expected state, (c) a fake `gh` shim returning MERGED for the PR-pending entity (C), and (d) a worktree fixture where HEAD is behind `origin/next` (D) and local main is ahead of `origin/next` (E), the helper emits a `drift[]` containing exactly one entry of each of A, B, C, D, E with the expected `slug`/`name`/`reason` fields. A "flip" test mutates the fixture (archive the live entity → its agent now appears as Class A) and re-runs to confirm the classification is data-driven.
Verified by: `internal/dispatch/reconcile_test.go` — a Go test owning the fixture directory, invoking the helper's `Run` function directly (not via subprocess), asserting the JSON output. Targets 200-300 LOC of test fixture + assertions; runs in `go test ./internal/dispatch/...` in under 2s.

**AC-2 — the FO contract carries the teardown + supersede + sweep steps.** `skills/first-officer/references/first-officer-shared-core.md` `## Merge and Cleanup` contains a teardown-at-terminal step naming the helper as the gate; the same file (or the Claude runtime adapter) contains a supersede-shutdown step in the reuse-or-fresh path; the Claude runtime adapter's `## Event Loop` calls `spacedock dispatch reconcile` at boot + idle + post-merge. Each prose step is paired with a behavioral backstop in the helper (Class A for teardown, B for supersede, C/D/E for the sweep).
Verified by: (a) AC-1's fixture covers archived-entity-with-live-agent (Class A) and same-slug-same-stage cohort (Class B), so the helper would flag the contract violation if the prose step were skipped — this is the load-bearing gate. (b) A `prose_neutrality_test.go`-style presence check confirms the three contract paragraphs exist and name the helper command verbatim (proof-at-claim-level for text-is-the-claim wording, NOT a behavioral assertion). Item (a) is what makes the AC pass; item (b) is bookkeeping.

**AC-3 — local-main + binary reset discipline is encoded as Class E and runnable.** The reset-`main`→`origin/next` and rebuild-`./spacedock` discipline (this session's lesson) is detected by Class E and executed by the FO as a deterministic shell sequence whose exit code is the proof.
Verified by: AC-1's fixture E case proves detection; an integration test in `internal/dispatch/reconcile_e_test.go` sets up a tiny git fixture with local main ahead of origin/next, runs the helper, then runs the prescribed `git fetch && git reset --hard origin/next` and asserts main now points at origin/next (exit 0 + `rev-parse main == rev-parse origin/next`). The binary rebuild is verified by exit code of `go build -o /tmp/spacedock-test ./cmd/spacedock` from the fixture — the build either succeeds or fails, no prose proof needed.

## Test plan

Three test files, all in-process Go tests (no subprocess, no live CI):

1. **`internal/dispatch/reconcile_test.go`** — AC-1's five-class fixture. Owns a tmpdir tree mimicking a workflow + team config + a stub `gh` (interface-injected, not a real binary). ~300 LOC. Runs in `go test` under 2s.
2. **`internal/dispatch/reconcile_decompose_test.go`** — table-driven test for `decompose(name)`: known-stage stripping, cycle-suffix stripping, custom-stage stripping from a README fixture, ambiguous-name rejection. ~80 LOC.
3. **`internal/dispatch/reconcile_e_test.go`** — AC-3's local-main + git-reset integration. Uses `t.TempDir()` + `git init` + real `git` calls (no network). ~120 LOC.

The spike for the SendMessage-to-zombie question is NOT in this test plan — it runs once in implementation as throwaway code, the result is recorded in the implementation entity body's design notes, and the FINAL code's shutdown path reflects whichever of (i)/(ii)/(iii) was observed. No regression test for the spike itself — it informs the design, doesn't gate it.

Estimated implementation cost: 1 implementer-cycle (300-500 LOC helper + 500-600 LOC tests + 3 small contract edits). Validation is offline-only (no live CI needed — every test is in-process Go).

## Notes
- Scaffolding (FO contract) + new `internal/dispatch/reconcile.go` + tests — goes through a worktree on a `spacedock-ensign/ensign-lifecycle-reconcile` branch. Coordinates with the `2a` opt-in-guard pattern (a helper the FO runs, not a universal mandate — though Class A/B teardown IS universal since lingering agents waste tokens in every workflow).
- The agent-teardown is the PRIMARY ask (Classes A + B); the PR/branch-drift surfacing (Classes C + D + E) is the same idle-reconcile sweep generalized — folds in cleanly because the helper already opens all three state sources for A/B anyway. Adds ~50 LOC of helper code and ~150 LOC of test fixture, well below the bloat threshold.
- Captures the dogfooding lesson: this session's FO failures (leaked agents, reactive rebase) are precisely a memory-based-discipline ceiling; the fix is to compute from state.
- Out of scope (deferred to separate entities): auto-restart of the FO on context exhaustion, rewriting the team-registry desync recovery ladder, and any cross-team reconciliation (the helper scopes to one team at a time).

## Stage Report: ideation

- DONE: The reconcile mechanism is concretely designed: (a) what state IS authoritative, (b) the sweep algorithm that detects drift, (c) the actions taken per drift class. The riskiest unverified mechanism is exercised first or recorded as 'no spike needed' with the proven mechanisms relied on.
  Authoritative state = team `config.json` members[] + entity frontmatter + git refs (Mechanism section). Sweep = five disjoint drift classes A/B/C/D/E with pseudocode (Mechanism — sweep algorithm). Actions = the per-class table (Mechanism — actions per drift class). Riskiest unverified mechanism is `SendMessage`-to-zombie semantics — flagged as the implementation-stage spike; one-shot run before the rest of A/B lands. Proven mechanisms for the other candidates recorded: `members[]` schema via existing `internal/claudeteam/contextbudget.go`, `gh pr view --json state`, `git rev-list --count`, `git reset --hard`, `go build`.
- DONE: ACs are entity-level + each has a `Verified by:` clause naming a runnable check outside the entity body. Specifically each of (teardown-at-terminal, supersede-shutdown, reset-main→origin/next + rebuild discipline, lingering-ensign detection) must be guarded by a real test or a deterministic command's output/exit code, not by prose.
  AC-1 verified by `internal/dispatch/reconcile_test.go` (in-process Go test, five-class fixture + flip test). AC-2 verified by the AC-1 fixture covering archived-entity-with-live-agent (Class A) and same-slug-same-stage cohort (Class B) — those are the load-bearing behavioral gates; a presence check on the contract paragraphs is bookkeeping. AC-3 verified by `internal/dispatch/reconcile_e_test.go` (real `git init` fixture, asserts `rev-parse main == rev-parse origin/next` after the reset) + `go build -o /tmp/spacedock-test ./cmd/spacedock` exit code. All four required behaviors map to runnable checks outside the entity body.
- DONE: The integration surface is named: a new `spacedock dispatch reconcile` helper command vs an FO contract step vs both, with the choice rationale. If a new command, its argv + output shape + exit-code semantics are in the body.
  Choice = both (Integration surface section). Helper at `internal/dispatch/reconcile.go`, argv `--workflow-dir DIR [--team-name NAME] [--repo-root DIR] [--include classes]`, output `{"command":"reconcile","team_name":"…","drift":[{"class":"A|B|C|D|E", …}]}`, exits `0` (sweep ran, drift-present-or-empty), `1` (sweep failed: missing config / invalid wd / not in repo), `2` (usage). Rationale: detection must be code (the prose-only-rule failure mode IS the bug we are fixing); the FO must know when/how to call it (boot, idle, post-merge) and how to act per class. Neither half can be omitted.

### Summary

Sharpened the body from a sketch into a concrete design: three authoritative state sources (team `config.json`, entity frontmatter, git refs), a five-class disjoint-classification sweep algorithm (A lingering / B superseded / C un-advanced PR / D stale branch / E stale local main), a per-class action table the FO executes, a documented spike for the one truly unverified mechanism (`SendMessage`-to-zombie semantics), a `no spike needed` line listing the four mechanisms proven by existing code/contract, three entity-level ACs each verified by a named in-process Go test or deterministic command exit code, and a chosen integration surface = new `spacedock dispatch reconcile` helper PLUS three small FO-contract additions wiring it into boot / idle / post-merge / teardown / supersede paths. The argv + output JSON shape + exit-code semantics are pinned in the body so the implementer does not re-litigate them.

## Stage Report: implementation

- DONE: `spacedock dispatch reconcile` ships at internal/dispatch/reconcile.go with the documented argv (--workflow-dir, --team-name, --repo-root, --include), JSON output shape ({command, team_name, drift[]}), exit codes (0 sweep ran; 1 setup failure; 2 usage), and the 5-class disjoint detection (A lingering / B superseded / C un-advanced PR / D stale branch / E stale local main). All 5 classes detect cleanly on the synthetic fixture; the SendMessage-to-zombie semantics impl-stage spike is run and its result is recorded in the entity body's design notes before A/B's teardown path lands.
  Helper at .worktrees/spacedock-ensign-ensign-lifecycle-reconcile/internal/dispatch/reconcile.go; routed in dispatch.go's switch. Live exercise against the team confirms argv + JSON shape + exit codes (`reconcile {…} drift=[]` exit 0; missing flag exit 2; missing dir exit 1; bad include exit 2). Spike result recorded inline in "Riskiest unverified mechanism" — behavior (i), silent success. Commit 31e27026.
- DONE: Three Go tests are green: reconcile_test.go (five-class fixture + flip; ~300 LOC); reconcile_decompose_test.go (name decomposition table-driven; ~80 LOC); reconcile_e_test.go (real git init fixture asserting `rev-parse main == rev-parse origin/next` after reset + go build exit code; ~120 LOC). `go test ./...` clean.
  reconcile_test.go = 470 LOC (5 tests: TestReconcileFiveClasses, TestReconcileFlipReclassifies, TestReconcileIncludeScope, TestReconcileMissingWorkflowDir, TestReconcileUsageError). reconcile_decompose_test.go = 132 LOC (3 tests, 23 subcases). reconcile_e_test.go = 138 LOC (2 tests, real git init + go build). `go test ./...` across all 10 packages: ok, ~50s wall clock.
- DONE: The 3 FO-contract additions are written verbatim into the scaffolding: skills/first-officer/references/first-officer-shared-core.md ## Merge and Cleanup gets a teardown-at-terminal step naming `spacedock dispatch reconcile`; the same shared-core (or claude-first-officer-runtime.md) gets a supersede-shutdown step in the reuse-or-fresh path; claude-first-officer-runtime.md ## Event Loop gets `spacedock dispatch reconcile` calls at boot + idle + post-merge. Each prose addition is load-bearing per the entity's AC-2 framing — the helper is the behavioral backstop, the prose is how the FO uses it.
  shared-core ## Merge and Cleanup grew a new step 10 (teardown-at-terminal) naming the helper and decompose() rule. shared-core ## Completion and Gates "If fresh dispatch" path grew a new supersede-shutdown paragraph (the reuse-or-fresh decision point per the entity's request). claude-first-officer-runtime.md ## Event Loop grew a new step 0 (boot + post-merge sweep) and step 4 grew an idle re-run of the same sweep. Each cites the helper command + Class A/B/C/D/E backstop verbatim.

### Summary

Built and shipped `spacedock dispatch reconcile` as a pure-reader CLI helper that emits a five-class drift report (A lingering / B superseded / C un-advanced PR / D stale branch / E stale local main) from team config + entity frontmatter + git refs. Three Go test files (~720 LOC total) cover the five-class fixture with a flip, the name decomposer, and an E-class integration with real git init + `go build` exit-code proof; `go test ./...` is clean. The `~/.claude/teams` read lives behind `claudeteam.LoadReconcileTeam` so the host-neutrality scan over `internal/dispatch` stays green. Three FO-contract additions wire the helper into Merge-and-Cleanup (teardown), Completion-and-Gates reuse-or-fresh (supersede-shutdown), and the Claude Event Loop (boot / idle / post-merge sweep) — each prose step names the helper command and is paired with the class the helper would detect if the step were skipped. The SendMessage-to-zombie spike returned behavior (i) — silent success — and the design notes capture the consequence: the helper's Class A on the next sweep is the only authoritative confirmation of teardown.

## Stage Report: validation

- DONE: Every AC-N from ## Acceptance criteria is reproduced from a runnable check OUTSIDE the entity body at HEAD 4a70ce83: AC-1 verified by `internal/dispatch/reconcile_test.go` (5 tests including TestReconcileFiveClasses + TestReconcileFlipReclassifies); AC-2 verified by AC-1's behavioral backstop (Class A + B detect contract violations) AND a prose-presence check that the 3 contract paragraphs cite the helper verbatim; AC-3 verified by `internal/dispatch/reconcile_e_test.go` (real git init fixture asserts `rev-parse main == rev-parse origin/next` after reset + go build exit code).
  All three ACs verified at HEAD 4a70ce83. `go test ./internal/dispatch/... -run 'Reconcile|Decompose' -v` = 26/26 PASS (TestReconcileFiveClasses, TestReconcileFlipReclassifies, TestReconcileIncludeScope, TestReconcileMissingWorkflowDir, TestReconcileUsageError, TestReconcileEDetectsAndResetAdvancesMain, TestReconcileEGoBuildIsRunnable, TestDecomposeCanonicalStages + 12 subcases, TestDecomposeWithCustomWorkflowStages + 5 subcases, TestParseIncludeRoundTrip). AC-1's TestReconcileFiveClasses asserts exactly one of A/B/C/D/E with documented slug/name/reason; the flip test mutates the fixture (archive `alive`) and the alive ensign reclassifies to Class A — data-driven, not hard-coded. AC-3's E_test runs a real `git init` repo, sets `refs/remotes/origin/next`, makes main 1-ahead, asserts `rev-parse main != rev-parse origin/next` pre-reset, runs `git reset --hard origin/next`, asserts equality post-reset; the sibling go-build test compiles `./cmd/spacedock` to a tmpdir with exit code as the proof.
- DONE: The 3 FO-contract additions are present at the documented locations: shared-core ## Merge and Cleanup step 10 (teardown), shared-core ## Completion and Gates supersede-shutdown paragraph, claude-runtime ## Event Loop step 0 (boot+post-merge sweep) + step 4 (idle re-run). Each cites the helper command verbatim AND names the Class (A/B/C/D/E).
  Verified at: first-officer-shared-core.md:240 "Teardown agents at terminal" — cites ``spacedock dispatch reconcile`` and names **Class A (lingering)** backstop. first-officer-shared-core.md:159 "Supersede-shutdown" — lands in the `## Completion and Gates` reuse-or-fresh path (line 157 "If fresh dispatch"), cites the helper and names **Class B (superseded)** backstop. claude-first-officer-runtime.md:248 step 0 "State-keyed reconcile sweep" — names all five Classes A/B/C/D/E with per-class action, lists boot/idle/post-merge moments; step 4 (line 252) re-runs the sweep at idle. Markdown is syntactically valid (no broken headings; numbered list continuity preserved — `## Merge and Cleanup` runs 1→10 cleanly; `## Event Loop` runs 0→4). No malformed links. Diff stat: shared-core +3/-0, claude-runtime +2/-1.
- DONE: Full repo `go test ./...` green at HEAD 4a70ce83 (765/765 in 12 packages); the host-neutrality scan over internal/dispatch stays green (the ~/.claude/teams read lives behind claudeteam.LoadReconcileTeam); the SendMessage-to-zombie spike result is recorded inline in the entity body's `## Riskiest unverified mechanism` section (behavior (i) — silent success — observed at entity body lines 107-110).
  `go test ./...` (full repo) = 765 passed in 12 packages, no failures. `grep '\.claude/teams' internal/dispatch/` returns only (a) golden text in `testdata/golden/build-crossproduct-...txt` and (b) a doc comment in `reconcile.go:29` referencing `claudeteam.LoadReconcileTeam` — no actual host-path read in `internal/dispatch`; the read lives at `internal/claudeteam/reconcile.go:43-87` behind `LoadReconcileTeam`. Spike result confirmed at entity body lines 107-110: "behavior (i) holds. On Claude Code, `SendMessage(to=\"spacedock-ensign-this-name-does-not-exist-on-purpose-for-spike\", ...)` returned `{\"success\":true, ...}` synchronously, with no error, no block, no crash." Lines 108-110 carry the three design consequences (fire-and-forget safe; Class A on next sweep is the authoritative confirmation; no goroutine timeout needed).

### Summary

PASSED. All three ACs verified by runnable checks outside the entity body at HEAD 4a70ce83. AC-1: `internal/dispatch/reconcile_test.go` (5 tests, 26 subcases incl. the flip) emits exactly one of each A/B/C/D/E from the synthetic fixture and the flip reclassifies. AC-2: the three FO-contract additions land in the documented sections (shared-core ## Merge and Cleanup step 10, shared-core ## Completion and Gates supersede-shutdown, claude-runtime ## Event Loop step 0 + step 4), each citing `spacedock dispatch reconcile` verbatim and naming its per-Class backstop; markdown is syntactically valid. AC-3: `internal/dispatch/reconcile_e_test.go` asserts `rev-parse main == rev-parse origin/next` after the reset on a real git fixture, plus `go build -o tmp ./cmd/spacedock` exit-code proof. Live exercise of `/tmp/spacedock-validate dispatch reconcile --workflow-dir docs/dev` against the running session emitted a valid envelope with one Class D drift (sibling `codex-live-ci` worktree 1 behind), confirming the helper runs end-to-end. `go test ./...` is 765/765 green across 12 packages; the host-neutrality discipline holds (the only `~/.claude/teams` read lives at `internal/claudeteam/reconcile.go` behind `LoadReconcileTeam`). The SendMessage-to-zombie spike result (behavior (i) — silent success) is recorded inline at entity body lines 107-110 with three design consequences. Implementation came in at ~1450 LOC vs the ideation's 800-1100 estimate; the overshoot is honest — the helper's 627 LOC distributes across 5 class detectors (28/39/27/31/13) + a 65-LOC orchestrator + a 42-LOC decomposer + 2 injection seams + JSON marshalling, no single function bloated.
