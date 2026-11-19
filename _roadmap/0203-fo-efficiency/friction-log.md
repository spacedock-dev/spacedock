# 0203 — As-if-implemented dogfood: friction log

Captain directive (2026-06-13): operate this FO session **as if** j9 (lazy-TeamCreate +
shallow-boot-then-greet + split contract), #344 (clean budget probe), and T3 (slimmed refs)
are already live in the FO contract; **file friction as surfaced** — every point where the
as-if behavior rubs against the real, unsplit contract / binary.

Each entry maps to the sprint lever or AC it validates, so the log doubles as live evidence
for j9's implementation (especially AC-1/AC-6) and T3's cut-list.

**Class key** — `KNOWN`: already in sprint scope, being addressed; the as-if run only
re-measures it live (corroboration, not discovery). `SURFACED`: new signal the as-if run
produced that no sprint task currently owns.

| ID | Phase | Class | Maps to | Severity |
|----|-------|-------|---------|----------|
| FR-1 | boot | KNOWN | j9 P2 lazy-TeamCreate (~89k) | confirms lever (live) |
| FR-2 | boot | KNOWN | j9 P1 split / P3 shallow-boot (~16k) | confirms lever (live) |
| FR-3 | boot | KNOWN | j9 P3 S4-slim README-frontmatter (~7.7k) | confirms lever (live) |
| FR-4 | gate | SURFACED | entity-body cycle-report weight — no sprint owner | new seed |
| FR-5 | boot | SURFACED (incidental) | stale README mods clause; also pro-FR-3 evidence | FO-fixable |
| FR-6 | dispatch | SURFACED | reconcile Class-E destructive vs main-is-trunk; orphan debt | needs decision |
| FR-7 | boot/dispatch | SURFACED (meta) | simulability split — P2 simulable, P1 bootstrap-locked | refines AC-6 scope |
| FR-8 | dispatch | SURFACED | per-transition manual state commit+push (no tool-owned atomic) | "cheap to run" seed |
| FR-9 | dispatch | SURFACED | the #344 budget-probe bug spuriously blocked reuse of the FO's own ensign | proves #344 load-bearing |
| FR-10 | filing/merge | SURFACED | shipped `spacedock new` untaught in the live contract; j9's stale base would clobber #362 | cheap-to-run + j9 pre-merge blocker |
| FR-11 | run-cost | SURFACED | dispatched subagents pinned to 5m cache (top-level=1h) → full-prefix re-cache storm | host-support: 1h TTL (~$14–73/run) |
| FR-12 | run-cost | SURFACED | RTK-hook rewrite-echo (~29.6k) + empty task_reminder nudges = ~30k prefix tax | RTK/harness config |
| FR-13 | dispatch coord | SURFACED | FO re-tasked a worker after dispatching its reviewer → dirty-worktree race | serialize fix→review |

**Bottom line:** FR-1/2/3 are the three documented cost levers — the dogfood *re-measured* them,
it did not find them. The genuinely new signal is FR-4, FR-6, and especially **FR-7** (what the
as-if assumption can and cannot prove). FR-5 is an incidental doc bug that doubles as evidence for
the FR-3 fix.

---

## FR-1 — Eager team-harness load at boot

**As-if expectation (j9 P2):** the FO creates no team and loads no team-harness at boot;
`Skill(spacedock:using-claude-team)` + `TeamCreate` defer to the first team-mode dispatch.

**What actually happened:** the live Claude runtime adapter `## Team Creation` (line 7) says
*"At startup (after reading the README, before dispatch), invoke … `Skill(spacedock:using-claude-team)`."*
Following the real contract, I loaded the full team-harness skill at boot. All three sprint
entities sit at gates with **zero dispatch pending**, so both the harness load and the ~89k
`TeamCreate` cache-creation were fully avoidable on this boot.

**Mitigated:** I did **not** call `TeamCreate` (held per the as-if directive), so the ~89k spike
itself was avoided. The harness *skill read* still happened — that is the contract-read half of
the cost the P1 split removes.

