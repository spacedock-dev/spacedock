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
sprint: 0199-pre-flip-mechanics
group: survey
sprint-readiness: ready
---

One coherent correctness pass over the survey skill (`skills/survey/SKILL.md` + `references/queries.sql`), on a single corrected agentsview model. Consolidates the three coupled survey members (69 codex-cwd, 1p scaffold-fact, 4t sandbox-probe) plus the preflight's B2 model fix — they all edit one file, share the agentsview model, and bottom out on one live drive, so they ship as one task.

## Deliverables

1. **agentsview model fix (preflight B2 — the new piece).** 69's spike proved (and this ideation re-confirmed against real data) that agentsview v0.32.1 keys `project` by GIT-ROOT basename (this repo's root, every `.worktrees/*`, and `.spacedock-state` all key to ONE project), contradicting `SKILL.md:64` + the `queries.sql` `scoping`/`folded_keys` rationale (which assume per-cwd-basename divergence — the model `xn` shipped in 0.19.7). Ideation resolves the two open questions: `folded_keys` is structurally dead (always 1) and is DROPPED; the cwd-prefix-union scoping STAYS but its rationale inverts (it excludes same-basename sibling repos, not merges divergent keys). See "The corrected agentsview model" below.
2. **codex-cwd presence + hint (absorbs 69).** Codex sessions land `cwd=''` (agentsview doesn't persist Codex cwd) → the cwd-scoped query misses them. A separate flagged codex-presence count + report hint (NOT a silent project union — 69's collision evidence: workspace→6 roots, spacedock→3 roots). See 69's ideation report.
3. **SCAFFOLD state-the-fact (absorbs 1p).** Drop the recovered-vs-installed taxonomy (`SKILL.md` §3); state the observed scaffold fact, keeping the "recovered from behavior, not files" fact.
4. **sandbox install-detect probe swap (absorbs 4t).** Replace `command -v agentsview` (`SKILL.md:27`, a stat/access builtin Seatbelt denies) with an `agentsview --version` exit-code probe (an execve Seatbelt allows). See 4t's ideation report (root-caused).

## Absorbed entities

Supersedes `survey-codex-cwd-workaround` (69), `survey-scaffold-state-the-fact` (1p), `survey-agentsview-detect-under-sandbox` (4t) — each deferred with `superseded-by` = this. Their ideation reports (designs + spikes) are the inputs for deliverables 2/3/4; this task's ideation adds deliverable 1 (the §64 model fix) and weaves all four into one coherent SKILL.md/queries plan.

## Proof

Survey skill (`skills/survey`) — built by a worker under test; proof = the query-smoke (`queries.sql` executed against the production-shaped fixture) + a LIVE DRIVE exercising all observable changes (codex-presence hint, the corrected scoping with `folded_keys` dropped, the reshaped SCAFFOLD, the sandbox probe non-prompt) in ONE pass. Not a high-stakes surface → normal validation. Per the survey discipline, a grep over SKILL.md never satisfies a behavioral AC.

## The corrected agentsview model (deliverable 1 — the spine the other three hang off)

**The fact.** agentsview v0.32.1 derives `sessions.project` from the **git-root basename** (normalized: non-alphanumerics → `_`), NOT the cwd basename. Every checkout of one repo — the root, every worktree (`.worktrees/*` and `.claude/worktrees/*`), and the `docs/dev/.spacedock-state` split-root checkout — keys to the SAME `project`. Confirmed against real data (see Spike), reusing 69's git-root-basename finding plus a confirming probe in THIS ideation: 29 distinct spacedock-v1 cwds → ONE project key `spacedock_v1`; no spacedock-v1 subpath ever keyed to anything else; no single cwd ever mapped to >1 project.

**What this breaks.** Three places teach the disproven cwd-basename-divergence model:
- `SKILL.md:64` — "agentsview keys each session by the basename of its working directory, so a subdir checkout, a worktree, or the split-root state dir each get a DIFFERENT `project` key." FALSE: they get the SAME key.
- `queries.sql:18-21` (the `:repo_root` rationale) — "a subdir checkout and a worktree-style path coalesce to ONE repo identity even though agentsview keyed them under divergent basename-derived `project` columns." The coalescing is real; the *reason given* (divergent keys to merge) is false.
- `queries.sql:30-34` (the `scoping` comment) — "regardless of the divergent basename-derived `project` agentsview assigned a subdir/worktree checkout … the prefix union recovers the subdir + worktree sessions the basename key drops."

