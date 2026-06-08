---
id: vh5rzxexn9wc7dedex7cepzy
title: Survey skill correctness pass — agentsview git-root model fix + codex-cwd + scaffold-fact + sandbox probe (consolidates 69/1p/4t)
status: ideation
source: "captain (2026-06-08) — preflight B2 + consolidation directive. The survey members 69/1p/4t all edit skills/survey/SKILL.md and share one agentsview model; 69's spike disproved the cwd-basename keying the survey teaches. Consolidate into ONE coherent survey-skill pass on the corrected git-root-basename model."
score: "0.3"
started: 2026-06-08T17:02:54Z
completed:
verdict:
worktree:
issue:
sprint: 0198-pre-flip-hardening
group: survey
sprint-readiness: ready
---

One coherent correctness pass over the survey skill (`skills/survey/SKILL.md` + `references/queries.sql`), on a single corrected agentsview model. Consolidates the three coupled survey members (69 codex-cwd, 1p scaffold-fact, 4t sandbox-probe) plus the preflight's B2 model fix — they all edit one file, share the agentsview model, and bottom out on one live drive, so they ship as one task.

## Deliverables

1. **agentsview model fix (preflight B2 — the new piece).** 69's spike proved agentsview v0.32.1 keys `project` by GIT-ROOT basename (this repo's root, every `.worktrees/*`, and `.spacedock-state` all key to ONE project), contradicting `SKILL.md:64` + the `queries.sql` `scoping`/`folded_keys` rationale (which assume per-cwd-basename divergence — the model `xn` shipped in 0.19.7). Correct §64 + the scoping/folded_keys prose to the git-root-basename reality; reconsider whether `folded_keys` is still meaningful and whether the cwd-prefix-union is the right scoping mechanism under the corrected model.
2. **codex-cwd presence + hint (absorbs 69).** Codex sessions land `cwd=''` (agentsview doesn't persist Codex cwd) → the cwd-scoped query misses them. A separate flagged codex-presence count + report hint (NOT a silent project union — 69's collision evidence: workspace→6 roots, spacedock→3 roots). See 69's ideation report.
3. **SCAFFOLD state-the-fact (absorbs 1p).** Drop the recovered-vs-installed taxonomy (`SKILL.md` §3); state the observed scaffold fact, keeping the "recovered from behavior, not files" fact.
4. **sandbox install-detect probe swap (absorbs 4t).** Replace `command -v agentsview` (`SKILL.md:27`, a stat/access builtin Seatbelt denies) with an `agentsview --version` exit-code probe (an execve Seatbelt allows). See 4t's ideation report (root-caused).

## Absorbed entities

Supersedes `survey-codex-cwd-workaround` (69), `survey-scaffold-state-the-fact` (1p), `survey-agentsview-detect-under-sandbox` (4t) — each deferred with `superseded-by` = this. Their ideation reports (designs + spikes) are the inputs for deliverables 2/3/4; this task's ideation adds deliverable 1 (the §64 model fix) and weaves all four into one coherent SKILL.md/queries plan.

## Proof

Survey skill (`skills/survey`) — built by a worker under test; proof = the query-smoke (`queries.sql` executed against the production-shaped fixture) + a LIVE DRIVE exercising all observable changes (codex-presence hint, the corrected scoping/folded_keys, the reshaped SCAFFOLD, the sandbox probe non-prompt) in ONE pass. Not a high-stakes surface → normal validation. Per the survey discipline, a grep over SKILL.md never satisfies a behavioral AC.

## Acceptance criteria (ideation firms — absorbing the 3 reports' ACs + the §64 fix)

Ideation produces the unified AC set: the §64/folded_keys correction (verified against the corrected model + the live drive), codex-presence (query-smoke fixture with `cwd=''` Codex rows), the SCAFFOLD reshape (live drive), and the sandbox probe (deterministic probe-behavior test + sandboxed live drive with the pre-fix `command -v` as negative control).