**Evidence value:** direct live confirmation of j9's headline premise — the one clause at
`claude-first-officer-runtime.md:7` is exactly what P2 retargets to "before the first team-mode
dispatch." A greet-and-stop boot pays the team prefix for nothing.

## FR-2 — Full FO contract read at boot

**As-if expectation (j9 P1/P3):** boot reads only the slimmed boot-resident core (~4.5k);
the dispatch and merge modules are deferred references, not read until first dispatch /
terminalization.

**What actually happened:** boot read the full `first-officer-shared-core.md` (321 lines) and
the full `claude-first-officer-runtime.md` (225 lines) — Dispatch, Completion-and-Gates
reuse-conditions, Merge-and-Cleanup, Worktree ownership, Event Loop, Degraded Mode,
Context-Budget probe, Terminal Teardown. A greet-and-stop boot uses none of it.

**Evidence value:** confirms the ~16k "defer contract reads at greet" lever and the
~11k / ~70%-of-contract-read-cost cut P1 claims. The boot-resident vs deferred partition j9
P1 specifies (Startup/Status/IDs/Write-Scope/gate-spine resident; Dispatch/Merge deferred)
matches exactly which sections I touched at boot vs. won't touch until a gate is approved.

## FR-3 — Full workflow-README body read at boot

**As-if expectation (j9 P3 S4-slim):** boot reads only the README **frontmatter** (~175 tokens:
entity-label, stage names/ordering, gate/terminal flags) and defers the body.

**What actually happened:** I read the full `docs/dev/README.md` (315 lines: per-stage prose,
proof policy, Runtime-Live-CI docs, task template ~8.1k) to obtain the stage taxonomy + gate
flags — all of which are in the frontmatter. Confirmed the boot JSON does **not** carry
labels/stage-taxonomy, so a frontmatter read is genuinely still needed (exactly as P3 states) —
but the ~7.7k body was not.

**Evidence value:** confirms the ~7.7k README-body lever and validates P3's claim that the boot
JSON cannot substitute for the frontmatter (only the body defers).

## FR-4 — Entity-body cycle-report weight at gate presentation (tangential)

**Observation:** presenting j9's ideation gate required reading a 258-line / ~27k-token entity
body — four accreted ideation-cycle stage reports plus "(Superseded by cycle 2)" historical
annotations. This is a gate action, not boot, and entity bodies are not contract, so it is
**not** a j9 lever. But it is a real recurring cost: a fully-iterated entity body is paid in
full at every gate.

**Does the gate need only a section? — CONFIRMED, two halves.** Checked against what my j9
gate presentation actually consumed: the `## Proposed approach` heading (for the chosen-direction
line), the *latest* `## Stage Report` (cycle 4, for the checklist), `## Acceptance criteria` (for
the AC cross-check), and the external `staff-review.md` (reviewer findings). The other ~80% of the
body — spike-status, implementation-pin, full test-plan, and the THREE superseded cycle reports —
went unused.

- **Half A — self-inflicted, already covered by existing discipline.** The contract's Probe-and-
  Ideation discipline already says *"Prefer Grep over Read for targeted entity-body inspection.
  Anchor on heading … Read only when you need the full text."* I full-Read the body (paged twice)
  instead of Grepping `## Stage Report` / `## Acceptance criteria` / `## Proposed approach`. That
  cost is on me, not the design — the gate provably needs only a few sections.
- **Half B — the real seed (out of 0203 scope).** Even a disciplined section-targeted reader pays
  to *disambiguate* four accreted `## Stage Report: ideation (cycle N)` sections plus their
  "(Superseded by cycle 2)" annotations to find the latest. Seed: compact superseded cycle reports
  once a stage is gate-approved (or have the body carry only the live report + an archive link).
  Non-blocking, filed.

## FR-5 — Workflow README stale on registered mods (FO-fixable)