**What is STILL correct (do NOT change the mechanism).** The cwd-prefix-union (`cwd = :repo_root OR cwd LIKE :repo_root || '/%'`) stays. Under the corrected model its job inverts but is still load-bearing: `project` is git-root-*basename*, so it COLLIDES across unrelated repos that share a basename (69's real evidence: `spacedock` → 3 different roots, `workspace` → 6). A `project = :repo_project` scope would therefore fold a *sibling* repo's sessions into this one. The cwd-prefix is the ONLY column that distinguishes "this repo's worktrees/subdirs" (all under the absolute repo-root prefix) from "a same-basename sibling elsewhere on disk." So the scoping is right; only its *rationale prose* is wrong. This same collision is exactly why deliverable 2 (Codex, blank cwd) must be a flagged presence count, not a project union.

**`folded_keys` is now structurally dead.** `COUNT(DISTINCT project)` over a cwd-prefix-scoped Claude set is **always 1** for a git repo (every in-prefix cwd shares the git-root key) and trivially ≤1 for the non-git fallback (`REPO_ROOT=$(pwd)`, one cwd). It can never be >1, so the `folded_keys > 1` branch (`SKILL.md:96`, report line `SKILL.md:143`) is unreachable. **Decision: drop `folded_keys` from the `scoping` SELECT and remove the `{if folded_keys>1: coalesced from {folded_keys} keys}` report line + the §96 prose that interprets it.** Keep `blank_cwd` (still meaningful for Claude — 0.5% of Claude sessions have blank cwd per 69's spike — and the honest-accounting hook §98/§208 still wants it). This is a real shape change the query-smoke pins (AC-1).

## Deliverable details (exact before/after)

**D1 — model fix.** Rewrite `SKILL.md:64` to the git-root-basename fact: agentsview keys every checkout of one repo (root, worktrees, split-root state dir) to ONE `project` = the git-root basename, so the risk is the INVERSE — that key COLLIDES with same-basename sibling repos; resolve the repo root and scope by the absolute cwd-prefix so siblings stay out. Rewrite the `queries.sql` `:repo_root` rationale (`:18-21`) and the `scoping` comment (`:30-34`) the same way: the prefix excludes same-basename siblings, it does not merge divergent keys. Drop `folded_keys` (above). Net: prose tells the true story; mechanism unchanged.

**D2 — codex-presence count + hint (absorbs 69).** Add a labeled `-- name: codex-presence` query to `queries.sql`, bound to a new `:repo_project` param, returning `COUNT(*)` of `agent='codex'` sessions whose `project = :repo_project` and `SUM(cwd='')` among them. SKILL.md step 2 derives `:repo_project` from `REPO_ROOT`'s basename with the `-`/non-alphanumeric → `_` normalization agentsview applies (e.g. `spacedock-v1` → `spacedock_v1`), runs the query, and step 4 renders a hint line ONLY when the count > 0: "N Codex sessions match this repo by project NAME only (agentsview does not record Codex cwd) — may include a same-named sibling repo; the Claude-scoped body above does not cover them." NOT a union into the Claude `scoping`/scaffold/work-by-area sets (collision risk above). Codex decision/scaffold/work-by-area surfacing stays the documented deferred follow-up (`SKILL.md:13,96,209`).

**D3 — SCAFFOLD state-the-fact (absorbs 1p).** Drop the recovered-vs-installed-vs-active 3-bucket *taxonomy* (`SKILL.md` §3 step, and the §3-driven SCAFFOLD report block). State the observed fact directly: per scaffold family, its invocation count and whether it is checked in on disk. The "recovered from behavior, not files" phrasing IS the fact to state (1p's desired-output example), not taxonomy overhead. The file-probe + behavioral-tally JOIN stays (both signals are still read); only the bucket-label classification framing is removed in favor of the plain factual statement. Desired report shape (1p's example):
> `superpowers was recovered from behavior, not files: 186 skill invocations, but no checked-in .claude/skills/superpowers. Other recovered one-offs: plan-writing, using-git-worktrees, systematic-debug, simplify, debugging.`

**D4 — sandbox install-detect probe swap (absorbs 4t).** `SKILL.md:27`:
```
- if ! command -v agentsview >/dev/null; then echo "AGENTSVIEW MISSING"; fi
+ if ! agentsview --version >/dev/null 2>&1; then echo "AGENTSVIEW MISSING"; fi
```
`command -v` is an FS-access PATH-walk (`access()`/`stat()`) Seatbelt can deny while allowing `execve`; `agentsview --version` is an `execve` that survives whatever the survey's real through-the-binary reads survive. `2>&1` suppresses both the present banner and the absent "not found" so only the `AGENTSVIEW MISSING` sentinel reaches the agent. Contract preserved (silent ⇒ present; sentinel ⇒ absent). Install fallback (`SKILL.md:30`) unchanged.

## Step-4 report-fence composition (where 69 and 1p overlap)

Both touch the step-4 fence (`SKILL.md:141-169`). Resolved together so the build worker has ONE target shape:
- The PROJECT header line loses its `{if folded_keys>1: coalesced from {folded_keys} keys}` clause (D1) and KEEPS the `{if blank_cwd>0: · {blank_cwd} uncaptured-cwd sessions}` clause.
- A NEW conditional hint line under PROJECT renders only when `codex-presence > 0` (D2).
- The SCAFFOLD block becomes the plain state-the-fact statement (D3), not the bucket list.

## Spike: no new spike needed

**Mechanism de-risked.** The riskiest unverified path (D1, the model the other three hang off) was exercised in THIS ideation against real agentsview v0.32.1, driven through the binary (raw `~/.agentsview` reads are TCC-denied, per `SKILL.md:21`): synced THIS repo's Claude sessions across root + a worktree + the `.spacedock-state` checkout into a readable `AGENTSVIEW_DATA_DIR`, then queried. Result: 29 distinct spacedock-v1 cwds → ONE `project='spacedock_v1'`; `folded_keys=1` for the real scoping query; no in-prefix cwd ever keyed to another project; cwd→project is a function. D2's premise (Codex `cwd=''`, populated git-root-basename `project`, cross-repo basename collision) and D4's premise (FS-access vs execve asymmetry under Seatbelt) were each de-risked by 69's and 4t's spikes respectively (recorded in their now-superseded reports) and the D4 probe contract was re-confirmed here (present → silent exit 0; absent → `AGENTSVIEW MISSING`, underlying exit 127). **No spike needed: D1 rests on the git-root-basename finding confirmed here (and originally by 69); D2 on 69's spike; D3 composes already-proven query reads; D4 on 4t's Seatbelt-asymmetry spike.** The build's first tests are the fixture rows AC-1/AC-2 seed from these shapes.

## Acceptance criteria

- **AC-1 — corrected `scoping` shape: `folded_keys` dropped, `blank_cwd` kept, cwd-prefix scope unchanged.** `queries.sql`'s `scoping` query returns `sessions | blank_cwd | span` (no `folded_keys` column); the cwd-prefix-union (`cwd = :repo_root OR cwd LIKE :repo_root || '/%'`) is unchanged. Verified by: the query-smoke (AC-2) asserts the corrected column shape and that the prefix still coalesces the in-repo cwds while excluding the out-of-prefix and blank-cwd sessions — fails if `folded_keys` survives in the SELECT, if the column count is wrong, or if the prefix scope regresses. Expected values come from the fixture rows (independent of skill prose).

- **AC-2 — query-smoke over a production-accurate fixture proves Codex-by-`project` is counted, blank-cwd Codex rows are surfaced, and no Codex row is silently folded into the Claude scope.** `skills/integration/testdata/survey/fixture-sessions.sql` is corrected to the git-root-basename model (the in-repo Claude sessions all carry ONE `project` matching the repo basename, NOT the disproven divergent keys `_spacedock_state`/`feature_x`) and gains `agent='codex'` rows with the production shape (`project` set, `cwd = ''`): at least one matching `:repo_project` (in-repo Codex) and one with a same-`project`-different-implied-root sibling shape. `survey_queries_test.go` asserts: (a) the new `codex-presence` query returns the matching Codex count with `blank_cwd > 0`; (b) `scoping.sessions` is UNCHANGED by the added Codex rows (Codex stays out of the Claude scope — no union); (c) the corrected scoping column shape (AC-1). Verified by: `go test ./skills/integration/ -run TestSurveyQuerySmoke` passes. Expected values come from the fixture rows.

- **AC-3 — the sandbox install-detect probe uses invocation exit code, and emits only the sentinel.** `SKILL.md` step 1 probes presence with `agentsview --version >/dev/null 2>&1` (exit-code semantics), never `command -v`/`which`/`test -x`/`stat`/PATH-walk. A deterministic probe-behavior test runs the exact step-1 one-liner twice — once with a stub `agentsview` on a synthesized PATH (expect empty output, exit 0), once with the name absent from PATH (expect sole line `AGENTSVIEW MISSING`) — asserting captured stdout+stderr against those two independent fixture conditions. Verified by: that probe-behavior test passes (the gate-able half); the behavioral no-false-negative claim is AC-5's sandboxed live drive. Fails if the probe reverts to an FS-access form or leaks `--version` output.

- **AC-4 — the live survey report renders the corrected observable surface.** Driving `/spacedock:survey` end-to-end on a repo whose survey DB carries project-matched blank-cwd Codex sessions produces a report that: (a) shows the Codex-presence hint line stating the count and that the match is by project name only (cwd unrecorded); (b) shows the SCAFFOLD section as the state-the-fact statement (family + invocation count + checked-in presence/absence) with NO recovered-vs-installed-vs-active taxonomy labels; (c) shows the PROJECT line WITHOUT a "coalesced from N keys" clause (folded_keys gone). A second drive against a Codex-free DB omits the hint line. Verified by: the live-drive rendered report (the survey's proof bar); a grep over SKILL.md does NOT satisfy this — only the rendered output does.

- **AC-5 — under sandbox with agentsview installed, the survey detects it present and does NOT prompt to install.** Driving `/spacedock:survey` under the sandbox where the pre-fix `command -v` false-negatives, the run proceeds past step 1 to the sync/scan without emitting the `brew install --cask agentsview` consent prompt. Verified by: a sandboxed live drive observing the survey continue without the install prompt, with the pre-fix `command -v` probe confirmed to false-negative in the same sandbox as the negative control (establishing the drive exercises the failing path, not a vacuous pass).

## Test plan

- **Query-smoke (AC-1, AC-2) — fixture, cheap (~seconds), reuses the existing sqlite3-driven harness.** Correct `fixture-sessions.sql` to the git-root-basename model: the three in-repo Claude sessions all key to one `project` matching the repo basename (replacing the divergent `_spacedock_state`/`feature_x` keys), and drop the fixture comment lines that assert divergent keys. Add blank-cwd Codex rows (`project` set, `cwd=''`): one in-repo (`project = :repo_project`) and one same-basename sibling shape. Update `survey_queries_test.go`: the `scoping` assertions move from `sessions|folded_keys|blank_cwd|span` (4 fields, folded_keys=3) to the corrected `sessions|blank_cwd|span` (3 fields) and assert the prefix still counts 3 in-repo sessions; add the `codex-presence` assertions (matched count, blank_cwd>0) and the `scoping.sessions` unchanged-by-Codex assertion. Cost: LOW.
- **Probe-behavior test (AC-3) — deterministic, the gate-able half.** Extract/run the exact step-1 probe line against two synthesized PATH conditions (stub binary present; name absent), assert captured output. If the survey harness is shell-fixture-based, add it there; otherwise a focused shell/Go test under the survey integration testdata. Expected values (`""`, `AGENTSVIEW MISSING`) come from the fixture conditions. Cost: LOW.
- **Live drive (AC-4, AC-5) — minutes; the ONLY proof for the rendered report and the sandbox non-prompt.** ONE end-to-end pass exercises all observable changes at once: seed a survey DB (via the spike's `agentsview session sync` / scoped `agentsview sync` into a readable `AGENTSVIEW_DATA_DIR`) with project-matched blank-cwd Codex sessions, drive `/spacedock:survey` under sandbox, observe (a) the Codex hint, (b) the state-the-fact SCAFFOLD, (c) no folded-keys line, (d) no install prompt. Then a second short drive against a Codex-free DB to confirm the hint is absent (AC-4 negative), and the pre-fix `command -v` confirmed false-negative in the same sandbox (AC-5 negative control). Needs the agentsview binary (present, v0.32.1) and the sandboxed harness.
- **Estimated total cost/complexity: LOW–MEDIUM.** Prose edits across `SKILL.md` (§1 probe, §2 scoping/model, §3 SCAFFOLD, §4 fence) + one labeled query + a `folded_keys` drop + fixture correction + Codex fixture rows + test assertions + one probe-behavior test + one combined live drive. No agentsview/binary/Go-runtime changes beyond the test files (agentsview changes are explicitly out of scope).

## Out of scope

- Fixing agentsview to persist Codex `cwd` (upstream `kenn-io/agentsview`).
- Surfacing Codex decision / scaffold / work-by-area signals (broadening the `agent='claude'` filter across the body) — the documented deferred follow-up; this pass adds only the codex-presence count + hint.
- The install-instruction path (`SKILL.md:30`) — unchanged; D4 fixes detection only.

## Stage Report: ideation

- DONE: Absorb + integrate the three superseded designs (69 codex-cwd, 1p scaffold-fact, 4t sandbox-probe) into ONE coherent skills/survey/SKILL.md + references/queries.sql plan — including the step-4 report-template fence where 69 and 1p overlap.
  Read all three ideation reports; wove D2 (codex-presence count + hint), D3 (SCAFFOLD state-the-fact), D4 (`agentsview --version` probe) into "Deliverable details" + the "Step-4 report-fence composition" section resolving the 69/1p overlap (folded_keys clause dropped, codex hint added, SCAFFOLD reshaped) into ONE target fence shape.
- DONE: Firm the B2 model fix (the new piece): correct SKILL.md:64 + the queries.sql scoping/folded_keys rationale to the git-root-basename reality; decide whether folded_keys stays meaningful and whether the cwd-prefix-union is still right.
  Recorded in "The corrected agentsview model" + D1. Decisions: `folded_keys` is structurally always 1 → DROP it (column + report line + §96 prose); the cwd-prefix-union STAYS but its rationale inverts (excludes same-basename sibling repos rather than merging divergent keys, since the keys are NOT divergent). Reused 69's git-root-basename finding AND re-confirmed with a fresh probe (Spike section) — 29 distinct spacedock-v1 cwds → ONE `project='spacedock_v1'`, real scoping query `folded_keys=1`.
- DONE: Produce the unified build-ready ACs + test plan: query-smoke (queries.sql vs production fixture incl. cwd='' Codex rows) + the deterministic sandbox-probe test + ONE live-drive pass exercising all observable changes. No tautological SKILL.md grep as a behavioral AC.
  Five ACs (AC-1 corrected scoping shape, AC-2 query-smoke incl. blank-cwd Codex rows + no-union, AC-3 deterministic probe-behavior test, AC-4 live-drive rendered report, AC-5 sandboxed non-prompt with `command -v` negative control) + a 4-bullet test plan. Each AC's proof is a test/command/rendered-output outside the task body; the two text-half changes (probe wording, model prose) are explicitly NOT standalone behavioral ACs — the behavioral half is the live drive.

### Summary

Consolidated the three coupled survey members plus the preflight's B2 model fix into one coherent ideation on the corrected git-root-basename agentsview model. The B2 fix is the spine: re-confirmed against real agentsview v0.32.1 data (29 spacedock-v1 cwds → one `project` key) that `folded_keys` is structurally always 1 (DROP it) while the cwd-prefix scoping is still load-bearing for the inverse reason (it excludes same-basename sibling repos, which is also why Codex's blank-cwd sessions get a flagged presence count, not a silent project union). No new spike needed — D1 re-confirmed here, D2/D4 rest on 69's and 4t's recorded spikes, D3 composes proven reads. ACs pin the corrected query shape and Codex no-union via query-smoke, the probe via a deterministic present/absent test, and the rendered report + sandbox non-prompt via ONE combined live drive.
