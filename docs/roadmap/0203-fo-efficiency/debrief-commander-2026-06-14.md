# 0203 Commander drive — debrief (2026-06-14 → 06-15)

Interactive Commander session driving `0203-fo-efficiency` to (and past) the v0.20.3 cut. Team `spacedock-v1-dev-20260614-1120-fwrqijif` (claude, team mode).

## Headline outcome

**Boot-resident FO contract cut in half.** `v0.20.2` → `0.20.3-candidate`:
- Boot-resident (SKILL + shared-core + claude runtime adapter): **~26,000 → ~13,000 tok (−50%)** — the dispatch/merge machinery was deferred off the boot path (j9 split + T3 + 5e).
- Full tree: ~30,250 → ~31,270 tok (+1k structural split overhead, off the hot path).

## Shipped to main this session — 14 tasks

j9 (#365) · T3 (#367) · **sr** (#369) · 8e (#366) · **87** (#370) · **ec** (#372) · **5e** (#373) · **58** (#374) · **xf** (#375) · **gwt** (#376) · **7d** (#377) · **5h0** (#378) · **j1** (#379) · **r0n** (#380). (Bold = merged by this Commander.)

Highlights:
- **sr/87** — trunk single-sourced (`resolveTrunk` + `dispatch trunk`); the **shipped** pr-merge template is now config-driven (every new workflow inherits) + `reconciled-from-shipped`/`fo-realm` governance frontmatter on the dev instance.
- **58** — the real broad-search prevention: reframed `native_runner.go:79`'s no-workflow message from "pass --workflow-dir or run inside a workflow" (invited the hunt the captain manually halted) to a terminal "report and stop, do NOT search". Detector + registered live scenario as guard.
- **xf** — captain reversed mid-flight (on-demand polish → KEEP residency, relocate mechanism): `dispatch spawn-standing-all` composes the lifecycle; shed 4 contract subsections → 1 trigger line; residency proven live (356s).
- **gwt** — retired the session-long reconcile footgun: classE never `reset`s unpushed/diverged main; classD only touches owned worktrees.
- **5e** — ethos hoisted to SKILL.md + mod-block collapse + cascade.
- **5h0** — dropped README-substring prose-grep tests + closed the boundary-guard detector gap.
- **r0n** — `spacedock new` auto-discovers the workflow from root (downward fallback) + per-command `--help` (the exact friction this Commander hit filing j1).
- **j1** — docs.yml: dropped the failing Pages deploy-API, push `site/` to `gh-pages` directly (the 404 on every main push). Operator follow-up: set repo Pages source → "branch: gh-pages" to serve.

## In flight at debrief — 7e (LEAVE FOR NEXT SESSION)

**7e (headless-dispatch-mode-intent)** is at `status=implementation` (worktree `.worktrees/spacedock-ensign-headless-dispatch-mode-intent`). Per captain: **do NOT advance it to validation when implementation returns — leave it for the next session to validate.** It's the live-flake fix + the two-mode `-p` contract determination (captain-designed):
- Interactive → greet-and-stop; headless `-p` → drive to first gate/terminal + exit; headless + given-the-conn (prose) → resolve gates per `## Completion and Gates` → terminal. **Net contract −80 words (−25%)** — a reduction (captain pushed step 9 to a terse pointer form).
- Implementation is **spike-first**: prove `-p` deterministically drives on live (no greet-stop coin-flip) BEFORE the prose edit — that repeated-run determinism is AC-1's load-bearing proof. Also raises the 1m dispatch-close quiet budget and retargets `TestLiveEnsignCycle` off the `isTeamCreate` coin.

## The live-e2e flake — root-caused this session (→ 7e fixes it)

Every PR's live gate flaked; each dismissed as unrelated-to-change (host-isolated + a test the change doesn't touch + other hosts green). **Two harness/contract mismatches, both in `TestLiveEnsignCycle`/the live-cycle harness:** (1) greet-vs-drive ambiguity — the contract is *silent on `-p`* so the model coin-flips ("FO exited before TeamCreate"); (2) the 1m no-progress quiet budget is too tight for live turns ("dispatch close did not close within 1m0s"). 7e fixes both.

## Process learnings (the expensive ones)

1. **A registered live scenario ≠ a run one.** 58 + xf both shipped `//go:build live` tests that broke on their *first* CI run (58: empty-dir `gitInit` fixture; xf: unregistered). A new `TestLive*` must be (a) in `runtime-live-e2e.yml`'s `-run` AND (b) run green once locally before merge. (Backlog: a contractlint guard that every `TestLive*` is in `-run`.)
2. **Run live tests locally** — the validator can only *compile* them (no auth). Caught 58's fixture in 0.05s; confirmed xf residency PASS (356s).
3. **Detached audits earned their keep** — caught 5e's dangling runtime:13 ref, ec's hollow guard, 8e's material hole. Mandatory on high-stakes validation.
4. **CI approval gotchas:** the **Runtime Live E2E is a SEPARATE workflow run** from offline/install — find it (`gh run list --commit {sha}`), approve *its* `pending_deployments` directly, verify jobs flip to `in_progress`. The background auto-approve watcher silently failed once (#374 sat 32min).

## Captain direction — "too many steps" (the next theme)

Captain flagged the procedure is heavy. The distinction:
- **Discipline (keep):** detached audit, run-live-locally, registered scenarios, AC-on-independent-oracle — these caught the real bugs.
- **Ceremony (automate):** the ~9-step merge ceremony (run ~14× this session), the CI-approval dance, path-scoped commits. Targets: **`p2`** (pr-complete binary command — collapse the merge ceremony into one call), **`vc`** (reconcile --act), a `spacedock ci approve <pr>` helper. **These are the highest-leverage step-count cuts on the board.**
- Next contract direction (captain): **"condense the contract further"** after 7e lands.

## Sprint state
- **Active:** 7e (implementation — leave for next session).
- **Done/archived:** the 14 above.
- **v0.20.3 tag:** NOT yet cut. The 13-then-14 are antipattern-clean (pre-cut audit + per-task re-audits all clean). Captain holding the tag.

## Backlog seeds spun
- **p2 / vc / ci-approve helper** — collapse the merge ceremony + reconcile + CI approval into binary commands (the "too many steps" fix).
- **live-scenario-registration guard** — contractlint: every `//go:build live TestLive*` in `-run`.
- **live_test.go front-door de-tee** — ec's scoped-out follow-on (still dumps jsonl).
- **contract-condense (further)** — captain's next direction after 7e.
- **rzp** — dev-README+templates slim (the README slim itself FO-direct, now on main as `48edae4c`); **the Pages-serving repo-settings toggle** (operator, captain-owned).