**Observation:** `docs/dev/README.md:35` asserts *"No PR merge flow, mods, or lifecycle hooks
are in scope for this bootstrap workflow."* The boot JSON contradicts it:
`startup=[comm-officer, pr-merge]`, `idle=[pr-merge]`, `merge=[pr-merge]`,
`shutdown=[comm-officer]`. The README body is stale.

**Note:** the README is an FO-owned process doc (in write scope), so this is directly fixable —
but I will not churn it mid-greet. It also compounds FR-3: reading the stale body cost ~8k *and*
returned wrong information that the frontmatter would not have.

## FR-6 — Reconcile Class-E remedy is destructive against a main-is-trunk repo

**Observed at the first-dispatch reconcile sweep** — which, under as-if-j9, is the *correct*
deferred timing (the sweep travels with the dispatch module, C3, not boot; I ran it only when
creating the team to dispatch j9, not at greet). `spacedock dispatch reconcile` reported:

- **1× Class E:** local `main` is **27 commits *ahead*** of `origin/next`; the prescribed remedy
  `reset --hard origin/next` would **destroy** those 27 commits (including the 0203 artifacts this
  session runs on). This repo uses `main` as the working trunk, intentionally ahead of `next` —
  the Class-E remedy assumes main *mirrors* next. **Declined the reset** (destructive; captain
  rules forbid discarding commits unasked). main-ahead only means j9's worktree gets the latest
  code, which is correct.
- **6× Class D:** dead orphan worktrees from prior sprints (`agy`, `mdschema`,
  `pre-cut-audit-0199`, `prefer-new`, `state-sync`, `survey`), each behind `origin/next` by 13.
  Not 0203 entities; j9 uses a fresh worktree. Did **not** pull-rebase them (no live owner →
  pointless churn). Cleanup debt.

**Independent of the sprint** — a latent reconcile-design / repo-branch-model mismatch; the as-if
timing only relocated *when* it was hit (first dispatch, not boot). Needs a captain decision:
(a) is `reset --hard origin/next` ever right here, or should Class-E be advisory when main is
*ahead*? (b) prune the 6 dead orphan worktrees?

**Captain ruling (2026-06-13):** `reset --hard origin/next` is **definitely not right** here.
Backlog seed for the reconcile classification logic: Class-E must not prescribe a destructive
reset when local main is *ahead* of `origin/next` (nothing to recover, commits to lose) — it
should be advisory, or scoped to the genuinely-behind/diverged case. Orphan-prune (b) still open;
low priority, not blocking.

## FR-7 — Simulability split: what the dogfood can prove, and the bootstrap limit (meta)

The as-if run's sharpest signal. Assuming the new behavior is only *possible* for the levers that
are pure behavior-timing; it is *impossible* for the lever that is the contract itself:

- **P2 lazy-TeamCreate — fully simulable, and demonstrated live this boot.** Holding `TeamCreate`
  to first dispatch needed zero code change — it is purely a trigger-timing choice the FO controls.
  This boot avoided the ~89k cache-creation entirely (no team existed until the j9 dispatch). Live
  proof of j9's "Needs the split? — no" for P2, and the strongest pre-validation the dogfood can
  give the headline lever.
- **P3 shallow-boot — partially simulable.** The deferrable ACTIONS I could and did defer
  (reconcile → first dispatch; human status-table render → never rendered for my own reasoning;
  mod-file reads → deferred). The READS I could not: Startup step 4 as written had me read the full
  README (FR-3), and I had to read the full contract to know the procedure at all.
- **P1 contract split — NOT simulable (the bootstrap paradox).** To behave as a slim-booted FO I
  must first read the contract that tells me how to behave — which is the very cost the split
  removes. The as-if cannot simulate its own enabler.

**Implication for the drive order + AC-6.** "Cheap lever first" is doubly right: P2 is the biggest
lever AND the only one this dogfood can prove standalone *right now*. The run pre-validates the
~89k (P2) lever strongly and the ~16k (P1) lever **not at all** — P1's saving is only measurable
*after* it ships, which is exactly what AC-6's post-ship live drive is for. This run is
corroboration for P2 and a scoping note on what corroboration is even possible.

**CLOSED by the formal live run (2026-06-13).** The `shallow-boot` scenario PASSED on the rotated
credential: greet-turn context **45,862** (< ~60k, ~14k margin), max pre-greet `cache_creation`
**16,889** (no ~89k spike), **zero `TeamCreate`** calls; AC-1/2/3/6 all green; 2/2 PASS. This
confirms FR-7's prediction exactly — the cheap lever clears `<60k` *while still carrying the full
un-split ~16k contract*, so P1's ~16k cut is **not load-bearing** for the headline. **Captain
decision:** P1 stays in 0.20.3 anyway — as the contract-cleanup goal, because T3 depends on it,
and because the work is ready + de-risked. The dogfood's informal P2 demonstration and the formal
measurement agree; FR-7 resolved.

## FR-8 — Every FO state transition needs a manual path-scoped commit + push

**Observed across this session's dispatches.** `spacedock status --set` mutates the entity
frontmatter but does NOT commit it (verified: after the j9 dispatch `--set`, the state-checkout
log head was unchanged until I committed manually). So each state transition costs the FO a
fixed sequence: `status --set` → `git -C {state} add {path}` → `git -C {state} commit -- {path}`
→ `git -C {state} push`. The shared core names a "preferred — tool-managed atomic state commits"
path *"when the status tool owns add+commit under a lock"* — but this build doesn't, so the FO is
always on the path-scoped fallback.

**Theme fit:** this is squarely the sprint's *other* half — "make the FO cheap to **run**." Every
dispatch/advance/merge boundary pays the 3-git-op tax. It is adjacent to the deferred
binary-simplification line (`p2` pr-complete, `vc` reconcile --act) but BROADER: those target the
merge/completion ceremony, whereas this is every `--set`. **Candidate:** a `status --set --commit`
(or tool-owned atomic commit+push) that collapses the transition to one call.

**Status:** surfaced, not blocking — the fallback works. Worth a backlog seed on the
binary-simplification line, not 0203 scope. Filed for the record.

## FR-9 — The #344 budget-probe bug bit the FO's OWN reuse decision (live)

**Observed when the captain asked me to reuse the implementation ensign for probe-prep.** The
reuse-condition-0 budget probe (`dispatch context-budget --name ...-implementation`) returned:

```
resident_tokens: 432415, model: <synthetic>, context_limit: 200000, usage_pct: 216.2,
reuse_ok: false, config_declared_model: claude-opus-4-8,
mixed_models_warning: "['<synthetic>','claude-opus-4-8'] — using smallest context window",
config_drift_warning: "team config requested claude-opus-4-8 but runtime is <synthetic>"
```

This is **the exact bug `#344` fixes** (the `<synthetic>` census skip + the 1M-window promotion):
the `<synthetic>` internal-turn entries made the probe pick the smallest (200k) window and emit
spurious `mixed_models`/`config_drift` warnings, computing `432,415 / 200,000 = 216%` → over
budget. The ensign is actually **opus-4-8 (1M window)**: `432,415 / 1,000,000 ≈ 43%` — under the
60% threshold, **genuinely reusable**. The un-merged probe nearly forced a context-losing fresh
dispatch of exactly the ensign the captain wanted reused.

**Resolution:** overrode the spurious `reuse_ok:false` with the true budget (43%) and reused the
ensign — transparently, not silently (the fail-safe guards against a *genuinely* over-budget or
*absent* reading; this is a *known-buggy false-positive* with an interpretable true value).

**Evidence value (load-bearing, in the FO's own loop):** the sprint's own `#344` fix is needed by
the FO's reuse machinery, demonstrated LIVE — the spurious warnings actively corrupt the reuse
decision, not just a log cosmetic. This is the strongest possible argument for shipping `#344` in
the batch: a held-back fix is currently degrading the running FO. (Staff review already noted a
detached audit found `#344`'s over-suppression guards load-bearing; this is a live confirmation.)

## FR-10 — Shipped `spacedock new` untaught in the running contract; j9's stale base would clobber it

**Surfaced by the captain asking "why are you using `--next-id` rather than `new`?"** I filed three
seeds this session (`rz`/`js`/`xf`) via the manual flow — `status --next-id` → hand-assemble
frontmatter → `Write` → commit — when `spacedock new <slug> [--folder] [--id-seed S] < body`
atomically mints the id + writes a stamped, valid entity in one call (existing since #242).

**Two distinct problems:**

1. **Cheap-to-run lever (untaught command).** `nd` (#362, "Prefer `spacedock new` over manual
   `--next-id`") MERGED to origin/main at **21:24 today** and taught `new` in the contract. I booted
   *before* 21:24, so the contract I loaded only taught the manual flow; and my local main is still
   behind origin/main, so the working contract still lacks it. `new` collapses the 4-step intake
   into one call (+ the path-scoped commit/push for split-root) — a direct "cheap to run" win the
   FO wasn't getting. Going forward: use `new`.

2. **j9 pre-merge BLOCKER (stale base).** #362 rewrote the exact two contract files
   (`first-officer-shared-core.md`, `claude-first-officer-runtime.md`) that **j9's P1 split also
   rewrote** — but j9's worktree branched off local main (`b6e7a6f3`), which predates #362, so j9's
   split never carried #362's `new` teaching (worktree greps **0** for `spacedock new`; its base
   also lacks #364). **Merging j9 as-is would clobber/conflict-with #362.** The validator's "no
   content dropped across the cut" is correct *relative to j9's stale base* — it cannot see content
   that was never in the base. **Fix: rebase j9 onto current origin/main and re-incorporate #362's
   `new` teaching into the split before merge** (recommended: now, in this feedback cycle, so
   re-validation runs once on the correct base). Same stale-base class as the post-flip trunk
   divergence (FR-6 / `sr`) — local main lagging origin/main is biting the deliverable, not just the
   tooling.

## FR-11 — Dispatched subagents are pinned to the 5-minute cache (top-level threads get 1h) → re-cache storm

**The #1 run-cost lever** (forensics over the implementation ensign's session jsonl, new Claude Code format).

**Finding (host behavior, proven):** the dispatched ensign's prompt-prefix cache is **100% `ephemeral_5m`,
0% `ephemeral_1h`**; the top-level FO session is **100% `ephemeral_1h`**. Both ran the *same* ~4-hour
window. The FO's work cadence has **14–45 min idle gaps** (the worker waits between dispatch
advancements) — every gap **exceeds the 300s 5m TTL but is under the 3600s 1h TTL**, so the host
re-writes the worker's full ~300k–526k prefix from cold on each gap. **12 full-prefix re-caches =
2.44M tokens = 38% of the run's 6.42M lifetime `cache_creation`** (~$14 de-duped, ~$73 gross at Opus
rates). The launcher writes **zero `cache_control`** anywhere, so the TTL is the host's Task/subagent
machinery, not ours.

**Evidence to attach to the host request:** parent session `787cc95f` = 100% `ephemeral_1h`, 5
re-caches / 313 turns; subagent `agent-a83260bb8104fd269.jsonl` = 100% `ephemeral_5m`, 14 re-caches /
665 turns — identical wall-clock window. Every large top-level session sampled (6/6) is 100% 1h; every
dispatched subagent is 100% 5m.

**Tips / potential remedies:**
1. **Primary (host-support ask to Claude Code):** apply `cache_control ttl="1h"` to long-lived
   dispatched/Task subagent threads (the host already does this for top-level threads), OR expose a
   per-dispatch cache-TTL knob the launcher can set in the `subagent_type` envelope
   (`internal/dispatch/build.go`). Recovers the ~$14–73/run.
2. **Launcher-side workaround (no host change, today):** nothing in the prefix changes during an idle
   gap — so cut the gaps. Keep dispatch cadence tight enough that a worker isn't parked >5m; for
   known-long waits, prefer a fresh dispatch over keeping a 500k worker warm (the re-cache premium can
   exceed the context-loss cost — see FR-8 / the reuse economics).
3. **Caveat:** 1h writes cost 2× vs 5m's 1.25×, so 1h only wins when the prefix is re-read within the
   hour (true here — all gaps <60m). For gaps >60m, the cadence/fresh-dispatch fix is the better play.
4. **Proof shape:** post-fix, assert the subagent run shows `ephemeral_1h > 0` and no post-idle-gap
   `cache_creation > 200k` where the gap is <3600s. Negative control = the current 5m-pinned shard.

## FR-12 — Injected-attachment prefix tax: RTK hook rewrite-echo + empty task_reminder nudges

**~30.4k tokens of steady-state prefix tax**, re-cached on every 5m expiry (FR-11 amplifies it ~12–14×)
— ~$8/run hard floor.

1. **`hook_success` — 141 entries ≈ 29.6k tok (98% of the tax), from the RTK hook.** Every Bash call
   fires a `PreToolUse:Bash` hook (`rtk hook claude`) that auto-rewrites the command and **echoes the
   full rewrite JSON back into context** — the entire `hookSpecificOutput.updatedInput.command` (avg
   840B, max 2.3KB). The irony: RTK exists to *save* tokens (it filters command **output**), but its
   PreToolUse hook re-injects the command **into** context on every call.
   - **Remedy:** make `rtk hook claude` context-silent — the rewrite applies via the hook's
     `updatedInput` regardless; the verbose `hookSpecificOutput` echo is the waste. Check whether the
     PreToolUse hook protocol has a `suppressOutput`-style field, or trim the rtk hook's stdout to
     minimal. Helps **every** agent, not just this sprint — an RTK-side config fix.

2. **`task_reminder` — 54 entries (~0.8k, the 2%), from the Claude Code HARNESS (not a hook).** The
   built-in "task tools haven't been used recently… consider using TaskCreate/TaskUpdate" nudge.
   Verified harness-native: the attachment carries **no source attribution** (unlike `hook_success`'s
   `hookName`), just `{type:"task_reminder", content:[], itemCount:0}` — the harness renders the text
   from it. Hardcoded; **no `settings.json` key or env var disables it.** It fires when the task/todo
   tools are unused and the list is empty — and the spacedock FO/ensign tracks state via
   `spacedock status`, not the Claude task tools, so it's a **perpetual false nudge** firing every ~N
   turns forever.
   - **Remedy:** host-support ask for a setting to suppress it (globally or for dispatched subagents
     that intentionally don't use the task tools). Low priority — the payload is small; the cost is
     the re-cache amplification, which FR-11's 1h-TTL fix removes anyway.

## FR-13 — Re-tasking a worker after dispatching its reviewer races the shared worktree

**Self-inflicted FO coordination slip (mine).** The FO sent the implementation ensign an AC-2
refinement (use the full `sonnet_teamdelete_hang` oracle) **after** already dispatching the validator
for the cycle-3 re-review. The ensign and validator **share the worktree**
(`.worktrees/spacedock-ensign-lazy-teamcreate-shallow-boot`), so the ensign's in-flight addendum left
the tree dirty (a partially-staged fixture relocation) exactly as the validator started — which the
validator reported as a START blocker.

**Mitigated cleanly, both disciplines held:** the validator validated the **committed HEAD on a
detached checkout** (never the dirty live worktree) and flagged the dirty tree for reconciliation —
the correct reviewer discipline. The ensign committed the addendum (`8328582e`, test-only); the FO
verified it was test-only (live behavior identical) and retargeted the cert to it with no live-run
waste.

**Lessons:** (1) **Serialize fix→review** — the FO should not re-task a worker while its reviewer is
running against the shared worktree; send refinements *before* dispatching the reviewer, or *after*
its verdict. (2) The reviewer's "**cert a committed ref on a detached checkout, never the live
worktree**" rule is the safety net that made the slip harmless — worth keeping as explicit validator
discipline. Two agents sharing one worktree is the structural root; a per-agent worktree or a
review-against-pushed-ref convention would remove the race entirely.
